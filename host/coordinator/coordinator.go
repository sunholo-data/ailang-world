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
	"sync"

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
	AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error)
	GetReceipt(ctx context.Context, id string) (store.Receipt, bool, error)
	Commit(ctx context.Context, c store.Commit) error
	// EffectSpend seeds an effectful dispatch's grants (row 134 §4.4).
	EffectSpend(ctx context.Context, episodeID string) (map[store.EffectKey]int64, error)
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
type Coordinator struct {
	cfg Config
	// inFlight protects one process's in-progress invocation IDs. The store's single-process writer premise bounds this guard.
	inFlight        sync.Map
	uncertainIntent sync.Map
	// episodeLocks serialises effectful dispatch per episode (row 134 §4.4):
	// episodeID → a one-slot channel. One writer process per store is
	// enforced by the store's flock (V42), so an in-process lock serialises
	// every spend-seeding reader against every effect-intent writer.
	episodeLocks sync.Map
}

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

// Surface names the transport an invocation arrived on. It is the namespace
// of the invocation id and the tag in the entry header's writtenBy (row 136).
type Surface string

const (
	SurfaceA2A Surface = "a2a"
	SurfaceMCP Surface = "mcp"
)

// Call is one invocation. Request is the ONE registry + capability snapshot
// pair the caller admitted against; Dispatch never re-reads the registry.
// Surface is required: an empty or unknown surface is an InvalidCallError.
type Call struct {
	Surface   Surface
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

// InvocationID names the journal row for one session's task on one surface.
// The "a2a:" and "mcp:" namespaces never collide with the journal's reserved
// "effect:" namespace or /v1/commit's "rest:".
func InvocationID(surface Surface, episodeID, taskID string) string {
	return string(surface) + ":" + episodeID + ":" + taskID
}

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
	if call.Surface != SurfaceA2A && call.Surface != SurfaceMCP {
		return nil, &InvalidCallError{Field: "surface"}
	}
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

// hashes reports whether payload hashes to ref under ref's own algorithm.
// The store returns a payload under the ref it was asked for without
// re-hashing it (store.GetObject), so every read-back is checked here.
func hashes(ref hashref.HashRef, payload []byte) bool {
	sum, err := hashref.Sum(ref.Algo(), payload)
	return err == nil && sum == ref
}

// committed reads a resolved invocation back from the journal: record,
// output and world are all immutable rows written by the committing
// transaction, so the answer is the one the first call produced. Each row is
// verified against the ref that names it, and the three are verified against
// each other, before anything is reported as reconciled (row 122).
func (c *Coordinator) committed(ctx context.Context, rc store.Receipt) (Result, error) {
	recObj, ok, err := c.cfg.Store.GetObject(ctx, rc.Intent.TransitionRef)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: record: %w", rc.InvocationID, err)
	}
	if !ok {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "absent"}
	}
	if !hashes(rc.Intent.TransitionRef, recObj.Payload) {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "mismatch"}
	}
	// v1 and v2 share every field committed() checks; a v2 record must also
	// name its plan object and effect records by well-formed refs.
	var rec recordV2
	if err := json.Unmarshal(recObj.Payload, &rec); err != nil {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "undecodable"}
	}
	if recObj.SemanticID == RecordV2 && !wellFormedV2(rec) {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "undecodable"}
	}
	if rec.InvocationID != rc.InvocationID {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "mismatch"}
	}
	outRef, err := hashref.Parse(rec.Output)
	if err != nil {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "record", Kind: "undecodable"}
	}
	outObj, ok, err := c.cfg.Store.GetObject(ctx, outRef)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: output: %w", rc.InvocationID, err)
	}
	if !ok {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "output", Kind: "absent"}
	}
	if !hashes(outRef, outObj.Payload) {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "output", Kind: "mismatch"}
	}
	w, ok, err := c.cfg.Store.GetWorld(ctx, rc.Intent.WorldRef)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: reconcile %s: world: %w", rc.InvocationID, err)
	}
	if !ok {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "world", Kind: "absent"}
	}
	// The store verifies no world ref (store/durable.go); coordinator worlds are
	// planInvocation's, so the row must re-derive its ref, and applyRevision
	// makes its state root the recorded output.
	if worldRef(w) != rc.Intent.WorldRef {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "world", Kind: "mismatch"}
	}
	if w.StateRoot != outRef {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "world", Kind: "mismatch"}
	}
	var obj map[string]any
	if err := json.Unmarshal(outObj.Payload, &obj); err != nil || obj == nil {
		return Result{}, &IntegrityError{InvocationID: rc.InvocationID, Object: "output", Kind: "undecodable"}
	}
	return Result{Output: obj, OutputBytes: outObj.Payload, InvocationID: rc.InvocationID,
		WorldRef: w.Ref, EntryIndex: w.Revision, RecordRef: recObj.Hash, Reconciled: true}, nil
}

// Dispatch runs one call through propose → verify → execute → commit. ctx is
// the caller's single bounded context; it is never replaced. It bounds every
// step, including receipt lookup, intent append and commit. A caller deadline
// bounds the durable acquisition and the pre-commit transaction body.
func (c *Coordinator) Dispatch(ctx context.Context, call Call) (Result, error) {
	input, err := c.validateCall(call) // R1
	if err != nil {
		return Result{}, err
	}
	id := InvocationID(call.Surface, call.EpisodeID, call.TaskID)
	if _, loaded := c.inFlight.LoadOrStore(id, struct{}{}); loaded {
		return Result{}, &InFlightError{InvocationID: id} // R17
	}
	defer c.inFlight.Delete(id)
	// An effectful descriptor's whole Dispatch holds its episode's lock, and
	// its grants are seeded from the episode's durable spend (row 134 §4.4).
	// The lookup is the same snapshot and ID Bind resolves below; an absent ID
	// falls through to Bind's R2 refusal.
	grants := call.Grants
	effectful := false
	if d, ok := call.Request.Registry.Lookup(call.SkillID); ok && len(d.DeclaredEffects) != 0 {
		effectful = true
		unlock, err := c.lockEpisode(ctx, call.EpisodeID)
		if err != nil {
			return Result{}, err
		}
		defer unlock()
	}
	// Reconcile before anything runs: a resent task id is answered from the
	// journal and never re-executed or re-committed.
	rc, seen, err := c.cfg.Store.GetReceipt(ctx, id)
	if err != nil {
		return Result{}, fmt.Errorf("coordinator: receipt: %w", err)
	}
	if seen {
		if rc.State == store.ReceiptResolved {
			c.uncertainIntent.Delete(id)
			return c.committed(ctx, rc)
		}
		if _, uncertain := c.uncertainIntent.Load(id); uncertain {
			return Result{}, &UnconfirmedError{InvocationID: id, Err: errors.New("intent outcome indeterminate")}
		}
		return Result{}, &NotCommittedError{InvocationID: id} // R15
	}
	if _, uncertain := c.uncertainIntent.Load(id); uncertain {
		return Result{}, &UnconfirmedError{InvocationID: id, Err: errors.New("intent receipt absent after uncertain append")}
	}
	if effectful {
		if grants, err = c.seedGrants(ctx, call.EpisodeID, call.Grants); err != nil {
			return Result{}, err
		}
	}
	// Propose: authorization + confinement (R2/R3/R4).
	bound, err := transitionreg.Bind(call.Request.Registry, call.SkillID, call.Request.Caps,
		c.cfg.Binder(call.EpisodeID, grants))
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
	if len(d.DeclaredEffects) != 0 {
		return c.dispatchEffectful(ctx, call.Surface, id, call.EpisodeID, input, bound, d)
	}
	src, world, err := c.loadSourceAndWorld(ctx, d)
	if err != nil {
		return Result{}, err
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
	pl := planInvocation(world, call.Surface, id, call.EpisodeID, d, input, outBytes, c.cfg.Now())
	return c.appendAndCommit(ctx, id, pl, outBytes, outObj)
}

// loadSourceAndWorld reads the descriptor's verified source (R6) and the
// selected world (R7).
func (c *Coordinator) loadSourceAndWorld(ctx context.Context, d transitionreg.Descriptor) (store.Object, store.World, error) {
	src, ok, err := c.cfg.Store.GetObject(ctx, d.TransitionFn) // R6
	if err != nil {
		return store.Object{}, store.World{}, fmt.Errorf("coordinator: load source: %w", err)
	}
	if !ok {
		return store.Object{}, store.World{}, &SourceError{Kind: "absent"}
	}
	if sum, err := hashref.Sum(d.TransitionFn.Algo(), src.Payload); err != nil || sum != d.TransitionFn {
		return store.Object{}, store.World{}, &SourceError{Kind: "corrupt"}
	}
	head, ok, err := c.cfg.Store.SelectedHead(ctx) // R7
	if err != nil {
		return store.Object{}, store.World{}, fmt.Errorf("coordinator: read head: %w", err)
	}
	if !ok {
		return store.Object{}, store.World{}, &WorldAbsentError{}
	}
	world, ok, err := c.cfg.Store.GetWorld(ctx, head)
	if err != nil {
		return store.Object{}, store.World{}, fmt.Errorf("coordinator: read world: %w", err)
	}
	if !ok { // a head naming no world row is store damage, not "no world"
		return store.Object{}, store.World{}, fmt.Errorf("coordinator: selected head %s has no world row", head)
	}
	return src, world, nil
}

// appendAndCommit is the durable tail (R13–R16) shared by every invocation:
// the intent, then the compare-and-append commit, both under ctx.
func (c *Coordinator) appendAndCommit(ctx context.Context, id string, pl plan, outBytes []byte, outObj map[string]any) (Result, error) {
	if _, _, err := c.cfg.Store.AppendIntent(ctx, id, pl.Intent); err != nil { // R13
		if store.IsUncertain(err) {
			// The intent may have landed. A same-ID resend must inspect its receipt.
			c.uncertainIntent.Store(id, struct{}{})
			return Result{}, &UnconfirmedError{InvocationID: id, Err: err}
		}
		return Result{}, err
	}
	if err := c.cfg.Store.Commit(ctx, pl.Commit); err != nil {
		if store.IsConflict(err) { // R14: compared before any write, rolled back
			return Result{}, err
		}
		if (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) && !store.IsUncertain(err) {
			return Result{}, err // cancellation before the COMMIT cutoff is definite
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
