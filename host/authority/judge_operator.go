package authority

// Row 243 V4 sprint 2b (design §3.1 [r1], §5.2, §6; parked from sprint 2 by its size cut, S2-D1): the
// operator-action judges beside the verb judges, on a read-only handle with judge.go's step names.
// Neither judge takes an approver: each reads the enrolled approver the artifact names from ro
// (sshApprover, passkeyApprover), so the law judges the store's key, format, enrolment, revocation and
// counter, never a caller's. VerifiedOperatorAction is the typed verified-action proof: its fields
// are unexported and only JudgeOperatorAction and RecheckOperatorAction construct one.

import (
	"context"
	"errors"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// VerifiedOperatorAction is a judged operator-action proof; the zero value binds nothing.
type VerifiedOperatorAction struct {
	signed    *approval.SignedOperatorAction
	rechecked bool
}

func verifiedOperatorAction(s *approval.SignedOperatorAction, rechecked bool) *VerifiedOperatorAction {
	return &VerifiedOperatorAction{signed: s, rechecked: rechecked}
}

// ErrOperatorApprover is a recheck whose artifact's approver is not the one the action recorded.
var ErrOperatorApprover = errors.New("authority: the action's recorded approver did not sign it")

// OperatorJudgeInput is a stored operator action at now (S3's row): its strictly decoded envelope, the
// digest column, the challenge, the SHA-256 of the stored bytes, its times and whether it is consumed.
type OperatorJudgeInput struct {
	Action                         approval.OperatorAction
	Digest, Challenge, BytesDigest [32]byte
	CreatedAt, ExpiresAt, Now      int64
	Consumed                       bool
}

func (in OperatorJudgeInput) admit() approval.AdmitInput {
	return approval.AdmitInput{Now: in.Now, CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed: in.Consumed,
		StoredDigest: in.Digest, BytesDigest: in.BytesDigest}
}

// operatorApprover is the artifact's format (an armored SSHSIG, else a passkey assertion) and the
// enrolled approver it names on ro, or a *JudgeError.
func operatorApprover(ctx context.Context, ro *store.Credentials, artifact []byte) (string, approval.Approver, error) {
	if sig, sigErr := approval.ParseSSHSIG(artifact); sigErr == nil {
		a, err := sshApprover(ctx, ro, sig)
		return approval.FormatSSH, a, err
	} else if resp, respErr := approval.ParseAssertionResponse(artifact); respErr == nil {
		a, err := passkeyApprover(ctx, ro, resp)
		return approval.FormatWebAuthn, a, err
	} else {
		return "", approval.Approver{}, &JudgeError{Step: JudgeStepParse, Err: errors.Join(ErrUnknownArtifact, sigErr, respErr)}
	}
}

// JudgeOperatorAction is the daemon's confirmation judge (S7): the approver the artifact names on ro,
// then approval.ConfirmOperatorAction over in. It returns the proof or a *JudgeError.
func JudgeOperatorAction(ctx context.Context, ro *store.Credentials, in OperatorJudgeInput, artifact []byte) (*VerifiedOperatorAction, error) {
	format, a, err := operatorApprover(ctx, ro, artifact)
	if err != nil {
		return nil, err
	}
	s, err := approval.ConfirmOperatorAction(approval.OperatorConfirmInput{Action: in.Action, Admit: in.admit(), Challenge: in.Challenge,
		Artifact: artifact, Format: format, Approver: a})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return verifiedOperatorAction(s, false), nil
}

// RecheckOperatorAction is the runner's recheck before an effect (design §3.1): its own frozen action,
// the recorded challenge and expiry, the current key on ro at now, then the recorded approver identity.
// The lifetime bound is the daemon's (C2 at confirmation); the runner rechecks expiry and revocation.
func RecheckOperatorAction(ctx context.Context, ro *store.Credentials, action approval.OperatorAction, challenge [32]byte, expiry, now int64,
	artifact []byte, approvedBy string) (*VerifiedOperatorAction, error) {
	format, a, err := operatorApprover(ctx, ro, artifact)
	if err != nil {
		return nil, err
	}
	s, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: action, Challenge: challenge, Expiry: expiry, Now: now,
		Artifact: artifact, Format: format, Approver: a})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	if s.ApprovedBy() != approvedBy {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: ErrOperatorApprover}
	}
	return verifiedOperatorAction(s, true), nil
}

// Binds, ApprovedBy and PasskeyUse are the law's, nil-safe; Rechecked says a runner recheck made it.
func (v *VerifiedOperatorAction) Binds(a approval.OperatorAction, challenge [32]byte, expiry int64) error {
	if v == nil {
		return approval.ErrNotVerified
	}
	return v.signed.Binds(a, challenge, expiry)
}

func (v *VerifiedOperatorAction) ApprovedBy() string {
	if v == nil {
		return ""
	}
	return v.signed.ApprovedBy()
}

func (v *VerifiedOperatorAction) PasskeyUse() ([]byte, uint32, bool) {
	if v == nil {
		return nil, 0, false
	}
	return v.signed.PasskeyUse()
}

func (v *VerifiedOperatorAction) Rechecked() bool { return v != nil && v.signed != nil && v.rechecked }

// OperatorConfirmation transports the sealed approval law proof to the store
// without making store import authority. The store rechecks its enrollment.
func (v *VerifiedOperatorAction) OperatorConfirmation() *approval.SignedOperatorAction {
	if v == nil {
		return nil
	}
	return v.signed
}

// RecheckClaimedOperatorAction uses the runner's frozen inputs. The receipt
// supplies only the artifact and the identity to compare after verification.
func RecheckClaimedOperatorAction(ctx context.Context, ro *store.Credentials, frozen approval.OperatorAction, challenge [32]byte,
	expiry, now int64, row store.OperatorActionRow) (*VerifiedOperatorAction, error) {
	return RecheckOperatorAction(ctx, ro, frozen, challenge, expiry, now, row.Artifact, row.ApprovedBy)
}
