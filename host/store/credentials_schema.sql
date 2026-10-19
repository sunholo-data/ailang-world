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
