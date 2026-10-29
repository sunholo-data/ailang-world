package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

type followGit struct {
	gate  <-chan struct{}
	calls atomic.Int64
}

func (g *followGit) Git(ctx context.Context, a []string) (string, error) {
	g.calls.Add(1)
	if g.gate != nil {
		select {
		case <-g.gate:
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	switch a[3] {
	case "rev-parse":
		return strings.Repeat("3", 40), nil
	case "log":
		return strings.Repeat("3", 40) + "\t" + strings.Repeat("1", 40) + " " + strings.Repeat("2", 40) + "\t2026-10-29T10:00:00Z\tMerge pull request #7 from owner/trial/row220\n", nil
	case "cat-file":
		return "220. **surface-follow** · goal\n", nil
	case "merge-base":
		return strings.Repeat("1", 40), nil
	}
	return "", fmt.Errorf("unadmitted git command")
}

func followFixture(t *testing.T, n int) (Config, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	cfg := surfacefollow.Config{Schema: surfacefollow.ConfigSchema}
	for i := 0; i < n; i++ {
		repo := filepath.Join(root, fmt.Sprintf("repo%d", i))
		if err := os.Mkdir(repo, 0700); err != nil {
			t.Fatal(err)
		}
		cfg.Repos = append(cfg.Repos, surfacefollow.Repo{ID: fmt.Sprintf("repo%d", i), Repo: repo, Branch: "dev", Charter: "charter.md", Limit: 200})
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "follow.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "state"), 0700); err != nil {
		t.Fatal(err)
	}
	return Config{DBPath: filepath.Join(root, "state", "world.db"), BindHost: DefaultBindHost, SurfaceFollow: path, ErrorLog: &bytes.Buffer{}}, root
}

func newFollowDaemon(t *testing.T, cfg Config) *Daemon {
	t.Helper()
	d, err := New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func followRequest(d *Daemon, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, httptest.NewRequest("GET", path, nil))
	return w
}

func waitFollow(t *testing.T, pred func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !pred() {
		if time.Now().After(deadline) {
			t.Fatal("follow did not reach expected state")
		}
		time.Sleep(time.Millisecond)
	}
}

type followClock struct{}
type followTimer struct{ ch chan time.Time }

func (followClock) Now() time.Time {
	return time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
}
func (followClock) NewTimer(time.Duration) surfacefollow.Timer {
	return followTimer{ch: make(chan time.Time)}
}
func (t followTimer) C() <-chan time.Time {
	return t.ch
}
func (followTimer) Stop() bool {
	return true
}

// Exercise the same start seam Run calls after Listen, without opening a socket.
func startFollowTest(t *testing.T, d *Daemon, g *followGit) {
	t.Helper()
	d.follow.collector = surfacefollow.NewCollector(g, followClock{})
	d.startSurfaceFollow(boundedTestContext(t))
}

func TestSurfaceFollowRoutes(t *testing.T) {
	cfg, root := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	w := followRequest(d, "/surface/index.json")
	if w.Code != 503 || !strings.Contains(w.Body.String(), "repo0: building") || !strings.Contains(w.Body.String(), "repo1: building") {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	// Refusals have an id and a path-free class before the first accepted partition.
	d.follow.publisher.Sink(surfacefollow.Status{ID: "repo0", LastFailure: &surface.FollowFailure{Class: "source", Text: "Local source unavailable"}})
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	done := make(chan struct{})
	go func() {
		defer close(done)
		d.follow.publisher.Run(ctx)
	}()
	waitFollow(t, func() bool { return d.follow.publisher.Snapshot().Status.Repos[0].State == "refused" })
	w = followRequest(d, "/surface/index.json")
	if w.Code != 503 || !strings.Contains(w.Body.String(), "repo0: refused: source") {
		t.Fatal(w.Body)
	}
	cancel()
	<-done
	startFollowTest(t, d, &followGit{})
	waitFollow(t, func() bool { return len(d.follow.publisher.Snapshot().Envelope.Repos) == 2 })
	for _, path := range []string{"/surface/index.json", "/surface/follow-status.json"} {
		w = followRequest(d, path)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body)
		}
		for _, secret := range []string{root, cfg.SurfaceFollow, d.follow.config.Repos[0].Repo} {
			if strings.Contains(w.Body.String(), secret) {
				t.Fatalf("path leaked: %s", w.Body)
			}
		}
		if path == "/surface/index.json" {
			ix, err := surface.DecodeFollow(w.Body.Bytes())
			if err != nil || len(ix.Repos) != 2 {
				t.Fatalf("%s: %v", w.Body, err)
			}
			if !bytes.Equal(w.Body.Bytes(), d.follow.publisher.Snapshot().Bytes) {
				t.Fatal("envelope differs")
			}
		} else {
			var status surfacefollow.StatusDocument
			if json.Unmarshal(w.Body.Bytes(), &status) != nil || status.Schema != surfacefollow.StatusSchema || status.Mode != "follow" || len(status.Repos) != 2 ||
				status.Repos[0].State != "fresh" || status.Generation == "" {
				t.Fatal(w.Body)
			}
		}
	}
}

func TestSurfaceFollowFeedsV1ConsumersForOneRepo(t *testing.T) {
	for _, n := range []int{1, 2} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			cfg, _ := followFixture(t, n)
			d := newFollowDaemon(t, cfg)
			w := followRequest(d, "/surface/model.json")
			if w.Code != 200 || !strings.Contains(w.Body.String(), "the follow index is building") {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
			startFollowTest(t, d, &followGit{})
			waitFollow(t, func() bool { return len(d.follow.publisher.Snapshot().Envelope.Repos) == n })
			w = followRequest(d, "/surface/model.json")
			if n == 1 {
				if w.Code != 200 || !strings.Contains(w.Body.String(), strings.Repeat("3", 40)) {
					t.Fatalf("%d %s", w.Code, w.Body)
				}
				_, ix, empty, err := d.readSurfaceIndex()
				if err != nil || empty != "" || ix.Source.Head != strings.Repeat("3", 40) {
					t.Fatalf("%s %v", empty, err)
				}
			} else if w.Code != 503 || !strings.Contains(w.Body.String(), groupedFollowReason) {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
			for _, path := range []string{"/surface", "/surface/static", "/surface/stack.json?brane=row-220&at=" + strings.Repeat("2", 40)} {
				want := 200
				if n == 2 {
					want = 503
				}
				w := followRequest(d, path)
				if w.Code != want {
					t.Fatalf("%s: %d %s", path, w.Code, w.Body)
				}
			}
		})
	}
	// Legacy bytes and the new status route's neutral modes.
	cfg, root := followFixture(t, 1)
	cfg.SurfaceFollow = ""
	d := newFollowDaemon(t, cfg)
	if w := followRequest(d, "/surface/follow-status.json"); !strings.Contains(w.Body.String(), `"mode":"none"`) {
		t.Fatal(w.Body)
	}
	path := filepath.Join(root, "legacy.json")
	const legacy = `{"schema":"world/surface-index/v1","built_at":"2026-10-21T00:00:00Z","source":{},"branes":[],"states":[]}`
	if err := os.WriteFile(path, []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	d.cfg.SurfaceIndex = path
	if w := followRequest(d, "/surface/index.json"); w.Code != 200 || w.Body.String() != legacy {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	if w := followRequest(d, "/surface/follow-status.json"); !strings.Contains(w.Body.String(), `"mode":"index"`) || !strings.Contains(w.Body.String(), `"repos":[]`) {
		t.Fatal(w.Body)
	}
}

func TestSurfaceFollowStartupRefusals(t *testing.T) {
	for _, kind := range []string{"both", "mode", "cache-in-repo", "cache-symlink"} {
		t.Run(kind, func(t *testing.T) {
			cfg, root := followFixture(t, 1)
			switch kind {
			case "both":
				cfg.SurfaceIndex = filepath.Join(root, "index.json")
			case "mode":
				if err := os.Chmod(cfg.SurfaceFollow, 0644); err != nil {
					t.Fatal(err)
				}
			case "cache-in-repo":
				cfg.DBPath = filepath.Join(root, "repo0", "world.db")
			case "cache-symlink":
				if err := os.Symlink(filepath.Join(root, "repo0"), filepath.Join(root, "state", "surface-follow")); err != nil {
					t.Fatal(err)
				}
			}
			d, err := New(boundedTestContext(t), cfg)
			if d != nil {
				d.Close()
				t.Fatal("startup returned a daemon")
			}
			var startup *StartupError
			if !errors.As(err, &startup) || !strings.Contains(err.Error(), "surface follow") || !strings.Contains(err.Error(), "next: ailang-worldd serve --help") {
				t.Fatalf("%T %v", err, err)
			}
		})
	}
}

func TestSurfaceFollowStopsOnClose(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	if d.follow.cancel != nil {
		t.Fatal("New started follow")
	}
	gate := make(chan struct{})
	g := &followGit{gate: gate}
	startFollowTest(t, d, g)
	waitFollow(t, func() bool { return g.calls.Load() > 0 })
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	for _, done := range []chan struct{}{d.follow.followerDone, d.follow.publisherDone} {
		select {
		case <-done:
		default:
			t.Fatal("follow goroutine outlived Close")
		}
	}
}
