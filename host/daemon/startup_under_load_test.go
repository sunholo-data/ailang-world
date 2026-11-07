package daemon

// Row 256 (design_docs/planned/w-daemon-startup-under-load.md):
// the exec startup probe and registry shared one 9 s startup clock, every
// directory profile was left out with "the startup budget is spent", and the
// registry bootstrap, which shared that clock, then failed with "context
// deadline exceeded", so the service never came up. These tests stand in a
// probe client that never answers for the load: it makes the probe spend the
// whole startup clock. This does not establish the incident's cause.

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// lockedLog is an operator log the tests read back.
type lockedLog struct {
	mu sync.Mutex
	b  strings.Builder
}

func (l *lockedLog) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedLog) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

// loadRig is an exec rig whose one profile comes from --exec-profiles-dir (as
// on the home service) and episode ep1 is mapped to it. slow makes the probe
// client a process that never answers, so the probe spends the startup clock.
func loadRig(t *testing.T, slow bool) (Config, *lockedLog) {
	t.Helper()
	r := newExecRig(t)
	dir := filepath.Join(r.f.outside, "profiles-dir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(r.profile)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "proj.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	r.cfg.ExecProfilesDir = dir
	r.cfg.ExecEpisodeProjects = []string{"ep1=proj"}
	client := ""
	if slow {
		client = filepath.Join(r.f.outside, "slow-client")
		if err := os.WriteFile(client, []byte("#!/bin/sh\nexec sleep 60\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	prev := execStartupHook
	execStartupHook = func(c *broker.ExecStartupConfig) {
		c.Pin = *r.pin
		c.ProbeClient = client
	}
	t.Cleanup(func() { execStartupHook = prev })
	log := &lockedLog{}
	cfg := r.cfg
	cfg.BindHost, cfg.ErrorLog = DefaultBindHost, log
	return cfg, log
}

// withStartupClocks sets serve's two startup clocks for one test.
func withStartupClocks(t *testing.T, startup, registry time.Duration) {
	t.Helper()
	prevS, prevR := startupBudget, registryPhaseBudget
	startupBudget, registryPhaseBudget = startup, registry
	t.Cleanup(func() { startupBudget, registryPhaseBudget = prevS, prevR })
}

// TestServeStartsWhenTheProbeSpendsTheStartupClock drives Run, serve's own
// lifecycle, with a probe that spends the whole startup clock. Run must
// announce and answer /v1/health, naming the excluded profile and late
// failure observations without attributing a cause. Mutants: Run back on HEAD's one shared clock
// (startBounded(ctx, cfg, startupBudget, New)), or startDaemon starting the
// registry clock from the startup clock, and Run returns
// "daemon startup failed at registry-bootstrap: … context deadline exceeded".
func TestServeStartsWhenTheProbeSpendsTheStartupClock(t *testing.T) {
	cfg, log := loadRig(t, true)
	withStartupClocks(t, 1500*time.Millisecond, 9*time.Second)

	pr, pw := io.Pipe()
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	defer cancel()
	ran := make(chan error, 1)
	begin := time.Now()
	go func() {
		err := Run(ctx, cfg, pw)
		_ = pw.Close()
		ran <- err
	}()
	lines := make(chan string, 8)
	go func() {
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()
	var first string
	select {
	case l, ok := <-lines:
		if !ok {
			t.Fatalf("Run ended before announcing: %v\nlog:\n%s", <-ran, log.String())
		}
		first = l
	case <-time.After(20 * time.Second):
		t.Fatalf("no announcement in 20 s\nlog:\n%s", log.String())
	}
	if !strings.HasPrefix(first, ListenAnnouncePrefix) {
		t.Fatalf("first line %q", first)
	}
	if took := time.Since(begin); took > 12*time.Second {
		t.Errorf("startup took %s, want under the two clocks (1.5 s + 9 s) plus margin", took)
	}
	url := strings.TrimSpace(strings.TrimPrefix(first, ListenAnnouncePrefix))
	resp, err := http.Get(url + "/v1/health")
	if err != nil {
		t.Fatal(err)
	}
	var h HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&h)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(h.ExecProfiles) != 0 || len(h.ExecProfilesExcluded) != 1 {
		t.Fatalf("profiles %+v excluded %+v, want proj left out and nothing served", h.ExecProfiles, h.ExecProfilesExcluded)
	}
	x := h.ExecProfilesExcluded[0]
	if x.Project != "proj" || !x.OutOfTime || !strings.Contains(x.Reason, "the startup budget is spent") {
		t.Fatalf("exclusion = %+v, want proj, out_of_time, the probe's own words", x)
	}
	if !strings.Contains(log.String(), `exec profile "proj": `+broker.ExecProbeDeadlineObservation+
		" Review the failure reasons and restart serve if appropriate.") {
		t.Fatalf("log must report observations only:\n%s", log.String())
	}
	cancel()
	select {
	case err := <-ran:
		if err != nil {
			t.Fatalf("Run = %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Run did not stop")
	}
}

// TestRegistryPhaseHasItsOwnClock spends the startup clock at the construction
// layer, with what the HTTP view cannot show: the registry phase's other two
// steps (the template report and the integrity sweep) share its clock, so
// the sweep is complete and finds no hole; and an exec call for the left-out
// profile refuses with its name and the probe's reason (candidate 4, already
// built by row 189: fail-closed, nothing runs). Control: the same rig with
// the probe client answering, where the pass-through srt shim fails real arms
// at once: the profile is excluded with failures recorded before the deadline.
// Mutants: the registry clock started from the startup clock (bootstrap
// fails); the sweep left on the startup clock (a "scan" hole); a profile
// served after an out-of-time probe; the deadline never handed to the prober
// or ignored by fail (out_of_time false).
func TestRegistryPhaseHasItsOwnClock(t *testing.T) {
	for _, tc := range []struct {
		name      string
		slow      bool
		startup   time.Duration
		wantSpent bool
	}{
		{"probe-spends-the-startup-clock", true, 1500 * time.Millisecond, true},
		{"control-an-arm-refuses", false, 9 * time.Second, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, _ := loadRig(t, tc.slow)
			d, err := startDaemon(boundedTestContext(t), cfg, tc.startup, 9*time.Second)
			if err != nil {
				t.Fatalf("startDaemon = %v", err)
			}
			t.Cleanup(func() { _ = d.Close() })
			if rep := d.IntegrityReport(); !rep.Complete || len(rep.Holes) != 0 {
				t.Fatalf("integrity sweep = %+v, want complete with no hole", rep)
			}
			rec := httptest.NewRecorder()
			d.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
			var h HealthResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &h); err != nil {
				t.Fatal(err)
			}
			if len(h.ExecProfiles) != 0 || len(h.ExecProfilesExcluded) != 1 || h.ExecProfilesExcluded[0].OutOfTime != tc.wantSpent {
				t.Fatalf("profiles %+v excluded %+v, want proj left out with out_of_time=%t", h.ExecProfiles, h.ExecProfilesExcluded, tc.wantSpent)
			}
			got := d.workspace.execHandler("ep1", filepath.Join(d.workspace.root, "ep1"))
			refusal, ok := got.(broker.ExecRefusalHandler)
			if !ok || !strings.Contains(refusal.Why, `exec profile "proj" of episode "ep1" was left out at serve startup: exec startup probe failed`) {
				t.Fatalf("exec handler for ep1 = %#v, want the left-out refusal naming proj", got)
			}
		})
	}
}

// TestRegistryPhaseClockIsBounded: the registry phase's clock is a real
// bound, not just a fresh context. A registry budget of 1 ns refuses startup
// at registry-bootstrap with the deadline; the control (9 s) starts. Mutant:
// the registry context without a timeout, and the 1 ns row starts.
func TestRegistryPhaseClockIsBounded(t *testing.T) {
	for _, tc := range []struct {
		budget  time.Duration
		wantErr bool
	}{{time.Nanosecond, true}, {9 * time.Second, false}} {
		cfg := Config{DBPath: filepath.Join(canonicalTemp(t), "world.db"), BindHost: DefaultBindHost, ErrorLog: io.Discard}
		d, err := startDaemon(boundedTestContext(t), cfg, 9*time.Second, tc.budget)
		if !tc.wantErr {
			if err != nil {
				t.Fatalf("registry budget %s: %v", tc.budget, err)
			}
			_ = d.Close()
			continue
		}
		var startup *StartupError
		if !errors.As(err, &startup) || startup.Stage != StageRegistry || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("registry budget %s: startup = %v, want a registry-bootstrap deadline", tc.budget, err)
		}
		if startup.Settled != nil {
			<-startup.Settled
		}
	}
}

// TestServeStartupClocksArePinned pins the production values the design
// derives (§3.2): HEAD's 9 s startup clock unchanged, and a 9 s registry
// clock, 18 s of ctx-aware work; ctx-free work is separately unbounded.
func TestServeStartupClocksArePinned(t *testing.T) {
	if startupBudget != 9*time.Second || registryPhaseBudget != 9*time.Second {
		t.Fatalf("startup clocks = %s + %s, want 9 s + 9 s", startupBudget, registryPhaseBudget)
	}
}

// The two phase parents may differ in deadline, but serve cancellation must
// still reach the fresh registry phase. A live startup context paired with a
// canceled serve parent isolates that obligation from pre-registry work.
func TestRegistryPhasePreservesParentCancellation(t *testing.T) {
	cfg := Config{DBPath: filepath.Join(canonicalTemp(t), "world.db"), BindHost: DefaultBindHost, ErrorLog: io.Discard}
	parent, cancel := context.WithCancel(boundedTestContext(t))
	cancel()
	d, err := newDaemon(boundedTestContext(t), parent, 9*time.Second, cfg)
	if d != nil {
		_ = d.Close()
		t.Fatal("canceled phase parent allowed startup")
	}
	var startup *StartupError
	if !errors.Is(err, context.Canceled) || !errors.As(err, &startup) || startup.Stage != StageRegistry {
		t.Fatalf("startup = %v, want parent cancellation at registry-bootstrap", err)
	}
	if startup.Settled != nil {
		select {
		case <-startup.Settled:
		case <-time.After(3 * time.Second):
			t.Fatal("canceled startup cleanup did not settle")
		}
	}
	s, openErr := store.Open(cfg.DBPath)
	if openErr != nil {
		t.Fatalf("writer lock after canceled startup: %v", openErr)
	}
	_ = s.Close()
}

func TestExecExclusionLogIsObservationOnly(t *testing.T) {
	var out strings.Builder
	x := broker.ExecProfileExclusion{Project: "late", File: "/p/late.json", Reason: "genuine refusal: sibling write succeeded", OutOfTime: true}
	logExecProfileExclusion(&out, x)
	want := "ailang-worldd: exec profile \"late\" (/p/late.json) left out: genuine refusal: sibling write succeeded\n" +
		"ailang-worldd: exec profile \"late\": " + broker.ExecProbeDeadlineObservation + " Review the failure reasons and restart serve if appropriate.\n"
	if out.String() != want {
		t.Fatalf("log must preserve refusal and report observations only: %q", out.String())
	}
}

func TestHealthExecExclusionIsAdditive(t *testing.T) {
	x := HealthExecExclusion{Project: "p", File: "/p/p.json", Reason: "refused"}
	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"project":"p","file":"/p/p.json","reason":"refused"}` {
		t.Fatalf("unchanged exclusion bytes: %s", b)
	}
	x.OutOfTime = true
	b, err = json.Marshal(x)
	if err != nil || !strings.Contains(string(b), `"out_of_time":true`) {
		t.Fatalf("additive true flag: %s, %v", b, err)
	}
}
