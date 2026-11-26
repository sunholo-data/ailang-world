package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/worktree"
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
	// CredentialsPath is the store paired with WorkspaceRoot. An unavailable
	// lock or ledger bars worktree removal.
	CredentialsPath    string
	lookupOperatorStep func(context.Context, string, string) (bool, error)
	removeRegistered   func(context.Context, string, string, string, *os.File) error
}

const (
	autoGCInitialDelay = 2 * time.Minute
	autoGCInterval     = 6 * time.Hour
	autoGCOlderThan    = 72 * time.Hour
	autoGCPassBudget   = 2 * time.Minute
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
				d.runAutoGCPass(wsRoot, cacheDir, time.Now())
				timer.Reset(autoGCInterval)
			}
		}
	}()
}

func (d *Daemon) runAutoGCPass(wsRoot, cacheDir string, now time.Time) EpisodeGCReport {
	ctx, cancel := context.WithTimeout(context.Background(), autoGCPassBudget)
	defer cancel()
	rep := RunEpisodeGC(ctx, EpisodeGCOptions{
		WorkspaceRoot:   wsRoot,
		ExecCacheDir:    cacheDir,
		OlderThan:       autoGCOlderThan,
		Now:             now,
		IsActive:        func(episode string) bool { return d.isEpisodeActive(episode, now.Unix()) },
		CredentialsPath: d.credentialsPath,
	})
	if len(rep.Removed) > 0 && d.errLog != nil {
		fmt.Fprintf(d.errLog, "ailang-worldd: auto-gc reclaimed %s (%d removed, %d skipped)\n",
			FormatGCBytes(rep.ReclaimedBytes), len(rep.Removed), len(rep.Skipped))
	}
	return rep
}

func (d *Daemon) isEpisodeActive(episode string, now int64) bool {
	snap := d.activity.snapshot(now)
	for _, s := range snap.Sessions {
		if s.Episode == episode {
			return true
		}
	}
	return false
}

// RunEpisodeGC safely reclaims disk space from stale, clean, merged episode
// worktrees and orphaned startup probe / per-episode exec caches.
func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
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
				if reason := collectEpisodeWorktree(ctx, opts, now, name, wtPath, info); reason != "" {
					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
					})
					continue
				}
				removedWorktrees[name] = true
				rep.ReclaimedBytes += size
				rep.Removed = append(rep.Removed, EpisodeGCItem{
					Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: "clean & merged",
				})
			}
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

func collectEpisodeWorktree(ctx context.Context, opts EpisodeGCOptions, now time.Time, name, path string, scanned os.FileInfo) string {
	if opts.CredentialsPath == "" {
		return "operator step lookup unavailable"
	}
	if !verbs.EpisodeIDPattern.MatchString(name) {
		return "operator episode lock unavailable"
	}
	canonical, err := store.CanonicalPath(opts.CredentialsPath)
	if err != nil {
		return "operator episode lock unavailable"
	}
	lock, err := store.OpenOperatorEpisodeLock(canonical, name)
	if err != nil {
		return "operator episode lock unavailable"
	}
	defer lock.Close()
	lookup := opts.lookupOperatorStep
	if lookup == nil {
		lookup = store.EpisodeHasUnresolvedOperatorSteps
	}
	unresolved, err := lookup(ctx, canonical, name)
	if err != nil || ctx.Err() != nil {
		return "operator step lookup unavailable"
	}
	if unresolved {
		return "non-terminal operator step"
	}

	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() || !os.SameFile(info, scanned) {
		return "worktree changed during scan"
	}
	if opts.IsActive != nil && opts.IsActive(name) {
		return "active session"
	}
	if age := now.Sub(info.ModTime()); age < opts.OlderThan {
		return fmt.Sprintf("too recent (%s < %s)", age.Round(time.Minute), opts.OlderThan)
	}
	repo, registered, err := registeredGCWorktree(ctx, path)
	if err != nil {
		return err.Error()
	}
	if reason := checkWorktreeUnsafe(ctx, registered); reason != "" {
		return reason
	}
	fd, err := lock.Verify(canonical, name)
	if err != nil {
		return "operator episode lock unavailable"
	}
	if !opts.DryRun {
		remove := opts.removeRegistered
		if remove == nil {
			remove = worktree.RemoveRegistered
		}
		if err := remove(ctx, "git", repo, registered, fd); err != nil {
			return "remove failed: " + err.Error()
		}
	}
	return ""
}

func registeredGCWorktree(ctx context.Context, path string) (string, string, error) {
	none := errors.New("worktree registration unavailable or changed")
	meta := filepath.Join(path, ".git")
	info, err := os.Lstat(meta)
	if err != nil || !info.Mode().IsRegular() {
		return "", "", errors.New("not a linked git worktree")
	}
	physical, err := filepath.EvalSymlinks(meta)
	if err != nil {
		return "", "", none
	}
	gitDir, err := worktree.Output(ctx, "git", []string{"-C", path, "rev-parse", "--absolute-git-dir"})
	if err != nil {
		return "", "", none
	}
	back, err := os.ReadFile(filepath.Join(strings.TrimSuffix(gitDir, "\n"), "gitdir"))
	if err != nil {
		return "", "", none
	}
	backPath, err := filepath.EvalSymlinks(strings.TrimSuffix(string(back), "\n"))
	if err != nil || backPath != physical {
		return "", "", none
	}
	common, err := worktree.Output(ctx, "git", []string{"-C", path, "rev-parse", "--path-format=absolute", "--git-common-dir"})
	if err != nil {
		return "", "", none
	}
	repo, err := filepath.EvalSymlinks(strings.TrimSuffix(common, "\n"))
	if err != nil {
		return "", "", none
	}
	list, err := worktree.Output(ctx, "git", []string{"--git-dir", repo, "worktree", "list", "--porcelain", "-z"})
	if err != nil {
		return "", "", none
	}
	match := ""
	for _, record := range strings.Split(strings.TrimSuffix(list, "\x00\x00"), "\x00\x00") {
		first, _, _ := strings.Cut(record, "\x00")
		dir, ok := strings.CutPrefix(first, "worktree ")
		if !ok {
			return "", "", none
		}
		candidate, err := filepath.EvalSymlinks(dir)
		if err == nil && candidate == filepath.Dir(physical) {
			if match != "" {
				return "", "", none
			}
			match = dir
		}
	}
	if match == "" {
		return "", "", none
	}
	return repo, match, nil
}

func checkWorktreeUnsafe(ctx context.Context, wtPath string) string {
	out, err := worktree.Output(ctx, "git", []string{"--no-optional-locks", "-C", wtPath, "status", "--porcelain"})
	if err != nil {
		return "git status failed"
	}
	if strings.TrimSpace(out) != "" {
		return "dirty worktree (uncommitted or untracked files)"
	}
	for _, ref := range []string{"@{u}", "origin/dev", "origin/main", "dev", "main"} {
		if _, err := worktree.Output(ctx, "git", []string{"-C", wtPath, "merge-base", "--is-ancestor", "HEAD", ref}); err == nil {
			return ""
		}
	}
	return "unmerged commits on HEAD"
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
