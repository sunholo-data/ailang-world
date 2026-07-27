//go:build unix

// Unix implementation of the RATIFIED single-writer lock (w-worldd-m2
// Decision 2, arm A).
//
// The primitive is syscall.Flock(LOCK_EX|LOCK_NB) from the Go standard
// library. stdlib was chosen over golang.org/x/sys deliberately: it adds no
// module dependency, so the M2.B daemon dependency-allowlist test stays a
// simple statement about stdlib + this repo + the already-pinned
// modernc.org/sqlite chain.
//
// Two properties of flock are load-bearing here:
//
//   - Ownership belongs to the OPEN FILE DESCRIPTION, not to the pathname and
//     not to the process name. The kernel releases it when the last descriptor
//     for that description closes — including on SIGKILL, where no user-space
//     cleanup runs. That is exactly why a stale lock FILE can never wedge the
//     database.
//   - Because ownership is per description rather than per process, two
//     separate open() calls conflict even inside a single process. Open must
//     therefore release the lock if the SQLite open fails after acquisition, or
//     an in-process retry would deadlock against itself.
package store

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// acquireWriterLock takes the non-waiting exclusive writer lock for the
// canonical database path dbPath. On contention it returns
// *WriterAlreadyActive IMMEDIATELY — it never retries and never blocks.
func acquireWriterLock(dbPath string) (*writerLock, error) {
	lockPath := dbPath + writerLockSuffix
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("store: open writer lock file %q: %w", lockPath, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
			return nil, &WriterAlreadyActive{DBPath: dbPath, LockPath: lockPath}
		}
		return nil, fmt.Errorf("store: acquire writer lock %q: %w", lockPath, err)
	}
	return &writerLock{dbPath: dbPath, lockPath: lockPath, file: f}, nil
}

// release drops the OS lock and closes the descriptor. The lock FILE is left in
// place on purpose: unlinking it would race another process that has already
// opened it, and its mere existence is never treated as ownership.
func (l *writerLock) release() error {
	if l == nil || l.file == nil {
		return nil
	}
	unlockErr := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	closeErr := l.file.Close()
	l.file = nil
	if unlockErr != nil {
		return fmt.Errorf("store: release writer lock %q: %w", l.lockPath, unlockErr)
	}
	if closeErr != nil {
		return fmt.Errorf("store: close writer lock %q: %w", l.lockPath, closeErr)
	}
	return nil
}
