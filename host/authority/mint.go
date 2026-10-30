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
//
// st is the credentials store (row 189), opened writable by the attended
// `session` verbs only, after their terminal fence and y/N.
func Mint(ctx context.Context, st *store.Credentials, episodeID string, grants []broker.Capability, ttl, now int64, out io.Writer) (rawToken, credentialID string, row store.SessionRow, err error) {
	return MintBound(ctx, st, episodeID, ProfileBinding{}, grants, ttl, now, out)
}

// ProfileBinding is the exec profile a credential is bound to at the mint
// (row 189 M2, D-189-2 A): the profile's project and Digest, the sha256 of the
// exact profile bytes the human confirmed at the y/N. The zero value is an
// unbound credential.
type ProfileBinding struct {
	Project string
	Digest  string
}

// MintBound is Mint recording a profile binding on the credential. A binding
// needs both fields; the zero binding mints exactly what Mint mints.
func MintBound(ctx context.Context, st *store.Credentials, episodeID string, binding ProfileBinding, grants []broker.Capability, ttl, now int64, out io.Writer) (rawToken, credentialID string, row store.SessionRow, err error) {
	return MintStarted(ctx, st, episodeID, binding, grants, ttl, now, out, nil)
}

// MintStarted keeps the credential and its attended start in one transaction.
func MintStarted(ctx context.Context, st *store.Credentials, episodeID string, binding ProfileBinding, grants []broker.Capability, ttl, now int64, out io.Writer, start *store.EpisodeStart) (rawToken, credentialID string, row store.SessionRow, err error) {
	return mintRecorded(ctx, st, episodeID, binding, grants, ttl, now, out, start, nil)
}

// MintResumed keeps the credential and its resume entry in one transaction.
func MintResumed(ctx context.Context, st *store.Credentials, episodeID string, binding ProfileBinding, grants []broker.Capability, ttl, now int64, out io.Writer, resume store.EpisodeResume) (rawToken, credentialID string, row store.SessionRow, err error) {
	return mintRecorded(ctx, st, episodeID, binding, grants, ttl, now, out, nil, &resume)
}

func mintRecorded(ctx context.Context, st *store.Credentials, episodeID string, binding ProfileBinding, grants []broker.Capability, ttl, now int64, out io.Writer, start *store.EpisodeStart, resume *store.EpisodeResume) (rawToken, credentialID string, row store.SessionRow, err error) {
	if (binding.Project == "") != (binding.Digest == "") {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: a profile binding needs both the project and its digest")
	}
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

	// A grant never outlives its session (row 134 fix): any grant whose
	// ExpiresAt is later than the session's own expiry is clamped to it. The
	// caller's slice is not mutated.
	expiresAt := now + ttl
	grants = clampGrantExpiry(grants, expiresAt)
	grantsJSON, err := json.Marshal(grants)
	if err != nil {
		return "", "", store.SessionRow{}, fmt.Errorf("authority: mint: marshal grants: %w", err)
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
	if resume != nil {
		err = st.MintSessionWithResume(ctx, row, *resume)
	} else {
		err = st.MintSessionWithStart(ctx, row, start)
	}
	if err != nil {
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

// clampGrantExpiry returns a copy of grants with every ExpiresAt later than
// sessionExpiry lowered to sessionExpiry, so no stored grant can outlive the
// session that carries it. Earlier expiries (including 0) are kept as given.
func clampGrantExpiry(grants []broker.Capability, sessionExpiry int64) []broker.Capability {
	out := append([]broker.Capability(nil), grants...)
	for i := range out {
		if out[i].ExpiresAt > sessionExpiry {
			out[i].ExpiresAt = sessionExpiry
		}
	}
	return out
}
