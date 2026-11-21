package daemon

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceCommitEvidenceUsesExactGitObject(t *testing.T) {
	repo := t.TempDir()
	git := func(input string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Stdin = strings.NewReader(input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
		return strings.TrimSpace(string(output))
	}
	git("", "init", "-q")
	git("", "config", "user.name", "World Test")
	git("", "config", "user.email", "world@example.invalid")
	git("", "config", "commit.gpgSign", "false")
	tree := git("", "mktree")
	base := git("Base\n", "commit-tree", tree)
	const subject = "feat(surface): keep  meaningful spacing"
	const body = "First paragraph.\n\nSecond paragraph is indented:\n\tcontent\n"
	branch := git(subject+"\n\n"+body, "commit-tree", tree, "-p", base)
	main := git("Main\n", "commit-tree", tree, "-p", base)
	merge := git("Merge pull request #42\n\nMerge message.\n", "commit-tree", tree, "-p", main, "-p", branch)
	git("", "update-ref", "refs/heads/dev", merge)
	source := surfacefollow.Repo{ID: "world", Repo: repo, Branch: "dev"}
	d := &Daemon{follow: &surfaceFollow{config: surfacefollow.Config{Repos: []surfacefollow.Repo{source}}}}
	ctx := boundedTestContext(t)
	stack := func(kind, address string) surfacepage.FollowStack {
		return surfacepage.FollowStack{Repo: "world", StackModel: surfacepage.StackModel{Front: surfacepage.StackFront{
			ModelState: surfacepage.ModelState{ID: address, Kind: kind}, Run: &surfacepage.NodeRun{Episode: "selected"},
		}}}
	}
	finish := &EpisodeFinish{Episode: "selected", Outcome: "done", Commit: branch}
	for _, tc := range []struct {
		name, kind, address string
		finish              *EpisodeFinish
		want                []surfacepage.NodeCommit
	}{
		{"merge", surface.KindShared, "world:commit:" + merge, finish,
			[]surfacepage.NodeCommit{{ID: "world:commit:" + merge, Kinds: []string{"merge"}, Subject: "Merge pull request #42", Body: "Merge message.\n"}}},
		{"tip-and-finish", surface.KindStep, "world:commit:" + branch, finish,
			[]surfacepage.NodeCommit{{ID: "world:commit:" + branch, Kinds: []string{"branch-tip", "episode-finish"}, Subject: subject, Body: body}}},
		{"finish-different-tip", surface.KindStep, "world:commit:" + branch, &EpisodeFinish{Episode: "selected", Outcome: "done", Commit: main},
			[]surfacepage.NodeCommit{{ID: "world:commit:" + branch, Kinds: []string{"branch-tip"}, Subject: subject, Body: body}}},
		{"open", surfacepage.KindOpen, "world:episode:selected", nil, nil},
		{"open-with-verified-finish", surfacepage.KindOpen, "world:episode:selected", finish,
			[]surfacepage.NodeCommit{{ID: "world:commit:" + branch, Kinds: []string{"episode-finish"}, Subject: subject, Body: body}}},
		{"cross-episode", surfacepage.KindOpen, "world:episode:other", finish, nil},
		{"revision-expression", surfacepage.KindOpen, "world:episode:selected", &EpisodeFinish{Episode: "selected", Outcome: "done", Commit: "dev^"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := d.readSurfaceCommits(ctx, stack(tc.kind, tc.address), tc.finish)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestSurfaceCommitMessageRedactsPrivateTextAndCapsBytes(t *testing.T) {
	commit := "tree 0000000000000000000000000000000000000000\n\nFix by admin@example.invalid\n\nOutput at /Users/example/private\n"
	got, ok := parseSurfaceCommitMessage(commit)
	if !ok || !got.Redacted || got.Subject != "Fix by [email]" || !strings.Contains(got.Body, "[path]") {
		t.Fatalf("private text did not redact: %+v", got)
	}
	if _, ok := parseSurfaceCommitMessage(commit + strings.Repeat("x", surfaceCommitMessageBytes)); ok {
		t.Fatal("accepted a commit message above the surface cap")
	}
}
