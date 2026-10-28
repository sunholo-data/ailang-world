// Package worktree is World's one bounded, environment-scrubbed git site (row
// 138 M5, design_docs/planned/w-worldd-developer-cli.md §3.7). `ailang-worldd
// session new` uses it to make an episode's worktree before the session is
// minted (`git worktree add`), and to remove that worktree again when the mint
// fails (`git worktree remove`). `ailang-worldd surface-index` (row 188 M0a)
// uses Output for its read-only git calls (rev-parse, log, merge-base). All
// three launch through the one runGit.
//
// Why a package: the AC10 census (host/broker registry_publish_test.go)
// requires every production subprocess site to be DRIVEN and to scrub the
// registry credential, so the site lives in one small host package with an
// injectable git path the census can point at its env-dumping probe.
//
// The child environment is the process environment minus every registry
// variable (childenv.Scrubbed) and minus every GIT_* variable, so an ambient
// GIT_DIR, GIT_WORK_TREE, GIT_INDEX_FILE or GIT_CONFIG_* cannot redirect the
// git operation at a different repository (MUT-WORKTREE-ENV).
package worktree

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	// MaxStdout caps the standard output Output returns (8 MiB): surface-index's
	// largest read, a git log of at most 2500 merges, sits well under it.
	// Output past it is refused, never truncated.
	MaxStdout = 8 << 20
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
	return runGit(ctx, git, args, nil)
}

// Output runs `git args...` and returns its standard output. Stderr is kept
// only for the error message; stdout over MaxStdout is an error.
func Output(ctx context.Context, git string, args []string) (string, error) {
	out := &capped{max: MaxStdout}
	if err := runGit(ctx, git, args, out); err != nil {
		return "", err
	}
	if out.over {
		return "", fmt.Errorf("worktree: git %s: stdout is over the bound of %d bytes", strings.Join(args, " "), MaxStdout)
	}
	return out.buf.String(), nil
}

// IsLinked reports whether dir has the `.git` FILE a linked worktree carries
// ("gitdir: …"). It is the one test of "an episode is a git worktree": a
// directory failing it is refused by `session new` unless --allow-non-git,
// and is then reported degraded (row 185) by `session list` and /v1/health.
func IsLinked(dir string) bool {
	f, err := os.Open(filepath.Join(dir, ".git"))
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	head := make([]byte, 8)
	n, _ := io.ReadFull(f, head)
	return string(head[:n]) == "gitdir: "
}

// Remove runs `git -C repo worktree remove --force dir`.
func Remove(ctx context.Context, git, repo, dir string) error {
	return runGit(ctx, git, []string{"-C", repo, "worktree", "remove", "--force", dir}, nil)
}

// capped keeps the first max bytes written to it and notes any beyond, which
// it discards so the child never blocks on a full pipe.
type capped struct {
	buf  bytes.Buffer
	max  int
	over bool
}

func (c *capped) Write(p []byte) (int, error) {
	keep := p
	if room := c.max - c.buf.Len(); len(p) > room {
		c.over = true
		keep = p[:max(room, 0)]
	}
	c.buf.Write(keep)
	return len(p), nil
}

// runGit is the package's one git launch. With stdout nil, git's stdout and
// stderr both go to the error message (Add, Remove); otherwise stdout goes to
// stdout and only stderr to the message (Output).
func runGit(ctx context.Context, git string, args []string, stdout *capped) error {
	runCtx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, git, args...)
	cmd.Env = Env(os.Environ())
	cmd.Stdin = nil
	cmd.WaitDelay = waitDelay
	if deadline, ok := runCtx.Deadline(); ok {
		// Short callers (session metadata) must not spend another two seconds
		// draining pipes inherited by a descendant after git is killed.
		cmd.WaitDelay = min(waitDelay, max(time.Millisecond, time.Until(deadline)/20))
	}
	msg := &capped{max: maxOutput + 1}
	cmd.Stdout, cmd.Stderr = msg, msg
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if err := cmd.Run(); err != nil {
		text := strings.TrimSpace(msg.buf.String())
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
