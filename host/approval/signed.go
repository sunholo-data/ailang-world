package approval

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.3): the pure signature
// law. A mint request is approved by signature only when it is admitted, the approver is
// enrolled, active and not ahead of the clock, and its key's SSHSIG over the request's exact
// message (stored digest, row challenge, row expiry) verifies. No I/O here (coding standard S2).

import (
	"bytes"
	"errors"
	"fmt"
)

// SignedApproval is a verified signature approval. Its fields are unexported, so only
// ApproveBySignature or ApproveByAssertion, in this package, can return one with ok set; the zero
// value binds nothing.
type SignedApproval struct {
	ok           bool
	digest       [32]byte
	challenge    [32]byte
	expiry       int64
	name         string
	format       string
	key          []byte
	credentialID []byte // webauthn only (row 201 A2): the passkey that signed
	signCount    uint32 // webauthn only: the assertion's sign counter
}

// SignatureInput is what ApproveBySignature judges.
type SignatureInput struct {
	Admit     AdmitInput // the existing admission input (row times, consumed, both digests)
	Challenge [32]byte   // the request row's challenge
	Signature []byte     // the armored SSHSIG
	Approver  Approver   // the enrolled row the signature's key names
}

// NotAdmittedError is admission's refusal on the signature path (P1): the request is consumed,
// expired, ahead of the clock, or its bytes do not match their digest.
type NotAdmittedError struct{ Verdict Verdict }

func (e *NotAdmittedError) Error() string {
	return fmt.Sprintf("approval: the request is not admitted (%s)", e.Verdict)
}

// ErrNotVerified is any use of a SignedApproval that ApproveBySignature did not return.
var ErrNotVerified = errors.New("approval: not a verified signature approval")

// BindingError is a verified approval used for a request it was not made for.
type BindingError struct{ Field string } // "digest" | "challenge" | "expiry"

func (e *BindingError) Error() string {
	return fmt.Sprintf("approval: the signature approval was made for another request (%s differs)", e.Field)
}

// ApproveBySignature is the signature law (P1-P6, in this order): admission first, so a replay
// says "already used"; the approver's state before the cryptography, so a revoked key says
// "revoked" even when its signature is good.
func ApproveBySignature(in SignatureInput) (*SignedApproval, error) {
	if v := Admit(in.Admit); v != Admitted { // P1
		return nil, &NotAdmittedError{Verdict: v}
	}
	if err := checkApprover(in.Admit, in.Approver, FormatSSH); err != nil { // P2-P5
		return nil, err
	}
	if err := VerifyApproval(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
		return nil, err
	}
	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: in.Approver.Format,
		key: bytes.Clone(in.Approver.Key.Blob)}, nil
}

// checkApprover is P2-P5, shared by both laws (row 201 A2 §2.9): the name, the format (a webauthn
// approver must carry its passkey, P3w), not revoked, and not enrolled after now.
func checkApprover(admit AdmitInput, a Approver, format string) error {
	if !ApproverNamePattern.MatchString(a.Name) { // P2
		return ErrApproverName
	}
	if format == FormatWebAuthn && (a.Format != format || a.Passkey == nil) { // P3w
		return ErrApproverNotPasskey
	}
	if a.Format != format { // P3
		return ErrApproverFormat
	}
	if a.RevokedAt != 0 { // P4
		return ErrApproverRevoked
	}
	if admit.Now < a.EnrolledAt { // P5
		return ErrApproverClockBehind
	}
	return nil
}

// Binds reports whether s is a verified approval of exactly this digest, challenge and expiry.
func (s *SignedApproval) Binds(digest, challenge [32]byte, expiry int64) error {
	if s == nil || !s.ok { // B1
		return ErrNotVerified
	}
	if s.digest != digest { // B2
		return &BindingError{Field: "digest"}
	}
	if s.challenge != challenge { // B3
		return &BindingError{Field: "challenge"}
	}
	if s.expiry != expiry { // B4
		return &BindingError{Field: "expiry"}
	}
	return nil
}

// ---- M4b ----

// ApprovedBy is the approved_by value, format ":" name ("ssh:mark-mac"); "" unless verified.
func (s *SignedApproval) ApprovedBy() string {
	if s == nil || !s.ok { // B5
		return ""
	}
	return s.format + ":" + s.name
}

// Use is what the consume binds in SQL: the approver's name, format and key (a copy).
func (s *SignedApproval) Use() (name, format string, key []byte, ok bool) {
	if s == nil || !s.ok { // B6
		return "", "", nil, false
	}
	return s.name, s.format, bytes.Clone(s.key), true
}

// PasskeyUse is what the consume binds for a passkey (row 201 A2 §2.9): the credential ID (a copy)
// and the assertion's counter; ok only for a verified webauthn approval.
func (s *SignedApproval) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
	if s == nil || !s.ok || s.format != FormatWebAuthn { // B7, B8
		return nil, 0, false
	}
	return bytes.Clone(s.credentialID), s.signCount, true
}
