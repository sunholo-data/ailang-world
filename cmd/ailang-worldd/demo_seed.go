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
	c := newClient(o.addr)
	c.timeout = 10 * time.Second
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
	return nil
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
	header := wireHeader{EntryIndex: index, SemanticsEpoch: epoch, TransitionFn: o.Hash, Interpreter: interpreter, PrevEntryHash: prev, WrittenBy: "demo"}
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
