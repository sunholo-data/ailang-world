package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 188 N1: GET /surface/model.json and GET /surface/stack.json, the
// surface's view model for the app (N2). The rules are surfacepage's; these
// handlers add only the daemon's live facts (build, live sessions) and the
// shared index read. Outside the frozen /v1 table, GET-only.

// surfaceReason is every non-model body: a positive empty state (200), a
// refusal (503), a malformed query (400) or a missing selection (404).
type surfaceReason struct {
	Schema  string `json:"schema"`
	Empty   bool   `json:"empty,omitempty"`
	Refused bool   `json:"refused,omitempty"`
	Reason  string `json:"reason"`
}

// writeSurfaceJSON writes v with the index route's headers, bounded by
// surface.MaxBytes like the index it is projected from.
func (d *Daemon) writeSurfaceJSON(w http.ResponseWriter, r *http.Request, status int, schema string, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		d.writeInternalErrorLog(r, err)
		status, buf = http.StatusInternalServerError, bytes.Buffer{}
		_ = json.NewEncoder(&buf).Encode(surfaceReason{Schema: schema, Reason: internalErrorMessage})
	} else if buf.Len() > surface.MaxBytes {
		status, buf = http.StatusServiceUnavailable, bytes.Buffer{}
		_ = json.NewEncoder(&buf).Encode(surfaceReason{Schema: schema, Refused: true,
			Reason: fmt.Sprintf("the surface view is over the 4 MiB (%d bytes) bound", surface.MaxBytes)})
	}
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

// surfaceIndexFor reads the index for a model route. ok is false when the
// response (empty or refused) has been written.
func (d *Daemon) surfaceIndexFor(w http.ResponseWriter, r *http.Request, schema string) (surface.Index, bool) {
	_, ix, empty, err := d.readSurfaceIndex()
	switch {
	case err != nil:
		d.writeSurfaceJSON(w, r, http.StatusServiceUnavailable, schema, surfaceReason{Schema: schema, Refused: true, Reason: d.surfaceIndexRefusal(err)})
		return ix, false
	case empty != "":
		d.writeSurfaceJSON(w, r, http.StatusOK, schema, surfaceReason{Schema: schema, Empty: true, Reason: empty})
		return ix, false
	}
	return ix, true
}

// handleSurfaceModel serves GET /surface/model.json; any query is a 400.
func (d *Daemon) handleSurfaceModel(w http.ResponseWriter, r *http.Request) {
	schema := surfacepage.ModelSchema
	if r.URL.RawQuery != "" {
		d.writeSurfaceJSON(w, r, http.StatusBadRequest, schema, surfaceReason{Schema: schema, Reason: unknownSurfaceKeyMessage})
		return
	}
	ix, ok := d.surfaceIndexFor(w, r, schema)
	if !ok {
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	build := CurrentBuild()
	path, _ := os.Executable()
	now := surfacepage.Now{At: time.Now().UTC(), Healthy: true, Revision: build.VCSRevision, Modified: build.VCSModified,
		Build: build.ExecutableSHA256, StartedAt: time.Unix(d.startedAt, 0).UTC(), Supervisor: serviceSupervisor(), InstalledSHA256: installedDigestAt(path)}
	now.Sessions, now.SessionsUnavailable = d.liveSessions(ctx, r, now.At)
	latest, err := d.reads.LogEntriesLatest(ctx, 1)
	if err != nil {
		now.LiveUnavailable = "the log could not be read"
	} else {
		now.Live = &surfacepage.LiveCursor{LastIndex: -1}
		if len(latest) > 0 {
			head := latest[0].EntryHash.String()
			now.Live.LastIndex, now.Live.LogHead = latest[0].Header.EntryIndex, &head
		}
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, schema, surfacepage.BuildModel(ix, now))
}

// handleSurfaceStack serves GET /surface/stack.json?brane=KEY&at=ID: M3's
// closed grammar, exactly one brane and one at.
func (d *Daemon) handleSurfaceStack(w http.ResponseWriter, r *http.Request) {
	schema := surfacepage.StackSchema
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1 {
		d.writeSurfaceJSON(w, r, http.StatusBadRequest, schema, surfaceReason{Schema: schema, Reason: unknownSurfaceKeyMessage})
		return
	}
	ix, ok := d.surfaceIndexFor(w, r, schema)
	if !ok {
		return
	}
	s, err := surfacepage.BuildStack(ix, q.Get("brane"), q.Get("at"), time.Now().UTC())
	if err != nil {
		d.writeSurfaceJSON(w, r, http.StatusNotFound, schema, surfaceReason{Schema: schema, Reason: err.Error()})
		return
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, schema, s)
}
