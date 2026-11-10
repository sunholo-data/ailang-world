package sshsign

// Row 201 A1b sprint 4 S4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md
// §3 S4): the signer's buffer bound, and the child environment it hands ssh-keygen. No real
// ssh-keygen and no key: the program is a temp shell script.

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

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

// The signer's child sees what ssh-keygen needs (SSH_AUTH_SOCK, HOME, PATH) and never the
// registry credential (AC10, childenv.Scrubbed).
func TestSignHandsTheSignerAScrubbedEnvironment(t *testing.T) {
	dir := t.TempDir()
	dump := filepath.Join(dir, "child.env")
	bin := filepath.Join(dir, "signer")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n/usr/bin/env > '"+dump+"'\ncat > /dev/null\necho signature\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv(childenv.CredentialVariable, "ac10-sentinel")
	t.Setenv("SSH_AUTH_SOCK", "/tmp/agent.sock")
	t.Setenv("HOME", dir)
	out, err := Sign(t.Context(), bin, filepath.Join(dir, "k"), []byte("msg"), io.Discard)
	if err != nil || string(out) != "signature\n" {
		t.Fatalf("Sign = (%q, %v), want the signer's output", out, err)
	}
	b, err := os.ReadFile(dump)
	if err != nil {
		t.Fatal(err)
	}
	env := strings.Split(strings.TrimSpace(string(b)), "\n")
	if childenv.Has(env, childenv.CredentialVariable) {
		t.Fatalf("the signer observed %s", childenv.CredentialVariable)
	}
	for _, want := range []string{"SSH_AUTH_SOCK=/tmp/agent.sock", "HOME=" + dir} {
		found := false
		for _, kv := range env {
			found = found || kv == want
		}
		if !found {
			t.Errorf("the signer's environment lacks %s", want)
		}
	}
}

// SG1 (row 201 A2 sprint 1 M0): a signer that prints something and then fails is refused (G1),
// and its output is never returned as a signature.
func TestSignRefusesAFailedSignerEvenWithOutput(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "signer")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\ncat > /dev/null\necho partial-signature\nexit 3\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	out, err := Sign(t.Context(), bin, "k", []byte("msg"), io.Discard)
	var exit *exec.ExitError
	if out != nil || !errors.As(err, &exit) || exit.ExitCode() != 3 {
		t.Fatalf("Sign = (%q, %v), want nil bytes and an error wrapping exit status 3", out, err)
	}
}
