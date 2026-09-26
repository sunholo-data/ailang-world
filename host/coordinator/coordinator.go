// Package coordinator is World's server-side propose → verify → commit path
// for registered transitions (w-transition-invocation-coordinator, row 106).
// It composes transitionreg.Bind (propose: authorization + confinement),
// Bound.Check (verify: descriptor pins), the capsule runner (execute) and the
// store journal + compare-and-append Commit (commit) under ONE caller-owned
// bounded context. It lives in the host because every step is an effect
// (S2); the laws it enforces are world/transitions.ail's (plan.go).
package coordinator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Store is the narrow store surface the coordinator needs (prod: *store.Store).
type Store interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	SelectedHead(ctx context.Context) (hashref.HashRef, bool, error)
	GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error)
	AppendIntent(id string, intent store.JournalIntent) (int64, hashref.HashRef, error)
	GetReceipt(id string) (store.Receipt, bool, error)
	Commit(c store.Commit) error
}

// Runner is the capsule seam (prod: *capsule.Runner).
type Runner interface {
	RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error)
}

// BinderFor opens a session-scoped binder (prod: broker.OpenBinder).
type BinderFor func(episodeID string, caps []broker.Capability) transitionreg.Binder

// Config is the construction surface; every field is required.
type Config struct {
	Store     Store
	Runner    Runner
	Binder    BinderFor
	Now       func() int64
	MaxInput  int
	MaxOutput int
}

// Coordinator dispatches authorized invocations.
type Coordinator struct{ cfg Config }

// New validates cfg: a missing seam or a non-positive cap is a construction
// error, never "unlimited".
func New(cfg Config) (*Coordinator, error) {
	switch {
	case cfg.Store == nil, cfg.Runner == nil, cfg.Binder == nil, cfg.Now == nil:
		return nil, errors.New("coordinator: Store, Runner, Binder and Now are required")
	case cfg.MaxInput <= 0 || cfg.MaxOutput <= 0:
		return nil, fmt.Errorf("coordinator: MaxInput/MaxOutput must be positive, got %d/%d", cfg.MaxInput, cfg.MaxOutput)
	}
	return &Coordinator{cfg: cfg}, nil
}

// Call is one invocation. Request is the ONE registry + capability snapshot
// pair the caller admitted against; Dispatch never re-reads the registry.
type Call struct {
	Request   transitionreg.Request
	EpisodeID string
	Grants    []broker.Capability
	SkillID   string
	TaskID    string
	Input     map[string]any
	PinnedFn  *hashref.HashRef
}

// Result is a committed invocation.
type Result struct {
	Output       map[string]any
	OutputBytes  []byte
	InvocationID string
	WorldRef     hashref.HashRef
	EntryIndex   int64
	RecordRef    hashref.HashRef
	// Reconciled is true when the result was read back from the journal for
	// a resent task id rather than produced by this call (nothing re-ran).
	Reconciled bool
}

// InvocationID names the journal row for one session's task. The "a2a:"
// namespace never collides with the journal's reserved "effect:" namespace.
func InvocationID(episodeID, taskID string) string { return "a2a:" + episodeID + ":" + taskID }

func validTaskID(id string) bool {
	if id == "" || len(id) > 128 {
		return false
	}
	for _, r := range id {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '_' || r == '-') {
			return false
		}
	}
	return true
}

func (c *Coordinator) validateCall(call Call) ([]byte, error) {
	if !validTaskID(call.TaskID) {
		return nil, &InvalidCallError{Field: "task id"}
	}
	if call.EpisodeID == "" {
		return nil, &InvalidCallError{Field: "episode"}
	}
	if call.Input == nil {
		return nil, &InvalidCallError{Field: "input"}
	}
	in, err := json.Marshal(call.Input) // map keys sort: canonical
	if err != nil || len(in) > c.cfg.MaxInput {
		return nil, &InvalidCallError{Field: "input"}
	}
	return in, nil
}

// incompatibleMarkers are the pinned interpreter's own diagnostics for a
// source whose entry cannot take one string argument (V10).
var incompatibleMarkers = []string{"ARG_DECODE_MISMATCH", "entrypoint 'main' not found"}

func classifyExec(err error) error {
	var ex *capsule.ExecError
	if errors.As(err, &ex) {
		for _, m := range incompatibleMarkers {
			if bytes.Contains(ex.Stderr, []byte(m)) {
				return &IncompatibleError{Err: err}
			}
		}
	}
	return &ExecutionError{Err: err}
}

func (c *Coordinator) parseOutput(stdout []byte) ([]byte, map[string]any, error) {
	out := []byte(strings.TrimSuffix(string(stdout), "\n"))
	if len(out) > c.cfg.MaxOutput {
		return nil, nil, &OutputError{Reason: "exceeds the output cap"}
	}
	var obj map[string]any
	if err := json.Unmarshal(out, &obj); err != nil || obj == nil {
		return nil, nil, &OutputError{Reason: "is not a JSON object"}
	}
	return out, obj, nil
}

// committed reads a resolved invocation back from the journal: record,
// output and world are all immutable rows written by the committing
// transaction, so the answer is the one the first call produced.
func (c *Coordinator) committed(ctx context.Context, rc store.Receipt) (Result, error) {
	recObj, ok, err := c.cfg.Store.GetObject(ctx, rc.Intent.TransitionRef)
	if err != nil || !ok {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: record: ok=%v: %w", rc.InvocationID, ok, err)
	}
	var rec record
	if err := json.Unmarshal(recObj.Payload, &rec); err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: decode record: %w", rc.InvocationID, err)
	}
	outRef, err := hashref.Parse(rec.Output)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: output ref: %w", rc.InvocationID, err)
	}
	outObj, ok, err := c.cfg.Store.GetObject(ctx, outRef)
	if err != nil || !ok {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: output: ok=%v: %w", rc.InvocationID, ok, err)
	}
	w, ok, err := c.cfg.Store.GetWorld(ctx, rc.Intent.WorldRef)
	if err != nil || !ok {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: world: ok=%v: %w", rc.InvocationID, ok, err)
	}
	var obj map[string]any
	if err := json.Unmarshal(outObj.Payload, &obj); err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: output: %w", rc.InvocationID, err)
	}
	return Result{Output: obj, OutputBytes: outObj.Payload, InvocationID: rc.InvocationID,
		WorldRef: w.Ref, EntryIndex: w.Revision, RecordRef: recObj.Hash, Reconciled: true}, nil
}

// Dispatch runs one call through propose → verify → execute → commit. ctx is
// the caller's single bounded context; it is never replaced. It bounds every
// step that accepts it; AppendIntent, GetReceipt and Commit take no context
// (row 23's policy tranche owns that), so the durable steps are bounded only
// by what the store itself bounds — see the design doc's "Bounded waits".
func (c *Coordinator) Dispatch(ctx context.Context, call Call) (Result, error) {
	input, err := c.validateCall(call) // R1
	if err != nil {
		return Result{}, err
	}
	id := InvocationID(call.EpisodeID, call.TaskID)
	// Reconcile before anything runs: a resent task id is answered from the
	// journal and never re-executed or re-committed.
	rc, seen, err := c.cfg.Store.GetReceipt(id)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: receipt: %w", err)
	}
	if seen {
		if rc.State == store.ReceiptResolved {
			return c.committed(ctx, rc)
		}
		return Result{}, &NotCommittedError{InvocationID: id} // R15
	}
	// Propose: authorization + confinement (R2/R3/R4).
	bound, err := transitionreg.Bind(call.Request.Registry, call.SkillID, call.Request.Caps,
		c.cfg.Binder(call.EpisodeID, call.Grants))
	if err != nil {
		return Result{}, err
	}
	d := bound.Descriptor()
	// Verify: the proposal's pins against the bound descriptor (R5).
	p := transitionreg.Proposal{
		TransitionFn: d.TransitionFn, Interpreter: d.Interpreter, SemanticsEpoch: d.SemanticsEpoch,
		RequiredCaps: d.Access, ExpectedEffects: d.DeclaredEffects,
	}
	if call.PinnedFn != nil {
		p.TransitionFn = *call.PinnedFn
	}
	if err := bound.Check(p); err != nil {
		return Result{}, err
	}
	if len(d.DeclaredEffects) != 0 { // R8
		return Result{}, &EffectsUnsupportedError{ID: d.ID}
	}
	src, ok, err := c.cfg.Store.GetObject(ctx, d.TransitionFn) // R6
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: load source: %w", err)
	}
	if !ok {
		return Result{}, &SourceError{Kind: "absent"}
	}
	if sum, err := hashref.Sum(d.TransitionFn.Algo(), src.Payload); err != nil || sum != d.TransitionFn {
		return Result{}, &SourceError{Kind: "corrupt"}
	}
	head, ok, err := c.cfg.Store.SelectedHead(ctx) // R7
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: read head: %w", err)
	}
	if !ok {
		return Result{}, &WorldAbsentError{}
	}
	world, ok, err := c.cfg.Store.GetWorld(ctx, head)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: read world: %w", err)
	}
	if !ok { // a head naming no world row is store damage, not "no world"
		return Result{}, fmt.Errorf("coordinator: selected head %s has no world row", head)
	}
	// Execute (R9/R10/R11).
	res, err := c.cfg.Runner.RunContext(ctx, capsule.Entry{
		Interpreter: d.Interpreter, Source: src.Payload, Args: mustJSON(string(input)),
	})
	if err != nil {
		if cerr := ctx.Err(); cerr != nil {
			return Result{}, fmt.Errorf("coordinator: execute: %w", cerr)
		}
		return Result{}, classifyExec(err)
	}
	outBytes, outObj, err := c.parseOutput(res.Stdout) // R12
	if err != nil {
		return Result{}, err
	}
	pl := planInvocation(world, id, call.EpisodeID, d, input, outBytes, c.cfg.Now())
	// Commit boundary: a cancelled/expired ctx here leaves no durable mutation.
	if err := ctx.Err(); err != nil { // R11
		return Result{}, fmt.Errorf("coordinator: before commit: %w", err)
	}
	if _, _, err := c.cfg.Store.AppendIntent(id, pl.Intent); err != nil { // R13
		return Result{}, err
	}
	if err := c.cfg.Store.Commit(pl.Commit); err != nil {
		if store.IsConflict(err) { // R14: compared before any write, rolled back
			return Result{}, err
		}
		return Result{}, &UnconfirmedError{InvocationID: id, Err: err} // R16
	}
	// Committed: success is reported even if ctx expired during the durable
	// steps — a landed commit is never answered as a failure.
	return Result{
		Output: outObj, OutputBytes: outBytes, InvocationID: id,
		WorldRef: pl.Commit.NextWorld.Ref, EntryIndex: pl.Commit.Entry.Header.EntryIndex, RecordRef: pl.Record,
	}, nil
}
