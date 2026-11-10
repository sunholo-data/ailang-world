package store

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.1, §2.4): the approvers
// table (credentials schema 5). A row is written only by `approver enroll` after a typed y; it is
// marked revoked, never deleted, so a name and a key are each enrolled at most once, ever. The
// store does not import host/approval: the name grammar and key parsing are the pure law's, run
// by the CLI first; the checks here are the store's own backstop.

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ---- M2a ----
const approversTable = "approvers"

// approversSchemaVersion is the first credentials schema with approvers.
const approversSchemaVersion = 5

// approversSQL is the version 4 -> 5 upgrade: credentials_schema.sql's approvers table, column
// for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
const approversSQL = `CREATE TABLE approvers (
    name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    format       TEXT NOT NULL CHECK (format <> ''),
    public_key   BLOB NOT NULL UNIQUE CHECK (length(public_key) > 0),
    enrolled_at  INTEGER NOT NULL CHECK (enrolled_at > 0),
    last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
    revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
)`

// ApproverRow is one approvers row. LastUsedAt and RevokedAt are 0 for NULL.
type ApproverRow struct {
	Name       string
	Format     string
	PublicKey  []byte // the SSH wire blob
	EnrolledAt int64
	LastUsedAt int64 // 0 = NULL (never used)
	RevokedAt  int64 // 0 = NULL (active)
}

var (
	ErrApproverNameTaken   = errors.New("store: an approver with that name exists (names are never reused, even after revoke)")
	ErrApproverKeyEnrolled = errors.New("store: that key is already enrolled (a key is never enrolled twice, even after revoke)")
	ErrApproverUnknown     = errors.New("store: no such approver")
	ErrApproverRevoked     = errors.New("store: the approver is revoked")
	ErrApproverKeyChanged  = errors.New("store: the approver's key is not the key that signed")
	ErrApproverClockBehind = errors.New("store: the clock is behind the approver's enrolment time")
	ErrBadApprover         = errors.New("store: invalid approver row")
)

// maxListedApprovers caps ListApprovers.
const maxListedApprovers = 256

// ---- M2b ----

// olderForApprovers is the refusal of every approver method below version 5.
func (c *Credentials) olderForApprovers() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: approversSchemaVersion, Feature: "approvers"}
}

const approverColumns = `name, format, public_key, enrolled_at, last_used_at, revoked_at`

func scanApprover(scan func(...any) error) (ApproverRow, error) {
	var r ApproverRow
	var lastUsed, revoked sql.NullInt64
	if err := scan(&r.Name, &r.Format, &r.PublicKey, &r.EnrolledAt, &lastUsed, &revoked); err != nil {
		return ApproverRow{}, err
	}
	r.LastUsedAt, r.RevokedAt = lastUsed.Int64, revoked.Int64
	return r, nil
}

// EnrolApprover stores one new, active approver. It refuses a name or a key that any row holds,
// revoked rows included; the UNIQUE constraints are the backstop and their error text is never
// parsed.
func (c *Credentials) EnrolApprover(ctx context.Context, row ApproverRow) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < approversSchemaVersion {
		return c.olderForApprovers()
	}
	if row.Name == "" {
		return fmt.Errorf("%w: empty name", ErrBadApprover)
	}
	if row.Format != "ssh" {
		return fmt.Errorf("%w: format %q is not ssh", ErrBadApprover, row.Format)
	}
	if len(row.PublicKey) == 0 {
		return fmt.Errorf("%w: empty public key", ErrBadApprover)
	}
	if row.EnrolledAt <= 0 {
		return fmt.Errorf("%w: enrolled_at %d is not positive", ErrBadApprover, row.EnrolledAt)
	}
	if row.LastUsedAt != 0 {
		return fmt.Errorf("%w: a new approver has never been used", ErrBadApprover)
	}
	if row.RevokedAt != 0 {
		return fmt.Errorf("%w: a new approver is not revoked", ErrBadApprover)
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: enrol approver: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM approvers WHERE name = ?`, row.Name).Scan(&n); err != nil {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if n != 0 {
		return fmt.Errorf("%w: %q", ErrApproverNameTaken, row.Name)
	}
	var holder string
	err = tx.QueryRowContext(ctx, `SELECT name FROM approvers WHERE public_key = ?`, row.PublicKey).Scan(&holder)
	if err == nil {
		return fmt.Errorf("%w: it is approver %q's key", ErrApproverKeyEnrolled, holder)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`,
		row.Name, row.Format, row.PublicKey, row.EnrolledAt); err != nil {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: enrol approver: commit: %w", err)
	}
	return nil
}

// Approver reads one approver by name. ok is false when no such row exists.
func (c *Credentials) Approver(ctx context.Context, name string) (ApproverRow, bool, error) {
	return c.approverWhere(ctx, `name = ?`, name)
}

// ApproverByKey reads the approver whose public_key is blob, byte for byte. ok is false when none.
func (c *Credentials) ApproverByKey(ctx context.Context, blob []byte) (ApproverRow, bool, error) {
	return c.approverWhere(ctx, `public_key = ?`, blob)
}

func (c *Credentials) approverWhere(ctx context.Context, where string, arg any) (ApproverRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return ApproverRow{}, false, err
	}
	if c.version < approversSchemaVersion {
		return ApproverRow{}, false, c.olderForApprovers()
	}
	row, err := scanApprover(c.db.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE `+where, arg).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return ApproverRow{}, false, nil
	}
	if err != nil {
		return ApproverRow{}, false, fmt.Errorf("store: read approver: %w", err)
	}
	return row, true, nil
}

// ---- M2c ----

// ListApprovers lists every approver, revoked ones included, ordered by name, at most 256.
func (c *Credentials) ListApprovers(ctx context.Context) ([]ApproverRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < approversSchemaVersion {
		return nil, c.olderForApprovers()
	}
	rows, err := c.db.QueryContext(ctx, `SELECT `+approverColumns+` FROM approvers ORDER BY name LIMIT ?`, maxListedApprovers)
	if err != nil {
		return nil, fmt.Errorf("store: list approvers: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []ApproverRow
	for rows.Next() {
		r, err := scanApprover(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("store: list approvers: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list approvers: %w", err)
	}
	return out, nil
}

// RevokeApprover marks the active approver name revoked at now. The row is kept. It refuses an
// unknown name, a row already revoked, and a clock behind the enrolment, each with its own error.
func (c *Credentials) RevokeApprover(ctx context.Context, name string, now int64) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < approversSchemaVersion {
		return c.olderForApprovers()
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: revoke approver: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, `UPDATE approvers SET revoked_at = ? WHERE name = ? AND revoked_at IS NULL AND enrolled_at <= ?`, now, name, now)
	if err != nil {
		return fmt.Errorf("store: revoke approver: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: revoke approver: %w", err)
	}
	if n != 1 {
		return classifyRevoke(ctx, tx, name, now)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: revoke approver: commit: %w", err)
	}
	return nil
}

// classifyRevoke names why a revoke changed no row, reading the row in the same transaction.
// The order matches approval.CheckRevoke: revoked before clock-behind.
func classifyRevoke(ctx context.Context, tx *sql.Tx, name string, now int64) error {
	row, err := scanApprover(tx.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE name = ?`, name).Scan)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("%w: %q", ErrApproverUnknown, name)
	case err != nil:
		return fmt.Errorf("store: revoke approver: %w", err)
	case row.RevokedAt != 0:
		return fmt.Errorf("%w: %q at %d", ErrApproverRevoked, name, row.RevokedAt)
	case row.EnrolledAt > now:
		return fmt.Errorf("%w: %q enrolled at %d, now %d", ErrApproverClockBehind, name, row.EnrolledAt, now)
	default:
		return fmt.Errorf("store: revoke approver %q: changed no row", name)
	}
}

// ---- M4d (sprint 3) ----

// useApprover is the consume's approver binding (row 201 A1b §2.5 C8): inside the consume's
// transaction, at its clock reading, it marks the approver used only while its row is active,
// holds this exact format and key, and was enrolled by at. Otherwise classifyApproverUse names
// why (C9-C12), and the caller's rollback undoes the request update.
func useApprover(ctx context.Context, tx *sql.Tx, a ApproverUse, at int64) error {
	res, err := tx.ExecContext(ctx, `UPDATE approvers SET last_used_at = ? WHERE name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?`,
		at, a.Name, a.Format, a.PublicKey, at)
	if err != nil {
		return fmt.Errorf("store: use approver: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: use approver: %w", err)
	}
	if n != 1 {
		return classifyApproverUse(ctx, tx, a, at)
	}
	return nil
}

// classifyApproverUse names why useApprover changed no row, reading the row in the same transaction.
func classifyApproverUse(ctx context.Context, tx *sql.Tx, a ApproverUse, at int64) error {
	row, err := scanApprover(tx.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE name = ?`, a.Name).Scan)
	switch {
	case errors.Is(err, sql.ErrNoRows): // C9
		return fmt.Errorf("%w: %q", ErrApproverUnknown, a.Name)
	case err != nil:
		return fmt.Errorf("store: use approver: %w", err)
	case row.RevokedAt != 0: // C10
		return fmt.Errorf("%w: %q at %d", ErrApproverRevoked, a.Name, row.RevokedAt)
	case row.Format != a.Format || !bytes.Equal(row.PublicKey, a.PublicKey): // C11
		return fmt.Errorf("%w: %q", ErrApproverKeyChanged, a.Name)
	case row.EnrolledAt > at: // C12
		return fmt.Errorf("%w: %q enrolled at %d, now %d", ErrApproverClockBehind, a.Name, row.EnrolledAt, at)
	default:
		return fmt.Errorf("store: use approver %q: changed no row", a.Name)
	}
}
