package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/plugindoc"
)

// Row 189 M6 (w-world-onboarding.md §8, AC7): the "Use World on your repo"
// page is the script the attended AC7 run follows verbatim, so every
// ailang-worldd line in it is checked against this binary's help with the
// plugin's checker (checkWorlddLine). Lines starting with "claude " are not
// ailang-worldd lines and are not checked.

const useWorldPage = "../../docs/USE_WORLD_ON_YOUR_REPO.md"

// useWorldPageLines is the number of ailang-worldd lines the page holds: the
// walk must find at least this many, so an empty or broken walk fails.
const useWorldPageLines = 10

// useWorldLines is plugindoc.Lines over the page, plus every line outside a
// fence that itself starts with "ailang-worldd ", so a command pasted onto
// the page as bare text is checked too.
func useWorldLines(md string) []plugindoc.Line {
	lines := plugindoc.Lines(useWorldPage, md, "ailang-worldd")
	inFence := false
	for i, l := range strings.Split(md, "\n") {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence && strings.HasPrefix(t, "ailang-worldd ") {
			lines = append(lines, plugindoc.Line{File: useWorldPage, Line: i + 1, Text: strings.Join(strings.Fields(t), " ")})
		}
	}
	return lines
}

// useWorldErrors checks every line and returns one message per refusal.
func useWorldErrors(lines []plugindoc.Line) []string {
	var out []string
	for _, l := range lines {
		if err := checkWorlddLine(l.Text); err != nil {
			out = append(out, fmt.Sprintf("%s:%d: %q: %v", l.File, l.Line, l.Text, err))
		}
	}
	return out
}

func TestUseWorldPageCommandsMatchTheCLIHelp(t *testing.T) {
	raw, err := os.ReadFile(useWorldPage)
	if err != nil {
		t.Fatal(err)
	}
	md := string(raw)
	lines := plugindoc.Lines(useWorldPage, md, "ailang-worldd")
	if len(lines) < useWorldPageLines {
		t.Fatalf("control: the page walk found %d ailang-worldd lines, want at least %d: %v", len(lines), useWorldPageLines, lines)
	}
	base := useWorldErrors(useWorldLines(md))
	for _, e := range base {
		t.Error(e)
	}

	// Controls: each mutant appended to the page, fenced or bare, adds one
	// refusal to the same check, naming the mutant.
	for _, mutant := range []string{
		"ailang-worldd episode frobnicate w1",
		"ailang-worldd episode start w1 --no-such-flag",
		"ailang-worldd session revoke --episode w1 --bogus",
	} {
		for _, form := range []string{"\n```sh\n" + mutant + "\n```\n", "\n" + mutant + "\n"} {
			errs := useWorldErrors(useWorldLines(md + form))
			if len(errs) != len(base)+1 || !strings.Contains(strings.Join(errs, "\n"), fmt.Sprintf("%q", mutant)) {
				t.Errorf("control: appending %q gave %d refusal(s) %v, want %d, one naming it", form, len(errs), errs, len(base)+1)
			}
		}
	}
}
