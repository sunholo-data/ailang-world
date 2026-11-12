package daemon

// Row 201 A3 sprint 2 M4d (design §2.3, §6 RT1, RT7): the approval routes and their Host check.

import (
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// approvalRoutes is every approval route Handler registers this sprint (the card page and its
// script are sprint 3's).
var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}"}

// RT1: exactly the approval routes are registered, under /surface, each behind approvalRoute; no
// other method is answered on them; the card page and its script are not served yet.
func TestApprovalRoutesRegistered(t *testing.T) {
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
	var got []string
	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/(?:requests|approve)[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
		if m[2] == "" {
			t.Errorf("%s is not behind approvalRoute", m[1])
		}
		got = append(got, m[1])
	}
	sort.Strings(got)
	if want := append([]string(nil), approvalRoutes...); !reflect.DeepEqual(got, want) {
		t.Fatalf("approval routes %v, want %v", got, want)
	}
	a := newApprovalRig(t)
	if rec := a.serve("GET", "/surface/requests.json", apLoopback, nil); rec.Code != 200 {
		t.Fatalf("requests.json = %d", rec.Code)
	}
	requireRefusal(t, "unknown card", a.serve("GET", "/surface/requests/mr-0123456789abcdef", apLoopback, nil), 404, "NotFound")
	for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef"} {
		for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
			if rec := a.serve(method, target, apLoopback, nil, "Origin", "http://localhost:7644"); rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Errorf("%s %s = %d, want 405 with no CORS header", method, target, rec.Code)
			}
		}
	}
	for _, target := range []string{"/surface/approve.js", "/surface/approve", "/surface/requests", "/surface/requests/mr-0123456789abcdef/x"} {
		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404 (not served this sprint)", target, rec.Code)
		}
	}
}

// RT7: a Host outside O is 421 on every approval route; loopback at the arriving port, an enrolled
// passkey's host and a configured origin's host pass; with no local address there is no loopback
// member. Control: the existing routes still answer a foreign Host (row 199's posture).
func TestApprovalRoutesRefuseForeignHost(t *testing.T) {
	a := newApprovalRig(t)
	row := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
	targets := []string{"/surface/requests.json", "/surface/requests/" + row.RequestID}
	for _, target := range targets {
		for _, host := range []string{"rebind.example:7644", "localhost:7645", "localhost", "127.0.0.1:7644.evil", "studio.tail0.ts.net", "cfg.tail0.ts.net"} {
			requireRefusal(t, host+target, a.serve("GET", target, host, nil), 421, "MisdirectedRequest")
		}
		for _, host := range []string{"127.0.0.1:7644", "localhost:7644", "LOCALHOST:7644"} {
			if rec := a.serve("GET", target, host, nil); rec.Code != 200 {
				t.Errorf("GET %s with Host %s = %d, want 200", target, host, rec.Code)
			}
		}
	}
	if r := requireRefusal(t, "ts.net", a.serve("GET", targets[0], "studio.tail0.ts.net", nil), 421, "MisdirectedRequest"); !strings.Contains(r.Reason, "--surface-origin https://studio.tail0.ts.net") {
		t.Errorf("the unlisted ts.net Host's refusal %q does not name the fix", r.Reason)
	}
	a.enrolPasskey("b-phone", approvaltest.NewPasskey(approval.AlgES256, 62), "https://studio.tail0.ts.net")
	a.d.cfg.SurfaceOrigins = []string{"https://cfg.tail0.ts.net"}
	for _, target := range targets {
		for _, host := range []string{"studio.tail0.ts.net", "cfg.tail0.ts.net"} {
			if rec := a.serve("GET", target, host, nil); rec.Code != 200 {
				t.Errorf("GET %s with Host %s = %d, want 200", target, host, rec.Code)
			}
		}
	}
	rec := httptest.NewRecorder() // no local address: no loopback member
	req := httptest.NewRequest("GET", targets[0], nil)
	req.Host = apLoopback
	a.d.Handler().ServeHTTP(rec, req)
	requireRefusal(t, "no local address", rec, 421, "MisdirectedRequest")
	for _, target := range []string{"/surface/model.json", "/surface/passkey", "/surface/app.js"} {
		if rec := a.serve("GET", target, "rebind.example:7644", nil); rec.Code != 200 {
			t.Errorf("control: GET %s with a foreign Host = %d, want 200 (existing routes unchanged)", target, rec.Code)
		}
	}
}
