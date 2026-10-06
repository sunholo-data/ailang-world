package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// quickstartCommandLines returns the `/tmp/world-publish` and
// `/tmp/ailang-worldd` command lines of a QUICKSTART section, backslash
// continuations joined and whitespace normalised.
func quickstartCommandLines(section string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(section, "\\\n", " "), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/tmp/world-publish ") || strings.HasPrefix(line, "/tmp/ailang-worldd ") {
			out = append(out, strings.Join(strings.Fields(line), " "))
		}
	}
	return out
}

// QUICKSTART §0 (the attended-steps box, 2026-10-06) is a copy of §9's two
// attended commands, and the phrase it tells you to type is the one the
// command will ask for. §9's lines are bound to the CLI by
// TestQuickstartSection9FlagsMatchTheCLI, so binding §0 to §9 binds it to the
// CLI too. And no section still tells an operator to append `< /dev/tty`.
func TestQuickstartAttendedStepsBoxMatchesSection9(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	guide := string(raw)
	_, box, ok := strings.Cut(guide, "## 0. Attended steps (the only commands a person must run)")
	if !ok {
		t.Fatal("QUICKSTART §0 (attended steps) is absent")
	}
	box, _, _ = strings.Cut(box, "## 1. ")
	_, s9, ok := strings.Cut(guide, "### 9. Software-engineering tools")
	if !ok {
		t.Fatal("QUICKSTART §9 is absent")
	}
	s9, _, _ = strings.Cut(s9, "### 10. ")

	got := quickstartCommandLines(box)
	if len(got) != 2 || !strings.HasPrefix(got[0], "/tmp/world-publish transitions ") ||
		!strings.HasPrefix(got[1], "/tmp/ailang-worldd session new ") {
		t.Fatalf("§0 commands = %q, want the publish then the session new", got)
	}
	in9 := map[string]bool{}
	for _, line := range quickstartCommandLines(s9) {
		in9[line] = true
	}
	for _, line := range got {
		if !in9[line] {
			t.Errorf("§0 command %q is not §9's", line)
		}
	}

	manifest, err := os.ReadFile("../../packages/se-tools/transitions.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries []json.RawMessage
	if err := json.Unmarshal(manifest, &entries); err != nil {
		t.Fatal(err)
	}
	phrase := fmt.Sprintf("`publish %d transitions to /tmp/se-world/world.db`", len(entries))
	for name, text := range map[string]string{"§0": box, "§9": s9} {
		if !strings.Contains(text, phrase) {
			t.Errorf("%s does not give the phrase the command asks for, %s", name, phrase)
		}
	}
	for _, stale := range []string{"append `< /dev/tty`", "append </dev/tty", "stdin-is-not-the-controlling-terminal`)"} {
		if strings.Contains(guide, stale) {
			t.Errorf("QUICKSTART still says %q", stale)
		}
	}
}
