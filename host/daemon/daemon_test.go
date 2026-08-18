package daemon

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// fakeInterpreter writes a synthetic executable that echoes version when
// invoked with --version, and returns its path plus the exact bytes it holds
// (which are what host/archive content-addresses). It mirrors the helper of the
// same name in host/archive's tests, so the daemon's health facts are exercised
// hermetically — without depending on AILANG_BIN being exported.
func fakeInterpreter(t *testing.T, dir, version string) (string, []byte) {
	t.Helper()
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then printf '%s' '" + version + "'; exit 0; fi\n"
	p := filepath.Join(dir, "ailang-fake")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake interpreter: %v", err)
	}
	return p, []byte(script)
}

// getBody performs a bounded GET and returns the status code and body text.
func getBody(t *testing.T, url string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("build request %s: %v", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		t.Fatalf("read %s body: %v", url, err)
	}
	return resp.StatusCode, string(b)
}

// seedHead stores a genesis world revision and selects it, so GET /v1/head has
// a real canonical value to round-trip. Mirrors host/store's own seedGenesis.
func seedHead(t *testing.T, s *store.Store) store.World {
	t.Helper()
	w := store.World{
		Ref:       hashref.SumSHA256([]byte("daemon-test-world-genesis")),
		Revision:  0,
		StateRoot: hashref.SumSHA256([]byte("daemon-test-state-genesis")),
		LogHead:   hashref.SumSHA256([]byte("daemon-test-log-genesis")),
	}
	if err := s.PutWorld(w); err != nil {
		t.Fatalf("seed PutWorld: %v", err)
	}
	if err := s.SelectHead(w.Ref); err != nil {
		t.Fatalf("seed SelectHead: %v", err)
	}
	return w
}

// ---------------------------------------------------------------------------
// Decision 4 — local-first is structural, not advisory
// ---------------------------------------------------------------------------

// TestIsLoopbackHostMirrorsSketchPredicate replays the frozen sketch's OWN test
// vectors (design_docs/sketches/worlddapi.ail, isLoopbackHost — Z3-proven with
// `ensures { result == (host == "127.0.0.1" || host == "::1" || host ==
// "localhost") }`) against the Go mirror, plus the refusal arm the design doc
// names explicitly. BOTH arms are asserted: a predicate that always returned
// false would pass an accept-only table, and one that always returned true would
// pass a refuse-only table.
func TestIsLoopbackHostMirrorsSketchPredicate(t *testing.T) {
	cases := []struct {
		host string
		want bool
	}{
		// Accepted arm. "127.0.0.1" and "::1" are the sketch's own true vectors.
		{"127.0.0.1", true},
		{"::1", true},
		{"localhost", true},
		// Refused arm. "0.0.0.0" and "example.com" are the sketch's own false
		// vectors; 192.168.1.10 is the design doc's LAN-exposure case.
		{"0.0.0.0", false},
		{"example.com", false},
		{"192.168.1.10", false},
		{"", false},
	}
	sawTrue, sawFalse := false, false
	for _, c := range cases {
		if got := isLoopbackHost(c.host); got != c.want {
			t.Errorf("isLoopbackHost(%q) = %v, want %v (drift from the Z3-proven sketch predicate)",
				c.host, got, c.want)
		}
		if c.want {
			sawTrue = true
		} else {
			sawFalse = true
		}
	}
	// Non-vacuity: the table must genuinely exercise both arms.
	if !sawTrue || !sawFalse {
		t.Fatalf("table is one-sided (sawTrue=%v sawFalse=%v); a constant predicate would pass it",
			sawTrue, sawFalse)
	}
}

// TestNewRefusesNonLoopbackBind proves the refusal is enforced at the lifecycle
// level, not merely available as a predicate: New must return a named,
// errors.Is-matchable failure for every non-loopback host and must succeed for
// every loopback one.
//
// It additionally asserts that a refused startup creates NO database file —
// proving the bind policy is checked BEFORE any writer lock is taken, so a
// misconfigured bind can never disturb a database another process is serving.
func TestNewRefusesNonLoopbackBind(t *testing.T) {
	refused := []string{"0.0.0.0", "example.com", "192.168.1.10"}
	for _, host := range refused {
		t.Run("refused/"+host, func(t *testing.T) {
			dbPath := filepath.Join(t.TempDir(), "world.db")
			d, err := New(Config{DBPath: dbPath, BindHost: host, BindPort: DefaultBindPort})
			if err == nil {
				_ = d.Close()
				t.Fatalf("New accepted non-loopback bind host %q — local-first must be structural", host)
			}
			if !errors.Is(err, ErrNonLoopbackBind) {
				t.Fatalf("New(%q) error = %v, want one wrapping ErrNonLoopbackBind", host, err)
			}
			var se *StartupError
			if !errors.As(err, &se) {
				t.Fatalf("New(%q) error = %T, want *StartupError", host, err)
			}
			if se.Stage != StageBindPolicy {
				t.Fatalf("New(%q) StartupError.Stage = %q, want %q", host, se.Stage, StageBindPolicy)
			}
			if _, statErr := os.Stat(dbPath); !os.IsNotExist(statErr) {
				t.Fatalf("New(%q) created %s despite refusing the bind; the bind policy must be "+
					"checked before the store is opened", host, dbPath)
			}
		})
	}

	accepted := []string{"127.0.0.1", "::1", "localhost"}
	for _, host := range accepted {
		t.Run("accepted/"+host, func(t *testing.T) {
			dbPath := filepath.Join(t.TempDir(), "world.db")
			d, err := New(Config{DBPath: dbPath, BindHost: host, BindPort: 0})
			if err != nil {
				t.Fatalf("New rejected loopback bind host %q: %v", host, err)
			}
			if err := d.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Decision 7 — bounded waits & allocations
// ---------------------------------------------------------------------------

// TestBoundedWaitsAndBodyLimit is the D7 non-vacuity gate.
//
//	part (a): the constructed server's four timeout fields EQUAL the named
//	          constants and are NON-ZERO. A zero value in any of them is a test
//	          failure, not a default.
//	part (c): maxCommitBytes equals the Z3-proven withinCommitBytes bound in
//	          design_docs/sketches/worlddapi.ail, so the Go constant and the
//	          frozen semantic bound cannot drift silently.
//
// Part (b) — the 413 on an oversized POST /v1/commit body followed by a valid
// small commit — is exercised in handlers_test.go, where the M2.B route exists.
//
// The named constants are themselves pinned against the D7 table, so editing a
// constant to match a drifted server (or vice versa) still turns this red.
func TestBoundedWaitsAndBodyLimit(t *testing.T) {
	// The D7 table (design_docs/planned/w-worldd-m2.md, Decision 7), restated
	// as literals so the constants cannot be silently retuned.
	t.Run("constants match the D7 table", func(t *testing.T) {
		for _, c := range []struct {
			name string
			got  time.Duration
			want time.Duration
		}{
			{"readHeaderTimeout", readHeaderTimeout, 5 * time.Second},
			{"readTimeout", readTimeout, 30 * time.Second},
			{"writeTimeout", writeTimeout, 30 * time.Second},
			{"idleTimeout", idleTimeout, 120 * time.Second},
			{"defaultClientTimeout", defaultClientTimeout, 30 * time.Second},
			{"shutdownTimeout", shutdownTimeout, 10 * time.Second},
		} {
			if c.got != c.want {
				t.Errorf("%s = %s, want %s (D7 table)", c.name, c.got, c.want)
			}
		}
		if DefaultClientTimeout != defaultClientTimeout {
			t.Errorf("DefaultClientTimeout = %s, want %s — the CLI must not carry a second deadline",
				DefaultClientTimeout, defaultClientTimeout)
		}
	})

	// part (a): every http.Server wait is bounded by its named constant.
	t.Run("a/server timeouts are set and non-zero", func(t *testing.T) {
		srv := newServer(http.NewServeMux())
		for _, c := range []struct {
			name string
			got  time.Duration
			want time.Duration
		}{
			{"ReadHeaderTimeout", srv.ReadHeaderTimeout, readHeaderTimeout},
			{"ReadTimeout", srv.ReadTimeout, readTimeout},
			{"WriteTimeout", srv.WriteTimeout, writeTimeout},
			{"IdleTimeout", srv.IdleTimeout, idleTimeout},
		} {
			if c.got == 0 {
				t.Errorf("http.Server.%s is zero — an unbounded wait is a test failure, not a default", c.name)
				continue
			}
			if c.got != c.want {
				t.Errorf("http.Server.%s = %s, want the named constant %s", c.name, c.got, c.want)
			}
		}
	})

	// part (b): both arms are load-bearing. The oversized request must trip the
	// MaxBytesReader cap, and a subsequent valid request must prove the cap is
	// not simply rejecting every commit.
	t.Run("b/commit body cap rejects only oversized requests", func(t *testing.T) {
		d := newHandlerDaemon(t)
		oversized := bytes.NewReader(bytes.Repeat([]byte{'x'}, maxCommitBytes+1))
		assertErrorClass(t, requestRecorder(t, d, http.MethodPost, "/v1/commit", oversized),
			http.StatusRequestEntityTooLarge, "PayloadTooLarge")

		genesis := seedGenesisEmbedded(t, d, "body-cap")
		valid := testCommit(genesis, 1, "body-cap")
		rec := requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(valid)))
		if rec.Code != http.StatusOK {
			t.Fatalf("small valid commit after oversized body: status=%d body=%s", rec.Code, rec.Body)
		}
	})

	// part (c): the body cap is the Z3-proven bound, not a tunable.
	t.Run("c/maxCommitBytes matches the Z3-proven sketch bound", func(t *testing.T) {
		const z3ProvenBound = 8388608 // withinCommitBytes: n >= 0 && n <= 8388608
		if maxCommitBytes != z3ProvenBound {
			t.Fatalf("maxCommitBytes = %d, want %d — the Go constant has drifted from the "+
				"Z3-proven withinCommitBytes bound in design_docs/sketches/worlddapi.ail",
				maxCommitBytes, z3ProvenBound)
		}
	})
}

// TestBoundedGracefulShutdownDrainsInFlightRequest proves (*Daemon).Shutdown —
// the shipped shutdown path, not http.Server in general — is a REAL drain and is
// bounded:
//
//   - it WAITS for an in-flight request (the response body arrives intact, and
//     the shutdown takes at least the handler's own work time — a bare
//     srv.Close() would fail both assertions), and
//   - it completes well under shutdownTimeout.
//
// It drives a real Daemon built by New and only swaps in a deliberately slow
// handler: the two M2.A routes complete in microseconds, so with them the drain
// window is unobservable and the assertions would be vacuous.
func TestBoundedGracefulShutdownDrainsInFlightRequest(t *testing.T) {
	const handlerWork = 250 * time.Millisecond
	const wantBody = "drained"

	d, err := New(Config{
		DBPath:   filepath.Join(t.TempDir(), "world.db"),
		BindHost: DefaultBindHost,
		BindPort: 0,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })

	var once sync.Once
	started := make(chan struct{})
	d.srv = newServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		once.Do(func() { close(started) })
		time.Sleep(handlerWork)
		_, _ = io.WriteString(w, wantBody)
	}))

	if err := d.Listen(); err != nil {
		t.Fatalf("Listen: %v", err)
	}
	ln := d.ln
	served := make(chan error, 1)
	go func() { served <- d.Serve() }()

	// NOTE: this runs on a non-test goroutine, so it reports through the channel
	// rather than calling t.Fatalf (which may only be called from the test
	// goroutine).
	body := make(chan string, 1)
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String() + "/slow") //nolint:noctx // bounded by the server's D7 timeouts
		if err != nil {
			body <- "REQUEST ERROR: " + err.Error()
			return
		}
		defer func() { _ = resp.Body.Close() }()
		b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			body <- "READ ERROR: " + err.Error()
			return
		}
		body <- string(b)
	}()

	<-started
	start := time.Now()
	if err := d.Shutdown(); err != nil {
		t.Fatalf("Shutdown of a single in-flight request failed: %v", err)
	}
	elapsed := time.Since(start)

	if got := <-body; got != wantBody {
		t.Fatalf("in-flight response body = %q, want %q — Shutdown hard-closed a live request", got, wantBody)
	}
	if elapsed < handlerWork/2 {
		t.Fatalf("Shutdown returned after %s, far under the %s of in-flight work — it did not actually drain",
			elapsed, handlerWork)
	}
	if elapsed >= shutdownTimeout {
		t.Fatalf("Shutdown took %s, which is not bounded by shutdownTimeout %s", elapsed, shutdownTimeout)
	}
	if err := <-served; !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("Serve returned %v, want http.ErrServerClosed", err)
	}
}

// TestDaemonShutdownIsBoundedByTheD7Constant proves the SHIPPED Shutdown path is
// bounded by the D7 constant, and that its expiry branch really hard-closes:
//
//   - New wires drainTimeout to shutdownTimeout (a raised bound is red here), and
//   - with the bound shrunk, a request that outlives it makes Shutdown return a
//     deadline error within the budget instead of waiting the handler out.
//
// Without the second half, replacing shutdownTimeout with an arbitrarily large
// value would still look green — the drain simply never expires in practice.
func TestDaemonShutdownIsBoundedByTheD7Constant(t *testing.T) {
	d, err := New(Config{
		DBPath:   filepath.Join(t.TempDir(), "world.db"),
		BindHost: DefaultBindHost,
		BindPort: 0,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })

	if d.drainTimeout != shutdownTimeout {
		t.Fatalf("New wired drainTimeout = %s, want the D7 shutdownTimeout %s", d.drainTimeout, shutdownTimeout)
	}

	const budget = 100 * time.Millisecond
	const handlerWork = 5 * time.Second // deliberately far beyond the budget
	d.drainTimeout = budget

	var once sync.Once
	started := make(chan struct{})
	d.srv = newServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		once.Do(func() { close(started) })
		time.Sleep(handlerWork)
		_, _ = io.WriteString(w, "too late")
	}))
	if err := d.Listen(); err != nil {
		t.Fatalf("Listen: %v", err)
	}
	served := make(chan error, 1)
	go func() { served <- d.Serve() }()
	go func() {
		resp, err := http.Get("http://" + d.Addr() + "/slow") //nolint:noctx // the server's D7 timeouts bound it
		if err == nil {
			_ = resp.Body.Close()
		}
	}()

	<-started
	start := time.Now()
	err = d.Shutdown()
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("Shutdown returned nil after waiting out a %s handler under a %s budget — "+
			"an incomplete drain must be reported", handlerWork, budget)
	}
	if elapsed >= handlerWork/2 {
		t.Fatalf("Shutdown took %s under a %s budget — the drain is not bounded by drainTimeout", elapsed, budget)
	}
	if err := <-served; !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("Serve returned %v, want http.ErrServerClosed", err)
	}
}

// TestNewReleasesWriterAuthorityOnLateFailure proves a REFUSED startup never
// strands writer authority: when the lifecycle fails after store.Open (here at
// the archive stage), the store is closed and the cross-process writer lock is
// released, so the next process can open the same database.
func TestNewReleasesWriterAuthorityOnLateFailure(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	d, err := New(Config{
		DBPath:    dbPath,
		BindHost:  DefaultBindHost,
		BindPort:  0,
		AilangBin: filepath.Join(t.TempDir(), "no-such-interpreter"),
	})
	if err == nil {
		_ = d.Close()
		t.Fatal("New succeeded with a nonexistent --ailang-bin; the archive step must be fatal")
	}
	var se *StartupError
	if !errors.As(err, &se) {
		t.Fatalf("New error = %T (%v), want *StartupError", err, err)
	}
	if se.Stage != StageArchive {
		t.Fatalf("StartupError.Stage = %q, want %q", se.Stage, StageArchive)
	}

	// The load-bearing assertion: writer authority was released on the way out.
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store.Open after a refused startup: %v — the failed New stranded writer authority", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// stuckServer never finishes its drain, so the deadline-expiry branch of drain
// can be exercised honestly. A production server should never reach this
// branch, which is exactly why it needs a test rather than an assumption.
type stuckServer struct {
	closed atomic.Bool
}

func (s *stuckServer) Shutdown(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }
func (s *stuckServer) Close() error                       { s.closed.Store(true); return nil }

// TestDrainHardClosesWhenDeadlineExpires proves the drain is never unbounded:
// when the graceful phase does not finish, drain hard-Closes and REPORTS the
// incomplete drain (so the process exits non-zero) instead of waiting forever.
func TestDrainHardClosesWhenDeadlineExpires(t *testing.T) {
	const budget = 50 * time.Millisecond
	s := &stuckServer{}

	start := time.Now()
	err := drain(s, budget)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("drain returned nil for a drain that never finished — an incomplete drain must be reported")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("drain error = %v, want one wrapping context.DeadlineExceeded", err)
	}
	if !s.closed.Load() {
		t.Fatal("drain did not hard-Close after the deadline expired — connections would leak")
	}
	if elapsed > 20*budget {
		t.Fatalf("drain took %s for a %s budget — the wait is not bounded by the deadline", elapsed, budget)
	}
}

// ---------------------------------------------------------------------------
// Routes — real values round-trip, not merely a 200
// ---------------------------------------------------------------------------

// TestHealthAndHeadRoundTrip drives the two M2.A routes over a real HTTP
// round-trip (httptest) and asserts the ACTUAL VALUES:
//
//   - /v1/health carries the daemon version, the configured db path, the
//     archived interpreter's HashRef (independently recomputed from the
//     interpreter bytes) and the verbatim `--version` text from the manifest;
//   - /v1/head is 404 before a head exists (the sketch's NotFound class) and,
//     once seeded, returns the exact canonical "algo:digest" text.
//
// A handler that returned 200 with an empty body would fail every assertion.
func TestHealthAndHeadRoundTrip(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only (mirrors host/archive's tests)")
	}
	dir := t.TempDir()
	const fakeVersion = "FakeAILANG v9.9.9\ncommit: deadbeef\n"
	execPath, execBytes := fakeInterpreter(t, dir, fakeVersion)
	wantRef := hashref.SumSHA256(execBytes)
	dbPath := filepath.Join(dir, "world.db")

	d, err := New(Config{DBPath: dbPath, BindHost: DefaultBindHost, BindPort: 0, AilangBin: execPath})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })

	ts := httptest.NewServer(d.Handler())
	t.Cleanup(ts.Close)

	// --- GET /v1/health -----------------------------------------------------
	code, raw := getBody(t, ts.URL+"/v1/health")
	if code != http.StatusOK {
		t.Fatalf("GET /v1/health status = %d (body %q), want 200", code, raw)
	}
	var health HealthResponse
	if err := json.Unmarshal([]byte(raw), &health); err != nil {
		t.Fatalf("decode health body %q: %v", raw, err)
	}
	if health.DaemonVersion != Version {
		t.Errorf("health.daemon_version = %q, want %q", health.DaemonVersion, Version)
	}
	if health.DBPath != dbPath {
		t.Errorf("health.db_path = %q, want %q", health.DBPath, dbPath)
	}
	if health.InterpreterRef != wantRef.String() {
		t.Errorf("health.interpreter_ref = %q, want %q (the SHA-256 of the archived bytes)",
			health.InterpreterRef, wantRef.String())
	}
	if health.InterpreterVersion != fakeVersion {
		t.Errorf("health.interpreter_version = %q, want the manifest's verbatim %q",
			health.InterpreterVersion, fakeVersion)
	}

	// --- GET /v1/head, before any head is selected ---------------------------
	code, body := getBody(t, ts.URL+"/v1/head")
	if code != http.StatusNotFound {
		t.Fatalf("GET /v1/head with no selected head: status = %d (body %q), want 404", code, body)
	}
	var headErr APIError
	if err := json.Unmarshal([]byte(body), &headErr); err != nil || headErr.Error.Class != "NotFound" {
		t.Fatalf("GET /v1/head error body = %q, want JSON NotFound envelope (decode err=%v)", body, err)
	}

	// --- GET /v1/head, after seeding ----------------------------------------
	world := seedHead(t, d.store)
	code, body = getBody(t, ts.URL+"/v1/head")
	if code != http.StatusOK {
		t.Fatalf("GET /v1/head status = %d (body %q), want 200", code, body)
	}
	if body != world.Ref.String() {
		t.Fatalf("GET /v1/head body = %q, want the canonical head text %q", body, world.Ref.String())
	}
	if !strings.HasPrefix(body, hashref.AlgoSHA256+":") {
		t.Fatalf("GET /v1/head body = %q, want canonical %q text", body, "algo:digest")
	}
}

// TestRunAnnouncesResolvedListenAddress proves the stdout announcement contract
// M2.C's ephemeral-port end-to-end test depends on: with --bind 127.0.0.1:0 the
// kernel picks the port, and the ONLY way a caller can learn it is this line.
//
// It asserts the announced address is real by serving a request against it, and
// that a cancelled context produces an orderly, bounded Run return.
func TestRunAnnouncesResolvedListenAddress(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	cfg := Config{DBPath: dbPath, BindHost: DefaultBindHost, BindPort: 0}

	pr, pw := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ran := make(chan error, 1)
	go func() {
		err := Run(ctx, cfg, pw)
		_ = pw.Close()
		ran <- err
	}()

	line, err := bufio.NewReader(pr).ReadString('\n')
	if err != nil {
		t.Fatalf("read listen announcement: %v", err)
	}
	if !strings.HasPrefix(line, ListenAnnouncePrefix) {
		t.Fatalf("announcement = %q, want the stable prefix %q", line, ListenAnnouncePrefix)
	}
	url := strings.TrimSpace(strings.TrimPrefix(line, ListenAnnouncePrefix))
	if !strings.HasPrefix(url, "http://"+DefaultBindHost+":") {
		t.Fatalf("announced url = %q, want http://%s:<port>", url, DefaultBindHost)
	}
	port := url[strings.LastIndex(url, ":")+1:]
	if port == "" || port == "0" {
		t.Fatalf("announced url = %q, want a RESOLVED ephemeral port, not the requested 0", url)
	}

	// The announced address must actually be serving before the line is written.
	code, body := getBody(t, url+"/v1/health")
	if code != http.StatusOK {
		t.Fatalf("GET %s/v1/health status = %d (body %q), want 200 — the announcement preceded a live socket",
			url, code, body)
	}

	cancel()
	select {
	case err := <-ran:
		if err != nil {
			t.Fatalf("Run returned %v, want a clean bounded shutdown", err)
		}
	case <-time.After(shutdownTimeout + 5*time.Second):
		t.Fatalf("Run did not return within %s of cancellation — the shutdown path is unbounded",
			shutdownTimeout+5*time.Second)
	}

	// Writer authority was released, so the database can be reopened.
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store.Open after Run returned: %v — Run leaked writer authority", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestNewBootstrapsEpochRegistryIdempotently proves the registry step of the
// lifecycle is real and repeatable: a second New over the same database (after
// the first has released writer authority) finds the SAME epoch-registry head
// rather than creating a divergent epoch 1.
func TestNewBootstrapsEpochRegistryIdempotently(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	ailangBin, _ := fakeInterpreter(t, dir, "AILANG v0.30.0\nCommit: e37b370\n")
	cfg := Config{
		DBPath:    dbPath,
		BindHost:  DefaultBindHost,
		BindPort:  0,
		AilangBin: ailangBin,
	}

	heads := make([]string, 0, 2)
	for i := range 2 {
		d, err := New(cfg)
		if err != nil {
			t.Fatalf("New #%d: %v", i+1, err)
		}
		head, ok, err := d.store.GetRegistryHead(context.Background(), registry.SemanticID)
		if err != nil || !ok {
			_ = d.Close()
			t.Fatalf("GetRegistryHead #%d: ok=%v err=%v — the lifecycle must bootstrap the registry", i+1, ok, err)
		}
		heads = append(heads, head.String())
		if err := d.Close(); err != nil {
			t.Fatalf("Close #%d: %v", i+1, err)
		}
	}
	if heads[0] != heads[1] {
		t.Fatalf("registry head moved across restarts: %q then %q — bootstrap is not idempotent",
			heads[0], heads[1])
	}
	if !strings.HasPrefix(heads[0], hashref.AlgoSHA256+":") {
		t.Fatalf("registry head = %q, want canonical %q text", heads[0], "algo:digest")
	}
}

// TestReleaseFromVersion pins the epoch-registry candidate derivation: the
// verbatim multi-line `ailang --version` output reduces to its first non-empty
// line, and an empty probe falls back to the documented sentinel.
func TestReleaseFromVersion(t *testing.T) {
	cases := []struct{ in, want string }{
		{"AILANG v0.30.0\nCommit: e37b370\n", "AILANG v0.30.0"},
		{"\n\n  AILANG v0.31.0  \nrest", "AILANG v0.31.0"},
		{"", unpinnedRelease},
		{"   \n\t\n", unpinnedRelease},
	}
	for _, c := range cases {
		if got := releaseFromVersion(c.in); got != c.want {
			t.Errorf("releaseFromVersion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Zero-cloud dependency allowlist (Decision 4: enforced, not asserted)
// ---------------------------------------------------------------------------

// allowedDepModules is the closed set of non-stdlib module roots the daemon
// core (host/daemon + cmd/ailang-worldd) may reach, directly or transitively.
// The charter guardrail "local-first is inviolable" is enforced here: a cloud
// SDK, telemetry shim or network client cannot enter the build graph without
// this literal being edited, which is a reviewable event.
//
// Everything below modernc.org/sqlite is the indirect chain that the pure-Go
// SQLite driver drags in; it is all local computation (libc shim, allocator,
// bignum, strftime, tty probe, uuid) with no outbound network capability.
var allowedDepModules = []string{
	"github.com/sunholo-data/ailang-world", // this repo
	"modernc.org/sqlite",                   // pure-Go SQLite driver (host/store)
	"modernc.org/libc",                     // indirect: sqlite
	"modernc.org/mathutil",                 // indirect: sqlite
	"modernc.org/memory",                   // indirect: sqlite
	"golang.org/x/sys",                     // indirect: sqlite
	"github.com/dustin/go-humanize",        // indirect: sqlite
	"github.com/google/uuid",               // indirect: sqlite
	"github.com/mattn/go-isatty",           // indirect: sqlite
	"github.com/ncruces/go-strftime",       // indirect: sqlite
	"github.com/remyoudompheng/bigfft",     // indirect: sqlite
}

// daemonCorePatterns are the two package trees the allowlist governs.
var daemonCorePatterns = []string{"./host/daemon/...", "./cmd/ailang-worldd/..."}

// isStdlibImportPath reports whether an import path belongs to the Go standard
// library, using the canonical rule: a stdlib path's FIRST segment contains no
// dot (real modules are rooted at a domain, e.g. "modernc.org/sqlite").
//
// We deliberately use this string rule rather than `go list -deps -f
// '{{.Standard}}'`. Both are correct here — verified against the toolchain,
// including the GOROOT-vendored packages ("vendor/golang.org/x/net/...",
// "crypto/internal/entropy/v1.0.0"), which the dot rule classifies as stdlib
// (first segments "vendor" and "crypto") exactly as .Standard=true does. The
// dot rule wins because it is a pure function: it is unit-testable below
// without a toolchain, so the classifier itself is gated, not just its output.
func isStdlibImportPath(p string) bool {
	first, _, _ := strings.Cut(p, "/")
	return !strings.Contains(first, ".")
}

// disallowedDeps partitions a dependency list against allowedDepModules and
// returns the offending package paths BY NAME, so an operator reading CI output
// knows precisely what leaked in.
//
// It returns an error on an empty list. Per coding-standards S6 a gate must
// fail on its null case: if `go list` ever emits nothing (wrong working
// directory, silently-broken patterns, a mis-scoped refactor) an allowlist that
// vacuously passed would be a false green, which is the exact defect class that
// bit V27 and B1.
func disallowedDeps(deps []string) ([]string, error) {
	if len(deps) == 0 {
		return nil, errors.New("dependency list is empty: the allowlist would pass vacuously — " +
			"`go list -deps` produced no packages for " + strings.Join(daemonCorePatterns, " "))
	}
	var bad []string
	for _, d := range deps {
		if isStdlibImportPath(d) {
			continue
		}
		allowed := false
		for _, m := range allowedDepModules {
			if d == m || strings.HasPrefix(d, m+"/") {
				allowed = true
				break
			}
		}
		if !allowed {
			bad = append(bad, d)
		}
	}
	return bad, nil
}

// goListDeps shells out to `go list -deps` over the given patterns, rooted at
// the repository. It returns an error — never a skip — when the toolchain is
// unreachable. `go` is resolved from PATH only, on purpose: a single, obvious
// resolution path means a CI failure here reads as "go is missing" rather than
// silently falling back to some other toolchain.
func goListDeps(repoRoot string, patterns ...string) ([]string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("cannot locate the `go` toolchain on PATH: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	args := append([]string{"list", "-deps"}, patterns...)
	cmd := exec.CommandContext(ctx, goBin, args...)
	cmd.Dir = repoRoot
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s list -deps %s (in %s): %w\nstderr: %s",
			goBin, strings.Join(patterns, " "), repoRoot, err, stderr.String())
	}
	var deps []string
	sc := bufio.NewScanner(strings.NewReader(stdout.String()))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			deps = append(deps, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan `go list` output: %w", err)
	}
	return deps, nil
}

// repoRootFromPackageDir returns the module root as seen from this package's
// working directory (<root>/host/daemon), failing loudly if go.mod is not
// there — a wrong root would yield an empty dep list, i.e. a vacuous gate.
func repoRootFromPackageDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("repo root %s has no go.mod: %v", root, err)
	}
	return root
}

// TestDaemonDependencyAllowlist enforces Decision 4's zero-cloud claim as a
// gate rather than prose: the daemon core's full transitive build graph must
// consist only of the Go standard library, this repository, and the explicitly
// enumerated local-only modules above. It also proves its own non-vacuity —
// the null case, the classifier and the by-name reporting are each asserted.
func TestDaemonDependencyAllowlist(t *testing.T) {
	t.Run("null case fails", func(t *testing.T) {
		if _, err := disallowedDeps(nil); err == nil {
			t.Fatal("disallowedDeps(nil) returned no error: the gate would pass vacuously on an empty dependency list")
		}
		if _, err := disallowedDeps([]string{}); err == nil {
			t.Fatal("disallowedDeps([]) returned no error: the gate would pass vacuously on an empty dependency list")
		}
	})

	t.Run("classifier separates stdlib from modules", func(t *testing.T) {
		cases := []struct {
			in   string
			want bool
		}{
			{"net/http", true},
			{"crypto/internal/entropy/v1.0.0", true},
			{"vendor/golang.org/x/net/dns/dnsmessage", true},
			{"unsafe", true},
			{"modernc.org/sqlite", false},
			{"cloud.google.com/go/storage", false},
			{"github.com/sunholo-data/ailang-world/host/store", false},
		}
		for _, c := range cases {
			if got := isStdlibImportPath(c.in); got != c.want {
				t.Errorf("isStdlibImportPath(%q) = %v, want %v", c.in, got, c.want)
			}
		}
	})

	t.Run("an intruder is reported by name", func(t *testing.T) {
		bad, err := disallowedDeps([]string{
			"net/http",
			"modernc.org/sqlite",
			"cloud.google.com/go/storage",
			"github.com/aws/aws-sdk-go-v2/service/s3",
		})
		if err != nil {
			t.Fatalf("disallowedDeps: %v", err)
		}
		want := []string{"cloud.google.com/go/storage", "github.com/aws/aws-sdk-go-v2/service/s3"}
		if strings.Join(bad, ",") != strings.Join(want, ",") {
			t.Fatalf("disallowedDeps = %v, want exactly %v", bad, want)
		}
	})

	t.Run("daemon core build graph", func(t *testing.T) {
		root := repoRootFromPackageDir(t)
		deps, err := goListDeps(root, daemonCorePatterns...)
		if err != nil {
			// Loud by construction: never t.Skip. A skipped allowlist is a
			// false green, and the zero-cloud guardrail would go unenforced.
			t.Fatalf("could not enumerate the daemon dependency graph, so the zero-cloud "+
				"allowlist could NOT be enforced (this is a failure, not a skip): %v", err)
		}
		bad, err := disallowedDeps(deps)
		if err != nil {
			t.Fatalf("dependency allowlist could not be evaluated: %v", err)
		}
		if len(bad) > 0 {
			t.Fatalf("%d package(s) outside the zero-cloud allowlist reached the daemon core "+
				"(patterns %s):\n  %s\nEach must be either removed or added to allowedDepModules "+
				"with a justification (charter guardrail: local-first is inviolable).",
				len(bad), strings.Join(daemonCorePatterns, " "), strings.Join(bad, "\n  "))
		}
		t.Logf("zero-cloud allowlist enforced over %d transitive packages for %s",
			len(deps), strings.Join(daemonCorePatterns, " "))
	})
}
