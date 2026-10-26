package daemon

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// TestSurfaceAppRoutes: row 188 N2b. GET /surface is the app shell (the
// same bytes whatever the index holds), its assets are the exact embedded
// bytes, and every other name under /surface/ is a 404.
func TestSurfaceAppRoutes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operator-index.json")
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	h := d.Handler()
	request := func(method, target string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, target, nil))
		return w
	}
	shell := func(target string, status int) {
		t.Helper()
		w := request("GET", target)
		if w.Code != status {
			t.Fatalf("%s status = %d, want %d", target, w.Code, status)
		}
		for k, v := range map[string]string{"Content-Security-Policy": surfaceCSP, "Content-Type": "text/html; charset=utf-8",
			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := w.Header().Get(k); got != v {
				t.Errorf("%s: %s = %q, want %q", target, k, got, v)
			}
		}
		if !bytes.Equal(w.Body.Bytes(), surfacepage.AppShell) {
			t.Errorf("%s is not the embedded shell", target)
		}
	}
	model := func(status int, key string) {
		t.Helper()
		w := request("GET", "/surface/model.json")
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || w.Code != status || body[key] != true || body["reason"] == "" {
			t.Errorf("model.json = %d %s, want %d with %s", w.Code, w.Body.String(), status, key)
		}
	}

	t.Run("shell", func(t *testing.T) {
		body := string(surfacepage.AppShell)
		scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(body, -1)
		want := []string{surfacepage.ScriptElement[:len(surfacepage.ScriptElement)-len("</script>")], `<script src="/surface/app.js" defer>`}
		if len(scripts) != 2 || scripts[0] != want[0] || scripts[1] != want[1] {
			t.Errorf("scripts = %q, want %q", scripts, want)
		}
		if strings.Count(body, "></script>") != 2 {
			t.Error("a script element has an inline body")
		}
		if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=`).FindString(body); m != "" {
			t.Errorf("handler attribute %q", m)
		}
		for _, banned := range []string{"<img", "<link", "SPIKE_DATA", "data.js", "spike-data", "http://", "https://"} {
			if strings.Contains(strings.ReplaceAll(body, `xmlns="http://www.w3.org/2000/svg"`, ""), banned) {
				t.Errorf("shell has %q", banned)
			}
		}
		if !strings.Contains(body, "<style>\n"+string(surfacepage.AppStyle)+"</style>") {
			t.Error("app.css is not inlined whole")
		}
		for _, marker := range []string{`<html lang="en" data-app>`, `role="img" aria-label="AILANG"`, `id="stage" data-view="timeline"`, `<span class="brand-name">World</span>`} {
			if !strings.Contains(body, marker) {
				t.Errorf("shell lacks %q", marker)
			}
		}
	})

	t.Run("index states", func(t *testing.T) {
		shell("/surface", 200) // not built
		model(200, "empty")
		d.cfg.SurfaceIndex = ""
		shell("/surface", 200) // not configured
		model(200, "empty")
		d.cfg.SurfaceIndex = path
		rig, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, rig, 0o600); err != nil {
			t.Fatal(err)
		}
		shell("/surface", 200)
		shell("/surface?brane=row-189&at=f36428abd9a72e5093b61cba14e716c37019f1ef", 200)
		shell("/surface?brane=row-189&at=missing", 404)
		shell("/surface?brane=unknown&at=f36428abd9a72e5093b61cba14e716c37019f1ef", 404)
		if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
			t.Fatal(err)
		}
		shell("/surface", 503) // refused
		model(503, "refused")
		if w := request("GET", "/surface"); strings.Contains(w.Body.String(), path) {
			t.Error("the shell leaked the configured path")
		}
	})

	t.Run("assets", func(t *testing.T) {
		for _, tc := range []struct {
			target, contentType string
			body                []byte
		}{
			{"/surface/app.js", "text/javascript; charset=utf-8", surfacepage.AppScript},
			{"/surface/app.css", "text/css; charset=utf-8", surfacepage.AppStyle},
			{"/surface/ailang-logo.svg", "image/svg+xml", surfacepage.AppLogo},
		} {
			w := request("GET", tc.target)
			if w.Code != 200 || len(tc.body) == 0 || !bytes.Equal(w.Body.Bytes(), tc.body) {
				t.Errorf("%s: %d, or not the embedded bytes", tc.target, w.Code)
			}
			for k, v := range map[string]string{"Content-Type": tc.contentType, "Cache-Control": "no-cache",
				"X-Content-Type-Options": "nosniff", "Content-Security-Policy": "default-src 'none'"} {
				if got := w.Header().Get(k); got != v {
					t.Errorf("%s: %s = %q, want %q", tc.target, k, got, v)
				}
			}
			if w := request("POST", tc.target); w.Code != 405 {
				t.Errorf("POST %s = %d", tc.target, w.Code)
			}
		}
		// Every font app.css names is served by N2a's route.
		fonts := regexp.MustCompile(`url\('(/surface/fonts/[^']+)'\)`).FindAllStringSubmatch(string(surfacepage.AppStyle), -1)
		if len(fonts) != 3 {
			t.Fatalf("app.css names %d fonts, want 3", len(fonts))
		}
		for _, f := range fonts {
			if w := request("GET", f[1]); w.Code != 200 || w.Header().Get("Content-Type") != "font/ttf" {
				t.Errorf("%s = %d", f[1], w.Code)
			}
		}
	})

	t.Run("404", func(t *testing.T) {
		for _, target := range []string{"/surface/", "/surface/nope", "/surface/app.jsx", "/surface/index.html", "/surface/data.js",
			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css"} {
			if w := request("GET", target); w.Code != 404 {
				t.Errorf("%s = %d, want 404", target, w.Code)
			}
		}
	})
}
