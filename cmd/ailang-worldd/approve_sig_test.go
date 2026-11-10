package main

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M5a-M5c): approve binds the stored challenge, and approves by an enrolled key's signature.
// Every store is under t.TempDir(); every key is approvaltest.NewSigner's or a checked-in fixture
// (D-WORLD-109).

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// CR11 (M4c2): a stored challenge that is not 64 hex is refused before any terminal opens (F4), and
// the request stays pending.
func TestApproveRefusesANonHexChallenge(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	raw := rawCredentials(t, a.credPath)
	if _, err := raw.Exec(`UPDATE mint_requests SET challenge = ? WHERE request_id = ?`, strings.Repeat("zz", 32), a.id); err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()
	var opened bool
	var stdout, stderr bytes.Buffer
	code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, &opened, atTime(approveT0+1)))
	if code != exitUsage || opened {
		t.Fatalf("rc %d, terminal opened %v; want exitUsage and no terminal\nstderr %s", code, opened, stderr.String())
	}
	want := "ailang-worldd approve: request " + a.id + "'s stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request\n"
	if stderr.String() != want {
		t.Fatalf("stderr %q, want %q", stderr.String(), want)
	}
	if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
		t.Fatal("a request with a non-hex challenge was consumed or minted")
	}
}

// ---- M5a ----

// sigSigner is the enrolled approver's key in the approve --signature tests; otherSigner is a key
// that is never enrolled. Both are approvaltest's fixed seeds (D-WORLD-109).
var (
	sigSigner   = approvaltest.NewSigner(11)
	otherSigner = approvaltest.NewSigner(12)
)

// requestMessage is the 90-byte message the approval of the rig's stored row signs.
func requestMessage(t *testing.T, a approveRig) []byte {
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
	msg, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return msg
}

// writeSignature writes sig to a new file under t.TempDir() and returns its path.
func writeSignature(t *testing.T, sig []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "request.sig")
	if err := os.WriteFile(p, sig, 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// signedRig is a rig for ep1 at approveT0 with sigSigner enrolled as mark-mac at approveT0-100, and
// the path of sigSigner's signature over the request's message.
func signedRig(t *testing.T) (approveRig, string) {
	t.Helper()
	a := newApproveRig(t, "ep1", approveT0)
	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
	return a, writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, a)))
}

// CR2 (M5a): a verified signature mints the stored request with no terminal: approved_by is
// ssh:mark-mac, the approver's last_used_at is the consume's reading, stdout shows the request and
// its approver, and every other effect is what the typed y produces (CL16).
func TestApproveBySignatureMintsWithNoTerminal(t *testing.T) {
	a, sig := signedRig(t)
	var opened bool
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("rc %d, want exitOK\nstderr %s", code, stderr.String())
	}
	if opened {
		t.Fatal("a terminal was opened on the signature path")
	}
	row := a.row(t)
	if !row.Consumed || row.ApprovedBy != "ssh:mark-mac" {
		t.Fatalf("request after the signed approve: consumed %v by %q, want ssh:mark-mac", row.Consumed, row.ApprovedBy)
	}
	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != approveT0+1 {
		t.Fatalf("approver rows %+v, want mark-mac last used at %d", rows, approveT0+1)
	}
	out := stdout.String()
	for _, want := range []string{"mint request " + a.id + "\n", "expiry +3600s: approved by ssh:mark-mac\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout lacks %q:\n%s", want, out)
		}
	}
	mr, err := approval.Decode(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	var probs []string
	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, out)) {
		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty; this path records ssh:mark-mac (asserted above)
			probs = append(probs, p)
		}
	}
	if len(probs) != 0 {
		t.Errorf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
	}
}
