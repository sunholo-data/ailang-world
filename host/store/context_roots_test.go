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
	"cmd/ailang-worldd/episode_launch.go|launchEpisode|Background": 3, // row 217 M2: HTTP preflight (8 s), agent (wall cap + margin), append (3 s)

	"cmd/ailang-worldd/cli.go|get|Background":                                1,
	"cmd/ailang-worldd/cli.go|execute|Background":                            1,
	"cmd/ailang-worldd/cli.go|executeWithAuth|Background":                    1,
	"cmd/ailang-worldd/cli.go|budgetContext|Background":                      1, // row 138: tools/call/why/provenance whole-command budget; no Store call
	"cmd/ailang-worldd/signal_context.go|serveSignalContext|Background":      1,
	"cmd/ailang-worldd/session.go|runSessionMint|Background":                 1,
	"cmd/ailang-worldd/session.go|runSessionRevoke|Background":               1,
	"cmd/ailang-worldd/doctor.go|doctorStore|Background":                     1, // row 138 M4: doctor's read-only store reads (3 s) through OpenReadOnly
	"cmd/ailang-worldd/session_new.go|provisionEpisode|Background":           1, // row 138 M5 / row 189 M3: worktree add + 3 s mint + rollback (75 s), shared by session new and episode start
	"cmd/ailang-worldd/session_new.go|checkBindingFits|Background":           1, // row 189 M2: read-only ListSessions of one episode before the fence (3 s)
	"cmd/ailang-worldd/preset.go|runPreset|Background":                       1, // row 217: read-only preset validation and git inspection (8 s)
	"cmd/ailang-worldd/episode_preset.go|runPresetEpisodeStart|Background":   1, // row 217: pre-fence validation and git inspection (8 s)
	"cmd/ailang-worldd/episode_finish.go|runEpisodeFinishWith|Background":    1, // row 218 M1: finish's reads, pinned checks (each capped), commit (2 h)
	"cmd/ailang-worldd/session_new.go|runSessionList|Background":             1, // row 138 M5: read-only ListSessions (3 s)
	"cmd/ailang-worldd/setup.go|runSetup|Background":                         1, // row 138 M3: setup's whole-command budget (12 min); no Store call
	"cmd/ailang-worldd/profile_init.go|runProfileInit|Background":            1, // row 189 M4: profile init's toolchain probes (75 s); no Store call
	"cmd/ailang-worldd/surface_index.go|runSurfaceIndex|Background":          1, // row 188 M0a: surface-index's git reads; no Store call
	"cmd/world-publish/main.go|runApprove|Background":                        1,
	"cmd/world-publish/main.go|runPublish|Background":                        1,
	"cmd/world-publish/main.go|runReconcile|Background":                      1,
	"cmd/world-publish/reconcile_probe.go|reconcileReadOnlyProbe|Background": 1,
	"cmd/world-publish/transitions.go|runTransitions|Background":             1,
	"cmd/world-publish/epoch.go|epochNominateContext|Background":             1, // row 139 M2: epoch nominate's read and write budgets (10 s each)
	"host/archive/archive.go|probeVersion|Background":                        1,
	// Row 134 §4.3: post-effect writes are detached from the caller's
	// cancellation and each feeds context.WithTimeout directly (4 s / 2 s
	// budgets), so an executed effect is never stranded unrecorded.
	"host/broker/broker.go|postDispatchContext|WithoutCancel": 1,
	"host/coordinator/effectful.go|detached|WithoutCancel":    1,
	// Row 184 D2: the after-call worktree snapshot is detached from the
	// caller's cancellation (the toolchain call already ran) and feeds
	// context.WithTimeout directly (worktreeAfterBudget, 60 s); no Store call.
	"host/broker/worktree_delta.go|finish|WithoutCancel": 1,
	// Row 156 M6b: the after-command protected check is detached from the
	// caller's cancellation, re-capped by the call's ORIGINAL deadline
	// (context.WithDeadline), and checkProtected bounds it by
	// protectedCheckTimeout (3 751 ms, design V24); Git only, no Store call.
	"host/broker/handlers_exec.go|Execute|WithoutCancel": 1,
	// Row 219 M3 round 2: re-acceptance's single-flight probe is detached from
	// the first caller's cancellation (waiters share its result) and feeds
	// context.WithTimeout directly (the node version budget plus a 120 s probe
	// ceiling); round 3 also cancels it through ExecToolState.Close on daemon
	// shutdown. No Store call.
	"host/broker/exec_tools.go|resolve|WithoutCancel": 1,
	"host/capsule/capsule.go|Run|Background":          1,
	"host/daemon/daemon.go|drain|Background":          1,
	// Row 134 M4b: an episode's tool handler is built inside the coordinator's
	// BinderFor, which carries no ctx; its one summary subprocess is bounded
	// by workspaceHandlerBudget (3 s) fed directly from this root, and it
	// makes no Store call.
	"host/daemon/workspace.go|episodeHandler|Background": 1,
	// Row 189 M2: the exec profile binding is read inside the same ctx-less
	// BinderFor, through a fresh read-only credentials handle bounded by the
	// credential budget (3 s, B4) fed directly from this root.
	"host/daemon/binding.go|episodeBinding|Background": 1,
	// Row 156 M6a: NewExecHandler (built inside the same ctx-less BinderFor)
	// captures the protected baseline under protectedBaselineTimeout (4 848 ms,
	// design V24) fed directly from this root; Git subprocesses only, no Store
	// call.
	"host/broker/handlers_exec.go|NewExecHandler|Background":                  1,
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
		// The marketing-capture verifier CLI (M7b, row 150): plan-mandated production
		// tooling outside host/ and cmd/; the exact two-pass PNG gate for the six-shot pass.
		"scripts/marketingcheck/main.go": false,
		"scripts/marketingcheck/png.go":  false,
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
