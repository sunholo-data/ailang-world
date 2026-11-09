package store

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Row 220 M3: keyset pages read every row, and the start projection never
// carries preset bytes, paths, the remote or checks.
func TestOverlayReadsPageEveryRow(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	w, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = w.Close() }()
	ctx := boundedTestContext(t)
	const n = MaxOverlayPage + 7
	for i := 0; i < n; i++ {
		ep := fmt.Sprintf("e%04d", i)
		row := SessionRow{CredentialID: fmt.Sprintf("c%04d", i), EpisodeID: ep, GrantsJSON: "[]", ExpiresAt: 100, CreatedAt: int64(i / 3)}
		start := &EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: "2026-10-29T12:00:00Z",
			Preset: StartPreset{Name: "p", Path: "/home/operator/presets/p.json", SHA256: "sha256:x", Bytes: `{"secret":"preset-bytes"}`},
			Repo:   "/home/operator/dev/world", Remote: StartRemote{Name: "origin", URL: "https://operator:tok@example.invalid/r.git"},
			Base: StartBase{Branch: "dev", SHA: strings.Repeat("a", 40)}, Branch: "trial/" + ep, Agent: "claude", Model: "opus",
			Task:   &StartTask{Path: "/home/operator/task.md", SHA256: "sha256:t", Title: "Task " + ep},
			Checks: []StartCheck{{Name: "go", Argv: []string{"go", "test"}, Cwd: "/home/operator/dev/world"}}, TTLS: 3600}
		if i == 0 {
			start.Resumes = []EpisodeResume{{CredentialID: "r", ResumedFrom: "c0000", ResumedAt: "2026-10-29T13:00:00Z", TTLS: 60}}
		}
		if err := w.MintSessionWithStart(ctx, row, start); err != nil {
			t.Fatal(err)
		}
	}
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()

	var sessions []SessionRow
	var cursor SessionCursor
	pages := 0
	for {
		page, next, err := ro.ListSessionsAfter(ctx, cursor, 0)
		if err != nil {
			t.Fatal(err)
		}
		pages++
		sessions = append(sessions, page...)
		cursor = next
		if len(page) < MaxOverlayPage {
			break
		}
	}
	all, err := ro.ListSessions(ctx, "", MaxListSessions)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != n || pages != 2 || !reflect.DeepEqual(sessions[:len(all)], all) {
		t.Fatalf("paged %d rows in %d pages", len(sessions), pages)
	}
	seen := map[string]bool{}
	for _, s := range sessions {
		if seen[s.CredentialID] {
			t.Fatalf("row %s read twice across equal created_at", s.CredentialID)
		}
		seen[s.CredentialID] = true
	}

	var starts []EpisodeStartSummary
	after := ""
	for {
		page, err := ro.ListEpisodeStartSummaries(ctx, after, 0)
		if err != nil {
			t.Fatal(err)
		}
		starts = append(starts, page...)
		if len(page) < MaxOverlayPage {
			break
		}
		after = page[len(page)-1].Episode
	}
	if len(starts) != n {
		t.Fatalf("starts %d, want %d", len(starts), n)
	}
	want := EpisodeStartSummary{Episode: "e0000", StartedAt: "2026-10-29T12:00:00Z", Repo: "/home/operator/dev/world",
		BaseSHA: strings.Repeat("a", 40), Branch: "trial/e0000", Agent: "claude", Model: "opus", TaskTitle: "Task e0000", TTLS: 3600,
		Resumes: []EpisodeResume{{CredentialID: "r", ResumedFrom: "c0000", ResumedAt: "2026-10-29T13:00:00Z", TTLS: 60}}}
	if !reflect.DeepEqual(starts[0], want) {
		t.Fatalf("summary\n got %+v\nwant %+v", starts[0], want)
	}
	text := fmt.Sprintf("%+v", starts)
	for _, secret := range []string{"preset-bytes", "presets/p.json", "operator:tok", "example.invalid", "task.md", "go test"} {
		if strings.Contains(text, secret) {
			t.Fatalf("summary carries %q", secret)
		}
	}
}

func TestOverlayStartReadOnOlderStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	db := rawDB(t, path)
	if _, err := db.Exec(canonicalCredentialsDDL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 2)
	_ = db.Close()
	ro, err := OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	if got, err := ro.ListEpisodeStartSummaries(boundedTestContext(t), "", 0); err != nil || got != nil {
		t.Fatalf("%v %v", got, err)
	}
	if _, err := ro.ListEpisodeStartSummaries(nil, "", 0); err == nil {
		t.Fatal("no deadline accepted")
	}
}
