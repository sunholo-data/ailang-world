package authority

// Row 201 A2 sprint 2 M5 (design §6, AU3): a verified passkey approval mints with approved_by
// webauthn:<name> and moves the passkey's counter in the consume (A6). The passkey is approvaltest's
// test authenticator (D-WORLD-109).

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// AU3: A6 hands the store the credential and counter of a verified passkey approval; a zero
// SignedApproval is still ErrNotVerified (A4) and writes nothing.
func TestMintForRequestDerivesPasskeyUse(t *testing.T) {
	const origin, rpID = "http://localhost:7644", "localhost"
	pk := approvaltest.NewPasskey(approval.AlgES256, 7)
	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: [32]byte{9},
		Response: pk.Register(origin, rpID, [32]byte{9}), Now: 900})
	if err != nil {
		t.Fatal(err)
	}
	p := a.Passkey
	st, mr, row := au1Store(t)
	if err := st.EnrolPasskey(boundedTestContext(t), store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: 900},
		store.PasskeyRow{Name: a.Name, CredentialID: p.CredentialID, RPID: p.RPID, Origin: p.Origin, Alg: p.Key.Alg, AAGUID: p.AAGUID[:], SignCount: int64(p.SignCount)}); err != nil {
		t.Fatal(err)
	}
	d, _ := hex.DecodeString(row.Digest)
	h, err := approval.MessageHash([32]byte(d), au1Challenge, row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	s, err := approval.ApproveByAssertion(approval.AssertionInput{
		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: row.ExpiresAt, StoredDigest: [32]byte(d), BytesDigest: [32]byte(d)},
		Challenge: au1Challenge, Response: pk.Assert(origin, rpID, h, 5), Approver: a,
	})
	if err != nil {
		t.Fatal(err)
	}
	zero := au1RequestMint(mr, row)
	zero.ExpiresAt, zero.Signed = row.ExpiresAt, &approval.SignedApproval{}
	if err := au1Mint(t, st, zero); !errors.Is(err, approval.ErrNotVerified) {
		t.Fatalf("a zero approval: %v, want ErrNotVerified", err)
	}
	assertNotMinted(t, st, row)
	rm := au1RequestMint(mr, row)
	rm.ExpiresAt, rm.Signed = row.ExpiresAt, s
	if err := au1Mint(t, st, rm); err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	got, _, err1 := st.MintRequest(ctx, row.RequestID)
	stored, _, err2 := st.Passkey(ctx, "mark-mac")
	sessions, err3 := st.ListSessions(ctx, "", store.MaxListSessions)
	if err1 != nil || err2 != nil || err3 != nil || got.ApprovedBy != "webauthn:mark-mac" || stored.SignCount != 5 || len(sessions) != 1 {
		t.Fatalf("after the passkey mint: approved_by %q, counter %d, %d sessions (%v %v %v)", got.ApprovedBy, stored.SignCount, len(sessions), err1, err2, err3)
	}
}
