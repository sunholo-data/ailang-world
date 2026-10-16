package daemon

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

// Row 172 decision-packet sentences (D6). Each refusal sentence is the design's frozen wording.
const (
	packetValidatedFormat              = "decision packet validated: committed with its invocation record at entry #%d; proposal, approval request and approvals-chain anchor agree; lifecycle fields are the v1 producer's"
	packetUniqueTruncation             = " (duplicate check bounded to the first 100 stored packets)"
	packetRefusedBytes                 = "not a validated decision packet: the bytes are not a canonical world/decision-packet/v1 encoding whose packetHash recomputes"
	packetRefusedCarried               = "not a validated decision packet: it was not committed together with its invocation record"
	packetRefusedProposal              = "not a validated decision packet: proposalHash is not the plan of the invocation record committed with it"
	packetRefusedRequest               = "not a validated decision packet: requestRef is not the approval request that invocation's Human.Approve effect returned"
	packetRefusedAnchor                = "not a validated decision packet: its approval request is not in the broker's approvals chain"
	packetRefusedAnchorBounded         = "not a validated decision packet: its approval request is not among the newest 40 approval heads, so its approvals-chain membership is not established"
	packetRefusedAnchorMalformedPrefix = "not a validated decision packet: the approvals chain is malformed at " // PD13
	packetRefusedLifecycle             = "not a validated decision packet: its lifecycle fields disagree with its request's logical time or are not the v1 producer's"
	packetRefusedDuplicate             = "not a validated decision packet: another stored packet names the same approval request and the same proposal"
)

// packetRefusalSentence maps a refusal code (1..7) to its sentence.
func packetRefusalSentence(res coordinator.PacketResult) string {
	switch res.Code {
	case 1:
		return packetRefusedBytes
	case 2:
		return packetRefusedCarried
	case 3:
		return packetRefusedProposal
	case 4:
		return packetRefusedRequest
	case 5:
		if res.MalformedRef != "" {
			return packetRefusedAnchorMalformedPrefix + res.MalformedRef
		}
		if res.AnchorTruncated {
			return packetRefusedAnchorBounded
		}
		return packetRefusedAnchor
	case 6:
		return packetRefusedLifecycle
	case 7:
		return packetRefusedDuplicate
	}
	return packetRefusedBytes
}

// packetProjection retains the validated roster internally until M5b adds its
// workbench fields and renderer. Existing subject grades still use row 105.
type packetProjection struct {
	view   *workbench.PacketView
	roster packetRoster
}

// packetView projects a world/decision-packet/v1 object for its page.
func (d *Daemon) packetView(ctx context.Context, obj store.Object) (*workbench.PacketView, error) {
	projection, err := d.projectPacket(ctx, obj)
	if err != nil {
		return nil, err
	}
	return projection.view, nil
}

func (d *Daemon) projectPacket(ctx context.Context, obj store.Object) (*packetProjection, error) {
	res, err := coordinator.ValidatePacket(ctx, d.reads, obj)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return &packetProjection{view: &workbench.PacketView{Sentence: packetRefusalSentence(res)}}, nil
	}
	sentence := fmt.Sprintf(packetValidatedFormat, res.Entry)
	if res.UniqueTruncated {
		sentence += packetUniqueTruncation
	}
	proposalGrade, err := d.objectGradeView(ctx, res.Packet.ProposalHash)
	if err != nil {
		return nil, err
	}
	tfGrade, err := d.objectGradeView(ctx, res.TransitionFn) // MV10 target
	if err != nil {
		return nil, err
	}
	status := res.Approval.Status
	if res.Approval.DecidedBy != "" {
		status += " by " + res.Approval.DecidedBy
	}
	return &packetProjection{roster: packetRosterFromValidated(ctx, d.reads, res), view: &workbench.PacketView{Validated: true, Sentence: sentence, EntryIndex: res.Entry, EntryHref: pageHref(res.Entry, res.Entry),
		ProposalRef: res.Packet.ProposalHash.String(), ProposalHref: "?object=" + res.Packet.ProposalHash.String(),
		RequestRef: res.Packet.RequestRef.String(), RequestHref: "?object=" + res.Packet.RequestRef.String(),
		CreatedAt: res.Packet.CreatedAt, DeadlineAt: res.Packet.DeadlineAt, Policy: res.Packet.Policy,
		Escalations: res.Packet.EscalationsRemaining, ApprovalStatus: status, TransitionFnRef: res.TransitionFn.String(),
		ProposalGrade: proposalGrade, TransitionGrade: tfGrade}}, nil
}
