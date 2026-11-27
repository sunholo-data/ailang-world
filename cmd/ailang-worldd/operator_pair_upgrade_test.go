package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

const wantStoppedRunnerRemedy = "The operator runner companion remains stopped. After resolving this, run ailang-worldd service runner start."

func TestOperatorProcessWitness(t *testing.T) {
	mode := os.Getenv("WORLD_TEST_OPERATOR_WITNESS")
	if mode == "" {
		return
	}
	wait := func() error {
		fmt.Println("ready")
		_, err := io.Copy(io.Discard, os.Stdin)
		return err
	}
	path := os.Getenv("WORLD_TEST_OPERATOR_WITNESS_PATH")
	switch mode {
	case "git-inherited":
		lock := os.NewFile(4, "Git inherited runner lock")
		defer lock.Close()
		if got, err := lock.Stat(); err != nil {
			t.Fatal(err)
		} else if want, err := os.Stat(path); err != nil || !os.SameFile(got, want) {
			t.Fatalf("Git received a different runner lock: %v", err)
		}
		if err := syscall.Flock(4, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path+".git-observed", nil, 0600); err != nil {
			t.Fatal(err)
		}
		return
	case "inherited":
		lock := os.NewFile(3, "inherited runner lock")
		defer lock.Close()
		if _, err := lock.Stat(); err != nil {
			t.Fatal(err)
		}
	case "old-terminal":
		// Pre-exclusion writable terminal entry. Hold the real writer
		// without acquiring either operator lock.
		if err := runnerWrite(path, func(*store.Credentials) error { return wait() }); err != nil {
			t.Fatal(err)
		}
		return
	case "episode":
		lock, err := terminalEpisodeLock(path, "ep", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Close()
	case "owner":
		lock, err := store.OpenOperatorEpisodeLock(path, "ep")
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Close()
	default:
		t.Fatalf("bad witness mode %q", mode)
	}
	if err := wait(); err != nil {
		t.Fatal(err)
	}
}

func operatorProcessWitness(t *testing.T, mode, path string, files ...*os.File) (int, func()) {
	t.Helper()
	ctx, cancel := budgetContext(8 * time.Second)
	bin, err := os.Executable()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorProcessWitness$")
	cmd.WaitDelay = time.Second
	cmd.ExtraFiles = files
	cmd.Env = append(os.Environ(), "WORLD_TEST_OPERATOR_WITNESS="+mode, "WORLD_TEST_OPERATOR_WITNESS_PATH="+path,
		"GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	var stderr syncBuffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		cancel()
		t.Fatal(err)
	}
	waited := false
	t.Cleanup(func() {
		cancel()
		if !waited {
			_ = cmd.Wait()
		}
	})
	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	if err != nil || line != "ready\n" {
		t.Fatalf("process barrier: %q %v (%s)", line, err, &stderr)
	}
	release := func() {
		t.Helper()
		if !waited {
			stdin.Close()
			io.Copy(io.Discard, reader)
			err := cmd.Wait()
			waited = true
			if err != nil {
				t.Fatalf("witness exit: %v: %s", err, &stderr)
			}
		}
	}
	return cmd.Process.Pid, release
}

// Pair fixture uses fake supervisor, HTTP transport and install roots from
// operatorService, all under t.TempDir().
type pairedUpgradeSupervisor struct {
	serviceController
	f            *operatorServiceFixture
	runnerLoaded bool
	bootFailure  bool
	stopFault    string
}

func (s *pairedUpgradeSupervisor) Print(domain, label string) (string, error) {
	if label == runnerLabel && !s.runnerLoaded {
		return "", errServiceNotLoaded
	}
	return s.serviceController.Print(domain, label)
}

func (s *pairedUpgradeSupervisor) Bootout(domain, label string) error {
	if label == runnerLabel {
		switch s.stopFault {
		case "error":
			return errors.New("fixture stop failed")
		case "reported-only":
			return nil
		}
		s.runnerLoaded = false
		return nil
	}
	return s.serviceController.Bootout(domain, label)
}

func (s *pairedUpgradeSupervisor) Bootstrap(domain, plist string) error {
	if plist != s.f.env.runnerPaths().Plist {
		return s.serviceController.Bootstrap(domain, plist)
	}
	c, _, err := s.f.env.installedRunner()
	if err != nil {
		return err
	}
	process, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
	if err != nil {
		return fmt.Errorf("companion restart could not acquire its process lock: %w", err)
	}
	process.Close()
	s.runnerLoaded = true
	if s.bootFailure {
		s.bootFailure = false
		return errors.New("refusing first companion boot")
	}
	s.f.heartbeat = store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "paired-" + serviceShort(c.Executables[s.f.env.paths.Binary])},
		Protocol: store.OperatorRunnerProtocol, Build: c.Executables[s.f.env.paths.Binary]}
	return s.f.env.runnerTick(c, &s.f.heartbeat)
}

func pairedUpgradeFixture(t *testing.T) (*operatorServiceFixture, *pairedUpgradeSupervisor, string) {
	t.Helper()
	f := operatorService(t)
	f.mustInstall(t)
	f.homeJob.runner = true
	f.homeJob.onKick = func() error {
		probe, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
		if probe != nil {
			probe.Close()
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("daemon replacement lost its runner drain guard: %v", err)
		}
		f.health.Build.ExecutableSHA256 = f.homeJob.running
		return nil
	}
	candidate := f.source
	if err := os.WriteFile(candidate, []byte("candidate fixture"), 0700); err != nil {
		t.Fatal(err)
	}
	s := &pairedUpgradeSupervisor{serviceController: f.homeJob, f: f, runnerLoaded: true}
	f.env.controller = s
	f.env.runnerProcesses = func() ([]worldProcess, error) {
		return []worldProcess{{os.Getpid(), "ailang-worldd"}}, nil
	}
	return f, s, candidate
}

func pairedUpgradeEvidence(f *operatorServiceFixture) [7]string {
	p := f.env.runnerPaths()
	return [7]string{serviceFileHash(f.env.paths.Binary), serviceFileHash(f.env.paths.Record), serviceFileHash(f.env.paths.Plist),
		serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials), serviceFileHash(f.config.World)}
}

func TestOperatorOldBuildTerminalDrain(t *testing.T) {
	f, s, candidate := pairedUpgradeFixture(t)
	pid, release := operatorProcessWitness(t, "old-terminal", f.config.Credentials)
	read := f.env.runnerProcesses
	f.env.runnerProcesses = func() ([]worldProcess, error) {
		rows, _ := read()
		if syscall.Kill(pid, 0) == nil {
			rows = append(rows, worldProcess{pid, "ailang-worldd.prev"})
		}
		return rows, nil
	}
	before := pairedUpgradeEvidence(f)
	record, err := readServiceRecord(f.env.paths)
	if err != nil {
		t.Fatal(err)
	}
	workspaceBefore := provisionFiles(t, record.Workspace)
	spool := f.env.runnerPaths().Log
	filesBefore := serviceFileHash(spool)
	code, msg := f.run("upgrade", "--binary", candidate)
	if code != 1 || !strings.Contains(msg, "old terminal process retirement not proved (pid "+strconv.Itoa(pid)+")") ||
		pairedUpgradeEvidence(f) != before || !reflect.DeepEqual(workspaceBefore, provisionFiles(t, record.Workspace)) || serviceFileHash(spool) != filesBefore || !s.runnerLoaded {
		t.Fatalf("upgrade passed a live old-build writer: %d %s", code, msg)
	}
	release()
	code, msg = f.run("upgrade", "--binary", candidate)
	if code != 0 || !s.runnerLoaded || f.health.Build.ExecutableSHA256 != serviceFileHash(candidate) {
		t.Fatalf("retired terminal and matched pair did not upgrade: %d %s", code, msg)
	}
}

func testOperatorUpgradeActiveChild(t *testing.T) {
	f, s, candidate := pairedUpgradeFixture(t)
	claimFixtureEpisode(t, f.config.Credentials, "ep").Close()
	old := pairedUpgradeEvidence(f)
	lock, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
	if err != nil {
		t.Fatal(err)
	}
	_, release := operatorProcessWitness(t, "inherited", "", lock)
	lock.Close()

	original := f.env.terminal.openTerminal
	term := &upgradeTerm{r: strings.NewReader("n\n")}
	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
		return term, nil
	}
	code, msg := f.run("upgrade", "--binary", candidate)
	if code != 1 || !strings.Contains(msg, "aborted (not confirmed)") || pairedUpgradeEvidence(f) != old || !s.runnerLoaded {
		t.Fatalf("upgrade crossed typed refusal: %d %s", code, msg)
	}
	for _, line := range []string{"approved companion config: " + old[3], "approved companion plist: " + old[4]} {
		if !strings.Contains(term.String(), "\n"+line+"\n") {
			t.Fatalf("typed prompt missing %q: %s", line, term.String())
		}
	}
	f.env.terminal.openTerminal = original
	code, msg = f.run("upgrade", "--binary", candidate)
	if code != 1 || !strings.Contains(msg, "retained child still holds the runner lock") || pairedUpgradeEvidence(f) != old ||
		s.runnerLoaded || serviceFileHash(candidate) == "" || serviceFileHash(f.env.paths.Binary+".prev") != "" {
		t.Fatalf("upgrade swapped before child exit: %d %s", code, msg)
	}
	if !strings.Contains(msg, "companion drain refused") || !strings.Contains(msg, wantStoppedRunnerRemedy+"\n") {
		t.Fatalf("drain refusal omitted stopped companion remedy: %s", msg)
	}
	release()

	// Restore the pre-upgrade admission state, then exercise rollback after
	// the candidate daemon answers but its companion fails to start.
	s.runnerLoaded, s.bootFailure = true, true
	code, msg = f.run("upgrade", "--binary", candidate)
	rollback := pairedUpgradeEvidence(f)
	if code != 2 || !strings.Contains(msg, "paired upgrade rolled back") || rollback[0] != old[0] ||
		rollback[3] != old[3] || rollback[4] != old[4] || !s.runnerLoaded ||
		serviceFileHash(f.env.paths.Binary+".failed") != serviceFileHash(candidate) {
		t.Fatalf("compatible pair rollback: %d %s", code, msg)
	}
}

func TestOperatorPairedUpgradeRefusals(t *testing.T) {
	f, supervisor, candidate := pairedUpgradeFixture(t)
	before := pairedUpgradeEvidence(f)
	check := func(t *testing.T, want string, evidence [7]string) {
		t.Helper()
		code, msg := f.run("upgrade", "--binary", candidate)
		if code != 1 || !strings.Contains(msg, want) || pairedUpgradeEvidence(f) != evidence ||
			!supervisor.runnerLoaded || f.homeJob.kicks != 0 || serviceFileHash(f.env.paths.Binary+".prev") != "" {
			t.Fatalf("paired refusal %q: %d %s", want, code, msg)
		}
		if strings.Contains(msg, wantStoppedRunnerRemedy) {
			t.Fatalf("loaded companion reported stopped: %s", msg)
		}
	}
	t.Run("schema", func(t *testing.T) {
		probe := f.env.probe
		t.Cleanup(func() { f.env.probe = probe })
		f.env.probe = func(path, home string) ([]byte, error) {
			b, err := probe(path, home)
			if err != nil || path != candidate {
				return b, err
			}
			var v serviceVersion
			if err := json.Unmarshal(b, &v); err != nil {
				return nil, err
			}
			v.CredentialsSchema++
			return json.Marshal(v)
		}
		check(t, "both builds to match the actual credentials schema", before)
	})
	t.Run("process_inventory", func(t *testing.T) {
		read := f.env.runnerProcesses
		t.Cleanup(func() { f.env.runnerProcesses = read })
		f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, errors.New("fixture inventory unavailable") }
		check(t, "cannot verify old terminal process retirement", before)
		f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, nil }
		check(t, "did not include this process", before)
	})
	for _, tc := range []struct{ fault, want string }{
		{"error", "could not stop the companion"},
		{"reported-only", "cannot verify that the operator runner is unloaded"},
	} {
		t.Run("stop_"+tc.fault, func(t *testing.T) {
			supervisor.stopFault = tc.fault
			t.Cleanup(func() { supervisor.stopFault = "" })
			check(t, tc.want, before)
		})
	}
	t.Run("approved_files_changed", func(t *testing.T) {
		terminal := f.env.terminal.openTerminal
		p := f.env.runnerPaths()
		data, err := os.ReadFile(p.Record)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			f.env.terminal.openTerminal = terminal
			if err := os.WriteFile(p.Record, data, 0600); err != nil {
				t.Error(err)
			}
		})
		var changed [7]string
		f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
			if err := os.WriteFile(p.Record, append(data, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
			changed = pairedUpgradeEvidence(f)
			return terminal()
		}
		code, msg := f.run("upgrade", "--binary", candidate)
		if code != 1 || !strings.Contains(msg, "changed during confirmation") ||
			pairedUpgradeEvidence(f) != changed || changed[3] == before[3] || !supervisor.runnerLoaded || f.homeJob.kicks != 0 {
			t.Fatalf("changed approved pair: %d %s", code, msg)
		}
	})
	t.Run("post_drain_census", func(t *testing.T) {
		f, supervisor, candidate := pairedUpgradeFixture(t)
		before := pairedUpgradeEvidence(f)
		read := f.env.runnerProcesses
		pid := os.Getpid() + 1000
		f.env.runnerProcesses = func() ([]worldProcess, error) {
			rows, err := read()
			if !supervisor.runnerLoaded {
				rows = append(rows, worldProcess{pid, "ailang-worldd.prev"})
			}
			return rows, err
		}
		code, msg := f.run("upgrade", "--binary", candidate)
		if code != 1 || !strings.Contains(msg, "terminal drain refused") ||
			!strings.Contains(msg, fmt.Sprintf("old terminal process retirement not proved (pid %d)", pid)) ||
			!strings.Contains(msg, wantStoppedRunnerRemedy+"\n") || supervisor.runnerLoaded ||
			pairedUpgradeEvidence(f) != before || f.homeJob.kicks != 0 ||
			serviceFileHash(f.env.paths.Binary+".new") != "" || serviceFileHash(f.env.paths.Binary+".prev") != "" {
			t.Fatalf("post-drain census did not preserve the installed pair: %d %s", code, msg)
		}
	})
}

func testOperatorUpgradeStoppedCompanion(t *testing.T) {
	f, supervisor, candidate := pairedUpgradeFixture(t)
	supervisor.runnerLoaded = false
	before := pairedUpgradeEvidence(f)

	code, msg := f.run("upgrade", "--binary", candidate)
	after := pairedUpgradeEvidence(f)
	to := serviceFileHash(candidate)
	cfg, _, err := f.env.installedRunner()
	if code != 0 || supervisor.runnerLoaded || f.homeJob.kicks != 1 || err != nil ||
		after[0] != to || cfg.Executables[f.env.paths.Binary] != to ||
		after[5] != before[5] || after[6] != before[6] {
		t.Fatalf("upgrade did not preserve the stopped companion: code=%d loaded=%t config=%v\n%s",
			code, supervisor.runnerLoaded, err, msg)
	}
}

func testOperatorUpgradeRollbackSchema(t *testing.T) {
	f, supervisor, candidate := pairedUpgradeFixture(t)
	kick := f.homeJob.onKick
	var pending [7]string
	f.homeJob.onKick = func() error {
		if err := kick(); err != nil {
			return err
		}
		if f.homeJob.kicks != 1 {
			return nil
		}
		// Simulate the candidate committing a newer schema before it fails.
		credSQL(t, f.config.Credentials, fmt.Sprintf("PRAGMA user_version = %d", store.CredentialsSchemaVersion+1))
		pending = pairedUpgradeEvidence(f)
		return errors.New("fixture candidate failed after schema change")
	}
	before := pairedUpgradeEvidence(f)
	code, msg := f.run("upgrade", "--binary", candidate)
	wantPending := before
	wantPending[0], wantPending[5] = serviceFileHash(candidate), pending[5]
	if code != 3 || !strings.Contains(msg, "paired rollback refused") ||
		!strings.Contains(msg, "is newer than this binary") || !strings.Contains(msg, wantStoppedRunnerRemedy+"\n") ||
		supervisor.runnerLoaded || f.homeJob.kicks != 1 || pending[5] == before[5] || pending[5] == "" ||
		pending != wantPending || pairedUpgradeEvidence(f) != pending ||
		serviceFileHash(f.env.paths.Binary+".prev") != before[0] || serviceFileHash(f.env.paths.Binary+".failed") != "" {
		t.Fatalf("rollback crossed the changed credentials schema: code=%d kicks=%d\n%s", code, f.homeJob.kicks, msg)
	}
}
