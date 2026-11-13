package main

// Row 201 A4 sprint 3 M6g (design §2.10, §10; S7): QUICKSTART §18 documents the decision mode
// with the CLI's own flags, the decision namespace and the message length the code uses.

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

func TestQuickstartSection18FlagsMatchTheCLI(t *testing.T) {
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "## 18. Decide an agent's approval request")
	if !ok {
		t.Fatal("QUICKSTART has no §18")
	}
	section = strings.ReplaceAll(section, "\\\n", " ")
	seen := map[string]bool{}
	lines := 0
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "ailang-worldd approve sha256:") {
			continue
		}
		lines++
		for _, flag := range regexp.MustCompile(`(?:^|\s)(--[a-z-]+)`).FindAllStringSubmatch(line, -1) {
			seen[flag[1]] = true
			if !strings.Contains(approveHelp, flag[1]+" ") && !strings.Contains(approveHelp, flag[1]+"]") {
				t.Errorf("§18 uses %s, which approve's help does not document: %s", flag[1], line)
			}
		}
	}
	for _, flag := range []string{"--deny", "--sign-with", "--write-message", "--signature", "--expiry", "--webauthn-link", "--webauthn-assertion", "--daemon", "--db"} {
		if !seen[flag] {
			t.Errorf("§18 shows no approve line with %s", flag)
		}
	}
	if lines != 6 {
		t.Errorf("§18 holds %d approve sha256: lines, want 6", lines)
	}
	for _, want := range []string{"-n " + approval.DecisionNamespace, fmt.Sprintf("%d-byte decision message", approval.DecisionMessageLen)} {
		if !strings.Contains(section, want) {
			t.Errorf("§18 does not say %q", want)
		}
	}
}
