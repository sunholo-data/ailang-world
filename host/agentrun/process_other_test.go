//go:build !darwin && !linux

package agentrun

import (
	"os"
	"testing"
)

func TestUnsupportedProcessIdentity(t *testing.T) {
	if birth, live := ProcessIdentity(os.Getpid()); live || birth != "" {
		t.Fatalf("unsupported process identity: %q %v", birth, live)
	}
}
