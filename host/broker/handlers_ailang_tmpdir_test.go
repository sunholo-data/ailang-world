package broker

// Row 197 (w-toolchain-tmpdir-in-episode-cache): every toolchain child
// (policy-tool, `run --policy`, the CLI ops) gets TMPDIR=<cache>/tmp, made
// 0700 before each launch. Measured on v0.54.15 with no TMPDIR: ai_check and
// the CLI ops that spawn a child first make a private cache dir with
// mkdir /tmp/ailang-policy-tool-cache-<n> (removed when they finish) —
// outside the episode, recorded nowhere; inside the exec sandbox, where /tmp
// is not writable, the call fails ("policy-tool: creating the child's
// private cache dir: mkdir /tmp/ailang-policy-tool-cache-…: operation not
// permitted").

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestToolchainChildTMPDIRIsTheEpisodeCache (fake tool, always runs): the
// summary at construction, a policy-tool op and a run each see TMPDIR =
// <cache>/tmp at 0700, recreated when the cache was cleaned between calls;
// a symlink in its place is refused before anything launches.
func TestToolchainChildTMPDIRIsTheEpisodeCache(t *testing.T) {
	h, tool, l := newFakeToolHandler(t, "admitted")
	tmp := filepath.Join(l.cacheDir, "tmp")
	check := func(what string) {
		t.Helper()
		if got := string(mustRead(t, filepath.Join(tool.dir, "tmpdir"))); got != tmp {
			t.Fatalf("%s: the child's TMPDIR = %q, want %q", what, got, tmp)
		}
		if info, err := os.Lstat(tmp); err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Fatalf("%s: %s = %v (%v), want a 0700 directory", what, tmp, info, err)
		}
	}
	check("construction summary")
	for _, c := range []struct{ effect, payload string }{
		{EffectWorkspaceRead, `{"op":"read","path":"a.txt"}`},
		{EffectAilangRun, `{"path":"m.ail"}`},
	} {
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatal(err)
		}
		if _, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: c.effect, Scope: WorkspaceScope, Cost: 1}, []byte(c.payload)); err != nil {
			t.Fatalf("%s: %v", c.effect, err)
		}
		check(c.effect)
	}
	if err := os.RemoveAll(tmp); err != nil {
		t.Fatal(err)
	}
	elsewhere := filepath.Join(l.base, "elsewhere")
	mustMkdir(t, elsewhere)
	if err := os.Symlink(elsewhere, tmp); err != nil {
		t.Fatal(err)
	}
	before := tool.dispatches(t)
	_, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"op":"read","path":"a.txt"}`))
	if err == nil || !strings.Contains(err.Error(), "is not a directory") || tool.dispatches(t) != before {
		t.Fatalf("symlinked TMPDIR: err=%v dispatched=%d, want a refusal before any launch", err, tool.dispatches(t)-before)
	}
}

// snapshotOutside fingerprints everything under the fixture base except the
// worktree and the episode cache: the sibling episode, the policy directory
// and the cache root itself.
func snapshotOutside(t *testing.T, f *realFixture) map[string]string {
	t.Helper()
	snap := map[string]string{}
	err := filepath.WalkDir(f.base, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == f.root || p == f.cacheDir {
			return filepath.SkipDir
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(f.base, p)
		snap[rel] = fmt.Sprintf("%v %d %d", info.Mode(), info.Size(), info.ModTime().UnixNano())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snap
}

// policyToolCacheDirsInTmp lists the private cache dirs v0.54.15 names
// under /tmp; ok is false when /tmp cannot be listed (that half is skipped).
func policyToolCacheDirsInTmp() (map[string]bool, bool) {
	ents, err := os.ReadDir("/tmp")
	if err != nil {
		return nil, false
	}
	out := map[string]bool{}
	for _, e := range ents {
		if strings.HasPrefix(e.Name(), "ailang-policy-tool-cache-") {
			out[e.Name()] = true
		}
	}
	return out, true
}

// TestToolchainCallsWriteOnlyInTheWorktreeAndEpisodeCache (real tool binary)
// is row 197's proof: (1) ai_check and a CLI op succeed through the handler —
// inside the exec sandbox they fail with no TMPDIR (measured above); (2) the
// binary itself names where it puts the private cache dir: with <cache>/tmp
// replaced by a file, the same child under the handler's exact env reports
// `mkdir <cache>/tmp/ailang-policy-tool-cache-…`; (3) after reads, a write, a
// delete, the toolchain ops and a run, nothing outside {worktree, episode
// cache} under the fixture changed, and no private cache dir is left in /tmp.
func TestToolchainCallsWriteOnlyInTheWorktreeAndEpisodeCache(t *testing.T) {
	f := newRealFixture(t, "init")
	tmp := filepath.Join(f.cacheDir, "tmp")
	outside := snapshotOutside(t, f)
	tmpBefore, listable := policyToolCacheDirsInTmp()

	for _, c := range []struct {
		effect  string
		payload map[string]any
	}{
		{EffectAilangCheck, map[string]any{"op": "ai_check", "path": "readit.ail"}},
		{EffectAilangCLI, map[string]any{"op": "check", "path": "readit.ail"}},
		{EffectWorkspaceRead, map[string]any{"op": "read", "path": "data.txt", "offset": 1, "limit": 1}},
		{EffectWorkspaceWrite, map[string]any{"op": "write", "path": "made.txt", "content": "x\n"}},
		{EffectWorkspaceDelete, map[string]any{"op": "delete", "path": "made.txt"}},
	} {
		if got := f.op(t, c.effect, c.payload); got["ok"] != true {
			t.Fatalf("%s %v = %v, want ok (a TMPDIR outside the sandbox's writable set fails here)", c.effect, c.payload, got)
		}
	}
	if res := f.run(t, f.h, "readit.ail", ""); !res.Admitted || res.ExitCode != 0 || res.Stdout != "hello from ep1\n\n" {
		t.Fatalf("run readit.ail = %+v", res)
	}

	if err := os.RemoveAll(tmp); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, tmp, "not a directory\n")
	stderr := newBoundedSink(f.h.bounds.maxOutputBytes)
	out, err := runBounded(boundedTestContext(t), f.h.bounds, handlerCommand{
		path: f.bin, args: []string{"policy-tool", "--policy", f.policyPath}, dir: f.root, env: f.h.childEnv(),
		stdin: []byte(`{"op":"ai_check","path":"readit.ail"}`), stderr: stderr,
	})
	if want := "mkdir " + filepath.Join(tmp, "ailang-policy-tool-cache-"); err != nil || !strings.Contains(string(out), want) {
		t.Fatalf("ai_check with <cache>/tmp a file = %s (err %v, stderr %q); want the binary to name %s…", out, err, stderr.Bytes(), want)
	}
	if err := os.Remove(tmp); err != nil {
		t.Fatal(err)
	}
	if got := f.op(t, EffectAilangCheck, map[string]any{"op": "ai_check", "path": "readit.ail"}); got["ok"] != true {
		t.Fatalf("ai_check after the cache was cleaned = %v, want ok (the handler remakes <cache>/tmp)", got)
	}

	if after := snapshotOutside(t, f); !reflect.DeepEqual(after, outside) {
		t.Fatalf("toolchain calls changed files outside the worktree and episode cache:\nbefore %v\nafter  %v", outside, after)
	}
	f.assertProtectedUnchanged(t)
	if !listable {
		t.Log("/tmp is not listable here; the leftover check is skipped (the in-sandbox ok calls above carry the proof)")
		return
	}
	tmpAfter, _ := policyToolCacheDirsInTmp()
	for name := range tmpAfter {
		if !tmpBefore[name] {
			t.Fatalf("/tmp/%s appeared during the toolchain calls", name)
		}
	}
}
