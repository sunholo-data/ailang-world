package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func testOperatorServiceStartupGate(t *testing.T) {
	for _, tc := range []struct {
		arm, want        string
		loaded, stopFail bool
		boots, stops     int
	}{
		{arm: "same generation", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "late", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "stale", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "future", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "wrong build", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "wrong protocol", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "finish only", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "empty owner", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "bad read", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "failed stop", want: "companion could not be stopped", boots: 1, stops: 1, stopFail: true},
		{arm: "already stopped", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
		{arm: "already loaded stale", want: "loaded companion heartbeat unavailable", loaded: true},
		{arm: "already loaded live", loaded: true},
		{arm: "fresh restart", boots: 1},
		{arm: "bad previous", want: "cannot determine the previous runner lease"},
		{arm: "unknown supervisor", want: "could not inspect the companion"},
	} {
		t.Run(tc.arm, func(t *testing.T) {
			f := operatorService(t)
			f.mustInstall(t)
			current := f.env.terminal.now()
			clock := time.Unix(current, 0)
			reads := 0
			f.env.runnerNow = func() time.Time {
				reads++
				if tc.arm == "late" && reads >= 3 {
					return clock.Add(f.env.wait)
				}
				return clock
			}
			sleeps := 0
			f.env.runnerSleep = func(d time.Duration) {
				sleeps++
				if sleeps > 3 || d <= 0 || d > 100*time.Millisecond {
					t.Fatalf("unbounded poll: %d sleeps, delay %v", sleeps, d)
				}
				clock = clock.Add(d)
			}
			if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
				t.Fatal(err)
			}
			creds := f.config.Credentials
			f.ctl.loaded = tc.loaded
			if tc.stopFail {
				f.ctl.stopErr = errors.New("fake stop failed")
			}
			if tc.arm == "already stopped" {
				f.ctl.stopErr = errServiceNotLoaded
			}
			if tc.arm == "unknown supervisor" {
				f.ctl.runnerPrintErr = errors.New("fake supervisor error")
			}
			if tc.arm == "already loaded stale" {
				credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=113")
			}
			if tc.arm == "bad previous" {
				credSQL(t, creds, "UPDATE operator_runner_state SET generation='broken'")
			}
			knownHash := serviceFileHash(creds)
			newBoot := f.ctl.onBoot
			f.ctl.onBoot = func() error {
				switch tc.arm {
				case "same generation", "failed stop", "already stopped":
					return nil
				}
				if err := newBoot(); err != nil {
					return err
				}
				switch tc.arm {
				case "stale":
					credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=113")
				case "future":
					credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=124")
				case "wrong build":
					credSQL(t, creds, "UPDATE operator_runner_state SET build='wrong'")
				case "wrong protocol":
					credSQL(t, creds, "UPDATE operator_runner_state SET protocol=0")
				case "finish only":
					credSQL(t, creds, "UPDATE operator_runner_state SET readiness='finish'")
				case "empty owner":
					credSQL(t, creds, "UPDATE operator_runner_state SET owner=''")
				case "bad read":
					credSQL(t, creds, "UPDATE operator_runner_state SET catalogue_generation='broken'")
				}
				knownHash = serviceFileHash(creds)
				return nil
			}
			f.writes = nil
			code, output := f.command("start")
			switch {
			case tc.want == "":
				h, err := runnerHeartbeat(creds, time.Second)
				wantGen := int64(1 + tc.boots)
				if code != 0 || !f.ctl.loaded || err != nil || h.Owner.Generation != wantGen || sleeps != 0 {
					t.Fatalf("fresh success: %d %q loaded=%v heartbeat=%+v err=%v", code, output, f.ctl.loaded, h, err)
				}
			case code != 1 || !strings.Contains(output, tc.want):
				t.Fatalf("startup refusal: %d %s", code, output)
			}
			if tc.arm == "same generation" {
				if sleeps != 3 || reads != 8 || !clock.Equal(time.Unix(current, 0).Add(f.env.wait)) || f.ctl.loaded {
					t.Fatalf("startup timeout drift: sleeps=%d clock reads=%d time=%v loaded=%v", sleeps, reads, clock, f.ctl.loaded)
				}
				h, err := runnerHeartbeat(creds, time.Second)
				if err != nil || h.Owner.Generation != 1 || h.SeenAt != current {
					t.Fatalf("old liveness was changed: %+v %v", h, err)
				}
			}
			wantLoaded := tc.loaded || tc.stopFail || tc.arm == "fresh restart"
			if f.ctl.loaded != wantLoaded || f.ctl.boots != tc.boots || f.ctl.stops != tc.stops || f.ctl.kicks != 0 ||
				len(f.writes) != 0 || serviceFileHash(creds) != knownHash {
				t.Fatalf("unexpected start effects: boots=%d stops=%d kicks=%d writes=%v", f.ctl.boots, f.ctl.stops, f.ctl.kicks, f.writes)
			}
		})
	}
}
