package main

// Row 243 V1 (design §6 XV1): the app's address bar runs episode report in a live Chrome, on the
// passkey pages' driver (testdata/passkey_cdp.mjs) and skip rule; the daemon is a real one on a test
// store seeded by VR3's seedReportCalls, served on the loopback.

import (
	"fmt"
	"io"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// XV1: `episode report w1` typed into #addr shows the CLI's own bytes in the verb panel and clears the
// bar; `why` and `log` are refused in the panel with their CLI lines (V2b A3); `world a b c`, a refusal
// with no CLI line, is said in the header, the bar and the place kept.
func TestVerbBarEpisodeReportInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{})
	db := filepath.Join(t.TempDir(), "world.db")
	seedReportCalls(t, db, [][5]string{{"w1", "1", "ailang-read", broker.EffectWorkspaceRead, "ok"}, {"w1", "2", "ailang-read", "", ""}})
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(func() { srv.Close(); _ = d.Close() })
	code, want, stderr := runCLI(t, srv.URL, "episode", "report", "w1")
	if code != exitOK || !strings.Contains(want, "episode w1: 2 calls") {
		t.Fatalf("the CLI: exit %d stdout %q stderr %q", code, want, stderr)
	}
	b.do(t, nil, "openCard", map[string]any{"url": srv.URL + "/surface", "raw": true})
	b.do(t, nil, "waitFor", map[string]any{"expression": `Boolean(window.World && document.getElementById('addr'))`})
	enter := func(line string) {
		t.Helper()
		b.do(t, nil, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const e = document.getElementById('addr'); e.focus(); e.value = %q; return document.activeElement === e; })()`, line)})
		b.do(t, nil, "key", map[string]any{"key": "Enter"})
	}
	panel := func(has string) string {
		t.Helper()
		var s string
		b.do(t, &s, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const p = document.querySelector('#verb pre'); return p && p.textContent.includes(%q) && p.textContent; })()`, has)})
		return s
	}
	enter("episode report w1")
	if got := panel("episode w1: 2 calls"); got != want {
		t.Fatalf("the bar shows\n%q\nthe CLI printed\n%q", got, want)
	}
	var state struct {
		Addr, Head, Location string
		Verb                 bool
	}
	read := func() {
		t.Helper()
		b.do(t, &state, "waitFor", map[string]any{"expression": `({ addr: document.getElementById('addr').value, head: (document.querySelector('#verb .mono') || {}).textContent || '',
			location: location.href, verb: Boolean(document.getElementById('verb')) })`})
	}
	if read(); state.Addr != "" || state.Head != "$ ailang-worldd episode report w1" {
		t.Fatalf("after episode report: %+v", state)
	}
	enter("why")
	if got := panel("a target is required"); got != "ailang-worldd why: a target is required" {
		t.Fatalf("why shows %q", got)
	}
	read()
	before := state.Location
	enter("log")
	if got := panel("ailang-worldd log"); got != "log is a CLI command the app does not run; at a terminal: ailang-worldd log" {
		t.Fatalf("log shows %q", got)
	}
	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('verb').remove(), 1)`})
	enter("world a b c")
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('hmsg').textContent.includes('unknown address: try')`})
	if read(); state.Verb || state.Addr != "world a b c" || state.Location != before {
		t.Fatalf("after world a b c: %+v (was at %s)", state, before)
	}
	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
}
