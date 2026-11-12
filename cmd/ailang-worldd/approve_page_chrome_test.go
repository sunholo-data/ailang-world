package main

// Row 201 A3 sprint 3 M8a-M8b (design §2.7, §6 XA1, XA2): the request card in a live Chrome.

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// XA1: mark-mac enrolled through the A2 page and the CLI (as XB1); the card at
// /surface/approve/<id> shows every field of the stored request as its text; Approve records the
// passkey's assertion (200) and approve --await mints it; no CSP violation and no page error; the
// card's get has the same challenge, rpId, credential and user verification as the A2 page's get for
// the same request (the CLI's --webauthn-link). Opened again after the mint, the card says it is used.
func TestApprovePageInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
	now := time.Now().Unix() // the daemon's routes read the wall clock
	a := newApproveRig(t, "ep1", now)
	origin := cardOrigin(t, a.credPath, "", nil)
	t.Logf("%s, card at %s", b.chrome, origin)

	dl := t.TempDir()
	var page pageCeremony
	var out, errw bytes.Buffer
	term := &ceremonyTerm{}
	term.answer = func(n int) (string, error) {
		if n > 1 {
			return "", os.ErrClosed
		}
		if n == 1 {
			return "y\n", nil
		}
		link, err := printedLink(out.String())
		if err == nil {
			err = b.call(nil, "open", map[string]any{"url": link})
		}
		if err == nil {
			err = b.call(&page, "ceremony", map[string]any{"download": dl})
		}
		if err != nil {
			t.Error(err)
			return "", err
		}
		return page.After.Code + "\n", nil
	}
	env := enrolEnv(term, nil)
	env.now = atTime(now - 100)
	args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", filepath.Join(dl, "world-registration-mark-mac.json"), "--db", a.r.f.db}
	if code := runApprover(args, &out, &errw, env); code != exitOK {
		t.Fatalf("enrol: rc %d\nstdout %s\nstderr %s", code, out.String(), errw.String())
	}

	var lout, lerr bytes.Buffer
	if code := runApprove(a.args("--webauthn-link"), &lout, &lerr, approveEnv(nil, nil, atTime(now))); code != exitOK {
		t.Fatalf("--webauthn-link: rc %d\nstderr %s", code, lerr.String())
	}
	link, err := printedLink(lout.String())
	if err != nil {
		t.Fatal(err)
	}
	b.do(t, nil, "open", map[string]any{"url": link})
	b.do(t, &page, "ceremony", nil)
	var a2 []cardGet
	b.do(t, &a2, "waitFor", map[string]any{"expression": "window.__gets"})

	var view cardView
	b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id})
	h := requestHashOf(t, a)
	if view.State != "ready" || view.ButtonHidden || view.ButtonDisabled || !strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1, expires ") ||
		!strings.HasSuffix(view.What, ", message hash "+hex.EncodeToString(h[:8])) || view.SaveHidden || view.SaveName != "world-approval-"+a.id+".msg" {
		t.Fatalf("the card shows %+v", view)
	}
	if want := cardFields(t, a.row(t).Request); len(want) < 28 || !reflect.DeepEqual(view.Fields, want) || view.TTLLabel != "Session lifetime (s) (1.0 h)" {
		t.Fatalf("the card's fields\n%v\nwant every leaf of the stored request\n%v", view.Fields, want)
	}
	var after cardView // a fresh value: decoding into view would keep its fields
	b.do(t, &after, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}})
	if after.State != "recorded" || after.StatusError || !strings.HasPrefix(after.Status, "Approved by webauthn:mark-mac: the minter (approve "+a.id+" --await) mints it") {
		t.Fatalf("after Approve the card shows %+v", after)
	}
	if rec, ok := surfaceRecorded(t, a.credPath, a.id); !ok || rec.Approver != "webauthn:mark-mac" || rec.Format != "webauthn" {
		t.Fatalf("recorded %+v %v", rec, ok)
	}
	if len(after.Violations) != 0 || len(after.Errors) != 0 {
		t.Fatalf("CSP violations %q, page errors %q", after.Violations, after.Errors)
	}
	if len(a2) != 1 || len(after.Gets) != 1 || !reflect.DeepEqual(after.Gets[0], a2[0]) || after.Gets[0].RPID != "localhost" ||
		after.Gets[0].UserVerification != "required" || after.Gets[0].Timeout != 240000 || !after.Gets[0].Signal || len(after.Gets[0].Allow) != 1 {
		t.Fatalf("the card's get %+v is not the A2 page's %+v", after.Gets, a2)
	}
	var opened bool
	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
	}
	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
	var used cardView
	if b.do(t, &used, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id}); used.State != "used" || used.Status != "Already used: approved by webauthn:mark-mac." || !used.ButtonHidden {
		t.Fatalf("the used card shows %+v", used)
	}
}
