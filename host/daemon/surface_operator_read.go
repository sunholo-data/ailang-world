package daemon

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"html/template"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
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
<html lang="en"><meta charset="utf-8"><title>Operator action</title>
<body><h1>{{.Status.ID}}</h1><p>{{.Status.State}}</p>{{if .Status.Reason}}<p>{{.Status.Reason}}</p>{{end}}
{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre>{{.Task}}</pre>{{end}}
{{if .Status.StatusHref}}<p><a href="{{.Status.StatusHref}}">Action status (JSON)</a></p>{{end}}
{{if .Status.ReportHref}}<p><a href="{{.Status.ReportHref}}">Episode report</a></p>{{end}}
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
			switch {
			case launch.NoSpawn:
				status.State = "spawn-refused"
			case launch.Phase == "reserved" || launch.Receipt.Status == "launched":
				status.State = "uncertain"
			default:
				status.State = launch.Receipt.Status
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
		Status       operatorStatus
		Review, Task string
	}{Status: status}
	if len(row.Action) != 0 && status.Href != "" {
		a, err := approval.DecodeOperatorAction(row.Action)
		digest := sha256.Sum256(row.Action)
		if err != nil || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || !bytes.Equal(digest[:], row.Digest) {
			d.operatorError(w, store.ErrOperatorConflict)
			return
		}
		card.Review, card.Task = a.Review, a.Task
	}
	setSurfaceHeaders(w)
	_ = operatorCard.Execute(w, card)
}
