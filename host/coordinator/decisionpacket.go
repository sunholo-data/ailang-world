package coordinator

import "github.com/sunholo-data/ailang-world/host/hashref"

// C1 — Go mirrors of design_docs/sketches/decisionpacket.ail (row 172). Each
// law below is bound by name to its sketch law by decisionpacket_drift_test.go,
// which evaluates every inline sketch row through the mirror of the same name.

// lawRef is the sketch's Ref record: a content reference as the law sees it.
type lawRef struct{ algo, digest string }

func refOf(h hashref.HashRef) lawRef { return lawRef{h.Algo(), h.Digest()} }

// policyCode maps a TimeoutPolicy wire token to its sketch code: Cancel 0,
// EscalateBounded 1, ExecuteIfGranted 2, anything else -1.
func policyCode(token string) int {
	switch token {
	case "Cancel":
		return 0
	case "EscalateBounded":
		return 1
	case "ExecuteIfGranted":
		return 2
	}
	return -1
}

func scheduleWellFormed(createdAt, deadlineAt, escalationsRemaining int64) bool {
	return createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0
}

func mintable(now, grantExpiresAt int64) bool { return now >= 0 && now < grantExpiresAt }

func mintedScheduleWellFormed(now, grantExpiresAt int64) bool {
	return scheduleWellFormed(now, grantExpiresAt, 0)
}

func producerFieldsOk(policy int, escalationsRemaining int64) bool {
	return policy == 0 && escalationsRemaining == 0
}

func sameRef(a, b lawRef) bool { return a.algo == b.algo && a.digest == b.digest }

// packetRefusal is the validator's ordered refusal law: the code of the first
// failing fact, 0 when every fact holds.
func packetRefusal(bytesOk, carriedOk, proposalOk, requestOk, anchoredOk, lifecycleFact, pairUniqueOk bool) int {
	switch {
	case !bytesOk:
		return 1
	case !carriedOk:
		return 2
	case !proposalOk:
		return 3
	case !requestOk:
		return 4
	case !anchoredOk:
		return 5
	case !lifecycleFact:
		return 6
	case !pairUniqueOk:
		return 7
	}
	return 0
}

// samePair: two packets are duplicates only when both the request and the proposal agree.
func samePair(aRequest, aProposal, bRequest, bProposal lawRef) bool {
	return sameRef(aRequest, bRequest) && sameRef(aProposal, bProposal)
}

func lifecycleOk(createdAt, requestNow, deadlineAt, escalationsRemaining int64, policy int) bool {
	return createdAt == requestNow && scheduleWellFormed(createdAt, deadlineAt, escalationsRemaining) &&
		producerFieldsOk(policy, escalationsRemaining)
}
