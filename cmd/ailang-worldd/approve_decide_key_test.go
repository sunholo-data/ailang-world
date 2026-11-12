package main

// Row 201 A4 sprint 2 M5b (design §2.10, §6 CD5; CD6-CD9 are sprint 3's, plan S2-D1): approve
// <request-ref> --sign-with against the fake daemon (approve_decide_rig_test.go): the 91-byte
// decision message through the signer seam, and the real ssh-keygen on a key made in a temp dir
// (D-WORLD-109: no real key, no agent).

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// keyRun is approve <ref> --db --daemon extra... with env.
func (f *fakeDecider) keyRun(t *testing.T, env sessionEnv, extra ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runApprove(append([]string{f.ref, "--db", f.db, "--daemon", f.addr()}, extra...), &stdout, &stderr, env)
	return code, stdout.String(), stderr.String()
}

// decisionMessageOf is the message deciding verb on f's subject at expiry.
func (f *fakeDecider) decisionMessageOf(t *testing.T, verb byte, expiry int64) []byte {
	t.Helper()
	m, err := approval.DecisionMessage(verb, refDigest(f.subject.Packet), refDigest(f.ref), expiry)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// CD5: --sign-with signs the 91-byte decision message for the window's max through the decision
// signer, bounded by that expiry, and posts its output as sshsig (no terminal, no proof); a signer
// error posts nothing; the production env's decision signer is signDecisionWithSSHKeygen.
func TestApproveDecideSignWith(t *testing.T) {
	f := newFakeDecider(t)
	fs := &fakeSigner{ret: func([]byte) ([]byte, error) { return []byte("ARMORED"), nil }}
	opened := false
	env := approveEnv(nil, &opened, atTime(1710))
	env.signDecision = fs.sign
	code, stdout, stderr := f.keyRun(t, env, "--deny", "--sign-with", "/k/id")
	if want := f.decisionMessageOf(t, approval.VerbDeny, 1760); code != exitOK || opened || fs.key != "/k/id" || !bytes.Equal(fs.msg, want) || len(fs.msg) != 91 ||
		fs.left <= 0 || fs.left > 50*time.Second || len(f.posts) != 1 {
		t.Fatalf("rc %d, opened %v, key %q, msg %x (left %v), %d POSTs\nstderr %s", code, opened, fs.key, fs.msg, fs.left, len(f.posts), stderr)
	}
	if b := f.postBodies[0]; b.Typed || b.SSHSIG == nil || *b.SSHSIG != "ARMORED" || b.Expiry != 1760 || b.Verb != "deny" || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
		t.Fatalf("POST %+v", b)
	}
	if _, wrote := f.proofFile(t); wrote || !strings.Contains(stdout, "decided deny") {
		t.Fatalf("proof written %v; stdout %q", wrote, stdout)
	}
	g := newFakeDecider(t)
	env.signDecision = (&fakeSigner{ret: func([]byte) ([]byte, error) { return nil, errors.New("agent refused") }}).sign
	if code, _, stderr := g.keyRun(t, env, "--sign-with", "/k/id"); code != exitUsage || len(g.posts) != 0 || !strings.Contains(stderr, "--sign-with /k/id: agent refused; nothing was decided") {
		t.Fatalf("signer error: rc %d, %d POSTs\nstderr %s", code, len(g.posts), stderr)
	}
	if got, want := reflect.ValueOf(realSessionEnv().signDecision).Pointer(), reflect.ValueOf(signDecisionWithSSHKeygen).Pointer(); got != want {
		t.Fatal("realSessionEnv().signDecision is not signDecisionWithSSHKeygen")
	}
	for _, flag := range []string{"--signature", "--write-message"} { // not a request ref's flags in this sprint
		if code, _, stderr := g.keyRun(t, env, flag, "/x"); code != exitUsage || g.gets != 1 || !strings.Contains(stderr, "the other key flags do not decide it") {
			t.Errorf("%s: rc %d, GETs %d\nstderr %s", flag, code, g.gets, stderr)
		}
	}
}

// CD5 with the real tool: the production decision signer runs the real ssh-keygen on a key made
// in a temp dir (no agent, temp HOME); what it posts verifies with the daemon's own check
// (VerifyDecisionSSHSIG) over the message for the window's max, and not over the mint namespace.
// With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen is a failure, not a skip (CI).
func TestApproveDecideSignWithSSHKeygen(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
		}
		t.Skip("SKIP COUNTED: ssh-keygen not found; the decision --sign-with end-to-end test did NOT run")
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SSH_AUTH_SOCK", "")
	key := filepath.Join(t.TempDir(), "k")
	sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "cd5", "-f", key)
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
	env.signDecision = signDecisionWithSSHKeygen
	if code, _, stderr := f.keyRun(t, env, "--sign-with", key); code != exitOK || len(f.posts) != 1 {
		t.Fatalf("rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
	armored := []byte(*f.postBodies[0].SSHSIG)
	if err := approval.VerifyDecisionSSHSIG(armored, approval.VerbApprove, refDigest(f.subject.Packet), refDigest(f.ref), 1760, k); err != nil {
		t.Fatalf("the real signature does not decide: %v", err)
	}
	if err := approval.VerifyApproval(armored, refDigest(f.ref), refDigest(f.subject.Packet), 1760, k); !errors.Is(err, approval.ErrNamespace) {
		t.Fatalf("checked as a mint: %v, want ErrNamespace", err)
	}
}
