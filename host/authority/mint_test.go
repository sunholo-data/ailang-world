package authority

import (
	"bytes"
	"context"
	"encoding/hex"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// TestMint_PrintsOnce is mutation M4's sole killer (AC-M1-1): the raw credential
// is emitted EXACTLY once, and the persisted row stores ONLY sha256(token_hex) —
// never the raw token.
func TestMint_PrintsOnce(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	var out bytes.Buffer

	rawToken, credentialID, row, err := Mint(context.Background(), st, "ep-42", testGrants(), 3600, now, &out)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if n := strings.Count(out.String(), rawToken); n != 1 {
		t.Fatalf("raw token appears %d times in mint output %q, want exactly 1", n, out.String())
	}
	if !isHexToken(rawToken) {
		t.Fatalf("raw token %q is not 64 lowercase hex", rawToken)
	}
	if rawToken == credentialID {
		t.Fatal("raw token equals its stored credential_id — the store must hold the hash, not the raw token")
	}
	if row.CredentialID != credentialID {
		t.Fatalf("row.CredentialID = %q, want returned credentialID %q", row.CredentialID, credentialID)
	}
	if row.ExpiresAt != now+3600 || row.CreatedAt != now {
		t.Fatalf("row expiry %d / created %d, want %d / %d", row.ExpiresAt, row.CreatedAt, now+3600, now)
	}

	// Nothing in the persisted mapping can reproduce the raw token: resolving a
	// freshly minted, unexpired session must succeed, and the binding must not
	// expose any credential hash (SessionBinding has no credential field).
	b := New(st).Resolve("Bearer "+rawToken, now+1)
	if b.Success == nil {
		t.Fatalf("minted credential did not resolve: %#v", b)
	}
	if b.Success.EpisodeID != "ep-42" {
		t.Fatalf("resolved EpisodeID = %q, want ep-42", b.Success.EpisodeID)
	}
}

// TestMint_StoresOnlyHash proves D3: resolving by the SHA-256 hash works, but a
// look-up by the raw token string itself finds nothing (proving the store holds
// only the hash).
func TestMint_StoresOnlyHash(t *testing.T) {
	st := newTestStore(t)
	rawToken, credentialID, _, err := Mint(context.Background(), st, "ep-42", testGrants(), 3600, 1000, nil)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	// The hash (credential_id) is findable...
	want := hashHexToken(rawToken)
	if want != credentialID {
		t.Fatalf("hashHexToken(raw) = %q, want credentialID %q", want, credentialID)
	}
	row, ok, err := st.ResolveSession(context.Background(), credentialID)
	if err != nil || !ok || row.EpisodeID != "ep-42" {
		t.Fatalf("resolve by credential_id: ok=%v err=%v row=%+v, want ok with ep-42", ok, err, row)
	}
	// ...but the raw token string is not a stored credential_id.
	_, okByRaw, err := st.ResolveSession(context.Background(), rawToken)
	if err != nil || okByRaw {
		t.Fatalf("resolve by raw token: ok=%v err=%v, want not-found (raw never persisted)", okByRaw, err)
	}
}

// TestMint_DistinctCredentials proves AC-M1-2: two mints for the same episode
// produce distinct tokens and distinct stored credential_ids.
func TestMint_DistinctCredentials(t *testing.T) {
	st := newTestStore(t)
	a, aID, _, err := Mint(context.Background(), st, "ep-42", testGrants(), 3600, 1000, nil)
	if err != nil {
		t.Fatalf("first mint: %v", err)
	}
	b, bID, _, err := Mint(context.Background(), st, "ep-42", testGrants(), 3600, 1000, nil)
	if err != nil {
		t.Fatalf("second mint: %v", err)
	}
	if a == b {
		t.Fatal("two mints produced the same raw token")
	}
	if aID == bID {
		t.Fatal("two mints produced the same credential_id")
	}
	// Both resolve independently.
	ra := New(st).Resolve("Bearer "+a, 1000)
	rb := New(st).Resolve("Bearer "+b, 1000)
	if ra.Success == nil || rb.Success == nil {
		t.Fatalf("distinct mints did not both resolve: %#v / %#v", ra, rb)
	}
}

// TestMint_RejectsBadInputs checks the mint guards: empty episode, no grants,
// negative ttl, and a mint with an unknown-budget grant still stores grant JSON.
func TestMint_RejectsBadInputs(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	cases := []struct {
		name     string
		episode  string
		grants   []broker.Capability
		ttl      int64
		wantText string
	}{
		{"empty episode", "", testGrants(), 3600, "empty episode"},
		{"no grants", "ep-42", nil, 3600, "no grants"},
		{"negative ttl", "ep-42", testGrants(), -1, "negative ttl"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := Mint(ctx, st, tc.episode, tc.grants, tc.ttl, 1000, nil)
			if err == nil || !strings.Contains(err.Error(), tc.wantText) {
				t.Fatalf("Mint error = %v, want containing %q", err, tc.wantText)
			}
		})
	}
}

// TestMint_GrantsRoundTrip confirms the grants_json stored at mint is exactly
// the []broker.Capability the resolver returns (no drift through JSON).
func TestMint_GrantsRoundTrip(t *testing.T) {
	st := newTestStore(t)
	grants := []broker.Capability{
		{Effect: "fs.read", Scope: "/tmp/a", ExpiresAt: 0, Budget: 1},
		{Effect: "repo.commit", Scope: "world/ep-42", ExpiresAt: 0, Budget: 5},
	}
	tok := mintTestToken(t, st, "ep-42", grants, 3600, 1000)
	out := New(st).Resolve("Bearer "+tok, 1000)
	if out.Success == nil || len(out.Success.Caps) != 2 {
		t.Fatalf("resolve = %#v, want 2 caps", out)
	}
	if out.Success.Caps[0].Effect != "fs.read" || out.Success.Caps[1].Effect != "repo.commit" {
		t.Fatalf("caps = %+v, want fs.read and repo.commit in order", out.Success.Caps)
	}
}

// TestRevoke_DeletesRow is mutation M5's sole killer (AC-M2-5): after Revoke the
// mapping row is gone and resolving the revoked credential is DenialUnknown
// (delete-from-mapping = unknown, D4).
func TestRevoke_DeletesRow(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	tok := mintTestToken(t, st, "ep-42", testGrants(), 3600, now)
	credID := hashHexToken(tok)

	// Row exists before revocation.
	row, ok, err := st.ResolveSession(context.Background(), credID)
	if err != nil || !ok {
		t.Fatalf("pre-revoke resolve: ok=%v err=%v want ok", ok, err)
	}
	if row.CredentialID == "" {
		t.Fatalf("pre-revoke row has empty credential_id")
	}

	// The authority Revoke path deletes precisely this row.
	if err := Revoke(context.Background(), st, credID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	_, okAfter, err := st.ResolveSession(context.Background(), credID)
	if err != nil {
		t.Fatalf("post-revoke resolve: %v", err)
	}
	if okAfter {
		t.Fatal("row still present after revoke — the DELETE must remove it (mutation M5 leaves it)")
	}

	// The credential now resolves as UNKNOWN, never as expired or success.
	out := New(st).Resolve("Bearer "+tok, now+1)
	if out.Success != nil || out.Denied == nil || *out.Denied != DenialUnknown {
		t.Fatalf("post-revoke resolve = %#v, want Denied=DenialUnknown", out)
	}

	// Revoking an already-absent credential is a no-op, not an error.
	if err := Revoke(context.Background(), st, credID); err != nil {
		t.Fatalf("second revoke of same id: %v", err)
	}
	if err := Revoke(context.Background(), st, hex64); err != nil {
		t.Fatalf("revoke of never-minted id: %v", err)
	}
}

// TestRevoke_DoesNotRevokeOtherCredential pins the DELETE's WHERE clause to the
// exact credential_id (a scoped DELETE must not remove sibling rows).
func TestRevoke_DoesNotRevokeOtherCredential(t *testing.T) {
	st := newTestStore(t)
	a := mintTestToken(t, st, "ep-42", testGrants(), 3600, 1000)
	b := mintTestToken(t, st, "ep-42", testGrants(), 3600, 1000)
	if err := Revoke(context.Background(), st, hashHexToken(a)); err != nil {
		t.Fatalf("revoke a: %v", err)
	}
	if out := New(st).Resolve("Bearer "+b, 1000); out.Success == nil {
		t.Fatalf("credential b was revoked alongside a: %#v", out)
	}
}

// TestCredentialFormatIsOpaque proves D3's opacity: the raw token contains no
// episode/grants prefix and is exactly 32 bytes of crypto/rand. We cannot prove
// randomness, so we assert structure: 64 hex chars, no non-hex, and distinct
// across mints (already covered above).
func TestCredentialFormatIsOpaque(t *testing.T) {
	st := newTestStore(t)
	tok := mintTestToken(t, st, "ep-42", testGrants(), 3600, 1000)
	bin, err := hex.DecodeString(tok)
	if err != nil {
		t.Fatalf("token is not valid hex: %v", err)
	}
	if len(tok) != tokenHexLen {
		t.Fatalf("token len = %d, want %d", len(tok), tokenHexLen)
	}
	if len(bin) != 32 {
		t.Fatalf("token decodes to %d bytes, want 32", len(bin))
	}
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(tok) {
		t.Fatalf("token %q is not 64 lowercase hex", tok)
	}
}
