//go:build darwin || linux

package main

import (
	"syscall"
	"unsafe"
)

// isTerminal reports whether fd is a terminal: the termios read succeeds only
// on a tty (a character device such as /dev/null fails it).
func isTerminal(fd uintptr) bool {
	var t syscall.Termios
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, ioctlGetTermios, uintptr(unsafe.Pointer(&t)))
	return e == 0
}
