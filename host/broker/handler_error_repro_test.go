package broker

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// CF-J-2 guards the attended, ratified third-arm law (charter c26b27d): an
// ordinary in-process handler error writes exactly one failure record before
// returning, and its debit stands. This does not close the dispatch→record
// crash window; process death in that interval remains outside M3.B0.
type cfj2FailingHandler struct{ calls int }

func (h *cfj2FailingHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	h.calls++
	return nil, errors.New("handler failed after partial execution")
}

type cfj2RecordingStore struct {
	base    *store.Store
	records []store.Object
}

func (s *cfj2RecordingStore) PutObject(obj store.Object) error {
	if err := s.base.PutObject(obj); err != nil {
		return err
	}
	if obj.SemanticID == EffectRecordV1 {
		s.records = append(s.records, obj)
	}
	return nil
}

func (s *cfj2RecordingStore) GetObject(ref hashref.HashRef) (store.Object, bool, error) {
	return s.base.GetObject(ref)
}

func (s *cfj2RecordingStore) AppendNextEffectIntent(
	episodeID string,
	intent store.EffectIntent,
) (string, int64, error) {
	return s.base.AppendNextEffectIntent(episodeID, intent)
}

func (s *cfj2RecordingStore) AppendClaimedEffectIntent(
	episodeID string,
	intent store.EffectIntent,
	approvalRef, requestRef hashref.HashRef,
) (string, int64, error) {
	return s.base.AppendClaimedEffectIntent(episodeID, intent, approvalRef, requestRef)
}

func (s *cfj2RecordingStore) AppendEffectOutcome(
	id string,
	outcome store.EffectOutcome,
) (int64, hashref.HashRef, error) {
	return s.base.AppendEffectOutcome(id, outcome)
}

func TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord(t *testing.T) {
	st := openTestStore(t)
	recording := &cfj2RecordingStore{base: st}
	handler := &cfj2FailingHandler{}
	s := newSession(recording, "handler-error-repro",
		[]Capability{{Effect: "FS.Write", Scope: "/p", ExpiresAt: 100, Budget: 5}},
		Registry{"FS.Write": handler}, Live, nil)

	before := s.grants[0].Budget
	_, ref, err := s.Invoke(context.Background(),
		EffectRequest{Effect: "FS.Write", Scope: "/p", Cost: 3, Now: 1}, []byte("x"))

	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("Invoke error = %v, want *EffectFailedError", err)
	}
	if handler.calls != 1 {
		t.Fatalf("handler dispatched %d times, want exactly 1", handler.calls)
	}
	if failed.Unwrap() == nil {
		t.Fatal("live failure must unwrap the handler error")
	}
	if ref.IsZero() || failed.RecordRef != ref {
		t.Fatalf("record refs = %s and %s, want same non-zero ref", ref, failed.RecordRef)
	}

	// The debit stands by ratification even if the handler partially executed.
	if got, want := s.grants[0].Budget, before-3; got != want {
		t.Fatalf("ledger budget = %d, want %d", got, want)
	}
	if got := len(recording.records); got != 1 {
		t.Fatalf("failure record count = %d, want exactly 1", got)
	}
	rec, err := DecodeRecord(recording.records[0].Payload)
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Allowed || !rec.Failed || !rec.ResultRef.IsZero() || rec.Denial != "" ||
		rec.BudgetAfter != rec.BudgetBefore-rec.Cost || !RecordConsistent(rec) {
		t.Fatalf("failure record = %#v", rec)
	}
}
