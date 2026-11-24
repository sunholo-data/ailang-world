-- Credentials schema 10. Used for both a fresh store and the v9 upgrade.
CREATE TABLE operator_actions (
    id          TEXT NOT NULL PRIMARY KEY,
    episode_id  TEXT NOT NULL,
    kind        TEXT NOT NULL CHECK (kind IN ('start', 'finish')),
    review      BLOB NOT NULL CHECK (length(review) = 32),
    digest      BLOB NOT NULL CHECK (length(digest) = 32),
    action      BLOB,
    challenge   BLOB CHECK (length(challenge) = 32),
    created_at  INTEGER NOT NULL,
    expires_at  INTEGER NOT NULL CHECK (expires_at > created_at),
    state       TEXT NOT NULL CHECK (state IN ('requested', 'confirmed', 'claimed', 'running',
                 'awaiting-pr-confirmation', 'succeeded', 'failed', 'uncertain', 'expired', 'cancelled')),
    generation  INTEGER NOT NULL,
    owner       TEXT,
    approved_by TEXT,
    artifact    BLOB,
    terminal_at INTEGER,
    v3_id       TEXT UNIQUE,
    CHECK (terminal_at IS NOT NULL OR (action IS NOT NULL AND challenge IS NOT NULL)),
    CHECK ((approved_by IS NULL) = (artifact IS NULL)),
    CHECK (state <> 'requested' OR (approved_by IS NULL AND owner IS NULL)),
    CHECK (state NOT IN ('confirmed', 'claimed', 'running', 'awaiting-pr-confirmation', 'uncertain')
           OR approved_by IS NOT NULL),
    CHECK ((owner IS NOT NULL) = (state IN ('claimed', 'running', 'awaiting-pr-confirmation', 'uncertain'))),
    CHECK ((terminal_at IS NOT NULL) = (state IN ('succeeded', 'failed', 'expired', 'cancelled'))),
    CHECK (v3_id IS NULL OR kind = 'finish'),
    CHECK (state <> 'awaiting-pr-confirmation' OR v3_id IS NOT NULL)
);
CREATE UNIQUE INDEX operator_active_episode ON operator_actions(episode_id)
WHERE terminal_at IS NULL;
CREATE INDEX operator_review ON operator_actions(episode_id, kind, review);

CREATE TABLE operator_runner_state (
    singleton  INTEGER PRIMARY KEY CHECK (singleton = 1),
    generation INTEGER NOT NULL,
    owner      TEXT NOT NULL,
    protocol   INTEGER NOT NULL,
    build      TEXT NOT NULL,
    seen_at    INTEGER NOT NULL,
    readiness  TEXT NOT NULL,
    catalogue_generation INTEGER NOT NULL,
    catalogue BLOB
);

-- Operations use the same identity ledger as effect reservations. An operation
-- alias carries its caller's operation id even if another tap created the action.
CREATE TABLE operator_action_steps (
    kind       TEXT NOT NULL CHECK (kind IN ('operation', 'worktree', 'mint', 'files', 'launch',
                 'herdr-handoff', 'checks', 'commit', 'v3-card', 'push', 'label', 'pull-request', 'record', 'cleanup')),
    identity   TEXT NOT NULL,
    action_id  TEXT NOT NULL REFERENCES operator_actions(id),
    phase      TEXT NOT NULL CHECK (phase IN ('reserved', 'observed')),
    version    INTEGER NOT NULL,
    owner      TEXT,
    created_at INTEGER NOT NULL,
    receipt    BLOB,
    PRIMARY KEY (kind, identity),
    CHECK (kind <> 'operation' OR (phase = 'observed' AND owner IS NULL AND receipt IS NULL)),
    CHECK (kind = 'operation' OR owner IS NOT NULL),
    CHECK (phase <> 'reserved' OR receipt IS NULL)
);
