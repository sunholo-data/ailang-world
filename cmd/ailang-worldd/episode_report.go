package main

// Row 198 (w-episode-report): `episode report` prints the daemon's one-pass
// summary of an episode (GET /v1/episodes/{episode}/report): calls by tool,
// effect decisions, budget per grant, plan refusals and handler ok:false
// results. --json prints the route's JSON unchanged; --follow re-prints it
// whenever the store's head moves, as `log tail --follow` streams entries.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

const episodeReportHelp = `usage: ailang-worldd [--addr <url>] episode report <episode> [--json] [--follow] [--interval 2s] [--addr <url>]

Summarises one episode from World's own store, in one request
(GET /v1/episodes/<episode>/report): total calls, calls by tool, effect
decisions by effect (allowed / failed / denied:<why>), budget per grant
(minted, used, left), plan refusals (refused before any effect) and handler
refusals (the effect was allowed and ran, but the handler answered ok:false),
the files read (successful ailang-read calls; --json lists them), the first
and last entries and their times.

  --json        print the route's JSON (one object; with --follow one compact
                line per change)
  --follow      re-print when the episode changes, until Ctrl-C; the store
                head is polled every --interval and the report is fetched
                only when it moved. A daemon restart is retried with backoff
  --interval D  poll interval with --follow (default 2s)
  --addr <url>  the daemon (default: the global --addr, ` + daemon.DefaultAddr + `);
                a bare host:port is accepted

Exit 1 when the episode has no committed call (without --follow).
`

func runEpisodeReport(globalAddr string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd episode report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	asJSON := fs.Bool("json", false, "print the route's JSON")
	follow := fs.Bool("follow", false, "re-print when the episode changes, until Ctrl-C")
	interval := fs.Duration("interval", 2*time.Second, "poll interval with --follow")
	addr := fs.String("addr", "", "the daemon (default: the global --addr)")
	positional, err := parseInterleaved(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, episodeReportHelp)
			return exitOK
		}
		fmt.Fprint(stderr, episodeReportHelp)
		return exitUsage
	}
	if len(positional) != 1 || positional[0] == "" {
		fmt.Fprint(stderr, episodeReportHelp)
		return exitUsage
	}
	if *interval <= 0 {
		fmt.Fprintln(stderr, "ailang-worldd episode report: --interval must be positive")
		return exitUsage
	}
	base := globalAddr
	if *addr != "" {
		base = *addr
	}
	o := reportOpts{episode: positional[0], json: *asJSON, interval: *interval}
	if !*follow {
		ctx, cancel := budgetContext(whyBudget)
		defer cancel()
		return reportOnce(ctx, newClient(base), o, stdout, stderr)
	}
	ctx, stop := serveSignalContext()
	defer stop()
	return reportFollow(ctx, newClient(base), o, stdout, stderr)
}

type reportOpts struct {
	episode  string
	json     bool
	interval time.Duration
}

// errNoCalls is the route's 404: the episode has no committed call yet.
var errNoCalls = errors.New("no committed calls")

// fetchReport GETs the report: the body as sent, and its decoded form.
func fetchReport(ctx context.Context, c *client, episode string) ([]byte, daemon.EpisodeReport, error) {
	var rep daemon.EpisodeReport
	path := "/v1/episodes/" + url.PathEscape(episode) + "/report"
	status, body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, rep, err
	}
	switch status {
	case http.StatusOK:
		if err := json.Unmarshal(body, &rep); err != nil {
			return nil, rep, fmt.Errorf("GET %s: malformed body: %v", path, err)
		}
		return bytes.TrimSpace(body), rep, nil
	case http.StatusNotFound:
		var apiErr daemon.APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Class == "NotFound" {
			return nil, rep, errNoCalls
		}
	}
	return nil, rep, fmt.Errorf("GET %s returned HTTP %d: %s", path, status, strings.TrimSpace(string(body)))
}

func reportOnce(ctx context.Context, c *client, o reportOpts, stdout, stderr io.Writer) int {
	body, rep, err := fetchReport(ctx, c, o.episode)
	if errors.Is(err, errNoCalls) {
		fmt.Fprintf(stderr, "ailang-worldd episode report: episode %s has no committed calls\n", o.episode)
		return exitUsage
	}
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd episode report: %v\n", err)
		return exitUsage
	}
	if o.json {
		fmt.Fprintln(stdout, string(body))
	} else {
		printReport(stdout, rep)
	}
	return exitOK
}

// reportFollow polls the head and re-prints the report when it changed.
func reportFollow(ctx context.Context, c *client, o reportOpts, stdout, stderr io.Writer) int {
	lastHead, lastBody := "", ""
	waiting := false
	backoff := 250 * time.Millisecond
	var down time.Time
	for {
		head, _, err := headState(ctx, c)
		if err == nil && head != lastHead {
			var body []byte
			var rep daemon.EpisodeReport
			body, rep, err = fetchReport(ctx, c, o.episode)
			switch {
			case errors.Is(err, errNoCalls):
				err = nil
				lastHead = head
				if !waiting {
					fmt.Fprintf(stderr, "ailang-worldd episode report: waiting for episode %s's first committed call\n", o.episode)
					waiting = true
				}
			case err == nil:
				lastHead = head
				if string(body) != lastBody {
					lastBody = string(body)
					printFollowed(stdout, o, body, rep)
				}
			}
		}
		if err != nil {
			if ctx.Err() != nil {
				return exitOK
			}
			if down.IsZero() {
				down = time.Now()
				fmt.Fprintf(stderr, "ailang-worldd episode report: daemon unreachable (%v); retrying\n", err)
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
			fmt.Fprintf(stderr, "ailang-worldd episode report: reconnected after %s of outage\n", time.Since(down).Round(time.Millisecond))
			down, backoff = time.Time{}, 250*time.Millisecond
		}
		if !sleepCtx(ctx, o.interval) {
			return exitOK
		}
	}
}

func printFollowed(w io.Writer, o reportOpts, body []byte, rep daemon.EpisodeReport) {
	if o.json {
		var compact bytes.Buffer
		if json.Compact(&compact, body) == nil {
			body = compact.Bytes()
		}
		fmt.Fprintln(w, string(body))
		return
	}
	fmt.Fprintf(w, "--- %s\n", time.Now().UTC().Format(time.RFC3339))
	printReport(w, rep)
}

// printReport renders the compact human table.
func printReport(w io.Writer, r daemon.EpisodeReport) {
	fmt.Fprintf(w, "episode %s: %d calls, entries %d-%d, %s .. %s\n",
		r.Episode, r.TotalCalls, r.Entries[0], r.Entries[1], r.FirstAt, r.LastAt)
	fmt.Fprintln(w, "calls by tool:")
	for _, kv := range sortedCounts(r.CallsByTool) {
		fmt.Fprintf(w, "  %-18s %6d\n", kv.name, kv.n)
	}
	fmt.Fprintln(w, readsLine(r.Reads))
	fmt.Fprintln(w, "effects:")
	for _, e := range r.Effects {
		fmt.Fprintf(w, "  %-22s %-16s %6d\n", e.Effect, e.Decision, e.Count)
	}
	fmt.Fprintln(w, "budgets (minted / used / left):")
	for _, b := range r.Budgets {
		fmt.Fprintf(w, "  %-22s %6d %6d %6d\n", b.Effect+"@"+b.Scope, b.Minted, b.Used, b.Left)
	}
	printRefusals(w, "plan refusals (refused before any effect)", r.PlanRefusals)
	printRefusals(w, "handler refusals (effect allowed, handler ok:false)", r.HandlerRefusals)
	if r.Uncommitted > 0 {
		fmt.Fprintf(w, "uncommitted intents: %d\n", r.Uncommitted)
	}
}

// readsLine is the one-line reads summary; --json carries the file list.
func readsLine(r daemon.ReportReads) string {
	line := fmt.Sprintf("reads: %d calls, %d files", r.Count, len(r.Files))
	if r.Truncated > 0 {
		line += fmt.Sprintf(", %d more not listed", r.Truncated)
	}
	if r.Unlisted > 0 {
		line += fmt.Sprintf(", %d outside the worktree", r.Unlisted)
	}
	return line
}

func printRefusals(w io.Writer, title string, r daemon.ReportRefusals) {
	parts := make([]string, 0, len(r.ByTool))
	for _, kv := range sortedCounts(r.ByTool) {
		parts = append(parts, fmt.Sprintf("%s %d", kv.name, kv.n))
	}
	line := fmt.Sprintf("%s: %d", title, r.Count)
	if len(parts) > 0 {
		line += " (" + strings.Join(parts, ", ") + ")"
	}
	fmt.Fprintln(w, line)
	for _, f := range r.First {
		fmt.Fprintf(w, "  #%d %s: %s\n", f.Entry, f.Tool, f.Reason)
	}
}

type namedCount struct {
	name string
	n    int
}

// sortedCounts orders a count map by count, then name.
func sortedCounts(m map[string]int) []namedCount {
	out := make([]namedCount, 0, len(m))
	for k, n := range m {
		out = append(out, namedCount{k, n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].n != out[j].n {
			return out[i].n > out[j].n
		}
		return out[i].name < out[j].name
	})
	return out
}
