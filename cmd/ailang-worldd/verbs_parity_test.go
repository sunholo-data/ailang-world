package main

// Row 243 V1 (design §6 VB4): host/verbs' three word sets are exactly run's switch arms but the
// navigation token world, each in one set, each with a help path — read with go/ast, never run().

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"testing"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

// switchArms is the case string literals of run's `switch verb := rest[0]`.
func switchArms(t *testing.T) (arms []string) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			return fn.Name.Name == "run"
		}
		sw, ok := n.(*ast.SwitchStmt)
		if !ok || sw.Init == nil {
			return true
		}
		if as, ok := sw.Init.(*ast.AssignStmt); !ok || as.Lhs[0].(*ast.Ident).Name != "verb" {
			return true
		}
		for _, s := range sw.Body.List {
			for _, e := range s.(*ast.CaseClause).List {
				if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					w, _ := strconv.Unquote(lit.Value)
					arms = append(arms, w)
				}
			}
		}
		return false
	})
	sort.Strings(arms)
	return arms
}

// classified: arms minus world equal the union of sets, each word in exactly one set.
func classified(arms []string, sets ...[]string) bool {
	seen, n := map[string]int{}, 0
	for _, set := range sets {
		for _, w := range set {
			seen[w]++
		}
	}
	for _, a := range arms {
		if a != "world" && seen[a] != 1 {
			return false
		} else if a != "world" {
			n++
		}
	}
	return n == len(seen)
}

func TestVerbTableMatchesTheCLISwitch(t *testing.T) {
	arms, op, nv, ot := switchArms(t), verbs.OperatorWords, verbs.NeverWords, verbs.OtherWords
	if len(arms) != 29 || !classified(arms, op, nv, ot) {
		t.Fatalf("run's switch arms %q (want 29) are not exactly operator ∪ never ∪ other, each once, world excluded", arms)
	}
	// Controls: an extra arm, a word the table lacks, a word in two sets each fail the check.
	if classified(append(arms, "nonsuch"), op, nv, ot) || classified(arms, op, nv, ot[1:]) || classified(arms, op, nv, append(ot, "why")) {
		t.Error("a control passed the parity check")
	}
	for _, w := range append(append(append([]string{}, op...), nv...), ot...) {
		if _, ok := genericVerbHelp([]string{w, "--help"}); !ownHelpVerbs[w] && !ok && w != "help" {
			t.Errorf("%s has no help path", w)
		}
	}
}
