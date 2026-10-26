package daemon

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 188 M2: the world timeline, server-rendered from the surface index;
// since N2b at GET /surface/static, the no-script fallback, while GET /surface
// is the app shell (N2b). Both are outside the frozen /v1 table, GET-only,
// under the surface's CSP. Their shared script is GET /surface/surface.js
// (M4a); the app's own are GET /surface/app.js, app.css and ailang-logo.svg.

const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"

func setSurfaceHeaders(w http.ResponseWriter) {
	setWorkbenchHeaders(w)
	w.Header().Set("Content-Security-Policy", surfaceCSP)
}

func writeSurfaceError(w http.ResponseWriter, status int, class, message string) {
	setSurfaceHeaders(w)
	w.WriteHeader(status)
	_ = workbenchErrorTemplate.Execute(w, struct {
		Class   string
		Message string
	}{class, message})
}

const (
	unknownSurfaceKeyMessage = "unsupported surface query parameter"
	liveSessionsUnavailable  = "live sessions could not be read"
)

// surfaceQuery is both pages' grammar: no query, or exactly one brane and
// one at (stack). ok is false, and an inert 400 written, for anything else.
func surfaceQuery(w http.ResponseWriter, r *http.Request) (q url.Values, stack, ok bool) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	stack = r.URL.RawQuery != ""
	if err != nil || stack && (len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1) {
		writeSurfaceError(w, http.StatusBadRequest, "BadRequest", unknownSurfaceKeyMessage)
		return nil, false, false
	}
	return q, stack, true
}

// handleSurface serves the app shell, the same bytes for every accepted
// query: the app reads the model and opens a deep link itself, and shows an
// empty or refused index's reason from /surface/model.json. The status keeps
// M2/M3's meaning for a client without script: 503 for a refused index, 404
// for a deep link to no state, else 200.
func (d *Daemon) handleSurface(w http.ResponseWriter, r *http.Request) {
	q, stack, ok := surfaceQuery(w, r)
	if !ok {
		return
	}
	status := http.StatusOK
	_, ix, empty, err := d.readSurfaceIndex()
	switch {
	case err != nil:
		status = http.StatusServiceUnavailable
	case empty != "":
	case stack:
		if _, err := surfacepage.LayoutStack(ix, q.Get("brane"), q.Get("at"), time.Now()); err != nil {
			status = http.StatusNotFound
		}
	}
	setSurfaceHeaders(w)
	w.WriteHeader(status)
	_, _ = w.Write(surfacepage.AppShell)
}

// handleSurfaceStatic serves the server-rendered timeline, or with brane and
// at its stack: M2 and M3's pages, unchanged.
func (d *Daemon) handleSurfaceStatic(w http.ResponseWriter, r *http.Request) {
	q, stackQuery, ok := surfaceQuery(w, r)
	if !ok {
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
		setSurfaceHeaders(w)
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
	setSurfaceHeaders(w)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// handleSurfaceScript serves the surface's one embedded script, with the
// headers of /workbench/live.js.
func (d *Daemon) handleSurfaceScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	_, _ = w.Write(surfacepage.Script)
}

// The app's assets: the exact embedded bytes, revalidated on every use
// (no-cache), and inert if opened as documents.
func (d *Daemon) handleSurfaceAppScript(w http.ResponseWriter, _ *http.Request) {
	writeSurfaceAsset(w, "text/javascript; charset=utf-8", surfacepage.AppScript)
}

func (d *Daemon) handleSurfaceAppStyle(w http.ResponseWriter, _ *http.Request) {
	writeSurfaceAsset(w, "text/css; charset=utf-8", surfacepage.AppStyle)
}

func (d *Daemon) handleSurfaceLogo(w http.ResponseWriter, _ *http.Request) {
	writeSurfaceAsset(w, "image/svg+xml", surfacepage.AppLogo)
}

func writeSurfaceAsset(w http.ResponseWriter, contentType string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
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
