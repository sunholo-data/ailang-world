package daemon

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"syscall"
)

func TestEpisodeGCOperatorSteps(t *testing.T) {
	const contenderEnv = "WORLD_TEST_GC_CONTENDER"
	if creds := os.Getenv(contenderEnv); creds != "" {
		contender, err := store.OpenOperatorEpisodeLock(creds, "ep")
		if contender != nil {
			contender.Close()
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("activation acquired at the GC removal boundary: %v", err)
		}
		return
	}
	root := t.TempDir()
	ws := filepath.Join(root, "workspaces")
	if err := os.Mkdir(ws, 0700); err != nil {
		t.Fatal(err)
	}
	wt := filepath.Join(ws, "ep")
	git := func(dir string, args ...string) []byte {
		t.Helper()
		cmd := exec.CommandContext(boundedTestContext(t), "git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.net",
			"GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.net")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
		return out
	}
	git(root, "init", "-b", "main")
	if err := os.WriteFile(filepath.Join(root, "tracked"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	git(root, "add", "tracked")
	git(root, "commit", "-m", "base")
	git(root, "worktree", "add", "-b", "world/ep", wt, "main")
	info, err := os.Stat(wt)
	if err != nil {
		t.Fatal(err)
	}
	physical := func(path string) string {
		t.Helper()
		got, err := filepath.EvalSymlinks(path)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	path := physical(wt)
	head, status := git(wt, "rev-parse", "HEAD"), git(wt, "--no-optional-locks", "status", "--porcelain=v1")
	refs, registrations := git(root, "show-ref"), git(root, "worktree", "list", "--porcelain")
	contents := func() []byte {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(wt, "tracked"))
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	originalContents := contents()
	if len(status) != 0 {
		t.Fatalf("fixture is dirty: %q", status)
	}
	creds := filepath.Join(root, "credentials.db")
	if _, err := store.UpgradeCredentials(creds); err != nil {
		t.Fatal(err)
	}
	opts := EpisodeGCOptions{WorkspaceRoot: ws, CredentialsPath: creds, Now: info.ModTime().Add(73 * time.Hour), OlderThan: 72 * time.Hour}
	checkSkip := func(t *testing.T, rep EpisodeGCReport, reason string) {
		t.Helper()
		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Name != "ep" ||
			rep.Skipped[0].Reason != reason || physical(rep.Skipped[0].Path) != path {
			t.Fatalf("want sole worktree skip %q: %+v", reason, rep)
		}
		if _, err := os.Stat(wt); err != nil || !bytes.Equal(head, git(wt, "rev-parse", "HEAD")) ||
			!bytes.Equal(status, git(wt, "--no-optional-locks", "status", "--porcelain=v1")) ||
			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(registrations, git(root, "worktree", "list", "--porcelain")) ||
			!bytes.Equal(originalContents, contents()) {
			t.Fatalf("refusal changed worktree or Git evidence: %v", err)
		}
	}
	t.Run("candidate", func(t *testing.T) {
		o := opts
		o.DryRun = true
		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) { return false, nil }
		index := filepath.Join(string(bytes.TrimSpace(git(wt, "rev-parse", "--absolute-git-dir"))), "index")
		before, err := os.ReadFile(index)
		if err != nil {
			t.Fatal(err)
		}
		tracked := filepath.Join(wt, "tracked")
		if err := os.Chtimes(tracked, o.Now, o.Now); err != nil {
			t.Fatal(err)
		}
		rep := RunEpisodeGC(boundedTestContext(t), o)
		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || rep.Removed[0].Reason != "clean & merged" ||
			physical(rep.Removed[0].Path) != path || !bytes.Equal(refs, git(root, "show-ref")) ||
			!bytes.Equal(registrations, git(root, "worktree", "list", "--porcelain")) || !bytes.Equal(originalContents, contents()) {
			t.Fatalf("successful negative lookup must allow the candidate: %+v", rep)
		}
		after, err := os.ReadFile(index)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatalf("GC eligibility refreshed the Git index: %v", err)
		}
	})
	for _, name := range []string{"protected", "error", "timeout", "canceled", "unavailable"} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(boundedTestContext(t))
			defer cancel()
			o := opts
			reason := "operator step lookup unavailable"
			switch name {
			case "protected":
				reason = "non-terminal operator step"
				o.lookupOperatorStep = func(_ context.Context, _, episode string) (bool, error) {
					return episode == "ep", nil
				}
			case "error":
				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
					return false, errors.New("fixture store error")
				}
			case "timeout":
				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
					return false, context.DeadlineExceeded
				}
			case "unavailable":
				o.CredentialsPath = ""
			case "canceled":
				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
					cancel()
					return false, nil
				}
			}
			checkSkip(t, RunEpisodeGC(ctx, o), reason)
		})
		if t.Failed() {
			return
		}
	}
	t.Run("active", func(t *testing.T) {
		d := &Daemon{credentialsPath: creds}
		d.activity.begin("ep", opts.Now.Unix()+1, opts.Now.Unix())()
		checkSkip(t, d.runAutoGCPass(ws, "", opts.Now), "active session")
	})
	t.Run("bad_store_path", func(t *testing.T) {
		o := opts
		o.CredentialsPath = filepath.Join(root, "missing", "credentials.db")
		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "operator episode lock unavailable")
	})
	t.Run("invalid_episode", func(t *testing.T) {
		before := git(root, "worktree", "list", "--porcelain")
		reason := collectEpisodeWorktree(boundedTestContext(t), opts, opts.Now, "INVALID", wt, info)
		if reason != "operator episode lock unavailable" || !bytes.Equal(before, git(root, "worktree", "list", "--porcelain")) ||
			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(originalContents, contents()) {
			t.Fatalf("malformed episode modified Git or returned wrong reason: %q", reason)
		}
	})
	t.Run("lock_changed_under_review", func(t *testing.T) {
		o := opts
		canonical, err := store.CanonicalPath(creds)
		if err != nil {
			t.Fatal(err)
		}
		path := store.OperatorEpisodeLockPath(canonical, "ep")
		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
			if err := os.Rename(path, path+".old"); err != nil {
				t.Fatal(err)
			}
			return false, os.WriteFile(path, nil, 0600)
		}
		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "operator episode lock unavailable")
	})
	t.Run("replaced_during_scan", func(t *testing.T) {
		o := opts
		alternative := filepath.Join(root, "moved")
		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
			if err := os.Rename(wt, alternative); err != nil {
				t.Fatal(err)
			}
			return false, os.Mkdir(wt, 0700)
		}
		rep := RunEpisodeGC(boundedTestContext(t), o)
		if len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree changed during scan" ||
			len(rep.Removed) != 0 || !bytes.Equal(refs, git(root, "show-ref")) {
			t.Fatalf("reused pathname was removed or returned wrong reason: %+v", rep)
		}
		if err := os.Remove(wt); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(alternative, wt); err != nil {
			t.Fatal(err)
		}
		checkSkip(t, rep, "worktree changed during scan")
	})
	t.Run("not_linked", func(t *testing.T) {
		_, _, err := registeredGCWorktree(boundedTestContext(t), root)
		if err == nil || err.Error() != "not a linked git worktree" ||
			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(originalContents, contents()) {
			t.Fatalf("primary repo was admitted or changed: %v", err)
		}
	})
	t.Run("activation_first", func(t *testing.T) {
		alias := filepath.Join(root, "alias.db")
		if err := os.Symlink(creds, alias); err != nil {
			t.Fatal(err)
		}
		held, err := store.OpenOperatorEpisodeLock(alias, "ep")
		if err != nil {
			t.Fatal(err)
		}
		defer held.Close()
		// Acquire the runner description before the daemon even scans.
		d := &Daemon{credentialsPath: creds}
		checkSkip(t, d.runAutoGCPass(ws, "", opts.Now), "operator episode lock unavailable")
		canonical, err := store.CanonicalPath(creds)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Lstat(store.OperatorEpisodeLockPath(canonical, "ep-wrong")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("contended pass created a different episode lock: %v", err)
		}
	})
	t.Run("recheck_after_ledger", func(t *testing.T) {
		o := opts
		o.lookupOperatorStep = func(_ context.Context, _, _ string) (bool, error) {
			if err := os.Chtimes(wt, o.Now, o.Now); err != nil {
				t.Fatal(err)
			}
			return false, nil
		}
		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "too recent (0s < 72h0m0s)")
		if err := os.Chtimes(wt, info.ModTime(), info.ModTime()); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("registration_rechecked", func(t *testing.T) {
		o := opts
		o.DryRun = true
		decoy := filepath.Join(root, "decoy")
		git(root, "worktree", "add", "--detach", decoy, "main")
		defer git(root, "worktree", "remove", decoy)
		meta := filepath.Join(wt, ".git")
		original, err := os.ReadFile(meta)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.WriteFile(meta, original, 0600); err != nil {
				t.Error(err)
			}
		}()
		foreign, err := os.ReadFile(filepath.Join(decoy, ".git"))
		if err != nil {
			t.Fatal(err)
		}
		registered := git(root, "worktree", "list", "--porcelain")
		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
			return false, os.WriteFile(meta, foreign, 0600)
		}
		rep := RunEpisodeGC(boundedTestContext(t), o)
		rewrite, err := os.ReadFile(meta)
		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
			err != nil || !bytes.Equal(rewrite, foreign) || !bytes.Equal(refs, git(root, "show-ref")) ||
			!bytes.Equal(registered, git(root, "worktree", "list", "--porcelain")) || !bytes.Equal(originalContents, contents()) {
			t.Fatalf("changed registration was accepted or modified: %+v, %v", rep, err)
		}
	})
	t.Run("unregistered", func(t *testing.T) {
		meta := filepath.Join(wt, ".git")
		original, err := os.ReadFile(meta)
		if err != nil {
			t.Fatal(err)
		}
		admin := strings.TrimSpace(strings.TrimPrefix(string(original), "gitdir: "))
		outside := filepath.Join(filepath.Dir(filepath.Dir(admin)), "retained", filepath.Base(admin))
		if err := os.MkdirAll(filepath.Dir(outside), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(admin, outside); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Rename(outside, admin); err != nil {
				t.Error(err)
			}
			if err := os.WriteFile(meta, original, 0600); err != nil {
				t.Error(err)
			}
		}()
		if err := os.WriteFile(meta, []byte("gitdir: "+outside+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
		// The checkout still works, and Git metadata points back to it.
		// The administrative directory is outside the live registrations.
		if got := git(wt, "rev-parse", "HEAD"); !bytes.Equal(got, head) {
			t.Fatalf("unregistered fixture has no original HEAD: %s", got)
		}
		o := opts
		o.DryRun = true
		rep := RunEpisodeGC(boundedTestContext(t), o)
		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
			!bytes.Equal(originalContents, contents()) || !bytes.Equal(refs, git(root, "show-ref")) {
			t.Fatalf("unregistered checkout changed or misclassified: %+v", rep)
		}
	})
	t.Run("duplicate_registration", func(t *testing.T) {
		decoy := filepath.Join(root, "duplicate")
		git(root, "worktree", "add", "--detach", decoy, "main")
		defer git(root, "worktree", "remove", decoy)
		admin := string(bytes.TrimSpace(git(decoy, "rev-parse", "--absolute-git-dir")))
		back := filepath.Join(admin, "gitdir")
		original, err := os.ReadFile(back)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.WriteFile(back, original, 0600); err != nil {
				t.Error(err)
			}
		}()
		if err := os.WriteFile(back, []byte(filepath.Join(wt, ".git")+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
		duplicate := git(root, "worktree", "list", "--porcelain")
		matching := 0
		for _, line := range strings.Split(string(duplicate), "\n") {
			if dir, ok := strings.CutPrefix(line, "worktree "); ok && physical(dir) == path {
				matching++
			}
		}
		if matching != 2 {
			t.Fatalf("Git did not report two spellings of the registration: %s", duplicate)
		}
		o := opts
		o.DryRun = true
		rep := RunEpisodeGC(boundedTestContext(t), o)
		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
			!bytes.Equal(originalContents, contents()) || !bytes.Equal(refs, git(root, "show-ref")) ||
			!bytes.Equal(duplicate, git(root, "worktree", "list", "--porcelain")) {
			t.Fatalf("duplicate registration changed or misclassified: %+v", rep)
		}
	})
	t.Run("gc_first", func(t *testing.T) {
		o := opts
		const refusal = "fixture stops before Git removal"
		calls := 0
		o.removeRegistered = func(ctx context.Context, git, repo, registered string, fd *os.File) error {
			calls++
			// The contender runs at the last boundary before the Git child.
			info, err := fd.Stat()
			lockPath, pathErr := store.CanonicalPath(creds)
			named, statErr := os.Stat(store.OperatorEpisodeLockPath(lockPath, "ep"))
			if err != nil || pathErr != nil || statErr != nil || !os.SameFile(info, named) {
				t.Fatalf("GC did not lend its canonical descriptor to Git: %v %v %v", err, pathErr, statErr)
			}
			bin, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			child := exec.CommandContext(childCtx, bin, "-test.run=^TestEpisodeGCOperatorSteps$")
			child.Env = append(os.Environ(), contenderEnv+"="+creds, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
			child.WaitDelay = time.Second
			if out, err := child.CombinedOutput(); err != nil || !bytes.Contains(out, []byte("PASS")) {
				t.Fatalf("independent activation contender: %v %s", err, out)
			}
			return errors.New(refusal)
		}
		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "remove failed: "+refusal)
		if calls != 1 {
			t.Fatalf("GC boundary fired %d times", calls)
		}
		contender, err := store.OpenOperatorEpisodeLock(creds, "ep")
		if err != nil {
			t.Fatalf("GC lock remained held after removal refusal: %v", err)
		}
		contender.Close()
	})
	t.Run("daemon_wiring", func(t *testing.T) {
		// Invoke the synchronous daemon pass with a temporary configured store.
		// The timers, server, launchctl and network are absent from this fixture.
		d := &Daemon{credentialsPath: creds}
		d.activity.begin("ep", opts.Now.Unix()-1, opts.Now.Unix()-2)()
		rep := d.runAutoGCPass(ws, "", opts.Now)
		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || filepath.Join(physical(filepath.Dir(rep.Removed[0].Path)), filepath.Base(rep.Removed[0].Path)) != path {
			t.Fatalf("daemon did not pass its readable credentials store to GC: %+v", rep)
		}
		if _, err := os.Stat(wt); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("daemon did not remove unprotected worktree: %v", err)
		}
	})
}
