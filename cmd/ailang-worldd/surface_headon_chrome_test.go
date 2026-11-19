package main

import (
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
