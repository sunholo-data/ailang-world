package main

// Row 243 V3 sprint 2 M6 (design §2.5, §6 XC1): the confirmation card in a live Chrome, on the A2/A3
// driver (testdata/passkey_cdp.mjs, unchanged), gated as XB1 (WORLD_CHROME,
// WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them). No daemon route serves the card
// before sprint 3, so a loopback test handler serves surfacepage.ConfirmPage and ConfirmScript under
// the headers the daemon gives its own page scripts (surfaceCSP among them), and the card's data, made
// by surfacepage.ConfirmationCard from approval.CanonicalVerb's bytes. No store is read.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// XC1: the card shows, in its <pre>, exactly broker.FinishPromptText of the request's scope, title,
// body, protected paths and committer (the terminal's lines before its question), and Go's
// VerbMessageHash, under the daemon's CSP with no violation.
func TestConfirmPageShowsTheFinishPromptInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{})
	now := time.Now().Unix()
	scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: strings.Repeat("ab", 20), Base: "dev",
		TitleSHA256: strings.Repeat("0", 64), BodySHA256: strings.Repeat("1", 64), Draft: true, ExpiresAt: now + 900}
	title, body, protected, committer := "Fix the <bar>", "## Summary\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7 is not in the prompt\n", []string{"host/x.go"}, ""
	prompt := broker.FinishPromptText(scope, title, body, protected, committer)
	var tty bytes.Buffer // control: the terminal prints exactly these lines, then its question
	finishPrompt(&tty, scope, title, body, protected, committer)
	if tty.String() != prompt+"Open this pull request? [y/N] " {
		t.Fatalf("FinishPromptText is not the terminal's prompt: %q", tty.String())
	}
	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", "w1", "--preset", "p"},
		Episode: "w1", Prompt: prompt, Scope: broker.FormatFinishApprovalScope(scope), Title: title, Body: body, Committer: committer, Protected: protected}
	req, err := approval.CanonicalVerb(r)
	if err == nil {
		err = r.Validate()
	}
	if err != nil {
		t.Fatal(err)
	}
	digest, challenge := sha256.Sum256(req), sha256.Sum256([]byte("xc1 challenge"))
	const id = "vr-00000000000000c1"
	card, err := json.Marshal(surfacepage.ConfirmationCard(surfacepage.ConfirmationInput{ID: id, Request: req, Digest: hex.EncodeToString(digest[:]),
		Challenge: hex.EncodeToString(challenge[:]), Created: now, Expires: now + 600, Now: now}, nil))
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := approval.VerbMessageHash(digest, challenge, now+600)

	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	own := httptest.NewRecorder() // the headers the daemon serves its passkey page's script with
	d.Handler().ServeHTTP(own, httptest.NewRequest("GET", "/surface/passkey.js", nil))
	if csp := own.Header().Get("Content-Security-Policy"); own.Code != 200 || !strings.Contains(csp, "script-src 'self'") || !strings.Contains(csp, "connect-src 'self'") {
		t.Fatalf("the daemon's own script headers: %d %q", own.Code, csp)
	}
	serve := func(w http.ResponseWriter, contentType string, b []byte) {
		for k, v := range own.Header() {
			w.Header()[k] = v
		}
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(b)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
		switch q.URL.Path {
		case "/surface/confirm/" + id:
			serve(w, "text/html; charset=utf-8", surfacepage.ConfirmPage)
		case "/surface/confirm.js":
			serve(w, "text/javascript; charset=utf-8", surfacepage.ConfirmScript)
		case "/surface/confirmations/" + id:
			serve(w, "application/json", card)
		default:
			http.NotFound(w, q)
		}
	}))
	t.Cleanup(srv.Close)
	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)

	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface/confirm/" + id, "raw": true})
	var view struct {
		State, Prompt, Argv, Hash, Save, SSHLine, Title string
		Violations                                      []string
	}
	b.do(t, &view, "waitFor", map[string]any{"expression": `document.body.dataset.state !== 'loading' && (() => { const $ = (id) => document.getElementById(id); return {
		state: document.body.dataset.state, prompt: $('prompt').textContent, argv: $('argv').textContent, hash: $('m-hash') ? $('m-hash').textContent : '',
		save: $('save').download, sshLine: $('sshline').textContent, title: $('f-title') ? $('f-title').textContent : '', violations: window.__cspv || [] }; })()`})
	if view.State != "ready" || view.Prompt != prompt || view.Title != title || view.Argv != "ailang-worldd episode finish w1 --preset p" {
		t.Fatalf("the card shows %+v; want the prompt %q", view, prompt)
	}
	if view.Hash != hex.EncodeToString(hash[:]) || view.Save != "world-verb-"+id+".msg" || view.SSHLine != "ssh-keygen -Y sign -n world-verb@sunholo -f <key> world-verb-"+id+".msg" || len(view.Violations) != 0 {
		t.Fatalf("the card's message, link or CSP: %+v", view)
	}
}
