package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

type launchFixture struct {
	t                              *testing.T
	root, sessions, ws, db, creds  string
	presets, bin, facts, token     string
	start                          *store.EpisodeStart
	expired, revoked, missingStart bool
	selfEpisode                    string
}

func newLaunchFixture(t *testing.T, agent string) *launchFixture {
	t.Helper()
	f := &launchFixture{t: t, root: t.TempDir(), token: strings.Repeat("1", 64), selfEpisode: "ep"}
	f.sessions, f.ws = filepath.Join(f.root, "sessions"), filepath.Join(f.root, "workspace")
	f.presets, f.bin = filepath.Join(f.root, "presets"), filepath.Join(f.root, "bin")
	f.db, f.creds, f.facts = filepath.Join(f.root, "world.db"), filepath.Join(f.root, "credentials.db"), filepath.Join(f.root, "facts")
	for _, dir := range []string{f.sessions, f.ws, f.presets, f.bin, f.facts, filepath.Join(f.root, "home")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", filepath.Join(f.root, "home"))
	t.Setenv("PATH", f.bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LAUNCH_FACTS", f.facts)
	t.Setenv("WORLD_SESSION", "inherited-must-be-removed")
	p := launchPreset{Schema: "world/preset/v1", Name: "p", Agent: agent, Launch: &presetLaunch{Mode: "headless", DeadlineS: 5}}
	p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: f.db, CredentialsDB: f.creds, WorkspaceRoot: f.ws, SessionsDir: f.sessions})
	data, _ := json.Marshal(p)
	if err := os.WriteFile(filepath.Join(f.presets, "p.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	task := []byte("# AC9 fixture\nMake the fixture change.\n")
	f.start = &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: "ep", Agent: agent, Model: "fixture-model",
		Preset: store.StartPreset{Name: "p", Bytes: string(data)}, Branch: "world/ep", Task: &store.StartTask{SHA256: serviceHash(task)}}
	c, err := store.OpenCredentials(f.creds)
	if err != nil {
		t.Fatal(err)
	}
	err = c.MintSessionWithStart(presetTestContext(t), store.SessionRow{CredentialID: "fixture", EpisodeID: "ep", GrantsJSON: "[]", CreatedAt: 1, ExpiresAt: 9999999999}, f.start)
	_ = c.Close()
	if err != nil {
		t.Fatal(err)
	}
	for suffix, data := range map[string][]byte{".session": []byte(f.token), ".task.md": task} {
		if err := os.WriteFile(filepath.Join(f.sessions, "ep"+suffix), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	f.fake(agent, "printf '%s\\n' '{\"type\":\"fixture\"}'\n")
	oldClient, oldHealth := launchClient, fetchDaemonHealth
	t.Cleanup(func() { launchClient, fetchDaemonHealth = oldClient, oldHealth })
	launchClient = func(string) *client { return memoryDaemon(f.serve) }
	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) { return daemon.HealthResponse{DBPath: f.db}, nil }
	return f
}

func (f *launchFixture) fake(agent, output string) {
	f.t.Helper()
	script := "#!/bin/sh\n" +
		"pwd > \"$LAUNCH_FACTS/cwd\"\n" +
		"/bin/ls -A > \"$LAUNCH_FACTS/listing\"\n" +
		"/usr/bin/env > \"$LAUNCH_FACTS/env\"\n" +
		"printf '%s\\000' \"$@\" > \"$LAUNCH_FACTS/argv\"\n" +
		"printf 'started\\n' >> \"$LAUNCH_FACTS/started\"\n" +
		"prev=''\nfor a in \"$@\"; do\n" +
		"  if [ \"$prev\" = '--mcp-config' ]; then /bin/cat \"$a\" > \"$LAUNCH_FACTS/mcp-during\"; fi\n" +
		"  prev=\"$a\"\ndone\n" + output
	if err := os.WriteFile(filepath.Join(f.bin, agent), []byte(script), 0o700); err != nil {
		f.t.Fatal(err)
	}
}

func (f *launchFixture) serve(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/v1/episodes/ep/report":
		var start *store.EpisodeStart
		if !f.missingStart {
			start = f.read()
		}
		_ = json.NewEncoder(w).Encode(daemon.EpisodeReport{Start: start})
	case "/v1/session/self":
		if f.revoked {
			w.WriteHeader(401)
			return
		}
		left := int64(3600)
		if f.expired {
			left = 0
		}
		_ = json.NewEncoder(w).Encode(daemon.SessionSelf{Episode: f.selfEpisode, Worktree: filepath.Join(f.ws, "ep"), SecondsLeft: left,
			Addr: "http://" + defaultEpisodeAddr, Grants: []daemon.SessionGrant{{Effect: "Workspace.Read", Scope: "worktree", Budget: 50, Left: 50}}})
	default:
		f.t.Errorf("unexpected route %s", r.URL.Path)
		w.WriteHeader(404)
	}
}

func (f *launchFixture) read() *store.EpisodeStart {
	f.t.Helper()
	c, err := store.OpenCredentialsReadOnly(f.creds)
	if err != nil {
		f.t.Fatal(err)
	}
	defer c.Close()
	start, err := c.EpisodeStart(presetTestContext(f.t), "ep")
	if err != nil {
		f.t.Fatal(err)
	}
	return start
}

// rewrite is a fixture-only new mint; production never rewrites mint fields.
func (f *launchFixture) rewrite(change func(*store.EpisodeStart)) {
	f.t.Helper()
	start := f.read()
	change(start)
	// Replace the temporary test database with a fresh mint.
	if err := os.Remove(f.creds); err != nil {
		f.t.Fatal(err)
	}
	c, err := store.OpenCredentials(f.creds)
	if err != nil {
		f.t.Fatal(err)
	}
	defer c.Close()
	if err := c.MintSessionWithStart(presetTestContext(f.t), store.SessionRow{CredentialID: "fixture", EpisodeID: "ep", GrantsJSON: "[]", ExpiresAt: 9999999999, CreatedAt: 1}, start); err != nil {
		f.t.Fatal(err)
	}
}

func (f *launchFixture) run(extra ...string) (int, string, string) {
	f.t.Helper()
	var out, errOut bytes.Buffer
	args := append([]string{"ep", "--preset", "p", "--presets-dir", f.presets}, extra...)
	code := runEpisodeLaunch(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func (f *launchFixture) contents(path string) string {
	f.t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		f.t.Fatal(err)
	}
	return string(data)
}

func TestLaunchHeadlessClaudeAC9AndSecondLaunch(t *testing.T) {
	f := newLaunchFixture(t, "claude")
	// The stale mint config is a regular private file and must be replaced.
	mcp := filepath.Join(f.sessions, "ep.mcp.json")
	if err := os.WriteFile(mcp, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := f.run(); code != 0 {
		t.Fatalf("%d: %s %s", code, out, errOut)
	} else if !strings.Contains(out, "episode finish ep --preset p") || !strings.Contains(out, "episode report ep") {
		t.Fatal(out)
	}
	// The child reports its resolved cwd; on macOS t.TempDir() is under /var,
	// a symlink to /private/var, so compare against the resolved root.
	root, err := filepath.EvalSymlinks(f.root)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "launch", "ep")
	if got := strings.TrimSpace(f.contents(filepath.Join(f.facts, "cwd"))); got != dir {
		t.Fatalf("cwd %q != %q", got, dir)
	}
	if got := f.contents(filepath.Join(f.facts, "listing")); got != "" {
		t.Fatal("launch directory was not empty", got)
	}
	if env := f.contents(filepath.Join(f.facts, "env")); strings.Contains(env, f.token) || strings.Contains(env, "WORLD_SESSION=") {
		t.Fatal("claude inherited token")
	}
	if !strings.Contains(f.contents(filepath.Join(f.facts, "mcp-during")), f.token) {
		t.Fatal("no live mcp config")
	}
	for _, path := range []string{mcp, filepath.Join(f.sessions, "ep.launch.lock")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatal("not deleted", path, err)
		}
	}
	transcript := filepath.Join(f.sessions, "ep.transcript.jsonl")
	stream := "{\"type\":\"fixture\"}\n"
	if got := f.contents(transcript); got != stream {
		t.Fatalf("AC9 stream changed: %q", got)
	}
	for _, path := range []string{transcript, filepath.Join(f.sessions, "ep.agent.err")} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatal("private stream permissions", path, err)
		}
	}
	for _, path := range []string{dir, filepath.Dir(dir)} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o700 {
			t.Fatal("launch directory permissions", path, err)
		}
	}
	first := f.read()
	if len(first.Launches) != 1 {
		t.Fatal(first)
	}
	entry := first.Launches[0]
	raw, _ := json.Marshal(entry)
	if entry.Exit != 0 || entry.FinishReason != "exit" || entry.TranscriptSHA256 != serviceHash([]byte(stream)) ||
		entry.TranscriptState != "kept" || entry.MCPConfigDeleted == nil || !*entry.MCPConfigDeleted ||
		strings.Contains(string(raw), f.token) || strings.Contains(string(raw), "Make the fixture change") {
		t.Fatal(string(raw))
	}
	if !strings.Contains(strings.Join(entry.ArgvRedacted, " "), "<prompt sha256:") {
		t.Fatal(entry.ArgvRedacted)
	}
	t.Logf("AC9 argv: %q", entry.ArgvRedacted)
	t.Logf("AC9 launch record: %s", raw)
	if code, _, errOut := f.run(); code != 0 {
		t.Fatal(code, errOut)
	}
	second := f.read()
	if len(second.Launches) != 2 || !reflect.DeepEqual(first.Launches[0], second.Launches[0]) {
		t.Fatal("earlier launch changed", second)
	}
	if got := f.contents(transcript); got != stream+"{\"world_launch\": 2}\n"+stream {
		t.Fatalf("second stream %q", got)
	}
	for _, suffix := range []string{".session", ".task.md"} {
		if _, err := os.Stat(filepath.Join(f.sessions, "ep"+suffix)); err != nil {
			t.Fatal("session was removed", err)
		}
	}
	entries, _ := os.ReadDir(f.sessions)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".ep.") {
			t.Fatal("staging file retained", entry.Name())
		}
	}
}

func TestLaunchHeadlessCodexTokenAndPinnedMode(t *testing.T) {
	f := newLaunchFixture(t, "codex")
	// Mutable mode and model edits do not change the pinned launch.
	data, _ := os.ReadFile(filepath.Join(f.presets, "p.json"))
	var p launchPreset
	_ = json.Unmarshal(data, &p)
	p.Launch.Mode = "interactive"
	p.Model = "changed"
	data, _ = json.Marshal(p)
	_ = os.WriteFile(filepath.Join(f.presets, "p.json"), data, 0o600)
	if code, _, errOut := f.run(); code != 0 {
		t.Fatal(code, errOut)
	}
	env := f.contents(filepath.Join(f.facts, "env"))
	if !strings.Contains(env, "WORLD_SESSION="+f.token+"\n") {
		t.Fatal("codex token missing")
	}
	argv := strings.Split(strings.TrimSuffix(f.contents(filepath.Join(f.facts, "argv")), "\x00"), "\x00")
	if strings.Contains(strings.Join(argv, " "), f.token) || !strings.HasPrefix(argv[len(argv)-2], "developer_instructions=") || argv[len(argv)-3] != "-c" {
		t.Fatal("codex argv", argv)
	}
	if !strings.Contains(strings.Join(argv, " "), "fixture-model") || strings.Contains(strings.Join(argv, " "), "changed") {
		t.Fatal("model was not pinned", argv)
	}
	entries, _ := os.ReadDir(f.sessions)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".mcp.json") {
			t.Fatal("codex wrote mcp config")
		}
		if strings.Contains(f.contents(filepath.Join(f.sessions, e.Name())), f.token) && e.Name() != "ep.session" {
			t.Fatal("codex leaked token to", e.Name())
		}
	}
}

func TestLaunchStreamTokenScrubAndNonzero(t *testing.T) {
	for _, agent := range []string{"claude", "codex"} {
		t.Run(agent, func(t *testing.T) {
			f := newLaunchFixture(t, agent)
			output := "printf '%s\\n' \"$WORLD_SESSION\"\nprintf '%s\\n' \"$WORLD_SESSION\" >&2\nexit 7\n"
			if agent == "claude" {
				output = "/bin/cat \"$LAUNCH_FACTS/mcp-during\"\n/bin/cat \"$LAUNCH_FACTS/mcp-during\" >&2\nexit 7\n"
			}
			f.fake(agent, output)
			if code, _, errOut := f.run(); code != 1 || !strings.Contains(errOut, "exit 7 (exit)") {
				t.Fatal(code, errOut)
			}
			for _, suffix := range []string{".transcript.jsonl", ".agent.err"} {
				kept := f.contents(filepath.Join(f.sessions, "ep"+suffix))
				if strings.Contains(kept, f.token) || strings.Count(kept, "\n") != 1 || !strings.Contains(kept, "scrubbed") {
					t.Fatal("unsafe stream", kept)
				}
			}
			entry := f.read().Launches[0]
			if entry.TranscriptState != "scrubbed" || entry.Exit != 7 || entry.TranscriptSHA256 != serviceHash([]byte(f.contents(entry.Transcript))) {
				t.Fatal(entry)
			}
		})
	}
}

func TestLaunchRefusalsStartNothing(t *testing.T) {
	cases := []struct {
		name, want string
		change     func(*launchFixture)
		extra      []string
	}{
		{"nonempty", "is not empty", func(f *launchFixture) {
			dir := filepath.Join(f.root, "launch", "ep")
			_ = os.MkdirAll(dir, 0o700)
			_ = os.WriteFile(filepath.Join(dir, "untouched"), []byte("keep"), 0o600)
		}, nil},
		{"symlink", "symlink", func(f *launchFixture) {
			_ = os.MkdirAll(filepath.Join(f.root, "launch"), 0o700)
			_ = os.Symlink(f.facts, filepath.Join(f.root, "launch", "ep"))
		}, nil},
		{"workspace", "overlaps", func(f *launchFixture) {
			f.ws = filepath.Join(f.root, "launch")
			var p launchPreset
			_ = json.Unmarshal([]byte(f.start.Preset.Bytes), &p)
			p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: f.db, CredentialsDB: f.creds, WorkspaceRoot: f.ws, SessionsDir: f.sessions})
			b, _ := json.Marshal(p)
			_ = os.WriteFile(filepath.Join(f.presets, "p.json"), b, 0o600)
		}, nil},
		{"lock", "pid", func(f *launchFixture) {
			_ = os.WriteFile(filepath.Join(f.sessions, "ep.launch.lock"), []byte(fmt.Sprint(os.Getpid())), 0o600)
		}, nil},
		{"expired", "--branch world/ep", func(f *launchFixture) { f.expired = true }, nil},
		{"revoked", "--branch world/ep", func(f *launchFixture) { f.revoked = true }, nil},
		{"wrong-session", "--branch world/ep", func(f *launchFixture) { f.selfEpisode = "other" }, nil},
		{"no-start", "has no start record", func(f *launchFixture) { f.missingStart = true }, nil},
		{"preset-mismatch", "does not match", func(f *launchFixture) { f.rewrite(func(s *store.EpisodeStart) { s.Preset.Name = "other" }) }, nil},
		{"no-task", "has no task", func(f *launchFixture) { f.rewrite(func(s *store.EpisodeStart) { s.Task = nil }) }, nil},
		{"missing-task-copy", "is missing", func(f *launchFixture) { _ = os.Remove(filepath.Join(f.sessions, "ep.task.md")) }, nil},
		{"task-digest", "digest does not match", func(f *launchFixture) {
			_ = os.WriteFile(filepath.Join(f.sessions, "ep.task.md"), []byte("tampered"), 0o600)
		}, nil},
		{"interactive", "interactive launch is not built yet", func(f *launchFixture) {
			f.rewrite(func(s *store.EpisodeStart) { s.Preset.Bytes = `{"launch":{"mode":"interactive"}}` })
		}, nil},
		{"default-interactive", "interactive launch is not built yet", func(f *launchFixture) {
			f.rewrite(func(s *store.EpisodeStart) { s.Preset.Bytes = "{}" })
		}, nil},
		{"binary", "missing from PATH", func(f *launchFixture) { f.t.Setenv("PATH", f.t.TempDir()) }, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newLaunchFixture(t, "claude")
			tc.change(f)
			before := launchFileSnapshot(t, f.root)
			if code, _, errOut := f.run(tc.extra...); code == 0 || !strings.Contains(errOut, tc.want) || !strings.Contains(errOut, "why:") || !strings.Contains(errOut, "next:") || strings.Contains(errOut, "--resume") {
				t.Fatal(code, errOut)
			}
			if after := launchFileSnapshot(t, f.root); !reflect.DeepEqual(before, after) {
				t.Fatal("refusal changed files or directories", before, after)
			}
			if _, err := os.Stat(filepath.Join(f.facts, "started")); !os.IsNotExist(err) {
				t.Fatal("agent started", err)
			}
			if len(f.read().Launches) != 0 {
				t.Fatal("refusal recorded a launch")
			}
			if tc.name == "nonempty" && f.contents(filepath.Join(f.root, "launch", "ep", "untouched")) != "keep" {
				t.Fatal("contents changed")
			}
		})
	}
}

func TestLaunchHeadlessOverrideAndTimeout(t *testing.T) {
	f := newLaunchFixture(t, "codex")
	f.rewrite(func(s *store.EpisodeStart) { s.Preset.Bytes = `{"launch":{"mode":"interactive","deadline_s":1}}` })
	f.fake("codex", "/bin/sleep 30 &\nchild=$!\nprintf '%s\\n' \"$child\" > \"$LAUNCH_FACTS/child\"\nwait\n")
	before := time.Now()
	if code, _, errOut := f.run("--headless"); code != 1 || !strings.Contains(errOut, "(timeout)") {
		t.Fatal(code, errOut)
	}
	if time.Since(before) > 4*time.Second {
		t.Fatal("wall cap not enforced")
	}
	entry := f.read().Launches[0]
	if entry.FinishReason != "timeout" {
		t.Fatal(entry)
	}
	// The package-level wall cap test also checks the child cannot keep writing.
}

func TestLaunchAndStartHelp(t *testing.T) {
	for _, verb := range []string{"start", "launch"} {
		var out, errOut bytes.Buffer
		if code := runEpisode("", []string{verb, "--help"}, &out, &errOut); code != 0 {
			t.Fatal(code, errOut.String())
		}
		for _, flag := range []string{"--preset", "--presets-dir", "--headless"} {
			if !strings.Contains(out.String(), flag) {
				t.Fatal("missing flag", flag, out.String())
			}
		}
		if verb == "start" && !strings.Contains(out.String(), "--launch") {
			t.Fatal(out.String())
		}
	}
}

func launchFileSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		value := info.Mode().String()
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += ":" + serviceHash(data)
		} else if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			value += ":" + target
		}
		out[rel] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestLaunchTokenScanAcrossChunksAndPartialLines(t *testing.T) {
	f := newLaunchFixture(t, "codex")
	stage, err := os.CreateTemp(f.sessions, ".scan-*")
	if err != nil {
		t.Fatal(err)
	}
	defer stage.Close()
	// The literal token crosses the scanner's 32 KiB read boundary.
	data := strings.Repeat("x", 32*1024-10) + f.token + "\n"
	_, _ = stage.WriteString(data)
	path := filepath.Join(f.sessions, "scan.jsonl")
	state, digest, err := keepLaunchStream(stage, path, f.token, 1)
	if err != nil || state != "scrubbed" || strings.Contains(f.contents(path), f.token) || digest != serviceHash([]byte(f.contents(path))) {
		t.Fatal(state, digest, err)
	}
	_ = stage.Truncate(0)
	_, _ = stage.Seek(0, 0)
	_, _ = stage.WriteString("partial")
	state, _, err = keepLaunchStream(stage, path, f.token, 1)
	if err != nil || state != "kept" || f.contents(path) != "partial" {
		t.Fatal(state, err)
	}
	state, _, err = keepLaunchStream(stage, path, f.token, 2)
	if err != nil || f.contents(path) != "partial\n{\"world_launch\": 2}\npartial" {
		t.Fatal(state, err, f.contents(path))
	}
}

func TestLaunchProfileCommandsAndRedaction(t *testing.T) {
	f := newPresetFixture(t)
	health := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
	commands, err := launchCommands(health, "proj")
	if err != nil || len(commands) == 0 {
		t.Fatal(commands, err)
	}
	for i := 1; i < len(commands); i++ {
		if commands[i-1] >= commands[i] {
			t.Fatal("commands not in byte order", commands)
		}
	}
	if _, err := launchCommands(health, "absent"); err == nil {
		t.Fatal("missing profile accepted")
	}
	health.ExecProfiles[0].Digest = "sha256:wrong"
	if _, err := launchCommands(health, "proj"); err == nil {
		t.Fatal("digest mismatch accepted")
	}
	got := redactLaunchArgv([]string{"secret", "Bearer value", "P", "context"}, "P", "secret")
	want := []string{"<redacted>", "<redacted>", "<prompt " + serviceHash([]byte("P")) + ">", "context"}
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
}
