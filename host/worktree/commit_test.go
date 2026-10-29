package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommitLiteralPathsIdentityAndTrailer(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	ctx := testCtx(t)
	for _, pair := range [][2]string{{"user.name", "Operator"}, {"user.email", "operator@example.invalid"}, {"commit.gpgsign", "false"}} {
		if _, err := Output(ctx, git, []string{"-C", repo, "config", pair[0], pair[1]}); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"a[1].go", "a1.go", "unrelated.txt"} {
		if err := os.WriteFile(filepath.Join(repo, name), []byte("hello\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// An existing unrelated index entry must not sneak into this commit.
	if _, err := Output(ctx, git, []string{"-C", repo, "add", "unrelated.txt"}); err != nil {
		t.Fatal(err)
	}
	trailer := "World-Provenance: store=/fixture/world.db episode=ep1 entries=1-6000"
	sha, err := Commit(ctx, git, repo, "codex", []string{"a[1].go"}, "Title", "Paragraph.", trailer)
	if err != nil || len(sha) != 40 {
		t.Fatalf("commit %s %v", sha, err)
	}
	out, err := Output(ctx, git, []string{"-C", repo, "show", "--format=", "--name-only", sha})
	if err != nil || strings.TrimSpace(out) != "a[1].go" {
		t.Fatalf("literal paths: %q %v", out, err)
	}
	out, _ = Output(ctx, git, []string{"-C", repo, "log", "-1", "--format=%an <%ae>|%cn <%ce>%n%B"})
	if !strings.Contains(out, "codex via AILANG World <ailang-world@invalid>|Operator <operator@example.invalid>") || !strings.HasSuffix(strings.TrimSpace(out), trailer) {
		t.Fatal(out)
	}
	if _, err := Commit(ctx, git, repo, "codex", []string{"../escape"}, "no"); err == nil {
		t.Fatal("nonlocal path accepted")
	}
}
