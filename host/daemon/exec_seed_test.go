package daemon

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// Row 257 (w-exec-cache-seed-at-bind M2a): the bind starts the exec cache seed,
// and GET /v1/episodes/{ep}/report carries its live state.

// makeSeedDir writes <seed>/ab/entry = "seeded\n" under a fresh directory.
func makeSeedDir(t *testing.T) string {
	t.Helper()
	seed := filepath.Join(evalTempDir(t), "seed")
	if err := os.MkdirAll(filepath.Join(seed, "ab"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(seed, "ab", "entry"), []byte("seeded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return seed
}

// evalTempDir is t.TempDir with symlinks resolved (the daemon canonicalises).
func evalTempDir(t *testing.T) string {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return base
}

// newSeedDaemon is a daemon whose workspace exec profile "proj" seeds GOCACHE
// from seed. It has a fake srt tree and node (as health_exec_tool_status_test
// does), one workspace directory per episode, and no bindingFor.
func newSeedDaemon(t *testing.T, seed string, episodes ...string) (d *Daemon, root, stateDir string) {
	t.Helper()
	base := evalTempDir(t)
	stateDir, root = filepath.Join(base, "state"), filepath.Join(base, "workspace")
	home, tree := filepath.Join(base, "home"), filepath.Join(base, "node_modules")
	pkg := filepath.Join(tree, "@anthropic-ai", "sandbox-runtime")
	dirs := []string{stateDir, home, filepath.Join(pkg, "dist")}
	for _, ep := range episodes {
		dirs = append(dirs, filepath.Join(root, ep))
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(p, body string) {
		t.Helper()
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cli := "// fake srt"
	write(filepath.Join(pkg, "package.json"), `{"name":"@anthropic-ai/sandbox-runtime","version":"0.0.0-test"}`)
	write(filepath.Join(pkg, "dist", "cli.js"), cli)
	node := filepath.Join(base, "node")
	write(node, "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo v26.12.1; fi\nexit 0\n")
	sum := sha256.Sum256([]byte(cli))
	sb, err := broker.ArchiveExecSandbox(broker.ExecSandboxConfig{NodeModules: tree, Node: node, StateDir: stateDir,
		Pin: broker.ExecSandboxPin{Version: "0.0.0-test", CLISHA256: hex.EncodeToString(sum[:])}})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(map[string]any{
		"profile": broker.ExecProfileVersion, "project": "proj", "root": ".", "path": []string{"/usr/bin", "/bin"},
		"timeout_ms": 9000, "caches": map[string]any{"GOCACHE": map[string]any{"seed": seed}},
		"commands": map[string]any{"test": map[string]any{"argv": []string{"/bin/sh", "-c", "exit 0"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	p, err := broker.ParseExecProfile(raw)
	if err != nil {
		t.Fatal(err)
	}
	d = newHandlerDaemon(t)
	d.workspace.root, d.workspace.stateDir = root, stateDir
	d.workspace.bindingFor = nil
	d.workspace.exec = &workspaceExec{profile: p, byProject: map[string]*broker.ExecProfile{"proj": p}, sandbox: sb, operatorHome: home}
	t.Cleanup(func() {
		deadline := time.Now().Add(20 * time.Second)
		for _, ep := range episodes {
			for time.Now().Before(deadline) && hasCopying(d.workspace.execCacheStatus(ep)) {
				time.Sleep(10 * time.Millisecond)
			}
		}
	})
	return d, root, stateDir
}

func hasCopying(entries []ReportExecCache) bool {
	for _, e := range entries {
		if e.State == "copying" {
			return true
		}
	}
	return false
}

// waitDaemonSeed polls the episode's exec_cache until no entry is copying (<= 20 s).
func waitDaemonSeed(t *testing.T, d *Daemon, ep string) []ReportExecCache {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for {
		entries := d.workspace.execCacheStatus(ep)
		if !hasCopying(entries) {
			return entries
		}
		if time.Now().After(deadline) {
			t.Fatalf("seed for %s still copying after 20s: %+v", ep, entries)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// reportBody is the report fields these tests read.
type reportBody struct {
	TotalCalls int               `json:"total_calls"`
	ExecCache  []ReportExecCache `json:"exec_cache"`
}

func TestEpisodeReportShowsBindSeedFailure(t *testing.T) {
	missing := filepath.Join(evalTempDir(t), "no-such-seed")
	d, root, _ := newSeedDaemon(t, missing, "ep1")
	d.workspace.execHandler("ep1", filepath.Join(root, "ep1"))
	waitDaemonSeed(t, d, "ep1")
	rec := requestRecorder(t, d, http.MethodGet, "/v1/episodes/ep1/report", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("report status %d: %s", rec.Code, rec.Body.String())
	}
	var got reportBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	want := "seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)"
	if got.TotalCalls != 0 || len(got.ExecCache) != 1 {
		t.Fatalf("report: %s", rec.Body.String())
	}
	e := got.ExecCache[0]
	if e.Project != "proj" || e.Cache != "GOCACHE" || e.State != "failed" || e.Attempts != 1 || e.By != "bind" || e.Failure != want {
		t.Fatalf("exec_cache = %+v, want failed/1/bind %q", e, want)
	}
}

func TestBindSeedsAFreshEpisode(t *testing.T) {
	seed := makeSeedDir(t)
	d, root, stateDir := newSeedDaemon(t, seed, "ep1")
	if _, err := os.Lstat(filepath.Join(stateDir, "exec-cache")); !os.IsNotExist(err) {
		t.Fatalf("exec-cache exists before the bind: %v", err)
	}
	h, ok := d.workspace.execHandler("ep1", filepath.Join(root, "ep1")).(*broker.ExecHandler)
	if !ok || h == nil {
		t.Fatalf("execHandler did not return a *broker.ExecHandler")
	}
	entries := waitDaemonSeed(t, d, "ep1")
	if len(entries) != 1 || entries[0].State != "ready" || entries[0].Attempts != 1 || entries[0].By != "bind" {
		t.Fatalf("exec_cache = %+v, want one ready/1/bind", entries)
	}
	cache := filepath.Join(stateDir, "exec-cache", "ep1", "proj", "GOCACHE")
	if b, err := os.ReadFile(filepath.Join(cache, "ab", "entry")); err != nil || string(b) != "seeded\n" {
		t.Fatalf("seeded entry = %q, %v", b, err)
	}
	if b, err := os.ReadFile(filepath.Join(cache, ".world-seed-complete")); err != nil || string(b) != "world-exec-seed v1 GOCACHE\n" {
		t.Fatalf("marker = %q, %v", b, err)
	}
}

func TestEpisodeReportRacesFirstDispatch(t *testing.T) {
	seed := makeSeedDir(t)
	var eps []string
	for i := 1; i <= 8; i++ {
		eps = append(eps, fmt.Sprintf("ep%d", i))
	}
	d, root, _ := newSeedDaemon(t, seed, eps...)
	start := make(chan struct{})
	var wg sync.WaitGroup
	problems := make(chan string, 256)
	wg.Add(2)
	go func() { // what registry() does on a first dispatch
		defer wg.Done()
		<-start
		for _, ep := range eps {
			d.workspace.execHandler(ep, filepath.Join(root, ep))
		}
	}()
	go func() { // the report reader
		defer wg.Done()
		<-start
		for round := 0; round < 8; round++ {
			for _, ep := range eps {
				rec := httptest.NewRecorder()
				d.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/episodes/"+ep+"/report", nil))
				switch rec.Code {
				case http.StatusNotFound:
					if !strings.Contains(rec.Body.String(), "NotFound") {
						problems <- fmt.Sprintf("%s: 404 without class NotFound: %s", ep, rec.Body.String())
					}
				case http.StatusOK:
					var got reportBody
					if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
						problems <- fmt.Sprintf("%s: bad JSON: %v", ep, err)
						continue
					}
					if len(got.ExecCache) > 1 {
						problems <- fmt.Sprintf("%s: %d exec_cache entries", ep, len(got.ExecCache))
						continue
					}
					if len(got.ExecCache) == 1 {
						e := got.ExecCache[0]
						if e.Cache != "GOCACHE" || (e.State != "copying" && e.State != "ready") {
							problems <- fmt.Sprintf("%s: exec_cache entry %+v", ep, e)
						}
					}
				default:
					problems <- fmt.Sprintf("%s: status %d", ep, rec.Code)
				}
			}
		}
	}()
	close(start)
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(60 * time.Second):
		t.Fatal("race rig did not finish within 60s")
	}
	close(problems)
	for p := range problems {
		t.Error(p)
	}
	for _, ep := range eps {
		entries := waitDaemonSeed(t, d, ep)
		if len(entries) != 1 || entries[0].State != "ready" {
			t.Errorf("%s: exec_cache = %+v, want one ready entry", ep, entries)
		}
	}
}
