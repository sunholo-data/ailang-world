package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestReadCharacterCap(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "readcap")
	setResponse := func(content string) []byte {
		t.Helper()
		raw, _ := json.Marshal(map[string]any{"ok": true, "content": content, "total_lines": 9999})
		if err := os.WriteFile(filepath.Join(tool.dir, "response"), raw, 0600); err != nil {
			t.Fatal(err)
		}
		return raw
	}
	read := func(payload string) map[string]json.RawMessage {
		t.Helper()
		raw, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceRead, Scope: WorkspaceScope}, []byte(payload))
		if err != nil {
			t.Fatal(err)
		}
		var got map[string]json.RawMessage
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatal(err)
		}
		return got
	}
	// AC3: 102,000 characters, including non-ASCII characters whose byte count
	// is larger. Parse every advertised call and actually read that slice.
	line := strings.Repeat("é", 1019) + "\n"
	content := strings.Repeat(line, 100)
	setResponse(content)
	got := read(`{"op":"read","path":"fixture.txt"}`)
	var refused string
	_ = json.Unmarshal(got["refused"], &refused)
	if string(got["ok"]) != "false" || !strings.Contains(refused, "102000 characters over 100 lines; the cap is 40000\n  why:") {
		t.Fatal(refused)
	}
	t.Log(refused)
	re := regexp.MustCompile(`offset: ([0-9]+), limit: ([0-9]+)`)
	chunks := re.FindAllStringSubmatch(refused, -1)
	if len(chunks) != 3 {
		t.Fatalf("chunks: %s", refused)
	}
	covered := 0
	for _, c := range chunks {
		offset, _ := strconv.Atoi(c[1])
		limit, _ := strconv.Atoi(c[2])
		if offset != covered+1 || limit < 1 {
			t.Fatalf("gap/overlap: %s", refused)
		}
		piece := strings.Repeat(line, limit)
		if utf8.RuneCountInString(piece) > 40000 || (covered+limit < 100 && utf8.RuneCountInString(piece+line) <= 40000) {
			t.Fatal("non-maximal or oversized chunk")
		}
		setResponse(piece)
		if result := read(fmt.Sprintf(`{"op":"read","path":"fixture.txt","offset":%d,"limit":%d}`, offset, limit)); string(result["ok"]) != "true" {
			t.Fatalf("chunk refused: %s", result)
		}
		covered += limit
	}
	if covered != 100 {
		t.Fatal(covered)
	}
	setResponse("small\n" + strings.Repeat("x", 50000) + "\n")
	got = read(`{"op":"read","path":"huge","offset":41}`)
	_ = json.Unmarshal(got["refused"], &refused)
	if !strings.Contains(refused, "line 42 cannot be read whole yet") {
		t.Fatal(refused)
	}
	for _, tc := range []struct {
		size int
		ok   string
	}{{6000, "false"}, {4000, "true"}} {
		setResponse(strings.Repeat("x", tc.size))
		got = read(`{"op":"read","path":"explicit","max_chars":5000}`)
		if string(got["ok"]) != tc.ok {
			t.Fatal(got)
		}
		input, err := os.ReadFile(filepath.Join(tool.dir, "stdin"))
		if err != nil || strings.Contains(string(input), "max_chars") {
			t.Fatalf("forwarded cap: %s %v", input, err)
		}
	}
	before := tool.dispatches(t)
	if _, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceRead, Scope: WorkspaceScope}, []byte(`{"op":"read","path":"a","MAX_CHARS":5}`)); err == nil || tool.dispatches(t) != before {
		t.Fatal("alias admitted")
	}
	// AC5: all pre-existing fields and values, including their raw JSON, stay.
	original := setResponse(strings.Repeat("hello\n", 10))
	got = read(`{"op":"read","path":"small"}`)
	var want map[string]json.RawMessage
	_ = json.Unmarshal(original, &want)
	for k, v := range want {
		if string(got[k]) != string(v) {
			t.Fatalf("%s changed: %s != %s", k, got[k], v)
		}
	}
	// The refusal is a successful handler outcome, with one committed read.
	setResponse(content)
	recording := &publishRecordingStore{base: openTestStore(t)}
	session := newSession(recording, "cap-test", []Capability{{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, ExpiresAt: 100, Budget: 1,
	}}, Registry{EffectWorkspaceRead: h}, Live, nil)
	result, ref, err := session.Invoke(boundedTestContext(t), EffectRequest{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1, Now: 1,
	}, []byte(`{"op":"read","path":"fixture.txt"}`))
	if err != nil || ref.IsZero() || !strings.Contains(string(result), "refused:") {
		t.Fatalf("%s %s %v", result, ref, err)
	}
	receipt := effectReceipt(t, recording, 0)
	if receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "succeeded" {
		t.Fatalf("refusal not recorded: %+v", receipt)
	}
	if _, _, err := session.Invoke(boundedTestContext(t), EffectRequest{
		Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1, Now: 1,
	}, []byte(`{"op":"read","path":"fixture.txt"}`)); err == nil {
		t.Fatal("read did not spend its budget")
	}
}
