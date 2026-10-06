//go:build darwin || linux

// Package ptytest drives the attended (TTY-fenced) commands end to end from a
// test: it opens a pseudo-terminal pair and starts a child process either with
// that pty as its CONTROLLING terminal or in a new session with NO controlling
// terminal at all.
//
// The two shapes are the two sides of the fence, measured rather than
// simulated:
//
//   - WithoutControllingTerminal is setsid(2) with no TIOCSCTTY. In the child,
//     opening /dev/tty fails with ENXIO, which is exactly what an agent
//     harness, cron or a CI runner sees.
//   - WithControllingTerminal is setsid(2) + TIOCSCTTY on the pty's slave, so
//     the child's /dev/tty IS the pty and whatever the test writes to the
//     master is what a person typed.
//
// Stdlib only (syscall ioctls), no new module dependency. Test support: no
// production code imports it.
package ptytest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Open returns a new pseudo-terminal pair. The caller closes both.
func Open() (master, slave *os.File, err error) {
	return open()
}

// WithControllingTerminal makes cmd start in a new session whose controlling
// terminal is slave. slave is passed as an extra file and named by its CHILD
// descriptor number, which is what SysProcAttr.Ctty means.
func WithControllingTerminal(cmd *exec.Cmd, slave *os.File) {
	cmd.ExtraFiles = append(cmd.ExtraFiles, slave)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 2 + len(cmd.ExtraFiles)}
}

// WithoutControllingTerminal makes cmd start in a new session with no
// controlling terminal: its open("/dev/tty") fails.
func WithoutControllingTerminal(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

// Buffer is a child's stdout/stderr sink that is safe to read while os/exec's
// copier goroutine may still be writing (the subprocess-sink gate's rule).
type Buffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *Buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// Transcript collects everything the child writes to the pty, in the
// background, so a test can wait for a prompt without a pollable descriptor.
type Transcript struct {
	mu   sync.Mutex
	buf  bytes.Buffer
	done chan struct{}
}

// Record starts copying master into a Transcript. It returns when the master
// reports EOF or EIO (every slave descriptor closed).
func Record(master *os.File) *Transcript {
	tr := &Transcript{done: make(chan struct{})}
	go func() {
		defer close(tr.done)
		chunk := make([]byte, 4096)
		for {
			n, err := master.Read(chunk)
			if n > 0 {
				tr.mu.Lock()
				tr.buf.Write(chunk[:n])
				tr.mu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()
	return tr
}

// String is everything recorded so far.
func (tr *Transcript) String() string {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return tr.buf.String()
}

// WaitFor polls until want has been written to the pty, or timeout passes.
func (tr *Transcript) WaitFor(want string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.Contains(tr.String(), want) {
			return nil
		}
		select {
		case <-tr.done:
			if strings.Contains(tr.String(), want) {
				return nil
			}
			return fmt.Errorf("ptytest: the pty closed before %q appeared; transcript %q", want, tr.String())
		case <-time.After(20 * time.Millisecond):
		}
	}
	return fmt.Errorf("ptytest: %q did not appear on the pty within %v; transcript %q", want, timeout, tr.String())
}
