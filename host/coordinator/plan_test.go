package coordinator

import (
	"encoding/json"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"testing"
)

func lawPlan() (store.World, transitionreg.Descriptor, plan) {
	w := store.World{Ref: hashref.SumSHA256([]byte("w")), Revision: 7,
		StateRoot: hashref.SumSHA256([]byte("s")), LogHead: hashref.SumSHA256([]byte("l"))}
	d := transitionreg.Descriptor{ID: "x", TransitionFn: hashref.SumSHA256([]byte("fn")),
		Interpreter: hashref.SumSHA256([]byte("i")), SemanticsEpoch: 1}
	return w, d, planInvocation(w, "a2a:e:t", "e", d, []byte(`{"a":1}`), []byte(`{"b":2}`), 5)
}

// The record's skillId is the descriptor's ID — the registry key the card
// emits as skills[].id (projection.go:240).
func TestPlanRecordSkillID(t *testing.T) {
	_, d, p := lawPlan()
	var rec record
	if err := json.Unmarshal(p.Objects[2].Payload, &rec); err != nil || p.Objects[2].SemanticID != RecordV1 {
		t.Fatalf("record object: %v / %s", err, p.Objects[2].SemanticID)
	}
	if rec.SkillID != d.ID || rec.SkillID == "" {
		t.Fatalf("record skillId = %q, want descriptor ID %q", rec.SkillID, d.ID)
	}
}

func TestPlanLawRevision(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.NextWorld.Revision != w.Revision+1 || p.Commit.Entry.Header.EntryIndex != w.Revision+1 {
		t.Fatalf("applyRevision: next revision %d / entry %d, want %d", p.Commit.NextWorld.Revision, p.Commit.Entry.Header.EntryIndex, w.Revision+1)
	}
}

func TestPlanLawStateRoot(t *testing.T) {
	_, _, p := lawPlan()
	if p.Commit.NextWorld.StateRoot != hashref.SumSHA256([]byte(`{"b":2}`)) {
		t.Fatal("applyRevision: stateRoot is not the output object")
	}
}

func TestPlanLawLogHead(t *testing.T) {
	_, _, p := lawPlan()
	if p.Commit.NextWorld.LogHead != p.Commit.Entry.EntryHash || p.Commit.Entry.TransitionRef != p.Record {
		t.Fatal("applyRevision: logHead is not the new entry / entry does not reference the record")
	}
}

func TestPlanLawPrevEntry(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.Entry.Header.PrevEntryHash != w.LogHead || p.Intent.PrevEntryHash != w.LogHead {
		t.Fatal("entry does not chain from the observed log head")
	}
}

func TestPlanLawObservedHead(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.ObservedHead != w.Ref || p.Intent.ObservedHead != w.Ref || p.Intent.WorldRef != p.Commit.NextWorld.Ref {
		t.Fatal("proposalMatchesWorld: the commit/intent do not observe exactly the planned-against world")
	}
}
