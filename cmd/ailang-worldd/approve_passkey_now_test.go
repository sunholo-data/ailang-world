package main

// Row 201 A3 sprint 1 M2b (planner, the mutant MUT-WP-JUDGE-NOW): through the shared judge, the
// passkey law judges the approver at approve's own clock reading, not at the request's creation —
// the passkey twin of TestApproveBySignatureJudgesTheApproverAtNow (A1b sprint 4 SIG-NOW).

import (
	"bytes"
	"testing"
)

// PK-NOW: a passkey enrolled after the request was made (+5) and before the approve (+10) approves it.
func TestApproveByPasskeyJudgesTheApproverAtNow(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0+5)
	file := writeAssertion(t, passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2))
	var opened bool
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+10))); code != exitOK || opened {
		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal for a passkey enrolled at +5 and an approve at +10\nstderr %s", code, opened, stderr.String())
	}
	if row := a.row(t); !row.Consumed || row.ApprovedBy != "webauthn:mark-mac" {
		t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", row.Consumed, row.ApprovedBy)
	}
}
