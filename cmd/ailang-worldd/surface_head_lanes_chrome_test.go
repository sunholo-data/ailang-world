package main

import (
	"path/filepath"
	"testing"
)

func TestSurfaceHeadLaneStripInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := rigSurface(t, filepath.Join(t.TempDir(), "credentials.db"))
	b.do(t, nil, "tab", nil)
	defer b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface#heads", "raw": true})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}

	// Make the strip narrow enough that at least one lane starts out clipped.
	wait(`(() => {
		const strip = document.querySelector('.hc-lanes:not([hidden])');
		const current = strip?.querySelector('[aria-current="true"]');
		if (!current || strip.children.length < 3) return false;
		strip.style.width = '240px';
		const bounds = strip.getBoundingClientRect();
		const stage = strip.closest('#stage');
		const lanes = [...strip.querySelectorAll('button')];
		const target = lanes.find(button => button !== current &&
			button.getBoundingClientRect().right > bounds.right);
		if (!target || strip.scrollWidth <= strip.clientWidth || bounds.top !== stage.getBoundingClientRect().top ||
			!lanes.every(button => button.getBoundingClientRect().top === lanes[0].getBoundingClientRect().top)) return false;
		window.laneStrip = {strip, target, scene: strip.closest('.hc-scene')};
		target.focus({preventScroll:true});
		return document.activeElement === target;
	})()`)
	b.do(t, nil, "key", map[string]any{"key": "Enter"})
	wait(`(() => {
		const {strip, target, scene} = window.laneStrip;
		const bounds = strip.getBoundingClientRect(), button = target.getBoundingClientRect();
		return scene === document.querySelector('.hc-scene') && scene.dataset.viewLane === target.dataset.cameraLane &&
			strip.querySelector('[aria-current="true"]') === target && button.left >= bounds.left && button.right <= bounds.right + 1;
	})()`)
	b.do(t, nil, "tap", map[string]any{"id": "timeline-angle"})
	wait(`document.querySelector('#stage')?.dataset.view === 'timeline' && window.laneStrip.strip.hidden`)
}
