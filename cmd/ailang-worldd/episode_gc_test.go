package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
	base := t.TempDir()
	repo := filepath.Join(base, "repo")
	wsRoot := filepath.Join(base, "world-ws")
	cacheDir := filepath.Join(base, "exec-cache")
	for _, d := range []string{repo, wsRoot, cacheDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Initialize base repo with main branch.
	runGit(t, repo, "init", "-b", "main")
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "init")

	// 1. Clean merged stale worktree -> should be removed.
	cleanWT := filepath.Join(wsRoot, "ep-clean")
	runGit(t, repo, "worktree", "add", "-b", "world/ep-clean", cleanWT, "main")

	// 2. Dirty stale worktree -> must be skipped.
	dirtyWT := filepath.Join(wsRoot, "ep-dirty")
	runGit(t, repo, "worktree", "add", "-b", "world/ep-dirty", dirtyWT, "main")
	if err := os.WriteFile(filepath.Join(dirtyWT, "untracked.txt"), []byte("wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 3. Unmerged commit stale worktree -> must be skipped.
	unmergedWT := filepath.Join(wsRoot, "ep-unmerged")
	runGit(t, repo, "worktree", "add", "-b", "world/ep-unmerged", unmergedWT, "main")
	if err := os.WriteFile(filepath.Join(unmergedWT, "work.txt"), []byte("commit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, unmergedWT, "add", "work.txt")
	runGit(t, unmergedWT, "commit", "-m", "unmerged work")

	// Backdate all 3 worktrees by 100h so they pass the 72h threshold.
	oldTime := time.Now().Add(-100 * time.Hour)
	for _, d := range []string{cleanWT, dirtyWT, unmergedWT} {
		if err := os.Chtimes(d, oldTime, oldTime); err != nil {
			t.Fatal(err)
		}
	}

	// 4. Orphaned probe cache (old) and orphaned episode exec-cache + active episode cache.
	probeCache := filepath.Join(cacheDir, "wep-abcdef123456")
	orphanCache := filepath.Join(cacheDir, "ep-gone")
	cleanCache := filepath.Join(cacheDir, "ep-clean")
	dirtyCache := filepath.Join(cacheDir, "ep-dirty")
	for _, d := range []string{probeCache, orphanCache, cleanCache, dirtyCache} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "blob"), bytes.Repeat([]byte("x"), 1024), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(d, oldTime, oldTime); err != nil {
			t.Fatal(err)
		}
	}

	// Dry run leaves everything intact.
	var out, errw bytes.Buffer
	code := run([]string{"episode", "gc", "--workspace-root", wsRoot, "--exec-cache-dir", cacheDir, "--dry-run", "--json"}, &out, &errw)
	if code != exitOK {
		t.Fatalf("dry-run exit %d, stderr %q", code, errw.String())
	}
	var dryRep daemon.EpisodeGCReport
	if err := json.Unmarshal(out.Bytes(), &dryRep); err != nil {
		t.Fatalf("decode dry-run: %v\n%s", err, out.String())
	}
	if !dryRep.DryRun || len(dryRep.Removed) != 4 || len(dryRep.Skipped) != 2 {
		t.Fatalf("dry-run = %+v, want 4 removed (clean wt + 3 caches) and 2 skipped (dirty, unmerged)", dryRep)
	}
	if _, err := os.Stat(cleanWT); err != nil {
		t.Fatalf("dry-run deleted cleanWT: %v", err)
	}

	// Live run removes ep-clean + 3 caches, preserves ep-dirty, ep-unmerged, and dirtyCache.
	out.Reset()
	errw.Reset()
	code = run([]string{"episode", "gc", "--workspace-root", wsRoot, "--exec-cache-dir", cacheDir, "--json"}, &out, &errw)
	if code != exitOK {
		t.Fatalf("live exit %d, stderr %q", code, errw.String())
	}
	if _, err := os.Stat(cleanWT); !os.IsNotExist(err) {
		t.Errorf("cleanWT still exists after gc")
	}
	if _, err := os.Stat(dirtyWT); err != nil {
		t.Errorf("dirtyWT was deleted: %v", err)
	}
	if _, err := os.Stat(unmergedWT); err != nil {
		t.Errorf("unmergedWT was deleted: %v", err)
	}
	if _, err := os.Stat(probeCache); !os.IsNotExist(err) {
		t.Errorf("probeCache still exists")
	}
	if _, err := os.Stat(orphanCache); !os.IsNotExist(err) {
		t.Errorf("orphanCache still exists")
	}
	if _, err := os.Stat(cleanCache); !os.IsNotExist(err) {
		t.Errorf("cleanCache still exists")
	}
	if _, err := os.Stat(dirtyCache); err != nil {
		t.Errorf("dirtyCache for surviving worktree was deleted: %v", err)
	}
}
