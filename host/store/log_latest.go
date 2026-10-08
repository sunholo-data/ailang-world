package store

import (
	"context"
	"fmt"
)

const logEntriesLatestSQL = `SELECT entry_index, entry_hash_ref, semantics_epoch, transition_fn_ref, interpreter_ref,
 prev_entry_hash_ref, written_by, transition_ref
 FROM log_entries ORDER BY entry_index DESC LIMIT ?`

// LogEntriesLatest returns at most limit committed rows, newest first, crossing
// index gaps. A finite deadline is required; quarantined stores refuse reads.
func (s *Store) LogEntriesLatest(ctx context.Context, limit int) ([]LogEntry, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if limit < 1 || limit > MaxLogEntryPage {
		return nil, &InvalidLimitError{Op: "LogEntriesLatest", Limit: limit, Max: MaxLogEntryPage}
	}
	rows, err := s.db.QueryContext(ctx, logEntriesLatestSQL, limit)
	if err != nil {
		return nil, fmt.Errorf("store: log entries latest query: %w", err)
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
		return nil, fmt.Errorf("store: log entries latest rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("store: close log entries latest: %w", err)
	}
	return entries, nil
}
