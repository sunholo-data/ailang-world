package transitionreg

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// putSource stores one transition-source object and returns its ref, so
// descriptors pin a source the store really holds (verifySources' honest case).
func putSource(t *testing.T, s *store.Store, payload []byte) hashref.HashRef {
	t.Helper()
	obj := store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("test/transition-source")),
		SemanticID: "test/transition-source", Provenance: "publish_test", Payload: payload,
	}
	if err := s.PutObject(obj); err != nil {
		t.Fatalf("put source: %v", err)
	}
	return obj.Hash
}

// storedSourceDescriptor builds a valid descriptor whose TransitionFn the
// store holds, whose interpreter is the store's ARCHIVED fake (so the epoch
// derivation and loadability checks pass), and whose epoch is the one the
// fixture bootstrapped.
func storedSourceDescriptor(t *testing.T, s *store.Store, interpreter hashref.HashRef, id string) Descriptor {
	t.Helper()
	d := validDescriptor()
	d.ID = id
	d.Title = "title-" + id
	d.Interpreter = interpreter
	d.SemanticsEpoch = 1
	d.TransitionFn = putSource(t, s, []byte("transition source for "+id))
	return d
}

// TestGenesisCannotUseBuildNextWithZeroExpectedHead is the F2 measurement:
// BuildNext over a synthetic revision 0 pins Parent to the hash of the ENCODED
// empty revision, and Publish's absent-head branch refuses any non-zero
// parent. This is WHY PublishSet's genesis path constructs revision 1
// directly. If this test ever fails because BuildNext/Publish gained a
// genesis-capable path, the direct-construction branch in buildNextRevision
// can be replaced by it.
func TestGenesisCannotUseBuildNextWithZeroExpectedHead(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	r := NewReader(s)
	d := storedSourceDescriptor(t, s, ref, "tools.echo")
	genesis := Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 0}
	next, err := BuildNext(genesis, []Change{{ID: d.ID, Descriptor: &d}})
	if err != nil {
		t.Fatalf("BuildNext over synthetic revision 0: %v", err)
	}
	if next.Parent.IsZero() {
		t.Fatal("measurement premise wrong: BuildNext set a zero parent over an absent head")
	}
	_, err = r.Publish(context.Background(), hashref.HashRef{}, next)
	if err == nil || !strings.Contains(err.Error(), "parent is not captured head (absent)") {
		t.Fatalf("Publish(zero, BuildNext(genesis)) error = %v, want the absent-head parent refusal", err)
	}
	// Control: the same descriptor set publishes fine through the direct
	// genesis construction (what PublishSet does).
	res, err := NewPublisher(s, arch).PublishSet(context.Background(), []Change{{ID: d.ID, Descriptor: &d}})
	if err != nil || res.Revision != 1 || res.Unchanged {
		t.Fatalf("PublishSet genesis = (%+v, %v), want revision 1 changed", res, err)
	}
}

func TestPublishSetGenesisThenGrowth(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	r := NewPublisher(s, arch)
	ctx := context.Background()

	first := storedSourceDescriptor(t, s, ref, "tools.echo")
	res, err := r.PublishSet(ctx, []Change{{ID: first.ID, Descriptor: &first}})
	if err != nil {
		t.Fatalf("genesis publish: %v", err)
	}
	if res.Revision != 1 || res.Unchanged {
		t.Fatalf("genesis result = %+v, want revision 1 changed", res)
	}
	snap, err := r.ReadSnapshot(ctx)
	if err != nil {
		t.Fatalf("read after genesis: %v", err)
	}
	if snap.Revision != 1 || len(snap.List()) != 1 || snap.List()[0].ID != "tools.echo" {
		t.Fatalf("snapshot after genesis = revision %d entries %+v", snap.Revision, snap.List())
	}

	second := storedSourceDescriptor(t, s, ref, "worlds.merge")
	res, err = r.PublishSet(ctx, []Change{{ID: second.ID, Descriptor: &second}})
	if err != nil {
		t.Fatalf("growth publish: %v", err)
	}
	if res.Revision != 2 {
		t.Fatalf("growth revision = %d, want 2 (BuildNext over the captured head)", res.Revision)
	}
	snap, err = r.ReadSnapshot(ctx)
	if err != nil {
		t.Fatalf("read after growth: %v", err)
	}
	if got := snap.List(); len(got) != 2 || got[0].ID != "tools.echo" || got[1].ID != "worlds.merge" {
		t.Fatalf("growth entries = %+v, want both IDs in bytewise order", got)
	}
}

func TestPublishSetIdempotentRepublishIsNoOp(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	r := NewPublisher(s, arch)
	ctx := context.Background()
	d := storedSourceDescriptor(t, s, ref, "tools.echo")
	if _, err := r.PublishSet(ctx, []Change{{ID: d.ID, Descriptor: &d}}); err != nil {
		t.Fatalf("first publish: %v", err)
	}
	headBefore, _, ok, err := r.CurrentRevision(ctx)
	if err != nil || !ok {
		t.Fatalf("read head after first publish: ok=%v err=%v", ok, err)
	}

	res, err := r.PublishSet(ctx, []Change{{ID: d.ID, Descriptor: &d}})
	if err != nil {
		t.Fatalf("republish: %v", err)
	}
	if !res.Unchanged || res.Revision != 1 || res.Head != headBefore {
		t.Fatalf("republish = %+v, want Unchanged no-op at revision 1 with the same head", res)
	}
	snap, err := r.ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Revision != 1 || snap.Head != headBefore {
		t.Fatalf("head moved on identical republish: revision %d head %q", snap.Revision, snap.Head)
	}
}

func TestPublishSetRefusesAbsentTransitionSource(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	d := validDescriptor() // pins testFn, which no object backs
	d.Interpreter = ref
	d.SemanticsEpoch = 1
	_, err := NewPublisher(s, arch).PublishSet(context.Background(), []Change{{ID: d.ID, Descriptor: &d}})
	var absent *TransitionSourceAbsentError
	if !errors.As(err, &absent) {
		t.Fatalf("absent source error = %v, want *TransitionSourceAbsentError", err)
	}
	if _, _, ok, _ := NewReader(s).CurrentRevision(context.Background()); ok {
		t.Fatal("a refused publish must leave no head")
	}
}

func TestPublishSetRefusesEmptyAndRemovalOnlyGenesis(t *testing.T) {
	s, arch, _ := publisherStore(t, testFakeRelease, testFakeRelease)
	r := NewPublisher(s, arch)
	if _, err := r.PublishSet(context.Background(), nil); !errors.As(err, new(*EmptyGenesisError)) {
		t.Fatalf("empty genesis error = %v, want *EmptyGenesisError", err)
	}
	if _, err := r.PublishSet(context.Background(), []Change{{ID: "tools.echo"}}); err == nil {
		t.Fatal("removal-only genesis must be refused")
	}
}

// racingStore delegates to the real store but, on its FIRST head CAS, first
// installs a competing revision (another publisher won the race) and then
// returns the conflict the real CAS would have returned. With always set, it
// synthesizes a conflict on EVERY call and never delegates — the second
// publisher always loses.
type racingStore struct {
	*store.Store
	casCalls int
	always   bool
	compete  func() error
}

func (r *racingStore) CompareAndSetRegistryHead(name string, expected, next hashref.HashRef) error {
	r.casCalls++
	if r.always {
		head, ok, err := r.Store.GetRegistryHead(context.Background(), name)
		if err != nil {
			return err
		}
		return &store.RegistryCASConflict{Name: name, Expected: expected, Actual: head, HadHead: ok}
	}
	if r.casCalls == 1 && r.compete != nil {
		if err := r.compete(); err != nil {
			return err
		}
	}
	return r.Store.CompareAndSetRegistryHead(name, expected, next)
}

func TestPublishSetCASRetryMergesWinner(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	ctx := context.Background()

	// A competing publisher's descriptor set (different ID, also honest).
	competitor := storedSourceDescriptor(t, s, ref, "worlds.merge")
	compete := func() error {
		c := NewPublisher(s, arch)
		if _, err := c.PublishSet(ctx, []Change{{ID: competitor.ID, Descriptor: &competitor}}); err != nil {
			return err
		}
		return nil
	}
	racing := &racingStore{Store: s, compete: compete}
	r := NewPublisher(racing, arch)

	ours := storedSourceDescriptor(t, s, ref, "tools.echo")
	res, err := r.PublishSet(ctx, []Change{{ID: ours.ID, Descriptor: &ours}})
	if err != nil {
		t.Fatalf("publish under race: %v", err)
	}
	if res.Revision != 2 {
		t.Fatalf("revision after merge = %d, want 2 (retry rebuilds over the winner)", res.Revision)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := snap.List()
	if len(got) != 2 || got[0].ID != "tools.echo" || got[1].ID != "worlds.merge" {
		t.Fatalf("merged entries = %+v, want BOTH the winner's and ours (no lost update)", got)
	}
	if racing.casCalls != 2 {
		t.Fatalf("CAS calls = %d, want exactly 2 (one conflict, one win; no retry storm)", racing.casCalls)
	}
}

func TestPublishSetSecondConflictSurfacesTypedError(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	ctx := context.Background()
	ours := storedSourceDescriptor(t, s, ref, "tools.echo")
	alwaysLoses := &racingStore{Store: s, always: true, compete: nil}
	r := NewPublisher(alwaysLoses, arch)
	_, err := r.PublishSet(ctx, []Change{{ID: ours.ID, Descriptor: &ours}})
	if err == nil {
		t.Fatal("second conflict must surface an error, not spin")
	}
	if !store.IsRegistryCASConflict(err) {
		t.Fatalf("surfaced error = %v, want a detectable registry CAS conflict", err)
	}
	if alwaysLoses.casCalls != 2 {
		t.Fatalf("CAS calls = %d, want exactly 2 (bounded: one retry, then a typed failure)", alwaysLoses.casCalls)
	}
	if _, _, ok, _ := NewReader(s).CurrentRevision(ctx); ok {
		t.Fatal("a failed publish must not have moved the head")
	}
}

// TestPublishSetSameIDConflictOnCASRetryRefuses is kimi's round-1 secondary
// pinned exactly: the CAS-retry winner and this publisher both set the SAME
// ID with DIFFERENT descriptor bytes. The retry must REFUSE with the typed
// SameIDConflictError naming the ID — nothing is written, the head stays at
// the winner's revision with the WINNER's bytes (BuildNext replaces by ID,
// so without this refusal the loser would silently clobber them).
func TestPublishSetSameIDConflictOnCASRetryRefuses(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	ctx := context.Background()

	winner := storedSourceDescriptor(t, s, ref, "tools.echo")
	winner.Title = "title-winner"
	compete := func() error {
		c := NewPublisher(s, arch)
		if _, err := c.PublishSet(ctx, []Change{{ID: winner.ID, Descriptor: &winner}}); err != nil {
			return err
		}
		return nil
	}
	racing := &racingStore{Store: s, compete: compete}

	ours := storedSourceDescriptor(t, s, ref, "tools.echo")
	ours.Title = "title-loser" // same ID, different bytes
	_, err := NewPublisher(racing, arch).PublishSet(ctx, []Change{{ID: ours.ID, Descriptor: &ours}})
	var conflict *SameIDConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("same-ID race error = %v, want *SameIDConflictError", err)
	}
	if conflict.ID != "tools.echo" || conflict.Revision != 1 {
		t.Fatalf("conflict error = %+v, want the ID and the winner's revision named", conflict)
	}
	if racing.casCalls != 1 {
		t.Fatalf("CAS calls = %d, want exactly 1 (the refusal happens before a second write)", racing.casCalls)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := snap.List()
	if len(got) != 1 || got[0].ID != "tools.echo" || got[0].Title != "title-winner" {
		t.Fatalf("entries after the refused same-ID race = %+v, want the WINNER's bytes intact", got)
	}
}

// TestPublishSetCASRetrySameIDIdenticalBytesIsNoOp is the same-ID merge's
// non-conflicting half: when the winner published EXACTLY the bytes we were
// about to write, the retry converges to Unchanged at the winner's revision —
// the entry is neither dropped nor rewritten.
func TestPublishSetCASRetrySameIDIdenticalBytesIsNoOp(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	ctx := context.Background()

	ours := storedSourceDescriptor(t, s, ref, "tools.echo")
	compete := func() error {
		c := NewPublisher(s, arch)
		if _, err := c.PublishSet(ctx, []Change{{ID: ours.ID, Descriptor: &ours}}); err != nil {
			return err
		}
		return nil
	}
	racing := &racingStore{Store: s, compete: compete}
	res, err := NewPublisher(racing, arch).PublishSet(ctx, []Change{{ID: ours.ID, Descriptor: &ours}})
	if err != nil {
		t.Fatalf("identical same-ID race: %v", err)
	}
	if !res.Unchanged || res.Revision != 1 {
		t.Fatalf("identical same-ID race result = %+v, want Unchanged at the winner's revision 1", res)
	}
	if racing.casCalls != 1 {
		t.Fatalf("CAS calls = %d, want exactly 1 (the convergence writes nothing)", racing.casCalls)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snap.List(); len(got) != 1 || got[0].ID != "tools.echo" || got[0].Title != ours.Title {
		t.Fatalf("entries after the identical same-ID race = %+v, want the winner's entry intact", got)
	}
}

func TestCanonicalSchemaExportMatchesValidate(t *testing.T) {
	raw := []byte(`{"z":1.0,"a":"x"}`)
	got, err := CanonicalSchema(raw)
	if err != nil {
		t.Fatal(err)
	}
	d := validDescriptor()
	d.InputSchema = got
	if err := d.Validate(); err != nil {
		t.Fatalf("CanonicalSchema output rejected by Validate: %v", err)
	}
	if string(got) != `{"a":"x","z":1}` {
		t.Fatalf("CanonicalSchema = %s, want the codec's canonical form", got)
	}
}
