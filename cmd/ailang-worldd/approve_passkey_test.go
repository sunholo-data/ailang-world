package main

// Row 201 A2 sprint 3 M6f (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW7): approve
// --webauthn-link and --webauthn-assertion. Every store is under t.TempDir(); every passkey is
// approvaltest.Passkey's fixed seed (D-WORLD-109).

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// passkeySigner is the enrolled passkey of the approve tests (mark-mac, stored counter 1).
var passkeySigner = approvaltest.NewPasskey(approval.AlgEdDSA, 51)

// requestHashOf is the message hash a passkey signs to approve the rig's stored row.
func requestHashOf(t *testing.T, a approveRig) [32]byte {
	t.Helper()
	row := a.row(t)
	d, err := hex.DecodeString(row.Digest)
	if err != nil {
		t.Fatal(err)
	}
	c, err := hex.DecodeString(row.Challenge)
	if err != nil {
		t.Fatal(err)
	}
	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// writeAssertion writes an assertion file under t.TempDir() and returns its path.
func writeAssertion(t *testing.T, file []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "world-assertion.json")
	if err := os.WriteFile(p, file, 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// passkeyOf reads name's passkey row.
func passkeyOf(t *testing.T, creds, name string) store.PasskeyRow {
	t.Helper()
	ro, err := store.OpenCredentialsReadOnly(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p, ok, err := ro.Passkey(ctx, name)
	if err != nil || !ok {
		t.Fatalf("passkey %q: %v %v", name, ok, err)
	}
	return p
}

// passkeyRig is a request for ep1 at approveT0 with passkeySigner enrolled as mark-mac at
// approveT0-100.
func passkeyRig(t *testing.T) approveRig {
	t.Helper()
	a := newApproveRig(t, "ep1", approveT0)
	seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0-100)
	return a
}

// CW7: --webauthn-link prints one get link per origin (the request, its expiry and message hash,
// and that origin's active credentials, by name), opens no terminal and writes nothing; with no
// passkey it refuses (WF1), and below schema 6 it names the restart. --webauthn-assertion with the
// passkey's assertion over the link's challenge mints with no terminal: approved_by
// webauthn:mark-mac, the counter advanced, the effects the typed y makes. The same file again is
// refused as already used, and revoke lists the approval (K8).
func TestApproveByPasskey(t *testing.T) {
	t.Run("WF1-none-enrolled", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		calls := countDaemonCalls(t)
		var opened bool
		var stdout, stderr bytes.Buffer
		if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitUsage || opened || *calls != 0 ||
			stderr.String() != "ailang-worldd approve: no passkey approver is enrolled; see ailang-worldd approver list\n" {
			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
		}
		raw := rawCredentials(t, a.credPath)
		for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
			if _, err := raw.Exec(stmt); err != nil {
				t.Fatal(err)
			}
		}
		_ = raw.Close()
		stderr.Reset()
		if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitUsage || !strings.Contains(stderr.String(), "passkeys need schema 6") {
			t.Fatalf("v5: rc %d, stderr %q; want exitUsage naming schema 6", code, stderr.String())
		}
	})

	a := passkeyRig(t)
	other, third, revoked := approvaltest.NewPasskey(approval.AlgES256, 52), approvaltest.NewPasskey(approval.AlgEdDSA, 53), approvaltest.NewPasskey(approval.AlgEdDSA, 54)
	seedPasskey(t, a.credPath, "b-other", other, approveT0-100)
	seedPasskey(t, a.credPath, "c-tailnet", third, approveT0-100)
	seedPasskey(t, a.credPath, "d-revoked", revoked, approveT0-100)
	revokeSeeded(t, a.credPath, "d-revoked", approveT0-50)
	raw := rawCredentials(t, a.credPath)
	if _, err := raw.Exec(`UPDATE webauthn_credentials SET origin = 'https://x.ts.net', rp_id = 'x.ts.net' WHERE name = 'c-tailnet'`); err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()
	calls := countDaemonCalls(t)
	var opened bool
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened || *calls != 0 {
		t.Fatalf("link: rc %d, terminal %v, daemon calls %d; stderr %q", code, opened, *calls, stderr.String())
	}
	h := requestHashOf(t, a)
	row := a.row(t)
	b64 := base64.RawURLEncoding.EncodeToString
	get := func(rp, origin string, allow ...[]byte) ceremonyLink {
		l := ceremonyLink{V: 1, Op: "get", RP: rp, Origin: origin, Request: a.id, Expires: row.ExpiresAt, Hash16: hex.EncodeToString(h[:8]), Challenge: b64(h[:])}
		for _, id := range allow {
			l.Allow = append(l.Allow, b64(id))
		}
		return l
	}
	want := []ceremonyLink{get(pkRPID, pkOrigin, other.CredentialID(), passkeySigner.CredentialID()), get("x.ts.net", "https://x.ts.net", third.CredentialID())}
	out := stdout.String()
	if n := strings.Count(out, "/surface/passkey#"); n != 2 {
		t.Fatalf("%d links, want 2:\n%s", n, out)
	}
	_, second, _ := strings.Cut(out, "/surface/passkey#")
	if got := []ceremonyLink{linkFragment(t, out), linkFragment(t, second)}; !reflect.DeepEqual(got, want) {
		t.Fatalf("links %+v, want %+v", got, want)
	}
	for _, w := range []string{"approve request " + a.id + " (expires " + approverTime(row.ExpiresAt) + ", message hash " + hex.EncodeToString(h[:8]) + ")",
		"\n  " + pkOrigin + "/surface/passkey#", "\n  https://x.ts.net/surface/passkey#", "ailang-worldd approve " + a.id + " --webauthn-assertion <the saved file> --db "} {
		if !strings.Contains(out, w) {
			t.Fatalf("stdout lacks %q:\n%s", w, out)
		}
	}
	if row := a.row(t); row.Consumed || len(a.sessionRows(t)) != 0 {
		t.Fatal("--webauthn-link consumed the request or minted")
	}

	file := writeAssertion(t, passkeySigner.Assert(pkOrigin, pkRPID, h, 2))
	stdout.Reset()
	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened {
		t.Fatalf("assertion: rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
	}
	row = a.row(t)
	if !row.Consumed || row.ApprovedBy != "webauthn:mark-mac" {
		t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", row.Consumed, row.ApprovedBy)
	}
	if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 2 {
		t.Fatalf("sign count %d, want 2", p.SignCount)
	}
	if !strings.Contains(stdout.String(), "approved by webauthn:mark-mac\n") {
		t.Fatalf("stdout lacks the approver:\n%s", stdout.String())
	}
	mr, err := approval.Decode(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	var probs []string
	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, stdout.String())) {
		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty; this path records webauthn:mark-mac (above)
			probs = append(probs, p)
		}
	}
	if len(probs) != 0 {
		t.Fatalf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
	}

	stderr.Reset()
	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitUsage || !strings.Contains(stderr.String(), "was already used") || opened {
		t.Fatalf("replay: rc %d, stderr %q; want exitUsage, already used", code, stderr.String())
	}

	term := &lockedTerm{in: strings.NewReader("y\n")}
	stdout.Reset()
	if code := runApprover([]string{"revoke", "mark-mac", "--db", a.r.f.db}, &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+3))); code != exitOK ||
		!strings.Contains(stdout.String(), "it approved 1 session(s)") || !strings.Contains(stdout.String(), "  "+a.id+"  episode ep1") {
		t.Fatalf("revoke: rc %d, stdout %q; want the passkey's approval listed (K8)", code, stdout.String())
	}
}
