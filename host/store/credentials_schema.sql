-- credentials_schema.sql — the session credentials store (row 189, D-WORLD-85
-- C, w-world-onboarding §4). It is a separate SQLite file from the world store
-- with its own writer lock: `session new/mint/revoke` take THIS lock at the
-- terminal, and a serving daemon only ever opens it read-only, per request.
-- Its user_version is its own (credentialsSchemaVersion), independent of the
-- world store's.

-- SESSION credentials -> (episode, grants, expiry) — w-session-authority D2.
-- This table is the inbound Authorization: Bearer credential boundary minted
-- and resolved locally-first. credential_id stores ONLY sha256(token_hex) of
-- the once-printed raw token (D3); the raw token is never persisted. The PK
-- is the unique index a point lookup needs; there is NO secondary index and NO
-- revoked column (D4 revokes by DELETE, making a revoked credential "unknown").
-- The first five columns are the world store's v3/v4 table, so the one-time
-- move (MigrateCredentials) copies them by name unchanged.
--
-- project and profile_digest (version 2, row 189 M2, D-189-2 A) bind the
-- credential to the exec profile the human confirmed at the y/N: the
-- profile's project and the sha256 of its exact bytes. Both are NULL for an
-- unbound credential (every row minted before version 2, and every mint
-- without --profile); profile_digest's CHECK keeps them set together. The two
-- column definitions are byte-identical to the version 1 -> 2 upgrade's
-- ALTER TABLE ... ADD COLUMN statements (credentialsUpgradeV2), so a fresh and
-- an upgraded store hold the same table.
CREATE TABLE IF NOT EXISTS session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL,
    project TEXT CHECK (project IS NULL OR project <> ''),
    profile_digest TEXT CHECK ((profile_digest IS NULL) = (project IS NULL) AND (profile_digest IS NULL OR profile_digest <> ''))
);

-- Attended start history is independent of revocable credentials (version 3).
CREATE TABLE IF NOT EXISTS episode_starts (
    episode TEXT PRIMARY KEY CHECK (episode <> ''),
    start_json TEXT NOT NULL
);

-- Mint requests (version 4, row 201 A0). A request is inert until approved; approving consumes
-- it in the same transaction that inserts its session row. Consumed rows stay as the audit of
-- who approved which session; unconsumed rows are pruned a day after they expire.
CREATE TABLE IF NOT EXISTS mint_requests (
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
);

-- Approvers (version 5, row 201 A1b). A key enrolled at the daemon's terminal (a typed y, parent
-- R3) may approve a mint request by signature. Revoking MARKS the row and never deletes it, so a
-- name in mint_requests.approved_by ('ssh:<name>') names one key forever, and neither a name nor
-- a key is ever enrolled twice. format is checked in code (host/approval), not here, so A2 can add
-- 'webauthn' (and sign_count by ADD COLUMN) without rebuilding the table.
CREATE TABLE IF NOT EXISTS approvers (
    name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    format       TEXT NOT NULL CHECK (format <> ''),
    public_key   BLOB NOT NULL UNIQUE CHECK (length(public_key) > 0),
    enrolled_at  INTEGER NOT NULL CHECK (enrolled_at > 0),
    last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
    revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
);
