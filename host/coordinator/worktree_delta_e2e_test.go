package coordinator

// Row 184 r2 (judge item 3): a one-effect, no-finish-phase toolchain
// transition's COMMITTED output carries the handler's worktree_delta
// (composeEffectOutput's one-effect rule copies the handler fields), so a
// walk's `object get <output>` sees which files the call changed. The handler
// is the real broker.ExecHandler, run through a pass-through srt stub.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

func stubExecHandler(t *testing.T, worktree string) broker.Handler {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	state := filepath.Join(base, "state")
	pkg := filepath.Join(base, "node_modules", "@anthropic-ai", "sandbox-runtime")
	for _, d := range []string{filepath.Join(pkg, "dist"), state} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cli := "// stub cli\n"
	mustWriteFile(t, filepath.Join(pkg, "package.json"), `{"name":"@anthropic-ai/sandbox-runtime","version":"0.0.0-test"}`)
	mustWriteFile(t, filepath.Join(pkg, "dist", "cli.js"), cli)
	node := filepath.Join(base, "stub-node")
	if err := os.WriteFile(node, []byte("#!/bin/sh\nshift 4\nexec \"$@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(cli))
	sb, err := broker.ArchiveExecSandbox(broker.ExecSandboxConfig{NodeModules: filepath.Join(base, "node_modules"), Node: node,
		StateDir: state, Pin: broker.ExecSandboxPin{Version: "0.0.0-test", CLISHA256: hex.EncodeToString(sum[:])}})
	if err != nil {
		t.Fatal(err)
	}
	prof, _ := json.Marshal(map[string]any{"profile": broker.ExecProfileVersion, "project": "proj", "path": []any{"/usr/bin", "/bin"},
		"timeout_ms": 5000, "commands": map[string]any{"mk": map[string]any{"argv": []any{"/bin/sh", "-c", "printf made > made.txt"}}}})
	p, err := broker.ParseExecProfile(prof)
	if err != nil {
		t.Fatal(err)
	}
	h, err := broker.NewExecHandler(broker.ExecHandlerConfig{Profile: p, Sandbox: sb, Episode: fxEp, Worktree: worktree,
		WorkspaceRoot: filepath.Dir(worktree), StateDir: state, OperatorHome: base})
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCommittedOutputCarriesTheWorktreeDelta(t *testing.T) {
	execReq := transitionreg.EffectRequirement{Effect: broker.EffectWorkspaceExec, Scope: fxScope, Cost: 1}
	r, _ := fxRig(t, false, execReq)
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	worktree := filepath.Join(base, "ws", fxEp)
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	plan := `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Exec","scope":"worktree","cost":1,` +
		`"payload":{"command":"mk"}}],"finish":false,"result":null}`
	c := r.fxCoordinator(r.st, &phaseRunner{plan: plan},
		broker.Registry{broker.EffectWorkspaceExec: stubExecHandler(t, worktree), fxRead: &probe{}})
	grants := append(fxGrants(5), broker.Capability{Effect: broker.EffectWorkspaceExec, Scope: fxScope, ExpiresAt: 1000, Budget: 5})
	res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t-delta", grants))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	// The committed output object, read back from the store by the record's ref.
	rec := r.recordV2(res.RecordRef)
	obj, ok, err := r.st.GetObject(boundedTestContext(t), hashref.MustParse(rec.Output))
	if err != nil || !ok || obj.SemanticID != OutputV1 {
		t.Fatalf("output object ok=%v err=%v sid=%q", ok, err, obj.SemanticID)
	}
	var out struct {
		Delta *broker.WorktreeDelta `json:"worktree_delta"`
		World json.RawMessage       `json:"world"`
	}
	if err := json.Unmarshal(obj.Payload, &out); err != nil {
		t.Fatal(err)
	}
	if out.Delta == nil || out.World == nil || out.Delta.Total != 1 || len(out.Delta.Changes) != 1 {
		t.Fatalf("committed output %s lacks the worktree_delta", obj.Payload)
	}
	ch := out.Delta.Changes[0]
	sum := sha256.Sum256([]byte("made"))
	if ch.Path != "made.txt" || ch.Op != "created" || ch.Size == nil || *ch.Size != 4 || ch.SHA256 != hex.EncodeToString(sum[:]) ||
		out.Delta.Basis != broker.WorktreeBasisWalk {
		t.Fatalf("committed delta = %+v", out.Delta)
	}
	if string(obj.Payload) != string(res.OutputBytes) {
		t.Fatal("the caller's output bytes are not the committed output")
	}
}
