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
	return newTestLog(t, false)
}

// testOperatorLog is testErrorLog plus the cold-compile tripwire (row 153
// AC3.6): a `capsule: cold compile` line means a rig ran a transition with no
// compile-cache template, i.e. the slow path row 153 removed, so it FAILS the
// test (t.Errorf, not Fatal: the line can arrive off the test goroutine).
// newWSDaemon installs it whenever a test supplies no ErrorLog of its own.
func testOperatorLog(t *testing.T) io.Writer {
	t.Helper()
	return newTestLog(t, true)
}

func newTestLog(t *testing.T, tripwire bool) io.Writer {
	t.Helper()
	w := &testLogWriter{t: t, tripwire: tripwire}
	t.Cleanup(func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		w.done = true
	})
	return w
}

// coldCompileMarker is the operator line a template-less capsule run writes.
const coldCompileMarker = "capsule: cold compile"

type testLogWriter struct {
	mu       sync.Mutex
	t        *testing.T
	done     bool
	tripwire bool
}

func (w *testLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.done {
		line := strings.TrimRight(string(p), "\n")
		w.t.Log(line)
		if w.tripwire && strings.Contains(line, coldCompileMarker) {
			w.t.Errorf("a rig ran a transition with no capsule template (row 153 tripwire): %s", line)
		}
	}
	return len(p), nil
}
