package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// makeV4Credentials writes a schema-4 credentials store at path: the version-4 table set
// (episode_starts, mint_requests, session_credentials), the given bound rows and user_version 4.
func makeV4Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	makeV3Credentials(t, path, rows)
	db := rawDB(t, path)
	if _, err := db.Exec(mintRequestsSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 4)
}

// approverCount counts the approvers rows of c's store.
func approverCount(t *testing.T, c *Credentials) int {
	t.Helper()
	var n int
	if err := c.db.QueryRow(`SELECT count(*) FROM approvers`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// testApprover is a valid new row: Name name, Format "ssh", PublicKey []byte("key-" + name), EnrolledAt at.
func testApprover(name string, at int64) ApproverRow {
	return ApproverRow{Name: name, Format: "ssh", PublicKey: []byte("key-" + name), EnrolledAt: at}
}

// AS5: the version 1, 3 and 4 stores upgrade to 5 in one step, keeping every session row, the
// start history and the mint request row, and gaining an empty approvers table.
func TestCredentialsUpgradeToV5(t *testing.T) {
	const reqID = "mr-00000000000000aa"
	for _, tc := range []struct {
		name string
		v    int
		make func(t *testing.T, path string) []SessionRow
	}{
		{"v1", 1, func(t *testing.T, path string) []SessionRow {
			rows := rigRows()
			makeV1Credentials(t, path, rows)
			return rows
		}},
		{"v3", 3, func(t *testing.T, path string) []SessionRow {
			rows := []SessionRow{oneBoundRow()}
			makeV3Credentials(t, path, rows)
			if _, err := rawDB(t, path).Exec(`INSERT INTO episode_starts (episode, start_json) VALUES ('w900', '{"episode":"w900","note":"kept"}')`); err != nil {
				t.Fatal(err)
			}
			return rows
		}},
		{"v4", 4, func(t *testing.T, path string) []SessionRow {
			rows := []SessionRow{oneBoundRow()}
			makeV4Credentials(t, path, rows)
			db := rawDB(t, path)
			if _, err := db.Exec(`INSERT INTO episode_starts (episode, start_json) VALUES ('w900', '{"episode":"w900","note":"kept"}')`); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(`INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, x'7b7d', ?, ?, 'ep', 1000, 1600)`,
				reqID, strings.Repeat("ab", 32), strings.Repeat("cd", 32)); err != nil {
				t.Fatal(err)
			}
			return rows
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), CredentialsFileName)
			rows := tc.make(t, path)

			// (a) A session verb never upgrades: the version and the table set are unchanged.
			c, err := OpenCredentials(path)
			if err != nil {
				t.Fatal(err)
			}
			if got := c.Version(); got != tc.v {
				t.Fatalf("Version() = %d, want %d", got, tc.v)
			}
			before := tableNamesAt(t, path)
			_ = c.Close()
			if got := userVersion(t, path); got != tc.v {
				t.Fatalf("user_version after a session open = %d, want %d", got, tc.v)
			}
			if after := tableNamesAt(t, path); !reflect.DeepEqual(after, before) {
				t.Fatalf("a session open changed the tables: %v -> %v", before, after)
			}

			// (b) One upgrade, then none.
			if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded {
				t.Fatalf("UpgradeCredentials: upgraded %v, err %v; want true, nil", upgraded, err)
			}
			if again, err := UpgradeCredentials(path); err != nil || again {
				t.Fatalf("second UpgradeCredentials: upgraded %v, err %v; want false, nil", again, err)
			}

			// (c) Schema 5: the exact table set, the canonical approvers table, no approvers, and every
			// original row unchanged.
			if got := userVersion(t, path); got != credentialsSchemaVersion {
				t.Fatalf("user_version = %d, want %d", got, credentialsSchemaVersion)
			}
			db := rawDB(t, path)
			ddl := tableDDL(t, db)
			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials"})
			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
			}
			ro, err := OpenCredentialsReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = ro.Close() }()
			if got := approverCount(t, ro); got != 0 {
				t.Fatalf("approvers after the upgrade = %d, want 0", got)
			}
			ctx := testCtx(t)
			for _, r := range rows {
				got, ok, err := ro.ResolveSession(ctx, r.CredentialID)
				if err != nil || !ok || !reflect.DeepEqual(got, r) {
					t.Fatalf("session %s after the upgrade: got %+v ok %v err %v; want %+v", r.CredentialID, got, ok, err, r)
				}
			}
			if tc.v >= 3 {
				var start string
				if err := db.QueryRow(`SELECT start_json FROM episode_starts WHERE episode = 'w900'`).Scan(&start); err != nil {
					t.Fatal(err)
				}
				if want := `{"episode":"w900","note":"kept"}`; start != want {
					t.Fatalf("start_json after the upgrade = %s, want %s", start, want)
				}
			}
			if tc.v == 4 {
				var (
					request, digest, challenge, episode string
					created, expires, unconsumed        int64
				)
				if err := db.QueryRow(`SELECT request, digest, challenge, episode_id, created_at, expires_at, consumed_at IS NULL
					FROM mint_requests WHERE request_id = ?`, reqID).Scan(&request, &digest, &challenge, &episode, &created, &expires, &unconsumed); err != nil {
					t.Fatal(err)
				}
				if request != "{}" || digest != strings.Repeat("ab", 32) || challenge != strings.Repeat("cd", 32) ||
					episode != "ep" || created != 1000 || expires != 1600 || unconsumed != 1 {
					t.Fatalf("mint request row after the upgrade: request %q digest %q challenge %q episode %q created %d expires %d unconsumed %d",
						request, digest, challenge, episode, created, expires, unconsumed)
				}
			}
		})
	}
}

// AS6: the older-store refusal names its Feature and Need, and the text is byte-identical at every
// mint-request site (a store's version is the only thing that varies).
func TestOlderCredentialsStoreErrorText(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV3Credentials(t, path, nil)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	want := fmt.Sprintf("store: credentials store %q is schema 3; mint requests need schema 4: restart "+
		"`ailang-worldd serve` with this build once (it upgrades the store at startup), then run this again", path)
	ctx := testCtx(t)
	for name, call := range map[string]func() error{
		"PutMintRequest": func() error {
			return c.PutMintRequest(ctx, testMintRequestRow(1, 1000), 1000)
		},
		"MintRequest": func() error {
			_, _, err := c.MintRequest(ctx, "mr-0000000000000001")
			return err
		},
		"PendingMintRequests": func() error {
			_, err := c.PendingMintRequests(ctx, 1000, 0)
			return err
		},
		"MintSessionForRequest": func() error {
			return c.MintSessionForRequest(ctx, consumeSession(1), okConsume(testMintRequestRow(1, 1000), 1010, new(int), new(int)))
		},
	} {
		err := call()
		if err == nil {
			t.Fatalf("%s on a schema 3 store: no error", name)
		}
		if err.Error() != want {
			t.Fatalf("%s error text:\n got: %s\nwant: %s", name, err.Error(), want)
		}
		var oe *OlderCredentialsStoreError
		if !errors.As(err, &oe) {
			t.Fatalf("%s error %T is not *OlderCredentialsStoreError", name, err)
		}
		if oe.Need != 4 || oe.Feature != "mint requests" || oe.Version != 3 {
			t.Fatalf("%s: Need %d Feature %q Version %d; want 4, mint requests, 3", name, oe.Need, oe.Feature, oe.Version)
		}
	}

	// The approvers wording, constructed directly (the approver methods arrive with M2b).
	got := (&OlderCredentialsStoreError{Path: "p", Version: 4, Need: 5, Feature: "approvers"}).Error()
	const approversWant = "store: credentials store \"p\" is schema 4; approvers need schema 5: restart " +
		"`ailang-worldd serve` with this build once (it upgrades the store at startup), then run this again"
	if got != approversWant {
		t.Fatalf("approvers refusal text:\n got: %s\nwant: %s", got, approversWant)
	}
}
