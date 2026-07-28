package broker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
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
}

// runBounded is the one timeout and allocation surface used by every subprocess
// handler. It reads at most limit+1 bytes so overflow is detected, never hidden
// by truncation.
func runBounded(ctx context.Context, bounds handlerBounds, spec handlerCommand) ([]byte, error) {
	bounds = bounds.normalized()
	runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, spec.path, spec.args...)
	cmd.Dir = spec.dir
	cmd.Env = spec.env
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("broker: handler stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("broker: start handler subprocess: %w", err)
	}

	limited := &io.LimitedReader{R: pipe, N: bounds.maxOutputBytes + 1}
	output, readErr := io.ReadAll(limited)
	if int64(len(output)) > bounds.maxOutputBytes {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, &HandlerOutputOverflowError{Limit: bounds.maxOutputBytes}
	}
	waitErr := cmd.Wait()
	if waitErr != nil && runCtx.Err() == context.DeadlineExceeded {
		return nil, &HandlerTimeoutError{Timeout: bounds.execTimeout}
	}
	if readErr != nil {
		return nil, fmt.Errorf("broker: read handler output: %w", readErr)
	}
	if waitErr != nil {
		return nil, &HandlerExitError{Err: waitErr, Output: output}
	}
	return output, nil
}
