package main

// Row 182 (w-episode-outcome-capture, ratified attended 2026-10-16): the
// `outcome record` verb. An operator harness records how one graded task
// turned out — and how it was judged — as a world/task-outcome/v1 object
// committed next to the episode's effects, so a provenance walk can reach the
// outcome, the grading rule and the task's entry range. Row 114 Q1 measured
// the gap: the store held every tool call but no task, outcome or rule.
//
// The verb adds no route and no authority: the object is committed through
// the shared one-object writer (commitOneObject) under the episode's own
// session credential, written by "harness:floor".

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	taskOutcomeSemanticID = "world/task-outcome/v1"
	taskOutcomeProvenance = "harness:floor"
	taskOutcomeWrittenBy  = "harness:floor"
	maxTaskOutcomeBytes   = 16 << 10
	maxOutcomeHeadBytes   = 400
	maxOutcomeFieldBytes  = 1024
	maxOutcomeArgv        = 32
	maxOutcomeInputFiles  = 32

	exitOutcomeRefused   = 2 // the daemon refused the commit
	exitOutcomeUncertain = 3 // committed or not: reconcile, never resend blindly
)

// The payload. Fields are declared in json-name order so encoding/json emits
// sorted keys: the canonical form is json.Marshal of these structs.
type taskOutcome struct {
	Agent     string          `json:"agent"`
	Arm       string          `json:"arm"`
	Entries   *outcomeEntries `json:"entries"`
	EpisodeID string          `json:"episodeId"`
	Grading   outcomeGrading  `json:"grading"`
	Label     string          `json:"label"`
	Outcome   outcomeResult   `json:"outcome"`
	Run       int64           `json:"run"`
	Task      string          `json:"task"`
}

type outcomeEntries struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type outcomeGrading struct {
	Argv           []string `json:"argv"`
	ExpectedStdout string   `json:"expectedStdout"`
	GraderVersion  string   `json:"graderVersion"`
	InputFiles     []string `json:"inputFiles"`
	Rule           string   `json:"rule"`
	ToolSha256     string   `json:"toolSha256"`
}

type outcomeResult struct {
	Class          string `json:"class"`
	CompileOk      bool   `json:"compileOk"`
	ErrorCategory  string `json:"errorCategory"`
	RuntimeOk      bool   `json:"runtimeOk"`
	SolutionSha256 string `json:"solutionSha256"`
	SolutionState  string `json:"solutionState"`
	StderrHead     string `json:"stderrHead"`
	StdoutHead     string `json:"stdoutHead"`
	StdoutOk       bool   `json:"stdoutOk"`
}

// parseTaskOutcome strict-decodes and validates one payload and returns its
// canonical bytes. Every refusal happens here, before anything is sent.
func parseTaskOutcome(data []byte) ([]byte, taskOutcome, error) {
	var o taskOutcome
	if len(data) > maxTaskOutcomeBytes {
		return nil, o, fmt.Errorf("outcome exceeds %d bytes", maxTaskOutcomeBytes)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		return nil, o, fmt.Errorf("invalid outcome JSON: %v", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, o, fmt.Errorf("invalid outcome JSON: trailing data")
	}
	// Presence: decoding into structs cannot tell an absent key from a zero
	// value, so the required keys are checked on a generic decode.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, o, fmt.Errorf("invalid outcome JSON: %v", err)
	}
	for _, k := range []string{"agent", "arm", "entries", "episodeId", "grading", "label", "outcome", "run", "task"} {
		if _, ok := raw[k]; !ok {
			return nil, o, fmt.Errorf("outcome: missing required field %q", k)
		}
	}
	for name, v := range map[string]string{"agent": o.Agent, "episodeId": o.EpisodeID, "task": o.Task, "label": o.Label,
		"outcome.class": o.Outcome.Class, "grading.rule": o.Grading.Rule, "grading.graderVersion": o.Grading.GraderVersion} {
		if v == "" {
			return nil, o, fmt.Errorf("outcome: %s is empty", name)
		}
	}
	switch o.Arm {
	case "shell":
		if o.Entries != nil {
			return nil, o, fmt.Errorf("outcome: a shell row has no World entries (entries must be null)")
		}
	case "world":
		if o.Entries == nil {
			return nil, o, fmt.Errorf("outcome: a world row needs entries {from, to}")
		}
		if o.Entries.From < 1 || o.Entries.To < o.Entries.From-1 {
			// to == from-1 is a world task that committed no entry.
			return nil, o, fmt.Errorf("outcome: bad entries range %d..%d", o.Entries.From, o.Entries.To)
		}
	default:
		return nil, o, fmt.Errorf("outcome: arm %q is not shell or world", o.Arm)
	}
	if o.Run < 1 {
		return nil, o, fmt.Errorf("outcome: run %d < 1", o.Run)
	}
	if len(o.Outcome.StdoutHead) > maxOutcomeHeadBytes || len(o.Outcome.StderrHead) > maxOutcomeHeadBytes {
		return nil, o, fmt.Errorf("outcome: stdout/stderr head exceeds %d bytes", maxOutcomeHeadBytes)
	}
	if len(o.Grading.Argv) == 0 || len(o.Grading.Argv) > maxOutcomeArgv || len(o.Grading.InputFiles) > maxOutcomeInputFiles {
		return nil, o, fmt.Errorf("outcome: grading.argv needs 1..%d items, inputFiles at most %d", maxOutcomeArgv, maxOutcomeInputFiles)
	}
	for _, s := range append(append([]string{o.Agent, o.EpisodeID, o.Task, o.Label, o.Outcome.Class, o.Outcome.ErrorCategory,
		o.Outcome.SolutionSha256, o.Outcome.SolutionState, o.Grading.Rule, o.Grading.GraderVersion, o.Grading.ToolSha256,
		o.Grading.ExpectedStdout}, o.Grading.Argv...), o.Grading.InputFiles...) {
		if len(s) > maxOutcomeFieldBytes {
			return nil, o, fmt.Errorf("outcome: a field exceeds %d bytes", maxOutcomeFieldBytes)
		}
	}
	canonical, err := json.Marshal(o)
	if err != nil {
		return nil, o, fmt.Errorf("outcome: encode: %v", err)
	}
	return canonical, o, nil
}

// taskOutcomeObject is the committed object for canonical payload bytes.
func taskOutcomeObject(canonical []byte) demoObject {
	return demoObject{
		Hash:          hashref.SumSHA256(canonical).String(),
		InterfaceHash: hashref.SumSHA256([]byte(taskOutcomeSemanticID)).String(),
		SemanticID:    taskOutcomeSemanticID,
		Provenance:    taskOutcomeProvenance,
		Payload:       canonical,
	}
}

// outcomeCommitter is the write seam; production uses commitOneObject.
type outcomeCommitter func(ctx context.Context, c *client, token string, obj demoObject, writtenBy, interpreter, op string) (int64, error)

func runOutcome(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runOutcomeWith(addr, args, stdin, stdout, stderr, commitOneObject)
}

func runOutcomeWith(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer, commit outcomeCommitter) int {
	if len(args) == 0 || args[0] != "record" {
		fmt.Fprintln(stderr, "usage: ailang-worldd [--addr <url>] outcome record --file <outcome.json|-> [--session <file|token>]")
		return exitUsage
	}
	fs := flag.NewFlagSet("ailang-worldd outcome record", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "", "outcome JSON file, or - for stdin")
	session := fs.String("session", "", "session credential: a file holding the 64-hex token (preferred) or the token itself; default $"+sessionEnvVar)
	if err := fs.Parse(args[1:]); err != nil || len(fs.Args()) != 0 || *file == "" {
		if *file == "" {
			fmt.Fprintln(stderr, "ailang-worldd outcome record: --file is required")
		}
		return exitUsage
	}
	var src io.Reader = stdin
	if *file != "-" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd outcome record: %v\n", err)
			return exitUsage
		}
		defer func() { _ = f.Close() }()
		src = f
	}
	data, err := io.ReadAll(io.LimitReader(src, maxTaskOutcomeBytes+1))
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd outcome record: read: %v\n", err)
		return exitUsage
	}
	canonical, o, err := parseTaskOutcome(data)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd outcome record: refused: %v\n", err)
		return exitUsage
	}
	token, err := resolveSession(*session, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd outcome record: %v\n", err)
		return exitUsage
	}
	obj := taskOutcomeObject(canonical)
	ctx, cancel := budgetContext(demoProofObjectBound)
	defer cancel()
	entry, err := commit(ctx, demoProofClient(addr), token, obj, taskOutcomeWrittenBy, "", "outcome")
	if err != nil {
		var uncertain *store.UncertainError
		if errors.As(err, &uncertain) {
			fmt.Fprintf(stderr, "ailang-worldd outcome record: UNCERTAIN (reconcile with object get %s before any resend): %v\n", obj.Hash, err)
			return exitOutcomeUncertain
		}
		fmt.Fprintf(stderr, "ailang-worldd outcome record: %v\n", err)
		return exitOutcomeRefused
	}
	fmt.Fprintf(stdout, "recorded %s %s/%s/%s run %d: object %s entry %d\n", taskOutcomeSemanticID, o.Agent, o.Arm, o.Task, o.Run, obj.Hash, entry)
	return exitOK
}
