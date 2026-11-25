-- Standing session allowances and draw ledger (credentials schema v11, row 224, D-WORLD-96).
CREATE TABLE IF NOT EXISTS standing_allowances (
    allowance_id    TEXT PRIMARY KEY CHECK (length(allowance_id) = 19 AND allowance_id GLOB 'sa-[0-9a-f]*'),
    preset_name     TEXT NOT NULL CHECK (length(preset_name) BETWEEN 1 AND 51),
    preset_sha256   TEXT NOT NULL CHECK (length(preset_sha256) = 64),
    profile_project TEXT CHECK (profile_project IS NULL OR profile_project <> ''),
    profile_digest  TEXT CHECK ((profile_digest IS NULL) = (profile_project IS NULL) AND (profile_digest IS NULL OR profile_digest <> '')),
    payload         BLOB NOT NULL CHECK (length(payload) BETWEEN 1 AND 131072),
    digest          TEXT NOT NULL CHECK (length(digest) = 64),
    challenge       TEXT NOT NULL CHECK (length(challenge) = 64),
    max_per_day     INTEGER NOT NULL CHECK (max_per_day BETWEEN 1 AND 10),
    ttl_s           INTEGER NOT NULL CHECK (ttl_s BETWEEN 1 AND 43200),
    approved_at     INTEGER NOT NULL CHECK (approved_at > 0),
    expires_at      INTEGER NOT NULL CHECK (expires_at > approved_at AND expires_at - approved_at <= 2592000),
    approved_by     TEXT NOT NULL CHECK (approved_by <> ''),
    artifact        BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
    revoked_at      INTEGER CHECK (revoked_at IS NULL OR revoked_at >= approved_at),
    revoked_by      TEXT CHECK ((revoked_at IS NULL) = (revoked_by IS NULL) AND (revoked_by IS NULL OR revoked_by <> '')),
    revoked_reason  TEXT CHECK (revoked_reason IS NULL OR revoked_at IS NOT NULL),
    CHECK ((approved_by = 'tty') = (artifact IS NULL))
);

CREATE UNIQUE INDEX IF NOT EXISTS standing_allowances_active_preset
    ON standing_allowances(preset_name)
    WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS allowance_draws (
    draw_id       TEXT PRIMARY KEY CHECK (length(draw_id) = 19 AND draw_id GLOB 'ad-[0-9a-f]*'),
    allowance_id  TEXT NOT NULL REFERENCES standing_allowances(allowance_id),
    day_key       TEXT NOT NULL CHECK (length(day_key) = 10 AND day_key GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]'),
    episode_id    TEXT NOT NULL CHECK (episode_id <> ''),
    credential_id TEXT NOT NULL UNIQUE CHECK (length(credential_id) = 64),
    repo          TEXT NOT NULL CHECK (repo <> ''),
    drawn_at      INTEGER NOT NULL CHECK (drawn_at > 0),
    ttl_s         INTEGER NOT NULL CHECK (ttl_s BETWEEN 1 AND 43200)
);

CREATE INDEX IF NOT EXISTS allowance_draws_daily
    ON allowance_draws(allowance_id, day_key);
