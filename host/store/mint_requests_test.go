package store

import (
	"context"
	"database/sql"
	"errors"
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

// testMintRequestRow is a valid unconsumed row: id "mr-" + fmt.Sprintf("%016x", n), Request
// []byte(`{"n":<n>}`), Digest strings.Repeat("ab", 32), Challenge strings.Repeat("cd", 32),
// EpisodeID "ep", CreatedAt created, ExpiresAt created+600.
func testMintRequestRow(n int, created int64) MintRequestRow {
	return MintRequestRow{
		RequestID: fmt.Sprintf("mr-%016x", n),
		Request:   []byte(fmt.Sprintf(`{"n":%d}`, n)),
		Digest:    strings.Repeat("ab", 32),
		Challenge: strings.Repeat("cd", 32),
		EpisodeID: "ep",
		CreatedAt: created,
		ExpiresAt: created + 600,
	}
}

// testCtx is a bounded context for store calls in these tests.
func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	return ctx
}

// TestPutMintRequestNeedsVersion4 (ST3): on a version-3 store the three request methods refuse
// with OlderCredentialsStoreError and leave user_version at 3.
func TestPutMintRequestNeedsVersion4(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV3Credentials(t, path, nil)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	ctx := testCtx(t)
	checkOlder := func(name string, err error) {
		t.Helper()
		var older *OlderCredentialsStoreError
		if !errors.As(err, &older) {
			t.Fatalf("%s on a version-3 store: err %v, want *OlderCredentialsStoreError", name, err)
		}
		if older.Version != 3 {
			t.Fatalf("%s: Version = %d, want 3", name, older.Version)
		}
		if !strings.Contains(err.Error(), "restart `ailang-worldd serve`") || !strings.Contains(err.Error(), "need schema 4") {
			t.Fatalf("%s: message %q lacks the restart hint or \"need schema 4\"", name, err)
		}
	}
	checkOlder("PutMintRequest", c.PutMintRequest(ctx, testMintRequestRow(1, 1000), 1000))
	_, _, err = c.MintRequest(ctx, "mr-0000000000000001")
	checkOlder("MintRequest", err)
	_, err = c.PendingMintRequests(ctx, 1000, 0)
	checkOlder("PendingMintRequests", err)
	if got := userVersion(t, path); got != 3 {
		t.Fatalf("user_version = %d after refused calls, want 3", got)
	}
}

// TestPutMintRequestCapsAndPrunes (ST7): the cap counts only live pending rows, the day-old expired
// unconsumed rows are pruned, consumed rows are kept, and the pending list is ordered and filtered.
func TestPutMintRequestCapsAndPrunes(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	mustPut := func(row MintRequestRow, now int64) {
		t.Helper()
		if err := c.PutMintRequest(ctx, row, now); err != nil {
			t.Fatalf("put %s at %d: %v", row.RequestID, now, err)
		}
	}
	a, b, cc := testMintRequestRow(1, 1000), testMintRequestRow(2, 999), testMintRequestRow(3, 1001)
	mustPut(a, 1000)
	mustPut(b, 1000)
	mustPut(cc, 1000)
	// Mark B consumed directly: the real consume is covered by ST4.
	if _, err := c.db.Exec(`UPDATE mint_requests SET consumed_at = 1100, approved_by = 'tty', credential_id = 'x' WHERE request_id = ?`, b.RequestID); err != nil {
		t.Fatal(err)
	}

	// Step 2: at 88000 the boundary 1600 = 1600 + 86400 - 86400 prunes A only.
	const later = 88000
	mustPut(testMintRequestRow(4, later), later)
	if _, ok, err := c.MintRequest(ctx, a.RequestID); err != nil || ok {
		t.Fatalf("A (expires 1600, now 88000): found %v, err %v; want pruned", ok, err)
	}
	if row, ok, err := c.MintRequest(ctx, b.RequestID); err != nil || !ok || !row.Consumed {
		t.Fatalf("B (consumed, expired): found %v consumed %v, err %v; want kept", ok, row.Consumed, err)
	}
	if _, ok, err := c.MintRequest(ctx, cc.RequestID); err != nil || !ok {
		t.Fatalf("C (expires 1601): found %v, err %v; want kept", ok, err)
	}
	if _, ok, err := c.MintRequest(ctx, testMintRequestRow(4, later).RequestID); err != nil || !ok {
		t.Fatalf("D: found %v, err %v; want found", ok, err)
	}

	// Step 3: D plus E1..E15 is 16 live rows (C is expired, so it does not count).
	for n := 10; n <= 24; n++ {
		mustPut(testMintRequestRow(n, later), later)
	}
	over := testMintRequestRow(30, later)
	if err := c.PutMintRequest(ctx, over, later); !errors.Is(err, ErrTooManyPendingRequests) {
		t.Fatalf("17th live request: err %v, want ErrTooManyPendingRequests", err)
	}
	if _, ok, err := c.MintRequest(ctx, over.RequestID); err != nil || ok {
		t.Fatalf("refused request: found %v, err %v; want not stored", ok, err)
	}

	// Step 4: the pending list is the 16 live rows, newest first, ties by id descending.
	pending, err := c.PendingMintRequests(ctx, later, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != MaxPendingMintRequests {
		t.Fatalf("PendingMintRequests returned %d rows, want %d", len(pending), MaxPendingMintRequests)
	}
	for _, row := range pending {
		if row.RequestID == a.RequestID || row.RequestID == b.RequestID || row.RequestID == cc.RequestID {
			t.Fatalf("PendingMintRequests lists %s, which is expired or consumed", row.RequestID)
		}
	}
	if got, want := pending[0].RequestID, testMintRequestRow(24, later).RequestID; got != want {
		t.Fatalf("first pending id = %s, want %s (ties by request_id descending)", got, want)
	}

	// Step 4b: the filters on a store with few rows, so no LIMIT hides them. At 1500: X is consumed
	// (created latest), Z is expired (1100), Y is live. Only Y is pending.
	few := openTestCredentials(t)
	x, y, z := testMintRequestRow(5, 1200), testMintRequestRow(6, 1000), testMintRequestRow(7, 500)
	for _, row := range []MintRequestRow{x, y, z} {
		if err := few.PutMintRequest(ctx, row, 1000); err != nil {
			t.Fatalf("put %s: %v", row.RequestID, err)
		}
	}
	if _, err := few.db.Exec(`UPDATE mint_requests SET consumed_at = 1300, approved_by = 'tty', credential_id = 'x' WHERE request_id = ?`, x.RequestID); err != nil {
		t.Fatal(err)
	}
	small, err := few.PendingMintRequests(ctx, 1500, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(small) != 1 || small[0].RequestID != y.RequestID {
		ids := make([]string, 0, len(small))
		for _, row := range small {
			ids = append(ids, row.RequestID)
		}
		t.Fatalf("pending at 1500 = %v, want only %s", ids, y.RequestID)
	}
}

// TestPutMintRequestRefusesABadRow (ST8): every bad arm alone is refused with ErrBadMintRequest and
// stores nothing. The weak table CHECK accepts the bad id (plan §1.4 F2), so the Go check is load-bearing.
func TestPutMintRequestRefusesABadRow(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	good := testMintRequestRow(1, 1000)
	arms := map[string]func(*MintRequestRow){
		"id zzzz":           func(r *MintRequestRow) { r.RequestID = "mr-0zzzzzzzzzzzzzzz" },
		"id short":          func(r *MintRequestRow) { r.RequestID = "mr-0123" },
		"nil request":       func(r *MintRequestRow) { r.Request = nil },
		"upper digest":      func(r *MintRequestRow) { r.Digest = strings.Repeat("AB", 32) },
		"63-char challenge": func(r *MintRequestRow) { r.Challenge = strings.Repeat("c", 63) },
		"empty episode":     func(r *MintRequestRow) { r.EpisodeID = "" },
		"expires = created": func(r *MintRequestRow) { r.ExpiresAt = r.CreatedAt },
		"consumed":          func(r *MintRequestRow) { r.Consumed = true },
	}
	for name, mutate := range arms {
		row := good
		mutate(&row)
		err := c.PutMintRequest(ctx, row, 1000)
		if !errors.Is(err, ErrBadMintRequest) {
			t.Fatalf("%s: err %v, want ErrBadMintRequest", name, err)
		}
		if _, ok, err := c.MintRequest(ctx, row.RequestID); err != nil || ok {
			t.Fatalf("%s: stored %v, err %v; want nothing stored", name, ok, err)
		}
	}
	if err := c.PutMintRequest(ctx, good, 1000); err != nil {
		t.Fatalf("clean row refused: %v", err)
	}

	// Instrument control: the table's own CHECK accepts the bad id, so only the Go check refuses it.
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c2, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c2.Close() }()
	if _, err := c2.db.Exec(`INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES ('mr-0zzzzzzzzzzzzzzz', x'00', ?, ?, 'ep', 1, 2)`,
		strings.Repeat("a", 64), strings.Repeat("a", 64)); err != nil {
		t.Fatalf("the weak CHECK refused the bad id (control broken): %v", err)
	}
}
