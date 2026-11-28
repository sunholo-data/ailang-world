package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestQuickstartSection22RunsVerbatim runs QUICKSTART §22's command block line by line
// (coding-standards S7): every `ailang-worldd` line through the CLI, `mkdir -p` and `cp`
// as the files they make, under a temporary HOME. A line of any other shape fails, so
// the block cannot grow a step this test does not run.
func TestQuickstartSection22RunsVerbatim(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, block, ok := strings.Cut(string(raw), "<!-- app-catalog-walkthrough -->\n```bash\n")
	if !ok {
		t.Fatal("§22's walkthrough block is missing")
	}
	block, _, _ = strings.Cut(block, "```")
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := filepath.Join("..", "..")
	expand := func(p string) string {
		if strings.HasPrefix(p, "~/") {
			return filepath.Join(home, p[2:])
		}
		return filepath.Join(root, p)
	}
	ran := 0
	var outputs []string
	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
		f := strings.Fields(line)
		switch {
		case len(f) >= 2 && f[0] == "ailang-worldd":
			args := append([]string{}, f[1:]...)
			for i, a := range args {
				if strings.HasPrefix(a, "docs/") {
					args[i] = filepath.Join(root, a)
				}
			}
			var out, errOut bytes.Buffer
			if code := run(args, &out, &errOut); code != exitOK {
				t.Fatalf("%q exited %d\n%s%s", line, code, out.String(), errOut.String())
			}
			outputs = append(outputs, out.String())
		case len(f) == 3 && f[0] == "mkdir" && f[1] == "-p":
			if err := os.MkdirAll(expand(f[2]), 0o755); err != nil {
				t.Fatal(err)
			}
		case len(f) == 3 && f[0] == "cp":
			data, err := os.ReadFile(expand(f[1]))
			if err != nil {
				t.Fatal(err)
			}
			dst := expand(f[2])
			if strings.HasSuffix(f[2], "/") {
				dst = filepath.Join(dst, filepath.Base(f[1]))
			}
			if err := os.WriteFile(dst, data, 0o644); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatalf("§22 line %q is not a step this test runs", line)
		}
		ran++
	}
	if ran != 5 || len(outputs) != 3 {
		t.Fatalf("ran %d lines (%d CLI), want 5 (3 CLI)", ran, len(outputs))
	}
	if !strings.HasPrefix(outputs[0], "VALID docparse-convert (sha256:") {
		t.Errorf("app check printed %q", outputs[0])
	}
	if !strings.Contains(outputs[1], "docparse-convert") || !strings.Contains(outputs[1], "Convert documents") {
		t.Errorf("app list printed %q", outputs[1])
	}
	if !strings.Contains(outputs[2], "grant   Workspace.Read=worktree:2000") || !strings.Contains(outputs[2], "output  file out/**/*.md (checked by ailang-test)") {
		t.Errorf("app show printed %q", outputs[2])
	}
}
