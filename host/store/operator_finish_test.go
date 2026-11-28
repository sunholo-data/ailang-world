package store

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestOperatorFinishCompletedLaunch(t *testing.T) {
	f, digest, prepared := preparedOperatorLaunchFixture(t, 0)
	id := prepared.CAS.ID
	now := func() int64 { return pinNow }
	secret := [32]byte{1}
	before := f.snapshot(t)
	notCompleted := func() {
		t.Helper()
		if ids, err := f.c.OperatorCompletedLaunches(f.ctx); err != nil || len(ids) != 0 {
			t.Fatalf("incomplete launch discovered: %v %v", ids, err)
		}
		if _, err := f.c.OperatorFinishLaunch(f.ctx, prepared.Episode); !errors.Is(err, ErrOperatorMissing) {
			t.Fatalf("incomplete episode selected: %v", err)
		}
	}
	notCompleted()
	if !reflect.DeepEqual(before, f.snapshot(t)) {
		t.Fatal("candidate read changed the prepared launch")
	}
	if err := f.c.ReserveOperatorLaunch(startRecoveryTestContext(t), prepared, secret, f.proof, now); err != nil {
		t.Fatal(err)
	}
	if err := f.c.AuthorizeOperatorSpawn(startRecoveryTestContext(t), id, secret, f.proof, now); err != nil {
		t.Fatal(err)
	}
	read := func() OperatorLaunchSnapshot {
		t.Helper()
		v, err := f.c.ReadOperatorLaunch(f.ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	receipt := OperatorLaunchReceipt{Run: id, Executable: digest, PID: 42, Birth: "kernel-start", Status: "launched"}
	if err := f.c.RecordOperatorLaunch(startRecoveryTestContext(t), read(), receipt, nil, pinNow); err != nil {
		t.Fatal(err)
	}
	notCompleted()
	receipt.Status = "stopped"
	if err := f.c.RecordOperatorLaunch(startRecoveryTestContext(t), read(), receipt, &EpisodeLaunch{Mode: "headless"}, pinNow); err != nil {
		t.Fatal(err)
	}
	ids, err := f.c.OperatorCompletedLaunches(f.ctx)
	got, readErr := f.c.OperatorFinishLaunch(f.ctx, prepared.Episode)
	if err != nil || !reflect.DeepEqual(ids, []string{id}) || readErr != nil || got.Receipt != receipt || got.Key.Version != 4 {
		t.Fatalf("completed candidate or exact receipt absent: %v %v %+v %v", ids, err, got, readErr)
	}
	v := OperatorFinishReview{Episode: prepared.Episode, Preset: "sol", Store: f.c.lock.dbPath, Report: "fixture"}
	raw, err := json.Marshal(map[string]any{OperatorFinishCatalogueKey: map[string]OperatorFinishReview{v.Episode: v}})
	if err != nil {
		t.Fatal(err)
	}
	h := OperatorHeartbeat{Catalogue: raw}
	review, err := h.FinishReview(v.Episode, v.Preset)
	if err != nil || !reflect.DeepEqual(review, v) {
		t.Fatalf("catalogue finish projection: %+v %v", review, err)
	}
	if _, err := h.FinishReview(v.Episode, "other"); !errors.Is(err, ErrOperatorRunner) {
		t.Fatalf("foreign preset accepted: %v", err)
	}
	if a, err := review.Action(4, true); err != nil || a.Kind != "finish" || !a.Options.Ready || a.Task != "" {
		t.Fatalf("taskless typed review: %+v %v", a, err)
	}
}
