package main

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestOperatorProductionDeadlines(t *testing.T) {
	f := operatorService(t)
	f.mustInstall(t)
	r, err := readServiceRecord(f.env.paths)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	for _, override := range []time.Duration{0, time.Minute} {
		e := f.env
		e.operatorBudget = override
		var requested []time.Duration
		begin := time.Now()
		e.operatorDeadline = func(raw time.Duration, ctx context.Context) {
			requested = append(requested, raw)
			want := raw
			if override != 0 {
				want = override
			}
			got, ok := ctx.Deadline()
			if !ok || got.Before(begin.Add(want)) || got.After(time.Now().Add(want)) {
				t.Errorf("requested %v effective %v: deadline %v (%t)", raw, want, got, ok)
			}
		}
		ticks := make(chan time.Time)
		close(ticks)
		h := f.heartbeat
		if rc := e.runOperatorPreparationLoop(io.Discard, f.config, h, operatorStartSource{workspace: r.Workspace, clock: e.terminal.now}, nil, ticks, lock); rc != exitOK {
			t.Fatalf("empty real loop: %d", rc)
		}
		if !reflect.DeepEqual(requested, []time.Duration{time.Second, 22 * time.Second}) {
			t.Fatalf("production child contexts: %v", requested)
		}
	}
}
