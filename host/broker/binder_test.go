package broker

import (
	"errors"
	"reflect"
	"testing"
)

// AC-BINDER: the binder's BoundInvoker refuses an undeclared triple before
// the broker pipeline (no store access is needed to observe the refusal).
func TestOpenBinder(t *testing.T) {
	b := OpenBinder(nil, "episode-1", []Capability{{Effect: "IO", Scope: "x", ExpiresAt: 10, Budget: 5}}, nil)
	inv, err := b.Bind(Manifest{Access: Requirement{Effect: "IO", Scope: "x", Cost: 1}})
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	_, _, err = inv.Request(boundedTestContext(t), EffectRequest{Effect: "IO", Scope: "x", Cost: 1}, nil)
	var undeclared *UndeclaredEffectError
	if !errors.As(err, &undeclared) {
		t.Fatalf("Request on an empty declaration set = %v, want *UndeclaredEffectError", err)
	}
	if _, err := b.Bind(Manifest{Access: Requirement{Cost: -1}}); err == nil {
		t.Fatal("Bind accepted a negative access cost; the binder does not reach Session.Bind")
	}
}

// Row 134 M4b: OpenBinder runs exactly the registry it was given, copied at
// open — a later change to the caller's map neither widens nor narrows it.
func TestOpenBinderUsesACopyOfTheGivenRegistry(t *testing.T) {
	declared := []Requirement{{"A", "s", 1}, {"B", "s", 1}}
	reg := Registry{"A": ProbeHandler{}}
	b := OpenBinder(nil, "ep", nil, reg)
	reg["B"] = ProbeHandler{}
	delete(reg, "A")
	inv, err := b.Bind(Manifest{Declared: declared})
	if err != nil {
		t.Fatal(err)
	}
	if got := inv.Unhandled(); !reflect.DeepEqual(got, []string{"B"}) {
		t.Fatalf("Unhandled = %v, want [B]: the binder must hold the registry as it was at open", got)
	}
	empty, _ := OpenBinder(nil, "ep", nil, Registry{}).Bind(Manifest{Declared: declared})
	if got := empty.Unhandled(); !reflect.DeepEqual(got, []string{"A", "B"}) {
		t.Fatalf("empty registry Unhandled = %v, want [A B]", got)
	}
}
