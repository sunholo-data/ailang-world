package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/plugindoc"
)

// Row 189 M5 (w-world-onboarding.md §3.3, AC6): every world-publish line the
// ailang-world plugin prints names a verb in usageText and only flags that
// verb's own usage names.

// pluginPublishLines is the number of world-publish lines the plugin holds:
// the walk must find at least this many, so an empty or broken walk fails.
const pluginPublishLines = 2

// verbUsage is the verb's block of usageText: its line and the continuation
// lines under it, up to the next verb or a blank line.
func verbUsage(path []string) string {
	head := "  world-publish " + strings.Join(path, " ") + " "
	var block []string
	in := false
	for _, l := range strings.Split(usageText, "\n") {
		switch {
		case strings.HasPrefix(l, "  world-publish "):
			in = strings.HasPrefix(l+" ", head)
		case strings.TrimSpace(l) == "" || !strings.HasPrefix(l, "   "):
			in = false
		}
		if in {
			block = append(block, l)
		}
	}
	return strings.Join(block, "\n")
}

// checkPublishLine checks one plugin command line against usageText.
func checkPublishLine(text string) error {
	toks := strings.Fields(strings.TrimPrefix(text, "world-publish "))
	path := plugindoc.VerbPath(toks, plugindoc.UsagePaths(usageText, "  world-publish "))
	if path == nil {
		return fmt.Errorf("verb %q is not in world-publish's usage", plugindoc.Verb(toks))
	}
	help := verbUsage(path)
	if help == "" {
		return fmt.Errorf("world-publish %s has no usage block", strings.Join(path, " "))
	}
	for _, flag := range plugindoc.Flags(toks[len(path):]) {
		if !plugindoc.FlagIn(help, flag) {
			return fmt.Errorf("flag %s is not in world-publish %s's usage", flag, strings.Join(path, " "))
		}
	}
	return nil
}

func TestPluginCommandsMatchTheCLIHelp(t *testing.T) {
	lines, err := plugindoc.Walk("../../plugins/ailang-world", "world-publish")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < pluginPublishLines {
		t.Fatalf("control: the plugin walk found %d world-publish lines, want at least %d: %v", len(lines), pluginPublishLines, lines)
	}
	for _, l := range lines {
		if err := checkPublishLine(l.Text); err != nil {
			t.Errorf("%s:%d: %q: %v", l.File, l.Line, l.Text, err)
		}
	}

	// Controls: the same checker refuses an unknown verb and a flag outside
	// the verb's own usage (--live is publish's, not reconcile's), and
	// accepts a continuation-line flag and the two-word verb.
	for _, tc := range []struct{ line, want string }{
		{"world-publish frobnicate --store s", `verb "frobnicate" is not in`},
		{"world-publish epoch frobnicate --store s", `verb "epoch frobnicate" is not in`},
		{"world-publish transitions --store s --no-such-flag", "flag --no-such-flag is not in world-publish transitions's usage"},
		{"world-publish reconcile --store s --live", "flag --live is not in world-publish reconcile's usage"},
		{"world-publish publish --store s --credential-file c --live", ""},
		{"world-publish epoch nominate --store s --binary b --candidate", ""},
	} {
		err := checkPublishLine(tc.line)
		switch {
		case tc.want == "" && err != nil:
			t.Errorf("control: %q refused: %v", tc.line, err)
		case tc.want != "" && (err == nil || !strings.Contains(err.Error(), tc.want)):
			t.Errorf("control: %q: err %v, want %q", tc.line, err, tc.want)
		}
	}
}
