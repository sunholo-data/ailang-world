package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
)

func TestSurfaceHeadLaneStripInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	index := filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json")
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), index, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/surface/model.json" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, headThreadBrowserModel)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}

	wait(`(() => {
		const strip = document.querySelector('#head-toolbar .hc-lanes:not([hidden])');
		const select = strip?.querySelector('select[aria-label="Thread"]');
		if (!select) return false;
		const stage = strip.closest('#stage');
		const bounds = strip.getBoundingClientRect(), toolbar = strip.parentElement.getBoundingClientRect();
		const keys = [...select.options].map(option => option.value);
		if (bounds.top !== toolbar.top || bounds.left !== toolbar.left || bounds.height !== 44 ||
			bounds.width > 230 || keys.join(',') !== 'primary:brane:shared,primary:brane:work,secondary:brane:shared,secondary:brane:work' ||
			select.options[0].text !== 'Primary world' || select.value !== keys[0] ||
			select.options[0].parentElement.label !== 'primary' || select.options[2].parentElement.label !== 'secondary') return false;
		window.thread = {strip, select, stage, scene:stage.querySelector('.hc-scene')};
		select.focus();
		return document.activeElement === select &&
			select.dispatchEvent(new KeyboardEvent('keydown', {key:'Escape', bubbles:true, cancelable:true})) &&
			stage.dataset.view === 'heads';
	})()`)
	wait(`(() => {
		const {select} = window.thread;
		select.value = 'secondary:brane:work';
		select.dispatchEvent(new Event('change', {bubbles:true}));
		return true;
	})()`)
	wait(`window.thread.scene === document.querySelector('.hc-scene') &&
		window.thread.scene.dataset.viewLane === 'secondary:brane:work' &&
		window.thread.select.value === 'secondary:brane:work' &&
		document.activeElement === window.thread.select`)
	pollSurface(t, b)
	wait(`window.thread.select.isConnected && window.thread.select.value === 'secondary:brane:work'`)
	wait(`(() => {
		const face = window.thread.scene.querySelector('a[data-head-state="b-tip"]');
		if (!face) return false;
		face.focus();
		face.dispatchEvent(new KeyboardEvent('keydown', {key:'ArrowUp', bubbles:true, cancelable:true}));
		return true;
	})()`)
	wait(`window.thread.scene.dataset.travelState === 'b-early' &&
		window.thread.select.value === 'secondary:brane:work' &&
		[...window.thread.select.options].map(option => option.value).join(',') === 'primary:brane:shared,secondary:brane:work'`)
	b.do(t, nil, "tap", map[string]any{"id": "timeline-angle"})
	wait(`window.thread.stage.dataset.view === 'timeline' && window.thread.strip.hidden`)
}

const headThreadBrowserModel = `{
	"schema": "world/surface-model/v1",
	"now": "2026-10-20T15:00:00Z",
	"nowEdge": {},
	"repos": [
		{"id":"primary","shared":"primary:brane:shared","branch":"main"},
		{"id":"secondary","shared":"secondary:brane:shared","branch":"main"}
	],
	"lanes": [
		{"key":"secondary:brane:work","repo":"secondary","name":"Secondary work","glyph":"square","colour":"#42669a","states":[
			{"id":"b-early","kind":"step","label":"First work","at":"2026-10-20T11:00:00Z","rank":2},
			{"id":"b-tip","kind":"step","label":"Latest work","at":"2026-10-20T14:30:00Z","rank":5}
		]},
		{"key":"primary:brane:work","repo":"primary","name":"Primary work","glyph":"square","colour":"#42669a","states":[
			{"id":"a-tip","kind":"step","label":"Project work","at":"2026-10-20T14:15:00Z","rank":4}
		]},
		{"key":"secondary:brane:shared","repo":"secondary","name":"Secondary world","glyph":"circle","colour":"#8496a5","states":[
			{"id":"b-shared","kind":"shared","label":"Shared update","at":"2026-10-20T13:00:00Z","rank":3}
		]},
		{"key":"primary:brane:shared","repo":"primary","name":"Primary world","glyph":"circle","colour":"#8496a5","states":[
			{"id":"a-shared","kind":"shared","label":"Shared update","at":"2026-10-20T09:00:00Z","rank":1}
		]}
	],
	"links":[]
}`
