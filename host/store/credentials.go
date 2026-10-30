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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
)

// credentialsSchemaSQL is the credentials store's whole schema.
//
//go:embed credentials_schema.sql
var credentialsSchemaSQL string

// credentialsSchemaVersion is the credentials store's own user_version,
// independent of the world store's currentSchemaVersion. Version 2 added
// profile bindings; version 3 adds immutable attended episode start history.
const credentialsSchemaVersion = 3

// CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
// store's version with it to decide whether to upgrade.
const CredentialsSchemaVersion = credentialsSchemaVersion

// credentialsSchemaV1 is the version before the profile binding: the five
// columns only. It stays readable and writable (unbound mints) so a `session`
// verb of this build never breaks a serving daemon of an older build, which
// refuses a newer user_version; only `serve` of this build upgrades it
// (UpgradeCredentials).
const credentialsSchemaV1 = 1

// credentialsUpgradeV2 is the version 1 -> 2 upgrade: two nullable columns, so
// every existing row stays byte-identical and is unbound (NULL, NULL). The
// column definitions match credentials_schema.sql exactly.
var credentialsUpgradeV2 = []string{
	`ALTER TABLE session_credentials ADD COLUMN project TEXT CHECK (project IS NULL OR project <> '')`,
	`ALTER TABLE session_credentials ADD COLUMN profile_digest TEXT CHECK ((profile_digest IS NULL) = (project IS NULL) AND (profile_digest IS NULL OR profile_digest <> ''))`,
}

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
//
// Project and ProfileDigest are the exec profile binding (row 189 M2): both
// empty for an unbound credential (stored as NULL), both set for a bound one.
type SessionRow struct {
	CredentialID  string
	EpisodeID     string
	GrantsJSON    string
	ExpiresAt     int64
	CreatedAt     int64
	Project       string
	ProfileDigest string
}

// Bound reports whether the credential carries an exec profile binding.
func (r SessionRow) Bound() bool { return r.Project != "" }

// UnboundCredentialsStoreError refuses a bound mint into a version-1
// credentials store, which has no binding columns yet. Only `serve` of a
// build with the binding upgrades it, so a serving daemon of an older build
// is never handed a store it refuses.
type UnboundCredentialsStoreError struct{ Path string }

func (e *UnboundCredentialsStoreError) Error() string {
	return fmt.Sprintf("store: credentials store %q has an older schema; "+
		"restart `ailang-worldd serve` with this build once (it upgrades the store at startup), then mint again", e.Path)
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
	// version is the schema version the handle was opened at: 1 (no binding
	// columns) or credentialsSchemaVersion.
	version int
	display string
	// upgradedFrom is the version UpgradeCredentials upgraded from (0: none).
	upgradedFrom int
}

// Version is the credentials store's schema version as opened.
func (c *Credentials) Version() int { return c.version }

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
//
// A fresh file is created at credentialsSchemaVersion. An existing version-1
// store is opened as it is, never upgraded here: unbound mints and revokes
// work on it, and a bound mint is refused (UnboundCredentialsStoreError).
func OpenCredentials(path string) (*Credentials, error) {
	return openCredentialsWritable(path, false)
}

// UpgradeCredentials brings an existing credentials store to
// credentialsSchemaVersion under its writer lock, in one transaction: the
// version 1 -> 2 upgrade adds the binding columns, every existing row unchanged
// and unbound. It reports whether it upgraded; a current store is left
// untouched (idempotent), and a missing one is created current. `serve`
// calls it at startup — after that, a daemon of an older build would refuse
// the store, which is why the session verbs never upgrade it.
func UpgradeCredentials(path string) (upgraded bool, err error) {
	c, err := openCredentialsWritable(path, true)
	if err != nil {
		return false, err
	}
	upgraded = c.upgradedFrom != 0
	return upgraded, c.Close()
}

func openCredentialsWritable(path string, upgrade bool) (*Credentials, error) {
	mode := credOpenWrite
	if upgrade {
		mode = credOpenUpgrade
	}
	if isInMemoryDSN(path) {
		c, err := openCredentialsSQLite(path, path, mode)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	canonical, params, err := resolveDSN(path)
	if err != nil {
		return nil, err
	}
	lock, err := acquireWriterLock(canonical)
	if err != nil {
		return nil, err
	}
	c, err := openCredentialsSQLite(path, writeDSN(canonical, params), mode)
	if err != nil {
		_ = lock.release()
		return nil, err
	}
	c.lock = lock
	return c, nil
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
	return openCredentialsSQLite(path, readOnlyDSN(canonical, params), credOpenRead)
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

// credOpenMode is what an open may do to the schema.
type credOpenMode int

const (
	credOpenRead    credOpenMode = iota // nothing: a fresh file is refused
	credOpenWrite                       // initialize a fresh file
	credOpenUpgrade                     // initialize a fresh file, upgrade version 1
)

// openCredentialsSQLite opens the driver handle and enforces the credentials
// schema: a fresh file is initialized (unless mode is credOpenRead), version 1
// or credentialsSchemaVersion must hold exactly the session_credentials table,
// version 1 is upgraded only with credOpenUpgrade, and anything else is
// refused unmodified.
func openCredentialsSQLite(display, dsn string, mode credOpenMode) (*Credentials, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: open credentials %q: %w", display, err)
	}
	db.SetMaxOpenConns(1)
	c := &Credentials{db: db, display: display}
	if err := c.enforceSchema(mode); err != nil {
		_ = db.Close()
		return nil, err
	}
	return c, nil
}

func (c *Credentials) enforceSchema(mode credOpenMode) error {
	display, db := c.display, c.db
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("store: open credentials %q: read schema version: %w", display, err)
	}
	tables, err := applicationTables(db)
	if err != nil {
		return fmt.Errorf("store: open credentials %q: inspect schema: %w", display, err)
	}
	known := version >= credentialsSchemaV1 && version <= credentialsSchemaVersion
	switch {
	case version == 0 && len(tables) == 0:
		if mode == credOpenRead {
			return &CredentialsSchemaError{Path: display, Detail: "it has no credentials schema yet (serve or session new creates it)"}
		}
		c.version = credentialsSchemaVersion
		return initCredentialsSchema(display, db)
	case known && ((version < 3 && len(tables) == 1 && tables[0] == credentialsTable) || (version == 3 && len(tables) == 2 && tables[0] == "episode_starts" && tables[1] == credentialsTable)):
		c.version = version
		if version < credentialsSchemaVersion && mode == credOpenUpgrade {
			if err := upgradeCredentialsV3(display, db, version); err != nil {
				return err
			}
			c.version, c.upgradedFrom = credentialsSchemaVersion, version
		}
		return nil
	case known && len(tables) == 0:
		// Its only table was dropped: recreating it empty would silently
		// revoke every credential, so it is refused like commit_objects is.
		return &CredentialsSchemaError{Path: display, Detail: "it is missing table " + credentialsTable}
	case version > credentialsSchemaVersion:
		return &CredentialsSchemaError{Path: display, Detail: fmt.Sprintf("user_version %d is newer than this binary's %d; use a compatible binary", version, credentialsSchemaVersion)}
	default:
		return &CredentialsSchemaError{Path: display, Detail: fmt.Sprintf("user_version %d with tables %v is not a credentials store (is it the world store?)", version, tables)}
	}
}

// upgradeCredentialsV3 adds start history and, when upgrading version 1,
// the binding columns, in one transaction. Existing credentials are unchanged.
func upgradeCredentialsV3(display string, db *sql.DB, version int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("store: upgrade credentials %q: begin: %w", display, err)
	}
	defer func() { _ = tx.Rollback() }()
	statements := []string{episodeStartsSQL}
	if version == 1 {
		statements = append(append([]string{}, credentialsUpgradeV2...), statements...)
	}
	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
		}
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", credentialsSchemaVersion)); err != nil {
		return fmt.Errorf("store: upgrade credentials %q: set version: %w", display, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: upgrade credentials %q: commit: %w", display, err)
	}
	return nil
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
	return c.MintSessionWithStart(ctx, row, nil)
}

// MintSessionWithStart atomically records the attended start and its credential.
func (c *Credentials) MintSessionWithStart(ctx context.Context, row SessionRow, start *EpisodeStart) error {
	return c.mintSession(ctx, row, start, nil)
}

func (c *Credentials) mintSession(ctx context.Context, row SessionRow, start *EpisodeStart, resume *EpisodeResume) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if row.CredentialID == "" || row.EpisodeID == "" {
		return fmt.Errorf("store: mint session: empty credential_id or episode_id")
	}
	if (row.Project == "") != (row.ProfileDigest == "") {
		return fmt.Errorf("store: mint session: a profile binding needs both the project and its digest")
	}
	if (row.Bound() || start != nil || resume != nil) && c.version < credentialsSchemaVersion {
		return &UnboundCredentialsStoreError{Path: c.display}
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: mint session: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if resume != nil {
		if err := checkResume(ctx, tx, row, *resume); err != nil {
			return err
		}
		resume.CredentialID = row.CredentialID
	}
	// An unbound row names only the five original columns, so the same
	// statement works on a version-1 store; the binding columns stay NULL.
	stmt := `INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?);`
	args := []any{row.CredentialID, row.EpisodeID, row.GrantsJSON, row.ExpiresAt, row.CreatedAt}
	if row.Bound() {
		stmt = `INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at, project, profile_digest)
		 VALUES (?, ?, ?, ?, ?, ?, ?);`
		args = append(args, row.Project, row.ProfileDigest)
	}
	if _, err := tx.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("store: mint session: %w", err)
	}
	if start != nil {
		if start.Episode != row.EpisodeID || start.Schema != "world/episode-start/v1" {
			return fmt.Errorf("store: mint session: invalid start record")
		}
		data, err := json.Marshal(start)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO episode_starts (episode, start_json) VALUES (?, ?)`, row.EpisodeID, string(data)); err != nil {
			return fmt.Errorf("store: mint start: %w", err)
		}
	}
	if resume != nil {
		if err := appendResume(ctx, tx, row.EpisodeID, *resume); err != nil {
			return fmt.Errorf("store: mint resume: %w", err)
		}
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
	var project, digest sql.NullString
	err := c.db.QueryRowContext(ctx,
		sessionColumns(c.version)+` FROM session_credentials WHERE credential_id = ?;`,
		credentialID,
	).Scan(&row.CredentialID, &row.EpisodeID, &row.GrantsJSON, &row.ExpiresAt, &row.CreatedAt, &project, &digest)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return SessionRow{}, false, nil
	case err != nil:
		return SessionRow{}, false, fmt.Errorf("store: resolve session: %w", err)
	}
	row.Project, row.ProfileDigest = project.String, digest.String
	return row, true, nil
}

// sessionColumns is the SELECT list of one SessionRow at a schema version: a
// version-1 table (and a schema-4 world store's) has no binding columns, so
// they read as NULL there — unbound.
func sessionColumns(version int) string {
	if version >= 2 {
		return `SELECT credential_id, episode_id, grants_json, expires_at, created_at, project, profile_digest`
	}
	return `SELECT credential_id, episode_id, grants_json, expires_at, created_at, NULL, NULL`
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

// ListSessions returns up to limit session_credentials rows (at most
// MaxListSessions), ordered by created_at then credential_id; a non-empty
// episodeID restricts it to that episode. Rows carry the credential_id HASH
// only — the raw token is never stored (D3). It is a plain read, so it works
// on an OpenCredentialsReadOnly handle beside a live writer.
func (c *Credentials) ListSessions(ctx context.Context, episodeID string, limit int) ([]SessionRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	return listSessions(ctx, c.db, c.version, episodeID, limit)
}

// listSessions is ListSessions at a schema version; LegacySessions reuses it
// on a schema-4 world store, whose table has the five version-1 columns.
func listSessions(ctx context.Context, db *sql.DB, version int, episodeID string, limit int) ([]SessionRow, error) {
	if limit <= 0 || limit > MaxListSessions {
		limit = MaxListSessions
	}
	q := sessionColumns(version) + ` FROM session_credentials`
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
		var project, digest sql.NullString
		if err := rows.Scan(&r.CredentialID, &r.EpisodeID, &r.GrantsJSON, &r.ExpiresAt, &r.CreatedAt, &project, &digest); err != nil {
			return nil, fmt.Errorf("store: list sessions: scan: %w", err)
		}
		r.Project, r.ProfileDigest = project.String, digest.String
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list sessions: %w", err)
	}
	return out, nil
}
