//go:build !unix

// Non-unix arm of the single-writer lock, present so the package builds on
// every platform Go targets.
//
// HONESTY NOTE: this arm is NEVER exercised by the project's gates. The dev rig
// is darwin/arm64 and CI runs ubuntu — both unix. It is not tested, and it is
// not claimed to be tested. It fails closed rather than silently opening an
// unguarded write handle: a missing lock primitive must look like a refusal,
// not like a granted lock. Wiring a real primitive here (LockFileEx on Windows)
// is a separate, tested change.
package store

import "runtime"

// acquireWriterLock refuses to open a write handle on a platform with no wired
// lock primitive.
func acquireWriterLock(dbPath string) (*writerLock, error) {
	return nil, &UnsupportedPlatformError{GOOS: runtime.GOOS}
}

// release is unreachable on this platform (acquireWriterLock never returns a
// lock) and exists only so the shared call sites compile.
func (l *writerLock) release() error { return nil }
