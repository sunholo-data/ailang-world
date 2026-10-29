package surfacefollow

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestMissingSourcesKeepLastPartition(t *testing.T) {
	for _, command := range []string{"rev-parse", "cat-file"} {
		t.Run(command, func(t *testing.T) {
			g, clock := newGit(), newClock()
			s, _ := startFollower(t, []Repo{testRepo()}, g, clock)
			until(t, func() bool { return s.accepted(sha(3)) })
			g.move("/repo", sha(6))
			g.setHook(func(ctx context.Context, a []string) (string, error, bool) {
				if a[3] == command {
					return "", fmt.Errorf("/private/source secret-token"), true
				}
				return "", nil, false
			})
			clock.advance(PollInterval)
			until(t, func() bool { return s.failed("source") })
			s.mu.Lock()
			defer s.mu.Unlock()
			last := s.all[len(s.all)-1]
			if last.Partition == nil || last.Partition.Head != sha(3) {
				t.Fatal("source failure discarded last good")
			}
			if strings.Contains(last.LastFailure.Text, "secret-token") ||
				strings.Contains(last.LastFailure.Text, "/private") ||
				!strings.Contains(last.LastFailure.Text, "git -C <repo> show <branch>:<charter>") {
				t.Fatal("source refusal lacks safe cause and next command")
			}
		})
	}
}
