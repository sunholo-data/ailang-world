package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/appregistry"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
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
	t.Setenv("HOME", t.TempDir()) // never read the operator's real world store
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
		!strings.Contains(out, "grant   Workspace.Exec=worktree:100") || !strings.Contains(out, "state   local (not published)") ||
		!strings.Contains(out, "NOT PUBLISHED docparse-convert: only published apps launch (D-WORLD-129)") {
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
		!strings.Contains(errOut, "no app missing in the registry or "+dir+" (have: docparse-convert, wide)") || !strings.Contains(errOut, "next:") {
		t.Fatalf("app show missing = %d %q", code, errOut)
	}
	if code, _, errOut := runApp("show", "Bad_Name", "--apps-dir", dir); code != exitUsage || !strings.Contains(errOut, "invalid app name") {
		t.Fatalf("app show Bad_Name = %d %q", code, errOut)
	}
	if code, _, errOut := runApp("install", "x"); code != exitUsage || !strings.Contains(errOut, "app verbs are list, show and check") {
		t.Fatalf("app install = %d %q", code, errOut)
	}
	if code, out, _ := runApp("list", "--apps-dir", filepath.Join(dir, "absent")); code != exitOK || !strings.Contains(out, "no world store at ") || !strings.Contains(out, "no apps (catalog ") {
		t.Fatalf("app list on a missing catalog = %d %q", code, out)
	}
	if code, out, _ := runApp("check", "wide", "--apps-dir", dir, "--json"); code != exitUsage || !strings.Contains(out, `"valid": false`) {
		t.Fatalf("app check --json = %d %q", code, out)
	}
}

// Row 274 M1b-3: app list/show read the app registry (read-only, beside a writer) and
// merge it with the local catalog.
func TestAppVerbsReadTheRegistry(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	example, err := os.ReadFile(filepath.Join("..", "..", "docs", "examples", "apps", "docparse-convert.json"))
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dbPath := filepath.Join(base, "world.db")
	writer, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = writer.Close() }() // held open: the CLI must read beside a writer
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	obj, _, err := appregistry.ManifestObject(example)
	if err != nil {
		t.Fatal(err)
	}
	rev, err := appregistry.BuildPublish(nil, hashref.HashRef{}, "docparse-convert", obj.Hash, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := appregistry.Publish(ctx, writer, hashref.HashRef{}, rev, &obj); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(base, "apps")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	draft := strings.Replace(string(example), `"budget": 2000`, `"budget": 1999`, 1)
	if err := os.WriteFile(filepath.Join(dir, "docparse-convert.json"), []byte(draft), 0o600); err != nil {
		t.Fatal(err)
	}
	runApp := func(args ...string) (int, string, string) {
		var out, errOut bytes.Buffer
		code := run(append([]string{"app"}, args...), &out, &errOut)
		return code, out.String(), errOut.String()
	}
	code, out, errOut := runApp("list", "--db", dbPath, "--apps-dir", dir)
	if code != exitOK || !strings.Contains(out, "app registry: revision 1 in "+dbPath) ||
		!strings.Contains(out, "published v1, local copy differs") {
		t.Fatalf("app list = %d\n%s%s", code, out, errOut)
	}
	code, out, _ = runApp("show", "docparse-convert", "--db", dbPath, "--apps-dir", dir)
	if code != exitOK || !strings.Contains(out, "pinned  "+obj.Hash.String()) ||
		!strings.Contains(out, "grant   Workspace.Read=worktree:2000") || !strings.Contains(out, "LAUNCHABLE docparse-convert v1") {
		t.Fatalf("app show = %d\n%s (the pinned manifest, budget 2000, is what launches)", code, out)
	}
	code, out, _ = runApp("list", "--db", dbPath, "--apps-dir", dir, "--json")
	if code != exitOK || !strings.Contains(out, `"state": "published-local-differs"`) || !strings.Contains(out, `"launchable": true`) {
		t.Fatalf("app list --json = %d\n%s", code, out)
	}
	// No store file at --db: no registry, every app local, and nothing is created.
	missing := filepath.Join(base, "none.db")
	code, out, _ = runApp("list", "--db", missing, "--apps-dir", dir)
	if code != exitOK || !strings.Contains(out, "no world store at "+missing) || !strings.Contains(out, "local (not published)") {
		t.Fatalf("app list with no store = %d\n%s", code, out)
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatal("app list created a world store")
	}
}
