package broker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/procbound"
)

const (
	handlerExecTimeout    = 30 * time.Second
	maxHandlerOutputBytes = 8 << 20
)

var (
	ErrHandlerTimeout  = errors.New("broker: handler subprocess timed out")
	ErrHandlerOverflow = errors.New("broker: handler subprocess output overflow")
)

type handlerBounds struct {
	execTimeout    time.Duration
	maxOutputBytes int64
}

func (b handlerBounds) normalized() handlerBounds {
	if b.execTimeout <= 0 {
		b.execTimeout = handlerExecTimeout
	}
	if b.maxOutputBytes <= 0 {
		b.maxOutputBytes = maxHandlerOutputBytes
	}
	return b
}

// HandlerTimeoutError reports expiry of the shared subprocess wall-clock bound.
type HandlerTimeoutError struct {
	Timeout time.Duration
}

func (e *HandlerTimeoutError) Error() string {
	return fmt.Sprintf("%v after %s", ErrHandlerTimeout, e.Timeout)
}

func (e *HandlerTimeoutError) Unwrap() error { return ErrHandlerTimeout }

// HandlerOutputOverflowError reports output larger than the shared byte bound.
type HandlerOutputOverflowError struct {
	Limit int64
}

func (e *HandlerOutputOverflowError) Error() string {
	return fmt.Sprintf("%v: limit %d bytes", ErrHandlerOverflow, e.Limit)
}

func (e *HandlerOutputOverflowError) Unwrap() error { return ErrHandlerOverflow }

// HandlerExitError preserves a subprocess's exit status and bounded diagnostic output.
type HandlerExitError struct {
	Err    error
	Output []byte
}

func (e *HandlerExitError) Error() string {
	return fmt.Sprintf("broker: handler subprocess failed: %v (output %q)", e.Err, e.Output)
}

func (e *HandlerExitError) Unwrap() error { return e.Err }

type handlerCommand struct {
	path string
	args []string
	dir  string
	env  []string
	// stdin, when non-nil, is the child's whole standard input.
	stdin []byte
	// stderr, when non-nil, receives the child's stderr SEPARATELY from the
	// returned stdout capture, under the same byte bound; nil keeps the
	// historical merged capture every pre-existing handler relies on.
	stderr *boundedSink
}

// boundedSink is a concurrency-safe stderr capture holding at most limit
// bytes. It keeps accepting (and discarding) writes past the limit so the
// child never blocks on a full pipe, and records the overflow, which
// runBounded reports as *HandlerOutputOverflowError — never a silent
// truncation.
type boundedSink struct {
	mu       sync.Mutex
	buf      []byte
	limit    int64
	overflow bool
}

func newBoundedSink(limit int64) *boundedSink { return &boundedSink{limit: limit} }

func (s *boundedSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := s.limit - int64(len(s.buf))
	if int64(len(p)) > room {
		if room > 0 {
			s.buf = append(s.buf, p[:room]...)
		}
		s.overflow = true
		return len(p), nil
	}
	s.buf = append(s.buf, p...)
	return len(p), nil
}

// Bytes returns a copy of the captured bytes.
func (s *boundedSink) Bytes() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]byte(nil), s.buf...)
}

func (s *boundedSink) overflowed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.overflow
}

// killGroup is the cancellation kill boundary, a package-level seam so the
// timeout tests can observe the kill's time, target and errno directly.
var killGroup = func(pgid int) error {
	return syscall.Kill(-pgid, syscall.SIGKILL)
}

// pipeCloseGrace is how long after the execTimeout the parent's read end is
// closed. exec.Cmd.WaitDelay cannot bound this: Wait runs after the read.
const pipeCloseGrace = time.Second

// runBounded is the one timeout and allocation surface used by every subprocess
// handler. It reads at most limit+1 bytes so overflow is detected, never hidden
// by truncation.
func runBounded(ctx context.Context, bounds handlerBounds, spec handlerCommand) ([]byte, error) {
	bounds = bounds.normalized()
	runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, spec.path, spec.args...)
	// The bound must survive a handler that forks. Killing only the direct child
	// leaves a grandchild holding the inherited stdout pipe, so the capped read
	// below blocks until the GRANDCHILD exits and the deadline is not enforced —
	// observed on linux CI as 5.002s against a 40ms bound, exactly the runtime of
	// the grandchild. Run the handler in its own process group and kill the whole
	// group, which is the same correction scripts/verify_ail.sh already carries
	// (V26) for the same reason.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return killGroup(cmd.Process.Pid)
	}
	cmd.Dir = spec.dir
	cmd.Env = spec.env
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("broker: handler stdout pipe: %w", err)
	}
	if spec.stdin != nil {
		cmd.Stdin = bytes.NewReader(spec.stdin)
	}
	if spec.stderr != nil {
		if spec.stderr.limit <= 0 {
			spec.stderr.limit = bounds.maxOutputBytes
		}
		cmd.Stderr = spec.stderr
		// A separate stderr is copied by an exec-owned goroutine that Wait
		// joins; a descendant holding that pipe must not stall Wait, so the
		// pipes are force-closed a grace period after the child exits.
		cmd.WaitDelay = pipeCloseGrace
	} else {
		cmd.Stderr = cmd.Stdout
	}
	release, err := procbound.Admit()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		release()
		// Start refuses an already-expired context with its error; under load
		// the allowance can run out before the child exists. That is still
		// the wall-clock bound, not a start failure.
		if errors.Is(err, context.DeadlineExceeded) && errors.Is(runCtx.Err(), context.DeadlineExceeded) {
			return nil, &HandlerTimeoutError{Timeout: bounds.execTimeout}
		}
		return nil, fmt.Errorf("broker: start handler subprocess: %w", err)
	}
	// Wait is bounded too, so a failed kill cannot hang the call: the child is
	// left to a background reaper and reported as ErrCleanupIncomplete.
	cleanupDeadline := time.Now().Add(bounds.execTimeout + 2*pipeCloseGrace)
	wait := func() error { return procbound.Wait(cmd.Wait, time.Until(cleanupDeadline), release) }
	// A descendant that left the group (setsid) survives the group kill and
	// holds the pipe; close the read end after the bound so ReadAll returns.
	closer := time.AfterFunc(bounds.execTimeout+pipeCloseGrace, func() { _ = pipe.Close() })
	defer closer.Stop()

	limited := &io.LimitedReader{R: pipe, N: bounds.maxOutputBytes + 1}
	output, readErr := io.ReadAll(limited)
	if int64(len(output)) > bounds.maxOutputBytes {
		// Group-wide, like Cancel: a grandchild holding the pipe would otherwise
		// outlive the call, because Wait returns once the direct child is reaped
		// and the deadline's group kill never fires (queue row 24).
		errs := []error{&HandlerOutputOverflowError{Limit: bounds.maxOutputBytes}}
		killErr := killGroup(cmd.Process.Pid)
		waitErr := wait()
		// The group leader may already be an unreaped zombie when the overflow
		// kill fires: kill(-pgid) then answers ESRCH (linux: an empty group)
		// or EPERM (darwin: zombies only). Neither is a failure when the group
		// is empty after the reap; every other kill error is still joined.
		if killErr != nil && !procbound.NothingToKill(killErr, func() bool { return procbound.GroupEmpty(cmd.Process.Pid) }) {
			errs = append(errs, fmt.Errorf("broker: overflow kill: %w", killErr))
		}
		if errors.Is(waitErr, procbound.ErrCleanupIncomplete) {
			errs = append(errs, waitErr)
		}
		if len(errs) == 1 {
			return nil, errs[0]
		}
		return nil, errors.Join(errs...)
	}
	waitErr := wait()
	if waitErr != nil && runCtx.Err() == context.DeadlineExceeded {
		if errors.Is(waitErr, procbound.ErrCleanupIncomplete) {
			return nil, errors.Join(&HandlerTimeoutError{Timeout: bounds.execTimeout}, waitErr)
		}
		return nil, &HandlerTimeoutError{Timeout: bounds.execTimeout}
	}
	if readErr != nil {
		return nil, fmt.Errorf("broker: read handler output: %w", readErr)
	}
	if spec.stderr != nil && spec.stderr.overflowed() {
		return nil, &HandlerOutputOverflowError{Limit: spec.stderr.limit}
	}
	if waitErr != nil {
		return nil, &HandlerExitError{Err: waitErr, Output: output}
	}
	return output, nil
}
