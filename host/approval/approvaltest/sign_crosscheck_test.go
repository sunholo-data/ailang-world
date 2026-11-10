package approvaltest_test

// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
// S3, XS2): the test signer is SSHSIG as the real ssh-keygen reads it, so every test that signs
// with it tests the format the operator's tool produces. Test keys only; no agent, temp HOME.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// keygen runs the real ssh-keygen with stdin, no agent and the test's HOME; it returns the
// combined output and the exit code. A run over 30 s or a failure to start is t.Fatal.
func keygen(t *testing.T, stdin []byte, args ...string) (string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")
	cmd.Stdin = bytes.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	switch {
	case err == nil:
		return string(out), 0
	case errors.As(err, &ee):
		return string(out), ee.ExitCode()
	default:
		t.Fatalf("ssh-keygen %v: %v", args, err)
		return "", 0
	}
}

// XS2: ssh-keygen -Y check-novalidate accepts the test signer's signature over an approval
// message in World's namespace and refuses it over a message one byte different; ssh-keygen -l
// prints the signer's own fingerprint for its AuthorizedKey. With WORLD_REQUIRE_SSH_KEYGEN=1 a
// missing ssh-keygen is a failure, not a skip (CI).
func TestTestSignerAgreesWithSSHKeygen(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
		}
		t.Skip("SKIP COUNTED: ssh-keygen not found; the test-signer cross-check did NOT run")
	}
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	s := approvaltest.NewSigner(7)
	msg, err := approval.SignedMessage([32]byte(bytes.Repeat([]byte{0x11}, 32)), [32]byte(bytes.Repeat([]byte{0x22}, 32)), 1800000600)
	if err != nil {
		t.Fatal(err)
	}
	sig := filepath.Join(dir, "msg.sig")
	pub := filepath.Join(dir, "k.pub")
	if err := os.WriteFile(sig, s.Sign(approval.Namespace, msg), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pub, s.AuthorizedKey(), 0o600); err != nil {
		t.Fatal(err)
	}
	out, rc := keygen(t, msg, "-Y", "check-novalidate", "-n", approval.Namespace, "-s", sig)
	if want := `Good "` + approval.Namespace + `" signature`; rc != 0 || !strings.Contains(out, want) {
		t.Fatalf("check-novalidate of the message: rc %d, output %q; want 0 and %q", rc, out, want)
	}
	flipped := bytes.Clone(msg)
	flipped[len(flipped)-1] ^= 1
	if out, rc := keygen(t, flipped, "-Y", "check-novalidate", "-n", approval.Namespace, "-s", sig); rc != 255 {
		t.Fatalf("check-novalidate of a flipped message: rc %d (%q), want 255", rc, out)
	}
	if out, rc := keygen(t, nil, "-l", "-f", pub); rc != 0 || !strings.Contains(out, s.PublicKey().Fingerprint()) {
		t.Fatalf("ssh-keygen -l: rc %d, output %q; want 0 and %s", rc, out, s.PublicKey().Fingerprint())
	}
}
