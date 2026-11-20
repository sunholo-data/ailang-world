package main

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestSurfaceHeadOnReleasesGroundDragOverCardInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", map[string]any{"viewport": map[string]int{"width": 1280, "height": 844}})
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	type point struct{ X, Y float64 }
	read := func(expression string, out any) {
		t.Helper()
		b.do(t, out, "waitFor", map[string]any{"expression": expression})
	}
	mouse := func(event string, at point, buttons int) {
		t.Helper()
		b.do(t, nil, "mouse", map[string]any{"type": event, "x": at.X, "y": at.Y, "buttons": buttons})
	}
	var start point
	read(`(() => {
		const scene = document.querySelector('#stage[data-view="heads"] .hc-scene');
		if (!scene || scene.dataset.travelState) return false;
		const scroll = scene.querySelector('.hc-scroll');
		const rect = scroll.getBoundingClientRect();
		for (const row of [.53, .60, .67]) for (const col of [.5, .42, .58, .35, .65]) {
			const x = Math.round(rect.left + rect.width * col);
			const y = Math.round(rect.top + rect.height * row);
			if (!document.elementFromPoint(x,y)?.matches('.hc-floor-hit,.hc-graph')) continue;
			window.floorDrag = {scene, scroll, y, x};
			return {x,y};
		}
		return false;
	})()`, &start)
	mouse("mouseMoved", start, 0)
	mouse("mousePressed", start, 1)
	diagonal := point{start.X + 60, start.Y + 75}
	mouse("mouseMoved", diagonal, 1)
	read(`(() => {
		const {scene, scroll} = window.floorDrag;
		if (!scene.dataset.travelState || !scene.classList.contains('is-floor-dragging')) return false;
		window.floorDrag.firstCut = scene.dataset.travelState;
		window.floorDrag.firstPan = scroll.scrollLeft;
		return true;
	})()`, nil)
	mouse("mouseMoved", point{start.X - 60, diagonal.Y}, 1)
	var pan struct {
		Cut, FirstCut                      string
		Left, FirstPan, Width, ClientWidth float64
	}
	read(`(() => {
		const {scene, scroll, firstCut, firstPan} = window.floorDrag;
		return {cut:scene.dataset.travelState, firstCut, left:scroll.scrollLeft, firstPan, width:scroll.scrollWidth, clientWidth:scroll.clientWidth};
	})()`, &pan)
	if pan.Cut != pan.FirstCut || pan.Left-pan.FirstPan <= 10 {
		t.Fatalf("horizontal drag should pan within the same cut: %+v", pan)
	}

	var release struct {
		point
		Diagnostics json.RawMessage
	}
	read(`(() => {
		const {scene} = window.floorDrag;
		const diagnostics = [];
		for (const el of scene.querySelectorAll('.hc-face[data-head-state]')) {
			const rect = el.getBoundingClientRect();
			const x = Math.round(rect.left + rect.width / 2), y = Math.round(rect.top + rect.height / 2);
			const hit = document.elementFromPoint(x,y);
			if (!el.contains(hit)) { diagnostics.push({x,y,rect:rect.toJSON(),hit:hit?.className}); continue; }
			window.floorDrag.card = el.dataset.headState;
			return {x,y};
		}
		return {diagnostics};
	})()`, &release)
	if release.X == 0 || release.Y == 0 {
		t.Fatalf("no release point over a card: %s", release.Diagnostics)
	}
	mouse("mouseReleased", release.point, 0)
	read(`(() => {
		const {scene, firstCut, card} = window.floorDrag;
		return !scene.classList.contains('is-floor-dragging') && scene.dataset.travelState === firstCut &&
			document.querySelector('#stage').dataset.view === 'heads' && !!card;
	})()`, nil)
	mouse("mouseMoved", start, 0)
	read(`window.floorDrag.scene.dataset.travelState === window.floorDrag.firstCut &&
		!window.floorDrag.scene.classList.contains('is-floor-dragging')`, nil)
	mouse("mouseMoved", release.point, 0)
	mouse("mousePressed", release.point, 1)
	mouse("mouseReleased", release.point, 0)
	read(`document.querySelector('#stage').dataset.view === 'brane' &&
		document.querySelector('#front')?.dataset.front === window.floorDrag.card`, nil)
}
