package daemon

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// Exercise the actual health handler without opening a listener.
func TestHealthExecToolsKeepsStartupRealpaths(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "26.0.0", "node")
	newPath := filepath.Join(dir, "26.12.1", "node")
	for _, path := range []string{old, newPath} {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(dir, "node")
	if err := os.Symlink(old, link); err != nil {
		t.Fatal(err)
	}
	e := &workspaceExec{byProject: map[string]*broker.ExecProfile{
		"project": {Commands: map[string]broker.ExecProfileCommand{
			"test": {ExecCommand: broker.ExecCommand{Argv: []string{"node"}}, Argv0: link},
		}},
	}}
	e.tools = captureExecTools(e)
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(newPath, link); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(old); err != nil {
		t.Fatal(err)
	}
	d := &Daemon{workspace: &workspaceTools{exec: e}}
	rec := httptest.NewRecorder()
	d.handleHealth(rec, httptest.NewRequest("GET", "/v1/health", nil))
	var h HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &h); err != nil {
		t.Fatal(err)
	}
	if h.ExecTools["node"] != old {
		t.Fatalf("startup path changed: %+v", h.ExecTools)
	}
	empty := (&workspaceTools{}).healthExecTools()
	if empty == nil || len(empty) != 0 {
		t.Fatalf("disabled exec: %+v", empty)
	}
}
