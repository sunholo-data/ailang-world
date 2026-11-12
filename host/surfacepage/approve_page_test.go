package surfacepage

// Row 201 A3 sprint 3 M6c (design §2.7, §5.2, §6 PG3, PG3b): the card's field table is the
// canonical request's, field for field.

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// pageField is one entry of approve.js's FIELDS.
type pageField struct {
	Path, Label, Group, Type, No, Yes string
}

// pageFields is approve.js's FIELDS literal, between its /* FIELDS */ markers, parsed as JSON.
func pageFields(t *testing.T) []pageField {
	t.Helper()
	_, rest, ok := bytes.Cut(ApproveScript, []byte("/* FIELDS */"))
	lit, _, ok2 := bytes.Cut(rest, []byte("/* /FIELDS */"))
	var out []pageField
	if !ok || !ok2 || json.Unmarshal(lit, &out) != nil {
		t.Fatalf("approve.js has no FIELDS JSON between its markers")
	}
	return out
}

// structFields is every JSON key path of t's json tags ([] for a slice's element) with its kind,
// and every tag with an option (",omitempty" or any other).
func structFields(typ reflect.Type, prefix string, paths map[string]string, options *[]string) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name, opts, _ := strings.Cut(f.Tag.Get("json"), ",")
		if opts != "" {
			*options = append(*options, typ.Name()+"."+f.Name+": "+opts)
		}
		p, ft := prefix+name, f.Type
		if ft.Kind() == reflect.Slice {
			p, ft = p+"[]", ft.Elem()
		}
		switch ft.Kind() {
		case reflect.Struct:
			structFields(ft, p+".", paths, options)
		case reflect.String:
			paths[p] = "string"
		case reflect.Bool:
			paths[p] = "bool"
		case reflect.Int, reflect.Int64:
			paths[p] = "int"
		default:
			paths[p] = "unsupported " + ft.Kind().String()
		}
	}
}

// PG3: FIELDS names exactly approval.MintRequest's JSON key paths, each with its type, a label and
// a group the page shows (a boolean with both of its words). Control: the struct has 21 top-level
// keys, the count the browser measured in the stored bytes (design V24).
func TestApprovePageRendersEveryCanonicalField(t *testing.T) {
	want := map[string]string{}
	var options []string
	structFields(reflect.TypeOf(approval.MintRequest{}), "", want, &options)
	if n := reflect.TypeOf(approval.MintRequest{}).NumField(); n != 21 {
		t.Fatalf("control: MintRequest has %d top-level keys, want 21 (design V24)", n)
	}
	got := map[string]string{}
	groups := map[string]bool{"what": true, "where": true, "bound": true, "stores": true, "daemon": true}
	for _, f := range pageFields(t) {
		if _, dup := got[f.Path]; dup {
			t.Errorf("FIELDS lists %s twice", f.Path)
		}
		got[f.Path] = f.Type
		if f.Label == "" || !groups[f.Group] || (f.Type == "bool") != (f.No != "" && f.Yes != "") {
			t.Errorf("FIELDS entry %+v: a label, a shown group, and both words for a boolean only", f)
		}
	}
	if !reflect.DeepEqual(got, want) {
		var missing, extra []string
		for p, k := range want {
			if got[p] != k {
				missing = append(missing, p+" "+k)
			}
		}
		for p, k := range got {
			if want[p] != k {
				extra = append(extra, p+" "+k)
			}
		}
		sort.Strings(missing)
		sort.Strings(extra)
		t.Fatalf("FIELDS does not match approval.MintRequest: no renderer for %q; not in the struct %q", missing, extra)
	}
}

// leafPaths is every leaf key path of a decoded JSON value ([] for an element).
func leafPaths(v any, prefix string, out map[string]bool) {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			leafPaths(e, p, out)
		}
	case []any:
		for _, e := range x {
			leafPaths(e, prefix+"[]", out)
		}
	default:
		out[prefix] = true
	}
}

// PG3b: no tag of the canonical request carries an option (so no omitempty: design V28), and a
// degraded, non-git request with reuse false and every string empty still serializes all 21 keys
// and every FIELDS path, so the page's missing-path refusal never fires on a canonical request;
// branch, which XA2 arm (c) deletes, is one of them.
func TestApprovePageFieldTagControls(t *testing.T) {
	paths := map[string]string{}
	var options []string
	structFields(reflect.TypeOf(approval.MintRequest{}), "", paths, &options)
	if len(options) != 0 {
		t.Fatalf("tags with options on the canonical request: %q", options)
	}
	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: "w201d", Degraded: true,
		Grants: []approval.Grant{{Effect: "Workspace.Read"}}})
	if err != nil {
		t.Fatal(err)
	}
	var top map[string]any
	if err := json.Unmarshal(b, &top); err != nil || len(top) != 21 {
		t.Fatalf("the degraded request has %d top-level keys (%v), want 21", len(top), err)
	}
	if top["reuse"] != false || top["degraded"] != true || top["branch"] != "" {
		t.Fatalf("the degraded request: %s", b)
	}
	have := map[string]bool{}
	leafPaths(top, "", have)
	for _, f := range pageFields(t) {
		if !have[f.Path] {
			t.Errorf("the degraded request lacks %s", f.Path)
		}
		delete(have, f.Path)
	}
	if len(have) != 0 || paths["branch"] != "string" {
		t.Errorf("leaf paths with no FIELDS entry %v; branch %q", have, paths["branch"])
	}
}
