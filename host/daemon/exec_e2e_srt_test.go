package daemon

// Row 140 M3 (w-workspace-exec-toolchain-effect §6 AC3.2–AC3.4) end to end:
// a daemon started with `--exec-*` configuration (the real pinned srt, so the
// startup probe runs for real) serves workspace-exec over /mcp/ and /a2a/.
// The plan runs in the pinned interpreter (AILANG_BIN); the AILANG tool
// binary is the fake (exec never runs it). These tests skip, by name, only
// when AILANG_BIN, the srt install or node is absent; CI's go job runs them
// verbosely and fails unless each prints --- PASS.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// newExecSrtRig serves the published se-tools from a daemon whose exec
// profile names three commands in the fixture worktree ep1.
func newExecSrtRig(t *testing.T) *seRig {
	t.Helper()
	interp := os.Getenv("AILANG_BIN")
	if interp == "" {
		t.Skipf("SKIP %s: AILANG_BIN unset: the workspace-exec plan runs in the pinned %s", t.Name(), seInterpreterRelease)
	}
	nm := realExecSrt(t)
	node := execTestNode(t)
	f := newWSFixture(t)
	home := filepath.Join(f.base, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	profile := map[string]any{"profile": broker.ExecProfileVersion, "project": "proj", "path": []any{"/usr/bin", "/bin"},
		"timeout_ms": 5000, "probe": []any{"sh", "-c", "echo toolchain-ok"},
		"commands": map[string]any{
			"test":  map[string]any{"argv": []any{"printf", "[%s]"}, "flags": map[string]any{"-run": "regex"}, "positional": "relpath", "max_args": 4},
			"where": map[string]any{"argv": []any{"sh", "-c", "pwd -P"}},
			"mark":  map[string]any{"argv": []any{"sh", "-c", "echo ran >> ran.txt; echo marked"}},
		}}
	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	pf := filepath.Join(f.outside, "profiles", "proj.json")
	if err := os.MkdirAll(filepath.Dir(pf), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pf, data, 0o600); err != nil {
		t.Fatal(err)
	}
	d := mustWSDaemon(t, Config{DBPath: f.db, AilangBin: interp, WorkspaceRoot: f.root,
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease),
		ExecProfiles:  []string{pf}, ExecSandbox: nm, ExecNode: node,
		// Row 153 AC1.6: a refusal's operator line reaches the test log.
		ErrorLog: testErrorLog(t)})
	if d.coord == nil || d.workspace.exec == nil {
		t.Fatal("the daemon did not configure the coordinator and workspace-exec")
	}
	seGenesis(t, d)
	publishSeTools(t, d, f.db)
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(srv.Close)
	return &seRig{t: t, f: f, d: d, srv: srv, client: &http.Client{Timeout: DefaultClientTimeout}}
}

// callOK is r.call for a call that must reach the tool: a JSON-RPC error on
// the wire fails here, with the error's code and message (row 153 AC1.6,
// P14) rather than as a downstream "<nil>" field mismatch.
func (r *seRig) callOK(token, name string, args map[string]any) (mcpWire, time.Duration) {
	r.t.Helper()
	wire, took := r.call(token, name, args)
	if wire.Error != nil {
		r.t.Fatalf("%s: wire error %+v", name, *wire.Error)
	}
	return wire, took
}

// errText renders the JSON-RPC error by value: %+v of the *struct field
// prints an address, which is what hid refusals before row 153.
func (w mcpWire) errText() string {
	if w.Error == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *w.Error)
}

// execSpawns is the ep1 exec handler's spawn count (0 before it is built).
func (r *seRig) execSpawns() int64 {
	r.t.Helper()
	r.d.workspace.mu.Lock()
	defer r.d.workspace.mu.Unlock()
	cached, ok := r.d.workspace.execH["ep1"]
	if !ok {
		return 0
	}
	h, ok := cached.h.(*broker.ExecHandler)
	if !ok {
		r.t.Fatalf("ep1's Workspace.Exec handler is %T, want *broker.ExecHandler", cached.h)
	}
	return h.Spawns()
}

// execRequest decodes the Workspace.Exec payload the record of out's one
// effect names: the request object is length-prefixed fields ending in the
// canonical payload JSON.
func (r *seRig) execRequest(out map[string]any) (broker.EffectRecord, map[string]any) {
	r.t.Helper()
	rec, _ := r.effectRequestAny(out)
	obj, ok, err := r.d.store.GetObject(boundedTestContext(r.t), rec.RequestRef)
	if err != nil || !ok {
		r.t.Fatalf("request object %s: ok=%v err=%v", rec.RequestRef, ok, err)
	}
	raw := string(obj.Payload)
	if !strings.HasPrefix(raw, fmt.Sprintf("%d:%s", len(rec.Effect), rec.Effect)) || strings.Index(raw, "{") < 0 {
		r.t.Fatalf("request object %q does not name %s", raw, rec.Effect)
	}
	var req map[string]any
	if err := json.Unmarshal([]byte(raw[strings.Index(raw, "{"):]), &req); err != nil {
		r.t.Fatalf("request object %s: %v", raw, err)
	}
	return rec, req
}

// AC3.2: a workspace-exec call over /mcp/ runs the profiled command under
// the real srt and commits one entry whose record names Workspace.Exec and
// whose request holds exactly {command, args}.
func TestExecSrtMCPEndToEnd(t *testing.T) {
	r := newExecSrtRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceExec)
	if got := r.toolsList(token); !reflect.DeepEqual(got, []string{"workspace-exec"}) {
		t.Fatalf("tools/list = %v, want [workspace-exec]", got)
	}
	before := r.entryCount()
	args := []string{"-run", "TestX", "a"}
	wire, took := r.callOK(token, "workspace-exec", map[string]any{"command": "test", "args": args})
	out := wire.Result.StructuredContent
	t.Logf("workspace-exec test: %v (%d ms)", out, took.Milliseconds())
	if wire.Error != nil || wire.Result.IsError || out["exit_code"] != float64(0) || out["stdout"] != "[-run=TestX][a]" {
		t.Fatalf("workspace-exec = %+v (error %s), want exit 0 and the profiled argv's output", wire.Result, wire.errText())
	}
	sandbox, _ := out["sandbox"].(map[string]any)
	if sandbox["version"] != broker.ExecSandboxRelease || sandbox["cli_sha256"] != broker.ExecSandboxCLISHA256 ||
		sandbox["read_fence"] != true || sandbox["network"] != "none" {
		t.Fatalf("sandbox block %v, want the pinned srt with the read fence and no network", sandbox)
	}
	r.assertEffectRecord("workspace-exec", out, broker.EffectWorkspaceExec)
	if after := r.entryCount(); after != before+1 {
		t.Fatalf("log entries %d -> %d, want one commit", before, after)
	}
	rec, req := r.execRequest(out)
	if rec.Effect != broker.EffectWorkspaceExec || rec.Scope != broker.WorkspaceScope || rec.Cost != 1 {
		t.Fatalf("record = %+v", rec)
	}
	var payload map[string]any
	for _, v := range []any{req, req["payload"]} { // the payload, at the top or wrapped
		if m, ok := v.(map[string]any); ok && m["command"] != nil {
			payload = m
		}
	}
	if payload == nil || !reflect.DeepEqual(keysOfAny(payload), []string{"args", "command"}) || payload["command"] != "test" ||
		fmt.Sprint(payload["args"]) != fmt.Sprint(args) {
		t.Fatalf("request object %v, want exactly {command:test, args:%q}", req, args)
	}
	// The cwd is the episode's worktree, through the sandbox.
	wire, _ = r.callOK(token, "workspace-exec", map[string]any{"command": "where"})
	if got := strings.TrimSpace(fmt.Sprint(wire.Result.StructuredContent["stdout"])); got != filepath.Join(r.f.root, "ep1") {
		t.Fatalf("where = %q, want the ep1 worktree", got)
	}
	if n := r.execSpawns(); n < 2 {
		t.Fatalf("the exec handler spawned %d times for two runs", n)
	}
}

func keysOfAny(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func sortStrings(xs []string) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j-1] > xs[j]; j-- {
			xs[j-1], xs[j] = xs[j], xs[j-1]
		}
	}
}

// AC3.3: without a Workspace.Exec grant the call is refused and nothing
// runs or commits; with a budget of 1 the second call is denied:budget,
// recorded, and spawns nothing.
func TestExecSrtMCPGrantAndBudget(t *testing.T) {
	r := newExecSrtRig(t)
	ran := filepath.Join(r.f.root, "ep1", "ran.txt")
	noGrant := r.mint("ep1", broker.EffectWorkspaceRead)
	before := r.entryCount()
	wire, _ := r.call(noGrant, "workspace-exec", map[string]any{"command": "mark"})
	if wire.Error == nil && !wire.Result.IsError {
		t.Fatalf("workspace-exec without a Workspace.Exec grant = %+v (error %s), want refused", wire.Result, wire.errText())
	}
	if _, err := os.Stat(ran); !os.IsNotExist(err) || r.execSpawns() != 0 || r.entryCount() != before {
		t.Fatalf("the ungranted call ran (%v), spawned %d or committed (%d -> %d)", err, r.execSpawns(), before, r.entryCount())
	}

	one := r.mintBudgets("ep1", map[string]int64{broker.EffectWorkspaceExec: 1})
	wire, _ = r.callOK(one, "workspace-exec", map[string]any{"command": "mark"})
	if out := wire.Result.StructuredContent; out["exit_code"] != float64(0) || strings.TrimSpace(fmt.Sprint(out["stdout"])) != "marked" {
		t.Fatalf("first call with budget 1 = %+v (error %s)", wire.Result, wire.errText())
	}
	spawned := r.execSpawns()
	if spawned == 0 || mustReadFile(t, ran) != "ran\n" {
		t.Fatalf("the first call did not run (spawns %d)", spawned)
	}
	before = r.entryCount()
	wire, _ = r.callOK(one, "workspace-exec", map[string]any{"command": "mark"})
	out := wire.Result.StructuredContent
	if w := worldOf(t, out); len(w.Effects) != 1 || w.Effects[0].Status != "denied" {
		t.Fatalf("second call with budget 1 = %+v (error %s), want the effect denied", wire.Result, wire.errText())
	}
	rec, _ := r.effectRequestAny(out)
	if rec.Allowed || !strings.Contains(rec.Denial, "budget") || rec.Effect != broker.EffectWorkspaceExec {
		t.Fatalf("second call record = %+v, want Workspace.Exec denied:budget", rec)
	}
	if r.execSpawns() != spawned || mustReadFile(t, ran) != "ran\n" {
		t.Fatalf("the budget-denied call ran: spawns %d -> %d, ran.txt %q", spawned, r.execSpawns(), mustReadFile(t, ran))
	}
	if after := r.entryCount(); after != before+1 {
		t.Fatalf("the denied call committed %d entries, want 1 (a denial is recorded)", after-before)
	}
}

// AC3.4: a committed workspace-exec invocation replays through a binder
// with no registry — nothing executes — to the committed output byte for
// byte (row 134 V10/V46).
func TestExecSrtReplayIsByteEqual(t *testing.T) {
	r := newExecSrtRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceExec)
	code, raw := r.post("/a2a/", token, a2aSendBody("row140-replay", "workspace-exec", map[string]any{"command": "mark"}))
	var wire struct {
		Result struct {
			Metadata struct {
				InvocationID string `json:"invocation_id"`
			} `json:"metadata"`
			Artifacts []struct {
				Parts []struct {
					Data map[string]any `json:"data"`
				} `json:"parts"`
			} `json:"artifacts"`
		} `json:"result"`
	}
	if code != http.StatusOK || json.Unmarshal(raw, &wire) != nil || wire.Result.Metadata.InvocationID == "" ||
		len(wire.Result.Artifacts) != 1 || len(wire.Result.Artifacts[0].Parts) != 1 {
		t.Fatalf("tasks/send = %d %s", code, raw)
	}
	live := wire.Result.Artifacts[0].Parts[0].Data
	if live["exit_code"] != float64(0) {
		t.Fatalf("live run %v", live)
	}
	spawned := r.execSpawns()
	ran := filepath.Join(r.f.root, "ep1", "ran.txt")
	before := mustReadFile(t, ran)
	replayed, err := r.d.coord.Replay(boundedTestContext(t), wire.Result.Metadata.InvocationID,
		func(grants []broker.Capability, refs []hashref.HashRef) transitionreg.Binder {
			return broker.OpenReplayBinder(r.d.store, grants, refs)
		})
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	var replayedOut map[string]any
	if err := json.Unmarshal(replayed, &replayedOut); err != nil || !reflect.DeepEqual(replayedOut, live) {
		t.Fatalf("replayed %s, live %v (%v)", replayed, live, err)
	}
	liveBytes, _ := json.Marshal(live)
	again, _ := json.Marshal(replayedOut)
	if string(liveBytes) != string(again) {
		t.Fatalf("replay is not byte-equal:\n%s\n%s", again, liveBytes)
	}
	if r.execSpawns() != spawned || mustReadFile(t, ran) != before {
		t.Fatalf("the replay executed: spawns %d -> %d", spawned, r.execSpawns())
	}
}
