package store

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
