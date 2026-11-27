package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func listWorldProcesses() ([]worldProcess, error) {
	return nil, errors.New("operator upgrade requires the macOS process inventory")
}

func operatorProcessIdentity(pid int) (string, bool) {
	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	boot, bootErr := os.ReadFile("/proc/sys/kernel/random/boot_id")
	end := strings.LastIndexByte(string(stat), ')')
	if err != nil || bootErr != nil || end < 0 {
		return "", false
	}
	fields := strings.Fields(string(stat[end+1:]))
	if len(fields) <= 19 {
		return "", false
	}
	return strings.TrimSpace(string(boot)) + ":" + fields[19], fields[0] != "Z"
}
