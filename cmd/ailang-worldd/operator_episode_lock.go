package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

func operatorEpisodeLockPath(canonical, episode string) string {
	return store.OperatorEpisodeLockPath(canonical, episode)
}

func terminalEpisodeLock(creds, episode string, held *os.File) (*os.File, error) {
	return openTerminalEpisodeLock(creds, episode, held, os.O_CREATE)
}

func openTerminalEpisodeLock(creds, episode string, held *os.File, flags int, budget ...time.Duration) (*os.File, error) {
	canonical, err := store.CanonicalPath(creds)
	if err != nil {
		return nil, err
	}
	return openCanonicalTerminalEpisodeLock(canonical, episode, held, flags, budget...)
}

func openCanonicalTerminalEpisodeLock(canonical, episode string, held *os.File, flags int, budget ...time.Duration) (*os.File, error) {
	path := operatorEpisodeLockPath(canonical, episode)
	if held != nil {
		if held.Name() != path {
			return nil, errors.New("operator store or episode changed during episode start")
		}
		return nil, nil // start retains and closes the same open file description
	}

	// Take the same lock even before the first companion startup. A first
	// runner must not race a terminal operation that saw an unmanaged store.
	lock, lockErr := openRunnerLock(path, flags)
	absent := flags == 0 && errors.Is(lockErr, os.ErrNotExist)
	if lockErr != nil && !absent && !errors.Is(lockErr, syscall.EWOULDBLOCK) {
		return nil, lockErr
	}
	timeout := 4 * time.Second
	if len(budget) != 0 {
		timeout = budget[0]
	}
	ctx, cancel := budgetContext(timeout)
	defer cancel()
	id, readErr := store.ReadOperatorEpisodeOwner(ctx, canonical, episode)
	if (lockErr == nil || absent) && readErr == nil && id == "" {
		return lock, nil
	}
	if lock != nil {
		lock.Close()
	}

	switch {
	case readErr != nil || id != "" && !verbs.ValidOperatorActionID(id):
		return nil, fmt.Errorf("%w: could not read action ownership; run ailang-worldd service runner status", store.ErrOperatorConflict)
	case id != "":
		return nil, fmt.Errorf("%w: active action %s; run ailang-worldd service runner status and reconcile at the terminal", store.ErrOperatorConflict, id)
	default:
		return nil, fmt.Errorf("%w: episode lock held; run ailang-worldd service runner status", store.ErrOperatorConflict)
	}
}
