package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

func (e serviceEnv) reconcileRunnerMint(id string, out, stderr io.Writer,
	acquire func(string, int) (*os.File, error)) int {
	refuse := func(err error) int {
		return serviceMessage(stderr, "mint recovery refused", err.Error(), "ailang-worldd service runner status", 1)
	}
	c, _, err := e.installedRunner()
	if err != nil {
		return refuse(err)
	}
	r, err := readRunnerMint(c.Credentials, id)
	if err != nil {
		return refuse(err)
	}
	lock, err := acquire(operatorEpisodeLockPath(c.Credentials, r.Episode), 0)
	if err != nil {
		return refuse(errors.New("episode recovery lock held or unavailable"))
	}
	defer lock.Close()
	current, err := readRunnerMint(c.Credentials, id)
	if err != nil {
		return refuse(err)
	}
	if current.Episode != r.Episode {
		return refuse(store.ErrOperatorConflict)
	}
	r = current
	present, err := inspectRunnerMintToken(c, r)
	if err != nil {
		return refuse(err)
	}
	if present {
		return runnerMintRetained(out, stderr, r, e.terminal.now())
	}
	remedy := "ailang-worldd service runner reconcile-mint " + id
	if !e.confirm(stderr, fmt.Sprintf("Revoke orphaned mint %s for action %s and release its owner? Start and worktree are retained.", r.CredentialID, id)) {
		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
		return exitUsage
	}
	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
		ctx, cancel := budgetContext(time.Second)
		defer cancel()
		present, err := inspectRunnerMintToken(c, r)
		if err != nil {
			return err
		}
		if present {
			return errors.New("token appeared after review; inspect it again")
		}
		return db.AbortOperatorMint(ctx, r, e.terminal.now())
	})
	if err != nil {
		return refuse(err)
	}
	fmt.Fprintln(out, "Orphaned mint revoked for", id)
	return exitOK
}

func readRunnerMint(path, id string) (store.OperatorMintReceipt, error) {
	db, err := store.OpenCredentialsReadOnly(runnerDB(path))
	if err != nil {
		return store.OperatorMintReceipt{}, err
	}
	defer db.Close()
	ctx, cancel := budgetContext(time.Second)
	defer cancel()
	return db.ReadOperatorMint(ctx, id)
}

func runnerMintRetained(out, stderr io.Writer, r store.OperatorMintReceipt, now int64) int {
	if now >= r.ExpiresAt {
		return serviceMessage(stderr, "mint recovery refused", "matching token expired", "ailang-worldd service runner status", 1)
	}
	fmt.Fprintln(out, "Matching private token retained for", r.CAS.ID, "(owner retained for file reconciliation)")
	return exitOK
}

func inspectRunnerMintToken(c runnerConfig, r store.OperatorMintReceipt) (bool, error) {
	path := r.TokenPath
	if r.Store != c.Credentials || !verbs.EpisodeIDPattern.MatchString(r.Episode) ||
		!filepath.IsAbs(path) || filepath.Clean(path) != path || filepath.Base(path) != r.Episode+".session" {
		return false, errors.New("mint token path or store identity differs")
	}
	if err := runnerPath(c.Home, filepath.Dir(path), c.UID, true, 0700); err != nil {
		return false, fmt.Errorf("cannot establish private token directory: %w", err)
	}
	return readMintToken(path, c.UID, r.CredentialID)
}

func readMintToken(path string, uid int, credentialID string) (bool, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, errors.New("cannot open mint token safely")
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return false, errors.New("mint token is not a regular file")
	}
	stat := info.Sys().(*syscall.Stat_t)
	if int(stat.Uid) != uid || stat.Nlink != 1 || info.Mode() != 0600 {
		return false, errors.New("mint token must be private, owned and unlinked")
	}
	body, err := io.ReadAll(io.LimitReader(f, 65))
	if err != nil || len(body) != 64 || fmt.Sprintf("%x", sha256.Sum256(body)) != credentialID {
		return false, errors.New("mint token bytes differ from receipt")
	}
	return true, nil
}
