package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Row 274 M1: GET /surface/apps.json and the `app list` bar verb read the local app
// catalog, validate each manifest, and never show a file path.
func TestSurfaceAppsCatalog(t *testing.T) {
	example, err := os.ReadFile(filepath.Join("..", "..", "docs", "examples", "apps", "docparse-convert.json"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, body := range map[string]string{
		"docparse-convert.json": string(example),
		"broken.json":           "{",
		"wide.json":             strings.Replace(strings.Replace(string(example), `"docparse-convert"`, `"wide"`, 1), `"Workspace.Read"`, `"Model.Infer"`, 1),
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	d := &Daemon{appsDir: dir}

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/apps.json", nil)
	rec := httptest.NewRecorder()
	d.handleSurfaceApps(rec, req, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("apps status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), dir) {
		t.Fatalf("the surface must not show the catalog's machine path: %s", rec.Body.String())
	}
	var resp surfaceAppsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Schema != SurfaceAppsSchema || len(resp.Apps) != 3 {
		t.Fatalf("apps resp = %+v", resp)
	}
	byName := map[string]surfaceAppEntry{}
	for _, a := range resp.Apps {
		byName[a.Name] = a
	}
	if b := byName["broken"]; b.Valid || b.Error == "" {
		t.Errorf("broken = %+v, want a parse error", b)
	}
	good := byName["docparse-convert"]
	if !good.Valid || good.Agent != "claude" || good.Mode != "headless" || len(good.Grants) != 4 ||
		len(good.Outputs) != 1 || good.Outputs[0] != "out/**/*.md" || len(good.Checks) != 1 || !strings.HasPrefix(good.Digest, "sha256:") {
		t.Errorf("docparse-convert = %+v", good)
	}
	if w := byName["wide"]; w.Valid || len(w.Problems) == 0 || w.Problems[0].Field != "preset.grants[0].effect" {
		t.Errorf("an app granting Model.Infer must be invalid on that grant: %+v", w)
	}

	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/apps.json?x=1", nil)
	rec = httptest.NewRecorder()
	d.handleSurfaceApps(rec, req, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("query string status = %d, want 400", rec.Code)
	}

	for q, want := range map[string]string{
		"app list":        "INVALID 1 problem(s): ailang-worldd app check wide",
		"app list --json": `"name": "docparse-convert"`,
	} {
		req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/verb.json?q="+url.QueryEscape(q), nil)
		rec = httptest.NewRecorder()
		d.handleSurfaceVerb(rec, req, nil)
		var ans verbAnswer
		if err := json.Unmarshal(rec.Body.Bytes(), &ans); err != nil {
			t.Fatalf("decode verb %q: %v", q, err)
		}
		if ans.Outcome != "result" || ans.Class != "read" || !strings.Contains(ans.Text+ans.JSON, want) {
			t.Errorf("verb %q = %+v, want a read result containing %q", q, ans, want)
		}
	}

	// An unknown flag is a 200 refusal, as for preset list.
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/verb.json?q="+url.QueryEscape("app list --all"), nil)
	rec = httptest.NewRecorder()
	d.handleSurfaceVerb(rec, req, nil)
	var refused verbAnswer
	_ = json.Unmarshal(rec.Body.Bytes(), &refused)
	if refused.Outcome == "result" || !strings.Contains(refused.Reason, "accepts only --json") {
		t.Errorf("app list --all = %+v, want a refusal", refused)
	}

	// An empty or missing catalog is a quiet, valid answer.
	empty := &Daemon{appsDir: filepath.Join(t.TempDir(), "absent")}
	rec = httptest.NewRecorder()
	empty.handleSurfaceApps(rec, httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/apps.json", nil), nil)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"apps":[]`) {
		t.Fatalf("missing catalog = %d %s, want 200 with an empty list", rec.Code, rec.Body.String())
	}
}
