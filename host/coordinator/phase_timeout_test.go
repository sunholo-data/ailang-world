package coordinator

// Row 153 M1: a phase that exhausts its OWN budget is a typed
// *PhaseTimeoutError, distinct from the caller's deadline.

import (
	"context"
	"encoding/json"
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

// slowPlanRunner holds the plan phase until 500 ms before the PHASE
// context's own deadline, then answers: the slowest plan the cap admits.
type slowPlanRunner struct {
	*phaseRunner
	returned atomic.Int64 // unix nanos at which the plan returned
}

func (s *slowPlanRunner) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	var arg string
	if err := json.Unmarshal(e.Args, &arg); err != nil {
		return capsule.Result{}, err
	}
	if strings.Contains(arg, `"phase":"plan"`) {
		dl, ok := ctx.Deadline()
		if !ok {
			return capsule.Result{}, errors.New("the plan phase has no deadline")
		}
		select {
		case <-time.After(time.Until(dl) - 500*time.Millisecond):
		case <-ctx.Done():
			return capsule.Result{}, ctx.Err()
		}
		defer func() { s.returned.Store(time.Now().UnixNano()) }()
	}
	return s.phaseRunner.RunContext(ctx, e)
}

// AC2.3 (MUT-STEAL, MUT-DEFAULT): a plan that uses all but 500 ms of its cap
// still leaves the handler its full HandlerCap inside the daemon's 20 s
// invocation deadline. Default budgets, no hook: this also pins New's
// defaults. The 20 s is the literal daemon invokeDeadline (pinned by
// daemon.TestPlanPhaseBudgetIsDerived and mcp_test.go:165).
func TestSlowPlanLeavesHandlerItsCap(t *testing.T) {
	const modelInvokeDeadline = 20 * time.Second
	r, _ := fxRig(t, false)
	var handlerDeadline atomic.Int64
	h := &probe{delay: func(ctx context.Context) error {
		if dl, ok := ctx.Deadline(); ok {
			handlerDeadline.Store(dl.UnixNano())
		}
		return nil
	}}
	runner := &slowPlanRunner{phaseRunner: &phaseRunner{plan: readPlan(false)}}
	c := r.fxCoordinator(r.st, runner, broker.Registry{fxRead: h})
	ctx, cancel := context.WithTimeout(context.Background(), modelInvokeDeadline)
	defer cancel()
	if _, err := c.Dispatch(ctx, r.call("fx", "t1", fxGrants(5))); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if h.n.Load() != 1 || handlerDeadline.Load() == 0 || runner.returned.Load() == 0 {
		t.Fatalf("handler ran %d times, deadline %d, plan returned %d", h.n.Load(), handlerDeadline.Load(), runner.returned.Load())
	}
	left := time.Duration(handlerDeadline.Load() - runner.returned.Load())
	if left < HandlerCap-100*time.Millisecond {
		t.Fatalf("the handler had %s after the slow plan, want >= HandlerCap−100ms = %s: the plan cap steals the handler's time", left, HandlerCap-100*time.Millisecond)
	}
}
