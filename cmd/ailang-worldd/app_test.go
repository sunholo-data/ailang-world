package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Row 274 M1: the read-only app verbs over a temporary catalog.
func TestAppVerbs(t *testing.T) {
	example := filepath.Join("..", "..", "docs", "examples", "apps", "docparse-convert.json")
	body, err := os.ReadFile(example)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "docparse-convert.json"), body, 0o600); err != nil {
		t.Fatal(err)
	}
	bad := strings.Replace(strings.Replace(string(body), `"docparse-convert"`, `"wide"`, 1), `"scope": "worktree", "budget": 2000`, `"scope": "/", "budget": 2000`, 1)
	if err := os.WriteFile(filepath.Join(dir, "wide.json"), []byte(bad), 0o600); err != nil {
		t.Fatal(err)
	}
	runApp := func(args ...string) (int, string, string) {
		var out, errOut bytes.Buffer
		code := run(append([]string{"app"}, args...), &out, &errOut)
		return code, out.String(), errOut.String()
	}

	if code, out, _ := runApp("--help"); code != exitOK || !strings.Contains(out, "ailang-worldd app list|show|check") {
		t.Fatalf("app --help = %d %q", code, out)
	}
	code, out, _ := runApp("list", "--apps-dir", dir)
	if code != exitOK || !strings.Contains(out, "docparse-convert") || !strings.Contains(out, "Convert documents") ||
		!strings.Contains(out, "wide") || !strings.Contains(out, "INVALID 1 problem(s): ailang-worldd app check wide") {
		t.Fatalf("app list = %d\n%s", code, out)
	}
	if code, out, _ := runApp("show", "docparse-convert", "--apps-dir", dir); code != exitOK ||
		!strings.Contains(out, "grant   Workspace.Exec=worktree:100") || !strings.Contains(out, "VALID docparse-convert (sha256:") {
		t.Fatalf("app show = %d\n%s", code, out)
	}
	if code, out, _ := runApp("check", example); code != exitOK || !strings.HasPrefix(out, "VALID docparse-convert") {
		t.Fatalf("app check <file> = %d %q", code, out)
	}
	code, out, _ = runApp("check", "wide", "--apps-dir", dir)
	if code != exitUsage || !strings.Contains(out, `preset.grants[0].scope: "/" is not worktree`) || !strings.Contains(out, "next: fix each field") {
		t.Fatalf("app check wide = %d\n%s", code, out)
	}
	if code, _, errOut := runApp("show", "missing", "--apps-dir", dir); code != exitUsage ||
		!strings.Contains(errOut, "no app missing in "+dir+" (have: docparse-convert, wide)") || !strings.Contains(errOut, "next:") {
		t.Fatalf("app show missing = %d %q", code, errOut)
	}
	if code, _, errOut := runApp("show", "Bad_Name", "--apps-dir", dir); code != exitUsage || !strings.Contains(errOut, "invalid app name") {
		t.Fatalf("app show Bad_Name = %d %q", code, errOut)
	}
	if code, _, errOut := runApp("install", "x"); code != exitUsage || !strings.Contains(errOut, "app verbs are list, show and check") {
		t.Fatalf("app install = %d %q", code, errOut)
	}
	if code, out, _ := runApp("list", "--apps-dir", filepath.Join(dir, "absent")); code != exitOK || !strings.HasPrefix(out, "no apps in ") {
		t.Fatalf("app list on a missing catalog = %d %q", code, out)
	}
	if code, out, _ := runApp("check", "wide", "--apps-dir", dir, "--json"); code != exitUsage || !strings.Contains(out, `"valid": false`) {
		t.Fatalf("app check --json = %d %q", code, out)
	}
}
