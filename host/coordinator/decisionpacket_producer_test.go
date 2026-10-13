package coordinator

import (
	"context"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Row 172 producer tests (C3): a Human.Approve invocation mints one packet in
// its own commit; every other invocation mints none.

var approveReq = transitionreg.EffectRequirement{Effect: broker.EffectHumanApprove, Scope: "demo", Cost: 1}

func approvePlan(id string) string {
	return `{"plan":"world/effect-plan/v1","effects":[{"id":"` + id + `","effect":"Human.Approve","scope":"demo","cost":1,"payload":{"requester":"coord-test"}}],"finish":false,"result":null}`
}

var invokeGrant = broker.Capability{Effect: "Invoke", Scope: accessScope, ExpiresAt: 1000, Budget: 10}

func approveGrant(expires int64) broker.Capability {
	return broker.Capability{Effect: broker.EffectHumanApprove, Scope: "demo", ExpiresAt: expires, Budget: 5}
}

// countingHandler counts the handler's executions (the reconcile arm asserts zero re-runs).
type countingHandler struct {
	h broker.Handler
	n *atomic.Int64
}

func (c countingHandler) Execute(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, error) {
	c.n.Add(1)
	return c.h.Execute(ctx, req, payload)
}

// packetObjects returns every stored world/decision-packet/v1 object.
func packetObjects(t *testing.T, r *rig) []store.Object {
	t.Helper()
	ctx := boundedTestContext(t)
	refs, err := r.st.ObjectsBySemanticID(ctx, DecisionPacketV1, "", 50)
	if err != nil {
		t.Fatalf("ObjectsBySemanticID: %v", err)
	}
	var objs []store.Object
	for _, ref := range refs {
		obj, ok, err := r.st.GetObject(ctx, ref.Hash)
		if err != nil || !ok {
			t.Fatalf("packet %s: ok=%v err=%v", ref.Hash, ok, err)
		}
		objs = append(objs, obj)
	}
	return objs
}

// packetsIn decodes every stored packet (decode recomputes packetHash).
func packetsIn(t *testing.T, r *rig) []DecisionPacket {
	t.Helper()
	var out []DecisionPacket
	for _, obj := range packetObjects(t, r) {
		p, err := DecodeDecisionPacket(obj.Payload)
		if err != nil {
			t.Fatalf("stored packet %s does not decode: %v", obj.Hash, err)
		}
		out = append(out, p)
	}
	return out
}

// approveRigWith publishes the Human.Approve descriptor with plan answered by the phase runner.
func approveRigWith(t *testing.T, plan string) (*rig, *Coordinator) {
	t.Helper()
	r, _ := fxRig(t, true, approveReq)
	c := r.fxCoordinator(r.st, &phaseRunner{plan: plan}, broker.Registry{broker.EffectHumanApprove: broker.NewHumanHandler(r.st)})
	return r, c
}

func approveRig(t *testing.T) (*rig, *Coordinator) {
	t.Helper()
	return approveRigWith(t, approvePlan("e1"))
}

func TestProducerMintsPacketForApprovalOnly(t *testing.T) {
	t.Run("approval", func(t *testing.T) {
		r, c := approveRig(t)
		if _, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", []broker.Capability{invokeGrant, approveGrant(1760)})); err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if n := len(packetObjects(t, r)); n != 1 {
			t.Fatalf("packets = %d, want 1", n)
		}
	})
	t.Run("denied", func(t *testing.T) {
		r, c := approveRig(t)
		res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", []broker.Capability{invokeGrant, approveGrant(now)}))
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if st := decodeWorld(t, res.OutputBytes).World.Effects[0].Status; st != StatusDenied {
			t.Fatalf("effect status = %q, want denied", st)
		}
		if n := len(packetObjects(t, r)); n != 0 {
			t.Fatalf("packets = %d, want 0 for a denied approval", n)
		}
	})
	t.Run("non-approval", func(t *testing.T) {
		r, _ := fxRig(t, false)
		c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: &probe{}})
		if _, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5))); err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if n := len(packetObjects(t, r)); n != 0 {
			t.Fatalf("packets = %d, want 0 for a non-approval effect", n)
		}
	})
	t.Run("refusal-plan", func(t *testing.T) {
		r, c := approveRigWith(t, `{"plan":"world/effect-plan/v1","effects":[],"finish":false,"result":{"ok":false,"refused":"no"}}`)
		if _, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", []broker.Capability{invokeGrant, approveGrant(1760)})); err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if n := len(packetObjects(t, r)); n != 0 {
			t.Fatalf("packets = %d, want 0 for a refusal plan", n)
		}
	})
	t.Run("two-approvals", func(t *testing.T) {
		r, _ := approveRig(t)
		dispatchTwoApprovals(t, r)
		packets := packetsIn(t, r)
		if len(packets) != 2 {
			t.Fatalf("packets = %d, want 2", len(packets))
		}
		if packets[0].RequestRef != packets[1].RequestRef || packets[0].ProposalHash == packets[1].ProposalHash {
			t.Fatalf("two approvals: request refs %s / %s, proposal hashes %s / %s; want equal request refs and distinct proposals",
				packets[0].RequestRef, packets[1].RequestRef, packets[0].ProposalHash, packets[1].ProposalHash)
		}
	})
}

// dispatchTwoApprovals runs two content-identical approvals with different plans
// (e1, e2) over one store and returns the two results in dispatch order.
func dispatchTwoApprovals(t *testing.T, r *rig) (Result, Result) {
	t.Helper()
	caps := []broker.Capability{invokeGrant, approveGrant(1760)}
	h := broker.Registry{broker.EffectHumanApprove: broker.NewHumanHandler(r.st)}
	a := r.fxCoordinator(r.st, &phaseRunner{plan: approvePlan("e1")}, h)
	b := r.fxCoordinator(r.st, &phaseRunner{plan: approvePlan("e2")}, h)
	res1, err := a.Dispatch(boundedTestContext(t), r.call("fx", "t1", caps))
	if err != nil {
		t.Fatalf("approval 1: %v", err)
	}
	res2, err := b.Dispatch(boundedTestContext(t), r.call("fx", "t2", caps))
	if err != nil {
		t.Fatalf("approval 2: %v", err)
	}
	return res1, res2
}

func TestProducerPacketBindsInvocationFacts(t *testing.T) {
	r, c := approveRig(t)
	caps := []broker.Capability{invokeGrant, approveGrant(now), approveGrant(1760), approveGrant(1900)}
	ctx := boundedTestContext(t)

	// PD22: the deadline is the first allowed grant after the expired one.
	probeReq := broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: "demo", Cost: 1, Now: now}
	var wantDeadline int64
	for _, g := range caps[1:] {
		if broker.Decide(g, probeReq).Allowed {
			wantDeadline = g.ExpiresAt
			break
		}
	}
	if wantDeadline != 1760 {
		t.Fatalf("test ladder: first allowed deadline = %d, want 1760", wantDeadline)
	}

	res, err := c.Dispatch(ctx, r.call("fx", "t1", caps))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	rec := r.recordV2(res.RecordRef)
	er := r.effectRecord(hashref.MustParse(rec.Effects[0]))
	pending, ok, err := r.st.GetObject(ctx, er.ResultRef)
	if err != nil || !ok {
		t.Fatalf("effect result object: ok=%v err=%v", ok, err)
	}
	q, err := broker.ParsePendingApproval(pending.Payload)
	if err != nil {
		t.Fatalf("pending output: %v", err)
	}
	reqObj, ok, err := r.st.GetObject(ctx, q)
	if err != nil || !ok {
		t.Fatalf("approval request object: ok=%v err=%v", ok, err)
	}
	areq, err := broker.DecodeApprovalRequest(reqObj.Payload)
	if err != nil {
		t.Fatalf("decode approval request: %v", err)
	}
	rc, _, err := r.st.GetReceipt(ctx, res.InvocationID)
	if err != nil || rc.Intent == nil {
		t.Fatalf("receipt: intent present=%v err=%v", rc.Intent != nil, err)
	}

	packets := packetsIn(t, r)
	if len(packets) != 1 {
		t.Fatalf("packets = %d, want 1", len(packets))
	}
	p := packets[0]
	if p.CreatedAt != areq.Now || areq.Now != rc.Intent.LogicalTime || p.CreatedAt != now {
		t.Fatalf("createdAt %d, request now %d, intent logicalTime %d, want all %d", p.CreatedAt, areq.Now, rc.Intent.LogicalTime, now)
	}
	if p.DeadlineAt != 1760 {
		t.Fatalf("deadlineAt %d, want 1760", p.DeadlineAt)
	}
	if p.ProposalHash.String() != rec.Plan {
		t.Fatalf("proposalHash %s != record plan %s", p.ProposalHash, rec.Plan)
	}
	if p.RequestRef != q {
		t.Fatalf("requestRef %s != pending requestRef %s", p.RequestRef, q)
	}
	if p.Policy != "Cancel" || p.EscalationsRemaining != 0 {
		t.Fatalf("policy %q escalations %d, want Cancel 0", p.Policy, p.EscalationsRemaining)
	}
}

func TestProducerPacketIsCarriedByTheInvocationCommit(t *testing.T) {
	ctx := boundedTestContext(t)
	t.Run("single", func(t *testing.T) {
		r, c := approveRig(t)
		res, err := c.Dispatch(ctx, r.call("fx", "t1", []broker.Capability{invokeGrant, approveGrant(1760)}))
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		objs := packetObjects(t, r)
		if len(objs) != 1 {
			t.Fatalf("packets = %d, want 1", len(objs))
		}
		got, err := r.st.ObjectCommits(ctx, objs[0].Hash, -1, 10)
		if err != nil || !reflect.DeepEqual(got, []int64{res.EntryIndex}) {
			t.Fatalf("packet carriers = %v, %v; want [%d]", got, err, res.EntryIndex)
		}
		carriers, err := r.st.ObjectCommits(ctx, res.RecordRef, -1, 10)
		if err != nil || len(carriers) == 0 || carriers[0] != res.EntryIndex {
			t.Fatalf("record carriers = %v, %v; want first %d", carriers, err, res.EntryIndex)
		}
	})
	t.Run("two-approvals", func(t *testing.T) {
		r, _ := approveRig(t)
		res1, res2 := dispatchTwoApprovals(t, r)
		byPlan := map[string]int64{}
		for _, res := range []Result{res1, res2} {
			byPlan[r.recordV2(res.RecordRef).Plan] = res.EntryIndex
		}
		objs := packetObjects(t, r)
		if len(objs) != 2 {
			t.Fatalf("packets = %d, want 2", len(objs))
		}
		for _, obj := range objs {
			p, err := DecodeDecisionPacket(obj.Payload)
			if err != nil {
				t.Fatal(err)
			}
			want, ok := byPlan[p.ProposalHash.String()]
			if !ok {
				t.Fatalf("packet proposal %s matches no invocation", p.ProposalHash)
			}
			carriers, err := r.st.ObjectCommits(ctx, obj.Hash, -1, 10)
			if err != nil || len(carriers) == 0 || carriers[0] != want {
				t.Fatalf("packet %s first carrier = %v, %v; want its own entry %d", obj.Hash, carriers, err, want)
			}
		}
	})
}

func TestReplayReconcileToleratesPacketObject(t *testing.T) {
	ctx := boundedTestContext(t)
	r, _ := fxRig(t, true, approveReq)
	n := &atomic.Int64{}
	reg := broker.Registry{broker.EffectHumanApprove: countingHandler{h: broker.NewHumanHandler(r.st), n: n}}
	c := r.fxCoordinator(r.st, &phaseRunner{plan: approvePlan("e1")}, reg)
	caps := []broker.Capability{invokeGrant, approveGrant(1760)}
	res, err := c.Dispatch(ctx, r.call("fx", "t1", caps))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	out, err := c.Replay(ctx, res.InvocationID, r.replayBinder())
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if string(out) != string(res.OutputBytes) {
		t.Fatalf("replay bytes differ from the committed output")
	}
	before := n.Load()
	again, err := c.Dispatch(ctx, r.call("fx", "t1", caps))
	if err != nil {
		t.Fatalf("resend: %v", err)
	}
	if !again.Reconciled || string(again.OutputBytes) != string(res.OutputBytes) {
		t.Fatalf("resend reconciled=%v, bytes equal=%v; want reconciled with equal bytes", again.Reconciled, string(again.OutputBytes) == string(res.OutputBytes))
	}
	if n.Load() != before {
		t.Fatalf("resend re-ran the handler (%d -> %d)", before, n.Load())
	}
	if got := len(packetObjects(t, r)); got != 1 {
		t.Fatalf("packets after resend = %d, want 1", got)
	}
}

// TestProducerMintsNoPacketForMultiEffectPlan (M8 guard): L1 refuses a two-effect
// plan before the producer's single-effect condition is reached, so no packet is
// minted for it. The L1 cap is what makes the `== 1` condition equivalent to `>= 1`.
func TestProducerMintsNoPacketForMultiEffectPlan(t *testing.T) {
	plan := `{"plan":"world/effect-plan/v1","effects":[` +
		`{"id":"e1","effect":"Human.Approve","scope":"demo","cost":1,"payload":{"requester":"coord-test"}},` +
		`{"id":"e2","effect":"Workspace.Read","scope":"worktree","cost":1,"payload":{"op":"read","path":"a.txt"}}],"finish":false,"result":null}`
	r, _ := fxRig(t, true, approveReq, fxReadReq)
	c := r.fxCoordinator(r.st, &phaseRunner{plan: plan}, broker.Registry{broker.EffectHumanApprove: broker.NewHumanHandler(r.st), fxRead: &probe{}})
	_, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", []broker.Capability{invokeGrant, approveGrant(1760), {Effect: fxRead, Scope: fxScope, ExpiresAt: 1000, Budget: 5}}))
	if err == nil || !strings.Contains(err.Error(), "violates L1") {
		t.Fatalf("Dispatch err = %v, want an L1 refusal of the two-effect plan", err)
	}
	if n := len(packetObjects(t, r)); n != 0 {
		t.Fatalf("packets = %d, want 0 for a two-effect plan", n)
	}
}
