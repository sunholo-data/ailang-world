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

func TestCrashWindowRecoveryReportsAndResumptionMintsFreshOrdinal(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	base, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	const episodeID = "crash-window-resume"
	dispatches := 0
	failing := &failRecordStore{base: base}
	crashHandler := HandlerFunc(func(
		_ context.Context,
		_ EffectRequest,
		payload []byte,
	) ([]byte, error) {
		dispatches++
		receipt, hasIntent, err := base.GetEffectReceipt(
			store.EffectInvocationID(episodeID, 0),
		)
		if err != nil || !hasIntent || receipt.State != store.ReceiptIndeterminate {
			t.Fatalf("receipt at crash-window dispatch = %#v, hasIntent %v, err %v; want indeterminate",
				receipt, hasIntent, err)
		}
		return append([]byte("echo:"), payload...), nil
	})
	session := newSession(
		failing,
		episodeID,
		[]Capability{{"probe", "/ok", 100, 10}},
		Registry{"probe": crashHandler},
		Live,
		nil,
	)
	req := EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 23}
	result, recordRef, invokeErr := session.Invoke(
		context.Background(), req, []byte("before-crash"),
	)
	if invokeErr == nil || result != nil || !recordRef.IsZero() {
		t.Fatalf("faulted Invoke = result %q, ref %s, err %v; want nil, zero, error",
			result, recordRef, invokeErr)
	}
	if dispatches != 1 {
		t.Fatalf("pre-crash dispatches = %d, want 1", dispatches)
	}
	if failing.effectOutcomeCalls != 0 {
		t.Fatalf("pre-crash outcome calls = %d, want 0", failing.effectOutcomeCalls)
	}
	if _, ok, err := base.GetObject(requestHash(req, []byte("before-crash"))); err != nil || !ok {
		t.Fatalf("durable request object = ok %v, err %v; want present", ok, err)
	}
	if err := base.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reopened.Close() }()
	firstID := store.EffectInvocationID(episodeID, 0)
	receipt, hasIntent, err := reopened.GetEffectReceipt(firstID)
	if err != nil || !hasIntent || receipt.State != store.ReceiptIndeterminate {
		t.Fatalf("reopened receipt = %#v, hasIntent %v, err %v; want indeterminate",
			receipt, hasIntent, err)
	}

	recoveryProbe := &recoveryCountingProbe{}
	findings, err := Recover(reopened, Registry{"probe": recoveryProbe})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("recovery findings = %d, want 1", len(findings))
	}
	finding := findings[0]
	if finding.Err.EpisodeID != episodeID || finding.Err.Ordinal != 0 ||
		finding.Err.Effect != "probe" || finding.Err.Scope != "/ok" ||
		finding.EffectIntent.InvocationID != firstID {
		t.Fatalf("effect finding = %#v", finding)
	}
	if recoveryProbe.dispatches != 0 {
		t.Fatalf("recovery dispatches = %d, want 0", recoveryProbe.dispatches)
	}
	receipt, hasIntent, err = reopened.GetEffectReceipt(firstID)
	if err != nil || !hasIntent || receipt.State != store.ReceiptIndeterminate {
		t.Fatalf("post-recovery receipt = %#v, hasIntent %v, err %v; recovery acted",
			receipt, hasIntent, err)
	}

	resumedDispatches := 0
	resumed := NewSession(
		reopened,
		episodeID,
		[]Capability{{"probe", "/ok", 100, 10}},
		echoRegistry(&resumedDispatches),
	)
	result, recordRef, err = resumed.Invoke(
		context.Background(),
		EffectRequest{Effect: "probe", Scope: "/ok", Cost: 1, Now: 24},
		[]byte("after-crash"),
	)
	if err != nil {
		t.Fatalf("resumed Invoke: %v", err)
	}
	if string(result) != "echo:after-crash" || recordRef.IsZero() || resumedDispatches != 1 {
		t.Fatalf("resumed Invoke = result %q, ref %s, dispatches %d",
			result, recordRef, resumedDispatches)
	}
	secondID := store.EffectInvocationID(episodeID, 1)
	if secondID <= firstID {
		t.Fatalf("resumed ID %q is not strictly past %q", secondID, firstID)
	}
	resumedReceipt, hasIntent, err := reopened.GetEffectReceipt(secondID)
	if err != nil || !hasIntent || resumedReceipt.State != store.ReceiptResolved ||
		resumedReceipt.EffectOutcome == nil || resumedReceipt.EffectOutcome.RecordRef != recordRef {
		t.Fatalf("resumed receipt = %#v, hasIntent %v, err %v; want resolved",
			resumedReceipt, hasIntent, err)
	}
}

func TestRecoverCountingProbeDispatchesZeroHandlers(t *testing.T) {
	s, _ := pendingRecoveryCommit(t, "counting-probe")
	effectID, _, err := s.AppendNextEffectIntent("counting-probe-episode", store.EffectIntent{
		EpisodeID:   "counting-probe-episode",
		Effect:      recoveryProbeEffect,
		Scope:       "/pending",
		Cost:        1,
		RequestRef:  hashref.SumSHA256([]byte("counting-probe-effect-request")),
		LogicalTime: 1,
	})
	if err != nil {
		t.Fatalf("AppendNextEffectIntent: %v", err)
	}
	probe := &recoveryCountingProbe{}
	findings, err := Recover(s, Registry{recoveryProbeEffect: probe})
	if err != nil {
		t.Fatalf("Recover: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("Recover findings=%d, want 2 (commit and effect)", len(findings))
	}
	if findings[1].EffectIntent.InvocationID != effectID {
		t.Fatalf("effect finding = %#v, want invocation %q", findings[1], effectID)
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
	pages           [][]store.PendingIntent
	fromIndex       []int64 // the cursor received on each call; -1 means "no cursor"
	calls           int
	maxCalls        int
	sawReceipt      map[string]bool
	effectPages     [][]store.PendingEffectIntent
	effectFromIndex []int64
	effectCalls     int
	sawEffect       map[string]bool
}

// neverDrainingRecoveryStore reuses fixed page buffers because these tests
// intentionally exercise all 2^20 pages. Allocating a page on every call would
// turn a bounded-loop assertion into an allocator benchmark. Receipts are
// resolved so findings do not accumulate; the fake's tripwire makes an
// unbounded mutant fail deterministically instead of hanging.
type neverDrainingRecoveryStore struct {
	commitPage  []store.PendingIntent
	effectPage  []store.PendingEffectIntent
	commitCalls int
	effectCalls int
	maxCalls    int
}

func newNeverDrainingRecoveryStore(commit, effect bool) *neverDrainingRecoveryStore {
	p := &neverDrainingRecoveryStore{maxCalls: maxRecoveryPages + 16}
	if commit {
		p.commitPage = make([]store.PendingIntent, store.MaxPendingIntentsPage)
		for i := range p.commitPage {
			p.commitPage[i].InvocationID = "never-draining-commit"
		}
	}
	if effect {
		p.effectPage = make([]store.PendingEffectIntent, store.MaxPendingIntentsPage)
		for i := range p.effectPage {
			p.effectPage[i].InvocationID = "effect:never-draining:0"
		}
	}
	return p
}

func (p *neverDrainingRecoveryStore) PendingIntents(
	_ int,
	fromIndex ...int64,
) ([]store.PendingIntent, error) {
	if len(p.commitPage) == 0 {
		return nil, nil
	}
	p.commitCalls++
	if p.commitCalls > p.maxCalls {
		return nil, fmt.Errorf("commit fake tripwire exceeded after %d calls", p.commitCalls)
	}
	cursor := int64(0)
	if len(fromIndex) > 0 {
		cursor = fromIndex[0]
	}
	for i := range p.commitPage {
		p.commitPage[i].Seq = cursor + int64(i) + 1
	}
	return p.commitPage, nil
}

func (p *neverDrainingRecoveryStore) GetReceipt(id string) (store.Receipt, bool, error) {
	return store.Receipt{InvocationID: id, State: store.ReceiptResolved}, true, nil
}

func (p *neverDrainingRecoveryStore) PendingEffectIntents(
	_ int,
	fromIndex ...int64,
) ([]store.PendingEffectIntent, error) {
	if len(p.effectPage) == 0 {
		return nil, nil
	}
	p.effectCalls++
	if p.effectCalls > p.maxCalls {
		return nil, fmt.Errorf("effect fake tripwire exceeded after %d calls", p.effectCalls)
	}
	cursor := int64(0)
	if len(fromIndex) > 0 {
		cursor = fromIndex[0]
	}
	for i := range p.effectPage {
		p.effectPage[i].Seq = cursor + int64(i) + 1
	}
	return p.effectPage, nil
}

func (p *neverDrainingRecoveryStore) GetEffectReceipt(
	id string,
) (store.Receipt, bool, error) {
	return store.Receipt{InvocationID: id, State: store.ReceiptResolved}, true, nil
}

func TestRecoverCommitStopsAtPageBound(t *testing.T) {
	probe := newNeverDrainingRecoveryStore(true, false)
	_, err := recoverCommitPending(probe)
	const want = "broker: recovery exceeded 1048576 pages"
	if err == nil || err.Error() != want {
		t.Fatalf("recoverCommitPending error = %v, want %q (calls=%d)",
			err, want, probe.commitCalls)
	}
	if probe.commitCalls != maxRecoveryPages {
		t.Fatalf("commit pending calls = %d, want %d", probe.commitCalls, maxRecoveryPages)
	}
}

func TestRecoverEffectStopsAtPageBound(t *testing.T) {
	probe := newNeverDrainingRecoveryStore(false, true)
	_, err := recoverEffectPending(probe, nil)
	const want = "broker: effect recovery exceeded 1048576 pages"
	if err == nil || err.Error() != want {
		t.Fatalf("recoverEffectPending error = %v, want %q (calls=%d)",
			err, want, probe.effectCalls)
	}
	if probe.effectCalls != maxRecoveryPages {
		t.Fatalf("effect pending calls = %d, want %d", probe.effectCalls, maxRecoveryPages)
	}
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

func (p *pagedRecoveryStore) PendingEffectIntents(
	limit int,
	fromIndex ...int64,
) ([]store.PendingEffectIntent, error) {
	if limit != store.MaxPendingIntentsPage {
		return nil, fmt.Errorf("effect limit=%d, want store.MaxPendingIntentsPage", limit)
	}
	cursor := int64(-1)
	if len(fromIndex) > 0 {
		cursor = fromIndex[0]
	}
	p.effectFromIndex = append(p.effectFromIndex, cursor)
	p.effectCalls++
	if p.effectCalls > p.maxCalls {
		return nil, fmt.Errorf("effect recovery did not terminate: %d calls, cursors=%v",
			p.effectCalls, p.effectFromIndex)
	}
	idx := 0
	for i, page := range p.effectPages {
		if len(page) == 0 {
			continue
		}
		if cursor < page[len(page)-1].Seq {
			idx = i
			break
		}
		idx = i + 1
	}
	if idx >= len(p.effectPages) {
		return nil, nil
	}
	return p.effectPages[idx], nil
}

func (p *pagedRecoveryStore) GetEffectReceipt(id string) (store.Receipt, bool, error) {
	if p.sawEffect == nil {
		p.sawEffect = map[string]bool{}
	}
	p.sawEffect[id] = true
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

func TestRecoverEffectPagesWithKeysetCursorAndKeepsFindingShapesSeparate(t *testing.T) {
	full := make([]store.PendingEffectIntent, store.MaxPendingIntentsPage)
	for i := range full {
		ordinal := int64(i)
		id := store.EffectInvocationID("effect-paged", ordinal)
		full[i] = store.PendingEffectIntent{
			Seq:          ordinal + 1,
			InvocationID: id,
			Intent: store.EffectIntent{
				InvocationID: id, EpisodeID: "effect-paged", Ordinal: ordinal,
				Effect: "probe", Scope: "/paged",
			},
		}
	}
	tailOrdinal := int64(store.MaxPendingIntentsPage)
	tailID := store.EffectInvocationID("effect-paged", tailOrdinal)
	tail := []store.PendingEffectIntent{{
		Seq:          tailOrdinal + 1,
		InvocationID: tailID,
		Intent: store.EffectIntent{
			InvocationID: tailID, EpisodeID: "effect-paged", Ordinal: tailOrdinal,
			Effect: "probe", Scope: "/paged",
		},
	}}
	probe := &pagedRecoveryStore{
		effectPages: [][]store.PendingEffectIntent{full, tail},
		maxCalls:    8,
	}
	findings, err := recoverPending(probe)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != store.MaxPendingIntentsPage+1 {
		t.Fatalf("effect findings = %d, want %d",
			len(findings), store.MaxPendingIntentsPage+1)
	}
	if len(probe.effectFromIndex) < 2 ||
		probe.effectFromIndex[0] != -1 ||
		probe.effectFromIndex[1] != int64(store.MaxPendingIntentsPage) {
		t.Fatalf("effect cursors = %v, want [-1 %d ...]",
			probe.effectFromIndex, store.MaxPendingIntentsPage)
	}
	last := findings[len(findings)-1]
	if !probe.sawEffect[tailID] || last.Intent.InvocationID != "" ||
		last.EffectIntent.InvocationID != tailID || last.Err.Ordinal != tailOrdinal {
		t.Fatalf("last effect finding = %#v, saw receipt %v", last, probe.sawEffect[tailID])
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
		{false, true, true},
		{true, false, false},
		{true, true, true},
	} {
		if got := retryAllowed(row.indeterminate, row.reconciled); got != row.want {
			t.Fatalf("retryAllowed(%v,%v)=%v, want %v",
				row.indeterminate, row.reconciled, got, row.want)
		}
	}
}
