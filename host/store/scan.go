package store

import (
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// MaxIntegrityScanPage is the kernel-owned upper bound for one scan query.
const MaxIntegrityScanPage = 1000

// InvalidLimitError reports a scan request outside the kernel-owned bound.
type InvalidLimitError struct {
	Op         string
	Limit, Max int
}

func (e *InvalidLimitError) Error() string {
	return fmt.Sprintf("store: %s: limit %d is outside 1..%d", e.Op, e.Limit, e.Max)
}

func IsInvalidLimit(err error) bool {
	var invalid *InvalidLimitError
	return errors.As(err, &invalid)
}

type UnreadableRow struct {
	Table  string
	Index  int64
	Ref    string
	Field  string
	Reason string
}

type ScanPage struct {
	Rows      []UnreadableRow
	Scanned   int
	NextIndex int64
	NextRef   string
	Done      bool
}

func invalidScanLimit(op string, limit int) error {
	if limit < 1 || limit > MaxIntegrityScanPage {
		return &InvalidLimitError{Op: op, Limit: limit, Max: MaxIntegrityScanPage}
	}
	return nil
}

// ScanUnreadableLog scans a bounded keyset page. Columns are parsed in their
// declared order and only the first invalid field in each row is reported.
func (s *Store) ScanUnreadableLog(fromIndex int64, limit int) (ScanPage, error) {
	const op = "ScanUnreadableLog"
	if err := invalidScanLimit(op, limit); err != nil {
		return ScanPage{}, err
	}
	rows, err := s.db.Query(`SELECT entry_index, entry_hash_ref, transition_fn_ref,
		interpreter_ref, prev_entry_hash_ref, transition_ref
		FROM log_entries WHERE entry_index >= ? ORDER BY entry_index LIMIT ?`, fromIndex, limit)
	if err != nil {
		return ScanPage{}, fmt.Errorf("store: scan unreadable log: %w", err)
	}
	defer rows.Close()
	page := ScanPage{NextIndex: fromIndex}
	for rows.Next() {
		var index int64
		var texts [5]string
		if err := rows.Scan(&index, &texts[0], &texts[1], &texts[2], &texts[3], &texts[4]); err != nil {
			return ScanPage{}, fmt.Errorf("store: scan unreadable log row: %w", err)
		}
		page.Scanned++
		page.NextIndex = index + 1
		fields := [...]string{"entryHash", "transitionFn", "interpreter", "prevEntryHash", "transitionRef"}
		for i, text := range texts {
			if _, err := hashref.Parse(text); err != nil {
				page.Rows = append(page.Rows, UnreadableRow{
					Table: "log_entries", Index: index, Field: fields[i], Reason: err.Error(),
				})
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		return ScanPage{}, fmt.Errorf("store: scan unreadable log rows: %w", err)
	}
	page.Done = page.Scanned < limit
	return page, nil
}

// ScanUnreadableWorlds scans lexicographically by the explicit TEXT primary
// key. It never depends on OFFSET or SQLite rowids.
func (s *Store) ScanUnreadableWorlds(afterRef string, limit int) (ScanPage, error) {
	const op = "ScanUnreadableWorlds"
	if err := invalidScanLimit(op, limit); err != nil {
		return ScanPage{}, err
	}
	rows, err := s.db.Query(`SELECT world_ref, state_root, log_head
		FROM worlds WHERE world_ref > ? ORDER BY world_ref LIMIT ?`, afterRef, limit)
	if err != nil {
		return ScanPage{}, fmt.Errorf("store: scan unreadable worlds: %w", err)
	}
	defer rows.Close()
	page := ScanPage{NextRef: afterRef}
	for rows.Next() {
		var texts [3]string
		if err := rows.Scan(&texts[0], &texts[1], &texts[2]); err != nil {
			return ScanPage{}, fmt.Errorf("store: scan unreadable world row: %w", err)
		}
		page.Scanned++
		page.NextRef = texts[0]
		fields := [...]string{"worldRef", "stateRoot", "logHead"}
		for i, text := range texts {
			if _, err := hashref.Parse(text); err != nil {
				page.Rows = append(page.Rows, UnreadableRow{
					Table: "worlds", Ref: texts[0], Field: fields[i], Reason: err.Error(),
				})
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		return ScanPage{}, fmt.Errorf("store: scan unreadable world rows: %w", err)
	}
	page.Done = page.Scanned < limit
	return page, nil
}
