package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func TestPresetStartReadParity(t *testing.T) {
	f := newPresetFixture(t)
	oldClient := launchClient
	t.Cleanup(func() { launchClient = oldClient })
	launchClient = func(string) *client {
		t.Fatal("a start without --launch attempted a launch")
		return nil
	}
	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
	task := readStartFixtureFile(t, f.task)
	before := observeProvision(t, f.r, creds)
	input := presetStartInput{
		episodes: []string{"ep"}, workspaceRoot: f.r.f.ws,
		overrides: []presetStartOverride{{"agent", "claude"}, {"branch", "pinned"}, {"ttl", "120"}},
		readTask:  func() (string, []byte, error) { return f.task, task, nil },
	}
	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
	if err != nil {
		t.Fatal(err)
	}
	req := read.mint
	if req.episode != "ep" || req.dir != filepath.Join(f.r.f.ws, "ep") || req.outPath != filepath.Join(f.r.sessions, "ep.session") ||
		read.taskCopy != filepath.Join(f.r.sessions, "ep.task.md") || read.mcp != filepath.Join(f.r.sessions, "ep.mcp.json") ||
		req.startPoint != f.base || req.branch != "pinned" || req.existingBranch || req.reuse || req.ttl != 120 ||
		req.start.Agent != "claude" || req.start.Model != "" || req.start.Base.SHA != f.base ||
		req.start.Preset.SHA256 != serviceHash(f.data) || req.prof.Digest != f.profile.Digest ||
		req.start.Task == nil || req.start.Task.SHA256 != serviceHash(task) || req.start.Task.Title != "Fix fixture" ||
		read.health.ToolRef == "" || read.health.CredentialsDB != creds || !bytes.Equal(read.task, task) ||
		!reflect.DeepEqual(observeProvision(t, f.r, creds), before) {
		t.Fatalf("read seam diverged or wrote state: %+v", req.start)
	}
	if _, err := os.Lstat(operatorEpisodeLockPath(creds, "ep")); !os.IsNotExist(err) {
		t.Fatalf("read seam touched episode lock: %v", err)
	}

	term := &promptTerm{in: strings.NewReader("y\n")}
	env := presetEnv(term)
	env.openTerminal = func() (io.ReadWriteCloser, error) {
		// Resolution has completed. Only the already-read bytes may reach
		// the task copy, start row or prompt.
		seedStartFixtureFile(t, f.task, []byte("changed after resolution"))
		return term, nil
	}
	var out bytes.Buffer
	if code := runEpisodeStart(f.args("ep", "--agent", "claude", "--branch", "pinned", "--ttl", "120", "--task", f.task), &out, &out, env); code != exitOK {
		t.Fatal(code, out.String())
	}
	if !strings.HasPrefix(term.out.String(), req.presetPrompt) ||
		!strings.Contains(req.presetPrompt, "model <default for claude>") ||
		!strings.Contains(req.presetPrompt, "task      "+f.task+" \"Fix fixture\" ("+serviceHash(task)+")") ||
		!bytes.Equal(readStartFixtureFile(t, read.taskCopy), task) {
		t.Fatal("the attended path did not use the resolved review and bytes")
	}
	got := readPresetStart(t, f, "ep")
	want := *req.start
	dir := filepath.Join(f.r.f.ws, "ep")
	head, err := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "rev-parse", "HEAD"})
	rows := f.r.f.rows(t)
	token := readStartFixtureFile(t, req.outPath)
	mcp := readStartFixtureFile(t, read.mcp)
	if len(rows) != 1 || got == nil || err != nil || strings.TrimSpace(head) != f.base || !isHex64(string(token)) ||
		!bytes.Contains(mcp, append([]byte("Bearer "), token...)) {
		t.Fatalf("worktree, credentials or files differ from the resolved plan: rows=%d, HEAD=%s, err=%v", len(rows), head, err)
	}
	if !reflect.DeepEqual(got, &want) || rows[0].ExpiresAt != presetClock()+120 || rows[0].ProfileDigest != f.profile.Digest {
		t.Fatalf("mint changed reviewed start or binding: %+v", got)
	}
}

func TestPresetStartReadHeadless(t *testing.T) {
	f := newPresetFixture(t)
	bin := startReadFakeAgent(t, f)
	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
	before := observeProvision(t, f.r, creds)
	task := []byte("# Literal\n$(touch should-never-run)\n")
	input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws,
		launching: true, headless: true, hostsDir: t.TempDir(),
		readTask: func() (string, []byte, error) { return "inline", task, nil },
	}
	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
	if err != nil {
		t.Fatal(err)
	}
	binary, err := filepath.EvalSymlinks(filepath.Join(bin, "codex"))
	if err != nil {
		t.Fatal(err)
	}
	if read.launchMode != "headless" || read.launchWall != time.Hour || read.launchBinary != binary ||
		read.mint.start.Task.Path != "inline" || read.mint.start.Task.SHA256 != serviceHash(task) ||
		!bytes.Equal(read.task, task) || read.mcp != "" || !reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
		t.Fatal("typed resolver lost headless inputs or changed durable state")
	}
	if _, err := os.Lstat(filepath.Join(bin, "ran")); !os.IsNotExist(err) {
		t.Fatal("a read started the executable")
	}
}

func TestPresetStartReadInteractive(t *testing.T) {
	host, hosts, _ := interactiveFixture(t, "claude", "exit 99")
	f := newPresetFixture(t)
	t.Setenv("HOME", filepath.Join(host.root, "home"))
	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
	before := observeProvision(t, f.r, creds)
	input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws, launching: true, hostsDir: hosts,
		overrides: []presetStartOverride{{"agent", "claude"}, {"model", "explicit"}}}
	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
	if err != nil {
		t.Fatal(err)
	}
	if read.launchMode != "interactive" || read.launchWall != 0 || read.mcp != "" || read.mint.start.Task != nil ||
		read.mint.start.Model != "explicit" || !reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
		t.Fatal("interactive resolution changed mint semantics or durable state")
	}
	term := &promptTerm{in: strings.NewReader("n\n")}
	var out bytes.Buffer
	code := runEpisodeStart(f.args("ep", "--launch", "--hosts-dir", hosts, "--agent", "claude", "--model", "explicit"), &out, &out, presetEnv(term))
	if code != exitUsage || !strings.Contains(out.String(), "aborted (not confirmed)") ||
		!strings.Contains(term.out.String(), "agent     claude explicit") ||
		!reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
		t.Fatalf("terminal interactive validation drifted: %d %s", code, &out)
	}
	if _, err := os.Lstat(filepath.Join(host.facts, "started")); !os.IsNotExist(err) {
		t.Fatal("interactive resolution started a child")
	}
}

func TestPresetStartReadRefusals(t *testing.T) {
	for _, tc := range []struct{ name, want string }{
		{"branch", "derived branch already exists"},
		{"checked_out", "checked out in another worktree"},
		{"worktree", "already has a worktree"},
		{"token", ".session already exists; choose another episode"},
		{"task", ".task.md already exists"},
		{"mcp", ".mcp.json already exists"},
		{"loop", syscall.ELOOP.Error()},
		{"project_race", "exec profile changed during validation"},
		{"task_directory", "task must be a regular file"},
		{"unreadable_task", "no such file"},
		{"bad_ttl", "ttl must be non-negative"},
		{"bad_agent", "agent must be"},
		{"two_episodes", "expected at most one"},
		{"bad_episode", "does not match"},
		{"bad_branch", "check-ref-format"},
		{"health_missing", "fixture daemon offline"},
		{"world", "(F6)"},
		{"build", "is build"},
		{"tool", "serves the tool binary"},
		{"credentials", "which it never reads"},
		{"profile", "changed since the daemon"},
		{"invalid_wall", "launch deadline is too large"},
		{"launch_lock", "only one launch may run"},
		{"nonprivate_sessions", "must be private"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newPresetFixture(t)
			h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
			h.DBPath = f.r.f.db
			if _, err := store.UpgradeCredentials(h.CredentialsDB); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(f.r.sessions, 0700); err != nil {
				t.Fatal(err)
			}
			args := f.args("ep", "--agent", "claude", "--task", f.task)
			var healthErr error
			healthCalls, wantHealth := 0, 0
			switch tc.name {
			case "branch":
				presetGit(t, f.r.f, "branch", "world/ep", "dev")
			case "checked_out":
				presetGit(t, f.r.f, "branch", "world/used", "dev")
				presetGit(t, f.r.f, "worktree", "add", filepath.Join(t.TempDir(), "used"), "world/used")
				args = append(args, "--branch", "world/used")
			case "worktree":
				presetGit(t, f.r.f, "worktree", "add", "--detach", filepath.Join(f.r.f.ws, "ep"), "dev")
			case "token":
				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.session"), []byte("keep"))
			case "task":
				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.task.md"), []byte("keep"))
			case "mcp":
				target := filepath.Join(t.TempDir(), "missing")
				if err := os.Symlink(target, filepath.Join(f.r.sessions, "ep.mcp.json")); err != nil {
					t.Fatal(err)
				}
			case "loop":
				link := filepath.Join(filepath.Dir(f.r.sessions), "loop")
				if err := os.Symlink(link, link); err != nil {
					t.Fatal(err)
				}
				var p launchPreset
				if err := json.Unmarshal(f.data, &p); err != nil {
					t.Fatal(err)
				}
				p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: f.r.f.db, CredentialsDB: h.CredentialsDB,
					WorkspaceRoot: f.r.f.ws, SessionsDir: filepath.Join(link, "sessions")})
				changed, _ := json.Marshal(p)
				seedStartFixtureFile(t, f.file, changed)
			case "project_race":
				profile := filepath.Join(f.r.profiles, "proj.json")
				raw := readStartFixtureFile(t, profile)
				input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws, readTask: func() (string, []byte, error) {
					seedStartFixtureFile(t, profile, append(raw, '\n'))
					return f.task, nil, nil
				}}
				before := observeProvision(t, f.r, h.CredentialsDB)
				if read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock); read != nil ||
					err == nil || !strings.Contains(err.Error(), tc.want) || !reflect.DeepEqual(observeProvision(t, f.r, h.CredentialsDB), before) {
					t.Fatalf("profile drift allowed or wrote state: %v", err)
				}
				return
			case "task_directory":
				args = append(args, "--task", filepath.Dir(f.task))
			case "unreadable_task":
				args = append(args, "--task", f.task+"-missing")
			case "bad_ttl":
				args = append(args, "--ttl", "-1")
			case "bad_agent":
				args = append(args, "--agent", "other")
			case "two_episodes":
				args = append(args, "other")
			case "bad_episode":
				// Enter launch preflight before prepareMint can repeat the
				// same grammar check. The resolver must reject the episode
				// even though the requested agent would also be missing.
				bin := t.TempDir()
				if err := os.Symlink(f.r.f.git, filepath.Join(bin, "git")); err != nil {
					t.Fatal(err)
				}
				t.Setenv("PATH", bin)
				args = f.args("BAD", "--task", f.task, "--launch", "--headless")
			case "bad_branch":
				args = append(args, "--branch", "bad branch")
			case "health_missing":
				healthErr, wantHealth = errors.New("fixture daemon offline"), 1
			case "world":
				h.DBPath, wantHealth = filepath.Join(t.TempDir(), "other.db"), 1
			case "build":
				h.Build, wantHealth = daemon.BuildIdentity{}, 1
			case "tool":
				h.ToolRef, wantHealth = "sha256:other", 1
			case "credentials":
				h.CredentialsDB, wantHealth = filepath.Join(t.TempDir(), "other.db"), 1
			case "profile":
				h.ExecProfiles[0].Digest, wantHealth = "sha256:other", 1
			case "invalid_wall":
				var p launchPreset
				if err := json.Unmarshal(f.data, &p); err != nil {
					t.Fatal(err)
				}
				p.Launch = &presetLaunch{Mode: "headless", DeadlineS: 1e12}
				data, _ := json.Marshal(p)
				seedStartFixtureFile(t, f.file, data)
				args = append(args, "--launch", "--headless")
			case "launch_lock":
				startReadFakeAgent(t, f)
				args = append(args, "--agent", "codex", "--launch", "--headless")
				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.launch.lock"), []byte("0"))
			case "nonprivate_sessions":
				if err := os.Chmod(f.r.sessions, 0750); err != nil {
					t.Fatal(err)
				}
			}
			stubDaemon(t, h, healthErr)
			fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) {
				healthCalls++
				return h, healthErr
			}
			creds, err := credentialsFor(f.r.f.db, "")
			if err != nil {
				t.Fatal(err)
			}
			before := observeProvision(t, f.r, creds)
			var opened bool
			var out bytes.Buffer
			code := runEpisodeStart(args, &out, &out, spyEnv(&opened))
			if code != exitUsage || !strings.Contains(out.String(), tc.want) || opened || healthCalls != wantHealth ||
				!reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
				t.Fatalf("read refusal changed state, effect count or class: code=%d opened=%v health=%d %s", code, opened, healthCalls, &out)
			}
		})
	}
}

func startReadFakeAgent(t *testing.T, f presetFixture) string {
	t.Helper()
	bin := t.TempDir()
	gitLink := filepath.Join(bin, "git")
	if err := os.Symlink(f.r.f.git, gitLink); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(bin, "codex")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\n: > \"$WORLD_FIXTURE_AGENT_RAN\"\nexit 99\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	t.Setenv("WORLD_FIXTURE_AGENT_RAN", filepath.Join(bin, "ran"))
	return bin
}

func seedStartFixtureFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func readStartFixtureFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
