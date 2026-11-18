package approval_test

// Row 243 V4 sprint 2b (judge-408 notes 2 and 3): the proof's passkey use is its own copy (X-PK-ALIAS), and
// the runner's recheck bounds the recorded expiry, not the request's lifetime (design §3.1).

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// TestOperatorPasskeyUseIsACopy: neither the caller's approver nor a returned slice, edited after the
// call, changes the credential ID a proof binds (S3's confirm CAS keys on it).
func TestOperatorPasskeyUseIsACopy(t *testing.T) {
	want := opPasskey.CredentialID()
	in := opIn(t, approval.FormatWebAuthn)
	confirmed, err := approval.ConfirmOperatorAction(in)
	if err != nil {
		t.Fatal(err)
	}
	in.Approver.Passkey.SignCount = 2
	rechecked, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: in.Action, Challenge: opC, Expiry: opExp, Now: opNow,
		Artifact: in.Artifact, Format: approval.FormatWebAuthn, Approver: in.Approver})
	if err != nil {
		t.Fatal(err)
	}
	in.Approver.Passkey.CredentialID[0] ^= 0xff
	for name, s := range map[string]*approval.SignedOperatorAction{"confirmed": confirmed, "rechecked": rechecked} {
		cid, _, ok := s.PasskeyUse()
		if !ok || !bytes.Equal(cid, want) {
			t.Fatalf("%s: the approver's slice, edited after the call, moved the proof's credential ID to %x", name, cid)
		}
		cid[1] ^= 0xff
		if again, _, _ := s.PasskeyUse(); !bytes.Equal(again, want) {
			t.Fatalf("%s: a returned slice, edited, moved the proof's credential ID to %x", name, again)
		}
	}
}

// TestOperatorRunnerRecheckBoundsExpiryNotLifetime pins judge-408 note 3. Design §3.1: the 600 s bound is
// the signed request's ("signed request expiry 600 s"), judged by the daemon at confirmation (C2), where
// the row's creation time is known; the runner "recheck[s] expiry and revocation again at each new effect
// boundary" from its own recorded values, and OperatorCheckInput carries no creation time. So a recorded
// expiry an hour away is the runner's to admit while unexpired, and the confirmation refuses it.
func TestOperatorRunnerRecheckBoundsExpiryNotLifetime(t *testing.T) {
	const exp = opCreated + 3600
	a := opAction("start")
	for format, art := range map[string][]byte{approval.FormatSSH: opSSH(t, opSigner, approval.OperatorActionNamespace, a, exp),
		approval.FormatWebAuthn: opAssert(t, a, exp, 2)} {
		in := opIn(t, format)
		in.Artifact, in.Admit.ExpiresAt = art, exp
		if s, err := approval.ConfirmOperatorAction(in); s != nil || !errors.Is(err, approval.ErrOperatorLifetime) {
			t.Errorf("%s: the confirmation admitted a 3600 s request: %v", format, err)
		}
		if in.Approver.Passkey != nil {
			in.Approver.Passkey.SignCount = 2
		}
		for _, now := range []int64{opNow, exp - 1} {
			s, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: a, Challenge: opC, Expiry: exp, Now: now, Artifact: art,
				Format: format, Approver: in.Approver})
			if err != nil || s.Binds(a, opC, exp) != nil {
				t.Errorf("%s at %d: the recheck refused an unexpired recorded expiry: %v", format, now, err)
			}
		}
	}
}
