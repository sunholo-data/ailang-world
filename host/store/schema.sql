-- schema.sql — SQLite schema for the M1 semantic world store (Decision 4).
--
-- Every HashRef occupies exactly one canonical TEXT column in "algo:digest"
-- form. The Go boundary parses each value into host/hashref.HashRef before use,
-- giving one atomic indexed identity per digest and avoiding split-column
-- comparison mistakes. Algorithm-specific validation stays in the dispatcher so
-- future tags coexist in the same tables.

-- Immutable content-addressed objects: the ratified ObjectEnvelope (Decision 3)
-- plus the exact payload bytes addressed by hash_ref. hash_ref is the primary
-- identity; interface_hash_ref is the hash of the typed interface/schema bytes.
-- semantic_id and provenance are UTF-8 labels, not digest fields.
CREATE TABLE IF NOT EXISTS objects (
    hash_ref           TEXT PRIMARY KEY,
    interface_hash_ref TEXT NOT NULL,
    semantic_id        TEXT NOT NULL,
    provenance         TEXT NOT NULL,
    payload            BLOB NOT NULL
);

-- Immutable world revisions. Each world_ref addresses one revision; state_root
-- is the state object and log_head is the append-only log head at that revision.
CREATE TABLE IF NOT EXISTS worlds (
    world_ref  TEXT PRIMARY KEY,
    revision   INTEGER NOT NULL,
    state_root TEXT NOT NULL,
    log_head   TEXT NOT NULL
);

-- Append-only log. The six frozen LogHeader fields are stored verbatim:
--   entry_index, semantics_epoch, transition_fn_ref, interpreter_ref,
--   prev_entry_hash_ref, written_by.
-- transition_ref points to the content-addressed transition body and is OUTSIDE
-- the frozen header. entry_hash_ref addresses the canonical encoded
-- header-plus-body-reference bytes and is UNIQUE across the log.
CREATE TABLE IF NOT EXISTS log_entries (
    entry_index         INTEGER PRIMARY KEY,
    entry_hash_ref      TEXT NOT NULL UNIQUE,
    semantics_epoch     INTEGER NOT NULL,
    transition_fn_ref   TEXT NOT NULL,
    interpreter_ref     TEXT NOT NULL,
    prev_entry_hash_ref TEXT NOT NULL,
    written_by          TEXT NOT NULL,
    transition_ref      TEXT NOT NULL
);

-- Current immutable registry object reference, keyed by registry name (for
-- example "world/epoch-registry/v1"). object_ref addresses the selected
-- revision's immutable registry object.
CREATE TABLE IF NOT EXISTS epoch_registry_heads (
    registry_name TEXT PRIMARY KEY,
    object_ref    TEXT NOT NULL
);

-- The store's mutable selected-world-head pointer, keyed by a fixed head_key.
-- Unlike every other table this is NOT content-addressed: it is the single
-- compare-and-append serialization point (Decision 4). Commit reads world_ref
-- here under the transaction and advances it; a stale observed head yields a
-- ConflictError. M1 uses exactly one row (head_key = "selected_world_head").
CREATE TABLE IF NOT EXISTS store_heads (
    head_key  TEXT PRIMARY KEY,
    world_ref TEXT NOT NULL
);

-- Cached typecheck/verify result, keyed EXACTLY by the pair
-- (transition_fn_ref, interpreter_ref). semantics_epoch is copied in as
-- diagnostic/migration metadata only; it is NOT part of the cache key, so an
-- epoch-only change preserves the selected row as metadata-compatible.
CREATE TABLE IF NOT EXISTS verification_cache (
    transition_fn_ref TEXT NOT NULL,
    interpreter_ref   TEXT NOT NULL,
    semantics_epoch   INTEGER NOT NULL,
    verified          INTEGER NOT NULL,
    result_detail     TEXT NOT NULL,
    PRIMARY KEY (transition_fn_ref, interpreter_ref)
);

-- Durable ordered index for content-addressed intent and outcome payloads.
-- UNIQUE(invocation_id, kind) is also the lookup index used by receipt reads;
-- no additional index is needed.
CREATE TABLE IF NOT EXISTS journal (
    seq           INTEGER PRIMARY KEY,
    kind          TEXT NOT NULL CHECK (kind IN ('intent','outcome')),
    invocation_id TEXT NOT NULL CHECK (invocation_id <> ''),
    object_ref    TEXT NOT NULL CHECK (object_ref <> ''),
    UNIQUE (invocation_id, kind)
);
