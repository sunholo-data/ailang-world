package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/plugindoc"
)

// Row 189 M5 (w-world-onboarding.md §3.3, AC6): every ailang-worldd line the
// ailang-world plugin prints names a verb in the top-level usage and only
// flags that verb's own --help names.

const pluginDir = "../../plugins/ailang-world"

// pluginWorlddLines is the number of ailang-worldd lines the plugin holds:
// the walk must find at least this many, so an empty or broken walk fails.
const pluginWorlddLines = 16

// checkWorlddLine checks one plugin command line against this binary's help.
func checkWorlddLine(text string) error {
	toks := strings.Fields(strings.TrimPrefix(text, "ailang-worldd "))
	// The global --addr before the verb (--addr X or --addr=X).
	for len(toks) > 0 && strings.HasPrefix(toks[0], "--addr") {
		if toks[0] == "--addr" {
			if len(toks) < 2 {
				return fmt.Errorf("the global --addr has no value")
			}
			toks = toks[2:]
		} else {
			toks = toks[1:]
		}
	}
	path := plugindoc.VerbPath(toks, plugindoc.UsagePaths(usage, "  ailang-worldd ", "[--addr <url>] "))
	if path == nil {
		return fmt.Errorf("verb %q is not in ailang-worldd's top-level usage", plugindoc.Verb(toks))
	}
	var out, errw bytes.Buffer
	if code := run(append(append([]string(nil), path...), "--help"), &out, &errw); code != exitOK || out.Len() == 0 {
		return fmt.Errorf("ailang-worldd %s --help: exit %d, no help on stdout (stderr %q)", strings.Join(path, " "), code, errw.String())
	}
	for _, flag := range plugindoc.Flags(toks[len(path):]) {
		if !plugindoc.FlagIn(out.String(), flag) {
			return fmt.Errorf("flag %s is not in 'ailang-worldd %s --help'", flag, strings.Join(path, " "))
		}
	}
	return nil
}

func TestPluginCommandsMatchTheCLIHelp(t *testing.T) {
	lines, err := plugindoc.Walk(pluginDir, "ailang-worldd")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < pluginWorlddLines {
		t.Fatalf("control: the plugin walk found %d ailang-worldd lines, want at least %d: %v", len(lines), pluginWorlddLines, lines)
	}
	for _, l := range lines {
		if err := checkWorlddLine(l.Text); err != nil {
			t.Errorf("%s:%d: %q: %v", l.File, l.Line, l.Text, err)
		}
	}

	// Controls: the same checker refuses an unknown verb, an unknown
	// subcommand of a known family and an unknown flag, and accepts the
	// global --addr before the verb.
	for _, tc := range []struct{ line, want string }{
		{"ailang-worldd episode frobnicate w1", `verb "episode frobnicate" is not in`},
		{"ailang-worldd profile frobnicate --repo .", `verb "profile frobnicate" is not in`},
		{"ailang-worldd frobnicate", `verb "frobnicate" is not in`},
		{"ailang-worldd episode report w1 --no-such-flag", "flag --no-such-flag is not in 'ailang-worldd episode report --help'"},
		{"ailang-worldd --addr 127.0.0.1:47330 why 3 --json=x --nope", "flag --nope is not in 'ailang-worldd why --help'"},
		{"ailang-worldd --addr 127.0.0.1:47330 why 3 --json", ""},
		{"ailang-worldd profile init --repo . --project p", ""},
	} {
		err := checkWorlddLine(tc.line)
		switch {
		case tc.want == "" && err != nil:
			t.Errorf("control: %q refused: %v", tc.line, err)
		case tc.want != "" && (err == nil || !strings.Contains(err.Error(), tc.want)):
			t.Errorf("control: %q: err %v, want %q", tc.line, err, tc.want)
		}
	}
	// The extraction finds a fenced line, a continued fenced line, an inline
	// span and a "$ " / VAR=value prefix, and nothing in prose.
	md := "prose ailang-worldd health\n\n```sh\n$ FOO=1 ailang-worldd episode report w1 \\\n  --no-such-flag\n```\nsee `ailang-worldd why 3 --json` here\n"
	got := plugindoc.Lines("x.md", md, "ailang-worldd")
	want := []plugindoc.Line{{File: "x.md", Line: 4, Text: "ailang-worldd episode report w1 --no-such-flag"},
		{File: "x.md", Line: 7, Text: "ailang-worldd why 3 --json"}}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("control: extraction %v, want %v", got, want)
	}
}

// The top-level usage lists profile init (the check above demands it for the
// plugin's profile init lines).
func TestUsageListsProfileInit(t *testing.T) {
	if !strings.Contains(usage, "\n  ailang-worldd profile init --repo <dir> --project <name>") {
		t.Error("the top-level usage does not list profile init")
	}
}
