package transitionreg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type EffectRequirement struct {
	Effect, Scope string
	Cost          int64
}
type Descriptor struct {
	ID                        string
	TransitionFn, Interpreter hashref.HashRef
	SemanticsEpoch            int64
	InputSchema, OutputSchema []byte
	Access                    EffectRequirement
	DeclaredEffects           []EffectRequirement
	Title, Description        string
}
type Revision struct {
	SemanticID    string
	InterfaceHash hashref.HashRef
	Revision      int64
	Parent        hashref.HashRef
	Entries       []Descriptor
}

// ObjectStore is the persistence boundary used by readers and publishers.
// Keeping it narrow makes store failures directly injectable in tests.
type ObjectStore interface {
	GetRegistryHead(string) (hashref.HashRef, bool, error)
	GetObject(hashref.HashRef) (store.Object, bool, error)
	PutObject(store.Object) error
	CompareAndSetRegistryHead(string, hashref.HashRef, hashref.HashRef) error
}

// Reader exposes an eager, immutable view of one registry revision.
type Reader interface {
	ReadSnapshot(context.Context) (Snapshot, error)
}

// Snapshot is a self-contained view of a single immutable registry head.
type Snapshot struct {
	Head     hashref.HashRef
	Revision int64
	entries  []Descriptor
}

// StoreReader reads transition revisions and caches validated parsed objects by
// their immutable content hash. The registry head is never cached.
type StoreReader struct {
	store ObjectStore
	mu    sync.RWMutex
	cache map[hashref.HashRef]Snapshot
}

func NewReader(objects ObjectStore) *StoreReader {
	return &StoreReader{store: objects, cache: make(map[hashref.HashRef]Snapshot)}
}

func (r *StoreReader) ReadSnapshot(ctx context.Context) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, fmt.Errorf("read transition registry: context: %w", err)
	}
	head, ok, err := r.store.GetRegistryHead(store.TransitionRegistryV1)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read transition registry head: %w", err)
	}
	if !ok {
		return Snapshot{}, errors.New("read transition registry: head is absent")
	}

	r.mu.RLock()
	cached, found := r.cache[head]
	r.mu.RUnlock()
	if found {
		return cloneSnapshot(cached), nil
	}

	object, ok, err := r.store.GetObject(head)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read transition registry object: %w", err)
	}
	if !ok {
		return Snapshot{}, fmt.Errorf("read transition registry: object %q is absent", head.String())
	}
	if got := hashref.SumSHA256(object.Payload); got != head || object.Hash != head {
		return Snapshot{}, fmt.Errorf("read transition registry: object hash mismatch for %q", head.String())
	}
	if object.SemanticID != SemanticIDV1 {
		return Snapshot{}, fmt.Errorf("read transition registry: semantic ID %q is not %q", object.SemanticID, SemanticIDV1)
	}
	if object.InterfaceHash != InterfaceHashV1 {
		return Snapshot{}, fmt.Errorf("read transition registry: interface hash %q is not %q", object.InterfaceHash, InterfaceHashV1)
	}
	revision, err := DecodeRevision(object.Payload)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read transition registry revision: %w", err)
	}
	if revision.Revision == 0 && !revision.Parent.IsZero() {
		return Snapshot{}, errors.New("read transition registry revision: revision 0 must have no parent")
	}
	if revision.Revision > 0 && revision.Parent.IsZero() && revision.Revision != 1 {
		return Snapshot{}, errors.New("read transition registry revision: revision after 1 must have a parent")
	}

	parsed := Snapshot{Head: head, Revision: revision.Revision, entries: cloneDescriptors(revision.Entries)}
	r.mu.Lock()
	if existing, exists := r.cache[head]; exists {
		parsed = existing
	} else {
		r.cache[head] = parsed
	}
	r.mu.Unlock()
	return cloneSnapshot(parsed), nil
}

func (s Snapshot) Lookup(id string) (Descriptor, bool) {
	i := sort.Search(len(s.entries), func(i int) bool { return CompareID(s.entries[i].ID, id) >= 0 })
	if i == len(s.entries) || s.entries[i].ID != id {
		return Descriptor{}, false
	}
	return cloneDescriptor(s.entries[i]), true
}

func (s Snapshot) List() []Descriptor { return cloneDescriptors(s.entries) }

func cloneSnapshot(s Snapshot) Snapshot {
	return Snapshot{Head: s.Head, Revision: s.Revision, entries: cloneDescriptors(s.entries)}
}

func cloneDescriptors(entries []Descriptor) []Descriptor {
	out := make([]Descriptor, len(entries))
	for i := range entries {
		out[i] = cloneDescriptor(entries[i])
	}
	return out
}

func cloneDescriptor(d Descriptor) Descriptor {
	d.InputSchema = append([]byte(nil), d.InputSchema...)
	d.OutputSchema = append([]byte(nil), d.OutputSchema...)
	d.DeclaredEffects = append([]EffectRequirement(nil), d.DeclaredEffects...)
	return d
}

// Change replaces the descriptor named by ID, or removes it when Descriptor is
// nil. A replacement descriptor must carry the same stable ID.
type Change struct {
	ID         string
	Descriptor *Descriptor
}

// BuildNext applies changes without mutating current or its descriptors.
func BuildNext(current Revision, changes []Change) (Revision, error) {
	if err := current.Validate(); err != nil {
		return Revision{}, fmt.Errorf("build next: current revision: %w", err)
	}
	currentBytes, err := EncodeRevision(current)
	if err != nil {
		return Revision{}, fmt.Errorf("build next: encode current: %w", err)
	}
	byID := make(map[string]Descriptor, len(current.Entries)+len(changes))
	for _, d := range current.Entries {
		byID[d.ID] = cloneDescriptor(d)
	}
	seenChanges := make(map[string]struct{}, len(changes))
	for _, change := range changes {
		if err := validateID(change.ID); err != nil {
			return Revision{}, fmt.Errorf("build next: change ID: %w", err)
		}
		if _, duplicate := seenChanges[change.ID]; duplicate {
			return Revision{}, fmt.Errorf("build next: duplicate change ID %q", change.ID)
		}
		seenChanges[change.ID] = struct{}{}
		if change.Descriptor == nil {
			delete(byID, change.ID)
			continue
		}
		if change.Descriptor.ID != change.ID {
			return Revision{}, fmt.Errorf("build next: replacement ID %q does not match change ID %q", change.Descriptor.ID, change.ID)
		}
		if err := change.Descriptor.Validate(); err != nil {
			return Revision{}, fmt.Errorf("build next: descriptor %q: %w", change.ID, err)
		}
		byID[change.ID] = cloneDescriptor(*change.Descriptor)
	}
	entries := make([]Descriptor, 0, len(byID))
	for _, d := range byID {
		entries = append(entries, d)
	}
	entries, err = SortedDescriptors(entries)
	if err != nil {
		return Revision{}, fmt.Errorf("build next: %w", err)
	}
	next := Revision{
		SemanticID:    SemanticIDV1,
		InterfaceHash: InterfaceHashV1,
		Revision:      current.Revision + 1,
		Parent:        hashref.SumSHA256(currentBytes),
		Entries:       entries,
	}
	if _, err := EncodeRevision(next); err != nil {
		return Revision{}, fmt.Errorf("build next: encode next: %w", err)
	}
	return next, nil
}

// Publish writes next as an immutable object and then advances only the
// transition registry head with CAS. A CAS failure intentionally leaves the
// content-addressed object in the store.
func (r *StoreReader) Publish(ctx context.Context, expectedHead hashref.HashRef, next Revision) (hashref.HashRef, error) {
	if err := ctx.Err(); err != nil {
		return hashref.HashRef{}, fmt.Errorf("publish transition registry: context: %w", err)
	}
	if expectedHead.IsZero() {
		if next.Revision != 1 {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: revision %d is not expected 1", next.Revision)
		}
		if !next.Parent.IsZero() {
			return hashref.HashRef{}, errors.New("publish transition registry: parent is not captured head (absent)")
		}
	} else {
		currentObject, ok, err := r.store.GetObject(expectedHead)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: read expected revision: %w", err)
		}
		if !ok {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: expected object %q is absent", expectedHead)
		}
		current, err := DecodeRevision(currentObject.Payload)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: decode expected revision: %w", err)
		}
		if next.Revision != current.Revision+1 {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: revision %d is not expected %d", next.Revision, current.Revision+1)
		}
		if next.Parent != expectedHead {
			return hashref.HashRef{}, fmt.Errorf("publish transition registry: parent %q is not captured head %q", next.Parent, expectedHead)
		}
	}
	payload, err := EncodeRevision(next)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("publish transition registry: encode: %w", err)
	}
	ref := hashref.SumSHA256(payload)
	object := store.Object{Hash: ref, InterfaceHash: InterfaceHashV1, SemanticID: SemanticIDV1, Provenance: "host/transitionreg", Payload: payload}
	if err := r.store.PutObject(object); err != nil {
		return hashref.HashRef{}, fmt.Errorf("publish transition registry: put object: %w", err)
	}
	if err := r.store.CompareAndSetRegistryHead(store.TransitionRegistryV1, expectedHead, ref); err != nil {
		return hashref.HashRef{}, fmt.Errorf("publish transition registry: compare and set head: %w", err)
	}
	return ref, nil
}

func (d Descriptor) Validate() error {
	if err := validateID(d.ID); err != nil {
		return fmt.Errorf("descriptor ID: %w", err)
	}
	if d.TransitionFn.IsZero() {
		return errors.New("descriptor transitionFn is zero")
	}
	if d.Interpreter.IsZero() {
		return errors.New("descriptor interpreter is zero")
	}
	if d.SemanticsEpoch < 0 {
		return errors.New("descriptor semanticsEpoch is negative")
	}
	in, err := canonicalSchema(d.InputSchema)
	if err != nil {
		return fmt.Errorf("input schema: %w", err)
	}
	if !bytes.Equal(in, d.InputSchema) {
		return errors.New("input schema is not canonical")
	}
	out, err := canonicalSchema(d.OutputSchema)
	if err != nil {
		return fmt.Errorf("output schema: %w", err)
	}
	if !bytes.Equal(out, d.OutputSchema) {
		return errors.New("output schema is not canonical")
	}
	if err := validateRequirement(d.Access); err != nil {
		return fmt.Errorf("access: %w", err)
	}
	seen := map[EffectRequirement]struct{}{}
	for _, e := range d.DeclaredEffects {
		if err := validateRequirement(e); err != nil {
			return fmt.Errorf("declared effect: %w", err)
		}
		if _, ok := seen[e]; ok {
			return errors.New("duplicate declared effect")
		}
		seen[e] = struct{}{}
	}
	return nil
}
func validateRequirement(e EffectRequirement) error {
	if e.Cost < 0 {
		return errors.New("cost is negative")
	}
	return nil
}
func validateID(id string) error {
	if len(id) < 1 || len(id) > 128 {
		return errors.New("length must be 1..128 bytes")
	}
	for _, seg := range strings.FieldsFunc(id, func(r rune) bool { return r == '.' || r == '/' }) {
		if len(seg) < 1 || len(seg) > 32 {
			return errors.New("segment length must be 1..32 bytes")
		}
	}
	if strings.Contains(id, "//") || strings.Contains(id, "..") || strings.Contains(id, "/.") || strings.Contains(id, "./") {
		return errors.New("empty segment")
	}
	for i, c := range []byte(id) {
		valid := c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-' || c == '.' || c == '/'
		if !valid {
			return errors.New("does not match stable ID grammar")
		}
		if (i == 0 || i == len(id)-1 || id[i-1] == '.' || id[i-1] == '/') && (c == '_' || c == '-' || c == '.' || c == '/') {
			return errors.New("does not match stable ID grammar")
		}
		if i+1 < len(id) && (id[i+1] == '.' || id[i+1] == '/') && (c == '_' || c == '-') {
			return errors.New("does not match stable ID grammar")
		}
	}
	return nil
}
func (r Revision) Validate() error {
	if r.SemanticID != SemanticIDV1 {
		return fmt.Errorf("semanticID %q is not %q", r.SemanticID, SemanticIDV1)
	}
	if r.InterfaceHash != InterfaceHashV1 {
		return errors.New("wrong interface hash")
	}
	if r.Revision < 0 {
		return errors.New("revision is negative")
	}
	if len(r.Entries) > maxEntries {
		return fmt.Errorf("entries exceeds %d", maxEntries)
	}
	for i := range r.Entries {
		if err := r.Entries[i].Validate(); err != nil {
			return fmt.Errorf("entry %d: %w", i, err)
		}
		if i > 0 && CompareID(r.Entries[i-1].ID, r.Entries[i].ID) >= 0 {
			return errors.New("entries are not strictly ordered by ID")
		}
	}
	return nil
}
func CompareID(a, b string) int { return bytes.Compare([]byte(a), []byte(b)) }
func SortedDescriptors(entries []Descriptor) ([]Descriptor, error) {
	out := append([]Descriptor(nil), entries...)
	sort.Slice(out, func(i, j int) bool { return CompareID(out[i].ID, out[j].ID) < 0 })
	for i := range out {
		if err := out[i].Validate(); err != nil {
			return nil, err
		}
		if i > 0 && out[i-1].ID == out[i].ID {
			return nil, fmt.Errorf("duplicate descriptor ID %q", out[i].ID)
		}
	}
	return out, nil
}
