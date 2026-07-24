// Package registry implements the epoch registry from Decision 5 of the
// w-world-library-m1 design. The registry has semantic ID
// "world/epoch-registry/v1". Each revision is an IMMUTABLE object containing
// ordered EpochRecord values; the store's epoch_registry_heads table maps that
// semantic ID to the selected revision's HashRef.
//
// A registry revision is stored through the same object mechanism as every other
// kernel value (store.PutObject for the immutable object, then SetRegistryHead
// to name it): registry updates are ordinary logged/named world state, NOT a
// privileged side channel.
//
// Candidate nomination is advisory compatibility metadata only. Authoritative
// replay resolves an interpreter from the log entry's own interpreter HashRef
// exclusively (Decision 6/7); the registry is for inspection and later
// nomination workflows.
package registry

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// SemanticID is the semantic ID / registry name of the epoch registry
// (Decision 5). It equals store.EpochRegistryV1 and is the key used in
// epoch_registry_heads.
const SemanticID = store.EpochRegistryV1

// firstEpoch is the epoch number created by Bootstrap.
const firstEpoch int64 = 1

// EpochRecord is one entry in an epoch-registry revision: a semantics epoch and
// its ordered list of nominated interpreter candidates. Nomination is advisory
// compatibility metadata; the first element is the primary/initial candidate.
type EpochRecord struct {
	// Epoch is the semantics-epoch number this record describes.
	Epoch int64 `json:"epoch"`
	// Candidates is the ordered list of nominated interpreter release strings.
	// Candidates[0] is the first (primary) nomination. This is advisory metadata:
	// authoritative replay uses the log entry's interpreter HashRef exclusively.
	Candidates []string `json:"candidates"`
}

// Registry is the decoded content of one immutable epoch-registry revision: the
// ordered EpochRecord values. Its canonical bytes (see Encode) are what the
// object's HashRef addresses.
type Registry struct {
	// SemanticID is the registry's semantic ID; always SemanticID for v1.
	SemanticID string `json:"semantic_id"`
	// Epochs is the ordered list of epoch records, ascending by epoch.
	Epochs []EpochRecord `json:"epochs"`
}

// Encode returns the canonical byte encoding of a registry revision: compact,
// deterministic JSON. Field order is fixed by the struct and Go's encoder emits
// object keys in struct-declaration order, so identical registries always
// produce identical bytes and therefore identical HashRefs.
func (r Registry) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return nil, fmt.Errorf("registry: encode: %w", err)
	}
	// json.Encoder appends a trailing newline; keep it as part of the canonical
	// bytes so decode/encode round-trips are byte-stable.
	return buf.Bytes(), nil
}

// Decode parses canonical registry bytes produced by Encode.
func Decode(payload []byte) (Registry, error) {
	var r Registry
	if err := json.Unmarshal(payload, &r); err != nil {
		return Registry{}, fmt.Errorf("registry: decode: %w", err)
	}
	return r, nil
}

// object builds the immutable store.Object for a registry revision, computing
// its content address over the canonical bytes.
func object(r Registry) (store.Object, error) {
	payload, err := r.Encode()
	if err != nil {
		return store.Object{}, err
	}
	return store.Object{
		Hash: hashref.SumSHA256(payload),
		// The interface hash pins the registry's schema identity. M1 uses the
		// semantic ID's hash as a stable, self-describing interface tag.
		InterfaceHash: hashref.SumSHA256([]byte(SemanticID)),
		SemanticID:    SemanticID,
		Provenance:    "epoch-registry-bootstrap",
		Payload:       payload,
	}, nil
}

// Bootstrap creates epoch 1 of the registry with releaseString as its first
// nominated candidate and names the revision through the store's registry head.
//
// It is IDEMPOTENT: if a registry head already exists, Bootstrap loads the named
// object and verifies it is the identical epoch-1 revision (returning it), rather
// than creating a divergent epoch 1. Because the revision is content-addressed,
// re-running with the same releaseString yields the same object and head, and
// PutObject/SetRegistryHead are themselves idempotent.
//
// Storage uses the ordinary object mechanism: PutObject for the immutable
// revision object, then SetRegistryHead to map SemanticID to it. Registry
// updates are ordinary named world state, not a privileged bypass.
func Bootstrap(s *store.Store, releaseString string) (Registry, hashref.HashRef, error) {
	want := Registry{
		SemanticID: SemanticID,
		Epochs: []EpochRecord{
			{Epoch: firstEpoch, Candidates: []string{releaseString}},
		},
	}
	obj, err := object(want)
	if err != nil {
		return Registry{}, hashref.HashRef{}, err
	}

	// Idempotence: if a head already exists, it must resolve to the identical
	// epoch-1 revision. A head naming different bytes is a real divergence and is
	// surfaced as an error rather than silently overwritten.
	if head, ok, err := s.GetRegistryHead(SemanticID); err != nil {
		return Registry{}, hashref.HashRef{}, err
	} else if ok {
		if head.String() != obj.Hash.String() {
			return Registry{}, hashref.HashRef{}, fmt.Errorf(
				"registry: existing head %q diverges from bootstrap revision %q",
				head.String(), obj.Hash.String())
		}
		existing, found, err := s.GetObject(head)
		if err != nil {
			return Registry{}, hashref.HashRef{}, err
		}
		if !found {
			return Registry{}, hashref.HashRef{}, fmt.Errorf(
				"registry: head %q names an absent object", head.String())
		}
		decoded, err := Decode(existing.Payload)
		if err != nil {
			return Registry{}, hashref.HashRef{}, err
		}
		return decoded, head, nil
	}

	// Fresh bootstrap: store the immutable revision object, then name it. Both
	// operations are idempotent, so a partial prior run self-heals on retry.
	if err := s.PutObject(obj); err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	if err := s.SetRegistryHead(SemanticID, obj.Hash); err != nil {
		return Registry{}, hashref.HashRef{}, err
	}
	return want, obj.Hash, nil
}
