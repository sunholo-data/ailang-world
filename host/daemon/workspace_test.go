package daemon

// Row 134 M4b: the production handler registry (w-software-engineering-domain
// §4.3) — workspaceRegistry(episodeID), the --workspace-root/--tool-ailang-bin
// startup checks (AC4.4, AC4.5), the sibling-episode sandbox
// (MUT-SANDBOX-ROOT) and one Workspace.Read effect through the daemon's own
// Binder.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const toolBinEnv = "WORLD_TOOL_AILANG_BIN"

// ---------------------------------------------------------------------------
// fixture
// ---------------------------------------------------------------------------

// wsFixture is <base>/state/world.db (the daemon's state) beside
// <base>/ws (the workspace root) holding ep1, ep2 and the escape shapes.
type wsFixture struct {
	base, db, stateDir, root, outside string
	logDir                            string
}

func canonicalTemp(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func newWSFixture(t *testing.T) wsFixture {
	t.Helper()
	base := canonicalTemp(t)
	f := wsFixture{
		base: base, stateDir: filepath.Join(base, "state"), root: filepath.Join(base, "ws"),
		outside: filepath.Join(base, "outside"), logDir: filepath.Join(base, "toollog"),
	}
	f.db = filepath.Join(f.stateDir, "world.db")
	for _, d := range []string{f.stateDir, f.outside, f.logDir,
		filepath.Join(f.root, "ep1"), filepath.Join(f.root, "ep2"), filepath.Join(f.root, "EP1")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(f.root, "ep1", "data.txt"), "hello from ep1\n")
	writeFile(t, filepath.Join(f.root, "ep2", "secret.txt"), "secret of ep2\n")
	writeFile(t, filepath.Join(f.root, "plainfile"), "not a directory\n")
	for link, target := range map[string]string{
		"escape": f.outside,                    // a symlink out of the root
		"alias":  filepath.Join(f.root, "ep2"), // a symlink onto a sibling episode
		"self":   f.root,                       // a symlink onto the root itself
	} {
		if err := os.Symlink(target, filepath.Join(f.root, link)); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// fakeToolBin is a /bin/sh stand-in for the tool binary, archived like the
// real one. --version names release; a policy-tool summary answers
// restricted with fs_sandbox = the cwd and logs its environment; every other
// launch is counted in <logDir>/dispatches (absolute: the archive copies the
// script elsewhere).
func fakeToolBin(t *testing.T, logDir, release string) string {
	t.Helper()
	script := fmt.Sprintf(`#!/bin/sh
log=%q
if [ "$1" = "--version" ]; then echo %q; exit 0; fi
if [ "$1" = "policy-tool" ]; then
  req=$(cat)
  case "$req" in
    *'"summary"'*)
      echo x >> "$log/summaries"
      env > "$log/env"
      printf '{"ok":true,"summary":{"security_mode":"restricted","policy_digest":"fakedigest","fs_sandbox":"%%s","cli":["check"]}}' "$(pwd -P)"
      exit 0 ;;
  esac
  echo x >> "$log/dispatches"
  printf '{"ok":true,"content":"fake:%%s"}' "$(pwd -P)"
  exit 0
fi
echo x >> "$log/dispatches"
exit 0
`, logDir, release)
	bin := filepath.Join(canonicalTemp(t), "ailang")
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return bin
}

func countLines(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(data), "x")
}

func (f wsFixture) dispatches(t *testing.T) int {
	return countLines(t, filepath.Join(f.logDir, "dispatches"))
}
func (f wsFixture) summaries(t *testing.T) int {
	return countLines(t, filepath.Join(f.logDir, "summaries"))
}

func newWSDaemon(t *testing.T, cfg Config) (*Daemon, error) {
	t.Helper()
	cfg.BindHost = DefaultBindHost
	if cfg.ErrorLog == nil {
		cfg.ErrorLog = io.Discard
	}
	d, err := New(boundedTestContext(t), cfg)
	if err == nil {
		t.Cleanup(func() {
			if err := d.Close(); err != nil {
				t.Error(err)
			}
		})
	}
	return d, err
}

func mustWSDaemon(t *testing.T, cfg Config) *Daemon {
	t.Helper()
	d, err := newWSDaemon(t, cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return d
}

func registryNames(reg broker.Registry) []string {
	names := make([]string, 0, len(reg))
	for name := range reg {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// wantWorkspaceEffects is the nine names (row 135 AC3: 6 -> 8; row 140 AC1.3:
// 8 -> 9, Workspace.Exec; always bound).
var wantWorkspaceEffects = []string{"Ailang.CLI", "Ailang.Check", "Ailang.Discover", "Ailang.Run", "Ailang.RunEnv",
	"Ailang.RunNet", "Workspace.Exec", "Workspace.Read", "Workspace.Write"}

// refusedEpisodes are AC4.4's shapes: each must yield an empty registry.
var refusedEpisodes = []string{
	"../x", "..", ".", "", "EP1", "Ep1", "-ep1", "ep1/../ep2", "ep1/sub", "/abs", "ep_1", "ep1 ",
	strings.Repeat("a", 65),
	"escape",    // symlink out of the root
	"alias",     // symlink onto the sibling episode ep2
	"self",      // symlink onto the root itself
	"plainfile", // not a directory
	"absent",    // no worktree provisioned
}

// ---------------------------------------------------------------------------
// AC4.4 — workspaceRegistry
// ---------------------------------------------------------------------------

func TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot(t *testing.T) {
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	for _, ep := range refusedEpisodes {
		if reg := d.workspace.registry(ep); len(reg) != 0 {
			t.Errorf("registry(%q) = %v, want empty", ep, registryNames(reg))
		}
	}
	if n := f.summaries(t); n != 0 {
		t.Fatalf("a refused episode constructed a handler (%d summary runs)", n)
	}
	// Positive control: the provisioned episode gets the nine effect names.
	reg := d.workspace.registry("ep1")
	if got := registryNames(reg); fmt.Sprint(got) != fmt.Sprint(wantWorkspaceEffects) {
		t.Fatalf("registry(ep1) = %v, want %v", got, wantWorkspaceEffects)
	}
	if n := f.dispatches(t); n != 0 {
		t.Fatalf("tool dispatched %d times without a request", n)
	}
	// One handler per episode, cached: no second summary subprocess.
	_ = d.workspace.registry("ep1")
	if n := f.summaries(t); n != 1 {
		t.Fatalf("summary runs = %d after two registry(ep1) calls, want 1 (cached)", n)
	}
}

// TestWorkspaceExecUnconfiguredSpawnsNothing is row 140 M1's binding: the
// Workspace.Exec name is always bound (R8: workspace-exec declares it), and
// with no exec profile configured its handler answers the typed refusal
// without running the tool binary or anything else.
func TestWorkspaceExecUnconfiguredSpawnsNothing(t *testing.T) {
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	h := d.workspace.registry("ep1")[broker.EffectWorkspaceExec]
	if h == nil {
		t.Fatal("registry(ep1) does not bind Workspace.Exec")
	}
	out, err := h.Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectWorkspaceExec, Scope: broker.WorkspaceScope, Cost: 1},
		[]byte(`{"command":"test"}`))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil || got["ok"] != false || got["refused"] != broker.NoExecProfileRefusal || len(got) != 2 {
		t.Fatalf("Workspace.Exec without a profile = %s (%v), want ok:false refused %q", out, err, broker.NoExecProfileRefusal)
	}
	if n := f.dispatches(t); n != 0 {
		t.Fatalf("the unconfigured refusal ran the tool binary %d times", n)
	}
}

// TestWorkspaceRegistryNeedsBothFlags: either flag alone serves nothing and
// touches no per-episode state — no rendered policy, no cache directory, no
// tool launch. The process cwd holds an "ep1" directory so that a registry
// resolving episodes against an unset (empty) root would find one.
func TestWorkspaceRegistryNeedsBothFlags(t *testing.T) {
	f := newWSFixture(t)
	cwd := canonicalTemp(t)
	if err := os.MkdirAll(filepath.Join(cwd, "ep1"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(cwd)
	noEpisodeState := func(arm string, dirs ...string) {
		t.Helper()
		for _, dir := range dirs {
			for _, sub := range []string{"policies", "cache"} {
				if _, err := os.Lstat(filepath.Join(dir, sub)); err == nil {
					t.Fatalf("%s: per-episode state %s was created", arm, filepath.Join(dir, sub))
				}
			}
		}
		if n := f.summaries(t) + f.dispatches(t); n != 0 {
			t.Fatalf("%s: a half-configured daemon launched the tool %d times", arm, n)
		}
	}

	var log strings.Builder
	rootOnly := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ErrorLog: &log})
	if reg := rootOnly.workspace.registry("ep1"); len(reg) != 0 {
		t.Fatalf("--workspace-root alone: registry(ep1) = %v, want empty", registryNames(reg))
	}
	if !strings.Contains(log.String(), "workspace tools disabled") || strings.Contains(log.String(), "unavailable for episode") {
		t.Fatalf("root-only startup log = %q, want only the disabled notice", log.String())
	}
	noEpisodeState("--workspace-root alone", f.stateDir, cwd)
	if err := rootOnly.Close(); err != nil {
		t.Fatal(err)
	}

	log.Reset()
	toolOnly := mustWSDaemon(t, Config{DBPath: f.db, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease), ErrorLog: &log})
	if reg := toolOnly.workspace.registry("ep1"); len(reg) != 0 {
		t.Fatalf("--tool-ailang-bin alone: registry(ep1) = %v, want empty", registryNames(reg))
	}
	if !strings.Contains(log.String(), "workspace tools disabled") || strings.Contains(log.String(), "unavailable for episode") {
		t.Fatalf("tool-only startup log = %q, want only the disabled notice", log.String())
	}
	noEpisodeState("--tool-ailang-bin alone", f.stateDir, cwd)
}

// TestWorkspacePolicyRenderedPerEpisode: the episode's handler runs the
// archived binary against <db-dir>/policies/<ep>.toml (0600, fs_sandbox =
// root/ep, never the workspace root) with AILANG_CACHE_DIR=<db-dir>/cache/<ep>.
func TestWorkspacePolicyRenderedPerEpisode(t *testing.T) {
	f := newWSFixture(t)
	configured := fakeToolBin(t, f.logDir, ToolBinaryRelease)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: configured})
	if len(d.workspace.registry("ep1")) == 0 {
		t.Fatal("registry(ep1) is empty")
	}
	epRoot := filepath.Join(f.root, "ep1")
	policyPath := filepath.Join(f.stateDir, "policies", "ep1.toml")
	got, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatal(err)
	}
	want, err := broker.RenderEpisodePolicy(epRoot)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("rendered policy =\n%s\nwant (fs_sandbox = %s)\n%s", got, epRoot, want)
	}
	if !strings.Contains(string(got), `fs_sandbox = "`+epRoot+`"`) {
		t.Fatalf("policy fs_sandbox is not the episode worktree: %s", got)
	}
	if info, err := os.Stat(policyPath); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("policy mode = %v (%v), want 0600", info.Mode().Perm(), err)
	}
	env, err := os.ReadFile(filepath.Join(f.logDir, "env"))
	if err != nil {
		t.Fatal(err)
	}
	if wantCache := "AILANG_CACHE_DIR=" + filepath.Join(f.stateDir, "cache", "ep1"); !strings.Contains(string(env), wantCache+"\n") {
		t.Fatalf("tool env lacks %s:\n%s", wantCache, env)
	}
	// The handler runs the archived bytes, never the configured path.
	if !strings.HasPrefix(d.workspace.bin, archive.New(f.db).Root()+string(filepath.Separator)) || d.workspace.bin == configured {
		t.Fatalf("tool binary path = %q, want a path inside the archive %s", d.workspace.bin, archive.New(f.db).Root())
	}
}

// TestWorkspaceToolBinaryIsArchivedAndVerified: a non-measured release is
// refused at startup; the archived bytes run even after the configured file
// is removed; bytes changed in the archive fail the call.
func TestWorkspaceToolBinaryIsArchivedAndVerified(t *testing.T) {
	f := newWSFixture(t)
	_, err := newWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, "AILANG v0.41.0")})
	var startup *StartupError
	if !errors.As(err, &startup) || startup.Stage != StageArchive || !strings.Contains(err.Error(), ToolBinaryRelease) {
		t.Fatalf("New with an unmeasured tool release = %v, want a %s StartupError naming %s", err, StageArchive, ToolBinaryRelease)
	}

	configured := fakeToolBin(t, f.logDir, ToolBinaryRelease)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: configured})
	if err := os.Remove(configured); err != nil {
		t.Fatal(err)
	}
	h := d.workspace.registry("ep1")[broker.EffectWorkspaceRead]
	req := broker.EffectRequest{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, Cost: 1}
	if out, err := h.Execute(boundedTestContext(t), req, []byte(`{"op":"read","path":"data.txt"}`)); err != nil || !strings.Contains(string(out), "fake:") {
		t.Fatalf("archived tool with the configured file removed: %s %v", out, err)
	}
	before := f.dispatches(t)
	if err := os.Chmod(d.workspace.bin, 0o755); err != nil {
		t.Fatal(err)
	}
	fh, err := os.OpenFile(d.workspace.bin, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fh.WriteString("# tampered\n"); err != nil {
		t.Fatal(err)
	}
	_ = fh.Close()
	if out, err := h.Execute(boundedTestContext(t), req, []byte(`{"op":"read","path":"data.txt"}`)); err == nil {
		t.Fatalf("tampered archived tool ran: %s", out)
	}
	if f.dispatches(t) != before {
		t.Fatal("the tampered binary was launched")
	}
}

// ---------------------------------------------------------------------------
// AC4.5 — the daemon's state outside the workspace root
// ---------------------------------------------------------------------------

func TestWorkspaceRootContainingDaemonStateRefusesStartup(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, f wsFixture) (root, db string)
	}{
		{"root is the store directory", func(t *testing.T, f wsFixture) (string, string) { return f.stateDir, f.db }},
		{"store under the root", func(t *testing.T, f wsFixture) (string, string) {
			return f.base, f.db
		}},
		{"store inside an episode worktree", func(t *testing.T, f wsFixture) (string, string) {
			return f.root, filepath.Join(f.root, "ep1", "world.db")
		}},
		{"root reached through a symlink", func(t *testing.T, f wsFixture) (string, string) {
			link := filepath.Join(f.outside, "link-to-base")
			if err := os.Symlink(f.base, link); err != nil {
				t.Fatal(err)
			}
			return link, f.db
		}},
		{"policy directory symlinked into the root", func(t *testing.T, f wsFixture) (string, string) {
			if err := os.Symlink(filepath.Join(f.root, "ep1"), filepath.Join(f.stateDir, "policies")); err != nil {
				t.Fatal(err)
			}
			return f.root, f.db
		}},
		{"tool cache symlinked into the root", func(t *testing.T, f wsFixture) (string, string) {
			if err := os.Symlink(filepath.Join(f.root, "ep2"), filepath.Join(f.stateDir, "cache")); err != nil {
				t.Fatal(err)
			}
			return f.root, f.db
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			root, db := tc.setup(t, f)
			_, err := newWSDaemon(t, Config{DBPath: db, WorkspaceRoot: root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
			var startup *StartupError
			var state *WorkspaceStateError
			var inside *broker.PolicyInsideWorkspaceError
			if !errors.As(err, &startup) || startup.Stage != StageConfig || !errors.As(err, &state) || !errors.As(err, &inside) {
				t.Fatalf("New = %v, want a %s StartupError wrapping *WorkspaceStateError and *broker.PolicyInsideWorkspaceError", err, StageConfig)
			}
			// Refused before writer authority was taken: the store opens.
			s, err := store.Open(db)
			if err != nil {
				t.Fatalf("store after the refusal: %v", err)
			}
			_ = s.Close()
		})
	}
	// Control: the fixture's own layout starts.
	f := newWSFixture(t)
	mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	// A root that is not a directory is refused too.
	if _, err := newWSDaemon(t, Config{DBPath: filepath.Join(f.outside, "w.db"), WorkspaceRoot: filepath.Join(f.root, "plainfile")}); err == nil {
		t.Fatal("a regular file was accepted as --workspace-root")
	}
}

// ---------------------------------------------------------------------------
// AC4.4 through the coordinator: R8 before the plan or any effect
// ---------------------------------------------------------------------------

// readPlanRunner answers the plan phase with one Workspace.Read of data.txt
// and counts how often it ran.
type readPlanRunner struct{ plans atomic.Int64 }

const wsReadPlan = `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Read","scope":"worktree","cost":1,"payload":{"op":"read","path":"data.txt"}}],"finish":false,"result":null}`

func (r *readPlanRunner) RunContext(_ context.Context, e capsule.Entry) (capsule.Result, error) {
	var arg string
	if err := json.Unmarshal(e.Args, &arg); err != nil {
		return capsule.Result{}, err
	}
	if !strings.Contains(arg, `"phase":"plan"`) {
		return capsule.Result{}, fmt.Errorf("unexpected phase input %s", arg)
	}
	r.plans.Add(1)
	return capsule.Result{Stdout: []byte(wsReadPlan + "\n")}, nil
}

type wsCaps []broker.Capability

func (c wsCaps) CapabilitySnapshot(n int64) broker.CapabilitySnapshot {
	return broker.NewCapabilitySnapshot(c, n)
}

var wsReadReq = transitionreg.EffectRequirement{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, Cost: 1}

// publishReadTool commits a genesis world and publishes "ws.read" (access
// Workspace.Read@worktree:0, declaring Workspace.Read@worktree:1) whose
// source is src, pinned to interp.
func publishReadTool(t *testing.T, d *Daemon, interp hashref.HashRef, src string) {
	t.Helper()
	ctx := boundedTestContext(t)
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "workspace_test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	if err := d.store.Commit(ctx, store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{
		EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: interp,
		PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "workspace_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}}); err != nil {
		t.Fatalf("genesis: %v", err)
	}
	source, err := canon.Source([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	srcObj := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
		SemanticID: "world/transition-source/v1", Provenance: "workspace_test", Payload: source}
	if err := d.store.PutObject(ctx, srcObj); err != nil {
		t.Fatal(err)
	}
	desc := transitionreg.Descriptor{ID: "ws.read", TransitionFn: srcObj.Hash, Interpreter: interp, SemanticsEpoch: 1,
		InputSchema: []byte(`{"type":"object"}`), OutputSchema: []byte(`{"type":"object"}`),
		Access:          transitionreg.EffectRequirement{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, Cost: 0},
		DeclaredEffects: []transitionreg.EffectRequirement{wsReadReq}, Title: "Read", Description: "workspace read fixture"}
	payload, err := transitionreg.EncodeRevision(transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1,
		InterfaceHash: transitionreg.InterfaceHashV1, Revision: 1, Entries: []transitionreg.Descriptor{desc}})
	if err != nil {
		t.Fatal(err)
	}
	rev := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: transitionreg.InterfaceHashV1,
		SemanticID: transitionreg.SemanticIDV1, Provenance: "workspace_test", Payload: payload}
	if err := d.store.PutObject(ctx, rev); err != nil {
		t.Fatal(err)
	}
	if err := d.store.CompareAndSetRegistryHead(ctx, store.TransitionRegistryV1, hashref.HashRef{}, rev.Hash); err != nil {
		t.Fatalf("registry head: %v", err)
	}
}

func readCall(t *testing.T, d *Daemon, episode, task string, now int64) coordinator.Call {
	t.Helper()
	grants := []broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, ExpiresAt: now + 3600, Budget: 5}}
	req, err := transitionreg.NewRequest(boundedTestContext(t), transitionreg.NewReader(d.store), wsCaps(grants), now)
	if err != nil {
		t.Fatal(err)
	}
	return coordinator.Call{Request: req, EpisodeID: episode, Grants: grants, SkillID: "ws.read", TaskID: task, Input: map[string]any{"path": "data.txt"}}
}

// TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect is AC4.4 end to end through
// the daemon's production Binder: a refused episode's empty registry makes
// the coordinator refuse (R8) before the plan runs, and the tool is never
// launched; the provisioned episode runs its one effect.
func TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect(t *testing.T) {
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	publishReadTool(t, d, hashref.SumSHA256([]byte("fake-interpreter")), "module host/capsule/main\n\nexport func main(input: string) -> string {\n  input\n}\n")
	runner := &readPlanRunner{}
	now := time.Now().Unix()
	coord, err := coordinator.New(coordinator.Config{Store: d.store, Runner: runner, Binder: d.binder,
		Now: func() int64 { return now }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	for i, ep := range []string{"../x", "EP1", "escape", "alias"} {
		_, err := coord.Dispatch(boundedTestContext(t), readCall(t, d, ep, fmt.Sprintf("t%d", i), now))
		var unsupported *coordinator.EffectsUnsupportedError
		if !errors.As(err, &unsupported) || fmt.Sprint(unsupported.Effects) != "[Workspace.Read]" {
			t.Fatalf("episode %q: Dispatch = %v, want R8 naming [Workspace.Read]", ep, err)
		}
	}
	if runner.plans.Load() != 0 || f.dispatches(t) != 0 || f.summaries(t) != 0 {
		t.Fatalf("refused episodes ran plans=%d tool dispatches=%d summaries=%d, want 0/0/0",
			runner.plans.Load(), f.dispatches(t), f.summaries(t))
	}
	res, err := coord.Dispatch(boundedTestContext(t), readCall(t, d, "ep1", "ok", now))
	if err != nil {
		t.Fatalf("ep1 Dispatch: %v", err)
	}
	if runner.plans.Load() != 1 || f.dispatches(t) != 1 || !strings.Contains(string(res.OutputBytes), `"content":"fake:`+filepath.Join(f.root, "ep1")+`"`) {
		t.Fatalf("ep1: plans=%d dispatches=%d output %s", runner.plans.Load(), f.dispatches(t), res.OutputBytes)
	}
}

// ---------------------------------------------------------------------------
// real tool binary (gated)
// ---------------------------------------------------------------------------

func realToolBin(t *testing.T) string {
	t.Helper()
	bin := os.Getenv(toolBinEnv)
	if bin == "" {
		t.Skipf("%s is unset: this real-binary test needs the %s tool binary "+
			"(e.g. $HOME/.pinned-ailang-tools/v0.52.1/ailang); skipping", toolBinEnv, ToolBinaryRelease)
	}
	return bin
}

// TestWorkspaceSandboxIsTheEpisodeNotTheRoot (MUT-SANDBOX-ROOT): through the
// daemon's registry, ep1 reads its own file and cannot read ep2's — neither
// by `..` nor by naming the sibling as if the root were the sandbox.
func TestWorkspaceSandboxIsTheEpisodeNotTheRoot(t *testing.T) {
	bin := realToolBin(t)
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: bin})
	h := d.workspace.registry("ep1")[broker.EffectWorkspaceRead]
	if h == nil {
		t.Fatal("registry(ep1) has no Workspace.Read handler")
	}
	read := func(path string) map[string]any {
		t.Helper()
		body, _ := json.Marshal(map[string]string{"op": "read", "path": path})
		out, err := h.Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, Cost: 1}, body)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var resp map[string]any
		if err := json.Unmarshal(out, &resp); err != nil {
			t.Fatal(err)
		}
		return resp
	}
	if resp := read("data.txt"); resp["ok"] != true || resp["content"] != "hello from ep1\n" {
		t.Fatalf("in-episode read = %v", resp)
	}
	for _, sibling := range []string{"../ep2/secret.txt", "ep2/secret.txt", filepath.Join(f.root, "ep2", "secret.txt")} {
		if resp := read(sibling); resp["ok"] != false || strings.Contains(fmt.Sprint(resp), "secret of ep2") {
			t.Fatalf("sibling-episode read %q = %v, want refused", sibling, resp)
		}
	}
}

const wsReadPlanSource = "module host/capsule/main\n\nexport func main(input: string) -> string {\n" +
	`  "{\"plan\":\"world/effect-plan/v1\",\"effects\":[{\"id\":\"e1\",\"effect\":\"Workspace.Read\",\"scope\":\"worktree\",\"cost\":1,\"payload\":{\"op\":\"read\",\"path\":\"data.txt\"}}],\"finish\":false,\"result\":null}"` +
	"\n}\n"

// TestWorkspaceReadThroughProductionWiring: one Workspace.Read effect through
// the daemon's own coordinator — the archived interpreter runs the plan, the
// daemon's Binder supplies the episode's registry, the archived tool binary
// reads the worktree, and the effect record and spend are durable.
func TestWorkspaceReadThroughProductionWiring(t *testing.T) {
	tool := realToolBin(t)
	interp := os.Getenv("AILANG_BIN")
	if interp == "" {
		t.Skip("AILANG_BIN is unset: this e2e needs the pinned interpreter for the plan phase; skipping")
	}
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, AilangBin: interp, WorkspaceRoot: f.root, ToolAilangBin: tool})
	if d.coord == nil {
		t.Fatal("no coordinator with --ailang-bin set")
	}
	publishReadTool(t, d, daemonInterpreterRef(t, d), wsReadPlanSource)
	now := time.Now().Unix()
	res, err := d.coord.Dispatch(boundedTestContext(t), readCall(t, d, "ep1", "e2e", now))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	var out struct {
		OK      bool   `json:"ok"`
		Content string `json:"content"`
		Tool    string `json:"tool"`
		World   struct {
			Effects []struct {
				ID, Status string
				Record     *string
			} `json:"effects"`
		} `json:"world"`
	}
	if err := json.Unmarshal(res.OutputBytes, &out); err != nil {
		t.Fatal(err)
	}
	if !out.OK || out.Content != "hello from ep1\n" || out.Tool != d.workspace.binRef.String() ||
		len(out.World.Effects) != 1 || out.World.Effects[0].Status != "ok" || out.World.Effects[0].Record == nil {
		t.Fatalf("output = %s", res.OutputBytes)
	}
	recObj, ok, err := d.store.GetObject(boundedTestContext(t), hashref.MustParse(*out.World.Effects[0].Record))
	if err != nil || !ok || recObj.SemanticID != broker.EffectRecordV1 {
		t.Fatalf("effect record ok=%v err=%v", ok, err)
	}
	rec, err := broker.DecodeRecord(recObj.Payload)
	if err != nil || !rec.Allowed || rec.Failed || rec.Effect != broker.EffectWorkspaceRead {
		t.Fatalf("effect record = %+v %v", rec, err)
	}
	spend, err := d.store.EffectSpend(boundedTestContext(t), "ep1")
	if err != nil || spend[store.EffectKey{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope}] != 1 {
		t.Fatalf("spend = %v %v, want Workspace.Read=1", spend, err)
	}
}

// TestWorkspaceExamplesDirReachesTheToolEnv (row 134 break-3 fix): a
// configured --examples-dir reaches the tool's minimal child env as
// AILANG_EXAMPLES=<canonical dir>; with none configured the variable is
// absent, and the rest of the env is unchanged.
func TestWorkspaceExamplesDirReachesTheToolEnv(t *testing.T) {
	envOf := func(t *testing.T, configure func(f wsFixture, cfg *Config)) (wsFixture, []string) {
		f := newWSFixture(t)
		cfg := Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)}
		configure(f, &cfg)
		d := mustWSDaemon(t, cfg)
		if len(d.workspace.registry("ep1")) == 0 {
			t.Fatal("registry(ep1) is empty")
		}
		raw, err := os.ReadFile(filepath.Join(f.logDir, "env"))
		if err != nil {
			t.Fatal(err)
		}
		var keys []string
		for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			if k, _, ok := strings.Cut(line, "="); ok && k != "PWD" && k != "SHLVL" && k != "_" && k != "OLDPWD" {
				keys = append(keys, line)
			}
		}
		sort.Strings(keys)
		return f, keys
	}
	f, with := envOf(t, func(f wsFixture, cfg *Config) {
		corpus := filepath.Join(f.outside, "examples")
		if err := os.MkdirAll(corpus, 0o755); err != nil {
			t.Fatal(err)
		}
		cfg.ExamplesDir = corpus
	})
	cache := filepath.Join(f.stateDir, "cache", "ep1")
	want := []string{"AILANG_CACHE_DIR=" + cache, "AILANG_EXAMPLES=" + filepath.Join(f.outside, "examples"),
		"HOME=" + cache, "LANG=C", "LC_ALL=C", "PATH=/usr/bin:/bin"}
	if strings.Join(with, "\n") != strings.Join(want, "\n") {
		t.Fatalf("tool env with --examples-dir =\n%s\nwant\n%s", strings.Join(with, "\n"), strings.Join(want, "\n"))
	}
	_, without := envOf(t, func(wsFixture, *Config) {})
	for _, kv := range without {
		if strings.HasPrefix(kv, "AILANG_EXAMPLES=") {
			t.Fatalf("tool env without --examples-dir carries %s", kv)
		}
	}
	if len(without) != len(want)-1 {
		t.Fatalf("tool env without --examples-dir = %v, want the %d minimal keys", without, len(want)-1)
	}
}

// TestWorkspaceExamplesDirInsideTheRootRefusesStartup: an examples corpus an
// agent can write (inside --workspace-root, directly or through a symlink) is
// refused at startup like the policy and cache dirs; so is one that is not a
// directory. A corpus outside the root starts.
func TestWorkspaceExamplesDirInsideTheRootRefusesStartup(t *testing.T) {
	for _, tc := range []struct {
		name  string
		dir   func(t *testing.T, f wsFixture) string
		state bool // wraps *WorkspaceStateError and *broker.PolicyInsideWorkspaceError
	}{
		{"corpus inside an episode worktree", func(t *testing.T, f wsFixture) string {
			dir := filepath.Join(f.root, "ep1", "examples")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			return dir
		}, true},
		{"corpus is the root", func(t *testing.T, f wsFixture) string { return f.root }, true},
		{"corpus symlinked into the root", func(t *testing.T, f wsFixture) string {
			link := filepath.Join(f.outside, "examples-link")
			if err := os.Symlink(filepath.Join(f.root, "ep2"), link); err != nil {
				t.Fatal(err)
			}
			return link
		}, true},
		{"corpus is a file", func(t *testing.T, f wsFixture) string {
			p := filepath.Join(f.outside, "examples.txt")
			writeFile(t, p, "x\n")
			return p
		}, false},
		{"corpus is absent", func(t *testing.T, f wsFixture) string { return filepath.Join(f.outside, "missing") }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			_, err := newWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root,
				ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease), ExamplesDir: tc.dir(t, f)})
			var startup *StartupError
			if !errors.As(err, &startup) || startup.Stage != StageConfig || !strings.Contains(err.Error(), "examples corpus") {
				t.Fatalf("New = %v, want a %s StartupError about the examples corpus", err, StageConfig)
			}
			var state *WorkspaceStateError
			var inside *broker.PolicyInsideWorkspaceError
			if tc.state && (!errors.As(err, &state) || !errors.As(err, &inside)) {
				t.Fatalf("New = %v, want it to wrap *WorkspaceStateError and *broker.PolicyInsideWorkspaceError", err)
			}
		})
	}
	f := newWSFixture(t)
	corpus := filepath.Join(f.outside, "examples")
	if err := os.MkdirAll(corpus, 0o755); err != nil {
		t.Fatal(err)
	}
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root,
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease), ExamplesDir: corpus})
	if d.workspace.examplesDir != corpus {
		t.Fatalf("examplesDir = %q, want %q", d.workspace.examplesDir, corpus)
	}
}
