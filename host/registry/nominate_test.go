package registry

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Row 139 M2 (AC2 registry half, AC3 registry half). Release strings here are
// opaque candidate strings; the manifest-shaped ones are only needed by
// NominateCommand.
const (
	relA = "AILANG test-a"
	relB = "AILANG test-b"
	relC = "AILANG test-c"
	relD = "AILANG test-d"
)

func headOf(t *testing.T, s *store.Store) hashref.HashRef {
	t.Helper()
	h, ok, err := s.GetRegistryHead(boundedTestContext(t), SemanticID)
	if err != nil || !ok {
		t.Fatalf("GetRegistryHead: ok=%v err=%v", ok, err)
	}
	return h
}

func objectBytes(t *testing.T, s *store.Store, ref hashref.HashRef) store.Object {
	t.Helper()
	obj, ok, err := s.GetObject(boundedTestContext(t), ref)
	if err != nil || !ok {
		t.Fatalf("GetObject(%s): ok=%v err=%v", ref, ok, err)
	}
	return obj
}

// TestBootstrapFreshStoreIsByteIdentical pins the fresh-store revision to a
// literal: the canonical bytes, their content address and the bootstrap
// provenance. Any change to the encoding or to the fresh path reds here.
func TestBootstrapFreshStoreIsByteIdentical(t *testing.T) {
	s := openMem(t)
	_, head, err := Bootstrap(boundedTestContext(t), s, m1Release)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"semantic_id":"world/epoch-registry/v1","epochs":[{"epoch":1,"candidates":["AILANG v0.30.0 (commit e37b370)"]}]}` + "\n"
	obj := objectBytes(t, s, head)
	if string(obj.Payload) != want {
		t.Fatalf("fresh bootstrap bytes = %q, want %q", obj.Payload, want)
	}
	if head != hashref.SumSHA256([]byte(want)) {
		t.Fatalf("head %s is not the sha256 of the canonical bytes", head)
	}
	if obj.Provenance != "epoch-registry-bootstrap" || obj.InterfaceHash != hashref.SumSHA256([]byte(SemanticID)) {
		t.Fatalf("object metadata drifted: provenance %q interface %s", obj.Provenance, obj.InterfaceHash)
	}
}

func TestNominateCandidateAndNewEpochShapes(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	_, h0, err := Bootstrap(ctx, s, relA)
	if err != nil {
		t.Fatal(err)
	}
	before := objectBytes(t, s, h0)

	reg1, h1, err := Nominate(ctx, s, h0, relB, ModeCandidate)
	if err != nil {
		t.Fatalf("candidate: %v", err)
	}
	if want := []EpochRecord{{Epoch: 1, Candidates: []string{relA, relB}}}; !reflect.DeepEqual(reg1.Epochs, want) {
		t.Fatalf("candidate shape = %+v, want %+v", reg1.Epochs, want)
	}
	reg2, h2, err := Nominate(ctx, s, h1, relC, ModeNewEpoch)
	if err != nil {
		t.Fatalf("new epoch: %v", err)
	}
	want2 := []EpochRecord{{Epoch: 1, Candidates: []string{relA, relB}}, {Epoch: 2, Candidates: []string{relC}}}
	if !reflect.DeepEqual(reg2.Epochs, want2) {
		t.Fatalf("new-epoch shape = %+v, want %+v", reg2.Epochs, want2)
	}
	// A candidate after a new epoch joins the CURRENT (highest) epoch.
	reg3, h3, err := Nominate(ctx, s, h2, relD, ModeCandidate)
	if err != nil {
		t.Fatal(err)
	}
	if got := reg3.Epochs[1].Candidates; !reflect.DeepEqual(got, []string{relC, relD}) || !reflect.DeepEqual(reg3.Epochs[0], want2[0]) {
		t.Fatalf("candidate after new epoch = %+v", reg3.Epochs)
	}
	if headOf(t, s) != h3 {
		t.Fatal("the head does not name the last revision")
	}

	// Immutability: every revision is its own object, each decodes to the
	// registry it was when written, and the bootstrap object is unchanged.
	seen := map[hashref.HashRef]bool{}
	for i, c := range []struct {
		ref  hashref.HashRef
		want int
	}{{h0, 1}, {h1, 1}, {h2, 2}, {h3, 2}} {
		if seen[c.ref] {
			t.Fatalf("revision %d reuses an earlier head %s", i, c.ref)
		}
		seen[c.ref] = true
		obj := objectBytes(t, s, c.ref)
		if obj.Hash != hashref.SumSHA256(obj.Payload) {
			t.Fatalf("revision %d is not content-addressed", i)
		}
		dec, err := Decode(obj.Payload)
		if err != nil || len(dec.Epochs) != c.want {
			t.Fatalf("revision %d decodes to %+v (%v)", i, dec, err)
		}
		if i > 0 && obj.Provenance != "epoch-registry-nominate" {
			t.Fatalf("revision %d provenance %q", i, obj.Provenance)
		}
	}
	if after := objectBytes(t, s, h0); string(after.Payload) != string(before.Payload) {
		t.Fatal("nomination rewrote the bootstrap revision")
	}
}

func TestNominateStaleHeadRefusedAndNothingWritten(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	_, h0, err := Bootstrap(ctx, s, relA)
	if err != nil {
		t.Fatal(err)
	}
	_, h1, err := Nominate(ctx, s, h0, relB, ModeCandidate)
	if err != nil {
		t.Fatal(err)
	}
	// CONTROL above: the fresh head nominates. Now the stale one must not.
	_, _, err = Nominate(ctx, s, h0, relC, ModeCandidate)
	var conflict *store.RegistryCASConflict
	if !errors.As(err, &conflict) || conflict.Expected != h0 || conflict.Actual != h1 {
		t.Fatalf("stale head = %v, want *store.RegistryCASConflict expected %s actual %s", err, h0, h1)
	}
	if headOf(t, s) != h1 {
		t.Fatal("a refused nomination moved the head")
	}
	reg, _, _, err := Load(ctx, s)
	if err != nil || len(reg.EpochsNominating(relC)) != 0 {
		t.Fatalf("the refused release leaked into the registry: %+v %v", reg, err)
	}
}

// racingStore moves the head between Nominate's read and its CAS (inside
// PutObject), so the refusal comes from the STORE's compare-and-set, not from
// the early stale-head read.
type racingStore struct {
	*store.Store
	t     *testing.T
	raced hashref.HashRef
}

func (r *racingStore) PutObject(ctx context.Context, obj store.Object) error {
	if err := r.Store.PutObject(ctx, obj); err != nil {
		return err
	}
	if r.raced.IsZero() {
		head, _, err := r.Store.GetRegistryHead(ctx, SemanticID)
		if err != nil {
			return err
		}
		_, raced, err := Nominate(ctx, r.Store, head, relD, ModeCandidate)
		if err != nil {
			r.t.Fatalf("racer: %v", err)
		}
		r.raced = raced
	}
	return nil
}

func TestNominateLostRaceIsRefusedByTheStoreCAS(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	_, h0, err := Bootstrap(ctx, s, relA)
	if err != nil {
		t.Fatal(err)
	}
	r := &racingStore{Store: s, t: t}
	_, _, err = nominate(ctx, r, h0, relB, ModeCandidate)
	if !store.IsRegistryCASConflict(err) {
		t.Fatalf("lost race = %v, want a registry CAS conflict", err)
	}
	if r.raced.IsZero() || headOf(t, s) != r.raced {
		t.Fatal("the racer's revision is not the head: the loser overwrote it")
	}
}

func TestNominateRefusals(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	_, h0, err := Bootstrap(ctx, s, relA)
	if err != nil {
		t.Fatal(err)
	}
	_, h1, err := Nominate(ctx, s, h0, relB, ModeNewEpoch)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name    string
		head    hashref.HashRef
		release string
		mode    Mode
		want    string
	}{
		{"already in the current epoch", h1, relB, ModeCandidate, "already nominated by epoch(s) [2]"},
		{"already in an EARLIER epoch", h1, relA, ModeNewEpoch, "already nominated by epoch(s) [1]"},
		{"blank release", h1, " ", ModeCandidate, "blank or untrimmed"},
		{"untrimmed release", h1, relC + " ", ModeCandidate, "blank or untrimmed"},
		{"unknown mode", h1, relC, Mode("sideways"), "unknown nomination mode"},
		{"zero head", hashref.HashRef{}, relC, ModeCandidate, "needs the current registry head"},
	} {
		_, _, err := Nominate(ctx, s, c.head, c.release, c.mode)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v, want it to contain %q", c.name, err, c.want)
		}
	}
	var already *AlreadyNominatedError
	if _, _, err := Nominate(ctx, s, h1, relA, ModeCandidate); !errors.As(err, &already) {
		t.Fatalf("already-nominated is not typed: %T", err)
	}
	if headOf(t, s) != h1 {
		t.Fatal("a refusal moved the head")
	}
	// No head at all: nominate never bootstraps.
	if _, _, err := Nominate(ctx, openMem(t), h1, relC, ModeCandidate); !store.IsRegistryCASConflict(err) {
		t.Fatalf("headless store = %v, want a CAS conflict (no head)", err)
	}
}

func TestNextIsPure(t *testing.T) {
	cur := Registry{SemanticID: SemanticID, Epochs: []EpochRecord{{Epoch: 1, Candidates: []string{relA}}}}
	snapshot, _ := cur.Encode()
	if _, err := Next(cur, relB, ModeCandidate); err != nil {
		t.Fatal(err)
	}
	if _, err := Next(cur, relC, ModeNewEpoch); err != nil {
		t.Fatal(err)
	}
	if again, _ := cur.Encode(); string(again) != string(snapshot) {
		t.Fatal("Next mutated its input")
	}
	for name, bad := range map[string]Registry{
		"wrong semantic id": {SemanticID: "x", Epochs: cur.Epochs},
		"no epochs":         {SemanticID: SemanticID},
		"gap":               {SemanticID: SemanticID, Epochs: []EpochRecord{{Epoch: 2, Candidates: []string{relA}}}},
		"empty epoch":       {SemanticID: SemanticID, Epochs: []EpochRecord{{Epoch: 1}}},
		"duplicate across epochs": {SemanticID: SemanticID, Epochs: []EpochRecord{
			{Epoch: 1, Candidates: []string{relA}}, {Epoch: 2, Candidates: []string{relA}}}},
	} {
		if _, err := Next(bad, relD, ModeCandidate); err == nil {
			t.Errorf("%s: Next accepted an invalid registry", name)
		}
	}
}

// AC3, registry half: an existing head admits ANY nominated release without
// writing, and refuses an un-nominated one with a typed error.
func TestBootstrapOnExistingHeadAdmitsNominatedReleases(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	_, h0, err := Bootstrap(ctx, s, relA)
	if err != nil {
		t.Fatal(err)
	}
	_, h1, err := Nominate(ctx, s, h0, relB, ModeCandidate)
	if err != nil {
		t.Fatal(err)
	}
	_, h2, err := Nominate(ctx, s, h1, relC, ModeNewEpoch)
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{relA, relB, relC} {
		reg, head, err := Bootstrap(ctx, s, rel)
		if err != nil {
			t.Fatalf("Bootstrap(%q) on a grown registry: %v", rel, err)
		}
		if head != h2 || len(reg.Epochs) != 2 {
			t.Fatalf("Bootstrap(%q) = head %s epochs %d, want the grown head %s", rel, head, len(reg.Epochs), h2)
		}
	}
	_, _, err = Bootstrap(ctx, s, relD)
	var nn *NotNominatedError
	if !errors.As(err, &nn) || nn.Release != relD || nn.Head != h2 {
		t.Fatalf("un-nominated = %v, want *NotNominatedError for %q at %s", err, relD, h2)
	}
	if headOf(t, s) != h2 {
		t.Fatal("startup on an existing head wrote")
	}
}

func TestNominateCommand(t *testing.T) {
	got, ok := NominateCommand("/w/world.db", "AILANG v0.99.1", "/bin/ailang")
	want := "world-publish epoch nominate --store /w/world.db --release v0.99.1 --candidate --gate-verdict <verdict.json> --binary /bin/ailang"
	if !ok || got != want {
		t.Fatalf("NominateCommand = %q %v, want %q", got, ok, want)
	}
	for _, line := range []string{"unpinned", "AILANG v0.99.1-3-gabc-dirty", "v0.99.1", ""} {
		if _, ok := NominateCommand("/w", line, "/bin/ailang"); ok {
			t.Errorf("NominateCommand accepted %q", line)
		}
	}
}
