package proctest

import (
	"os"
	"syscall"
	"testing"
)

func TestSignalRefusesDangerousPidsWithoutSending(t *testing.T) {
	var sent []int
	record := func(pid int, _ syscall.Signal) error { sent = append(sent, pid); return nil }
	for _, pid := range []int{-1, 0, 1, os.Getpid(), os.Getppid()} {
		if err := Signal(pid, syscall.SIGKILL, record); err == nil {
			t.Errorf("Signal(%d) accepted", pid)
		}
	}
	if len(sent) != 0 {
		t.Fatalf("refused pids reached the send func: %v", sent)
	}
	// Positive control: an ordinary pid is passed through unchanged.
	if err := Signal(424242, syscall.SIGKILL, record); err != nil || len(sent) != 1 || sent[0] != 424242 {
		t.Fatalf("control: err=%v sent=%v", err, sent)
	}
}
