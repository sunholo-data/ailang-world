package authority

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Mint creates a new session credential (D1 primitive / D3 format) bound to
// episodeID with the given grants. The raw 32-byte token is rendered as 64
// lowercase hex and written to out EXACTLY once (AC-M1-1); the store persists
// ONLY sha256(token_hex) as credential_id (D3), never the raw token. It
// returns the raw token (so the CLI/--out caller can write it at mode 0600 for
// off-broadcast handoff) and the persisted row's metadata, including the stored
// credential_id hash for diagnostics.
//
// Errors/outcome messaging must never include the raw token — only the
// credential_id hash may leave this function in an error path.
func Mint(ctx context.Context, st *store.Store, episodeID string, grants []broker.Capability, ttl, now int64, out io.Writer) (rawToken, credentialID string, row store.SessionRow, err error) {
	if episodeID == "" {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: empty episode id")
	}
	if ttl < 0 {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: negative ttl %d", ttl)
	}
	if len(grants) == 0 {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: no grants")
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: read random: %w", err)
	}
	// 32 bytes -> 64 lowercase hex (encoding/hex lowercases).
	rawToken = hex.EncodeToString(raw)
	credentialID = hashHexToken(rawToken)

	grantsJSON, err := json.Marshal(grants)
	if err != nil {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: marshal grants: %w", err)
	}
	row = store.SessionRow{
		CredentialID: credentialID,
		EpisodeID:    episodeID,
		GrantsJSON:   string(grantsJSON),
		ExpiresAt:    now + ttl,
		CreatedAt:    now,
	}
	if err := st.MintSession(row); err != nil {
		return "", "", store.SessionRow{}, err
	}

	// Raw token printed exactly once. Nothing else in this path logs it: the
	// returned row carries only the hash, and no error payload carries rawToken.
	if out != nil {
		if _, err := io.WriteString(out, rawToken); err != nil {
			return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: print credential: %w", err)
		}
	}
	return rawToken, credentialID, row, nil
}
