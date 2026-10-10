package daemon

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// Exercise publication, session admission and the real HTTP projection, not a reconstructed payload.
func TestSeToolsDescriptorLiveExactEight(t *testing.T) {
	r := newSeRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite,
		broker.EffectAilangCheck, broker.EffectAilangRun, broker.EffectAilangRunEnv, broker.EffectAilangRunNet,
		broker.EffectAilangDiscover, broker.EffectAilangCLI)
	code, raw := r.post("/mcp/", token, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	payload := strings.TrimSpace(strings.TrimPrefix(string(raw), "event: message\ndata: "))
	var wire struct {
		Result struct {
			Tools []map[string]any `json:"tools"`
		} `json:"result"`
		Error any `json:"error"`
	}
	if code != http.StatusOK || json.Unmarshal([]byte(payload), &wire) != nil || wire.Error != nil {
		t.Fatalf("live tools/list: HTTP %d %s", code, raw)
	}
	manifest, err := os.ReadFile(seToolsManifest)
	if err != nil {
		t.Fatal(err)
	}
	var entries []map[string]any
	if err := json.Unmarshal(manifest, &entries); err != nil {
		t.Fatal(err)
	}
	expected := map[string]map[string]any{}
	for _, e := range entries {
		if e["id"] == "workspace-exec" {
			continue
		}
		expected[e["id"].(string)] = map[string]any{"name": e["id"], "description": e["description"], "inputSchema": e["inputSchema"], "outputSchema": e["outputSchema"]}
	}
	names := []string{"ailang-check", "ailang-cli", "ailang-edit", "ailang-read", "ailang-run", "ailang-write", "builtins-search", "examples-search"}
	if len(wire.Result.Tools) != len(names) {
		t.Fatalf("live exact-eight census: %d tools", len(wire.Result.Tools))
	}
	seen := map[string]bool{}
	for _, tool := range wire.Result.Tools {
		name, _ := tool["name"].(string)
		if seen[name] || !reflect.DeepEqual(tool, expected[name]) {
			t.Fatalf("live descriptor equality (description/name/input/output): %s", name)
		}
		seen[name] = true
	}
	for _, name := range names {
		if !seen[name] {
			t.Fatalf("live exact-eight census: missing %s", name)
		}
	}
	// Opt-in bank retains the exact SSE wire, separately from canonical descriptor size.
	if path := os.Getenv("WORLD_DESCRIPTOR_WIRE_CAPTURE"); path != "" {
		if err := os.WriteFile(path, raw, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("live tools/list: exactly eight descriptors; wire bytes=%d", len(raw))
	t.Log("each description and complete input/output schema equals published manifest")
}
