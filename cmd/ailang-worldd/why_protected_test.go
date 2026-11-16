package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

const protectedFixture = `{"digest":"sha256:patterns","patterns":["*_test.go","tests/**"],"enforced":"dirs-only","base":"123456789abcdef","base_digest":"sha256:tree","base_files":12,"base_ignored_digest":"sha256:ignored","base_ignored_files":3,"intact_before":null,"intact_after":null,"changes":[{"path":"tests/a","status":"untracked","phase":"after","path_truncated":true}],"truncated":true,"total":40,"error":"snapshot failed"}`

// The protected line end to end (its pure cases are host/why's TestWhyProtectedLine).
func TestWhyProtectedLine(t *testing.T) {
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
	t.Run("e2e-tampered", func(t *testing.T) {
		// A result whose bytes fail verification is never decoded: a forged
		// intact block on a Workspace.Exec result prints no protected line
		// and carries no protected key, while the effect link is ✗.
		r := newWalkRig(t)
		inv := r.invokeEffectResult(coordinator.SurfaceA2A, "ep", "tool", "task-0", "", []byte(`{"protected":`+protectedFixture+`}`), broker.EffectWorkspaceExec)
		forged := `{"protected":{"intact_before":true,"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}}`
		addr := tamperProxy(t, r.url, inv.result, []byte(forged))
		code, out, _ := whyRun(t, addr, "", inv.id)
		if code != exitIntegrity || !strings.Contains(out, "✗ effect") {
			t.Fatalf("tampered exec result: exit %d\n%s", code, out)
		}
		if strings.Contains(out, "protected:") {
			t.Fatalf("a tampered result printed a protected line:\n%s", out)
		}
		code, out, _ = whyRun(t, addr, "", inv.id, "--json")
		var ch whyChain
		if code != exitIntegrity || json.Unmarshal([]byte(out), &ch) != nil {
			t.Fatalf("--json over a tampered result: exit %d: %s", code, out)
		}
		for _, l := range ch.Links {
			if l.Protected != nil {
				t.Fatalf("a tampered result carried a protected block: %s", l.Protected)
			}
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
