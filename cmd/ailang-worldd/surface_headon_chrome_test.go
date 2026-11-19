package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

// Use the daemon's real index and stack endpoints. Retaining live DOM references
// detects any projection change that replaces the scene or reloads the page.
func TestSurfaceHeadOnUsesSameSceneInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)

	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	open := func(path string) {
		t.Helper()
		b.do(t, nil, "openCard", map[string]any{"url": origin + path, "raw": true})
	}
	enter := func() {
		t.Helper()
		b.do(t, nil, "key", map[string]any{"key": "Enter"})
	}

	open("/surface")
	wait(`document.documentElement.dataset.ready === '1' &&
		document.querySelector('#stage')?.dataset.view === 'timeline'`)
	wait(`(() => {
		const control = document.querySelector('#head-angle');
		if (!control?.getClientRects().length) return false;
		window.sceneBeforeTurn = { document, stage: document.querySelector('#stage'), timeline: document.querySelector('#tl') };
		control.focus();
		return document.activeElement === control;
	})()`)
	enter()

	const sameScene = `document === window.sceneBeforeTurn?.document &&
		document.querySelector('#stage') === window.sceneBeforeTurn?.stage &&
		document.querySelector('#tl') === window.sceneBeforeTurn?.timeline`
	const fitsViewport = `document.documentElement.scrollWidth <= document.documentElement.clientWidth + 1`
	const visibleFocus = `document.activeElement !== document.body &&
		document.activeElement.checkVisibility({visibilityProperty:true, opacityProperty:true})`
	wait(sameScene + ` && document.querySelector('#stage').dataset.view === 'heads' &&
		document.activeElement.matches('#stage a[data-head-state]') && ` + visibleFocus + ` && ` + fitsViewport)
	wait(`(() => {
		const face = document.activeElement;
		if (!face.matches('#stage a[data-head-state][href]')) return false;
		const href = new URL(face.href);
		const state = face.dataset.headState;
		if (href.pathname !== '/surface' || href.searchParams.get('at') !== state ||
			href.searchParams.get('brane') !== face.dataset.headLane || !state) return false;
		window.chosenHeadState = state;
		return true;
	})()`)
	enter()
	wait(sameScene + ` && document.querySelector('#stage').dataset.view === 'brane' &&
		document.querySelector('#front')?.dataset.front === window.chosenHeadState`)

	wait(`document.querySelector('#head-angle')?.getClientRects().length > 0`)
	b.do(t, nil, "tap", map[string]any{"id": "head-angle"})
	wait(sameScene + ` && document.querySelector('#stage').dataset.view === 'heads' && ` + visibleFocus + ` && ` + fitsViewport)

	wait(`document.querySelector('#timeline-angle')?.getClientRects().length > 0`)
	b.do(t, nil, "tap", map[string]any{"id": "timeline-angle"})
	wait(sameScene + ` && document.querySelector('#stage').dataset.view === 'timeline' &&
		document.querySelector('#tl .node')?.getClientRects().length > 0 &&
		window.__cspv?.length === 0 && ` + visibleFocus + ` && ` + fitsViewport)

	// Set the real persisted motion preference, then exercise the direct URL.
	wait(`(localStorage.setItem('world.appearance', JSON.stringify({motion:'reduced'})), true)`)
	open("/surface#heads")
	wait(`document.documentElement.dataset.ready === '1' && location.hash === '#heads' &&
		document.querySelector('#root').classList.contains('reduced') &&
		document.querySelector('#stage')?.dataset.view === 'heads' &&
		document.querySelector('#stage a[data-head-state]')?.getClientRects().length > 0 &&
		window.__cspv?.length === 0 && ` + fitsViewport)
	// Close the live page before the daemon fixture waits on its event stream.
	b.do(t, nil, "tab", nil)
}

func TestSurfaceHeadOnTravelsHistoryInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})

	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	tap := func(id string) {
		t.Helper()
		b.do(t, nil, "tap", map[string]any{"id": id})
	}
	tapDepth := func(name string) {
		t.Helper()
		var id string
		b.do(t, &id, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => {
			const button = [...document.querySelectorAll('button')].find(el =>
				(el.getAttribute('aria-label') || el.textContent.trim()) === %q &&
				!el.disabled && el.getAttribute('aria-disabled') !== 'true' && el.getClientRects().length);
			if (!button) return false;
			if (!button.id) button.id = 'chrome-depth-' + %q.replaceAll(' ', '-').toLowerCase();
			return button.id;
		})()`, name, name)})
		tap(id)
	}

	wait(`document.documentElement.dataset.ready === '1' &&
		(localStorage.setItem('world.appearance', JSON.stringify({motion:'reduced'})), true)`)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait(`document.documentElement.dataset.ready === '1' &&
		document.querySelector('#root').classList.contains('reduced') &&
		document.querySelector('#stage')?.dataset.view === 'heads'`)
	wait(`(() => {
		window.travelScene = document.querySelector('.hc-scene');
		window.travelStage = document.querySelector('#stage');
		window.getTravelStop = () => window.travelScene.dataset.travelState;
		window.visibleTravelFaces = () => [...window.travelScene.querySelectorAll('a[data-head-state][href]')]
			.filter(link => link.checkVisibility({visibilityProperty:true, opacityProperty:true}));
		fetch('/surface/model.json').then(response => response.json()).then(model => {
			window.travelRanks = new Map(model.lanes.flatMap(lane => lane.states.map(state => [state.id, state.rank])));
		});
		return Boolean(window.travelScene);
	})()`)
	wait(`window.travelRanks instanceof Map && window.getTravelStop() === '' &&
		window.visibleTravelFaces().length > 0`)

	assertStop := func(slot, newer string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const stop = window.getTravelStop();
			const cutoff = window.travelRanks.get(stop);
			const faces = window.visibleTravelFaces();
			const selected = faces.find(link => link.dataset.headState === stop);
			if (!selected || !Number.isFinite(cutoff) || stop === window[%q] ||
				(window[%q] && cutoff > window.travelRanks.get(window[%q]))) return false;
			const bounds = selected.getBoundingClientRect();
			if (bounds.right <= 0 || bounds.bottom <= 0 || bounds.left >= innerWidth || bounds.top >= innerHeight) return false;
			const url = new URL(selected.href);
			if (url.pathname !== '/surface' || url.searchParams.get('at') !== stop ||
				url.searchParams.get('brane') !== selected.dataset.headLane) return false;
			if (!faces.every(face => window.travelRanks.get(face.dataset.headState) <= cutoff)) return false;
			window[%q] = stop;
			return document.querySelector('.hc-scene') === window.travelScene &&
				document.querySelector('#stage[data-view="heads"]') === window.travelStage;
		})()`, newer, newer, newer, slot))
	}

	tapDepth("Move into history")
	assertStop("firstDepthStop", "")
	tapDepth("Move into history")
	assertStop("secondDepthStop", "firstDepthStop")

	// The browser harness has no physical +/- keys. Dispatch keydowns from a
	// focused head to exercise the same production keyboard handling.
	depthKey := func(key string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const link = window.visibleTravelFaces().find(face => face.dataset.headState === window.getTravelStop());
			if (!link) return false;
			link.focus();
			link.dispatchEvent(new KeyboardEvent('keydown', {key:%q, bubbles:true, cancelable:true}));
			return true;
		})()`, key))
	}

	// Open the selected historical stop by its native link and check the
	// camera resumes at that stop, with the reverse path still intact.
	wait(`(() => {
		const face = window.visibleTravelFaces().find(link => link.dataset.headState === window.secondDepthStop);
		if (!face) return false;
		face.focus();
		return document.activeElement === face;
	})()`)
	b.do(t, nil, "key", map[string]any{"key": "Enter"})
	wait(`document.querySelector('#stage[data-view="brane"]') === window.travelStage &&
		document.querySelector('#front')?.dataset.front === window.secondDepthStop`)
	tap("head-angle")
	wait(`document.querySelector('#stage[data-view="heads"]') === window.travelStage &&
		window.getTravelStop() === window.secondDepthStop`)

	tapDepth("Move toward present")
	wait(`window.getTravelStop() === window.firstDepthStop`)
	depthKey("+")
	wait(`window.getTravelStop() === window.secondDepthStop`)
	depthKey("-")
	wait(`window.getTravelStop() === window.firstDepthStop`)
	tapDepth("Return to present")
	wait(`window.getTravelStop() === '' && window.visibleTravelFaces().length > 0 &&
		document.querySelector('#stage[data-view="heads"]') === window.travelStage &&
		window.__cspv?.length === 0`)
}

func TestSurfaceHeadOnInspectsExactStateInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	wait(`document.documentElement.dataset.ready === '1' &&
		document.querySelector('#stage')?.dataset.view === 'heads'`)
	wait(`(() => {
		const scene = document.querySelector('.hc-scene');
		const head = scene.querySelector('.hc-face');
		if (!head) return false;
		head.focus();
		document.querySelector('#head-angle').focus();
		window.inspectionBaseline = {scene, url:location.href, selected:head.dataset.headState};
		fetch('/surface/model.json').then(res => res.json()).then(model => window.inspectionModel = model);
		return true;
	})()`)
	wait(`(() => {
		const model = window.inspectionModel;
		if (!model) return false;
		const {scene, selected} = window.inspectionBaseline;
		const steps = new Set(model.lanes.flatMap(lane => lane.states.filter(state => state.kind === 'step').map(state => state.id)));
		const panes = [...scene.querySelectorAll('.hc-rear[data-rear-state]'), ...scene.querySelectorAll('.hc-face[data-head-state]')];
		for (const pane of panes) {
			const id = pane.dataset.rearState || pane.dataset.headState;
			if (id === selected || !steps.has(id)) continue;
			const incident = model.links.filter(link => link.from === id || link.to === id);
			if (!incident.length) continue;
			pane.scrollIntoView({block:'nearest', inline:'nearest'});
			const rect = pane.getBoundingClientRect();
			const x = rect.left + rect.width / 2, y = rect.top + rect.height / 2;
			const hit = document.elementFromPoint(x, y);
			if (!pane.contains(hit)) continue;
			window.inspectedPane = {id, incident};
			hit.dispatchEvent(new PointerEvent('pointermove', {bubbles:true, clientX:x, clientY:y, pointerType:'mouse'}));
			return true;
		}
		return false;
	})()`)
	wait(`(() => {
		const {scene, selected, url} = window.inspectionBaseline;
		const {id, incident} = window.inspectedPane;
		if (scene.dataset.inspectState !== id ||
			Number(scene.querySelector('.hc-inspection').dataset.incidentCount) !== incident.length ||
			scene.querySelector('.hc-face[aria-current="true"]')?.dataset.headState !== selected ||
			location.href !== url || document.querySelector('#stage').dataset.view !== 'heads') return false;
		const neighbors = new Set(incident.map(edge => edge.from === id ? edge.to : edge.from));
		return [...scene.querySelectorAll('[data-head-state],[data-rear-state]')].every(el => {
			const state = el.dataset.headState || el.dataset.rearState;
			return el.classList.contains('is-related') === neighbors.has(state) &&
				el.classList.contains('is-inspected') === (state === id);
		});
	})()`)
	wait(`(window.inspectionBaseline.scene.dispatchEvent(new PointerEvent('pointerleave', {pointerType:'mouse'})), true)`)
	wait(`(() => {
		const {scene, selected, url} = window.inspectionBaseline;
		return !scene.dataset.inspectState && Number(scene.querySelector('.hc-inspection').dataset.incidentCount) === 0 &&
			!scene.querySelector('.is-inspected,.is-related') &&
			scene.querySelector('.hc-face[aria-current="true"]')?.dataset.headState === selected &&
			location.href === url && document.querySelector('#stage').dataset.view === 'heads' &&
			window.__cspv?.length === 0;
	})()`)
}
