package broker

// Row 184 AC1 + AC6 on the REAL tool binary (WORLD_TOOL_AILANG_BIN, the
// measured v0.52.1; skipped by name when unset, like every real-binary test
// in this package — CI never sets it, row 160). The binary-free companions
// are TestAilangRunWorktreeDeltaOnTheFakeTool and
// TestPolicyToolWorktreeDeltaOnlyForToolchainOps.

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// The trial's program (row 114 Q3, episode tps-a): any std import compiles
// std modules, and v0.52.1's `run --policy` writes that compile cache into
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

// AC6: the trial's case. The cache files the v0.52.1 run leaks into the
// worktree appear in the run's own delta, so a walk can answer "where did
// .ailang/ come from" from the store. A re-run that changes nothing carries
// no delta (AC1's second half on the real binary).
func TestAilangRunWorktreeDeltaRecordsTheV0521CacheLeak(t *testing.T) {
	f := newRealFixture(t, "worktree")
	mustWrite(t, filepath.Join(f.root, "zz_ls.ail"), progZzLs)
	out := f.runRaw(t, map[string]any{"path": "zz_ls.ail", "stdin": ".\n"})
	var res runResult
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatal(err)
	}
	if !res.Admitted || res.ExitCode != 0 || !strings.Contains(res.Stdout, "zz_ls.ail") {
		t.Fatalf("zz_ls run = %+v", res)
	}
	d := res.WorktreeDelta
	if d == nil || d.Basis != WorktreeBasisGit || d.Error != "" {
		t.Fatalf("delta = %+v, want the git-status basis", d)
	}
	cache := 0
	for _, c := range d.Changes {
		if !strings.HasPrefix(c.Path, ".ailang/cache/") {
			t.Fatalf("unexpected change %+v: the run wrote only its compile cache", c)
		}
		if c.Op != "created" || c.Size == nil || len(c.SHA256) != 64 {
			t.Fatalf("cache change %+v lacks created/size/sha256", c)
		}
		cache++
	}
	if _, ok := changeSet(d)[".ailang/cache/compile/manifest.json"]; !ok || cache != d.Total || d.Truncated {
		t.Fatalf("delta lists %d cache files of %d (truncated %v), manifest present %v", cache, d.Total, d.Truncated, ok)
	}
	// Design §1 V5 measured 36 files for this program on this binary.
	if cache != 36 {
		t.Fatalf("delta lists %d .ailang/cache files, V5 measured 36", cache)
	}
	t.Logf("v0.52.1 zz_ls run: %d .ailang/cache files in worktree_delta", cache)

	// A re-run hits the cache and writes nothing: no worktree_delta.
	out = f.runRaw(t, map[string]any{"path": "zz_ls.ail", "stdin": ".\n"})
	if strings.Contains(string(out), "worktree_delta") {
		t.Fatalf("an unchanged re-run carries a delta: %s", out)
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
