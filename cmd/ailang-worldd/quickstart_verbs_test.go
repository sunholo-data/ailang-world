package main

// Row 243 V1 (design §7, §10): QUICKSTART §19 says what the bar does with each line it shows, and
// host/verbs, the page and the daemon agree.

import (
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

func TestQuickstartSection19MatchesTheBar(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	guide := string(raw)
	s18, s19 := strings.Index(guide, "\n## 18. "), strings.Index(guide, "\n## 19. Operator verbs in the address bar *(attended — pending first verbatim run)*\n")
	if s18 < 0 || s19 < s18 {
		t.Fatalf("§19 is missing or not after §18 (%d, %d)", s18, s19)
	}
	section := guide[s19:]
	_, block, _ := strings.Cut(section, "```text\n")
	block, _, _ = strings.Cut(block, "```")
	want := map[string]verbs.Class{"episode report": verbs.Read, "why": verbs.Read, "approve": verbs.Navigate, "episode finish": verbs.Mutate, "": verbs.Address}
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
		argv, err := verbs.Split(line)
		if err != nil {
			t.Fatalf("§19 line %q: %v", line, err)
		}
		v, _ := verbs.Lookup(argv)
		if c, ok := want[v.Words]; !ok || c != v.Class {
			t.Errorf("§19 line %q is %q %s, which §19 does not describe", line, v.Words, v.Class)
		}
		seen[v.Words] = true
	}
	if len(seen) != len(want) {
		t.Errorf("§19's lines cover %v, want each of %v", seen, want)
	}
	for _, w := range verbs.NeverWords {
		if !strings.Contains(section, "`"+w+"`") {
			t.Errorf("§19 does not name %s among the commands that stay at a terminal", w)
		}
	}
	if !strings.Contains(section, "15 seconds") || !strings.Contains(string(surfacepage.AppScript), "const VERB_ROUTE_BUDGET_MS = 15000;") {
		t.Error("§19's 15 seconds is not the page's budget")
	}
	if !strings.Contains(string(surfacepage.AppShell), "episode report") || !strings.Contains(section, "GET /surface/verb.json") {
		t.Error("the bar's placeholder or §19 does not name the verb route's first verb")
	}
}
