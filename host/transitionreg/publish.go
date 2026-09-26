package transitionreg

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// SetResult reports the outcome of one PublishSet call. Unchanged is true when
// the requested descriptor set already IS the head revision's content, so no
// revision was written (idempotent republish is a no-op, not a new revision).
type SetResult struct {
	Head      hashref.HashRef
	Revision  int64
	Unchanged bool
}

// CurrentRevision reads and decodes the head revision. ok=false when no head
// exists yet. It is the publisher's re-read seam for CAS retries.
func (r *StoreReader) CurrentRevision(ctx context.Context) (hashref.HashRef, Revision, bool, error) {
	head, ok, err := r.store.GetRegistryHead(ctx, store.TransitionRegistryV1)
	if err != nil {
		return hashref.HashRef{}, Revision{}, false, fmt.Errorf("publish transition registry: read head: %w", err)
	}
	if !ok {
		return hashref.HashRef{}, Revision{}, false, nil
	}
	object, found, err := r.store.GetObject(ctx, head)
	if err != nil {
		return hashref.HashRef{}, Revision{}, false, fmt.Errorf("publish transition registry: read head object: %w", err)
	}
	if !found {
		return hashref.HashRef{}, Revision{}, false, fmt.Errorf("publish transition registry: head object %q is absent", head.String())
	}
	current, err := DecodeRevision(object.Payload)
	if err != nil {
		return hashref.HashRef{}, Revision{}, false, fmt.Errorf("publish transition registry: decode head revision: %w", err)
	}
	return head, current, true, nil
}

// PublishSet publishes the descriptor changes as the next registry revision
// through the landed CAS (Publish). It is the production publisher core behind
// operator CLI verbs and, later, the package-install path (P9): every write
// goes through BuildNext + Publish's expected-head CAS — never a blind head
// overwrite. Honesty: every non-nil descriptor's TransitionFn must name an
// object the store holds (verifySources), every descriptor's SemanticsEpoch
// must be an epoch the registry nominates for the pinned interpreter
// (verifyEpochs — derived or validated, never defaulted), and every stored
// source payload must LOAD under that pinned archived interpreter
// (verifyLoadable). All three run before any write, so refusal leaves the
// head untouched. Idempotence: when the head already holds exactly the
// requested entries, nothing is written. Races: one bounded CAS retry
// (re-read, rebuild, re-publish) that MERGES different IDs and REFUSES a
// same-ID/different-bytes clobber; a second conflict surfaces the typed
// store error for the operator to re-run.
func (r *StoreReader) PublishSet(ctx context.Context, changes []Change) (SetResult, error) {
	if r.arch == nil {
		return SetResult{}, &PublisherArchiveRequiredError{}
	}
	if err := verifySources(ctx, r.store, changes); err != nil {
		return SetResult{}, err
	}
	if err := r.verifyEpochs(ctx, changes); err != nil {
		return SetResult{}, err
	}
	if err := r.verifyLoadable(ctx, changes); err != nil {
		return SetResult{}, err
	}
	head, current, ok, err := r.CurrentRevision(ctx)
	if err != nil {
		return SetResult{}, err
	}
	next, err := buildNextRevision(ok, current, changes)
	if err != nil {
		return SetResult{}, err
	}
	if ok && sameEntries(current.Entries, next.Entries) {
		return SetResult{Head: head, Revision: current.Revision, Unchanged: true}, nil
	}
	expected := head
	if !ok {
		expected = hashref.HashRef{}
	}
	ref, err := r.Publish(ctx, expected, next)
	if err == nil {
		return SetResult{Head: ref, Revision: next.Revision}, nil
	}
	if !store.IsRegistryCASConflict(err) {
		return SetResult{}, err
	}
	// One bounded retry: the head moved between our read and the CAS. Re-read,
	// rebuild over the winner (or over a still-absent head), re-publish. A
	// publisher whose changes reproduce the winner's entries converges to
	// Unchanged instead of racing forever; a publisher whose bytes CONFLICT
	// with the winner's fresh same-ID entry is refused — BuildNext replaces
	// by ID, and silently clobbering the winner is exactly the lost update
	// the CAS exists to prevent.
	head, current, ok, err = r.CurrentRevision(ctx)
	if err != nil {
		return SetResult{}, err
	}
	if err := refuseConflictingSameIDs(current, changes); err != nil {
		return SetResult{}, err
	}
	next, err = buildNextRevision(ok, current, changes)
	if err != nil {
		return SetResult{}, err
	}
	if ok && sameEntries(current.Entries, next.Entries) {
		return SetResult{Head: head, Revision: current.Revision, Unchanged: true}, nil
	}
	expected = head
	if !ok {
		expected = hashref.HashRef{}
	}
	ref, err = r.Publish(ctx, expected, next)
	if err != nil {
		return SetResult{}, err
	}
	return SetResult{Head: ref, Revision: next.Revision}, nil
}

// sameEntries compares two descriptor sets by their canonical encoding, so
// idempotence is decided by exactly the bytes the store will hold.
func sameEntries(a, b []Descriptor) bool {
	left, err := EncodeRevision(Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 1, Entries: a})
	if err != nil {
		return false
	}
	right, err := EncodeRevision(Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 1, Entries: b})
	if err != nil {
		return false
	}
	return bytes.Equal(left, right)
}
