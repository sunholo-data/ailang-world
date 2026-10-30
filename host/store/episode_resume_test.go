package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMintSessionWithResumeAtomicAndPreservesJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ctx := boundedTestContext(t)
	original := `{"schema":"world/episode-start/v1","episode":"ep","future":{"n":9007199254740993,"null":null},"task":null,"preset":{"bytes":"exact\\nbytes"},"launches":[{"exit":7}]}`
	if _, err := c.db.ExecContext(ctx, "INSERT INTO episode_starts VALUES (?, ?)", "ep", original); err != nil {
		t.Fatal(err)
	}
	resume := EpisodeResume{ResumedFrom: "old", ResumedAt: "2026-10-30T00:00:00Z", TTLS: 50, PresetFileSHA256: "sha256:now"}
	for i, id := range []string{"new", "newer"} {
		row := SessionRow{CredentialID: id, EpisodeID: "ep", GrantsJSON: "[]", ExpiresAt: 100, CreatedAt: int64(i + 1)}
		if i == 1 {
			row.Project, row.ProfileDigest = "p", "sha256:p"
		}
		if err := c.MintSessionWithResume(ctx, row, resume); err != nil {
			t.Fatal(err)
		}
		got, ok, err := c.ResolveSession(ctx, id)
		if err != nil || !ok || !reflect.DeepEqual(row, got) {
			t.Fatal(got, ok, err)
		}
	}
	var data string
	if err := c.db.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = 'ep'").Scan(&data); err != nil {
		t.Fatal(err)
	}
	var before, after map[string]json.RawMessage
	_ = json.Unmarshal([]byte(original), &before)
	_ = json.Unmarshal([]byte(data), &after)
	var resumes []EpisodeResume
	if err := json.Unmarshal(after["resumes"], &resumes); err != nil || len(resumes) != 2 || resumes[0] != resume || resumes[1] != resume {
		t.Fatal(resumes, err)
	}
	delete(after, "resumes")
	if !reflect.DeepEqual(before, after) {
		t.Fatal("mint fields changed", before, after)
	}
	missing := SessionRow{CredentialID: "missing", EpisodeID: "missing", GrantsJSON: "[]"}
	if err := c.MintSessionWithResume(ctx, missing, resume); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	if _, ok, err := c.ResolveSession(ctx, "missing"); err != nil || ok {
		t.Fatal("failed mint escaped transaction", ok, err)
	}
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	if err := ro.MintSessionWithResume(ctx, missing, resume); err == nil {
		t.Fatal("read-only handle accepted resume")
	}
}
