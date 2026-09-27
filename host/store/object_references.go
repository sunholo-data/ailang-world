package store

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

const MaxObjectReferencePage = 500

type ReferenceKind int

const (
	ReferenceTransitionRef ReferenceKind = iota
	ReferenceTransitionFn
	ReferenceInterpreter
	ReferenceStateRoot
)

type ObjectReferenceCursor struct {
	Kind       ReferenceKind
	EntryIndex int64
	WorldRef   hashref.HashRef
}

type ObjectReference struct{ Cursor ObjectReferenceCursor }

type InvalidObjectReferenceCursorError struct {
	Cursor ObjectReferenceCursor
	Reason string
}

func (e *InvalidObjectReferenceCursorError) Error() string {
	return fmt.Sprintf("store: invalid object reference cursor: %s", e.Reason)
}

func validateObjectReferences(ref hashref.HashRef, after *ObjectReferenceCursor, limit int) error {
	const op = "ObjectReferences"
	if limit < 1 || limit > MaxObjectReferencePage {
		return &InvalidLimitError{Op: op, Limit: limit, Max: MaxObjectReferencePage}
	}
	if _, err := hashref.Parse(ref.String()); err != nil {
		return &InvalidRefError{Op: op, Field: "ref", Text: ref.String(), Err: err}
	}
	if after == nil {
		return nil
	}
	if after.Kind < ReferenceTransitionRef || after.Kind > ReferenceStateRoot {
		return &InvalidObjectReferenceCursorError{*after, "unknown kind"}
	}
	if after.Kind == ReferenceStateRoot {
		if after.EntryIndex != 0 {
			return &InvalidObjectReferenceCursorError{*after, "state-root cursor has entry index"}
		}
		if _, err := hashref.Parse(after.WorldRef.String()); err != nil {
			return &InvalidObjectReferenceCursorError{*after, "invalid world ref"}
		}
	} else {
		if after.EntryIndex < 0 {
			return &InvalidObjectReferenceCursorError{*after, "negative entry index"}
		}
		if !after.WorldRef.IsZero() {
			return &InvalidObjectReferenceCursorError{*after, "log cursor has world ref"}
		}
	}
	return nil
}

// ObjectReferences returns a bounded, relation-ordered page of entry/world edges.
func (s *Store) ObjectReferences(ctx context.Context, ref hashref.HashRef, after *ObjectReferenceCursor, limit int) ([]ObjectReference, error) {
	if err := validateObjectReferences(ref, after, limit); err != nil {
		return nil, err
	}
	return []ObjectReference{}, nil
}
