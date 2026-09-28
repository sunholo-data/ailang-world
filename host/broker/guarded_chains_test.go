package broker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestGuardedApprovalChain(t *testing.T) {
	db := openTestStore(t)
	f := newPublishFixture(t, "https://registry.example", "guard-approval")
	plan := attendedPlanFor(f, "guard-approval")
	if _, err := MintAttendedApproval(context.Background(), db, plan); !errors.Is(err, store.ErrNoDeadline) {
		t.Fatalf("unbounded mint = %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ref, err := MintAttendedApproval(ctx, db, plan)
	if err != nil {
		t.Fatalf("bounded mint = %v", err)
	}
	approval, ok, err := db.GetObject(ctx, ref)
	if err != nil || !ok {
		t.Fatalf("read minted approval: ok=%v err=%v", ok, err)
	}
	if err := db.PutObject(context.Background(), approval); !errors.Is(err, store.ErrNoDeadline) {
		t.Fatalf("unbounded approval write = %v", err)
	}
}

func TestGuardedPublishValidationChain(t *testing.T) {
	db := openTestStore(t)
	f := newPublishFixture(t, "https://registry.example", "guard-validation")
	plan := attendedPlanFor(f, "guard-validation")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ref, err := MintAttendedApproval(ctx, db, plan)
	if err != nil {
		t.Fatal(err)
	}
	req := EffectRequest{Effect: EffectRegistryPublish, Scope: plan.PublishScope(), Cost: PublishCost, Now: plan.PublishAt}
	if _, err := validatePublishApproval(context.Background(), db, plan.Payload(ref), req); !errors.Is(err, store.ErrNoDeadline) {
		t.Fatalf("unbounded validation = %v", err)
	}
	if _, err := validatePublishApproval(ctx, db, plan.Payload(ref), req); err != nil {
		t.Fatalf("bounded validation = %v", err)
	}
}
