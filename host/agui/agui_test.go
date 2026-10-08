package agui

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The fixture has a deliberate gap. Values are frozen log-route-shaped JSON;
// the golden is an independent, checked-in wire specification, not regenerated
// by the encoder under test.
func fixture() []Entry {
	var out []Entry
	for _, i := range []int64{0, 1, 2, 5, 6} {
		h := fmt.Sprintf("sha256:%064x", i+1)
		v := fmt.Sprintf(`{"header":{"entryIndex":%d,"semanticsEpoch":1,"transitionFn":"sha256:%064x","interpreter":"sha256:%064x","prevEntryHash":"sha256:%064x","writtenBy":"fixture <>&"},"entryHash":"%s","transitionRef":"sha256:%064x"}`, i, 100, 101, i+100, h, i+200)
		out = append(out, Entry{Index: i, EntryHash: h, Value: json.RawMessage(v)})
	}
	return out
}

var fixtureRun = Run{ThreadID: "t", RunID: "r"}

func prefix() []byte {
	b := Start(fixtureRun, -1, nil)
	for _, e := range fixture() {
		b = append(b, EntryFrames(e)...)
	}
	return b
}
func allFrames() []byte {
	b := prefix()
	b = append(b, Finished(fixtureRun, 6)...)
	return append(b, Error("Internal", "store read failed")...)
}
func readGolden(t *testing.T) []byte {
	t.Helper()
	b, e := os.ReadFile("testdata/stream_fixture.golden")
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestGolden(t *testing.T) {
	if got := prefix(); !bytes.Equal(got, readGolden(t)) {
		t.Fatalf("golden bytes differ:\n%s", got)
	}
	if got := string(Finished(fixtureRun, 6)); got != "data: {\"type\":\"RUN_FINISHED\",\"threadId\":\"t\",\"runId\":\"r\",\"result\":{\"lastIndex\":6}}\n\n" {
		t.Fatalf("finished bytes: %q", got)
	}
	if got := string(Error("Internal", "store read failed")); got != "data: {\"type\":\"RUN_ERROR\",\"message\":\"store read failed\",\"code\":\"Internal\"}\n\n" {
		t.Fatalf("error bytes: %q", got)
	}
	if ProtocolVersion != "1.0" || StateSchema != "world/agui-state/v1" || !reflect.DeepEqual(customNames, []string{"world.entry.committed"}) {
		t.Fatal("protocol/schema/closed CUSTOM table changed")
	}
	// A resumed snapshot describes the cursor entry, including the null genesis arm.
	e := fixture()[3]
	want := fmt.Sprintf("data: {\"type\":\"RUN_STARTED\",\"threadId\":\"t\",\"runId\":\"r\",\"protocolVersion\":\"1.0\"}\n\ndata: {\"type\":\"STATE_SNAPSHOT\",\"snapshot\":{\"schema\":\"world/agui-state/v1\",\"lastIndex\":5,\"logHead\":\"%s\"}}\n\n", e.EntryHash)
	if string(Start(fixtureRun, 5, &e)) != want {
		t.Fatal("cursor snapshot bytes differ")
	}
}
func TestDeterministic(t *testing.T) {
	if !bytes.Equal(allFrames(), allFrames()) {
		t.Fatal("two encodes differ")
	}
	for _, ev := range stockEvents(t, allFrames()) {
		if _, ok := ev["timestamp"]; ok {
			t.Fatal("timestamp must be absent")
		}
	}
	if !bytes.Equal(prefix(), readGolden(t)) {
		t.Fatal("deterministic golden bytes differ")
	}
}

// Exact stock-client LF LF splitting, data lines only; incomplete tail dropped.
func stockEvents(t *testing.T, b []byte) []map[string]any {
	t.Helper()
	var events []map[string]any
	parts := strings.Split(string(b), "\n\n")
	for _, frame := range parts[:len(parts)-1] {
		var lines []string
		for _, line := range strings.Split(frame, "\n") {
			if strings.HasPrefix(line, "data: ") {
				lines = append(lines, strings.TrimPrefix(line, "data: "))
			}
		}
		var ev map[string]any
		if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &ev); err != nil {
			t.Fatal(err)
		}
		events = append(events, ev)
	}
	return events
}
func TestFramesSplitLikeStockClient(t *testing.T) {
	b := allFrames()
	if bytes.ContainsRune(b, '\r') {
		t.Fatal("CR forbidden: stock client uses LF LF")
	}
	expected := append(append([]byte(nil), readGolden(t)...), []byte("data: {\"type\":\"RUN_FINISHED\",\"threadId\":\"t\",\"runId\":\"r\",\"result\":{\"lastIndex\":6}}\n\ndata: {\"type\":\"RUN_ERROR\",\"message\":\"store read failed\",\"code\":\"Internal\"}\n\n")...)
	got, want := stockEvents(t, b), stockEvents(t, expected)
	if len(got) != 14 || !reflect.DeepEqual(got, want) {
		t.Fatalf("stock event list differs: %d", len(got))
	}
	if truncated := stockEvents(t, b[:len(b)-1]); !reflect.DeepEqual(truncated, want[:len(want)-1]) {
		t.Fatal("unterminated tail was not dropped")
	}
}

// Pinned AG-UI 1.0 schema: MIT, ag-ui-protocol/ag-ui commit
// 903a9ab9a154a164e61c216d8206a760f44780dc, spec/1.0/schema.json.
// This checks the emitted subset, not a general JSON Schema validator.
func TestEventsConformToPinnedSchema(t *testing.T) {
	b, err := os.ReadFile("testdata/agui-1.0-schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(b)) != "4b5c93226838a0e72d88e6c5df20633c686c49fcb75d9be2815c6cbf9e48e71a" {
		t.Fatal("schema SHA-256 differs")
	}
	var schema struct {
		Defs map[string]struct {
			Properties map[string]struct {
				Const string `json:"const"`
			}
			Required []string
			Enum     []string
			AllOf    []struct {
				Ref string `json:"$ref"`
			}
		} `json:"$defs"`
	}
	if err = json.Unmarshal(b, &schema); err != nil {
		t.Fatal(err)
	}
	names := map[string]string{"RUN_STARTED": "RunStartedEvent", "RUN_FINISHED": "RunFinishedEvent", "RUN_ERROR": "RunErrorEvent", "STATE_SNAPSHOT": "StateSnapshotEvent", "STATE_DELTA": "StateDeltaEvent", "CUSTOM": "CustomEvent"}
	seen := map[string]int{}
	for _, ev := range stockEvents(t, allFrames()) {
		typ, ok := ev["type"].(string)
		if !ok {
			t.Fatal("missing type")
		}
		name, ok := names[typ]
		if !ok {
			t.Fatalf("unexpected type %s", typ)
		}
		def, ok := schema.Defs[name]
		if !ok || def.Properties["type"].Const != typ {
			t.Fatalf("missing/mismatched def %s", name)
		}
		enumOK := false
		for _, v := range schema.Defs["EventType"].Enum {
			enumOK = enumOK || v == typ
		}
		if !enumOK {
			t.Fatal("type absent from enum")
		}
		allowed := map[string]bool{}
		for k := range def.Properties {
			allowed[k] = true
		}
		required := append([]string(nil), def.Required...)
		for _, ref := range def.AllOf {
			base := schema.Defs[strings.TrimPrefix(ref.Ref, "#/$defs/")]
			for k := range base.Properties {
				allowed[k] = true
			}
			required = append(required, base.Required...)
		}
		for _, k := range required {
			if _, ok := ev[k]; !ok {
				t.Fatalf("%s missing required %s", typ, k)
			}
		}
		for k := range ev {
			if !allowed[k] {
				t.Fatalf("%s disallowed property %s", typ, k)
			}
		}
		if typ == "STATE_DELTA" {
			ops := ev["delta"].([]any)
			if len(ops) != 2 {
				t.Fatal("want two delta ops")
			}
			for _, op := range ops {
				m := op.(map[string]any)
				for _, k := range []string{"op", "path", "value"} {
					if _, ok := m[k]; !ok {
						t.Fatalf("delta missing %s", k)
					}
				}
				if m["op"] != "replace" || (m["path"] != "/lastIndex" && m["path"] != "/logHead") {
					t.Fatal("unexpected delta op/path")
				}
			}
		}
		seen[typ]++
	}
	if len(seen) != 6 {
		t.Fatalf("checked %d event types, want six", len(seen))
	}
}
func TestDeltaIdempotent(t *testing.T) {
	for _, e := range fixture() {
		evs := stockEvents(t, EntryFrames(e))
		if len(evs) != 2 {
			t.Fatal("want entry pair")
		}
		state := map[string]any{"lastIndex": float64(e.Index), "logHead": e.EntryHash}
		want := map[string]any{"lastIndex": float64(e.Index), "logHead": e.EntryHash}
		for n := 0; n < 2; n++ {
			for _, v := range evs[1]["delta"].([]any) {
				op := v.(map[string]any)
				if op["op"] != "replace" {
					t.Fatal("every delta op must be replace")
				}
				path, ok := op["path"].(string)
				if !ok {
					t.Fatal("path not string")
				}
				key := strings.TrimPrefix(path, "/")
				if _, ok := state[key]; !ok {
					t.Fatalf("unknown replace path %s", path)
				}
				state[key] = op["value"]
			}
			if !reflect.DeepEqual(state, want) {
				t.Fatalf("delta application %d changed post-entry state: %v", n, state)
			}
		}
	}
}
func TestAGUIDependencyBoundary(t *testing.T) {
	// Go's cache and any tool scratch writes must stay in this test's TempDir.
	dir := t.TempDir()
	for _, args := range [][]string{{"list", "-deps", "./host/agui"}, {"vet", "./host/agui"}} {
		cmd := exec.Command("go", args...)
		cmd.Dir = "../.."
		cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(dir, "cache"), "GOTMPDIR="+dir, "TMPDIR="+dir, "GOPROXY=off", "GOTOOLCHAIN=local", "GOTELEMETRY=off", "GOENV=off", "HOME="+dir)
		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("go %v: %v\n%s", args, err, b)
		}
		for _, dep := range strings.Fields(string(b)) {
			if dep == "net/http" || strings.HasSuffix(dep, "/host/store") {
				t.Fatalf("forbidden dependency %s", dep)
			}
		}
	}
}
