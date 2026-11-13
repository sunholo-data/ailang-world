package authority

// Row 201 A4 sprint 2 (design §2.7): the decision judge beside the mint judge (judge.go, unchanged),
// on a read-only handle, with judge.go's step names, over a decision's subject rather than a
// MintRequestRow. JudgeDecisionAssertion (sprint 3) defers the counter (W15) to the credentials
// writer's transaction (design §2.7 [r2]).

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// DecisionJudgeInput is a decision's subject at now: the verb byte, the packet digest P, the
// request digest Q and the signed expiry.
type DecisionJudgeInput struct {
	Verb            byte
	Packet, Request [32]byte
	Expiry, Now     int64
}

// JudgeDecisionSSH judges an armored SSHSIG decision: ParseSSHSIG, the approver its key names
// (ApproverByKey on ro), then approval.DecideBySignature. It returns the verified decision or a
// *JudgeError.
func JudgeDecisionSSH(ctx context.Context, ro *store.Credentials, in DecisionJudgeInput, armored []byte) (*approval.SignedDecision, error) {
	sig, err := approval.ParseSSHSIG(armored)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return nil, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint: sig.PublicKey.Fingerprint(), Err: ErrNotEnrolled}
	}
	signed, err := approval.DecideBySignature(approval.DecisionSignatureInput{Verb: in.Verb, Packet: in.Packet, Request: in.Request,
		Expiry: in.Expiry, Now: in.Now, Signature: armored,
		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt}})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}

// JudgeDecisionAssertion judges a passkey decision assertion: ParseAssertionResponse, the passkey its
// credential ID names (PasskeyByCredentialID on ro), ParseStoredCOSEKey, then
// approval.DecideByAssertion, which runs W1-W14 and not W15: the counter is judged by
// RecordKeyDecision against the stored counter, off the exact-match retry. It returns the verified
// decision or a *JudgeError, with judge.go's step names.
func JudgeDecisionAssertion(ctx context.Context, ro *store.Credentials, in DecisionJudgeInput, response []byte) (*approval.SignedDecision, error) {
	resp, err := approval.ParseAssertionResponse(response)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
	}
	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
	}
	signed, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: in.Verb, Packet: in.Packet, Request: in.Request,
		Expiry: in.Expiry, Now: in.Now, Response: response,
		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}
