package main

// Row 201 A3 sprint 1 M3a (design_docs/planned/w-world-approvals-by-key-a3.md §2.6, §6 AW1, AW5):
// approve --await, the minter. Every store is under t.TempDir(); every key is approvaltest's;
// the surface's recording is the store's RecordSurfaceApproval, called by the test mid-wait
// (D-WORLD-109: no daemon route writes the table in this sprint).

import (
	"bytes"
	"context"
	"encoding/hex"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// awaitClock is approve --await's fake clock: env.after advances now by its duration and counts
// the waits, then calls onWait(waits, now) before it fires.
type awaitClock struct {
	now, waits int64
	onWait     func(n, now int64)
}

// awaitEnv is the minter's env at start: the fake clock, and a terminal that records into *opened
// and would answer y, so any path that reaches the typed y is seen (AW1, AW7).
func awaitEnv(start int64, opened *bool, onWait func(n, now int64)) (sessionEnv, *awaitClock) {
	c := &awaitClock{now: start, onWait: onWait}
	env := approveEnv(&promptTerm{in: strings.NewReader("y\n")}, opened, func() int64 { return c.now })
	env.after = func(d time.Duration) <-chan time.Time {
		c.waits++
		c.now += int64(d / time.Second)
		if c.onWait != nil {
			c.onWait(c.waits, c.now)
		}
		ch := make(chan time.Time, 1)
		ch <- time.Unix(c.now, 0)
		return ch
	}
	return env, c
}

// recordSurface records artifact for the rig's request as the surface would, at now.
func recordSurface(t *testing.T, a approveRig, format, approver string, artifact []byte, now int64) {
	t.Helper()
	st, err := store.OpenCredentials(a.credPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: a.id, Format: format, Approver: approver, Artifact: artifact, Digest: a.row(t).Digest}, now); err != nil {
		t.Error(err)
	}
}

// requireMintedBy fails unless the rig's request was consumed by approvedBy with the effects the
// typed y makes (CL16), and stdout shows the approver.
func requireMintedBy(t *testing.T, a approveRig, approvedBy, stdout string) {
	t.Helper()
	row := a.row(t)
	if !row.Consumed || row.ApprovedBy != approvedBy || len(a.sessionRows(t)) != 1 {
		t.Fatalf("request consumed %v by %q, %d sessions; want one, by %s", row.Consumed, row.ApprovedBy, len(a.sessionRows(t)), approvedBy)
	}
	if !strings.Contains(stdout, "approved by "+approvedBy+"\n") {
		t.Fatalf("stdout lacks the approver:\n%s", stdout)
	}
	mr, err := approval.Decode(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	var probs []string
	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, stdout)) {
		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty
			probs = append(probs, p)
		}
	}
	if len(probs) != 0 {
		t.Fatalf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
	}
}

// AW1: a passkey assertion recorded mid-wait (after the third wait) is judged and mints with no
// terminal: one session, approved_by webauthn:mark-mac, the counter advanced, the token in the
// request's session file; the minter printed the request's message hash and the card's address
// once for the two passkeys' one origin, which is also the request's loopback daemon, and polled
// exactly four times.
func TestApproveAwaitMintsARecordedPasskeyApproval(t *testing.T) {
	a := passkeyRig(t)
	seedPasskey(t, a.credPath, "b-other", approvaltest.NewPasskey(approval.AlgES256, 52), approveT0-100)
	h := requestHashOf(t, a)
	var opened bool
	env, c := awaitEnv(approveT0+1, &opened, func(n, now int64) {
		if n == 3 {
			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 2), now)
		}
	})
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened {
		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
	}
	if c.waits != 3 {
		t.Fatalf("%d waits, want 3 (the fourth poll finds the approval)", c.waits)
	}
	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
	if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 2 {
		t.Fatalf("sign count %d, want 2", p.SignCount)
	}
	mr, err := approval.Decode(a.row(t).Request)
	if err != nil {
		t.Fatal(err)
	}
	_, port, err := net.SplitHostPort(mr.Daemon.Addr)
	if err != nil {
		t.Fatal(err)
	}
	if pkOrigin != "http://localhost:"+port {
		t.Fatalf("the rig's passkey origin %s is not the request's loopback daemon %s", pkOrigin, mr.Daemon.Addr)
	}
	for _, want := range []string{"waiting for an approval of request " + a.id + " (expires " + approverTime(a.row(t).ExpiresAt) + ", message hash " + hex.EncodeToString(h[:8]) + ") from the surface",
		"\n  open " + pkOrigin + "/surface/approve/" + a.id + "\n"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout lacks %q:\n%s", want, stdout.String())
		}
	}
	if n := strings.Count(stdout.String(), "/surface/approve/"); n != 1 {
		t.Fatalf("%d card addresses, want 1 (one origin):\n%s", n, stdout.String())
	}
	if strings.Contains(stderr.String(), "refused the surface's approval") {
		t.Fatalf("a verified approval was refused:\n%s", stderr.String())
	}
}

// AW5: --await is a sixth exclusive mode (F1); --await-for needs --await and at least 1s; --list
// refuses --await (F2). Each prints the help and exits exitUsage before any store is read.
func TestApproveAwaitFlagRules(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	for _, extra := range [][]string{
		{"--await", "--signature", "x.sig"}, {"--await", "--sign-with", "k"}, {"--await", "--write-message", "m"},
		{"--await", "--webauthn-link"}, {"--await", "--webauthn-assertion", "a.json"},
		{"--await-for", "30s"}, {"--await", "--await-for", "0s"}, {"--await", "--await-for", "999ms"}, {"--await", "--await-for", "-5s"},
	} {
		var opened bool
		var stdout, stderr bytes.Buffer
		env, _ := awaitEnv(approveT0+1, &opened, nil) // a fake clock: a rule that let a mode through waits at most to the expiry
		if code := runApprove(a.args(extra...), &stdout, &stderr, env); code != exitUsage || opened ||
			!strings.HasPrefix(stderr.String(), "usage: ailang-worldd approve") {
			t.Errorf("%v: rc %d, terminal %v, stderr %q; want exitUsage and the help", extra, code, opened, stderr.String())
		}
	}
	var stdout, stderr bytes.Buffer
	if code := runApprove([]string{"--list", "--await", "--db", a.r.f.db}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitUsage ||
		!strings.HasPrefix(stderr.String(), "usage: ailang-worldd approve") {
		t.Errorf("--list --await: rc %d, stderr %q; want exitUsage and the help", code, stderr.String())
	}
	if !strings.Contains(approveHelp, "ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]\n") {
		t.Error("the help lacks the --await usage line")
	}
	if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
		t.Fatal("a refused flag combination consumed the request or minted")
	}
	// Control: --await-for 1s with --await is accepted (no approval is recorded, so it waits one
	// second and stops, naming the time).
	env, c := awaitEnv(approveT0+1, nil, nil)
	stderr.Reset()
	if code := runApprove(a.args("--await", "--await-for", "1s"), &stdout, &stderr, env); code != exitUsage || c.waits != 1 ||
		!strings.Contains(stderr.String(), "no approval for request "+a.id+" by "+approverTime(approveT0+2)+"; it stays pending until ") {
		t.Fatalf("control: rc %d, %d waits, stderr %q", code, c.waits, stderr.String())
	}
}
