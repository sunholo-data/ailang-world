package daemon

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

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

// TestDaemonInvocationMintsDecisionPacket (AC3): the production binder reaches Human.Approve.
func TestDaemonInvocationMintsDecisionPacket(t *testing.T) {
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	c := humanCoordinator(t, d, d.binder, humanPlan("e1"))
	ctx := boundedTestContext(t)
	res, err := c.Dispatch(ctx, humanCall(t, d, "ep1", "t1"))
	if err != nil {
		t.Fatal(err)
	}
	if n := packetCount(t, d); n != 1 {
		t.Fatalf("packetCount = %d, want 1", n)
	}
	k := packetForEntry(t, d, res.EntryIndex)
	entries, err := d.store.ObjectCommits(ctx, res.RecordRef, -1, 1)
	if err != nil || len(entries) != 1 || entries[0] != res.EntryIndex {
		t.Fatalf("record ObjectCommits = %v (err=%v), want [%d]", entries, err, res.EntryIndex)
	}
	p, err := coordinator.DecodeDecisionPacket(k.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if p.CreatedAt != packetNow || p.DeadlineAt != packetExpiry {
		t.Fatalf("packet times = %d/%d, want %d/%d", p.CreatedAt, p.DeadlineAt, packetNow, packetExpiry)
	}
	page, err := broker.RecentApprovals(ctx, d.store, 40)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, summary := range page.Summaries {
		if summary.RequestRef == p.RequestRef.String() && summary.Status == "pending" {
			found = true
		}
	}
	if !found {
		t.Fatalf("recent approvals lack pending request %s: %+v", p.RequestRef, page.Summaries)
	}
	sh, ok := d.human.(serializedHandler)
	if !ok || sh.mu != &d.humanMu {
		t.Fatalf("Human.Approve handler = %T, want serializedHandler with daemon mutex", d.human)
	}
	c2 := humanCoordinator(t, d, d.binder, humanPlan("e2"))
	res2, err := c2.Dispatch(ctx, humanCall(t, d, "ep1", "t2"))
	if err != nil {
		t.Fatal(err)
	}
	if n := packetCount(t, d); n != 2 {
		t.Fatalf("packetCount = %d, want 2", n)
	}
	k1 := packetForEntry(t, d, res.EntryIndex)
	k2 := packetForEntry(t, d, res2.EntryIndex)
	if k1.Hash != k.Hash || k2.Hash == k1.Hash || res2.EntryIndex == res.EntryIndex {
		t.Fatalf("two packets not first-carried by their own entries: %s at %d, %s at %d", k1.Hash, res.EntryIndex, k2.Hash, res2.EntryIndex)
	}
}

type concurrencyProbe struct {
	mu       sync.Mutex
	cur, max int
	second   chan struct{}
	once     sync.Once
}

func (p *concurrencyProbe) Execute(ctx context.Context, _ broker.EffectRequest, _ []byte) ([]byte, error) {
	p.mu.Lock()
	p.cur++
	p.max = max(p.max, p.cur)
	if p.cur == 2 {
		p.once.Do(func() { close(p.second) })
	}
	p.mu.Unlock()
	select {
	case <-p.second:
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
	}
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
	return []byte(`{"status":"pending","requestRef":"` + hashref.SumSHA256([]byte("ac8")).String() + `"}`), nil
}

// TestHumanApproveIsSerializedAcrossEpisodes (AC8): one daemon mutex covers every episode.
func TestHumanApproveIsSerializedAcrossEpisodes(t *testing.T) {
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	sh := d.human.(serializedHandler)
	probe := &concurrencyProbe{second: make(chan struct{})}
	d.human = serializedHandler{mu: sh.mu, h: probe}
	dispatchPair := func(episodeA, episodeB, taskA, taskB string) [2]error {
		c := humanCoordinator(t, d, d.binder, humanPlan("e1"))
		calls := []coordinator.Call{humanCall(t, d, episodeA, taskA), humanCall(t, d, episodeB, taskB)}
		ctx := boundedTestContext(t)
		var wg sync.WaitGroup
		var errs [2]error
		for i := range calls {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, errs[i] = c.Dispatch(ctx, calls[i])
			}(i)
		}
		done := make(chan struct{})
		go func() { wg.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("concurrent Human.Approve dispatches did not finish")
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
		return errs
	}
	errs := dispatchPair("epA", "epB", "a", "b")
	if probe.max != 1 {
		t.Fatalf("serialized probe max = %d, want 1", probe.max)
	}
	t.Logf("serialized probe max = %d", probe.max)
	requireOneCommitted(t, errs)
	// CONTROL: an unwrapped fresh probe must observe both episodes concurrently.
	probe2 := &concurrencyProbe{second: make(chan struct{})}
	d.human = probe2
	errs = dispatchPair("epC", "epD", "c", "d")
	if probe2.max != 2 {
		t.Fatalf("unwrapped probe max = %d, want 2", probe2.max)
	}
	t.Logf("CONTROL unwrapped probe max = %d", probe2.max)
	requireOneCommitted(t, errs)
}

// requireOneCommitted accepts a world-head conflict on the slower of two
// concurrent dispatches: both plan against the same selected head, so the store
// rejects the second commit as stale (store.ConflictError: callers re-plan).
// That is the store's optimistic concurrency, not Human.Approve's; the probe max
// above is what AC8 measures. Any other error, or no commit at all, fails.
func requireOneCommitted(t *testing.T, errs [2]error) {
	t.Helper()
	committed := 0
	for _, err := range errs {
		switch {
		case err == nil:
			committed++
		case store.IsConflict(err):
			t.Logf("world-head conflict (expected for one of two concurrent commits): %v", err)
		default:
			t.Error(err)
		}
	}
	if committed == 0 {
		t.Error("neither concurrent Human.Approve dispatch committed")
	}
}
