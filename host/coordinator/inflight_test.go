package coordinator

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/store"
)

type flightRunner func(context.Context, capsule.Entry) (capsule.Result, error)

func (f flightRunner) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	return f(ctx, e)
}

type receiptErrorStore struct {
	*store.Store
	fail atomic.Bool
}

func (s *receiptErrorStore) GetReceipt(id string) (store.Receipt, bool, error) {
	if s.fail.Swap(false) {
		return store.Receipt{}, false, errors.New("receipt failed")
	}
	return s.Store.GetReceipt(id)
}

func flightRig(t *testing.T) (*rig, Call) {
	t.Helper()
	r := newRig(t, false)
	r.genesis()
	r.describe("plain", r.source(echoSrc), "Invoke")
	r.publish()
	return r, r.call("plain", "same", grants)
}

func TestInFlightGuard(t *testing.T) {
	t.Run("R17_concurrent_same_id", func(t *testing.T) {
		r, call := flightRig(t)
		entered := make(chan struct{})
		release := make(chan struct{})
		done := make(chan error, 1)
		var runs atomic.Int32
		runner := flightRunner(func(context.Context, capsule.Entry) (capsule.Result, error) {
			runs.Add(1)
			close(entered)
			<-release
			return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
		})
		c := r.coordinator(r.st, runner)
		go func() { _, err := c.Dispatch(context.Background(), call); done <- err }()
		<-entered
		_, err := c.Dispatch(context.Background(), call)
		var inFlight *InFlightError
		if !errors.As(err, &inFlight) || inFlight.InvocationID != InvocationID(call.EpisodeID, call.TaskID) {
			t.Fatalf("retry = %T %v", err, err)
		}
		if got := runs.Load(); got != 1 {
			t.Fatalf("runs=%d, want 1", got)
		}
		close(release)
		if err := <-done; err != nil {
			t.Fatal(err)
		}
		again, err := c.Dispatch(context.Background(), call)
		if err != nil || !again.Reconciled || runs.Load() != 1 {
			t.Fatalf("reconciled=%+v err=%v runs=%d", again, err, runs.Load())
		}
	})
	t.Run("release_after_error", func(t *testing.T) {
		r, call := flightRig(t)
		var runs atomic.Int32
		c := r.coordinator(r.st, flightRunner(func(context.Context, capsule.Entry) (capsule.Result, error) {
			runs.Add(1)
			return capsule.Result{}, errors.New("run failed")
		}))
		for i := 0; i < 2; i++ {
			_, err := c.Dispatch(context.Background(), call)
			var inFlight *InFlightError
			if err == nil || errors.As(err, &inFlight) {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		if runs.Load() != 2 {
			t.Fatalf("runs=%d, want 2", runs.Load())
		}
	})
	t.Run("release_after_panic", func(t *testing.T) {
		r, call := flightRig(t)
		var runs atomic.Int32
		c := r.coordinator(r.st, flightRunner(func(context.Context, capsule.Entry) (capsule.Result, error) {
			if runs.Add(1) == 1 {
				panic("runner panic")
			}
			return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
		}))
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			_, _ = c.Dispatch(context.Background(), call)
		}()
		if _, err := c.Dispatch(context.Background(), call); err != nil {
			t.Fatalf("retry after panic: %v", err)
		}
		if runs.Load() != 2 {
			t.Fatalf("runs=%d, want 2", runs.Load())
		}
	})
	t.Run("release_after_receipt_error", func(t *testing.T) {
		r, call := flightRig(t)
		st := &receiptErrorStore{Store: r.st}
		st.fail.Store(true)
		c := r.coordinator(st, fakeRunner{stdout: `{"ok":true}`})
		if _, err := c.Dispatch(context.Background(), call); err == nil {
			t.Fatal("receipt error was ignored")
		}
		if _, err := c.Dispatch(context.Background(), call); err != nil {
			t.Fatalf("retry after receipt error: %v", err)
		}
	})
}
