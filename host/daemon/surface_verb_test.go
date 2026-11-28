package daemon

// Row 243 V1 (design §2.2, §6 VR1, VR4, VS1): GET /surface/verb.json.

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

// verbRoutes is every verb route Handler registers.
var verbRoutes = []string{"GET /surface/verb.json"}

func verbRig(t *testing.T) (approvalRig, *reportSeed) {
	t.Helper()
	s := newReportSeed(t)
	s.call("mcp", "w1", "1", "ailang-read", read(50, map[string]any{"ok": true, "content": "x"}), nil)
	s.call("mcp", "w1", "2", `<img src=x onerror=alert(1)>`, nil, map[string]any{"ok": false, "refused": `<b>hostile</b> "reason"`})
	return approvalRig{t: t, d: s.d}, s
}

func verbGet(a approvalRig, q string) (*verbAnswerT, string) {
	a.t.Helper()
	rec := a.serve("GET", "/surface/verb.json?q="+url.QueryEscape(q), apLoopback, nil)
	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
		a.t.Fatalf("%q: %d %s %v", q, rec.Code, rec.Body, err)
	}
	var v verbAnswerT
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil || v.V != 1 || v.Schema != VerbSchema || v.Q != q || v.Argv == nil {
		a.t.Fatalf("%q: not a v1 answer: %s", q, rec.Body)
	}
	return &v, rec.Body.String()
}

// verbAnswerT decodes an answer with every member optional, so a missing one is seen.
type verbAnswerT struct {
	V, Exit                                            int
	Schema, Q, Verb, Class, Outcome, Reason, CLI, Href string
	Argv                                               []string
	Text, JSON, Stderr                                 *string
}

// VR1: one route, behind approvalRoute: a foreign Host is 421; GET only; nothing else under it.
func TestSurfaceVerbRouteRegistered(t *testing.T) {
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
	var got []string
	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/verb[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
		if m[2] == "" {
			t.Errorf("%s is not behind approvalRoute", m[1])
		}
		got = append(got, m[1])
	}
	sort.Strings(got)
	if !reflect.DeepEqual(got, verbRoutes) {
		t.Fatalf("verb routes %v, want %v", got, verbRoutes)
	}
	a := newApprovalRig(t)
	for _, host := range []string{"rebind.example:7644", "localhost:7645", "studio.tail0.ts.net"} {
		requireRefusal(t, host, a.serve("GET", "/surface/verb.json?q=why", host, nil), 421, "MisdirectedRequest")
	}
	for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
		if rec := a.serve(method, "/surface/verb.json?q=why", apLoopback, nil, "Origin", "http://localhost:7644"); rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("%s = %d, want 405 with no CORS header", method, rec.Code)
		}
	}
	for _, target := range []string{"/surface/verb", "/surface/verb.json/x", "/surface/verb.jsonx"} {
		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404", target, rec.Code)
		}
	}
}

// VR4: q is required; the line bound is an operator refusal, checked on the raw value; a store
// failure or an expired read deadline answers a refusal with the /v1 route's text, never the error.
func TestSurfaceVerbBoundsAndErrors(t *testing.T) {
	a, s := verbRig(t)
	requireRefusal(t, "no q", a.serve("GET", "/surface/verb.json", apLoopback, nil), 400, "BadRequest")
	requireRefusal(t, "other key", a.serve("GET", "/surface/verb.json?line=why", apLoopback, nil), 400, "BadRequest")
	if v, _ := verbGet(a, ""); v.Outcome != "address" {
		t.Errorf("empty q = %+v", *v)
	}
	over := "why " + strings.Repeat("x", verbs.MaxLine-3)
	if v, _ := verbGet(a, over); v.Outcome != "refused" || v.Reason != verbs.ErrTooLong.Error() || len(v.Argv) != 0 || v.Class != "refused" {
		t.Errorf("1,025 bytes = %+v", *v)
	}
	if v, _ := verbGet(a, over[:verbs.MaxLine]); v.Outcome != "refused" || v.Reason == verbs.ErrTooLong.Error() || len(v.Argv) != 2 {
		t.Errorf("1,024 bytes = %+v", *v)
	}
	if v, _ := verbGet(a, `why "a`); v.Outcome != "refused" || v.Reason != verbs.ErrUnclosed.Error() {
		t.Errorf("unclosed = %+v", *v)
	}
	var log bytes.Buffer
	s.d.errLog = &log
	s.d.readDeadline = expiredReadDeadline
	if v, raw := verbGet(a, "episode report w1"); v.Outcome != "refused" || v.Reason != "read deadline (-1ns) exceeded" || v.Text != nil {
		t.Errorf("expired deadline = %s", raw)
	}
	s.d.readDeadline = readDeadline
	if err := s.d.store.Close(); err != nil {
		t.Fatal(err)
	}
	v, raw := verbGet(a, "episode report w1")
	if v.Outcome != "refused" || v.Reason != internalErrorMessage || strings.Contains(raw, "closed") {
		t.Errorf("store failure = %s", raw)
	}
	if !strings.Contains(log.String(), "internal error: GET /surface/verb.json") {
		t.Errorf("the failure is not logged for the operator: %q", log.String())
	}
	// The /v1 route still answers its own classes from the extracted builder.
	assertErrorClass(t, requestRecorder(t, s.d, "GET", "/v1/episodes/w1/report", nil), 500, "Internal")
}

// VS1: host/verbs stays pure. The bar delegates operator review to one
// read-only adapter, and approve's id pattern is host/store's.
func TestSurfaceVerbIsReadOnly(t *testing.T) {
	importsOf := func(path string) (map[string]bool, *ast.File) {
		t.Helper()
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		out := map[string]bool{}
		for _, im := range f.Imports {
			p, _ := strconv.Unquote(im.Path.Value)
			out[p] = true
		}
		return out, f
	}
	files, err := filepath.Glob(filepath.Join("..", "verbs", "*.go"))
	if err != nil || len(files) == 0 {
		t.Fatalf("host/verbs: %v %v", files, err)
	}
	pure := map[string]bool{"errors": true, "flag": true, "fmt": true, "io": true, "regexp": true, "strings": true, "time": true, "unicode/utf8": true}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		ims, _ := importsOf(f)
		for p := range ims {
			if !pure[p] {
				t.Errorf("%s imports %s", f, p)
			}
		}
	}
	ims, f := importsOf("surface_verb.go")
	for p := range ims {
		if !map[string]bool{"bytes": true, "encoding/json": true, "fmt": true, "net/http": true, "strings": true,
			"github.com/sunholo-data/ailang-world/host/verbs": true}[p] {
			t.Errorf("surface_verb.go imports %s", p)
		}
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
	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true, "verbWhy": true, "verbProfileList": true, "verbPresetList": true, "verbAppList": true, "operatorVerb": true}; !reflect.DeepEqual(calls, want) {
		t.Errorf("surface_verb.go reaches d.%v, want exactly %v", calls, want)
	}
	for _, id := range []string{"mr-0123456789abcdef", "mr-0123456789ABCDEF", "mr-0123456789abcde", "mr-0123456789abcdef0", "xr-0123456789abcdef", "mr-012345678 abcdef"} {
		if href, _ := verbs.ApproveTarget([]string{id}); (href != "") != store.MintRequestIDPattern.MatchString(id) {
			t.Errorf("approve %q: href %q, but store's pattern says %v", id, href, store.MintRequestIDPattern.MatchString(id))
		}
	}
}
