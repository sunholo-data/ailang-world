package store

import (
	"errors"
	"sync"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// openMem opens a fresh in-memory store for a test and registers cleanup.
func openMem(t *testing.T) *Store {
	t.Helper()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// obj builds an Object whose Hash correctly addresses its payload, so it passes
// content verification. interfaceHash/semanticId/provenance are labels.
func obj(payload string, semanticID string) Object {
	p := []byte(payload)
	return Object{
		Hash:          hashref.SumSHA256(p),
		InterfaceHash: hashref.SumSHA256([]byte("iface:" + semanticID)),
		SemanticID:    semanticID,
		Provenance:    "test",
		Payload:       p,
	}
}

func TestObjectPersistenceRoundTrip(t *testing.T) {
	s := openMem(t)
	o := obj("hello world payload", "state/v1")

	if err := s.PutObject(o); err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	got, ok, err := s.GetObject(o.Hash)
	if err != nil || !ok {
		t.Fatalf("GetObject: ok=%v err=%v", ok, err)
	}
	if got.Hash.String() != o.Hash.String() ||
		got.InterfaceHash.String() != o.InterfaceHash.String() ||
		got.SemanticID != o.SemanticID ||
		got.Provenance != o.Provenance ||
		string(got.Payload) != string(o.Payload) {
		t.Fatalf("object round-trip mismatch:\n got=%+v\nwant=%+v", got, o)
	}

	// Re-inserting the identical object is idempotent.
	if err := s.PutObject(o); err != nil {
		t.Fatalf("PutObject (idempotent): %v", err)
	}
}

func TestObjectContentVerificationRejectsMismatch(t *testing.T) {
	s := openMem(t)
	bad := obj("real payload", "state/v1")
	// Corrupt the payload so Hash no longer addresses it.
	bad.Payload = []byte("tampered payload")

	if err := s.PutObject(bad); err == nil {
		t.Fatal("PutObject accepted an object whose hash does not match its payload")
	}
}

func TestWorldPersistenceRoundTrip(t *testing.T) {
	s := openMem(t)
	w := World{
		Ref:       hashref.SumSHA256([]byte("world-rev-1")),
		Revision:  1,
		StateRoot: hashref.SumSHA256([]byte("state-1")),
		LogHead:   hashref.SumSHA256([]byte("log-head-1")),
	}
	if err := s.PutWorld(w); err != nil {
		t.Fatalf("PutWorld: %v", err)
	}
	got, ok, err := s.GetWorld(w.Ref)
	if err != nil || !ok {
		t.Fatalf("GetWorld: ok=%v err=%v", ok, err)
	}
	if got != w {
		t.Fatalf("world round-trip mismatch:\n got=%+v\nwant=%+v", got, w)
	}
}

// TestFrozenHeaderRoundTrip stores a log entry through a commit and reads it back,
// asserting every one of the six frozen header fields plus the separate
// transition-body reference come back byte-identical.
func TestFrozenHeaderRoundTrip(t *testing.T) {
	s := openMem(t)

	genesis := seedGenesis(t, s)

	entryHeader := LogHeader{
		EntryIndex:     1,
		SemanticsEpoch: 1,
		TransitionFn:   hashref.SumSHA256([]byte("transition-fn-canonical-bytes")),
		Interpreter:    hashref.SumSHA256([]byte("interpreter-bytes")),
		PrevEntryHash:  genesis.LogHead,
		WrittenBy:      "sha256:deadbeef ailang v0.30.0",
	}
	transitionBody := obj("transition body payload", "transition/body")
	entryHash := hashref.SumSHA256([]byte("encoded-header-plus-body-ref"))

	next := World{
		Ref:       hashref.SumSHA256([]byte("world-rev-2")),
		Revision:  genesis.Revision + 1,
		StateRoot: hashref.SumSHA256([]byte("state-2")),
		LogHead:   entryHash,
	}

	if err := s.Commit(Commit{
		ObservedHead: genesis.Ref,
		Objects:      []Object{transitionBody},
		NextWorld:    next,
		Entry: LogEntry{
			Header:        entryHeader,
			EntryHash:     entryHash,
			TransitionRef: transitionBody.Hash,
		},
	}); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	got, ok, err := s.GetLogEntry(1)
	if err != nil || !ok {
		t.Fatalf("GetLogEntry(1): ok=%v err=%v", ok, err)
	}
	if got.Header != entryHeader {
		t.Fatalf("frozen header not stored verbatim:\n got=%+v\nwant=%+v", got.Header, entryHeader)
	}
	if got.EntryHash.String() != entryHash.String() {
		t.Fatalf("entry hash mismatch: got %q want %q", got.EntryHash, entryHash)
	}
	if got.TransitionRef.String() != transitionBody.Hash.String() {
		t.Fatalf("transition ref mismatch: got %q want %q", got.TransitionRef, transitionBody.Hash)
	}

	// The selected head must have advanced to the new world.
	sel, ok, err := s.SelectedHead()
	if err != nil || !ok {
		t.Fatalf("SelectedHead: ok=%v err=%v", ok, err)
	}
	if sel.String() != next.Ref.String() {
		t.Fatalf("selected head not advanced: got %q want %q", sel, next.Ref)
	}
}

// seedGenesis creates and selects a genesis world head so subsequent commits
// have a non-nil observed head to compare against. Returns the genesis world.
func seedGenesis(t *testing.T, s *Store) World {
	t.Helper()
	g := World{
		Ref:       hashref.SumSHA256([]byte("world-genesis")),
		Revision:  0,
		StateRoot: hashref.SumSHA256([]byte("state-genesis")),
		LogHead:   hashref.SumSHA256([]byte("log-genesis")),
	}
	if err := s.PutWorld(g); err != nil {
		t.Fatalf("seed PutWorld: %v", err)
	}
	if err := s.SelectHead(g.Ref); err != nil {
		t.Fatalf("seed SelectHead: %v", err)
	}
	return g
}

// TestCommitSingleTransactionConflictOnStaleHead proves the compare-and-append:
// a caller planning against the genesis head succeeds; a second caller who also
// planned against the (now stale) genesis head gets a structured ConflictError,
// and the store is left untouched by the failed commit.
func TestCommitConflictOnStaleHead(t *testing.T) {
	s := openMem(t)
	genesis := seedGenesis(t, s)

	body := obj("body-a", "t/a")
	entryHash := hashref.SumSHA256([]byte("entry-1"))
	world2 := World{
		Ref:       hashref.SumSHA256([]byte("world-2")),
		Revision:  1,
		StateRoot: hashref.SumSHA256([]byte("state-2")),
		LogHead:   entryHash,
	}
	commitA := Commit{
		ObservedHead: genesis.Ref,
		Objects:      []Object{body},
		NextWorld:    world2,
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex:    1,
				TransitionFn:  hashref.SumSHA256([]byte("fn-a")),
				Interpreter:   hashref.SumSHA256([]byte("interp")),
				PrevEntryHash: genesis.LogHead,
				WrittenBy:     "writer-a",
			},
			EntryHash:     entryHash,
			TransitionRef: body.Hash,
		},
	}
	if err := s.Commit(commitA); err != nil {
		t.Fatalf("first Commit: %v", err)
	}

	// Second caller still observes the genesis head — now stale.
	body2 := obj("body-b", "t/b")
	commitB := Commit{
		ObservedHead: genesis.Ref, // stale
		Objects:      []Object{body2},
		NextWorld: World{
			Ref:       hashref.SumSHA256([]byte("world-3")),
			Revision:  2,
			StateRoot: hashref.SumSHA256([]byte("state-3")),
			LogHead:   hashref.SumSHA256([]byte("entry-2")),
		},
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex:    2,
				TransitionFn:  hashref.SumSHA256([]byte("fn-b")),
				Interpreter:   hashref.SumSHA256([]byte("interp")),
				PrevEntryHash: hashref.SumSHA256([]byte("entry-2")),
				WrittenBy:     "writer-b",
			},
			EntryHash:     hashref.SumSHA256([]byte("entry-2")),
			TransitionRef: body2.Hash,
		},
	}
	err := s.Commit(commitB)
	if err == nil {
		t.Fatal("stale-head Commit succeeded; expected ConflictError")
	}
	if !IsConflict(err) {
		t.Fatalf("expected *ConflictError, got %T: %v", err, err)
	}
	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("errors.As failed for ConflictError")
	}
	if ce.ObservedHead.String() != genesis.Ref.String() {
		t.Fatalf("ConflictError.ObservedHead = %q, want %q", ce.ObservedHead, genesis.Ref)
	}
	if ce.SelectedHead.String() != world2.Ref.String() {
		t.Fatalf("ConflictError.SelectedHead = %q, want %q", ce.SelectedHead, world2.Ref)
	}

	// The failed commit must have written nothing: log entry 2 absent, and the
	// selected head unchanged from commitA's world.
	if _, ok, _ := s.GetLogEntry(2); ok {
		t.Fatal("failed commit left log entry 2 behind; transaction did not roll back")
	}
	sel, _, _ := s.SelectedHead()
	if sel.String() != world2.Ref.String() {
		t.Fatalf("selected head changed after failed commit: got %q want %q", sel, world2.Ref)
	}
}

func TestRegistryHeadRoundTrip(t *testing.T) {
	s := openMem(t)
	reg := hashref.SumSHA256([]byte("epoch-1-registry-object"))
	if err := s.SetRegistryHead(EpochRegistryV1, reg); err != nil {
		t.Fatalf("SetRegistryHead: %v", err)
	}
	got, ok, err := s.GetRegistryHead(EpochRegistryV1)
	if err != nil || !ok {
		t.Fatalf("GetRegistryHead: ok=%v err=%v", ok, err)
	}
	if got.String() != reg.String() {
		t.Fatalf("registry head mismatch: got %q want %q", got, reg)
	}

	// Updating the head replaces it in place (one row per registry name).
	reg2 := hashref.SumSHA256([]byte("epoch-2-registry-object"))
	if err := s.SetRegistryHead(EpochRegistryV1, reg2); err != nil {
		t.Fatalf("SetRegistryHead update: %v", err)
	}
	got2, _, _ := s.GetRegistryHead(EpochRegistryV1)
	if got2.String() != reg2.String() {
		t.Fatalf("registry head not updated: got %q want %q", got2, reg2)
	}
}

func TestCompareAndSetRegistryHead(t *testing.T) {
	put := func(t *testing.T, s *Store, payload string) Object {
		t.Helper()
		o := obj(payload, TransitionRegistryV1)
		if err := s.PutObject(o); err != nil {
			t.Fatalf("PutObject(%q): %v", payload, err)
		}
		return o
	}
	assertHead := func(t *testing.T, s *Store, name string, want hashref.HashRef, wantOK bool) {
		t.Helper()
		got, ok, err := s.GetRegistryHead(name)
		if err != nil || ok != wantOK || (ok && got != want) {
			t.Fatalf("GetRegistryHead(%q) = (%q,%v,%v), want (%q,%v,nil)", name, got, ok, err, want, wantOK)
		}
	}

	t.Run("absent_head_accepts_zero_expected", func(t *testing.T) {
		s := openMem(t)
		next := put(t, s, "next")
		if err := s.CompareAndSetRegistryHead(TransitionRegistryV1, hashref.HashRef{}, next.Hash); err != nil {
			t.Fatalf("CompareAndSetRegistryHead: %v", err)
		}
		assertHead(t, s, TransitionRegistryV1, next.Hash, true)
	})

	t.Run("absent_head_rejects_nonzero_expected", func(t *testing.T) {
		s := openMem(t)
		expected := put(t, s, "expected")
		next := put(t, s, "next")
		err := s.CompareAndSetRegistryHead(TransitionRegistryV1, expected.Hash, next.Hash)
		var conflict *RegistryCASConflict
		if !errors.As(err, &conflict) || !IsRegistryCASConflict(err) || conflict.HadHead || conflict.Expected != expected.Hash || !conflict.Actual.IsZero() {
			t.Fatalf("conflict = %#v, err=%v", conflict, err)
		}
		assertHead(t, s, TransitionRegistryV1, hashref.HashRef{}, false)
	})

	t.Run("stale_expected_conflicts", func(t *testing.T) {
		s := openMem(t)
		actual := put(t, s, "actual")
		expected := put(t, s, "expected")
		next := put(t, s, "next")
		if err := s.SetRegistryHead(TransitionRegistryV1, actual.Hash); err != nil {
			t.Fatal(err)
		}
		err := s.CompareAndSetRegistryHead(TransitionRegistryV1, expected.Hash, next.Hash)
		var conflict *RegistryCASConflict
		if !errors.As(err, &conflict) || !conflict.HadHead || conflict.Expected != expected.Hash || conflict.Actual != actual.Hash {
			t.Fatalf("conflict = %#v, err=%v", conflict, err)
		}
		assertHead(t, s, TransitionRegistryV1, actual.Hash, true)
	})

	t.Run("dangling_next_refused", func(t *testing.T) {
		s := openMem(t)
		actual := put(t, s, "actual")
		if err := s.SetRegistryHead(TransitionRegistryV1, actual.Hash); err != nil {
			t.Fatal(err)
		}
		dangling := hashref.SumSHA256([]byte("absent"))
		if err := s.CompareAndSetRegistryHead(TransitionRegistryV1, actual.Hash, dangling); err == nil {
			t.Fatal("CompareAndSetRegistryHead accepted a dangling next object")
		}
		assertHead(t, s, TransitionRegistryV1, actual.Hash, true)
	})

	t.Run("rollback_leaves_head_unchanged", func(t *testing.T) {
		s := openMem(t)
		actual := put(t, s, "actual")
		stale := put(t, s, "stale")
		next := put(t, s, "next")
		if err := s.SetRegistryHead(TransitionRegistryV1, actual.Hash); err != nil {
			t.Fatal(err)
		}
		if err := s.CompareAndSetRegistryHead(TransitionRegistryV1, stale.Hash, next.Hash); err == nil {
			t.Fatal("stale CAS succeeded")
		}
		assertHead(t, s, TransitionRegistryV1, actual.Hash, true)
		if err := s.CompareAndSetRegistryHead(TransitionRegistryV1, actual.Hash, hashref.SumSHA256([]byte("missing"))); err == nil {
			t.Fatal("dangling CAS succeeded")
		}
		assertHead(t, s, TransitionRegistryV1, actual.Hash, true)
	})

	t.Run("epoch_registry_isolation", func(t *testing.T) {
		s := openMem(t)
		epoch := put(t, s, "epoch")
		transition := put(t, s, "transition")
		epochNext := put(t, s, "epoch-next")
		if err := s.SetRegistryHead(EpochRegistryV1, epoch.Hash); err != nil {
			t.Fatal(err)
		}
		if err := s.CompareAndSetRegistryHead(TransitionRegistryV1, hashref.HashRef{}, transition.Hash); err != nil {
			t.Fatal(err)
		}
		assertHead(t, s, EpochRegistryV1, epoch.Hash, true)
		assertHead(t, s, TransitionRegistryV1, transition.Hash, true)
		if err := s.CompareAndSetRegistryHead(EpochRegistryV1, epoch.Hash, epochNext.Hash); err != nil {
			t.Fatal(err)
		}
		assertHead(t, s, EpochRegistryV1, epochNext.Hash, true)
		assertHead(t, s, TransitionRegistryV1, transition.Hash, true)
	})

	t.Run("concurrent_racers_one_winner", func(t *testing.T) {
		s := openMem(t)
		const racers = 8
		next := make([]Object, racers)
		for i := range next {
			next[i] = put(t, s, string(rune('a'+i)))
		}
		start := make(chan struct{})
		errs := make(chan error, racers)
		var wg sync.WaitGroup
		for i := range next {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				errs <- s.CompareAndSetRegistryHead(TransitionRegistryV1, hashref.HashRef{}, next[i].Hash)
			}(i)
		}
		close(start)
		wg.Wait()
		close(errs)
		wins, conflicts := 0, 0
		for err := range errs {
			if err == nil {
				wins++
			} else if IsRegistryCASConflict(err) {
				conflicts++
			} else {
				t.Fatalf("unexpected racer error: %v", err)
			}
		}
		if wins != 1 || conflicts != racers-1 {
			t.Fatalf("wins=%d conflicts=%d, want 1/%d", wins, conflicts, racers-1)
		}
		head, ok, err := s.GetRegistryHead(TransitionRegistryV1)
		if err != nil || !ok {
			t.Fatalf("final head: ok=%v err=%v", ok, err)
		}
		found := false
		for _, candidate := range next {
			if candidate.Hash == head {
				found = true
			}
		}
		if !found {
			t.Fatalf("final head %q was not a racer", head)
		}
	})
}

// TestVerificationCacheKeyIsExactlyThePair proves the cache key is
// (transition_fn_ref, interpreter_ref) EXCLUSIVELY:
//   - a different transitionFn is a distinct row (cache miss on lookup),
//   - a different interpreter is a distinct row (cache miss),
//   - an epoch-ONLY change on the same pair updates the SAME row in place
//     (metadata-compatible), leaving exactly one row for that pair.
func TestVerificationCacheKeyIsExactlyThePair(t *testing.T) {
	s := openMem(t)

	fn := hashref.SumSHA256([]byte("transition-fn"))
	interp := hashref.SumSHA256([]byte("interpreter"))
	other := hashref.SumSHA256([]byte("other"))

	base := VerifyResult{
		TransitionFn:   fn,
		Interpreter:    interp,
		SemanticsEpoch: 1,
		Verified:       true,
		Detail:         "typecheck ok @ epoch 1",
	}
	if err := s.PutVerifyResult(base); err != nil {
		t.Fatalf("PutVerifyResult base: %v", err)
	}

	// Exact pair hits.
	got, ok, err := s.GetVerifyResult(fn, interp)
	if err != nil || !ok {
		t.Fatalf("GetVerifyResult exact pair: ok=%v err=%v", ok, err)
	}
	if !got.Verified || got.Detail != base.Detail || got.SemanticsEpoch != 1 {
		t.Fatalf("cached row mismatch: %+v", got)
	}

	// Different transitionFn — must miss (distinct key).
	if _, ok, _ := s.GetVerifyResult(other, interp); ok {
		t.Fatal("different transitionFn hit the cache; key is not the exact pair")
	}
	// Different interpreter — must miss (distinct key).
	if _, ok, _ := s.GetVerifyResult(fn, other); ok {
		t.Fatal("different interpreter hit the cache; key is not the exact pair")
	}

	// Epoch-only change on the SAME pair: must update the same row, not add one.
	epochChanged := base
	epochChanged.SemanticsEpoch = 2
	epochChanged.Detail = "typecheck ok @ epoch 2"
	if err := s.PutVerifyResult(epochChanged); err != nil {
		t.Fatalf("PutVerifyResult epoch-only change: %v", err)
	}

	// Still exactly one row for the pair, and it is metadata-updated in place.
	if n := s.countCacheRows(t, fn, interp); n != 1 {
		t.Fatalf("epoch-only change produced %d rows for the pair; want exactly 1", n)
	}
	after, ok, _ := s.GetVerifyResult(fn, interp)
	if !ok {
		t.Fatal("pair lookup miss after epoch-only change; the selected row was not preserved")
	}
	if after.SemanticsEpoch != 2 || after.Detail != "typecheck ok @ epoch 2" {
		t.Fatalf("epoch metadata not carried into the same row: %+v", after)
	}
}

// countCacheRows counts verification_cache rows for a given pair via the
// exported handle's underlying DB (test-only introspection).
func (s *Store) countCacheRows(t *testing.T, fn, interp hashref.HashRef) int {
	t.Helper()
	var n int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM verification_cache
		  WHERE transition_fn_ref = ? AND interpreter_ref = ?;`,
		fn.String(), interp.String(),
	).Scan(&n)
	if err != nil {
		t.Fatalf("count cache rows: %v", err)
	}
	return n
}
