package broker

// Row 191a (w-workspace-delete): Workspace.Delete is policy-tool's `delete`
// op (v0.54.15, ailang#1875) behind its own effect name. The op allowlist
// matrix (TestAilangToolOpAllowlistIsFixedPerEffect) holds delete to
// Workspace.Delete alone; this file is the confinement half on the real tool
// binary: the protected paths, the escapes and the non-files refuse, and a
// refused delete leaves the bytes where they were.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func (f *realFixture) deleteOp(t *testing.T, path string) map[string]any {
	t.Helper()
	return f.op(t, EffectWorkspaceDelete, map[string]any{"op": "delete", "path": path})
}

// TestAilangDeleteConfinementOnTheToolBinary runs for both shipped
// topologies: every protected target of the write matrix (the git metadata
// and each fs_deny_write glob) plus case-folded spellings refuse with the
// file unchanged; `..`, absolute and symlink-target escapes refuse; a
// symlink is removed itself, never its target; a directory and a missing
// file refuse; an ordinary file is removed.
func TestAilangDeleteConfinementOnTheToolBinary(t *testing.T) {
	for _, topology := range []string{"init", "worktree"} {
		t.Run(topology, func(t *testing.T) {
			f := newRealFixture(t, topology)
			targets := append(f.protectedTargets(), ".GITHUB/workflows/y.yml", ".Claude/settings.local.json", ".GITMODULES")
			seeded := map[string][]byte{}
			for _, rel := range targets {
				p := filepath.Join(f.root, rel)
				if data, err := os.ReadFile(p); err == nil {
					seeded[rel] = data // the git metadata the fixture made
					continue
				}
				if strings.HasPrefix(strings.ToLower(rel), ".git/") || strings.EqualFold(rel, ".git") {
					continue // a .git path the topology lacks: refused all the same
				}
				mustMkdir(t, filepath.Dir(p))
				mustWrite(t, p, "protected "+rel+"\n")
				seeded[rel] = mustRead(t, p)
			}
			for _, rel := range targets {
				got := f.deleteOp(t, rel)
				if got["ok"] != false || fmt.Sprint(got["refused"]) == "" {
					t.Fatalf("delete %s = %v, want refused", rel, got)
				}
				if want, ok := seeded[rel]; ok {
					if now, err := os.ReadFile(filepath.Join(f.root, rel)); err != nil || !bytes.Equal(now, want) {
						t.Fatalf("refused delete %s changed the file: %q, %v", rel, now, err)
					}
				}
			}
			f.assertProtectedUnchanged(t)

			for _, out := range []string{"../ep2/secret.txt", filepath.Join(f.sibling, "secret.txt")} {
				if got := f.deleteOp(t, out); got["ok"] != false || !strings.Contains(fmt.Sprint(got["refused"]), "escapes sandbox") {
					t.Fatalf("delete %s = %v, want refused escapes sandbox", out, got)
				}
			}
			// link.txt -> ../ep2/secret.txt: the link goes, the target stays.
			if got := f.deleteOp(t, "link.txt"); got["ok"] != true {
				t.Fatalf("delete link.txt = %v", got)
			}
			if _, err := os.Lstat(filepath.Join(f.root, "link.txt")); !os.IsNotExist(err) {
				t.Fatalf("link.txt still present: %v", err)
			}
			if data := mustRead(t, filepath.Join(f.sibling, "secret.txt")); string(data) != "secret\n" {
				t.Fatalf("the symlink's target changed: %q", data)
			}

			mustMkdir(t, filepath.Join(f.root, "dir"))
			mustWrite(t, filepath.Join(f.root, "dir", "f.txt"), "x")
			for rel, why := range map[string]string{"dir": "is a directory", "nope.txt": "no such file or directory"} {
				if got := f.deleteOp(t, rel); got["ok"] != false || !strings.Contains(fmt.Sprint(got["refused"]), why) {
					t.Fatalf("delete %s = %v, want refused %q", rel, got, why)
				}
			}
			if got := f.deleteOp(t, "dir/f.txt"); got["ok"] != true {
				t.Fatalf("delete dir/f.txt = %v", got)
			}
			if _, err := os.Lstat(filepath.Join(f.root, "dir", "f.txt")); !os.IsNotExist(err) {
				t.Fatalf("dir/f.txt still present: %v", err)
			}
			if info, err := os.Stat(filepath.Join(f.root, "dir")); err != nil || !info.IsDir() {
				t.Fatalf("deleting dir/f.txt removed its directory: %v", err)
			}
		})
	}
}

// TestWorkspaceWriteNeverDeletes is the authority split on the real binary:
// the delete op under Workspace.Write is refused by the handler before
// policy-tool runs, and the file stays.
func TestWorkspaceWriteNeverDeletes(t *testing.T) {
	f := newRealFixture(t, "worktree")
	_, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceWrite, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"op":"delete","path":"data.txt"}`))
	var refusal *AilangToolRefusalError
	if !errors.As(err, &refusal) || refusal.Op != "delete" {
		t.Fatalf("Workspace.Write delete: err = %v, want the op-allowlist refusal", err)
	}
	if data := mustRead(t, filepath.Join(f.root, "data.txt")); string(data) != "hello from ep1\n" {
		t.Fatalf("data.txt changed: %q", data)
	}
}
