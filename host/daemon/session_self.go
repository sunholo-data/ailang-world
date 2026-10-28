package daemon

// Session self is a host transport projection, not a new authority or package.
import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

type SessionGrant struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Budget int64  `json:"budget"`
	Used   int64  `json:"used"`
	Left   int64  `json:"left"`
}

type SessionSelf struct {
	Git         *SessionGit     `json:"git,omitempty"`
	Episode     string          `json:"episode"`
	Worktree    string          `json:"worktree"`
	ExpiresAt   int64           `json:"expiresAt"`
	SecondsLeft int64           `json:"secondsLeft"`
	Grants      []SessionGrant  `json:"grants"`
	Entries     *SessionEntries `json:"entries,omitempty"`
	Addr        string          `json:"addr"`
	Surface     string          `json:"surface"`
}

type SessionEntries struct {
	First int64 `json:"first"`
	Last  int64 `json:"last"`
}

func (d *Daemon) handleSessionSelf(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := d.readCtx(r)
	defer cancel()
	out, err := d.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())
	if err != nil {
		if timedOut(ctx, err) {
			writeReadTimeout(w, d.readDeadline)
		} else {
			d.writeInternalError(w, r, err)
		}
		return
	}
	if out.Denied != nil {
		// The self route always returns 401, preserving MCP's fixed messages.
		writeAPIError(w, classFor(*out.Denied), msgFor(*out.Denied), http.StatusUnauthorized)
		return
	}
	b := out.Success
	rep, known, err := buildEpisodeReport(ctx, d.intents, d.reads, b.EpisodeID)
	if err != nil {
		d.writeInternalError(w, r, err)
		return
	}
	addr := "http://" + r.Host
	body := SessionSelf{Episode: b.EpisodeID, ExpiresAt: b.ExpiresAt,
		SecondsLeft: max(int64(0), b.ExpiresAt-time.Now().Unix()),
		Grants:      []SessionGrant{}, Addr: addr, Surface: addr + "/surface"}
	if d.workspace != nil && d.workspace.root != "" {
		body.Worktree = filepath.Join(d.workspace.root, b.EpisodeID)
	}
	if known {
		body.Entries = &SessionEntries{First: rep.Entries[0], Last: rep.Entries[1]}
	}
	spend, err := d.store.EffectSpend(ctx, b.EpisodeID)
	if err != nil {
		d.writeInternalError(w, r, err)
		return
	}
	for _, g := range b.Caps {
		used := spend[store.EffectKey{Effect: g.Effect, Scope: g.Scope}]
		body.Grants = append(body.Grants, SessionGrant{g.Effect, g.Scope, g.Budget, used, max(int64(0), g.Budget-used)})
	}
	body.Git = readSessionGit(ctx, body.Worktree, "git")
	writeJSON(w, http.StatusOK, body)
}
