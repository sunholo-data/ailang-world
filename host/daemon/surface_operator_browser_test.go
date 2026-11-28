package daemon

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestOperatorCardAssets(t *testing.T) {
	f := newOperatorHTTP(t, time.Minute)
	before := f.evidence()
	s := f.request(opRequest(1))
	seeded := f.evidence()
	if seeded.actions != before.actions+1 || seeded.operations != before.operations+1 {
		t.Fatalf("pending row control: %+v to %+v", before, seeded)
	}
	card := f.do("GET", s.Href, nil, nil)
	html := card.Body.Bytes()
	if bytes.Contains(html, []byte(`rel="stylesheet"`)) || bytes.Count(html, []byte("<style>")) != 1 ||
		!bytes.Contains(html, []byte("<style>"+string(surfacepage.OperatorStyle)+"</style>")) {
		t.Fatal("card must inline the embedded stylesheet exactly once")
	}
	const script = "/surface/operator-card.js"
	if !bytes.Contains(html, []byte(script)) {
		t.Fatalf("card omitted %s", script)
	}
	r := f.do("GET", script, nil, nil)
	if r.Code != 200 || !bytes.Equal(r.Body.Bytes(), surfacepage.OperatorScript) || r.Header().Get("Content-Type") != "text/javascript; charset=utf-8" ||
		r.Header().Get("X-Content-Type-Options") != "nosniff" || r.Header().Get("Content-Security-Policy") != "default-src 'none'" {
		t.Fatalf("asset %s: %d %v", script, r.Code, r.Header())
	}
	if r := f.do("POST", script, nil, nil); r.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST %s: %d", script, r.Code)
	}
	if r := f.do("GET", script, nil, func(r *http.Request) { r.Host = "other.invalid" }); r.Code != http.StatusMisdirectedRequest {
		t.Fatalf("bad Host %s: %d", script, r.Code)
	}
	if after := f.evidence(); after != seeded {
		t.Fatalf("asset reads changed action or operation counts: %+v to %+v", seeded, after)
	}
}
