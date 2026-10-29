package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

type upgradeTerm struct {
	r io.Reader
	bytes.Buffer
}

func (t *upgradeTerm) Read(p []byte) (int, error) { return t.r.Read(p) }
func (t *upgradeTerm) Close() error               { return nil }

func upgradeFixture(t *testing.T) (*serviceFixture, string, *upgradeTerm) {
	t.Helper()
	f := newServiceFixture(t)
	if code, msg := f.run(f.args...); code != 0 {
		t.Fatalf("install: %d %s", code, msg)
	}
	f.writes = nil
	candidate := filepath.Join(t.TempDir(), "candidate")
	if err := os.WriteFile(candidate, []byte("upgrade candidate"), 0755); err != nil {
		t.Fatal(err)
	}
	term := &upgradeTerm{r: strings.NewReader("y\n")}
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
	return f, candidate, term
}

func TestHomeServiceUpgradeSuccess(t *testing.T) {
	f, candidate, term := upgradeFixture(t)
	old := serviceFileHash(f.env.paths.Binary)
	before, _ := os.Stat(f.env.paths.Binary)
	code, msg := f.run("upgrade", "--binary", candidate)
	if code != 0 || !strings.Contains(msg, "upgraded "+serviceShort(old)+" → "+serviceShort(serviceFileHash(candidate))) {
		t.Fatalf("%d %s", code, msg)
	}
	after, _ := os.Stat(f.env.paths.Binary)
	if os.SameFile(before, after) {
		t.Fatal("binary inode was overwritten")
	}
	to := serviceFileHash(candidate)
	if serviceFileHash(f.env.paths.Binary) != to || serviceFileHash(f.env.paths.Binary+".prev") != old || f.ctl.running != to {
		t.Fatal("binary, previous build or running build differs")
	}
	r, err := readServiceRecord(f.env.paths)
	if err != nil || r.BinarySHA256 != to || r.PreviousSHA256 != old || len(r.Upgrades) != 1 ||
		r.Upgrades[0] != (serviceUpgrade{From: old, To: to, At: 123, Outcome: "ok"}) {
		t.Fatalf("%+v %v", r, err)
	}
	if f.ctl.kicks != 1 || f.ctl.stops != 0 {
		t.Fatalf("kicks=%d stops=%d", f.ctl.kicks, f.ctl.stops)
	}
	for _, want := range []string{"from " + serviceShort(old) + " → to " + serviceShort(to), "tool AILANG: " + toolRelease() + " (unchanged)", "interpreter: " + interpreterRelease + " (unchanged)", "plist: unchanged", "10 s drain plus a restart", "Proceed? [y/N]"} {
		if !strings.Contains(term.String(), want) {
			t.Fatalf("prompt missing %q: %s", want, term.String())
		}
	}
	for _, p := range f.writes {
		if p != filepath.Dir(f.env.paths.Record) && p != f.env.paths.Record && p != f.env.paths.Lock &&
			p != filepath.Dir(f.env.paths.Binary) && !within(filepath.Dir(f.env.paths.Binary), p) &&
			p != filepath.Dir(f.env.paths.Log) && !within(filepath.Dir(f.env.paths.Log), p) {
			t.Fatalf("upgrade wrote %s", p)
		}
		if p == f.env.paths.Plist || p == filepath.Dir(f.env.paths.Plist) {
			t.Fatalf("upgrade touched plist: %s", p)
		}
	}
	t.Logf("prompt:\n%s\nresult:\n%s", term.String(), msg)
	_, status := f.run("status")
	t.Logf("status:\n%s", status)
}

func TestHomeServiceUpgradeActive(t *testing.T) {
	for _, inFlight := range []int{0, 1} {
		t.Run(fmt.Sprint(inFlight), func(t *testing.T) {
			f, candidate, term := upgradeFixture(t)
			f.active = daemon.HealthActive{InFlight: inFlight, Sessions: []daemon.HealthActiveSession{{Episode: "w223e", LastCallUnix: 83}}}
			old := serviceFileHash(f.env.paths.Binary)
			code, msg := f.run("upgrade", "--binary", candidate)
			if code != 1 || !strings.Contains(msg, "w223e (last call 40 s ago)") || !strings.Contains(msg, "service upgrade --pause") ||
				len(f.writes) != 0 || serviceFileHash(f.env.paths.Binary) != old || f.ctl.kicks != 0 {
				t.Fatalf("%d %s writes=%v", code, msg, f.writes)
			}
			t.Log(msg)
			if code, msg := f.run("upgrade", "--binary", candidate, "--pause"); code != 0 {
				t.Fatalf("%d %s", code, msg)
			}
			if !strings.Contains(term.String(), "--pause will interrupt: 1 episodes are working: w223e (last call 40 s ago)") {
				t.Fatal(term.String())
			}
		})
	}
	t.Run("in flight without recent session", func(t *testing.T) {
		f, candidate, _ := upgradeFixture(t)
		f.active.InFlight = 1
		if code, msg := f.run("upgrade", "--binary", candidate); code != 1 || !strings.Contains(msg, "1 calls in flight") || len(f.writes) != 0 {
			t.Fatalf("%d %s %v", code, msg, f.writes)
		}
	})
}

func TestHomeServiceUpgradeRollback(t *testing.T) {
	for _, failsOld := range []bool{false, true} {
		t.Run(fmt.Sprint(failsOld), func(t *testing.T) {
			f, candidate, _ := upgradeFixture(t)
			old := serviceFileHash(f.env.paths.Binary)
			os.MkdirAll(filepath.Dir(f.env.paths.Log), 0700)
			os.WriteFile(f.env.paths.Log, []byte("candidate refused to serve\n"), 0600)
			f.ctl.onKick = func() error {
				if f.ctl.kicks == 1 || failsOld {
					f.ctl.running = "sha256:wrong"
				}
				return nil
			}
			code, msg := f.run("upgrade", "--binary", candidate)
			want, outcome := 2, "rolled back"
			if failsOld {
				want, outcome = 3, "rollback failed"
			}
			if code != want || f.ctl.kicks != 2 || serviceFileHash(f.env.paths.Binary) != old ||
				serviceFileHash(f.env.paths.Binary+".failed") != serviceFileHash(candidate) {
				t.Fatalf("%d %s kicks=%d", code, msg, f.ctl.kicks)
			}
			r, err := readServiceRecord(f.env.paths)
			if err != nil || r.BinarySHA256 != old || len(r.Upgrades) != 1 || r.Upgrades[0].Outcome != outcome {
				t.Fatalf("%+v %v", r, err)
			}
			if failsOld {
				if f.ctl.stops != 1 || f.ctl.loaded || !strings.Contains(msg, f.env.paths.Binary+" (old)") ||
					!strings.Contains(msg, f.env.paths.Binary+".failed (new)") || !strings.Contains(msg, "next: "+f.env.paths.Binary+" service start") {
					t.Fatal(msg)
				}
			} else if f.ctl.stops != 0 || f.ctl.running != old || !strings.Contains(msg, "which is serving") || !strings.Contains(msg, "candidate refused to serve") {
				t.Fatal(msg)
			}
			for _, want := range []string{"  why:", "  next:"} {
				if !strings.Contains(msg, want) {
					t.Fatal(msg)
				}
			}
			t.Log(msg)
		})
	}
}

func TestHomeServiceUpgradeRefusals(t *testing.T) {
	for _, kind := range []string{"missing", "directory", "symlink", "platform", "tool", "interpreter", "bad JSON", "trailing JSON", "probe error", "not installed", "uninstalled", "plist", "file sha", "health sha", "health down", "already", "no terminal", "answer", "lock", "flags"} {
		t.Run(kind, func(t *testing.T) {
			f, candidate, term := upgradeFixture(t)
			want, contains := 1, "  why:"
			switch kind {
			case "missing":
				os.Remove(candidate)
				contains = candidate
			case "directory":
				os.Remove(candidate)
				os.Mkdir(candidate, 0700)
			case "symlink":
				os.Remove(candidate)
				os.Symlink(f.source, candidate)
			case "platform", "tool", "interpreter":
				probe := f.env.probe
				f.env.probe = func(binary, home string) ([]byte, error) {
					b, err := probe(binary, home)
					if binary == candidate {
						var v serviceVersion
						json.Unmarshal(b, &v)
						switch kind {
						case "platform":
							v.Platform = "wrong/platform"
						case "tool":
							v.ToolRelease = "v999"
						case "interpreter":
							v.InterpreterRelease = "v999"
						}
						b, _ = json.Marshal(v)
					}
					return b, err
				}
				contains = map[string]string{"platform": "is not an ailang-worldd for", "tool": "pins tool AILANG v999", "interpreter": "world-publish epoch nominate"}[kind]
			case "bad JSON":
				f.env.probe = func(string, string) ([]byte, error) { return []byte("{}"), nil }
			case "trailing JSON":
				probe := f.env.probe
				f.env.probe = func(b, h string) ([]byte, error) {
					data, err := probe(b, h)
					return append(data, []byte("{}")...), err
				}
			case "probe error":
				f.env.probe = func(string, string) ([]byte, error) { return nil, errors.New("timeout or nonzero exit") }
			case "not installed":
				os.Remove(f.env.paths.Record)
				contains = "no World service is installed"
			case "uninstalled":
				r, _ := readServiceRecord(f.env.paths)
				r.Uninstalled = true
				b, _ := json.Marshal(r)
				os.WriteFile(f.env.paths.Record, b, 0600)
				contains = "no World service is installed"
			case "plist":
				os.WriteFile(f.env.paths.Plist, []byte("edited"), 0644)
				contains = "service start shows the diff"
			case "file sha":
				os.WriteFile(f.env.paths.Binary, []byte("drift"), 0755)
				contains = "installed file sha256 differs from service.json"
			case "health sha":
				f.mismatch = true
				contains = "restart pending"
			case "health down":
				f.serving = false
				contains = "does not answer /v1/health"
			case "already":
				candidate = f.source
				want, contains = 0, "already installed (sha256:"
			case "no terminal":
				f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return nil, errors.New("no tty") }
				contains = "refusing: no controlling terminal"
			case "answer":
				term.r = strings.NewReader("yes\n")
				contains = "aborted"
			case "lock":
				lock, err := f.env.lock()
				if err != nil {
					t.Fatal(err)
				}
				defer lock.Close()
				f.writes = nil
				contains = "a service change is in progress"
			case "flags":
				code, msg := f.run("upgrade", "--replace")
				if code != 1 || len(f.writes) != 0 {
					t.Fatalf("%d %s", code, msg)
				}
				return
			}
			before := serviceFileHash(f.env.paths.Binary)
			code, msg := f.run("upgrade", "--binary", candidate)
			if code != want || !strings.Contains(msg, contains) || len(f.writes) != 0 ||
				serviceFileHash(f.env.paths.Binary) != before || f.ctl.kicks != 0 {
				t.Fatalf("%d %s writes=%v", code, msg, f.writes)
			}
			if kind != "already" && (!strings.Contains(msg, "  why:") || !strings.Contains(msg, "  next:")) {
				t.Fatal(msg)
			}
		})
	}
}

func TestHomeServiceUpgradeProductionProbe(t *testing.T) {
	f := newServiceFixture(t)
	f.env.probe = nil
	dump := filepath.Join(t.TempDir(), "environment")
	candidate := filepath.Join(t.TempDir(), "candidate")
	// Quote the temp path as shell data; only this fixture script invokes env.
	quoted := "'" + strings.ReplaceAll(dump, "'", "'\\''") + "'"
	script := "#!/bin/sh\n[ \"$1\" = version ] && [ \"$2\" = --json ] && [ \"$#\" = 2 ] || exit 9\n/usr/bin/env > " + quoted + "\nprintf '%s\\n' '{\"platform\":\"wrong/platform\"}'\n"
	if err := os.WriteFile(candidate, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("UPGRADE_TEST_AMBIENT", "must not reach candidate")
	code, msg := f.run("upgrade", "--binary", candidate)
	if code != 1 || !strings.Contains(msg, candidate+" is not an ailang-worldd for "+runtime.GOOS+"/"+runtime.GOARCH) || len(f.writes) != 0 {
		t.Fatalf("%d %s", code, msg)
	}
	b, err := os.ReadFile(dump)
	if err != nil {
		t.Fatal(err)
	}
	env := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		p := strings.SplitN(line, "=", 2)
		if len(p) != 2 {
			t.Fatal(line)
		}
		env[p[0]] = p[1]
		switch p[0] {
		case "PATH", "HOME", "PWD", "SHLVL", "_":
		default:
			t.Fatalf("ambient variable reached candidate: %s", p[0])
		}
	}
	if env["HOME"] != f.home || env["PATH"] != "/usr/bin:/bin:/usr/sbin:/sbin" {
		t.Fatal(env)
	}
}

func TestHomeServiceUpgradeHistory(t *testing.T) {
	f, candidate, _ := upgradeFixture(t)
	for _, omit := range []bool{false, true} {
		b, _ := os.ReadFile(f.env.paths.Record)
		var record map[string]any
		json.Unmarshal(b, &record)
		if omit {
			delete(record, "upgrades")
		} else {
			record["upgrades"] = []string{}
		}
		b, _ = json.Marshal(record)
		os.WriteFile(f.env.paths.Record, b, 0600)
		r, err := readServiceRecord(f.env.paths)
		if err != nil || len(r.Upgrades) != 0 {
			t.Fatalf("%+v %v", r, err)
		}
	}
	var firstTo string
	for i := 0; i < 11; i++ {
		os.WriteFile(candidate, []byte(fmt.Sprintf("build %d", i)), 0755)
		if i == 0 {
			firstTo = serviceFileHash(candidate)
		}
		f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
			return &fakeTerm{r: strings.NewReader("y\n")}, nil
		}
		if code, msg := f.run("upgrade", "--binary", candidate); code != 0 {
			t.Fatalf("%d %s", code, msg)
		}
	}
	r, err := readServiceRecord(f.env.paths)
	if err != nil || len(r.Upgrades) != 10 || r.Upgrades[0].From != firstTo || r.Upgrades[9].To != serviceFileHash(candidate) {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestHomeServiceUpgradeInstallMissingDirectory(t *testing.T) {
	f := newServiceFixture(t)
	root := filepath.Join(f.home, "workspaces")
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	code, msg := f.run(f.args...)
	if code != 1 || !strings.Contains(msg, "service install: "+root+" does not exist") ||
		!strings.Contains(msg, "setup creates both") || !strings.Contains(msg, "next: ailang-worldd setup, then ailang-worldd service install") || len(f.writes) != 0 {
		t.Fatalf("%d %s writes=%v", code, msg, f.writes)
	}
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("missing root was created")
	}
}

func TestHomeServiceUpgradeStatusVerbose(t *testing.T) {
	f, _, _ := upgradeFixture(t)
	for _, serving := range []bool{true, false} {
		f.serving = serving
		for _, mode := range []string{"human", "--verbose", "--json"} {
			args := []string{"status"}
			if mode != "human" {
				args = append(args, mode)
			}
			code, msg := f.run(args...)
			want := 0
			if !serving {
				want = 10
			}
			if code != want || strings.Contains(msg, "pid = 123") != (mode != "human") {
				t.Fatalf("%s: %d %s", mode, code, msg)
			}
			if mode == "--json" && !strings.Contains(msg, "\"supervisor_state\"") {
				t.Fatal(msg)
			}
		}
	}
	if len(f.writes) != 0 {
		t.Fatal(f.writes)
	}
}
