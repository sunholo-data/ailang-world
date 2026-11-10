package main

// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
// S2, S3): approve --sign-with. Every store is under t.TempDir(); the fake signer signs with
// approvaltest's fixed seeds, and the real-tool tests make their keys with ssh-keygen in
// t.TempDir() under a temp HOME with no agent (D-WORLD-109).

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// fakeSigner is a sign seam: it records what approve passed and returns ret's answer.
type fakeSigner struct {
	calls int
	key   string
	msg   []byte
	left  time.Duration // the context's deadline, measured from the call
	ret   func(msg []byte) ([]byte, error)
}

func (f *fakeSigner) sign(ctx context.Context, key string, msg []byte, _ io.Writer) ([]byte, error) {
	f.calls++
	f.key, f.msg = key, bytes.Clone(msg)
	if d, ok := ctx.Deadline(); ok {
		f.left = time.Until(d)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return f.ret(msg)
}

// signEnv is approveEnv with f as the sign seam (nil f: no seam at all).
func signEnv(f *fakeSigner, opened *bool, clock func() int64) sessionEnv {
	env := approveEnv(nil, opened, clock)
	if f != nil {
		env.sign = f.sign
	}
	return env
}

// countDaemonCalls counts fetchDaemonHealth calls for the rest of the test.
func countDaemonCalls(t *testing.T) *int {
	t.Helper()
	calls := 0
	inner := fetchDaemonHealth
	t.Cleanup(func() { fetchDaemonHealth = inner })
	fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) { calls++; return inner(addr) }
	return &calls
}

// CR8: --sign-with signs the request's exact message through the seam and judges the output as
// --signature does: no seam is an internal error (F7), a signer error is refused (F8), another
// approver's valid signature over another message and garbage are refused (F14, F11), and the
// right signature mints as ssh:mark-mac with no terminal.
func TestApproveSignWithFakeSigner(t *testing.T) {
	for _, arm := range []struct {
		name string
		ret  func(msg []byte) ([]byte, error) // nil: no sign seam
		code int
		want string // exact stderr when it ends in "\n", else a fragment
	}{
		{name: "no-signer", code: exitFatal, want: "ailang-worldd approve: internal error: no signer (sessionEnv.sign)\n"},
		{name: "signer-error", ret: func([]byte) ([]byte, error) { return nil, errors.New("agent refused") }, code: exitUsage,
			want: "ailang-worldd approve: --sign-with /keys/k: agent refused; nothing was approved\n"},
		{name: "other-approver-other-message", ret: func([]byte) ([]byte, error) {
			return otherSigner.Sign(approval.Namespace, []byte("another message")), nil
		}, code: exitUsage, want: "the signature does not approve request "},
		{name: "garbage", ret: func([]byte) ([]byte, error) { return []byte("not a signature\n"), nil }, code: exitUsage,
			want: "ssh-keygen's output for --sign-with /keys/k is not an SSH signature: approval: sshsig: bad armor"},
		{name: "good", ret: func(msg []byte) ([]byte, error) { return sigSigner.Sign(approval.Namespace, msg), nil }, code: exitOK},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a, _ := signedRig(t)
			seedApprover(t, a.credPath, "other", otherSigner.PublicKey(), approveT0-100)
			calls := countDaemonCalls(t)
			var f *fakeSigner
			if arm.ret != nil {
				f = &fakeSigner{ret: arm.ret}
			}
			var opened bool
			var stdout, stderr bytes.Buffer
			code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, signEnv(f, &opened, atTime(approveT0+1)))
			if code != arm.code || opened {
				t.Fatalf("rc %d, terminal opened %v; want %d and no terminal\nstderr %s", code, opened, arm.code, stderr.String())
			}
			row := a.row(t)
			if arm.code != exitOK {
				if strings.HasSuffix(arm.want, "\n") && stderr.String() != arm.want || !strings.Contains(stderr.String(), arm.want) {
					t.Fatalf("stderr %q, want %q", stderr.String(), arm.want)
				}
				if row.Consumed || len(a.sessionRows(t)) != 0 || *calls != 0 {
					t.Fatalf("consumed %v, daemon calls %d; want pending, no session, no network call", row.Consumed, *calls)
				}
				return
			}
			if !row.Consumed || row.ApprovedBy != "ssh:mark-mac" {
				t.Fatalf("consumed %v by %q, want ssh:mark-mac", row.Consumed, row.ApprovedBy)
			}
			if f.calls != 1 || f.key != "/keys/k" || !bytes.Equal(f.msg, requestMessage(t, a)) {
				t.Fatalf("the seam: %d call(s), key %q, msg %x; want 1, /keys/k and the request's message", f.calls, f.key, f.msg)
			}
		})
	}
}

// CR9: the signer's context ends when the request expires: 37 s left gives a 37 s deadline, 1 s
// gives 1 s, and a request that expired after admission is refused with the standard text before
// the signer is called (F7a).
func TestApproveSignWithDeadline(t *testing.T) {
	stop := func([]byte) ([]byte, error) { return nil, errors.New("stop") }
	for _, arm := range []struct {
		name     string
		clock    func() int64
		calls    int
		min, max time.Duration
	}{
		{"37-seconds-left", atTime(approveT0 + 600 - 37), 1, 36 * time.Second, 37 * time.Second},
		{"1-second-left", clockSeq(1, approveT0+598, approveT0+599), 1, 1, time.Second},
		{"expired-since-admission", clockSeq(1, approveT0+599, approveT0+600), 0, 0, 0},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a, _ := signedRig(t)
			f := &fakeSigner{ret: stop}
			var stdout, stderr bytes.Buffer
			if code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, signEnv(f, nil, arm.clock)); code != exitUsage {
				t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
			}
			if f.calls != arm.calls {
				t.Fatalf("the signer was called %d time(s), want %d\nstderr %s", f.calls, arm.calls, stderr.String())
			}
			if arm.calls == 0 {
				want := "ailang-worldd approve: request " + a.id + " expired at 2027-01-15T08:10:00Z; fix: make a new request with episode start --request\n"
				if stderr.String() != want {
					t.Fatalf("stderr %q, want %q", stderr.String(), want)
				}
			} else if f.left < arm.min || f.left > arm.max {
				t.Fatalf("the signer's deadline was %s ahead, want within [%s, %s]", f.left, arm.min, arm.max)
			}
			if a.row(t).Consumed {
				t.Fatal("the request was consumed")
			}
		})
	}
}
