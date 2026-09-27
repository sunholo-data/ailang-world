package capsule

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/proctest"
)

const overflowLoop = `i=0; while [ $i -lt 200 ]; do echo 0123456789abcdef0123456789abcdef; i=$((i+1)); done`

// TestEscapeeHelper is the helper mode behind proctest.EscapeeShell; in a
// normal run it returns at once.
func TestEscapeeHelper(t *testing.T) { proctest.RunEscapeeIfRequested() }

func blockedWriteHelper() {
	if os.Getenv("AILANG_BLOCKED_WRITE") != "1" {
		return
	}
	dir, n := flag.Arg(0), flag.Arg(1)
	size, err := strconv.Atoi(n)
	if err != nil || size < 1 {
		os.Exit(2)
	}
	if err := os.WriteFile(filepath.Join(dir, "child.pid"), []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		os.Exit(3)
	}
	written, err := syscall.Write(1, make([]byte, size))
	marker := "wrote-all"
	if err != nil {
		marker = "write-error"
	} else if written != size {
		marker = "short-write"
	}
	_ = os.WriteFile(filepath.Join(dir, marker), []byte(fmt.Sprintf("%d/%d: %v", written, size, err)), 0o600)
	os.Exit(0)
}

func pipeCapacity(t *testing.T) int {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	fd := int(w.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		t.Fatal(err)
	}
	capacity := 0
	chunk := make([]byte, 4096)
	for {
		n, err := syscall.Write(fd, chunk)
		if n > 0 {
			capacity += n
		}
		if errors.Is(err, syscall.EAGAIN) {
			break
		}
		if err != nil || n <= 0 || capacity > int(^uint(0)>>1)/4 {
			t.Fatalf("measure pipe capacity: bytes=%d, write=%d, error=%v", capacity, n, err)
		}
	}
	if capacity <= 0 {
		t.Fatalf("invalid pipe capacity %d", capacity)
	}
	return capacity
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }

// The child's single write exceeds all space in the pipe plus the 65 bytes
// readCapped can drain. Its own markers distinguish a blocked writer from one
// that finished before the overflow kill.
func TestBlockedChildOverflowKill(t *testing.T) {
	blockedWriteHelper()
	const limit = 64
	capacity := pipeCapacity(t)
	payload := 4 * capacity
	if payload <= capacity+limit+1 {
		t.Fatalf("payload %d does not exceed pipe %d plus cap %d", payload, capacity, limit+1)
	}
	t.Logf("pipe capacity=%d payload=%d limit=%d", capacity, payload, limit)
	dir := t.TempDir()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(dir, "ailang-blocked")
	script := fmt.Sprintf("#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo 'AILANG v0.30.0 blocked probe'; exit 0; fi\nAILANG_BLOCKED_WRITE=1 exec %s -test.run='^TestBlockedChildOverflowKill$' %s %d\n", shellQuote(exe), shellQuote(dir), payload)
	if err := os.WriteFile(wrapper, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	fx := archiveExecutable(t, wrapper)
	pidPath := filepath.Join(dir, "child.pid")
	killChildGroup := func() {
		data, err := os.ReadFile(pidPath)
		if err != nil {
			return
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil && proctest.Check(pid) == nil {
			_ = syscall.Kill(-pid, syscall.SIGKILL)
		}
	}
	t.Cleanup(killChildGroup)
	type outcome struct {
		result Result
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := New(fx.archive, Config{MaxOutputBytes: limit, ExecTimeout: 120 * time.Second}).Run(Entry{
			Interpreter: fx.ref, Source: source(`export func main() -> string { "x" }`),
		})
		done <- outcome{result, err}
	}()
	var got outcome
	select {
	case got = <-done:
	case <-time.After(30 * time.Second):
		noMarker := true
		for _, marker := range []string{"wrote-all", "short-write", "write-error"} {
			if data, err := os.ReadFile(filepath.Join(dir, marker)); err == nil {
				noMarker = false
				t.Errorf("watchdog found child %s marker before cleanup: %s", marker, data)
			} else if !errors.Is(err, os.ErrNotExist) {
				noMarker = false
				t.Errorf("watchdog read %s marker: %v", marker, err)
			}
		}
		if noMarker {
			t.Log("watchdog: no child completion marker before cleanup")
		}
		killChildGroup()
		reaped := false
		select {
		case <-done:
			reaped = true
		case <-time.After(5 * time.Second):
			t.Error("watchdog cleanup did not finish Run")
		}
		if data, err := os.ReadFile(pidPath); err == nil {
			pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
			if parseErr == nil && proctest.Check(pid) == nil {
				if !proctest.DeadWithin(t, pid, time.Second) {
					t.Errorf("watchdog cleanup left child %d alive", pid)
				} else {
					t.Logf("watchdog cleanup: child %d dead, Run reaped=%t", pid, reaped)
				}
			}
		} else {
			t.Errorf("watchdog found no child PID marker: %v", err)
		}
		t.Fatal("30 s watchdog fired: blocked child was not killed by overflow")
	}
	pid := proctest.ReadPid(t, pidPath)
	proctest.ReapOnCleanup(t, pid)
	for _, marker := range []string{"wrote-all", "short-write", "write-error"} {
		if data, err := os.ReadFile(filepath.Join(dir, marker)); err == nil {
			t.Errorf("child recorded %s: %s", marker, data)
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("read %s marker: %v", marker, err)
		}
	}
	var overflow *OutputLimitError
	if !errors.As(got.err, &overflow) {
		t.Errorf("error = %T %v, want *OutputLimitError", got.err, got.err)
	}
	var timeout *TimeoutError
	if errors.As(got.err, &timeout) {
		t.Errorf("unexpected *TimeoutError: %v", got.err)
	}
	if len(got.result.Stdout) > limit || len(got.result.Stderr) > limit {
		t.Errorf("captured bytes stdout=%d stderr=%d, limit=%d", len(got.result.Stdout), len(got.result.Stderr), limit)
	}
	if !proctest.DeadWithin(t, pid, time.Second) {
		t.Errorf("child %d still alive after Run returned", pid)
	}
}

func scriptedInterpreter(t *testing.T, before, after string) archivedFixture {
	t.Helper()
	p := filepath.Join(t.TempDir(), "ailang-fork")
	script := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo \"AILANG v0.30.0 fork probe\"; exit 0; fi\n" +
		before + "\n" + overflowLoop + "\n" + after + "\n"
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return archiveExecutable(t, p)
}

func runScripted(t *testing.T, fx archivedFixture, execTimeout time.Duration) error {
	t.Helper()
	_, err := New(fx.archive, Config{ExecTimeout: execTimeout, MaxOutputBytes: 1024}).Run(Entry{
		Interpreter: fx.ref, Source: source(`export func main() -> string { "x" }`)})
	return err
}

// recordFailingKill replaces killGroup with one that records the (cmd.Process
// derived) pgid and fails with EPERM WITHOUT signalling.
func recordFailingKill(t *testing.T) func() int {
	t.Helper()
	var mu sync.Mutex
	var pgids []int
	orig := killGroup
	t.Cleanup(func() { killGroup = orig })
	killGroup = func(pgid int) error {
		mu.Lock()
		defer mu.Unlock()
		pgids = append(pgids, pgid)
		return syscall.EPERM
	}
	return func() int {
		mu.Lock()
		defer mu.Unlock()
		if len(pgids) == 0 {
			t.Fatal("killGroup was never called")
		}
		return pgids[0]
	}
}

// AC1. The overflow kill must reach a forked grandchild holding the inherited
// pipes. Observable: the grandchild's own state. It records its pid BEFORE the
// interpreter writes, and writes a survival marker only if its sleep
// completes. ExecTimeout (120 s) is 12x the grandchild's lifetime (10 s), so a
// direct-child-only kill returns after the grandchild finished: marker present.
func TestOverflowKillReachesForkedGrandchild(t *testing.T) {
	dir := t.TempDir()
	pidf, survived := filepath.Join(dir, "gc.pid"), filepath.Join(dir, "survived")
	fx := scriptedInterpreter(t, fmt.Sprintf("(sleep 10; : > %s) &\necho $! > %s", survived, pidf), "wait")
	base := procbound.Outstanding()
	err := runScripted(t, fx, 120*time.Second)
	if got := procbound.Outstanding(); got != base {
		t.Fatalf("reservations after a reaped Run = %d, want %d (released at reap)", got, base)
	}
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError", err, err)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("cleanup reported incomplete although the child exited: %v", err)
	}
	if _, statErr := os.Stat(survived); statErr == nil {
		t.Fatalf("grandchild %d ran to completion: the overflow kill did not reach it", gc)
	}
	if !proctest.DeadWithin(t, gc, time.Second) {
		t.Fatalf("grandchild %d still alive after Run returned", gc)
	}
}

// AC2. A failed overflow kill is surfaced, and the typed overflow stays primary.
func TestOverflowKillErrorJoinedBehindTypedError(t *testing.T) {
	injected := errors.New("injected kill failure")
	orig := killGroup
	t.Cleanup(func() { killGroup = orig })
	killGroup = func(pgid int) error { return errors.Join(orig(pgid), injected) }
	err := runScripted(t, scriptedInterpreter(t, ":", "wait"), 120*time.Second)
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) || !errors.Is(err, injected) {
		t.Fatalf("error = %v, want *OutputLimitError joined with the kill failure", err)
	}
}

// AC4. A descendant that LEFT the group (setsid) is out of the kill's reach;
// the pipe-close bound must still return Run. Observable: Run returned while
// the escapee is ALIVE (without the bound Run waits for it to exit).
func TestPipeCloseBoundsEscapedDescendant(t *testing.T) {
	pidf := filepath.Join(t.TempDir(), "gc.pid")
	err := runScripted(t, scriptedInterpreter(t, proctest.EscapeeShell(t, pidf), "wait"), 3*time.Second)
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError", err, err)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("cleanup reported incomplete although the child exited: %v", err)
	}
	if !proctest.Alive(t, gc) {
		t.Fatalf("escapee %d not alive at return: Run waited for it instead of closing the pipes", gc)
	}
}

// AC7. Both kills fail (EPERM, nothing signalled) and the direct child lives
// on: Run must still return, with the typed overflow, the kill failure and the
// incomplete cleanup all discoverable, while the child is ALIVE. A guarded
// cleanup then kills it and the background waiter reaps it.
func TestFailedKillBoundsRunAndChildIsReapedLater(t *testing.T) {
	pgidOf := recordFailingKill(t)
	base := procbound.Outstanding()
	err := runScripted(t, scriptedInterpreter(t, ":", "exec sleep 30"), time.Second)
	pgid := pgidOf()
	killed := false
	t.Cleanup(func() {
		if !killed {
			_ = proctest.Signal(pgid, syscall.SIGKILL, syscall.Kill)
		}
	})
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) || !errors.Is(err, syscall.EPERM) || !errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("error = %v, want *OutputLimitError + EPERM + ErrCleanupIncomplete", err)
	}
	if !proctest.Alive(t, pgid) {
		t.Fatalf("direct child %d not alive at return: Run waited for it", pgid)
	}
	if got := procbound.Outstanding(); got != base+1 {
		t.Fatalf("outstanding waiters = %d, want %d", got, base+1)
	}
	if err := proctest.Signal(pgid, syscall.SIGKILL, syscall.Kill); err != nil {
		t.Fatal(err)
	}
	killed = true
	for start := time.Now(); procbound.Outstanding() != base && time.Since(start) < 5*time.Second; time.Sleep(time.Millisecond) {
	}
	if got := procbound.Outstanding(); got != base {
		t.Fatalf("background waiter did not reap the killed child: outstanding %d, want %d", got, base)
	}
}

// AC8. A full set of reservations refuses a new Run before Start.
func TestFullCleanupBacklogRefusesRun(t *testing.T) {
	var releases []func()
	defer func() {
		for _, release := range releases {
			release()
		}
	}()
	for procbound.Outstanding() < procbound.MaxOutstanding {
		release, err := procbound.Admit()
		if err != nil {
			break
		}
		releases = append(releases, release)
	}
	err := runScripted(t, scriptedInterpreter(t, ":", "wait"), 120*time.Second)
	if !errors.Is(err, procbound.ErrCleanupBacklog) {
		t.Fatalf("error = %v, want ErrCleanupBacklog", err)
	}
}

// A Start failure releases the reservation Admit took (the archived
// interpreter is made non-executable after its hash was recorded).
func TestStartFailureReleasesReservation(t *testing.T) {
	fx := scriptedInterpreter(t, ":", "wait")
	if err := os.Chmod(fx.path, 0o644); err != nil {
		t.Fatal(err)
	}
	base := procbound.Outstanding()
	err := runScripted(t, fx, 120*time.Second)
	var execErr *ExecError
	if !errors.As(err, &execErr) {
		t.Fatalf("error = %T %v, want *ExecError from Start", err, err)
	}
	if got := procbound.Outstanding(); got != base {
		t.Fatalf("reservations after a failed Start = %d, want %d", got, base)
	}
}

// AC9 (r2). An ordinary overflow, real killGroup, child exits on its own:
// the error is exactly the typed overflow and no ESRCH is anywhere in it.
func TestOrdinaryOverflowCarriesNoESRCH(t *testing.T) {
	err := runScripted(t, scriptedInterpreter(t, ":", ""), 120*time.Second)
	if _, ok := err.(*OutputLimitError); !ok {
		t.Fatalf("error = %T %v, want exactly *OutputLimitError", err, err)
	}
	if errors.Is(err, syscall.ESRCH) {
		t.Fatalf("ESRCH reached the caller: %v", err)
	}
}

// AC9 (r2, TestZombieGroupKillNotJoined). kill(-pgid) = ESRCH means the group
// was already empty; it must not be joined. The real kill still runs first.
func TestZombieGroupKillNotJoined(t *testing.T) {
	orig := killGroup
	t.Cleanup(func() { killGroup = orig })
	for name, esrch := range map[string]error{
		"bare":    syscall.ESRCH,
		"wrapped": fmt.Errorf("kill group: %w", syscall.ESRCH),
	} {
		killGroup = func(pgid int) error { _ = orig(pgid); return esrch }
		err := runScripted(t, scriptedInterpreter(t, ":", "wait"), 120*time.Second)
		if _, ok := err.(*OutputLimitError); !ok {
			t.Fatalf("%s: error = %T %v, want exactly *OutputLimitError", name, err, err)
		}
	}
}
