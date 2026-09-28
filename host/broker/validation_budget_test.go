package broker

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type blockedApprovalRead struct{ *store.Store }

func (s blockedApprovalRead) GetObject(ctx context.Context, _ hashref.HashRef) (store.Object, bool, error) {
	<-ctx.Done()
	return store.Object{}, false, ctx.Err()
}

func TestPublishValidationBudgetOnly(t *testing.T) {
	body, err := os.ReadFile("broker.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(body)
	for _, phrase := range []string{"context.WithTimeout(ctx, 3*time.Second)", "validatePublishApproval(validationCtx, s.store, payload, req)", "cancelValidation()"} {
		if !strings.Contains(src, phrase) {
			t.Fatalf("missing validation-only budget step %q", phrase)
		}
	}
	t.Run("validation expires", func(t *testing.T) {
		base := openTestStore(t)
		f := newPublishFixture(t, "https://registry.example", "validation-budget")
		plan := attendedPlanFor(f, "validation-budget")
		ref, err := MintAttendedApproval(context.Background(), base, plan)
		if err != nil {
			t.Fatal(err)
		}
		req := EffectRequest{Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: plan.PublishAt}
		s := newSession(blockedApprovalRead{base}, plan.EpisodeID+"-publish", []Capability{{Effect: req.Effect, Scope: req.Scope, ExpiresAt: plan.ExpiresAt, Budget: PublishCost}}, Registry{req.Effect: HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
			t.Fatal("handler reached after validation expiry")
			return nil, nil
		})}, Live, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		start := time.Now()
		_, _, err = s.Invoke(ctx, req, plan.Payload(ref))
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("validation result = %v", err)
		}
		if elapsed := time.Since(start); elapsed > 4*time.Second {
			t.Fatalf("validation took %v; B3 is 3s", elapsed)
		}
		if ctx.Err() != nil {
			t.Fatal("validation consumed enclosing caller budget")
		}
	})
	t.Run("dispatch remains independent", func(t *testing.T) {
		base := openTestStore(t)
		f := newPublishFixture(t, "https://registry.example", "dispatch-budget")
		plan := attendedPlanFor(f, "dispatch-budget")
		ref, err := MintAttendedApproval(context.Background(), base, plan)
		if err != nil {
			t.Fatal(err)
		}
		req := EffectRequest{Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: plan.PublishAt}
		s := newSession(base, plan.EpisodeID+"-publish", []Capability{{Effect: req.Effect, Scope: req.Scope, ExpiresAt: plan.ExpiresAt, Budget: PublishCost}}, Registry{req.Effect: HandlerFunc(func(ctx context.Context, _ EffectRequest, _ []byte) ([]byte, error) {
			time.Sleep(3100 * time.Millisecond)
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		})}, Live, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()
		got, _, err := s.Invoke(ctx, req, plan.Payload(ref))
		if err != nil || string(got) != "ok" {
			t.Fatalf("dispatch after >B3 = %q, %v", got, err)
		}
	})
}
