// Package authority owns the inbound Authorization: Bearer credential boundary
// (w-session-authority D5). It mints, revokes and resolves session credentials
// locally-first over the store's persisted session_credentials table.
//
// # Timing safety — equality is the indexed lookup, NOT a Go-level compare
//
// The resolved equality of a presented credential is performed by SQLite's
// indexed point lookup on the TEXT PRIMARY KEY (WHERE credential_id = ? where
// credential_id is sha256(token_hex)) — never by a Go byte-for-byte compare of
// token or hash material. crypto/subtle is deliberately NOT used: after an
// index hit a constant-time compare is vacuous (it would compare the query
// parameter against the row that matched it — a value against itself), and the
// stored hashes of unknown 256-bit random tokens leak nothing via a B-tree
// prefix timing difference. Guarded by TestAuthority_NoDirectTokenCompare.
package authority

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// DenialKind is the typed, fail-closed classification of why a credential did
// not resolve to a session. Each kind is distinct and independently
// distinguishable so a consumer can tell absent from malformed from unknown
// from expired (F6).
type DenialKind uint8

const (
	// DenialAbsent — no Authorization: Bearer header at all.
	DenialAbsent DenialKind = iota
	// DenialMalformed — header present but not "Bearer <64-hex-token>".
	DenialMalformed
	// DenialUnknown — token is present-shaped but has no mapping row (or its
	// row was deleted by revocation).
	DenialUnknown
	// DenialExpired — mapping row exists but expires_at < now.
	DenialExpired
)

// String maps a denial kind to the operator-facing class name so HTTP wiring
// (M2 middleware) and diagnostics share one vocabulary.
func (k DenialKind) String() string {
	switch k {
	case DenialAbsent:
		return "SessionAbsent"
	case DenialMalformed:
		return "InvalidSession"
	case DenialUnknown:
		return "SessionUnknown"
	case DenialExpired:
		return "SessionExpired"
	default:
		return fmt.Sprintf("DenialKind(%d)", uint8(k))
	}
}

// SessionBinding is the credential->(episode, grants, expiry) mapping this
// sprint hands to the row-40 consumer (F6). CredentialID (the stored hash) is
// deliberately NOT exposed: no consumer needs the secret's hash as a capability
// signal.
type SessionBinding struct {
	EpisodeID string
	Caps      []broker.Capability
	ExpiresAt int64
	CreatedAt int64
}

// ResolveOutcome is the typed result of credential resolution. Exactly one of
// Success / Denied is set per Resolve call; it is never (nil, nil), so the
// unauthenticated-degrade path is structurally impossible.
type ResolveOutcome struct {
	Success *SessionBinding
	Denied  *DenialKind
}

// CredentialStore is the minimal store seam Resolve needs. *store.Store
// satisfies it. It exists so the bounded-lookup test can inject a counting
// observer over exactly the one resolve query.
type CredentialStore interface {
	ResolveSession(ctx context.Context, credentialID string) (store.SessionRow, bool, error)
}

// Resolver maps an opaque Authorization: Bearer credential to a binding.
type Resolver interface {
	// Resolve maps the ENTIRE Authorization header value (or "" if absent)
	// plus a caller-supplied now (Unix seconds) to a typed outcome.
	Resolve(header string, now int64) ResolveOutcome
	// ResolveContext is Resolve with a caller-supplied context threading
	// one bounded request deadline into the store's single-connection wait
	// (w-a2a-session-projection P6.A-CTX). Credential policy is IDENTICAL to
	// Resolve; the ONE behavioural difference is that a store read failure
	// (including context.Canceled / context.DeadlineExceeded while the query
	// waits for the pooled connection) is returned as a non-nil error with a
	// zero outcome — never collapsed into DenialUnknown — so a transport
	// deadline can never be misreported to the caller as a bad credential.
	ResolveContext(ctx context.Context, header string, now int64) (ResolveOutcome, error)
}

// New builds a Resolver over a CredentialStore.
func New(st CredentialStore) Resolver {
	return &resolver{st: st}
}

type resolver struct {
	st CredentialStore
}

// tokenHexLen is the length of the minted 64-lowercase-hex credential token
// (32 bytes from crypto/rand, D3).
const tokenHexLen = 64

// Resolve implements Resolver with the bounded hash-then-index lookup. An
// absent header is refused before any store touch (bounded by construction);
// a present-shaped token is hashed once and looked up by the single indexed
// PK query (bounded by design).
//
// Resolve is exactly ResolveContext over context.Background(): the credential
// policy lives in ONE place so the two entry points can never drift apart
// (P6.A-CTX policy equivalence), and the /v1/commit middleware keeps today's
// context-less behaviour byte-for-byte, including the store-error ->
// DenialUnknown collapse (a residual the projection does not inherit).
func (r *resolver) Resolve(header string, now int64) ResolveOutcome {
	out, err := r.ResolveContext(context.Background(), header, now)
	if err != nil {
		// A store read failure is not a denial the caller can act on; fail
		// closed to the same surface as an unknown credential.
		denied := DenialUnknown
		return ResolveOutcome{Denied: &denied}
	}
	return out
}

// ResolveContext implements the resolve policy once for both entry points.
// Everything up to the store call (header shape, token shape, hash) is pure
// and identical to Resolve; the store call receives the caller's ctx so a
// bounded request deadline bounds the wait for the store's single pooled
// connection (store.SetMaxOpenConns(1)), and a store error surfaces as an
// error, not a credential denial.
func (r *resolver) ResolveContext(ctx context.Context, header string, now int64) (ResolveOutcome, error) {
	if header == "" {
		denied := DenialAbsent
		return ResolveOutcome{Denied: &denied}, nil
	}
	fields := splitHeader(header)
	if len(fields) != 2 || fields[0] != "Bearer" {
		denied := DenialMalformed
		return ResolveOutcome{Denied: &denied}, nil
	}
	token := fields[1]
	if !isHexToken(token) {
		denied := DenialMalformed
		return ResolveOutcome{Denied: &denied}, nil
	}
	credentialID := hashHexToken(token)
	row, ok, err := r.st.ResolveSession(ctx, credentialID)
	if err != nil {
		return ResolveOutcome{}, fmt.Errorf("authority: resolve: %w", err)
	}
	if !ok {
		denied := DenialUnknown
		return ResolveOutcome{Denied: &denied}, nil
	}
	if now > row.ExpiresAt {
		denied := DenialExpired
		return ResolveOutcome{Denied: &denied}, nil
	}
	grants, err := unmarshalGrants(row.GrantsJSON)
	if err != nil {
		denied := DenialUnknown
		return ResolveOutcome{Denied: &denied}, nil
	}
	return ResolveOutcome{Success: &SessionBinding{
		EpisodeID: row.EpisodeID,
		Caps:      grants,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}}, nil
}

// splitHeader tokenizes the Authorization header value.
func splitHeader(header string) []string {
	var fields []string
	cur := ""
	for i := 0; i < len(header); i++ {
		c := header[i]
		switch {
		case c == ' ' || c == '\t':
			if cur != "" {
				fields = append(fields, cur)
				cur = ""
			}
		default:
			cur += string(c)
		}
	}
	if cur != "" {
		fields = append(fields, cur)
	}
	return fields
}

// isHexToken reports whether s is exactly 64 hex characters (case-insensitive).
func isHexToken(s string) bool {
	if len(s) != tokenHexLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		valid := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !valid {
			return false
		}
	}
	return true
}

// hashHexToken returns the 64-lowercase-hex SHA-256 of the presented token's
// hex string (D3's stored form). No comparison against any other value happens
// here; equality is the SQLite indexed PK lookup done by the store.
func hashHexToken(tokenHex string) string {
	sum := sha256.Sum256([]byte(tokenHex))
	return hex.EncodeToString(sum[:])
}

// unmarshalGrants decodes the opaque grants_json persisted at mint back into a
// []broker.Capability. Any decode failure is treated as an unknown credential
// at resolve (the row is unusable), never as a success.
func unmarshalGrants(jsonText string) ([]broker.Capability, error) {
	var grants []broker.Capability
	if err := json.Unmarshal([]byte(jsonText), &grants); err != nil {
		return nil, fmt.Errorf("authority: resolve: unmarshal grants: %w", err)
	}
	return grants, nil
}

// sessionBindingCtxKey is the unexported, typed context key that carries the
// resolved binding from the middleware to the handler (D6). Never a global
// string bag.
type sessionBindingCtxKey struct{}

// WithBinding returns ctx carrying the resolved session binding. It is the
// middleware-side setter (M2 wires it on success).
func WithBinding(ctx context.Context, b *SessionBinding) context.Context {
	return context.WithValue(ctx, sessionBindingCtxKey{}, b)
}

// FromContext returns the resolved session binding carried in ctx, if any.
// Consumers (handleCommit, M2) call it to read episode/grants without ever
// seeing a credential hash.
func FromContext(ctx context.Context) (*SessionBinding, bool) {
	b, ok := ctx.Value(sessionBindingCtxKey{}).(*SessionBinding)
	return b, ok
}
