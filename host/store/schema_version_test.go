package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

const frozenFutureSchemaVersion = 3

const expectedCurrentSchemaVersion = 2

// schemaV1SQL is copied verbatim from ad619d8:host/store/schema.sql. The
// artifact SHA-256 is 13893a296c394cc2c5b3997e8fff729467dc9ac83a03b458796634aa52fb5436.
// Never mechanically update this fixture to follow current DDL.
const schemaV1SQL = `-- schema.sql — SQLite schema for the M1 semantic world store (Decision 4).
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
`

// schemaV2SQL is an independently authored ledger entry: the frozen v1 schema
// plus the reviewed v2 approval-claims DDL. It must never reference schemaSQL.
const schemaV2SQL = schemaV1SQL + `

CREATE TABLE IF NOT EXISTS approval_claims (
    approval_ref TEXT PRIMARY KEY,
    request_ref  TEXT NOT NULL,
    invocation_id TEXT NOT NULL UNIQUE
);
`

const applicationObjectCountSQL = `SELECT count(*) FROM sqlite_master WHERE name NOT LIKE 'sqlite\_%' ESCAPE '\'`

func rawDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func schemaState(t *testing.T, db *sql.DB) ([]string, int) {
	t.Helper()
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE name NOT LIKE 'sqlite\_%' ESCAPE '\' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	return names, version
}

func requireFreshControl(t *testing.T) {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "fresh-control.db"))
	if err != nil {
		t.Fatalf("fresh positive control: %v", err)
	}
	_ = s.Close()
}

func TestApplicationObjectPredicateUsesLiteralReservedPrefix(t *testing.T) {
	db := rawDB(t, filepath.Join(t.TempDir(), "predicate.db"))
	if _, err := db.Exec(`CREATE TABLE sqliteX_probe(id INTEGER PRIMARY KEY AUTOINCREMENT); CREATE TABLE realapp(id INTEGER);`); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRow(applicationObjectCountSQL).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("application object count = %d, want 2 (sqliteX_probe and realapp; exclude sqlite_sequence and autoindexes)", count)
	}
}

func TestSchemaVersionErrorMessages(t *testing.T) {
	cases := []struct {
		err  error
		want string
	}{
		{&LegacySchemaVersionError{Path: "x", Found: 1, Current: 2}, `store: schema version legacy: "x" has user_version 1 with application schema; binary requires 2; refusing to modify`},
		{&FutureSchemaVersionError{Path: "x", Found: 3, Current: 2}, `store: schema version future: "x" has user_version 3; binary supports 2; use a compatible binary`},
		{&InvalidSchemaVersionError{Path: "x", Found: -1, Current: 2}, `store: schema version invalid: "x" has negative user_version -1; binary requires 2; refusing to modify`},
		{&UninitializedReadOnlyStoreError{Path: "x"}, `store: schema uninitialized: read-only store "x" has no application schema; open writable once to initialize`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("error %q lacks exact substring %q", tc.err, tc.want)
		}
	}
}

func TestFreshWriterInitializesSchemaAndVersionTwo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fresh.db")
	s, err := Open(path)
	if err != nil {
		db := rawDB(t, path)
		_, version := schemaState(t, db)
		if version != 2 {
			t.Fatalf("fresh user_version = %d, want 2 (Open error: %v)", version, err)
		}
		t.Fatal(err)
	}
	defer s.Close()
	ddl := tableDDL(t, s.db)
	if len(ddl) != 8 {
		t.Fatalf("fresh table count = %d, want 8", len(ddl))
	}
	_, version := schemaState(t, s.db)
	if version != 2 {
		t.Fatalf("fresh user_version = %d, want 2", version)
	}
}

func TestFreshInitializationIsAtomic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "atomic.db")
	db := rawDB(t, path)
	sentinel := errors.New("induced version write failure")
	err := freshInitTx(db, func(*sql.Tx) error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("freshInitTx error = %v, want wrapped sentinel", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	var count, version int
	if err := reopened.QueryRow(applicationObjectCountSQL).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if err := reopened.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if count != 0 || version != 0 {
		t.Fatalf("rolled-back state = (%d objects, version %d), want (0, 0)", count, version)
	}
}

func setPositiveFixtureVersion(t *testing.T, db *sql.DB, version int) {
	t.Helper()
	if version < 1 || version > 2147483647 {
		t.Fatalf("positive fixture version %d outside 1..2147483647", version)
	}
	if _, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", version)); err != nil {
		t.Fatal(err)
	}
}

func setNegativeFixtureVersion(t *testing.T, db *sql.DB, version int) {
	t.Helper()
	if version >= 0 {
		t.Fatalf("negative fixture version = %d", version)
	}
	if _, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", version)); err != nil {
		t.Fatal(err)
	}
	var got int
	if err := db.QueryRow("PRAGMA user_version").Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != version {
		t.Fatalf("negative fixture user_version = %d, want %d", got, version)
	}
}

func TestSupportedVersionTwoOpensWithoutRewritingPragma(t *testing.T) {
	path := filepath.Join(t.TempDir(), "supported.db")
	db := rawDB(t, path)
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 2)
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	afterNames, afterVersion := schemaState(t, s.db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion || afterVersion != 2 {
		t.Fatalf("supported state changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
}

func TestVersionOneStoreIsRejectedUnmodifiedByWriterAndReader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "version-one.db")
	db := rawDB(t, path)
	if _, err := db.Exec(schemaV1SQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 1)
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	for _, open := range []struct {
		name string
		fn   func(string) (*Store, error)
	}{{"Open", Open}, {"OpenReadOnly", OpenReadOnly}} {
		t.Run(open.name, func(t *testing.T) {
			s, err := open.fn(path)
			if s != nil {
				_ = s.Close()
				t.Fatalf("%s returned a store", open.name)
			}
			var legacy *LegacySchemaVersionError
			if !errors.As(err, &legacy) || legacy.Found != 1 || legacy.Current != 2 {
				t.Fatalf("%s error = %#v, want legacy Found=1 Current=2", open.name, err)
			}
		})
	}
	db = rawDB(t, path)
	afterNames, afterVersion := schemaState(t, db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
		t.Fatalf("version-one fixture changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
	setPositiveFixtureVersion(t, db, 2)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("Open after correcting fixture (writer-lock release probe): %v", err)
	}
	if err := reopened.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFutureVersionIsRejected(t *testing.T) {
	requireFreshControl(t)
	path := filepath.Join(t.TempDir(), "future.db")
	db := rawDB(t, path)
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, frozenFutureSchemaVersion)
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if s != nil {
		_ = s.Close()
		t.Fatal("future Open returned a store")
	}
	var future *FutureSchemaVersionError
	if !errors.As(err, &future) || future.Found != 3 || future.Current != 2 {
		t.Fatalf("Open error = %#v, want future Found=3 Current=2", err)
	}
	if !strings.Contains(err.Error(), "schema version future") {
		t.Fatalf("future message = %q", err)
	}
	db = rawDB(t, path)
	afterNames, afterVersion := schemaState(t, db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
		t.Fatalf("future changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
}

func TestNegativeVersionsAreInvalid(t *testing.T) {
	for _, version := range []int{-1, -2147483648} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			requireFreshControl(t)
			path := filepath.Join(t.TempDir(), "invalid.db")
			db := rawDB(t, path)
			if _, err := db.Exec(schemaSQL); err != nil {
				t.Fatal(err)
			}
			setNegativeFixtureVersion(t, db, version)
			beforeNames, beforeVersion := schemaState(t, db)
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			s, err := Open(path)
			if s != nil {
				_ = s.Close()
				t.Fatal("invalid Open returned a store")
			}
			var invalid *InvalidSchemaVersionError
			if !errors.As(err, &invalid) || invalid.Found != version || invalid.Current != 2 {
				t.Fatalf("Open error = %#v, want invalid Found=%d Current=2", err, version)
			}
			if !strings.Contains(err.Error(), "schema version invalid") {
				t.Fatalf("invalid message = %q", err)
			}
			db = rawDB(t, path)
			afterNames, afterVersion := schemaState(t, db)
			if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
				t.Fatalf("invalid changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
			}
		})
	}
}

func TestSchemaVersionRangeBounds(t *testing.T) {
	if currentSchemaVersion < 1 || currentSchemaVersion > 2147483647 {
		t.Fatalf("currentSchemaVersion = %d, want 1..2147483647", currentSchemaVersion)
	}
	if frozenFutureSchemaVersion < 1 || frozenFutureSchemaVersion > 2147483647 {
		t.Fatalf("frozen future version = %d, want 1..2147483647", frozenFutureSchemaVersion)
	}
}

func TestOpenReadOnlyEnforcesVersionTwo(t *testing.T) {
	lockControlPath := filepath.Join(t.TempDir(), "writer-lock-control.db")
	control, err := Open(lockControlPath)
	if err != nil {
		t.Fatalf("writable lock positive control: %v", err)
	}
	if _, err := os.Stat(lockControlPath + writerLockSuffix); err != nil {
		_ = control.Close()
		t.Fatalf("writable Open did not create expected lock file: %v", err)
	}
	if err := control.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		setup     func(*testing.T, *sql.DB)
		emptyFile bool
		wantType  string
		accept    bool
	}{
		{name: "version2", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(schemaSQL)
			if err != nil {
				t.Fatal(err)
			}
			setPositiveFixtureVersion(t, db, 2)
		}, accept: true},
		{name: "legacy", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(preJournalSchemaV0)
			if err != nil {
				t.Fatal(err)
			}
		}, wantType: "legacy"},
		{name: "sqliteX_probe", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(`CREATE TABLE sqliteX_probe(id INTEGER)`)
			if err != nil {
				t.Fatal(err)
			}
		}, wantType: "legacy"},
		{name: "future", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(schemaSQL)
			if err != nil {
				t.Fatal(err)
			}
			setPositiveFixtureVersion(t, db, 3)
		}, wantType: "future"},
		{name: "negative1", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(schemaSQL)
			if err != nil {
				t.Fatal(err)
			}
			setNegativeFixtureVersion(t, db, -1)
		}, wantType: "invalid"},
		{name: "int32min", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(schemaSQL)
			if err != nil {
				t.Fatal(err)
			}
			setNegativeFixtureVersion(t, db, -2147483648)
		}, wantType: "invalid"},
		{name: "empty", emptyFile: true, wantType: "uninitialized"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "readonly.db")
			var beforeNames []string
			var beforeVersion int
			if tc.emptyFile {
				if err := os.WriteFile(path, nil, 0o600); err != nil {
					t.Fatal(err)
				}
			} else {
				db := rawDB(t, path)
				tc.setup(t, db)
				beforeNames, beforeVersion = schemaState(t, db)
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}
			ro, err := OpenReadOnly(path)
			if tc.accept {
				if err != nil {
					t.Fatal(err)
				}
				if ro == nil {
					t.Fatal("OpenReadOnly returned nil store")
				}
				_ = ro.Close()
			} else {
				if ro != nil {
					_ = ro.Close()
					t.Fatal("OpenReadOnly returned a store for rejected fixture")
				}
				switch tc.wantType {
				case "legacy":
					var target *LegacySchemaVersionError
					if !errors.As(err, &target) || !strings.Contains(err.Error(), "schema version legacy") {
						t.Fatalf("error = %T %v, want legacy", err, err)
					}
				case "future":
					var target *FutureSchemaVersionError
					if !errors.As(err, &target) || !strings.Contains(err.Error(), "schema version future") {
						t.Fatalf("error = %T %v, want future", err, err)
					}
				case "invalid":
					var target *InvalidSchemaVersionError
					if !errors.As(err, &target) || !strings.Contains(err.Error(), "schema version invalid") {
						t.Fatalf("error = %T %v, want invalid", err, err)
					}
				case "uninitialized":
					var target *UninitializedReadOnlyStoreError
					if !errors.As(err, &target) || !strings.Contains(err.Error(), "schema uninitialized") {
						t.Fatalf("error = %T %v, want uninitialized", err, err)
					}
				}
				db := rawDB(t, path)
				afterNames, afterVersion := schemaState(t, db)
				if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
					t.Fatalf("read-only rejection changed state: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
				}
			}
			if _, err := os.Stat(path + writerLockSuffix); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("OpenReadOnly created writer lock %q: %v", path+writerLockSuffix, err)
			}
		})
	}
}

func TestSchemaVersionLedgerIsIndependent(t *testing.T) {
	source, err := os.ReadFile("schema_version_test.go")
	if err != nil {
		t.Fatalf("read ledger source: %v", err)
	}
	// The authored-ledger check must be ANCHORED TO LINE START, and its needle must
	// not be a single literal.
	//
	// This file greps ITSELF, so any needle written as one literal is present in the
	// source by virtue of the check that looks for it — a self-referential control
	// that passes no matter what the declaration says. The two negative checks below
	// are split (`"var schemaV2SQL = "+"schemaSQL"`) because a whole literal there
	// would make them fire ALWAYS; the positive check had no such pressure and was
	// first written as one literal, which made it fire NEVER.
	//
	// Measured, not assumed: with the positive check written as the plain literal
	// `"const schemaV2SQL = schemaV1SQL +"`, replacing the declaration with
	// `var schemaV2SQL = string(schemaSQL)` — the ledger becoming the very file it
	// is supposed to independently attest — left this test reporting `ok` in 0.290s.
	// Both negative needles were dodged by the `string(...)` conversion, and the DDL
	// comparison below then compares schema.sql against itself and trivially agrees.
	//
	// The repair is the same one the mission's charter rotation uses: anchor to
	// `^`, so prose and check-lines (which are indented inside this func) cannot
	// satisfy it. `(?m)` makes `^` match at each line start.
	ledgerDecl := regexp.MustCompile(`(?m)^const schemaV2SQL = ` + `schemaV1SQL \+`)
	if !ledgerDecl.Match(source) {
		t.Fatal("schemaV2SQL is not an authored constant extending the frozen v1 ledger")
	}
	if strings.Contains(string(source), "var schemaV2SQL = "+"schemaSQL") {
		t.Fatal("schemaV2SQL is derived from schemaSQL")
	}
	// Semantic backstop, immune to every needle game above: a ledger that IS the
	// schema under test cannot attest to it. This catches any derivation the
	// source-text checks miss, including forms nobody has thought of yet.
	if schemaV2SQL == schemaSQL {
		t.Fatal("schemaV2SQL is byte-identical to schemaSQL — the ledger is the file under test, so this gate is a no-op")
	}
	if strings.Contains(string(source), "frozenDB.Exec("+"schemaSQL)") {
		t.Fatal("version-2 ledger fixture is derived from schemaSQL")
	}
	if currentSchemaVersion != expectedCurrentSchemaVersion {
		t.Fatalf("currentSchemaVersion = %d, frozen expectation = %d", currentSchemaVersion, expectedCurrentSchemaVersion)
	}
	frozenDB := rawDB(t, filepath.Join(t.TempDir(), "frozen-v2.db"))
	if _, err := frozenDB.Exec(schemaV2SQL); err != nil {
		t.Fatal(err)
	}
	frozenDDL := tableDDL(t, frozenDB)

	current, err := Open(filepath.Join(t.TempDir(), "current.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer current.Close()
	currentDDL := tableDDL(t, current.db)

	frozenNames := sortedDDLNames(frozenDDL)
	currentNames := sortedDDLNames(currentDDL)
	requireExactTableNames(t, "frozen version 2", frozenDDL, currentNames)
	requireExactTableNames(t, "current schema", currentDDL, frozenNames)
	if !reflect.DeepEqual(frozenNames, currentNames) {
		t.Fatalf("version-2 table names = %v, current names = %v", frozenNames, currentNames)
	}
	for _, name := range frozenNames {
		got, want := normalizeDDL(currentDDL[name]), normalizeDDL(frozenDDL[name])
		if got != want {
			t.Fatalf("version-2 ledger DDL mismatch for table %q:\n got current: %s\nwant frozen: %s", name, got, want)
		}
	}
}

func TestLegacyVersionZeroStoreIsRejectedUnmodified(t *testing.T) {
	requireFreshControl(t)
	path := filepath.Join(t.TempDir(), "legacy.db")
	db := rawDB(t, path)
	if _, err := db.Exec(preJournalSchemaV0); err != nil {
		t.Fatal(err)
	}
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if s != nil {
		_ = s.Close()
		t.Fatal("legacy Open returned a store")
	}
	var legacy *LegacySchemaVersionError
	if !errors.As(err, &legacy) || legacy.Found != 0 || legacy.Current != 2 {
		t.Fatalf("Open error = %#v, want legacy Found=0 Current=2", err)
	}
	if !strings.Contains(err.Error(), "schema version legacy") {
		t.Fatalf("legacy message = %q", err)
	}
	db = rawDB(t, path)
	afterNames, afterVersion := schemaState(t, db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
		t.Fatalf("legacy changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
	if _, ok := tableDDL(t, db)["journal"]; ok {
		t.Fatal("journal created in legacy store")
	}
	// Correct this fixture out-of-band only to prove the rejected writer did
	// not strand its lock. This is a test-local cleanup probe, not a supported
	// operator remedy for legacy stores.
	setPositiveFixtureVersion(t, db, 2)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("Open after correcting test fixture (writer-lock release probe): %v", err)
	}
	if err := reopened.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestSqliteXProbeVersionZeroStoreIsLegacy(t *testing.T) {
	requireFreshControl(t)
	path := filepath.Join(t.TempDir(), "probe.db")
	db := rawDB(t, path)
	if _, err := db.Exec(`CREATE TABLE sqliteX_probe(id INTEGER)`); err != nil {
		t.Fatal(err)
	}
	beforeNames, beforeVersion := schemaState(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if s != nil {
		_ = s.Close()
		t.Fatal("sqliteX_probe Open returned a store")
	}
	var legacy *LegacySchemaVersionError
	if !errors.As(err, &legacy) {
		t.Fatalf("Open error = %T %v, want legacy", err, err)
	}
	db = rawDB(t, path)
	afterNames, afterVersion := schemaState(t, db)
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion {
		t.Fatalf("probe changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
	}
	if _, ok := tableDDL(t, db)["sqliteX_probe"]; !ok {
		t.Fatal("sqliteX_probe did not survive rejected Open")
	}
}
