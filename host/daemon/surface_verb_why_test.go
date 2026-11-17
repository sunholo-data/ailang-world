package daemon

// Row 243 V2b (design §2.4, §2.5 e, §6 VR5, VR7, VR8, VS2): why's in-process getter and the A2 cap.

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/why"
)

func methodName(h http.HandlerFunc) string {
	n := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	return strings.TrimSuffix(n[strings.LastIndex(n, ".")+1:], "-fm")
}

// VR5: a path outside the seven is an error (never a 404, never another handler); a 404 is an answer;
// the body cap is 8 MiB exactly; a 503 is a timeout and another 5xx a failure.
func TestWhyGetterUnroutedAndOversize(t *testing.T) {
	d := newHandlerDaemon(t)
	var log bytes.Buffer
	d.errLog = &log
	ctx := boundedTestContext(t)
	r := httptest.NewRequest("GET", "/surface/verb.json?q=why", nil)
	for _, p := range []string{"/v1/session/self", "/v1/nowhere", "/v1/registry/x", "/v1/objects/by-semantic-id/x", "/v1/episodes/w1/report"} {
		var tally whyTally
		if st, body, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, p); err == nil || err.Error() != "unrouted in-process read: "+p || st != 0 || body != nil || !tally.failed || tally.timeout {
			t.Errorf("%s = %d %q %v %+v", p, st, body, err, tally)
		}
	}
	if !strings.Contains(log.String(), "internal error: GET /surface/verb.json: unrouted in-process read: /v1/session/self") {
		t.Errorf("not logged: %q", log.String())
	}
	var tally whyTally
	if _, _, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, "/v1/log/%zz"); err == nil || !tally.failed {
		t.Errorf("an unparseable path = %v %+v", err, tally)
	}
	tally = whyTally{}
	if st, _, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, "/v1/objects/sha256:"+strings.Repeat("0", 64)+"?payload=true"); st != 404 || err != nil || tally != (whyTally{}) {
		t.Errorf("absent object = %d %v %+v", st, err, tally)
	}
	stub := http.NewServeMux()
	for _, c := range []struct {
		n      int
		status int
	}{{8388608, 200}, {8388609, 0}} {
		n := c.n
		stub.HandleFunc(fmt.Sprintf("GET /big/%d", n), func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(bytes.Repeat([]byte("x"), n)) })
		var tally whyTally
		st, body, err := d.whyGetter(r, stub, &tally)(ctx, fmt.Sprintf("/big/%d", n))
		if c.status == 200 && (st != 200 || len(body) != n || err != nil || tally.failed) {
			t.Errorf("%d bytes = %d %d %v", n, st, len(body), err)
		}
		if c.status == 0 && (err == nil || err.Error() != "response from /big/8388609 exceeds 8388608 bytes" || !tally.failed) {
			t.Errorf("%d bytes = %d %v %+v", n, st, err, tally)
		}
	}
	for status, want := range map[int]whyTally{503: {timeout: true}, 500: {failed: true}, 502: {failed: true}, 404: {}, 400: {}} {
		code := status
		stub.HandleFunc(fmt.Sprintf("GET /s/%d", code), func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(code) })
		var tally whyTally
		if st, _, err := d.whyGetter(r, stub, &tally)(ctx, fmt.Sprintf("/s/%d", code)); st != code || err != nil || tally != want {
			t.Errorf("status %d = %d %v %+v", code, st, err, tally)
		}
	}
}

// VR7 (unit): the A2 cap on text+json+stderr together, at 1 MiB exactly; members cleared; only a result
// is capped. The JSON-encoded body of an answer at the cap is larger than the cap and bounded by 6×.
func TestVerbAnswerCap(t *testing.T) {
	third := strings.Repeat("y", 1048576/2)
	for _, c := range []struct {
		a    verbAnswer
		over bool
	}{
		{verbAnswer{Outcome: "result", Text: strings.Repeat("x", 1048576)}, false},
		{verbAnswer{Outcome: "result", Text: strings.Repeat("x", 1048577)}, true},
		{verbAnswer{Outcome: "result", JSON: strings.Repeat("x", 1048577)}, true},
		{verbAnswer{Outcome: "result", Stderr: strings.Repeat("x", 1048576) + "\n", Exit: 1}, true},
		{verbAnswer{Outcome: "result", Text: third, JSON: third, Stderr: "z", Exit: 3}, true},
		{verbAnswer{Outcome: "result", Text: third, JSON: third}, false},
		{verbAnswer{Outcome: "refused", Reason: "r", Text: strings.Repeat("x", 1048577)}, false},
	} {
		a := c.a
		a.CLI = "ailang-worldd why 1"
		n := len(a.Text) + len(a.JSON) + len(a.Stderr)
		capVerbAnswer(&a)
		if over := a.Outcome == "refused" && a.Reason == fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: ailang-worldd why 1", n); over != c.over {
			t.Errorf("%d bytes: %s %q", n, a.Outcome, a.Reason)
		}
		if c.over && (a.Text != "" || a.JSON != "" || a.Stderr != "" || a.Exit != 0) {
			t.Errorf("%d bytes: members left on the refusal", n)
		}
		if !c.over && c.a.Outcome == "result" && (a.Outcome != "result" || a.Text != c.a.Text) {
			t.Errorf("%d bytes refused", n)
		}
	}
	w := httptest.NewRecorder()
	writeApprovalJSON(w, 200, verbAnswer{Outcome: "result", Text: strings.Repeat("\x01", 1048576)})
	if n := w.Body.Len(); n <= 1048576 || n > 6*1048576+1024 {
		t.Errorf("the encoded body is %d bytes", n)
	}
}

// VR8: no handler the private mux dispatches reads middleware state (a context value), so the
// in-process path cannot silently read nil where HTTP would read something; control: approvalOrigins's
// LocalAddrContextKey read is found by the same instrument.
func TestWhyHandlersReadNoMiddlewareState(t *testing.T) {
	fset := token.NewFileSet()
	files, _ := filepath.Glob("*.go")
	bodies := map[string]*ast.FuncDecl{}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		af, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range af.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv != nil {
				bodies[fd.Name.Name] = fd
			}
		}
	}
	readsState := func(fd *ast.FuncDecl) bool {
		found := false
		ast.Inspect(fd, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				found = found || x.Sel.Name == "Value"
			case *ast.Ident:
				found = found || x.Name == "activityRequestKey"
			}
			return true
		})
		return found
	}
	if fd := bodies["approvalOrigins"]; fd == nil || !readsState(fd) {
		t.Fatal("control: approvalOrigins's context read is not found")
	}
	d := newHandlerDaemon(t)
	for p, h := range d.whyRoutes() {
		if fd := bodies[methodName(h)]; fd == nil || readsState(fd) {
			t.Errorf("%s: %s reads middleware state (or is not found)", p, methodName(h))
		}
	}
}

// VS2: surface_verb_why.go's imports and d.* reach are exact sets: the seven read handlers (VR4/VR6
// prove which), the read deadline, the error log and its own three helpers; each pattern is bound to the
// method Handler registers for it. No store, credentials store, process or write.
func TestSurfaceVerbWhyIsReadOnly(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "surface_verb_why.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ims := map[string]bool{}
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		ims[p] = true
	}
	if want := map[string]bool{"bytes": true, "context": true, "errors": true, "fmt": true, "net/http": true,
		"github.com/sunholo-data/ailang-world/host/verbs": true, "github.com/sunholo-data/ailang-world/host/why": true}; !reflect.DeepEqual(ims, want) {
		t.Errorf("surface_verb_why.go imports %v", ims)
	}
	calls := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "d" {
				calls[sel.Sel.Name] = true
			}
		}
		return true
	})
	want := map[string]bool{"handleHealth": true, "handleHead": true, "handleWorld": true, "handleObject": true, "handleLogEntry": true,
		"handleLogRange": true, "handleReceipt": true, "whyRoutes": true, "whyMux": true, "whyGetter": true, "readCtx": true,
		"readDeadline": true, "writeInternalErrorLog": true}
	if !reflect.DeepEqual(calls, want) {
		t.Errorf("surface_verb_why.go reaches d.%v, want exactly %v", calls, want)
	}
	d := newHandlerDaemon(t)
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	for p, h := range d.whyRoutes() {
		if !regexp.MustCompile(`mux\.HandleFunc\("` + regexp.QuoteMeta(p) + `", d\.` + methodName(h) + `\)`).Match(src) {
			t.Errorf("%s is bound to %s, which Handler does not register for it", p, methodName(h))
		}
	}
}

// VR1 (unit, design §2.5 d and f): the walk's answer is shown only in one of the CLI's three shapes —
// stdout with exit 0 or 3 (text, or json with --json), stderr with exit 1; anything else is (d).
func TestWhyResult(t *testing.T) {
	const d = "the walk answered something the app cannot show; at a terminal: ailang-worldd why 1"
	for _, c := range []struct {
		ans    why.Answer
		asJSON bool
		want   verbAnswer
	}{
		{why.Answer{Stdout: "o\n"}, false, verbAnswer{Outcome: "result", Text: "o\n"}},
		{why.Answer{Stdout: "o\n", Exit: 3}, false, verbAnswer{Outcome: "result", Text: "o\n", Exit: 3}},
		{why.Answer{Stdout: "{}\n", Exit: 3}, true, verbAnswer{Outcome: "result", JSON: "{}\n", Exit: 3}},
		{why.Answer{Stderr: "e\n", Exit: 1}, true, verbAnswer{Outcome: "result", Stderr: "e\n", Exit: 1}},
		{why.Answer{}, false, verbAnswer{Outcome: "refused", Reason: d}},
		{why.Answer{Stdout: "o", Stderr: "e", Exit: 1}, false, verbAnswer{Outcome: "refused", Reason: d}},
		{why.Answer{Stdout: "o", Exit: 1}, false, verbAnswer{Outcome: "refused", Reason: d}},
		{why.Answer{Stdout: "o", Exit: 2}, false, verbAnswer{Outcome: "refused", Reason: d}},
		{why.Answer{Stderr: "e", Exit: 3}, false, verbAnswer{Outcome: "refused", Reason: d}},
		{why.Answer{Stderr: "e"}, false, verbAnswer{Outcome: "refused", Reason: d}},
	} {
		a := verbAnswer{Outcome: "refused", CLI: "ailang-worldd why 1"}
		whyResult(&a, c.ans, c.asJSON)
		c.want.CLI = a.CLI
		if !reflect.DeepEqual(a, c.want) {
			t.Errorf("%+v json=%v: %+v, want %+v", c.ans, c.asJSON, a, c.want)
		}
	}
}
