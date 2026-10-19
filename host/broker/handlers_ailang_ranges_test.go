package broker

// Row 190 (w-ailang-read-ranges): ailang-read's offset/limit. The broker adds
// no code for them — every payload key but op and flags passes through to
// policy-tool verbatim — so these tests pin the pass-through on the fake tool
// and the upstream shapes (ailang#1872) on the real one.
//
// policyToolRequestKeys stays {op, flags}: the alias guard protects the keys
// the handler itself reads or polices, and it reads and polices neither
// offset nor limit, so an alias of them bypasses no broker check. v0.54.15's
// policy-tool also refuses any key that is not exactly one of its own
// (TestPolicyToolRefusesCaseVariantRangeKeysItself), so "OFFSET" never
// reaches a read as a range either.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAilangReadRangeKeysPassThrough: a ranged read's payload reaches
// policy-tool's stdin with its integer offset and limit unchanged.
func TestAilangReadRangeKeysPassThrough(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "ok")
	if _, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"op":"read","path":"docs/big.md","offset":201,"limit":100}`)); err != nil {
		t.Fatal(err)
	}
	var sent map[string]any
	if err := json.Unmarshal(mustRead(t, filepath.Join(tool.dir, "stdin")), &sent); err != nil {
		t.Fatal(err)
	}
	if sent["op"] != "read" || sent["path"] != "docs/big.md" || sent["offset"] != float64(201) || sent["limit"] != float64(100) || len(sent) != 4 {
		t.Fatalf("policy-tool got %v, want the ranged read verbatim", sent)
	}
}

// TestAilangReadRangesOnTheToolBinary is ailang#1872 as World serves it: a
// whole read carries no range keys; a ranged read returns its lines with
// total_lines, offset and limit; limit is clamped (to 2000, here to the
// file); an offset past the end is ok with no content; and the bounds the
// plan enforces are policy-tool's own too.
func TestAilangReadRangesOnTheToolBinary(t *testing.T) {
	f := newRealFixture(t, "worktree")
	mustWrite(t, filepath.Join(f.root, "a.txt"), "l1\nl2\nl3\nl4\n")
	read := func(extra map[string]any) map[string]any {
		t.Helper()
		payload := map[string]any{"op": "read", "path": "a.txt"}
		for k, v := range extra {
			payload[k] = v
		}
		return f.op(t, EffectWorkspaceRead, payload)
	}
	whole := read(nil)
	if whole["ok"] != true || whole["content"] != "l1\nl2\nl3\nl4\n" {
		t.Fatalf("whole read = %v", whole)
	}
	for _, k := range []string{"total_lines", "offset", "limit"} {
		if _, present := whole[k]; present {
			t.Fatalf("a whole read carries %q: %v", k, whole)
		}
	}
	for _, tc := range []struct {
		name                 string
		args                 map[string]any
		content              any // nil = no content key
		total, offset, limit float64
	}{
		{"offset 2 limit 2", map[string]any{"offset": 2, "limit": 2}, "l2\nl3\n", 4, 2, 2},
		{"offset only", map[string]any{"offset": 3}, "l3\nl4\n", 4, 3, 2},
		{"limit clamped", map[string]any{"limit": 5000}, "l1\nl2\nl3\nl4\n", 4, 1, 4},
		{"offset past the end", map[string]any{"offset": 9}, nil, 4, 9, 0},
	} {
		got := read(tc.args)
		if got["ok"] != true || got["content"] != tc.content || got["total_lines"] != tc.total ||
			got["offset"] != tc.offset || got["limit"] != tc.limit {
			t.Fatalf("%s: %v; want content %q total_lines %v offset %v limit %v", tc.name, got, tc.content, tc.total, tc.offset, tc.limit)
		}
	}
	for args, why := range map[string]string{
		`{"offset":0}`:   "offset is 1-based and must be >= 1 (got 0)",
		`{"offset":-1}`:  "offset is 1-based and must be >= 1 (got -1)",
		`{"limit":0}`:    "limit must be >= 1 (got 0)",
		`{"offset":1.5}`: "must be an integer",
	} {
		var extra map[string]any
		if err := json.Unmarshal([]byte(args), &extra); err != nil {
			t.Fatal(err)
		}
		if got := read(extra); got["ok"] != false || !strings.Contains(fmt.Sprint(got["refused"]), why) {
			t.Fatalf("read %s = %v, want refused %q", args, got, why)
		}
	}
}

// TestPolicyToolRefusesCaseVariantRangeKeysItself is the premise of leaving
// offset/limit out of policyToolRequestKeys: the handler passes a case variant
// through (it guards only op and flags), and v0.54.15 refuses it by name
// instead of reading it as a range. If upstream ever decodes keys loosely
// again, this reds and the keys must join the guard.
func TestPolicyToolRefusesCaseVariantRangeKeysItself(t *testing.T) {
	f := newRealFixture(t, "worktree")
	if err := os.WriteFile(filepath.Join(f.root, "a.txt"), []byte("l1\nl2\nl3\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"Offset", "OFFSET", "Limit", "LIMIT"} {
		got := f.op(t, EffectWorkspaceRead, map[string]any{"op": "read", "path": "a.txt", key: 2})
		want := fmt.Sprintf("unknown request key %q", key)
		if got["ok"] != false || !strings.Contains(fmt.Sprint(got["refused"]), want) || got["content"] != nil {
			t.Fatalf("read with %q = %v, want policy-tool's %s", key, got, want)
		}
	}
}
