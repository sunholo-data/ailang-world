package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestMarketingInventoryCompareBranches pins each CompareInventory refusal on its own. The
// host/daemon end-to-end test changes head and entries together, so it cannot tell which
// branch refused (controller, iteration 310).
func TestMarketingInventoryCompareBranches(t *testing.T) {
	base := func() Inventory {
		return Inventory{
			Schema:  "world/demo/capture-inventory/v1",
			Head:    json.RawMessage(`"sha256:head"`),
			Entries: []inventoryEntry{{"e0", "t0"}, {"e1", "t1"}},
			Objects: []inventoryObject{{"t0", "s0"}, {"t1", "s1"}},
		}
	}
	cases := []struct {
		name, code, detail string
		edit               func(before, after *Inventory)
	}{
		{"equal", "", "", func(_, _ *Inventory) {}},
		{"head", "inventory-mutated", "head differs", func(_, a *Inventory) { a.Head = json.RawMessage(`"sha256:other"`) }},
		{"entry", "inventory-mutated", "entry 1 differs", func(_, a *Inventory) { a.Entries[1].EntryHash = "e9" }},
		{"entry-ref", "inventory-mutated", "entry 0 differs", func(_, a *Inventory) { a.Entries[0].TransitionRef = "t9" }},
		{"count", "inventory-mutated", "entry count 2 != 3", func(_, a *Inventory) { a.Entries = append(a.Entries, inventoryEntry{"e2", "t1"}) }},
		{"object-hash", "inventory-mutated", "object set differs", func(_, a *Inventory) { a.Objects[1].SHA256 = "s9" }},
		{"object-ref", "inventory-mutated", "object set differs", func(_, a *Inventory) { a.Objects[0].Ref = "t9" }},
		{"object-missing", "inventory-mutated", "object set differs", func(_, a *Inventory) { a.Objects = a.Objects[:1] }},
		{"object-extra", "inventory-mutated", "object set differs", func(_, a *Inventory) { a.Objects = append(a.Objects, inventoryObject{"t2", "s2"}) }},
		{"empty-head-before", "inventory-empty", "head or entries empty", func(b, _ *Inventory) { b.Head = nil }},
		{"null-head-after", "inventory-empty", "head or entries empty", func(_, a *Inventory) { a.Head = json.RawMessage(`null`) }},
		{"blank-head-after", "inventory-empty", "head or entries empty", func(_, a *Inventory) { a.Head = json.RawMessage(`""`) }},
		{"no-entries-before", "inventory-empty", "head or entries empty", func(b, _ *Inventory) { b.Entries = nil }},
		{"no-entries-after", "inventory-empty", "head or entries empty", func(_, a *Inventory) { a.Entries = nil }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before, after := base(), base()
			tc.edit(&before, &after)
			err := CompareInventory(before, after)
			if tc.code == "" {
				if err != nil {
					t.Fatalf("identical inventories must compare equal: %v", err)
				}
				return
			}
			var r *Refusal
			if !errors.As(err, &r) || r.Code != tc.code || r.Detail != tc.detail {
				t.Fatalf("%s must refuse %s (%s): got %v", tc.name, tc.code, tc.detail, err)
			}
		})
	}
}

// TestMarketingInventoryReadPaging serves a log whose page cap (2) is below the requested
// limit, so ReadInventory must keep paging; it also pins the object digests, the text head
// and the non-200 refusal (controller, iteration 310).
func TestMarketingInventoryReadPaging(t *testing.T) {
	refs := []string{"t1", "t0", "t1", "t0", "t1"}
	bodies := map[string]string{"t0": `{"hash":"t0"}`, "t1": `{"hash":"t1"}`}
	serve := func(missing string) *httptest.Server {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /v1/head", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "sha256:head") })
		mux.HandleFunc("GET /v1/log", func(w http.ResponseWriter, r *http.Request) {
			from, _ := strconv.Atoi(r.URL.Query().Get("from"))
			if r.URL.Query().Get("limit") != "500" {
				http.Error(w, "unexpected limit", http.StatusBadRequest)
				return
			}
			items := []inventoryEntry{}
			for i := from; i < len(refs) && i < from+2; i++ {
				items = append(items, inventoryEntry{"e" + strconv.Itoa(i), refs[i]})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
		})
		mux.HandleFunc("GET /v1/objects/{ref}", func(w http.ResponseWriter, r *http.Request) {
			if r.PathValue("ref") == missing {
				http.Error(w, "object not found", http.StatusNotFound)
				return
			}
			fmt.Fprint(w, bodies[r.PathValue("ref")])
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)
		return srv
	}
	t.Run("pages", func(t *testing.T) {
		inv, err := ReadInventory(serve("").URL)
		if err != nil {
			t.Fatal(err)
		}
		if string(inv.Head) != `"sha256:head"` {
			t.Fatalf("text head must be kept as a JSON string: %s", inv.Head)
		}
		if len(inv.Entries) != len(refs) {
			t.Fatalf("every page must be read: got %d entries, want %d", len(inv.Entries), len(refs))
		}
		for i, e := range inv.Entries {
			if e.EntryHash != "e"+strconv.Itoa(i) || e.TransitionRef != refs[i] {
				t.Fatalf("entries must keep log order: %d = %+v", i, e)
			}
		}
		var want []inventoryObject
		for _, ref := range []string{"t0", "t1"} {
			sum := sha256.Sum256([]byte(bodies[ref]))
			want = append(want, inventoryObject{ref, hex.EncodeToString(sum[:])})
		}
		if fmt.Sprint(inv.Objects) != fmt.Sprint(want) {
			t.Fatalf("objects must be the distinct refs, sorted, with body digests: got %v want %v", inv.Objects, want)
		}
	})
	t.Run("object-404", func(t *testing.T) {
		if _, err := ReadInventory(serve("t0").URL); err == nil || !strings.Contains(err.Error(), "404") {
			t.Fatalf("a non-200 object read must fail the inventory: %v", err)
		}
	})
}
