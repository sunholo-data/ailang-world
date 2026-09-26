// Package proctest holds the guarded process-observation helpers the host
// subprocess tests use. It exists because three design iterations (190, 191)
// were killed silently by a probe that passed pid -1, parsed from a pid file
// that was never written, to kill(2): kill(-1, SIGKILL) signals every process
// the user owns. Every pid that reaches a signal goes through Check first.
package proctest

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Check refuses a pid that must never be signalled from a test: <= 1 (-1 is
// "every process", 0 is "my group", 1 is init), this process and its parent.
func Check(pid int) error {
	if pid <= 1 || pid == os.Getpid() || pid == os.Getppid() {
		return fmt.Errorf("proctest: refusing to signal pid %d", pid)
	}
	return nil
}

// Signal sends sig through send only after Check accepts pid.
func Signal(pid int, sig syscall.Signal, send func(int, syscall.Signal) error) error {
	if err := Check(pid); err != nil {
		return err
	}
	return send(pid, sig)
}

// Alive reports whether pid exists (kill -0), fatally refusing a bad pid.
func Alive(t testing.TB, pid int) bool {
	t.Helper()
	err := Signal(pid, 0, syscall.Kill)
	if err != nil && Check(pid) != nil {
		t.Fatalf("%v", err)
	}
	return err == nil || errors.Is(err, syscall.EPERM)
}

// DeadWithin polls pid's own state; it never infers death from elapsed time.
func DeadWithin(t testing.TB, pid int, window time.Duration) bool {
	t.Helper()
	for start := time.Now(); time.Since(start) < window; time.Sleep(time.Millisecond) {
		if !Alive(t, pid) {
			return true
		}
	}
	return !Alive(t, pid)
}

// ReadPid reads a fixture-written pid. A missing or malformed file is a test
// failure, never a sentinel that could flow on to kill(2).
func ReadPid(t testing.TB, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture did not record grandchild pid: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("fixture pid %q: %v", data, err)
	}
	if err := Check(pid); err != nil {
		t.Fatalf("%v", err)
	}
	return pid
}

// ReapOnCleanup SIGKILLs a checked pid when the test ends, so no fixture
// sleep outlives the test.
func ReapOnCleanup(t testing.TB, pid int) {
	t.Helper()
	if err := Check(pid); err != nil {
		t.Fatalf("%v", err)
	}
	t.Cleanup(func() { _ = Signal(pid, syscall.SIGKILL, syscall.Kill) })
}
