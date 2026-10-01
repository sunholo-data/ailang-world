package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang/serveapi/protocol"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const invokeEchoSource = "module quickstart/echo\n\nexport func main(input: string) -> string {\n  \"{\\\"echo\\\":${input}}\"\n}\n"

func invocationDaemon(t *testing.T) (*Daemon, string, hashref.HashRef) {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Skip("AILANG_BIN is required for the real interpreter E2E")
	}
	db := filepath.Join(t.TempDir(), "world.db")
	d, err := New(boundedTestContext(t), Config{DBPath: db, BindHost: DefaultBindHost, AilangBin: bin})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Error(err)
		}
	})
	interp := daemonInterpreterRef(t, d)
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")), SemanticID: "test/genesis", Provenance: "invoke_e2e_test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	err = d.store.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: interp, PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "invoke_e2e_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}})
	if err != nil {
		t.Fatalf("genesis: %v", err)
	}
	source, err := canon.Source([]byte(invokeEchoSource))
	if err != nil {
		t.Fatal(err)
	}
	src := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")), SemanticID: "world/transition-source/v1", Provenance: "invoke_e2e_test", Payload: source}
	if err := d.store.PutObject(boundedTestContext(t), src); err != nil {
		t.Fatal(err)
	}
	desc := transitionreg.Descriptor{ID: "tools.echo", TransitionFn: src.Hash, Interpreter: interp, SemanticsEpoch: 1, InputSchema: []byte(`{"type":"object"}`), OutputSchema: []byte(`{"type":"object"}`), Access: transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1}, Title: "Echo", Description: "E2E echo"}
	if _, err := transitionreg.NewPublisher(d.store, archive.New(db)).PublishSet(boundedTestContext(t), []transitionreg.Change{{ID: desc.ID, Descriptor: &desc}}); err != nil {
		t.Fatalf("PublishSet: %v", err)
	}
	return d, db, interp
}

func invokeBody(task string) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":9,"method":"tasks/send","params":{"id":%q,"metadata":{"skill_id":"tools.echo"},"message":{"role":"user","parts":[{"type":"data","data":{"msg":"hi"}}]}}}`, task)
}

func parseInvokeTask(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var wire map[string]any
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["id"] != float64(9) || wire["error"] != nil {
		t.Fatalf("A2A response: %s", body)
	}
	task, ok := wire["result"].(map[string]any)
	if !ok {
		t.Fatalf("missing result: %s", body)
	}
	return task
}

func TestA2AInvokesPublishedTransition(t *testing.T) {
	d, db, interp := invocationDaemon(t)
	var executions atomic.Int32
	replaceInvocation(t, d, db, d.store, countingInvocationRunner{real: capsule.New(archive.New(db), capsule.Config{}), count: &executions}, invokeDeadline)
	token := mintSessionGrants(t, d, "ep-invoke", "world.apply")
	call := func() map[string]any {
		rec := requestRecorderAuth(t, d, "Bearer "+token, http.MethodPost, "/a2a/", strings.NewReader(invokeBody("task-1")))
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d: %s", rec.Code, rec.Body)
		}
		return parseInvokeTask(t, rec.Body.Bytes())
	}
	task := call()
	if task["id"] != "task-1" || task["status"].(map[string]any)["state"] != "completed" {
		t.Fatalf("task: %+v", task)
	}
	meta := task["metadata"].(map[string]any)
	if meta["entry_index"] != float64(1) || meta["invocation_id"] == "" || meta["world_ref"] == "" {
		t.Fatalf("metadata: %+v", meta)
	}
	parts := task["artifacts"].([]any)[0].(map[string]any)["parts"].([]any)
	if len(parts) != 1 || parts[0].(map[string]any)["type"] != "data" {
		t.Fatalf("parts: %+v", parts)
	}
	output := parts[0].(map[string]any)["data"].(map[string]any)
	if output["echo"].(map[string]any)["msg"] != "hi" {
		t.Fatalf("output: %+v", output)
	}
	log := requestRecorderAuth(t, d, "Bearer "+token, http.MethodGet, "/v1/log/1", nil)
	if log.Code != http.StatusOK {
		t.Fatalf("GET log: %d %s", log.Code, log.Body)
	}
	var entry logEntryResponse
	if err := json.Unmarshal(log.Body.Bytes(), &entry); err != nil {
		t.Fatal(err)
	}
	if entry.Header.EntryIndex != 1 || entry.Header.Interpreter != interp.String() {
		t.Fatalf("entry: %+v", entry)
	}
	fn, err := hashref.Parse(entry.Header.TransitionFn)
	if err != nil {
		t.Fatal(err)
	}
	src, ok, err := d.store.GetObject(boundedTestContext(t), fn)
	if err != nil || !ok {
		t.Fatalf("source: %v %v", ok, err)
	}
	input, err := json.Marshal(map[string]any{"msg": "hi"})
	if err != nil {
		t.Fatal(err)
	}
	arg, err := json.Marshal(string(input))
	if err != nil {
		t.Fatal(err)
	}
	replay, err := capsule.New(archive.New(db), capsule.Config{}).RunContext(boundedTestContext(t), capsule.Entry{Interpreter: interp, Source: src.Payload, Args: arg})
	if err != nil {
		t.Fatalf("replay: %v stderr=%s", err, replay.Stderr)
	}
	want, _ := json.Marshal(output)
	if !bytes.Equal(bytes.TrimSpace(replay.Stdout), want) {
		t.Fatalf("replay bytes %q != output %q", replay.Stdout, want)
	}
	second := call()
	if got := executions.Load(); got != 1 {
		t.Fatalf("same task executed %d times, want once", got)
	}
	if second["metadata"].(map[string]any)["entry_index"] != float64(1) {
		t.Fatalf("resend: %+v", second)
	}
	if _, ok, err := d.store.GetLogEntry(boundedTestContext(t), 2); err != nil || ok {
		t.Fatalf("second entry: ok=%v err=%v", ok, err)
	}
	rc, ok, err := d.store.GetReceipt(boundedTestContext(t), meta["invocation_id"].(string))
	if err != nil || !ok || rc.State != store.ReceiptResolved {
		t.Fatalf("receipt: %+v ok=%v err=%v", rc, ok, err)
	}
}

type countingInvocationRunner struct {
	real  *capsule.Runner
	count *atomic.Int32
}

func (r countingInvocationRunner) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	r.count.Add(1)
	return r.real.RunContext(ctx, e)
}

type invocationStore struct {
	*store.Store
	committed chan struct{}
	release   <-chan struct{}
}

func (s invocationStore) Commit(ctx context.Context, c store.Commit) error {
	err := s.Store.Commit(ctx, c)
	if err == nil && s.committed != nil {
		close(s.committed)
		<-s.release
	}
	return err
}

type invocationRunner struct {
	real   *capsule.Runner
	exited chan struct{}
}

func (r invocationRunner) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	if bytes.Contains(e.Source, []byte("fib(32)")) {
		result, err := r.real.RunContext(ctx, e)
		close(r.exited)
		return result, err
	}
	return capsule.Result{Stdout: []byte("{\"echo\":{\"msg\":\"hi\"}}\n")}, nil
}

func publishSlowInvocation(t *testing.T, d *Daemon, db string, interp hashref.HashRef) {
	t.Helper()
	const slow = `module quickstart/slow
func fib(n: int) -> int { if n < 2 then n else fib(n - 1) + fib(n - 2) }
export func main(input: string) -> string { if fib(32) >= 0 then input else input }
`
	source, err := canon.Source([]byte(slow))
	if err != nil {
		t.Fatal(err)
	}
	src := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")), SemanticID: "world/transition-source/v1", Provenance: "invoke_e2e_test", Payload: source}
	if err := d.store.PutObject(boundedTestContext(t), src); err != nil {
		t.Fatal(err)
	}
	desc := transitionreg.Descriptor{ID: "tools.slow", TransitionFn: src.Hash, Interpreter: interp, SemanticsEpoch: 1,
		InputSchema: []byte(`{"type":"object"}`), OutputSchema: []byte(`{"type":"object"}`),
		Access: transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1}, Title: "Slow", Description: "deadline fixture"}
	if _, err := transitionreg.NewPublisher(d.store, archive.New(db)).PublishSet(boundedTestContext(t), []transitionreg.Change{{ID: desc.ID, Descriptor: &desc}}); err != nil {
		t.Fatalf("publish slow: %v", err)
	}
}

func replaceInvocation(t *testing.T, d *Daemon, db string, st coordinator.Store, runner coordinator.Runner, wait time.Duration) {
	t.Helper()
	coord, err := coordinator.New(coordinator.Config{Store: st, Runner: runner,
		Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(d.store, ep, caps)
		},
		Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	p, err := projection.New(projection.Config{CredentialWait: credentialBudget, CallbackTimeout: invokeDeadline, MaxCallbacks: 8, WriteWait: writeTimeout, Resolver: d.resolver, Reader: transitionreg.NewReader(d.store), Heads: d.reads,
		Deny: writeSessionDenial, Fail: writeAPIError, ErrorLog: d.errLog, Agent: protocol.AgentInfo{Name: "ailang-worldd", Version: Version},
		MaxWait: readDeadline, InvokeWait: wait, Coordinator: coord})
	if err != nil {
		t.Fatal(err)
	}
	d.projection = p
	_ = db
}

func socketCall(client *http.Client, url, token, task string) ([]byte, error) {
	return socketCallBody(client, url, token, invokeBody(task))
}

func socketCallBody(client *http.Client, url, token, body string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url+"/a2a/", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func TestA2AResendAfterWriteTimeout(t *testing.T) {
	d, db, _ := invocationDaemon(t)
	committed := make(chan struct{})
	release := make(chan struct{})
	replaceInvocation(t, d, db, invocationStore{Store: d.store, committed: committed, release: release}, capsule.New(archive.New(db), capsule.Config{}), 3*time.Second)
	server := httptest.NewUnstartedServer(d.Handler())
	server.Config.WriteTimeout = 2 * time.Second
	server.Start()
	defer server.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	token := mintSessionGrants(t, d, "ep-resend", "world.apply")
	first := make(chan error, 1)
	go func() { _, err := socketCall(client, server.URL, token, "task-resend"); first <- err }()
	select {
	case <-committed:
	case <-time.After(5 * time.Second):
		t.Fatal("commit did not land")
	}
	// Keep the response blocked beyond the server's WriteTimeout, after the
	// durable commit has already landed.
	time.Sleep(2100 * time.Millisecond)
	close(release)
	select {
	case err := <-first:
		if err == nil {
			t.Fatal("first response unexpectedly delivered after write timeout")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("write timeout did not end first call")
	}
	var body []byte
	var err error
	for deadline := time.After(3 * time.Second); ; {
		body, err = socketCall(client, server.URL, token, "task-resend")
		if err == nil && bytes.Contains(body, []byte(`"result"`)) {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("resend did not reconcile: %v %s", err, body)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	task := parseInvokeTask(t, body)
	if task["metadata"].(map[string]any)["entry_index"] != float64(1) {
		t.Fatalf("resend: %+v", task)
	}
	if _, ok, err := d.store.GetLogEntry(boundedTestContext(t), 2); err != nil || ok {
		t.Fatalf("second entry: %v %v", ok, err)
	}
}

func TestA2AInvokeDeadlineClosesSocket(t *testing.T) {
	d, db, interp := invocationDaemon(t)
	publishSlowInvocation(t, d, db, interp)
	exited := make(chan struct{})
	replaceInvocation(t, d, db, d.store, invocationRunner{real: capsule.New(archive.New(db), capsule.Config{}), exited: exited}, time.Second)
	closed := make(chan struct{}, 4)
	server := httptest.NewUnstartedServer(d.Handler())
	var mu sync.Mutex
	seen := map[net.Conn]bool{}
	server.Config.ConnState = func(c net.Conn, state http.ConnState) {
		mu.Lock()
		defer mu.Unlock()
		if state == http.StateNew {
			seen[c] = true
		}
		if state == http.StateClosed && seen[c] {
			delete(seen, c)
			select {
			case closed <- struct{}{}:
			default:
			}
		}
	}
	server.Start()
	defer server.Close()
	token := mintSessionGrants(t, d, "ep-deadline", "world.apply")
	start := time.Now()
	request := strings.Replace(invokeBody("task-stall"), "tools.echo", "tools.slow", 1)
	body, err := socketCallBody(server.Client(), server.URL, token, request)
	if err != nil {
		t.Fatalf("deadline call: %v", err)
	}
	if !bytes.Contains(body, []byte(`"code":-32603`)) || time.Since(start) > 3*time.Second {
		t.Fatalf("deadline result after %s: %s", time.Since(start), body)
	}
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("capsule runner did not exit")
	}
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("server did not close timed-out connection")
	}
	good, err := socketCall(server.Client(), server.URL, token, "task-normal")
	if err != nil {
		t.Fatal(err)
	}
	if parseInvokeTask(t, good)["status"].(map[string]any)["state"] != "completed" {
		t.Fatalf("normal: %s", good)
	}
}

func TestA2AErrorLogWiring(t *testing.T) {
	var log bytes.Buffer
	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, ErrorLog: &log})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Error(err)
		}
	})
	seedTransitionRegistry(t, d.store, "tools.echo", "alpha")
	token := mintSessionGrants(t, d, "ep-log", "alpha")
	rec := requestRecorderAuth(t, d, "Bearer "+token, http.MethodPost, "/a2a/", strings.NewReader(invokeBody("t1")))
	if !strings.Contains(rec.Body.String(), "transition invocation is not available in this daemon") {
		t.Fatalf("wire: %s", rec.Body.String())
	}
	if bytes.Count(log.Bytes(), []byte("\n")) != 1 || !strings.Contains(log.String(), "a2a refusal: tasks/send -:") {
		t.Fatalf("daemon ErrorLog=%q", log.String())
	}
}
