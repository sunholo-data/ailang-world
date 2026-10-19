package coordinator

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"
)

// Baseline is the byte-verbatim manifest at reviewed parent 5ca111547fdd80480a8d42578fb59942ee8a0b63.
const descriptorBaselineSHA = "240840ca9a0f8f4c60ac159dd78fce26a6c681200e8a10c64fe0e08c0a228868"

var descriptorFloorIDs = []string{"ailang-check", "ailang-cli", "ailang-edit", "ailang-read", "ailang-run", "ailang-write", "builtins-search", "examples-search"}

// descriptorAddedIDs are the tools added after the baseline: row 191a's
// ailang-delete. Each is checked field by field in TestSeToolsDescriptorCompaction.
var descriptorAddedIDs = []string{"ailang-delete"}

// applyPostBaselineRows patches the immutable baseline with the reviewed
// non-description changes made after it, each written out here independently
// of the manifest, so any other drift still reds the equivalence check.
func applyPostBaselineRows(t *testing.T, base []map[string]any) {
	t.Helper()
	patch := func(e map[string]any, schema, prop, raw string) {
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			t.Fatal(err)
		}
		e[schema].(map[string]any)["properties"].(map[string]any)[prop] = v
	}
	for _, e := range base {
		switch e["id"] {
		case "ailang-read": // row 190: ranged reads
			patch(e, "inputSchema", "offset", `{"type":"integer","minimum":1,"description":"First line, 1-based"}`)
			patch(e, "inputSchema", "limit", `{"type":"integer","minimum":1,"description":"Line count, at most 2000"}`)
			for _, p := range []string{"total_lines", "offset", "limit"} {
				patch(e, "outputSchema", p, `{"type":"integer"}`)
			}
		}
	}
}

// checkAddedDescriptor is row 191a's descriptor: ailang-delete takes exactly
// the baseline read's {path} schema, answers the baseline write's {ok,
// refused, world}, and declares only its own Workspace.Delete triple.
func checkAddedDescriptor(t *testing.T, e map[string]any, byID map[string]map[string]any) {
	t.Helper()
	want := map[string]any{
		"id": "ailang-delete", "title": "Delete (sandboxed)", "transitionFnFile": "packages/se-tools/se_tools/delete.ail",
		"inputSchema": byID["ailang-read"]["inputSchema"], "outputSchema": byID["ailang-write"]["outputSchema"],
		"access":          map[string]any{"effect": "Workspace.Delete", "scope": "worktree", "cost": float64(0)},
		"declaredEffects": []any{map[string]any{"effect": "Workspace.Delete", "scope": "worktree", "cost": float64(1)}},
	}
	if desc, ok := e["description"].(string); !ok || !strings.Contains(desc, "Workspace.Delete") {
		t.Errorf("ailang-delete description %q does not name its effect", e["description"])
	}
	got := map[string]any{}
	for k, v := range e {
		if k != "description" {
			got[k] = v
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ailang-delete descriptor %v, want %v", got, want)
	}
}

func descriptorEntries(t *testing.T, path string) []map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var entries []map[string]any
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatal(err)
	}
	return entries
}

// descriptorCensus holds a manifest to the eight floor tools, workspace-exec
// and exactly the added tools given.
func descriptorCensus(t *testing.T, entries []map[string]any, added []string) {
	t.Helper()
	ids, floor := []string{}, []string{}
	extra := map[string]bool{"workspace-exec": true}
	for _, id := range added {
		extra[id] = true
	}
	for _, e := range entries {
		id, ok := e["id"].(string)
		if !ok {
			t.Fatal("manifest census: id must be a string")
		}
		ids = append(ids, id)
		if !extra[id] {
			floor = append(floor, id)
		}
	}
	sort.Strings(ids)
	sort.Strings(floor)
	want := append(append(append([]string{}, descriptorFloorIDs...), "workspace-exec"), added...)
	sort.Strings(want)
	if !reflect.DeepEqual(ids, want) || !reflect.DeepEqual(floor, descriptorFloorIDs) {
		t.Fatalf("manifest/floor census: got %v / %v, want %v / %v", ids, floor, want, descriptorFloorIDs)
	}
}

func descriptorBytes(t *testing.T, entries []map[string]any) (int, int) {
	t.Helper()
	tools, prose := []map[string]any{}, 0
	for _, e := range entries {
		if e["id"] == "workspace-exec" {
			continue
		}
		desc, ok := e["description"].(string)
		if !ok || strings.TrimSpace(desc) == "" {
			t.Fatal("empty tool description")
		}
		prose += len([]byte(desc))
		tools = append(tools, map[string]any{"name": e["id"], "description": desc, "inputSchema": e["inputSchema"], "outputSchema": e["outputSchema"]})
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(tools); err != nil {
		t.Fatal(err)
	}
	return len(bytes.TrimSuffix(b.Bytes(), []byte("\n"))), prose
}

func TestSeToolsDescriptorCompaction(t *testing.T) {
	raw, err := os.ReadFile("testdata/se-tools-descriptor-baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != descriptorBaselineSHA {
		t.Fatal("immutable baseline SHA mismatch")
	}
	base := descriptorEntries(t, "testdata/se-tools-descriptor-baseline.json")
	candidate := descriptorEntries(t, filepath.Join(seToolsRepoRoot(t), seToolsManifest))
	descriptorCensus(t, base, nil)
	descriptorCensus(t, candidate, descriptorAddedIDs)
	if size, prose := descriptorBytes(t, base); size != 13947 || prose != 7430 {
		t.Fatalf("parent descriptor bytes = %d/%d, want 13947/7430", size, prose)
	}
	// The compaction budget binds the eight floor tools, as it did at the
	// baseline; a tool added since is measured, not budgeted against them.
	// Row 190 moved the size bound 10500 -> 10600 for ailang-read's range
	// schema (offset/limit in, total_lines/offset/limit out), which takes the
	// eight tools past 10500: new arguments need schema bytes. The prose bound is unchanged, so the
	// range is described by the schema, not the description.
	floor := []map[string]any{}
	for _, e := range candidate {
		if !slices.Contains(descriptorAddedIDs, e["id"].(string)) {
			floor = append(floor, e)
		}
	}
	if size, prose := descriptorBytes(t, floor); size > 10600 || prose > 3800 {
		t.Errorf("candidate byte budget: %d/%d exceeds 10600/3800", size, prose)
	} else {
		all, allProse := descriptorBytes(t, candidate)
		t.Logf("canonical eight-tool bytes=%d; description UTF-8 bytes=%d (with added tools %d/%d)", size, prose, all, allProse)
	}
	pristine := map[string]map[string]any{}
	for _, e := range descriptorEntries(t, "testdata/se-tools-descriptor-baseline.json") {
		pristine[e["id"].(string)] = e
	}
	applyPostBaselineRows(t, base)
	byID := map[string]map[string]any{}
	for _, e := range base {
		byID[e["id"].(string)] = e
	}
	for _, e := range candidate {
		id := e["id"].(string)
		if slices.Contains(descriptorAddedIDs, id) {
			checkAddedDescriptor(t, e, pristine)
			continue
		}
		old := byID[id]
		if id != "workspace-exec" {
			if e["description"] == old["description"] {
				t.Errorf("description unchanged: %s", id)
			}
			delete(e, "description")
			delete(old, "description")
		}
		if !reflect.DeepEqual(e, old) {
			t.Errorf("complete non-description equivalence: %s", id)
		}
	}
}

func TestSeToolsDescriptorSourceEquality(t *testing.T) {
	// Raw byte hashes independently banked from the same reviewed parent; no git or external binary required in CI.
	// Moved since, each by its own row: read.ail by row 190 (ranged reads);
	// delete.ail is row 191a's new module.
	hashes := map[string]string{
		"packages/se-tools/se_tools/read.ail":            "00213aecddeb49c62d4aec429682cf114b58d53aadd057229ef1662e751b75cf",
		"packages/se-tools/se_tools/delete.ail":          "084d90f35b3184af04d77893309d8453b2360dc1eb070947cb5b4d1fda47c834",
		"packages/se-tools/se_tools/write.ail":           "2d90757c61274f6f172551f179b254178f3ccfcedd21ee7955ab2e47d364deda",
		"packages/se-tools/se_tools/edit.ail":            "2d12856ab07f987c1d26c82b9320503e830817f27b625e7a80320b22ae32f286",
		"packages/se-tools/se_tools/check.ail":           "440b6491de07656df8161bb69b09af82e74306fd52b2152744921720c2cbd03b",
		"packages/se-tools/se_tools/run.ail":             "add740e9a157821fbef3de116ad60c2a33fd7fd9ca1d4f103c6efc8106317286",
		"packages/se-tools/se_tools/builtins_search.ail": "7ed1f8436865c194b2ce30611c0de92a73b0d86bf4776b2a28bf1c8d192e988b",
		"packages/se-tools/se_tools/examples_search.ail": "bd6d41cfae833e636cd079352ad4dc63fde37800dd8bdda4da39abbbb2e0e6c7",
		"packages/se-tools/se_tools/cli.ail":             "ce221e5ebf0c12929a3a36535ac1297ca695ffb0531d79e256c255697c962a70",
		"packages/se-tools/se_tools/exec.ail":            "51d93ce680eea3031eaf786530b3c2a4d61b3b547440a4e542786695b05c9a05",
	}
	for path, want := range hashes {
		raw, err := os.ReadFile(filepath.Join(seToolsRepoRoot(t), filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(raw)); got != want {
			t.Errorf("parent transition source-byte equivalence: %s got %s want %s", path, got, want)
		}
	}
	t.Logf("compared %d transition sources to parent raw-byte SHA-256", len(hashes))
}
