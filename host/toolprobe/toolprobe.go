// Package toolprobe is the one bounded, environment-scrubbed subprocess site
// `ailang-worldd profile init` (row 189 M4) uses to ask a repo's own
// toolchain where it lives: `go env GOROOT GOMODCACHE GOVERSION` and `uv
// python find`.
//
// Why a package: the AC10 census (host/broker registry_publish_test.go)
// requires every production subprocess site to be DRIVEN and to scrub the
// registry credential, and a site in package main cannot be driven from
// there, so the site lives here (as host/worktree's does) with an injectable
// binary path the census can point at its env-dumping probe.
//
// The child environment is the process environment minus every registry
// variable (childenv.Scrubbed).
package toolprobe

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// Timeout bounds one probe; waitDelay bounds Wait after a kill.
const (
	Timeout   = 30 * time.Second
	waitDelay = 2 * time.Second
)

// Output runs bin with args in dir and returns its standard output, as
// exec.Cmd.Output does: a non-zero exit is an *exec.ExitError carrying the
// child's stderr.
func Output(ctx context.Context, bin, dir string, args ...string) ([]byte, error) {
	return OutputEnv(ctx, bin, dir, os.Environ(), args...)
}

// OutputEnv shares the AC10 driven subprocess site and always scrubs credentials.
// Doctor supplies only PATH and HOME; launcher tests supply their explicit fixture env.
func OutputEnv(ctx context.Context, bin, dir string, env []string, args ...string) ([]byte, error) {
	runCtx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, bin, args...)
	cmd.Dir = dir
	cmd.Env = childenv.Scrubbed(env)
	cmd.Stdin = nil
	cmd.WaitDelay = waitDelay
	return cmd.Output()
}
