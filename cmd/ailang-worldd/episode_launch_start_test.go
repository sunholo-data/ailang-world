package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

func TestStartLaunchValidatesBeforeFence(t *testing.T) {
	f := newPresetFixture(t)
	bin := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := os.WriteFile(filepath.Join(bin, "codex"), []byte("#!/bin/sh\nif [ \"$1\" = --version ]; then echo fixture; exit 0; fi\nexit 99\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	var p launchPreset
	_ = json.Unmarshal(f.data, &p)
	p.Launch = &presetLaunch{Mode: "interactive"}
	data, _ := json.Marshal(p)
	_ = os.WriteFile(f.file, data, 0o600)
	var opened bool
	var out, errOut bytes.Buffer
	beforeRows := f.r.f.rows(t)
	beforeTrees := presetGit(t, f.r.f, "worktree", "list", "--porcelain")
	if code := runEpisodeStart(f.args("--task", f.task, "--launch", "--hosts-dir", t.TempDir()), &out, &errOut, spyEnv(&opened)); code == 0 ||
		!strings.Contains(errOut.String(), "hosts init") || opened {
		t.Fatal(code, opened, errOut.String())
	}
	f.r.assertNoSessionFiles(t)
	if after := presetGit(t, f.r.f, "worktree", "list", "--porcelain"); beforeTrees != after {
		t.Fatal("worktree changed")
	}
	if len(f.r.f.rows(t)) != len(beforeRows) {
		t.Fatal("credentials changed")
	}
}

func TestStartLaunchHeadlessMintsAndLaunchesOnce(t *testing.T) {
	f := newPresetFixture(t)
	home, bin, facts := t.TempDir(), t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LAUNCH_FACTS", facts)
	script := "#!/bin/sh\nprintf 'started\\n' >> \"$LAUNCH_FACTS/started\"\nprintf '{\"type\":\"fixture\"}\\n'\n"
	if err := os.WriteFile(filepath.Join(bin, "codex"), []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	old := launchClient
	t.Cleanup(func() { launchClient = old })
	launchClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
			ep := "p-2610281352"
			switch r.URL.Path {
			case "/v1/episodes/" + ep + "/report":
				_ = json.NewEncoder(w).Encode(daemon.EpisodeReport{Start: readPresetStart(t, f, ep)})
			case "/v1/session/self":
				_ = json.NewEncoder(w).Encode(daemon.SessionSelf{Episode: ep, SecondsLeft: 1000})
			default:
				t.Error(r.URL.Path)
				w.WriteHeader(404)
			}
		})
	}
	term := &promptTerm{in: strings.NewReader("y\n")}
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--task", f.task, "--launch", "--headless"), &out, &errOut, presetEnv(term)); code != 0 {
		t.Fatal(code, out.String(), errOut.String())
	}
	data, err := os.ReadFile(filepath.Join(facts, "started"))
	if err != nil || string(data) != "started\n" {
		t.Fatal("launch count", string(data), err)
	}
	if len(f.r.f.rows(t)) != 1 || len(readPresetStart(t, f, "p-2610281352").Launches) != 1 {
		t.Fatal("mint/record missing")
	}
}

func TestStartLaunchFlagRefusals(t *testing.T) {
	for _, args := range [][]string{
		{"--launch"}, {"--headless"}, {"--launch", "--headless", "--preset", "se-tools"},
	} {
		var out, errOut bytes.Buffer
		var opened bool
		if code := runEpisodeStart(args, &out, &errOut, spyEnv(&opened)); code == 0 || opened ||
			!strings.Contains(errOut.String(), "why:") || !strings.Contains(errOut.String(), "next:") {
			t.Fatal(args, code, opened, errOut.String())
		}
	}
}

func TestStartLaunchPreFenceTaskBinaryAndDirectory(t *testing.T) {
	for _, which := range []string{"task", "binary", "directory"} {
		t.Run(which, func(t *testing.T) {
			f := newPresetFixture(t)
			bin := t.TempDir()
			t.Setenv("HOME", t.TempDir())
			t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
			if which != "binary" {
				_ = os.WriteFile(filepath.Join(bin, "codex"), []byte("#!/bin/sh\nexit 99\n"), 0o700)
			} else {
				_ = os.Symlink(f.r.f.git, filepath.Join(bin, "git"))
				t.Setenv("PATH", bin)
			}
			args := f.args("--launch", "--headless")
			want := "--task"
			if which != "task" {
				args = append(args, "--task", f.task)
				want = "missing from PATH"
			}
			if which == "directory" {
				dir := filepath.Join(filepath.Dir(f.r.sessions), "launch", "p-2610281352")
				_ = os.MkdirAll(dir, 0o700)
				_ = os.WriteFile(filepath.Join(dir, "keep"), []byte("keep"), 0o600)
				want = "is not empty"
			}
			var opened bool
			var out, errOut bytes.Buffer
			env := spyEnv(&opened)
			env.now = presetClock
			if code := runEpisodeStart(args, &out, &errOut, env); code == 0 || opened || !strings.Contains(errOut.String(), want) {
				t.Fatal(code, opened, errOut.String())
			}
			f.r.assertNoSessionFiles(t)
			if len(f.r.f.rows(t)) != 0 {
				t.Fatal("mint before validation")
			}
		})
	}
}
