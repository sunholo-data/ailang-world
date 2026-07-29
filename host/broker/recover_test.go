package broker

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const recoveryProbeEffect = "Model.Infer"

type recoveryCountingProbe struct {
	dispatches int
}

func (p *recoveryCountingProbe) Execute(
	context.Context,
	EffectRequest,
	[]byte,
) ([]byte, error) {
	p.dispatches++
	return nil, nil
}

type recoveryStoreProbe struct {
	*store.Store
	limits []int
}

func (p *recoveryStoreProbe) PendingIntents(
	limit int,
	fromIndex ...int64,
) ([]store.PendingIntent, error) {
	p.limits = append(p.limits, limit)
	return p.Store.PendingIntents(limit, fromIndex...)
}

func openRecoveryStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatalf("open recovery store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func recoveryCommitFixture(t *testing.T, s *store.Store, label string) store.Commit {
	t.Helper()
	genesis := store.World{
		Ref:       hashref.SumSHA256([]byte("recovery-genesis-world-" + label)),
		Revision:  0,
		StateRoot: hashref.SumSHA256([]byte("recovery-genesis-state-" + label)),
		LogHead:   hashref.SumSHA256([]byte("recovery-genesis-log-" + label)),
	}
	if err := s.PutWorld(genesis); err != nil {
		t.Fatalf("PutWorld genesis: %v", err)
	}
	if err := s.SelectHead(genesis.Ref); err != nil {
		t.Fatalf("SelectHead genesis: %v", err)
	}
	transitionPayload := []byte("planned-transition-" + label)
	transition := store.Object{
		Hash:          hashref.SumSHA256(transitionPayload),
		InterfaceHash: hashref.SumSHA256([]byte("world/recovery-transition/v1")),
		SemanticID:    "world/recovery-transition/v1",
		Provenance:    "host/broker/recover_test",
		Payload:       transitionPayload,
	}
	entryHash := hashref.SumSHA256([]byte("planned-entry-" + label))
	return store.Commit{
		InvocationID: "recover-" + label,
		ObservedHead: genesis.Ref,
		Objects:      []store.Object{transition},
		NextWorld: store.World{
			Ref:       hashref.SumSHA256([]byte("planned-world-" + label)),
			Revision:  1,
			StateRoot: hashref.SumSHA256([]byte("planned-state-" + label)),
			LogHead:   entryHash,
		},
		Entry: store.LogEntry{
			Header: store.LogHeader{
				EntryIndex:     1,
				SemanticsEpoch: 1,
				TransitionFn:   hashref.SumSHA256([]byte("transition-fn-" + label)),
				Interpreter:    hashref.SumSHA256([]byte("interpreter-" + label)),
				PrevEntryHash:  genesis.LogHead,
				WrittenBy:      "broker-recovery-test",
			},
			EntryHash:     entryHash,
			TransitionRef: transition.Hash,
		},
	}
}

func appendRecoveryIntent(t *testing.T, s *store.Store, c store.Commit) {
	t.Helper()
	intent := store.JournalIntent{
		InvocationID:  c.InvocationID,
		WorldRef:      c.NextWorld.Ref,
		EntryHash:     c.Entry.EntryHash,
		ObservedHead:  c.ObservedHead,
		PrevEntryHash: c.Entry.Header.PrevEntryHash,
		TransitionFn:  c.Entry.Header.TransitionFn,
		TransitionRef: c.Entry.TransitionRef,
		Interpreter:   c.Entry.Header.Interpreter,
		LogicalTime:   41,
	}
	if _, _, err := s.AppendIntent(c.InvocationID, intent); err != nil {
		t.Fatalf("AppendIntent: %v", err)
	}
}

func pendingRecoveryCommit(t *testing.T, label string) (*store.Store, store.Commit) {
	t.Helper()
	s := openRecoveryStore(t)
	c := recoveryCommitFixture(t, s, label)
	appendRecoveryIntent(t, s, c)
	return s, c
}

func TestRecoverCountingProbeDispatchesZeroHandlers(t *testing.T) {
	s, _ := pendingRecoveryCommit(t, "counting-probe")
	probe := &recoveryCountingProbe{}
	findings, err := Recover(s, Registry{recoveryProbeEffect: probe})
	if err != nil {
		t.Fatalf("Recover: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("Recover findings=%d, want 1", len(findings))
	}
	if probe.dispatches != 0 {
		t.Fatalf("recovery dispatched %d handlers, want 0", probe.dispatches)
	}
}

func TestRecoverSurfacesNeverLieLaw(t *testing.T) {
	s, c := pendingRecoveryCommit(t, "never-lie")
	findings, err := Recover(s)
	if err != nil {
		t.Fatalf("Recover: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("Recover findings=%d, want 1", len(findings))
	}
	var indeterminate *IndeterminateEffectError
	if !errors.As(findings[0].Err, &indeterminate) {
		t.Fatalf("recovery finding=%T, want *IndeterminateEffectError", findings[0].Err)
	}
	if indeterminate.InvocationID != c.InvocationID {
		t.Fatalf("InvocationID=%q, want %q", indeterminate.InvocationID, c.InvocationID)
	}
	if indeterminate.PlannedWorldRef != c.NextWorld.Ref ||
		indeterminate.PlannedEntryHash != c.Entry.EntryHash {
		t.Fatalf("planned refs=(%s,%s), want (%s,%s)",
			indeterminate.PlannedWorldRef, indeterminate.PlannedEntryHash,
			c.NextWorld.Ref, c.Entry.EntryHash)
	}
	receipt, ok, err := s.GetReceipt(c.InvocationID)
	if err != nil || !ok || receipt.State == store.ReceiptNotStarted {
		t.Fatalf("receipt=(ok=%v,state=%s,err=%v), durable intent reported not-started",
			ok, receipt.State, err)
	}
}

func TestRecoverModelInferNeverRedispatchesAfterResolution(t *testing.T) {
	s, c := pendingRecoveryCommit(t, "model-infer")
	probe := &recoveryCountingProbe{}
	registry := Registry{recoveryProbeEffect: probe}
	findings, err := Recover(s, registry)
	if err != nil || len(findings) != 1 {
		t.Fatalf("Recover before resolution=(findings=%d,err=%v), want (1,nil)",
			len(findings), err)
	}
	outcome := store.JournalOutcome{
		InvocationID: c.InvocationID,
		Status:       "abandoned",
		ResultRef:    hashref.SumSHA256([]byte("operator-abandoned-" + c.InvocationID)),
		LogicalTime:  42,
	}
	if _, _, err := s.AppendOutcome(c.InvocationID, outcome); err != nil {
		t.Fatalf("AppendOutcome: %v", err)
	}
	findings, err = Recover(s, registry)
	if err != nil || len(findings) != 0 {
		t.Fatalf("Recover after resolution=(findings=%d,err=%v), want (0,nil)",
			len(findings), err)
	}
	if probe.dispatches != 0 {
		t.Fatalf("recovery dispatched %d handlers after resolution, want 0", probe.dispatches)
	}
}

func TestRecoverCommitPathPlannedStateAbsentWithoutOutcome(t *testing.T) {
	s, c := pendingRecoveryCommit(t, "not-committed")
	findings, err := Recover(s)
	if err != nil || len(findings) != 1 {
		t.Fatalf("Recover=(findings=%d,err=%v), want (1,nil)", len(findings), err)
	}
	if _, ok, err := s.GetWorld(c.NextWorld.Ref); err != nil || ok {
		t.Fatalf("GetWorld(planned)=(ok=%v,err=%v), want (false,nil)", ok, err)
	}
	if _, ok, err := s.GetLogEntry(c.Entry.Header.EntryIndex); err != nil || ok {
		t.Fatalf("GetLogEntry(planned)=(ok=%v,err=%v), want (false,nil)", ok, err)
	}
}

func TestRecoverUsesKernelPagingBound(t *testing.T) {
	s, _ := pendingRecoveryCommit(t, "paging")
	probe := &recoveryStoreProbe{Store: s}
	if _, err := recoverPending(probe); err != nil {
		t.Fatalf("recoverPending: %v", err)
	}
	if len(probe.limits) == 0 {
		t.Fatal("PendingIntents was not called")
	}
	for i, limit := range probe.limits {
		if limit != store.MaxPendingIntentsPage {
			t.Fatalf("PendingIntents call %d limit=%d, want store.MaxPendingIntentsPage=%d",
				i, limit, store.MaxPendingIntentsPage)
		}
	}
	for _, limit := range []int{0, store.MaxPendingIntentsPage + 1} {
		t.Run(fmt.Sprintf("invalid-%d", limit), func(t *testing.T) {
			_, err := s.PendingIntents(limit)
			var invalid *store.InvalidLimitError
			if !errors.As(err, &invalid) {
				t.Fatalf("PendingIntents(%d) error=%T %v, want *store.InvalidLimitError",
					limit, err, err)
			}
		})
	}
}

func TestRecoveryConsumerContractMirrorsSketch(t *testing.T) {
	if mayReportNotStarted(true) || !mayReportNotStarted(false) {
		t.Fatal("mayReportNotStarted does not mirror the authoritative sketch")
	}
	for _, row := range []struct {
		indeterminate, reconciled, want bool
	}{
		{false, false, true},
		{true, false, false},
		{true, true, true},
	} {
		if got := retryAllowed(row.indeterminate, row.reconciled); got != row.want {
			t.Fatalf("retryAllowed(%v,%v)=%v, want %v",
				row.indeterminate, row.reconciled, got, row.want)
		}
	}
}
