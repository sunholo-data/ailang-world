package authority

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestMintResumedStoresAndRollsBack(t *testing.T) {
	st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	ctx := boundedTestContext(t)
	start := &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: "ep"}
	_, old, _, err := MintStarted(ctx, st, "ep", ProfileBinding{}, testGrants(), 10, 1000, nil, start)
	if err != nil {
		t.Fatal(err)
	}
	resume := store.EpisodeResume{ResumedFrom: old, ResumedAt: "2026-10-30T00:00:00Z", TTLS: 50, PresetFileSHA256: "sha256:now"}
	raw, id, row, err := MintResumed(ctx, st, "ep", ProfileBinding{}, testGrants(), 50, 2000, nil, resume)
	if err != nil || row.CredentialID != id || id != hashHexToken(raw) {
		t.Fatal(row, err)
	}
	if _, ok, err := st.ResolveSession(ctx, id); err != nil || !ok {
		t.Fatal(ok, err)
	}
	got, err := st.EpisodeStart(ctx, "ep")
	if err != nil || len(got.Resumes) != 1 || got.Resumes[0] != resume {
		t.Fatal(got, err)
	}
	before, err := st.ListSessions(ctx, "", store.MaxListSessions)
	if err != nil {
		t.Fatal(err)
	}
	if raw, id, _, err := MintResumed(ctx, st, "missing", ProfileBinding{}, testGrants(), 50, 2000, nil, resume); !errors.Is(err, sql.ErrNoRows) || raw != "" || id != "" {
		t.Fatal(raw, id, err)
	}
	after, err := st.ListSessions(ctx, "", store.MaxListSessions)
	if err != nil || len(before) != len(after) {
		t.Fatal("failed mint persisted", err)
	}
}
