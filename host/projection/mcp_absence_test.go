package projection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"strings"
	"testing"
	"time"
)

type stageReader struct {
	inner, failure transitionreg.Reader
	calls          int
	failAt         int
}

func (r *stageReader) ReadSnapshot(ctx context.Context) (transitionreg.Snapshot, error) {
	r.calls++
	if r.calls == r.failAt {
		return r.failure.ReadSnapshot(ctx)
	}
	return r.inner.ReadSnapshot(ctx)
}
func TestMCPAbsentPrecheckSnapshotFailure(t *testing.T) {
	for _, leg := range []string{"list", "invoke"} {
		for _, kind := range []string{"direct", "real_reader", "same_text", "deadline"} {
			t.Run(leg+"/"+kind, func(t *testing.T) {
				cfg, tok, st, runner := mcpFixture(t)
				cfg.Heads = fixedHeads{ok: false}
				cause := error(errors.New("private MCP reached snapshot fault"))
				if kind == "same_text" {
					cause = errors.New("read transition registry: head is absent")
				}
				if kind == "deadline" {
					cause = context.DeadlineExceeded
				}
				direct := &absenceReader{err: cause}
				real := &racedObjectStore{ObjectStore: st, err: cause}
				failure := transitionreg.Reader(direct)
				if kind == "real_reader" {
					failure = transitionreg.NewReader(real)
				}
				reader := &stageReader{inner: cfg.Reader, failure: failure, failAt: 1}
				if leg == "invoke" {
					reader.failAt = 2
				}
				cfg.Reader = reader
				var sink bytes.Buffer
				cfg.ErrorLog = &sink
				h := mustHandler(t, cfg)
				mints := 0
				h.mintTask = func() (string, error) { mints++; return strings.Repeat("a", 64), nil }
				body := `{"jsonrpc":"2.0","id":7,"method":"tools/list"}`
				if leg == "invoke" {
					body = `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`
				}
				w := postMCP(t, h, "Bearer "+tok, body)
				msg := "host callback failed"
				if kind == "deadline" {
					msg = "host callback timed out"
				}
				want := `{"jsonrpc":"2.0","id":7,"error":{"code":-32603,"message":"` + msg + `"}}` + "\n"
				var wire struct {
					ID    int
					Error struct {
						Code    int
						Message string
					}
				}
				if err := json.Unmarshal(w.Body.Bytes(), &wire); err != nil || wire.ID != 7 || wire.Error.Code != -32603 || wire.Error.Message != msg || w.Code != 200 || w.Body.String() != want || strings.Contains(w.Header().Get("Content-Type"), "event-stream") {
					t.Fatalf("whole host envelope=%d %s want %s", w.Code, w.Body, want)
				}
				if runner.runs != 0 || mints != 0 || reader.calls != reader.failAt || strings.Count(sink.String(), cause.Error()) != 1 {
					t.Fatalf("failure attribution runs=%d mint=%d reads=%d sink=%s", runner.runs, mints, reader.calls, &sink)
				}
				if kind == "real_reader" {
					if real.heads != 1 || real.objects != 1 {
						t.Fatalf("real reader fault not reached %d/%d", real.heads, real.objects)
					}
				} else if direct.calls != 1 {
					t.Fatalf("direct fault reached=%d", direct.calls)
				}
			})
		}
	}
}
func TestMCPSurfaceRefusalRecovers(t *testing.T) {
	cfg, tok, st, runner := mcpFixture(t)
	head, _, err := st.GetRegistryHead(boundedTestContext(t), store.TransitionRegistryV1)
	if err != nil {
		t.Fatal(err)
	}
	// Retain the admitted source/interpreter while changing only its stable ID.
	snap, err := cfg.Reader.ReadSnapshot(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	ds := snap.List()
	ds[0].ID = strings.Repeat("a", 32) + "/" + strings.Repeat("b", 32)
	next := publishRevision(t, st, transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1, Revision: 2, Parent: head, Entries: ds}, head)
	h := mustHandler(t, cfg)
	binding, err := cfg.Resolver.ResolveContext(boundedTestContext(t), "Bearer "+tok, time.Now().Unix())
	if err != nil || binding.Success == nil {
		t.Fatalf("binding=%v/%v", binding, err)
	}
	if tools, err := (mcpAdapter{h}).Tools(boundedTestContext(t), binding.Success); err == nil || len(tools) != 0 {
		t.Fatalf("adapter must refuse long surface before upstream validation: tools=%v err=%v", tools, err)
	}
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/list"}`)
	if strings.Contains(w.Header().Get("Content-Type"), "event-stream") || !strings.Contains(w.Body.String(), `"code":-32603`) || runner.runs != 0 {
		t.Fatalf("long surface not refused %s", w.Body)
	}
	ds[0].ID = "tools.echo"
	publishRevision(t, st, transitionreg.Revision{SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1, Revision: 3, Parent: next, Entries: ds}, next)
	if names := mcpNames(t, postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/list"}`)); len(names) != 1 || names[0] != "tools_decho" {
		t.Fatalf("repaired list=%v", names)
	}
	w = postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if len(mcpPayload(t, w)["error"]) != 0 || runner.runs != 1 {
		t.Fatalf("repaired call=%s runs=%d", w.Body, runner.runs)
	}
}

type snapshotFunc func(context.Context) (transitionreg.Snapshot, error)

func (f snapshotFunc) ReadSnapshot(ctx context.Context) (transitionreg.Snapshot, error) {
	return f(ctx)
}
func TestMCPInvokeCancelAfterAdmission(t *testing.T) {
	cfg, _, _, runner := mcpFixture(t)
	snap, err := cfg.Reader.ReadSnapshot(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	defer cancel()
	cfg.Reader = snapshotFunc(func(context.Context) (transitionreg.Snapshot, error) { cancel(); return snap, nil })
	h := mustHandler(t, cfg)
	mints := 0
	h.mintTask = func() (string, error) { mints++; return strings.Repeat("a", 64), nil }
	_, err = (mcpAdapter{h}).Invoke(ctx, &authority.SessionBinding{EpisodeID: "ep-a", Caps: []broker.Capability{liveGrant("alpha")}}, protocol.Invocation{Name: "tools_decho", Arguments: json.RawMessage(`{}`)})
	if !errors.Is(err, context.Canceled) || mints != 0 || runner.runs != 0 {
		t.Fatalf("after-admission cancellation err=%v mint=%d runs=%d", err, mints, runner.runs)
	}
}
func TestMCPBoundsValidation(t *testing.T) {
	base := testConfig(openStore(t))
	for _, change := range []func(*Config){func(c *Config) { c.CredentialWait = 0 }, func(c *Config) { c.CallbackTimeout = 0 }, func(c *Config) { c.MaxCallbacks = 0 }, func(c *Config) { c.WriteWait = 0 }, func(c *Config) { c.InvokeWait = c.WriteWait }} {
		c := base
		change(&c)
		if _, err := New(c); err == nil {
			t.Fatal("invalid scalar accepted")
		}
	}
	if _, err := New(base); err != nil {
		t.Fatal(err)
	}
}

type deadlineResolver struct{ observed time.Duration }

func (r *deadlineResolver) ResolveContext(ctx context.Context, _ string, _ int64) (authority.ResolveOutcome, error) {
	d, ok := ctx.Deadline()
	if !ok {
		return authority.ResolveOutcome{}, errors.New("deadline missing")
	}
	r.observed = time.Until(d)
	return authority.ResolveOutcome{Success: &authority.SessionBinding{EpisodeID: "ep-a", Caps: []broker.Capability{liveGrant("alpha")}}}, nil
}
func TestMCPResolverInnerBudget(t *testing.T) {
	cfg := testConfig(openStore(t))
	cfg.CredentialWait = 30 * time.Millisecond
	r := &deadlineResolver{}
	cfg.Resolver = r
	h := mustHandler(t, cfg)
	postMCP(t, h, "Bearer anything", `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if r.observed <= 0 || r.observed > cfg.CredentialWait {
		t.Fatalf("credential deadline=%s want <=%s", r.observed, cfg.CredentialWait)
	}
}

func TestMCPListedThenTrulyAbsent(t *testing.T) {
	cfg, tok, _, runner := mcpFixture(t)
	cfg.Heads = fixedHeads{ok: false}
	reader := &stageReader{inner: cfg.Reader, failure: &absenceReader{err: transitionreg.RegistryHeadAbsentError{}}, failAt: 2}
	cfg.Reader = reader
	h := mustHandler(t, cfg)
	mints := 0
	h.mintTask = func() (string, error) { mints++; return strings.Repeat("a", 64), nil }
	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_decho","arguments":{}}}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"code":-32603`) || strings.Contains(w.Header().Get("Content-Type"), "event-stream") || reader.calls != 2 || mints != 0 || runner.runs != 0 {
		t.Fatalf("listed then absent reads=%d mint=%d runs=%d wire=%s", reader.calls, mints, runner.runs, w.Body)
	}
}

type mcpFailResolver struct{ cause error }

func (r mcpFailResolver) ResolveContext(context.Context, string, int64) (authority.ResolveOutcome, error) {
	return authority.ResolveOutcome{}, r.cause
}
func TestMCPResolverHostFailure(t *testing.T) {
	for _, cause := range []error{errors.New("private resolver fault"), context.DeadlineExceeded} {
		cfg, tok, _, runner := mcpFixture(t)
		cfg.Resolver = mcpFailResolver{cause}
		reader := &countingReader{inner: cfg.Reader}
		cfg.Reader = reader
		h := mustHandler(t, cfg)
		w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/list"}`)
		message := "host callback failed"
		if errors.Is(cause, context.DeadlineExceeded) {
			message = "host callback timed out"
		}
		want := `{"jsonrpc":"2.0","id":7,"error":{"code":-32603,"message":"` + message + `"}}` + "\n"
		if w.Code != 200 || w.Body.String() != want || reader.calls != 0 || runner.runs != 0 {
			t.Fatalf("resolver carrier reads=%d runs=%d wire=%d %s", reader.calls, runner.runs, w.Code, w.Body)
		}
	}
}
