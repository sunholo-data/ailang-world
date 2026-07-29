package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	EffectHumanApprove      = "Human.Approve"
	EffectHumanPollApproval = "Human.PollApproval"

	ApprovalRequestV1  = "world/approval-request/v1"
	ApprovalDecisionV1 = "world/approval-decision/v1"
	ApprovalsV1        = "world/approvals/v1"
)

var (
	ErrApprovalRequestNotFound = errors.New("broker: approval request not found")
	ErrInvalidApprovalDecision = errors.New("broker: invalid approval decision")
)

type approvalStore interface {
	objectStore
	SetRegistryHead(string, hashref.HashRef) error
	GetRegistryHead(string) (hashref.HashRef, bool, error)
}

// HumanHandler implements the synchronous approval request and poll handlers.
// Authorization, debit, result persistence, and effect-record persistence stay
// in Session.Invoke, exactly as for every other handler.
type HumanHandler struct {
	store approvalStore
}

func NewHumanHandler(s *store.Store) *HumanHandler {
	return newHumanHandler(s)
}

func newHumanHandler(s approvalStore) *HumanHandler {
	return &HumanHandler{store: s}
}

type approvalRequestWire struct {
	Effect    string `json:"effect"`
	Scope     string `json:"scope"`
	Cost      int64  `json:"cost"`
	Requester string `json:"requester"`
	Now       int64  `json:"now"`
}

type approvalDecisionWire struct {
	RequestRef string `json:"requestRef"`
	Decision   string `json:"decision"`
	DecidedBy  string `json:"decidedBy"`
	Now        int64  `json:"now"`
}

type approvalHeadWire struct {
	PreviousHead string `json:"previousHead"`
	RequestRef   string `json:"requestRef"`
	DecisionRef  string `json:"decisionRef"`
}

type approvalInputWire struct {
	RequestRef string `json:"requestRef"`
	Requester  string `json:"requester"`
}

type pendingWire struct {
	Status     string `json:"status"`
	RequestRef string `json:"requestRef"`
}

type observedDecisionWire struct {
	Status   string          `json:"status"`
	Decision json.RawMessage `json:"decision"`
}

func (h *HumanHandler) Execute(_ context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	switch req.Effect {
	case EffectHumanApprove:
		var input approvalInputWire
		if len(payload) != 0 {
			if err := decodeApprovalJSON(payload, &input); err != nil {
				return nil, err
			}
		}
		requestPayload := mustApprovalJSON(approvalRequestWire{
			Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
			Requester: input.Requester, Now: req.Now,
		})
		requestObj := brokerObject(ApprovalRequestV1, requestPayload)
		if err := h.store.PutObject(requestObj); err != nil {
			return nil, fmt.Errorf("broker: put approval request: %w", err)
		}
		if err := appendApprovalHead(h.store, requestObj.Hash, hashref.HashRef{}); err != nil {
			return nil, err
		}
		return pendingBytes(requestObj.Hash), nil

	case EffectHumanPollApproval:
		var input approvalInputWire
		if err := decodeApprovalJSON(payload, &input); err != nil {
			return nil, err
		}
		requestRef, err := hashref.Parse(input.RequestRef)
		if err != nil {
			return nil, fmt.Errorf("broker: poll approval requestRef: %w", err)
		}
		decision, found, err := findApprovalDecision(h.store, requestRef)
		if err != nil {
			return nil, err
		}
		if !found {
			return pendingBytes(requestRef), nil
		}
		return mustApprovalJSON(observedDecisionWire{
			Status: "decided", Decision: json.RawMessage(decision.Payload),
		}), nil
	default:
		return nil, fmt.Errorf("broker: Human handler does not implement %q", req.Effect)
	}
}

// DecideApproval is an operator entry point, not an effect. It creates one
// immutable decision object and moves only the approvals registry head.
func DecideApproval(
	s *store.Store,
	requestRef hashref.HashRef,
	decision string,
	decidedBy string,
	now int64,
) (hashref.HashRef, error) {
	return decideApproval(s, requestRef, decision, decidedBy, now)
}

func decideApproval(
	s approvalStore,
	requestRef hashref.HashRef,
	decision string,
	decidedBy string,
	now int64,
) (hashref.HashRef, error) {
	if decision != "approve" && decision != "deny" {
		return hashref.HashRef{}, ErrInvalidApprovalDecision
	}
	request, ok, err := s.GetObject(requestRef)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: read approval request: %w", err)
	}
	if !ok || request.SemanticID != ApprovalRequestV1 {
		return hashref.HashRef{}, ErrApprovalRequestNotFound
	}
	if _, found, err := findApprovalRequest(s, requestRef); err != nil {
		return hashref.HashRef{}, err
	} else if !found {
		return hashref.HashRef{}, ErrApprovalRequestNotFound
	}
	payload := mustApprovalJSON(approvalDecisionWire{
		RequestRef: requestRef.String(), Decision: decision, DecidedBy: decidedBy, Now: now,
	})
	obj := brokerObject(ApprovalDecisionV1, payload)
	if err := s.PutObject(obj); err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: put approval decision: %w", err)
	}
	if err := appendApprovalHead(s, requestRef, obj.Hash); err != nil {
		return hashref.HashRef{}, err
	}
	return obj.Hash, nil
}

func appendApprovalHead(s approvalStore, requestRef, decisionRef hashref.HashRef) error {
	previous, ok, err := s.GetRegistryHead(ApprovalsV1)
	if err != nil {
		return fmt.Errorf("broker: read approvals head: %w", err)
	}
	wire := approvalHeadWire{RequestRef: requestRef.String()}
	if ok {
		wire.PreviousHead = previous.String()
	}
	if !decisionRef.IsZero() {
		wire.DecisionRef = decisionRef.String()
	}
	obj := brokerObject(ApprovalsV1, mustApprovalJSON(wire))
	if err := s.PutObject(obj); err != nil {
		return fmt.Errorf("broker: put approvals head: %w", err)
	}
	if err := s.SetRegistryHead(ApprovalsV1, obj.Hash); err != nil {
		return fmt.Errorf("broker: move approvals head: %w", err)
	}
	return nil
}

func findApprovalRequest(s approvalStore, requestRef hashref.HashRef) (store.Object, bool, error) {
	return walkApprovalHead(s, requestRef, false)
}

func findApprovalDecision(s approvalStore, requestRef hashref.HashRef) (store.Object, bool, error) {
	return walkApprovalHead(s, requestRef, true)
}

func walkApprovalHead(s approvalStore, requestRef hashref.HashRef, wantDecision bool) (store.Object, bool, error) {
	head, ok, err := s.GetRegistryHead(ApprovalsV1)
	if err != nil || !ok {
		return store.Object{}, false, err
	}
	// This is O(all approval-head objects) per poll. A cycle cannot exist in
	// this content-addressed chain: creating one would require an object's hash
	// to contain itself. A future indexed approval surface should replace this
	// linear walk rather than weakening that immutable-chain invariant.
	for !head.IsZero() {
		obj, found, getErr := s.GetObject(head)
		if getErr != nil {
			return store.Object{}, false, getErr
		}
		if !found || obj.SemanticID != ApprovalsV1 {
			return store.Object{}, false, errors.New("broker: approvals head is missing or invalid")
		}
		var wire approvalHeadWire
		if err := decodeApprovalJSON(obj.Payload, &wire); err != nil {
			return store.Object{}, false, err
		}
		if wire.RequestRef == requestRef.String() {
			if wantDecision && wire.DecisionRef != "" {
				ref, err := hashref.Parse(wire.DecisionRef)
				if err != nil {
					return store.Object{}, false, err
				}
				decision, found, err := s.GetObject(ref)
				if err != nil || !found {
					return store.Object{}, false, err
				}
				if decision.SemanticID != ApprovalDecisionV1 {
					return store.Object{}, false, errors.New("broker: approval decision has wrong semantic ID")
				}
				return decision, true, nil
			}
			if !wantDecision {
				request, found, err := s.GetObject(requestRef)
				return request, found, err
			}
		}
		if wire.PreviousHead == "" {
			break
		}
		head, err = hashref.Parse(wire.PreviousHead)
		if err != nil {
			return store.Object{}, false, err
		}
	}
	return store.Object{}, false, nil
}

func pendingBytes(requestRef hashref.HashRef) []byte {
	return mustApprovalJSON(pendingWire{Status: "pending", RequestRef: requestRef.String()})
}

func mustApprovalJSON(value any) []byte {
	payload, err := json.Marshal(value)
	if err != nil {
		panic("broker: fixed approval value cannot fail JSON encoding: " + err.Error())
	}
	return payload
}

func decodeApprovalJSON(payload []byte, value any) error {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	if err := dec.Decode(value); err != nil {
		return fmt.Errorf("broker: decode approval payload: %w", err)
	}
	return nil
}
