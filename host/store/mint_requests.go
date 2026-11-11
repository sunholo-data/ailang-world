package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
)

// Row 201 A0 (design_docs/planned/w-world-approvals-by-key-a0a1.md §2.4): stored, single-use
// mint requests. `episode start --request` puts one; `approve` consumes it in the same
// transaction that inserts its session row.

// ---- M4a ----
const mintRequestsTable = "mint_requests"

// mintRequestsSQL is the version 3 -> 4 upgrade: credentials_schema.sql's mint_requests table,
// column for column (TestCredentialsUpgradeToV4 compares the two).
const mintRequestsSQL = `CREATE TABLE mint_requests (
    request_id    TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
    request       BLOB NOT NULL,
    digest        TEXT NOT NULL CHECK (length(digest) = 64),
    challenge     TEXT NOT NULL CHECK (length(challenge) = 64),
    episode_id    TEXT NOT NULL CHECK (episode_id <> ''),
    created_at    INTEGER NOT NULL,
    expires_at    INTEGER NOT NULL CHECK (expires_at > created_at),
    consumed_at   INTEGER,
    approved_by   TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
    credential_id TEXT,
    CHECK ((consumed_at IS NULL) = (approved_by IS NULL) AND (consumed_at IS NULL) = (credential_id IS NULL))
)`

// ---- M4b ----
// mintRequestsSchemaVersion is the first credentials schema with mint_requests.
const mintRequestsSchemaVersion = 4

// MaxPendingMintRequests caps the unconsumed, unexpired requests one store holds (design §4, R4).
const MaxPendingMintRequests = 16

// mintRequestPruneAfter: an unconsumed request is deleted once it has been expired this long.
const mintRequestPruneAfter = 86400

// MintRequestIDPattern is a request id: "mr-" and 16 lowercase hex (8 random bytes). The table's
// CHECK is weaker (plan §1.4 F2), so PutMintRequest enforces this.
var MintRequestIDPattern = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)

var (
	ErrTooManyPendingRequests = errors.New("store: 16 mint requests are already pending; approve them or let them expire (ailang-worldd approve --list)")
	ErrRequestUnknown         = errors.New("store: no such mint request")
	ErrRequestConsumed        = errors.New("store: the mint request was already used")
	ErrRequestExpired         = errors.New("store: the mint request expired before it was consumed")
	ErrRequestClockBehind     = errors.New("store: the clock is behind the mint request's creation time")
	ErrRequestDigestChanged   = errors.New("store: the mint request's digest is not the one approved")
	ErrBadMintRequest         = errors.New("store: invalid mint request row")
	ErrIncompleteConsume      = errors.New("store: a mint request consume needs an id, a digest, a challenge, approved_by, a clock and a check")
	// ErrRequestChallengeChanged: the row's challenge is not the one the approval was made for
	// (row 201 A1b §2.5 C6, C7).
	ErrRequestChallengeChanged = errors.New("store: the mint request's challenge is not the one signed")
	// ErrInconsistentApproval: approved_by is "tty" exactly when no approver is bound, and
	// "<format>:<name>" of the bound approver otherwise (row 201 A1b §2.5 C3, C4).
	ErrInconsistentApproval = errors.New("store: approved_by does not match the approver the consume binds")
)

// OlderCredentialsStoreError refuses an operation on a store below the version its Feature needs:
// mint requests need 4, approvers need 5 (row 201 A1b §2.1).
type OlderCredentialsStoreError struct {
	Path    string
	Version int    // the store's version
	Need    int    // the first version with Feature
	Feature string // "mint requests" or "approvers"
}

func (e *OlderCredentialsStoreError) Error() string {
	return fmt.Sprintf("store: credentials store %q is schema %d; %s need schema %d: "+
		"restart `ailang-worldd serve` with this build once (it upgrades the store at startup), then run this again",
		e.Path, e.Version, e.Feature, e.Need)
}

// MintRequestRow is one mint_requests row.
type MintRequestRow struct {
	RequestID    string
	Request      []byte // the canonical request bytes (approval.Canonical)
	Digest       string // 64 lowercase hex: SHA-256 of Request when it was made
	Challenge    string // 64 lowercase hex
	EpisodeID    string
	CreatedAt    int64
	ExpiresAt    int64
	Consumed     bool   // consumed_at IS NOT NULL
	ConsumedAt   int64  // 0 unless Consumed
	ApprovedBy   string // "" unless Consumed
	CredentialID string // "" unless Consumed
}

// PutMintRequest stores one new, unconsumed request. It prunes requests expired a day ago, and
// refuses when MaxPendingMintRequests live ones are already held.
func (c *Credentials) PutMintRequest(ctx context.Context, row MintRequestRow, now int64) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < mintRequestsSchemaVersion {
		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
	}
	if !MintRequestIDPattern.MatchString(row.RequestID) {
		return fmt.Errorf("%w: request id %q", ErrBadMintRequest, row.RequestID)
	}
	if len(row.Request) == 0 {
		return fmt.Errorf("%w: empty request", ErrBadMintRequest)
	}
	if !isLowerHex64(row.Digest) {
		return fmt.Errorf("%w: digest is not 64 lowercase hex", ErrBadMintRequest)
	}
	if !isLowerHex64(row.Challenge) {
		return fmt.Errorf("%w: challenge is not 64 lowercase hex", ErrBadMintRequest)
	}
	if row.EpisodeID == "" {
		return fmt.Errorf("%w: empty episode id", ErrBadMintRequest)
	}
	if row.ExpiresAt <= row.CreatedAt {
		return fmt.Errorf("%w: expires_at %d is not after created_at %d", ErrBadMintRequest, row.ExpiresAt, row.CreatedAt)
	}
	if row.Consumed {
		return fmt.Errorf("%w: a new request is unconsumed", ErrBadMintRequest)
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: put mint request: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM mint_requests WHERE consumed_at IS NULL AND expires_at <= ?`, now-mintRequestPruneAfter); err != nil {
		return fmt.Errorf("store: put mint request: prune: %w", err)
	}
	var live int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM mint_requests WHERE consumed_at IS NULL AND expires_at > ?`, now).Scan(&live); err != nil {
		return fmt.Errorf("store: put mint request: count: %w", err)
	}
	if live >= MaxPendingMintRequests {
		return ErrTooManyPendingRequests
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		row.RequestID, row.Request, row.Digest, row.Challenge, row.EpisodeID, row.CreatedAt, row.ExpiresAt); err != nil {
		return fmt.Errorf("store: put mint request: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: put mint request: commit: %w", err)
	}
	return nil
}

// MintRequest reads one request by id. ok is false when no such row exists.
func (c *Credentials) MintRequest(ctx context.Context, id string) (MintRequestRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return MintRequestRow{}, false, err
	}
	if c.version < mintRequestsSchemaVersion {
		return MintRequestRow{}, false, &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
	}
	row, err := scanMintRequest(c.db.QueryRowContext(ctx,
		`SELECT `+mintRequestColumns+` FROM mint_requests WHERE request_id = ?`, id).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return MintRequestRow{}, false, nil
	}
	if err != nil {
		return MintRequestRow{}, false, fmt.Errorf("store: read mint request: %w", err)
	}
	return row, true, nil
}

// PendingMintRequests lists the live, unconsumed requests, newest first. limit <= 0 or above
// MaxPendingMintRequests means MaxPendingMintRequests.
func (c *Credentials) PendingMintRequests(ctx context.Context, now int64, limit int) ([]MintRequestRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < mintRequestsSchemaVersion {
		return nil, &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
	}
	if limit <= 0 || limit > MaxPendingMintRequests {
		limit = MaxPendingMintRequests
	}
	rows, err := c.db.QueryContext(ctx,
		`SELECT `+mintRequestColumns+` FROM mint_requests WHERE consumed_at IS NULL AND expires_at > ? ORDER BY created_at DESC, request_id DESC LIMIT ?`,
		now, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list mint requests: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []MintRequestRow
	for rows.Next() {
		row, err := scanMintRequest(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("store: list mint requests: %w", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list mint requests: %w", err)
	}
	return out, nil
}

// MaxDecidedMintRequests caps one DecidedMintRequests listing (row 201 A1b design §2.4 [r2]).
const MaxDecidedMintRequests = 64

// DecidedMintRequests lists the consumed requests, newest consumed_at first and request_id
// descending among equal times (design §2.4 [r2]), at most MaxDecidedMintRequests (limit <= 0
// or above it means the cap). A non-empty approvedBy keeps only the rows it approved
// (approver revoke's K8). It reads every store with mint requests, schema 4 included.
func (c *Credentials) DecidedMintRequests(ctx context.Context, approvedBy string, limit int) ([]MintRequestRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < mintRequestsSchemaVersion {
		return nil, &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
	}
	if limit <= 0 || limit > MaxDecidedMintRequests {
		limit = MaxDecidedMintRequests
	}
	q, args := `SELECT `+mintRequestColumns+` FROM mint_requests WHERE consumed_at IS NOT NULL`, []any{}
	if approvedBy != "" {
		q, args = q+` AND approved_by = ?`, append(args, approvedBy)
	}
	rows, err := c.db.QueryContext(ctx, q+` ORDER BY consumed_at DESC, request_id DESC LIMIT ?`, append(args, limit)...)
	if err != nil {
		return nil, fmt.Errorf("store: list decided mint requests: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []MintRequestRow
	for rows.Next() {
		row, err := scanMintRequest(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("store: list decided mint requests: %w", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list decided mint requests: %w", err)
	}
	return out, nil
}

// isLowerHex64 reports whether s is exactly 64 lowercase hexadecimal characters.
func isLowerHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if !(ch >= '0' && ch <= '9' || ch >= 'a' && ch <= 'f') {
			return false
		}
	}
	return true
}

const mintRequestColumns = `request_id, request, digest, challenge, episode_id, created_at, expires_at, consumed_at, approved_by, credential_id`

// scanMintRequest reads one mintRequestColumns row through scan (a Row.Scan or Rows.Scan).
func scanMintRequest(scan func(...any) error) (MintRequestRow, error) {
	var (
		row         MintRequestRow
		consumedAt  sql.NullInt64
		approvedBy  sql.NullString
		credentials sql.NullString
	)
	if err := scan(&row.RequestID, &row.Request, &row.Digest, &row.Challenge, &row.EpisodeID,
		&row.CreatedAt, &row.ExpiresAt, &consumedAt, &approvedBy, &credentials); err != nil {
		return MintRequestRow{}, err
	}
	row.Consumed = consumedAt.Valid
	row.ConsumedAt = consumedAt.Int64
	row.ApprovedBy = approvedBy.String
	row.CredentialID = credentials.String
	return row, nil
}

// ---- M4c ----

// RequestConsume is what MintSessionForRequest consumes. Clock is read exactly once, inside the
// transaction, after BeginTx (design §2.4 [r3]); that reading is consumed_at and the expiry
// bound. Check runs with that reading before the guarded UPDATE; an error rolls everything back
// (plan §1.3 D2: authority.MintForRequest puts approval.MintMatches here). Challenge is the row's
// challenge (64 lowercase hex), bound in the guarded UPDATE (row 201 A1b §2.5 C1, C6).
type RequestConsume struct {
	ID, Digest, ApprovedBy string
	Challenge              string
	Approver               *ApproverUse // nil: the typed y; else the approver whose signature approved (§2.5 C8)
	Clock                  func() int64
	Check                  func(consumedAt int64) error
}

// ApproverUse is the approver a signature consume binds: the consume marks that approver's active
// row used, with this exact format and key, in its own transaction, or refuses.
type ApproverUse struct {
	Name, Format string
	PublicKey    []byte
	Passkey      *PasskeyUse // non-nil exactly when Format is "webauthn" (row 201 A2 C13)
}

// MintSessionForRequest consumes request rc and inserts session in one transaction: the request
// is marked consumed (consumed_at, approved_by, credential_id) only if it is still pending, and
// the session row lands with it or not at all.
func (c *Credentials) MintSessionForRequest(ctx context.Context, session SessionRow, rc RequestConsume) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if rc.ID == "" || rc.Digest == "" || !isLowerHex64(rc.Challenge) || rc.ApprovedBy == "" || rc.Clock == nil || rc.Check == nil { // C1
		return ErrIncompleteConsume
	}
	if a := rc.Approver; a != nil && (a.Name == "" || a.Format == "" || len(a.PublicKey) == 0) { // C2
		return ErrIncompleteConsume
	}
	if rc.Approver == nil && rc.ApprovedBy != "tty" { // C3
		return ErrInconsistentApproval
	}
	if rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format+":"+rc.Approver.Name { // C4
		return ErrInconsistentApproval
	}
	if session.CredentialID == "" || session.EpisodeID == "" {
		return fmt.Errorf("store: mint session: empty credential_id or episode_id")
	}
	if c.version < mintRequestsSchemaVersion {
		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
	}
	if rc.Approver != nil && c.version < approversSchemaVersion { // C5
		return c.olderForApprovers()
	}
	if a := rc.Approver; a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13
		return ErrIncompleteConsume
	}
	if a := rc.Approver; a != nil && a.Passkey != nil {
		if p := a.Passkey; len(p.CredentialID) < 1 || len(p.CredentialID) > 1023 || p.SignCount < 0 || p.SignCount > maxSignCount { // C14
			return ErrIncompleteConsume
		}
		if c.version < passkeysSchemaVersion { // C15
			return c.olderForPasskeys()
		}
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: consume mint request: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	at := rc.Clock()
	if err := rc.Check(at); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `UPDATE mint_requests SET consumed_at = ?, approved_by = ?, credential_id = ? WHERE request_id = ? AND digest = ? AND challenge = ? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?`,
		at, rc.ApprovedBy, session.CredentialID, rc.ID, rc.Digest, rc.Challenge, at, at) // C6
	if err != nil {
		return fmt.Errorf("store: consume mint request: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: consume mint request: rows: %w", err)
	}
	if n != 1 {
		return classifyConsume(ctx, tx, rc, at, n)
	}
	if rc.Approver != nil { // C8: the approver's row, in this transaction, or nothing is consumed
		if err := useApprover(ctx, tx, *rc.Approver, at); err != nil {
			return err
		}
		if p := rc.Approver.Passkey; p != nil { // C16-C18: and its counter, forward only
			if err := usePasskey(ctx, tx, rc.Approver.Name, *p); err != nil {
				return err
			}
		}
	}
	if err := insertSessionRow(ctx, tx, session); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: consume mint request: commit: %w", err)
	}
	return nil
}

// classifyConsume names why the guarded UPDATE changed no row, reading the request inside tx.
func classifyConsume(ctx context.Context, tx *sql.Tx, rc RequestConsume, at int64, changed int64) error {
	var (
		digest, challenge string
		createdAt, expAt  int64
		consumed          bool
	)
	err := tx.QueryRowContext(ctx, `SELECT digest, challenge, created_at, expires_at, consumed_at IS NOT NULL FROM mint_requests WHERE request_id = ?`, rc.ID).
		Scan(&digest, &challenge, &createdAt, &expAt, &consumed)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrRequestUnknown, rc.ID)
	}
	if err != nil {
		return fmt.Errorf("store: consume mint request: classify: %w", err)
	}
	switch {
	case consumed:
		return fmt.Errorf("%w: %s", ErrRequestConsumed, rc.ID)
	case at < createdAt:
		return fmt.Errorf("%w: %s", ErrRequestClockBehind, rc.ID)
	case at >= expAt:
		return fmt.Errorf("%w: %s", ErrRequestExpired, rc.ID)
	case digest != rc.Digest:
		return fmt.Errorf("%w: %s", ErrRequestDigestChanged, rc.ID)
	case challenge != rc.Challenge: // C7
		return fmt.Errorf("%w: %s", ErrRequestChallengeChanged, rc.ID)
	}
	return fmt.Errorf("store: mint request %s: the guarded update changed %d rows", rc.ID, changed)
}
