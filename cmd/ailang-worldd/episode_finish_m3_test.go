package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

// This seam retains remote state across invocations. A timeout can happen on
// either side of the effect; querying that state is the only safe retry.
type statefulFinishDispatch struct {
	calls           []string
	all             []string
	head            string
	pr, label       bool
	fail            string
	timeout, landed bool
	pushes, creates int
	prepareErr      error
}

func (s *statefulFinishDispatch) event(name string, effect func()) error {
	s.calls = append(s.calls, name)
	s.all = append(s.all, name)
	if s.fail != name || s.landed {
		effect()
	}
	if s.fail == name {
		if s.timeout {
			return context.DeadlineExceeded
		}
		return fmt.Errorf("injected %s failure", name)
	}
	return nil
}

func (s *statefulFinishDispatch) Prepare(context.Context, string, string) (bool, error) {
	s.all = append(s.all, "prepare")
	return s.label, s.prepareErr
}

func (s *statefulFinishDispatch) RemoteHead(context.Context, string, string, string, string) (string, error) {
	err := s.event("ls-remote", func() {})
	return s.head, err
}

func (s *statefulFinishDispatch) Ancestor(ctx context.Context, git, wt, older, sha string) (bool, error) {
	err := s.event("ancestor", func() {})
	if err != nil {
		return false, err
	}
	_, err = worktree.Output(ctx, git, []string{"-C", wt, "merge-base", "--is-ancestor", older, sha})
	return err == nil, nil
}

func (s *statefulFinishDispatch) Push(_ context.Context, _, _, _, sha, _ string) error {
	return s.event("push", func() { s.head = sha; s.pushes++ })
}

func (s *statefulFinishDispatch) CreateLabel(context.Context, string, string) error {
	return s.event("label-create", func() { s.label = true })
}

func (s *statefulFinishDispatch) ViewPR(context.Context, string, string, string) (daemon.FinishPR, bool, error) {
	err := s.event("pr-view", func() {})
	return m3PR(), s.pr, err
}

func (s *statefulFinishDispatch) CreatePR(context.Context, string, string, string, string, string, string, bool) (daemon.FinishPR, error) {
	err := s.event("pr-create", func() { s.pr = true; s.creates++ })
	return m3PR(), err
}

func m3PR() daemon.FinishPR {
	return daemon.FinishPR{URL: "https://github.com/fixture/repo/pull/1", Number: 1, Draft: true}
}

func m3Fixture(t *testing.T) (*finishFixture, *statefulFinishDispatch) {
	t.Helper()
	f := newFinishFixture(t)
	f.allowRevoked = true
	s := &statefulFinishDispatch{label: true}
	finishDispatchFactory = func() finishDispatch { return s }
	// Keep M1/M2's fixture stable. This wrapper gives each approval a distinct
	// identity, as the real route does, so fresh approvals are observable.
	finishClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
			rec := httptest.NewRecorder()
			f.serve(rec, r)
			body := rec.Body.Bytes()
			if r.URL.Path == "/v1/episodes/ep1/finish" {
				var approval daemon.FinishApproval
				if json.Unmarshal(body, &approval) == nil && approval.Approval != "" {
					approval.Approval = fmt.Sprintf("sha256:%064x", len(f.scopes)+3)
					body, _ = json.Marshal(approval)
				}
			}
			for key, values := range rec.Header() {
				w.Header()[key] = values
			}
			w.WriteHeader(rec.Code)
			w.Write(body)
		})
	}
	old := finishAfterStep
	finishAfterStep = nil
	t.Cleanup(func() { finishAfterStep = old })
	return f, s
}

func m3Run(t *testing.T, f *finishFixture, answer string, want int, extra ...string) (string, string, string) {
	t.Helper()
	terminal := &captureFinishTerminal{Reader: strings.NewReader(answer)}
	env := testMintEnv("", 1)
	env.openTerminal = func() (io.ReadWriteCloser, error) { return terminal, nil }
	args := append([]string{"--include-unattributed"}, extra...)
	rc, out, errs := f.run(env, args...)
	if rc != want {
		t.Fatalf("exit %d, want %d\nstdout: %s\nstderr: %s\nterminal: %s", rc, want, out, errs, terminal.String())
	}
	return out, errs, terminal.String()
}

func m3Closed(t *testing.T, f *finishFixture) {
	t.Helper()
	st, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, exists, err := st.ResolveSession(boundedTestContext(t), f.credentialID); err != nil || exists {
		t.Fatalf("credential remains: %v %v", exists, err)
	}
}

func m3Calls(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("calls %v; want %v", got, want)
	}
}

func m3Kill(t *testing.T, f *finishFixture, step string) {
	t.Helper()
	sentinel := &struct{}{}
	finishAfterStep = func(at string) {
		if at == step {
			panic(sentinel)
		}
	}
	defer func() {
		finishAfterStep = nil
		if p := recover(); p != sentinel {
			t.Fatalf("kill %s: panic %v", step, p)
		}
	}()
	rc, out, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
	t.Fatalf("kill hook was not reached: exit %d, stdout %s, stderr %s", rc, out, errs)
}

func TestEpisodeFinishResumeAC13(t *testing.T) {
	for _, step := range []string{"commit", "approve", "dispatch", "record"} {
		t.Run(step, func(t *testing.T) {
			f, s := m3Fixture(t)
			f.rep.Start.Checks = []store.StartCheck{{Name: "pinned-check", Argv: []string{"fixture", "test"}, Cwd: ".", TimeoutS: 3, Exec: "pinned", Args: []string{"test"}}}
			m3Kill(t, f, step)
			first := append([]string(nil), s.calls...)
			n := len(s.calls)
			m3Run(t, f, "y\n", 0)
			if got := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-list", "--count", f.before+"..HEAD")); got != "1" {
				t.Fatal("commit count " + got)
			}
			if s.pushes != 1 || s.creates != 1 || len(f.finishRecords) != 1 || f.finishRecords[0].Outcome != "done" {
				t.Fatalf("pushes %d creates %d records %+v", s.pushes, s.creates, f.finishRecords)
			}
			full := []string{"ls-remote", "push", "pr-view", "pr-create"}
			switch step {
			case "commit", "approve":
				m3Calls(t, first, nil)
				m3Calls(t, s.calls[n:], full)
			case "dispatch":
				m3Calls(t, first, full)
				m3Calls(t, s.calls[n:], []string{"ls-remote", "pr-view"})
			case "record":
				m3Calls(t, first, full)
				m3Calls(t, s.calls[n:], []string{})
			}
			if f.checks != 1 || f.outcomes != 1 || len(f.finishRecords[0].Checks) != 1 || !reflect.DeepEqual(f.finishRecords[0].Unattributed, []string{"b.go"}) {
				t.Fatalf("resume lost evidence or repeated checks: checks %d outcomes %d record %+v", f.checks, f.outcomes, f.finishRecords[0])
			}
			approvals := 1
			if step == "approve" || step == "dispatch" {
				approvals = 2 // an orphan is never reusable
			}
			if len(f.scopes) != approvals {
				t.Fatalf("approvals %d want %d", len(f.scopes), approvals)
			}
			m3Closed(t, f)
			t.Logf("first %v; rerun %v", first, s.calls[n:])
		})
	}
}

func TestEpisodeFinishRerunAC10(t *testing.T) {
	for _, name := range []string{"push-failed", "pr-failed", "ambiguous-push-landed", "ambiguous-push-lost", "ambiguous-pr-created", "diverged"} {
		t.Run(name, func(t *testing.T) {
			f, s := m3Fixture(t)
			s.fail = "push"
			first := []string{"ls-remote", "push"}
			second := []string{"ls-remote", "push", "pr-view", "pr-create"}
			class := "push-failed"
			switch name {
			case "pr-failed":
				s.fail, s.label = "pr-create", false
				class = "pr-failed"
				first = []string{"ls-remote", "push", "label-create", "pr-view", "pr-create"}
				second = []string{"ls-remote", "pr-view", "pr-create"}
			case "ambiguous-push-landed", "ambiguous-push-lost":
				s.timeout, s.landed = true, name == "ambiguous-push-landed"
				class = "ambiguous"
				if s.landed {
					second = []string{"ls-remote", "pr-view", "pr-create"}
				}
			case "ambiguous-pr-created":
				s.fail, s.timeout, s.landed = "pr-create", true, true
				class = "ambiguous"
				first = []string{"ls-remote", "push", "pr-view", "pr-create"}
				second = []string{"ls-remote", "pr-view"}
			case "diverged":
				s.fail, s.head = "", strings.Repeat("f", 40)
				first = []string{"ls-remote", "ancestor"}
				second = first
			}
			m3Run(t, f, "y\n", 2)
			m3Calls(t, s.calls, first)
			if f.finishRecords[0].Outcome != class {
				t.Fatalf("class %+v", f.finishRecords[0])
			}
			s.fail = ""
			n := len(s.calls)
			want := 0
			if name == "diverged" {
				want = 2
			}
			_, _, terminal := m3Run(t, f, "", want)
			m3Calls(t, s.calls[n:], second)
			if len(f.scopes) != 1 || len(f.finishRecords) != 2 || f.finishRecords[0].Approval != f.finishRecords[1].Approval || !strings.Contains(terminal, "resuming under approval ") {
				t.Fatalf("reuse: scopes %v records %+v terminal %s", f.scopes, f.finishRecords, terminal)
			}
			if name != "diverged" && (s.pushes != 1 || s.creates != 1) {
				t.Fatalf("pushes %d creates %d", s.pushes, s.creates)
			}
			m3Closed(t, f)
			t.Logf("first %v; rerun %v", first, s.calls[n:])
		})
	}
}

func TestEpisodeFinishApprovalReuse(t *testing.T) {
	for _, name := range []string{"near-expiry", "ready-changed", "body-changed", "label-missing-unscoped"} {
		t.Run(name, func(t *testing.T) {
			f, s := m3Fixture(t)
			s.fail = "pr-create"
			m3Run(t, f, "y\n", 2)
			s.fail = ""
			n, all := len(s.calls), len(s.all)
			var extra []string
			switch name {
			case "near-expiry":
				scope, _ := broker.ParseFinishApprovalScope(f.finishRecords[0].Scope)
				scope.ExpiresAt = time.Now().Unix() + 299
				f.finishRecords[0].Scope = broker.FormatFinishApprovalScope(scope)
			case "ready-changed":
				extra = []string{"--ready"}
			case "body-changed":
				path := filepath.Join(f.sessions, "ep1.pr-body.md")
				data, _ := os.ReadFile(path)
				os.WriteFile(path, append([]byte("changed\n"), data...), 0o600)
			case "label-missing-unscoped":
				s.label = false
			}
			want := 0
			if name == "body-changed" {
				want = 1
			}
			_, errs, terminal := m3Run(t, f, "y\n", want, extra...)
			if want == 1 {
				if len(s.all) != all || len(f.finishRecords) != 1 || !strings.Contains(errs, "changed since it was recorded") {
					t.Fatalf("body guard: %s %v", errs, s.all[all:])
				}
				return
			}
			if len(f.scopes) != 2 || !strings.Contains(terminal, "Open this pull request?") {
				t.Fatalf("no new prompt: %s scopes %v", terminal, f.scopes)
			}
			scope, _ := broker.ParseFinishApprovalScope(f.scopes[1])
			if name == "label-missing-unscoped" && !scope.CreateLabel {
				t.Fatal("label creation not scoped")
			}
			if name == "ready-changed" && scope.Draft {
				t.Fatal("still draft")
			}
			if len(s.calls[n:]) == 0 {
				t.Fatal("no retry")
			}
		})
	}
}

func m3MintAgain(t *testing.T, f *finishFixture) {
	t.Helper()
	st, err := store.OpenCredentials(f.health.CredentialsDB)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	_, f.credentialID, _, err = authority.Mint(boundedTestContext(t), st, "ep1", []broker.Capability{{Effect: "Workspace.Write", Scope: "world", Budget: 10, ExpiresAt: time.Now().Unix() + 3600}}, 3600, time.Now().Unix(), nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEpisodeFinishFixupRound(t *testing.T) {
	t.Run("same-episode", func(t *testing.T) {
		f, s := m3Fixture(t)
		m3Run(t, f, "y\n", 0)
		first := f.finishRecords[0].Commit
		f.write("a.go", "fix-up\n")
		m3MintAgain(t, f)
		n := len(s.calls)
		m3Run(t, f, "y\n", 0)
		m3Calls(t, s.calls[n:], []string{"ls-remote", "ancestor", "push", "pr-view"})
		if len(f.finishRecords) != 2 || f.finishRecords[1].Approval == f.finishRecords[0].Approval || f.finishRecords[1].Trailer != f.finishRecords[0].Trailer {
			t.Fatalf("records %+v", f.finishRecords)
		}
		if got := strings.TrimSpace(f.gitOutput("-C", f.wt, "show", "--format=", "--name-only", "HEAD")); got != "a.go" {
			t.Fatal(got)
		}
		if parent := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD^")); parent != first {
			t.Fatal(parent)
		}
		t.Logf("change base %s; fix-up only a.go", first)
	})
	t.Run("branch-freed", func(t *testing.T) {
		for _, failed := range []bool{false, true} {
			f, s := m3Fixture(t)
			if failed {
				s.fail = "pr-create"
			}
			want := 0
			if failed {
				want = 2
			}
			m3Run(t, f, "y\n", want)
			branch, err := worktree.Output(boundedTestContext(t), f.git, []string{"-C", f.wt, "symbolic-ref", "--short", "HEAD"})
			list := f.gitOutput("-C", f.wt, "worktree", "list", "--porcelain")
			if failed {
				if err != nil || strings.TrimSpace(branch) != f.rep.Start.Branch || !strings.Contains(list, "branch refs/heads/"+f.rep.Start.Branch) {
					t.Fatalf("failed branch lost: %s %v %s", branch, err, list)
				}
				continue
			}
			if err == nil || strings.Contains(list, "branch refs/heads/"+f.rep.Start.Branch) {
				t.Fatalf("branch not freed: %s %v %s", branch, err, list)
			}
			if head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")); head != f.finishRecords[0].Commit {
				t.Fatal("detach changed HEAD")
			}
			for path, want := range map[string]string{"a.go": "agent\n", "gen.txt": "produced\n", "b.go": "hand\n"} {
				data, err := os.ReadFile(filepath.Join(f.wt, path))
				if err != nil || string(data) != want {
					t.Fatalf("kept file %s %q %v", path, data, err)
				}
			}
			n := len(s.all)
			m3Run(t, f, "", 0)
			if len(s.all) != n {
				t.Fatal(s.all[n:])
			}
		}
	})
	t.Run("new-episode", func(t *testing.T) {
		f, _ := m3Fixture(t)
		f.write("ep0.txt", "previous episode\n")
		f.gitOutput("-C", f.wt, "add", "ep0.txt")
		f.gitOutput("-C", f.wt, "commit", "-qm", "ep0\n\nWorld-Provenance: store=fixture episode=ep0 entries=1-2")
		previous := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
		f.rep.Start.BranchExisted = true
		m3Run(t, f, "y\n", 0)
		paths := f.gitOutput("-C", f.wt, "show", "--format=", "--name-only", "HEAD")
		if strings.Contains(paths, "ep0.txt") || containsString(f.finishRecords[0].Unattributed, "ep0.txt") {
			t.Fatalf("previous round included: %s %+v", paths, f.finishRecords)
		}
		if parent := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD^")); parent != previous {
			t.Fatal(parent)
		}
		t.Logf("change base %s; ep0.txt excluded", previous)
	})
}

func TestEpisodeFinishFinishedRerun(t *testing.T) {
	for _, name := range []string{"finished", "revoke-retry"} {
		t.Run(name, func(t *testing.T) {
			f, s := m3Fixture(t)
			var lock *store.Credentials
			if name == "revoke-retry" {
				f.recordHook = func() {
					var err error
					lock, err = store.OpenCredentials(f.health.CredentialsDB)
					if err != nil {
						t.Fatal(err)
					}
				}
				m3Run(t, f, "y\n", 2)
				lock.Close()
				f.recordHook = nil
			} else {
				m3Run(t, f, "y\n", 0)
			}
			n, requests := len(s.all), len(f.requests)
			out, _, _ := m3Run(t, f, "", 0)
			if !strings.Contains(out, "episode ep1 is already finished: "+m3PR().URL) || len(s.all) != n || len(f.finishRecords) != 1 {
				t.Fatalf("%s calls %v records %+v", out, s.all[n:], f.finishRecords)
			}
			for _, path := range f.requests[requests:] {
				if path == "/v1/episodes/ep1/finish" || path == "/v1/session/self" {
					t.Fatal("finished called " + path)
				}
			}
			m3Closed(t, f)
		})
	}
}

func TestEpisodeFinishResumeGuards(t *testing.T) {
	for _, name := range []string{"no-session", "stale-body", "title-mismatch", "dry-run"} {
		t.Run(name, func(t *testing.T) {
			f, s := m3Fixture(t)
			s.fail = "push"
			m3Run(t, f, "y\n", 2)
			s.fail = ""
			n, requests := len(s.all), len(f.requests)
			var extra []string
			switch name {
			case "no-session":
				f.expired = true
				os.Remove(filepath.Join(f.sessions, "ep1.session"))
				f.rep.Uncommitted = 1 // F3 is skipped for the pending commit.
				f.write("later.txt", "not in this finish\n")
			case "stale-body":
				os.WriteFile(filepath.Join(f.sessions, "ep1.pr-body.md"), []byte("stale\n"), 0o600)
			case "title-mismatch":
				extra = []string{"--title", "different"}
			case "dry-run":
				extra = []string{"--dry-run"}
			}
			want := 1
			if name == "no-session" {
				want = 0
			}
			out, errs, _ := m3Run(t, f, "", want, extra...)
			if want == 1 && (len(s.all) != n || len(f.finishRecords) != 1) {
				t.Fatalf("guard dispatched: %v %+v %s", s.all[n:], f.finishRecords, errs)
			}
			if name == "no-session" {
				if !strings.Contains(out, "1 uncommitted paths are not part of this finish") {
					t.Fatal(out)
				}
				for _, path := range f.requests[requests:] {
					if path == "/v1/session/self" {
						t.Fatal("resume checked session")
					}
				}
			}
		})
	}
}
