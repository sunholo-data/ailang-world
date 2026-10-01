package daemon

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestMCPMountMethodSource(t *testing.T) {
	tree, err := parser.ParseFile(token.NewFileSet(), "daemon.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	mcp, a2a := 0, 0
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Handler" {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 2 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "HandleFunc" {
				return true
			}
			literal, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}
			pattern, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatal(err)
			}
			handler, ok := call.Args[1].(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if handler.Sel.Name == "MCP" {
				mcp++
				if pattern != "POST /mcp/" {
					t.Errorf("MCP mount=%q want exactly POST /mcp/", pattern)
				}
			}
			if pattern == "POST /a2a/" && handler.Sel.Name == "A2A" {
				a2a++
			}
			return true
		})
	}
	if mcp != 1 || a2a != 1 {
		t.Fatalf("executable MCP/A2A registrations=%d/%d", mcp, a2a)
	}
}
func TestMCPGetRefused(t *testing.T) {
	d := newHandlerDaemon(t)
	server := httptest.NewServer(d.Handler())
	defer server.Close()
	client := &http.Client{Timeout: DefaultClientTimeout}
	r, err := client.Get(server.URL + "/mcp/")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	if r.StatusCode != 405 || r.Header.Get("Allow") != "POST" {
		t.Fatalf("GET MCP=%d Allow=%q", r.StatusCode, r.Header.Get("Allow"))
	}
}
func TestMCPMountedBearerDenial(t *testing.T) {
	d := newHandlerDaemon(t)
	r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(`{}`))
	if d.isProtected(r) {
		t.Fatal("MCP entered REST session middleware")
	}
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, r)
	if w.Code != 401 || strings.Contains(w.Body.String(), `"error"`) {
		t.Fatalf("MCP denial=%d %s", w.Code, w.Body)
	}
}
