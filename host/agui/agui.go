// Package agui formats deterministic AG-UI 1.0 SSE bytes from committed entries.
// It performs no I/O and reads no clock.
package agui

import (
	"encoding/json"
	"strconv"
)

const (
	ProtocolVersion = "1.0"
	StateSchema     = "world/agui-state/v1"
)

// World CUSTOM names are a closed table; extensions require a new golden test.
var customNames = []string{"world.entry.committed"}

// Entry is a committed log row. Value must be valid JSON, marshaled from the
// frozen log-route representation by the caller. EntryHash is its canonical ref.
type Entry struct {
	Index     int64
	EntryHash string
	Value     json.RawMessage
}

// Run carries the opaque identifiers echoed by run-scoped events.
type Run struct{ ThreadID, RunID string }

type runStarted struct {
	Type            string `json:"type"`
	ThreadID        string `json:"threadId"`
	RunID           string `json:"runId"`
	ProtocolVersion string `json:"protocolVersion"`
}
type state struct {
	Schema    string  `json:"schema"`
	LastIndex int64   `json:"lastIndex"`
	LogHead   *string `json:"logHead"`
}
type snapshot struct {
	Type     string `json:"type"`
	Snapshot state  `json:"snapshot"`
}
type custom struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}
type patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}
type delta struct {
	Type  string  `json:"type"`
	Delta []patch `json:"delta"`
}
type result struct {
	LastIndex int64 `json:"lastIndex"`
}
type runFinished struct {
	Type     string `json:"type"`
	ThreadID string `json:"threadId"`
	RunID    string `json:"runId"`
	Result   result `json:"result"`
}
type runError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func frame(event any) []byte {
	b, err := json.Marshal(event)
	if err != nil {
		panic("agui: invalid entry JSON: " + err.Error())
	}
	return append(append([]byte("data: "), b...), '\n', '\n')
}

// Start opens the run with a snapshot of the resume cursor, not the latest head.
// A nil snap represents the genesis cursor (-1), whose log head is null.
func Start(run Run, after int64, snap *Entry) []byte {
	var head *string
	if snap != nil {
		h := snap.EntryHash
		head = &h
	}
	b := frame(runStarted{"RUN_STARTED", run.ThreadID, run.RunID, ProtocolVersion})
	return append(b, frame(snapshot{"STATE_SNAPSHOT", state{StateSchema, after, head}})...)
}

// EntryFrames emits CUSTOM followed by STATE_DELTA. Only the last frame carries
// an SSE id, so that id acknowledges the entire entry pair. Replace patches
// permit safe replay after a truncated response.
func EntryFrames(e Entry) []byte {
	b := frame(custom{"CUSTOM", customNames[0], e.Value})
	b = append(b, []byte("id: "+strconv.FormatInt(e.Index, 10)+"\n")...)
	return append(b, frame(delta{"STATE_DELTA", []patch{{"replace", "/lastIndex", e.Index}, {"replace", "/logHead", e.EntryHash}}})...)
}

// Finished closes a successful run, returning its final resume cursor.
func Finished(run Run, lastIndex int64) []byte {
	return frame(runFinished{"RUN_FINISHED", run.ThreadID, run.RunID, result{lastIndex}})
}

// Error closes a failed run. The boundary supplies a sanitized message and code;
// AG-UI RUN_ERROR has no threadId or runId properties.
func Error(code, msg string) []byte { return frame(runError{"RUN_ERROR", msg, code}) }
