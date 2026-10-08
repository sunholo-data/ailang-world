package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"modernc.org/sqlite"
)

func logAfterFixture(t *testing.T, indexes []int64) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "log.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	w := seedGenesis(t, s)
	for _, i := range indexes {
		prev := w
		prev.Revision = i - 1
		next, c := carryingCommit(prev)
		c.Entry.Header.SemanticsEpoch = 7
		c.Entry.Header.WrittenBy = "log-after fixture <>&"
		if err := s.Commit(boundedTestContext(t), c); err != nil {
			t.Fatalf("commit %d: %v", i, err)
		}
		w = next
	}
	return s, path
}
func logIndexes(entries []LogEntry) []int64 {
	out := make([]int64, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Header.EntryIndex)
	}
	return out
}
func TestLogEntriesAfterCrossesGap(t *testing.T) {
	s, _ := logAfterFixture(t, []int64{0, 1, 5})
	for _, tc := range []struct {
		after int64
		want  []int64
	}{{-1, []int64{0, 1, 5}}, {1, []int64{5}}, {5, []int64{}}} {
		got, err := s.LogEntriesAfter(boundedTestContext(t), tc.after, 100)
		if err != nil || !reflect.DeepEqual(logIndexes(got), tc.want) {
			t.Fatalf("after %d: got %v err=%v want %v (exclusive cursor)", tc.after, logIndexes(got), err, tc.want)
		}
	}
	var indexes []int64
	for i := int64(0); i < 207; i++ {
		if i != 2 && i != 3 && i != 4 {
			indexes = append(indexes, i)
		}
	}
	large, _ := logAfterFixture(t, indexes)
	for _, limit := range []int{1, 100, 500} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			after := int64(-1)
			var got []int64
			for n := 0; n <= len(indexes); n++ {
				page, err := large.LogEntriesAfter(boundedTestContext(t), after, limit)
				if err != nil {
					t.Fatal(err)
				}
				if len(page) > limit {
					t.Fatalf("limit %d exceeded", limit)
				}
				if len(page) == 0 {
					break
				}
				for _, e := range page {
					if e.Header.EntryIndex <= after {
						t.Fatal("ordering/cursor not strictly increasing")
					}
					after = e.Header.EntryIndex
					got = append(got, after)
				}
			}
			if !reflect.DeepEqual(got, indexes) {
				t.Fatalf("limit %d paging: got %v want %v", limit, got, indexes)
			}
		})
	}
}
func TestLogEntriesAfterBounds(t *testing.T) {
	s, _ := logAfterFixture(t, []int64{0, 1, 5})
	ctx := boundedTestContext(t)
	for _, limit := range []int{0, 501} {
		var invalid *InvalidLimitError
		_, err := s.LogEntriesAfter(ctx, -1, limit)
		if !errors.As(err, &invalid) || invalid.Op != "LogEntriesAfter" || invalid.Limit != limit || invalid.Max != 500 {
			t.Fatalf("limit %d: %v", limit, err)
		}
	}
	if _, err := s.LogEntriesAfter(ctx, -2, 1); err == nil {
		t.Fatal("cursor below -1 accepted")
	}
	for _, ctx := range []context.Context{nil, context.Background()} {
		if _, err := s.LogEntriesAfter(ctx, -1, 1); !errors.Is(err, ErrNoDeadline) {
			t.Fatalf("missing deadline: %v", err)
		}
	}
	expired, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	if _, err := s.LogEntriesAfter(expired, -1, 1); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expired context: %v", err)
	}
	canceled, cancel2 := context.WithTimeout(context.Background(), time.Second)
	cancel2()
	if _, err := s.LogEntriesAfter(canceled, -1, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled context: %v", err)
	}
	Quarantine(s)
	for _, ctx := range []context.Context{ctx, nil} {
		if _, err := s.LogEntriesAfter(ctx, -1, 1); !errors.Is(err, ErrQuarantined) {
			t.Fatalf("quarantine: %v", err)
		}
	}
}

// A driver wrapper observes the actual statements, not the source text. The
// underlying sqlite dependency is already used by Store; no dependency is added.
type logQueryDriver struct {
	driver.Driver
	mu      sync.Mutex
	queries []string
}
type logQueryConn struct {
	driver.Conn
	owner *logQueryDriver
}

func (d *logQueryDriver) Open(name string) (driver.Conn, error) {
	c, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &logQueryConn{c, d}, nil
}
func (c *logQueryConn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	c.owner.mu.Lock()
	c.owner.queries = append(c.owner.queries, q)
	c.owner.mu.Unlock()
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, q, args)
}

var logQueryDriverID atomic.Uint64

func TestLogEntriesAfterMatchesGetLogEntry(t *testing.T) {
	s, path := logAfterFixture(t, []int64{0, 1, 5})
	var want []LogEntry
	for _, i := range []int64{0, 1, 5} {
		e, ok, err := s.GetLogEntry(boundedTestContext(t), i)
		if err != nil || !ok {
			t.Fatal(err)
		}
		want = append(want, e)
	}
	if err := s.db.Close(); err != nil {
		t.Fatal(err)
	}
	d := &logQueryDriver{Driver: &sqlite.Driver{}}
	name := fmt.Sprintf("log-after-query-%d", logQueryDriverID.Add(1))
	sql.Register(name, d)
	db, err := sql.Open(name, path)
	if err != nil {
		t.Fatal(err)
	}
	s.db = db
	s.db.SetMaxOpenConns(1)
	got, err := s.LogEntriesAfter(boundedTestContext(t), -1, 100)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("row representation mismatch:\ngot=%+v\nwant=%+v\nerr=%v", got, want, err)
	}
	d.mu.Lock()
	queries := append([]string(nil), d.queries...)
	d.mu.Unlock()
	if len(queries) != 1 {
		t.Fatalf("want one SQL query, got %v", queries)
	}
	q := strings.Join(strings.Fields(strings.ToUpper(queries[0])), " ")
	if !strings.Contains(q, "WHERE ENTRY_INDEX > ? ORDER BY ENTRY_INDEX LIMIT ?") {
		t.Fatalf("not ordered exclusive bounded keyset SQL: %s", q)
	}
	for _, column := range []string{"entry_hash_ref", "transition_fn_ref", "interpreter_ref", "prev_entry_hash_ref", "transition_ref"} {
		t.Run(column, func(t *testing.T) {
			// Only the column name (closed test table) is interpolated; values are bound.
			var original string
			if err := s.db.QueryRow("SELECT " + column + " FROM log_entries WHERE entry_index=1").Scan(&original); err != nil {
				t.Fatal(err)
			}
			if _, err := s.db.Exec("UPDATE log_entries SET "+column+"=? WHERE entry_index=1", "not-a-hash"); err != nil {
				t.Fatal(err)
			}
			_, _, getterErr := s.GetLogEntry(boundedTestContext(t), 1)
			_, afterErr := s.LogEntriesAfter(boundedTestContext(t), 0, 1)
			if getterErr == nil || afterErr == nil || getterErr.Error() != afterErr.Error() {
				t.Fatalf("malformed %s: getter=%v keyset=%v", column, getterErr, afterErr)
			}
			if _, err := s.db.Exec("UPDATE log_entries SET "+column+"=? WHERE entry_index=1", original); err != nil {
				t.Fatal(err)
			}
		})
	}
}
