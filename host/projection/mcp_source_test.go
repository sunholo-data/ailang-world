package projection

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"
)

func TestMCPWireOwnershipSource(t *testing.T) {
	sources := projectionSources(t)
	calls, imports, wrappers := 0, 0, 0
	for file, src := range sources {
		tree, err := parser.ParseFile(token.NewFileSet(), file, src, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, im := range tree.Imports {
			path, _ := strconv.Unquote(im.Path.Value)
			if path == "github.com/sunholo-data/ailang/serveapi/protocol/mcphttp" {
				imports++
			}
		}
		ast.Inspect(tree, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "mcphttp" && sel.Sel.Name == "NewHandler" {
				calls++
			}
			return true
		})
		for _, decl := range tree.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "MCP" {
				continue
			}
			wrappers++
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					allowed := false
					switch f := call.Fun.(type) {
					case *ast.Ident:
						allowed = f.Name == "cancel"
					case *ast.SelectorExpr:
						allowed = f.Sel.Name == "WithTimeout" || f.Sel.Name == "Context" || f.Sel.Name == "ServeHTTP" || f.Sel.Name == "WithContext"
					}
					if !allowed {
						t.Errorf("aggregate wrapper has non-context/dispatch call")
					}
				}
				return true
			})
		}
		if strings.HasSuffix(file, "mcp.go") {
			for _, bad := range []string{"WriteMCPEnvelope", "json.NewEncoder", "event: message", "ResponseController", "SetWriteDeadline", "WriteHeader"} {
				if strings.Contains(src, bad) {
					t.Errorf("local wire/deadline operation: %s", bad)
				}
			}
		}
	}
	if imports != 1 || calls != 1 || wrappers != 1 {
		t.Fatalf("upstream import/call/wrapper census=%d/%d/%d", imports, calls, wrappers)
	}
}
