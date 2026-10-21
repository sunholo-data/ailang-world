package daemon

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestSurfaceIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operator-index.json")
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	h := d.Handler()
	request := func(method string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "/surface/index.json", nil))
		return w
	}
	check := func(w *httptest.ResponseRecorder, status int, body string) {
		t.Helper()
		if w.Code != status || w.Body.String() != body {
			t.Fatalf("response = %d %s, want %d %s", w.Code, w.Body.String(), status, body)
		}
		for k, v := range map[string]string{"Content-Type": "application/json", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff"} {
			if w.Header().Get(k) != v {
				t.Errorf("%s = %q, want %q", k, w.Header().Get(k), v)
			}
		}
	}
	empty := func(reason string) string {
		return `{"schema":"world/surface-index/v1","empty":true,"reason":"` + reason + `"}`
	}
	t.Run("absent", func(t *testing.T) {
		check(request("GET"), 200, empty("the surface index has not been built yet: run ailang-worldd surface-index"))
	})
	t.Run("not configured", func(t *testing.T) {
		d.cfg.SurfaceIndex = ""
		defer func() { d.cfg.SurfaceIndex = path }()
		check(request("GET"), 200, empty("no surface index is configured: start serve with --surface-index <file>"))
	})
	write := func(data string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	valid := `{"schema":"world/surface-index/v1","built_at":"2026-10-21T00:00:00Z","source":{},"branes":[],"states":[]}`
	t.Run("valid and rebuilt", func(t *testing.T) {
		write(valid)
		check(request("GET"), 200, valid)
		rebuilt := valid + "\n  "
		write(rebuilt)
		check(request("GET"), 200, rebuilt)
	})
	var ix surface.Index
	if err := json.Unmarshal([]byte(valid), &ix); err != nil {
		t.Fatal(err)
	}
	state := surface.State{ID: strings.Repeat("a", 40), Kind: surface.KindStep, At: ix.BuiltAt}
	ix.States = []surface.State{state, state}
	invalid, err := json.Marshal(ix)
	if err != nil {
		t.Fatal(err)
	}
	validationErr := surface.Validate(ix)
	if validationErr == nil {
		t.Fatal("duplicate fixture is valid")
	}
	for _, tc := range []struct {
		name, data, reason string
	}{
		{"oversize", valid + strings.Repeat(" ", surface.MaxBytes+1-len(valid)), "4 MiB"},
		{"undecodable", "{", "unexpected end"},
		{"duplicate state", string(invalid), validationErr.Error()},
		{"path in validation error", strings.Replace(valid, surface.Schema, path, 1), "schema"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			write(tc.data)
			w := request("GET")
			if w.Code != 503 || !strings.Contains(w.Body.String(), tc.reason) || !strings.Contains(w.Body.String(), "SurfaceIndexRefused") {
				t.Fatalf("response = %d %s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), path) {
				t.Fatalf("configured path leaked: %s", w.Body.String())
			}
		})
	}
	t.Run("read error", func(t *testing.T) {
		d.cfg.SurfaceIndex = filepath.Dir(path)
		defer func() { d.cfg.SurfaceIndex = path }()
		w := request("GET")
		if w.Code != 503 || strings.Contains(w.Body.String(), d.cfg.SurfaceIndex) {
			t.Fatalf("read error = %d %s", w.Code, w.Body.String())
		}
	})
	t.Run("POST", func(t *testing.T) {
		if w := request(http.MethodPost); w.Code != 405 {
			t.Fatalf("POST = %d", w.Code)
		}
	})
}

func TestSurfaceIndexStartupRefusal(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(root, "index.json"), filepath.Join(link, "index.json")} {
		t.Run(path, func(t *testing.T) {
			d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, WorkspaceRoot: root, SurfaceIndex: path})
			if d != nil {
				d.Close()
			}
			var se *StartupError
			if !errors.As(err, &se) || se.Stage != StageConfig || !strings.Contains(err.Error(), "surface index") || !strings.Contains(err.Error(), "--workspace-root") {
				t.Fatalf("startup = %v, want surface index StageConfig refusal", err)
			}
		})
	}
}
