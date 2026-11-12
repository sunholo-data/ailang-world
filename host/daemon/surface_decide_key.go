package daemon

// Row 201 A4 sprint 2 (design §2.5 D11, §2.7): the POST's key arm, in sprint 1's arm slot (S1-D9):
// after the law, under humanMu, the decision judge on a read-only handle, Binds, then a non-waiting
// writer for RecordKeyDecision (K0-K5). SSHSIG only; the passkey arm is sprint 3's (plan S2-D1).

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// isApproverRace is a key decision's approver that changed since the read-only judge: revoked,
// re-keyed, gone or enrolled after now (K3, in the writer's transaction).
func isApproverRace(err error) bool {
	return errors.Is(err, store.ErrApproverRevoked) || errors.Is(err, store.ErrApproverKeyChanged) ||
		errors.Is(err, store.ErrApproverUnknown) || errors.Is(err, store.ErrApproverClockBehind)
}

// keyArm is D11 for a key format. It returns the open credentials writer (the caller closes it
// after LandKeyDecision) and the decidedBy the verified key derives, or ok false when it has
// answered a refusal: nothing is decided, and nothing is recorded unless RecordKeyDecision committed.
func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Request, f broker.DecisionFacts, body surfacepage.DecideBody, now int64) (*store.Credentials, string, bool) {
	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
	artifact := []byte(*body.SSHSIG)
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return nil, "", false
	}
	signed, err := authority.JudgeDecisionSSH(ctx, ro, authority.DecisionJudgeInput{Verb: verb, Packet: f.Packet, Request: f.Request, Expiry: body.Expiry, Now: now}, artifact)
	_ = ro.Close()
	if err == nil {
		err = signed.Binds(verb, f.Packet, f.Request, body.Expiry)
	}
	var je *authority.JudgeError
	switch {
	case errors.As(err, &je) && je.Step == authority.JudgeStepStore:
		d.approvalStoreError(w, r, err)
		return nil, "", false
	case err != nil:
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the SSH signature does not decide request "+f.Ref+": "+err.Error())
		return nil, "", false
	}
	wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy
	if err != nil {
		d.approvalStoreError(w, r, err)
		return nil, "", false
	}
	name, format, key, _ := signed.Use()
	err = wr.RecordKeyDecision(ctx, store.KeyDecisionRow{RequestRef: f.Ref, PacketRef: "sha256:" + hex.EncodeToString(f.Packet[:]), Verb: body.Verb,
		Format: format, Approver: signed.DecidedBy(), Artifact: artifact, Expiry: body.Expiry}, store.ApproverUse{Name: name, Format: format, PublicKey: key}, now)
	if err != nil {
		_ = wr.Close()
		if isApproverRace(err) {
			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the approver changed since the signature was judged: "+err.Error())
		} else {
			d.approvalStoreError(w, r, err)
		}
		return nil, "", false
	}
	return wr, signed.DecidedBy(), true
}
