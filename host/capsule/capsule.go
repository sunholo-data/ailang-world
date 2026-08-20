// Package capsule runs untrusted AILANG transition source beneath the M3
// process isolation floor.
//
// The floor is intentionally narrow. It does not contain a malicious native
// interpreter, provide OS-enforced network isolation, impose memory or CPU
// limits, or use chroot, containers, or microVMs. Its trust anchor is the
// content-hash-pinned released interpreter; the untrusted input is transition
// source.
package capsule

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	capsuleExecTimeout    = 60 * time.Second
	maxCapsuleOutputBytes = int64(8 << 20)
	entryModulePath       = "host/capsule/main.ail"
	entryFn               = "main"
)

// Config supplies the injectable bounds used by the shipped execution path.
// Zero values select the production defaults.
type Config struct {
	ExecTimeout    time.Duration
	MaxOutputBytes int64
}

// Entry pins both the interpreter and the canonical transition source.
type Entry struct {
	Interpreter hashref.HashRef
	Source      []byte
}

// Result contains the exact stdout and stderr streams produced by a capsule.
type Result struct {
	Stdout []byte
	Stderr []byte
}

// childProcess is the narrow slice of *os.Process the output lifecycle needs.
// Production supplies the real child; tests supply a fake whose Kill and Wait
// write their own observables.
type childProcess interface {
	Kill() error
	Wait() error
}

type cmdChild struct {
	cmd *exec.Cmd
}

func (c cmdChild) Kill() error { return c.cmd.Process.Kill() }
func (c cmdChild) Wait() error { return c.cmd.Wait() }

// HashMismatchError means the resolved interpreter bytes no longer match the
// entry's content address. Execution is refused before the child is started.
type HashMismatchError struct {
	Ref  hashref.HashRef
	Path string
	Got  hashref.HashRef
}

func (e *HashMismatchError) Error() string {
	return fmt.Sprintf("capsule: interpreter hash mismatch at %q: got %q, want %q",
		e.Path, e.Got.String(), e.Ref.String())
}

// TimeoutError means the capsule exceeded its wall-clock allowance.
type TimeoutError struct {
	Limit time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("capsule: execution exceeded wall-clock limit %s", e.Limit)
}

// errOutputLimit is the internal sentinel a capped read returns, so an overflow
// stays distinguishable from an unrelated pipe failure all the way out.
var errOutputLimit = errors.New("capsule: output limit exceeded")

// OutputLimitError means stdout or stderr exceeded the byte cap.
type OutputLimitError struct {
	Limit int64
}

func (e *OutputLimitError) Error() string {
	return fmt.Sprintf("capsule: output exceeded %d-byte limit", e.Limit)
}

func (e *OutputLimitError) Unwrap() error { return errOutputLimit }

// ExecError reports an interpreter failure without discarding its bounded
// output streams.
type ExecError struct {
	Path   string
	Stderr []byte
	Err    error
}

func (e *ExecError) Error() string {
	return fmt.Sprintf("capsule: interpreter run failed at %q (stderr: %q): %v",
		e.Path, string(e.Stderr), e.Err)
}

func (e *ExecError) Unwrap() error { return e.Err }

// Runner resolves interpreters from the artifact archive and applies F1-F6.
type Runner struct {
	archive        *archive.Archive
	execTimeout    time.Duration
	maxOutputBytes int64
	caps           string
}

// New constructs a Runner. The archive is authoritative: callers cannot
// provide an ambient executable path at execution time.
func New(a *archive.Archive, cfg Config) *Runner {
	timeout := cfg.ExecTimeout
	if timeout == 0 {
		timeout = capsuleExecTimeout
	}
	outputLimit := cfg.MaxOutputBytes
	if outputLimit == 0 {
		outputLimit = maxCapsuleOutputBytes
	}
	return &Runner{archive: a, execTimeout: timeout, maxOutputBytes: outputLimit, caps: ""}
}

// Run stages and executes one pinned transition under the six-part floor.
func (r *Runner) Run(entry Entry) (Result, error) {
	execPath, err := r.archive.Resolve(entry.Interpreter)
	if err != nil {
		return Result{}, err
	}
	if err := verifyExecutable(execPath, entry.Interpreter); err != nil {
		return Result{}, err
	}

	root, err := os.MkdirTemp("", "world-capsule-*")
	if err != nil {
		return Result{}, fmt.Errorf("capsule: create root: %w", err)
	}
	defer func() { _ = os.RemoveAll(root) }()

	srcPath := filepath.Join(root, filepath.FromSlash(entryModulePath))
	if err := os.MkdirAll(filepath.Dir(srcPath), 0o755); err != nil {
		return Result{}, fmt.Errorf("capsule: create source directory: %w", err)
	}
	if err := os.WriteFile(srcPath, entry.Source, 0o644); err != nil {
		return Result{}, fmt.Errorf("capsule: stage source: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.execTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, execPath,
		"run", "--quiet", "--caps", r.caps, "--entry", entryFn, entryModulePath)
	cmd.Dir = root
	cmd.Env = []string{"AILANG_FS_SANDBOX=" + root}
	// Same correction as host/broker's runBounded: kill the whole process group,
	// or a forked grandchild keeps the inherited pipes open and outlives F5/F6.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return Result{}, fmt.Errorf("capsule: stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return Result{}, fmt.Errorf("capsule: stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return Result{}, &ExecError{Path: execPath, Err: err}
	}
	res, runErr, err := collectOutput(ctx, stdoutPipe, stderrPipe, r.maxOutputBytes, r.execTimeout, cmdChild{cmd})
	if err != nil {
		return res, err
	}
	if runErr != nil {
		return res, &ExecError{Path: execPath, Stderr: res.Stderr, Err: runErr}
	}
	return res, nil
}

// collectOutput coordinates the post-Start output lifecycle. The caller owns
// the finite cleanup bound and must arrange that cancellation makes both
// supplied readers and Wait return. collectOutput neither makes an arbitrary
// io.Reader cancellable nor makes Wait() error bounded. Production supplies
// that property via context.WithTimeout, exec.CommandContext, and cmd.Cancel's
// group-wide SIGKILL. The residual lifecycle work is owned by queue row 24,
// w-host-subprocess-cleanup-boundary.
func collectOutput(ctx context.Context, stdoutPipe, stderrPipe io.Reader, limit int64, execTimeout time.Duration, child childProcess) (Result, error, error) {
	var stdout, stderr []byte
	var stdoutErr, stderrErr error
	var killOnce sync.Once
	var wg sync.WaitGroup
	// F6 must not decay into F5. Past the cap nothing drains the pipe, so an
	// unkilled child blocks in write() until the wall clock expires and the
	// caller is told "timeout" for what is really an overflow. Kill on overflow,
	// exactly as host/broker's runBounded already does for handler subprocesses.
	drain := func(pipe io.Reader, dst *[]byte, dstErr *error) {
		defer wg.Done()
		*dst, *dstErr = readCapped(pipe, limit)
		if errors.Is(*dstErr, errOutputLimit) {
			killOnce.Do(func() { _ = child.Kill() })
		}
	}
	wg.Add(2)
	go drain(stdoutPipe, &stdout, &stdoutErr)
	go drain(stderrPipe, &stderr, &stderrErr)
	// Reads must complete before Wait, which closes the pipes (os/exec).
	wg.Wait()
	runErr := child.Wait()

	// Overflow outranks the deadline: killing the child is what let Wait return,
	// and that must not be reported as a wall-clock expiry.
	if errors.Is(stdoutErr, errOutputLimit) || errors.Is(stderrErr, errOutputLimit) {
		return Result{Stdout: stdout, Stderr: stderr}, runErr, &OutputLimitError{Limit: limit}
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return Result{Stdout: stdout, Stderr: stderr}, runErr, &TimeoutError{Limit: execTimeout}
	}
	if stdoutErr != nil || stderrErr != nil {
		return Result{Stdout: stdout, Stderr: stderr}, runErr,
			fmt.Errorf("capsule: read output: %w", errors.Join(stdoutErr, stderrErr))
	}
	return Result{Stdout: stdout, Stderr: stderr}, runErr, nil
}

func verifyExecutable(path string, ref hashref.HashRef) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	got, err := hashref.Sum(ref.Algo(), data)
	if err != nil {
		return err
	}
	if got.Digest() != ref.Digest() {
		return &HashMismatchError{Ref: ref, Path: path, Got: got}
	}
	return nil
}

func readCapped(pipe io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(pipe, limit+1))
	// Overflow is decided before err: ReadAll can report both, and an over-cap
	// read is an overflow whatever else went wrong on the way.
	if int64(len(data)) > limit {
		return data[:limit], errOutputLimit
	}
	if err != nil {
		return data, err
	}
	return bytes.Clone(data), nil
}
