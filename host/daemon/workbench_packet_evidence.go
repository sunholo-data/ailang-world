package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	packetRosterInputLimit = 64 << 10
	packetRosterReadLimit  = 1 << 20
	packetRosterClaimLimit = 4 << 10
	// Bounds apply to decoded stdout and its named entries.
	packetRosterStdoutLimit    = 1 << 20
	packetRosterNamedTestLimit = 256
	packetRosterReadTimeout    = 10 * time.Second
)

// These transport projections carry no object-grade authority. The package's
// evidence.ail owns kind admission and constructor grades (executable drift gate).
// Only the configured row-105 validator can supply PROVEN.
type packetRoster struct {
	rows        []packetEvidenceRow
	unavailable string
}
type packetEvidenceRow struct {
	kind, ref, label, producer, unavailable string
	subject, compiler, version, verdict     string
	passed                                  int
}
type packetRosterReader interface {
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
	GetLogEntry(context.Context, int64) (store.LogEntry, bool, error)
	ObjectCommits(context.Context, hashref.HashRef, int64, int) ([]int64, error)
}
type packetRosterPair struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}
type packetRosterInput struct {
	Goal       string             `json:"goal"`
	Candidates []string           `json:"candidates"`
	Roster     []packetRosterPair `json:"roster"`
}

// strictRosterJSON rejects duplicate keys as well as unknown keys and trailing
// values. encoding/json's default last-key-wins behavior is unsuitable here.
func strictRosterJSON(raw []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	var value func() error
	value = func() error {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for dec.More() {
				key, err := dec.Token()
				if err != nil {
					return err
				}
				s, ok := key.(string)
				if !ok || seen[s] {
					return errors.New("duplicate JSON key")
				}
				// encoding/json accepts case-insensitive field aliases. The
				// producer's strict shape admits only the exact JSON keys.
				for _, field := range []string{"goal", "candidates", "roster", "kind", "ref", "author", "text"} {
					if s != field && strings.EqualFold(s, field) {
						return errors.New("noncanonical JSON key")
					}
				}
				seen[s] = true
				if err = value(); err != nil {
					return err
				}
			}
		case '[':
			for dec.More() {
				if err = value(); err != nil {
					return err
				}
			}
		default:
			return errors.New("unexpected JSON delimiter")
		}
		_, err = dec.Token()
		return err
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("trailing JSON")
	}
	dec = json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// packetRosterFromValidated follows only the carrying entry of an already
// validated packet. No semantic-ID lookup can discover or substitute evidence.
// Read/TESTED refusals discard the roster; proof refusals remain explicit rows.
func packetRosterFromValidated(parent context.Context, r packetRosterReader, packet coordinator.PacketResult, grader proofValidator) packetRoster {
	refuse := func(reason string) packetRoster { return packetRoster{unavailable: reason} }
	if packet.Code != 0 {
		return refuse("packet not validated")
	}
	ctx, cancel := context.WithTimeout(parent, packetRosterReadTimeout)
	defer cancel()
	read := func(ref hashref.HashRef, role, id string, limit int) (store.Object, string) {
		if ctx.Err() != nil {
			return store.Object{}, "read deadline"
		}
		o, ok, err := r.GetObject(ctx, ref)
		if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return store.Object{}, "read deadline"
		}
		if err != nil || !ok {
			return store.Object{}, role + " missing"
		}
		if len(o.Payload) > limit {
			return store.Object{}, role + " size"
		}
		if o.Hash != ref || hashref.SumSHA256(o.Payload) != ref {
			return store.Object{}, role + " hash"
		}
		wantInterface := hashref.SumSHA256([]byte(id))
		if id == evidence.ProofSemanticID {
			wantInterface = evidence.InterfaceHashV1
		}
		if o.SemanticID != id || o.InterfaceHash != wantInterface {
			return store.Object{}, role + " kind"
		}
		return o, ""
	}
	entry, ok, err := r.GetLogEntry(ctx, packet.Entry)
	if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return refuse("read deadline")
	}
	if err != nil || !ok {
		return refuse("carrying entry missing")
	}
	record, reason := read(entry.TransitionRef, "record", coordinator.RecordV2, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	var rec struct {
		Input string `json:"input"`
	}
	if json.Unmarshal(record.Payload, &rec) != nil {
		return refuse("record schema")
	}
	ref, err := hashref.Parse(rec.Input)
	if err != nil {
		return refuse("input ref")
	}
	// The smaller input cap is checked on decoded object bytes, before parsing.
	input, reason := read(ref, "input", coordinator.InputV1, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	if len(input.Payload) > packetRosterInputLimit {
		return refuse("input size")
	}
	var in packetRosterInput
	if strictRosterJSON(input.Payload, &in) != nil {
		return refuse("input shape")
	}
	if in.Goal == "" || len(in.Candidates) != 2 || in.Candidates[0] == "" || in.Candidates[1] == "" || in.Candidates[0] == in.Candidates[1] {
		return refuse("input schema")
	}
	if len(in.Roster) != 4 {
		return refuse("roster count")
	}
	seen := map[string]bool{}
	refs := map[string]hashref.HashRef{}
	for _, pair := range in.Roster {
		if seen[pair.Kind] {
			return refuse("duplicate kind")
		}
		seen[pair.Kind] = true
		// Transport transcription of world-demo/evidence.rosterKindAllowed;
		// TestWorldDemoLegacyGradeDrift runs the package's law against the kernel.
		switch pair.Kind {
		case "test-invocation", "recorded-effect", "authored-claim", "proof-envelope":
		default:
			return refuse("unknown kind")
		}
		ref, err := hashref.Parse(pair.Ref)
		if err != nil {
			return refuse("ref syntax")
		}
		refs[pair.Kind] = ref
	}
	effect, reason := read(refs["recorded-effect"], "effect", broker.EffectRecordV1, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	er, err := broker.DecodeRecord(effect.Payload)
	if err != nil {
		return refuse("effect decode")
	}
	if !broker.RecordConsistent(er) {
		return refuse("effect inconsistent")
	}
	if er.ResultRef.IsZero() {
		return refuse("effect result absent")
	}
	if _, reason = read(er.RequestRef, "request", broker.EffectRequestV1, packetRosterReadLimit); reason != "" {
		return refuse(reason)
	}
	if !er.ResultRef.IsZero() {
		if _, reason = read(er.ResultRef, "result", broker.EffectResultV1, packetRosterReadLimit); reason != "" {
			return refuse(reason)
		}
	}
	claim, reason := read(refs["authored-claim"], "claim", "world/demo/authored-claim/v1", packetRosterClaimLimit)
	if reason != "" {
		return refuse(reason)
	}
	var authored struct {
		Kind   string `json:"kind"`
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	if strictRosterJSON(claim.Payload, &authored) != nil || authored.Text == "" {
		return refuse("claim shape")
	}
	if authored.Kind != "authored-claim" {
		return refuse("claim kind")
	}
	if authored.Author != "demo-authored" || claim.Provenance != "demo-authored" {
		return refuse("claim author")
	}
	if ctx.Err() != nil {
		return refuse("read deadline")
	}
	rows := make([]packetEvidenceRow, 0, 4)
	for _, pair := range in.Roster {
		row := packetEvidenceRow{kind: pair.Kind, ref: pair.Ref}
		switch pair.Kind {
		case "recorded-effect":
			row.label = "ATTESTED"
			row.producer = effect.Provenance
		case "authored-claim":
			row.label = "CLAIMED"
			row.producer = authored.Author
		case "test-invocation":
			row, reason = packetTestedRow(ctx, r, read, refs[pair.Kind])
			if reason != "" {
				return refuse(reason)
			}
		case "proof-envelope":
			row = packetProvenRow(ctx, read, refs[pair.Kind], grader)
		}
		rows = append(rows, row)
	}
	if ctx.Err() != nil {
		return refuse("read deadline")
	}
	return packetRoster{rows: rows}
}

// The record wire is coordinator.RecordV2 (host/coordinator/plan.go).
type packetTestRecord struct {
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
type packetObjectRead func(hashref.HashRef, string, string, int) (store.Object, string)

func packetTestedRow(ctx context.Context, r packetRosterReader, read packetObjectRead, ref hashref.HashRef) (packetEvidenceRow, string) {
	row := packetEvidenceRow{kind: "test-invocation", ref: ref.String()}
	obj, reason := read(ref, "test record", coordinator.RecordV2, packetRosterReadLimit)
	if reason != "" {
		return row, reason
	}
	var rec packetTestRecord
	if strictRosterJSON(obj.Payload, &rec) != nil {
		return row, "test record schema"
	}
	commits, err := r.ObjectCommits(ctx, ref, -1, 1)
	if err != nil || len(commits) == 0 {
		return row, "test carrying commit missing"
	}
	entry, ok, err := r.GetLogEntry(ctx, commits[0])
	if err != nil || !ok || entry.TransitionRef != ref {
		return row, "test carrying commit mismatch"
	}
	if len(rec.Effects) != 1 {
		return row, "test effect count"
	}
	effectRef, err := hashref.Parse(rec.Effects[0])
	if err != nil {
		return row, "test effect ref"
	}
	obj, reason = read(effectRef, "test effect", broker.EffectRecordV1, packetRosterReadLimit)
	if reason != "" {
		return row, reason
	}
	er, err := broker.DecodeRecord(obj.Payload)
	if err != nil || !broker.RecordConsistent(er) {
		return row, "test effect inconsistent"
	}
	req, reason := read(er.RequestRef, "test request", broker.EffectRequestV1, packetRosterReadLimit)
	if reason != "" {
		return row, reason
	}
	result, reason := read(er.ResultRef, "test result", broker.EffectResultV1, packetRosterReadLimit)
	if reason != "" {
		return row, reason
	}
	// Broker requests are length-prefixed bytes, NOT JSON envelopes (requestBytes).
	prefix := fmt.Sprintf("%d:%s%d:%s%d:", len(er.Effect), er.Effect, len(er.Scope), er.Scope, er.Cost)
	if er.Effect != broker.EffectAilangCLI || !bytes.HasPrefix(req.Payload, []byte(prefix)) {
		return row, "test effect not Ailang.CLI"
	}
	now, payload, found := strings.Cut(string(req.Payload[len(prefix):]), ":")
	_, timeErr := strconv.ParseInt(now, 10, 64)
	var request map[string]json.RawMessage
	if !found || timeErr != nil || strictRosterJSON([]byte(payload), &request) != nil {
		return row, "test request schema"
	}
	var op string
	_ = json.Unmarshal(request["op"], &op)
	if op != "test" {
		return row, "test op"
	}
	// Results also carry host stamps; decode their known outcome fields separately.
	var out map[string]json.RawMessage
	if strictRosterJSON(result.Payload, &out) != nil || out == nil {
		return row, "test result schema"
	}
	var argv []string
	if json.Unmarshal(out["argv"], &argv) != nil || len(argv) == 0 || argv[0] != "test" {
		return row, "test argv"
	}
	var exit *int
	var good bool
	if json.Unmarshal(out["exit_code"], &exit) != nil || exit == nil || *exit != 0 {
		return row, "test exit code"
	}
	if json.Unmarshal(out["ok"], &good) != nil || !good {
		return row, "test not ok"
	}
	var stdout string
	if json.Unmarshal(out["stdout"], &stdout) != nil {
		return row, "test stdout schema"
	}
	if len(stdout) > packetRosterStdoutLimit {
		return row, "test stdout size"
	}
	count, reason := packetTestReport([]byte(stdout))
	if reason != "" {
		return row, reason
	}
	row.label, row.producer, row.verdict, row.passed = "TESTED", result.Provenance, "pass", count
	return row, ""
}

// These are the fields emitted by the real test --json report, including
// property diagnostics. Decode strictly before projecting the required verdict.
type packetNamedTest struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Duration  string `json:"duration"`
	Location  string `json:"location"`
	Error     string `json:"error,omitempty"`
	Seed      string `json:"seed,omitempty"`
	SkipKind  string `json:"skip_kind,omitempty"`
	Discarded int    `json:"discarded_inputs,omitempty"`
	Generated int    `json:"generated_inputs,omitempty"`
	Runs      int    `json:"tests_run,omitempty"`
}
type packetReportWire struct {
	Success    bool              `json:"success"`
	Failed     int               `json:"failed_tests"`
	Passed     int               `json:"passed_tests"`
	Skipped    int               `json:"skipped_tests"`
	Vacuous    int               `json:"vacuous_skips"`
	Total      int               `json:"total_tests"`
	Module     string            `json:"module_path"`
	Duration   string            `json:"total_duration"`
	Seed       string            `json:"seed"`
	Derivation string            `json:"seed_derivation"`
	Mode       string            `json:"seed_mode"`
	Tests      []packetNamedTest `json:"tests"`
	Properties []packetNamedTest `json:"properties"`
}

func packetTestReport(raw []byte) (int, string) {
	var report map[string]json.RawMessage
	var wire packetReportWire
	if !utf8.Valid(raw) || strictRosterJSON(raw, &wire) != nil || strictRosterJSON(raw, &report) != nil {
		return 0, "test report JSON"
	}
	if report == nil {
		return 0, "test report null"
	}
	var success bool
	if json.Unmarshal(report["success"], &success) != nil || !success {
		return 0, "test success"
	}
	counts := make([]int, 4)
	for i, key := range []string{"failed_tests", "skipped_tests", "vacuous_skips", "passed_tests"} {
		var n *int
		if json.Unmarshal(report[key], &n) != nil || n == nil {
			return 0, "test report counters"
		}
		counts[i] = *n
	}
	if counts[0] != 0 {
		return 0, "test failed"
	}
	if counts[1] != 0 {
		return 0, "test skipped"
	}
	if counts[2] != 0 {
		return 0, "test vacuous"
	}
	if counts[3] < 1 {
		return 0, "test zero passes"
	}
	var named []map[string]json.RawMessage
	for _, key := range []string{"tests", "properties"} {
		var entries []map[string]json.RawMessage
		_ = json.Unmarshal(report[key], &entries) // strict wire decode already checked types
		named = append(named, entries...)
	}
	if len(named) == 0 || len(named) > packetRosterNamedTestLimit {
		return 0, "test named count"
	}
	for _, entry := range named {
		var name, status string
		_ = json.Unmarshal(entry["name"], &name)
		_ = json.Unmarshal(entry["status"], &status)
		if name == "" {
			return 0, "test empty name"
		}
		if status != "pass" {
			return 0, "test status"
		}
	}
	if counts[3] != len(named) {
		return 0, "test passed count"
	}
	return len(named), ""
}

func packetProvenRow(ctx context.Context, read packetObjectRead, ref hashref.HashRef, grader proofValidator) packetEvidenceRow {
	row := packetEvidenceRow{kind: "proof-envelope", ref: ref.String(), producer: "unavailable"}
	refuse := func(reason string) packetEvidenceRow { row.unavailable = reason; return row }
	obj, reason := read(ref, "proof envelope", evidence.ProofSemanticID, evidence.MaxBytes)
	if reason != "" {
		return refuse(reason)
	}
	row.producer = obj.Provenance
	envelope, err := evidence.DecodeAuthenticatedEnvelope(obj.Payload)
	if err != nil {
		return refuse(gradeReasonStructuralCandidate)
	}
	report, _ := evidence.DecodeProofReportV1(envelope.Report) // envelope decoder already checks it
	if _, reason = read(report.Subject, "proof subject", "world/source/ailang/v1", packetRosterReadLimit); reason != "" {
		return refuse(reason)
	}
	if grader == nil {
		return refuse(objectGradeUnavailableReason)
	}
	result := grader.ValidateProof(ctx, ref, report.Subject)
	sealed, valid := result.Validated()
	if !valid {
		reason, _ := result.Unsupported()
		return refuse(gradeReasonSentence(reason))
	}
	if _, proven := grader.Resolve(sealed).Proven(); !proven {
		return refuse(gradeReasonNoneValidated)
	}
	// The subject is shown only once the MAC and Resolve have authenticated it.
	row.label, row.subject, row.compiler, row.version = "PROVEN", report.Subject.String(), report.Compiler.String(), report.CompilerVersion
	return row
}
