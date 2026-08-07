package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	// The five publish-approval refusal classes (SM.B2b / AC9). They are
	// separate sentinels rather than one because a caller that cannot tell
	// "there is no such approval" from "the approval you named says DENY"
	// cannot write the runbook step that follows.
	//
	// ErrApprovalAlreadyConsumed is deliberately NOT among them: single use is
	// enforced durably by store.AppendClaimedEffectIntent, not here, and
	// re-deriving it in memory is exactly the mistake AC8 exists to refuse.
	ErrPublishApprovalMissing   = errors.New("broker: publish approval decision is missing")
	ErrPublishApprovalDenied    = errors.New("broker: publish approval decision is not an approval")
	ErrPublishApprovalMalformed = errors.New("broker: publish approval is malformed")
	ErrPublishApprovalScope     = errors.New("broker: publish approval does not stamp this packet")
	ErrPublishApprovalExpired   = errors.New("broker: publish approval has expired")
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

// ---------------------------------------------------------------------------
// SM.B2b — the publish-bound approval scope, and its validation
//
// Everything below rides the LANDED approval surface. Human.Approve,
// Human.PollApproval, HumanHandler, DecideApproval, ApprovalRequestV1,
// ApprovalDecisionV1, approvalRequestWire and approvalDecisionWire are all
// unchanged: the publish binding is carried in the EXISTING
// approvalRequestWire.Scope STRING. There is no new wire type, no parallel
// approval codec, and no widening of the EffectRecord codec.
// ---------------------------------------------------------------------------

const (
	// publishApprovalScopeMark separates the frozen PublishScope grammar from
	// the publish-bound terms.
	//
	// Two independent things make '#' unambiguous, and the second does not
	// depend on the first:
	//
	//  1. validatePublishOrigin REFUSES any origin carrying a fragment
	//     (registry_publish.go: "carries a query or fragment"), so an accepted
	//     registry origin can never contain one — verified at HEAD, and
	//     exercised by TestLoopbackConstructorRefusesEveryNonLoopbackOrigin.
	//  2. Even if a payload smuggled one in, the parse splits at the FIRST
	//     mark, so PublishApprovalScope.Publish never contains a '#' — while
	//     the wantPublish it is compared against would. The comparison
	//     therefore fails CLOSED rather than colliding. This is what
	//     TestPublishApprovalScopeRefusesASmuggledFragment measures.
	//
	// The version stays part of the head — an approval for 0.1.0 still cannot
	// authorize 0.1.1.
	publishApprovalScopeMark = "#"
	publishApprovalTermSep   = "&"
	publishApprovalKeySep    = "="
)

// publishApprovalScopeTerms is the FROZEN term order of the publish-bound
// suffix. It is stated once, here, and both the formatter and the parser are
// driven from it, so a reordering is a single edit that changes every minted
// scope — and therefore invalidates every approval minted under the old order.
var publishApprovalScopeTerms = []string{
	"effect", "manifest", "tarball", "content", "interface", "expires",
}

// PublishApprovalScope is the parsed canonical scope of an attended
// Registry.Publish approval:
//
//	registry:<origin>/package:<vendor>/<name>/version:<version>#effect=Registry.Publish&manifest=<ref>&tarball=<sha>&content=<sha>&interface=<sha>&expires=<logical-time>
//
// The suffix is what makes the stamp bind BYTES rather than a name: an
// approval whose tarball digest differs from the digest recomputed at dispatch
// authorizes nothing, even for the identical package and version.
//
// Effect is carried HERE rather than read out of approvalRequestWire.Effect
// because the landed HumanHandler stamps that field with the effect that
// MINTED the request — always "Human.Approve" — never with the effect being
// approved. Measured, not assumed: minting through the landed
// Human.Approve/DecideApproval path yields
// approvalRequestWire{Effect:"Human.Approve", ...}, so a gate that demanded
// Effect == Registry.Publish there would refuse the ONLY path an operator has.
// The design's requirement that the broker "check the request's effect" is
// therefore discharged against this content-bound term, which the requester
// commits to at request time and cannot change afterwards.
type PublishApprovalScope struct {
	// Publish is the frozen PublishScope grammar, verbatim.
	Publish string
	// Effect is the effect this stamp authorizes. Only EffectRegistryPublish
	// is accepted by validatePublishApproval.
	Effect        string
	ManifestRef   string
	TarballSHA256 string
	ContentHash   string
	InterfaceHash string
	// ExpiresAt is the LAST logical time at which the approval may be used.
	ExpiresAt int64
}

// FormatPublishApprovalScope renders the canonical scope. It is the only
// producer; PublishApprovalScopeFor is the convenience wrapper handlers use.
func FormatPublishApprovalScope(scope PublishApprovalScope) string {
	values := []string{
		scope.Effect, scope.ManifestRef, scope.TarballSHA256, scope.ContentHash,
		scope.InterfaceHash, strconv.FormatInt(scope.ExpiresAt, 10),
	}
	if len(values) != len(publishApprovalScopeTerms) {
		// A term added to the frozen list without a value here would otherwise
		// index out of range at the first mint. Fail loudly and immediately.
		panic("broker: publish approval scope term/value arity drift")
	}
	terms := make([]string, 0, len(publishApprovalScopeTerms))
	for i, key := range publishApprovalScopeTerms {
		terms = append(terms, key+publishApprovalKeySep+values[i])
	}
	return scope.Publish + publishApprovalScopeMark + strings.Join(terms, publishApprovalTermSep)
}

// PublishApprovalScopeFor builds the canonical scope an attended approval must
// carry to authorize publishing exactly these bytes at exactly this version.
func PublishApprovalScopeFor(id PublishIdentity, hashes PublishHashes, expiresAt int64) string {
	return FormatPublishApprovalScope(PublishApprovalScope{
		Publish:       PublishScope(id.RegistryOrigin, id.Vendor, id.Name, id.Version),
		Effect:        EffectRegistryPublish,
		ManifestRef:   id.ManifestRef.String(),
		TarballSHA256: hashes.TarballSHA256,
		ContentHash:   hashes.ContentHash,
		InterfaceHash: hashes.InterfaceHash,
		ExpiresAt:     expiresAt,
	})
}

// ParsePublishApprovalScope is strict in every direction: exactly one mark,
// exactly len(publishApprovalScopeTerms) terms, in exactly that order, each
// with exactly one key separator, and a decimal expiry. A scope this function
// cannot fully account for must not authorize an irreversible write, so every
// deviation is an error rather than a tolerated extra.
func ParsePublishApprovalScope(raw string) (PublishApprovalScope, error) {
	mark := strings.Index(raw, publishApprovalScopeMark)
	if mark < 0 {
		return PublishApprovalScope{}, fmt.Errorf(
			"%w: scope %q carries no %q mark", ErrPublishApprovalMalformed, raw, publishApprovalScopeMark)
	}
	head, tail := raw[:mark], raw[mark+len(publishApprovalScopeMark):]
	if strings.Contains(tail, publishApprovalScopeMark) {
		return PublishApprovalScope{}, fmt.Errorf(
			"%w: scope %q carries more than one %q mark", ErrPublishApprovalMalformed, raw, publishApprovalScopeMark)
	}
	terms := strings.Split(tail, publishApprovalTermSep)
	if len(terms) != len(publishApprovalScopeTerms) {
		return PublishApprovalScope{}, fmt.Errorf(
			"%w: scope %q carries %d terms, want exactly %d (%v)",
			ErrPublishApprovalMalformed, raw, len(terms), len(publishApprovalScopeTerms),
			publishApprovalScopeTerms)
	}
	values := make([]string, len(terms))
	for i, term := range terms {
		key, value, found := strings.Cut(term, publishApprovalKeySep)
		if !found || strings.Contains(value, publishApprovalKeySep) {
			return PublishApprovalScope{}, fmt.Errorf(
				"%w: scope term %q is not exactly one %q assignment",
				ErrPublishApprovalMalformed, term, publishApprovalKeySep)
		}
		if key != publishApprovalScopeTerms[i] {
			return PublishApprovalScope{}, fmt.Errorf(
				"%w: scope term %d is %q, want %q (the term order is frozen)",
				ErrPublishApprovalMalformed, i, key, publishApprovalScopeTerms[i])
		}
		if value == "" {
			return PublishApprovalScope{}, fmt.Errorf(
				"%w: scope term %q is empty", ErrPublishApprovalMalformed, key)
		}
		values[i] = value
	}
	expiry := values[len(values)-1]
	expiresAt, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return PublishApprovalScope{}, fmt.Errorf(
			"%w: scope expiry %q is not a decimal logical time", ErrPublishApprovalMalformed, expiry)
	}
	if head == "" {
		return PublishApprovalScope{}, fmt.Errorf(
			"%w: scope %q has an empty publish grammar", ErrPublishApprovalMalformed, raw)
	}
	// POSITIONAL literal, deliberately: adding a field to PublishApprovalScope
	// without deciding where it comes from fails to compile rather than
	// silently minting a scope with a zero term.
	return PublishApprovalScope{
		head, values[0], values[1], values[2], values[3], values[4], expiresAt,
	}, nil
}

// approvalObjectReader is the read surface publish validation needs. It is
// exactly objectStore's GetObject, narrowed, so Session can pass its own store
// without widening objectStore for every implementer.
type approvalObjectReader interface {
	GetObject(hashref.HashRef) (store.Object, bool, error)
}

// validatePublishApproval is the SM.B2b gate: it decides whether a LANDED
// attended approval authorizes THIS publish request, and returns the approval
// reference the durable single-use claim is keyed on.
//
// It is called by Session.Invoke BEFORE store.AppendClaimedEffectIntent, which
// is BEFORE the handler runs — so every refusal below is strictly earlier than
// the credential load and the POST, and the dispatch counter is still zero.
//
// The traversal is content-addressed the whole way: payload.approvalRef names
// an ApprovalDecisionV1 object, whose RequestRef names an ApprovalRequestV1
// object, whose canonical Scope names the exact bytes. Nothing here trusts a
// configuration field, and nothing here is re-derivable from memory.
func validatePublishApproval(
	s approvalObjectReader,
	payload []byte,
	req EffectRequest,
) (hashref.HashRef, error) {
	id, hashes, err := DecodePublishPayload(payload)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("%w: %s", ErrPublishApprovalMalformed, err.Error())
	}

	decisionObj, ok, err := s.GetObject(id.ApprovalRef)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: read publish approval decision: %w", err)
	}
	if !ok {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: %s names no object", ErrPublishApprovalMissing, id.ApprovalRef.String())
	}
	if decisionObj.SemanticID != ApprovalDecisionV1 {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: %s is a %q object, want %q",
			ErrPublishApprovalMalformed, id.ApprovalRef.String(), decisionObj.SemanticID, ApprovalDecisionV1)
	}
	var decision approvalDecisionWire
	if err := decodeApprovalJSON(decisionObj.Payload, &decision); err != nil {
		return hashref.HashRef{}, fmt.Errorf("%w: decision: %s", ErrPublishApprovalMalformed, err.Error())
	}
	if decision.Decision != "approve" {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: decision %s says %q", ErrPublishApprovalDenied, id.ApprovalRef.String(), decision.Decision)
	}

	requestRef, err := hashref.Parse(decision.RequestRef)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: decision requestRef: %s", ErrPublishApprovalMalformed, err.Error())
	}
	requestObj, ok, err := s.GetObject(requestRef)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: read publish approval request: %w", err)
	}
	if !ok {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: decision %s references request %s",
			ErrApprovalRequestNotFound, id.ApprovalRef.String(), requestRef.String())
	}
	if requestObj.SemanticID != ApprovalRequestV1 {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: request %s is a %q object, want %q",
			ErrPublishApprovalMalformed, requestRef.String(), requestObj.SemanticID, ApprovalRequestV1)
	}
	var request approvalRequestWire
	if err := decodeApprovalJSON(requestObj.Payload, &request); err != nil {
		return hashref.HashRef{}, fmt.Errorf("%w: request: %s", ErrPublishApprovalMalformed, err.Error())
	}

	// The approval must have been minted by the landed attended request effect,
	// and it must be priced as exactly one publish unit. request.Effect names
	// the effect that MINTED the request (always Human.Approve — see
	// PublishApprovalScope.Effect); the effect being AUTHORIZED is a
	// content-bound term of the scope and is checked immediately below.
	if request.Effect != EffectHumanApprove {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: request %s was minted by effect %q, want %q",
			ErrPublishApprovalScope, requestRef.String(), request.Effect, EffectHumanApprove)
	}
	if request.Cost != PublishCost {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: request %s approves cost %d, want exactly %d",
			ErrPublishApprovalScope, requestRef.String(), request.Cost, PublishCost)
	}

	scope, err := ParsePublishApprovalScope(request.Scope)
	if err != nil {
		return hashref.HashRef{}, err
	}
	if scope.Effect != EffectRegistryPublish {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: approval stamps effect %q, want %q",
			ErrPublishApprovalScope, scope.Effect, EffectRegistryPublish)
	}
	wantPublish := PublishScope(id.RegistryOrigin, id.Vendor, id.Name, id.Version)
	// Three separate identities, each named on failure: what the approval
	// stamps, what the payload publishes, and what the capability was decided
	// against. Collapsing them would let a grant for one package ride an
	// approval for another.
	if scope.Publish != wantPublish {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: approval stamps %q but the payload publishes %q",
			ErrPublishApprovalScope, scope.Publish, wantPublish)
	}
	if req.Scope != wantPublish {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: effect scope %q does not describe the payload package %q",
			ErrPublishApprovalScope, req.Scope, wantPublish)
	}
	for _, arm := range []struct{ name, stamped, packet string }{
		{"manifest ref", scope.ManifestRef, id.ManifestRef.String()},
		{"tarball hash", scope.TarballSHA256, hashes.TarballSHA256},
		{"content hash", scope.ContentHash, hashes.ContentHash},
		{"interface hash", scope.InterfaceHash, hashes.InterfaceHash},
	} {
		if arm.stamped != arm.packet {
			return hashref.HashRef{}, fmt.Errorf(
				"%w: approval stamps %s %s but the payload carries %s",
				ErrPublishApprovalScope, arm.name, arm.stamped, arm.packet)
		}
	}

	// Logical time, in both directions. An approval requested or decided AFTER
	// the publish it authorizes is not evidence of anything, and an approval
	// used after its stamped expiry is spent authority.
	if request.Now > req.Now {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: request %s was made at logical time %d, after the publish at %d",
			ErrPublishApprovalMalformed, requestRef.String(), request.Now, req.Now)
	}
	if decision.Now < request.Now || decision.Now > req.Now {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: decision %s was made at logical time %d, outside [request %d, publish %d]",
			ErrPublishApprovalMalformed, id.ApprovalRef.String(), decision.Now, request.Now, req.Now)
	}
	if scope.ExpiresAt < request.Now {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: approval expiry %d precedes its own request time %d",
			ErrPublishApprovalMalformed, scope.ExpiresAt, request.Now)
	}
	if req.Now > scope.ExpiresAt {
		return hashref.HashRef{}, fmt.Errorf(
			"%w: approval %s expired at logical time %d; the publish is at %d",
			ErrPublishApprovalExpired, id.ApprovalRef.String(), scope.ExpiresAt, req.Now)
	}
	return id.ApprovalRef, nil
}
