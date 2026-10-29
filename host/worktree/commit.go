package worktree

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// Commit stages only the named local paths, in bounded literal-pathspec
// batches. Hooks and fsmonitor are disabled; signing and committer identity
// remain the repository's own configuration.
func Commit(ctx context.Context, git, dir, agent string, paths []string, messageParts ...string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("worktree: no paths to commit")
	}
	for _, p := range paths {
		if !filepath.IsLocal(p) || filepath.Clean(p) == "." {
			return "", fmt.Errorf("worktree: invalid commit path %q", p)
		}
	}
	if strings.ContainsAny(agent, "\r\n<>") {
		return "", fmt.Errorf("worktree: invalid author")
	}
	// Do not accidentally commit a previously staged, unrelated path.
	if err := runGit(ctx, git, []string{"-C", dir, "-c", "core.fsmonitor=false", "reset", "-q", "--mixed", "HEAD"}, nil); err != nil {
		return "", err
	}
	for lo := 0; lo < len(paths); lo += 64 {
		args := []string{"--literal-pathspecs", "-C", dir, "-c", "core.fsmonitor=false", "add", "-A", "--"}
		args = append(args, paths[lo:min(lo+64, len(paths))]...)
		if err := runGit(ctx, git, args, nil); err != nil {
			return "", err
		}
	}
	args := []string{"-C", dir, "-c", "core.fsmonitor=false", "-c", "core.hooksPath=/dev/null",
		"commit", "--allow-empty", "--cleanup=verbatim", "--author", agent + " via AILANG World <ailang-world@invalid>"}
	for _, part := range messageParts {
		args = append(args, "-m", part)
	}
	if err := runGit(ctx, git, args, nil); err != nil {
		return "", err
	}
	sha, err := Output(ctx, git, []string{"-C", dir, "rev-parse", "HEAD"})
	return strings.TrimSpace(sha), err
}
