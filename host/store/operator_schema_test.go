package store

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"
)

func operatorSchema(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT type || ':' || name || ':' || coalesce(sql, '') FROM sqlite_schema ORDER BY type, name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var schema []string
	for rows.Next() {
		var entry string
		if err := rows.Scan(&entry); err != nil {
			t.Fatal(err)
		}
		schema = append(schema, normalizeDDL(entry))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return schema
}

func TestOperatorSchema10(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	raw := rawDB(t, path)
	if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, raw, 9)
	old, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	row := oneBoundRow()
	if old.Version() != 9 || old.MintSession(testCtx(t), row) != nil {
		t.Fatal("v9 terminal open must retain its version and bound writes")
	}
	if err := old.Close(); err != nil {
		t.Fatal(err)
	}
	// Fail after CREATE TABLE succeeds. Partial DDL and version must roll back.
	if _, err := raw.Exec(`CREATE INDEX operator_active_episode ON session_credentials(episode_id)`); err != nil {
		t.Fatal(err)
	}
	before := operatorSchema(t, raw)
	if ok, err := UpgradeCredentials(path); err == nil || ok {
		t.Fatalf("colliding migration: %v %v", ok, err)
	}
	if userVersion(t, path) != 9 || !reflect.DeepEqual(before, operatorSchema(t, raw)) {
		t.Fatal("failed migration changed schema")
	}
	if _, err := raw.Exec(`DROP INDEX operator_active_episode`); err != nil {
		t.Fatal(err)
	}
	if ok, err := UpgradeCredentials(path); err != nil || !ok || userVersion(t, path) != 10 {
		t.Fatalf("v9 upgrade: %v %v", ok, err)
	}
	if ok, err := UpgradeCredentials(path); err != nil || ok {
		t.Fatalf("repeat upgrade: %v %v", ok, err)
	}
	fresh := openTestCredentials(t)
	if fresh.Version() != 10 || !reflect.DeepEqual(operatorSchema(t, raw), operatorSchema(t, fresh.db)) {
		t.Fatal("fresh and upgraded SQLite schema differ")
	}
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	got, ok, err := ro.ResolveSession(testCtx(t), row.CredentialID)
	if err != nil || !ok || got != row {
		t.Fatalf("lost v9 credential: %+v %v %v", got, ok, err)
	}
}
