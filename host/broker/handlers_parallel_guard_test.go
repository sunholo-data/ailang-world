package broker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBrokerTestsDoNotCallParallel(t *testing.T) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the broker test package")
	}
	packageDir := filepath.Dir(source)
	entries, err := os.ReadDir(packageDir)
	if err != nil {
		t.Fatalf("read broker package directory %s: %v", packageDir, err)
	}

	fset := token.NewFileSet()
	var files []string
	anchors := map[string]bool{
		"handlers_test.go":    false,
		filepath.Base(source): false,
	}
	var helperSelectors int
	var offenders []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		files = append(files, entry.Name())
		if _, anchored := anchors[entry.Name()]; anchored {
			anchors[entry.Name()] = true
		}
		path := filepath.Join(packageDir, entry.Name())
		parsed, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			selector, ok := node.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			switch selector.Sel.Name {
			case "Helper":
				helperSelectors++
			case "Parallel":
				pos := fset.Position(selector.Sel.Pos())
				offenders = append(offenders,
					fmt.Sprintf("%s:%d:%d", filepath.Base(pos.Filename), pos.Line, pos.Column))
			}
			return true
		})
	}

	if len(files) == 0 {
		t.Fatal("broker test enumeration found zero _test.go files")
	}
	for anchor, found := range anchors {
		if !found {
			t.Fatalf("broker test enumeration missed required anchor %s", anchor)
		}
	}
	if helperSelectors == 0 {
		t.Fatal("broker test AST walk found zero Helper selectors; Parallel census is vacuous")
	}
	if len(offenders) != 0 {
		t.Fatalf("broker tests call Parallel at %v; parallel tests would race the t.Cleanup restore of the killGroup package global", offenders)
	}
	t.Logf("enumerated=%d files, Helper selectors=%d, offenders=%v", len(files), helperSelectors, offenders)
}
