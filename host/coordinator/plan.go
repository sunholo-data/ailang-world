package coordinator

import (
	"encoding/json"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Semantic IDs of the three immutable objects one invocation records.
const (
	InputV1   = "world/invocation-input/v1"
	OutputV1  = "world/invocation-output/v1"
	RecordV1  = "world/invocation-record/v1"
	writtenBy = "coordinator:a2a"
)

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

// plan is the complete, pure description of one invocation commit.
type plan struct {
	Objects []store.Object
	Commit  store.Commit
	Intent  store.JournalIntent
	Record  hashref.HashRef
}

func object(semanticID string, payload []byte) store.Object {
	return store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID: semanticID, Provenance: writtenBy, Payload: payload,
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
func planInvocation(w store.World, id, episodeID string, d transitionreg.Descriptor,
	input, output []byte, logicalTime int64) plan {
	in, out := object(InputV1, input), object(OutputV1, output)
	rec := object(RecordV1, mustJSON(record{
		InvocationID: id, EpisodeID: episodeID, SkillID: d.ID,
		TransitionFn: d.TransitionFn.String(), Interpreter: d.Interpreter.String(),
		SemanticsEpoch: d.SemanticsEpoch, Input: in.Hash.String(), Output: out.Hash.String(),
	}))
	header := store.LogHeader{
		EntryIndex: w.Revision + 1, SemanticsEpoch: d.SemanticsEpoch,
		TransitionFn: d.TransitionFn, Interpreter: d.Interpreter,
		PrevEntryHash: w.LogHead, WrittenBy: writtenBy,
	}
	entryHash := hashref.SumSHA256(mustJSON(entryWire{
		EntryIndex: header.EntryIndex, SemanticsEpoch: header.SemanticsEpoch,
		TransitionFn: header.TransitionFn.String(), Interpreter: header.Interpreter.String(),
		PrevEntryHash: header.PrevEntryHash.String(), WrittenBy: header.WrittenBy,
		TransitionRef: rec.Hash.String(),
	}))
	next := store.World{Revision: w.Revision + 1, StateRoot: out.Hash, LogHead: entryHash}
	next.Ref = hashref.SumSHA256(mustJSON(worldWire{
		Revision: next.Revision, StateRoot: next.StateRoot.String(), LogHead: next.LogHead.String(),
	}))
	objects := []store.Object{in, out, rec}
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
