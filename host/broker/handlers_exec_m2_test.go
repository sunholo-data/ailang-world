package broker

// Row 140 M2 (w-workspace-exec-toolchain-effect §4.5, §6 AC2.4, AC2.6, AC2.7,
// AC2.9, §7): ExecHandler against a PASS-THROUGH srt stub. The stub stands
// where node runs srt's cli.js: it records each launch (argv, host env),
// refuses a launch whose srt options are not closed by "--" (srt would parse
// them, V12: MUT-NO-DASHDASH) and then runs the trampoline UNSANDBOXED. It
// proves the launch line, the grammar gate, the capture, the deadline and the
// per-call re-verification without srt; the confinement itself is proved
// against the real srt in handlers_exec_srt_test.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

type execFixture struct {
	h                                *ExecHandler
	base, root, worktree, sibling    string
	state, cache, home               string
	calls, hostEnv, argvLog, stubLog string
	profile                          *ExecProfile
	sandbox                          *ExecSandbox
}

type execFixtureOpts struct {
	node    string                 // overrides the stub as the host node
	tree    string                 // a node_modules dir; "" = a fake tree
	pin     *ExecSandboxPin        // nil = the fake tree's own pin
	profile func(m map[string]any) // edits the default profile
	// profileF edits it knowing the fixture's paths (set before it runs).
	profileF func(m map[string]any, f *execFixture)
	setup    func(f *execFixture) // runs before the handler is built
	home     string               // the operator HOME; "" = a scratch dir
}

// stubNode writes the pass-through srt stub.
func stubNode(t *testing.T, dir string) (node, calls, hostEnv, argvLog string) {
	t.Helper()
	calls, hostEnv, argvLog = filepath.Join(dir, "stub.calls"), filepath.Join(dir, "stub.hostenv"), filepath.Join(dir, "stub.argv")
	node = filepath.Join(dir, "stub-node")
	body := "#!/bin/sh\n" +
		"printf 'call\\n' >> " + shellQuote(calls) + "\n" +
		"/usr/bin/env > " + shellQuote(hostEnv) + "\n" +
		": > " + shellQuote(argvLog) + "\n" +
		"for a in \"$@\"; do printf '%s\\0' \"$a\" >> " + shellQuote(argvLog) + "; done\n" +
		"if [ \"$2\" != \"--settings\" ] || [ \"$4\" != \"--\" ]; then\n" +
		"  echo \"STUB-SRT: srt would parse '$2' '$4' as its own options\" >&2; exit 99\nfi\n" +
		"shift 4\nexec \"$@\"\n"
	if err := os.WriteFile(node, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return node, calls, hostEnv, argvLog
}

func defaultExecProfile(f *execFixture) map[string]any {
	return map[string]any{
		"profile": ExecProfileVersion, "project": "proj", "root": ".",
		"path": []any{"/usr/bin", "/bin"}, "timeout_ms": 9000,
		"commands": map[string]any{
			"test": map[string]any{"argv": []any{"printf", "[%s]"}, "flags": map[string]any{
				"-run": "regex", "-v": "bool", "-k": map[string]any{"class": "regex", "form": "sep"}},
				"positional": "relpath", "max_args": 8},
			"exit":  map[string]any{"argv": []any{"/bin/sh", "-c", "exit 3"}},
			"term":  map[string]any{"argv": []any{"/bin/sh", "-c", "kill -TERM $$"}},
			"big":   map[string]any{"argv": []any{"/bin/sh", "-c", "yes 0123456789abcdef | head -c 2000000; printf 'tail-err' >&2"}},
			"ctl":   map[string]any{"argv": []any{"/bin/sh", "-c", "head -c 300000 /dev/zero; head -c 300000 /dev/zero >&2"}},
			"env":   map[string]any{"argv": []any{"/usr/bin/env"}},
			"sleep": map[string]any{"argv": []any{"/bin/sh", "-c", "sleep 3011 & sleep 3011 & echo started; wait"}},
		},
	}
}

func newExecFixture(t *testing.T, o execFixtureOpts) *execFixture {
	t.Helper()
	base := canonicalTempDir(t)
	f := &execFixture{base: base, root: filepath.Join(base, "workspace"), state: filepath.Join(base, "state")}
	f.worktree, f.sibling = filepath.Join(f.root, "ep1"), filepath.Join(f.root, "ep2")
	f.home = o.home
	if f.home == "" {
		f.home = filepath.Join(base, "home")
	}
	for _, d := range []string{f.worktree, f.sibling, f.state, f.home} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	stub, calls, hostEnv, argvLog := stubNode(t, base)
	f.calls, f.hostEnv, f.argvLog = calls, hostEnv, argvLog
	node := stub
	if o.node != "" {
		node = o.node
	}
	tree := o.tree
	if tree == "" {
		tree = fakeSrtTree(t, "0.0.0-test", "// stub cli\n")
	}
	pin := fakePin(t, tree)
	if o.pin != nil {
		pin = *o.pin
	}
	sb, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: tree, Node: node, StateDir: f.state, Pin: pin})
	if err != nil {
		t.Fatalf("ArchiveExecSandbox: %v", err)
	}
	f.sandbox = sb
	m := defaultExecProfile(f)
	if o.profile != nil {
		o.profile(m)
	}
	if o.profileF != nil {
		o.profileF(m, f)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if f.profile, err = ParseExecProfile(data); err != nil {
		t.Fatalf("ParseExecProfile: %v", err)
	}
	f.cache = filepath.Join(f.state, "exec-cache", "ep1", f.profile.Project)
	if o.setup != nil {
		o.setup(f)
	}
	f.h, err = NewExecHandler(ExecHandlerConfig{Profile: f.profile, Sandbox: sb, Episode: "ep1", Worktree: f.worktree,
		WorkspaceRoot: f.root, StateDir: f.state, OperatorHome: f.home})
	if err != nil {
		t.Fatalf("NewExecHandler: %v", err)
	}
	return f
}

func (f *execFixture) call(t *testing.T, command string, args ...string) ([]byte, error) {
	t.Helper()
	payload := map[string]any{"command": command}
	if args != nil {
		payload["args"] = args
	}
	data, _ := json.Marshal(payload)
	return f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1}, data)
}

// result decodes a run's output strictly into the §4.5 shape.
func execDecode(t *testing.T, out []byte) execResult {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(out))
	dec.DisallowUnknownFields()
	var r execResult
	if err := dec.Decode(&r); err != nil {
		t.Fatalf("result %s is not the §4.5 shape: %v", out, err)
	}
	return r
}

func (f *execFixture) mustRun(t *testing.T, command string, args ...string) execResult {
	t.Helper()
	out, err := f.call(t, command, args...)
	if err != nil {
		t.Fatalf("Execute(%s %q): %v", command, args, err)
	}
	return execDecode(t, out)
}

// stubCalls is how many times the host node (the stub) was launched.
func (f *execFixture) stubCalls(t *testing.T) int {
	data, err := os.ReadFile(f.calls)
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(data), "call\n")
}

func TestExecHandlerResultShape(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	out, err := f.call(t, "test", "-run", "TestX", "a")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	want := []string{"argv", "duration_ms", "exit_code", "limit", "profile", "sandbox", "stderr", "stderr_bytes",
		"stderr_sha256", "stderr_truncated", "stdout", "stdout_bytes", "stdout_sha256", "stdout_truncated", "timed_out"}
	if got := keysOf(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("result keys = %v, want %v", got, want)
	}
	if got := keysOf(m["profile"].(map[string]any)); !reflect.DeepEqual(got, []string{"argv0_sha256", "command", "digest", "project"}) {
		t.Fatalf("profile keys = %v", got)
	}
	if got := keysOf(m["sandbox"].(map[string]any)); !reflect.DeepEqual(got,
		[]string{"cli_sha256", "network", "read_fence", "runtime", "settings_digest", "version"}) {
		t.Fatalf("sandbox keys = %v", got)
	}
	r := execDecode(t, out)
	if r.ExitCode == nil || *r.ExitCode != 0 || r.TimedOut || r.Limit != nil || r.Stdout != "[-run=TestX][a]" {
		t.Fatalf("result = %+v", r)
	}
	c := f.profile.Commands["test"]
	if r.Profile != (execResultProfile{Project: "proj", Digest: f.profile.Digest, Command: "test", Argv0SHA256: c.Argv0SHA256}) {
		t.Fatalf("profile block = %+v", r.Profile)
	}
	settings, _ := os.ReadFile(filepath.Join(f.state, "exec", "ep1.proj.srt.json"))
	if r.Sandbox != (execResultSandbox{Runtime: ExecSandboxPackage, Version: "0.0.0-test", CLISHA256: f.sandbox.CLISHA256,
		SettingsDigest: "sha256:" + sha256Hex(settings), ReadFence: true, Network: "none"}) {
		t.Fatalf("sandbox block = %+v", r.Sandbox)
	}
	if !reflect.DeepEqual(r.Argv, []string{"printf", "[%s]", "-run=TestX", "a"}) {
		t.Fatalf("argv = %q, want the profile prefix plus the emitted args", r.Argv)
	}
}

// The exact §4.5 launch line: srt's options closed by "--", the trampoline,
// World's K=V set, the absolute argv0, the prefix and the emitted args — and
// the host node's environment exactly PATH=/usr/bin:/bin (MUT-NO-DASHDASH,
// MUT-NO-TRAMPOLINE, MUT-HOST-ENV, MUT-SHELL-JOIN).
func TestExecHandlerLaunchLineAndHostEnv(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
		m["env"] = map[string]any{"NODE_OPTIONS": "--require=./zz.cjs", "GOFLAGS": "-mod=readonly"}
		m["caches"] = map[string]any{"GOCACHE": map[string]any{}}
	}})
	r := f.mustRun(t, "test", "-k", "a b", "-v", "$(id)", "x;y")
	if r.ExitCode == nil || *r.ExitCode != 0 || r.Stdout != "[-k][a b][-v][$(id)][x;y]" {
		t.Fatalf("stdout = %q exit %v: the args reached the command altered (MUT-SHELL-JOIN) or srt took them", r.Stdout, r.ExitCode)
	}
	raw, err := os.ReadFile(f.argvLog)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(raw), "\x00"), "\x00")
	kv := []string{"HOME=" + f.cache + "/home", "TMPDIR=" + f.cache + "/tmp", "PATH=/usr/bin:/bin", "LANG=C.UTF-8",
		"LC_ALL=C.UTF-8", "PYTHONPYCACHEPREFIX=" + f.cache + "/pycache", "npm_config_cache=" + f.cache + "/npm", "CI=1",
		"NO_COLOR=1", "GOFLAGS=-mod=readonly", "NODE_OPTIONS=--require=./zz.cjs", "GOCACHE=" + f.cache + "/GOCACHE"}
	want := append([]string{f.sandbox.CLI(), "--settings", filepath.Join(f.state, "exec", "ep1.proj.srt.json"), "--",
		"/bin/sh", "-c", `/usr/bin/env -i "$@"; exit $?`, "world-exec"}, kv...)
	want = append(want, "/usr/bin/printf", "[%s]", "-k", "a b", "-v", "$(id)", "x;y")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("launch line =\n%q\nwant\n%q", got, want)
	}
	hostEnv := readEnvDump(t, f.hostEnv)
	var kept []string
	for _, kv := range hostEnv {
		// /bin/sh exports these itself when it runs the stub.
		if k, _, _ := strings.Cut(kv, "="); k != "PWD" && k != "SHLVL" && k != "_" && k != "OLDPWD" {
			kept = append(kept, kv)
		}
	}
	if !reflect.DeepEqual(kept, []string{"PATH=/usr/bin:/bin"}) {
		t.Fatalf("the host node's env = %q, want exactly PATH=/usr/bin:/bin (V40, MUT-HOST-ENV)", kept)
	}
	// The child's env is exactly the K=V set (env -i inside the trampoline).
	r = f.mustRun(t, "env")
	child := strings.Split(strings.TrimSpace(r.Stdout), "\n")
	sort.Strings(child)
	wantKV := append([]string(nil), kv...)
	sort.Strings(wantKV)
	if !reflect.DeepEqual(child, wantKV) {
		t.Fatalf("child env = %q, want exactly %q", child, wantKV)
	}
}

// AC2.4: every grammar refusal is answered before anything is spawned —
// neither the cache seeding nor the host node runs.
func TestExecHandlerRefusalsSpawnNothing(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	for _, tc := range []struct {
		name, command string
		args          []string
		want          string
	}{
		{"unprofiled command", "deploy", nil, `"deploy" is not a command`},
		{"-exec=sh", "test", []string{"-exec=sh"}, "does not list this flag"},
		{"--", "test", []string{"--", "a"}, `"--" is admitted only`},
		{"-c", "test", []string{"-c"}, "does not list this flag"},
		{"-s=x.json", "test", []string{"-s=x.json"}, "does not list this flag"},
		{"../x", "test", []string{"../x"}, "not a positional"},
		{"/abs", "test", []string{"/abs"}, "not a positional"},
		{"newline in a flag value", "test", []string{"-run=a\nb"}, "does not match the flag's class"},
		{"--settings", "test", []string{"--settings"}, "does not list this flag"},
		{"bad command id", "Test", nil, "command"},
		{"unknown payload key", "test", nil, `"cwd"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]any{"command": tc.command}
			if tc.args != nil {
				payload["args"] = tc.args
			}
			if tc.name == "unknown payload key" {
				payload["cwd"] = "/"
			}
			data, _ := json.Marshal(payload)
			out, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1}, data)
			if err != nil {
				t.Fatalf("Execute: %v (a refusal is an answer, not a failure)", err)
			}
			var got struct {
				OK      bool   `json:"ok"`
				Refused string `json:"refused"`
			}
			if json.Unmarshal(out, &got) != nil || got.OK || !strings.Contains(got.Refused, tc.want) {
				t.Fatalf("output = %s, want ok:false refused naming %q", out, tc.want)
			}
			if n := f.h.spawns.Load(); n != 0 {
				t.Fatalf("exec spawn counter = %d after a refusal, want 0", n)
			}
			if n := f.stubCalls(t); n != 0 {
				t.Fatalf("the host node ran %d times after a refusal", n)
			}
		})
	}
}

func TestExecHandlerExitStatusThroughTheTrampoline(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	if r := f.mustRun(t, "exit"); r.ExitCode == nil || *r.ExitCode != 3 {
		t.Fatalf("exit 3 -> %v", r.ExitCode)
	}
	if r := f.mustRun(t, "term"); r.ExitCode == nil || *r.ExitCode != 143 {
		t.Fatalf("SIGTERM self-kill -> %v, want 143", r.ExitCode)
	}
}

// AC2.6: 2 MB of output is a truncated SUCCESS with the full stream's total
// and digest (MUT-OVERFLOW-ERROR), and the worst-case stream — every byte a
// six-byte JSON escape, on both streams — still fits the coordinator's 1 MiB
// MaxOutput with room for the envelope.
func TestExecHandlerTruncatesLargeOutput(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	out, err := f.call(t, "big")
	if err != nil {
		t.Fatalf("a 2 MB run failed: %v (the exec handler must not use the strict capture)", err)
	}
	r := execDecode(t, out)
	if !r.StdoutTruncated || r.StdoutBytes != 2000000 || r.StdoutSHA256 != sha256Hex(repeatStream(2000000)) || r.Limit != nil {
		t.Fatalf("2 MB stdout: truncated=%t bytes=%d sha=%s limit=%v", r.StdoutTruncated, r.StdoutBytes, r.StdoutSHA256, r.Limit)
	}
	if r.Stderr != "tail-err" || r.StderrTruncated {
		t.Fatalf("stderr = %q", r.Stderr)
	}
	if len(out) >= 1<<20 {
		t.Fatalf("result is %d bytes, want < 1 MiB", len(out))
	}
	out, err = f.call(t, "ctl")
	if err != nil {
		t.Fatal(err)
	}
	r = execDecode(t, out)
	if !r.StdoutTruncated || !r.StderrTruncated || r.StdoutBytes != 300000 {
		t.Fatalf("control-char run: %+v", r)
	}
	// coordinator MaxOutput (daemon.go: MaxOutput 1 << 20) less a 128 KiB
	// allowance for the world block the coordinator adds around it.
	if limit := 1<<20 - 128<<10; len(out) > limit {
		t.Fatalf("worst-case result is %d bytes, over %d", len(out), limit)
	}
	t.Logf("worst-case all-control-character result: %d bytes", len(out))
}

// A stream past the operator's hard cap is killed and recorded limit:"output".
func TestExecHandlerHardCapRecordsLimitOutput(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	f.h.maxOutput = 1 << 20
	r := f.mustRun(t, "big")
	if r.Limit == nil || *r.Limit != "output" || r.ExitCode != nil || r.TimedOut {
		t.Fatalf("over the cap: limit=%v exit=%v timed_out=%t", r.Limit, r.ExitCode, r.TimedOut)
	}
}

// AC2.7 (stub half; the real-srt half is TestExecSrtTimeoutKillsTheGroup).
func TestExecHandlerTimeoutKillsTheGroup(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) { m["timeout_ms"] = 800 }})
	start := time.Now()
	r := f.mustRun(t, "sleep")
	elapsed := time.Since(start)
	if !r.TimedOut || r.ExitCode != nil || !strings.Contains(r.Stdout, "started") {
		t.Fatalf("result = %+v, want timed_out with the partial head", r)
	}
	if elapsed > 800*time.Millisecond+3*time.Second {
		t.Fatalf("returned after %s; the deadline is 800 ms", elapsed)
	}
	if n := countProcs(t, "sleep 3011"); n != 0 {
		t.Fatalf("%d group members alive after the deadline (MUT-NO-GROUP-KILL)", n)
	}
}

// countProcs counts processes whose full command line is exactly cmdline,
// polling up to a second for the kill to land.
func countProcs(t *testing.T, cmdline string) int {
	t.Helper()
	n := 0
	for i := 0; i < 10; i++ {
		out, _ := exec.Command("pgrep", "-f", "^"+cmdline+"$").Output()
		if n = len(strings.Fields(string(out))); n == 0 {
			return 0
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = exec.Command("pkill", "-9", "-f", "^"+cmdline+"$").Run()
	return n
}

// The handler's deadline is min(timeout_ms, the call's remaining budget).
func TestExecHandlerDeadlineIsTheSmallerBudget(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()
	data, _ := json.Marshal(map[string]any{"command": "sleep"})
	start := time.Now()
	out, err := f.h.Execute(ctx, EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope, Cost: 1}, data)
	if err != nil {
		t.Fatal(err)
	}
	if r := execDecode(t, out); !r.TimedOut {
		t.Fatalf("result = %+v, want timed_out on the 700 ms call budget", r)
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("returned after %s; the call budget was 700 ms, timeout_ms 9000", elapsed)
	}
	countProcs(t, "sleep 3011")
}

// AC2.9: after startup a flipped byte of the archived cli.js, or a swapped
// argv0, is a recorded failure naming the file, with nothing spawned
// (MUT-STACK-UNVERIFIED).
func TestExecHandlerReverifiesTheStackBeforeEachCall(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, f *execFixture) string
	}{
		{"archived cli.js byte flipped", func(t *testing.T, f *execFixture) string {
			flipByte(t, f.sandbox.CLI())
			return f.sandbox.CLI()
		}},
		{"argv0 swapped", func(t *testing.T, f *execFixture) string {
			p := f.profile.Commands["test"].Argv0
			flipByte(t, p)
			return p
		}},
		{"node swapped", func(t *testing.T, f *execFixture) string {
			flipByte(t, f.sandbox.Node)
			return f.sandbox.Node
		}},
		{"settings rewritten", func(t *testing.T, f *execFixture) string {
			p := filepath.Join(f.state, "exec", "ep1.proj.srt.json")
			if err := os.WriteFile(p, []byte(`{"filesystem":{}}`), 0o600); err != nil {
				t.Fatal(err)
			}
			return p
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bin := t.TempDir()
			printf := filepath.Join(bin, "printf")
			// A script, not a copied Mach-O: macOS scans a fresh binary's
			// signature on first exec, which can stall for seconds.
			if err := os.WriteFile(printf, []byte("#!/bin/sh\nexec /usr/bin/printf \"$@\"\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
				m["path"] = []any{bin, "/usr/bin", "/bin"}
			}})
			f.mustRun(t, "test", "a") // verifies clean, then the stack drifts
			before, calls := f.h.spawns.Load(), f.stubCalls(t)
			want := tc.mutate(t, f)
			out, err := f.call(t, "test", "a")
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("Execute = %s, %v; want a failure naming %s", out, err, want)
			}
			if f.h.spawns.Load() != before || f.stubCalls(t) != calls {
				t.Fatal("something was spawned after the stack drifted")
			}
		})
	}
}

func TestExecHandlerProjectRootAndCaches(t *testing.T) {
	seed := filepath.Join(canonicalTempDir(t), "seed")
	if err := os.MkdirAll(filepath.Join(seed, "ab"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(seed, "ab", "entry"), []byte("seeded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f := newExecFixture(t, execFixtureOpts{
		profile: func(m map[string]any) {
			m["root"] = "sub"
			m["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": seed}, "EMPTYCACHE": map[string]any{}}
			m["commands"].(map[string]any)["pwd"] = map[string]any{"argv": []any{"/bin/sh", "-c",
				`pwd -P; cat "$GOCACHE/ab/entry"; echo new > "$GOCACHE/ab/entry"; ls -d "$EMPTYCACHE" "$HOME" "$TMPDIR"`}}
		},
		setup: func(f *execFixture) {
			if err := os.MkdirAll(filepath.Join(f.worktree, "sub"), 0o755); err != nil {
				t.Fatal(err)
			}
		},
	})
	r := f.mustRun(t, "pwd")
	lines := strings.Split(strings.TrimSpace(r.Stdout), "\n")
	if len(lines) != 5 || lines[0] != filepath.Join(f.worktree, "sub") || lines[1] != "seeded" {
		t.Fatalf("result %+v, want cwd <worktree>/sub and the seeded entry", r)
	}
	if got, _ := os.ReadFile(filepath.Join(seed, "ab", "entry")); string(got) != "seeded\n" {
		t.Fatalf("the seed changed to %q; the cache must be a clone", got)
	}
	// The episode cache persists across calls and is not re-seeded.
	r = f.mustRun(t, "pwd")
	if lines := strings.Split(r.Stdout, "\n"); len(lines) < 2 || lines[1] != "new" {
		t.Fatalf("second call stdout = %q; the episode cache was re-seeded", r.Stdout)
	}
}

func TestExecHandlerRefusesOtherEffectsAndScopes(t *testing.T) {
	f := newExecFixture(t, execFixtureOpts{})
	for _, req := range []EffectRequest{
		{Effect: EffectWorkspaceExec, Scope: "/", Cost: 1},
		{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
	} {
		if out, err := f.h.Execute(boundedTestContext(t), req, []byte(`{"command":"test"}`)); err == nil {
			t.Fatalf("Execute(%+v) = %s, want an error", req, out)
		}
	}
}

// driveExecHostNode is the AC10(a) credential arm's first Workspace.Exec
// environment (MUT-ENV-INHERIT): the probe stands where node runs srt, so it
// records the HOST-side environment.
func driveExecHostNode(t *testing.T, probe string) {
	t.Helper()
	f := newExecFixture(t, execFixtureOpts{node: probe})
	_, _ = f.call(t, "test", "a")
}

// driveExecChild is the second: the probe is the profiled command, run by the
// pass-through stub through the trampoline, so it records the CHILD's
// environment.
func driveExecChild(t *testing.T, probe string) {
	t.Helper()
	f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
		m["path"] = []any{filepath.Dir(probe), "/usr/bin", "/bin"}
		m["commands"].(map[string]any)["probe"] = map[string]any{"argv": []any{filepath.Base(probe)}}
	}})
	if _, err := f.call(t, "probe"); err != nil {
		t.Fatalf("drive the exec child: %v", err)
	}
}
