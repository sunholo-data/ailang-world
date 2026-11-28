package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const (
	runnerPoll                = 2 * time.Second
	operatorPreparationBudget = 22 * time.Second
)

type operatorBudgetKey struct{}

func operatorChildBudget(ctx context.Context, budget time.Duration) (context.Context, context.CancelFunc) {
	if override, ok := ctx.Value(operatorBudgetKey{}).(time.Duration); ok {
		budget = override
	}
	return context.WithTimeout(ctx, budget)
}

func (e serviceEnv) operatorContext(budget time.Duration) (result context.Context, done context.CancelFunc) {
	if e.operatorDeadline != nil {
		defer func() { e.operatorDeadline(budget, result) }()
	}
	if e.operatorBudget == 0 {
		return budgetContext(budget)
	}
	ctx, cancel := budgetContext(e.operatorBudget)
	return context.WithValue(worktree.WithCommandTimeout(ctx, e.operatorBudget), operatorBudgetKey{}, e.operatorBudget), cancel
}

func runnerKernelLock(path string) (*os.File, error) {
	return openRunnerLock(path, os.O_CREATE)
}

func openRunnerLock(path string, flags int) (*os.File, error) {
	f, err := os.OpenFile(path, flags|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Mode().Perm() != 0600 || int(st.Sys().(*syscall.Stat_t).Uid) != os.Getuid() {
		f.Close()
		return nil, errors.New("unsafe runner lock file")
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func runOperatorRunner(args []string, errOut io.Writer, e serviceEnv) int {
	if len(args) != 0 {
		return exitUsage
	}
	c, b, err := e.installedRunner()
	if err != nil || os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b) ||
		daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary] {
		return serviceMessage(errOut, "runner startup refused", "installation or executable identity mismatch", "ailang-worldd service runner status", 2)
	}
	if err := e.checkRunnerPair(c); err != nil {
		return serviceMessage(errOut, "runner startup refused", err.Error(), "ailang-worldd service runner status", 2)
	}
	lock, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
	if err != nil {
		return serviceMessage(errOut, "runner startup refused", "cannot acquire the kernel runner lock", "ailang-worldd service runner status", 2)
	}
	defer lock.Close()
	h := store.OperatorHeartbeat{
		Owner: store.OperatorOwner{Nonce: rand.Text()}, Protocol: store.OperatorRunnerProtocol,
		Build: c.Executables[e.paths.Binary],
	}
	stops := make(chan os.Signal, 1)
	signal.Notify(stops, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(stops)
	ticker := time.NewTicker(runnerPoll)
	defer ticker.Stop()
	source := operatorStartSource{defaultPresetsDir(), defaultExecProfilesDir(), "", defaultHostsDir(), e.terminal.now}
	if len(c.FinishRemotes) != 0 {
		ctx, cancel := e.operatorContext(30 * time.Second)
		h.Readiness = "both"
		for _, remote := range c.FinishRemotes {
			if err := prepareOperatorFinish(ctx, c, remote); err != nil {
				fmt.Fprintln(errOut, err)
				h.Readiness = "start"
				break
			}
		}
		cancel()
	}
	e.operatorCatalogue = func(c runnerConfig) []byte { return e.readOperatorCatalogue(c, source) }
	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
}

func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig, h store.OperatorHeartbeat,
	src operatorStartSource, stops <-chan os.Signal, ticks <-chan time.Time, lock *os.File) int {
	var recoveryCursor string
	for {
		select {
		case <-stops:
			return 0
		default:
		}

		if err := e.runnerTick(c, &h); err != nil {
			return serviceMessage(errOut, "runner halted", "pair or heartbeat unavailable", "ailang-worldd service runner status", 2)
		}
		record, err := readServiceRecord(e.paths)
		if err != nil {
			return serviceMessage(errOut, "runner halted", "service record unavailable", "ailang-worldd service runner status", 2)
		}
		ctx, cancel := e.operatorContext(operatorPreparationBudget)
		src.workspace = record.Workspace
		launch := func(ro *store.Credentials, row store.OperatorActionRow, current store.OperatorStartRecovery, episode *store.OperatorEpisodeLock) error {
			err := e.launchPreparedOperatorStart(ctx, c, h, src, stops, lock, errOut, ro, row, current, episode)
			if errors.Is(err, errOperatorDraining) {
				return err
			}
			if err != nil {
				fmt.Fprintf(errOut, "operator launch %s: %s\n", row.ID, operatorLaunchFailure(err))
			}
			return errOperatorLaunchHandled
		}
		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock, launch)
		cancel()
		switch {
		case errors.Is(err, errOperatorDraining):
			return 0
		case err != nil && !errors.Is(err, store.ErrOperatorMissing) && !errors.Is(err, errOperatorLaunchHandled):
			var stale *approval.BindingError
			var remedy operatorRecoveryNeeded
			switch {
			case errors.As(err, &stale):
				fmt.Fprintf(errOut, "operator start binding stale (%s); request fresh review\n", stale.Field)
			case errors.As(err, &remedy):
				fmt.Fprintf(errOut, "operator start blocked; attended recovery: %s\n", remedy.command)
			default:
				fmt.Fprintln(errOut, "operator start blocked; inspect service runner status and reconcile durable evidence")
			}
		}
		if h.Readiness == "both" {
			err := e.runFinishTurn(c, &h, src, stops, ticks)
			var auth operatorGitHubAuth
			var busy operatorFinishBusy
			switch {
			case errors.Is(err, errOperatorDraining):
				return 0
			case errors.As(err, &auth):
				h.Readiness = "start"
				if err := e.touchRunner(c, &h); err != nil {
					return serviceMessage(errOut, "runner halted", "heartbeat unavailable", "ailang-worldd service runner status", 2)
				}
				fmt.Fprintln(errOut, auth)
			case errors.As(err, &busy):
				fmt.Fprintln(errOut, busy.Error())
			case err != nil && !errors.Is(err, store.ErrOperatorMissing):
				fmt.Fprintln(errOut, finishRefusal(err))
			}
		}
		select {
		case _, open := <-ticks:
			if !open {
				return 0
			}
		case <-stops:
			return 0
		}
	}
}

func finishRefusal(err error) string {
	for _, class := range []error{errOperatorFinishReview, errOperatorFinishEligibility, errOperatorFinishTerminal,
		errOperatorFinishFence, errOperatorFinishCheck, errOperatorFinishUncertain, errOperatorCommitFailed} {
		if errors.Is(err, class) {
			return class.Error()
		}
	}
	return "operator finish blocked; inspect action status and durable evidence"
}

func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error {
	if err := e.checkRunnerConfig(c); err != nil {
		return err
	}
	if err := e.checkRunnerPair(c); err != nil {
		return err
	}
	if h.Readiness == "" {
		h.Readiness = "start"
	}
	if e.operatorCatalogue != nil {
		h.Catalogue = e.operatorCatalogue(c)
	} else {
		h.Catalogue, _ = json.Marshal(c.PresetNames)
	}
	return e.touchRunner(c, h)
}

func (e serviceEnv) touchRunner(c runnerConfig, h *store.OperatorHeartbeat) error {
	h.SeenAt = e.terminal.now()
	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
		ctx, cancel := e.operatorContext(time.Second)
		defer cancel()
		generation, err := db.PutOperatorHeartbeat(ctx, *h)
		if err == nil {
			h.Owner.Generation = generation
		}
		return err
	})
}

func runnerHeartbeat(path string, budget time.Duration) (store.OperatorHeartbeat, error) {
	ctx, cancel := budgetContext(budget)
	defer cancel()
	return store.ReadOperatorHeartbeat(ctx, path)
}

func (e serviceEnv) waitForRunner(c runnerConfig, previous int64, read func(string, time.Duration) (store.OperatorHeartbeat, error)) bool {
	deadline := e.runnerNow().Add(e.wait)
	for {
		remaining := deadline.Sub(e.runnerNow())
		if remaining <= 0 {
			return false
		}
		h, err := read(c.Credentials, min(remaining, time.Second))
		now := e.runnerNow()
		if !now.Before(deadline) {
			return false
		}
		if err == nil && h.Owner.Generation > previous &&
			h.Available(c.Executables[e.paths.Binary], "start", now.Unix()) {
			return true
		}
		e.runnerSleep(min(deadline.Sub(now), 100*time.Millisecond))
	}
}

func (e serviceEnv) runnerStatus(w io.Writer) error {
	c, _, err := e.installedRunner()
	if err != nil {
		return err
	}
	db, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return errors.New("credentials disconnected")
	}
	defer db.Close()
	ctx, cancel := budgetContext(time.Second)
	defer cancel()
	h, err := db.OperatorHeartbeat(ctx)
	if err != nil {
		return errors.New("runner heartbeat unavailable")
	}
	fmt.Fprintf(w, "runner last seen %d, generation %d\n", h.SeenAt, h.Owner.Generation)
	if err := e.checkRunnerPair(c); err != nil {
		return fmt.Errorf("daemon disconnected: %w", err)
	}
	if _, err := e.controller.Print(e.domain, runnerLabel); err != nil {
		return errors.New("runner disconnected")
	}
	build, now := c.Executables[e.paths.Binary], e.terminal.now()
	if !h.Available(build, "start", now) {
		return errors.New("runner heartbeat stale or incompatible")
	}
	fmt.Fprintln(w, "start: ready")
	if !h.Available(build, "finish", now) {
		fmt.Fprintln(w, "finish blocked: github-auth; "+operatorFinishRemedy)
	} else {
		fmt.Fprintln(w, "finish: ready")
	}
	return nil
}

func runnerWrite(path string, write func(*store.Credentials) error) error {
	db, err := store.OpenCredentials(runnerDB(path))
	if err != nil {
		return err
	}
	defer db.Close()
	return write(db)
}

func runnerCancel(c runnerConfig, id string, now int64) error {
	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
		ctx, cancel := budgetContext(time.Second)
		defer cancel()
		return db.CancelUnstartedOperatorAction(ctx, id, now)
	})
}

// Bound schema reads on open and require an existing database for writes.
func runnerDB(path string) string {
	u := url.URL{Scheme: "file", Path: path, RawQuery: "mode=rw&_pragma=busy_timeout%28250%29"}
	return u.String()
}
