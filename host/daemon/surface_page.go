package daemon

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 188 M2: GET /surface, the world timeline, server-rendered from the
// surface index. It is outside the frozen /v1 table, GET-only, under the
// workbench's CSP; it carries no script until M4 adds /surface/surface.js.

const (
	unknownSurfaceKeyMessage = "unsupported surface query parameter"
	liveSessionsUnavailable  = "live sessions could not be read"
)

// handleSurface accepts either no query or exactly one brane and one at.
func (d *Daemon) handleSurface(w http.ResponseWriter, r *http.Request) {
	q, queryErr := url.ParseQuery(r.URL.RawQuery)
	stackQuery := r.URL.RawQuery != ""
	if queryErr != nil || stackQuery && (len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1) {
		writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", unknownSurfaceKeyMessage)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	build := CurrentBuild()
	now := surfacepage.Now{At: time.Now(), Revision: build.VCSRevision, Modified: build.VCSModified}
	now.Sessions, now.SessionsUnavailable = d.liveSessions(ctx, r, now.At)

	status := http.StatusOK
	var tl surfacepage.Timeline
	_, ix, empty, err := d.readSurfaceIndex()
	if stackQuery {
		s := surfacepage.Stack{}
		switch {
		case err != nil:
			status = http.StatusServiceUnavailable
			s.Reason = d.surfaceIndexRefusal(err)
		case empty != "":
			s.Reason = empty
		default:
			s, err = surfacepage.LayoutStack(ix, q.Get("brane"), q.Get("at"), now.At)
			if err != nil {
				status = http.StatusNotFound
				s.Reason = err.Error()
			}
		}
		body, renderErr := surfacepage.RenderStack(s)
		if renderErr != nil {
			d.writeWorkbenchInternalError(w, r, renderErr)
			return
		}
		setWorkbenchHeaders(w)
		w.WriteHeader(status)
		_, _ = w.Write(body)
		return
	}
	switch {
	case err != nil:
		status = http.StatusServiceUnavailable
		tl = surfacepage.Timeline{Reason: d.surfaceIndexRefusal(err)}
	case empty != "":
		tl = surfacepage.Timeline{Reason: empty}
	default:
		// Healthy claims only this: the daemon is serving and the index
		// validated.
		now.Healthy = true
		tl = surfacepage.LayoutTimeline(ix, now.At)
	}
	body, err := surfacepage.RenderTimeline(tl, now)
	if err != nil {
		d.writeWorkbenchInternalError(w, r, err)
		return
	}
	setWorkbenchHeaders(w)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// liveSessions reads the credentials store read-only, as `session list`
// does, and returns the sessions live at now, oldest first. A failure is
// logged and named on the page, never fatal to it.
func (d *Daemon) liveSessions(ctx context.Context, r *http.Request, now time.Time) ([]surfacepage.Session, string) {
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.writeInternalErrorLog(r, err)
		return nil, liveSessionsUnavailable
	}
	defer func() { _ = ro.Close() }()
	rows, err := ro.ListSessions(ctx, "", store.MaxListSessions)
	if err != nil {
		d.writeInternalErrorLog(r, err)
		return nil, liveSessionsUnavailable
	}
	if len(rows) >= store.MaxListSessions {
		// One page is the oldest credentials; a live one may lie past it.
		return nil, "over 500 session credentials: live sessions not counted"
	}
	var out []surfacepage.Session
	for _, row := range rows {
		if now.Unix() < row.ExpiresAt {
			out = append(out, surfacepage.Session{Episode: row.EpisodeID, Started: time.Unix(row.CreatedAt, 0)})
		}
	}
	return out, ""
}
