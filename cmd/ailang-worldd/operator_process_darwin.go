package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func listWorldProcesses() ([]worldProcess, error) {
	rows, err := unix.SysctlKinfoProcSlice("kern.proc.uid", os.Getuid())
	if err != nil {
		return nil, err
	}
	out := make([]worldProcess, 0, len(rows))
	for _, row := range rows {
		p := row.Proc
		out = append(out, worldProcess{int(p.P_pid), unix.ByteSliceToString(p.P_comm[:])})
	}
	return out, nil
}

// operatorProcessIdentity reads kernel birth time and excludes a reused PID.
func operatorProcessIdentity(pid int) (string, bool) {
	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil || int(row.Proc.P_pid) != pid {
		return "", false
	}
	p := row.Proc
	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
}
