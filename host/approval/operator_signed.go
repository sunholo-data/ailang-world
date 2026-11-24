package approval

// Row 243 V4 sprint 2 M3 (design §3.1, §5.2): the pure law of an operator-action confirmation, verb_signed.go's
// shape over the operator message. ConfirmOperatorAction is the daemon's confirmation; CheckOperatorAction
// is the runner's recheck from its own frozen action (no P1, the counter already moved, the key judged at
// now). No I/O; now is an input.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

// SignedOperatorAction is a verified operator-action confirmation. Its fields are unexported, so only
// ConfirmOperatorAction and CheckOperatorAction return one with ok set; the zero value binds nothing.
type SignedOperatorAction struct {
	ok                bool
	digest, challenge [32]byte
	expiry            int64
	name, format      string
	credentialID      []byte
	signCount         uint32
	confirmation      *OperatorApprovalBinding
	artifactHash      [32]byte
}

// OperatorConfirmInput is what ConfirmOperatorAction judges: the stored action (strictly decoded), the
// row's admission input (StoredDigest is the row's digest column, BytesDigest the SHA-256 of its stored
// bytes), its challenge, the artifact, its format and the enrolled approver it names.
type OperatorConfirmInput struct {
	Action    OperatorAction
	Admit     AdmitInput
	Challenge [32]byte
	Artifact  []byte
	Format    string // FormatSSH or FormatWebAuthn
	Approver  Approver
}

// OperatorCheckInput is what CheckOperatorAction rechecks: the runner's own frozen action, the
// challenge and expiry it recorded, the clock, the stored artifact, its format and the approver.
type OperatorCheckInput struct {
	Action      OperatorAction
	Challenge   [32]byte
	Expiry, Now int64
	Artifact    []byte
	Format      string
	Approver    Approver
}

var (
	ErrOperatorLifetime = errors.New("approval: the action's request outlives 600 s")
	ErrOperatorCounter  = errors.New("approval: the confirmation's passkey counter was never recorded")
)

// ConfirmOperatorAction is the operator confirmation law, in this order: the action is valid (its
// signature digest), P1 admission (consumed, clock behind, expiry at or before now, stored bytes), the
// row's digest is this action's, the lifetime is at most OperatorActionLifetime, P2-P5, then the
// signature over the operator message, and for a passkey W1-W14 and W15.
func ConfirmOperatorAction(in OperatorConfirmInput) (*SignedOperatorAction, error) {
	d, err := OperatorSignatureDigest(in.Action)
	if err != nil {
		return nil, err
	}
	if v := Admit(in.Admit); v != Admitted { // P1
		return nil, &NotAdmittedError{Verdict: v}
	}
	if in.Admit.StoredDigest != d { // C1
		return nil, &BindingError{Field: "digest"}
	}
	if in.Admit.ExpiresAt-in.Admit.CreatedAt > int64(OperatorActionLifetime/time.Second) { // C2
		return nil, ErrOperatorLifetime
	}
	return operatorProof(d, in.Challenge, in.Admit.ExpiresAt, in.Admit.Now, in.Artifact, in.Format, in.Approver, true)
}

// CheckOperatorAction is the runner's recheck before an effect (design §3.1): its own action's digest,
// the recorded challenge and expiry, expiry at or before now refused (R1), P2-P5 at now (a key revoked
// since is refused), the signature, and for a passkey W1-W14 and a counter the confirmation recorded
// (R2: the stored counter is at or past the artifact's; never advanced again).
func CheckOperatorAction(in OperatorCheckInput) (*SignedOperatorAction, error) {
	d, err := OperatorSignatureDigest(in.Action)
	if err != nil {
		return nil, err
	}
	if in.Now >= in.Expiry { // R1
		return nil, &NotAdmittedError{Verdict: RefuseExpired}
	}
	return operatorProof(d, in.Challenge, in.Expiry, in.Now, in.Artifact, in.Format, in.Approver, false)
}

// operatorProof is P2-P6 / W1-W15 for both laws; confirm selects W15 (the counter must move) over R2.
func operatorProof(d, challenge [32]byte, expiry, now int64, artifact []byte, format string, ap Approver, confirm bool) (*SignedOperatorAction, error) {
	if format != FormatSSH && format != FormatWebAuthn {
		return nil, fmt.Errorf("%w: %q", ErrVerbFormat, format)
	}
	if err := checkApprover(AdmitInput{Now: now}, ap, format); err != nil { // P2-P5, P3w
		return nil, err
	}
	s := &SignedOperatorAction{ok: true, digest: d, challenge: challenge, expiry: expiry, name: ap.Name, format: format}
	if confirm {
		s.confirmation = &OperatorApprovalBinding{Name: ap.Name, Format: format, PublicKey: string(ap.Key.Blob)}
		s.artifactHash = sha256.Sum256(artifact)
	}
	if format == FormatSSH {
		sig, err := ParseSSHSIG(artifact)
		if err != nil {
			return nil, err
		}
		msg, err := OperatorActionMessage(d, challenge, expiry)
		if err != nil {
			return nil, err
		}
		if err := VerifySSHSIG(sig, OperatorActionNamespace, msg, ap.Key); err != nil { // P6
			return nil, err
		}
		return s, nil
	}
	h, hErr := OperatorActionMessageHash(d, challenge, expiry)
	ad, err := assertionSteps(artifact, ap, h, hErr) // W1-W14
	if err != nil {
		return nil, err
	}
	p := ap.Passkey
	if confirm && (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
			ErrSignCountRegression, p.SignCount, ad.SignCount, ap.Name)
	}
	if !confirm && ad.SignCount > p.SignCount { // R2
		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d", ErrOperatorCounter, p.SignCount, ad.SignCount)
	}
	if confirm {
		s.confirmation.PublicKey = string(p.Key.Canonical())
		s.confirmation.RPID, s.confirmation.Origin, s.confirmation.BackupEligible = p.RPID, p.Origin, p.BackupEligible
	}
	s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCount
	return s, nil
}

// Binds reports whether s is a verified confirmation of exactly this action, challenge and expiry.
func (s *SignedOperatorAction) Binds(a OperatorAction, challenge [32]byte, expiry int64) error {
	if s == nil || !s.ok {
		return ErrNotVerified
	}
	d, err := OperatorSignatureDigest(a)
	switch {
	case err != nil:
		return err
	case s.digest != d:
		return &BindingError{Field: "digest"}
	case s.challenge != challenge:
		return &BindingError{Field: "challenge"}
	case s.expiry != expiry:
		return &BindingError{Field: "expiry"}
	}
	return nil
}

// ApprovedBy is format ":" name ("webauthn:mark-mac"); "" unless verified.
func (s *SignedOperatorAction) ApprovedBy() string {
	if s == nil || !s.ok {
		return ""
	}
	return s.format + ":" + s.name
}

// PasskeyUse is what S3's confirm CAS binds for a passkey: the credential ID (a copy) and the
// assertion's counter; ok only for a verified webauthn confirmation.
func (s *SignedOperatorAction) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
	if s == nil || !s.ok || s.format != FormatWebAuthn {
		return nil, 0, false
	}
	return bytes.Clone(s.credentialID), s.signCount, true
}

// OperatorApprovalBinding describes the verified enrollment snapshot for the
// store's transactional comparison. It is obtained from a signed proof.
type OperatorApprovalBinding struct {
	Name, Format, PublicKey, RPID, Origin string
	BackupEligible                        bool
}

// ConfirmationFor binds the proof to the exact artifact to be stored and
// distinguishes initial confirmation from a runner recheck.
func (s *SignedOperatorAction) ConfirmationFor(artifact []byte) (OperatorApprovalBinding, error) {
	if s == nil || !s.ok || s.confirmation == nil || sha256.Sum256(artifact) != s.artifactHash {
		return OperatorApprovalBinding{}, ErrNotVerified
	}
	return *s.confirmation, nil
}
