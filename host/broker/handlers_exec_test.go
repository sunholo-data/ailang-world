package broker

// Row 140 M1 (w-workspace-exec-toolchain-effect §4.1, §4.2 gate 2): until an
// exec profile can be configured (M2/M3), Workspace.Exec is bound to the typed
// unconfigured refusal, in examples_search's no-corpus shape (V31).

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
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

// MUT-SECOND-LIFECYCLE: handlers_exec.go owns no process lifecycle — no
// exec.Cmd, no Setpgid, no kill, no procbound — so runBounded stays the one
// (§4.5 [REVISED r1]). Like TestProjection_NoDeadlineTampering, a source scan;
// identifiers are read from the syntax tree, so a comment cannot trip it.
func TestExecHandlerOwnsNoProcessLifecycle(t *testing.T) {
	forbidden := func(file string) []string {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		var hits []string
		for _, imp := range f.Imports {
			switch strings.Trim(imp.Path.Value, `"`) {
			case "os/exec", "syscall", "github.com/sunholo-data/ailang-world/host/procbound":
				hits = append(hits, "import "+imp.Path.Value)
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if pkg, ok := x.X.(*ast.Ident); ok {
					switch pkg.Name + "." + x.Sel.Name {
					case "exec.Cmd", "exec.Command", "exec.CommandContext", "syscall.Kill", "syscall.SysProcAttr":
						hits = append(hits, fmt.Sprintf("%s.%s at %s", pkg.Name, x.Sel.Name, fset.Position(x.Pos())))
					}
				}
			case *ast.Ident:
				switch x.Name {
				case "Setpgid", "Kill", "procbound", "killGroup":
					hits = append(hits, fmt.Sprintf("%s at %s", x.Name, fset.Position(x.Pos())))
				}
			}
			return true
		})
		return hits
	}
	if hits := forbidden("handlers_exec.go"); len(hits) != 0 {
		t.Fatalf("handlers_exec.go owns process-lifecycle code (MUT-SECOND-LIFECYCLE): %v", hits)
	}
	// Instrument health: the scanner sees runBounded's own lifecycle.
	if hits := forbidden("handlers.go"); len(hits) < 4 {
		t.Fatalf("the scanner found only %v in handlers.go; it is blind", hits)
	}
}
