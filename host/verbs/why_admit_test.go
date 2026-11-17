package verbs

// Row 243 V2b (design §2.2 VA1, §2.3 VA2): why's one flag set and parse, and the bar's admission.

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

// VA1: why's flag set exactly, its two positional rules and its --scan bound.
func TestWhyFlags(t *testing.T) {
	fs := flag.NewFlagSet("x", flag.ContinueOnError)
	WhyFlags(fs)
	got := map[string]string{}
	fs.VisitAll(func(f *flag.Flag) { got[f.Name] = f.DefValue + " " + f.Usage })
	if want := map[string]string{"expect-store": " refuse unless the daemon serves this store (World-Provenance store= path)",
		"result": " a 'call --json-out' output file", "scan": "500 entries to scan back from head (max 5000)",
		"json": "false print the chain as JSON"}; !reflect.DeepEqual(got, want) {
		t.Errorf("flags = %v", got)
	}
	for _, c := range []struct {
		args   []string
		target string
		rest   int
	}{{[]string{"3", "--json"}, "3", 1}, {[]string{"-", "--scan", "2"}, "-", 2}, {[]string{"--json", "3"}, "", 2}, {[]string{"-3"}, "", 1}, {nil, "", 0}} {
		if tg, rest := WhyLeadingTarget(c.args); tg != c.target || len(rest) != c.rest {
			t.Errorf("WhyLeadingTarget(%q) = %q %q", c.args, tg, rest)
		}
	}
	for _, c := range []struct {
		args         []string
		target, want string
	}{{[]string{"3"}, "", "3"}, {nil, "head", "head"}, {[]string{"2"}, "1", `unexpected argument "2"`}, {[]string{"1", "2"}, "", `unexpected argument "2"`}} {
		fs := flag.NewFlagSet("x", flag.ContinueOnError)
		_ = fs.Parse(c.args)
		tg, err := WhyTrailingTarget(fs, c.target)
		if (err == nil && tg != c.want) || (err != nil && err.Error() != c.want) {
			t.Errorf("WhyTrailingTarget(%q, %q) = %q, %v", c.args, c.target, tg, err)
		}
	}
	if WhyScanOK(0) || !WhyScanOK(1) || !WhyScanOK(5000) || WhyScanOK(5001) || WhyScanDefault != 500 || WhyScanMax != 5000 {
		t.Error("WhyScanOK bound")
	}
}

// VA2: one row per argument form of design §2.3, the reason exactly; the target, --scan and --json
// of every admitted form.
func TestAdmitWhy(t *testing.T) {
	const at = "; at a terminal: ailang-worldd why "
	ref := "sha256:" + strings.Repeat("ab", 32)
	for _, c := range []struct {
		line, target string
		scan         int
		json         bool
		reason       string
	}{
		{"", "", 500, false, ""}, // no target: the route's why.ParseTarget answers "a target is required"
		{"3", "3", 500, false, ""},
		{"head", "head", 500, false, ""},
		{ref, ref, 500, false, ""},
		{"sha256:zz", "sha256:zz", 500, false, ""}, // malformed: why.ParseTarget's, at the route
		{"a2a:e:t", "a2a:e:t", 500, false, ""},
		{"mcp:e:t", "mcp:e:t", 500, false, ""},
		{"rest:x", "rest:x", 500, false, ""},
		{"x", "x", 500, false, ""},
		{"head --json", "head", 500, true, ""},
		{"--json head", "head", 500, true, ""},
		{"a2a:e:t --scan 1", "a2a:e:t", 1, false, ""},
		{"--scan 5000 a2a:e:t", "a2a:e:t", 5000, false, ""},
		{"--scan 0 head", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
		{"head --scan 5001", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
		{"head --scan -1", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
		{"1 2", "", 0, false, `ailang-worldd why: unexpected argument "2"`},
		{"--json 1 2", "", 0, false, `ailang-worldd why: unexpected argument "2"`},
		{"-3", "", 0, false, "flag provided but not defined: -3"},
		{"--bogus", "", 0, false, "flag provided but not defined: -bogus"},
		{"--scan x head", "", 0, false, `invalid value "x" for flag -scan: parse error`},
		{"-h", "", 0, false, "see ailang-worldd why --help at a terminal"},
		{"--help", "", 0, false, "see ailang-worldd why --help at a terminal"},
		{"head --help --result f", "", 0, false, "see ailang-worldd why --help at a terminal"},
		{"-", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "-"},
		{"- --scan 3", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "- --scan 3"},
		{"--result f", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "--result f"},
		{"--result 'a b'", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "--result 'a b'"}, // VB8: Quote's form, not a plain join
		{"--expect-store p head", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "--expect-store p head"},
		{"- --result f --expect-store p", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "- --result f --expect-store p"},
		{"- --expect-store p", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "- --expect-store p"},
		{"--result '' head", "head", 500, false, ""}, // an empty --result is no --result, as at a terminal
	} {
		args, err := Split(c.line)
		if err != nil {
			t.Fatal(err)
		}
		target, scan, asJSON, reason := AdmitWhy(args)
		if reason != c.reason || target != c.target || scan != c.scan || asJSON != c.json {
			t.Errorf("AdmitWhy(%q) = %q %d %v %q; want %q %d %v %q", c.line, target, scan, asJSON, reason, c.target, c.scan, c.json, c.reason)
		}
		if v, _ := Lookup(append([]string{"why"}, args...)); strings.Contains(reason, "; at a terminal: ") && !strings.HasSuffix(reason, "; at a terminal: "+v.CLI) {
			t.Errorf("AdmitWhy(%q): the CLI line %q is not Lookup's %q", c.line, reason, v.CLI)
		}
	}
}
