package main

import (
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
