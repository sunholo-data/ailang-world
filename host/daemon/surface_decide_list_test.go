package daemon

// Row 201 A4 sprint 4 (design §2.8, DMD1–DMD3 end to end): GET /surface/requests.json lists the
// pending Human.Approve requests a decision card decides, through d.Handler(). Temp stores only.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// requestsJSON is GET /surface/requests.json from the loopback, decoded, and its raw body.
func (r decideRig) requestsJSON() (surfacepage.RequestsJSON, string) {
	r.t.Helper()
	rec := r.serve("GET", "/surface/requests.json", apLoopback, nil)
	var out surfacepage.RequestsJSON
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &out) != nil {
		r.t.Fatalf("requests.json: %d %s", rec.Code, rec.Body.String())
	}
	return out, rec.Body.String()
}

// TestSurfaceRequestsListsDecisions: a decidable request with its packet and deadline, a reserved
// publish-shaped one with the law's D5 reason, a shared one (two packets) with D4b and no packet; a
// packet no commit carries listed with D4b's packet sentence; a decided request, a request no packet
// carries, are absent; the mint lists are unchanged; past the
// deadline the row carries D7.
func TestSurfaceRequestsListsDecisions(t *testing.T) {
	r := newDecideRig(t)
	if out, raw := r.requestsJSON(); !strings.Contains(raw, `"decisions":[]`) || !strings.Contains(raw, `"decisionsTruncated":false`) || len(out.Decisions) != 0 {
		t.Fatalf("an empty world: %s", raw)
	}
	q1, p1 := r.packet("list-ok", broker.WorkspaceScope, "agent-1")
	qr, pr := r.packet("list-reserved", "registry:pkg#effect=Registry.Publish", "agent-r")
	qs, _ := r.packet("list-shared-1", broker.WorkspaceScope, "agent-s")
	r.packet("list-shared-2", broker.WorkspaceScope, "agent-s")
	qd, _ := r.packet("list-decided", broker.WorkspaceScope, "agent-d")
	qz := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-z", packetNow)
	qu := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-u", packetNow) // its packet is stored, never committed (code 2)
	pu := forgeObject(coordinator.DecisionPacketV1, coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("list-uncarried")),
		RequestRef: qu, CreatedAt: packetNow, DeadlineAt: packetNow + 60, Policy: "Cancel"}))
	if err := r.d.store.PutObject(boundedTestContext(t), pu); err != nil {
		t.Fatal(err)
	}
	if _, err := broker.DecideApproval(boundedTestContext(t), r.d.store, qd, "approve", typedDecider, decNow); err != nil {
		t.Fatal(err)
	}
	out, raw := r.requestsJSON()
	if out.DecisionsTruncated || out.Pending == nil || out.Decided == nil || !strings.Contains(raw, `"pending":[]`) {
		t.Fatalf("truncated or mint lists changed: %s", raw)
	}
	by := map[string]surfacepage.PendingDecision{}
	for _, d := range out.Decisions {
		by[d.Request] = d
	}
	if len(out.Decisions) != 4 || len(by) != 4 {
		t.Fatalf("decisions = %+v, want the decidable, the reserved, the shared and the uncarried one", out.Decisions)
	}
	want := surfacepage.PendingDecision{Request: q1.String(), Packet: p1.Hash.String(), Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope,
		Cost: 1, Requester: "agent-1", Deadline: packetNow + 60}
	if by[q1.String()] != want {
		t.Errorf("decidable: %+v, want %+v", by[q1.String()], want)
	}
	if got := by[qr.String()]; got.Packet != pr.Hash.String() || got.Reserved != "publish" || !strings.Contains(got.Refusal, "world-publish approve") {
		t.Errorf("reserved: %+v", got)
	}
	if got := by[qs.String()]; got.Packet != "" || got.Deadline != 0 || !strings.Contains(got.Refusal, "is carried by 2 decision packets") {
		t.Errorf("shared: %+v", got)
	}
	if got := by[qu.String()]; got.Packet != pu.Hash.String() || !strings.Contains(got.Refusal, "the decision packet carrying request "+qu.String()+" is refused: ") {
		t.Errorf("uncarried: %+v", got)
	}
	for _, q := range []string{qd.String(), qz.String()} {
		if _, ok := by[q]; ok {
			t.Errorf("%s is listed", q)
		}
	}
	r.d.approveNow = func() int64 { return packetNow + 60 }
	out, _ = r.requestsJSON()
	for _, d := range out.Decisions {
		if d.Request == q1.String() && !strings.Contains(d.Refusal, "the packet's deadline passed") {
			t.Errorf("past the deadline: %+v", d)
		}
	}
	if len(out.Decisions) != 4 {
		t.Errorf("past the deadline: %+v", out.Decisions)
	}
}

// TestSurfaceRequestsDecisionsTruncated: more than 100 stored packets, or more than 20 requests on
// the approvals chain, mark the list truncated.
func TestSurfaceRequestsDecisionsTruncated(t *testing.T) {
	t.Run("packet-scan", func(t *testing.T) {
		r := newDecideRig(t)
		r.packet("scan-ok", broker.WorkspaceScope, "agent")
		r.fillPackets("scan", 100)
		if out, _ := r.requestsJSON(); !out.DecisionsTruncated {
			t.Errorf("101 packets: not truncated: %+v", out)
		}
	})
	t.Run("packet-scan-at-the-bound", func(t *testing.T) {
		r := newDecideRig(t)
		r.packet("bound-ok", broker.WorkspaceScope, "agent")
		r.fillPackets("bound", 99)
		if out, _ := r.requestsJSON(); out.DecisionsTruncated || len(out.Decisions) != 1 {
			t.Errorf("100 packets: %+v", out)
		}
	})
	t.Run("approvals-chain", func(t *testing.T) {
		r := newDecideRig(t)
		for i := 0; i < 21; i++ {
			approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-"+string(rune('a'+i)), packetNow)
		}
		if out, _ := r.requestsJSON(); !out.DecisionsTruncated {
			t.Errorf("21 requests: not truncated: %+v", out)
		}
	})
}

// TestSurfaceRequestsDecisionsUnreadable: a world store that cannot be read lists no decision (an
// empty list, never null) and keeps the mint lists.
func TestSurfaceRequestsDecisionsUnreadable(t *testing.T) {
	r := newDecideRig(t)
	r.packet("unreadable", broker.WorkspaceScope, "agent")
	if err := r.d.store.Close(); err != nil {
		t.Fatal(err)
	}
	_, raw := r.requestsJSON()
	if !strings.Contains(raw, `"decisions":[]`) || !strings.Contains(raw, `"pending":[]`) || !strings.Contains(raw, `"decisionsTruncated":false`) {
		t.Fatalf("unreadable world store: %s", raw)
	}
}
