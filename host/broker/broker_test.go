package broker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func openTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("close store: %v", err)
		}
	})
	return s
}

func echoRegistry(counter *int) Registry {
	return Registry{"probe": HandlerFunc(func(
		_ context.Context,
		_ EffectRequest,
		payload []byte,
	) ([]byte, error) {
		*counter++
		return append([]byte("echo:"), payload...), nil
	})}
}

func TestDeniedInvokeWritesOneRecord(t *testing.T) {
	s := openTestStore(t)
	count := 0
	session := NewSession(s, "denied-record", []Capability{{"probe", "/ok", 10, 5}}, echoRegistry(&count))
	req := EffectRequest{Effect: "probe", Scope: "/wrong", Cost: 2, Now: 1}
	result, recordRef, err := session.Invoke(context.Background(), req, []byte("x"))
	var denial *DenialError
	if !errors.As(err, &denial) {
		t.Fatalf("Invoke error = %v, want *DenialError", err)
	}
	if result != nil {
		t.Fatalf("denied result = %q, want nil", result)
	}
	if recordRef.IsZero() || denial.RecordRef != recordRef {
		t.Fatalf("record refs = %s and %s, want same non-zero ref", recordRef, denial.RecordRef)
	}
	if count != 0 {
		t.Fatalf("handler dispatch count = %d, want 0", count)
	}
	obj, ok, err := s.GetObject(recordRef)
	if err != nil || !ok {
		t.Fatalf("GetObject = ok %v, err %v", ok, err)
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Allowed || rec.Denial != LabelDeniedScope || rec.BudgetAfter != rec.BudgetBefore ||
		!rec.ResultRef.IsZero() || !RecordConsistent(rec) {
		t.Fatalf("denial record = %#v", rec)
	}
	receipt, hasIntent, err := s.GetEffectReceipt(store.EffectInvocationID("denied-record", 0))
	if err != nil || hasIntent || receipt.State != store.ReceiptNotStarted {
		t.Fatalf("denied receipt = %#v, hasIntent %v, err %v; want not-started", receipt, hasIntent, err)
	}
}

func TestAllowedInvokeWritesResultAndRecord(t *testing.T) {
	s := openTestStore(t)
	count := 0
	session := NewSession(s, "allowed-record", []Capability{{"probe", "/ok", 10, 5}}, echoRegistry(&count))
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 2, Now: 1}
	result, recordRef, err := session.Invoke(context.Background(), req, []byte("x"))
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != "echo:x" || count != 1 {
		t.Fatalf("result %q, dispatches %d", result, count)
	}
	obj, ok, err := s.GetObject(recordRef)
	if err != nil || !ok {
		t.Fatalf("GetObject = ok %v, err %v", ok, err)
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Allowed || rec.BudgetBefore != 5 || rec.BudgetAfter != 3 ||
		rec.ResultRef.IsZero() || !RecordConsistent(rec) {
		t.Fatalf("allowed record = %#v", rec)
	}
	resultObj, ok, err := s.GetObject(rec.ResultRef)
	if err != nil || !ok || !bytes.Equal(resultObj.Payload, result) {
		t.Fatalf("result object = ok %v, payload %q, err %v", ok, resultObj.Payload, err)
	}
}

func TestLedgerUsesRemainingBudget(t *testing.T) {
	s := openTestStore(t)
	count := 0
	session := NewSession(s, "ledger-budget", []Capability{{"probe", "/ok", 10, 5}}, echoRegistry(&count))
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 3, Now: 1}
	if _, _, err := session.Invoke(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
	_, ref, err := session.Invoke(context.Background(), req, nil)
	var denial *DenialError
	if !errors.As(err, &denial) || denial.Decision.Label != LabelDeniedBudget {
		t.Fatalf("second Invoke error = %v, want denied:budget", err)
	}
	if count != 1 {
		t.Fatalf("dispatch count = %d, want 1", count)
	}
	obj, ok, getErr := s.GetObject(ref)
	if getErr != nil || !ok {
		t.Fatalf("GetObject = ok %v, err %v", ok, getErr)
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if rec.BudgetBefore != 2 || rec.BudgetAfter != 2 {
		t.Fatalf("denial budgets = %d -> %d, want 2 -> 2", rec.BudgetBefore, rec.BudgetAfter)
	}
}

type failRecordStore struct {
	base               *store.Store
	effectOutcomeCalls int
}

func (s *failRecordStore) PutObject(obj store.Object) error {
	if obj.SemanticID == EffectRecordV1 {
		return errors.New("injected record write failure")
	}
	return s.base.PutObject(obj)
}

func (s *failRecordStore) GetObject(ref hashref.HashRef) (store.Object, bool, error) {
	return s.base.GetObject(ref)
}

func (s *failRecordStore) AppendNextEffectIntent(
	episodeID string,
	intent store.EffectIntent,
) (string, int64, error) {
	return s.base.AppendNextEffectIntent(episodeID, intent)
}

func (s *failRecordStore) AppendEffectOutcome(
	id string,
	outcome store.EffectOutcome,
) (int64, hashref.HashRef, error) {
	s.effectOutcomeCalls++
	return s.base.AppendEffectOutcome(id, outcome)
}

func TestRecordFailureDeliversNoResult(t *testing.T) {
	base := openTestStore(t)
	count := 0
	failing := &failRecordStore{base: base}
	session := newSession(
		failing, "record-failure", []Capability{{"probe", "/ok", 10, 5}},
		echoRegistry(&count), Live, nil,
	)
	result, ref, err := session.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1},
		[]byte("secret"),
	)
	if err == nil || result != nil || !ref.IsZero() {
		t.Fatalf("Invoke = result %q, ref %s, err %v; want nil, zero, error", result, ref, err)
	}
	if count != 1 {
		t.Fatalf("dispatch count = %d, want 1", count)
	}
	if failing.effectOutcomeCalls != 0 {
		t.Fatalf("AppendEffectOutcome calls = %d, want 0 after record failure",
			failing.effectOutcomeCalls)
	}
	receipt, hasIntent, receiptErr := base.GetEffectReceipt(
		store.EffectInvocationID("record-failure", 0),
	)
	if receiptErr != nil || !hasIntent || receipt.State != store.ReceiptIndeterminate {
		t.Fatalf("record-failure receipt = %#v, hasIntent %v, err %v; want indeterminate",
			receipt, hasIntent, receiptErr)
	}
}

func TestIntentIsDurableBeforeDispatch(t *testing.T) {
	s := openTestStore(t)
	const episodeID = "intent-before-dispatch"
	handler := HandlerFunc(func(
		_ context.Context,
		_ EffectRequest,
		_ []byte,
	) ([]byte, error) {
		receipt, hasIntent, err := s.GetEffectReceipt(store.EffectInvocationID(episodeID, 0))
		if err != nil || !hasIntent || receipt.State != store.ReceiptIndeterminate {
			t.Fatalf("receipt at dispatch = %#v, hasIntent %v, err %v; want indeterminate",
				receipt, hasIntent, err)
		}
		return []byte("done"), nil
	})
	session := NewSession(s, episodeID,
		[]Capability{{"probe", "/ok", 10, 5}}, Registry{"probe": handler})
	if _, _, err := session.Invoke(context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 7}, []byte("request")); err != nil {
		t.Fatal(err)
	}
}

func TestFailedInvokeJournalsResolvedOutcome(t *testing.T) {
	s := openTestStore(t)
	const episodeID = "failed-outcome"
	handlerErr := errors.New("injected handler failure")
	session := NewSession(s, episodeID,
		[]Capability{{"probe", "/ok", 10, 5}}, Registry{"probe": HandlerFunc(
			func(context.Context, EffectRequest, []byte) ([]byte, error) {
				return nil, handlerErr
			},
		)})
	_, recordRef, err := session.Invoke(context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 7}, nil)
	var failed *EffectFailedError
	if !errors.As(err, &failed) || !errors.Is(err, handlerErr) {
		t.Fatalf("Invoke error = %v, want EffectFailedError wrapping handler error", err)
	}
	receipt, hasIntent, err := s.GetEffectReceipt(store.EffectInvocationID(episodeID, 0))
	if err != nil || !hasIntent || receipt.State != store.ReceiptResolved ||
		receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "failed" ||
		receipt.EffectOutcome.RecordRef != recordRef || receipt.EffectOutcome.LogicalTime != 7 {
		t.Fatalf("failed receipt = %#v, hasIntent %v, err %v", receipt, hasIntent, err)
	}
}

func TestAllowedLiveSessionRequiresEpisodeID(t *testing.T) {
	s := openTestStore(t)
	session := NewSession(s, "",
		[]Capability{{"probe", "/ok", 10, 5}}, echoRegistry(new(int)))
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 9}
	_, _, err := session.Invoke(context.Background(), req, []byte("request"))
	if err == nil || err.Error() != "broker: live allowed effect requires an episode ID" {
		t.Fatalf("Invoke error = %v, want empty-episode failure", err)
	}
	if _, ok, getErr := s.GetObject(requestHash(req, []byte("request"))); getErr != nil || ok {
		t.Fatalf("request object after empty episode = ok %v, err %v; want absent", ok, getErr)
	}
}

func TestAllowedLiveSessionWithoutHandlerDoesNotDebitBudget(t *testing.T) {
	s := openTestStore(t)
	session := NewSession(s, "missing-handler",
		[]Capability{{"missing", "/ok", 10, 5}}, Registry{})
	req := EffectRequest{Effect: "missing", Scope: "/ok", Cost: 2, Now: 9}
	_, _, err := session.Invoke(context.Background(), req, []byte("request"))
	if err == nil || err.Error() != `broker: no handler registered for "missing"` {
		t.Fatalf("Invoke error = %v, want missing-handler failure", err)
	}
	if got := session.grants[0].Budget; got != 5 {
		t.Fatalf("budget after missing handler = %d, want 5 (not debited)", got)
	}
}

func TestRecordGoldenBytes(t *testing.T) {
	rec := EffectRecord{
		Effect: "FS.Read", Scope: "/project/a", Cost: 2,
		BudgetBefore: 5, BudgetAfter: 3, Allowed: true,
		RequestRef: hashref.SumSHA256([]byte("request")),
		ResultRef:  hashref.SumSHA256([]byte("result")),
	}
	const want = `{"effect":"FS.Read","scope":"/project/a","cost":2,"budgetBefore":5,"budgetAfter":3,"allowed":true,"failed":false,"denial":"","requestRef":"sha256:1f58b9145b24d108d7ac38887338b3ea3229833b9c1e418250343f907bfd1047","resultRef":"sha256:f6a214f7a5fcda0c2cee9660b7fc29f5649e3c68aad48e20e950137c98913a68"}`
	if got := string(EncodeRecord(rec)); got != want {
		t.Fatalf("golden bytes differ:\ngot  %s\nwant %s", got, want)
	}
	roundTrip, err := DecodeRecord([]byte(want))
	if err != nil {
		t.Fatal(err)
	}
	if roundTrip != rec {
		t.Fatalf("round trip = %#v, want %#v", roundTrip, rec)
	}

	failed := EffectRecord{
		Effect: "FS.Read", Scope: "/project/a", Cost: 2,
		BudgetBefore: 5, BudgetAfter: 3, Allowed: true, Failed: true,
		RequestRef: hashref.SumSHA256([]byte("request")),
	}
	const wantFailed = `{"effect":"FS.Read","scope":"/project/a","cost":2,"budgetBefore":5,"budgetAfter":3,"allowed":true,"failed":true,"denial":"","requestRef":"sha256:1f58b9145b24d108d7ac38887338b3ea3229833b9c1e418250343f907bfd1047","resultRef":""}`
	if got := string(EncodeRecord(failed)); got != wantFailed {
		t.Fatalf("failed golden bytes differ:\ngot  %s\nwant %s", got, wantFailed)
	}
	failedRoundTrip, err := DecodeRecord([]byte(wantFailed))
	if err != nil {
		t.Fatal(err)
	}
	if failedRoundTrip != failed {
		t.Fatalf("failed round trip = %#v, want %#v", failedRoundTrip, failed)
	}
}

func TestRecordConsistentAllSketchArms(t *testing.T) {
	resultRef := hashref.SumSHA256([]byte("result"))
	cases := []struct {
		rec  EffectRecord
		want bool
	}{
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Denial: "budget"}, true},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Denial: "budget"}, false},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Failed: true, Denial: "budget"}, false},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, ResultRef: resultRef}, true},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Allowed: true, ResultRef: resultRef}, false},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true}, false},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, Failed: true}, true},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Allowed: true, Failed: true}, false},
		{EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, Failed: true, ResultRef: resultRef}, false},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got := RecordConsistent(tc.rec); got != tc.want {
				t.Fatalf("RecordConsistent = %v, want %v", got, tc.want)
			}
		})
	}
}

func decodeStoredRecord(t *testing.T, s *store.Store, ref hashref.HashRef) EffectRecord {
	t.Helper()
	obj, ok, err := s.GetObject(ref)
	if err != nil || !ok {
		t.Fatalf("GetObject(%s) = ok %v, err %v", ref, ok, err)
	}
	rec, err := DecodeRecord(obj.Payload)
	if err != nil {
		t.Fatal(err)
	}
	return rec
}

func TestThreeArmsDistinguishableFromBytesAlone(t *testing.T) {
	s := openTestStore(t)
	invoke := func(scope string, handler Handler) hashref.HashRef {
		t.Helper()
		session := NewSession(
			s, "decision-cases",
			[]Capability{{Effect: "probe", Scope: "/ok", ExpiresAt: 10, Budget: 5}},
			Registry{"probe": handler},
		)
		_, ref, _ := session.Invoke(
			context.Background(),
			EffectRequest{Effect: "probe", Scope: scope, Cost: 2, Now: 1},
			[]byte("input"),
		)
		if ref.IsZero() {
			t.Fatalf("Invoke scope %q returned zero record ref", scope)
		}
		return ref
	}

	deniedRef := invoke("/denied", ProbeHandler{})
	successRef := invoke("/ok", ProbeHandler{})
	failedRef := invoke("/ok", HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		return nil, errors.New("partial failure")
	}))
	records := []EffectRecord{
		decodeStoredRecord(t, s, deniedRef),
		decodeStoredRecord(t, s, successRef),
		decodeStoredRecord(t, s, failedRef),
	}
	pairs := [][2]bool{
		{records[0].Allowed, records[0].Failed},
		{records[1].Allowed, records[1].Failed},
		{records[2].Allowed, records[2].Failed},
	}
	want := [][2]bool{{false, false}, {true, false}, {true, true}}
	if !reflect.DeepEqual(pairs, want) {
		t.Fatalf("(allowed, failed) pairs decoded from record bytes = %v, want %v", pairs, want)
	}
}

func TestLedgerReconstructibleFromRecordStreamOnFailure(t *testing.T) {
	s := openTestStore(t)
	calls := 0
	session := NewSession(
		s, "budget-irrelevant",
		[]Capability{{Effect: "probe", Scope: "/ok", ExpiresAt: 10, Budget: 5}},
		Registry{"probe": HandlerFunc(func(_ context.Context, _ EffectRequest, _ []byte) ([]byte, error) {
			calls++
			if calls == 2 {
				return nil, errors.New("failed after spending")
			}
			return []byte("ok"), nil
		})},
	)
	_, successRef, err := session.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 2, Now: 1},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, failedRef, err := session.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 3, Now: 1},
		nil,
	)
	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("second Invoke error = %v, want *EffectFailedError", err)
	}

	reconstructed := int64(5)
	for _, ref := range []hashref.HashRef{successRef, failedRef} {
		rec := decodeStoredRecord(t, s, ref)
		if rec.Allowed {
			reconstructed -= rec.Cost
		}
	}
	if got := session.grants[0].Budget; reconstructed != got || got != 0 {
		t.Fatalf("record-stream budget = %d, live budget = %d, want both 0", reconstructed, got)
	}
}

func TestReplayOfFailedRecordReproducesTheFailure(t *testing.T) {
	s := openTestStore(t)
	grants := []Capability{{Effect: "probe", Scope: "/ok", ExpiresAt: 10, Budget: 5}}
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 2, Now: 1}
	live := NewSession(s, "replay-failed-source", grants, Registry{"probe": HandlerFunc(
		func(context.Context, EffectRequest, []byte) ([]byte, error) {
			return nil, errors.New("platform-specific handler detail")
		},
	)})
	_, recordRef, liveErr := live.Invoke(context.Background(), req, nil)
	var liveFailed *EffectFailedError
	if !errors.As(liveErr, &liveFailed) {
		t.Fatalf("live error = %v, want *EffectFailedError", liveErr)
	}

	replayDispatches := 0
	replay := NewReplaySession(
		s,
		grants,
		Registry{"probe": ProbeHandler{Dispatches: &replayDispatches}},
		[]hashref.HashRef{recordRef},
	)
	result, replayRef, replayErr := replay.Invoke(context.Background(), req, nil)
	var replayFailed *EffectFailedError
	if result != nil || replayRef != recordRef {
		t.Errorf("replay = result %q, ref %s; want nil, %s", result, replayRef, recordRef)
	}
	if replayDispatches != 0 {
		t.Errorf("replay dispatches = %d, want 0", replayDispatches)
	}
	if !errors.As(replayErr, &replayFailed) {
		t.Fatalf("replay error = %v, want *EffectFailedError (dispatches stayed %d)", replayErr, replayDispatches)
	}
	if liveFailed.Effect != replayFailed.Effect || liveFailed.Scope != replayFailed.Scope ||
		liveFailed.RecordRef != replayFailed.RecordRef {
		t.Fatalf("live failure = %#v, replay failure = %#v", liveFailed, replayFailed)
	}
	if replayFailed.Unwrap() != nil {
		t.Fatalf("replay failure unwrap = %v, want nil", replayFailed.Unwrap())
	}
}

func TestReplayReturnsRecordedBytesWithoutDispatch(t *testing.T) {
	s := openTestStore(t)
	liveDispatches := 0
	grants := []Capability{{"probe", "/ok", 10, 5}}
	const episodeID = "replay-bytes-source"
	live := NewSession(s, episodeID, grants, Registry{
		"probe": ProbeHandler{Dispatches: &liveDispatches},
	})
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 2, Now: 1}
	want, recordRef, err := live.Invoke(context.Background(), req, []byte("input"))
	if err != nil {
		t.Fatal(err)
	}
	if liveDispatches != 1 {
		t.Fatalf("live dispatches = %d, want 1", liveDispatches)
	}

	replayDispatches := 0
	replay := NewReplaySession(s, grants, Registry{
		"probe": ProbeHandler{Dispatches: &replayDispatches},
	}, []hashref.HashRef{recordRef})
	got, gotRef, err := replay.Invoke(context.Background(), req, []byte("ignored"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) || gotRef != recordRef {
		t.Fatalf("replay = %q, %s; want %q, %s", got, gotRef, want, recordRef)
	}
	if replayDispatches != 0 {
		t.Fatalf("replay dispatches = %d, want 0", replayDispatches)
	}
	nextReceipt, hasNext, err := s.GetEffectReceipt(store.EffectInvocationID(episodeID, 1))
	if err != nil || hasNext || nextReceipt.State != store.ReceiptNotStarted {
		t.Fatalf("replay-created receipt = %#v, hasIntent %v, err %v; want none",
			nextReceipt, hasNext, err)
	}
}

func TestReplayGapNeverFallsBackToLive(t *testing.T) {
	s := openTestStore(t)
	dispatches := 0
	missing := hashref.SumSHA256([]byte("missing-record"))
	replay := NewReplaySession(
		s,
		[]Capability{{"probe", "/ok", 10, 5}},
		Registry{"probe": ProbeHandler{Dispatches: &dispatches}},
		[]hashref.HashRef{missing},
	)
	result, ref, err := replay.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1},
		nil,
	)
	var gap *ReplayGapError
	if !errors.As(err, &gap) || result != nil || !ref.IsZero() {
		t.Errorf("Invoke = result %q, ref %s, err %v; want nil, zero, ReplayGapError", result, ref, err)
	}
	if dispatches != 0 {
		t.Fatalf("replay fallback dispatched %d times", dispatches)
	}
}

func TestReplayRejectsMismatchedRequest(t *testing.T) {
	s := openTestStore(t)
	dispatches := 0
	grants := []Capability{{"probe", "/ok", 10, 5}}
	live := NewSession(s, "replay-mismatch-source", grants, Registry{"probe": ProbeHandler{Dispatches: &dispatches}})
	_, recordRef, err := live.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	replay := NewReplaySession(s, grants, nil, []hashref.HashRef{recordRef})
	_, _, err = replay.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 2, Now: 1},
		nil,
	)
	var gap *ReplayGapError
	if !errors.As(err, &gap) {
		t.Fatalf("mismatch error = %v, want ReplayGapError", err)
	}
}
