package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommitExplicitCommitter(t *testing.T) {
	git := gitBin(t)
	repo := newRepo(t, git)
	ctx := testCtx(t)
	for _, kv := range [][2]string{{"user.name", "Repository Test"}, {"user.email", "repo@example.invalid"}, {"commit.gpgsign", "false"}} {
		if _, err := Output(ctx, git, []string{"-C", repo, "config", kv[0], kv[1]}); err != nil {
			t.Fatal(err)
		}
	}
	for _, identity := range []string{"Test Operator <operator@example.invalid>", ""} {
		if err := os.WriteFile(filepath.Join(repo, "identity.txt"), []byte(identity), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := CommitWithCommitter(ctx, git, repo, "codex", identity, []string{"identity.txt"}, "Identity"); err != nil {
			t.Fatal(err)
		}
		got, err := Output(ctx, git, []string{"-C", repo, "log", "-1", "--format=%an <%ae>|%cn <%ce>"})
		want := identity
		if want == "" {
			want = "Repository Test <repo@example.invalid>"
		}
		if err != nil || strings.TrimSpace(got) != "codex via AILANG World <ailang-world@invalid>|"+want {
			t.Fatalf("%q %v", got, err)
		}
	}
	if _, err := CommitWithCommitter(ctx, git, repo, "codex", "Bad\nName <bad@example.invalid>", []string{"identity.txt"}, "No"); err == nil {
		t.Fatal("newline accepted")
	}
}

func TestParseCommitter(t *testing.T) {
	for _, value := range []string{"Name", "<email>", "Name <>", "Name <email> tail", "Name <<email>>", "Name >email<", "Name\r <email>", "Name\n <email>"} {
		if _, _, err := ParseCommitter(value); err == nil {
			t.Errorf("accepted %q", value)
		}
	}
}
