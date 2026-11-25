package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

func reconcileWorktreeCommand(id string) string {
	return "ailang-worldd service runner reconcile-worktree " + id
}

func readRunnerWorktree(credentials, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
	var target operatorWorktreeTarget
	db, err := store.OpenCredentialsReadOnly(runnerDB(credentials))
	if err != nil {
		return store.OperatorStepKey{}, target, "", err
	}
	defer db.Close()
	ctx, cancel := budgetContext(time.Second)
	defer cancel()
	key, signed, err := db.ReadOperatorWorktree(ctx, id)
	if err != nil {
		return key, target, "", err
	}
	if err := json.Unmarshal([]byte(key.Identity), &target); err != nil {
		return key, target, "", store.ErrOperatorConflict
	}
	path, err := operatorPhysicalTarget(operatorWorktreeFact(signed))
	if err != nil {
		return key, target, "", err
	}
	if signed.Store != credentials || signed.Kind != "start" || !signed.Options.Headless ||
		!verbs.EpisodeIDPattern.MatchString(signed.Episode) || filepath.Base(path) != signed.Episode ||
		target.Action != id || target.Path != path || target.Branch != "world/"+signed.Episode ||
		!filepath.IsAbs(target.Repo) || !filepath.IsAbs(target.Common) {
		return key, target, "", store.ErrOperatorConflict
	}
	return key, target, signed.Episode, nil
}

func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer,
	readReservation func(string, string) (store.OperatorStepKey, operatorWorktreeTarget, string, error)) int {
	remedy := reconcileWorktreeCommand(id)
	refuse := func(err error) int {
		return serviceMessage(stderr, "worktree recovery refused", err.Error(), remedy, exitUsage)
	}
	c, _, err := e.installedRunner()
	if err != nil {
		return refuse(err)
	}
	key, target, episode, err := readReservation(c.Credentials, id)
	if err != nil {
		return refuse(err)
	}
	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, episode), 0)
	if err != nil {
		return refuse(errors.New("episode recovery lock held or unavailable"))
	}
	defer lock.Close()
	current, _, currentEpisode, err := readReservation(c.Credentials, id)
	if err != nil {
		return refuse(err)
	}
	if current != key || currentEpisode != episode {
		return refuse(store.ErrOperatorConflict)
	}
	git, err := sessionGit()
	if err != nil {
		return refuse(err)
	}
	ctx, cancel := budgetContext(5 * time.Second)
	before, err := target.inspect(ctx, git)
	cancel()
	if err != nil {
		return refuse(err)
	}
	disposition := store.OperatorWorktreeRetained
	if before.state == treeAbsent {
		disposition = store.OperatorWorktreeAbsent
	}
	if !e.confirm(stderr, fmt.Sprintf("Stop action %s (%s)? Keep its worktree, branch and files as found.", id, disposition)) {
		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
		return exitUsage
	}
	ctx, cancel = budgetContext(5 * time.Second)
	after, err := target.inspect(ctx, git)
	cancel()
	if err != nil {
		return refuse(err)
	}
	if before != after {
		return refuse(errors.New("worktree evidence changed after review"))
	}
	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
		ctx, cancel := budgetContext(time.Second)
		defer cancel()
		return db.DisposeOperatorWorktree(ctx, key, disposition, e.terminal.now())
	})
	if err != nil {
		return refuse(err)
	}
	fmt.Fprintln(out, "Worktree reservation disposed for", id, "("+disposition+")")
	return exitOK
}
