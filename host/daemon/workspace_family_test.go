package daemon

// Row 152: the two handler families an episode registry binds — the eight
// AILANG tool names (need a built AILANG handler: module root, package cache,
// policy, summary subprocess) and Workspace.Exec (needs only the worktree) —
// fail independently, and a stuck AILANG build never blocks exec construction
// or another episode.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// gatedToolBin is fakeToolBin plus per-episode control of the policy-tool
// summary. A summary request appends "x" to <log>/summaries and
// <log>/summaries.<ep> (ep = basename of the cwd, the episode worktree: these
// tests use no module root), waits while <log>/hold.<ep> exists, then answers
// {"ok":false} (exit 1) if <log>/fail.<ep> exists, else the healthy summary.
// Releasing a hold is removing its file; a hold never released plus a low
// budget is a summary that never answers.
func gatedToolBin(t *testing.T, logDir string) string {
	t.Helper()
	script := fmt.Sprintf(`#!/bin/sh
log=%q
if [ "$1" = "--version" ]; then echo %q; exit 0; fi
if [ "$1" = "policy-tool" ]; then
  req=$(cat)
  case "$req" in
    *'"summary"'*)
      ep=$(basename "$(pwd -P)")
      echo x >> "$log/summaries"
      echo x >> "$log/summaries.$ep"
      while [ -e "$log/hold.$ep" ]; do sleep 0.02; done
      if [ -e "$log/fail.$ep" ]; then printf '{"ok":false}\n'; exit 1; fi
      printf '{"ok":true,"summary":{"security_mode":"restricted","policy_digest":"fakedigest","fs_sandbox":"%%s","cli":["check"]}}' "$(pwd -P)"
      exit 0 ;;
  esac
  echo x >> "$log/dispatches"
  printf '{"ok":true,"content":"fake:%%s"}' "$(pwd -P)"
  exit 0
fi
echo x >> "$log/dispatches"
exit 0
`, logDir, ToolBinaryRelease)
	bin := filepath.Join(canonicalTemp(t), "ailang")
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return bin
}

// withHandlerBudget sets the workspaceHandlerBudget seam for the test. Tests
// using it are serial and join every build goroutine before it is restored.
func withHandlerBudget(t *testing.T, d time.Duration) {
	t.Helper()
	old := workspaceHandlerBudget
	workspaceHandlerBudget = d
	t.Cleanup(func() { workspaceHandlerBudget = old })
}

// execPlanRunner answers the plan phase with one Workspace.Exec `where`.
type execPlanRunner struct{ plans atomic.Int64 }

const wsExecPlan = `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Exec","scope":"worktree","cost":1,"payload":{"command":"where"}}],"finish":false,"result":null}`

func (r *execPlanRunner) RunContext(_ context.Context, e capsule.Entry) (capsule.Result, error) {
	var arg string
	if err := json.Unmarshal(e.Args, &arg); err != nil {
		return capsule.Result{}, err
	}
	if !strings.Contains(arg, `"phase":"plan"`) {
		return capsule.Result{}, fmt.Errorf("unexpected phase input %s", arg)
	}
	r.plans.Add(1)
	return capsule.Result{Stdout: []byte(wsExecPlan + "\n")}, nil
}

// execDispatchCall is readCall's shape for ws.exec (name: execCall exists).
func execDispatchCall(t *testing.T, d *Daemon, episode, task string, now int64) coordinator.Call {
	t.Helper()
	grants := []broker.Capability{{Effect: broker.EffectWorkspaceExec, Scope: broker.WorkspaceScope, ExpiresAt: now + 3600, Budget: 5}}
	req, err := transitionreg.NewRequest(boundedTestContext(t), transitionreg.NewReader(d.store), wsCaps(grants), now)
	if err != nil {
		t.Fatal(err)
	}
	return coordinator.Call{Surface: coordinator.SurfaceA2A, Request: req, EpisodeID: episode, Grants: grants, SkillID: "ws.exec", TaskID: task, Input: map[string]any{}}
}

// familyRig is a daemon with the gated tool binary, exec configured, ws.exec
// and ws.read published, and one coordinator (own fake plan runner) per tool
// over the one store and binder.
type familyRig struct {
	t                    *testing.T
	f                    wsFixture
	d                    *Daemon
	log                  *syncLog
	execRunner           *execPlanRunner
	readRunner           *readPlanRunner
	execCoord, readCoord *coordinator.Coordinator
	now                  int64
}

func newFamilyRig(t *testing.T, cfg Config) *familyRig {
	t.Helper()
	f := newWSFixture(t)
	r := &familyRig{t: t, f: f, log: &syncLog{}, execRunner: &execPlanRunner{}, readRunner: &readPlanRunner{}, now: time.Now().Unix()}
	cfg.DBPath, cfg.WorkspaceRoot, cfg.ErrorLog = f.db, f.root, r.log
	cfg.ToolAilangBin = gatedToolBin(t, f.logDir)
	r.d = mustWSDaemon(t, cfg)
	configureStubExec(t, r.d, f)
	publishTools(t, r.d, hashref.SumSHA256([]byte("fake-interpreter")),
		"module host/capsule/main\n\nexport func main(input: string) -> string {\n  input\n}\n",
		wsTool{ID: "ws.exec", Effect: broker.EffectWorkspaceExec, Title: "Exec", Description: "workspace exec fixture"},
		wsTool{ID: "ws.read", Effect: broker.EffectWorkspaceRead, Title: "Read", Description: "workspace read fixture"})
	mk := func(run coordinator.Runner) *coordinator.Coordinator {
		c, err := coordinator.New(coordinator.Config{Store: r.d.store, Runner: run, Binder: r.d.binder,
			Now: func() int64 { return r.now }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	r.execCoord, r.readCoord = mk(r.execRunner), mk(r.readRunner)
	return r
}

func (r *familyRig) mark(kind, ep string) string { return filepath.Join(r.f.logDir, kind+"."+ep) }

// set creates the hold/fail marker for ep.
func (r *familyRig) set(kind, ep string) {
	r.t.Helper()
	writeFile(r.t, r.mark(kind, ep), "x\n")
}

// clear removes the marker (releasing a hold).
func (r *familyRig) clear(kind, ep string) {
	r.t.Helper()
	if err := os.Remove(r.mark(kind, ep)); err != nil && !os.IsNotExist(err) {
		r.t.Fatal(err)
	}
}

func (r *familyRig) summaries(ep string) int {
	return countLines(r.t, filepath.Join(r.f.logDir, "summaries."+ep))
}

// requireExecCommit dispatches ws.exec in ep and requires a committed result
// with exit_code 0 and stdout == the ep worktree.
func (r *familyRig) requireExecCommit(ep, task string) {
	r.t.Helper()
	res, err := r.execCoord.Dispatch(boundedTestContext(r.t), execDispatchCall(r.t, r.d, ep, task, r.now))
	if err != nil {
		r.t.Fatalf("ws.exec Dispatch in %s: %v", ep, err)
	}
	var out struct {
		ExitCode *int   `json:"exit_code"`
		Stdout   string `json:"stdout"`
	}
	if err := json.Unmarshal(res.OutputBytes, &out); err != nil || out.ExitCode == nil || *out.ExitCode != 0 ||
		strings.TrimSpace(out.Stdout) != filepath.Join(r.f.root, ep) {
		r.t.Fatalf("ws.exec output in %s = %s (%v), want exit_code 0 and stdout %s", ep, res.OutputBytes, err, filepath.Join(r.f.root, ep))
	}
}

func (r *familyRig) requireReadR8(ep, task string) {
	r.t.Helper()
	_, err := r.readCoord.Dispatch(boundedTestContext(r.t), readCall(r.t, r.d, ep, task, r.now))
	var unsupported *coordinator.EffectsUnsupportedError
	if !errors.As(err, &unsupported) || fmt.Sprint(unsupported.Effects) != "[Workspace.Read]" {
		r.t.Fatalf("ws.read Dispatch in %s = %v, want R8 naming [Workspace.Read]", ep, err)
	}
	if n := r.readRunner.plans.Load(); n != 0 {
		r.t.Fatalf("a refused ws.read ran %d plans, want 0", n)
	}
}

// AC1.4: the row's pin end to end. A policy summary that never answers (hold
// never released, 300 ms budget) leaves Workspace.Exec bound and committing,
// while the AILANG tool is refused (R8) before its plan.
func TestWorkspaceExecSurvivesPolicySummaryTimeout(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		withHandlerBudget(t, 300*time.Millisecond)
		r := newFamilyRig(t, Config{})
		r.set("hold", "ep1")
		// (i) exec commits.
		r.requireExecCommit("ep1", "t-exec")
		// (ii) the AILANG tool is R8, with zero plans.
		r.requireReadR8("ep1", "t-read")
		// (iii) the operator line names the AILANG cause.
		text := r.log.String()
		if !strings.Contains(text, `workspace tools unavailable for episode "ep1": broker: policy summary:`) || !strings.Contains(text, "timed out") {
			t.Fatalf("operator log = %q, want the policy-summary timeout line", text)
		}
		// (iv) the cause is the typed handler timeout.
		_, err := r.d.workspace.episodeHandler("ep1", filepath.Join(r.f.root, "ep1"))
		var hte *broker.HandlerTimeoutError
		if !errors.As(err, &hte) {
			t.Fatalf("episodeHandler error = %v, want a *broker.HandlerTimeoutError", err)
		}
	})
	t.Run("healthy", func(t *testing.T) {
		// Control: the same fixture without a hold commits both calls.
		withHandlerBudget(t, 5*time.Second)
		r := newFamilyRig(t, Config{})
		r.requireExecCommit("ep1", "t-exec")
		res, err := r.readCoord.Dispatch(boundedTestContext(t), readCall(t, r.d, "ep1", "t-read", r.now))
		if err != nil {
			t.Fatalf("healthy ws.read Dispatch: %v", err)
		}
		if want := "fake:" + filepath.Join(r.f.root, "ep1"); !strings.Contains(string(res.OutputBytes), want) {
			t.Fatalf("healthy read output = %s, want it to contain %s", res.OutputBytes, want)
		}
		if r.log.String() != "" {
			t.Fatalf("healthy operator log = %q, want empty", r.log.String())
		}
	})
}

// AC1.9: a fault before the policy write (a regular file where the policies
// directory belongs: MkdirAll fails ENOTDIR for every uid, root included)
// leaves exactly [Workspace.Exec], and exec still commits.
func TestWorkspaceExecSurvivesPolicyDirFault(t *testing.T) {
	r := newFamilyRig(t, Config{})
	policies := filepath.Join(r.f.stateDir, "policies")
	if err := os.RemoveAll(policies); err != nil {
		t.Fatal(err)
	}
	writeFile(t, policies, "plant\n")
	if info, err := os.Lstat(policies); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("the planted fault is not a regular file (%v, %v)", info, err)
	}
	reg := r.d.workspace.registry("ep1")
	requireOnlyExec(t, reg)
	if n := r.summaries("ep1"); n != 0 {
		t.Fatalf("a policy-directory fault still ran %d summaries", n)
	}
	r.requireExecCommit("ep1", "t-exec")
	if !strings.Contains(r.log.String(), "not a directory") {
		t.Fatalf("operator log = %q, want the mkdir fault", r.log.String())
	}
}

// pendingRegistry is a registry(ep) call running in a goroutine. done closes
// when it returns; reg is valid only after that. Callers register a cleanup
// that joins it after the budget seam, so it runs before the seam is restored.
type pendingRegistry struct {
	done chan struct{}
	reg  broker.Registry
}

func (r *familyRig) startRegistry(ep string) *pendingRegistry {
	p := &pendingRegistry{done: make(chan struct{})}
	go func() {
		defer close(p.done)
		p.reg = r.d.workspace.registry(ep)
	}()
	return p
}

// AC1.7 (exec half): exec construction for another episode returns while
// ep1's AILANG build is still held.
func TestWorkspaceExecNeverWaitsOnAilangBuild(t *testing.T) {
	withHandlerBudget(t, 5*time.Second)
	r := newFamilyRig(t, Config{})
	r.set("hold", "ep1")
	ep1 := r.startRegistry("ep1")
	t.Cleanup(func() { // runs before the budget restore: join the build
		r.clear("hold", "ep1")
		<-ep1.done
	})
	waitFor(t, "ep1's summary to start", 30*time.Second, func() bool { return r.summaries("ep1") == 1 })

	h := r.d.workspace.execHandler("ep2", filepath.Join(r.f.root, "ep2"))
	out, err := execCall(t, h)
	if err != nil || !strings.Contains(string(out), filepath.Join(r.f.root, "ep2")) {
		t.Fatalf("exec in ep2 = %s, %v", out, err)
	}
	select {
	case <-ep1.done:
		t.Fatal("ep1's AILANG build had already finished when ep2's exec returned: the hold did not hold, or exec waited on the build")
	default:
	}
	r.clear("hold", "ep1")
	<-ep1.done
	if got := registryNames(ep1.reg); strings.Join(got, ",") != strings.Join(wantWorkspaceEffects, ",") {
		t.Fatalf("ep1 registry after release = %v, want all nine", got)
	}
}

// AC1.7 r3 / AC2.3: two episodes, one failing and one healthy, hit
// registry(ep) and execHandler(ep) from 8 goroutines each. Asserted after the
// join, never on timing; run it under -race.
func TestWorkspaceTwoEpisodeConcurrentRegistryExecRace(t *testing.T) {
	r := newFamilyRig(t, Config{})
	r.set("fail", "ep1")
	const n = 8
	episodes := []string{"ep1", "ep2"}
	regs := map[string][]broker.Registry{}
	execs := map[string][]broker.Handler{}
	for _, ep := range episodes {
		regs[ep], execs[ep] = make([]broker.Registry, n), make([]broker.Handler, n)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, ep := range episodes {
		for i := 0; i < n; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				regs[ep][i] = r.d.workspace.registry(ep)
			}()
			go func() {
				defer wg.Done()
				<-start
				execs[ep][i] = r.d.workspace.execHandler(ep, filepath.Join(r.f.root, ep))
			}()
		}
	}
	close(start)
	wg.Wait()

	for _, ep := range episodes {
		r.d.workspace.execMu.Lock()
		cached := r.d.workspace.execH[ep].h
		r.d.workspace.execMu.Unlock()
		if cached == nil {
			t.Fatalf("%s: no exec handler cached", ep)
		}
		for i := 0; i < n; i++ {
			names := strings.Join(registryNames(regs[ep][i]), ",")
			switch ep {
			case "ep1":
				if names != broker.EffectWorkspaceExec {
					t.Fatalf("ep1 registry %d = %s, want only Workspace.Exec", i, names)
				}
			default:
				if names != strings.Join(wantWorkspaceEffects, ",") {
					t.Fatalf("ep2 registry %d = %s, want all nine", i, names)
				}
			}
			if regs[ep][i][broker.EffectWorkspaceExec] != cached {
				t.Fatalf("%s registry %d binds a different exec handler than the cache", ep, i)
			}
			if execs[ep][i] != cached {
				t.Fatalf("%s execHandler call %d returned a different handler than the cache", ep, i)
			}
		}
	}
	if n := r.summaries("ep2"); n != 1 {
		t.Fatalf("healthy ep2 ran %d summaries, want 1", n)
	}
}
