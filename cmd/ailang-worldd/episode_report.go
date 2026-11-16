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
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const episodeReportHelp = `usage: ailang-worldd [--addr <url>] episode report <episode> [--json | --markdown] [--follow] [--interval 2s] [--addr <url>]

Summarises one episode from World's own store, in one request
(GET /v1/episodes/<episode>/report): total calls, calls by tool, effect
decisions by effect (allowed / failed / denied:<why>), budget per grant
(minted, used, left), plan refusals (refused before any effect) and handler
refusals (the effect was allowed and ran, but the handler answered ok:false),
the files read (successful ailang-read calls; --json lists them), the first
and last entries and their times.

  --markdown    render a scrubbed PR body (cannot combine with --json or --follow)
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
	flags := verbs.EpisodeReportFlags(fs) // row 243: the bar's admission parses the same set
	markdown, asJSON, follow, interval, addr := flags.Markdown, flags.JSON, flags.Follow, flags.Interval, flags.Addr
	positional, err := verbs.ParseInterleaved(fs, args)
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
	if *markdown && (*asJSON || *follow) {
		fmt.Fprintln(stderr, "refusing: --markdown cannot combine with --json or --follow")
		return exitUsage
	}
	base := globalAddr
	if *addr != "" {
		base = *addr
	}
	o := reportOpts{episode: positional[0], json: *asJSON, markdown: *markdown, interval: *interval}
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
	markdown bool
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
	if o.markdown {
		h, err := (worldReader{ctx: ctx, c: c}).health()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
		text, err := renderReportMarkdown(rep, h.DBPath, reportAdditions{})
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
		fmt.Fprint(stdout, text)
	} else if o.json {
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

// printReport renders the compact human table: the daemon's own text, which the app's bar shows too.
func printReport(w io.Writer, r daemon.EpisodeReport) { daemon.WriteEpisodeReportText(w, r) }
