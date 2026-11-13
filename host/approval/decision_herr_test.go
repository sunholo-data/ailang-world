package approval_test

// Row 201 A4 sprint 4 M7a (the sprint-3 judge's non-blocking finding 2): a message hash that cannot
// be built (hErr) is refused at W6, in W6's place, never before W1-W5: a decision with no message
// AND an earlier failing step is refused by that earlier step, as the A2 order has it.

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

func TestDecideByAssertionRefusesAMissingMessageAtW6(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(t *testing.T, in *approval.DecisionAssertionInput)
		want error
	}{
		{"W1-envelope", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse},
		{"W2-raw-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2, approvaltest.Edit{RawID: []byte("another")})
		}, approval.ErrCredentialIDMismatch},
		{"W5-type", func(t *testing.T, in *approval.DecisionAssertionInput) {
			in.Response = dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2, approvaltest.Edit{Type: "webauthn.create"})
		}, approval.ErrClientDataType},
		{"W6-itself", func(*testing.T, *approval.DecisionAssertionInput) {}, approval.ErrChallenge},
	} {
		in := dmaInput(t, 2)
		in.Verb = 3 // no decision message: hErr
		tc.edit(t, &in)
		if d, err := approval.DecideByAssertion(in); d != nil || !errors.Is(err, tc.want) {
			t.Errorf("%s: (%v, %v), want %v", tc.name, d, err, tc.want)
		}
	}
}
