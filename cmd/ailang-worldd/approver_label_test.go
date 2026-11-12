package main

// Row 201 A2 sprint 3 M6b (design_docs/planned/w-world-approvals-by-key-a2.md §2.12b, CW9): approver
// list and revoke show a passkey approver's passkey. Every store is under t.TempDir(); every key is
// approvaltest.Passkey's fixed seed (D-WORLD-109).

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// The test passkeys' origin and relying party ID: the daemon's own page on this machine.
const pkOrigin, pkRPID = "http://localhost:7644", "localhost"

// seedPasskey enrols name with pk's key at `at` straight through the store (no verb, no terminal),
// as approver enroll --format webauthn stores it.
func seedPasskey(t *testing.T, creds, name string, pk *approvaltest.Passkey, at int64) {
	t.Helper()
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	k := pk.Key()
	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: at},
		store.PasskeyRow{Name: name, RPID: pkRPID, Origin: pkOrigin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
		t.Fatal(err)
	}
}

// CW9: list labels a passkey (L4) and a hand-made webauthn row with no passkey (L5); revoke shows
// the passkey's label; a schema-5 store lists its ssh approvers and is asked for no passkey.
func TestApproverListAndRevokeShowPasskeys(t *testing.T) {
	es, ed := approvaltest.NewPasskey(approval.AlgES256, 21), approvaltest.NewPasskey(approval.AlgEdDSA, 22)
	esLabel := "passkey ES256 localhost cred " + passkeyCred(es.CredentialID())
	db, creds := approverStore(t)
	seedPasskey(t, creds, "a", es, approveT0)
	seedPasskey(t, creds, "b", ed, approveT0)
	revokeSeeded(t, creds, "b", approveT0+200)
	raw := rawCredentials(t, creds)
	if _, err := raw.Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('c', 'webauthn', x'a0', 1800000000)`); err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()
	var out, errw bytes.Buffer
	if code := runApprover([]string{"list", "--db", db}, &out, &errw, sessionEnv{}); code != exitOK {
		t.Fatalf("list: code %d; stderr %q", code, errw.String())
	}
	want := "a  webauthn  " + esLabel + "  enrolled 2027-01-15T08:00:00Z  last used never  active\n" +
		"b  webauthn  passkey EdDSA localhost cred " + passkeyCred(ed.CredentialID()) + "  enrolled 2027-01-15T08:00:00Z  last used never  revoked 2027-01-15T08:03:20Z\n" +
		"c  webauthn  invalid-passkey  enrolled 2027-01-15T08:00:00Z  last used never  active\n"
	if out.String() != want {
		t.Fatalf("list:\n%q\nwant:\n%q", out.String(), want)
	}
	if len(passkeyCred(es.CredentialID())) != 16 || passkeyCred(es.CredentialID()) == passkeyCred(ed.CredentialID()) {
		t.Fatalf("cred labels %q %q: want 16 distinct hex", passkeyCred(es.CredentialID()), passkeyCred(ed.CredentialID()))
	}

	term := &lockedTerm{in: strings.NewReader("y\n")}
	out.Reset()
	if code := runApprover([]string{"revoke", "a", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+300))); code != exitOK {
		t.Fatalf("revoke: code %d; stderr %q", code, errw.String())
	}
	if wantPrompt := "Revoke approver a (" + esLabel + ")? "; !strings.HasPrefix(term.Output(), wantPrompt) {
		t.Fatalf("revoke prompt %q, want the prefix %q", term.Output(), wantPrompt)
	}
	if out.String() != "revoked approver a\nno session was approved by it\n" {
		t.Fatalf("revoke stdout %q", out.String())
	}

	// A schema-5 store: list asks for no passkey of its ssh approvers, and a hand-made webauthn row
	// there names the restart.
	db5, creds5 := approverStore(t)
	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
	raw = rawCredentials(t, creds5)
	for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
		if _, err := raw.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	_ = raw.Close()
	out.Reset()
	if code := runApprover([]string{"list", "--db", db5}, &out, &errw, sessionEnv{}); code != exitOK || !strings.HasPrefix(out.String(), "s  ssh  ssh-ed25519 "+edFingerprint) {
		t.Fatalf("v5 list: code %d, stdout %q; stderr %q", code, out.String(), errw.String())
	}
	raw = rawCredentials(t, creds5)
	if _, err := raw.Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('w', 'webauthn', x'a0', 1800000000)`); err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()
	errw.Reset()
	if code := runApprover([]string{"list", "--db", db5}, &out, &errw, sessionEnv{}); code != exitUsage || !strings.Contains(errw.String(), "passkeys need schema 6") {
		t.Fatalf("v5 list with a webauthn row: code %d, stderr %q; want exitUsage naming schema 6", code, errw.String())
	}
}
