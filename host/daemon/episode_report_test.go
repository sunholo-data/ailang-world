package daemon

// Row 198: the episode report on a seeded store. reportSeed commits real
// coordinator-shaped invocations (journal intent, then a Commit carrying the
// v2 record, plan, effect records and results, whose outcome lands in the
// same transaction) so the report reads exactly what a served episode leaves.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// seededEffect is one planned effect of a seeded call; result nil means no
// result object (a denied or failed effect).
type seededEffect struct {
	effect, scope, denial string
	cost, before          int64
	allowed, failed       bool
	result                map[string]any
}

type reportSeed struct {
	t    *testing.T
	d    *Daemon
	head store.World
	// input, when set, is the next call's input object (readCall sets it);
	// otherwise the input is {"task": <invocation id>}.
	input map[string]any
}

func reportObject(semantic string, payload []byte) store.Object {
	return store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(semantic)),
		SemanticID: semantic, Provenance: "episode-report-test", Payload: payload}
}

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func newReportSeed(t *testing.T) *reportSeed {
	d := newHandlerDaemon(t)
	return &reportSeed{t: t, d: d, head: seedGenesisEmbedded(t, d, "report")}
}

// intentFor builds the commit and its journal intent for one call.
func (s *reportSeed) build(surface, episode, task, tool string, effects []seededEffect, planResult map[string]any) (store.Commit, store.JournalIntent) {
	t := s.t
	id := string(surface) + ":" + episode + ":" + task
	index := s.head.Revision + 1
	in := map[string]any{"task": id}
	if s.input != nil {
		in = s.input
	}
	input := reportObject(coordinator.InputV1, mustMarshal(t, in))
	var objects []store.Object
	var recordRefs []string
	var planned []map[string]any
	for i, e := range effects {
		req := reportObject(broker.EffectRequestV1, mustMarshal(t, map[string]any{"req": id, "i": i}))
		rec := broker.EffectRecord{Effect: e.effect, Scope: e.scope, Cost: e.cost, BudgetBefore: e.before,
			BudgetAfter: e.before, Allowed: e.allowed, Failed: e.failed, Denial: e.denial, RequestRef: req.Hash}
		if e.allowed {
			rec.BudgetAfter = e.before - e.cost
		}
		objects = append(objects, req)
		if e.result != nil {
			res := reportObject(broker.EffectResultV1, mustMarshal(t, e.result))
			rec.ResultRef = res.Hash
			objects = append(objects, res)
		}
		if !broker.RecordConsistent(rec) {
			t.Fatalf("seeded effect record %+v is not consistent", rec)
		}
		ro := reportObject(broker.EffectRecordV1, broker.EncodeRecord(rec))
		objects = append(objects, ro)
		recordRefs = append(recordRefs, ro.Hash.String())
		planned = append(planned, map[string]any{"id": fmt.Sprintf("e%d", i+1), "effect": e.effect, "scope": e.scope, "cost": e.cost})
	}
	if recordRefs == nil {
		recordRefs = []string{}
	}
	if planned == nil {
		planned = []map[string]any{}
	}
	plan := reportObject(coordinator.EffectPlanV1, mustMarshal(t, map[string]any{
		"plan": "world/effect-plan/v1", "effects": planned, "finish": false, "result": planResult}))
	output := reportObject(coordinator.OutputV1, mustMarshal(t, map[string]any{"out": id}))
	fn := hashref.SumSHA256([]byte("fn-" + tool))
	interp := hashref.SumSHA256([]byte("interpreter"))
	record := reportObject(coordinator.RecordV2, mustMarshal(t, map[string]any{
		"invocationId": id, "episodeId": episode, "skillId": tool, "transitionFn": fn.String(),
		"interpreter": interp.String(), "semanticsEpoch": 1, "input": input.Hash.String(),
		"output": output.Hash.String(), "plan": plan.Hash.String(), "effects": recordRefs}))
	objects = append(objects, input, output, plan, record)
	entryHash := hashref.SumSHA256([]byte(fmt.Sprintf("entry-%d-%s", index, id)))
	next := store.World{Ref: hashref.SumSHA256([]byte(fmt.Sprintf("world-%d-%s", index, id))), Revision: index,
		StateRoot: output.Hash, LogHead: entryHash}
	c := store.Commit{InvocationID: id, ObservedHead: s.head.Ref, Objects: objects, NextWorld: next,
		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: index, SemanticsEpoch: 1, TransitionFn: fn,
			Interpreter: interp, PrevEntryHash: s.head.LogHead, WrittenBy: "coordinator:" + surface},
			EntryHash: entryHash, TransitionRef: record.Hash}}
	intent := store.JournalIntent{InvocationID: id, WorldRef: next.Ref, EntryHash: entryHash, ObservedHead: s.head.Ref,
		PrevEntryHash: s.head.LogHead, TransitionFn: fn, TransitionRef: record.Hash, Interpreter: interp,
		LogicalTime: 1_800_000_000 + index}
	return c, intent
}

// call commits one invocation and returns its entry index.
func (s *reportSeed) call(surface, episode, task, tool string, effects []seededEffect, planResult map[string]any) int64 {
	s.t.Helper()
	c, intent := s.build(surface, episode, task, tool, effects, planResult)
	ctx := boundedTestContext(s.t)
	if _, _, err := s.d.store.AppendIntent(ctx, intent.InvocationID, intent); err != nil {
		s.t.Fatalf("AppendIntent %s: %v", intent.InvocationID, err)
	}
	if err := s.d.store.Commit(ctx, c); err != nil {
		s.t.Fatalf("Commit %s: %v", intent.InvocationID, err)
	}
	s.head = c.NextWorld
	return c.Entry.Header.EntryIndex
}

// readCall commits one ailang-read call whose input object is input, as the
// MCP surface records a read's arguments.
func (s *reportSeed) readCall(episode, task string, input map[string]any, effects []seededEffect, planResult map[string]any) int64 {
	s.t.Helper()
	s.input = input
	defer func() { s.input = nil }()
	return s.call("mcp", episode, task, "ailang-read", effects, planResult)
}

// pending journals an intent whose commit never lands.
func (s *reportSeed) pending(surface, episode, task, tool string) {
	s.t.Helper()
	_, intent := s.build(surface, episode, task, tool, nil, map[string]any{"ok": true})
	if _, _, err := s.d.store.AppendIntent(boundedTestContext(s.t), intent.InvocationID, intent); err != nil {
		s.t.Fatalf("AppendIntent %s: %v", intent.InvocationID, err)
	}
}

func read(before int64, result map[string]any) []seededEffect {
	return []seededEffect{{effect: broker.EffectWorkspaceRead, scope: "worktree", cost: 1, before: before, allowed: true, result: result}}
}

// seedTwoEpisodes interleaves ep-a, ep-b and ep-a:x (whose ids share ep-a's
// "mcp:ep-a:" key range) on one daemon, the friction #6 shape.
func seedTwoEpisodes(t *testing.T) *reportSeed {
	s := newReportSeed(t)
	okRead := map[string]any{"ok": true, "content": "x"}
	ran := func(effect string, before int64, result map[string]any) []seededEffect {
		return []seededEffect{{effect: effect, scope: "worktree", cost: 1, before: before, allowed: true, result: result}}
	}
	// Entry 1.
	s.call("mcp", "ep-a", "1", "ailang-read", read(50, okRead), nil)
	// Entry 2: the other episode.
	s.call("mcp", "ep-b", "1", "ailang-write", ran(broker.EffectWorkspaceWrite, 50, map[string]any{"ok": true}), nil)
	// Entry 3: a refusal plan, no effect.
	s.call("mcp", "ep-a", "2", "ailang-read", nil, map[string]any{"ok": false, "refused": `missing required argument "path"`})
	// Entry 4: allowed, but the handler answered ok:false.
	s.call("mcp", "ep-a", "3", "ailang-cli", ran(broker.EffectAilangCLI, 50,
		map[string]any{"ok": false, "refused": "fmt write refuses: save formatted stdout with ailang-write"}), nil)
	// Entry 5: the other episode.
	s.call("mcp", "ep-b", "2", "ailang-read", read(50, okRead), nil)
	// Entry 6: a failing command is not a refusal.
	s.call("mcp", "ep-a", "4", "workspace-exec", ran(broker.EffectWorkspaceExec, 600, map[string]any{"exit_code": 1, "stdout": "FAIL"}), nil)
	// Entry 7.
	s.call("mcp", "ep-a", "5", "ailang-read", read(49, okRead), nil)
	// Entry 8: denied by the broker.
	s.call("mcp", "ep-a", "6", "workspace-exec",
		[]seededEffect{{effect: broker.EffectWorkspaceExec, scope: "worktree", cost: 1, before: 599, denial: "scope"}}, nil)
	// Entry 9: episode ep-a:x, inside ep-a's key range.
	s.call("mcp", "ep-a:x", "1", "ailang-read", read(50, okRead), nil)
	// Entry 10: ep-a over A2A.
	s.call("a2a", "ep-a", "7", "ailang-edit", ran(broker.EffectWorkspaceWrite, 50, map[string]any{"ok": true}), nil)
	// An intent whose commit never landed.
	s.pending("mcp", "ep-a", "8", "ailang-read")
	return s
}

func reportFor(t *testing.T, d *Daemon, episode string) EpisodeReport {
	t.Helper()
	ctx, cancel := context.WithTimeout(boundedTestContext(t), readDeadline)
	defer cancel()
	rep, ok, err := buildEpisodeReport(ctx, d.intents, d.reads, episode)
	if err != nil || !ok {
		t.Fatalf("report %s: ok=%v err=%v", episode, ok, err)
	}
	return rep
}

func TestEpisodeReportCountsCallsByTool(t *testing.T) {
	s := seedTwoEpisodes(t)
	rep := reportFor(t, s.d, "ep-a")
	want := map[string]int{"ailang-read": 3, "ailang-cli": 1, "workspace-exec": 2, "ailang-edit": 1}
	if rep.TotalCalls != 7 || !reflect.DeepEqual(rep.CallsByTool, want) {
		t.Fatalf("total %d calls_by_tool %v, want 7 %v", rep.TotalCalls, rep.CallsByTool, want)
	}
	if rep.Entries != [2]int64{1, 10} {
		t.Errorf("entries = %v, want [1 10]", rep.Entries)
	}
	if !reflect.DeepEqual(rep.BySurface, map[string]int{"mcp": 6, "a2a": 1}) {
		t.Errorf("calls_by_surface = %v", rep.BySurface)
	}
	if rep.Uncommitted != 1 {
		t.Errorf("uncommitted_intents = %d, want 1 (the pending intent is not a call)", rep.Uncommitted)
	}
	if rep.FirstAt != "2027-01-15T08:00:01Z" || rep.LastAt != "2027-01-15T08:00:10Z" {
		t.Errorf("first_at/last_at = %s/%s", rep.FirstAt, rep.LastAt)
	}
}

func TestEpisodeReportEffectsAndBudgets(t *testing.T) {
	s := seedTwoEpisodes(t)
	rep := reportFor(t, s.d, "ep-a")
	wantEffects := []ReportEffect{
		{broker.EffectAilangCLI, "allowed", 1},
		{broker.EffectWorkspaceExec, "allowed", 1},
		{broker.EffectWorkspaceExec, "denied:scope", 1},
		{broker.EffectWorkspaceRead, "allowed", 2},
		{broker.EffectWorkspaceWrite, "allowed", 1},
	}
	if !reflect.DeepEqual(rep.Effects, wantEffects) {
		t.Errorf("effects = %+v\nwant %+v", rep.Effects, wantEffects)
	}
	wantBudgets := []ReportBudget{
		{broker.EffectAilangCLI, "worktree", 50, 1, 49},
		{broker.EffectWorkspaceExec, "worktree", 600, 1, 599},
		{broker.EffectWorkspaceRead, "worktree", 50, 2, 48},
		{broker.EffectWorkspaceWrite, "worktree", 50, 1, 49},
	}
	if !reflect.DeepEqual(rep.Budgets, wantBudgets) {
		t.Errorf("budgets = %+v\nwant %+v", rep.Budgets, wantBudgets)
	}
}

func TestEpisodeReportPlanRefusal(t *testing.T) {
	s := seedTwoEpisodes(t)
	rep := reportFor(t, s.d, "ep-a")
	want := ReportRefusals{Count: 1, ByTool: map[string]int{"ailang-read": 1},
		First: []ReportRefusal{{Entry: 3, Tool: "ailang-read", Reason: `missing required argument "path"`}}}
	if !reflect.DeepEqual(rep.PlanRefusals, want) {
		t.Errorf("plan_refusals = %+v, want %+v", rep.PlanRefusals, want)
	}
}

func TestEpisodeReportHandlerOkFalse(t *testing.T) {
	s := seedTwoEpisodes(t)
	rep := reportFor(t, s.d, "ep-a")
	// Entry 4's effect record says allowed; only its result says ok:false.
	// Entry 6's exec exited 1 with no ok field: not a refusal.
	want := ReportRefusals{Count: 1, ByTool: map[string]int{"ailang-cli": 1},
		First: []ReportRefusal{{Entry: 4, Tool: "ailang-cli", Reason: "fmt write refuses: save formatted stdout with ailang-write"}}}
	if !reflect.DeepEqual(rep.HandlerRefusals, want) {
		t.Errorf("handler_refusals = %+v, want %+v", rep.HandlerRefusals, want)
	}
}

func TestEpisodeReportKeepsInterleavedEpisodesApart(t *testing.T) {
	s := seedTwoEpisodes(t)
	b := reportFor(t, s.d, "ep-b")
	if b.TotalCalls != 2 || b.Entries != [2]int64{2, 5} ||
		!reflect.DeepEqual(b.CallsByTool, map[string]int{"ailang-write": 1, "ailang-read": 1}) {
		t.Errorf("ep-b: total %d entries %v tools %v, want 2 [2 5] {write 1, read 1}", b.TotalCalls, b.Entries, b.CallsByTool)
	}
	x := reportFor(t, s.d, "ep-a:x")
	if x.TotalCalls != 1 || x.Entries != [2]int64{9, 9} {
		t.Errorf("ep-a:x: total %d entries %v, want 1 [9 9]", x.TotalCalls, x.Entries)
	}
	if a := reportFor(t, s.d, "ep-a"); a.CallsByTool["ailang-write"] != 0 || a.TotalCalls != 7 {
		t.Errorf("ep-a absorbed another episode's calls: %v", a.CallsByTool)
	}
	ctx, cancel := context.WithTimeout(boundedTestContext(t), readDeadline)
	defer cancel()
	if _, ok, err := buildEpisodeReport(ctx, s.d.intents, s.d.reads, "ep"); ok || err != nil {
		t.Errorf("episode ep (a prefix of ep-a, ep-b) reported calls: ok=%v err=%v", ok, err)
	}
}

func TestEpisodeReportRESTRoute(t *testing.T) {
	s := seedTwoEpisodes(t)
	rec := requestRecorder(t, s.d, http.MethodGet, "/v1/episodes/ep-a/report", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body)
	}
	want, err := json.Marshal(reportFor(t, s.d, "ep-a"))
	if err != nil {
		t.Fatal(err)
	}
	var got, wantAny any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	_ = json.Unmarshal(want, &wantAny)
	if !reflect.DeepEqual(got, wantAny) {
		t.Errorf("REST body differs from the report:\n got %s\nwant %s", rec.Body, want)
	}
	for _, key := range []string{"episode", "entries", "total_calls", "calls_by_tool", "effects", "budgets",
		"plan_refusals", "handler_refusals", "reads", "first_at", "last_at"} {
		if _, ok := got.(map[string]any)[key]; !ok {
			t.Errorf("REST body lacks %q", key)
		}
	}
	// A colon in the episode is one path segment, escaped or not.
	if rec := requestRecorder(t, s.d, http.MethodGet, "/v1/episodes/ep-a%3Ax/report", nil); rec.Code != http.StatusOK {
		t.Errorf("ep-a:x: status %d body=%s", rec.Code, rec.Body)
	}
	assertErrorClass(t, requestRecorder(t, s.d, http.MethodGet, "/v1/episodes/nobody/report", nil), http.StatusNotFound, "NotFound")
	if rec := requestRecorder(t, s.d, http.MethodPost, "/v1/episodes/ep-a/report", nil); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST: status %d, want 405", rec.Code)
	}
	// R1: as open as its GET siblings, never behind the commit gate.
	if s.d.isProtected(httptest.NewRequest(http.MethodGet, "/v1/episodes/ep-a/report", nil)) {
		t.Error("the report route is protected; row 199 closes R1 for every /v1 read together")
	}
	s.d.readDeadline = expiredReadDeadline
	assertErrorClass(t, requestRecorder(t, s.d, http.MethodGet, "/v1/episodes/ep-a/report", nil), http.StatusServiceUnavailable, "Timeout")
}

// Row 188 M0b: reads come from the ailang-read calls' recorded input objects.

var readOK = map[string]any{"ok": true, "content": "x"}

func assertReads(t *testing.T, got, want ReportReads) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("reads = %+v\nwant %+v", got, want)
	}
}

func TestEpisodeReportReadsCountPerPathInFirstEntryOrder(t *testing.T) {
	s := newReportSeed(t)
	b1 := s.readCall("r", "1", map[string]any{"path": "b.go"}, read(50, readOK), nil)
	a := s.readCall("r", "2", map[string]any{"path": "a.go"}, read(49, readOK), nil)
	s.readCall("r", "3", map[string]any{"path": "b.go"}, read(48, readOK), nil)
	// Another tool's input naming a path is not a read.
	s.input = map[string]any{"path": "c.go", "content": "y"}
	s.call("mcp", "r", "4", "ailang-write", []seededEffect{{effect: broker.EffectWorkspaceWrite, scope: "worktree",
		cost: 1, before: 50, allowed: true, result: map[string]any{"ok": true}}}, nil)
	s.input = nil
	assertReads(t, reportFor(t, s.d, "r").Reads, ReportReads{Count: 3,
		Files: []ReportRead{{Path: "b.go", Count: 2, FirstEntry: b1}, {Path: "a.go", Count: 1, FirstEntry: a}}})
}

func TestEpisodeReportRangedReadCountsAsItsPath(t *testing.T) {
	s := newReportSeed(t)
	e := s.readCall("r", "1", map[string]any{"path": "host/x.go", "offset": 10, "limit": 40}, read(50, readOK), nil)
	assertReads(t, reportFor(t, s.d, "r").Reads, ReportReads{Count: 1,
		Files: []ReportRead{{Path: "host/x.go", Count: 1, FirstEntry: e}}})
}

func TestEpisodeReportHandlerRefusedReadIsNotARead(t *testing.T) {
	s := newReportSeed(t)
	s.readCall("r", "1", map[string]any{"path": "missing.go"}, read(50, map[string]any{"ok": false, "refused": "no such file"}), nil)
	rep := reportFor(t, s.d, "r")
	assertReads(t, rep.Reads, ReportReads{Files: []ReportRead{}})
	if rep.HandlerRefusals.Count != 1 {
		t.Errorf("handler_refusals = %d, want 1", rep.HandlerRefusals.Count)
	}
}

func TestEpisodeReportPlanRefusedReadIsNotARead(t *testing.T) {
	s := newReportSeed(t)
	s.readCall("r", "1", map[string]any{"path": "a.go", "offset": -1}, nil, map[string]any{"ok": false, "refused": "offset must be positive"})
	e := s.readCall("r", "2", map[string]any{"path": "a.go"}, read(50, readOK), nil)
	rep := reportFor(t, s.d, "r")
	assertReads(t, rep.Reads, ReportReads{Count: 1, Files: []ReportRead{{Path: "a.go", Count: 1, FirstEntry: e}}})
	if rep.PlanRefusals.Count != 1 {
		t.Errorf("plan_refusals = %d, want 1", rep.PlanRefusals.Count)
	}
}

func TestEpisodeReportNonLocalReadIsUnlistedAndNeverEchoed(t *testing.T) {
	for _, path := range []string{"/etc/passwd", "../x"} {
		t.Run(path, func(t *testing.T) {
			s := newReportSeed(t)
			s.readCall("r", "1", map[string]any{"path": path}, read(50, readOK), nil)
			rep := reportFor(t, s.d, "r")
			assertReads(t, rep.Reads, ReportReads{Files: []ReportRead{}, Unlisted: 1})
			body := string(mustMarshal(t, rep))
			if strings.Contains(body, `"`+path+`"`) {
				t.Errorf("the report echoes %q: %s", path, body)
			}
		})
	}
}

// The cap is tested at a lowered value: reportReadFiles is a package
// variable, so the test commits 4 calls rather than 1001.
func TestEpisodeReportReadsCapTruncates(t *testing.T) {
	saved := reportReadFiles
	reportReadFiles = 3
	t.Cleanup(func() { reportReadFiles = saved })
	s := newReportSeed(t)
	var want []ReportRead
	for i := 0; i < reportReadFiles+1; i++ {
		p := fmt.Sprintf("f%d.go", i)
		e := s.readCall("r", fmt.Sprint(i), map[string]any{"path": p}, read(int64(50-i), readOK), nil)
		want = append(want, ReportRead{Path: p, Count: 1, FirstEntry: e})
	}
	got := reportFor(t, s.d, "r").Reads
	if len(got.Files) != reportReadFiles {
		t.Errorf("len(files) = %d, want %d", len(got.Files), reportReadFiles)
	}
	assertReads(t, got, ReportReads{Count: reportReadFiles, Files: want[:reportReadFiles], Truncated: 1})
}

func TestEpisodeReportNoReadsMarshalsAnEmptyFileList(t *testing.T) {
	s := newReportSeed(t)
	s.call("mcp", "r", "1", "ailang-write", []seededEffect{{effect: broker.EffectWorkspaceWrite, scope: "worktree",
		cost: 1, before: 50, allowed: true, result: map[string]any{"ok": true}}}, nil)
	rep := reportFor(t, s.d, "r")
	assertReads(t, rep.Reads, ReportReads{Files: []ReportRead{}})
	body := string(mustMarshal(t, rep))
	if !strings.Contains(body, `"reads":{"count":0,"files":[],"truncated":0,"unlisted":0}`) {
		t.Errorf(`want "files":[] in %s`, body)
	}
}
