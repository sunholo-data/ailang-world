package coordinator

import (
	"encoding/json"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Semantic IDs of the three immutable objects one invocation records.
const (
	InputV1  = "world/invocation-input/v1"
	OutputV1 = "world/invocation-output/v1"
	RecordV1 = "world/invocation-record/v1"
	RecordV2 = "world/invocation-record/v2"
)

// writtenBy is the entry header's WrittenBy and the objects' (advisory,
// first-writer-wins) Provenance for an invocation on surface.
func writtenBy(surface Surface) string { return "coordinator:" + string(surface) }

// record is the canonical invocation record (the log entry's TransitionRef).
// Field order is the canonical encoding order.
type record struct {
	InvocationID   string `json:"invocationId"`
	EpisodeID      string `json:"episodeId"`
	SkillID        string `json:"skillId"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	Input          string `json:"input"`
	Output         string `json:"output"`
}

// recordV2 is the effectful invocation record (row 134 §4.1): the v1 fields
// in v1 order, then the plan object's ref and the ordered effect-record refs
// (never null: an empty plan records []).
type recordV2 struct {
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

// worldRef is the content address of a coordinator-planned world row.
func worldRef(w store.World) hashref.HashRef {
	return hashref.SumSHA256(mustJSON(worldWire{
		Revision: w.Revision, StateRoot: w.StateRoot.String(), LogHead: w.LogHead.String(),
	}))
}

// plan is the complete, pure description of one invocation commit.
type plan struct {
	Objects []store.Object
	Commit  store.Commit
	Intent  store.JournalIntent
	Record  hashref.HashRef
}

func object(by, semanticID string, payload []byte) store.Object {
	return store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID: semanticID, Provenance: by, Payload: payload,
	}
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // unreachable: only fixed string/int structs are encoded
	}
	return b
}

// planInvocation transcribes world/transitions.ail's laws onto store types:
// applyRevision (revision+1, stateRoot = output, logHead = the new entry) and
// proposalMatchesWorld (the commit observes exactly w, enforced by the
// store's compare-and-append at the boundary). Pure: no clock, no I/O.
func planInvocation(w store.World, surface Surface, id, episodeID string, d transitionreg.Descriptor,
	input, output []byte, logicalTime int64) plan {
	by := writtenBy(surface)
	in, out := object(by, InputV1, input), object(by, OutputV1, output)
	rec := object(by, RecordV1, mustJSON(record{
		InvocationID: id, EpisodeID: episodeID, SkillID: d.ID,
		TransitionFn: d.TransitionFn.String(), Interpreter: d.Interpreter.String(),
		SemanticsEpoch: d.SemanticsEpoch, Input: in.Hash.String(), Output: out.Hash.String(),
	}))
	return planCommit(w, by, id, d, []store.Object{in, out, rec}, rec, out, logicalTime)
}

// planEffectInvocation is planInvocation for an effectful descriptor (row
// 134): the record is world/invocation-record/v2 — the v1 fields plus the
// plan object's ref and the ordered effect-record refs — and the plan object
// is committed alongside input and output. The world laws are unchanged.
func planEffectInvocation(w store.World, surface Surface, id, episodeID string, d transitionreg.Descriptor,
	input, output, planObject []byte, effects []hashref.HashRef, logicalTime int64) plan {
	by := writtenBy(surface)
	in, out, po := object(by, InputV1, input), object(by, OutputV1, output), object(by, EffectPlanV1, planObject)
	refs := make([]string, len(effects))
	for i, e := range effects {
		refs[i] = e.String()
	}
	rec := object(by, RecordV2, mustJSON(recordV2{
		InvocationID: id, EpisodeID: episodeID, SkillID: d.ID,
		TransitionFn: d.TransitionFn.String(), Interpreter: d.Interpreter.String(),
		SemanticsEpoch: d.SemanticsEpoch, Input: in.Hash.String(), Output: out.Hash.String(),
		Plan: po.Hash.String(), Effects: refs,
	}))
	return planCommit(w, by, id, d, []store.Object{in, out, po, rec}, rec, out, logicalTime)
}

func planCommit(w store.World, by, id string, d transitionreg.Descriptor, objects []store.Object,
	rec, out store.Object, logicalTime int64) plan {
	header := store.LogHeader{
		EntryIndex: w.Revision + 1, SemanticsEpoch: d.SemanticsEpoch,
		TransitionFn: d.TransitionFn, Interpreter: d.Interpreter,
		PrevEntryHash: w.LogHead, WrittenBy: by,
	}
	entryHash := hashref.SumSHA256(mustJSON(entryWire{
		EntryIndex: header.EntryIndex, SemanticsEpoch: header.SemanticsEpoch,
		TransitionFn: header.TransitionFn.String(), Interpreter: header.Interpreter.String(),
		PrevEntryHash: header.PrevEntryHash.String(), WrittenBy: header.WrittenBy,
		TransitionRef: rec.Hash.String(),
	}))
	next := store.World{Revision: w.Revision + 1, StateRoot: out.Hash, LogHead: entryHash}
	next.Ref = worldRef(next)
	return plan{
		Objects: objects,
		Record:  rec.Hash,
		Commit: store.Commit{
			InvocationID: id, ObservedHead: w.Ref, Objects: objects, NextWorld: next,
			Entry: store.LogEntry{Header: header, EntryHash: entryHash, TransitionRef: rec.Hash},
		},
		Intent: store.JournalIntent{
			InvocationID: id, WorldRef: next.Ref, EntryHash: entryHash, ObservedHead: w.Ref,
			PrevEntryHash: header.PrevEntryHash, TransitionFn: header.TransitionFn,
			TransitionRef: rec.Hash, Interpreter: header.Interpreter, LogicalTime: logicalTime,
		},
	}
}
