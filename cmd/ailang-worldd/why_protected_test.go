package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

const protectedFixture = `{"digest":"sha256:patterns","patterns":["*_test.go","tests/**"],"enforced":"dirs-only","base":"123456789abcdef","base_digest":"sha256:tree","base_files":12,"base_ignored_digest":"sha256:ignored","base_ignored_files":3,"intact_before":null,"intact_after":null,"changes":[{"path":"tests/a","status":"untracked","phase":"after","path_truncated":true}],"truncated":true,"total":40,"error":"snapshot failed"}`

func TestWhyProtectedLine(t *testing.T) {
	for _, tc := range []struct {
		name, block, want string
		exec              bool
	}{
		{"intact", `{"intact_before":true,"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: intact (12 files @ 1234567, digest sha256:tree, enforced glob)", true},
		{"changed-before", `{"intact_before":false,"intact_after":true,"changes":[{"path":"a_test.go","status":"modified","phase":"before"}],"total":1}`, "protected: CHANGED before the call — a_test.go (modified)", true},
		{"changed-after", `{"intact_before":true,"intact_after":false,"changes":[{"path":"tests/a","status":"untracked","phase":"after"}],"total":1}`, "protected: CHANGED after the call — tests/a (untracked)", true},
		{"changed-both", `{"intact_before":false,"intact_after":false,"changes":[{"path":"a","status":"deleted","phase":"before"},{"path":"b","status":"ignored","phase":"after"}],"total":2}`, "protected: CHANGED before and after the call — a (deleted), b (ignored)", true},
		{"index-flag", `{"intact_before":false,"changes":[{"path":"a_test.go","status":"index-flag"}],"total":1}`, "protected: CHANGED before the call — a_test.go (index-flag)", true},
		{"head-moved", `{"intact_after":false,"changes":[{"path":"HEAD","status":"head-moved"}],"total":1}`, "protected: CHANGED after the call — HEAD (head-moved)", true},
		{"truncated", `{"intact_after":false,"changes":[{"path":"a","status":"modified"},{"path":"b","status":"deleted"},{"path":"c","status":"untracked"},{"path":"d","status":"ignored"},{"path":"e","status":"modified"},{"path":"f","status":"deleted"}],"truncated":true,"total":40}`, "protected: CHANGED after the call — a (modified), b (deleted), c (untracked), d (ignored), e (modified) (+35 more)", true},
		{"error", protectedFixture, "protected: not checked — snapshot failed", true},
		{"error-beside-changed", `{"intact_after":false,"changes":[{"path":"a","status":"modified"}],"total":1,"error":"before failed"}`, "protected: CHANGED after the call — a (modified); not checked — before failed", true},
		{"error-beside-true", `{"intact_before":true,"intact_after":true,"error":"failed"}`, "protected: not checked — failed", true},
		{"missing-verdict", `{"intact_before":true}`, "protected: not checked — no verdict recorded", true},
		{"missing-before", `{"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: not checked — no verdict recorded", true},
		{"malformed", `"broken"`, "protected: not checked — malformed protected block", true},
		{"malformed-null", `null`, "protected: not checked — malformed protected block", true},
		{"malformed-field", `{"intact_before":"true"}`, "protected: not checked — malformed protected block", true},
		{"none-declared", "", "protected: none declared", true},
		{"non-exec", "", "", false},
		{"old-record", "", "protected: none declared", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload := []byte(`{"content":"old"}`)
			if tc.block != "" {
				payload = []byte(`{"protected":` + tc.block + `}`)
			}
			if got := protectedLine(resultProtected(payload), tc.exec); got != tc.want {
				t.Fatalf("protected line = %q, want %q", got, tc.want)
			}
		})
	}
	t.Run("e2e-exec", func(t *testing.T) {
		r := newWalkRig(t) // Every REST fixture owns its bounded context in this closure.
		for i, tc := range []struct{ effect, result, want string }{
			{broker.EffectWorkspaceExec, `{"protected":` + protectedFixture + `}`, "protected: not checked — snapshot failed"},
			{broker.EffectWorkspaceExec, `{}`, "protected: none declared"},
			{"Workspace.Read", `{}`, ""},
		} {
			inv := r.invokeEffectResult(coordinator.SurfaceA2A, "ep", "tool", fmt.Sprintf("task-%d", i), "", []byte(tc.result), tc.effect)
			code, out, errOut := whyRun(t, r.url, "", inv.id)
			assertChainOK(t, "protected", code, out, errOut, int64(i+1))
			var lines []string
			for _, line := range strings.Split(out, "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "protected:") {
					lines = append(lines, line)
				}
			}
			want := ""
			if tc.want != "" {
				want = "            " + tc.want
			}
			if got := strings.Join(lines, "\n"); got != want {
				t.Fatalf("e2e protected line = %q, want %q", got, want)
			}
			code, out, _ = whyRun(t, r.url, "", inv.id, "--json")
			var ch whyChain
			if code != exitOK || json.Unmarshal([]byte(out), &ch) != nil {
				t.Fatalf("--json exit %d: %s", code, out)
			}
			var got json.RawMessage
			for _, l := range ch.Links {
				if l.Name == "effect" {
					got = l.Protected
				}
			}
			var wantBlock json.RawMessage
			if i == 0 {
				wantBlock = json.RawMessage(protectedFixture)
			}
			assertProtectedBlock(t, got, wantBlock)
		}
	})
}

func assertProtectedBlock(t *testing.T, got, want json.RawMessage) {
	t.Helper()
	var g, w any
	if len(got) != 0 {
		if err := json.Unmarshal(got, &g); err != nil {
			t.Fatal(err)
		}
	}
	if len(want) != 0 {
		if err := json.Unmarshal(want, &w); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("protected JSON = %s, want %s", got, want)
	}
}

func TestWhyProtectedJSON(t *testing.T) {
	raw := resultProtected([]byte(`{"protected":` + protectedFixture + `}`))
	b := mustJSONT(t, whyLink{Name: "effect", OK: true, Detail: "result", Protected: raw})
	var link map[string]json.RawMessage
	if err := json.Unmarshal(b, &link); err != nil {
		t.Fatal(err)
	}
	assertProtectedBlock(t, link["protected"], json.RawMessage(protectedFixture))
	var block map[string]json.RawMessage
	if err := json.Unmarshal(link["protected"], &block); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"intact_before", "intact_after"} {
		if string(block[key]) != "null" {
			t.Fatalf("%s = %s, want null", key, block[key])
		}
	}
	old := whyLink{Name: "effect", OK: true, Ref: "hash", Detail: "result", Protected: resultProtected([]byte(`{"content":"old"}`)), ExecResult: true}
	b = mustJSONT(t, old)
	if want := `{"name":"effect","ok":true,"ref":"hash","detail":"result"}`; string(b) != want {
		t.Fatalf("old link JSON = %s, want %s", b, want)
	}
	var decoded whyLink
	if err := json.Unmarshal(b, &decoded); err != nil || decoded.Protected != nil {
		t.Fatalf("old record: %+v, %v", decoded, err)
	}
	if bytes.Contains(b, []byte(`"protected"`)) {
		t.Fatalf("old link carries protected: %s", b)
	}
}
