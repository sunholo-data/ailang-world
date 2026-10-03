package coordinator

// Row 135 M1 (w-ailang-run-stdin-argv-caps §4.1, §6 AC1.1/AC1.3): ailang-run's
// stdin / argv / caps plans, run through the REAL capsule under the pinned
// interpreter. Each expected plan is written out here by hand, independently
// of the module's own inline tests and of the verify gate's Python, and each
// is then admitted by parsePlan against the manifest's three declared triples
// (a refusal is a lawful zero-effect plan, L7). The two cases too large for
// an inline-test literal (a 1 025-byte argv item, a 65 537-byte stdin) are
// pinned only here.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// seToolsRunDeclared is ailang-run's declared triple set (§4.2 L3).
var seToolsRunDeclared = []transitionreg.EffectRequirement{
	{Effect: "Ailang.Run", Scope: "worktree", Cost: 1},
	{Effect: "Ailang.RunEnv", Scope: "worktree", Cost: 1},
	{Effect: "Ailang.RunNet", Scope: "worktree", Cost: 1},
}

func runPlanBytes(effect, payload string) string {
	return `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"` + effect +
		`","scope":"worktree","cost":1,"payload":` + payload + `}],"finish":false,"result":null}`
}

func runRefusalBytes(msg string) string {
	return `{"plan":"world/effect-plan/v1","effects":[],"finish":false,"result":{"ok":false,"refused":"` +
		strings.ReplaceAll(msg, `"`, `\"`) + `"}}`
}

func TestSeToolsRunPlansStdinArgvCaps(t *testing.T) {
	manifest := loadSeToolsManifest(t)
	e, ok := manifest["ailang-run"]
	if !ok {
		t.Fatal("manifest has no ailang-run")
	}
	declared := e.declared()
	if fmt.Sprint(declared) != fmt.Sprint(seToolsRunDeclared) {
		t.Fatalf("ailang-run declares %+v, want %+v", declared, seToolsRunDeclared)
	}
	rig := newSeToolsRig(t)

	items := func(n int, item string) string {
		return "[" + strings.TrimSuffix(strings.Repeat(`"`+item+`",`, n), ",") + "]"
	}
	bigItem := strings.Repeat("a", 1025)
	bigStdin := strings.Repeat("x", 65537)
	fullStdin := strings.Repeat("x", 65536)
	unknownCap := `unknown capability "Process"; caps admits only: Declassify, Env, FS, IO, Net`

	for _, tc := range []struct {
		name, args, want string
	}{
		// No new argument: byte-identical to row 134's plan.
		{"no new args", `{"path":"args.ail"}`, runPlanBytes("Ailang.Run", `{"path":"args.ail"}`)},
		{"args_json only", `{"path":"args.ail","args_json":"\"data.txt\""}`,
			runPlanBytes("Ailang.Run", `{"path":"args.ail","args_json":"\"data.txt\""}`)},
		{"stdin", `{"path":"pipeline.ail","stdin":"1\n2\n3\n4\n5\n"}`,
			runPlanBytes("Ailang.Run", `{"path":"pipeline.ail","stdin":"1\n2\n3\n4\n5\n"}`)},
		{"stdin at the bound", `{"path":"p.ail","stdin":"` + fullStdin + `"}`,
			runPlanBytes("Ailang.Run", `{"path":"p.ail","stdin":"`+fullStdin+`"}`)},
		{"argv", `{"path":"sum.ail","argv":["numbers.txt"]}`,
			runPlanBytes("Ailang.Run", `{"path":"sum.ail","argv":["numbers.txt"]}`)},
		{"Env -> RunEnv", `{"path":"sum.ail","caps":["IO","FS","Env"],"argv":["numbers.txt"]}`,
			runPlanBytes("Ailang.RunEnv", `{"path":"sum.ail","argv":["numbers.txt"],"caps":["Env","FS","IO"]}`)},
		{"Net -> RunNet", `{"path":"api.ail","caps":["Net","IO"]}`,
			runPlanBytes("Ailang.RunNet", `{"path":"api.ail","caps":["IO","Net"]}`)},
		{"Declassify stays Run", `{"path":"pi.ail","caps":["IO","Declassify"]}`,
			runPlanBytes("Ailang.Run", `{"path":"pi.ail","caps":["Declassify","IO"]}`)},
		{"argv verbatim", `{"path":"argv.ail","argv":["--caps=IO,Net","-","--","a b"]}`,
			runPlanBytes("Ailang.Run", `{"path":"argv.ail","argv":["--caps=IO,Net","-","--","a b"]}`)},
		{"32 argv items", `{"path":"x.ail","argv":` + items(32, "a") + `}`,
			runPlanBytes("Ailang.Run", `{"path":"x.ail","argv":`+items(32, "a")+`}`)},
		{"Env+Net refused", `{"path":"x.ail","caps":["Env","Net","IO"]}`,
			runRefusalBytes("caps names both Env and Net; one run takes at most one of them")},
		{"unknown cap", `{"path":"x.ail","caps":["IO","Process"]}`, runRefusalBytes(unknownCap)},
		{"duplicate cap", `{"path":"x.ail","caps":["IO","IO"]}`, runRefusalBytes(`capability "IO" is listed twice`)},
		{"caps as a string", `{"path":"x.ail","caps":"IO,Net"}`,
			runRefusalBytes(`argument "caps" must be an array of capability names`)},
		{"33 argv items", `{"path":"x.ail","argv":` + items(33, "a") + `}`,
			runRefusalBytes(`argument "argv" has 33 items; the limit is 32`)},
		{"1025-byte argv item", `{"path":"x.ail","argv":["` + bigItem + `"]}`,
			runRefusalBytes(`an argv item is 1025 bytes; the limit is 1024`)},
		{"NUL in argv", `{"path":"x.ail","argv":["a\u0000b"]}`, runRefusalBytes(`an argv item contains a NUL byte`)},
		{"65537-byte stdin", `{"path":"x.ail","stdin":"` + bigStdin + `"}`,
			runRefusalBytes(`argument "stdin" is 65537 bytes; the limit is 65536`)},
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
			if len(plan.Effects) != 1 || plan.Effects[0].Requirement.Cost != 1 {
				t.Fatalf("plan effects = %+v", plan.Effects)
			}
		})
	}
}
