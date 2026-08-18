package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// This file holds the store-layer half of w-daemon-read-cancellation (queue
// item 18, M1): the read path's elapsed-time bound is the CONTEXT, and the
// lock layer's retry policy is busy_timeout. The two are different mechanisms
// at different layers and each has its own test here, because the item's whole
// thesis is that the existing gates could not see the layer where the wait
// happens.

// -----------------------------------------------------------------------------
// T1.5 — the context reaches the blocking database/sql call.
// -----------------------------------------------------------------------------

// watchdogBound is how long a getter that HONOURS its context may take before
// the arm is declared failed. Every arm under test answers in microseconds when
// the propagation is present; a mutant that drops the context parks in the pool
// wait forever, so any value comfortably above scheduling noise discriminates.
const watchdogBound = 2 * time.Second

// TestReadGettersHonorContext proves that each of the five daemon-path read
// getters carries its caller's context all the way to the blocking
// database/sql call.
//
// The stimulus is the real unbounded wait this item exists to bound: the store
// runs with SetMaxOpenConns(1), so a test that takes the sole pool connection
// makes every other read queue behind it with no bound of its own. A getter
// given an ALREADY-EXPIRED context must refuse immediately; a getter that
// passes context.Background() to database/sql instead (MU4a-e) parks in the
// pool wait and never returns.
//
// R2 / watchdog discipline: the red path does not merely t.Error. It RELEASES
// the blocked getter by closing the occupying *sql.Conn and then drains the
// result channel under a second bound. Without that release the parked getter
// holds the sole connection, openMem's deferred Close queues behind it, and the
// arm converts a clean failure into a suite-wide hang.
func TestReadGettersHonorContext(t *testing.T) {
	cases := []struct {
		name string
		call func(context.Context, *Store) error
	}{
		{"GetObject", func(ctx context.Context, s *Store) error {
			_, _, err := s.GetObject(ctx, hashref.SumSHA256([]byte("absent-object")))
			return err
		}},
		{"GetWorld", func(ctx context.Context, s *Store) error {
			_, _, err := s.GetWorld(ctx, hashref.SumSHA256([]byte("absent-world")))
			return err
		}},
		{"GetLogEntry", func(ctx context.Context, s *Store) error {
			_, _, err := s.GetLogEntry(ctx, 0)
			return err
		}},
		{"GetRegistryHead", func(ctx context.Context, s *Store) error {
			_, _, err := s.GetRegistryHead(ctx, "world/epoch-registry/v1")
			return err
		}},
		{"SelectedHead", func(ctx context.Context, s *Store) error {
			_, _, err := s.SelectedHead(ctx)
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)

			// Occupy the sole pool connection. Every subsequent read must wait
			// for it, which is exactly the wait the context has to bound.
			conn, err := s.db.Conn(context.Background())
			if err != nil {
				t.Fatalf("take the sole pool connection: %v", err)
			}
			released := false
			release := func() {
				if !released {
					released = true
					_ = conn.Close()
				}
			}
			defer release()

			// Known-positive control: with the connection held, an ordinary
			// context-carrying query on the SAME pool must fail to get a
			// connection and return its context's error. This proves the
			// occupation is real, so the arm below is non-vacuous.
			//
			// The control deliberately does NOT call the getter under test. A
			// control routed through the code under test cannot bound itself:
			// under MU4a-e the getter ignores its context, so such a control
			// blocks forever in the MAIN goroutine where no watchdog can reach
			// it, and the arm dies by global test timeout (measured: 180s)
			// instead of by its own named assertion.
			ctlCtx, ctlCancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			controlStart := time.Now()
			var one int
			ctlErr := s.db.QueryRowContext(ctlCtx, "SELECT 1").Scan(&one)
			ctlCancel()
			if !errors.Is(ctlErr, context.DeadlineExceeded) {
				t.Fatalf("control: with the sole connection held, a plain SELECT 1 under a "+
					"live 150ms context returned err=%v after %v; the connection is not "+
					"actually occupied, so the expired-context arm below would pass "+
					"vacuously", ctlErr, time.Since(controlStart))
			}

			expired, cancelExpired := context.WithDeadline(
				context.Background(), time.Now().Add(-time.Hour))
			defer cancelExpired()

			done := make(chan error, 1)
			go func() { done <- tc.call(expired, s) }()

			select {
			case err := <-done:
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("%s with an already-expired context returned err=%v, "+
						"want an error wrapping context.DeadlineExceeded", tc.name, err)
				}
			case <-time.After(watchdogBound):
				t.Errorf("%s did not return within %v with an already-expired context: "+
					"the context does not reach the blocking database/sql call",
					tc.name, watchdogBound)
				// RELEASE, then drain — see the R2 note above.
				release()
				select {
				case <-done:
				case <-time.After(watchdogBound):
					t.Errorf("%s did not unblock within %v even after the occupying "+
						"connection was released", tc.name, watchdogBound)
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// T1.6 — the lock policy exists on the LIVE CONNECTION, not just in a Go literal.
// -----------------------------------------------------------------------------

// wantBusyTimeoutMS is the busy_timeout policy value the tests below expect on
// a live production connection. It is a TEST-LOCAL LITERAL on purpose.
//
// Comparing the PRAGMA readback against the production constant instead would
// be a tautology: mutating that constant moves both sides of the comparison at
// once and the assertion can never fail. MU6 (2000 -> 0) SURVIVED that shape
// when it was first written, and was only killed once the expected value stopped
// being the thing under test.
const wantBusyTimeoutMS = 2000

// TestProductionDSNSetsBusyTimeout asserts busy_timeout by reading it back from
// the driver with `PRAGMA busy_timeout` on both production handles.
//
// The readback is the point. Asserting the Go constant, or grepping the DSN
// string, would pass under a driver that silently ignored the parameter and
// under an injection that never reached the connection. Reading driver state is
// the only form of this assertion that can fail for the right reason, which is
// what makes MU6 (value 2000 -> 0) a real kill rather than a tautology.
func TestProductionDSNSetsBusyTimeout(t *testing.T) {
	path := filepath.Join(t.TempDir(), "busytimeout.db")

	write, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%q): %v", path, err)
	}
	var writeTimeoutMS int
	if err := write.db.QueryRow("PRAGMA busy_timeout").Scan(&writeTimeoutMS); err != nil {
		_ = write.Close()
		t.Fatalf("read back PRAGMA busy_timeout on the write handle: %v", err)
	}
	// Close before opening read-only: Open holds the cross-process writer lock.
	if err := write.Close(); err != nil {
		t.Fatalf("close write handle: %v", err)
	}
	if writeTimeoutMS != wantBusyTimeoutMS {
		t.Errorf("write handle PRAGMA busy_timeout = %d, want %d — the lock policy "+
			"did not reach the live connection", writeTimeoutMS, wantBusyTimeoutMS)
	}

	readOnly, err := OpenReadOnly(path)
	if err != nil {
		t.Fatalf("OpenReadOnly(%q): %v", path, err)
	}
	defer func() { _ = readOnly.Close() }()
	var readTimeoutMS int
	if err := readOnly.db.QueryRow("PRAGMA busy_timeout").Scan(&readTimeoutMS); err != nil {
		t.Fatalf("read back PRAGMA busy_timeout on the read-only handle: %v", err)
	}
	if readTimeoutMS != wantBusyTimeoutMS {
		t.Errorf("read-only handle PRAGMA busy_timeout = %d, want %d", readTimeoutMS, wantBusyTimeoutMS)
	}

	// The production constant is pinned separately and explicitly, so a
	// deliberate policy change reds in exactly one obvious place.
	if busyTimeoutMillis != wantBusyTimeoutMS {
		t.Errorf("busyTimeoutMillis = %d, pinned at %d", busyTimeoutMillis, wantBusyTimeoutMS)
	}

	// A caller who set their own busy_timeout keeps it: the injection fills a
	// gap, it never overrides an explicit choice. The DSN must be a file: URI —
	// resolveDSN parses query parameters ONLY on the URI spelling, and a plain
	// path DSN silently discards everything after the "?".
	canonical, params, err := resolveDSN("file:" + path + "?_pragma=busy_timeout(1234)")
	if err != nil {
		t.Fatalf("resolveDSN with an explicit busy_timeout: %v", err)
	}
	if len(params["_pragma"]) == 0 {
		t.Fatalf("control: resolveDSN returned no _pragma parameters for an explicit "+
			"busy_timeout DSN (%v) — the override arm below would pass vacuously", params)
	}
	dsn := writeDSN(canonical, params)
	if !strings.Contains(dsn, "busy_timeout%281234%29") {
		t.Errorf("writeDSN dropped the caller's explicit busy_timeout: %q", dsn)
	}
	if strings.Contains(dsn, fmt.Sprintf("busy_timeout%%28%d%%29", wantBusyTimeoutMS)) {
		t.Errorf("writeDSN overrode the caller's explicit busy_timeout: %q", dsn)
	}
}

// -----------------------------------------------------------------------------
// T1.7 — the lock policy rides out a TRANSIENT exclusive lock.
// -----------------------------------------------------------------------------

// TestReadRetriesUnderTransientExclusiveLock puts a real SQLite EXCLUSIVE
// transaction in the read's way from a second raw driver connection, releases
// it at 300ms, and bounds the read at 3s while its own context allows 5s.
//
// PRE-REGISTERED OUTCOMES (the design deliberately does not claim which the
// driver produces, and BOTH satisfy the bound):
//
//	(a) retry-wins    — busy_timeout's retry loop outlasts the lock and the read
//	                    returns the row at roughly 300ms.
//	(b) interrupt-wins — the driver's busy sleep is context-interruptible and the
//	                    read returns a context/interrupt error before 2000ms.
//
// MEASURED ON THIS RIG (darwin, modernc.org/sqlite v1.54.0, journal_mode=delete):
// **(a) retry-wins** — the read returned the row at ~341ms in the planner's
// probe and again in this test. The assertion accepts either outcome and names
// which one it saw in its log line, so a future driver change that flips the
// answer is recorded rather than silently absorbed.
//
// The arm that makes this a gate rather than a stopwatch is MU5: with the
// busy_timeout injection deleted the read fails INSTANTLY with SQLITE_BUSY —
// under 3s, but with neither the row nor a context error.
func TestReadRetriesUnderTransientExclusiveLock(t *testing.T) {
	const (
		lockHeldFor = 300 * time.Millisecond
		readBound   = 3 * time.Second
		readCtxFor  = 5 * time.Second
	)

	path := filepath.Join(t.TempDir(), "exclusive.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%q): %v", path, err)
	}
	defer func() { _ = s.Close() }()

	seeded := obj("transient exclusive lock payload", "state/v1")
	if err := s.PutObject(seeded); err != nil {
		t.Fatalf("seed object: %v", err)
	}

	// A SECOND raw driver connection, not a second Store: store.Open takes the
	// non-waiting cross-process writer lock and would refuse, and the lock we
	// want here is SQLite's, not the store's.
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("raw driver open: %v", err)
	}
	defer func() { _ = raw.Close() }()
	locker, err := raw.Conn(context.Background())
	if err != nil {
		t.Fatalf("raw driver conn: %v", err)
	}
	if _, err := locker.ExecContext(context.Background(), "BEGIN EXCLUSIVE"); err != nil {
		_ = locker.Close()
		t.Fatalf("BEGIN EXCLUSIVE: %v", err)
	}

	// Known-positive control: the lock really does exclude this store's reads.
	// Without it, a driver or journal-mode change that stopped EXCLUSIVE from
	// blocking readers would make the timing assertion below pass vacuously.
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if _, _, err := s.GetObject(probeCtx, seeded.Hash); err == nil {
		probeCancel()
		_ = locker.Close()
		t.Fatalf("control: the read succeeded while a BEGIN EXCLUSIVE transaction was " +
			"open, so the lock does not exclude readers and this arm proves nothing")
	}
	probeCancel()

	released := make(chan struct{})
	go func() {
		defer close(released)
		time.Sleep(lockHeldFor)
		_, _ = locker.ExecContext(context.Background(), "ROLLBACK")
		_ = locker.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), readCtxFor)
	defer cancel()
	start := time.Now()
	got, ok, err := s.GetObject(ctx, seeded.Hash)
	elapsed := time.Since(start)
	<-released

	if elapsed > readBound {
		t.Fatalf("read under a transient exclusive lock took %v, want <= %v", elapsed, readBound)
	}

	// Outcome (b) is an INTERRUPT, and the discriminator against MU5 is the
	// error's KIND, not its timing: a busy/locked failure is also "an error
	// before 2000ms", so a purely temporal (b) clause would ACCEPT the very
	// mutant this arm exists to kill.
	interrupted := errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) ||
		strings.Contains(strings.ToLower(fmt.Sprint(err)), "interrupt")
	lockRefused := err != nil && !interrupted &&
		(strings.Contains(strings.ToLower(fmt.Sprint(err)), "busy") ||
			strings.Contains(strings.ToLower(fmt.Sprint(err)), "locked"))

	switch {
	case err == nil && ok && got.Hash.String() == seeded.Hash.String():
		t.Logf("PRE-REGISTERED OUTCOME (a) retry-wins: the row came back at %v "+
			"(lock held %v, busy_timeout %dms)", elapsed, lockHeldFor, wantBusyTimeoutMS)
	case interrupted && elapsed < time.Duration(wantBusyTimeoutMS)*time.Millisecond:
		t.Logf("PRE-REGISTERED OUTCOME (b) interrupt-wins: the read returned "+
			"context/interrupt err=%v at %v, before the %dms busy_timeout",
			err, elapsed, wantBusyTimeoutMS)
	case lockRefused:
		t.Fatalf("the read was REFUSED by the lock (err=%v) after %v instead of riding "+
			"it out: this is the MU5 shape — the busy_timeout injection in "+
			"writeDSN/readOnlyDSN is gone, so the driver fails at the first "+
			"SQLITE_BUSY instead of retrying for %dms", err, elapsed, wantBusyTimeoutMS)
	default:
		t.Fatalf("read under a transient exclusive lock matched NEITHER pre-registered "+
			"outcome: elapsed=%v ok=%v err=%v. Outcome (a) wants the row back; outcome "+
			"(b) wants a context/interrupt error before %dms", elapsed, ok, err, wantBusyTimeoutMS)
	}
}

// -----------------------------------------------------------------------------
// T1.8 — the §2.8 ratchet.
// -----------------------------------------------------------------------------

// deadlineFreeReadPins is the EXACT set of production store reads this item
// leaves deadline-free, by ratified deferral DR-2. Each is today's behaviour
// made visible at the call site rather than hidden in a context-free signature.
//
// The set may SHRINK — threading a real context through one of these sites is a
// one-line edit here in the same diff — but it may never GROW: a new
// deadline-free store read anywhere under host/ or cmd/ reds this test. That is
// what makes the follow-on item's progress mechanically observable, 11 -> 0, and
// what lets the store-boundary reject land exactly when this reads zero.
var deadlineFreeReadPins = map[string]int{
	"host/broker/approve.go":    8,
	"host/registry/registry.go": 2,
	"host/replay/replay.go":     1,
}

// deadlineFreeReadCall matches a call of one of the five context-first read
// getters whose context argument is the deadline-free literal.
var deadlineFreeReadCall = regexp.MustCompile(
	`\.(GetObject|GetWorld|GetLogEntry|GetRegistryHead|SelectedHead)\(\s*context\.Background\(\)`)

// TestNoNewDeadlineFreeStoreReads pins the deadline-free residue.
//
// It scans the PRODUCTION (non-_test) Go sources under host/ and cmd/, rooted
// via runtime.Caller so the result does not depend on the test's working
// directory, and requires every file's count to equal its pin — zero for every
// file not named above.
func TestNoNewDeadlineFreeStoreReads(t *testing.T) {
	root := repoRootFromCaller(t)

	// Non-vacuity guard 1: every pinned file must exist. A rename that made the
	// scan miss a pinned file would otherwise read as "count 0, pin 0 — fine"
	// only because the pin was never consulted.
	for rel := range deadlineFreeReadPins {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("pinned file %s is not at %s: %v — the pin is stale, not satisfied", rel, root, err)
		}
	}

	counts := map[string]int{}
	scanned := 0
	for _, top := range []string{"host", "cmd"} {
		base := filepath.Join(root, top)
		err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			src, rerr := os.ReadFile(path)
			if rerr != nil {
				return rerr
			}
			scanned++
			rel, rerr := filepath.Rel(root, path)
			if rerr != nil {
				return rerr
			}
			if n := len(deadlineFreeReadCall.FindAllIndex(src, -1)); n > 0 {
				counts[filepath.ToSlash(rel)] = n
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", base, err)
		}
	}

	// Non-vacuity guard 2: the walk must actually have read files. A rooting
	// bug that scanned an empty tree would report every count as 0 and every
	// "all else: 0" clause as satisfied.
	if scanned < 20 {
		t.Fatalf("scanned only %d production .go files under %s — the scan is not "+
			"reaching the tree it claims to pin", scanned, root)
	}

	total := 0
	for rel, want := range deadlineFreeReadPins {
		got := counts[rel]
		total += got
		if got != want {
			t.Errorf("%s has %d deadline-free store read(s), pinned at %d. "+
				"Fewer: delete or lower this pin in the SAME diff that threaded the "+
				"context. More: a new unbounded store read landed — thread a real "+
				"context instead of raising the pin", rel, got, want)
		}
	}
	for rel, got := range counts {
		if _, pinned := deadlineFreeReadPins[rel]; pinned {
			continue
		}
		total += got
		t.Errorf("%s has %d deadline-free store read(s) and is pinned at 0: a new "+
			"unbounded store read landed outside the ratified DR-2 residue", rel, got)
	}

	wantTotal := 0
	for _, n := range deadlineFreeReadPins {
		wantTotal += n
	}
	if total != wantTotal {
		t.Errorf("deadline-free store reads total %d, pinned at %d", total, wantTotal)
	}
	t.Logf("ratchet: %d deadline-free production store reads across %d scanned files "+
		"(pins %v)", total, scanned, deadlineFreeReadPins)
}

// repoRootFromCaller locates the repository root from this source file's own
// path, so the scan is independent of the test process's working directory.
func repoRootFromCaller(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed: cannot root the source scan")
	}
	// <root>/host/store/context_read_test.go -> <root>
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("rooted at %s but there is no go.mod there: %v", root, err)
	}
	return root
}
