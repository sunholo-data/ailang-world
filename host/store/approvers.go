package store

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.1, §2.4): the approvers
// table (credentials schema 5). A row is written only by `approver enroll` after a typed y; it is
// marked revoked, never deleted, so a name and a key are each enrolled at most once, ever. The
// store does not import host/approval: the name grammar and key parsing are the pure law's, run
// by the CLI first; the checks here are the store's own backstop.

import (
	"errors"
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
