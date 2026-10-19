//go:build unix

package main

import (
	"bytes"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// preSplitRow turns the fresh store at db into a schema-4 store (the v3/v4
// session_credentials table, relabelled 4) holding one credential for
// episode, and returns its credential_id.
func preSplitRow(t *testing.T, db, episode string) string {
	t.Helper()
	raw, err := sql.Open("sqlite", db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = raw.Close() }()
	id := strings.Repeat("0", 63) + "1"
	for _, stmt := range []string{
		`CREATE TABLE session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL
)`,
		`INSERT INTO session_credentials VALUES ('` + id + `', '` + episode + `', '[]', 4000000000, 1)`,
		`PRAGMA user_version = 4`,
	} {
		if _, err := raw.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	return id
}

// holdWorldLock takes db's world writer lock the way a daemon of an older
// build does (flock on "<db>.writer.lock") — store.Open would refuse a
// schema-4 store — and returns its release.
func holdWorldLock(t *testing.T, db string) func() {
	t.Helper()
	canonical, err := store.CanonicalPath(db)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(canonical+".writer.lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	released := false
	release := func() {
		if !released {
			released = true
			_ = f.Close()
		}
	}
	t.Cleanup(release)
	return release
}

// A schema-4 world store (from before the split) served by a daemon of an
// older build: session new names the one-time move and its fix, and
// provisions nothing; session list reads the old credentials in place,
// read-only; with the lock free session new moves them and mints beside them.
func TestSessionVerbsOnAPreSplitStore(t *testing.T) {
	f := newEpisodeFixture(t)
	old := preSplitRow(t, f.db, "ep-old")
	held := holdWorldLock(t, f.db)
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitFatal || !strings.Contains(stderr.String(), "schema version 4") || !strings.Contains(stderr.String(), "start `ailang-worldd serve` with this build") {
		t.Fatalf("pre-split + held lock: rc %d stderr %q", code, stderr.String())
	}
	if _, err := os.Lstat(filepath.Join(f.ws, "ep1")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a worktree was created before the refusal")
	}
	var out, errw bytes.Buffer
	if code := runSessionList([]string{"--db", f.db}, &out, &errw, testMintEnv("", 1000)); code != exitOK || !strings.Contains(out.String(), old) {
		t.Fatalf("list of a pre-split store beside its daemon: rc %d out %q err %q", code, out.String(), errw.String())
	}
	held()
	stdout.Reset()
	stderr.Reset()
	if code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", 1000)); code != exitOK {
		t.Fatalf("pre-split, lock free: rc %d stderr %q", code, stderr.String())
	}
	rows := f.rows(t)
	if len(rows) != 2 || rows[0].CredentialID != old {
		t.Fatalf("rows after the move and the mint = %+v", rows)
	}
	if v, err := store.WorldSchemaVersion(f.db); err != nil || v != 5 {
		t.Fatalf("world schema after the move = %d %v", v, err)
	}
}
