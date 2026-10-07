package daemon

import (
	"testing"

	"github.com/sunholo-data/ailang-world/host/coordinator"
)

// TestPlanPhaseBudgetIsDerived is row 153 AC2.1: the plan cap is the largest
// one that leaves the handler its full cap inside the daemon's invocation
// deadline, so the number cannot drift from the budgets it is derived from
// (the derivation is also the comment on coordinator.PlanPhaseBudget).
func TestPlanPhaseBudgetIsDerived(t *testing.T) {
	want := invokeDeadline - coordinator.HandlerCap - coordinator.HandlerHeadroom
	if coordinator.PlanPhaseBudget != want {
		t.Fatalf("PlanPhaseBudget = %s, want invokeDeadline − HandlerCap − HandlerHeadroom = %s − %s − %s = %s",
			coordinator.PlanPhaseBudget, invokeDeadline, coordinator.HandlerCap, coordinator.HandlerHeadroom, want)
	}
}
