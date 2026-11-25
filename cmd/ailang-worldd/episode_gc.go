package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

const episodeGCHelp = `usage: ailang-worldd episode gc [--workspace-root D] [--exec-cache-dir D]
    [--older-than 72h] [--dry-run] [--json]

Safely reclaim disk space from finished/stale episode worktrees and orphaned
exec caches. Never removes dirty worktrees (uncommitted or untracked files)
or worktrees with unmerged commits. Also runs automatically inside ailang-worldd
serve (2m after boot, then every 6h).

  --workspace-root D   episode worktrees root (default ~/.ailang/world-ws)
  --exec-cache-dir D   per-episode exec caches (default ~/.ailang/world/exec-cache)
  --older-than DUR     minimum age before a clean merged worktree is removed (default 72h)
  --dry-run            report what would be removed without deleting anything
  --json               emit JSON report
`

func defaultWorkspaceRootDir() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, ".ailang", "world-ws")
	}
	return ""
}

func defaultExecCacheDir() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, ".ailang", "world", "exec-cache")
	}
	return ""
}

func runEpisodeGC(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd episode gc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	wsRoot := fs.String("workspace-root", defaultWorkspaceRootDir(), "episode worktrees directory")
	cacheDir := fs.String("exec-cache-dir", defaultExecCacheDir(), "per-episode exec cache directory")
	olderThan := fs.Duration("older-than", 72*time.Hour, "minimum age before removing clean merged worktrees")
	dryRun := fs.Bool("dry-run", false, "preview without deleting")
	asJSON := fs.Bool("json", false, "emit JSON")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, episodeGCHelp)
			return exitOK
		}
		fmt.Fprint(stderr, episodeGCHelp)
		return exitUsage
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd episode gc: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	if *olderThan < 0 {
		fmt.Fprintln(stderr, "ailang-worldd episode gc: --older-than must be non-negative")
		return exitUsage
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	rep := daemon.RunEpisodeGC(ctx, daemon.EpisodeGCOptions{
		WorkspaceRoot: *wsRoot,
		ExecCacheDir:  *cacheDir,
		OlderThan:     *olderThan,
		DryRun:        *dryRun,
		Now:           time.Now(),
	})
	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
		return exitOK
	}
	verb := "reclaimed"
	if rep.DryRun {
		verb = "would reclaim"
	}
	fmt.Fprintf(stdout, "episode gc: %s %s (%d removed, %d skipped)\n",
		verb, daemon.FormatGCBytes(rep.ReclaimedBytes), len(rep.Removed), len(rep.Skipped))
	for _, it := range rep.Removed {
		fmt.Fprintf(stdout, "  REMOVE %-11s %-28s %8s  (%s)\n", it.Kind, it.Name, daemon.FormatGCBytes(it.Bytes), it.Reason)
	}
	for _, it := range rep.Skipped {
		fmt.Fprintf(stdout, "  SKIP   %-11s %-28s %8s  (%s)\n", it.Kind, it.Name, daemon.FormatGCBytes(it.Bytes), it.Reason)
	}
	return exitOK
}
