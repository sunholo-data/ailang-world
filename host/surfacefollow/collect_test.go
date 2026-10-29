package surfacefollow

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeTimer struct {
	clock *fakeClock
	at    time.Time
	ch    chan time.Time
	stop  bool
}

func (t *fakeTimer) C() <-chan time.Time { return t.ch }
func (t *fakeTimer) Stop() bool {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	old := t.stop
	t.stop = true
	return !old
}

type fakeClock struct {
	mu     sync.Mutex
	now    time.Time
	timers []*fakeTimer
}

func newClock() *fakeClock {
	return &fakeClock{now: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
}
func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}
func (c *fakeClock) NewTimer(d time.Duration) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &fakeTimer{clock: c, at: c.now.Add(d), ch: make(chan time.Time, 1)}
	c.timers = append(c.timers, t)
	if d <= 0 {
		t.stop = true
		t.ch <- c.now
	}
	return t
}
func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	for _, t := range c.timers {
		if !t.stop && !t.at.After(c.now) {
			t.stop = true
			t.ch <- c.now
		}
	}
}

func sha(n int) string { return fmt.Sprintf("%040x", n) }

type fakeGit struct {
	mu    sync.Mutex
	heads map[string]string
	calls [][]string
	hook  func(context.Context, []string) (string, error, bool)
}

func newGit() *fakeGit { return &fakeGit{heads: map[string]string{"/repo": sha(3)}} }
func (g *fakeGit) Git(ctx context.Context, a []string) (string, error) {
	g.mu.Lock()
	g.calls = append(g.calls, append([]string(nil), a...))
	hook := g.hook
	head := g.heads[a[1]]
	g.mu.Unlock()
	if hook != nil {
		if out, err, handled := hook(ctx, a); handled {
			return out, err
		}
	}
	switch a[3] {
	case "rev-parse":
		return head, nil
	case "log":
		h := a[len(a)-1]
		return h + "\t" + sha(1) + " " + sha(2) + "\t2026-10-29T10:00:00Z\tMerge pull request #7 from owner/trial/row220\n", nil
	case "cat-file":
		return "220. **surface-follow** · goal\n", nil
	case "merge-base":
		return sha(1), nil
	}
	return "", fmt.Errorf("unadmitted command")
}
func (g *fakeGit) move(repo, head string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.heads[repo] = head
}
func (g *fakeGit) recorded() [][]string {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([][]string, len(g.calls))
	for i, a := range g.calls {
		out[i] = append([]string(nil), a...)
	}
	return out
}
func (g *fakeGit) setHook(h func(context.Context, []string) (string, error, bool)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.hook = h
}

func testRepo() Repo {
	return Repo{ID: "world", Repo: "/repo", Charter: "charter.md", Branch: "dev", Limit: 200}
}

func until(t *testing.T, pred func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !pred() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for test synchronization")
		}
		runtime.Gosched()
	}
}

func countCommand(g *fakeGit, command string) int {
	n := 0
	for _, a := range g.recorded() {
		if a[3] == command {
			n++
		}
	}
	return n
}

type statuses struct {
	mu  sync.Mutex
	all []Status
}

func (s *statuses) sink(v Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.all = append(s.all, v)
}
func (s *statuses) accepted(head string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.all {
		if v.Partition != nil && v.Partition.Head == head {
			return true
		}
	}
	return false
}
func (s *statuses) failed(class string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.all {
		if v.LastFailure != nil && v.LastFailure.Class == class {
			return true
		}
	}
	return false
}
func startFollower(t *testing.T, repos []Repo, g *fakeGit, c *fakeClock) (*statuses, func()) {
	t.Helper()
	s := &statuses{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	f := &Follower{Config: Config{Schema: ConfigSchema, Repos: repos}, Collector: NewCollector(g, c), Sink: s.sink}
	go func() { done <- f.Run(ctx) }()
	stop := func() {
		cancel()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Run did not join workers")
		}
	}
	t.Cleanup(stop)
	return s, func() {
		// Cleanup is the sole consumer of done.
		cancel()
	}
}

func TestCollectPinsTheHead(t *testing.T) {
	g, clock := newGit(), newClock()
	g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
		if a[3] == "rev-parse" {
			g.move("/repo", sha(6))
			return sha(3), nil, true
		}
		return "", nil, false
	})
	p, err := NewCollector(g, clock).Collect(context.Background(), testRepo(), nil)
	if err != nil || p.Head != sha(3) {
		t.Fatalf("pinned collection: %v", err)
	}
	for _, a := range g.recorded() {
		if a[3] == "log" && a[len(a)-1] != sha(3) {
			t.Fatal("log used moving branch")
		}
		if a[3] == "cat-file" && a[5] != sha(3)+":charter.md" {
			t.Fatal("charter used moving branch")
		}
	}
	t.Logf("cold argv: %v", g.recorded())
}

func TestUnchangedSourcesProbeOnly(t *testing.T) {
	g, clock := newGit(), newClock()
	repos := []Repo{}
	for i := 0; i < 4; i++ {
		r := testRepo()
		r.ID, r.Repo = fmt.Sprintf("repo%d", i), fmt.Sprintf("/repo%d", i)
		g.move(r.Repo, sha(3+i*3))
		repos = append(repos, r)
	}
	s, _ := startFollower(t, repos, g, clock)
	for i := 0; i < 4; i++ {
		head := sha(3 + i*3)
		until(t, func() bool { return s.accepted(head) })
	}
	initial := len(g.recorded())
	for tick := 1; tick <= 4; tick++ {
		clock.advance(PollInterval)
		want := 4 + 4*tick
		until(t, func() bool { return countCommand(g, "rev-parse") >= want })
	}
	calls := g.recorded()[initial:]
	if len(calls) > 16 {
		t.Fatalf("unchanged minute: %d calls", len(calls))
	}
	for _, a := range calls {
		if a[3] != "rev-parse" {
			t.Fatalf("unexpected rebuild: %v", a)
		}
	}
	t.Logf("unchanged argv: %v", calls[0])
}

func movingScenario(t *testing.T, moves int) {
	g, clock := newGit(), newClock()
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
		if a[3] == "log" && a[len(a)-1] == sha(3) {
			once.Do(func() { close(entered) })
			// Model a child returning late even after cancellation.
			<-release
		}
		return "", nil, false
	})
	s, _ := startFollower(t, []Repo{testRepo()}, g, clock)
	until(t, func() bool {
		select {
		case <-entered:
			return true
		default:
			return false
		}
	})
	for i := 1; i <= moves; i++ {
		g.move("/repo", sha(3+3*i))
		clock.advance(PollInterval)
		want := i + 1
		until(t, func() bool { return countCommand(g, "rev-parse") >= want })
		// Wait until the event loop observes the probe, not just the fake.
		until(t, func() bool {
			s.mu.Lock()
			defer s.mu.Unlock()
			for _, st := range s.all {
				if st.Head == sha(3+3*i) {
					return true
				}
			}
			return false
		})
	}
	close(release)
	latest := sha(3 + 3*moves)
	until(t, func() bool { return s.accepted(latest) })
	if s.accepted(sha(3)) {
		t.Fatal("stale job delivered")
	}
	if n := countCommand(g, "log"); n != 2 {
		t.Fatalf("got %d builds, want first plus newest", n)
	}
}
func TestMovingRefCoalesces(t *testing.T)    { movingScenario(t, 3) }
func TestStaleJobCannotDeliver(t *testing.T) { movingScenario(t, 1) }

func TestRewindRebuilds(t *testing.T) {
	g, clock := newGit(), newClock()
	s, _ := startFollower(t, []Repo{testRepo()}, g, clock)
	until(t, func() bool { return s.accepted(sha(3)) })
	g.move("/repo", sha(1))
	// Use an older but distinct full head; the log fixture must have unique
	// merge and step IDs, so supply an empty merge window for this head.
	g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
		if a[3] == "log" && a[len(a)-1] == sha(1) {
			return "", nil, true
		}
		return "", nil, false
	})
	clock.advance(PollInterval)
	until(t, func() bool { return s.accepted(sha(1)) })
	if countCommand(g, "log") != 2 {
		t.Fatal("rewind did not rebuild")
	}
}

func TestMergeBaseCache(t *testing.T) {
	g, clock := newGit(), newClock()
	c := NewCollector(g, clock)
	for _, head := range []string{sha(3), sha(6)} {
		if _, err := c.BuildPinned(context.Background(), testRepo(), head); err != nil {
			t.Fatal(err)
		}
	}
	if countCommand(g, "merge-base") != 1 {
		t.Fatal("immutable parent pair cache missed")
	}
}

func TestOnlyReadCommands(t *testing.T) {
	g, clock := newGit(), newClock()
	c := NewCollector(g, clock)
	p, err := c.Collect(context.Background(), testRepo(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Collect(context.Background(), testRepo(), &p); err != nil {
		t.Fatal(err)
	}
	for _, a := range g.recorded() {
		valid := len(a) >= 4 && reflect.DeepEqual(a[:3], []string{"-C", "/repo", "--no-optional-locks"})
		switch a[3] {
		case "rev-parse":
			valid = valid && reflect.DeepEqual(a[3:], []string{"rev-parse", "--verify", "refs/heads/dev^{commit}"})
		case "log":
			valid = valid && len(a) == 10 && reflect.DeepEqual(a[3:9], []string{"log", "--first-parent", "--merges", "--format=%x1e%H%x1f%P%x1f%cI%x1f%s%x1f%b", "-n", strconv.Itoa(testRepo().Limit)}) && validHead(a[9])
		case "merge-base":
			valid = valid && len(a) == 6 && validHead(a[4]) && validHead(a[5])
		case "cat-file":
			valid = valid && len(a) == 6 && a[4] == "blob" && a[5] == sha(3)+":charter.md"
		default:
			valid = false
		}
		if !valid {
			t.Fatalf("unadmitted argv: %v", a)
		}
		for _, arg := range a {
			if strings.HasPrefix(arg, "world:") {
				t.Fatal("qualified id passed to git")
			}
		}
	}
}

func TestCallBudget(t *testing.T) {
	g, clock := newGit(), newClock()
	c := NewCollector(g, clock)
	ctx := context.Background()
	for i := 0; i < CallLimit; i++ {
		_, _ = c.call(ctx, testRepo(), false, "merge-base", sha(1), sha(2))
	}
	if len(g.recorded()) != CallLimit-16 {
		t.Fatal("rebuilds consumed probe reservation")
	}
	for i := 0; i < 16; i++ {
		if _, err := c.Probe(ctx, testRepo()); err != nil {
			t.Fatal("probe starved", err)
		}
	}
	if _, err := c.Probe(ctx, testRepo()); err != errBudget {
		t.Fatal("rolling budget exceeded")
	}
	if len(g.recorded()) != CallLimit {
		t.Fatal("incorrect launch census")
	}
	clock.advance(time.Minute)
	if _, err := c.Probe(ctx, testRepo()); err != nil {
		t.Fatal("budget did not expire", err)
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	for _, command := range []string{"rev-parse", "log"} {
		t.Run(command, func(t *testing.T) {
			g, clock := newGit(), newClock()
			entered, childDone := make(chan struct{}, 1), make(chan struct{})
			g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
				if a[3] != command {
					return "", nil, false
				}
				entered <- struct{}{}
				<-ctx.Done()
				close(childDone)
				return "", ctx.Err(), true
			})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan struct{})
			go func() {
				defer close(done)
				_ = (&Follower{Config: Config{Repos: []Repo{testRepo()}}, Collector: NewCollector(g, clock)}).Run(ctx)
			}()
			<-entered
			cancel()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("worker or probe remained running")
			}
			select {
			case <-childDone:
			default:
				t.Fatal("Run returned before its child finished")
			}
		})
	}
}
