package daemon

// Row 243 V1 (design §2.4, §6 VR2): every outcome class end to end, and the page's replayed answers.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

// verbGolden is the route's real answers the page's node test replays (VP6): regenerate with
// WORLD_UPDATE_VERB_GOLDEN=1; otherwise every answer must equal it byte for byte.
var verbGolden = filepath.Join("..", "surfacepage", "testdata", "verb-route.golden.json")

// VR2: each outcome end to end through the handler, and the page's replayed answers are these.
func TestSurfaceVerbOutcomes(t *testing.T) {
	a, _ := verbRig(t)
	const mr, ref = "mr-0123456789abcdef", "sha256:abababababababababababababababababababababababababababababababab"
	for _, c := range []struct{ q, outcome, class, verb, cli, href, reason string }{
		{"row-189@6", "address", "address", "", "", "", ""},
		{"world shared@0ffbb9c", "address", "address", "", "", "", ""},
		{"back", "address", "address", "", "", "", ""},
		{"world why", "address", "address", "", "", "", ""},
		{"world session", "address", "address", "", "", "", ""},
		{"world get", "address", "address", "", "", "", ""},
		{"world episode", "address", "address", "", "", "", ""},
		{"why", "refused", "read", "why", "ailang-worldd why", "", "ailang-worldd why: a target is required"},
		{"episode finish w1 --preset p", "terminal", "mutate", "episode finish", "ailang-worldd episode finish w1 --preset p", "", "changes World"},
		{"episode start w1", "terminal", "mutate", "episode start", "ailang-worldd episode start w1", "", "changes World"},
		{"episode launch w1 --preset p", "terminal", "host", "episode launch", "ailang-worldd episode launch w1 --preset p", "", "needs a host process"},
		{"approve " + mr, "navigate", "navigate", "approve", "ailang-worldd approve " + mr, "/surface/approve/" + mr, ""},
		{"ailang-worldd approve " + ref, "navigate", "navigate", "approve", "ailang-worldd approve " + ref, "/surface/decide/" + ref, ""},
		{"approve " + mr + " --deny", "refused", "navigate", "approve", "ailang-worldd approve " + mr + " --deny", "", "the app opens approve"},
		{"approver enroll --name x", "refused", "never", "approver", "ailang-worldd approver enroll --name x", "", "never run from the app"},
		{"session frob", "refused", "never", "session", "ailang-worldd session frob", "", "never run from the app"},
		{"log", "refused", "other", "log", "ailang-worldd log", "", "log is a CLI command the app does not run; at a terminal: ailang-worldd log"},
		{"log tail", "refused", "other", "log", "ailang-worldd log tail", "", "CLI command"},
		{"ailang-worldd world get world-1", "refused", "other", "world get", "ailang-worldd world get world-1", "", "world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1"},
		{"episode", "refused", "refused", "episode", "ailang-worldd episode", "", verbs.EpisodeNeedsSubcommand},
		{"episode frobnicate", "refused", "refused", "episode", "ailang-worldd episode frobnicate", "", verbs.EpisodeNeedsSubcommand},
		{"episode w1", "refused", "refused", "episode", "ailang-worldd episode w1", "", verbs.EpisodeNeedsSubcommand},
		{"world a b c", "refused", "refused", "world", "", "", verbs.AddressHelp},
		{"--addr x episode report w1", "refused", "refused", "", "ailang-worldd --addr x episode report w1", "", "global flags are not accepted"},
		{"episode report w1 --follow", "refused", "read", "episode report", "ailang-worldd episode report w1 --follow", "", "--follow is not run from the app"},
		{"episode report w1 --markdown", "refused", "read", "episode report", "ailang-worldd episode report w1 --markdown", "", "--markdown is not run from the app"},
		{"episode report " + strings.Repeat("e", 257), "refused", "read", "episode report", "ailang-worldd episode report " + strings.Repeat("e", 257), "", "episode must be 1..256 bytes"},
		{"episode report w1; rm -rf ~", "refused", "refused", "", "", "", "Shell:"},
		{"episode report nobody", "result", "read", "episode report", "ailang-worldd episode report nobody", "", ""},
		{"episode report w1", "result", "read", "episode report", "ailang-worldd episode report w1", "", ""},
		{"ailang-worldd episode report w1 --json", "result", "read", "episode report", "ailang-worldd episode report w1 --json", "", ""},
	} {
		v, _ := verbGet(a, c.q)
		if v.Outcome != c.outcome || v.Class != c.class || v.Verb != c.verb || v.CLI != c.cli || v.Href != c.href || !strings.Contains(v.Reason, c.reason) || (c.reason == "" && v.Reason != "") {
			t.Errorf("%q = %+v, want outcome %s class %s verb %q cli %q href %q reason %q", c.q, *v, c.outcome, c.class, c.verb, c.cli, c.href, c.reason)
		}
		if (v.Text != nil || v.JSON != nil || v.Stderr != nil) != (c.outcome == "result") {
			t.Errorf("%q: result members on a %s answer", c.q, v.Outcome)
		}
	}
	if v, _ := verbGet(a, "episode report nobody"); v.Stderr == nil || *v.Stderr != "ailang-worldd episode report: episode nobody has no committed calls\n" || v.Exit != 1 || v.Text != nil {
		t.Errorf("not found = %+v", *v)
	}
	if v, _ := verbGet(a, "episode report w1"); v.Text == nil || !strings.Contains(*v.Text, "episode w1: 2 calls") || !strings.Contains(*v.Text, `#2 <img src=x onerror=alert(1)>: <b>hostile</b> "reason"`) || v.JSON != nil || v.Exit != 0 {
		t.Errorf("text = %+v", *v)
	}
	if v, _ := verbGet(a, "episode report w1 --json"); v.JSON == nil || !strings.HasPrefix(*v.JSON, `{"episode":"w1",`) || !strings.HasSuffix(*v.JSON, "}\n") || strings.Count(*v.JSON, "\n") != 1 {
		t.Errorf("json = %+v", *v)
	}
	// The page's node test (VP6, VP1–VP3) replays these real answers, never hand-typed ones.
	golden := map[string]string{}
	for _, q := range []string{"row-189@6", "world shared@0ffbb9c", "back", "world why", "why", "log", "world session",
		"ailang-worldd world get world-1", "episode", "world episode", "world a b c", "episode finish w1 --preset p",
		"approve mr-0123456789abcdef", "approve " + ref, "episode report nobody", "episode report w1", "episode report w1 --json",
		strings.Repeat("x", verbs.MaxLine+1), "why 2", "why head --json", "why 9", "why --result f"} {
		_, raw := verbGet(a, q)
		golden[q] = raw
	}
	want, err := json.MarshalIndent(golden, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want = append(want, '\n')
	if os.Getenv("WORLD_UPDATE_VERB_GOLDEN") == "1" {
		if err := os.WriteFile(verbGolden, want, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got, err := os.ReadFile(verbGolden); err != nil || !bytes.Equal(got, want) {
		t.Errorf("%s is not the route's answers (WORLD_UPDATE_VERB_GOLDEN=1 regenerates it): %v", verbGolden, err)
	}
}
