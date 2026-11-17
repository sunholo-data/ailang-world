package approval

// Row 243 V3 sprint 1 M1c (design §6 VL4, the recorded half): Chrome's recorded assertions (A2's
// fixtures, testdata/webauthn) signed a mint request's message hash; the verb laws refuse each at
// W6, while the mint law (the control) approves it.

import (
	"errors"
	"testing"
)

func TestChromeMintAssertionIsNotAVerb(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		in := chromeAssertionInput(t, alg)
		if _, err := ApproveByAssertion(in); err != nil {
			t.Fatalf("%s: control: the recorded assertion does not mint: %v", alg, err)
		}
		v, err := ConfirmByAssertion(VerbAssertionInput{Admit: in.Admit, Challenge: in.Challenge, Response: in.Response, Approver: in.Approver})
		if v != nil || !errors.Is(err, ErrChallenge) {
			t.Fatalf("%s: ConfirmByAssertion = (%v, %v), want ErrChallenge", alg, v, err)
		}
		v, err = CheckConfirmation(VerbCheckInput{Digest: in.Admit.StoredDigest, Challenge: in.Challenge, Expiry: in.Admit.ExpiresAt, Now: in.Admit.Now,
			Artifact: in.Response, Format: FormatWebAuthn, Approver: in.Approver})
		if v != nil || !errors.Is(err, ErrChallenge) {
			t.Fatalf("%s: CheckConfirmation = (%v, %v), want ErrChallenge", alg, v, err)
		}
	}
}
