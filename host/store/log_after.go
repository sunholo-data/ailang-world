package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// MaxLogEntryPage bounds each keyset log read.
const MaxLogEntryPage = 500

// InvalidLogCursorError reports a cursor below the genesis sentinel.
type InvalidLogCursorError struct{ After int64 }

func (e *InvalidLogCursorError) Error() string {
	return fmt.Sprintf("store: invalid log cursor %d: want -1 (start) or an entry index", e.After)
}

const logEntriesAfterSQL = `SELECT entry_index, entry_hash_ref, semantics_epoch, transition_fn_ref, interpreter_ref,
 prev_entry_hash_ref, written_by, transition_ref
 FROM log_entries WHERE entry_index > ? ORDER BY entry_index LIMIT ?`

// LogEntriesAfter returns at most limit committed rows, in ascending index order,
// strictly after after (-1 starts at genesis). Keyset reads cross index gaps.
// The caller must supply a finite deadline; quarantined stores refuse all reads.
func (s *Store) LogEntriesAfter(ctx context.Context, after int64, limit int) ([]LogEntry, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if limit < 1 || limit > MaxLogEntryPage {
		return nil, &InvalidLimitError{Op: "LogEntriesAfter", Limit: limit, Max: MaxLogEntryPage}
	}
	if after < -1 {
		return nil, &InvalidLogCursorError{After: after}
	}
	rows, err := s.db.QueryContext(ctx, logEntriesAfterSQL, after, limit)
	if err != nil {
		return nil, fmt.Errorf("store: log entries after query: %w", err)
	}
	defer rows.Close()
	entries := make([]LogEntry, 0, limit)
	for rows.Next() {
		entry, err := scanLogAfterRow(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: log entries after rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("store: close log entries after: %w", err)
	}
	return entries, nil
}

// Keep this local parser faithful to GetLogEntry's representation and error
// labels. The equality and malformed-row tests bind the two algorithms without
// refactoring the frozen store.go surface (planner decision PD4).
func scanLogAfterRow(rows *sql.Rows) (LogEntry, error) {
	var index, epoch int64
	var entryHashText, fnText, interpText, prevText, writtenBy, transRefText string
	if err := rows.Scan(&index, &entryHashText, &epoch, &fnText, &interpText, &prevText, &writtenBy, &transRefText); err != nil {
		return LogEntry{}, fmt.Errorf("store: get log entry %d: %w", index, err)
	}
	entryHash, err := hashref.Parse(entryHashText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("store: log entry %d hash: %w", index, err)
	}
	fn, err := hashref.Parse(fnText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("store: log entry %d transitionFn: %w", index, err)
	}
	interp, err := hashref.Parse(interpText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("store: log entry %d interpreter: %w", index, err)
	}
	prev, err := hashref.Parse(prevText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("store: log entry %d prevEntryHash: %w", index, err)
	}
	transRef, err := hashref.Parse(transRefText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("store: log entry %d transitionRef: %w", index, err)
	}
	return LogEntry{
		Header: LogHeader{
			EntryIndex:     index,
			SemanticsEpoch: epoch,
			TransitionFn:   fn,
			Interpreter:    interp,
			PrevEntryHash:  prev,
			WrittenBy:      writtenBy,
		},
		EntryHash:     entryHash,
		TransitionRef: transRef,
	}, nil
}
