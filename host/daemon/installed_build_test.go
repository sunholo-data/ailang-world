package daemon

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestInstalledDigestCachedBySizeAndMtime(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "binary")
	write := func(s string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(s), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("first")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	first := installedDigestAt(path)
	if first == "" || first != executableDigestAt(path) {
		t.Fatal("initial digest", first)
	}
	write("other")
	if err := os.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	if got := installedDigestAt(path); got != first {
		t.Fatalf("same size and mtime: %s != %s", got, first)
	}
	write("different size")
	if got := installedDigestAt(path); got == first || got != executableDigestAt(path) {
		t.Fatal("size change was not rehashed", got)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if got := installedDigestAt(path); got != "" {
		t.Fatal("unreadable path", got)
	}
}

func TestSurfaceModelServiceFields(t *testing.T) {
	t.Setenv("XPC_SERVICE_NAME", "")
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	rig, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "index.json")
	if err := os.WriteFile(path, rig, 0600); err != nil {
		t.Fatal(err)
	}
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(root, "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/surface/model.json", nil))
	var m surfacepage.Model
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	s := m.NowEdge.Service
	want := strings.TrimPrefix(CurrentBuild().ExecutableSHA256, "sha256:")[:7]
	if w.Code != 200 || s == nil || s.Build != want || s.Supervisor != "none" || s.RestartPending || !s.Warn || s.StartedAt == "" || s.Line != "World is not running as a service. It stops when its terminal does. Run ailang-worldd service install." {
		t.Fatalf("service: %+v, response %s", s, w.Body.String())
	}
	raw, _ := json.Marshal(s)
	t.Logf("served nowEdge.service: %s", raw)
}
