package store

// Row 201 A2 sprint 2 M4b (design §6, PS2, PS5): schema 6's passkey rows, enrolled only with their
// approver, and never twice.

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// testPasskey is a valid new webauthn approver and passkey: name, COSE key "cose-"+name, enrolled
// at 1000; credential ID 32 bytes of id, rp localhost, ES256, counter 1.
func testPasskey(name string, id byte) (ApproverRow, PasskeyRow) {
	return ApproverRow{Name: name, Format: "webauthn", PublicKey: []byte("cose-" + name), EnrolledAt: 1000},
		PasskeyRow{Name: name, CredentialID: bytes.Repeat([]byte{id}, 32), RPID: "localhost", Origin: "http://localhost:7644",
			Alg: -7, AAGUID: make([]byte, 16), SignCount: 1}
}

// rowCounts is the number of approvers and webauthn_credentials rows in c's store.
func rowCounts(t *testing.T, c *Credentials) string {
	t.Helper()
	var a, p int
	if err := c.db.QueryRow(`SELECT (SELECT count(*) FROM approvers), (SELECT count(*) FROM webauthn_credentials)`).Scan(&a, &p); err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%d approvers, %d passkeys", a, p)
}

// openV5 opens a schema-5 store (not upgraded).
func openV5(t *testing.T) *Credentials {
	t.Helper()
	path := filepath.Join(t.TempDir(), CredentialsFileName)
	makeV5Credentials(t, path, nil)
	c, err := OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// isOlderPasskeys reports whether err is the passkeys refusal of a schema-5 store.
func isOlderPasskeys(err error) bool {
	var oe *OlderCredentialsStoreError
	return errors.As(err, &oe) && oe.Version == 5 && oe.Need == 6 && oe.Feature == "passkeys" && strings.Contains(err.Error(), "passkeys need schema 6")
}

// PS2: every EnrolPasskey refusal returns its typed error and writes nothing to either table.
func TestEnrolPasskeyRefusals(t *testing.T) {
	ctx := testCtx(t)
	a, p := testPasskey("mark", 1)
	c5 := openV5(t)
	if err := c5.EnrolPasskey(ctx, a, p); !isOlderPasskeys(err) {
		t.Fatalf("v5 enrol: %v, want the passkeys refusal", err)
	}
	if n := approverCount(t, c5); n != 0 {
		t.Fatalf("v5 enrol wrote %d approvers", n)
	}

	c := openTestCredentials(t)
	bad := map[string]func(a *ApproverRow, p *PasskeyRow){
		"S1-format-ssh":    func(a *ApproverRow, p *PasskeyRow) { a.Format = "ssh" },
		"S1-format-empty":  func(a *ApproverRow, p *PasskeyRow) { a.Format = "" },
		"S2-names-differ":  func(a *ApproverRow, p *PasskeyRow) { p.Name = "other" },
		"S3-empty-name":    func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "", "" },
		"S3-nil-key":       func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = nil },
		"S4-enrolled-0":    func(a *ApproverRow, p *PasskeyRow) { a.EnrolledAt = 0 },
		"S4-enrolled--1":   func(a *ApproverRow, p *PasskeyRow) { a.EnrolledAt = -1 },
		"S5-last-used":     func(a *ApproverRow, p *PasskeyRow) { a.LastUsedAt = 1500 },
		"S5-revoked":       func(a *ApproverRow, p *PasskeyRow) { a.RevokedAt = 1500 },
		"S6-cred-0":        func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = nil },
		"S6-cred-1024":     func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = make([]byte, 1024) },
		"S7-alg-rs256":     func(a *ApproverRow, p *PasskeyRow) { p.Alg = -257 },
		"S7-alg-0":         func(a *ApproverRow, p *PasskeyRow) { p.Alg = 0 },
		"S8-aaguid-15":     func(a *ApproverRow, p *PasskeyRow) { p.AAGUID = make([]byte, 15) },
		"S8-aaguid-17":     func(a *ApproverRow, p *PasskeyRow) { p.AAGUID = make([]byte, 17) },
		"S9-rp-empty":      func(a *ApproverRow, p *PasskeyRow) { p.RPID = "" },
		"S9-origin-empty":  func(a *ApproverRow, p *PasskeyRow) { p.Origin = "" },
		"S10-count--1":     func(a *ApproverRow, p *PasskeyRow) { p.SignCount = -1 },
		"S10-count-2^32":   func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1 << 32 },
		"S10-count-2^63-1": func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1<<63 - 1 },
	}
	for name, edit := range bad {
		a, p := testPasskey("bad", 2)
		edit(&a, &p)
		if err := c.EnrolPasskey(ctx, a, p); !errors.Is(err, ErrBadApprover) {
			t.Fatalf("%s: err %v, want ErrBadApprover", name, err)
		}
		if got := rowCounts(t, c); got != "0 approvers, 0 passkeys" {
			t.Fatalf("%s wrote: %s", name, got)
		}
	}
	// The bounds are inclusive: each of these enrols.
	for i, edit := range []func(a *ApproverRow, p *PasskeyRow){
		func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = []byte{0xee} },
		func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = bytes.Repeat([]byte{0xef}, 1023) },
		func(a *ApproverRow, p *PasskeyRow) { p.SignCount, p.Alg = 0, -8 },
		func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1<<32 - 1 },
	} {
		a, p := testPasskey(fmt.Sprint("edge", i), byte(0xf0+i))
		edit(&a, &p)
		if err := c.EnrolPasskey(ctx, a, p); err != nil {
			t.Fatalf("edge %d: %v", i, err)
		}
	}

	if err := c.EnrolPasskey(ctx, a, p); err != nil {
		t.Fatalf("good enrol: %v", err)
	}
	if err := c.EnrolApprover(ctx, testApprover("ssh-one", 1000)); err != nil {
		t.Fatal(err)
	}
	const after = "6 approvers, 5 passkeys"
	if got := rowCounts(t, c); got != after {
		t.Fatalf("after the good enrols: %s, want %s", got, after)
	}
	check := func(when string) {
		t.Helper()
		for _, r := range []struct {
			name   string
			edit   func(a *ApproverRow, p *PasskeyRow)
			want   error
			holder string
		}{
			{"name-taken", func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "mark", "mark" }, ErrApproverNameTaken, "mark"},
			{"name-taken-by-ssh", func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "ssh-one", "ssh-one" }, ErrApproverNameTaken, "ssh-one"},
			{"key-enrolled", func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = []byte("cose-mark") }, ErrApproverKeyEnrolled, "mark"},
			{"key-enrolled-by-ssh", func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = []byte("key-ssh-one") }, ErrApproverKeyEnrolled, "ssh-one"},
			{"credential-enrolled", func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = bytes.Repeat([]byte{1}, 32) }, ErrPasskeyCredentialEnrolled, "mark"},
		} {
			a, p := testPasskey("new", 3)
			r.edit(&a, &p)
			err := c.EnrolPasskey(ctx, a, p)
			if !errors.Is(err, r.want) || !strings.Contains(err.Error(), `"`+r.holder+`"`) {
				t.Fatalf("%s, %s: err %v, want %v naming %q", when, r.name, err, r.want, r.holder)
			}
			if got := rowCounts(t, c); got != after {
				t.Fatalf("%s, %s wrote: %s", when, r.name, got)
			}
		}
	}
	check("active")
	if _, err := c.db.Exec(`UPDATE approvers SET revoked_at = 1200 WHERE name IN ('mark', 'ssh-one')`); err != nil {
		t.Fatal(err)
	}
	check("revoked")
}

// PS5: EnrolApprover still refuses every format but ssh, so a passkey enters only with its row.
func TestEnrolApproverStillSSHOnly(t *testing.T) {
	c := openTestCredentials(t)
	a, _ := testPasskey("mark", 1)
	if err := c.EnrolApprover(testCtx(t), a); !errors.Is(err, ErrBadApprover) || rowCounts(t, c) != "0 approvers, 0 passkeys" {
		t.Fatalf("EnrolApprover(webauthn): %v, %s", err, rowCounts(t, c))
	}
}
