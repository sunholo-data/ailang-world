package broker

// Row 184 AC1 + AC6 on the REAL tool binary (WORLD_TOOL_AILANG_BIN, the
// manifest's tool release; skipped by name when unset, like every real-binary
// test in this package — CI never sets it, row 160). The binary-free
// companions are TestAilangRunWorktreeDeltaOnTheFakeTool and
// TestPolicyToolWorktreeDeltaOnlyForToolchainOps.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The trial's program (row 114 Q3, episode tps-a): any std import compiles
// std modules, and v0.52.1's `run --policy` wrote that compile cache into
// the sandbox root (upstream sunholo-data/ailang#1547, design §1 V5).
const progZzLs = "module zz_ls\n\nimport std/io (println, readLine)\nimport std/fs (listDir)\nimport std/string (trim, join)\n\n" +
	"export func main() -> () ! {IO, FS} {\n  let path = trim(readLine(()));\n  println(join(\"\\n\", listDir(path)))\n}\n"

func (f *realFixture) runRaw(t *testing.T, payload map[string]any) []byte {
	t.Helper()
	body, _ := json.Marshal(payload)
	out, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatalf("run %s: %v", body, err)
	}
	return out
}

// AC6: the trial's case, now fixed upstream. On v0.52.1 this run leaked its
// compile cache into the worktree: design §1 V5 measured 36 files under
// .ailang/cache/ (compile/manifest.json among them), each recorded `created`
// with size and sha256 in the run's own delta, which is how a walk answered
// "where did .ailang/ come from" (ailang#1547). v0.54.15, the tool release
// since row 139 D7, writes no cache into the sandbox: the same run changes
// nothing in the worktree, so it carries no worktree_delta, and no .ailang/
// directory appears. A re-run changes nothing either (AC1's second half).
// The delta itself is still live on this binary:
// TestAilangRunWorktreeDeltaRecordsAWrittenFile is the positive control.
func TestAilangRunLeavesNoCacheInTheWorktreeDelta(t *testing.T) {
	f := newRealFixture(t, "worktree")
	mustWrite(t, filepath.Join(f.root, "zz_ls.ail"), progZzLs)
	for _, run := range []string{"first", "re-run"} {
		out := f.runRaw(t, map[string]any{"path": "zz_ls.ail", "stdin": ".\n"})
		var res runResult
		if err := json.Unmarshal(out, &res); err != nil {
			t.Fatal(err)
		}
		if !res.Admitted || res.ExitCode != 0 || !strings.Contains(res.Stdout, "zz_ls.ail") {
			t.Fatalf("%s zz_ls run = %+v", run, res)
		}
		if d := res.WorktreeDelta; d != nil || strings.Contains(string(out), "worktree_delta") {
			t.Fatalf("%s zz_ls run carries a worktree delta (the ailang#1547 cache leak is back?): %+v", run, d)
		}
		if _, err := os.Lstat(filepath.Join(f.root, ".ailang")); !os.IsNotExist(err) {
			t.Fatalf("%s zz_ls run: .ailang/ exists in the worktree (Lstat err %v)", run, err)
		}
	}
	f.assertProtectedUnchanged(t)
}

// AC1 on the real binary: a program that writes a file records it as created
// with size and sha256.
func TestAilangRunWorktreeDeltaRecordsAWrittenFile(t *testing.T) {
	f := newRealFixture(t, "init")
	res := f.run(t, f.h, "wr.ail", `"made.txt"`)
	if !res.Admitted || res.ExitCode != 0 || res.WorktreeDelta == nil {
		t.Fatalf("wr run = %+v", res)
	}
	assertChange(t, res.WorktreeDelta, "made.txt", "created", "pwned")
	f.assertProtectedUnchanged(t)
}
