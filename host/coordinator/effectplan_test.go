package coordinator

import (
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

var planDeclared = []transitionreg.EffectRequirement{
	{Effect: "Workspace.Read", Scope: "worktree", Cost: 1},
	{Effect: "Ailang.Check", Scope: "worktree", Cost: 1},
}

// planJSON builds a world/effect-plan/v1 object; effects and result are raw
// JSON fragments.
func planJSON(effects, finish, result string) string {
	return `{"plan":"world/effect-plan/v1","effects":[` + effects + `],"finish":` + finish + `,"result":` + result + `}`
}

func effectJSON(id, effect, scope, cost, payload string) string {
	return `{"id":"` + id + `","effect":"` + effect + `","scope":"` + scope + `","cost":` + cost + `,"payload":` + payload + `}`
}

var readEffect = effectJSON("e1", "Workspace.Read", "worktree", "1", `{"path":"src/main.ail","op":"read"}`)

// payloadOfCanonicalSize returns a payload object whose canonical encoding is
// exactly n bytes: {"c":"xxx…"} has 8 bytes of framing.
func payloadOfCanonicalSize(n int) string {
	return `{"c":"` + strings.Repeat("x", n-8) + `"}`
}

// TestParsePlanLaws is AC1.1: one row per law (and per arm of a law), each
// refusal typed as *PlanLawError with the violated law named.
func TestParsePlanLaws(t *testing.T) {
	rows := []struct {
		name   string
		output string
		law    PlanLaw // "" = admitted
	}{
		// Admitted plans.
		{"admit/one-declared-effect", planJSON(readEffect, "false", "null"), ""},
		{"admit/one-effect-with-finish", planJSON(effectJSON("e1", "Ailang.Check", "worktree", "1", `{"op":"ai_check","path":"a.ail"}`), "true", "null"), ""},
		{"admit/zero-effect-refusal", planJSON("", "false", `{"ok":false,"refused":"unknown argument key"}`), ""},
		{"admit/id-16-bytes", planJSON(effectJSON("abcdefghijklmn09", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), ""},
		{"admit/payload-exactly-1MiB", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "1", payloadOfCanonicalSize(1<<20)), "false", "null"), ""},

		// Shape (codec): not a world/effect-plan/v1 object.
		{"shape/not-json", `{"plan":`, LawShape},
		{"shape/array-root", `[]`, LawShape},
		{"shape/duplicate-key", `{"plan":"world/effect-plan/v1","plan":"x","effects":[],"finish":false,"result":{"ok":true}}`, LawShape},
		{"shape/wrong-discriminator", strings.Replace(planJSON("", "false", `{"ok":true}`), "effect-plan/v1", "effect-plan/v2", 1), LawShape},
		{"shape/unknown-top-key", strings.Replace(planJSON("", "false", `{"ok":true}`), `"finish"`, `"extra":1,"finish"`, 1), LawShape},
		{"shape/missing-result", `{"plan":"world/effect-plan/v1","effects":[],"finish":false}`, LawShape},
		{"shape/finish-not-bool", planJSON("", "null", `{"ok":true}`), LawShape},
		{"shape/effects-not-array", `{"plan":"world/effect-plan/v1","effects":{},"finish":false,"result":{"ok":true}}`, LawShape},
		{"shape/effect-unknown-key", planJSON(strings.Replace(readEffect, `"cost"`, `"order":1,"cost"`, 1), "false", "null"), LawShape},
		{"shape/cost-not-integer", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "1.5", `{}`), "false", "null"), LawShape},
		{"shape/cost-string", planJSON(effectJSON("e1", "Workspace.Read", "worktree", `"1"`, `{}`), "false", "null"), LawShape},

		// L1: at most one effect in v1.
		{"L1/two-effects", planJSON(readEffect+","+effectJSON("e2", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), LawL1},

		// L2: id grammar [a-z0-9]{1,16} (uniqueness is unreachable through
		// parsePlan while L1 = 1; its law is pinned by the drift test).
		{"L2/uppercase-id", planJSON(effectJSON("E1", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), LawL2},
		{"L2/hyphen-id", planJSON(effectJSON("e-1", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), LawL2},
		{"L2/empty-id", planJSON(effectJSON("", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), LawL2},
		{"L2/id-17-bytes", planJSON(effectJSON("abcdefghijklmnopq", "Workspace.Read", "worktree", "1", `{}`), "false", "null"), LawL2},

		// L3: the exact (effect, scope, cost) triple is declared.
		{"L3/undeclared-effect", planJSON(effectJSON("e1", "Workspace.Write", "worktree", "1", `{}`), "false", "null"), LawL3},
		{"L3/undeclared-scope", planJSON(effectJSON("e1", "Workspace.Read", "/", "1", `{}`), "false", "null"), LawL3},
		{"L3/undeclared-cost", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "0", `{}`), "false", "null"), LawL3},

		// L4: payload is a JSON object, canonical ≤ 1 MiB.
		{"L4/payload-array", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "1", `[]`), "false", "null"), LawL4},
		{"L4/payload-null", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "1", `null`), "false", "null"), LawL4},
		{"L4/payload-over-1MiB", planJSON(effectJSON("e1", "Workspace.Read", "worktree", "1", payloadOfCanonicalSize(1<<20+1)), "false", "null"), LawL4},

		// L6: `world` is reserved.
		{"L6/world-in-result", planJSON("", "false", `{"ok":false,"world":{}}`), LawL6},
		{"L6/world-top-level", strings.Replace(planJSON("", "false", `{"ok":true}`), `"finish"`, `"world":{},"finish"`, 1), LawL6},

		// L7: zero effects ⇔ non-null result object; finish needs an effect.
		{"L7/zero-effects-null-result", planJSON("", "false", "null"), LawL7},
		{"L7/effects-with-result", planJSON(readEffect, "false", `{"ok":true}`), LawL7},
		{"L7/result-not-object", planJSON("", "false", `"refused"`), LawL7},
		{"L7/finish-without-effect", planJSON("", "true", `{"ok":true}`), LawL7},
	}
	seen := map[PlanLaw]bool{}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			plan, err := parsePlan([]byte(row.output), planDeclared)
			if row.law == "" {
				if err != nil {
					t.Fatalf("parsePlan refused a lawful plan: %v", err)
				}
				return
			}
			var lawErr *PlanLawError
			if !errors.As(err, &lawErr) {
				t.Fatalf("parsePlan = (%+v, %v), want a *PlanLawError for %s", plan.Effects, err, row.law)
			}
			if lawErr.Law != row.law {
				t.Fatalf("refused under %s (%v), want %s", lawErr.Law, lawErr, row.law)
			}
			seen[row.law] = true
		})
	}
	for _, law := range []PlanLaw{LawShape, LawL1, LawL2, LawL3, LawL4, LawL6, LawL7} {
		if !seen[law] {
			t.Errorf("no refusal row exercised %s", law)
		}
	}
}

// TestParsePlanAdmittedShape pins what an admitted plan carries: the matched
// requirement, the canonically re-encoded payload, no result, and the
// canonical plan bytes.
func TestParsePlanAdmittedShape(t *testing.T) {
	out := ` { "result":null, "finish":false, "effects":[` + readEffect + `], "plan":"world/effect-plan/v1" } `
	plan, err := parsePlan([]byte(out), planDeclared)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Effects) != 1 || plan.Finish || plan.Result != nil {
		t.Fatalf("plan = %+v", plan)
	}
	e := plan.Effects[0]
	if e.ID != "e1" || e.Requirement != planDeclared[0] {
		t.Fatalf("effect = %+v", e)
	}
	if got := string(e.Payload); got != `{"op":"read","path":"src/main.ail"}` {
		t.Fatalf("payload = %s, want canonical key order", got)
	}
	want := `{"effects":[{"cost":1,"effect":"Workspace.Read","id":"e1","payload":{"op":"read","path":"src/main.ail"},"scope":"worktree"}],"finish":false,"plan":"world/effect-plan/v1","result":null}`
	if string(plan.Canonical) != want {
		t.Fatalf("canonical plan = %s\nwant            %s", plan.Canonical, want)
	}

	refusal, err := parsePlan([]byte(planJSON("", "false", `{"refused":"x","ok":false}`)), planDeclared)
	if err != nil {
		t.Fatal(err)
	}
	if len(refusal.Effects) != 0 || string(refusal.Result) != `{"ok":false,"refused":"x"}` {
		t.Fatalf("refusal plan = %+v (result %s)", refusal, refusal.Result)
	}
}

// TestParseFinishLaws covers L6 for the finish phase.
func TestParseFinishLaws(t *testing.T) {
	rows := []struct {
		name   string
		output string
		law    PlanLaw
	}{
		{"admit/object", `{"passed":true,"error_count":0,"errors":[]}`, ""},
		{"shape/not-json", `{`, LawShape},
		{"L6/not-object", `[1]`, LawL6},
		{"L6/requests-effects", `{"effects":[]}`, LawL6},
		{"L6/world-key", `{"passed":true,"world":{}}`, LawL6},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			got, err := parseFinish([]byte(row.output))
			if row.law == "" {
				if err != nil || string(got) != `{"error_count":0,"errors":[],"passed":true}` {
					t.Fatalf("parseFinish = (%s, %v)", got, err)
				}
				return
			}
			var lawErr *PlanLawError
			if !errors.As(err, &lawErr) || lawErr.Law != row.law {
				t.Fatalf("parseFinish error = %v, want *PlanLawError %s", err, row.law)
			}
		})
	}
}
