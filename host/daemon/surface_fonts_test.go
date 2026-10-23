package daemon

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSurfaceFontRoute(t *testing.T) {
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	h := d.Handler()
	for _, name := range []string{"Montserrat-Variable.ttf", "IBMPlexMono-Regular.ttf", "IBMPlexMono-Medium.ttf", "OFL-Montserrat.txt", "OFL-IBMPlexMono.txt"} {
		t.Run(name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join("..", "surfacepage", "fonts", name))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", "/surface/fonts/"+name, nil))
			if w.Code != 200 || !bytes.Equal(w.Body.Bytes(), want) {
				t.Fatalf("status %d or bytes differ", w.Code)
			}
			ct := "font/ttf"
			if strings.HasSuffix(name, ".txt") {
				ct = "text/plain; charset=utf-8"
			}
			for k, v := range map[string]string{"Content-Type": ct, "Cache-Control": "public, max-age=31536000, immutable", "X-Content-Type-Options": "nosniff"} {
				if got := w.Header().Get(k); got != v {
					t.Errorf("%s = %q, want %q", k, got, v)
				}
			}
		})
	}
	for _, suffix := range []string{"", "..", "../Montserrat-Variable.ttf", "%2e%2e", "%2e%2e%2fMontserrat-Variable.ttf", "Montserrat-Variable.ttf%2Fextra", "%2FMontserrat-Variable.ttf", "montserrat-variable.ttf", "missing.ttf", "OFL-Montserrat.TXT", "Montserrat-Variable.ttf/extra", "./Montserrat-Variable.ttf", "/Montserrat-Variable.ttf"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/surface/fonts/"+suffix, nil))
		if w.Code != 404 || !strings.Contains(w.Body.String(), "surface font not found") {
			t.Errorf("%q: status %d, body %q", suffix, w.Code, w.Body.String())
		}
	}
	for _, suffix := range []string{"?x=1", "?", "?%zz", "?brane=x&at=y"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/surface/fonts/Montserrat-Variable.ttf"+suffix, nil))
		if w.Code != 400 || strings.Contains(w.Body.String(), "<script") {
			t.Errorf("query %q: %d", suffix, w.Code)
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/surface/fonts/Montserrat-Variable.ttf", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST = %d", w.Code)
	}
}

func TestSurfaceSecurityHeaders(t *testing.T) {
	const wantCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
	path := filepath.Join(t.TempDir(), "index.json")
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	h := d.Handler()
	check := func(target string, status int) {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", target, nil))
		if w.Code != status {
			t.Fatalf("%s status = %d, want %d", target, w.Code, status)
		}
		for k, v := range map[string]string{"Content-Security-Policy": wantCSP, "Content-Type": "text/html; charset=utf-8", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := w.Header().Get(k); got != v {
				t.Errorf("%s: %s = %q, want %q", target, k, got, v)
			}
		}
	}
	check("/surface", 200) // not built
	d.cfg.SurfaceIndex = ""
	check("/surface", 200) // unconfigured
	d.cfg.SurfaceIndex = path
	rig, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rig, 0600); err != nil {
		t.Fatal(err)
	}
	check("/surface", 200)
	check("/surface?brane=row-189&at=f36428abd9a72e5093b61cba14e716c37019f1ef", 200)
	check("/surface?x=1", 400)
	check("/surface?brane=missing&at=missing", 404)
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	check("/surface", 503)
	check("/surface?brane=row-189&at=missing", 503)
}
