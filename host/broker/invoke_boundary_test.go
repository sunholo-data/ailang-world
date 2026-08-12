package broker

// TR.C is deliberately a structural, fail-closed gate. It scans syntax rather
// than types, so any selector named Invo+ke is rejected outside host/broker,
// even when its receiver is unrelated. That over-approximation avoids a new
// go/packages dependency and makes future exceptions explicit in review.

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type bindingFinding struct {
	path, kind, enclosingFunc string
	line                      int
}

func (f bindingFinding) String() string {
	return fmt.Sprintf("%s:%d kind=%s fn=%s", f.path, f.line, f.kind, f.enclosingFunc)
}

func enclosingFunction(tree *ast.File, pos token.Pos) string {
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Pos() <= pos && pos <= fn.End() {
			return fn.Name.Name
		}
	}
	return "<file>"
}

func scanFile(rel string, tree *ast.File, fset *token.FileSet, insideBroker bool) []bindingFinding {
	invokeSel := "Invo" + "ke"
	sessType := "Sess" + "ion"
	ctorLive := "New" + "Session"
	ctorReplay := "New" + "Replay" + "Session"
	brokerPath := "github.com/sunholo-data/ailang-world/host/" + "broker"
	locals := map[string]bool{}
	var findings []bindingFinding
	add := func(pos token.Pos, kind string) {
		findings = append(findings, bindingFinding{rel, kind, enclosingFunction(tree, pos), fset.Position(pos).Line})
	}
	for _, spec := range tree.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil || path != brokerPath {
			continue
		}
		local := "broker"
		if spec.Name != nil {
			local = spec.Name.Name
		}
		if local == "." {
			add(spec.Pos(), "dot-import")
		} else if local != "_" {
			locals[local] = true
		}
	}
	ast.Inspect(tree, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if sel.Sel.Name == invokeSel {
				add(sel.Pos(), "invoke-call")
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || !locals[id.Name] {
				return true
			}
			if sel.Sel.Name == ctorLive {
				add(sel.Pos(), "ctor-live")
			}
			if sel.Sel.Name == ctorReplay {
				add(sel.Pos(), "ctor-replay")
			}
		case *ast.SelectorExpr:
			id, ok := n.X.(*ast.Ident)
			if ok && locals[id.Name] && n.Sel.Name == sessType {
				add(n.Pos(), "session-type")
			}
		}
		return true
	})
	_ = insideBroker // classification is asserted by the caller, not hidden in the detector.
	return findings
}

func parseAndScan(root, rel string, inside bool) ([]bindingFinding, error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(rel)), nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", rel, err)
	}
	return scanFile(rel, tree, fset, inside), nil
}

type walkStats struct {
	skippedTests         int
	skippedNestedModules int
}

func enumerateProductionGoFiles(root string) ([]string, walkStats, error) {
	var files []string
	var stats walkStats
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			if path != root {
				if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
					stats.skippedNestedModules++
					return filepath.SkipDir
				} else if !os.IsNotExist(err) {
					return err
				}
			}
			return nil
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			stats.skippedTests++
			return nil
		}
		if filepath.Ext(entry.Name()) != ".go" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, stats, err
}

func goListProductionFiles(t *testing.T, root string) []string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("locate go for enumeration cross-check: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goBin, "list", "-f", `{{.Dir}}|{{join .GoFiles "|"}}`, "./...")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list production files: %v", err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}
		for _, name := range parts[1:] {
			if name == "" {
				continue
			}
			rel, err := filepath.Rel(root, filepath.Join(parts[0], name))
			if err != nil {
				t.Fatalf("relativize go-list file: %v", err)
			}
			files = append(files, filepath.ToSlash(rel))
		}
	}
	sort.Strings(files)
	return files
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate binding-gate source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func assertEnumeration(t *testing.T, root string, files []string, stats walkStats) {
	t.Helper()
	if len(files) == 0 {
		t.Fatal("walked ZERO production .go files")
	}
	if len(files) < 30 {
		t.Fatalf("walked only %d production .go files, want at least 30", len(files))
	}
	walked := make(map[string]bool, len(files))
	for _, file := range files {
		walked[file] = true
		if strings.HasSuffix(file, "_test.go") {
			t.Fatalf("production walk included test file %s", file)
		}
	}
	if stats.skippedTests == 0 {
		t.Fatal("production walk skipped zero _test.go files; exclusion is vacuous")
	}
	requiredAnchors := []string{
		"host/broker/publish_op.go", "host/broker/broker.go", "host/broker/confined.go",
		"host/transitionreg/bind.go", "cmd/world-publish/main.go", "cmd/ailang-worldd/main.go",
		"host/store/writer_lock_other.go",
	}
	for _, anchor := range requiredAnchors {
		if !walked[anchor] {
			t.Fatalf("production walk missed required anchor %s", anchor)
		}
	}
	listed := goListProductionFiles(t, root)
	if len(listed) < 30 {
		t.Fatalf("go list enumerated %d files, want at least 30", len(listed))
	}
	for _, file := range listed {
		if !walked[file] {
			t.Fatalf("filesystem walk is not a go-list superset; missing %s", file)
		}
	}
	t.Logf("walked=%d skipped_tests=%d skipped_nested_modules=%d golist=%d", len(files), stats.skippedTests, stats.skippedNestedModules, len(listed))
}

func TestRegistryDispatchBindingBoundary(t *testing.T) {
	root := repositoryRoot(t)
	files, stats, err := enumerateProductionGoFiles(root)
	if err != nil {
		t.Fatalf("enumerate production Go files: %v", err)
	}
	t.Run("enumeration", func(t *testing.T) { assertEnumeration(t, root, files, stats) })

	var inside, outside []bindingFinding
	for _, rel := range files {
		isInside := strings.HasPrefix(rel, "host/broker/")
		found, err := parseAndScan(root, rel, isInside)
		if err != nil {
			t.Fatalf("repository scan: %v", err)
		}
		if isInside {
			inside = append(inside, found...)
		} else {
			outside = append(outside, found...)
		}
	}
	t.Run("outside_broker_is_clean", func(t *testing.T) {
		if len(outside) != 0 {
			t.Fatalf("raw broker binding outside host/broker (Invoke matching is deliberately receiver-type-blind):\n%v", outside)
		}
	})
	t.Run("inside_broker_exemption", func(t *testing.T) {
		const wantCount = 3
		if len(inside) != wantCount {
			t.Fatalf("inside-broker exemption count=%d want %d: %v", len(inside), wantCount, inside)
		}
		got := map[string]int{}
		for _, finding := range inside {
			got[finding.path+"|"+finding.enclosingFunc+"|"+finding.kind]++
		}
		want := map[string]int{
			"host/broker/publish_op.go|mintAttendedApproval|invoke-call":  2,
			"host/broker/publish_op.go|invokeAttendedPublish|invoke-call": 1,
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("inside-broker exemption identity mismatch: got=%v want=%v", got, want)
		}
	})
}
