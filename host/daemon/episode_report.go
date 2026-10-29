package daemon

// Row 198 (w-episode-report): GET /v1/episodes/{episode}/report summarises
// one episode from World's own store in one pass — calls by tool, effect
// decisions, budget per grant, plan refusals and handler ok:false results —
// so an operator (and the trial ledger) no longer runs one `why --json`
// subprocess per entry and parses its prose.
//
// It reads the structured records the why chain is built from, nothing else:
//
//	journal intent  "<surface>:<episode>:<task>" (+ outcome: the commit landed)
//	  → invocation record (world/invocation-record/v1|v2): episodeId, skillId
//	  → v2 plan object (zero effects: the refusal plan's result)
//	  → v2 effect records (broker.DecodeRecord): effect, scope, cost, budget,
//	    allowed/failed/denial
//	  → effect result (allowed, not failed): the handler's own JSON
//
// THE SCAN IS EPISODE-SCOPED. The coordinator journals every invocation under
// "<surface>:<episode>:<task>" (coordinator.InvocationID), and the journal's
// UNIQUE (invocation_id, kind) index makes [prefix, prefix with ";") one key
// range (store.InvocationIntents, the access path store.EffectSpend uses). So
// a report reads only that episode's journal rows, then a handful of
// primary-key lookups per call (record, world, plan or effect records, the
// effect result); the log and every other episode are never walked. The cost
// is linear in the episode's own call count, under the route's read deadline.
//
// Why not a package (S3)? It is a read projection over host-layer records
// (journal rows, broker effect records), served by the daemon that owns the
// store handle; an AILANG package cannot reach the store.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// EpisodeReport is the GET /v1/episodes/{episode}/report body and
// `ailang-worldd episode report --json`'s output.
type EpisodeReport struct {
	Start   *store.EpisodeStart `json:"start,omitempty"`
	Episode string              `json:"episode"`
	// Entries is [first, last]: the episode's first and last log entries.
	Entries         [2]int64          `json:"entries"`
	TotalCalls      int               `json:"total_calls"`
	CallsByTool     map[string]int    `json:"calls_by_tool"`
	Effects         []ReportEffect    `json:"effects"`
	Budgets         []ReportBudget    `json:"budgets"`
	PlanRefusals    ReportRefusals    `json:"plan_refusals"`
	HandlerRefusals ReportRefusals    `json:"handler_refusals"`
	Writes          ReportWrites      `json:"writes"`
	Reads           ReportReads       `json:"reads"`
	FirstAt         string            `json:"first_at"`
	LastAt          string            `json:"last_at"`
	BySurface       map[string]int    `json:"calls_by_surface"`
	Uncommitted     int               `json:"uncommitted_intents"`
	Notes           map[string]string `json:"notes"`
}

// ReportEffect counts one (effect, decision) pair. Decision is why's verdict:
// "allowed", "failed" or "denied:<denial>".
type ReportEffect struct {
	Effect   string `json:"effect"`
	Decision string `json:"decision"`
	Count    int    `json:"count"`
}

// ReportBudget is one grant's budget as the effect records show it: Left is
// the latest record's budgetAfter, Used the cost of every allowed record, and
// Minted their sum.
type ReportBudget struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Minted int64  `json:"minted"`
	Used   int64  `json:"used"`
	Left   int64  `json:"left"`
}

// ReportRefusals is a refusal class: its count, per tool, and the first few.
type ReportRefusals struct {
	Count  int             `json:"count"`
	ByTool map[string]int  `json:"by_tool"`
	First  []ReportRefusal `json:"first"`
}

// ReportRefusal is one refused call.
type ReportRefusal struct {
	Entry  int64  `json:"entry"`
	Tool   string `json:"tool"`
	Reason string `json:"reason"`
}

// ReportReads is the files the episode read, from its committed ailang-read
// calls' invocation records (row 188 M0b): the record's input object names
// the path, so the source is World's own log, the same for every agent.
type ReportReads struct {
	Count     int          `json:"count"`     // read calls listed below, summed
	Files     []ReportRead `json:"files"`     // one per distinct path, first read first
	Truncated int          `json:"truncated"` // distinct paths beyond the cap, not listed
	Unlisted  int          `json:"unlisted"`  // read calls whose path is not a local relative path
}

// ReportWrites separates explicit edits from paths produced by effect results.
type ReportWrites struct {
	Written               []ReportWrite      `json:"written"`
	Produced              []ReportRead       `json:"produced"`
	DeltasTruncated       []int64            `json:"deltas_truncated"`
	DeltaErrors           []ReportDeltaError `json:"delta_errors"`
	WrittenTruncated      int                `json:"written_truncated"`
	ProducedTruncated     int                `json:"produced_truncated"`
	Unlisted              int                `json:"unlisted"`
	DeltaTruncatedOmitted int                `json:"deltas_truncated_omitted"`
	DeltaErrorsOmitted    int                `json:"delta_errors_omitted"`
}
type ReportWrite struct {
	Path       string `json:"path"`
	Op         string `json:"op"`
	Count      int    `json:"count"`
	FirstEntry int64  `json:"first_entry"`
	lastEntry  int64
}
type ReportDeltaError struct {
	Entry int64  `json:"entry"`
	Error string `json:"error"`
}

// ReportRead is one distinct path: how often it was read or produced, and first when.
type ReportRead struct {
	Path       string `json:"path"`
	Count      int    `json:"count"`
	FirstEntry int64  `json:"first_entry"`
}

// reportReadTool is the tool whose calls are reads.
const reportReadTool = "ailang-read"

// reportReadFiles caps the distinct paths listed. A variable, not a const,
// so a test can lower it rather than commit a thousand calls.
var reportReadFiles = 1000

// reportFirstRefusals is how many refusals each class lists; reportReasonMax
// bounds one reason's length.
const (
	reportFirstRefusals = 5
	reportReasonMax     = 240
	reportIntentPage    = 500
)

// reportSurfaces are the coordinator's invocation-id namespaces.
var reportSurfaces = []coordinator.Surface{coordinator.SurfaceA2A, coordinator.SurfaceMCP}

// reportNotes say where each derived field comes from, in the body itself.
var reportNotes = map[string]string{
	"writes":           "written: successful committed ailang-write/edit/delete calls, from input.path; produced: allowed non-failed effect result worktree_delta.changes; deltas_truncated and delta_errors name evidence gaps; local relative paths only, each path list ordered by first entry and capped at 1000; excess counts and unlisted are explicit",
	"entries":          "first and last log entries of the episode's committed calls",
	"plan_refusals":    "calls whose plan named no effect and answered ok:false (refused before any effect ran)",
	"handler_refusals": "calls whose effect was allowed and ran, but whose handler result is ok:false (the effect record alone shows these as allowed)",
	"budgets":          "per (effect, scope) with at least one effect record: left = latest budgetAfter, used = cost of allowed records, minted = left + used",
	"first_at":         "coordinator logical time (unix seconds) of the first and last committed call, UTC",
	"reads":            "committed ailang-read calls with a v2 record whose every effect was allowed and did not fail and whose handler result is not ok:false; the path is the record input's \"path\" (a ranged read counts); a path that is not local and relative counts in unlisted and is never shown; files are ordered by first entry, at most 1000",
}

// intentReader is the report's journal seam; *store.Store satisfies it.
type intentReader interface {
	InvocationIntents(ctx context.Context, prefix string, afterSeq int64, limit int) ([]store.InvocationIntent, error)
}

// objectWorldReader is the subset of readStore the report reads through.
type objectWorldReader interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error)
}

// reportRecord is the part of an invocation record (v1 or v2) the report reads.
type reportRecord struct {
	EpisodeID string   `json:"episodeId"`
	SkillID   string   `json:"skillId"`
	Plan      string   `json:"plan"`
	Input     string   `json:"input"`
	Effects   []string `json:"effects"`
}

type budgetKey struct{ effect, scope string }

type budgetState struct {
	used      int64
	left      int64
	lastEntry int64
	lastPos   int
	seen      bool
}

type reportBuilder struct {
	rep      EpisodeReport
	effects  map[[2]string]int
	budgets  map[budgetKey]*budgetState
	written  map[string]*ReportWrite
	produced map[string]*ReportRead
	reads    map[string]*ReportRead
	firstAt  int64
	lastAt   int64
}

// buildEpisodeReport scans episode's journal key ranges and aggregates its
// committed calls. ok is false when the episode has no committed call.
func buildEpisodeReport(ctx context.Context, intents intentReader, reads objectWorldReader, episode string) (EpisodeReport, bool, error) {
	b := &reportBuilder{
		rep: EpisodeReport{Episode: episode, CallsByTool: map[string]int{}, BySurface: map[string]int{},
			Effects: []ReportEffect{}, Budgets: []ReportBudget{},
			PlanRefusals:    ReportRefusals{ByTool: map[string]int{}, First: []ReportRefusal{}},
			HandlerRefusals: ReportRefusals{ByTool: map[string]int{}, First: []ReportRefusal{}},
			Reads:           ReportReads{Files: []ReportRead{}},
			Writes:          ReportWrites{Written: []ReportWrite{}, Produced: []ReportRead{}, DeltasTruncated: []int64{}, DeltaErrors: []ReportDeltaError{}},
			Notes:           reportNotes},
		written: map[string]*ReportWrite{}, produced: map[string]*ReportRead{},
		effects: map[[2]string]int{},
		budgets: map[budgetKey]*budgetState{},
		reads:   map[string]*ReportRead{},
	}
	for _, surface := range reportSurfaces {
		prefix := string(surface) + ":" + episode + ":"
		after := int64(0)
		for {
			page, err := intents.InvocationIntents(ctx, prefix, after, reportIntentPage)
			if err != nil {
				return EpisodeReport{}, false, err
			}
			for _, it := range page {
				after = it.Seq
				if err := b.add(ctx, reads, string(surface), it); err != nil {
					return EpisodeReport{}, false, err
				}
			}
			if len(page) < reportIntentPage {
				break
			}
		}
	}
	if b.rep.TotalCalls == 0 {
		return b.rep, false, nil
	}
	b.finish()
	return b.rep, true, nil
}

// add folds one journal intent into the report.
func (b *reportBuilder) add(ctx context.Context, reads objectWorldReader, surface string, it store.InvocationIntent) error {
	if it.Outcome == nil {
		b.rep.Uncommitted++
		return nil
	}
	recObj, ok, err := reads.GetObject(ctx, it.Intent.TransitionRef)
	if err != nil {
		return err
	}
	if !ok || (recObj.SemanticID != coordinator.RecordV1 && recObj.SemanticID != coordinator.RecordV2) {
		return fmt.Errorf("invocation %s: transitionRef %s is not an invocation record", it.Intent.InvocationID, it.Intent.TransitionRef)
	}
	var rec reportRecord
	if err := json.Unmarshal(recObj.Payload, &rec); err != nil {
		return fmt.Errorf("invocation %s: undecodable record: %w", it.Intent.InvocationID, err)
	}
	// The key range is a prefix: "mcp:a:" also holds episode "a:b"'s ids.
	if rec.EpisodeID != b.rep.Episode {
		return nil
	}
	w, ok, err := reads.GetWorld(ctx, it.Outcome.ResultRef)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invocation %s: its committed world %s is absent", it.Intent.InvocationID, it.Outcome.ResultRef)
	}
	entry, tool := w.Revision, rec.SkillID
	b.call(entry, tool, surface, it.Intent.LogicalTime)
	if recObj.SemanticID != coordinator.RecordV2 {
		return nil
	}
	if len(rec.Effects) == 0 {
		return b.planRefusal(ctx, reads, entry, tool, rec.Plan)
	}
	handlerRefused := ""
	// A read is an ailang-read call whose every effect was allowed and did
	// not fail and whose handler did not answer ok:false (the refusal test
	// below); a plan refusal returned above, before any effect.
	isRead := tool == reportReadTool
	isWrite := tool == "ailang-write" || tool == "ailang-edit" || tool == "ailang-delete"
	writeEffect := false
	for pos, text := range rec.Effects {
		ref, err := hashref.Parse(text)
		if err != nil {
			return fmt.Errorf("entry %d: effect record ref: %w", entry, err)
		}
		obj, ok, err := reads.GetObject(ctx, ref)
		if err != nil {
			return err
		}
		if !ok || obj.SemanticID != broker.EffectRecordV1 {
			return fmt.Errorf("entry %d: effect record %s does not resolve", entry, text)
		}
		d, err := broker.DecodeRecord(obj.Payload)
		if err != nil {
			return fmt.Errorf("entry %d: %w", entry, err)
		}
		b.effect(entry, pos, d)
		if d.Effect == broker.EffectWorkspaceWrite || d.Effect == broker.EffectWorkspaceDelete {
			writeEffect = true
		}
		if d.Allowed && !d.Failed {
			if err := b.delta(ctx, reads, entry, d.ResultRef); err != nil {
				return err
			}
		}
		if !d.Allowed || d.Failed {
			isRead = false
			isWrite = false
		}
		if d.Allowed && !d.Failed && handlerRefused == "" {
			if handlerRefused, err = handlerRefusal(ctx, reads, d.ResultRef); err != nil {
				return fmt.Errorf("entry %d: %w", entry, err)
			}
		}
	}
	if handlerRefused != "" {
		addRefusal(&b.rep.HandlerRefusals, entry, tool, handlerRefused)
		return nil
	}
	if isWrite && writeEffect {
		return b.write(ctx, reads, entry, rec.Input, strings.TrimPrefix(tool, "ailang-"))
	}
	if isRead {
		return b.read(ctx, reads, entry, rec.Input)
	}
	return nil
}

// read folds one read call: the path is the record's input object's "path"
// (offset/limit and any other key are ignored, so a ranged read counts). A
// path that is not local and relative is counted, never echoed.
func (b *reportBuilder) read(ctx context.Context, reads objectWorldReader, entry int64, inputText string) error {
	ref, err := hashref.Parse(inputText)
	if err != nil {
		return fmt.Errorf("entry %d: input ref: %w", entry, err)
	}
	obj, ok, err := reads.GetObject(ctx, ref)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("entry %d: input %s does not resolve", entry, inputText)
	}
	var in struct {
		Path *string `json:"path"`
	}
	if json.Unmarshal(obj.Payload, &in) != nil || in.Path == nil || !filepath.IsLocal(*in.Path) {
		b.rep.Reads.Unlisted++
		return nil
	}
	r := b.reads[*in.Path]
	if r == nil {
		r = &ReportRead{Path: *in.Path, FirstEntry: entry}
		b.reads[*in.Path] = r
	}
	r.Count++
	if entry < r.FirstEntry {
		r.FirstEntry = entry
	}
	return nil
}

func (b *reportBuilder) write(ctx context.Context, reads objectWorldReader, entry int64, input, op string) error {
	ref, err := hashref.Parse(input)
	if err != nil {
		return err
	}
	obj, ok, err := reads.GetObject(ctx, ref)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("entry %d: missing input", entry)
	}
	var in struct {
		Path string `json:"path"`
	}
	if json.Unmarshal(obj.Payload, &in) != nil || !filepath.IsLocal(in.Path) {
		b.rep.Writes.Unlisted++
		return nil
	}
	r := b.written[in.Path]
	if r == nil {
		r = &ReportWrite{Path: in.Path, FirstEntry: entry}
		b.written[in.Path] = r
	}
	r.Count++
	if entry < r.FirstEntry {
		r.FirstEntry = entry
	}
	if entry >= r.lastEntry {
		r.Op, r.lastEntry = op, entry
	}
	return nil
}

func (b *reportBuilder) delta(ctx context.Context, reads objectWorldReader, entry int64, ref hashref.HashRef) error {
	if ref.IsZero() {
		return nil
	}
	obj, ok, err := reads.GetObject(ctx, ref)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("entry %d: missing effect result", entry)
	}
	var res struct {
		Delta *broker.WorktreeDelta `json:"worktree_delta"`
	}
	if json.Unmarshal(obj.Payload, &res) != nil || res.Delta == nil {
		return nil
	}
	d := res.Delta
	if d.Truncated {
		b.rep.Writes.DeltasTruncated = append(b.rep.Writes.DeltasTruncated, entry)
	}
	if d.Error != "" {
		b.rep.Writes.DeltaErrors = append(b.rep.Writes.DeltaErrors, ReportDeltaError{Entry: entry, Error: clipReason(d.Error)})
	}
	for _, c := range d.Changes {
		if c.PathTruncated || !filepath.IsLocal(c.Path) {
			b.rep.Writes.Unlisted++
			continue
		}
		r := b.produced[c.Path]
		if r == nil {
			r = &ReportRead{Path: c.Path, FirstEntry: entry}
			b.produced[c.Path] = r
		}
		r.Count++
		if entry < r.FirstEntry {
			r.FirstEntry = entry
		}
	}
	return nil
}

func (b *reportBuilder) call(entry int64, tool, surface string, at int64) {
	r := &b.rep
	if r.TotalCalls == 0 || entry < r.Entries[0] {
		r.Entries[0] = entry
	}
	if r.TotalCalls == 0 || entry > r.Entries[1] {
		r.Entries[1] = entry
	}
	if r.TotalCalls == 0 || at < b.firstAt {
		b.firstAt = at
	}
	if r.TotalCalls == 0 || at > b.lastAt {
		b.lastAt = at
	}
	r.TotalCalls++
	r.CallsByTool[tool]++
	r.BySurface[surface]++
}

// effect folds one effect record into the effect and budget tallies.
func (b *reportBuilder) effect(entry int64, pos int, d broker.EffectRecord) {
	verdict := "allowed"
	switch {
	case !d.Allowed:
		verdict = "denied:" + d.Denial
	case d.Failed:
		verdict = "failed"
	}
	b.effects[[2]string{d.Effect, verdict}]++
	k := budgetKey{d.Effect, d.Scope}
	s := b.budgets[k]
	if s == nil {
		s = &budgetState{}
		b.budgets[k] = s
	}
	if d.Allowed {
		s.used += d.Cost
	}
	if !s.seen || entry > s.lastEntry || (entry == s.lastEntry && pos >= s.lastPos) {
		s.left, s.lastEntry, s.lastPos, s.seen = d.BudgetAfter, entry, pos, true
	}
}

// planRefusal reads a zero-effect call's plan object: a result that is not
// ok:true is a refusal before any effect.
func (b *reportBuilder) planRefusal(ctx context.Context, reads objectWorldReader, entry int64, tool, planText string) error {
	ref, err := hashref.Parse(planText)
	if err != nil {
		return fmt.Errorf("entry %d: plan ref: %w", entry, err)
	}
	obj, ok, err := reads.GetObject(ctx, ref)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("entry %d: plan %s does not resolve", entry, planText)
	}
	var plan struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(obj.Payload, &plan); err != nil {
		return fmt.Errorf("entry %d: undecodable plan: %w", entry, err)
	}
	okField, reason := resultVerdict(plan.Result)
	if okField != nil && *okField {
		return nil
	}
	if reason == "" {
		reason = "the plan named no effect"
	}
	addRefusal(&b.rep.PlanRefusals, entry, tool, reason)
	return nil
}

// handlerRefusal reads an allowed effect's result and returns its reason when
// the handler answered ok:false ("" otherwise).
func handlerRefusal(ctx context.Context, reads objectWorldReader, ref hashref.HashRef) (string, error) {
	if ref.IsZero() {
		return "", nil
	}
	obj, ok, err := reads.GetObject(ctx, ref)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("effect result %s does not resolve", ref)
	}
	okField, reason := resultVerdict(obj.Payload)
	if okField == nil || *okField {
		return "", nil
	}
	if reason == "" {
		reason = "ok:false"
	}
	return reason, nil
}

// resultVerdict decodes a result object's top-level "ok" and the first
// reason-like string field.
func resultVerdict(raw []byte) (*bool, string) {
	var res map[string]json.RawMessage
	if json.Unmarshal(raw, &res) != nil {
		return nil, ""
	}
	var ok *bool
	if v, has := res["ok"]; has {
		var b bool
		if json.Unmarshal(v, &b) == nil {
			ok = &b
		}
	}
	for _, k := range []string{"refused", "error", "message", "reason", "stderr", "status"} {
		var s string
		if v, has := res[k]; has && json.Unmarshal(v, &s) == nil && strings.TrimSpace(s) != "" {
			return ok, clipReason(s)
		}
	}
	return ok, ""
}

func clipReason(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > reportReasonMax {
		s = s[:reportReasonMax] + "…"
	}
	return s
}

// addRefusal counts one refusal and keeps the reportFirstRefusals earliest.
func addRefusal(r *ReportRefusals, entry int64, tool, reason string) {
	r.Count++
	r.ByTool[tool]++
	r.First = append(r.First, ReportRefusal{Entry: entry, Tool: tool, Reason: reason})
	sort.SliceStable(r.First, func(i, j int) bool { return r.First[i].Entry < r.First[j].Entry })
	if len(r.First) > reportFirstRefusals {
		r.First = r.First[:reportFirstRefusals]
	}
}

// finish renders the tallies in a stable order.
func (b *reportBuilder) finish() {
	for k, n := range b.effects {
		b.rep.Effects = append(b.rep.Effects, ReportEffect{Effect: k[0], Decision: k[1], Count: n})
	}
	sort.Slice(b.rep.Effects, func(i, j int) bool {
		x, y := b.rep.Effects[i], b.rep.Effects[j]
		if x.Effect != y.Effect {
			return x.Effect < y.Effect
		}
		return x.Decision < y.Decision
	})
	for k, s := range b.budgets {
		b.rep.Budgets = append(b.rep.Budgets, ReportBudget{Effect: k.effect, Scope: k.scope, Minted: s.left + s.used, Used: s.used, Left: s.left})
	}
	sort.Slice(b.rep.Budgets, func(i, j int) bool {
		x, y := b.rep.Budgets[i], b.rep.Budgets[j]
		if x.Effect != y.Effect {
			return x.Effect < y.Effect
		}
		return x.Scope < y.Scope
	})
	for _, r := range b.reads {
		b.rep.Reads.Files = append(b.rep.Reads.Files, *r)
	}
	sort.Slice(b.rep.Reads.Files, func(i, j int) bool {
		x, y := b.rep.Reads.Files[i], b.rep.Reads.Files[j]
		if x.FirstEntry != y.FirstEntry {
			return x.FirstEntry < y.FirstEntry
		}
		return x.Path < y.Path
	})
	if n := len(b.rep.Reads.Files); n > reportReadFiles {
		b.rep.Reads.Truncated = n - reportReadFiles
		b.rep.Reads.Files = b.rep.Reads.Files[:reportReadFiles]
	}
	for _, r := range b.rep.Reads.Files {
		b.rep.Reads.Count += r.Count
	}
	for _, r := range b.written {
		b.rep.Writes.Written = append(b.rep.Writes.Written, *r)
	}
	for _, r := range b.produced {
		b.rep.Writes.Produced = append(b.rep.Writes.Produced, *r)
	}
	w := &b.rep.Writes
	sort.Slice(w.Written, func(i, j int) bool {
		if w.Written[i].FirstEntry != w.Written[j].FirstEntry {
			return w.Written[i].FirstEntry < w.Written[j].FirstEntry
		}
		return w.Written[i].Path < w.Written[j].Path
	})
	sort.Slice(w.Produced, func(i, j int) bool {
		if w.Produced[i].FirstEntry != w.Produced[j].FirstEntry {
			return w.Produced[i].FirstEntry < w.Produced[j].FirstEntry
		}
		return w.Produced[i].Path < w.Produced[j].Path
	})
	if n := len(w.Written); n > reportReadFiles {
		w.WrittenTruncated = n - reportReadFiles
		w.Written = w.Written[:reportReadFiles]
	}
	if n := len(w.Produced); n > reportReadFiles {
		w.ProducedTruncated = n - reportReadFiles
		w.Produced = w.Produced[:reportReadFiles]
	}
	sort.Slice(w.DeltasTruncated, func(i, j int) bool { return w.DeltasTruncated[i] < w.DeltasTruncated[j] })
	sort.Slice(w.DeltaErrors, func(i, j int) bool { return w.DeltaErrors[i].Entry < w.DeltaErrors[j].Entry })
	unique := w.DeltasTruncated[:0]
	for _, entry := range w.DeltasTruncated {
		if len(unique) == 0 || unique[len(unique)-1] != entry {
			unique = append(unique, entry)
		}
	}
	w.DeltasTruncated = unique
	if n := len(w.DeltasTruncated); n > reportReadFiles {
		w.DeltaTruncatedOmitted = n - reportReadFiles
		w.DeltasTruncated = w.DeltasTruncated[:reportReadFiles]
	}
	if n := len(w.DeltaErrors); n > reportReadFiles {
		w.DeltaErrorsOmitted = n - reportReadFiles
		w.DeltaErrors = w.DeltaErrors[:reportReadFiles]
	}
	b.rep.FirstAt = time.Unix(b.firstAt, 0).UTC().Format(time.RFC3339)
	b.rep.LastAt = time.Unix(b.lastAt, 0).UTC().Format(time.RFC3339)
}

// maxReportEpisodeBytes bounds the {episode} path value.
const maxReportEpisodeBytes = 256

// handleEpisodeReport serves GET /v1/episodes/{episode}/report (row 198). It
// is a read route like its GET siblings: unauthenticated under residual R1
// (isProtected; row 199 closes R1 for every /v1 read together).
func (d *Daemon) handleEpisodeReport(w http.ResponseWriter, r *http.Request) {
	episode := r.PathValue("episode")
	if episode == "" || len(episode) > maxReportEpisodeBytes {
		writeAPIError(w, "BadRequest", fmt.Sprintf("episode must be 1..%d bytes", maxReportEpisodeBytes), http.StatusBadRequest)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	rep, ok, err := buildEpisodeReport(ctx, d.intents, d.reads, episode)
	if err != nil {
		if timedOut(ctx, err) {
			writeReadTimeout(w, d.readDeadline)
			return
		}
		d.writeInternalError(w, r, err)
		return
	}
	start, startErr := d.readEpisodeStart(ctx, episode)
	if startErr != nil {
		d.writeInternalError(w, r, startErr)
		return
	}
	rep.Start = start
	if !ok && start == nil {
		writeAPIError(w, "NotFound", "episode "+episode+" has no committed calls", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, rep)
}
