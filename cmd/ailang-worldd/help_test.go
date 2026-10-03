package main

import (
	"bytes"
	"strings"
	"testing"
)

// everyVerb is every verb line of the CLI, as argv.
var everyVerb = [][]string{
	{"serve"}, {"health"}, {"head"}, {"world", "get"}, {"object", "get"}, {"object", "find"},
	{"log", "get"}, {"log", "range"}, {"log", "tail"}, {"registry", "get"}, {"commit"},
	{"tools", "list"}, {"call"}, {"why"}, {"provenance"}, {"setup"}, {"doctor"},
	{"session"}, {"session", "new"}, {"session", "list"}, {"session", "revoke"}, {"session", "mint"},
}

// AC5.7: `<verb> --help` and `help <verb>` exit 0 with that verb's help on
// stdout for every verb, and the top-level usage is complete.
func TestEveryVerbHasHelp(t *testing.T) {
	for _, verb := range everyVerb {
		name := strings.Join(verb, " ")
		for _, argv := range [][]string{append(append([]string{}, verb...), "--help"), append([]string{"help"}, verb...)} {
			var out, errw bytes.Buffer
			code := run(argv, &out, &errw)
			if code != exitOK {
				t.Errorf("%q: exit %d, want 0 (stderr %q)", argv, code, errw.String())
				continue
			}
			if !strings.Contains(out.String(), "ailang-worldd") || !strings.Contains(out.String(), verb[0]) {
				t.Errorf("%q printed no help for %s:\n%s", argv, name, out.String())
			}
		}
	}
	var out, errw bytes.Buffer
	if code := run([]string{"help"}, &out, &errw); code != exitOK {
		t.Fatalf("help: exit %d", code)
	}
	for _, want := range []string{
		"serve --db", "] health", "] head", "world get", "object get", "object find", "log get", "log range",
		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("top-level usage lacks %q", want)
		}
	}
	// The duplicated default ("(default 3600) (default 3600)") is gone.
	out.Reset()
	run([]string{"session", "mint", "--help"}, &out, &errw)
	var flags bytes.Buffer
	runSessionMint([]string{"-h"}, &out, &flags, testMintEnv("", 0))
	if strings.Contains(flags.String()+out.String(), "(default 3600) (default 3600)") {
		t.Error("session mint still prints its default twice")
	}
}
