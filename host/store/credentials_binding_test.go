package store

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// makeV1Credentials writes a version-1 credentials store (the five columns,
// user_version 1) holding rows: the shape of every credentials.db M1 wrote,
// the operator rig's ten live credentials among them.
func makeV1Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	db := rawDB(t, path)
	if _, err := db.Exec(canonicalCredentialsV1DDL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, credentialsSchemaV1)
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?)`,
			r.CredentialID, r.EpisodeID, r.GrantsJSON, r.ExpiresAt, r.CreatedAt); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

// tenRigRows is rigRows widened to the rig's ten live credentials.
func tenRigRows() []SessionRow {
	rows := rigRows()
	for i := 8; i < 10; i++ {
		rows = append(rows, SessionRow{
			CredentialID: fmt.Sprintf("%064x", i+1),
			EpisodeID:    fmt.Sprintf("w%d", 180+i),
			GrantsJSON:   fmt.Sprintf(`[{"Effect":"Workspace.Exec","Scope":"worktree","ExpiresAt":%d,"Budget":%d}]`, 2000+i, 50+i),
			ExpiresAt:    int64(2000 + i),
			CreatedAt:    int64(1000 + i),
		})
	}
	return rows
}

func userVersion(t *testing.T, path string) int {
	t.Helper()
	db := rawDB(t, path)
	defer func() { _ = db.Close() }()
	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

func tableColumns(t *testing.T, db *sql.DB) string {
	t.Helper()
	rows, err := db.Query(`SELECT name, type, "notnull", pk FROM pragma_table_info('session_credentials') ORDER BY cid`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var name, typ string
		var notNull, pk int
		if err := rows.Scan(&name, &typ, &notNull, &pk); err != nil {
			t.Fatal(err)
		}
		out = append(out, fmt.Sprintf("%s:%s:%d:%d", name, typ, notNull, pk))
	}
	return strings.Join(out, " ")
}

// Row 189 M2: the version 1 -> 2 upgrade keeps every existing row byte for
// byte and unbound (NULL, NULL), reaches the fresh store's columns, and is
// idempotent. Mutant: drop the ADD COLUMN statements (or stamp the version
// without them) and the bound mint after the upgrade fails; rewrite a row in
// the upgrade and the DeepEqual fires.
func TestCredentialsUpgradeV1KeepsEveryRowUnboundAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	want := tenRigRows()
	makeV1Credentials(t, path, want)

	// Control: before the upgrade a version-1 store reads, unbound.
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatalf("read-only open of a version-1 store: %v", err)
	}
	got, err := ro.ListSessions(boundedTestContext(t), "", 0)
	_ = ro.Close()
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("version-1 rows = %+v %v", got, err)
	}

	for i, wantUpgraded := range []bool{true, false, false} {
		upgraded, err := UpgradeCredentials(path)
		if err != nil || upgraded != wantUpgraded {
			t.Fatalf("UpgradeCredentials run %d = %v, %v; want %v", i, upgraded, err, wantUpgraded)
		}
		if v := userVersion(t, path); v != credentialsSchemaVersion {
			t.Fatalf("run %d: user_version %d, want %d", i, v, credentialsSchemaVersion)
		}
		ro, err := OpenCredentialsReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ro.ListSessions(boundedTestContext(t), "", 0)
		_ = ro.Close()
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d: rows after the upgrade =\n%+v\nwant\n%+v (%v)", i, got, want, err)
		}
	}
	db := rawDB(t, path)
	var unbound int
	if err := db.QueryRow(`SELECT count(*) FROM session_credentials WHERE project IS NULL AND profile_digest IS NULL`).Scan(&unbound); err != nil || unbound != len(want) {
		t.Fatalf("%d unbound rows (%v), want all %d", unbound, err, len(want))
	}
	upgradedCols := tableColumns(t, db)
	_ = db.Close()
	fresh := openTestCredentials(t)
	if freshCols := tableColumns(t, fresh.db); freshCols != upgradedCols {
		t.Fatalf("upgraded columns %q differ from a fresh store's %q", upgradedCols, freshCols)
	}

	// The upgraded store takes a bound mint and resolves the binding.
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	bound := SessionRow{CredentialID: fmt.Sprintf("%064x", 99), EpisodeID: "w199", GrantsJSON: "[]",
		ExpiresAt: 3000, CreatedAt: 2500, Project: "ailang-world", ProfileDigest: "sha256:" + strings.Repeat("ab", 32)}
	if err := c.MintSession(boundedTestContext(t), bound); err != nil {
		t.Fatalf("bound mint after the upgrade: %v", err)
	}
	if r, ok, err := c.ResolveSession(boundedTestContext(t), bound.CredentialID); err != nil || !ok || !reflect.DeepEqual(r, bound) {
		t.Fatalf("resolve bound = %+v %v %v", r, ok, err)
	}
}

// A session verb never upgrades a version-1 store (an older daemon serving it
// would refuse a newer user_version): an unbound mint lands in it as before,
// and a bound mint is refused naming the serve restart. Mutant: let
// OpenCredentials upgrade and user_version reads 2 here.
func TestCredentialsSessionVerbsLeaveAVersion1StoreAlone(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV1Credentials(t, path, rigRows())
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatalf("writable open of a version-1 store: %v", err)
	}
	defer func() { _ = c.Close() }()
	if c.Version() != credentialsSchemaV1 {
		t.Fatalf("opened at version %d, want 1", c.Version())
	}
	unbound := SessionRow{CredentialID: fmt.Sprintf("%064x", 77), EpisodeID: "w200", GrantsJSON: "[]", ExpiresAt: 3000, CreatedAt: 2500}
	if err := c.MintSession(boundedTestContext(t), unbound); err != nil {
		t.Fatalf("unbound mint into a version-1 store: %v", err)
	}
	bound := unbound
	bound.CredentialID, bound.Project, bound.ProfileDigest = fmt.Sprintf("%064x", 78), "p", "sha256:x"
	err = c.MintSession(boundedTestContext(t), bound)
	var refused *UnboundCredentialsStoreError
	if !errors.As(err, &refused) || !strings.Contains(err.Error(), "restart `ailang-worldd serve`") {
		t.Fatalf("bound mint into a version-1 store: %v, want UnboundCredentialsStoreError", err)
	}
	_ = c.Close()
	if v := userVersion(t, path); v != credentialsSchemaV1 {
		t.Fatalf("a session verb moved user_version to %d", v)
	}
}

// The binding is both-or-neither, in Go and in the table itself.
func TestCredentialsBindingNeedsProjectAndDigest(t *testing.T) {
	c := openTestCredentials(t)
	half := SessionRow{CredentialID: fmt.Sprintf("%064x", 5), EpisodeID: "w1", GrantsJSON: "[]", ExpiresAt: 2, CreatedAt: 1, Project: "p"}
	if err := c.MintSession(boundedTestContext(t), half); err == nil {
		t.Fatal("a project without a digest was minted")
	}
	if _, err := c.db.Exec(`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at, project) VALUES ('x', 'w1', '[]', 2, 1, 'p')`); err == nil {
		t.Fatal("the table accepted a project without a digest")
	}
	// Control: the full binding is accepted by the same table.
	if _, err := c.db.Exec(`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at, project, profile_digest) VALUES ('y', 'w1', '[]', 2, 1, 'p', 'sha256:d')`); err != nil {
		t.Fatalf("the table refused a full binding: %v", err)
	}
}
