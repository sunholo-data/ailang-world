package daemon

import (
	"bytes"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMCPMountMethodSource(t *testing.T) {
	tree, err := parser.ParseFile(token.NewFileSet(), "daemon.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	mcp, a2a := 0, 0
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Handler" || fn.Recv == nil || len(fn.Recv.List) != 1 {
			continue
		}
		receiver, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		typ, ok := receiver.X.(*ast.Ident)
		if !ok || typ.Name != "Daemon" {
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
			mux, ok := sel.X.(*ast.Ident)
			if !ok || mux.Name != "mux" {
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
			projection, ok := handler.X.(*ast.SelectorExpr)
			if !ok || projection.Sel.Name != "projection" {
				return true
			}
			owner, ok := projection.X.(*ast.Ident)
			if !ok || owner.Name != "d" {
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

func TestMCPDaemonBudgetWiring(t *testing.T) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, "daemon.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"CredentialWait": "credentialBudget", "MaxWait": "readDeadline", "InvokeWait": "invokeDeadline", "CallbackTimeout": "invokeDeadline", "MaxCallbacks": "procbound.MaxOutstanding", "WriteWait": "writeTimeout"}
	seen := 0
	ast.Inspect(tree, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "New" {
			return true
		}
		owner, ok := sel.X.(*ast.Ident)
		if !ok || owner.Name != "projection" {
			return true
		}
		cfg, ok := call.Args[0].(*ast.CompositeLit)
		if !ok {
			t.Fatal("projection config is not literal")
		}
		seen++
		found := map[string]string{}
		for _, elt := range cfg.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			var b bytes.Buffer
			if err := format.Node(&b, fset, kv.Value); err != nil {
				t.Fatal(err)
			}
			found[key.Name] = b.String()
		}
		for key, value := range want {
			if found[key] != value {
				t.Errorf("production %s=%s expected%s", key, found[key], value)
			}
		}
		return true
	})
	if seen != 1 {
		t.Fatalf("production projection configs=%d", seen)
	}
	if credentialBudget != 3*time.Second || readDeadline != 10*time.Second || invokeDeadline != 20*time.Second || writeTimeout != 30*time.Second || procbound.MaxOutstanding != 8 {
		t.Fatal("frozen production bounds changed")
	}
}
