package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type budgetResolver struct {
	inner  authority.Resolver
	calls  atomic.Int32
	held   atomic.Int64
	budget atomic.Int64
}

func (r *budgetResolver) ResolveContext(ctx context.Context, header string, now int64) (authority.ResolveOutcome, error) {
	if r.calls.Add(1) == 1 {
		d, ok := ctx.Deadline()
		if !ok {
			return authority.ResolveOutcome{}, errors.New("resolver has no deadline")
		}
		r.budget.Store(int64(time.Until(d)))
		start := time.Now()
		timer := time.NewTimer(2800 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return authority.ResolveOutcome{}, ctx.Err()
		case <-timer.C:
		}
		r.held.Store(int64(time.Since(start)))
	}
	return r.inner.ResolveContext(ctx, header, now)
}

type budgetHeads struct {
	inner  projection.HeadReader
	calls  atomic.Int32
	held   atomic.Int64
	budget atomic.Int64
}

func (r *budgetHeads) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	if r.calls.Add(1) == 1 {
		d, ok := ctx.Deadline()
		if !ok {
			return hashref.HashRef{}, false, errors.New("tools has no deadline")
		}
		r.budget.Store(int64(time.Until(d)))
		start := time.Now()
		timer := time.NewTimer(9800 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return hashref.HashRef{}, false, ctx.Err()
		case <-timer.C:
		}
		r.held.Store(int64(time.Since(start)))
	}
	return r.inner.GetRegistryHead(ctx, name)
}

type budgetStore struct {
	*store.Store
	mu            sync.Mutex
	gets, intents []string
	commits       int
}

func (s *budgetStore) GetReceipt(ctx context.Context, id string) (store.Receipt, bool, error) {
	s.mu.Lock()
	s.gets = append(s.gets, id)
	s.mu.Unlock()
	return s.Store.GetReceipt(ctx, id)
}
func (s *budgetStore) AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error) {
	s.mu.Lock()
	s.intents = append(s.intents, id)
	s.mu.Unlock()
	return s.Store.AppendIntent(ctx, id, intent)
}
func (s *budgetStore) Commit(ctx context.Context, c store.Commit) error {
	err := s.Store.Commit(ctx, c)
	if err == nil {
		s.mu.Lock()
		s.commits++
		s.mu.Unlock()
	}
	return err
}

type budgetRunner struct {
	batch     bool
	calls     atomic.Int32
	recovered atomic.Bool
	returned  chan struct{}
	remaining atomic.Int64
}

func (r *budgetRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	n := r.calls.Add(1)
	if r.recovered.Load() {
		return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
	}
	if r.batch && n == 1 {
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return capsule.Result{}, ctx.Err()
		case <-timer.C:
		}
		return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
	}
	if d, ok := ctx.Deadline(); ok {
		r.remaining.Store(int64(time.Until(d)))
	}
	<-ctx.Done()
	close(r.returned)
	return capsule.Result{}, ctx.Err()
}

type budgetWrite struct {
	http.ResponseWriter
	written *atomic.Int64
	start   time.Time
}

func (w budgetWrite) WriteHeader(status int) {
	w.written.CompareAndSwap(0, int64(time.Since(w.start)))
	w.ResponseWriter.WriteHeader(status)
}
func (w budgetWrite) Write(b []byte) (int, error) {
	w.written.CompareAndSwap(0, int64(time.Since(w.start)))
	return w.ResponseWriter.Write(b)
}
func TestMCPPostBudgetProductionConstants(t *testing.T) {
	for _, batch := range []bool{false, true} {
		name := "single"
		if batch {
			name = "batch"
		}
		t.Run(name, func(t *testing.T) {
			d, _, _ := invocationDaemon(t)
			tok := mintSessionGrants(t, d, "ep-budget", "world.apply")
			resolver := &budgetResolver{inner: d.resolver}
			heads := &budgetHeads{inner: d.reads}
			record := &budgetStore{Store: d.store}
			runner := &budgetRunner{batch: batch, returned: make(chan struct{})}
			coord, err := coordinator.New(coordinator.Config{Store: record, Runner: runner, Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
				return broker.OpenBinder(d.store, ep, caps)
			}, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
			if err != nil {
				t.Fatal(err)
			}
			d.projection, err = projection.New(projection.Config{Resolver: resolver, Reader: transitionreg.NewReader(d.store), Heads: heads, Deny: writeSessionDenial, Fail: writeAPIError, ErrorLog: io.Discard, Agent: protocol.AgentInfo{Name: "ailang-worldd", Version: Version}, MaxWait: readDeadline, InvokeWait: invokeDeadline, CredentialWait: credentialBudget, CallbackTimeout: invokeDeadline, MaxCallbacks: 8, WriteWait: writeTimeout, Coordinator: coord})
			if err != nil {
				t.Fatal(err)
			}
			var written atomic.Int64
			start := time.Now()
			handler := d.Handler()
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(budgetWrite{w, &written, start}, r) }))
			server.Config = newServer(server.Config.Handler)
			server.Start()
			defer server.Close()
			client := &http.Client{Timeout: 40 * time.Second}
			body := `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{"x":1}}}`
			if batch {
				body = "[" + body + "," + body + "," + body + "]"
			}
			request, err := http.NewRequest("POST", server.URL+"/mcp/", strings.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Authorization", "Bearer "+tok)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json, text/event-stream")
			response, callerr := client.Do(request)
			elapsed := time.Since(start)
			var payload []byte
			if response != nil {
				payload, err = io.ReadAll(response.Body)
				response.Body.Close()
			}
			select {
			case <-runner.returned:
			case <-time.After(time.Second):
				t.Fatal("timed callback did not return within finite teardown")
			}
			calls, admissions, commits := int32(1), int32(2), 0
			if batch {
				calls = 2
				admissions = 3
				commits = 1
			}
			if runner.calls.Load() != calls || heads.calls.Load() != admissions || record.commits != commits || len(record.intents) != commits || len(record.gets) != int(calls) {
				t.Fatalf("stage accounting calls=%d admission=%d commits=%d intents=%d gets=%d remaining=%s", runner.calls.Load(), heads.calls.Load(), record.commits, len(record.intents), len(record.gets), time.Duration(runner.remaining.Load()))
			}
			for i, id := range record.gets {
				rc, ok, err := d.store.GetReceipt(boundedTestContext(t), id)
				if err != nil {
					t.Fatal(err)
				}
				want := batch && i == 0
				if ok != want || (want && rc.State != store.ReceiptResolved) {
					t.Fatalf("item%d receipt=%+v/%t expected=%t", i, rc, ok, want)
				}
			}
			t.Logf("production clocks3/10/20/30 elapsed=%s write=%s resolverHeld=%s toolsHeld=%s N/K/J=%d/%d/%d admissions=%d commits=%d", elapsed, time.Duration(written.Load()), time.Duration(resolver.held.Load()), time.Duration(heads.held.Load()), map[bool]int{false: 1, true: 3}[batch], calls, calls, heads.calls.Load(), record.commits)
			if time.Duration(resolver.held.Load()) < 2750*time.Millisecond || time.Duration(resolver.held.Load()) >= credentialBudget || time.Duration(heads.held.Load()) < 9750*time.Millisecond || time.Duration(heads.held.Load()) >= readDeadline || time.Duration(resolver.budget.Load()) > credentialBudget || time.Duration(heads.budget.Load()) > readDeadline {
				t.Fatalf("inner bounds held=%s/%s budgets=%s/%s", time.Duration(resolver.held.Load()), time.Duration(heads.held.Load()), time.Duration(resolver.budget.Load()), time.Duration(heads.budget.Load()))
			}
			if callerr != nil || err != nil || elapsed > invokeDeadline+time.Second || elapsed < invokeDeadline-100*time.Millisecond || time.Duration(written.Load()) >= writeTimeout || written.Load() == 0 {
				t.Fatalf("aggregate/wire elapsed=%s written=%s err=%v/%v payload=%s", elapsed, time.Duration(written.Load()), callerr, err, payload)
			}
			if time.Duration(runner.remaining.Load()) >= 8*time.Second {
				t.Fatalf("callback did not inherit aggregate remainder: %s", time.Duration(runner.remaining.Load()))
			}
			var wire struct {
				ID    any
				Error struct {
					Code    int
					Message string
				}
			}
			if json.Unmarshal(payload, &wire) != nil || wire.Error.Code != -32603 || wire.Error.Message != "host callback timed out" || response.StatusCode != 200 {
				t.Fatalf("frozen timeout envelope=%s", payload)
			}
			if batch && wire.ID != nil {
				t.Fatalf("batch failure ID=%v", wire.ID)
			}

			runner.recovered.Store(true)
			request, err = http.NewRequest("POST", server.URL+"/mcp/", strings.NewReader(`{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`))
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Authorization", "Bearer "+tok)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json, text/event-stream")
			response, err = client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			payload, err = io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil || !strings.Contains(string(payload), "event: message") || strings.Contains(string(payload), `"error"`) {
				t.Fatalf("slot recovery=%s err=%v", payload, err)
			}
		})
	}
}

// A slow body can outlive the aggregate context because the upstream decoder is
// synchronous. The frozen transport deadline bounds that residual exposure.
func TestMCPLateBodySocketResidual(t *testing.T) {
	d, _, _ := invocationDaemon(t)
	tok := mintSessionGrants(t, d, "ep-late", "world.apply")
	handler := d.Handler()
	done := make(chan struct{})
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { defer close(done); handler.ServeHTTP(w, r) }))
	server.Config = newServer(server.Config.Handler)
	server.Start()
	defer server.Close()
	conn, err := net.DialTimeout("tcp", server.Listener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(35 * time.Second)); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = fmt.Fprintf(conn, "POST /mcp/ HTTP/1.1\r\nHost: localhost\r\nAuthorization: Bearer %s\r\nContent-Type: application/json\r\nAccept: application/json, text/event-stream\r\nContent-Length: 100\r\n\r\n{", tok)
	if err != nil {
		t.Fatal(err)
	}
	data, readerr := io.ReadAll(conn)
	elapsed := time.Since(start)
	if readerr != nil && !errors.Is(readerr, io.EOF) {
		t.Fatalf("read failed before server closure: %v", readerr)
	}
	if elapsed < writeTimeout-time.Second || elapsed > writeTimeout+2*time.Second {
		t.Fatalf("frozen socket elapsed=%s data=%s", elapsed, data)
	}
	if strings.Contains(string(data), "event: message") {
		t.Fatalf("late body produced success: %s", data)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler failed finite teardown")
	}
	if _, ok, err := d.store.GetLogEntry(boundedTestContext(t), 1); err != nil || ok {
		t.Fatalf("late body committed: %t/%v", ok, err)
	}
	t.Logf("late-body residual elapsed=%s bytes=%d: no universal20s wall-clock claim", elapsed, len(data))
}
