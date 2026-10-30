package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

func pointerHome(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", root)
	return root
}

func pointerWrite(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func pointerRecord(t *testing.T, home, addr string) serviceRecord {
	t.Helper()
	p := homeServicePaths(home)
	pointerWrite(t, p.Binary, []byte("installed fixture"))
	pointerWrite(t, p.Plist, []byte("approved plist"))
	r := serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel, Binary: p.Binary, Plist: p.Plist,
		BinarySHA256: serviceFileHash(p.Binary), PlistSHA256: serviceFileHash(p.Plist), Address: addr,
		Store: filepath.Join(home, "data", "world.db"), Workspace: filepath.Join(home, "ws"), Log: p.Log}
	pointerSave(t, home, r)
	return r
}

func pointerSave(t *testing.T, home string, r serviceRecord) {
	t.Helper()
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	pointerWrite(t, homeServicePaths(home).Record, data)
}

// Route HTTP through a recorder: the exec sandbox cannot bind even loopback.
func pointerHTTP(t *testing.T, handler http.Handler) string {
	t.Helper()
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	http.DefaultTransport = serviceTransport(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w.Result(), nil
	})
	return "http://127.0.0.1:7999"
}

func pointerContains(t *testing.T, out string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func pointerThreePart(t *testing.T, out string) {
	t.Helper()
	if len(strings.Split(strings.TrimSpace(out), "\n")) < 3 {
		t.Errorf("not three parts: %s", out)
	}
	pointerContains(t, out, "why:", "next:")
	for _, bad := range []string{"HTTP 0", "exit status 134"} {
		if strings.Contains(out, bad) {
			t.Errorf("leaked %q: %s", bad, out)
		}
	}
}

func TestR1OneDefaultAddress(t *testing.T) {
	home := pointerHome(t)
	r := pointerRecord(t, home, "127.0.0.1:7999")
	oldTransport, oldFetch := whoamiTransport, fetchDaemonHealth
	t.Cleanup(func() { whoamiTransport, fetchDaemonHealth = oldTransport, oldFetch })
	t.Setenv("WORLD_SESSION", strings.Repeat("a", 64))
	for _, tc := range []struct {
		name, want string
		record     *serviceRecord
	}{
		{"absent", "127.0.0.1:7644", nil},
		{"installed", r.Address, &r},
		{"uninstalled", "127.0.0.1:7644", func() *serviceRecord { x := r; x.Uninstalled = true; return &x }()},
		{"other home", "127.0.0.1:7644", func() *serviceRecord { x := r; x.Binary = filepath.Join(home, "other", "ailang-worldd"); return &x }()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.record == nil {
				_ = os.Remove(homeServicePaths(home).Record)
			} else {
				pointerSave(t, home, *tc.record)
			}
			if defaultServiceBind() != tc.want || defaultServiceAddr() != "http://"+tc.want || defaultEpisodeAddr != daemon.DefaultBind {
				t.Fatal("default helpers disagree")
			}
			for _, explicit := range []bool{false, true} {
				want := tc.want
				if explicit {
					want = "127.0.0.1:7888"
				}
				var b bytes.Buffer
				seen := ""
				whoamiTransport = serviceTransport(func(req *http.Request) (*http.Response, error) {
					seen = req.URL.Host
					return nil, errors.New("test offline")
				})
				args := []string{"whoami"}
				if explicit {
					args = append([]string{"--addr", "http://" + want}, args...)
				}
				if code := run(args, &b, &b); code != 10 || seen != want {
					t.Fatalf("global default: code %d host %s want %s; %s", code, seen, want, b.String())
				}
				dir := filepath.Join(home, tc.name, map[bool]string{false: "default", true: "explicit"}[explicit])
				hostArgs := []string{}
				if explicit {
					hostArgs = []string{"--addr", "http://" + want}
				}
				if code, out := hostsInitTest(t, dir, hostArgs...); code != 0 {
					t.Fatal(out)
				}
				data, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
				if err != nil {
					t.Fatal(err)
				}
				var manifest hostsManifest
				if json.Unmarshal(data, &manifest) != nil || manifest.Addr != "http://"+want {
					t.Fatalf("hosts default: %s", data)
				}
				// Observe episode's actual parsed flag through its health seam.
				er := newEpisodeRun(t)
				// newEpisodeRun isolates HOME; restore this case's record.
				t.Setenv("HOME", home)
				fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) {
					seen = addr
					return daemon.HealthResponse{}, &daemonUnreachableError{err: errors.New("test offline")}
				}
				extra := []string{}
				if explicit {
					extra = []string{"--addr", want}
				}
				seen = ""
				b.Reset()
				opened := false
				code := runEpisodeStart(er.args("default-probe", extra...), &b, &b, spyEnv(&opened))
				if code != 1 || seen != want || opened {
					t.Fatalf("episode default: %d %s want %s opened %v: %s", code, seen, want, opened, b.String())
				}
			}
		})
	}
}

type pointerController struct {
	fakeServiceController
	prints int
}

func (c *pointerController) Print(domain, label string) (string, error) {
	c.prints++
	if domain != "gui/test" || label != worldServiceLabel {
		return "", errors.New("wrong supervisor identity")
	}
	return c.fakeServiceController.Print(domain, label)
}

func pointerDoctorEnv(t *testing.T, home string, c *pointerController) {
	t.Helper()
	old := doctorServiceEnv
	t.Cleanup(func() { doctorServiceEnv = old })
	doctorServiceEnv = func() serviceEnv {
		return serviceEnv{paths: homeServicePaths(home), controller: c, domain: "gui/test", label: worldServiceLabel}
	}
}

func TestDoctorQuickServiceFix(t *testing.T) {
	home := pointerHome(t)
	r := pointerRecord(t, home, "127.0.0.1:7999")
	oldTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = oldTransport })
	http.DefaultTransport = serviceTransport(func(*http.Request) (*http.Response, error) { return nil, errors.New("closed test endpoint") })
	c := &pointerController{}
	pointerDoctorEnv(t, home, c)
	for _, plist := range []bool{true, false} {
		if !plist {
			if err := os.Remove(r.Plist); err != nil {
				t.Fatal(err)
			}
		}
		var b bytes.Buffer
		if code := runDoctor("http://"+r.Address, []string{"--quick", "--json"}, &b, &b); code != 10 {
			t.Fatalf("exit %d: %s", code, b.String())
		}
		fix := "ailang-worldd service start"
		if !plist {
			fix = "ailang-worldd service install"
		}
		pointerContains(t, b.String(), "\"check\":\"daemon\"", "\"mark\":\"fail\"", fix)
		for _, line := range strings.Split(strings.TrimSpace(b.String()), "\n") {
			var row map[string]any
			if err := json.Unmarshal([]byte(line), &row); err != nil {
				t.Fatal(err)
			}
			if name, ok := row["check"]; ok && name != "daemon" && name != "session" && name != "host" {
				t.Fatalf("quick added check %s", name)
			}
		}
		if c.prints != 0 {
			t.Fatal("quick queried supervisor")
		}
	}
}

func TestDoctorServiceCheck(t *testing.T) {
	home := pointerHome(t)
	c := &pointerController{}
	pointerDoctorEnv(t, home, c)
	var b bytes.Buffer
	report := &doctorReport{w: &b}
	doctorService(report, "http://127.0.0.1:7998")
	pointerContains(t, b.String(), "! service", "no World service is installed", "ailang-worldd service install")
	if c.prints != 0 {
		t.Fatal("no record called supervisor")
	}
	r := pointerRecord(t, home, "127.0.0.1:7999")
	installed := serviceFileHash(r.Binary)
	for _, tc := range []struct {
		name, addr, build, want string
		loaded, edited          bool
		calls, code             int
	}{
		{"other address", "http://127.0.0.1:7998", "", "✓ service", false, false, 0, 0},
		{"not loaded", "http://" + r.Address, "", "ailang-worldd service start", false, false, 1, 10},
		{"loaded", "http://" + r.Address, installed, "pid 123", true, false, 1, 0},
		{"pending", "http://" + r.Address, "sha256:other", "restart pending", true, false, 1, 0},
		{"edited", "http://" + r.Address, installed, "the plist differs from the one installed", true, true, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b.Reset()
			c.prints, c.loaded = 0, tc.loaded
			if tc.edited {
				pointerWrite(t, r.Plist, []byte("edited"))
			}
			report := &doctorReport{w: &b, serviceBuild: daemon.BuildIdentity{ExecutableSHA256: tc.build}}
			doctorService(report, tc.addr)
			pointerContains(t, b.String(), tc.want, "this ailang-worldd (", "is not the installed one", "service upgrade")
			if c.prints != tc.calls || report.classCode != tc.code {
				t.Fatalf("prints %d code %d: %s", c.prints, report.classCode, b.String())
			}
		})
	}
	// Drive the full pass: it reuses doctorDaemon's one health answer.
	c.prints, c.loaded = 0, true
	addr := pointerHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"daemon_version": "test", "db_path": r.Store, "build": daemon.BuildIdentity{ExecutableSHA256: installed}})
	}))
	r.Address = strings.TrimPrefix(addr, "http://")
	pointerSave(t, home, r)
	b.Reset()
	_ = runDoctor(addr, []string{"--db", r.Store}, &b, &b)
	if c.prints != 1 {
		t.Fatalf("full pass: %d supervisor calls: %s", c.prints, b.String())
	}
	class, code := doctorClass("service")
	if class != "daemon" || code != 10 {
		t.Fatal(class, code)
	}
}

func TestWhoamiServiceMessages(t *testing.T) {
	home := pointerHome(t)
	t.Setenv("WORLD_SESSION", strings.Repeat("a", 64))
	old := whoamiTransport
	t.Cleanup(func() { whoamiTransport = old })
	whoamiTransport = serviceTransport(func(*http.Request) (*http.Response, error) { return nil, errors.New("offline") })
	addr := "http://127.0.0.1:7999"
	check := func(wants ...string) {
		t.Helper()
		var b bytes.Buffer
		if code := runWhoami(addr, nil, &b, &b); code != 10 {
			t.Fatal(code, b.String())
		}
		pointerContains(t, b.String(), wants...)
		pointerThreePart(t, b.String())
	}
	check("what: daemon unavailable", "nothing answered at "+addr+", and no World service is installed (no ~/.ailang-world/service.json)", "next: ailang-worldd service install (or serve by hand for a scratch store: QUICKSTART §1)")
	r := pointerRecord(t, home, "127.0.0.1:7999")
	check("why: the World service is not serving at "+addr, "next: ailang-worldd service start (ailang-worldd service status shows why it stopped)")
	r.Address = "127.0.0.1:7998"
	pointerSave(t, home, r)
	check("why: nothing answered at "+addr+" (the World service is at "+r.Address+")", "next: pass --addr http://"+r.Address+", or start the daemon at "+addr)
}

func TestEpisodeStartServicePointers(t *testing.T) {
	for _, scenario := range []string{"down", "cli differs", "restart pending", "other address"} {
		t.Run(scenario, func(t *testing.T) {
			er := newEpisodeRun(t)
			home := pointerHome(t)
			r := pointerRecord(t, home, "127.0.0.1:7999")
			oldFetch, oldSelf := fetchDaemonHealth, episodeSelfBuild
			t.Cleanup(func() { fetchDaemonHealth, episodeSelfBuild = oldFetch, oldSelf })
			self := daemon.BuildIdentity{ExecutableSHA256: "sha256:abcdef0123456789"}
			if scenario == "restart pending" {
				self.ExecutableSHA256 = serviceFileHash(r.Binary)
			}
			episodeSelfBuild = func() daemon.BuildIdentity { return self }
			fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
				if scenario == "down" || scenario == "other address" {
					return daemon.HealthResponse{}, &daemonUnreachableError{err: errors.New("offline")}
				}
				return daemon.HealthResponse{Build: daemon.BuildIdentity{ExecutableSHA256: "sha256:123456789abc"}}, nil
			}
			if scenario == "other address" {
				r.Address = "127.0.0.1:7998"
				pointerSave(t, home, r)
			}
			var b bytes.Buffer
			opened := false
			code := runEpisodeStart(er.args("service-probe", "--addr", "127.0.0.1:7999"), &b, &b, spyEnv(&opened))
			if code != 1 || opened {
				t.Fatalf("exit %d opened %v: %s", code, opened, b.String())
			}
			switch scenario {
			case "down":
				pointerContains(t, b.String(), "refusing: the World service is not serving at 127.0.0.1:7999", "ailang-worldd service start")
				if strings.Contains(b.String(), " serve ") {
					t.Fatal(b.String())
				}
			case "cli differs":
				pointerContains(t, b.String(), "refusing: this ailang-worldd (", "sha abcdef0", "is not the installed one", r.Binary+" episode start …", "service upgrade")
			case "restart pending":
				pointerContains(t, b.String(), "running build 1234567", "(restart pending)", "service stop, then ailang-worldd service start")
			case "other address":
				pointerContains(t, b.String(), "no daemon is serving on", " serve ")
			}
			if scenario != "other address" {
				pointerThreePart(t, b.String())
			}
		})
	}
}

func pointerHostSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			out[rel] = string(b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestHostsInitUsesTheInstalledService(t *testing.T) {
	home := pointerHome(t)
	r := pointerRecord(t, home, "127.0.0.1:7999")
	dir := filepath.Join(home, "hosts")
	code, out := hostsInitTest(t, dir)
	if code != 0 {
		t.Fatal(code, out)
	}
	exe, _ := os.Executable()
	pointerContains(t, out, "! the scripts call the installed World binary "+r.Binary+", not this one ("+exe+")")
	files := pointerHostSnapshot(t, dir)
	for _, name := range []string{"bin/world-context", "bin/world-statusline", "pi/world-status.ts"} {
		pointerContains(t, files[name], r.Binary)
	}
	for _, name := range []string{"claude-world.mcp.json", "codex/config.toml", "pi-world-mcp.json"} {
		pointerContains(t, files[name], "http://"+r.Address+"/mcp/")
	}
	for name, b := range files {
		if strings.Contains(b, exe) {
			t.Errorf("%s calls current executable", name)
		}
	}
}

func TestHostsInitStableAcrossUpgrade(t *testing.T) {
	home := pointerHome(t)
	r := pointerRecord(t, home, "127.0.0.1:7999")
	dir := filepath.Join(home, "hosts")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatal(code, out)
	}
	before := pointerHostSnapshot(t, dir)
	pointerWrite(t, r.Binary, []byte("new installed bytes after an upgrade"))
	r.BinarySHA256 = serviceFileHash(r.Binary)
	pointerSave(t, home, r)
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatal(code, out)
	}
	if after := pointerHostSnapshot(t, dir); !reflect.DeepEqual(before, after) {
		t.Fatal("upgrade changed generated bytes")
	}
}

func TestWhyExpectStore(t *testing.T) {
	root := pointerHome(t)
	a, b := filepath.Join(root, "a.db"), filepath.Join(root, "b.db")
	pointerWrite(t, a, nil)
	alias := filepath.Join(root, "alias.db")
	if err := os.Symlink(a, alias); err != nil {
		t.Fatal(err)
	}
	health, logs := 0, 0
	allowLog := false
	addr := pointerHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/health":
			health++
			_ = json.NewEncoder(w).Encode(map[string]string{"db_path": a})
		default:
			logs++
			if !allowLog {
				t.Error("store mismatch read the log", req.URL.Path)
			}
			http.NotFound(w, req)
		}
	}))
	var out, errOut bytes.Buffer
	code := runWhy(addr, []string{"3", "--expect-store", b}, strings.NewReader(""), &out, &errOut)
	if code != 1 || out.Len() != 0 || logs != 0 || health != 1 {
		t.Fatalf("code %d health %d logs %d out %s err %s", code, health, logs, out.String(), errOut.String())
	}
	want := "ailang-worldd why: the daemon at " + addr + " serves store " + a + ", not " + b + "\n" +
		"  why: entry indexes are per store, so the same index in another store is another entry (a World-Provenance trailer names its store as store=<path>)\n" +
		"  next: pass the --addr of the daemon serving " + b + " (the frozen rig: start it from daemon-cmdline.txt on 127.0.0.1:47330), or drop --expect-store\n"
	if errOut.String() != want {
		t.Fatalf("got %s want %s", errOut.String(), want)
	}
	pointerThreePart(t, errOut.String())
	allowLog = true
	for _, path := range []string{a, alias} {
		out.Reset()
		errOut.Reset()
		before := logs
		_ = runWhy(addr, []string{"3", "--expect-store", path}, strings.NewReader(""), &out, &errOut)
		if logs == before || strings.Contains(errOut.String(), "not "+path) {
			t.Fatalf("matching store did not reach locate: %s", errOut.String())
		}
	}
}

func TestServicePointerMessagesAreThreePart(t *testing.T) {
	// These drivers check the complete messages and exit codes at their real
	// command paths, with every remote effect stubbed.
	t.Run("whoami F2 F3", TestWhoamiServiceMessages)
	t.Run("episode F3 F14", TestEpisodeStartServicePointers)
	t.Run("why F17", TestWhyExpectStore)
}
