package broker

import (
	"context"
	"errors"
	"testing"
)

// AC-BINDER: the binder's BoundInvoker refuses an undeclared triple before
// the broker pipeline (no store access is needed to observe the refusal).
func TestOpenBinder(t *testing.T) {
	b := OpenBinder(nil, "episode-1", []Capability{{Effect: "IO", Scope: "x", ExpiresAt: 10, Budget: 5}})
	inv, err := b.Bind(Manifest{Access: Requirement{Effect: "IO", Scope: "x", Cost: 1}})
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	_, _, err = inv.Request(context.Background(), EffectRequest{Effect: "IO", Scope: "x", Cost: 1}, nil)
	var undeclared *UndeclaredEffectError
	if !errors.As(err, &undeclared) {
		t.Fatalf("Request on an empty declaration set = %v, want *UndeclaredEffectError", err)
	}
	if _, err := b.Bind(Manifest{Access: Requirement{Cost: -1}}); err == nil {
		t.Fatal("Bind accepted a negative access cost; the binder does not reach Session.Bind")
	}
}
