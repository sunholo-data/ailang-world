package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"modernc.org/sqlite"
	"reflect"
	"strings"
	"testing"
	"time"
)

type latestLogReader interface {
	LogEntriesLatest(context.Context, int) ([]LogEntry, error)
}

func latestReader(t *testing.T, s *Store) latestLogReader {
	t.Helper()
	r, ok := any(s).(latestLogReader)
	if !ok {
		t.Fatal("Store.LogEntriesLatest is missing")
	}
	return r
}
func TestLogEntriesLatestNewestFirst(t *testing.T) {
	s, path := logAfterFixture(t, []int64{0, 1, 5})
	r := latestReader(t, s)
	for _, tc := range []struct {
		limit int
		want  []int64
	}{{10, []int64{5, 1, 0}}, {1, []int64{5}}} {
		got, err := r.LogEntriesLatest(boundedTestContext(t), tc.limit)
		if err != nil || !reflect.DeepEqual(logIndexes(got), tc.want) {
			t.Fatalf("limit %d: indexes %v want %v err=%v", tc.limit, logIndexes(got), tc.want, err)
		}
		for _, e := range got {
			want, ok, err := s.GetLogEntry(boundedTestContext(t), e.Header.EntryIndex)
			if err != nil || !ok || !reflect.DeepEqual(e, want) {
				t.Fatalf("row differs from GetLogEntry: %+v %+v %v", e, want, err)
			}
		}
	}
	empty, _ := logAfterFixture(t, nil)
	got, err := latestReader(t, empty).LogEntriesLatest(boundedTestContext(t), 10)
	if err != nil || len(got) != 0 {
		t.Fatalf("empty store: %v %v", got, err)
	}
	// Observe the executed statement through the existing query driver.
	if err := s.db.Close(); err != nil {
		t.Fatal(err)
	}
	driver := &logQueryDriver{Driver: &sqlite.Driver{}}
	name := fmt.Sprintf("log-latest-query-%d", logQueryDriverID.Add(1))
	sql.Register(name, driver)
	db, err := sql.Open(name, path)
	if err != nil {
		t.Fatal(err)
	}
	s.db = db
	s.db.SetMaxOpenConns(1)
	if _, err := r.LogEntriesLatest(boundedTestContext(t), 10); err != nil {
		t.Fatal(err)
	}
	driver.mu.Lock()
	queries := append([]string(nil), driver.queries...)
	driver.mu.Unlock()
	if len(queries) != 1 {
		t.Fatalf("want one indexed query, got %v", queries)
	}
	q := strings.Join(strings.Fields(strings.ToUpper(queries[0])), " ")
	if !strings.Contains(q, "ORDER BY ENTRY_INDEX DESC LIMIT ?") {
		t.Fatalf("not newest-first bounded SQL: %s", q)
	}
}
func TestLogEntriesLatestBounds(t *testing.T) {
	s, _ := logAfterFixture(t, []int64{0, 1, 5})
	r := latestReader(t, s)
	ctx := boundedTestContext(t)
	for _, limit := range []int{0, 501} {
		var bad *InvalidLimitError
		_, err := r.LogEntriesLatest(ctx, limit)
		if !errors.As(err, &bad) || bad.Op != "LogEntriesLatest" || bad.Limit != limit || bad.Max != 500 {
			t.Fatalf("limit %d: %v", limit, err)
		}
	}
	for _, c := range []context.Context{nil, context.Background()} {
		if _, err := r.LogEntriesLatest(c, 1); !errors.Is(err, ErrNoDeadline) {
			t.Fatalf("no deadline: %v", err)
		}
	}
	expired, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	if _, err := r.LogEntriesLatest(expired, 1); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expired: %v", err)
	}
	canceled, stop := context.WithTimeout(context.Background(), 30*time.Second)
	stop()
	if _, err := r.LogEntriesLatest(canceled, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled: %v", err)
	}
	Quarantine(s)
	for _, c := range []context.Context{ctx, nil} {
		if _, err := r.LogEntriesLatest(c, 1); !errors.Is(err, ErrQuarantined) {
			t.Fatalf("quarantine: %v", err)
		}
	}
}
