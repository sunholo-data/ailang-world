package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sunholo-data/ailang-world/host/store"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

type fakeServiceController struct {
	loaded       bool
	boots, stops int
	key          string
	bootErr      error
	kicks        int
	binary       string
	running      string
	onKick       func() error
}

func (f *fakeServiceController) Bootstrap(domain, plist string) error {
	f.boots++
	f.loaded = true
	f.running = serviceFileHash(f.binary)
	return f.bootErr
}
func (f *fakeServiceController) Kickstart(domain, label string) error {
	f.kicks++
	f.loaded = true
	f.running = serviceFileHash(f.binary)
	if f.onKick != nil {
		return f.onKick()
	}
	return nil
}
func (f *fakeServiceController) Bootout(domain, label string) error {
	f.stops++
	f.loaded = false
	return nil
}
func (f *fakeServiceController) Print(domain, label string) (string, error) {
	if !f.loaded {
		return "not loaded", errServiceNotLoaded
	}
	return "pid = 123\nruns = 1\nlast exit code = 0", nil
}
func (f *fakeServiceController) Getenv(name string) (string, error) { return f.key, nil }

type serviceTransport func(*http.Request) (*http.Response, error)

func (f serviceTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type serviceFixture struct {
	env               serviceEnv
	ctl               *fakeServiceController
	args              []string
	home, db, source  string
	serving, mismatch bool
	writes            []string
	active            daemon.HealthActive
}

func newServiceFixture(t *testing.T) *serviceFixture {
	t.Helper()
	f := &serviceFixture{home: t.TempDir(), serving: true}
	f.env = serviceEnv{paths: homeServicePaths(f.home), controller: nil, terminal: sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }, now: func() int64 { return 123 }},
		goos: "darwin", domain: "gui/123", label: worldServiceLabel, wait: 0, portFree: func(string) error { return nil }}
	f.ctl = &fakeServiceController{binary: f.env.paths.Binary}
	f.env.controller = f.ctl
	f.env.probe = func(binary, home string) ([]byte, error) {
		return json.Marshal(serviceVersion{Build: daemon.BuildIdentity{ExecutableSHA256: serviceFileHash(binary)}, Platform: runtime.GOOS + "/" + runtime.GOARCH, ToolRelease: toolRelease(), InterpreterRelease: interpreterRelease})
	}
	f.env.observe = func(p string) { f.writes = append(f.writes, p) }
	f.env.client = &http.Client{Transport: serviceTransport(func(r *http.Request) (*http.Response, error) {
		if !f.serving {
			return nil, errors.New("stopped")
		}
		hash := f.ctl.running
		if f.mismatch {
			hash = "sha256:other"
		}
		b, _ := json.Marshal(daemon.HealthResponse{Status: "ok", Build: daemon.BuildIdentity{ExecutableSHA256: hash}, Active: f.active})
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(b)), Header: make(http.Header)}, nil
	})}
	f.db = filepath.Join(f.home, "data", "world.db")
	root := filepath.Join(f.home, "workspaces")
	interp := filepath.Join(f.home, "interpreter")
	tools := filepath.Join(f.home, "tools")
	for _, p := range []string{filepath.Dir(f.db), root, interp, filepath.Join(tools, toolRelease())} {
		if err := os.MkdirAll(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	pinBytes := []byte("fixture pin")
	old := pinTable
	pinTable = append(pinTable[:0:0], pinTable...)
	for i := range pinTable {
		pinTable[i].BinarySHA256 = strings.TrimPrefix(serviceHash(pinBytes), "sha256:")
	}
	t.Cleanup(func() { pinTable = old })
	for _, r := range pinRoles(interp, tools) {
		if err := os.WriteFile(r.dest, pinBytes, 0700); err != nil {
			t.Fatal(err)
		}
	}
	f.source = filepath.Join(f.home, "candidate")
	os.WriteFile(f.source, []byte("binary fixture"), 0700)
	f.args = []string{"install", "--binary", f.source, "--db", f.db, "--workspace-root", root, "--interpreter-dir", interp, "--tools-dir", tools}
	return f
}
func (f *serviceFixture) run(args ...string) (int, string) {
	var out bytes.Buffer
	code := runServiceWith(args, &out, &out, f.env)
	return code, out.String()
}
func TestHomeServiceFence(t *testing.T) {
	for _, answer := range []string{"no tty", "n\n", "yes\n", "\n"} {
		t.Run(strings.TrimSpace(answer), func(t *testing.T) {
			f := newServiceFixture(t)
			f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
				if answer == "no tty" {
					return nil, os.ErrNotExist
				}
				return &fakeTerm{r: strings.NewReader(answer)}, nil
			}
			code, msg := f.run(f.args...)
			if code != 1 {
				t.Fatalf("%d %s", code, msg)
			}
			if answer == "no tty" && !strings.Contains(msg, "refusing: no controlling terminal") {
				t.Fatal(msg)
			}
			if len(f.writes) != 0 {
				t.Fatal(f.writes)
			}
			for _, p := range []string{f.env.paths.Plist, f.env.paths.Record, f.env.paths.Lock, f.env.paths.Binary} {
				if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("wrote %s", p)
				}
			}
		})
	}
}
func TestHomeServiceInstallAndPlist(t *testing.T) {
	f := newServiceFixture(t)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	if f.ctl.boots != 1 {
		t.Fatal(f.ctl.boots)
	}
	for p, mode := range map[string]os.FileMode{f.env.paths.Plist: 0644, f.env.paths.Record: 0600, f.env.paths.Binary: 0755, f.env.paths.Lock: 0600} {
		st, err := os.Lstat(p)
		if err != nil || !st.Mode().IsRegular() || st.Mode().Perm() != mode {
			t.Fatalf("%s %v %v", p, st, err)
		}
	}
	b, _ := os.ReadFile(f.env.paths.Plist)
	s := string(b)
	for _, want := range []string{"<string>" + f.env.paths.Binary + "</string>", "<string>127.0.0.1:7644</string>", "<key>KeepAlive</key><true/>", "<key>ExitTimeOut</key><integer>30</integer>", "<key>EnvironmentVariables</key><dict><key>PATH</key><string>/usr/bin:/bin:/usr/sbin:/sbin</string></dict>"} {
		if !strings.Contains(s, want) {
			t.Fatal(want)
		}
	}
	if strings.Contains(s, "AILANG_REGISTRY_API_KEY") {
		t.Fatal(s)
	}
	if tool, err := exec.LookPath("plutil"); err == nil {
		if b, err := exec.Command(tool, "-lint", "--", f.env.paths.Plist).CombinedOutput(); err != nil {
			t.Fatalf("%s %v", b, err)
		}
	} else {
		t.Log("plutil unavailable; lint skipped")
	}
	r, err := readServiceRecord(f.env.paths)
	if err != nil || r.Schema != "world/service/v1" || r.BinarySHA256 != serviceFileHash(f.env.paths.Binary) {
		t.Fatalf("%+v %v", r, err)
	}
}
func TestHomeServiceTimeout(t *testing.T) {
	f := newServiceFixture(t)
	f.serving = false
	code, msg := f.run(f.args...)
	if code != 2 || f.ctl.boots != 1 || f.ctl.stops != 1 || !strings.Contains(msg, "files were kept") {
		t.Fatalf("%d %s %+v", code, msg, f.ctl)
	}
	if _, err := os.Stat(f.env.paths.Plist); err != nil {
		t.Fatal(err)
	}
}
func TestHomeServiceStatusStates(t *testing.T) {
	f := newServiceFixture(t)
	if code, msg := f.run("status"); code != 10 || !strings.Contains(msg, "no World service") {
		t.Fatalf("%d %s", code, msg)
	}
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	for _, tc := range []struct {
		name            string
		serve, mismatch bool
		code            int
	}{{"matching", true, false, 0}, {"stopped", false, false, 10}, {"restart pending", true, true, 10}} {
		f.serving, f.mismatch = tc.serve, tc.mismatch
		if code, msg := f.run("status", "--json"); code != tc.code || !strings.Contains(msg, tc.name) && tc.name == "restart pending" {
			t.Fatalf("%s %d %s", tc.name, code, msg)
		}
	}
	f.serving, f.mismatch = true, false
	file, _ := os.OpenFile(f.env.paths.Plist, os.O_APPEND|os.O_WRONLY, 0)
	file.WriteString("\n")
	file.Close()
	if code, msg := f.run("status"); code != 0 || !strings.Contains(msg, "plist differs") {
		t.Fatalf("%d %s", code, msg)
	}
	if code, _ := f.run("status", "--binary", f.source); code != 1 {
		t.Fatal(code)
	}
}
func TestHomeServiceUninstallAndWriteSet(t *testing.T) {
	f := newServiceFixture(t)
	creds := filepath.Join(filepath.Dir(f.db), "credentials.db")
	os.WriteFile(f.db, []byte("db bytes"), 0600)
	os.WriteFile(creds, []byte("credential bytes"), 0600)
	dbHash, credHash := serviceFileHash(f.db), serviceFileHash(creds)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	os.MkdirAll(filepath.Dir(f.env.paths.Log), 0700)
	os.WriteFile(f.env.paths.Log, []byte("keep log"), 0600)
	os.WriteFile(f.source, []byte("write-set upgrade"), 0700)
	if code, msg := f.run("upgrade", "--binary", f.source); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	for _, v := range []string{"status", "stop", "start", "uninstall"} {
		if code, msg := f.run(v); code != 0 {
			t.Fatalf("%s %d %s", v, code, msg)
		}
	}
	if serviceFileHash(f.db) != dbHash || serviceFileHash(creds) != credHash {
		t.Fatal("database changed")
	}
	for _, p := range []string{f.env.paths.Binary, f.env.paths.Log} {
		if _, err := os.Stat(p); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(f.env.paths.Plist); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	r, err := readServiceRecord(f.env.paths)
	if err != nil || !r.Uninstalled {
		t.Fatalf("%+v %v", r, err)
	}
	// Only the named files, binary/log contents and their exact initialization
	// directories are admitted. An arbitrary ancestor (such as HOME) is not.
	allowed := []string{f.env.paths.Plist, f.env.paths.Record, f.env.paths.Lock,
		filepath.Dir(f.env.paths.Plist), filepath.Dir(f.env.paths.Record),
		filepath.Dir(f.env.paths.Binary), filepath.Dir(f.env.paths.Log)}
	for _, p := range f.writes {
		ok := within(filepath.Dir(f.env.paths.Binary), p) || within(filepath.Dir(f.env.paths.Log), p)
		for _, a := range allowed {
			if p == a {
				ok = true
			}
		}
		if !ok {
			t.Fatalf("unexpected write %s", p)
		}
	}
}
func TestHomeServiceRefusals(t *testing.T) {
	for _, kind := range []string{"platform", "label", "key", "port", "pin", "writer", "service lock", "bind", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			f := newServiceFixture(t)
			switch kind {
			case "platform":
				f.env.goos = "linux"
			case "label":
				f.env.label = "dev.ailang.mission-world"
			case "key":
				f.ctl.key = "secret"
			case "port":
				f.env.portFree = func(string) error { return errors.New("busy") }
			case "pin":
				os.Remove(pinRoles(filepath.Join(f.home, "interpreter"), filepath.Join(f.home, "tools"))[0].dest)
			case "writer": // Probe a held lock using the same host API as serve.
				st, err := storeForServiceTest(f.db)
				if err != nil {
					t.Fatal(err)
				}
				defer st()
			case "service lock":
				os.MkdirAll(filepath.Dir(f.env.paths.Lock), 0700)
				lock, err := os.OpenFile(f.env.paths.Lock, os.O_CREATE|os.O_RDWR, 0600)
				if err != nil {
					t.Fatal(err)
				}
				defer lock.Close()
				if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
					t.Fatal(err)
				}
			case "bind":
				f.args = append(f.args, "--bind", "0.0.0.0:7644")
			case "symlink":
				os.MkdirAll(filepath.Dir(f.env.paths.Binary), 0700)
				os.Symlink(f.source, f.env.paths.Binary)
			}
			code, msg := f.run(f.args...)
			if code != 1 || !strings.Contains(msg, "  why:") || !strings.Contains(msg, "  next:") {
				t.Fatalf("%d %s", code, msg)
			}
			if len(f.writes) > 0 {
				t.Fatal(f.writes)
			}
		})
	}
}
func TestHomeServiceEditedStartAndRotation(t *testing.T) {
	f := newServiceFixture(t)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	f.run("stop")
	b, _ := os.ReadFile(f.env.paths.Plist)
	os.WriteFile(f.env.paths.Plist, append(b, '\n'), 0644)
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
	before := f.ctl.boots
	if code, _ := f.run("start"); code != 1 || f.ctl.boots != before {
		t.Fatal(code)
	}
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }
	os.MkdirAll(filepath.Dir(f.env.paths.Log), 0700)
	log, err := os.Create(f.env.paths.Log)
	if err != nil {
		t.Fatal(err)
	}
	log.Truncate((10 << 20) + 1)
	log.Close()
	if code, msg := f.run("start"); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	if _, err := os.Stat(f.env.paths.Log + ".1"); err != nil {
		t.Fatal(err)
	}
}
func TestHomeServiceSeamOnly(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range files {
		if strings.HasSuffix(p, "_test.go") || p == "service_launchctl.go" {
			continue
		}
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(b, []byte("launchctl")) {
			t.Fatalf("supervisor process outside seam: %s", p)
		}
	}
}
func TestHomeVersionJSON(t *testing.T) {
	var b bytes.Buffer
	if code := run([]string{"version", "--json"}, &b, &b); code != 0 {
		t.Fatalf("%d %s", code, b.String())
	}
	var v struct {
		Build                                     daemon.BuildIdentity `json:"build"`
		Platform, ToolRelease, InterpreterRelease string
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"build", "platform", "tool_release", "interpreter_release"} {
		if len(raw[key]) == 0 {
			t.Fatal(key)
		}
	}
	json.Unmarshal(raw["build"], &v.Build)
	for key, want := range map[string]string{"platform": runtime.GOOS + "/" + runtime.GOARCH, "tool_release": toolRelease(), "interpreter_release": interpreterRelease} {
		var got string
		if err := json.Unmarshal(raw[key], &got); err != nil || got != want {
			t.Fatalf("%s = %s, want %s", key, got, want)
		}
	}
	if v.Build.ExecutableSHA256 == "" || v.Build.GoVersion == "" {
		t.Fatal(b.String())
	}
}

func storeForServiceTest(path string) (func(), error) {
	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	return func() { s.Close() }, nil
}

func TestHomeServiceStatusReadOnly(t *testing.T) {
	f := newServiceFixture(t)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	f.writes = nil
	for _, args := range [][]string{{"status"}, {"status", "--json"}} {
		if code, msg := f.run(args...); code != 0 {
			t.Fatalf("%d %s", code, msg)
		}
	}
	if len(f.writes) != 0 {
		t.Fatal(f.writes)
	}
}

func TestHomeServiceEditedNonLoopbackRefused(t *testing.T) {
	f := newServiceFixture(t)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	b, _ := os.ReadFile(f.env.paths.Plist)
	b = bytes.ReplaceAll(b, []byte("127.0.0.1:7644"), []byte("0.0.0.0:7644"))
	os.WriteFile(f.env.paths.Plist, b, 0644)
	f.writes = nil
	before := f.ctl.boots
	if code, msg := f.run("start"); code != 1 || !strings.Contains(msg, "loopback") {
		t.Fatalf("%d %s", code, msg)
	}
	if len(f.writes) != 0 || before != f.ctl.boots {
		t.Fatal("edited non-loopback plist started")
	}
}
