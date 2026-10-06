package projection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"github.com/sunholo-data/ailang/serveapi/protocol/hostcall"
	"github.com/sunholo-data/ailang/serveapi/protocol/mcphttp"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const mcpCallItem = `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{"x":1}}}`

type accountingStore struct {
	*store.Store
	mu      sync.Mutex
	ids     []string
	tasks   []string
	commits int
}

func (s *accountingStore) AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error) {
	s.mu.Lock()
	s.ids = append(s.ids, id)
	s.mu.Unlock()
	return s.Store.AppendIntent(ctx, id, intent)
}
func (s *accountingStore) Commit(ctx context.Context, c store.Commit) error {
	err := s.Store.Commit(ctx, c)
	if err == nil {
		s.mu.Lock()
		s.commits++
		s.mu.Unlock()
	}
	return err
}

type batchRunner struct {
	runs, active, maximum atomic.Int32
	failAt                int32
	afterFirst            func()
}

func (r *batchRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	n := r.runs.Add(1)
	active := r.active.Add(1)
	defer r.active.Add(-1)
	for {
		old := r.maximum.Load()
		if active <= old || r.maximum.CompareAndSwap(old, active) {
			break
		}
	}
	if n == r.failAt {
		return capsule.Result{}, errors.New("batch injected capsule failure")
	}
	if n == 1 && r.afterFirst != nil {
		r.afterFirst()
	}
	return capsule.Result{Stdout: []byte(`{"ok":true}`)}, nil
}
func batchFixture(t *testing.T) (*Handler, string, *accountingStore, *countingReader, *batchRunner) {
	cfg, tok, st, _ := mcpFixture(t)
	record := &accountingStore{Store: st}
	runner := &batchRunner{}
	reader := &countingReader{inner: cfg.Reader}
	cfg.Reader = reader
	coord, err := coordinator.New(coordinator.Config{Store: record, Runner: runner, Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
		return broker.OpenBinder(st, ep, caps, nil)
	}, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	cfg.Coordinator = coord
	handler := mustHandler(t, cfg)
	handler.mintTask = func() (string, error) {
		task, err := mintMCPTask()
		if err == nil {
			record.tasks = append(record.tasks, task)
		}
		return task, err
	}
	return handler, tok, record, reader, runner
}
func batchResponses(t *testing.T, w *httptest.ResponseRecorder) []map[string]json.RawMessage {
	t.Helper()
	prefix := "event: message\ndata: "
	if w.Code != 200 || !strings.HasPrefix(w.Body.String(), prefix) {
		t.Fatalf("batch SSE=%d %s", w.Code, w.Body)
	}
	var items []map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(w.Body.String(), prefix))), &items); err != nil {
		t.Fatal(err)
	}
	return items
}
func assertReceipts(t *testing.T, s *accountingStore, want int) {
	t.Helper()
	if len(s.ids) != want || s.commits != want {
		t.Fatalf("intent/commit counts=%d/%d want %d", len(s.ids), s.commits, want)
	}
	if len(s.tasks) < want {
		t.Fatalf("minted tasks=%d want%d", len(s.tasks), want)
	}
	taskSeen := map[string]bool{}
	for i, task := range s.tasks {
		if len(task) != 64 || strings.Trim(task, "0123456789abcdef") != "" || taskSeen[task] {
			t.Fatalf("task shape/uniqueness=%v", s.tasks)
		}
		taskSeen[task] = true
		if i < len(s.ids) && coordinator.InvocationID(coordinator.SurfaceMCP, "ep-a", task) != s.ids[i] {
			t.Fatalf("task/receipt identity mismatch item%d", i)
		}
	}
	unique := map[string]bool{}
	for _, id := range s.ids {
		if unique[id] {
			t.Fatalf("duplicate invocation IDs: %v", s.ids)
		}
		unique[id] = true
		task := strings.TrimPrefix(id, "mcp:ep-a:")
		if len(task) != 64 || strings.Trim(task, "0123456789abcdef") != "" {
			t.Fatalf("task not 64hex: %q", task)
		}
		rc, ok, err := s.GetReceipt(boundedTestContext(t), id)
		if err != nil || !ok || rc.State != store.ReceiptResolved {
			t.Fatalf("receipt %q=%+v/%t/%v", id, rc, ok, err)
		}
	}
	for i := 1; i <= want; i++ {
		if _, ok, err := s.GetLogEntry(boundedTestContext(t), int64(i)); err != nil || !ok {
			t.Fatalf("entry %d absent: %t/%v", i, ok, err)
		}
	}
	if _, ok, err := s.GetLogEntry(boundedTestContext(t), int64(want+1)); err != nil || ok {
		t.Fatalf("unexpected extra entry %d: %t/%v", want+1, ok, err)
	}
}
func TestMCPBatchAccounting(t *testing.T) {
	t.Run("repeated_rpc_id", func(t *testing.T) {
		h, tok, s, reader, runner := batchFixture(t)
		w := postMCP(t, h, "Bearer "+tok, "["+mcpCallItem+","+mcpCallItem+"]")
		responses := batchResponses(t, w)
		if len(responses) != 2 || string(responses[0]["id"]) != "7" || string(responses[1]["id"]) != "7" {
			t.Fatalf("batch ids=%v", responses)
		}
		if reader.calls != 3 || runner.runs.Load() != 2 || len(s.tasks) != 2 || runner.maximum.Load() != 1 {
			t.Fatalf("K+1 admissions/runs/overlap=%d/%d/%d", reader.calls, runner.runs.Load(), runner.maximum.Load())
		}
		assertReceipts(t, s, 2)
	})
	t.Run("removed_between_items", func(t *testing.T) {
		h, tok, s, reader, runner := batchFixture(t)
		runner.afterFirst = func() {
			head, ok, err := s.GetRegistryHead(boundedTestContext(t), store.TransitionRegistryV1)
			if err != nil || !ok {
				t.Error("registry head missing")
				return
			}
			publishRevision(t, s.Store, transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1, Revision: 2, Parent: head, Entries: []transitionreg.Descriptor{}}, head)
		}
		w := postMCP(t, h, "Bearer "+tok, "["+mcpCallItem+","+mcpCallItem+"]")
		if !strings.Contains(w.Body.String(), `"code":-32603`) || strings.Contains(w.Header().Get("Content-Type"), "event-stream") || runner.runs.Load() != 1 || reader.calls != 3 || len(s.tasks) != 1 {
			t.Fatalf("stale batch invocation runs=%d reads=%d wire=%s", runner.runs.Load(), reader.calls, w.Body)
		}
		assertReceipts(t, s, 1)
	})
}
func TestMCPBatchConformance(t *testing.T) {
	t.Run("old_versions_and_item_errors", func(t *testing.T) {
		for _, version := range []string{"", "2025-03-26"} {
			h, tok, s, _, runner := batchFixture(t)
			r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(`[{"jsonrpc":"2.0","id":1,"method":"ping"},`+mcpCallItem+`,{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"unknown","arguments":{}}},{"jsonrpc":"2.0","id":4,"method":"unknown"},{"jsonrpc":"2.0","id":5,"method":"ping"}]`))
			r.Header.Set("Authorization", "Bearer "+tok)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept", "application/json, text/event-stream")
			r.Header.Set("MCP-Protocol-Version", version)
			w := httptest.NewRecorder()
			h.MCP(w, r)
			responses := batchResponses(t, w)
			if len(responses) != 5 || len(responses[2]["error"]) == 0 || len(responses[3]["error"]) == 0 || len(responses[4]["result"]) == 0 || runner.runs.Load() != 1 {
				t.Fatalf("item continuation responses=%v runs=%d", responses, runner.runs.Load())
			}
			assertReceipts(t, s, 1)
		}
	})
	t.Run("whole_host_failure", func(t *testing.T) {
		h, tok, s, _, runner := batchFixture(t)
		runner.failAt = 2
		w := postMCP(t, h, "Bearer "+tok, "["+mcpCallItem+","+mcpCallItem+","+mcpCallItem+"]")
		var wire struct {
			ID    any
			Error struct {
				Code    int
				Message string
			}
		}
		if err := json.Unmarshal(w.Body.Bytes(), &wire); err != nil || wire.ID != nil || wire.Error.Code != -32603 || wire.Error.Message != "host callback failed" || w.Code != 200 || strings.Contains(w.Header().Get("Content-Type"), "event-stream") || runner.runs.Load() != 2 || len(s.tasks) != 2 {
			t.Fatalf("whole failure runs=%d wire=%d %s", runner.runs.Load(), w.Code, w.Body)
		}
		assertReceipts(t, s, 1)
	})
	t.Run("notifications_and_bad_arrays", func(t *testing.T) {
		h, tok, _, _, runner := batchFixture(t)
		w := postMCP(t, h, "Bearer "+tok, `[{"jsonrpc":"2.0","method":"notifications/initialized"}]`)
		if w.Code != 202 || w.Body.Len() != 0 || runner.runs.Load() != 0 {
			t.Fatalf("notifications=%d %s", w.Code, w.Body)
		}
		for _, body := range []string{`[]`, `[bad]`, `[{"jsonrpc":"1.0","id":1,"method":"ping"}]`} {
			w := postMCP(t, h, "Bearer "+tok, body)
			if w.Code != 400 {
				t.Fatalf("bad array %s status=%d", body, w.Code)
			}
		}
	})
	t.Run("newer_versions_refuse", func(t *testing.T) {
		h, tok, _, _, runner := batchFixture(t)
		for _, version := range []string{"2025-06-18", "2025-11-25"} {
			r := httptest.NewRequest("POST", "/mcp/", strings.NewReader("["+mcpCallItem+"]"))
			r.Header.Set("Authorization", "Bearer "+tok)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept", "application/json, text/event-stream")
			r.Header.Set("MCP-Protocol-Version", version)
			w := httptest.NewRecorder()
			h.MCP(w, r)
			if w.Code != 400 || runner.runs.Load() != 0 {
				t.Fatalf("newer batch=%d runs=%d", w.Code, runner.runs.Load())
			}
		}
	})
}

type goldenAdapter struct{}

func (goldenAdapter) ResolveSession(context.Context, *http.Request) (protocol.Session, error) {
	return "fixture", nil
}
func (goldenAdapter) Tools(context.Context, protocol.Session) ([]protocol.ToolDescriptor, error) {
	return []protocol.ToolDescriptor{}, nil
}
func (goldenAdapter) Invoke(context.Context, protocol.Session, protocol.Invocation) (protocol.InvocationResult, error) {
	return protocol.InvocationResult{}, errors.New("unexpected invoke")
}
func TestMCPWireConformance(t *testing.T) {
	cfg := testConfig(openStore(t))
	tok := mintToken(t, cfg.Heads.(*store.Store), "ep-wire", []broker.Capability{liveGrant("alpha")})
	h := mustHandler(t, cfg)
	body := `{"jsonrpc":"2.0","id":7,"method":"ping"}`
	actual := postMCP(t, h, "Bearer "+tok, body)
	runner, err := hostcall.New(time.Second, 1)
	if err != nil {
		t.Fatal(err)
	}
	seam, err := mcphttp.NewHandler(mcphttp.Config{Agent: cfg.Agent, Resolver: goldenAdapter{}, Tools: goldenAdapter{}, Invoker: goldenAdapter{}, Runner: runner})
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/event-stream")
	reference := httptest.NewRecorder()
	seam.ServeHTTP(reference, r)
	if os.Getenv("WORLD_UPDATE_MCP_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("testdata/mcp_success.sse", reference.Body.Bytes(), 0644); err != nil {
			t.Fatal(err)
		}
	}
	golden, err := os.ReadFile("testdata/mcp_success.sse")
	if err != nil {
		t.Fatal(err)
	}
	if actual.Body.String() != string(golden) || actual.Body.String() != reference.Body.String() || actual.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("upstream SSE differs: %q / %q", actual.Body.String(), reference.Body.String())
	}
	initialized := mcpPayload(t, postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":2,"method":"initialize","params":{"protocolVersion":"2099-01-01"}}`))
	if !strings.Contains(string(initialized["result"]), `"protocolVersion":"2025-11-25"`) {
		t.Fatalf("version negotiation=%s", initialized["result"])
	}
	for _, headers := range []map[string]string{{"Content-Type": "text/plain"}, {"Accept": "application/json"}, {"MCP-Protocol-Version": "invalid"}} {
		r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+tok)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json, text/event-stream")
		for k, v := range headers {
			r.Header.Set(k, v)
		}
		w := httptest.NewRecorder()
		h.MCP(w, r)
		if w.Code != 400 {
			t.Fatalf("transport headers%v=%d", headers, w.Code)
		}
	}
	w := postMCP(t, h, "Bearer "+tok, strings.Repeat(" ", mcphttp.MaxRequestBodyBytes+1))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "invalid MCP request body") {
		t.Fatalf("body cap=%d %s", w.Code, w.Body)
	}
}
func TestMCPInvalidInvocationArguments(t *testing.T) {
	for _, args := range []string{`null`, `[]`, `1`, `"text"`} {
		h, tok, s, _, runner := batchFixture(t)
		body := strings.Replace(mcpCallItem, `{"x":1}`, args, 1)
		w := postMCP(t, h, "Bearer "+tok, body)
		if !strings.Contains(w.Body.String(), `"code":-32603`) || runner.runs.Load() != 0 || len(s.ids) != 0 {
			t.Fatalf("invalid input %s runs=%d intents=%d wire=%s", args, runner.runs.Load(), len(s.ids), w.Body)
		}
	}
}
func TestMCPTrueAbsenceAndPublication(t *testing.T) {
	st := openStore(t)
	cfg := testConfig(st)
	tok := mintToken(t, st, "ep-empty", []broker.Capability{liveGrant("alpha")})
	h := mustHandler(t, cfg)
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if names := mcpNames(t, w); len(names) != 0 || !strings.Contains(w.Body.String(), `"tools":[]`) {
		t.Fatalf("empty tools=%v wire=%s", names, w.Body)
	}
	cfg.Heads = fixedHeads{ok: false}
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h = mustHandler(t, cfg)
	if names := mcpNames(t, postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)); !reflect.DeepEqual(names, []string{"tools_decho"}) {
		t.Fatalf("raced publication=%v", names)
	}
}
func TestMCPAbsentPrecheckSuccessfulInvocation(t *testing.T) {
	cfg, tok, st, runner := mcpFixture(t)
	cfg.Heads = fixedHeads{ok: false}
	h := mustHandler(t, cfg)
	task := ""
	h.mintTask = func() (string, error) { var err error; task, err = mintMCPTask(); return task, err }
	w := postMCP(t, h, "Bearer "+tok, mcpCallItem)
	if len(mcpPayload(t, w)["error"]) != 0 || runner.runs != 1 {
		t.Fatalf("raced publication invoke runs=%d wire=%s", runner.runs, w.Body)
	}
	receipt, ok, err := st.GetReceipt(boundedTestContext(t), coordinator.InvocationID(coordinator.SurfaceMCP, "ep-a", task))
	if err != nil || !ok || receipt.State != store.ReceiptResolved {
		t.Fatalf("published invocation receipt=%+v/%t/%v", receipt, ok, err)
	}
}

func TestMCPTaskEntropyFailure(t *testing.T) {
	h, tok, s, _, runner := batchFixture(t)
	h.mintTask = func() (string, error) { return "", fmt.Errorf("entropy unavailable") }
	w := postMCP(t, h, "Bearer "+tok, mcpCallItem)
	if runner.runs.Load() != 0 || len(s.ids) != 0 || !strings.Contains(w.Body.String(), `"code":-32603`) {
		t.Fatalf("entropy failure executed runs=%d intents=%d wire=%s", runner.runs.Load(), len(s.ids), w.Body)
	}
}

// lockedSink is an ErrorLog written from the hostcall goroutine.
type lockedSink struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *lockedSink) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *lockedSink) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

// TestMCPRefusalLineNamesSurfaceAndCause is row 136 AC1.1a — the
// binary-free companion of AC1.1 that CI runs (kills MUT-MCP-LOG-LABEL and
// MUT-MCP-LOG-NOCAUSE without WORLD_TOOL_AILANG_BIN): a failed MCP
// tools/call writes exactly one operator line, labelled mcp, naming the
// minted invocation id and carrying the cause.
func TestMCPRefusalLineNamesSurfaceAndCause(t *testing.T) {
	h, tok, s, _, runner := batchFixture(t)
	var sink lockedSink
	h.errorLog = &sink
	runner.failAt = 1
	w := postMCP(t, h, "Bearer "+tok, mcpCallItem)
	if !strings.Contains(w.Body.String(), `"host callback failed"`) || len(s.tasks) != 1 {
		t.Fatalf("wire=%s tasks=%d", w.Body, len(s.tasks))
	}
	want := "ailang-worldd: mcp refusal: tools/call " + coordinator.InvocationID(coordinator.SurfaceMCP, "ep-a", s.tasks[0]) + `: "`
	if got := sink.String(); strings.Count(got, "\n") != 1 || !strings.HasPrefix(got, want) || !strings.Contains(got, "batch injected capsule failure") {
		t.Fatalf("operator log = %q, want one line starting %q carrying the cause", got, want)
	}
}

// failAfterHeads delegates the first okCalls head reads, then fails.
type failAfterHeads struct {
	inner   HeadReader
	okCalls int
	err     error
	mu      sync.Mutex
	n       int
}

func (f *failAfterHeads) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	f.mu.Lock()
	f.n++
	fail := f.n > f.okCalls
	f.mu.Unlock()
	if fail {
		return hashref.HashRef{}, false, f.err
	}
	return f.inner.GetRegistryHead(ctx, name)
}

// TestMCPPreDispatchRefusalLinesAreLabelled is the sibling of
// TestMCPRefusalLineNamesSurfaceAndCause for the three MCP refusal sites that
// have no invocation id yet (row 136 round 2): the registry read behind
// tools/list, a descriptor set that cannot be projected to MCP, and the
// admission read behind tools/call. Each writes exactly one operator line
// labelled mcp with the method and the "-" id.
func TestMCPPreDispatchRefusalLinesAreLabelled(t *testing.T) {
	const list = `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	cases := []struct {
		name, body, prefix string
		setup              func(t *testing.T, h *Handler, st *store.Store)
		cause              string
	}{
		{"tools_list_registry_read", list, "ailang-worldd: mcp refusal: tools/list -: \"",
			func(t *testing.T, h *Handler, _ *store.Store) {
				h.heads = errHeads{err: errors.New("injected head failure")}
			},
			"injected head failure"},
		{"tools_list_unprojectable", list, "ailang-worldd: mcp refusal: tools/list -: \"",
			func(t *testing.T, h *Handler, st *store.Store) {
				head, ok, err := st.GetRegistryHead(boundedTestContext(t), store.TransitionRegistryV1)
				if err != nil || !ok {
					t.Fatalf("registry head ok=%v err=%v", ok, err)
				}
				d := descriptor("tools.echo", "alpha")
				d.InputSchema = []byte(`{"type":"array"}`)
				publishRevision(t, st, transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1,
					InterfaceHash: transitionreg.InterfaceHashV1, Revision: 2, Parent: head, Entries: []transitionreg.Descriptor{d}}, head)
			},
			"input schema type must be object"},
		// tools/call lists first (name check), then admits inside Invoke: the
		// head read fails only from its second call on.
		{"tools_call_admission", mcpCallItem, "ailang-worldd: mcp refusal: tools/call -: \"",
			func(t *testing.T, h *Handler, _ *store.Store) {
				h.heads = &failAfterHeads{inner: h.heads, okCalls: 1, err: errors.New("injected head failure")}
			},
			"injected head failure"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h, tok, s, _, _ := batchFixture(t)
			var sink lockedSink
			h.errorLog = &sink
			c.setup(t, h, s.Store)
			w := postMCP(t, h, "Bearer "+tok, c.body)
			got := sink.String()
			if strings.Count(got, "\n") != 1 || !strings.HasPrefix(got, c.prefix) || !strings.Contains(got, c.cause) {
				t.Fatalf("operator log = %q, want one line starting %q carrying %q (wire=%s)", got, c.prefix, c.cause, w.Body)
			}
		})
	}
}
