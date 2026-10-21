package daemon

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// TestSurfaceScriptRoute: GET /surface/surface.js serves the embedded
// script with /workbench/live.js's headers, outside the frozen /v1 table.
func TestSurfaceScriptRoute(t *testing.T) {
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	h := d.Handler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surface/surface.js", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	for k, v := range map[string]string{"Content-Type": "text/javascript; charset=utf-8", "Content-Security-Policy": "default-src 'none'",
		"X-Content-Type-Options": "nosniff", "Cache-Control": "no-store"} {
		if got := w.Header().Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	if len(surfacepage.Script) == 0 || w.Body.String() != string(surfacepage.Script) {
		t.Fatal("served script is not the embedded asset")
	}
	for _, seam := range []string{"laneStates", "parseAddress", "keyAction", "shouldHandle", "buildModel", "loadLook", "saveLook", "resolveLook"} {
		if !strings.Contains(w.Body.String(), seam) {
			t.Errorf("served script lacks seam %s", seam)
		}
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/surface/surface.js", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST = %d", w.Code)
	}
}

func TestSurfacePage(t *testing.T) {
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
	check := func(w *httptest.ResponseRecorder, status int, want ...string) string {
		t.Helper()
		body := w.Body.String()
		if w.Code != status {
			t.Fatalf("status = %d, want %d: %s", w.Code, status, body)
		}
		for k, v := range map[string]string{"Content-Security-Policy": workbenchCSP, "X-Content-Type-Options": "nosniff",
			"Referrer-Policy": "no-referrer", "Content-Type": "text/html; charset=utf-8", "Cache-Control": "no-store"} {
			if got := w.Header().Get(k); got != v {
				t.Errorf("%s = %q, want %q", k, got, v)
			}
		}
		for _, s := range want {
			if !strings.Contains(body, s) {
				t.Errorf("body lacks %q", s)
			}
		}
		// A 400 is the workbench's plain error page; every surface page
		// carries the one script.
		scripts := 1
		if status == http.StatusBadRequest {
			scripts = 0
		}
		inertSurfacePage(t, body, scripts)
		return body
	}
	header := `<span class="name">World</span>`

	t.Run("not built", func(t *testing.T) {
		check(request("GET", "/surface"), 200, header, "the surface index has not been built yet", "Nothing needs you")
	})
	t.Run("not configured", func(t *testing.T) {
		d.cfg.SurfaceIndex = ""
		defer func() { d.cfg.SurfaceIndex = path }()
		body := check(request("GET", "/surface"), 200, header, "no surface index is configured", "World serving, no index shown")
		if strings.Contains(body, "World healthy") {
			t.Error("an absent index reads as healthy")
		}
	})

	rig, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	creds, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	for _, row := range []store.SessionRow{
		{CredentialID: strings.Repeat("a", 64), EpisodeID: "w188e", GrantsJSON: "[]", CreatedAt: now - 360, ExpiresAt: now + 3600},
		{CredentialID: strings.Repeat("b", 64), EpisodeID: "w-old", GrantsJSON: "[]", CreatedAt: now - 7200, ExpiresAt: now - 1},
	} {
		if err := creds.MintSession(boundedTestContext(t), row); err != nil {
			t.Fatal(err)
		}
	}
	if err := creds.Close(); err != nil {
		t.Fatal(err)
	}

	t.Run("fixture index", func(t *testing.T) {
		w189d := "f36428abd9a72e5093b61cba14e716c37019f1ef"
		body := check(request("GET", "/surface"), 200,
			header, "World healthy · checked ",
			`<g class="lane" data-brane="shared" data-glyph="circle">`,
			`<g class="lane" data-brane="row-189" data-glyph="triangle"><title>w-world-onboarding · row 189</title>`,
			`<a class="state" data-state="`+w189d+`" data-kind="step" data-brane="row-189" href="/surface?brane=row-189&amp;at=`+w189d+`">`,
			`<a class="state" data-state="90035557b2012895551867c0ebe2f36b853aac09" data-kind="shared" data-brane="shared" href="/surface?brane=shared&amp;at=90035557b2012895551867c0ebe2f36b853aac09">`,
			`<path class="link" data-from="0ee1945ffba795812f06dce95998b646adfac343" data-to="`+w189d+`"`,
			`<path class="link" data-from="`+w189d+`" data-to="90035557b2012895551867c0ebe2f36b853aac09"`,
			`<text class="badge" data-state="`+w189d+`"`,
			`<text class="tick" x="288" y="14">17 Oct</text>`,
			"1 live session", "running w188e · started 6 min ago", "Nothing needs you")
		if got := strings.Count(body, `<a class="state" `); got != 48 {
			t.Errorf("%d states, want 48", got)
		}
		if strings.Contains(body, "w-old") {
			t.Error("an expired session is listed as live")
		}
	})

	t.Run("stack", func(t *testing.T) {
		check(request("GET", "/surface?brane=row-189&at=f36428abd9a72e5093b61cba14e716c37019f1ef"), 200, header, `class="front"`, `data-col="inputs"`)
		for _, target := range []string{
			"/surface?brane=unknown&at=f36428abd9a72e5093b61cba14e716c37019f1ef",
			"/surface?brane=row-189&at=missing",
			"/surface?brane=shared&at=f36428abd9a72e5093b61cba14e716c37019f1ef",
		} {
			check(request("GET", target), 404, header, `href="/surface"`)
		}
		if w := request("POST", "/surface?brane=shared&at=x"); w.Code != 405 {
			t.Fatalf("stack POST = %d", w.Code)
		}
	})

	t.Run("query", func(t *testing.T) {
		for _, target := range []string{"/surface?x=1", "/surface?brane=row-189", "/surface?at=x", "/surface?brane=x&at=y&other=z", "/surface?brane=x&brane=y&at=z", "/surface?brane=x&at=y&at=z", "/surface?brane=%zz&at=x", "/surface?world"} {
			w := request("GET", target)
			check(w, 400, "unsupported surface query parameter")
			if strings.Contains(w.Body.String(), `class="lane"`) {
				t.Errorf("%s rendered the timeline", target)
			}
		}
	})

	t.Run("refused", func(t *testing.T) {
		for name, data := range map[string]string{
			"path in validation error": strings.Replace(string(rig), surface.Schema, path, 1),
			"undecodable":              "{",
		} {
			if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
				t.Fatal(err)
			}
			w := request("GET", "/surface")
			body := check(w, 503, header, "the configured surface index: ", "Nothing needs you")
			if strings.Contains(body, path) || strings.Contains(body, filepath.Dir(path)) {
				t.Errorf("%s: configured path leaked", name)
			}
			if strings.Contains(body, `class="lane"`) || strings.Contains(body, "World healthy") {
				t.Errorf("%s: a refused index rendered as healthy lanes", name)
			}
		}
		d.cfg.SurfaceIndex = filepath.Dir(path)
		defer func() { d.cfg.SurfaceIndex = path }()
		if body := check(request("GET", "/surface"), 503, header); strings.Contains(body, d.cfg.SurfaceIndex) {
			t.Error("read error leaked the configured path")
		}
	})

	t.Run("POST", func(t *testing.T) {
		if w := request(http.MethodPost, "/surface"); w.Code != 405 {
			t.Fatalf("POST = %d", w.Code)
		}
	})
}

// inertSurfacePage is the page rule (M2, M4a): scripts script elements, each
// exactly /surface/surface.js with no inline body; no on*= handler; and
// every href local.
func inertSurfacePage(t *testing.T, body string, scripts int) {
	t.Helper()
	if n := strings.Count(strings.ToLower(body), "<script"); n != scripts || strings.Count(body, surfacepage.ScriptElement) != scripts {
		t.Errorf("page has %d <script elements, want %d of %s", n, scripts, surfacepage.ScriptElement)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=`).FindString(body); m != "" {
		t.Errorf("page has a handler attribute %q", m)
	}
	for _, h := range regexp.MustCompile(`(?i)\bhref\s*=\s*["']?([^"'\s>]*)`).FindAllStringSubmatch(body, -1) {
		if !(strings.HasPrefix(h[1], "/") || strings.HasPrefix(h[1], "#")) || strings.HasPrefix(h[1], "//") {
			t.Errorf("href %q is not local", h[1])
		}
	}
}
