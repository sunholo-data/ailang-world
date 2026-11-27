package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

// VerbSchema names GET /surface/verb.json's answer (row 243 V1, design §2.2).
const VerbSchema = "world/surface-verb/v1"

// verbAnswer is one answer of the verb route. Every answer carries v, schema, q, argv, verb, class
// and outcome; the rest belong to their outcome.
type verbAnswer struct {
	V                   int      `json:"v"`
	Schema              string   `json:"schema"`
	Q                   string   `json:"q"`
	Argv                []string `json:"argv"`
	Verb                string   `json:"verb"`
	Class               string   `json:"class"`
	Outcome             string   `json:"outcome"`
	Reason              string   `json:"reason,omitempty"`
	CLI                 string   `json:"cli,omitempty"`
	Href                string   `json:"href,omitempty"`
	Text                string   `json:"text,omitempty"`
	JSON                string   `json:"json,omitempty"`
	Stderr              string   `json:"stderr,omitempty"`
	Exit                int      `json:"exit,omitempty"`
	CatalogueGeneration int64    `json:"catalogueGeneration,omitempty"`
}

// handleSurfaceVerb is GET /surface/verb.json?q=<line>, behind approvalRoute: the bar's line through
// host/verbs' one grammar. It reads (episode report, the /v1 route's own builder; why, host/why's walk) and writes nothing:
// no store write or process. An operator's typo is a 200 refusal.
func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []string) {
	query := r.URL.Query()
	if !query.Has("q") {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "q is required")
		return
	}
	a := verbAnswer{V: 1, Schema: VerbSchema, Q: query.Get("q"), Argv: []string{}, Class: string(verbs.Refused), Outcome: "refused"}
	defer func() { capVerbAnswer(&a); writeApprovalJSON(w, http.StatusOK, a) }()
	argv, err := verbs.Split(a.Q) // its first check is the 1,024-byte bound: no oversized line is tokenized
	if err != nil {
		a.Reason = err.Error()
		return
	}
	v, args := verbs.Lookup(argv)
	a.Argv, a.Verb, a.Class, a.CLI, a.Reason = argv, v.Words, string(v.Class), v.CLI, v.Reason
	if d.operatorVerb(r, argv, &a) {
		return
	}
	switch {
	case v.Class == verbs.Address:
		a.Outcome = "address"
	case v.Words == "episode report":
		d.verbEpisodeReport(r, args, &a)
	case v.Words == "why":
		d.verbWhy(r, args, &a)
	case v.Words == "profile list":
		d.verbProfileList(args, &a)
	case v.Words == "preset list":
		d.verbPresetList(args, &a)
	case v.Class == verbs.Navigate:
		if a.Href, a.Reason = verbs.ApproveTarget(args); a.Href != "" {
			a.Outcome = "navigate"
		}
	case v.Class == verbs.Read || v.Class == verbs.Mutate || v.Class == verbs.Host:
		a.Outcome = "terminal"
	}
}

// verbEpisodeReport runs episode report for the bar: the CLI's flag set and parse, the bar's
// admission, then the /v1 route's builder; the text is the CLI's own printer's.
func (d *Daemon) verbEpisodeReport(r *http.Request, args []string, a *verbAnswer) {
	episode, asJSON, reason := verbs.AdmitEpisodeReport(args)
	if reason != "" {
		a.Reason = reason
		return
	}
	if len(episode) > maxReportEpisodeBytes {
		a.Reason = fmt.Sprintf("episode must be 1..%d bytes", maxReportEpisodeBytes)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	rep, found, err := d.episodeReport(ctx, episode)
	if err != nil {
		d.writeInternalErrorLog(r, err)
		a.Reason = internalErrorMessage
		if reportTimedOut(ctx, err) {
			a.Reason = fmt.Sprintf("read deadline (%s) exceeded", d.readDeadline)
		}
		return
	}
	a.Outcome = "result"
	switch {
	case !found:
		a.Exit, a.Stderr = 1, "ailang-worldd episode report: episode "+episode+" has no committed calls\n"
	case asJSON:
		var b bytes.Buffer
		_ = json.NewEncoder(&b).Encode(rep) // the /v1 route's bytes: compact JSON and one newline, as the CLI prints them
		a.JSON = b.String()
	default:
		var b strings.Builder
		WriteEpisodeReportText(&b, rep)
		a.Text = b.String()
	}
}
