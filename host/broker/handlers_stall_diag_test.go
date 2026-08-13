package broker

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type storeCallTiming struct {
	op     string
	offset time.Duration
}

type stallDiagnosis struct {
	err          error
	elapsed      time.Duration
	executeEnter time.Duration
	executeExit  time.Duration
	storeCalls   []storeCallTiming
}

func (d stallDiagnosis) preHandlerWindow() time.Duration {
	return d.executeEnter
}

func (d stallDiagnosis) executeWindow() time.Duration {
	return d.executeExit - d.executeEnter
}

func (d stallDiagnosis) postHandlerWindow() time.Duration {
	return d.elapsed - d.executeExit
}

func (d stallDiagnosis) String() string {
	return fmt.Sprintf("elapsed=%s pre-handler=%s Execute=%s post-handler=%s error=%+v cause=%+v store-calls=%v",
		d.elapsed, d.preHandlerWindow(), d.executeWindow(), d.postHandlerWindow(), d.err,
		errors.Unwrap(d.err), d.storeCalls)
}

type timingHandler struct {
	handler     Handler
	invokeStart time.Time
	enter       time.Duration
	exit        time.Duration
}

func (h *timingHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	h.enter = time.Since(h.invokeStart)
	result, err := h.handler.Execute(ctx, req, payload)
	h.exit = time.Since(h.invokeStart)
	return result, err
}

func invokeWithStallDiagnosis(
	t *testing.T,
	session *Session,
	store *handlerRecordingStore,
	timing *timingHandler,
	req EffectRequest,
	payload []byte,
) stallDiagnosis {
	t.Helper()
	start := time.Now()
	timing.invokeStart = start
	var calls []storeCallTiming
	priorStoreHook := store.onStoreCall
	store.onStoreCall = func(op string) {
		calls = append(calls, storeCallTiming{op: op, offset: time.Since(start)})
		if priorStoreHook != nil {
			priorStoreHook(op)
		}
	}
	_, _, invokeErr := session.Invoke(context.Background(), req, payload)
	return stallDiagnosis{
		err:          invokeErr,
		elapsed:      time.Since(start),
		executeEnter: timing.enter,
		executeExit:  timing.exit,
		storeCalls:   calls,
	}
}

type warmUpRunner func(ctx context.Context, path string) error

type warmUpCall struct {
	count       int
	hadDeadline bool
	timeout     time.Duration
}

func warmUpFixture(t *testing.T, path string, run warmUpRunner) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := run(ctx, path); err != nil {
		t.Fatalf("warm fixture %s: %v", path, err)
	}
}

func execWarmUpRunner(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, path)
	cmd.Env = []string{"W16_WARM=1"}
	return cmd.Run()
}

func TestWarmUpRunsExactlyOnceUnderABoundedContext(t *testing.T) {
	var call warmUpCall
	warmUpFixture(t, "fixture", func(ctx context.Context, path string) error {
		call.count++
		deadline, ok := ctx.Deadline()
		call.hadDeadline = ok
		if ok {
			call.timeout = time.Until(deadline)
		}
		if path != "fixture" {
			return fmt.Errorf("path = %q, want fixture", path)
		}
		return nil
	})
	if call.count != 1 {
		t.Fatalf("warm-up call count = %d, want 1", call.count)
	}
	if !call.hadDeadline {
		t.Fatal("warm-up context had no deadline")
	}
	if call.timeout <= 0 {
		t.Fatalf("warm-up timeout = %s, want positive", call.timeout)
	}
}

func TestStallDiagnosisAttributesHandlerWindow(t *testing.T) {
	scope := t.TempDir()
	handler := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		time.Sleep(1500 * time.Millisecond)
		return nil, fmt.Errorf("stub: %w", ErrHandlerTimeout)
	})
	timing := &timingHandler{handler: handler}
	session, recording := handlerSession(t, EffectGitCommit, scope, timing)
	diagnosis := invokeWithStallDiagnosis(t, session, recording, timing,
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, []byte("handler stall"))
	if diagnosis.executeWindow() < 1200*time.Millisecond {
		t.Fatalf("handler stall was not attributed to Execute: %s", diagnosis)
	}
	if !strings.Contains(diagnosis.String(), "stub: broker: handler subprocess timed out") {
		t.Fatalf("handler diagnosis omitted error text: %s", diagnosis)
	}
}

func TestStallDiagnosisAttributesStoreWindow(t *testing.T) {
	scope := t.TempDir()
	handler := HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		return nil, fmt.Errorf("stub: %w", ErrHandlerTimeout)
	})
	timing := &timingHandler{handler: handler}
	session, recording := handlerSession(t, EffectGitCommit, scope, timing)
	recording.onStoreCall = func(op string) {
		if op == "AppendEffectOutcome" {
			time.Sleep(1500 * time.Millisecond)
		}
	}
	diagnosis := invokeWithStallDiagnosis(t, session, recording, timing,
		EffectRequest{Effect: EffectGitCommit, Scope: scope, Cost: 2, Now: 1}, []byte("store stall"))
	if len(diagnosis.storeCalls) == 0 || diagnosis.postHandlerWindow() < 1200*time.Millisecond {
		t.Fatalf("store stall was not attributed to post-handler work: %s", diagnosis)
	}
}
