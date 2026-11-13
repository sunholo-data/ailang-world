package daemon

// Row 201 A4 sprint 3 M6a (the sprint-2 judge's finding 3, JM01-JM04): the writer's approver race
// (K3 in RecordKeyDecision's transaction) answers 422 NotApproved; any other store error is not a
// race. Through the POST the read-only judge refuses such an approver first (plan §3.1, D-RACE), so
// the mapping is pinned here, operand by operand, as the errors RecordKeyDecision wraps them.

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestIsApproverRaceOperands(t *testing.T) {
	for _, base := range []error{store.ErrApproverRevoked, store.ErrApproverKeyChanged, store.ErrApproverUnknown, store.ErrApproverClockBehind} {
		if err := fmt.Errorf("%w: %q", base, "mark"); !isApproverRace(err) {
			t.Errorf("%v: not a race", err)
		}
	}
	for _, err := range []error{store.ErrBadKeyDecision, store.ErrKeyDecisionLanded, errors.New("store: record key decision: disk I/O error")} {
		if isApproverRace(err) {
			t.Errorf("%v: a race", err)
		}
	}
}
