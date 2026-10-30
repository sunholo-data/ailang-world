package main

import (
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/plugindoc"
)

func quickstartServiceLines(md string) []plugindoc.Line {
	return plugindoc.Lines("../../docs/QUICKSTART.md", md, "ailang-worldd")
}

func TestQuickstartServiceSectionMatchesTheCLIHelp(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, md, found := strings.Cut(string(raw), "## Run World as a service")
	if !found {
		t.Fatal("missing service section")
	}
	md, _, found = strings.Cut(md, "## Run your agents through World")
	if !found {
		t.Fatal("missing end of service section")
	}
	lines := plugindoc.Lines("../../docs/QUICKSTART.md", md, "ailang-worldd")
	if len(lines) < 10 {
		t.Fatalf("found only %d command lines", len(lines))
	}
	base := useWorldErrors(quickstartServiceLines(md))
	if len(base) != 0 {
		t.Fatal(base)
	}
	for _, mutant := range []string{"ailang-worldd frobnicate", "ailang-worldd why 3 --expect-shop x", "ailang-worldd service status --bogus"} {
		for _, form := range []string{"\n```sh\n" + mutant + "\n```\n", "\n`" + mutant + "`\n"} {
			errs := useWorldErrors(quickstartServiceLines(md + form))
			if len(errs) != len(base)+1 || !strings.Contains(strings.Join(errs, "\n"), mutant) {
				t.Errorf("mutant %q: %v", mutant, errs)
			}
		}
	}
}
