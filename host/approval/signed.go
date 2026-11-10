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
// ApproveBySignature, in this package, can return one with ok set; the zero value binds nothing.
type SignedApproval struct {
	ok        bool
	digest    [32]byte
	challenge [32]byte
	expiry    int64
	name      string
	format    string
	key       []byte
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
	if !ApproverNamePattern.MatchString(in.Approver.Name) { // P2
		return nil, ErrApproverName
	}
	if in.Approver.Format != FormatSSH { // P3
		return nil, ErrApproverFormat
	}
	if in.Approver.RevokedAt != 0 { // P4
		return nil, ErrApproverRevoked
	}
	if in.Admit.Now < in.Approver.EnrolledAt { // P5
		return nil, ErrApproverClockBehind
	}
	if err := VerifyApproval(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
		return nil, err
	}
	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: in.Approver.Format,
		key: bytes.Clone(in.Approver.Key.Blob)}, nil
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
