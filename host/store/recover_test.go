package store

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// IndeterminateEffectError is test-local because SD.C needs to prove the
// consumer contract, while the kernel deliberately exposes facts only through
// GetReceipt and PendingIntents. Adding a production broker error before M3
// would put consumer policy into host/store.
type IndeterminateEffectError struct {
	InvocationID string
}

func (e *IndeterminateEffectError) Error() string {
	return fmt.Sprintf("effect %q is indeterminate: durable intent has no outcome", e.InvocationID)
}

func mayReportNotStarted(hasIntent bool) bool { return !hasIntent }

func retryAllowed(indeterminate, reconciled bool) bool {
	return !indeterminate || reconciled
}

type countingProbe struct {
	dispatches int
}

func (p *countingProbe) dispatch() { p.dispatches++ }

// recoverIndeterminate surfaces the durable ambiguity and never dispatches.
// MUT-AUTO-RETRY changes this function to call probe.dispatch(); two independent
// tests below then red.
func recoverIndeterminate(s *Store, id string, probe *countingProbe) error {
	receipt, ok, err := s.GetReceipt(id)
	if err != nil {
		return err
	}
	if ok && receipt.State == ReceiptIndeterminate {
		return &IndeterminateEffectError{InvocationID: id}
	}
	return nil
}

func pendingRecoveryFixture(t *testing.T) (*Store, Commit) {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	c := journalCommitFixture(t, s, "recover-pending")
	if _, _, err := s.AppendIntent(c.InvocationID, testCommitIntent(c.InvocationID, c)); err != nil {
		t.Fatal(err)
	}
	return s, c
}

func TestRecoverIndeterminateSurfacesNeverLieLaw(t *testing.T) {
	s, c := pendingRecoveryFixture(t)
	probe := &countingProbe{}
	err := recoverIndeterminate(s, c.InvocationID, probe)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("recovery error=%T %v, want *IndeterminateEffectError", err, err)
	}
	if mayReportNotStarted(true) {
		t.Fatal("mayReportNotStarted(true)=true: durable intent was reported not-started")
	}
	if !mayReportNotStarted(false) {
		t.Fatal("mayReportNotStarted(false)=false: no-intent case must permit not-started")
	}
	if probe.dispatches != 0 {
		t.Fatalf("indeterminate recovery dispatched %d times, want zero", probe.dispatches)
	}
}

func TestRecoverRetryAllowedMirrorsAllSketchRows(t *testing.T) {
	// This store-local copy mirrors the sketch only. Production retryAllowed
	// lives in broker, so broker's MUT-RETRY-XOR cannot reach this witness.
	rows := []struct {
		name                      string
		indeterminate, reconciled bool
		want                      bool
	}{
		{"not-indeterminate", false, false, true},
		{"not-indeterminate-reconciled", false, true, true},
		{"indeterminate-unreconciled", true, false, false},
		{"indeterminate-reconciled", true, true, true},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			if got := retryAllowed(row.indeterminate, row.reconciled); got != row.want {
				t.Fatalf("retryAllowed(%v,%v)=%v, want %v",
					row.indeterminate, row.reconciled, got, row.want)
			}
		})
	}
}

func reconcileCommitNotExecuted(t *testing.T, s *Store, c Commit) {
	t.Helper()
	if _, ok, err := s.GetWorld(context.Background(), c.NextWorld.Ref); err != nil || ok {
		t.Fatalf("planned world before reconciliation: ok=%v err=%v, want absent", ok, err)
	}
	if _, ok, err := s.GetLogEntry(context.Background(), c.Entry.Header.EntryIndex); err != nil || ok {
		t.Fatalf("planned entry before reconciliation: ok=%v err=%v, want absent", ok, err)
	}
	outcome := JournalOutcome{
		InvocationID: c.InvocationID,
		Status:       "not-executed",
		ResultRef:    hashref.SumSHA256([]byte("reconciled-not-executed-" + c.InvocationID)),
		LogicalTime:  43,
	}
	if _, _, err := s.AppendOutcome(c.InvocationID, outcome); err != nil {
		t.Fatalf("append reconciling outcome: %v", err)
	}
}

func TestRecoverCommitPathDeterministicallyReconcilesUntouchedStore(t *testing.T) {
	s, c := pendingRecoveryFixture(t)
	reconcileCommitNotExecuted(t, s, c)
	receipt, ok, err := s.GetReceipt(c.InvocationID)
	if err != nil || !ok || receipt.State != ReceiptResolved {
		t.Fatalf("receipt after reconcile=(ok=%v,state=%s,err=%v), want resolved",
			ok, receipt.State, err)
	}
	if receipt.Outcome == nil || receipt.Outcome.Status != "not-executed" {
		t.Fatalf("reconciling outcome=%+v, want not-executed", receipt.Outcome)
	}
	if pending, err := s.PendingIntents(MaxPendingIntentsPage); err != nil || len(pending) != 0 {
		t.Fatalf("pending after reconcile=%d err=%v, want zero", len(pending), err)
	}
}

func TestRecoverModelInferNeverRedispatchesEvenWhenResolutionOffered(t *testing.T) {
	s, c := pendingRecoveryFixture(t)
	probe := &countingProbe{}
	err := recoverIndeterminate(s, c.InvocationID, probe)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("Model.Infer recovery error=%T %v, want indeterminate", err, err)
	}

	// An operator may resolve/abandon a paid inference, but Model.Infer has no
	// deterministic reconciler. Offering that resolution is never permission
	// to dispatch the paid effect again.
	outcome := JournalOutcome{
		InvocationID: c.InvocationID,
		Status:       "abandoned",
		ResultRef:    hashref.SumSHA256([]byte("operator-abandoned-" + c.InvocationID)),
		LogicalTime:  44,
	}
	if _, _, err := s.AppendOutcome(c.InvocationID, outcome); err != nil {
		t.Fatal(err)
	}
	if probe.dispatches != 0 {
		t.Fatalf("Model.Infer re-dispatched %d times after resolution, want zero", probe.dispatches)
	}
}
