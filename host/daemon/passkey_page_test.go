package daemon

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// PG1 (row 201 A2 M7, design §2.14, §6): GET /surface/passkey and /surface/passkey.js serve the
// embedded page and script, byte-equal to host/surfacepage/passkey/, under surfaceCSP and the
// surface headers, with the world store closed (they read nothing); POST is refused.
func TestPasskeyPageRoutes(t *testing.T) {
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
	if err != nil {
		t.Fatal(err)
	}
	h := d.Handler()
	if err := d.Close(); err != nil { // no store read: the routes answer with the store closed
		t.Fatal(err)
	}
	for _, tc := range []struct{ target, file, contentType string }{
		{"/surface/passkey", "index.html", "text/html; charset=utf-8"},
		{"/surface/passkey.js", "passkey.js", "text/javascript; charset=utf-8"},
	} {
		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "passkey", tc.file))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.target, nil))
		if w.Code != 200 || len(want) == 0 || !bytes.Equal(w.Body.Bytes(), want) {
			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/passkey/%s", tc.target, w.Code, tc.file)
		}
		for k, v := range map[string]string{"Content-Security-Policy": surfaceCSP, "Content-Type": tc.contentType,
			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := w.Header().Get(k); got != v {
				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
			}
		}
		w = httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("POST", tc.target, nil))
		if w.Code != 405 {
			t.Errorf("POST %s = %d, want 405", tc.target, w.Code)
		}
	}

	// The page loads exactly its own script, with no inline script or handler and no other URL.
	page := string(surfacepage.PasskeyPage)
	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || scripts[0] != `<script src="/surface/passkey.js">` || !strings.Contains(page, `<script src="/surface/passkey.js"></script>`) {
		t.Errorf("scripts %q", scripts)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
		t.Errorf("the page has %q", m)
	}
	// Labelled instrument-health controls (not load-bearing; XB1 and XB2 are): the script makes no
	// request and writes no markup.
	for _, banned := range []string{"fetch(", "XMLHttpRequest", "innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "sendBeacon", "WebSocket"} {
		if bytes.Contains(surfacepage.PasskeyScript, []byte(banned)) {
			t.Errorf("passkey.js uses %s", banned)
		}
	}
}
