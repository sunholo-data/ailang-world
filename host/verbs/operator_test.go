package verbs

// Row 243 V4 S1 (design §5.1's grammar and identifier rows): the admitted operator subset. Every line
// is cut by Split, as the bar cuts it; every refusal asserts its class AND that no tuple is admitted.
// The expectations are this file's literals, never the production tables.

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var allYes = OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}

// admitLine splits line and admits it; Intent must be set exactly for a review.
func admitLine(t *testing.T, line string, p OperatorPolicy) OperatorAnswer {
	t.Helper()
	argv, err := Split(line)
	if err != nil {
		t.Fatalf("Split(%q): %v", line, err)
	}
	return admitArgv(t, argv, p)
}

func admitArgv(t *testing.T, argv []string, p OperatorPolicy) OperatorAnswer {
	t.Helper()
	a := AdmitOperator(argv, p)
	if (a.Class == OperatorReview) != (a.Intent != nil) {
		t.Fatalf("AdmitOperator(%q) = class %q with intent %+v", argv, a.Class, a.Intent)
	}
	return a
}

// refused requires class c, field (when not "-") and a reason containing has, with no tuple.
func refused(t *testing.T, a OperatorAnswer, what string, c OperatorClass, field, has string) {
	t.Helper()
	if a.Class != c || a.Intent != nil || (field != "-" && a.Field != field) || !strings.Contains(a.Reason, has) {
		t.Errorf("%s: got class %q field %q reason %q intent %+v; want %q %q containing %q, no intent", what, a.Class, a.Field, a.Reason, a.Intent, c, field, has)
	}
}

const (
	start  = "episode start w1 --preset p --launch --headless"
	finish = "episode finish w1 --preset p"
)

func TestOperatorGrammarCase(t *testing.T) {
	if a := admitLine(t, start, allYes); a.Class != OperatorReview {
		t.Fatalf("control: %q is %q", start, a.Class)
	}
	for _, line := range []string{"Episode start w1 --preset p --launch --headless", "episode Start w1 --preset p --launch --headless",
		"EPISODE START w1 --preset p --launch --headless", "episode FINISH w1 --preset p", "episode Finish w1 --preset p"} {
		argv, _ := Split(line)
		v, _ := Lookup(argv)
		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
	}
	for _, line := range []string{"episode start w1 --Preset p --launch --headless", "episode start w1 --preset p --LAUNCH --headless",
		"episode start w1 --preset p --launch --Headless", "episode start w1 --PRESET=p --launch --headless",
		"episode finish w1 --preset p --Ready", "episode finish w1 --preset p --Confirm-On-Surface"} {
		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag provided but not defined")
	}
	refused(t, admitLine(t, "episode start W1 --preset p --launch --headless", allYes), "W1", OperatorUsage, "<episode>", "does not match")
	refused(t, admitLine(t, "episode start w1 --preset P --launch --headless", allYes), "P", OperatorUsage, "--preset", "not a launch preset name")
	refused(t, admitLine(t, "episode finish wA --preset p", allYes), "wA", OperatorUsage, "<episode>", "does not match")
}

func TestOperatorGrammarArity(t *testing.T) {
	for _, line := range []string{"episode start --preset p --launch --headless", "episode start w1 w2 --preset p --launch --headless",
		"episode start w1 --preset p --launch --headless extra", "episode start --preset p w1 --launch w2 --headless",
		"episode finish --preset p", "episode finish w1 w2 --preset p", "episode finish w1 --preset p --ready x"} {
		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "<episode>", "needs exactly one <episode>")
	}
	for _, line := range []string{"episode start w1 --launch --headless", "episode start w1 --preset '' --launch --headless",
		"episode finish w1", "episode finish w1 --preset= --ready"} {
		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "--preset", "--preset <name> is required")
	}
	for _, line := range []string{"episode start w1 --launch --headless --preset", "episode finish w1 --preset", "episode finish w1 --ready --preset"} {
		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag needs an argument")
	}
	refused(t, admitLine(t, "episode start '' --preset p --launch --headless", allYes), "empty episode", OperatorUsage, "<episode>", "does not match")
}

func TestOperatorGrammarExactWords(t *testing.T) {
	for _, line := range []string{"episode starts w1 --preset p --launch --headless", "episode star w1 --preset p --launch --headless",
		"episode start-x w1 --preset p --launch --headless", "episodes start w1 --preset p --launch --headless",
		"xepisode start w1 --preset p --launch --headless", "episode finis w1 --preset p", "episode finished w1 --preset p",
		"episode launch w1 --preset p --headless", "episode report w1", "episode", "start w1 --preset p --launch --headless"} {
		argv, _ := Split(line)
		v, _ := Lookup(argv)
		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
	}
	for _, x := range []string{"--pres p", "--presets p", "--preset-dir p", "--head", "--headles", "--headlessx", "--launc", "--launched",
		"--read", "--readyy", "--confirm", "--confirm-on-surface-x", "--dry", "--ttl-x 1"} {
		line := start + " " + x
		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag provided but not defined")
	}
	refused(t, admitLine(t, "episode start w1 ---preset p --launch --headless", allYes), "---preset", OperatorUsage, "", "bad flag syntax")
}

func TestOperatorGrammarCanonical(t *testing.T) {
	head := []string{Binary, "episode", "start", "w1", "--preset", "p", "--launch", "--headless"}
	cases := []struct {
		lines []string
		want  OperatorIntent
	}{
		{[]string{start, "ailang-worldd " + start, "episode start --preset=p --launch --headless w1", "episode start -preset p -launch -headless w1",
			"episode start --headless --launch w1 --preset 'p'", `episode start "w1" --preset="p" --launch=true --headless=1`,
			"episode start --launch --preset p w1 --headless", "episode start w1 --preset p --launch=T --headless=TRUE"},
			OperatorIntent{Kind: "start", Episode: "w1", Preset: "p", Headless: true, Argv: head}},
		{[]string{"episode start w1 --preset p --launch", "episode start --launch=true w1 -preset=p"},
			OperatorIntent{Kind: "start", Episode: "w1", Preset: "p", Argv: head[:7]}},
		{[]string{finish, "episode finish --preset=p w1", "ailang-worldd episode finish -preset p w1"},
			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p"}}},
		{[]string{"episode finish --confirm-on-surface w1 --preset p", "episode finish w1 --preset p --confirm-on-surface=true"},
			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", ConfirmOnSurface: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface"}}},
		{[]string{"episode finish w1 --ready --preset p --confirm-on-surface", "episode finish --confirm-on-surface=1 --ready=true w1 --preset=p"},
			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", ConfirmOnSurface: true, Ready: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface", "--ready"}}},
		{[]string{"episode finish w1 --ready --preset my-preset-2"},
			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "my-preset-2", Ready: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "my-preset-2", "--ready"}}},
	}
	for _, c := range cases {
		for _, line := range c.lines {
			a := admitLine(t, line, allYes)
			if a.Class != OperatorReview || !reflect.DeepEqual(*a.Intent, c.want) {
				t.Errorf("%q: class %q intent %+v, want %+v", line, a.Class, a.Intent, c.want)
				continue
			}
			argv, _ := Split(line)
			if v, _ := Lookup(argv); a.CLI != v.CLI || v.CLI == "" {
				t.Errorf("%q: CLI %q", line, a.CLI)
			}
			if again := admitArgv(t, a.Intent.Argv, allYes); again.Class != OperatorReview || !reflect.DeepEqual(*again.Intent, c.want) {
				t.Errorf("%q: the canonical line %q re-admits as %+v", line, a.Intent.Argv, again.Intent)
			}
		}
	}
}

func TestOperatorGrammarDuplicates(t *testing.T) {
	for _, c := range []struct{ line, field string }{
		{"episode start w1 --preset p --preset p --launch --headless", "--preset"},
		{"episode start w1 --preset p -preset=q --launch --headless", "--preset"},
		{"episode start w1 --preset=p --launch --launch --headless", "--launch"},
		{"episode start w1 --preset p --launch=true --launch --headless", "--launch"},
		{"episode start w1 --preset p --launch --headless -headless", "--headless"},
		{"episode finish w1 --preset p --ready --ready", "--ready"},
		{"episode finish w1 --preset p --confirm-on-surface --confirm-on-surface=true", "--confirm-on-surface"},
		{"episode finish w1 --preset p --title a --title b", "--title"},
	} {
		refused(t, admitLine(t, c.line, allYes), c.line, OperatorUsage, c.field, "is given more than once")
	}
}

func TestOperatorGrammarOptions(t *testing.T) {
	startOut := map[string]bool{"db": false, "credentials-db": false, "workspace-root": false, "repo": false, "branch": false, "profile": false,
		"exec-profiles-dir": false, "task": false, "agent": false, "model": false, "ttl": false, "presets-dir": false, "resume": false,
		"hosts-dir": false, "request": true, "grant": false, "budget": false, "sessions-dir": false, "addr": false}
	finishOut := map[string]bool{"presets-dir": false, "title": false, "include-unattributed": true, "exec-profiles-dir": false,
		"dry-run": true, "allow-failing": false}
	for _, set := range []struct {
		base string
		out  map[string]bool
	}{{start, startOut}, {finish, finishOut}} {
		for name, isBool := range set.out {
			x, eq := "--"+name+" v", "--"+name+"=v" // =v is not a boolean: a value flag mirrored as a bool fails here
			if isBool {
				x, eq = "--"+name, "--"+name+"=true"
			}
			for _, line := range []string{set.base + " " + x, strings.Replace(set.base, " w1 ", " "+x+" w1 ", 1), set.base + " " + eq} {
				a := admitLine(t, line, allYes)
				refused(t, a, line, OperatorTerminal, "--"+name, "--"+name+" is not run from the app; at a terminal: "+a.CLI)
				if verb := strings.Fields(set.base)[1]; !strings.HasPrefix(a.Reason, "ailang-worldd episode "+verb+": --"+name+" ") {
					t.Errorf("%q: reason %q does not lead with the command", line, a.Reason)
				}
			}
		}
	}
	for _, c := range []struct{ line, field string }{
		{"episode start w1 --preset p --launch=false --headless", "--launch"},
		{"episode start w1 --preset p --launch --headless=false", "--headless"},
		{"episode start w1 --preset p --launch=f --headless=0", "--launch"},
		{"episode finish w1 --preset p --ready=false", "--ready"},
		{"episode finish w1 --preset p --confirm-on-surface=FALSE", "--confirm-on-surface"},
		{"episode finish w1 --preset p --dry-run=false", "--dry-run"},
	} {
		has := "is given false"
		if c.field == "--dry-run" {
			has = "is not run from the app"
		}
		refused(t, admitLine(t, c.line, allYes), c.line, OperatorTerminal, c.field, has)
	}
	refused(t, admitLine(t, "episode start w1 --preset p", allYes), "no --launch", OperatorTerminal, "--launch", "only with --launch; at a terminal: ")
	refused(t, admitLine(t, "episode start w1 --preset p --headless", allYes), "standalone --headless", OperatorUsage, "--headless", "--headless requires --launch")
	refused(t, admitLine(t, "episode start w1 --preset p --launch=maybe", allYes), "not a bool", OperatorUsage, "",
		`ailang-worldd episode start: invalid boolean value "maybe" for -launch: parse error`) // the flag package's words
	refused(t, admitLine(t, "episode start w1 --preset se-tools --launch --headless", allYes), "se-tools", OperatorUsage, "--preset", "not a launch preset name")
	// Every boolean spelling decides as strconv.ParseBool (the flag package) does: true admits, false is
	// refused as false, anything else is a parse error.
	for _, w := range []string{"1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false", "False", "", "2", "yes", "tRUE", "true ", "on"} {
		a, line := admitLine(t, "episode start w1 --preset p --launch --headless='"+w+"'", allYes), "--headless="+w
		switch b, err := strconv.ParseBool(w); {
		case err != nil:
			refused(t, a, line, OperatorUsage, "", "invalid boolean value")
		case b:
			if a.Class != OperatorReview {
				t.Errorf("%s: %q %q", line, a.Class, a.Reason)
			}
		default:
			refused(t, a, line, OperatorTerminal, "--headless", "is given false")
		}
	}
	for _, h := range []string{"-h", "--help", "-help"} {
		refused(t, admitLine(t, start+" "+h, allYes), h, OperatorUsage, "--help", "--help at a terminal")
	}
}
