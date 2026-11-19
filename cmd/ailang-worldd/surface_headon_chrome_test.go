package main

import (
	"fmt"
	"net/http"
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
		document.activeElement.matches('#stage a:is([data-head-state], [data-rear-state])') && ` + visibleFocus + ` && ` + fitsViewport)
	wait(`(() => {
		const face = document.activeElement;
		if (!face.matches('#stage a[href]:is([data-head-state], [data-rear-state])')) return false;
		const href = new URL(face.getAttribute('href'), location.href);
		const state = face.dataset.headState || face.dataset.rearState;
		const lane = href.searchParams.get('brane');
		if (href.pathname !== '/surface' || href.searchParams.get('at') !== state || !lane ||
			(face.dataset.headLane && lane !== face.dataset.headLane) || !state) return false;
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
			window.travelTimes = new Map(model.lanes.flatMap(lane => lane.states.map(state => [state.id, Date.parse(state.at)])));
		});
		return Boolean(window.travelScene);
	})()`)
	wait(`window.travelTimes instanceof Map && window.getTravelStop() === '' &&
		window.visibleTravelFaces().length > 0`)

	assertStop := func(slot, newer string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const stop = window.getTravelStop();
			const cutoff = window.travelTimes.get(stop);
			const faces = window.visibleTravelFaces();
			const selected = faces.find(link => link.dataset.headState === stop);
			if (!selected || !Number.isFinite(cutoff) || stop === window[%q] ||
				(window[%q] && cutoff > window.travelTimes.get(window[%q]))) return false;
			const bounds = selected.getBoundingClientRect();
			if (bounds.right <= 0 || bounds.bottom <= 0 || bounds.left >= innerWidth || bounds.top >= innerHeight) return false;
			const url = new URL(selected.href);
			if (url.pathname !== '/surface' || url.searchParams.get('at') !== stop ||
				url.searchParams.get('brane') !== selected.dataset.headLane) return false;
			if (!faces.every(face => window.travelTimes.get(face.dataset.headState) <= cutoff)) return false;
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

func TestSurfaceHeadOnRenderedBuildingsConvergeInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.documentElement.dataset.ready === '1' &&
		[...document.querySelectorAll('#stage[data-view="heads"] .hc-roof polygon')]
			.some(roof => roof.checkVisibility() && roof.points.length === 4)`})
	// Observe the strokes the page actually submits to Chrome on a repaint.
	// Keep this independent of WorldHeadCamera and its projection math.
	b.do(t, nil, "waitFor", map[string]any{"expression": `(() => {
		const canvas = document.querySelector('.hc-graph');
		const ctx = canvas.getContext('2d');
		const originals = {};
		let path = [];
		for (const [method, record] of Object.entries({
			beginPath: () => { path = []; },
			moveTo: (x, y) => { path.push(['M', x, y]); },
			lineTo: (x, y) => { path.push(['L', x, y]); },
			stroke: () => {
				if (path.map(point => point[0]).join('') !== 'MLML') return;
				const [[, x, y], [, left, endY], [, x2, y2], [, right, rightY]] = path;
				if (x !== x2 || y !== y2 || left >= x || right <= x || endY <= y || endY !== rightY) return;
				const pixel = ctx.getTransform().transformPoint({x, y});
				const rect = canvas.getBoundingClientRect();
				window.drawnRoadVanishingPoint = [
					rect.x + pixel.x * rect.width / canvas.width, rect.y + pixel.y * rect.height / canvas.height
				];
				for (const key of Object.keys(originals)) delete ctx[key];
			},
		})) {
			originals[method] = ctx[method];
			ctx[method] = (...args) => {
				record(...args);
				return originals[method].apply(ctx, args);
			};
		}
		document.querySelector('.hc-scroll').dispatchEvent(new Event('scroll'));
		return true;
	})()`})
	var result struct {
		Left, Right, Center    int
		RoadOffset, FaceOffset float64
		Failures               []string
	}
	b.do(t, &result, "waitFor", map[string]any{"expression": `(() => {
		const vp = window.drawnRoadVanishingPoint;
		if (!vp) return false;
		const viewport = document.querySelector('.hc-scroll');
		const center = viewport.getBoundingClientRect().left + viewport.clientWidth / 2;
		const result = {left:0, right:0, center:0, roadOffset:Math.abs(vp[0] - center), faceOffset:10000, failures:[]};
		const vertices = polygon => {
			const transform = polygon.getScreenCTM();
			return [...polygon.points].map(point => {
				const {x, y} = point.matrixTransform(transform);
				return [x, y];
			});
		};
		const distance = (a, b) => Math.hypot(b[0] - a[0], b[1] - a[1]);
		const nearest = (target, points) => points.reduce((best, pt) => distance(pt, target) < distance(best, target) ? pt : best);
		const miss = (a, b) => Math.abs((b[0] - a[0]) * (vp[1] - a[1]) - (b[1] - a[1]) * (vp[0] - a[0])) / distance(a, b);
		for (const face of document.querySelectorAll('.hc-face[data-head-state]')) {
			const {left, right, top, bottom} = face.getBoundingClientRect();
			result.faceOffset = Math.min(result.faceOffset, Math.abs((left + right) / 2 - center));
			const roof = face.querySelector('.hc-roof polygon'), cheek = face.querySelector('.hc-cheek polygon');
			if (!roof?.checkVisibility() || roof.points.length !== 4) continue;
			const corners = vertices(roof);
			const front = [[left, top], [right, top]].map(target => nearest(target, corners));
			const back = corners.filter(pt => !front.includes(pt)).sort((a, b) => a[0] - b[0]);
			if (back.length !== 2 || distance(front[0], [left, top]) > 1 || distance(front[1], [right, top]) > 1 ||
				distance(back[0], back[1]) >= distance(front[0], front[1]) ||
				miss(front[0], back[0]) > 1 || miss(front[1], back[1]) > 1) {
				result.failures.push(face.dataset.headState);
				continue;
			}
			const bank = right < vp[0] - 1 ? 'left' : left > vp[0] + 1 ? 'right' : '';
			if (!bank) {
				if (left < vp[0] - 1 && right > vp[0] + 1 && cheek?.checkVisibility()) result.failures.push(face.dataset.headState);
				else result.center++;
				continue;
			}
			if (!cheek?.checkVisibility() || cheek.points.length !== 4) {
				result.failures.push(face.dataset.headState);
				continue;
			}
			const walls = vertices(cheek);
			const inward = bank === 'left' ? 1 : 0;
			const foot = [front[inward][0], bottom];
			const frontFoot = nearest(foot, walls);
			const backFoot = walls.find(pt => pt !== frontFoot && distance(pt, front[inward]) > 1 && distance(pt, back[inward]) > 1);
			if (distance(frontFoot, foot) > 1 || !backFoot || miss(frontFoot, backFoot) > 1 ||
				distance(nearest(front[inward], walls), front[inward]) > 1 ||
				distance(nearest(back[inward], walls), back[inward]) > 1) {
				result.failures.push(face.dataset.headState);
			} else {
				result[bank]++;
			}
		}
		return result;
	})()`})
	if result.Left == 0 || result.Right == 0 || result.RoadOffset > 1 || result.FaceOffset > 1 || len(result.Failures) != 0 {
		t.Fatalf("expected centered road and converging rendered roofs and inward walls: %+v", result)
	}
}

func TestSurfaceHeadOnOpensOlderRecordedPaneInChrome(t *testing.T) {
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
		window.longLaneScene = document.querySelector('.hc-scene');
		fetch('/surface/model.json').then(r => r.json()).then(model => {
			window.recordedLanes = model.lanes;
		});
		return Boolean(window.longLaneScene);
	})()`)
	var inventory struct {
		Lanes                     int
		Missing, Duplicate, Wrong []string
	}
	b.do(t, &inventory, "waitFor", map[string]any{"expression": `(() => {
		const lanes = window.recordedLanes;
		if (!lanes) return false;
		const byState = new Map();
		for (const link of window.longLaneScene.querySelectorAll('a[data-head-state], a[data-rear-state]')) {
			const id = link.dataset.headState || link.dataset.rearState;
			const copies = byState.get(id) || [];
			copies.push(link);
			byState.set(id, copies);
		}
		const report = {lanes:lanes.filter(lane=>lane.states.length).length,missing:[],duplicate:[],wrong:[]};
		for (const lane of lanes) {
			for (const state of lane.states) {
				const copies = byState.get(state.id) || [];
				if (copies.length === 0) { report.missing.push(state.id); continue; }
				if (copies.length > 1) { report.duplicate.push(state.id); continue; }
				const link = copies[0];
				const url = new URL(link.getAttribute('href'), location.href);
				const rect = link.getBoundingClientRect();
				if (url.pathname !== '/surface' || url.searchParams.get('at') !== state.id ||
					url.searchParams.get('brane') !== lane.key || rect.width <= 0 || rect.height <= 0) report.wrong.push(state.id);
			}
		}
		return report;
	})()`})
	if inventory.Lanes < 2 || len(inventory.Missing)+len(inventory.Duplicate)+len(inventory.Wrong) > 0 {
		t.Fatalf("expected one rendered native pane with the exact stack address per recorded state: %+v", inventory)
	}
	wait(`(() => {
		const states = window.recordedLanes.find(lane => lane.key === 'row-156').states
			.slice().sort((a, b) => b.rank - a.rank);
		if (states.length < 5) return false;
		const oldest = states.at(-1);
		const pane = [...window.longLaneScene.querySelectorAll('a[data-rear-state]')]
			.find(link => link.dataset.rearState === oldest.id);
		if (!pane) return false;
		pane.focus();
		if (document.activeElement !== pane) return false;
		window.oldNativePane = oldest.id;
		return true;
	})()`)
	wait(`window.longLaneScene.dataset.inspectState === window.oldNativePane &&
		location.hash === '#heads' && !location.search`)
	b.do(t, nil, "key", map[string]any{"key": "Enter"})
	wait(`document.querySelector('#stage[data-view="brane"]') &&
		document.querySelector('#front')?.dataset.front === window.oldNativePane`)
	b.do(t, nil, "tap", map[string]any{"id": "head-angle"})
	wait(`document.querySelector('#stage[data-view="heads"] .hc-scene') === window.longLaneScene &&
		window.longLaneScene.querySelector('a[data-rear-state][aria-current="true"]')?.dataset.rearState === window.oldNativePane`)
}

func TestSurfaceHeadOnUsesSharedRecordedTimeInChrome(t *testing.T) {
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
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	wait(`document.documentElement.dataset.ready === '1' && document.querySelector('#stage')?.dataset.view === 'timeline'`)
	var timeline struct {
		Tips   int
		DeltaX float64
		Ticks  int
	}
	b.do(t, &timeline, "waitFor", map[string]any{"expression": `(() => {
		const tips = [...document.querySelectorAll('#tl .node:is([data-state="a-tip"],[data-state="b-tip"])')];
		const ticks = [...document.querySelectorAll('#tl .tick-t')]
			.filter(tick => tick.checkVisibility() && tick.textContent.trim());
		const center = el => { const r = el.getBoundingClientRect(); return r.x + r.width / 2; };
		return {tips:tips.length, deltaX:tips.length === 2 ? Math.abs(center(tips[0])-center(tips[1])) : -1, ticks:ticks.length};
	})()`})
	if timeline.Tips != 2 || timeline.DeltaX > 1 || timeline.Ticks == 0 {
		t.Fatalf("expected equal-time Timeline nodes and a visible time axis: %+v", timeline)
	}
	b.do(t, nil, "tap", map[string]any{"id": "head-angle"})
	const setupScene = `(() => {
		const laneX = el => { const r = el.getBoundingClientRect(); return r.left + r.width / 2; };
		const stage = document.querySelector('#stage[data-view="heads"]');
		if (!stage) return false;
		const axis = stage.querySelector('.hc-time-axis');
		const ticks = axis && [...axis.querySelectorAll('time[datetime]')];
		if (!axis?.checkVisibility() || !axis.textContent.includes('Time') ||
			!ticks.some(tick => tick.checkVisibility() && tick.textContent.trim() && Number.isFinite(Date.parse(tick.dateTime)))) return false;
		window.timedScene = stage.querySelector('.hc-scene');
		window.paneFor = id => [...window.timedScene.querySelectorAll('a[data-head-state], a[data-rear-state]')]
			.find(el => (el.dataset.headState || el.dataset.rearState) === id);
		window.alignedPanes = (a, b) => {
			const first = window.paneFor(a), second = window.paneFor(b);
			if (!first || !second) return false;
			const midpoint = r => r.top + r.height / 2;
			const front = link => link.querySelector('.hc-mini-front') || link;
			return Math.abs(midpoint(front(first).getBoundingClientRect()) - midpoint(front(second).getBoundingClientRect())) <= 1;
		};
		window.lateralStep = (id, key) => {
			const link = window.paneFor(id);
			link.focus();
			if (document.activeElement !== link) return false;
			link.dispatchEvent(new KeyboardEvent('keydown', {key, bubbles:true, cancelable:true}));
			return true;
		};
		window.centeredAt = (id, lane) => {
			const face = window.paneFor(id);
			const scroll = window.timedScene.querySelector('.hc-scroll');
			const center = scroll.getBoundingClientRect().left + scroll.clientWidth / 2;
			return face?.matches('a[data-head-state][aria-current="true"]') &&
				window.timedScene.dataset.viewLane === lane && document.activeElement === face &&
				Math.abs(laneX(face) - center) <= 1;
		};
		window.axisEndsAt = instant => {
			const ticks = [...window.timedScene.querySelectorAll('.hc-time-axis time[datetime]')]
				.filter(tick => tick.checkVisibility());
			return ticks.at(-1)?.dateTime === instant && ticks.every(tick => Date.parse(tick.dateTime) <= Date.parse(instant));
		};
		return window.alignedPanes('a-tip', 'b-tip');
	})()`
	wait(setupScene)
	wait(`(() => {
		const tip = window.paneFor('a-tip');
		tip.focus();
		if (document.activeElement !== tip) return false;
		tip.dispatchEvent(new KeyboardEvent('keydown', {key:'+', bubbles:true, cancelable:true}));
		return true;
	})()`)
	wait(`(() => {
		const scene = window.timedScene;
		if (scene.dataset.travelState !== 'a-stop' || document.querySelector('#stage').dataset.view !== 'heads') return false;
		window.sameRecordedCut = () => {
			const expected = ['a-stop', 'b-earlier', 'b-tie', 'shared-cut', 'shared-early'];
			const actual = [...scene.querySelectorAll('a[data-head-state], a[data-rear-state]')]
				.map(el => el.dataset.headState || el.dataset.rearState).sort();
			return JSON.stringify(actual) === JSON.stringify(expected);
		};
		return window.sameRecordedCut() && window.alignedPanes('a-stop', 'b-tie');
	})()`)
	var movement struct {
		Key   string
		Stops []struct{ ID, Lane string }
	}
	b.do(t, &movement, "waitFor", map[string]any{"expression": `(() => {
		const laneX = el => { const r = el.getBoundingClientRect(); return r.left + r.width / 2; };
		const ordered = [...window.timedScene.querySelectorAll('a[data-head-state]')]
			.sort((a, b) => laneX(a) - laneX(b));
		const first = ordered.indexOf(window.paneFor('a-stop')), second = ordered.indexOf(window.paneFor('b-tie'));
		if (first < 0 || second < 0 || first === second) return false;
		window.historicalCenter = (id, lane) => window.centeredAt(id, lane) &&
			window.timedScene.dataset.travelState === 'a-stop' && window.sameRecordedCut() &&
			window.alignedPanes('a-stop', 'b-tie') && window.axisEndsAt('2026-10-20T12:00:00.000Z');
		const stops = ordered.slice(Math.min(first, second), Math.max(first, second) + 1);
		if (first > second) stops.reverse();
		return {key:first < second ? 'ArrowRight' : 'ArrowLeft',
			stops:stops.map(el=>({id:el.dataset.headState, lane:el.dataset.headLane}))};
	})()`})
	for i := 1; i < len(movement.Stops); i++ {
		from, to := movement.Stops[i-1], movement.Stops[i]
		wait(fmt.Sprintf("window.lateralStep(%q, %q)", from.ID, movement.Key))
		wait(fmt.Sprintf("window.historicalCenter(%q, %q)", to.ID, to.Lane))
	}
	reverseKey := "ArrowLeft"
	if movement.Key == reverseKey {
		reverseKey = "ArrowRight"
	}
	for i := len(movement.Stops) - 1; i > 0; i-- {
		from, to := movement.Stops[i], movement.Stops[i-1]
		wait(fmt.Sprintf("window.lateralStep(%q, %q)", from.ID, reverseKey))
		wait(fmt.Sprintf("window.historicalCenter(%q, %q)", to.ID, to.Lane))
	}

	// Measure the near plane at the historical cut, then launch a fresh
	// Head-on view from the present with no explicit cutoff.
	var nearPlane struct {
		Top, Width, Height float64
	}
	b.do(t, &nearPlane, "waitFor", map[string]any{"expression": `(() => {
		const {top, width, height} = window.paneFor('a-stop').getBoundingClientRect();
		return {top, width, height};
	})()`})
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait(`document.documentElement.dataset.ready === '1'`)
	wait(setupScene)
	var neighbor struct {
		ID, Lane, Out, Back string
	}
	b.do(t, &neighbor, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => {
		const scene = window.timedScene;
		if (scene.dataset.travelState || scene.dataset.viewLane !== 'shared' ||
			!scene.querySelector('.hc-time-axis').checkVisibility()) return false;
		const heads = [...scene.querySelectorAll('a[data-head-state]')].sort((a, b) =>
			a.getBoundingClientRect().left - b.getBoundingClientRect().left);
		const start = heads.indexOf(window.paneFor('shared-late'));
		if (start < 0) return false;
		const side = start > 0 ? -1 : 1;
		const face = heads[start + side];
		if (!face || !['a-tip', 'b-tip'].includes(face.dataset.headState)) return false;
		window.presentFaceHeight = face.getBoundingClientRect().height;
		window.nearPlane = {top:%f, width:%f, height:%f};
		window.atNearPlane = id => {
			const face = window.paneFor(id);
			if (!face || !face.checkVisibility()) return false;
			const rect = face.getBoundingClientRect();
			return ['top', 'width', 'height'].every(edge => Math.abs(rect[edge] - window.nearPlane[edge]) <= 1);
		};
		return {id:face.dataset.headState, lane:face.dataset.headLane,
			out:side < 0 ? 'ArrowLeft' : 'ArrowRight', back:side < 0 ? 'ArrowRight' : 'ArrowLeft'};
	})()`, nearPlane.Top, nearPlane.Width, nearPlane.Height)})
	wait(fmt.Sprintf("window.lateralStep('shared-late', %q)", neighbor.Out))
	wait(fmt.Sprintf(`window.centeredAt(%q, %q) && window.atNearPlane(%q) &&
		window.paneFor(%q).getBoundingClientRect().height > window.presentFaceHeight + 1 &&
		window.axisEndsAt('2026-10-20T14:00:00.000Z') && !window.paneFor('shared-late') &&
		window.alignedPanes('a-tip', 'b-tip')`, neighbor.ID, neighbor.Lane, neighbor.ID, neighbor.ID))

	// A resize or model refresh can settle the flight before its last frame.
	// The destination must still receive keyboard focus.
	wait(fmt.Sprintf("window.lateralStep(%q, %q) && window.dispatchEvent(new Event('resize'))", neighbor.ID, neighbor.Back))
	wait(`window.centeredAt('shared-late', 'shared') && window.atNearPlane('shared-late') &&
		window.axisEndsAt('2026-10-20T14:30:00.000Z') &&
		window.paneFor('a-tip').getBoundingClientRect().height < window.nearPlane.height - 1 &&
		window.alignedPanes('a-tip', 'b-tip')`)
}

// The stop shares an instant with b-tie despite different wire offsets.
// b-earlier and b-later have ranks that deliberately disagree with their time.
const recordedTimeBrowserModel = `{
	"schema": "world/surface-model/v1",
	"now": "2026-10-20T15:00:00Z",
	"nowEdge": {},
	"lanes": [
		{"key":"shared","name":"Shared world","glyph":"circle","colour":"#8496a5","states":[
			{"id":"shared-early","kind":"shared","label":"Shared early","at":"2026-10-20T09:00:00Z","rank":1},
			{"id":"shared-cut","kind":"shared","label":"Shared cut","at":"2026-10-20T12:00:00Z","rank":2},
			{"id":"shared-late","kind":"shared","label":"Shared late","at":"2026-10-20T14:30:00Z","rank":3}
		]},
		{"key":"row-156","row":156,"name":"First branch","glyph":"hex","colour":"#93675d","states":[
			{"id":"a-stop","kind":"step","label":"Historical stop","at":"2026-10-20T12:00:00Z","rank":20},
			{"id":"a-tip","kind":"step","label":"First tip","at":"2026-10-20T14:00:00Z","rank":50}
		]},
		{"key":"row-189","row":189,"name":"Second branch","glyph":"square","colour":"#42669a","states":[
			{"id":"b-earlier","kind":"step","label":"Earlier time, higher rank","at":"2026-10-20T10:00:00Z","rank":30},
			{"id":"b-tie","kind":"step","label":"Same instant","at":"2026-10-20T14:00:00+02:00","rank":40},
			{"id":"b-later","kind":"step","label":"Later time, lower rank","at":"2026-10-20T13:00:00Z","rank":10},
			{"id":"b-tip","kind":"step","label":"Second tip","at":"2026-10-20T16:00:00+02:00","rank":51}
		]}
	],
	"links": [
		{"from":"a-stop","to":"a-tip"},
		{"from":"b-earlier","to":"b-tie"},
		{"from":"b-tie","to":"b-later"},
		{"from":"b-later","to":"b-tip"}
	]
}`
