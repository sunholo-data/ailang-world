package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
			window.floorDrag = {scene};
			return {x,y};
		}
		return false;
	})()`, &start)
	mouse("mouseMoved", start, 0)
	mouse("mousePressed", start, 1)
	mouse("mouseMoved", point{start.X, start.Y + 75}, 1)
	read(`(() => {
		const {scene} = window.floorDrag;
		if (!scene.dataset.travelState || !scene.classList.contains('is-floor-dragging')) return false;
		window.floorDrag.firstCut = scene.dataset.travelState;
		return true;
	})()`, nil)

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

func TestSurfaceHeadOnDragsBetweenThreadsLikeArrowKeysInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	index := filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json")
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), index, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/surface/model.json" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, recordedTimeBrowserModel)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	b.do(t, nil, "tab", map[string]any{"viewport": map[string]int{"width": 1280, "height": 844}})
	defer b.do(t, nil, "tab", nil)
	read := func(expression string, out any) {
		t.Helper()
		b.do(t, out, "waitFor", map[string]any{"expression": expression})
	}
	type view struct{ Lane, Node, Thread, Cut, Time string }
	open := func(historical bool) view {
		t.Helper()
		b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
		read(`(() => {
			const stage = document.querySelector('#stage[data-view="heads"]');
			const scene = stage?.querySelector('.hc-scene');
			if (scene?.dataset.viewLane !== 'shared' || scene.dataset.travelState) return false;
			const scroll = scene.querySelector('.hc-scroll');
			const face = () => scene.querySelector('a[data-head-state][aria-current="true"]');
			window.threadView = () => {
				const node = face();
				const marker = scene.querySelector('.hc-time-axis:not([hidden]) time:last-child');
				const thread = stage.querySelector('.hc-lanes select');
				if (!node || document.activeElement !== node || !marker || thread?.value !== scene.dataset.viewLane) return false;
				return {lane:scene.dataset.viewLane, node:node.dataset.headState, thread:thread.value,
					cut:scene.dataset.travelState, time:marker.dateTime};
			};
			window.threadCentered = () => {
				const node = face()?.getBoundingClientRect();
				return node && Math.abs(node.left + node.width / 2 - scroll.getBoundingClientRect().left - scroll.clientWidth / 2) <= 1;
			};
			window.threadArrow = key => {
				const node = key === 'ArrowUp' ? scene.querySelector('a[data-head-state="shared-late"]') : face();
				if (!node) return false;
				node.focus({preventScroll:true});
				node.dispatchEvent(new KeyboardEvent('keydown', {key, bubbles:true, cancelable:true}));
				return true;
			};
			const overflow = getComputedStyle(scroll);
			return overflow.overflowX === 'hidden' && ['auto','scroll'].includes(overflow.overflowY);
		})()`, nil)
		var initial view
		if historical {
			read(`window.threadArrow('ArrowUp')`, nil)
			read(`window.threadView()?.cut === 'shared-cut' && window.threadView()`, &initial)
		}
		return initial
	}
	for _, historical := range []bool{false, true} {
		t.Logf("verify lateral gestures with explicit historical cut = %t", historical)
		initial := open(historical)
		var right, middle, left view
		read(`window.threadArrow('ArrowRight')`, nil)
		read(`window.threadView()?.lane !== 'shared' && window.threadView()`, &right)
		read(`window.threadArrow('ArrowLeft')`, nil)
		read(`window.threadView()?.lane === 'shared' && window.threadView()`, &middle)
		read(`window.threadArrow('ArrowLeft')`, nil)
		read(fmt.Sprintf(`window.threadView()?.lane !== 'shared' && window.threadView()?.lane !== %q && window.threadView()`, right.Lane), &left)
		if historical {
			for _, position := range []view{right, middle, left} {
				if position.Cut != initial.Cut || position.Time != initial.Time {
					t.Fatalf("arrows changed historical scope: before %+v, after %+v", initial, position)
				}
			}
		}

		open(historical)
		var start struct{ X, Y float64 }
		read(`(() => {
			const scroll = document.querySelector('.hc-scroll');
			const rect = scroll.getBoundingClientRect();
			for (const row of [.55, .65, .75]) for (const col of [.5, .4, .6, .35, .65]) {
				const x = Math.round(rect.left + rect.width * col), y = Math.round(rect.top + rect.height * row);
				if (document.elementFromPoint(x,y)?.matches('.hc-floor-hit,.hc-graph')) return {x,y};
			}
			return false;
		})()`, &start)
		mouse := func(event string, dx float64, buttons int) {
			t.Helper()
			b.do(t, nil, "mouse", map[string]any{"type": event, "x": start.X + dx, "y": start.Y, "buttons": buttons})
		}
		mouse("mouseMoved", 0, 0)
		mouse("mousePressed", 0, 1)
		mouse("mouseMoved", -35, 1)
		read(`document.querySelector('.hc-scene').dataset.viewLane === 'shared'`, nil)
		mouse("mouseMoved", -105, 1)
		var actual view
		read(`window.threadView()?.lane !== 'shared' && window.threadView()`, &actual)
		if actual != right {
			t.Fatalf("drag left %+v, ArrowRight %+v", actual, right)
		}
		mouse("mouseMoved", 0, 1)
		read(`window.threadView()?.lane === 'shared' && window.threadView()`, &actual)
		if actual != middle {
			t.Fatalf("reversing drag %+v, ArrowLeft %+v", actual, middle)
		}
		mouse("mouseMoved", 105, 1)
		read(fmt.Sprintf(`window.threadView()?.lane === %q && window.threadView()`, left.Lane), &actual)
		if actual != left {
			t.Fatalf("drag right %+v, ArrowLeft %+v", actual, left)
		}
		mouse("mouseMoved", -105, 1)
		mouse("mouseReleased", -105, 0)
		read(fmt.Sprintf(`window.threadCentered() && window.threadView()?.lane === %q && window.threadView()`, right.Lane), &actual)
		if actual != right {
			t.Fatalf("drag across two lanes %+v, ArrowRight %+v", actual, right)
		}
	}
}
