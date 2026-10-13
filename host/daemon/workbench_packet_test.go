package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

const packetSectionStart = `<section aria-label="decision packet">`

// packetSection returns the decision-packet section of body; t.Fatal if absent.
func packetSection(t *testing.T, body string) string {
	t.Helper()
	sec, ok := workbenchRegion(body, packetSectionStart, "</section>")
	if !ok {
		t.Fatalf("decision packet section absent; body=%s", body)
	}
	return sec
}

// wantPacketGradeLine renders the expected evidence line for one grade.
func wantPacketGradeLine(kind, ref string, g workbench.GradeView) string {
	var grade string
	if g.Available {
		grade = "<span>" + string(g.Label) + "</span>"
		if g.HasVerdict {
			grade += " verdict: " + string(g.Verdict)
		}
	} else {
		reason := g.Unavailable
		if reason == "" {
			reason = "no grade reason was supplied"
		}
		grade = "GRADE UNAVAILABLE — " + html.EscapeString(reason)
	}
	return "<p>evidence: " + kind + ` <span class="hash">` + ref + "</span> — " + grade + "</p>"
}

// producePacket dispatches one Human.Approve invocation with the explicit test binder
// and returns its entry and the packet first carried by that entry.
func producePacket(t *testing.T, d *Daemon, plan, task string) (coordinator.Result, store.Object) {
	t.Helper()
	c := humanCoordinator(t, d, testHumanBinder(d), plan)
	res, err := c.Dispatch(boundedTestContext(t), humanCall(t, d, "ep1", task))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	return res, packetForEntry(t, d, res.EntryIndex)
}

// recordTransitionFn reads the transitionFn of an invocation record.
func recordTransitionFn(t *testing.T, d *Daemon, recordRef hashref.HashRef) hashref.HashRef {
	t.Helper()
	obj, ok, err := d.store.GetObject(boundedTestContext(t), recordRef)
	if err != nil || !ok {
		t.Fatalf("GetObject record %s: ok=%v err=%v", recordRef, ok, err)
	}
	var m map[string]any
	if err := json.Unmarshal(obj.Payload, &m); err != nil {
		t.Fatal(err)
	}
	return hashref.MustParse(m["transitionFn"].(string))
}

// countingPacketScanStore counts packet scans through the page's read surface.
type countingPacketScanStore struct {
	readStore
	calls *atomic.Int32
}

func (s countingPacketScanStore) ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error) {
	if id == coordinator.DecisionPacketV1 {
		s.calls.Add(1)
	}
	return s.readStore.ObjectsBySemanticID(ctx, id, after, limit)
}

// TestWorkbenchPacketValidated (AC4).
func TestWorkbenchPacketValidated(t *testing.T) {
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	res, k := producePacket(t, d, humanPlan("e1"), "t1")
	p, err := coordinator.DecodeDecisionPacket(k.Payload)
	if err != nil {
		t.Fatal(err)
	}
	code, body := gradePage(t, d, k.Hash)
	if code != http.StatusOK {
		t.Fatalf("status = %d", code)
	}
	sec := packetSection(t, body)
	ctx := boundedTestContext(t)
	want := []string{
		"<p>" + html.EscapeString(fmt.Sprintf(packetValidatedFormat, res.EntryIndex)) + "</p>",
		`href="/workbench?object=` + p.ProposalHash.String() + `"`,
		`href="/workbench?object=` + p.RequestRef.String() + `"`,
		fmt.Sprintf(`href="/workbench?from=%d&amp;entry=%d"`, res.EntryIndex, res.EntryIndex),
		"approvals chain: pending",
		"policy Cancel · escalationsRemaining 0",
		"createdAt 1700 · deadlineAt 1760 — logical times; this host emits no Timeout transition (item 7 owns the emitter), so this page does not state whether the deadline has passed",
	}
	for _, w := range want {
		if strings.Count(sec, w) != 1 {
			t.Fatalf("section missing %q exactly once; section=%s", w, sec)
		}
	}
	pg, err := d.objectGradeView(ctx, p.ProposalHash)
	if err != nil {
		t.Fatal(err)
	}
	tg, err := d.objectGradeView(ctx, recordTransitionFn(t, d, res.RecordRef))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sec, wantPacketGradeLine("proposal", p.ProposalHash.String(), pg)) {
		t.Fatalf("proposal grade line absent; section=%s", sec)
	}
	if !strings.Contains(sec, wantPacketGradeLine("transition function", recordTransitionFn(t, d, res.RecordRef).String(), tg)) {
		t.Fatalf("transition grade line absent; section=%s", sec)
	}

	// Non-packet arm (PD6): the plan page is not scanned for packets.
	var n atomic.Int32
	d.reads = countingPacketScanStore{readStore: d.reads, calls: &n}
	code, body = gradePage(t, d, p.ProposalHash)
	if code != http.StatusOK {
		t.Fatalf("plan page status = %d", code)
	}
	for _, s := range []string{packetSectionStart, "decision packet validated", "not a validated decision packet"} {
		if strings.Contains(body, s) {
			t.Fatalf("plan page contains %q", s)
		}
	}
	if n.Load() != 0 {
		t.Fatalf("plan page scanned packets %d times, want 0", n.Load())
	}
	// Control: the packet page scans once.
	if code, _ := gradePage(t, d, k.Hash); code != http.StatusOK || n.Load() != 1 {
		t.Fatalf("packet page: status=%d scans=%d, want 200/1", code, n.Load())
	}
}
