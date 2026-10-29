package surfacefollow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func publishRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func publisherFor(t *testing.T, repos ...Repo) *Publisher {
	t.Helper()
	root := publishRoot(t)
	repos = append([]Repo(nil), repos...)
	for i := range repos {
		repos[i].Repo = filepath.Join(root, repos[i].ID)
		if err := os.Mkdir(repos[i].Repo, 0700); err != nil {
			t.Fatal(err)
		}
	}
	p, err := NewPublisher(Config{Schema: ConfigSchema, Repos: repos}, []byte("config"), filepath.Join(root, "state"), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func partitionFor(t *testing.T, r Repo, n int) surface.FollowPartition {
	t.Helper()
	p, err := NewCollector(newGit(), newClock()).BuildPinned(context.Background(), r, sha(n))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func deliver(p *Publisher, part surface.FollowPartition) {
	p.Sink(Status{ID: part.ID, Head: part.Head, BuiltAt: part.BuiltAt, CheckedAt: part.CheckedAt, Partition: &part})
	p.publish()
}

func cacheBytes(t *testing.T, p *Publisher) []byte {
	t.Helper()
	b, err := os.ReadFile(p.path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestPublishAtomicGeneration(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	deliver(p, partitionFor(t, r, 3))
	old := p.Snapshot()
	newPart := partitionFor(t, r, 6)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	errs := make(chan string, 8)
	ready := make(chan struct{}, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = p.Snapshot()
			ready <- struct{}{}
			for {
				select {
				case <-stop:
					return
				default:
				}
				s := p.Snapshot()
				var ix surface.FollowIndex
				if json.Unmarshal(s.Bytes, &ix) != nil || ix.Generation != s.Status.Generation || ix.Generation != s.Envelope.Generation ||
					ix.Repos[0].Head != s.Status.Repos[0].AcceptedHead ||
					(ix.Repos[0].Head != old.Envelope.Repos[0].Head && ix.Repos[0].Head != newPart.Head) {
					errs <- "mixed snapshot"
					return
				}
			}
		}()
	}
	for i := 0; i < 8; i++ {
		<-ready
	}
	deliver(p, newPart)
	close(stop)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	if old.Status.Generation == p.Snapshot().Status.Generation {
		t.Fatal("generation did not change")
	}
	if old.Envelope.Repos[0].Head != sha(3) {
		t.Fatal("mutated old snapshot")
	}
}

type faultFile struct {
	cacheFile
	stage string
}

func (f faultFile) Write(b []byte) (int, error) {
	if f.stage == "write" {
		return 0, errors.New("write")
	}
	if f.stage == "short-write" {
		return len(b) - 1, nil
	}
	return f.cacheFile.Write(b)
}

func (f faultFile) Sync() error {
	if f.stage == "sync" {
		return errors.New("sync")
	}
	return f.cacheFile.Sync()
}

func (f faultFile) Close() error {
	err := f.cacheFile.Close()
	if f.stage == "close" {
		return errors.New("close")
	}
	return err
}

func TestPublishFailuresKeepLastGood(t *testing.T) {
	for _, stage := range []string{"create", "write", "short-write", "sync", "close", "rename", "dir-sync", "validation", "oversize"} {
		t.Run(stage, func(t *testing.T) {
			r := testRepo()
			p := publisherFor(t, r)
			deliver(p, partitionFor(t, r, 3))
			old, disk := p.Snapshot(), cacheBytes(t, p)
			next := partitionFor(t, r, 6)
			class := "cache-write"
			switch stage {
			case "validation":
				next.Index.Schema = "invalid"
				class = "candidate"
			case "oversize":
				// Source.rev has no v1 length cap: keep the fixture valid so
				// this exercises the aggregate byte refusal after validation.
				next.Index.Source.Rev = strings.Repeat("x", surface.MaxBytes+1)
				class = "oversize"
			case "create":
				p.ops.create = func(string, string) (cacheFile, error) { return nil, errors.New("create") }
			case "rename":
				p.ops.rename = func(string, string) error { return errors.New("rename") }
			case "dir-sync":
				p.ops.syncDir = func(string) error { return errors.New("sync") }
			default:
				create := p.ops.create
				p.ops.create = func(dir, pattern string) (cacheFile, error) {
					f, err := create(dir, pattern)
					if err != nil {
						return nil, err
					}
					return faultFile{cacheFile: f, stage: stage}, nil
				}
			}
			deliver(p, next)
			got := p.Snapshot()
			if got.Status.Generation != old.Status.Generation || got.Envelope.Repos[0].Head != sha(3) ||
				!bytes.Equal(disk, cacheBytes(t, p)) {
				t.Fatal("lost last good generation")
			}
			if got.Status.Repos[0].State != "stale" || got.Status.Repos[0].LastFailure.Class != class {
				t.Fatalf("status: %+v", got.Status.Repos[0])
			}
		})
	}
}

func TestStalePartitionRetained(t *testing.T) {
	a, b := testRepo(), testRepo()
	b.ID = "tools"
	p := publisherFor(t, a, b)
	a, b = p.config.Repos[0], p.config.Repos[1]
	g, clock := newGit(), newClock()
	g.move(a.Repo, sha(3))
	g.move(b.Repo, sha(6))
	ctx, cancel := context.WithCancel(context.Background())
	fd, pd := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(pd)
		p.Run(ctx)
	}()
	go func() {
		defer close(fd)
		_ = (&Follower{Config: p.config, Collector: NewCollector(g, clock), Sink: p.Sink}).Run(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		for _, done := range []chan struct{}{fd, pd} {
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Error("follower or publisher did not stop")
			}
		}
	})
	until(t, func() bool { return len(p.Snapshot().Envelope.Repos) == 2 })
	g.setHook(func(_ context.Context, args []string) (string, error, bool) {
		if args[1] == a.Repo && args[3] == "log" {
			return "", errors.New("private source diagnostic"), true
		}
		return "", nil, false
	})
	g.move(a.Repo, sha(9))
	g.move(b.Repo, sha(12))
	clock.advance(PollInterval)
	until(t, func() bool {
		s := p.Snapshot()
		return s.Status.Repos[0].State == "stale" && s.Status.Repos[1].State == "fresh" &&
			s.Status.Repos[1].AcceptedHead == sha(12)
	})
	s := p.Snapshot()
	if len(s.Envelope.Repos) != 2 || s.Envelope.Repos[0].Head != sha(3) || s.Envelope.Repos[1].Head != sha(12) ||
		s.Envelope.Repos[0].LastFailure.Class != "source" {
		t.Fatalf("%+v", s.Status)
	}
}

func TestStatusOnlyChangeKeepsTheCache(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	part := partitionFor(t, r, 3)
	deliver(p, part)
	disk, gen := cacheBytes(t, p), p.Snapshot().Status.Generation
	writes := 0
	create := p.ops.create
	p.ops.create = func(d, s string) (cacheFile, error) {
		writes++
		return create(d, s)
	}
	g, clock := newGit(), newClock()
	collector := NewCollector(g, clock)
	clock.advance(time.Minute)
	part, err := collector.Collect(context.Background(), r, &part)
	if err != nil {
		t.Fatal(err)
	}
	deliver(p, part)
	if writes != 0 || !bytes.Equal(disk, cacheBytes(t, p)) || p.Snapshot().Status.Generation != gen ||
		p.Snapshot().Envelope.Repos[0].CheckedAt != part.CheckedAt {
		t.Fatal("status-only publication wrote cache or lost status")
	}
	g.move(r.Repo, sha(6))
	next, err := collector.Collect(context.Background(), r, &part)
	if err != nil {
		t.Fatal(err)
	}
	deliver(p, next)
	if writes != 1 || p.Snapshot().Status.Generation == gen {
		t.Fatal("new head did not write exactly once")
	}
}

func TestRestartLoadsCacheAsStale(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	deliver(p, partitionFor(t, r, 3))
	restarted, err := NewPublisher(p.config, []byte("config"), filepath.Dir(filepath.Dir(p.path)), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	st := restarted.Snapshot().Status.Repos[0]
	if st.State != "stale" || st.LastFailure.Class != "restart" || len(restarted.Snapshot().Bytes) == 0 || restarted.Snapshot().Status.Generation != p.Snapshot().Status.Generation {
		t.Fatalf("%+v", st)
	}
	g, c := newGit(), newClock()
	g.move(p.config.Repos[0].Repo, sha(3))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fd, pd := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(pd)
		restarted.Run(ctx)
	}()
	go func() {
		defer close(fd)
		_ = (&Follower{Config: p.config, Collector: NewCollector(g, c), Sink: restarted.Sink, Seed: restarted.Seed()}).Run(ctx)
	}()
	until(t, func() bool { return restarted.Snapshot().Status.Repos[0].State == "fresh" })
	cancel()
	for _, done := range []chan struct{}{fd, pd} {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("did not stop")
		}
	}
	for _, cmd := range []string{"log", "cat-file", "merge-base"} {
		if countCommand(g, cmd) != 0 {
			t.Fatalf("restart rebuilt: %v", g.recorded())
		}
	}
}

func TestCacheKeyedByConfig(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	oldPath := filepath.Join(filepath.Dir(p.path), "index-old.json")
	if err := os.WriteFile(oldPath, []byte("old config"), 0600); err != nil {
		t.Fatal(err)
	}
	deliver(p, partitionFor(t, r, 3))
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatal("old cache was not removed after the first successful publish")
	}
	for path, mode := range map[string]os.FileMode{filepath.Dir(p.path): 0700, p.path: 0600} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != mode {
			t.Fatalf("cache permissions: %v %v", info, err)
		}
	}
	dir := filepath.Dir(filepath.Dir(p.path))
	other, err := NewPublisher(p.config, []byte("different"), dir, nil, nil)
	if err != nil || len(other.Snapshot().Envelope.Repos) != 0 || other.path == p.path {
		t.Fatal("loaded another config")
	}
	if err := os.WriteFile(p.path, []byte("corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	cold, err := NewPublisher(p.config, []byte("config"), dir, nil, nil)
	if err != nil || len(cold.Snapshot().Envelope.Repos) != 0 {
		t.Fatal("corrupt cache accepted")
	}
	// A cache for a different branch is ignored even with a matching filename.
	deliver(p, partitionFor(t, r, 6))
	cfg := p.config
	cfg.Repos = append([]Repo(nil), cfg.Repos...)
	cfg.Repos[0].Branch = "main"
	cold, err = NewPublisher(cfg, []byte("config"), dir, nil, nil)
	if err != nil || len(cold.Snapshot().Envelope.Repos) != 0 {
		t.Fatal("branch mismatch accepted")
	}
}

func TestRepoStates(t *testing.T) {
	part := partitionFor(t, testRepo(), 3)
	for _, tc := range []struct {
		want   string
		status Status
	}{
		{"building", Status{}}, {"refused", Status{LastFailure: failure("source")}},
		{"fresh", Status{Head: part.Head, Partition: &part}},
		{"updating", Status{Head: sha(6), Partition: &part}},
		{"stale", Status{Head: part.Head, Partition: &part, LastFailure: failure("source")}},
	} {
		if got := repoStatus(tc.status).State; got != tc.want {
			t.Errorf("got %s want %s", got, tc.want)
		}
	}
}

// A blocked cache writer must not hold the follower's event loop or queue
// every intermediate observation.
func TestPublisherSinkLatestWins(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	first, latest := partitionFor(t, r, 3), partitionFor(t, r, 9)
	create := p.ops.create
	entered, release := make(chan struct{}), make(chan struct{})
	var unblock sync.Once
	writes := 0
	p.ops.create = func(dir, pattern string) (cacheFile, error) {
		writes++
		if writes == 1 {
			close(entered)
			<-release
		}
		return create(dir, pattern)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.Run(ctx)
	}()
	t.Cleanup(func() {
		unblock.Do(func() { close(release) })
		cancel()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("publisher did not stop")
		}
	})
	p.Sink(Status{ID: r.ID, Head: first.Head, CheckedAt: first.CheckedAt, Partition: &first})
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("publisher did not enter cache write")
	}
	sent := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			p.Sink(Status{ID: r.ID, Head: latest.Head, CheckedAt: latest.CheckedAt, Partition: &latest})
		}
		close(sent)
	}()
	select {
	case <-sent:
	case <-time.After(time.Second):
		t.Fatal("Sink waited for disk I/O")
	}
	unblock.Do(func() { close(release) })
	until(t, func() bool { return p.Snapshot().Status.Repos[0].AcceptedHead == latest.Head })
	cancel()
	<-done
	if writes != 2 {
		t.Fatalf("queued intermediate statuses: %d writes", writes)
	}
}
