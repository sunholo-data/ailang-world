package main

// Row 201 A3 sprint 3 M8d (design §2.8, §6 XA3's list half): Needs you in the surface app, in a
// live Chrome.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// XA3, the list half (Decided is sprint 4's, plan Appendix C): /surface's now card shows the pending
// request under "Needs you · 1", from GET /surface/requests.json; a second request appears on the
// app's next refresh (no timer of its own); the first one's link opens its card, ready.
func TestNeedsYouListsPendingRequests(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", nil)
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(index, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	origin := cardOrigin(t, a.credPath, index, nil)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	var needs string
	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent; })()`})
	if want := "Needs you · 1ep1 · claude · 10 grants · expires in 10 minNotify me"; needs != want {
		t.Fatalf("Needs you shows %q, want %q", needs, want)
	}
	newApproveRigIn(t, a.r, "ep2", now)
	b.do(t, nil, "waitFor", map[string]any{"expression": `(window.World.poll(), 1)`})
	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent.startsWith('Needs you · 2') && e.textContent; })()`})
	if !strings.Contains(needs, "ep1 · claude") || !strings.Contains(needs, "ep2 · claude") {
		t.Fatalf("after a refresh Needs you shows %q", needs)
	}
	var href string
	b.do(t, &href, "waitFor", map[string]any{"expression": `document.getElementById('needs-` + a.id + `').getAttribute('href')`})
	var at struct{ Location string }
	b.do(t, &at, "press", map[string]any{"id": "needs-" + a.id, "page": true})
	if href != "/surface/approve/"+a.id || at.Location != origin+"/surface/approve/"+a.id {
		t.Fatalf("the link %q opened %q", href, at.Location)
	}
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
	var view cardView
	if b.do(t, &view, "card", nil); !strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1") || len(view.Errors) != 0 || len(view.Violations) != 0 {
		t.Fatalf("the card shows %+v", view)
	}
}
