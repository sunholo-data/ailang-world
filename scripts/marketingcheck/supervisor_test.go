package marketingcheck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

type supervisorVerdict struct {
	Verdict     string `json:"verdict"`
	Detail      string `json:"detail"`
	ShotVerdict string `json:"shot_verdict"`
}

func pythonSupervisor(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal("python3 is mandatory: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, py, append([]string{"../marketing_capture.py"}, args...)...)
	// A real input payload makes dropping DEVNULL observable, even in headless CI.
	cmd.Stdin = strings.NewReader("input must never reach Chrome\n")
	return cmd
}

func runSupervisor(t *testing.T, want string, args ...string) supervisorVerdict {
	t.Helper()
	cmd := pythonSupervisor(t, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	var got supervisorVerdict
	if e := json.Unmarshal(out, &got); e != nil {
		t.Fatalf("supervisor JSON: %v, stdout=%q stderr=%q err=%v", e, out, stderr.String(), err)
	}
	if got.Verdict != want || (err == nil) != (want == "ok") {
		t.Errorf("verdict must be %s with matching exit status: got %+v err=%v stderr=%s", want, got, err, stderr.String())
	}
	return got
}

func fakeScript(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0700); err != nil {
		t.Fatal(err)
	}
	return path
}

const screenshotPath = `for arg do
 case "$arg" in --screenshot=*) png=${arg#--screenshot=} ;; esac
done
`
const writePNG = `printf '\211PNG\r\n\032\nfixture' > "$png"
`

func shotArgs(chrome, out, scheme string) []string {
	return []string{"shot", "--chrome", chrome, "--out", out, "--url", "http://example.invalid/workbench?object=fake", "--width", "390", "--height", "844", "--scheme", scheme}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestMarketingSupervisorFakeChromeSuccess(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "record")
	t.Setenv("FAKE_RECORD", log)
	chrome := fakeScript(t, dir, "chrome", `printf '%s\n' "$@" > "$FAKE_RECORD.args"
printf '%s' "$HOME" > "$FAKE_RECORD.home"
cat > "$FAKE_RECORD.stdin"
`+screenshotPath+writePNG)
	profiles := map[string]bool{}
	for _, scheme := range []string{"light", "dark"} {
		out := filepath.Join(dir, scheme+".png")
		runSupervisor(t, "ok", shotArgs(chrome, out, scheme)...)
		args := strings.Split(strings.TrimSuffix(readFile(t, log+".args"), "\n"), "\n")
		if len(args) != 11 {
			t.Errorf("complete Chrome argv must have 11 arguments: %q", args)
		}
		profile := strings.TrimPrefix(args[4], "--user-data-dir=")
		if !filepath.IsAbs(profile) || profiles[profile] {
			t.Errorf("profile must be absolute and unique: %q", profile)
		}
		profiles[profile] = true
		schemeValue := "1"
		if scheme == "dark" {
			schemeValue = "0"
		}
		want := []string{"--headless=new", "--no-first-run", "--no-default-browser-check", "--disable-background-networking", "--user-data-dir=" + profile, "--screenshot=" + out, "--window-size=390,844", "--force-device-scale-factor=1", "--lang=en-US", "--blink-settings=preferredColorScheme=" + schemeValue, "http://example.invalid/workbench?object=fake"}
		if !reflect.DeepEqual(args, want) {
			t.Errorf("%s argv must exactly pin scheme, viewport, scale and language:\ngot %q\nwant %q", scheme, args, want)
		}
		if got := readFile(t, log+".home"); got != os.Getenv("HOME") {
			t.Errorf("HOME must be preserved: got %q want %q", got, os.Getenv("HOME"))
		}
		if got := readFile(t, log+".stdin"); got != "" {
			t.Errorf("stdin must be EOF: got %q", got)
		}
		if got := readFile(t, out); !strings.HasPrefix(got, "\x89PNG\r\n\x1a\n") {
			t.Errorf("fresh PNG fixture missing: %q", got)
		}
	}
	stale := filepath.Join(dir, "stale.png")
	if err := os.WriteFile(stale, []byte("old PNG"), 0600); err != nil {
		t.Fatal(err)
	}
	noWrite := fakeScript(t, dir, "no-write", "exit 0\n")
	got := runSupervisor(t, "DestinationExists", shotArgs(noWrite, stale, "light")...)
	if got.Verdict != "DestinationExists" || readFile(t, stale) != "old PNG" {
		t.Error("stale PNG must never satisfy a new shot")
	}
}

func TestMarketingSupervisorFakeChromeNoFile(t *testing.T) {
	for _, tc := range []struct{ name, body, verdict string }{
		{"missing", "exit 0\n", "MissingPNG"},
		{"nonzero", writePNG + "exit 9\n", "NonzeroExit"},
		{"empty", `: > "$png"` + "\n", "EmptyPNG"},
		{"directory", `mkdir "$png"` + "\n", "NotRegularFile"},
		{"symlink", `printf data > "$png.target"; ln -s "$png.target" "$png"` + "\n", "NotRegularFile"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			chrome := fakeScript(t, dir, "chrome", screenshotPath+tc.body)
			got := runSupervisor(t, tc.verdict, shotArgs(chrome, filepath.Join(dir, "shot.png"), "light")...)
			if tc.name == "missing" && got.Verdict != "MissingPNG" {
				t.Error("rc0 with no file must be MissingPNG")
			}
		})
	}
}

// Linux containers can retain adopted zombies; they are dead, even before PID 1 reaps.
func pidDead(pid int) bool {
	if err := syscall.Kill(pid, 0); err == syscall.ESRCH {
		return true
	}
	if b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		if end := bytes.LastIndexByte(b, ')'); end >= 0 && strings.HasPrefix(string(b[end+1:]), " Z") {
			return true
		}
	}
	return false
}

func recordedPIDs(path string) []int {
	b, _ := os.ReadFile(path)
	var pids []int
	for _, s := range strings.Fields(string(b)) {
		if pid, err := strconv.Atoi(s); err == nil && pid > 1 {
			pids = append(pids, pid)
		}
	}
	return pids
}

func cleanupPIDs(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() {
		for _, pid := range recordedPIDs(path) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
		}
	})
}

func waitFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("readiness file not published within bound: %s", path)
}

func assertDead(t *testing.T, pid int, label string) {
	t.Helper()
	deadline := time.Now().Add(600 * time.Millisecond)
	for time.Now().Before(deadline) {
		if pidDead(pid) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("%s must be dead after bounded cleanup (pid %d)", label, pid)
}

func TestMarketingSupervisorFakeChromeTimeoutKillsGroup(t *testing.T) {
	dir := t.TempDir()
	pids := filepath.Join(dir, "pids")
	t.Setenv("FAKE_PIDS", pids)
	cleanupPIDs(t, pids)
	// Each shell records itself; the middle shell waits on its own child.
	grandchild := fakeScript(t, dir, "grandchild", `echo "$$" >> "$FAKE_PIDS"
exec sleep 60
`)
	t.Setenv("FAKE_GRANDCHILD", grandchild)
	child := fakeScript(t, dir, "child", `echo "$$" >> "$FAKE_PIDS"
"$FAKE_GRANDCHILD" &
wait
`)
	t.Setenv("FAKE_CHILD", child)
	chrome := fakeScript(t, dir, "chrome", screenshotPath+writePNG+`"$FAKE_CHILD" &
wait
`)
	start := time.Now()
	out := filepath.Join(dir, "shot.png")
	args := append(shotArgs(chrome, out, "light"), "--timeout", "1.5")
	got := runSupervisor(t, "Timeout", args...)
	if len(readFile(t, out)) == 0 {
		t.Error("timeout fixture must have written a PNG before hanging")
	}
	if got.Verdict != "Timeout" {
		t.Error("timeout must fail even when PNG exists")
	}
	if elapsed := time.Since(start); elapsed > 7*time.Second {
		t.Errorf("timeout must return within deadline + 5s reap + margin: %s", elapsed)
	}
	ids := recordedPIDs(pids)
	if len(ids) != 2 {
		t.Fatalf("fake must record child and grandchild: %v", ids)
	}
	// Grandchild first so the child-only mutant reaches the required assertion.
	assertDead(t, ids[1], "grandchild")
	assertDead(t, ids[0], "child")
}

func daemonFixture(t *testing.T) (string, string, string, string) {
	t.Helper()
	dir := t.TempDir()
	pids, ready, started := filepath.Join(dir, "pids"), filepath.Join(dir, "ready"), filepath.Join(dir, "started")
	t.Setenv("FAKE_PIDS", pids)
	t.Setenv("FAKE_READY", ready)
	cleanupPIDs(t, pids)
	daemon := fakeScript(t, dir, "daemon", `trap '' TERM
echo "$$" > "$FAKE_PIDS"
printf ready > "$FAKE_READY"
while :; do sleep 0.02; done
`)
	return daemon, pids, ready, started
}

func TestMarketingSupervisorFakeDaemonCleanup(t *testing.T) {
	// Assert production values independently of shortened lifecycle unit controls.
	py := pythonSupervisor(t, "--help")
	help, err := py.CombinedOutput()
	if err != nil {
		t.Fatalf("help: %v %s", err, help)
	}
	for _, text := range []string{"shot=30s", "pass=300s", "startup=30s", "TERM=5s", "KILL/reap=5s"} {
		if !strings.Contains(string(help), text) {
			t.Errorf("production default missing: %s", text)
		}
	}
	t.Run("escalation", func(t *testing.T) {
		daemon, pids, ready, started := daemonFixture(t)
		start := time.Now()
		runSupervisor(t, "ok", "daemon", "--ready-file", ready, "--started-file", started, "--lifetime", "0.05", "--term-wait", "0.15", "--kill-wait", "0.3", "--", daemon)
		elapsed := time.Since(start)
		if elapsed < 150*time.Millisecond || elapsed > 2*time.Second {
			t.Errorf("TERM wait then KILL/reap must respect unit bounds: %s", elapsed)
		}
		ids := recordedPIDs(pids)
		if len(ids) != 1 {
			t.Fatalf("daemon PID missing: %v", ids)
		}
		assertDead(t, ids[0], "TERM-ignoring daemon")
	})
	t.Run("cleanup-overrides-shot", func(t *testing.T) {
		daemon, pids, ready, started := daemonFixture(t)
		dir := t.TempDir()
		chrome := fakeScript(t, dir, "chrome", screenshotPath+writePNG)
		out := filepath.Join(dir, "shot.png")
		got := runSupervisor(t, "CleanupFailed", "daemon", "--ready-file", ready, "--started-file", started, "--term-wait", "0.05", "--kill-wait", "0.3", "--unit-cleanup-failure", "--chrome", chrome, "--out", out, "--url", "http://example.invalid/", "--", daemon)
		if got.Verdict != "CleanupFailed" || got.ShotVerdict != "ok" {
			t.Errorf("cleanup failure must override earlier shot success: %+v", got)
		}
		if len(readFile(t, out)) == 0 {
			t.Error("shot must actually succeed before cleanup failure")
		}
		for _, pid := range recordedPIDs(pids) {
			assertDead(t, pid, "injected-failure daemon")
		}
	})
	for _, mode := range []string{"exit", "wait"} {
		t.Run("shell-"+mode, func(t *testing.T) {
			daemon, pids, ready, started := daemonFixture(t)
			ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "bash", "../capture_marketing.sh", "--unit-daemon", mode, daemon, ready, started, "0.15", "0.3")
			outputPath := filepath.Join(t.TempDir(), "shell.log")
			output, err := os.Create(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { output.Close() })
			cmd.Stdout, cmd.Stderr = output, output
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			// Always stop the shell; PID cleanup remains registered even on mutant failures.
			t.Cleanup(func() {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
				b, _ := os.ReadFile(started)
				var ids struct {
					Supervisor int `json:"supervisor_pid"`
				}
				if json.Unmarshal(b, &ids) == nil && ids.Supervisor > 1 {
					_ = syscall.Kill(ids.Supervisor, syscall.SIGKILL)
				}
			})
			waitFile(t, started)
			start := time.Now()
			if mode == "wait" {
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					t.Fatal(err)
				}
			}
			err = cmd.Wait()
			if mode == "exit" && err != nil {
				t.Errorf("normal shell cleanup: %v %s", err, readFile(t, outputPath))
			}
			if mode == "wait" {
				exitErr, ok := err.(*exec.ExitError)
				if !ok || exitErr.ExitCode() != 143 {
					t.Errorf("TERM trap must exit 143 after supervised cleanup: %v", err)
				}
			}
			if elapsed := time.Since(start); elapsed > 2*time.Second {
				t.Errorf("shell trap cleanup exceeded bound: %s", elapsed)
			}
			ids := recordedPIDs(pids)
			if len(ids) != 1 {
				t.Fatalf("shell daemon PID missing: %v, %s", ids, readFile(t, outputPath))
			}
			assertDead(t, ids[0], "shell "+mode+" daemon")
		})
	}
}
