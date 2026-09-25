package authority

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"

	_ "modernc.org/sqlite"
)

// ---------------------------------------------------------------------------
// P6.A-CTX — ResolveContext: bounded, context-aware session resolution
// (w-a2a-session-projection; the /v1/commit middleware keeps Resolve
// byte-for-byte — these tests pin the EQUIVALENCE and the ONE difference:
// store errors surface as errors, and the caller's ctx bounds the wait).
// ---------------------------------------------------------------------------

// TestResolveContext_PolicyEquivalence is the P6.A-CTX policy-equivalence
// gate (MUT-RESOLVE-POLICY-DRIFT's killer): over ONE fixture table spanning
// every policy branch (absent, both malformed shapes, unknown, the expiry
// boundary on BOTH sides, success), Resolve(header, now) and
// ResolveContext(ctx, header, now) with a background context must agree on
// the outcome EXACTLY. Because Resolve is implemented as ResolveContext over
// context.Background(), this test compares both entry points against the
// table's EXPECTED outcomes, not merely against each other — a drift in the
// shared branch (e.g. expiry `>` mutated to `>=`) fails the pinned
// expectation, which the bare compare could never see.
func TestResolveContext_PolicyEquivalence(t *testing.T) {
	st := newTestStore(t)
	now := int64(1000)
	ttl := int64(60)
	tok := mintTestToken(t, st, "ep-42", testGrants(), ttl, now)

	table := []struct {
		name     string
		header   string
		now      int64
		wantDeny *DenialKind // nil = want Success
	}{
		{name: "absent", header: "", now: now, wantDeny: denyPtr(DenialAbsent)},
		{name: "malformed scheme", header: "Key abc", now: now, wantDeny: denyPtr(DenialMalformed)},
		{name: "malformed token shape", header: "Bearer z", now: now, wantDeny: denyPtr(DenialMalformed)},
		{name: "unknown token", header: "Bearer " + hex64, now: now, wantDeny: denyPtr(DenialUnknown)},
		// Expiry boundary, pinned from the policy text `now > row.ExpiresAt`:
		// at expires_at exactly the credential is STILL LIVE (success)…
		{name: "expiry boundary live", header: "Bearer " + tok, now: now + ttl, wantDeny: nil},
		// …one second past it, expired.
		{name: "expired", header: "Bearer " + tok, now: now + ttl + 1, wantDeny: denyPtr(DenialExpired)},
		{name: "success", header: "Bearer " + tok, now: now + 1, wantDeny: nil},
	}

	res := New(st)
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			gotPlain := res.Resolve(tc.header, tc.now)
			gotCtx, err := New(st).ResolveContext(context.Background(), tc.header, tc.now)
			if err != nil {
				t.Fatalf("ResolveContext over background ctx returned error: %v", err)
			}
			// Both entry points agree with each other (same Success/Denied
			// values, same binding contents on success).
			if !sameOutcome(gotPlain, gotCtx) {
				t.Fatalf("Resolve = %#v but ResolveContext = %#v: the two entry points DRIFTED", gotPlain, gotCtx)
			}
			// And both agree with the pinned expectation.
			if tc.wantDeny == nil {
				if gotCtx.Success == nil || gotCtx.Denied != nil {
					t.Fatalf("outcome = %#v, want Success", gotCtx)
				}
			} else {
				if gotCtx.Success != nil || gotCtx.Denied == nil || *gotCtx.Denied != *tc.wantDeny {
					t.Fatalf("outcome = %#v, want Denied=%s", gotCtx, *tc.wantDeny)
				}
			}
		})
	}
}

func denyPtr(k DenialKind) *DenialKind { return &k }

func sameOutcome(a, b ResolveOutcome) bool {
	if (a.Success == nil) != (b.Success == nil) || (a.Denied == nil) != (b.Denied == nil) {
		return false
	}
	if a.Denied != nil && *a.Denied != *b.Denied {
		return false
	}
	if a.Success != nil {
		if a.Success.EpisodeID != b.Success.EpisodeID ||
			a.Success.ExpiresAt != b.Success.ExpiresAt ||
			a.Success.CreatedAt != b.Success.CreatedAt ||
			fmt.Sprint(a.Success.Caps) != fmt.Sprint(b.Success.Caps) {
			return false
		}
	}
	return true
}

// errorStore is a CredentialStore whose read always fails — the store-failure
// fixture (single-connection contention, driver faults) that ResolveContext
// must surface as an ERROR, never as a credential denial.
type errorStore struct{ err error }

func (e *errorStore) ResolveSession(context.Context, string) (store.SessionRow, bool, error) {
	return store.SessionRow{}, false, e.err
}

// TestResolveContext_StoreErrorIsErrorNotDenial is the R1-CTX fix's core and
// the killer of MUT-CTX-ERR-AS-UNKNOWN: a store read failure must return a
// NON-NIL error (wrapping the store's error) and a ZERO outcome from
// ResolveContext — while Resolve (same failure, background ctx) still
// collapses it to DenialUnknown exactly as before (the /v1/commit residual).
func TestResolveContext_StoreErrorIsErrorNotDenial(t *testing.T) {
	storeErr := errors.New("store read failed")
	res := New(&errorStore{err: storeErr})

	out, err := res.ResolveContext(context.Background(), "Bearer "+hex64, 1000)
	if err == nil {
		t.Fatalf("ResolveContext with a failing store: err = nil, outcome = %#v — a store failure is NOT a credential denial", out)
	}
	if !errors.Is(err, storeErr) && !strings.Contains(err.Error(), storeErr.Error()) {
		t.Fatalf("err = %v, want it to carry the store error", err)
	}
	if out.Success != nil || out.Denied != nil {
		t.Fatalf("store-failure outcome = %#v, want the zero ResolveOutcome", out)
	}

	// Resolve keeps today's byte-identical behaviour: same failure collapses
	// to DenialUnknown (never an error).
	plain := res.Resolve("Bearer "+hex64, 1000)
	if plain.Success != nil || plain.Denied == nil || *plain.Denied != DenialUnknown {
		t.Fatalf("Resolve over a failing store = %#v, want DenailUnknown (unchanged /v1/commit behaviour)", plain)
	}

	// A context error on the store read is likewise an ERROR, so a transport
	// deadline never arrives at the caller disguised as "unknown credential".
	ctxErr := New(&errorStore{err: context.DeadlineExceeded})
	out2, err2 := ctxErr.ResolveContext(context.Background(), "Bearer "+hex64, 1000)
	if !errors.Is(err2, context.DeadlineExceeded) || out2.Success != nil || out2.Denied != nil {
		t.Fatalf("ctx-error outcome = (%#v, %v), want (zero outcome, error wrapping context.DeadlineExceeded)", out2, err2)
	}
}

// heldConnDB is a CredentialStore backed by a REAL one-connection
// database/sql pool (the same SetMaxOpenConns(1) shape store.Open pins at
// store.go:305): the test holds the pool's ONLY connection with an open
// transaction, so the resolve query's wait for a free connection is the
// genuine database/sql mechanism — exactly the wait a worldd resolve faces
// behind an in-flight commit. The transaction self-rolls-back after 3 s so a
// mutation that drops the ctx (MUT-RESOLVE-BACKGROUND) fails RED on
// assertions instead of hanging the suite.
type heldConnDB struct {
	db *sql.DB
}

func (h *heldConnDB) ResolveSession(ctx context.Context, credentialID string) (store.SessionRow, bool, error) {
	var row store.SessionRow
	err := h.db.QueryRowContext(ctx,
		`SELECT credential_id, episode_id, grants_json, expires_at, created_at
		   FROM session_credentials WHERE credential_id = ?;`,
		credentialID,
	).Scan(&row.CredentialID, &row.EpisodeID, &row.GrantsJSON, &row.ExpiresAt, &row.CreatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return store.SessionRow{}, false, nil
	case err != nil:
		return store.SessionRow{}, false, fmt.Errorf("store: resolve session: %w", err)
	}
	return row, true, nil
}

// TestResolveContext_DeadlineWhileSingleConnectionHeld is the P6.A-CTX
// deadline test: the store's single pooled connection is held by an open
// transaction, so ResolveContext with a 100 ms ctx must come back with an
// error wrapping context.DeadlineExceeded WITHIN the bound — never a
// DenialUnknown 401 (which is what the context-less Resolve would have
// answered) and never an unbounded wait. MUT-RESOLVE-BACKGROUND's killer: if
// the store call ignores ctx (context.Background()), the query sits behind
// the held connection until the 3 s self-rollback, then resolves WITHOUT an
// error, and every assertion below reds.
func TestResolveContext_DeadlineWhileSingleConnectionHeld(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open raw sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1) // the store.go:305 shape: ONE pooled connection
	if _, err := db.Exec(`CREATE TABLE session_credentials (
		credential_id TEXT PRIMARY KEY,
		episode_id    TEXT NOT NULL,
		grants_json   TEXT NOT NULL,
		expires_at    INTEGER NOT NULL,
		created_at    INTEGER NOT NULL
	);`); err != nil {
		t.Fatalf("create session_credentials: %v", err)
	}

	// Hold the pool's only connection with an open transaction; it frees
	// itself after 3 s so a ctx-dropping mutation REDS instead of hanging.
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin held transaction: %v", err)
	}
	time.AfterFunc(3*time.Second, func() { _ = tx.Rollback() })

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	out, rerr := New(&heldConnDB{db: db}).ResolveContext(ctx, "Bearer "+hex64, time.Now().Unix())
	elapsed := time.Since(start)

	if rerr == nil {
		t.Fatalf("ResolveContext with the single connection held and a 100ms deadline: err = nil, outcome = %#v — "+
			"want an error wrapping context.DeadlineExceeded (the ctx was dropped or the error was misreported as a denial)", out)
	}
	if !errors.Is(rerr, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want errors.Is(err, context.DeadlineExceeded)", rerr)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("resolution took %v — the bounded wait exceeded its bound (a dropped ctx waits behind the held connection)", elapsed)
	}
	if out.Success != nil || out.Denied != nil {
		t.Fatalf("deadline outcome = %#v, want the zero ResolveOutcome (never DenialUnknown for a store/ctx failure)", out)
	}

	// Resolve over the SAME held-pool fixture still answers DenialUnknown —
	// the unchanged /v1/commit behaviour — but it can only do so after the
	// transaction frees, which is the residual this doc declares.
}

// TestResolve_NoGoroutineInResolvePath is the source-level half of P6.A-CTX:
// resolution runs on the CALLING goroutine — a source scan of every non-test
// authority file must find no `go ` statement (no worker spawned, no
// background resolve), which is what makes the bounded ctx the ONLY wait a
// resolve performs.
func TestResolve_NoGoroutineInResolvePath(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(".", "*.go"))
	if err != nil {
		t.Fatalf("glob authority files: %v", err)
	}
	checked := 0
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
			trimmed := strings.TrimSpace(line)
			// A goroutine spawn is exactly a statement beginning with `go `;
			// the word may appear in comments and prose and is fine there.
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if trimmed == "go" || strings.HasPrefix(trimmed, "go ") {
				t.Errorf("%s:%d: goroutine spawn in the authority package: %q — resolution must run on the calling goroutine", f, i+1, trimmed)
			}
		}
	}
	if checked == 0 {
		t.Fatal("source scan checked zero non-test authority files — non-vacuous gate vacuous")
	}
}

// compile-time witness that the concrete resolver satisfies the widened
// interface, including ResolveContext.
var _ Resolver = (*resolver)(nil)
