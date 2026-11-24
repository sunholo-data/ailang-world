package store

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestOperatorRunnerUnavailable(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	if _, err := c.OperatorHeartbeat(ctx); !errors.Is(err, ErrOperatorRunner) {
		t.Fatalf("empty heartbeat: %v", err)
	}
	h := OperatorHeartbeat{Owner: OperatorOwner{Generation: 1, Nonce: "first"}, Protocol: OperatorRunnerProtocol, Build: "test", Readiness: "start", SeenAt: 10}
	for _, tc := range []struct {
		name, build, kind string
		now               int64
		change            func(*OperatorHeartbeat)
		want              bool
	}{
		{"now", "test", "start", 10, nil, true},
		{"future", "test", "start", 9, nil, false},
		{"just fresh", "test", "start", 19, nil, true},
		{"exactly stale", "test", "start", 20, nil, false},
		{"wrong build", "other", "start", 10, nil, false},
		{"wrong protocol", "test", "start", 10, func(h *OperatorHeartbeat) { h.Protocol++ }, false},
		{"no generation", "test", "start", 10, func(h *OperatorHeartbeat) { h.Owner.Generation = 0 }, false},
		{"no owner", "test", "start", 10, func(h *OperatorHeartbeat) { h.Owner.Nonce = "" }, false},
		{"finish not ready", "test", "finish", 10, nil, false},
		{"both ready", "test", "finish", 10, func(h *OperatorHeartbeat) { h.Readiness = "both" }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			x := h
			if tc.change != nil {
				tc.change(&x)
			}
			if got := x.Available(tc.build, tc.kind, tc.now); got != tc.want {
				t.Fatalf("available = %t", got)
			}
		})
	}
}

func TestOperatorCatalogueAdmission(t *testing.T) {
	c := openTestCredentials(t)
	ctx := testCtx(t)
	a := OperatorHeartbeat{Owner: OperatorOwner{Nonce: "first"}, Protocol: OperatorRunnerProtocol, Build: "test", SeenAt: 11, Readiness: "start", Catalogue: []byte("one")}
	publish := func(h OperatorHeartbeat, gen, catalogueGen int64) OperatorHeartbeat {
		t.Helper()
		gotGen, err := c.PutOperatorHeartbeat(ctx, h)
		got, readErr := c.OperatorHeartbeat(ctx)
		h.Owner.Generation, h.CatalogueGeneration = gen, catalogueGen
		if gotGen != gen || err != nil || readErr != nil || !reflect.DeepEqual(got, h) {
			t.Fatalf("heartbeat %+v, generation %d: %v %v; want %+v", got, gotGen, err, readErr, h)
		}
		return h
	}
	first := publish(a, 1, 1)
	a.Owner.Generation, a.SeenAt, a.Catalogue = 1, 11, []byte("two")
	publish(a, 1, 2) // same-second tick is a valid refresh
	a.SeenAt++
	second := publish(a, 1, 3)
	if string(first.Catalogue) != "one" {
		t.Fatal("published generation changed in-place")
	}
	for _, damage := range []func(*OperatorHeartbeat){
		func(h *OperatorHeartbeat) { h.Owner.Nonce = "stale" },
		func(h *OperatorHeartbeat) { h.Owner.Generation++ },
		func(h *OperatorHeartbeat) { h.Protocol++ },
		func(h *OperatorHeartbeat) { h.Build = "other" },
		func(h *OperatorHeartbeat) { h.SeenAt-- },
	} {
		x := second
		damage(&x)
		if _, err := c.PutOperatorHeartbeat(ctx, x); !errors.Is(err, ErrOperatorRunner) {
			t.Fatalf("stale update admitted: %+v %v", x, err)
		}
		got, err := c.OperatorHeartbeat(ctx)
		if err != nil || !reflect.DeepEqual(got, second) {
			t.Fatalf("rejected write changed row: %+v %v", got, err)
		}
	}
	a.Owner.Nonce, a.Owner.Generation = "restart", 0
	publish(a, 2, 4)
}

func TestOperatorFinishReadiness(t *testing.T) {
	c := openTestCredentials(t)
	cas, runner := operatorClaimFixture(t, c)
	ctx := testCtx(t)
	operatorExec(t, c, `UPDATE operator_actions SET kind='finish'`)
	before := getOperator(t, c, cas.ID, 1590)
	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); !errors.Is(err, ErrOperatorRunner) {
		t.Fatalf("finish claim: %v", err)
	}
	if got := getOperator(t, c, cas.ID, 1590); !reflect.DeepEqual(got, before) {
		t.Fatalf("blocked claim changed row: %+v", got)
	}
	operatorExec(t, c, `UPDATE operator_actions SET kind='start'`)
	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
		t.Fatalf("start claim: %v", err)
	}
}

func TestOperatorCancelInflight(t *testing.T) {
	c := openTestCredentials(t)
	cas, runner := operatorClaimFixture(t, c)
	ctx := testCtx(t)
	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
		t.Fatal(err)
	}
	operatorExec(t, c, `INSERT INTO operator_action_steps
		SELECT 'launch','attempt',id,'reserved',1,owner,1590,NULL FROM operator_actions`)
	refuse := func(callCtx context.Context, id string, now int64, want error) {
		t.Helper()
		before := getOperator(t, c, id, now)
		err := c.CancelUnstartedOperatorAction(callCtx, id, now)
		if got := getOperator(t, c, id, now); !errors.Is(err, want) || !reflect.DeepEqual(got, before) {
			t.Fatalf("cancel refusal changed %s: %+v, %v", id, got, err)
		}
	}
	refuse(ctx, cas.ID, 1591, ErrOperatorConflict)
	operatorExec(t, c, "UPDATE operator_actions SET state='confirmed', owner=NULL")
	refuse(ctx, cas.ID, 1591, ErrOperatorConflict)
	var owner string
	if err := c.db.QueryRow(`SELECT owner FROM operator_action_steps WHERE kind='launch' AND phase='reserved'`).Scan(&owner); err != nil || owner != runner.Nonce {
		t.Fatalf("reservation lost: %s %v", owner, err)
	}
	pending := putOperator(t, c, operatorRequest(t, c, 2, "pending", "start"), 1591)
	refuse(ctx, pending.ID, pending.ExpiresAt, ErrOperatorConflict)
	refuse(context.Background(), pending.ID, 1591, ErrNoDeadline)
	if err := c.CancelUnstartedOperatorAction(ctx, "missing", 1591); !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("cancel of missing action: %v", err)
	}
	if err := c.CancelUnstartedOperatorAction(ctx, pending.ID, 1591); err != nil {
		t.Fatal(err)
	}
	if got := getOperator(t, c, pending.ID, 1591); got.State != "cancelled" || got.Generation != pending.Generation+1 || got.TerminalAt != 1591 {
		t.Fatalf("cancellation readback: %+v", got)
	}
	refuse(ctx, pending.ID, 1591, ErrOperatorConflict)
}

func TestOperatorCancelConfirmed(t *testing.T) {
	c := openTestCredentials(t)
	cas, _ := operatorClaimFixture(t, c)
	const now = 1591
	before := getOperator(t, c, cas.ID, now)
	if before.State != "confirmed" || before.Owner != "" || before.ExpiresAt <= now {
		t.Fatalf("expected confirmed, unclaimed and unexpired action: %+v", before)
	}
	if err := c.CancelUnstartedOperatorAction(testCtx(t), cas.ID, now); err != nil {
		t.Fatal(err)
	}
	if got := getOperator(t, c, cas.ID, now); got.State != "cancelled" || got.Generation != before.Generation+1 || got.TerminalAt != now {
		t.Fatalf("cancellation readback: %+v", got)
	}
}
