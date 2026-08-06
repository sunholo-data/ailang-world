// Package boundary holds executable tests for dependencies that cross host
// package boundaries. It intentionally contains no production code.
package boundary

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type goGroup struct {
	name, pattern, dir, mutantFile, mutantImport string
	// extraForbidden holds import prefixes rejected for THIS group only, on top
	// of forbiddenImportPrefixes. It exists because the loopback exception below
	// is true of exactly one group, and a single shared list silently granted it
	// to all three: measured at the SM.B2a boundary, bare "net/http" blank-imported
	// into host/store/store.go left this gate green (rc=0) while the
	// "net/http/httputil" arm correctly red-ed. Baseline net/http presence in each
	// group's dependency closure is host/store 0, host/replay 0,
	// cmd/ailang-worldd 1 — so store and replay were exempt from a transport they
	// never used, in the gate whose whole purpose is confining network to
	// host/broker.
	extraForbidden []string
}

var protectedGoGroups = []goGroup{
	{"host/store", "./host/store/...", "host/store", "host/store/store.go", "net/http/httputil", []string{"net/http"}},
	{"host/replay", "./host/replay/...", "host/replay", "host/replay/replay.go", "net/http/httputil", []string{"net/http"}},
	// ailang-worldd already uses net/http for its loopback-only daemon client.
	// Keep that inherited transport visible as a narrow exception while rejecting
	// an additional HTTP surface (and all direct registry/cloud imports). This is
	// the ONLY group the exception is true of, hence the empty extraForbidden.
	{"cmd/ailang-worldd", "./cmd/ailang-worldd/...", "cmd/ailang-worldd", "cmd/ailang-worldd/main.go", "net/http/httputil", nil},
}

var forbiddenImportPrefixes = []string{
	"cloud.google.com/",
	"github.com/Azure/",
	"github.com/aws/",
	"github.com/sunholo-data/ailang-world/host/registry",
	"net/http/httptest",
	"net/http/httputil",
}

var forbiddenAILMarkers = []string{".ailang/cache", "package-cache", "package_cache", "http://", "https://"}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// goListDeps deliberately mirrors host/broker's bounded go-list helper. Go
// dependency closures and world/*.ail sources are different enumerations: Go
// can enumerate the former, but cannot see AILANG modules at all.
func goListDeps(root string, patterns ...string) ([]string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("cannot locate the `go` toolchain on PATH: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goBin, append([]string{"list", "-deps"}, patterns...)...)
	cmd.Dir = root
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s list -deps %s: %w\nstderr: %s", goBin, strings.Join(patterns, " "), err, stderr.String())
	}
	var deps []string
	s := bufio.NewScanner(strings.NewReader(stdout.String()))
	for s.Scan() {
		if line := strings.TrimSpace(s.Text()); line != "" {
			deps = append(deps, line)
		}
	}
	return deps, s.Err()
}

func enumerateAIL(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(filepath.Join(root, "world"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".ailang" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".ail" {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}

func forbiddenImport(path string, extra ...string) bool {
	for _, prefix := range append(append([]string{}, forbiddenImportPrefixes...), extra...) {
		if path == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func checkGoGroup(root string, group goGroup) error {
	deps, err := goListDeps(root, group.pattern)
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		return fmt.Errorf("%s dependency enumeration is empty: guard would pass vacuously", group.name)
	}
	return filepath.WalkDir(filepath.Join(root, group.dir), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// go list -deps (without -test) describes production dependencies, so the
		// attribution scan must use the same scope.
		if d.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		for _, spec := range parsed.Imports {
			imp, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return fmt.Errorf("parse import in %s: %w", path, err)
			}
			if forbiddenImport(imp, group.extraForbidden...) {
				rel, _ := filepath.Rel(root, path)
				return fmt.Errorf("%s: forbidden registry/HTTP/cloud dependency %q", filepath.ToSlash(rel), imp)
			}
		}
		return nil
	})
}

func checkAILGroup(root string) error {
	files, err := enumerateAIL(root)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("world AILANG enumeration is empty: guard would pass vacuously")
	}
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return err
		}
		lower := strings.ToLower(string(data))
		for _, marker := range forbiddenAILMarkers {
			if strings.Contains(lower, marker) {
				return fmt.Errorf("%s: forbidden registry/package-cache marker %q", rel, marker)
			}
		}
	}
	return nil
}

func digest(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

func mutateAndRestore(t *testing.T, root, rel string, mutate func([]byte) []byte, check func() error) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	baseline, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	baseSHA := digest(baseline)
	restored := false
	defer func() {
		if !restored {
			_ = os.WriteFile(path, baseline, info.Mode().Perm())
		}
	}()
	mutant := mutate(baseline)
	if err := os.WriteFile(path, mutant, info.Mode().Perm()); err != nil {
		t.Fatal(err)
	}
	mutantSHA := digest(mutant)
	if mutantSHA == baseSHA {
		t.Fatalf("mutation did not apply to %s: sha256 remained %s", rel, baseSHA)
	}
	err = check()
	if err == nil {
		t.Fatalf("mutation in %s passed boundary guard", rel)
	}
	if !strings.Contains(err.Error(), rel) {
		t.Fatalf("guard failure did not name exact path %s: %v", rel, err)
	}
	t.Logf("MUTATION path=%s baseline_sha256=%s mutant_sha256=%s guard=%q", rel, baseSHA, mutantSHA, err)
	if err := os.WriteFile(path, baseline, info.Mode().Perm()); err != nil {
		t.Fatal(err)
	}
	restoredBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	restoredSHA := digest(restoredBytes)
	if restoredSHA != baseSHA {
		t.Fatalf("restore mismatch for %s: got %s want %s", rel, restoredSHA, baseSHA)
	}
	restored = true
	t.Logf("RESTORE path=%s restored_sha256=%s byte_identical=true", rel, restoredSHA)
}

func TestWorldBoundaryDependencyAllowlist(t *testing.T) {
	root := repoRoot(t)
	// Enumerate and print every protected group before any mutation. Go groups
	// are dependency closures; world is the distinct AILANG source enumeration.
	for _, group := range protectedGoGroups {
		deps, err := goListDeps(root, group.pattern)
		if err != nil {
			t.Fatal(err)
		}
		if len(deps) == 0 {
			t.Fatalf("%s dependency enumeration is empty", group.name)
		}
		t.Logf("ENUMERATION group=%s exact_count=%d dependencies=%q", group.name, len(deps), deps)
		if err := checkGoGroup(root, group); err != nil {
			t.Fatal(err)
		}
	}
	ailFiles, err := enumerateAIL(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ailFiles) == 0 {
		t.Fatal("world AILANG enumeration is empty")
	}
	t.Logf("ENUMERATION group=world exact_count=%d files=%q", len(ailFiles), ailFiles)
	if err := checkAILGroup(root); err != nil {
		t.Fatal(err)
	}

	for _, group := range protectedGoGroups {
		group := group
		t.Run("mutation_"+strings.ReplaceAll(group.name, "/", "_"), func(t *testing.T) {
			mutateAndRestore(t, root, group.mutantFile, func(baseline []byte) []byte {
				anchor := []byte("import (\n")
				if bytes.Count(baseline, anchor) != 1 {
					t.Fatalf("mutation anchor count for %s is %d, want 1", group.mutantFile, bytes.Count(baseline, anchor))
				}
				insert := []byte(fmt.Sprintf("import (\n\t_ %q // boundary mutation: compiling HTTP import\n", group.mutantImport))
				return bytes.Replace(baseline, anchor, insert, 1)
			}, func() error { return checkGoGroup(root, group) })
		})
	}
	t.Run("mutation_world", func(t *testing.T) {
		mutateAndRestore(t, root, "world/types.ail", func(baseline []byte) []byte {
			return append(append([]byte(nil), baseline...), []byte("\n-- boundary mutation: package-cache .ailang/cache\n")...)
		}, func() error { return checkAILGroup(root) })
	})

	// Green control: the broker is intentionally the network boundary. Its
	// non-empty closure must remain permitted by this inverse protected-group guard.
	broker, err := goListDeps(root, "./host/broker/...")
	if err != nil {
		t.Fatal(err)
	}
	if len(broker) == 0 {
		t.Fatal("host/broker dependency enumeration is empty")
	}
	t.Logf("GREEN_CONTROL group=host/broker exact_count=%d result=PASS", len(broker))
}

func TestWorldBoundaryNullCases(t *testing.T) {
	if err := func(files []string) error {
		if len(files) == 0 {
			return fmt.Errorf("empty")
		}
		return nil
	}(nil); err == nil {
		t.Fatal("empty AILANG enumeration passed vacuously")
	}
	if forbiddenImport("fmt") {
		t.Fatal("stdlib control unexpectedly forbidden")
	}
	if !forbiddenImport("net/http/httputil") {
		t.Fatal("HTTP mutation control unexpectedly permitted")
	}
}

// TestBareNetHTTPExemptionIsPerGroup pins the asymmetry that the single shared
// forbiddenImportPrefixes list used to erase. The loopback-daemon exception is
// true of cmd/ailang-worldd and of nothing else, so bare "net/http" must be
// rejected for host/store and host/replay while remaining permitted for
// cmd/ailang-worldd. Collapsing extraForbidden back into one global list makes
// this test fail rather than silently re-granting the exemption.
func TestBareNetHTTPExemptionIsPerGroup(t *testing.T) {
	byName := make(map[string]goGroup, len(protectedGoGroups))
	for _, g := range protectedGoGroups {
		byName[g.name] = g
	}
	if len(byName) != 3 {
		t.Fatalf("protected group enumeration is %d, want 3: the guard below would not cover what it claims", len(byName))
	}

	for _, want := range []struct {
		group     string
		forbidden bool
	}{
		{"host/store", true},
		{"host/replay", true},
		{"cmd/ailang-worldd", false}, // documented loopback-IPC exception
	} {
		g, ok := byName[want.group]
		if !ok {
			t.Fatalf("protected group %q absent: enumeration no longer covers it", want.group)
		}
		if got := forbiddenImport("net/http", g.extraForbidden...); got != want.forbidden {
			t.Errorf("forbiddenImport(\"net/http\") for %s = %v, want %v", want.group, got, want.forbidden)
		}
		// net/http/httputil stays forbidden everywhere, per the shared list.
		if !forbiddenImport("net/http/httputil", g.extraForbidden...) {
			t.Errorf("net/http/httputil unexpectedly permitted for %s", want.group)
		}
	}
}
