package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// The test adds a driver to the existing same-origin app asset. The page
// keeps its production shell and CSP; only this test server has report URLs.
// Chrome is driven and reaped by the workbench's AC10-covered launcher.
func TestSurfaceLiveChrome(t *testing.T) {
	chrome := os.Getenv("WORLD_CHROME")
	if chrome == "" {
		t.Skip("WORLD_CHROME unset; operator must run the real Chrome drill")
	}
	d := surfaceLiveDaemon(t)
	c := genesisCommit("surface-chrome")
	genesisIndex := c.Entry.Header.EntryIndex
	firstCommit := true
	if err := d.store.Commit(boundedTestContext(t), c); err != nil {
		t.Fatal(err)
	}
	d.aguiBudget, d.aguiTick = time.Second, 10*time.Millisecond
	handler := d.Handler()
	var mu sync.Mutex
	models, stacks, posts := 0, 0, 0
	var firstCursor, modelCursor int64
	var traceErr error
	reports := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/surface/app.js":
			w.Header().Set("Content-Type", "text/javascript")
			_, _ = w.Write(surfacepage.AppScript)
			_, _ = io.WriteString(w, surfaceLiveChromeDriver)
			return
		case "/__surface_live":
			stage := r.URL.Query().Get("stage")
			if stage == "commit" {
				mu.Lock()
				if firstCommit {
					if posts == 0 || firstCursor != modelCursor || firstCursor != genesisIndex {
						traceErr = fmt.Errorf("first POST=%d model=%d posts=%d", firstCursor, modelCursor, posts)
					}
					firstCommit = false
				}
				mu.Unlock()
				next := testCommit(c.NextWorld, c.Entry.Header.EntryIndex+1, "surface-chrome")
				if err := d.store.Commit(boundedTestContext(t), next); err != nil {
					mu.Lock()
					traceErr = err
					mu.Unlock()
				}
				c = next
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]int64{"index": c.Entry.Header.EntryIndex})
				return
			}
			reports <- stage + ":" + r.URL.Query().Get("detail")
			w.WriteHeader(204)
			return
		case "/surface/model.json":
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, r)
			var m surfacepage.Model
			err := json.Unmarshal(rec.Body.Bytes(), &m)
			mu.Lock()
			models++
			if err != nil || m.NowEdge.Live == nil {
				traceErr = fmt.Errorf("model cursor missing: %s", rec.Body.String())
			} else if models == 1 {
				modelCursor = m.NowEdge.Live.LastIndex
			}
			mu.Unlock()
			for k, vs := range rec.Header() {
				w.Header()[k] = vs
			}
			w.WriteHeader(rec.Code)
			_, _ = w.Write(rec.Body.Bytes())
			return
		case "/surface/stack.json":
			mu.Lock()
			stacks++
			mu.Unlock()
		case "/agui/":
			b, err := io.ReadAll(r.Body)
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(b))
			var input struct {
				State struct {
					LastIndex int64 `json:"lastIndex"`
				} `json:"state"`
			}
			if err == nil {
				err = json.Unmarshal(b, &input)
			}
			mu.Lock()
			if err != nil {
				traceErr = err
			}
			if posts == 0 {
				firstCursor = input.State.LastIndex
			}
			posts++
			mu.Unlock()
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	startWorkbenchChrome(t, chrome, server.URL+"/surface", false)
	guard := time.NewTimer(20 * time.Second)
	defer guard.Stop()
	for {
		select {
		case report := <-reports:
			if strings.HasPrefix(report, "error:") {
				t.Fatal(report)
			}
			mu.Lock()
			err, nm, ns := traceErr, models, stacks
			mu.Unlock()
			if err != nil {
				t.Fatal(err)
			}
			if report == "timeline:" && nm < 2 {
				t.Fatalf("model not re-fetched: %d", nm)
			}
			if report == "done:" {
				if nm < 3 || ns < 2 {
					t.Fatalf("missing refresh: models=%d stacks=%d", nm, ns)
				}
				return
			}
		case <-guard.C:
			t.Fatal("timed out waiting for Chrome live refresh assertions")
		}
	}
}

const surfaceLiveChromeDriver = `
;(async () => {
  const sleep = () => new Promise(r => setTimeout(r, 25));
  async function until(f, label) {
    for (let i=0; i<600; i++) { if (f()) return; await sleep(); }
    throw Error('timeout: ' + label);
  }
  const report = (stage, detail='') => fetch('/__surface_live?stage='+stage+'&detail='+encodeURIComponent(detail));
  window.addEventListener('error', e => report('error', e.message));
  window.addEventListener('unhandledrejection', e => report('error', String(e.reason)));
  document.addEventListener('securitypolicyviolation', e => report('error', 'CSP: ' + e.violatedDirective));
  try {
    await until(() => document.documentElement.dataset.ready === '1' && document.querySelector('[data-now]'), 'now card');
    await until(() => /^Live · up to entry /.test(document.querySelector('[data-live]').textContent), 'initial live status');
    World.zoom(0.5); World.pan(World.view().L - 180);
    const before = JSON.stringify(World.view());
    const commit = async () => (await (await report('commit')).json()).index;
    const checked = index => new RegExp('^Live · up to entry '+index+' · checked [0-9]{2}:[0-9]{2}$').test(document.querySelector('[data-live]').textContent);
    const index = await commit();
    await until(() => checked(index), 'timeline refreshed');
    if (JSON.stringify(World.view()) !== before) throw Error('view changed during refresh');
    await report('timeline');
    const model = await (await fetch('/surface/model.json')).json();
    const state = model.lanes.flatMap(lane => lane.states)
      .filter(state => state.kind === 'step')
      .sort((a, b) => Date.parse(b.at) - Date.parse(a.at))[0];
    if (!state || !World.open(state.id)) throw Error('fixture brane missing');
    await until(() => document.querySelector('[data-front]'), 'brane open');
    const front = document.querySelector('[data-front]').dataset.front, url = location.href;
    const next = await commit();
    await until(() => checked(next), 'brane refreshed');
    if (document.querySelector('[data-front]').dataset.front !== front) throw Error('front changed');
    if (location.href !== url) throw Error('URL changed');
    if (JSON.stringify(World.view()) !== before) throw Error('brane refresh changed view');
    await report('done');
  } catch (e) { await report('error', e.stack || e.message); }
})();
`
