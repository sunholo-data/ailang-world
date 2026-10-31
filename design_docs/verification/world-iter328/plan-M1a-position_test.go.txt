package main

// Row 231 M1 (design_docs/planned/w-marketing-position-checker.md §2.2–§2.4, §3 M1, §4 P1–P17):
// `marketingcheck position` on the V64/V65 corpora and ten in-test arms, its parity with the
// banked oracle, its exit codes, and the deadline-enforcing subprocess runner both parity tests use.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// ---- the deadline-enforcing subprocess runner (design §2.2 [r2], §4 P17) ----

const (
	subprocBuildLimit = 120 * time.Second // oracle and CLI builds
	subprocRunLimit   = 30 * time.Second  // one checker run
	subprocDrainLimit = 5 * time.Second   // output draining after exit or kill
)

type subprocResult struct {
	Stdout, Stderr string
	Code           int
}

// boundedRun runs name with args in dir (empty = the package directory) under limit. The child
// starts in its own process group; on expiry the whole group is killed and the error is
// `subprocess-timeout: <command> exceeded <limit>`. Output draining is bounded by
// subprocDrainLimit, and the group is reaped after every run. It never skips or substitutes a
// result: a non-zero exit is returned as Code, anything else as an error.
func boundedRun(limit time.Duration, dir string, env []string, name string, args ...string) (subprocResult, error) {
	what := strings.Join(append([]string{name}, args...), " ")
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = subprocDrainLimit
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	res := subprocResult{Stdout: stdout.String(), Stderr: stderr.String(), Code: -1}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return res, fmt.Errorf("subprocess-timeout: %s exceeded %s", what, limit)
	}
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.Code = ee.ExitCode()
			return res, nil
		}
		return res, fmt.Errorf("subprocess %s: %v", what, err)
	}
	res.Code = 0
	return res, nil
}

// mustRun is boundedRun that fails the test on a timeout or a start/drain error.
func mustRun(t *testing.T, limit time.Duration, dir string, env []string, name string, args ...string) subprocResult {
	t.Helper()
	res, err := boundedRun(limit, dir, env, name, args...)
	if err != nil {
		t.Fatalf("%v\nstdout: %s\nstderr: %s", err, res.Stdout, res.Stderr)
	}
	return res
}

// procAlive reports whether pid is a live (non-zombie) process.
func procAlive(t *testing.T, pid int) bool {
	t.Helper()
	if syscall.Kill(pid, 0) != nil {
		return false
	}
	res, err := boundedRun(10*time.Second, "", nil, "ps", "-o", "stat=", "-p", strconv.Itoa(pid))
	if err != nil {
		t.Fatalf("ps: %v", err)
	}
	stat := strings.TrimSpace(res.Stdout)
	return stat != "" && !strings.HasPrefix(stat, "Z")
}

// TestMarketingSubprocessDeadline is the runner's own control (design §4 P17, AC-4): a fake
// checker that sleeps past the bound, with a grandchild in its process group.
func TestMarketingSubprocessDeadline(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "hang.sh")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nsleep 20 &\necho $! > \"$1\"\nwait\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	pidFile := filepath.Join(dir, "grandchild.pid")
	const limit = 2 * time.Second
	start := time.Now()
	_, err := boundedRun(limit, dir, nil, fake, pidFile)
	elapsed := time.Since(start)
	if err == nil || !strings.Contains(err.Error(), "subprocess-timeout: "+fake) {
		t.Fatalf("want a subprocess-timeout failure naming %s, got %v after %s", fake, err, elapsed)
	}
	if elapsed > limit+subprocDrainLimit+3*time.Second {
		t.Fatalf("the runner returned after %s, bound %s", elapsed, limit+subprocDrainLimit)
	}
	b, rerr := os.ReadFile(pidFile)
	pid, perr := strconv.Atoi(strings.TrimSpace(string(b)))
	if rerr != nil || perr != nil || pid <= 0 {
		t.Fatalf("grandchild pid not recorded: %v %v %q", rerr, perr, b)
	}
	deadline := time.Now().Add(3 * time.Second)
	for procAlive(t, pid) {
		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("grandchild %d survived the process-group kill", pid)
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Logf("subprocess-timeout after %s; grandchild %d reaped", elapsed.Round(time.Millisecond), pid)
}
