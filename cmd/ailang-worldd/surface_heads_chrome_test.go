package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"sync/atomic"
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
	stackURL := func(state string) string { return "/surface?brane=" + lane + "&at=" + state }
	stackPath := func(state string) string { return stackURL(state) + "#heads" }
	timelinePath := func(state string) string { return "/surface#brane=" + lane + "&at=" + state }
	returnPath := func(state string) string { return "/surface/heads?at=" + state + "#" + lane }
	wait(`location.pathname === '/surface/heads' &&
		[...document.querySelectorAll('h1')].some(h => h.textContent.trim() === 'Branch heads' && h.getClientRects().length)`)

	// An older lane may be folded in the present slice. Recall it through
	// the same navigator an operator uses before inspecting its history.
	wait(fmt.Sprintf(`(() => {
		const pick = document.querySelector('.wh-picker select');
		if (!pick || !Array.from(pick.options).some(o => o.value === %q)) return false;
		pick.value = %q;
		pick.dispatchEvent(new Event('change', {bubbles:true}));
		return true;
	})()`, lane, lane))
	wait(fmt.Sprintf(`document.querySelector('.wh-card[data-head=%q].is-selected .wh-open')?.getAttribute('href') === %q &&
		document.querySelector('.wh-card[data-head=%q].is-selected .wh-map')?.getAttribute('href') === %q`,
		lane, stackPath(tip), lane, timelinePath(tip)))

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

	assertHeadsReturn := func(state string) {
		t.Helper()
		b.do(t, nil, "tap", map[string]any{"id": "heads-link"})
		pane := fmt.Sprintf(`.wh-card[data-head=%q].is-selected `, lane)
		if state == tip {
			pane += fmt.Sprintf(`.wh-open[data-head-state=%q]`, state)
		} else {
			pane += fmt.Sprintf(`.wh-rear[data-head-state=%q][aria-current="location"]`, state)
		}
		wait(fmt.Sprintf(`location.pathname + location.search + location.hash === %q &&
			document.querySelector(%q)?.getClientRects().length > 0 && window.__cspv?.length === 0`, returnPath(state), pane))
	}
	assertStackAndReturn := func(state string, fromHeads bool) {
		t.Helper()
		wait(fmt.Sprintf(`location.pathname + location.search === %q &&
			document.documentElement.dataset.ready === '1' && document.querySelector('#stage')?.dataset.view === 'brane' &&
			document.querySelector('#front')?.dataset.front === %q &&
			document.querySelector('#heads-link')?.getAttribute('href') === %q`, stackURL(state), state, returnPath(state)))
		if fromHeads {
			wait(`location.hash === '#heads' && document.querySelector('#heads-link')?.getAttribute('aria-label') === 'Back to Branch heads'`)
		}
		assertHeadsReturn(state)
	}
	assertExactTimeline := func(state string) {
		t.Helper()
		wait(fmt.Sprintf(`location.pathname + location.search + location.hash === %q &&
			document.documentElement.dataset.ready === '1' && document.querySelector('#stage')?.dataset.view === 'timeline' &&
			document.querySelector('#state-status')?.hidden && document.querySelector('#front')?.dataset.front !== %q &&
			document.querySelector('#tl .node[data-state=%q]')?.checkVisibility({visibilityProperty:true, opacityProperty:true}) &&
			document.querySelector('#heads-link')?.getAttribute('href') === %q`, timelinePath(state), state, state, returnPath(state)))
	}

	// The front face's timeline control opens the current tip on the timeline.
	focus(fmt.Sprintf(`.wh-card[data-head=%q] .wh-map[href=%q]`, lane, timelinePath(tip)))
	enter()
	assertExactTimeline(tip)
	assertHeadsReturn(tip)

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
	assertStackAndReturn(rearState, true)

	// Browse the full lane history and choose a different older state by keyboard.
	focus(fmt.Sprintf(`.wh-card[data-head=%q] .wh-all`, lane))
	enter()
	focus(fmt.Sprintf(`.wh-history[open] a[href=%q]`, stackPath(earlierState)))
	enter()
	assertStackAndReturn(earlierState, true)

	// Follow the timestamp's timeline link and return directly to its exact rear pane.
	focus(fmt.Sprintf(`.wh-card[data-head=%q] .wh-all`, lane))
	enter()
	focus(fmt.Sprintf(`.wh-history[open] .wh-history-time[href=%q]`, timelinePath(earlierState)))
	enter()
	assertExactTimeline(earlierState)
	assertHeadsReturn(earlierState)

	// The header's timeline control inherits the returned position. The exact
	// timeline node opens the stack for that state.
	focus(fmt.Sprintf(`#timeline-link[href=%q]`, timelinePath(earlierState)))
	enter()
	assertExactTimeline(earlierState)
	wait(fmt.Sprintf(`(() => {
		const node = document.querySelector('#tl .node[data-state=%q]');
		if (!node?.checkVisibility({visibilityProperty:true, opacityProperty:true})) return false;
		const box = node.getBoundingClientRect();
		if (!node.contains(document.elementFromPoint(box.x + box.width / 2, box.y + box.height / 2))) return false;
		node.id = 'chrome-exact-timeline-node';
		return true;
	})()`, earlierState))
	b.do(t, nil, "tap", map[string]any{"id": "chrome-exact-timeline-node"})
	assertStackAndReturn(earlierState, false)
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

// A folded head keeps its activity trace when expanded, dragged, deselected,
// and updated by the next model response.
func TestSurfaceHeadActivityMovesWithWindowInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	var responses atomic.Int32
	const last = "2026-11-06T09:32:11Z"
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), "", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/surface/model.json" {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{
				"schema":"world/surface-model/v1", "now":"2026-11-19T15:00:00Z", "nowEdge":{"sessions":[{"episode":"running-episode"}]},
				"lanes":[
					{"key":"running","name":"Running","states":[
						{"id":"running-tip","rank":1,"at":"2026-09-01T10:00:00Z","episode":"running-episode"}
					]},
					{"key":"older","name":"Older","states":[
						{"id":"older-tip","rank":3,"title":"Recorded model %d"},
						{"id":"older-last","rank":2,"at":%q},
						{"id":"older-first","rank":1,"at":"2026-11-02T13:00:00Z"}
					]}
				]
			}`, responses.Add(1), last)
		})
	})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface/heads", "raw": true})
	wait(fmt.Sprintf(`(() => {
		const card = document.querySelector('.wh-card[data-head="older"]');
		window.visibleRecency = card => [...(card?.querySelectorAll('.wh-recency') || [])]
			.find(part => part.getClientRects().length && getComputedStyle(part).visibility === 'visible');
		const running = document.querySelector('.wh-card[data-head="running"]');
		const runningAge = window.visibleRecency(running)?.querySelector('time.wh-age');
		if (!running?.querySelector('.wh-live')?.getClientRects().length ||
			Date.parse(runningAge?.dateTime) !== Date.parse('2026-09-01T10:00:00Z')) return false;
		const recency = window.visibleRecency(card);
		const trace = recency?.querySelector('.wh-trace');
		const stamp = recency?.querySelector('time.wh-age');
		if (!trace || !stamp || Date.parse(stamp.dateTime) !== Date.parse(%q) ||
			!card.querySelector('.wh-restore')?.getClientRects().length) return false;
		window.headTrace = trace.outerHTML;
		window.headStamp = stamp.dateTime;
		window.foldedPosition = [card.style.left, card.style.top];
		card.querySelector('.wh-restore').id = 'chrome-restore-older';
		return true;
	})()`, last))
	b.do(t, nil, "tap", map[string]any{"id": "chrome-restore-older"})
	wait(`(() => {
		const card = document.querySelector('.wh-card[data-head="older"].is-selected');
		const grip = card?.querySelector('.wh-grip');
		const recency = window.visibleRecency(card);
		if (!grip?.getClientRects().length || recency?.querySelector('.wh-trace')?.outerHTML !== window.headTrace ||
			recency?.querySelector('time.wh-age')?.dateTime !== window.headStamp ||
			Math.abs(parseFloat(card.style.left) - parseFloat(window.foldedPosition[0])) > 2 ||
			card.style.top !== window.foldedPosition[1]) return false;
		window.expandedPosition = [card.offsetLeft, card.offsetTop];
		grip.id = 'chrome-drag-older';
		return true;
	})()`)
	b.do(t, nil, "drag", map[string]any{"id": "chrome-drag-older", "dx": 40, "dy": -30})
	wait(`(() => {
		const card = document.querySelector('.wh-card[data-head="older"]');
		const [left, top] = window.expandedPosition;
		if (Math.hypot(card.offsetLeft - left, card.offsetTop - top) < 10) return false;
		if (window.visibleRecency(card)?.querySelector('.wh-trace')?.outerHTML !== window.headTrace) return false;
		window.draggedPosition = [card.style.left, card.style.top];
		const picker = document.querySelector('.wh-picker select');
		picker.value = 'running';
		picker.dispatchEvent(new Event('change', {bubbles:true}));
		return true;
	})()`)
	const retained = `(() => {
		const card = document.querySelector('.wh-card[data-head="older"]:not(.is-selected)');
		const recency = window.visibleRecency(card);
		return recency?.querySelector('.wh-trace')?.outerHTML === window.headTrace &&
			recency?.querySelector('time.wh-age')?.dateTime === window.headStamp &&
			card.style.left === window.draggedPosition[0] && card.style.top === window.draggedPosition[1] &&
			window.__cspv?.length === 0;
	})()`
	wait(retained)
	wait(`(() => {
		const card = document.querySelector('.wh-card[data-head="older"]');
		window.modelTitleBeforeRefresh = card.querySelector('.wh-open').textContent;
		document.dispatchEvent(new Event('visibilitychange'));
		return true;
	})()`)
	wait(retained + ` && document.querySelector('.wh-card[data-head="older"] .wh-open').textContent !== window.modelTitleBeforeRefresh`)
}
