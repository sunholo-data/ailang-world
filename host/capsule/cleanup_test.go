package capsule

import (
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

// AC1. The overflow kill must reach a forked grandchild holding the inherited
// pipes. Observable: the grandchild's own state. It records its pid BEFORE the
// interpreter writes, and writes a survival marker only if its sleep
// completes. ExecTimeout (120 s) is 12x the grandchild's lifetime (10 s), so a
// direct-child-only kill returns after the grandchild finished: marker present.
func TestOverflowKillReachesForkedGrandchild(t *testing.T) {
	dir := t.TempDir()
	pidf, survived := filepath.Join(dir, "gc.pid"), filepath.Join(dir, "survived")
	fx := scriptedInterpreter(t, fmt.Sprintf("(sleep 10; : > %s) &\necho $! > %s", survived, pidf), "wait")
	err := runScripted(t, fx, 120*time.Second)
	gc := proctest.ReadPid(t, pidf)
	proctest.ReapOnCleanup(t, gc)
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError", err, err)
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
