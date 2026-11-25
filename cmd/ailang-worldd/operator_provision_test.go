package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorProvisionEpisodeExclusion(t *testing.T) {
	gc := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(gc) })

	t.Run("no_terminal", func(t *testing.T) {
		r := newEpisodeRun(t)
		creds, err := credentialsFor(r.f.db, "")
		if err != nil {
			t.Fatal(err)
		}
		before := observeProvision(t, r, creds)
		var out bytes.Buffer
		code := runSessionNew(r.f.args("ep", "--preset", "se-tools"), &out, &out, noTerminalEnv(1800))
		_, sidecarErr := os.Lstat(operatorEpisodeLockPath(creds, "ep"))
		if code != exitUsage || !strings.Contains(out.String(), "no controlling terminal") ||
			!os.IsNotExist(sidecarErr) || !reflect.DeepEqual(observeProvision(t, r, creds), before) {
			t.Fatalf("session new bypassed the terminal fence: %d %s", code, &out)
		}
	})

	t.Run("session_new", func(t *testing.T) {
		r := newEpisodeRun(t)
		cred, err := credentialsFor(r.f.db, "")
		if err != nil {
			t.Fatal(err)
		}
		held := claimFixtureEpisode(t, cred, "ep1")
		if err := held.Close(); err != nil {
			t.Fatal(err)
		}
		// The lease is stale and the synthetic holder has exited. The
		// durable action must still take precedence over a typed yes.
		var out bytes.Buffer
		term := &promptTerm{in: strings.NewReader("y\n")}
		env := testMintEnv("y\n", 1800)
		env.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
		token := filepath.Join(r.sessions, "custom.session")
		if err := os.Mkdir(r.sessions, 0700); err != nil {
			t.Fatal(err)
		}
		before := observeProvision(t, r, cred)
		code := runSessionNew(r.f.args("ep1", "--preset", "se-tools", "--branch", "new", "--out", token), &out, &out, env)
		expectEpisodeRefusal(t, code, out.String(), before.credentials, cred)
		if term.out.Len() != 0 || !reflect.DeepEqual(observeProvision(t, r, cred), before) {
			t.Fatal("contention opened a prompt or changed Git, store or session state")
		}

		// A different episode still mints, and the normal exit releases
		// its own descriptor. Do not rely on GC to close a leaked file.
		out.Reset()
		if code := runSessionNew(r.f.args("other", "--preset", "se-tools", "--out", token), &out, &out, env); code != exitOK {
			t.Fatalf("independent episode: %d %s", code, &out)
		}
		if rows := r.f.rows(t); len(rows) != 1 || rows[0].EpisodeID != "other" {
			t.Fatalf("independent mint did not persist: %+v", rows)
		}
		requireProvisionLockReleased(t, cred, "other")
	})
	t.Run("approve", func(t *testing.T) {
		a, sig := signedRig(t)
		held := claimFixtureEpisode(t, a.credPath, a.ep)
		before := observeProvision(t, a.r, a.credPath)
		request, approvers := a.row(t), approverRows(t, a.credPath)
		var out bytes.Buffer
		var opened bool
		env := approveEnv(nil, &opened, atTime(approveT0+1))
		code := runApprove(a.args("--signature", sig), &out, &out, env)
		expectEpisodeRefusal(t, code, out.String(), before.credentials, a.credPath)
		if opened || !reflect.DeepEqual(observeProvision(t, a.r, a.credPath), before) ||
			!reflect.DeepEqual(a.row(t), request) || !reflect.DeepEqual(approverRows(t, a.credPath), approvers) {
			t.Fatal("contended signature changed request, approver, Git, store or session state")
		}
		if err := held.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("approve_releases", func(t *testing.T) {
		a, sig := signedRig(t)
		var out bytes.Buffer
		var opened bool
		code := runApprove(a.args("--signature", sig), &out, &out, approveEnv(nil, &opened, atTime(approveT0+1)))
		if code != exitOK || opened || !a.row(t).Consumed {
			t.Fatalf("uncontended approval: %d %s", code, &out)
		}
		requireProvisionLockReleased(t, a.credPath, a.ep)
	})
}

type provisionSnapshot struct {
	credentials, worktrees, branches string
	workspace, sessions              map[string]string
}

func observeProvision(t *testing.T, r episodeRun, cred string) provisionSnapshot {
	t.Helper()
	return provisionSnapshot{serviceFileHash(cred), presetGit(t, r.f, "worktree", "list", "--porcelain"),
		presetGit(t, r.f, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads"),
		provisionFiles(t, r.f.ws), provisionFiles(t, r.sessions)}
}

func provisionFiles(t *testing.T, dir string) map[string]string {
	t.Helper()
	if _, err := os.Lstat(dir); os.IsNotExist(err) {
		return nil
	}
	return launchFileSnapshot(t, dir)
}

func requireProvisionLockReleased(t *testing.T, cred, ep string) {
	t.Helper()
	id, err := store.CanonicalPath(cred)
	if err != nil {
		t.Fatal(err)
	}
	path := operatorEpisodeLockPath(id, ep)
	fd, err := openRunnerLock(path, 0)
	if err != nil {
		t.Fatalf("cannot reacquire %s: %v", path, err)
	}
	fd.Close()
}
