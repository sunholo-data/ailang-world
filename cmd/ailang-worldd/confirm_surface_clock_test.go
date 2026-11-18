package main

// Row 243 V3 sprint 6 (judge-404 survivors E-AWAIT-NOW-STALE and E-RC-INPUT-NOW, design §2.7 steps 5 and
// 6(a)): the confirm-on-surface wait reads its own injected clock on every poll and ends at its deadline on
// that clock, and the re-check judges the approver at this process's clock, never the row's. This file
// sorts before the finish rig's, so a wait that never ends fails here, by name, in under a second.

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// The wait polls at expires_at+9 and at +10 (its deadline) and refuses there, on the injected clock: once
// that clock passes the deadline the test fails and the wait is parked (never spun), whatever the wall clock.
func TestAwaitVerbConfirmationEndsOnItsClock(t *testing.T) {
	creds, in := surfaceUnit(t)
	own := storeOwn(t, in, "vr-00000000000000e1", 0xe1, sfT0)
	deadline := own.expires + int64(verbLandGrace/time.Second)
	now, waits, late := deadline-1, 0, make(chan int64, 1)
	env := sessionEnv{now: func() int64 { return now }, after: func(d time.Duration) <-chan time.Time {
		waits++
		if now += int64(d / time.Second); now > deadline+1 {
			late <- now
			return nil // parks the wait: a nil channel never fires
		}
		ch := make(chan time.Time, 1)
		ch <- time.Unix(now, 0)
		return ch
	}}
	done := make(chan string, 1)
	go func() { _, refusal := awaitVerbConfirmation(env, creds, own); done <- refusal }()
	select {
	case refusal := <-done:
		if want := "no confirmation of request vr-00000000000000e1 by " + approverTime(deadline); refusal != want || waits != 1 || now != deadline {
			t.Fatalf("%q after %d waits at %d, want %q after 1 wait at %d", refusal, waits, now, want, deadline)
		}
	case at := <-late:
		t.Fatalf("the wait was still polling at %d on its own clock, past its deadline %d", at, deadline)
	case <-time.After(20 * time.Second):
		t.Fatal("the wait did not end within 20 s of wall time")
	}
}

// The approver is judged at the CLI's own clock: a row that says it was consumed at +50 does not admit a
// key enrolled at +20 to a poll at +19; the same row re-checks at +20 (P5's boundary, the CLI's clock).
func TestPollVerbRequestRechecksAtItsOwnClock(t *testing.T) {
	for _, c := range []struct {
		at   int64
		want string
	}{{sfT0 + 19, approval.ErrApproverClockBehind.Error()}, {sfT0 + 20, ""}} {
		t.Run(fmt.Sprint(c.at-sfT0), func(t *testing.T) {
			creds, in := surfaceUnit(t)
			own := storeOwn(t, in, "vr-00000000000000e2", 0xe2, sfT0)
			approveByHand(t, creds, own.id, verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires), "ssh:mark-mac", sfT0+50)
			credSQL(t, creds, fmt.Sprintf("UPDATE approvers SET enrolled_at = %d WHERE name = 'mark-mac'", sfT0+20))
			row, refusal, err := pollVerbRequest(creds, own, c.at)
			if err != nil || (c.want == "") != (refusal == "") || !strings.Contains(refusal, c.want) || (c.want == "" && row.ConsumedAt != sfT0+50) ||
				(c.want != "" && !strings.HasPrefix(refusal, "the confirmation of request vr-00000000000000e2 does not verify: ")) {
				t.Fatalf("poll at %d: %q %v, want %q", c.at, refusal, err, c.want)
			}
		})
	}
}
