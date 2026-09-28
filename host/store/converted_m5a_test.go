package store

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// TestConvertedStoreMethodsM5a exercises the pool wait and the durable cutoff
// separately for each newly converted method. Every cutoff arm reads the row
// after the worker has settled, so an uncertain result is reconciled.
func TestConvertedStoreMethodsM5a(t *testing.T) {
	type method struct {
		name  string
		setup func(*testing.T, *Store) (func(context.Context) error, func(*testing.T))
		write bool
	}
	methods := []method{
		{"PutObject", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			o := obj("m5a-object", "test/object")
			return func(ctx context.Context) error { return s.PutObject(ctx, o) }, func(t *testing.T) {
				_, ok, err := s.GetObject(boundedTestContext(t), o.Hash)
				if err != nil || !ok {
					t.Fatalf("object not durable: ok=%v err=%v", ok, err)
				}
			}
		}, true},
		{"SetRegistryHead", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			ref := hashref.SumSHA256([]byte("m5a-head"))
			return func(ctx context.Context) error { return s.SetRegistryHead(ctx, "m5a/head", ref) }, func(t *testing.T) {
				got, ok, err := s.GetRegistryHead(boundedTestContext(t), "m5a/head")
				if err != nil || !ok || got != ref {
					t.Fatalf("head not durable: %v %v %v", got, ok, err)
				}
			}
		}, true},
		{"AppendNextEffectIntent", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			return func(ctx context.Context) error {
					_, _, err := s.AppendNextEffectIntent(ctx, "m5a-next", effectIntentFixture("m5a-next", 1))
					return err
				}, func(t *testing.T) {
					rc, ok, err := s.GetEffectReceipt(boundedTestContext(t), EffectInvocationID("m5a-next", 0))
					if err != nil || !ok || rc.State != ReceiptIndeterminate {
						t.Fatalf("intent not durable: %+v %v %v", rc, ok, err)
					}
				}
		}, true},
		{"AppendClaimedEffectIntent", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			approval := hashref.SumSHA256([]byte("m5a-approval"))
			request := hashref.SumSHA256([]byte("m5a-request"))
			intent := effectIntentFixture("m5a-claimed", 1)
			intent.RequestRef = request
			return func(ctx context.Context) error {
					_, _, err := s.AppendClaimedEffectIntent(ctx, "m5a-claimed", intent, approval, request)
					return err
				}, func(t *testing.T) {
					rc, ok, err := s.GetEffectReceipt(boundedTestContext(t), EffectInvocationID("m5a-claimed", 0))
					if err != nil || !ok || rc.State != ReceiptIndeterminate {
						t.Fatalf("claimed intent not durable: %+v %v %v", rc, ok, err)
					}
				}
		}, true},
		{"AppendEffectOutcome", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			id, _, err := s.AppendNextEffectIntent(boundedTestContext(t), "m5a-outcome", effectIntentFixture("m5a-outcome", 1))
			if err != nil {
				t.Fatal(err)
			}
			outcome := EffectOutcome{InvocationID: id, Status: "succeeded", RecordRef: hashref.SumSHA256([]byte("m5a-record")), LogicalTime: 2}
			return func(ctx context.Context) error { _, _, err := s.AppendEffectOutcome(ctx, id, outcome); return err }, func(t *testing.T) {
				rc, ok, err := s.GetEffectReceipt(boundedTestContext(t), id)
				if err != nil || !ok || rc.State != ReceiptResolved {
					t.Fatalf("outcome not durable: %+v %v %v", rc, ok, err)
				}
			}
		}, true},
		{"GetEffectReceipt", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			id, _, err := s.AppendNextEffectIntent(boundedTestContext(t), "m5a-receipt", effectIntentFixture("m5a-receipt", 1))
			if err != nil {
				t.Fatal(err)
			}
			return func(ctx context.Context) error { _, _, err := s.GetEffectReceipt(ctx, id); return err }, nil
		}, false},
	}
	for _, tc := range methods {
		t.Run("held/"+tc.name, func(t *testing.T) {
			s := openFileStore(t)
			call, _ := tc.setup(t, s)
			conn, err := s.db.Conn(boundedTestContext(t))
			if err != nil {
				t.Fatal(err)
			}
			released := false
			release := func() {
				if !released {
					released = true
					_ = conn.Close()
				}
			}
			defer release()
			ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
			defer cancel()
			got, elapsed := boundedCall(t, func() error { return call(ctx) }, release)
			if !errors.Is(got, context.DeadlineExceeded) || IsUncertain(got) {
				t.Fatalf("held result = %v, want definite deadline", got)
			}
			if elapsed > durableBudget+durableEpsilon {
				t.Fatalf("held elapsed = %v", elapsed)
			}
		})
		if !tc.write {
			continue
		}
		t.Run("cutoff/"+tc.name, func(t *testing.T) {
			s := openFileStore(t)
			call, reconcile := tc.setup(t, s)
			release := make(chan struct{})
			prev := durableAfterCommitHook
			durableAfterCommitHook = func() { <-release }
			defer func() { releaseOnce(release); s.workers.Wait(); durableAfterCommitHook = prev }()
			ctx, cancel := context.WithTimeout(context.Background(), durableBudget)
			defer cancel()
			got, elapsed := boundedCall(t, func() error { return call(ctx) }, func() { releaseOnce(release) })
			if !IsUncertain(got) {
				t.Fatalf("post-cutoff result = %v, want uncertain", got)
			}
			if elapsed > durableBudget+durableEpsilon {
				t.Fatalf("cutoff elapsed = %v", elapsed)
			}
			releaseOnce(release)
			s.workers.Wait()
			reconcile(t)
		})
	}
}
