package approval

// Row 201 A4 sprint 2 (design §2.3, §2.5 D11): the pure law of a decision signed by an enrolled SSH
// key: the mint law's P2-P5 and VerifySSHSIG over the 91-byte message under DecisionNamespace (the
// mint law is unchanged). Sprint 3 adds the passkey half, DecideByAssertion: the assertion law's
// W1-W14 over the decision message's hash, with W15 deferred to the credentials writer (design
// §2.7 [r2]). No I/O; now is an input.

import "bytes"

// SignedDecision is a verified key decision. Its fields are unexported, so only DecideBySignature
// and DecideByAssertion can return one with ok set; the zero value binds nothing.
type SignedDecision struct {
	ok              bool
	verb            byte
	packet, request [32]byte
	expiry          int64
	name, format    string
	key             []byte
	credentialID    []byte // webauthn: the passkey that signed
	signCount       uint32 // webauthn: the assertion's counter, not yet judged (W15 is the writer's)
}

// DecisionSignatureInput is what DecideBySignature judges: the decision's subject (verb, packet
// digest P, request digest Q, expiry), the clock, the armored SSHSIG and the enrolled approver
// its key names.
type DecisionSignatureInput struct {
	Verb            byte
	Packet, Request [32]byte
	Expiry, Now     int64
	Signature       []byte
	Approver        Approver
}

// VerifyDecisionSSHSIG parses an armored SSHSIG and verifies it over DecisionMessage(verb, packet,
// request, expiry) under DecisionNamespace with key. It returns the first error.
func VerifyDecisionSSHSIG(armored []byte, verb byte, packet, request [32]byte, expiry int64, key SSHPublicKey) error {
	s, err := ParseSSHSIG(armored)
	if err != nil {
		return err
	}
	msg, err := DecisionMessage(verb, packet, request, expiry)
	if err != nil {
		return err
	}
	return VerifySSHSIG(s, DecisionNamespace, msg, key)
}

// DecideBySignature is the decision signature law: the approver's state first (P2-P5, the mint
// law's own checks, so a revoked key says "revoked" even when its signature is good), then the
// signature over the decision message. The expiry window and the subject are the decision law's
// (D8-D10), checked before this runs.
func DecideBySignature(in DecisionSignatureInput) (*SignedDecision, error) {
	if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatSSH); err != nil { // P2-P5
		return nil, err
	}
	if err := VerifyDecisionSSHSIG(in.Signature, in.Verb, in.Packet, in.Request, in.Expiry, in.Approver.Key); err != nil {
		return nil, err
	}
	return &SignedDecision{ok: true, verb: in.Verb, packet: in.Packet, request: in.Request, expiry: in.Expiry,
		name: in.Approver.Name, format: in.Approver.Format, key: bytes.Clone(in.Approver.Key.Blob)}, nil
}

// DecisionAssertionInput is what DecideByAssertion judges: the decision's subject, the clock, the
// assertion file and the enrolled approver its credential ID names, with Passkey set.
type DecisionAssertionInput struct {
	Verb            byte
	Packet, Request [32]byte
	Expiry, Now     int64
	Response        []byte
	Approver        Approver
}

// DecideByAssertion is the decision assertion law (design §2.3, §2.7 [r2]): the approver's state
// (P2-P5, P3w), then the assertion law's W1-W14 with the challenge DecisionMessageHash(verb, packet,
// request, expiry). It does not run W15: the counter is judged in RecordKeyDecision's transaction
// against the stored counter (usePasskey), and skipped only for the exact artifact of an unlanded
// row, so a retry after a failed world write is not refused by its own first use.
func DecideByAssertion(in DecisionAssertionInput) (*SignedDecision, error) {
	if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
		return nil, err
	}
	h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry)
	ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14
	if err != nil {
		return nil, err
	}
	p := in.Approver.Passkey
	return &SignedDecision{ok: true, verb: in.Verb, packet: in.Packet, request: in.Request, expiry: in.Expiry,
		name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
}

// Binds reports whether s is a verified decision of exactly this verb, packet, request and expiry.
func (s *SignedDecision) Binds(verb byte, packet, request [32]byte, expiry int64) error {
	switch {
	case s == nil || !s.ok:
		return ErrNotVerified
	case s.verb != verb:
		return &BindingError{Field: "verb"}
	case s.packet != packet:
		return &BindingError{Field: "packet"}
	case s.request != request:
		return &BindingError{Field: "request"}
	case s.expiry != expiry:
		return &BindingError{Field: "expiry"}
	}
	return nil
}

// DecidedBy is the decision's decidedBy, format ":" name ("ssh:mark-mac"); "" unless verified.
func (s *SignedDecision) DecidedBy() string {
	if s == nil || !s.ok {
		return ""
	}
	return s.format + ":" + s.name
}

// Use is what the evidence row binds in SQL: the approver's name, format and key (a copy).
func (s *SignedDecision) Use() (name, format string, key []byte, ok bool) {
	if s == nil || !s.ok {
		return "", "", nil, false
	}
	return s.name, s.format, bytes.Clone(s.key), true
}

// PasskeyUse is what the evidence writer binds for a passkey decision: the credential ID (a copy)
// and the assertion's counter, which the writer judges (W15) against the stored one; ok only for a
// verified webauthn decision.
func (s *SignedDecision) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
	if s == nil || !s.ok || s.format != FormatWebAuthn {
		return nil, 0, false
	}
	return bytes.Clone(s.credentialID), s.signCount, true
}
