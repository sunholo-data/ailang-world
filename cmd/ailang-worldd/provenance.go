package main

// D-WORLD-60 (row 138): `provenance` prints the World-Provenance trailer a
// PR or commit made through World carries —
//
//	World-Provenance: store=<db> episode=<ep> entries=<from>-<to>
//
// — so `ailang-worldd why <entry>` can walk any of those entries back to its
// tool calls, plans and effect records. from/to are the first and last
// coordinator entries of the episode in the scanned range, found with why's
// log scanner.

import (
	"flag"
	"fmt"
	"io"
)

const provenanceHelp = `usage: ailang-worldd [--addr <url>] provenance [--since <entry>] [--episode <ep>] [--scan N]

Prints the World-Provenance trailer (D-WORLD-60) for the coordinator entries
an episode committed:

  World-Provenance: store=<db path from /v1/health> episode=<ep> entries=<from>-<to>

from/to are that episode's first and last entries from --since to head.
--since defaults to head-499; the range may hold at most --scan entries
(default 500, max 5000). Without --episode, one trailer is printed per
episode found in the range, in order of first appearance. Exit 1 when the
range holds no coordinator entry (of that episode).
`

type episodeRange struct {
	episode       string
	from, to, num int64
}

func runProvenance(addr string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd provenance", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	since := fs.Int64("since", -1, "first entry of the range (default head-499)")
	episode := fs.String("episode", "", "only this episode")
	scan := fs.Int("scan", whyScanDefault, "maximum entries in the range (max 5000)")
	if code, done := parseVerbFlags(fs, args, provenanceHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd provenance: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	if *scan < 1 || *scan > whyScanMax {
		fmt.Fprintf(stderr, "ailang-worldd provenance: --scan must be 1..%d\n", whyScanMax)
		return exitUsage
	}
	ctx, cancel := budgetContext(whyBudget)
	defer cancel()
	r := worldReader{ctx: ctx, c: newClient(addr)}
	h, err := r.health()
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd provenance: %v\n", err)
		return exitUsage
	}
	head, ok, err := r.headIndex()
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd provenance: %v\n", err)
		return exitUsage
	}
	if !ok {
		fmt.Fprintln(stderr, "ailang-worldd provenance: no world head has been selected yet: nothing committed")
		return exitUsage
	}
	from := *since
	if from < 0 {
		from = head - int64(*scan) + 1
		if from < 0 {
			from = 0
		}
	}
	if from > head {
		fmt.Fprintf(stderr, "ailang-worldd provenance: --since %d is past head (entry %d)\n", from, head)
		return exitUsage
	}
	if head-from+1 > int64(*scan) {
		fmt.Fprintf(stderr, "ailang-worldd provenance: entries %d-%d hold %d entries, more than --scan %d; narrow --since or raise --scan (max %d)\n",
			from, head, head-from+1, *scan, whyScanMax)
		return exitUsage
	}
	ranges, err := episodeRanges(r, from, head, *scan)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd provenance: %v\n", err)
		return exitUsage
	}
	printed := 0
	for _, er := range ranges {
		if *episode != "" && er.episode != *episode {
			continue
		}
		fmt.Fprintf(stdout, "World-Provenance: store=%s episode=%s entries=%d-%d\n", h.DBPath, er.episode, er.from, er.to)
		printed++
	}
	if printed == 0 {
		what := "no coordinator entries"
		if *episode != "" {
			what = "no coordinator entries for episode " + *episode
		}
		fmt.Fprintf(stderr, "ailang-worldd provenance: %s in entries %d-%d\n", what, from, head)
		return exitUsage
	}
	return exitOK
}

// episodeRanges scans [from, to] forwards and groups the coordinator entries
// by episode, in order of first appearance.
func episodeRanges(r worldReader, from, to int64, max int) ([]episodeRange, error) {
	var out []episodeRange
	idx := map[string]int{}
	_, err := r.scanEntries(from, to, false, max, func(e wireEntry, rec *invRecord) (bool, error) {
		if rec == nil {
			return false, nil
		}
		i, seen := idx[rec.EpisodeID]
		if !seen {
			idx[rec.EpisodeID] = len(out)
			out = append(out, episodeRange{episode: rec.EpisodeID, from: e.Header.EntryIndex, to: e.Header.EntryIndex, num: 1})
			return false, nil
		}
		out[i].to = e.Header.EntryIndex
		out[i].num++
		return false, nil
	})
	return out, err
}
