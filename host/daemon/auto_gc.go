package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// EpisodeGCItem describes one worktree or cache directory removed or skipped by RunEpisodeGC.
type EpisodeGCItem struct {
	Kind   string `json:"kind"` // "probe-cache" | "exec-cache" | "worktree"
	Name   string `json:"name"`
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	Reason string `json:"reason,omitempty"`
}

// EpisodeGCReport is the structured result of one GC pass.
type EpisodeGCReport struct {
	DryRun         bool            `json:"dry_run"`
	OlderThan      string          `json:"older_than"`
	ReclaimedBytes int64           `json:"reclaimed_bytes"`
	Removed        []EpisodeGCItem `json:"removed"`
	Skipped        []EpisodeGCItem `json:"skipped"`
}

// EpisodeGCOptions configures a fail-safe disk claw-back pass.
type EpisodeGCOptions struct {
	WorkspaceRoot string
	ExecCacheDir  string
	OlderThan     time.Duration
	DryRun        bool
	Now           time.Time
	IsActive      func(episode string) bool
}

const (
	autoGCInitialDelay = 2 * time.Minute
	autoGCInterval     = 6 * time.Hour
	autoGCOlderThan    = 72 * time.Hour
)

// startAutoGC runs periodic fail-safe disk claw-back inside the daemon when
// WorkspaceRoot is configured, stopping cleanly when d.aguiStop closes.
func (d *Daemon) startAutoGC(dbPath string) {
	if d.workspace.root == "" || dbPath == "" {
		return
	}
	wsRoot := d.workspace.root
	cacheDir := filepath.Join(filepath.Dir(dbPath), "exec-cache")
	go func() {
		timer := time.NewTimer(autoGCInitialDelay)
		defer timer.Stop()
		for {
			select {
			case <-d.aguiStop:
				return
			case <-timer.C:
				rep := RunEpisodeGC(EpisodeGCOptions{
					WorkspaceRoot: wsRoot,
					ExecCacheDir:  cacheDir,
					OlderThan:     autoGCOlderThan,
					Now:           time.Now(),
					IsActive:      d.isEpisodeActive,
				})
				if len(rep.Removed) > 0 && d.errLog != nil {
					fmt.Fprintf(d.errLog, "ailang-worldd: auto-gc reclaimed %s (%d removed, %d skipped)\n",
						FormatGCBytes(rep.ReclaimedBytes), len(rep.Removed), len(rep.Skipped))
				}
				timer.Reset(autoGCInterval)
			}
		}
	}()
}

func (d *Daemon) isEpisodeActive(episode string) bool {
	snap := d.activity.snapshot(time.Now().Unix())
	for _, s := range snap.Sessions {
		if s.Episode == episode {
			return true
		}
	}
	return false
}

// RunEpisodeGC safely reclaims disk space from stale, clean, merged episode
// worktrees and orphaned startup probe / per-episode exec caches.
func RunEpisodeGC(opts EpisodeGCOptions) EpisodeGCReport {
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	rep := EpisodeGCReport{
		DryRun:    opts.DryRun,
		OlderThan: opts.OlderThan.String(),
		Removed:   []EpisodeGCItem{},
		Skipped:   []EpisodeGCItem{},
	}
	removedWorktrees := map[string]bool{}
	reposToPrune := map[string]bool{}

	// 1. Scan episode worktrees in WorkspaceRoot.
	if opts.WorkspaceRoot != "" {
		if ents, err := os.ReadDir(opts.WorkspaceRoot); err == nil {
			for _, e := range ents {
				if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
					continue
				}
				name := e.Name()
				wtPath := filepath.Join(opts.WorkspaceRoot, name)
				info, err := e.Info()
				if err != nil {
					continue
				}
				size := gcDirSizeBytes(wtPath)
				if opts.IsActive != nil && opts.IsActive(name) {
					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: "active session",
					})
					continue
				}
				age := now.Sub(info.ModTime())
				if age < opts.OlderThan {
					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size,
						Reason: fmt.Sprintf("too recent (%s < %s)", age.Round(time.Minute), opts.OlderThan),
					})
					continue
				}
				if reason := checkWorktreeUnsafe(wtPath); reason != "" {
					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
					})
					continue
				}
				commonDir := gitCommonDir(wtPath)
				if !opts.DryRun {
					if err := removeWorktree(wtPath, commonDir); err != nil {
						rep.Skipped = append(rep.Skipped, EpisodeGCItem{
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
				rep.Removed = append(rep.Removed, EpisodeGCItem{
					Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: "clean & merged",
				})
			}
		}
	}
	if !opts.DryRun {
		for commonDir := range reposToPrune {
			_ = exec.Command("git", "--git-dir", commonDir, "worktree", "prune").Run()
		}
	}

	// 2. Scan exec-cache dir for orphaned startup probe caches (wep-*) and orphaned episode caches.
	if opts.ExecCacheDir != "" {
		if ents, err := os.ReadDir(opts.ExecCacheDir); err == nil {
			for _, e := range ents {
				if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
					continue
				}
				name := e.Name()
				full := filepath.Join(opts.ExecCacheDir, name)
				info, err := e.Info()
				if err != nil {
					continue
				}
				size := gcDirSizeBytes(full)
				if strings.HasPrefix(name, "wep-") {
					if now.Sub(info.ModTime()) < time.Hour {
						rep.Skipped = append(rep.Skipped, EpisodeGCItem{
							Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: "recent probe (<1h)",
						})
						continue
					}
					if !opts.DryRun {
						if err := os.RemoveAll(full); err != nil {
							rep.Skipped = append(rep.Skipped, EpisodeGCItem{
								Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: err.Error(),
							})
							continue
						}
					}
					rep.ReclaimedBytes += size
					rep.Removed = append(rep.Removed, EpisodeGCItem{
						Kind: "probe-cache", Name: name, Path: full, Bytes: size, Reason: "orphaned startup probe cache",
					})
					continue
				}
				if opts.IsActive != nil && opts.IsActive(name) {
					continue
				}
				wtExists := false
				if opts.WorkspaceRoot != "" && !removedWorktrees[name] {
					if st, err := os.Stat(filepath.Join(opts.WorkspaceRoot, name)); err == nil && st.IsDir() {
						wtExists = true
					}
				}
				if wtExists {
					continue
				}
				if !opts.DryRun {
					if err := os.RemoveAll(full); err != nil {
						rep.Skipped = append(rep.Skipped, EpisodeGCItem{
							Kind: "exec-cache", Name: name, Path: full, Bytes: size, Reason: err.Error(),
						})
						continue
					}
				}
				rep.ReclaimedBytes += size
				rep.Removed = append(rep.Removed, EpisodeGCItem{
					Kind: "exec-cache", Name: name, Path: full, Bytes: size, Reason: "worktree absent",
				})
			}
		}
	}

	sort.Slice(rep.Removed, func(i, j int) bool { return rep.Removed[i].Path < rep.Removed[j].Path })
	sort.Slice(rep.Skipped, func(i, j int) bool { return rep.Skipped[i].Path < rep.Skipped[j].Path })
	return rep
}

func checkWorktreeUnsafe(wtPath string) string {
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

func gcDirSizeBytes(root string) int64 {
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

// FormatGCBytes formats a byte count for human-readable GC reports.
func FormatGCBytes(n int64) string {
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
