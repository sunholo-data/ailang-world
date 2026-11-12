package broker

// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4–§2.6): the law an operator's
// decision of a Human.Approve request passes before anything is written. Pure: the daemon reads the
// facts (the store, the clock, the request body) and this file only orders them; the first failing
// fact wins (D1–D10), and the format arm (D11) runs after it in the daemon.

import (
	"fmt"
	"strings"
	"time"
)

// MaxDecisionWindow is how far ahead of now a decision's expiry may be (D10), in seconds.
const MaxDecisionWindow = 600

// Discovery is the result of the request-wide decision-packet discovery (design §2.4 [r1]).
// Complete is true only when the keyset scan ran out before its overflow sentinel; Carrying counts
// every decoded packet naming the request; Undecodable counts packets that did not decode.
type Discovery struct {
	Complete                       bool
	Scanned, Carrying, Undecodable int
}

// DecisionFacts are what the daemon read, under humanMu, for one decision of request Ref.
type DecisionFacts struct {
	Ref            string    // the request ref as given (the reasons name it)
	RefOK          bool      // D1: sha256: and 64 lowercase hex
	RequestOK      bool      // D2: a strictly decoding world/approval-request/v1 of Human.Approve
	OnChain        bool      // D3: on the approvals chain
	Discovery      Discovery // D4a
	PacketCode     int       // D4b: the one carrying packet's ValidatePacket code
	PacketSentence string    // D4b: that code's sentence when it is not 0
	Scope          string    // D5: the request's scope
	Decided        bool      // D6: the chain holds a decision of the request
	DecidedVerb    string    // D6
	DecidedBy      string    // D6
	Now            int64     // D7, D10: the daemon's clock, Unix seconds
	Deadline       int64     // D7, D10: the packet's deadlineAt
	Packet         [32]byte  // D8: the stored packet object digest P
	Request        [32]byte  // D8: the stored request object digest Q
	PostedPacket   [32]byte  // D8: the P the caller saw
	PostedRequest  [32]byte  // D8: the Q the caller saw
	Verb           string    // D9
	Expiry         int64     // D10: the decision's expiry
}

// LawRefusal is the law's verdict: Step "" (and Status 0) accepts; otherwise the first failing
// step, its HTTP status and a reason that names the fix.
type LawRefusal struct {
	Step   string
	Status int
	Reason string
}

// DecisionWindowMax is the latest expiry a decision made at now may carry: min(now + 600, deadline).
func DecisionWindowMax(now, deadline int64) int64 {
	return min(now+MaxDecisionWindow, deadline)
}

// SubjectRefusal is the law's D1–D7: what refuses a request before any decision of it is posted.
// GET /surface/decisions/{ref} shows it; DecisionRefusal starts with it.
func SubjectRefusal(f DecisionFacts) LawRefusal {
	switch {
	case !f.RefOK:
		return LawRefusal{"D1", 400, fmt.Sprintf("%q is not a request ref (sha256: and 64 lowercase hex)", f.Ref)}
	case !f.RequestOK:
		return LawRefusal{"D2", 404, "no such approval request " + f.Ref}
	case !f.OnChain:
		return LawRefusal{"D3", 404, "request " + f.Ref + " is not an approval request of this world (it is not on the approvals chain)"}
	case !f.Discovery.Complete || f.Discovery.Undecodable > 0:
		return LawRefusal{"D4a", 409, fmt.Sprintf("the decision packets carrying request %s cannot be counted completely (more than %d stored packets, or one that does not decode); it cannot be decided here", f.Ref, DecisionDiscoveryLimit)}
	case f.Discovery.Carrying == 0:
		return LawRefusal{"D4b", 409, "no committed decision packet carries request " + f.Ref + "; it cannot be decided here"}
	case f.Discovery.Carrying != 1:
		return LawRefusal{"D4b", 409, fmt.Sprintf("request %s is carried by %d decision packets; deciding it would decide all of them; it cannot be decided here", f.Ref, f.Discovery.Carrying)}
	case f.PacketCode != 0:
		return LawRefusal{"D4b", 409, "the decision packet carrying request " + f.Ref + " is refused: " + f.PacketSentence}
	}
	if owner, reserved := ReservedApprovalScope(f.Scope); reserved {
		return LawRefusal{"D5", 403, reservedReason(owner)}
	}
	switch {
	case f.Decided:
		return LawRefusal{"D6", 409, fmt.Sprintf("request %s was already decided %s by %s", f.Ref, f.DecidedVerb, f.DecidedBy)}
	case f.Now >= f.Deadline:
		return LawRefusal{"D7", 409, "the packet's deadline passed at " + time.Unix(f.Deadline, 0).UTC().Format(time.RFC3339) + "; its policy is Cancel"}
	}
	return LawRefusal{}
}

// DecisionRefusal is the whole law, D1–D10, over a posted decision.
func DecisionRefusal(f DecisionFacts) LawRefusal {
	if r := SubjectRefusal(f); r.Step != "" {
		return r
	}
	switch {
	case f.PostedPacket != f.Packet || f.PostedRequest != f.Request:
		return LawRefusal{"D8", 409, "the request changed since it was shown; reload it"}
	case f.Verb != "approve" && f.Verb != "deny":
		return LawRefusal{"D9", 400, fmt.Sprintf("bad verb %q: a decision is approve or deny", f.Verb)}
	case f.Expiry <= f.Now || f.Expiry > DecisionWindowMax(f.Now, f.Deadline):
		return LawRefusal{"D10", 409, "the decision's expiry is stale or too far ahead (at most 10 minutes, and not after the packet's deadline)"}
	}
	return LawRefusal{}
}

// DecisionDiscoveryLimit is how many stored decision packets the discovery reads before its
// overflow sentinel (design §2.4 [r1]); one more means the set is not established.
const DecisionDiscoveryLimit = 2000

// ReservedApprovalScope reports whether scope belongs to a typed owner that is not A4's card
// (design §2.6 L1): publish-shaped scopes (world-publish approve, R4) and finish-shaped scopes
// (episode finish). Near-miss shapes fail closed: a future parser change cannot make a refused
// shape decidable here.
func ReservedApprovalScope(scope string) (owner string, reserved bool) {
	if _, err := ParsePublishApprovalScope(scope); err == nil ||
		strings.HasPrefix(scope, "registry:") || strings.Contains(scope, "effect="+EffectRegistryPublish) {
		return "publish", true
	}
	if _, err := ParseFinishApprovalScope(scope); err == nil || strings.HasPrefix(scope, "github-pr:") {
		return "finish", true
	}
	return "", false
}

func reservedReason(owner string) string {
	if owner == "publish" {
		return "publish approvals are minted by `world-publish approve` at a terminal (R4); this request cannot be decided here"
	}
	return "finish approvals are made by `ailang-worldd episode finish`; this request cannot be decided here"
}
