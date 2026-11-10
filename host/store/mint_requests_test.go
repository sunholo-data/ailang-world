package store

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// makeV2Credentials writes a schema-2 credentials store at path: the version-2 table, the given
// bound rows and the version-2 user_version.
func makeV2Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	db := rawDB(t, path)
	if _, err := db.Exec(canonicalCredentialsDDL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 2)
	insertBoundRawRows(t, db, rows)
}

// makeV3Credentials writes a schema-3 credentials store at path: the version-3 table set
// (episode_starts, session_credentials), the given bound rows and the version-3 user_version.
func makeV3Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	db := rawDB(t, path)
	if _, err := db.Exec(episodeStartsSQL); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(canonicalCredentialsDDL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 3)
	insertBoundRawRows(t, db, rows)
}

// insertBoundRawRows inserts rows with their profile binding, as the version 2 and 3 stores hold them.
func insertBoundRawRows(t *testing.T, db *sql.DB, rows []SessionRow) {
	t.Helper()
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at, project, profile_digest) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			r.CredentialID, r.EpisodeID, r.GrantsJSON, r.ExpiresAt, r.CreatedAt, r.Project, r.ProfileDigest); err != nil {
			t.Fatal(err)
		}
	}
}

// oneBoundRow is one bound credential, for the stores that have the binding columns.
func oneBoundRow() SessionRow {
	return SessionRow{
		CredentialID:  fmt.Sprintf("%064x", 0xb0),
		EpisodeID:     "w900",
		GrantsJSON:    `[{"Effect":"Workspace.Exec","Scope":"worktree","ExpiresAt":3000,"Budget":9}]`,
		ExpiresAt:     3000,
		CreatedAt:     1100,
		Project:       "proj-a",
		ProfileDigest: strings.Repeat("a", 64),
	}
}

// tableNamesAt lists the schema objects of the store at path, read through its own handle.
func tableNamesAt(t *testing.T, path string) []string {
	t.Helper()
	names, _ := schemaState(t, rawDB(t, path))
	return names
}

// TestCredentialsUpgradeToV4 (ST2): a version 1, 2 or 3 store is upgraded to version 4 by
// UpgradeCredentials only, keeps every credential, and gains exactly the mint_requests table.
func TestCredentialsUpgradeToV4(t *testing.T) {
	for _, v := range []int{1, 2, 3} {
		t.Run(fmt.Sprintf("v%d", v), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), CredentialsFileName)
			var rows []SessionRow
			const startJSON = `{"episode":"w900","note":"kept"}`
			switch v {
			case 1:
				rows = rigRows()
				makeV1Credentials(t, path, rows)
			case 2:
				rows = []SessionRow{oneBoundRow()}
				makeV2Credentials(t, path, rows)
			case 3:
				rows = []SessionRow{oneBoundRow()}
				makeV3Credentials(t, path, rows)
				db := rawDB(t, path)
				if _, err := db.Exec(`INSERT INTO episode_starts (episode, start_json) VALUES (?, ?)`, "w900", startJSON); err != nil {
					t.Fatal(err)
				}
			}
			tablesBefore := tableNamesAt(t, path)

			// (a) A session verb opens the store at its version and never upgrades it.
			c, err := OpenCredentials(path)
			if err != nil {
				t.Fatal(err)
			}
			if c.Version() != v {
				t.Fatalf("OpenCredentials Version() = %d, want %d", c.Version(), v)
			}
			if err := c.Close(); err != nil {
				t.Fatal(err)
			}
			if got := userVersion(t, path); got != v {
				t.Fatalf("a session open changed user_version to %d, want %d", got, v)
			}
			if got := tableNamesAt(t, path); !reflect.DeepEqual(got, tablesBefore) {
				t.Fatalf("a session open changed the tables: %v -> %v", tablesBefore, got)
			}

			// (b) The upgrade runs once.
			if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded {
				t.Fatalf("UpgradeCredentials: upgraded %v, err %v; want true, nil", upgraded, err)
			}
			if again, err := UpgradeCredentials(path); err != nil || again {
				t.Fatalf("second UpgradeCredentials: upgraded %v, err %v; want false, nil", again, err)
			}

			// (c) The result is schema 4: the exact table set, the canonical mint_requests table, and
			// every original credential (and the start row) unchanged.
			if got := userVersion(t, path); got != credentialsSchemaVersion {
				t.Fatalf("user_version = %d, want %d", got, credentialsSchemaVersion)
			}
			db := rawDB(t, path)
			ddl := tableDDL(t, db)
			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"episode_starts", "mint_requests", "session_credentials"})
			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
			}
			ro, err := OpenCredentialsReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = ro.Close() }()
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			for _, r := range rows {
				got, ok, err := ro.ResolveSession(ctx, r.CredentialID)
				if err != nil || !ok || !reflect.DeepEqual(got, r) {
					t.Fatalf("credential %s after upgrade: ok %v, err %v, got %+v; want %+v", r.CredentialID, ok, err, got, r)
				}
			}
			if v == 3 {
				var raw string
				if err := db.QueryRow(`SELECT start_json FROM episode_starts WHERE episode = ?`, "w900").Scan(&raw); err != nil || raw != startJSON {
					t.Fatalf("episode start after upgrade: %q, err %v; want %q", raw, err, startJSON)
				}
			}
		})
	}
}
