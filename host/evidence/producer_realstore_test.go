package evidence_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/store"
)

// busyWriter wraps a real store as an ObjectWriter and reports the busy window
// it was given. It forwards PutObject unchanged, so only the advertised window
// differs from the store's own.
type busyWriter struct {
	s    *store.Store
	busy time.Duration
}

func (w busyWriter) PutObject(ctx context.Context, o store.Object) error { return w.s.PutObject(ctx, o) }
func (w busyWriter) BusyTimeout() time.Duration                          { return w.busy }

// TestProducerConstructorPinsBusyTimeoutBelowBothTimeouts pins the ordering of
// the producer's two object timeouts against the REAL store windows of both
// seams. Every window below comes from a real store or a forwarding wrapper.
func TestProducerConstructorPinsBusyTimeoutBelowBothTimeouts(t *testing.T) {
	checker := evidence.CheckerConfig{ExecPath: "/nonexistent/ailang", Compiler: testCompiler, CompilerVersion: "AILANG v0.41.0"}

	construct := func(reader evidence.ObjectReader, writer evidence.ObjectWriter, read, write time.Duration) error {
		_, err := evidence.NewProducer(testKey, reader, writer, checker, evidence.ProducerConfig{
			ObjectReadTimeout:  read,
			ObjectWriteTimeout: write,
			ExecTimeout:        time.Minute,
			MaxOutputBytes:     1 << 20,
		})
		return err
	}

	t.Run("ordered", func(t *testing.T) {
		s, _ := openEvidenceStore(t)
		if err := construct(s, s, 3*time.Second, 3*time.Second); err != nil {
			t.Fatalf("ordered configuration (2s default window, 3s timeouts) must construct: %v", err)
		}
	})

	t.Run("read-unordered", func(t *testing.T) {
		s, _ := openEvidenceStore(t)
		err := construct(s, openEvidenceStoreWithoutBusyWait(t), time.Second, 3*time.Second)
		if !errors.Is(err, evidence.ErrUnorderedTimeouts) || errors.Is(err, evidence.ErrInvalidProducerConfig) {
			t.Fatalf("want ErrUnorderedTimeouts; got %v", err)
		}
		if !containsAll(err.Error(), "ObjectReadTimeout") {
			t.Fatalf("detail must name ObjectReadTimeout: %v", err)
		}
	})

	t.Run("write-unordered", func(t *testing.T) {
		s, _ := openEvidenceStore(t)
		err := construct(openEvidenceStoreWithoutBusyWait(t), s, 3*time.Second, time.Second)
		if !errors.Is(err, evidence.ErrUnorderedTimeouts) {
			t.Fatalf("want ErrUnorderedTimeouts; got %v", err)
		}
		if !containsAll(err.Error(), "ObjectWriteTimeout") {
			t.Fatalf("detail must name ObjectWriteTimeout: %v", err)
		}
	})

	t.Run("writer-window-only", func(t *testing.T) {
		s, _ := openEvidenceStore(t)
		err := construct(openEvidenceStoreWithoutBusyWait(t), s, time.Second, time.Second)
		if !errors.Is(err, evidence.ErrUnorderedTimeouts) {
			t.Fatalf("the writer's 2s window must be read: want ErrUnorderedTimeouts; got %v", err)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		s, _ := openEvidenceStore(t)
		err := construct(s, busyWriter{s: s, busy: -time.Nanosecond}, 3*time.Second, 3*time.Second)
		if !errors.Is(err, evidence.ErrInvalidProducerConfig) {
			t.Fatalf("a negative BusyTimeout is unknown: want ErrInvalidProducerConfig; got %v", err)
		}
	})

	t.Run("wrapper-preserves-non-default-window", func(t *testing.T) {
		custom, _ := openEvidenceStoreWithBusyWindow(t, 1500*time.Millisecond)
		if got := custom.BusyTimeout(); got != 1500*time.Millisecond {
			t.Fatalf("instrument: custom store window = %s; want 1500ms", got)
		}
		writer := busyWriter{s: custom, busy: custom.BusyTimeout()}
		if err := construct(openEvidenceStoreWithoutBusyWait(t), writer, 1800*time.Millisecond, 1800*time.Millisecond); err != nil {
			t.Fatalf("1800ms timeouts exceed the forwarded 1500ms window: must construct: %v", err)
		}
		err := construct(openEvidenceStoreWithoutBusyWait(t), writer, 1400*time.Millisecond, 1400*time.Millisecond)
		if !errors.Is(err, evidence.ErrUnorderedTimeouts) {
			t.Fatalf("1400ms timeouts do not exceed the forwarded 1500ms window: want ErrUnorderedTimeouts; got %v", err)
		}
	})
}

func containsAll(s string, needles ...string) bool {
	for _, n := range needles {
		if !strings.Contains(s, n) {
			return false
		}
	}
	return true
}
