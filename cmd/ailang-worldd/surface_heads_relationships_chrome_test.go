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
		{From: "mixed-base", To: "mixed-tip"},
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
		window.station = () => document.querySelector('.wh-topology-station');
		const on = key => headCard(key)?.classList.contains('has-wh-link');
		return headCard('selected')?.classList.contains('is-selected') &&
			on('base') && on('derived') && on('mixed') && !on('unrelated') &&
			headCard('unrelated').classList.contains('wh-unrelated') &&
			station()?.querySelector('.wh-topology-in [data-relation="fork"] button[data-locate="base"]') &&
			station()?.querySelector('.wh-topology-out [data-relation="base"] button[data-locate="derived"]') &&
			station()?.querySelector('.wh-topology-out [data-relation="merge"] button[data-locate="mixed"]') &&
			station()?.querySelector('.wh-topology-in button[data-locate="mixed"]') &&
			document.querySelectorAll('.wh-link.is-active').length === 3 &&
			document.querySelector('.wh-neighborhood header')?.textContent.includes('1 off page');
	})()`)
	tap(`[data-topology-state="selected-tip"]`)
	wait(`(() => {
		const evidence = station().querySelector('.wh-topology-evidence:not([hidden])');
		return evidence?.querySelector('a[data-head-timeline="mixed-base"]')?.getAttribute('href') === '/surface#brane=mixed&at=mixed-base' &&
			evidence.querySelector('a[data-head-state="mixed-base"]')?.getAttribute('href') === '/surface?brane=mixed&at=mixed-base#heads' &&
			evidence.textContent.includes('Earlier state');
	})()`)
	tap(`[data-topology-state="selected-tip"]`)

	// Following a folded peer moves the focus and redraws the neighborhood.
	tap(`.wh-card[data-head="mixed"] .wh-restore`)
	wait(`headCard('mixed').classList.contains('is-selected') &&
		headCard('selected').classList.contains('has-wh-link') &&
		!headCard('base').classList.contains('has-wh-link') &&
		!headCard('derived').classList.contains('has-wh-link') &&
		station().querySelector('[data-topology-state="mixed-tip"]') &&
		station().querySelector('.wh-topology-in [data-relation="continues"] a[data-head-timeline="mixed-base"]')?.getAttribute('href') === '/surface#brane=mixed&at=mixed-base' &&
		document.querySelectorAll('.wh-link.is-active').length === 1`)
	tap(`.wh-topology-earlier > summary`)
	wait(`document.querySelector('.wh-topology-earlier[open] [data-topology-state="mixed-base"]')?.getClientRects().length`)
	tap(`.wh-neighborhood > .wh-topology-panel > .wh-topology button[data-locate="selected"]`)
	wait(`headCard('selected').classList.contains('is-selected') && headCard('base').classList.contains('has-wh-link')`)

	// Navigate via the topology to a peer on the second page.
	tap(`.wh-neighborhood button[data-locate="remote"]`)
	wait(`headCard('remote')?.classList.contains('is-selected') &&
		!headCard('selected') && !document.querySelector('.wh-card.has-wh-link') &&
		document.querySelector('.wh-neighborhood header')?.textContent.includes('1 off page') &&
		document.querySelector('.wh-neighborhood button[data-locate="selected"]')`)

	tap(`.wh-card[data-head="isolated"] .wh-restore`)
	wait(`headCard('isolated').classList.contains('is-selected') &&
		!document.querySelector('.wh-card.has-wh-link, .wh-card.wh-unrelated, .wh-link.is-active') &&
		document.querySelector('.wh-neighborhood header')?.textContent.includes('No branch links') &&
		station().querySelector('[data-topology-state="isolated-tip"]') &&
		window.__cspv?.length === 0`)
}
