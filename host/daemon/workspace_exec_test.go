package daemon

// Row 140 M2 (w-workspace-exec-toolchain-effect §4.5, §10): the real exec
// handler is bound only behind a configured profile. With none, the M1 typed
// refusal stays (TestWorkspaceExecUnconfiguredSpawnsNothing, unedited). The
// `serve --exec-*` flags that fill workspaceTools.exec are M3; this file sets
// the field directly, as they will.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

// stubExecSandbox archives a minimal srt-shaped tree run by a pass-through
// stub node (`shift 4; exec "$@"`): the trampoline runs unsandboxed, which
// is all the wiring needs to show.
func stubExecSandbox(t *testing.T, stateDir string) *broker.ExecSandbox {
	t.Helper()
	dir := canonicalTemp(t)
	nm := filepath.Join(dir, "node_modules")
	pkg := filepath.Join(nm, "@anthropic-ai", "sandbox-runtime")
	if err := os.MkdirAll(filepath.Join(pkg, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, "package.json"), `{"name":"@anthropic-ai/sandbox-runtime","version":"0.0.0-test"}`)
	cli := "// stub cli\n"
	writeFile(t, filepath.Join(pkg, "dist", "cli.js"), cli)
	node := filepath.Join(dir, "stub-node")
	if err := os.WriteFile(node, []byte("#!/bin/sh\nshift 4\nexec \"$@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(cli))
	sb, err := broker.ArchiveExecSandbox(broker.ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: stateDir,
		Pin: broker.ExecSandboxPin{Version: "0.0.0-test", CLISHA256: hex.EncodeToString(sum[:])}})
	if err != nil {
		t.Fatal(err)
	}
	return sb
}

func execProfile(t *testing.T, readRoots ...string) *broker.ExecProfile {
	t.Helper()
	m := map[string]any{"profile": broker.ExecProfileVersion, "project": "proj", "path": []any{"/usr/bin", "/bin"},
		"timeout_ms": 5000, "read_roots": readRoots,
		"commands": map[string]any{"where": map[string]any{"argv": []any{"/bin/sh", "-c", "pwd -P"}}}}
	if readRoots == nil {
		delete(m, "read_roots")
	}
	data, _ := json.Marshal(m)
	p, err := broker.ParseExecProfile(data)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func execCall(t *testing.T, h broker.Handler) ([]byte, error) {
	t.Helper()
	return h.Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectWorkspaceExec, Scope: broker.WorkspaceScope, Cost: 1},
		[]byte(`{"command":"where"}`))
}

func TestWorkspaceExecConfiguredBindsTheEpisodeHandler(t *testing.T) {
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	d.workspace.exec = &workspaceExec{profile: execProfile(t), sandbox: stubExecSandbox(t, f.stateDir)}

	h1 := d.workspace.registry("ep1")[broker.EffectWorkspaceExec]
	if _, ok := h1.(*broker.ExecHandler); !ok {
		t.Fatalf("configured registry(ep1) binds %T, want *broker.ExecHandler", h1)
	}
	out, err := execCall(t, h1)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var res struct {
		ExitCode *int   `json:"exit_code"`
		Stdout   string `json:"stdout"`
	}
	if err := json.Unmarshal(out, &res); err != nil || res.ExitCode == nil || *res.ExitCode != 0 ||
		strings.TrimSpace(res.Stdout) != filepath.Join(f.root, "ep1") {
		t.Fatalf("Workspace.Exec in ep1 = %s (%v), want exit 0 run in the ep1 worktree", out, err)
	}
	if _, err := os.Stat(filepath.Join(f.stateDir, "exec", "ep1.proj.srt.json")); err != nil {
		t.Fatalf("the episode's settings were not rendered under the state dir: %v", err)
	}
	// Per episode: ep2 gets its own worktree; ep1's handler is cached.
	out, err = execCall(t, d.workspace.registry("ep2")[broker.EffectWorkspaceExec])
	if err != nil || !strings.Contains(string(out), filepath.Join(f.root, "ep2")) {
		t.Fatalf("Workspace.Exec in ep2 = %s, %v", out, err)
	}
	if again := d.workspace.registry("ep1")[broker.EffectWorkspaceExec]; again != h1 {
		t.Fatal("registry(ep1) rebuilt the exec handler instead of reusing it")
	}
	if n := f.dispatches(t); n != 0 {
		t.Fatalf("Workspace.Exec ran the AILANG tool binary %d times", n)
	}
}

// A handler that cannot be built is a recorded failure, never an unbound
// name (R8 would refuse the whole plan) and never a run.
func TestWorkspaceExecUnbuildableHandlerFails(t *testing.T) {
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	d.workspace.exec = &workspaceExec{profile: execProfile(t, filepath.Join(f.base, "no-such-read-root")),
		sandbox: stubExecSandbox(t, f.stateDir)}
	h := d.workspace.registry("ep1")[broker.EffectWorkspaceExec]
	if h == nil {
		t.Fatal("Workspace.Exec is unbound")
	}
	if out, err := execCall(t, h); err == nil || !strings.Contains(err.Error(), "no-such-read-root") {
		t.Fatalf("Execute = %s, %v; want the construction error", out, err)
	}
}

// ExecMaxTimeoutMS is HandlerCap − 1 000 ms (D-140-4 = A); the broker cannot
// import the coordinator, so the binding lives here.
func TestExecMaxTimeoutIsHandlerCapLessOneSecond(t *testing.T) {
	if want := (coordinator.HandlerCap - time.Second).Milliseconds(); int64(broker.ExecMaxTimeoutMS) != want {
		t.Fatalf("broker.ExecMaxTimeoutMS = %d, want HandlerCap − 1 s = %d", broker.ExecMaxTimeoutMS, want)
	}
}
