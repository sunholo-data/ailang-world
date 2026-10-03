// Package worktree is World's one bounded, environment-scrubbed `git worktree
// add|remove` site (row 138 M5, design_docs/planned/w-worldd-developer-cli.md
// §3.7). `ailang-worldd session new` uses it to make an episode's worktree
// before the session is minted, and to remove that worktree again when the
// mint fails.
//
// Why a package: the AC10 census (host/broker registry_publish_test.go)
// requires every production subprocess site to be DRIVEN and to scrub the
// registry credential, so the site lives in one small host package with an
// injectable git path the census can point at its env-dumping probe.
//
// The child environment is the process environment minus every registry
// variable (childenv.Scrubbed) and minus every GIT_* variable, so an ambient
// GIT_DIR, GIT_WORK_TREE, GIT_INDEX_FILE or GIT_CONFIG_* cannot redirect the
// worktree operation at a different repository (MUT-WORKTREE-ENV).
package worktree

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// Timeout bounds one git invocation; WaitDelay bounds Wait after a kill.
const (
	Timeout   = 30 * time.Second
	waitDelay = 2 * time.Second
	// maxOutput caps the git output kept for an error message.
	maxOutput = 4096
)

// Env is the child environment: environ without registry or GIT_* variables.
// Like childenv.Scrubbed it is never nil, so exec never falls back to
// inheriting the process environment.
func Env(environ []string) []string {
	scrubbed := childenv.Scrubbed(environ)
	kept := make([]string, 0, len(scrubbed))
	for _, kv := range scrubbed {
		if strings.HasPrefix(kv, "GIT_") {
			continue
		}
		kept = append(kept, kv)
	}
	return kept
}

// Add runs `git -C repo worktree add --detach dir`, or with branch != ""
// `git -C repo worktree add -b branch dir`.
func Add(ctx context.Context, git, repo, dir, branch string) error {
	args := []string{"-C", repo, "worktree", "add"}
	if branch != "" {
		args = append(args, "-b", branch, dir)
	} else {
		args = append(args, "--detach", dir)
	}
	return runGit(ctx, git, args)
}

// Remove runs `git -C repo worktree remove --force dir`.
func Remove(ctx context.Context, git, repo, dir string) error {
	return runGit(ctx, git, []string{"-C", repo, "worktree", "remove", "--force", dir})
}

func runGit(ctx context.Context, git string, args []string) error {
	runCtx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, git, args...)
	cmd.Env = Env(os.Environ())
	cmd.Stdin = nil
	cmd.WaitDelay = waitDelay
	out, err := cmd.CombinedOutput()
	if err != nil {
		text := strings.TrimSpace(string(out))
		if len(text) > maxOutput {
			text = text[:maxOutput] + "…"
		}
		if runCtx.Err() != nil {
			return fmt.Errorf("worktree: git %s: timed out after %v: %w (%s)", strings.Join(args, " "), Timeout, runCtx.Err(), text)
		}
		return fmt.Errorf("worktree: git %s: %w: %s", strings.Join(args, " "), err, text)
	}
	return nil
}
