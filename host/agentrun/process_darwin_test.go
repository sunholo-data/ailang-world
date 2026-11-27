package agentrun

import (
	"os"
	"testing"
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
