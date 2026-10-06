package coordinator

// Effectful dispatch (queue row 134, M2a/M2b/M3b; design
// design_docs/planned/w-software-engineering-domain.md §4.1, §4.3, §4.4).
//
// A descriptor with declared effects runs plan → execute → finish → commit:
// the pure transition answers {"phase":"plan","args":…} with a
// world/effect-plan/v1 object (parsePlan, effectplan.go); the coordinator
// requests the planned effect through the descriptor-bound broker invoker
// (TR.C: dispatch stays Bound.Request → BoundInvoker.Request); an optional
// {"phase":"finish",…,"results":[…]} run shapes the result; and the call
// commits a world/invocation-record/v2 naming the plan and the effect
// records. Everything after the effect runs detached from the caller's
// cancellation, and any failure there is the typed *EffectsUnrecordedError —
// never a retry, never an untyped error.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Per-phase deadline budget (§4.3). The handler gets min(HandlerCap,
// remaining − HandlerHeadroom); the headroom is the finish phase plus the
// post-effect writes, so a call whose handler uses its whole cap still
// commits inside the caller's deadline.
const (
	PlanPhaseBudget   = 2 * time.Second
	FinishPhaseBudget = 2 * time.Second
	HandlerCap        = 10 * time.Second
	HandlerHeadroom   = FinishPhaseBudget + PostEffectBudget
	PostEffectBudget  = 4 * time.Second
)

// EffectResult statuses (§4.1). "skipped" is reserved for multi-effect plans
// (R-SE-10); a v1 plan has at most one effect.
const (
	StatusOK     = "ok"
	StatusFailed = "failed"
	StatusDenied = "denied"
)

// phaseInput is the calling convention v2 argument (§4.1).
type phaseInput struct {
	Phase   string          `json:"phase"`
	Args    json.RawMessage `json:"args"`
	Results []EffectResult  `json:"results,omitempty"`
}

// EffectResult is one executed (or refused) planned effect, as the finish
// phase sees it.
type EffectResult struct {
	ID     string          `json:"id"`
	Effect string          `json:"effect"`
	Status string          `json:"status"`
	Record *string         `json:"record"`
	Output json.RawMessage `json:"output"`
}

type worldEffect struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Record *string `json:"record"`
}

// worldBlock is the reserved `world` provenance member of every effectful
// output.
type worldBlock struct {
	Effects []worldEffect `json:"effects"`
	Plan    string        `json:"plan"`
}

// lockEpisode acquires the episode's one-slot lock, or gives up when ctx ends.
func (c *Coordinator) lockEpisode(ctx context.Context, episodeID string) (func(), error) {
	v, _ := c.episodeLocks.LoadOrStore(episodeID, make(chan struct{}, 1))
	slot := v.(chan struct{})
	select {
	case slot <- struct{}{}:
		return func() { <-slot }, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("coordinator: episode %s lock: %w", episodeID, ctx.Err())
	}
}

// seedGrants returns grants with each budget reduced by the episode's
// durable spend for that (effect, scope), floored at 0 (§4.4). Spend counts
// executed effects only: a denial appends no intent (V55).
func (c *Coordinator) seedGrants(ctx context.Context, episodeID string, grants []broker.Capability) ([]broker.Capability, error) {
	spend, err := c.cfg.Store.EffectSpend(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("coordinator: effect spend: %w", err)
	}
	seeded := append([]broker.Capability(nil), grants...)
	for i := range seeded {
		if spent := spend[store.EffectKey{Effect: seeded[i].Effect, Scope: seeded[i].Scope}]; spent > 0 {
			seeded[i].Budget = max(seeded[i].Budget-spent, 0)
		}
	}
	return seeded, nil
}

// handlerBudget is min(HandlerCap, remaining − HandlerHeadroom).
func handlerBudget(ctx context.Context) time.Duration {
	dl, ok := ctx.Deadline()
	if !ok {
		return HandlerCap
	}
	return min(HandlerCap, time.Until(dl)-HandlerHeadroom)
}

// detached returns a context that ignores the caller's cancellation and is
// bounded by its own budget: the post-effect context (§4.3).
func detached(ctx context.Context, budget time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), budget)
}

// runPhase runs one convention-v2 phase in the capsule under its own cap.
func (c *Coordinator) runPhase(ctx context.Context, budget time.Duration, interp hashref.HashRef, src []byte, in phaseInput) ([]byte, error) {
	pctx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()
	res, err := c.cfg.Runner.RunContext(pctx, capsule.Entry{
		Interpreter: interp, Source: src, Args: mustJSON(string(mustJSON(in))),
	})
	if err != nil {
		if cerr := pctx.Err(); cerr != nil {
			return nil, fmt.Errorf("coordinator: execute %s phase: %w", in.Phase, cerr)
		}
		return nil, classifyExec(err)
	}
	out, _, err := c.parseOutput(res.Stdout)
	return out, err
}

// requestEffect requests one planned effect through the bound invoker and
// classifies the broker's answer. ran reports that a handler was (or may have
// been) dispatched; a non-nil err with ran == false is a pre-dispatch refusal
// (nothing executed), and with ran == true an effect whose accounting is
// incomplete.
type requester interface {
	Request(context.Context, broker.EffectRequest, []byte) ([]byte, hashref.HashRef, error)
}

func requestEffect(ctx context.Context, inv requester, e PlannedEffect, now int64) (res EffectResult, ref hashref.HashRef, ran bool, err error) {
	out, ref, err := inv.Request(ctx, broker.EffectRequest{
		Effect: e.Requirement.Effect, Scope: e.Requirement.Scope, Cost: e.Requirement.Cost, Now: now,
	}, e.Payload)
	res = EffectResult{ID: e.ID, Effect: e.Requirement.Effect}
	var denial *broker.DenialError
	var failed *broker.EffectFailedError
	var unwritten *broker.OutcomeWriteError
	var indeterminate *broker.IndeterminateEffectError
	switch {
	case err == nil:
		res.Status, res.Output = StatusOK, out
	case errors.As(err, &denial):
		res.Status, ref = StatusDenied, denial.RecordRef
	case errors.As(err, &failed):
		res.Status, ref = StatusFailed, failed.RecordRef
	case errors.As(err, &unwritten):
		return res, unwritten.RecordRef, true, err
	case errors.As(err, &indeterminate):
		return res, hashref.HashRef{}, true, err
	default:
		return res, hashref.HashRef{}, false, err
	}
	s := ref.String()
	res.Record = &s
	return res, ref, true, nil
}

// composeEffectOutput builds the caller's output (§4.1): the zero-effect rule
// (the plan's own result), the one-effect rule (the handler output's fields),
// or the finish output — each plus the reserved `world` block. finishCtx
// bounds the finish phase.
func (c *Coordinator) composeEffectOutput(finishCtx context.Context, interp hashref.HashRef, src, input []byte,
	pl Plan, planRef hashref.HashRef, results []EffectResult) ([]byte, map[string]any, error) {
	wb := worldBlock{Effects: []worldEffect{}, Plan: planRef.String()}
	for _, r := range results {
		wb.Effects = append(wb.Effects, worldEffect{ID: r.ID, Status: r.Status, Record: r.Record})
	}
	var body []byte
	switch {
	case len(results) == 0:
		body = pl.Result
	case pl.Finish:
		fed := make([]EffectResult, len(results))
		for i, r := range results {
			fed[i] = r
			if r.Output != nil {
				canonical, err := transitionreg.CanonicalJSON(r.Output, maxPhaseOutputRaw)
				if err != nil {
					return nil, nil, &OutputError{Reason: "of effect " + r.ID + " is not JSON"}
				}
				fed[i].Output = canonical
			}
		}
		out, err := c.runPhase(finishCtx, FinishPhaseBudget, interp, src, phaseInput{Phase: "finish", Args: input, Results: fed})
		if err != nil {
			return nil, nil, err
		}
		if body, err = parseFinish(out); err != nil {
			return nil, nil, err
		}
	case len(results) == 1 && results[0].Status == StatusOK:
		canonical, err := transitionreg.CanonicalJSON(results[0].Output, maxPhaseOutputRaw)
		if err != nil || kindOf(canonical) != '{' {
			return nil, nil, &OutputError{Reason: "of effect " + results[0].ID + " is not a JSON object"}
		}
		body = canonical
	case len(results) == 1:
		body = []byte(`{}`) // denied or failed: the world block carries the status and record
	default: // unreachable under L1; kept so a widened L1 cannot silently drop results
		return nil, nil, &OutputError{Reason: "has several effects and no finish phase"}
	}
	members, err := objectMembers(body)
	if err != nil {
		return nil, nil, &OutputError{Reason: "is not a JSON object"}
	}
	if _, ok := members[reservedWorldKey]; ok {
		return nil, nil, lawErr(LawL6, "the output carries the reserved key %q", reservedWorldKey)
	}
	members[reservedWorldKey] = mustJSON(wb)
	outBytes, err := transitionreg.CanonicalJSON(mustJSON(members), maxPhaseOutputRaw)
	if err != nil {
		return nil, nil, &OutputError{Reason: "cannot be encoded"}
	}
	if len(outBytes) > c.cfg.MaxOutput {
		return nil, nil, &OutputError{Reason: "exceeds the output cap"}
	}
	var obj map[string]any
	if err := json.Unmarshal(outBytes, &obj); err != nil || obj == nil {
		return nil, nil, &OutputError{Reason: "is not a JSON object"}
	}
	return outBytes, obj, nil
}

// dispatchEffectful runs an effectful descriptor's invocation. The caller
// holds the episode lock and bound the descriptor over seeded grants.
func (c *Coordinator) dispatchEffectful(ctx context.Context, surface Surface, id, episodeID string, input []byte,
	bound *transitionreg.Bound, d transitionreg.Descriptor) (Result, error) {
	// R8: every declared effect must have a handler in this binder, decided
	// before any effect — or the plan — runs.
	if missing := bound.Unhandled(); len(missing) != 0 {
		return Result{}, &EffectsUnsupportedError{ID: d.ID, Effects: missing}
	}
	src, world, err := c.loadSourceAndWorld(ctx, d)
	if err != nil {
		return Result{}, err
	}
	now := c.cfg.Now() // one logical time: the effect requests and the intent
	planOut, err := c.runPhase(ctx, PlanPhaseBudget, d.Interpreter, src.Payload, phaseInput{Phase: "plan", Args: input})
	if err != nil {
		return Result{}, err
	}
	pl, err := parsePlan(planOut, d.DeclaredEffects)
	if err != nil {
		return Result{}, err
	}
	planRef := hashref.SumSHA256(pl.Canonical)
	if len(pl.Effects) == 0 {
		// A refusal plan: no effect runs, so the durable tail is the pure one.
		outBytes, outObj, err := c.composeEffectOutput(ctx, d.Interpreter, src.Payload, input, pl, planRef, nil)
		if err != nil {
			return Result{}, err
		}
		p := planEffectInvocation(world, surface, id, episodeID, d, input, outBytes, pl.Canonical, nil, now)
		return c.appendAndCommit(ctx, id, p, outBytes, outObj)
	}

	hctx, cancelHandler := context.WithTimeout(ctx, handlerBudget(ctx))
	var results []EffectResult
	var records []hashref.HashRef
	for _, e := range pl.Effects {
		res, ref, ran, err := requestEffect(hctx, bound, e, now)
		if !ref.IsZero() {
			records = append(records, ref)
		}
		if err != nil {
			cancelHandler()
			if ran {
				return Result{}, &EffectsUnrecordedError{InvocationID: id, EffectRecords: records, Cause: err}
			}
			return Result{}, err // refused before dispatch: nothing ran
		}
		results = append(results, res)
	}
	cancelHandler()

	// Post-effect: detached from the caller, bounded per step (§4.3).
	unrecorded := func(cause error) (Result, error) {
		return Result{}, &EffectsUnrecordedError{InvocationID: id, EffectRecords: records, Cause: cause}
	}
	fctx, cancelFinish := detached(ctx, FinishPhaseBudget)
	outBytes, outObj, err := c.composeEffectOutput(fctx, d.Interpreter, src.Payload, input, pl, planRef, results)
	cancelFinish()
	if err != nil {
		return unrecorded(err)
	}
	p := planEffectInvocation(world, surface, id, episodeID, d, input, outBytes, pl.Canonical, records, now)
	pctx, cancelPost := detached(ctx, PostEffectBudget)
	defer cancelPost()
	res, err := c.appendAndCommit(pctx, id, p, outBytes, outObj)
	if err != nil {
		return unrecorded(err)
	}
	return res, nil
}

// ReplayBinderFor opens a replay binder over a committed invocation's effect
// records (prod: broker.OpenReplayBinder, whose registry is nil).
type ReplayBinderFor func(grants []broker.Capability, records []hashref.HashRef) transitionreg.Binder

// Replay re-executes a committed effectful invocation from its v2 record
// (§4.1 Replay, AC-REPLAY-EFFECT): the plan phase is re-run and must equal
// the recorded plan object; the plan's requests are answered by a replay
// session over the recorded effect records (no handler runs); the finish
// phase, if any, is re-run on the reconstructed results; and the recomposed
// output must equal the committed output byte for byte. It returns the
// replayed output bytes.
func (c *Coordinator) Replay(ctx context.Context, invocationID string, open ReplayBinderFor) ([]byte, error) {
	diverged := func(stage, format string, args ...any) error {
		return &ReplayDivergenceError{InvocationID: invocationID, Stage: stage, Reason: fmt.Sprintf(format, args...)}
	}
	rc, ok, err := c.cfg.Store.GetReceipt(ctx, invocationID)
	if err != nil {
		return nil, fmt.Errorf("coordinator: replay %s: receipt: %w", invocationID, err)
	}
	if !ok || rc.State != store.ReceiptResolved || rc.Intent == nil {
		return nil, diverged("receipt", "the invocation is not committed")
	}
	committed, err := c.committed(ctx, rc) // verifies record, output and world
	if err != nil {
		return nil, err
	}
	recObj, _, err := c.cfg.Store.GetObject(ctx, committed.RecordRef)
	if err != nil {
		return nil, fmt.Errorf("coordinator: replay %s: record: %w", invocationID, err)
	}
	if recObj.SemanticID != RecordV2 {
		return nil, diverged("record", "record is %s, not an effectful %s record", recObj.SemanticID, RecordV2)
	}
	var rec recordV2
	if err := json.Unmarshal(recObj.Payload, &rec); err != nil || !wellFormedV2(rec) {
		return nil, diverged("record", "record does not decode")
	}
	get := func(name, text string) ([]byte, error) {
		ref, err := hashref.Parse(text)
		if err != nil {
			return nil, diverged("record", "%s ref: %v", name, err)
		}
		o, ok, err := c.cfg.Store.GetObject(ctx, ref)
		if err != nil {
			return nil, fmt.Errorf("coordinator: replay %s: %s: %w", invocationID, name, err)
		}
		if !ok || !hashes(ref, o.Payload) {
			return nil, diverged("record", "%s object is absent or does not match its ref", name)
		}
		return o.Payload, nil
	}
	src, err := get("source", rec.TransitionFn)
	if err != nil {
		return nil, err
	}
	input, err := get("input", rec.Input)
	if err != nil {
		return nil, err
	}
	recordedPlan, err := get("plan", rec.Plan)
	if err != nil {
		return nil, err
	}
	interp, err := hashref.Parse(rec.Interpreter)
	if err != nil {
		return nil, diverged("record", "interpreter ref: %v", err)
	}

	// (1) the plan phase reproduces the recorded plan object.
	planOut, err := c.runPhase(ctx, PlanPhaseBudget, interp, src, phaseInput{Phase: "plan", Args: input})
	if err != nil {
		return nil, err
	}
	declared := planRequirements(recordedPlan)
	pl, err := parsePlan(planOut, declared)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(pl.Canonical, recordedPlan) {
		return nil, diverged("plan", "the re-run plan differs from the recorded plan object")
	}
	if len(rec.Effects) != len(pl.Effects) {
		return nil, diverged("effects", "%d recorded effect records for %d planned effects", len(rec.Effects), len(pl.Effects))
	}

	// (2) the same requests, answered from the recorded effect records.
	refs := make([]hashref.HashRef, len(rec.Effects))
	var grants []broker.Capability
	seeded := map[store.EffectKey]bool{}
	for i, text := range rec.Effects {
		payload, err := get("effect record", text)
		if err != nil {
			return nil, err
		}
		er, err := broker.DecodeRecord(payload)
		if err != nil {
			return nil, diverged("effects", "effect record %d: %v", i, err)
		}
		refs[i] = hashref.MustParse(text)
		if key := (store.EffectKey{Effect: er.Effect, Scope: er.Scope}); !seeded[key] {
			seeded[key] = true
			grants = append(grants, broker.ReplayGrant(er, rc.Intent.LogicalTime))
		}
	}
	manifest := broker.Manifest{Declared: make([]broker.Requirement, len(declared))}
	for i, r := range declared {
		manifest.Declared[i] = broker.Requirement{Effect: r.Effect, Scope: r.Scope, Cost: r.Cost}
	}
	inv, err := open(grants, refs).Bind(manifest)
	if err != nil {
		return nil, fmt.Errorf("coordinator: replay %s: bind: %w", invocationID, err)
	}
	var results []EffectResult
	for i, e := range pl.Effects {
		res, ref, _, err := requestEffect(ctx, inv, e, rc.Intent.LogicalTime)
		if err != nil {
			return nil, diverged("effects", "effect %s: %v", e.ID, err)
		}
		if ref != refs[i] {
			return nil, diverged("effects", "effect %s answered from %s, recorded %s", e.ID, ref, refs[i])
		}
		results = append(results, res)
	}

	// (3) finish on the reconstructed results, (4) recompose and compare.
	outBytes, _, err := c.composeEffectOutput(ctx, interp, src, input, pl, hashref.SumSHA256(pl.Canonical), results)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(outBytes, committed.OutputBytes) {
		return nil, diverged("output", "the recomposed output differs from the committed output")
	}
	return outBytes, nil
}

// planRequirements reads the (effect, scope, cost) triples a recorded plan
// object names. A replay has no registry snapshot; the recorded plan is the
// authority being reproduced, and byte equality with it is checked.
func planRequirements(planObject []byte) []transitionreg.EffectRequirement {
	var wire struct {
		Effects []struct {
			Effect string `json:"effect"`
			Scope  string `json:"scope"`
			Cost   int64  `json:"cost"`
		} `json:"effects"`
	}
	if json.Unmarshal(planObject, &wire) != nil {
		return nil
	}
	var out []transitionreg.EffectRequirement
	for _, e := range wire.Effects {
		r := transitionreg.EffectRequirement{Effect: e.Effect, Scope: e.Scope, Cost: e.Cost}
		dup := false
		for _, have := range out {
			dup = dup || have == r
		}
		if !dup {
			out = append(out, r)
		}
	}
	return out
}

// wellFormedV2 reports that a v2 record names its plan and effect records by
// well-formed refs.
func wellFormedV2(rec recordV2) bool {
	if _, err := hashref.Parse(rec.Plan); err != nil || rec.Effects == nil {
		return false
	}
	for _, e := range rec.Effects {
		if _, err := hashref.Parse(e); err != nil {
			return false
		}
	}
	return true
}
