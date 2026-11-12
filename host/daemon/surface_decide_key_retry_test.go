package daemon

// Row 201 A4 sprint 2 M5b (design §2.7, §6 RD12, RD13 (a), R4e; sprint 1's carried H-LOCK): the key
// arm's busy writer, the retry after a failed world write, the reserved scopes by key, and humanMu
// under concurrent POSTs. The rig's key helpers are surface_decide_key_test.go's.

import (
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// RD12: a held credentials writer is 503 Busy with Retry-After, nothing decided or recorded; once
// it is released the same body decides.
func TestDecidePOSTBusyWriter(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	q, p := r.packet("busy", broker.WorkspaceScope, "agent")
	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
	held, err := store.OpenCredentials(r.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	rec := r.postKey(q, body)
	requireRefusal(t, "busy", rec, 503, "Busy")
	if n, _ := r.decisions(); n != 0 || rec.Header().Get("Retry-After") != "1" {
		t.Fatalf("busy: %d decisions, Retry-After %q", n, rec.Header().Get("Retry-After"))
	}
	_ = held.Close()
	if _, ok := r.keyRow(q); ok {
		t.Fatal("busy: evidence was recorded")
	}
	if rec := r.postKey(q, body); rec.Code != 200 {
		t.Fatalf("released: %d %s", rec.Code, rec.Body.String())
	}
}

// RD13 (a): the world write fails after RecordKeyDecision (a trigger refuses the decision object)
// → 500, one unlanded row, 0 decisions; the identical body retried → 200, one decision, the row
// landed and kept as first recorded (K0's exact match); a replay after that is D6's 409.
func TestDecidePOSTRetryAfterAFailedWorldWrite(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	q, p := r.packet("retry", broker.WorkspaceScope, "agent")
	tq, tp := r.packet("retry-typed", broker.WorkspaceScope, "typed")
	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "deny", p, q, decNow+50)
	db, err := sql.Open("sqlite", r.d.cfg.DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if _, err := db.Exec(`CREATE TRIGGER fail_decision BEFORE INSERT ON objects WHEN NEW.semantic_id = '` + broker.ApprovalDecisionV1 + `' BEGIN SELECT RAISE(ABORT, 'injected world-write failure'); END`); err != nil {
		t.Fatal(err)
	}
	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
	requireRefusal(t, "failed world write, typed", r.post(tq.String(), decideBody("approve", tp, tq, decNow+50), r.proof(tq, "decide-approve")), 500, "Internal")
	first, ok := r.keyRow(q)
	if n, _ := r.decisions(); n != 0 || !ok || first.DecisionRef != "" {
		t.Fatalf("after the failure: %d decisions, row %+v %v", n, first, ok)
	}
	if _, err := db.Exec(`DROP TRIGGER fail_decision`); err != nil {
		t.Fatal(err)
	}
	r.d.approveNow = func() int64 { return decNow + 5 }
	rec := r.postKey(q, body)
	var made surfacepage.DecisionMade
	_ = json.Unmarshal(rec.Body.Bytes(), &made)
	row, _ := r.keyRow(q)
	if n, _ := r.decisions(); rec.Code != 200 || n != 1 || row.DecisionRef != made.Decision || row.DecidedAt != first.DecidedAt {
		t.Fatalf("retry: %d %s; %d decisions; row %+v (first %+v)", rec.Code, rec.Body.String(), n, row, first)
	}
	requireRefusal(t, "replay", r.postKey(q, body), 409, "Conflict")
}

// R4e: publish- and finish-shaped requests with a valid key are refused by L1 (403 with the
// owner), nothing decided or recorded. With L1 compiled out the key decision would land; the
// publish consumer's L2 refuses it (R4b), the must-die pair the mutant ledger records (LAW-D5).
func TestDecideRefusesReservedScopesByKey(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	for scope, owner := range map[string]string{"registry:pkg#effect=Registry.Publish": "world-publish approve", "github-pr:o/r#x": "ailang-worldd episode finish"} {
		q, p := r.packet("reserved-"+owner, scope, "agent")
		got := requireRefusal(t, scope, r.postKey(q, decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)), 403, "Forbidden")
		if !strings.Contains(got.Reason, owner) {
			t.Errorf("%s: %q", scope, got.Reason)
		}
		if _, ok := r.keyRow(q); ok {
			t.Fatalf("%s: evidence was recorded", scope)
		}
	}
	if n, _ := r.decisions(); n != 0 {
		t.Fatalf("%d decisions", n)
	}
}

// H-LOCK (carried from sprint 1): eight concurrent POSTs of one valid key decision are serialized
// by humanMu: exactly one decides and every other sees it decided (D6's 409), never a busy writer
// or a refused second write.
func TestDecideConcurrentPOSTsDecideOnce(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	q, p := r.packet("concurrent", broker.WorkspaceScope, "agent")
	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
	codes := make([]int, 8)
	var wg sync.WaitGroup
	for i := range codes {
		wg.Add(1)
		go func(i int) { defer wg.Done(); codes[i] = r.postKey(q, body).Code }(i)
	}
	wg.Wait()
	ok := 0
	for _, c := range codes {
		if c == 200 {
			ok++
		} else if c != 409 {
			t.Errorf("a concurrent POST answered %d (codes %v)", c, codes)
		}
	}
	if n, _ := r.decisions(); ok != 1 || n != 1 {
		t.Fatalf("%d decided, %d decisions (codes %v)", ok, n, codes)
	}
}
