package daemon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// TestSurfaceFollowChrome is the row 220 M4 browser drill (AC13), in real
// headless Chrome against the production shell and CSP: AG-UI is refused
// throughout, yet a new merge in repo1 reaches the page by the model poll,
// with pan and zoom kept; a model/stack generation mismatch keeps the brane
// view and says so; and the 15 s visible-page timer fires on its own.
func TestSurfaceFollowChrome(t *testing.T) {
	chrome := os.Getenv("WORLD_CHROME")
	if chrome == "" {
		t.Skip("WORLD_CHROME unset; operator must run the real Chrome drill")
	}
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	full := rigFollowIndex(t)
	// Generation g1: repo1 lacks its newest merge (the last step and shared state).
	short := full
	short.States = append([]surface.State(nil), full.States[:len(full.States)-2]...)
	short.Source.Head, short.Source.Merges = short.States[len(short.States)-1].ID, full.Source.Merges-1
	snap := func(gen string, repo1 surface.Index) *surfacefollow.Snapshot {
		s := followSnapshotFor(t, gen, []string{"repo0", "repo1"}, nil)
		s.Envelope.Repos[1].Index, s.Envelope.Repos[1].Head = repo1, repo1.Source.Head
		s.Envelope.Open = nil
		if err := surface.ValidateFollow(s.Envelope); err != nil {
			t.Fatal(err)
		}
		s.Bytes, _ = json.Marshal(s.Envelope)
		return s
	}
	// The page's requests are sequential, so a stack request can be served
	// from a different (newer) snapshot than the model before it.
	var mu sync.Mutex
	current, skew, override := snap("g1", short), (*surfacefollow.Snapshot)(nil), (*surfacefollow.Snapshot)(nil)
	d.follow.snapshot = func() *surfacefollow.Snapshot {
		mu.Lock()
		defer mu.Unlock()
		if override != nil {
			return override
		}
		return current
	}
	handler := d.Handler()
	reports := make(chan string, 16)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/surface/app.js":
			w.Header().Set("Content-Type", "text/javascript")
			_, _ = w.Write(surfacepage.AppScript)
			_, _ = io.WriteString(w, strings.ReplaceAll(surfaceFollowChromeDriver, "__MERGE__", "repo1:commit:"+full.Source.Head))
			return
		case "/agui/":
			w.Header().Set("Retry-After", "30")
			http.Error(w, "AG-UI disconnected for this drill", http.StatusServiceUnavailable)
			return
		case "/__follow":
			mu.Lock()
			switch r.URL.Query().Get("stage") {
			case "merge":
				current = snap("g2", full)
			case "skew":
				current, skew = snap("g3", full), snap("g4", full)
			case "unskew":
				skew = nil
			case "timer":
				current = snap("g5", full)
			default:
				reports <- r.URL.Query().Get("stage") + ":" + r.URL.Query().Get("detail")
			}
			mu.Unlock()
			w.WriteHeader(204)
			return
		case "/surface/stack.json":
			// While skewed, the stack is always one publication ahead.
			mu.Lock()
			override = skew
			mu.Unlock()
			handler.ServeHTTP(w, r)
			mu.Lock()
			override = nil
			mu.Unlock()
			return
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	startWorkbenchChrome(t, chrome, server.URL+"/surface", false)
	guard := time.NewTimer(45 * time.Second)
	defer guard.Stop()
	for {
		select {
		case report := <-reports:
			if strings.HasPrefix(report, "error:") {
				t.Fatal(report)
			}
			if report == "done:" {
				return
			}
			t.Log(report)
		case <-guard.C:
			t.Fatal("timed out waiting for the follow drill")
		}
	}
}

const surfaceFollowChromeDriver = `
;(async () => {
  const sleep = ms => new Promise(r => setTimeout(r, ms || 25));
  async function until(f, label, n) {
    for (let i=0; i<(n || 600); i++) { if (f()) return; await sleep(); }
    throw Error('timeout: ' + label);
  }
  const report = (stage, detail='') => fetch('/__follow?stage='+stage+'&detail='+encodeURIComponent(detail));
  window.addEventListener('error', e => report('error', e.message));
  window.addEventListener('unhandledrejection', e => report('error', String(e.reason)));
  document.addEventListener('securitypolicyviolation', e => report('error', 'CSP: ' + e.violatedDirective));
  const q = s => document.querySelector(s);
  try {
    await until(() => document.documentElement.dataset.ready === '1' && q('[data-now]'), 'now card');
    if (World.generation() !== 'g1') throw Error('first generation ' + World.generation());
    for (const k of ['repo0:brane:shared', 'repo1:brane:shared']) if (!q('[data-lane="' + k + '"]')) throw Error('lane ' + k + ' missing');
    if (!q('[data-sources] [data-repo="repo0"]') || !q('[data-sources] [data-repo="repo1"]')) throw Error('sources missing');
    await until(() => /^(Paused|Live updates unavailable)/.test(q('[data-live]').textContent), 'AG-UI refused');
    World.zoom(0.5); World.pan(World.view().L - 180);
    const keep = () => { const v = World.view(); return JSON.stringify([v.u, v.L]); }, before = keep();
    if (q('[data-state="__MERGE__"]')) throw Error('the new merge is already drawn');
    await report('merge');
    await World.poll();
    await until(() => World.generation() === 'g2' && q('[data-state="__MERGE__"]'), 'merge drawn');
    if (keep() !== before) throw Error('view changed by the poll: ' + before + ' → ' + keep());
    await report('step', 'merge visible with AG-UI down');
    // A brane open while the stack straddles two publications: retry once, then keep the coherent view.
    const model = await (await fetch('/surface/model.json')).json();
    const step = model.lanes.filter(l => l.repo === 'repo0').flatMap(l => l.states).filter(s => s.kind === 'step')[0];
    if (!World.open(step.id)) throw Error('no repo0 step');
    await until(() => q('[data-front]') && q('[data-front]').dataset.front === step.id, 'brane open');
    const url = location.href;
    await report('skew');
    await World.poll();
    await until(() => q('[data-poll]') && /graph changed during the refresh/.test(q('[data-poll]').textContent), 'kept view note');
    if (World.generation() !== 'g2' || q('[data-front]').dataset.front !== step.id || location.href !== url) throw Error('mismatch replaced the view');
    await report('unskew');
    await World.poll();
    await until(() => World.generation() === 'g3' && !q('[data-poll]'), 'coherent pair applied');
    if (q('[data-front]').dataset.front !== step.id || location.href !== url) throw Error('front or URL changed');
    await report('step', 'mismatch kept, then coherent');
    // The visible-page timer: no manual poll.
    await report('timer');
    await until(() => World.generation() === 'g5', '15 s poll', 800);
    await report('done');
  } catch (e) { await report('error', e.stack || e.message); }
})();
`
