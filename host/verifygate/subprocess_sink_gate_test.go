package verifygate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
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

// sinkKind records how a declared bare sink is held, because the two forms bind
// to cmd.Stderr differently: a value needs &x, a pointer is assigned directly.
type sinkKind int

const (
	sinkValue sinkKind = iota + 1
	sinkPointer
)

// bareSinkType reports whether expr names bytes.Buffer or strings.Builder — the
// two sinks in this repo that carry no internal synchronisation.
func bareSinkType(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	return (pkg.Name == "bytes" && selector.Sel.Name == "Buffer") ||
		(pkg.Name == "strings" && selector.Sel.Name == "Builder")
}

// bareSinkExpr classifies an initialiser: bytes.Buffer{} is a value, and both
// &bytes.Buffer{} and new(bytes.Buffer) are pointers. Anything else is not a
// bare sink and is ignored.
func bareSinkExpr(expr ast.Expr) (sinkKind, bool) {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		if bareSinkType(expr.Type) {
			return sinkValue, true
		}
	case *ast.UnaryExpr:
		if expr.Op == token.AND {
			if inner, ok := expr.X.(*ast.CompositeLit); ok && bareSinkType(inner.Type) {
				return sinkPointer, true
			}
		}
	case *ast.CallExpr:
		if ident, ok := expr.Fun.(*ast.Ident); ok && ident.Name == "new" && len(expr.Args) == 1 {
			if bareSinkType(expr.Args[0]) {
				return sinkPointer, true
			}
		}
	}
	return 0, false
}

// inspectSubprocessFunction analyses ONE function INCLUDING every nested function
// literal. Scoping the walk to a single function body instead would let the
// defect be reassembled across a closure boundary — declare and bind the sink in
// the outer function, call Start() inside `go func(){}()` — and evade detection
// with the pattern fully intact. Over-approximating across the subtree can only
// produce a LOUD false positive; under-approximating produces the silent miss
// this gate exists to prevent.
func inspectSubprocessFunction(fset *token.FileSet, body *ast.BlockStmt) ([]subprocessSinkViolation, bool) {
	declared := make(map[string]sinkKind)
	hasStart := false
	ast.Inspect(body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.ValueSpec:
			if node.Type != nil && bareSinkType(node.Type) {
				for _, name := range node.Names {
					declared[name.Name] = sinkValue
				}
			}
			for i, value := range node.Values {
				if i >= len(node.Names) {
					break
				}
				if kind, ok := bareSinkExpr(value); ok {
					declared[node.Names[i].Name] = kind
				}
			}
		case *ast.AssignStmt:
			// `buf := bytes.Buffer{}` and `buf := new(bytes.Buffer)` never produce
			// a ValueSpec, and := is the more idiomatic form for a local sink.
			if node.Tok != token.DEFINE {
				return true
			}
			for i, lhs := range node.Lhs {
				if i >= len(node.Rhs) {
					break
				}
				name, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if kind, ok := bareSinkExpr(node.Rhs[i]); ok {
					declared[name.Name] = kind
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
			rhs := assignment.Rhs[i]
			// An inline &bytes.Buffer{} / new(bytes.Buffer) is the same defect with
			// no name to look up.
			if _, ok := bareSinkExpr(rhs); ok {
				violations = append(violations, subprocessSinkViolation{
					position: fset.Position(lhs.Pos()),
					variable: types.ExprString(rhs),
				})
				continue
			}
			switch rhs := rhs.(type) {
			case *ast.UnaryExpr:
				if rhs.Op != token.AND {
					continue
				}
				if name, ok := rhs.X.(*ast.Ident); ok && declared[name.Name] == sinkValue {
					violations = append(violations, subprocessSinkViolation{
						position: fset.Position(lhs.Pos()),
						variable: name.Name,
					})
				}
			case *ast.Ident:
				if declared[rhs.Name] == sinkPointer {
					violations = append(violations, subprocessSinkViolation{
						position: fset.Position(lhs.Pos()),
						variable: rhs.Name,
					})
				}
			}
		}
		return true
	})
	return violations, true
}

// detectSubprocessSinks analyses each top-level function once. Returning false
// from the FuncDecl/FuncLit cases stops ast.Inspect descending, so a nested
// literal is never analysed a second time on its own and cannot double-count.
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
		return false
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
				name: "positive short-assign composite",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { buf := bytes.Buffer{}; cmd := exec.Command("helper"); cmd.Stderr = &buf; _ = cmd.Start() }
`,
				want: 1,
			},
			{
				name: "positive new pointer",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { buf := new(bytes.Buffer); cmd := exec.Command("helper"); cmd.Stderr = buf; _ = cmd.Start() }
`,
				want: 1,
			},
			{
				name: "positive inline address of literal",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { cmd := exec.Command("helper"); cmd.Stderr = &bytes.Buffer{}; _ = cmd.Start() }
`,
				want: 1,
			},
			{
				name: "positive start inside closure",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { var buf bytes.Buffer; cmd := exec.Command("helper"); cmd.Stderr = &buf; go func() { _ = cmd.Start() }() }
`,
				want: 1,
			},
			{
				name: "positive strings.Builder sink",
				src: `package fixture
import ("os/exec"; "strings")
func test() { var out strings.Builder; cmd := exec.Command("helper"); cmd.Stdout = &out; _ = cmd.Start() }
`,
				want: 1,
			},
			{
				name: "negative guarded sink by name",
				src: `package fixture
import "os/exec"
type syncBuf struct{}
func (*syncBuf) Write(p []byte) (int, error) { return len(p), nil }
func test() { var sink syncBuf; cmd := exec.Command("helper"); cmd.Stderr = &sink; _ = cmd.Start() }
`,
				want: 0,
			},
			{
				name: "negative guarded sink composite",
				src: `package fixture
import "os/exec"
type syncBuf struct{}
func (*syncBuf) Write(p []byte) (int, error) { return len(p), nil }
func test() { cmd := exec.Command("helper"); cmd.Stderr = &syncBuf{}; _ = cmd.Start() }
`,
				want: 0,
			},
			{
				name: "negative bare buffer with no Start",
				src: `package fixture
import ("bytes"; "os/exec")
func test() { var buf bytes.Buffer; cmd := exec.Command("helper"); cmd.Stderr = &buf; _ = cmd.Run() }
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
