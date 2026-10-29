package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
)

type fakeFinishDispatch struct {
	calls    []string
	prepares int
	label    bool
	existing bool
	fail     string
	timeout  bool
	diverged bool
	sha      string
	base     string
	title    string
	bodyFile string
	body     string
	draft    bool
	remote   string
	branch   string
	repo     string
}

func (f *fakeFinishDispatch) event(name string) error {
	f.calls = append(f.calls, name)
	if f.fail == name {
		if f.timeout {
			return context.DeadlineExceeded
		}
		return fmt.Errorf("injected %s failure", name)
	}
	return nil
}

func (f *fakeFinishDispatch) Prepare(context.Context, string, string) (bool, error) {
	f.prepares++
	return f.label, nil
}

func (f *fakeFinishDispatch) RemoteHead(context.Context, string, string, string, string) (string, error) {
	err := f.event("ls-remote")
	if f.diverged {
		return strings.Repeat("f", 40), err
	}
	return "", err
}

func (f *fakeFinishDispatch) Ancestor(context.Context, string, string, string, string) (bool, error) {
	return false, f.event("ancestor")
}

func (f *fakeFinishDispatch) Push(_ context.Context, _, _, remote, sha, branch string) error {
	f.sha, f.remote, f.branch = sha, remote, branch
	return f.event("push")
}

func (f *fakeFinishDispatch) CreateLabel(context.Context, string, string) error {
	return f.event("label-create")
}

func (f *fakeFinishDispatch) ViewPR(_ context.Context, _, repo, _ string) (daemon.FinishPR, bool, error) {
	f.repo = repo
	err := f.event("pr-view")
	return daemon.FinishPR{URL: "https://github.com/fixture/repo/pull/1", Number: 1, Draft: true}, f.existing, err
}

func (f *fakeFinishDispatch) CreatePR(_ context.Context, _, _, base, _, title, body string, draft bool) (daemon.FinishPR, error) {
	f.base, f.title, f.bodyFile, f.draft = base, title, body, draft
	data, err := os.ReadFile(body)
	if err != nil {
		return daemon.FinishPR{}, err
	}
	f.body = string(data)
	err = f.event("pr-create")
	return daemon.FinishPR{URL: "https://github.com/fixture/repo/pull/1", Number: 1, Draft: draft}, err
}

func TestEpisodeFinishDeclineAC7(t *testing.T) {
	f := newFinishFixture(t)
	rc, _, errs := f.run(testMintEnv("n\n", 1), "--include-unattributed")
	if rc != 1 || len(f.seam.calls) != 0 || len(f.scopes) != 0 || len(f.finishRecords) != 0 ||
		strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) == f.before {
		t.Fatalf("decline: %d %s calls %v", rc, errs, f.seam.calls)
	}
	if !strings.Contains(errs, "next:") || !strings.Contains(errs, "gh pr create") {
		t.Fatal(errs)
	}
}

func TestEpisodeFinishBytesAC8(t *testing.T) {
	for _, change := range []string{"body", "HEAD", "remote"} {
		t.Run(change, func(t *testing.T) {
			f := newFinishFixture(t)
			finishBeforeDispatch = func() {
				switch change {
				case "body":
					path := filepath.Join(f.sessions, "ep1.pr-body.md")
					b, _ := os.ReadFile(path)
					b[0] ^= 1
					os.WriteFile(path, b, 0o600)
				case "HEAD":
					f.gitOutput("-C", f.wt, "commit", "--allow-empty", "-qm", "changed after approval")
				case "remote":
					f.gitOutput("-C", f.wt, "remote", "set-url", "origin", "https://github.com/other/repo")
				}
			}
			rc, _, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
			if rc != 1 || !strings.Contains(errs, "the commit or body changed since approval") || len(f.seam.calls) != 0 ||
				len(f.scopes) != 1 || len(f.finishRecords) != 0 {
				t.Fatalf("%d %s calls %v", rc, errs, f.seam.calls)
			}
		})
	}
}

func TestEpisodeFinishDispatchAC9AC11AC12(t *testing.T) {
	for _, name := range []string{"missing-label", "existing-pr", "ready"} {
		t.Run(name, func(t *testing.T) {
			f := newFinishFixture(t)
			f.seam.label = name != "missing-label"
			f.seam.existing = name == "existing-pr"
			args := []string{"--include-unattributed", "--title", "Approved fixture"}
			if name == "ready" {
				args = append(args, "--ready")
			}
			env := testMintEnv("y\n", 1)
			rc, out, errs := f.run(env, args...)
			want := []string{"ls-remote", "push"}
			if !f.seam.label {
				want = append(want, "label-create")
			}
			want = append(want, "pr-view")
			if !f.seam.existing {
				want = append(want, "pr-create")
			}
			if rc != 0 || !reflect.DeepEqual(f.seam.calls, want) || len(f.finishRecords) != 1 {
				t.Fatalf("%d %s calls %v", rc, errs, f.seam.calls)
			}
			record := f.finishRecords[0]
			if record.Outcome != "done" || record.PR.URL != "https://github.com/fixture/repo/pull/1" || !strings.Contains(out, record.PR.URL) {
				t.Fatalf("%+v %s", record, out)
			}
			scope, err := broker.ParseFinishApprovalScope(f.scopes[0])
			if err != nil || scope.SHA != f.seam.sha || scope.Base != "dev" || scope.Draft != (name != "ready") ||
				scope.TitleSHA256 != finishSHA([]byte("Approved fixture")) || scope.CreateLabel != !f.seam.label || f.seam.remote != "origin" || f.seam.branch != scope.Head || f.seam.repo != scope.Repository {
				t.Fatalf("scope %+v %v", scope, err)
			}
			if !f.seam.existing {
				if f.seam.base != scope.Base || f.seam.title != "Approved fixture" || f.seam.draft != scope.Draft ||
					finishSHA([]byte(f.seam.body)) != scope.BodySHA256 || !strings.HasSuffix(strings.TrimSpace(f.seam.body), record.Trailer) {
					t.Fatal("approved bytes differ from PR")
				}
				for _, part := range []string{fmt.Sprintf("%d calls.", f.rep.TotalCalls), "Calls by tool:", "Plan refusals: 1", "#12 ailang-read: missing path", "Handler refusals: 1", "#15 ailang-cli: fmt write refuses", "Workspace.Exec@worktree: 91 / 509", "Budgets (used / left):", "### Not made through World"} {
					if !strings.Contains(f.seam.body, part) {
						t.Error("body missing " + part)
					}
				}
			}
			creds, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
			if err != nil {
				t.Fatal(err)
			}
			_, exists, err := creds.ResolveSession(boundedTestContext(t), f.credentialID)
			creds.Close()
			if err != nil || exists {
				t.Fatalf("not revoked: %v %v", exists, err)
			}
			t.Logf("scope: %s\ncalls: %v", f.scopes[0], f.seam.calls)
		})
	}
}

func TestEpisodeFinishOutcomesAC10(t *testing.T) {
	for _, name := range []string{"push-failed", "pr-failed", "ambiguous", "diverged"} {
		t.Run(name, func(t *testing.T) {
			f := newFinishFixture(t)
			f.seam.fail = "push"
			want := []string{"ls-remote", "push"}
			class := name
			switch name {
			case "pr-failed":
				f.seam.fail = "pr-create"
				want = append(want, "pr-view", "pr-create")
			case "ambiguous":
				f.seam.timeout = true
			case "diverged":
				class = "push-failed"
				f.seam.fail = ""
				f.seam.diverged = true
				want = []string{"ls-remote", "ancestor"}
			}
			rc, _, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
			if rc != 2 || len(f.finishRecords) != 1 || f.finishRecords[0].Outcome != class || !reflect.DeepEqual(f.seam.calls, want) {
				t.Fatalf("%d %s records %+v calls %v", rc, errs, f.finishRecords, f.seam.calls)
			}
			t.Logf("%s: %v", class, f.seam.calls)
		})
	}
}

func TestFinishRepository(t *testing.T) {
	for _, remote := range []string{"https://github.com/o/r.git", "git@github.com:o/r", "ssh://git@github.com/o/r.git"} {
		if got, err := finishRepository(remote); err != nil || got != "o/r" {
			t.Fatalf("%s: %s %v", remote, got, err)
		}
	}
	for _, remote := range []string{"https://evil.example/o/r", "https://github.com/o/r#other", "https://github.com/o/r/extra"} {
		if _, err := finishRepository(remote); err == nil {
			t.Fatal(remote)
		}
	}
}

type captureFinishTerminal struct {
	io.Reader
	bytes.Buffer
}

func (t *captureFinishTerminal) Read(p []byte) (int, error) {
	return t.Reader.Read(p)
}

func (t *captureFinishTerminal) Close() error {
	return nil
}

func TestEpisodeFinishProtected(t *testing.T) {
	for _, changedDigest := range []bool{false, true} {
		t.Run(fmt.Sprint(changedDigest), func(t *testing.T) {
			f := newFinishFixture(t)
			profiles := filepath.Join(f.root, "profiles")
			os.MkdirAll(profiles, 0o700)
			raw, _ := json.Marshal(map[string]any{
				"profile":    broker.ExecProfileVersion,
				"project":    "fixture",
				"path":       []string{"/usr/bin", "/bin"},
				"timeout_ms": 9000,
				"protected":  []string{"*_test.go"},
				"commands":   map[string]any{"check": map[string]any{"argv": []string{"true"}, "max_args": 1}},
			})
			profile, err := broker.ParseExecProfile(raw)
			if err != nil {
				t.Fatal(err)
			}
			file := filepath.Join(profiles, "fixture.json")
			os.WriteFile(file, raw, 0o600)
			creds, err := store.OpenCredentials(f.health.CredentialsDB)
			if err != nil {
				t.Fatal(err)
			}
			if err := authority.Revoke(boundedTestContext(t), creds, f.credentialID); err != nil {
				t.Fatal(err)
			}
			_, f.credentialID, _, err = authority.MintBound(boundedTestContext(t), creds, "ep1",
				authority.ProfileBinding{Project: "fixture", Digest: profile.Digest},
				[]broker.Capability{{Effect: "Workspace.Write", Scope: "world", Budget: 10, ExpiresAt: 4000000000}}, 3600, time.Now().Unix(), nil)
			creds.Close()
			if err != nil {
				t.Fatal(err)
			}
			f.rep.Start.Preset.Bytes = `{"profile":"fixture"}`
			f.write("a_test.go", "protected edit\n")
			if changedDigest {
				os.WriteFile(file, append(raw, '\n'), 0o600)
			}
			terminal := &captureFinishTerminal{Reader: strings.NewReader("y\n")}
			env := testMintEnv("", 1)
			env.openTerminal = func() (io.ReadWriteCloser, error) { return terminal, nil }
			rc, _, errs := f.run(env, "--include-unattributed", "--exec-profiles-dir", profiles)
			if changedDigest {
				if rc != 1 || !strings.Contains(errs, "the exec profile fixture changed since episode ep1 started") || len(f.seam.calls) != 0 {
					t.Fatalf("%d %s", rc, errs)
				}
				return
			}
			if rc != 0 {
				t.Fatalf("%d %s", rc, errs)
			}
			for where, text := range map[string]string{
				"commit": f.gitOutput("-C", f.wt, "log", "-1", "--format=%B"),
				"prompt": terminal.String(), "body": f.seam.body,
			} {
				if !strings.Contains(text, "rotected paths changed: a_test.go") {
					t.Errorf("%s: %s", where, text)
				}
			}
			t.Log("prompt:\n" + terminal.String())
		})
	}
}

func TestFinishProductionGHArgv(t *testing.T) {
	f := newFinishFixture(t)
	bin := t.TempDir()
	log := filepath.Join(bin, "argv")
	script := `#!/bin/sh
test -z "$WORLD_SESSION" || exit 19
printf '%s\n' "$@" >> "$FINISH_ARGV"
case "$1 $2" in
  "label list") printf 'true\n' ;;
  "pr view") printf '{"url":"https://github.com/fixture/repo/pull/1","number":1,"isDraft":true,"state":"OPEN"}\n' ;;
  "pr create") printf 'https://github.com/fixture/repo/pull/1\n' ;;
  "label create") exit 0 ;;
  *) exit 21 ;;
esac
`
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	t.Setenv("FINISH_ARGV", log)
	t.Setenv("WORLD_SESSION", "fixture-sentinel")
	p := &productionFinishDispatch{}
	if found, err := p.Prepare(boundedTestContext(t), f.wt, "fixture/repo"); err != nil || !found {
		t.Fatalf("prepare %v %v", found, err)
	}
	if err := p.CreateLabel(boundedTestContext(t), f.wt, "fixture/repo"); err != nil {
		t.Fatal(err)
	}
	if _, found, err := p.ViewPR(boundedTestContext(t), f.wt, "fixture/repo", "custom-branch"); err != nil || !found {
		t.Fatalf("view %v %v", found, err)
	}
	body := filepath.Join(f.sessions, "body.md")
	os.WriteFile(body, []byte("fixture"), 0o600)
	if _, err := p.CreatePR(boundedTestContext(t), f.wt, "fixture/repo", "dev", "custom-branch", "Fixture title", body, true); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(log)
	text := string(data)
	if strings.Count(text, "-R\nfixture/repo\n") != 4 {
		t.Fatal(text)
	}
	for _, want := range []string{"--base\ndev\n--head\ncustom-branch\n--title\nFixture title\n--body-file\n" + body,
		"--label\nailang-world\n--draft\n", "--color\ne73c17\n--description\n" + finishLabelDescription} {
		if !strings.Contains(text, want) {
			t.Fatal(text)
		}
	}
}

func TestEpisodeFinishRecordFailureKeepsSession(t *testing.T) {
	f := newFinishFixture(t)
	f.recordFailure = true
	rc, _, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
	if rc != 2 || !strings.Contains(errs, "session kept") || len(f.finishRecords) != 0 {
		t.Fatalf("%d %s", rc, errs)
	}
	creds, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
	if err != nil {
		t.Fatal(err)
	}
	defer creds.Close()
	if _, exists, err := creds.ResolveSession(boundedTestContext(t), f.credentialID); err != nil || !exists {
		t.Fatalf("record failure revoked: %v %v", exists, err)
	}
}

func TestEpisodeFinishRevokeFailureF10(t *testing.T) {
	f := newFinishFixture(t)
	var lock *store.Credentials
	f.recordHook = func() {
		var err error
		lock, err = store.OpenCredentials(f.health.CredentialsDB)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { lock.Close() })
	}
	rc, _, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
	if rc != 2 || !strings.Contains(errs, "PR recorded, session NOT closed") ||
		!strings.Contains(errs, "session revoke --episode ep1") || len(f.finishRecords) != 1 {
		t.Fatalf("%d %s", rc, errs)
	}
	if _, exists, err := lock.ResolveSession(boundedTestContext(t), f.credentialID); err != nil || !exists {
		t.Fatalf("revoke failure: %v %v", exists, err)
	}
}

func TestEpisodeFinishCredentialsMismatch(t *testing.T) {
	f := newFinishFixture(t)
	f.health.CredentialsDB = filepath.Join(f.root, "other.db")
	rc, _, errs := f.run(testMintEnv("y\n", 1), "--include-unattributed")
	if rc != 1 || !strings.Contains(errs, f.health.CredentialsDB) ||
		!strings.Contains(errs, filepath.Join(f.root, "credentials.db")) || len(f.seam.calls) != 0 || len(f.scopes) != 0 {
		t.Fatalf("%d %s", rc, errs)
	}
}
