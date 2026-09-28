package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestObjectReferencesValidation(t *testing.T) {
	ref := hashref.SumSHA256([]byte("target"))
	world := hashref.SumSHA256([]byte("world"))
	t.Run("limits", func(t *testing.T) {
		for _, limit := range []int{-1, 0, 501, 1000} {
			err := validateObjectReferences(ref, nil, limit)
			var bad *InvalidLimitError
			if !errors.As(err, &bad) {
				t.Fatalf("limit %d: %v", limit, err)
			}
		}
		for _, limit := range []int{1, 500} {
			if err := validateObjectReferences(ref, nil, limit); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("ref", func(t *testing.T) {
		err := validateObjectReferences(hashref.HashRef{}, nil, 1)
		var bad *InvalidRefError
		if !errors.As(err, &bad) {
			t.Fatalf("ref error: %v", err)
		}
	})
	t.Run("cursor", func(t *testing.T) {
		good := []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 0}, {Kind: ReferenceTransitionFn, EntryIndex: 2}, {Kind: ReferenceInterpreter, EntryIndex: 10}, {Kind: ReferenceStateRoot, WorldRef: world}}
		for _, c := range good {
			if err := validateObjectReferences(ref, &c, 1); err != nil {
				t.Fatalf("good %+v: %v", c, err)
			}
		}
		bad := []ObjectReferenceCursor{{Kind: -1}, {Kind: 4}, {Kind: ReferenceTransitionRef, EntryIndex: -1}, {Kind: ReferenceTransitionFn, WorldRef: world}, {Kind: ReferenceStateRoot}, {Kind: ReferenceStateRoot, EntryIndex: 1, WorldRef: world}}
		for _, c := range bad {
			err := validateObjectReferences(ref, &c, 1)
			var invalid *InvalidObjectReferenceCursorError
			if !errors.As(err, &invalid) {
				t.Fatalf("bad %+v: %v", c, err)
			}
		}
	})
	s := openMem(t)
	got, err := s.ObjectReferences(context.Background(), ref, nil, 1)
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("empty page=%v %v", got, err)
	}
}

func refTestHash(label string) hashref.HashRef { return hashref.SumSHA256([]byte(label)) }

func insertReferenceEntry(t *testing.T, s *Store, index int64, transition, fn, interp hashref.HashRef) {
	t.Helper()
	_, err := s.db.Exec(`INSERT INTO log_entries(entry_index,entry_hash_ref,semantics_epoch,transition_fn_ref,interpreter_ref,prev_entry_hash_ref,written_by,transition_ref) VALUES(?,?,?,?,?,?,?,?)`, index, refTestHash(fmt.Sprintf("entry-%d", index)).String(), 1, fn.String(), interp.String(), refTestHash("prev").String(), "test", transition.String())
	if err != nil {
		t.Fatal(err)
	}
}
func insertReferenceWorld(t *testing.T, s *Store, label string, state hashref.HashRef) hashref.HashRef {
	t.Helper()
	ref := refTestHash(label)
	_, err := s.db.Exec(`INSERT INTO worlds(world_ref,revision,state_root,log_head) VALUES(?,?,?,?)`, ref.String(), 1, state.String(), refTestHash("log-head").String())
	if err != nil {
		t.Fatal(err)
	}
	return ref
}
func cursors(items []ObjectReference) []ObjectReferenceCursor {
	out := make([]ObjectReferenceCursor, len(items))
	for i, x := range items {
		out[i] = x.Cursor
	}
	return out
}
func wantPage(t *testing.T, s *Store, ref hashref.HashRef, after *ObjectReferenceCursor, limit int, want []ObjectReferenceCursor) {
	t.Helper()
	got, err := s.ObjectReferences(context.Background(), ref, after, limit)
	if err != nil || !reflect.DeepEqual(cursors(got), want) {
		t.Fatalf("page after=%+v limit=%d: got=%+v err=%v want=%+v", after, limit, cursors(got), err, want)
	}
}

func TestObjectReferencesPage(t *testing.T) {
	ref, other := refTestHash("target"), refTestHash("other")
	t.Run("transition-cursor", func(t *testing.T) {
		s := openMem(t)
		for _, i := range []int64{0, 2, 10} {
			insertReferenceEntry(t, s, i, ref, other, other)
		}
		after := ObjectReferenceCursor{Kind: ReferenceTransitionRef, EntryIndex: 2}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 10}})
	})
	t.Run("function-cursor", func(t *testing.T) {
		s := openMem(t)
		for _, i := range []int64{0, 2, 10} {
			insertReferenceEntry(t, s, i, other, ref, other)
		}
		after := ObjectReferenceCursor{Kind: ReferenceTransitionFn, EntryIndex: 2}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceTransitionFn, EntryIndex: 10}})
	})
	t.Run("interpreter-cursor", func(t *testing.T) {
		s := openMem(t)
		for _, i := range []int64{0, 2, 10} {
			insertReferenceEntry(t, s, i, other, other, ref)
		}
		after := ObjectReferenceCursor{Kind: ReferenceInterpreter, EntryIndex: 2}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceInterpreter, EntryIndex: 10}})
	})
	t.Run("world-cursor", func(t *testing.T) {
		s := openMem(t)
		a := insertReferenceWorld(t, s, "a", ref)
		b := insertReferenceWorld(t, s, "b", ref)
		keys := []hashref.HashRef{a, b}
		sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
		after := ObjectReferenceCursor{Kind: ReferenceStateRoot, WorldRef: keys[0]}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceStateRoot, WorldRef: keys[1]}})
	})
	t.Run("mixed-kind-bound", func(t *testing.T) {
		s := openMem(t)
		insertReferenceEntry(t, s, 1, ref, ref, other)
		insertReferenceEntry(t, s, 2, other, ref, other)
		wantPage(t, s, ref, nil, 2, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 1}, {Kind: ReferenceTransitionFn, EntryIndex: 1}})
	})
	t.Run("within-kind", func(t *testing.T) {
		s := openMem(t)
		insertReferenceEntry(t, s, 1, ref, other, other)
		insertReferenceEntry(t, s, 2, ref, other, other)
		after := ObjectReferenceCursor{Kind: ReferenceTransitionRef, EntryIndex: 1}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 2}})
	})
	t.Run("numeric-order", func(t *testing.T) {
		s := openMem(t)
		insertReferenceEntry(t, s, 10, ref, other, other)
		insertReferenceEntry(t, s, 2, ref, other, other)
		wantPage(t, s, ref, nil, 2, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 2}, {Kind: ReferenceTransitionRef, EntryIndex: 10}})
	})
	t.Run("index-zero", func(t *testing.T) {
		s := openMem(t)
		insertReferenceEntry(t, s, 0, ref, other, other)
		wantPage(t, s, ref, nil, 1, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 0}})
	})
	for _, n := range []int{0, 1, 100, 101, 501} {
		t.Run(fmt.Sprintf("count-%d", n), func(t *testing.T) {
			s := openMem(t)
			for i := 0; i < n; i++ {
				insertReferenceEntry(t, s, int64(i), ref, other, other)
			}
			var all []ObjectReferenceCursor
			var after *ObjectReferenceCursor
			for {
				page, err := s.ObjectReferences(context.Background(), ref, after, 100)
				if err != nil {
					t.Fatal(err)
				}
				all = append(all, cursors(page)...)
				if len(page) < 100 {
					break
				}
				last := page[len(page)-1].Cursor
				after = &last
			}
			if len(all) != n {
				t.Fatalf("walk got %d want %d", len(all), n)
			}
			for i, c := range all {
				if c.Kind != ReferenceTransitionRef || c.EntryIndex != int64(i) {
					t.Fatalf("edge %d: %+v", i, c)
				}
			}
		})
	}
	t.Run("absent-cursor", func(t *testing.T) {
		s := openMem(t)
		insertReferenceEntry(t, s, 2, ref, other, other)
		after := ObjectReferenceCursor{Kind: ReferenceTransitionRef, EntryIndex: 1}
		wantPage(t, s, ref, &after, 1, []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 2}})
	})
}

func TestObjectReferencesScope(t *testing.T) {
	ref, other := refTestHash("target"), refTestHash("other")
	tests := []struct {
		name string
		kind ReferenceKind
	}{{"transition", ReferenceTransitionRef}, {"function", ReferenceTransitionFn}, {"interpreter", ReferenceInterpreter}, {"world", ReferenceStateRoot}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			if tc.kind == ReferenceStateRoot {
				good := insertReferenceWorld(t, s, "good", ref)
				insertReferenceWorld(t, s, "decoy", other)
				wantPage(t, s, ref, nil, 10, []ObjectReferenceCursor{{Kind: ReferenceStateRoot, WorldRef: good}})
				return
			}
			parts := [3]hashref.HashRef{other, other, other}
			parts[tc.kind] = ref
			insertReferenceEntry(t, s, 1, parts[0], parts[1], parts[2])
			insertReferenceEntry(t, s, 2, other, other, other)
			wantPage(t, s, ref, nil, 10, []ObjectReferenceCursor{{Kind: tc.kind, EntryIndex: 1}})
		})
	}
}

func TestObjectReferencesOutOfScope(t *testing.T) {
	s := openMem(t)
	ref, other := refTestHash("target"), refTestHash("other")
	o := obj("interface-source", "interface")
	o.InterfaceHash = ref
	if err := s.PutObject(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	if err := s.SetRegistryHead(context.Background(), "test", ref); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec("INSERT INTO journal(kind,invocation_id,object_ref) VALUES('intent','scope',?)", ref.String()); err != nil {
		t.Fatal(err)
	}
	insertReferenceEntry(t, s, 1, other, other, other)
	if _, err := s.db.Exec("UPDATE log_entries SET prev_entry_hash_ref=? WHERE entry_index=1", ref.String()); err != nil {
		t.Fatal(err)
	}
	insertReferenceWorld(t, s, "out-of-scope-world", other)
	if _, err := s.db.Exec("UPDATE worlds SET log_head=?", ref.String()); err != nil {
		t.Fatal(err)
	}
	wantPage(t, s, ref, nil, 10, []ObjectReferenceCursor{})
}

func TestObjectReferencesGuard(t *testing.T) {
	t.Run("no-sql", func(t *testing.T) {
		s := openMem(t)
		s.referenceIndexesAvailable = false
		calls := 0
		objectReferencesBeforeQuery = func() { calls++ }
		defer func() { objectReferencesBeforeQuery = nil }()
		_, err := s.ObjectReferences(context.Background(), refTestHash("target"), nil, 10)
		var unavailable *ReferenceIndexUnavailableError
		if !errors.As(err, &unavailable) || calls != 0 {
			t.Fatalf("guard err=%v calls=%d", err, calls)
		}
	})
}

func TestObjectReferencesCancellation(t *testing.T) {
	t.Run("blocked-connection", func(t *testing.T) {
		s := openMem(t)
		ref := refTestHash("target")
		other := refTestHash("other")
		insertReferenceEntry(t, s, 0, ref, other, other)
		conn, err := s.db.Conn(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		timer := time.AfterFunc(100*time.Millisecond, func() { conn.Close() })
		defer timer.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		_, err = s.ObjectReferences(ctx, ref, nil, 1)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("blocked query err=%v", err)
		}
		conn.Close()
	})
}

func TestObjectReferencesPlans(t *testing.T) {
	s := openMem(t)
	for i, q := range []string{transitionReferencesSQL, functionReferencesSQL, interpreterReferencesSQL, worldReferencesSQL} {
		var after any = int64(-1)
		if i == 3 {
			after = ""
		}
		rows, err := s.db.Query("EXPLAIN QUERY PLAN "+q, refTestHash("target").String(), after, 10)
		if err != nil {
			t.Fatal(err)
		}
		details := ""
		for rows.Next() {
			var a, b, c int
			var detail string
			if err := rows.Scan(&a, &b, &c, &detail); err != nil {
				t.Fatal(err)
			}
			details += detail
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		rows.Close()
		if !strings.Contains(details, "USING COVERING INDEX "+referenceIndexSpecs[i].name) || strings.Contains(details, "TEMP B-TREE") {
			t.Fatalf("plan %d: %s", i, details)
		}
	}
}

// TestObjectReferencesMeasuredAtScale measures the final file-store access path.
func TestObjectReferencesMeasuredAtScale(t *testing.T) {
	rawN := os.Getenv("WORLD103_N")
	if rawN == "" {
		t.Skip("set WORLD103_N for file-store measurement")
	}
	n, err := strconv.Atoi(rawN)
	if err != nil || n < 1 {
		t.Fatalf("WORLD103_N=%q", rawN)
	}
	path := filepath.Join(t.TempDir(), "scale.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	fn, interp, state := obj("scale-function", "fn"), obj("scale-interpreter", "interpreter"), obj("scale-state", "state")
	for _, o := range []Object{fn, interp, state} {
		if err := s.PutObject(context.Background(), o); err != nil {
			t.Fatal(err)
		}
	}
	current := seedGenesis(t, s)
	current.StateRoot = state.Hash
	current.Ref = refTestHash("scale-genesis")
	if err := s.PutWorld(current); err != nil {
		t.Fatal(err)
	}
	if err := s.SelectHead(current.Ref); err != nil {
		t.Fatal(err)
	}
	var sparse hashref.HashRef
	startWrite := time.Now()
	for i := 0; i < n; i++ {
		body := obj(fmt.Sprintf("scale-body-%d", i), "body")
		entryHash := refTestHash(fmt.Sprintf("scale-entry-%d", i))
		next := World{Ref: refTestHash(fmt.Sprintf("scale-world-%d", i)), Revision: int64(i + 1), StateRoot: state.Hash, LogHead: entryHash}
		commit := Commit{ObservedHead: current.Ref, Objects: []Object{body}, NextWorld: next, Entry: LogEntry{Header: LogHeader{EntryIndex: int64(i), SemanticsEpoch: 1, TransitionFn: fn.Hash, Interpreter: interp.Hash, PrevEntryHash: current.LogHead, WrittenBy: "scale"}, EntryHash: entryHash, TransitionRef: body.Hash}}
		if err := s.Commit(commit); err != nil {
			t.Fatal(err)
		}
		current = next
		sparse = body.Hash
	}
	t.Logf("production_writes_n=%d elapsed_ms=%.3f", n, float64(time.Since(startWrite).Nanoseconds())/1e6)
	oracle := func(ref hashref.HashRef) []ObjectReferenceCursor {
		var all []ObjectReferenceCursor
		for kind, q := range []string{transitionReferencesSQL, functionReferencesSQL, interpreterReferencesSQL, worldReferencesSQL} {
			var after any = int64(-1)
			if kind == 3 {
				after = ""
			}
			rows, err := s.db.Query(q, ref.String(), after, n+2)
			if err != nil {
				t.Fatal(err)
			}
			for rows.Next() {
				c := ObjectReferenceCursor{Kind: ReferenceKind(kind)}
				if kind == 3 {
					var text string
					if err := rows.Scan(&text); err != nil {
						t.Fatal(err)
					}
					c.WorldRef, err = hashref.Parse(text)
					if err != nil {
						t.Fatal(err)
					}
				} else {
					if err := rows.Scan(&c.EntryIndex); err != nil {
						t.Fatal(err)
					}
				}
				all = append(all, c)
			}
			if err := rows.Err(); err != nil {
				t.Fatal(err)
			}
			rows.Close()
		}
		return all
	}
	for _, tc := range []struct {
		name string
		ref  hashref.HashRef
	}{{"sparse", sparse}, {"hot", fn.Hash}, {"world", state.Hash}, {"absent", refTestHash("scale-absent")}} {
		want := oracle(tc.ref)
		var got []ObjectReferenceCursor
		var after *ObjectReferenceCursor
		for len(got) < len(want)+1 {
			page, err := s.ObjectReferences(context.Background(), tc.ref, after, 101)
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, cursors(page)...)
			if len(page) < 101 {
				break
			}
			last := page[len(page)-1].Cursor
			after = &last
		}
		if !reflect.DeepEqual(got, want) && !(len(got) == 0 && len(want) == 0) {
			t.Fatalf("%s ordered pages: got=%d want=%d", tc.name, len(got), len(want))
		}
		samples := make([]int64, 101)
		for i := range samples {
			start := time.Now()
			_, err := s.ObjectReferences(context.Background(), tc.ref, nil, 101)
			if err != nil {
				t.Fatal(err)
			}
			samples[i] = time.Since(start).Nanoseconds()
		}
		sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
		t.Logf("%s rows=%d p50_us=%.3f pmax_us=%.3f", tc.name, len(want), float64(samples[50])/1000, float64(samples[100])/1000)
	}
	for i, q := range []string{transitionReferencesSQL, functionReferencesSQL, interpreterReferencesSQL, worldReferencesSQL} {
		var after any = int64(-1)
		if i == 3 {
			after = ""
		}
		rows, err := s.db.Query("EXPLAIN QUERY PLAN "+q, sparse.String(), after, 101)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var a, b, c int
			var detail string
			if err := rows.Scan(&a, &b, &c, &detail); err != nil {
				t.Fatal(err)
			}
			t.Logf("plan_%d=%s", i, detail)
		}
		rows.Close()
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("existing_index_Open_us=%.3f", float64(time.Since(start).Nanoseconds())/1000)
	for _, spec := range referenceIndexSpecs {
		if _, err := s.db.Exec("DROP INDEX " + spec.name); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	start = time.Now()
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("provisioning_Open_us=%.3f", float64(time.Since(start).Nanoseconds())/1000)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	afterStat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("disk_before=%d disk_after=%d delta=%d", before.Size(), afterStat.Size(), afterStat.Size()-before.Size())
}
