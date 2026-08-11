package broker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Mode selects physical dispatch or deterministic record replay.
type Mode uint8

const (
	Live Mode = iota
	Replay
)

// Handler executes an already-authorized effect. It contains no authority
// logic.
type Handler interface {
	Execute(context.Context, EffectRequest, []byte) ([]byte, error)
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(context.Context, EffectRequest, []byte) ([]byte, error)

func (f HandlerFunc) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	return f(ctx, req, payload)
}

// Registry maps exact effect names to handlers.
type Registry map[string]Handler

type objectStore interface {
	PutObject(store.Object) error
	GetObject(hashref.HashRef) (store.Object, bool, error)
	AppendNextEffectIntent(string, store.EffectIntent) (string, int64, error)
	AppendClaimedEffectIntent(string, store.EffectIntent, hashref.HashRef, hashref.HashRef) (string, int64, error)
	AppendEffectOutcome(string, store.EffectOutcome) (int64, hashref.HashRef, error)
}

// Session is one serial capability ledger and effect-record stream.
type Session struct {
	mu        sync.Mutex
	store     objectStore
	episodeID string
	grants    []Capability
	registry  Registry
	mode      Mode
	replay    []hashref.HashRef
	next      int
	epoch     int64
}

// CapabilitySnapshot is an immutable copy of a session capability ledger at
// one instant. Now is supplied by the caller; the broker never reads a clock.
type CapabilitySnapshot struct {
	Epoch  int64
	Now    int64
	grants []Capability
}

// CapabilitySnapshot returns a detached view of the current ledger.
func (s *Session) CapabilitySnapshot(now int64) CapabilitySnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return CapabilitySnapshot{Epoch: s.epoch, Now: now, grants: append([]Capability(nil), s.grants...)}
}

// Grants returns a fresh copy of the snapshot's grants.
func (c CapabilitySnapshot) Grants() []Capability {
	return append([]Capability(nil), c.grants...)
}

// Len reports the number of grants in the snapshot.
func (c CapabilitySnapshot) Len() int { return len(c.grants) }

func (s *Session) debitGrant(i int, remaining int64) {
	s.grants[i].Budget = remaining
	s.epoch++
}

// NewSession constructs a live session over an injected store handle.
func NewSession(s *store.Store, episodeID string, grants []Capability, registry Registry) *Session {
	return newSession(s, episodeID, grants, registry, Live, nil)
}

// NewReplaySession constructs a replay session consuming recordRefs in order.
func NewReplaySession(
	s *store.Store,
	grants []Capability,
	registry Registry,
	recordRefs []hashref.HashRef,
) *Session {
	return newSession(s, "", grants, registry, Replay, recordRefs)
}

func newSession(
	s objectStore,
	episodeID string,
	grants []Capability,
	registry Registry,
	mode Mode,
	recordRefs []hashref.HashRef,
) *Session {
	return &Session{
		store: s, episodeID: episodeID,
		grants: append([]Capability(nil), grants...), registry: registry,
		mode: mode, replay: append([]hashref.HashRef(nil), recordRefs...),
	}
}

// DenialError is returned only after the denial record has been persisted.
type DenialError struct {
	Decision  Decision
	RecordRef hashref.HashRef
}

func (e *DenialError) Error() string {
	return fmt.Sprintf("broker: effect denied: %s (record %s)", e.Decision.Label, e.RecordRef.String())
}

// EffectFailedError reports a dispatched handler failure after its immutable
// failure record has been persisted. Live errors unwrap the handler's detail;
// replay errors are reconstructed from the record alone and do not.
type EffectFailedError struct {
	Effect    string
	Scope     string
	RecordRef hashref.HashRef
	cause     error
}

func (e *EffectFailedError) Error() string {
	return fmt.Sprintf("broker: effect %q failed for scope %q (record %s)", e.Effect, e.Scope, e.RecordRef.String())
}

func (e *EffectFailedError) Unwrap() error {
	return e.cause
}

// ReplayGapError reports a missing or inconsistent next replay record.
type ReplayGapError struct {
	Index int
	Why   string
}

func (e *ReplayGapError) Error() string {
	return fmt.Sprintf("broker: replay gap at record %d: %s", e.Index, e.Why)
}

// Invoke runs the one frozen decision/debit/dispatch/record pipeline.
func (s *Session) Invoke(
	ctx context.Context,
	req EffectRequest,
	payload []byte,
) (result []byte, recordRef hashref.HashRef, err error) {
	return s.invoke(ctx, req, payload)
}

func (s *Session) invoke(
	ctx context.Context,
	req EffectRequest,
	payload []byte,
) (result []byte, recordRef hashref.HashRef, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	grantIndex, decision := s.decide(req)
	if s.mode == Replay {
		return s.invokeReplay(req, decision)
	}
	requestObj := requestObject(req, payload)
	requestRef := requestObj.Hash
	budgetBefore := int64(0)
	if grantIndex >= 0 {
		budgetBefore = s.grants[grantIndex].Budget
	}
	if !decision.Allowed {
		rec := EffectRecord{
			Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
			BudgetBefore: budgetBefore, BudgetAfter: budgetBefore,
			Denial: decision.Label, RequestRef: requestRef,
		}
		ref, putErr := s.putRecord(rec)
		if putErr != nil {
			return nil, hashref.HashRef{}, putErr
		}
		return nil, ref, &DenialError{Decision: decision, RecordRef: ref}
	}

	if s.episodeID == "" {
		// Keep this broker-boundary guard for an actionable error. The store
		// independently rejects the same input in
		// TestAppendNextEffectIntentValidationAndOrdinalDerivation, while
		// TestAllowedLiveSessionRequiresEpisodeID pins this boundary message.
		return nil, hashref.HashRef{}, errors.New("broker: live allowed effect requires an episode ID")
	}
	handler, ok := s.registry[req.Effect]
	if !ok {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: no handler registered for %q", req.Effect)
	}
	if err := s.store.PutObject(requestObj); err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: put effect request: %w", err)
	}
	intent := store.EffectIntent{
		EpisodeID: s.episodeID, Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
		RequestRef: requestRef, LogicalTime: req.Now,
	}
	var (
		effectID string
		ordinal  int64
	)
	if req.Effect == EffectRegistryPublish {
		// Registry.Publish is the one irreversible effect, so its durable
		// intent is written by the transaction that ALSO consumes the attended
		// approval (SM.B1's AppendClaimedEffectIntent). A separate claim and
		// intent would let a restart or a second session re-dispatch the same
		// approval; here neither the claim, the journal row nor the intent
		// object becomes visible without the other.
		//
		// SM.B2b: the attended stamp is VALIDATED against the landed
		// ApprovalDecisionV1/ApprovalRequestV1 objects first, here rather than
		// in the handler, for one reason — "already consumed" is one of AC9's
		// seven refusal classes and it is decided by the claim transaction on
		// the very next line. Putting the other six on the far side of that
		// transaction would split one refusal set across two layers, and would
		// require handing the handler a store. Both refusal families therefore
		// land BEFORE the handler runs, i.e. before the credential is read and
		// before any request can leave this process.
		approvalRef, approvalErr := validatePublishApproval(s.store, payload, req)
		if approvalErr != nil {
			return nil, hashref.HashRef{}, approvalErr
		}
		effectID, ordinal, err = s.store.AppendClaimedEffectIntent(
			s.episodeID, intent, approvalRef, requestRef)
	} else {
		effectID, ordinal, err = s.store.AppendNextEffectIntent(s.episodeID, intent)
	}
	if err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: append effect intent: %w", err)
	}
	s.debitGrant(grantIndex, decision.Remaining)
	result, err = handler.Execute(ctx, req, payload)
	if err != nil {
		// THE ONE NARROW SPECIAL CASE. Only the typed *IndeterminateEffectError
		// suppresses the outcome, and only after the intent is already durable.
		// Any other error — including a subprocess timeout, an overflow or a
		// non-zero exit — keeps the landed resolved-failed behaviour verbatim
		// in the block below. MUT-SM-ALL-ERRORS-PENDING widens this arm and
		// reds the landed handler-failure tests, which is what proves it narrow.
		var indeterminate *IndeterminateEffectError
		if errors.As(err, &indeterminate) {
			indeterminate.InvocationID = effectID
			indeterminate.EpisodeID = s.episodeID
			indeterminate.Ordinal = ordinal
			if indeterminate.Effect == "" {
				indeterminate.Effect = req.Effect
			}
			if indeterminate.Scope == "" {
				indeterminate.Scope = req.Scope
			}
			// No EffectRecord, no outcome. The budget stays debited: an
			// ambiguous dispatch consumed the one attended attempt, and only
			// SM.C's read-only reconciliation may resolve the receipt.
			return nil, hashref.HashRef{}, err
		}
		rec := EffectRecord{
			Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
			BudgetBefore: budgetBefore, BudgetAfter: decision.Remaining,
			Allowed: true, Failed: true, RequestRef: requestRef,
		}
		ref, putErr := s.putRecord(rec)
		if putErr != nil {
			return nil, hashref.HashRef{}, putErr
		}
		if _, _, outcomeErr := s.store.AppendEffectOutcome(effectID, store.EffectOutcome{
			InvocationID: effectID, Status: "failed", RecordRef: ref, LogicalTime: req.Now,
		}); outcomeErr != nil {
			return nil, hashref.HashRef{}, fmt.Errorf("broker: append failed effect outcome: %w", outcomeErr)
		}
		return nil, ref, &EffectFailedError{
			Effect: req.Effect, Scope: req.Scope, RecordRef: ref, cause: err,
		}
	}
	resultObj := resultObject(result)
	if err := s.store.PutObject(resultObj); err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: put effect result: %w", err)
	}
	rec := EffectRecord{
		Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
		BudgetBefore: budgetBefore, BudgetAfter: decision.Remaining, Allowed: true,
		RequestRef: requestRef, ResultRef: resultObj.Hash,
	}
	recordRef, err = s.putRecord(rec)
	if err != nil {
		return nil, hashref.HashRef{}, err
	}
	if _, _, err := s.store.AppendEffectOutcome(effectID, store.EffectOutcome{
		InvocationID: effectID, Status: "succeeded", RecordRef: recordRef, LogicalTime: req.Now,
	}); err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: append succeeded effect outcome: %w", err)
	}
	return result, recordRef, nil
}

func (s *Session) decide(req EffectRequest) (int, Decision) {
	return decideOver(s.grants, req)
}

func denialRank(label string) int {
	switch label {
	case LabelDeniedEffectName:
		return 0
	case LabelDeniedScope:
		return 1
	case LabelDeniedExpired:
		return 2
	case LabelDeniedBudget:
		return 3
	default:
		return 4
	}
}

func requestHash(req EffectRequest, payload []byte) hashref.HashRef {
	return hashref.SumSHA256(requestBytes(req, payload))
}

func requestBytes(req EffectRequest, payload []byte) []byte {
	text := fmt.Sprintf("%d:%s%d:%s%d:%d:", len(req.Effect), req.Effect, len(req.Scope), req.Scope, req.Cost, req.Now)
	return append([]byte(text), payload...)
}

func requestObject(req EffectRequest, payload []byte) store.Object {
	return brokerObject(EffectRequestV1, requestBytes(req, payload))
}

func (s *Session) putRecord(rec EffectRecord) (hashref.HashRef, error) {
	if !RecordConsistent(rec) {
		return hashref.HashRef{}, errors.New("broker: refuses inconsistent effect record")
	}
	obj := recordObject(rec)
	if err := s.store.PutObject(obj); err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: put effect record: %w", err)
	}
	return obj.Hash, nil
}

func (s *Session) invokeReplay(
	req EffectRequest,
	decision Decision,
) ([]byte, hashref.HashRef, error) {
	index := s.next
	if index >= len(s.replay) {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "record stream exhausted"}
	}
	recordRef := s.replay[index]
	obj, ok, err := s.store.GetObject(recordRef)
	if err != nil {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "read record: " + err.Error()}
	}
	if !ok {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "record object is missing"}
	}
	if obj.SemanticID != EffectRecordV1 {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "object is not an effect record"}
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: err.Error()}
	}
	grantIndex, _ := s.decide(req)
	expectedBudget := int64(0)
	if grantIndex >= 0 {
		expectedBudget = s.grants[grantIndex].Budget
	}
	if rec.Effect != req.Effect || rec.Scope != req.Scope || rec.Cost != req.Cost ||
		rec.BudgetBefore != expectedBudget ||
		rec.Allowed != decision.Allowed ||
		rec.Allowed && (rec.Denial != "" || rec.BudgetAfter != decision.Remaining) ||
		!rec.Allowed && rec.Denial != decision.Label ||
		!RecordConsistent(rec) {
		return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "record does not match request and decision"}
	}
	if decision.Allowed {
		if rec.Failed {
			s.debitGrant(grantIndex, decision.Remaining)
			s.next++
			return nil, recordRef, &EffectFailedError{
				Effect: rec.Effect, Scope: rec.Scope, RecordRef: recordRef,
			}
		}
		resultObj, ok, err := s.store.GetObject(rec.ResultRef)
		if err != nil {
			return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "read result: " + err.Error()}
		}
		if !ok {
			return nil, hashref.HashRef{}, &ReplayGapError{Index: index, Why: "result object is missing"}
		}
		s.debitGrant(grantIndex, decision.Remaining)
		s.next++
		return append([]byte(nil), resultObj.Payload...), recordRef, nil
	}
	s.next++
	return nil, recordRef, &DenialError{Decision: decision, RecordRef: recordRef}
}
