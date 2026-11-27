package store

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

func preparedOperatorLaunchFixture(t *testing.T, ttl int64) (operatorPinFixture, string, OperatorStartRecovery) {
	t.Helper()
	digest := "sha256:" + hex64("fixture executable")
	prefix := []approval.OperatorResolved{{Name: "launch_binary_sha256", Value: digest}}
	f := newOperatorPinFixtureWithFacts(t, false, prefix, startRecoveryTestContext(t))
	if ttl > 0 {
		row := getOperator(t, f.c, f.in.Attempt.CAS.ID, pinNow)
		a, err := approval.DecodeOperatorAction(row.Action)
		if err != nil {
			t.Fatal(err)
		}
		f.in.Start.TTLS = ttl
		facts, err := f.in.Resolution()
		if err != nil {
			t.Fatal(err)
		}
		a.Resolution = append(prefix, facts...)
		f.proof = signOperatorPinAction(t, f.c, row, a)
		f.in.Attempt.CAS.Digest, err = approval.OperatorSignatureDigest(a)
		if err != nil {
			t.Fatal(err)
		}
	}
	f.mint(t)
	id := f.in.Attempt.CAS.ID
	for _, value := range []string{"files-identity", OperatorFilesPublished} {
		files, _, err := f.c.ReadOperatorFiles(f.ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.c.RecordOperatorFiles(f.ctx, files, value, f.proof, func() int64 { return pinNow }); err != nil {
			t.Fatal(err)
		}
	}
	prepared, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || prepared.Window != OperatorStartPrepared {
		t.Fatalf("prepared: %+v %v", prepared, err)
	}
	return f, digest, prepared
}

func TestOperatorLaunchReceiptLaw(t *testing.T) {
	f, digest, prepared := preparedOperatorLaunchFixture(t, 0)
	id, before := prepared.CAS.ID, f.snapshot(t)
	handoff := [32]byte{1}
	now := func() int64 { return pinNow }
	f.assertRefused(t, before, f.c.ReserveOperatorStep(f.ctx, operatorLaunchKey(prepared, handoff), pinNow), ErrOperatorInput)
	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, nil, now), ErrOperatorInput)
	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, func() int64 { return before.action.ExpiresAt }), ErrOperatorExpired)
	wrong := prepared
	wrong.CAS.Generation++
	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, wrong, handoff, f.proof, now), ErrOperatorConflict)
	if err := f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, now); err != nil {
		t.Fatal(err)
	}
	reserved, err := f.c.ReadOperatorLaunch(f.ctx, id)
	if err != nil || reserved.Phase != "reserved" || reserved.Receipt != (OperatorLaunchReceipt{}) {
		t.Fatalf("reserved: %+v %v", reserved, err)
	}
	before = f.snapshot(t)
	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, now), ErrOperatorConflict)
	f.assertRefused(t, before, f.c.ObserveOperatorStep(f.ctx, reserved.Key, []byte("fake")), ErrOperatorInput)
	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, [32]byte{2}, f.proof, now), ErrOperatorConflict)
	if files, action, err := f.c.ReadOperatorSpawnInput(f.ctx, id); err != nil || files.Step.Action != prepared.CAS || action.ID != id {
		t.Fatalf("spawn input readback: %+v %s %v", files.Step, action.ID, err)
	}
	if err := f.c.AuthorizeOperatorSpawn(f.ctx, id, handoff, f.proof, now); err != nil {
		t.Fatal(err)
	}
	before = f.snapshot(t)
	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, handoff, f.proof, now), ErrOperatorConflict)
	reserved, err = f.c.ReadOperatorLaunch(f.ctx, id)
	if err != nil || reserved.Key.Version != 2 || reserved.Phase != "reserved" {
		t.Fatalf("dispatch grant: %+v %v", reserved, err)
	}

	r := OperatorLaunchReceipt{Run: id, Executable: digest, PID: 42, Birth: "kernel-start", Status: "launched"}
	wrongBirthAtStart := r
	wrongBirthAtStart.Birth = ""
	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, wrongBirthAtStart, nil, pinNow), ErrOperatorInput)
	wrongRun := r
	wrongRun.Run += "x"
	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, wrongRun, nil, pinNow), ErrOperatorConflict)
	if err := f.c.RecordOperatorLaunch(f.ctx, reserved, r, nil, before.action.ExpiresAt+1); err != nil {
		t.Fatal(err)
	}
	before = f.snapshot(t)
	live, err := f.c.ReadOperatorLaunch(f.ctx, id)
	if err != nil || live.Receipt != r || before.action.State != "succeeded" || !reflect.DeepEqual(before, f.snapshot(t)) {
		t.Fatalf("observed: %+v %v", live, err)
	}
	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, r, nil, pinNow), ErrOperatorConflict)
	bad := r
	bad.Executable = "sha256:" + hex64("wrong observed executable")
	raw, err := json.Marshal(bad)
	if err != nil {
		t.Fatal(err)
	}
	operatorExec(t, f.c, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", raw)
	corrupt := f.snapshot(t)
	_, err = f.c.ReadOperatorLaunch(f.ctx, id)
	f.assertRefused(t, corrupt, err, ErrOperatorConflict)
	raw, _ = json.Marshal(r)
	operatorExec(t, f.c, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", raw)
	if owner, err := ReadOperatorEpisodeOwner(f.ctx, f.c.lock.dbPath, prepared.Episode); err != nil || owner != id {
		t.Fatalf("observed child lost terminal effect exclusion: %q %v", owner, err)
	}
	if pending, err := EpisodeHasUnresolvedOperatorSteps(f.ctx, f.c.lock.dbPath, prepared.Episode); err != nil || !pending {
		t.Fatalf("observed child lost GC effect exclusion: %t %v", pending, err)
	}
	if err := f.c.PruneOperatorActions(f.ctx, pinNow+31*86400); err != nil || !reflect.DeepEqual(before, f.snapshot(t)) {
		t.Fatalf("pruner changed an unresolved receipt: %v", err)
	}
	r.Status = "capped"
	r.Exit = -1
	history := EpisodeLaunch{Mode: "headless", Exit: r.Exit}
	wrongBirth := r
	wrongBirth.Birth += "x"
	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, live, wrongBirth, &history, pinNow), ErrOperatorConflict)
	if err := f.c.RecordOperatorLaunch(f.ctx, live, r, &history, pinNow); err != nil {
		t.Fatal(err)
	}
	finished, err := f.c.ReadOperatorLaunch(f.ctx, id)
	start, startErr := f.c.EpisodeStart(f.ctx, before.action.Episode)
	if err != nil || startErr != nil || finished.Receipt != r || len(start.Launches) != 1 || start.Launches[0].Exit != r.Exit {
		t.Fatalf("completed: %+v %v, history %v %v", finished, err, start, startErr)
	}
	if err := f.c.RecordOperatorLaunch(f.ctx, live, r, &history, pinNow); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("completion replay: %v", err)
	}
}

func TestOperatorLaunchSession(t *testing.T) {
	f, _, ready := preparedOperatorLaunchFixture(t, 5)
	id, secret := ready.CAS.ID, [32]byte{1}
	now := func() int64 { return pinNow }

	before := f.snapshot(t)
	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, ready, secret, f.proof, func() int64 { return pinNow + 5 }), ErrOperatorExpired)
	if err := f.c.ReserveOperatorLaunch(f.ctx, ready, secret, f.proof, now); err != nil {
		t.Fatal(err)
	}
	reserved, err := f.c.ReadOperatorLaunch(f.ctx, id)
	if err != nil || reserved.Receipt != (OperatorLaunchReceipt{}) {
		t.Fatalf("reserved: %+v %v", reserved, err)
	}
	before = f.snapshot(t)
	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, func() int64 { return pinNow + 5 }), ErrOperatorExpired)
	operatorExec(t, f.c, "UPDATE session_credentials SET project='changed'")
	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
	operatorExec(t, f.c, "UPDATE session_credentials SET project=?", f.in.Profile.Project)
	operatorExec(t, f.c, `INSERT INTO session_credentials
		(credential_id,episode_id,grants_json,expires_at,created_at,project,profile_digest)
		SELECT ?,episode_id,grants_json,expires_at,created_at,project,profile_digest FROM session_credentials`, hex64("other credential"))
	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
	operatorExec(t, f.c, "DELETE FROM session_credentials WHERE credential_id=?", hex64("other credential"))
	operatorExec(t, f.c, "UPDATE approvers SET revoked_at=?", pinNow)
	before = f.snapshot(t)
	if err := f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now); err == nil || !reflect.DeepEqual(before, f.snapshot(t)) {
		t.Fatalf("revoked approver changed a reservation: %v", err)
	}
	operatorExec(t, f.c, "UPDATE approvers SET revoked_at=NULL")
	operatorExec(t, f.c, "DELETE FROM session_credentials")
	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
}
