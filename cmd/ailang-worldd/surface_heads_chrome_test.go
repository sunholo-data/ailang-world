package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
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
		new URL(document.querySelector('#heads-link').href).pathname === '/surface/heads'`)
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

	// An older lane may be folded in the present slice. Recall it through
	// the same navigator an operator uses before inspecting its history.
	wait(fmt.Sprintf(`(() => {
		const pick = document.querySelector('.wh-picker select');
		if (!pick || !Array.from(pick.options).some(o => o.value === %q)) return false;
		pick.value = %q;
		pick.dispatchEvent(new Event('change', {bubbles:true}));
		return true;
	})()`, lane, lane))

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

// The model can arrive first. Hold the stack response until Chrome has built
// the timeline, then verify that the user still sees only the stack loader.
func TestSurfaceHeadsDirectStackWhileLoadingInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	resume := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(resume) }) }
	index := filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json")
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), index, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/surface/stack.json" {
				select {
				case <-resume:
				case <-r.Context().Done():
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	})
	defer release()
	b.do(t, nil, "tab", nil)
	const id = "f36428abd9a72e5093b61cba14e716c37019f1ef"
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface?brane=row-189&at=" + id + "#heads", "raw": true})
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.querySelector('#tl .node') &&
		document.querySelector('#stage').dataset.view === 'stack-loading' &&
		document.querySelector('#L0').inert && getComputedStyle(document.querySelector('#L0')).visibility === 'hidden'`})
	release()
	b.do(t, nil, "waitFor", map[string]any{"expression": fmt.Sprintf(`document.querySelector('#stage').dataset.view === 'brane' &&
		document.querySelector('#front')?.dataset.front === %q`, id)})
}
