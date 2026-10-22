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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/evidence"
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
	if manifest.Version != h.InterpreterVersion || !strings.HasPrefix(manifest.Version, demoProofVersion+"\n") {
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
			if candidate == demoProofVersion {
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
	// The credential lives in the credentials store beside the world store
	// (row 189), read-only as the daemon reads it; spend stays in the world.
	credPath, err := store.DefaultCredentialsPath(db)
	if err != nil {
		return false, fmt.Errorf("cannot locate the credentials store")
	}
	out, err := authority.New(authority.ReadOnlyAt(credPath)).ResolveContext(authCtx, "Bearer "+token, now)
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
func demoSeedDomain(ctx context.Context, o demoOptions, token string) (err error) {
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
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w; inspect history before any further action; never blindly retry", err)
		}
	}()
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
	// writes a success manifest.
	stages, err := demoRunToolStages(ctx, c, token, source)
	if err != nil {
		return fmt.Errorf("S1 domain commits recorded; S2/S3 refused: %v; no success manifest written", err)
	}
	// S4 proves the exact source object committed by S1, after S2/S3 succeed.
	envelope, err := demoSeedProof(ctx, c, token, h, o.key, hashref.SumSHA256(source))
	if err != nil {
		return err
	}
	return demoSeedApproval(ctx, c, token, objects, stages, envelope, o.manifest)
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
	return ref, nil
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
	EntryIndex   int64    `json:"entryIndex"`
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
	index, present := wire.Result.Metadata["entry_index"].(float64)
	if !present || index < 0 || index != float64(int64(index)) { // M4b metadata-index
		return run, nil, fmt.Errorf("%s: entry_index must be present and integral", task)
	}
	run.EntryIndex = int64(index)
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

// S5 passes descriptions and references from the completed stages to the
// published convention-v2 approval transition. Authority and packet production
// remain in the coordinator/broker; the seed neither decides nor retries.
type demoEvidenceRef struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

func demoApprovalRoster(test, effect, claim string, envelope hashref.HashRef) []demoEvidenceRef {
	return []demoEvidenceRef{
		{Kind: "test-invocation", Ref: test},
		{Kind: "recorded-effect", Ref: effect},
		{Kind: "authored-claim", Ref: claim},
		{Kind: "proof-envelope", Ref: envelope.String()},
	}
}

func demoApprovalInput(goal string, candidates []string, roster []demoEvidenceRef) (map[string]any, error) {
	if len(roster) != 4 {
		return nil, fmt.Errorf("S5 roster must have exactly four entries")
	}
	for i, kind := range []string{"test-invocation", "recorded-effect", "authored-claim", "proof-envelope"} {
		if roster[i].Kind != kind {
			return nil, fmt.Errorf("S5 roster kind at %d must be %s", i, kind)
		}
	}
	return map[string]any{"goal": goal, "candidates": candidates, "roster": roster}, nil
}

func demoSeedApproval(ctx context.Context, c *client, token string, objects []demoObject, stages demoToolStages, envelope hashref.HashRef, manifest string) error {
	// These are the seed's own fixed S1 payloads, already committed unchanged.
	var goal struct{ Goal string }
	var a, b struct{ Description string }
	for i, dst := range []any{&goal, &a, &b} {
		if err := json.Unmarshal(objects[i].Payload, dst); err != nil { // M4b s1-decode
			return fmt.Errorf("S5 S1 payload %d decode refused: %w", i, err)
		}
	}
	var test struct {
		TransitionRef string `json:"transitionRef"`
		Header        struct {
			EntryIndex *int64 `json:"entryIndex"`
		} `json:"header"`
	}
	if err := demoGet(worldReader{ctx: ctx, c: c}, fmt.Sprintf("/v1/log/%d", stages.Test.EntryIndex), &test); err != nil {
		return fmt.Errorf("S1-S4 recorded; S5 test invocation unavailable: %w; no success manifest written", err)
	}
	if test.Header.EntryIndex == nil || *test.Header.EntryIndex != stages.Test.EntryIndex { // M4b test-index
		return fmt.Errorf("S5 test entry entryIndex must be present and match")
	}
	roster := demoApprovalRoster(test.TransitionRef, stages.Read.EffectRecord, objects[4].Hash, envelope)
	input, err := demoApprovalInput(goal.Goal, []string{a.Description, b.Description}, roster)
	if err != nil {
		return err
	}
	approval, _, err := demoA2ASend(ctx, c, token, "demo-approval", "demo-approve", input)
	if err != nil {
		return fmt.Errorf("S5 refused: %w; no success manifest written", err)
	}
	r := worldReader{ctx: ctx, c: c}
	m, err := demoSeedS6(r, token, objects, stages, approval, roster)
	if err != nil {
		return fmt.Errorf("S6 refused: %w; no success manifest written", err)
	}
	return demoFinalizeSeed(r, m, token, manifest)
}

func demoRequestApproval(ctx context.Context, c *client, token, goal string, candidates []string, roster []demoEvidenceRef) error {
	input, err := demoApprovalInput(goal, candidates, roster)
	if err != nil {
		return err
	}
	_, _, err = demoA2ASend(ctx, c, token, "demo-approval", "demo-approve", input)
	if err != nil {
		return fmt.Errorf("S1-S4 recorded; S5 refused: %w; S6 incomplete; no success manifest written; inspect history before any further action", err)
	}
	return demoApprovalIncomplete()
}

func demoApprovalIncomplete() error {
	return fmt.Errorf("S1-S5 recorded; S6 incomplete; no success manifest written; inspect history before any further action")
}

// demoValidateWhy validates only the real CLI report. Coordinator is deliberately
// first, including for empty/broken reports, before any packet inspection.
func demoValidateWhy(ch whyChain) error {
	if !ch.Coordinator { // M4b coordinator
		return fmt.Errorf("entry %d: coordinator must be true", ch.Entry)
	}
	if !ch.OK { // M4b why-ok
		return fmt.Errorf("entry %d: ok must be true", ch.Entry)
	}
	if len(ch.Links) == 0 { // M4b zero-links
		return fmt.Errorf("entry %d: zero links refused", ch.Entry)
	}
	expected := []string{"effect", "world", "entry", "record", "input", "plan", "output"}
	for _, name := range expected {
		found := false
		for _, link := range ch.Links {
			if link.Name == name {
				found = true
				if !link.OK { // M4b link-false
					return fmt.Errorf("entry %d: expected %s link must exist and be true", ch.Entry, name)
				}
			}
		}
		if !found { // M4b link-absent
			return fmt.Errorf("entry %d: expected %s link must exist and be true", ch.Entry, name)
		}
	}
	return nil
}

func demoWhy(r worldReader, index int64) (whyChain, []byte, error) {
	e, how, err := r.locate(whyTarget{kind: "index", index: index}, whyScanDefault)
	if err != nil {
		return whyChain{}, nil, err
	}
	ch, err := r.buildChain(e, nil)
	if err != nil {
		return ch, nil, err
	}
	ch.MatchedBy = how
	raw, err := json.MarshalIndent(ch, "", "  ")
	raw = append(raw, '\n') // byte-for-byte runWhy --json output
	return ch, raw, err
}

type demoManifestRef struct {
	Role  string `json:"role"`
	Kind  string `json:"kind"`
	Ref   string `json:"ref"`
	Index int64  `json:"index,omitempty"`
}
type demoManifestWhy struct {
	Role   string `json:"role"`
	Entry  int64  `json:"entry"`
	SHA256 string `json:"sha256"`
}
type demoPublicManifest struct {
	Schema       string            `json:"schema"`
	Version      int               `json:"version"`
	ProofVersion string            `json:"proofVersion"`
	ToolVersion  string            `json:"toolVersion"`
	References   []demoManifestRef `json:"references"`
	Why          []demoManifestWhy `json:"why"`
}

// Assembly retains references, never raw reports (which contain input bytes).
// Each declared hash is resolved again independently before returning to seed.
func demoAssembleManifest(r worldReader, objects []demoObject, runs []demoToolRun, roster []demoEvidenceRef, packet, request, proposal string) (demoPublicManifest, error) {
	m := demoPublicManifest{Schema: "world/demo/public-manifest/v1", Version: 1, ProofVersion: demoProofVersion, ToolVersion: daemon.ToolBinaryRelease}
	add := func(role, kind, ref string, index int64) {
		m.References = append(m.References, demoManifestRef{role, kind, ref, index})
	}
	head, _, err := headState(r.ctx, r.c)
	if err != nil {
		return m, err
	}
	add("frozen-head", "world", head, 0)
	entries, err := r.page(0, 500)
	if err != nil {
		return m, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Header.EntryIndex < entries[j].Header.EntryIndex })
	for _, e := range entries {
		add("entry", "entry", e.EntryHash, e.Header.EntryIndex)
		add("object", "object", e.TransitionRef, 0)
	}
	for i, o := range objects {
		add([]string{"goal", "candidate-a", "candidate-b", "source", "claim"}[i], "object", o.Hash, 0)
	}
	for _, e := range roster {
		add(e.Kind, "object", e.Ref, 0)
	}
	for _, pair := range []struct{ role, ref string }{{"packet", packet}, {"request", request}, {"proposal", proposal}} {
		add(pair.role, "object", pair.ref, 0)
	}
	for i, run := range runs {
		ch, raw, err := demoWhy(r, run.EntryIndex)
		if err != nil {
			return m, err
		}
		if err := demoValidateWhy(ch); err != nil {
			return m, err
		}
		m.Why = append(m.Why, demoManifestWhy{[]string{"read", "test", "check", "approval"}[i], run.EntryIndex, sum(raw)})
		e, _, err := r.entry(run.EntryIndex)
		if err != nil {
			return m, err
		}
		rec, _, err := r.record(e)
		if err != nil {
			return m, err
		}
		add("invocation", "object", e.TransitionRef, run.EntryIndex)
		add("transition", "object", rec.TransitionFn, run.EntryIndex)
		add("interpreter", "pin", rec.Interpreter, run.EntryIndex)
		for _, link := range ch.Links {
			if link.Name != "entry" {
				kind := "object"
				if link.Name == "world" {
					kind = "world"
				}
				add(link.Name, kind, link.Ref, run.EntryIndex)
			}
		}
		for _, ref := range rec.Effects {
			obj, _, err := r.object(ref)
			if err != nil {
				return m, err
			}
			fx, err := broker.DecodeRecord(obj.Payload)
			if err != nil {
				return m, err
			}
			add("effect-request", "object", fx.RequestRef.String(), run.EntryIndex)
			if !fx.ResultRef.IsZero() {
				add("effect-result", "object", fx.ResultRef.String(), run.EntryIndex)
			}
		}
	}
	return m, nil
}

var demoPrivatePublicPattern = regexp.MustCompile(`(?i)(bearer |session|private.?key|operator.?key|credential|(?:[a-z]:[\\/]|(?:"|\s)/)|\d{4}-\d\d-\d\d[T ][0-9:]|capture.?time)`)

func demoValidatePublicBytes(raw []byte, token string) error {
	if (token != "" && bytes.Contains(raw, []byte(token))) || demoPrivatePublicPattern.Match(raw) { // M4b public-bytes
		return fmt.Errorf("public bytes contain private material, absolute path or capture time")
	}
	return nil
}

func demoValidateManifest(r worldReader, m demoPublicManifest, token string) error {
	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := demoValidatePublicBytes(raw, token); err != nil {
		return err
	}
	if m.Schema != "world/demo/public-manifest/v1" || m.Version != 1 || m.ProofVersion != demoProofVersion || m.ToolVersion != daemon.ToolBinaryRelease { // M4b manifest-schema
		return fmt.Errorf("public manifest schema/version/pins must match")
	}
	roles := map[string]bool{}
	for _, role := range strings.Fields("frozen-head entry object goal candidate-a candidate-b source claim test-invocation recorded-effect authored-claim proof-envelope packet request proposal invocation transition interpreter world record input plan effect output effect-request effect-result") {
		roles[role] = true
	}
	for _, ref := range m.References {
		if !roles[ref.Role] { // M4b manifest-role
			return fmt.Errorf("public manifest role must be a declared reference name")
		}
		valid := false
		switch ref.Kind {
		case "object":
			obj, ok, err := r.object(ref.Ref)
			valid = err == nil && ok && obj.Hash == ref.Ref && sum(obj.Payload) == ref.Ref
		case "world":
			w, ok, err := r.world(ref.Ref)
			valid = err == nil && ok && w.Ref == ref.Ref && sumJSON(worldWire{w.Revision, w.StateRoot, w.LogHead}) == ref.Ref
		case "entry":
			e, ok, err := r.entry(ref.Index)
			flat := sumJSON(entryWire{e.Header.EntryIndex, e.Header.SemanticsEpoch, e.Header.TransitionFn, e.Header.Interpreter, e.Header.PrevEntryHash, e.Header.WrittenBy, e.TransitionRef})
			nested := sumJSON(struct {
				Header        wireHeader `json:"header"`
				TransitionRef string     `json:"transitionRef"`
			}{e.Header, e.TransitionRef})
			valid = err == nil && ok && e.EntryHash == ref.Ref && (flat == ref.Ref || nested == ref.Ref)
		case "pin":
			h, err := r.health()
			valid = err == nil && h.InterpreterRef == ref.Ref
		}
		if !valid { // M4b manifest-ref
			return fmt.Errorf("every declared manifest reference resolves to matching content: %s %s", ref.Role, ref.Ref)
		}
	}
	for _, w := range m.Why {
		ch, raw, err := demoWhy(r, w.Entry)
		if err != nil {
			return err
		}
		if err := demoValidateWhy(ch); err != nil {
			return err
		}
		if sum(raw) != w.SHA256 { // M4b why-hash
			return fmt.Errorf("entry %d: why raw JSON hash mismatch", w.Entry)
		}
	}
	if err := demoManifestCompleteness(r, m); err != nil { // M4b manifest-completeness
		return err
	}
	return nil
}

// Required singleton roles and per-run references are checked independently of
// assembly. Matching the actual invocation/effects prevents one run's refs from
// standing in for another's, even when their content hashes happen to coincide.
func demoManifestCompleteness(r worldReader, m demoPublicManifest) error {
	counts := map[string]int{}
	for _, ref := range m.References {
		counts[ref.Role]++
	}
	for _, role := range strings.Fields("frozen-head goal candidate-a candidate-b source claim test-invocation recorded-effect authored-claim proof-envelope packet request proposal") {
		if counts[role] != 1 {
			return fmt.Errorf("manifest completeness: role %s requires exactly 1 reference, got %d", role, counts[role])
		}
	}
	seen := map[int64]bool{}
	expected := map[demoManifestRef]int{}
	for _, role := range []string{"read", "test", "check", "approval"} {
		var matches []demoManifestWhy
		for _, w := range m.Why {
			if w.Role == role {
				matches = append(matches, w)
			}
		}
		if len(matches) != 1 {
			return fmt.Errorf("manifest completeness: role why-%s requires exactly 1 hash, got %d", role, len(matches))
		}
		w := matches[0]
		if seen[w.Entry] {
			return fmt.Errorf("manifest completeness: role why-%s duplicates entry %d", role, w.Entry)
		}
		seen[w.Entry] = true
		ch, _, err := demoWhy(r, w.Entry)
		if err != nil {
			return err
		}
		e, _, err := r.entry(w.Entry)
		if err != nil {
			return err
		}
		rec, _, err := r.record(e)
		if err != nil {
			return err
		}
		if rec == nil {
			return fmt.Errorf("manifest completeness: role invocation at entry %d absent", w.Entry)
		}
		add := func(role, kind, ref string) { expected[demoManifestRef{role, kind, ref, w.Entry}]++ }
		add("invocation", "object", e.TransitionRef)
		add("transition", "object", rec.TransitionFn)
		add("interpreter", "pin", rec.Interpreter)
		for _, link := range ch.Links {
			if link.Name == "entry" {
				continue
			}
			kind := "object"
			if link.Name == "world" {
				kind = "world"
			}
			add(link.Name, kind, link.Ref)
		}
		for _, ref := range rec.Effects {
			obj, _, err := r.object(ref)
			if err != nil {
				return err
			}
			fx, err := broker.DecodeRecord(obj.Payload)
			if err != nil {
				return err
			}
			add("effect-request", "object", fx.RequestRef.String())
			if !fx.ResultRef.IsZero() {
				add("effect-result", "object", fx.ResultRef.String())
			}
		}
	}
	if len(m.Why) != 4 {
		return fmt.Errorf("manifest completeness: role why requires exactly 4 hashes, got %d", len(m.Why))
	}
	actual := map[demoManifestRef]int{}
	for _, ref := range m.References {
		switch ref.Role {
		case "invocation", "transition", "interpreter", "world", "record", "input", "plan", "effect", "output", "effect-request", "effect-result":
			actual[ref]++
		}
	}
	// Iterate in manifest/why order for stable, role-specific diagnostics.
	for _, w := range m.Why {
		for _, role := range strings.Fields("invocation transition interpreter world record input plan effect output effect-request effect-result") {
			for ref, n := range expected {
				if ref.Index == w.Entry && ref.Role == role && actual[ref] != n {
					return fmt.Errorf("manifest completeness: role %s for %s entry %d requires %d reference(s), got %d", role, w.Role, w.Entry, n, actual[ref])
				}
			}
		}
	}
	for ref, n := range actual {
		if expected[ref] != n {
			return fmt.Errorf("manifest completeness: extra role %s at entry %d", ref.Role, ref.Index)
		}
	}
	return nil
}

func demoSeedS6(r worldReader, token string, objects []demoObject, stages demoToolStages, approval demoToolRun, roster []demoEvidenceRef) (demoPublicManifest, error) {
	var m demoPublicManifest
	runs := []demoToolRun{stages.Read, stages.Test, stages.Check, approval}
	// Must precede packet count: an ordinary REST entry must fail coordinator.
	for _, run := range runs {
		ch, _, err := demoWhy(r, run.EntryIndex)
		if err != nil {
			return m, err
		}
		if err := demoValidateWhy(ch); err != nil {
			return m, err
		}
	}
	var page struct{ Items []wireObject }
	if err := demoGet(r, "/v1/objects/by-semantic-id/"+coordinator.DecisionPacketV1+"?limit=500", &page); err != nil {
		return m, err
	}
	if len(page.Items) != 1 { // M4b packet-count
		return m, fmt.Errorf("entry %d: exactly one decision packet required", approval.EntryIndex)
	}
	packetRef := ""
	for _, item := range page.Items {
		packetRef = item.Hash
	}
	obj, ok, reason := r.checkedObject(packetRef, coordinator.DecisionPacketV1)
	if !ok { // M4b packet-object
		return m, fmt.Errorf("packet: %s", reason)
	}
	p, err := coordinator.DecodeDecisionPacket(obj.Payload)
	if err != nil { // M4b packet-decode
		return m, err
	}
	e, _, err := r.entry(approval.EntryIndex)
	if err != nil {
		return m, err
	}
	rec, _, err := r.record(e)
	if err != nil {
		return m, err
	}
	status, body, err := r.c.do(r.ctx, "GET", "/workbench?object="+obj.Hash, nil)
	anchor := fmt.Sprintf("decision packet validated: committed with its invocation record at entry #%d;", approval.EntryIndex)
	if err := demoValidatePacketAnchor(approval.EntryIndex, err == nil && status == 200 && bytes.Contains(body, []byte(anchor)), p.ProposalHash.String(), rec.Plan); err != nil {
		return m, err
	}
	m, err = demoAssembleManifest(r, objects, runs, roster, obj.Hash, p.RequestRef.String(), p.ProposalHash.String())
	if err != nil {
		return m, err
	}
	return m, demoValidateManifest(r, m, token)
}

// Keep the final kind refusal independent of the projection's UNAVAILABLE guard.
func demoValidateRosterKind(slot demoEvidenceRef, obj demoObject) error {
	ids := map[string]string{"test-invocation": coordinator.RecordV2, "recorded-effect": broker.EffectRecordV1,
		"authored-claim": "world/demo/authored-claim/v1", "proof-envelope": evidence.ProofSemanticID}
	id := ids[slot.Kind]
	iface := sum([]byte(id))
	if slot.Kind == "proof-envelope" {
		iface = evidence.InterfaceHashV1.String()
	}
	if id == "" || obj.SemanticID != id || obj.InterfaceHash != iface {
		return fmt.Errorf("roster kind: %s.ref must resolve to %s", slot.Kind, id)
	}
	return nil
}

func demoValidateRosterPage(body []byte) error {
	page := string(body)
	if strings.Contains(page, "<p>evidence roster UNAVAILABLE — ") {
		return fmt.Errorf("roster page: evidence roster UNAVAILABLE")
	}
	if strings.Count(page, `<div aria-label="evidence roster">`) != 1 {
		return fmt.Errorf("roster page: evidence roster absent")
	}
	for _, marker := range []string{"<p>recorded-effect: ATTESTED · producer/author: ",
		"<p>authored-claim: CLAIMED · producer/author: demo-authored", "<p>test-invocation: TESTED · producer/author: ",
		"<p>proof-envelope: PROVEN · producer/author: ", " · verdict pass · "} {
		if strings.Count(page, marker) != 1 {
			return fmt.Errorf("roster page: missing resolved marker %s", marker)
		}
	}
	return nil
}

func demoValidateSeedRoster(r worldReader, m demoPublicManifest) error {
	refs := map[string]string{}
	for _, ref := range m.References {
		refs[ref.Role] = ref.Ref
	}
	var input struct {
		Roster []demoEvidenceRef `json:"roster"`
	}
	for _, why := range m.Why {
		if why.Role != "approval" {
			continue
		}
		e, ok, err := r.entry(why.Entry)
		if err != nil || !ok {
			return fmt.Errorf("roster approval entry unavailable: %v", err)
		}
		rec, _, err := r.record(e)
		if err != nil || rec == nil {
			return fmt.Errorf("roster approval record unavailable: %v", err)
		}
		obj, ok, reason := r.checkedObject(rec.Input, coordinator.InputV1)
		if !ok {
			return fmt.Errorf("roster input: %s", reason)
		}
		if err := json.Unmarshal(obj.Payload, &input); err != nil {
			return err
		}
	}
	kinds := []string{"test-invocation", "recorded-effect", "authored-claim", "proof-envelope"}
	if len(input.Roster) != len(kinds) {
		return fmt.Errorf("roster count: require four stored slots")
	}
	for i, slot := range input.Roster {
		if slot.Kind != kinds[i] || refs[slot.Kind] != slot.Ref {
			return fmt.Errorf("roster slot %d: manifest/input mismatch", i)
		}
		var obj demoObject
		err := demoGet(r, "/v1/objects/"+slot.Ref+"?payload=true", &obj)
		if err != nil || obj.Hash != slot.Ref || sum(obj.Payload) != slot.Ref {
			return fmt.Errorf("roster %s.ref unavailable or hash mismatch", slot.Kind)
		}
		if err := demoValidateRosterKind(slot, obj); err != nil {
			return err
		}
	}
	status, body, err := r.c.do(r.ctx, "GET", "/workbench?object="+refs["packet"], nil)
	if err != nil || status != 200 {
		return fmt.Errorf("roster page unavailable: HTTP %d: %v", status, err)
	}
	return demoValidateRosterPage(body)
}

func demoFinalizeSeed(r worldReader, m demoPublicManifest, token, path string) error {
	if err := demoValidateManifest(r, m, token); err != nil {
		return err
	}
	if err := demoValidateSeedRoster(r, m); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return demoPublishManifest(path, append(raw, '\n'))
}

// Link publishes complete, synced public bytes atomically and cannot replace a
// target created since preflight. The scratch file is on the same filesystem.
func demoPublishManifest(path string, raw []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".demo-manifest-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = f.Chmod(0644); err == nil {
		_, err = f.Write(raw)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Link(f.Name(), path); err != nil {
		return fmt.Errorf("publish manifest requires absent target: %w", err)
	}
	return nil
}

// The daemon's existing packet validator establishes carrying-commit membership;
// the seed also binds that verdict to its own S5 index and invocation plan.
func demoValidatePacketAnchor(index int64, carried bool, proposal, plan string) error {
	if !carried || proposal != plan { // M4b packet-anchor
		return fmt.Errorf("entry %d: packet carrying entry/proposal mismatch", index)
	}
	return nil
}
