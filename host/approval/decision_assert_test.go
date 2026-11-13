package approval_test

// Row 201 A4 sprint 3 M6b (design §2.3, §2.7 [r2], §6 DM7, the passkey halves of DM3-DM6): the
// decision assertion law with the approvaltest authenticator (D-WORLD-109) and, for the real tool,
// Chrome's recorded assertion (A2's fixture): P2-P5, P3w and W1-W14 over the decision message's
// hash, and no W15 (the counter is the credentials writer's). The A2 matrix stays unedited.

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

var dmPasskey = approvaltest.NewPasskey(approval.AlgES256, 61)

// dmHash is the decision message hash of (verb, dmP, dmQ, expiry).
func dmHash(t *testing.T, verb byte, expiry int64) [32]byte {
	t.Helper()
	h, err := approval.DecisionMessageHash(verb, dmP, dmQ, expiry)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// dmaInput is a valid decision assertion: mark-mac's passkey (enrolled at 900, stored counter 1)
// signs approve of (P, Q, expiry) with counter count at paOrigin.
func dmaInput(t *testing.T, count uint32, e ...approvaltest.Edit) approval.DecisionAssertionInput {
	t.Helper()
	a := paApprover(t, dmPasskey)
	return approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: dmP, Request: dmQ, Expiry: dmExpiry, Now: dmNow,
		Response: dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), count, e...), Approver: a}
}

// DM7: each of P2-P5, P3w and W1-W14 refuses an otherwise valid decision assertion, with the
// decision challenge; a nil decision on every refusal.
func TestDecideByAssertionSteps(t *testing.T) {
	mint, err := approval.MessageHash(dmQ, dmP, dmExpiry)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		edit func(t *testing.T, in *approval.DecisionAssertionInput)
		want error
	}{
		{"P2-name", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName},
		{"P3w-format", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey},
		{"P3w-no-passkey", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Passkey = nil }, approval.ErrApproverNotPasskey},
		{"P4-revoked", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked},
		{"P5-clock-behind", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.EnrolledAt = dmNow + 1 }, approval.ErrApproverClockBehind},
		{"W1-envelope", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse},
		{"W2-raw-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{RawID: []byte("another")})
		}, approval.ErrCredentialIDMismatch},
		{"W3-user-handle", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{UserHandle: []byte("mallory")})
		}, approval.ErrUserHandle},
		{"W5-type", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Type: "webauthn.create"})
		}, approval.ErrClientDataType},
		{"W6-the-mint-message", func(_ *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, mint, 2)
		}, approval.ErrChallenge},
		{"W6-the-other-verb", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Verb = approval.VerbDeny }, approval.ErrChallenge},
		{"W6-another-packet", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Packet[0] ^= 1 }, approval.ErrChallenge},
		{"W6-another-request", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Request[0] ^= 1 }, approval.ErrChallenge},
		{"W6-another-expiry", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Expiry++ }, approval.ErrChallenge},
		{"W6-no-message", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Verb = 3 }, approval.ErrChallenge},
		{"W6-no-message-signed-zero", func(_ *testing.T, in *approval.DecisionAssertionInput) { // the zero hash an absent message would leave
			in.Verb, in.Response = 3, dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2)
		}, approval.ErrChallenge},
		{"W7-origin", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert("https://x.ts.net", paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2)
		}, approval.ErrOriginMismatch},
		{"W8-cross-origin", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{CrossOrigin: true})
		}, approval.ErrCrossOrigin},
		{"W9-flags", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Set: 0x40})
		}, approval.ErrAuthDataFlags},
		{"W10-rp-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, "y.ts.net", dmHash(t, approval.VerbApprove, dmExpiry), 2)
		}, approval.ErrRPIDHash},
		{"W11-user-presence", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Clear: 0x01})
		}, approval.ErrUserPresence},
		{"W12-user-verification", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Clear: 0x04})
		}, approval.ErrUserVerification},
		{"W13-backup-eligibility", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Set: 0x08})
		}, approval.ErrBackupEligibility},
		{"W14-bad-signature", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{BadSignature: true})
		}, approval.ErrBadSignature},
	} {
		in := dmaInput(t, 2)
		tc.edit(t, &in)
		got, err := approval.DecideByAssertion(in)
		if !errors.Is(err, tc.want) || got != nil {
			t.Errorf("%s: %v (decision %v), want %v", tc.name, err, got, tc.want)
		}
	}
}

// DM7 (W15 deferred) + the control: a valid decision assertion decides as webauthn:mark-mac, binds
// its subject, and hands the writer the credential ID and the assertion's counter, which the law
// did not judge: counter 1 (= the stored 1) and 0 are verified here too; the writer refuses them
// (KD2). The mint law over the same response is refused (DM4's passkey half).
func TestDecideByAssertionDefersTheCounter(t *testing.T) {
	for _, count := range []uint32{2, 1, 0} {
		in := dmaInput(t, count)
		got, err := approval.DecideByAssertion(in)
		if err != nil {
			t.Fatalf("counter %d: %v", count, err)
		}
		id, n, ok := got.PasskeyUse()
		if got.DecidedBy() != "webauthn:mark-mac" || got.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry) != nil || !ok || n != count || !bytes.Equal(id, dmPasskey.CredentialID()) {
			t.Fatalf("counter %d: %q, binds %v, use %x %d %v", count, got.DecidedBy(), got.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry), id, n, ok)
		}
		if _, format, key, ok := got.Use(); format != approval.FormatWebAuthn || !ok || !bytes.Equal(key, in.Approver.Passkey.Key.Canonical()) {
			t.Fatalf("counter %d: use %s %x %v", count, format, key, ok)
		}
	}
	in := dmaInput(t, 2)
	mint := approval.AssertionInput{Admit: approval.AdmitInput{Now: dmNow, CreatedAt: dmNow - 10, ExpiresAt: dmExpiry, StoredDigest: dmQ, BytesDigest: dmQ},
		Challenge: dmP, Response: in.Response, Approver: in.Approver}
	if got, err := approval.ApproveByAssertion(mint); !errors.Is(err, approval.ErrChallenge) || got != nil {
		t.Fatalf("a decision assertion as a mint: %v", err)
	}
	if _, _, ok := (&approval.SignedDecision{}).PasskeyUse(); ok {
		t.Fatal("the zero decision has a passkey use")
	}
	s := approvaltest.NewSigner(41)
	ssh, err := approval.DecideBySignature(dmInput(t, s, approval.VerbApprove))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, ok := ssh.PasskeyUse(); ok {
		t.Fatal("an ssh decision has a passkey use")
	}
}
