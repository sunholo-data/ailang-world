package main

// Row 201 A4 sprint 4 M7a (the sprint-3 judge's non-blocking findings 3 and 4): the packet and
// request digests the CLI posts with --signature and --webauthn-assertion are the subject's, and
// both artifact reads refuse a file one byte over their bound before any POST. The fake daemon
// (approve_decide_rig_test.go) answers any packet, so only these checks pin the body's digests.

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// The --signature and --webauthn-assertion bodies carry P and Q as the subject served them.
func TestApproveDecideArtifactPostsTheSubjectsDigests(t *testing.T) {
	f := newFakeDecider(t)
	env := approveEnv(nil, nil, atTime(1710))
	f.enrolDecisionPasskey(t, true)
	sigPath := filepath.Join(t.TempDir(), "decide.sig")
	if err := os.WriteFile(sigPath, []byte("SIGNED-ELSEWHERE"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbApprove, 1760))
	asrPath := filepath.Join(t.TempDir(), "world-assertion.json")
	if err := os.WriteFile(asrPath, cdPasskey.Assert(cdOrigin, cdRP, h, 2), 0o600); err != nil {
		t.Fatal(err)
	}
	p, q := strings.TrimPrefix(f.subject.Packet, "sha256:"), strings.TrimPrefix(f.ref, "sha256:")
	for i, args := range [][]string{{"--signature", sigPath, "--expiry", "1760"}, {"--webauthn-assertion", asrPath, "--expiry", "1760"}} {
		if code, _, stderr := f.keyRun(t, env, args...); code != exitOK || len(f.posts) != i+1 {
			t.Fatalf("%v: rc %d, %d POSTs\nstderr %s", args, code, len(f.posts), stderr)
		}
		if b := f.postBodies[i]; b.Packet != p || b.Request != q || b.Verb != "approve" || b.V != 1 {
			t.Errorf("%v posted packet %q request %q verb %q v %d, want %q %q approve 1", args, b.Packet, b.Request, b.Verb, b.V, p, q)
		}
	}
}

// Each artifact flag reads at most its bound: one byte over is refused with the bound named and
// nothing is posted; exactly the bound is read (and, not being an artifact, refused or posted).
func TestApproveDecideArtifactReadBounds(t *testing.T) {
	f := newFakeDecider(t)
	env := approveEnv(nil, nil, atTime(1710))
	for _, c := range []struct {
		flag string
		max  int
	}{{"--signature", approveMaxSignatureBytes}, {"--webauthn-assertion", approval.MaxPasskeyResponseBytes}} {
		over := filepath.Join(t.TempDir(), "over")
		if err := os.WriteFile(over, bytes.Repeat([]byte("x"), c.max+1), 0o600); err != nil {
			t.Fatal(err)
		}
		code, _, stderr := f.keyRun(t, env, c.flag, over, "--expiry", "1760")
		if code != exitUsage || len(f.posts) != 0 || !strings.Contains(stderr, "over the "+itoa(c.max)+"-byte bound") {
			t.Errorf("%s one byte over %d: rc %d, %d POSTs\nstderr %s", c.flag, c.max, code, len(f.posts), stderr)
		}
	}
	at := filepath.Join(t.TempDir(), "at")
	if err := os.WriteFile(at, bytes.Repeat([]byte("x"), approveMaxSignatureBytes), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := f.keyRun(t, env, "--signature", at, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 || len(*f.postBodies[0].SSHSIG) != approveMaxSignatureBytes {
		t.Errorf("--signature at the bound: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
	}
}
