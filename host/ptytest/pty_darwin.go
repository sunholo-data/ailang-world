//go:build darwin

package ptytest

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// open is posix_openpt + grantpt + unlockpt + ptsname, spelled as the three
// darwin ioctls libc uses for them.
func open() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	fd := master.Fd()
	for _, req := range []uintptr{syscall.TIOCPTYGRANT, syscall.TIOCPTYUNLK} {
		if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, req, 0); e != 0 {
			_ = master.Close()
			return nil, nil, fmt.Errorf("ptytest: ioctl %#x: %w", req, e)
		}
	}
	var name [128]byte
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&name[0]))); e != 0 {
		_ = master.Close()
		return nil, nil, fmt.Errorf("ptytest: TIOCPTYGNAME: %w", e)
	}
	path := string(name[:bytes.IndexByte(name[:], 0)])
	slave, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}
