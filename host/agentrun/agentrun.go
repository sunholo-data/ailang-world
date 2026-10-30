// Package agentrun owns the bounded agent subprocess site.
// Output goes to private staging files or the caller\'s controlling terminal.
package agentrun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/procbound"
)

type Result struct {
	ExitCode int
	Reason   string
}

// Run takes an already scrubbed environment, allowing Codex's child-only token.
// Stdin is /dev/null. Cancellation and the wall cap kill the entire group;
// received terminal signals are forwarded, then bounded by the same wall cap.
func Run(ctx context.Context, argv []string, dir string, env []string, stdout, stderr *os.File, timeout time.Duration) (Result, error) {
	return run(ctx, argv, dir, env, stdout, stderr, nil, timeout)
}

// RunInteractive runs in the terminal foreground group. Terminal signals already
// reach the child; only SIGTERM is forwarded. Cancellation kills only its pid.
func RunInteractive(ctx context.Context, argv []string, dir string, env []string, tty *os.File, timeout time.Duration) (Result, error) {
	if tty == nil {
		return Result{ExitCode: -1, Reason: "exit"}, errors.New("agentrun: terminal required")
	}
	return run(ctx, argv, dir, env, tty, tty, tty, timeout)
}

func run(ctx context.Context, argv []string, dir string, env []string, stdout, stderr, tty *os.File, timeout time.Duration) (Result, error) {
	r := Result{ExitCode: -1, Reason: "exit"}
	if len(argv) == 0 || timeout <= 0 || stdout == nil || stderr == nil {
		return r, fmt.Errorf("agentrun: argv, output files and positive wall cap required")
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, argv[0], argv[1:]...)
	cmd.Dir, cmd.Env = dir, env
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if tty == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	} else {
		cmd.Stdin = tty
		cmd.Cancel = func() error { return syscall.Kill(cmd.Process.Pid, syscall.SIGKILL) }
	}
	cmd.WaitDelay = 100 * time.Millisecond
	signals := make(chan os.Signal, 2)
	if tty == nil {
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	} else {
		signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM)
	}
	defer signal.Stop(signals)
	release, err := procbound.Admit()
	if err != nil {
		return r, err
	}
	if err := cmd.Start(); err != nil {
		release()
		return r, err
	}
	done := make(chan struct{})
	forwarded := make(chan bool, 1)
	go func() {
		seen := false
		for {
			select {
			case sig := <-signals:
				if tty == nil {
					seen = true
					_ = syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
				} else if sig == syscall.SIGTERM {
					seen = true
					_ = syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
				}
			case <-done:
				forwarded <- seen
				return
			}
		}
	}()
	err = procbound.Wait(cmd.Wait, timeout+time.Second, release)
	close(done)
	seen := <-forwarded
	if cmd.ProcessState != nil {
		r.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runCtx.Err() != nil && (tty == nil || errors.Is(runCtx.Err(), context.DeadlineExceeded)) {
		r.Reason = "timeout"
	} else if seen || r.ExitCode < 0 {
		r.Reason = "signal"
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		return r, err
	}
	if tty != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) || runCtx.Err() != nil {
			return r, nil
		}
		return r, err
	}
	// A child that backgrounds descendants must not leave them running.
	killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	// darwin: a descendant just killed by the wall cap is reparented and stays a
	// zombie until launchd reaps it, so give the reap a second before calling a
	// zombie-only group (EPERM) a cleanup failure.
	groupGone := func() bool {
		for deadline := time.Now().Add(time.Second); ; time.Sleep(10 * time.Millisecond) {
			if procbound.GroupEmpty(cmd.Process.Pid) {
				return true
			}
			if time.Now().After(deadline) {
				return false
			}
		}
	}
	if killErr != nil && !procbound.NothingToKill(killErr, groupGone) {
		return r, errors.Join(err, killErr)
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) || runCtx.Err() != nil {
		return r, nil
	}
	return r, err
}
