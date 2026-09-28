package store

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// MaxObjectCommitPage bounds one ObjectCommits read.
const MaxObjectCommitPage = 500

// InvalidObjectCommitCursorError reports an afterEntry below the -1 start sentinel.
type InvalidObjectCommitCursorError struct{ AfterEntry int64 }

func (e *InvalidObjectCommitCursorError) Error() string {
	return fmt.Sprintf("store: invalid object commit cursor %d: want -1 (start) or an entry index", e.AfterEntry)
}

const objectCommitsSQL = "SELECT entry_index FROM commit_objects WHERE object_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?"

// ObjectCommits returns at most limit entry indexes whose commit carried ref,
// ascending, strictly after afterEntry (-1 starts at genesis). An empty result
// means no commit carried ref: membership is recorded from genesis on every
// store this binary opens (schema v4; older stores are refused at Open).
func (s *Store) ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	const op = "ObjectCommits"
	if limit < 1 || limit > MaxObjectCommitPage {
		return nil, &InvalidLimitError{Op: op, Limit: limit, Max: MaxObjectCommitPage}
	}
	if _, err := hashref.Parse(ref.String()); err != nil {
		return nil, &InvalidRefError{Op: op, Field: "ref", Text: ref.String(), Err: err}
	}
	if afterEntry < -1 {
		return nil, &InvalidObjectCommitCursorError{AfterEntry: afterEntry}
	}
	rows, err := s.db.QueryContext(ctx, objectCommitsSQL, ref.String(), afterEntry, limit)
	if err != nil {
		return nil, fmt.Errorf("store: object commits query: %w", err)
	}
	entries := make([]int64, 0, limit)
	for rows.Next() {
		var index int64
		if err = rows.Scan(&index); err != nil {
			break
		}
		if index < 0 {
			err = fmt.Errorf("negative stored entry index %d", index)
			break
		}
		entries = append(entries, index)
	}
	if err == nil {
		err = rows.Err()
	}
	closeErr := rows.Close()
	if err != nil {
		return nil, fmt.Errorf("store: object commits row: %w", err)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("store: close object commits: %w", closeErr)
	}
	return entries, nil
}
