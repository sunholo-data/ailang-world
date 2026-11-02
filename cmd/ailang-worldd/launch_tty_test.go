package main

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/ptytest"
)

func fixtureOpener(t *testing.T) (func() (*os.File, error), **os.File) {
	t.Helper()
	var opened *os.File
	return func() (*os.File, error) {
		f, err := os.OpenFile(filepath.Join(t.TempDir(), "tty"), os.O_RDWR|os.O_CREATE, 0o600)
		opened = f
		return f, err
	}, &opened
}

func TestLaunchTerminalPrefersStdinTerminal(t *testing.T) {
	master, slave, err := ptytest.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer slave.Close()
	open, opened := fixtureOpener(t)
	got, err := launchTerminal(slave, open)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Close()
	if got.Fd() == slave.Fd() {
		t.Fatal("the agent must get a duplicate, never our own stdin (closing it would close ours)")
	}
	var a, b syscall.Stat_t
	if syscall.Fstat(int(got.Fd()), &a) != nil || syscall.Fstat(int(slave.Fd()), &b) != nil || a.Rdev != b.Rdev {
		t.Fatalf("got device %d, want the stdin pty %d", a.Rdev, b.Rdev)
	}
	if _, err := (*opened).Stat(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("the /dev/tty handle must be closed once stdin is chosen: %v", err)
	}
}

func TestLaunchTerminalFallsBackWithoutStdinTerminal(t *testing.T) {
	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer null.Close()
	open, opened := fixtureOpener(t)
	got, err := launchTerminal(null, open)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Close()
	if got != *opened {
		t.Fatal("with a non-terminal stdin (/dev/null is a character device, not a tty) the opened terminal is used")
	}
}

func TestLaunchTerminalRefusesWithoutControllingTerminal(t *testing.T) {
	master, slave, err := ptytest.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer slave.Close()
	want := errors.New("open /dev/tty: device not configured")
	if _, err := launchTerminal(slave, func() (*os.File, error) { return nil, want }); !errors.Is(err, want) {
		t.Fatalf("no controlling terminal must still refuse, even with a terminal stdin: %v", err)
	}
}
