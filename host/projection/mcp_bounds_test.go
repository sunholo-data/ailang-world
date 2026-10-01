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
	block    atomic.Bool
	returned chan struct{}
	runs     atomic.Int32
}

func (r *cooperativeRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	r.runs.Add(1)
	if r.block.Load() {
		<-ctx.Done()
		close(r.returned)
		return capsule.Result{}, ctx.Err()
	}
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
	cfg.Resolver = delayedResolver{inner: cfg.Resolver, delay: 30 * time.Millisecond}
	t.Cleanup(func() {
		select {
		case <-runner.returned:
		case <-time.After(200 * time.Millisecond):
			t.Error("callback survived teardown")
		}
	})
	cfg.InvokeWait = 80 * time.Millisecond
	cfg.CallbackTimeout = time.Second
	cfg.MaxCallbacks = 1
	h := mustHandler(t, cfg)
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if !strings.Contains(w.Body.String(), `"host callback timed out"`) {
		t.Fatalf("deadline wire=%s", w.Body)
	}
	select {
	case <-runner.returned:
	case <-time.After(20 * time.Millisecond):
		t.Fatal("callback did not cooperate with aggregate deadline")
	}
	runner.block.Store(false)
	w = postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if len(mcpPayload(t, w)["error"]) != 0 || runner.runs.Load() != 2 {
		t.Fatalf("slot recovery runs=%d wire=%s", runner.runs.Load(), w.Body)
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
