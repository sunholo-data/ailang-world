package daemon

// Row 201 A4 sprint 1 (design §2.4 [r1], §6 RD9, RD9b-RD9d): the POST counts the packets carrying
// a request over every stored packet, never over ValidatePacket's 100-packet prefix.

import (
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

// RD9: two invocations whose requests collide (one ref) and two packets → 409 D4b, 0 decisions.
func TestDecideRefusesASharedRequest(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("shared-1", broker.WorkspaceScope, "agent")
	q2, _ := r.packet("shared-2", broker.WorkspaceScope, "agent")
	if q2 != q {
		t.Fatalf("identical requests did not collide: %v %v", q, q2)
	}
	_, head := r.decisions()
	got := requireRefusal(t, "shared", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
	if !strings.Contains(got.Reason, "is carried by 2 decision packets; deciding it would decide all of them") {
		t.Errorf("shared: %q", got.Reason)
	}
	r.requireNothingDecided("shared", q, head)
}

// RD9b: 150 packets, two carrying the request, the second sorting after the first 101 (where
// ValidatePacket's duplicate scan stops; it validates the first with code 0) → 409 D4b, 0
// decisions, the head unchanged.
func TestDecideDiscoveryFindsAHiddenSecondPacket(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("hidden-1", broker.WorkspaceScope, "agent")
	fillers := r.fillPackets("hidden", 147)
	late := r.packetAfter("hidden-2", q, broker.WorkspaceScope, fillers, 102)
	if pos := r.packetPosition(late); pos < 101 {
		t.Fatalf("the second packet is at position %d, inside the first 101", pos)
	}
	if res, err := coordinator.ValidatePacket(boundedTestContext(t), r.d.store, p); err != nil || res.Code != 0 {
		t.Fatalf("ValidatePacket of the first = %d, %v; want 0 (the scan cannot see the second)", res.Code, err)
	}
	_, head := r.decisions()
	got := requireRefusal(t, "hidden", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
	if !strings.Contains(got.Reason, "is carried by 2 decision packets") {
		t.Errorf("hidden: %q", got.Reason)
	}
	r.requireNothingDecided("hidden", q, head)
}

// RD9c: 2,001 packets, exactly one carrying the request → 409 D4a (the overflow sentinel), 0
// decisions.
func TestDecideDiscoveryRefusesAnOverflow(t *testing.T) {
	r := newDecideRig(t)
	r.d.commitBudget = time.Minute // 2,001 payload reads in one POST: the budget is not what this test measures (a loaded rig exceeded 3 s)
	q, p := r.packet("overflow", broker.WorkspaceScope, "agent")
	r.fillPackets("overflow", broker.DecisionDiscoveryLimit)
	_, head := r.decisions()
	got := requireRefusal(t, "overflow", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
	if !strings.Contains(got.Reason, "cannot be counted completely (more than 2000 stored packets") {
		t.Errorf("overflow: %q", got.Reason)
	}
	r.requireNothingDecided("overflow", q, head)
}

// RD9d (must live): 150 packets, exactly one carrying the request, sorting after the first 101 →
// the typed POST lands one decision.
func TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix(t *testing.T) {
	r := newDecideRig(t)
	fillers := r.fillPackets("beyond", 149)
	q := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent", packetNow)
	p := r.packetAfter("beyond", q, broker.WorkspaceScope, fillers, 102)
	if pos := r.packetPosition(p); pos < 101 {
		t.Fatalf("the packet is at position %d, inside the first 101", pos)
	}
	if rec := r.post(q.String(), decideBody("deny", p, q, decNow+50), r.proof(q, "decide-deny")); rec.Code != 200 {
		t.Fatalf("beyond the prefix: %d %s", rec.Code, rec.Body.String())
	}
	if n, _ := r.decisions(); n != 1 {
		t.Fatalf("%d decisions, want 1", n)
	}
}
