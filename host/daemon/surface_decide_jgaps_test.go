package daemon

// Row 201 A4 sprint 2 M4a: the sprint-1 judge's test gaps J31-J33 (iteration 389), end to end.

import (
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

// J31/J32: D1's grammar is anchored and lower-case on both routes: 65 hex, 63 hex and upper-case
// hex are each 400 D1 (a 65-hex ref that slipped past D1 would make hashref.MustParse panic), 0
// decisions; the 64-hex ref is the control.
func TestDecideRefGrammarIsAnchored(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("anchored", broker.WorkspaceScope, "agent")
	proof := r.proof(q, "decide-approve")
	_, head := r.decisions()
	for name, ref := range map[string]string{"65 hex": q.String() + "0", "63 hex": q.String()[:70], "upper hex": "sha256:" + strings.ToUpper(q.Digest())} {
		got := requireRefusal(t, "POST "+name, r.post(ref, decideBody("approve", p, q, decNow+50), proof), 400, "BadRequest")
		if !strings.Contains(got.Reason, "is not a request ref") {
			t.Errorf("POST %s: %q", name, got.Reason)
		}
		requireRefusal(t, "GET "+name, r.serve("GET", "/surface/decisions/"+ref, apLoopback, nil), 400, "BadRequest")
		r.requireNothingDecided(name, q, head)
	}
	if rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil); rec.Code != 200 {
		t.Fatalf("control GET: %d %s", rec.Code, rec.Body.String())
	}
}

// J33 (and sprint 1's D4a-undecodable row, unit-only until now): a stored packet that does not
// decode makes the discovery incomplete, so GET answers D4a's 409 with no subject and the POST
// refuses with 0 decisions.
func TestDecisionSubjectRefusesAnUndecodablePacket(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("undecodable", broker.WorkspaceScope, "agent")
	if err := r.d.store.PutObject(boundedTestContext(t), forgeObject(coordinator.DecisionPacketV1, []byte(`{"not":"a packet"}`))); err != nil {
		t.Fatal(err)
	}
	_, head := r.decisions()
	got := requireRefusal(t, "GET", r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil), 409, "Conflict")
	if !strings.Contains(got.Reason, "cannot be counted completely") {
		t.Errorf("GET: %q", got.Reason)
	}
	requireRefusal(t, "POST", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
	r.requireNothingDecided("undecodable", q, head)
}
