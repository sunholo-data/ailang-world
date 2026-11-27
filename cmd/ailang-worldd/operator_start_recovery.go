package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

type operatorRecoveryNeeded struct{ command string }

func (r operatorRecoveryNeeded) Error() string { return r.command }

type operatorPreparedLaunch func(*store.Credentials, store.OperatorActionRow, store.OperatorStartRecovery, *store.OperatorEpisodeLock) error

func (e serviceEnv) recoverOperatorStart(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource, cursor *string, stops <-chan os.Signal, launch ...operatorPreparedLaunch) error {
	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return err
	}
	defer ro.Close()
	row, err := ro.NextOwnedOperatorStart(ctx, *cursor)
	if errors.Is(err, store.ErrOperatorMissing) {
		*cursor = ""
		return nil
	}
	if err != nil {
		return err
	}
	*cursor = row.ID
	return e.recoverOwnedStart(ctx, ro, c, h, src, row, stops, launch)
}

func (e serviceEnv) recoverOwnedStart(ctx context.Context, ro *store.Credentials, c runnerConfig, h store.OperatorHeartbeat,
	src operatorStartSource, row store.OperatorActionRow, stops <-chan os.Signal, launch []operatorPreparedLaunch) error {
	if !verbs.EpisodeIDPattern.MatchString(row.Episode) {
		return store.ErrOperatorConflict
	}
	if _, err := ro.ReadOperatorLaunch(ctx, row.ID); !errors.Is(err, store.ErrOperatorMissing) {
		if err != nil {
			return err
		}
		return operatorRecoveryNeeded{"ailang-worldd service runner launch-status " + row.ID}
	}
	lock, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
	if err != nil {
		return fmt.Errorf("%w: live episode owner or effect", store.ErrOperatorConflict)
	}
	defer lock.Close()
	current, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
	if err != nil || current.CAS != operatorPreparationCAS(row) || current.Owner != row.Owner || current.Episode != row.Episode {
		return store.ErrOperatorConflict
	}
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil || a.Store != c.Credentials {
		return store.ErrOperatorConflict
	}
	// A heartbeat can make an old review ineligible. It cannot prove an
	// episode process has exited. The runner and episode locks cover this.
	switch current.Window {
	case store.OperatorStartGitPending:
		return operatorRecoveryNeeded{reconcileWorktreeCommand(row.ID)}
	case store.OperatorStartMintOnly:
		return operatorRecoveryNeeded{"ailang-worldd service runner reconcile-mint " + row.ID}
	case store.OperatorStartFilesPending:
		return operatorRecoveryNeeded{reconcileFilesCommand(row.ID)}
	case store.OperatorStartTreeOnly, store.OperatorStartMintPending, store.OperatorStartMintUncertain:
		return operatorRecoveryNeeded{reconcileStartCommand(row.ID)}
	case store.OperatorStartPrepared:
		if len(launch) != 0 {
			return launch[0](ro, row, current, lock)
		}
		return verifyPreparedStart(ctx, c, current, ro, a, src.clock())
	case store.OperatorStartUntouched:
		if err := e.proveOldTerminalQuiet(); err != nil {
			return err
		}
		catalogue, err := strconv.ParseInt(operatorFileFact(a, "mint_catalogue"), 10, 64)
		if err != nil {
			return &approval.BindingError{Field: "mint_catalogue"}
		}
		live, err := ro.OperatorHeartbeat(ctx)
		if err != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
			return &approval.BindingError{Field: "mint_catalogue"}
		}
		if err := e.checkStartPair(c, src.workspace); err != nil {
			return err
		}
		if _, err := authority.OperatorStartLaunchDigest(a); err != nil {
			return err
		}
		p, err := authority.RecheckClaimedOperatorAction(ctx, ro, a, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), row)
		if err != nil {
			return err
		}
		attempt := store.OperatorMintAttempt{CAS: current.CAS, Runner: h.Owner, Catalogue: catalogue}
		in, err := reconstructOperatorStart(ctx, a, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
		if err != nil {
			return err
		}
		git, err := sessionGit()
		if err != nil {
			return err
		}
		target, err := newOperatorWorktreeTarget(ctx, git, a, in.Start)
		if err != nil {
			return err
		}
		tree, err := target.inspect(ctx, git)
		if err != nil || tree.state != treeAbsent {
			return store.ErrOperatorConflict
		}
		for _, path := range []string{in.TokenPath, in.Start.Task.Path, operatorFileFact(a, "mcp_path")} {
			if path != "" {
				if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
					return store.ErrOperatorConflict
				}
			}
		}
		admit, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
		if err != nil {
			return errOperatorDraining
		}
		defer admit.Close()
		select {
		case <-stops:
			return errOperatorDraining
		default:
		}
		return runnerWrite(c.Credentials, func(db *store.Credentials) error {
			live, err := db.OperatorHeartbeat(ctx)
			if err != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
				return store.ErrOperatorRunner
			}
			return db.DisposeOperatorStartGap(ctx, current, p.OperatorConfirmation().MintProof(), src.clock())
		})
	}
	return store.ErrOperatorConflict
}

func verifyPreparedStart(ctx context.Context, c runnerConfig, current store.OperatorStartRecovery,
	ro *store.Credentials, a approval.OperatorAction, now int64) error {
	s, _, err := ro.ReadOperatorFiles(ctx, a.ID)
	if err != nil || s.Phase != "observed" || s.Mint.CAS != current.CAS || s.Worktree != current.Worktree ||
		now >= s.Session.ExpiresAt || store.ValidateOperatorStartSession(s, a) != nil {
		return store.ErrOperatorConflict
	}
	if _, err := recoveryWorktree(ctx, current, a); err != nil {
		return err
	}
	files, err := readStartRecoveryFiles(c, s, a)
	if err != nil || files.Files[0].Status != "exact" || files.Files[1].Status != "exact" {
		return store.ErrOperatorConflict
	}
	return nil
}

func readStartRecoveryFiles(c runnerConfig, s store.OperatorFileState, a approval.OperatorAction) (operatorFilesProof, error) {
	want, start, err := operatorFilePaths(s, a)
	var recorded operatorFileSet
	if err != nil || json.Unmarshal([]byte(s.Step.Identity), &recorded) != nil || recorded.Inode == 0 {
		return operatorFilesProof{}, store.ErrOperatorConflict
	}
	want.Device, want.Inode = recorded.Device, recorded.Inode
	if recorded != want {
		return operatorFilesProof{}, store.ErrOperatorConflict
	}
	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
	if err != nil {
		return operatorFilesProof{}, err
	}
	files := readOperatorFileSet(recorded, addr)
	if files.Dir.Status != "private directory" {
		return operatorFilesProof{}, store.ErrOperatorConflict
	}
	return files, nil
}

type operatorRecoveryObserved struct {
	tree    operatorWorktreeEvidence
	files   operatorFilesProof
	session store.OperatorFileState
}

func (e serviceEnv) reconcileOperatorStart(id string, out, stderr io.Writer) int {
	refuse := func(err error) int {
		return serviceMessage(stderr, "start recovery refused", err.Error(), "stop the companion and inspect durable evidence before retrying "+reconcileStartCommand(id), exitUsage)
	}
	lifecycle, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
	if err != nil {
		return refuse(err)
	}
	defer lifecycle.Close()
	c, _, err := e.installedRunner()
	if err != nil {
		return refuse(err)
	}
	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
		return refuse(errors.New("companion must be unloaded"))
	}
	process, err := e.waitRunnerExit(c)
	if err != nil {
		return refuse(err)
	}
	defer process.Close()
	if err := e.proveOldTerminalQuiet(); err != nil {
		return refuse(err)
	}
	ctx, cancel := e.operatorContext(5 * time.Second)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return refuse(err)
	}
	defer ro.Close()
	row, err := ro.OperatorAction(ctx, id, 0)
	if err != nil || !verbs.EpisodeIDPattern.MatchString(row.Episode) {
		return refuse(store.ErrOperatorConflict)
	}
	episode, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
	if err != nil {
		return refuse(err)
	}
	defer episode.Close()
	reviewed, err := ro.ReadOperatorStartRecovery(ctx, id)
	if err != nil || reviewed.CAS != operatorPreparationCAS(row) || reviewed.Owner != row.Owner {
		return refuse(store.ErrOperatorConflict)
	}
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil || a.Store != c.Credentials {
		return refuse(store.ErrOperatorConflict)
	}
	check := func(ctx context.Context) (operatorRecoveryObserved, error) {
		var seen operatorRecoveryObserved
		switch reviewed.Window {
		case store.OperatorStartUntouched:
			return seen, nil
		case store.OperatorStartTreeOnly, store.OperatorStartMintPending, store.OperatorStartMintUncertain, store.OperatorStartFilesPending:
			tree, err := recoveryWorktree(ctx, reviewed, a)
			if err != nil {
				return seen, err
			}
			seen.tree = tree
			if reviewed.Window != store.OperatorStartFilesPending {
				return seen, nil
			}
			seen.session, _, err = ro.ReadOperatorFiles(ctx, id)
			if err != nil {
				return seen, err
			}
			seen.files, err = readStartRecoveryFiles(c, seen.session, a)
			if err != nil || seen.files.Files[1].Status != "absent" ||
				(seen.files.Files[0].Status != "absent" && seen.files.Files[0].Status != "exact") {
				return seen, store.ErrOperatorConflict
			}
			return seen, nil
		default:
			return seen, store.ErrOperatorConflict
		}
	}
	before, err := check(ctx)
	if err != nil {
		return refuse(err)
	}
	cancel()
	change := fmt.Sprintf("Stop action %s (%s) and retain its existing evidence?", id, reviewed.Window)
	if reviewed.Window == store.OperatorStartFilesPending {
		change = fmt.Sprintf("Stop %s and revoke orphan credential %s, preserving all files and the worktree?", id, before.session.Session.CredentialID)
	}
	if reviewed.Window == store.OperatorStartMintUncertain {
		change = fmt.Sprintf("Mark %s uncertain and retain its owner and all evidence?", id)
	}
	if !e.confirm(stderr, change) {
		return exitUsage
	}
	ctx, done := e.operatorContext(5 * time.Second)
	defer done()
	after, err := check(ctx)
	if err != nil || before != after {
		return refuse(store.ErrOperatorConflict)
	}
	if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
		switch reviewed.Window {
		case store.OperatorStartFilesPending:
			return db.RecordOperatorFiles(ctx, before.session, store.OperatorFilesStoppedOrphan, nil, e.terminal.now)
		case store.OperatorStartMintUncertain:
			return db.MarkOperatorMintUncertain(ctx, reviewed)
		default:
			return db.DisposeOperatorStartGap(ctx, reviewed, nil, e.terminal.now())
		}
	}); err != nil {
		return refuse(err)
	}
	if reviewed.Window == store.OperatorStartMintUncertain {
		fmt.Fprintln(out, "Mint commit uncertain; owner retained:", id)
		return exitOK
	}
	fmt.Fprintln(out, "Start action stopped; existing evidence retained:", id)
	return exitOK
}
