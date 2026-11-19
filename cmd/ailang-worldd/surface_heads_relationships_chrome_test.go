package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceHeadsTracksSelectedRelationshipsInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	model := surfacepage.Model{
		Schema: surfacepage.ModelSchema,
		Now:    "2026-11-19T15:00:00Z",
	}
	add := func(key, name, kind, at string) {
		model.Lanes = append(model.Lanes, surfacepage.ModelLane{
			Key: key, Name: name, Glyph: "square", Colour: "#597291",
			States: []surfacepage.ModelState{{ID: key + "-tip", Kind: kind, Rank: len(model.Lanes), At: at, Title: name + " state"}},
		})
	}
	add("selected", "Selected", "step", "2026-11-19T14:00:00Z")
	add("base", "Base", "shared", "2026-11-04T12:00:00Z")
	add("derived", "Derived", "step", "2026-11-03T12:00:00Z")
	add("mixed", "Mixed", "shared", "2026-11-02T12:00:00Z")
	add("unrelated", "Unrelated", "step", "2026-11-01T12:00:00Z")
	for i := 0; i < 11; i++ {
		add(fmt.Sprintf("filler-%d", i), fmt.Sprintf("Filler %d", i), "step", fmt.Sprintf("2026-10-%02dT12:00:00Z", i+1))
	}
	add("remote", "Remote", "step", "2026-09-02T12:00:00Z")
	add("isolated", "Isolated", "step", "2026-09-01T12:00:00Z")
	model.Lanes[3].States = append(model.Lanes[3].States,
		surfacepage.ModelState{ID: "mixed-base", Kind: "shared", Rank: 0, At: "2026-10-30T12:00:00Z", Title: "Earlier shared state"})
	model.Links = []surfacepage.ModelLink{
		{From: "base-tip", To: "selected-tip"},
		{From: "selected-tip", To: "derived-tip"},
		{From: "mixed-base", To: "selected-tip"},
		{From: "selected-tip", To: "mixed-tip"},
		{From: "selected-tip", To: "remote-tip"},
	}
	origin := cardOrigin(t, filepath.Join(t.TempDir(), "credentials.db"), "", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/surface/model.json" {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model)
		})
	})
	wait := func(expression string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": expression})
	}
	tap := func(selector string) {
		t.Helper()
		wait(fmt.Sprintf(`(() => {
			const control = document.querySelector(%q);
			if (!control?.getClientRects().length) return false;
			control.scrollIntoView({block:'center'});
			document.getElementById('chrome-head-relationship-control')?.removeAttribute('id');
			control.id = 'chrome-head-relationship-control';
			return true;
		})()`, selector))
		b.do(t, nil, "tap", map[string]any{"id": "chrome-head-relationship-control"})
	}
	b.do(t, nil, "tab", nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface/heads#selected", "raw": true})
	wait(`(() => {
		window.headCard = key => document.querySelector('.wh-card[data-head="' + key + '"]');
		window.relationshipBadge = key => Array.from(headCard(key)?.querySelectorAll('.wh-link-marker') || [])
			.find(badge => badge.getClientRects().length)?.textContent;
		const base = headCard('base')?.querySelector('.wh-restore');
		const badge = base?.querySelector('.wh-link-marker');
		const recency = base?.querySelector('.wh-recency');
		const description = document.getElementById(base?.getAttribute('aria-describedby'));
		return headCard('selected')?.classList.contains('is-selected') &&
			relationshipBadge('base') === 'In' && relationshipBadge('derived') === 'Out' &&
			relationshipBadge('mixed') === 'Both' && !relationshipBadge('unrelated') &&
			headCard('unrelated').classList.contains('wh-unrelated') &&
			headCard('base').classList.contains('has-wh-link') &&
			description?.textContent === 'Recorded links to the selected branch' &&
			badge.getBoundingClientRect().bottom <= recency.getBoundingClientRect().top &&
			document.querySelectorAll('.wh-link.is-active').length === 3 &&
			document.querySelector('.wh-neighborhood > summary')?.textContent.includes('1 off page');
	})()`)

	// Following a folded peer moves the focus and redraws the neighborhood.
	tap(`.wh-card[data-head="mixed"] .wh-restore`)
	wait(`headCard('mixed').classList.contains('is-selected') &&
		relationshipBadge('selected') === 'Both' &&
		!headCard('base').querySelector('.wh-link-marker') && !relationshipBadge('derived') &&
		document.querySelectorAll('.wh-link.is-active').length === 1`)
	tap(`.wh-neighborhood > summary`)
	tap(`.wh-neighborhood button[data-locate="selected"]`)
	wait(`headCard('selected').classList.contains('is-selected') && relationshipBadge('base') === 'In'`)

	// Navigate via the relationship list to a peer on the second page.
	tap(`.wh-neighborhood > summary`)
	tap(`.wh-neighborhood button[data-locate="remote"]`)
	wait(`headCard('remote')?.classList.contains('is-selected') &&
		!headCard('selected') && !document.querySelector('.wh-link-marker') &&
		document.querySelector('.wh-neighborhood > summary')?.textContent.includes('1 off page') &&
		document.querySelector('.wh-neighborhood button[data-locate="selected"]')`)

	tap(`.wh-card[data-head="isolated"] .wh-restore`)
	wait(`headCard('isolated').classList.contains('is-selected') &&
		!document.querySelector('.wh-link-marker, .wh-card.wh-unrelated, .wh-link.is-active') &&
		document.querySelector('.wh-link-note')?.textContent.includes('No recorded links') &&
		window.__cspv?.length === 0`)
}
