package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.CommandContext(boundedTestContext(t), "git", append([]string{"-C", dir}, args...)...)
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
	home := filepath.Join(base, "home")
	t.Setenv("HOME", home)
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
	oldTime := time.Unix(1_000_000_000, 0)
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

	// An unrelated default store must never clear a custom workspace.
	defaultCreds := filepath.Join(home, ".ailang", "world", store.CredentialsFileName)
	if err := os.MkdirAll(filepath.Dir(defaultCreds), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := store.UpgradeCredentials(defaultCreds); err != nil {
		t.Fatal(err)
	}
	creds := filepath.Join(base, "custom.db")
	c, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ctx := boundedTestContext(t)
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-00000000000000000000000000000001",
		Kind: "start", Episode: "ep-clean", Store: creds, Task: "fixture", Review: "fixture",
		Options:    approval.OperatorOptions{Preset: "sol", Headless: true},
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
	a.Store, err = store.CanonicalPath(creds)
	if err != nil {
		t.Fatal(err)
	}
	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	if _, err := c.RequestOperatorAction(ctx, store.OperatorRequest{Action: a, OperationID: "00000000000000000000000000000001",
		CreatedAt: 1000, ExpiresAt: 1600}, 1000); err != nil {
		t.Fatal(err)
	}
	// A reserved effect can survive its parent action. Use the real migrated
	// schema to model that case without claiming or running an action.
	raw, err := sql.Open("sqlite", creds)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	if _, err := raw.ExecContext(ctx, `UPDATE operator_actions SET state='failed', terminal_at=1001 WHERE id=?;
		INSERT INTO operator_action_steps (kind, identity, action_id, phase, version, owner, created_at)
		VALUES ('worktree','fixture-wt',?,'reserved',1,'fixture',1001)`, a.ID, a.ID); err != nil {
		t.Fatal(err)
	}
	gitState := func(args ...string) []byte {
		t.Helper()
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = cleanWT
		got, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, got)
		}
		return got
	}
	head, status := gitState("rev-parse", "HEAD"), gitState("--no-optional-locks", "status", "--porcelain=v1")
	refs, registration := gitState("show-ref"), gitState("worktree", "list", "--porcelain")
	physical := func(path string) string {
		t.Helper()
		got, err := filepath.EvalSymlinks(path)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	record := homeServicePaths(home)
	if err := os.MkdirAll(filepath.Dir(record.Record), 0700); err != nil {
		t.Fatal(err)
	}
	setRecord := func(root string) {
		t.Helper()
		data, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel,
			Address: "127.0.0.1:7644", Plist: record.Plist, Binary: record.Binary, Credentials: creds, Workspace: root})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(record.Record, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	assertSkip := func(reason string, extra ...string) {
		t.Helper()
		var out, errw bytes.Buffer
		args := []string{"episode", "gc", "--workspace-root", wsRoot, "--exec-cache-dir", filepath.Join(base, "absent-cache"), "--json"}
		code := run(append(args, extra...), &out, &errw)
		var rep daemon.EpisodeGCReport
		if err := json.Unmarshal(out.Bytes(), &rep); code != exitOK || err != nil || len(rep.Removed) != 0 {
			t.Fatalf("GC refusal: code %d, report %+v, stderr %q, decode %v", code, rep, errw.String(), err)
		}
		found := false
		for _, item := range rep.Skipped {
			if item.Name == "ep-clean" && item.Reason == reason && physical(item.Path) == physical(cleanWT) {
				found = true
			}
		}
		if _, err := os.Stat(cleanWT); err != nil || !found || !bytes.Equal(head, gitState("rev-parse", "HEAD")) ||
			!bytes.Equal(status, gitState("--no-optional-locks", "status", "--porcelain=v1")) ||
			!bytes.Equal(refs, gitState("show-ref")) || !bytes.Equal(registration, gitState("worktree", "list", "--porcelain")) ||
			string(readStartFixtureFile(t, filepath.Join(cleanWT, "README.md"))) != "hello\n" {
			t.Fatalf("want unchanged ep-clean and skip %q: %+v, stat %v", reason, rep, err)
		}
	}
	assertSkip("operator step lookup unavailable")
	setRecord(repo)
	if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
		t.Fatalf("different configured workspace selected credentials %q", got)
	}
	assertSkip("operator step lookup unavailable", "--credentials-db", creds)
	alias := filepath.Join(base, "workspace-alias")
	if err := os.Symlink(wsRoot, alias); err != nil {
		t.Fatal(err)
	}
	setRecord(wsRoot)
	assertSkip("operator step lookup unavailable", "--credentials-db", defaultCreds)
	held, err := store.OpenOperatorEpisodeLock(creds, "ep-clean")
	if err != nil {
		t.Fatal(err)
	}
	assertSkip("operator episode lock unavailable", "--credentials-db", creds)
	if err := held.Close(); err != nil {
		t.Fatal(err)
	}
	assertSkip("non-terminal operator step", "--credentials-db", creds)
	if got := episodeGCCredentialsPath(alias, ""); got != creds {
		t.Fatalf("symlink scan root selected %q, want %q", got, creds)
	}
	setRecord(alias)
	assertSkip("non-terminal operator step")
	if _, err := raw.ExecContext(ctx, "UPDATE operator_action_steps SET phase='observed',receipt=? WHERE kind='worktree'", []byte("stopped: retained")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(record.Record); err != nil {
		t.Fatal(err)
	}
	if codePath := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); codePath != "" {
		t.Fatalf("missing workspace selected a credentials path: %q", codePath)
	}
	if err := os.MkdirAll(defaultWorkspaceRootDir(), 0700); err != nil {
		t.Fatal(err)
	}
	if got := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); got != defaultCreds {
		t.Fatalf("default credentials path = %q, want %q", got, defaultCreds)
	}
	if err := os.WriteFile(record.Record, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); got != "" {
		t.Fatalf("malformed record fell back to credentials %q", got)
	}
	t.Run("missing_home", func(t *testing.T) {
		t.Setenv("HOME", "")
		t.Chdir(t.TempDir())
		p := homeServicePaths("")
		if err := os.MkdirAll(filepath.Dir(p.Record), 0700); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel,
			Address: "127.0.0.1:7644", Plist: p.Plist, Binary: p.Binary, Credentials: creds, Workspace: wsRoot})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p.Record, data, 0600); err != nil {
			t.Fatal(err)
		}
		if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
			t.Fatalf("missing HOME selected credentials from cwd: %q", got)
		}
	})
	setRecord(alias)
	if got := episodeGCCredentialsPath(filepath.Join(base, "missing"), ""); got != "" {
		t.Fatalf("unresolvable scan root selected credentials %q", got)
	}
	setRecord(filepath.Join(base, "missing"))
	if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
		t.Fatalf("unresolvable recorded root selected credentials %q", got)
	}
	setRecord(alias)
	t.Run("relative_root", func(t *testing.T) {
		t.Chdir(base)
		if got := episodeGCCredentialsPath(filepath.Base(wsRoot), ""); got != creds {
			t.Fatalf("relative scan selected %q, want %q", got, creds)
		}
	})

	// Dry run leaves everything intact and a terminal-only ledger allows GC.
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
	if _, err := os.Stat(cleanWT); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("cleanWT still exists after gc")
	}
	if _, err := os.Stat(dirtyWT); err != nil {
		t.Errorf("dirtyWT was deleted: %v", err)
	}
	if _, err := os.Stat(unmergedWT); err != nil {
		t.Errorf("unmergedWT was deleted: %v", err)
	}
	if _, err := os.Stat(probeCache); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("probeCache still exists")
	}
	if _, err := os.Stat(orphanCache); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("orphanCache still exists")
	}
	if _, err := os.Stat(cleanCache); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("cleanCache still exists")
	}
	if _, err := os.Stat(dirtyCache); err != nil {
		t.Errorf("dirtyCache for surviving worktree was deleted: %v", err)
	}
}
