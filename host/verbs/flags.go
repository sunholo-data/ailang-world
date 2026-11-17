package verbs

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

// ParseInterleaved parses fs over args, allowing flags and positional
// arguments in any order (Go's flag package stops at the first positional).
func ParseInterleaved(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return positional, nil
		}
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

// EpisodeReportArgs are episode report's flags.
type EpisodeReportArgs struct {
	Markdown, JSON, Follow *bool
	Interval               *time.Duration
	Addr                   *string
}

// EpisodeReportFlags registers episode report's five flags on fs; the CLI and the bar share it.
func EpisodeReportFlags(fs *flag.FlagSet) *EpisodeReportArgs {
	return &EpisodeReportArgs{
		Markdown: fs.Bool("markdown", false, "render a scrubbed PR body"),
		JSON:     fs.Bool("json", false, "print the route's JSON"),
		Follow:   fs.Bool("follow", false, "re-print when the episode changes, until Ctrl-C"),
		Interval: fs.Duration("interval", 2*time.Second, "poll interval with --follow"),
		Addr:     fs.String("addr", "", "the daemon (default: the global --addr)"),
	}
}

// barRefusal is the bar profile's reason for each flag the bar does not run.
var barRefusal = []struct{ flag, why string }{
	{"follow", "the app shows one report; it re-reads on Enter"},
	{"interval", "the app shows one report; it re-reads on Enter"},
	{"markdown", "the PR body is rendered by episode finish"},
	{"addr", "the app reads its own daemon"},
}

// AdmitEpisodeReport runs episode report's flag set and parse over args, then the bar profile:
// --json is allowed, every other flag and -h/--help refused, exactly one episode required.
func AdmitEpisodeReport(args []string) (episode string, asJSON bool, reason string) {
	fs := flag.NewFlagSet(Binary+" episode report", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	a := EpisodeReportFlags(fs)
	pos, err := ParseInterleaved(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		return "", false, "see " + Binary + " episode report --help at a terminal"
	}
	if err != nil {
		return "", false, Binary + " episode report: " + err.Error()
	}
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	for _, r := range barRefusal {
		if set[r.flag] {
			return "", false, fmt.Sprintf("--%s is not run from the app: %s", r.flag, r.why)
		}
	}
	if len(pos) != 1 || pos[0] == "" {
		return "", false, Binary + " episode report needs exactly one <episode>"
	}
	return pos[0], *a.JSON, ""
}

// WhyScanDefault and WhyScanMax are why's --scan default and bound (host/why's ScanDefault, ScanMax).
const (
	WhyScanDefault = 500
	WhyScanMax     = 5000
)

// WhyArgs are why's flags.
type WhyArgs struct {
	ExpectStore, Result *string
	Scan                *int
	JSON                *bool
}

// WhyFlags registers why's four flags on fs; the CLI and the bar share it.
func WhyFlags(fs *flag.FlagSet) *WhyArgs {
	return &WhyArgs{
		ExpectStore: fs.String("expect-store", "", "refuse unless the daemon serves this store (World-Provenance store= path)"),
		Result:      fs.String("result", "", "a 'call --json-out' output file"),
		Scan:        fs.Int("scan", WhyScanDefault, "entries to scan back from head (max 5000)"),
		JSON:        fs.Bool("json", false, "print the chain as JSON"),
	}
}

// WhyLeadingTarget takes a leading target: a first argument that is not a flag, or exactly -.
func WhyLeadingTarget(args []string) (target string, rest []string) {
	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
		return args[0], args[1:]
	}
	return "", args
}

// WhyTrailingTarget takes the one positional left after the flags, unless a target was already given.
func WhyTrailingTarget(fs *flag.FlagSet, target string) (string, error) {
	if fs.NArg() == 0 {
		return target, nil
	}
	if target != "" || fs.NArg() > 1 {
		return "", fmt.Errorf("unexpected argument %q", fs.Arg(fs.NArg()-1))
	}
	return fs.Arg(0), nil
}

// WhyScanOK reports whether n is a --scan the CLI runs.
func WhyScanOK(n int) bool { return n >= 1 && n <= WhyScanMax }

// whyBarRefusal is the bar profile's reason for each of why's reads the app does not make (design §2.3),
// in the order the bar checks them.
var whyBarRefusal = []struct{ what, why string }{
	{"--result", "the app reads no file on the daemon's host"},
	{"--expect-store", "the app reads its own daemon's store"},
	{"why -", "the app has no stdin"},
}

// AdmitWhy runs why's flag set and parse over args in the CLI's order (leading target, flags, trailing
// target, the --scan bound), then the bar profile: --result, --expect-store and - are refused with the
// CLI line. The target is then why.ParseTarget's (the route's), as at a terminal.
func AdmitWhy(args []string) (target string, scan int, asJSON bool, reason string) {
	cli := Quote(append([]string{Binary, "why"}, args...))
	target, rest := WhyLeadingTarget(args)
	fs := flag.NewFlagSet(Binary+" why", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	a := WhyFlags(fs)
	if err := fs.Parse(rest); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return "", 0, false, "see " + Binary + " why --help at a terminal"
		}
		return "", 0, false, err.Error()
	}
	target, err := WhyTrailingTarget(fs, target)
	if err != nil {
		return "", 0, false, Binary + " why: " + err.Error()
	}
	if !WhyScanOK(*a.Scan) {
		return "", 0, false, fmt.Sprintf("%s why: --scan must be 1..%d", Binary, WhyScanMax)
	}
	for i, used := range []bool{*a.Result != "", *a.ExpectStore != "", target == "-"} {
		if used {
			r := whyBarRefusal[i]
			return "", 0, false, r.what + " is not run from the app: " + r.why + "; at a terminal: " + cli
		}
	}
	return target, *a.Scan, *a.JSON, ""
}
