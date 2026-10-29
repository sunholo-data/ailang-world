package surfacefollow

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/worktree"
)

func TestRealGitLeavesSourceUnchanged(t *testing.T) {
	repo := t.TempDir()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal("git required for production seam proof:", err)
	}
	command := func(args ...string) string {
		t.Helper()
		cmd := exec.Command(git, append([]string{"-C", repo, "-c", "user.name=Fixture",
			"-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false"}, args...)...)
		cmd.Env = worktree.Env(os.Environ())
		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture git %v: %v: %s", args, err, b)
		}
		return string(b)
	}
	command("init", "-b", "main")
	if err := os.WriteFile(filepath.Join(repo, "charter.md"), []byte("220. **follow** · source collector\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	command("add", "charter.md")
	command("commit", "-m", "base")
	command("checkout", "-b", "trial/row220")
	if err := os.WriteFile(filepath.Join(repo, "work.txt"), []byte("step\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	command("add", "work.txt")
	command("commit", "-m", "step")
	command("checkout", "main")
	command("merge", "--no-ff", "trial/row220", "-m", "Merge pull request #7 from owner/trial/row220")
	snapshot := func() map[string]string {
		t.Helper()
		out := map[string]string{"refs": command("for-each-ref", "--format=%(refname) %(objectname)")}
		for _, name := range []string{".git/index", ".git/HEAD"} {
			b, err := os.ReadFile(filepath.Join(repo, name))
			if err != nil {
				t.Fatal(err)
			}
			out[name] = string(b)
		}
		err := filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(repo, path)
			out[rel] = string(b)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	before := snapshot()
	canonical, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	r := Repo{ID: "world", Repo: canonical, Charter: "charter.md", Branch: "main", Limit: 200}
	c := NewCollector(ReadGit{}, newClock())
	p, err := c.Collect(context.Background(), r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Index.States) != 2 || p.Head != strings.TrimSpace(command("rev-parse", "HEAD")) {
		t.Fatal("real merge not projected at its pinned head")
	}
	if _, err := c.Collect(context.Background(), r, &p); err != nil {
		t.Fatal(err)
	}
	after := snapshot()
	if !reflect.DeepEqual(before, after) {
		t.Fatal("source refs, git index, HEAD or worktree bytes changed")
	}
}
