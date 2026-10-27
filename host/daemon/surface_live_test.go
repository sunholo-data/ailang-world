package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

type surfaceLatestFailure struct{ readStore }

func (s surfaceLatestFailure) LogEntriesLatest(context.Context, int) ([]store.LogEntry, error) {
	return nil, errors.New("private database path and error")
}

func surfaceLiveDaemon(t *testing.T) *Daemon {
	t.Helper()
	rig, err := os.ReadFile("../surface/testdata/index-rig.json")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(path, rig, 0600); err != nil {
		t.Fatal(err)
	}
	d := newHandlerDaemon(t)
	d.cfg.SurfaceIndex = path
	return d
}

func TestSurfaceModelLiveCursor(t *testing.T) {
	d := surfaceLiveDaemon(t)
	read := func() surfacepage.NowEdge {
		t.Helper()
		w := httptest.NewRecorder()
		d.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/surface/model.json", nil))
		if w.Code != 200 {
			t.Fatalf("model: %d %s", w.Code, w.Body.String())
		}
		var m surfacepage.Model
		if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
			t.Fatal(err)
		}
		return m.NowEdge
	}
	if e := read(); e.Live == nil || e.Live.LastIndex != -1 || e.Live.LogHead != nil {
		t.Fatalf("empty log: %+v", e)
	}
	c := genesisCommit("surface-live")
	if err := d.store.Commit(boundedTestContext(t), c); err != nil {
		t.Fatal(err)
	}
	next := testCommit(c.NextWorld, c.Entry.Header.EntryIndex+1, "surface-live")
	if err := d.store.Commit(boundedTestContext(t), next); err != nil {
		t.Fatal(err)
	}
	if e := read(); e.Live == nil || e.Live.LastIndex != next.Entry.Header.EntryIndex || e.Live.LogHead == nil || *e.Live.LogHead != next.Entry.EntryHash.String() {
		t.Fatalf("newest entry: %+v", e.Live)
	}
	d.reads = surfaceLatestFailure{d.reads}
	e := read()
	b, _ := json.Marshal(e)
	if e.Live != nil || e.LiveUnavailable != "the log could not be read" || strings.Contains(string(b), "private") {
		t.Fatalf("read failure: %s", b)
	}
}
