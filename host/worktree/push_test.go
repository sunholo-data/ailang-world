package worktree

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPushLocalNeverForces(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git unavailable")
	}
	ctx := context.Background()
	root := t.TempDir()
	repo, remote := filepath.Join(root, "repo"), filepath.Join(root, "remote.git")
	query := func(args ...string) string {
		t.Helper()
		out, err := Output(ctx, git, args)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(out)
	}
	query("init", "-q", repo)
	query("init", "--bare", "-q", remote)
	query("-C", repo, "config", "user.name", "Fixture")
	query("-C", repo, "config", "user.email", "fixture@example.invalid")
	query("-C", repo, "config", "commit.gpgsign", "false")
	query("-C", repo, "commit", "--allow-empty", "-qm", "base")
	base := query("-C", repo, "rev-parse", "HEAD")
	query("-C", repo, "commit", "--allow-empty", "-qm", "remote tip")
	tip := query("-C", repo, "rev-parse", "HEAD")
	if err := Push(ctx, git, repo, remote, tip, "world/ep1"); err != nil {
		t.Fatal(err)
	}
	query("-C", repo, "reset", "--hard", base)
	if err := os.WriteFile(filepath.Join(repo, "different"), []byte("local branch"), 0o600); err != nil {
		t.Fatal(err)
	}
	query("-C", repo, "add", "different")
	query("-C", repo, "commit", "-qm", "diverged")
	sha := query("-C", repo, "rev-parse", "HEAD")
	if err := Push(ctx, git, repo, remote, sha, "world/ep1"); err == nil {
		t.Fatal("non-fast-forward push succeeded")
	}
	if got := query("--git-dir", remote, "rev-parse", "refs/heads/world/ep1"); got != tip {
		t.Fatal("remote was forced")
	}
	for _, branch := range []string{"+main", "-f", "a:b"} {
		if err := Push(ctx, git, repo, remote, sha, branch); err == nil {
			t.Fatal("accepted unsafe branch")
		}
	}
}
