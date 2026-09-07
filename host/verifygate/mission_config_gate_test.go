package verifygate

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

const missionDisclosure = "World decision-ledger validated; centralized runtime currency is NOT certified here."

type missionSyncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *missionSyncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *missionSyncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func runMissionCommand(t *testing.T, dir string, env []string, name string, args ...string) (int, string, bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output missionSyncBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		select {
		case <-done:
		case <-time.After(time.Second):
			return -1, "instrument failure: context deadline exceeded; process-group cleanup also exceeded 1s\n" + output.String(), true
		}
		return -1, "instrument failure: context deadline exceeded\n" + output.String(), true
	}
	if err == nil {
		return 0, output.String(), false
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String(), false
	}
	t.Fatalf("start %s: %v", name, err)
	return -1, output.String(), false
}

func missionEnv(extra ...string) []string {
	var env []string
	for _, entry := range os.Environ() {
		key := strings.SplitN(entry, "=", 2)[0]
		if key == "AILANG_BIN" || strings.HasPrefix(key, "AILANG_FLEET") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, extra...)
}

func exactReplace(t *testing.T, src, old, replacement string) string {
	t.Helper()
	if count := strings.Count(src, old); count != 1 {
		t.Fatalf("instrument failure: mutation needle count=%d, want 1 for %q", count, old)
	}
	return strings.Replace(src, old, replacement, 1)
}

func missionGateSource(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join(repoRoot, "scripts", "verify_go.sh")
	if override := os.Getenv("MISSION_CONFIG_GATE_SOURCE"); override != "" {
		path = override
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func writeFile(t *testing.T, path string, raw []byte, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, mode); err != nil {
		t.Fatal(err)
	}
}

func fixtureLedger(rows string) string {
	return "# Fixture\n\n<!-- decision-ledger:start -->\n" + rows + "<!-- decision-ledger:end -->\n"
}

func newMissionFixture(t *testing.T, ledger string, controlledStop bool) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "world")
	script := string(missionGateSource(t))
	switch os.Getenv("MISSION_CONFIG_TEST_MUTATION") {
	case "M1":
		script = exactReplace(t, script, "\ncheck_mission_config\n", "\n")
	case "M2":
		script = exactReplace(t, script,
			"  /bin/bash \"$validator\" --check --file \"$charter\"\n",
			"  /bin/bash \"$validator\" --check --file \"$charter\" || true\n")
	case "M3":
		body := `check_mission_config() {
  validator="scripts/mission_decisions.sh"
  charter="design_docs/world-mission.md"
  if [ ! -r "$validator" ]; then
    echo "verify_go.sh: FATAL: mission validator is missing or unreadable: $validator" >&2
    return 1
  fi
  if [ ! -r "$charter" ]; then
    echo "verify_go.sh: FATAL: World charter is missing or unreadable: $charter" >&2
    return 1
  fi
  /bin/bash -n "$validator"
  /bin/bash "$validator" --check --file "$charter"
  echo "   World decision-ledger validated; centralized runtime currency is NOT certified here."
}`
		script = exactReplace(t, script, body, "check_mission_config() {\n  return 0\n}")
	case "M9":
		refusal := `if [ "$#" -gt 0 ]; then
  echo "verify_go.sh: unrecognized first argument: $1" >&2
  exit 2
}

`
		script = exactReplace(t, script, refusal, "")
	}
	if controlledStop {
		script = exactReplace(t, script, "\ncheck_mission_config\n", "\ncheck_mission_config\necho 'CONTROLLED_STAGE_AFTER_MISSION'\nexit 0\n")
	}
	writeFile(t, filepath.Join(root, "scripts", "verify_go.sh"), []byte(script), 0o755)
	validator, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "mission_decisions.sh"))
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "scripts", "mission_decisions.sh"), validator, 0o755)
	writeFile(t, filepath.Join(root, "design_docs", "world-mission.md"), []byte(ledger), 0o644)
	if controlledStop {
		shim := "#!/bin/sh\nif [ \"${1:-}\" = --version ]; then echo 'AILANG v0.30.0'; exit 0; fi\nexit 0\n"
		writeFile(t, filepath.Join(root, "bin", "ailang"), []byte(shim), 0o755)
		for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "-qm", "fixture"}} {
			rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "git", args...)
			if timedOut || rc != 0 {
				t.Fatalf("fixture git command failed: rc=%d timeout=%v\n%s", rc, timedOut, out)
			}
		}
	}
	return root
}

const validFixtureRows = "| ID | Status | Decision | Evidence |\n|---|---|---|---|\n| D-WORLD-5 | RESOLVED | yes | measured |\n"

func TestMissionConfigRealCharter(t *testing.T) {
	rc, out, timedOut := runMissionCommand(t, repoRoot, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
	if timedOut || rc != 0 {
		t.Fatalf("real mission-config mode failed: rc=%d timeout=%v\n%s", rc, timedOut, out)
	}
	match := regexp.MustCompile(`decision ledger valid: ([1-9][0-9]*) rows`).FindStringSubmatch(out)
	if match == nil || !strings.Contains(out, missionDisclosure) {
		t.Fatalf("mission-config output lacks nonzero count or disclosure:\n%s", out)
	}
}

func TestMissionConfigRejectsBadLedger(t *testing.T) {
	cases := []struct{ name, ledger, want string }{
		{"absent_ledger_block", "# no ledger\n", "expected exactly one decision-ledger block (start=0 end=0)"},
		{"zero_row_ledger", fixtureLedger("| ID | Status | Decision | Evidence |\n|---|---|---|---|\n"), "decision ledger has no rows"},
		{"duplicate_id", fixtureLedger(validFixtureRows + "| D-WORLD-5 | OPEN | again | measured |\n"), "duplicate decision ID: D-WORLD-5"},
		{"invalid_status", fixtureLedger("| D-WORLD-5 | MAYBE | yes | measured |\n"), "invalid status for D-WORLD-5: MAYBE"},
		{"empty_required_field", fixtureLedger("| D-WORLD-5 | OPEN |  | measured |\n"), "empty decision/evidence field for D-WORLD-5"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newMissionFixture(t, tc.ledger, false)
			rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
			if timedOut || rc == 0 || !strings.Contains(out, tc.want) {
				t.Fatalf("bad ledger not propagated: rc=%d timeout=%v want=%q\n%s", rc, timedOut, tc.want, out)
			}
		})
	}
	t.Run("missing_charter", func(t *testing.T) {
		root := newMissionFixture(t, fixtureLedger(validFixtureRows), false)
		if err := os.Rename(filepath.Join(root, "design_docs", "world-mission.md"), filepath.Join(root, "design_docs", "elsewhere.md")); err != nil {
			t.Fatal(err)
		}
		rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
		if timedOut || rc == 0 || !strings.Contains(out, "World charter is missing or unreadable") {
			t.Fatalf("missing charter not refused: rc=%d timeout=%v\n%s", rc, timedOut, out)
		}
	})
	t.Run("missing_validator", func(t *testing.T) {
		root := newMissionFixture(t, fixtureLedger(validFixtureRows), false)
		if err := os.Rename(filepath.Join(root, "scripts", "mission_decisions.sh"), filepath.Join(root, "scripts", "validator-away.sh")); err != nil {
			t.Fatal(err)
		}
		rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
		if timedOut || rc == 0 || !strings.Contains(out, "mission validator is missing or unreadable") {
			t.Fatalf("missing validator not refused: rc=%d timeout=%v\n%s", rc, timedOut, out)
		}
	})
	t.Run("valid_control", func(t *testing.T) {
		root := newMissionFixture(t, fixtureLedger(validFixtureRows), false)
		rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
		if timedOut || rc != 0 || !strings.Contains(out, "decision ledger valid: 1 rows") {
			t.Fatalf("valid control failed: rc=%d timeout=%v\n%s", rc, timedOut, out)
		}
	})
}

func TestMissionConfigBlockingTimeoutClassified(t *testing.T) {
	root := newMissionFixture(t, fixtureLedger(validFixtureRows), false)
	writeFile(t, filepath.Join(root, "scripts", "mission_decisions.sh"), []byte("#!/bin/bash\nsleep 30\n"), 0o755)
	_, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "scripts/verify_go.sh", "--mission-config-check")
	if !timedOut || !strings.Contains(out, "instrument failure: context deadline exceeded") {
		t.Fatalf("blocking fixture was not classified as instrument failure: timeout=%v\n%s", timedOut, out)
	}
}

func TestMissionConfigBashSyntax(t *testing.T) {
	root := newMissionFixture(t, fixtureLedger(validFixtureRows), false)
	rc, out, timedOut := runMissionCommand(t, root, missionEnv(), "/bin/bash", "-n", "scripts/verify_go.sh")
	if timedOut || rc != 0 {
		t.Fatalf("copied verify_go.sh fails bash -n: rc=%d timeout=%v\n%s", rc, timedOut, out)
	}
}

func TestMissionConfigDefaultFlowIntegration(t *testing.T) {
	for _, tc := range []struct {
		name, ledger string
		valid        bool
	}{
		{"valid_ledger", fixtureLedger(validFixtureRows), true},
		{"invalid_ledger", fixtureLedger("| D-WORLD-5 | MAYBE | yes | measured |\n"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newMissionFixture(t, tc.ledger, true)
			bin := filepath.Join(root, "bin")
			env := missionEnv("AILANG_BIN="+filepath.Join(bin, "ailang"), "PATH="+bin+":"+os.Getenv("PATH"))
			rc, out, timedOut := runMissionCommand(t, root, env, "/bin/bash", "scripts/verify_go.sh")
			if timedOut {
				t.Fatalf("default-flow fixture timed out\n%s", out)
			}
			marker := strings.Index(out, "CONTROLLED_STAGE_AFTER_MISSION")
			rows := strings.Index(out, "decision ledger valid:")
			if tc.valid && (rc != 0 || rows < 0 || marker <= rows) {
				t.Fatalf("valid flow did not validate before controlled stage: rc=%d rows=%d marker=%d\n%s", rc, rows, marker, out)
			}
			if !tc.valid && (rc == 0 || marker >= 0 || strings.Contains(out, "go version")) {
				t.Fatalf("invalid flow reached controlled/Go stage: rc=%d marker=%d\n%s", rc, marker, out)
			}
		})
	}
}

func activeLines(raw string) []string {
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			lines = append(lines, line)
		}
	}
	return lines
}

func TestActiveWiringHasNoDriverConsumers(t *testing.T) {
	ciRaw, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	scriptRaw, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
	if err != nil {
		t.Fatal(err)
	}
	ci := string(ciRaw)
	if len(activeLines(ci)) == 0 || !strings.Contains(ci, "./scripts/verify_go.sh") || !strings.Contains(ci, "./scripts/verify_ail.sh") {
		t.Fatal("instrument failure: CI active-command controls are absent")
	}
	jobsStart := strings.Index(ci, "\njobs:\n")
	if jobsStart < 0 {
		t.Fatal("instrument failure: ci.yml jobs anchor absent")
	}
	jobs := regexp.MustCompile(`(?m)^  [a-z0-9-]+:$`).FindAllString(ci[jobsStart:], -1)
	if len(jobs) != 2 {
		t.Errorf(".github/workflows/ci.yml: enumerated %d jobs, want 2: %v", len(jobs), jobs)
	}
	for i, line := range strings.Split(ci, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") && strings.Contains(line, "tools/launchd/") {
			t.Errorf(".github/workflows/ci.yml:%d: active tools/launchd consumer: %s", i+1, strings.TrimSpace(line))
		}
	}
	// Only inspect the default-flow region; function definitions for the diagnostic
	// still exist until MS4 and are deliberately not active default consumers.
	defaultStart := strings.Index(string(scriptRaw), `if [ -z "${AILANG_BIN:-}" ]`)
	if defaultStart < 0 {
		t.Fatal("instrument failure: verify_go.sh default-flow anchor absent")
	}
	for i, line := range strings.Split(string(scriptRaw)[defaultStart:], "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") && strings.Contains(line, "tools/launchd/") {
			t.Errorf("scripts/verify_go.sh:%d: active tools/launchd consumer: %s", strings.Count(string(scriptRaw)[:defaultStart], "\n")+i+1, strings.TrimSpace(line))
		}
	}
}

func TestRetirementCompleteness(t *testing.T) {
	paths := []string{"host", "scripts", ".github"}
	needle := regexp.MustCompile("check_driver_" + "fleet|driver-" + "fleet-check|AILANG_" + "FLEET_REPO|REQUIRED_" + "FLEET_PATHS|driver_" + "fleet_scope")
	positive, negative := 0, 0
	for _, rel := range paths {
		err := filepath.Walk(filepath.Join(repoRoot, rel), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			raw, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			if needle.Match(raw) {
				t.Errorf("retired fleet diagnostic remains in %s", path)
			}
			positive += strings.Count(string(raw), "--mission-config-check")
			negative += strings.Count(string(raw), "zzq-nonsense-"+"literal-9f3")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if positive == 0 || negative != 0 {
		t.Fatalf("retirement scan controls invalid: positive=%d negative=%d", positive, negative)
	}
	retiredFixture := "driver_" + "fleet_scope_gate_test.go"
	if _, err := os.Stat(filepath.Join(repoRoot, "host", "verifygate", retiredFixture)); !os.IsNotExist(err) {
		t.Fatalf("retired fixture is still present or stat failed: %v", err)
	}
}

func TestRetiredModeIsRefused(t *testing.T) {
	mutation := os.Getenv("MISSION_CONFIG_TEST_MUTATION")
	fleetFlag := "--driver-" + "fleet-check"
	nonDashFleet := "driver-" + "fleet-check"
	if mutation != "M9" {
		for _, arg := range []string{fleetFlag, "--zzq-unknown-mode-9f3", nonDashFleet} {
			t.Run("real_"+strings.ReplaceAll(strings.TrimLeft(arg, "-"), "-", "_"), func(t *testing.T) {
				rc, out, timedOut := runMissionCommand(t, repoRoot, missionEnv("AILANG_BIN="+os.Getenv("AILANG_BIN")), "/bin/bash", "scripts/verify_go.sh", arg)
				if timedOut || rc == 0 || !strings.Contains(out, arg) || strings.Contains(out, "── AILANG_BIN=") || strings.Contains(out, "✓ go gate PASSED") {
					t.Fatalf("unknown mode was not refused before preflight: rc=%d timeout=%v arg=%s\n%s", rc, timedOut, arg, out)
				}
			})
		}
	}
	t.Run("driver_fleet_check", func(t *testing.T) {
		root := newMissionFixture(t, fixtureLedger(validFixtureRows), true)
		bin := filepath.Join(root, "bin")
		env := missionEnv("AILANG_BIN="+filepath.Join(bin, "ailang"), "PATH="+bin+":"+os.Getenv("PATH"))
		rc, out, timedOut := runMissionCommand(t, root, env, "/bin/bash", "scripts/verify_go.sh", fleetFlag)
		if timedOut || rc == 0 || !strings.Contains(out, fleetFlag) || strings.Contains(out, "── AILANG_BIN=") || strings.Contains(out, "✓ go gate PASSED") {
			t.Fatalf("fixture retired mode was not refused before preflight: rc=%d timeout=%v\n%s", rc, timedOut, out)
		}
	})
}
