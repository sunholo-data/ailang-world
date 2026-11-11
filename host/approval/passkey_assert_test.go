package approval

// Row 201 A2 sprint 2 M3b (design §6, PA2-PA4, XW2): the assertion law over Chrome's recorded
// assertions, and what a verified passkey approval hands the consume. The refusal arms are
// passkey_law_test.go's (package approval_test: they need approvaltest, which imports this package).

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

// chromeAssertionInput is the recorded assertion (counter 2) against the passkey the recorded
// registration enrols (counter 1), admitted for the fixture request whose message hash both signed.
func chromeAssertionInput(t *testing.T, alg string) AssertionInput {
	t.Helper()
	c, reg := chromeRegistration(t, alg)
	a, err := NewPasskeyApprover(reg)
	if err != nil {
		t.Fatal(err)
	}
	d := fixtureDigest(t)
	return AssertionInput{Admit: AdmitInput{Now: fixtureExpiry - 300, CreatedAt: fixtureExpiry - 600, ExpiresAt: fixtureExpiry,
		StoredDigest: d, BytesDigest: d}, Challenge: fixtureChallenge(), Response: c.Assertion, Approver: a}
}

// chromeApproval is ApproveByAssertion(chromeAssertionInput(t, "es256")), which must succeed.
func chromeApproval(t *testing.T) *SignedApproval {
	t.Helper()
	s, err := ApproveByAssertion(chromeAssertionInput(t, "es256"))
	if err != nil || s == nil {
		t.Fatalf("the recorded assertion: (%v, %v), want a verified approval", s, err)
	}
	return s
}

// PA2 (XW2): both recorded assertions approve the fixture request through the production law.
func TestApproveByAssertionAcceptsChrome(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		in := chromeAssertionInput(t, alg)
		s, err := ApproveByAssertion(in)
		if err != nil {
			t.Fatalf("%s: %v", alg, err)
		}
		if err := s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry); err != nil {
			t.Fatalf("%s: Binds(the fixture request) = %v", alg, err)
		}
		if got := s.ApprovedBy(); got != "webauthn:mark-mac" {
			t.Fatalf("%s: ApprovedBy() = %q", alg, got)
		}
		n, f, k, ok := s.Use()
		if !ok || n != "mark-mac" || f != FormatWebAuthn || !bytes.Equal(k, in.Approver.Passkey.Key.Canonical()) {
			t.Fatalf("%s: Use() = (%q, %q, %x, %v), want the canonical COSE key", alg, n, f, k, ok)
		}
		id, count, ok := s.PasskeyUse()
		if !ok || !bytes.Equal(id, chromeFixture(t, alg).RegRawID) || count != 2 {
			t.Fatalf("%s: PasskeyUse() = (%x, %d, %v), want the credential ID, 2, true", alg, id, count, ok)
		}
	}
}

// PA3 (XW2, last arm): the recorded assertion (counter 2) is refused once the stored counter is
// the second recorded assertion's (3, design §1.4 P2 of sprint 1).
func TestChromeAssertionAfterAssertion2IsRefused(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		in := chromeAssertionInput(t, alg)
		in.Approver.Passkey.SignCount = binary.BigEndian.Uint32(chromeFixture(t, alg).Asr2AuthData[33:37])
		if in.Approver.Passkey.SignCount != 3 {
			t.Fatalf("%s: assertion2's counter %d, want 3", alg, in.Approver.Passkey.SignCount)
		}
		if s, err := ApproveByAssertion(in); s != nil || !errors.Is(err, ErrSignCountRegression) {
			t.Fatalf("%s: (%v, %v), want ErrSignCountRegression", alg, s, err)
		}
	}
}

// PA4: PasskeyUse is ok only for a verified webauthn approval (B7, B8), and neither the input
// approver nor a returned slice can change it afterwards (judge carry-forward, iteration 379).
func TestSignedApprovalPasskeyUse(t *testing.T) {
	var nilApproval *SignedApproval
	for name, s := range map[string]*SignedApproval{
		"nil":                     nilApproval,
		"zero":                    {},
		"unverified-webauthn":     {format: FormatWebAuthn, name: "mark-mac", credentialID: []byte{1}, signCount: 2},
		"verified-ssh":            verifiedApproval(t),
		"verified-ssh-with-count": func() *SignedApproval { s := verifiedApproval(t); s.credentialID, s.signCount = []byte{1}, 2; return s }(),
	} {
		if id, n, ok := s.PasskeyUse(); ok || id != nil || n != 0 {
			t.Errorf("%s.PasskeyUse() = (%x, %d, %v), want (nil, 0, false)", name, id, n, ok)
		}
	}
	in := chromeAssertionInput(t, "es256")
	want := bytes.Clone(in.Approver.Passkey.CredentialID)
	wantKey := in.Approver.Passkey.Key.Canonical()
	s, err := ApproveByAssertion(in)
	if err != nil {
		t.Fatal(err)
	}
	id, _, _ := s.PasskeyUse()
	id[0] ^= 0xff                               // the caller's copy
	in.Approver.Passkey.CredentialID[1] ^= 0xff // the input ApproveByAssertion read
	in.Approver.Passkey.Key.X[0] ^= 0xff
	if id2, n, ok := s.PasskeyUse(); !ok || n != 2 || !bytes.Equal(id2, want) {
		t.Fatalf("PasskeyUse() after editing a returned ID and the input = (%x, %d, %v), want (%x, 2, true)", id2, n, ok, want)
	}
	if _, _, k, _ := s.Use(); !bytes.Equal(k, wantKey) {
		t.Fatalf("Use() after editing the input key = %x, want %x", k, wantKey)
	}
}

// The registration law's result does not alias its input (judge carry-forward, iteration 379).
func TestNewPasskeyApproverDoesNotAliasItsInput(t *testing.T) {
	_, in := chromeRegistration(t, "es256")
	in.Response = bytes.Clone(in.Response)
	a, err := NewPasskeyApprover(in)
	if err != nil {
		t.Fatal(err)
	}
	p := *a.Passkey
	want := Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, Passkey: &p}
	p.CredentialID, p.Key.X, p.Key.Y = bytes.Clone(p.CredentialID), bytes.Clone(p.Key.X), bytes.Clone(p.Key.Y)
	for i := range in.Response {
		in.Response[i] = 'x'
	}
	if !reflect.DeepEqual(a, want) {
		t.Fatalf("the approver changed when its input did: %+v, want %+v", a.Passkey, want.Passkey)
	}
}
