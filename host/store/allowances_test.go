package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func hex64(seed string) string {
	sum := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(sum[:])
}

func TestStandingAllowancesLifecycleAndDailyCap(t *testing.T) {
	ctx := testCtx(t)
	cs := openTestCredentials(t)
	if cs.Version() != credentialsSchemaVersion {
		t.Fatalf("version = %d, want %d", cs.Version(), credentialsSchemaVersion)
	}

	payload := []byte(`{"schema":"world/standing-allowance/v1","preset_name":"docparse"}`)
	row := StandingAllowanceRow{
		AllowanceID:    "sa-0123456789abcdef",
		PresetName:     "docparse",
		PresetSHA256:   hex64("preset-bytes"),
		ProfileProject: "docparse",
		ProfileDigest:  "sha256:" + hex64("profile-bytes"),
		Payload:        payload,
		Digest:         hex64(string(payload)),
		Challenge:      hex64("challenge-1"),
		MaxPerDay:      2,
		TTLS:           43200,
		ApprovedAt:     1_800_000_000,
		ExpiresAt:      1_800_000_000 + 30*86400,
		ApprovedBy:     "tty",
	}
	if err := cs.GrantAllowance(ctx, row); err != nil {
		t.Fatalf("GrantAllowance: %v", err)
	}

	// Second active allowance for the same preset is refused.
	dup := row
	dup.AllowanceID = "sa-fedcba9876543210"
	if err := cs.GrantAllowance(ctx, dup); !errors.Is(err, ErrAllowanceExists) {
		t.Fatalf("duplicate active grant = %v, want ErrAllowanceExists", err)
	}

	got, ok, err := cs.GetActiveAllowance(ctx, "docparse")
	if err != nil || !ok || got.AllowanceID != row.AllowanceID {
		t.Fatalf("GetActiveAllowance = %+v %v %v", got, ok, err)
	}

	// Draw 1 and Draw 2 succeed atomically; Draw 3 on the same day hits the daily cap (2).
	for i := 1; i <= 3; i++ {
		ep := fmt.Sprintf("ep-draw-%d", i)
		sess := SessionRow{
			CredentialID:  hex64("cred-" + ep),
			EpisodeID:     ep,
			GrantsJSON:    `[{"effect":"Workspace.Read","scope":"worktree","budget":5000}]`,
			ExpiresAt:     1_800_043_200,
			CreatedAt:     1_800_000_000 + int64(i),
			Project:       "docparse",
			ProfileDigest: row.ProfileDigest,
		}
		start := EpisodeStart{
			Schema:    "world/episode-start/v1",
			Episode:   ep,
			StartedAt: "2026-11-25T10:00:00Z",
			Preset:    StartPreset{Name: "docparse", Path: "/tmp/docparse.json", SHA256: "sha256:" + row.PresetSHA256, Bytes: "{}"},
			Repo:      "/tmp/docparse",
			Branch:    "world/" + ep,
			Agent:     "agy",
			TTLS:      43200,
		}
		draw := AllowanceDrawRow{
			DrawID:       fmt.Sprintf("ad-%016x", i),
			AllowanceID:  row.AllowanceID,
			DayKey:       "2026-11-25",
			EpisodeID:    ep,
			CredentialID: sess.CredentialID,
			Repo:         "/tmp/docparse",
			DrawnAt:      1_800_000_000 + int64(i),
			TTLS:         43200,
		}
		num, err := cs.MintSessionFromAllowance(ctx, sess, start, draw, row.MaxPerDay)
		if i <= 2 {
			if err != nil || num != i {
				t.Fatalf("draw %d: num=%d err=%v", i, num, err)
			}
			st, err := cs.EpisodeStart(ctx, ep)
			if err != nil || st == nil || st.Allowance == nil || st.Allowance.DrawNumber != i || st.Allowance.AllowanceID != row.AllowanceID {
				t.Fatalf("EpisodeStart %s: %+v %v", ep, st, err)
			}
		} else if !errors.Is(err, ErrAllowanceDailyCap) {
			t.Fatalf("draw 3 err = %v, want ErrAllowanceDailyCap", err)
		}
	}

	if n, err := cs.CountAllowanceDraws(ctx, row.AllowanceID, "2026-11-25"); err != nil || n != 2 {
		t.Fatalf("CountAllowanceDraws = %d %v, want 2", n, err)
	}

	// Revoke by preset name; active lookup disappears, subsequent draw on next day refuses,
	// and a new allowance for the same preset can now be granted.
	rev, err := cs.RevokeAllowance(ctx, "docparse", 1_800_001_000, "tty", "test revoke")
	if err != nil || rev.RevokedAt == nil || !strings.Contains(rev.RevokedReason, "test revoke") {
		t.Fatalf("RevokeAllowance = %+v %v", rev, err)
	}
	if _, ok, err := cs.GetActiveAllowance(ctx, "docparse"); err != nil || ok {
		t.Fatalf("GetActiveAllowance after revoke: ok=%v err=%v", ok, err)
	}
	if _, err := cs.MintSessionFromAllowance(ctx, SessionRow{
		CredentialID: hex64("cred-after-rev"), EpisodeID: "ep-after-rev", GrantsJSON: "[]",
		ExpiresAt: 1_800_100_000, CreatedAt: 1_800_090_000,
	}, EpisodeStart{Schema: "world/episode-start/v1", Episode: "ep-after-rev"}, AllowanceDrawRow{
		DrawID: "ad-0000000000000009", AllowanceID: row.AllowanceID, DayKey: "2026-11-26",
		EpisodeID: "ep-after-rev", CredentialID: hex64("cred-after-rev"), Repo: "/tmp/docparse",
		DrawnAt: 1_800_090_000, TTLS: 43200,
	}, 2); !errors.Is(err, ErrAllowanceRevoked) {
		t.Fatalf("draw after revoke = %v, want ErrAllowanceRevoked", err)
	}
	if err := cs.GrantAllowance(ctx, dup); err != nil {
		t.Fatalf("re-grant after revoke: %v", err)
	}
}
