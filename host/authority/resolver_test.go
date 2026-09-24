package authority

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// newTestStore opens an in-memory store with the full (v3) schema.
func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open in-memory store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// mintTestToken mints a credential for a test and returns just the raw token.
func mintTestToken(t *testing.T, st *store.Store, episode string, grants []broker.Capability, ttl, now int64) string {
	t.Helper()
	tok, _, _, err := Mint(context.Background(), st, episode, grants, ttl, now, io.Discard)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	return tok
}

func testGrants() []broker.Capability {
	return []broker.Capability{{Effect: "fs.read", Scope: "/tmp/b", Budget: 1}}
}

var hex64 = strings.Repeat("a", 64)

// TestResolve_TypedDenialTaxonomy drives all four denial kinds and the success
// path through one resolver, asserting that absent / malformed / unknown /
// expired are all DISTINCT kinds (mutation M2's sole killer: collapsing unknown
// into absent must red this test).
func TestResolve_TypedDenialTaxonomy(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	ttl := int64(60)
	tok := mintTestToken(t, st, "ep-42", testGrants(), ttl, now)
	res := New(st)

	// Absent: empty header.
	absent := res.Resolve("", now)
	if absent.Success != nil || absent.Denied == nil || *absent.Denied != DenialAbsent {
		t.Fatalf("empty header outcome = %#v, want Denied=DenialAbsent", absent)
	}
	// Malformed: present header but not "Bearer <64-hex>".
	malformed := res.Resolve("Key abc", now)
	if malformed.Denied == nil || *malformed.Denied != DenialMalformed {
		t.Fatalf("\"Key abc\" outcome = %#v, want Denied=DenialMalformed", malformed)
	}
	// Unknown: valid-shaped token, no mapping row.
	unknown := res.Resolve("Bearer "+hex64, now)
	if unknown.Denied == nil || *unknown.Denied != DenialUnknown {
		t.Fatalf("unknown token outcome = %#v, want Denied=DenialUnknown", unknown)
	}
	// Expired: real row, but now past expires_at.
	expired := res.Resolve("Bearer "+tok, now+ttl+1)
	if expired.Denied == nil || *expired.Denied != DenialExpired {
		t.Fatalf("expired token outcome = %#v, want Denied=DenialExpired", expired)
	}
	// Success: valid, unexpired.
	ok := res.Resolve("Bearer "+tok, now+1)
	if ok.Success == nil || ok.Denied != nil {
		t.Fatalf("valid token outcome = %#v, want Success", ok)
	}

	// Every denial kind is distinct and, together with Success, covers the
	// outcome space of interest without overlapping.
	kinds := []DenialKind{*absent.Denied, *malformed.Denied, *unknown.Denied, *expired.Denied}
	seen := make(map[DenialKind]bool)
	for _, k := range kinds {
		if seen[k] {
			t.Fatalf("denial kind %d (%s) repeated across the taxonomy — kinds must be distinct", uint8(k), k)
		}
		seen[k] = true
	}
	if len(seen) != 4 {
		t.Fatalf("expected exactly 4 distinct denial kinds, got %d", len(seen))
	}
}

// TestResolve_MalformedHeader covers the AC-M2-2 malformed shapes: a non-Bearer
// scheme, a bare "Bearer", and a token that is not 64 hex — all DenialMalformed
// (never 401-unknown).
func TestResolve_MalformedHeader(t *testing.T) {
	st := newTestStore(t)
	res := New(st)
	for _, header := range []string{"Key abc", "Bearer", "Bearer z", "Bearer " + strings.Repeat("z", 64)} {
		out := res.Resolve(header, 1000)
		if out.Success != nil || out.Denied == nil || *out.Denied != DenialMalformed {
			t.Fatalf("header %q outcome = %#v, want Denied=DenialMalformed", header, out)
		}
	}
}

// TestResolve_SuccessBinding proves AC-M1-3: resolve returns episode+grants+expiry
// and never the raw token (SessionBinding has no credential field).
func TestResolve_SuccessBinding(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	ttl := int64(3600)
	episode := "ep-42"
	grants := []broker.Capability{{Effect: "fs.read", Scope: "/tmp/b", Budget: 1}}
	tok := mintTestToken(t, st, episode, grants, ttl, now)

	b := New(st).Resolve("Bearer "+tok, now+5)
	if b.Success == nil {
		t.Fatal("expected success")
	}
	s := b.Success
	if s.EpisodeID != episode {
		t.Fatalf("EpisodeID = %q, want %q", s.EpisodeID, episode)
	}
	if len(s.Caps) != 1 || s.Caps[0].Effect != "fs.read" || s.Caps[0].Scope != "/tmp/b" || s.Caps[0].Budget != 1 {
		t.Fatalf("Caps = %+v, want the minted fs.read grant", s.Caps)
	}
	if s.ExpiresAt != now+ttl {
		t.Fatalf("ExpiresAt = %d, want %d", s.ExpiresAt, now+ttl)
	}
	if s.CreatedAt != now {
		t.Fatalf("CreatedAt = %d, want %d", s.CreatedAt, now)
	}
}

// TestResolve_ExpiryEnforced proves AC-M2-4 / mutation M3's sole killer: a row
// with expires_at in the past resolves to DenialExpired, distinct from unknown.
func TestResolve_ExpiryEnforced(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	ttl := int64(1) // --ttl 1
	tok := mintTestToken(t, st, "ep-42", testGrants(), ttl, now)
	res := New(st)

	// Before expiry: success.
	if out := res.Resolve("Bearer "+tok, now); out.Success == nil {
		t.Fatalf("resolve before expiry = %#v, want success", out)
	}
	// After expiry (now advanced past expires_at): DenialExpired.
	expired := res.Resolve("Bearer "+tok, now+ttl+1)
	if expired.Success != nil || expired.Denied == nil || *expired.Denied != DenialExpired {
		t.Fatalf("resolve after expiry = %#v, want Denied=DenialExpired", expired)
	}
}

// TestResolve_UnknownVsExpired distinguishes unknown from expired on the same
// store: a random valid-shaped token is Unknown (never Expired).
func TestResolve_UnknownVsExpired(t *testing.T) {
	st := newTestStore(t)
	res := New(st)
	out := res.Resolve("Bearer "+hex64, 1000)
	if out.Denied == nil || *out.Denied != DenialUnknown {
		t.Fatalf("unknown-shape token = %#v, want Denied=DenialUnknown", out)
	}
}

// TestNoAlternateHeader is mutation M6's sole killer: the resolver only ever
// reads the Authorization header it is passed. An empty Authorization (the only
// signal an X-World-Session-only request can present to the resolver, which has
// no alternate header concept by construction) yields DenialAbsent — it never
// attempts any alternate-header or fallback lookup. The HTTP middleware half of
// AC-M2-6 wires in M2.
func TestNoAlternateHeader(t *testing.T) {
	st := newTestStore(t)
	tok := mintTestToken(t, st, "ep-42", testGrants(), 3600, 1000)
	res := New(st)
	// The resolver receives "" — an X-World-Session-only request contributes no
	// Authorization value, so it must be absent, never resolved.
	out := res.Resolve("", 1000)
	if out.Success != nil || out.Denied == nil || *out.Denied != DenialAbsent {
		t.Fatalf("empty header (alternate-header-only request) = %#v, want Denied=DenialAbsent", out)
	}
	// And the valid token still resolves only via Authorization, as a control.
	if out := res.Resolve("Bearer "+tok, 1000); out.Success == nil {
		t.Fatalf("valid token control did not resolve: %#v", out)
	}
}

// countingStore is the bounded-lookup seam (AC-M3-2): it wraps a real
// CredentialStore and counts how many resolve queries the resolver issues.
type countingStore struct {
	credentialStore CredentialStore
	resolveCalls    int
}

func (c *countingStore) ResolveSession(ctx context.Context, credentialID string) (store.SessionRow, bool, error) {
	c.resolveCalls++
	return c.credentialStore.ResolveSession(ctx, credentialID)
}

// TestResolve_SingleBoundedLookup asserts the resolve path issues EXACTLY one
// indexed query for a present token and ZERO queries for an absent header or a
// malformed header (bounded by construction, D3/D5 / AC-M3-2). Each
// CredentialStore.ResolveSession is a single indexed WHERE credential_id = ?
// PK lookup inside the store, so the seam call count equates to the SQL query
// count.
func TestResolve_SingleBoundedLookup(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	tok := mintTestToken(t, st, "ep-42", testGrants(), 3600, now)
	counter := &countingStore{credentialStore: st}
	res := New(counter)

	// A present token resolves with exactly one query.
	out := res.Resolve("Bearer "+tok, now)
	if out.Success == nil {
		t.Fatalf("expected success, got %#v", out)
	}
	if counter.resolveCalls != 1 {
		t.Fatalf("resolve query count = %d, want exactly 1", counter.resolveCalls)
	}

	// Absent header never touches the store (zero queries).
	counter.resolveCalls = 0
	if out := res.Resolve("", now); out.Denied == nil || *out.Denied != DenialAbsent {
		t.Fatalf("absent header outcome = %#v", out)
	}
	if counter.resolveCalls != 0 {
		t.Fatalf("absent header issued %d queries, want 0", counter.resolveCalls)
	}

	// Malformed header never touches the store (zero queries).
	counter.resolveCalls = 0
	if out := res.Resolve("Key abc", now); out.Denied == nil || *out.Denied != DenialMalformed {
		t.Fatalf("malformed outcome = %#v", out)
	}
	if counter.resolveCalls != 0 {
		t.Fatalf("malformed header issued %d queries, want 0", counter.resolveCalls)
	}
}

// TestAuthority_NoDirectTokenCompare is mutation M1's sole killer (AC-M3-1). It
// source-scans every non-test Go file in host/authority and FAILS if any line
// performs a direct equality/contains comparison of token or credential-hash
// material. The resolver's equality is the SQLite indexed PK lookup inside the
// store, never a Go byte-for-byte compare of secret material — unknown-vs-known
// tokens traverse the identical single-query shape.
func TestAuthority_NoDirectTokenCompare(t *testing.T) {
	tokenish := `(?:presented|stored|token|secret|digest|credential|hash)`
	// A tokenish identifier on BOTH sides of a ==/!= — the shape a direct
	// secret comparison takes (e.g. `if presentedHash == storedHash`).
	bothSides := regexp.MustCompile(`(?i)\b` + tokenish + `\w*\s*(?:==|!=)\s*(?:` + tokenish + `)\w*`)
	// An equality/contains helper applied to secret material.
	helper := regexp.MustCompile(`(?i)(?:bytes\.Equal|EqualFold|Contains)\([^)]*\b` + tokenish + `[^)]*\)`)

	files, err := filepath.Glob(filepath.Join(".", "*.go"))
	if err != nil {
		t.Fatalf("glob authority files: %v", err)
	}
	checked := 0
	anyBad := false
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		checked++
		for i, line := range strings.Split(string(data), "\n") {
			if bothSides.MatchString(line) || helper.MatchString(line) {
				t.Errorf("%s:%d: direct token/credential-hash comparison: %q", f, i+1, strings.TrimSpace(line))
				anyBad = true
			}
		}
	}
	if checked == 0 {
		t.Fatal("source scan checked zero non-test authority files — non-vacuous gate vacuous")
	}
	if anyBad {
		t.Fatal("authority resolve path must not compare token/hash material directly; equality is the indexed store lookup")
	}
}
