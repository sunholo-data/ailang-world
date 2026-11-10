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
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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

// ---- S3 ----

// fakeSSHKeygen puts an executable `ssh-keygen` running body (sh) first on PATH for the rest of the
// test, and returns its directory. The real tool is never run by it.
func fakeSSHKeygen(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ssh-keygen"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return dir
}

// G1-G3 and the exec line, with a fake ssh-keygen first on PATH: the production signer runs
// `ssh-keygen -Y sign -n world-approval@sunholo -f <key>` with the message on stdin and its stderr
// passed through, and refuses a failed run (G1), over 16 KiB of output (G2), no output (G3) and a
// signer still running when the request expires. Each is exitUsage, pending.
func TestSignWithSSHKeygenRefusesBadOutput(t *testing.T) {
	for _, arm := range []struct {
		name, body string
		now        int64
		want       []string
	}{
		{"exec-line-and-no-output", "printf '%s\\n' \"$@\" > \"$(dirname \"$0\")/args\"\ncat > \"$(dirname \"$0\")/stdin\"\n", approveT0 + 1,
			[]string{"--sign-with /keys/k: ssh-keygen -Y sign printed nothing; nothing was approved"}},
		{"failed", "cat > /dev/null\necho 'No private key found' >&2\nexit 255\n", approveT0 + 1,
			[]string{"No private key found\n", "--sign-with /keys/k: ssh-keygen -Y sign: exit status 255; nothing was approved"}},
		{"over-16k", "cat > /dev/null\nhead -c 16385 /dev/zero | tr '\\000' a\n", approveT0 + 1,
			[]string{"--sign-with /keys/k: ssh-keygen -Y sign printed over 16384 bytes; nothing was approved"}},
		{"outlives-the-request", "cat > /dev/null\nsleep 20\n", approveT0 + 599, // 1 s left
			[]string{"--sign-with /keys/k: ssh-keygen -Y sign: signal: killed; nothing was approved"}},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a, _ := signedRig(t)
			dir := fakeSSHKeygen(t, arm.body)
			env := approveEnv(nil, nil, atTime(arm.now))
			env.sign = signWithSSHKeygen
			var stdout, stderr bytes.Buffer
			start := time.Now()
			if code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, env); code != exitUsage {
				t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
			}
			if took := time.Since(start); took > 10*time.Second {
				t.Fatalf("approve took %s: the signer outlived the request's deadline", took)
			}
			for _, want := range arm.want {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("stderr %q lacks %q", stderr.String(), want)
				}
			}
			if a.row(t).Consumed {
				t.Fatal("the request was consumed")
			}
			if arm.name != "exec-line-and-no-output" {
				return
			}
			args, err := os.ReadFile(filepath.Join(dir, "args"))
			if want := "-Y\nsign\n-n\nworld-approval@sunholo\n-f\n/keys/k\n"; err != nil || string(args) != want {
				t.Fatalf("ssh-keygen's arguments %q (%v), want %q", args, err, want)
			}
			if in, err := os.ReadFile(filepath.Join(dir, "stdin")); err != nil || !bytes.Equal(in, requestMessage(t, a)) {
				t.Fatalf("ssh-keygen's stdin %x (%v), want the request's message", in, err)
			}
		})
	}
}

// cappedBuffer holds at most its bound, however much is written, and records the overflow (G2).
func TestCappedBufferKeepsItsBound(t *testing.T) {
	c := &cappedBuffer{max: 4}
	for _, w := range []string{"ab", "cdef", "g"} {
		if n, err := c.Write([]byte(w)); n != len(w) || err != nil {
			t.Fatalf("Write(%q) = (%d, %v), want (%d, nil)", w, n, err, len(w))
		}
	}
	if c.buf.String() != "abcd" || !c.over {
		t.Fatalf("buffer %q, over %v; want \"abcd\", true", c.buf.String(), c.over)
	}
}

// The production env signs with signWithSSHKeygen: without it, --sign-with is the F7 internal error.
func TestRealSessionEnvSignsWithSSHKeygen(t *testing.T) {
	if got, want := reflect.ValueOf(realSessionEnv().sign).Pointer(), reflect.ValueOf(signWithSSHKeygen).Pointer(); got != want {
		t.Fatalf("realSessionEnv().sign is %#x, want signWithSSHKeygen (%#x)", got, want)
	}
}

// sshKeygen runs the real ssh-keygen with no agent (SSH_AUTH_SOCK emptied) and stdin from
// /dev/null, under the test's temp HOME; a run over 30 s or a failure to start is t.Fatal.
func sshKeygen(t *testing.T, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("ssh-keygen %v: %v\n%s", args, err, out)
	}
}

// XS1: approve --sign-with with the production signer and the REAL ssh-keygen, for each key type:
// a key made in t.TempDir() under a temp HOME, with no agent, signs and mints as ssh:mark-mac; its
// .pub with the private file moved away and no agent is refused, naming ssh-keygen's failure (G1).
// With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen is a failure, not a skip (CI).
func TestSignWithSSHKeygenEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
		}
		t.Skip("SKIP COUNTED: ssh-keygen not found; the --sign-with end-to-end test did NOT run")
	}
	for _, kt := range []struct {
		name string
		args []string
	}{{"ed25519", []string{"-t", "ed25519"}}, {"ecdsa-p256", []string{"-t", "ecdsa", "-b", "256"}}} {
		t.Run(kt.name, func(t *testing.T) {
			a := newApproveRig(t, "ep1", approveT0) // its episode run sets HOME to a temp dir
			t.Setenv("SSH_AUTH_SOCK", "")
			key := filepath.Join(t.TempDir(), "k")
			sshKeygen(t, append(append([]string{"-q"}, kt.args...), "-N", "", "-C", "xs1", "-f", key)...)
			pub, err := os.ReadFile(key + ".pub")
			if err != nil {
				t.Fatal(err)
			}
			k, err := approval.ParseAuthorizedKey(pub)
			if err != nil {
				t.Fatal(err)
			}
			seedApprover(t, a.credPath, "mark-mac", k, approveT0-100)
			env := approveEnv(nil, nil, atTime(approveT0+1))
			env.sign = signWithSSHKeygen
			var stdout, stderr bytes.Buffer
			if code := runApprove(a.args("--sign-with", key), &stdout, &stderr, env); code != exitOK {
				t.Fatalf("approve --sign-with %s: rc %d, want exitOK\nstderr %s", key, code, stderr.String())
			}
			if row := a.row(t); row.ApprovedBy != "ssh:mark-mac" {
				t.Fatalf("approved_by %q, want ssh:mark-mac", row.ApprovedBy)
			}

			b := newApproveRigIn(t, a.r, "ep2", approveT0)
			if err := os.Rename(key, key+".moved"); err != nil {
				t.Fatal(err)
			}
			stderr.Reset()
			if code := runApprove(b.args("--sign-with", key+".pub"), &stdout, &stderr, env); code != exitUsage {
				t.Fatalf("approve --sign-with k.pub, no private key and no agent: rc %d, want exitUsage\nstderr %s", code, stderr.String())
			}
			for _, want := range []string{"No private key found", "ssh-keygen -Y sign: exit status 255; nothing was approved"} {
				if !strings.Contains(stderr.String(), want) {
					t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
				}
			}
			if b.row(t).Consumed {
				t.Fatal("the request was consumed with no key to sign it")
			}
		})
	}
}
