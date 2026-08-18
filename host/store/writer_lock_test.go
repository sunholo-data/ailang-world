package store

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// This file is the load-bearing artifact of the RATIFIED single-writer decision
// (w-worldd-m2 Decision 2 r3 — arm A, ratified attended by Mark, 2026-07-27).
//
// The proof is CROSS-PROCESS by construction. A subprocess helper — a real
// second OS process, re-exec'd from this very test binary — opens the writer
// and only then signals readiness; the parent then attempts the same database
// and must be refused. An in-process test (two Opens in one goroutine, a
// sync.Mutex, a package-level map of held paths) would NOT prove the ratified
// property, and the subprocess shape is chosen precisely because such an
// implementation is structurally incapable of passing it.
//
// Each case carries its own anti-vacuity hook, because a lock test is easy to
// write in a way that cannot fail:
//
//   - the parent FAILS LOUDLY if READY never arrives (a crashed helper would
//     otherwise make "Open was refused" trivially true);
//   - a NEGATIVE CONTROL re-Opens after the helper exits, so an implementation
//     that always errors cannot look green;
//   - the read-only case asserts the exact VALUE the writer wrote, not merely
//     the absence of an error;
//   - the crash case asserts the stale lock FILE is still on disk, so it is
//     genuinely exercising the existence-is-not-ownership path.

const (
	// helperDBEnv names the database the subprocess helper should open as writer.
	helperDBEnv = "WORLDD_LOCK_HELPER_DB"
	// helperSeedEnv, when set, is a payload the helper writes BEFORE signalling
	// readiness, so a concurrent reader has a known value to read back.
	helperSeedEnv = "WORLDD_LOCK_HELPER_SEED"
	// helperReady is printed on stdout only after the writer is genuinely held.
	helperReady = "READY"
	// helperSemanticID labels the seeded object.
	helperSemanticID = "helper-seed"
	// readyTimeout bounds how long the parent waits for the helper to come up.
	readyTimeout = 60 * time.Second
	// nonWaitingBudget is the ratified fail-closed budget: a contended Open must
	// return essentially immediately, not block-then-fail.
	nonWaitingBudget = 2 * time.Second
)

// TestWriterLockHelperProcess is not a test of this package: it is the body of
// the subprocess used by the cross-process cases below. It runs only when
// re-exec'd with WORLDD_LOCK_HELPER_DB set, and it never returns to the test
// framework — it exits explicitly so its stdout carries nothing but READY.
func TestWriterLockHelperProcess(t *testing.T) {
	dbPath := os.Getenv(helperDBEnv)
	if dbPath == "" {
		t.Skip("subprocess helper body for the cross-process writer-lock proof; " +
			"runs only when re-exec'd with " + helperDBEnv)
	}
	s, err := Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "helper: Open(%q) failed: %v\n", dbPath, err)
		os.Exit(3)
	}
	if seed := os.Getenv(helperSeedEnv); seed != "" {
		if err := s.PutObject(obj(seed, helperSemanticID)); err != nil {
			fmt.Fprintf(os.Stderr, "helper: seed write failed: %v\n", err)
			_ = s.Close()
			os.Exit(4)
		}
	}
	// Readiness is signalled ONLY on a successful Open: the parent may assume
	// that after READY the writer lock is genuinely held by this process.
	fmt.Println(helperReady)
	// Hold the writer until the parent closes our stdin (or SIGKILLs us).
	_, _ = io.Copy(io.Discard, os.Stdin)
	_ = s.Close()
	os.Exit(0)
}

// syncBuffer is a concurrency-safe sink for the helper's stderr; os/exec writes
// to it from its own goroutine while the test reads it for diagnostics.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

// helperProc is a running writer-holding subprocess.
type helperProc struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stderr *syncBuffer
	waited bool
}

// startWriterHelper re-execs this test binary as a SECOND OS PROCESS that opens
// dbPath as writer, and returns only after that process has signalled READY.
// Failure to see READY is a hard failure: without it, every "the writer lock
// refused me" assertion below would be vacuous.
func startWriterHelper(t *testing.T, dbPath string, extraEnv ...string) *helperProc {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestWriterLockHelperProcess$")
	cmd.Env = append(os.Environ(), helperDBEnv+"="+dbPath)
	cmd.Env = append(cmd.Env, extraEnv...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("helper stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("helper stdout pipe: %v", err)
	}
	stderr := &syncBuffer{}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}

	h := &helperProc{cmd: cmd, stdin: stdin, stderr: stderr}
	t.Cleanup(func() { h.kill(t) })

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		line, rerr := bufio.NewReader(stdout).ReadString('\n')
		if rerr != nil {
			errCh <- rerr
			return
		}
		lineCh <- strings.TrimSpace(line)
	}()

	select {
	case line := <-lineCh:
		if line != helperReady {
			t.Fatalf("helper first stdout line = %q, want %q; stderr:\n%s",
				line, helperReady, stderr.String())
		}
	case rerr := <-errCh:
		t.Fatalf("helper never signalled %s (stdout closed: %v) — it failed to take the "+
			"writer lock, so any contention assertion would be VACUOUS; stderr:\n%s",
			helperReady, rerr, stderr.String())
	case <-time.After(readyTimeout):
		t.Fatalf("helper never signalled %s within %v — refusing to assert contention "+
			"against a helper that may not hold the writer; stderr:\n%s",
			helperReady, readyTimeout, stderr.String())
	}
	return h
}

// stop closes the helper's stdin and waits for a clean exit.
func (h *helperProc) stop(t *testing.T) {
	t.Helper()
	if h.waited {
		return
	}
	_ = h.stdin.Close()
	if err := h.cmd.Wait(); err != nil {
		t.Fatalf("helper exited with error: %v; stderr:\n%s", err, h.stderr.String())
	}
	h.waited = true
}

// kill SIGKILLs the helper so no deferred cleanup runs, then reaps it. After
// Wait returns, the OS has closed the helper's descriptors and therefore
// dropped its flock ownership.
func (h *helperProc) kill(t *testing.T) {
	t.Helper()
	if h.waited {
		return
	}
	_ = h.cmd.Process.Kill()
	_ = h.stdin.Close()
	_, _ = h.cmd.Process.Wait()
	_ = h.cmd.Wait()
	h.waited = true
}

// TestWriterLockCrossProcessContention is the ratified proof: while a SECOND
// REAL OS PROCESS holds the writer, Open on the same canonical path fails
// immediately with the structured *WriterAlreadyActive — and succeeds again
// once that process is gone.
func TestWriterLockCrossProcessContention(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	h := startWriterHelper(t, dbPath)

	start := time.Now()
	s, err := Open(dbPath)
	elapsed := time.Since(start)
	if err == nil {
		_ = s.Close()
		t.Fatalf("Open(%q) SUCCEEDED while OS process %d holds the writer lock: "+
			"the ratified single-writer invariant is not enforced",
			dbPath, h.cmd.Process.Pid)
	}
	var active *WriterAlreadyActive
	if !errors.As(err, &active) {
		t.Fatalf("contended Open returned %T (%v), want *store.WriterAlreadyActive", err, err)
	}
	if !IsWriterAlreadyActive(err) {
		t.Fatalf("IsWriterAlreadyActive(%v) = false, want true", err)
	}
	if active.DBPath == "" || active.LockPath == "" {
		t.Fatalf("structured error carries empty paths: %+v", active)
	}
	if want := active.DBPath + writerLockSuffix; active.LockPath != want {
		t.Fatalf("LockPath = %q, want %q", active.LockPath, want)
	}
	// NON-WAITING: fail closed means refuse now, not refuse eventually.
	if elapsed >= nonWaitingBudget {
		t.Fatalf("contended Open took %v (budget %v): the ratified semantics are "+
			"NON-WAITING, not block-then-fail", elapsed, nonWaitingBudget)
	}
	if _, serr := os.Stat(active.LockPath); serr != nil {
		t.Fatalf("lock file %q missing while the writer is held: %v", active.LockPath, serr)
	}

	// NEGATIVE CONTROL — without this, an implementation that always returned
	// WriterAlreadyActive would pass everything above.
	h.stop(t)
	after, err := Open(dbPath)
	if err != nil {
		t.Fatalf("negative control FAILED: Open(%q) after the writer process exited: %v",
			dbPath, err)
	}
	if err := after.Close(); err != nil {
		t.Fatalf("Close after reacquiring the writer: %v", err)
	}
}

// TestWriterLockReadOnlyConcurrent proves OpenReadOnly works while another
// process holds the writer, reads the exact value that process wrote, is
// genuinely read-only, and takes no writer lock of its own.
func TestWriterLockReadOnlyConcurrent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	const seed = "row-written-by-the-holding-writer-process"
	h := startWriterHelper(t, dbPath, helperSeedEnv+"="+seed)

	ro, err := OpenReadOnly(dbPath)
	if err != nil {
		t.Fatalf("OpenReadOnly(%q) while the writer is held by another process: %v", dbPath, err)
	}

	// Assert the VALUE, not merely the absence of an error.
	want := obj(seed, helperSemanticID)
	got, ok, err := ro.GetObject(context.Background(), want.Hash)
	if err != nil {
		t.Fatalf("read-only GetObject: %v", err)
	}
	if !ok {
		t.Fatalf("read-only handle did not see the object %q written by the writer process",
			want.Hash.String())
	}
	if string(got.Payload) != seed {
		t.Fatalf("read-only payload = %q, want %q", string(got.Payload), seed)
	}

	// It must really be read-only, or "no writer lock" would be unsound.
	if err := ro.PutObject(obj("written-through-the-read-only-handle", "should-fail")); err == nil {
		t.Fatal("read-only handle accepted a write; mode=ro is not in effect")
	}

	// The read-only handle holds NO writer lock: with it still open, a writer
	// must be able to take the lock the moment the holding process exits.
	h.stop(t)
	w, err := Open(dbPath)
	if err != nil {
		t.Fatalf("an open read-only handle blocked a writer, so OpenReadOnly took the "+
			"writer lock: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close writer: %v", err)
	}
	if err := ro.Close(); err != nil {
		t.Fatalf("Close read-only handle: %v", err)
	}
}

// TestWriterLockCrashRecovery proves a SIGKILLed writer — no deferred cleanup,
// no unlock, no file removal — cannot wedge the database. Lock-file EXISTENCE
// is never ownership.
func TestWriterLockCrashRecovery(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	lockPath := dbPath + writerLockSuffix
	h := startWriterHelper(t, dbPath)

	h.kill(t) // SIGKILL: nothing in the helper gets a chance to tidy up.

	// Assert the stale pathname really is on disk, otherwise this test would be
	// silently exercising the easy case instead of the recovery case.
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("lock file %q is absent after SIGKILL, so this test is not exercising "+
			"the stale-lock-file case: %v", lockPath, err)
	}

	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("a leftover lock FILE from a SIGKILLed writer wedged the database: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close after crash recovery: %v", err)
	}
}

// TestInMemoryOpenTakesNoLock proves the ratified in-memory carve-out: an
// in-memory database takes no lock and leaves no lock file anywhere, and two
// in-memory handles coexist. The two landed :memory: call sites
// (host/store/store_test.go, host/registry/registry_test.go) depend on this.
func TestInMemoryOpenTakesNoLock(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir) // any relative lock file would land here

	for _, dsn := range []string{":memory:", "file::memory:", "file:named-mem?mode=memory"} {
		if !isInMemoryDSN(dsn) {
			t.Fatalf("isInMemoryDSN(%q) = false, want true", dsn)
		}
		s, err := Open(dsn)
		if err != nil {
			t.Fatalf("Open(%q): %v", dsn, err)
		}
		if s.lock != nil {
			t.Fatalf("Open(%q) took a writer lock (%q); in-memory databases must take none",
				dsn, s.lock.lockPath)
		}
		if err := s.Close(); err != nil {
			t.Fatalf("Close(%q): %v", dsn, err)
		}
	}

	// Two concurrent in-memory handles must BOTH succeed.
	a, err := Open(":memory:")
	if err != nil {
		t.Fatalf("first Open(:memory:): %v", err)
	}
	defer func() { _ = a.Close() }()
	b, err := Open(":memory:")
	if err != nil {
		t.Fatalf("second concurrent Open(:memory:): %v — the in-memory carve-out is broken "+
			"and the landed :memory: call sites would go red", err)
	}
	defer func() { _ = b.Close() }()

	// No lock file anywhere under the working directory.
	var found []string
	if err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), writerLockSuffix) {
			found = append(found, p)
		}
		return nil
	}); err != nil {
		t.Fatalf("walk %q: %v", dir, err)
	}
	if len(found) != 0 {
		t.Fatalf("in-memory Opens created lock files: %v", found)
	}
}

// TestCanonicalDBPathCollapsesSpellings asserts the canonicalization itself:
// different spellings of one database — relative, dot-segmented, through a
// symlinked directory — must resolve to a single identity, both before and
// after the database file exists.
//
// This is the case that is actually SENSITIVE to canonicalDBPath. The
// cross-process alias test below is an end-to-end property check, but it cannot
// isolate this code: the kernel resolves symlinks during open(2), so an aliased
// spelling reaches the same inode and therefore the same flock even if the
// canonicalization were removed. Stating that plainly is more useful than
// claiming a stronger proof than the test delivers.
func TestCanonicalDBPathCollapsesSpellings(t *testing.T) {
	root := t.TempDir()
	realDir := filepath.Join(root, "real")
	if err := os.Mkdir(realDir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	linkDir := filepath.Join(root, "link")
	if err := os.Symlink(realDir, linkDir); err != nil {
		t.Skipf("symlinks unavailable on this filesystem: %v", err)
	}
	dbPath := filepath.Join(realDir, "world.db")

	spellings := []string{
		dbPath,
		filepath.Join(linkDir, "world.db"),
		filepath.Join(linkDir, "sub", "..", "world.db"),
		filepath.Join(root, "real", ".", "world.db"),
	}

	// Before the database exists: canonicalization runs on the PARENT.
	want, err := canonicalDBPath(spellings[0])
	if err != nil {
		t.Fatalf("canonicalDBPath(%q) before creation: %v", spellings[0], err)
	}
	if want == filepath.Join(linkDir, "world.db") {
		t.Fatalf("canonical identity %q still names the symlinked directory", want)
	}
	for _, spelling := range spellings[1:] {
		got, err := canonicalDBPath(spelling)
		if err != nil {
			t.Fatalf("canonicalDBPath(%q): %v", spelling, err)
		}
		if got != want {
			t.Fatalf("canonicalDBPath(%q) = %q, want %q (one database, one identity)",
				spelling, got, want)
		}
	}

	// After creation: canonicalization runs on the TARGET.
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	for _, spelling := range spellings {
		got, err := canonicalDBPath(spelling)
		if err != nil {
			t.Fatalf("canonicalDBPath(%q) after creation: %v", spelling, err)
		}
		if got != want {
			t.Fatalf("canonicalDBPath(%q) after creation = %q, want %q", spelling, got, want)
		}
	}

	// A relative spelling from inside the directory must agree too.
	t.Chdir(linkDir)
	got, err := canonicalDBPath("world.db")
	if err != nil {
		t.Fatalf("canonicalDBPath(relative): %v", err)
	}
	if got != want {
		t.Fatalf("canonicalDBPath(relative) = %q, want %q", got, want)
	}
}

// TestWriterLockCanonicalIdentityAcrossSpellings is the end-to-end companion:
// a second OS process reaching one database through a symlinked, dot-segmented
// spelling is still refused. (Inode identity does most of the work here — see
// TestCanonicalDBPathCollapsesSpellings for the case that isolates the
// canonicalization code itself.)
func TestWriterLockCanonicalIdentityAcrossSpellings(t *testing.T) {
	root := t.TempDir()
	realDir := filepath.Join(root, "real")
	if err := os.Mkdir(realDir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	linkDir := filepath.Join(root, "link")
	if err := os.Symlink(realDir, linkDir); err != nil {
		t.Skipf("symlinks unavailable on this filesystem: %v", err)
	}
	dbPath := filepath.Join(realDir, "world.db")
	alias := filepath.Join(linkDir, ".", "sub", "..", "world.db")

	h := startWriterHelper(t, dbPath)

	s, err := Open(alias)
	if err == nil {
		_ = s.Close()
		t.Fatalf("Open(%q) succeeded while another process holds %q: the lock is keyed by "+
			"path SPELLING, not by canonical database identity", alias, dbPath)
	}
	if !IsWriterAlreadyActive(err) {
		t.Fatalf("Open via an aliased spelling returned %v, want *store.WriterAlreadyActive", err)
	}

	// Negative control for the aliased spelling too.
	h.stop(t)
	after, err := Open(alias)
	if err != nil {
		t.Fatalf("negative control: Open(%q) after the writer exited: %v", alias, err)
	}
	if err := after.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestWriterLockReleasedOnClose covers release-on-Close. The RATIFIED
// contention proof is TestWriterLockCrossProcessContention above; this case is
// deliberately in-process and only asserts that Close hands the lock back.
// (flock ownership is per open file description, so the same-process conflict
// asserted here is genuine OS behaviour, not a Go-level mutex.)
func TestWriterLockReleasedOnClose(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")

	first, err := Open(dbPath)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	if first.lock == nil {
		t.Fatal("file-backed Open took no writer lock")
	}

	if second, err := Open(dbPath); err == nil {
		_ = second.Close()
		t.Fatal("a second Open succeeded while the first handle was still open")
	} else if !IsWriterAlreadyActive(err) {
		t.Fatalf("second Open returned %v, want *store.WriterAlreadyActive", err)
	}

	if err := first.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	third, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open after Close did not reacquire the writer lock: %v", err)
	}
	if err := third.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestOpenReleasesWriterLockOnSQLiteFailure proves the lock is handed back when
// SQLite fails AFTER acquisition — otherwise a failed Open would strand writer
// authority for the life of the process.
func TestOpenReleasesWriterLockOnSQLiteFailure(t *testing.T) {
	// A directory is a resolvable path that SQLite cannot open as a database.
	bad := filepath.Join(t.TempDir(), "not-a-database")
	if err := os.Mkdir(bad, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	first, err := Open(bad)
	if err == nil {
		_ = first.Close()
		t.Fatalf("Open(%q) unexpectedly succeeded on a directory", bad)
	}
	if IsWriterAlreadyActive(err) {
		t.Fatalf("first Open reported contention rather than a SQLite failure: %v", err)
	}
	if _, serr := os.Stat(bad + writerLockSuffix); serr != nil {
		t.Fatalf("the lock file was never created, so this test is not exercising the "+
			"release-after-acquisition path: %v", serr)
	}

	second, err := Open(bad)
	if err == nil {
		_ = second.Close()
		t.Fatalf("Open(%q) unexpectedly succeeded on retry", bad)
	}
	if IsWriterAlreadyActive(err) {
		t.Fatalf("retry reported *WriterAlreadyActive: the failed Open leaked the writer "+
			"lock instead of releasing it: %v", err)
	}
}

// TestOpenReadOnlyRejectsMissingAndInMemory pins the two OpenReadOnly refusals:
// it never creates a database, and it refuses an in-memory DSN outright rather
// than handing back a silently empty store.
func TestOpenReadOnlyRejectsMissingAndInMemory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent.db")
	if s, err := OpenReadOnly(missing); err == nil {
		_ = s.Close()
		t.Fatalf("OpenReadOnly(%q) succeeded on a missing database", missing)
	}
	if _, err := os.Stat(missing); err == nil {
		t.Fatalf("OpenReadOnly created %q; a read-only handle must never create a database", missing)
	}
	if _, err := os.Stat(missing + writerLockSuffix); err == nil {
		t.Fatalf("OpenReadOnly created a writer lock file for %q", missing)
	}

	for _, dsn := range []string{":memory:", "file::memory:"} {
		if s, err := OpenReadOnly(dsn); err == nil {
			_ = s.Close()
			t.Fatalf("OpenReadOnly(%q) succeeded; an in-memory database has no shared "+
				"contents to read", dsn)
		}
	}
}

// TestIsInMemoryDSN pins the classification the carve-out rests on, including
// the arm that matters most: query parameters on a NON-file: DSN are dropped by
// the driver, so "x.db?mode=memory" is a real FILE and must be locked.
func TestIsInMemoryDSN(t *testing.T) {
	cases := []struct {
		dsn  string
		want bool
	}{
		{":memory:", true},
		{"file::memory:", true},
		{"file::memory:?cache=shared", true},
		{"file:anything?mode=memory", true},
		{"file:/tmp/world.db?mode=memory&cache=shared", true},
		{":memory:?_pragma=busy_timeout(1000)", true},
		{"world.db", false},
		{"/var/tmp/world.db", false},
		{"file:/var/tmp/world.db", false},
		{"file:/var/tmp/world.db?mode=ro", false},
		{"world.db?mode=memory", false}, // driver drops the query: a real file
		{"memory", false},
	}
	for _, c := range cases {
		if got := isInMemoryDSN(c.dsn); got != c.want {
			t.Errorf("isInMemoryDSN(%q) = %v, want %v", c.dsn, got, c.want)
		}
	}
}
