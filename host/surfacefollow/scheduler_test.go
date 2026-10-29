package surfacefollow

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDeadlines(t *testing.T) {
	t.Run("call-keeps-last-good", func(t *testing.T) {
		g, clock := newGit(), newClock()
		s, _ := startFollower(t, []Repo{testRepo()}, g, clock)
		until(t, func() bool { return s.accepted(sha(3)) })
		entered := make(chan struct{}, 1)
		g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
			if a[3] == "log" {
				entered <- struct{}{}
				<-ctx.Done()
				return "", ctx.Err(), true
			}
			return "", nil, false
		})
		g.move("/repo", sha(6))
		clock.advance(PollInterval)
		<-entered
		clock.advance(CallTimeout + time.Nanosecond)
		until(t, func() bool { return s.failed("timeout") })
		s.mu.Lock()
		defer s.mu.Unlock()
		last := s.all[len(s.all)-1]
		if last.Partition == nil || last.Partition.Head != sha(3) {
			t.Fatal("timeout replaced last good")
		}
	})
	t.Run("queue-counts-toward-job", func(t *testing.T) {
		g, clock := newGit(), newClock()
		repos := []Repo{}
		for i := 0; i < 3; i++ {
			r := testRepo()
			r.ID, r.Repo = fmt.Sprintf("repo%d", i), fmt.Sprintf("/repo%d", i)
			g.move(r.Repo, sha(3+i*3))
			repos = append(repos, r)
		}
		s, _ := startFollower(t, repos, g, clock)
		for i := 0; i < 3; i++ {
			head := sha(3 + i*3)
			until(t, func() bool { return s.accepted(head) })
		}
		release := make(chan struct{})
		entered := make(chan struct{}, 2)
		g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
			if a[3] == "log" {
				entered <- struct{}{}
				<-release
			}
			return "", nil, false
		})
		for i := 0; i < 3; i++ {
			g.move(repos[i].Repo, sha(30+i*3))
		}
		clock.advance(PollInterval)
		<-entered
		<-entered
		// The third repo has one queued job, with no worker yet.
		clock.advance(JobTimeout + time.Nanosecond)
		close(release)
		until(t, func() bool {
			s.mu.Lock()
			defer s.mu.Unlock()
			for _, st := range s.all {
				if st.ID == "repo2" && st.LastFailure != nil && st.LastFailure.Class == "timeout" &&
					st.Partition != nil && st.Partition.Head == sha(9) {
					return true
				}
			}
			return false
		})
		if countCommand(g, "log") != 5 {
			t.Fatal("expired queued job launched git")
		}
	})
}

func TestBackoff(t *testing.T) {
	g, clock := newGit(), newClock()
	var mu sync.Mutex
	attempts := []time.Time{}
	g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
		if a[3] == "log" {
			mu.Lock()
			attempts = append(attempts, clock.Now())
			mu.Unlock()
			return "", fmt.Errorf("source missing"), true
		}
		return "", nil, false
	})
	s, _ := startFollower(t, []Repo{testRepo()}, g, clock)
	waitAttempts := func(n int) {
		t.Helper()
		until(t, func() bool {
			s.mu.Lock()
			defer s.mu.Unlock()
			mu.Lock()
			count := len(attempts)
			mu.Unlock()
			if count < n {
				return false
			}
			last := s.all[len(s.all)-1]
			expected := retryDelay(n)
			if n == 5 {
				expected = retryDelay(1)
			}
			return last.LastFailure != nil && last.RetryAt.Equal(clock.Now().Add(expected))
		})
	}
	waitAttempts(1)
	for tick := 1; tick <= 7; tick++ {
		clock.advance(PollInterval)
		wantProbes := tick + 1
		until(t, func() bool { return countCommand(g, "rev-parse") >= wantProbes })
		if tick == 1 {
			waitAttempts(2)
		}
		if tick == 3 {
			waitAttempts(3)
		}
		if tick == 7 {
			waitAttempts(4)
		}
	}
	mu.Lock()
	got := append([]time.Time(nil), attempts...)
	mu.Unlock()
	if len(got) != 4 || got[1].Sub(got[0]) != 15*time.Second ||
		got[2].Sub(got[1]) != 30*time.Second || got[3].Sub(got[2]) != 60*time.Second {
		t.Fatalf("backoff attempts: %v", got)
	}
	g.move("/repo", sha(6))
	clock.advance(PollInterval)
	waitAttempts(5)
	mu.Lock()
	defer mu.Unlock()
	if len(attempts) != 5 {
		t.Fatal("changed head did not reset backoff")
	}
}
