package store

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func resumeStoreR2(t *testing.T) *Credentials {
	t.Helper()
	c, err := OpenCredentials(filepath.Join(t.TempDir(), CredentialsFileName))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	start := &EpisodeStart{Schema: "world/episode-start/v1", Episode: "ep"}
	row := SessionRow{CredentialID: "original", EpisodeID: "ep", GrantsJSON: "[]", CreatedAt: 1, ExpiresAt: 10}
	if err := c.MintSessionWithStart(boundedTestContext(t), row, start); err != nil {
		t.Fatal(err)
	}
	return c
}

func assertResumeRefusalR2(t *testing.T, c *Credentials, row SessionRow, resume EpisodeResume, want string) {
	t.Helper()
	ctx := boundedTestContext(t)
	before, err := c.EpisodeStart(ctx, "ep")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := c.ListSessions(ctx, "ep", MaxListSessions)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.MintSessionWithResume(ctx, row, resume); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatal("wrong refusal", err)
	}
	after, err := c.EpisodeStart(ctx, "ep")
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("history changed", after, err)
	}
	got, err := c.ListSessions(ctx, "ep", MaxListSessions)
	if err != nil || !reflect.DeepEqual(rows, got) {
		t.Fatal("credential escaped transaction", got, err)
	}
}

func TestMintSessionWithResumeRefusesLive(t *testing.T) {
	for _, now := range []int64{9, 10} {
		name := "before-expiry"
		if now == 10 {
			name = "at-expiry"
		}
		t.Run(name, func(t *testing.T) {
			c := resumeStoreR2(t)
			row := SessionRow{CredentialID: "new", EpisodeID: "ep", GrantsJSON: "[]", CreatedAt: now, ExpiresAt: 50}
			assertResumeRefusalR2(t, c, row, EpisodeResume{ResumedFrom: "original"}, "F12")
		})
	}
}

func TestMintSessionWithResumeRefusesStaleCurrent(t *testing.T) {
	c := resumeStoreR2(t)
	ctx := boundedTestContext(t)
	first := SessionRow{CredentialID: "first", EpisodeID: "ep", GrantsJSON: "[]", CreatedAt: 11, ExpiresAt: 20}
	if err := c.MintSessionWithResume(ctx, first, EpisodeResume{ResumedFrom: "original"}); err != nil {
		t.Fatal(err)
	}
	// Even after all credentials expire, an old prompt cannot advance the chain.
	next := SessionRow{CredentialID: "stale", EpisodeID: "ep", GrantsJSON: "[]", CreatedAt: 21, ExpiresAt: 30}
	assertResumeRefusalR2(t, c, next, EpisodeResume{ResumedFrom: "original"}, "current credential changed")
	// Revoking current cannot expose the original expired credential.
	if err := c.RevokeSession(ctx, "first"); err != nil {
		t.Fatal(err)
	}
	assertResumeRefusalR2(t, c, next, EpisodeResume{ResumedFrom: "first"}, "revoke ends an episode")
}
