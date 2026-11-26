package daemon

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestEpisodeGCOperatorSteps(t *testing.T) {
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
	git(root, "commit", "--allow-empty", "-m", "base")
	git(root, "worktree", "add", "--detach", wt, "main")
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
	head, status := git(wt, "rev-parse", "HEAD"), git(wt, "status", "--porcelain=v1")
	if len(status) != 0 {
		t.Fatalf("fixture is dirty: %q", status)
	}
	opts := EpisodeGCOptions{WorkspaceRoot: ws, Now: info.ModTime().Add(73 * time.Hour), OlderThan: 72 * time.Hour}
	checkSkip := func(t *testing.T, rep EpisodeGCReport, reason string) {
		t.Helper()
		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Name != "ep" ||
			rep.Skipped[0].Reason != reason || physical(rep.Skipped[0].Path) != path {
			t.Fatalf("want sole worktree skip %q: %+v", reason, rep)
		}
		if _, err := os.Stat(wt); err != nil || !bytes.Equal(head, git(wt, "rev-parse", "HEAD")) ||
			!bytes.Equal(status, git(wt, "status", "--porcelain=v1")) {
			t.Fatalf("refusal changed worktree or Git evidence: %v", err)
		}
	}
	t.Run("candidate", func(t *testing.T) {
		o := opts
		o.DryRun = true
		o.HasOpenOperatorStep = func(context.Context, string) (bool, error) { return false, nil }
		rep := RunEpisodeGC(boundedTestContext(t), o)
		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || rep.Removed[0].Reason != "clean & merged" ||
			physical(rep.Removed[0].Path) != path {
			t.Fatalf("successful negative lookup must allow the candidate: %+v", rep)
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
				o.HasOpenOperatorStep = func(_ context.Context, episode string) (bool, error) {
					return episode == "ep", nil
				}
			case "error":
				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
					return false, errors.New("fixture store error")
				}
			case "timeout":
				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
					return false, context.DeadlineExceeded
				}
			case "canceled":
				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
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
	t.Run("daemon_wiring", func(t *testing.T) {
		creds := filepath.Join(root, "credentials.db")
		if _, err := store.UpgradeCredentials(creds); err != nil {
			t.Fatal(err)
		}
		// Invoke the synchronous daemon pass with a temporary configured store.
		// The timers, server, launchctl and network are absent from this fixture.
		d := &Daemon{credentialsPath: creds}
		rep := d.runAutoGCPass(ws, "", opts.Now)
		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || filepath.Join(physical(filepath.Dir(rep.Removed[0].Path)), filepath.Base(rep.Removed[0].Path)) != path {
			t.Fatalf("daemon did not pass its readable credentials store to GC: %+v", rep)
		}
		if _, err := os.Stat(wt); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("daemon did not remove unprotected worktree: %v", err)
		}
	})
}
