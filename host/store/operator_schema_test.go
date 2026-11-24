package store

import (
	"bytes"
	"database/sql"
	"errors"
	"os"
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
	_, err = old.OperatorAction(testCtx(t), "missing", 1000)
	if oldSchema, ok := errors.AsType[*OlderCredentialsStoreError](err); !ok || oldSchema.Version != 9 || oldSchema.Need != 10 {
		t.Fatalf("schema 9 operator read: %v", err)
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

func TestOperatorEpisodeOwnerReadSchema(t *testing.T) {
	ctx := testCtx(t)
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
		t.Fatalf("absent store: %q %v", id, err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("absent store changed: %v", err)
	}
	if _, err := ReadOperatorEpisodeOwner(ctx, "", "ep"); err == nil {
		t.Fatal("accepted an empty database path")
	}
	if _, err := ReadOperatorEpisodeOwner(ctx, filepath.Dir(path), "ep"); err == nil || err.Error() != "expected a regular operator credentials file" {
		t.Fatal("accepted a directory as a database")
	}
	raw := rawDB(t, path)
	if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, raw, 9)
	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
		t.Fatalf("v9 owner lookup: %q %v", id, err)
	}
	if _, err := UpgradeCredentials(path); err != nil {
		t.Fatal(err)
	}
	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
		t.Fatalf("empty v10 owner lookup: %q %v", id, err)
	}
	ro, _, err := openOperatorReadOnly(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	var busy int
	busyErr := ro.QueryRowContext(ctx, "PRAGMA busy_timeout").Scan(&busy)
	_, writeErr := ro.ExecContext(ctx, "DELETE FROM operator_actions")
	ro.Close()
	if busyErr != nil || busy != 0 || writeErr == nil {
		t.Fatalf("deadline or read-only open drift: busy=%d %v write=%v", busy, busyErr, writeErr)
	}
	setPositiveFixtureVersion(t, raw, 11)
	if _, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); err == nil || userVersion(t, path) != 11 {
		t.Fatalf("future schema accepted or changed: %v", err)
	}
	corrupt := filepath.Join(t.TempDir(), "corrupt.db")
	if _, err := UpgradeCredentials(corrupt); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(corrupt)
	if err != nil || len(data) == 0 {
		t.Fatalf("read migrated fixture: %v", err)
	}
	data[0] ^= 0xff
	if err := os.WriteFile(corrupt, data, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadOperatorEpisodeOwner(ctx, corrupt, "ep"); err == nil {
		t.Fatal("corrupt database accepted")
	}
	if got, err := os.ReadFile(corrupt); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("corrupt database changed: %v", err)
	}
}
