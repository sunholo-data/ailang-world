package main

// Row 201 A4 sprint 3 M6e (design §2.10, §6 CD6-CD9): approve <request-ref> with an artifact made
// elsewhere, against the fake daemon (approve_decide_rig_test.go): --write-message and --signature
// --expiry (with the real ssh-keygen signing the written file, on a key made in a temp dir), and
// --webauthn-link and --webauthn-assertion --expiry (an approvaltest passkey signing the link's
// challenge, judged by the daemon's own law). D-WORLD-109: temp dirs, no real key, no agent.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// CD6, CD7: --write-message writes the message for the window's max to a NEW file (never over an
// existing one) and prints the signing line in the decision namespace and the --signature step with
// its expiry; --signature posts the file with the --expiry it was signed for, as sshsig; --signature
// without --expiry, --expiry alone and --expiry with a mint id are refused before any daemon call.
func TestApproveDecideWriteMessageAndSignature(t *testing.T) {
	f := newFakeDecider(t)
	env := approveEnv(nil, nil, atTime(1710))
	msgPath := filepath.Join(t.TempDir(), "decide.msg")
	code, stdout, stderr := f.keyRun(t, env, "--deny", "--write-message", msgPath)
	got, _ := os.ReadFile(msgPath)
	if code != exitOK || !bytes.Equal(got, f.decisionMessageOf(t, approval.VerbDeny, 1760)) || len(got) != 91 || len(f.posts) != 0 ||
		!strings.Contains(stdout, "ssh-keygen -Y sign -n world-decision@sunholo -f <your key> < ") || !strings.Contains(stdout, f.ref+" --deny --signature ") ||
		!strings.Contains(stdout, "--expiry 1760 --daemon "+f.addr()) || !strings.Contains(stdout, "deciding DENY") {
		t.Fatalf("write-message: rc %d, %d POSTs\nstdout %s\nstderr %s", code, len(f.posts), stdout, stderr)
	}
	if code, _, stderr := f.keyRun(t, env, "--write-message", msgPath); code != exitUsage || !strings.Contains(stderr, "nothing was written") {
		t.Fatalf("an existing file: rc %d %s", code, stderr)
	}
	sigPath := filepath.Join(t.TempDir(), "decide.sig")
	if err := os.WriteFile(sigPath, []byte("SIGNED-ELSEWHERE"), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := f.keyRun(t, env, "--signature", sigPath, "--expiry", "1755"); code != exitOK || len(f.posts) != 1 ||
		f.postBodies[0].SSHSIG == nil || *f.postBodies[0].SSHSIG != "SIGNED-ELSEWHERE" || f.postBodies[0].Expiry != 1755 || f.postBodies[0].Typed ||
		f.postBodies[0].WebAuthn != nil || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
		t.Fatalf("--signature: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
	missing := filepath.Join(t.TempDir(), "absent")
	for _, flag := range []string{"--signature", "--webauthn-assertion"} {
		if code, _, stderr := f.keyRun(t, env, flag, missing, "--expiry", "1755"); code != exitUsage || len(f.posts) != 1 || !strings.Contains(stderr, flag+": ") {
			t.Fatalf("%s of a missing file: rc %d, %d POSTs\nstderr %s", flag, code, len(f.posts), stderr)
		}
	}
	gets := f.gets
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{f.ref, "--signature", sigPath}, "go with --expiry <t>"},
		{[]string{f.ref, "--expiry", "1755"}, "go with --expiry <t>"},
		{[]string{"mr-0123456789abcdef", "--expiry", "1755"}, "--expiry take a Human.Approve request ref"},
	} {
		var stdout, stderr bytes.Buffer
		if code := runApprove(append(tc.args, "--db", f.db), &stdout, &stderr, env); code != exitUsage || f.gets != gets || !strings.Contains(stderr.String(), tc.want) { // no --daemon: it would trip CD3 itself
			t.Errorf("%v: rc %d, GETs %d\nstderr %s", tc.args, code, f.gets, stderr.String())
		}
	}
}

// CD6 + CD7 with the real tool: the file --write-message wrote, signed by the real ssh-keygen with
// the printed namespace on a key made in a temp dir (no agent, temp HOME), posted by --signature
// with the printed expiry, verifies with the daemon's own check (VerifyDecisionSSHSIG) and is not
// a mint signature. With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen fails (CI).
func TestApproveDecideWriteMessageSignedBySSHKeygen(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
		}
		t.Skip("SKIP COUNTED: ssh-keygen not found; the decision --write-message end-to-end test did NOT run")
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SSH_AUTH_SOCK", "")
	dir := t.TempDir()
	key := filepath.Join(dir, "k")
	sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "cd7", "-f", key)
	pub, err := os.ReadFile(key + ".pub")
	if err != nil {
		t.Fatal(err)
	}
	k, err := approval.ParseAuthorizedKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	f := newFakeDecider(t)
	env := approveEnv(nil, nil, atTime(1710))
	msgPath := filepath.Join(dir, "decide.msg")
	code, stdout, stderr := f.keyRun(t, env, "--write-message", msgPath)
	m := regexp.MustCompile(`ssh-keygen -Y sign -n (\S+) .* --expiry (\d+) `).FindStringSubmatch(strings.ReplaceAll(stdout, "\n", " "))
	if code != exitOK || m == nil {
		t.Fatalf("write-message: rc %d\nstdout %s\nstderr %s", code, stdout, stderr)
	}
	sshKeygen(t, "-Y", "sign", "-n", m[1], "-f", key, msgPath) // writes decide.msg.sig
	if code, _, stderr := f.keyRun(t, env, "--signature", msgPath+".sig", "--expiry", m[2]); code != exitOK || len(f.posts) != 1 {
		t.Fatalf("--signature: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
	armored := []byte(*f.postBodies[0].SSHSIG)
	if err := approval.VerifyDecisionSSHSIG(armored, approval.VerbApprove, refDigest(f.subject.Packet), refDigest(f.ref), 1760, k); err != nil {
		t.Fatalf("the real signature of the written file does not decide: %v", err)
	}
	if err := approval.VerifyApproval(armored, refDigest(f.ref), refDigest(f.subject.Packet), 1760, k); err == nil {
		t.Fatal("the decision signature verifies as a mint")
	}
}

var cdPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 73)

const cdOrigin, cdRP = "https://studio.x.ts.net", "studio.x.ts.net"

// enrolDecisionPasskey makes f's credentials store (replacing its empty file) and, with enrol,
// enrols mark-mac's passkey in it.
func (f *fakeDecider) enrolDecisionPasskey(t *testing.T, enrol bool) {
	t.Helper()
	if fi, err := os.Stat(f.creds); err == nil && fi.Size() == 0 {
		if err := os.Remove(f.creds); err != nil {
			t.Fatal(err)
		}
	}
	c, err := store.OpenCredentials(f.creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	k := cdPasskey.Key()
	if !enrol {
		return
	}
	if err := c.EnrolPasskey(ctx, store.ApproverRow{Name: "mark-mac", Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: 1000},
		store.PasskeyRow{Name: "mark-mac", RPID: cdRP, Origin: cdOrigin, CredentialID: cdPasskey.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
		t.Fatal(err)
	}
}

// CD8, CD9: --webauthn-link prints the passkey page's get link for the active passkey, whose
// challenge is the decision message's hash for the window's max (the verb in the link's request and
// in capitals, the hash the page shows) and the --webauthn-assertion step, and posts nothing; the
// assertion an authenticator makes over that challenge, posted by --webauthn-assertion with the
// printed expiry, is the daemon's law's decision (DecideByAssertion) as webauthn:mark-mac. A file
// that is not an assertion is refused before the POST; no passkey enrolled is refused.
func TestApproveDecideWebAuthnLinkAndAssertion(t *testing.T) {
	f := newFakeDecider(t)
	env := approveEnv(nil, nil, atTime(1710))
	f.enrolDecisionPasskey(t, false)
	if code, _, stderr := f.keyRun(t, env, "--webauthn-link"); code != exitUsage || !strings.Contains(stderr, "no passkey approver is enrolled") {
		t.Fatalf("no passkey: rc %d %s", code, stderr)
	}
	f.enrolDecisionPasskey(t, true)
	code, stdout, stderr := f.keyRun(t, env, "--deny", "--webauthn-link")
	m := regexp.MustCompile(`(?m)^  (https://\S+/surface/passkey#(\S+))$`).FindStringSubmatch(stdout)
	if code != exitOK || m == nil || len(f.posts) != 0 || !strings.Contains(stdout, "DENY request "+f.ref) ||
		!strings.Contains(stdout, f.ref+" --deny --webauthn-assertion <the saved file> --expiry 1760 --daemon "+f.addr()) {
		t.Fatalf("link: rc %d, %d POSTs\nstdout %s\nstderr %s", code, len(f.posts), stdout, stderr)
	}
	raw, err := base64.RawURLEncoding.DecodeString(m[2])
	if err != nil {
		t.Fatal(err)
	}
	var l ceremonyLink
	if err := json.Unmarshal(raw, &l); err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbDeny, 1760))
	ref16 := strings.TrimPrefix(f.ref, "sha256:")[:16]
	if l.Op != "get" || l.Origin != cdOrigin || l.RP != cdRP || l.Challenge != base64.RawURLEncoding.EncodeToString(h[:]) || l.Expires != 1760 ||
		l.Request != "decide-deny-"+ref16 || len(l.Allow) != 1 || l.Allow[0] != base64.RawURLEncoding.EncodeToString(cdPasskey.CredentialID()) ||
		!strings.Contains(stdout, "message hash "+l.Hash16) || l.Hash16 != base64Hex(h[:8]) {
		t.Fatalf("link %+v\nstdout %s", l, stdout)
	}
	asrPath := filepath.Join(t.TempDir(), "world-assertion-"+l.Request+".json")
	if err := os.WriteFile(asrPath, append(cdPasskey.Assert(cdOrigin, cdRP, h, 2), '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := f.keyRun(t, env, "--deny", "--webauthn-assertion", asrPath, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 {
		t.Fatalf("assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
	b := f.postBodies[0]
	if b.WebAuthn == nil || b.SSHSIG != nil || b.Typed || b.Verb != "deny" || b.Expiry != 1760 || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
		t.Fatalf("POST %+v", b)
	}
	k := cdPasskey.Key()
	got, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbDeny, Packet: refDigest(f.subject.Packet), Request: refDigest(f.ref),
		Expiry: b.Expiry, Now: 1710, Response: *b.WebAuthn,
		Approver: approval.Approver{Name: "mark-mac", Format: approval.FormatWebAuthn, EnrolledAt: 1000,
			Passkey: &approval.Passkey{RPID: cdRP, Origin: cdOrigin, CredentialID: cdPasskey.CredentialID(), Key: k, SignCount: 1}}})
	if err != nil || got.DecidedBy() != "webauthn:mark-mac" {
		t.Fatalf("the posted assertion does not decide: %v", err)
	}
	notAsr := filepath.Join(t.TempDir(), "x.json")
	if err := os.WriteFile(notAsr, []byte(`{"id":"x"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := f.keyRun(t, env, "--webauthn-assertion", notAsr, "--expiry", "1760"); code != exitUsage || len(f.posts) != 1 || !strings.Contains(stderr, "is not a passkey assertion") {
		t.Fatalf("not an assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
}

func base64Hex(b []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, 0, 2*len(b))
	for _, c := range b {
		out = append(out, digits[c>>4], digits[c&15])
	}
	return string(out)
}
