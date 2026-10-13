package coordinator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

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

// C2 — the canonical world/decision-packet/v1 codec (row 172 D2). The byte form
// is FROZEN: compact JSON in the member order below, packetHash the SumSHA256 of
// the same object without packetHash. Amendments are /v2 semantic IDs.

// DecisionPacket is the decoded decision packet (the v1 wire, D2).
type DecisionPacket struct {
	PacketHash, ProposalHash, RequestRef        hashref.HashRef
	CreatedAt, DeadlineAt, EscalationsRemaining int64
	Policy                                      string
}

// packetWire is the wire object; member order IS the canonical order.
type packetWire struct {
	PacketHash           string `json:"packetHash"`
	ProposalHash         string `json:"proposalHash"`
	RequestRef           string `json:"requestRef"`
	CreatedAt            int64  `json:"createdAt"`
	DeadlineAt           int64  `json:"deadlineAt"`
	EscalationsRemaining int64  `json:"escalationsRemaining"`
	Policy               string `json:"policy"`
}

// packetBodyWire is packetWire without packetHash, same order and tags.
type packetBodyWire struct {
	ProposalHash         string `json:"proposalHash"`
	RequestRef           string `json:"requestRef"`
	CreatedAt            int64  `json:"createdAt"`
	DeadlineAt           int64  `json:"deadlineAt"`
	EscalationsRemaining int64  `json:"escalationsRemaining"`
	Policy               string `json:"policy"`
}

var errPacketEncoding = errors.New("coordinator: not a canonical world/decision-packet/v1 encoding")

var errPacketSchedule = errors.New("coordinator: the authorizing grant does not admit a well-formed packet schedule")

func packetBodyHash(p DecisionPacket) hashref.HashRef {
	return hashref.SumSHA256(mustJSON(packetBodyWire{ProposalHash: p.ProposalHash.String(), RequestRef: p.RequestRef.String(),
		CreatedAt: p.CreatedAt, DeadlineAt: p.DeadlineAt, EscalationsRemaining: p.EscalationsRemaining, Policy: p.Policy}))
}

// EncodeDecisionPacket returns the canonical world/decision-packet/v1 payload
// and sets p.PacketHash to the hash it carries (the body hash, D2).
func EncodeDecisionPacket(p DecisionPacket) []byte {
	p.PacketHash = packetBodyHash(p)
	return mustJSON(packetWire{PacketHash: p.PacketHash.String(), ProposalHash: p.ProposalHash.String(),
		RequestRef: p.RequestRef.String(), CreatedAt: p.CreatedAt, DeadlineAt: p.DeadlineAt,
		EscalationsRemaining: p.EscalationsRemaining, Policy: p.Policy})
}

// DecodeDecisionPacket strictly decodes a world/decision-packet/v1 payload: one
// object, no unknown members, the exact canonical bytes, and a recomputing packetHash.
func DecodeDecisionPacket(payload []byte) (DecisionPacket, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	var w packetWire
	if err := dec.Decode(&w); err != nil {
		return DecisionPacket{}, fmt.Errorf("%w: %v", errPacketEncoding, err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return DecisionPacket{}, fmt.Errorf("%w: trailing data", errPacketEncoding)
	}
	var refs [3]hashref.HashRef
	for i, text := range []string{w.PacketHash, w.ProposalHash, w.RequestRef} {
		ref, err := hashref.Parse(text)
		if err != nil {
			return DecisionPacket{}, fmt.Errorf("%w: %v", errPacketEncoding, err)
		}
		refs[i] = ref
	}
	if policyCode(w.Policy) < 0 {
		return DecisionPacket{}, fmt.Errorf("%w: unknown policy %q", errPacketEncoding, w.Policy)
	}
	p := DecisionPacket{PacketHash: refs[0], ProposalHash: refs[1], RequestRef: refs[2], CreatedAt: w.CreatedAt,
		DeadlineAt: w.DeadlineAt, EscalationsRemaining: w.EscalationsRemaining, Policy: w.Policy}
	if !bytes.Equal(mustJSON(w), payload) { // MV1 target line
		return DecisionPacket{}, fmt.Errorf("%w: not the canonical byte form", errPacketEncoding)
	}
	if packetBodyHash(p) != p.PacketHash { // MV2 target line
		return DecisionPacket{}, fmt.Errorf("%w: packetHash does not recompute", errPacketEncoding)
	}
	return p, nil
}

// mintPacket builds the world/decision-packet/v1 object for a proposal and its
// approval request. The schedule is the mint window (D3): createdAt = now,
// deadlineAt = the authorizing grant's expiry, escalationsRemaining = 0, Policy Cancel.
func mintPacket(now, grantExpiresAt int64, proposal, request hashref.HashRef, by string) (store.Object, error) {
	if !mintable(now, grantExpiresAt) || !mintedScheduleWellFormed(now, grantExpiresAt) {
		return store.Object{}, errPacketSchedule
	}
	payload := EncodeDecisionPacket(DecisionPacket{ProposalHash: proposal, RequestRef: request, CreatedAt: now,
		DeadlineAt: grantExpiresAt, EscalationsRemaining: 0, Policy: "Cancel"}) // MP6 target: "Cancel"
	return object(by, DecisionPacketV1, payload), nil
}
