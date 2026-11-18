package main

// Row 243 V3 sprint 5 M14 (design §6 CQ1; the S7 rule that QUICKSTART describes what the code does):
// QUICKSTART §20 runs the flag the CLI parses, quotes the lines finish prints and the namespace the
// card signs under, and states the wait the code waits.

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

func TestQuickstartSection20MatchesTheFinishFlag(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	guide := string(raw)
	s19, s20 := strings.Index(guide, "\n## 19. "), strings.Index(guide, "\n## 20. Confirm a finish's pull request with your key *(attended — pending first verbatim run)*\n")
	if s19 < 0 || s20 < s19 {
		t.Fatalf("§20 is missing or not after §19 (%d, %d)", s19, s20)
	}
	section := guide[s20:]
	_, block, _ := strings.Cut(section, "```bash\n")
	block, _, _ = strings.Cut(block, "```")
	argv, err := verbs.Split(strings.TrimSpace(block))
	if err != nil || len(argv) < 3 || argv[0] != "ailang-worldd" || argv[1] != "episode" || argv[2] != "finish" || argv[len(argv)-1] != "--confirm-on-surface" {
		t.Fatalf("§20's command %q (%v)", block, err)
	}
	if !strings.Contains(episodeFinishHelp, "--confirm-on-surface") {
		t.Error("finish's help does not list the flag §20 runs")
	}
	src, err := os.ReadFile("episode_finish_surface.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, quoted := range []string{"confirm request %s on the surface (expires %s, message hash %s)", "the typed y is not read in this mode", "confirmed by %s (decision %s)"} {
		shown := strings.NewReplacer("%s on", "vr-<16 hex> on", "expires %s", "expires <time>", "hash %s", "hash <16 hex>", "by %s", "by <approver>", "decision %s", "decision sha256:…").Replace(quoted)
		if !strings.Contains(string(src), quoted) || !strings.Contains(section, shown) {
			t.Errorf("§20 does not quote finish's line %q as %q", quoted, shown)
		}
	}
	if !strings.Contains(section, "-n "+approval.VerbNamespace) || approval.VerbRequestLifetime != 10*time.Minute || !strings.Contains(section, "up to 10 minutes (and 10 seconds' grace)") || verbLandGrace != 10*time.Second {
		t.Error("§20's namespace, wait or grace is not the code's")
	}
}
