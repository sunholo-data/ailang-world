package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A pinned local base is used even if HEAD moves; an existing branch checks
// out directly and remains protected by git's one-worktree-per-branch rule.
func TestAddStartPointAndExistingBranch(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	ctx := testCtx(t)
	query := func(args ...string) string {
		t.Helper()
		s, err := Output(ctx, git, append([]string{"-C", repo}, args...))
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(s)
	}
	base := query("rev-parse", "HEAD")
	query("branch", "dev")
	query("-c", "user.email=test@example.invalid", "-c", "user.name=test", "commit", "-q", "--allow-empty", "-m", "feature")
	dir := filepath.Join(t.TempDir(), "ep")
	if err := Add(ctx, git, repo, dir, "world/ep", base); err != nil {
		t.Fatal(err)
	}
	head, err := Output(ctx, git, []string{"-C", dir, "rev-parse", "HEAD"})
	if err != nil || strings.TrimSpace(head) != base {
		t.Fatalf("%s %v", head, err)
	}
	other := filepath.Join(t.TempDir(), "existing")
	if err := AddExisting(ctx, git, repo, other, "dev"); err != nil {
		t.Fatal(err)
	}
	branch, err := Output(ctx, git, []string{"-C", other, "branch", "--show-current"})
	if err != nil || strings.TrimSpace(branch) != "dev" {
		t.Fatalf("%s %v", branch, err)
	}
	if err := AddExisting(ctx, git, repo, filepath.Join(t.TempDir(), "busy"), "dev"); err == nil {
		t.Fatal("busy branch accepted")
	}
	if _, err := os.Stat(filepath.Join(other, ".git")); err != nil {
		t.Fatal(err)
	}
}
