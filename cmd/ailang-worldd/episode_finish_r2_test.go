package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func r2Head(f *finishFixture) string {
	return strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
}

func r2Detached(t *testing.T, f *finishFixture) {
	t.Helper()
	if branch, err := worktree.Output(boundedTestContext(t), f.git, []string{"-C", f.wt, "symbolic-ref", "--short", "HEAD"}); err == nil {
		t.Fatalf("HEAD attached to %s", branch)
	}
}

func r2Fixup(t *testing.T) (*finishFixture, *statefulFinishDispatch) {
	t.Helper()
	f, s := m3Fixture(t)
	m3Run(t, f, "y\n", 0)
	m3MintAgain(t, f)
	f.write("a.go", "fix-up\n")
	return f, s
}

func TestEpisodeFinishFixupKilledAfterCommitNotFinished(t *testing.T) {
	f, s := r2Fixup(t)
	c1 := r2Head(f)
	m3Kill(t, f, "commit")
	c2 := r2Head(f)
	if c1 == c2 {
		t.Fatal("fix-up did not commit")
	}
	n := len(s.all)
	out, _, terminal := m3Run(t, f, "y\n", 0)
	if strings.Contains(out, "already finished") || !strings.Contains(terminal, "[y/N]") || s.head != c2 || len(s.all) == n || len(f.finishRecords) != 2 || f.finishRecords[1].Commit != c2 {
		t.Fatalf("pending fix-up skipped: out %s terminal %s calls %v records %+v", out, terminal, s.all[n:], f.finishRecords)
	}
}

func TestEpisodeFinishStaleRoundBodyRefused(t *testing.T) {
	for _, name := range []string{"commit-line", "trailer-line"} {
		t.Run(name, func(t *testing.T) {
			f, s := r2Fixup(t)
			path := filepath.Join(f.sessions, "ep1.pr-body.md")
			old, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			c1 := r2Head(f)
			m3Kill(t, f, "commit")
			c2 := r2Head(f)
			body := string(old)
			if name == "trailer-line" {
				body = strings.Replace(body, "Commit: "+c1, "Commit: "+c2, 1)
				body = strings.Replace(body, f.finishRecords[0].Trailer, f.finishRecords[0].Trailer+" wrong", 1)
			}
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			n := len(s.all)
			_, errs, _ := m3Run(t, f, "y\n", 1)
			if len(s.all) != n || !strings.Contains(errs, "is not for "+c2) || len(f.finishRecords) != 1 {
				t.Fatalf("stale body accepted: %s calls %v", errs, s.all[n:])
			}
		})
	}
}

func TestEpisodeFinishDetachedPushFailedRefuses(t *testing.T) {
	for _, name := range []string{"push-failed", "done-other-commit"} {
		t.Run(name, func(t *testing.T) {
			f, s := m3Fixture(t)
			s.fail = "push"
			m3Run(t, f, "y\n", 2)
			f.gitOutput("-C", f.wt, "-c", "core.hooksPath=/dev/null", "checkout", "--detach")
			if name == "done-other-commit" {
				f.finishRecords[0].Outcome = "done"
				f.finishRecords[0].Commit = f.before
			}
			s.fail = ""
			n := len(s.all)
			_, errs, _ := m3Run(t, f, "y\n", 1)
			if len(s.all) != n || !strings.Contains(errs, "HEAD is not on "+f.rep.Start.Branch) {
				t.Fatalf("detached pending finish accepted: %s %v", errs, s.all[n:])
			}
		})
	}
}

func TestEpisodeFinishRevokeRealErrorIsF10(t *testing.T) {
	f, _ := m3Fixture(t)
	old := finishRevokeCredentials
	t.Cleanup(func() { finishRevokeCredentials = old })
	called := false
	finishRevokeCredentials = func(ctx context.Context, _ *store.Credentials, ep string, out io.Writer) error {
		called = true
		ro, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
		if err != nil {
			t.Fatal(err)
		}
		defer ro.Close()
		err = revokeEpisodeCredentials(ctx, ro, ep, out)
		if err == nil || err == errNoEpisodeCredential {
			t.Fatalf("need real revoke failure, got %v", err)
		}
		return err
	}
	out, errs, _ := m3Run(t, f, "y\n", 2)
	if !called || !strings.Contains(errs, "F10") || strings.Contains(out, "session already closed") || len(f.finishRecords) != 1 {
		t.Fatalf("revoke error swallowed: %s %s", out, errs)
	}
}

func TestEpisodeFinishFixupBranchMovedRefuses(t *testing.T) {
	f, s := r2Fixup(t)
	f.gitOutput("-C", f.wt, "update-ref", "refs/heads/"+f.rep.Start.Branch, f.before)
	n := len(s.all)
	_, errs, _ := m3Run(t, f, "y\n", 1)
	if !strings.Contains(errs, "recorded branch moved") || len(s.all) != n {
		t.Fatalf("moved branch accepted: %s %v", errs, s.all[n:])
	}
	r2Detached(t, f)
}

func TestEpisodeFinishChangeBaseNewestFirstParent(t *testing.T) {
	f, _ := m3Fixture(t)
	commit := func(title string) string {
		f.gitOutput("-C", f.wt, "commit", "--allow-empty", "-qm", title+"\n\nWorld-Provenance: store=fixture episode=other entries=1-2")
		return r2Head(f)
	}
	oldest := commit("oldest World")
	newest := commit("newest first-parent World")
	// Git orders by committer seconds. Make the side commit strictly newer so
	// dropping --first-parent deterministically selects it instead.
	time.Sleep(time.Until(time.Unix(time.Now().Unix()+1, 0)))
	f.gitOutput("-C", f.wt, "checkout", "-qb", "side", oldest)
	side := commit("side World")
	f.gitOutput("-C", f.wt, "checkout", f.rep.Start.Branch)
	f.gitOutput("-C", f.wt, "merge", "--no-ff", "-m", "merge side", "side")
	query := func(args ...string) (string, error) {
		return worktree.Output(boundedTestContext(t), f.git, append([]string{"-C", f.wt}, args...))
	}
	base, err := finishChangeBase(query, f.rep.Start)
	if err != nil || base != newest || base == side || base == oldest {
		t.Fatalf("base %s want newest first-parent %s (old %s side %s): %v", base, newest, oldest, side, err)
	}
}

func TestEpisodeFinishDetachAfterRecord(t *testing.T) {
	f, s := m3Fixture(t)
	f.recordFailure = true
	m3Run(t, f, "y\n", 2)
	if branch := strings.TrimSpace(f.gitOutput("-C", f.wt, "symbolic-ref", "--short", "HEAD")); branch != f.rep.Start.Branch {
		t.Fatalf("record failure released branch: %s", branch)
	}
	f.recordFailure = false
	n := len(s.calls)
	m3Run(t, f, "y\n", 0)
	m3Calls(t, s.calls[n:], []string{"ls-remote", "pr-view"})
	if len(f.finishRecords) != 1 || f.finishRecords[0].Outcome != "done" {
		t.Fatalf("resume records %+v", f.finishRecords)
	}
	r2Detached(t, f)
}

func TestEpisodeFinishReuseTermsEach(t *testing.T) {
	f, s := m3Fixture(t)
	s.fail = "push"
	m3Run(t, f, "y\n", 2)
	record := f.finishRecords[0]
	want, err := broker.ParseFinishApprovalScope(record.Scope)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	if _, ok := finishReusableApproval(&record, "resume-dispatch", want, true, now); !ok {
		t.Fatal("control approval not reusable")
	}
	for _, name := range []string{"sha", "body", "title", "base", "head", "remote", "repo", "mode"} {
		t.Run(name, func(t *testing.T) {
			changed, mode := want, "resume-dispatch"
			switch name {
			case "sha":
				changed.SHA = strings.Repeat("f", 40)
			case "body":
				changed.BodySHA256 = strings.Repeat("f", 64)
			case "title":
				changed.TitleSHA256 = strings.Repeat("f", 64)
			case "base":
				changed.Base += "-other"
			case "head":
				changed.Head += "-other"
			case "remote":
				changed.Remote += "-other"
			case "repo":
				changed.Repository += "-other"
			case "mode":
				mode = "resume"
			}
			if _, ok := finishReusableApproval(&record, mode, changed, true, now); ok {
				t.Fatal("changed " + name + " reused approval")
			}
		})
	}
}

func TestEpisodeFinishFixupReattach(t *testing.T) {
	for _, name := range []string{"hooks-off", "dry-run-stays-detached"} {
		t.Run(name, func(t *testing.T) {
			f, s := r2Fixup(t)
			c1 := r2Head(f)
			hooks := filepath.Join(f.root, "hooks")
			if err := os.MkdirAll(hooks, 0o700); err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(f.root, "checkout-marker")
			if err := os.WriteFile(filepath.Join(hooks, "post-checkout"), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0o700); err != nil {
				t.Fatal(err)
			}
			f.gitOutput("-C", f.wt, "config", "core.hooksPath", hooks)
			n := len(s.all)
			if name == "hooks-off" {
				m3Run(t, f, "y\n", 0)
				if r2Head(f) == c1 || s.head != r2Head(f) {
					t.Fatal("fix-up did not re-attach and finish")
				}
			} else {
				out, _, _ := m3Run(t, f, "", 0, "--dry-run")
				if !strings.Contains(out, "would re-attach") || r2Head(f) != c1 || len(s.all) != n {
					t.Fatalf("dry run changed state: %s %v", out, s.all[n:])
				}
				r2Detached(t, f)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("post-checkout hook ran: %v", err)
			}
		})
	}
}
