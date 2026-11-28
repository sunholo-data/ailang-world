package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

const fakeLaunchMarker = "#!/bin/sh\nprintf x >> \"${0%/*}/ran\"\n"

type operatorLaunchRig struct {
	*worktreeStartRig
	c      runnerConfig
	h      store.OperatorHeartbeat
	runner *os.File
	ctx    context.Context
	db     *store.Credentials
	self   daemon.SessionSelf
	now    int64
	token  string
}

func newOperatorLaunchRig(t *testing.T, script string) *operatorLaunchRig {
	t.Helper()
	return newOperatorLaunchRigWithTask(t, script, operatorStartTask)
}

func newOperatorLaunchRigWithTask(t *testing.T, script, task string) *operatorLaunchRig {
	t.Helper()
	return newOperatorLaunchRigWithRemote(t, script, task, "")
}

func newOperatorLaunchRigWithRemote(t *testing.T, script, task, remote string) *operatorLaunchRig {
	t.Helper()
	r, c, h, runner := operatorActivationFixtureWithRemote(t, script, task, remote)
	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: activationTestContext(t), now: r.src.clock()}
	r.src.clock = func() int64 { return f.now }
	r.env.terminal.now = r.src.clock
	var err error
	f.db, err = store.OpenCredentialsReadOnly(c.Credentials)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.db.Close() })
	health := r.env.client.Transport
	r.env.client.Transport = serviceTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/session/self" {
			return health.RoundTrip(req)
		}
		if req.URL.Host != c.Address || f.token == "" || req.Header.Get("Authorization") != "Bearer "+f.token {
			t.Error("launch self did not use the pinned session")
			return nil, errors.New("fixture incorrectly bound")
		}
		body, _ := json.Marshal(f.self)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	})
	r.env.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
		s := f.receipt(t)
		if s.Phase != "reserved" || s.Key.Version != 1 || s.Receipt.PID != 0 {
			t.Fatalf("supervisor crossed reservation boundary: %+v", s)
		}
		return r.env.executeOperatorLaunch(f.ctx, job, log)
	}
	if out := f.tick(t); out != "" {
		t.Fatalf("preparation tick: %s", out)
	}
	f.readPrepared(t)
	return f
}

func (f *operatorLaunchRig) readPrepared(t *testing.T) {
	t.Helper()
	s, err := f.db.ReadOperatorStartRecovery(f.ctx, f.action.ID)
	if err != nil || s.Window != store.OperatorStartPrepared {
		t.Fatalf("not prepared: %+v %v", s, err)
	}
	files, _, err := f.db.ReadOperatorFiles(f.ctx, f.action.ID)
	if err != nil {
		t.Fatal(err)
	}
	token, err := os.ReadFile(files.Mint.TokenPath)
	if err != nil || len(token) != 64 {
		t.Fatalf("real mint did not write a token: %v", err)
	}
	f.token = string(token)
	f.self = daemon.SessionSelf{Episode: f.action.Episode, Worktree: f.path(), ExpiresAt: files.Session.ExpiresAt, SecondsLeft: files.Session.ExpiresAt - f.now}
}

func (f *operatorLaunchRig) tick(t *testing.T) string {
	t.Helper()
	ticks := make(chan time.Time)
	close(ticks)
	var out syncBuffer
	if rc := f.env.runOperatorPreparationLoop(&out, f.c, f.h, f.src, nil, ticks, f.runner); rc != exitOK {
		t.Fatalf("runner tick %d: %s", rc, out.String())
	}
	return out.String()
}

func (f *operatorLaunchRig) receipt(t *testing.T) store.OperatorLaunchSnapshot {
	t.Helper()
	s, err := f.db.ReadOperatorLaunch(f.ctx, f.action.ID)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func (f *operatorLaunchRig) assertUnspawned(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(f.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fake child ran: %v", err)
	}
	if _, err := f.db.ReadOperatorLaunch(f.ctx, f.action.ID); !errors.Is(err, store.ErrOperatorMissing) {
		t.Fatalf("unexpected reservation or receipt: %v", err)
	}
}

func TestOperatorSpawnRefused(t *testing.T) {
	f := newOperatorLaunchRig(t, "#!/no-world-fixture-interpreter\nprintf unexpected\n")
	before := activationEvidence(t, f.worktreeStartRig)
	out := f.tick(t)
	s := f.receipt(t)
	after := activationEvidence(t, f.worktreeStartRig)
	if !strings.Contains(out, ": spawn-refused") || !s.NoSpawn || s.Receipt != (store.OperatorLaunchReceipt{}) ||
		!reflect.DeepEqual(before[1:], after[1:]) {
		t.Fatalf("known exec refusal: %s %+v", out, s)
	}
	if _, err := os.Lstat(filepath.Join(f.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("exec refusal spawned a fixture: %v", err)
	}
	start, err := f.db.EpisodeStart(f.ctx, f.action.Episode)
	if err != nil || len(start.Launches) != 0 {
		t.Fatalf("exec refusal created launch history: %v %v", start, err)
	}
}

func TestOperatorLaunchDispatchRechecks(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		drift      func(binary, dir string) error
	}{
		{
			name: "binary_digest",
			want: (&approval.BindingError{Field: "launch_binary_sha256"}).Error(),
			drift: func(binary, _ string) error {
				return os.WriteFile(binary, []byte(fakeLaunchMarker+"# changed after preflight\n"), 0700)
			},
		},
		{
			name: "nonempty_directory",
			want: "is not empty (F8)",
			drift: func(_, dir string) error {
				return os.WriteFile(filepath.Join(dir, "unexpected"), []byte("drift"), 0600)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newOperatorLaunchRig(t, fakeLaunchMarker)
			f.assertUnspawned(t)
			beforeAction, err := f.db.OperatorAction(f.ctx, f.action.ID, f.now)
			if err != nil {
				t.Fatal(err)
			}
			var reserved store.OperatorLaunchSnapshot
			var dispatchErr error
			supervise := f.env.operatorSupervisor
			f.env.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
				reserved = f.receipt(t)
				dispatchErr = supervise(job, log)
				return dispatchErr
			}
			f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string,
				out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
				// Reservation and worker preflight have passed. Drift here exercises Before.
				if err := tc.drift(argv[0], dir); err != nil {
					t.Fatal(err)
				}
				return agentrun.RunObserved(ctx, argv, dir, env, out, stderr, wall, hooks)
			}
			out := f.tick(t)
			if dispatchErr == nil || !strings.Contains(dispatchErr.Error(), tc.want) || !strings.Contains(out, "blocked or uncertain") {
				t.Errorf("dispatch refusal = %v, output = %q, want %q", dispatchErr, out, tc.want)
			}
			if _, err := os.Stat(filepath.Join(f.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("fake child ran: %v", err)
			}
			if after := f.receipt(t); reserved.Key.Version != 1 || after != reserved || after.Receipt != (store.OperatorLaunchReceipt{}) {
				t.Errorf("dispatch changed reservation or recorded a receipt: before %+v, after %+v", reserved, after)
			}
			afterAction, err := f.db.OperatorAction(f.ctx, f.action.ID, f.now)
			if err != nil || !reflect.DeepEqual(beforeAction, afterAction) {
				t.Errorf("dispatch changed action state: %s -> %s (%v)", beforeAction.State, afterAction.State, err)
			}
			start, err := f.db.EpisodeStart(f.ctx, f.action.Episode)
			if err != nil || start == nil || len(start.Launches) != 0 {
				t.Errorf("dispatch refusal created launch history: %v (%v)", start, err)
			}
		})
	}
}

func TestOperatorLaunchCrashWindow(t *testing.T) {
	s := store.OperatorLaunchSnapshot{Phase: "observed", Receipt: store.OperatorLaunchReceipt{Run: fixtureActionID, PID: 42, Birth: "first", Status: "launched"}}
	for _, birth := range []string{"reused", "first"} {
		want := "uncertain"
		if birth == "first" {
			want = "launched"
		}
		if got := operatorLaunchProcessStatus(s, func(pid int) (string, bool) { return birth, pid == 42 }); got != want {
			t.Fatalf("birth %q classified %q, want %q", birth, got, want)
		}
	}
}

func TestOperatorLaunchExit(t *testing.T) {
	t.Run("failed", func(t *testing.T) {
		dir := t.TempDir()
		script := filepath.Join(dir, "child")
		if err := os.WriteFile(script, []byte("#!/bin/sh\nexit 7\n"), 0700); err != nil {
			t.Fatal(err)
		}
		stream, err := os.CreateTemp(dir, "stream")
		if err != nil {
			t.Fatal(err)
		}
		ctx := activationTestContext(t)
		result, err := agentrun.RunObserved(ctx, []string{script}, dir, []string{}, stream, stream, time.Minute, agentrun.LaunchHooks{})
		_ = stream.Close()
		history := operatorLaunchHistory(operatorLaunchInputs{sessions: dir}, store.OperatorLaunchReceipt{Status: "failed"}, "ep", stream, stream, time.Time{}, time.Time{})
		if err != nil || operatorLaunchExitStatus(ctx, result) != "failed" || result.ExitCode != 7 || history.FinishReason != "output-failed" {
			t.Fatalf("nonzero real child: %+v %v", result, err)
		}
	})
	t.Run("capped_and_scrubbed", func(t *testing.T) {
		dir := t.TempDir()
		script := filepath.Join(dir, "child")
		if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf 'out %s\\n' \"$WORLD_SESSION\"\nprintf 'err %s\\n' \"$WORLD_SESSION\" >&2\nprintf x >&5\nwhile :; do :; done\n"), 0700); err != nil {
			t.Fatal(err)
		}
		stdout, err := operatorLaunchStage(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer stdout.Close()
		stderr, err := operatorLaunchStage(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer stderr.Close()
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()
		proc := activationTestContext(t)
		deadline, _ := proc.Deadline()
		if err := r.SetReadDeadline(deadline); err != nil {
			t.Fatal(err)
		}
		ctx, cap := context.WithCancel(proc)
		defer cap()
		canary := strings.Repeat("1a", 32)
		hooks := agentrun.LaunchHooks{Files: []*os.File{nil, nil, w}, Started: func(int) error {
			var ready [1]byte
			_, err := io.ReadFull(r, ready[:])
			cap()
			return err
		}}
		argv := []string{script, "sensitive task"}
		result, err := agentrun.RunObserved(ctx, argv, dir, []string{"WORLD_SESSION=" + canary}, stdout, stderr, time.Minute, hooks)
		status := operatorLaunchExitStatus(proc, result)
		if err != nil || status != "capped" {
			t.Fatalf("capped launch: %+v %v", result, err)
		}
		if raw, err := filepath.Glob(filepath.Join(dir, ".operator-output-*")); err != nil || len(raw) != 0 {
			t.Fatalf("unscrubbed stream retained a pathname: %v (%v)", raw, err)
		}
		in := operatorLaunchInputs{sessions: dir, token: canary, argv: argv, prompt: argv[1]}
		history := operatorLaunchHistory(in, store.OperatorLaunchReceipt{Status: status}, "ep", stdout, stderr, time.Time{}, time.Time{})
		for _, path := range []string{"ep.transcript.jsonl", "ep.agent.err"} {
			body, err := os.ReadFile(filepath.Join(dir, path))
			if err != nil || strings.Contains(string(body), canary) || !strings.Contains(string(body), "scrubbed") {
				t.Fatalf("child canary reached %s or scrub control absent: %v", path, err)
			}
		}
		raw, err := json.Marshal(history)
		if err != nil || strings.Contains(string(raw), canary) || strings.Contains(string(raw), argv[1]) ||
			history.TranscriptState != "scrubbed" || history.FinishReason != "capped" {
			t.Fatalf("unredacted launch history: %v", err)
		}
	})
}

func testOperatorStartOnceLaunch(t *testing.T) {
	t.Setenv("WORLD_AMBIENT_TEST_CANARY", "must stay out")
	f := newOperatorLaunchRig(t, fakeLaunchMarker+"[ -e /dev/fd/3 ] && [ -e /dev/fd/4 ]\n")
	f.assertUnspawned(t)
	supervise := f.env.operatorSupervisor
	f.env.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
		before := f.receipt(t)
		var nested syncBuffer
		ticks := make(chan time.Time, 1)
		ticks <- time.Time{}
		close(ticks)
		restarted := f.env
		restarted.operatorDeadline = nil
		if rc := restarted.runOperatorPreparationLoop(&nested, f.c, f.h, f.src, nil, ticks, f.runner); rc != exitOK ||
			strings.Count(nested.String(), "launch-status "+f.action.ID) != 1 || f.receipt(t) != before {
			t.Fatalf("reserved restart crossed spawn boundary: %d %s", rc, &nested)
		}
		return supervise(job, log)
	}
	var budgets []time.Duration
	f.env.operatorDeadline = func(want time.Duration, ctx context.Context) {
		budgets = append(budgets, want)
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > activationTestBudget {
			t.Error("launch missed shared bounded context override")
		}
	}
	f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
		commands, err := launchCommands(f.health, f.in.Profile.Project)
		body, _ := json.Marshal(f.self)
		text, _ := launchContext(body, f.c.Address, f.token)
		want := launchArgv("codex", f.in.Start.Model, launchPrompt(f.action.Task, commands), dir, "", f.c.Address, text, "", nil)
		want[0] = filepath.Join(f.bin, "codex")
		if err != nil || !slices.Equal(want, argv) || !slices.Contains(argv, "read-only") ||
			!slices.Contains(env, "WORLD_SESSION="+f.token) || !slices.Contains(env, "HOME="+f.c.Home) ||
			slices.ContainsFunc(env, func(v string) bool { return strings.Contains(v, "WORLD_AMBIENT_TEST_CANARY") }) ||
			wall != time.Hour || len(hooks.Files) != 2 {
			t.Fatal("managed launch diverged from reviewed V3 confinement or process bounds")
		}
		return agentrun.RunObserved(ctx, argv, dir, env, out, stderr, wall, hooks)
	}
	out := f.tick(t)
	if !strings.Contains(out, ": stopped") {
		t.Fatalf("installed tick did not complete launch: %s", out)
	}
	if !reflect.DeepEqual(budgets, []time.Duration{time.Second, 22 * time.Second, 5 * time.Second, 5 * time.Second, 5 * time.Second}) {
		t.Fatalf("launch context budget calls: %v", budgets)
	}
	before := f.receipt(t)
	f.tick(t)
	ran, err := os.ReadFile(filepath.Join(f.bin, "ran"))
	if err != nil || string(ran) != "x" || f.receipt(t) != before || before.Receipt.Run != f.action.ID ||
		before.Receipt.Executable != operatorFileFact(f.action, "launch_binary_sha256") {
		t.Fatalf("spawn count or identity differed: %q %v", ran, err)
	}
}

func TestOperatorStartBrowser(t *testing.T) {
	// Same browser-supplied task as TestOperatorConfirmHTTP. The launch rig
	// signs and freezes it before preparation reads the durable row.
	task := "# Task\n<script>alert(1)</script> ; $(literal) &"
	f := newOperatorLaunchRigWithTask(t, fakeLaunchMarker+"while [ ! -e release ]; do :; done\nexit 0\n", task)
	f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
		promptArgs := 0
		for _, arg := range argv {
			if strings.Contains(arg, task) {
				promptArgs++
			}
		}
		if promptArgs != 1 || !strings.HasPrefix(argv[len(argv)-1], task) {
			t.Errorf("browser task did not survive launch as one literal prompt argument: %q", argv)
		}
		start := hooks.Started
		hooks.Started = func(pid int) error {
			defer func() { _ = os.WriteFile(filepath.Join(dir, "release"), nil, 0600) }()
			if err := start(pid); err != nil {
				return err
			}
			f.now = f.self.ExpiresAt + 1
			f.env.client.Transport = serviceTransport(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("fixture daemon restarting")
			})
			var status, stderr syncBuffer
			if rc := runServiceRunner([]string{"launch-status", f.action.ID}, &status, &stderr, f.env); rc != exitOK ||
				status.String() != fmt.Sprintf("operator launch %s: launched\n", f.action.ID) {
				t.Errorf("restart failed to observe live receipt: %s %s", &status, &stderr)
			}
			return nil
		}
		return agentrun.RunObserved(ctx, argv, dir, env, out, stderr, wall, hooks)
	}
	before, _, err := f.db.ReadOperatorFiles(f.ctx, f.action.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out := f.tick(t); !strings.Contains(out, ": launched") || !strings.Contains(out, ": stopped") {
		t.Fatalf("launch could not finish after expiry: %s", out)
	}
	s, ok, err := f.db.ResolveSession(f.ctx, before.Mint.CredentialID)
	h, histErr := f.db.EpisodeStart(f.ctx, f.action.Episode)
	if err != nil || histErr != nil || !ok || s != before.Session || len(h.Launches) != 1 || f.receipt(t).Receipt.Run != f.action.ID {
		t.Fatalf("restart changed start history or duplicated mint: %v %v %t", err, histErr, ok)
	}
}
