package main

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M5a-M5c): approve binds the stored challenge, and approves by an enrolled key's signature.
// Every store is under t.TempDir(); every key is approvaltest.NewSigner's or a checked-in fixture
// (D-WORLD-109).

import (
	"bytes"
	"strings"
	"testing"
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
