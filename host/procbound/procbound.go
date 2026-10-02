// Package procbound bounds the last step of a subprocess call: waiting for the
// direct child to exit. A process-group SIGKILL makes that wait short in every
// case the kill reaches; when the kill fails (or its target does not die) the
// wait must still return, so the call can report the failure instead of
// hanging. The abandoned wait keeps running in the background and reaps the
// child whenever it exits.
//
// Every subprocess holds one reservation from Admit (before cmd.Start) until
// its direct child is reaped. MaxOutstanding bounds active children plus
// retained background waiters process-wide, exactly: Admit is a nonblocking
// compare-and-swap reservation, not an observation.
package procbound

import (
	"errors"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ErrCleanupIncomplete means the direct child had not exited by the cleanup
// deadline. A background waiter still owns its eventual reap.
var ErrCleanupIncomplete = errors.New("subprocess cleanup incomplete: child not reaped by the cleanup deadline")

// ErrCleanupBacklog refuses a new subprocess while MaxOutstanding reservations
// (active children plus children that survived their kill) are held.
var ErrCleanupBacklog = errors.New("subprocess refused: too many unreaped children outstanding")

// MaxOutstanding is the process-wide limit on reservations: admitted children
// not yet reaped, including those whose reap a background waiter owns.
const MaxOutstanding = 8

var outstanding atomic.Int64

// Outstanding reports the reservations currently held.
func Outstanding() int64 { return outstanding.Load() }

// Admit reserves one slot before cmd.Start, or refuses at once with
// ErrCleanupBacklog when MaxOutstanding are held; it never waits. The returned
// release is idempotent: the first call frees the slot, later calls do nothing.
// The caller releases on Start failure; otherwise it hands release to Wait.
func Admit() (release func(), err error) {
	for {
		n := outstanding.Load()
		if n >= MaxOutstanding {
			return nil, ErrCleanupBacklog
		}
		if outstanding.CompareAndSwap(n, n+1) {
			break
		}
	}
	var once sync.Once
	return func() { once.Do(func() { outstanding.Add(-1) }) }, nil
}

// Wait runs wait (the direct child's reap). If it finishes within d, Wait
// releases the reservation and returns its result. Otherwise Wait returns
// ErrCleanupIncomplete and a background goroutine keeps the reservation until
// wait returns, then releases it: the slot stays counted until the reap.
func Wait(wait func() error, d time.Duration, release func()) error {
	done := make(chan error, 1)
	go func() { done <- wait() }()
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case err := <-done:
		release()
		return err
	case <-timer.C:
	}
	go func() {
		<-done
		release()
	}()
	return ErrCleanupIncomplete
}

// NothingToKill reports whether a failed process-group kill found nothing to
// kill, so it must not be reported as a cleanup failure. Call it only after
// the direct child's wait, so the reaped leader no longer holds the group.
//   - ESRCH: the group was already empty.
//   - Exactly EPERM (the kill's own bare errno), and groupEmpty() now: darwin
//     answers kill(-pgid, SIGKILL) with EPERM, not ESRCH, when every member of
//     the group is an unreaped zombie (queue row 115, measured), and an empty
//     group after the reap shows no member outlived the call. A wrapped or
//     joined EPERM carries more than that answer and is still reported.
func NothingToKill(killErr error, groupEmpty func() bool) bool {
	if errors.Is(killErr, syscall.ESRCH) {
		return true
	}
	return killErr == syscall.EPERM && groupEmpty()
}

// GroupEmpty reports whether process group pgid has no members left:
// kill(-pgid, 0) answers ESRCH. Signal 0 sends nothing; any other answer (a
// live, unsignallable or zombie member) reports false, as does pgid <= 1.
func GroupEmpty(pgid int) bool {
	if pgid <= 1 {
		return false
	}
	return errors.Is(syscall.Kill(-pgid, 0), syscall.ESRCH)
}
