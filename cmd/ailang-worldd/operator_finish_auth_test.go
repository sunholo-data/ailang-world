package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/worktree"
)

const fixtureFinishRemote = "https://github.com/fixture/approved.git"

func writeFinishTestTool(t *testing.T, c *runnerConfig, dir, name, text string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+text), 0700); err != nil {
		t.Fatal(err)
	}
	c.Executables[path] = serviceFileHash(path)
}

func TestOperatorFinishAuth(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bare := filepath.Join(dir, "bare")
	ctx := boundedTestContext(t)
	gitRun := func(args ...string) {
		t.Helper()
		if _, err := worktree.RunIsolated(ctx, git, args, []string{"HOME=" + dir, "PATH=" + runnerSystemPath,
			"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}, nil); err != nil {
			t.Fatal(err)
		}
	}
	gitRun("-c", "init.defaultBranch=main", "init", "--bare", "-q", bare)
	gitRun("--git-dir="+bare, "--work-tree="+dir, "-c", "core.bare=false", "-c", "user.name=Fixture", "-c", "user.email=test@invalid",
		"commit", "--allow-empty", "-qm", "fixture")
	t.Setenv("GH_TOKEN", "private-fixture-token")
	t.Setenv("GIT_CONFIG_COUNT", "forbidden")
	t.Setenv("AILANG_REGISTRY_API_KEY", "registry-fixture-token")
	for _, tc := range []struct {
		name, remote, ghFailure, permission, gitResult string
		step                                           operatorGitHubAuth
	}{
		{"missing-unpinned-tool", fixtureFinishRemote, "", "WRITE", "", 2},
		{"login", fixtureFinishRemote, "auth", "WRITE", "", 3},
		{"label-read", fixtureFinishRemote, "label", "WRITE", "", 4},
		{"repo-read", fixtureFinishRemote, "repo", "WRITE", "", 5},
		{"permission", fixtureFinishRemote, "", "READ", "", 5},
		{"https-helper", fixtureFinishRemote, "", "WRITE", "exit 1", 6},
		{"ssh-read", "git@github.com:fixture/approved.git", "", "WRITE", "exit 1", 6},
		{"empty", fixtureFinishRemote, "", "WRITE", "exit 0", 6},
		{"write", fixtureFinishRemote, "", "WRITE", "", 0},
		{"maintain", fixtureFinishRemote, "", "MAINTAIN", "", 0},
		{"keychain", fixtureFinishRemote, "", "WRITE", "", 0},
		{"admin", "ssh://git@github.com/fixture/approved.git", "", "ADMIN", "", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tools := filepath.Join(dir, tc.name)
			if err := os.Mkdir(tools, 0700); err != nil {
				t.Fatal(err)
			}
			trace := filepath.Join(tools, "trace")
			c := runnerConfig{Home: dir, UID: os.Getuid(), ToolDirs: []string{tools}, Executables: map[string]string{}}
			envGuard := `[ "${GH_TOKEN-unset}|${GIT_CONFIG_COUNT-unset}|${AILANG_REGISTRY_API_KEY-unset}" = unset\|unset\|unset ]
[ "$GH_PROMPT_DISABLED|$GIT_TERMINAL_PROMPT" = 1\|0 ]
`
			writeFinishTestTool(t, &c, tools, "gh", envGuard+fmt.Sprintf(`
case "$*" in
  'auth status --hostname github.com') ;;
  'label list -R fixture/approved --search ailang-world --json name') ;;
  'repo view fixture/approved --json viewerPermission --jq .viewerPermission') ;;
  *) printf forbidden >> %s; exit 9;;
esac
printf '%%s ' "$1" >> %s
if [ "$1" = repo ]; then printf '%%s\n' %s; fi
[ "$1" != %s ]
`, shellQuote(trace), shellQuote(trace), shellQuote(tc.permission), shellQuote(tc.ghFailure)))
			if tc.name == "keychain" {
				c.FinishHTTPSHelper = "osxkeychain"
				writeFinishTestTool(t, &c, tools, "git-credential-osxkeychain", "exit 1\n")
			}
			gitArgs := "--git-dir=/dev/null -c credential.interactive=never -c credential.helper= "
			if strings.HasPrefix(tc.remote, "https://") {
				gitArgs += "-c credential.helper=!" + shellQuote(filepath.Join(tools, "gh")) + " auth git-credential"
				if tc.name == "keychain" {
					gitArgs = "--git-dir=/dev/null -c credential.interactive=never -c credential.helper= -c credential.helper=!" +
						shellQuote(filepath.Join(tools, "git-credential-osxkeychain"))
				}
			} else {
				gitArgs += "-c core.sshCommand=/usr/bin/ssh -F /dev/null -o BatchMode=yes -o StrictHostKeyChecking=yes"
			}
			gitArgs += " ls-remote --symref -- " + tc.remote + " HEAD"
			if tc.gitResult == "" {
				tc.gitResult = "exec " + shellQuote(git) + " --git-dir=/dev/null ls-remote --symref -- " + shellQuote(bare) + " HEAD"
			}
			writeFinishTestTool(t, &c, tools, "git", envGuard+fmt.Sprintf(`
[ "$GIT_CONFIG_NOSYSTEM|$GIT_CONFIG_GLOBAL" = 1\|/dev/null ]
[ "${SSH_AUTH_SOCK-unset}" = unset ]
[ "$*" = %s ] || { printf forbidden >> %s; exit 9; }
printf 'git ' >> %s
%s
`, shellQuote(gitArgs), shellQuote(trace), shellQuote(trace), tc.gitResult))
			if tc.step == 2 {
				delete(c.Executables, filepath.Join(tools, "gh"))
			}
			err := prepareOperatorFinish(boundedTestContext(t), c, tc.remote)
			if tc.step == 0 && err != nil || tc.step != 0 && (!errors.Is(err, tc.step) || !strings.Contains(err.Error(), "blocked: github-auth")) {
				t.Fatalf("want readiness %d, got %v", tc.step, err)
			}
			want := map[operatorGitHubAuth]string{2: "", 3: "auth ", 4: "auth label ", 5: "auth label repo "}[tc.step]
			if tc.step == 0 || tc.step == 6 {
				want = "auth label repo git "
			}
			if got, err := os.ReadFile(trace); string(got) != want || err != nil && !(tc.step == 2 && errors.Is(err, os.ErrNotExist)) {
				t.Fatalf("child calls %q (%v), want %q", got, err, want)
			}
		})
	}
}

func TestOperatorFinishAuthInstall(t *testing.T) {
	f := operatorService(t)
	f.config.FinishRemotes = []string{fixtureFinishRemote}
	tools := f.config.ToolDirs[0]
	trace := filepath.Join(f.home, "readiness")
	writeFinishTestTool(t, &f.config, tools, "gh", "printf ready > "+shellQuote(trace)+"\nexit 1\n")
	writeFinishTestTool(t, &f.config, tools, "git", "exit 99\n")
	observe := f.env.observe
	f.env.observe = func(path string) {
		got, err := os.ReadFile(trace)
		if string(got) != "ready" || err != nil {
			t.Fatalf("installation wrote %s before auth was observed", path)
		}
		observe(path)
	}
	if rc, out := f.install(t); rc != exitOK || !strings.Contains(out, "finish blocked: github-auth (readiness step 3)") {
		t.Fatalf("install %d: %s", rc, out)
	}
	if err := os.Remove(filepath.Join(tools, "gh")); err != nil {
		t.Fatal(err)
	}
	if rc, out := f.command("start"); rc != exitOK {
		t.Fatalf("start after finish auth failure %d: %s", rc, out)
	}
	if rc, out := f.command("status"); rc != exitOK || !strings.Contains(out, "start: ready") || !strings.Contains(out, "finish blocked: github-auth") {
		t.Fatalf("status %d: %s", rc, out)
	}
}

func TestOperatorFinishToolPins(t *testing.T) {
	f := operatorService(t)
	f.config.FinishRemotes = []string{fixtureFinishRemote}
	tools := f.config.ToolDirs[0]
	trace := filepath.Join(f.home, "unexpected-readiness")
	for _, name := range []string{"gh", "git"} {
		writeFinishTestTool(t, &f.config, tools, name, "exit 17\n")
	}
	if err := f.env.checkRunnerConfig(f.config); err != nil {
		t.Fatalf("pinned configuration: %v", err)
	}
	for _, path := range []string{filepath.Join(tools, "git"), f.env.paths.Binary, filepath.Join(tools, "gh")} {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("#!/bin/sh\ntouch "+shellQuote(trace)+"\nexit 1\n"), 0700); err != nil {
			t.Fatal(err)
		}
		if filepath.Base(path) == "gh" {
			if err := f.env.checkRunnerConfig(f.config); err != nil {
				t.Fatalf("finish-only gh must leave start available: %v", err)
			}
			if err := prepareOperatorFinish(boundedTestContext(t), f.config, fixtureFinishRemote); !errors.Is(err, operatorGitHubAuth(2)) {
				t.Fatalf("changed readiness executable was accepted: %v", err)
			}
		} else if err := f.env.checkRunnerConfig(f.config); err == nil || !strings.Contains(err.Error(), "executable changed") {
			t.Fatalf("changed shared executable %s: %v", path, err)
		}
		if err := os.WriteFile(path, b, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a changed tool executed or sentinel could not be read: %v", err)
	}
}
