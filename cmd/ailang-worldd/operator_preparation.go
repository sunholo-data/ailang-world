package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

var errOperatorDraining = errors.New("operator runner draining")

func (e serviceEnv) consumeOperatorPreparation(ctx context.Context, c runnerConfig, heartbeat store.OperatorHeartbeat,
	src operatorStartSource, stops <-chan os.Signal, process *os.File, recoveryCursor *string,
	acquire func(string, string) (*store.OperatorEpisodeLock, error), launch ...operatorPreparedLaunch) error {
	select {
	case <-stops:
		return errOperatorDraining
	default:
	}
	if err := e.recoverOperatorStart(ctx, c, heartbeat, src, recoveryCursor, stops, launch...); err != nil {
		return err
	}
	return e.prepareOperatorStart(ctx, c, heartbeat.Owner, heartbeat.Build, src, stops, process, acquire)
}

func (e serviceEnv) prepareOperatorStart(ctx context.Context, c runnerConfig, owner store.OperatorOwner, build string,
	src operatorStartSource, stops <-chan os.Signal, process *os.File,
	acquire func(string, string) (*store.OperatorEpisodeLock, error)) error {
	return withOperatorStartCandidate(ctx, c.Credentials, src.clock, store.OpenCredentialsReadOnly, acquire, func(held *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
		ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
		if err != nil {
			return err
		}
		defer ro.Close()
		h, err := ro.OperatorHeartbeat(ctx)
		if err != nil {
			return err
		}
		if !preparationHeartbeat(h, owner, build, h.CatalogueGeneration, src.clock()) {
			return store.ErrOperatorRunner
		}
		action, err := approval.DecodeOperatorAction(row.Action)
		if err != nil {
			return err
		}
		binary := operatorFileFact(action, "launch_binary")
		sha, err := authority.OperatorStartLaunchDigest(action)
		if err != nil {
			return err
		}
		if !slices.Contains(c.PresetNames, action.Options.Preset) || c.Backend != "headless" ||
			filepath.Base(binary) != "codex" || c.Executables[binary] != sha ||
			operatorFileFact(action, "served_address") != c.Address || len(row.Challenge) != 32 {
			return store.ErrOperatorConflict
		}
		if err := e.proveOldTerminalQuiet(); err != nil {
			return err
		}
		newEffect := func() error {
			select {
			case <-stops:
				return errOperatorDraining
			default:
			}
			if serviceFileHash(binary) != sha {
				return &approval.BindingError{Field: "launch_binary_sha256"}
			}
			if err := e.checkStartPair(c, src.workspace); err != nil {
				return err
			}
			live, err := ro.OperatorHeartbeat(ctx)
			if err != nil || !preparationHeartbeat(live, owner, build, h.CatalogueGeneration, src.clock()) {
				return store.ErrOperatorRunner
			}
			current, err := ro.OperatorAction(ctx, row.ID, src.clock())
			if err != nil {
				return err
			}
			_, err = authority.RecheckClaimedOperatorAction(ctx, ro, action, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), current)
			return err
		}
		attempt := store.OperatorMintAttempt{CAS: operatorPreparationCAS(row), Runner: owner, Catalogue: h.CatalogueGeneration}
		in, err := reconstructOperatorStart(ctx, action, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
		if err != nil {
			return err
		}
		if in.Start.Agent != "codex" || operatorFileFact(action, "mcp_path") != "" {
			return store.ErrOperatorConflict
		}
		if err := runnerPath(c.Home, filepath.Dir(in.TokenPath), c.UID, true, 0700); err != nil {
			return fmt.Errorf("private session directory: %w", err)
		}
		// Hold lifecycle exclusion across the last stop check and claim.
		// The exclusive episode description covers all later effects.
		admission, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
		if err != nil {
			return errOperatorDraining
		}
		err = newEffect()
		if err == nil {
			err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
				fresh, err := db.OperatorHeartbeat(ctx)
				if err != nil || !preparationHeartbeat(fresh, owner, build, h.CatalogueGeneration, src.clock()) {
					return store.ErrOperatorRunner
				}
				return db.ClaimOperatorAction(ctx, attempt.CAS, owner, src.clock())
			})
		}
		admission.Close()
		if err != nil {
			return err
		}
		attempt.CAS.Generation++
		checkpoint := func(name string) error {
			if strings.HasPrefix(name, "before ") && name != "before readback" && name != "before observe" {
				return newEffect()
			}
			return nil
		}
		if err := reserveOperatorStartWorktree(ctx, c.Credentials, held, action, attempt, src, checkpoint, process); err != nil {
			return err
		}
		if err := newEffect(); err != nil {
			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
		}
		row, err = ro.OperatorAction(ctx, row.ID, src.clock())
		if err != nil {
			return err
		}
		in.Attempt.CAS, in.Attempt.StepVersion = operatorPreparationCAS(row), 1
		var token string
		if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
			key := store.OperatorStepKey{Action: in.Attempt.CAS, Owner: row.Owner, Kind: "mint", Identity: row.ID, Version: 1}
			if err := db.ReserveOperatorStep(ctx, key, src.clock()); err != nil {
				return err
			}
			var err error
			token, err = authority.MintOperatorStart(ctx, db, ro, in, action, row, src.clock)
			return err
		}); err != nil {
			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
		}
		if err := publishOperatorStartFiles(ctx, c, held, row.ID, token, src.clock, checkpoint); err != nil {
			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
		}
		prepared, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
		if err != nil || prepared.Window != store.OperatorStartPrepared {
			return fmt.Errorf("%w; %s", store.ErrOperatorConflict, reconcileStartCommand(row.ID))
		}
		return nil
	})
}

func operatorPreparationCAS(row store.OperatorActionRow) store.OperatorCAS {
	return store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
}

func preparationHeartbeat(h store.OperatorHeartbeat, owner store.OperatorOwner, build string, catalogue, now int64) bool {
	return h.Owner == owner && h.CatalogueGeneration == catalogue && h.Available(build, "start", now)
}

func (e serviceEnv) checkStartPair(c runnerConfig, workspace string) error {
	installed, _, err := e.installedRunner()
	if err != nil || !reflect.DeepEqual(installed, c) {
		return store.ErrOperatorConflict
	}
	if err := e.checkRunnerPair(c); err != nil {
		return err
	}
	r, err := readServiceRecord(e.paths)
	if err != nil {
		return err
	}
	physical, err := filepath.EvalSymlinks(r.Workspace)
	current, currentErr := filepath.EvalSymlinks(workspace)
	if err != nil || currentErr != nil || physical != current {
		return store.ErrOperatorConflict
	}
	return nil
}

func reconcileStartCommand(id string) string {
	return "ailang-worldd service runner reconcile-start " + id
}

func recoveryWorktree(ctx context.Context, current store.OperatorStartRecovery, a approval.OperatorAction) (operatorWorktreeEvidence, error) {
	var target operatorWorktreeTarget
	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
	if err != nil || json.Unmarshal([]byte(current.Worktree), &target) != nil || target.Action != a.ID || target.Path != path {
		return operatorWorktreeEvidence{}, store.ErrOperatorConflict
	}
	git, err := sessionGit()
	if err != nil {
		return operatorWorktreeEvidence{}, err
	}
	seen, err := target.inspect(ctx, git)
	if err != nil || seen.state != treeExact {
		return operatorWorktreeEvidence{}, store.ErrOperatorConflict
	}
	return seen, nil
}
