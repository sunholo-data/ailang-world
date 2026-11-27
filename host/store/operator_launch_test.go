package store

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

func TestOperatorLaunchReceiptLaw(t *testing.T) {
	digest := "sha256:" + hex64("fixture executable")
	f := newOperatorPinFixtureWithFacts(t, false, []approval.OperatorResolved{{Name: "launch_binary_sha256", Value: digest}}, startRecoveryTestContext(t))
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
	before := f.snapshot(t)
	prepared, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
	if err != nil || prepared.Window != OperatorStartPrepared {
		t.Fatalf("prepared: %+v %v", prepared, err)
	}
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
