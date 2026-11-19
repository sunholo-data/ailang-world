package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

// The browser follows the production links and uses the real index-backed
// model and stack. The selected state predates the lane tip, so opening the
// latest state for every history selection cannot pass this drill.
func TestSurfaceHeadsOpensExactHistoryInChrome(t *testing.T) {
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
	open("/surface")
	wait(`document.documentElement.dataset.ready === '1' &&
		document.querySelector('#heads-link')?.getAttribute('href') === '/surface/heads'`)
	b.do(t, nil, "tap", map[string]any{"id": "heads-link"})

	const lane = "row-189"
	const rearState = "f36428abd9a72e5093b61cba14e716c37019f1ef"
	const earlierState = "59bbded417388861f3acae45c07eb9046f3cacaa"
	const tip = "652c341d5072c5f49fefe29f9409c93a5119a5be"
	stackPath := func(state string) string { return "/surface?brane=" + lane + "&at=" + state + "#heads" }
	wait(fmt.Sprintf(`location.pathname === '/surface/heads' &&
		[...document.querySelectorAll('h1')].some(h => h.textContent.trim() === 'Branch heads' && h.getClientRects().length) &&
		document.querySelector('.wh-card[data-head=%q] .wh-open')?.getAttribute('href') === %q`,
		lane, stackPath(tip)))

	focus := func(selector string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const control = document.querySelector(%q);
			if (!control || !control.getClientRects().length) return false;
			control.focus();
			return document.activeElement === control;
		})()`, selector))
	}
	enter := func() {
		t.Helper()
		b.do(t, nil, "key", map[string]any{"key": "Enter"})
	}

	assertStackAndReturn := func(state string) {
		t.Helper()
		wait(fmt.Sprintf(`location.pathname + location.search + location.hash === %q &&
			document.documentElement.dataset.ready === '1' && document.querySelector('#stage')?.dataset.view === 'brane' &&
			document.querySelector('#front')?.dataset.front === %q &&
			document.querySelector('#heads-link')?.getAttribute('aria-label') === 'Back to Branch heads' &&
			document.querySelector('#heads-link')?.getAttribute('href') === %q`, stackPath(state), state, "/surface/heads#"+lane))
		b.do(t, nil, "tap", map[string]any{"id": "heads-link"})
		wait(fmt.Sprintf(`location.pathname === '/surface/heads' && location.hash === %q &&
			document.querySelector('.wh-card[data-head=%q].is-selected .wh-open')?.getClientRects().length > 0 &&
			window.__cspv?.length === 0`, "#"+lane, lane))
	}

	// Verify the rear pane can receive an actual pointer click above the front
	// card. The driver takes an element ID, so tag the matched link for the tap.
	rear := fmt.Sprintf(`.wh-card[data-head=%q] a.wh-rear[href=%q]`, lane, stackPath(rearState))
	wait(fmt.Sprintf(`(() => {
		const link = document.querySelector(%q);
		if (!link?.checkVisibility({visibilityProperty:true, opacityProperty:true})) return false;
		link.scrollIntoView({block:'center'});
		const box = link.getBoundingClientRect();
		if (!link.contains(document.elementFromPoint(box.x + box.width / 2, box.y + box.height / 2))) return false;
		link.id = 'chrome-rear-pane';
		return true;
	})()`, rear))
	b.do(t, nil, "tap", map[string]any{"id": "chrome-rear-pane"})
	assertStackAndReturn(rearState)

	// Browse the full lane history and choose a different older state by keyboard.
	focus(fmt.Sprintf(`.wh-card[data-head=%q] .wh-all`, lane))
	enter()
	focus(fmt.Sprintf(`.wh-history[open] a[href=%q]`, stackPath(earlierState)))
	enter()
	assertStackAndReturn(earlierState)
}
