package store

import (
	"errors"
	"sync"
	"testing"
)

func operatorClaimFixture(t *testing.T, c *Credentials) (OperatorCAS, OperatorOwner) {
	t.Helper()
	row := putOperator(t, c, operatorRequest(t, c, 1, "claim", "start"), 1000)
	// This store-only fixture exercises the SQL claim. The actual proof and
	// signed artifact are exercised end-to-end in host/authority.
	operatorExec(t, c, `UPDATE operator_actions SET state = 'confirmed', approved_by = 'ssh:fixture', artifact = action`)
	operatorExec(t, c, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1590, 'start', 1, NULL)`)
	return OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}, OperatorOwner{Nonce: "runner", Generation: 7}
}

func TestOperatorClaimBeforeExpiry(t *testing.T) {
	c := openTestCredentials(t)
	cas, runner := operatorClaimFixture(t, c)
	ctx := testCtx(t)
	for _, tc := range []struct {
		name   string
		cas    OperatorCAS
		runner OperatorOwner
		now    int64
		want   error
	}{
		{"exact expiry", cas, runner, 1600, ErrOperatorExpired},
		{"past expiry", cas, runner, 1601, ErrOperatorExpired},
		{"clock behind", cas, runner, 999, ErrOperatorConflict},
		{"future heartbeat", cas, runner, 1589, ErrOperatorRunner},
		{"readiness", cas, runner, 1599, ErrOperatorRunner},
		{"past heartbeat", cas, runner, 1599, ErrOperatorRunner},
		{"owner", cas, OperatorOwner{Nonce: "other", Generation: 7}, 1590, ErrOperatorRunner},
		{"runner generation", cas, OperatorOwner{Nonce: "runner", Generation: 6}, 1590, ErrOperatorRunner},
		{"action generation", OperatorCAS{ID: cas.ID, Digest: cas.Digest, Generation: 2}, OperatorOwner{}, 1590, ErrOperatorConflict},
		{"digest", OperatorCAS{ID: cas.ID, Generation: 1}, OperatorOwner{}, 1590, ErrOperatorConflict},
		{"id", OperatorCAS{ID: "missing", Digest: cas.Digest, Generation: 1}, runner, 1590, ErrOperatorMissing},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "readiness" {
				operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'none'`)
				defer operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'start'`)
			}
			if tc.name == "past heartbeat" {
				operatorExec(t, c, `UPDATE operator_runner_state SET seen_at = 1589`)
				defer operatorExec(t, c, `UPDATE operator_runner_state SET seen_at = 1590`)
			}
			if err := c.ClaimOperatorAction(ctx, tc.cas, tc.runner, tc.now); !errors.Is(err, tc.want) {
				t.Fatalf("claim: %v want %v", err, tc.want)
			}
			got := getOperator(t, c, cas.ID, 1599)
			if got.State != "confirmed" || got.Owner != "" || got.Generation != cas.Generation {
				t.Fatalf("claim refusal changed ownership: %+v", got)
			}
			operatorCounts(t, c, 1, 1)
		})
	}
	operatorExec(t, c, `UPDATE operator_actions SET state = 'cancelled', terminal_at = 1590`)
	if err := c.ClaimOperatorAction(ctx, cas, OperatorOwner{}, 1599); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("cancelled: %v", err)
	}
	operatorCounts(t, c, 1, 1)
	if got := getOperator(t, c, cas.ID, 1599); got.State != "cancelled" || got.Owner != "" {
		t.Fatal("cancel won without retaining state")
	}
}

func TestOperatorClaimRace(t *testing.T) {
	c := openTestCredentials(t)
	cas, runner := operatorClaimFixture(t, c)
	operatorExec(t, c, `UPDATE operator_actions SET created_at = 1599; UPDATE operator_runner_state SET seen_at = 1599`)
	const contenders = 2
	var ready sync.WaitGroup
	ready.Add(contenders)
	start, done := make(chan struct{}), make(chan error, contenders)
	ctx := testCtx(t)
	for range contenders {
		go func() {
			ready.Done()
			<-start
			done <- c.ClaimOperatorAction(ctx, cas, runner, 1599)
		}()
	}
	ready.Wait()
	close(start)
	a, b := <-done, <-done
	if a != nil {
		a, b = b, a
	}
	if a != nil || !errors.Is(b, ErrOperatorConflict) {
		t.Fatalf("race: %v / %v", a, b)
	}
	got := getOperator(t, c, cas.ID, 1599)
	if got.State != "claimed" || got.Owner != runner.Nonce || got.Generation != cas.Generation+1 {
		t.Fatalf("claim result: %+v", got)
	}
	operatorCounts(t, c, 1, 1)
	// Simulate corruption of the active-episode index. The CAS's owned-
	// sibling predicate must still refuse a second claim.
	operatorExec(t, c, `DROP INDEX operator_active_episode`)
	sibling := operatorRequest(t, c, 2, "claim", "finish")
	operatorExec(t, c, `INSERT INTO operator_actions
		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation, approved_by, artifact)
		SELECT ?, episode_id, 'finish', review, digest, action, challenge, created_at, expires_at,
			'confirmed', 1, approved_by, artifact FROM operator_actions WHERE id = ?`, sibling.Action.ID, cas.ID)
	operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'both'`)
	cas.ID = sibling.Action.ID
	if err := c.ClaimOperatorAction(ctx, cas, runner, 1599); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("claim alongside owned sibling: %v", err)
	}
	if second := getOperator(t, c, cas.ID, 1599); second.State != "confirmed" || second.Owner != "" {
		t.Fatalf("second claim changed state: %+v", second)
	}

	operatorCounts(t, c, 2, 1)
	operatorExec(t, c, `UPDATE operator_actions SET kind = 'finish', v3_id = 'vr-0000000000000001' WHERE state = 'claimed'`)
	if _, err := c.db.Exec(`UPDATE operator_actions SET v3_id = 'vr-0000000000000001' WHERE state = 'confirmed'`); err == nil {
		t.Fatal("the same V3 card was linked to two actions")
	}
}
