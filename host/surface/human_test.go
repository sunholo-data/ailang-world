package surface

import (
	"strings"
	"testing"
)

func TestHumanSources(t *testing.T) {
	h := strings.Repeat("a", 40)
	record := func(subject, body string) string {
		return "\x1e" + h + "\x1f" + h + " " + h + "\x1f2026-10-20T00:00:00Z\x1f" + subject + "\x1f" + body
	}
	merges, err := ParseMerges(record("Merge pull request #1 from o/trial/row1-x", "\nfeat(x): A measured title\nMore text") +
		record("Merge pull request #2 from o/trial/row1-y", "") + record("Merge branch topic", "No title"))
	if err != nil {
		t.Fatal(err)
	}
	if merges[0].Title != "feat(x): A measured title" || merges[1].Title != "" || merges[2].Title != "" {
		t.Fatalf("titles: %+v", merges)
	}
	goals := ChartRowGoals("1. **w-world-onboarding** · **USING WORLD AND AILANG WITH THE CLI IS TOO HARD**\n2. **w-other** · No statement.")
	if goals[1] != "Using World and AILANG with the CLI is too hard." || goals[2] != "" || humanName("w-world-onboarding") != "World onboarding" {
		t.Fatalf("goals: %v", goals)
	}
	first, full := summaries("**I completed the requested work.** Contact a@example.test at /tmp/private or ~/private; keep cmd/tool/x.go. Next.")
	if first != "I completed the requested work." || full != "I completed the requested work. Contact [email] at [path] or [path] keep cmd/tool/x.go. Next." {
		t.Fatalf("summaries: %q / %q", first, full)
	}
	if a, b := summaries(""); a != "" || b != "" {
		t.Fatal("invented summary")
	}
	if len([]rune(capHuman(strings.Repeat("word ", 100), 200))) > 200 {
		t.Fatal("title cap")
	}
}

func TestChartRowGoalsCodeSpans(t *testing.T) {
	charter := "156. **w-worktree** · Literal `**` before **AN AGENT CAN EDIT THE TESTS THAT JUDGE IT: `workspace-exec`'s WHOLE WORKTREE IS WRITABLE**\n" +
		"157. **w-code-only** · `**THIS UPPERCASE STATEMENT IS ONLY INLINE CODE**`"
	goals := ChartRowGoals(charter)
	want := "An agent can edit the tests that judge it: `workspace-exec`'s whole worktree is writable."
	if goals[156] != want {
		t.Fatalf("row 156 goal: got %q, want %q", goals[156], want)
	}
	if goal, ok := goals[157]; ok {
		t.Fatalf("row 157 invented goal: %q", goal)
	}
}

func TestRegradeRecords(t *testing.T) {
	good := `{"schema":"trial-regrade/1","episode":"synthetic","grader":"test","passed":true,"rc":0}`
	for _, tc := range []struct {
		name string
		data string
		want bool
	}{
		{"counts", good, true},
		{"wrong episode", strings.ReplaceAll(good, "synthetic", "other"), false},
		{"failed", strings.ReplaceAll(good, "true", "false"), false},
		{"oversize", good + strings.Repeat(" ", 64<<10), false},
		{"missing rc", strings.ReplaceAll(good, `,"rc":0`, ""), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := RegradePassed([]byte(tc.data), "synthetic", "test"); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHumanValidateRefusals(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*Index)
	}{
		{"title cap", func(x *Index) { x.States[0].Title = strings.Repeat("x", 201) }},
		{"summary cap", func(x *Index) { x.Steps[0].Summary = strings.Repeat("x", 301) }},
		{"full cap", func(x *Index) { x.Steps[0].SummaryFull = strings.Repeat("x", 4001) }},
		{"goal cap", func(x *Index) { x.Branes[0].Goal = strings.Repeat("x", 301) }},
		{"title path", func(x *Index) { x.States[0].Title = "/" + "Users/synthetic/x" }},
		{"summary path", func(x *Index) { x.Steps[0].Summary = "/" + "home/synthetic/x" }},
		{"full email", func(x *Index) { x.Steps[0].SummaryFull = "a@example.test" }},
		{"goal email", func(x *Index) { x.Branes[0].Goal = "a@example.test" }},
		{"passed", func(x *Index) { x.Steps[0].Delivered.Passed = x.Steps[0].Delivered.Total + 1 }},
		{"total", func(x *Index) { x.Steps[0].Delivered.Total++ }},
		{"regraded", func(x *Index) { x.Steps[0].Delivered.Regraded = []string{"never-failed"} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			x, _ := rigIndex(t)
			tc.mutate(&x)
			if Validate(x) == nil {
				t.Fatal("accepted invalid field")
			}
		})
	}
}

func TestRoundHumanFacts(t *testing.T) {
	first := EpisodeRecord{AgentSummary: "I completed the first round successfully. More details."}
	last := EpisodeRecord{AgentSummary: "Wrong round.", Graders: []recordGrader{{Name: "test", Passed: false}, {Name: "vet", Passed: true}}, Regraded: map[string]bool{"test": true}}
	s := stepOf(Episode{ID: "synthetic", Rounds: []EpisodeRecord{first, last}}, State{}, Report{}, false)
	if s.Summary != "I completed the first round successfully." || s.Graders.Passed != 1 || s.Delivered.Passed != 2 || len(s.Failed) != 1 || s.Delivered.Regraded[0] != "test" {
		t.Fatalf("%+v", s)
	}
}

func TestHumanPrivateTextRefusals(t *testing.T) {
	for _, field := range []string{"title", "summary", "summary_full", "goal"} {
		for _, private := range []string{"/" + "Users/synthetic/file", "/" + "home/synthetic/file", "synthetic@example.test"} {
			t.Run(field+private, func(t *testing.T) {
				x, _ := rigIndex(t)
				switch field {
				case "title":
					x.States[0].Title = private
				case "summary":
					x.Steps[0].Summary = private
				case "summary_full":
					x.Steps[0].SummaryFull = private
				case "goal":
					x.Branes[0].Goal = private
				}
				if Validate(x) == nil {
					t.Fatal("accepted private text")
				}
			})
		}
	}
}
