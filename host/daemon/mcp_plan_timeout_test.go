package daemon

// Row 153 AC1.5: the /mcp/ wire for a plan-phase timeout is the released
// handler's frozen envelope. This pins it so the day ailang#1602 types the
// callback error (row 159), this test flips deliberately instead of silently.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// planBlocksRunner blocks every plan on the phase context. A mutant that
// drops the phase budget fails after the bound, never hangs the suite.
type planBlocksRunner struct{ bound time.Duration }

func (r planBlocksRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	select {
	case <-ctx.Done():
		return capsule.Result{}, ctx.Err()
	case <-time.After(r.bound):
		return capsule.Result{}, errors.New("the plan budget never cancelled the capsule")
	}
}

func TestMCPPlanTimeoutFrozenWire(t *testing.T) {
	f := newWSFixture(t)
	var log syncLog
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root,
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease), ErrorLog: &log})
	publishReadTool(t, d, hashref.SumSHA256([]byte("fake-interpreter")),
		"module host/capsule/main\n\nexport func main(input: string) -> string {\n  input\n}\n")
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(srv.Close)
	r := &seRig{t: t, f: f, d: d, srv: srv, client: &http.Client{Timeout: DefaultClientTimeout}}
	r.serveWith(d.store, planBlocksRunner{bound: coordinator.PlanPhaseBudget + 5*time.Second}, &log)

	token := r.mint("ep1", "Workspace.Read")
	before := r.entryCount()
	wire, took := r.call(token, "ws_dread", map[string]any{"path": "data.txt"})
	if wire.Error == nil || wire.Error.Code != -32603 || wire.Error.Message != "host callback timed out" {
		t.Fatalf("tools/call under a plan timeout error = %s, want the frozen -32603 \"host callback timed out\" (row 159 / ailang#1602 flips this)", wire.errText())
	}
	if took > coordinator.PlanPhaseBudget+4*time.Second {
		t.Fatalf("call took %s: the plan budget (%s) did not end it", took, coordinator.PlanPhaseBudget)
	}
	text := log.String()
	if strings.Count(text, "\n") != 1 || !strings.Contains(text, "plan phase exceeded its") {
		t.Fatalf("operator log = %q, want exactly one line naming the plan phase timeout", text)
	}
	if after := r.entryCount(); after != before {
		t.Fatalf("the timed-out call committed (%d -> %d)", before, after)
	}
}
