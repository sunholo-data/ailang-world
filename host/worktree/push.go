package worktree

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Push sends exactly sha to branch without force, through the bounded git site.
func Push(ctx context.Context, git, dir, remote, sha, branch string) error {
	if remote == "" || strings.HasPrefix(remote, "-") || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(sha) ||
		branch == "" || strings.ContainsAny(branch, " :+~^?*[\\\r\n") || strings.HasPrefix(branch, "-") {
		return fmt.Errorf("worktree: invalid push destination")
	}
	return runGit(ctx, git, []string{"-C", dir, "-c", "core.hooksPath=/dev/null", "push", remote, sha + ":refs/heads/" + branch}, nil)
}
