package daemon

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type startupReference struct {
	File, Owner, Target, Kind string
	Line                      int
	Test                      bool
}

// Inventory references, including function values passed to startBounded and
// imported aliases (f := renamed.New is caught at the binding, before f()).
// Parse all Go files, including platform-specific files, without build tags.
// Local variables shadowing New are excluded by the parser's object binding.
func startupReferences(root string) ([]startupReference, int, error) {
	var refs []startupReference
	files := 0
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "vendor" || entry.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		files++
		set := token.NewFileSet()
		f, err := parser.ParseFile(set, path, nil, 0)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		local := filepath.ToSlash(filepath.Dir(rel)) == "host/daemon" && f.Name.Name == "daemon"
		aliases := map[string]bool{}
		for _, im := range f.Imports {
			p, _ := strconv.Unquote(im.Path.Value)
			if p != "github.com/sunholo-data/ailang-world/host/daemon" {
				continue
			}
			alias := "daemon"
			if im.Name != nil {
				alias = im.Name.Name
			}
			if alias == "." {
				return fmt.Errorf("dot daemon import in %s requires explicit caller review", rel)
			}
			aliases[alias] = true
		}
		for _, decl := range f.Decls {
			owner := "<package>"
			if fn, ok := decl.(*ast.FuncDecl); ok {
				owner = fn.Name.Name
			}
			calls := map[ast.Node]bool{}
			selectorNames := map[ast.Node]bool{}
			ast.Inspect(decl, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					calls[call.Fun] = true
				}
				if sel, ok := n.(*ast.SelectorExpr); ok {
					selectorNames[sel.Sel] = true
				}
				return true
			})
			ast.Inspect(decl, func(n ast.Node) bool {
				target := ""
				switch x := n.(type) {
				case *ast.SelectorExpr:
					id, ok := x.X.(*ast.Ident)
					if ok && aliases[id.Name] && id.Obj == nil && x.Sel.Name == "New" {
						target = "New"
					}
					if target == "" {
						return true
					}
				case *ast.Ident:
					if selectorNames[x] {
						return true
					}
					if !local {
						return true
					}
					if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name == x {
						return true
					}
					if x.Obj != nil && x.Obj.Kind != ast.Fun {
						return true
					}
					switch x.Name {
					case "New", "newDaemon", "startDaemon", "startBounded":
						target = x.Name
					}
				}
				if target != "" {
					kind := "function-value"
					if calls[n] {
						kind = "call"
					}
					refs = append(refs, startupReference{rel, owner, target, kind, set.Position(n.Pos()).Line, strings.HasSuffix(rel, "_test.go")})
				}
				return true
			})
		}
		return nil
	})
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].File != refs[j].File {
			return refs[i].File < refs[j].File
		}
		return refs[i].Line < refs[j].Line
	})
	return refs, files, err
}

func TestDaemonConstructorCallerInventory(t *testing.T) {
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("no caller path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(here), "../.."))
	refs, files, err := startupReferences(root)
	if err != nil {
		t.Fatal(err)
	}
	if files < 200 {
		t.Fatalf("inventory saw only %d Go files", files)
	}
	allowed := map[string]bool{
		"host/daemon/daemon.go|Run|startBounded":         true, // baseline construction
		"host/daemon/daemon.go|Run|New":                  true, // baseline Run -> startBounded(..., New)
		"host/daemon/daemon.go|startBounded|New":         true,
		"host/daemon/daemon.go|Run|startDaemon":          true,
		"host/daemon/daemon.go|startDaemon|startBounded": true,
		"host/daemon/daemon.go|startDaemon|newDaemon":    true,
		"host/daemon/daemon.go|New|newDaemon":            true,
	}
	newTests, newProd, prod := 0, 0, 0
	for _, r := range refs {
		if r.Target == "New" {
			if r.Test {
				newTests++
			} else {
				newProd++
			}
		}
		if !r.Test {
			prod++
			if !allowed[r.File+"|"+r.Owner+"|"+r.Target] {
				t.Errorf("unreviewed production startup caller: %+v", r)
			}
		}
		t.Logf("test=%t %s:%d owner=%s target=%s kind=%s", r.Test, r.File, r.Line, r.Owner, r.Target, r.Kind)
	}
	if newTests == 0 || prod == 0 {
		t.Fatal("caller inventory missed known positive New test or production path")
	}
	t.Logf("Go files=%d New test references=%d New production references=%d all production startup references=%d", files, newTests, newProd, prod)
}

func TestDaemonCallerInventoryAliasControl(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"cmd/control.go": `package main
import renamed "github.com/sunholo-data/ailang-world/host/daemon"
var indirect = renamed.New
func direct() { renamed.New(nil, renamed.Config{}) }
func nested() { sink(renamed.New).Method() }
`,
		"cmd/control_test.go": `package main
import d "github.com/sunholo-data/ailang-world/host/daemon"
var testConstructor = d.New
`,
		"host/daemon/control.go": `package daemon
func local() { f := New; f(nil, Config{}) }
func shadow() { New := func() {}; New() }
func field() { var obj struct{ New func() }; obj.New() }
`,
	}
	for name, data := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	refs, count, err := startupReferences(root)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 || len(refs) != 5 {
		t.Fatalf("alias/direct/local/shadow control: files=%d refs=%+v", count, refs)
	}
	got := map[string]int{}
	for _, r := range refs {
		got[fmt.Sprintf("%t|%s|%s", r.Test, r.Owner, r.Kind)]++
	}
	for _, key := range []string{"false|<package>|function-value", "false|direct|call", "false|nested|function-value", "true|<package>|function-value", "false|local|function-value"} {
		if got[key] != 1 {
			t.Fatalf("control missing %s: %+v", key, refs)
		}
	}
}
