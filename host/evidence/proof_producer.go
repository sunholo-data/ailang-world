package evidence

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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

// ErrUntrustedTool refuses a checker whose bytes or version line are not the
// pinned ones. It is reported before the checker runs on any input.
var ErrUntrustedTool = errors.New("evidence: checker is not the pinned tool")

// ErrUnverifiedRun refuses an ai-check run that did not verify its module clean:
// a failed check, an unavailable verifier, or any error, counterexample or skip.
var ErrUnverifiedRun = errors.New("evidence: ai-check run is not a clean verification")

// ErrVacuousCheck refuses a run that produced output but verified no function.
// A green run that proves nothing is refused, never minted.
var ErrVacuousCheck = errors.New("evidence: ai-check run verified nothing")

// ErrRequiredIdentityMissing refuses a proof that does not verify every
// identity the caller named, and refuses an empty requirement outright.
var ErrRequiredIdentityMissing = errors.New("evidence: required identity not verified")

// ErrEnvelopeTooLarge refuses an envelope the validator could not read: its
// read cap is MaxBytes, so a larger envelope is minted nowhere.
var ErrEnvelopeTooLarge = errors.New("evidence: proof envelope exceeds the validator read cap")

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

// versionProbeMaxBytes caps the --version probe's output on each stream.
const versionProbeMaxBytes = 4 << 10

// verifyTool proves the checker's bytes and then its version line. The bytes
// are hashed before anything executes, so a tampered binary never runs.
func (p *Producer) verifyTool(ctx context.Context) error {
	data, err := os.ReadFile(p.checker.ExecPath)
	if err != nil {
		return fmt.Errorf("%w: read executable %s: %v", ErrUntrustedTool, p.checker.ExecPath, err)
	}
	got, err := hashref.Sum(p.checker.Compiler.Algo(), data)
	if err != nil {
		return fmt.Errorf("%w: hash executable %s: %v", ErrUntrustedTool, p.checker.ExecPath, err)
	}
	if got.Digest() != p.checker.Compiler.Digest() {
		return fmt.Errorf("%w: executable %s bytes hash to %s; want %s", ErrUntrustedTool, p.checker.ExecPath, got, p.checker.Compiler)
	}
	out, err := p.runBounded(ctx, "", childenv.Scrubbed(os.Environ()), versionProbeMaxBytes, "--version")
	if err != nil {
		return fmt.Errorf("%w: version probe: %w", ErrUntrustedTool, err)
	}
	line := strings.TrimRight(strings.SplitN(string(out), "\n", 2)[0], "\r")
	if line != p.checker.CompilerVersion {
		return fmt.Errorf("%w: version line %q; want %q", ErrUntrustedTool, line, p.checker.CompilerVersion)
	}
	return nil
}

// aiCheckReport is the strict subset of ai-check's JSON the producer reads.
// It deliberately has no field for verify.verified: the verified set comes only
// from the results records (D4 step 8).
type aiCheckReport struct {
	Check struct {
		Passed bool `json:"passed"`
	} `json:"check"`
	Verify struct {
		Available      bool `json:"available"`
		Counterexample int  `json:"counterexample"`
		Skipped        int  `json:"skipped"`
		Errors         int  `json:"errors"`
		Results        []struct {
			Function string `json:"function"`
			Status   string `json:"status"`
		} `json:"results"`
	} `json:"verify"`
}

// interpretCheck turns captured ai-check stdout into the sorted set of verified
// functions. Each refusal is one branch with its own detail fragment, applied in
// a fixed order. Unknown JSON fields are tolerated.
func interpretCheck(stdout []byte) ([]string, error) {
	if len(bytes.TrimSpace(stdout)) == 0 {
		return nil, fmt.Errorf("%w: empty ai-check stdout", ErrVacuousCheck)
	}
	var d aiCheckReport
	if err := json.Unmarshal(stdout, &d); err != nil {
		return nil, fmt.Errorf("%w: unparseable ai-check JSON: %v", ErrUnverifiedRun, err)
	}
	if !d.Check.Passed {
		return nil, fmt.Errorf("%w: check.passed is not true", ErrUnverifiedRun)
	}
	if !d.Verify.Available {
		return nil, fmt.Errorf("%w: verify.available is not true", ErrUnverifiedRun)
	}
	if d.Verify.Errors > 0 {
		return nil, fmt.Errorf("%w: verify.errors=%d", ErrUnverifiedRun, d.Verify.Errors)
	}
	if d.Verify.Counterexample > 0 {
		return nil, fmt.Errorf("%w: verify.counterexample=%d", ErrUnverifiedRun, d.Verify.Counterexample)
	}
	if d.Verify.Skipped > 0 {
		return nil, fmt.Errorf("%w: verify.skipped=%d", ErrUnverifiedRun, d.Verify.Skipped)
	}
	var verified []string
	for _, r := range d.Verify.Results {
		if r.Status == "verified" {
			verified = append(verified, r.Function)
		}
	}
	if len(verified) == 0 {
		if len(d.Verify.Results) == 0 {
			return nil, fmt.Errorf("%w: 0 results records from non-empty output", ErrVacuousCheck)
		}
		return nil, fmt.Errorf("%w: 0 verified among %d results records", ErrVacuousCheck, len(d.Verify.Results))
	}
	sort.Strings(verified)
	return verified, nil
}

// requireIdentities admits a verified set only when it contains every named
// identity. An empty requirement is refused here and nowhere else (PD10).
func requireIdentities(verified, required []string) error {
	if len(required) == 0 {
		return fmt.Errorf("%w: no required identities named", ErrRequiredIdentityMissing)
	}
	for _, id := range required {
		if !containsIdentity(verified, id) {
			return fmt.Errorf("%w: required identity %q not verified", ErrRequiredIdentityMissing, id)
		}
	}
	return nil
}

// mint encodes a clean report over the verified set, MACs it under the producer
// key, and refuses any envelope the validator could not read back.
func (p *Producer) mint(subject hashref.HashRef, verified []string) ([]byte, hashref.HashRef, error) {
	report := ProofReportV1{
		Schema:          ProofReportSchemaV1,
		Subject:         subject,
		Compiler:        p.checker.Compiler,
		CompilerVersion: p.checker.CompilerVersion,
		Verified:        verified,
		CheckPassed:     true,
		ProofSucceeded:  true,
	}
	rb, err := EncodeProofReportV1(report)
	if err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("evidence: encode proof report: %w", err)
	}
	mac := hmac.New(sha256.New, p.key[:])
	mac.Write(rb)
	env, err := EncodeAuthenticatedEnvelope(AuthenticatedEnvelope{Report: rb, MAC: mac.Sum(nil), MACValid: true})
	if err != nil {
		return nil, hashref.HashRef{}, fmt.Errorf("evidence: encode envelope: %w", err)
	}
	if len(env) > MaxBytes {
		return nil, hashref.HashRef{}, fmt.Errorf("%w: envelope is %d bytes; validator read cap is %d", ErrEnvelopeTooLarge, len(env), MaxBytes)
	}
	return env, hashref.SumSHA256(env), nil
}
