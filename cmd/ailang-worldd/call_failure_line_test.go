package main

// Row 255 (design_docs/planned/w-exec-failure-carries-a-reason.md §5, TR):
// `ailang-worldd call` prints a failed effect's failure line beside its
// record, and an effect with no line prints as before.

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestCallPrintsTheFailureLine(t *testing.T) {
	const line = "seed-cache: GOCACHE: timed out after 9s"
	for _, tc := range []struct{ text, want string }{
		{`{"world":{"effects":[{"id":"e1","reason":"` + line + `","record":"sha256:1111","status":"failed"}],"plan":"sha256:2222"}}`,
			"  effect     e1 failed sha256:1111 — " + line + "\n"},
		{`{"world":{"effects":[{"id":"e1","record":"sha256:1111","status":"failed"}],"plan":"sha256:2222"}}`,
			"  effect     e1 failed sha256:1111\n"},
	} {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(tc.text), &fields); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		printCallResult(&out, tc.text, fields)
		if !strings.Contains(out.String(), tc.want) {
			t.Fatalf("call printed:\n%s\nwant the line %q", out.String(), tc.want)
		}
	}
}
