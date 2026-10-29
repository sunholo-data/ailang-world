package broker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// B's tool is transiently unresolvable (dangling intermediate link) while node
// changes; B is excluded. The link comes back to the SAME accepted target.
// Is B ever probed against the new node before it runs on it?
func TestExecToolExcludedProfileIsProbedWithTheNewNode(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	base := canonicalTempDir(t)
	tool := filepath.Join(base, "real-tool")
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	mid := filepath.Join(base, "mid")
	repointTool(t, mid, tool)
	link := filepath.Join(base, "link")
	repointTool(t, link, mid)
	second := *m.f.profile
	second.Project = "second"
	second.Commands = map[string]ExecProfileCommand{}
	c := m.f.profile.Commands["test"]
	c.Argv, c.Argv0 = []string{link}, link
	c.Argv0SHA256, _ = sha256File(tool)
	second.Commands["test"], second.Probe = c, c
	var mu sync.Mutex
	probed := map[string][]string{}
	state, err := NewExecToolState(ExecToolStateConfig{Node: m.link, Sandbox: m.f.sandbox,
		Profiles:    map[string]*ExecProfile{"proj": m.f.profile, "second": &second},
		ProbeConfig: m.state.cfg.ProbeConfig,
		Probe: func(_ context.Context, cfg ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
			mu.Lock()
			defer mu.Unlock()
			for _, p := range ps {
				probed[cfg.Sandbox.Node] = append(probed[cfg.Sandbox.Node], p.Project)
			}
			return make([]error, len(ps)), nil
		}})
	if err != nil {
		t.Fatal(err)
	}
	m.state, m.f.h.tools = state, state
	os.Remove(mid) // transient: B's chain dangles
	repointTool(t, m.link, m.b)
	m.f.mustRun(t, "test", "a")
	repointTool(t, mid, tool) // back to the same accepted target
	sb, _, why, err := state.resolve(boundedTestContext(t), "second", m.f.worktree)
	secondProbedWithB := false
	for _, p := range probed[m.b] {
		if p == "second" {
			secondProbedWithB = true
		}
	}
	if why == "" && err == nil && sb != nil && sb.Node == m.b && !secondProbedWithB {
		t.Fatalf("profile second accepted on new node %s without ever being probed with it (probed=%v)", m.b, probed)
	}
}

func TestExecToolEveryProfileRecordsItsProbedNode(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	second := *m.f.profile
	second.Project = "second"
	state, err := NewExecToolState(ExecToolStateConfig{Node: m.link, Sandbox: m.f.sandbox,
		Profiles:    map[string]*ExecProfile{"proj": m.f.profile, "second": &second},
		ProbeConfig: m.state.cfg.ProbeConfig,
		Probe: func(_ context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
			m.probes.Add(1)
			if len(ps) != 2 {
				t.Errorf("want both profiles, got %d", len(ps))
			}
			return make([]error, len(ps)), nil
		}})
	if err != nil {
		t.Fatal(err)
	}
	m.state, m.f.h.tools = state, state
	for p := range state.profiles {
		if !sameTool(state.probedNode[p], state.tools["node"].identity) {
			t.Fatalf("startup omitted probed node for %s", p)
		}
	}
	repointTool(t, m.link, m.b)
	m.f.mustRun(t, "test", "a")
	for _, p := range []string{"proj", "second"} {
		state.mu.RLock()
		id := state.probedNode[p]
		node := state.tools["node"].identity
		state.mu.RUnlock()
		if id.path != m.b || !sameTool(id, node) {
			t.Fatalf("%s did not record new node: %s", p, id.path)
		}
		_, _, why, err := state.resolve(boundedTestContext(t), p, m.f.worktree)
		if why != "" || err != nil {
			t.Fatalf("%s: %s, %v", p, why, err)
		}
	}
	if m.probes.Load() != 1 {
		t.Fatalf("later call repeated probe: %d", m.probes.Load())
	}
}

func TestExecToolProbeTimeoutIsNotRemembered(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	m.state.cfg.ProbeBudget = 5 * time.Millisecond
	m.state.cfg.Probe = func(ctx context.Context, _ ExecProbeConfig, _ []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		<-ctx.Done()
		return nil, ctx.Err()
	}
	repointTool(t, m.link, m.b)
	_, _, why, err := m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
	want := "exec node: re-checking the changed target\n  why: the startup probe did not finish within 5ms\n  next: retry; if it keeps timing out, check host load (no restart needed)"
	if err != nil || why != want {
		t.Fatalf("timeout refusal: %q, %v", why, err)
	}
	m.state.mu.RLock()
	for key, tool := range m.state.tools {
		if tool.failed != nil || tool.why != "" {
			t.Errorf("timeout remembered for %s: %+v", key, tool)
		}
	}
	m.state.mu.RUnlock()
	m.state.cfg.Probe = func(_ context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		return make([]error, len(ps)), nil
	}
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 2 {
		t.Fatalf("timeout did not retry: %d", m.probes.Load())
	}
}

func TestExecToolArgv0SteeringLinkInsideWorkspaceRefuses(t *testing.T) {
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
	// The configured path and final executable are outside forbidden roots;
	// only the intermediate link is agent-writable.
	steering := filepath.Join(m.f.root, "steering")
	repointTool(t, steering, b)
	repointTool(t, link, steering)
	marker := filepath.Join(base, "version-ran")
	body, err := os.ReadFile(m.a)
	if err != nil {
		t.Fatal(err)
	}
	body = []byte(strings.Replace(string(body), "then echo v26.12.1", "then touch "+shellQuote(marker)+"; echo v26.12.1", 1))
	if err := os.WriteFile(m.a, body, 0755); err != nil {
		t.Fatal(err)
	}
	// Refresh this fixture's digest so a version-marker rewrite is not drift.
	sum, err := sha256File(m.a)
	if err != nil {
		t.Fatal(err)
	}
	m.state.tools["node"] = func() execAcceptedTool {
		n := m.state.tools["node"]
		n.sum = sum
		return n
	}()
	m.state.sandbox.NodeSHA256 = sum
	out, err := m.f.call(t, "test")
	if err != nil || !strings.Contains(string(out), "argv[0]") || !strings.Contains(string(out), "the workspace root") {
		t.Fatalf("argv[0] steering did not refuse: %s, %v", out, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) || m.probes.Load() != 0 || m.f.h.Spawns() != 0 {
		t.Fatalf("steering caused version/probe/launch: %v", err)
	}
}

func TestExecToolRecheckPanicReleasesWaiters(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	entered, release := make(chan struct{}), make(chan struct{})
	m.state.cfg.Probe = func(context.Context, ExecProbeConfig, []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		close(entered)
		<-release
		panic("probe seam")
	}
	repointTool(t, m.link, m.b)
	result := make(chan error, 1)
	go func() {
		_, _, _, err := m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
		result <- err
	}()
	waitExecSignal(t, entered)
	m.state.mu.RLock()
	flight := m.state.flight
	m.state.mu.RUnlock()
	close(release)
	waitExecSignal(t, flight.done)
	if err := <-result; err == nil || !strings.Contains(err.Error(), "probe seam") {
		t.Fatalf("panic not reported: %v", err)
	}
	m.state.cfg.Probe = func(_ context.Context, _ ExecProbeConfig, ps []*ExecProfile) ([]error, error) {
		m.probes.Add(1)
		return make([]error, len(ps)), nil
	}
	m.f.mustRun(t, "test", "a")
	if m.probes.Load() != 2 {
		t.Fatal("panic remembered")
	}
}

func TestExecToolCloseCancelsFlight(t *testing.T) {
	m := newMovingExec(t, execFixtureOpts{})
	entered, cancelled := make(chan struct{}), make(chan struct{})
	m.state.cfg.Probe = func(ctx context.Context, _ ExecProbeConfig, _ []*ExecProfile) ([]error, error) {
		close(entered)
		<-ctx.Done()
		close(cancelled)
		return nil, ctx.Err()
	}
	repointTool(t, m.link, m.b)
	result := make(chan struct{})
	go func() {
		defer close(result)
		m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
	}()
	waitExecSignal(t, entered)
	m.state.mu.RLock()
	flight := m.state.flight
	m.state.mu.RUnlock()
	m.state.Close()
	m.state.Close()
	waitExecSignal(t, cancelled)
	waitExecSignal(t, flight.done)
	waitExecSignal(t, result)
	if m.state.Accepted()["node"] != m.a {
		t.Fatal("closed flight published changed node")
	}
	_, _, _, err := m.state.resolve(boundedTestContext(t), "proj", m.f.worktree)
	if err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("closed state started another flight: %v", err)
	}
}
