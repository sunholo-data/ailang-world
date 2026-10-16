package main

// Demo orchestration is a host adapter: all public writes use the sole-writer
// daemon REST API. The only database handle below is an existing read-only
// authority view, after daemon/database identity has been checked.
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type demoObject struct {
	Hash          string `json:"hash"`
	InterfaceHash string `json:"interfaceHash"`
	SemanticID    string `json:"semanticId"`
	Provenance    string `json:"provenance"`
	Payload       []byte `json:"payload"`
}
type demoCommit struct {
	ObservedHead string       `json:"observedHead"`
	Objects      []demoObject `json:"objects"`
	NextWorld    wireWorld    `json:"nextWorld"`
	Entry        wireEntry    `json:"entry"`
}

func demoGet(r worldReader, path string, v any) error {
	ok, err := r.getJSON(path, v)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("required resource absent: %s", path)
	}
	return nil
}
func demoPreflight(ctx context.Context, o demoOptions, token string, r worldReader) (daemon.HealthResponse, int64, error) {
	var h daemon.HealthResponse
	if _, err := os.Lstat(o.manifest); !os.IsNotExist(err) {
		return h, 0, fmt.Errorf("manifest must be a new absent path")
	}
	if err := demoGet(r, "/v1/health", &h); err != nil {
		return h, 0, err
	}
	want, e := filepath.Abs(o.db)
	if e != nil {
		return h, 0, e
	}
	got, e := filepath.Abs(h.DBPath)
	if e != nil {
		return h, 0, e
	}
	if want != got {
		return h, 0, fmt.Errorf("db/health identity mismatch")
	}
	pin, err := os.ReadFile(o.pin)
	if err != nil {
		return h, 0, fmt.Errorf("cannot read interpreter pin")
	}
	pinRef := hashref.SumSHA256(pin)
	if pinRef.String() != h.InterpreterRef {
		return h, 0, fmt.Errorf("interpreter pin mismatch")
	}

	hasWorld, err := demoGrants(ctx, want, token)
	if err != nil {
		return h, 0, err
	}
	// Startup already obtained the bounded, credential-scrubbed --version
	// output. Cross-check its existing read-only sidecar against the matching
	// local bytes and health; the seed launches no additional subprocess.
	manifest, err := archive.New(want).ReadManifest(pinRef)
	if err != nil {
		return h, 0, fmt.Errorf("interpreter archive manifest unavailable")
	}
	if manifest.Hash != pinRef.String() || manifest.Size != int64(len(pin)) {
		return h, 0, fmt.Errorf("interpreter archive identity mismatch")
	}
	if manifest.Version != h.InterpreterVersion || !strings.HasPrefix(manifest.Version, "AILANG v0.41.0\n") {
		return h, 0, fmt.Errorf("interpreter version mismatch")
	}

	ref, hasHead, err := headState(ctx, r.c)
	if err != nil {
		return h, 0, err
	}
	page, err := r.page(0, 1)
	if err != nil {
		return h, 0, err
	}
	if hasWorld || hasHead || ref != "" || len(page) != 0 {
		return h, 0, fmt.Errorf("nonempty public world/log; inspect history before seeding")
	}
	var rh struct {
		Head string `json:"head"`
	}
	if err = demoGet(r, "/v1/registry/"+store.EpochRegistryV1, &rh); err != nil {
		return h, 0, err
	}
	obj, ok, err := r.object(rh.Head)
	if err != nil {
		return h, 0, err
	}
	if !ok || hashref.SumSHA256(obj.Payload).String() != rh.Head {
		return h, 0, fmt.Errorf("epoch registry object mismatch")
	}
	// Read the registry's existing wire shape; nomination authority remains
	// in the daemon's content-addressed registry, not a new CLI dependency.
	var reg struct {
		SemanticID string `json:"semantic_id"`
		Epochs     []struct {
			Epoch      int64    `json:"epoch"`
			Candidates []string `json:"candidates"`
		} `json:"epochs"`
	}
	if err := json.Unmarshal(obj.Payload, &reg); err != nil {
		return h, 0, fmt.Errorf("registry: decode: %w", err)
	}
	for _, ep := range reg.Epochs {
		for _, candidate := range ep.Candidates {
			if candidate == "AILANG v0.41.0" {
				return h, ep.Epoch, nil
			}
		}
	}
	return h, 0, fmt.Errorf("interpreter has no nominated epoch")
}
func demoGrants(ctx context.Context, db, token string) (bool, error) {
	ro, err := store.OpenReadOnly(db)
	if err != nil {
		return false, fmt.Errorf("cannot open read-only authority view")
	}
	defer ro.Close()
	authCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	now := time.Now().Unix()
	out, err := authority.New(ro).ResolveContext(authCtx, "Bearer "+token, now)
	if err != nil {
		return false, fmt.Errorf("read-only session resolution failed")
	}
	if out.Denied != nil {
		return false, fmt.Errorf("%s", out.Denied.String())
	}
	if out.Success == nil || out.Success.EpisodeID != "demo" {
		return false, fmt.Errorf("demo episode session required")
	}
	requirements := []broker.Requirement{{Effect: "world.apply", Scope: "world", Cost: 10}, {Effect: "Workspace.Read", Scope: "worktree", Cost: 1}, {Effect: "Ailang.CLI", Scope: "worktree", Cost: 2}, {Effect: "Human.Approve", Scope: "worktree", Cost: 1}}
	spend, err := ro.EffectSpend(authCtx, out.Success.EpisodeID)
	if err != nil {
		return false, fmt.Errorf("cannot inspect remaining grant budgets")
	}
	caps := append([]broker.Capability(nil), out.Success.Caps...)
	for i := range caps {
		caps[i].Budget -= spend[store.EffectKey{Effect: caps[i].Effect, Scope: caps[i].Scope}]
	}
	snap := broker.NewCapabilitySnapshot(caps, now)
	for _, req := range requirements {
		if !broker.Allows(snap, req).Allowed {
			return false, fmt.Errorf("missing, expired or insufficient %s@%s grant", req.Effect, req.Scope)
		}
	}
	worlds, err := ro.ScanUnreadableWorlds(authCtx, "", 1)
	if err != nil {
		return false, fmt.Errorf("cannot inspect public worlds")
	}
	return worlds.Scanned > 0, nil
}
func demoSeedDomain(ctx context.Context, o demoOptions, token string) error {
	c := demoProofClient(o.addr)
	r := worldReader{ctx: ctx, c: c}
	h, epoch, err := demoPreflight(ctx, o, token, r)
	if err != nil {
		return err
	}
	source, err := os.ReadFile(filepath.Join("world", "logepoch.ail"))
	if err != nil {
		return fmt.Errorf("prepare worktree containing unchanged world/logepoch.ail; run from repository root")
	}
	objects := demoDomainObjects(source)
	observed := ""
	prev := hashref.SumSHA256([]byte("world-demo-genesis")).String()
	for i, obj := range objects {
		actual, has, err := headState(ctx, c)
		if err != nil {
			return err
		}
		if actual != observed || has != (observed != "") {
			return fmt.Errorf("unexpected CAS head; stop and inspect history")
		}
		commit := demoMakeCommit(obj, int64(i), epoch, h.InterpreterRef, observed, prev)
		data, err := json.Marshal(commit)
		if err != nil {
			return err
		}
		status, body, err := c.doAuth(ctx, "POST", "/v1/commit", bytes.NewReader(data), "Bearer "+token)
		if err != nil {
			return fmt.Errorf("commit outcome uncertain; inspect history, never retry automatically")
		}
		if status != 200 {
			var failure daemon.APIError
			_ = json.Unmarshal(body, &failure)
			return fmt.Errorf("commit refused HTTP %d %s; stop and inspect history", status, failure.Error.Class)
		}
		var result struct {
			SelectedHead string `json:"selectedHead"`
		}
		if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {
			return fmt.Errorf("commit response mismatch; inspect history")
		}
		observed = commit.NextWorld.Ref
		prev = commit.Entry.EntryHash
	}
	// S2/S3 run after S1 is recorded. A refusal keeps S1 on record and never
	// writes a manifest; the caller's S2-S6 boundary text stays truthful.
	if _, err := demoRunToolStages(ctx, c, token, source); err != nil {
		return fmt.Errorf("S1 domain commits recorded; S2/S3 refused: %v; S2-S6 incomplete; no success manifest written", err)
	}
	// S4 proves the exact source object committed by S1, after S2/S3 succeed.
	_, err = demoSeedProof(ctx, c, token, h, o.key, hashref.SumSHA256(source))
	return err
}

// demoSeedProve forwards to the real producer. The seed integration control
// observes the constructed adapters here because HTTP cannot expose client settings.
var demoSeedProve = demoProve

var demoSeedRequiredIdentities = []string{"sameRef", "servesEntry"}

func demoSeedProof(ctx context.Context, c *client, token string, h daemon.HealthResponse, keyPath string, source hashref.HashRef) (hashref.HashRef, error) {
	key, err := demoLoadProofKey(keyPath)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 key refused: %w; no success manifest written", err)
	}
	ref, err := demoSeedProve(ctx, demoProofReader{c: c}, demoProofWriter{c: c, token: token, interpreter: h.InterpreterRef}, h, key, source, demoSeedRequiredIdentities)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 refused: %w; stop and inspect history; never retry or reseed; no success manifest written", err)
	}
	return ref, fmt.Errorf("S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action")
}
func demoDomainObjects(source []byte) []demoObject {
	goal := []byte(`{"goal":"Review two candidate approaches using real recorded evidence","author":"demo"}`)
	a := []byte(`{"candidate":"A","description":"Keep the current verified log epoch module","author":"demo"}`)
	b := []byte(`{"candidate":"B","description":"Review a future revision after evidence is available","author":"demo"}`)
	claim := []byte(`{"kind":"authored-claim","author":"demo-authored","text":"Demonstration commentary authored for this demo; not an AI review or proof."}`)
	ids := []string{"world/demo/goal/v1", "world/demo/proposal/v1", "world/demo/proposal/v1", "world/source/ailang/v1", "world/demo/authored-claim/v1"}
	payloads := [][]byte{goal, a, b, source, claim}
	out := make([]demoObject, len(ids))
	for i, id := range ids {
		out[i] = demoObject{Hash: hashref.SumSHA256(payloads[i]).String(), InterfaceHash: hashref.SumSHA256([]byte(id)).String(), SemanticID: id, Provenance: "demo-authored", Payload: payloads[i]}
	}
	return out
}
func demoMakeCommit(o demoObject, index, epoch int64, interpreter, observed, prev string) demoCommit {
	return makeOneObjectCommit(o, index, epoch, interpreter, observed, prev, "demo")
}

// makeOneObjectCommit builds the commit that appends one object as the next
// log entry (its transitionRef), written by writtenBy. Shared by the demo seed
// and proof (writtenBy "demo") and `outcome record` (row 182, "harness:floor").
func makeOneObjectCommit(o demoObject, index, epoch int64, interpreter, observed, prev, writtenBy string) demoCommit {
	header := wireHeader{EntryIndex: index, SemanticsEpoch: epoch, TransitionFn: o.Hash, Interpreter: interpreter, PrevEntryHash: prev, WrittenBy: writtenBy}
	encoded, _ := json.Marshal(struct {
		Header        wireHeader `json:"header"`
		TransitionRef string     `json:"transitionRef"`
	}{header, o.Hash})
	entryHash := hashref.SumSHA256(encoded).String()
	state := hashref.SumSHA256([]byte(observed + "\n" + o.Hash)).String()
	worldBytes, _ := json.Marshal(struct {
		Revision           int64
		StateRoot, LogHead string
	}{index, state, entryHash})
	return demoCommit{ObservedHead: observed, Objects: []demoObject{o}, NextWorld: wireWorld{Ref: hashref.SumSHA256(worldBytes).String(), Revision: index, StateRoot: state, LogHead: entryHash}, Entry: wireEntry{Header: header, EntryHash: entryHash, TransitionRef: o.Hash}}
}

// S2/S3 (M2b): read and tool test/check stages dispatched through the real
// /a2a/ endpoint. The seed holds no tool binary and launches no subprocess;
// every invocation is the daemon's (sole-writer) work and is retained by ID.
const (
	demoLogepochPath = "world/logepoch.ail"
	demoReadTask     = "demo-read"
	demoTestTask     = "demo-test"
	demoCheckTask    = "demo-check"
)

// demoToolRun is one retained A2A invocation: the daemon's invocation ID,
// the world effect record reference, and the CLI result it carried.
type demoToolRun struct {
	Task         string   `json:"task"`
	Skill        string   `json:"skill"`
	Op           string   `json:"op,omitempty"`
	InvocationID string   `json:"invocationId"`
	EffectRecord string   `json:"effectRecord"`
	Content      string   `json:"content,omitempty"`
	Argv         []string `json:"argv,omitempty"`
	ExitCode     int      `json:"exitCode"`
	OK           bool     `json:"ok"`
	Stdout       string   `json:"stdout,omitempty"`
	Stderr       string   `json:"stderr,omitempty"`
}

type demoToolStages struct {
	Read, Test, Check demoToolRun
}

// demoTestReport is the subset of `ailang test --json` the seed
// relies on. The validator decides; it never trusts a bare success flag alone.
type demoTestReport struct {
	Success      bool `json:"success"`
	Passed       int  `json:"passed_tests"`
	Failed       int  `json:"failed_tests"`
	Skipped      int  `json:"skipped_tests"`
	VacuousSkips int  `json:"vacuous_skips"`
}

// demoValidateTestReport refuses a null, failing, skipped, vacuous or
// zero-pass report. The null check decodes into a pointer first so it stays a
// named refusal rather than falling out of the zero-value checks below.
func demoValidateTestReport(stdout string) (demoTestReport, error) {
	var probe *demoTestReport
	if err := json.Unmarshal([]byte(stdout), &probe); err != nil {
		return demoTestReport{}, fmt.Errorf("test report is not JSON: %w", err)
	}
	if probe == nil {
		return demoTestReport{}, fmt.Errorf("null report is refused")
	}
	var rep demoTestReport
	if err := json.Unmarshal([]byte(stdout), &rep); err != nil {
		return demoTestReport{}, fmt.Errorf("test report is not JSON: %w", err)
	}
	switch {
	case !rep.Success:
		return rep, fmt.Errorf("test report success=false is refused")
	case rep.Failed != 0:
		return rep, fmt.Errorf("failed_tests=%d is refused", rep.Failed)
	case rep.Skipped != 0:
		return rep, fmt.Errorf("skipped_tests=%d is refused", rep.Skipped)
	case rep.VacuousSkips != 0:
		return rep, fmt.Errorf("vacuous_skips=%d is refused", rep.VacuousSkips)
	case rep.Passed == 0:
		return rep, fmt.Errorf("zero named passes are refused")
	}
	return rep, nil
}

// demoA2ASend dispatches one tasks/send over real /a2a/ and returns the
// retained run plus the artifact data. Any non-completed task is a refusal.
func demoA2ASend(ctx context.Context, c *client, token, task, skill string, data map[string]any) (demoToolRun, map[string]any, error) {
	run := demoToolRun{Task: task, Skill: skill}
	body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 9, "method": "tasks/send", "params": map[string]any{
		"id": task, "metadata": map[string]any{"skill_id": skill},
		"message": map[string]any{"role": "user", "parts": []any{map[string]any{"type": "data", "data": data}}}}})
	if err != nil {
		return run, nil, err
	}
	status, raw, err := c.doAuth(ctx, "POST", "/a2a/", bytes.NewReader(body), "Bearer "+token)
	if err != nil {
		return run, nil, fmt.Errorf("%s dispatch outcome uncertain; inspect history, never retry automatically", task)
	}
	var wire struct {
		Result *struct {
			Status    struct{ State string } `json:"status"`
			Artifacts []struct {
				Parts []struct {
					Data map[string]any `json:"data"`
				} `json:"parts"`
			} `json:"artifacts"`
			Metadata map[string]any `json:"metadata"`
		} `json:"result"`
		Error any `json:"error"`
	}
	if status != 200 || json.Unmarshal(raw, &wire) != nil || wire.Error != nil || wire.Result == nil ||
		wire.Result.Status.State != "completed" || len(wire.Result.Artifacts) != 1 || len(wire.Result.Artifacts[0].Parts) != 1 {
		return run, nil, fmt.Errorf("%s: tasks/send refused (HTTP %d)", task, status)
	}
	out := wire.Result.Artifacts[0].Parts[0].Data
	run.InvocationID, _ = wire.Result.Metadata["invocation_id"].(string)
	run.EffectRecord = demoEffectRecord(out)
	if run.InvocationID == "" || run.EffectRecord == "" {
		return run, nil, fmt.Errorf("%s: invocation or effect record reference absent", task)
	}
	return run, out, nil
}

// demoEffectRecord returns world.effects[0].record of a tool output, or "".
func demoEffectRecord(out map[string]any) string {
	w, _ := out["world"].(map[string]any)
	effects, _ := w["effects"].([]any)
	if len(effects) != 1 {
		return ""
	}
	e, _ := effects[0].(map[string]any)
	rec, _ := e["record"].(string)
	return rec
}

func demoFillCLI(run *demoToolRun, out map[string]any) {
	if raw, ok := out["argv"].([]any); ok {
		for _, a := range raw {
			s, _ := a.(string)
			run.Argv = append(run.Argv, s)
		}
	}
	if code, ok := out["exit_code"].(float64); ok {
		run.ExitCode = int(code)
	}
	run.OK, _ = out["ok"].(bool)
	run.Stdout, _ = out["stdout"].(string)
	run.Stderr, _ = out["stderr"].(string)
}

// demoTestStage runs the S3 test (op "test") or check (op "ai_check") through
// the ailang-cli skill. The op is retained before any dispatch, so a refused
// op is still on record.
func demoTestStage(ctx context.Context, c *client, token, task, op, path string) (demoToolRun, error) {
	run := demoToolRun{Task: task, Skill: "ailang-cli", Op: op}
	if op != "test" && op != "ai_check" {
		return run, fmt.Errorf("refuses unlisted op %q", op)
	}
	data := map[string]any{"op": op, "path": path}
	if op == "test" {
		data["flags"] = map[string]any{"json": ""}
	}
	call, out, err := demoA2ASend(ctx, c, token, task, "ailang-cli", data)
	if err != nil {
		return run, err
	}
	call.Op = op
	demoFillCLI(&call, out)
	if !call.OK || call.ExitCode != 0 {
		return call, fmt.Errorf("%s: CLI refused or exited %d", task, call.ExitCode)
	}
	return call, nil
}

// demoRunToolStages executes S2 (read, bytes must equal the committed source)
// and S3 (test then ai_check). It stops at the first refusal.
func demoRunToolStages(ctx context.Context, c *client, token string, source []byte) (demoToolStages, error) {
	return demoRunStages(ctx, c, token, source, "")
}

// demoRunStages is demoRunToolStages with a task-ID suffix, so a control arm
// can use fresh A2A task IDs. The S2 mismatch is named before S3 is dispatched.
func demoRunStages(ctx context.Context, c *client, token string, source []byte, suffix string) (demoToolStages, error) {
	var s demoToolStages
	read, out, err := demoA2ASend(ctx, c, token, demoReadTask+suffix, "ailang-read", map[string]any{"path": demoLogepochPath})
	if err != nil {
		return s, err
	}
	read.Content, _ = out["content"].(string)
	s.Read = read
	if read.Content != string(source) {
		return s, fmt.Errorf("S2 read bytes differ from committed %s", demoLogepochPath)
	}
	s.Test, err = demoTestStage(ctx, c, token, demoTestTask+suffix, "test", demoLogepochPath)
	if err != nil {
		return s, err
	}
	if _, err = demoValidateTestReport(s.Test.Stdout); err != nil {
		return s, err
	}
	s.Check, err = demoTestStage(ctx, c, token, demoCheckTask+suffix, "ai_check", demoLogepochPath)
	if err != nil {
		return s, err
	}
	return s, demoValidateAICheck(s.Check.Stdout)
}

// demoValidateAICheck requires a non-null ai_check object whose check passed
// and whose verifier ran with at least one verified function and no failures.
func demoValidateAICheck(stdout string) error {
	var probe *map[string]any
	if err := json.Unmarshal([]byte(stdout), &probe); err != nil {
		return fmt.Errorf("ai_check output is not JSON: %w", err)
	}
	if probe == nil {
		return fmt.Errorf("null ai_check report is refused")
	}
	var rep struct {
		Check struct {
			Passed     bool `json:"passed"`
			ErrorCount int  `json:"error_count"`
		} `json:"check"`
		Verify struct {
			Available      bool `json:"available"`
			Verified       int  `json:"verified"`
			Counterexample int  `json:"counterexample"`
			Errors         int  `json:"errors"`
		} `json:"verify"`
	}
	if err := json.Unmarshal([]byte(stdout), &rep); err != nil {
		return fmt.Errorf("ai_check output is not JSON: %w", err)
	}
	switch {
	case !rep.Check.Passed:
		return fmt.Errorf("ai_check check.passed=false is refused")
	case rep.Check.ErrorCount != 0:
		return fmt.Errorf("ai_check check.error_count=%d is refused", rep.Check.ErrorCount)
	case !rep.Verify.Available:
		return fmt.Errorf("ai_check verify unavailable is refused")
	case rep.Verify.Verified == 0:
		return fmt.Errorf("ai_check verify.verified=0 is refused")
	case rep.Verify.Counterexample != 0 || rep.Verify.Errors != 0:
		return fmt.Errorf("ai_check verify counterexample or errors is refused")
	}
	return nil
}
