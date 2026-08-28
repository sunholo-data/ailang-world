package verifygate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type subprocessSinkViolation struct {
	position token.Position
	variable string
}

func bareSinkType(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok || (selector.Sel.Name != "Buffer" && selector.Sel.Name != "Builder") {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && ((pkg.Name == "bytes" && selector.Sel.Name == "Buffer") ||
		(pkg.Name == "strings" && selector.Sel.Name == "Builder"))
}

func inspectSubprocessFunction(fset *token.FileSet, body *ast.BlockStmt) ([]subprocessSinkViolation, bool) {
	declared := make(map[string]bool)
	hasStart := false
	ast.Inspect(body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.FuncLit:
			// A function literal is inspected separately from its enclosing function.
			return node.Body == body
		case *ast.ValueSpec:
			if node.Type != nil && bareSinkType(node.Type) {
				for _, name := range node.Names {
					declared[name.Name] = true
				}
			}
		case *ast.CallExpr:
			if selector, ok := node.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Start" {
				hasStart = true
			}
		}
		return true
	})
	if !hasStart {
		return nil, false
	}

	var violations []subprocessSinkViolation
	ast.Inspect(body, func(node ast.Node) bool {
		if literal, ok := node.(*ast.FuncLit); ok && literal.Body != body {
			return false
		}
		assignment, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for i, lhs := range assignment.Lhs {
			if i >= len(assignment.Rhs) {
				break
			}
			target, ok := lhs.(*ast.SelectorExpr)
			if !ok || (target.Sel.Name != "Stderr" && target.Sel.Name != "Stdout") {
				continue
			}
			address, ok := assignment.Rhs[i].(*ast.UnaryExpr)
			if !ok || address.Op != token.AND {
				continue
			}
			name, ok := address.X.(*ast.Ident)
			if ok && declared[name.Name] {
				violations = append(violations, subprocessSinkViolation{
					position: fset.Position(lhs.Pos()),
					variable: name.Name,
				})
			}
		}
		return true
	})
	return violations, true
}

func detectSubprocessSinks(fset *token.FileSet, file *ast.File) ([]subprocessSinkViolation, []token.Position) {
	var violations []subprocessSinkViolation
	var starts []token.Position
	ast.Inspect(file, func(node ast.Node) bool {
		var body *ast.BlockStmt
		switch node := node.(type) {
		case *ast.FuncDecl:
			body = node.Body
		case *ast.FuncLit:
			body = node.Body
		default:
			return true
		}
		if body == nil {
			return true
		}
		found, hasStart := inspectSubprocessFunction(fset, body)
		if hasStart {
			starts = append(starts, fset.Position(body.Pos()))
			violations = append(violations, found...)
		}
		return true
	})
	return violations, starts
}

func subprocessGateRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		} else if !os.IsNotExist(err) {
			t.Fatalf("inspect %s for go.mod: %v", dir, err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find go.mod walking up from test working directory %s", dir)
		}
		dir = parent
	}
}

func TestSubprocessSinkGate(t *testing.T) {
	t.Run("detector controls", func(t *testing.T) {
		fixtures := []struct {
			name string
			src  string
			want int
		}{
			{
				name: "positive bare buffer",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { var buf bytes.Buffer; cmd := exec.Command("helper"); cmd.Stderr = &buf; _ = cmd.Start() }
`,
				want: 1,
			},
			{
				name: "negative guarded sink",
				src: `package fixture
import "os/exec"
type syncBuf struct{}
func (*syncBuf) Write(p []byte) (int, error) { return len(p), nil }
func test() { cmd := exec.Command("helper"); cmd.Stderr = &syncBuf{}; _ = cmd.Start() }
`,
				want: 0,
			},
		}
		for _, fixture := range fixtures {
			t.Run(fixture.name, func(t *testing.T) {
				fset := token.NewFileSet()
				file, err := parser.ParseFile(fset, fixture.name+".go", fixture.src, 0)
				if err != nil {
					t.Fatalf("DETECTOR is broken: parse control: %v", err)
				}
				violations, _ := detectSubprocessSinks(fset, file)
				if len(violations) != fixture.want {
					t.Fatalf("DETECTOR is broken: %s produced %d violations, want %d", fixture.name, len(violations), fixture.want)
				}
			})
		}
	})

	root := subprocessGateRepoRoot(t)
	filesParsed := 0
	var violations []subprocessSinkViolation
	var startSites []token.Position
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && (entry.Name() == ".git" || entry.Name() == "vendor" || entry.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		filesParsed++
		found, starts := detectSubprocessSinks(fset, file)
		violations = append(violations, found...)
		startSites = append(startSites, starts...)
		return nil
	})
	if err != nil {
		t.Fatalf("walk and parse repository Go source: %v", err)
	}

	startFilesSet := make(map[string]bool)
	for _, site := range startSites {
		rel, err := filepath.Rel(root, site.Filename)
		if err != nil {
			t.Fatal(err)
		}
		startFilesSet[filepath.ToSlash(rel)] = true
	}
	startFiles := make([]string, 0, len(startFilesSet))
	for file := range startFilesSet {
		startFiles = append(startFiles, file)
	}
	sort.Strings(startFiles)
	t.Logf("subprocess sink gate: violations=%d filesParsed=%d startSites=%d controls=positive:ok,negative:ok", len(violations), filesParsed, len(startSites))
	if filesParsed < 50 {
		t.Fatalf("subprocess sink gate parsed only %d Go files, want at least 50", filesParsed)
	}
	if len(startSites) < 3 {
		t.Fatalf("subprocess sink gate examined only %d functions containing Start, want at least 3; files=%s", len(startSites), strings.Join(startFiles, ", "))
	}
	for _, violation := range violations {
		rel, err := filepath.Rel(root, violation.position.Filename)
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("%s:%d: subprocess sink %q is a bare bytes.Buffer or strings.Builder; an unsynchronised read can race os/exec's copier goroutine, which remains live until Wait returns; use the syncBuffer pattern in host/store/writer_lock_test.go", filepath.ToSlash(rel), violation.position.Line, violation.variable)
	}
}
