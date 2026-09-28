package store

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// TestConvertedStoreMethodsM5b proves bounded pool waits and reconciles all six
// writes after the commit cutoff using the settled durable state.
func TestConvertedStoreMethodsM5b(t *testing.T) {
	type method struct {
		name  string
		setup func(*testing.T, *Store) (func(context.Context) error, func(*testing.T))
		write bool
	}
	methods := []method{
		{"CompareAndSetRegistryHead", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			o := obj("m5b-cas", "test/object")
			if err := s.PutObject(context.Background(), o); err != nil {
				t.Fatal(err)
			}
			return func(ctx context.Context) error {
					return s.CompareAndSetRegistryHead(ctx, "m5b/cas", hashref.HashRef{}, o.Hash)
				}, func(t *testing.T) {
					got, ok, err := s.GetRegistryHead(context.Background(), "m5b/cas")
					if err != nil || !ok || got != o.Hash {
						t.Fatalf("CAS not durable: %v %v %v", got, ok, err)
					}
				}
		}, true},
		{"MintSession", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			row := SessionRow{CredentialID: "m5b-credential", EpisodeID: "ep", GrantsJSON: "[]", ExpiresAt: 100, CreatedAt: 1}
			return func(ctx context.Context) error { return s.MintSession(ctx, row) }, func(t *testing.T) {
				got, ok, err := s.ResolveSession(context.Background(), row.CredentialID)
				if err != nil || !ok || got != row {
					t.Fatalf("mint not durable: %+v %v %v", got, ok, err)
				}
			}
		}, true},
		{"PendingIntents", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			return func(ctx context.Context) error { _, err := s.PendingIntents(ctx, 1); return err }, nil
		}, false},
		{"PendingEffectIntents", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			return func(ctx context.Context) error { _, err := s.PendingEffectIntents(ctx, 1); return err }, nil
		}, false},
		{"ScanUnreadableLog", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			return func(ctx context.Context) error { _, err := s.ScanUnreadableLog(ctx, 0, 1); return err }, nil
		}, false},
		{"ScanUnreadableWorlds", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			return func(ctx context.Context) error { _, err := s.ScanUnreadableWorlds(ctx, "", 1); return err }, nil
		}, false},
		{"PutVerifyResult", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			row := VerifyResult{TransitionFn: hashref.SumSHA256([]byte("m5b-fn")), Interpreter: hashref.SumSHA256([]byte("m5b-interp")), SemanticsEpoch: 1, Verified: true, Detail: "m5b"}
			return func(ctx context.Context) error { return s.PutVerifyResult(ctx, row) }, func(t *testing.T) {
				got, ok, err := s.GetVerifyResult(context.Background(), row.TransitionFn, row.Interpreter)
				if err != nil || !ok || got != row {
					t.Fatalf("verify result not durable: %+v %v %v", got, ok, err)
				}
			}
		}, true},
		{"PutWorld", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			w := World{Ref: hashref.SumSHA256([]byte("m5b-world")), Revision: 1, StateRoot: hashref.SumSHA256([]byte("m5b-state")), LogHead: hashref.SumSHA256([]byte("m5b-log"))}
			return func(ctx context.Context) error { return s.PutWorld(ctx, w) }, func(t *testing.T) {
				got, ok, err := s.GetWorld(context.Background(), w.Ref)
				if err != nil || !ok || got != w {
					t.Fatalf("world not durable: %+v %v %v", got, ok, err)
				}
			}
		}, true},
		{"SelectHead", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			ref := hashref.SumSHA256([]byte("m5b-selected"))
			return func(ctx context.Context) error { return s.SelectHead(ctx, ref) }, func(t *testing.T) {
				got, ok, err := s.SelectedHead(context.Background())
				if err != nil || !ok || got != ref {
					t.Fatalf("head not durable: %v %v %v", got, ok, err)
				}
			}
		}, true},
		{"AppendOutcome", func(t *testing.T, s *Store) (func(context.Context) error, func(*testing.T)) {
			id := "m5b-outcome"
			c := journalCommitFixture(t, s, id)
			if _, _, err := s.AppendIntent(context.Background(), id, testCommitIntent(id, c)); err != nil {
				t.Fatal(err)
			}
			outcome := JournalOutcome{InvocationID: id, Status: "committed", ResultRef: c.NextWorld.Ref, LogicalTime: 2}
			return func(ctx context.Context) error { _, _, err := s.AppendOutcome(ctx, id, outcome); return err }, func(t *testing.T) {
				rc, ok, err := s.GetReceipt(context.Background(), id)
				if err != nil || !ok || rc.State != ReceiptResolved {
					t.Fatalf("outcome not durable: %+v %v %v", rc, ok, err)
				}
			}
		}, true},
	}
	for _, tc := range methods {
		t.Run("held/"+tc.name, func(t *testing.T) {
			s := openFileStore(t)
			call, _ := tc.setup(t, s)
			conn, err := s.db.Conn(context.Background())
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
