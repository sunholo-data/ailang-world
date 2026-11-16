package verbs

// Row 243 V1 (design §6 VB5): episode report's one flag set and the bar's admission.

import (
	"flag"
	"reflect"
	"strings"
	"testing"
	"time"
)

// VB5: the CLI's flag set exactly, and the bar profile's one allow and its refusals.
func TestEpisodeReportFlags(t *testing.T) {
	fs := flag.NewFlagSet("x", flag.ContinueOnError)
	EpisodeReportFlags(fs)
	got := map[string]string{}
	fs.VisitAll(func(f *flag.Flag) { got[f.Name] = f.DefValue })
	if want := map[string]string{"markdown": "false", "json": "false", "follow": "false", "interval": (2 * time.Second).String(), "addr": ""}; !reflect.DeepEqual(got, want) {
		t.Errorf("flags = %v, want %v", got, want)
	}
	for _, c := range []struct {
		args          []string
		episode       string
		json          bool
		reasonContain string
	}{
		{[]string{"w1"}, "w1", false, ""},
		{[]string{"--json", "w1"}, "w1", true, ""},
		{[]string{"w1", "--json=false"}, "w1", false, ""},
		{[]string{"w1", "--follow"}, "", false, "--follow is not run from the app: the app shows one report"},
		{[]string{"w1", "--interval", "5s"}, "", false, "--interval is not run from the app"},
		{[]string{"w1", "--markdown"}, "", false, "--markdown is not run from the app: the PR body"},
		{[]string{"w1", "--addr", "http://x"}, "", false, "--addr is not run from the app: the app reads its own daemon"},
		{[]string{"w1", "-h"}, "", false, "see ailang-worldd episode report --help"},
		{[]string{"--help"}, "", false, "see ailang-worldd episode report --help"},
		{[]string{"w1", "--frob"}, "", false, "flag provided but not defined: -frob"},
		{[]string{"w1", "--interval", "x"}, "", false, "invalid value"},
		{[]string{}, "", false, "needs exactly one <episode>"},
		{[]string{"w1", "w2"}, "", false, "needs exactly one <episode>"},
		{[]string{""}, "", false, "needs exactly one <episode>"},
	} {
		ep, asJSON, reason := AdmitEpisodeReport(c.args)
		if ep != c.episode || asJSON != c.json || !strings.Contains(reason, c.reasonContain) || (c.reasonContain == "") != (reason == "") {
			t.Errorf("AdmitEpisodeReport(%q) = %q, %v, %q", c.args, ep, asJSON, reason)
		}
	}
}

// ParseInterleaved keeps the CLI's order-free parse (its CLI callers' tests pin the rest).
func TestParseInterleaved(t *testing.T) {
	fs := flag.NewFlagSet("x", flag.ContinueOnError)
	j := fs.Bool("json", false, "")
	pos, err := ParseInterleaved(fs, []string{"a", "--json", "b"})
	if err != nil || !*j || !reflect.DeepEqual(pos, []string{"a", "b"}) {
		t.Errorf("= %q %v %v", pos, *j, err)
	}
}
