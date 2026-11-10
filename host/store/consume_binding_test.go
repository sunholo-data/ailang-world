package store

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M4d1, M4d2): the consume binds the request's challenge and, on the signature path, its
// approver's active row, in the transaction that consumes the request and inserts the session.

import (
	"errors"
	"strings"
	"testing"
)

// putTestRequest stores testMintRequestRow(1, 1000) in c and returns it.
func putTestRequest(t *testing.T, c *Credentials) MintRequestRow {
	t.Helper()
	r := testMintRequestRow(1, 1000)
	if err := c.PutMintRequest(testCtx(t), r, 1000); err != nil {
		t.Fatal(err)
	}
	return r
}

// assertPendingNoSession fails unless r is still pending in c and c holds no session row.
func assertPendingNoSession(t *testing.T, c *Credentials, r MintRequestRow) {
	t.Helper()
	ctx := testCtx(t)
	if got, ok, err := c.MintRequest(ctx, r.RequestID); err != nil || !ok || got.Consumed {
		t.Fatalf("request %s after a refused consume: %+v ok %v err %v, want it pending", r.RequestID, got, ok, err)
	}
	if rows, err := c.ListSessions(ctx, "", MaxListSessions); err != nil || len(rows) != 0 {
		t.Fatalf("%d session rows after a refused consume (err %v), want 0", len(rows), err)
	}
}

// CS1 (M4c2): a consume whose challenge is not 64 lowercase hex is incomplete (C1), refused
// before the clock is read.
func TestMintSessionForRequestNeedsTheChallenge(t *testing.T) {
	for name, ch := range map[string]string{
		"empty":   "",
		"63-hex":  strings.Repeat("c", 63),
		"65-hex":  strings.Repeat("c", 65),
		"upper":   strings.Repeat("CD", 32),
		"not-hex": strings.Repeat("zz", 32),
	} {
		t.Run(name, func(t *testing.T) {
			c := openTestCredentials(t)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := okConsume(r, 1010, &clocks, &checks)
			rc.Challenge = ch
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) {
				t.Fatalf("err %v, want ErrIncompleteConsume", err)
			}
			if clocks != 0 {
				t.Fatalf("Clock read %d times before an incomplete-consume refusal, want 0", clocks)
			}
			assertPendingNoSession(t, c, r)
		})
	}
}

// CS5 (M4c2): the guarded UPDATE binds the row's challenge (C6), and the refusal names it (C7),
// after the digest: a consume with both wrong names the digest.
func TestConsumeBindsTheChallenge(t *testing.T) {
	other := strings.Repeat("ef", 32)
	for _, a := range []struct {
		name   string
		mutate func(*RequestConsume)
		want   error
	}{
		{"other-challenge", func(rc *RequestConsume) { rc.Challenge = other }, ErrRequestChallengeChanged},
		{"other-digest-and-challenge", func(rc *RequestConsume) { rc.Digest, rc.Challenge = other, other }, ErrRequestDigestChanged},
	} {
		t.Run(a.name, func(t *testing.T) {
			c := openTestCredentials(t)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := okConsume(r, 1010, &clocks, &checks)
			a.mutate(&rc)
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
				t.Fatalf("err %v, want %v", err, a.want)
			}
			assertPendingNoSession(t, c, r)
		})
	}
}
