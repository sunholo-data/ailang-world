package store

// Row 201 A2 sprint 2 M5 (design §6, CS7-CS9): a passkey consume binds the credential that signed
// and moves its counter forward, in the transaction that consumes the request, or nothing is written.

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"
)

// passkeyAt is a fresh schema-6 store, its path, and webauthn approver x enrolled at 900 with
// credential cred and counter stored.
func passkeyAt(t *testing.T, cred []byte, stored int64) (*Credentials, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	a, p := testPasskey("x", 1)
	a.EnrolledAt, p.CredentialID, p.SignCount = 900, cred, stored
	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
		t.Fatal(err)
	}
	return c, path
}

var credX = bytes.Repeat([]byte{1}, 32)

// passkeyConsume is okConsume(r, at, ...) for x's passkey assertion with counter count.
func passkeyConsume(r MintRequestRow, at int64, clocks, checks *int, cred []byte, count int64) RequestConsume {
	rc := okConsume(r, at, clocks, checks)
	rc.Approver = &ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: cred, SignCount: count}}
	rc.ApprovedBy = "webauthn:x"
	return rc
}

// passkeyState is x's counter and last_used_at.
func passkeyState(t *testing.T, c *Credentials) (int64, int64) {
	t.Helper()
	p, ok, err := c.Passkey(testCtx(t), "x")
	if err != nil || !ok {
		t.Fatalf("x's passkey: ok %v err %v", ok, err)
	}
	return p.SignCount, approverX(t, c).LastUsedAt
}

// CS7: C13 (the passkey is set exactly for webauthn), C14 (its bounds) and C15 (schema 6), each
// before the clock is read, leaving the request pending, no session, and x untouched.
func TestPasskeyConsumeCompleteness(t *testing.T) {
	for name, edit := range map[string]func(*RequestConsume){
		"C13-webauthn-without-passkey": func(rc *RequestConsume) { rc.Approver.Passkey = nil },
		"C13-ssh-with-passkey":         func(rc *RequestConsume) { rc.Approver.Format, rc.ApprovedBy = "ssh", "ssh:x" },
		"C14-empty-credential":         func(rc *RequestConsume) { rc.Approver.Passkey.CredentialID = nil },
		"C14-1024-byte-credential":     func(rc *RequestConsume) { rc.Approver.Passkey.CredentialID = make([]byte, 1024) },
		"C14-counter-2^32":             func(rc *RequestConsume) { rc.Approver.Passkey.SignCount = 1 << 32 },
		"C14-counter--1":               func(rc *RequestConsume) { rc.Approver.Passkey.SignCount = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			c, _ := passkeyAt(t, credX, 1)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := passkeyConsume(r, 1010, &clocks, &checks, credX, 2)
			edit(&rc)
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) || clocks != 0 {
				t.Fatalf("err %v after %d clock reads, want ErrIncompleteConsume before the clock", err, clocks)
			}
			assertPendingNoSession(t, c, r)
			if n, used := passkeyState(t, c); n != 1 || used != 0 {
				t.Fatalf("x after a refusal: counter %d, last_used_at %d", n, used)
			}
		})
	}
	t.Run("C15-v5-store", func(t *testing.T) {
		c := openV5(t)
		r := putTestRequest(t, c)
		var clocks, checks int
		err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(r, 1010, &clocks, &checks, credX, 2))
		if !isOlderPasskeys(err) || clocks != 0 {
			t.Fatalf("err %v after %d clock reads, want the passkeys refusal before the clock", err, clocks)
		}
		assertPendingNoSession(t, c, r)
	})
}

// CS8: the counter moves to the assertion's only forward (both zero: unchanged), with last_used_at,
// approved_by and the session; otherwise ErrSignCountRegression and nothing changes.
func TestPasskeyConsumeAdvancesTheCounter(t *testing.T) {
	for _, tc := range []struct {
		stored, count int64
		cred          []byte
		ok            bool
	}{
		{1, 2, credX, true}, {0, 0, credX, true}, {0, 7, credX, true}, {1, 1<<32 - 1, credX, true},
		{1, 2, []byte{7}, true}, {1, 2, bytes.Repeat([]byte{7}, 1023), true},
		{2, 2, credX, false}, {3, 2, credX, false}, {5, 0, credX, false}, {1<<32 - 1, 1<<32 - 1, credX, false},
	} {
		c, _ := passkeyAt(t, tc.cred, tc.stored)
		r := putTestRequest(t, c)
		var clocks, checks int
		err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(r, 1010, &clocks, &checks, tc.cred, tc.count))
		n, used := passkeyState(t, c)
		if !tc.ok {
			if !errors.Is(err, ErrSignCountRegression) || n != tc.stored || used != 0 {
				t.Fatalf("stored %d, new %d: err %v, counter %d, last_used_at %d; want the regression and nothing changed", tc.stored, tc.count, err, n, used)
			}
			assertPendingNoSession(t, c, r)
			continue
		}
		got, _, gerr := c.MintRequest(testCtx(t), r.RequestID)
		sessions, _ := c.ListSessions(testCtx(t), "", MaxListSessions)
		if err != nil || n != tc.count || used != 1010 || gerr != nil || got.ApprovedBy != "webauthn:x" || len(sessions) != 1 {
			t.Fatalf("stored %d, new %d (cred %d bytes): err %v, counter %d, last_used_at %d, approved_by %q, %d sessions",
				tc.stored, tc.count, len(tc.cred), err, n, used, got.ApprovedBy, len(sessions))
		}
	}
	// The same counter again, on a fresh request, is a regression: the first consume moved it.
	c, _ := passkeyAt(t, credX, 1)
	var clocks, checks int
	if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(putTestRequest(t, c), 1010, &clocks, &checks, credX, 2)); err != nil {
		t.Fatal(err)
	}
	r2 := testMintRequestRow(2, 1000)
	if err := c.PutMintRequest(testCtx(t), r2, 1000); err != nil {
		t.Fatal(err)
	}
	err := c.MintSessionForRequest(testCtx(t), consumeSession(2), passkeyConsume(r2, 1020, &clocks, &checks, credX, 2))
	if n, used := passkeyState(t, c); !errors.Is(err, ErrSignCountRegression) || n != 2 || used != 1010 {
		t.Fatalf("the replayed counter: err %v, counter %d, last_used_at %d; want the regression, 2, 1010", err, n, used)
	}
}

// CS9: the passkey row is judged in the consume's transaction. A commit made on a second connection
// after BeginTx (from the Clock hook, as A1b's CS3) that revokes x, rewrites its credential ID,
// deletes its passkey row or raises its counter refuses the consume and rolls back everything.
func TestPasskeyConsumeChecksTheRowInTheTransaction(t *testing.T) {
	for _, a := range []struct {
		name, sql string
		want      error
		counter   int64
	}{
		{"revoked", `UPDATE approvers SET revoked_at = 1005 WHERE name = 'x'`, ErrApproverRevoked, 1},
		{"credential-rewritten", `UPDATE webauthn_credentials SET credential_id = x'02' WHERE name = 'x'`, ErrPasskeyChanged, 1},
		{"passkey-deleted", `DELETE FROM webauthn_credentials WHERE name = 'x'`, ErrApproverUnknown, -1},
		{"counter-raised", `UPDATE webauthn_credentials SET sign_count = 9 WHERE name = 'x'`, ErrSignCountRegression, 9},
	} {
		t.Run(a.name, func(t *testing.T) {
			c, path := passkeyAt(t, credX, 1)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := passkeyConsume(r, 1010, &clocks, &checks, credX, 2)
			rc.Clock = func() int64 {
				if _, err := rawDB(t, path).Exec(a.sql); err != nil {
					t.Errorf("the second connection: %v", err)
				}
				return 1010
			}
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
				t.Fatalf("err %v, want %v", err, a.want)
			}
			assertPendingNoSession(t, c, r)
			if x := approverX(t, c); x.LastUsedAt != 0 {
				t.Fatalf("x's last_used_at %d after a refusal, want 0", x.LastUsedAt)
			}
			if p, ok, _ := c.Passkey(testCtx(t), "x"); a.counter >= 0 && (!ok || p.SignCount != a.counter) {
				t.Fatalf("x's counter %d (ok %v), want %d", p.SignCount, ok, a.counter)
			}
		})
	}
}
