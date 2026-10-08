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
// interpreter: the bytes are staged in a scratch root at CapsuleEntryPath (the
// capsule's own staging path) and `<archived interpreter> check <file>` runs
// there, bounded by procbound (a process-wide reservation plus a bounded reap)
// and a wall-clock context, under the scrubbed child environment — the
// probeVersion pattern this package already uses for `--version`. `check`
// parses and type-checks; it never executes the module. ok=false means the
// interpreter refused the source (non-zero exit, or the binary failed to
// start); err is reserved for host-side infrastructure failures (scratch
// root, admission, and — row 153 — a passing check whose compile-cache
// template could not be built, a *TemplateBuildError). Callers that need the
// interpreter's verdict keep Output for their own typed errors.
//
// A passing check is also the template build (row 153): the interpreter is
// pointed at a private cache directory (AILANG_CACHE_DIR, appended last so an
// ambient value cannot win), and on success that directory is promoted to the
// (interpreter, source) template every capsule run copies. A pass that wrote
// no compile/manifest.json (an interpreter that ignores the variable, a full
// disk) is a *TemplateBuildError: publishing a transition that would run cold
// forever is refused rather than silently accepted.
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
	checkFile := filepath.FromSlash(CapsuleEntryPath)
	if err := os.MkdirAll(filepath.Join(root, filepath.Dir(checkFile)), 0o755); err != nil {
		return CheckResult{}, fmt.Errorf("check source: create source directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(root, checkFile), source, 0o644); err != nil {
		return CheckResult{}, fmt.Errorf("check source: stage source: %w", err)
	}
	tmp, err := a.NewCapsuleTemplateTmp()
	if err != nil {
		return CheckResult{}, fmt.Errorf("check source: create template tmp: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmp) }() // a no-op once PromoteTemplate has renamed it away
	runCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, execPath, "check", checkFile)
	cmd.Dir = root
	// AILANG_RELAX_MODULES=1: a declared module path never matches the staging
	// path; without the relax flag the verdict depends on whether the
	// interpreter recognises the scratch root as a temp dir (macOS resolves it
	// under /private/var and does not) — measured, iter-195. Staging at the
	// capsule's path (row 153) changes the name, not that reason.
	cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1", CacheDirEnv+"="+tmp)
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
	passed := CheckResult{Output: outputHead(stdout.String() + stderr.String()), Passed: true}
	if info, err := os.Stat(filepath.Join(tmp, filepath.FromSlash(manifestRel))); err != nil || !info.Mode().IsRegular() {
		return passed, &TemplateBuildError{Interpreter: interpreter, Source: hashref.SumSHA256(source),
			Reason: "check passed but wrote no compile/manifest.json"}
	}
	if err := a.PromoteTemplate(tmp, interpreter, source); err != nil {
		return passed, err
	}
	return passed, nil
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
