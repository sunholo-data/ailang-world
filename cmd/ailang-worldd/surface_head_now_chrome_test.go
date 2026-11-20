package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceHeadOnSharesNowCardWithTimelineInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"inject": "notify"})
	defer b.do(t, nil, "tab", nil)
	now := time.Now().Unix()
	pending := newApproveRig(t, "ep-pending", now)
	raw := rawCredentials(t, pending.credPath)
	defer raw.Close()
	for i := 1; i <= 4; i++ {
		decided := newApproveRigIn(t, pending.r, fmt.Sprintf("ep-decided%d", i), now)
		if _, err := raw.Exec("UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?", now-int64(60*i), decided.id); err != nil {
			t.Fatal(err)
		}
	}
	origin := nowCardOrigin(t, pending.credPath)
	visit := func(path string) {
		t.Helper()
		b.do(t, nil, "openCard", map[string]any{"url": origin + path, "raw": true})
	}
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	tap := func(id string) {
		t.Helper()
		b.do(t, nil, "tap", map[string]any{"id": id})
	}
	cardAt := func(view, container string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const card = document.querySelector('[data-now]');
			if (document.querySelector('#stage').dataset.view !== %q ||
				document.querySelectorAll('[data-now]').length !== 1 || !card?.closest(%q)) return false;
			return [...card.querySelectorAll('[id]')].every(el =>
				document.querySelectorAll('#' + CSS.escape(el.id)).length === 1);
		})()`, view, container))
	}
	expanded := func(open bool) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const toolbar = document.querySelector('#head-space > #head-toolbar');
			const scene = document.querySelector('#head-space > .hc-scene');
			const button = toolbar?.querySelector('#head-status-toggle');
			const panel = toolbar?.querySelector('#head-status-detail');
			return scene && toolbar.getBoundingClientRect().bottom <= scene.getBoundingClientRect().top + 1 &&
				button instanceof HTMLButtonElement && button.getAttribute('aria-controls') === panel?.id &&
				button.getAttribute('aria-expanded') === String(%t) && panel.checkVisibility() === %t;
		})()`, open, open))
	}

	visit("/surface#heads")
	cardAt("heads", "#head-status-detail")
	expanded(true)
	tap("head-status-toggle")
	expanded(false)
	tap("head-status-toggle")
	expanded(true)
	wait(fmt.Sprintf(`(() => {
		const panel = document.querySelector('#head-status-detail');
		const running = panel.querySelector('[data-running] button[data-go]');
		const needs = panel.querySelector(%q);
		const source = panel.querySelector('[data-sources] [data-repo="rig"]');
		return running?.checkVisibility() && running.textContent.includes('fixture-run') &&
			needs?.checkVisibility() && needs.textContent.includes('ep-pending') &&
			source?.checkVisibility() && source.textContent.includes('local branch checked') &&
			panel.querySelector('[data-decided] #decided-all')?.checkVisibility();
	})()`, "#needs-"+pending.id))

	tap("decided-all")
	wait(`document.querySelector('#head-status-detail #decided-all')?.getAttribute('aria-expanded') === 'true' &&
		document.querySelectorAll('#head-status-detail [data-decided] .nl.sub').length === 4`)
	b.do(t, nil, "grantNotifications", map[string]any{"origin": origin})
	tap("notify-me")
	wait(`document.querySelector('#head-status-detail [data-notify]')?.textContent === 'Notifications on' &&
		window.__asks?.length === 1 && window.__asks[0].active`)

	wait(`(window.nowCardBeforeTurn = document.querySelector('[data-now]'), true)`)
	tap("timeline-angle")
	cardAt("timeline", "#timeline-now")
	wait(`window.nowCardBeforeTurn === document.querySelector('#timeline-now [data-now]')`)
	tap("head-angle")
	cardAt("heads", "#head-status-detail")
	wait(`window.nowCardBeforeTurn === document.querySelector('#head-status-detail [data-now]')`)
	expanded(true)
	visit("/surface#heads")
	cardAt("heads", "#head-status-detail")
	expanded(true)
	tap("head-status-toggle")
	expanded(false)
	visit("/surface#heads")
	cardAt("heads", "#head-status-detail")
	expanded(false)
	tap("head-status-toggle")
	expanded(true)

	wait(`(() => {
		const running = document.querySelector('#head-status-detail [data-running] button[data-go]');
		if (!running) return false;
		running.id = 'fixture-running-action';
		window.expectedRunningState = running.dataset.go;
		return true;
	})()`)
	tap("fixture-running-action")
	wait(`document.querySelector('#stage').dataset.view === 'brane' &&
		document.querySelector('#front')?.dataset.front === window.expectedRunningState`)
	visit("/surface#heads")
	cardAt("heads", "#head-status-detail")
	expanded(true)
	tap("needs-" + pending.id)
	wait(fmt.Sprintf(`location.pathname === %q && document.body.dataset.state === 'ready'`, "/surface/approve/"+pending.id))
}

// Use the daemon's graph and approval store. The legacy rig has no live
// follower, so give its model the running and source records it would serve.
func nowCardOrigin(t *testing.T, creds string) string {
	t.Helper()
	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(index, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	return cardOrigin(t, creds, index, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/surface/model.json" {
				next.ServeHTTP(w, r)
				return
			}
			response := httptest.NewRecorder()
			next.ServeHTTP(response, r)
			var model struct {
				surfacepage.Model
				Repos []surfacepage.ModelRepo `json:"repos"`
			}
			if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &model) != nil {
				http.Error(w, "could not read fixture model", http.StatusInternalServerError)
				return
			}
			for _, lane := range model.Lanes {
				if len(lane.States) == 0 || lane.Key == "shared" {
					continue
				}
				model.NowEdge.Open = &surfacepage.ModelOpen{Status: "complete", Live: []surfacepage.ModelRunning{
					{Episode: "fixture-run", ID: lane.States[len(lane.States)-1].ID, Repo: "rig", Label: "Running"},
				}}
				break
			}
			model.Repos = []surfacepage.ModelRepo{{ID: "rig", State: "fresh", Line: "rig · local branch checked"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model)
		})
	})
}
