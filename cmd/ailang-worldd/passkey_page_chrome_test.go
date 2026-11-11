package main

// Row 201 A2 sprint 4 M7, XB1 (design §2.14, §6): the page in a live Chrome, end to end through the
// CLI. The links it opens are the ones the CLI printed; the files it saves go back through the CLI.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// XB1: the CLI's create link enrols, through approver enroll --format webauthn, the file the page
// saved, with the code the page showed; the CLI's get link approves, through approve
// --webauthn-assertion, the file the page saved, after a get on an authenticator without user
// verification saved nothing; the finished attempt's timer never fires. A create on an authenticator
// without user verification saves nothing, and a link for another page is refused before any button.
// (Whether the button is disabled while pending is XB2's: a virtual authenticator answers too fast to
// observe it.)
func TestPasskeyPageInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := passkeyPageOrigin(t)
	t.Logf("%s, page at %s", b.chrome, origin)

	t.Run("enrol-then-approve", func(t *testing.T) {
		b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
		a := newApproveRig(t, "ep1", approveT0)
		dl := t.TempDir()
		reg := filepath.Join(dl, "world-registration-mark-mac.json")
		var view pageView
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
				err = b.call(&view, "open", map[string]any{"url": link})
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
		args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", reg, "--db", a.r.f.db}
		if code := runApprover(args, &out, &errw, enrolEnv(term, nil)); code != exitOK {
			t.Fatalf("enrol: rc %d\nstdout %s\nstderr %s", code, out.String(), errw.String())
		}
		if view.State != "ready" || view.What != "Create a passkey for approver mark-mac at localhost" || view.Button != "Create passkey" || view.ButtonDisabled {
			t.Fatalf("the create page shows %+v", view)
		}
		file := saved(t, page, origin, dl, "world-registration-mark-mac.json")
		if page.After.Code != codeOf(t, file) {
			t.Fatalf("the page's code %q is not Go's %q for the saved file", page.After.Code, codeOf(t, file))
		}
		row, pk := enrolled(t, a.credPath)
		if row.Name != "mark-mac" || pk.Origin != origin || pk.RPID != "localhost" || pk.Alg != approval.AlgES256 || !strings.Contains(out.String(), "enrolled approver mark-mac: passkey ") {
			t.Fatalf("enrolled %+v %+v\nstdout %s", row, pk, out.String())
		}

		var lout, lerr bytes.Buffer
		var opened bool
		if code := runApprove(a.args("--webauthn-link"), &lout, &lerr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened {
			t.Fatalf("--webauthn-link: rc %d, terminal %v\nstderr %s", code, opened, lerr.String())
		}
		link, err := printedLink(lout.String())
		if err != nil {
			t.Fatal(err)
		}
		b.do(t, &view, "open", map[string]any{"url": link})
		h := requestHashOf(t, a)
		if view.State != "ready" || !strings.HasPrefix(view.What, "Approve request "+a.id+", expires ") || !strings.HasSuffix(view.What, ", message hash "+hex.EncodeToString(h[:8])) || view.Button != "Approve with passkey" {
			t.Fatalf("the get page shows %+v", view)
		}
		var moved struct{ Moved int }
		b.do(t, &moved, "uv", map[string]any{"verified": false}) // an authenticator with no user verification: nothing to save
		b.do(t, &page, "ceremony", map[string]any{"download": dl})
		if files := worldFiles(t, dl); page.After.State != "failed" || !page.After.StatusError || !page.After.SaveHidden || page.SaveClicked || len(page.Downloads) != 0 || len(files) != 1 || page.After.ButtonDisabled {
			t.Fatalf("a get without user verification: %+v, downloads %+v, files %q", page.After, page.Downloads, files)
		}
		if moved.Moved != 1 {
			t.Fatalf("%d credentials moved, want 1", moved.Moved)
		}
		b.do(t, nil, "uv", map[string]any{"verified": true}) // the retry, on the same page
		b.do(t, &page, "ceremony", map[string]any{"download": dl})
		name := "world-assertion-" + a.id + ".json"
		saved(t, page, origin, dl, name)
		b.do(t, &view, "advance", map[string]any{"budget": 250000}) // the finished attempt's timer was cleared
		if view.State != "done" || view.StatusError || view.SaveHidden {
			t.Fatalf("240 s after the result the page shows %+v", view)
		}
		if code := runApprove(a.args("--webauthn-assertion", filepath.Join(dl, name)), &lout, &lerr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitOK || opened {
			t.Fatalf("--webauthn-assertion: rc %d, terminal %v\nstderr %s", code, opened, lerr.String())
		}
		if r := a.row(t); !r.Consumed || r.ApprovedBy != "webauthn:mark-mac" {
			t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", r.Consumed, r.ApprovedBy)
		}
		if after := passkeyOf(t, a.credPath, "mark-mac"); after.SignCount <= pk.SignCount {
			t.Fatalf("sign count %d after the approval, %d at enrolment", after.SignCount, pk.SignCount)
		}
	})

	t.Run("user-verification-required", func(t *testing.T) {
		b.do(t, nil, "tab", map[string]any{"authenticator": "no-uv"})
		db, creds := approverStore(t)
		dl := t.TempDir()
		var page pageCeremony
		var out, errw bytes.Buffer
		term := &ceremonyTerm{}
		term.answer = func(n int) (string, error) {
			if n > 0 {
				return "", os.ErrClosed
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
			}
			return page.After.Code + "\n", err
		}
		args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", filepath.Join(dl, "world-registration-mark-mac.json"), "--db", db}
		code := runApprover(args, &out, &errw, enrolEnv(term, nil))
		if page.After.State != "failed" || !page.After.StatusError || !page.After.SaveHidden || page.SaveClicked || len(page.Downloads) != 0 || page.After.ButtonDisabled {
			t.Fatalf("a UV-false authenticator: %+v, downloads %+v", page.After, page.Downloads)
		}
		if files := worldFiles(t, dl); code != exitUsage || len(files) != 0 || len(approverRows(t, creds)) != 0 {
			t.Fatalf("rc %d, files saved %q, %d approvers; stderr %s", code, files, len(approverRows(t, creds)), errw.String())
		}
	})

	t.Run("links-for-another-page", func(t *testing.T) {
		b.do(t, nil, "tab", nil)
		ok := ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: base64.RawURLEncoding.EncodeToString(make([]byte, 32))}
		frag := func(l ceremonyLink) string { _, f, _ := strings.Cut(l.link(), "#"); return f }
		edit := func(e func(*ceremonyLink)) string { l := ok; e(&l); return frag(l) }
		for _, tc := range []struct{ name, fragment, state, what string }{
			{"control", frag(ok), "ready", "Create a passkey for approver mark-mac at localhost"},
			{"version", edit(func(l *ceremonyLink) { l.V = 2 }), "refused", "(not a version 1 create or get)"},
			{"op", edit(func(l *ceremonyLink) { l.Op = "sign" }), "refused", "(not a version 1 create or get)"},
			{"origin", edit(func(l *ceremonyLink) { l.Origin = "http://localhost:7644" }), "refused", "This link is for http://localhost:7644 (passkey site localhost): open it there, not at " + origin + "."},
			{"rp", edit(func(l *ceremonyLink) { l.RP = "x.ts.net" }), "refused", "This link is for " + origin + " (passkey site x.ts.net)"},
			{"not-base64url", "e30=", "refused", "(not base64url)"},
			{"not-json", base64.RawURLEncoding.EncodeToString([]byte("{")), "refused", "This is not a World passkey link ("},
		} {
			var view pageView
			b.do(t, &view, "open", map[string]any{"url": origin + "/surface/passkey#" + tc.fragment})
			if view.State != tc.state || !strings.Contains(view.What, tc.what) || view.ButtonHidden != (tc.state == "refused") {
				t.Errorf("%s: state %q, button hidden %v, what %q; want %q containing %q", tc.name, view.State, view.ButtonHidden, view.What, tc.state, tc.what)
			}
		}
	})
}
