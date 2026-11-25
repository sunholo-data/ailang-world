package authority

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// AllowanceMintRequest carries the preset and start facts needed to draw a
// session from an approved standing allowance without prompting the operator
// (row 224, D-WORLD-96).
type AllowanceMintRequest struct {
	PresetName   string
	PresetSHA256 string // 64 lowercase hex (without "sha256:" prefix)
	Repo         string
	Caller       string // "" or host adapter; "mission-loop" is always refused
	Start        store.EpisodeStart
}

// MintFromAllowance validates the active standing allowance for req.PresetName
// against every D1-D8 bound (approval.AdmitAllowanceDraw) and atomically writes
// session_credentials, episode_starts (with Allowance provenance), and
// allowance_draws in one BEGIN IMMEDIATE transaction.
func MintFromAllowance(
	ctx context.Context,
	st *store.Credentials,
	episodeID string,
	binding ProfileBinding,
	grants []broker.Capability,
	ttl, now int64,
	out io.Writer,
	req AllowanceMintRequest,
) (rawToken, credentialID string, row store.SessionRow, drawNum int, allowanceID string, err error) {
	if (binding.Project == "") != (binding.Digest == "") {
		return "", "", store.SessionRow{}, 0, "", fmt.Errorf("authority: mint from allowance: a profile binding needs both project and digest")
	}
	if episodeID == "" {
		return "", "", store.SessionRow{}, 0, "", fmt.Errorf("authority: mint from allowance: empty episode id")
	}
	if len(grants) == 0 {
		return "", "", store.SessionRow{}, 0, "", fmt.Errorf("authority: mint from allowance: no grants")
	}

	al, ok, err := st.GetActiveAllowance(ctx, req.PresetName)
	if err != nil {
		return "", "", store.SessionRow{}, 0, "", err
	}
	if !ok {
		return "", "", store.SessionRow{}, 0, "", fmt.Errorf("%w: no active standing allowance for preset %q", store.ErrAllowanceNotFound, req.PresetName)
	}

	dayKey := approval.AllowanceDayKey(now)
	drawsToday, err := st.CountAllowanceDraws(ctx, al.AllowanceID, dayKey)
	if err != nil {
		return "", "", store.SessionRow{}, 0, "", err
	}

	approverRevoked := false
	if al.ApprovedBy != "tty" {
		ap, found, err := st.Approver(ctx, al.ApprovedBy)
		if err != nil {
			return "", "", store.SessionRow{}, 0, "", err
		}
		if !found || ap.RevokedAt != 0 {
			approverRevoked = true
		}
	}

	apGrants := make([]approval.Grant, len(grants))
	for i, g := range grants {
		apGrants[i] = approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget}
	}
	if err := approval.AdmitAllowanceDraw(approval.AllowanceDrawInput{
		RawPayload:      al.Payload,
		StoredDigestHex: al.Digest,
		ApprovedAt:      al.ApprovedAt,
		ExpiresAt:       al.ExpiresAt,
		RevokedAt:       al.RevokedAt,
		ApproverRevoked: approverRevoked,
		Now:             now,
		Caller:          req.Caller,
		PresetName:      req.PresetName,
		PresetSHA256:    req.PresetSHA256,
		Profile:         approval.ProfileRef{Project: binding.Project, Digest: binding.Digest},
		Repo:            req.Repo,
		Grants:          apGrants,
		TTLSeconds:      ttl,
		DrawsToday:      drawsToday,
	}); err != nil {
		return "", "", store.SessionRow{}, 0, al.AllowanceID, err
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", store.SessionRow{}, 0, al.AllowanceID, fmt.Errorf("authority: mint from allowance: read random: %w", err)
	}
	rawToken = hex.EncodeToString(raw)
	credentialID = hashHexToken(rawToken)

	drawRand := make([]byte, 8)
	if _, err := rand.Read(drawRand); err != nil {
		return "", "", store.SessionRow{}, 0, al.AllowanceID, fmt.Errorf("authority: mint from allowance: read draw id: %w", err)
	}
	drawID := "ad-" + hex.EncodeToString(drawRand)

	expiresAt := now + ttl
	clamped := clampGrantExpiry(grants, expiresAt)
	grantsJSON, err := json.Marshal(clamped)
	if err != nil {
		return "", "", store.SessionRow{}, 0, al.AllowanceID, fmt.Errorf("authority: mint from allowance: marshal grants: %w", err)
	}
	row = store.SessionRow{
		CredentialID:  credentialID,
		EpisodeID:     episodeID,
		GrantsJSON:    string(grantsJSON),
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		Project:       binding.Project,
		ProfileDigest: binding.Digest,
	}
	drawRow := store.AllowanceDrawRow{
		DrawID:       drawID,
		AllowanceID:  al.AllowanceID,
		DayKey:       dayKey,
		EpisodeID:    episodeID,
		CredentialID: credentialID,
		Repo:         req.Repo,
		DrawnAt:      now,
		TTLS:         ttl,
	}
	drawNum, err = st.MintSessionFromAllowance(ctx, row, req.Start, drawRow, al.MaxPerDay)
	if err != nil {
		return "", "", store.SessionRow{}, 0, al.AllowanceID, err
	}
	if out != nil {
		if _, err := fmt.Fprintln(out, rawToken); err != nil {
			return "", "", store.SessionRow{}, 0, al.AllowanceID, fmt.Errorf("authority: mint from allowance: write token: %w", err)
		}
	}
	return rawToken, credentialID, row, drawNum, al.AllowanceID, nil
}
