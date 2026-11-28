package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

var (
	errOperatorFinishTerminal  = errors.New("unsandboxed check requires terminal finish")
	errOperatorFinishUncertain = errors.New("finish outcome uncertain; inspect durable receipts")
	errOperatorFinishCheck     = errors.New("sandboxed check failed; local commit withheld")
	errOperatorCommitFailed    = errors.New("local commit failed; retained the check receipts")
)

type operatorFinishWork struct {
	prep       operatorFinishPrep
	generation int64
	manifest   store.OperatorFinishReview
	git        operatorFinishGit
	report     daemon.EpisodeReport
	spool      string
	read       *store.Credentials
	fresh      func(context.Context) error
	now        func() int64
	budget     func(time.Duration) (context.Context, context.CancelFunc)
	run        func(context.Context, store.StartCheck) (checkrun.Result, error)
	record     func(context.Context, demoObject) (int64, error)
	evidence   func(context.Context, store.OperatorFinishCheck) bool
	commit     func(context.Context, store.OperatorLocalCommit, [][]byte, string) error
}

func (w operatorFinishWork) step(kind, identity string) store.OperatorStepKey {
	r := w.prep.row
	return store.OperatorStepKey{Action: operatorPreparationCAS(r), Owner: r.Owner, Kind: kind, Identity: r.ID + ":" + identity, Version: 1}
}

func (w operatorFinishWork) reserve(k store.OperatorStepKey) error {
	ctx, cancel := w.budget(time.Second)
	defer cancel()
	return runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
		return db.ReserveOperatorStep(ctx, k, w.now())
	})
}

func (w operatorFinishWork) observe(k store.OperatorStepKey, receipt any) error {
	raw, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	ctx, cancel := w.budget(time.Second)
	defer cancel()
	return runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
		return db.ObserveOperatorStep(ctx, k, raw)
	})
}

func (w operatorFinishWork) inputs(ctx context.Context) ([][]byte, error) {
	snap, b, err := w.git.snapshot(ctx, w.report.Start)
	if err != nil || !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {
		return nil, errOperatorFinishReview
	}
	return b, nil
}

// runEpisodeFinishWithPrep is the private entry beside terminal finish. Its
// caller has the episode and runner locks, and has completed Git auth before
// the claim. All effects use the exact inputs in the locally verified claim.
func runEpisodeFinishWithPrep(ctx context.Context, w operatorFinishWork) error {
	if !w.prep.claimed() {
		return approval.ErrNotVerified
	}
	m := w.manifest
	if !w.boundManifest(w.generation) {
		return errOperatorFinishReview
	}
	for _, check := range m.Checks {
		if check.Unsandboxed {
			return errOperatorFinishTerminal
		}
	}
	if m.Episode != w.prep.action.Episode || m.Store != w.prep.action.Store || w.report.Start == nil ||
		len(m.Unattributed) != 0 || len(m.Protected) != 0 || len(m.Checks) > 31 || len(m.Commit.Snapshot.Paths) == 0 {
		return errOperatorFinishFence
	}
	j, err := w.read.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(w.prep.row))
	if err != nil {
		return err
	}
	if j.CommitPhase == "reserved" {
		return w.reconcileCommit(j.CommitKey)
	}
	if j.CommitPhase == "observed" {
		if j.Commit.State == store.OperatorLocalCommitDone {
			return nil
		}
		return errOperatorCommitFailed
	}
	if j.OpenCheck || len(j.Checks) > len(m.Checks) {
		return errOperatorFinishUncertain
	}
	for i, old := range j.Checks {
		if old.Name != m.Checks[i].Name || !w.evidence(ctx, old) {
			return errOperatorFinishUncertain
		}
		if old.Result != "passed" {
			return errOperatorFinishCheck
		}
	}
	contents, err := w.inputs(ctx)
	if err != nil {
		return err
	}
	for i, check := range m.Checks[len(j.Checks):] {
		if err := w.fresh(ctx); err != nil {
			return err
		}
		if err := w.check(ctx, check, i+len(j.Checks)); err != nil {
			return err
		}
		contents, err = w.inputs(ctx)
		if err != nil {
			return err
		}
	}
	if err := w.fresh(ctx); err != nil {
		return err
	}
	tree, err := w.git.planTree(ctx, w.spool, m.Commit, contents)
	if err != nil {
		return err
	}
	contents, err = w.inputs(ctx)
	if err != nil {
		return err
	}
	if err := w.fresh(ctx); err != nil {
		return err
	}
	key := w.step("commit", tree)
	if err := w.reserve(key); err != nil {
		return err
	}
	commitErr := w.commit(ctx, m.Commit, contents, tree)
	if err := w.reconcileCommit(key); err == nil {
		return nil
	}
	if errors.Is(commitErr, errOperatorCommitNoRef) {
		budget, cancel := w.budget(3 * time.Second)
		defer cancel()
		if _, err := w.inputs(budget); err == nil {
			if err := w.observe(key, store.OperatorFinishLocalCommit{State: "commit-failed"}); err == nil {
				return errOperatorCommitFailed
			}
		}
	}
	return errOperatorFinishUncertain
}

func (w operatorFinishWork) check(ctx context.Context, check store.StartCheck, index int) error {
	key := w.step("checks", fmt.Sprintf("%02d", index))
	if err := w.reserve(key); err != nil {
		return err
	}
	runCtx, stop := operatorChildBudget(ctx, time.Duration(check.TimeoutS*float64(time.Second)))
	result, err := w.run(runCtx, check)
	stop()
	if err != nil {
		// No trustworthy terminal readback. Retain the reservation and let
		// an operator inspect evidence. A later runner must not retry it.
		return fmt.Errorf("%w: sandboxed check interrupted", errOperatorFinishUncertain)
	}
	passed := result.ExitCode == 0
	outcome := finishOutcome(w.report, check, result, passed)
	raw, _ := json.Marshal(outcome)
	canonical, _, err := parseTaskOutcome(raw)
	if err != nil {
		return errOperatorFinishUncertain
	}
	obj := taskOutcomeObject(canonical)
	writeCtx, cancel := w.budget(3 * time.Second)
	defer cancel()
	entry, err := w.record(writeCtx, obj)
	if err != nil || entry < 1 {
		return errOperatorFinishUncertain
	}
	class := "passed"
	if !passed {
		class = "failed"
	}
	if err := w.observe(key, store.OperatorFinishCheck{Name: check.Name, Result: class, Ref: obj.Hash, Entry: entry}); err != nil {
		return err
	}
	if !passed {
		return errOperatorFinishCheck
	}
	return nil
}

func (w operatorFinishWork) reconcileCommit(key store.OperatorStepKey) error {
	tree, ok := strings.CutPrefix(key.Identity, w.prep.row.ID+":")
	if !ok || !operatorFinishSHA.MatchString(tree) {
		return errOperatorFinishUncertain
	}
	ctx, cancel := w.budget(3 * time.Second)
	defer cancel()
	if _, err := w.git.config(ctx, w.report.Start); err != nil {
		return errOperatorFinishUncertain
	}
	sha, err := w.git.inspectCommit(ctx, w.manifest.Commit, tree, w.report.Start.Agent)
	if err != nil {
		return errOperatorFinishUncertain
	}
	return w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: sha})
}

func (w operatorFinishWork) boundManifest(generation int64) bool {
	a, err := w.manifest.Action(generation, w.prep.action.Options.Ready)
	return err == nil && slices.Equal(a.Resolution, w.prep.action.Resolution) && a.Review == w.prep.action.Review
}
