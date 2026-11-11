package store

// Row 201 A2 sprint 2 M4c (design §6, PS3, RC6): the passkey lookups read back exactly what was
// enrolled, and miss without an error.

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// PS3: the three lookups return the enrolled rows exactly, miss without an error, and refuse a v5 store.
func TestPasskeyLookups(t *testing.T) {
	ctx := testCtx(t)
	c := openTestCredentials(t)
	for i, name := range []string{"c", "a", "b"} {
		a, p := testPasskey(name, byte(10+i))
		if err := c.EnrolPasskey(ctx, a, p); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.EnrolApprover(ctx, testApprover("s", 1000)); err != nil {
		t.Fatal(err)
	}
	if err := c.RevokeApprover(ctx, "b", 1100); err != nil {
		t.Fatal(err)
	}
	if _, err := c.db.Exec(`UPDATE approvers SET last_used_at = 1050 WHERE name = 'a'`); err != nil {
		t.Fatal(err)
	}
	wantA, wantP := testPasskey("a", 11)
	wantA.LastUsedAt = 1050
	if b, _, ok, err := c.PasskeyByCredentialID(ctx, bytes.Repeat([]byte{12}, 32)); err != nil || !ok || b.Name != "b" || b.RevokedAt != 1100 {
		t.Fatalf("PasskeyByCredentialID(b's), revoked: %+v ok %v err %v", b, ok, err)
	}
	gotA, gotP, ok, err := c.PasskeyByCredentialID(ctx, wantP.CredentialID)
	if err != nil || !ok || !reflect.DeepEqual(gotA, wantA) || !reflect.DeepEqual(gotP, wantP) {
		t.Fatalf("PasskeyByCredentialID(a's): %+v %+v ok %v err %v", gotA, gotP, ok, err)
	}
	for name, id := range map[string][]byte{"prefix": wantP.CredentialID[:31], "longer": append(bytes.Clone(wantP.CredentialID), 0), "none": nil} {
		if a, p, ok, err := c.PasskeyByCredentialID(ctx, id); ok || err != nil || a.Name != "" || p.Name != "" {
			t.Fatalf("PasskeyByCredentialID(%s): %+v %+v ok %v err %v; want a miss", name, a, p, ok, err)
		}
	}
	if p, ok, err := c.Passkey(ctx, "b"); err != nil || !ok || p.Name != "b" || p.CredentialID[0] != 12 {
		t.Fatalf("Passkey(b), revoked: %+v ok %v err %v", p, ok, err)
	}
	for _, name := range []string{"s", "", "a%", "zz"} {
		if p, ok, err := c.Passkey(ctx, name); ok || err != nil || p.Name != "" {
			t.Fatalf("Passkey(%q): %+v ok %v err %v; want a miss", name, p, ok, err)
		}
	}
	names := func() string {
		t.Helper()
		got, err := c.ActivePasskeys(ctx)
		if err != nil {
			t.Fatal(err)
		}
		var n []string
		for _, e := range got {
			if e.Approver.Name != e.Passkey.Name || e.Approver.Format != "webauthn" {
				t.Fatalf("ActivePasskeys pair %+v", e)
			}
			n = append(n, e.Passkey.Name)
		}
		return strings.Join(n, ",")
	}
	if got := names(); got != "a,c" {
		t.Fatalf("ActivePasskeys = %s, want a,c (b revoked, s ssh)", got)
	}
	// A passkey row beside an ssh approver (made by hand) is not an active passkey.
	if _, err := c.db.Exec(`INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count)
		VALUES ('s', x'99', 'localhost', 'http://localhost', -7, zeroblob(16), 0, 0)`); err != nil {
		t.Fatal(err)
	}
	if got := names(); got != "a,c" {
		t.Fatalf("ActivePasskeys with an ssh row's passkey = %s, want a,c", got)
	}
	for i := 0; i < 70; i++ {
		a, p := testPasskey(fmt.Sprintf("m%02d", i), 0)
		p.CredentialID = []byte(a.Name)
		if err := c.EnrolPasskey(ctx, a, p); err != nil {
			t.Fatal(err)
		}
	}
	if got := strings.Split(names(), ","); len(got) != 64 || got[0] != "a" || got[1] != "c" || got[63] != "m61" || !sort.StringsAreSorted(got) {
		t.Fatalf("ActivePasskeys over 72 active: %d names, first %v, last %v; want 64 by name", len(got), got[:2], got[len(got)-1])
	}

	c5 := openV5(t)
	_, _, _, e1 := c5.PasskeyByCredentialID(ctx, wantP.CredentialID)
	_, _, e2 := c5.Passkey(ctx, "a")
	_, e3 := c5.ActivePasskeys(ctx)
	for i, err := range []error{e1, e2, e3} {
		if !isOlderPasskeys(err) {
			t.Fatalf("v5 lookup %d: %v, want the passkeys refusal", i, err)
		}
	}
}

// passkeyRowColumns maps every PasskeyRow field to its webauthn_credentials column.
var passkeyRowColumns = map[string]string{
	"Name": "name", "CredentialID": "credential_id", "RPID": "rp_id", "Origin": "origin", "Alg": "alg",
	"AAGUID": "aaguid", "BackupEligible": "backup_eligible", "SignCount": "sign_count",
}

// RC6: PasskeyRow, the table and the rows the three lookups scan agree field for field, each a
// distinct non-zero value.
func TestPasskeyRowMatchesTheTable(t *testing.T) {
	ctx := testCtx(t)
	c := openTestCredentials(t)
	typ := reflect.TypeOf(PasskeyRow{})
	var fields, keys, cols, want []string
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	for k, v := range passkeyRowColumns {
		keys, want = append(keys, k), append(want, v)
	}
	rows, err := c.db.Query(`PRAGMA table_info(webauthn_credentials)`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatal(err)
		}
		cols = append(cols, name)
	}
	_ = rows.Close()
	for _, s := range [][]string{fields, keys, cols, want} {
		sort.Strings(s)
	}
	if !reflect.DeepEqual(fields, keys) || !reflect.DeepEqual(cols, want) {
		t.Fatalf("PasskeyRow fields %v vs map %v; columns %v vs map %v", fields, keys, cols, want)
	}
	a := ApproverRow{Name: "rc6", Format: "webauthn", PublicKey: []byte("rc6-key"), EnrolledAt: 1000}
	p := PasskeyRow{Name: "rc6", CredentialID: []byte{9, 8, 7}, RPID: "rp.ts.net", Origin: "https://rp.ts.net", Alg: -8,
		AAGUID: bytes.Repeat([]byte{5}, 16), BackupEligible: true, SignCount: 4294967290}
	if err := c.EnrolPasskey(ctx, a, p); err != nil {
		t.Fatal(err)
	}
	_, byID, _, err1 := c.PasskeyByCredentialID(ctx, p.CredentialID)
	byName, _, err2 := c.Passkey(ctx, "rc6")
	active, err3 := c.ActivePasskeys(ctx)
	if err1 != nil || err2 != nil || err3 != nil || len(active) != 1 {
		t.Fatalf("lookups: %v %v %v %d", err1, err2, err3, len(active))
	}
	for _, got := range []PasskeyRow{byID, byName, active[0].Passkey} {
		gv, wv := reflect.ValueOf(got), reflect.ValueOf(p)
		for i := 0; i < typ.NumField(); i++ {
			if gv.Field(i).IsZero() || !reflect.DeepEqual(gv.Field(i).Interface(), wv.Field(i).Interface()) {
				t.Fatalf("field %s = %v, want %v", typ.Field(i).Name, gv.Field(i).Interface(), wv.Field(i).Interface())
			}
		}
	}
}
