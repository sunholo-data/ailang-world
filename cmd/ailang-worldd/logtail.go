package main

// Row 138 M2 (design_docs/planned/w-worldd-developer-cli.md §3.4): `log
// tail` prints log entries one line each, and with --follow polls for new
// ones. The cursor is the next entry index and advances by the number of
// items RECEIVED, never by the page size (MUT-FOLLOW-SKIP): a short page
// followed by new commits must print every one of them exactly once.

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/coordinator"
)

const logTailHelp = `usage: ailang-worldd [--addr <url>] log tail [--from N] [--follow] [--interval 1s] [--raw]

Prints log entries, one line each: index, short entry hash, writtenBy, and for
a coordinator invocation its episode, skill and effect statuses.

  --from N      first entry (default: head-19, the last 20 entries)
  --follow      keep polling for new entries until Ctrl-C; a daemon restart
                is retried with backoff (up to 5 s) and the outage is named
  --interval D  poll interval with --follow (default 1s)
  --raw         print each entry as its JSON
`

const tailBackoffMax = 5 * time.Second

type tailOpts struct {
	from     int64
	follow   bool
	interval time.Duration
	raw      bool
}

func runLogTail(addr string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd log tail", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	from := fs.Int64("from", -1, "first entry (default head-19)")
	follow := fs.Bool("follow", false, "keep polling for new entries")
	interval := fs.Duration("interval", time.Second, "poll interval with --follow")
	raw := fs.Bool("raw", false, "print each entry as JSON")
	if code, done := parseVerbFlags(fs, args, logTailHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd log tail: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	if *interval <= 0 {
		fmt.Fprintln(stderr, "ailang-worldd log tail: --interval must be positive")
		return exitUsage
	}
	ctx, stop := serveSignalContext()
	defer stop()
	return logTail(ctx, newClient(addr), tailOpts{from: *from, follow: *follow, interval: *interval, raw: *raw}, stdout, stderr)
}

// sleepCtx waits d or until ctx ends; false means ctx ended.
func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func logTail(ctx context.Context, c *client, o tailOpts, stdout, stderr io.Writer) int {
	r := worldReader{ctx: ctx, c: c}
	next := o.from
	if next < 0 {
		h, ok, err := r.headIndex()
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd log tail: %v\n", err)
			return exitUsage
		}
		next = 0
		if ok && h > 19 {
			next = h - 19
		}
	}
	backoff := 250 * time.Millisecond
	var down time.Time
	for {
		items, err := r.page(next, logPageMax)
		if err != nil {
			if ctx.Err() != nil {
				return exitOK
			}
			if !o.follow {
				fmt.Fprintf(stderr, "ailang-worldd log tail: %v\n", err)
				return exitUsage
			}
			if down.IsZero() {
				down = time.Now()
				fmt.Fprintf(stderr, "ailang-worldd log tail: daemon unreachable (%v); retrying\n", err)
			}
			if !sleepCtx(ctx, backoff) {
				return exitOK
			}
			if backoff *= 2; backoff > tailBackoffMax {
				backoff = tailBackoffMax
			}
			continue
		}
		if !down.IsZero() {
			fmt.Fprintf(stderr, "ailang-worldd log tail: reconnected after %s of outage; resuming at entry %d (nothing skipped)\n",
				time.Since(down).Round(time.Millisecond), next)
			down, backoff = time.Time{}, 250*time.Millisecond
		}
		for _, e := range items {
			fmt.Fprintln(stdout, tailLine(r, e, o.raw))
		}
		// The cursor advances by what was RECEIVED, never by the page size.
		next += int64(len(items))
		if len(items) == logPageMax {
			continue
		}
		if !o.follow {
			return exitOK
		}
		if !sleepCtx(ctx, o.interval) {
			return exitOK
		}
	}
}

// tailLine renders one entry.
func tailLine(r worldReader, e wireEntry, raw bool) string {
	if raw {
		b, _ := json.Marshal(e)
		return string(b)
	}
	short := strings.TrimPrefix(e.EntryHash, "sha256:")
	if len(short) > 8 {
		short = short[:8]
	}
	line := fmt.Sprintf("#%d %s %s", e.Header.EntryIndex, short, e.Header.WrittenBy)
	rec, _, err := r.record(e)
	if err != nil || rec == nil {
		return line
	}
	line += fmt.Sprintf(" %s %s", rec.EpisodeID, rec.SkillID)
	if rec.SemanticID != coordinator.RecordV2 {
		return line
	}
	statuses := make([]string, 0, len(rec.Effects))
	for _, ref := range rec.Effects {
		d, ok := r.effectRecord(ref)
		switch {
		case !ok:
			statuses = append(statuses, "?")
		case !d.Allowed:
			statuses = append(statuses, d.Effect+" denied:"+d.Denial)
		case d.Failed:
			statuses = append(statuses, d.Effect+" failed")
		default:
			statuses = append(statuses, d.Effect+" ok")
		}
	}
	if len(statuses) == 0 {
		return line + " [no effects]"
	}
	return line + " [" + strings.Join(statuses, ", ") + "]"
}
