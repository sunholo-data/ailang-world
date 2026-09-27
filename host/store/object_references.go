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

const transitionReferencesSQL = "SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?"
const functionReferencesSQL = "SELECT entry_index FROM log_entries WHERE transition_fn_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?"
const interpreterReferencesSQL = "SELECT entry_index FROM log_entries WHERE interpreter_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?"
const worldReferencesSQL = "SELECT world_ref FROM worlds WHERE state_root = ? AND world_ref > ? ORDER BY world_ref LIMIT ?"

// objectReferencesBeforeQuery observes an actual branch query in store tests.
var objectReferencesBeforeQuery func()

// ObjectReferences returns at most limit edges, ordered by relation then source key.
func (s *Store) ObjectReferences(ctx context.Context, ref hashref.HashRef, after *ObjectReferenceCursor, limit int) ([]ObjectReference, error) {
	if err := validateObjectReferences(ref, after, limit); err != nil {
		return nil, err
	}
	if !s.referenceIndexesAvailable {
		return nil, &ReferenceIndexUnavailableError{}
	}
	items := make([]ObjectReference, 0, limit)
	queries := [...]string{transitionReferencesSQL, functionReferencesSQL, interpreterReferencesSQL, worldReferencesSQL}
	for kind := ReferenceTransitionRef; kind <= ReferenceStateRoot && len(items) < limit; kind++ {
		if after != nil && kind < after.Kind {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		query := queries[kind]
		var afterKey any
		if kind == ReferenceStateRoot {
			afterKey = ""
			if after != nil && kind == after.Kind {
				afterKey = after.WorldRef.String()
			}
		} else {
			afterKey = int64(-1)
			if after != nil && kind == after.Kind {
				afterKey = after.EntryIndex
			}
		}
		remaining := limit - len(items)
		if objectReferencesBeforeQuery != nil {
			objectReferencesBeforeQuery()
		}
		rows, err := s.db.QueryContext(ctx, query, ref.String(), afterKey, remaining)
		if err != nil {
			return nil, fmt.Errorf("store: object references query: %w", err)
		}
		for rows.Next() {
			cursor := ObjectReferenceCursor{Kind: kind}
			if kind == ReferenceStateRoot {
				var text string
				if err = rows.Scan(&text); err == nil {
					cursor.WorldRef, err = hashref.Parse(text)
				}
			} else {
				err = rows.Scan(&cursor.EntryIndex)
				if err == nil && cursor.EntryIndex < 0 {
					err = fmt.Errorf("negative stored entry index %d", cursor.EntryIndex)
				}
			}
			if err != nil {
				break
			}
			items = append(items, ObjectReference{Cursor: cursor})
		}
		if err == nil {
			err = rows.Err()
		}
		closeErr := rows.Close()
		if err != nil {
			return nil, fmt.Errorf("store: object references row: %w", err)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("store: close object references: %w", closeErr)
		}
	}
	return items, nil
}
