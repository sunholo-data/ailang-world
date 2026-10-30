package main

import (
	"os"
	"strings"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

// homeServiceRecord is the read-only lookup shared by client pointers.
// Invalid, absent and uninstalled records all mean there is no home service.
func homeServiceRecord() (serviceRecord, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return serviceRecord{}, false
	}
	r, err := readServiceRecord(homeServicePaths(home))
	return r, err == nil && !r.Uninstalled
}

func defaultServiceBind() string {
	if r, ok := homeServiceRecord(); ok {
		return r.Address
	}
	return defaultEpisodeAddr
}

func defaultServiceAddr() string {
	if r, ok := homeServiceRecord(); ok {
		return "http://" + r.Address
	}
	return daemon.DefaultAddr
}

func serviceAddressMatches(addr string, r serviceRecord) bool {
	return strings.TrimSuffix(strings.TrimPrefix(addr, "http://"), "/") == r.Address
}

func pointerSHA7(s string) string {
	s = strings.TrimPrefix(s, "sha256:")
	if len(s) > 7 {
		return s[:7]
	}
	return s
}
