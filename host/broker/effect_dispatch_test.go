package broker

// Row 134 M2a/M2b broker seams: BoundInvoker.Unhandled, OpenReplayBinder,
// ReplayGrant, and the detached post-dispatch writes.

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestBoundInvokerUnhandled(t *testing.T) {
	declared := []Requirement{{"A", "s", 1}, {"B", "s", 1}, {"B", "t", 1}, {"C", "s", 0}}
	live := newSession(nil, "ep", nil, Registry{"A": ProbeHandler{}}, Live, nil)
	inv, err := live.Bind(Manifest{Declared: declared})
	if err != nil {
		t.Fatal(err)
	}
	if got := inv.Unhandled(); !reflect.DeepEqual(got, []string{"B", "C"}) {
		t.Fatalf("live Unhandled = %v, want [B C] (declaration order, no duplicates)", got)
	}
	empty, _ := OpenBinder(nil, "ep", nil, nil).Bind(Manifest{Declared: declared})
	if got := empty.Unhandled(); !reflect.DeepEqual(got, []string{"A", "B", "C"}) {
		t.Fatalf("OpenBinder (nil registry) Unhandled = %v, want all", got)
	}
	replay, _ := OpenReplayBinder(nil, nil, nil).Bind(Manifest{Declared: declared})
	if got := replay.Unhandled(); got != nil {
		t.Fatalf("replay Unhandled = %v, want none: replay never dispatches", got)
	}
}

// ReplayGrant reproduces every recorded decision label under Decide.
func TestReplayGrantReproducesRecordedDecision(t *testing.T) {
	const now = 50
	req := EffectRequest{Effect: "E", Scope: "s", Cost: 1, Now: now}
	for _, live := range []Capability{
		{"E", "s", 100, 3},     // allowed
		{"E", "s", 100, 0},     // denied:budget
		{"E", "s", now, 3},     // denied:expired
		{"E", "other", 100, 3}, // denied:scope
		{"F", "s", 100, 3},     // denied:effect-name
	} {
		want := Decide(live, req)
		rec := EffectRecord{Effect: "E", Scope: "s", Cost: 1, BudgetBefore: live.Budget, Allowed: want.Allowed}
		if !want.Allowed {
			rec.Denial = want.Label
		}
		if got := Decide(ReplayGrant(rec, now), req); got != want {
			t.Fatalf("ReplayGrant for %+v decides %+v, want %+v", live, got, want)
		}
	}
}

// OpenReplayBinder answers from the recorded stream with a nil registry.
func TestOpenReplayBinderNeverDispatches(t *testing.T) {
	s := openTestStore(t)
	count := 0
	reg := Registry{"probe": HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		count++
		return []byte(`{"ok":true}`), nil
	})}
	live := newSession(s, "ep", []Capability{{"probe", "/ok", 10, 5}}, reg, Live, nil)
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1}
	out, ref, err := live.Invoke(boundedTestContext(t), req, []byte("p"))
	if err != nil {
		t.Fatal(err)
	}
	rec := decodeStoredRecord(t, s, ref)
	inv, err := OpenReplayBinder(s, []Capability{ReplayGrant(rec, 1)}, []hashref.HashRef{ref}).
		Bind(Manifest{Declared: []Requirement{{"probe", "/ok", 1}}})
	if err != nil {
		t.Fatal(err)
	}
	got, gotRef, err := inv.Request(boundedTestContext(t), req, []byte("p"))
	if err != nil || string(got) != string(out) || gotRef != ref || count != 1 {
		t.Fatalf("replay = %q %s %v (dispatches %d), want recorded bytes and no dispatch", got, gotRef, err, count)
	}
}

// The writes after a handler runs survive the caller's cancellation (V48 was
// "on the request's ctx"): a handler that cancels its ctx and fails, or
// succeeds, still leaves its record and resolved outcome.
func TestPostDispatchWritesAreDetached(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "succeeded", true: "failed"}[fail], func(t *testing.T) {
			s := openTestStore(t)
			ctx, cancel := context.WithCancel(boundedTestContext(t))
			reg := Registry{"probe": HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
				cancel()
				if fail {
					return nil, errors.New("handler failed")
				}
				return []byte(`{"ok":true}`), nil
			})}
			session := newSession(s, "ep", []Capability{{"probe", "/ok", 10, 5}}, reg, Live, nil)
			_, ref, err := session.Invoke(ctx, EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1}, nil)
			var failed *EffectFailedError
			if fail && !errors.As(err, &failed) || !fail && err != nil {
				t.Fatalf("Invoke err = %T %v", err, err)
			}
			if rec := decodeStoredRecord(t, s, ref); rec.Failed != fail || !rec.Allowed {
				t.Fatalf("record = %+v", rec)
			}
			rc, ok, err := s.GetEffectReceipt(boundedTestContext(t), store.EffectInvocationID("ep", 0))
			if err != nil || !ok || rc.State != store.ReceiptResolved || rc.EffectOutcome.RecordRef != ref {
				t.Fatalf("receipt = %+v ok=%v err=%v, want resolved naming %s", rc, ok, err, ref)
			}
		})
	}
}

// A post-dispatch write failure is the typed OutcomeWriteError, keeping the
// landed message.
func TestPostDispatchWriteFailureIsTyped(t *testing.T) {
	count := 0
	session := newSession(&failRecordStore{base: openTestStore(t)}, "ep",
		[]Capability{{"probe", "/ok", 10, 5}}, echoRegistry(&count), Live, nil)
	_, _, err := session.Invoke(boundedTestContext(t), EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 1}, nil)
	var unwritten *OutcomeWriteError
	if !errors.As(err, &unwritten) || unwritten.Effect != "probe" || count != 1 ||
		err.Error() != "broker: put effect record: injected record write failure" {
		t.Fatalf("err = %T %v, want *OutcomeWriteError after one dispatch", err, err)
	}
}
