package surfacefollow

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestLandedFilesFromActualMerge(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		git("add", "--", name)
	}
	git("init", "-q", "-b", "dev")
	git("config", "user.name", "World Test")
	git("config", "user.email", "world@example.invalid")
	git("config", "commit.gpgSign", "false")
	write("charter.md", "220. **surface-follow** · goal\n")
	git("commit", "-qm", "base")
	git("switch", "-qc", "topic")
	write("side-only.txt", "never lands\n")
	git("commit", "-qm", "topic")
	git("switch", "-q", "dev")
	write("already-on-dev.txt", "unchanged by merge\n")
	git("commit", "-qm", "dev update")
	git("merge", "--no-ff", "--no-commit", "topic")
	git("rm", "-qf", "--", "side-only.txt")
	path := "visible\tfile\n"
	write(path, "created in merge\n")
	git("commit", "-qm", "Merge pull request #7 from owner/topic/row220")
	landed := git("rev-parse", "HEAD")
	git("switch", "-qc", "empty-topic")
	git("commit", "--allow-empty", "-qm", "empty")
	git("switch", "-q", "dev")
	git("merge", "--no-ff", "-qm", "Merge pull request #8 from owner/topic/row220", "empty-topic")
	empty := git("rev-parse", "HEAD")

	r := testRepo()
	r.Repo = dir
	got, err := NewCollector(nil, nil).BuildPinned(context.Background(), r, empty)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]surface.FileChanges{
		landed: {Changes: []surface.Change{{Path: path, Status: "A"}}},
		empty:  {},
	}
	if !got.LandedFilesChecked || !reflect.DeepEqual(got.Index.LandedFiles, want) {
		t.Fatalf("landed files = %#v, want %#v", got.Index.LandedFiles, want)
	}
}

func TestLandedFilesTruncationAndInvalidOutput(t *testing.T) {
	id := sha(3)
	var out strings.Builder
	out.WriteString(id + "\x00")
	for i := 0; i < surface.MaxChanges+1; i++ {
		out.WriteString("M\x00file/" + strconv.Itoa(i) + "\x00")
	}
	files := parseLandedFiles(out.String(), []string{id})[id]
	if len(files.Changes) != surface.MaxChanges || files.Truncated != 1 {
		t.Fatalf("listed %d, truncated %d", len(files.Changes), files.Truncated)
	}
	for _, malformed := range []string{out.String() + "A\x00unterminated", out.String() + id + "\x00", id + "\x00A\x00../outside\x00"} {
		if got := parseLandedFiles(malformed, []string{id}); got != nil {
			t.Fatal("accepted malformed record stream")
		}
	}
}

func TestLandedFilesUnavailableLeavesGraph(t *testing.T) {
	g := newGit()
	g.setHook(func(_ context.Context, args []string) (string, error, bool) {
		if args[3] == "--no-lazy-fetch" {
			return "", errors.New("output over cap"), true
		}
		return "", nil, false
	})
	got, err := NewCollector(g, newClock()).BuildPinned(context.Background(), testRepo(), sha(3))
	if err != nil || !got.LandedFilesChecked || len(got.Index.States) != 2 || len(got.Index.LandedFiles) != 0 {
		t.Fatalf("expected usable graph with unknown landed files: %+v %v", got, err)
	}
}

func TestLandedFilesBackfillsLegacyCacheAtSameHead(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	legacy := partitionFor(t, r, 3)
	legacy.LandedFilesChecked = false
	legacy.Index.LandedFiles = nil
	deliver(p, legacy)
	oldGeneration := p.Snapshot().Status.Generation
	restarted, err := NewPublisher(p.config, []byte("config"), filepath.Dir(filepath.Dir(p.path)), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	g, clock := newGit(), newClock()
	g.move(p.config.Repos[0].Repo, legacy.Head)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = (&Follower{Config: p.config, Collector: NewCollector(g, clock),
			Sink: func(s Status) {
				restarted.Sink(s)
				restarted.publish()
			}, Seed: restarted.Seed()}).Run(ctx)
	}()
	until(t, func() bool {
		s := restarted.Snapshot()
		return len(s.Envelope.Repos) == 1 && len(s.Envelope.Repos[0].Index.LandedFiles) == 1
	})
	current := restarted.Snapshot()
	if current.Status.Generation == oldGeneration || current.Envelope.Repos[0].Head != legacy.Head {
		t.Fatal("backfill did not publish a new generation of the same head")
	}
	clock.advance(PollInterval)
	until(t, func() bool { return countCommand(g, "rev-parse") == 2 })
	cancel()
	<-done
	for _, command := range []string{"log", "merge-base", "cat-file"} {
		if countCommand(g, command) != 0 {
			t.Fatalf("legacy backfill reread graph: %v", g.recorded())
		}
	}
	if countCommand(g, "show") != 1 {
		t.Fatalf("expected exactly one backfill Git read: %v", g.recorded())
	}
	reloaded, err := NewPublisher(p.config, []byte("config"), filepath.Dir(filepath.Dir(p.path)), nil, nil)
	if err != nil || len(reloaded.Seed()[r.ID].Index.LandedFiles) != 1 || reloaded.Snapshot().Status.Generation != current.Status.Generation {
		t.Fatal("backfilled evidence did not survive reload", err)
	}
}

func TestLandedFilesYieldToManifestCap(t *testing.T) {
	part := partitionFor(t, testRepo(), 3)
	// Adjust an unrelated field to put the envelope at the byte cap without
	// file evidence. Validation of the graph is covered by the publisher tests.
	const generation = "candidate"
	ix := surface.FollowIndex{Schema: surface.FollowSchema, Generation: generation, Repos: []surface.FollowPartition{part}}
	ix.Repos[0].Index.LandedFiles = nil
	baseline, class, _ := boundedManifest(&ix)
	if class != "" {
		t.Fatal(class)
	}
	ix.Generation = generation + strings.Repeat(" ", surface.MaxBytes-len(baseline))
	ix.Repos[0].Index.LandedFiles = part.Index.LandedFiles
	data, class, dropped := boundedManifest(&ix)
	if class != "" || !dropped || len(data) != surface.MaxBytes || len(part.Index.LandedFiles) != 1 {
		t.Fatalf("size %d, class %q, dropped %t, source files %d", len(data), class, dropped, len(part.Index.LandedFiles))
	}
}

func TestLandedFilesRetriesTransientReadOnUnchangedHead(t *testing.T) {
	g, clock := newGit(), newClock()
	g.setHook(func(_ context.Context, args []string) (string, error, bool) {
		if args[3] == "--no-lazy-fetch" {
			return "", context.DeadlineExceeded, true
		}
		return "", nil, false
	})
	statuses, _ := startFollower(t, []Repo{testRepo()}, g, clock)
	until(t, func() bool { return statuses.accepted(sha(3)) })
	g.setHook(nil)
	latest := func() Status {
		statuses.mu.Lock()
		defer statuses.mu.Unlock()
		return statuses.all[len(statuses.all)-1]
	}
	if part := latest().Partition; part.LandedFilesChecked || len(part.Index.LandedFiles) != 0 {
		t.Fatal("a transient error permanently settled merge evidence")
	}
	for elapsed := PollInterval; elapsed < landedFilesRetry; elapsed += PollInterval {
		clock.advance(PollInterval)
		check := clock.Now().Format("2006-01-02T15:04:05Z07:00")
		until(t, func() bool { return latest().CheckedAt == check })
		if countCommand(g, "show") != 1 {
			t.Fatal("optional read retried before cooldown")
		}
	}
	clock.advance(PollInterval)
	until(t, func() bool { return len(latest().Partition.Index.LandedFiles) == 1 })
	if countCommand(g, "log") != 1 || countCommand(g, "show") != 2 {
		t.Fatal("retry rebuilt the graph or did not batch file evidence")
	}
}
