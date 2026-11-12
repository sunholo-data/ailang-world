package main

// Row 201 A3 sprint 2 M4a (iteration 384's carried obligation): the minter judges a recorded
// artifact at the poll's own clock reading, never at the request's creation. A judge at
// p.row.CreatedAt left every sprint-1 await test green (controller-reproduced, iteration 384).

import (
	"bytes"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// AW-NOW: an approver enrolled after the request was made (+5) and before the poll (+10) approves
// it through the surface, in both formats: the await twin of PK-NOW and A1b's SIG-NOW.
func TestApproveAwaitJudgesTheApproverAtNow(t *testing.T) {
	for _, format := range []string{"webauthn", "ssh"} {
		t.Run(format, func(t *testing.T) {
			a := newApproveRig(t, "ep1", approveT0)
			var artifact []byte
			if format == "ssh" {
				seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0+5)
				artifact = sigSigner.Sign(approval.Namespace, requestMessage(t, a))
			} else {
				seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0+5)
				artifact = passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
			}
			var opened bool
			env, _ := awaitEnv(approveT0+10, &opened, func(n, now int64) {
				if n == 1 {
					recordSurface(t, a, format, format+":mark-mac", artifact, now)
				}
			})
			var stdout, stderr bytes.Buffer
			if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env); code != exitOK || opened {
				t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal for an approver enrolled at +5 and polls from +10\nstderr %s", code, opened, stderr.String())
			}
			requireMintedBy(t, a, format+":mark-mac", stdout.String())
		})
	}
}
