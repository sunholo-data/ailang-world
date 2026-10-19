package store

// The one-time move of session credentials out of a schema-4 world store
// (row 189 M1a). Schema 5 is schema 4 without session_credentials; the rows
// move to the credentials store and none stay behind, so there is exactly one
// source of truth.
//
// Crash safety: the copy, the DROP and the version stamp are ONE SQLite
// transaction spanning both files (the credentials store is ATTACHed to the
// world store's connection). In rollback-journal mode — both stores use
// journal_mode=delete, and the move refuses any other mode — SQLite commits a
// multi-file transaction atomically through a super-journal: after a crash
// either the world store is still schema 4 with every row and the credentials
// store holds none of the copies, or the world store is schema 5 and the
// credentials store holds every row. Nothing is lost and nothing is
// duplicated (credential_id is the PRIMARY KEY on both sides, and a re-run of
// an interrupted move starts again from the untouched schema-4 rows).
//
// Who runs it: `serve` at startup (daemon.New), and `session new/mint/revoke`,
// which need the credentials store authoritative before they write it. A
// plain writable Open of a schema-4 store refuses with *PreSplitStoreError
// instead, because only those callers know where the credentials store is
// (--credentials-db); read-only opens of schema 4 are accepted, since the
// world tables are unchanged.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

// preSplitSchemaVersion is the last world schema that held session_credentials.
const preSplitSchemaVersion = 4

// PreSplitSchemaVersion exports preSplitSchemaVersion for the session verbs,
// which read a pre-split store's credentials in place (LegacySessions).
const PreSplitSchemaVersion = preSplitSchemaVersion

// PreSplitStoreError reports a writable Open of a schema-4 world store, whose
// session credentials have not moved to their own store yet. The store is not
// modified.
type PreSplitStoreError struct{ Path string }

func (e *PreSplitStoreError) Error() string {
	return fmt.Sprintf("store: schema version 4: %q still holds its session credentials, which this binary keeps in a separate credentials store; "+
		"start `ailang-worldd serve` with this build once (or run session new/mint/revoke while no daemon serves it) to move them; refusing to open it for writing before then", e.Path)
}

// CredentialsNotMovedError reports that a schema-4 world store's credentials
// could not be moved because another process holds its writer lock — in
// practice a daemon from an older build still serving it.
type CredentialsNotMovedError struct {
	World, Credentials string
	Err                error
}

func (e *CredentialsNotMovedError) Error() string {
	return fmt.Sprintf("store: %q is at schema version 4 and still holds its session credentials; they move to %q once, under its writer lock, "+
		"which another process holds (%v). fix: stop the daemon serving it, then start `ailang-worldd serve` with this build — it moves them at startup", e.World, e.Credentials, e.Err)
}

func (e *CredentialsNotMovedError) Unwrap() error { return e.Err }

// migrateCredentialsHook is a no-op in production. Tests replace it to
// interrupt the move inside its transaction ("copied": rows inserted into the
// credentials store, world table not yet dropped; "dropped": table dropped and
// version stamped, not yet committed).
var migrateCredentialsHook = func(stage string) error { return nil }

// WorldSchemaVersion reads a world store's user_version without a lock and
// without any write. A missing file reports 0. It is how the session verbs
// tell a schema-4 store (credentials not moved yet) from a current one before
// deciding which store to read.
func WorldSchemaVersion(path string) (int, error) {
	if isInMemoryDSN(path) {
		return 0, nil
	}
	canonical, params, err := resolveDSN(path)
	if err != nil {
		return 0, err
	}
	if _, err := os.Stat(canonical); errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	db, err := sql.Open("sqlite", readOnlyDSN(canonical, params))
	if err != nil {
		return 0, fmt.Errorf("store: read schema version of %q: %w", path, err)
	}
	defer func() { _ = db.Close() }()
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return 0, fmt.Errorf("store: read schema version of %q: %w", path, err)
	}
	return version, nil
}

// LegacySessions lists the credentials still held by a schema-4 world store,
// read-only and lock-free, so `session list` works beside a daemon of an
// older build that has not moved them yet. It refuses any other version.
func LegacySessions(ctx context.Context, worldPath, episodeID string, limit int) ([]SessionRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	canonical, params, err := resolveDSN(worldPath)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", readOnlyDSN(canonical, params))
	if err != nil {
		return nil, fmt.Errorf("store: legacy sessions %q: %w", worldPath, err)
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return nil, fmt.Errorf("store: legacy sessions %q: %w", worldPath, err)
	}
	if version != preSplitSchemaVersion {
		return nil, fmt.Errorf("store: legacy sessions %q: user_version %d is not the pre-split schema %d", worldPath, version, preSplitSchemaVersion)
	}
	return listSessions(ctx, db, episodeID, limit)
}

// MigrateCredentials moves every session_credentials row of the schema-4 world
// store at worldPath into the credentials store at credPath, drops the table
// and stamps the world store schema 5, in one transaction (see the file
// comment). It returns how many rows moved.
//
// It is a no-op (0, nil) for an in-memory, missing, fresh or already-current
// world store, and it checks that BEFORE taking any lock, so a session verb
// beside a serving daemon of this build never contends for the world lock.
// For a schema-4 store it takes the world store's writer lock and then the
// credentials store's; a held world lock is *CredentialsNotMovedError.
func MigrateCredentials(ctx context.Context, worldPath, credPath string) (moved int, err error) {
	if err := requireDeadline(ctx); err != nil {
		return 0, err
	}
	if isInMemoryDSN(worldPath) {
		return 0, nil
	}
	same, err := SameDatabase(worldPath, credPath)
	if err != nil {
		return 0, err
	}
	if same {
		return 0, fmt.Errorf("store: the credentials store %q is the world store itself; name a separate file", credPath)
	}
	// A read-only peek that cannot read (a hot journal left by an
	// interrupted move, say) falls through to the locked check, whose
	// writable connection rolls the journal back.
	if version, perr := WorldSchemaVersion(worldPath); perr == nil && version != preSplitSchemaVersion {
		return 0, nil
	}
	canonical, params, err := resolveDSN(worldPath)
	if err != nil {
		return 0, err
	}
	lock, err := acquireWriterLock(canonical)
	if err != nil {
		if IsWriterAlreadyActive(err) {
			return 0, &CredentialsNotMovedError{World: worldPath, Credentials: credPath, Err: err}
		}
		return 0, err
	}
	defer func() { _ = lock.release() }()

	db, err := sql.Open("sqlite", writeDSN(canonical, params))
	if err != nil {
		return 0, fmt.Errorf("store: move credentials: open %q: %w", worldPath, err)
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return 0, fmt.Errorf("store: move credentials: read schema version: %w", err)
	}
	if version != preSplitSchemaVersion {
		return 0, nil // another process moved them between the peek and the lock
	}

	creds, err := OpenCredentials(credPath)
	if err != nil {
		return 0, fmt.Errorf("store: move credentials: %w", err)
	}
	defer func() { _ = creds.Close() }()
	credCanonical, _, err := resolveDSN(credPath)
	if err != nil {
		return 0, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("store: move credentials: %w", err)
	}
	defer func() { _ = conn.Close() }()
	// The attachment lives on this one connection; db.Close (deferred above)
	// discards both, so there is no DETACH to forget.
	if _, err := conn.ExecContext(ctx, `ATTACH DATABASE ? AS creds`, credCanonical); err != nil {
		return 0, fmt.Errorf("store: move credentials: attach %q: %w", credPath, err)
	}
	for _, schema := range []string{"main", "creds"} {
		var mode string
		if err := conn.QueryRowContext(ctx, "PRAGMA "+schema+".journal_mode").Scan(&mode); err != nil {
			return 0, fmt.Errorf("store: move credentials: journal mode: %w", err)
		}
		if !strings.EqualFold(mode, "delete") && !strings.EqualFold(mode, "truncate") && !strings.EqualFold(mode, "persist") {
			return 0, fmt.Errorf("store: move credentials: %s journal_mode is %q; the atomic two-file move needs a rollback journal", schema, mode)
		}
	}
	return moveCredentialRows(ctx, conn, worldPath)
}

// moveCredentialRows is the single transaction: copy, verify, drop, stamp.
func moveCredentialRows(ctx context.Context, conn *sql.Conn, worldPath string) (int, error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("store: move credentials: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var total int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM main.session_credentials`).Scan(&total); err != nil {
		return 0, fmt.Errorf("store: move credentials: count: %w", err)
	}
	const columns = `credential_id, episode_id, grants_json, expires_at, created_at`
	if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO creds.session_credentials (`+columns+`)
		SELECT `+columns+` FROM main.session_credentials`); err != nil {
		return 0, fmt.Errorf("store: move credentials: copy: %w", err)
	}
	// Every world row must now exist in the credentials store field for
	// field. INSERT OR IGNORE kept a pre-existing row with the same id; if it
	// differs in any column the move is refused and rolled back.
	var unmatched int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM main.session_credentials m
		WHERE NOT EXISTS (SELECT 1 FROM creds.session_credentials c
			WHERE c.credential_id = m.credential_id AND c.episode_id = m.episode_id
			  AND c.grants_json = m.grants_json AND c.expires_at = m.expires_at
			  AND c.created_at = m.created_at)`).Scan(&unmatched); err != nil {
		return 0, fmt.Errorf("store: move credentials: verify: %w", err)
	}
	if unmatched != 0 {
		return 0, fmt.Errorf("store: move credentials: %d credential(s) of %q conflict with a different row of the same id in the credentials store; nothing was moved", unmatched, worldPath)
	}
	if err := migrateCredentialsHook("copied"); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `DROP TABLE main.session_credentials`); err != nil {
		return 0, fmt.Errorf("store: move credentials: drop: %w", err)
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA main.user_version = %d", currentSchemaVersion)); err != nil {
		return 0, fmt.Errorf("store: move credentials: stamp version: %w", err)
	}
	if err := migrateCredentialsHook("dropped"); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("store: move credentials: commit: %w", err)
	}
	return total, nil
}
