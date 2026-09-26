package broker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
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

func runScript(t *testing.T, script string, execTimeout time.Duration) error {
	t.Helper()
	_, err := runBounded(context.Background(), handlerBounds{execTimeout: execTimeout, maxOutputBytes: 1024},
		handlerCommand{path: "/bin/sh", args: []string{"-c", script}, dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}})
	return err
}

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

// AC1. Broker returns promptly on overflow even with a direct-child-only kill,
// so the defect is a LEAKED grandchild, observed by its own state.
func TestOverflowKillReachesForkedGrandchild(t *testing.T) {
	dir := t.TempDir()
	pidf, survived := filepath.Join(dir, "gc.pid"), filepath.Join(dir, "survived")
	base := procbound.Outstanding()
	err := runScript(t, fmt.Sprintf("(sleep 10; : > %s) &\necho $! > %s\n%s\nwait", survived, pidf, overflowLoop), 120*time.Second)
	if got := procbound.Outstanding(); got != base {
		t.Fatalf("reservations after a reaped call = %d, want %d (released at reap)", got, base)
	}
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var overflow *HandlerOutputOverflowError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *HandlerOutputOverflowError", err, err)
	}
	if !proctest.DeadWithin(t, gc, time.Second) {
		t.Fatalf("grandchild %d outlived the overflow return", gc)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("cleanup reported incomplete although the child exited: %v", err)
	}
	if _, statErr := os.Stat(survived); statErr == nil {
		t.Fatalf("grandchild %d ran to completion", gc)
	}
}

// AC2.
func TestOverflowKillErrorJoinedBehindTypedError(t *testing.T) {
	injected := errors.New("injected kill failure")
	orig := killGroup
	t.Cleanup(func() { killGroup = orig })
	killGroup = func(pgid int) error { return errors.Join(orig(pgid), injected) }
	err := runScript(t, overflowLoop, 120*time.Second)
	var overflow *HandlerOutputOverflowError
	if !errors.As(err, &overflow) || !errors.Is(err, injected) || !errors.Is(err, ErrHandlerOverflow) {
		t.Fatalf("error = %v, want *HandlerOutputOverflowError joined with the kill failure", err)
	}
}

// AC4. Timeout path (no overflow): an escapee holds the pipe past the group
// kill; the pipe-close bound returns runBounded as a timeout while it is ALIVE.
func TestPipeCloseBoundsEscapedDescendant(t *testing.T) {
	pidf := filepath.Join(t.TempDir(), "gc.pid")
	err := runScript(t, proctest.EscapeeShell(t, pidf)+"\necho small\nwait", 3*time.Second)
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var timeout *HandlerTimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("error = %T %v, want *HandlerTimeoutError", err, err)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("cleanup reported incomplete although the child exited: %v", err)
	}
	if !proctest.Alive(t, gc) {
		t.Fatalf("escapee %d not alive at return: runBounded waited for it", gc)
	}
}

// AC7 on both broker paths: every kill fails (EPERM, nothing signalled) and
// the direct child lives on; the call returns with the typed error and the
// failures discoverable while the child is ALIVE; a guarded cleanup kills it
// and the background waiter reaps it.
func TestFailedKillBoundsCallAndChildIsReapedLater(t *testing.T) {
	for _, tc := range []struct {
		name, script string
		check        func(error) bool
	}{
		{"overflow", overflowLoop + "\nexec sleep 30", func(err error) bool {
			var o *HandlerOutputOverflowError
			return errors.As(err, &o) && errors.Is(err, syscall.EPERM)
		}},
		{"timeout", "exec sleep 30", func(err error) bool {
			var o *HandlerTimeoutError
			return errors.As(err, &o)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pgidOf := recordFailingKill(t)
			base := procbound.Outstanding()
			err := runScript(t, tc.script, time.Second)
			pgid := pgidOf()
			killed := false
			t.Cleanup(func() {
				if !killed {
					_ = proctest.Signal(pgid, syscall.SIGKILL, syscall.Kill)
				}
			})
			if !tc.check(err) || !errors.Is(err, procbound.ErrCleanupIncomplete) {
				t.Fatalf("error = %v, want the typed error + ErrCleanupIncomplete", err)
			}
			if !proctest.Alive(t, pgid) {
				t.Fatalf("direct child %d not alive at return: the call waited for it", pgid)
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
		})
	}
}

// AC8. A full set of reservations refuses a new call before Start.
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
	err := runScript(t, "echo ok", 120*time.Second)
	if !errors.Is(err, procbound.ErrCleanupBacklog) {
		t.Fatalf("error = %v, want ErrCleanupBacklog", err)
	}
}

// A Start failure releases the reservation Admit took.
func TestStartFailureReleasesReservation(t *testing.T) {
	base := procbound.Outstanding()
	_, err := runBounded(context.Background(), handlerBounds{execTimeout: time.Minute, maxOutputBytes: 1024},
		handlerCommand{path: filepath.Join(t.TempDir(), "no-such-handler"), dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}})
	if err == nil || !strings.Contains(err.Error(), "start handler subprocess") {
		t.Fatalf("error = %v, want a Start failure", err)
	}
	if got := procbound.Outstanding(); got != base {
		t.Fatalf("reservations after a failed Start = %d, want %d", got, base)
	}
}

// AC9 (r2). An ordinary overflow, real killGroup, child exits on its own:
// the error is exactly the typed overflow and no ESRCH is anywhere in it.
func TestOrdinaryOverflowCarriesNoESRCH(t *testing.T) {
	err := runScript(t, overflowLoop, 120*time.Second)
	if _, ok := err.(*HandlerOutputOverflowError); !ok {
		t.Fatalf("error = %T %v, want exactly *HandlerOutputOverflowError", err, err)
	}
	if errors.Is(err, syscall.ESRCH) || !errors.Is(err, ErrHandlerOverflow) {
		t.Fatalf("error = %v: ESRCH present or ErrHandlerOverflow lost", err)
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
		err := runScript(t, overflowLoop+"\nwait", 120*time.Second)
		if _, ok := err.(*HandlerOutputOverflowError); !ok {
			t.Fatalf("%s: error = %T %v, want exactly *HandlerOutputOverflowError", name, err, err)
		}
	}
}

// AC10 (r2). More barrier-synchronized callers than the limit, every kill
// failing (EPERM, nothing signalled), every child outliving the call: exactly
// MaxOutstanding are admitted, the rest refused without starting; the count
// never exceeds the limit; the abandoned waiters RETAIN their reservations
// (a further Admit is refused); after a guarded kill each is released exactly
// once, so the count returns to zero and not below.
func TestReservationIsExactUnderConcurrency(t *testing.T) {
	for start := time.Now(); procbound.Outstanding() != 0 && time.Since(start) < 5*time.Second; time.Sleep(time.Millisecond) {
	}
	if got := procbound.Outstanding(); got != 0 {
		t.Fatalf("reservations held at start = %d, want 0", got)
	}
	var mu sync.Mutex
	var pgids []int
	killed := map[int]bool{}
	orig := killGroup
	t.Cleanup(func() {
		killGroup = orig
		mu.Lock()
		defer mu.Unlock()
		for _, pgid := range pgids {
			if !killed[pgid] {
				_ = proctest.Signal(pgid, syscall.SIGKILL, syscall.Kill)
			}
		}
	})
	killGroup = func(pgid int) error {
		mu.Lock()
		defer mu.Unlock()
		pgids = append(pgids, pgid)
		return syscall.EPERM
	}

	var peak atomic.Int64
	stopSampling := make(chan struct{})
	sampled := make(chan struct{})
	go func() {
		defer close(sampled)
		for {
			if n := procbound.Outstanding(); n > peak.Load() {
				peak.Store(n)
			}
			select {
			case <-stopSampling:
				return
			default:
				time.Sleep(50 * time.Microsecond)
			}
		}
	}()

	const callers = procbound.MaxOutstanding + 4
	barrier := make(chan struct{})
	errs := make([]error, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-barrier
			errs[i] = runScript(t, "exec sleep 30", time.Second)
		}(i)
	}
	close(barrier)
	wg.Wait()
	close(stopSampling)
	<-sampled

	var admitted, refused int
	for i, err := range errs {
		var timeout *HandlerTimeoutError
		switch {
		case errors.Is(err, procbound.ErrCleanupBacklog):
			refused++
		case errors.As(err, &timeout) && errors.Is(err, procbound.ErrCleanupIncomplete):
			admitted++
		default:
			t.Fatalf("caller %d: error = %v, want ErrCleanupBacklog or timeout+ErrCleanupIncomplete", i, err)
		}
	}
	if admitted != procbound.MaxOutstanding || refused != callers-procbound.MaxOutstanding {
		t.Fatalf("admitted=%d refused=%d, want %d/%d", admitted, refused, procbound.MaxOutstanding, callers-procbound.MaxOutstanding)
	}
	if p := peak.Load(); p > procbound.MaxOutstanding {
		t.Fatalf("reservations peaked at %d, limit %d", p, procbound.MaxOutstanding)
	}
	if got := procbound.Outstanding(); got != procbound.MaxOutstanding {
		t.Fatalf("reservations after return = %d, want %d retained by the background waiters", got, procbound.MaxOutstanding)
	}
	if release, err := procbound.Admit(); !errors.Is(err, procbound.ErrCleanupBacklog) {
		if release != nil {
			release()
		}
		t.Fatalf("Admit with every slot retained = %v, want ErrCleanupBacklog", err)
	}

	mu.Lock()
	distinct := map[int]bool{}
	for _, pgid := range pgids {
		distinct[pgid] = true
	}
	for pgid := range distinct {
		if err := proctest.Signal(pgid, syscall.SIGKILL, syscall.Kill); err != nil {
			mu.Unlock()
			t.Fatal(err)
		}
		killed[pgid] = true
	}
	mu.Unlock()
	if len(distinct) != procbound.MaxOutstanding {
		t.Fatalf("killGroup saw %d distinct children, want %d (refused callers must not start one)", len(distinct), procbound.MaxOutstanding)
	}
	for start := time.Now(); procbound.Outstanding() > 0 && time.Since(start) < 5*time.Second; time.Sleep(time.Millisecond) {
	}
	time.Sleep(50 * time.Millisecond) // let a double release (if any) land
	if got := procbound.Outstanding(); got != 0 {
		t.Fatalf("reservations after reaping = %d, want exactly 0 (each released once)", got)
	}
}
