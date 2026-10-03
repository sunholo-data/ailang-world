package verifygate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"
)

// Every test-side file write in this package goes through one of the wrappers
// below so that no fork in the test process can inherit a write fd (golang/go
// #22315: a forked child holds a copy of every open fd until its own execve,
// so a later execve of a just-written script fails with ETXTBSY on Linux).
// Every fork takes syscall.ForkLock for WRITING, so a writer holding RLock
// from open to close excludes all forks. See
// design_docs/planned/w-verifygate-etxtbsy.md.

// forkLockedDo holds syscall.ForkLock.RLock across open, fill, hold and Close.
// NOTHING inside the region may fork (exec.Command.Start, os.StartProcess, ...):
// a fork calls ForkLock.Lock and would deadlock against this goroutine's own
// RLock. fill must write bytes only. hold is a test-only seam, nil in every
// production-of-tests caller.
func forkLockedDo(open func() (*os.File, error), fill func(*os.File) error, hold func()) error {
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()
	f, err := open()
	if err != nil {
		return err
	}
	if fill != nil {
		if err := fill(f); err != nil {
			f.Close()
			return err
		}
	}
	if hold != nil {
		hold()
	}
	return f.Close()
}

func forkLockedWrite(path string, flag int, mode os.FileMode, fill func(*os.File) error, hold func()) error {
	return forkLockedDo(func() (*os.File, error) { return os.OpenFile(path, flag, mode) }, fill, hold)
}

// writeFileForkLocked has os.WriteFile semantics (it does not chmod an existing file).
func writeFileForkLocked(path string, data []byte, mode os.FileMode) error {
	return forkLockedWrite(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, nil)
}

// copyFileForkLocked copies src into the fresh path dst (O_EXCL). Error texts
// match the original copyGateFileErr.
func copyFileForkLocked(src *os.File, dst, rel string, mode os.FileMode) error {
	filled := false
	var copyErr error
	err := forkLockedWrite(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode, func(f *os.File) error {
		filled = true
		_, copyErr = io.Copy(f, src)
		return copyErr
	}, nil)
	if err == nil {
		return nil
	}
	if !filled {
		return fmt.Errorf("create copy target %s: %v", rel, err)
	}
	if copyErr != nil {
		return fmt.Errorf("copy %s: %v", rel, copyErr)
	}
	return err
}

// createTempForkLocked replaces CreateTemp+WriteString+Close. On a fill or
// close error the file is removed and "" is returned.
func createTempForkLocked(dir, pattern string, data []byte) (string, error) {
	var name string
	err := forkLockedDo(func() (*os.File, error) {
		f, err := os.CreateTemp(dir, pattern)
		if f != nil {
			name = f.Name()
		}
		return f, err
	}, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, nil)
	if err != nil {
		if name != "" {
			os.Remove(name)
		}
		return "", err
	}
	return name, nil
}

// createForkLocked is the os.Create shape: create/truncate, empty, closed.
func createForkLocked(path string) error {
	return forkLockedWrite(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666, nil, nil)
}

func TestForkLockedWriteBlocksConcurrentFork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "locked.txt")
	inside := make(chan struct{})
	release := make(chan struct{})
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- forkLockedWrite(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644, nil, func() {
			close(inside)
			<-release
		})
	}()
	<-inside
	started := make(chan error, 1)
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	go func() { started <- cmd.Start() }()
	select {
	case <-started:
		close(release)
		t.Fatal("Start returned while a fork-locked write fd was open")
	case <-time.After(250 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-started:
		if err != nil {
			t.Fatalf("Start failed after release: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return within 5 s of release")
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := <-writeDone; err != nil {
		t.Fatal(err)
	}
}

func TestKernelRefusesExecOfWriterOpenFile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("ETXTBSY is a Linux kernel rule; darwin permits exec of a writer-open file (design V7)")
	}
	path := filepath.Join(t.TempDir(), "writer-open.sh")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("#!/bin/sh\nexit 0\n"); err != nil {
		f.Close()
		t.Fatal(err)
	}
	err = exec.Command(path).Start()
	if !errors.Is(err, syscall.ETXTBSY) {
		f.Close()
		t.Fatalf("expected ETXTBSY exec of a writer-open file, got %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command(path).Run(); err != nil {
		t.Fatalf("exec after Close failed: %v", err)
	}
}
