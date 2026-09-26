package transitionreg

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// EmptyGenesisError reports a first revision that would contain no entries.
type EmptyGenesisError struct{}

func (e *EmptyGenesisError) Error() string {
	return "publish transition registry: a first revision needs at least one descriptor"
}

// SameIDConflictError reports the CAS-retry merge refusing to continue: the
// concurrent winner's head already holds this entry ID with DIFFERENT
// canonical bytes than this publisher was about to write. Nothing was
// written by the loser — the head stays at the winner's revision — and the
// error names the ID so the operator can resolve it by hand.
type SameIDConflictError struct {
	ID       string
	Revision int64
}

func (e *SameIDConflictError) Error() string {
	return fmt.Sprintf(
		"publish transition registry: entry %q is already published with different bytes in revision %d by a concurrent publisher; nothing was written",
		e.ID, e.Revision)
}

// refuseConflictingSameIDs is the same-ID retry rule (kimi's round-1
// objection): it compares each incoming change's canonical descriptor bytes
// with the winner's head entry for that ID. Different bytes → typed
// SameIDConflictError (nothing written). Identical bytes are NOT a conflict:
// the retry converges to Unchanged because the winner already wrote them.
func refuseConflictingSameIDs(winner Revision, changes []Change) error {
	ours := make(map[string][]byte, len(changes))
	for _, change := range changes {
		if change.Descriptor == nil {
			continue
		}
		encoded, err := descriptorBytes(*change.Descriptor)
		if err != nil {
			return fmt.Errorf("publish transition registry: %w", err)
		}
		ours[change.ID] = encoded
	}
	for i := range winner.Entries {
		mine, requested := ours[winner.Entries[i].ID]
		if !requested {
			continue
		}
		theirs, err := descriptorBytes(winner.Entries[i])
		if err != nil {
			return fmt.Errorf("publish transition registry: %w", err)
		}
		if !bytes.Equal(mine, theirs) {
			return &SameIDConflictError{ID: winner.Entries[i].ID, Revision: winner.Revision}
		}
	}
	return nil
}

// descriptorBytes is one descriptor's canonical encoding — the same bytes the
// store holds — so same-ID comparisons are byte-exact, not field-by-field.
func descriptorBytes(d Descriptor) ([]byte, error) {
	return EncodeRevision(Revision{
		SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1,
		Revision: 1, Entries: []Descriptor{d},
	})
}

// buildNextRevision constructs the next revision. The genesis case (no head)
// is built DIRECTLY as revision 1 with a zero parent: BuildNext over a
// synthetic revision 0 pins Parent to the hash of the ENCODED empty revision,
// which Publish's absent-head branch refuses ("parent is not captured head").
func buildNextRevision(ok bool, current Revision, changes []Change) (Revision, error) {
	if ok {
		next, err := BuildNext(current, changes)
		if err != nil {
			return Revision{}, fmt.Errorf("publish transition registry: %w", err)
		}
		return next, nil
	}
	entries := make([]Descriptor, 0, len(changes))
	for _, change := range changes {
		if change.Descriptor == nil {
			return Revision{}, errors.New("publish transition registry: genesis cannot remove " + change.ID)
		}
		if err := change.Descriptor.Validate(); err != nil {
			return Revision{}, fmt.Errorf("publish transition registry: descriptor %q: %w", change.ID, err)
		}
		entries = append(entries, *change.Descriptor)
	}
	if len(entries) == 0 {
		return Revision{}, &EmptyGenesisError{}
	}
	sorted, err := SortedDescriptors(entries)
	if err != nil {
		return Revision{}, fmt.Errorf("publish transition registry: %w", err)
	}
	return Revision{
		SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1,
		Revision: 1, Parent: hashref.HashRef{}, Entries: sorted,
	}, nil
}
