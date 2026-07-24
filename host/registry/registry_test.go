package registry

import (
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// openMem opens a fresh in-memory store for a test and registers cleanup.
func openMem(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("store.Open(:memory:): %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

const m1Release = "AILANG v0.30.0 (commit e37b370)"

func TestBootstrapCreatesEpochOneWithReleaseCandidate(t *testing.T) {
	s := openMem(t)

	reg, head, err := Bootstrap(s, m1Release)
	if err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}
	if head.IsZero() {
		t.Fatal("Bootstrap returned a zero head")
	}

	// Epoch 1 exists and names the M1 release string as its FIRST candidate.
	if len(reg.Epochs) != 1 {
		t.Fatalf("expected 1 epoch, got %d", len(reg.Epochs))
	}
	if reg.Epochs[0].Epoch != 1 {
		t.Fatalf("first epoch number = %d, want 1", reg.Epochs[0].Epoch)
	}
	if len(reg.Epochs[0].Candidates) == 0 || reg.Epochs[0].Candidates[0] != m1Release {
		t.Fatalf("epoch-1 first candidate = %v, want %q", reg.Epochs[0].Candidates, m1Release)
	}
	if reg.SemanticID != SemanticID {
		t.Fatalf("semantic ID = %q, want %q", reg.SemanticID, SemanticID)
	}
}

func TestBootstrapStoresThroughObjectAndRegistryHead(t *testing.T) {
	s := openMem(t)

	_, head, err := Bootstrap(s, m1Release)
	if err != nil {
		t.Fatalf("Bootstrap: %v", err)
	}

	// The head is named via the ordinary store registry-head mechanism.
	gotHead, ok, err := s.GetRegistryHead(SemanticID)
	if err != nil || !ok {
		t.Fatalf("GetRegistryHead: ok=%v err=%v", ok, err)
	}
	if gotHead.String() != head.String() {
		t.Fatalf("registry head = %q, want %q", gotHead.String(), head.String())
	}

	// The revision itself is an ordinary immutable object addressable by the head.
	obj, ok, err := s.GetObject(head)
	if err != nil || !ok {
		t.Fatalf("GetObject(head): ok=%v err=%v", ok, err)
	}
	if obj.SemanticID != SemanticID {
		t.Fatalf("object semantic ID = %q, want %q", obj.SemanticID, SemanticID)
	}
	// The object's content address addresses its own canonical bytes, and those
	// bytes decode back to the same registry (round-trip stable).
	decoded, err := Decode(obj.Payload)
	if err != nil {
		t.Fatalf("Decode(object payload): %v", err)
	}
	if decoded.Epochs[0].Candidates[0] != m1Release {
		t.Fatalf("decoded first candidate = %q, want %q", decoded.Epochs[0].Candidates[0], m1Release)
	}
}

func TestBootstrapIsIdempotent(t *testing.T) {
	s := openMem(t)

	reg1, head1, err := Bootstrap(s, m1Release)
	if err != nil {
		t.Fatalf("first Bootstrap: %v", err)
	}
	// Running twice does not create a divergent epoch 1: same head, same content.
	reg2, head2, err := Bootstrap(s, m1Release)
	if err != nil {
		t.Fatalf("second Bootstrap (idempotent): %v", err)
	}
	if head1.String() != head2.String() {
		t.Fatalf("idempotent bootstrap changed head: %q -> %q", head1.String(), head2.String())
	}
	if len(reg2.Epochs) != 1 || reg2.Epochs[0].Epoch != 1 {
		t.Fatalf("second bootstrap did not yield a single epoch 1: %+v", reg2.Epochs)
	}
	// Encodings are byte-identical across runs (content addressing is stable).
	b1, err := reg1.Encode()
	if err != nil {
		t.Fatalf("encode reg1: %v", err)
	}
	b2, err := reg2.Encode()
	if err != nil {
		t.Fatalf("encode reg2: %v", err)
	}
	if string(b1) != string(b2) {
		t.Fatalf("idempotent bootstrap produced different canonical bytes")
	}
}

func TestBootstrapDetectsDivergentHead(t *testing.T) {
	s := openMem(t)
	if _, _, err := Bootstrap(s, m1Release); err != nil {
		t.Fatalf("first Bootstrap: %v", err)
	}
	// A second bootstrap with a DIFFERENT release string would produce different
	// content-addressed bytes; the existing head diverges, which must be an error
	// rather than a silent overwrite of epoch 1.
	if _, _, err := Bootstrap(s, "AILANG v0.31.0 (commit deadbee)"); err == nil {
		t.Fatal("expected divergent-head error, got nil")
	}
}

func TestEncodeIsDeterministic(t *testing.T) {
	r := Registry{
		SemanticID: SemanticID,
		Epochs: []EpochRecord{
			{Epoch: 1, Candidates: []string{"a", "b"}},
			{Epoch: 2, Candidates: []string{"c"}},
		},
	}
	first, err := r.Encode()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	for i := 0; i < 20; i++ {
		again, err := r.Encode()
		if err != nil {
			t.Fatalf("encode iter %d: %v", i, err)
		}
		if string(again) != string(first) {
			t.Fatalf("encode is nondeterministic at iter %d", i)
		}
	}
	// Round-trips back to the same structure.
	decoded, err := Decode(first)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(decoded.Epochs) != 2 || decoded.Epochs[1].Candidates[0] != "c" {
		t.Fatalf("round-trip mismatch: %+v", decoded)
	}
}
