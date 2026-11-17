package authority

// Row 243 V3 sprint 1 M2 (design §2.3): the verb judges beside the mint and decision judges, on a
// read-only handle, with judge.go's step names. JudgeVerbSSH and JudgeVerbAssertion are the daemon's
// C11 over a stored verb request's primitives (no store type: sprint 1 adds none); RecheckVerb is
// the CLI's re-check of a recorded confirmation (approval.CheckConfirmation, design §2.7 step 6(a)).

import (
	"context"
	"errors"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// VerbJudgeInput is a stored verb request at now: its digest column and challenge, the SHA-256 of
// its stored bytes, its times and whether it is consumed.
type VerbJudgeInput struct {
	Digest, Challenge, BytesDigest [32]byte
	CreatedAt, ExpiresAt, Now      int64
	Consumed                       bool
}

func (in VerbJudgeInput) admit() approval.AdmitInput {
	return approval.AdmitInput{Now: in.Now, CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed: in.Consumed,
		StoredDigest: in.Digest, BytesDigest: in.BytesDigest}
}

// sshApprover is the enrolled approver an SSHSIG's key names (ApproverByKey on ro), or a *JudgeError.
func sshApprover(ctx context.Context, ro *store.Credentials, sig *approval.SSHSig) (approval.Approver, error) {
	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
	if err != nil {
		return approval.Approver{}, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return approval.Approver{}, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint: sig.PublicKey.Fingerprint(), Err: ErrNotEnrolled}
	}
	return approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
		EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt}, nil
}

// passkeyApprover is the enrolled approver, with its passkey, an assertion's credential ID names
// (PasskeyByCredentialID on ro, ParseStoredCOSEKey), or a *JudgeError.
func passkeyApprover(ctx context.Context, ro *store.Credentials, resp approval.AssertionResponse) (approval.Approver, error) {
	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
	if err != nil {
		return approval.Approver{}, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return approval.Approver{}, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
	}
	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
	if err != nil {
		return approval.Approver{}, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
	}
	return approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
		Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
			BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}, nil
}

// JudgeVerbSSH judges an armored SSHSIG confirmation of in: ParseSSHSIG, the approver its key names,
// then approval.ConfirmBySignature. It returns the verified confirmation or a *JudgeError.
func JudgeVerbSSH(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, armored []byte) (*approval.SignedVerb, error) {
	sig, err := approval.ParseSSHSIG(armored)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	a, err := sshApprover(ctx, ro, sig)
	if err != nil {
		return nil, err
	}
	signed, err := approval.ConfirmBySignature(approval.VerbSignatureInput{Admit: in.admit(), Challenge: in.Challenge, Signature: armored, Approver: a})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}

// JudgeVerbAssertion judges a passkey confirmation of in: ParseAssertionResponse, the passkey its
// credential ID names, its stored key, then approval.ConfirmByAssertion (W15 in full). It returns the
// verified confirmation or a *JudgeError.
func JudgeVerbAssertion(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, response []byte) (*approval.SignedVerb, error) {
	resp, err := approval.ParseAssertionResponse(response)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	a, err := passkeyApprover(ctx, ro, resp)
	if err != nil {
		return nil, err
	}
	signed, err := approval.ConfirmByAssertion(approval.VerbAssertionInput{Admit: in.admit(), Challenge: in.Challenge, Response: response, Approver: a})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}

// ErrUnknownArtifact is a recorded artifact that is neither an armored SSHSIG nor an assertion.
var ErrUnknownArtifact = errors.New("authority: the artifact is neither an SSH signature nor a passkey assertion")

// RecheckVerb is the CLI's re-check of a recorded confirmation (design §2.3 [r1], §2.7 step 6(a)):
// the artifact's kind (an armored SSHSIG, else a passkey assertion), the enrolled approver it names
// on ro, then approval.CheckConfirmation over in's digest, challenge and expiry at in.Now — the
// caller's own values, never the row's. in.CreatedAt, BytesDigest and Consumed are not read: the
// row is consumed by then (no P1) and the counter has moved (no W15). It returns the verified
// confirmation or a *JudgeError.
func RecheckVerb(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, artifact []byte) (*approval.SignedVerb, error) {
	check := approval.VerbCheckInput{Digest: in.Digest, Challenge: in.Challenge, Expiry: in.ExpiresAt, Now: in.Now, Artifact: artifact}
	if sig, sigErr := approval.ParseSSHSIG(artifact); sigErr == nil {
		a, err := sshApprover(ctx, ro, sig)
		if err != nil {
			return nil, err
		}
		check.Format, check.Approver = approval.FormatSSH, a
	} else if resp, respErr := approval.ParseAssertionResponse(artifact); respErr == nil {
		a, err := passkeyApprover(ctx, ro, resp)
		if err != nil {
			return nil, err
		}
		check.Format, check.Approver = approval.FormatWebAuthn, a
	} else {
		return nil, &JudgeError{Step: JudgeStepParse, Err: errors.Join(ErrUnknownArtifact, sigErr, respErr)}
	}
	signed, err := approval.CheckConfirmation(check)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}
