package coordinator

// Row 153 M1: a phase that exhausts its OWN budget is a typed
// *PhaseTimeoutError, distinct from the caller's deadline.

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
)

// AC1.1/AC1.2 (MUT-RELABEL, MUT-NOWRAP): the four classification rows. No
// sleeps: the function is pure over its inputs.
func TestClassifyPhaseErrMatrix(t *testing.T) {
	const budget = 2 * time.Second
	if got := classifyPhaseErr("plan", budget, time.Second, nil, nil); got != nil {
		t.Fatalf("(nil,nil) = %v, want nil", got)
	}

	got := classifyPhaseErr("plan", budget, 2010*time.Millisecond, nil, context.DeadlineExceeded)
	var pt *PhaseTimeoutError
	if !errors.As(got, &pt) || pt.Phase != "plan" || pt.Budget != budget || pt.Elapsed != 2010*time.Millisecond {
		t.Fatalf("(nil,DE) = %T %v, want *PhaseTimeoutError{plan, 2s, 2.01s}", got, got)
	}
	if !errors.Is(got, context.DeadlineExceeded) {
		t.Fatalf("(nil,DE) %v does not unwrap to context.DeadlineExceeded", got)
	}
	if want := "coordinator: plan phase exceeded its 2s budget after 2.01s"; got.Error() != want {
		t.Fatalf("(nil,DE) text = %q, want %q", got.Error(), want)
	}

	got = classifyPhaseErr("plan", budget, time.Second, context.DeadlineExceeded, context.DeadlineExceeded)
	if errors.As(got, &pt) {
		t.Fatalf("(DE,DE) = %v, want the caller's deadline, not a phase timeout", got)
	}
	if !errors.Is(got, context.DeadlineExceeded) || !strings.Contains(got.Error(), "execute plan phase") {
		t.Fatalf("(DE,DE) = %v, want the unchanged 'execute plan phase' deadline", got)
	}

	got = classifyPhaseErr("finish", budget, time.Second, context.Canceled, context.Canceled)
	if errors.As(got, &pt) || !errors.Is(got, context.Canceled) {
		t.Fatalf("(Canceled,Canceled) = %T %v, want a plain cancel", got, got)
	}
}

// firstPlanBlocks blocks the FIRST plan on ctx.Done() and answers every
// later phase from the embedded phaseRunner.
type firstPlanBlocks struct {
	*phaseRunner
	plans atomic.Int32
}

func (f *firstPlanBlocks) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	res, err := f.phaseRunner.RunContext(ctx, e)
	if err != nil {
		return res, err
	}
	if f.plans.Add(1) == 1 {
		<-ctx.Done()
		return capsule.Result{}, ctx.Err()
	}
	return res, nil
}

// AC1.4 (V28, MUT-EARLY-INTENT): a plan timeout runs nothing and records
// nothing, so a resend of the SAME task id is a fresh invocation that commits
// exactly once.
func TestA2APlanTimeoutRetrySameTaskID(t *testing.T) {
	r, w := fxRig(t, false)
	h := &probe{}
	runner := &firstPlanBlocks{phaseRunner: &phaseRunner{plan: readPlan(false)}}
	c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
	c.planBudget = 100 * time.Millisecond

	_, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	var pt *PhaseTimeoutError
	if !errors.As(err, &pt) || pt.Phase != "plan" {
		t.Fatalf("dispatch #1 err = %T %v, want *PhaseTimeoutError{plan}", err, err)
	}
	if h.n.Load() != 0 || r.spend() != 0 {
		t.Fatalf("after the timeout: dispatches %d spend %d, want 0/0", h.n.Load(), r.spend())
	}
	r.assertUntouched(w.Ref, "t1")

	res, err := c.Dispatch(boundedTestContext(t), r.call("fx", "t1", fxGrants(5)))
	if err != nil {
		t.Fatalf("dispatch #2 (same task id) = %v, want success", err)
	}
	if h.n.Load() != 1 || r.spend() != 1 {
		t.Fatalf("after the resend: dispatches %d spend %d, want 1/1", h.n.Load(), r.spend())
	}
	if res.EntryIndex != w.Revision+1 {
		t.Fatalf("entry index %d, want exactly one log entry (%d)", res.EntryIndex, w.Revision+1)
	}
	if got := decodeWorld(t, res.OutputBytes).World.Effects; len(got) != 1 || got[0].Status != StatusOK {
		t.Fatalf("effects = %+v, want one ok effect record", got)
	}
}
