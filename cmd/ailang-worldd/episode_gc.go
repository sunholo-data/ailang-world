package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const episodeGCHelp = `usage: ailang-worldd episode gc [--workspace-root D] [--exec-cache-dir D]
    [--older-than 72h] [--dry-run] [--json]

Safely reclaim disk space from finished/stale episode worktrees and orphaned
exec caches. Never removes dirty worktrees (uncommitted or untracked files)
or worktrees with unmerged commits.

  --workspace-root D   episode worktrees root (default ~/.ailang/world-ws)
  --exec-cache-dir D   per-episode exec caches (default ~/.ailang/world/exec-cache)
  --older-than DUR     minimum age before a clean merged worktree is removed (default 72h)
  --dry-run            report what would be removed without deleting anything
  --json               emit JSON report
`

type gcItem struct {
	Kind   string `json:"kind"` // "probe-cache" | "exec-cache" | "worktree"
	Name   string `json:"name"`
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	Reason string `json:"reason,omitempty"`
}

type gcReport struct {
	DryRun         bool     `json:"dry_run"`
	OlderThan      string   `json:"older_than"`
	ReclaimedBytes int64    `json:"reclaimed_bytes"`
	Removed        []gcItem `json:"removed"`
	Skipped        []gcItem `json:"skipped"`
}

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

	rep := executeEpisodeGC(*wsRoot, *cacheDir, *olderThan, *dryRun, time.Now())
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
		verb, formatBytes(rep.ReclaimedBytes), len(rep.Removed), len(rep.Skipped))
	for _, it := range rep.Removed {
		fmt.Fprintf(stdout, "  REMOVE %-11s %-28s %8s  (%s)\n", it.Kind, it.Name, formatBytes(it.Bytes), it.Reason)
	}
	for _, it := range rep.Skipped {
		fmt.Fprintf(stdout, "  SKIP   %-11s %-28s %8s  (%s)\n", it.Kind, it.Name, formatBytes(it.Bytes), it.Reason)
	}
	return exitOK
}

func executeEpisodeGC(wsRoot, cacheDir string, olderThan time.Duration, dryRun bool, now time.Time) gcReport {
	rep := gcReport{
		DryRun:    dryRun,
		OlderThan: olderThan.String(),
		Removed:   []gcItem{},
		Skipped:   []gcItem{},
	}
	removedWorktrees := map[string]bool{}
	reposToPrune := map[string]bool{}

	// 1. Scan episode worktrees in wsRoot.
	if wsRoot != "" {
		if ents, err := os.ReadDir(wsRoot); err == nil {
			for _, e := range ents {
				if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
					continue
				}
				name := e.Name()
				wtPath := filepath.Join(wsRoot, name)
				info, err := e.Info()
				if err != nil {
					continue
				}
				age := now.Sub(info.ModTime())
				size := dirSizeBytes(wtPath)
				if age < olderThan {
					rep.Skipped = append(rep.Skipped, gcItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size,
						Reason: fmt.Sprintf("too recent (%s < %s)", age.Round(time.Minute), olderThan),
					})
					continue
				}
				if reason := checkWorktreeUnsafe(wtPath); reason != "" {
					rep.Skipped = append(rep.Skipped, gcItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
					})
					continue
				}
				commonDir := gitCommonDir(wtPath)
				if !dryRun {
					if err := removeWorktree(wtPath, commonDir); err != nil {
						rep.Skipped = append(rep.Skipped, gcItem{
							Kind: "worktree", Name: name, Path: wtPath, Bytes: size,
							Reason: "remove failed: " + err.Error(),
						})
						continue
					}
					if commonDir != "" {
						reposToPrune[commonDir] = true
					}
				}
				removedWorktrees[name] = true
				rep.ReclaimedBytes += size
				rep.Removed = append(rep.Removed, gcItem{
					Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: "clean & merged",
				})
			}
		}
	}
	if !dryRun {
		for commonDir := range reposToPrune {
			_ = exec.Command("git", "--git-dir", commonDir, "worktree", "prune").Run()
		}
	}

	// 2. Scan exec-cache dir for orphaned startup probe caches (wep-*) and orphaned episode caches.
	if cacheDir != "" {
		if ents, err := os.ReadDir(cacheDir); err == nil {
			for _, e := range ents {
				if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
					continue
				}
				name := e.Name()
				full := filepath.Join(cacheDir, name)
				info, err := e.Info()
				if err != nil {
					continue
				}
				size := dirSizeBytes(full)
				if strings.HasPrefix(name, "wep-") {
					if now.Sub(info.ModTime()) < time.Hour {
						rep.Skipped = append(rep.Skipped, gcItem{
							Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: "recent probe (<1h)",
						})
						continue
					}
					if !dryRun {
						if err := os.RemoveAll(full); err != nil {
							rep.Skipped = append(rep.Skipped, gcItem{
								Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: err.Error(),
							})
							continue
						}
					}
					rep.ReclaimedBytes += size
					rep.Removed = append(rep.Removed, gcItem{
						Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: "orphaned startup probe cache",
					})
					continue
				}
				wtExists := false
				if wsRoot != "" && !removedWorktrees[name] {
					if st, err := os.Stat(filepath.Join(wsRoot, name)); err == nil && st.IsDir() {
						wtExists = true
					}
				}
				if wtExists {
					continue
				}
				if !dryRun {
					if err := os.RemoveAll(full); err != nil {
						rep.Skipped = append(rep.Skipped, gcItem{
							Kind: "exec-cache", Name: name, Path: full, Bytes: size, Reason: err.Error(),
						})
						continue
					}
				}
				rep.ReclaimedBytes += size
				rep.Removed = append(rep.Removed, gcItem{
					Kind: "exec-cache", Name: name, Path: full, Bytes: size, Reason: "worktree absent",
				})
			}
		}
	}

	sort.Slice(rep.Removed, func(i, j int) bool { return rep.Removed[i].Path < rep.Removed[j].Path })
	sort.Slice(rep.Skipped, func(i, j int) bool { return rep.Skipped[i].Path < rep.Skipped[j].Path })
	return rep
}

// checkWorktreeUnsafe returns a non-empty skip reason if wtPath is not a clean,
// merged git worktree.
func checkWorktreeUnsafe(wtPath string) string {
	// Must be a git worktree (has .git file, never a standalone repo .git dir).
	gitMeta, err := os.Lstat(filepath.Join(wtPath, ".git"))
	if err != nil || gitMeta.IsDir() {
		return "not a linked git worktree"
	}
	out, err := exec.Command("git", "-C", wtPath, "status", "--porcelain").Output()
	if err != nil {
		return "git status failed"
	}
	if strings.TrimSpace(string(out)) != "" {
		return "dirty worktree (uncommitted or untracked files)"
	}
	// Check whether HEAD is merged into any local/remote base branch (dev/main) or its upstream.
	for _, ref := range []string{"@{u}", "origin/dev", "origin/main", "dev", "main"} {
		if exec.Command("git", "-C", wtPath, "merge-base", "--is-ancestor", "HEAD", ref).Run() == nil {
			return ""
		}
	}
	return "unmerged commits on HEAD"
}

func gitCommonDir(wtPath string) string {
	out, err := exec.Command("git", "-C", wtPath, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func removeWorktree(wtPath, commonDir string) error {
	if commonDir != "" {
		if err := exec.Command("git", "--git-dir", commonDir, "worktree", "remove", wtPath).Run(); err == nil {
			return nil
		}
	}
	return os.RemoveAll(wtPath)
}

func dirSizeBytes(root string) int64 {
	var total int64
	_ = filepath.WalkDir(root, func(_ string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}

func formatBytes(n int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)
	switch {
	case n >= gb:
		return fmt.Sprintf("%.1f GiB", float64(n)/float64(gb))
	case n >= mb:
		return fmt.Sprintf("%.1f MiB", float64(n)/float64(mb))
	case n >= kb:
		return fmt.Sprintf("%.1f KiB", float64(n)/float64(kb))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
