package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func semanticObject(t *testing.T, s *Store, id, label string) Object {
	t.Helper()
	payload := []byte("sid-payload-" + label)
	o := Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("sid-iface")),
		SemanticID: id, Provenance: "sid-test-" + label, Payload: payload,
	}
	if err := s.PutObject(context.Background(), o); err != nil {
		t.Fatalf("PutObject %s: %v", label, err)
	}
	return o
}

func hashesOf(objects []Object) []string {
	out := make([]string, len(objects))
	for i, o := range objects {
		out[i] = o.Hash.String()
	}
	return out
}

// seedSemanticIDs inserts seven objects sharing one semanticId in DESCENDING
// hash order (so insertion order and hash order disagree), plus one decoy with
// a different id. It returns the shared id's hashes in ascending order.
func seedSemanticIDs(t *testing.T, s *Store) (string, []string) {
	t.Helper()
	const id = "world/mission/incident/shared"
	var want []string
	var labels []string
	for i := 0; i < 7; i++ {
		labels = append(labels, fmt.Sprintf("n%d", i))
	}
	sort.Slice(labels, func(i, j int) bool {
		return hashref.SumSHA256([]byte("sid-payload-"+labels[i])).String() >
			hashref.SumSHA256([]byte("sid-payload-"+labels[j])).String()
	})
	for _, label := range labels {
		want = append(want, semanticObject(t, s, id, label).Hash.String())
	}
	semanticObject(t, s, id+"-decoy", "decoy")
	semanticObject(t, s, "world/mission/incident/shar", "prefix-decoy")
	sort.Strings(want)
	return id, want
}

// TestObjectsBySemanticIDNonUniqueOrderedAndPaged pins AC-1..AC-3: every object
// carrying the id is returned and no other, in ascending hash order that does
// not follow insertion order, and keyset pages partition the set exactly.
func TestObjectsBySemanticIDNonUniqueOrderedAndPaged(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "sid.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx := context.Background()
	id, want := seedSemanticIDs(t, s)

	all, err := s.ObjectsBySemanticID(ctx, id, "", MaxSemanticIDPage)
	if err != nil {
		t.Fatal(err)
	}
	if got := hashesOf(all); !reflect.DeepEqual(got, want) {
		t.Fatalf("full page = %v, want ascending %v", got, want)
	}
	for _, o := range all {
		if o.SemanticID != id || o.Payload != nil || !strings.HasPrefix(o.Provenance, "sid-test-") {
			t.Fatalf("item %+v: want semanticId %q, nil payload, stored provenance", o, id)
		}
	}
	again, err := s.ObjectsBySemanticID(ctx, id, "", MaxSemanticIDPage)
	if err != nil || !reflect.DeepEqual(hashesOf(again), want) {
		t.Fatalf("second call = %v (%v), want identical %v", hashesOf(again), err, want)
	}

	var paged []string
	after := ""
	for pages := 0; ; pages++ {
		if pages > len(want) {
			t.Fatalf("paging did not terminate: %v", paged)
		}
		page, err := s.ObjectsBySemanticID(ctx, id, after, 3)
		if err != nil {
			t.Fatal(err)
		}
		if len(page) > 3 {
			t.Fatalf("page after %q has %d items, want at most the limit 3", after, len(page))
		}
		paged = append(paged, hashesOf(page)...)
		if len(page) < 3 {
			break
		}
		after = page[len(page)-1].Hash.String()
	}
	if !reflect.DeepEqual(paged, want) {
		t.Fatalf("pages concatenated = %v, want %v", paged, want)
	}

	none, err := s.ObjectsBySemanticID(ctx, "world/mission/incident/never", "", 10)
	if err != nil || none == nil || len(none) != 0 {
		t.Fatalf("unknown id = %#v, %v; want empty non-nil slice", none, err)
	}
	for _, bad := range []int{0, MaxSemanticIDPage + 1} {
		if _, err := s.ObjectsBySemanticID(ctx, id, "", bad); !IsInvalidLimit(err) {
			t.Fatalf("limit %d: err = %v, want InvalidLimitError", bad, err)
		}
	}
}

// TestObjectsBySemanticIDUsesIndex pins AC-7 and AC-8: a fresh store carries
// the index and the planner searches it; unguarded SQL gives the same result.
func TestObjectsBySemanticIDUsesIndex(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "sid-index.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	id, want := seedSemanticIDs(t, s)

	var plan []string
	rows, err := s.db.Query("EXPLAIN QUERY PLAN "+objectsBySemanticIDSQL, id, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var a, b, c int
		var detail string
		if err := rows.Scan(&a, &b, &c, &detail); err != nil {
			t.Fatal(err)
		}
		plan = append(plan, detail)
	}
	rows.Close()
	joined := strings.Join(plan, " | ")
	if !strings.Contains(joined, "USING INDEX objects_by_semantic_id") || strings.Contains(joined, "TEMP B-TREE") {
		t.Fatalf("query plan = %q, want an ordered search of objects_by_semantic_id with no sort step", joined)
	}
	indexed, err := rawSemanticObjects(s.db, id, "", MaxSemanticIDPage)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := s.db.Exec("DROP INDEX objects_by_semantic_id"); err != nil {
		t.Fatal(err)
	}
	unindexed, err := rawSemanticObjects(s.db, id, "", MaxSemanticIDPage)
	if err != nil || !reflect.DeepEqual(unindexed, indexed) || !reflect.DeepEqual(hashesOf(unindexed), want) {
		t.Fatalf("without the index = %v (%v), want byte-identical rows %v", unindexed, err, indexed)
	}
}

// TestSemanticIDIndexBuiltOnWritableReopen pins B4's no-version-bump path: a
// version-3 store that lacks the index gains it on its next writable Open.
func TestSemanticIDIndexBuiltOnWritableReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sid-reopen.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec("DROP INDEX objects_by_semantic_id"); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatalf("reopen v3 store without the index: %v", err)
	}
	defer s.Close()
	var n int
	if err := s.db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='index' AND name='objects_by_semantic_id'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("index count after writable reopen = %d (%v), want 1", n, err)
	}
}

// rawSemanticObjects exercises result equivalence without the serving guard.
func rawSemanticObjects(db *sql.DB, id, after string, limit int) ([]Object, error) {
	rows, err := db.Query(objectsBySemanticIDSQL, id, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Object
	for rows.Next() {
		var h, iface, semantic, prov string
		if err := rows.Scan(&h, &iface, &semantic, &prov); err != nil {
			return nil, err
		}
		hr, err := hashref.Parse(h)
		if err != nil {
			return nil, err
		}
		ir, err := hashref.Parse(iface)
		if err != nil {
			return nil, err
		}
		out = append(out, Object{Hash: hr, InterfaceHash: ir, SemanticID: semantic, Provenance: prov})
	}
	return out, rows.Err()
}

func TestLookupIndexGuardReadOnlyThenProvision(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guard.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	o := semanticObject(t, s, "guard/name", "one")
	if _, err := s.db.Exec("DROP INDEX objects_by_semantic_id"); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	ro, err := OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	objectsBySemanticIDBeforeQuery = func() { calls++ }
	t.Cleanup(func() { objectsBySemanticIDBeforeQuery = nil })
	_, err = ro.ObjectsBySemanticID(context.Background(), "guard/name", "", 10)
	var unavailable *LookupIndexUnavailableError
	if !errors.As(err, &unavailable) || calls != 0 {
		t.Fatalf("unindexed read-only lookup = %v, hook calls=%d; want typed refusal and zero calls", err, calls)
	}
	ro.Close()
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s.Close()
	ro, err = OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	got, err := ro.ObjectsBySemanticID(context.Background(), "guard/name", "", 10)
	if err != nil || len(got) != 1 || got[0].Hash != o.Hash {
		t.Fatalf("reopened lookup=%v (%v), want %s", got, err, o.Hash)
	}
	rows, err := ro.db.Query("EXPLAIN QUERY PLAN "+objectsBySemanticIDSQL, "guard/name", "", 10)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	plan := ""
	for rows.Next() {
		var a, b, c int
		var detail string
		if err := rows.Scan(&a, &b, &c, &detail); err != nil {
			t.Fatal(err)
		}
		plan += detail
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plan, "USING INDEX objects_by_semantic_id") || strings.Contains(plan, "TEMP B-TREE") {
		t.Fatalf("read-only query plan=%q", plan)
	}
}

func TestLookupIndexProvisionFailureIsSurfaced(t *testing.T) {
	before := lookupIndexProvisionDeadline
	lookupIndexProvisionDeadline = -time.Nanosecond
	t.Cleanup(func() { lookupIndexProvisionDeadline = before })
	s, err := Open(filepath.Join(t.TempDir(), "deadline.db"))
	if s != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Open = %v, %v; want nil store and deadline", s, err)
	}
}

func TestLookupIndexIncompatibleIsRefused(t *testing.T) {
	defs := []struct{ name, ddl string }{
		{"wrong-column-order", "CREATE INDEX objects_by_semantic_id ON objects(hash_ref, semantic_id)"},
		{"partial", "CREATE INDEX objects_by_semantic_id ON objects(semantic_id, hash_ref) WHERE semantic_id <> ''"},
		{"nocase", "CREATE INDEX objects_by_semantic_id ON objects(semantic_id COLLATE NOCASE, hash_ref)"},
		{"descending", "CREATE INDEX objects_by_semantic_id ON objects(semantic_id, hash_ref DESC)"},
	}
	for _, tc := range defs {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "incompatible.db")
			s, err := Open(path)
			if err != nil {
				t.Fatal(err)
			}
			semanticObject(t, s, "x", "one")
			if _, err := s.db.Exec("DROP INDEX objects_by_semantic_id"); err != nil {
				t.Fatal(err)
			}
			if _, err := s.db.Exec(tc.ddl); err != nil {
				t.Fatal(err)
			}
			var before string
			if err := s.db.QueryRow("SELECT sql FROM sqlite_master WHERE name='objects_by_semantic_id'").Scan(&before); err != nil {
				t.Fatal(err)
			}
			s.Close()
			// These column names and positions come from modernc.org/sqlite, not sqlite3 CLI.
			db, err := sql.Open("sqlite", path)
			if err != nil {
				t.Fatal(err)
			}
			for _, probe := range []struct {
				query   string
				columns []string
			}{
				{"PRAGMA index_list('objects')", []string{"seq", "name", "unique", "origin", "partial"}},
				{"PRAGMA index_xinfo('objects_by_semantic_id')", []string{"seqno", "cid", "name", "desc", "coll", "key"}},
			} {
				rows, err := db.Query(probe.query)
				if err != nil {
					t.Fatal(err)
				}
				cols, err := rows.Columns()
				rows.Close()
				if err != nil || !reflect.DeepEqual(cols, probe.columns) {
					t.Fatalf("%s columns=%v (%v), want %v", probe.query, cols, err, probe.columns)
				}
			}
			db.Close()
			writable, openErr := Open(path)
			if writable != nil {
				writable.Close()
			}
			db, err = sql.Open("sqlite", path)
			if err != nil {
				t.Fatal(err)
			}
			var after string
			if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE name='objects_by_semantic_id'").Scan(&after); err != nil {
				t.Fatal(err)
			}
			db.Close()
			if after != before {
				t.Fatalf("index SQL changed: before %q, after %q", before, after)
			}
			if writable != nil || openErr == nil {
				t.Fatalf("incompatible writable Open=%v (%v), want refusal", writable, openErr)
			}
			ro, err := OpenReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			defer ro.Close()
			calls := 0
			objectsBySemanticIDBeforeQuery = func() { calls++ }
			t.Cleanup(func() { objectsBySemanticIDBeforeQuery = nil })
			_, err = ro.ObjectsBySemanticID(context.Background(), "x", "", 10)
			var unavailable *LookupIndexUnavailableError
			if !errors.As(err, &unavailable) || calls != 0 {
				t.Fatalf("read-only lookup=%v, hook calls=%d; want typed refusal, zero calls", err, calls)
			}
		})
	}
}
