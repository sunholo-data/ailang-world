package coordinator

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// PacketReader is the read surface ValidatePacket needs; the daemon's readStore satisfies it.
type PacketReader interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error)
	GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)
	ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error)
	ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error)
}

// PacketResult is ValidatePacket's verdict: Code is packetRefusal's (0 = validated).
type PacketResult struct {
	Code            int
	Entry           int64                  // the packet's first carrying entry (Code 0 or ≥ 3)
	Packet          DecisionPacket         // Code != 1
	TransitionFn    hashref.HashRef        // the record's transitionFn (Code 0 or ≥ 3)
	Approval        broker.ApprovalSummary // the chain's summary for RequestRef (Code 0 or ≥ 6)
	AnchorTruncated bool                   // code 5 failed and the 40-head read was truncated
	MalformedRef    string                 // code 5 failed on a malformed approvals chain
	UniqueTruncated bool                   // code 7's scan saw more than packetScanLimit packets
}

const packetScanLimit = 100

// ValidatePacket checks a stored world/decision-packet/v1 object against the seven
// refusal checks in order (short-circuit). Store errors are returned, never a code.
func ValidatePacket(ctx context.Context, r PacketReader, obj store.Object) (PacketResult, error) {
	var res PacketResult
	var f [7]bool
	refuse := func() (PacketResult, error) {
		res.Code = packetRefusal(f[0], f[1], f[2], f[3], f[4], f[5], f[6])
		return res, nil
	}

	// 1. bytes
	p, err := DecodeDecisionPacket(obj.Payload)
	if err != nil || obj.SemanticID != DecisionPacketV1 || obj.InterfaceHash != hashref.SumSHA256([]byte(DecisionPacketV1)) {
		return refuse()
	}
	f[0] = true
	res.Packet = p

	// 2. carried
	es, err := r.ObjectCommits(ctx, obj.Hash, -1, 1)
	if err != nil {
		return PacketResult{}, err
	}
	if len(es) == 0 {
		return refuse()
	}
	e := es[0]
	entry, ok, err := r.GetLogEntry(ctx, e)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok {
		return refuse()
	}
	recObj, ok, err := r.GetObject(ctx, entry.TransitionRef)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok || recObj.SemanticID != RecordV2 {
		return refuse()
	}
	var rec recordV2
	if json.Unmarshal(recObj.Payload, &rec) != nil || !wellFormedV2(rec) {
		return refuse()
	}
	tf, err := hashref.Parse(rec.TransitionFn)
	if err != nil {
		return refuse()
	}
	rcs, err := r.ObjectCommits(ctx, recObj.Hash, -1, 1)
	if err != nil {
		return PacketResult{}, err
	}
	if len(rcs) == 0 || rcs[0] != e {
		return refuse()
	}
	f[1] = true
	res.Entry = e
	res.TransitionFn = tf

	// 3. proposal
	if rec.Plan != p.ProposalHash.String() {
		return refuse()
	}
	planObj, ok, err := r.GetObject(ctx, p.ProposalHash)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok || planObj.SemanticID != EffectPlanV1 {
		return refuse()
	}
	f[2] = true

	// 4. request
	if len(rec.Effects) != 1 {
		return refuse()
	}
	erRef, err := hashref.Parse(rec.Effects[0])
	if err != nil {
		return refuse()
	}
	erObj, ok, err := r.GetObject(ctx, erRef)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok || erObj.SemanticID != broker.EffectRecordV1 {
		return refuse()
	}
	er, err := broker.DecodeRecord(erObj.Payload)
	if err != nil {
		return refuse()
	}
	if er.Effect != broker.EffectHumanApprove || !er.Allowed || er.Failed || er.ResultRef.IsZero() {
		return refuse()
	}
	resObj, ok, err := r.GetObject(ctx, er.ResultRef)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok || resObj.SemanticID != broker.EffectResultV1 {
		return refuse()
	}
	pending, err := broker.ParsePendingApproval(resObj.Payload)
	if err != nil {
		return refuse()
	}
	if pending != p.RequestRef {
		return refuse()
	}
	reqObj, ok, err := r.GetObject(ctx, p.RequestRef)
	if err != nil {
		return PacketResult{}, err
	}
	if !ok || reqObj.SemanticID != broker.ApprovalRequestV1 {
		return refuse()
	}
	req, err := broker.DecodeApprovalRequest(reqObj.Payload)
	if err != nil {
		return refuse()
	}
	if req.Effect != broker.EffectHumanApprove || req.Scope != er.Scope {
		return refuse()
	}
	f[3] = true

	// 5. anchored
	page, err := broker.RecentApprovals(ctx, r, 40)
	if errors.Is(err, broker.ErrApprovalChainMalformed) {
		res.MalformedRef = page.MalformedRef
		return refuse()
	}
	if err != nil {
		return PacketResult{}, err
	}
	found := false
	for _, s := range page.Summaries {
		if s.RequestRef == p.RequestRef.String() {
			found, res.Approval = true, s
			break
		}
	}
	if !found {
		res.AnchorTruncated = page.Truncated
		return refuse()
	}
	f[4] = true

	// 6. lifecycle
	if !lifecycleOk(p.CreatedAt, req.Now, p.DeadlineAt, p.EscalationsRemaining, policyCode(p.Policy)) {
		return refuse()
	}
	f[5] = true

	// 7. pair-unique
	others, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, "", packetScanLimit+1)
	if err != nil {
		return PacketResult{}, err
	}
	res.UniqueTruncated = len(others) > packetScanLimit
	for _, o := range others[:min(len(others), packetScanLimit)] {
		if o.Hash == obj.Hash {
			continue
		}
		full, ok, err := r.GetObject(ctx, o.Hash)
		if err != nil {
			return PacketResult{}, err
		}
		if !ok {
			continue
		}
		q, err := DecodeDecisionPacket(full.Payload)
		if err != nil {
			continue
		}
		if samePair(refOf(p.RequestRef), refOf(p.ProposalHash), refOf(q.RequestRef), refOf(q.ProposalHash)) {
			return refuse()
		}
	}
	f[6] = true
	return refuse()
}
