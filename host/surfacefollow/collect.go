package surfacefollow

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const (
	PollInterval = 15 * time.Second
	CallTimeout  = 2 * time.Second
	JobTimeout   = 20 * time.Second
	CallLimit    = 1200
)

// Git is the only process seam. Arguments include the canonical -C prefix.
type Git interface {
	Git(context.Context, []string) (string, error)
}

// ReadGit delegates to the existing bounded, environment-scrubbed launch site.
type ReadGit struct{}

func (ReadGit) Git(ctx context.Context, args []string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}
	return worktree.Output(ctx, git, args)
}

// Timer and Clock permit deterministic polling and deadline tests.
type Timer interface {
	C() <-chan time.Time
	Stop() bool
}

type Clock interface {
	Now() time.Time
	NewTimer(time.Duration) Timer
}

type realClock struct{}
type realTimer struct{ *time.Timer }

func (realClock) Now() time.Time                 { return time.Now() }
func (realClock) NewTimer(d time.Duration) Timer { return realTimer{time.NewTimer(d)} }
func (t realTimer) C() <-chan time.Time          { return t.Timer.C }

// bounded derives from the caller and also honors the injected clock. Closing
// the scope joins its timer goroutine, including on cancellation.
func bounded(parent context.Context, clock Clock, d time.Duration) (context.Context, func()) {
	ctx, cancel := context.WithTimeout(parent, max(d, 0))
	timer := clock.NewTimer(max(d, 0))
	done, exited := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(exited)
		select {
		case <-timer.C():
			cancel()
		case <-ctx.Done():
		case <-done:
		}
	}()
	return ctx, func() {
		cancel()
		timer.Stop()
		close(done)
		<-exited
	}
}

var errBudget = errors.New("git call budget exhausted")

// Collector holds only immutable-parent merge-base memoization and a rolling
// launch census. It is safe for the follower's two workers and probe loop.
type Collector struct {
	Git   Git
	Clock Clock
	mu    sync.Mutex
	bases map[string]string
	calls []time.Time
}

func NewCollector(git Git, clock Clock) *Collector {
	if git == nil {
		git = ReadGit{}
	}
	if clock == nil {
		clock = realClock{}
	}
	return &Collector{Git: git, Clock: clock, bases: map[string]string{}}
}

func (c *Collector) call(ctx context.Context, r Repo, probe bool, args ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	c.mu.Lock()
	now := c.Clock.Now()
	first := 0
	for first < len(c.calls) && !c.calls[first].After(now.Add(-time.Minute)) {
		first++
	}
	c.calls = c.calls[first:]
	limit := CallLimit
	if !probe {
		// Reserve a full minute of four-repo probes. Rebuilds cannot
		// consume the slots that keep unchanged sources observable.
		limit -= 16
	}
	if len(c.calls) >= limit {
		c.mu.Unlock()
		return "", errBudget
	}
	c.calls = append(c.calls, now)
	c.mu.Unlock()
	callCtx, closeScope := bounded(ctx, c.Clock, CallTimeout)
	defer closeScope()
	out, err := c.Git.Git(callCtx, append([]string{"-C", r.Repo, "--no-optional-locks"}, args...))
	if callCtx.Err() != nil {
		return "", context.DeadlineExceeded
	}
	return out, err
}

// Failure emits fixed browser-safe text. The diagnostic command uses placeholders
// so neither paths nor untrusted child output enter status.
func failure(class string) *surface.FollowFailure {
	text := map[string]string{
		"source":    "Local branch or charter unavailable; next: git -C <repo> show <branch>:<charter>",
		"timeout":   "Git call exceeded 2 s or queued job exceeded 20 s; next: ailang-worldd service status",
		"budget":    "1200 git calls per minute bound reached; retry next poll",
		"candidate": "Candidate refused; lower this repo limit; next: ailang-worldd service status",
		"cancelled": "Collection cancelled; retry next poll",
	}[class]
	return &surface.FollowFailure{Class: class, Text: text}
}

func classify(err error) *surface.FollowFailure {
	switch {
	case errors.Is(err, errBudget):
		return failure("budget")
	case errors.Is(err, context.DeadlineExceeded):
		return failure("timeout")
	case errors.Is(err, context.Canceled):
		return failure("cancelled")
	default:
		return failure("source")
	}
}

func validHead(head string) bool {
	_, err := surface.QualifyID("check", "commit", head)
	return err == nil
}

// Probe resolves only the configured local branch.
func (c *Collector) Probe(ctx context.Context, r Repo) (string, error) {
	out, err := c.call(ctx, r, true, "rev-parse", "--verify", "refs/heads/"+r.Branch+"^{commit}")
	head := strings.TrimSpace(out)
	if err != nil {
		return "", err
	}
	if !validHead(head) {
		return "", fmt.Errorf("invalid probe head")
	}
	return head, nil
}

// Collect probes once and returns the accepted partition unchanged at equal head.
func (c *Collector) Collect(ctx context.Context, r Repo, previous *surface.FollowPartition) (surface.FollowPartition, error) {
	ctx, done := bounded(ctx, c.Clock, JobTimeout)
	defer done()
	head, err := c.Probe(ctx, r)
	if err != nil {
		return surface.FollowPartition{}, err
	}
	if previous != nil && previous.Head == head {
		p := *previous
		p.CheckedAt = c.Clock.Now().UTC().Format(time.RFC3339)
		p.LastFailure = nil
		return p, nil
	}
	return c.BuildPinned(ctx, r, head)
}

// BuildPinned never resolves a branch. Every object read is anchored at head.
func (c *Collector) BuildPinned(ctx context.Context, r Repo, head string) (surface.FollowPartition, error) {
	if !validHead(head) {
		return surface.FollowPartition{}, fmt.Errorf("invalid pinned head")
	}
	log, err := c.call(ctx, r, false, "log", "--first-parent", "--merges",
		"--format=%x1e%H%x1f%P%x1f%cI%x1f%s%x1f%b", "-n", strconv.Itoa(r.Limit), head)
	if err != nil {
		return surface.FollowPartition{}, err
	}
	merges, err := surface.ParseMerges(log)
	if err != nil || len(merges) > r.Limit {
		return surface.FollowPartition{}, fmt.Errorf("candidate merge log refused")
	}
	charter, err := c.call(ctx, r, false, "cat-file", "blob", head+":"+r.Charter)
	if err != nil {
		return surface.FollowPartition{}, err
	}
	bases := map[string]string{}
	for _, m := range merges {
		key := r.Repo + "\x00" + m.Parents[0] + "\x00" + m.Parents[1]
		c.mu.Lock()
		base := c.bases[key]
		c.mu.Unlock()
		if base == "" {
			base, err = c.call(ctx, r, false, "merge-base", m.Parents[0], m.Parents[1])
			if err != nil {
				return surface.FollowPartition{}, err
			}
			base = strings.TrimSpace(base)
			if !validHead(base) {
				return surface.FollowPartition{}, fmt.Errorf("invalid merge base")
			}
			c.mu.Lock()
			// Bound memory across arbitrarily long daemon lifetimes.
			if len(c.bases) >= 8000 {
				c.bases = map[string]string{}
			}
			c.bases[key] = base
			c.mu.Unlock()
		}
		bases[m.Hash] = base
	}
	now := c.Clock.Now()
	ix, err := surface.Build(surface.BuildInput{Merges: merges, Bases: bases,
		Rev: r.Branch, Head: head, BuiltAt: now,
		Titles: surface.ChartRowTitles(charter), Goals: surface.ChartRowGoals(charter)})
	if err == nil {
		err = surface.Validate(ix)
	}
	if err != nil {
		return surface.FollowPartition{}, fmt.Errorf("candidate validation refused")
	}
	p := surface.FollowPartition{ID: r.ID, Head: head, BuiltAt: ix.BuiltAt,
		CheckedAt: now.UTC().Format(time.RFC3339), Index: ix}
	data, err := surface.Marshal(ix)
	if err != nil || len(data) > surface.MaxBytes {
		return surface.FollowPartition{}, fmt.Errorf("candidate bytes refused")
	}
	if ctx.Err() != nil {
		return surface.FollowPartition{}, context.DeadlineExceeded
	}
	return p, nil
}
