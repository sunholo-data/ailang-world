package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

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
	ctx := boundedTestContext(t)
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
	ctx := boundedTestContext(t)
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
	head, err := g.inspectCommit(ctx, spec, tree)
	if err != nil || head == initial.Head {
		t.Fatalf("exact descendant missing: %s %v", head, err)
	}
	raw := f.gitOutput("-C", f.wt, "cat-file", "commit", head)
	for _, field := range []string{"parent", "tree", "message", "trailer", "committer"} {
		bad, want := spec, tree
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
		}
		if operatorFinishCommitMatches(raw, bad, want) {
			t.Fatalf("%s mismatch observed", field)
		}
	}
	f.write("extra", "partial stage")
	f.gitOutput("-C", f.wt, "add", "--", "extra")
	if _, err := g.inspectCommit(ctx, spec, tree); err == nil || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
		t.Fatal("uncertain partial stage adopted or rewrote HEAD")
	}
}
