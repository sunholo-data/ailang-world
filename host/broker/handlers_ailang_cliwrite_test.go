package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Row 134 V67 (independent judge finding P2): Ailang.CLI never writes the
// worktree. On v0.51.0 policy-tool runs `fmt --write` under the rendered
// episode policy, and the formatter rewrites the file without the
// fs_deny_write check. These tests pin the handler refusal (binary-independent
// and on the real tool binary) and the audited per-op flag table. The AC4.1
// matrix (handlers_ailang_test.go) calls assertCLIFmtWriteRefused.

// TestAilangCLIRefusesWriteFlagsWithoutDispatch is the binary-independent
// half of V67 (MUT-CLI-WRITE-FLAG): a write-capable flag never reaches
// policy-tool, under any value, aliased key or duplicated flags object; the
// read-only fmt forms still dispatch, byte for byte.
func TestAilangCLIRefusesWriteFlagsWithoutDispatch(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "ok")
	ctx := boundedTestContext(t)
	for _, payload := range []string{
		`{"op":"fmt","path":"a.ail","flags":{"write":""}}`,
		`{"op":"fmt","path":"a.ail","flags":{"write":"false"}}`,
		`{"op":"fmt","path":"a.ail","flags":{"check":"","write":""}}`,
		`{"op":"fmt","path":"a.ail","flags":{"write":""},"flags":{"write":""}}`,
		`{"op":"fmt","path":"a.ail","FLAGS":{"write":""}}`,
		`{"op":"fmt","path":"a.ail","Flags":{}}`,
		"{\"op\":\"fmt\",\"path\":\"a.ail\",\"flagſ\":{\"write\":\"\"}}",
		`{"op":"fmt","path":"a.ail","flags":"write"}`,
	} {
		before := tool.dispatches(t)
		_, err := h.Execute(ctx, EffectRequest{Effect: EffectAilangCLI, Scope: WorkspaceScope, Cost: 1}, []byte(payload))
		var refusal *AilangToolRefusalError
		if !errors.As(err, &refusal) || tool.dispatches(t) != before {
			t.Errorf("%s: err=%v dispatched=%d, want a refusal and no dispatch", payload, err, tool.dispatches(t)-before)
		}
	}
	for _, payload := range []string{
		`{"flags":{"check":""},"op":"fmt","path":"a.ail"}`,
		`{"op":"fmt","path":"a.ail"}`,
		`{"flags":null,"op":"fmt","path":"a.ail"}`,
		`{"flags":{"write":""},"op":"check","path":"a.ail"}`, // not fmt: policy-tool's own flag refusal answers
	} {
		before := tool.dispatches(t)
		if _, err := h.Execute(ctx, EffectRequest{Effect: EffectAilangCLI, Scope: WorkspaceScope, Cost: 1}, []byte(payload)); err != nil || tool.dispatches(t) != before+1 {
			t.Errorf("%s: err=%v, want one dispatch", payload, err)
		}
		if got, _ := os.ReadFile(filepath.Join(tool.dir, "stdin")); string(got) != payload {
			t.Errorf("%s: policy-tool read %s", payload, got)
		}
	}
}

// unformattedAil is a valid module that `ailang fmt` rewrites (the extra
// spaces after `func`), so a write would change its bytes.
const unformattedAil = "module x\nexport func   main() -> int { 1 }\n"

// assertCLIFmtWriteRefused is the AC4.1 fmt rows: every payload that would
// make policy-tool run `fmt --write` is refused by the handler with no
// dispatch, and the target's bytes are unchanged. It returns the row count.
func (f *realFixture) assertCLIFmtWriteRefused(t *testing.T) int {
	t.Helper()
	rows := 0
	for _, target := range []string{".claude/fmt.ail", ".ailang/fmt.ail", "plain_fmt.ail"} {
		mustMkdir(t, filepath.Dir(filepath.Join(f.root, target)))
		mustWrite(t, filepath.Join(f.root, target), unformattedAil)
		for _, payload := range []string{
			fmt.Sprintf(`{"op":"fmt","path":%q,"flags":{"write":""}}`, target),
			fmt.Sprintf(`{"op":"fmt","path":%q,"flags":{"write":"false"}}`, target),
			fmt.Sprintf(`{"op":"fmt","path":%q,"FLAGS":{"write":""}}`, target),
			fmt.Sprintf(`{"op":"fmt","path":%q,"flags":{},"flags":{"write":""}}`, target),
		} {
			_, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangCLI, Scope: WorkspaceScope, Cost: 1}, []byte(payload))
			var refusal *AilangToolRefusalError
			if !errors.As(err, &refusal) {
				t.Errorf("Ailang.CLI %s: err = %v, want *AilangToolRefusalError", payload, err)
			}
			if got, _ := f.snapshot(target); got != unformattedAil {
				t.Fatalf("Ailang.CLI %s rewrote %s: %q", payload, target, got)
			}
			rows++
		}
		// policy-tool MERGES duplicate "flags" objects (measured: a write in
		// the first survives a trailing {}), but the handler forwards only
		// its own re-encode, where the last duplicate wins: this payload
		// reaches the binary as plain fmt, and the file is unchanged.
		dup := fmt.Sprintf(`{"op":"fmt","path":%q,"flags":{"write":""},"flags":{}}`, target)
		if out, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangCLI, Scope: WorkspaceScope, Cost: 1}, []byte(dup)); err != nil ||
			!strings.Contains(string(out), `"argv":["fmt",`+strconv.Quote(target)+`]`) {
			t.Fatalf("Ailang.CLI %s: err=%v out=%s, want the collapsed plain fmt", dup, err, out)
		}
		if got, _ := f.snapshot(target); got != unformattedAil {
			t.Fatalf("Ailang.CLI %s rewrote %s: %q", dup, target, got)
		}
		rows++
		// The read-only forms stay admitted: plain fmt returns the text,
		// fmt --check reports; neither touches the file.
		if resp := f.op(t, EffectAilangCLI, map[string]any{"op": "fmt", "path": target}); resp["ok"] != true ||
			!strings.Contains(fmt.Sprint(resp["stdout"]), "export func main() -> int {") {
			t.Fatalf("Ailang.CLI fmt %s = %v, want the formatted text", target, resp)
		}
		if resp := f.op(t, EffectAilangCLI, map[string]any{"op": "fmt", "path": target, "flags": map[string]string{"check": ""}}); resp["ok"] != false || resp["exit_code"] != float64(1) {
			t.Fatalf("Ailang.CLI fmt --check %s = %v, want rc 1 (needs formatting)", target, resp)
		}
		if got, _ := f.snapshot(target); got != unformattedAil {
			t.Fatalf("a read-only fmt rewrote %s: %q", target, got)
		}
		rows += 2
	}
	f.assertProtectedUnchanged(t)
	return rows
}

// measuredCLIFlags is each `cli` op's admitted flag set as the tool binary's
// policy-tool states it in the refusal for an unknown flag (V67 on v0.51.0,
// re-baselined on v0.52.1 by row-134 design §13 V78: builtins_list gained the
// read-only filters --by-effect --by-module --module --query --verbose, and
// examples_list's --tag became --tags; a write watch over the worktree and the
// cache saw no file change for any of them). Of these, only fmt's --write
// writes the worktree (cliWriteFlags); the table is a tripwire: a release that
// admits a new flag reds here before it can reach an agent unaudited.
var measuredCLIFlags = map[string]string{
	"agent_prompt": "none", "ai_check": "--timeout", "axioms": "none", "builtins_list": "--by-effect --by-module --json --module --query --verbose",
	"builtins_show": "none", "check": "--json --quiet --strict-syntax", "devtools_prompt": "none",
	"docs_search": "--json --limit", "examples_list": "--status --tags", "examples_search": "none",
	"examples_show": "none", "examples_tags": "none", "fmt": "--check --write", "iface": "--compact",
	"pkg_docs": "none", "policy_check": "none", "prompt": "none",
	"test": "--allow-skips --json --no-color --package", "tree": "none", "version": "none",
}

// TestAilangCLIAdmittedFlagsAreTheAuditedSet is the V67 audit on the real
// binary: the cli op list and every op's admitted flags equal the measured
// table, and the handler refuses exactly the write-capable ones.
func TestAilangCLIAdmittedFlagsAreTheAuditedSet(t *testing.T) {
	f := newRealFixture(t, "init")
	if len(f.h.cliOps) != len(measuredCLIFlags) {
		t.Fatalf("cli op list has %d ops, the audit covers %d", len(f.h.cliOps), len(measuredCLIFlags))
	}
	for op, admitted := range measuredCLIFlags {
		if !f.h.cliOps[op] {
			t.Fatalf("audited op %q is not in the binary's cli list", op)
		}
		out, err := f.h.policyTool(boundedTestContext(t), []byte(fmt.Sprintf(`{"op":%q,"path":"data.txt","flags":{"zzbogus":""}}`, op)))
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("op %s does not admit flag --zzbogus (admitted: %s)", op, admitted)
		var resp map[string]any
		if json.Unmarshal(out, &resp) != nil || resp["refused"] != want {
			t.Errorf("%s: policy-tool said %s, want refused %q", op, out, want)
		}
		for _, flag := range strings.Fields(admitted) {
			name := strings.TrimPrefix(flag, "--")
			if flag == "none" {
				continue
			}
			writes := false
			for _, w := range cliWriteFlags[op] {
				writes = writes || w == name
			}
			if writes != (op == "fmt" && name == "write") {
				t.Errorf("%s --%s: write-capable=%v disagrees with the V67 audit", op, name, writes)
			}
		}
	}
}
