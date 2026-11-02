//go:build darwin

package main

import (
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/ptytest"
)

// The w-episode-finish M4 exit run (2026-11-02): Codex panicked "reader source
// not set" because kqueue refuses a descriptor opened from /dev/tty (EINVAL).
// The terminal the agent inherits must be one kqueue can poll.
func TestLaunchTerminalIsKqueuePollable(t *testing.T) {
	master, slave, err := ptytest.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer slave.Close()
	open, _ := fixtureOpener(t)
	got, err := launchTerminal(slave, open)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Close()
	kq, err := syscall.Kqueue()
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(kq)
	ev := []syscall.Kevent_t{{Ident: uint64(got.Fd()), Filter: syscall.EVFILT_READ, Flags: syscall.EV_ADD | syscall.EV_RECEIPT}}
	out := make([]syscall.Kevent_t, 1)
	n, err := syscall.Kevent(kq, ev, out, nil)
	if err != nil || n != 1 || (out[0].Flags&syscall.EV_ERROR != 0 && out[0].Data != 0) {
		t.Fatalf("kqueue cannot poll the agent's terminal: n=%d err=%v data=%d", n, err, out[0].Data)
	}
}
