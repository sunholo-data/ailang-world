package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// commitCarrying appends entry prev.Revision+1 carrying objs on top of prev.
func commitCarrying(t testing.TB, s *Store, prev World, objs ...Object) World {
	t.Helper()
	next, c := carryingCommit(prev, objs...)
	if err := s.Commit(boundedTestContext(t), c); err != nil {
		t.Fatalf("Commit entry %d: %v", next.Revision, err)
	}
	return next
}

func carryingCommit(prev World, objs ...Object) (World, Commit) {
	index := prev.Revision + 1
	entryHash := hashref.SumSHA256([]byte(fmt.Sprintf("membership-entry-%d", index)))
	next := World{
		Ref: hashref.SumSHA256([]byte(fmt.Sprintf("membership-world-%d", index))), Revision: index,
		StateRoot: hashref.SumSHA256([]byte(fmt.Sprintf("membership-state-%d", index))), LogHead: entryHash,
	}
	return next, Commit{
		ObservedHead: prev.Ref, Objects: objs, NextWorld: next,
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex: index, SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("membership-fn")),
				Interpreter:   hashref.SumSHA256([]byte("membership-interp")),
				PrevEntryHash: prev.LogHead, WrittenBy: "membership-test",
			},
			EntryHash: entryHash, TransitionRef: hashref.SumSHA256([]byte(fmt.Sprintf("membership-body-%d", index))),
		},
	}
}

func membershipRows(t *testing.T, s *Store) []string {
	t.Helper()
	rows, err := s.db.Query(`SELECT object_ref || '@' || entry_index FROM commit_objects ORDER BY 1`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			t.Fatal(err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

func membershipFor(t *testing.T, s *Store, ref hashref.HashRef) []int64 {
	t.Helper()
	rows, err := s.db.Query(`SELECT entry_index FROM commit_objects WHERE object_ref = ? ORDER BY entry_index`, ref.String())
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var index int64
		if err := rows.Scan(&index); err != nil {
			t.Fatal(err)
		}
		out = append(out, index)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// AC1: membership and the entry are one transaction, in both directions.
func TestCommitMembershipRollsBackWithCommit(t *testing.T) {
	for _, tc := range []struct{ name, trigger string }{
		{"membership-fails", `CREATE TRIGGER inject BEFORE INSERT ON commit_objects BEGIN SELECT RAISE(ABORT, 'injected membership failure'); END`},
		{"later-step-fails", `CREATE TRIGGER inject BEFORE INSERT ON store_heads BEGIN SELECT RAISE(ABORT, 'injected head failure'); END`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			w := seedGenesis(t, s)
			if _, err := s.db.Exec(`DELETE FROM store_heads`); err != nil { // force the INSERT branch of Step 5
				t.Fatal(err)
			}
			if _, err := s.db.Exec(tc.trigger); err != nil {
				t.Fatal(err)
			}
			_, c := carryingCommit(World{Ref: hashref.HashRef{}, LogHead: w.LogHead}, obj("rollback", "t/rb"))
			if err := s.Commit(boundedTestContext(t), c); err == nil {
				t.Fatal("Commit succeeded despite injected failure")
			}
			if _, ok, err := s.GetLogEntry(boundedTestContext(t), 1); err != nil || ok {
				t.Fatalf("entry 1 after failed commit: ok=%v err=%v", ok, err)
			}
			if rows := membershipRows(t, s); len(rows) != 0 {
				t.Fatalf("membership after failed commit: %v", rows)
			}
		})
	}
}

// AC2: the FK fence refuses membership naming an absent entry, even under OR IGNORE.
func TestCommitMembershipForeignKeys(t *testing.T) {
	s := openMem(t)
	o := obj("fk", "t/fk")
	if err := s.PutObject(boundedTestContext(t), o); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT OR IGNORE INTO commit_objects (object_ref, entry_index) VALUES (?, 999)`, o.Hash.String()); err == nil {
		t.Fatal("membership naming absent entry 999 was accepted")
	}
	if rows := membershipRows(t, s); len(rows) != 0 {
		t.Fatalf("rows = %v", rows)
	}
}

// AC3: a hash listed twice in one commit (equal input/output bytes) is one row.
func TestCommitMembershipDedupesWithinCommit(t *testing.T) {
	s := openMem(t)
	w := seedGenesis(t, s)
	in, out := obj("same-bytes", "invocation/input"), obj("same-bytes", "invocation/output")
	commitCarrying(t, s, w, in, out)
	if got := membershipFor(t, s, in.Hash); !reflect.DeepEqual(got, []int64{1}) {
		t.Fatalf("membership = %v, want [1]", got)
	}
}

// AC4: a resolved replay returns nil and writes no membership.
func TestCommitReplayDoesNotDuplicateMembership(t *testing.T) {
	s := openMem(t)
	c := journalCommitFixture(t, s, "replay")
	if _, _, err := s.AppendIntent(boundedTestContext(t), "replay", testCommitIntent("replay", c)); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(boundedTestContext(t), c); err != nil {
		t.Fatal(err)
	}
	before := membershipRows(t, s)
	if len(before) != 1 {
		t.Fatalf("membership after first commit = %v", before)
	}
	if err := s.Commit(boundedTestContext(t), c); err != nil {
		t.Fatalf("replay: %v", err)
	}
	if after := membershipRows(t, s); !reflect.DeepEqual(after, before) {
		t.Fatalf("replay changed membership: before=%v after=%v", before, after)
	}
}

// AC5: a stale-head commit writes no membership.
func TestConflictWritesNoMembership(t *testing.T) {
	s := openMem(t)
	g := seedGenesis(t, s)
	commitCarrying(t, s, g, obj("first", "t/first"))
	stale := obj("stale", "t/stale")
	_, c := carryingCommit(World{Ref: g.Ref, Revision: 1, LogHead: g.LogHead}, stale)
	if err := s.Commit(boundedTestContext(t), c); !IsConflict(err) {
		t.Fatalf("want ConflictError, got %v", err)
	}
	if got := membershipFor(t, s, stale.Hash); len(got) != 0 {
		t.Fatalf("conflicting commit recorded membership %v", got)
	}
}

// AC8: a v3 store (no membership table) is refused unmodified by writer and reader.
func TestVersionThreeStoreIsRefusedUnmodified(t *testing.T) {
	path := filepath.Join(t.TempDir(), "v3.db")
	db := rawDB(t, path)
	if _, err := db.Exec(schemaV3SQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 3)
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, open := range []struct {
		name string
		fn   func(string) (*Store, error)
	}{{"Open", Open}, {"OpenReadOnly", OpenReadOnly}} {
		s, err := open.fn(path)
		if s != nil {
			_ = s.Close()
		}
		var legacy *LegacySchemaVersionError
		if !errors.As(err, &legacy) || legacy.Found != 3 || legacy.Current != 4 {
			t.Fatalf("%s: err = %v, want legacy Found=3 Current=4", open.name, err)
		}
	}
	db = rawDB(t, path)
	afterNames, afterVersion := schemaState(t, db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != 3 || beforeVersion != 3 {
		t.Fatalf("v3 store changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
}

// AC16: an existing v4 store whose commit_objects was dropped is refused by
// both open modes before any DDL, and the table stays absent (never recreated empty).
func TestDroppedMembershipTableIsRefusedNotRecreated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dropped.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	carrier := obj("dropped-carrier", "t/dropped")
	commitCarrying(t, s, seedGenesis(t, s), carrier)
	if got := membershipFor(t, s, carrier.Hash); !reflect.DeepEqual(got, []int64{1}) {
		t.Fatalf("carrier before drop = %v", got)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	db := rawDB(t, path)
	if _, err := db.Exec(`DROP TABLE commit_objects`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, open := range []struct {
		name string
		fn   func(string) (*Store, error)
	}{{"Open", Open}, {"OpenReadOnly", OpenReadOnly}} {
		got, err := open.fn(path)
		if got != nil {
			_ = got.Close()
		}
		var integrity *SchemaIntegrityError
		if !errors.As(err, &integrity) || integrity.Table != "commit_objects" {
			t.Errorf("%s: err = %v, want SchemaIntegrityError for commit_objects", open.name, err)
		}
		db := rawDB(t, path)
		var tables int
		if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'commit_objects'`).Scan(&tables); err != nil {
			t.Fatal(err)
		}
		_ = db.Close()
		if tables != 0 {
			t.Fatalf("%s recreated commit_objects (count %d)", open.name, tables)
		}
	}
}

func objectCommits(t *testing.T, s *Store, ref hashref.HashRef) []int64 {
	t.Helper()
	got, err := s.ObjectCommits(boundedTestContext(t), ref, -1, MaxObjectCommitPage)
	if err != nil {
		t.Fatalf("ObjectCommits: %v", err)
	}
	return got
}

// AC6 + AC7: one, several and no carrying commits; PutObject-then-Commit.
func TestCommitRecordsMembershipPerCarryingCommit(t *testing.T) {
	s := openMem(t)
	w := seedGenesis(t, s)
	shared, once, stored, later := obj("shared", "t/shared"), obj("once", "t/once"), obj("stored-only", "t/stored"), obj("put-then-commit", "t/later")
	for _, o := range []Object{stored, later} {
		if err := s.PutObject(boundedTestContext(t), o); err != nil {
			t.Fatal(err)
		}
	}
	w = commitCarrying(t, s, w, shared)
	w = commitCarrying(t, s, w, once)
	w = commitCarrying(t, s, w, shared, later)
	_ = w
	for _, tc := range []struct {
		name string
		ref  hashref.HashRef
		want []int64
	}{
		{"multi", shared.Hash, []int64{1, 3}},
		{"single", once.Hash, []int64{2}},
		{"none", stored.Hash, []int64{}},
		{"put-then-commit", later.Hash, []int64{3}},
	} {
		if got := objectCommits(t, s, tc.ref); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: ObjectCommits = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// AC10: ascending order and strict-after continuation reassemble the one-shot read.
func TestObjectCommitsOrderAndContinuation(t *testing.T) {
	s := openMem(t)
	w := seedGenesis(t, s)
	shared := obj("carried-seven-times", "t/seven")
	var want []int64
	for i := 1; i <= 10; i++ {
		if i%3 == 0 {
			w = commitCarrying(t, s, w, obj(fmt.Sprintf("other-%d", i), "t/other"))
			continue
		}
		w = commitCarrying(t, s, w, shared)
		want = append(want, int64(i))
	}
	if len(want) != 7 {
		t.Fatalf("fixture carriers = %d", len(want))
	}
	var pages [][]int64
	var all []int64
	after := int64(-1)
	for round := 0; round < 5; round++ { // bounded: a non-advancing cursor reds, not hangs
		page, err := s.ObjectCommits(boundedTestContext(t), shared.Hash, after, 3)
		if err != nil {
			t.Fatal(err)
		}
		if len(page) == 0 {
			break
		}
		pages = append(pages, page)
		all = append(all, page...)
		after = page[len(page)-1]
	}
	if !reflect.DeepEqual(all, want) || len(pages) != 3 || len(pages[2]) != 1 {
		t.Fatalf("pages = %v, want %v in pages of 3,3,1", pages, want)
	}
	if !sort.SliceIsSorted(all, func(i, j int) bool { return all[i] < all[j] }) {
		t.Fatalf("not ascending: %v", all)
	}
}

// AC11: limit and cursor bounds.
func TestObjectCommitsValidation(t *testing.T) {
	s := openMem(t)
	ref := obj("v", "t/v").Hash
	ctx := boundedTestContext(t)
	for _, limit := range []int{0, MaxObjectCommitPage + 1} {
		var invalid *InvalidLimitError
		if _, err := s.ObjectCommits(ctx, ref, -1, limit); !errors.As(err, &invalid) {
			t.Errorf("limit %d: want InvalidLimitError, got %v", limit, err)
		}
	}
	var cursor *InvalidObjectCommitCursorError
	if _, err := s.ObjectCommits(ctx, ref, -2, 1); !errors.As(err, &cursor) {
		t.Errorf("afterEntry -2: want InvalidObjectCommitCursorError, got %v", err)
	}
	var badRef *InvalidRefError
	if _, err := s.ObjectCommits(ctx, hashref.HashRef{}, -1, 1); !errors.As(err, &badRef) {
		t.Errorf("zero ref: want InvalidRefError, got %v", err)
	}
	if _, err := s.ObjectCommits(ctx, ref, -1, MaxObjectCommitPage); err != nil {
		t.Errorf("limit %d refused: %v", MaxObjectCommitPage, err)
	}
}

// AC9: a v4 store opened read-only answers the same membership; AC6 existence:
// every returned entry is a stored log entry.
func TestObjectCommitsReadOnlyMatchesWriterAndEntriesExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "membership.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	w := seedGenesis(t, s)
	shared := obj("ro-shared", "t/ro")
	w = commitCarrying(t, s, w, shared)
	w = commitCarrying(t, s, w, obj("ro-other", "t/ro"))
	commitCarrying(t, s, w, shared)
	writerView := objectCommits(t, s, shared.Hash)
	ro, err := OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	readerView := objectCommits(t, ro, shared.Hash)
	if !reflect.DeepEqual(readerView, []int64{1, 3}) || !reflect.DeepEqual(readerView, writerView) {
		t.Fatalf("reader %v writer %v, want [1 3]", readerView, writerView)
	}
	for _, index := range readerView {
		if _, ok, err := ro.GetLogEntry(boundedTestContext(t), index); err != nil || !ok {
			t.Fatalf("membership names entry %d: ok=%v err=%v", index, ok, err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

// AC11 (store damage): a v4 store missing the table answers an error, never [].
func TestObjectCommitsMissingTableIsAnError(t *testing.T) {
	s := openMem(t)
	if _, err := s.db.Exec(`DROP TABLE commit_objects`); err != nil {
		t.Fatal(err)
	}
	if got, err := s.ObjectCommits(boundedTestContext(t), obj("x", "t/x").Hash, -1, 1); err == nil {
		t.Fatalf("missing table answered %v, nil", got)
	}
}

// AC13: N=10,000 timing. Opt-in: MEMBERSHIP_TIMING=1.
func TestObjectCommitsTimingAt10k(t *testing.T) {
	if testing.Short() || !timingEnabled() {
		t.Skip("set MEMBERSHIP_TIMING=1")
	}
	path := filepath.Join(t.TempDir(), "timing.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	w := seedGenesis(t, s)
	hot := obj("hot", "t/hot")
	const n = 10000
	commitTimes := make([]time.Duration, 0, n)
	for i := 1; i <= n; i++ {
		objs := []Object{obj(fmt.Sprintf("in-%d", i), "invocation/input"), obj(fmt.Sprintf("out-%d", i), "invocation/output"), obj(fmt.Sprintf("rec-%d", i), "invocation/record")}
		if i%100 == 0 {
			objs = append(objs, hot)
		}
		next, c := carryingCommit(w, objs...)
		start := time.Now()
		if err := s.Commit(boundedTestContext(t), c); err != nil {
			t.Fatal(err)
		}
		commitTimes = append(commitTimes, time.Since(start))
		w = next
	}
	cold := obj(fmt.Sprintf("out-%d", n/2), "invocation/output").Hash
	stored := obj("never-committed", "t/none")
	if err := s.PutObject(boundedTestContext(t), stored); err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	for _, tc := range []struct {
		name  string
		ref   hashref.HashRef
		limit int
		want  int
	}{{"single", cold, 101, 1}, {"hot-100", hot.Hash, 101, 100}, {"none", stored.Hash, 101, 0}} {
		samples := make([]time.Duration, 0, 200)
		for i := 0; i < 200; i++ {
			start := time.Now()
			got, err := s.ObjectCommits(ctx, tc.ref, -1, tc.limit)
			samples = append(samples, time.Since(start))
			if err != nil || len(got) != tc.want {
				t.Fatalf("%s: %d rows err=%v", tc.name, len(got), err)
			}
		}
		t.Logf("read %s p50=%s p95=%s", tc.name, pct(samples, 50), pct(samples, 95))
	}
	t.Logf("commit (3 objects, file-backed, N=%d) p50=%s p95=%s", n, pct(commitTimes, 50), pct(commitTimes, 95))
}

func timingEnabled() bool { return os.Getenv("MEMBERSHIP_TIMING") == "1" }

func pct(samples []time.Duration, p int) time.Duration {
	sorted := append([]time.Duration(nil), samples...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted[(len(sorted)-1)*p/100]
}
