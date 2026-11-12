package store

// Row 201 A3 sprint 1 M1 (design_docs/planned/w-world-approvals-by-key-a3.md §2.2, §6 SA1-SA3):
// credentials schema 7's surface_approvals. Every store is under t.TempDir() (D-WORLD-109).

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// canonicalSurfaceApprovalsDDL is the current (version 7, row 201 A3) surface_approvals table.
var canonicalSurfaceApprovalsDDL = `CREATE TABLE surface_approvals (
    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
    approver    TEXT NOT NULL CHECK (approver <> ''),
    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
    digest      TEXT NOT NULL CHECK (length(digest) = 64),
    received_at INTEGER NOT NULL
)`

// makeV6Credentials writes a schema-6 credentials store at path: makeV5Credentials, the
// webauthn_credentials table and user_version 6.
func makeV6Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	makeV5Credentials(t, path, rows)
	db := rawDB(t, path)
	if _, err := db.Exec(passkeysSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 6)
}

// testSurfaceApproval is a well-formed row for testMintRequestRow(n, ...).
func testSurfaceApproval(n int) SurfaceApprovalRow {
	r := testMintRequestRow(n, 1000)
	return SurfaceApprovalRow{RequestID: r.RequestID, Format: "webauthn", Approver: "webauthn:mark", Artifact: []byte(`{"id":"x"}`), Digest: r.Digest}
}

// surfaceApprovalCount counts c's surface_approvals rows.
func surfaceApprovalCount(t *testing.T, c *Credentials) int {
	t.Helper()
	var n int
	if err := c.db.QueryRow(`SELECT count(*) FROM surface_approvals`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// SA1: a fresh store is schema 7 with surface_approvals as credentials_schema.sql and the upgrade
// write it; a schema-6 store is upgraded once by UpgradeCredentials only, its passkey kept and the
// new table empty; a session open of the v6 store neither upgrades it nor records an approval,
// which names the serve restart.
func TestCredentialsSchema7SurfaceApprovals(t *testing.T) {
	c := openTestCredentials(t)
	ddl := tableDDL(t, c.db)
	for name, got := range map[string]string{"schema file": ddl["surface_approvals"], "surfaceApprovalsSQL": surfaceApprovalsSQL} {
		if normalizeDDL(got) != normalizeDDL(canonicalSurfaceApprovalsDDL) {
			t.Fatalf("surface_approvals DDL (%s):\n got: %s\nwant: %s", name, got, canonicalSurfaceApprovalsDDL)
		}
	}
	if got, want := credentialsTables(7), []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("credentialsTables(7) = %v, want %v", got, want)
	}
	if got, want := credentialsTables(6), []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("credentialsTables(6) = %v, want %v", got, want)
	}

	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
	a, p := testPasskey("kept", 3)
	raw := rawDB(t, path)
	for _, ins := range []struct {
		q    string
		args []any
	}{
		{`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`, []any{a.Name, a.Format, a.PublicKey, a.EnrolledAt}},
		{`INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count) VALUES (?, ?, ?, ?, ?, ?, 0, 1)`,
			[]any{p.Name, p.CredentialID, p.RPID, p.Origin, p.Alg, p.AAGUID}},
		{`INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, x'7b7d', ?, ?, 'ep', 1000, 1600)`,
			[]any{testSurfaceApproval(1).RequestID, testSurfaceApproval(1).Digest, strings.Repeat("cd", 32)}},
	} {
		if _, err := raw.Exec(ins.q, ins.args...); err != nil {
			t.Fatal(err)
		}
	}
	v6, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	err = v6.RecordSurfaceApproval(testCtx(t), testSurfaceApproval(1), 1010)
	var oe *OlderCredentialsStoreError
	if !errors.As(err, &oe) || oe.Version != 6 || oe.Need != 7 || oe.Feature != "surface approvals" ||
		!strings.Contains(err.Error(), "surface approvals need schema 7: restart `ailang-worldd serve`") {
		t.Fatalf("v6 record: %v, want the restart refusal naming schema 7", err)
	}
	if _, _, err := v6.SurfaceApproval(testCtx(t), testSurfaceApproval(1).RequestID); !errors.As(err, &oe) || oe.Need != 7 {
		t.Fatalf("v6 read: %v, want OlderCredentialsStoreError need 7", err)
	}
	_ = v6.Close()
	if got := userVersion(t, path); got != 6 {
		t.Fatalf("a session open changed user_version to %d, want 6", got)
	}
	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded {
		t.Fatalf("UpgradeCredentials: upgraded %v, err %v; want true, nil", upgraded, err)
	}
	if again, err := UpgradeCredentials(path); err != nil || again {
		t.Fatalf("second UpgradeCredentials: upgraded %v, err %v; want false, nil", again, err)
	}
	if got := userVersion(t, path); got != credentialsSchemaVersion {
		t.Fatalf("user_version after the upgrade = %d, want %d", got, credentialsSchemaVersion)
	}
	up := tableDDL(t, rawDB(t, path))
	requireExactTableNames(t, "upgraded v6 store", up, credentialsTables(credentialsSchemaVersion))
	if got, want := normalizeDDL(up["surface_approvals"]), normalizeDDL(canonicalSurfaceApprovalsDDL); got != want {
		t.Fatalf("upgraded surface_approvals DDL:\n got: %s\nwant: %s", got, want)
	}
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	if _, gotP, ok, err := ro.PasskeyByCredentialID(testCtx(t), p.CredentialID); err != nil || !ok || gotP.SignCount != 1 {
		t.Fatalf("the passkey after the upgrade: ok %v err %v %+v", ok, err, gotP)
	}
	if _, ok, err := ro.SurfaceApproval(testCtx(t), testSurfaceApproval(1).RequestID); err != nil || ok {
		t.Fatalf("a recorded approval right after the upgrade: ok %v err %v; want none", ok, err)
	}
}

// SA2: each guard of RecordSurfaceApproval refuses with its sentinel and writes nothing; a valid row
// is read back with received_at = now; a later valid row replaces it while the request is pending.
func TestRecordSurfaceApprovalGuards(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	pending := testMintRequestRow(1, 1000)
	if err := c.PutMintRequest(ctx, pending, 1000); err != nil {
		t.Fatal(err)
	}
	used := testMintRequestRow(2, 1000)
	if err := c.PutMintRequest(ctx, used, 1000); err != nil {
		t.Fatal(err)
	}
	if err := c.MintSessionForRequest(ctx, consumeSession(1), okConsume(used, 1005, new(int), new(int))); err != nil {
		t.Fatal(err)
	}
	edit := func(f func(*SurfaceApprovalRow)) SurfaceApprovalRow {
		r := testSurfaceApproval(1)
		f(&r)
		return r
	}
	for _, arm := range []struct {
		name string
		row  SurfaceApprovalRow
		now  int64
		want error
	}{
		{"S2-id", edit(func(r *SurfaceApprovalRow) { r.RequestID = "mr-00000000000000A1" }), 1010, ErrBadSurfaceApproval},
		{"S2-format", edit(func(r *SurfaceApprovalRow) { r.Format, r.Approver = "tty", "tty:mark" }), 1010, ErrBadSurfaceApproval},
		{"S2-approver-other-format", edit(func(r *SurfaceApprovalRow) { r.Approver = "ssh:mark" }), 1010, ErrBadSurfaceApproval},
		{"S2-approver-no-name", edit(func(r *SurfaceApprovalRow) { r.Approver = "webauthn:" }), 1010, ErrBadSurfaceApproval},
		{"S2-artifact-empty", edit(func(r *SurfaceApprovalRow) { r.Artifact = nil }), 1010, ErrBadSurfaceApproval},
		{"S2-artifact-over", edit(func(r *SurfaceApprovalRow) { r.Artifact = bytes.Repeat([]byte("a"), 16385) }), 1010, ErrBadSurfaceApproval},
		{"S2-digest", edit(func(r *SurfaceApprovalRow) { r.Digest = strings.Repeat("AB", 32) }), 1010, ErrBadSurfaceApproval},
		{"S3-unknown", edit(func(r *SurfaceApprovalRow) { r.RequestID = testMintRequestRow(9, 1000).RequestID }), 1010, ErrRequestUnknown},
		{"S4-consumed", edit(func(r *SurfaceApprovalRow) { r.RequestID = used.RequestID }), 1010, ErrRequestConsumed},
		{"S4-clock-behind", testSurfaceApproval(1), 999, ErrRequestClockBehind},
		{"S4-expired", testSurfaceApproval(1), 1600, ErrRequestExpired},
		{"S4-digest-changed", edit(func(r *SurfaceApprovalRow) { r.Digest = strings.Repeat("ef", 32) }), 1010, ErrRequestDigestChanged},
	} {
		t.Run(arm.name, func(t *testing.T) {
			if err := c.RecordSurfaceApproval(ctx, arm.row, arm.now); !errors.Is(err, arm.want) {
				t.Fatalf("RecordSurfaceApproval = %v, want %v", err, arm.want)
			}
			if n := surfaceApprovalCount(t, c); n != 0 {
				t.Fatalf("%d rows after a refusal, want 0", n)
			}
		})
	}
	// Controls: the bounds themselves are admitted (an artifact of exactly 16384 bytes, the
	// request's creation second, the last second before expiry).
	for _, ok := range []struct {
		row SurfaceApprovalRow
		now int64
	}{
		{edit(func(r *SurfaceApprovalRow) { r.Artifact = bytes.Repeat([]byte("a"), 16384) }), 1000},
		{edit(func(r *SurfaceApprovalRow) { r.Artifact = []byte("b") }), 1599},
	} {
		if err := c.RecordSurfaceApproval(ctx, ok.row, ok.now); err != nil {
			t.Fatalf("control at %d: %v", ok.now, err)
		}
	}
	want := SurfaceApprovalRow{RequestID: pending.RequestID, Format: "ssh", Approver: "ssh:mark-mac", Artifact: []byte("-----BEGIN SSH SIGNATURE-----\n"), Digest: pending.Digest, ReceivedAt: 1010}
	in := want
	in.ReceivedAt = 0
	if err := c.RecordSurfaceApproval(ctx, in, 1010); err != nil {
		t.Fatal(err)
	}
	got, ok, err := c.SurfaceApproval(ctx, pending.RequestID)
	if err != nil || !ok || !reflect.DeepEqual(got, want) || surfaceApprovalCount(t, c) != 1 {
		t.Fatalf("after S5's replace: %+v ok %v err %v (%d rows); want %+v alone", got, ok, err, surfaceApprovalCount(t, c), want)
	}
	if _, ok, err := c.SurfaceApproval(ctx, used.RequestID); err != nil || ok {
		t.Fatalf("the consumed request has a recorded approval: ok %v err %v", ok, err)
	}
}

// SA3: PutMintRequest's prune deletes the recorded approval of a request it prunes, in the same
// transaction, and keeps a consumed request's; on a schema-6 store it prunes and puts as before.
func TestSurfaceApprovalPrunedWithItsRequest(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	gone, kept := testMintRequestRow(1, 1000), testMintRequestRow(2, 1000)
	for _, r := range []MintRequestRow{gone, kept} {
		if err := c.PutMintRequest(ctx, r, 1000); err != nil {
			t.Fatal(err)
		}
		row := testSurfaceApproval(1)
		row.RequestID = r.RequestID
		if err := c.RecordSurfaceApproval(ctx, row, 1010); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.MintSessionForRequest(ctx, consumeSession(1), okConsume(kept, 1020, new(int), new(int))); err != nil {
		t.Fatal(err)
	}
	if err := c.PutMintRequest(ctx, testMintRequestRow(3, 1600+86400), 1600+86400); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := c.SurfaceApproval(ctx, gone.RequestID); err != nil || ok {
		t.Fatalf("the pruned request's approval: ok %v err %v; want gone", ok, err)
	}
	if _, ok, err := c.SurfaceApproval(ctx, kept.RequestID); err != nil || !ok {
		t.Fatalf("the consumed request's approval: ok %v err %v; want kept", ok, err)
	}
	if n := surfaceApprovalCount(t, c); n != 1 {
		t.Fatalf("%d rows after the prune, want 1", n)
	}

	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV6Credentials(t, path, nil)
	v6, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = v6.Close() }()
	for _, r := range []MintRequestRow{gone, testMintRequestRow(3, 1600+86400)} {
		if err := v6.PutMintRequest(ctx, r, r.CreatedAt); err != nil {
			t.Fatalf("v6 put %s: %v", r.RequestID, err)
		}
	}
	if _, ok, err := v6.MintRequest(ctx, gone.RequestID); err != nil || ok {
		t.Fatalf("v6: the expired request was not pruned: ok %v err %v", ok, err)
	}
}
