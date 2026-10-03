package worktree

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func gitBin(t *testing.T) string {
	t.Helper()
	p, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not on PATH")
	}
	return p
}

// newRepo makes a git repo with one commit, configured locally so no global
// config is needed.
func newRepo(t *testing.T, git string) string {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "proj")
	for _, args := range [][]string{
		{"init", "-q", repo},
		{"-C", repo, "-c", "user.email=wt@example.invalid", "-c", "user.name=wt", "commit", "-q", "--allow-empty", "-m", "init"},
	} {
		cmd := exec.Command(git, args...)
		cmd.Env = Env(os.Environ())
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	return repo
}

func TestAddAndRemoveAWorktree(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	dir := filepath.Join(t.TempDir(), "ws", "ep1")
	_ = os.MkdirAll(filepath.Dir(dir), 0o755)
	if err := Add(testCtx(t), git, repo, dir, ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, ".git"))
	if err != nil || !strings.HasPrefix(string(b), "gitdir: ") {
		t.Fatalf("%s is not a linked worktree: %q %v", dir, b, err)
	}
	if err := Remove(testCtx(t), git, repo, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("worktree still present after Remove: %v", err)
	}
}

func TestAddWithBranch(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	dir := filepath.Join(t.TempDir(), "ep2")
	if err := Add(testCtx(t), git, repo, dir, "ep2-work"); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(git, "-C", dir, "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Env = Env(os.Environ())
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) != "ep2-work" {
		t.Fatalf("branch = %q (%v), want ep2-work", out, err)
	}
}

// AC5.8 (with the AC10 driver in host/broker): an ambient GIT_DIR must not
// redirect the worktree add at another repository, and no GIT_* or registry
// variable reaches the child. Kills MUT-WORKTREE-ENV.
func TestAddIgnoresAmbientGitAndRegistryEnvironment(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	decoy := newRepo(t, git)
	t.Setenv("GIT_DIR", filepath.Join(decoy, ".git"))
	t.Setenv("GIT_WORK_TREE", decoy)
	t.Setenv("AILANG_REGISTRY_API_KEY", "sentinel")
	dir := filepath.Join(t.TempDir(), "ep3")
	if err := Add(testCtx(t), git, repo, dir, ""); err != nil {
		t.Fatalf("Add under an ambient GIT_DIR: %v", err)
	}
	list := exec.Command(git, "-C", repo, "worktree", "list", "--porcelain")
	list.Env = Env(os.Environ())
	out, err := list.Output()
	if err != nil {
		t.Fatal(err)
	}
	real, _ := filepath.EvalSymlinks(dir)
	if !strings.Contains(string(out), "worktree "+real) && !strings.Contains(string(out), "worktree "+dir) {
		t.Fatalf("the worktree was not added to the named repo:\n%s", out)
	}
	for _, kv := range Env(os.Environ()) {
		if strings.HasPrefix(kv, "GIT_") || strings.HasPrefix(kv, "AILANG_REGISTRY") {
			t.Fatalf("Env kept %q", strings.SplitN(kv, "=", 2)[0])
		}
	}
	if Env(nil) == nil {
		t.Fatal("Env(nil) is nil: exec would inherit the process environment")
	}
}

func TestAddFailureIsReported(t *testing.T) {
	git := gitBin(t)
	err := Add(testCtx(t), git, filepath.Join(t.TempDir(), "not-a-repo"), filepath.Join(t.TempDir(), "x"), "")
	if err == nil || !strings.Contains(err.Error(), "worktree add") {
		t.Fatalf("err = %v, want a reported git failure", err)
	}
}
