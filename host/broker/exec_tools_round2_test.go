package broker

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func appendExecBytes(t *testing.T, path string) {
	t.Helper()
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("\n# rewritten in place\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil || !os.SameFile(before, after) {
		t.Fatalf("rewrite changed inode: %v", err)
	}
}

func requireExecDrift(t *testing.T, m *movingExec, path string) {
	t.Helper()
	before := m.f.h.Spawns()
	out, err := m.f.call(t, "test", "a")
	var drift *ExecStackDriftError
	if !errors.As(err, &drift) || drift.Path != path || out != nil || m.probes.Load() != 0 || m.f.h.Spawns() != before {
		t.Fatalf("want drift %s without probe/launch: %s, %v, probes %d", path, out, err, m.probes.Load())
	}
}

func TestExecToolInPlaceRewriteAfterFailureRefuses(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.state.nodeFailed(m.state.sandbox, "dependency disappeared")
	appendExecBytes(t, m.a)
	requireExecDrift(t, m, m.a)
}

func TestExecToolRepointWithLibnodeRewrittenInPlaceRefuses(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	lib := filepath.Join(filepath.Dir(m.link), "lib", "libnode.fixture")
	mustMkdir(t, filepath.Dir(lib))
	if err := os.WriteFile(lib, []byte("shared libnode"), 0644); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(filepath.Dir(lib), "libnode.earlier")
	if err := os.WriteFile(other, []byte("another shared library"), 0644); err != nil {
		t.Fatal(err)
	}
	// Accept these shared libraries before changing node. Both nodes'
	// ../lib resolve to them. Digest ordering is not part of file identity.
	m.f.mustRun(t, "test", "a")
	m.state.mu.Lock()
	sb := *m.state.sandbox
	sb.Libnode = append([]execFileDigest(nil), sb.Libnode...)
	sb.Libnode[0], sb.Libnode[1] = sb.Libnode[1], sb.Libnode[0]
	m.state.sandbox = &sb
	m.state.mu.Unlock()
	m.probes.Store(0)
	appendExecBytes(t, lib)
	repointTool(t, m.link, m.b)
	requireExecDrift(t, m, lib)
}

func waitExecSignal(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for probe synchronization")
	}
}

func testExecToolCallerLeavesProbeRunning(t *testing.T) {
	t.Helper()
	m := newMovingExec(t, execFixtureOpts{})
	// Warm the binary before giving the caller a short deadline (macOS's
	// first launch can take 1.7 seconds).
	if err := checkExecNodeVersion(boundedTestContext(t), m.b); err != nil {
		t.Fatal(err)
	}
	entered, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	m.state.cfg.Probe = func(ctx context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		close(entered)
		select {
		case <-release:
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			close(finished)
			return make([]error, len(ps)), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	repointTool(t, m.link, m.b)
	ctx, cancel := context.WithTimeout(boundedTestContext(t), 300*time.Millisecond)
	defer cancel()
	type result struct {
		why string
		err error
	}
	resultCh := make(chan result, 1)
	go func() {
		_, _, why, err := m.state.resolve(ctx, "proj", m.f.worktree)
		resultCh <- result{why, err}
	}()
	waitExecSignal(t, entered)
	r := <-resultCh
	want := "exec node: re-checking the changed target\n  why: the startup probe is still running after " + m.link + " changed\n  next: retry in a few seconds (no restart needed)"
	if r.err != nil || r.why != want {
		t.Fatalf("deadline refusal: %q, %v", r.why, r.err)
	}
	m.state.mu.RLock()
	flight := m.state.flight
	remembered := m.state.tools["node"].why
	m.state.mu.RUnlock()
	if flight == nil || remembered != "" {
		t.Fatalf("deadline remembered or flight stopped: %q", remembered)
	}
	close(release)
	waitExecSignal(t, finished)
	waitExecSignal(t, flight.done)
	running := m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 1 || running.Sandbox.NodeTarget != m.b {
		t.Fatalf("probe was repeated or not accepted: probes %d, %+v", m.probes.Load(), running)
	}
}

func TestExecToolSlowProbeOutlivesTheCaller(t *testing.T) {
	testExecToolCallerLeavesProbeRunning(t)
}

func TestExecToolCallerDeadlineRefusesNotRemembered(t *testing.T) {
	testExecToolCallerLeavesProbeRunning(t)
}

func TestExecToolConcurrentCallersShareOneProbe(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "accepted", true: "failed"}[fail], func(t *testing.T) {
			m := newMovingExec(t, execFixtureOpts{})
			entered, release := make(chan struct{}), make(chan struct{})
			m.state.cfg.Probe = func(ctx context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
				if m.probes.Add(1) != 1 {
					return nil, errors.New("second probe")
				}
				close(entered)
				select {
				case <-release:
					if fail {
						return []error{errors.New("arm7-toolchain: rejected")}, nil
					}
					return make([]error, len(ps)), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			repointTool(t, m.link, m.b)
			const n = 8
			results := make(chan string, n)
			var wg sync.WaitGroup
			ready := make(chan struct{}, n)
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					ready <- struct{}{}
					sb, _, why, err := m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
					if err != nil {
						results <- err.Error()
					} else if why != "" {
						results <- why
					} else {
						results <- sb.Node
					}
				}()
			}
			waitExecSignal(t, entered)
			for i := 0; i < n; i++ {
				<-ready
			}
			close(release)
			wg.Wait()
			close(results)
			var first string
			for got := range results {
				if first == "" {
					first = got
				}
				if got != first || (!fail && got != m.b) || (fail && !strings.Contains(got, "arm7-toolchain: rejected")) {
					t.Fatalf("inconsistent flight result: %q, first %q", got, first)
				}
			}
			if m.probes.Load() != 1 {
				t.Fatalf("probes %d", m.probes.Load())
			}
		})
	}
}

func TestExecToolLateCallCannotClearNodeFailed(t *testing.T) {
	base := canonicalTempDir(t)
	a, b, link := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "tool")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	repointTool(t, link, a)
	m := newMovingExec(t, execFixtureOpts{profile: func(p map[string]any) {
		p["commands"] = map[string]any{"test": map[string]any{"argv": []string{link}, "positional": "relpath"}}
	}})
	old := m.state.sandbox
	m.state.cfg.Probe = func(context.Context, ExecProbeConfig, []*ExecProfile) ([]error, error) {
		// This late result arrives after resolve took its immutable snapshot,
		// before that snapshot's success is published.
		m.state.nodeFailed(old, "late dyld failure")
		return []error{nil}, nil
	}
	repointTool(t, link, b)
	_, _, why, err := m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
	if err != nil || why != "" {
		t.Fatalf("resolve: %s, %v", why, err)
	}
	_, status := m.state.Health()
	if !strings.Contains(status["node"].Probe, "late dyld failure") {
		t.Fatalf("late publication cleared failure: %+v", status["node"])
	}
}

func TestExecToolNodeFailedStaleSnapshotDoesNotPoison(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	old := m.state.sandbox
	repointTool(t, m.link, m.b)
	m.f.mustRun(t, "test", "a")
	m.state.nodeFailed(old, "stale dyld failure")
	_, status := m.state.Health()
	if status["node"].Probe != "ok" || status["node"].Accepted != m.b {
		t.Fatalf("stale result poisoned current node: %+v", status["node"])
	}
}

func TestExecToolPostProbeIdentityRecheck(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.state.cfg.Probe = func(context.Context, ExecProbeConfig, []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		repointTool(t, m.link, m.a)
		return []error{nil}, nil
	}
	repointTool(t, m.link, m.b)
	out, err := m.f.call(t, "test", "a")
	if err != nil || !strings.Contains(string(out), "was interrupted; retry") {
		t.Fatalf("post-probe target change did not refuse: %s, %v", out, err)
	}
	if m.state.Accepted()["node"] != m.a || m.probes.Load() != 1 || m.f.h.Spawns() != 0 {
		t.Fatal("moving target was published or launched")
	}
}

func TestExecToolOneProfilesBrokenToolExcludesOnlyThatProfile(t *testing.T) {
	for _, mode := range []string{"gone", "not-executable", "probe-failed"} {
		t.Run(mode, func(t *testing.T) {
			m := newMovingExec(t, execFixtureOpts{})
			base := canonicalTempDir(t)
			tool := filepath.Join(base, "broken-other-tool")
			if err := os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
				t.Fatal(err)
			}
			second := *m.f.profile
			second.Project = "second"
			second.Commands = make(map[string]ExecProfileCommand)
			c := m.f.profile.Commands["test"]
			c.Argv = []string{tool}
			c.Argv0 = tool
			c.Argv0SHA256, _ = sha256File(tool)
			second.Commands["test"] = c
			second.Probe = c
			state, err := NewExecToolState(ExecToolStateConfig{Node: m.link, Sandbox: m.f.sandbox,
				Profiles:    map[string]*ExecProfile{"proj": m.f.profile, "second": &second},
				ProbeConfig: m.state.cfg.ProbeConfig,
				Probe: func(_ context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
					m.probes.Add(1)
					if mode == "probe-failed" {
						if len(ps) != 2 {
							return nil, errors.New("missing profile")
						}
						return []error{nil, errors.New("arm7-toolchain: broken other tool")}, nil
					}
					if len(ps) != 1 || ps[0].Project != "proj" {
						return nil, errors.New("broken profile included in probe")
					}
					return []error{nil}, nil
				}})
			if err != nil {
				t.Fatal(err)
			}
			m.state, m.f.h.tools = state, state
			if mode == "gone" {
				if err := os.Remove(tool); err != nil {
					t.Fatal(err)
				}
			} else if mode == "not-executable" {
				// A new inode, failing the executable check.
				replacement := filepath.Join(base, "replacement")
				if err := os.WriteFile(replacement, []byte("not executable"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(replacement, tool); err != nil {
					t.Fatal(err)
				}
			}
			repointTool(t, m.link, m.b)
			m.f.mustRun(t, "test", "a")
			_, _, why, err := state.resolve(boundedTestContext(t), "second", m.f.worktree)
			if err != nil || !strings.Contains(why, "exec "+tool+" ") || !strings.Contains(why, "\n  why: ") || !strings.Contains(why, "\n  next: ") {
				t.Fatalf("wrong project's refusal: %q, %v", why, err)
			}
			m.f.mustRun(t, "test", "a")
			if state.Accepted()["node"] != m.b || m.probes.Load() != 1 {
				t.Fatalf("healthy project blocked: probes %d", m.probes.Load())
			}
		})
	}
}

func TestExecToolSymlinkInsideWorkspaceRefuses(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	steering := filepath.Join(m.f.root, "steering")
	repointTool(t, steering, m.b)
	repointTool(t, m.link, steering)
	// The version check itself leaves this marker, so zero probe calls alone
	// cannot hide an unsandboxed --version launch.
	marker := filepath.Join(filepath.Dir(m.link), "version-ran")
	body, err := os.ReadFile(m.b)
	if err != nil {
		t.Fatal(err)
	}
	body = []byte(strings.Replace(string(body), "then echo v26.12.1", "then touch "+shellQuote(marker)+"; echo v26.12.1", 1))
	if err := os.WriteFile(m.b, body, 0755); err != nil {
		t.Fatal(err)
	}
	requireF15(t, m.f, "the workspace root")
	if _, err := os.Stat(marker); !os.IsNotExist(err) || m.probes.Load() != 0 {
		t.Fatalf("forbidden steering link ran --version: %v", err)
	}
}

func TestExecToolStartupChecksEveryLink(t *testing.T) {
	for _, kind := range []string{"configured-node", "intermediate-node", "parent-node", "argv0"} {
		t.Run(kind, func(t *testing.T) {
			base, paths := hostPathsFixture(t)
			node, _ := movingNode(t, filepath.Join(base, "node"))
			marker := filepath.Join(base, "version-ran")
			body, err := os.ReadFile(node)
			if err != nil {
				t.Fatal(err)
			}
			body = []byte(strings.Replace(string(body), "then echo v26.12.1", "then touch "+shellQuote(marker)+"; echo v26.12.1", 1))
			if err := os.WriteFile(node, body, 0755); err != nil {
				t.Fatal(err)
			}
			steering := filepath.Join(paths.WorkspaceRoot, "steering")
			target := node
			if kind == "argv0" {
				target = "/bin/echo"
			}
			if kind == "parent-node" {
				target = filepath.Dir(node)
			}
			repointTool(t, steering, target)
			configured := steering
			if kind != "configured-node" {
				configured = filepath.Join(base, "configured")
				repointTool(t, configured, steering)
			}
			if kind == "parent-node" {
				configured = filepath.Join(configured, filepath.Base(node))
			}
			if kind != "argv0" {
				_, err = resolveExecNode(boundedTestContext(t), configured, []namedRoot{{"the workspace root", paths.WorkspaceRoot}})
			} else {
				p := filepath.Join(base, "profile.json")
				body := "{\"profile\":\"" + ExecProfileVersion + "\",\"project\":\"proj\",\"root\":\".\",\"path\":[\"/bin\",\"/usr/bin\"],\"timeout_ms\":9000,\"commands\":{\"test\":{\"argv\":[\"" + configured + "\"]}}}"
				if err := os.WriteFile(p, []byte(body), 0600); err != nil {
					t.Fatal(err)
				}
				_, err = loadExecProfileFile(p, paths, []namedRoot{{"the workspace root", paths.WorkspaceRoot}})
			}
			if err == nil || !strings.Contains(err.Error(), "the workspace root") {
				t.Fatalf("startup accepted writable chain: %v", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("startup ran --version before placement check: %v", err)
			}
		})
	}
}

func TestExecToolUnknownProjectRefuses(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.state.mu.Lock()
	delete(m.state.profiles, "proj")
	m.state.mu.Unlock()
	out, err := m.f.call(t, "test", "a")
	if err != nil || !strings.Contains(string(out), "unknown exec project") || m.f.h.Spawns() != 0 {
		t.Fatalf("unknown project: %s, %v", out, err)
	}
}
