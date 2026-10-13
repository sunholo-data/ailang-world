package daemon

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"sync/atomic"
)

// Row 172 decision-packet daemon tests (C4.1 reservation; shared helpers for C5-C7).

const packetNow, packetExpiry = int64(1700), int64(1760) // the design's Examples

// humanPlan is a Human.Approve plan with one effect, requested by "packet-test".
func humanPlan(id string) string {
	return `{"plan":"world/effect-plan/v1","effects":[{"id":"` + id + `","effect":"Human.Approve","scope":"worktree","cost":1,"payload":{"requester":"packet-test"}}],"finish":false,"result":null}`
}

// humanPlanRunner answers the plan phase with a fixed plan.
type humanPlanRunner struct {
	plan  string
	plans atomic.Int64
}

func (r *humanPlanRunner) RunContext(_ context.Context, _ capsule.Entry) (capsule.Result, error) {
	r.plans.Add(1)
	return capsule.Result{Stdout: []byte(r.plan + "\n")}, nil
}

// publishHumanTool publishes one skill whose effect is Human.Approve.
func publishHumanTool(t *testing.T, d *Daemon) {
	t.Helper()
	publishTools(t, d, hashref.SumSHA256([]byte("fake-interpreter")), "module host/capsule/main\n\nexport func main(input: string) -> string {\n  input\n}\n", wsTool{ID: "hu.approve", Effect: broker.EffectHumanApprove, Title: "Approve", Description: "decision packet fixture"})
}

// humanGrants is one live Human.Approve grant, expiring at packetExpiry.
func humanGrants() []broker.Capability {
	return []broker.Capability{{Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope, ExpiresAt: packetExpiry, Budget: 5}}
}

// humanCall is the A2A call for hu.approve on episode/task at packetNow.
func humanCall(t *testing.T, d *Daemon, episode, task string) coordinator.Call {
	t.Helper()
	req, err := transitionreg.NewRequest(boundedTestContext(t), transitionreg.NewReader(d.store), wsCaps(humanGrants()), packetNow)
	if err != nil {
		t.Fatal(err)
	}
	return coordinator.Call{Surface: coordinator.SurfaceA2A, Request: req, EpisodeID: episode, Grants: humanGrants(), SkillID: "hu.approve", TaskID: task, Input: map[string]any{}}
}

// testHumanBinder binds an explicit registry carrying Human.Approve (PD25: not d.binder).
func testHumanBinder(d *Daemon) coordinator.BinderFor {
	return func(episodeID string, caps []broker.Capability) transitionreg.Binder {
		reg := map[string]broker.Handler{broker.EffectHumanApprove: broker.NewHumanHandler(d.store)}
		return broker.OpenBinder(d.store, episodeID, caps, reg)
	}
}

// humanCoordinator is a coordinator over d's store with the fixed plan and packetNow.
func humanCoordinator(t *testing.T, d *Daemon, binder coordinator.BinderFor, plan string) *coordinator.Coordinator {
	t.Helper()
	c, err := coordinator.New(coordinator.Config{Store: d.store, Runner: &humanPlanRunner{plan: plan}, Binder: binder,
		Now: func() int64 { return packetNow }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// packetForEntry returns the one decision packet first carried by entry.
func packetForEntry(t *testing.T, d *Daemon, entry int64) store.Object {
	t.Helper()
	ctx := boundedTestContext(t)
	refs, err := d.store.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, "", 100)
	if err != nil {
		t.Fatal(err)
	}
	var found []store.Object
	for _, ref := range refs {
		es, err := d.store.ObjectCommits(ctx, ref.Hash, -1, 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(es) == 1 && es[0] == entry {
			full, ok, err := d.store.GetObject(ctx, ref.Hash)
			if err != nil || !ok {
				t.Fatalf("GetObject %s: ok=%v err=%v", ref.Hash, ok, err)
			}
			found = append(found, full)
		}
	}
	if len(found) != 1 {
		t.Fatalf("packets first carried by entry %d = %d, want 1", entry, len(found))
	}
	return found[0]
}

// packetCount is the number of stored decision packets.
func packetCount(t *testing.T, d *Daemon) int {
	t.Helper()
	refs, err := d.store.ObjectsBySemanticID(boundedTestContext(t), coordinator.DecisionPacketV1, "", 100)
	if err != nil {
		t.Fatal(err)
	}
	return len(refs)
}

// TestRESTCommitRefusesDecisionPacketObjects (AC7): REST may not store the packet label.
func TestRESTCommitRefusesDecisionPacketObjects(t *testing.T) {
	d := newHandlerDaemon(t)
	g := seedGenesisEmbedded(t, d, "ac7")
	c := testCommit(g, 1, "ac7")
	c.Objects[0].SemanticID = coordinator.DecisionPacketV1
	rec := requestRecorderAuth(t, d, authHeader(t, d), http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(c)))
	body := assertErrorClass(t, rec, http.StatusBadRequest, "BadRequest")
	if want := "objects[0].semanticId: world/decision-packet/v1 is minted only by the invocation coordinator"; body.Error.Message != want {
		t.Fatalf("message = %q, want %q", body.Error.Message, want)
	}
	if n := packetCount(t, d); n != 0 {
		t.Fatalf("packetCount = %d, want 0", n)
	}
	head, ok, err := d.store.SelectedHead(boundedTestContext(t))
	if err != nil || !ok || head.String() != g.Ref.String() {
		t.Fatalf("selected head = %v (ok=%v err=%v), want genesis %s", head, ok, err, g.Ref)
	}
	control := testCommit(g, 1, "ac7-control")
	if rec := requestRecorderAuth(t, d, authHeader(t, d), http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(control))); rec.Code != http.StatusOK {
		t.Fatalf("control status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
}
