package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func assertNoLaunchStageR2(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".ep.") || strings.Contains(entry.Name(), "stage") {
			t.Fatal("private staging copy remains", entry.Name())
		}
	}
}

func TestHostTranscriptRefusesSymlink(t *testing.T) {
	for _, agent := range []string{"claude", "codex"} {
		t.Run(agent, func(t *testing.T) {
			script := claudeSessionScript
			if agent == "codex" {
				script = codexSessionScript
			}
			script += "mv \"$source\" \"$source.real\"\nln -s \"$source.real\" \"$source\"\n"
			f, hosts, _ := interactiveFixture(t, agent, script)
			if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
				t.Fatal(code, errOut)
			}
			entry := f.read().Launches[0]
			if entry.TranscriptState != "missing" || !strings.Contains(entry.TranscriptReason, "not a regular file") {
				t.Fatal(entry)
			}
			if _, err := os.Stat(entry.Transcript); !os.IsNotExist(err) {
				t.Fatal("kept a symlink transcript", err)
			}
			// The copy boundary itself must refuse a leaf substituted after lookup.
			source := f.contents(filepath.Join(f.facts, "source"))
			stage, err := os.CreateTemp(t.TempDir(), "stage")
			if err != nil {
				t.Fatal(err)
			}
			defer stage.Close()
			if err := copyHostTranscript(hostTranscript{path: source}, stage); err == nil || !strings.Contains(err.Error(), "not a regular file") {
				t.Fatal("copy followed leaf", err)
			}
			assertNoLaunchStageR2(t, f.sessions)
		})
	}
}

func TestHostTranscriptStageRemoved(t *testing.T) {
	for _, kind := range []string{"kept", "scrubbed", "missing", "error"} {
		t.Run(kind, func(t *testing.T) {
			script := claudeSessionScript
			switch kind {
			case "scrubbed":
				script = strings.Replace(script, `printf '%s\n' '{"type":"fixture"}' > "$source"`, `printf '%s\n' "$WORLD_SESSION" > "$source"`, 1)
			case "missing":
				script = "echo no-transcript\n"
			case "error":
				script += "chmod 000 \"$source\"\n"
			}
			f, hosts, _ := interactiveFixture(t, "claude", script)
			if kind == "error" {
				// Ensure temp cleanup can remove the deliberately unreadable source.
				t.Cleanup(func() {
					b, _ := os.ReadFile(filepath.Join(f.facts, "source"))
					_ = os.Chmod(string(b), 0o600)
				})
			}
			if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
				t.Fatal(code, errOut)
			}
			entry := f.read().Launches[0]
			want := kind
			if kind == "error" {
				want = "missing"
			}
			if entry.TranscriptState != want {
				t.Fatal(entry)
			}
			assertNoLaunchStageR2(t, f.sessions)
		})
	}
}

func TestHostTranscriptSizeChild(t *testing.T) {
	if os.Getenv("WORLD_TRANSCRIPT_SIZE_CHILD") != "1" {
		return
	}
	data, err := os.ReadFile(filepath.Join(os.Getenv("LAUNCH_FACTS"), "source"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(string(data), os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate((64 << 20) + 1); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestHostTranscriptSizeCap(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("WORLD_TEST_BINARY", binary)
	t.Setenv("WORLD_TRANSCRIPT_SIZE_CHILD", "1")
	script := claudeSessionScript + "\"$WORLD_TEST_BINARY\" -test.run=^TestHostTranscriptSizeChild$\n"
	f, hosts, _ := interactiveFixture(t, "claude", script)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entry := f.read().Launches[0]
	if entry.TranscriptState != "missing" || !strings.Contains(entry.TranscriptReason, "larger than 64 MiB") {
		t.Fatal(entry)
	}
	if _, err := os.Stat(entry.Transcript); !os.IsNotExist(err) {
		t.Fatal("oversize transcript kept", err)
	}
	info, err := os.Stat(entry.HostTranscript)
	if err != nil || info.Size() != (64<<20)+1 {
		t.Fatal("sparse fixture not created", info, err)
	}
	assertNoLaunchStageR2(t, f.sessions)
}

func TestLaunchInteractiveWallCapIsSessionLeft(t *testing.T) {
	f, hosts, _ := interactiveFixture(t, "claude", "exec /bin/sleep 30\n")
	f.rewrite(func(s *store.EpisodeStart) { s.Preset.Bytes = `{"launch":{"mode":"interactive","deadline_s":60}}` })
	launchClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/session/self" {
				f.serve(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(daemon.SessionSelf{
				Episode:     "ep",
				Worktree:    filepath.Join(f.ws, "ep"),
				SecondsLeft: 2,
				Addr:        "http://" + defaultEpisodeAddr,
				Grants:      []daemon.SessionGrant{{Effect: "Workspace.Read", Scope: "worktree", Budget: 50, Left: 50}},
			})
		})
	}
	started := time.Now()
	if code, _, errOut := f.run("--hosts-dir", hosts); code == 0 {
		t.Fatal("sleep escaped session wall cap", errOut)
	}
	entry := f.read().Launches[0]
	elapsed := time.Since(started)
	// The cap is the session's 2 s, not the preset's 60 s or the fake's 30 s sleep; the
	// upper bound leaves room for a loaded grader (it failed at 8.6 s with a 5-8 s window).
	if entry.FinishReason != "timeout" || elapsed < 2*time.Second || elapsed > 20*time.Second {
		t.Fatal(entry, elapsed)
	}
	assertNoLaunchStageR2(t, f.sessions)
}
