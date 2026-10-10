package coordinator

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Baseline is the byte-verbatim manifest at reviewed parent 5ca111547fdd80480a8d42578fb59942ee8a0b63.
const descriptorBaselineSHA = "240840ca9a0f8f4c60ac159dd78fce26a6c681200e8a10c64fe0e08c0a228868"

var descriptorFloorIDs = []string{"ailang-check", "ailang-cli", "ailang-edit", "ailang-read", "ailang-run", "ailang-write", "builtins-search", "examples-search"}

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

func descriptorCensus(t *testing.T, entries []map[string]any) {
	t.Helper()
	ids, floor := []string{}, []string{}
	for _, e := range entries {
		id, ok := e["id"].(string)
		if !ok {
			t.Fatal("manifest census: id must be a string")
		}
		ids = append(ids, id)
		if id != "workspace-exec" {
			floor = append(floor, id)
		}
	}
	sort.Strings(ids)
	sort.Strings(floor)
	want := append(append([]string{}, descriptorFloorIDs...), "workspace-exec")
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
	descriptorCensus(t, base)
	descriptorCensus(t, candidate)
	if size, prose := descriptorBytes(t, base); size != 13947 || prose != 7430 {
		t.Fatalf("parent descriptor bytes = %d/%d, want 13947/7430", size, prose)
	}
	if size, prose := descriptorBytes(t, candidate); size > 10500 || prose > 3800 {
		t.Errorf("candidate byte budget: %d/%d exceeds 10500/3800", size, prose)
	} else {
		t.Logf("canonical eight-tool bytes=%d; description UTF-8 bytes=%d", size, prose)
	}
	byID := map[string]map[string]any{}
	for _, e := range base {
		byID[e["id"].(string)] = e
	}
	for _, e := range candidate {
		id := e["id"].(string)
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
	hashes := map[string]string{
		"packages/se-tools/se_tools/read.ail":            "076bd38e834b90d34bcacc3f64f5a863e2dbec04e0aa20d58bb630e874474672",
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
