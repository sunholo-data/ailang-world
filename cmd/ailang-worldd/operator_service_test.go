package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

type operatorServiceFixture struct {
	*serviceFixture
	config     runnerConfig
	configFile string
	heartbeat  store.OperatorHeartbeat
	homeJob    *fakeServiceController
	health     daemon.HealthResponse
}

func operatorService(t *testing.T) *operatorServiceFixture {
	t.Helper()
	f := newServiceFixture(t)
	closeWorld, err := storeForServiceTest(f.db)
	if err != nil {
		t.Fatal(err)
	}
	closeWorld()
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("World fixture: %d %s", code, msg)
	}
	r, err := readServiceRecord(f.env.paths)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.UpgradeCredentials(r.Credentials); err != nil {
		t.Fatal(err)
	}
	world, creds := r.Store, r.Credentials // fixture paths are physically named
	f.env.runnerLogin = func() (int, string, error) { return os.Getuid(), f.home, nil }
	o := &operatorServiceFixture{serviceFixture: f, homeJob: f.ctl, configFile: filepath.Join(f.home, "proposed.json"),
		config: runnerConfig{Home: f.home, UID: os.Getuid(), Address: r.Address, World: world, Credentials: creds,
			Executables: map[string]string{r.Binary: r.BinarySHA256}, PresetNames: []string{"sol"},
			ToolDirs: []string{filepath.Dir(r.Binary)}, Backend: "headless"},
		heartbeat: store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fake-runner"}, Protocol: store.OperatorRunnerProtocol, Build: r.BinarySHA256},
		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
	f.env.controller = f.ctl
	f.env.wait = 250 * time.Millisecond
	clock := time.Unix(f.env.terminal.now(), 0)
	f.env.runnerNow = func() time.Time { return clock }
	f.env.runnerSleep = func(d time.Duration) { clock = clock.Add(d) }
	f.ctl.onBoot = func() error {
		o.heartbeat.Owner.Generation = 0
		return f.env.runnerTick(o.config, &o.heartbeat)
	}
	f.env.client = &http.Client{Transport: serviceTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != r.Address || req.URL.Path != "/v1/health" {
			t.Fatalf("unexpected health URL: %s", req.URL)
		}
		if !f.serving {
			return nil, errors.New("fake daemon went away")
		}
		b, _ := json.Marshal(o.health)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(b)), Header: make(http.Header)}, nil
	})}
	f.writes = nil
	return o
}

func (f *operatorServiceFixture) command(args ...string) (int, string) {
	return f.run(append([]string{"runner"}, args...)...)
}
func (f *operatorServiceFixture) install(t *testing.T) (int, string) {
	t.Helper()
	b, err := json.Marshal(f.config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.configFile, b, 0600); err != nil {
		t.Fatal(err)
	}
	return f.command("install", "--config", f.configFile)
}
func (f *operatorServiceFixture) mustInstall(t *testing.T) {
	t.Helper()
	if c, s := f.install(t); c != 0 {
		t.Fatalf("install: %d %s", c, s)
	}
}

func (f *operatorServiceFixture) refused(t *testing.T, want string, code int, message string) {
	t.Helper()
	if code != 1 || !strings.Contains(message, want) || len(f.writes) != 0 || f.ctl.boots != 0 {
		t.Fatalf("refusal %q: %d %s writes=%v", want, code, message, f.writes)
	}
}

func TestOperatorServiceFence(t *testing.T) {
	f := operatorService(t)
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
	c, msg := f.install(t)
	f.refused(t, "aborted", c, msg)
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }
	f.mustInstall(t)
	p := f.env.runnerPaths()
	b, err := os.ReadFile(p.Plist)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.Plist, append(b, ' '), 0600); err != nil {
		t.Fatal(err)
	}
	before := serviceFileHash(p.Plist)
	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "differs from approved") || f.ctl.boots != 0 || serviceFileHash(p.Plist) != before {
		t.Fatalf("%d %s", c, s)
	}
	for _, words := range [][]string{{"service", "runner", "start"}, {"operator-runner"}} {
		got := verbs.AdmitOperator(words, verbs.OperatorPolicy{Start: "A", Finish: "A"})
		if got.Intent != nil {
			t.Fatalf("bar admitted %q", words)
		}
	}
}

func TestOperatorServiceOwnership(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		damage     func(*testing.T, *operatorServiceFixture)
	}{
		{"root", "non-root", func(t *testing.T, f *operatorServiceFixture) {
			f.env.runnerLogin = func() (int, string, error) { return 0, f.home, nil }
			f.config.UID = 0
		}},
		{"uid", "non-root", func(t *testing.T, f *operatorServiceFixture) { f.config.UID++ }},
		{"home", "login home", func(t *testing.T, f *operatorServiceFixture) { f.config.Home = t.TempDir() }},
		{"writable tool dir", "writable by others", func(t *testing.T, f *operatorServiceFixture) {
			p := filepath.Join(f.home, "loose")
			if err := os.Mkdir(p, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(p, 0770); err != nil {
				t.Fatal(err)
			}
			f.config.ToolDirs = []string{p}
		}},
		{"other writable", "writable by others", func(t *testing.T, f *operatorServiceFixture) {
			p := filepath.Join(f.home, "other-writable")
			if err := os.Mkdir(p, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(p, 0702); err != nil {
				t.Fatal(err)
			}
			f.config.ToolDirs = []string{p}
		}},
		{"symlink tool", "wrong type", func(t *testing.T, f *operatorServiceFixture) {
			p := filepath.Join(f.home, "linked")
			if err := os.Symlink(f.env.paths.Binary, p); err != nil {
				t.Fatal(err)
			}
			f.config.Executables[p] = serviceFileHash(f.env.paths.Binary)
		}},
		{"changed hash", "executable changed", func(t *testing.T, f *operatorServiceFixture) {
			f.config.Executables[f.env.paths.Binary] = "sha256:other"
		}},
		{"missing executable", "missing installed", func(t *testing.T, f *operatorServiceFixture) { delete(f.config.Executables, f.env.paths.Binary) }},
		{"non executable", "executable changed", func(t *testing.T, f *operatorServiceFixture) {
			if err := os.Chmod(f.env.paths.Binary, 0600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := operatorService(t)
			tc.damage(t, f)
			c, msg := f.install(t)
			f.refused(t, tc.want, c, msg)
		})
	}
}

func TestRunnerPathRefusalClasses(t *testing.T) {
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	uid := os.Getuid()
	if err := runnerPath(home, home, uid, true, 0); err != nil {
		t.Fatalf("valid target directory refused: %v", err)
	}
	link := filepath.Join(home, "linked-dir")
	if err := os.Symlink(home, link); err != nil {
		t.Fatal(err)
	}
	if err := runnerPath(home, link, uid, true, 0); err == nil || err.Error() != "path has wrong type or is a symlink" {
		t.Errorf("symlink refusal = %v, want wrong type", err)
	}

	file := filepath.Join(home, "writable-file")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(file, 0666); err != nil {
		t.Fatal(err)
	}
	if err := runnerPath(home, file, uid, false, 0); err == nil || err.Error() != "path writable by others" {
		t.Errorf("regular file refusal = %v, want writable by others", err)
	}
}

func TestOperatorServiceConfig(t *testing.T) {
	f := operatorService(t)
	var b bytes.Buffer
	if c := runOperatorRunner([]string{"episode", "start", "evil"}, &b, f.env); c != exitUsage || b.Len() != 0 {
		t.Fatalf("accepted verb argv: %d %q", c, b.String())
	}
	raw, err := json.Marshal(f.config)
	if err != nil {
		t.Fatal(err)
	}
	for _, payload := range []string{string(raw[:len(raw)-1]) + `,"token":"canary-secret"}`, string(raw) + " {}", "{"} {
		if err := os.WriteFile(f.configFile, []byte(payload), 0600); err != nil {
			t.Fatal(err)
		}
		if c, s := f.command("install", "--config", f.configFile); c != 1 || !strings.Contains(s, "invalid runner setup") || strings.Contains(s, "canary-secret") || len(f.writes) != 0 {
			t.Fatalf("%d %s", c, s)
		}
	}
	f.config.Backend = "shell"
	c, msg := f.install(t)
	f.refused(t, "backend", c, msg)
}

func TestOperatorServicePair(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		change     func(*operatorServiceFixture)
	}{
		{"world", "World store mismatch", func(f *operatorServiceFixture) { f.health.DBPath = f.config.Credentials }},
		{"credentials", "credentials store mismatch", func(f *operatorServiceFixture) { f.health.CredentialsDB = f.config.World }},
		{"build", "restart pending", func(f *operatorServiceFixture) { f.health.Build.ExecutableSHA256 = "sha256:other" }},
		{"loopback", "numeric loopback", func(f *operatorServiceFixture) { f.config.Address = "192.0.2.1:7644" }},
		{"named loopback", "numeric loopback", func(f *operatorServiceFixture) { f.config.Address = "localhost:7644" }},
		{"distinct", "distinct canonical files", func(f *operatorServiceFixture) { f.config.Credentials = f.config.World }},
		{"health status", "daemon unavailable", func(f *operatorServiceFixture) { f.health.Status = "starting" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := operatorService(t)
			f.mustInstall(t)
			f.writes = nil
			tc.change(f)
			if tc.name == "loopback" || tc.name == "named loopback" || tc.name == "distinct" {
				c, msg := f.install(t)
				f.refused(t, tc.want, c, msg)
				return
			}
			c, msg := f.command("start")
			f.refused(t, tc.want, c, msg)
			if err := f.env.runnerTick(f.config, &f.heartbeat); err == nil || f.heartbeat.Owner.Generation != 0 {
				t.Fatal("incompatible daemon allowed heartbeat")
			}
		})
	}
}

func TestOperatorServiceStartup(t *testing.T) {
	f := operatorService(t)
	f.mustInstall(t)
	f.serving = false
	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "daemon unavailable") || f.ctl.boots != 0 {
		t.Fatalf("%d %s", c, s)
	}
	f.serving = true
	credSQL(t, f.config.Credentials, "DROP TABLE IF EXISTS allowance_draws; DROP TABLE IF EXISTS standing_allowances; DROP TABLE operator_action_steps; DROP TABLE operator_actions; DROP TABLE operator_runner_state; PRAGMA user_version=9")
	if c, msg := f.command("start"); c != 1 || !strings.Contains(msg, "schema mismatch") || f.ctl.boots != 0 {
		t.Fatalf("schema skew: %d %s", c, msg)
	}
	if _, err := store.UpgradeCredentials(f.config.Credentials); err != nil {
		t.Fatal(err)
	}
	f.ctl.bootErr = errors.New("fake bootstrap failure")
	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "bootstrap") || f.ctl.boots != 1 || f.ctl.loaded {
		t.Fatalf("%d %s", c, s)
	}
	f.ctl.bootErr = nil
	if c, s := f.command("start"); c != 0 || !f.ctl.loaded {
		t.Fatalf("%d %s", c, s)
	}
	t.Run("gate", testOperatorServiceStartupGate)
}

func TestOperatorServiceUpgradePairFence(t *testing.T) {
	for _, tc := range []struct{ arm, want string }{
		{"installed", "operator runner is still installed"},
		{"plist-only", "operator runner is still installed"},
		{"registered", "cannot verify that the operator runner is unloaded"},
		{"locked", "operator runner lifecycle is busy or unavailable"},
		{"inspect error", "cannot verify that the operator runner is unloaded"},
	} {
		t.Run(tc.arm, func(t *testing.T) {
			f := operatorService(t)
			f.mustInstall(t)
			p := f.env.runnerPaths()
			switch tc.arm {
			case "installed":
				if err := os.Remove(p.Plist); err != nil {
					t.Fatal(err)
				}
			case "plist-only":
				if err := os.Remove(p.Record); err != nil {
					t.Fatal(err)
				}
			default:
				if c, msg := f.command("uninstall"); c != 0 {
					t.Fatal(c, msg)
				}
			}
			f.env.controller = f.homeJob
			f.homeJob.runner = tc.arm == "registered"
			if tc.arm == "inspect error" {
				f.homeJob.runnerPrintErr = errors.New("supervisor query failed")
			}
			f.homeJob.onKick = func() error {
				f.health.Build.ExecutableSHA256 = f.homeJob.running
				return nil
			}
			if tc.arm == "locked" {
				lock, err := runnerLifecycleLock(filepath.Dir(f.env.paths.Record))
				if err != nil {
					t.Fatal(err)
				}
				defer lock.Close()
			}
			if err := os.WriteFile(f.source, []byte("upgrade candidate"), 0700); err != nil {
				t.Fatal(err)
			}
			before := serviceFileHash(f.env.paths.Binary)
			f.writes = nil
			c, msg := f.run("upgrade", "--binary", f.source)
			after := serviceFileHash(f.env.paths.Binary)
			if c != 1 || !strings.Contains(msg, tc.want) || f.homeJob.kicks != 0 || after != before || len(f.writes) != 0 {
				t.Fatalf("want refusal %q; got %d %s kicks=%d binary=%s → %s writes=%v", tc.want, c, msg, f.homeJob.kicks, before, after, f.writes)
			}
		})
	}
}

func TestOperatorServiceCancelConfirmation(t *testing.T) {
	f := operatorService(t)
	f.mustInstall(t)
	action := approval.OperatorAction{Schema: approval.OperatorActionSchema,
		ID: "oa-00000000000000000000000000000001", Kind: "start", Episode: "pending",
		Store: f.config.Credentials, Review: "review", Task: "test task",
		Options:    approval.OperatorOptions{Preset: "sol", Headless: true},
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
	action.Argv = approval.OperatorArgv(action.Kind, action.Episode, action.Options)
	ctx, cancel := budgetContext(5 * time.Second)
	defer cancel()
	now := f.env.terminal.now()
	if err := runnerWrite(f.config.Credentials, func(c *store.Credentials) error {
		_, err := c.RequestOperatorAction(ctx, store.OperatorRequest{Action: action,
			OperationID: "00000000000000000000000000000001", CreatedAt: now, ExpiresAt: now + 600}, now)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	ro := openCredsReadOnly(t, f.config.World)
	defer ro.Close()
	checkState := func(want string) {
		t.Helper()
		row, err := ro.OperatorAction(ctx, action.ID, now)
		if err != nil || row.State != want {
			t.Fatalf("action state = %q, want %q: %v", row.State, want, err)
		}
	}
	checkState("requested")
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
	if c, msg := f.command("cancel", action.ID); c != 1 || !strings.Contains(msg, "aborted (not confirmed)") {
		t.Errorf("declined cancel: %d %s", c, msg)
	}
	checkState("requested")
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }
	if c, msg := f.command("cancel", action.ID); c != 0 {
		t.Fatalf("confirmed cancel: %d %s", c, msg)
	}
	checkState("cancelled")
}

func TestOperatorRunnerKernelLock(t *testing.T) {
	p := filepath.Join(t.TempDir(), "runner.lock")
	f, err := runnerKernelLock(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	done := make(chan error, 1)
	go func() {
		second, err := runnerKernelLock(p)
		if second != nil {
			second.Close()
		}
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("contender: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("runner lock waited")
	}
}

func TestOperatorServiceWriteScope(t *testing.T) {
	f := operatorService(t)
	base := f.env.paths
	p := f.env.runnerPaths()
	read := func() [5]string {
		return [5]string{serviceFileHash(base.Plist), serviceFileHash(base.Record), serviceFileHash(base.Binary),
			serviceFileHash(f.config.World), serviceFileHash(f.config.Credentials)}
	}
	want := read()
	f.mustInstall(t)
	for _, path := range []string{p.Plist, p.Record, p.Log} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0600 {
			t.Fatalf("private mode %s: %v %v", path, info, err)
		}
	}
	if info, err := os.Stat(filepath.Dir(p.Log)); err != nil || info.Mode().Perm() != 0700 {
		t.Fatalf("spool: %v %v", info, err)
	}
	if c, s := f.command("start"); c != 0 {
		t.Fatalf("start: %d %s", c, s)
	}
	want[4] = serviceFileHash(f.config.Credentials)
	for _, verb := range []string{"stop", "uninstall"} {
		if c, s := f.command(verb); c != 0 {
			t.Fatalf("%s: %d %s", verb, c, s)
		}
	}
	if f.homeJob.stops != 0 || f.ctl.loaded || want != read() || serviceFileHash(p.Plist) != "" || serviceFileHash(p.Record) != "" || serviceFileHash(p.Log) == "" {
		t.Fatal("companion changed base files, jobs, or retained evidence")
	}
	for _, path := range f.writes {
		if !within(f.home, path) {
			t.Fatalf("write escaped fake HOME: %s", path)
		}
	}
}

func TestOperatorServiceEnvironment(t *testing.T) {
	f := operatorService(t)
	t.Setenv("GH_TOKEN", "canary-secret")
	t.Setenv("PATH", "/canary-path")
	f.mustInstall(t)
	p := f.env.runnerPaths()
	b, err := os.ReadFile(p.Plist)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{runnerLabel, "<string>operator-runner</string>", "<key>HOME</key><string>" + f.home + "</string>", "<key>PATH</key><string>" + f.config.ToolDirs[0] + ":" + runnerSystemPath + "</string>", "GH_PROMPT_DISABLED", "GIT_TERMINAL_PROMPT", serviceFileHash(p.Record)} {
		if !strings.Contains(s, want) {
			t.Fatalf("plist lacks %q: %s", want, s)
		}
	}
	for _, path := range []string{p.Plist, p.Record, p.Log} {
		raw, err := os.ReadFile(path)
		if err != nil || bytes.Contains(raw, []byte("canary")) {
			t.Fatalf("ambient environment leaked to %s: %v", path, err)
		}
	}
}

func TestOperatorStatusDisconnected(t *testing.T) {
	f := operatorService(t)
	f.mustInstall(t)
	if c, s := f.command("start"); c != 0 {
		t.Fatalf("%d %s", c, s)
	}
	if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
		t.Fatal(err)
	}
	if string(f.heartbeat.Catalogue) != `["sol"]` {
		t.Fatalf("unexpected catalogue: %q", f.heartbeat.Catalogue)
	}
	known := "runner last seen " + strconv.FormatInt(f.heartbeat.SeenAt, 10) + ", generation 1"
	if c, s := f.command("status"); c != 0 || !strings.Contains(s, known) || !strings.Contains(s, "start: ready") || !strings.Contains(s, "finish blocked: github-auth") {
		t.Fatalf("%d %s", c, s)
	}
	f.serving = false
	if c, s := f.command("status"); c != 1 || !strings.Contains(s, known) || !strings.Contains(s, "daemon disconnected") || strings.Contains(s, "start: ready") {
		t.Fatalf("%d %s", c, s)
	}
}
