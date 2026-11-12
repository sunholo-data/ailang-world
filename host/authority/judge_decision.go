package authority

// Row 201 A4 sprint 2 (design §2.7): the decision judge beside the mint judge (judge.go, unchanged),
// on a read-only handle, with judge.go's step names, over a decision's subject rather than a
// MintRequestRow. The passkey judge (W15 deferred to the writer) is sprint 3's (plan S2-D1).

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
