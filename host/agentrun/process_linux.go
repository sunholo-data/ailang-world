package agentrun

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
)

func ProcessIdentity(pid int) (string, bool) {
	return linuxProcessIdentity(os.DirFS("/proc"), pid)
}

func linuxProcessIdentity(proc fs.FS, pid int) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	stat, err := fs.ReadFile(proc, strconv.Itoa(pid)+"/stat")
	boot, bootErr := fs.ReadFile(proc, "sys/kernel/random/boot_id")
	end := strings.LastIndexByte(string(stat), ')')
	if err != nil || bootErr != nil || end < 0 {
		return "", false
	}
	fields := strings.Fields(string(stat[end+1:]))
	if len(fields) <= 19 || len(strings.TrimSpace(string(boot))) == 0 {
		return "", false
	}
	ticks, err := strconv.ParseUint(fields[19], 10, 64)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(boot)) + ":" + strconv.FormatUint(ticks, 10), fields[0] != "Z"
}
