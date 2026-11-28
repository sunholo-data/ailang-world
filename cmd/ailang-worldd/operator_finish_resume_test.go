package main

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorFinishResumeEvidence(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	ctx := activationTestContext(t)
	if err := w.check(ctx, w.manifest.Checks[0], 0); err != nil {
		t.Fatal(err)
	}
	w.evidence = func(ctx context.Context, c store.OperatorFinishCheck) bool { return finishWorldEvidence(ctx, f.db, c) }
	before := operatorFinishJournal(t, w)
	if err := runEpisodeFinishWithPrep(ctx, w); !errors.Is(err, errOperatorFinishUncertain) ||
		!reflect.DeepEqual(before, operatorFinishJournal(t, w)) || len(before.Checks) != 1 || before.Checks[0].Result != "passed" ||
		f.checks != 1 || f.outcomes != 1 || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
		t.Fatalf("resume trusted a prior pass without matching World evidence: %v", err)
	}
}

func TestOperatorCommitFailureInputRecheck(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
		seedStartFixtureFile(t, filepath.Join(f.wt, "unexpected"), []byte("changed after reservation"))
		return errOperatorCommitNoRef
	}
	err := runEpisodeFinishWithPrep(activationTestContext(t), w)
	j := operatorFinishJournal(t, w)
	if !errors.Is(err, errOperatorFinishUncertain) || j.CommitPhase != "reserved" || j.Commit.State != "" ||
		len(j.Checks) != 1 || f.checks != 1 || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
		t.Fatalf("changed inputs were recorded as a definite commit failure: %v %+v", err, j)
	}
}

func TestOperatorCommitResumeReserved(t *testing.T) {
	f, first, now := operatorLocalFinishFixture(t, nil)
	ctx := activationTestContext(t)
	b, err := first.inputs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	spec := first.manifest.Commit
	tree, err := first.git.planTree(ctx, first.spool, spec, b)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.reserve(first.step("commit", tree)); err != nil {
		t.Fatal(err)
	}
	if err := first.commit(ctx, spec, b, tree); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.OpenCredentialsReadOnly(first.manifest.Store)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	p := first.prep
	p, err = verifyOperatorFinishPrep(ctx, reopened, p.action, *now, operatorPreparationCAS(p.row), p.row.Owner)
	if err != nil {
		t.Fatal(err)
	}
	resumed := operatorFinishWork{prep: p, manifest: first.manifest, generation: first.generation, read: reopened,
		git: first.git, report: first.report, budget: first.budget}
	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
	index := filepath.Join(first.git.admin, "index")
	config := filepath.Join(first.git.common, "config")
	oldIndex, oldConfig := serviceFileHash(index), readStartFixtureFile(t, config)
	seedStartFixtureFile(t, config, append(oldConfig, []byte("\n[credential]\n\thelper = false\n")...))
	before := operatorFinishJournal(t, resumed)
	if err := runEpisodeFinishWithPrep(ctx, resumed); !errors.Is(err, errOperatorFinishUncertain) ||
		!reflect.DeepEqual(before, operatorFinishJournal(t, resumed)) || before.CommitPhase != "reserved" ||
		serviceFileHash(index) != oldIndex || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
		t.Fatalf("restart adopted a commit despite unsafe Git config: %v", err)
	}
	seedStartFixtureFile(t, config, oldConfig)
	if err := runEpisodeFinishWithPrep(ctx, resumed); err != nil {
		t.Fatal(err)
	}
	if j := operatorFinishJournal(t, resumed); j.Commit != (store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: head}) ||
		head == f.before || serviceFileHash(index) != oldIndex || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head || f.checks != 0 {
		t.Fatalf("restart failed to adopt the exact reserved descendant read-only: %+v", j)
	}
}
