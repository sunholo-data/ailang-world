package daemon

// Row 201 A3 sprint 3 M6b (design §2.7, §6 PG2): the request card page and its script.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// PG2: GET /surface/approve/{id} and /surface/approve.js serve the embedded bytes, byte-equal to
// host/surfacepage/approve/, under the unchanged surfaceCSP and the surface headers; a malformed id
// is 400; the A2 page's bytes are unchanged. The page loads exactly its own script.
func TestApprovePageRoutes(t *testing.T) {
	const csp = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
	if surfaceCSP != csp {
		t.Fatalf("surfaceCSP changed: %q", surfaceCSP)
	}
	a := newApprovalRig(t)
	for _, tc := range []struct{ target, file, contentType string }{
		{"/surface/approve/mr-0123456789abcdef", "index.html", "text/html; charset=utf-8"},
		{"/surface/approve/mr-00000000000000b1", "index.html", "text/html; charset=utf-8"},
		{"/surface/approve.js", "approve.js", "text/javascript; charset=utf-8"},
	} {
		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "approve", tc.file))
		if err != nil {
			t.Fatal(err)
		}
		rec := a.serve("GET", tc.target, apLoopback, nil)
		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/approve/%s", tc.target, rec.Code, tc.file)
		}
		for k, v := range map[string]string{"Content-Security-Policy": csp, "Content-Type": tc.contentType,
			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := rec.Header().Get(k); got != v {
				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
			}
		}
		if rec.Header().Get("Set-Cookie") != "" || rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("GET %s sets a cookie or a CORS header", tc.target)
		}
	}
	for _, id := range []string{"mr-0123", "mr-0123456789ABCDEF", "x"} {
		requireRefusal(t, id, a.serve("GET", "/surface/approve/"+id, apLoopback, nil), 400, "BadRequest")
	}
	for file, sum := range map[string]string{"index.html": "1695ebc2ae20f303dbfa3b617b3bec4c76035447ec290deba91e695fec9cc06d",
		"passkey.js": "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f"} {
		b, err := os.ReadFile(filepath.Join("..", "surfacepage", "passkey", file))
		if h := sha256.Sum256(b); err != nil || hex.EncodeToString(h[:]) != sum {
			t.Errorf("the A2 page's %s changed (%v)", file, err)
		}
	}
	page := string(surfacepage.ApprovePage)
	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/approve.js"></script>`) {
		t.Errorf("scripts %q", scripts)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
		t.Errorf("the page has %q", m)
	}
	// Labelled instrument-health controls (not load-bearing; XA1 and XA2 are): the script writes no
	// markup, evaluates no string and reads no cookie.
	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
		if bytes.Contains(surfacepage.ApproveScript, []byte(banned)) {
			t.Errorf("approve.js uses %s", banned)
		}
	}
}
