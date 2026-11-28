package daemon

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"html/template"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

type operatorStatus struct {
	V            int    `json:"v"`
	ID           string `json:"id"`
	Episode      string `json:"episode"`
	Kind         string `json:"kind"`
	State        string `json:"state"`
	Href         string `json:"href,omitempty"`
	StatusHref   string `json:"statusHref,omitempty"`
	ReportHref   string `json:"reportHref,omitempty"`
	ConfirmHref  string `json:"confirmHref,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Now          int64  `json:"now"`
	CreatedAt    int64  `json:"createdAt"`
	ExpiresAt    int64  `json:"expiresAt"`
	TerminalAt   int64  `json:"terminalAt,omitempty"`
	RunnerSeenAt int64  `json:"runnerSeenAt,omitempty"`
}

func operatorStatusFor(row store.OperatorActionRow, now int64) operatorStatus {
	out := operatorStatus{V: 1, ID: row.ID, Episode: row.Episode, Kind: row.Kind, State: row.State, Now: now,
		CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, TerminalAt: row.TerminalAt}
	switch row.State {
	case "requested", "confirmed", "claimed", "running", "awaiting-pr-confirmation", "succeeded", "failed", "expired", "cancelled", "uncertain":
	default:
		out.State, out.Reason = "unavailable", "unknown stored action state"
		return out
	}
	if !verbs.ValidOperatorActionID(row.ID) || !verbs.EpisodeIDPattern.MatchString(row.Episode) ||
		(row.Kind != "start" && row.Kind != "finish") {
		out.State, out.Reason = "unavailable", "invalid stored action identity"
		return out
	}
	out.Href = operatorPath + "/" + row.ID
	out.StatusHref = out.Href + "/status.json"
	if row.State == "succeeded" || row.State == "failed" || row.State == "uncertain" {
		out.ReportHref = "/v1/episodes/" + row.Episode + "/report"
	}

	return out
}

var operatorCard = template.Must(template.New("operator action").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Operator action</title><link rel="stylesheet" href="/surface/operator-card.css">
<script defer src="/surface/operator-card.js"></script></head><body>
<h1 id="operator-title" tabindex="-1">Operator action {{.Status.ID}}</h1>
<p id="operator-state" data-state="{{.Status.State}}">{{.Status.State}}</p>
<p id="operator-progress" role="status">Status at {{.Status.Now}}{{if .Status.Reason}} · {{.Status.Reason}}{{end}}</p>
{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
<section id="operator-controls" hidden aria-label="Confirm this action">
<button id="operator-key" type="button" hidden>Confirm with passkey</button>
<details><summary>Confirm with SSH signature</summary>
<p>Review the action and sign the exact message with an enrolled SSH key in the ceremony's namespace.</p>
{{if .Ceremony}}<pre id="operator-ceremony">{{.Ceremony}}</pre>{{end}}
<label for="operator-signature">Armored signature</label><textarea id="operator-signature" rows="4"></textarea>
<button id="operator-sign" type="button">Submit signature</button></details></section>
<p id="operator-answer" role="status" tabindex="-1"></p>
<button id="operator-refresh" type="button">Refresh status</button>
{{if .Status.StatusHref}}<p><a href="{{.Status.StatusHref}}">Action status (JSON)</a></p>{{end}}
<p><a id="operator-report" {{if .Status.ReportHref}}href="{{.Status.ReportHref}}"{{else}}hidden{{end}}>Episode report</a></p>
<p><a id="operator-back" href="/surface">Back to World</a></p>
</body></html>`))

func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id string, statusJSON bool) {
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	defer ro.Close()
	now := d.approveNow()
	row, err := ro.OperatorAction(ctx, id, now)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	status := operatorStatusFor(row, now)
	h, heartbeatErr := ro.OperatorHeartbeat(ctx)
	status.RunnerSeenAt = h.SeenAt
	if (heartbeatErr != nil || !h.Available(CurrentBuild().ExecutableSHA256, row.Kind, now)) && status.Reason == "" {
		status.Reason = "runner disconnected; showing durable action status"
	}
	if status.Href != "" && row.Kind == "start" && (row.State == "running" || row.State == "succeeded" || row.State == "uncertain" || row.State == "failed") {
		launch, err := ro.ReadOperatorLaunch(ctx, id)
		if err == nil {
			if launch.NoSpawn {
				status.State = "spawn-refused"
			} else {
				status.State = launch.ProcessStatus(agentrun.ProcessIdentity)
			}
			if status.ReportHref == "" {
				status.ReportHref = "/v1/episodes/" + row.Episode + "/report"
			}
			if status.State == "uncertain" {
				status.Reason = "launch fenced; inspect ailang-worldd service runner launch-status " + id
			}
		} else if !errors.Is(err, store.ErrOperatorMissing) {
			d.operatorError(w, err)
			return
		}
	}
	if statusJSON {
		writeApprovalJSON(w, http.StatusOK, status)
		return
	}
	card := struct {
		Status                 operatorStatus
		Review, Task, Ceremony string
	}{Status: status}
	if len(row.Action) != 0 && status.Href != "" {
		a, err := approval.DecodeOperatorAction(row.Action)
		digest := sha256.Sum256(row.Action)
		if err != nil || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || !bytes.Equal(digest[:], row.Digest) {
			d.operatorError(w, store.ErrOperatorConflict)
			return
		}
		card.Review, card.Task = a.Review, a.Task
		if row.State == "requested" && len(row.Challenge) == 32 && row.CreatedAt <= now {
			card.Ceremony, err = operatorCeremony(ctx, ro, row, status.Href+"/confirm")
			if err != nil {
				d.operatorError(w, err)
				return
			}
		}
	}
	setSurfaceHeaders(w)
	_ = operatorCard.Execute(w, card)
}

func operatorScript(w http.ResponseWriter, _ *http.Request, _ []string) {
	writeSurfaceAsset(w, "text/javascript; charset=utf-8", surfacepage.OperatorScript)
}

func operatorStyle(w http.ResponseWriter, _ *http.Request, _ []string) {
	writeSurfaceAsset(w, "text/css; charset=utf-8", surfacepage.OperatorStyle)
}
