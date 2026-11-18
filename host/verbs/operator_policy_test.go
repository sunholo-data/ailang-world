package verbs

// Row 243 V4 S1 (design §5.1's policy, never-dispatch, identifier and other-verb rows), on operator_test.go's
// helpers: the answers as an input, the families never admitted, the three identifier grammars at each
// boundary class, and every other line keeping Lookup's classification.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// policyOracle is design §4.1/§10's answer table, written out here (not read from the code).
func policyOracle(p OperatorPolicy, what string) bool {
	ab := func(a Answer) bool { return a == "A" || a == "B" }
	switch what {
	case "headless":
		return ab(p.Start)
	case "interactive":
		return p.Start == "A" && p.HandoffProven
	case "finish":
		return ab(p.Start) && ab(p.Finish)
	}
	return ab(p.Start) && p.Finish == "A" // ready
}

func TestOperatorGrammarPolicy(t *testing.T) {
	lines := []struct{ what, line, field string }{
		{"headless", start, "--headless"},
		{"interactive", "episode start w1 --preset p --launch", "--launch"},
		{"finish", finish + " --confirm-on-surface", ""},
		{"ready", finish + " --ready", ""},
	}
	n, admitted := 0, 0
	for _, s := range []Answer{"", "A", "B", "C", "a", "b", "D", "A "} {
		for _, f := range []Answer{"", "A", "B", "C", "a", "D"} {
			for _, h := range []bool{false, true} {
				p := OperatorPolicy{Start: s, Finish: f, HandoffProven: h}
				for _, l := range lines {
					a, want := admitLine(t, l.line, p), policyOracle(p, l.what)
					n++
					if want {
						admitted++
						if a.Class != OperatorReview {
							t.Errorf("%+v %s: %q %q, want review", p, l.what, a.Class, a.Reason)
						}
						continue
					}
					field := l.field
					if l.what == "ready" && policyOracle(p, "finish") {
						field = "--ready"
					}
					refused(t, a, l.what+" "+string(s)+"/"+string(f), OperatorTerminal, field, "; at a terminal: ")
				}
			}
		}
	}
	if n != 384 || admitted != 42 {
		t.Fatalf("%d cases, %d admitted; want 384 and 42", n, admitted)
	}
	for _, l := range lines { // the zero value answers nothing and admits nothing
		refused(t, admitLine(t, l.line, OperatorPolicy{}), "zero policy "+l.what, OperatorTerminal, "-", "D-WORLD-1")
	}
}

func TestOperatorNeverDispatch(t *testing.T) {
	for _, line := range []string{"approver enroll --name mark-key", "approver revoke mark-key", "session mint w1", "session revoke w1",
		"service install", "service", "serve", "setup", "ailang-worldd approver enroll", "ailang-worldd service uninstall",
		"world-publish approve sha256:x", "world-publish"} {
		a := admitLine(t, line, allYes)
		refused(t, a, line, OperatorNever, "-", "is never run from the app")
		if !strings.Contains(a.Reason, "; at a terminal: ") || a.CLI == "" {
			t.Errorf("%q: reason %q CLI %q", line, a.Reason, a.CLI)
		}
	}
	for _, line := range []string{"sh -c ls", "bash", "episode launch w1 --preset p --headless", "episode report w1", "call x", "log", "why head",
		"approve mr-0123456789abcdef", "world-publishx", "commit"} {
		argv, _ := Split(line)
		v, _ := Lookup(argv)
		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
	}
	for _, line := range []string{start + "; rm -rf x", start + " $(id)", start + " | sh", "episode start w1 --preset `p` --launch --headless"} {
		if argv, err := Split(line); err != ErrShell {
			t.Errorf("Split(%q) = %q, %v; want ErrShell (the bar runs no shell)", line, argv, err)
		}
	}
}

func TestOperatorIdentifierAnchors(t *testing.T) {
	ep := func(id string) OperatorAnswer {
		return admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes)
	}
	pr := func(id string) OperatorAnswer {
		return admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes)
	}
	if ep("w1").Class != OperatorReview || pr("p").Class != OperatorReview || !ValidOperatorActionID("oa-"+strings.Repeat("0a", 16)) {
		t.Fatal("controls: w1, p and oa-0a… must be admitted")
	}
	for _, id := range []string{"_w1", "w1_", ".w1", "w1.", "w1/", "/w1", " w1", "w1 ", "é1", "w1é"} {
		refused(t, ep(id), "episode "+id, OperatorUsage, "<episode>", "does not match")
		refused(t, pr(id), "preset "+id, OperatorUsage, "--preset", "not a launch preset name")
	}
	hex := strings.Repeat("0a", 16)
	for _, id := range []string{"x" + "oa-" + hex, "oa-" + hex + "x", " oa-" + hex, "oa-" + hex + " ", "-oa-" + hex, "oa-" + hex + "-", "oa" + hex, "oa_" + hex, "oa-" + "-" + hex[1:]} {
		if ValidOperatorActionID(id) {
			t.Errorf("action id %q admitted", id)
		}
	}
}

func TestOperatorIdentifierFlags(t *testing.T) {
	hex := strings.Repeat("0a", 16)
	for _, id := range []string{"w1\nzz", "zz\nw1", "w1\n", "\nw1", "W1", "wQ"} {
		refused(t, admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes), "episode "+id, OperatorUsage, "<episode>", "does not match")
		refused(t, admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes), "preset "+id, OperatorUsage, "--preset", "not a launch preset name")
	}
	for _, id := range []string{"oa-" + hex + "\nx", "x\noa-" + hex, "oa-" + hex + "\n", "oa-" + strings.ToUpper(hex), "OA-" + hex, "Oa-" + hex} {
		if ValidOperatorActionID(id) {
			t.Errorf("action id %q admitted", id)
		}
	}
}

func TestOperatorIdentifierLength(t *testing.T) {
	ep := func(id string) OperatorAnswer {
		return admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes)
	}
	pr := func(id string) OperatorAnswer {
		return admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes)
	}
	for _, id := range []string{"a", "0", "a" + strings.Repeat("-", 63), strings.Repeat("a", 64)} {
		if a := ep(id); a.Class != OperatorReview || a.Intent.Episode != id {
			t.Errorf("episode of %d characters: %q %q", len(id), a.Class, a.Reason)
		}
	}
	for _, id := range []string{"", strings.Repeat("a", 65), "a" + strings.Repeat("-", 64), "-"} {
		refused(t, ep(id), fmt.Sprintf("episode of %d characters", len(id)), OperatorUsage, "<episode>", "does not match")
	}
	for _, id := range []string{"p", strings.Repeat("p", 51)} {
		if a := pr(id); a.Class != OperatorReview || a.Intent.Preset != id {
			t.Errorf("preset of %d characters: %q %q", len(id), a.Class, a.Reason)
		}
	}
	for _, id := range []string{strings.Repeat("p", 52), strings.Repeat("p", 64), strings.Repeat("p", 65)} {
		refused(t, pr(id), fmt.Sprintf("preset of %d characters", len(id)), OperatorUsage, "--preset", "at most 51 characters")
	}
	hex := strings.Repeat("0a", 16)
	if !ValidOperatorActionID("oa-"+hex) || len("oa-"+hex) != 35 {
		t.Fatal("control: oa- and 32 hex digits")
	}
	for _, id := range []string{"oa-" + hex[:31], "oa-" + hex + "0", "oa-", "oa-" + hex[:30], "oa-" + hex + "00"} {
		if ValidOperatorActionID(id) {
			t.Errorf("action id of %d hex digits admitted", len(id)-3)
		}
	}
}

// TestOperatorOtherVerbParity: every line of GET /surface/verb.json's golden corpus that is not episode
// start or episode finish keeps Lookup's classification — never a review, the same reason and CLI line.
func TestOperatorOtherVerbParity(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "surfacepage", "testdata", "verb-route.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]string
	if err := json.Unmarshal(b, &golden); err != nil {
		t.Fatal(err)
	}
	lines := []string{}
	for q := range golden {
		lines = append(lines, q)
	}
	for _, w := range append(append(append([]string{}, OperatorWords...), NeverWords...), OtherWords...) {
		lines = append(lines, w, w+" x")
	}
	n := 0
	for _, line := range lines {
		argv, err := Split(line)
		if err != nil {
			continue
		}
		v, _ := Lookup(argv)
		a := admitArgv(t, argv, allYes)
		if v.Words == "episode start" || v.Words == "episode finish" {
			if a.Class == OperatorOther || a.Class == OperatorNever {
				t.Errorf("%q: %q is an operator line", line, a.Class)
			}
			continue
		}
		want := OperatorOther
		if v.Class == Never {
			want = OperatorNever
		}
		if a.Class != want || a.Reason != v.Reason || a.CLI != v.CLI || a.Field != "" {
			t.Errorf("%q: %+v, want class %q with Lookup's reason %q and CLI %q", line, a, want, v.Reason, v.CLI)
		}
		n++
	}
	if n < 70 {
		t.Fatalf("%d lines kept Lookup's classification; want at least 70", n)
	}
	if a := admitLine(t, "episode finish w1 --preset p", allYes); a.Class != OperatorReview {
		t.Errorf("the golden's finish line under every answer A: %q", a.Class)
	}
}
