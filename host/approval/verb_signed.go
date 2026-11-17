package approval

// Row 243 V3 sprint 1 (design §2.3): the pure law of a verb confirmation by an enrolled key, over the
// 86-byte verb message (verb.go). ConfirmBySignature and ConfirmByAssertion are the mint laws'
// shape (signed.go, passkey_assert.go) with the verb message: P1 admission, P2-P5 checkApprover,
// then the signature, and for a passkey W1-W14 (assertionSteps) and W15 in full. CheckConfirmation
// is the CLI's re-check of a recorded confirmation: the same law without P1 and W15. No I/O; now is
// an input.

import (
	"bytes"
	"errors"
	"fmt"
)

// SignedVerb is a verified verb confirmation. Its fields are unexported, so only ConfirmBySignature,
// ConfirmByAssertion and CheckConfirmation can return one with ok set; the zero value binds nothing.
type SignedVerb struct {
	ok           bool
	digest       [32]byte
	challenge    [32]byte
	expiry       int64
	name, format string
	key          []byte
	credentialID []byte // webauthn: the passkey that signed
	signCount    uint32 // webauthn: the assertion's counter
}

// VerbSignatureInput is what ConfirmBySignature judges: the stored request's admission input, its
// challenge, the armored SSHSIG and the enrolled approver its key names.
type VerbSignatureInput struct {
	Admit     AdmitInput // row times, consumed, the stored digest and the SHA-256 of the stored bytes
	Challenge [32]byte   // the request row's challenge
	Signature []byte     // the armored SSHSIG
	Approver  Approver
}

// VerbAssertionInput is what ConfirmByAssertion judges: the stored request's admission input, its
// challenge, the assertion file and the enrolled approver its credential ID names, Passkey set.
type VerbAssertionInput struct {
	Admit     AdmitInput
	Challenge [32]byte
	Response  []byte
	Approver  Approver
}

// VerbCheckInput is what CheckConfirmation re-checks: the CLI's own digest, challenge and expiry,
// the clock, the recorded artifact, its format and the enrolled approver it names.
type VerbCheckInput struct {
	Digest, Challenge [32]byte
	Expiry, Now       int64
	Artifact          []byte
	Format            string // FormatSSH or FormatWebAuthn
	Approver          Approver
}

// verifyVerbSSHSIG parses an armored SSHSIG and verifies it over VerbMessage(digest, challenge,
// expiry) under VerbNamespace with key.
func verifyVerbSSHSIG(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error {
	s, err := ParseSSHSIG(armored)
	if err != nil {
		return err
	}
	msg, err := VerbMessage(digest, challenge, expiry)
	if err != nil {
		return err
	}
	return VerifySSHSIG(s, VerbNamespace, msg, key)
}

// ConfirmBySignature is the verb signature law (P1-P6, in this order): admission first, so a replay
// says "already used"; the approver's state before the cryptography; then the SSHSIG over the verb
// message of the stored digest, the row's challenge and the row's expiry.
func ConfirmBySignature(in VerbSignatureInput) (*SignedVerb, error) {
	if v := Admit(in.Admit); v != Admitted { // P1
		return nil, &NotAdmittedError{Verdict: v}
	}
	if err := checkApprover(in.Admit, in.Approver, FormatSSH); err != nil { // P2-P5
		return nil, err
	}
	if err := verifyVerbSSHSIG(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
		return nil, err
	}
	return &SignedVerb{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge, expiry: in.Admit.ExpiresAt,
		name: in.Approver.Name, format: FormatSSH, key: bytes.Clone(in.Approver.Key.Blob)}, nil
}

// ConfirmByAssertion is the verb assertion law (P1-P5, P3w, W1-W15, in this order), as
// ApproveByAssertion with the challenge VerbMessageHash. W15 runs in full: the counter moves in the
// consume, which runs after the world write, so a retry of the same artifact still passes (design
// §2.3).
func ConfirmByAssertion(in VerbAssertionInput) (*SignedVerb, error) {
	if v := Admit(in.Admit); v != Admitted { // P1
		return nil, &NotAdmittedError{Verdict: v}
	}
	if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
		return nil, err
	}
	h, hErr := VerbMessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
	ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14
	if err != nil {
		return nil, err
	}
	p := in.Approver.Passkey
	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
	}
	return &SignedVerb{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge, expiry: in.Admit.ExpiresAt,
		name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
}

// ErrVerbFormat is a CheckConfirmation format other than FormatSSH and FormatWebAuthn.
var ErrVerbFormat = errors.New("approval: a confirmation's format is ssh or webauthn")

// CheckConfirmation is the CLI's re-check of a recorded confirmation (design §2.3 [r1], §2.7 step
// 6(a)): the verb law without P1 (the row is consumed by then) and without W15 (the counter was
// judged and moved by the daemon), over the caller's own digest, challenge and expiry, as
// DecideByAssertion runs W1-W14 without W15. The approver's state (P2-P5) is judged at now, so a
// key revoked since the confirmation is refused (fail-safe).
func CheckConfirmation(in VerbCheckInput) (*SignedVerb, error) {
	switch in.Format {
	case FormatSSH:
		if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatSSH); err != nil { // P2-P5
			return nil, err
		}
		if err := verifyVerbSSHSIG(in.Artifact, in.Digest, in.Challenge, in.Expiry, in.Approver.Key); err != nil { // P6
			return nil, err
		}
		return &SignedVerb{ok: true, digest: in.Digest, challenge: in.Challenge, expiry: in.Expiry,
			name: in.Approver.Name, format: FormatSSH, key: bytes.Clone(in.Approver.Key.Blob)}, nil
	case FormatWebAuthn:
		if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
			return nil, err
		}
		h, hErr := VerbMessageHash(in.Digest, in.Challenge, in.Expiry)
		ad, err := assertionSteps(in.Artifact, in.Approver, h, hErr) // W1-W14; no W15
		if err != nil {
			return nil, err
		}
		p := in.Approver.Passkey
		return &SignedVerb{ok: true, digest: in.Digest, challenge: in.Challenge, expiry: in.Expiry,
			name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrVerbFormat, in.Format)
}

// Binds reports whether s is a verified confirmation of exactly this digest, challenge and expiry.
func (s *SignedVerb) Binds(digest, challenge [32]byte, expiry int64) error {
	switch {
	case s == nil || !s.ok:
		return ErrNotVerified
	case s.digest != digest:
		return &BindingError{Field: "digest"}
	case s.challenge != challenge:
		return &BindingError{Field: "challenge"}
	case s.expiry != expiry:
		return &BindingError{Field: "expiry"}
	}
	return nil
}

// ApprovedBy is the confirmation's approved_by and the finish approval's decidedBy, format ":" name
// ("webauthn:mark-mac"); "" unless verified.
func (s *SignedVerb) ApprovedBy() string {
	if s == nil || !s.ok {
		return ""
	}
	return s.format + ":" + s.name
}

// Use is what the consume binds in SQL: the approver's name, format and key (a copy).
func (s *SignedVerb) Use() (name, format string, key []byte, ok bool) {
	if s == nil || !s.ok {
		return "", "", nil, false
	}
	return s.name, s.format, bytes.Clone(s.key), true
}

// PasskeyUse is what the consume binds for a passkey: the credential ID (a copy) and the
// assertion's counter; ok only for a verified webauthn confirmation.
func (s *SignedVerb) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
	if s == nil || !s.ok || s.format != FormatWebAuthn {
		return nil, 0, false
	}
	return bytes.Clone(s.credentialID), s.signCount, true
}
