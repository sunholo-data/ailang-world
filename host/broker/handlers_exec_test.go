package broker

// Row 140 M1 (w-workspace-exec-toolchain-effect §4.1, §4.2 gate 2): until an
// exec profile can be configured (M2/M3), Workspace.Exec is bound to the typed
// unconfigured refusal, in examples_search's no-corpus shape (V31).

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExecUnconfiguredHandlerRefuses(t *testing.T) {
	var h Handler = ExecUnconfiguredHandler{}
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"command":"test","args":["-run=X"]}`))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("output %q: %v", out, err)
	}
	want := map[string]any{"ok": false, "refused": NoExecProfileRefusal}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("output = %s, want %v", out, want)
	}
	if NoExecProfileRefusal != "no exec profile configured: start ailang-worldd serve with --exec-profile FILE" {
		t.Fatalf("refusal text moved: %q", NoExecProfileRefusal)
	}
	if EffectWorkspaceExec != "Workspace.Exec" {
		t.Fatalf("EffectWorkspaceExec = %q", EffectWorkspaceExec)
	}

	// Anything else is a handler failure, never a refusal that reads as served.
	for _, req := range []EffectRequest{
		{Effect: EffectWorkspaceExec, Scope: "/", Cost: 1},
		{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
	} {
		if out, err := h.Execute(boundedTestContext(t), req, []byte(`{}`)); err == nil {
			t.Fatalf("Execute(%+v) = %s, want an error", req, out)
		}
	}
}
