package agentrun

import (
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestDarwinProcessIdentity(t *testing.T) {
	pid := os.Getpid()
	birth, live := ProcessIdentity(pid)
	if !live || birth == "" {
		t.Fatalf("live process missing birth: %q %v", birth, live)
	}
	if again, live := ProcessIdentity(pid); !live || again != birth {
		t.Fatalf("unstable birth: %q %v", again, live)
	}
	for _, pid := range []int{0, -1, 1 << 30} {
		if birth, live := ProcessIdentity(pid); live || birth != "" {
			t.Fatalf("invalid or absent PID %d: %q %v", pid, birth, live)
		}
	}
}

func TestDarwinZombieIdentity(t *testing.T) {
	var row unix.KinfoProc
	row.Proc.P_pid = 123
	row.Proc.P_starttime.Sec = 456
	read := func(string, ...int) (*unix.KinfoProc, error) { return &row, nil }
	row.Proc.P_stat = 2 // SRUN in Darwin <sys/proc.h>
	if birth, live := darwinProcessIdentity(123, read); birth != "456.0" || !live {
		t.Fatalf("live control: %q %t", birth, live)
	}
	row.Proc.P_stat = 5 // SZOMB in Darwin <sys/proc.h>
	if birth, live := darwinProcessIdentity(123, read); birth != "456.0" || live {
		t.Fatalf("a zombie was live: %q %t", birth, live)
	}
}
