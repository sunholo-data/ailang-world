package daemon

// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decisions
// GET /surface/requests.json lists and approve --list prints. One read of the approvals chain
// (broker.RecentApprovals: at most 20 summaries of 40 heads) and one packet scan of at most
// decisionListScan packets; each request one scanned packet carries gets ValidatePacket's verdict and
// the law's display reason at now. Display only: the POST runs the complete discovery and the law.

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decisionListScan is the listing's packet scan: the packet validator's 100-packet bound plus one
// sentinel, whose presence marks the list truncated.
const decisionListScan = 101

// PendingDecisions reads the decisions list from the world store s at now (the daemon's handle or
// the CLI's read-only one). Store errors and a malformed approvals chain are returned, never a row.
func PendingDecisions(ctx context.Context, s *store.Store, now int64) ([]surfacepage.PendingDecision, bool, error) {
	page, err := broker.RecentApprovals(ctx, s, 40)
	if err != nil {
		return nil, false, err
	}
	scanned, err := s.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, "", decisionListScan)
	if err != nil {
		return nil, false, err
	}
	scanTruncated := len(scanned) == decisionListScan
	carrying := map[string][]store.Object{}
	for _, o := range scanned[:min(len(scanned), decisionListScan-1)] {
		full, ok, err := s.GetObject(ctx, o.Hash)
		if err != nil {
			return nil, false, err
		}
		if p, err := coordinator.DecodeDecisionPacket(full.Payload); ok && err == nil {
			carrying[p.RequestRef.String()] = append(carrying[p.RequestRef.String()], full)
		}
	}
	var cands []surfacepage.DecisionCandidate
	for _, sm := range page.Summaries {
		c := surfacepage.DecisionCandidate{Request: sm.RequestRef, Effect: sm.Effect, Scope: sm.Scope, Requester: sm.Requester,
			Cost: sm.Cost, Decided: sm.Status != "pending"}
		packets := carrying[sm.RequestRef]
		for _, p := range packets {
			c.Packets = append(c.Packets, p.Hash.String())
		}
		if !c.Decided && sm.Effect == broker.EffectHumanApprove && len(packets) != 0 {
			f := broker.DecisionFacts{Ref: sm.RequestRef, RefOK: true, RequestOK: true, OnChain: true, Scope: sm.Scope, Now: now,
				Discovery: broker.Discovery{Complete: true, Scanned: len(scanned), Carrying: len(packets)}}
			if len(packets) == 1 {
				res, err := coordinator.ValidatePacket(ctx, s, packets[0])
				if err != nil {
					return nil, false, err
				}
				f.PacketCode, f.PacketSentence, f.Deadline = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt
				c.Deadline = res.Packet.DeadlineAt
			}
			c.Reserved, _ = broker.ReservedApprovalScope(sm.Scope)
			c.Refusal = broker.SubjectRefusal(f).Reason
		}
		cands = append(cands, c)
	}
	rows, truncated := surfacepage.FoldDecisions(cands, page.Truncated, scanTruncated)
	return rows, truncated, nil
}

// readDecisionList fills out's decisions; a read that fails is logged by the caller and lists none.
func (d *Daemon) readDecisionList(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
	rows, truncated, err := PendingDecisions(ctx, d.store, now)
	if err != nil {
		return err
	}
	out.Decisions, out.DecisionsTruncated = rows, truncated
	return nil
}
