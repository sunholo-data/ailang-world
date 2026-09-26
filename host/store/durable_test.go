package store

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// durableBudget is the caller deadline these tests hold the store past;
// durableEpsilon is the scheduling allowance on top of it.
const (
	durableBudget  = 100 * time.Millisecond
	durableEpsilon = 250 * time.Millisecond
)

func openFileStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// settledGoroutines polls until the goroutine count is at most base, or 2 s.
func settledGoroutines(base int) int {
	deadline := time.Now().Add(2 * time.Second)
	n := runtime.NumGoroutine()
	for n > base && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		n = runtime.NumGoroutine()
	}
	return n
}

// boundedCall runs f; if it has not returned within durableBudget+1s, unblock
// runs (so a ctx-ignoring mutant fails instead of hanging) and the test fails.
func boundedCall(t *testing.T, f func() error, unblock func()) (error, time.Duration) {
	t.Helper()
	start := time.Now()
	done := make(chan error, 1)
	go func() { done <- f() }()
	select {
	case err := <-done:
		return err, time.Since(start)
	case <-time.After(durableBudget + time.Second):
		unblock()
		err := <-done
		t.Fatalf("still blocked %v after a %v deadline (err after unblock: %v)", time.Since(start), durableBudget, err)
		return err, 0
	}
}

// TestDurableOpsHonourHeldConnection is the row-106 M5 shape: the sole pooled
// connection is held past the caller deadline; every durable operation returns
// the definite not-committed ctx error within budget+epsilon, leaves no
// goroutine behind, and leaves the store untouched.
func TestDurableOpsHonourHeldConnection(t *testing.T) {
	s := openFileStore(t)
	c := journalCommitFixture(t, s, "held")
	if _, _, err := s.AppendIntent("held", testCommitIntent("held", c)); err != nil {
		t.Fatalf("seed intent: %v", err)
	}
	ops := map[string]func(ctx context.Context) error{
		"CommitContext": func(ctx context.Context) error { return s.CommitContext(ctx, c) },
		"AppendIntentContext": func(ctx context.Context) error {
			_, _, err := s.AppendIntentContext(ctx, "held-2", testCommitIntent("held-2", c))
			return err
		},
		"GetReceiptContext": func(ctx context.Context) error {
			_, _, err := s.GetReceiptContext(ctx, "held")
			return err
		},
	}
	for name, op := range ops {
		t.Run(name, func(t *testing.T) {
			base := runtime.NumGoroutine()
			conn, err := s.db.Conn(context.Background())
			if err != nil {
				t.Fatalf("hold conn: %v", err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
			err, elapsed := boundedCall(t, func() error { return op(ctx) }, func() { _ = conn.Close() })
			cancel()
			_ = conn.Close()
			if !errors.Is(err, context.DeadlineExceeded) || IsUncertain(err) {
				t.Fatalf("err = %v; want definite DeadlineExceeded", err)
			}
			if elapsed > durableBudget+durableEpsilon {
				t.Fatalf("returned after %v; bound %v", elapsed, durableBudget+durableEpsilon)
			}
			if n := settledGoroutines(base); n > base {
				t.Fatalf("goroutines %d > baseline %d: a blocked worker accumulated", n, base)
			}
		})
	}
	rc, _, err := s.GetReceiptContext(context.Background(), "held")
	if err != nil || rc.State != ReceiptIndeterminate {
		t.Fatalf("receipt after refused ops = %+v, %v; want indeterminate (intent only)", rc, err)
	}
	if rc2, _, _ := s.GetReceipt("held-2"); rc2.State != ReceiptNotStarted {
		t.Fatalf("held-2 intent landed: %+v", rc2)
	}
	if err := s.CommitContext(context.Background(), c); err != nil {
		t.Fatalf("control: released store must commit: %v", err)
	}
}

// TestCommitCancelledMidBodyIsAllOrNothing cancels between the head/world/log
// writes and the outcome write: the result is definite not-committed and none
// of the transaction's writes are visible.
func TestCommitCancelledMidBodyIsAllOrNothing(t *testing.T) {
	s := openFileStore(t)
	c := journalCommitFixture(t, s, "mid")
	if _, _, err := s.AppendIntent("mid", testCommitIntent("mid", c)); err != nil {
		t.Fatalf("seed intent: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	prev := commitBeforeOutcomeHook
	// Cancel, then yield so database/sql's asynchronous rollback (awaitDone)
	// lands before the next ctx-free statement: that statement then reports
	// sql.ErrTxDone, the arm notCommitted must translate.
	commitBeforeOutcomeHook = func() { cancel(); time.Sleep(50 * time.Millisecond) }
	defer func() { commitBeforeOutcomeHook = prev }()
	err := s.CommitContext(ctx, c)
	if !errors.Is(err, context.Canceled) || IsUncertain(err) {
		t.Fatalf("err = %v; want definite Canceled", err)
	}
	head, _, herr := s.SelectedHead(context.Background())
	if herr != nil || head.String() != c.ObservedHead.String() {
		t.Fatalf("head moved to %v (%v); want untouched %v", head, herr, c.ObservedHead)
	}
	if _, ok, _ := s.GetWorld(context.Background(), c.NextWorld.Ref); ok {
		t.Fatal("next world row landed from a cancelled commit")
	}
	if rc, _, _ := s.GetReceipt("mid"); rc.State != ReceiptIndeterminate {
		t.Fatalf("receipt = %v; want indeterminate", rc.State)
	}
}

func TestCommitBodyHookPhases(t *testing.T) {
	s := openFileStore(t)
	c := journalCommitFixture(t, s, "body-phases")
	if _, _, err := s.AppendIntent("body-phases", testCommitIntent("body-phases", c)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	previous := commitBodyHook
	visits := 0
	commitBodyHook = func(got context.Context) {
		if got != ctx {
			t.Fatal("body hook did not receive caller context")
		}
		visits++
		if visits == 2 {
			cancel()
			time.Sleep(50 * time.Millisecond)
		}
	}
	defer func() { commitBodyHook = previous }()
	err := s.CommitContext(ctx, c)
	if visits != 2 {
		t.Fatalf("body hook visits = %d; want head read and world insert", visits)
	}
	if !errors.Is(err, context.Canceled) || IsUncertain(err) {
		t.Fatalf("err = %v; want definite Canceled", err)
	}
	if _, ok, err := s.GetWorld(context.Background(), c.NextWorld.Ref); err != nil || ok {
		t.Fatalf("cancelled world: found=%v err=%v", ok, err)
	}
	if rc, _, err := s.GetReceipt("body-phases"); err != nil || rc.State != ReceiptIndeterminate {
		t.Fatalf("receipt = %v, err = %v; want indeterminate", rc.State, err)
	}
}

func TestAppendIntentCancelledMidBodyIsNotCommitted(t *testing.T) {
	s := openFileStore(t)
	c := journalCommitFixture(t, s, "append-mid-body")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	previous := appendIntentBodyHook
	visits := 0
	appendIntentBodyHook = func(got context.Context) {
		if got != ctx {
			t.Fatal("append hook did not receive caller context")
		}
		visits++
		cancel()
		time.Sleep(50 * time.Millisecond)
	}
	defer func() { appendIntentBodyHook = previous }()
	_, _, err := s.AppendIntentContext(ctx, "append-mid-body", testCommitIntent("append-mid-body", c))
	if visits != 1 {
		t.Fatalf("append hook visits = %d; want one", visits)
	}
	if !errors.Is(err, context.Canceled) || IsUncertain(err) {
		t.Fatalf("err = %v; want definite Canceled", err)
	}
	rc, hasIntent, err := s.GetReceipt("append-mid-body")
	if err != nil || hasIntent || rc.State != ReceiptNotStarted {
		t.Fatalf("receipt = %v, hasIntent = %v, err = %v; want not started with no intent", rc.State, hasIntent, err)
	}
}

// TestDurableStallAtCutoffIsUncertainAndReconciles stalls the durable step of
// each durable write past the caller deadline: the caller gets *UncertainError
// (not a ctx error) within budget+epsilon, the finishing goroutine exits once
// the stall releases, and a later receipt read settles the outcome.
func TestDurableStallAtCutoffIsUncertainAndReconciles(t *testing.T) {
	cases := []struct {
		name string
		id   string // the receipt read that reconciles
		op   func(ctx context.Context, s *Store, c Commit) error
		want ReceiptState // deadline passed before the cutoff: rolled back
	}{
		{"CommitContext", "stall", func(ctx context.Context, s *Store, c Commit) error {
			return s.CommitContext(ctx, c)
		}, ReceiptIndeterminate},
		{"AppendIntentContext", "stall-2", func(ctx context.Context, s *Store, c Commit) error {
			_, _, err := s.AppendIntentContext(ctx, "stall-2", testCommitIntent("stall-2", c))
			return err
		}, ReceiptNotStarted},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openFileStore(t)
			c := journalCommitFixture(t, s, "stall")
			if _, _, err := s.AppendIntent("stall", testCommitIntent("stall", c)); err != nil {
				t.Fatalf("seed intent: %v", err)
			}
			release := make(chan struct{})
			// Release on every exit: a stalled worker would otherwise block
			// the Cleanup Close (which waits for durable workers) forever.
			defer func() {
				select {
				case <-release:
				default:
					close(release)
				}
			}()
			prev := durableCommitHook
			durableCommitHook = func() { <-release }
			defer func() { durableCommitHook = prev }()
			base := runtime.NumGoroutine()
			ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
			defer cancel()
			err, elapsed := boundedCall(t, func() error { return tc.op(ctx, s, c) }, func() { close(release) })
			if !IsUncertain(err) || errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("err = %v; want *UncertainError that is not a ctx error", err)
			}
			if elapsed > durableBudget+durableEpsilon {
				t.Fatalf("returned after %v; bound %v", elapsed, durableBudget+durableEpsilon)
			}
			select {
			case <-release:
			default:
				close(release)
			}
			if n := settledGoroutines(base); n > base {
				t.Fatalf("goroutines %d > baseline %d after release", n, base)
			}
			rctx, rcancel := context.WithTimeout(context.Background(), time.Second)
			defer rcancel()
			rc, _, err := s.GetReceiptContext(rctx, tc.id)
			if err != nil || rc.State != tc.want {
				t.Fatalf("reconciled %s = %v, %v; want %v (not committed)", tc.id, rc.State, err, tc.want)
			}
		})
	}
}

// bareCommit builds a commit with no InvocationID (the /v1/commit shape) that
// appends log entry index on top of head.
func bareCommit(head World, index int64, tag string) Commit {
	body := obj("bare-body-"+tag, "transition/body")
	entryHash := hashref.SumSHA256([]byte("bare-entry-" + tag))
	return Commit{
		ObservedHead: head.Ref, Objects: []Object{body},
		NextWorld: World{Ref: hashref.SumSHA256([]byte("bare-world-" + tag)),
			Revision: head.Revision + 1, StateRoot: body.Hash, LogHead: entryHash},
		Entry: LogEntry{Header: LogHeader{EntryIndex: index, SemanticsEpoch: 1,
			TransitionFn: body.Hash, Interpreter: body.Hash, PrevEntryHash: head.LogHead,
			WrittenBy: "bare-test"}, EntryHash: entryHash, TransitionRef: body.Hash},
	}
}

// logWitness is the bare-commit reconciliation rule (M2 promotes it to
// production): the log row at c's entry index is written only by Commit and
// never updated or deleted, and entry_index is the primary key. Absent or
// different ⇒ c did not commit; equal in every stored field ⇒ a commit
// carrying c's exact log row is durable. The selected head is NOT consulted.
func logWitness(ctx context.Context, s *Store, c Commit) (landed bool, err error) {
	got, ok, err := s.GetLogEntry(ctx, c.Entry.Header.EntryIndex)
	if err != nil || !ok {
		return false, err
	}
	return reflect.DeepEqual(got, c.Entry), nil
}

// TestUncertainBareCommitReconcilesByLogIndexNotHead is quorum r1 objection 1's
// test: commit A reports *UncertainError, commit B lands before A is
// reconciled, and reconciliation must still classify A correctly. A head
// comparison would say "A failed" in the landed arm; the log witness does not.
func TestUncertainBareCommitReconcilesByLogIndexNotHead(t *testing.T) {
	for _, tc := range []struct {
		name       string
		stallAfter bool // true: A lands, then the caller times out
		wantLanded bool
	}{
		{"A_landed_then_B_on_top", true, true},
		{"A_rolled_back_then_B_takes_the_index", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := openFileStore(t)
			genesis := seedGenesis(t, s)
			a := bareCommit(genesis, 1, "A")
			release := make(chan struct{})
			defer func() {
				select {
				case <-release:
				default:
					close(release)
				}
			}()
			hook := &durableCommitHook
			if tc.stallAfter {
				hook = &durableAfterCommitHook
			}
			prev := *hook
			*hook = func() { <-release }
			defer func() { *hook = prev }()
			ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
			defer cancel()
			err, _ := boundedCall(t, func() error { return s.CommitContext(ctx, a) }, func() { close(release) })
			if !IsUncertain(err) {
				t.Fatalf("A: err = %v; want *UncertainError", err)
			}
			select {
			case <-release:
			default:
				close(release)
			}
			// B is planned against whatever head is now selected (A's world if
			// A landed, genesis otherwise), so it commits either way.
			bctx, bcancel := context.WithTimeout(context.Background(), time.Second)
			defer bcancel()
			base := genesis
			if tc.stallAfter {
				base = a.NextWorld
			}
			b := bareCommit(base, base.Revision+1, "B")
			if err := s.CommitContext(bctx, b); err != nil {
				t.Fatalf("B: %v", err)
			}
			head, _, err := s.SelectedHead(bctx)
			if err != nil || head.String() == a.NextWorld.Ref.String() {
				t.Fatalf("head = %v (%v); the scenario needs A's world NOT to be the head", head, err)
			}
			landed, err := logWitness(bctx, s, a)
			if err != nil || landed != tc.wantLanded {
				t.Fatalf("logWitness(A) = %v, %v; want %v", landed, err, tc.wantLanded)
			}
		})
	}
}

// TestCloseKeepsWriterLockUntilDurableWorkerSettles is quorum r2's lifecycle
// test at the store layer: commit A lands and its worker stalls after COMMIT
// past A's deadline (A reports *UncertainError). Close must not release the
// writer lock while that worker is outstanding: a second writer Open is
// refused while the worker is stalled, and succeeds only after it settles.
func TestCloseKeepsWriterLockUntilDurableWorkerSettles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	genesis := seedGenesis(t, s)
	release := make(chan struct{})
	prev := durableAfterCommitHook
	durableAfterCommitHook = func() { <-release }
	defer func() { durableAfterCommitHook = prev }()
	unblock := func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}
	defer unblock()
	ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
	defer cancel()
	err, _ = boundedCall(t, func() error { return s.CommitContext(ctx, bareCommit(genesis, 1, "L")) }, unblock)
	if !IsUncertain(err) {
		t.Fatalf("A: err = %v; want *UncertainError", err)
	}
	closed := make(chan error, 1)
	go func() { closed <- s.Close() }()
	time.Sleep(durableBudget) // give an early-releasing Close time to act
	if other, err := Open(path); err == nil {
		_ = other.Close()
		t.Fatal("second writer Open succeeded while a durable worker was outstanding: writer lock released early")
	} else if !IsWriterAlreadyActive(err) {
		t.Fatalf("second Open: %v; want the writer-already-active refusal", err)
	}
	unblock()
	select {
	case err := <-closed:
		if err != nil {
			t.Fatalf("Close after release: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return after the worker settled")
	}
	again, err := Open(path)
	if err != nil {
		t.Fatalf("re-Open after release: %v", err)
	}
	defer again.Close()
	if landed, err := logWitness(ctx2s(t), again, bareCommit(genesis, 1, "L")); err != nil || !landed {
		t.Fatalf("A's entry after cleanup: landed=%v err=%v; want landed", landed, err)
	}
}

func ctx2s(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)
	return ctx
}
