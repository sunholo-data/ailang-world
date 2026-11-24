package authority

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"modernc.org/sqlite"
)

type operatorStoreRig struct {
	path string
	*store.Credentials
	action   approval.OperatorAction
	cas      store.OperatorCAS
	proof    *VerifiedOperatorAction
	artifact []byte
}

func newOperatorStoreRig(t *testing.T, passkey bool) operatorStoreRig {
	t.Helper()
	path, _ := jgRig(t)
	ctx := boundedTestContext(t)
	c := operatorStoreOpen(t, path, store.OpenCredentials)
	canonical, err := store.CanonicalPath(path)
	if err != nil {
		t.Fatal(err)
	}
	a := ojAction()
	a.Store = canonical
	req := store.OperatorRequest{Action: a, OperationID: "1234567890abcdef1234567890abcdef", CreatedAt: 1000, ExpiresAt: 1600, Challenge: sha256.Sum256([]byte(t.Name()))}
	row, err := c.RequestOperatorAction(ctx, req, 1000)
	if err != nil {
		t.Fatal(err)
	}
	in := OperatorJudgeInput{Action: a, Digest: [32]byte(row.Digest), BytesDigest: sha256.Sum256(row.Action),
		Challenge: [32]byte(row.Challenge), CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: 1010}
	var artifact []byte
	if passkey {
		artifact = ojAssert(t, jgPasskey, in, 2)
	} else {
		artifact = ojSSH(t, jgSigner, approval.OperatorActionNamespace, in)
	}
	ro := operatorStoreOpen(t, path, store.OpenCredentialsReadOnly)
	proof, err := JudgeOperatorAction(ctx, ro, in, artifact)
	if err != nil {
		t.Fatal(err)
	}
	return operatorStoreRig{path, c, a, store.OperatorCAS{ID: row.ID, Digest: in.Digest, Generation: row.Generation}, proof, artifact}
}

func operatorStoreOpen(t *testing.T, path string, open func(string) (*store.Credentials, error)) *store.Credentials {
	t.Helper()
	c, err := open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func (r operatorStoreRig) confirm(t *testing.T) store.OperatorActionRow {
	t.Helper()
	row, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1010)
	if err != nil || row.State != "confirmed" || row.ApprovedBy != r.proof.ApprovedBy() {
		t.Fatalf("confirm: %s %q %v", row.State, row.ApprovedBy, err)
	}
	return row
}

func TestOperatorConfirmReplay(t *testing.T) {
	for _, passkey := range []bool{false, true} {
		r := newOperatorStoreRig(t, passkey)
		want := r.confirm(t)
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		c := operatorStoreOpen(t, r.path, store.OpenCredentials)
		ctx := boundedTestContext(t)
		got, err := c.ConfirmOperatorAction(ctx, r.cas, r.artifact, r.proof, 1020)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("confirm replay changed the row or refused: %v", err)
		}
		other := bytes.Clone(r.artifact)
		other[len(other)-1] ^= 1
		_, mismatch := c.ConfirmOperatorAction(ctx, r.cas, other, r.proof, 1021)
		if !errors.Is(mismatch, store.ErrOperatorConflict) {
			t.Errorf("different replay artifact returned %v", mismatch)
		}
		after, err := c.OperatorAction(ctx, r.cas.ID, 1021)
		if err != nil || !reflect.DeepEqual(after, want) {
			t.Fatal("invalid replay changed durable approval")
		}
		name := "mark-ssh"
		if passkey {
			name = "mark-mac"
		}
		a, ok, err := c.Approver(ctx, name)
		if err != nil || !ok || a.LastUsedAt != 1010 {
			t.Fatalf("replay changed approver use: %+v %v", a, err)
		}
		p, ok, err := c.Passkey(ctx, "mark-mac")
		if err != nil || !ok || (passkey && p.SignCount != 2) || (!passkey && p.SignCount != 1) {
			t.Fatalf("unexpected counter: %d %v", p.SignCount, err)
		}
	}
}

func TestOperatorConfirmCrash(t *testing.T) {
	const drift = `CREATE TRIGGER drift BEFORE UPDATE OF last_used_at ON approvers BEGIN UPDATE operator_actions SET `
	for _, tc := range []struct {
		name, sql string
		want      error
		counter   int64
	}{
		{"update-fault", `CREATE TRIGGER crash_operator BEFORE UPDATE OF approved_by ON operator_actions
		   BEGIN SELECT raise(abort, 'test: interrupted before state commit'); END`, nil, 1},
		{"generation drift", drift + `generation = 9; END`, store.ErrOperatorConflict, 1},
		{"digest drift", drift + `digest = zeroblob(32); END`, store.ErrOperatorConflict, 1},
		{"expiry drift", drift + `expires_at = 1010; END`, store.ErrOperatorConflict, 1},
		{"counter-race", `UPDATE webauthn_credentials SET sign_count = 2`, store.ErrSignCountRegression, 2},
		{"revoked", `UPDATE approvers SET revoked_at = 1010 WHERE name = 'mark-mac'`, store.ErrApproverRevoked, 1},
		{"changed-rp", `UPDATE webauthn_credentials SET rp_id = 'replacement.example'`, store.ErrApproverKeyChanged, 1},
		{"changed-origin", `UPDATE webauthn_credentials SET origin = 'https://replacement.example'`, store.ErrApproverKeyChanged, 1},
		{"changed-backup", `UPDATE webauthn_credentials SET backup_eligible = 1`, store.ErrApproverKeyChanged, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newOperatorStoreRig(t, true)
			jgExec(t, r.path, tc.sql)
			ctx := boundedTestContext(t)
			_, err := r.ConfirmOperatorAction(ctx, r.cas, r.artifact, r.proof, 1010)
			if tc.want == nil {
				operatorSQLiteCode(t, err, 1811) // SQLITE_CONSTRAINT_TRIGGER
			} else if !errors.Is(err, tc.want) {
				t.Fatalf("confirmation fault: %v, want %v", err, tc.want)
			}
			r.assertUnconfirmed(t, tc.counter)
		})
	}
}

func operatorSQLiteCode(t *testing.T, err error, code int) {
	t.Helper()
	var se *sqlite.Error
	if !errors.As(err, &se) || se.Code() != code {
		t.Fatalf("confirmation error %v, want SQLite code %d", err, code)
	}
}

func TestOperatorConfirmReadOnly(t *testing.T) {
	r := newOperatorStoreRig(t, true)
	ro := operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly)
	_, err := ro.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1010)
	operatorSQLiteCode(t, err, 8) // SQLITE_READONLY
	r.assertUnconfirmed(t, 1)
}

func (r operatorStoreRig) assertUnconfirmed(t *testing.T, count int64) {
	t.Helper()
	ctx := boundedTestContext(t)
	row, err := r.OperatorAction(ctx, r.cas.ID, 1010)
	if err != nil || row.State != "requested" || row.ApprovedBy != "" || row.Artifact != nil || row.Generation != r.cas.Generation || !bytes.Equal(row.Digest, r.cas.Digest[:]) || row.ExpiresAt != 1600 {
		t.Fatalf("failure changed request: %s %q %v", row.State, row.ApprovedBy, err)
	}
	p, ok, err := r.Passkey(ctx, "mark-mac")
	if err != nil || !ok || p.SignCount != count {
		t.Fatalf("failure changed counter: %d %v", p.SignCount, err)
	}
	a, ok, err := r.Approver(ctx, "mark-mac")
	if err != nil || !ok || a.LastUsedAt != 0 {
		t.Fatalf("failure changed approver use: %d %v", a.LastUsedAt, err)
	}
}

func TestOperatorConfirmStoreBinding(t *testing.T) {
	r := newOperatorStoreRig(t, true)
	for _, tc := range []struct {
		name     string
		cas      store.OperatorCAS
		artifact []byte
		proof    *VerifiedOperatorAction
		now      int64
		want     error
	}{
		{"id", store.OperatorCAS{ID: "missing", Digest: r.cas.Digest, Generation: r.cas.Generation}, r.artifact, r.proof, 1010, store.ErrOperatorMissing},
		{"digest", store.OperatorCAS{ID: r.cas.ID, Generation: r.cas.Generation}, r.artifact, r.proof, 1600, store.ErrOperatorConflict},
		{"generation", store.OperatorCAS{ID: r.cas.ID, Digest: r.cas.Digest, Generation: r.cas.Generation + 1}, r.artifact, nil, 1010, store.ErrOperatorConflict},
		{"nil", r.cas, r.artifact, nil, 1010, approval.ErrNotVerified},
		{"artifact substitution", r.cas, bytes.Clone(r.artifact[:len(r.artifact)-1]), r.proof, 1010, approval.ErrNotVerified},
		{"at expiry", r.cas, r.artifact, r.proof, 1600, store.ErrOperatorExpired},
		{"clock behind", r.cas, r.artifact, r.proof, 999, store.ErrOperatorConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := r.ConfirmOperatorAction(boundedTestContext(t), tc.cas, tc.artifact, tc.proof, tc.now); !errors.Is(err, tc.want) {
				t.Fatalf("confirmation: %v, want %v", err, tc.want)
			}
			r.assertUnconfirmed(t, 1)
		})
	}
	foreign := newOperatorStoreRig(t, true)
	_, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, foreign.artifact, foreign.proof, 1010)
	if binding, ok := errors.AsType[*approval.BindingError](err); !ok || binding.Field != "digest" {
		t.Fatalf("proof for different store: %v", err)
	}
	r.assertUnconfirmed(t, 1)

	if row, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1000); err != nil || row.State != "confirmed" {
		t.Fatalf("confirm at creation: %s %v", row.State, err)
	}
}

func testOperatorClaimedRecheck(t *testing.T) {
	r := newOperatorStoreRig(t, true)
	confirmed := r.confirm(t)
	ctx := boundedTestContext(t)
	jgExec(t, r.path, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1010, 'start', 1, NULL)`)
	if err := r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{Nonce: "runner", Generation: 7}, 1010); err != nil {
		t.Fatal(err)
	}
	row, err := r.OperatorAction(ctx, r.cas.ID, 1010)
	if err != nil || row.State != "claimed" {
		t.Fatalf("claim read: %s %v", row.State, err)
	}
	ro := operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly)
	challenge := [32]byte(confirmed.Challenge)
	check := func(a approval.OperatorAction, now int64, row store.OperatorActionRow) (*VerifiedOperatorAction, error) {
		return RecheckClaimedOperatorAction(ctx, ro, a, challenge, confirmed.ExpiresAt, now, row)
	}
	// The row may disagree with the runner's frozen evidence.
	row.Digest, row.Challenge, row.ExpiresAt = nil, nil, 1601
	proof, err := check(r.action, 1599, row)
	if err != nil || !proof.Rechecked() || proof.ApprovedBy() != r.proof.ApprovedBy() {
		t.Fatalf("claimed receipt proof or exact approver lost: %v", err)
	}
	if _, err := proof.OperatorConfirmation().ConfirmationFor(row.Artifact); !errors.Is(err, approval.ErrNotVerified) {
		t.Fatalf("runner proof was accepted as an initial confirmation: %v", err)
	}
	wrong := r.action
	wrong.Task += "changed"
	if proof, err := check(wrong, 1599, row); proof != nil || err == nil {
		t.Fatal("row supplied authority for a different frozen intent")
	}
	if proof, err := check(r.action, 1600, row); proof != nil || err == nil {
		t.Fatal("row supplied authority after recorded expiry")
	}
	row.ApprovedBy = "webauthn:replacement"
	if proof, err := check(r.action, 1599, row); proof != nil || !errors.Is(err, ErrOperatorApprover) {
		t.Fatalf("trusted row identity: %v", err)
	}
	row.ApprovedBy = confirmed.ApprovedBy
	row.Artifact = nil
	if proof, err := check(r.action, 1599, row); proof != nil || !errors.Is(err, ErrUnknownArtifact) {
		t.Fatalf("claimed state supplied authority without an artifact: %v", err)
	}
}
