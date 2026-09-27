package store

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func semanticObject(t *testing.T, s *Store, id, label string) Object {
	t.Helper()
	payload := []byte("sid-payload-" + label)
	o := Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("sid-iface")),
		SemanticID: id, Provenance: "sid-test-" + label, Payload: payload,
	}
	if err := s.PutObject(o); err != nil {
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
