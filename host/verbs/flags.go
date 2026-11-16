package verbs

import (
	"errors"
	"flag"
	"fmt"
	"io"
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
