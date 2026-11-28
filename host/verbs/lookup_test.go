package verbs

// Row 243 V1 (design §6 VB3): the verb table and approve's navigate arm.

import (
	"reflect"
	"strings"
	"testing"
)

var lookupCases = []struct {
	line, words string
	class       Class
	reason      string
	cli         bool
	rest        []string
}{
	{"row-189@6", "", Address, "", false, nil},
	{"world", "", Address, "", false, nil},
	{"back", "", Address, "", false, nil},
	{"back x", "", Address, "", false, nil},
	{"world why", "", Address, "", false, nil},
	{"world session", "", Address, "", false, nil},
	{"world get", "", Address, "", false, nil},
	{"World a b c", "", Address, "", false, nil},
	{"ailang-worldd", "", Address, "", false, nil},
	{"why", "why", Read, "", true, []string{}},
	{"why e1 --json", "why", Read, "", true, []string{"e1", "--json"}},
	{"episode report w1", "episode report", Read, "", true, []string{"w1"}},
	{"ailang-worldd episode report w1 --json", "episode report", Read, "", true, []string{"w1", "--json"}},
	{"profile list", "profile list", Read, "", true, []string{}},
	{"profile list --json", "profile list", Read, "", true, []string{"--json"}},
	{"preset list", "preset list", Read, "", true, []string{}},
	{"preset list --json", "preset list", Read, "", true, []string{"--json"}},
	{"app list", "app list", Read, "", true, []string{}},
	{"app list --json", "app list", Read, "", true, []string{"--json"}},
	{"approve mr-0123456789abcdef", "approve", Navigate, "", true, []string{"mr-0123456789abcdef"}},
	{"episode start w1", "episode start", Mutate, "episode start changes World. The app cannot run it yet.", true, []string{"w1"}},
	{"episode finish w1 --preset p", "episode finish", Mutate, "episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key.", true, []string{"w1", "--preset", "p"}},
	{"episode launch w1", "episode launch", Host, "needs a host process", true, []string{"w1"}},
	{"allowance grant --preset p", "allowance", Never, "allowance is never run from the app: granting and revoking a standing allowance", true, []string{"grant", "--preset", "p"}},
	{"approver enroll --name x", "approver", Never, "approver is never run from the app: enrolling", true, []string{"enroll", "--name", "x"}},
	{"session frob", "session", Never, "minting and revoking a session", true, []string{"frob"}},
	{"serve", "serve", Never, "host administration", true, []string{}},
	{"log", "log", Other, "log is a CLI command the app does not run; at a terminal: ailang-worldd log", true, []string{}},
	{"log tail", "log", Other, "at a terminal: ailang-worldd log tail", true, []string{"tail"}},
	{"run world-agy fix bug", "run", Other, "run is a CLI command the app does not run; at a terminal: ailang-worldd run world-agy fix bug", true, []string{"world-agy", "fix", "bug"}},
	{"ailang-worldd world get world-1", "world get", Other, "world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1", true, []string{"world-1"}},
	{"world get a b", "world get", Other, "world get is a CLI command", true, []string{"a", "b"}},
	{"--addr x episode report w1", "", Refused, "global flags are not accepted", true, nil},
	{"ailang-worldd -h", "", Refused, "global flags are not accepted", true, nil},
	{"episode", "episode", Refused, EpisodeNeedsSubcommand, true, []string{}},
	{"episode frobnicate", "episode", Refused, EpisodeNeedsSubcommand, true, []string{"frobnicate"}},
	{"episode w1", "episode", Refused, EpisodeNeedsSubcommand, true, []string{"w1"}},
	{"world a b c", "world", Refused, AddressHelp, false, []string{"a", "b", "c"}},
}

// VB3: every row of the table and of the trace (design §2.1(b), §2.4), the shadow set, step 4b.
func TestLookup(t *testing.T) {
	for _, c := range lookupCases {
		argv, err := Split(c.line)
		if err != nil {
			t.Fatal(err)
		}
		if back, err := Split(Quote(argv)); err != nil || !reflect.DeepEqual(back, argv) { // VB2 over every VB3 argv
			t.Errorf("Split(Quote(%q)) = %q, %v", argv, back, err)
		}
		v, rest := Lookup(argv)
		cli := strings.TrimPrefix(c.line, "ailang-worldd ")
		if v.Words != c.words || v.Class != c.class || !strings.Contains(v.Reason, c.reason) || (c.reason == "" && v.Reason != "") ||
			(v.CLI != "") != c.cli || (c.cli && v.CLI != "ailang-worldd "+cli) || !reflect.DeepEqual(rest, c.rest) {
			t.Errorf("Lookup(%q) = %+v, %q", c.line, v, rest)
		}
	}
	// The shadow set: every one of the 31 words is a verb, bare and with an argument; never an address.
	all := append(append(append([]string{}, OperatorWords...), NeverWords...), OtherWords...)
	if len(all) != 31 {
		t.Fatalf("shadow set has %d words, want 31", len(all))
	}
	for _, w := range all {
		for _, argv := range [][]string{{w}, {w, "zzz"}, {Binary, w, "zzz"}} {
			if v, _ := Lookup(argv); v.Class == Address || v.CLI == "" {
				t.Errorf("Lookup(%q) = %+v: a shadow-set word navigates or has no CLI line", argv, v)
			}
		}
		if v, _ := Lookup([]string{"world", w}); v.Class != Address {
			t.Errorf("world %s = %+v: a shadowed name is not reachable as world NAME", w, v)
		}
	}
}

// ApproveTarget: an id or a ref opens its card; anything else is refused.
func TestApproveTarget(t *testing.T) {
	ref := "sha256:" + strings.Repeat("ab", 32)
	for args, want := range map[string]string{"mr-0123456789abcdef": "/surface/approve/mr-0123456789abcdef", ref: "/surface/decide/" + ref,
		"mr-0123456789abcdef --deny": "", "--list": "", "": "", ref + "0": "", "sha256:" + strings.Repeat("AB", 32): "", "mr-0123456789abcdef x": "", "mr-0123456789ABCDEF": "",
		"vr-0123456789abcdef": "/surface/confirm/vr-0123456789abcdef", "vr-0123": "", "vr-0123456789abcdef0": "", "vr-0123456789ABCDEF": "", "xvr-0123456789abcdef": "",
		"vr-0123456789abcdef x": "", "vr-0123456789abcdef --deny": "", "--deny vr-0123456789abcdef": ""} {
		href, reason := ApproveTarget(strings.Fields(args))
		if href != want || (href == "") == (reason == "") {
			t.Errorf("ApproveTarget(%q) = %q, %q", args, href, reason)
		}
	}
}
