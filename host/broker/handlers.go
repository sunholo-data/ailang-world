package broker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
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
	// capture is the capture strategy. Its zero value (nil) is the strict
	// capture every handler but Workspace.Exec uses, byte-for-byte: stdout
	// over the bound is an error that kills the group, stderr is merged or
	// goes to the boundedSink above. A non-nil captureHeadTail keeps each
	// stream's head and tail in its own sinks instead (row 140 §4.5); stdin
	// and stderr above are then unused (the child's stdin is /dev/null).
	capture *captureHeadTail
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
	if spec.capture != nil {
		return runHeadTail(runCtx, cmd, bounds, spec.capture)
	}
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
		return nil, overflowKill(cmd, wait, bounds.maxOutputBytes)
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

// overflowKill is the overflow path both capture strategies share.
// Group-wide, like Cancel: a grandchild holding the pipe would otherwise
// outlive the call, because Wait returns once the direct child is reaped and
// the deadline's group kill never fires (queue row 24).
func overflowKill(cmd *exec.Cmd, wait func() error, limit int64) error {
	errs := []error{&HandlerOutputOverflowError{Limit: limit}}
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
		return errs[0]
	}
	return errors.Join(errs...)
}

// Row 140 §4.5: the Workspace.Exec output bounds, per stream.
const (
	execHeadBytes = 8 << 10
	execTailBytes = 56 << 10
)

// captureHeadTail is runBounded's second capture strategy (row 140 §4.5):
// per stream it keeps the first head and the last tail bytes, counts every
// byte and hashes the whole stream, and never blocks the child. A stream past
// hardCap kills the group through overflowKill, exactly as the strict
// capture's overflow does, and runBounded returns
// *HandlerOutputOverflowError{Limit: hardCap}. Whatever the outcome — exit,
// deadline or overflow — the sinks keep what was read.
type captureHeadTail struct {
	hardCap        int64
	stdout, stderr *headTailSink
	over           chan struct{}
	overOnce       sync.Once
}

func newCaptureHeadTail(head, tail int, hardCap int64) *captureHeadTail {
	c := &captureHeadTail{hardCap: hardCap, over: make(chan struct{})}
	signal := func() { c.overOnce.Do(func() { close(c.over) }) }
	c.stdout = newHeadTailSink(head, tail, hardCap, signal)
	c.stderr = newHeadTailSink(head, tail, hardCap, signal)
	return c
}

// runHeadTail is runBounded after cmd is built, under captureHeadTail. Both
// streams go to their sinks through exec's own copying goroutines, which Wait
// joins (WaitDelay force-closes a pipe a descendant still holds). The
// lifecycle — Setpgid, the group kill on the deadline, procbound, the
// overflow kill — is runBounded's own.
func runHeadTail(runCtx context.Context, cmd *exec.Cmd, bounds handlerBounds, c *captureHeadTail) ([]byte, error) {
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	cmd.WaitDelay = pipeCloseGrace
	release, err := procbound.Admit()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		release()
		if errors.Is(err, context.DeadlineExceeded) && errors.Is(runCtx.Err(), context.DeadlineExceeded) {
			return nil, &HandlerTimeoutError{Timeout: bounds.execTimeout}
		}
		return nil, fmt.Errorf("broker: start handler subprocess: %w", err)
	}
	cleanupDeadline := time.Now().Add(bounds.execTimeout + 2*pipeCloseGrace)
	done := make(chan error, 1)
	go func() { done <- procbound.Wait(cmd.Wait, time.Until(cleanupDeadline), release) }()
	var waitErr error
	select {
	case <-c.over:
		return nil, overflowKill(cmd, func() error { return <-done }, c.hardCap)
	case waitErr = <-done:
	}
	if waitErr != nil && runCtx.Err() == context.DeadlineExceeded {
		if errors.Is(waitErr, procbound.ErrCleanupIncomplete) {
			return nil, errors.Join(&HandlerTimeoutError{Timeout: bounds.execTimeout}, waitErr)
		}
		return nil, &HandlerTimeoutError{Timeout: bounds.execTimeout}
	}
	// The cap may be crossed by the last writes before the exit: still an
	// overflow, never a silently capped success.
	select {
	case <-c.over:
		return nil, &HandlerOutputOverflowError{Limit: c.hardCap}
	default:
	}
	if waitErr != nil {
		return nil, &HandlerExitError{Err: waitErr}
	}
	return nil, nil
}

// headTailSink is one stream of captureHeadTail. Write never fails and blocks
// only on its own mutex, so the child never stalls on a full pipe; past
// hardCap it keeps counting and hashing and signals onOver.
type headTailSink struct {
	mu      sync.Mutex
	headMax int
	tailMax int
	head    []byte
	tail    []byte // the last ≤ 2×tailMax bytes; read as the last tailMax
	total   int64
	hash    hash.Hash
	hardCap int64
	onOver  func()
}

func newHeadTailSink(head, tail int, hardCap int64, onOver func()) *headTailSink {
	return &headTailSink{headMax: head, tailMax: tail, hardCap: hardCap, hash: sha256.New(), onOver: onOver}
}

func (s *headTailSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hash.Write(p)
	s.total += int64(len(p))
	if room := s.headMax - len(s.head); room > 0 {
		s.head = append(s.head, p[:min(room, len(p))]...)
	}
	keep := p
	if len(keep) > s.tailMax {
		keep = keep[len(keep)-s.tailMax:]
	}
	s.tail = append(s.tail, keep...)
	if len(s.tail) > 2*s.tailMax {
		s.tail = append(s.tail[:0], s.tail[len(s.tail)-s.tailMax:]...)
	}
	if s.total > s.hardCap && s.onOver != nil {
		s.onOver()
	}
	return len(p), nil
}

// streamCapture is one stream as the exec result records it.
type streamCapture struct {
	Text      string
	Bytes     int64
	Truncated bool
	SHA256    string
}

// truncationMarker joins a truncated stream's head and tail.
func truncationMarker(omitted int64) string {
	return fmt.Sprintf("\n[... %d bytes omitted ...]\n", omitted)
}

func (s *headTailSink) snapshot() streamCapture {
	s.mu.Lock()
	defer s.mu.Unlock()
	tail := s.tail
	if len(tail) > s.tailMax {
		tail = tail[len(tail)-s.tailMax:]
	}
	c := streamCapture{Bytes: s.total, SHA256: hex.EncodeToString(s.hash.Sum(nil))}
	if s.total <= int64(s.headMax+s.tailMax) {
		// Whole: the head plus the bytes after it, which the tail holds.
		rest := int(s.total) - len(s.head)
		c.Text = string(s.head) + string(tail[len(tail)-rest:])
		return c
	}
	c.Truncated = true
	c.Text = string(s.head) + truncationMarker(s.total-int64(len(s.head))-int64(len(tail))) + string(tail)
	return c
}

// handlerExitStatus is a *HandlerExitError's exit status: the code of a
// normal exit, 128+n for a death by signal n (the shell's convention).
func handlerExitStatus(err error) (int, bool) {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return 0, false
	}
	if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		return 128 + int(ws.Signal()), true
	}
	return exitErr.ExitCode(), true
}

// cleanupIncomplete reports a run whose children were not all reaped.
func cleanupIncomplete(err error) bool { return errors.Is(err, procbound.ErrCleanupIncomplete) }
