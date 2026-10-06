//go:build linux

package ptytest

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// open is posix_openpt + unlockpt + ptsname, spelled as the two linux ioctls
// libc uses for them (grantpt is a no-op with devpts).
func open() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	fd := master.Fd()
	var unlock int32
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&unlock))); e != 0 {
		_ = master.Close()
		return nil, nil, fmt.Errorf("ptytest: TIOCSPTLCK: %w", e)
	}
	var n uint32
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n))); e != 0 {
		_ = master.Close()
		return nil, nil, fmt.Errorf("ptytest: TIOCGPTN: %w", e)
	}
	slave, err := os.OpenFile("/dev/pts/"+strconv.FormatUint(uint64(n), 10), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}
