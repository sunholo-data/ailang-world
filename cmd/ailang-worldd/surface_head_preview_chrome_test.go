package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSurfaceHeadOnPreviewsRecordedNodeInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(index, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), index, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/surface/model.json" {
				next.ServeHTTP(w, r)
				return
			}
			response := httptest.NewRecorder()
			next.ServeHTTP(response, r)
			var model map[string]json.RawMessage
			var lanes []map[string]json.RawMessage
			if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &model) != nil ||
				json.Unmarshal(model["lanes"], &lanes) != nil {
				http.Error(w, "could not read preview fixture", http.StatusInternalServerError)
				return
			}
			for _, lane := range lanes {
				var states []map[string]json.RawMessage
				if err := json.Unmarshal(lane["states"], &states); err != nil {
					http.Error(w, "could not read fixture states", http.StatusInternalServerError)
					return
				}
				for _, state := range states {
					state["summary"] = json.RawMessage(`"Preserve the literal <img> text."`)
				}
				lane["states"], _ = json.Marshal(states)
			}
			model["lanes"], _ = json.Marshal(lanes)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model)
		})
	})
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	wait(`document.documentElement.dataset.ready === '1' &&
		document.querySelector('#stage')?.dataset.view === 'heads' &&
		(window.previewModelPromise ||= fetch('/surface/model.json').then(response => response.json()).then(model => { window.previewModel = model; }), true)`)
	wait(`(() => {
		const model = window.previewModel;
		if (!model) return false;
		const nodes = new Map(model.lanes.flatMap(lane => lane.states.map(state => [state.id, state])));
		for (const el of document.querySelectorAll('.hc-scene a:is([data-head-state], [data-rear-state], [data-lineage-state])')) {
			const id = el.dataset.headState || el.dataset.rearState || el.dataset.lineageState;
			const state = nodes.get(id);
			if (!state?.summary || !state.at || !state.pr ||
				!model.links.some(link => (link.from === id && nodes.has(link.to)) || (link.to === id && nodes.has(link.from)))) continue;
			el.scrollIntoView({block:'nearest',inline:'nearest'});
			const box = el.getBoundingClientRect();
			const x = box.x + box.width / 2, y = box.y + box.height / 2;
			const hit = document.elementFromPoint(x, y);
			if (!el.contains(hit)) continue;
			el.id = 'preview-test-node';
			window.previewTarget = {state, el, hit, x, y, url:location.href};
			hit.dispatchEvent(new PointerEvent('pointermove', {bubbles:true,clientX:x,clientY:y,pointerType:'mouse'}));
			return true;
		}
		return false;
	})()`)
	const richPreview = `(() => {
		const {el, state, url} = window.previewTarget;
		const tip = document.querySelector('.hc-scene [role="tooltip"]:not([hidden])');
		if (!tip || !el.getAttribute('aria-describedby')?.split(/\s+/).includes(tip.id) ||
			tip.querySelector('.hc-preview-summary')?.textContent !== state.summary.replace(/\s+-\s+/g,' ').replace(/\s{2,}/g,' ').trim() ||
			tip.querySelector('time')?.dateTime !== state.at ||
			!tip.querySelector('.hc-preview-title')?.textContent.trim() ||
			!tip.querySelector('.hc-preview-relationship')?.textContent.trim() ||
			![...tip.querySelectorAll('.hc-preview-meta')].some(meta => meta.textContent.includes('PR #' + state.pr)) ||
			tip.querySelector('script, iframe, img') ||
			location.href !== url || window.__cspv?.length !== 0) return false;
		const r = tip.getBoundingClientRect();
		return r.width > 0 && r.height > 0 && r.left >= 0 && r.top >= 0 && r.right <= innerWidth && r.bottom <= innerHeight;
	})()`
	wait(richPreview)
	wait(`(document.querySelector('.hc-scene').dispatchEvent(new PointerEvent('pointerleave', {pointerType:'mouse'})),
		document.querySelector('#head-angle').focus(), true)`)
	b.do(t, nil, "key", map[string]any{"key": "Tab"})
	wait(`(window.previewTarget.el.focus(), document.activeElement === window.previewTarget.el)`)
	wait(richPreview)
	b.do(t, nil, "key", map[string]any{"key": "Escape"})
	wait(`document.activeElement === window.previewTarget.el &&
		document.querySelector('.hc-scene [role="tooltip"]').hidden &&
		location.href === window.previewTarget.url`)
	b.do(t, nil, "tap", map[string]any{"id": "preview-test-node"})
	wait(`document.querySelector('#stage').dataset.view === 'brane' &&
		document.querySelector('#front')?.dataset.front === window.previewTarget.state.id`)
}

func TestSurfaceRecordKindsUseConsistentMarkersInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	check := func(view, selector string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => {
			if (document.querySelector('#stage').dataset.view !== %q) return false;
			const records = [...document.querySelectorAll(%q)];
			const marker = link => link.querySelector('.g');
			const labelled = (kind, label) => records.filter(el =>
				el.dataset.recordKind === kind && el.getAttribute('aria-label')?.includes(label) && marker(el));
			const merges = labelled('github-merge', 'GitHub merge');
			const episodes = labelled('episode', 'World episode');
			return merges.length > 0 && episodes.length > 0 &&
				merges.some(el => {
					const svg = marker(el);
					const paths = svg.querySelectorAll('path');
					const carrier = svg.closest('#tl, .hc-face-bar');
					return carrier && svg.matches('svg.merge-glyph[aria-hidden="true"]') && paths.length === 1 &&
						paths[0].getAttribute('d')?.split('M').length === 3 &&
						getComputedStyle(paths[0]).stroke === getComputedStyle(carrier).color;
				}) &&
				episodes.every(el => !marker(el).matches('.merge-glyph'));
		})()`, view, selector)})
	}
	check("timeline", "#tl .node[data-record-kind]")
	b.do(t, nil, "tap", map[string]any{"id": "head-angle"})
	check("heads", ".hc-scene :is(.hc-face,.hc-rear)[data-record-kind]")
}
