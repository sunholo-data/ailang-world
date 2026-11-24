package main

// Row 243 V3 sprint 6 (S5 §1.4 F8): the CLI's production code opens the credentials writer at exactly
// these sites, as host/daemon's TestDaemonCredentialWriterSites pins the daemon's: a session's open
// (mints, approver verbs), launch's record, finish's revoke and S5's verb-request put. Control: the
// read-only opens are found by the same walk.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCLICredentialWriterSites(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	sites, readOnly := map[string]string{}, 0
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), f, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			ast.Inspect(fn, func(n ast.Node) bool {
				if sel, ok := n.(*ast.SelectorExpr); ok {
					if pkg, isIdent := sel.X.(*ast.Ident); isIdent && pkg.Name == "store" {
						switch sel.Sel.Name {
						case "OpenCredentials", "UpgradeCredentials":
							sites[fn.Name.Name] += sel.Sel.Name
						case "OpenCredentialsReadOnly":
							readOnly++
						}
					}
				}
				return true
			})
		}
	}
	want := map[string]string{"runnerWrite": "OpenCredentials", "openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials"}
	if !reflect.DeepEqual(sites, want) || readOnly < 20 {
		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 20 read-only opens", sites, readOnly, want)
	}
	t.Logf("read-only opens: %d", readOnly)
}
