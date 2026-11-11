package store

// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.10): passkeys (credentials
// schema 6). A webauthn approver is an approvers row (format webauthn, public_key the canonical
// COSE key) and one webauthn_credentials row, written together by EnrolPasskey after the typed
// enrolment code and y. Like approvers, a passkey row is never deleted; its sign_count only moves
// forward, inside the consume that uses it. The store does not import host/approval: these checks
// are the store's own backstop for the pure law the CLI runs first.

// ---- M4a ----

const passkeysTable = "webauthn_credentials"

// passkeysSchemaVersion is the first credentials schema with webauthn_credentials.
const passkeysSchemaVersion = 6

// passkeysSQL is the version 5 -> 6 upgrade: credentials_schema.sql's webauthn_credentials table,
// column for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
const passkeysSQL = `CREATE TABLE webauthn_credentials (
    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
    origin          TEXT NOT NULL CHECK (origin <> ''),
    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
)`
