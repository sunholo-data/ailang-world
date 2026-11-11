package store

// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
// M4c2, M4d1, M4d2): the consume binds the request's challenge and, on the signature path, its
// approver's active row, in the transaction that consumes the request and inserts the session.

import (
	"errors"
	"path/filepath"
	"reflect"
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

// ---- M4d1 ----

// credentialsAt is a fresh schema-5 credentials store, its path, and approver x
// (testApprover("x", 900): format ssh, key "key-x") enrolled in it.
func credentialsAt(t *testing.T) (*Credentials, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	if err := c.EnrolApprover(testCtx(t), testApprover("x", 900)); err != nil {
		t.Fatal(err)
	}
	return c, path
}

// signedConsume is okConsume(r, at, ...) for a signature approval by x: Approver x, approved_by ssh:x.
func signedConsume(r MintRequestRow, at int64, clocks, checks *int) RequestConsume {
	rc := okConsume(r, at, clocks, checks)
	rc.Approver = &ApproverUse{Name: "x", Format: "ssh", PublicKey: []byte("key-x")}
	rc.ApprovedBy = "ssh:x"
	return rc
}

// approverX is approver x's row in c.
func approverX(t *testing.T, c *Credentials) ApproverRow {
	t.Helper()
	row, ok, err := c.Approver(testCtx(t), "x")
	if err != nil || !ok {
		t.Fatalf("approver x: ok %v err %v", ok, err)
	}
	return row
}

// CS2 (M4d1): approved_by is tty exactly when no approver is bound (C3) and <format>:<name> of the
// bound approver otherwise (C4); an approver needs schema 5 (C5). Each refusal comes before the
// clock is read and leaves the request pending and no session.
func TestApprovedByMatchesTheApprover(t *testing.T) {
	for _, a := range []struct {
		name   string
		mutate func(*RequestConsume)
	}{
		{"no-approver-ssh", func(rc *RequestConsume) { rc.Approver = nil }},
		{"approver-tty", func(rc *RequestConsume) { rc.ApprovedBy = "tty" }},
		{"approver-other-name", func(rc *RequestConsume) { rc.ApprovedBy = "ssh:y" }},
		{"approver-other-format", func(rc *RequestConsume) { rc.ApprovedBy = "webauthn:x" }},
		{"approver-format-not-ssh", func(rc *RequestConsume) { rc.Approver.Format = "webauthn" }},
	} {
		t.Run(a.name, func(t *testing.T) {
			c, _ := credentialsAt(t)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := signedConsume(r, 1010, &clocks, &checks)
			a.mutate(&rc)
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrInconsistentApproval) {
				t.Fatalf("err %v, want ErrInconsistentApproval", err)
			}
			if clocks != 0 {
				t.Fatalf("Clock read %d times, want 0", clocks)
			}
			assertPendingNoSession(t, c, r)
			if got := approverX(t, c); got.LastUsedAt != 0 {
				t.Fatalf("x's last_used_at %d after a refusal, want 0", got.LastUsedAt)
			}
		})
	}
	t.Run("v4-store", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), CredentialsFileName)
		makeV4Credentials(t, path, nil)
		c, err := OpenCredentials(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = c.Close() }()
		r := putTestRequest(t, c)
		var clocks, checks int
		err = c.MintSessionForRequest(testCtx(t), consumeSession(1), signedConsume(r, 1010, &clocks, &checks))
		var older *OlderCredentialsStoreError
		if !errors.As(err, &older) || older.Need != 5 || older.Feature != "approvers" {
			t.Fatalf("err %v, want *OlderCredentialsStoreError{Need: 5, Feature: approvers}", err)
		}
		if clocks != 0 {
			t.Fatalf("Clock read %d times on a schema-4 store, want 0", clocks)
		}
		assertPendingNoSession(t, c, r)
	})
}

// requestConsumeFields is the disposition of every RequestConsume field and of its element type
// ApproverUse's fields (RC3). "required" fields are zeroed one at a time below.
var requestConsumeFields = map[string]string{
	"RequestConsume.ID":         "required",
	"RequestConsume.Digest":     "required",
	"RequestConsume.Challenge":  "required (C1)",
	"RequestConsume.ApprovedBy": "derived and checked: tty with no approver, <format>:<name> with one (C3, C4; TestApprovedByMatchesTheApprover)",
	"RequestConsume.Approver":   "optional: nil is the typed y; set, it is bound by C2, C5 and C8",
	"RequestConsume.Clock":      "required",
	"RequestConsume.Check":      "required",
	"ApproverUse.Name":          "required (C2)",
	"ApproverUse.Format":        "required (C2)",
	"ApproverUse.PublicKey":     "required (C2)",
	"ApproverUse.Passkey":       "optional with a meaning: set exactly for a webauthn approver (C13), bound by C14-C18",
	"PasskeyUse.CredentialID":   "required (C14): 1 to 1023 bytes",
	"PasskeyUse.SignCount":      "bounded (C14): 0 to 2^32-1; 0 is a counter (TestPasskeyConsumeCompleteness)",
}

// RC3 (M4d1): every RequestConsume and ApproverUse field has a disposition, and zeroing each
// required one in an otherwise valid signed consume is ErrIncompleteConsume, before the clock.
func TestRequestConsumeCompleteness(t *testing.T) {
	n := 0
	for _, typ := range []reflect.Type{reflect.TypeOf(RequestConsume{}), reflect.TypeOf(ApproverUse{}), reflect.TypeOf(PasskeyUse{})} {
		for i := 0; i < typ.NumField(); i++ {
			if key := typ.Name() + "." + typ.Field(i).Name; requestConsumeFields[key] == "" {
				t.Errorf("%s has no disposition in requestConsumeFields", key)
			}
		}
		n += typ.NumField()
	}
	if n != len(requestConsumeFields) {
		t.Fatalf("RequestConsume, ApproverUse and PasskeyUse have %d fields, requestConsumeFields %d", n, len(requestConsumeFields))
	}
	zero := map[string]func(*RequestConsume){
		"RequestConsume.ID":        func(rc *RequestConsume) { rc.ID = "" },
		"RequestConsume.Digest":    func(rc *RequestConsume) { rc.Digest = "" },
		"RequestConsume.Challenge": func(rc *RequestConsume) { rc.Challenge = "" },
		"RequestConsume.Clock":     func(rc *RequestConsume) { rc.Clock = nil },
		"RequestConsume.Check":     func(rc *RequestConsume) { rc.Check = nil },
		"ApproverUse.Name":         func(rc *RequestConsume) { rc.Approver.Name, rc.ApprovedBy = "", "ssh:" },
		"ApproverUse.Format":       func(rc *RequestConsume) { rc.Approver.Format, rc.ApprovedBy = "", ":x" },
		"ApproverUse.PublicKey":    func(rc *RequestConsume) { rc.Approver.PublicKey = nil },
		"PasskeyUse.CredentialID": func(rc *RequestConsume) {
			rc.Approver.Format, rc.ApprovedBy, rc.Approver.Passkey = "webauthn", "webauthn:x", &PasskeyUse{SignCount: 2}
		},
	}
	for key, d := range requestConsumeFields {
		if _, ok := zero[key]; strings.HasPrefix(d, "required") != ok {
			t.Errorf("%s: disposition %q, zeroing arm present %v", key, d, ok)
		}
	}
	for key, z := range zero {
		t.Run(key, func(t *testing.T) {
			c, _ := credentialsAt(t)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := signedConsume(r, 1010, &clocks, &checks)
			z(&rc)
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) {
				t.Fatalf("err %v, want ErrIncompleteConsume", err)
			}
			if clocks != 0 {
				t.Fatalf("Clock read %d times, want 0", clocks)
			}
			assertPendingNoSession(t, c, r)
		})
	}
}

// ---- M4d2 ----

// CS3 (M4d2): a signed consume marks its approver used at the in-transaction reading (C8). It is
// refused, the request pending, no session and last_used_at unchanged, when the approver is
// revoked by a commit made after BeginTx (from inside the Clock hook: the transaction is deferred,
// so it holds no lock until its first statement, design V36/V37), unknown (C9), revoked (C10),
// holding another key or format (C11), or enrolled after the reading (C12).
func TestSignedConsumeChecksTheApproverInTheTransaction(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		c, _ := credentialsAt(t)
		r := putTestRequest(t, c)
		var clocks, checks int
		if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), signedConsume(r, 1010, &clocks, &checks)); err != nil {
			t.Fatal(err)
		}
		got, ok, err := c.MintRequest(testCtx(t), r.RequestID)
		if err != nil || !ok || !got.Consumed || got.ApprovedBy != "ssh:x" {
			t.Fatalf("request after the signed consume: %+v ok %v err %v, want consumed by ssh:x", got, ok, err)
		}
		if x := approverX(t, c); x.LastUsedAt != 1010 {
			t.Fatalf("x's last_used_at %d, want 1010 (the in-transaction reading)", x.LastUsedAt)
		}
		if rows, err := c.ListSessions(testCtx(t), "", MaxListSessions); err != nil || len(rows) != 1 {
			t.Fatalf("%d session rows (err %v), want 1", len(rows), err)
		}
	})
	for _, a := range []struct {
		name  string
		setup func(t *testing.T, c *Credentials, path string, rc *RequestConsume)
		want  error
	}{
		{"revoked-in-the-window", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
			rc.Clock = func() int64 {
				if _, err := rawDB(t, path).Exec(`UPDATE approvers SET revoked_at = 1005 WHERE name = 'x'`); err != nil {
					t.Errorf("revoke on a second connection: %v", err)
				}
				return 1010
			}
		}, ErrApproverRevoked},
		{"unknown", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
			rc.Approver.Name, rc.ApprovedBy = "y", "ssh:y"
		}, ErrApproverUnknown},
		{"other-key", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
			rc.Approver.PublicKey = []byte("key-y")
		}, ErrApproverKeyChanged},
		{"other-format", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
			rc.Approver.Format, rc.ApprovedBy = "webauthn", "webauthn:x"
			rc.Approver.Passkey = &PasskeyUse{CredentialID: make([]byte, 32), SignCount: 1} // past C13-C15 (row 201 A2 §2.10)
		}, ErrApproverKeyChanged},
		{"enrolled-after-the-reading", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
			if _, err := rawDB(t, path).Exec(`UPDATE approvers SET enrolled_at = 1500 WHERE name = 'x'`); err != nil {
				t.Fatal(err)
			}
		}, ErrApproverClockBehind},
	} {
		t.Run(a.name, func(t *testing.T) {
			c, path := credentialsAt(t)
			r := putTestRequest(t, c)
			var clocks, checks int
			rc := signedConsume(r, 1010, &clocks, &checks)
			a.setup(t, c, path, &rc)
			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
				t.Fatalf("err %v, want %v", err, a.want)
			}
			assertPendingNoSession(t, c, r)
			if x := approverX(t, c); x.LastUsedAt != 0 {
				t.Fatalf("x's last_used_at %d after a refusal, want 0", x.LastUsedAt)
			}
		})
	}
}

// CS4 (M4d2): a typed consume (no approver, approved_by tty) mints on a schema-5 store that holds an
// approver, and touches no approver row.
func TestTypedConsumeTouchesNoApprover(t *testing.T) {
	c, _ := credentialsAt(t)
	r := putTestRequest(t, c)
	var clocks, checks int
	if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), okConsume(r, 1010, &clocks, &checks)); err != nil {
		t.Fatal(err)
	}
	if x := approverX(t, c); x.LastUsedAt != 0 {
		t.Fatalf("x's last_used_at %d after a typed consume, want 0", x.LastUsedAt)
	}
}
