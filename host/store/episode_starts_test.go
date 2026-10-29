package store

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

// Version 2 already has bound credentials. Only serve upgrades it, preserving
// every byte while adding independent attended start history.
func TestEpisodeStartsUpgradeV2AndAtomicMint(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	db := rawDB(t, path)
	if _, err := db.Exec(canonicalCredentialsDDL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 2)
	want := SessionRow{CredentialID: "live", EpisodeID: "old", GrantsJSON: "[]", ExpiresAt: 9999999999, CreatedAt: 1, Project: "p", ProfileDigest: "sha256:old"}
	if _, err := db.Exec("INSERT INTO session_credentials VALUES (?, ?, ?, ?, ?, ?, ?)", want.CredentialID, want.EpisodeID, want.GrantsJSON, want.ExpiresAt, want.CreatedAt, want.Project, want.ProfileDigest); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.Version() != 2 {
		t.Fatal("CLI upgraded older store")
	}
	start := &EpisodeStart{Schema: "world/episode-start/v1", Episode: "new"}
	row := SessionRow{CredentialID: "new", EpisodeID: "new", GrantsJSON: "[]", ExpiresAt: 20, CreatedAt: 10}
	var refused *UnboundCredentialsStoreError
	if err := c.MintSessionWithStart(boundedTestContext(t), row, start); !errors.As(err, &refused) {
		t.Fatalf("old store accepted start: %v", err)
	}
	_ = c.Close()
	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded {
		t.Fatalf("upgrade %v %v", upgraded, err)
	}
	if upgraded, err := UpgradeCredentials(path); err != nil || upgraded {
		t.Fatalf("repeat upgrade %v %v", upgraded, err)
	}
	c, err = OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	got, ok, err := c.ResolveSession(boundedTestContext(t), want.CredentialID)
	if err != nil || !ok || !reflect.DeepEqual(want, got) {
		t.Fatalf("%+v %v %v", got, ok, err)
	}
	if err := c.MintSessionWithStart(boundedTestContext(t), row, start); err != nil {
		t.Fatal(err)
	}
	// A duplicate episode must roll back the credential insert as well.
	duplicate := row
	duplicate.CredentialID = "must-not-land"
	if err := c.MintSessionWithStart(boundedTestContext(t), duplicate, start); err == nil {
		t.Fatal("overwrote start")
	}
	if _, ok, err := c.ResolveSession(boundedTestContext(t), duplicate.CredentialID); err != nil || ok {
		t.Fatalf("credential escaped rollback: %v %v", ok, err)
	}
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	record, err := ro.EpisodeStart(boundedTestContext(t), "new")
	if err != nil || !reflect.DeepEqual(record, start) {
		t.Fatalf("%+v %v", record, err)
	}
	if err := c.RevokeSession(boundedTestContext(t), row.CredentialID); err != nil {
		t.Fatal(err)
	}
	if record, err := ro.EpisodeStart(boundedTestContext(t), "new"); err != nil || record == nil {
		t.Fatalf("revoke lost start: %v", err)
	}
}
