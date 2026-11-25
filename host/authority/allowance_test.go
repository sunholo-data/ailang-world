package authority

import (
	"bytes"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

func seedAllowanceFixture(t *testing.T, st *store.Credentials, maxPerDay int, now int64) (presetSHA, profileDigest string, grants []broker.Capability) {
	t.Helper()
	presetSHA = strings.Repeat("a", 64)
	profileDigest = "sha256:" + strings.Repeat("b", 64)
	grants = []broker.Capability{
		{Effect: "Workspace.Read", Scope: "*", Budget: 50},
		{Effect: "Workspace.Exec", Scope: "*", Budget: 20},
	}
	al := approval.StandingAllowance{
		Schema:       approval.AllowanceSchema,
		PresetName:   "world-agy",
		PresetSHA256: presetSHA,
		Profile:      approval.ProfileRef{Project: "ailang-world", Digest: profileDigest},
		Repos:        []string{"/repo/ailang-world"},
		Grants: []approval.Grant{
			{Effect: "Workspace.Read", Scope: "*", Budget: 50},
			{Effect: "Workspace.Exec", Scope: "*", Budget: 20},
		},
		TTLSeconds:         7200,
		MaxPerDay:          maxPerDay,
		ValidDays:          30,
		WorldDB:            "/state/world.db",
		CredentialsDB:      "/state/credentials.db",
		WorkspaceRoot:      "/state/workspaces",
		ExcludeMissionLoop: true,
	}
	payload, err := approval.CanonicalAllowance(al)
	if err != nil {
		t.Fatalf("CanonicalAllowance: %v", err)
	}
	digest, err := approval.AllowanceDigest(al)
	if err != nil {
		t.Fatalf("AllowanceDigest: %v", err)
	}
	err = st.GrantAllowance(boundedTestContext(t), store.StandingAllowanceRow{
		AllowanceID:    "sa-0123456789abcdef",
		PresetName:     "world-agy",
		PresetSHA256:   presetSHA,
		ProfileProject: "ailang-world",
		ProfileDigest:  profileDigest,
		Payload:        payload,
		Digest:         hex.EncodeToString(digest[:]),
		Challenge:      strings.Repeat("c", 64),
		MaxPerDay:      maxPerDay,
		TTLS:           7200,
		ApprovedAt:     now - 60,
		ExpiresAt:      (now - 60) + 30*86400,
		ApprovedBy:     "tty",
	})
	if err != nil {
		t.Fatalf("GrantAllowance: %v", err)
	}
	return presetSHA, profileDigest, grants
}

func TestMintFromAllowanceBoundsAndProvenance(t *testing.T) {
	ctx := boundedTestContext(t)
	credPath := filepath.Join(t.TempDir(), "credentials.db")
	st, err := store.OpenCredentials(credPath)
	if err != nil {
		t.Fatalf("OpenCredentials: %v", err)
	}
	defer func() { _ = st.Close() }()

	const now int64 = 1_776_000_000
	presetSHA, profileDigest, grants := seedAllowanceFixture(t, st, 2, now)
	binding := ProfileBinding{Project: "ailang-world", Digest: profileDigest}
	baseReq := AllowanceMintRequest{
		PresetName:   "world-agy",
		PresetSHA256: presetSHA,
		Repo:         "/repo/ailang-world",
		Start: store.EpisodeStart{
			Episode: "ep-1",
			Repo:    "/repo/ailang-world",
			Branch:  "world/ep-1",
			Agent:   "agy",
		},
	}

	// 1. Mission-loop caller is always refused (D-WORLD-96).
	mlReq := baseReq
	mlReq.Caller = "mission-loop"
	if _, _, _, _, _, err := MintFromAllowance(ctx, st, "ep-ml", binding, grants, 7200, now, &bytes.Buffer{}, mlReq); err == nil {
		t.Fatalf("expected mission-loop draw to be refused")
	}

	// 2. Preset digest drift is refused.
	driftReq := baseReq
	driftReq.PresetSHA256 = strings.Repeat("f", 64)
	if _, _, _, _, _, err := MintFromAllowance(ctx, st, "ep-drift", binding, grants, 7200, now, &bytes.Buffer{}, driftReq); err == nil {
		t.Fatalf("expected drifted preset digest to be refused")
	}

	// 3. Draw 1 succeeds and records Allowance provenance on EpisodeStart.
	var out bytes.Buffer
	tok1, credID1, row1, draw1, alID1, err := MintFromAllowance(ctx, st, "ep-1", binding, grants, 7200, now, &out, baseReq)
	if err != nil {
		t.Fatalf("draw 1 failed: %v", err)
	}
	if len(tok1) != 64 || credID1 == "" || row1.EpisodeID != "ep-1" || draw1 != 1 || alID1 != "sa-0123456789abcdef" {
		t.Fatalf("unexpected draw 1 result: tokLen=%d credID=%q ep=%q draw=%d alID=%q", len(tok1), credID1, row1.EpisodeID, draw1, alID1)
	}
	start1, err := st.EpisodeStart(ctx, "ep-1")
	if err != nil || start1 == nil {
		t.Fatalf("EpisodeStart(ep-1): start=%v err=%v", start1, err)
	}
	if start1.Allowance == nil || start1.Allowance.AllowanceID != alID1 || start1.Allowance.DrawNumber != 1 || start1.Allowance.ApprovedBy != "tty" || start1.Allowance.DrawID == "" {
		t.Fatalf("unexpected start1.Allowance: %+v", start1.Allowance)
	}

	// 4. Draw 2 succeeds.
	req2 := baseReq
	req2.Start.Episode = "ep-2"
	req2.Start.Branch = "world/ep-2"
	_, _, _, draw2, _, err := MintFromAllowance(ctx, st, "ep-2", binding, grants, 7200, now+10, &out, req2)
	if err != nil {
		t.Fatalf("draw 2 failed: %v", err)
	}
	if draw2 != 2 {
		t.Fatalf("draw 2 got drawNum=%d, want 2", draw2)
	}

	// 5. Draw 3 exceeds daily cap (2/day) and refuses.
	req3 := baseReq
	req3.Start.Episode = "ep-3"
	req3.Start.Branch = "world/ep-3"
	if _, _, _, _, _, err := MintFromAllowance(ctx, st, "ep-3", binding, grants, 7200, now+20, &out, req3); err == nil {
		t.Fatalf("expected draw 3 over daily cap to be refused")
	}

	// 6. Revoking invalidates future draws even on the next day.
	if _, err := st.RevokeAllowance(ctx, alID1, now+30, "tty", "test revoke"); err != nil {
		t.Fatalf("RevokeAllowance: %v", err)
	}
	if _, _, _, _, _, err := MintFromAllowance(ctx, st, "ep-4", binding, grants, 7200, now+86400, &out, req3); !errors.Is(err, store.ErrAllowanceNotFound) {
		t.Fatalf("expected ErrAllowanceNotFound after revoke, got %v", err)
	}
}
