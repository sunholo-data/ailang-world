package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func movingNode(t *testing.T, dir string) (string, string) {
	t.Helper()
	mustMkdir(t, dir)
	node, calls, _, _ := stubNode(t, dir)
	body, err := os.ReadFile(node)
	if err != nil {
		t.Fatal(err)
	}
	body = []byte(strings.Replace(string(body), "#!/bin/sh\n", "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo v26.12.1; exit 0; fi\n", 1))
	if err := os.WriteFile(node, body, 0755); err != nil {
		t.Fatal(err)
	}
	return node, calls
}

func repointTool(t *testing.T, link, target string) {
	t.Helper()
	if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

type movingExec struct {
	f                  *execFixture
	state              *ExecToolState
	link, a, b, callsB string
	probes             atomic.Int32
}

func newMovingExec(t *testing.T, opts execFixtureOpts) *movingExec {
	t.Helper()
	base := canonicalTempDir(t)
	a, _ := movingNode(t, filepath.Join(base, "a"))
	b, callsB := movingNode(t, filepath.Join(base, "b"))
	link := filepath.Join(base, "node")
	repointTool(t, link, a)
	opts.node = link
	f := newExecFixture(t, opts)
	f.calls, f.hostEnv, f.argvLog = filepath.Join(base, "a", "stub.calls"), filepath.Join(base, "a", "stub.hostenv"), filepath.Join(base, "a", "stub.argv")
	m := &movingExec{f: f, link: link, a: a, b: b, callsB: callsB}
	state, err := NewExecToolState(ExecToolStateConfig{Node: link, Sandbox: f.sandbox,
		Profiles:    map[string]*ExecProfile{f.profile.Project: f.profile},
		ProbeConfig: ExecProbeConfig{Paths: ExecHostPaths{WorkspaceRoot: f.root, StateDir: f.state, OperatorHome: f.home}},
		Probe: func(ctx context.Context, cfg ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
			m.probes.Add(1)
			if cfg.Sandbox.Node == link || len(ps) != 1 {
				return nil, errors.New("probe did not receive accepted targets")
			}
			return make([]error, len(ps)), nil
		}})
	if err != nil {
		t.Fatal(err)
	}
	m.state, f.h.tools = state, state
	return m
}

// f15RefusalProblem says why out is not a pre-launch F15 refusal ("" when it is). out must be
// exactly one complete JSON value (trailing whitespace allowed): the two-key envelope. A launch is
// seen by structure, never by text (row 234): a launched reply is execResult, which always
// carries exit_code, so a strict decode into the two-key envelope fails; a post-launch refusal
// (a 134 exit with a broken node, handlers_exec.go:326) has the envelope's shape but moved the
// spawn counter.
func f15RefusalProblem(out []byte, before, after int64) string {
	dec := json.NewDecoder(bytes.NewReader(out))
	dec.DisallowUnknownFields()
	var r struct {
		OK      *bool   `json:"ok"`
		Refused *string `json:"refused"`
	}
	if err := dec.Decode(&r); err != nil {
		return "not the refusal envelope: " + err.Error()
	}
	if err := dec.Decode(&json.RawMessage{}); err != io.EOF {
		return fmt.Sprintf("trailing content after the envelope: %v", err)
	}
	if r.OK == nil || *r.OK {
		return "ok is not false"
	}
	if r.Refused == nil || *r.Refused == "" {
		return "refused is empty"
	}
	if after != before {
		return "the call spawned"
	}
	return ""
}

func requireF15(t *testing.T, f *execFixture, parts ...string) string {
	t.Helper()
	before := f.h.Spawns()
	out, err := f.call(t, "test", "a")
	if err != nil {
		t.Fatal(err)
	}
	after := f.h.Spawns()
	if problem := f15RefusalProblem(out, before, after); problem != "" {
		t.Fatalf("not a pre-launch refusal (%s): %s (spawns %d -> %d)", problem, out, before, after)
	}
	var refusal struct {
		Refused string `json:"refused"`
	}
	if err := json.Unmarshal(out, &refusal); err != nil {
		t.Fatal(err)
	}
	for _, part := range append([]string{"does not start:", "\n  why: ", "\n  next: "}, parts...) {
		if !strings.Contains(refusal.Refused, part) {
			t.Fatalf("missing %q: %s", part, out)
		}
	}
	return refusal.Refused
}

func TestExecToolFollowsARepointedNode(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.f.mustRun(t, "test", "a")
	repointTool(t, m.link, m.b)
	r := m.f.mustRun(t, "test", "a")
	sum, _ := sha256File(m.b)
	if r.Sandbox.Node != m.link || r.Sandbox.NodeTarget != m.b || r.Sandbox.NodeSHA256 != sum || m.probes.Load() != 1 {
		t.Fatalf("after repoint: %+v; probes %d", r.Sandbox, m.probes.Load())
	}
	body, _ := os.ReadFile(m.callsB)
	if strings.Count(string(body), "call\n") != 1 {
		t.Fatalf("B did not run once: %q", body)
	}
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 1 {
		t.Fatal("accepted identity was probed again")
	}
}

func TestExecToolDeletedTargetRefusesF15(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	if err := os.Remove(m.a); err != nil {
		t.Fatal(err)
	}
	requireF15(t, m.f, m.link, m.a, "the target is gone", "no restart needed")
	if m.probes.Load() != 0 || m.f.stubCalls(t) != 0 {
		t.Fatal("deleted target spawned something")
	}
	repointTool(t, m.link, m.b)
	r := m.f.mustRun(t, "test", "a")
	if r.Sandbox.NodeTarget != m.b || m.probes.Load() != 1 {
		t.Fatalf("did not recover: %+v", r)
	}
}

func TestExecToolFailedProbeRefusesF15(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.state.cfg.Probe = func(context.Context, ExecProbeConfig, []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		return []error{&ExecProbeError{Project: "proj", Failed: []ExecProbeArmFailure{{Arm: ExecProbeArmToolchain, Why: "fixture failed"}}}}, nil
	}
	repointTool(t, m.link, m.b)
	first := requireF15(t, m.f, ExecProbeArmToolchain, "the startup probe failed after the target changed", "no restart needed")
	second := requireF15(t, m.f, ExecProbeArmToolchain)
	if first != second || m.probes.Load() != 1 {
		t.Fatalf("failure not remembered: probes %d; %s vs %s", m.probes.Load(), first, second)
	}
	c, _ := movingNode(t, filepath.Join(filepath.Dir(m.link), "c"))
	repointTool(t, m.link, c)
	requireF15(t, m.f, ExecProbeArmToolchain)
	if m.probes.Load() != 2 {
		t.Fatal("new identity was not probed")
	}
}

func TestExecToolArgv0FollowsARepointedLink(t *testing.T) {
	base := canonicalTempDir(t)
	a, b, link := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "tool")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("#!/bin/sh\nprintf '"+filepath.Base(p)+"'\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	repointTool(t, link, a)
	m := newMovingExec(t, execFixtureOpts{profile: func(p map[string]any) {
		p["commands"] = map[string]any{"test": map[string]any{"argv": []any{link}, "positional": "relpath", "max_args": 8}}
	}})
	repointTool(t, link, b)
	r := m.f.mustRun(t, "test", "a")
	sum, _ := sha256File(b)
	log, _ := os.ReadFile(m.f.argvLog)
	if r.Stdout != "b" || r.Profile.Argv0 != link || r.Profile.Argv0Target != b || r.Profile.Argv0SHA256 != sum ||
		!strings.Contains(string(log), b+"\x00") || strings.Contains(string(log), link+"\x00") || m.probes.Load() != 1 {
		t.Fatalf("wrong command target: %+v, launch %q", r, log)
	}
}

func TestExecToolPlacementReapplied(t *testing.T) {
	for _, tool := range []string{"node", "argv0"} {
		for _, root := range []string{"workspace", "state"} {
			t.Run(tool+"/"+root, func(t *testing.T) {
				base := canonicalTempDir(t)
				bin := filepath.Join(base, "tool")
				if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
					t.Fatal(err)
				}
				link := filepath.Join(base, "link")
				repointTool(t, link, bin)
				m := newMovingExec(t, execFixtureOpts{profile: func(p map[string]any) {
					p["commands"] = map[string]any{"test": map[string]any{"argv": []any{link}, "positional": "relpath", "max_args": 8}}
				}})
				dir, want := m.f.root, "the workspace root"
				if root == "state" {
					dir, want = m.f.state, "the state dir"
				}
				target, _ := movingNode(t, filepath.Join(dir, "forbidden"))
				if tool == "node" {
					link = m.link
				}
				repointTool(t, link, target)
				requireF15(t, m.f, want)
				if m.probes.Load() != 0 || m.f.stubCalls(t) != 0 {
					t.Fatal("placement refusal spawned or probed")
				}
			})
		}
	}
	t.Run("startup", func(t *testing.T) {
		base, paths := hostPathsFixture(t)
		node, _ := movingNode(t, filepath.Join(base, "node"))
		tool := filepath.Join(paths.WorkspaceRoot, "tool")
		if err := os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(base, "profile.json")
		body, _ := json.Marshal(map[string]any{"profile": ExecProfileVersion, "project": "proj", "root": ".",
			"path": []string{"/usr/bin", "/bin"}, "timeout_ms": 9000,
			"commands": map[string]any{"test": map[string]any{"argv": []string{tool}}}})
		if err := os.WriteFile(p, body, 0600); err != nil {
			t.Fatal(err)
		}
		tree := fakeSrtTree(t, "0.0.0-test", "// fixture")
		_, err := StartExec(boundedTestContext(t), ExecStartupConfig{ProfileFiles: []string{p}, Node: node,
			SandboxDir: tree, Pin: fakePin(t, tree), Paths: paths})
		if err == nil || !strings.Contains(err.Error(), "the workspace root") {
			t.Fatalf("startup accepted forbidden argv0: %v", err)
		}
	})
}

func TestExecToolExit134ChecksNode(t *testing.T) {
	for _, fail := range []bool{true, false} {
		t.Run(map[bool]string{true: "broken-node", false: "command-aborted"}[fail], func(t *testing.T) {
			dir := canonicalTempDir(t)
			node := filepath.Join(dir, "node")
			version := "echo v26.12.1; exit 0"
			if fail {
				version = "echo 'dyld: Library not loaded: libada.3.dylib' >&2; exit 1"
			}
			body := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then " + version + "; fi\nexit 134\n"
			if err := os.WriteFile(node, []byte(body), 0755); err != nil {
				t.Fatal(err)
			}
			f := newExecFixture(t, execFixtureOpts{node: node})
			before := f.h.Spawns()
			out, err := f.call(t, "test", "a")
			if err != nil {
				t.Fatal(err)
			}
			if f15RefusalProblem(out, before, f.h.Spawns()) == "" {
				t.Fatalf("a launch that exited 134 passes requireF15: %s", out)
			}
			if fail {
				if !strings.Contains(string(out), "dyld: Library not loaded") || !strings.Contains(string(out), "no restart needed") || strings.Contains(string(out), "exit status 134") {
					t.Fatalf("not F15: %s", out)
				}
			} else if r := execDecode(t, out); r.ExitCode == nil || *r.ExitCode != 134 {
				t.Fatalf("command abort lost: %s", out)
			}
		})
	}
}

func TestExecToolOneReacceptanceUnderConcurrency(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	entered, release := make(chan struct{}), make(chan struct{})
	m.state.cfg.Probe = func(ctx context.Context, cfg ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		close(entered)
		select {
		case <-release:
			return make([]error, len(ps)), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	repointTool(t, m.link, m.b)
	var wg sync.WaitGroup
	results := make(chan []byte, 8)
	errs := make(chan error, 8)
	ctx := boundedTestContext(t)
	for i := 0; i < 8; i++ {
		// Separate episode handlers share state and immutable snapshots.
		h, err := NewExecHandler(ExecHandlerConfig{Tools: m.state, Profile: m.f.profile, Sandbox: m.f.sandbox,
			Episode: "ep" + string(rune('a'+i)), Worktree: m.f.worktree, WorkspaceRoot: m.f.root,
			StateDir: m.f.state, OperatorHome: m.f.home})
		if err != nil {
			t.Fatal(err)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := h.Execute(ctx, EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope},
				[]byte("{\"command\":\"test\",\"args\":[\"a\"]}"))
			results <- out
			errs <- err
		}()
	}
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("probe did not start")
	}
	close(release)
	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	for out := range results {
		r := execDecode(t, out)
		if r.Sandbox.NodeTarget != m.b || r.ExitCode == nil || *r.ExitCode != 0 {
			t.Fatalf("wrong target: %s", out)
		}
	}
	if m.probes.Load() != 1 {
		t.Fatalf("probes %d, want one", m.probes.Load())
	}
	body, _ := os.ReadFile(m.callsB)
	if strings.Count(string(body), "call\n") != 8 {
		t.Fatalf("not all calls ran B: %q", body)
	}
}

func TestExecToolCellarProfilePathRefusesByName(t *testing.T) {
	base := canonicalTempDir(t)
	tool := filepath.Join(base, "Cellar", "go", "1.26", "bin", "go")
	mustMkdir(t, filepath.Dir(tool))
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	m := newMovingExec(t, execFixtureOpts{profile: func(p map[string]any) {
		p["commands"] = map[string]any{"test": map[string]any{"argv": []any{tool}, "positional": "relpath", "max_args": 8}}
	}})
	if err := os.Remove(tool); err != nil {
		t.Fatal(err)
	}
	requireF15(t, m.f, tool, "\n  next: ailang-worldd profile init --repo ", "--project proj")
}

func TestExecToolInPlaceRewriteStaysDrift(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	flipByte(t, m.a)
	out, err := m.f.call(t, "test", "a")
	var drift *ExecStackDriftError
	if !errors.As(err, &drift) || drift.Path != m.a || out != nil || m.probes.Load() != 0 {
		t.Fatalf("in-place rewrite reaccepted: %s, %v", out, err)
	}
}

// Caller cancellation no longer cancels the shared probe (M3 round 2).
func TestExecToolInterruptedProbeIsNotRemembered(t *testing.T) {
	testExecToolCallerLeavesProbeRunning(t)
}

func TestExecToolLibnodeIdentityAndDrift(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	lib := filepath.Join(filepath.Dir(m.link), "lib", "libnode.fixture")
	mustMkdir(t, filepath.Dir(lib))
	if err := os.WriteFile(lib, []byte("libnode A"), 0644); err != nil {
		t.Fatal(err)
	}
	m.f.mustRun(t, "test", "a")
	replacement := lib + ".new"
	if err := os.WriteFile(replacement, []byte("libnode B"), 0644); err != nil {
		t.Fatal(err)
	}
	// Keep the temporary file outside the libnode.* glob.
	other := filepath.Join(filepath.Dir(m.link), "replacement")
	if err := os.Rename(replacement, other); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(other, lib); err != nil {
		t.Fatal(err)
	}
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 2 || len(m.f.sandbox.Libnode) != 0 {
		t.Fatalf("libnode not reaccepted independently, or startup snapshot mutated: %d", m.probes.Load())
	}
	flipByte(t, lib)
	out, err := m.f.call(t, "test", "a")
	var drift *ExecStackDriftError
	if !errors.As(err, &drift) || drift.Path != lib || out != nil || m.probes.Load() != 2 {
		t.Fatalf("libnode tampering reaccepted: %s, %v", out, err)
	}
}

func TestExecToolStartFailureIsRetriedWithoutChangingIdentity(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	marker := filepath.Join(filepath.Dir(m.link), "broken")
	body, err := os.ReadFile(m.b)
	if err != nil {
		t.Fatal(err)
	}
	guard := "if [ -f " + shellQuote(marker) + " ]; then echo 'dyld: broken dependency' >&2; exit 1; fi; "
	body = []byte(strings.Replace(string(body), "then echo v26.12.1", "then "+guard+"echo v26.12.1", 1))
	if err := os.WriteFile(m.b, body, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(marker, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	repointTool(t, m.link, m.b)
	requireF15(t, m.f, "dyld: broken dependency", "--version failed")
	requireF15(t, m.f, "dyld: broken dependency")
	if m.probes.Load() != 0 {
		t.Fatal("a node that does not start was probed")
	}
	_, health := m.state.Health()
	if !strings.HasPrefix(health["node"].Probe, "does-not-start:") {
		t.Fatalf("start failure not visible in health: %+v", health)
	}
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 1 {
		t.Fatal("fixed dependency did not recover without a new node identity")
	}
}

func TestExecToolNodeChangeProbesEveryProfile(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	second := *m.f.profile
	second.Project = "second"
	state, err := NewExecToolState(ExecToolStateConfig{Node: m.link, Sandbox: m.f.sandbox,
		Profiles:    map[string]*ExecProfile{"proj": m.f.profile, "second": &second},
		ProbeConfig: m.state.cfg.ProbeConfig,
		Probe: func(_ context.Context, cfg ExecProbeConfig, profiles []*ExecProfile) ([]error, error) {
			m.probes.Add(1)
			if len(profiles) != 2 || cfg.Sandbox.Node != m.b {
				return nil, errors.New("node change did not probe both profiles with B")
			}
			return []error{nil, nil}, nil
		}})
	if err != nil {
		t.Fatal(err)
	}
	m.f.h.tools = state
	repointTool(t, m.link, m.b)
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 1 {
		t.Fatal("node did not trigger a shared probe")
	}
}

func TestExecToolArchiveDriftPrecedesReacceptance(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	flipByte(t, m.f.sandbox.CLI())
	repointTool(t, m.link, m.b)
	out, err := m.f.call(t, "test", "a")
	var drift *ExecStackDriftError
	if !errors.As(err, &drift) || drift.Path != m.f.sandbox.CLI() || out != nil || m.probes.Load() != 0 || m.f.h.Spawns() != 0 {
		t.Fatalf("archive drift reached reacceptance: %s, %v", out, err)
	}
}

// T2 (row 234): the six requireF15 callers' refusals quote t.TempDir paths, so a temp root
// holding "134" must not red a real pre-launch refusal. The root is forced through GOTMPDIR
// (testing reads it first, one parent per test); the precondition makes the hazard real.
func TestExecToolF15RefusalUnderTempRoot134(t *testing.T) {
	root, err := os.MkdirTemp("", "w234-x134x-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	t.Setenv("GOTMPDIR", root)
	m := newMovingExec(t, execFixtureOpts{})
	if !strings.Contains(m.link, "x134x") {
		t.Fatalf("temp root not honoured: %s is not under %s", m.link, root)
	}
	if err := os.Remove(m.a); err != nil {
		t.Fatal(err)
	}
	requireF15(t, m.f, m.link, m.a, "the target is gone")
}

// T1 (row 234): f15RefusalProblem judges structure and the spawn counter, never text. The
// accepted rows quote temp paths holding "134"; every reject row fails on one property.
func TestF15RefusalProblemIsStructural(t *testing.T) {
	env := `{"ok":false,"refused":"exec /tmp/x134x/TestExecToolCellarProfilePathRefusesByName1349870625/001/node does not start: the target is gone"}`
	for _, row := range []struct {
		name          string
		out           string
		before, after int64
		accepted      bool
	}{
		{"refusal-temp-root-134", env, 0, 0, true},
		{"trailing-whitespace", env + " \n\t", 0, 0, true},
		{"post-launch-f15", env, 0, 1, false},
		{"launched-exit-134", `{"exit_code":134,"timed_out":false,"limit":null,"duration_ms":3,"stdout":"","stderr":"/tmp/x134x/node: Abort trap: 6"}`, 0, 1, false},
		{"refusal-plus-exit-code", `{"ok":false,"refused":"exec /tmp/x134x/node does not start","exit_code":134}`, 0, 0, false},
		{"ok-true", `{"ok":true,"refused":"exec /tmp/x134x/node does not start"}`, 0, 0, false},
		{"refused-empty", `{"ok":false,"refused":""}`, 0, 0, false},
		{"concatenated-json", env + `{"exit_code":134}`, 0, 0, false},
		{"trailing-garbage", env + " garbage", 0, 0, false},
	} {
		t.Run(row.name, func(t *testing.T) {
			problem := f15RefusalProblem([]byte(row.out), row.before, row.after)
			if (problem == "") != row.accepted {
				t.Fatalf("problem %q, want accepted=%v: %s", problem, row.accepted, row.out)
			}
		})
	}
}
