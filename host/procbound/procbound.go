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
