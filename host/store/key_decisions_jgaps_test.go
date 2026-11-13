package store

// Row 201 A4 sprint 3 M6a (the sprint-2 judge's findings 1 and 2): LandKeyDecision lands only the
// request it names, and K0's exact match needs the same approver. The fixtures are
// key_decisions_test.go's; every store is under t.TempDir() (D-WORLD-109).

import (
	"reflect"
	"testing"
)

// JM22: with two unlanded rows, landing one leaves the other unlanded (the UPDATE's request_ref
// filter), and the other then lands with its own decision.
func TestLandKeyDecisionLandsOnlyItsRequest(t *testing.T) {
	c := kdStore(t)
	other := kdRow()
	other.RequestRef = kdD
	for _, row := range []KeyDecisionRow{kdRow(), other} {
		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1000); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.LandKeyDecision(testCtx(t), kdQ, kdP); err != nil {
		t.Fatal(err)
	}
	got, ok, err := c.KeyDecision(testCtx(t), kdD)
	if err != nil || !ok || got.DecisionRef != "" {
		t.Fatalf("the other request after landing %s: %+v %v %v, want it unlanded", kdQ, got, ok, err)
	}
	if first, _ := kdRead(t, c); first.DecisionRef != kdP {
		t.Fatalf("the landed request: %+v", first)
	}
	if err := c.LandKeyDecision(testCtx(t), kdD, kdQ); err != nil {
		t.Fatalf("the other request lands on its own: %v", err)
	}
}

// Finding 2 (K0's approver operand at the store API): the same artifact and signed fields recorded
// for another approver is not the exact match: the unlanded row is replaced (the new approver,
// the new decided_at), never kept as the first approver's.
func TestRecordKeyDecisionExactMatchNeedsTheSameApprover(t *testing.T) {
	c := kdStore(t)
	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "ann-ssh", Format: "ssh", PublicKey: []byte("k2"), EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
		t.Fatal(err)
	}
	row := kdRow()
	row.Approver = "ssh:ann-ssh"
	if err := c.RecordKeyDecision(testCtx(t), row, ApproverUse{Name: "ann-ssh", Format: "ssh", PublicKey: []byte("k2")}, 1010); err != nil {
		t.Fatal(err)
	}
	row.DecidedAt = 1010
	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, row) {
		t.Fatalf("another approver's record: %+v, want the row replaced by %+v", got, row)
	}
}
