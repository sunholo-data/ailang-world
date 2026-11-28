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
	for _, asset := range []struct {
		path, media string
		body        []byte
	}{
		{"/surface/operator-card.js", "text/javascript; charset=utf-8", surfacepage.OperatorScript},
		{"/surface/operator-card.css", "text/css; charset=utf-8", surfacepage.OperatorStyle},
	} {
		if !bytes.Contains(card.Body.Bytes(), []byte(asset.path)) {
			t.Fatalf("card omitted %s", asset.path)
		}
		r := f.do("GET", asset.path, nil, nil)
		if r.Code != 200 || !bytes.Equal(r.Body.Bytes(), asset.body) || r.Header().Get("Content-Type") != asset.media ||
			r.Header().Get("X-Content-Type-Options") != "nosniff" || r.Header().Get("Content-Security-Policy") != "default-src 'none'" {
			t.Fatalf("asset %s: %d %v", asset.path, r.Code, r.Header())
		}
		if r := f.do("POST", asset.path, nil, nil); r.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST %s: %d", asset.path, r.Code)
		}
		if r := f.do("GET", asset.path, nil, func(r *http.Request) { r.Host = "other.invalid" }); r.Code != http.StatusMisdirectedRequest {
			t.Fatalf("bad Host %s: %d", asset.path, r.Code)
		}
	}
	if after := f.evidence(); after != seeded {
		t.Fatalf("asset reads changed action or operation counts: %+v to %+v", seeded, after)
	}
}
