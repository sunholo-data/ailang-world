package archive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
)

// checkTimeout bounds one `<interpreter> check <source>` subprocess — the
// same wall clock the version probe gets (probeTimeout).
const checkTimeout = 10 * time.Second

// checkBuffer is a concurrency-safe sink for the checked interpreter's
// streams: os/exec's copier goroutine keeps writing after procbound.Wait can
// return (ErrCleanupIncomplete leaves the reap to a background waiter), so
// the sinks must stay safe to read concurrently — the syncBuffer pattern the
// subprocess-sink gate (host/verifygate) demands for Start-based launches.
type checkBuffer struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (s *checkBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *checkBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

// CheckResult reports one bounded source check: the interpreter's combined
// output (whitespace-trimmed, capped) and whether it accepted the source.
type CheckResult struct {
	Output string
	Passed bool
}

// CheckSource proves canonical source bytes load under the archived
// interpreter: the bytes are staged in a scratch root and
// `<archived interpreter> check <file>` runs there, bounded by procbound
// (a process-wide reservation plus a bounded reap) and a wall-clock context,
// under the scrubbed child environment — the probeVersion pattern this
// package already uses for `--version`. `check` parses and type-checks; it
// never executes the module. ok=false means the interpreter refused the
// source (non-zero exit, or the binary failed to start); err is reserved for
// host-side infrastructure failures (scratch root, admission). Callers that
// need the interpreter's verdict keep Output for their own typed errors.
func (a *Archive) CheckSource(ctx context.Context, interpreter hashref.HashRef, source []byte) (CheckResult, error) {
	execPath, err := a.Resolve(interpreter)
	if err != nil {
		return CheckResult{}, fmt.Errorf("check source: resolve pinned interpreter %q: %w", interpreter.String(), err)
	}
	root, err := os.MkdirTemp("", "world-source-check-*")
	if err != nil {
		return CheckResult{}, fmt.Errorf("check source: create scratch root: %w", err)
	}
	defer func() { _ = os.RemoveAll(root) }()
	const checkFile = "entry.ail"
	if err := os.WriteFile(filepath.Join(root, checkFile), source, 0o644); err != nil {
		return CheckResult{}, fmt.Errorf("check source: stage source: %w", err)
	}
	runCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, execPath, "check", checkFile)
	cmd.Dir = root
	// AILANG_RELAX_MODULES=1: the staged file is entry.ail, so a declared
	// module path never matches it; without the relax flag the verdict depends
	// on whether the interpreter recognises the scratch root as a temp dir
	// (macOS resolves it under /private/var and does not) — measured, iter-195.
	cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1")
	var stdout, stderr checkBuffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	release, err := procbound.Admit()
	if err != nil {
		return CheckResult{}, fmt.Errorf("check source: %w", err)
	}
	if err := cmd.Start(); err != nil {
		release()
		return CheckResult{Output: outputHead(stdout.String() + stderr.String()), Passed: false}, nil
	}
	if waitErr := procbound.Wait(cmd.Wait, checkTimeout, release); waitErr != nil {
		output := outputHead(stdout.String() + stderr.String())
		if waitErr == procbound.ErrCleanupIncomplete {
			output = "interpreter check did not finish within " + checkTimeout.String() + "; " + output
		}
		return CheckResult{Output: output, Passed: false}, nil
	}
	return CheckResult{Output: outputHead(stdout.String() + stderr.String()), Passed: true}, nil
}

// outputHead keeps captured interpreter output readable in one error line:
// whitespace-trimmed, capped at 400 bytes.
func outputHead(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 400 {
		return s[:400] + "…"
	}
	return s
}
