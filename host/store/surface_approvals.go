package store

// Row 201 A3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.2): the surface's recorded
// approvals (credentials schema 7). The daemon stores a POSTed signed artifact only after the A1b/A2
// law verified it; the minter (approve --await) reads it and judges it again. A row is never a
// credential and never trusted: its approver and digest columns are display only.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ---- M1a ----

const surfaceApprovalsTable = "surface_approvals"

// surfaceApprovalsSchemaVersion is the first credentials schema with surface_approvals.
const surfaceApprovalsSchemaVersion = 7

// surfaceApprovalsSQL is the version 6 -> 7 upgrade: credentials_schema.sql's surface_approvals
// table, column for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
const surfaceApprovalsSQL = `CREATE TABLE surface_approvals (
    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
    approver    TEXT NOT NULL CHECK (approver <> ''),
    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
    digest      TEXT NOT NULL CHECK (length(digest) = 64),
    received_at INTEGER NOT NULL
)`

// ---- M1b ----

// MaxSurfaceApprovalArtifactBytes bounds a recorded artifact: the A2 assertion envelope's and the
// --signature file's own 16 KiB bound (approval.MaxPasskeyResponseBytes, approveMaxSignatureBytes).
const MaxSurfaceApprovalArtifactBytes = 16 << 10

// ErrBadSurfaceApproval refuses a row that breaks the table's grammar (S2).
var ErrBadSurfaceApproval = errors.New("store: invalid surface approval row")

// SurfaceApprovalRow is one surface_approvals row: the artifact the daemon verified for a request.
// Approver ("<format>:<name>") and Digest (the request's digest when it was recorded) are display
// only; the minter judges Artifact again against the request it reads.
type SurfaceApprovalRow struct {
	RequestID  string
	Format     string // "webauthn" or "ssh"
	Approver   string // "<format>:<name>"
	Artifact   []byte // the assertion JSON or the armored SSHSIG, 1..16384 bytes
	Digest     string // 64 lowercase hex
	ReceivedAt int64  // set by RecordSurfaceApproval to its now
}

func (c *Credentials) olderForSurfaceApprovals() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: surfaceApprovalsSchemaVersion, Feature: "surface approvals"}
}

// checkSurfaceApprovalRow is RecordSurfaceApproval's grammar (S2): each refusal is ErrBadSurfaceApproval.
func checkSurfaceApprovalRow(row SurfaceApprovalRow) error {
	switch {
	case !MintRequestIDPattern.MatchString(row.RequestID):
		return fmt.Errorf("%w: request id %q", ErrBadSurfaceApproval, row.RequestID)
	case row.Format != "webauthn" && row.Format != "ssh":
		return fmt.Errorf("%w: format %q is not webauthn or ssh", ErrBadSurfaceApproval, row.Format)
	case !strings.HasPrefix(row.Approver, row.Format+":") || len(row.Approver) == len(row.Format)+1:
		return fmt.Errorf("%w: approver %q is not %s:<name>", ErrBadSurfaceApproval, row.Approver, row.Format)
	case len(row.Artifact) < 1 || len(row.Artifact) > MaxSurfaceApprovalArtifactBytes:
		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadSurfaceApproval, MaxSurfaceApprovalArtifactBytes, len(row.Artifact))
	case !isLowerHex64(row.Digest):
		return fmt.Errorf("%w: digest is not 64 lowercase hex", ErrBadSurfaceApproval)
	}
	return nil
}

// RecordSurfaceApproval stores row for a pending request in one transaction, received at now: the
// request must exist (S3), be unconsumed, created at or before now, unexpired at now, and carry
// row's digest (S4), or nothing is written and the request's own sentinel is returned. A later row
// for the same request replaces the earlier one (S5): the minter may have refused it.
func (c *Credentials) RecordSurfaceApproval(ctx context.Context, row SurfaceApprovalRow, now int64) error {
	if err := requireDeadline(ctx); err != nil { // S1
		return err
	}
	if c.version < surfaceApprovalsSchemaVersion { // S1
		return c.olderForSurfaceApprovals()
	}
	if err := checkSurfaceApprovalRow(row); err != nil { // S2
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: record surface approval: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var (
		digest           string
		createdAt, expAt int64
		consumed         bool
	)
	err = tx.QueryRowContext(ctx, `SELECT digest, created_at, expires_at, consumed_at IS NOT NULL FROM mint_requests WHERE request_id = ?`, row.RequestID).
		Scan(&digest, &createdAt, &expAt, &consumed)
	if errors.Is(err, sql.ErrNoRows) { // S3
		return fmt.Errorf("%w: %s", ErrRequestUnknown, row.RequestID)
	}
	if err != nil {
		return fmt.Errorf("store: record surface approval: %w", err)
	}
	switch { // S4
	case consumed:
		return fmt.Errorf("%w: %s", ErrRequestConsumed, row.RequestID)
	case now < createdAt:
		return fmt.Errorf("%w: %s", ErrRequestClockBehind, row.RequestID)
	case now >= expAt:
		return fmt.Errorf("%w: %s", ErrRequestExpired, row.RequestID)
	case digest != row.Digest:
		return fmt.Errorf("%w: %s", ErrRequestDigestChanged, row.RequestID)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO surface_approvals (request_id, format, approver, artifact, digest, received_at) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(request_id) DO UPDATE SET format = excluded.format, approver = excluded.approver, artifact = excluded.artifact, digest = excluded.digest, received_at = excluded.received_at`,
		row.RequestID, row.Format, row.Approver, row.Artifact, row.Digest, now); err != nil { // S5
		return fmt.Errorf("store: record surface approval: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: record surface approval: commit: %w", err)
	}
	return nil
}

// SurfaceApproval reads the recorded approval for request id. ok is false when there is none.
func (c *Credentials) SurfaceApproval(ctx context.Context, id string) (SurfaceApprovalRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return SurfaceApprovalRow{}, false, err
	}
	if c.version < surfaceApprovalsSchemaVersion {
		return SurfaceApprovalRow{}, false, c.olderForSurfaceApprovals()
	}
	var row SurfaceApprovalRow
	err := c.db.QueryRowContext(ctx, `SELECT request_id, format, approver, artifact, digest, received_at FROM surface_approvals WHERE request_id = ?`, id).
		Scan(&row.RequestID, &row.Format, &row.Approver, &row.Artifact, &row.Digest, &row.ReceivedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return SurfaceApprovalRow{}, false, nil
	}
	if err != nil {
		return SurfaceApprovalRow{}, false, fmt.Errorf("store: read surface approval: %w", err)
	}
	return row, true, nil
}
