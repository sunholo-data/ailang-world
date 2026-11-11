package store

// Row 201 A3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.2): the surface's recorded
// approvals (credentials schema 7). The daemon stores a POSTed signed artifact only after the A1b/A2
// law verified it; the minter (approve --await) reads it and judges it again. A row is never a
// credential and never trusted: its approver and digest columns are display only.

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
