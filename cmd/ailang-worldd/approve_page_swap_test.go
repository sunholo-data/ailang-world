package main

// Row 201 A3 sprint 3 M8b (design §2.7, §4 T4, §6 XA2; AW6's page half): the card refuses what it
// cannot bind, and its SSH path, in a live Chrome.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// XA2: the card's response rewritten in the browser. (a) bytes changed, digest kept (or another
// request's id): refused, no get; (b) bytes and digest forged together: the page signs, the daemon
// refuses the POST (409, R9: the body's digest is not the request's), nothing is recorded; (c) a key
// added, branch deleted, reuse a string, ttl_s a fraction, or a byte that is not UTF-8, each re-hashed:
// refused naming the path, no get and no message link; a non-hex nonce, a negative expiry or bytes that
// are not base64url: refused.
// Control: no rewrite, ready.
func TestApprovePageRefusesASwappedRequest(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	origin := cardOrigin(t, a.credPath, "", nil)
	virtualPasskey(t, b, a.credPath, origin, "mark-mac", 81, now-100)
	card := origin + "/surface/approve/" + a.id
	for _, tc := range []struct{ mode, state, what string }{
		{"", "ready", "Approve request " + a.id},
		{"bytes", "refused", "The daemon's request does not match its digest; nothing to approve."},
		{"forge", "ready", "Approve request " + a.id},
		{"extra", "refused", "This request has a field the page cannot show (extra); nothing to approve; approve it at a terminal."},
		{"drop:branch", "refused", "This request has a field the page cannot show (branch); nothing to approve; approve it at a terminal."},
		{"type", "refused", "This request has a field the page cannot show (reuse); nothing to approve; approve it at a terminal."},
		{"float", "refused", "This request has a field the page cannot show (ttl_s); nothing to approve; approve it at a terminal."},
		{"utf8", "refused", "The daemon's request is not UTF-8 JSON ("},
		{"id", "refused", "The daemon's request does not match its digest; nothing to approve."},
		{"nonce", "refused", "The daemon's challenge or expiry is malformed; nothing to approve."},
		{"expiry", "refused", "The daemon's challenge or expiry is malformed; nothing to approve."},
		{"b64", "refused", "The daemon's request is not base64url; nothing to approve."},
	} {
		var view cardView // a fresh value: decoding into an old one would keep its fields
		b.do(t, nil, "fetchRewrite", map[string]any{"mode": tc.mode})
		b.do(t, &view, "openCard", map[string]any{"url": card})
		if view.State != tc.state || !strings.HasPrefix(view.What, tc.what) {
			t.Fatalf("%q: state %q, what %q; want %q, %q", tc.mode, view.State, view.What, tc.state, tc.what)
		}
		if tc.state == "refused" && (len(view.Gets) != 0 || !view.ButtonHidden || !view.SaveHidden || view.SaveHref != "" || len(view.Fields) != 0) {
			t.Fatalf("%q: a refused card offers something: %+v", tc.mode, view)
		}
		if tc.mode != "forge" {
			continue
		}
		b.do(t, &view, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}})
		if view.State != "rejected" || len(view.Gets) != 1 || !strings.Contains(view.Status, "HTTP 409 Conflict): the request changed since the card was shown") {
			t.Fatalf("forged: the page shows %+v", view)
		}
		if _, ok := surfaceRecorded(t, a.credPath, a.id); ok {
			t.Fatal("forged: an approval was recorded")
		}
	}
	var n struct{ Rewritten int }
	b.do(t, &n, "fetchRewrite", map[string]any{"mode": ""})
	if n.Rewritten != 11 {
		t.Fatalf("%d card responses rewritten, want 11", n.Rewritten)
	}
}

// The card's SSH path (planner, the page half of AW6): Save the message saves exactly the request's
// 90-byte message (approval.SignedMessage) through a blob: download of the card's origin; its
// approvaltest SSHSIG, pasted, is recorded as ssh:mark-mac (200), and approve --await mints it. A
// signature over another message is refused (422) and records nothing. The only passkey is enrolled
// at another origin: the card names it and hides the passkey button.
func TestApprovePageSSHSIGInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), now-100)
	origin := cardOrigin(t, a.credPath, "", nil)
	virtualPasskey(t, b, a.credPath, "https://studio.tail0.ts.net", "mark-phone", 83, now-100)
	var view cardView
	b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id})
	if view.State != "ready" || !view.ButtonHidden || view.Note != "No passkey is enrolled for "+origin+": open this card at https://studio.tail0.ts.net." || view.SaveHidden || view.SSHLine != "ssh-keygen -Y sign -n world-approval@sunholo -f <key> world-approval-"+a.id+".msg" {
		t.Fatalf("the card shows %+v", view)
	}
	dl := t.TempDir()
	var downloads []struct{ Event, State, Name, URL string }
	b.do(t, &downloads, "save", map[string]any{"id": "save", "dir": dl})
	name := "world-approval-" + a.id + ".msg"
	msg, err := os.ReadFile(filepath.Join(dl, name))
	if want := requestMessage(t, a); err != nil || !bytes.Equal(msg, want) || len(downloads) == 0 || !strings.HasPrefix(downloads[0].URL, "blob:"+origin+"/") {
		t.Fatalf("saved %x (%v) via %+v, want the request's message %x", msg, err, downloads, want)
	}
	var after cardView
	b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sigSigner.Sign("world-approval@sunholo", append([]byte("x"), msg...)))})
	b.do(t, &after, "press", map[string]any{"id": "post-sig", "states": []string{"recorded", "rejected", "unknown"}})
	if after.State != "rejected" || !strings.Contains(after.Status, "(HTTP 422 NotApproved): the SSH signature does not approve request "+a.id) {
		t.Fatalf("a signature over another message: %+v", after)
	}
	if _, ok := surfaceRecorded(t, a.credPath, a.id); ok {
		t.Fatal("a refused signature was recorded")
	}
	b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sigSigner.Sign("world-approval@sunholo", msg))})
	var done cardView
	b.do(t, &done, "press", map[string]any{"id": "post-sig", "states": []string{"recorded", "rejected", "unknown"}})
	if done.State != "recorded" || !strings.HasPrefix(done.Status, "Approved by ssh:mark-mac: ") || len(done.Violations) != 0 {
		t.Fatalf("the signature: %+v", done)
	}
	var opened bool
	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
	}
	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
}
