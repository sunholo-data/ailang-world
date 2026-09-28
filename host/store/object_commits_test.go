package store

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// commitCarrying appends entry prev.Revision+1 carrying objs on top of prev.
func commitCarrying(t testing.TB, s *Store, prev World, objs ...Object) World {
	t.Helper()
	next, c := carryingCommit(prev, objs...)
	if err := s.Commit(c); err != nil {
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
			if err := s.Commit(c); err == nil {
				t.Fatal("Commit succeeded despite injected failure")
			}
			if _, ok, err := s.GetLogEntry(context.Background(), 1); err != nil || ok {
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
	if err := s.PutObject(o); err != nil {
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
	if _, _, err := s.AppendIntent("replay", testCommitIntent("replay", c)); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(c); err != nil {
		t.Fatal(err)
	}
	before := membershipRows(t, s)
	if len(before) != 1 {
		t.Fatalf("membership after first commit = %v", before)
	}
	if err := s.Commit(c); err != nil {
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
	if err := s.Commit(c); !IsConflict(err) {
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
