package main

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorFinishManifestCapture(t *testing.T) {
	f := newFinishFixture(t)
	f.gitOutput("-C", f.wt, "branch", f.rep.Start.Base.Branch, f.before)
	ctx := activationTestContext(t)
	g := operatorFinishGitFixture(t, f)
	in := operatorFinishInputs{report: f.rep, git: g, profile: &broker.ExecProfile{Protected: []string{"a.go"}}, world: f.db,
		manifest: store.OperatorFinishReview{Base: f.rep.Start.Base, Commit: store.OperatorLocalCommit{Message: "Reviewed finish"}}}
	m, err := in.frozen(ctx)
	if err != nil || !reflect.DeepEqual(m.Protected, []string{"a.go"}) || !reflect.DeepEqual(m.Unattributed, []string{"b.go"}) ||
		!strings.HasPrefix(m.Commit.Message, "Reviewed finish\n\n") || m.Commit.Snapshot.Head != f.before {
		t.Fatalf("signed snapshot and provenance: %+v %v", m, err)
	}
	before := serviceFileHash(filepath.Join(g.admin, "index"))
	in.manifest.Base.SHA = strings.Repeat("d", 40)
	if _, err := in.frozen(ctx); !errors.Is(err, errOperatorFinishReview) || serviceFileHash(filepath.Join(g.admin, "index")) != before {
		t.Fatalf("wrong base changed evidence: %v", err)
	}
}

func TestOperatorFinishWorldEvidence(t *testing.T) {
	ctx := activationTestContext(t)
	path := filepath.Join(t.TempDir(), "world.db")
	world, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer world.Close()
	h := func(s string) hashref.HashRef { return hashref.SumSHA256([]byte(s)) }
	prev := store.World{Ref: h("genesis"), StateRoot: h("state"), LogHead: h("log")}
	if err := world.PutWorld(activationTestContext(t), prev); err != nil {
		t.Fatal(err)
	}
	if err := world.SelectHead(activationTestContext(t), prev.Ref); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(finishOutcome(ac5Report(), store.StartCheck{Name: "short", Argv: []string{"echo", "ok"}}, checkrun.Result{}, true))
	if err != nil {
		t.Fatal(err)
	}
	body, _, err := parseTaskOutcome(raw)
	if err != nil {
		t.Fatal(err)
	}
	result := taskOutcomeObject(body)
	ref := hashref.SumSHA256(result.Payload)
	obj := store.Object{Hash: ref, InterfaceHash: h("interface"), SemanticID: result.SemanticID, Provenance: result.Provenance, Payload: result.Payload}
	if err := world.PutObject(activationTestContext(t), obj); err != nil {
		t.Fatal(err)
	}
	r := store.OperatorFinishCheck{Name: "short", Result: "passed", Ref: ref.String(), Entry: 1}
	if finishWorldEvidence(ctx, path, r) {
		t.Fatal("loose object is not a check receipt")
	}
	next := store.World{Ref: h("next"), Revision: 1, StateRoot: h("next-state"), LogHead: h("next-log")}
	if err := world.Commit(activationTestContext(t), store.Commit{ObservedHead: prev.Ref, Objects: []store.Object{obj}, NextWorld: next,
		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 1, SemanticsEpoch: 1, TransitionFn: ref,
			Interpreter: h("interpreter"), PrevEntryHash: prev.LogHead, WrittenBy: taskOutcomeWrittenBy},
			EntryHash: next.LogHead, TransitionRef: ref}}); err != nil {
		t.Fatal(err)
	}
	if !finishWorldEvidence(ctx, path, r) {
		t.Fatal("exact committed outcome was not readable")
	}
	for _, wrong := range []store.OperatorFinishCheck{
		{Name: r.Name, Result: r.Result, Ref: r.Ref, Entry: 2},
		{Name: "other", Result: r.Result, Ref: r.Ref, Entry: 1},
		{Name: r.Name, Result: "failed", Ref: r.Ref, Entry: 1},
	} {
		if finishWorldEvidence(ctx, path, wrong) {
			t.Fatalf("another check/entry reused a valid outcome object: %+v", wrong)
		}
	}
}
