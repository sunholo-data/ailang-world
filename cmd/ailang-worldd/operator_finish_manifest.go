package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type operatorFinishInputs struct {
	manifest store.OperatorFinishReview
	report   daemon.EpisodeReport
	git      operatorFinishGit
	token    string
	profile  *broker.ExecProfile
	world    string
}

func (i operatorFinishInputs) frozen(ctx context.Context) (store.OperatorFinishReview, error) {
	v := i.manifest
	tree, _, err := i.git.snapshot(ctx, i.report.Start)
	if err != nil {
		return v, err
	}
	if !i.baseMatches(ctx) || len(tree.Paths) == 0 {
		return v, errOperatorFinishReview
	}
	if err := validateFinishChecks(v.Checks, nil, i.git.dir); err != nil {
		return v, errOperatorFinishFence
	}
	var paths []string
	for _, p := range tree.Paths {
		paths = append(paths, p.Path)
	}
	v.Protected = finishProtectedPaths(i.profile, paths)
	v.Unattributed = finishUnattributed(paths, i.report)
	v.Commit.Snapshot = tree
	v.Commit.Message += fmt.Sprintf("\n\n%s %s; %d calls; %d/%d checks passed; %d files changed.\n\n%s\n",
		i.report.Start.Agent, i.report.Start.Model, i.report.TotalCalls, len(v.Checks), len(v.Checks), len(paths), reportTrailer(i.report, i.world))
	if scrubMarkdown(v.Commit.Message) != nil {
		return v, errOperatorFinishFence
	}
	return v, nil
}

func (i operatorFinishInputs) baseMatches(ctx context.Context) bool {
	if _, err := i.git.config(ctx, i.report.Start); err != nil {
		return false
	}
	base, err := i.git.git(ctx, nil, nil, "rev-parse", "--verify", "refs/heads/"+i.manifest.Base.Branch+"^{commit}")
	return err == nil && strings.TrimSpace(base) == i.manifest.Base.SHA
}

func finishManifestForRow(c runnerConfig, row store.OperatorActionRow) (approval.OperatorAction, store.OperatorFinishReview, int64, error) {
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil || a.Kind != "finish" || a.Episode != row.Episode || a.Store != c.Credentials {
		return approval.OperatorAction{}, store.OperatorFinishReview{}, 0, approval.ErrNotVerified
	}
	var m store.OperatorFinishReview
	generation, genErr := strconv.ParseInt(operatorFileFact(a, "finish_catalogue"), 10, 64)
	if json.Unmarshal([]byte(operatorFileFact(a, "finish_manifest")), &m) != nil || genErr != nil || generation <= 0 ||
		!slices.Contains(c.FinishRemotes, m.Remote.URL) {
		return approval.OperatorAction{}, store.OperatorFinishReview{}, 0, approval.ErrNotVerified
	}
	return a, m, generation, nil
}

func finishWorldEvidence(bound context.Context, path string, receipt store.OperatorFinishCheck) bool {
	hash, err := hashref.Parse(receipt.Ref)
	world, openErr := store.OpenReadOnly(path)
	if openErr != nil {
		return false
	}
	defer world.Close()
	object, found, readErr := world.GetObject(bound, hash)
	return err == nil && readErr == nil && found && hashref.SumSHA256(object.Payload) == hash
}
