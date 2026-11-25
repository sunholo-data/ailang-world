package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestSurfaceOnboardingProfilesAndPresets(t *testing.T) {
	base := t.TempDir()
	profilesDir := filepath.Join(base, "profiles")
	presetsDir := filepath.Join(base, "presets")
	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(presetsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write one valid proposed profile and one malformed active profile.
	validProfile := `{
  "profile": "world/exec-profile/v1",
  "project": "docparse",
  "root": ".",
  "path": ["/bin", "/usr/bin"],
  "timeout_ms": 9000,
  "probe": ["/bin/sh", "-c", "true"],
  "commands": {
    "ailang-check": {"argv": ["/bin/sh", "-c", "true"], "positional": "relpath", "max_args": 4}
  }
}`
	if err := os.WriteFile(filepath.Join(profilesDir, "docparse.proposed.json"), []byte(validProfile), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profilesDir, "broken.json"), []byte("{not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write one valid preset.
	validPreset := `{
  "schema": "world/preset/v1",
  "name": "docparse",
  "repo": "/tmp/docparse",
  "remote": "origin",
  "base": "main",
  "agent": "agy",
  "profile": "docparse",
  "ttl_s": 43200,
  "grants": [{"effect": "Workspace.Read", "scope": "worktree", "budget": 5000}],
  "launch": {"mode": "headless", "deadline_s": 7200},
  "checks": [{"name": "check", "argv": ["ailang", "check"]}]
}`
	if err := os.WriteFile(filepath.Join(presetsDir, "docparse.json"), []byte(validPreset), 0o600); err != nil {
		t.Fatal(err)
	}

	d := &Daemon{
		execProfilesDir: profilesDir,
		presetsDir:      presetsDir,
	}

	// 1. GET /surface/profiles.json
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/profiles.json", nil)
	rec := httptest.NewRecorder()
	d.handleSurfaceProfiles(rec, req, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("profiles status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var profResp surfaceProfilesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &profResp); err != nil {
		t.Fatalf("decode profiles: %v", err)
	}
	if profResp.Schema != SurfaceProfilesSchema || len(profResp.Profiles) != 2 {
		t.Fatalf("profiles resp = %+v", profResp)
	}
	if profResp.Profiles[0].Project != "broken" || profResp.Profiles[0].Error == "" {
		t.Errorf("broken profile = %+v, want error", profResp.Profiles[0])
	}
	if profResp.Profiles[1].Project != "docparse" || profResp.Profiles[1].State != "proposed" || profResp.Profiles[1].Root != "." || len(profResp.Profiles[1].Commands) != 1 || profResp.Profiles[1].Commands[0] != "ailang-check" {
		t.Errorf("docparse profile = %+v", profResp.Profiles[1])
	}

	// 2. GET /surface/presets.json
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/presets.json", nil)
	rec = httptest.NewRecorder()
	d.handleSurfacePresets(rec, req, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("presets status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var presResp surfacePresetsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &presResp); err != nil {
		t.Fatalf("decode presets: %v", err)
	}
	if presResp.Schema != SurfacePresetsSchema || len(presResp.Presets) != 1 {
		t.Fatalf("presets resp = %+v", presResp)
	}
	p := presResp.Presets[0]
	if p.Name != "docparse" || p.Agent != "agy" || p.Mode != "headless" || p.Profile != "docparse" || p.GrantsCount != 1 || len(p.Checks) != 1 || p.Checks[0] != "check" || p.Digest == "" {
		t.Errorf("docparse preset = %+v", p)
	}

	// 3. Unknown query parameter is refused with 400.
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/presets.json?foo=bar", nil)
	rec = httptest.NewRecorder()
	d.handleSurfacePresets(rec, req, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("query param status = %d, want 400", rec.Code)
	}

	// 4. Address-bar Read verbs: profile list and preset list (text and --json).
	for _, q := range []string{"profile list", "profile list --json", "preset list", "preset list --json"} {
		req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/verb.json?q="+url.QueryEscape(q), nil)
		rec = httptest.NewRecorder()
		d.handleSurfaceVerb(rec, req, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("verb %q status = %d, body = %s", q, rec.Code, rec.Body.String())
		}
		var ans verbAnswer
		if err := json.Unmarshal(rec.Body.Bytes(), &ans); err != nil {
			t.Fatalf("decode verb %q: %v", q, err)
		}
		if ans.Outcome != "result" || ans.Class != "read" {
			t.Errorf("verb %q = %+v, want outcome=result class=read", q, ans)
		}
	}

	// 5. Standalone /surface/repos page and script.
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/repos", nil)
	rec = httptest.NewRecorder()
	d.handleSurfaceReposPage(rec, req, nil)
	if rec.Code != http.StatusOK || rec.Body.Len() == 0 {
		t.Fatalf("repos page status = %d len = %d", rec.Code, rec.Body.Len())
	}
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:7644/surface/repos.js", nil)
	rec = httptest.NewRecorder()
	d.handleSurfaceReposScript(rec, req, nil)
	if rec.Code != http.StatusOK || rec.Body.Len() == 0 {
		t.Fatalf("repos script status = %d len = %d", rec.Code, rec.Body.Len())
	}
}
