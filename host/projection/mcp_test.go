package projection

import (
	"context"
	"encoding/json"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func mcpFixture(t *testing.T) (Config, string, *store.Store, *countingRunner) {
	st := openStore(t)
	genesis := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "projection-test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	world := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: genesis.Hash, LogHead: entryHash}
	interp := hashref.SumSHA256([]byte("interpreter"))
	err := st.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{genesis}, NextWorld: world,
		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")),
			Interpreter: interp, PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "projection-test"},
			EntryHash: entryHash, TransitionRef: genesis.Hash}})
	if err != nil {
		t.Fatal(err)
	}
	source, err := canon.Source([]byte("module transitions/echo\nexport func main(input: string) -> string { input }\n"))
	if err != nil {
		t.Fatal(err)
	}
	src := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
		SemanticID: "world/transition-source/v1", Provenance: "projection-test", Payload: source}
	if err := st.PutObject(boundedTestContext(t), src); err != nil {
		t.Fatal(err)
	}
	d := descriptor("tools.echo", "alpha")
	d.TransitionFn = src.Hash
	d.Interpreter = interp
	d.DeclaredEffects = nil
	seedRegistry(t, st, d)
	runner := &countingRunner{}
	coord, err := coordinator.New(coordinator.Config{Store: st, Runner: runner,
		Binder: func(ep string, grants []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(st, ep, grants, nil)
		},
		Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(st)
	cfg.Coordinator = coord
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
	return cfg, tok, st, runner
}
func postMCP(t *testing.T, h *Handler, auth, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(body))
	r.Header.Set("Authorization", auth)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/event-stream")
	w := httptest.NewRecorder()
	h.MCP(w, r)
	return w
}
func mcpPayload(t *testing.T, w *httptest.ResponseRecorder) map[string]json.RawMessage {
	t.Helper()
	body := w.Body.String()
	if !strings.HasPrefix(body, "event: message\ndata: ") {
		t.Fatalf("not upstream SSE: %q", body)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(body, "event: message\ndata: "))), &out); err != nil {
		t.Fatal(err)
	}
	return out
}
func mcpNames(t *testing.T, w *httptest.ResponseRecorder) []string {
	t.Helper()
	payload := mcpPayload(t, w)
	var result struct{ Tools []struct{ Name string } }
	if err := json.Unmarshal(payload["result"], &result); err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(result.Tools))
	for _, d := range result.Tools {
		names = append(names, d.Name)
	}
	return names
}
func TestMCPListExactSetPerSession(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.alpha", "alpha"), descriptor("tools.beta", "beta"))
	cfg := testConfig(st)
	h := mustHandler(t, cfg)
	for _, effect := range []string{"alpha", "beta"} {
		tok := mintToken(t, st, "ep-"+effect, []broker.Capability{liveGrant(effect)})
		w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
		if names := mcpNames(t, w); !reflect.DeepEqual(names, []string{"tools_d" + effect}) {
			t.Fatalf("session set=%v", names)
		}
	}
}
func TestMCPDenialMatrix(t *testing.T) {
	st := openStore(t)
	cfg := testConfig(st)
	heads := &countingHeads{inner: st}
	reader := &countingReader{inner: cfg.Reader}
	cfg.Heads = heads
	cfg.Reader = reader
	h := mustHandler(t, cfg)
	expired := mintExpiredToken(t, st)
	for _, auth := range []string{"", "garbage", "Bearer " + strings.Repeat("f", 64), "Bearer " + expired} {
		w := postMCP(t, h, auth, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
		if heads.calls != 0 || reader.calls != 0 {
			t.Fatalf("denied registry reads %d/%d", heads.calls, reader.calls)
		}
		if w.Code != 401 {
			t.Fatalf("denied status=%d body=%s", w.Code, w.Body)
		}
	}
	r := httptest.NewRequest("POST", "/mcp/", strings.NewReader(`{}`))
	tok := mintToken(t, st, "ep-live", []broker.Capability{liveGrant("alpha")})
	r.Header.Set("X-API-Key", "Bearer "+tok)
	w := httptest.NewRecorder()
	h.MCP(w, r)
	if w.Code != 401 {
		t.Fatalf("API key admitted: %d", w.Code)
	}
	if heads.calls != 0 || reader.calls != 0 {
		t.Fatalf("denied registry reads %d/%d", heads.calls, reader.calls)
	}
	postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if heads.calls == 0 || reader.calls == 0 {
		t.Fatal("authorized read control absent")
	}
}
func TestMCPCallDispatch(t *testing.T) {
	cfg, tok, st, runner := mcpFixture(t)
	h := mustHandler(t, cfg)
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{"x":1}}}`)
	payload := mcpPayload(t, w)
	if len(payload["error"]) != 0 || runner.runs != 1 {
		t.Fatalf("dispatch result=%s runs=%d", w.Body, runner.runs)
	}
	world, ok, err := st.SelectedHead(boundedTestContext(t))
	if err != nil || !ok || world.IsZero() {
		t.Fatalf("committed head=%v/%t/%v", world, ok, err)
	}
}
func TestMCPInvokeListedThenRevoked(t *testing.T) {
	cfg, tok, st, runner := mcpFixture(t)
	real := cfg.Reader
	reader := &revokeReader{inner: real, st: st, t: t}
	cfg.Reader = reader
	h := mustHandler(t, cfg)
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if !strings.Contains(w.Body.String(), `"code":-32603`) || runner.runs != 0 || reader.calls != 2 {
		t.Fatalf("revoked dispatch=%d reads=%d wire=%s", runner.runs, reader.calls, w.Body)
	}
}

type revokeReader struct {
	inner transitionreg.Reader
	st    *store.Store
	t     *testing.T
	calls int
}

func (r *revokeReader) ReadSnapshot(ctx context.Context) (transitionreg.Snapshot, error) {
	r.calls++
	if r.calls == 2 {
		head, ok, err := r.st.GetRegistryHead(ctx, store.TransitionRegistryV1)
		if err != nil || !ok {
			r.t.Fatalf("head control %v/%t", err, ok)
		}
		publishRevision(r.t, r.st, transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1, Revision: 2, Parent: head, Entries: []transitionreg.Descriptor{}}, head)
	}
	return r.inner.ReadSnapshot(ctx)
}
func TestMCPInvokeCanceledBeforeMint(t *testing.T) {
	cfg, _, _, runner := mcpFixture(t)
	h := mustHandler(t, cfg)
	mints := 0
	h.mintTask = func() (string, error) { mints++; return strings.Repeat("a", 64), nil }
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	cancel()
	_, err := (mcpAdapter{h}).Invoke(ctx, &authority.SessionBinding{EpisodeID: "ep-a", Caps: []broker.Capability{liveGrant("alpha")}}, protocol.Invocation{Name: "tools_decho", Arguments: json.RawMessage(`{}`)})
	if err == nil || mints != 0 || runner.runs != 0 {
		t.Fatalf("expired invoke err=%v mint=%d runs=%d", err, mints, runner.runs)
	}
}
