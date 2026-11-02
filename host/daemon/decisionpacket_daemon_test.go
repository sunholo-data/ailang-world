package daemon

import (
	"bytes"
	"context"
	"fmt"
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

// probeWaitBound bounds how long the first Human.Approve execution waits for
// the second dispatch (row 235). It returns on the arrival, so it costs nothing
// when healthy; it sits below dispatchPair's bound (+5 s) and the test context
// (30 s), so a missing second dispatch fails by name, never as "max = 1".
const probeWaitBound = 10 * time.Second

// serializedHold is how long the first holder keeps the mutex after the second
// dispatch has reached the arrivalGate: a mutex-less wrapper lets the second
// in, which the probe sees as max = 2 (secondary to the TryLock check).
const serializedHold = 50 * time.Millisecond

// lateSecond is the regression arm's stagger: the 250-400 ms gap that read as
// "unwrapped probe max = 1" on the 200 ms probe (iteration 334 measurement).
// It is measured from the first probe entry (entered), not from goroutine start.
const lateSecond = 300 * time.Millisecond

// arrivalGate counts Human.Approve arrivals OUTSIDE the serialization, so the
// serialized arm can show the second dispatch arrived while the first was
// still inside Execute (it does not observe the second calling Lock).
type arrivalGate struct {
	h      broker.Handler
	n      atomic.Int32
	second chan struct{}
}

func (g *arrivalGate) Execute(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, error) {
	if g.n.Add(1) == 2 {
		close(g.second)
	}
	return g.h.Execute(ctx, req, payload)
}

// concurrencyProbe records how many executions overlap. Its first entry closes
// entered and holds until the second dispatch arrives: inside the probe
// (CONTROL, arrived nil) or at the gate (serialized, then serializedHold).
// held is set only in the serialized arm: each entry then checks it is held
// and counts in unheld any entry where TryLock succeeded. CONTROL's probe has
// held == nil, so the check and the counter are off there.
type concurrencyProbe struct {
	mu                sync.Mutex
	cur, max, entries int
	unheld            int
	missed            string
	inside2           chan struct{}
	entered           chan struct{}
	once              sync.Once
	arrived           <-chan struct{}
	held              *sync.Mutex
	timeout           <-chan time.Time // test-only; nil uses the real probeWaitBound
}

func (p *concurrencyProbe) Execute(ctx context.Context, _ broker.EffectRequest, _ []byte) ([]byte, error) {
	unheld := p.held != nil && p.held.TryLock()
	if unheld {
		p.held.Unlock()
	}
	p.mu.Lock()
	if unheld {
		p.unheld++
	}
	first := p.entries == 0
	if first {
		close(p.entered)
	}
	p.entries++
	p.cur++
	p.max = max(p.max, p.cur)
	if p.cur == 2 {
		p.once.Do(func() { close(p.inside2) })
	}
	p.mu.Unlock()
	if first {
		p.awaitSecond(ctx)
	}
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
	return []byte(`{"status":"pending","requestRef":"` + hashref.SumSHA256([]byte("ac8")).String() + `"}`), nil
}

func (p *concurrencyProbe) awaitSecond(ctx context.Context) {
	var wait <-chan struct{} = p.inside2
	if p.arrived != nil {
		wait = p.arrived
	}
	timeout := p.timeout
	if timeout == nil {
		timeout = time.After(probeWaitBound)
	}
	why := ""
	select {
	case <-wait:
		if p.arrived != nil {
			select {
			case <-p.inside2:
			case <-time.After(serializedHold):
			case <-ctx.Done():
			}
		}
		return
	case <-timeout:
		why = fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)
	case <-ctx.Done():
		why = "test context ended before the second Human.Approve dispatch arrived: " + ctx.Err().Error()
	}
	p.mu.Lock()
	p.missed = why
	p.mu.Unlock()
}

// TestHumanApproveIsSerializedAcrossEpisodes (AC8): one daemon mutex covers every episode.
// Row 235: each arm holds its first execution until the second dispatch arrives,
// and "late-second" replays the CI gap: the second dispatch starts lateSecond
// after the first has entered the probe.
func TestHumanApproveIsSerializedAcrossEpisodes(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stagger time.Duration
	}{{"together", 0}, {"late-second", lateSecond}} {
		t.Run(tc.name, func(t *testing.T) { humanApproveSerializedArms(t, tc.stagger) })
	}
}

func humanApproveSerializedArms(t *testing.T, stagger time.Duration) {
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	sh := d.human.(serializedHandler)
	dispatchPair := func(entered <-chan struct{}, episodeA, episodeB, taskA, taskB string) [2]error {
		c := humanCoordinator(t, d, d.binder, humanPlan("e1"))
		calls := []coordinator.Call{humanCall(t, d, episodeA, taskA), humanCall(t, d, episodeB, taskB)}
		ctx := boundedTestContext(t)
		var wg sync.WaitGroup
		var errs [2]error
		for i := range calls {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if i == 1 && stagger > 0 {
					select {
					case <-entered:
					case <-time.After(probeWaitBound):
						errs[i] = fmt.Errorf("first dispatch did not enter the probe within %s", probeWaitBound)
						return
					case <-ctx.Done():
						errs[i] = ctx.Err()
						return
					}
					time.Sleep(stagger)
				}
				_, errs[i] = c.Dispatch(ctx, calls[i])
			}(i)
		}
		done := make(chan struct{})
		go func() { wg.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(probeWaitBound + 5*time.Second):
			t.Fatal("concurrent Human.Approve dispatches did not finish")
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
		return errs
	}
	gate := &arrivalGate{second: make(chan struct{})}
	probe := &concurrencyProbe{inside2: make(chan struct{}), entered: make(chan struct{}), arrived: gate.second, held: sh.mu}
	gate.h = serializedHandler{mu: sh.mu, h: probe}
	d.human = gate
	errs := dispatchPair(probe.entered, "epA", "epB", "a", "b")
	if probe.missed != "" {
		t.Fatalf("serialized arm measured nothing: %s (dispatch errors %v)", probe.missed, errs)
	}
	if probe.unheld != 0 {
		t.Fatalf("serialized probe ran %d of %d executions without the daemon mutex held", probe.unheld, probe.entries)
	}
	if probe.max != 1 {
		t.Fatalf("serialized probe max = %d, want 1", probe.max)
	}
	t.Logf("serialized probe max = %d, gate arrivals = %d", probe.max, gate.n.Load())
	requireOneCommitted(t, errs)
	// CONTROL: an unwrapped fresh probe must observe both episodes concurrently.
	// held is nil here: no TryLock check, no unheld assertion; missed and max only.
	probe2 := &concurrencyProbe{inside2: make(chan struct{}), entered: make(chan struct{})}
	d.human = probe2
	errs = dispatchPair(probe2.entered, "epC", "epD", "c", "d")
	if probe2.missed != "" {
		t.Fatalf("CONTROL measured nothing: %s (dispatch errors %v)", probe2.missed, errs)
	}
	if probe2.max != 2 {
		t.Fatalf("unwrapped probe max = %d, want 2", probe2.max)
	}
	t.Logf("CONTROL unwrapped probe max = %d", probe2.max)
	requireOneCommitted(t, errs)
}

// TestConcurrencyProbeNamesAMissedArrival drives the failure bound without
// depending on scheduling a second dispatch or waiting for the real 10 s timer.
func TestConcurrencyProbeNamesAMissedArrival(t *testing.T) {
	ctx := boundedTestContext(t)
	timeout := make(chan time.Time, 1)
	probe := &concurrencyProbe{
		inside2: make(chan struct{}),
		entered: make(chan struct{}),
		arrived: make(chan struct{}),
		timeout: timeout,
	}
	bound := time.NewTimer(5 * time.Second)
	defer bound.Stop()
	done := make(chan struct{})
	go func() {
		_, _ = probe.Execute(ctx, broker.EffectRequest{}, nil)
		close(done)
	}()
	select {
	case <-probe.entered:
	case <-bound.C:
		t.Fatalf("%s: Execute did not enter within 5s", t.Name())
	}
	// Buffered so even an ignored timeout seam cannot strand the sender.
	timeout <- time.Time{}
	select {
	case <-done:
	case <-bound.C:
		t.Fatalf("%s: Execute did not return after the injected timeout within 5s", t.Name())
	}
	want := fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)
	if probe.missed != want {
		t.Fatalf("missed = %q, want %q", probe.missed, want)
	}
	if probe.entries != 1 {
		t.Fatalf("entries = %d, want 1", probe.entries)
	}
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
