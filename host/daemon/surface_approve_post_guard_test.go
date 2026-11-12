package daemon

// Row 201 A3 sprint 2 M5b (design §2.4, §4 T2, T3, T10, T12, §6 RT8-RT12): the POST's cross-origin
// guards, replay, the writer it never keeps, no session or cookie, and the writer-site inventory.

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// RT8: a POST from outside O, cross-site or not JSON is refused and records nothing; a preflight
// is 405 with no CORS header. A configured origin's page may POST through the loopback Host
// (behind a proxy); unconfigured, the same origin is refused naming --surface-origin.
func TestApprovePOSTCrossOrigin(t *testing.T) {
	a, row := newPostRig(t)
	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
	for _, tc := range []struct {
		hdr    []string
		status int
	}{
		{[]string{"Origin", ""}, 403}, {[]string{"Origin", "null"}, 403}, {[]string{"Origin", "http://localhost:7645"}, 403},
		{[]string{"Sec-Fetch-Site", "same-site"}, 403}, {[]string{"Sec-Fetch-Site", "cross-site"}, 403},
		{[]string{"Content-Type", "text/plain"}, 415}, {[]string{"Content-Type", "application/x-www-form-urlencoded"}, 415},
	} {
		if rec := a.post(row.RequestID, body, tc.hdr...); rec.Code != tc.status {
			t.Errorf("%v = %d, want %d", tc.hdr, rec.Code, tc.status)
		}
	}
	if r := requireRefusal(t, "unconfigured ts.net", a.post(row.RequestID, body, "Origin", "https://cfg.tail0.ts.net"), 403, "CrossOrigin"); !strings.Contains(r.Reason, "--surface-origin https://cfg.tail0.ts.net") {
		t.Errorf("reason %q does not name the fix", r.Reason)
	}
	if rec := a.serve("OPTIONS", "/surface/approve/"+row.RequestID, apLoopback, nil, "Origin", "http://localhost:7645", "Access-Control-Request-Method", "POST"); rec.Code != 405 ||
		rec.Header().Get("Access-Control-Allow-Origin")+rec.Header().Get("Access-Control-Allow-Methods") != "" {
		t.Fatalf("preflight = %d %v, want 405 and no CORS header", rec.Code, rec.Header())
	}
	if n, _ := a.recordedAndSessions(); n != 0 {
		t.Fatalf("%d recorded after refusals", n)
	}
	a.d.cfg.SurfaceOrigins = []string{"https://cfg.tail0.ts.net"}
	if rec := a.post(row.RequestID, body, "Origin", "https://cfg.tail0.ts.net"); rec.Code != 200 {
		t.Fatalf("a configured origin through the loopback Host = %d %s, want 200", rec.Code, rec.Body.String())
	}
}

// RT9: the same body twice while pending is recorded once (replaced); after the request is
// consumed the same body is refused as already used.
func TestApprovePOSTReplay(t *testing.T) {
	a, row := newPostRig(t)
	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
	for i := 0; i < 2; i++ {
		if rec := a.post(row.RequestID, body); rec.Code != 200 {
			t.Fatalf("POST %d = %d %s", i+1, rec.Code, rec.Body.String())
		}
	}
	if n, _ := a.recordedAndSessions(); n != 1 {
		t.Fatalf("%d recorded, want 1", n)
	}
	a.consume(row, apT0-1, strings.Repeat("cd", 32))
	if r := requireRefusal(t, "after the consume", a.post(row.RequestID, body), 409, "Conflict"); !strings.Contains(r.Reason, "already used") {
		t.Fatalf("reason %q", r.Reason)
	}
}

// RT10: a held writer is 503 Busy with Retry-After and nothing recorded, never a wait; released,
// the POST records; after every response the daemon holds no writer.
func TestApprovePOSTWriterBusy(t *testing.T) {
	a, row := newPostRig(t)
	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
	held, err := store.OpenCredentials(a.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	requireRefusal(t, "busy", a.post(row.RequestID, body), 503, "Busy") // RT4 R13-busy pins Retry-After
	_ = held.Close()
	if n, _ := a.recordedAndSessions(); n != 0 {
		t.Fatalf("%d recorded while busy", n)
	}
	if rec := a.post(row.RequestID, body); rec.Code != 200 {
		t.Fatalf("released: %d %s", rec.Code, rec.Body.String())
	}
	w, err := store.OpenCredentials(a.d.credentialsPath)
	if err != nil {
		t.Fatalf("the daemon kept the writer after the POST: %v", err)
	}
	_ = w.Close()
}

// RT11: neither a valid session Bearer nor a cookie approves: the handler answers its own
// refusals; isProtected is unchanged (POST /v1/commit only); the handler reads neither header.
func TestApprovePOSTNoSessionNoCookie(t *testing.T) {
	a, row := newPostRig(t)
	var tok string
	a.write(func(ctx context.Context, c *store.Credentials) (err error) {
		tok, _, _, err = authority.Mint(ctx, c, "w201b", grantForAuth(), 3600, time.Now().Unix(), nil)
		return err
	})
	auth := []string{"Authorization", "Bearer " + tok, "Cookie", "session=" + tok}
	requireRefusal(t, "session, garbage", a.post(row.RequestID, approveBody(row.Digest, []byte(`{"id":"x"}`)), auth...), 422, "NotApproved")
	requireRefusal(t, "session, no body", a.post(row.RequestID, nil, auth...), 400, "BadRequest")
	if n, s := a.recordedAndSessions(); n != 0 || s != 1 {
		t.Fatalf("%d recorded, %d sessions; want 0 and the one minted here", n, s)
	}
	if a.d.isProtected(httptest.NewRequest("POST", "/surface/approve/"+row.RequestID, nil)) || !a.d.isProtected(httptest.NewRequest("POST", "/v1/commit", nil)) {
		t.Fatal("isProtected changed")
	}
	src, err := os.ReadFile("surface_approve.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, read := range []string{`"Cookie"`, `"Authorization"`, `.Cookie(`, `.Cookies()`, "authority.FromContext"} {
		if strings.Contains(string(src), read) {
			t.Errorf("surface_approve.go reads %s", read)
		}
	}
}

// RT12: the daemon's production code opens the credentials writer at exactly two sites:
// ensureCredentialsStore's upgrade at serve start and the POST's record. Control: the read-only
// opens are found by the same walk.
func TestDaemonCredentialWriterSites(t *testing.T) {
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
				sel, ok := n.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if pkg, isIdent := sel.X.(*ast.Ident); isIdent && pkg.Name == "store" {
					switch sel.Sel.Name {
					case "OpenCredentials", "UpgradeCredentials":
						sites[fn.Name.Name] += sel.Sel.Name
					case "OpenCredentialsReadOnly":
						readOnly++
					}
				}
				return true
			})
		}
	}
	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 5 read-only opens", sites, readOnly, want)
	}
}
