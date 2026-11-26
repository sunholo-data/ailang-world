package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const pinNow int64 = 1590

type operatorPinFixture struct {
	c     *Credentials
	ctx   context.Context
	in    OperatorMintInput
	proof *approval.OperatorMintProof
}

func newOperatorPinFixture(t *testing.T, allowance bool) operatorPinFixture {
	t.Helper()
	c, ctx := openTestCredentials(t), testCtx(t)
	cas, runner := operatorClaimFixture(t, c)
	row := getOperator(t, c, cas.ID, pinNow)
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(c.lock.dbPath)
	in := OperatorMintInput{
		Attempt:      OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1},
		TokenPath:    filepath.Join(dir, a.Episode+".session"),
		Profile:      approval.ProfileRef{Project: "world", Digest: hashref.SumSHA256([]byte("pin profile")).String()},
		Grants:       []approval.Grant{{Effect: "exec", Scope: "repo", Budget: 4}},
		CredentialID: hex64("operator pin"),
		Start: EpisodeStart{Schema: "world/episode-start/v1", Episode: a.Episode, Agent: "codex",
			Preset: StartPreset{Name: a.Options.Preset}, Repo: dir, TTLS: 600,
			Task: &StartTask{Path: filepath.Join(dir, a.Episode+".task.md"), SHA256: hashref.SumSHA256([]byte(a.Task)).String()}},
	}
	if allowance {
		// Give the fixture a signed, internally consistent mixed-origin mint.
		// A normal allowance draw alone lacks a matching operator mint receipt.
		in.Start.Allowance = &StartAllowance{AllowanceID: "sa-0123456789abcdef", DrawID: "ad-0123456789abcdef", ApprovedBy: "tty", DayKey: "1970-01-01", DrawNumber: 1}
	}
	a.Resolution, err = in.Resolution()
	if err != nil {
		t.Fatal(err)
	}
	proof := signOperatorPinAction(t, c, row, a)
	in.Attempt.CAS.Digest, err = approval.OperatorSignatureDigest(a)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.ClaimOperatorAction(ctx, in.Attempt.CAS, runner, pinNow); err != nil {
		t.Fatal(err)
	}
	in.Attempt.CAS.Generation++
	return operatorPinFixture{c, ctx, in, proof}
}

func signOperatorPinAction(t *testing.T, c *Credentials, row OperatorActionRow, a approval.OperatorAction) *approval.OperatorMintProof {
	t.Helper()
	body, err := approval.CanonicalOperatorAction(a)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := approval.OperatorSignatureDigest(a)
	if err != nil {
		t.Fatal(err)
	}
	challenge := [32]byte(row.Challenge)
	message, err := approval.OperatorActionMessage(digest, challenge, row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	signer := approvaltest.NewSigner(42)
	ap := approval.Approver{Name: "pin", Format: approval.FormatSSH, Key: signer.PublicKey(), EnrolledAt: 900}
	_, enrolled, err := c.Approver(testCtx(t), ap.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !enrolled {
		if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: ap.Name, Format: ap.Format, PublicKey: ap.Key.Blob, EnrolledAt: ap.EnrolledAt}); err != nil {
			t.Fatal(err)
		}
	}
	artifact := signer.Sign(approval.OperatorActionNamespace, message)
	signed, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: a, Challenge: challenge,
		Expiry: row.ExpiresAt, Now: pinNow, Artifact: artifact, Format: ap.Format, Approver: ap})
	if err != nil {
		t.Fatal(err)
	}
	operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, artifact=?, approved_by=? WHERE id=?",
		body, digest[:], artifact, signed.ApprovedBy(), row.ID)
	return signed.MintProof()
}

func (f operatorPinFixture) reserveWorktree(t *testing.T) OperatorStepKey {
	t.Helper()
	key, err := f.c.ReserveOperatorWorktree(f.ctx, f.in.Attempt, "fixture-worktree", f.proof, func() int64 { return pinNow })
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func (f operatorPinFixture) mint(t *testing.T) {
	t.Helper()
	worktree := f.reserveWorktree(t)
	if err := f.c.ObserveOperatorStep(f.ctx, worktree, []byte("created")); err != nil {
		t.Fatal(err)
	}
	in := f.in
	in.Attempt.CAS, in.Attempt.StepVersion = worktree.Action, 1
	key := OperatorStepKey{in.Attempt.CAS, worktree.Owner, "mint", worktree.Action.ID, 1}
	if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
		t.Fatal(err)
	}
	if err := f.c.CommitOperatorMint(f.ctx, in, f.proof, func() int64 { return pinNow }); err != nil {
		t.Fatal(err)
	}
}

type operatorPinState struct {
	action   OperatorActionRow
	steps    string
	sessions []SessionRow
	start    *EpisodeStart
}

func (f operatorPinFixture) snapshot(t *testing.T) operatorPinState {
	t.Helper()
	id := f.in.Attempt.CAS.ID
	s := operatorPinState{action: getOperator(t, f.c, id, pinNow)}
	err := f.c.db.QueryRowContext(f.ctx, `SELECT json_group_array(json_array(kind, identity, phase, version, owner, created_at,
		CASE WHEN receipt IS NULL THEN NULL ELSE hex(receipt) END))
		FROM (SELECT * FROM operator_action_steps WHERE action_id=? ORDER BY kind, identity)`, id).Scan(&s.steps)
	if err != nil {
		t.Fatal(err)
	}
	s.sessions, err = f.c.ListSessions(f.ctx, f.in.Start.Episode, 20)
	if err != nil {
		t.Fatal(err)
	}
	s.start, err = f.c.EpisodeStart(f.ctx, f.in.Start.Episode)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func (f operatorPinFixture) assertRefused(t *testing.T, before operatorPinState, got, want error) {
	t.Helper()
	if after := f.snapshot(t); !errors.Is(got, want) || !reflect.DeepEqual(before, after) {
		t.Fatalf("refusal = %v, want %v; durable rows unchanged = %t", got, want, reflect.DeepEqual(before, after))
	}
}

func TestOperatorFileStoreReadPins(t *testing.T) {
	for _, name := range []string{"valid", "headless", "duplicate_worktree"} {
		t.Run(name, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			f.mint(t)
			id := f.in.Attempt.CAS.ID
			var want error
			switch name {
			case "headless":
				row := getOperator(t, f.c, id, pinNow)
				a, err := approval.DecodeOperatorAction(row.Action)
				if err != nil {
					t.Fatal(err)
				}
				a.Options.Headless = false
				a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
				signOperatorPinAction(t, f.c, row, a)
				want = ErrOperatorConflict
			case "duplicate_worktree":
				row := getOperator(t, f.c, id, pinNow)
				key := OperatorStepKey{OperatorCAS{id, [32]byte(row.Digest), row.Generation}, row.Owner, "worktree", "duplicate", 1}
				if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
					t.Fatal(err)
				}
				want = sql.ErrNoRows
			}
			before := f.snapshot(t)
			got, _, err := f.c.ReadOperatorFiles(f.ctx, id)
			f.assertRefused(t, before, err, want)
			if want == nil && (got.Worktree != "fixture-worktree" || got.Phase != "" || got.Mint.CredentialID != f.in.CredentialID) {
				t.Fatalf("positive read = %+v", got)
			}
		})
	}
}

func TestOperatorFileStoreReservePins(t *testing.T) {
	for _, name := range []string{"valid", "foreign_lock", "empty_identity", "before_created", "different_approver", "allowance", "action_expired"} {
		t.Run(name, func(t *testing.T) {
			f := newOperatorPinFixture(t, name == "allowance")
			f.mint(t)
			id := f.in.Attempt.CAS.ID
			state, action, err := f.c.ReadOperatorFiles(f.ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			value, now, want := "files-identity", pinNow, ErrOperatorConflict
			switch name {
			case "foreign_lock":
				dbPath := f.c.lock.dbPath
				f.c.lock.dbPath = filepath.Join(filepath.Dir(dbPath), "other.db")
				defer func() { f.c.lock.dbPath = dbPath }()
			case "empty_identity":
				value, want = "", ErrOperatorInput
			case "different_approver":
				operatorExec(t, f.c, "UPDATE operator_actions SET approved_by='ssh:other'")
			case "action_expired":
				now, want = getOperator(t, f.c, id, pinNow).ExpiresAt, ErrOperatorExpired
				if now >= state.Session.ExpiresAt {
					t.Fatal("fixture requires a live credential at action expiry")
				}
			}
			before := f.snapshot(t)
			if name == "before_created" {
				// The reservation SQL also rejects a future created_at. Pin
				// the earlier admission independently of that second guard.
				row := before.action
				row.CreatedAt = now + 1
				tx, err := f.c.db.BeginTx(f.ctx, &sql.TxOptions{ReadOnly: true})
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()
				err = admitOperatorFiles(f.ctx, tx, row, state, action, f.proof, now)
				if closeErr := tx.Rollback(); closeErr != nil {
					t.Fatal(closeErr)
				}
				f.assertRefused(t, before, err, want)
				return
			}
			err = f.c.RecordOperatorFiles(f.ctx, state, value, f.proof, func() int64 { return now })
			if name != "valid" {
				f.assertRefused(t, before, err, want)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			key := state.Step
			key.Identity, key.Version = value, 1
			assertReservedOperatorStep(t, f.c, key)
		})
	}
}

func TestOperatorWorktreeStoreReservePins(t *testing.T) {
	for _, name := range []string{"heartbeat_edge", "existing_start", "existing_credential"} {
		t.Run(name, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			switch name {
			case "heartbeat_edge":
				operatorExec(t, f.c, "UPDATE operator_runner_state SET seen_at=?", pinNow-10)
			case "existing_start":
				tx, err := f.c.db.BeginTx(f.ctx, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()
				if _, err := insertEpisodeStart(f.ctx, tx, f.in.Start.Episode, f.in.Start); err != nil {
					t.Fatal(err)
				}
				if err := tx.Commit(); err != nil {
					t.Fatal(err)
				}
			case "existing_credential":
				row := SessionRow{CredentialID: f.in.CredentialID, EpisodeID: f.in.Start.Episode, GrantsJSON: "[]",
					CreatedAt: pinNow, ExpiresAt: pinNow + 600}
				if err := f.c.MintSession(f.ctx, row); err != nil {
					t.Fatal(err)
				}
			}
			before := f.snapshot(t)
			key, err := f.c.ReserveOperatorWorktree(f.ctx, f.in.Attempt, "fixture-worktree", f.proof, func() int64 { return pinNow })
			f.assertRefused(t, before, err, ErrOperatorConflict)
			if key != (OperatorStepKey{}) {
				t.Fatalf("refusal returned worktree key %+v", key)
			}
		})
	}
}

func TestOperatorWorktreeStoreRecoveryPins(t *testing.T) {
	for _, name := range []string{"observed", "read_other_kind", "dispose_other_kind"} {
		t.Run(name, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			key := f.reserveWorktree(t)
			if name == "observed" {
				if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
					t.Fatal(err)
				}
			} else {
				other := OperatorStepKey{key.Action, key.Owner, "mint", key.Action.ID, 1}
				if err := f.c.ReserveOperatorStep(f.ctx, other, pinNow); err != nil {
					t.Fatal(err)
				}
			}
			before := f.snapshot(t)
			if name == "dispose_other_kind" {
				err := f.c.DisposeOperatorWorktree(f.ctx, key, OperatorWorktreeRetained, pinNow)
				f.assertRefused(t, before, err, ErrOperatorConflict)
				return
			}
			_, _, err := f.c.ReadOperatorWorktree(f.ctx, key.Action.ID)
			f.assertRefused(t, before, err, ErrOperatorMissing)
		})
	}
}
