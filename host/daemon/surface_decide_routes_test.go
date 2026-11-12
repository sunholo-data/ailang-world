package daemon

// Row 201 A4 sprint 1 (design §2.8, §6 RD1): the decision routes, under /surface, behind
// approvalRoute; no other method; no CORS; /v1 is untouched (TestFrozenV1RouteTableMatchesMux).

import (
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// RD1: exactly the two decision routes are registered, each behind approvalRoute; a foreign Host is
// 421 on both; other methods are 405 with no CORS header; no other name under them is served.
func TestDecisionRoutesRegistered(t *testing.T) {
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
	var got []string
	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/deci(?:de|sions)[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
		if m[2] == "" {
			t.Errorf("%s is not behind approvalRoute", m[1])
		}
		got = append(got, m[1])
	}
	sort.Strings(got)
	if want := []string{"GET /surface/decisions/{ref}", "POST /surface/decide/{ref}"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("decision routes %v, want %v", got, want)
	}
	r := newApprovalRig(t)
	ref := "sha256:" + strings.Repeat("ab", 32)
	for _, tc := range []struct{ method, target string }{{"GET", "/surface/decisions/" + ref}, {"POST", "/surface/decide/" + ref}} {
		requireRefusal(t, "foreign "+tc.target, r.serve(tc.method, tc.target, "rebind.example:7644", nil, "Origin", "http://127.0.0.1:7644"), 421, "MisdirectedRequest")
	}
	requireRefusal(t, "unknown subject", r.serve("GET", "/surface/decisions/"+ref, apLoopback, nil), 404, "NotFound")
	for _, tc := range []struct {
		target  string
		methods []string
	}{
		{"/surface/decisions/" + ref, []string{"POST", "PUT", "DELETE", "OPTIONS"}},
		{"/surface/decide/" + ref, []string{"GET", "PUT", "DELETE", "OPTIONS"}},
	} {
		for _, method := range tc.methods {
			rec := r.serve(method, tc.target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
			if rec.Code != 405 {
				t.Errorf("%s %s = %d, want 405", method, tc.target, rec.Code)
			}
			for k := range rec.Header() {
				if strings.HasPrefix(strings.ToLower(k), "access-control-") {
					t.Errorf("%s %s wrote %s", method, tc.target, k)
				}
			}
		}
	}
	for _, target := range []string{"/surface/decisions", "/surface/decisions/", "/surface/decisions/" + ref + "/x", "/surface/decide", "/surface/decide/" + ref + "/x"} {
		if rec := r.serve("GET", target, apLoopback, nil); rec.Code != 404 && rec.Code != 405 {
			t.Errorf("GET %s = %d, want 404 or 405", target, rec.Code)
		}
	}
}
