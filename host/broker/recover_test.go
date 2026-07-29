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

// pagedRecoveryStore scripts multi-page PendingIntents responses so the keyset
// cursor and the loop's termination are exercised. The temp-file fixtures above
// hold ONE pending intent, so they only ever produce a single short page: with
// them, dropping the cursor entirely changes nothing observable and
// MUT-PENDING-UNBOUNDED cannot red. Measured at iteration 35 — the mutation was
// green against all five original cases. This fake is what gives the paging
// discipline something to prove.
type pagedRecoveryStore struct {
	pages      [][]store.PendingIntent
	fromIndex  []int64 // the cursor received on each call; -1 means "no cursor"
	calls      int
	maxCalls   int
	sawReceipt map[string]bool
}

func (p *pagedRecoveryStore) PendingIntents(
	limit int,
	fromIndex ...int64,
) ([]store.PendingIntent, error) {
	if limit != store.MaxPendingIntentsPage {
		return nil, fmt.Errorf("limit=%d, want store.MaxPendingIntentsPage", limit)
	}
	cursor := int64(-1)
	if len(fromIndex) > 0 {
		cursor = fromIndex[0]
	}
	p.fromIndex = append(p.fromIndex, cursor)
	p.calls++
	if p.calls > p.maxCalls {
		return nil, fmt.Errorf(
			"recovery did not terminate: %d calls, cursors=%v", p.calls, p.fromIndex)
	}
	// Serve pages positionally, but ONLY when the cursor actually advanced the
	// way a keyset scan requires; a stuck cursor re-serves page 0 forever,
	// which is precisely the unbounded-rescan defect.
	idx := 0
	for i, page := range p.pages {
		if len(page) == 0 {
			continue
		}
		if cursor < page[len(page)-1].Seq {
			idx = i
			break
		}
		idx = i + 1
	}
	if idx >= len(p.pages) {
		return nil, nil
	}
	return p.pages[idx], nil
}

func (p *pagedRecoveryStore) GetReceipt(id string) (store.Receipt, bool, error) {
	if p.sawReceipt == nil {
		p.sawReceipt = map[string]bool{}
	}
	p.sawReceipt[id] = true
	return store.Receipt{InvocationID: id, State: store.ReceiptIndeterminate}, true, nil
}

func TestRecoverPagesWithKeysetCursorAcrossFullPages(t *testing.T) {
	full := make([]store.PendingIntent, store.MaxPendingIntentsPage)
	for i := range full {
		seq := int64(i + 1)
		full[i] = store.PendingIntent{
			Seq:          seq,
			InvocationID: fmt.Sprintf("paged-%d", seq),
		}
	}
	tail := []store.PendingIntent{{Seq: int64(store.MaxPendingIntentsPage + 1),
		InvocationID: fmt.Sprintf("paged-%d", store.MaxPendingIntentsPage+1)}}

	probe := &pagedRecoveryStore{
		pages:    [][]store.PendingIntent{full, tail},
		maxCalls: 8, // a stuck cursor blows this long before maxRecoveryPages
	}
	findings, err := recoverPending(probe)
	if err != nil {
		t.Fatalf("recoverPending across pages: %v", err)
	}
	want := store.MaxPendingIntentsPage + 1
	if len(findings) != want {
		t.Fatalf("findings=%d, want %d (one per pending intent across both pages)",
			len(findings), want)
	}
	// The FIRST call carries no cursor; every later call must carry the last
	// Seq of the page before it. That is the keyset contract.
	if len(probe.fromIndex) < 2 {
		t.Fatalf("PendingIntents called %d times, want at least 2 pages",
			len(probe.fromIndex))
	}
	if probe.fromIndex[0] != -1 {
		t.Errorf("first PendingIntents cursor=%d, want none", probe.fromIndex[0])
	}
	if probe.fromIndex[1] != int64(store.MaxPendingIntentsPage) {
		t.Errorf("second PendingIntents cursor=%d, want %d — the keyset cursor did not advance",
			probe.fromIndex[1], store.MaxPendingIntentsPage)
	}
	if !probe.sawReceipt[fmt.Sprintf("paged-%d", store.MaxPendingIntentsPage+1)] {
		t.Error("the intent on the SECOND page was never read: recovery stopped after one page")
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
