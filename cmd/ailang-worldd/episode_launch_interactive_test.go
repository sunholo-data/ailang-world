package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

func interactiveFixture(t *testing.T, agent, output string) (*launchFixture, string, string) {
	t.Helper()
	f := newLaunchFixture(t, agent)
	f.rewrite(func(s *store.EpisodeStart) { s.Preset.Bytes = `{"launch":{"mode":"interactive"}}` })
	oldWorldd := hostsWorlddPath
	t.Cleanup(func() { hostsWorlddPath = oldWorldd })
	worldd := filepath.Join(f.bin, "fixture-worldd")
	if err := os.WriteFile(worldd, []byte("#!/bin/sh\nexit 1\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	hostsWorlddPath = func() string { return worldd }
	hosts := filepath.Join(f.root, "hosts")
	for _, name := range []string{"claude", "codex"} {
		script := "#!/bin/sh\nif [ \"$1\" = --version ]; then echo fixture; exit 0; fi\n" +
			"pwd -P > \"$LAUNCH_FACTS/cwd\"\n/usr/bin/env > \"$LAUNCH_FACTS/env\"\nprintf '%s\\000' \"$@\" > \"$LAUNCH_FACTS/argv\"\necho started >> \"$LAUNCH_FACTS/started\"\n" + output
		if err := os.WriteFile(filepath.Join(f.bin, name), []byte(script), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	var out strings.Builder
	if code := runHosts([]string{"init", "--out-dir", hosts, "--addr", "http://" + defaultEpisodeAddr}, &out, &out); code != 0 {
		t.Fatal(code, out.String())
	}
	tty := filepath.Join(t.TempDir(), "tty")
	old := openLaunchTerminal
	t.Cleanup(func() { openLaunchTerminal = old })
	openLaunchTerminal = func() (*os.File, error) { return os.OpenFile(tty, os.O_RDWR|os.O_CREATE, 0o600) }
	return f, hosts, tty
}

const claudeSessionScript = `prev=''
for a in "$@"; do
  if [ "$prev" = --session-id ]; then id="$a"; fi
  prev="$a"
done
slug=$(pwd -P | LC_ALL=C sed 's/[^A-Za-z0-9]/-/g')
root="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"
mkdir -p "$root/projects/$slug"
source="$root/projects/$slug/$id.jsonl"
printf '%s\n' '{"type":"fixture"}' > "$source"
printf '%s' "$source" > "$LAUNCH_FACTS/source"
echo hello
`

const codexSessionScript = `id=12345678-1234-4234-8234-123456789abc
day="$CODEX_HOME/sessions/$(date +%Y/%m/%d)"
mkdir -p "$day"
source="$day/rollout-$(date +%Y-%m-%dT%H-%M-%S)-$id.jsonl"
printf '{"type":"session_meta","payload":{"id":"%s","cwd":"%s"}}\n' "$id" "$(pwd -P)" > "$source"
printf '%s' "$source" > "$LAUNCH_FACTS/source"
decoy="$day/rollout-$(date +%Y-%m-%dT%H-%M-%S)-other.jsonl"
printf '{"type":"session_meta","payload":{"id":"other","cwd":"/another/cwd"}}\n' > "$decoy"
# A previous-minute matching cwd cannot belong to this launch.
oldstamp=$(date -v-1M +%Y-%m-%dT%H-%M-%S 2>/dev/null || date -d '1 minute ago' +%Y-%m-%dT%H-%M-%S)
old="$day/rollout-$oldstamp-00000000-0000-4000-8000-000000000000.jsonl"
printf '{"type":"session_meta","payload":{"id":"old","cwd":"%s"}}\n' "$(pwd -P)" > "$old"
echo hello
`

func TestLaunchInteractiveClaudeAC10(t *testing.T) {
	f, hosts, tty := interactiveFixture(t, "claude", claudeSessionScript)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entry := f.read().Launches[0]
	source := f.contents(filepath.Join(f.facts, "source"))
	got := f.contents(filepath.Join(f.sessions, "ep.transcript.jsonl"))
	if got != f.contents(source) || entry.TranscriptSHA256 != serviceHash([]byte(got)) || entry.TranscriptState != "kept" || entry.Mode != "interactive" || entry.HostTranscript != source || entry.MCPConfigDeleted != nil {
		t.Fatal(entry, got)
	}
	if info, err := os.Stat(entry.Transcript); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatal(info, err)
	}
	args := strings.Split(strings.TrimSuffix(f.contents(filepath.Join(f.facts, "argv")), "\x00"), "\x00")
	value := func(flag string) string {
		for i, a := range args {
			if a == flag && i+1 < len(args) {
				return args[i+1]
			}
		}
		return ""
	}
	grants := []broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 50}}
	if value("--session-id") != entry.HostSessionID || len(entry.HostSessionID) != 36 || entry.HostSessionID[14] != '4' || value("--allowedTools") != strings.Join(allowedTools(grants), ",") || value("--append-system-prompt") == "" || !strings.HasPrefix(args[len(args)-1], "# AC9 fixture") {
		t.Fatal(args, entry)
	}
	env := f.contents(filepath.Join(f.facts, "env"))
	if !strings.Contains(env, "WORLD_SESSION="+f.token+"\n") || f.contents(tty) != "hello\n" || strings.Contains(got, "hello") {
		t.Fatal(env, got)
	}
	if _, err := os.Lstat(filepath.Join(f.sessions, "ep.mcp.json")); !os.IsNotExist(err) {
		t.Fatal("wrote MCP config", err)
	}
	if entry.Host != filepath.Join(hosts, "bin", "claude-world") || entry.ArgvRedacted[0] != entry.Host {
		t.Fatal(entry)
	}
	data, _ := json.Marshal(entry)
	t.Log("claude launch: " + string(data))
	t.Logf("claude argv: %q", entry.ArgvRedacted)
}

func TestLaunchInteractiveCodexAC10(t *testing.T) {
	f, hosts, tty := interactiveFixture(t, "codex", codexSessionScript)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entry := f.read().Launches[0]
	source := f.contents(filepath.Join(f.facts, "source"))
	if entry.TranscriptState != "kept" || entry.HostTranscript != source || f.contents(entry.Transcript) != f.contents(source) || entry.HostSessionID != "12345678-1234-4234-8234-123456789abc" || f.contents(tty) != "hello\n" {
		t.Fatal(entry)
	}
	env := f.contents(filepath.Join(f.facts, "env"))
	if !strings.Contains(env, "CODEX_HOME="+filepath.Join(hosts, "codex")+"\n") || !strings.Contains(env, "WORLD_SESSION="+f.token+"\n") {
		t.Fatal(env)
	}
	args := strings.Split(f.contents(filepath.Join(f.facts, "argv")), "\x00")
	if len(args) < 8 || args[2] != "-m" || args[3] != "fixture-model" || args[4] != "-c" || !strings.HasPrefix(args[5], "developer_instructions=") || !strings.HasPrefix(args[6], "# AC9 fixture") {
		t.Fatal(args)
	}
	t.Logf("codex argv: %q", entry.ArgvRedacted)
}

func TestLaunchInteractiveScrub(t *testing.T) {
	script := strings.Replace(claudeSessionScript, `printf '%s\n' '{"type":"fixture"}' > "$source"`, `printf '%s\n' "$WORLD_SESSION" > "$source"`, 1)
	f, hosts, _ := interactiveFixture(t, "claude", script)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entry := f.read().Launches[0]
	if entry.TranscriptState != "scrubbed" || strings.Contains(f.contents(entry.Transcript), f.token) || strings.Count(f.contents(entry.Transcript), "\n") != 1 {
		t.Fatal(entry)
	}
	if !strings.Contains(f.contents(entry.HostTranscript), f.token) {
		t.Fatal("host source changed")
	}
}

func TestLaunchInteractiveMissingAndNoTask(t *testing.T) {
	f, hosts, _ := interactiveFixture(t, "codex", "echo hello\n")
	f.rewrite(func(s *store.EpisodeStart) { s.Task = nil; s.Model = "" })
	path := filepath.Join(f.sessions, "ep.transcript.jsonl")
	_ = os.WriteFile(path, []byte("previous\n"), 0o600)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entry := f.read().Launches[0]
	if entry.TranscriptState != "missing" || !strings.Contains(entry.TranscriptReason, filepath.Join(hosts, "codex", "sessions")) || f.contents(path) != "previous\n" || len(entry.ArgvRedacted) != 3 {
		t.Fatal(entry)
	}
}

func TestLaunchInteractiveSecondAfterHeadless(t *testing.T) {
	f, hosts, _ := interactiveFixture(t, "claude", claudeSessionScript)
	if code, _, errOut := f.run("--headless"); code != 0 {
		t.Fatal(code, errOut)
	}
	first := f.read().Launches[0]
	before := f.contents(first.Transcript)
	if code, _, errOut := f.run("--hosts-dir", hosts); code != 0 {
		t.Fatal(code, errOut)
	}
	entries := f.read().Launches
	if len(entries) != 2 || !reflect.DeepEqual(first, entries[0]) || !strings.HasPrefix(f.contents(first.Transcript), before) || !strings.Contains(f.contents(first.Transcript), `{"world_launch": 2}`) {
		t.Fatal(entries)
	}
}

func TestLaunchInteractiveRefusals(t *testing.T) {
	for _, kind := range []string{"terminal", "missing", "stale", "address"} {
		t.Run(kind, func(t *testing.T) {
			f, hosts, _ := interactiveFixture(t, "claude", claudeSessionScript)
			want := "hosts init"
			switch kind {
			case "terminal":
				openLaunchTerminal = func() (*os.File, error) { return nil, fmt.Errorf("fixture tty absent") }
				want = "interactive launch needs a terminal"
			case "missing":
				hosts = t.TempDir()
			case "stale":
				_ = os.WriteFile(filepath.Join(hosts, "bin", "claude-world"), []byte("#!/bin/sh\nexit 99\n"), 0o700)
			case "address":
				// Build a complete, internally consistent launcher for another daemon.
				hosts = filepath.Join(t.TempDir(), "other-hosts")
				var out strings.Builder
				if code := runHosts([]string{"init", "--out-dir", hosts, "--addr", "http://localhost:" + strings.Split(defaultEpisodeAddr, ":")[1]}, &out, &out); code != 0 {
					t.Fatal(code, out.String())
				}
				want = "differs from preset address"
			}
			before := launchFileSnapshot(t, f.root)
			if code, _, errOut := f.run("--hosts-dir", hosts); code == 0 || !strings.Contains(errOut, want) {
				t.Fatal(code, errOut)
			}
			if !reflect.DeepEqual(before, launchFileSnapshot(t, f.root)) {
				t.Fatal("refusal changed files")
			}
			if _, err := os.Stat(filepath.Join(f.facts, "started")); !os.IsNotExist(err) {
				t.Fatal("agent started", err)
			}
		})
	}
}

func TestHostTranscriptPathClaude(t *testing.T) {
	if got := claudeTranscriptSlug("/private/var/x/.ailang/a_b"); got != "-private-var-x--ailang-a-b" {
		t.Fatal(got)
	}
	dir := t.TempDir()
	real, _ := filepath.EvalSymlinks(dir)
	root := t.TempDir()
	id := "fixture"
	path := filepath.Join(root, "projects", claudeTranscriptSlug(real), id+".jsonl")
	_ = os.MkdirAll(filepath.Dir(path), 0o700)
	_ = os.WriteFile(path, []byte("{}"), 0o600)
	got := claudeTranscriptPath(dir, id, []string{"HOME=/unused", "CLAUDE_CONFIG_DIR=" + root})
	if got.path != path {
		t.Fatal(got)
	}
	long := filepath.Join(dir, strings.Repeat("a", 100), strings.Repeat("b", 100))
	_ = os.MkdirAll(long, 0o700)
	real, _ = filepath.EvalSymlinks(long)
	slug := claudeTranscriptSlug(real)
	path = filepath.Join(root, "projects", slug[:200]+"-hosthash", id+".jsonl")
	_ = os.MkdirAll(filepath.Dir(path), 0o700)
	_ = os.WriteFile(path, []byte("{}"), 0o600)
	if got := claudeTranscriptPath(long, id, []string{"CLAUDE_CONFIG_DIR=" + root}); got.path != path {
		t.Fatal(got)
	}
}

func TestHostTranscriptPathCodexMidnightAndAmbiguity(t *testing.T) {
	root, dir := t.TempDir(), t.TempDir()
	start := time.Date(2026, 10, 29, 23, 59, 59, 0, time.Local)
	end := start.Add(2 * time.Second)
	write := func(stamp time.Time, id, cwd string) string {
		path := filepath.Join(root, "codex", "sessions", stamp.Format("2006/01/02"), "rollout-"+stamp.Format("2006-01-02T15-04-05")+"-"+id+".jsonl")
		_ = os.MkdirAll(filepath.Dir(path), 0o700)
		data, _ := json.Marshal(map[string]any{"type": "session_meta", "payload": map[string]string{"id": id, "cwd": cwd}})
		_ = os.WriteFile(path, append(data, '\n'), 0o644)
		return path
	}
	path := write(end, "one", dir)
	write(start.Add(-time.Minute), "old", dir)
	write(start, "other", "/other")
	if got := codexTranscriptPath(root, dir, start, end); got.path != path || got.id != "one" {
		t.Fatal(got)
	}
	write(start, "two", dir)
	if got := codexTranscriptPath(root, dir, start, end); got.path != "" || !strings.Contains(got.reason, "2 rollouts") {
		t.Fatal(got)
	}
}
