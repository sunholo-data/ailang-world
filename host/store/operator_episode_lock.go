package store

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// OperatorEpisodeLockPath requires a canonical credentials path and a validated
// episode. Terminal callers share this on-disk key with the runner and GC.
func OperatorEpisodeLockPath(canonical, episode string) string {
	return canonical + ".operator-episode-" + episode + ".lock"
}

// OperatorEpisodeLock attests to a non-waiting flock taken on this exact open
// description. The owner alone closes it. The sidecar is never unlinked.
type OperatorEpisodeLock struct{ file *os.File }

// OpenOperatorEpisodeLock takes the shared kernel lock without waiting.
func OpenOperatorEpisodeLock(credentials, episode string) (*OperatorEpisodeLock, error) {
	canonical, err := CanonicalPath(credentials)
	if err != nil {
		return nil, err
	}
	if episode == "" || strings.ContainsAny(episode, "/\\\x00") {
		return nil, ErrOperatorConflict
	}
	path := OperatorEpisodeLockPath(canonical, episode)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}
	lock := &OperatorEpisodeLock{f}
	if _, err := lock.Verify(canonical, episode); err != nil {
		lock.Close()
		return nil, err
	}
	return lock, nil
}

func (l *OperatorEpisodeLock) Close() error { return l.file.Close() }

// Verify lends the exact descriptor for a synchronous effect and, if needed,
// for exec inheritance. Its borrower must never close or unlock the file.
func (l *OperatorEpisodeLock) Verify(canonical, episode string) (*os.File, error) {
	return l.verify(canonical, episode, os.Getuid())
}

func (l *OperatorEpisodeLock) verify(canonical, episode string, uid int) (*os.File, error) {
	path := OperatorEpisodeLockPath(canonical, episode)
	if l == nil || l.file == nil || l.file.Name() != path {
		return nil, fmt.Errorf("%w: wrong episode lock", ErrOperatorConflict)
	}
	opened, err := l.file.Stat()
	if err != nil {
		return nil, fmt.Errorf("%w: closed episode lock: %w", ErrOperatorConflict, err)
	}
	current, err := os.Lstat(path)
	if err != nil || !os.SameFile(opened, current) {
		return nil, fmt.Errorf("%w: episode lock pathname changed", ErrOperatorConflict)
	}
	stat := opened.Sys().(*syscall.Stat_t)
	if !opened.Mode().IsRegular() || opened.Mode().Perm() != 0600 || int(stat.Uid) != uid ||
		stat.Nlink != 1 {
		return nil, fmt.Errorf("%w: unsafe episode lock", ErrOperatorConflict)
	}
	return l.file, nil
}
