package projection

// Row 255 (design_docs/planned/w-exec-failure-carries-a-reason.md §5, TP): an
// MCP tools/call whose effect failed with a failure line returns that line in
// the tool result's world block and logs it exactly once on the operator's
// error log. On A2A, the same task sent twice runs once and logs once: the
// resend is read back from the journal (Reconciled) and logs nothing.

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const failureLine = "exec-stack-drift: /opt/node: digest changed"

// planRunner answers the plan phase with one alpha@world effect.
type planRunner struct{}

func (planRunner) RunContext(context.Context, capsule.Entry) (capsule.Result, error) {
	return capsule.Result{Stdout: []byte(`{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"alpha","scope":"world","cost":1,"payload":{}}],"finish":false,"result":null}` + "\n")}, nil
}

type classifiedFailure struct{}

func (classifiedFailure) Execute(context.Context, broker.EffectRequest, []byte) ([]byte, error) {
	return nil, &broker.HandlerFailure{Class: broker.FailExecStackDrift, Detail: "/opt/node: digest changed",
		Err: errors.New("broker: exec confinement stack drifted: /Users/op/secret/node: sha256 …")}
}

func TestMCPEffectFailureLineReturnedAndLoggedOnce(t *testing.T) {
	st := openStore(t)
	genesis := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "projection-test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	world := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: genesis.Hash, LogHead: entryHash}
	interp := hashref.SumSHA256([]byte("interpreter"))
	if err := st.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{genesis}, NextWorld: world,
		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")),
			Interpreter: interp, PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "projection-test"},
			EntryHash: entryHash, TransitionRef: genesis.Hash}}); err != nil {
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
	d := descriptor("tools.fx", "alpha") // declares alpha@world, cost 1
	d.TransitionFn, d.Interpreter = src.Hash, interp
	seedRegistry(t, st, d)
	coord, err := coordinator.New(coordinator.Config{Store: st, Runner: planRunner{},
		Binder: func(ep string, grants []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(st, ep, grants, broker.Registry{"alpha": classifiedFailure{}})
		},
		Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	var log bytes.Buffer
	cfg := testConfig(st)
	cfg.Coordinator, cfg.ErrorLog = coord, &log
	// The handler's budget is what the call has left after HandlerHeadroom (6 s).
	cfg.InvokeWait, cfg.CallbackTimeout = 20*time.Second, 20*time.Second
	h := mustHandler(t, cfg)
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	w := postMCP(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"tools_dfx","arguments":{"x":1}}}`)
	payload := mcpPayload(t, w)
	result := string(payload["result"])
	if len(payload["error"]) != 0 || !strings.Contains(result, `\"reason\":\"`+failureLine+`\"`) || strings.Contains(result, "/Users/op") {
		t.Fatalf("tools/call result %s, want the failure line in the world block and no cause text", w.Body)
	}
	want := regexp.MustCompile(`^ailang-worldd: mcp effect failed: mcp:ep-a:[0-9a-f]{64} e1 sha256:[0-9a-f]{64}: ` + regexp.QuoteMeta(failureLine) + "\n$")
	if !want.MatchString(log.String()) {
		t.Fatalf("error log = %q, want exactly one line naming the call, e1, its record and the failure line", log.String())
	}

	log.Reset()
	body := `{"jsonrpc":"2.0","id":8,"method":"tasks/send","params":{"id":"t1","metadata":{"skill_id":"tools.fx"},"message":{"parts":[{"type":"data","data":{"x":1}}]}}}`
	for i := 0; i < 2; i++ {
		if rec := postA2A(t, h, "Bearer "+tok, body); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"reason":"`+failureLine+`"`) {
			t.Fatalf("a2a send %d: %d %s", i, rec.Code, rec.Body)
		}
	}
	once := regexp.MustCompile(`^ailang-worldd: a2a effect failed: a2a:ep-a:t1 e1 sha256:[0-9a-f]{64}: ` + regexp.QuoteMeta(failureLine) + "\n$")
	if !once.MatchString(log.String()) {
		t.Fatalf("error log after the same a2a task twice = %q, want exactly one line", log.String())
	}
}
