package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeFakePreset(t *testing.T, dir, name, repo, agent string) {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"name":  name,
		"repo":  repo,
		"agent": agent,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".json"), b, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunWrapperBuildsEpisodeStartArgs(t *testing.T) {
	origStart := runEpisodeStartFunc
	origGetwd := runGetwdFunc
	defer func() {
		runEpisodeStartFunc = origStart
		runGetwdFunc = origGetwd
	}()

	tmpRoot := t.TempDir()
	presetsDir := filepath.Join(tmpRoot, "presets")
	repoWorld := filepath.Join(tmpRoot, "ailang-world")
	repoWorldSub := filepath.Join(repoWorld, "cmd", "ailang-worldd")
	repoDocparse := filepath.Join(tmpRoot, "docparse")
	repoOther := filepath.Join(tmpRoot, "unrelated")
	for _, d := range []string{presetsDir, repoWorldSub, repoDocparse, repoOther} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFakePreset(t, presetsDir, "world-agy", repoWorld, "agy")
	writeFakePreset(t, presetsDir, "world-agy-interactive", repoWorld, "agy")
	writeFakePreset(t, presetsDir, "world", repoWorld, "claude")
	writeFakePreset(t, presetsDir, "docparse", repoDocparse, "agy")

	var gotArgs []string
	var gotTaskBody string
	runEpisodeStartFunc = func(args []string, stdout, stderr io.Writer, env sessionEnv) int {
		gotArgs = append([]string(nil), args...)
		gotTaskBody = ""
		for i := 0; i+1 < len(args); i++ {
			if args[i] == "--task" {
				b, err := os.ReadFile(args[i+1])
				if err == nil {
					gotTaskBody = string(b)
				}
			}
		}
		return exitOK
	}
	runGetwdFunc = func() (string, error) { return repoOther, nil }

	// 1. Explicit preset + inline task -> headless + --from-allowance + --task <tmp>
	var out, errOut bytes.Buffer
	if rc := runWrapper([]string{"--presets-dir", presetsDir, "docparse", "fix", "the", "bug"}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("inline run failed rc=%d stderr=%s", rc, errOut.String())
	}
	if gotTaskBody != "fix the bug\n" {
		t.Fatalf("got task body %q, want %q", gotTaskBody, "fix the bug\n")
	}
	wantPrefix := []string{"--preset", "docparse", "--launch", "--headless", "--from-allowance", "--task"}
	if len(gotArgs) < 7 || !reflect.DeepEqual(gotArgs[:6], wantPrefix) {
		t.Fatalf("got args %q, want prefix %q", gotArgs, wantPrefix)
	}
	if _, err := os.Stat(gotArgs[6]); !os.IsNotExist(err) {
		t.Fatalf("expected temp task file %q to be cleaned up, err=%v", gotArgs[6], err)
	}

	// 2. Bare `run` inside repo subdirectory -> auto-detects shorter agy preset ("world-agy") + interactive!
	runGetwdFunc = func() (string, error) { return repoWorldSub, nil }
	gotArgs = nil
	out.Reset()
	errOut.Reset()
	if rc := runWrapper([]string{"--presets-dir", presetsDir}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("cwd bare run failed rc=%d stderr=%s", rc, errOut.String())
	}
	wantCwdInteractive := []string{"--preset", "world-agy", "--launch", "--interactive", "--from-allowance", "--presets-dir", presetsDir}
	if !reflect.DeepEqual(gotArgs, wantCwdInteractive) {
		t.Fatalf("cwd bare args = %q, want %q", gotArgs, wantCwdInteractive)
	}

	// 3. `run "fix bug"` inside repo -> auto-detects "world-agy" + headless inline task!
	gotArgs = nil
	out.Reset()
	errOut.Reset()
	if rc := runWrapper([]string{"--presets-dir", presetsDir, "fix", "bug"}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("cwd inline run failed rc=%d stderr=%s", rc, errOut.String())
	}
	if gotTaskBody != "fix bug\n" || gotArgs[1] != "world-agy" || gotArgs[3] != "--headless" {
		t.Fatalf("cwd inline got args=%q task=%q", gotArgs, gotTaskBody)
	}

	// 4. `run --agent claude --prompt --episode my-ep` inside repo -> prefers claude preset ("world")!
	gotArgs = nil
	out.Reset()
	errOut.Reset()
	if rc := runWrapper([]string{"--presets-dir", presetsDir, "--agent", "claude", "--prompt", "--episode", "my-ep"}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("cwd claude run failed rc=%d stderr=%s", rc, errOut.String())
	}
	wantClaude := []string{"--preset", "world", "--launch", "--interactive", "--agent", "claude", "--presets-dir", presetsDir, "my-ep"}
	if !reflect.DeepEqual(gotArgs, wantClaude) {
		t.Fatalf("cwd claude args = %q, want %q", gotArgs, wantClaude)
	}

	// 5. Refusals: unrelated dir with no preset, --headless without task, both --task and inline text, --interactive + --headless.
	runGetwdFunc = func() (string, error) { return repoOther, nil }
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"no preset in unrelated dir", []string{"--presets-dir", presetsDir}, "no launch preset"},
		{"headless without task", []string{"--presets-dir", presetsDir, "docparse", "--headless"}, "headless run requires a task"},
		{"both task and inline", []string{"--presets-dir", presetsDir, "docparse", "--task", "/tmp/x.md", "inline"}, "cannot pass both --task"},
		{"interactive and headless", []string{"--presets-dir", presetsDir, "docparse", "--interactive", "--headless"}, "mutually exclusive"},
	} {
		out.Reset()
		errOut.Reset()
		if rc := runWrapper(tc.args, &out, &errOut, sessionEnv{}); rc == exitOK {
			t.Errorf("%s: expected refusal", tc.name)
		} else if !strings.Contains(errOut.String(), tc.want) {
			t.Errorf("%s: stderr %q does not contain %q", tc.name, errOut.String(), tc.want)
		}
	}
}
