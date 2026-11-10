package approval

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// ---- M4a ----

// sigInput is a valid SignatureInput over the checked-in fixtures: the ed25519 approver
// "mark-mac", enrolled at fixtureExpiry-1000, admitted at fixtureExpiry-300 for a request made at
// fixtureExpiry-600 that expires at fixtureExpiry, signed by testdata/sshsig/ed25519.sig.
func sigInput(t *testing.T) SignatureInput {
	t.Helper()
	d := fixtureDigest(t)
	return SignatureInput{
		Admit:     AdmitInput{Now: fixtureExpiry - 300, CreatedAt: fixtureExpiry - 600, ExpiresAt: fixtureExpiry, StoredDigest: d, BytesDigest: d},
		Challenge: fixtureChallenge(),
		Signature: fixtureBytes(t, "ed25519.sig"),
		Approver:  Approver{Name: "mark-mac", Format: FormatSSH, Key: fixtureKey(t, "ed25519.pub"), EnrolledAt: fixtureExpiry - 1000},
	}
}

// verifiedApproval is ApproveBySignature(sigInput(t)), which must succeed.
func verifiedApproval(t *testing.T) *SignedApproval {
	t.Helper()
	s, err := ApproveBySignature(sigInput(t))
	if err != nil || s == nil {
		t.Fatalf("the clean input: (%v, %v), want a verified approval", s, err)
	}
	return s
}

// SA1: one arm per branch of ApproveBySignature, in its order (P1-P6).
func TestApproveBySignatureOrder(t *testing.T) {
	other := fixtureDigest(t)
	other[0] ^= 1
	notAdmitted := func(v Verdict) func(error) bool {
		return func(err error) bool {
			var na *NotAdmittedError
			return errors.As(err, &na) && na.Verdict == v
		}
	}
	is := func(target error) func(error) bool { return func(err error) bool { return errors.Is(err, target) } }
	arms := []struct {
		name   string
		mutate func(*SignatureInput)
		want   func(error) bool
	}{
		{"consumed", func(in *SignatureInput) { in.Admit.Consumed = true }, notAdmitted(RefuseConsumed)},
		{"consumed-and-revoked", func(in *SignatureInput) { in.Admit.Consumed = true; in.Approver.RevokedAt = fixtureExpiry - 900 }, notAdmitted(RefuseConsumed)},
		{"clock-behind-request", func(in *SignatureInput) { in.Admit.Now = fixtureExpiry - 601 }, notAdmitted(RefuseClockBehind)},
		{"expired", func(in *SignatureInput) { in.Admit.Now = fixtureExpiry }, notAdmitted(RefuseExpired)},
		{"bytes-changed", func(in *SignatureInput) { in.Admit.BytesDigest = other }, notAdmitted(RefuseDigestMismatch)},
		{"name-upper", func(in *SignatureInput) { in.Approver.Name = "Mark" }, is(ErrApproverName)},
		{"name-empty", func(in *SignatureInput) { in.Approver.Name = "" }, is(ErrApproverName)},
		{"format-webauthn", func(in *SignatureInput) { in.Approver.Format = FormatWebAuthn }, is(ErrApproverFormat)},
		{"revoked-good-signature", func(in *SignatureInput) { in.Approver.RevokedAt = fixtureExpiry - 900 }, is(ErrApproverRevoked)},
		{"revoked-bad-signature", func(in *SignatureInput) {
			in.Approver.RevokedAt = fixtureExpiry - 900
			in.Signature = withLastSigByteFlipped(t, in.Signature)
		}, is(ErrApproverRevoked)},
		{"enrolled-after-now", func(in *SignatureInput) { in.Approver.EnrolledAt = in.Admit.Now + 1 }, is(ErrApproverClockBehind)},
		{"git-namespace", func(in *SignatureInput) { in.Signature = fixtureBytes(t, "git-ns-ed25519.sig") }, is(ErrNamespace)},
		{"p256-signature", func(in *SignatureInput) { in.Signature = fixtureBytes(t, "p256.sig") }, is(ErrKeyMismatch)},
		{"flipped-signature", func(in *SignatureInput) { in.Signature = withLastSigByteFlipped(t, in.Signature) }, is(ErrBadSignature)},
		{"not-armored", func(in *SignatureInput) { in.Signature = []byte("x") }, is(ErrArmor)},
		{"other-digest", func(in *SignatureInput) { in.Admit.StoredDigest, in.Admit.BytesDigest = other, other }, is(ErrBadSignature)},
		{"other-challenge", func(in *SignatureInput) { in.Challenge[0] ^= 1 }, is(ErrBadSignature)},
		{"other-expiry", func(in *SignatureInput) { in.Admit.ExpiresAt = fixtureExpiry + 1 }, is(ErrBadSignature)},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			in := sigInput(t)
			a.mutate(&in)
			s, err := ApproveBySignature(in)
			if s != nil || !a.want(err) {
				t.Fatalf("(%v, %v): want a nil approval and the %s refusal", s, err, a.name)
			}
		})
	}
	t.Run("enrolled-at-now", func(t *testing.T) {
		in := sigInput(t)
		in.Approver.EnrolledAt = in.Admit.Now
		if s, err := ApproveBySignature(in); err != nil || s == nil {
			t.Fatalf("(%v, %v): an approver enrolled at now is active", s, err)
		}
	})
	t.Run("clean", func(t *testing.T) {
		if err := verifiedApproval(t).Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry); err != nil {
			t.Fatalf("the clean approval does not bind its own request: %v", err)
		}
	})
}

// SA2: Binds compares each of digest, challenge and expiry, and refuses an unverified value.
func TestSignedApprovalBinds(t *testing.T) {
	s := verifiedApproval(t)
	d, c, e := fixtureDigest(t), fixtureChallenge(), int64(fixtureExpiry)
	if err := s.Binds(d, c, e); err != nil {
		t.Fatalf("Binds(own request) = %v", err)
	}
	var nilApproval *SignedApproval
	for name, u := range map[string]*SignedApproval{"zero": {}, "nil": nilApproval} {
		if err := u.Binds(d, c, e); !errors.Is(err, ErrNotVerified) {
			t.Errorf("%s.Binds = %v, want ErrNotVerified", name, err)
		}
	}
	d2, c2 := d, c
	d2[31] ^= 1
	c2[31] ^= 1
	for _, a := range []struct {
		field string
		err   error
	}{
		{"digest", s.Binds(d2, c, e)},
		{"challenge", s.Binds(d, c2, e)},
		{"expiry", s.Binds(d, c, e+1)},
		{"expiry", s.Binds(d, c, e-1)},
	} {
		var be *BindingError
		if !errors.As(a.err, &be) || be.Field != a.field {
			t.Errorf("Binds with %s changed = %v, want a *BindingError naming %s", a.field, a.err, a.field)
		}
	}
}

// ---- M4b ----

// SA3: the accessors fail closed on an unverified value, and Use hands out copies of the key.
func TestSignedApprovalAccessorsFailClosed(t *testing.T) {
	var nilApproval *SignedApproval
	for name, u := range map[string]*SignedApproval{"zero": {}, "nil": nilApproval} {
		if got := u.ApprovedBy(); got != "" {
			t.Errorf("%s.ApprovedBy() = %q, want \"\"", name, got)
		}
		if n, f, k, ok := u.Use(); ok || n != "" || f != "" || k != nil {
			t.Errorf("%s.Use() = (%q, %q, %x, %v), want zero values and false", name, n, f, k, ok)
		}
	}
	in := sigInput(t)
	s, err := ApproveBySignature(in)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.ApprovedBy(); got != "ssh:mark-mac" {
		t.Errorf("ApprovedBy() = %q, want ssh:mark-mac", got)
	}
	want := bytes.Clone(in.Approver.Key.Blob)
	n, f, k, ok := s.Use()
	if !ok || n != "mark-mac" || f != FormatSSH || !bytes.Equal(k, want) {
		t.Fatalf("Use() = (%q, %q, %x, %v), want (mark-mac, ssh, the ed25519 blob, true)", n, f, k, ok)
	}
	k[0] ^= 0xff                    // the caller's copy
	in.Approver.Key.Blob[1] ^= 0xff // the input ApproveBySignature read
	if _, _, k2, _ := s.Use(); !bytes.Equal(k2, want) {
		t.Fatalf("Use() after editing a returned key and the input key = %x, want %x (copies)", k2, want)
	}
}

// signedApprovalFields is the disposition of every SignedApproval field: what reads it.
var signedApprovalFields = map[string]string{
	"ok":        "B1/B5/B6: an unverified value binds nothing and has no approver",
	"digest":    "B2: Binds",
	"challenge": "B3: Binds",
	"expiry":    "B4: Binds",
	"name":      "ApprovedBy and Use",
	"format":    "ApprovedBy and Use",
	"key":       "Use",
}

// RC1: every SignedApproval field has a disposition, and changing any one of them in a verified
// value changes what Binds, ApprovedBy or Use report.
func TestSignedApprovalFieldsAreAllBound(t *testing.T) {
	typ := reflect.TypeOf(SignedApproval{})
	for i := 0; i < typ.NumField(); i++ {
		if _, ok := signedApprovalFields[typ.Field(i).Name]; !ok {
			t.Errorf("SignedApproval.%s has no disposition in signedApprovalFields", typ.Field(i).Name)
		}
	}
	if typ.NumField() != len(signedApprovalFields) {
		t.Fatalf("SignedApproval has %d fields, signedApprovalFields %d", typ.NumField(), len(signedApprovalFields))
	}
	perturb := map[string]func(*SignedApproval){
		"ok":        func(s *SignedApproval) { s.ok = false },
		"digest":    func(s *SignedApproval) { s.digest[0] ^= 1 },
		"challenge": func(s *SignedApproval) { s.challenge[0] ^= 1 },
		"expiry":    func(s *SignedApproval) { s.expiry++ },
		"name":      func(s *SignedApproval) { s.name += "x" },
		"format":    func(s *SignedApproval) { s.format += "x" },
		"key":       func(s *SignedApproval) { s.key = bytes.Clone(s.key); s.key[0] ^= 0xff },
	}
	if len(perturb) != len(signedApprovalFields) {
		t.Fatalf("%d perturbations for %d fields", len(perturb), len(signedApprovalFields))
	}
	observe := func(s *SignedApproval) string {
		n, f, k, ok := s.Use()
		return fmt.Sprintf("%v|%s|%s|%s|%x|%v", s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry), s.ApprovedBy(), n, f, k, ok)
	}
	base := observe(verifiedApproval(t))
	for field, p := range perturb {
		s := verifiedApproval(t)
		p(s)
		if observe(s) == base {
			t.Errorf("changing %s changed nothing Binds, ApprovedBy or Use report", field)
		}
	}
}
