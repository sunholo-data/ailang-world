package agentrun

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// ProcessIdentity reads kernel birth time and excludes a reused PID.
func ProcessIdentity(pid int) (string, bool) {
	return darwinProcessIdentity(pid, unix.SysctlKinfoProc)
}

func darwinProcessIdentity(pid int, read func(string, ...int) (*unix.KinfoProc, error)) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	row, err := read("kern.proc.pid", pid)
	if err != nil || int(row.Proc.P_pid) != pid {
		return "", false
	}
	p := row.Proc
	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
}
