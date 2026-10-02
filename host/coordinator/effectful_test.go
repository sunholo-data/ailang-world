package coordinator

// Row 134 M2a/M2b/M3b: effectful dispatch through the broker.
//
// The tests inject a handler registry through the production binder,
// broker.OpenBinder(store, episode, grants, reg) (row 134 M4b): the
// coordinator still sees only a transitionreg.Binder, the raw Session never
// crosses into coordinator code, and every effect still goes Bound.Request →
// BoundInvoker.Request.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const (
	fxRead  = "Workspace.Read"
	fxWrite = "Workspace.Write"
	fxScope = "worktree"
	fxEp    = "ep1"
)

var fxReadReq = transitionreg.EffectRequirement{Effect: fxRead, Scope: fxScope, Cost: 1}

func fxGrants(budget int64) []broker.Capability {
	return []broker.Capability{{Effect: fxRead, Scope: fxScope, ExpiresAt: 1000, Budget: budget}}
}

// probe is a JSON-answering handler with a race-safe dispatch counter.
type probe struct {
	n     atomic.Int64
	delay func(ctx context.Context) error
	hook  func()
}

func (p *probe) Execute(ctx context.Context, _ broker.EffectRequest, payload []byte) ([]byte, error) {
	p.n.Add(1)
	if p.hook != nil {
		p.hook()
	}
	if p.delay != nil {
		if err := p.delay(ctx); err != nil {
			return nil, err
		}
	}
	return []byte(`{"content":"probe","ok":true,"payload":` + string(payload) + `}`), nil
}

// phaseRunner answers convention-v2 phases.
type phaseRunner struct {
	plan     string
	finish   func(in map[string]any) (string, error)
	planHook func()
	mu       sync.Mutex
	phases   []string
}

func (f *phaseRunner) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	var arg string
	if err := json.Unmarshal(e.Args, &arg); err != nil {
		return capsule.Result{}, err
	}
	var in map[string]any
	if err := json.Unmarshal([]byte(arg), &in); err != nil {
		return capsule.Result{}, err
	}
	phase, _ := in["phase"].(string)
	f.mu.Lock()
	f.phases = append(f.phases, phase)
	f.mu.Unlock()
	switch phase {
	case "plan":
		if f.planHook != nil {
			f.planHook()
		}
		return capsule.Result{Stdout: []byte(f.plan + "\n")}, nil
	case "finish":
		out, err := f.finish(in)
		return capsule.Result{Stdout: []byte(out)}, err
	}
	return capsule.Result{}, errors.New("unknown phase " + phase)
}

func readPlan(finish bool) string {
	f := "false"
	if finish {
		f = "true"
	}
	return `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Read","scope":"worktree","cost":1,"payload":{"op":"read","path":"a.txt"}}],"finish":` + f + `,"result":null}`
}

// fxRig publishes one effectful descriptor "fx" (access Invoke@world:0 when
// splitAccess, else Workspace.Read@worktree:0) declaring declared.
func fxRig(t *testing.T, splitAccess bool, declared ...transitionreg.EffectRequirement) (*rig, store.World) {
	t.Helper()
	r := newRig(t, false)
	w := r.genesis()
	access := transitionreg.EffectRequirement{Effect: fxRead, Scope: fxScope, Cost: 0}
	if splitAccess {
		access = transitionreg.EffectRequirement{Effect: "Invoke", Scope: accessScope, Cost: 0}
	}
	if len(declared) == 0 {
		declared = []transitionreg.EffectRequirement{fxReadReq}
	}
	r.descs = append(r.descs, transitionreg.Descriptor{
		ID: "fx", TransitionFn: r.source(echoSrc), Interpreter: r.interp, SemanticsEpoch: 1,
		InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`), Access: access,
		DeclaredEffects: declared, Title: "fx", Description: "fx",
	})
	r.publish()
	return r, w
}

func (r *rig) fxCoordinator(st Store, runner Runner, reg broker.Registry) *Coordinator {
	r.t.Helper()
	c, err := New(Config{
		Store: st, Runner: runner, Now: func() int64 { return now }, MaxInput: 1 << 16, MaxOutput: 1 << 16,
		Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(r.st, ep, caps, reg)
		},
	})
	if err != nil {
		r.t.Fatalf("New: %v", err)
	}
	return c
}

func (r *rig) replayBinder() ReplayBinderFor {
	return func(grants []broker.Capability, refs []hashref.HashRef) transitionreg.Binder {
		return broker.OpenReplayBinder(r.st, grants, refs)
	}
}

func (r *rig) effectRecord(ref hashref.HashRef) broker.EffectRecord {
	r.t.Helper()
	obj, ok, err := r.st.GetObject(boundedTestContext(r.t), ref)
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
		r.t.Fatalf("effect record %s: ok=%v err=%v sid=%q", ref, ok, err, obj.SemanticID)
	}
	rec, err := broker.DecodeRecord(obj.Payload)
	if err != nil {
		r.t.Fatalf("decode effect record: %v", err)
	}
	return rec
}

func (r *rig) recordV2(ref hashref.HashRef) recordV2 {
	r.t.Helper()
	obj, ok, err := r.st.GetObject(boundedTestContext(r.t), ref)
	if err != nil || !ok || obj.SemanticID != RecordV2 {
		r.t.Fatalf("invocation record %s: ok=%v err=%v sid=%q, want %s", ref, ok, err, obj.SemanticID, RecordV2)
	}
	var rec recordV2
	if err := json.Unmarshal(obj.Payload, &rec); err != nil {
		r.t.Fatal(err)
	}
	return rec
}

func (r *rig) spend() int64 {
	r.t.Helper()
	s, err := r.st.EffectSpend(boundedTestContext(r.t), fxEp)
	if err != nil {
		r.t.Fatal(err)
	}
	return s[store.EffectKey{Effect: fxRead, Scope: fxScope}]
}

// effectReceiptRecord resolves the episode's ordinal-th effect through the
// store's effect receipt API and returns the outcome's record ref.
func (r *rig) effectReceiptRecord(ordinal int64) hashref.HashRef {
	r.t.Helper()
	rc, ok, err := r.st.GetEffectReceipt(boundedTestContext(r.t), store.EffectInvocationID(fxEp, ordinal))
	if err != nil || !ok || rc.State != store.ReceiptResolved || rc.EffectOutcome == nil {
		r.t.Fatalf("effect receipt %d = %+v ok=%v err=%v, want resolved", ordinal, rc, ok, err)
	}
	return rc.EffectOutcome.RecordRef
}

type worldOut struct {
	World struct {
		Effects []struct {
			ID     string  `json:"id"`
			Status string  `json:"status"`
			Record *string `json:"record"`
		} `json:"effects"`
		Plan string `json:"plan"`
	} `json:"world"`
}

func decodeWorld(t *testing.T, out []byte) worldOut {
	t.Helper()
	var w worldOut
	if err := json.Unmarshal(out, &w); err != nil {
		t.Fatalf("output %s: %v", out, err)
	}
	return w
}

// AC2.1: one effect → the handler's fields plus world.effects[0].record,
// which resolves to an allowed EffectRecordV1; the v2 record names it.
func TestEffectDispatchProbe(t *testing.T) {
	r, w := fxRig(t, false)
	h := &probe{}
	c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
	res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if h.n.Load() != 1 {
		t.Fatalf("handler dispatches = %d, want 1", h.n.Load())
	}
	if !strings.Contains(string(res.OutputBytes), `"content":"probe"`) || !strings.Contains(string(res.OutputBytes), `"payload":{"op":"read","path":"a.txt"}`) {
		t.Fatalf("output lacks the handler bytes: %s", res.OutputBytes)
	}
	wo := decodeWorld(t, res.OutputBytes)
	if len(wo.World.Effects) != 1 || wo.World.Effects[0].Status != StatusOK || wo.World.Effects[0].Record == nil || wo.World.Effects[0].ID != "e1" {
		t.Fatalf("world block = %+v", wo.World)
	}
	recRef := hashref.MustParse(*wo.World.Effects[0].Record)
	if er := r.effectRecord(recRef); !er.Allowed || er.Failed || er.Effect != fxRead {
		t.Fatalf("effect record = %+v, want allowed %s", er, fxRead)
	}
	if got := r.effectReceiptRecord(0); got != recRef {
		t.Fatalf("effect receipt record = %s, want %s", got, recRef)
	}
	rec := r.recordV2(res.RecordRef)
	if len(rec.Effects) != 1 || rec.Effects[0] != recRef.String() || rec.Plan != wo.World.Plan {
		t.Fatalf("v2 record = %+v, want effects [%s] and plan %s", rec, recRef, wo.World.Plan)
	}
	planObj, ok, err := r.st.GetObject(boundedTestContext(t), hashref.MustParse(rec.Plan))
	if err != nil || !ok || planObj.SemanticID != EffectPlanV1 {
		t.Fatalf("plan object ok=%v err=%v sid=%q", ok, err, planObj.SemanticID)
	}
	r.assertHead(res.WorldRef)
	if res.EntryIndex != w.Revision+1 {
		t.Fatalf("entry index %d", res.EntryIndex)
	}
	// committed() decodes v2: a resend is reconciled byte-for-byte, nothing re-runs.
	again, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if err != nil || !again.Reconciled || string(again.OutputBytes) != string(res.OutputBytes) || h.n.Load() != 1 {
		t.Fatalf("resend = %+v %v (dispatches %d), want reconciled same bytes", again, err, h.n.Load())
	}
}

// AC2.2 (MUT-PLAN-UNDECLARED): an undeclared triple in a plan is refused
// before any dispatch.
func TestEffectPlanUndeclaredTriple(t *testing.T) {
	for name, plan := range map[string]string{
		"cost":   strings.Replace(readPlan(false), `"cost":1`, `"cost":0`, 1),
		"effect": strings.Replace(readPlan(false), fxRead, fxWrite, 1),
		"scope":  strings.Replace(readPlan(false), `"scope":"worktree"`, `"scope":"/"`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			r, w := fxRig(t, false)
			h := &probe{}
			c := r.fxCoordinator(r.st, &phaseRunner{plan: plan}, broker.Registry{fxRead: h, fxWrite: h})
			_, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
			var law *PlanLawError
			if !errors.As(err, &law) || law.Law != LawL3 {
				t.Fatalf("err = %T %v, want L3 *PlanLawError", err, err)
			}
			if h.n.Load() != 0 || r.spend() != 0 {
				t.Fatalf("dispatches %d spend %d, want 0/0", h.n.Load(), r.spend())
			}
			r.assertUntouched(w.Ref, "t1")
		})
	}
}

// AC2.3 (MUT-R8-PARTIAL): a declared effect with no registered handler is
// refused before any effect — and before the plan — runs, even when the
// plan would only use a registered one.
func TestEffectUnregisteredDeclared(t *testing.T) {
	r, w := fxRig(t, false, fxReadReq, transitionreg.EffectRequirement{Effect: fxWrite, Scope: fxScope, Cost: 1})
	h := &probe{}
	runner := &phaseRunner{plan: readPlan(false)}
	c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
	_, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	var unsupported *EffectsUnsupportedError
	if !errors.As(err, &unsupported) || len(unsupported.Effects) != 1 || unsupported.Effects[0] != fxWrite {
		t.Fatalf("err = %T %v, want *EffectsUnsupportedError naming %s", err, err, fxWrite)
	}
	if h.n.Load() != 0 || r.spend() != 0 {
		t.Fatalf("dispatches %d spend %d, want 0/0", h.n.Load(), r.spend())
	}
	r.assertUntouched(w.Ref, "t1")
}

// The zero-effect rule (§4.1): a refusal plan commits, and its result reaches
// the caller with an empty world block.
func TestEffectZeroEffectRefusalCommits(t *testing.T) {
	r, _ := fxRig(t, false)
	h := &probe{}
	plan := `{"plan":"world/effect-plan/v1","effects":[],"finish":false,"result":{"ok":false,"refused":"unknown key limit"}}`
	c := r.fxCoordinator(r.st, &phaseRunner{plan: plan}, broker.Registry{fxRead: h})
	res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	want := `{"ok":false,"refused":"unknown key limit","world":{"effects":[],"plan":"` + hashref.SumSHA256(mustCanon(t, plan)).String() + `"}}`
	if string(res.OutputBytes) != want {
		t.Fatalf("output = %s\nwant     %s", res.OutputBytes, want)
	}
	if rec := r.recordV2(res.RecordRef); len(rec.Effects) != 0 || rec.Effects == nil {
		t.Fatalf("v2 record effects = %#v, want []", rec.Effects)
	}
	if h.n.Load() != 0 {
		t.Fatal("refusal plan dispatched a handler")
	}
}

func mustCanon(t *testing.T, raw string) []byte {
	t.Helper()
	b, err := transitionreg.CanonicalJSON([]byte(raw), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// The finish rule: the finish phase sees the EffectResult and its output is
// returned plus the world block.
func TestEffectFinishPhase(t *testing.T) {
	r, _ := fxRig(t, false)
	h := &probe{}
	var seen []any
	runner := &phaseRunner{plan: readPlan(true), finish: func(in map[string]any) (string, error) {
		seen, _ = in["results"].([]any)
		return `{"passed":true}`, nil
	}}
	c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
	res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if !strings.HasPrefix(string(res.OutputBytes), `{"passed":true,"world":{"effects":[{"id":"e1","record":"sha256:`) {
		t.Fatalf("output = %s", res.OutputBytes)
	}
	if len(seen) != 1 {
		t.Fatalf("finish results = %v", seen)
	}
	got, _ := seen[0].(map[string]any)
	if got["status"] != StatusOK || got["effect"] != fxRead || got["output"].(map[string]any)["content"] != "probe" {
		t.Fatalf("finish saw %v", got)
	}
}

type conflictCommitStore struct{ *store.Store }

func (conflictCommitStore) Commit(context.Context, store.Commit) error {
	return &store.ConflictError{}
}

// AC2.4 (MUT-CONFLICT-REEXEC): an R14 after the effect returns the typed
// EffectsUnrecordedError naming the effect record, which the store's effect
// receipt resolves; the effect ran exactly once and is debited once.
func TestEffectConflictUnrecorded(t *testing.T) {
	r, w := fxRig(t, false)
	h := &probe{}
	c := r.fxCoordinator(conflictCommitStore{r.st}, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
	_, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if h.n.Load() != 1 || r.spend() != 1 {
		t.Fatalf("dispatches %d spend %d, want 1/1 (no re-execution on R14)", h.n.Load(), r.spend())
	}
	var unrec *EffectsUnrecordedError
	if !errors.As(err, &unrec) || !store.IsConflict(err) || unrec.InvocationID != InvocationID(fxEp, "t1") {
		t.Fatalf("err = %T %v, want *EffectsUnrecordedError wrapping R14", err, err)
	}
	if len(unrec.EffectRecords) != 1 || r.effectReceiptRecord(0) != unrec.EffectRecords[0] {
		t.Fatalf("effect records %v do not resolve through the effect receipt", unrec.EffectRecords)
	}
	if er := r.effectRecord(unrec.EffectRecords[0]); !er.Allowed || er.Failed {
		t.Fatalf("effect record = %+v", er)
	}
	if h.n.Load() != 1 || r.spend() != 1 {
		t.Fatalf("dispatches %d spend %d, want 1/1", h.n.Load(), r.spend())
	}
	r.assertHead(w.Ref)
}

// AC2.5 (MUT-NO-HEADROOM): a handler that sleeps past its cap is recorded
// failed, and the call still commits inside the caller's deadline.
func TestEffectHandlerCapCommitsFailed(t *testing.T) {
	r, _ := fxRig(t, false)
	h := &probe{delay: func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second):
			return nil
		}
	}}
	c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
	deadline := HandlerHeadroom + time.Second // handler cap = 1 s
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	start := time.Now()
	res, err := c.Dispatch(ctx, r.call("fx", "t1", fxGrants(5)))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if elapsed >= HandlerHeadroom || ctx.Err() != nil {
		t.Fatalf("returned after %v (ctx %v), want inside the deadline with the headroom intact", elapsed, ctx.Err())
	}
	wo := decodeWorld(t, res.OutputBytes)
	if len(wo.World.Effects) != 1 || wo.World.Effects[0].Status != StatusFailed {
		t.Fatalf("world block = %+v, want failed", wo.World)
	}
	rec := r.recordV2(res.RecordRef)
	if er := r.effectRecord(hashref.MustParse(rec.Effects[0])); !er.Allowed || !er.Failed {
		t.Fatalf("effect record = %+v, want failed", er)
	}
}

type failCommitStore struct{ *store.Store }

func (failCommitStore) Commit(context.Context, store.Commit) error {
	return errors.New("injected commit failure")
}

// AC2.5c (MUT-ATTACHED-COMMIT): the caller's ctx is cancelled after the
// effect. A healthy detached tail commits a v2 record; a failing finish or
// commit returns the typed error with a resolvable effect record — never an
// untyped error, never a lost effect record.
func TestEffectCallerCancelledAfterEffect(t *testing.T) {
	cases := []struct {
		name   string
		finish bool
		fail   bool // finish phase fails
		store  func(*store.Store) Store
	}{
		{"healthy_commit", false, false, nil},
		{"healthy_commit_with_finish", true, false, nil},
		{"finish_fails", true, true, nil},
		{"commit_errors", false, false, func(s *store.Store) Store { return failCommitStore{s} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, w := fxRig(t, false)
			ctx, cancel := context.WithCancel(boundedTestContext(t))
			h := &probe{hook: cancel}
			runner := &phaseRunner{plan: readPlan(tc.finish), finish: func(map[string]any) (string, error) {
				if tc.fail {
					return "", &capsule.ExecError{Stderr: []byte("finish panicked")}
				}
				return `{"done":true}`, nil
			}}
			var st Store = r.st
			if tc.store != nil {
				st = tc.store(r.st)
			}
			res, err := r.fxCoordinator(st, runner, broker.Registry{fxRead: h}).Dispatch(ctx, r.call("fx", "t1", fxGrants(5)))
			if ctx.Err() == nil {
				t.Fatal("the caller ctx was not cancelled by the effect")
			}
			recorded := r.effectReceiptRecord(0)
			if tc.fail || tc.store != nil {
				var unrec *EffectsUnrecordedError
				if !errors.As(err, &unrec) || len(unrec.EffectRecords) != 1 || unrec.EffectRecords[0] != recorded {
					t.Fatalf("err = %T %v, want typed *EffectsUnrecordedError naming %s", err, err, recorded)
				}
				r.assertHead(w.Ref)
				return
			}
			if err != nil {
				t.Fatalf("detached commit failed: %T %v", err, err)
			}
			rec := r.recordV2(res.RecordRef)
			if len(rec.Effects) != 1 || rec.Effects[0] != recorded.String() {
				t.Fatalf("v2 record effects %v, want [%s]", rec.Effects, recorded)
			}
			r.assertHead(res.WorldRef)
		})
	}
}

// AC2.6 (AC-REPLAY-EFFECT; MUT-RECORD-V1, MUT-REPLAY-LIVE): a committed v2
// record re-executes through a nil-registry Replay session to the same
// output bytes, and no handler runs.
func TestEffectReplay(t *testing.T) {
	for _, finish := range []bool{false, true} {
		t.Run(map[bool]string{false: "no_finish", true: "finish"}[finish], func(t *testing.T) {
			r, _ := fxRig(t, false)
			h := &probe{}
			runner := &phaseRunner{plan: readPlan(finish), finish: func(in map[string]any) (string, error) {
				res := in["results"].([]any)[0].(map[string]any)
				return `{"status":"` + res["status"].(string) + `"}`, nil
			}}
			c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
			res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
			if err != nil {
				t.Fatalf("Dispatch: %v", err)
			}
			before := h.n.Load()
			out, err := c.Replay(boundedTestContext(t), res.InvocationID, r.replayBinder())
			if err != nil {
				t.Fatalf("Replay: %v", err)
			}
			if string(out) != string(res.OutputBytes) {
				t.Fatalf("replay output %s\nwant          %s", out, res.OutputBytes)
			}
			if h.n.Load() != before {
				t.Fatalf("replay dispatched a handler (%d → %d)", before, h.n.Load())
			}
		})
	}
	t.Run("denied_and_failed_records_replay", func(t *testing.T) {
		r, _ := fxRig(t, false)
		fail := false
		h := &probe{delay: func(context.Context) error {
			if fail {
				return errors.New("handler failed")
			}
			return nil
		}}
		c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
		fail = true
		failed, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
		if err != nil {
			t.Fatalf("failed dispatch: %v", err)
		}
		denied, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t2", fxGrants(1)))
		if err != nil {
			t.Fatalf("denied dispatch: %v", err)
		}
		for _, res := range []Result{failed, denied} {
			out, err := c.Replay(boundedTestContext(t), res.InvocationID, r.replayBinder())
			if err != nil || string(out) != string(res.OutputBytes) {
				t.Fatalf("replay %s = %s %v, want %s", res.InvocationID, out, err, res.OutputBytes)
			}
		}
		if h.n.Load() != 1 {
			t.Fatalf("dispatches = %d, want 1 (the failed one)", h.n.Load())
		}
	})
	t.Run("pure_record_is_not_replayed_as_effectful", func(t *testing.T) {
		r := newRig(t, false)
		r.genesis()
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.publish()
		c := r.coordinator(r.st, fakeRunner{stdout: `{"ok":true}`})
		res, err := c.Dispatch(boundedTestContext(t), r.call("plain", "t1", grants))
		if err != nil {
			t.Fatal(err)
		}
		var div *ReplayDivergenceError
		if _, err := c.Replay(boundedTestContext(t), res.InvocationID, r.replayBinder()); !errors.As(err, &div) || div.Stage != "record" {
			t.Fatalf("replay of a v1 record = %v, want record divergence", err)
		}
	})
}

// AC3.1 (MUT-BUDGET-FRESH): budget 2 persists across fresh Dispatches: ok,
// ok, then a recorded denied:budget.
func TestEffectBudgetPersistsAcrossDispatches(t *testing.T) {
	r, _ := fxRig(t, false)
	h := &probe{}
	c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
	var statuses []string
	var last Result
	for _, task := range []string{"t1", "t2", "t3"} {
		res, err := c.Dispatch(boundedTestContext(t), r.call("fx", task, fxGrants(2)))
		if err != nil {
			t.Fatalf("Dispatch %s: %v", task, err)
		}
		statuses = append(statuses, decodeWorld(t, res.OutputBytes).World.Effects[0].Status)
		last = res
	}
	if strings.Join(statuses, ",") != "ok,ok,denied" {
		t.Fatalf("statuses = %v, want ok,ok,denied", statuses)
	}
	er := r.effectRecord(hashref.MustParse(*decodeWorld(t, last.OutputBytes).World.Effects[0].Record))
	if er.Allowed || er.Denial != broker.LabelDeniedBudget || er.BudgetBefore != 0 {
		t.Fatalf("third record = %+v, want denied:budget from budget 0", er)
	}
	if h.n.Load() != 2 || r.spend() != 2 {
		t.Fatalf("dispatches %d spend %d, want 2/2", h.n.Load(), r.spend())
	}
}

// AC3.4: denials (scope, then liveness) are not spend; the allowed call after
// them makes spend exactly 1.
func TestEffectDenialsAreNotSpend(t *testing.T) {
	r, _ := fxRig(t, true)
	h := &probe{}
	c := r.fxCoordinator(r.st, &phaseRunner{plan: readPlan(false)}, broker.Registry{fxRead: h})
	invoke := broker.Capability{Effect: "Invoke", Scope: accessScope, ExpiresAt: 1000, Budget: 10}
	rows := []struct {
		task  string
		grant broker.Capability
		want  string
		label string
	}{
		{"t1", broker.Capability{Effect: fxRead, Scope: "elsewhere", ExpiresAt: 1000, Budget: 5}, StatusDenied, broker.LabelDeniedScope},
		{"t2", broker.Capability{Effect: fxRead, Scope: fxScope, ExpiresAt: now, Budget: 5}, StatusDenied, broker.LabelDeniedExpired},
		{"t3", broker.Capability{Effect: fxRead, Scope: fxScope, ExpiresAt: 1000, Budget: 5}, StatusOK, ""},
	}
	for _, row := range rows {
		res, err := c.Dispatch(boundedTestContext(t), r.call("fx", row.task, []broker.Capability{invoke, row.grant}))
		if err != nil {
			t.Fatalf("Dispatch %s: %v", row.task, err)
		}
		e := decodeWorld(t, res.OutputBytes).World.Effects[0]
		if e.Status != row.want || r.effectRecord(hashref.MustParse(*e.Record)).Denial != row.label {
			t.Fatalf("%s: status %s label %q, want %s %q", row.task, e.Status, r.effectRecord(hashref.MustParse(*e.Record)).Denial, row.want, row.label)
		}
		// Every denial label replays from its record alone (broker.ReplayGrant).
		if out, err := c.Replay(boundedTestContext(t), res.InvocationID, r.replayBinder()); err != nil || string(out) != string(res.OutputBytes) {
			t.Fatalf("%s replay = %s %v", row.task, out, err)
		}
	}
	if r.spend() != 1 || h.n.Load() != 1 {
		t.Fatalf("spend %d dispatches %d, want 1/1", r.spend(), h.n.Load())
	}
}

// AC3.3 (MUT-BUDGET-NOLOCK): two concurrent same-episode dispatches with
// budget 1 — exactly one is allowed.
func TestEffectConcurrentSameEpisodeBudget(t *testing.T) {
	r, _ := fxRig(t, false)
	h := &probe{}
	// The plan phase is slow, so without the episode lock both calls would
	// seed from spend 0 before either appends its effect intent.
	runner := &phaseRunner{plan: readPlan(false), planHook: func() { time.Sleep(150 * time.Millisecond) }}
	c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
	calls := []Call{r.call("fx", "t1", fxGrants(1)), r.call("fx", "t2", fxGrants(1))}
	var wg sync.WaitGroup
	start := make(chan struct{})
	statuses := make([]string, len(calls))
	errs := make([]error, len(calls))
	for i := range calls {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			res, err := c.Dispatch(boundedTestContext(t), calls[i])
			errs[i] = err
			if err == nil {
				statuses[i] = decodeWorld(t, res.OutputBytes).World.Effects[0].Status
			}
		}(i)
	}
	close(start)
	wg.Wait()
	// The budget law first: exactly one effect executed and was debited.
	if h.n.Load() != 1 || r.spend() != 1 {
		t.Fatalf("dispatches %d spend %d (errs %v), want exactly one allowed", h.n.Load(), r.spend(), errs)
	}
	for i, err := range errs {
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	ok := 0
	for _, s := range statuses {
		if s == StatusOK {
			ok++
		}
	}
	if ok != 1 || h.n.Load() != 1 || r.spend() != 1 {
		t.Fatalf("statuses %v dispatches %d spend %d, want exactly one allowed", statuses, h.n.Load(), r.spend())
	}
}

// The pure path is untouched by the effectful branch: a pure descriptor
// takes no episode lock (a held lock does not block it).
func TestPureDispatchTakesNoEpisodeLock(t *testing.T) {
	r := newRig(t, false)
	r.genesis()
	r.describe("plain", r.source(echoSrc), "Invoke")
	r.publish()
	c := r.coordinator(r.st, fakeRunner{stdout: `{"ok":true}`})
	unlock, err := c.lockEpisode(boundedTestContext(t), "ep1")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.Dispatch(ctx, r.call("plain", "t1", grants)); err != nil {
		t.Fatalf("pure Dispatch under a held episode lock: %v", err)
	}
}
