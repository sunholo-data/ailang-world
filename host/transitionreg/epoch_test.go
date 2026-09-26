package transitionreg

import (
	"context"
	"errors"
	"testing"
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
