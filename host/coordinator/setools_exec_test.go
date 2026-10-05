package coordinator

// Row 140 M1 (w-workspace-exec-toolchain-effect §4.1, §6 AC1.1/AC1.3):
// workspace-exec's plans, run through the REAL capsule under the pinned
// interpreter. Each expected plan is written out here by hand, independently
// of the module's own inline tests and of the verify gate's Python, and each
// is then admitted by parsePlan against the manifest's one declared triple (a
// refusal is a lawful zero-effect plan, L7). The cases too large for an
// inline-test literal (a 513-byte args item, and the byte-not-character
// bound on a multi-byte item) are pinned only here.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// seToolsExecDeclared is workspace-exec's declared triple set (§4.1).
var seToolsExecDeclared = []transitionreg.EffectRequirement{
	{Effect: "Workspace.Exec", Scope: "worktree", Cost: 1},
}

func execPlanBytes(payload string) string {
	return `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Exec","scope":"worktree","cost":1,"payload":` +
		payload + `}],"finish":false,"result":null}`
}

func execIDRefusal(id string) string {
	return runRefusalBytes(`command "` + id + `" refused: a command id must match ^[a-z][a-z0-9-]{0,31}$ (the id of a command in the operator's exec profile)`)
}

// TestSeToolsExecDescriptor is AC1.3's descriptor half: the 9th tool declares
// exactly one triple, its access is the same effect at cost 0, and its input
// schema is closed (additionalProperties:false) over exactly {command, args}.
func TestSeToolsExecDescriptor(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	if len(manifest) != 9 {
		t.Fatalf("manifest has %d tools, want 9 (row 140 adds workspace-exec)", len(manifest))
	}
	e, ok := manifest["workspace-exec"]
	if !ok {
		t.Fatal("manifest has no workspace-exec")
	}
	if e.TransitionFnFile != "packages/se-tools/se_tools/exec.ail" {
		t.Fatalf("transitionFnFile = %q", e.TransitionFnFile)
	}
	if !reflect.DeepEqual(e.declared(), seToolsExecDeclared) {
		t.Fatalf("workspace-exec declares %+v, want %+v", e.declared(), seToolsExecDeclared)
	}
	if e.Access != (seToolsRequirement{Effect: "Workspace.Exec", Scope: "worktree", Cost: 0}) {
		t.Fatalf("access = %+v, want (Workspace.Exec, worktree, 0)", e.Access)
	}
	var schema struct {
		Type                 string                     `json:"type"`
		Properties           map[string]json.RawMessage `json:"properties"`
		Required             []string                   `json:"required"`
		AdditionalProperties *bool                      `json:"additionalProperties"`
	}
	if err := json.Unmarshal(e.InputSchema, &schema); err != nil {
		t.Fatal(err)
	}
	var props []string
	for k := range schema.Properties {
		props = append(props, k)
	}
	sort.Strings(props)
	if schema.Type != "object" || schema.AdditionalProperties == nil || *schema.AdditionalProperties ||
		fmt.Sprint(props) != "[args command]" || fmt.Sprint(schema.Required) != "[command]" {
		t.Fatalf("input schema = %s; want a closed object over {command, args}, command required", e.InputSchema)
	}
	for _, phrase := range []string{"exec profile", "no network", "sandbox"} {
		if !strings.Contains(e.Description, phrase) {
			t.Fatalf("description lacks %q: %s", phrase, e.Description)
		}
	}
}

// TestSeToolsExecPlans is AC1.1 through the real capsule: the admitted plans'
// exact bytes, and each id / bounds / L-ARGS refusal as a zero-effect plan.
// MUT-PREFIX-FROM-AGENT's plan half: an id with an uppercase letter or a
// leading "-" (an executable or a flag) never plans an effect.
func TestSeToolsExecPlans(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	e, ok := manifest["workspace-exec"]
	if !ok {
		t.Fatal("manifest has no workspace-exec")
	}
	declared := e.declared()
	rig := newSeToolsRig(t)

	items := func(n int, item string) string {
		return "[" + strings.TrimSuffix(strings.Repeat(`"`+item+`",`, n), ",") + "]"
	}
	at512 := strings.Repeat("a", 512)
	over512 := strings.Repeat("a", 513)
	wide256 := strings.Repeat("é", 256) // 512 UTF-8 bytes, 256 characters
	wide257 := strings.Repeat("é", 257) // 514 UTF-8 bytes, 257 characters

	for _, tc := range []struct {
		name, args, want string
	}{
		{"command only", `{"command":"test"}`, execPlanBytes(`{"command":"test"}`)},
		{"command and args", `{"command":"test","args":["-run=X","./a/..."]}`,
			execPlanBytes(`{"command":"test","args":["-run=X","./a/..."]}`)},
		{"empty args kept", `{"command":"test","args":[]}`, execPlanBytes(`{"command":"test","args":[]}`)},
		{"16 args", `{"command":"test","args":` + items(16, "a") + `}`,
			execPlanBytes(`{"command":"test","args":` + items(16, "a") + `}`)},
		{"512-byte item", `{"command":"test","args":["` + at512 + `"]}`,
			execPlanBytes(`{"command":"test","args":["` + at512 + `"]}`)},
		{"512 bytes of 2-byte characters", `{"command":"test","args":["` + wide256 + `"]}`,
			execPlanBytes(`{"command":"test","args":["` + wide256 + `"]}`)},
		{"32-char id", `{"command":"a-345678901234567890123456789012"}`,
			execPlanBytes(`{"command":"a-345678901234567890123456789012"}`)},
		{"args verbatim (the grammar is the handler's)", `{"command":"test","args":["--settings","x.json","--"]}`,
			execPlanBytes(`{"command":"test","args":["--settings","x.json","--"]}`)},
		{"uppercase id", `{"command":"Test"}`, execIDRefusal("Test")},
		{"flag as id", `{"command":"-x"}`, execIDRefusal("-x")},
		{"executable path as id", `{"command":"/bin/sh"}`, execIDRefusal("/bin/sh")},
		{"33-char id", `{"command":"abcdefghijklmnopqrstuvwxyzabcdefg"}`, execIDRefusal("abcdefghijklmnopqrstuvwxyzabcdefg")},
		{"non-ASCII id", `{"command":"tést"}`, execIDRefusal("tést")},
		{"17 args", `{"command":"test","args":` + items(17, "a") + `}`,
			runRefusalBytes(`argument "args" has 17 items; the limit is 16`)},
		{"513-byte item", `{"command":"test","args":["` + over512 + `"]}`,
			runRefusalBytes(`an args item is 513 bytes; each must be 1 to 512`)},
		{"514 bytes of 2-byte characters", `{"command":"test","args":["` + wide257 + `"]}`,
			runRefusalBytes(`an args item is 514 bytes; each must be 1 to 512`)},
		{"empty item", `{"command":"test","args":[""]}`, runRefusalBytes(`an args item is 0 bytes; each must be 1 to 512`)},
		{"NUL in an item", `{"command":"test","args":["-run=a\u0000b"]}`, runRefusalBytes(`an args item contains a NUL byte`)},
		{"args as a string", `{"command":"test","args":"-run=X"}`, runRefusalBytes(`argument "args" must be an array of strings`)},
		{"non-string item", `{"command":"test","args":[true]}`, runRefusalBytes(`argument "args" must be an array of strings`)},
		{"unknown key", `{"command":"test","env":{"GOFLAGS":"-exec=sh"}}`,
			runRefusalBytes(`unknown argument "env"; workspace-exec admits only: command, args`)},
		{"non-object arguments", `"test"`, runRefusalBytes(`arguments must be a JSON object`)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := string(rig.run(e, planInput(tc.args)))
			if got != tc.want {
				g, w := got, tc.want
				if len(g) > 400 {
					g = g[:400] + "…"
				}
				if len(w) > 400 {
					w = w[:400] + "…"
				}
				t.Fatalf("plan bytes:\n got %s\nwant %s", g, w)
			}
			plan, err := parsePlan([]byte(got), declared)
			if err != nil {
				t.Fatalf("parsePlan refused the plan: %v", err)
			}
			if strings.Contains(tc.want, `"effects":[]`) {
				if len(plan.Effects) != 0 {
					t.Fatalf("refusal planned %d effects", len(plan.Effects))
				}
				return
			}
			if len(plan.Effects) != 1 || plan.Effects[0].Requirement != seToolsExecDeclared[0] {
				t.Fatalf("plan effects = %+v", plan.Effects)
			}
		})
	}
}
