// The session credentials store (row 189 M1a, D-WORLD-85 = C,
// w-world-onboarding §4 C).
//
// session_credentials used to be a table of the world store, so minting a
// session needed the world store's writer lock, which a serving daemon holds:
// every mint meant stopping the daemon. It now lives in its OWN SQLite file
// with its OWN writer lock ("<credentials-db>.writer.lock"):
//
//   - `session new/mint/revoke` open it with OpenCredentials, at the terminal,
//     after the y/N — the mint path's order is unchanged; only the lock it
//     takes is a different file;
//   - a serving daemon never holds it: it resolves each request's bearer token
//     through a fresh OpenCredentialsReadOnly handle (host/authority
//     ReadOnlyAt), so a mint or a revoke is visible on the next request;
//   - `session list` reads it read-only, beside a live daemon or not.
//
// The default location is DefaultCredentialsPath: credentials.db in the world
// store's own directory, so a rig's store and its credentials travel together.
// The world store keeps its single-writer lock for everything else.
package store

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
)

// credentialsSchemaSQL is the credentials store's whole schema.
//
//go:embed credentials_schema.sql
var credentialsSchemaSQL string

// credentialsSchemaVersion is the credentials store's own user_version,
// independent of the world store's currentSchemaVersion.
const credentialsSchemaVersion = 1

// CredentialsFileName is the default credentials store's file name, created in
// the world store's directory.
const CredentialsFileName = "credentials.db"

// credentialsTable is the one application table a credentials store holds.
const credentialsTable = "session_credentials"

// SessionRow is the persisted credential->(episode, grants, expiry) mapping row for
// session_credentials (w-session-authority D2). GrantsJSON is an OPAQUE JSON
// encoding of the []broker.Capability the authority package marshalled; the store
// deliberately does not interpret it (store must not import broker — broker imports
// store), so the field is typed as a plain string and round-trips unchanged.
type SessionRow struct {
	CredentialID string
	EpisodeID    string
	GrantsJSON   string
	ExpiresAt    int64
	CreatedAt    int64
}

// CredentialsSchemaError reports a file that is not a credentials store this
// binary can use: a future or invalid version, foreign tables (a world store
// passed as --credentials-db, say), or a missing session_credentials table.
// The file is never modified.
type CredentialsSchemaError struct {
	Path   string
	Detail string
}

func (e *CredentialsSchemaError) Error() string {
	return fmt.Sprintf("store: credentials store %q: %s; refusing to use it", e.Path, e.Detail)
}

// Credentials is a handle on a credentials store. A writable handle
// (OpenCredentials) holds the store's writer lock until Close; a read-only
// handle (OpenCredentialsReadOnly) holds no lock and refuses writes at the
// SQLite layer.
type Credentials struct {
	db   *sql.DB
	lock *writerLock
}

// DefaultCredentialsPath is the credentials store beside the world store at
// worldDB: <dir of worldDB>/credentials.db. For the default world store
// ~/.ailang/world/world.db that is ~/.ailang/world/credentials.db, the
// location w-world-onboarding §4 names. An in-memory world store has no
// directory, so the caller must name the credentials store itself.
func DefaultCredentialsPath(worldDB string) (string, error) {
	if isInMemoryDSN(worldDB) {
		return "", fmt.Errorf("store: in-memory world store %q has no directory for a credentials store; name one explicitly", worldDB)
	}
	canonical, _, err := resolveDSN(worldDB)
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(canonical), CredentialsFileName), nil
}

// CanonicalPath is the absolute, symlink-free file a file-backed DSN names
// (its query parameters dropped) — the identity its writer lock is keyed by.
func CanonicalPath(dsn string) (string, error) {
	if isInMemoryDSN(dsn) {
		return "", fmt.Errorf("store: %q is an in-memory database with no file", dsn)
	}
	canonical, _, err := resolveDSN(dsn)
	return canonical, err
}

// SameDatabase reports whether two DSNs name the same database file, so a
// caller can refuse a credentials store that IS the world store.
func SameDatabase(a, b string) (bool, error) {
	if isInMemoryDSN(a) || isInMemoryDSN(b) {
		return false, nil
	}
	ca, _, err := resolveDSN(a)
	if err != nil {
		return false, err
	}
	cb, _, err := resolveDSN(b)
	if err != nil {
		return false, err
	}
	return ca == cb, nil
}

// OpenCredentials opens (or creates) the credentials store at path, taking its
// NON-WAITING writer lock first, exactly as Open does for the world store:
// contention returns *WriterAlreadyActive immediately. The lock file is
// "<canonical credentials path>.writer.lock", a different file from the world
// store's, so it never contends with a serving daemon.
func OpenCredentials(path string) (*Credentials, error) {
	if isInMemoryDSN(path) {
		db, err := openCredentialsSQLite(path, path, true)
		if err != nil {
			return nil, err
		}
		return &Credentials{db: db}, nil
	}
	canonical, params, err := resolveDSN(path)
	if err != nil {
		return nil, err
	}
	lock, err := acquireWriterLock(canonical)
	if err != nil {
		return nil, err
	}
	db, err := openCredentialsSQLite(path, writeDSN(canonical, params), true)
	if err != nil {
		_ = lock.release()
		return nil, err
	}
	return &Credentials{db: db, lock: lock}, nil
}

// OpenCredentialsReadOnly opens an EXISTING credentials store read-only. It
// takes no lock and runs no DDL, so it works beside a writer — the daemon's
// per-request resolver and `session list` use it.
func OpenCredentialsReadOnly(path string) (*Credentials, error) {
	if isInMemoryDSN(path) {
		return nil, fmt.Errorf("store: open credentials read-only %q: an in-memory database is per-connection and has no shared contents to read", path)
	}
	canonical, params, err := resolveDSN(path)
	if err != nil {
		return nil, err
	}
	db, err := openCredentialsSQLite(path, readOnlyDSN(canonical, params), false)
	if err != nil {
		return nil, err
	}
	return &Credentials{db: db}, nil
}

// Close closes the database, then releases the writer lock if held.
func (c *Credentials) Close() error {
	err := c.db.Close()
	if c.lock != nil {
		if rerr := c.lock.release(); rerr != nil && err == nil {
			err = rerr
		}
		c.lock = nil
	}
	return err
}

// openCredentialsSQLite opens the driver handle and enforces the credentials
// schema: a fresh file is initialized (only when apply is set), version
// credentialsSchemaVersion must hold exactly the session_credentials table,
// and anything else is refused unmodified.
func openCredentialsSQLite(display, dsn string, apply bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: open credentials %q: %w", display, err)
	}
	db.SetMaxOpenConns(1)
	if err := enforceCredentialsSchema(display, db, apply); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func enforceCredentialsSchema(display string, db *sql.DB, apply bool) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("store: open credentials %q: read schema version: %w", display, err)
	}
	tables, err := applicationTables(db)
	if err != nil {
		return fmt.Errorf("store: open credentials %q: inspect schema: %w", display, err)
	}
	switch {
	case version == 0 && len(tables) == 0:
		if !apply {
			return &CredentialsSchemaError{Path: display, Detail: "it has no credentials schema yet (serve or session new creates it)"}
		}
		return initCredentialsSchema(display, db)
	case version == credentialsSchemaVersion && len(tables) == 1 && tables[0] == credentialsTable:
		return nil
	case version == credentialsSchemaVersion && len(tables) == 0:
		// Its only table was dropped: recreating it empty would silently
		// revoke every credential, so it is refused like commit_objects is.
		return &CredentialsSchemaError{Path: display, Detail: "it is missing table " + credentialsTable}
	case version > credentialsSchemaVersion:
		return &CredentialsSchemaError{Path: display, Detail: fmt.Sprintf("user_version %d is newer than this binary's %d; use a compatible binary", version, credentialsSchemaVersion)}
	default:
		return &CredentialsSchemaError{Path: display, Detail: fmt.Sprintf("user_version %d with tables %v is not a credentials store (is it the world store?)", version, tables)}
	}
}

// initCredentialsSchema applies the schema and stamps the version in one
// transaction, so a crash leaves either a fresh file or a complete store.
func initCredentialsSchema(display string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("store: initialize credentials %q: begin: %w", display, err)
	}
	if _, err := tx.Exec(credentialsSchemaSQL); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("store: initialize credentials %q: %w", display, err)
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", credentialsSchemaVersion)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("store: initialize credentials %q: set version: %w", display, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: initialize credentials %q: commit: %w", display, err)
	}
	return nil
}

// applicationTables lists the non-internal tables of db, sorted by name.
func applicationTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite\_%' ESCAPE '\' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// MintSession inserts one session_credentials row. credential_id is the stored
// sha256(token_hex) — never the raw token (D3). The INSERT is a plain write: a
// collision on the PK (two mints hashing to the same id) errors loudly rather than
// being silently absorbed.
func (c *Credentials) MintSession(ctx context.Context, row SessionRow) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if row.CredentialID == "" || row.EpisodeID == "" {
		return fmt.Errorf("store: mint session: empty credential_id or episode_id")
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: mint session: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?);`,
		row.CredentialID, row.EpisodeID, row.GrantsJSON, row.ExpiresAt, row.CreatedAt,
	); err != nil {
		return fmt.Errorf("store: mint session: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: mint session: commit: %w", err)
	}
	return nil
}

// ResolveSession performs the single indexed point lookup by credential_id (the
// TEXT PRIMARY KEY) that the authority resolver relies on (D3/D5 bounded
// lookup). It returns ok=false when no row exists (unknown/revoked). No scan, no
// LIKE, no fallback table.
func (c *Credentials) ResolveSession(ctx context.Context, credentialID string) (SessionRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return SessionRow{}, false, err
	}
	var row SessionRow
	err := c.db.QueryRowContext(ctx,
		`SELECT credential_id, episode_id, grants_json, expires_at, created_at
		   FROM session_credentials WHERE credential_id = ?;`,
		credentialID,
	).Scan(&row.CredentialID, &row.EpisodeID, &row.GrantsJSON, &row.ExpiresAt, &row.CreatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return SessionRow{}, false, nil
	case err != nil:
		return SessionRow{}, false, fmt.Errorf("store: resolve session: %w", err)
	}
	return row, true, nil
}

// RevokeSession deletes the mapping row for a credential_id in a transaction
// (D4 delete-revokes: a deleted row makes the next resolve of that credential an
// UNKNOWN credential). Deleting a credential_id that is not present is a no-op,
// not an error. The daemon opens a fresh read-only handle per request, so the
// revoke is seen by the next request after this commit.
func (c *Credentials) RevokeSession(ctx context.Context, credentialID string) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: revoke session: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM session_credentials WHERE credential_id = ?;`, credentialID); err != nil {
		return fmt.Errorf("store: revoke session: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: revoke session: commit: %w", err)
	}
	return nil
}

// MaxListSessions caps one ListSessions page (row 138 M5, `session list`).
const MaxListSessions = 500

// listSessionsSQL is the ListSessions query; LegacySessions reuses it on a
// schema-4 world store, whose table has the same columns.
const listSessionsSQL = `SELECT credential_id, episode_id, grants_json, expires_at, created_at
        FROM session_credentials`

// ListSessions returns up to limit session_credentials rows (at most
// MaxListSessions), ordered by created_at then credential_id; a non-empty
// episodeID restricts it to that episode. Rows carry the credential_id HASH
// only — the raw token is never stored (D3). It is a plain read, so it works
// on an OpenCredentialsReadOnly handle beside a live writer.
func (c *Credentials) ListSessions(ctx context.Context, episodeID string, limit int) ([]SessionRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	return listSessions(ctx, c.db, episodeID, limit)
}

func listSessions(ctx context.Context, db *sql.DB, episodeID string, limit int) ([]SessionRow, error) {
	if limit <= 0 || limit > MaxListSessions {
		limit = MaxListSessions
	}
	q := listSessionsSQL
	args := []any{}
	if episodeID != "" {
		q += ` WHERE episode_id = ?`
		args = append(args, episodeID)
	}
	q += ` ORDER BY created_at, credential_id LIMIT ?;`
	args = append(args, limit)
	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list sessions: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []SessionRow
	for rows.Next() {
		var r SessionRow
		if err := rows.Scan(&r.CredentialID, &r.EpisodeID, &r.GrantsJSON, &r.ExpiresAt, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: list sessions: scan: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list sessions: %w", err)
	}
	return out, nil
}
