package broker

import (
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"testing"
	"time"
)

type contextWitnessStore struct {
	*store.Store
	t              *testing.T
	want           context.Context
	objects, heads int
}

func (s *contextWitnessStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ctx != s.want {
		s.t.Fatalf("GetObject did not receive exact caller ctx")
	}
	s.objects++
	return s.Store.GetObject(ctx, ref)
}
func (s *contextWitnessStore) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	if ctx != s.want {
		s.t.Fatalf("GetRegistryHead did not receive exact caller ctx")
	}
	s.heads++
	return s.Store.GetRegistryHead(ctx, name)
}
func TestApprovalContextIdentity(t *testing.T) {
	// Both deadline-free and bounded callers reach identical pipelines; a
	// deadline is supplied only by the test's second control, never production.
	for _, bounded := range []bool{false, true} {
		t.Run(map[bool]string{false: "deadlineFree", true: "bounded"}[bounded], func(t *testing.T) {
			ctx := context.WithValue(context.Background(), struct{}{}, "caller")
			if bounded {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Hour)
				defer cancel()
			}
			base := openTestStore(t)
			w := &contextWitnessStore{Store: base, t: t, want: ctx}
			f := newPublishFixture(t, "https://registry.example", "thread-context")
			plan := attendedPlanFor(f, "ctx-identity")
			ref, err := mintAttendedApproval(ctx, w, plan)
			if err != nil {
				t.Fatal(err)
			}
			// Mint traverses Execute, append, decide, findRequest, walk's request and
			// decision branches, and both sessions. Exact read counts prevent dead arms.
			if w.objects != 5 || w.heads != 4 {
				t.Fatalf("mint read census objects=%d heads=%d, want 5/4", w.objects, w.heads)
			}
			w.objects = 0
			w.heads = 0
			req := EffectRequest{Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: plan.PublishAt}
			sess := newSession(w, plan.EpisodeID+"-publish", []Capability{{Effect: req.Effect, Scope: req.Scope, ExpiresAt: plan.ExpiresAt, Budget: PublishCost}}, Registry{req.Effect: HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) { return []byte("ok"), nil })}, Live, nil)
			if _, _, err := sess.Invoke(ctx, req, plan.Payload(ref)); err != nil {
				t.Fatal(err)
			}
			if w.objects != 2 || w.heads != 0 {
				t.Fatalf("publish read census objects=%d heads=%d, want 2/0", w.objects, w.heads)
			}
		})
	}
}
func TestApprovalCallerCancellation(t *testing.T) {
	base := openTestStore(t)
	h := NewHumanHandler(base)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f := newPublishFixture(t, "https://registry.example", "cancel-context")
	plan := attendedPlanFor(f, "cancel-context")
	ref, err := MintAttendedApproval(context.Background(), base, plan)
	if err != nil {
		t.Fatalf("live mint control: %v", err)
	}
	decision, ok, err := base.GetObject(context.Background(), ref)
	if err != nil || !ok {
		t.Fatal(err)
	}
	var wire approvalDecisionWire
	if err := decodeApprovalJSON(decision.Payload, &wire); err != nil {
		t.Fatal(err)
	}
	request, _ := hashref.Parse(wire.RequestRef)
	tests := []struct {
		name string
		call func() error
	}{
		{"approve", func() error { _, e := h.Execute(ctx, EffectRequest{Effect: EffectHumanApprove}, nil); return e }},
		{"poll", func() error {
			_, e := h.Execute(ctx, EffectRequest{Effect: EffectHumanPollApproval}, mustApprovalJSON(approvalInputWire{RequestRef: request.String()}))
			return e
		}},
		{"decide", func() error { _, e := DecideApproval(ctx, base, request, "approve", "operator", 3); return e }},
		{"mint", func() error { _, e := MintAttendedApproval(ctx, base, plan); return e }},
		{"publish", func() error {
			_, e := validatePublishApproval(ctx, base, plan.Payload(ref), EffectRequest{Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: plan.PublishAt})
			return e
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); !errors.Is(err, context.Canceled) {
				t.Fatalf("read ignored caller cancellation: %v", err)
			}
		})
	}
}
