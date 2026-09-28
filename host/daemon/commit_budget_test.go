package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
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

func (f *commitSeam) CommitContext(ctx context.Context, c store.Commit) error {
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
		return f.real.CommitContext(ctx, c)
	case "land-then-uncertain": // COMMIT completed, the caller's deadline won the race
		lctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := f.real.CommitContext(lctx, c); err != nil {
			return err
		}
		waitCtx()
		return &store.UncertainError{Op: "commit", Cause: context.DeadlineExceeded}
	default: // "uncertain-not-landed": rolled back, the caller saw only the expiry
		waitCtx()
		return &store.UncertainError{Op: "commit", Cause: context.DeadlineExceeded}
	}
}

// AppendIntentContext passes through to the real store; "block" models the
// held connection for the intent step too.
func (f *commitSeam) AppendIntentContext(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error) {
	if f.mode == "block" {
		select {
		case <-ctx.Done():
		case <-time.After(seamEscape):
		}
	}
	return f.real.AppendIntentContext(ctx, id, intent)
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

// TestCommitBudgetAndUncertainReconcile is AC5 (row 23 M2): /v1/commit is
// bounded by B6, the two 503s are distinct classes, and an uncertain commit
// reconciles by the log row at its index, never by the head.
func TestCommitBudgetAndUncertainReconcile(t *testing.T) {
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
