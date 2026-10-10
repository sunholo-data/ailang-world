package broker

// Row 140 M2 (w-workspace-exec-toolchain-effect §4.5 "One process lifecycle"):
// runBounded's capture strategy. The zero value is today's strict capture and
// is pinned by every pre-existing runBounded and handler test, unedited
// (MUT-STRICT-REGRESSED); these arms cover only the new captureHeadTail, whose
// one user is ExecHandler.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/procbound"
)

// repeatStream is what `yes 0123456789abcdef | head -c n` writes.
func repeatStream(n int) []byte {
	line := []byte("0123456789abcdef\n")
	out := bytes.Repeat(line, n/len(line)+1)
	return out[:n]
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func TestRunBoundedHeadTailKeepsHeadTailTotalAndDigest(t *testing.T) {
	const n = 2000000
	script := writeExecutable(t, "yes 0123456789abcdef | head -c 2000000\nprintf 'ERR-HEAD' >&2\nyes e | head -c 300000 >&2")
	capture := newCaptureHeadTail(execHeadBytes, execTailBytes, 64<<20)
	out, err := runBounded(boundedTestContext(t), handlerBounds{execTimeout: 20 * time.Second},
		handlerCommand{path: script, dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}, capture: capture})
	if err != nil {
		t.Fatalf("runBounded: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("head/tail capture returned %d stdout bytes from runBounded; the sinks hold the streams", len(out))
	}
	full := repeatStream(n)
	got := capture.stdout.snapshot()
	if got.Bytes != n || !got.Truncated || got.SHA256 != sha256Hex(full) {
		t.Fatalf("stdout snapshot bytes=%d truncated=%t sha=%s, want %d true %s", got.Bytes, got.Truncated, got.SHA256, n, sha256Hex(full))
	}
	wantText := string(full[:execHeadBytes]) + truncationMarker(n-execHeadBytes-execTailBytes) + string(full[n-execTailBytes:])
	if got.Text != wantText {
		t.Fatalf("stdout text is not head + marker + tail (len %d, want %d)", len(got.Text), len(wantText))
	}
	errSnap := capture.stderr.snapshot()
	if errSnap.Bytes != 300008 || !errSnap.Truncated || !strings.HasPrefix(errSnap.Text, "ERR-HEAD") {
		t.Fatalf("stderr snapshot bytes=%d truncated=%t head=%q", errSnap.Bytes, errSnap.Truncated, errSnap.Text[:16])
	}
}

func TestRunBoundedHeadTailShortStreamIsWholeAndUntruncated(t *testing.T) {
	script := writeExecutable(t, "printf 'hello\\n'; printf 'warn' >&2; exit 3")
	capture := newCaptureHeadTail(execHeadBytes, execTailBytes, 64<<20)
	_, err := runBounded(boundedTestContext(t), handlerBounds{execTimeout: 20 * time.Second},
		handlerCommand{path: script, dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}, capture: capture})
	var exitErr *HandlerExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("err = %v, want *HandlerExitError", err)
	}
	if code, ok := handlerExitStatus(err); !ok || code != 3 {
		t.Fatalf("handlerExitStatus = %d, %t; want 3, true", code, ok)
	}
	if s := capture.stdout.snapshot(); s.Text != "hello\n" || s.Bytes != 6 || s.Truncated || s.SHA256 != sha256Hex([]byte("hello\n")) {
		t.Fatalf("stdout snapshot = %+v", s)
	}
	if s := capture.stderr.snapshot(); s.Text != "warn" || s.Truncated {
		t.Fatalf("stderr snapshot = %+v", s)
	}
	// Exactly head+tail bytes is still whole: nothing is dropped, so no marker.
	exact := repeatStream(execHeadBytes + execTailBytes)
	sink := newHeadTailSink(execHeadBytes, execTailBytes, 64<<20, nil)
	for i := 0; i < len(exact); i += 1000 {
		_, _ = sink.Write(exact[i:min(i+1000, len(exact))])
	}
	if s := sink.snapshot(); s.Truncated || s.Text != string(exact) {
		t.Fatalf("a head+tail-sized stream was truncated or altered (truncated=%t len=%d)", s.Truncated, len(s.Text))
	}
}

// A runaway stream is killed through the strict path's overflow kill and
// reported as *HandlerOutputOverflowError: the cap is per stream.
func TestRunBoundedHeadTailHardCapKillsTheGroup(t *testing.T) {
	for _, stream := range []string{"stdout", "stderr"} {
		t.Run(stream, func(t *testing.T) {
			redirect := ""
			if stream == "stderr" {
				redirect = " >&2"
			}
			script := writeExecutable(t, "yes runaway"+redirect)
			capture := newCaptureHeadTail(execHeadBytes, execTailBytes, 1<<20)
			start := time.Now()
			_, err := runBounded(boundedTestContext(t), handlerBounds{execTimeout: 20 * time.Second},
				handlerCommand{path: script, dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}, capture: capture})
			var overflow *HandlerOutputOverflowError
			if !errors.As(err, &overflow) || overflow.Limit != 1<<20 {
				t.Fatalf("err = %v, want *HandlerOutputOverflowError{1 MiB}", err)
			}
			if elapsed := time.Since(start); elapsed > 10*time.Second {
				t.Fatalf("overflow took %s; the kill did not fire on the cap", elapsed)
			}
		})
	}
}

// On the deadline the partial head/tail stays readable.
func TestRunBoundedHeadTailTimeoutKeepsPartialOutput(t *testing.T) {
	ctx := boundedTestContext(t)
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) <= 0 {
		t.Fatal("finite caller deadline missing/expired")
	}
	dir := t.TempDir()
	ready := filepath.Join(dir, "ready")
	script := writeExecutable(t, `printf 'PARTIAL-OUT'; printf 'PARTIAL-ERR' >&2
: > "$READY"
if [ "$HOLD" -gt 0 ]; then sleep "$HOLD"; fi`)
	assertStreams := func(c *captureHeadTail) {
		t.Helper()
		for name, want := range map[string]string{"stdout": "PARTIAL-OUT", "stderr": "PARTIAL-ERR"} {
			sink := c.stdout
			if name == "stderr" {
				sink = c.stderr
			}
			if s := sink.snapshot(); s.Text != want || s.Bytes != int64(len(want)) || s.Truncated || s.SHA256 != sha256Hex([]byte(want)) {
				t.Fatalf("exact partial %s Text/Bytes/SHA/untruncated: %+v", name, s)
			}
		}
	}
	originalKill := killGroup
	var mu sync.Mutex
	type observation struct {
		pgid  int
		ready bool
		err   error
	}
	var kills []observation
	// Wait owns the reservation until the actual cmd.Wait reap/pipe-copy join.
	// On any failure keep the seam installed until that reservation is released.
	reap := func(base int64) bool {
		t.Helper()
		if procbound.Outstanding() != base {
			mu.Lock()
			for _, k := range kills {
				_ = originalKill(k.pgid)
			}
			mu.Unlock()
			for procbound.Outstanding() != base && time.Until(deadline) > 0 {
				time.Sleep(time.Millisecond)
			}
		}
		if got := procbound.Outstanding(); got != base {
			t.Errorf("incomplete actual reap: Outstanding=%d want %d", got, base)
			return false
		}
		return true
	}
	run := func(bound, hold time.Duration) (*captureHeadTail, error, time.Duration) {
		t.Helper()
		if err := os.Remove(ready); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		base := procbound.Outstanding()
		defer reap(base)
		capture := newCaptureHeadTail(execHeadBytes, execTailBytes, 64<<20)
		start := time.Now()
		_, err := runBounded(ctx, handlerBounds{execTimeout: bound}, handlerCommand{path: script, dir: dir,
			env: []string{"PATH=/usr/bin:/bin", "READY=" + ready, fmt.Sprintf("HOLD=%d", hold/time.Second)}, capture: capture})
		return capture, err, time.Since(start)
	}
	var C time.Duration
	for i := 0; i < 3; i++ {
		capture, err, elapsed := run(20*time.Second, 0)
		t.Logf("no-sleep control %d elapsed=%s", i+1, elapsed)
		if err != nil {
			t.Fatalf("no-sleep control %d: %v", i+1, err)
		}
		assertStreams(capture)
		if _, err := os.Stat(ready); err != nil {
			t.Fatalf("control readiness after both printf: %v", err)
		}
		C = max(C, elapsed)
	}
	const maxDuration = time.Duration(1<<63 - 1)
	if C <= 0 || C > (maxDuration-2*pipeCloseGrace-time.Millisecond)/10 {
		t.Fatalf("invalid/overflow calibration C=%s", C)
	}
	B := ((10*C + 2*pipeCloseGrace + time.Millisecond - 1) / time.Millisecond) * time.Millisecond
	if B <= 0 || B > (maxDuration-2*pipeCloseGrace-time.Second)/3 {
		t.Fatalf("invalid/overflow timeout B=%s", B)
	}
	H := ((3*B + 2*pipeCloseGrace + time.Second - 1) / time.Second) * time.Second
	R := time.Until(deadline)
	t.Logf("calibration C=%s B=%s H=%s R=%s grace=%s", C, B, H, R, pipeCloseGrace)
	// Subtractions avoid overflow when checking the full startup/hold/reap allowance.
	if H < 3*B+2*pipeCloseGrace || B >= R || H >= R-B || pipeCloseGrace >= R-B-H {
		t.Fatalf("finite fixture guard: B+H+grace must fit remaining caller deadline (C=%s B=%s H=%s R=%s)", C, B, H, R)
	}
	base := procbound.Outstanding()
	t.Cleanup(func() {
		if reap(base) {
			killGroup = originalKill
		}
	})
	killGroup = func(pgid int) error {
		mu.Lock()
		defer mu.Unlock()
		_, statErr := os.Stat(ready) // Observe before delegating the real process-group kill.
		err := originalKill(pgid)
		kills = append(kills, observation{pgid, statErr == nil, err})
		return err
	}
	capture, err, elapsed := run(B, H)
	mu.Lock()
	observed := append([]observation(nil), kills...)
	mu.Unlock()
	t.Logf("deadline arm elapsed=%s err=%v kills=%+v Outstanding=%d baseline=%d", elapsed, err, observed, procbound.Outstanding(), base)
	t.Logf("deadline snapshots stdout=%+v stderr=%+v", capture.stdout.snapshot(), capture.stderr.snapshot())
	var timeout *HandlerTimeoutError
	if !errors.Is(err, ErrHandlerTimeout) || !errors.As(err, &timeout) || timeout.Timeout != B {
		t.Fatalf("exact timeout B=%s required: err=%v real-kills=%d", B, err, len(observed))
	}
	if len(observed) == 0 {
		t.Fatal("real deadline group kill missing")
	}
	for _, k := range observed {
		if !k.ready {
			t.Fatalf("readiness before real deadline kill missing: %+v", k)
		}
	}
	assertStreams(capture)
}
