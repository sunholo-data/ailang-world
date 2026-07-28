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
}

// Session is one serial capability ledger and effect-record stream.
type Session struct {
	mu       sync.Mutex
	store    objectStore
	grants   []Capability
	registry Registry
	mode     Mode
	replay   []hashref.HashRef
	next     int
}

// NewSession constructs a live session over an injected store handle.
func NewSession(s *store.Store, grants []Capability, registry Registry) *Session {
	return newSession(s, grants, registry, Live, nil)
}

// NewReplaySession constructs a replay session consuming recordRefs in order.
func NewReplaySession(
	s *store.Store,
	grants []Capability,
	registry Registry,
	recordRefs []hashref.HashRef,
) *Session {
	return newSession(s, grants, registry, Replay, recordRefs)
}

func newSession(
	s objectStore,
	grants []Capability,
	registry Registry,
	mode Mode,
	recordRefs []hashref.HashRef,
) *Session {
	return &Session{
		store: s, grants: append([]Capability(nil), grants...), registry: registry,
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
	s.mu.Lock()
	defer s.mu.Unlock()

	grantIndex, decision := s.decide(req)
	if s.mode == Replay {
		return s.invokeReplay(req, decision)
	}
	requestRef := requestHash(req, payload)
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

	s.grants[grantIndex].Budget = decision.Remaining
	handler, ok := s.registry[req.Effect]
	if !ok {
		return nil, hashref.HashRef{}, fmt.Errorf("broker: no handler registered for %q", req.Effect)
	}
	result, err = handler.Execute(ctx, req, payload)
	if err != nil {
		rec := EffectRecord{
			Effect: req.Effect, Scope: req.Scope, Cost: req.Cost,
			BudgetBefore: budgetBefore, BudgetAfter: decision.Remaining,
			Allowed: true, Failed: true, RequestRef: requestRef,
		}
		ref, putErr := s.putRecord(rec)
		if putErr != nil {
			return nil, hashref.HashRef{}, putErr
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
	return result, recordRef, nil
}

func (s *Session) decide(req EffectRequest) (int, Decision) {
	if len(s.grants) == 0 {
		return -1, Decision{Label: LabelDeniedEffectName}
	}
	bestIndex := 0
	best := Decide(s.grants[0], req)
	for i := 1; i < len(s.grants); i++ {
		next := Decide(s.grants[i], req)
		if next.Allowed {
			return i, next
		}
		if denialRank(next.Label) > denialRank(best.Label) {
			bestIndex, best = i, next
		}
	}
	return bestIndex, best
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
	text := fmt.Sprintf("%d:%s%d:%s%d:%d:", len(req.Effect), req.Effect, len(req.Scope), req.Scope, req.Cost, req.Now)
	input := append([]byte(text), payload...)
	return hashref.SumSHA256(input)
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
			s.grants[grantIndex].Budget = decision.Remaining
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
		s.grants[grantIndex].Budget = decision.Remaining
		s.next++
		return append([]byte(nil), resultObj.Payload...), recordRef, nil
	}
	s.next++
	return nil, recordRef, &DenialError{Decision: decision, RecordRef: recordRef}
}
