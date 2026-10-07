package daemon

import (
	"io"
	"strings"
	"sync"
	"testing"
)

// testErrorLog is a daemon ErrorLog that writes each operator line to t.Log,
// so a rig failure shows WHY a call refused (row 153 AC1.6) instead of the
// io.Discard default hiding it. It is mutex-guarded (the hostcall goroutine
// writes), and drops writes after the test ends: t.Log after completion
// panics, and a daemon goroutine can outlive its test (C7).
func testErrorLog(t *testing.T) io.Writer {
	t.Helper()
	w := &testLogWriter{t: t}
	t.Cleanup(func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		w.done = true
	})
	return w
}

type testLogWriter struct {
	mu   sync.Mutex
	t    *testing.T
	done bool
}

func (w *testLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.done {
		w.t.Log(strings.TrimRight(string(p), "\n"))
	}
	return len(p), nil
}
