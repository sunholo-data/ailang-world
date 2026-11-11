package approval_test

// Row 201 A2 sprint 2 M3b3 (design §6, PA1, PA3): the assertion law's refusals, in its order, and
// its counter law, made by the test authenticator (passkey_testauth_test.go).

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// paApprover is pk's registration enrolled as mark-mac at 900 (stored counter 1).
func paApprover(t *testing.T, pk *approvaltest.Passkey) approval.Approver {
	t.Helper()
	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: paOrigin, Challenge: paEnrol,
		Response: pk.Register(paOrigin, paRPID, paEnrol), Now: 900})
	if err != nil {
		t.Fatalf("enrol the test passkey: %v", err)
	}
	return a
}

// paInput is a valid AssertionInput: pk's assertion (counter 2) of the PA request (created 1000,
// expires 1600), admitted at 1010.
func paInput(t *testing.T, pk *approvaltest.Passkey, e ...approvaltest.Edit) approval.AssertionInput {
	t.Helper()
	return approval.AssertionInput{
		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: 1600, StoredDigest: paDigest, BytesDigest: paDigest},
		Challenge: paChallenge,
		Response:  pk.Assert(paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), 2, e...),
		Approver:  paApprover(t, pk),
	}
}

// paCase is one PA1 arm: the assertion is signer's (pk; another authenticator when other, under
// pk's credential ID when underPK) at origin and rpID ("" is paOrigin, paRPID) over hash (nil is
// the PA request's message hash) with counter count (0 is 2) and Edit e; then edit changes the
// input; want is the refusal (verdict: admission's).
type paCase struct {
	name         string
	e            approvaltest.Edit
	origin, rpID string
	hash         *[32]byte
	count        uint32
	other        bool
	underPK      bool
	edit         func(*approval.AssertionInput)
	want         error
	verdict      approval.Verdict
}

// PA1: one arm per P1-P5, P3w and W1-W15 on an otherwise valid assertion, in the law's order, for
// both algorithms; each refusal returns a nil approval. Pairs pin the order.
func TestApproveByAssertionOrder(t *testing.T) {
	row, rowHash, zero := paChallenge, paHash(t, paChallenge, paChallenge, 1600), [32]byte{}
	later := paHash(t, paDigest, paChallenge, 1601)
	cases := []paCase{
		{name: "P1-consumed", edit: func(in *approval.AssertionInput) { in.Admit.Consumed = true }, verdict: approval.RefuseConsumed},
		{name: "P1-consumed-and-revoked", edit: func(in *approval.AssertionInput) { in.Admit.Consumed, in.Approver.RevokedAt = true, 950 }, verdict: approval.RefuseConsumed},
		{name: "P1-clock-behind", edit: func(in *approval.AssertionInput) { in.Admit.Now = 999 }, verdict: approval.RefuseClockBehind},
		{name: "P1-expired", edit: func(in *approval.AssertionInput) { in.Admit.Now = 1600 }, verdict: approval.RefuseExpired},
		{name: "P1-bytes-changed", edit: func(in *approval.AssertionInput) { in.Admit.BytesDigest[0] ^= 1 }, verdict: approval.RefuseDigestMismatch},
		{name: "P2-name", edit: func(in *approval.AssertionInput) { in.Approver.Name = "Mark" }, want: approval.ErrApproverName},
		{name: "P3w-format-ssh", edit: func(in *approval.AssertionInput) { in.Approver.Format = approval.FormatSSH }, want: approval.ErrApproverNotPasskey},
		{name: "P3w-no-passkey", edit: func(in *approval.AssertionInput) { in.Approver.Passkey = nil }, want: approval.ErrApproverNotPasskey},
		{name: "P4-revoked-good-assertion", edit: func(in *approval.AssertionInput) { in.Approver.RevokedAt = 950 }, want: approval.ErrApproverRevoked},
		{name: "P4-revoked-bad-assertion", e: approvaltest.Edit{BadSignature: true}, edit: func(in *approval.AssertionInput) { in.Approver.RevokedAt = 950 }, want: approval.ErrApproverRevoked},
		{name: "P5-enrolled-after-now", edit: func(in *approval.AssertionInput) { in.Approver.EnrolledAt = 1011 }, want: approval.ErrApproverClockBehind},
		{name: "W1-envelope", edit: func(in *approval.AssertionInput) {
			in.Response = bytes.Replace(in.Response, []byte(`"public-key"`), []byte(`"x"`), 1)
		}, want: approval.ErrPasskeyResponse},
		{name: "W2-another-credentials-good-assertion", other: true, want: approval.ErrCredentialIDMismatch},
		{name: "W2-raw-id", e: approvaltest.Edit{RawID: []byte("another credential")}, want: approval.ErrCredentialIDMismatch},
		{name: "W3-user-handle", e: approvaltest.Edit{UserHandle: []byte("mallory")}, want: approval.ErrUserHandle},
		{name: "W3-empty-user-handle", e: approvaltest.Edit{UserHandle: []byte{}}, want: approval.ErrUserHandle},
		{name: "W4-client-data", e: approvaltest.Edit{ClientData: []byte(`{"type":"webauthn.get"}`)}, want: approval.ErrClientData},
		{name: "W5-type-create", e: approvaltest.Edit{Type: "webauthn.create"}, want: approval.ErrClientDataType},
		{name: "W5-before-W6", e: approvaltest.Edit{Type: "webauthn.create"}, hash: &row, want: approval.ErrClientDataType},
		{name: "W6-the-row-challenge-not-the-hash", hash: &row, want: approval.ErrChallenge},
		{name: "W6-another-expiry", hash: &later, want: approval.ErrChallenge},
		{name: "W6-another-digest", hash: &rowHash, want: approval.ErrChallenge},
		{name: "W6-negative-expiry-has-no-hash", hash: &zero, edit: func(in *approval.AssertionInput) {
			in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Approver.EnrolledAt = -10, -20, -1, -30
		}, want: approval.ErrChallenge},
		{name: "W7-origin", origin: "https://x.ts.net", want: approval.ErrOriginMismatch},
		{name: "W7-before-W10", origin: "https://y.ts.net", rpID: "y.ts.net", want: approval.ErrOriginMismatch},
		{name: "W8-cross-origin", e: approvaltest.Edit{CrossOrigin: true}, want: approval.ErrCrossOrigin},
		{name: "W8-top-origin", e: approvaltest.Edit{TopOrigin: paOrigin}, want: approval.ErrCrossOrigin},
		{name: "W9-attested-data", e: approvaltest.Edit{Set: 0x40}, want: approval.ErrAuthDataFlags},
		{name: "W9-backup-state-without-eligibility", e: approvaltest.Edit{Set: 0x10}, want: approval.ErrAuthDataFlags},
		{name: "W10-rp-id", rpID: "y.ts.net", want: approval.ErrRPIDHash},
		{name: "W14-bad-signature", e: approvaltest.Edit{BadSignature: true}, want: approval.ErrBadSignature},
		{name: "W14-another-key-same-credential-id", other: true, underPK: true, want: approval.ErrBadSignature},
		{name: "W14-before-W15", count: 1, e: approvaltest.Edit{BadSignature: true}, want: approval.ErrBadSignature},
		{name: "W15-counter-not-advanced", count: 1, want: approval.ErrSignCountRegression},
	}
	// W11-W13 (flag-state coverage, iteration 379's judge): each refusal with every other flag
	// state that reaches it, so a mutant that also needs another bit dies. BE marks the enrolled
	// passkey backup-eligible.
	for _, f := range []struct {
		name               string
		set, clear         byte
		enrolledBE, badSig bool
		want               error
	}{
		{"W11-no-UP", 0, 0x01, false, false, approval.ErrUserPresence},
		{"W11-no-UP-or-UV", 0, 0x05, false, false, approval.ErrUserPresence},
		{"W11-no-UP-BE", 0x08, 0x01, true, false, approval.ErrUserPresence},
		{"W11-no-UP-BE-BS", 0x18, 0x01, true, false, approval.ErrUserPresence},
		{"W11-no-UP-ED", 0x80, 0x01, false, false, approval.ErrUserPresence},
		{"W12-no-UV", 0, 0x04, false, false, approval.ErrUserVerification},
		{"W12-no-UV-BE", 0x08, 0x04, true, false, approval.ErrUserVerification},
		{"W12-no-UV-BE-BS", 0x18, 0x04, true, false, approval.ErrUserVerification},
		{"W12-no-UV-BE-not-enrolled", 0x08, 0x04, false, false, approval.ErrUserVerification},
		{"W12-no-UV-ED", 0x80, 0x04, false, false, approval.ErrUserVerification},
		{"W12-before-W14", 0, 0x04, false, true, approval.ErrUserVerification},
		{"W13-BE-not-enrolled", 0x08, 0, false, false, approval.ErrBackupEligibility},
		{"W13-BE-BS-not-enrolled", 0x18, 0, false, false, approval.ErrBackupEligibility},
		{"W13-enrolled-BE-absent", 0, 0, true, false, approval.ErrBackupEligibility},
		{"W13-enrolled-BE-absent-ED", 0x80, 0, true, false, approval.ErrBackupEligibility},
		{"W13-before-W14", 0x08, 0, false, true, approval.ErrBackupEligibility},
	} {
		be := f.enrolledBE
		cases = append(cases, paCase{name: f.name, want: f.want, edit: func(in *approval.AssertionInput) { in.Approver.Passkey.BackupEligible = be },
			e: approvaltest.Edit{Set: f.set &^ 0x80, Clear: f.clear, Extensions: f.set&0x80 != 0, BadSignature: f.badSig}})
	}
	// Acceptance controls: each knob the arms use, alone and harmless, still verifies.
	enrolBE := func(in *approval.AssertionInput) { in.Approver.Passkey.BackupEligible = true }
	controls := []paCase{
		{name: "clean"},
		{name: "enrolled-at-now", edit: func(in *approval.AssertionInput) { in.Approver.EnrolledAt = in.Admit.Now }},
		{name: "user-handle-is-the-name", e: approvaltest.Edit{UserHandle: []byte("mark-mac")}},
		{name: "extensions", e: approvaltest.Edit{Extensions: true}},
		{name: "BE-enrolled", e: approvaltest.Edit{Set: 0x08}, edit: enrolBE},
		{name: "BE-BS-enrolled", e: approvaltest.Edit{Set: 0x18}, edit: enrolBE},
	}
	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
		pk, other := approvaltest.NewPasskey(alg, 7), approvaltest.NewPasskey(approval.AlgEdDSA, 9)
		for i, c := range append(cases, controls...) {
			t.Run(fmt.Sprint(alg, "/", c.name), func(t *testing.T) {
				in := paInput(t, pk)
				signer, origin, rpID, hash, count := pk, paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), uint32(2)
				if c.other {
					signer = other
				}
				if c.underPK {
					c.e.RawID = pk.CredentialID()
				}
				if c.origin != "" {
					origin = c.origin
				}
				if c.rpID != "" {
					rpID = c.rpID
				}
				if c.hash != nil {
					hash = *c.hash
				}
				if c.count != 0 {
					count = c.count
				}
				in.Response = signer.Assert(origin, rpID, hash, count, c.e)
				if c.edit != nil {
					c.edit(&in)
				}
				s, err := approval.ApproveByAssertion(in)
				var na *approval.NotAdmittedError
				switch {
				case i >= len(cases):
					if err != nil || s.ApprovedBy() != "webauthn:mark-mac" || s.Binds(paDigest, paChallenge, 1600) != nil {
						t.Fatalf("control: (%v, %v), want a verified approval of the PA request", s, err)
					}
				case s != nil:
					t.Fatalf("a refused assertion returned %v", s)
				case c.verdict != 0 && (!errors.As(err, &na) || na.Verdict != c.verdict):
					t.Fatalf("err %v, want not admitted (%s)", err, c.verdict)
				case c.verdict == 0 && !errors.Is(err, c.want):
					t.Fatalf("err %v, want %v", err, c.want)
				}
			})
		}
	}
}

// PA3: the counter law (W15): both zero passes; otherwise the new counter must be strictly greater
// than the enrolled one; the verified approval carries the new counter.
func TestSignCountLaw(t *testing.T) {
	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 7)
	for _, c := range []struct {
		stored, new uint32
		ok          bool
	}{
		{1, 2, true}, {0, 0, true}, {0, 7, true}, {0, 1, true}, {7, 4294967295, true}, {0, 4294967295, true},
		{2, 2, false}, {3, 2, false}, {5, 0, false}, {1, 0, false}, {4294967295, 0, false}, {4294967295, 4294967295, false},
	} {
		in := paInput(t, pk)
		in.Approver.Passkey.SignCount = c.stored
		in.Response = pk.Assert(paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), c.new)
		s, err := approval.ApproveByAssertion(in)
		if !c.ok {
			if s != nil || !errors.Is(err, approval.ErrSignCountRegression) || !strings.Contains(err.Error(), "approver revoke mark-mac") {
				t.Errorf("stored %d, new %d: (%v, %v), want ErrSignCountRegression naming approver revoke mark-mac", c.stored, c.new, s, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("stored %d, new %d: %v, want a verified approval", c.stored, c.new, err)
			continue
		}
		if _, n, ok := s.PasskeyUse(); !ok || n != c.new {
			t.Errorf("stored %d, new %d: PasskeyUse counter %d ok %v, want %d", c.stored, c.new, n, ok, c.new)
		}
	}
}
