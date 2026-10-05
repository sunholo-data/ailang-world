package broker

// Row 140 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.1,
// §4.2): the Workspace.Exec effect, which runs one command of the operator's
// exec profile in the episode worktree under the sandbox runtime. M1 binds
// the name (R8: workspace-exec declares it, so it is always bound) to the
// typed unconfigured refusal below; M2 adds ExecHandler, which matches the
// call against the profile (MatchExecArgs) and runs it through runBounded.

import (
	"context"
	"encoding/json"
	"fmt"
)

// EffectWorkspaceExec is the effect name of the workspace-exec tool, scoped
// to WorkspaceScope like every software-engineering effect.
const EffectWorkspaceExec = "Workspace.Exec"

// NoExecProfileRefusal is the Workspace.Exec answer when the daemon has no
// exec profile (§4.2 gate 2), in examples_search's no-corpus shape (V31).
const NoExecProfileRefusal = "no exec profile configured: start ailang-worldd serve with --exec-profile FILE"

// ExecUnconfiguredHandler serves Workspace.Exec when no exec profile is
// configured: it answers {ok:false, refused: NoExecProfileRefusal} and runs
// nothing. Any other effect or scope is a handler failure.
type ExecUnconfiguredHandler struct{}

func (ExecUnconfiguredHandler) Execute(_ context.Context, req EffectRequest, _ []byte) ([]byte, error) {
	if req.Effect != EffectWorkspaceExec {
		return nil, fmt.Errorf("broker: exec handler does not implement %q", req.Effect)
	}
	if req.Scope != WorkspaceScope {
		return nil, fmt.Errorf("broker: exec handler: scope %q is not %q", req.Scope, WorkspaceScope)
	}
	return json.Marshal(map[string]any{"ok": false, "refused": NoExecProfileRefusal})
}
