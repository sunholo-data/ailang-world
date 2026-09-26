package transitionreg

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// TestEpochsForInterpreterRefusesAbsentAndUnnominated is M3's direct pin of
// the derivation (before PublishSet exists): an interpreter whose release no
// epoch nominates derives ZERO epochs even though epoch 1 exists (no
// "every epoch" / default-to-1 derivation — MUT-10), and a store with no
// epoch-registry head is the typed EpochRegistryAbsentError, never a
// synthetic epoch-1 registry (MUT-11's absent-head arm).
func TestEpochsForInterpreterRefusesAbsentAndUnnominated(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, "OTHER-RELEASE v1")
	epochs, release, err := NewPublisher(s, arch).EpochsForInterpreter(context.Background(), ref)
	if err != nil {
		t.Fatalf("EpochsForInterpreter: %v", err)
	}
	if release != testFakeRelease || len(epochs) != 0 {
		t.Fatalf("derived (%v, %q), want no epoch for release %q (epoch 1 nominates another release)", epochs, release, testFakeRelease)
	}

	absentStore, absentArch, absentRef := publisherStore(t, testFakeRelease, "")
	_, _, err = NewPublisher(absentStore, absentArch).EpochsForInterpreter(context.Background(), absentRef)
	if !errors.As(err, new(*EpochRegistryAbsentError)) {
		t.Fatalf("absent epoch registry = %v, want *EpochRegistryAbsentError (never a default)", err)
	}

	if _, _, err := NewReader(s).EpochsForInterpreter(context.Background(), ref); !errors.As(err, new(*PublisherArchiveRequiredError)) {
		t.Fatalf("bare-reader derivation = %v, want *PublisherArchiveRequiredError", err)
	}
}

// TestPublishSetEpochTwinInterpretersKeepDistinctPins is AC-EPOCH-TWIN
// (quorum r2, gpt6-astra; adopted in the form that pins D1's true
// behaviour). The epoch check is ADVISORY release nomination: two archived
// interpreters with DIFFERENT hashes and an IDENTICAL first --version line
// are both epoch-eligible under the same nomination — and each published
// descriptor keeps its OWN Interpreter HashRef (the authoritative pin, D1),
// never a release-level stand-in, so card, invocation and replay can never
// conflate them. (b) A release the registry does not nominate is refused
// with the typed mismatch and the head is unchanged.
func TestPublishSetEpochTwinInterpretersKeepDistinctPins(t *testing.T) {
	const release = "TWIN-FAKE v1"
	s, arch, first := publisherStore(t, release, release)
	second, err := arch.Archive(fakeInterpreterScript(t, release+"\nCommit: twin-second", 0))
	if err != nil {
		t.Fatalf("archive twin interpreter: %v", err)
	}
	if first == second {
		t.Fatal("premise: the twin interpreters must have different hashes")
	}
	ctx := context.Background()
	pub := NewPublisher(s, arch)
	for _, ref := range []hashref.HashRef{first, second} {
		epochs, got, err := pub.EpochsForInterpreter(ctx, ref)
		if err != nil || got != release || len(epochs) != 1 || epochs[0] != 1 {
			t.Fatalf("EpochsForInterpreter(%s) = (%v, %q, %v), want ([1], %q)", ref, epochs, got, err, release)
		}
	}

	alpha := storedSourceDescriptor(t, s, first, "tools.alpha")
	beta := storedSourceDescriptor(t, s, second, "tools.beta")
	res, err := pub.PublishSet(ctx, []Change{{ID: alpha.ID, Descriptor: &alpha}, {ID: beta.ID, Descriptor: &beta}})
	if err != nil || res.Revision != 1 || res.Unchanged {
		t.Fatalf("twin publish = (%+v, %v), want revision 1 with both entries", res, err)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := snap.List()
	if len(got) != 2 || got[0].ID != "tools.alpha" || got[1].ID != "tools.beta" {
		t.Fatalf("twin entries = %+v, want tools.alpha and tools.beta", got)
	}
	if got[0].Interpreter != first || got[1].Interpreter != second {
		t.Fatalf("published pins = (%s, %s), want each descriptor's OWN interpreter (%s, %s): the epoch nominates a release, the HashRef stays authoritative",
			got[0].Interpreter, got[1].Interpreter, first, second)
	}

	other, err := arch.Archive(fakeInterpreterScript(t, "TWIN-OTHER v2", 0))
	if err != nil {
		t.Fatal(err)
	}
	gamma := storedSourceDescriptor(t, s, other, "tools.gamma")
	_, err = pub.PublishSet(ctx, []Change{{ID: gamma.ID, Descriptor: &gamma}})
	var mismatch *InterpreterEpochMismatchError
	if !errors.As(err, &mismatch) || mismatch.ID != "tools.gamma" || mismatch.Release != "TWIN-OTHER v2" || len(mismatch.Nominating) != 0 {
		t.Fatalf("non-nominated release = %v, want *InterpreterEpochMismatchError for tools.gamma with no nominating epoch", err)
	}
	after, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if after.Head != snap.Head || after.Revision != 1 || len(after.List()) != 2 {
		t.Fatalf("after refusal: head %s rev %d entries %d, want unchanged (%s, 1, 2)", after.Head, after.Revision, len(after.List()), snap.Head)
	}
}

// TestPublishSetEpochCheckIsPerDescriptorInterpreter closes the evaluator's
// surviving mutation (iter-195 eval r1, E-finding 1): every earlier batch test
// pinned interpreters nominated by the SAME epoch, so a verifyEpochs that
// looked every descriptor up under the FIRST descriptor's interpreter passed
// the whole suite. Here the epoch registry nominates two disjoint releases
// under two epochs, and one batch publishes a descriptor for each: the batch
// is accepted only if each descriptor's epoch is checked against ITS OWN
// interpreter, and a crossed pairing is refused with the head unchanged.
func TestPublishSetEpochCheckIsPerDescriptorInterpreter(t *testing.T) {
	const relA, relB = "EPOCH-A v1", "EPOCH-B v2"
	s, arch, refA := publisherStore(t, relA, "")
	refB, err := arch.Archive(fakeInterpreterScript(t, relB, 0))
	if err != nil {
		t.Fatalf("archive second interpreter: %v", err)
	}
	payload, err := registry.Registry{SemanticID: registry.SemanticID, Epochs: []registry.EpochRecord{
		{Epoch: 1, Candidates: []string{relA}},
		{Epoch: 2, Candidates: []string{relB}},
	}}.Encode()
	if err != nil {
		t.Fatal(err)
	}
	obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(registry.SemanticID)),
		SemanticID: registry.SemanticID, Provenance: "test-two-epoch-registry", Payload: payload}
	if err := s.PutObject(obj); err != nil {
		t.Fatal(err)
	}
	if err := s.CompareAndSetRegistryHead(registry.SemanticID, hashref.HashRef{}, obj.Hash); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	pub := NewPublisher(s, arch)

	crossed := storedSourceDescriptor(t, s, refA, "tools.crossed")
	crossed.SemanticsEpoch = 2
	var mismatch *InterpreterEpochMismatchError
	if _, err := pub.PublishSet(ctx, []Change{{ID: crossed.ID, Descriptor: &crossed}}); !errors.As(err, &mismatch) || mismatch.ID != "tools.crossed" {
		t.Fatalf("crossed pairing (release %q under epoch 2) = %v, want *InterpreterEpochMismatchError", relA, err)
	}
	if _, _, ok, _ := NewReader(s).CurrentRevision(ctx); ok {
		t.Fatal("a refused publish must leave no head")
	}

	alpha := storedSourceDescriptor(t, s, refA, "tools.alpha")
	alpha.SemanticsEpoch = 1
	beta := storedSourceDescriptor(t, s, refB, "tools.beta")
	beta.SemanticsEpoch = 2
	res, err := pub.PublishSet(ctx, []Change{{ID: alpha.ID, Descriptor: &alpha}, {ID: beta.ID, Descriptor: &beta}})
	if err != nil || res.Revision != 1 {
		t.Fatalf("mixed-epoch batch = (%+v, %v), want revision 1: each descriptor's epoch must be checked against its OWN interpreter", res, err)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snap.List(); len(got) != 2 || got[0].SemanticsEpoch != 1 || got[1].SemanticsEpoch != 2 {
		t.Fatalf("published entries = %+v, want tools.alpha@1 and tools.beta@2", got)
	}
}
