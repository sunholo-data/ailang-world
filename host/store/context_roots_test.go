package store

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// Every surviving production root is named, including the explicit CLI mint
// root. These are compatibility exceptions, NOT evidence of bounded store waits.
var contextRootPins = map[string]int{
	"cmd/ailang-worldd/cli.go|get|Background":                                 1,
	"cmd/ailang-worldd/cli.go|execute|Background":                             1,
	"cmd/ailang-worldd/cli.go|executeWithAuth|Background":                     1,
	"cmd/ailang-worldd/signal_context.go|serveSignalContext|Background":       1,
	"cmd/ailang-worldd/session.go|runSessionMint|Background":                  1,
	"cmd/ailang-worldd/session.go|runSessionRevoke|Background":                1,
	"cmd/world-publish/main.go|runApprove|Background":                         1,
	"cmd/world-publish/main.go|runPublish|Background":                         1,
	"cmd/world-publish/main.go|runReconcile|Background":                       1,
	"cmd/world-publish/reconcile_probe.go|reconcileReadOnlyProbe|Background":  1,
	"cmd/world-publish/transitions.go|runTransitions|Background":              1,
	"host/archive/archive.go|probeVersion|Background":                         1,
	"host/capsule/capsule.go|Run|Background":                                  1,
	"host/daemon/daemon.go|drain|Background":                                  1,
	"host/replay/replay.go|runPinnedTransition|Background":                    1,
	"host/store/lookup_index.go|provisionLookupIndex|Background":              1,
	"host/store/reference_index.go|provisionReferenceIndexes|Background":      1,
	"host/store/reference_index.go|verifyReadOnlyReferenceIndexes|Background": 1,
}

// These files create contexts for process or client lifetime and make no Store
// call. Their roots are still counted above.
var nonStoreReachingRoots = map[string]bool{
	"cmd/ailang-worldd/cli.go":             true,
	"cmd/ailang-worldd/signal_context.go":  true,
	"cmd/world-publish/reconcile_probe.go": true,
	"host/capsule/capsule.go":              true,
}

func isContextDeadlineConstructor(expr ast.Expr, alias string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok || (sel.Sel.Name != "WithTimeout" && sel.Sel.Name != "WithDeadline") {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == alias
}

func TestProductionContextRoots(t *testing.T) {
	root := repoRootFromCaller(t)
	got := map[string]int{}
	seen := map[string]bool{}
	for _, top := range []string{"host", "cmd"} {
		err := filepath.Walk(filepath.Join(root, top), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			seen[rel] = true
			f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				return err
			}
			alias := ""
			for _, im := range f.Imports {
				p, _ := strconv.Unquote(im.Path.Value)
				if p == "context" {
					alias = "context"
					if im.Name != nil {
						alias = im.Name.Name
					}
					if alias == "." {
						t.Errorf("dot context import in %s", rel)
					}
				}
			}
			if alias == "" {
				return nil
			}
			for _, d := range f.Decls {
				owner := "<package>"
				if fn, ok := d.(*ast.FuncDecl); ok {
					owner = fn.Name.Name
				}
				var visit func(ast.Node, ast.Node)
				visit = func(n ast.Node, parent ast.Node) {
					if n == nil {
						return
					}
					if call, ok := n.(*ast.CallExpr); ok && !nonStoreReachingRoots[rel] {
						if root, ok := call.Fun.(*ast.SelectorExpr); ok && root.Sel.Name == "Background" {
							if id, ok := root.X.(*ast.Ident); ok && id.Name == alias {
								outer, ok := parent.(*ast.CallExpr)
								if !ok || len(outer.Args) == 0 || outer.Args[0] != call || !isContextDeadlineConstructor(outer.Fun, alias) {
									t.Errorf("%s|%s: Background must directly feed WithTimeout/WithDeadline", rel, owner)
								}
							}
						}
					}
					x, ok := n.(*ast.SelectorExpr)
					if ok {
						if id, ok := x.X.(*ast.Ident); ok && id.Name == alias {
							switch x.Sel.Name {
							case "Background", "TODO", "WithoutCancel":
								got[rel+"|"+owner+"|"+x.Sel.Name]++
							}
						}
					}
					ast.Inspect(n, func(child ast.Node) bool {
						if child != n {
							visit(child, n)
							return false
						}
						return true
					})
				}
				visit(d, nil)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(seen) < 58 {
		t.Fatalf("root scan saw %d production files, want >=58", len(seen))
	}
	for _, p := range []string{"host/broker/approve.go", "host/broker/publish_op.go", "host/registry/registry.go", "host/replay/replay.go", "host/daemon/daemon.go", "cmd/world-publish/main.go", "host/store/store.go"} {
		if !seen[p] {
			t.Fatalf("root scan missed %s", p)
		}
	}
	totals := map[string]int{}
	for key, count := range got {
		parts := strings.Split(key, "|")
		totals[parts[len(parts)-1]] += count
	}
	t.Logf("root references: Background=%d TODO=%d WithoutCancel=%d", totals["Background"], totals["TODO"], totals["WithoutCancel"])
	if !reflect.DeepEqual(got, contextRootPins) {
		t.Fatalf("production context roots changed: got %v; want %v", got, contextRootPins)
	}
	t.Logf("root census: %d named roots, %d production files", len(got), len(seen))
}

// Keep the host/cmd census honest when new source trees are introduced. Walking
// also checks untracked files, so the guard fires before a new tree is committed.
func TestProductionGoSurface(t *testing.T) {
	root := repoRootFromCaller(t)
	allowed := map[string]bool{
		"design_docs/verification/w-race-gate-blindspot/racecontrol/main.go": false,
		"design_docs/verification/w-race-gate-blindspot/repro/main.go":       false,
	}
	scanned, reproducers := 0, 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" || info.Name() == ".snap" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "host/") || strings.HasPrefix(rel, "cmd/") {
			scanned++
		} else if _, ok := allowed[rel]; ok {
			allowed[rel] = true
			reproducers++
		} else {
			t.Errorf("non-test Go file outside context-root scan and explicit reproducer allow-list: %s", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for path, seen := range allowed {
		if !seen {
			t.Errorf("allow-listed reproducer missing: %s", path)
		}
	}
	t.Logf("Go surface: host/cmd=%d allow-listed reproducers=%d", scanned, reproducers)
}
