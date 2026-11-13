package main

// Row 201 A4 sprint 4 (design §2.8, §10 step 4; S7): QUICKSTART §18's card paragraph names the
// card's own controls, the card's route and an approve --list line the CLI documents.

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestQuickstartSection18CardMatchesThePage(t *testing.T) {
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	page, err := os.ReadFile("../../host/surfacepage/decide/index.html")
	if err != nil {
		t.Fatal(err)
	}
	routes, err := os.ReadFile("../../host/daemon/daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "## 18. Decide an agent's approval request")
	if !ok {
		t.Fatal("QUICKSTART has no §18")
	}
	_, card, ok := strings.Cut(section, "**From the surface.**")
	if !ok {
		t.Fatal("§18 has no card paragraph")
	}
	labels := regexp.MustCompile(`\*\*([^*]+)\*\*`).FindAllStringSubmatch(card, -1)
	if len(labels) < 6 {
		t.Fatalf("the card paragraph names %d controls, want the card's six and Needs you", len(labels))
	}
	for _, l := range labels {
		if l[1] != "Needs you" && !strings.Contains(string(page), ">"+l[1]+"<") {
			t.Errorf("§18 names **%s**, which the card page does not show", l[1])
		}
	}
	if !strings.Contains(card, "/surface/decide/<request-ref>") || !strings.Contains(string(routes), `"GET /surface/decide/{ref}"`) {
		t.Error("§18 does not name the card's route, or the daemon does not serve it")
	}
	lists := 0
	for _, line := range strings.Split(card, "\n") {
		if !strings.HasPrefix(line, "ailang-worldd approve --list") {
			continue
		}
		lists++
		for _, flag := range regexp.MustCompile(`\s(--[a-z-]+)`).FindAllStringSubmatch(line, -1) {
			if !strings.Contains(approveHelp, flag[1]+" ") && !strings.Contains(approveHelp, flag[1]+"]") {
				t.Errorf("§18 uses %s, which approve's help does not document: %s", flag[1], line)
			}
		}
	}
	if lists != 1 {
		t.Errorf("the card paragraph holds %d approve --list lines, want 1", lists)
	}
}
