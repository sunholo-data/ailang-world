package main

// Row 201 A4 sprint 3 M6f (design §2.10, §6 CD8-CD9; the real-tool sibling of DM7 and CD8-CD9, plan
// §0 retro 11): approve <request-ref> --webauthn-link in a live Chrome, gated as XB1 (WORLD_CHROME,
// WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them). Chrome's virtual authenticator
// (a test key, D-WORLD-109, temp profile) creates a passkey through the A2 page, the decision link
// the CLI printed signs the decision message's hash through the same page, unchanged, and the file
// it saved, posted by --webauthn-assertion, is the daemon's law's decision: Chrome's bytes, not
// approvaltest's.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestApproveDecideWebAuthnLinkInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := passkeyPageOrigin(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
	enrol := sha256.Sum256([]byte("decide-enrol"))
	dl := t.TempDir()
	var page pageCeremony
	b.do(t, nil, "open", map[string]any{"url": ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac",
		Challenge: base64.RawURLEncoding.EncodeToString(enrol[:])}.link()})
	b.do(t, &page, "ceremony", map[string]any{"download": dl})
	reg := saved(t, page, origin, dl, "world-registration-mark-mac.json")
	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: enrol, Response: reg, Now: 1000})
	if err != nil {
		t.Fatal(err)
	}
	f := newFakeDecider(t)
	f.enrolDecisionPasskey(t, false)
	c, err := store.OpenCredentials(f.creds)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p := a.Passkey
	if err := c.EnrolPasskey(ctx, store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: a.EnrolledAt},
		store.PasskeyRow{Name: a.Name, RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, AAGUID: p.AAGUID[:], Alg: p.Key.Alg,
			BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)}); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	env := approveEnv(nil, nil, atTime(1710))
	code, stdout, stderr := f.keyRun(t, env, "--webauthn-link")
	link, err := printedLink(stdout)
	if code != exitOK || err != nil {
		t.Fatalf("--webauthn-link: rc %d, %v\nstdout %s\nstderr %s", code, err, stdout, stderr)
	}
	var view pageView
	b.do(t, &view, "open", map[string]any{"url": link})
	ref16 := strings.TrimPrefix(f.ref, "sha256:")[:16]
	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbApprove, 1760))
	if view.State != "ready" || !strings.Contains(view.What, "decide-approve-"+ref16) || !strings.Contains(view.What, "message hash "+base64Hex(h[:8])) {
		t.Fatalf("the page shows %+v", view)
	}
	dl2 := t.TempDir()
	b.do(t, &page, "ceremony", map[string]any{"download": dl2})
	asr := saved(t, page, origin, dl2, "world-assertion-decide-approve-"+ref16+".json")
	path := filepath.Join(dl2, "world-assertion-decide-approve-"+ref16+".json")
	if code, _, stderr := f.keyRun(t, env, "--webauthn-assertion", path, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 || f.postBodies[0].WebAuthn == nil {
		t.Fatalf("--webauthn-assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
	got, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: refDigest(f.subject.Packet), Request: refDigest(f.ref),
		Expiry: 1760, Now: 1710, Response: *f.postBodies[0].WebAuthn, Approver: a})
	if !bytes.Equal(bytes.TrimSpace(asr), []byte(*f.postBodies[0].WebAuthn)) {
		t.Fatalf("the POST did not carry the saved file:\n%s\n%s", asr, *f.postBodies[0].WebAuthn)
	}
	if err != nil || got.DecidedBy() != "webauthn:mark-mac" {
		t.Fatalf("Chrome's decision assertion does not decide: %v", err)
	}
}
