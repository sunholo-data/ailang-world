package main

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M5a-M5c): approve binds the stored challenge, and approves by an enrolled key's signature.
// Every store is under t.TempDir(); every key is approvaltest.NewSigner's or a checked-in fixture
// (D-WORLD-109).

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/daemon"
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

// ---- M5b ----

// CR3 (M5b): with an approver enrolled, the typed path still opens the terminal, records tty and
// leaves the approver unused; --list refuses --signature (F2).
func TestApproveTypedPathUnchangedWithAnApprover(t *testing.T) {
	a, sig := signedRig(t)
	var opened bool
	var stdout, stderr bytes.Buffer
	term := &promptTerm{in: strings.NewReader("y\n")}
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, &opened, atTime(approveT0+1))); code != exitOK || !opened {
		t.Fatalf("typed approve: rc %d, terminal opened %v; want exitOK through the terminal\nstderr %s", code, opened, stderr.String())
	}
	if row := a.row(t); row.ApprovedBy != "tty" {
		t.Fatalf("approved_by %q, want tty", row.ApprovedBy)
	}
	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != 0 {
		t.Fatalf("approver rows %+v, want mark-mac never used", rows)
	}
	stderr.Reset()
	listArgs := []string{"--list", "--signature", sig, "--db", a.r.f.db, "--credentials-db", a.credPath}
	if code := runApprove(listArgs, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitUsage || stderr.String() != approveHelp {
		t.Fatalf("--list --signature: rc %d, stderr %q; want exitUsage and the help", code, stderr.String())
	}
}

// CR4, CR5 (M5b): each signature refusal exits exitUsage with its text, before any terminal and
// before any network call, and leaves the request as it was. A consumed or expired request is
// refused by admission before the signature file is read (the path does not exist).
func TestApproveBySignatureRefusalsBeforeAnyEffect(t *testing.T) {
	good := func(t *testing.T, a approveRig) string {
		return writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, a)))
	}
	absent := func(t *testing.T, a approveRig) string { return filepath.Join(t.TempDir(), "absent.sig") }
	for _, arm := range []struct {
		name     string
		setup    func(t *testing.T, a approveRig)
		sig      func(t *testing.T, a approveRig) string
		now      int64 // 0: approveT0+1
		want     string
		consumed bool // the setup consumed the request (one session row)
	}{
		{name: "fifo", sig: func(t *testing.T, a approveRig) string {
			p := filepath.Join(t.TempDir(), "sig.fifo")
			if err := syscall.Mkfifo(p, 0o600); err != nil {
				t.Fatal(err)
			}
			return p
		}, want: "sig.fifo is not a regular file"},
		{name: "over-16k", sig: func(t *testing.T, a approveRig) string {
			return writeSignature(t, bytes.Repeat([]byte("a"), 16<<10+1))
		}, want: "is 16385 bytes, over the 16384-byte bound"},
		{name: "missing", sig: absent, want: "--signature: open "},
		{name: "garbage", sig: func(t *testing.T, a approveRig) string { return writeSignature(t, []byte("not a signature\n")) },
			want: "is not an SSH signature: approval: sshsig: bad armor"},
		{name: "unenrolled-key", sig: func(t *testing.T, a approveRig) string {
			return writeSignature(t, otherSigner.Sign(approval.Namespace, requestMessage(t, a)))
		}, want: "the signature's key " + otherSigner.PublicKey().Fingerprint() + " is not an enrolled approver; see ailang-worldd approver list"},
		{name: "revoked", setup: func(t *testing.T, a approveRig) { revokeSeeded(t, a.credPath, "mark-mac", approveT0-50) },
			sig: good, want: "approval: the approver is revoked"},
		{name: "other-request", sig: func(t *testing.T, a approveRig) string {
			msg := requestMessage(t, a)
			msg[len(msg)-9] ^= 1 // the challenge's last byte
			return writeSignature(t, sigSigner.Sign(approval.Namespace, msg))
		}, want: "approval: sshsig: signature does not verify"},
		{name: "enrolled-after-now", setup: func(t *testing.T, a approveRig) {
			raw := rawCredentials(t, a.credPath)
			defer func() { _ = raw.Close() }()
			if _, err := raw.Exec(`UPDATE approvers SET enrolled_at = ? WHERE name = 'mark-mac'`, approveT0+50); err != nil {
				t.Fatal(err)
			}
		}, sig: good, want: "approval: the clock is behind the approver's enrolment time"},
		{name: "git-namespace", sig: func(t *testing.T, a approveRig) string {
			return writeSignature(t, sigSigner.Sign("git", requestMessage(t, a)))
		}, want: "approval: sshsig: wrong namespace"},
		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
			raw := rawCredentials(t, a.credPath)
			defer func() { _ = raw.Close() }()
			for _, stmt := range []string{`DROP TABLE approvers`, `PRAGMA user_version = 4`} {
				if _, err := raw.Exec(stmt); err != nil {
					t.Fatal(err)
				}
			}
		}, sig: good, want: "approvers need schema 5"},
		{name: "consumed-first", setup: func(t *testing.T, a approveRig) {
			var stdout, stderr bytes.Buffer
			if code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+1))); code != exitOK {
				t.Fatalf("typed approve: rc %d\nstderr %s", code, stderr.String())
			}
		}, sig: absent, want: "was already used", consumed: true},
		{name: "expired-first", sig: absent, now: approveT0 + 600, want: "expired at"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a := newApproveRig(t, "ep1", approveT0)
			seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
			if arm.setup != nil {
				arm.setup(t, a)
			}
			calls := 0
			inner := fetchDaemonHealth
			t.Cleanup(func() { fetchDaemonHealth = inner })
			fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) { calls++; return inner(addr) }
			now := arm.now
			if now == 0 {
				now = approveT0 + 1
			}
			var opened bool
			var stdout, stderr bytes.Buffer
			code := runApprove(a.args("--signature", arm.sig(t, a)), &stdout, &stderr, approveEnv(nil, &opened, atTime(now)))
			if code != exitUsage || opened || calls != 0 {
				t.Fatalf("rc %d, terminal opened %v, daemon calls %d; want exitUsage, none, 0\nstderr %s", code, opened, calls, stderr.String())
			}
			if !strings.Contains(stderr.String(), arm.want) {
				t.Fatalf("stderr %q lacks %q", stderr.String(), arm.want)
			}
			if row := a.row(t); row.Consumed != arm.consumed || row.ApprovedBy == "ssh:mark-mac" {
				t.Fatalf("request consumed %v by %q, want consumed %v and never by ssh:mark-mac", row.Consumed, row.ApprovedBy, arm.consumed)
			}
			want := 0
			if arm.consumed {
				want = 1
			}
			if n := len(a.sessionRows(t)); n != want {
				t.Fatalf("%d session rows, want %d", n, want)
			}
		})
	}
}

// CR6 (M5b): a revoke that commits after approve verified the signature, inside the mint's own
// transaction (the consume's in-transaction clock reading: design V36/V37), refuses the mint at
// the write: exitFatal, the request pending, no session, the worktree this run made removed.
func TestApproveRevokeRaceIsRefusedAtTheWrite(t *testing.T) {
	a, sig := signedRig(t)
	calls := 0
	clock := func() int64 {
		calls++
		if calls == 3 { // 1 admission, 2 provisionEpisode, 3 the consume's reading
			raw := rawCredentials(t, a.credPath)
			defer func() { _ = raw.Close() }()
			if _, err := raw.Exec(`UPDATE approvers SET revoked_at = ? WHERE name = 'mark-mac'`, approveT0+1); err != nil {
				t.Errorf("the racing revoke: %v", err)
			}
		}
		return approveT0 + 1
	}
	var opened bool
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, clock)); code != exitFatal {
		t.Fatalf("rc %d, want exitFatal\nstderr %s", code, stderr.String())
	}
	for _, want := range []string{"store: the approver is revoked", "removed the worktree"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
		}
	}
	if calls != 3 || opened {
		t.Errorf("clock read %d times (want 3), terminal opened %v", calls, opened)
	}
	if row := a.row(t); row.Consumed {
		t.Error("the request was consumed by a revoked approver")
	}
	if n := len(a.sessionRows(t)); n != 0 {
		t.Errorf("%d session rows, want 0", n)
	}
	if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree exists (err %v)", err)
	}
	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != 0 || rows[0].RevokedAt != approveT0+1 {
		t.Errorf("approver rows %+v, want mark-mac revoked and never used", rows)
	}
}
