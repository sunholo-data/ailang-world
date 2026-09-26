package broker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

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

// AC1. Broker returns promptly on overflow even with a direct-child-only kill,
// so the defect is a LEAKED grandchild, observed by its own state.
func TestOverflowKillReachesForkedGrandchild(t *testing.T) {
	dir := t.TempDir()
	pidf, survived := filepath.Join(dir, "gc.pid"), filepath.Join(dir, "survived")
	err := runScript(t, fmt.Sprintf("(sleep 10; : > %s) &\necho $! > %s\n%s\nwait", survived, pidf, overflowLoop), 120*time.Second)
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var overflow *HandlerOutputOverflowError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *HandlerOutputOverflowError", err, err)
	}
	if !proctest.DeadWithin(t, gc, time.Second) {
		t.Fatalf("grandchild %d outlived the overflow return", gc)
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
	if !proctest.Alive(t, gc) {
		t.Fatalf("escapee %d not alive at return: runBounded waited for it", gc)
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
