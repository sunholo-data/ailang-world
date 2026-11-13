package daemon

// Row 201 A4 sprint 2 (design §2.5 D11, §2.7): the POST's key arm, in sprint 1's arm slot (S1-D9):
// after the law, under humanMu, the decision judge on a read-only handle, Binds, then a non-waiting
// writer for RecordKeyDecision (K0-K5). An SSHSIG (sprint 2) or a passkey assertion (sprint 3),
// whose counter the writer judges (K4; the judge defers W15, design §2.7 [r2]).

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

// isCounterRefusal is a passkey decision the writer's counter refused (K4, the deferred W15): a
// counter that did not advance past the stored one, or a passkey replaced since the judge.
func isCounterRefusal(err error) bool {
	return errors.Is(err, store.ErrSignCountRegression) || errors.Is(err, store.ErrPasskeyChanged)
}

// keyArm is D11 for a key format. It returns the open credentials writer (the caller closes it
// after LandKeyDecision) and the decidedBy the verified key derives, or ok false when it has
// answered a refusal: nothing is decided, and nothing is recorded unless RecordKeyDecision committed.
func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Request, f broker.DecisionFacts, body surfacepage.DecideBody, now int64) (*store.Credentials, string, bool) {
	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
	in := authority.DecisionJudgeInput{Verb: verb, Packet: f.Packet, Request: f.Request, Expiry: body.Expiry, Now: now}
	judge, what := authority.JudgeDecisionSSH, "SSH signature"
	var artifact []byte
	if body.WebAuthn != nil { // the grammar admits exactly one key format
		judge, artifact, what = authority.JudgeDecisionAssertion, []byte(*body.WebAuthn), "passkey assertion"
	} else {
		artifact = []byte(*body.SSHSIG)
	}
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return nil, "", false
	}
	signed, err := judge(ctx, ro, in, artifact)
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
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not decide request "+f.Ref+": "+err.Error())
		return nil, "", false
	}
	wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy
	if err != nil {
		d.approvalStoreError(w, r, err)
		return nil, "", false
	}
	name, format, key, _ := signed.Use()
	use := store.ApproverUse{Name: name, Format: format, PublicKey: key}
	if id, count, ok := signed.PasskeyUse(); ok { // K4: the writer judges the counter, off the exact-match retry
		use.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)}
	}
	err = wr.RecordKeyDecision(ctx, store.KeyDecisionRow{RequestRef: f.Ref, PacketRef: "sha256:" + hex.EncodeToString(f.Packet[:]), Verb: body.Verb,
		Format: format, Approver: signed.DecidedBy(), Artifact: artifact, Expiry: body.Expiry}, use, now)
	if err != nil {
		_ = wr.Close()
		switch {
		case isApproverRace(err):
			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the approver changed since the "+what+" was judged: "+err.Error())
		case isCounterRefusal(err):
			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not decide request "+f.Ref+": "+err.Error())
		default:
			d.approvalStoreError(w, r, err)
		}
		return nil, "", false
	}
	return wr, signed.DecidedBy(), true
}
