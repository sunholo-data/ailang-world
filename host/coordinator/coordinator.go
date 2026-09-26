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

