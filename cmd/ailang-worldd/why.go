package main

// Row 138 M2 (design_docs/planned/w-worldd-developer-cli.md §3.3): `why`
// walks a result, an object, an invocation id or an entry back to its whole
// provenance chain, recomputing and checking every link.
//
// An MCP result carries no entry index, invocation id or record ref (V12), and
// /v1/receipts refuses a2a: ids (V13), so a target that is not an index is
// found by scanning the log backwards from head — one log page of 500 at a
// time and one object GET per entry — bounded by --scan (default 500, max
// 5000) and a 60 s budget. The reads are unauthenticated GETs (V15); `why`
// trusts that read surface (R-CLI-2) but verifies every content address it is
// handed against the bytes it is handed.

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	whyBudget       = 60 * time.Second
	whyScanDefault  = 500
	whyScanMax      = 5000
	logPageMax      = 500
	coordinatorAuth = "coordinator:a2a" // the coordinator's writtenBy (host/coordinator/plan.go)
)

// ---------------------------------------------------------------------------
// the read surface
// ---------------------------------------------------------------------------

type wireHeader struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
}

type wireEntry struct {
	Header        wireHeader `json:"header"`
	EntryHash     string     `json:"entryHash"`
	TransitionRef string     `json:"transitionRef"`
}

type wireObject struct {
	Hash       string `json:"hash"`
	SemanticID string `json:"semanticId"`
	Payload    []byte `json:"payload"` // base64 on the wire
}

type wireWorld struct {
	Ref       string `json:"ref"`
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

// invRecord is a world/invocation-record/v1 or /v2 payload.
type invRecord struct {
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

// worldReader is the GET side of the daemon, under one command budget.
type worldReader struct {
	ctx context.Context
	c   *client
}

// getJSON decodes a 200 into v; a 404 is (false, nil); anything else errors.
func (r worldReader) getJSON(path string, v any) (bool, error) {
	if err := r.ctx.Err(); err != nil {
		return false, err
	}
	status, body, err := r.c.do(r.ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, err
	}
	switch status {
	case http.StatusOK:
		if err := json.Unmarshal(body, v); err != nil {
			return false, fmt.Errorf("GET %s: malformed body: %v", path, err)
		}
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("GET %s returned HTTP %d: %s", path, status, strings.TrimSpace(string(body)))
	}
}

func (r worldReader) entry(i int64) (wireEntry, bool, error) {
	var e wireEntry
	ok, err := r.getJSON("/v1/log/"+strconv.FormatInt(i, 10), &e)
	return e, ok, err
}

func (r worldReader) page(from int64, limit int) ([]wireEntry, error) {
	var p struct {
		Items []wireEntry `json:"items"`
	}
	_, err := r.getJSON(fmt.Sprintf("/v1/log?from=%d&limit=%d", from, limit), &p)
	return p.Items, err
}

func (r worldReader) object(ref string) (wireObject, bool, error) {
	var o wireObject
	ok, err := r.getJSON("/v1/objects/"+url.PathEscape(ref)+"?payload=true", &o)
	return o, ok, err
}

func (r worldReader) world(ref string) (wireWorld, bool, error) {
	var w wireWorld
	ok, err := r.getJSON("/v1/worlds/"+url.PathEscape(ref), &w)
	return w, ok, err
}

// headIndex is the selected world's revision (= its log head's entry index).
func (r worldReader) headIndex() (int64, bool, error) {
	ref, ok, err := headState(r.ctx, r.c)
	if err != nil || !ok {
		return 0, ok, err
	}
	w, found, err := r.world(ref)
	if err != nil {
		return 0, false, err
	}
	if !found {
		return 0, false, fmt.Errorf("the selected head %s has no world row", ref)
	}
	return w.Revision, true, nil
}

type healthInfo struct {
	DBPath             string `json:"db_path"`
	InterpreterRef     string `json:"interpreter_ref"`
	InterpreterVersion string `json:"interpreter_version"`
}

func (r worldReader) health() (healthInfo, error) {
	var h healthInfo
	_, err := r.getJSON("/v1/health", &h)
	return h, err
}

// record fetches an entry's transitionRef and decodes it when it is a
// coordinator invocation record; rec is nil for any other entry.
func (r worldReader) record(e wireEntry) (*invRecord, wireObject, error) {
	obj, ok, err := r.object(e.TransitionRef)
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("entry %d: transitionRef %s does not resolve", e.Header.EntryIndex, e.TransitionRef)
		}
		return nil, obj, err
	}
	if obj.SemanticID != coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2 {
		return nil, obj, nil
	}
	var rec invRecord
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
func (r worldReader) scanEntries(from, to int64, reverse bool, max int,
	visit func(e wireEntry, rec *invRecord) (bool, error)) (int, error) {
	seen := 0
	for lo, hi := from, to; lo <= hi && seen < max; {
		pFrom := lo
		if reverse {
			pFrom = hi - logPageMax + 1
			if pFrom < lo {
				pFrom = lo
			}
		}
		n := hi - pFrom + 1
		if !reverse && n > logPageMax {
			n = logPageMax
		}
		items, err := r.page(pFrom, int(n))
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
			rec, _, err := r.record(items[i])
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
// why
// ---------------------------------------------------------------------------

const whyHelp = `usage: ailang-worldd [--addr <url>] why <target> [--scan N] [--json]
       ailang-worldd [--addr <url>] why --result <file> [--scan N] [--json]

Walks one committed invocation back to its whole provenance chain, checking
every link (✓ or ✗):

  world    the world ref, recomputed from the entry and the output, and
           confirmed by GET /v1/worlds/<ref>
  entry    index, writtenBy, prev, transitionFn and the interpreter (with its
           pin name); the entry hash is recomputed
  record   invocation id, episode, skill (its hash and its fields agree with
           the entry)
  input    your arguments, decoded and elided
  plan     the effects planned, or the refusal result
  effect   per effect record: allowed/failed/denial, budget before→after,
           request and result sizes
  output   the committed output; it must equal world.stateRoot

Targets:
  <index> | head               a log entry
  sha256:<ref>                 any object (dispatched on its semanticId: a
                               record, input, output, plan, effect record,
                               effect request or result), or a world ref
  a2a:<id> | rest:<id>         an invocation id
  - | --result <file>          the output of 'call --json-out': its sha256 is
                               the output ref (the world.plan is the fallback)

A target that is not an index is found by scanning the log backwards from
head: --scan N entries (default 500, max 5000) within a 60 s budget; past
either it prints "not found in the last N entries" and exits 1.

Exit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or
usage. --json prints the chain as JSON.
`

type whyTarget struct {
	kind   string // index, head, ref, invocation, result
	index  int64
	text   string
	result []byte
}

func parseWhyTarget(arg, resultFile string, stdin io.Reader) (whyTarget, error) {
	if resultFile != "" {
		if arg != "" {
			return whyTarget{}, errors.New("give either a target or --result, not both")
		}
		b, err := os.ReadFile(resultFile)
		if err != nil {
			return whyTarget{}, fmt.Errorf("read --result %s: %v", resultFile, errors.Unwrap(err))
		}
		return whyTarget{kind: "result", result: trimOneNewline(b)}, nil
	}
	switch {
	case arg == "":
		return whyTarget{}, errors.New("a target is required")
	case arg == "-":
		b, err := io.ReadAll(io.LimitReader(stdin, maxClientResponseBytes+1))
		if err != nil {
			return whyTarget{}, fmt.Errorf("read stdin: %v", err)
		}
		return whyTarget{kind: "result", result: trimOneNewline(b)}, nil
	case arg == "head":
		return whyTarget{kind: "head"}, nil
	case strings.HasPrefix(arg, "a2a:") || strings.HasPrefix(arg, "rest:"):
		return whyTarget{kind: "invocation", text: arg}, nil
	case strings.HasPrefix(arg, "sha256:"):
		if _, err := hashref.Parse(arg); err != nil {
			return whyTarget{}, fmt.Errorf("%q is not a content address: %v", arg, err)
		}
		return whyTarget{kind: "ref", text: arg}, nil
	}
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || n < 0 {
		return whyTarget{}, fmt.Errorf("%q is not an entry index, head, sha256:<ref>, a2a:/rest:<id>, or -", arg)
	}
	return whyTarget{kind: "index", index: n}, nil
}

// trimOneNewline drops the one newline `call --json-out` appends.
func trimOneNewline(b []byte) []byte {
	if n := len(b); n > 0 && b[n-1] == '\n' {
		return b[:n-1]
	}
	return b
}

func runWhy(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	target := ""
	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
		target, args = args[0], args[1:]
	}
	fs := flag.NewFlagSet("ailang-worldd why", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	resultFile := fs.String("result", "", "a 'call --json-out' output file")
	scan := fs.Int("scan", whyScanDefault, "entries to scan back from head (max 5000)")
	asJSON := fs.Bool("json", false, "print the chain as JSON")
	if code, done := parseVerbFlags(fs, args, whyHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() > 0 {
		if target != "" || fs.NArg() > 1 {
			fmt.Fprintf(stderr, "ailang-worldd why: unexpected argument %q\n", fs.Arg(fs.NArg()-1))
			return exitUsage
		}
		target = fs.Arg(0)
	}
	if *scan < 1 || *scan > whyScanMax {
		fmt.Fprintf(stderr, "ailang-worldd why: --scan must be 1..%d\n", whyScanMax)
		return exitUsage
	}
	t, err := parseWhyTarget(target, *resultFile, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd why: %v\n%s", err, whyHelp)
		return exitUsage
	}
	ctx, cancel := budgetContext(whyBudget)
	defer cancel()
	r := worldReader{ctx: ctx, c: newClient(addr)}
	e, matchedBy, err := r.locate(t, *scan)
	if err != nil {
		if errors.Is(err, errScanExhausted) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Fprintf(stderr, "ailang-worldd why: not found in the last %d entries (%s)\n", *scan, matchedBy)
			return exitUsage
		}
		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
		return exitUsage
	}
	ch, err := r.buildChain(e, t.result)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
		return exitUsage
	}
	ch.MatchedBy = matchedBy
	if *asJSON {
		b, _ := json.MarshalIndent(ch, "", "  ")
		fmt.Fprintln(stdout, string(b))
	} else {
		printChain(stdout, ch)
	}
	if !ch.OK {
		return exitIntegrity
	}
	return exitOK
}

// locate resolves a target to one log entry. matchedBy names how.
func (r worldReader) locate(t whyTarget, scan int) (wireEntry, string, error) {
	byIndex := func(i int64, how string) (wireEntry, string, error) {
		e, ok, err := r.entry(i)
		if err != nil {
			return e, how, err
		}
		if !ok {
			return e, how, fmt.Errorf("log entry %d not found", i)
		}
		return e, how, nil
	}
	switch t.kind {
	case "index":
		return byIndex(t.index, "entry index")
	case "head":
		h, ok, err := r.headIndex()
		if err != nil {
			return wireEntry{}, "head", err
		}
		if !ok {
			return wireEntry{}, "head", errors.New("no world head has been selected yet")
		}
		return byIndex(h, "head")
	case "invocation":
		if strings.HasPrefix(t.text, "rest:") {
			var rc struct {
				State     string `json:"state"`
				ResultRef string `json:"resultRef"`
			}
			if _, err := r.getJSON("/v1/receipts/"+url.PathEscape(t.text), &rc); err != nil {
				return wireEntry{}, "rest receipt", err
			}
			if rc.State != "resolved" || rc.ResultRef == "" {
				return wireEntry{}, "rest receipt", fmt.Errorf("receipt %s is %s: nothing committed", t.text, rc.State)
			}
			return r.byWorld(rc.ResultRef, "rest receipt → world")
		}
		return r.scanFor(scan, "invocationId", func(e wireEntry, rec *invRecord) bool {
			return rec != nil && rec.InvocationID == t.text
		})
	case "result":
		ref := hashref.SumSHA256(t.result).String()
		e, how, err := r.scanFor(scan, "output (sha256 of the result)", func(e wireEntry, rec *invRecord) bool {
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
		if json.Unmarshal(t.result, &out) != nil || out.World.Plan == "" {
			return e, how, err
		}
		return r.scanFor(scan, "world.plan (fallback: the result bytes are not a committed output)", func(e wireEntry, rec *invRecord) bool {
			return rec != nil && rec.Plan == out.World.Plan
		})
	}
	// A content address: dispatch on what it names.
	obj, ok, err := r.object(t.text)
	if err != nil {
		return wireEntry{}, "object", err
	}
	if !ok {
		return r.byWorld(t.text, "world ref")
	}
	ref := t.text
	switch obj.SemanticID {
	case coordinator.RecordV1, coordinator.RecordV2:
		return r.scanFor(scan, "record", func(e wireEntry, _ *invRecord) bool { return e.TransitionRef == ref })
	case coordinator.OutputV1:
		return r.scanFor(scan, "output", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Output == ref })
	case coordinator.InputV1:
		return r.scanFor(scan, "input", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Input == ref })
	case coordinator.EffectPlanV1:
		return r.scanFor(scan, "plan", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Plan == ref })
	case broker.EffectRecordV1:
		return r.scanFor(scan, "effect record", func(_ wireEntry, rec *invRecord) bool { return rec != nil && contains(rec.Effects, ref) })
	case broker.EffectResultV1, broker.EffectRequestV1:
		return r.scanFor(scan, "effect "+strings.TrimSuffix(strings.TrimPrefix(obj.SemanticID, "world/effect-"), "/v1"),
			func(_ wireEntry, rec *invRecord) bool {
				if rec == nil {
					return false
				}
				for _, er := range rec.Effects {
					if d, ok := r.effectRecord(er); ok && (d.ResultRef.String() == ref || d.RequestRef.String() == ref) {
						return true
					}
				}
				return false
			})
	default:
		return r.scanFor(scan, "transitionRef ("+obj.SemanticID+")", func(e wireEntry, _ *invRecord) bool { return e.TransitionRef == ref })
	}
}

func (r worldReader) effectRecord(ref string) (broker.EffectRecord, bool) {
	obj, ok, err := r.object(ref)
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
		return broker.EffectRecord{}, false
	}
	d, err := broker.DecodeRecord(obj.Payload)
	return d, err == nil
}

// byWorld resolves a world ref to the entry its log head names.
func (r worldReader) byWorld(ref, how string) (wireEntry, string, error) {
	w, ok, err := r.world(ref)
	if err != nil {
		return wireEntry{}, how, err
	}
	if !ok {
		return wireEntry{}, how, fmt.Errorf("%s is neither an object nor a world in this store", ref)
	}
	e, ok, err := r.entry(w.Revision)
	if err != nil {
		return e, how, err
	}
	if !ok || e.EntryHash != w.LogHead {
		return e, how, fmt.Errorf("world %s names log head %s, which is not entry %d", ref, w.LogHead, w.Revision)
	}
	return e, how, nil
}

// scanFor scans back from head for the first entry match accepts.
func (r worldReader) scanFor(scan int, how string, match func(wireEntry, *invRecord) bool) (wireEntry, string, error) {
	h, ok, err := r.headIndex()
	if err != nil {
		return wireEntry{}, how, err
	}
	if !ok {
		return wireEntry{}, how, errScanExhausted
	}
	var found *wireEntry
	_, err = r.scanEntries(0, h, true, scan, func(e wireEntry, rec *invRecord) (bool, error) {
		if match(e, rec) {
			found = &e
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return wireEntry{}, how, err
	}
	if found == nil {
		return wireEntry{}, how, errScanExhausted
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

type whyLink struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Ref    string `json:"ref,omitempty"`
	Detail string `json:"detail"`
}

type whyChain struct {
	Entry       int64     `json:"entry"`
	MatchedBy   string    `json:"matchedBy"`
	Coordinator bool      `json:"coordinator"`
	Note        string    `json:"note,omitempty"`
	Links       []whyLink `json:"links"`
	OK          bool      `json:"ok"`
}

func (ch *whyChain) add(name string, ok bool, ref, detail string, args ...any) {
	ch.Links = append(ch.Links, whyLink{Name: name, OK: ok, Ref: ref, Detail: fmt.Sprintf(detail, args...)})
	if !ok {
		ch.OK = false
	}
}

// entryWire / worldWire mirror host/coordinator/plan.go's canonical encodings
// (field order is the encoding order), so the entry hash and the world ref
// are recomputed, not trusted.
type entryWire struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
	TransitionRef  string `json:"transitionRef"`
}

type worldWire struct {
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

func sumJSON(v any) string {
	b, _ := json.Marshal(v)
	return hashref.SumSHA256(b).String()
}

func sum(b []byte) string { return hashref.SumSHA256(b).String() }

// checkedObject fetches ref and reports whether it resolves, carries the
// wanted semantic id, and its payload hashes to ref.
func (r worldReader) checkedObject(ref, semantic string) (wireObject, bool, string) {
	obj, ok, err := r.object(ref)
	switch {
	case err != nil:
		return obj, false, err.Error()
	case !ok:
		return obj, false, "does not resolve"
	case semantic != "" && obj.SemanticID != semantic:
		return obj, false, fmt.Sprintf("semanticId %s, want %s", obj.SemanticID, semantic)
	case sum(obj.Payload) != ref:
		return obj, false, fmt.Sprintf("payload hashes to %s, not its ref", sum(obj.Payload))
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
func (r worldReader) buildChain(e wireEntry, result []byte) (whyChain, error) {
	ch := whyChain{Entry: e.Header.EntryIndex, OK: true}
	h, err := r.health()
	if err != nil {
		return ch, err
	}
	pin := "not the daemon's pinned interpreter"
	if e.Header.Interpreter == h.InterpreterRef && h.InterpreterRef != "" {
		pin = strings.TrimSpace(strings.SplitN(h.InterpreterVersion, "\n", 2)[0]) + ", the daemon's pin"
	}
	rec, recObj, err := r.record(e)
	if err != nil {
		return ch, err
	}
	entryLine := fmt.Sprintf("#%d writtenBy %s prev %s transitionFn %s interpreter %s (%s)",
		e.Header.EntryIndex, e.Header.WrittenBy, e.Header.PrevEntryHash, e.Header.TransitionFn, e.Header.Interpreter, pin)
	if rec == nil {
		ch.Note = fmt.Sprintf("entry %d is not a coordinator invocation (writtenBy %s; its transitionRef is a %s): "+
			"there is no record chain to walk, so only the entry and its object are shown", e.Header.EntryIndex, e.Header.WrittenBy, recObj.SemanticID)
		ch.add("entry", true, e.EntryHash, "%s (a client-supplied entry hash: not recomputed)", entryLine)
		ok := sum(recObj.Payload) == e.TransitionRef
		ch.add("object", ok, e.TransitionRef, "%s, %d bytes%s", recObj.SemanticID, len(recObj.Payload), failNote(ok, "payload does not hash to transitionRef"))
		return ch, nil
	}
	ch.Coordinator = true

	// world: recomputed from the entry and the output, confirmed by the route.
	wref := sumJSON(worldWire{Revision: e.Header.EntryIndex, StateRoot: rec.Output, LogHead: e.EntryHash})
	w, wok, werr := r.world(wref)
	worldOK := werr == nil && wok && w.Revision == e.Header.EntryIndex && w.StateRoot == rec.Output && w.LogHead == e.EntryHash
	why := "GET /v1/worlds/<recomputed ref> does not return this world"
	if werr != nil {
		why = werr.Error()
	}
	ch.add("world", worldOK, wref, "revision %d, stateRoot %s, logHead %s (recomputed)%s", e.Header.EntryIndex, rec.Output, e.EntryHash, failNote(worldOK, why))

	// entry: the hash is recomputed over the canonical header + transitionRef.
	eh := sumJSON(entryWire{EntryIndex: e.Header.EntryIndex, SemanticsEpoch: e.Header.SemanticsEpoch,
		TransitionFn: e.Header.TransitionFn, Interpreter: e.Header.Interpreter, PrevEntryHash: e.Header.PrevEntryHash,
		WrittenBy: e.Header.WrittenBy, TransitionRef: e.TransitionRef})
	ch.add("entry", eh == e.EntryHash, e.EntryHash, "%s%s", entryLine, failNote(eh == e.EntryHash, "entry hash recomputes to "+eh))

	// record: its hash, and its fields agree with the entry.
	recOK := sum(rec.payload) == e.TransitionRef && rec.TransitionFn == e.Header.TransitionFn &&
		rec.Interpreter == e.Header.Interpreter && rec.SemanticsEpoch == e.Header.SemanticsEpoch
	ch.add("record", recOK, e.TransitionRef, "%s invocation %s episode %s skill %s%s", rec.SemanticID, rec.InvocationID, rec.EpisodeID, rec.SkillID,
		failNote(recOK, "the record's hash or its transitionFn/interpreter/epoch disagree with the entry"))

	// input.
	in, inOK, inWhy := r.checkedObject(rec.Input, coordinator.InputV1)
	ch.add("input", inOK, rec.Input, "%s%s", elide(string(in.Payload)), failNote(inOK, inWhy))

	// plan and effects (v2 only).
	var planned []struct {
		ID     string `json:"id"`
		Effect string `json:"effect"`
		Scope  string `json:"scope"`
		Cost   int64  `json:"cost"`
	}
	if rec.SemanticID == coordinator.RecordV2 {
		po, pOK, pWhy := r.checkedObject(rec.Plan, coordinator.EffectPlanV1)
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
				detail = "no effects; result " + elide(string(plan.Result))
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
			obj, ok, oWhy := r.checkedObject(ref, broker.EffectRecordV1)
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
			req, reqOK, _ := r.checkedObject(d.RequestRef.String(), broker.EffectRequestV1)
			okAll = okAll && reqOK
			resText := "no result"
			if !d.ResultRef.IsZero() {
				res, resOK, _ := r.checkedObject(d.ResultRef.String(), broker.EffectResultV1)
				okAll = okAll && resOK
				resText = fmt.Sprintf("result %d B", len(res.Payload))
			}
			ch.add("effect", okAll, ref, "%s %s %s@%s cost %d budget %d→%d request %d B %s%s", id, verdict, d.Effect, d.Scope, d.Cost,
				d.BudgetBefore, d.BudgetAfter, len(req.Payload), resText,
				failNote(okAll, "the record is inconsistent or its request/result does not resolve to its hash"))
		}
	}

	// output: hashes to record.output, equals world.stateRoot, and its world
	// block names this record's plan and effect records.
	out, oOK, oWhy := r.checkedObject(rec.Output, coordinator.OutputV1)
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
		if sum(result) == rec.Output {
			detail += "; your result is these exact bytes"
		} else {
			oOK, oWhy = false, "your result ("+sum(result)+") is not the committed output bytes"
		}
	}
	ch.add("output", oOK, rec.Output, "%s%s", detail, failNote(oOK, oWhy))
	return ch, nil
}

func printChain(w io.Writer, ch whyChain) {
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
	}
	if bad == 0 {
		fmt.Fprintf(w, "all %d link(s) verified\n", len(ch.Links))
	} else {
		fmt.Fprintf(w, "%d of %d link(s) BROKEN — integrity refusal (exit 3)\n", bad, len(ch.Links))
	}
}
