package approval_test

// Row 201 A4 sprint 2 M4d (design §2.3, §6 DM3-DM6, DM8; the passkey halves are sprint 3's): the
// decision signature law with approvaltest keys (D-WORLD-109). An external test package because
// approvaltest imports host/approval.

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

var dmP, dmQ = [32]byte{0x50}, [32]byte{0x51}

const dmNow, dmExpiry = int64(1710), int64(1760)

// dmInput is a valid decision input: mark-mac's approve of (P, Q, expiry), signed by s's key over
// the decision message, enrolled before now.
func dmInput(t *testing.T, s *approvaltest.Signer, verb byte) approval.DecisionSignatureInput {
	t.Helper()
	msg, err := approval.DecisionMessage(verb, dmP, dmQ, dmExpiry)
	if err != nil {
		t.Fatal(err)
	}
	return approval.DecisionSignatureInput{Verb: verb, Packet: dmP, Request: dmQ, Expiry: dmExpiry, Now: dmNow,
		Signature: s.Sign(approval.DecisionNamespace, msg),
		Approver:  approval.Approver{Name: "mark-mac", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: dmNow - 100}}
}

// DM8 + the control: P2-P5 each refuse a good signature (name, format, revoked, clock behind);
// the valid input decides as ssh:mark-mac, binds its subject and hands the key to the evidence row.
func TestDecideBySignatureApprover(t *testing.T) {
	s := approvaltest.NewSigner(41)
	for _, tc := range []struct {
		name string
		edit func(a *approval.Approver)
		want error
	}{
		{"P2 name", func(a *approval.Approver) { a.Name = "Mark" }, approval.ErrApproverName},
		{"P3 format", func(a *approval.Approver) { a.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat},
		{"P4 revoked", func(a *approval.Approver) { a.RevokedAt = dmNow - 1 }, approval.ErrApproverRevoked},
		{"P5 clock behind", func(a *approval.Approver) { a.EnrolledAt = dmNow + 1 }, approval.ErrApproverClockBehind},
	} {
		in := dmInput(t, s, approval.VerbApprove)
		tc.edit(&in.Approver)
		if d, err := approval.DecideBySignature(in); d != nil || !errors.Is(err, tc.want) {
			t.Errorf("%s: %v, %v; want %v", tc.name, d, err, tc.want)
		}
	}
	d, err := approval.DecideBySignature(dmInput(t, s, approval.VerbApprove))
	if err != nil || d.DecidedBy() != "ssh:mark-mac" || d.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry) != nil {
		t.Fatalf("control: %v, %v", d, err)
	}
	if name, format, key, ok := d.Use(); !ok || name != "mark-mac" || format != approval.FormatSSH || string(key) != string(s.PublicKey().Blob) {
		t.Fatalf("Use = %q %q %x %v", name, format, key, ok)
	}
}

// DM3 (SSH half): the mint message signed under the mint namespace, and the mint message signed
// under the decision namespace, are refused by the decision law; DM4: a decision signature is
// refused by the mint law (namespace), so neither artifact crosses.
func TestMintSignatureIsNotADecision(t *testing.T) {
	s := approvaltest.NewSigner(42)
	mint, err := approval.SignedMessage(dmQ, dmP, dmExpiry)
	if err != nil {
		t.Fatal(err)
	}
	decision, _ := approval.DecisionMessage(approval.VerbApprove, dmP, dmQ, dmExpiry)
	for name, sig := range map[string][]byte{"mint namespace": s.Sign(approval.Namespace, mint), "mint message": s.Sign(approval.DecisionNamespace, mint),
		"decision message, mint namespace": s.Sign(approval.Namespace, decision)} {
		in := dmInput(t, s, approval.VerbApprove)
		in.Signature = sig
		if _, err := approval.DecideBySignature(in); err == nil {
			t.Errorf("%s: a mint signature decided", name)
		}
	}
	in := dmInput(t, s, approval.VerbApprove)
	if err := approval.VerifyApproval(in.Signature, dmQ, dmP, dmExpiry, s.PublicKey()); !errors.Is(err, approval.ErrNamespace) {
		t.Fatalf("DM4: a decision signature checked as a mint: %v, want ErrNamespace", err)
	}
}

// DM5: an approve signature posted as a deny (and the reverse) does not verify.
func TestDecisionBindsTheVerb(t *testing.T) {
	s := approvaltest.NewSigner(43)
	for _, v := range [][2]byte{{approval.VerbApprove, approval.VerbDeny}, {approval.VerbDeny, approval.VerbApprove}} {
		in := dmInput(t, s, v[0])
		in.Verb = v[1]
		if _, err := approval.DecideBySignature(in); !errors.Is(err, approval.ErrBadSignature) {
			t.Errorf("signed %d, posted %d: %v, want ErrBadSignature", v[0], v[1], err)
		}
	}
	in := dmInput(t, s, approval.VerbApprove)
	in.Expiry++ // the expiry is signed too
	if _, err := approval.DecideBySignature(in); !errors.Is(err, approval.ErrBadSignature) {
		t.Errorf("another expiry: %v", err)
	}
}

// DM6: Binds refuses each field that differs and the zero value; DecidedBy and Use are empty for
// an unverified decision.
func TestSignedDecisionBinds(t *testing.T) {
	d, err := approval.DecideBySignature(dmInput(t, approvaltest.NewSigner(44), approval.VerbDeny))
	if err != nil {
		t.Fatal(err)
	}
	for field, err := range map[string]error{
		"verb":    d.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry),
		"packet":  d.Binds(approval.VerbDeny, dmQ, dmQ, dmExpiry),
		"request": d.Binds(approval.VerbDeny, dmP, dmP, dmExpiry),
		"expiry":  d.Binds(approval.VerbDeny, dmP, dmQ, dmExpiry+1),
	} {
		var be *approval.BindingError
		if !errors.As(err, &be) || be.Field != field {
			t.Errorf("%s: %v", field, err)
		}
	}
	var zero *approval.SignedDecision
	if err := zero.Binds(approval.VerbDeny, dmP, dmQ, dmExpiry); !errors.Is(err, approval.ErrNotVerified) || zero.DecidedBy() != "" {
		t.Fatalf("nil decision: %v %q", err, zero.DecidedBy())
	}
	if _, _, _, ok := (&approval.SignedDecision{}).Use(); ok || (&approval.SignedDecision{}).DecidedBy() != "" || !errors.Is((&approval.SignedDecision{}).Binds(0, dmP, dmQ, 0), approval.ErrNotVerified) {
		t.Fatal("the zero decision is usable")
	}
}
