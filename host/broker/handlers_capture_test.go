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
	"strings"
	"testing"
	"time"
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
	script := writeExecutable(t, "printf 'PARTIAL-OUT'; printf 'PARTIAL-ERR' >&2; sleep 30")
	capture := newCaptureHeadTail(execHeadBytes, execTailBytes, 64<<20)
	_, err := runBounded(boundedTestContext(t), handlerBounds{execTimeout: 300 * time.Millisecond},
		handlerCommand{path: script, dir: t.TempDir(), env: []string{"PATH=/usr/bin:/bin"}, capture: capture})
	if !errors.Is(err, ErrHandlerTimeout) {
		t.Fatalf("err = %v, want ErrHandlerTimeout", err)
	}
	if s := capture.stdout.snapshot(); s.Text != "PARTIAL-OUT" {
		t.Fatalf("partial stdout = %q", s.Text)
	}
	if s := capture.stderr.snapshot(); s.Text != "PARTIAL-ERR" {
		t.Fatalf("partial stderr = %q", s.Text)
	}
}
