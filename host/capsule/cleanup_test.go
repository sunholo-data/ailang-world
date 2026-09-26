package capsule

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
