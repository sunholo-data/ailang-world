package main

// Row 201 A3 sprint 3 M8d and sprint 4 (design §2.8, §6 XA3): Needs you and the decision list in the
// surface app, in a live Chrome.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// rigSurface serves the daemon over credentials store creds with a copy of the rig's surface index.
func rigSurface(t *testing.T, creds string) string {
	t.Helper()
	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(index, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	return cardOrigin(t, creds, index, nil)
}

// sectionText waits until the now card's section sel (data-needs, data-decided) has text starting
// with prefix, and returns it.
func sectionText(t *testing.T, b *passkeyBrowser, sel, prefix string) string {
	t.Helper()
	var s string
	b.do(t, &s, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const e = document.querySelector('[%s]'); return e && e.textContent.startsWith(%q) && e.textContent; })()`, sel, prefix)})
	return s
}

// pollSurface runs the app's own refresh (window.World.poll) and waits until it has finished.
func pollSurface(t *testing.T, b *passkeyBrowser) {
	t.Helper()
	b.do(t, nil, "waitFor", map[string]any{"expression": `(window.__polled = 0, window.World.poll().then(() => { window.__polled = 1; }), 1)`})
	b.do(t, nil, "waitFor", map[string]any{"expression": `window.__polled === 1`})
}

// keyView is the driver's key answer.
type keyView struct{ Location, Active, View string }

// XA3: /surface's now card shows the pending request under "Needs you · 1", from GET
// /surface/requests.json; a second request appears on the app's next refresh (no timer of its own).
// N focuses the first one's link, from the timeline and from a brane, and Enter opens its card.
// Approved there and minted by approve --await, the request is under Decided as "You · Minted · ep1 ·
// webauthn:mark-mac · <hh:mm>"; with four decided, Decided shows three and "Show all 4" toggles.
func TestNeedsYouListsPendingRequests(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	origin := rigSurface(t, a.credPath)
	virtualPasskey(t, b, a.credPath, origin, "mark-mac", 7, now-100)
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	if needs, want := sectionText(t, b, "data-needs", "Needs you · 1"), "Needs you · 1ep1 · claude · 10 grants · expires in 10 minNotify me"; needs != want {
		t.Fatalf("Needs you shows %q, want %q", needs, want)
	}
	var k keyView
	if b.do(t, &k, "key", map[string]any{"key": "n"}); k.Active != "needs-"+a.id || k.View != "timeline" {
		t.Fatalf("N on the timeline: %+v", k)
	}
	var found bool
	if b.do(t, &found, "waitFor", map[string]any{"expression": `window.World.open('e2')`}); !found {
		t.Fatal("the rig has no episode e2")
	}
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('stage').dataset.view === 'brane'`})
	if b.do(t, &k, "key", map[string]any{"key": "N"}); k.Active != "needs-"+a.id || k.View != "timeline" {
		t.Fatalf("N in a brane: %+v", k)
	}
	second := newApproveRigIn(t, a.r, "ep2", now)
	pollSurface(t, b)
	if needs := sectionText(t, b, "data-needs", "Needs you · 2"); !strings.Contains(needs, "ep1 · claude") || !strings.Contains(needs, "ep2 · claude") {
		t.Fatalf("after a refresh Needs you shows %q", needs)
	}
	var href string
	b.do(t, &href, "waitFor", map[string]any{"expression": `document.getElementById('needs-` + a.id + `').getAttribute('href')`})
	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('needs-` + a.id + `').focus(), 1)`})
	if b.do(t, &k, "key", map[string]any{"key": "Enter", "page": true}); href != "/surface/approve/"+a.id || k.Location != origin+"/surface/approve/"+a.id {
		t.Fatalf("the link %q opened %q", href, k.Location)
	}
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
	var view cardView
	if b.do(t, &view, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}}); view.State != "recorded" ||
		!strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1") || len(view.Errors) != 0 || len(view.Violations) != 0 {
		t.Fatalf("the card shows %+v", view)
	}
	var opened bool
	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
	}
	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	at := time.Unix(a.row(t).ConsumedAt, 0).UTC().Format("15:04")
	if d, want := sectionText(t, b, "data-decided", "Decided"), "DecidedYou · Minted · ep1 · webauthn:mark-mac · "+at; d != want {
		t.Fatalf("Decided shows %q, want %q", d, want)
	}
	if needs := sectionText(t, b, "data-needs", "Needs you · 1"); !strings.HasPrefix(needs, "Needs you · 1ep2 · claude") {
		t.Fatalf("after the mint Needs you shows %q", needs)
	}
	raw := rawCredentials(t, a.credPath)
	defer func() { _ = raw.Close() }()
	for i, r := range []approveRig{second, newApproveRigIn(t, a.r, "ep3", now), newApproveRigIn(t, a.r, "ep4", now)} {
		if _, err := raw.Exec(`UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?`, now-int64(600*(i+1)), r.id); err != nil {
			t.Fatal(err)
		}
	}
	pollSurface(t, b)
	lines := func(want int, button string) {
		t.Helper()
		var got struct {
			Lines  int
			Button string
		}
		b.do(t, &got, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const d = document.querySelector('[data-decided]'), x = document.getElementById('decided-all');
			return d && x && d.querySelectorAll('.nl.sub').length === %d && { lines: %d, button: x.textContent }; })()`, want, want)})
		if got.Lines != want || got.Button != button {
			t.Fatalf("Decided shows %d lines and %q, want %d and %q", got.Lines, got.Button, want, button)
		}
	}
	lines(3, "Show all 4")
	b.do(t, nil, "tap", map[string]any{"id": "decided-all"})
	lines(4, "Show fewer")
	b.do(t, nil, "tap", map[string]any{"id": "decided-all"})
	lines(3, "Show all 4")
	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
}
