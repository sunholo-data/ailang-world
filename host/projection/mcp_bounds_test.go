package projection

import (
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type cooperativeRunner struct {
	block     atomic.Bool
	returned  chan struct{}
	runs      atomic.Int32
	returns   atomic.Int32
	remaining atomic.Int64
}

func (r *cooperativeRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	n := r.runs.Add(1)
	if n == 1 {
		deadline, ok := ctx.Deadline()
		if ok {
			r.remaining.Store(int64(time.Until(deadline)))
		}
	}
	if r.block.Load() {
		<-ctx.Done()
		r.returns.Add(1)
		close(r.returned)
		return capsule.Result{}, ctx.Err()
	}
	r.returns.Add(1)
	return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
}
func TestMCPInvokeDeadlineFreesSlot(t *testing.T) {
	cfg, tok, st, _ := mcpFixture(t)
	runner := &cooperativeRunner{returned: make(chan struct{})}
	runner.block.Store(true)
	coord, err := coordinator.New(coordinator.Config{Store: st, Runner: runner, Binder: func(ep string, caps []broker.Capability) transitionreg.Binder { return broker.OpenBinder(st, ep, caps) }, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	cfg.Coordinator = coord
	resolver := &initialDelayResolver{inner: cfg.Resolver, delay: 80 * time.Millisecond}
	cfg.Resolver = resolver
	t.Cleanup(func() {
		select {
		case <-runner.returned:
		case <-time.After(3 * time.Second):
			t.Error("callback survived teardown")
		}
	})
	const initialBudget = 300 * time.Millisecond
	cfg.InvokeWait = initialBudget
	cfg.CallbackTimeout = 2 * time.Second
	cfg.MaxCallbacks = 1
	h := mustHandler(t, cfg)
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if !strings.Contains(w.Body.String(), `"host callback timed out"`) {
		t.Fatalf("deadline wire=%s", w.Body)
	}
	select {
	case <-runner.returned:
	case <-time.After(3 * time.Second):
		t.Fatal("callback did not cooperate with aggregate deadline")
	}
	if runner.runs.Load() != 1 || runner.returns.Load() != 1 {
		t.Fatalf("initial callback entered=%d returned=%d", runner.runs.Load(), runner.returns.Load())
	}
	remaining := time.Duration(runner.remaining.Load())
	if remaining <= 0 || remaining > initialBudget-40*time.Millisecond {
		t.Fatalf("callback lost aggregate context: remaining=%s initial=%s", remaining, initialBudget)
	}
	// The same handler and one-slot Runner recover with an independent live
	// test budget. Only the initial request carries the artificial delay.
	h.invokeWait = 3 * time.Second
	runner.block.Store(false)
	w = postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	t.Logf("initial aggregate=%s callback=%s entry remaining=%s; recovery=%s same-handler capacity=%d entered=%d return-path=%d resolver=%d", initialBudget, cfg.CallbackTimeout, remaining, h.invokeWait, cfg.MaxCallbacks, runner.runs.Load(), runner.returns.Load(), resolver.calls.Load())
	if len(mcpPayload(t, w)["error"]) != 0 || runner.runs.Load() != 2 || runner.returns.Load() != 2 || resolver.calls.Load() != 2 {
		t.Fatalf("slot recovery entered=%d returned=%d resolver=%d wire=%s", runner.runs.Load(), runner.returns.Load(), resolver.calls.Load(), w.Body)
	}
}

type deadlineHeads struct {
	inner     HeadReader
	remaining time.Duration
	calls     int
}

func (r *deadlineHeads) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	r.calls++
	deadline, ok := ctx.Deadline()
	if !ok {
		return hashref.HashRef{}, false, errors.New("deadline absent")
	}
	r.remaining = time.Until(deadline)
	return r.inner.GetRegistryHead(ctx, name)
}
func TestMCPToolsInnerBudget(t *testing.T) {
	cfg, _, _, _ := mcpFixture(t)
	cfg.MaxWait = 30 * time.Millisecond
	heads := &deadlineHeads{inner: cfg.Heads}
	cfg.Heads = heads
	h := mustHandler(t, cfg)
	_, err := (mcpAdapter{h}).Tools(boundedTestContext(t), &authority.SessionBinding{EpisodeID: "ep-a", Caps: []broker.Capability{liveGrant("alpha")}})
	if err != nil || heads.calls != 1 || heads.remaining <= 0 || heads.remaining > cfg.MaxWait {
		t.Fatalf("tools deadline=%s calls=%d err=%v", heads.remaining, heads.calls, err)
	}
}

type delayedResolver struct {
	inner authority.Resolver
	delay time.Duration
}

func (r delayedResolver) ResolveContext(ctx context.Context, header string, now int64) (authority.ResolveOutcome, error) {
	timer := time.NewTimer(r.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return authority.ResolveOutcome{}, ctx.Err()
	case <-timer.C:
		return r.inner.ResolveContext(ctx, header, now)
	}
}

// initialDelayResolver confines deadline stimulus to the first request.
type initialDelayResolver struct {
	inner authority.Resolver
	delay time.Duration
	calls atomic.Int32
}

func (r *initialDelayResolver) ResolveContext(ctx context.Context, header string, now int64) (authority.ResolveOutcome, error) {
	if r.calls.Add(1) == 1 {
		return (delayedResolver{inner: r.inner, delay: r.delay}).ResolveContext(ctx, header, now)
	}
	return r.inner.ResolveContext(ctx, header, now)
}
