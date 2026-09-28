package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

type startupErrorWriter struct {
	bytes.Buffer
	printed chan struct{}
}

func (w *startupErrorWriter) Write(p []byte) (int, error) {
	n, err := w.Buffer.Write(p)
	close(w.printed)
	return n, err
}

func TestRunServeWaitsForStartupSettled(t *testing.T) {
	settled := make(chan struct{})
	stderr := &startupErrorWriter{printed: make(chan struct{})}
	result := make(chan int, 1)
	go func() {
		result <- serveResult(&daemon.StartupError{Stage: daemon.StageRegistry, Err: errors.New("bootstrap stalled"), Settled: settled}, stderr)
	}()
	select {
	case <-stderr.printed:
	case <-time.After(time.Second):
		t.Fatal("error not printed before settlement")
	}
	if !strings.Contains(stderr.String(), "bootstrap stalled") {
		t.Fatal("missing startup error")
	}
	select {
	case <-result:
		t.Fatal("serve exited before settlement")
	default:
	}
	close(settled)
	select {
	case code := <-result:
		if code != exitFatal {
			t.Fatalf("exit=%d", code)
		}
	case <-time.After(time.Second):
		t.Fatal("serve did not exit")
	}
}
