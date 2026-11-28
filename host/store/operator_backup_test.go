package store

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

type storeOperatorProof struct {
	signed *approval.SignedOperatorAction
}

func (p storeOperatorProof) OperatorConfirmation() *approval.SignedOperatorAction { return p.signed }

func TestOperatorBackupEligibilityDrift(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	key := approvaltest.NewPasskey(approval.AlgEdDSA, 82)
	a, p := testPasskey("fixture", 82)
	challenge := [32]byte{1}
	enrolled, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: a.Name, Origin: p.Origin, Challenge: challenge,
		Response: key.Register(p.Origin, p.RPID, challenge), Now: a.EnrolledAt})
	if err != nil {
		t.Fatal(err)
	}
	a.PublicKey, p.CredentialID, p.Alg, p.AAGUID = enrolled.Passkey.Key.Canonical(), key.CredentialID(), approval.AlgEdDSA, enrolled.Passkey.AAGUID[:]
	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
		t.Fatal(err)
	}
	req := operatorRequest(t, c, 1, "drift", "start")
	row := putOperator(t, c, req, 1000)
	cas := OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
	hash, err := approval.OperatorActionMessageHash(cas.Digest, req.Challenge, req.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	artifact := key.Assert(p.Origin, p.RPID, hash, 2)
	signed, err := approval.ConfirmOperatorAction(approval.OperatorConfirmInput{Action: req.Action, Challenge: req.Challenge, Artifact: artifact,
		Format: a.Format, Approver: enrolled, Admit: approval.AdmitInput{Now: 1001, CreatedAt: req.CreatedAt, ExpiresAt: req.ExpiresAt,
			StoredDigest: cas.Digest, BytesDigest: cas.Digest}})
	if err != nil {
		t.Fatal(err)
	}
	setBackup := func(eligible bool) {
		t.Helper()
		if _, err := c.db.ExecContext(testCtx(t), `UPDATE webauthn_credentials SET backup_eligible = ? WHERE name = ?`, eligible, a.Name); err != nil {
			t.Fatal(err)
		}
	}
	setBackup(true)
	read := func() [2]int64 {
		t.Helper()
		var seen [2]int64
		if err := c.db.QueryRowContext(ctx, `SELECT COALESCE(last_used_at, 0), sign_count FROM approvers JOIN webauthn_credentials USING(name) WHERE name = ?`, a.Name).
			Scan(&seen[0], &seen[1]); err != nil {
			t.Fatal(err)
		}
		return seen
	}
	before := read()
	if _, err := c.ConfirmOperatorAction(testCtx(t), cas, artifact, storeOperatorProof{signed}, 1001); !errors.Is(err, ErrApproverKeyChanged) {
		t.Fatalf("backup drift: %v", err)
	}
	if !reflect.DeepEqual(row, getOperator(t, c, row.ID, 1001)) || read() != before {
		t.Fatal("backup drift wrote the action or key counter")
	}
	operatorCounts(t, c, 1, 1)
	setBackup(false)
	if confirmed, err := c.ConfirmOperatorAction(testCtx(t), cas, artifact, storeOperatorProof{signed}, 1001); err != nil || confirmed.State != "confirmed" || read() != [2]int64{1001, 2} {
		t.Fatalf("matching control: state=%s key=%v err=%v", confirmed.State, read(), err)
	}
}
