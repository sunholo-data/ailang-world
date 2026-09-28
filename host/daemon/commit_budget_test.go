package daemon

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// commitSeam wraps the daemon's REAL store behind the durableStore seam and
// substitutes only the outcome a real store cannot produce on demand. Every
// handler line still runs the production path.
type commitSeam struct {
	real *store.Store
	mode string // "block", "land-then-uncertain", "uncertain-not-landed"
}

// seamEscape bounds "block": a ctx-ignoring mutant is released after it and
// then commits for real, so it is RED on status and elapsed time, not a hang.
const seamEscape = 2 * time.Second

func (f *commitSeam) Commit(ctx context.Context, c store.Commit) error {
	// waitCtx models the caller's deadline ending. Every real uncertain
	// outcome carries an ENDED ctx (the caller stopped waiting), so the fake
	// must too; returning early would let a handler that consults timedOut
	// before IsUncertain pass (MUT-M2-TIMEOUT-BEFORE-UNCERTAIN).
	waitCtx := func() {
		select {
		case <-ctx.Done():
		case <-time.After(seamEscape):
		}
	}
	switch f.mode {
	case "block": // a held sole connection: nothing happens until ctx ends
		waitCtx()
		return f.real.Commit(ctx, c)
	case "land-then-uncertain": // COMMIT completed, the caller's deadline won the race
		lctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := f.real.Commit(lctx, c); err != nil {
			return err
		}
		waitCtx()
		return &store.UncertainError{Op: "commit", Cause: context.DeadlineExceeded}
	default: // "uncertain-not-landed": rolled back, the caller saw only the expiry
		waitCtx()
		return &store.UncertainError{Op: "commit", Cause: context.DeadlineExceeded}
	}
}

// AppendIntent passes through to the real store; "block" models the
// held connection for the intent step too.
func (f *commitSeam) AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error) {
	if f.mode == "block" {
		select {
		case <-ctx.Done():
		case <-time.After(seamEscape):
		}
	}
	return f.real.AppendIntent(ctx, id, intent)
}

func postCommitRec(t *testing.T, d *Daemon, auth string, c store.Commit) (*httptest.ResponseRecorder, time.Duration) {
	t.Helper()
	start := time.Now()
	rec := requestRecorderAuth(t, d, auth, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(c)))
	return rec, time.Since(start)
}

func headOf(t *testing.T, d *Daemon) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	head, _, err := d.store.SelectedHead(ctx)
	if err != nil {
		t.Fatalf("SelectedHead: %v", err)
	}
	return head.String()
}

type unusedCoordinatorRunner struct{}

func (unusedCoordinatorRunner) RunContext(context.Context, capsule.Entry) (capsule.Result, error) {
	panic("coordinator runner reached before receipt")
}

// TestCommitBudgetAndUncertainReconcile is AC5 (row 23 M2): /v1/commit is
// bounded by B6, the two 503s are distinct classes, and an uncertain commit
// reconciles by the log row at its index, never by the head.
func TestCommitBudgetAndUncertainReconcile(t *testing.T) {
	t.Run("coordinator_held_connection_honors_caller_ctx", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "world.db")
		st, err := store.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer st.Close()
		blocker, err := sql.Open("sqlite", path)
		if err != nil {
			t.Fatal(err)
		}
		defer blocker.Close()
		conn, err := blocker.Conn(boundedTestContext(t))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if _, err := conn.ExecContext(boundedTestContext(t), "BEGIN EXCLUSIVE"); err != nil {
			t.Fatal(err)
		}
		released := make(chan struct{})
		go func() {
			time.Sleep(700 * time.Millisecond)
			_, _ = conn.ExecContext(boundedTestContext(t), "ROLLBACK")
			close(released)
		}()
		readDone := make(chan struct{})
		go func() {
			_, _, _ = st.GetObject(boundedTestContext(t), hashref.SumSHA256([]byte("held")))
			close(readDone)
		}()
		time.Sleep(30 * time.Millisecond)
		c, err := coordinator.New(coordinator.Config{Store: st, Runner: unusedCoordinatorRunner{}, Binder: func(string, []broker.Capability) transitionreg.Binder { return nil }, Now: func() int64 { return 1 }, MaxInput: 1024, MaxOutput: 1024})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err = c.Dispatch(ctx, coordinator.Call{EpisodeID: "ep", TaskID: "held", Input: map[string]any{}})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("coordinator held receipt result=%v, want deadline", err)
		}
		if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
			t.Errorf("coordinator held receipt took %v, want caller bound", elapsed)
		}
		<-released
		select {
		case <-readDone:
		case <-time.After(3 * time.Second):
			t.Fatal("held read did not settle")
		}
	})
	t.Run("i_held_connection_is_503_timeout_within_budget", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "held")
		auth := authHeader(t, d)
		d.commits = &commitSeam{real: d.store, mode: "block"}
		d.commitBudget = 100 * time.Millisecond
		rec, elapsed := postCommitRec(t, d, auth, testCommit(genesis, 1, "held"))
		body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
		if !strings.Contains(body.Error.Message, "not committed") {
			t.Fatalf("timeout message = %q; want it to say not committed", body.Error.Message)
		}
		if elapsed > d.commitBudget+500*time.Millisecond {
			t.Fatalf("answered after %v; bound %v + 500ms", elapsed, d.commitBudget)
		}
		if got := headOf(t, d); got != genesis.Ref.String() {
			t.Fatalf("head moved to %s after a not-committed 503", got)
		}
	})

	for _, tc := range []struct {
		name       string
		mode       string
		bOnTop     bool // B planned on A's world (A landed) vs on genesis
		skipB      bool
		wantLanded bool
	}{
		{"ii_A_landed_uncertain_then_B_on_top", "land-then-uncertain", true, false, true},
		{"iii_A_rolled_back_then_B_takes_the_index", "uncertain-not-landed", false, false, false},
		{"iii_b_A_rolled_back_and_nothing_else", "uncertain-not-landed", false, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := newHandlerDaemon(t)
			genesis := seedGenesisEmbedded(t, d, "unc")
			auth := authHeader(t, d)
			a := testCommit(genesis, 1, "A")
			d.commits = &commitSeam{real: d.store, mode: tc.mode}
			d.commitBudget = 100 * time.Millisecond
			rec, _ := postCommitRec(t, d, auth, a)
			body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "CommitUncertain")
			if !strings.Contains(body.Error.Message, "/v1/log/1") {
				t.Fatalf("uncertain message = %q; want the reconcile route GET /v1/log/1", body.Error.Message)
			}
			d.commits = d.store
			d.commitBudget = commitBudget // B is an ordinary commit: the ratified B6, not the arm's 100 ms
			if !tc.skipB {
				base := genesis
				if tc.bOnTop {
					base = a.NextWorld
				}
				b := testCommit(base, base.Revision+1, "B")
				if rec, _ := postCommitRec(t, d, auth, b); rec.Code != http.StatusOK {
					t.Fatalf("commit B: status %d body %s", rec.Code, rec.Body)
				}
				if got := headOf(t, d); got == a.NextWorld.Ref.String() {
					t.Fatalf("head = A's world; the scenario needs another head")
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			landed, err := d.store.CommitLanded(ctx, a)
			if err != nil || landed != tc.wantLanded {
				t.Fatalf("CommitLanded(A) = %v, %v; want %v", landed, err, tc.wantLanded)
			}
			// The client-side form QUICKSTART documents: GET /v1/log/{index}
			// and compare with the entry that was sent.
			got := requestRecorder(t, d, http.MethodGet, "/v1/log/1", nil)
			sent, _ := json.Marshal(logJSON(a.Entry))
			equal := got.Code == http.StatusOK && strings.TrimSpace(got.Body.String()) == string(sent)
			if equal != tc.wantLanded {
				t.Fatalf("wire witness: GET /v1/log/1 = %d %s; equal-to-sent=%v, want %v", got.Code, got.Body, equal, tc.wantLanded)
			}
		})
	}

	t.Run("iv_reconcile_read_timeout_gives_no_resend_advice", func(t *testing.T) {
		d := newHandlerDaemon(t)
		d.readDeadline = expiredReadDeadline
		rec := requestRecorder(t, d, http.MethodGet, "/v1/log/1", nil)
		body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
		if strings.Contains(strings.ToLower(body.Error.Message), "resend") {
			t.Fatalf("reconcile-read timeout carries resend advice: %q", body.Error.Message)
		}
	})
}

// TestCommitUncertainStatusMirrorsSketch: the new class's status is the frozen
// sketch's httpStatus arm, parsed rather than restated.
func TestCommitUncertainStatusMirrorsSketch(t *testing.T) {
	want, ok := sketchHTTPStatusVectors(t)["CommitUncertain"]
	if !ok {
		t.Fatalf("the sketch declares no CommitUncertain arm")
	}
	rec := httptest.NewRecorder()
	writeCommitUncertain(rec, store.Commit{})
	if rec.Code != want {
		t.Fatalf("CommitUncertain renders %d; the sketch says %d", rec.Code, want)
	}
}

func postCommitWithID(t *testing.T, d *Daemon, auth, id string, c store.Commit) (*httptest.ResponseRecorder, time.Duration) {
	t.Helper()
	var wire map[string]any
	if err := json.Unmarshal(encodeCommit(c), &wire); err != nil {
		t.Fatal(err)
	}
	wire["invocationId"] = id
	raw, _ := json.Marshal(wire)
	start := time.Now()
	rec := requestRecorderAuth(t, d, auth, http.MethodPost, "/v1/commit", bytes.NewReader(raw))
	return rec, time.Since(start)
}

func getReceipt(t *testing.T, d *Daemon, id string) receiptResponse {
	t.Helper()
	rec := requestRecorder(t, d, http.MethodGet, "/v1/receipts/"+id, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /v1/receipts/%s = %d %s", id, rec.Code, rec.Body)
	}
	var body receiptResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode receipt: %v; body %s", err, rec.Body)
	}
	return body
}

// TestCommitInvocationReceiptIdentity is AC5b (row 23 M3): with an opt-in
// "rest:" invocationId the receipt, not the head and not the log row, is the
// full identity witness of the commit, world included.
func TestCommitInvocationReceiptIdentity(t *testing.T) {
	for _, tc := range []struct {
		name      string
		mode      string
		wantState string
		wantWorld bool
	}{
		{"A_landed_uncertain_then_B_on_top", "land-then-uncertain", "resolved", true},
		{"A_rolled_back_uncertain", "uncertain-not-landed", "indeterminate", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := newHandlerDaemon(t)
			genesis := seedGenesisEmbedded(t, d, "rcpt")
			auth := authHeader(t, d)
			a := testCommit(genesis, 1, "A")
			d.commits = &commitSeam{real: d.store, mode: tc.mode}
			d.commitBudget = 100 * time.Millisecond
			rec, _ := postCommitWithID(t, d, auth, "rest:A", a)
			body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "CommitUncertain")
			if !strings.Contains(body.Error.Message, "/v1/receipts/rest:A") {
				t.Fatalf("uncertain message = %q; want the receipt route", body.Error.Message)
			}
			d.commits = d.store
			d.commitBudget = commitBudget
			if tc.wantWorld {
				b := testCommit(a.NextWorld, 2, "B")
				if rec, _ := postCommitRec(t, d, auth, b); rec.Code != http.StatusOK {
					t.Fatalf("commit B: %d %s", rec.Code, rec.Body)
				}
			}
			got := getReceipt(t, d, "rest:A")
			if got.State != tc.wantState {
				t.Fatalf("receipt state = %q, want %q", got.State, tc.wantState)
			}
			if tc.wantWorld && got.ResultRef != a.NextWorld.Ref.String() {
				t.Fatalf("receipt resultRef = %q, want A's world %s", got.ResultRef, a.NextWorld.Ref)
			}
		})
	}

	t.Run("held_connection_on_the_intent_step_is_503_timeout", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "rcpt-held")
		auth := authHeader(t, d)
		d.commits = &commitSeam{real: d.store, mode: "block"}
		d.commitBudget = 100 * time.Millisecond
		rec, elapsed := postCommitWithID(t, d, auth, "rest:held", testCommit(genesis, 1, "held"))
		assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
		if elapsed > d.commitBudget+500*time.Millisecond {
			t.Fatalf("answered after %v; bound %v + 500ms", elapsed, d.commitBudget)
		}
		if got := getReceipt(t, d, "rest:held"); got.State != "not-started" {
			t.Fatalf("receipt after a pre-cutoff intent timeout = %q, want not-started", got.State)
		}
	})

	t.Run("foreign_namespaces_are_refused", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "rcpt-ns")
		auth := authHeader(t, d)
		for _, id := range []string{"effect:ep:1", "a2a:ep:task", "rest:"} {
			rec, _ := postCommitWithID(t, d, auth, id, testCommit(genesis, 1, "ns"))
			body := assertErrorClass(t, rec, http.StatusBadRequest, "BadRequest")
			if !strings.Contains(body.Error.Message, "rest:<id>") {
				t.Fatalf("%s: message %q does not name the rest: namespace", id, body.Error.Message)
			}
			got := requestRecorder(t, d, http.MethodGet, "/v1/receipts/"+id, nil)
			assertErrorClass(t, got, http.StatusBadRequest, "BadRequest")
		}
		if got := headOf(t, d); got != genesis.Ref.String() {
			t.Fatalf("a refused namespace moved the head to %s", got)
		}
	})
}

// mismatchSeam returns the store's own *store.InvocationMismatchError from the
// intent step. Through the real handler that error cannot come from a real
// store: the handler derives the intent from the commit, so an id with
// different content always collides at the intent append first
// (DuplicateInvocationError). The seam pins the mapping of the second type.
type mismatchSeam struct{ real *store.Store }

func (f mismatchSeam) Commit(ctx context.Context, c store.Commit) error {
	return f.real.Commit(ctx, c)
}

func (f mismatchSeam) AppendIntent(_ context.Context, id string, _ store.JournalIntent) (int64, hashref.HashRef, error) {
	return 0, hashref.HashRef{}, &store.InvocationMismatchError{ID: id, Field: "WorldRef", Want: "sha256:want", Got: "sha256:got"}
}

// TestCommitInvocationReuseIsBadRequest pins writeCommitError's
// duplicate/mismatch arm (judge survivor MY-2): reusing a "rest:" id for a
// different commit is the client's own input error, so it is 400 BadRequest
// naming the reuse, not a sanitized 500, and it never moves the head.
func TestCommitInvocationReuseIsBadRequest(t *testing.T) {
	assertReuse := func(t *testing.T, rec *httptest.ResponseRecorder, errLog *bytes.Buffer, id, detail string) {
		t.Helper()
		body := assertErrorClass(t, rec, http.StatusBadRequest, "BadRequest")
		msg := body.Error.Message
		if !strings.HasPrefix(msg, "invocationId already names a different commit: ") || !strings.Contains(msg, id) || !strings.Contains(msg, detail) {
			t.Fatalf("reuse message = %q; want the reuse named with id %q and %q", msg, id, detail)
		}
		if strings.Contains(msg, "internal") || errLog.Len() != 0 {
			t.Fatalf("a client reuse error went down the internal path: message %q, error log %q", msg, errLog.String())
		}
	}

	t.Run("duplicate_id_different_content_real_store", func(t *testing.T) {
		d := newHandlerDaemon(t)
		var errLog bytes.Buffer
		d.errLog = &errLog
		genesis := seedGenesisEmbedded(t, d, "reuse")
		auth := authHeader(t, d)
		a := testCommit(genesis, 1, "A")
		// Live control: the first use of the id commits.
		if rec, _ := postCommitWithID(t, d, auth, "rest:reuse", a); rec.Code != http.StatusOK {
			t.Fatalf("first commit with rest:reuse: status %d body %s", rec.Code, rec.Body)
		}
		// Same id, different commit content (another commit on A's world).
		rec, _ := postCommitWithID(t, d, auth, "rest:reuse", testCommit(a.NextWorld, 2, "B"))
		assertReuse(t, rec, &errLog, "rest:reuse", "duplicate intent")
		if got := headOf(t, d); got != a.NextWorld.Ref.String() {
			t.Fatalf("a refused id reuse moved the head to %s", got)
		}
	})

	t.Run("invocation_mismatch_seam", func(t *testing.T) {
		d := newHandlerDaemon(t)
		var errLog bytes.Buffer
		d.errLog = &errLog
		genesis := seedGenesisEmbedded(t, d, "mismatch")
		auth := authHeader(t, d)
		d.commits = mismatchSeam{real: d.store}
		rec, _ := postCommitWithID(t, d, auth, "rest:mismatch", testCommit(genesis, 1, "M"))
		assertReuse(t, rec, &errLog, "rest:mismatch", "field WorldRef mismatch")
		if got := headOf(t, d); got != genesis.Ref.String() {
			t.Fatalf("a refused mismatch moved the head to %s", got)
		}
	})
}
