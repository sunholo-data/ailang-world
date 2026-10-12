package evidence

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
	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ErrInvalidProducerConfig refuses a producer whose configuration cannot make
// a bounded, pinned proof. Each refusal names only its own field.
var ErrInvalidProducerConfig = errors.New("evidence: invalid producer configuration")

// ErrExecTimeout reports that the checker outlived ExecTimeout. Its process
// group was killed before the producer returned.
var ErrExecTimeout = errors.New("evidence: checker exceeded its execution timeout")

// ErrOutputLimit reports that a checker stream exceeded MaxOutputBytes. The
// group was killed, and no bytes are returned.
var ErrOutputLimit = errors.New("evidence: checker output exceeded its limit")

// ObjectWriter is the write seam of the producer. BusyTimeout reports the
// store's lock-retry window so the write timeout can be ordered above it.
type ObjectWriter interface {
	PutObject(context.Context, store.Object) error
	BusyTimeout() time.Duration
}

// CheckerConfig pins the external checker by its bytes and its version line.
type CheckerConfig struct {
	ExecPath        string
	Compiler        hashref.HashRef
	CompilerVersion string
}

// ProducerConfig bounds every leg of a proof: the two object calls and the
// checker's wall clock and captured output.
type ProducerConfig struct {
	ObjectReadTimeout  time.Duration
	ObjectWriteTimeout time.Duration
	ExecTimeout        time.Duration
	MaxOutputBytes     int64
}

// Producer mints authenticated proof envelopes from a pinned checker run.
type Producer struct {
	key     *[32]byte
	reader  ObjectReader
	writer  ObjectWriter
	checker CheckerConfig
	cfg     ProducerConfig
}

// NewProducer refuses a configuration it cannot bound. Each refusal is a
// separate branch so each is separately observable and separately killable.
func NewProducer(key [32]byte, reader ObjectReader, writer ObjectWriter, checker CheckerConfig, cfg ProducerConfig) (*Producer, error) {
	if reader == nil {
		return nil, fmt.Errorf("%w: reader is nil", ErrInvalidProducerConfig)
	}
	if writer == nil {
		return nil, fmt.Errorf("%w: writer is nil", ErrInvalidProducerConfig)
	}
	if checker.ExecPath == "" {
		return nil, fmt.Errorf("%w: ExecPath is empty", ErrInvalidProducerConfig)
	}
	if checker.Compiler.IsZero() {
		return nil, fmt.Errorf("%w: Compiler identity is unset", ErrInvalidProducerConfig)
	}
	if checker.CompilerVersion == "" {
		return nil, fmt.Errorf("%w: CompilerVersion is empty", ErrInvalidProducerConfig)
	}
	if cfg.ObjectReadTimeout <= 0 {
		return nil, fmt.Errorf("%w: ObjectReadTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ObjectReadTimeout)
	}
	if cfg.ObjectWriteTimeout <= 0 {
		return nil, fmt.Errorf("%w: ObjectWriteTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ObjectWriteTimeout)
	}
	if cfg.ExecTimeout <= 0 {
		return nil, fmt.Errorf("%w: ExecTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ExecTimeout)
	}
	if cfg.MaxOutputBytes <= 0 {
		return nil, fmt.Errorf("%w: MaxOutputBytes must be positive, got %d", ErrInvalidProducerConfig, cfg.MaxOutputBytes)
	}
	seams := []struct {
		name string
		busy time.Duration
	}{{"reader", reader.BusyTimeout()}, {"writer", writer.BusyTimeout()}}
	for _, s := range seams {
		if s.busy < 0 {
			return nil, fmt.Errorf("%w: %s BusyTimeout is unknown (%s)", ErrInvalidProducerConfig, s.name, s.busy)
		}
		if s.busy > 0 && cfg.ObjectReadTimeout <= s.busy {
			return nil, fmt.Errorf("%w: ObjectReadTimeout %s must exceed %s BusyTimeout %s", ErrUnorderedTimeouts, cfg.ObjectReadTimeout, s.name, s.busy)
		}
		if s.busy > 0 && cfg.ObjectWriteTimeout <= s.busy {
			return nil, fmt.Errorf("%w: ObjectWriteTimeout %s must exceed %s BusyTimeout %s", ErrUnorderedTimeouts, cfg.ObjectWriteTimeout, s.name, s.busy)
		}
	}
	keyCopy := new([32]byte)
	*keyCopy = key
	return &Producer{key: keyCopy, reader: reader, writer: writer, checker: checker, cfg: cfg}, nil
}

// pipeCloseGrace bounds the one case the group kill cannot reach: a descendant
// that left the process group and still holds the pipes.
const pipeCloseGrace = time.Second

// killGroup SIGKILLs the checker's whole process group. It is a package-level
// seam so tests can observe or fail the kill.
var killGroup = func(pgid int) error {
	return syscall.Kill(-pgid, syscall.SIGKILL)
}

var errStreamOverflow = errors.New("stream exceeds its output limit")

// readCapped reads at most limit bytes and reports an overflow when the stream
// holds more. The overflow decision precedes err: ReadAll can report both, and
// an over-cap stream is an overflow whatever else went wrong.
func readCapped(pipe io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(pipe, limit+1))
	if int64(len(data)) > limit {
		return data[:limit], errStreamOverflow
	}
	if err != nil {
		return data, err
	}
	return bytes.Clone(data), nil
}

// runBounded runs the pinned checker in its own process group under one
// wall-clock bound and a per-stream output cap. Exit codes are advisory: only
// the captured stdout and the bounded classification decide the outcome.
func (p *Producer) runBounded(ctx context.Context, dir string, env []string, limit int64, args ...string) ([]byte, error) {
	release, err := procbound.Admit()
	if err != nil {
		return nil, err
	}
	execCtx, cancel := context.WithTimeout(ctx, p.cfg.ExecTimeout)
	defer cancel()
	cmd := exec.CommandContext(execCtx, p.checker.ExecPath, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return killGroup(cmd.Process.Pid)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		release()
		return nil, fmt.Errorf("evidence: stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		release()
		return nil, fmt.Errorf("evidence: stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		release()
		if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: start: %v", ErrExecTimeout, err)
		}
		return nil, fmt.Errorf("evidence: start %s: %w", p.checker.ExecPath, err)
	}
	started := time.Now()
	closer := time.AfterFunc(p.cfg.ExecTimeout+pipeCloseGrace, func() {
		_ = stdoutPipe.Close()
		_ = stderrPipe.Close()
	})
	defer closer.Stop()

	var stdout, stderr []byte
	var stdoutErr, stderrErr error
	var killOnce sync.Once
	var killErr error
	var wg sync.WaitGroup
	drain := func(pipe io.Reader, dst *[]byte, dstErr *error) {
		defer wg.Done()
		*dst, *dstErr = readCapped(pipe, limit)
		if errors.Is(*dstErr, errStreamOverflow) {
			killOnce.Do(func() { killErr = killGroup(cmd.Process.Pid) })
		}
	}
	wg.Add(2)
	go drain(stdoutPipe, &stdout, &stdoutErr)
	go drain(stderrPipe, &stderr, &stderrErr)
	// Reads finish before Wait: the pipes close when the reader returns (os/exec).
	wg.Wait()
	runErr := procbound.Wait(cmd.Wait, time.Until(started.Add(p.cfg.ExecTimeout+2*pipeCloseGrace)), release)

	switch {
	case errors.Is(stdoutErr, errStreamOverflow) || errors.Is(stderrErr, errStreamOverflow):
		// The group leader may already be an unreaped zombie when the kill fires:
		// NothingToKill separates that benign answer from a real kill failure.
		stream := "stderr"
		if errors.Is(stdoutErr, errStreamOverflow) {
			stream = "stdout"
		}
		var killFailure error
		if killErr != nil && !procbound.NothingToKill(killErr, func() bool { return procbound.GroupEmpty(cmd.Process.Pid) }) {
			killFailure = fmt.Errorf("evidence: overflow kill: %w", killErr)
		}
		return nil, withRunFailures(fmt.Errorf("%w: %s exceeded %d bytes", ErrOutputLimit, stream, limit), killFailure, runErr)
	case errors.Is(execCtx.Err(), context.DeadlineExceeded):
		return nil, withRunFailures(fmt.Errorf("%w after %s", ErrExecTimeout, p.cfg.ExecTimeout), nil, runErr)
	case stdoutErr != nil || stderrErr != nil:
		return nil, fmt.Errorf("evidence: read checker output: %w", errors.Join(stdoutErr, stderrErr))
	}
	var exitErr *exec.ExitError
	if runErr != nil && !errors.As(runErr, &exitErr) {
		return nil, fmt.Errorf("evidence: wait for checker: %w", runErr)
	}
	return stdout, nil
}

// withRunFailures keeps the typed error primary and joins a failed kill and an
// incomplete cleanup behind it. With neither, it returns primary itself.
func withRunFailures(primary, killFailure, waitErr error) error {
	errs := []error{primary}
	if killFailure != nil {
		errs = append(errs, killFailure)
	}
	if errors.Is(waitErr, procbound.ErrCleanupIncomplete) {
		errs = append(errs, waitErr)
	}
	if len(errs) == 1 {
		return primary
	}
	return errors.Join(errs...)
}

const (
	// stagedModule is the only module the checker sees: the source under test,
	// copied into a scratch root so the checker never reads the caller's tree.
	stagedModule = "staged/module.ail"
	// aiCheckZ3Timeout is the checker's own solver budget, inside ExecTimeout.
	aiCheckZ3Timeout = "5s"
)

// runCheck stages source in a fresh scratch root and returns the checker's
// captured stdout. The checker's temp and cache writes are confined to that
// root, which is removed on return. Relax and cache are set here, after the
// scrubbed ambient environment, so an inherited value cannot win.
func (p *Producer) runCheck(ctx context.Context, source []byte) ([]byte, error) {
	root, err := os.MkdirTemp("", "evidence-producer-*")
	if err != nil {
		return nil, fmt.Errorf("evidence: scratch root: %w", err)
	}
	defer os.RemoveAll(root)
	for _, dir := range []string{"staged", "tmp", "cache"} {
		if err := os.Mkdir(filepath.Join(root, dir), 0o755); err != nil {
			return nil, fmt.Errorf("evidence: scratch %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, stagedModule), source, 0o644); err != nil {
		return nil, fmt.Errorf("evidence: stage source: %w", err)
	}
	env := append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1",
		"TMPDIR="+filepath.Join(root, "tmp"), archive.CacheDirEnv+"="+filepath.Join(root, "cache"))
	return p.runBounded(ctx, root, env, p.cfg.MaxOutputBytes, "ai-check", "-timeout", aiCheckZ3Timeout, stagedModule)
}
