package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// TestSurfaceModelRoutes: GET /surface/model.json and /surface/stack.json
// (row 188 N1) over the rig fixture: headers, bodies against surfacepage's
// own projection, the closed grammar, empty and refused indexes, and POST.
func TestSurfaceModelRoutes(t *testing.T) {
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
	check := func(w *httptest.ResponseRecorder, status int, schema string) map[string]any {
		t.Helper()
		if w.Code != status {
			t.Fatalf("status = %d, want %d: %s", w.Code, status, w.Body.String())
		}
		for k, v := range map[string]string{"Content-Type": "application/json", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff"} {
			if got := w.Header().Get(k); got != v {
				t.Errorf("%s = %q, want %q", k, got, v)
			}
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("body is not JSON: %v", err)
		}
		if body["schema"] != schema {
			t.Errorf("schema = %v, want %s", body["schema"], schema)
		}
		return body
	}
	const model, stack = surfacepage.ModelSchema, surfacepage.StackSchema
	w189d := "f36428abd9a72e5093b61cba14e716c37019f1ef"
	stackOf := func(brane, at string) string { return "/surface/stack.json?brane=" + brane + "&at=" + at }

	t.Run("not built", func(t *testing.T) {
		for target, schema := range map[string]string{"/surface/model.json": model, stackOf("row-189", w189d): stack} {
			body := check(request("GET", target), 200, schema)
			if body["empty"] != true || !strings.Contains(body["reason"].(string), "has not been built yet") {
				t.Errorf("%s: %v", target, body)
			}
		}
	})
	t.Run("not configured", func(t *testing.T) {
		d.cfg.SurfaceIndex = ""
		defer func() { d.cfg.SurfaceIndex = path }()
		for target, schema := range map[string]string{"/surface/model.json": model, stackOf("row-189", w189d): stack} {
			body := check(request("GET", target), 200, schema)
			if body["empty"] != true || !strings.Contains(body["reason"].(string), "no surface index is configured") {
				t.Errorf("%s: %v", target, body)
			}
		}
	})

	rig, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rig, 0o600); err != nil {
		t.Fatal(err)
	}
	ix, err := surface.Decode(rig)
	if err != nil {
		t.Fatal(err)
	}
	creds, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	for _, row := range []store.SessionRow{
		{CredentialID: strings.Repeat("a", 64), EpisodeID: "w188i", GrantsJSON: "[]", CreatedAt: now - 360, ExpiresAt: now + 3600},
		{CredentialID: strings.Repeat("b", 64), EpisodeID: "w-old", GrantsJSON: "[]", CreatedAt: now - 7200, ExpiresAt: now - 1},
	} {
		if err := creds.MintSession(boundedTestContext(t), row); err != nil {
			t.Fatal(err)
		}
	}
	if err := creds.Close(); err != nil {
		t.Fatal(err)
	}

	t.Run("model", func(t *testing.T) {
		w := request("GET", "/surface/model.json")
		check(w, 200, model)
		var got surfacepage.Model
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		at, err := time.Parse(time.RFC3339, got.Now)
		if err != nil || at.Location() != time.UTC {
			t.Fatalf("now = %q", got.Now)
		}
		// The served lanes, states and links are surfacepage's for the
		// same index; only the clock-dependent words may move.
		want := surfacepage.BuildModel(ix, surfacepage.Now{At: at})
		for i := range want.Lanes {
			for j := range want.Lanes[i].States {
				want.Lanes[i].States[j].Rel = got.Lanes[i].States[j].Rel
			}
		}
		if !reflect.DeepEqual(got.Lanes, want.Lanes) || !reflect.DeepEqual(got.Links, want.Links) {
			t.Error("served lanes or links differ from BuildModel")
		}
		e := got.NowEdge
		if e.Health != "World healthy" || e.Checked != at.Format("15:04") || e.Daemon == "" || len(e.NeedsYou) != 0 {
			t.Errorf("now edge = %+v", e)
		}
		if len(e.Sessions) != 1 || e.Sessions[0].Episode != "w188i" || e.Sessions[0].Started != "6 min ago" {
			t.Errorf("sessions = %+v (an expired session is not live)", e.Sessions)
		}
	})

	t.Run("stack", func(t *testing.T) {
		for _, x := range ix.States {
			brane := x.Brane
			if x.Kind == surface.KindShared || brane == "" {
				brane = surfacepage.SharedKey
			}
			w := request("GET", stackOf(brane, x.ID))
			check(w, 200, stack)
			var got surfacepage.StackModel
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			want, err := surfacepage.BuildStack(ix, brane, x.ID, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			if got.Front.ID != x.ID || got.Front.Position != want.Front.Position || len(got.Rail) != len(want.Rail) ||
				len(got.Earlier) != len(want.Earlier) || len(got.Later) != len(want.Later) {
				t.Errorf("%s: served stack differs from BuildStack", x.ID)
			}
		}
		for _, target := range []string{stackOf("unknown", w189d), stackOf("row-189", "missing"), stackOf("shared", w189d)} {
			if body := check(request("GET", target), 404, stack); body["reason"] == "" {
				t.Errorf("%s: no reason", target)
			}
		}
	})

	t.Run("query", func(t *testing.T) {
		for _, target := range []string{"/surface/model.json?x=1", "/surface/model.json?brane=row-189&at=" + w189d, "/surface/model.json?world"} {
			if body := check(request("GET", target), 400, model); body["reason"] != unknownSurfaceKeyMessage {
				t.Errorf("%s: %v", target, body)
			}
		}
		for _, q := range []string{"", "?x=1", "?brane=row-189", "?at=x", "?brane=x&at=y&other=z", "?brane=x&brane=y&at=z", "?brane=x&at=y&at=z", "?brane=%zz&at=x", "?world"} {
			if body := check(request("GET", "/surface/stack.json"+q), 400, stack); body["reason"] != unknownSurfaceKeyMessage {
				t.Errorf("%s: %v", q, body)
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
			for target, schema := range map[string]string{"/surface/model.json": model, stackOf("row-189", w189d): stack} {
				w := request("GET", target)
				body := check(w, 503, schema)
				if body["refused"] != true || !strings.HasPrefix(body["reason"].(string), "the configured surface index: ") {
					t.Errorf("%s %s: %v", name, target, body)
				}
				if strings.Contains(w.Body.String(), path) || strings.Contains(w.Body.String(), filepath.Dir(path)) {
					t.Errorf("%s %s: configured path leaked", name, target)
				}
			}
		}
		d.cfg.SurfaceIndex = filepath.Dir(path)
		defer func() { d.cfg.SurfaceIndex = path }()
		if w := request("GET", "/surface/model.json"); w.Code != 503 || strings.Contains(w.Body.String(), d.cfg.SurfaceIndex) {
			t.Errorf("read error: %d %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST", func(t *testing.T) {
		for _, target := range []string{"/surface/model.json", stackOf("row-189", w189d)} {
			if w := request(http.MethodPost, target); w.Code != http.StatusMethodNotAllowed {
				t.Errorf("POST %s = %d", target, w.Code)
			}
		}
	})
}
