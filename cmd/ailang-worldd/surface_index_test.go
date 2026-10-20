package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Row 188 M0a: `surface-index`. Building needs git, so these arms cover the
// flags, the refusals made before git runs, and --validate on the golden.

const surfaceGolden = "../../host/surface/testdata/index.json"

func runSurface(args ...string) (int, string, string) {
	var out, errw bytes.Buffer
	code := run(append([]string{"surface-index"}, args...), &out, &errw)
	return code, out.String(), errw.String()
}

func TestSurfaceIndexHelp(t *testing.T) {
	for _, args := range [][]string{{"help", "surface-index"}, {"surface-index", "--help"}} {
		var out, errw bytes.Buffer
		if code := run(args, &out, &errw); code != exitOK || out.String() != surfaceIndexHelp {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", args, code, out.String(), errw.String())
		}
	}
	for _, flag := range []string{"--repo", "--charter", "--out", "--rev", "--limit", "--validate"} {
		if !strings.Contains(surfaceIndexHelp, "  "+flag+" ") {
			t.Errorf("help does not describe %s", flag)
		}
	}
	if !strings.Contains(usage, "\n  ailang-worldd surface-index --repo <dir> --charter <file> --out <file>") ||
		!strings.Contains(usage, "\n  ailang-worldd surface-index --validate <file>") {
		t.Error("the top-level usage does not list surface-index")
	}
	if err := checkWorlddLine("ailang-worldd surface-index --repo . --charter c.md --out /tmp/i.json --rev dev --limit 50"); err != nil {
		t.Errorf("help check: %v", err)
	}
}

func TestSurfaceIndexFlagRefusals(t *testing.T) {
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "index.json")
	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, "--repo, --charter and --out are required"},
		{[]string{"--repo", repo, "--charter", "c.md"}, "--repo, --charter and --out are required"},
		{[]string{"--validate", surfaceGolden, "--repo", repo}, "--validate takes no other flag"},
		{[]string{"--validate", surfaceGolden, "--rev", "dev"}, "--validate takes no other flag"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--limit", "0"}, "--limit 0 is outside 1..2500"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--limit", "2501"}, "--limit 2501 is outside 1..2500"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--rev", "--upload-pack=x"}, "is not a revision"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "extra"}, `unexpected argument "extra"`},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", filepath.Join(repo, "sub", "index.json")}, "is inside --repo"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", repo}, "is inside --repo"},
		{[]string{"--repo", filepath.Join(repo, "nope"), "--charter", "c.md", "--out", outside}, "--repo"},
		{[]string{"--repo", repo, "--charter", filepath.Join(repo, "nope.md"), "--out", outside}, "--charter"},
	} {
		code, stdout, stderr := runSurface(tc.args...)
		if code != exitUsage || stdout != "" || !strings.Contains(stderr, "refused: ") || !strings.Contains(stderr, tc.want) {
			t.Errorf("%v: exit %d, stdout %q, stderr %q; want %q", tc.args, code, stdout, stderr, tc.want)
		}
	}
	if _, err := os.Stat(outside); err == nil {
		t.Error("a refused build wrote --out")
	}
	var out, errw bytes.Buffer
	if code := run([]string{"--addr", "127.0.0.1:1", "surface-index", "--validate", surfaceGolden}, &out, &errw); code != exitUsage {
		t.Errorf("--addr surface-index: exit %d", code)
	}
}

func TestSurfaceIndexValidate(t *testing.T) {
	code, stdout, stderr := runSurface("--validate", surfaceGolden)
	if code != exitOK || !strings.Contains(stdout, "24 shared states, 24 steps, 6 branes") {
		t.Fatalf("golden: exit %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	raw, err := os.ReadFile(surfaceGolden)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, tc := range map[string]struct{ from, to, want string }{
		"duplicate id": {`"id": "f36428abd9a72e5093b61cba14e716c37019f1ef"`, `"id": "90035557b2012895551867c0ebe2f36b853aac09"`, "duplicates"},
		"bad kind":     {`"kind": "loop"`, `"kind": "nightly"`, `kind "nightly"`},
		"schema":       {surfaceSchemaLine, `"schema": "world/surface-index/v0"`, "schema"},
		"unknown key":  {surfaceSchemaLine, surfaceSchemaLine + `, "path": "x"`, "does not decode"},
	} {
		if !strings.Contains(string(raw), tc.from) {
			t.Fatalf("%s: the golden has no %q", name, tc.from)
		}
		p := filepath.Join(dir, strings.ReplaceAll(name, " ", "-")+".json")
		if err := os.WriteFile(p, []byte(strings.Replace(string(raw), tc.from, tc.to, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		code, stdout, stderr := runSurface("--validate", p)
		if code == exitOK || stdout != "" || !strings.Contains(stderr, tc.want) {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want %q", name, code, stdout, stderr, tc.want)
		}
	}
	if code, _, stderr := runSurface("--validate", filepath.Join(dir, "missing.json")); code == exitOK || !strings.Contains(stderr, "missing.json") {
		t.Errorf("missing file: exit %d, %q", code, stderr)
	}
}

const surfaceSchemaLine = `"schema": "world/surface-index/v1"`

func TestWriteFileAtomic(t *testing.T) {
	p := filepath.Join(t.TempDir(), "d", "index.json")
	for _, body := range []string{"one\n", "two\n"} {
		if err := writeFileAtomic(p, []byte(body)); err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(p); string(got) != body {
			t.Errorf("read %q, want %q", got, body)
		}
	}
	ents, _ := os.ReadDir(filepath.Dir(p))
	if len(ents) != 1 {
		t.Errorf("temporary files left: %v", ents)
	}
}
