package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

const frozenFutureSchemaVersion = 2

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
		{&LegacySchemaVersionError{Path: "x", Found: 0, Current: 1}, `store: schema version legacy: "x" has user_version 0 with application schema; binary requires 1; refusing to modify`},
		{&FutureSchemaVersionError{Path: "x", Found: 2, Current: 1}, `store: schema version future: "x" has user_version 2; binary supports 1; use a compatible binary`},
		{&InvalidSchemaVersionError{Path: "x", Found: -1, Current: 1}, `store: schema version invalid: "x" has negative user_version -1; binary requires 1; refusing to modify`},
		{&UninitializedReadOnlyStoreError{Path: "x"}, `store: schema uninitialized: read-only store "x" has no application schema; open writable once to initialize`},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("error %q lacks exact substring %q", tc.err, tc.want)
		}
	}
}

func TestFreshWriterInitializesSchemaAndVersionOne(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fresh.db")
	s, err := Open(path)
	if err != nil {
		db := rawDB(t, path)
		_, version := schemaState(t, db)
		if version != 1 {
			t.Fatalf("fresh user_version = %d, want 1 (Open error: %v)", version, err)
		}
		t.Fatal(err)
	}
	defer s.Close()
	ddl := tableDDL(t, s.db)
	if len(ddl) != 7 {
		t.Fatalf("fresh table count = %d, want 7", len(ddl))
	}
	_, version := schemaState(t, s.db)
	if version != 1 {
		t.Fatalf("fresh user_version = %d, want 1", version)
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

func TestSupportedVersionOneOpensWithoutRewritingPragma(t *testing.T) {
	path := filepath.Join(t.TempDir(), "supported.db")
	db := rawDB(t, path)
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 1)
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
	if !reflect.DeepEqual(afterNames, beforeNames) || afterVersion != beforeVersion || afterVersion != 1 {
		t.Fatalf("supported state changed: before=(%v,%d) after=(%v,%d)", beforeNames, beforeVersion, afterNames, afterVersion)
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
	if !errors.As(err, &future) || future.Found != 2 || future.Current != 1 {
		t.Fatalf("Open error = %#v, want future Found=2 Current=1", err)
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
			if !errors.As(err, &invalid) || invalid.Found != version || invalid.Current != 1 {
				t.Fatalf("Open error = %#v, want invalid Found=%d Current=1", err, version)
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

func TestOpenReadOnlyEnforcesVersionOne(t *testing.T) {
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
		{name: "version1", setup: func(t *testing.T, db *sql.DB) {
			_, err := db.Exec(schemaSQL)
			if err != nil {
				t.Fatal(err)
			}
			setPositiveFixtureVersion(t, db, 1)
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
			setPositiveFixtureVersion(t, db, 2)
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
	if !errors.As(err, &legacy) || legacy.Found != 0 || legacy.Current != 1 {
		t.Fatalf("Open error = %#v, want legacy Found=0 Current=1", err)
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
