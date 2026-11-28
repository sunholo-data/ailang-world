package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"golang.org/x/sys/unix"
)

func operatorFinishGitFixture(t *testing.T, f *finishFixture) operatorFinishGit {
	t.Helper()
	g, err := openOperatorFinishGit(operatorFinishAuth{Git: f.git, Home: f.root, Path: runnerSystemPath},
		f.wt, filepath.Join(f.root, "repo"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestOperatorFinishWorktreeFence(t *testing.T) {
	f := newFinishFixture(t)
	g := operatorFinishGitFixture(t, f)
	ctx := activationTestContext(t)
	path := filepath.Join(g.common, "config")
	config := readStartFixtureFile(t, path)
	if _, _, err := g.snapshot(ctx, f.rep.Start); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(f.root, "sentinel")
	payload := "touch " + shellQuote(sentinel)
	for _, setting := range []string{"core.hooksPath", "filter.fixture.clean", "credential.helper", "core.fsmonitor"} {
		f.gitOutput("-C", f.wt, "config", setting, "!"+payload)
		if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
			t.Fatalf("%s was accepted: %v", setting, err)
		}
		seedStartFixtureFile(t, path, config)
	}
	bad := filepath.Join(f.wt, "fixture-link")
	if err := os.Symlink(sentinel, bad); err != nil {
		t.Fatal(err)
	}
	if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
		t.Fatalf("dirty symlink was accepted: %v", err)
	}
	if err := os.Remove(bad); err != nil {
		t.Fatal(err)
	}
	if err := unix.Mkfifo(bad, 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
		t.Fatalf("dirty special file was accepted: %v", err)
	}
	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("untrusted program executed or sentinel unobservable: %v", err)
	}
}

func TestOperatorCommitIdentity(t *testing.T) {
	f := newFinishFixture(t)
	g := operatorFinishGitFixture(t, f)
	ctx := activationTestContext(t)
	initial, contents, err := g.snapshot(ctx, f.rep.Start)
	if err != nil {
		t.Fatal(err)
	}
	spec := store.OperatorLocalCommit{Snapshot: initial, Message: "fixture\n\nWorld-Provenance: episode=ep1\n", Committer: "Fixture Operator <fixture@example.invalid>"}
	tree, err := g.planTree(ctx, f.root, spec, contents)
	if err != nil {
		t.Fatal(err)
	}
	seen, _, err := g.snapshot(ctx, f.rep.Start)
	if err != nil || !reflect.DeepEqual(seen, initial) {
		t.Fatalf("private staging changed production index or refs: %v", err)
	}
	sentinel := filepath.Join(f.root, "git-hook")
	for _, hook := range []string{"pre-commit", "reference-transaction"} {
		if err := os.WriteFile(filepath.Join(g.common, "hooks", hook), []byte("#!/bin/sh\ntouch "+shellQuote(sentinel)+"\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.commit(ctx, f.root, spec, contents, tree, f.rep.Start); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("repo-controlled hook executed or sentinel unreadable: %v", err)
	}
	head, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent)
	if err != nil || head == initial.Head {
		t.Fatalf("exact descendant missing: %s %v", head, err)
	}
	raw := f.gitOutput("-C", f.wt, "cat-file", "commit", head)
	for _, field := range []string{"parent", "tree", "message", "trailer", "committer", "author-name", "author-email", "author-time"} {
		bad, want, candidate := spec, tree, raw
		switch field {
		case "parent":
			bad.Snapshot.Head = strings.Repeat("d", 40)
		case "tree":
			want = head
		case "message":
			bad.Message = "bad title\n\nWorld-Provenance: episode=ep1\n"
		case "trailer":
			bad.Message = strings.Replace(spec.Message, "episode=ep1", "episode=other", 1)
		case "committer":
			bad.Committer = "Someone Else <different@invalid>"
		case "author-name":
			candidate = strings.Replace(raw, "author codex via AILANG World ", "author other via AILANG World ", 1)
		case "author-email":
			candidate = strings.Replace(raw, "ailang-world@invalid", "other@invalid", 1)
		case "author-time":
			headers := strings.SplitN(raw, "\n", 4)
			identity, _, _ := strings.Cut(headers[2], "> ")
			headers[2] = identity + "> 1"
			candidate = strings.Join(headers, "\n")
		}
		if (strings.HasPrefix(field, "author") && candidate == raw) || operatorFinishCommitMatches(candidate, bad, want, f.rep.Start.Agent) {
			t.Fatalf("%s mismatch observed", field)
		}
	}
	f.gitOutput("-C", f.wt, "branch", "other-branch", head)
	f.gitOutput("-C", f.wt, "symbolic-ref", "HEAD", "refs/heads/other-branch")
	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); !errors.Is(err, errOperatorFinishFence) {
		t.Fatalf("adopted a commit on the wrong branch: %v", err)
	}
	f.gitOutput("-C", f.wt, "symbolic-ref", "HEAD", initial.Branch)
	f.write("extra", "partial stage")
	f.gitOutput("-C", f.wt, "add", "--", "extra")
	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); err == nil || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
		t.Fatal("uncertain partial stage adopted or rewrote HEAD")
	}
}

func TestOperatorFinishCaptureBounds(t *testing.T) {
	f := newFinishFixture(t)
	g := operatorFinishGitFixture(t, f)
	ctx := activationTestContext(t)
	capture := func() store.OperatorFinishTree {
		t.Helper()
		snap, _, err := g.snapshot(ctx, f.rep.Start)
		if err != nil {
			t.Fatal(err)
		}
		return snap
	}
	refuse := func() {
		t.Helper()
		before := serviceFileHash(filepath.Join(g.admin, "index"))
		if got, _, err := g.snapshot(ctx, f.rep.Start); !errors.Is(err, errOperatorFinishFence) || len(got.Paths) != 0 ||
			serviceFileHash(filepath.Join(g.admin, "index")) != before || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
			t.Fatalf("capture fence allowed a changed Git state: %v", err)
		}
	}
	first := capture()
	f.write("b.go", "edit\n")
	second := capture()
	if reflect.DeepEqual(first, second) || !slices.Contains(second.Paths, store.OperatorFinishPath{Path: "b.go", Mode: "100644", Digest: finishSHA([]byte("edit\n"))}) {
		t.Fatal("untracked same-size edit did not change the content digest")
	}
	for _, path := range []string{filepath.Join(g.admin, "config.worktree"), filepath.Join(g.common, "objects", "info", "alternates")} {
		seedStartFixtureFile(t, path, []byte("\n"))
		refuse()
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	f.gitOutput("-C", f.wt, "add", "--", "b.go")
	refuse()
	f.gitOutput("-C", f.wt, "read-tree", first.Head)

	extra := filepath.Join(f.wt, "bounds")
	if err := os.Mkdir(extra, 0700); err != nil {
		t.Fatal(err)
	}
	for n := range 254 {
		seedStartFixtureFile(t, filepath.Join(extra, fmt.Sprint(n)), nil)
	}
	if got := len(capture().Paths); got != 257 {
		t.Fatalf("exact file cap: %d", got)
	}
	seedStartFixtureFile(t, filepath.Join(extra, "overflow"), nil)
	refuse()
	if err := os.RemoveAll(extra); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.go", "b.go", "gen.txt"} {
		if err := os.Remove(filepath.Join(f.wt, name)); err != nil {
			t.Fatal(err)
		}
	}
	for _, p := range []struct {
		name string
		size int64
	}{{"large", 16 << 20}, {"overflow", 1}} {
		path := filepath.Join(f.wt, p.name)
		seedStartFixtureFile(t, path, nil)
		if err := os.Truncate(path, p.size); err != nil {
			t.Fatal(err)
		}
		if p.name == "large" {
			capture()
		}
	}
	refuse()
	files := fstest.MapFS{}
	tracked := map[string]bool{}
	for n := range 16383 {
		name := fmt.Sprintf("fixture-%05d", n)
		files[name] = &fstest.MapFile{}
		tracked[name] = true
	}
	if names, err := operatorFinishUntracked(ctx, files, tracked, nil); err != nil || len(names) != 0 {
		t.Fatalf("exact traversal bound including root: %v", err)
	}
	files["overflow"] = &fstest.MapFile{}
	if names, err := operatorFinishUntracked(ctx, files, tracked, nil); !errors.Is(err, errOperatorFinishFence) || len(names) != 0 {
		t.Fatalf("traversal beyond the bound: %v", err)
	}
}

func TestOperatorFinishReadReplacement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checked")
	replacement := filepath.Join(dir, "replacement")
	for _, p := range []string{path, replacement} {
		seedStartFixtureFile(t, p, []byte("same"))
		stamp := time.Unix(1000, 0)
		if err := os.Chtimes(p, stamp, stamp); err != nil {
			t.Fatal(err)
		}
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if b, _, err := readOperatorFinishFile(f, path, 4); err != nil || string(b) != "same" {
		t.Fatalf("held file control: %q %v", b, err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(replacement, path); err != nil {
		t.Fatal(err)
	}
	if b, _, err := readOperatorFinishFile(f, path, 4); !errors.Is(err, errOperatorFinishFence) || b != nil {
		t.Fatalf("same bytes and different inode escaped fence: %q %v", b, err)
	}
}

func TestOperatorCommitReservationFence(t *testing.T) {
	f := newFinishFixture(t)
	g := operatorFinishGitFixture(t, f)
	ctx := activationTestContext(t)
	initial, contents, err := g.snapshot(ctx, f.rep.Start)
	if err != nil {
		t.Fatal(err)
	}
	spec := store.OperatorLocalCommit{Snapshot: initial, Message: "fixture\n", Committer: "Fixture Operator <fixture@example.invalid>"}
	// Leave the valid object in this temp repo so an unchecked mismatch would commit it.
	tree, err := g.stage(ctx, spec, contents, []string{"GIT_INDEX_FILE=" + filepath.Join(f.root, "plan-index")})
	if err != nil {
		t.Fatal(err)
	}
	changed := slices.Clone(contents)
	changed[0] = []byte("different tree")
	if err := g.commit(ctx, f.root, spec, changed, tree, f.rep.Start); !errors.Is(err, errOperatorCommitNoRef) {
		t.Fatalf("candidate tree mismatch: %v", err)
	}
	wrap := func(g operatorFinishGit, script string) operatorFinishGit {
		t.Helper()
		path := filepath.Join(t.TempDir(), "git")
		if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+script+"\nexec "+shellQuote(f.git)+" \"$@\"\n"), 0700); err != nil {
			t.Fatal(err)
		}
		g.bin = path
		return g
	}
	unexpectedPath := filepath.Join(f.wt, "changed-after-commit-tree")
	afterObject := wrap(g, "for x; do if [ \"$x\" = commit-tree ]; then touch "+shellQuote(unexpectedPath)+"; fi; done")
	if err := afterObject.commit(ctx, f.root, spec, contents, tree, f.rep.Start); !errors.Is(err, errOperatorFinishReview) {
		t.Fatalf("post-object drift ignored: %v", err)
	}
	if err := os.Remove(unexpectedPath); err != nil {
		t.Fatal(err)
	}
	if seen, _, err := g.snapshot(ctx, f.rep.Start); err != nil || !reflect.DeepEqual(initial, seen) {
		t.Fatalf("refusal changed Git state: %v", err)
	}
	// Advance the pinned branch at the update-ref seam, after the final snapshot.
	foreign := strings.TrimSpace(f.gitOutput("-C", f.wt, "commit-tree", tree, "-p", initial.Head, "-m", "foreign"))
	refArgs := " --git-dir=" + shellQuote(g.admin) + " update-ref " + shellQuote(initial.Branch) + " " + shellQuote(foreign)
	race := wrap(g, "for x; do if [ \"$x\" = update-ref ]; then "+shellQuote(f.git)+refArgs+"; fi; done")
	if err := race.commit(ctx, f.root, spec, contents, tree, f.rep.Start); err == nil ||
		strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != foreign {
		t.Fatalf("old-value CAS overwrote a different HEAD: %v", err)
	}
	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); !errors.Is(err, errOperatorFinishFence) {
		t.Fatalf("recovery adopted the competing commit: %v", err)
	}
}

func TestOperatorFinishConfigReplacement(t *testing.T) {
	f := newFinishFixture(t)
	g := operatorFinishGitFixture(t, f)
	ctx := activationTestContext(t)
	path := filepath.Join(g.common, "config")
	if _, err := g.config(ctx, f.rep.Start); err != nil {
		t.Fatal(err)
	}
	swapped := filepath.Join(f.root, "replacement")
	body := readStartFixtureFile(t, path)
	seedStartFixtureFile(t, swapped, append(body, []byte("\n[core]\n\tfsmonitor = /no-fixture-program\n")...))
	sentinel := filepath.Join(f.root, "outside-config-parser")
	fake := filepath.Join(f.root, "git-config-race")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nset -eu\n"+
		"if [ \"$1\" = --git-dir=/dev/null ]; then mv "+shellQuote(swapped)+" "+shellQuote(path)+
		"; else touch "+shellQuote(sentinel)+"; exit 93; fi\nexec "+shellQuote(f.git)+" \"$@\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	g.bin = fake
	before := serviceFileHash(filepath.Join(g.admin, "index"))
	if _, _, err := g.snapshot(ctx, f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
		t.Fatalf("replacement was accepted: %v", err)
	}
	spec := store.OperatorLocalCommit{Committer: "Fixture Operator <fixture@example.invalid>"}
	if err := g.commit(ctx, f.root, spec, nil, "", f.rep.Start); !errors.Is(err, errOperatorFinishReview) {
		t.Fatalf("unsafe configuration reached commit staging: %v", err)
	}
	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) ||
		before != serviceFileHash(filepath.Join(g.admin, "index")) || bytes.Equal(body, readStartFixtureFile(t, path)) {
		t.Fatalf("replacement did not stop before using the live Git configuration: %v", err)
	}
}
