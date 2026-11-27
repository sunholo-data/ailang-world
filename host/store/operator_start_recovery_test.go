package store

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestOperatorStartRecoveryWindows(t *testing.T) {
	for i, want := range []OperatorStartWindow{OperatorStartUntouched, OperatorStartGitPending,
		OperatorStartTreeOnly, OperatorStartMintPending, OperatorStartMintOnly, OperatorStartFilesPending, OperatorStartPrepared} {
		t.Run([]string{"claimed", "reserved_worktree", "observed_worktree", "reserved_mint", "observed_mint", "reserved_files", "prepared"}[i], func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			id := f.in.Attempt.CAS.ID
			if i >= 4 {
				f.mint(t)
			} else if i >= 1 {
				worktree := f.reserveWorktree(t)
				if i >= 2 {
					if err := f.c.ObserveOperatorStep(f.ctx, worktree, []byte("created")); err != nil {
						t.Fatal(err)
					}
					if i == 3 {
						key := OperatorStepKey{worktree.Action, worktree.Owner, "mint", id, 1}
						if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			if i >= 5 {
				read, _, err := f.c.ReadOperatorFiles(f.ctx, id)
				if err != nil {
					t.Fatal(err)
				}
				if err := f.c.RecordOperatorFiles(f.ctx, read, "files-identity", f.proof, func() int64 { return pinNow }); err != nil {
					t.Fatal(err)
				}
				if i == 6 {
					read, _, err = f.c.ReadOperatorFiles(f.ctx, id)
					if err != nil {
						t.Fatal(err)
					}
					if err := f.c.RecordOperatorFiles(f.ctx, read, OperatorFilesPublished, f.proof, func() int64 { return pinNow }); err != nil {
						t.Fatal(err)
					}
				}
			}

			before := f.snapshot(t)
			got, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
			if err != nil || got.Window != want || got.CAS.Generation != before.action.Generation ||
				got.Owner != before.action.Owner || !reflect.DeepEqual(before, f.snapshot(t)) {
				t.Fatalf("read recovery = %q %v, want %q with unchanged state", got.Window, err, want)
			}
			stale := got
			stale.CAS.Generation++
			f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, stale, nil, pinNow), ErrOperatorConflict)
			switch want {
			case OperatorStartUntouched, OperatorStartTreeOnly, OperatorStartMintPending:
				if want == OperatorStartUntouched {
					f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, got, f.proof, before.action.ExpiresAt), ErrOperatorExpired)
					if err := f.c.DisposeOperatorStartGap(f.ctx, got, f.proof, pinNow); err != nil {
						t.Fatal(err)
					}
				} else if err := f.c.DisposeOperatorStartGap(f.ctx, got, nil, pinNow); err != nil {
					t.Fatal(err)
				}
				after := f.snapshot(t)
				state := "failed"
				if want == OperatorStartUntouched {
					state = "confirmed"
				}
				if after.action.State != state || after.action.Owner != "" || after.action.Generation != got.CAS.Generation+1 ||
					before.steps != after.steps || !reflect.DeepEqual(before.start, after.start) ||
					!reflect.DeepEqual(before.sessions, after.sessions) {
					t.Fatalf("disposition = %s with owner %q and changed evidence", after.action.State, after.action.Owner)
				}
			default:
				f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, got, nil, pinNow), ErrOperatorConflict)
			}
		})
	}
}

func TestOperatorStartRecoveryRejectsContradictoryEvidence(t *testing.T) {
	for _, kind := range []string{"unexpected", "foreign", "duplicate", "nonempty", "wrong_kind", "version_zero", "unowned", "bad_receipt", "wrong_mint_identity"} {
		t.Run(kind, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			id := f.in.Attempt.CAS.ID
			if kind == "nonempty" {
				f.mint(t)
				operatorExec(t, f.c, "UPDATE operator_action_steps SET phase='reserved', receipt=NULL WHERE kind='mint'")
			} else if kind == "wrong_kind" {
				operatorExec(t, f.c, "UPDATE operator_actions SET kind='finish'")
			} else if kind == "unowned" {
				operatorExec(t, f.c, "UPDATE operator_actions SET owner=''")
			} else {
				key := f.reserveWorktree(t)
				switch kind {
				case "wrong_mint_identity":
					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
						t.Fatal(err)
					}
					mint := OperatorStepKey{key.Action, key.Owner, "mint", id + "-other", 1}
					if err := f.c.ReserveOperatorStep(f.ctx, mint, pinNow); err != nil {
						t.Fatal(err)
					}
				case "version_zero":
					operatorExec(t, f.c, "UPDATE operator_action_steps SET version=0 WHERE kind='worktree'")
				case "bad_receipt":
					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("other")); err != nil {
						t.Fatal(err)
					}
				case "unexpected":
					operatorExec(t, f.c, "UPDATE operator_action_steps SET kind='launch' WHERE kind='worktree'")
				case "foreign":
					operatorExec(t, f.c, "UPDATE operator_action_steps SET owner='foreign' WHERE kind='worktree'")
				case "duplicate":
					key.Identity += "-second"
					if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
						t.Fatal(err)
					}
				}
			}
			before := f.snapshot(t)
			current, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
			if kind == "nonempty" {
				if err != nil || current.Window != OperatorStartMintUncertain {
					t.Fatalf("contradictory mint commit evidence: %q %v", current.Window, err)
				}
				f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, current, nil, pinNow), ErrOperatorConflict)
				return
			}
			if !errors.Is(err, ErrOperatorConflict) || !reflect.DeepEqual(before, f.snapshot(t)) {
				t.Fatalf("contradictory evidence: window=%q err=%v", current.Window, err)
			}
		})
	}
}

func TestOperatorStartDispositionRejectsBlobOwner(t *testing.T) {
	f := newOperatorPinFixture(t, false)
	id := f.in.Attempt.CAS.ID
	reviewed, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || reviewed.Window != OperatorStartUntouched || reviewed.Owner == "" {
		t.Fatalf("untouched start: %+v %v", reviewed, err)
	}

	// The TEXT column accepts a BLOB. Go scans both as the same string, but
	// SQLite must still match the actual stored owner for disposition.
	operatorExec(t, f.c, "UPDATE operator_actions SET owner=CAST(owner AS BLOB) WHERE id=?", id)
	before := f.snapshot(t)
	current, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || current != reviewed {
		t.Fatalf("BLOB owner changed the reviewed window: %+v %v", current, err)
	}
	f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, reviewed, nil, pinNow), ErrOperatorConflict)

	operatorExec(t, f.c, "UPDATE operator_actions SET owner=CAST(owner AS TEXT) WHERE id=?", id)
	if err := f.c.DisposeOperatorStartGap(f.ctx, reviewed, nil, pinNow); err != nil {
		t.Fatalf("disposition with matching TEXT owner: %v", err)
	}
}

func TestOperatorPartialFilesOrphanRevoke(t *testing.T) {
	f := newOperatorPinFixture(t, false)
	f.mint(t)
	id := f.in.Attempt.CAS.ID
	state, _, err := f.c.ReadOperatorFiles(f.ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.c.RecordOperatorFiles(f.ctx, state, "file-identity", f.proof, func() int64 { return pinNow }); err != nil {
		t.Fatal(err)
	}
	state, _, err = f.c.ReadOperatorFiles(f.ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	other := SessionRow{CredentialID: hex64("other"), EpisodeID: f.in.Start.Episode,
		GrantsJSON: "[]", ExpiresAt: pinNow + 1, CreatedAt: pinNow}
	if err := f.c.MintSession(f.ctx, other); err != nil {
		t.Fatal(err)
	}
	before := f.snapshot(t)
	stale := state
	stale.Step.Version++
	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, stale, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
	operatorExec(t, f.c, `CREATE TRIGGER refuse_orphan_delete BEFORE DELETE ON session_credentials
		WHEN OLD.credential_id != '`+other.CredentialID+`' BEGIN SELECT RAISE(IGNORE); END`)
	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
	operatorExec(t, f.c, "DROP TRIGGER refuse_orphan_delete")
	operatorExec(t, f.c, `CREATE TRIGGER move_orphan_episode AFTER UPDATE OF state ON operator_actions
		WHEN NEW.state='failed' BEGIN UPDATE session_credentials SET episode_id='other'
		WHERE credential_id <> '`+other.CredentialID+`'; END`)
	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
	operatorExec(t, f.c, "DROP TRIGGER move_orphan_episode")
	if err := f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }); err != nil {
		t.Fatal(err)
	}
	after := f.snapshot(t)
	if after.action.State != "failed" || after.action.Owner != "" || len(after.sessions) != 1 ||
		after.sessions[0].CredentialID != hex64("other") || !reflect.DeepEqual(before.start, after.start) {
		t.Fatalf("orphan disposition did not preserve unrelated credential and exact start: %+v", after)
	}
}

func TestNextOwnedOperatorStart(t *testing.T) {
	c, ctx := openTestCredentials(t), testCtx(t)
	add := func(n int, episode, kind string, owner bool) OperatorActionRow {
		t.Helper()
		req := operatorRequest(t, c, n, episode, kind)
		row := putOperator(t, c, req, req.CreatedAt)
		if owner {
			signOperatorPinActionAt(t, c, row, req.Action, req.CreatedAt)
			operatorExec(t, c, "UPDATE operator_actions SET state='claimed', owner='runner' WHERE id=?", row.ID)
		}
		return getOperator(t, c, row.ID, req.CreatedAt)
	}
	last := add(4, "last", "start", true)
	first := add(3, "first", "start", true)
	add(2, "unowned", "start", false)
	add(1, "finish", "finish", true)
	before, err := os.ReadFile(c.lock.dbPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, run := range []struct{ after, want string }{{"", first.ID}, {first.ID, last.ID}, {last.ID, ""}} {
		got, err := c.NextOwnedOperatorStart(ctx, run.after)
		if got.ID != run.want || (run.want == "" && !errors.Is(err, ErrOperatorMissing)) || (run.want != "" && err != nil) {
			t.Fatalf("owned start after %q: got %q (%v), want %q", run.after, got.ID, err, run.want)
		}
		assertGCLedgerBytes(t, c.lock.dbPath, before)
	}
}

func TestOperatorMarkMintUncertain(t *testing.T) {
	f := newOperatorPinFixture(t, false)
	id := f.in.Attempt.CAS.ID
	tree := f.reserveWorktree(t)
	if err := f.c.ObserveOperatorStep(f.ctx, tree, []byte("created")); err != nil {
		t.Fatal(err)
	}
	key := OperatorStepKey{tree.Action, tree.Owner, "mint", id, 1}
	if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
		t.Fatal(err)
	}
	pending, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || pending.Window != OperatorStartMintPending {
		t.Fatalf("pending mint: %+v %v", pending, err)
	}
	f.assertRefused(t, f.snapshot(t), f.c.MarkOperatorMintUncertain(f.ctx, pending), ErrOperatorConflict)
	if err := f.c.MintSession(f.ctx, SessionRow{CredentialID: f.in.CredentialID, EpisodeID: f.in.Start.Episode,
		GrantsJSON: "[]", ExpiresAt: pinNow + 1, CreatedAt: pinNow}); err != nil {
		t.Fatal(err)
	}
	reviewed, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || reviewed.Window != OperatorStartMintUncertain {
		t.Fatalf("uncertain mint: %+v %v", reviewed, err)
	}
	before := f.snapshot(t)
	f.assertRefused(t, before, f.c.MarkOperatorMintUncertain(f.ctx, pending), ErrOperatorConflict)
	if err := f.c.MarkOperatorMintUncertain(f.ctx, reviewed); err != nil {
		t.Fatal(err)
	}
	after := f.snapshot(t)
	if after.action.State != "uncertain" || after.action.Owner != before.action.Owner || after.action.Generation != before.action.Generation+1 ||
		before.steps != after.steps || !reflect.DeepEqual(before.sessions, after.sessions) || !reflect.DeepEqual(before.start, after.start) {
		t.Fatalf("uncertain state lost evidence: %+v", after)
	}
	current, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.c.MarkOperatorMintUncertain(f.ctx, current); err != nil || !reflect.DeepEqual(after, f.snapshot(t)) {
		t.Fatalf("uncertain disposition replay changed evidence: %v", err)
	}
}
