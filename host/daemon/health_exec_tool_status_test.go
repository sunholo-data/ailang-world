package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestHealthExecToolStatusShowsTargetAndProbe(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	stateDir, root, home := filepath.Join(base, "state"), filepath.Join(base, "workspace"), filepath.Join(base, "home")
	tree := filepath.Join(base, "node_modules")
	pkg := filepath.Join(tree, "@anthropic-ai", "sandbox-runtime")
	for _, dir := range []string{stateDir, filepath.Join(root, "ep1"), home, filepath.Join(pkg, "dist")} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(p, body string) {
		t.Helper()
		if err := os.WriteFile(p, []byte(body), 0755); err != nil {
			t.Fatal(err)
		}
	}
	cli := "// fake srt"
	write(filepath.Join(pkg, "package.json"), `{"name":"@anthropic-ai/sandbox-runtime","version":"0.0.0-test"}`)
	write(filepath.Join(pkg, "dist", "cli.js"), cli)
	a, b, c, link := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "c"), filepath.Join(base, "node")
	for _, p := range []string{a, b, c} {
		write(p, "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo v26.12.1; fi\nexit 0\n")
	}
	repoint := func(target string) {
		t.Helper()
		if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
	}
	repoint(a)
	sum := sha256.Sum256([]byte(cli))
	sb, err := broker.ArchiveExecSandbox(broker.ExecSandboxConfig{NodeModules: tree, Node: a, StateDir: stateDir,
		Pin: broker.ExecSandboxPin{Version: "0.0.0-test", CLISHA256: hex.EncodeToString(sum[:])}})
	if err != nil {
		t.Fatal(err)
	}
	raw := []byte(`{"profile":"world-exec-profile/v1","project":"proj","root":".","path":["/usr/bin","/bin"],"timeout_ms":9000,"commands":{"test":{"argv":["/bin/sh","-c","exit 0"]}}}`)
	// Use the exported profile version, keeping this fixture independent of a schema bump.
	raw = []byte(strings.Replace(string(raw), "world-exec-profile/v1", broker.ExecProfileVersion, 1))
	p, err := broker.ParseExecProfile(raw)
	if err != nil {
		t.Fatal(err)
	}
	fail := false
	toolState, err := broker.NewExecToolState(broker.ExecToolStateConfig{Node: link, Sandbox: sb,
		Profiles:    map[string]*broker.ExecProfile{"proj": p},
		ProbeConfig: broker.ExecProbeConfig{Paths: broker.ExecHostPaths{StateDir: stateDir, WorkspaceRoot: root, OperatorHome: home}},
		Probe: func(context.Context, broker.ExecProbeConfig, []*broker.ExecProfile) ([]error, error) {
			if fail {
				return []error{errors.New("arm7 fixture")}, nil
			}
			return []error{nil}, nil
		}})
	if err != nil {
		t.Fatal(err)
	}
	d := &Daemon{workspace: &workspaceTools{exec: &workspaceExec{toolState: toolState}}}
	health := func(target, accepted, probe string) {
		t.Helper()
		rec := httptest.NewRecorder()
		d.handleHealth(rec, httptest.NewRequest("GET", "/v1/health", nil))
		var wire struct {
			Tools  map[string]string                `json:"exec_tools"`
			Status map[string]broker.ExecToolStatus `json:"exec_tool_status"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &wire); err != nil {
			t.Fatal(err)
		}
		s := wire.Status["node"]
		if wire.Tools["node"] != accepted || s.Configured != link || s.Target != target || s.Accepted != accepted ||
			!strings.HasPrefix(s.Probe, probe) || s.CheckedAt == "" {
			t.Fatalf("health: %s", rec.Body.Bytes())
		}
	}
	health(a, a, "ok")
	repoint(b)
	health(b, a, "pending")
	h, err := broker.NewExecHandler(broker.ExecHandlerConfig{Tools: toolState, Profile: p, Sandbox: sb, Episode: "ep1",
		Worktree: filepath.Join(root, "ep1"), WorkspaceRoot: root, StateDir: stateDir, OperatorHome: home})
	if err != nil {
		t.Fatal(err)
	}
	call := func() []byte {
		t.Helper()
		ctx := boundedTestContext(t)
		out, err := h.Execute(ctx, broker.EffectRequest{Effect: broker.EffectWorkspaceExec, Scope: broker.WorkspaceScope}, []byte(`{"command":"test"}`))
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	if out := call(); strings.Contains(string(out), "refused") {
		t.Fatal(string(out))
	}
	health(b, b, "ok")
	fail = true
	repoint(c)
	if out := call(); !strings.Contains(string(out), "refused") {
		t.Fatal(string(out))
	}
	health(c, b, "failed: ")
}
