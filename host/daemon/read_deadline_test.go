package daemon

import (
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
		d.readDeadline = 1 * time.Nanosecond

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
	d.readDeadline = 1 * time.Nanosecond

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
