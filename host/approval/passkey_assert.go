package approval

// Row 201 A2 (design §2.9): the assertion law. A mint request is approved by a passkey only when
// it is admitted, the approver is an active webauthn approver, and its passkey signed, with user
// verification, this request's message hash at the enrolled origin and relying party ID, with a
// counter that moved forward. No I/O.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

// AssertionInput is what ApproveByAssertion judges.
type AssertionInput struct {
	Admit     AdmitInput // the existing admission input (row times, consumed, both digests)
	Challenge [32]byte   // the request row's challenge
	Response  []byte     // the assertion file
	Approver  Approver   // the enrolled row the assertion's credential ID names, with Passkey set
}

var (
	ErrApproverNotPasskey  = errors.New("approval: the approver is not a webauthn approver with an enrolled passkey")
	ErrUserHandle          = errors.New("approval: the assertion's userHandle is not the approver's name")
	ErrBackupEligibility   = errors.New("approval: the assertion's backup eligibility (BE) is not the enrolled passkey's")
	ErrSignCountRegression = errors.New("approval: the passkey's sign counter did not advance (a cloned authenticator, or an older assertion)")
)

// ApproveByAssertion is the assertion law (P1-P5, P3w, W1-W15, in this order): admission and the
// approver's state first, the credential ID before any signed content is parsed, and the counter
// last, so it is judged only on a genuine signature.
func ApproveByAssertion(in AssertionInput) (*SignedApproval, error) {
	if v := Admit(in.Admit); v != Admitted { // P1
		return nil, &NotAdmittedError{Verdict: v}
	}
	if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
		return nil, err
	}
	p := in.Approver.Passkey
	resp, err := ParseAssertionResponse(in.Response) // W1
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(resp.RawID, p.CredentialID) { // W2
		return nil, ErrCredentialIDMismatch
	}
	if resp.HasUserHandle && !bytes.Equal(resp.UserHandle, []byte(in.Approver.Name)) { // W3
		return nil, ErrUserHandle
	}
	cd, err := ParseClientData(resp.ClientDataJSON) // W4
	if err != nil {
		return nil, err
	}
	if cd.Type != "webauthn.get" { // W5
		return nil, ErrClientDataType
	}
	h, err := MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
	if err != nil || cd.Challenge != b64Challenge(h[:]) { // W6
		return nil, ErrChallenge
	}
	if cd.Origin != p.Origin { // W7
		return nil, ErrOriginMismatch
	}
	if cd.CrossOrigin || cd.HasTopOrigin { // W8
		return nil, ErrCrossOrigin
	}
	ad, err := ParseAuthData(resp.AuthenticatorData, false) // W9
	if err != nil {
		return nil, err
	}
	if ad.RPIDHash != sha256.Sum256([]byte(p.RPID)) { // W10
		return nil, ErrRPIDHash
	}
	if !ad.UP { // W11
		return nil, ErrUserPresence
	}
	if !ad.UV { // W12
		return nil, ErrUserVerification
	}
	if ad.BE != p.BackupEligible { // W13
		return nil, ErrBackupEligibility
	}
	cdh := sha256.Sum256(resp.ClientDataJSON)
	if err := p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature); err != nil { // W14
		return nil, err
	}
	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
	}
	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(),
		credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
}
