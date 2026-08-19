package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// expiredReadDeadline is the deterministic timeout stimulus for the read-
// deadline tests: any NON-POSITIVE duration makes context.WithTimeout take
// context.WithDeadline's `dur <= 0` branch, which cancels the context
// SYNCHRONOUSLY at construction — no timer, no goroutine, no race — so the
// store read is refused at connection acquisition and the 503/Timeout branch
// runs on every route, every run. A small POSITIVE duration (the previous
// `1 * time.Nanosecond`) is a FUTURE deadline: it arms a time.AfterFunc, and a
// fast read can complete before the timer goroutine runs, answering 200 with a
// real body (measured at base: ~0.65–0.8% of runs). Do not "shrink" this back
// to a positive value; TestExpiredReadDeadlineExpiresAtConstruction reds on
// the sign, and the design doc for w-daemon-timeout-test-flake holds the
// measurements.
const expiredReadDeadline = -1 * time.Nanosecond

// ---------------------------------------------------------------------------
// The route table under test
// ---------------------------------------------------------------------------

// readRoute is one of the SIX /v1 GET routes that reach the store. Six routes,
// FIVE distinct getters: GetLogEntry serves both /v1/log/{index} and the
// bounded loop of /v1/log. A test that drives fewer than six routes cannot see
// a deadline that was installed in four handlers and forgotten in two, which is
// exactly the drift shape this item is about.
type readRoute struct {
	name   string
	target string
	getter string
}

// seedReadRoutes seeds d's store with a genesis world plus one commit and
// returns the six route targets, every one of which answers 200 against the
// unmutated daemon (asserted by the normal-deadline subtest below — without
// that arm, a 503-on-everything daemon would pass the timeout assertions).
func seedReadRoutes(t *testing.T, d *Daemon, label string) []readRoute {
	t.Helper()
	genesis := seedGenesisEmbedded(t, d, label)
	commit := testCommit(genesis, 1, label)
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("seed Commit: %v", err)
	}
	return []readRoute{
		{"head", "/v1/head", "SelectedHead"},
		{"world", "/v1/worlds/" + commit.NextWorld.Ref.String(), "GetWorld"},
		{"object", "/v1/objects/" + commit.Objects[0].Hash.String(), "GetObject"},
		{"log entry", "/v1/log/1", "GetLogEntry"},
		{"log range", "/v1/log?from=0&limit=5", "GetLogEntry"},
		{"registry", "/v1/registry/world/epoch-registry/v1", "GetRegistryHead"},
	}
}

// ---------------------------------------------------------------------------
// blockingStore — the hang-shaped stimulus, WRAPPING the real store
// ---------------------------------------------------------------------------

// errStoreInterrupted is what every blocked getter returns when the caller's
// context is cancelled.
//
// It deliberately does NOT wrap context.DeadlineExceeded, because the real
// driver's interrupt path does not either: a read cancelled mid-flight surfaces
// as SQLITE_INTERRUPT. That is precisely why the handler classifies on
// ctx.Err() FIRST and only then on errors.Is — and it is what makes MU3
// (`if ctx.Err() != nil` -> `if false && ctx.Err() != nil`) a real kill instead
// of a survivor. A fake that returned ctx.Err() here would let the errors.Is
// arm classify the mutant correctly and MU3 would survive.
var errStoreInterrupted = errors.New("store: query interrupted (SQLITE_INTERRUPT)")

// errStoreEscaped is returned when the WATCHDOG released a parked getter. It is
// a distinct sentinel so a released-by-watchdog run can never be misread as a
// pass: it is not a context error, so the handler classifies it as Internal and
// every 503 assertion in this file reds on it.
var errStoreEscaped = errors.New("blockingStore: released by the test watchdog — this arm never passes")

// blockingStore embeds the real *store.Store and overrides ALL FIVE getters to
// block. All five, not one: the six routes reach the store through five
// different getters, so a one-getter fake would let the other routes fall
// through to the embedded real store, answer 200 in microseconds, and red the
// 503 assertion for a reason that has nothing to do with the mutation.
//
// It WRAPS rather than REPLACES (the iteration-80 vacuity trap): only the store
// BODIES are substituted, so every line of handler code under test — readCtx,
// the defer cancel(), timedOut, writeReadTimeout — is the production path.
type blockingStore struct {
	*store.Store

	escape    chan struct{}
	escapeOne sync.Once

	enteredOnce sync.Once
	entered     chan struct{}

	mu       sync.Mutex
	released []string
}

func newBlockingStore(s *store.Store) *blockingStore {
	return &blockingStore{Store: s, escape: make(chan struct{}), entered: make(chan struct{})}
}

// release is the watchdog's red path: it unparks every blocked getter so the
// handler returns, ServeHTTP completes and cleanup runs against a quiescent
// pool. A bare t.Error unblocks nothing, and with the store's single connection
// held by a parked getter a deferred Close() behind it turns a clean red into a
// suite-wide hang.
func (b *blockingStore) release() { b.escapeOne.Do(func() { close(b.escape) }) }

func (b *blockingStore) block(ctx context.Context) error {
	b.enteredOnce.Do(func() { close(b.entered) })
	select {
	case <-ctx.Done():
		b.record("ctx-done")
		return errStoreInterrupted
	case <-b.escape:
		b.record("escape")
		return errStoreEscaped
	}
}

func (b *blockingStore) record(how string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.released = append(b.released, how)
}

func (b *blockingStore) releases() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.released...)
}

func (b *blockingStore) GetObject(ctx context.Context, _ hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, b.block(ctx)
}

func (b *blockingStore) GetWorld(ctx context.Context, _ hashref.HashRef) (store.World, bool, error) {
	return store.World{}, false, b.block(ctx)
}

func (b *blockingStore) GetLogEntry(ctx context.Context, _ int64) (store.LogEntry, bool, error) {
	return store.LogEntry{}, false, b.block(ctx)
}

func (b *blockingStore) GetRegistryHead(ctx context.Context, _ string) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, b.block(ctx)
}

func (b *blockingStore) SelectedHead(ctx context.Context) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, b.block(ctx)
}

// ---------------------------------------------------------------------------
// recordingStore — records the context each getter received, then delegates
// ---------------------------------------------------------------------------

// recordingStore also wraps the real store and overrides all five getters, but
// it BLOCKS NOTHING: it records the context it was handed and delegates to the
// embedded real store, so every route serves a normal 200. Its only observable
// is the recorded context, read AFTER ServeHTTP has returned.
type recordingStore struct {
	*store.Store

	mu  sync.Mutex
	got []context.Context
}

func (r *recordingStore) note(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.got = append(r.got, ctx)
}

func (r *recordingStore) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.got = nil
}

func (r *recordingStore) recorded() []context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]context.Context(nil), r.got...)
}

func (r *recordingStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	r.note(ctx)
	return r.Store.GetObject(ctx, ref)
}

func (r *recordingStore) GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error) {
	r.note(ctx)
	return r.Store.GetWorld(ctx, ref)
}

func (r *recordingStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	r.note(ctx)
	return r.Store.GetLogEntry(ctx, index)
}

func (r *recordingStore) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	r.note(ctx)
	return r.Store.GetRegistryHead(ctx, name)
}

func (r *recordingStore) SelectedHead(ctx context.Context) (hashref.HashRef, bool, error) {
	r.note(ctx)
	return r.Store.SelectedHead(ctx)
}

// ---------------------------------------------------------------------------
// AC3 layer 2 + layer 3 — the deadline is installed and renders as 503/Timeout
// ---------------------------------------------------------------------------

// TestDaemonReadDeadline proves the read deadline is REAL on all six
// store-reading GET routes, in three arms that fail for different reasons:
//
//   - normal-deadline-answers-200 is the non-vacuity control. Without it a
//     daemon that answered 503 to everything would pass the two timeout arms.
//   - real-store-expired-deadline is FAKE-FREE: a real seeded store and a 1 ns
//     deadline. Remove the deadline from readCtx (MU1) and the store answers
//     normally — 200 with a real body — which no timeout branch can write.
//   - blocking-store supplies the stimulus a real store cannot produce on
//     demand: a getter that is still waiting when the deadline fires. It is the
//     only arm that can see the CLASSIFIER (MU3), because its error is
//     interrupt-shaped and does not wrap context.DeadlineExceeded.
func TestDaemonReadDeadline(t *testing.T) {
	t.Run("normal-deadline-answers-200", func(t *testing.T) {
		d := newHandlerDaemon(t)
		routes := seedReadRoutes(t, d, "normal-deadline")
		if d.readDeadline != readDeadline {
			t.Fatalf("d.readDeadline = %s, want the shipped default %s", d.readDeadline, readDeadline)
		}
		for _, route := range routes {
			rec := requestRecorder(t, d, http.MethodGet, route.target, nil)
			if rec.Code != http.StatusOK {
				t.Errorf("%s (%s): status = %d, want 200 under the shipped %s deadline; body=%s",
					route.name, route.target, rec.Code, readDeadline, rec.Body)
			}
		}
	})

	t.Run("real-store-expired-deadline", func(t *testing.T) {
		d := newHandlerDaemon(t)
		routes := seedReadRoutes(t, d, "expired-deadline")
		// Shrinking the FIELD is the only way to exercise the shipped timeout
		// branch without a ten-second test — which is why New wires a field
		// rather than referencing the constant directly.
		d.readDeadline = expiredReadDeadline

		for _, route := range routes {
			rec := requestRecorder(t, d, http.MethodGet, route.target, nil)
			body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
			if !strings.Contains(body.Error.Message, "read deadline") {
				t.Errorf("%s: timeout message = %q, want it to name the deadline",
					route.name, body.Error.Message)
			}
		}
	})

	t.Run("blocking-store", func(t *testing.T) {
		d := newHandlerDaemon(t)
		routes := seedReadRoutes(t, d, "blocking-store")
		blocking := newBlockingStore(d.store)
		d.reads = blocking
		d.readDeadline = 50 * time.Millisecond
		t.Cleanup(blocking.release)

		for _, route := range routes {
			route := route
			done := make(chan *httptest.ResponseRecorder, 1)
			go func() {
				req := httptest.NewRequest(http.MethodGet, route.target, nil)
				rec := httptest.NewRecorder()
				d.Handler().ServeHTTP(rec, req)
				done <- rec
			}()

			var rec *httptest.ResponseRecorder
			select {
			case rec = <-done:
			case <-time.After(2 * time.Second):
				// Red path: report, then RELEASE, then confirm the goroutine
				// actually exited before this subtest's cleanup runs.
				t.Errorf("%s (%s): handler still blocked after 2s — the %s deadline never fired",
					route.name, route.target, d.readDeadline)
				blocking.release()
				select {
				case <-done:
				case <-time.After(2 * time.Second):
					t.Fatalf("%s: handler did not exit even after the escape channel closed", route.name)
				}
				continue
			}
			assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
		}

		// Every release must have come from ctx-done. An "escape" release means
		// the watchdog fired, and the distinct sentinel guarantees that shows up
		// as a red above rather than as a quiet pass.
		for i, how := range blocking.releases() {
			if how != "ctx-done" {
				t.Errorf("release[%d] = %q, want %q — the getter was unparked by the watchdog, not the deadline",
					i, how, "ctx-done")
			}
		}
	})
}

// ---------------------------------------------------------------------------
// AC3/AC4' layer 3b — the context is the REQUEST's, not Background's
// ---------------------------------------------------------------------------

// TestDaemonReadDisconnect is the arm that discriminates r.Context() from
// context.Background() in readCtx, and it is the behavioural authority for
// AC4' (a grep can only see a spelling; this sees the propagation).
//
// The read deadline is left at its shipped 10 s default, so the ONLY thing that
// can unblock the parked getter inside the 2 s bound is cancellation flowing
// from the request. Under MU2 (r.Context() -> context.Background()) the derived
// context is orphaned: the unblock signal arrives at the 10 s deadline, long
// after the watchdog has red — and released.
func TestDaemonReadDisconnect(t *testing.T) {
	d := newHandlerDaemon(t)
	routes := seedReadRoutes(t, d, "disconnect")
	blocking := newBlockingStore(d.store)
	d.reads = blocking
	t.Cleanup(blocking.release)

	if d.readDeadline != readDeadline {
		t.Fatalf("d.readDeadline = %s, want the shipped default %s — this test's whole "+
			"discriminating power comes from the deadline being far outside the watchdog",
			d.readDeadline, readDeadline)
	}

	route := routes[1] // /v1/worlds/{ref} -> GetWorld
	reqCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		req := httptest.NewRequest(http.MethodGet, route.target, nil).WithContext(reqCtx)
		d.Handler().ServeHTTP(httptest.NewRecorder(), req)
	}()

	// Wait for the getter to actually be parked before disconnecting; cancelling
	// before the store call starts would prove nothing about propagation.
	select {
	case <-blocking.entered:
	case <-time.After(2 * time.Second):
		t.Fatalf("the blocking getter was never entered — the request never reached the store")
	}
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Errorf("the store call did not unblock within 2s of the client disconnecting — "+
			"readCtx is not derived from r.Context() (the %s deadline is the only other writer)",
			d.readDeadline)
		blocking.release()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("handler did not exit even after the escape channel closed")
		}
	}

	releases := blocking.releases()
	if len(releases) == 0 {
		t.Fatalf("no release was recorded — the getter never returned")
	}
	if releases[0] != "ctx-done" {
		t.Errorf("first release = %q, want %q — the disconnect did not reach the store's context",
			releases[0], "ctx-done")
	}
}

// ---------------------------------------------------------------------------
// AC3 layer 4 — every handler releases its CancelFunc on return
// ---------------------------------------------------------------------------

// TestReadCtxCancelledAfterHandler is the gate on the `defer cancel()` mandate.
// A context.WithTimeout whose CancelFunc is never called leaks its timer and its
// context subtree until the deadline fires — once per store-reading request.
//
// The observable is the context the getter itself received, read AFTER
// ServeHTTP returns: it must already report context.Canceled. With readDeadline
// at its 10 s default, nothing else can write that value inside test time, so
// deleting one route's `defer cancel()` (MU11) leaves that route's recorded
// Err() nil and reds here. The 200 assertion is the paired control: a cancel
// that fired EARLY would error the delegated store call and red it.
func TestReadCtxCancelledAfterHandler(t *testing.T) {
	d := newHandlerDaemon(t)
	routes := seedReadRoutes(t, d, "cancel-release")
	recording := &recordingStore{Store: d.store}
	d.reads = recording

	if d.readDeadline != readDeadline {
		t.Fatalf("d.readDeadline = %s, want the shipped default %s", d.readDeadline, readDeadline)
	}

	for _, route := range routes {
		recording.reset()
		rec := requestRecorder(t, d, http.MethodGet, route.target, nil)
		if rec.Code != http.StatusOK {
			t.Errorf("%s (%s): status = %d, want 200; body=%s", route.name, route.target, rec.Code, rec.Body)
			continue
		}
		got := recording.recorded()
		if len(got) == 0 {
			t.Errorf("%s: the route recorded no store call — the seam is not on this route's path", route.name)
			continue
		}
		for i, ctx := range got {
			if err := ctx.Err(); !errors.Is(err, context.Canceled) {
				t.Errorf("%s: recorded ctx[%d].Err() = %v, want context.Canceled — %s's "+
					"`defer cancel()` is missing, so its timer leaks until the %s deadline",
					route.name, i, err, route.getter, d.readDeadline)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// AC3/AC7 — the Go 503 and the frozen sketch cannot drift
// ---------------------------------------------------------------------------

// sketchHTTPStatusVectors parses the httpStatus arms out of the frozen,
// compiler-checked sketch (design_docs/sketches/worlddapi.ail) rather than
// restating them. Restating them would make this test a mirror of itself: the
// point of the idiom (TestIsLoopbackHostMirrorsSketchPredicate) is that the
// AILANG side is the authority and the Go side is checked against it.
//
// Rooted via runtime.Caller so the result does not depend on the test's working
// directory — `ai-check` and `go test` disagree about cwd often enough that a
// relative path here would be an instrument failure, not a result.
func sketchHTTPStatusVectors(t *testing.T) map[string]int {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller(0) failed — cannot root the sketch path")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	path := filepath.Join(repoRoot, "design_docs", "sketches", "worlddapi.ail")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the frozen sketch at %s: %v", path, err)
	}
	vectors := map[string]int{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "=>") {
			continue
		}
		arm, status, found := strings.Cut(line, "=>")
		if !found {
			continue
		}
		name, _, hasArgs := strings.Cut(strings.TrimSpace(arm), "(")
		if !hasArgs {
			continue
		}
		var code int
		if _, err := fmt.Sscanf(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(status), ",")), "%d", &code); err != nil {
			continue
		}
		vectors[strings.TrimSpace(name)] = code
	}
	// Non-vacuity: a parser that silently matched nothing would make every
	// assertion below pass. The sketch's own arms are the pin.
	for _, required := range []string{"HeadConflict", "NotFound", "BadRequest", "PayloadTooLarge", "Internal", "Timeout"} {
		if _, present := vectors[required]; !present {
			t.Fatalf("the sketch's httpStatus has no %s arm (parsed %d arms from %s) — "+
				"either the sketch drifted or this parser is broken; both are red",
				required, len(vectors), path)
		}
	}
	return vectors
}

// TestTimeoutStatusMirrorsSketch replays the sketch's own httpStatus vector for
// the Timeout arm against what the daemon ACTUALLY writes, so the frozen AILANG
// surface and the Go transport cannot drift apart.
//
// Two arms, and both are load-bearing:
//
//   - the sketch's Timeout arm must still say 503 (source-of-truth pin);
//   - a real timed-out read must answer exactly that status, with class
//     "Timeout". MU9 (503 -> 500 in the Go timeout branch) reds this second arm
//     while the sketch is untouched — which is the drift direction the gate can
//     actually enforce, since Leg 2 never runs the sketch's inline tests.
func TestTimeoutStatusMirrorsSketch(t *testing.T) {
	vectors := sketchHTTPStatusVectors(t)

	wantTimeout, ok := vectors["Timeout"]
	if !ok {
		t.Fatalf("the sketch declares no Timeout arm")
	}
	if wantTimeout != 503 {
		t.Fatalf("the sketch's Timeout arm = %d, want 503 (504 is a GATEWAY's statement; this daemon is the origin)",
			wantTimeout)
	}
	// The other arms are asserted too, so a sketch edit that "fixed" the Timeout
	// row by renumbering its neighbours cannot pass.
	for _, c := range []struct {
		arm  string
		want int
	}{
		{"HeadConflict", 409}, {"NotFound", 404}, {"BadRequest", 400},
		{"PayloadTooLarge", 413}, {"Internal", 500}, {"Timeout", 503},
	} {
		if got := vectors[c.arm]; got != c.want {
			t.Errorf("sketch httpStatus %s => %d, want %d", c.arm, got, c.want)
		}
	}

	d := newHandlerDaemon(t)
	routes := seedReadRoutes(t, d, "mirrors-sketch")
	d.readDeadline = expiredReadDeadline

	for _, route := range routes {
		rec := requestRecorder(t, d, http.MethodGet, route.target, nil)
		if rec.Code != wantTimeout {
			t.Errorf("%s (%s): a timed-out read answered %d, want the sketch's %d for Timeout; body=%s",
				route.name, route.target, rec.Code, wantTimeout, rec.Body)
			continue
		}
		var body APIError
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Errorf("%s: decode API error: %v; body=%s", route.name, err, rec.Body)
			continue
		}
		if body.Error.Class != "Timeout" {
			t.Errorf("%s: class = %q, want %q — the sketch's ApiError constructor name is the wire class",
				route.name, body.Error.Class, "Timeout")
		}
	}
}

// ---------------------------------------------------------------------------
// AC5 — a 500 sanitizes the wire and keeps the detail for the operator
// ---------------------------------------------------------------------------

// internalDetailSentinel is a string PRODUCTION CODE CANNOT PRODUCE.
//
// That property is the whole test, and it is the direct lesson of M2's MU3
// survival: there, the design's own prescribed fake returned ctx.Err(), which is
// exactly the value the surviving arm of a two-arm classifier needed, so the
// mutant passed on a value that came from the FIXTURE rather than from the code
// under test. A sanitize test seeded with any error text the daemon could have
// written itself has the same hole — "the body does not contain the detail"
// would be satisfiable by accident.
//
// So the seeded detail is a random token, verified absent from the entire
// repository before it was chosen (`grep -rc kQ7v .` -> no file matched, with a
// same-scope known-positive control that did match). If it appears in a response
// body, it can ONLY have travelled there from errSentinelInternal through the
// 500 branch; if it appears on the error log, likewise. Neither write can be
// faked into existence by anything else in the process.
const internalDetailSentinel = "kQ7v-store-detail-9f3c1d82"

// errSentinelInternal is shaped like the real leak this AC closes: store errors
// interpolate the display DSN path (`store: open %q`), which is host filesystem
// detail with no bearing on the caller's request.
//
// It is deliberately NOT a context error and does NOT wrap
// context.DeadlineExceeded, so timedOut() classifies it false and the INTERNAL
// branch runs — not the 503 branch. An arm that accidentally rendered 503 would
// be caught by assertErrorClass below, not silently pass.
var errSentinelInternal = errors.New(
	`store: open "/private/var/folders/` + internalDetailSentinel + `/world.db": disk I/O error`)

// failingStore wraps the real store and overrides ALL FIVE getters to fail with
// the sentinel error. It WRAPS rather than replaces (the iteration-80 vacuity
// trap): every line of the handler under test — readCtx, defer cancel, timedOut,
// writeInternalError — is the production path.
//
// All five, not one: the six read routes reach the store through five distinct
// getters, so a one-getter fake would let the other routes fall through to the
// embedded real store and answer 200, and assertErrorClass would red for a
// reason that has nothing to do with sanitization.
type failingStore struct {
	*store.Store
}

func (failingStore) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, errSentinelInternal
}

func (failingStore) GetWorld(context.Context, hashref.HashRef) (store.World, bool, error) {
	return store.World{}, false, errSentinelInternal
}

func (failingStore) GetLogEntry(context.Context, int64) (store.LogEntry, bool, error) {
	return store.LogEntry{}, false, errSentinelInternal
}

func (failingStore) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, errSentinelInternal
}

func (failingStore) SelectedHead(context.Context) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, errSentinelInternal
}

// TestInternalErrorsAreSanitized is AC5's persistent form.
//
// It asserts TWO WRITES SEPARATELY, and the separation is the design:
//
//	(a) the RESPONSE BODY must NOT carry the sentinel, and must carry
//	    internalErrorMessage and nothing else;
//	(b) the ERROR LOG must carry the sentinel VERBATIM, on one line, next to the
//	    route it came from.
//
// Neither assertion implies the other. A mutation that restores the raw error to
// the body (MU7) still writes the log line and so still satisfies (b) — it dies
// on (a). A mutation that drops the log write still sanitizes the body and so
// still satisfies (a) — it dies on (b). A single combined assertion would be
// killable by either mutation and would tell you nothing about which.
//
// The route enumeration matters too: it drives all SIX read routes plus
// POST /v1/commit, i.e. every 500 branch in the package, so a sanitizer
// installed in five handlers and forgotten in the sixth reds here.
func TestInternalErrorsAreSanitized(t *testing.T) {
	t.Run("error-log-defaults-to-stderr", func(t *testing.T) {
		// The nil-means-os.Stderr contract, asserted on the CONSTRUCTED field
		// rather than inferred. New resolves it once; a daemon that resolved
		// lazily (or not at all) would leave this nil and panic on the first
		// 500, which is precisely the moment nothing may panic.
		d := newHandlerDaemon(t)
		if d.errLog == nil {
			t.Fatalf("d.errLog is nil on a constructed daemon — a 500 would panic writing its detail line")
		}
		if d.errLog != os.Stderr {
			t.Errorf("d.errLog = %#v, want os.Stderr — Config.ErrorLog was nil, and nil means the process owner's stderr", d.errLog)
		}
	})

	t.Run("read-routes", func(t *testing.T) {
		d := newHandlerDaemon(t)
		routes := seedReadRoutes(t, d, "sanitize")

		var errLog bytes.Buffer
		d.errLog = &errLog
		d.reads = failingStore{Store: d.store}

		for _, route := range routes {
			errLog.Reset()
			rec := requestRecorder(t, d, http.MethodGet, route.target, nil)
			body := assertErrorClass(t, rec, http.StatusInternalServerError, "Internal")

			// (a) THE RESPONSE WRITE. Independent of (b).
			if strings.Contains(rec.Body.String(), internalDetailSentinel) {
				t.Errorf("%s (%s): the 500 body leaked the store's internal detail %q; body=%s",
					route.name, route.target, internalDetailSentinel, rec.Body)
			}
			if body.Error.Message != internalErrorMessage {
				t.Errorf("%s: 500 message = %q, want the fixed %q — a 500 body carries no host state",
					route.name, body.Error.Message, internalErrorMessage)
			}

			// (b) THE ERROR-LOG WRITE. Independent of (a).
			assertErrorLogLine(t, errLog.String(), route.name, http.MethodGet, route.target)
		}
	})

	t.Run("commit-route", func(t *testing.T) {
		// POST /v1/commit's 500 is the one internal branch NOT on the read seam
		// (it is store.Commit, the write path, which this item deliberately does
		// not put behind an interface). Closing the store is how the existing
		// suite already produces a real store failure on a live route, so this
		// arm exercises the shipped branch with a REAL error rather than a fake.
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "sanitize-commit")
		payload := encodeCommit(testCommit(genesis, 1, "sanitize-commit"))
		if err := d.store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}

		var errLog bytes.Buffer
		d.errLog = &errLog

		rec := requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(payload))
		body := assertErrorClass(t, rec, http.StatusInternalServerError, "Internal")
		if body.Error.Message != internalErrorMessage {
			t.Errorf("commit 500 message = %q, want the fixed %q; body=%s",
				body.Error.Message, internalErrorMessage, rec.Body)
		}
		// The detail must still reach the operator: an empty log here would mean
		// the branch sanitized by DELETING the information rather than routing it.
		line := strings.TrimRight(errLog.String(), "\n")
		if line == "" {
			t.Fatalf("POST /v1/commit wrote a sanitized 500 but logged nothing — the detail was destroyed, not routed")
		}
		if !strings.Contains(line, "POST /v1/commit") {
			t.Errorf("commit error-log line = %q, want it to name the route %q", line, "POST /v1/commit")
		}
	})
}

// assertErrorLogLine is (b): the operator stream got ONE line, it names the
// route, and it carries the seeded detail verbatim.
//
// The one-line property is asserted, not assumed. "One line per error" is what
// makes the stream greppable under a 500 storm, and a %+v dump would satisfy a
// bare Contains check while burying the log in stack traces.
func assertErrorLogLine(t *testing.T, logged, label, method, target string) {
	t.Helper()
	if logged == "" {
		t.Errorf("%s: the error log is EMPTY — the 500's detail was destroyed rather than routed to the operator", label)
		return
	}
	if !strings.HasSuffix(logged, "\n") {
		t.Errorf("%s: error-log write %q is not newline-terminated", label, logged)
	}
	line := strings.TrimRight(logged, "\n")
	if strings.Contains(line, "\n") {
		t.Errorf("%s: the error log wrote %d lines for one error, want exactly 1:\n%s",
			label, strings.Count(line, "\n")+1, logged)
	}
	if !strings.Contains(line, internalDetailSentinel) {
		t.Errorf("%s: error-log line = %q, want it to carry the VERBATIM detail %q that the body no longer does",
			label, line, internalDetailSentinel)
	}
	// The route, not the raw target: r.URL.Path is logged, deliberately without
	// RawQuery, so /v1/log?from=0&limit=5 logs as GET /v1/log.
	//
	// The trailing colon is LOAD-BEARING and is why this is not a bare Contains
	// on the route. The producer's format is "%s %s: %v", so the route is
	// always followed immediately by ": ". Without the colon, "GET /v1/log" is
	// a PREFIX of "GET /v1/log?from=0&limit=5", so logging r.URL.String()
	// instead of r.URL.Path would leak client-supplied query text into the
	// operator stream and this assertion would still pass — measured: that
	// exact mutation SURVIVED the whole host/daemon suite (rc=0) before the
	// colon was added, and reds here with it.
	path, _, _ := strings.Cut(target, "?")
	if want := method + " " + path + ":"; !strings.Contains(line, want) {
		t.Errorf("%s: error-log line = %q, want it to name the route %q — a detail line with no route is unattributable, and an UNANCHORED route match would accept a line carrying client-supplied query text",
			label, line, want)
	}
}

// TestExpiredReadDeadlineExpiresAtConstruction pins the property every 503
// assertion in this file now rests on: the stimulus context is DEAD BEFORE any
// store read can begin. Two assertions, two mutations:
//   - the sign check kills the "shrink it back to a positive nanosecond"
//     mutation deterministically (no timing anywhere);
//   - the ctx.Err() check goes through the production readCtx, so a readCtx
//     that ignores d.readDeadline (or re-derives from the 10s constant) reds
//     here in one run.
func TestExpiredReadDeadlineExpiresAtConstruction(t *testing.T) {
	if expiredReadDeadline >= 0 {
		t.Fatalf("expiredReadDeadline = %s, want a negative duration — a positive value arms "+
			"a timer and re-creates the 200-vs-503 race this constant exists to remove",
			expiredReadDeadline)
	}
	d := newHandlerDaemon(t)
	d.readDeadline = expiredReadDeadline
	ctx, cancel := d.readCtx(httptest.NewRequest(http.MethodGet, "/v1/head", nil))
	defer cancel()
	if ctx.Err() == nil {
		t.Fatalf("readCtx under an already-expired deadline returned a LIVE context — the " +
			"stimulus must be expired at construction, before any store read can begin")
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("ctx.Err() = %v, want context.DeadlineExceeded", ctx.Err())
	}
}
