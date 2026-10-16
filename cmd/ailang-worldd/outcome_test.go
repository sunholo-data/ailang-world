package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// validOutcome is a world row shaped like the row-114 Q1 case (codex r2).
func validOutcome() map[string]any {
	return map[string]any{
		"agent": "codex", "arm": "world", "episodeId": "fl-cx-r2", "run": 2, "label": "final-paired-r2",
		"task": "config_file_parser", "entries": map[string]any{"from": 259, "to": 261},
		"outcome": map[string]any{"class": "FAIL", "errorCategory": "logic_error", "compileOk": true, "runtimeOk": true,
			"stdoutOk": false, "stdoutHead": "configuration file not found\n", "stderrHead": "",
			"solutionSha256": strings.Repeat("a", 64), "solutionState": "present"},
		"grading": map[string]any{"rule": "only benchmark/solution.ail is copied into an empty temporary directory, plus the task's declared input_files; then run",
			"argv":       []any{"run", "--entry", "main", "--quiet", "--relax-modules", "--caps", "IO,FS", "benchmark/solution.ail"},
			"inputFiles": []any{}, "expectedStdout": "Loaded MyApp v1.0.0 on port 8080 with 3 features\n",
			"graderVersion": "floor-grade/1", "toolSha256": strings.Repeat("b", 64)},
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestOutcomeParseRefusals: every AC2 refusal happens before anything is sent.
func TestOutcomeParseRefusals(t *testing.T) {
	if _, _, err := parseTaskOutcome(mustJSON(t, validOutcome())); err != nil {
		t.Fatalf("control: the valid outcome is refused: %v", err)
	}
	type mut func(m map[string]any)
	cases := map[string]mut{
		"unknown top field":     func(m map[string]any) { m["extra"] = 1 },
		"unknown nested field":  func(m map[string]any) { m["outcome"].(map[string]any)["why"] = "x" },
		"bad arm":               func(m map[string]any) { m["arm"] = "both" },
		"world without entries": func(m map[string]any) { m["entries"] = nil },
		"shell with entries":    func(m map[string]any) { m["arm"] = "shell" },
		"from > to+1":           func(m map[string]any) { m["entries"] = map[string]any{"from": 262, "to": 259} },
		"from 0":                func(m map[string]any) { m["entries"] = map[string]any{"from": 0, "to": 3} },
		"run 0":                 func(m map[string]any) { m["run"] = 0 },
		"non-integer run":       func(m map[string]any) { m["run"] = 2.5 },
		"empty class":           func(m map[string]any) { m["outcome"].(map[string]any)["class"] = "" },
		"empty rule":            func(m map[string]any) { m["grading"].(map[string]any)["rule"] = "" },
		"empty argv":            func(m map[string]any) { m["grading"].(map[string]any)["argv"] = []any{} },
		"stdout head over 400":  func(m map[string]any) { m["outcome"].(map[string]any)["stdoutHead"] = strings.Repeat("x", 401) },
		"field over 1 KiB":      func(m map[string]any) { m["label"] = strings.Repeat("x", 1025) },
	}
	for _, k := range []string{"agent", "arm", "entries", "episodeId", "grading", "label", "outcome", "run", "task"} {
		k := k
		cases["missing "+k] = func(m map[string]any) { delete(m, k) }
	}
	for name, f := range cases {
		m := validOutcome()
		f(m)
		if _, _, err := parseTaskOutcome(mustJSON(t, m)); err == nil {
			t.Errorf("%s: accepted, want a refusal", name)
		}
	}
	// A shell row with entries:null is valid (the control for "shell with entries").
	m := validOutcome()
	m["arm"], m["entries"] = "shell", nil
	if _, _, err := parseTaskOutcome(mustJSON(t, m)); err != nil {
		t.Errorf("shell row with null entries refused: %v", err)
	}
	// A world task that committed nothing: to == from-1.
	m = validOutcome()
	m["entries"] = map[string]any{"from": 259, "to": 258}
	if _, _, err := parseTaskOutcome(mustJSON(t, m)); err != nil {
		t.Errorf("empty world range refused: %v", err)
	}
	if _, _, err := parseTaskOutcome(append(mustJSON(t, validOutcome()), []byte(" {}")...)); err == nil {
		t.Error("trailing JSON value accepted")
	}
	if _, _, err := parseTaskOutcome(bytes.Repeat([]byte(" "), maxTaskOutcomeBytes+1)); err == nil {
		t.Error("oversized payload accepted")
	}
}

// TestOutcomeCanonical: key order and whitespace do not change the committed
// bytes, and the canonical form is a fixed point.
func TestOutcomeCanonical(t *testing.T) {
	a, _, err := parseTaskOutcome(mustJSON(t, validOutcome()))
	if err != nil {
		t.Fatal(err)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, mustJSON(t, validOutcome()), "", "   "); err != nil {
		t.Fatal(err)
	}
	b, _, err := parseTaskOutcome(pretty.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("whitespace changed the canonical bytes:\n%s\n%s", a, b)
	}
	c, _, err := parseTaskOutcome(a)
	if err != nil || !bytes.Equal(a, c) {
		t.Fatalf("canonical form is not a fixed point: %v", err)
	}
	// Keys are sorted: re-encoding through a generic map (Go sorts map keys) is identical.
	var generic any
	if err := json.Unmarshal(a, &generic); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, mustJSON(t, generic)) {
		t.Fatalf("canonical bytes are not key-sorted:\n%s", a)
	}
	o := taskOutcomeObject(a)
	if o.Hash != hashref.SumSHA256(a).String() || o.SemanticID != taskOutcomeSemanticID || o.Provenance != "harness:floor" {
		t.Fatalf("object: %+v", o)
	}
}

func writeOutcomeFile(t *testing.T, v any) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "outcome.json")
	if err := os.WriteFile(p, mustJSON(t, v), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// TestOutcomeRecordRoundTrip (AC1, AC6): the verb commits one object over the
// real daemon; it reads back byte-equal, is found by its semantic id, is the
// transitionRef of an entry written by harness:floor, and why/provenance still
// work over a log whose head is that entry.
func TestOutcomeRecordRoundTrip(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	r := worldReader{ctx: ctx, c: rig.client(forwardRT{})}
	headBefore, _, err := headState(ctx, r.c)
	if err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	file := writeOutcomeFile(t, validOutcome())
	if rc := runOutcome(rig.url, []string{"record", "--file", file, "--session", rig.token}, nil, &out, &errb); rc != exitOK {
		t.Fatalf("outcome record rc=%d stderr=%s", rc, errb.String())
	}
	canonical, _, _ := parseTaskOutcome(mustJSON(t, validOutcome()))
	ref := hashref.SumSHA256(canonical).String()
	if !strings.Contains(out.String(), "object "+ref+" entry 1") {
		t.Fatalf("stdout does not name the object and entry 1: %q", out.String())
	}
	obj, ok, err := r.object(ref)
	if err != nil || !ok {
		t.Fatalf("object %s: ok=%v err=%v", ref, ok, err)
	}
	if obj.SemanticID != taskOutcomeSemanticID || !bytes.Equal(obj.Payload, canonical) {
		t.Fatalf("read back %s %q, want %s %q", obj.SemanticID, obj.Payload, taskOutcomeSemanticID, canonical)
	}
	e, ok, err := r.entry(1)
	if err != nil || !ok {
		t.Fatalf("entry 1: ok=%v err=%v", ok, err)
	}
	if e.Header.WrittenBy != "harness:floor" || e.TransitionRef != ref {
		t.Fatalf("entry 1 writtenBy=%q transitionRef=%s, want harness:floor %s", e.Header.WrittenBy, e.TransitionRef, ref)
	}
	headAfter, _, _ := headState(ctx, r.c)
	if headAfter == headBefore {
		t.Fatal("head did not move")
	}
	out.Reset()
	errb.Reset()
	if rc := runObject(rig.url, []string{"find", taskOutcomeSemanticID}, &out, &errb); rc != exitOK || !strings.Contains(out.String(), ref) {
		t.Fatalf("object find rc=%d out=%q err=%q", rc, out.String(), errb.String())
	}
	out.Reset()
	errb.Reset()
	if rc := runWhy(rig.url, []string{"head"}, nil, &out, &errb); rc != exitOK {
		t.Fatalf("why head over an outcome entry: rc=%d out=%q err=%q", rc, out.String(), errb.String())
	}
	out.Reset()
	errb.Reset()
	if rc := runWhy(rig.url, []string{ref}, nil, &out, &errb); rc != exitOK {
		t.Fatalf("why <outcome ref>: rc=%d out=%q err=%q", rc, out.String(), errb.String())
	}
}

// TestOutcomeRecordRefusals (AC2, AC3): a refused payload sends nothing; a
// missing session is the daemon's refusal (exit 2) and commits nothing.
func TestOutcomeRecordRefusals(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	c := rig.client(forwardRT{})
	head0, _, _ := headState(ctx, c)

	bad := validOutcome()
	bad["arm"] = "both"
	var out, errb bytes.Buffer
	if rc := runOutcome(rig.url, []string{"record", "--file", writeOutcomeFile(t, bad), "--session", rig.token}, nil, &out, &errb); rc != exitUsage {
		t.Fatalf("bad payload rc=%d, want %d", rc, exitUsage)
	}
	wrong := strings.Repeat("0", 64)
	errb.Reset()
	if rc := runOutcome(rig.url, []string{"record", "--file", writeOutcomeFile(t, validOutcome()), "--session", wrong}, nil, &out, &errb); rc != exitOutcomeRefused {
		t.Fatalf("unknown session rc=%d, want %d; stderr=%s", rc, exitOutcomeRefused, errb.String())
	}
	if !strings.Contains(errb.String(), "401") {
		t.Fatalf("refusal does not surface the daemon's 401: %s", errb.String())
	}
	if head1, _, _ := headState(ctx, c); head1 != head0 {
		t.Fatalf("a refused record moved the head %s -> %s", head0, head1)
	}
}

// TestOutcomeUncertainExit (AC3): an uncertain commit is exit 3 and names the
// object to reconcile; a plain refusal is exit 2. The seam replaces the writer.
func TestOutcomeUncertainExit(t *testing.T) {
	file := writeOutcomeFile(t, validOutcome())
	uncertain := func(context.Context, *client, string, demoObject, string, string, string) (int64, error) {
		return 0, &store.UncertainError{Op: "outcome commit", Cause: errors.New("connection reset")}
	}
	var out, errb bytes.Buffer
	if rc := runOutcomeWith("http://127.0.0.1:1", []string{"record", "--file", file, "--session", strings.Repeat("1", 64)}, nil, &out, &errb, uncertain); rc != exitOutcomeUncertain {
		t.Fatalf("uncertain rc=%d, want %d", rc, exitOutcomeUncertain)
	}
	if !strings.Contains(errb.String(), "UNCERTAIN") || !strings.Contains(errb.String(), "object get sha256:") {
		t.Fatalf("uncertain stderr: %s", errb.String())
	}
	var gotWrittenBy, gotInterp string
	refused := func(_ context.Context, _ *client, _ string, _ demoObject, writtenBy, interp, _ string) (int64, error) {
		gotWrittenBy, gotInterp = writtenBy, interp
		return 0, errors.New("outcome writer: commit refused HTTP 409 Conflict; not retried")
	}
	errb.Reset()
	if rc := runOutcomeWith("http://127.0.0.1:1", []string{"record", "--file", file, "--session", strings.Repeat("1", 64)}, nil, &out, &errb, refused); rc != exitOutcomeRefused {
		t.Fatalf("refused rc=%d, want %d", rc, exitOutcomeRefused)
	}
	if gotWrittenBy != "harness:floor" || gotInterp != "" {
		t.Fatalf("writer called with writtenBy=%q interpreter=%q, want harness:floor and the head's own", gotWrittenBy, gotInterp)
	}
}

// TestOutcomeStdin: --file - reads the payload from stdin.
func TestOutcomeStdin(t *testing.T) {
	var sent demoObject
	seam := func(_ context.Context, _ *client, _ string, obj demoObject, _, _, _ string) (int64, error) {
		sent = obj
		return 7, nil
	}
	var out, errb bytes.Buffer
	if rc := runOutcomeWith("http://127.0.0.1:1", []string{"record", "--file", "-", "--session", strings.Repeat("1", 64)},
		bytes.NewReader(mustJSON(t, validOutcome())), &out, &errb, seam); rc != exitOK {
		t.Fatalf("stdin rc=%d stderr=%s", rc, errb.String())
	}
	canonical, _, _ := parseTaskOutcome(mustJSON(t, validOutcome()))
	if !bytes.Equal(sent.Payload, canonical) || !strings.Contains(out.String(), "entry 7") {
		t.Fatalf("sent %q out %q", sent.Payload, out.String())
	}
}
