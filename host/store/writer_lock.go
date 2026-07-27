// Writer-lock support for the RATIFIED single-writer process model
// (w-worldd-m2 Decision 2, revision r3 — arm A, ratified attended by Mark,
// 2026-07-27).
//
// This file holds the platform-independent half: the structured contention
// error, DSN classification (the in-memory carve-out), and the canonical
// database identity used by BOTH the SQLite handle and the lock file. The
// OS primitive itself lives in writer_lock_unix.go / writer_lock_other.go.
//
// Shape of the invariant:
//
//   - Open(path) on a FILE-BACKED database takes a NON-WAITING OS-backed
//     exclusive lock on "<canonical-db>.writer.lock" BEFORE any SQLite write
//     handle exists. Contention returns *WriterAlreadyActive immediately; it
//     never waits, never retries, never blocks.
//   - An in-memory database is per-connection and physically unreachable from
//     another process, so there is nothing to exclude: no lock is taken and no
//     lock file is created. This is decided on the resolved DSN BEFORE any
//     canonicalization.
//   - Lock-file EXISTENCE is never ownership. The OS drops ownership when the
//     holding process dies, so a leftover pathname from a killed writer cannot
//     wedge the database — a fresh process reopens the same file and reacquires.
package store

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// writerLockSuffix is appended to the canonical database path to name the lock
// file. The lock is process/OS state, never database state: schema.sql is
// unchanged by the single-writer model.
const writerLockSuffix = ".writer.lock"

// WriterAlreadyActive reports that another process (or another open file
// description in this process) already holds the exclusive writer lock for a
// database. It is returned immediately on contention — Open never waits for the
// incumbent writer to go away.
type WriterAlreadyActive struct {
	// DBPath is the canonical database path whose writer authority is held.
	DBPath string
	// LockPath is the lock file guarding DBPath.
	LockPath string
}

func (e *WriterAlreadyActive) Error() string {
	return fmt.Sprintf(
		"store: another process already holds the writer lock for %q (lock file %q); "+
			"only one writer may be active per database", e.DBPath, e.LockPath)
}

// IsWriterAlreadyActive reports whether err is (or wraps) a
// *WriterAlreadyActive. It mirrors the landed IsConflict idiom.
func IsWriterAlreadyActive(err error) bool {
	var w *WriterAlreadyActive
	return errors.As(err, &w)
}

// UnsupportedPlatformError is returned when the writer lock cannot be taken
// because this build has no OS lock primitive wired up. Refusing to open is
// deliberate: silently skipping the lock would turn the fail-closed guarantee
// into an unenforced comment.
type UnsupportedPlatformError struct {
	GOOS string
}

func (e *UnsupportedPlatformError) Error() string {
	return fmt.Sprintf("store: single-writer lock is not implemented for GOOS=%q; "+
		"refusing to open a write handle without it", e.GOOS)
}

// writerLock is a held OS-level exclusive lock. The lock is a property of the
// OPEN FILE DESCRIPTION, so the *os.File must stay open for the Store's whole
// lifetime — closing it releases the lock.
type writerLock struct {
	dbPath   string
	lockPath string
	file     *os.File
}

// isInMemoryDSN reports whether dsn names a SQLite in-memory database, using the
// same DSN reading as the pinned modernc.org/sqlite driver:
//
//   - ":memory:" (bare, or as the path part of a file: URI: "file::memory:");
//   - any file: URI carrying mode=memory.
//
// Query parameters on a NON-file: DSN are dropped by the driver before it ever
// sees them, so they are dropped here too — otherwise "x.db?mode=memory" would
// be classified in-memory while the driver opened a real, unlocked file.
func isInMemoryDSN(dsn string) bool {
	pathPart, rawQuery, hasQuery := strings.Cut(dsn, "?")
	if strings.HasPrefix(pathPart, "file:") {
		pathPart = strings.TrimPrefix(pathPart, "file:")
	} else {
		hasQuery = false
	}
	if pathPart == ":memory:" {
		return true
	}
	if !hasQuery {
		return false
	}
	q, err := url.ParseQuery(rawQuery)
	return err == nil && q.Get("mode") == "memory"
}

// resolveDSN splits a file-backed DSN into its canonical absolute database path
// and the driver query parameters carried by a file: URI DSN (nil for a plain
// path DSN).
//
// The canonical path is what makes the invariant hold across path spellings:
// two processes reaching one database by different names (relative, symlinked,
// dot-segmented) resolve to the same identity and therefore to the same lock
// file, so they collide as they must.
func resolveDSN(dsn string) (canonical string, params url.Values, err error) {
	pathPart, rawQuery, hasQuery := strings.Cut(dsn, "?")
	if strings.HasPrefix(pathPart, "file:") {
		pathPart = strings.TrimPrefix(pathPart, "file:")
		// "file:///abs" and "file://localhost/abs" both denote /abs.
		pathPart = strings.TrimPrefix(pathPart, "//localhost")
		if strings.HasPrefix(pathPart, "//") {
			pathPart = strings.TrimPrefix(pathPart, "//")
		}
		if unescaped, uerr := url.PathUnescape(pathPart); uerr == nil {
			pathPart = unescaped
		}
		if hasQuery {
			params, err = url.ParseQuery(rawQuery)
			if err != nil {
				return "", nil, fmt.Errorf("store: parse DSN query %q: %w", rawQuery, err)
			}
		}
	}
	if pathPart == "" {
		return "", nil, fmt.Errorf("store: empty database path in DSN %q", dsn)
	}
	canonical, err = canonicalDBPath(pathPart)
	if err != nil {
		return "", nil, err
	}
	return canonical, params, nil
}

// canonicalDBPath resolves path to an absolute, symlink-free database identity.
// The target usually does not exist yet on first open, so symlinks are resolved
// on the PARENT directory and the base name is rejoined afterwards.
func canonicalDBPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("store: resolve %q: %w", path, err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("store: resolve %q: %w", path, err)
	}
	// Not created yet: canonicalize the directory that will hold it.
	dir, base := filepath.Split(abs)
	resolvedDir, derr := filepath.EvalSymlinks(filepath.Clean(dir))
	if derr != nil {
		return "", fmt.Errorf("store: resolve parent of %q: %w", path, derr)
	}
	return filepath.Join(resolvedDir, base), nil
}

// writeDSN renders the DSN handed to the driver for a WRITE handle. A plain
// path DSN stays a plain path (byte-identical driver behaviour to M1, modulo
// canonicalization); a URI DSN keeps its parameters.
func writeDSN(canonical string, params url.Values) string {
	if len(params) == 0 {
		return canonical
	}
	return fileURI(canonical, params)
}

// readOnlyDSN renders the DSN handed to the driver for a READ-ONLY handle. The
// driver always passes SQLITE_OPEN_READWRITE|CREATE|URI, and SQLite's mode=
// parameter may only RESTRICT those flags, never elevate them — so mode=ro
// yields a genuinely read-only connection.
func readOnlyDSN(canonical string, params url.Values) string {
	q := url.Values{}
	for k, v := range params {
		q[k] = append([]string(nil), v...)
	}
	q.Set("mode", "ro")
	return fileURI(canonical, q)
}

// fileURI builds a "file://<path>?<params>" DSN with proper escaping.
func fileURI(canonical string, params url.Values) string {
	u := url.URL{Scheme: "file", Path: canonical}
	if len(params) > 0 {
		u.RawQuery = params.Encode()
	}
	return u.String()
}
