package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

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
	world, err := store.CanonicalPath(r.Store)
	if err != nil {
		t.Fatal(err)
	}
	creds, err := store.CanonicalPath(r.Credentials)
	if err != nil {
		t.Fatal(err)
	}
	f.env.runnerLogin = func() (int, string, error) { return os.Getuid(), f.home, nil }
	o := &operatorServiceFixture{serviceFixture: f, homeJob: f.ctl, configFile: filepath.Join(f.home, "proposed.json"),
		config:    runnerConfig{f.home, os.Getuid(), r.Address, world, creds, map[string]string{r.Binary: r.BinarySHA256}, []string{"sol"}, []string{filepath.Dir(r.Binary)}, "headless"},
		heartbeat: store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fake-runner"}, Protocol: store.OperatorRunnerProtocol, Build: r.BinarySHA256},
		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
	f.env.controller = f.ctl
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

func TestOperatorServiceFence(t *testing.T) {
	f := operatorService(t)
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
	if c, s := f.install(t); c != 1 || !strings.Contains(s, "aborted") || len(f.writes) != 0 || f.ctl.boots != 0 {
		t.Fatalf("%d %s writes=%v", c, s, f.writes)
	}
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
			if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
				t.Fatalf("%d %s writes=%v", c, s, f.writes)
			}
		})
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
	if c, s := f.install(t); c != 1 || !strings.Contains(s, "backend") || len(f.writes) != 0 {
		t.Fatalf("%d %s", c, s)
	}
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := operatorService(t)
			f.mustInstall(t)
			f.writes = nil
			tc.change(f)
			if tc.name == "loopback" {
				if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 {
					t.Fatalf("%d %s", c, s)
				}
				return
			}
			if c, s := f.command("start"); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
				t.Fatalf("%d %s", c, s)
			}
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
	f.ctl.bootErr = errors.New("fake bootstrap failure")
	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "bootstrap") || f.ctl.boots != 1 || f.ctl.loaded {
		t.Fatalf("%d %s", c, s)
	}
	f.ctl.bootErr = nil
	if c, s := f.command("start"); c != 0 || !f.ctl.loaded {
		t.Fatalf("%d %s", c, s)
	}
}

func TestOperatorServiceWriteScope(t *testing.T) {
	f := operatorService(t)
	base := f.env.paths
	p := f.env.runnerPaths()
	files := []string{base.Plist, base.Record, base.Binary, f.config.World, f.config.Credentials}
	read := func() []string {
		var h []string
		for _, p := range files {
			h = append(h, serviceFileHash(p))
		}
		return h
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
	for _, verb := range []string{"start", "stop", "uninstall"} {
		if c, s := f.command(verb); c != 0 {
			t.Fatalf("%s: %d %s", verb, c, s)
		}
	}
	if f.homeJob.stops != 0 || f.ctl.loaded || !reflect.DeepEqual(want, read()) || serviceFileHash(p.Plist) != "" || serviceFileHash(p.Record) != "" || serviceFileHash(p.Log) == "" {
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
	if c, s := f.command("status"); c != 1 || !strings.Contains(s, "heartbeat unavailable") {
		t.Fatalf("%d %s", c, s)
	}
	if c, s := f.command("start"); c != 0 {
		t.Fatalf("%d %s", c, s)
	}
	if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
		t.Fatal(err)
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
