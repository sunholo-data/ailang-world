package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
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

// forgeObject is a store object whose hashes are honest for its payload; the forgery is in its semantics.
func forgeObject(semanticID string, payload []byte) store.Object {
	return store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID: semanticID, Provenance: "packet-forge", Payload: payload}
}

// forgedRecordV2 mirrors coordinator's recordV2 wire (plan.go) field for field.
type forgedRecordV2 struct {
	InvocationID   string   `json:"invocationId"`
	EpisodeID      string   `json:"episodeId"`
	SkillID        string   `json:"skillId"`
	TransitionFn   string   `json:"transitionFn"`
	Interpreter    string   `json:"interpreter"`
	SemanticsEpoch int64    `json:"semanticsEpoch"`
	Input          string   `json:"input"`
	Output         string   `json:"output"`
	Plan           string   `json:"plan"`
	Effects        []string `json:"effects"`
}

// approvalRequestInChain files one pending Human.Approve request in the approvals chain and returns its ref.
func approvalRequestInChain(t *testing.T, d *Daemon, scope, requester string, now int64) hashref.HashRef {
	t.Helper()
	out, err := broker.NewHumanHandler(d.store).Execute(boundedTestContext(t),
		broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: scope, Cost: 1, Now: now}, []byte(`{"requester":"`+requester+`"}`))
	if err != nil {
		t.Fatal(err)
	}
	ref, err := broker.ParsePendingApproval(out)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

// forgedChain describes one forged invocation commit (record, effect record, result, packet).
type forgedChain struct {
	label       string
	request     hashref.HashRef            // the request the effect result's pending wire names
	scope       string                     // effect-record scope (must equal the request's scope to pass code 4)
	tf          hashref.HashRef            // record transitionFn
	planRef     hashref.HashRef            // record plan
	includePlan bool                       // whether the plan object is committed alongside
	packet      coordinator.DecisionPacket // the packet this chain carries
	rawPacket   []byte                     // when non-nil, replaces EncodeDecisionPacket(packet)
	extra       []store.Object             // more objects in the same commit
	failed      bool                       // the effect record's Failed bit
}

func forgedPlanPayload(label string) []byte { return []byte(`{"forged":"` + label + `"}`) }

func newForgedChain(label string, request hashref.HashRef, scope string) *forgedChain {
	plan := forgeObject(coordinator.EffectPlanV1, forgedPlanPayload(label))
	return &forgedChain{label: label, request: request, scope: scope,
		tf: hashref.SumSHA256([]byte("forge-tf-" + label)), planRef: plan.Hash, includePlan: true,
		packet: coordinator.DecisionPacket{ProposalHash: plan.Hash, RequestRef: request, CreatedAt: packetNow,
			DeadlineAt: packetNow + 60, Policy: "Cancel"}}
}

func (f *forgedChain) planObject() store.Object {
	return forgeObject(coordinator.EffectPlanV1, forgedPlanPayload(f.label))
}

// commitObjects appends one commit chained from the selected head and returns its entry index.
func commitObjects(t *testing.T, d *Daemon, label string, objects []store.Object, transitionRef hashref.HashRef) int64 {
	t.Helper()
	ctx := boundedTestContext(t)
	head, ok, err := d.store.SelectedHead(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var w store.World
	if !ok {
		w = seedGenesisEmbedded(t, d, label+"-genesis")
	} else {
		var found bool
		w, found, err = d.store.GetWorld(ctx, head)
		if err != nil || !found {
			t.Fatalf("GetWorld head %s: found=%v err=%v", head, found, err)
		}
	}
	c := testCommit(w, w.Revision+1, label)
	c.Objects = objects
	c.Entry.TransitionRef = transitionRef
	if err := d.store.Commit(boundedTestContext(t), c); err != nil {
		t.Fatal(err)
	}
	return c.Entry.Header.EntryIndex
}

// commitForgedChain commits the forged chain and returns its packet object and entry index.
func commitForgedChain(t *testing.T, d *Daemon, f *forgedChain) (store.Object, int64) {
	t.Helper()
	result := forgeObject(broker.EffectResultV1, []byte(`{"status":"pending","requestRef":"`+f.request.String()+`"}`))
	er := forgeObject(broker.EffectRecordV1, broker.EncodeRecord(broker.EffectRecord{Effect: broker.EffectHumanApprove, Scope: f.scope,
		Cost: 1, BudgetBefore: 5, BudgetAfter: 4, Allowed: true, Failed: f.failed, RequestRef: hashref.SumSHA256([]byte("forge-effect-request-" + f.label)), ResultRef: result.Hash}))
	recPayload, err := json.Marshal(forgedRecordV2{InvocationID: "forge-" + f.label, EpisodeID: "forge", SkillID: "forge",
		TransitionFn: f.tf.String(), Interpreter: hashref.SumSHA256([]byte("forge-interp")).String(), SemanticsEpoch: 1,
		Input: hashref.SumSHA256([]byte("forge-in-" + f.label)).String(), Output: hashref.SumSHA256([]byte("forge-out-" + f.label)).String(),
		Plan: f.planRef.String(), Effects: []string{er.Hash.String()}})
	if err != nil {
		t.Fatal(err)
	}
	record := forgeObject(coordinator.RecordV2, recPayload)
	packetPayload := f.rawPacket
	if packetPayload == nil {
		packetPayload = coordinator.EncodeDecisionPacket(f.packet)
	}
	packet := forgeObject(coordinator.DecisionPacketV1, packetPayload)
	objects := []store.Object{result, er, record, packet}
	if f.includePlan {
		objects = append(objects, f.planObject())
	}
	objects = append(objects, f.extra...)
	entry := commitObjects(t, d, f.label, objects, record.Hash)
	return packet, entry
}

// wantRefused asserts the packet page of ref shows exactly the refusal sentence and no validated-only content.
func wantRefused(t *testing.T, d *Daemon, ref hashref.HashRef, sentence string) {
	t.Helper()
	code, body := gradePage(t, d, ref)
	if code != http.StatusOK {
		t.Fatalf("status = %d", code)
	}
	sec := packetSection(t, body)
	if !strings.Contains(sec, "<p>"+html.EscapeString(sentence)+"</p>") {
		t.Fatalf("want refusal %q; section=%s", sentence, sec)
	}
	for _, s := range []string{"approvals chain:", "evidence: proposal", "committed at:"} {
		if strings.Contains(body, s) {
			t.Fatalf("refused page contains %q", s)
		}
	}
}

// wantValidated asserts the packet page of ref shows the validated sentence for entry.
func wantValidated(t *testing.T, d *Daemon, ref hashref.HashRef, entry int64) {
	t.Helper()
	code, body := gradePage(t, d, ref)
	if code != http.StatusOK {
		t.Fatalf("status = %d", code)
	}
	sec := packetSection(t, body)
	if want := "<p>" + html.EscapeString(fmt.Sprintf(packetValidatedFormat, entry)) + "</p>"; !strings.Contains(sec, want) {
		t.Fatalf("want validated at entry %d; section=%s", entry, sec)
	}
}

// TestWorkbenchPacketRefusals (AC5): one arm per refusal, each on a fresh daemon.
func TestWorkbenchPacketRefusals(t *testing.T) {
	scope := "forge"
	t.Run("0-forged-control", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		p, entry := commitForgedChain(t, d, newForgedChain("c0", q, scope))
		wantValidated(t, d, p.Hash, entry)
	})
	t.Run("1-whitespace", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c1", q, scope)
		f.rawPacket = []byte(strings.Replace(string(coordinator.EncodeDecisionPacket(f.packet)), ",\"createdAt\"", ", \"createdAt\"", 1))
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedBytes)
	})
	t.Run("1b-wrong-hash", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c1b", q, scope)
		enc := coordinator.EncodeDecisionPacket(f.packet)
		dec, err := coordinator.DecodeDecisionPacket(enc)
		if err != nil {
			t.Fatal(err)
		}
		f.rawPacket = []byte(strings.Replace(string(enc), dec.PacketHash.String(), hashref.SumSHA256([]byte("wrong")).String(), 1))
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedBytes)
	})
	t.Run("2-recarry", func(t *testing.T) {
		d := newHandlerDaemon(t)
		publishHumanTool(t, d)
		res, k := producePacket(t, d, humanPlan("e1"), "t1")
		p, err := coordinator.DecodeDecisionPacket(k.Payload)
		if err != nil {
			t.Fatal(err)
		}
		p.DeadlineAt++
		k2 := forgeObject(coordinator.DecisionPacketV1, coordinator.EncodeDecisionPacket(p))
		commitObjects(t, d, "recarry", []store.Object{k2}, res.RecordRef)
		wantRefused(t, d, k2.Hash, packetRefusedCarried)
	})
	t.Run("3-other-plan", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c3", q, scope)
		p2 := forgeObject(coordinator.EffectPlanV1, []byte(`{"forged":"other"}`))
		f.extra = []store.Object{p2}
		f.packet.ProposalHash = p2.Hash
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedProposal)
	})
	t.Run("3b-missing-plan", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c3b", q, scope)
		m := hashref.SumSHA256([]byte("missing-plan"))
		f.includePlan = false
		f.planRef = m
		f.packet.ProposalHash = m
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedProposal)
	})
	t.Run("4-other-request", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		q2 := approvalRequestInChain(t, d, scope, "forge-other", packetNow)
		f := newForgedChain("c4", q, scope)
		f.packet.RequestRef = q2
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedRequest)
	})
	t.Run("4b-scope-mismatch", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, "other-scope", "forge", packetNow)
		p, _ := commitForgedChain(t, d, newForgedChain("c4b", q, scope))
		wantRefused(t, d, p.Hash, packetRefusedRequest)
	})
	t.Run("4c-failed-effect", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c4c", q, scope)
		f.failed = true
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedRequest)
	})
	t.Run("5-unanchored", func(t *testing.T) {
		d := newHandlerDaemon(t)
		approvalRequestInChain(t, d, "unrelated", "forge", packetNow)
		q3 := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.Approve","scope":"forge","cost":1,"requester":"forge-unanchored","now":1700}`))
		if err := d.store.PutObject(boundedTestContext(t), q3); err != nil {
			t.Fatal(err)
		}
		p, _ := commitForgedChain(t, d, newForgedChain("c5", q3.Hash, scope))
		wantRefused(t, d, p.Hash, packetRefusedAnchor)
	})
	t.Run("5b-bounded", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		for i := 0; i < 41; i++ {
			approvalRequestInChain(t, d, fmt.Sprintf("newer-%d", i), "forge", packetNow)
		}
		p, _ := commitForgedChain(t, d, newForgedChain("c5b", q, scope))
		wantRefused(t, d, p.Hash, packetRefusedAnchorBounded)
	})
	t.Run("6-created-at", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c6", q, scope)
		f.packet.CreatedAt = packetNow + 1
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedLifecycle)
	})
	t.Run("6b-policy", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c6b", q, scope)
		f.packet.Policy = "EscalateBounded"
		p, _ := commitForgedChain(t, d, f)
		wantRefused(t, d, p.Hash, packetRefusedLifecycle)
	})
	t.Run("7-duplicate-pair", func(t *testing.T) {
		d := newHandlerDaemon(t)
		q := approvalRequestInChain(t, d, scope, "forge", packetNow)
		f := newForgedChain("c7", q, scope)
		p, entry := commitForgedChain(t, d, f)
		wantValidated(t, d, p.Hash, entry)
		dup := f.packet
		dup.DeadlineAt++
		if err := d.store.PutObject(boundedTestContext(t), forgeObject(coordinator.DecisionPacketV1, coordinator.EncodeDecisionPacket(dup))); err != nil {
			t.Fatal(err)
		}
		wantRefused(t, d, p.Hash, packetRefusedDuplicate)
	})
	t.Run("7b-collapsed-request", func(t *testing.T) {
		d := newHandlerDaemon(t)
		publishHumanTool(t, d)
		res1, k1 := producePacket(t, d, humanPlan("e1"), "t1")
		res2, k2 := producePacket(t, d, humanPlan("e2"), "t2")
		p1, err := coordinator.DecodeDecisionPacket(k1.Payload)
		if err != nil {
			t.Fatal(err)
		}
		p2, err := coordinator.DecodeDecisionPacket(k2.Payload)
		if err != nil {
			t.Fatal(err)
		}
		if p1.RequestRef != p2.RequestRef || p1.ProposalHash == p2.ProposalHash {
			t.Fatalf("want equal requests and distinct proposals; got %s/%s and %s/%s", p1.RequestRef, p1.ProposalHash, p2.RequestRef, p2.ProposalHash)
		}
		wantValidated(t, d, k1.Hash, res1.EntryIndex)
		wantValidated(t, d, k2.Hash, res2.EntryIndex)
	})
}

// TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive (AC6): packet grades come only from row 105's projection.
func TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive(t *testing.T) {
	ctx := boundedTestContext(t)
	t.Run("configured", func(t *testing.T) {
		d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
		tf := storeGradeSubject(t, d, "packet-tf")
		storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, tf)))
		f := newForgedChain("ac6", approvalRequestInChain(t, d, "forge", "forge", packetNow), "forge")
		f.tf = tf
		p, entry := commitForgedChain(t, d, f)
		wantValidated(t, d, p.Hash, entry)
		_, body := gradePage(t, d, p.Hash)
		sec := packetSection(t, body)
		tg, err := d.objectGradeView(ctx, tf)
		if err != nil {
			t.Fatal(err)
		}
		if !tg.Available || tg.Label != workbench.GradePROVEN {
			t.Fatalf("transition grade = %+v, want PROVEN", tg)
		}
		if !strings.Contains(sec, wantPacketGradeLine("transition function", tf.String(), tg)) || !strings.Contains(sec, "— <span>PROVEN</span></p>") {
			t.Fatalf("transition line lacks PROVEN; section=%s", sec)
		}
		pg, err := d.objectGradeView(ctx, f.planRef)
		if err != nil {
			t.Fatal(err)
		}
		if pg.Unavailable != gradeReasonBoundElsewhere {
			t.Fatalf("proposal grade reason = %q, want %q", pg.Unavailable, gradeReasonBoundElsewhere)
		}
		if !strings.Contains(sec, wantPacketGradeLine("proposal", f.planRef.String(), pg)) {
			t.Fatalf("proposal line absent; section=%s", sec)
		}
	})
	t.Run("unconfigured", func(t *testing.T) {
		d := newHandlerDaemon(t)
		f := newForgedChain("ac6u", approvalRequestInChain(t, d, "forge", "forge", packetNow), "forge")
		p, entry := commitForgedChain(t, d, f)
		wantValidated(t, d, p.Hash, entry)
		_, body := gradePage(t, d, p.Hash)
		sec := packetSection(t, body)
		unavailable := "GRADE UNAVAILABLE — " + html.EscapeString(objectGradeUnavailableReason)
		if strings.Count(sec, unavailable) != 2 {
			t.Fatalf("want two unavailable grade lines; section=%s", sec)
		}
	})
}

// TestQuickstartDecisionPacketSection (S7; PD15, PD32, PD33): the executed QUICKSTART section describes what the daemon does.
func TestQuickstartDecisionPacketSection(t *testing.T) {
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "**Inspect a decision packet** (row 172).")
	if !ok {
		t.Fatal("QUICKSTART decision-packet section not found")
	}
	for _, want := range []string{"D-WORLD-73", `{"items":[]}`, "http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1",
		"/workbench?object=$PACKET", "decision packet validated", "not a validated decision packet", "approvals chain: pending", "400 BadRequest"} {
		if !strings.Contains(section, want) {
			t.Fatalf("section lacks %q", want)
		}
	}
	if !strings.Contains(section, "--grant Human.Approve=worktree:1") {
		t.Fatal("section lacks --grant Human.Approve=worktree:1")
	}
	for _, line := range strings.Split(section, "\n") {
		if strings.HasPrefix(line, "/tmp/ailang-worldd serve ") || strings.HasPrefix(line, "|") {
			t.Fatalf("section line %q is a served-daemon command or a table row", line)
		}
	}
	route := regexp.MustCompile(`http://127\.0\.0\.1:7644(/v1/objects/by-semantic-id/[^\s|"')]+)`).FindStringSubmatch(section)
	if len(route) != 2 || route[1] != "/v1/objects/by-semantic-id/world/decision-packet/v1" {
		t.Fatalf("documented route = %v", route)
	}

	// Fresh-store arm: a Human.Approve skill alone mints nothing, so the documented listing is empty.
	d0 := newHandlerDaemon(t)
	publishHumanTool(t, d0)
	if rec := requestRecorder(t, d0, http.MethodGet, route[1], nil); rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"items":[]}` {
		t.Fatalf("empty listing = %d %s", rec.Code, rec.Body)
	}

	// Produced-packet arm: the documented listing, then the documented object page.
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	res, k := producePacket(t, d, humanPlan("e1"), "t1")
	rec := requestRecorder(t, d, http.MethodGet, route[1], nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("listing status = %d; body=%s", rec.Code, rec.Body)
	}
	var listing struct {
		Items []struct {
			Hash string `json:"hash"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listing); err != nil {
		t.Fatal(err)
	}
	if len(listing.Items) == 0 || listing.Items[0].Hash != k.Hash.String() {
		t.Fatalf("listing = %+v, want first item %s", listing, k.Hash)
	}
	if !strings.Contains(section, "http://127.0.0.1:7644/workbench?object=$PACKET") {
		t.Fatal("section lacks the documented object-page URL")
	}
	page := requestRecorder(t, d, http.MethodGet, strings.Replace("/workbench?object=$PACKET", "$PACKET", listing.Items[0].Hash, 1), nil)
	if page.Code != http.StatusOK {
		t.Fatalf("object page status = %d", page.Code)
	}
	if !strings.Contains(page.Body.String(), html.EscapeString(fmt.Sprintf(packetValidatedFormat, res.EntryIndex))) {
		t.Fatalf("object page lacks the validated sentence for entry %d", res.EntryIndex)
	}
}
