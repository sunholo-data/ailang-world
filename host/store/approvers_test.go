package store

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
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

// AS2: every EnrolApprover refusal returns its typed error and writes nothing.
func TestEnrolApproverRefusals(t *testing.T) {
	ctx := testCtx(t)

	// v4: the refusal names approvers, schema 5, and the restart; the approvers table is not created.
	v4 := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV4Credentials(t, v4, nil)
	c4, err := OpenCredentials(v4)
	if err != nil {
		t.Fatal(err)
	}
	err = c4.EnrolApprover(ctx, testApprover("mark-mac", 1000))
	var older *OlderCredentialsStoreError
	if !errors.As(err, &older) {
		t.Fatalf("v4 enrol: err %v, want *OlderCredentialsStoreError", err)
	}
	if older.Version != 4 || older.Need != 5 || older.Feature != "approvers" {
		t.Fatalf("v4 enrol: Version %d Need %d Feature %q; want 4, 5, approvers", older.Version, older.Need, older.Feature)
	}
	if !strings.Contains(err.Error(), "approvers need schema 5") || !strings.Contains(err.Error(), "restart `ailang-worldd serve`") {
		t.Fatalf("v4 enrol text: %s", err.Error())
	}
	_ = c4.Close()
	if got := userVersion(t, v4); got != 4 {
		t.Fatalf("user_version after the refused enrol = %d, want 4", got)
	}
	var tables int
	if err := rawDB(t, v4).QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'approvers'`).Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if tables != 0 {
		t.Fatalf("a refused enrol created the approvers table (%d objects named approvers)", tables)
	}

	// Bad rows on a fresh (version 5) store: each refusal is ErrBadApprover and writes nothing.
	c := openTestCredentials(t)
	bad := []struct {
		name string
		row  ApproverRow
	}{
		{"empty name", ApproverRow{Name: "", Format: "ssh", PublicKey: []byte("key-"), EnrolledAt: 1000}},
		{"format gpg", ApproverRow{Name: "bad", Format: "gpg", PublicKey: []byte("key-bad"), EnrolledAt: 1000}},
		{"nil key", ApproverRow{Name: "bad", Format: "ssh", PublicKey: nil, EnrolledAt: 1000}},
		{"enrolled 0", ApproverRow{Name: "bad", Format: "ssh", PublicKey: []byte("key-bad"), EnrolledAt: 0}},
		{"enrolled -1", ApproverRow{Name: "bad", Format: "ssh", PublicKey: []byte("key-bad"), EnrolledAt: -1}},
		{"last used set", ApproverRow{Name: "bad", Format: "ssh", PublicKey: []byte("key-bad"), EnrolledAt: 1000, LastUsedAt: 1500}},
		{"revoked set", ApproverRow{Name: "bad", Format: "ssh", PublicKey: []byte("key-bad"), EnrolledAt: 1000, RevokedAt: 1500}},
	}
	for _, tc := range bad {
		before := approverCount(t, c)
		if err := c.EnrolApprover(ctx, tc.row); !errors.Is(err, ErrBadApprover) {
			t.Fatalf("%s: err %v, want ErrBadApprover", tc.name, err)
		}
		if after := approverCount(t, c); after != before {
			t.Fatalf("%s: rows %d -> %d, want no write", tc.name, before, after)
		}
	}

	// Good: one row.
	if err := c.EnrolApprover(ctx, testApprover("mark-mac", 1000)); err != nil {
		t.Fatalf("good enrol: %v", err)
	}
	if got := approverCount(t, c); got != 1 {
		t.Fatalf("after the good enrol: %d rows, want 1", got)
	}

	// Name taken, and key enrolled: typed, one row still.
	nameTaken := ApproverRow{Name: "mark-mac", Format: "ssh", PublicKey: []byte("key-other"), EnrolledAt: 1100}
	keyTaken := ApproverRow{Name: "other", Format: "ssh", PublicKey: []byte("key-mark-mac"), EnrolledAt: 1100}
	checkRefusals := func(when string) {
		t.Helper()
		if err := c.EnrolApprover(ctx, nameTaken); !errors.Is(err, ErrApproverNameTaken) {
			t.Fatalf("%s: name taken: err %v, want ErrApproverNameTaken", when, err)
		}
		err := c.EnrolApprover(ctx, keyTaken)
		if !errors.Is(err, ErrApproverKeyEnrolled) {
			t.Fatalf("%s: key enrolled: err %v, want ErrApproverKeyEnrolled", when, err)
		}
		if !strings.Contains(err.Error(), `"mark-mac"`) {
			t.Fatalf("%s: key enrolled text %q does not name the holder", when, err.Error())
		}
		if got := approverCount(t, c); got != 1 {
			t.Fatalf("%s: %d rows after the refusals, want 1", when, got)
		}
	}
	checkRefusals("active")

	// After a revoke the same two refusals still hold: names and keys are never reused.
	if _, err := c.db.Exec(`UPDATE approvers SET revoked_at = 1200 WHERE name = 'mark-mac'`); err != nil {
		t.Fatal(err)
	}
	checkRefusals("revoked")

	// Control: a new name and key enrol.
	if err := c.EnrolApprover(ctx, testApprover("third", 1300)); err != nil {
		t.Fatalf("control enrol: %v", err)
	}
	if got := approverCount(t, c); got != 2 {
		t.Fatalf("after the control enrol: %d rows, want 2", got)
	}
}

// AS3: the two lookups return the enrolled row exactly, and miss without an error.
func TestApproverLookups(t *testing.T) {
	ctx := testCtx(t)
	c := openTestCredentials(t)
	if err := c.EnrolApprover(ctx, testApprover("a", 1000)); err != nil {
		t.Fatal(err)
	}
	if err := c.EnrolApprover(ctx, testApprover("b", 1100)); err != nil {
		t.Fatal(err)
	}

	row, ok, err := c.Approver(ctx, "a")
	if err != nil || !ok {
		t.Fatalf("Approver(a): ok %v err %v", ok, err)
	}
	if want := (ApproverRow{Name: "a", Format: "ssh", PublicKey: []byte("key-a"), EnrolledAt: 1000}); !reflect.DeepEqual(row, want) {
		t.Fatalf("Approver(a) = %+v, want %+v", row, want)
	}
	byKey, ok, err := c.ApproverByKey(ctx, []byte("key-b"))
	if err != nil || !ok || byKey.Name != "b" || byKey.EnrolledAt != 1100 {
		t.Fatalf("ApproverByKey(key-b) = %+v ok %v err %v; want approver b", byKey, ok, err)
	}
	if row, ok, err := c.Approver(ctx, "zzz"); err != nil || ok || !reflect.DeepEqual(row, ApproverRow{}) {
		t.Fatalf("Approver(zzz) = %+v ok %v err %v; want zero, false, nil", row, ok, err)
	}
	if _, ok, err := c.ApproverByKey(ctx, []byte("key-")); err != nil || ok {
		t.Fatalf("ApproverByKey(prefix of both keys): ok %v err %v; want false, nil", ok, err)
	}
	if _, ok, err := c.ApproverByKey(ctx, []byte("key-a\x00")); err != nil || ok {
		t.Fatalf("ApproverByKey(key-a + NUL): ok %v err %v; want false, nil", ok, err)
	}

	// A version 4 store refuses both lookups with the approvers refusal.
	v4 := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV4Credentials(t, v4, nil)
	c4, err := OpenCredentials(v4)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c4.Close() }()
	if _, _, err := c4.Approver(ctx, "a"); !isOlderApprovers(err) {
		t.Fatalf("v4 Approver: err %v, want the approvers refusal", err)
	}
	if _, _, err := c4.ApproverByKey(ctx, []byte("key-a")); !isOlderApprovers(err) {
		t.Fatalf("v4 ApproverByKey: err %v, want the approvers refusal", err)
	}
}

// isOlderApprovers reports whether err is the approvers refusal (Need 5).
func isOlderApprovers(err error) bool {
	var oe *OlderCredentialsStoreError
	return errors.As(err, &oe) && oe.Need == 5 && oe.Feature == "approvers"
}

// approverRowColumns maps every ApproverRow field to its approvers column.
var approverRowColumns = map[string]string{
	"Name": "name", "Format": "format", "PublicKey": "public_key",
	"EnrolledAt": "enrolled_at", "LastUsedAt": "last_used_at", "RevokedAt": "revoked_at",
}

// RC2: the struct, the table and the scanned rows agree field for field (a swapped or dropped
// column is caught by value, not only by name).
func TestApproverRowMatchesTheTable(t *testing.T) {
	ctx := testCtx(t)
	c := openTestCredentials(t)

	// (1) The struct's field names are the map's keys.
	typ := reflect.TypeOf(ApproverRow{})
	var fields []string
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	var keys []string
	for k := range approverRowColumns {
		keys = append(keys, k)
	}
	sort.Strings(fields)
	sort.Strings(keys)
	if !reflect.DeepEqual(fields, keys) {
		t.Fatalf("ApproverRow fields %v, map keys %v", fields, keys)
	}

	// (2) The table's column names are the map's values.
	rows, err := c.db.Query(`PRAGMA table_info(approvers)`)
	if err != nil {
		t.Fatal(err)
	}
	var cols []string
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatal(err)
		}
		cols = append(cols, name)
	}
	_ = rows.Close()
	var want []string
	for _, v := range approverRowColumns {
		want = append(want, v)
	}
	sort.Strings(cols)
	sort.Strings(want)
	if !reflect.DeepEqual(cols, want) {
		t.Fatalf("approvers columns %v, want %v", cols, want)
	}

	// (3) Every field of a fully populated row reads back, field by field.
	if err := c.EnrolApprover(ctx, ApproverRow{Name: "rc2", Format: "ssh", PublicKey: []byte("rc2-key"), EnrolledAt: 1000}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.db.Exec(`UPDATE approvers SET last_used_at = 1500, revoked_at = 1700 WHERE name = 'rc2'`); err != nil {
		t.Fatal(err)
	}
	wantRow := ApproverRow{Name: "rc2", Format: "ssh", PublicKey: []byte("rc2-key"), EnrolledAt: 1000, LastUsedAt: 1500, RevokedAt: 1700}
	byName, ok, err := c.Approver(ctx, "rc2")
	if err != nil || !ok {
		t.Fatalf("Approver(rc2): ok %v err %v", ok, err)
	}
	byKey, ok, err := c.ApproverByKey(ctx, []byte("rc2-key"))
	if err != nil || !ok {
		t.Fatalf("ApproverByKey(rc2-key): ok %v err %v", ok, err)
	}
	for _, got := range []ApproverRow{byName, byKey} {
		gv, wv := reflect.ValueOf(got), reflect.ValueOf(wantRow)
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			if gv.Field(i).IsZero() {
				t.Fatalf("field %s reads as zero", name)
			}
			if !reflect.DeepEqual(gv.Field(i).Interface(), wv.Field(i).Interface()) {
				t.Fatalf("field %s = %v, want %v", name, gv.Field(i).Interface(), wv.Field(i).Interface())
			}
		}
	}
}
