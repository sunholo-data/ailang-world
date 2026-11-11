package main

// Row 201 A1b sprint 2: the approver verb tests' rig. Every store is under t.TempDir(); every
// key is a checked-in fixture (D-WORLD-109).

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// The checked-in key fixtures' fingerprints (host/approval/testdata/sshsig/SOURCE.txt).
const (
	edFingerprint   = "SHA256:pHMYJm+CndqyZGrzMHQQBH0tNeNGj1U9sy0JrVsYdSQ"
	p256Fingerprint = "SHA256:yMwGgMZdq1aVNYRG/HO2T69s/Fb14yPZW0GBWtwRyEo"
)

// fixturePub is the path of a checked-in public key, host/approval/testdata/sshsig/<name>.
func fixturePub(name string) string {
	return filepath.Join("..", "..", "host", "approval", "testdata", "sshsig", name)
}

// fixtureKeyOf parses a checked-in public key with the production parser.
func fixtureKeyOf(t *testing.T, name string) approval.SSHPublicKey {
	t.Helper()
	b, err := os.ReadFile(fixturePub(name))
	if err != nil {
		t.Fatal(err)
	}
	k, err := approval.ParseAuthorizedKey(b)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// approverStore is an existing world store and a fresh schema-5 credentials store beside it.
func approverStore(t *testing.T) (db, creds string) {
	t.Helper()
	db = dbPathOf(t)
	creds, err := credentialsFor(db, "")
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return db, creds
}

// approverStoreV4 is approverStore taken back to schema 4: the webauthn_credentials and approvers
// tables dropped and user_version 4, which is exactly the version-4 table set.
func approverStoreV4(t *testing.T) (db, creds string) {
	t.Helper()
	db, creds = approverStore(t)
	raw := rawCredentials(t, creds)
	defer func() { _ = raw.Close() }()
	for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
		if _, err := raw.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	return db, creds
}

// approverRows lists creds' approvers read-only.
func approverRows(t *testing.T, creds string) []store.ApproverRow {
	t.Helper()
	ro, err := store.OpenCredentialsReadOnly(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := ro.ListApprovers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// seedApprover enrols name with key at `at` straight through the store (no verb, no terminal).
func seedApprover(t *testing.T, creds, name string, key approval.SSHPublicKey, at int64) {
	t.Helper()
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.EnrolApprover(ctx, store.ApproverRow{Name: name, Format: approval.FormatSSH, PublicKey: key.Blob, EnrolledAt: at}); err != nil {
		t.Fatal(err)
	}
}

// revokeSeeded revokes name at `at` straight through the store.
func revokeSeeded(t *testing.T, creds, name string, at int64) {
	t.Helper()
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.RevokeApprover(ctx, name, at); err != nil {
		t.Fatal(err)
	}
}

// enrollArgs is `approver enroll --name name --format ssh --pubkey pub --db db`.
func enrollArgs(db, name, pub string) []string {
	return []string{"enroll", "--name", name, "--format", "ssh", "--pubkey", pub, "--db", db}
}
