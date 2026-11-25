package store

import (
	"crypto/sha256"
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

func TestReserveOperatorWorktreeGuards(t *testing.T) {
	for _, name := range []string{"valid", "kind", "headless", "store", "approver", "challenge", "claimed", "owner", "identity", "expiry"} {
		t.Run(name, func(t *testing.T) {
			c := openTestCredentials(t)
			ctx := testCtx(t)
			cas, runner := operatorClaimFixture(t, c)
			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
				t.Fatal(err)
			}
			row := getOperator(t, c, cas.ID, 1590)
			action, err := approval.DecodeOperatorAction(row.Action)
			if err != nil {
				t.Fatal(err)
			}
			switch name {
			case "kind":
				action.Kind, action.Task, action.Options.Headless = "finish", "", false
			case "headless":
				action.Options.Headless = false
			case "store":
				action.Store += ".other"
			}
			action.Argv = approval.OperatorArgv(action.Kind, action.Episode, action.Options)
			body, err := approval.CanonicalOperatorAction(action)
			if err != nil {
				t.Fatal(err)
			}
			digest, err := approval.OperatorSignatureDigest(action)
			if err != nil {
				t.Fatal(err)
			}
			challenge := [32]byte(row.Challenge)
			message, err := approval.OperatorActionMessage(digest, challenge, row.ExpiresAt)
			if err != nil {
				t.Fatal(err)
			}
			signer := approvaltest.NewSigner(42)
			ap := approval.Approver{Name: "fixture", Format: approval.FormatSSH, Key: signer.PublicKey(), EnrolledAt: 1000}
			if err := c.EnrolApprover(ctx, ApproverRow{Name: ap.Name, Format: ap.Format, PublicKey: ap.Key.Blob, EnrolledAt: ap.EnrolledAt}); err != nil {
				t.Fatal(err)
			}
			artifact := signer.Sign(approval.OperatorActionNamespace, message)
			signed, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: action, Challenge: challenge,
				Expiry: row.ExpiresAt, Now: 1590, Artifact: artifact, Format: ap.Format, Approver: ap})
			if err != nil {
				t.Fatal(err)
			}
			proof := signed.MintProof()
			cas.Digest, cas.Generation = digest, row.Generation
			operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, kind=?, artifact=?, approved_by=?", body, digest[:], action.Kind, artifact, signed.ApprovedBy())
			identity, now, want := "fixture", int64(1590), ErrOperatorConflict
			switch name {
			case "owner":
				operatorExec(t, c, "UPDATE operator_actions SET owner='other'")
			case "identity":
				identity = ""
			case "expiry":
				now, want = row.ExpiresAt, ErrOperatorExpired
				operatorExec(t, c, "UPDATE operator_runner_state SET seen_at=?", now)
			case "kind", "claimed":
				if name == "claimed" {
					// The SQL CAS also rejects running. Make reaching Bind observable.
					proof = nil
					operatorExec(t, c, "UPDATE operator_actions SET state='running'")
				} else {
					// A finish must refuse before decoding even when its stored bytes are damaged.
					broken := body[:1]
					cas.Digest = sha256.Sum256(broken)
					operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?", broken, cas.Digest[:])
				}
			case "approver":
				operatorExec(t, c, "UPDATE operator_actions SET approved_by='ssh:other'")
			case "challenge":
				// Model a damaged row; the original 32 bytes still bind to the proof.
				operatorExec(t, c, "PRAGMA ignore_check_constraints = ON")
				operatorExec(t, c, "UPDATE operator_actions SET challenge=?", append(row.Challenge, 1))
				operatorExec(t, c, "PRAGMA ignore_check_constraints = OFF")
			}
			before := getOperator(t, c, cas.ID, 1590)
			key, err := c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1}, identity, proof, func() int64 { return now })
			if name == "valid" {
				if err != nil {
					t.Fatal(err)
				}
				assertReservedOperatorStep(t, c, key)
				operatorCounts(t, c, 1, 2)
				if after := getOperator(t, c, cas.ID, 1590); after.State != "running" || after.Generation != cas.Generation+1 || after.Owner != runner.Nonce {
					t.Fatalf("incorrect admitted state: %s/%d/%s", after.State, after.Generation, after.Owner)
				}
				return
			}
			if !errors.Is(err, want) || key != (OperatorStepKey{}) ||
				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
				t.Fatalf("refusal changed action or returned the wrong class: %v", err)
			}
			operatorCounts(t, c, 1, 1)
			sessions, err := c.ListSessions(ctx, "", 20)
			start, startErr := c.EpisodeStart(ctx, action.Episode)
			if err != nil || len(sessions) != 0 || startErr != nil || start != nil {
				t.Fatalf("reserve changed mint evidence: sessions=%d start present=%t: %v %v", len(sessions), start != nil, err, startErr)
			}
		})
	}
}

func TestDisposeOperatorUncertainStepGuards(t *testing.T) {
	for _, name := range []string{"kind", "owner", "version", "action_kind"} {
		t.Run(name, func(t *testing.T) {
			c := openTestCredentials(t)
			ctx := testCtx(t)
			cas, runner := operatorClaimFixture(t, c)
			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
				t.Fatal(err)
			}
			cas.Generation++
			key := OperatorStepKey{cas, runner.Nonce, "worktree", "fixture", 1}
			operatorExec(t, c, "UPDATE operator_actions SET state='running'")
			if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
				t.Fatal(err)
			}
			before := getOperator(t, c, cas.ID, 1590)
			bad := key
			want := ErrOperatorConflict
			switch name {
			case "kind":
				bad.Kind, want = "files", ErrOperatorInput
			case "owner":
				operatorExec(t, c, "UPDATE operator_action_steps SET owner='different' WHERE kind='worktree'")
				key.Owner = "different"
			case "version":
				bad.Version++
			case "action_kind":
				operatorExec(t, c, "UPDATE operator_actions SET kind='finish'")
				before = getOperator(t, c, cas.ID, 1590)
			}
			if err := c.DisposeOperatorWorktree(ctx, bad, OperatorWorktreeUncertain, 1591); !errors.Is(err, want) ||
				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
				t.Fatalf("disposition changed action or returned the wrong class: %v", err)
			}
			assertReservedOperatorStep(t, c, key)
		})
	}
}

func TestOperatorStepReceiptAtLimit(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	cas, runner := operatorClaimFixture(t, c)
	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
		t.Fatal(err)
	}
	cas.Generation++
	key := OperatorStepKey{cas, runner.Nonce, "checks", "fixture", 1}
	if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
		t.Fatal(err)
	}
	body := make([]byte, MaxSurfaceApprovalArtifactBytes)
	body[0], body[len(body)-1] = 1, 2
	if err := c.ObserveOperatorStep(ctx, key, body); err != nil {
		t.Fatal(err)
	}
	after, err := c.OperatorStep(ctx, key)
	if err != nil || !reflect.DeepEqual(after, OperatorStepReceipt{"observed", key.Owner, 2, body}) {
		t.Fatalf("at-limit receipt was not stored exactly: %v, size %d", err, len(after.Receipt))
	}
}

func TestReadOperatorWorktreeGuards(t *testing.T) {
	for _, name := range []string{"uncertain", "duplicate"} {
		t.Run(name, func(t *testing.T) {
			c := openTestCredentials(t)
			ctx := testCtx(t)
			cas, runner := operatorClaimFixture(t, c)
			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
				t.Fatal(err)
			}
			cas.Generation++
			key := OperatorStepKey{cas, runner.Nonce, "worktree", "first", 1}
			if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
				t.Fatal(err)
			}
			operatorExec(t, c, "UPDATE operator_actions SET state='running'")
			var second OperatorStepKey
			var want error
			if name == "uncertain" {
				if err := c.DisposeOperatorWorktree(ctx, key, OperatorWorktreeUncertain, 1591); err != nil {
					t.Fatal(err)
				}
				key.Action.Generation++
				key.Version++
			} else {
				second = key
				second.Identity = "second"
				if err := c.ReserveOperatorStep(ctx, second, 1590); err != nil {
					t.Fatal(err)
				}
				want = ErrOperatorMissing
			}
			before := getOperator(t, c, cas.ID, 1591)
			got, signed, err := c.ReadOperatorWorktree(ctx, cas.ID)
			if !errors.Is(err, want) || !reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1591)) ||
				(want == nil && (got != key || signed.ID != cas.ID || signed.Kind != "start")) {
				t.Fatalf("wrong reservation read or action changed: %+v, %v", got, err)
			}
			assertReservedOperatorStep(t, c, key)
			if name == "duplicate" {
				assertReservedOperatorStep(t, c, second)
			}
		})
	}
}
