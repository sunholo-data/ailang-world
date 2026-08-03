package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

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
