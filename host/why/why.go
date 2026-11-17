// Package why is the walk behind `ailang-worldd why`, moved from cmd/ailang-worldd (row 243 V2a, design §2.1):
// it reads World only through an injected Getter, so the CLI (HTTP) and the daemon (in process) run one walk.
package why

// Row 138 M2 (design_docs/planned/w-worldd-developer-cli.md §3.3): `why`
// walks a result, an object, an invocation id or an entry back to its whole
// provenance chain, recomputing and checking every link.
//
// An MCP result carries no entry index, invocation id or record ref (V12), and
// /v1/receipts refuses a2a:/mcp: ids (V13), so a target that is not an index is
// found by scanning the log backwards from head — one log page of 500 at a
// time and one object GET per entry — bounded by --scan (default 500, max
// 5000) and a 60 s budget. The reads are unauthenticated GETs (V15); `why`
// trusts that read surface (R-CLI-2) but verifies every content address it is
// handed against the bytes it is handed.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	ScanDefault = 500
	ScanMax     = 5000
	LogPageMax  = 500
	// MaxResponseBytes is the CLI client's cap on one response (maxClientResponseBytes).
	MaxResponseBytes = 8388608
	// The two statuses a read distinguishes (host/why imports no net/http).
	statusOK       = 200
	statusNotFound = 404
	exitOK         = 0 // the CLI's exits: every link ✓
	exitUsage      = 1 // usage, or not found
	exitIntegrity  = 3 // any link ✗
)

// ---------------------------------------------------------------------------
// the read surface
// ---------------------------------------------------------------------------

type Header struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
}

type Entry struct {
	Header        Header `json:"header"`
	EntryHash     string `json:"entryHash"`
	TransitionRef string `json:"transitionRef"`
}

type Object struct {
	Hash       string `json:"hash"`
	SemanticID string `json:"semanticId"`
	Payload    []byte `json:"payload"` // base64 on the wire
}

type World struct {
	Ref       string `json:"ref"`
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

// Record is a world/invocation-record/v1 or /v2 payload.
type Record struct {
	SemanticID     string   `json:"-"`
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
	payload        []byte
}

// Getter reads one GET path of the daemon's read surface: its status and body. A 404 is an
// answer, not an error; err is a transport or bound failure.
type Getter func(ctx context.Context, path string) (status int, body []byte, err error)

// Reader is the GET side of the daemon, under one command budget.
type Reader struct {
	Ctx         context.Context
	Get         Getter
	HealthCache *Health
}

// getJSON decodes a 200 into v; a 404 is (false, nil); anything else errors.
func (r Reader) GetJSON(path string, v any) (bool, error) {
	if err := r.Ctx.Err(); err != nil {
		return false, err
	}
	status, body, err := r.Get(r.Ctx, path)
	if err != nil {
		return false, err
	}
	switch status {
	case statusOK:
		if err := json.Unmarshal(body, v); err != nil {
			return false, fmt.Errorf("GET %s: malformed body: %v", path, err)
		}
		return true, nil
	case statusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("GET %s returned HTTP %d: %s", path, status, strings.TrimSpace(string(body)))
	}
}

func (r Reader) Entry(i int64) (Entry, bool, error) {
	var e Entry
	ok, err := r.GetJSON("/v1/log/"+strconv.FormatInt(i, 10), &e)
	return e, ok, err
}

func (r Reader) Page(from int64, limit int) ([]Entry, error) {
	var p struct {
		Items []Entry `json:"items"`
	}
	_, err := r.GetJSON(fmt.Sprintf("/v1/log?from=%d&limit=%d", from, limit), &p)
	return p.Items, err
}

func (r Reader) Object(ref string) (Object, bool, error) {
	var o Object
	ok, err := r.GetJSON("/v1/objects/"+url.PathEscape(ref)+"?payload=true", &o)
	return o, ok, err
}

func (r Reader) World(ref string) (World, bool, error) {
	var w World
	ok, err := r.GetJSON("/v1/worlds/"+url.PathEscape(ref), &w)
	return w, ok, err
}

// headIndex is the selected world's revision (= its log head's entry index).
func (r Reader) HeadIndex() (int64, bool, error) {
	ref, ok, err := Head(r.Ctx, r.Get)
	if err != nil || !ok {
		return 0, ok, err
	}
	w, found, err := r.World(ref)
	if err != nil {
		return 0, false, err
	}
	if !found {
		return 0, false, fmt.Errorf("the selected head %s has no world row", ref)
	}
	return w.Revision, true, nil
}

// Head reads GET /v1/head: the selected world ref, or present=false on
// the 404 "no world head has been selected yet".
func Head(ctx context.Context, get Getter) (ref string, present bool, err error) {
	status, body, err := get(ctx, "/v1/head")
	if err != nil {
		return "", false, err
	}
	switch status {
	case statusOK:
		return strings.TrimSpace(string(body)), true, nil
	case statusNotFound:
		return "", false, nil
	default:
		return "", false, fmt.Errorf("GET /v1/head returned HTTP %d: %s", status, strings.TrimSpace(string(body)))
	}
}

type Health struct {
	DBPath             string `json:"db_path"`
	InterpreterRef     string `json:"interpreter_ref"`
	InterpreterVersion string `json:"interpreter_version"`
}

func (r Reader) Health() (Health, error) {
	if r.HealthCache != nil {
		return *r.HealthCache, nil
	}
	var h Health
	_, err := r.GetJSON("/v1/health", &h)
	return h, err
}

// record fetches an entry's transitionRef and decodes it when it is a
// coordinator invocation record; rec is nil for any other entry.
func (r Reader) Record(e Entry) (*Record, Object, error) {
	obj, ok, err := r.Object(e.TransitionRef)
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("entry %d: transitionRef %s does not resolve", e.Header.EntryIndex, e.TransitionRef)
		}
		return nil, obj, err
	}
	if obj.SemanticID != coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2 {
		return nil, obj, nil
	}
	var rec Record
	if err := json.Unmarshal(obj.Payload, &rec); err != nil {
		return nil, obj, fmt.Errorf("entry %d: undecodable %s: %v", e.Header.EntryIndex, obj.SemanticID, err)
	}
	rec.SemanticID, rec.payload = obj.SemanticID, obj.Payload
	return &rec, obj, nil
}

// errScanExhausted reports a bounded scan that ran out of entries or budget.
var errScanExhausted = errors.New("scan exhausted")

// scanEntries visits the entries [from, to] — backwards from `to` when
// reverse — one page of up to 500 at a time, decoding each entry's record
// (one object GET per entry). It stops at the first visit that returns
// true, after max entries, or when the budget ends.
func (r Reader) ScanEntries(from, to int64, reverse bool, max int,
	visit func(e Entry, rec *Record) (bool, error)) (int, error) {
	seen := 0
	for lo, hi := from, to; lo <= hi && seen < max; {
		pFrom := lo
		if reverse {
			pFrom = hi - LogPageMax + 1
			if pFrom < lo {
				pFrom = lo
			}
		}
		n := hi - pFrom + 1
		if !reverse && n > LogPageMax {
			n = LogPageMax
		}
		items, err := r.Page(pFrom, int(n))
		if err != nil {
			return seen, err
		}
		if len(items) == 0 {
			return seen, nil
		}
		for k := range items {
			i := k
			if reverse {
				i = len(items) - 1 - k
			}
			if seen >= max {
				break
			}
			seen++
			rec, _, err := r.Record(items[i])
			if err != nil {
				return seen, err
			}
			stop, err := visit(items[i], rec)
			if err != nil || stop {
				return seen, err
			}
		}
		if reverse {
			hi = pFrom - 1
		} else {
			lo = pFrom + int64(len(items))
		}
	}
	return seen, nil
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// why
// ---------------------------------------------------------------------------

// Target is what a walk starts from. Kind is index, head, ref, invocation or
// result; a result's bytes are read by the CLI (--result, -), never here.
type Target struct {
	Kind   string
	Index  int64
	Text   string
	Result []byte
}

// ParseTarget parses a target argument other than - and --result (the two
// reads the CLI makes at its boundary), with the CLI's error texts.
func ParseTarget(arg string) (Target, error) {
	switch {
	case arg == "":
		return Target{}, errors.New("a target is required")
	case arg == "head":
		return Target{Kind: "head"}, nil
	case strings.HasPrefix(arg, "a2a:") || strings.HasPrefix(arg, "mcp:") || strings.HasPrefix(arg, "rest:"):
		return Target{Kind: "invocation", Text: arg}, nil
	case strings.HasPrefix(arg, "sha256:"):
		if _, err := hashref.Parse(arg); err != nil {
			return Target{}, fmt.Errorf("%q is not a content address: %v", arg, err)
		}
		return Target{Kind: "ref", Text: arg}, nil
	}
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || n < 0 {
		return Target{}, fmt.Errorf("%q is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -", arg)
	}
	return Target{Kind: "index", Index: n}, nil
}

// Options are a walk's --scan and --json.
type Options struct {
	Scan int
	JSON bool
}

// Answer is exactly what the CLI writes: its stdout, its stderr, its exit.
type Answer struct {
	Stdout, Stderr string
	Exit           int
}

// Run walks t to its chain and renders it as `ailang-worldd why` does: exit 0
// when every link is ✓, 3 when any is ✗, 1 when t is not found.
func Run(r Reader, t Target, o Options) (a Answer) {
	var stdout, stderr strings.Builder
	defer func() { a.Stdout, a.Stderr = stdout.String(), stderr.String() }()
	e, matchedBy, err := r.Locate(t, o.Scan)
	if err != nil {
		if errors.Is(err, errScanExhausted) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Fprintf(&stderr, "ailang-worldd why: not found in the last %d entries (%s)\n", o.Scan, matchedBy)
			return Answer{Exit: exitUsage}
		}
		fmt.Fprintf(&stderr, "ailang-worldd why: %v\n", err)
		return Answer{Exit: exitUsage}
	}
	ch, err := r.BuildChain(e, t.Result)
	if err != nil {
		fmt.Fprintf(&stderr, "ailang-worldd why: %v\n", err)
		return Answer{Exit: exitUsage}
	}
	ch.MatchedBy = matchedBy
	if o.JSON {
		b, _ := json.MarshalIndent(ch, "", "  ")
		fmt.Fprintln(&stdout, string(b))
	} else {
		printChain(&stdout, ch)
	}
	if !ch.OK {
		return Answer{Exit: exitIntegrity}
	}
	return Answer{Exit: exitOK}
}

// locate resolves a target to one log entry. matchedBy names how.
func (r Reader) Locate(t Target, scan int) (Entry, string, error) {
	byIndex := func(i int64, how string) (Entry, string, error) {
		e, ok, err := r.Entry(i)
		if err != nil {
			return e, how, err
		}
		if !ok {
			return e, how, fmt.Errorf("log entry %d not found", i)
		}
		return e, how, nil
	}
	switch t.Kind {
	case "index":
		return byIndex(t.Index, "entry index")
	case "head":
		h, ok, err := r.HeadIndex()
		if err != nil {
			return Entry{}, "head", err
		}
		if !ok {
			return Entry{}, "head", errors.New("no world head has been selected yet")
		}
		return byIndex(h, "head")
	case "invocation":
		if strings.HasPrefix(t.Text, "rest:") {
			var rc struct {
				State     string `json:"state"`
				ResultRef string `json:"resultRef"`
			}
			if _, err := r.GetJSON("/v1/receipts/"+url.PathEscape(t.Text), &rc); err != nil {
				return Entry{}, "rest receipt", err
			}
			if rc.State != "resolved" || rc.ResultRef == "" {
				return Entry{}, "rest receipt", fmt.Errorf("receipt %s is %s: nothing committed", t.Text, rc.State)
			}
			return r.ByWorld(rc.ResultRef, "rest receipt → world")
		}
		return r.ScanFor(scan, "invocationId", func(e Entry, rec *Record) bool {
			return rec != nil && rec.InvocationID == t.Text
		})
	case "result":
		ref := hashref.SumSHA256(t.Result).String()
		e, how, err := r.ScanFor(scan, "output (sha256 of the result)", func(e Entry, rec *Record) bool {
			return rec != nil && rec.Output == ref
		})
		if err == nil || !errors.Is(err, errScanExhausted) {
			return e, how, err
		}
		var out struct {
			World struct {
				Plan string `json:"plan"`
			} `json:"world"`
		}
		if json.Unmarshal(t.Result, &out) != nil || out.World.Plan == "" {
			return e, how, err
		}
		return r.ScanFor(scan, "world.plan (fallback: the result bytes are not a committed output)", func(e Entry, rec *Record) bool {
			return rec != nil && rec.Plan == out.World.Plan
		})
	}
	// A content address: dispatch on what it names.
	obj, ok, err := r.Object(t.Text)
	if err != nil {
		return Entry{}, "object", err
	}
	if !ok {
		return r.ByWorld(t.Text, "world ref")
	}
	ref := t.Text
	switch obj.SemanticID {
	case coordinator.RecordV1, coordinator.RecordV2:
		return r.ScanFor(scan, "record", func(e Entry, _ *Record) bool { return e.TransitionRef == ref })
	case coordinator.OutputV1:
		return r.ScanFor(scan, "output", func(_ Entry, rec *Record) bool { return rec != nil && rec.Output == ref })
	case coordinator.InputV1:
		return r.ScanFor(scan, "input", func(_ Entry, rec *Record) bool { return rec != nil && rec.Input == ref })
	case coordinator.EffectPlanV1:
		return r.ScanFor(scan, "plan", func(_ Entry, rec *Record) bool { return rec != nil && rec.Plan == ref })
	case broker.EffectRecordV1:
		return r.ScanFor(scan, "effect record", func(_ Entry, rec *Record) bool { return rec != nil && contains(rec.Effects, ref) })
	case broker.EffectResultV1, broker.EffectRequestV1:
		return r.ScanFor(scan, "effect "+strings.TrimSuffix(strings.TrimPrefix(obj.SemanticID, "world/effect-"), "/v1"),
			func(_ Entry, rec *Record) bool {
				if rec == nil {
					return false
				}
				for _, er := range rec.Effects {
					if d, ok := r.EffectRecord(er); ok && (d.ResultRef.String() == ref || d.RequestRef.String() == ref) {
						return true
					}
				}
				return false
			})
	default:
		return r.ScanFor(scan, "transitionRef ("+obj.SemanticID+")", func(e Entry, _ *Record) bool { return e.TransitionRef == ref })
	}
}

func (r Reader) EffectRecord(ref string) (broker.EffectRecord, bool) {
	obj, ok, err := r.Object(ref)
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
		return broker.EffectRecord{}, false
	}
	d, err := broker.DecodeRecord(obj.Payload)
	return d, err == nil
}

// byWorld resolves a world ref to the entry its log head names.
func (r Reader) ByWorld(ref, how string) (Entry, string, error) {
	w, ok, err := r.World(ref)
	if err != nil {
		return Entry{}, how, err
	}
	if !ok {
		return Entry{}, how, fmt.Errorf("%s is neither an object nor a world in this store", ref)
	}
	e, ok, err := r.Entry(w.Revision)
	if err != nil {
		return e, how, err
	}
	if !ok || e.EntryHash != w.LogHead {
		return e, how, fmt.Errorf("world %s names log head %s, which is not entry %d", ref, w.LogHead, w.Revision)
	}
	return e, how, nil
}

// scanFor scans back from head for the first entry match accepts.
func (r Reader) ScanFor(scan int, how string, match func(Entry, *Record) bool) (Entry, string, error) {
	h, ok, err := r.HeadIndex()
	if err != nil {
		return Entry{}, how, err
	}
	if !ok {
		return Entry{}, how, errScanExhausted
	}
	var found *Entry
	_, err = r.ScanEntries(0, h, true, scan, func(e Entry, rec *Record) (bool, error) {
		if match(e, rec) {
			found = &e
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return Entry{}, how, err
	}
	if found == nil {
		return Entry{}, how, errScanExhausted
	}
	return *found, how, nil
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// the chain
// ---------------------------------------------------------------------------

type Link struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Ref    string `json:"ref,omitempty"`
	Detail string `json:"detail"`
	// Worktree (row 184 D3) is an effect result's worktree_delta: the files
	// the toolchain call created, modified or deleted in the worktree.
	Worktree   *broker.WorktreeDelta `json:"worktree,omitempty"`
	Protected  json.RawMessage       `json:"protected,omitempty"`
	ExecResult bool                  `json:"-"`
}

// Keep the complete stored block for JSON; nil alone means the key is absent.
func resultProtected(payload []byte) json.RawMessage {
	var res struct {
		Protected json.RawMessage `json:"protected"`
	}
	if json.Unmarshal(payload, &res) != nil {
		return nil
	}
	return res.Protected
}

// whyProtected mirrors broker's private protectedRecord for text decoding.
type whyProtected struct {
	Digest            string   `json:"digest"`
	Patterns          []string `json:"patterns"`
	Enforced          string   `json:"enforced"`
	Base              string   `json:"base"`
	BaseDigest        string   `json:"base_digest"`
	BaseFiles         int      `json:"base_files"`
	BaseIgnoredDigest string   `json:"base_ignored_digest"`
	BaseIgnoredFiles  int      `json:"base_ignored_files"`
	IntactBefore      *bool    `json:"intact_before"`
	IntactAfter       *bool    `json:"intact_after"`
	Changes           []struct {
		Path          string `json:"path"`
		Status        string `json:"status"`
		Phase         string `json:"phase"`
		PathTruncated bool   `json:"path_truncated,omitempty"`
	} `json:"changes"`
	Truncated bool   `json:"truncated"`
	Total     int    `json:"total"`
	Error     string `json:"error"`
}

func protectedLine(raw json.RawMessage, exec bool) string {
	if raw == nil {
		if exec {
			return "protected: none declared"
		}
		return ""
	}
	var p *whyProtected
	if json.Unmarshal(raw, &p) != nil || p == nil {
		return "protected: not checked — malformed protected block"
	}
	var phases []string
	if p.IntactBefore != nil && !*p.IntactBefore {
		phases = append(phases, "before")
	}
	if p.IntactAfter != nil && !*p.IntactAfter {
		phases = append(phases, "after")
	}
	if len(phases) > 0 {
		var paths []string
		for _, c := range p.Changes {
			if len(paths) == whyWorktreePaths {
				break
			}
			paths = append(paths, fmt.Sprintf("%s (%s)", c.Path, c.Status))
		}
		line := "protected: CHANGED " + strings.Join(phases, " and ") + " the call — " + strings.Join(paths, ", ")
		if p.Total > len(paths) {
			line += fmt.Sprintf(" (+%d more)", p.Total-len(paths))
		}
		if p.Error != "" {
			line += "; not checked — " + p.Error
		}
		return line
	}
	if p.IntactBefore == nil || p.IntactAfter == nil || p.Error != "" {
		reason := p.Error
		if reason == "" {
			reason = "no verdict recorded"
		}
		return "protected: not checked — " + reason
	}
	return fmt.Sprintf("protected: intact (%d files @ %.7s, digest %s, enforced %s)", p.BaseFiles, p.Base, p.BaseDigest, p.Enforced)
}

// whyWorktreePaths is how many changed paths the text form prints.
const whyWorktreePaths = 5

// resultWorktreeDelta decodes an effect result's worktree_delta, or nil when
// the result has none (or is not a JSON object).
func resultWorktreeDelta(payload []byte) *broker.WorktreeDelta {
	var res struct {
		Delta *broker.WorktreeDelta `json:"worktree_delta"`
	}
	if json.Unmarshal(payload, &res) != nil {
		return nil
	}
	return res.Delta
}

// worktreeLine renders a delta as `worktree: N created, M modified, K
// deleted` plus the first paths (D3), or the snapshot's error or cap.
func worktreeLine(d *broker.WorktreeDelta) string {
	switch {
	case d.Error != "":
		return "worktree: not recorded — " + d.Error
	case d.Basis == broker.WorktreeBasisWalkCapped:
		return fmt.Sprintf("worktree: not recorded — the non-git worktree has more than %d files (basis %s)", d.CappedAt, d.Basis)
	}
	// Counts are of the listed changes; a truncated delta says so.
	var created, modified, deleted int
	paths := make([]string, 0, whyWorktreePaths)
	for _, c := range d.Changes {
		switch c.Op {
		case "created":
			created++
		case "modified":
			modified++
		case "deleted":
			deleted++
		}
		if len(paths) < whyWorktreePaths {
			paths = append(paths, c.Op[:1]+" "+c.Path)
		}
	}
	line := fmt.Sprintf("worktree: %d created, %d modified, %d deleted (basis %s)", created, modified, deleted, d.Basis)
	if d.Truncated {
		line += fmt.Sprintf(" — first %d of %d changes listed", len(d.Changes), d.Total)
	}
	if len(paths) > 0 {
		line += ": " + strings.Join(paths, ", ")
		if more := d.Total - len(paths); more > 0 {
			line += fmt.Sprintf(" (+%d more)", more)
		}
	}
	return line
}

type Chain struct {
	Entry       int64  `json:"entry"`
	MatchedBy   string `json:"matchedBy"`
	Coordinator bool   `json:"coordinator"`
	Note        string `json:"note,omitempty"`
	Links       []Link `json:"links"`
	OK          bool   `json:"ok"`
}

func (ch *Chain) add(name string, ok bool, ref, detail string, args ...any) {
	ch.Links = append(ch.Links, Link{Name: name, OK: ok, Ref: ref, Detail: fmt.Sprintf(detail, args...)})
	if !ok {
		ch.OK = false
	}
}

// EntryWire / WorldWire mirror host/coordinator/plan.go's canonical encodings
// (field order is the encoding order), so the entry hash and the world ref
// are recomputed, not trusted.
type EntryWire struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
	TransitionRef  string `json:"transitionRef"`
}

type WorldWire struct {
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

func SumJSON(v any) string {
	b, _ := json.Marshal(v)
	return hashref.SumSHA256(b).String()
}

func Sum(b []byte) string { return hashref.SumSHA256(b).String() }

// checkedObject fetches ref and reports whether it resolves, carries the
// wanted semantic id, and its payload hashes to ref.
func (r Reader) CheckedObject(ref, semantic string) (Object, bool, string) {
	obj, ok, err := r.Object(ref)
	switch {
	case err != nil:
		return obj, false, err.Error()
	case !ok:
		return obj, false, "does not resolve"
	case semantic != "" && obj.SemanticID != semantic:
		return obj, false, fmt.Sprintf("semanticId %s, want %s", obj.SemanticID, semantic)
	case Sum(obj.Payload) != ref:
		return obj, false, fmt.Sprintf("payload hashes to %s, not its ref", Sum(obj.Payload))
	}
	return obj, true, ""
}

func failNote(ok bool, why string) string {
	if ok {
		return ""
	}
	return " — ✗ " + why
}

// buildChain renders and checks one entry's provenance chain. result is the
// caller's result bytes (why - / --result), or nil.
func (r Reader) BuildChain(e Entry, result []byte) (Chain, error) {
	ch := Chain{Entry: e.Header.EntryIndex, OK: true}
	h, err := r.Health()
	if err != nil {
		return ch, err
	}
	pin := "not the daemon's pinned interpreter"
	if e.Header.Interpreter == h.InterpreterRef && h.InterpreterRef != "" {
		pin = strings.TrimSpace(strings.SplitN(h.InterpreterVersion, "\n", 2)[0]) + ", the daemon's pin"
	}
	rec, recObj, err := r.Record(e)
	if err != nil {
		return ch, err
	}
	entryLine := fmt.Sprintf("#%d writtenBy %s prev %s transitionFn %s interpreter %s (%s)",
		e.Header.EntryIndex, e.Header.WrittenBy, e.Header.PrevEntryHash, e.Header.TransitionFn, e.Header.Interpreter, pin)
	if rec == nil {
		ch.Note = fmt.Sprintf("entry %d is not a coordinator invocation (writtenBy %s; its transitionRef is a %s): "+
			"there is no record chain to walk, so only the entry and its object are shown", e.Header.EntryIndex, e.Header.WrittenBy, recObj.SemanticID)
		ch.add("entry", true, e.EntryHash, "%s (a client-supplied entry hash: not recomputed)", entryLine)
		ok := Sum(recObj.Payload) == e.TransitionRef
		ch.add("object", ok, e.TransitionRef, "%s, %d bytes%s", recObj.SemanticID, len(recObj.Payload), failNote(ok, "payload does not hash to transitionRef"))
		return ch, nil
	}
	ch.Coordinator = true

	// world: recomputed from the entry and the output, confirmed by the route.
	wref := SumJSON(WorldWire{Revision: e.Header.EntryIndex, StateRoot: rec.Output, LogHead: e.EntryHash})
	w, wok, werr := r.World(wref)
	worldOK := werr == nil && wok && w.Revision == e.Header.EntryIndex && w.StateRoot == rec.Output && w.LogHead == e.EntryHash
	why := "GET /v1/worlds/<recomputed ref> does not return this world"
	if werr != nil {
		why = werr.Error()
	}
	ch.add("world", worldOK, wref, "revision %d, stateRoot %s, logHead %s (recomputed)%s", e.Header.EntryIndex, rec.Output, e.EntryHash, failNote(worldOK, why))

	// entry: the hash is recomputed over the canonical header + transitionRef.
	eh := SumJSON(EntryWire{EntryIndex: e.Header.EntryIndex, SemanticsEpoch: e.Header.SemanticsEpoch,
		TransitionFn: e.Header.TransitionFn, Interpreter: e.Header.Interpreter, PrevEntryHash: e.Header.PrevEntryHash,
		WrittenBy: e.Header.WrittenBy, TransitionRef: e.TransitionRef})
	ch.add("entry", eh == e.EntryHash, e.EntryHash, "%s%s", entryLine, failNote(eh == e.EntryHash, "entry hash recomputes to "+eh))

	// record: its hash, and its fields agree with the entry.
	recOK := Sum(rec.payload) == e.TransitionRef && rec.TransitionFn == e.Header.TransitionFn &&
		rec.Interpreter == e.Header.Interpreter && rec.SemanticsEpoch == e.Header.SemanticsEpoch
	ch.add("record", recOK, e.TransitionRef, "%s invocation %s episode %s skill %s%s", rec.SemanticID, rec.InvocationID, rec.EpisodeID, rec.SkillID,
		failNote(recOK, "the record's hash or its transitionFn/interpreter/epoch disagree with the entry"))

	// input.
	in, inOK, inWhy := r.CheckedObject(rec.Input, coordinator.InputV1)
	ch.add("input", inOK, rec.Input, "%s%s", whyInputLine(in.Payload), failNote(inOK, inWhy))

	// plan and effects (v2 only).
	var planned []struct {
		ID     string `json:"id"`
		Effect string `json:"effect"`
		Scope  string `json:"scope"`
		Cost   int64  `json:"cost"`
	}
	if rec.SemanticID == coordinator.RecordV2 {
		po, pOK, pWhy := r.CheckedObject(rec.Plan, coordinator.EffectPlanV1)
		var plan struct {
			Effects []struct {
				ID     string `json:"id"`
				Effect string `json:"effect"`
				Scope  string `json:"scope"`
				Cost   int64  `json:"cost"`
			} `json:"effects"`
			Result json.RawMessage `json:"result"`
		}
		detail := "(undecodable)"
		if pOK && json.Unmarshal(po.Payload, &plan) == nil {
			planned = plan.Effects
			if len(plan.Effects) == 0 {
				detail = "no effects; result " + Elide(string(plan.Result))
			} else {
				parts := make([]string, 0, len(plan.Effects))
				for _, pe := range plan.Effects {
					parts = append(parts, fmt.Sprintf("%s %s@%s cost %d", pe.ID, pe.Effect, pe.Scope, pe.Cost))
				}
				detail = fmt.Sprintf("%d effect(s): %s", len(plan.Effects), strings.Join(parts, "; "))
			}
		}
		ch.add("plan", pOK, rec.Plan, "%s%s", detail, failNote(pOK, pWhy))
		for i, ref := range rec.Effects {
			obj, ok, oWhy := r.CheckedObject(ref, broker.EffectRecordV1)
			id := fmt.Sprintf("#%d", i)
			if i < len(planned) {
				id = planned[i].ID
			}
			if !ok {
				ch.add("effect", false, ref, "%s — ✗ %s", id, oWhy)
				continue
			}
			d, err := broker.DecodeRecord(obj.Payload)
			if err != nil {
				ch.add("effect", false, ref, "%s — ✗ %v", id, err)
				continue
			}
			verdict := "allowed"
			switch {
			case !d.Allowed:
				verdict = "denied:" + d.Denial
			case d.Failed:
				verdict = "failed"
			}
			okAll := broker.RecordConsistent(d)
			req, reqOK, _ := r.CheckedObject(d.RequestRef.String(), broker.EffectRequestV1)
			okAll = okAll && reqOK
			resText := "no result"
			if d.Failure != "" {
				// Row 255: the failed record's bounded failure line.
				resText = "no result: " + d.Failure
			}
			var delta *broker.WorktreeDelta
			var protected json.RawMessage
			execResult := false
			if !d.ResultRef.IsZero() {
				res, resOK, _ := r.CheckedObject(d.ResultRef.String(), broker.EffectResultV1)
				okAll = okAll && resOK
				resText = fmt.Sprintf("result %d B", len(res.Payload))
				if resOK {
					delta = resultWorktreeDelta(res.Payload)
					protected = resultProtected(res.Payload)
					execResult = d.Effect == broker.EffectWorkspaceExec
				}
			}
			ch.add("effect", okAll, ref, "%s %s %s@%s cost %d budget %d→%d request %d B %s%s", id, verdict, d.Effect, d.Scope, d.Cost,
				d.BudgetBefore, d.BudgetAfter, len(req.Payload), resText,
				failNote(okAll, "the record is inconsistent or its request/result does not resolve to its hash"))
			ch.Links[len(ch.Links)-1].Worktree = delta
			ch.Links[len(ch.Links)-1].Protected = protected
			ch.Links[len(ch.Links)-1].ExecResult = execResult
		}
	}

	// output: hashes to record.output, equals world.stateRoot, and its world
	// block names this record's plan and effect records.
	out, oOK, oWhy := r.CheckedObject(rec.Output, coordinator.OutputV1)
	if oOK && !(worldOK && w.StateRoot == rec.Output) {
		oOK, oWhy = false, "the output is not the committed world's stateRoot"
	}
	if oOK && rec.SemanticID == coordinator.RecordV2 {
		var ob struct {
			World struct {
				Plan    string `json:"plan"`
				Effects []struct {
					Record *string `json:"record"`
				} `json:"effects"`
			} `json:"world"`
		}
		_ = json.Unmarshal(out.Payload, &ob)
		var recs []string
		for _, fx := range ob.World.Effects {
			if fx.Record != nil {
				recs = append(recs, *fx.Record)
			}
		}
		if ob.World.Plan != rec.Plan || strings.Join(recs, ",") != strings.Join(rec.Effects, ",") {
			oOK, oWhy = false, "the output's world block does not name this record's plan and effect records"
		}
	}
	detail := fmt.Sprintf("%d B = world.stateRoot", len(out.Payload))
	if result != nil {
		if Sum(result) == rec.Output {
			detail += "; your result is these exact bytes"
		} else {
			oOK, oWhy = false, "your result ("+Sum(result)+") is not the committed output bytes"
		}
	}
	ch.add("output", oOK, rec.Output, "%s%s", detail, failNote(oOK, oWhy))
	return ch, nil
}

func printChain(w io.Writer, ch Chain) {
	fmt.Fprintf(w, "why: entry %d (matched by %s)\n", ch.Entry, ch.MatchedBy)
	if ch.Note != "" {
		fmt.Fprintf(w, "  note: %s\n", ch.Note)
	}
	bad := 0
	for _, l := range ch.Links {
		mark := "✓"
		if !l.OK {
			mark, bad = "✗", bad+1
		}
		ref := ""
		if l.Ref != "" {
			ref = l.Ref + "  "
		}
		fmt.Fprintf(w, "  %s %-7s %s%s\n", mark, l.Name, ref, l.Detail)
		if l.Worktree != nil {
			fmt.Fprintf(w, "            %s\n", worktreeLine(l.Worktree))
		}
		if line := protectedLine(l.Protected, l.ExecResult); line != "" {
			fmt.Fprintf(w, "            %s\n", line)
		}
	}
	if bad == 0 {
		fmt.Fprintf(w, "all %d link(s) verified\n", len(ch.Links))
	} else {
		fmt.Fprintf(w, "%d of %d link(s) BROKEN — integrity refusal (exit 3)\n", bad, len(ch.Links))
	}
}

// whyKeyFieldMax bounds one key field's encoded value; whyKeyFieldCount the
// number of key fields shown.
const (
	whyKeyFieldMax   = 120
	whyKeyFieldCount = 6
)

// whyInputLine renders an invocation input for `why`. A payload that fits in
// ElideAt is printed whole, unchanged. A longer one is elided — and because
// canonical JSON orders keys, the field that names WHAT a call touched can
// sit past the cut (ailang-edit's `path` follows `new_text`; row 114 F-7, a
// walk missed an edited file this way). So the short top-level fields of an
// object input — each scalar or small array whose encoding is at most
// whyKeyFieldMax bytes, in key order, at most whyKeyFieldCount — are printed
// first, then the elided payload.
func whyInputLine(payload []byte) string {
	one := strings.ReplaceAll(string(payload), "\n", `\n`)
	if len(one) <= ElideAt {
		return one
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(payload, &obj) != nil || len(obj) == 0 {
		return Elide(string(payload))
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var fields []string
	for _, k := range keys {
		v := obj[k]
		if len(v) == 0 || len(v) > whyKeyFieldMax || v[0] == '{' {
			continue
		}
		fields = append(fields, k+"="+strings.ReplaceAll(string(v), "\n", `\n`))
		if len(fields) == whyKeyFieldCount {
			break
		}
	}
	if len(fields) == 0 {
		return Elide(string(payload))
	}
	return strings.Join(fields, " ") + " · " + Elide(string(payload))
}

// ElideAt is how many bytes of a value Elide keeps.
const ElideAt = 160

// Elide shortens a value to ElideAt bytes plus its byte count, on one line.
func Elide(s string) string {
	one := strings.ReplaceAll(s, "\n", `\n`)
	if len(one) <= ElideAt {
		return one
	}
	return fmt.Sprintf("%s… (%d bytes)", strings.ToValidUTF8(one[:ElideAt], ""), len(s))
}
