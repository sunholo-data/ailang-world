package approval

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
)

// runSSHKeygen runs ssh-keygen with the given arguments and stdin (nil for none). SSH_AUTH_SOCK is
// emptied so that no real agent is consulted (row 201 constraint 1). It returns stdout and the exit
// code; a timeout or a failure to start is t.Fatal.
func runSSHKeygen(t *testing.T, stdin []byte, args ...string) (stdout []byte, rc int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("ssh-keygen %v timed out after 30s", args)
	}
	if err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("ssh-keygen %v failed to run: %v", args, err)
		}
		rc = ee.ExitCode()
	}
	if rc != 0 {
		t.Logf("ssh-keygen %v exited %d: %s", args, rc, strings.TrimSpace(errBuf.String()))
	}
	return out.Bytes(), rc
}

// TestSSHSIGCrossCheckWithSSHKeygen cross-checks the SSHSIG verifier against real ssh-keygen (row 201
// test XC1). For each key type, real keys are made in t.TempDir(), a message is signed from stdin,
// and five arms are verified by both ssh-keygen -Y verify and VerifySSHSIG: they must agree.
func TestSSHSIGCrossCheckWithSSHKeygen(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
		}
		t.Skip("SKIP COUNTED: ssh-keygen not found; the SSHSIG cross-check did NOT run")
	}
	t.Run("ed25519", func(t *testing.T) {
		crossCheckSSHKeygen(t, []string{"-t", "ed25519"}, "ssh-ed25519")
	})
	t.Run("ecdsa-p256", func(t *testing.T) {
		crossCheckSSHKeygen(t, []string{"-t", "ecdsa", "-b", "256"}, "ecdsa-sha2-nistp256")
	})
}

// crossCheckSSHKeygen runs the XC1 steps for one key type. Arms are loop iterations, not subtests.
func crossCheckSSHKeygen(t *testing.T, keyArgs []string, wantType string) {
	dir := t.TempDir()
	key := filepath.Join(dir, "key")
	other := filepath.Join(dir, "other")

	// 1. Make the two keys (test keys only, no passphrase, in t.TempDir()).
	for _, p := range []struct{ name, path string }{{"xc1", key}, {"other", other}} {
		args := append(append([]string{"-q"}, keyArgs...), "-N", "", "-C", "xc1-"+p.name, "-f", p.path)
		if _, rc := runSSHKeygen(t, nil, args...); rc != 0 {
			t.Fatalf("ssh-keygen keygen %s: rc %d, want 0", p.name, rc)
		}
	}

	// 2. Sign the fixture message from stdin, under the world namespace, the git namespace and with
	// the other key. Each must give rc 0, a non-empty stdout and parse with the wanted key type.
	msg := fixtureMessage(t)
	good, goodRC := runSSHKeygen(t, msg, "-Y", "sign", "-n", Namespace, "-f", key)
	gitSig, gitRC := runSSHKeygen(t, msg, "-Y", "sign", "-n", "git", "-f", key)
	otherSig, otherRC := runSSHKeygen(t, msg, "-Y", "sign", "-n", Namespace, "-f", other)
	good256, rc256 := runSSHKeygen(t, msg, "-Y", "sign", "-n", Namespace, "-O", "hashalg=sha256", "-f", key)
	for name, sig := range map[string]struct {
		out []byte
		rc  int
	}{"good": {good, goodRC}, "git": {gitSig, gitRC}, "other": {otherSig, otherRC}, "good-sha256": {good256, rc256}} {
		if sig.rc != 0 {
			t.Fatalf("ssh-keygen -Y sign (%s): rc %d, want 0", name, sig.rc)
		}
		if len(sig.out) == 0 {
			t.Fatalf("ssh-keygen -Y sign (%s): empty stdout", name)
		}
		s, err := ParseSSHSIG(sig.out)
		if err != nil {
			t.Fatalf("ssh-keygen -Y sign (%s) does not parse: %v", name, err)
		}
		if s.PublicKey.Type != wantType {
			t.Fatalf("ssh-keygen -Y sign (%s): key type %q, want %q", name, s.PublicKey.Type, wantType)
		}
	}
	// The default signature is SHA-512 and the -O hashalg=sha256 one is SHA-256: the two arms below
	// are the only ones with a real SHA-256 signature.
	if s, err := ParseSSHSIG(good); err != nil || s.HashAlg != "sha512" {
		t.Fatalf("ssh-keygen default signature: err %v; want hash sha512", err)
	}
	if s, err := ParseSSHSIG(good256); err != nil || s.HashAlg != "sha256" {
		t.Fatalf("ssh-keygen -O hashalg=sha256 signature: err %v; want hash sha256", err)
	}

	// 3. Build the allowed-signers file for the key and parse the public key with the production parser.
	pubData, err := os.ReadFile(key + ".pub")
	if err != nil {
		t.Fatal(err)
	}
	pub, err := ParseAuthorizedKey(pubData)
	if err != nil {
		t.Fatalf("ParseAuthorizedKey(key.pub): %v", err)
	}
	if pub.Type != wantType {
		t.Fatalf("key.pub type %q, want %q", pub.Type, wantType)
	}
	fields := strings.Fields(string(pubData))
	if len(fields) < 2 {
		t.Fatalf("key.pub has %d fields, want at least 2", len(fields))
	}
	allowed := filepath.Join(dir, "allowed")
	line := "xc1 namespaces=\"" + Namespace + "\" " + fields[0] + " " + fields[1] + "\n"
	if err := os.WriteFile(allowed, []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}

	// 4. Arms. Each is verified by ssh-keygen -Y verify and by VerifySSHSIG; both must agree.
	arms := []struct {
		name    string
		sig     []byte
		message []byte
		accept  bool // the expected verdict, so that the good arm is non-vacuous
	}{
		{"good", good, msg, true},
		{"tampered-message", good, flipByte(msg, 0), false},
		{"git-namespace", gitSig, msg, false},
		{"another-key", otherSig, msg, false},
		{"flipped-signature", withLastSigByteFlipped(t, good), msg, false},
		{"good-sha256", good256, msg, true},
		{"tampered-message-sha256", good256, flipByte(msg, 0), false},
		{"flipped-signature-sha256", withLastSigByteFlipped(t, good256), msg, false},
	}
	for _, arm := range arms {
		sigPath := filepath.Join(dir, arm.name+".sig")
		if err := os.WriteFile(sigPath, arm.sig, 0o600); err != nil {
			t.Fatal(err)
		}
		_, rc := runSSHKeygen(t, arm.message, "-Y", "verify", "-f", allowed, "-I", "xc1",
			"-n", Namespace, "-s", sigPath)

		ours := func() error {
			s, err := ParseSSHSIG(arm.sig)
			if err != nil {
				return err
			}
			return VerifySSHSIG(s, Namespace, arm.message, pub)
		}()

		if (rc == 0) != (ours == nil) {
			t.Errorf("arm %s: ssh-keygen rc %d but VerifySSHSIG err %v; they must agree", arm.name, rc, ours)
		}
		if arm.accept {
			if rc != 0 || ours != nil {
				t.Errorf("arm %s (good): rc %d, VerifySSHSIG err %v; want both accept", arm.name, rc, ours)
			}
		} else if rc == 0 {
			t.Errorf("arm %s: ssh-keygen accepted a signature that must be refused (instrument sanity)", arm.name)
		}
	}

	// 5. The fingerprint that ssh-keygen prints must be the one the production parser computes.
	fp, rc := runSSHKeygen(t, nil, "-l", "-f", key+".pub")
	if rc != 0 {
		t.Fatalf("ssh-keygen -l: rc %d, want 0", rc)
	}
	if !strings.Contains(string(fp), pub.Fingerprint()) {
		t.Fatalf("ssh-keygen -l %q does not contain %q", fp, pub.Fingerprint())
	}
}
