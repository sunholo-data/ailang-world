package main

// Row 243 V3 sprint 5 M13 (design §2.8, §6 CB4): `approve vr-<16 hex>` at the CLI points at the
// finish confirmation's card, as the bar's ApproveTarget does; it reads no store and mints nothing.

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestApproveVerbRequestPointsAtItsCard(t *testing.T) {
	db := filepath.Join(t.TempDir(), "world.db")
	var stdout, stderr bytes.Buffer
	code := runApprove([]string{"vr-0123456789abcdef", "--db", db}, &stdout, &stderr, noTerminalEnv(1))
	want := "ailang-worldd approve: request vr-0123456789abcdef is a finish confirmation: confirm it on its card, /surface/confirm/vr-0123456789abcdef; the waiting episode finish opens the pull request\n"
	if code != exitUsage || stderr.String() != want || stdout.Len() != 0 {
		t.Fatalf("rc %d\nstderr %q\nwant   %q", code, stderr.String(), want)
	}
	stderr.Reset()
	if code := runApprove([]string{"vr-0123", "--db", db}, &stdout, &stderr, noTerminalEnv(1)); code != exitUsage || !bytes.Contains(stderr.Bytes(), []byte(`"vr-0123" is not a mint request id`)) {
		t.Fatalf("a short id: rc %d %q", code, stderr.String())
	}
}
