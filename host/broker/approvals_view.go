package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ApprovalReader is the read-only surface needed by the workbench projection.
type ApprovalReader interface {
	GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error)
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
}

// ApprovalHead reuses the producer's wire shape; heads are newest first.
type ApprovalHead = approvalHeadWire

type ApprovalSummary struct {
	RequestRef, DecisionRef  string
	Effect, Scope, Requester string
	Cost                     int64
	Status, DecidedBy        string
}
type ApprovalPage struct {
	Summaries    []ApprovalSummary
	Truncated    bool
	MalformedRef string
}

var ErrApprovalChainMalformed = errors.New("broker: approvals chain is malformed")

// FoldApprovals preserves newest relevant encounter order, retaining the
// newest decision reference for each request without reading any bodies.
func FoldApprovals(heads []ApprovalHead) []ApprovalSummary {
	var rows []ApprovalSummary
	indexes := map[string]int{}
	for _, h := range heads {
		i, ok := indexes[h.RequestRef]
		if !ok {
			i = len(rows)
			indexes[h.RequestRef] = i
			rows = append(rows, ApprovalSummary{RequestRef: h.RequestRef, Status: "pending"})
		}
		if rows[i].DecisionRef == "" && h.DecisionRef != "" {
			rows[i].DecisionRef = h.DecisionRef
		}
	}
	return rows
}

// RecentApprovals reads at most 40 heads, folds them, and reads bodies only
// for the first 20 retained summaries (at most 80 object reads in total).
func RecentApprovals(ctx context.Context, r ApprovalReader, maxHeads int) (ApprovalPage, error) {
	p := ApprovalPage{}
	if maxHeads < 1 {
		return p, fmt.Errorf("broker: maxHeads must be positive")
	}
	maxHeads = min(maxHeads, 40)
	malformed := func(ref string) (ApprovalPage, error) { p.MalformedRef = ref; return p, ErrApprovalChainMalformed }
	head, ok, err := r.GetRegistryHead(ctx, ApprovalsV1)
	if err != nil || !ok {
		return p, err
	}
	seen := map[hashref.HashRef]bool{}
	var heads []ApprovalHead
	for !head.IsZero() && len(heads) < maxHeads {
		if seen[head] {
			return malformed(head.String())
		}
		seen[head] = true
		o, found, err := r.GetObject(ctx, head)
		if err != nil {
			return p, err
		}
		var h ApprovalHead
		if !found || o.SemanticID != ApprovalsV1 || decodeApprovalJSON(o.Payload, &h) != nil {
			return malformed(head.String())
		}
		if _, err := hashref.Parse(h.RequestRef); err != nil {
			return malformed(head.String())
		}
		if h.DecisionRef != "" {
			if _, err := hashref.Parse(h.DecisionRef); err != nil {
				return malformed(head.String())
			}
		}
		heads = append(heads, h)
		if h.PreviousHead == "" {
			head = hashref.HashRef{}
			break
		}
		next, err := hashref.Parse(h.PreviousHead)
		if err != nil {
			return malformed(head.String())
		}
		head = next
	}
	p.Truncated = !head.IsZero()
	rows := FoldApprovals(heads)
	if len(rows) > 20 {
		p.Truncated = true
		rows = rows[:20]
	}
	// Head references were validated before folding. Missing or mistyped body
	// objects are malformed data; transport failures retain their I/O identity.
	readBody := func(ref, id string, out any) (bool, error) {
		h, err := hashref.Parse(ref)
		if err != nil {
			return false, nil
		}
		o, found, err := r.GetObject(ctx, h)
		if err != nil {
			return false, err
		}
		return found && o.SemanticID == id && decodeApprovalJSON(o.Payload, out) == nil, nil
	}
	for _, row := range rows {
		var req approvalRequestWire
		valid, err := readBody(row.RequestRef, ApprovalRequestV1, &req)
		if err != nil {
			return p, err
		}
		if !valid {
			return malformed(row.RequestRef)
		}
		row.Effect, row.Scope, row.Requester, row.Cost = req.Effect, req.Scope, req.Requester, req.Cost
		if row.DecisionRef != "" {
			var dec approvalDecisionWire
			valid, err := readBody(row.DecisionRef, ApprovalDecisionV1, &dec)
			if err != nil {
				return p, err
			}
			if !valid || dec.RequestRef != row.RequestRef || (dec.Decision != "approve" && dec.Decision != "deny") {
				return malformed(row.DecisionRef)
			}
			row.Status, row.DecidedBy = dec.Decision, dec.DecidedBy
		}
		p.Summaries = append(p.Summaries, row)
	}
	return p, nil
}
