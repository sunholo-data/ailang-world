package broker

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func execReadOp(t *testing.T, h *AilangToolHandler, payload map[string]any) map[string]any {
	t.Helper()
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := h.Execute(boundedTestContext(t), EffectRequest{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1,
	}, rawPayload)
	if err != nil {
		t.Fatalf("Execute(%s): %v", rawPayload, err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal(%s): %v", raw, err)
	}
	return out
}

func assertThreePartRefusal(t *testing.T, out map[string]any, substr string) {
	t.Helper()
	if out["ok"] != false {
		t.Fatalf("got ok=%v, want false (%v)", out["ok"], out)
	}
	refused, _ := out["refused"].(string)
	if !strings.HasPrefix(refused, "refused: ") || !strings.Contains(refused, "\n  why: ") || !strings.Contains(refused, "\n  next: ") {
		t.Fatalf("not a three-part refusal: %q", refused)
	}
	if substr != "" && !strings.Contains(refused, substr) {
		t.Fatalf("refusal %q does not contain %q", refused, substr)
	}
}

func TestWorkspaceSearchCoreAndWindowing(t *testing.T) {
	h, _, _ := newFakeToolHandler(t, "ok")
	root := h.root

	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Normal files, binary file, .git file, and a huge single-line file.
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte("package a\n// needle here\nfunc Alpha() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "b.md"), []byte("hello NEEDLE world\nsecond line\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin.dat"), []byte("abc\x00needle\x00def"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "HEAD"), []byte("needle inside git\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hugeLine := strings.Repeat("x", 3000) + "ROW228_MARKER" + strings.Repeat("y", 3000) + "\n"
	if err := os.WriteFile(filepath.Join(root, "charter.md"), []byte("header\n"+hugeLine), 0o644); err != nil {
		t.Fatal(err)
	}

	// Fixed string search skips binary and .git.
	out := execReadOp(t, h, map[string]any{"op": "search", "pattern": "needle"})
	if out["ok"] != true {
		t.Fatalf("search failed: %v", out)
	}
	matches, _ := out["matches"].([]any)
	if len(matches) != 1 {
		t.Fatalf("got %d matches, want 1: %v", len(matches), matches)
	}
	m0 := matches[0].(map[string]any)
	if m0["path"] != "a.go" || m0["line"] != float64(2) || m0["col"] != float64(4) {
		t.Fatalf("unexpected match: %v", m0)
	}
	if out["truncated"] != false {
		t.Fatalf("unexpected truncation: %v", out)
	}

	// ignore_case + include glob.
	out = execReadOp(t, h, map[string]any{"op": "search", "pattern": "needle", "ignore_case": true, "include": "*.md"})
	matches = out["matches"].([]any)
	if len(matches) != 1 || matches[0].(map[string]any)["path"] != "sub/b.md" {
		t.Fatalf("ignore_case + include got %v", matches)
	}

	// RE2 regex search on a single huge line: window is <= 240 runes with ellipses.
	out = execReadOp(t, h, map[string]any{"op": "search", "pattern": "ROW[0-9]+_MARKER", "regex": true, "path": "charter.md"})
	matches = out["matches"].([]any)
	if len(matches) != 1 {
		t.Fatalf("regex search got %v", out)
	}
	text := matches[0].(map[string]any)["text"].(string)
	if utf8.RuneCountInString(text) > searchWindowRunes {
		t.Fatalf("window too long: %d runes", utf8.RuneCountInString(text))
	}
	if !strings.HasPrefix(text, "…") || !strings.HasSuffix(text, "…") || !strings.Contains(text, "ROW228_MARKER") {
		t.Fatalf("window malformed: %q", text)
	}
}

func TestWorkspaceSearchCapsAndConfinement(t *testing.T) {
	h, _, _ := newFakeToolHandler(t, "ok")
	root := h.root

	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("SECRET_OUTSIDE_TOKEN\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideFile, filepath.Join(root, "escape_file.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideDir, filepath.Join(root, "escape_dir")); err != nil {
		t.Fatal(err)
	}

	// Walking root never reads through symlinks to outside files/dirs.
	out := execReadOp(t, h, map[string]any{"op": "search", "pattern": "SECRET_OUTSIDE_TOKEN"})
	if out["ok"] != true || len(out["matches"].([]any)) != 0 {
		t.Fatalf("symlink escaped sandbox: %v", out)
	}

	// Direct path escapes, symlinks, .git, unknown args, and bad values all refuse in 3 parts.
	refusals := []struct {
		payload map[string]any
		want    string
	}{
		{map[string]any{"op": "search", "pattern": "x", "path": "../outside"}, ".."},
		{map[string]any{"op": "search", "pattern": "x", "path": "/etc/hosts"}, "sandbox"},
		{map[string]any{"op": "search", "pattern": "x", "path": "-flag"}, "sandbox"},
		{map[string]any{"op": "search", "pattern": "x", "path": ".git"}, ".git"},
		{map[string]any{"op": "search", "pattern": "x", "path": ".GIT/HEAD"}, ".git"},
		{map[string]any{"op": "search", "pattern": "x", "path": "escape_file.txt"}, "symlink"},
		{map[string]any{"op": "search", "pattern": "x", "path": "escape_dir/secret.txt"}, "symlink"},
		{map[string]any{"op": "search", "pattern": "x", "exclude": "*.go"}, "unknown argument"},
		{map[string]any{"op": "search", "pattern": "x", "PATH": "."}, "unknown argument"},
		{map[string]any{"op": "search", "pattern": ""}, "1 to 512"},
		{map[string]any{"op": "search", "pattern": "[unclosed", "regex": true}, "invalid regex"},
		{map[string]any{"op": "search", "pattern": "x", "include": "[bad"}, "invalid include"},
		{map[string]any{"op": "search", "pattern": "x", "max_results": 501}, "out of range"},
	}
	for _, tc := range refusals {
		assertThreePartRefusal(t, execReadOp(t, h, tc.payload), tc.want)
	}

	// max_results cap sets truncated:true and next.
	var sb strings.Builder
	for i := 0; i < 30; i++ {
		sb.WriteString("hit_line\n")
	}
	if err := os.WriteFile(filepath.Join(root, "many.txt"), []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	out = execReadOp(t, h, map[string]any{"op": "search", "pattern": "hit_line", "path": "many.txt", "max_results": 5})
	if out["ok"] != true || out["truncated"] != true || len(out["matches"].([]any)) != 5 || out["next"] != "narrow path or include" {
		t.Fatalf("max_results cap failed: %v", out)
	}

	// Whole-result 40,000-char cap truncates before exceeding defaultReadMaxChars.
	var big strings.Builder
	for i := 0; i < 400; i++ {
		big.WriteString("CAPHIT_" + strings.Repeat("z", 220) + "\n")
	}
	if err := os.WriteFile(filepath.Join(root, "big.txt"), []byte(big.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	raw, err := h.Execute(boundedTestContext(t), EffectRequest{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1,
	}, []byte(`{"op":"search","pattern":"CAPHIT_","path":"big.txt","max_results":500}`))
	if err != nil {
		t.Fatal(err)
	}
	if utf8.RuneCount(raw) > defaultReadMaxChars {
		t.Fatalf("search response exceeded %d runes: got %d", defaultReadMaxChars, utf8.RuneCount(raw))
	}
	var capOut map[string]any
	_ = json.Unmarshal(raw, &capOut)
	if capOut["ok"] != true || capOut["truncated"] != true || capOut["next"] != "narrow path or include" {
		t.Fatalf("expected truncated cap response, got %v", capOut)
	}
}

func TestWorkspaceListBoundedSymlinkSafeAndRecorded(t *testing.T) {
	h, _, _ := newFakeToolHandler(t, "ok")
	root := h.root

	outsideDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outsideDir, "outside.txt"), []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".git", "objects"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "HEAD"), []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "a", "b", "c"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "a", "one.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "a", "b", "c", "deep.go"), []byte("package c\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideDir, filepath.Join(root, "linkdir")); err != nil {
		t.Fatal(err)
	}

	// depth=6 lists linkdir as symlink with NO entries beneath it, and skips .git entirely.
	out := execReadOp(t, h, map[string]any{"op": "list", "path": ".", "depth": 6})
	if out["ok"] != true || out["truncated"] != false {
		t.Fatalf("list failed: %v", out)
	}
	entries := out["entries"].([]any)
	sawLink := false
	prev := ""
	for _, raw := range entries {
		e := raw.(map[string]any)
		p := e["path"].(string)
		if p < prev {
			t.Fatalf("entries not sorted: %q after %q", p, prev)
		}
		prev = p
		if strings.HasPrefix(p, ".git") {
			t.Fatalf("list leaked .git entry: %q", p)
		}
		if strings.HasPrefix(p, "linkdir/") {
			t.Fatalf("list followed symlink: %q", p)
		}
		if p == "linkdir" {
			sawLink = true
			if e["type"] != "symlink" {
				t.Fatalf("linkdir type=%v, want symlink", e["type"])
			}
		}
	}
	if !sawLink {
		t.Fatalf("expected linkdir symlink in entries: %v", entries)
	}

	// max_entries cap returns exact count and truncated:true.
	out = execReadOp(t, h, map[string]any{"op": "list", "path": ".", "depth": 6, "max_entries": 2})
	if out["ok"] != true || out["truncated"] != true || len(out["entries"].([]any)) != 2 {
		t.Fatalf("max_entries cap failed: %v", out)
	}

	// include filter.
	out = execReadOp(t, h, map[string]any{"op": "list", "path": ".", "depth": 6, "include": "*.go"})
	entries = out["entries"].([]any)
	if len(entries) != 2 || entries[0].(map[string]any)["path"] != "a/b/c/deep.go" || entries[1].(map[string]any)["path"] != "a/one.go" {
		t.Fatalf("include filter got %v", entries)
	}

	// Refusals: path escape, symlink target, .git, unknown key, out-of-range depth/max_entries.
	for _, tc := range []struct {
		payload map[string]any
		want    string
	}{
		{map[string]any{"op": "list", "path": "../x"}, ".."},
		{map[string]any{"op": "list", "path": "/tmp"}, "sandbox"},
		{map[string]any{"op": "list", "path": "linkdir"}, "symlink"},
		{map[string]any{"op": "list", "path": ".git"}, ".git"},
		{map[string]any{"op": "list", "recursive": true}, "unknown argument"},
		{map[string]any{"op": "list", "DEPTH": 2}, "unknown argument"},
		{map[string]any{"op": "list", "depth": 7}, "out of range"},
		{map[string]any{"op": "list", "max_entries": 1001}, "out of range"},
	} {
		assertThreePartRefusal(t, execReadOp(t, h, tc.payload), tc.want)
	}

	// AC8: both search and list commit one Workspace.Read effect (cost 1) through Broker.
	recording := &publishRecordingStore{base: openTestStore(t)}
	session := newSession(recording, "search-list-record", []Capability{{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, ExpiresAt: 100, Budget: 2,
	}}, Registry{EffectWorkspaceRead: h}, Live, nil)
	for i, payload := range []string{`{"op":"list","path":"."}`, `{"op":"search","pattern":"package"}`} {
		res, ref, err := session.Invoke(boundedTestContext(t), EffectRequest{
			Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1, Now: 1,
		}, []byte(payload))
		if err != nil || ref.IsZero() || !strings.Contains(string(res), `"ok":true`) {
			t.Fatalf("call %d failed: res=%s ref=%s err=%v", i, res, ref, err)
		}
		receipt := effectReceipt(t, recording, i)
		if receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "succeeded" {
			t.Fatalf("call %d not recorded as succeeded: %+v", i, receipt)
		}
	}
	// Budget of 2 is now exhausted.
	if _, _, err := session.Invoke(boundedTestContext(t), EffectRequest{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1, Now: 1,
	}, []byte(`{"op":"list"}`)); err == nil {
		t.Fatal("expected Workspace.Read budget exhaustion after 2 calls")
	}
}
