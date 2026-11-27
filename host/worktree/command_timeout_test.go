package worktree

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestWithCommandTimeoutClampsGitCap(t *testing.T) {
	for _, tc := range []struct{ requested, want time.Duration }{
		{time.Hour, 30 * time.Second},
		{2 * time.Second, 2 * time.Second},
	} {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		ctx = WithCommandTimeout(ctx, tc.requested)
		_, err := Output(ctx, "/usr/bin/false", []string{"noop"})
		if err == nil || !strings.Contains(err.Error(), "timed out after "+tc.want.String()+":") {
			t.Fatalf("git cap for %v: %v", tc.requested, err)
		}
	}
}
