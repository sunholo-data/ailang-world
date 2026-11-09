package daemon

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 220 M4: in follow mode every surface route reads one published
// snapshot pointer per request and serves the repo-qualified v2 documents.
// The 503 refusal is only for "no repo has an accepted graph"; a partial
// surface is a 200 whose repos carry their own degraded status.

// followSnapshot is the one read of the published snapshot for a request.
func (d *Daemon) followSnapshot() *surfacefollow.Snapshot {
	if d.follow.snapshot != nil {
		return d.follow.snapshot()
	}
	return d.follow.publisher.Snapshot()
}

// followSources pairs the snapshot's statuses with the config's branches,
// in config order. Neither carries a path.
func (d *Daemon) followSources(s *surfacefollow.Snapshot) []surfacepage.FollowSource {
	byID := map[string]surfacefollow.RepoStatus{}
	for _, r := range s.Status.Repos {
		byID[r.ID] = r
	}
	out := make([]surfacepage.FollowSource, 0, len(d.follow.config.Repos))
	for _, r := range d.follow.config.Repos {
		st, ok := byID[r.ID]
		if !ok {
			st = surfacefollow.RepoStatus{ID: r.ID, State: surfacepage.SourceBuilding}
		}
		out = append(out, surfacepage.FollowSource{ID: r.ID, Branch: r.Branch, State: st.State,
			CheckedAt: st.CheckedAt, BuiltAt: st.BuiltAt, Failure: st.LastFailure})
	}
	return out
}

// followRefusal is the structured 503 body when no repo has an accepted
// graph: the reason names every repo's state and the repos carry the lines.
func (d *Daemon) followRefusal(s *surfacefollow.Snapshot, schema string, now time.Time) surfaceReason {
	var states []string
	var repos []surfacepage.ModelRepo
	for _, src := range d.followSources(s) {
		line, warn := surfacepage.SourceLine(src, now)
		state := src.State
		if src.Failure != nil {
			state += ": " + src.Failure.Class
		}
		states = append(states, src.ID+": "+state)
		repos = append(repos, surfacepage.ModelRepo{ID: src.ID, Branch: src.Branch, State: src.State, Line: line, Warn: warn})
	}
	return surfaceReason{Schema: schema, Refused: true, Repos: repos,
		Reason: "no repo has an accepted graph yet: " + strings.Join(states, "; ")}
}

func (d *Daemon) followNow(r *http.Request) surfacepage.Now {
	ctx, cancel := d.readCtx(r)
	defer cancel()
	build := CurrentBuild()
	path, _ := os.Executable()
	now := surfacepage.Now{At: time.Now().UTC(), Healthy: true, Revision: build.VCSRevision, Modified: build.VCSModified,
		Build: build.ExecutableSHA256, StartedAt: time.Unix(d.startedAt, 0).UTC(), Supervisor: serviceSupervisor(), InstalledSHA256: installedDigestAt(path)}
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
	return now
}

// handleFollowModel serves GET /surface/model.json in follow mode.
func (d *Daemon) handleFollowModel(w http.ResponseWriter, r *http.Request) {
	schema := surfacepage.FollowModelSchema
	s := d.followSnapshot()
	if len(s.Envelope.Repos) == 0 {
		d.writeSurfaceJSON(w, r, http.StatusServiceUnavailable, schema, d.followRefusal(s, schema, time.Now().UTC()))
		return
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, schema, surfacepage.BuildFollowModel(s.Envelope, d.followSources(s), d.followNow(r)))
}

// handleFollowStack serves GET /surface/stack.json in follow mode, from the
// same kind of snapshot read, carrying its generation.
func (d *Daemon) handleFollowStack(w http.ResponseWriter, r *http.Request, q url.Values) {
	schema := surfacepage.FollowStackSchema
	s := d.followSnapshot()
	if len(s.Envelope.Repos) == 0 {
		d.writeSurfaceJSON(w, r, http.StatusServiceUnavailable, schema, d.followRefusal(s, schema, time.Now().UTC()))
		return
	}
	st, err := surfacepage.BuildFollowStack(s.Envelope, q.Get("brane"), q.Get("at"), time.Now().UTC())
	if err != nil {
		d.writeSurfaceJSON(w, r, http.StatusNotFound, schema, surfaceReason{Schema: schema, Reason: err.Error()})
		return
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, schema, st)
}

// followShellStatus is the app shell's status for a client without script:
// 503 with no accepted graph, 404 for a selection that does not resolve.
func (d *Daemon) followShellStatus(q url.Values, stack bool) int {
	s := d.followSnapshot()
	if len(s.Envelope.Repos) == 0 {
		return http.StatusServiceUnavailable
	}
	if stack {
		if _, err := surfacepage.BuildFollowStack(s.Envelope, q.Get("brane"), q.Get("at"), time.Now()); err != nil {
			return http.StatusNotFound
		}
	}
	return http.StatusOK
}

// handleFollowStatic serves /surface/static in follow mode: one timeline per
// repo with each repo's freshness, or a stack page with qualified links.
func (d *Daemon) handleFollowStatic(w http.ResponseWriter, r *http.Request, q url.Values, stackQuery bool) {
	s := d.followSnapshot()
	now := d.followNow(r)
	sources := d.followSources(s)
	status := http.StatusOK
	var body []byte
	var err error
	if stackQuery {
		var st surfacepage.Stack
		if len(s.Envelope.Repos) == 0 {
			status, st.Reason = http.StatusServiceUnavailable, d.followRefusal(s, surfacepage.FollowStackSchema, now.At).Reason
		} else {
			lines := map[string]surfacepage.FollowSource{}
			for _, src := range sources {
				lines[src.ID] = src
			}
			st, err = surfacepage.LayoutFollowStack(s.Envelope, q.Get("brane"), q.Get("at"), now.At, func(repo string) (string, bool) {
				return surfacepage.SourceLine(lines[repo], now.At)
			})
			if err != nil {
				status, st = http.StatusNotFound, surfacepage.Stack{Reason: err.Error()}
			}
		}
		body, err = surfacepage.RenderStack(st)
	} else {
		if len(s.Envelope.Repos) == 0 {
			status = http.StatusServiceUnavailable
		}
		body, err = surfacepage.RenderFollowTimeline(s.Envelope, surfacepage.BuildFollowModel(s.Envelope, sources, now), now)
	}
	if err != nil {
		d.writeWorkbenchInternalError(w, r, err)
		return
	}
	setSurfaceHeaders(w)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
