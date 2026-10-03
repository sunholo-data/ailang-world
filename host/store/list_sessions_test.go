package store

import (
	"path/filepath"
	"testing"
)

// TestListSessionsReadsBesideALiveWriter pins row 138 M5's session list read:
// ordered by created_at, filtered by episode, and served by a read-only
// handle while a writer holds the lock (AC5.6).
func TestListSessionsReadsBesideALiveWriter(t *testing.T) {
	db := filepath.Join(t.TempDir(), "world.db")
	w, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = w.Close() }()
	for i, row := range []SessionRow{
		{CredentialID: "c3", EpisodeID: "ep2", GrantsJSON: "[]", ExpiresAt: 300, CreatedAt: 30},
		{CredentialID: "c1", EpisodeID: "ep1", GrantsJSON: "[]", ExpiresAt: 100, CreatedAt: 10},
		{CredentialID: "c2", EpisodeID: "ep1", GrantsJSON: `[{"Effect":"x"}]`, ExpiresAt: 200, CreatedAt: 20},
	} {
		if err := w.MintSession(boundedTestContext(t), row); err != nil {
			t.Fatalf("mint %d: %v", i, err)
		}
	}
	ro, err := OpenReadOnly(db)
	if err != nil {
		t.Fatalf("read-only open beside a live writer: %v", err)
	}
	defer func() { _ = ro.Close() }()
	all, err := ro.ListSessions(boundedTestContext(t), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 || all[0].CredentialID != "c1" || all[1].CredentialID != "c2" || all[2].CredentialID != "c3" {
		t.Fatalf("all = %+v, want c1,c2,c3 by created_at", all)
	}
	ep1, err := ro.ListSessions(boundedTestContext(t), "ep1", 0)
	if err != nil || len(ep1) != 2 || ep1[1].GrantsJSON != `[{"Effect":"x"}]` {
		t.Fatalf("ep1 = %+v %v", ep1, err)
	}
	one, err := ro.ListSessions(boundedTestContext(t), "", 1)
	if err != nil || len(one) != 1 {
		t.Fatalf("limit 1 = %+v %v", one, err)
	}
	none, err := ro.ListSessions(boundedTestContext(t), "nope", 0)
	if err != nil || len(none) != 0 {
		t.Fatalf("unknown episode = %+v %v", none, err)
	}
}
