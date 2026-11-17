package daemon

// Row 243 V2b (design_docs/planned/w-world-address-bar-verbs-v2.md §2.4, §2.5): why in the bar. The
// walk is host/why's, the CLI's own; it reads World through a private mux over the seven /v1 GET
// handlers it needs, in-process. No route, POST, store write or process is added.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/why"
)

// MaxVerbAnswerBytes caps a result answer's text, json and stderr together (design A2): 1 MiB of raw
// bytes. The JSON-encoded body can be larger (a control byte or <, >, & escapes to six bytes, \u00XX)
// and is still bounded, by six times the cap plus the answer's other members.
const MaxVerbAnswerBytes = 1 << 20

// verbAnswerCap is the cap the route applies; a test lowers it to drive the refusal end to end.
var verbAnswerCap = MaxVerbAnswerBytes

// capVerbAnswer turns a result over the cap into a refusal with its size and the CLI line (design §2.5 e).
func capVerbAnswer(a *verbAnswer) {
	if n := len(a.Text) + len(a.JSON) + len(a.Stderr); a.Outcome == "result" && n > verbAnswerCap {
		a.Outcome, a.Text, a.JSON, a.Stderr, a.Exit = "refused", "", "", "", 0
		a.Reason = fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: %s", n, a.CLI)
	}
}

// whyRoutes is the seven /v1 GET patterns why's walk reads, each to the method value Handler registers.
func (d *Daemon) whyRoutes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /v1/health":        d.handleHealth,
		"GET /v1/head":          d.handleHead,
		"GET /v1/worlds/{ref}":  d.handleWorld,
		"GET /v1/objects/{ref}": d.handleObject,
		"GET /v1/log/{index}":   d.handleLogEntry,
		"GET /v1/log":           d.handleLogRange,
		"GET /v1/receipts/{id}": d.handleReceipt,
	}
}

// whyMux is a private mux over exactly whyRoutes, built per walk.
func (d *Daemon) whyMux() *http.ServeMux {
	m := http.NewServeMux()
	for pattern, h := range d.whyRoutes() {
		m.HandleFunc(pattern, h)
	}
	return m
}

// whyTally is what the walk's reads saw: a read deadline (503) or another failure.
type whyTally struct{ timeout, failed bool }

// boundedRecorder is one in-process response, its body capped at max bytes.
type boundedRecorder struct {
	header http.Header
	status int
	body   bytes.Buffer
	max    int
	over   bool
}

func (b *boundedRecorder) Header() http.Header { return b.header }
func (b *boundedRecorder) WriteHeader(status int) {
	if b.status == 0 {
		b.status = status
	}
}
func (b *boundedRecorder) Write(p []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	if b.body.Len()+len(p) > b.max {
		b.over = true
		return 0, errors.New("response over the in-process cap")
	}
	return b.body.Write(p)
}

// whyGetter is why's Getter over mux, in-process: a 404 is an answer; an unrouted path, a request that
// cannot be built or a body over why.MaxResponseBytes is an error, logged for the operator. tally
// records each 503 (the handlers' read deadline) and each other failure.
func (d *Daemon) whyGetter(r *http.Request, mux *http.ServeMux, tally *whyTally) why.Getter {
	return func(ctx context.Context, path string) (int, []byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err == nil {
			if _, pattern := mux.Handler(req); pattern == "" {
				err = fmt.Errorf("unrouted in-process read: %s", path)
			}
		}
		var rec *boundedRecorder
		if err == nil {
			rec = &boundedRecorder{header: http.Header{}, max: why.MaxResponseBytes}
			mux.ServeHTTP(rec, req)
			if rec.over {
				err = fmt.Errorf("response from %s exceeds %d bytes", path, why.MaxResponseBytes)
			}
		}
		if err != nil {
			tally.failed = true
			d.writeInternalErrorLog(r, err)
			return 0, nil, err
		}
		switch {
		case rec.status == http.StatusServiceUnavailable:
			tally.timeout = true
		case rec.status >= 500:
			tally.failed = true
		}
		return rec.status, rec.body.Bytes(), nil
	}
}

// verbWhy runs why for the bar (design §2.5): the CLI's parse and the bar's admission, then the walk
// under the daemon's read deadline, then the outcome by precedence.
func (d *Daemon) verbWhy(r *http.Request, args []string, a *verbAnswer) {
	target, scan, asJSON, reason := verbs.AdmitWhy(args)
	if reason != "" {
		a.Reason = reason
		return
	}
	t, err := why.ParseTarget(target)
	if err != nil {
		a.Reason = "ailang-worldd why: " + err.Error()
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	var tally whyTally
	ans := why.Run(why.Reader{Ctx: ctx, Get: d.whyGetter(r, d.whyMux(), &tally)}, t, why.Options{Scan: scan, JSON: asJSON})
	switch {
	case r.Context().Err() != nil: // (a) the page went away: nobody reads this
		a.Reason = "the request ended before the walk did"
	case errors.Is(ctx.Err(), context.DeadlineExceeded) || tally.timeout: // (b)
		a.Reason = fmt.Sprintf("read deadline (%s) exceeded: the walk did not finish in the app; at a terminal it has 60 s: %s", d.readDeadline, a.CLI)
	case tally.failed: // (c)
		a.Reason = "the daemon could not read World during the walk, so no chain is shown; at a terminal: " + a.CLI
	default: // (d), (f)
		whyResult(a, ans, asJSON)
	}
}

// whyResult is the walk's answer as the bar's result (design §2.5 f), or (d) a refusal when it is not
// exactly one of the CLI's three shapes: stdout with exit 0 or 3, stderr with exit 1.
func whyResult(a *verbAnswer, ans why.Answer, asJSON bool) {
	out, errOut := ans.Stdout != "", ans.Stderr != ""
	if out == errOut || (out && ans.Exit != 0 && ans.Exit != 3) || (errOut && ans.Exit != 1) {
		a.Reason = "the walk answered something the app cannot show; at a terminal: " + a.CLI
		return
	}
	a.Outcome, a.Exit = "result", ans.Exit
	switch {
	case errOut:
		a.Stderr = ans.Stderr
	case asJSON:
		a.JSON = ans.Stdout
	default:
		a.Text = ans.Stdout
	}
}
