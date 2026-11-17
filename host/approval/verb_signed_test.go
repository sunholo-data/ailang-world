package approval_test

// Row 243 V3 sprint 1 M1b (design §2.3, §6 VL4-VL7): the verb law with approvaltest keys and
// passkeys (D-WORLD-109): domain separation both directions, the signature law's P1-P6, the
// assertion law's P1-P5, P3w, W1-W15, and what a SignedVerb binds and hands its consume.

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

var vD, vC = [32]byte{0x76, 1}, [32]byte{0x76, 2}

const vCreated, vExpiry, vNow = int64(1000), int64(1600), int64(1010)

var vPasskey = approvaltest.NewPasskey(approval.AlgES256, 71)

// vAdmit is an admitted request: created 1000, expires 1600, now 1010, bytes = digest vD.
func vAdmit() approval.AdmitInput {
	return approval.AdmitInput{Now: vNow, CreatedAt: vCreated, ExpiresAt: vExpiry, StoredDigest: vD, BytesDigest: vD}
}

// vMsg and vHash are what a confirmation of (digest, challenge, expiry) signs.
func vMsg(t *testing.T, d, c [32]byte, e int64) []byte {
	t.Helper()
	m, err := approval.VerbMessage(d, c, e)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func vHash(t *testing.T, d, c [32]byte, e int64) [32]byte {
	t.Helper()
	h, err := approval.VerbMessageHash(d, c, e)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// vsInput is a valid signature confirmation: s's SSHSIG over the request's verb message, s
// enrolled as mark-ssh at 900.
func vsInput(t *testing.T, s *approvaltest.Signer) approval.VerbSignatureInput {
	t.Helper()
	return approval.VerbSignatureInput{Admit: vAdmit(), Challenge: vC, Signature: s.Sign(approval.VerbNamespace, vMsg(t, vD, vC, vExpiry)),
		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
}

// vaInput is a valid assertion confirmation: vPasskey (enrolled as mark-mac at 900, stored counter
// 1) asserts the request's verb message hash with counter count.
func vaInput(t *testing.T, count uint32, e ...approvaltest.Edit) approval.VerbAssertionInput {
	t.Helper()
	return approval.VerbAssertionInput{Admit: vAdmit(), Challenge: vC,
		Response: vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), count, e...), Approver: paApprover(t, vPasskey)}
}

// VL4: domain separation, both directions, for both key kinds. A mint or decision SSHSIG never
// confirms a verb (wrong namespace); a verb SSHSIG never approves a mint nor decides (wrong
// namespace); a mint or decision assertion never confirms (W6), nor a verb assertion mints or
// decides (W6).
func TestVerbDomainSeparation(t *testing.T) {
	s := approvaltest.NewSigner(72)
	mintMsg, _ := approval.SignedMessage(vD, vC, vExpiry)
	decMsg, _ := approval.DecisionMessage(approval.VerbApprove, vD, vC, vExpiry)
	for _, tc := range []struct {
		name string
		sig  []byte
		want error
	}{
		{"a mint SSHSIG", s.Sign(approval.Namespace, mintMsg), approval.ErrNamespace},
		{"a decision SSHSIG", s.Sign(approval.DecisionNamespace, decMsg), approval.ErrNamespace},
		{"the verb message, mint namespace", s.Sign(approval.Namespace, vMsg(t, vD, vC, vExpiry)), approval.ErrNamespace},
		{"the verb message, decision namespace", s.Sign(approval.DecisionNamespace, vMsg(t, vD, vC, vExpiry)), approval.ErrNamespace},
		{"the mint message, verb namespace", s.Sign(approval.VerbNamespace, mintMsg), approval.ErrBadSignature},
		{"the decision message, verb namespace", s.Sign(approval.VerbNamespace, decMsg), approval.ErrBadSignature},
		{"the verb message without its zero byte", s.Sign(approval.VerbNamespace, append([]byte(approval.VerbPrefix), vMsg(t, vD, vC, vExpiry)[14:]...)), approval.ErrBadSignature},
	} {
		in := vsInput(t, s)
		in.Signature = tc.sig
		if v, err := approval.ConfirmBySignature(in); v != nil || !errors.Is(err, tc.want) {
			t.Errorf("%s: (%v, %v), want %v", tc.name, v, err, tc.want)
		}
	}
	verb := vsInput(t, s).Signature
	mint := approval.SignatureInput{Admit: vAdmit(), Challenge: vC, Signature: verb,
		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
	if a, err := approval.ApproveBySignature(mint); a != nil || !errors.Is(err, approval.ErrNamespace) {
		t.Errorf("a verb SSHSIG minted: (%v, %v)", a, err)
	}
	dec := approval.DecisionSignatureInput{Verb: approval.VerbApprove, Packet: vD, Request: vC, Expiry: vExpiry, Now: vNow, Signature: verb, Approver: mint.Approver}
	if d, err := approval.DecideBySignature(dec); d != nil || !errors.Is(err, approval.ErrNamespace) {
		t.Errorf("a verb SSHSIG decided: (%v, %v)", d, err)
	}
	if v, err := approval.ConfirmBySignature(vsInput(t, s)); err != nil || v.ApprovedBy() != "ssh:mark-ssh" {
		t.Fatalf("control: (%v, %v)", v, err)
	}

	mintHash, _ := approval.MessageHash(vD, vC, vExpiry)
	decHash, _ := approval.DecisionMessageHash(approval.VerbApprove, vD, vC, vExpiry)
	for name, h := range map[string][32]byte{"a mint assertion": mintHash, "a decision assertion": decHash} {
		in := vaInput(t, 2)
		in.Response = vPasskey.Assert(paOrigin, paRPID, h, 2)
		if v, err := approval.ConfirmByAssertion(in); v != nil || !errors.Is(err, approval.ErrChallenge) {
			t.Errorf("%s confirmed: (%v, %v), want ErrChallenge (W6)", name, v, err)
		}
	}
	va := vaInput(t, 2)
	if a, err := approval.ApproveByAssertion(approval.AssertionInput{Admit: va.Admit, Challenge: vC, Response: va.Response, Approver: va.Approver}); a != nil || !errors.Is(err, approval.ErrChallenge) {
		t.Errorf("a verb assertion minted: (%v, %v)", a, err)
	}
	if d, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: vD, Request: vC, Expiry: vExpiry, Now: vNow,
		Response: va.Response, Approver: va.Approver}); d != nil || !errors.Is(err, approval.ErrChallenge) {
		t.Errorf("a verb assertion decided: (%v, %v)", d, err)
	}
	if v, err := approval.ConfirmByAssertion(va); err != nil || v.ApprovedBy() != "webauthn:mark-mac" {
		t.Fatalf("control: (%v, %v)", v, err)
	}
}

// VL5: the signature law's order: P1 each verdict (before the approver), P2-P5 (before the
// cryptography), then P6: a bad signature, a signature over another digest, challenge or expiry.
func TestConfirmBySignatureLaw(t *testing.T) {
	s, other := approvaltest.NewSigner(73), approvaltest.NewSigner(74)
	for _, tc := range []struct {
		name    string
		edit    func(in *approval.VerbSignatureInput)
		want    error
		verdict approval.Verdict
	}{
		{"P1 consumed", func(in *approval.VerbSignatureInput) { in.Admit.Consumed = true }, nil, approval.RefuseConsumed},
		{"P1 consumed before P4", func(in *approval.VerbSignatureInput) { in.Admit.Consumed, in.Approver.RevokedAt = true, 950 }, nil, approval.RefuseConsumed},
		{"P1 clock behind", func(in *approval.VerbSignatureInput) { in.Admit.Now = vCreated - 1 }, nil, approval.RefuseClockBehind},
		{"P1 expired", func(in *approval.VerbSignatureInput) { in.Admit.Now = vExpiry }, nil, approval.RefuseExpired},
		{"P1 digest", func(in *approval.VerbSignatureInput) { in.Admit.BytesDigest[0] ^= 1 }, nil, approval.RefuseDigestMismatch},
		{"P2 name", func(in *approval.VerbSignatureInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName, 0},
		{"P3 format", func(in *approval.VerbSignatureInput) { in.Approver.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat, 0},
		{"P4 revoked", func(in *approval.VerbSignatureInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked, 0},
		{"P4 revoked before P6", func(in *approval.VerbSignatureInput) { in.Approver.RevokedAt, in.Signature = 950, []byte("x") }, approval.ErrApproverRevoked, 0},
		{"P5 clock behind", func(in *approval.VerbSignatureInput) { in.Approver.EnrolledAt = vNow + 1 }, approval.ErrApproverClockBehind, 0},
		{"P6 not armored", func(in *approval.VerbSignatureInput) { in.Signature = []byte("y\n") }, approval.ErrArmor, 0},
		{"P6 another key", func(in *approval.VerbSignatureInput) { in.Approver.Key = other.PublicKey() }, approval.ErrKeyMismatch, 0},
		{"P6 another digest", func(in *approval.VerbSignatureInput) {
			d := vD
			d[1] ^= 1
			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d
		}, approval.ErrBadSignature, 0},
		{"P6 another challenge", func(in *approval.VerbSignatureInput) { in.Challenge[1] ^= 1 }, approval.ErrBadSignature, 0},
		{"P6 another expiry", func(in *approval.VerbSignatureInput) { in.Admit.ExpiresAt++ }, approval.ErrBadSignature, 0},
	} {
		in := vsInput(t, s)
		tc.edit(&in)
		v, err := approval.ConfirmBySignature(in)
		var na *approval.NotAdmittedError
		switch {
		case v != nil:
			t.Errorf("%s: confirmed", tc.name)
		case tc.verdict != 0 && (!errors.As(err, &na) || na.Verdict != tc.verdict):
			t.Errorf("%s: %v, want verdict %v", tc.name, err, tc.verdict)
		case tc.want != nil && !errors.Is(err, tc.want):
			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
		}
	}
	v, err := approval.ConfirmBySignature(vsInput(t, s))
	if err != nil || v.Binds(vD, vC, vExpiry) != nil || v.ApprovedBy() != "ssh:mark-ssh" {
		t.Fatalf("control: (%v, %v)", v, err)
	}
	if name, format, key, ok := v.Use(); !ok || name != "mark-ssh" || format != approval.FormatSSH || !bytes.Equal(key, s.PublicKey().Blob) {
		t.Fatalf("Use = %q %q %x %v", name, format, key, ok)
	}
	if _, _, ok := v.PasskeyUse(); ok {
		t.Fatal("an ssh confirmation has a passkey use")
	}
}

// VL6: the assertion law's order: P1, P2-P5 and P3w, W1-W14 over the verb message hash, then W15 in
// full (a counter that does not advance is refused here); the control hands the consume the
// credential ID and the assertion's counter.
func TestConfirmByAssertionLaw(t *testing.T) {
	for _, tc := range []struct {
		name    string
		edit    func(t *testing.T, in *approval.VerbAssertionInput)
		want    error
		verdict approval.Verdict
	}{
		{"P1 consumed", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Consumed = true }, nil, approval.RefuseConsumed},
		{"P1 clock behind", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Now = vCreated - 1 }, nil, approval.RefuseClockBehind},
		{"P1 expired", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Now = vExpiry + 1 }, nil, approval.RefuseExpired},
		{"P1 digest", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.BytesDigest[0] ^= 1 }, nil, approval.RefuseDigestMismatch},
		{"P2 name", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName, 0},
		{"P3w format", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey, 0},
		{"P3w no passkey", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Passkey = nil }, approval.ErrApproverNotPasskey, 0},
		{"P4 revoked", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked, 0},
		{"P5 clock behind", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.EnrolledAt = vNow + 1 }, approval.ErrApproverClockBehind, 0},
		{"W1 envelope", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse, 0},
		{"W2 raw id", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{RawID: []byte("another")})
		}, approval.ErrCredentialIDMismatch, 0},
		{"W3 user handle", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{UserHandle: []byte("mallory")})
		}, approval.ErrUserHandle, 0},
		{"W5 type", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Type: "webauthn.create"})
		}, approval.ErrClientDataType, 0},
		{"W6 another digest", func(_ *testing.T, in *approval.VerbAssertionInput) {
			d := vD
			d[1] ^= 1
			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d
		}, approval.ErrChallenge, 0},
		{"W6 another challenge", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Challenge[1] ^= 1 }, approval.ErrChallenge, 0},
		{"W6 another expiry", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.ExpiresAt++ }, approval.ErrChallenge, 0},
		{"W7 origin", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert("https://y.ts.net", paRPID, vHash(t, vD, vC, vExpiry), 2)
		}, approval.ErrOriginMismatch, 0},
		{"W8 cross origin", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{CrossOrigin: true})
		}, approval.ErrCrossOrigin, 0},
		{"W9 flags", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Set: 0x40})
		}, approval.ErrAuthDataFlags, 0},
		{"W10 rp id", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, "y.ts.net", vHash(t, vD, vC, vExpiry), 2)
		}, approval.ErrRPIDHash, 0},
		{"W11 user presence", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Clear: 0x01})
		}, approval.ErrUserPresence, 0},
		{"W12 user verification", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Clear: 0x04})
		}, approval.ErrUserVerification, 0},
		{"W13 backup eligibility", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Set: 0x08})
		}, approval.ErrBackupEligibility, 0},
		{"W14 bad signature", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{BadSignature: true})
		}, approval.ErrBadSignature, 0},
		{"W15 the stored counter", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 1)
		}, approval.ErrSignCountRegression, 0},
		{"W15 a zero counter over a stored 1", func(t *testing.T, in *approval.VerbAssertionInput) {
			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 0)
		}, approval.ErrSignCountRegression, 0},
	} {
		in := vaInput(t, 2)
		tc.edit(t, &in)
		v, err := approval.ConfirmByAssertion(in)
		var na *approval.NotAdmittedError
		switch {
		case v != nil:
			t.Errorf("%s: confirmed", tc.name)
		case tc.verdict != 0 && (!errors.As(err, &na) || na.Verdict != tc.verdict):
			t.Errorf("%s: %v, want verdict %v", tc.name, err, tc.verdict)
		case tc.want != nil && !errors.Is(err, tc.want):
			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
		}
	}
	in := vaInput(t, 2)
	v, err := approval.ConfirmByAssertion(in)
	if err != nil || v.Binds(vD, vC, vExpiry) != nil || v.ApprovedBy() != "webauthn:mark-mac" {
		t.Fatalf("control: (%v, %v)", v, err)
	}
	id, n, ok := v.PasskeyUse()
	if !ok || n != 2 || !bytes.Equal(id, vPasskey.CredentialID()) {
		t.Fatalf("PasskeyUse = %x %d %v", id, n, ok)
	}
	if _, format, key, ok := v.Use(); !ok || format != approval.FormatWebAuthn || !bytes.Equal(key, in.Approver.Passkey.Key.Canonical()) {
		t.Fatalf("Use = %s %x %v", format, key, ok)
	}
	zero := vaInput(t, 0) // a counter-less authenticator (stored 0, asserted 0) passes W15
	zero.Approver.Passkey.SignCount = 0
	if _, err := approval.ConfirmByAssertion(zero); err != nil {
		t.Fatalf("a counter-less authenticator: %v", err)
	}
}

// VL7: Binds refuses each field that differs, the nil and the zero value; ApprovedBy, Use and
// PasskeyUse are empty unless verified.
func TestSignedVerbBinds(t *testing.T) {
	v, err := approval.ConfirmBySignature(vsInput(t, approvaltest.NewSigner(75)))
	if err != nil {
		t.Fatal(err)
	}
	d, c := vD, vC
	d[2], c[2] = 9, 9
	for field, err := range map[string]error{
		"digest":    v.Binds(d, vC, vExpiry),
		"challenge": v.Binds(vD, c, vExpiry),
		"expiry":    v.Binds(vD, vC, vExpiry+1),
	} {
		var be *approval.BindingError
		if !errors.As(err, &be) || be.Field != field {
			t.Errorf("%s: %v", field, err)
		}
	}
	if err := v.Binds(vD, vC, vExpiry); err != nil {
		t.Fatalf("control: %v", err)
	}
	var none *approval.SignedVerb
	zero := &approval.SignedVerb{}
	for name, s := range map[string]*approval.SignedVerb{"nil": none, "zero": zero} {
		if err := s.Binds([32]byte{}, [32]byte{}, 0); !errors.Is(err, approval.ErrNotVerified) {
			t.Errorf("%s: Binds = %v", name, err)
		}
		if s.ApprovedBy() != "" {
			t.Errorf("%s: ApprovedBy = %q", name, s.ApprovedBy())
		}
		if _, _, _, ok := s.Use(); ok {
			t.Errorf("%s: Use ok", name)
		}
		if _, _, ok := s.PasskeyUse(); ok {
			t.Errorf("%s: PasskeyUse ok", name)
		}
	}
	_, _, key, _ := v.Use()
	key[0] ^= 1
	if _, _, again, _ := v.Use(); again[0] == key[0] {
		t.Fatal("Use hands out its own key")
	}
}
