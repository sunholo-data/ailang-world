package approval_test

// Row 243 V3 sprint 1 M1c (design §2.3 [r1], §6 VL8): CheckConfirmation, the CLI's re-check of a
// recorded confirmation. It verifies an artifact over the caller's own digest, challenge and expiry
// with no P1 (the row is consumed and may be past its expiry) and no W15 (the daemon moved the
// counter to this very assertion's); it judges the approver at now, so a key revoked since the
// confirmation is refused (fail-safe; audit note 5, carried to sprint 5).

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// vcSSH is a recorded SSHSIG confirmation of (vD, vC, vExpiry) re-checked at now: after the
// expiry, as the CLI's poll may run (the 10 s grace).
func vcSSH(t *testing.T, s *approvaltest.Signer) approval.VerbCheckInput {
	t.Helper()
	return approval.VerbCheckInput{Digest: vD, Challenge: vC, Expiry: vExpiry, Now: vExpiry + 5, Format: approval.FormatSSH,
		Artifact: s.Sign(approval.VerbNamespace, vMsg(t, vD, vC, vExpiry)),
		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
}

// vcPasskey is a recorded assertion (counter 2) re-checked after the daemon moved the stored
// counter to 2.
func vcPasskey(t *testing.T) approval.VerbCheckInput {
	t.Helper()
	a := paApprover(t, vPasskey)
	a.Passkey.SignCount = 2
	return approval.VerbCheckInput{Digest: vD, Challenge: vC, Expiry: vExpiry, Now: vExpiry + 5, Format: approval.FormatWebAuthn,
		Artifact: vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2), Approver: a}
}

func TestCheckConfirmation(t *testing.T) {
	s := approvaltest.NewSigner(76)
	mintMsg, _ := approval.SignedMessage(vD, vC, vExpiry)
	mintHash, _ := approval.MessageHash(vD, vC, vExpiry)
	other := vD
	other[3] ^= 1
	for _, tc := range []struct {
		name string
		in   func(t *testing.T) approval.VerbCheckInput
		edit func(in *approval.VerbCheckInput)
		want error
	}{
		{"ssh revoked approver", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.RevokedAt = vExpiry + 1 }, approval.ErrApproverRevoked},
		{"ssh approver name", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName},
		{"ssh enrolled after now", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.EnrolledAt = vExpiry + 6 }, approval.ErrApproverClockBehind},
		{"ssh a webauthn approver", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat},
		{"ssh another digest", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Digest = other }, approval.ErrBadSignature},
		{"ssh another challenge", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Challenge = other }, approval.ErrBadSignature},
		{"ssh another expiry", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Expiry-- }, approval.ErrBadSignature},
		{"ssh a mint-namespace SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Artifact = s.Sign(approval.Namespace, mintMsg) }, approval.ErrNamespace},
		{"ssh a decision-namespace SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) {
			in.Artifact = s.Sign(approval.DecisionNamespace, vMsg(t, vD, vC, vExpiry))
		}, approval.ErrNamespace},
		{"ssh a negative expiry", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Expiry = -1 }, approval.ErrNegativeExpiry},
		{"ssh an assertion as an SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Artifact = vcPasskey(t).Artifact }, approval.ErrArmor},
		{"ssh another key", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Key = approvaltest.NewSigner(77).PublicKey() }, approval.ErrKeyMismatch},
		{"webauthn revoked approver", vcPasskey, func(in *approval.VerbCheckInput) { in.Approver.RevokedAt = vExpiry + 1 }, approval.ErrApproverRevoked},
		{"webauthn an ssh approver", vcPasskey, func(in *approval.VerbCheckInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey},
		{"webauthn another digest", vcPasskey, func(in *approval.VerbCheckInput) { in.Digest = other }, approval.ErrChallenge},
		{"webauthn another challenge", vcPasskey, func(in *approval.VerbCheckInput) { in.Challenge = other }, approval.ErrChallenge},
		{"webauthn another expiry", vcPasskey, func(in *approval.VerbCheckInput) { in.Expiry++ }, approval.ErrChallenge},
		{"webauthn a negative expiry", vcPasskey, func(in *approval.VerbCheckInput) { in.Expiry = -1 }, approval.ErrChallenge},
		{"webauthn a mint assertion", vcPasskey, func(in *approval.VerbCheckInput) { in.Artifact = vPasskey.Assert(paOrigin, paRPID, mintHash, 2) }, approval.ErrChallenge},
		{"webauthn bad signature", vcPasskey, func(in *approval.VerbCheckInput) {
			in.Artifact = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{BadSignature: true})
		}, approval.ErrBadSignature},
		{"no format", vcPasskey, func(in *approval.VerbCheckInput) { in.Format = "" }, approval.ErrVerbFormat},
		{"an unknown format", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Format = "typed" }, approval.ErrVerbFormat},
	} {
		in := tc.in(t)
		tc.edit(&in)
		if v, err := approval.CheckConfirmation(in); v != nil || !errors.Is(err, tc.want) {
			t.Errorf("%s: (%v, %v), want %v", tc.name, v, err, tc.want)
		}
	}
	// The verified re-checks: no P1 (past the expiry, a consumed row has no admission input at all)
	// and no W15 (the stored counter is already this assertion's, and a counter-less 0 over 0).
	for name, in := range map[string]approval.VerbCheckInput{"ssh": vcSSH(t, s), "webauthn, counter moved to 2": vcPasskey(t)} {
		v, err := approval.CheckConfirmation(in)
		if err != nil || v.Binds(vD, vC, vExpiry) != nil {
			t.Fatalf("%s: (%v, %v)", name, v, err)
		}
		if want := in.Format + ":" + in.Approver.Name; v.ApprovedBy() != want {
			t.Fatalf("%s: ApprovedBy = %q, want %q", name, v.ApprovedBy(), want)
		}
	}
	ahead := vcPasskey(t)
	ahead.Approver.Passkey.SignCount = 7 // a later assertion moved the counter on; the recorded one still verifies
	v, err := approval.CheckConfirmation(ahead)
	if err != nil {
		t.Fatalf("a counter already past the artifact's: %v", err)
	}
	id, n, ok := v.PasskeyUse()
	if !ok || n != 2 || !bytes.Equal(id, vPasskey.CredentialID()) {
		t.Fatalf("PasskeyUse = %x %d %v", id, n, ok)
	}
	if _, format, key, ok := v.Use(); !ok || format != approval.FormatWebAuthn || !bytes.Equal(key, ahead.Approver.Passkey.Key.Canonical()) {
		t.Fatalf("Use = %s %x %v", format, key, ok)
	}
	ssh, _ := approval.CheckConfirmation(vcSSH(t, s))
	if _, _, key, ok := ssh.Use(); !ok || !bytes.Equal(key, s.PublicKey().Blob) {
		t.Fatalf("ssh Use = %x %v", key, ok)
	}
}
