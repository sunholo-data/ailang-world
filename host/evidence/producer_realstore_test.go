package evidence_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// busyWriter wraps a real store as an ObjectWriter and reports the busy window
// it was given. It forwards PutObject unchanged, so only the advertised window
// differs from the store's own.
type busyWriter struct {
	s    *store.Store
	busy time.Duration
}

func (w busyWriter) PutObject(ctx context.Context, o store.Object) error {
	return w.s.PutObject(ctx, o)
}
func (w busyWriter) BusyTimeout() time.Duration { return w.busy }

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

// pinnedChecker pins the real AILANG_BIN by its bytes and its version line. The
// external-package twin of the internal helper: it fails loudly and never skips.
func pinnedChecker(t *testing.T) evidence.CheckerConfig {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the producer's binary arms need the pinned v0.41.0 release; never skip")
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("pinned binary --version: %v", err)
	}
	if line := strings.SplitN(string(out), "\n", 2)[0]; line != "AILANG v0.41.0" {
		t.Fatalf("AILANG_BIN reports %q; want AILANG v0.41.0", line)
	}
	data, err := os.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	return evidence.CheckerConfig{ExecPath: bin, Compiler: hashref.SumSHA256(data), CompilerVersion: "AILANG v0.41.0"}
}

// countingWriter forwards to a real store and counts attempts. A non-nil fail is
// returned without touching the store, so the write leg is refused before commit.
type countingWriter struct {
	s     *store.Store
	calls atomic.Int32
	fail  error
}

func (w *countingWriter) PutObject(ctx context.Context, o store.Object) error {
	w.calls.Add(1)
	if w.fail != nil {
		return w.fail
	}
	return w.s.PutObject(ctx, o)
}
func (w *countingWriter) BusyTimeout() time.Duration { return w.s.BusyTimeout() }

// tamperReader returns the real metadata with the payload's last byte flipped:
// a source whose bytes no longer hash to the ref that was asked for.
type tamperReader struct{ s *store.Store }

func (r tamperReader) ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (store.ObjectMeta, []byte, error) {
	meta, payload, err := r.s.ReadObject(ctx, ref, maxBytes)
	if err != nil || len(payload) == 0 {
		return meta, payload, err
	}
	tampered := append([]byte(nil), payload...)
	tampered[len(tampered)-1] ^= 0xff
	return meta, tampered, nil
}
func (r tamperReader) BusyTimeout() time.Duration { return r.s.BusyTimeout() }

// uncertainWriter commits through the real store and then reports the outcome as
// uncertain, the way a timed-out commit can. It records the hash it was handed.
type uncertainWriter struct {
	s     *store.Store
	calls atomic.Int32
	hash  hashref.HashRef
}

func (w *uncertainWriter) PutObject(ctx context.Context, o store.Object) error {
	w.calls.Add(1)
	w.hash = o.Hash
	if err := w.s.PutObject(ctx, o); err != nil {
		return err
	}
	return &store.UncertainError{Op: "put object", Cause: context.DeadlineExceeded}
}
func (w *uncertainWriter) BusyTimeout() time.Duration { return w.s.BusyTimeout() }

// putSource stores a source object the producer can read back and verify.
func putSource(t *testing.T, s *store.Store, src []byte) hashref.HashRef {
	t.Helper()
	ref := hashref.SumSHA256(src)
	if err := s.PutObject(boundedTestContext(t), store.Object{
		Hash: ref, InterfaceHash: hashref.SumSHA256([]byte("test/ail-source")),
		SemanticID: "test/ail-source", Provenance: "producer-test", Payload: src,
	}); err != nil {
		t.Fatal(err)
	}
	return ref
}

// producerWith builds the ordered production configuration: 3s object timeouts
// above the 2s default store window, a 2 minute checker, and 1 MiB of output.
func producerWith(t *testing.T, reader evidence.ObjectReader, writer evidence.ObjectWriter, checker evidence.CheckerConfig) *evidence.Producer {
	t.Helper()
	p, err := evidence.NewProducer(testKey, reader, writer, checker, evidence.ProducerConfig{
		ObjectReadTimeout:  3 * time.Second,
		ObjectWriteTimeout: 3 * time.Second,
		ExecTimeout:        2 * time.Minute,
		MaxOutputBytes:     1 << 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// TestProducerEndToEndMintsEnvelopeValidatorAccepts runs the whole producer against
// the real store and the pinned binary, then has the production validator accept
// what it minted. Each refusal arm must leave the writer untouched.
func TestProducerEndToEndMintsEnvelopeValidatorAccepts(t *testing.T) {
	s, _ := openEvidenceStore(t)
	pinned := pinnedChecker(t)
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	srcRef := putSource(t, s, source)
	ctx := context.Background()
	required := []string{"twicePlusOne"}

	t.Run("mint", func(t *testing.T) {
		writer := &countingWriter{s: s}
		ref, err := producerWith(t, s, writer, pinned).GenerateProof(ctx, srcRef, required)
		if err != nil || ref.IsZero() || writer.calls.Load() != 1 {
			t.Fatalf("the verified source must mint once: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
		v, err := evidence.NewValidator(testKey, s, evidence.CompilerConfig{
			Compiler: pinned.Compiler, CompilerVersion: "AILANG v0.41.0", ObjectReadTimeout: 3 * time.Second,
		}, required)
		if err != nil {
			t.Fatal(err)
		}
		sealed, ok := v.ValidateProof(boundedTestContext(t), ref, srcRef).Validated()
		if !ok {
			t.Fatal("the production validator must accept the producer's envelope")
		}
		if grade, proven := v.Resolve(sealed).Proven(); !proven || grade != evidence.ResolvedGradeProven {
			t.Fatalf("the sealed evidence must resolve as proven: grade=%v proven=%v", grade, proven)
		}
	})
	t.Run("idempotent", func(t *testing.T) {
		first, err := producerWith(t, s, &countingWriter{s: s}, pinned).GenerateProof(ctx, srcRef, required)
		if err != nil {
			t.Fatal(err)
		}
		second, err := producerWith(t, s, &countingWriter{s: s}, pinned).GenerateProof(ctx, srcRef, required)
		if err != nil || second != first {
			t.Fatalf("a repeated mint must return the identical ref: first=%s second=%s err=%v", first, second, err)
		}
	})
	t.Run("tampered-compiler", func(t *testing.T) {
		bad := pinned
		bad.Compiler = hashref.SumSHA256([]byte("not the pinned binary"))
		writer := &countingWriter{s: s}
		ref, err := producerWith(t, s, writer, bad).GenerateProof(ctx, srcRef, required)
		if !errors.Is(err, evidence.ErrUntrustedTool) || !ref.IsZero() || writer.calls.Load() != 0 {
			t.Fatalf("a tampered checker must be refused with no write: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
	})
	t.Run("absent-identity", func(t *testing.T) {
		writer := &countingWriter{s: s}
		ref, err := producerWith(t, s, writer, pinned).GenerateProof(ctx, srcRef, []string{"absentIdentity"})
		if !errors.Is(err, evidence.ErrRequiredIdentityMissing) || !ref.IsZero() || writer.calls.Load() != 0 {
			t.Fatalf("an unverified identity must be refused with no write: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
	})
	t.Run("source-hash-mismatch", func(t *testing.T) {
		writer := &countingWriter{s: s}
		ref, err := producerWith(t, tamperReader{s: s}, writer, pinned).GenerateProof(ctx, srcRef, required)
		if err == nil || !strings.Contains(err.Error(), "hash to") || !ref.IsZero() || writer.calls.Load() != 0 {
			t.Fatalf("source bytes that miss their ref must be refused with no write: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
	})
	t.Run("missing-source", func(t *testing.T) {
		writer := &countingWriter{s: s}
		missing := hashref.SumSHA256([]byte("never stored"))
		ref, err := producerWith(t, s, writer, pinned).GenerateProof(ctx, missing, required)
		if err == nil || !strings.Contains(err.Error(), "not found") || !ref.IsZero() || writer.calls.Load() != 0 {
			t.Fatalf("an absent source must be refused by name with no write: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
	})
	t.Run("oversize-source", func(t *testing.T) {
		big := putSource(t, s, make([]byte, evidence.MaxSourceBytes+1))
		writer := &countingWriter{s: s}
		ref, err := producerWith(t, s, writer, pinned).GenerateProof(ctx, big, required)
		var tooLarge *store.ObjectTooLargeError
		if !errors.As(err, &tooLarge) || !ref.IsZero() || writer.calls.Load() != 0 {
			t.Fatalf("a source past MaxSourceBytes must be refused as too large: ref=%s err=%v calls=%d", ref, err, writer.calls.Load())
		}
	})
	t.Run("write-error", func(t *testing.T) {
		injected := errors.New("injected not-committed")
		writer := &countingWriter{s: s, fail: injected}
		ref, err := producerWith(t, s, writer, pinned).GenerateProof(ctx, srcRef, required)
		if !errors.Is(err, injected) || !ref.IsZero() {
			t.Fatalf("a failed write must surface its error with a zero ref: ref=%s err=%v", ref, err)
		}
	})
}

// TestProducerMintsValidatorReadableEnvelope reads the minted object back through the
// store and checks every field the validator and a reviewer rely on.
func TestProducerMintsValidatorReadableEnvelope(t *testing.T) {
	s, _ := openEvidenceStore(t)
	pinned := pinnedChecker(t)
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	srcRef := putSource(t, s, source)
	writer := &countingWriter{s: s}
	ref, err := producerWith(t, s, writer, pinned).GenerateProof(context.Background(), srcRef, []string{"twicePlusOne"})
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	meta, payload, err := s.ReadObject(boundedTestContext(t), ref, evidence.MaxBytes)
	if err != nil {
		t.Fatalf("the minted object must read back: %v", err)
	}
	if meta.SemanticID != evidence.ProofSemanticID || meta.InterfaceHash != evidence.InterfaceHashV1 {
		t.Fatalf("stored identity %q/%s is not the validator's proof identity", meta.SemanticID, meta.InterfaceHash)
	}
	if hashref.SumSHA256(payload) != ref {
		t.Fatalf("the payload does not hash to the returned ref %s", ref)
	}
	env, err := evidence.DecodeAuthenticatedEnvelope(payload)
	if err != nil {
		t.Fatalf("payload must decode as an envelope: %v", err)
	}
	report, err := evidence.DecodeProofReportV1(env.Report)
	if err != nil {
		t.Fatalf("report must decode: %v", err)
	}
	if report.Subject != srcRef || report.Compiler != pinned.Compiler || report.CompilerVersion != "AILANG v0.41.0" ||
		strings.Join(report.Verified, ",") != "twicePlusOne" || report.Errors != 0 || report.Counterexamples != 0 ||
		!report.CheckPassed || !report.ProofSucceeded {
		t.Fatalf("minted report fields are wrong: %+v", report)
	}
	mac := hmac.New(sha256.New, testKey[:])
	mac.Write(env.Report)
	if !hmac.Equal(mac.Sum(nil), env.MAC) {
		t.Fatal("the envelope MAC does not authenticate its report")
	}
	if writer.calls.Load() != 1 {
		t.Fatalf("a clean mint writes exactly once; calls=%d", writer.calls.Load())
	}
}

// TestProducerSurfacesUncertainWriteOutcome: an uncertain write is reported as
// uncertain, never retried, and the committed object is reachable for reconcile.
func TestProducerSurfacesUncertainWriteOutcome(t *testing.T) {
	s, _ := openEvidenceStore(t)
	pinned := pinnedChecker(t)
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	srcRef := putSource(t, s, source)
	writer := &uncertainWriter{s: s}
	ref, err := producerWith(t, s, writer, pinned).GenerateProof(context.Background(), srcRef, []string{"twicePlusOne"})
	var uncertain *store.UncertainError
	if !errors.As(err, &uncertain) || !ref.IsZero() {
		t.Fatalf("an uncertain write must surface as *store.UncertainError with a zero ref: ref=%s err=%v", ref, err)
	}
	if writer.calls.Load() != 1 {
		t.Fatalf("an uncertain write must never be retried blind; calls=%d", writer.calls.Load())
	}
	if _, ok, err := s.GetObject(boundedTestContext(t), writer.hash); err != nil || !ok {
		t.Fatalf("the uncertain commit must be present for the caller's reconcile: ok=%v err=%v", ok, err)
	}
}

// decoyStimulus pins the store's one pooled connection with a decoy read of a
// measured hold, so a store call made inside the window must wait behind it.
type decoyStimulus struct {
	t      *testing.T
	s      *store.Store
	decoy  hashref.HashRef
	settle time.Duration
}

// attemptRecord is what an armed wrapper observed. It is written on the
// producer's goroutine, so it never calls t.Fatal: the test classifies it after
// the producer returns.
type attemptRecord struct {
	decoyCh     <-chan time.Duration
	decoyTook   time.Duration
	decoyJoined bool
	callElapsed time.Duration
	instrument  string
	wouldBe     hashref.HashRef
}

// join waits for the decoy to finish, once, and returns how long it held.
func (r *attemptRecord) join() time.Duration {
	if !r.decoyJoined {
		r.decoyTook = <-r.decoyCh
		r.decoyJoined = true
	}
	return r.decoyTook
}

// arm starts the decoy, lets it settle, and checks the stimulus is live: the
// decoy must still be in flight and the caller's deadline must be intact.
func (d decoyStimulus) arm(ctx context.Context, rec *attemptRecord) {
	rec.decoyCh = startDecoyRead(d.t, d.s, d.decoy)
	time.Sleep(d.settle)
	select {
	case took := <-rec.decoyCh:
		rec.decoyTook, rec.decoyJoined = took, true
		rec.instrument = fmt.Sprintf("decoy finished in %s before the store call; settle %s is too long for its hold", took, d.settle)
	default:
	}
	if rec.instrument == "" && ctx.Err() != nil {
		rec.instrument = "caller deadline was spent before the store call; the settle consumed the bound"
	}
}

// decoyReader is the real store behind a read seam. Unarmed, it is transparent.
// Armed, it pins the connection first and times ONLY the real ReadObject call.
type decoyReader struct {
	s   *store.Store
	d   decoyStimulus
	rec *attemptRecord
}

func (r decoyReader) ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (store.ObjectMeta, []byte, error) {
	if r.rec == nil {
		return r.s.ReadObject(ctx, ref, maxBytes)
	}
	r.d.arm(ctx, r.rec)
	start := time.Now()
	meta, payload, err := r.s.ReadObject(ctx, ref, maxBytes)
	r.rec.callElapsed = time.Since(start)
	return meta, payload, err
}
func (r decoyReader) BusyTimeout() time.Duration { return r.s.BusyTimeout() }

// decoyWriter is the real store behind a write seam. Armed, it records the hash it
// was handed, pins the connection, and times ONLY the real PutObject call.
type decoyWriter struct {
	s   *store.Store
	d   decoyStimulus
	rec *attemptRecord
}

func (w decoyWriter) PutObject(ctx context.Context, o store.Object) error {
	if w.rec == nil {
		return w.s.PutObject(ctx, o)
	}
	w.rec.wouldBe = o.Hash
	w.d.arm(ctx, w.rec)
	start := time.Now()
	err := w.s.PutObject(ctx, o)
	w.rec.callElapsed = time.Since(start)
	return err
}
func (w decoyWriter) BusyTimeout() time.Duration { return w.s.BusyTimeout() }

// generateWithin runs GenerateProof on its own goroutine under a watchdog. The
// watchdog catches only a true hang: a correct mint includes hashing the binary
// and two subprocesses, so it is set far above any deadline the test measures.
func generateWithin(t *testing.T, p *evidence.Producer, src hashref.HashRef, required []string, watchdog time.Duration) (hashref.HashRef, error) {
	t.Helper()
	type result struct {
		ref hashref.HashRef
		err error
	}
	done := make(chan result, 1)
	go func() {
		ref, err := p.GenerateProof(context.Background(), src, required)
		done <- result{ref, err}
	}()
	select {
	case r := <-done:
		return r.ref, r.err
	case <-time.After(watchdog):
		t.Fatalf("GenerateProof did not return within the %s watchdog: a true hang, not a deadline", watchdog)
	}
	return hashref.HashRef{}, nil
}

// requireDeadlineRefusal is the only PASS for an armed attempt: a deadline refusal
// that is not a missing-deadline programming error, with a zero ref.
func requireDeadlineRefusal(t *testing.T, ref hashref.HashRef, err error) {
	t.Helper()
	if !errors.Is(err, context.DeadlineExceeded) || errors.Is(err, store.ErrNoDeadline) || !ref.IsZero() {
		t.Fatalf("a blocked store leg must refuse with DeadlineExceeded and a zero ref: ref=%s err=%v", ref, err)
	}
}

// TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout pins the source read
// leg against a real store whose one connection a decoy holds. The bound is
// derived from the decoy's measured hold, so the stimulus holds for more than
// 20x the bound on any machine. The decoy holds the connection and claims nothing
// about SQLite lock contention.
// controlBound is the generous bound for no-decoy control arms. A control that
// shares the derived bound's ~10 ms scale measures scheduler load under the full
// parallel suite, not the producer.
const controlBound = 30 * time.Second

func TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout(t *testing.T) {
	// The stimulus needs at least two Ps, and that is a measurement rather than
	// a precaution: this driver's SQLite is pure Go, so at GOMAXPROCS=1 the decoy
	// monopolizes the single P for its whole ~53 ms read and the blocked reader
	// is never scheduled to join the pool queue at all -- it runs only after the
	// connection is free, wins it, and exhausts the retry budget. Measured on
	// unmutated code: 2/10 retry-exhaustion reds at GOMAXPROCS=1, 0/10 at the
	// default. Declare the instrument's precondition and enforce it instead of
	// inheriting whatever the harness happens to set.
	if prev := runtime.GOMAXPROCS(0); prev < 2 {
		// Announce it. Every other precondition in this file is loud, and a
		// precondition that silently corrects itself is one nobody can audit.
		t.Logf("raising GOMAXPROCS %d -> 2 for the duration of this test: the stimulus needs two Ps to exist concurrently", prev)
		defer runtime.GOMAXPROCS(prev)
		runtime.GOMAXPROCS(2)
	}
	s := openEvidenceStoreWithoutBusyWait(t)
	pinned := pinnedChecker(t)
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	srcRef := putSource(t, s, source)
	decoyRef := putEvidenceObject(t, s, []byte(strings.Repeat("d", 256<<20)))
	hold := measureGetObject(t, s, decoyRef)
	if hold < minDecoyHold {
		t.Fatalf("instrument failure: decoy held the connection for only %v; want at least %v for this stimulus to exist", hold, minDecoyHold)
	}
	bound := hold / decoyHoldRatio
	settle := hold / 40
	watchdog := 20*hold + 60*time.Second
	required := []string{"twicePlusOne"}
	t.Logf("decoy hold=%v -> derived ObjectReadTimeout=%v settle=%v watchdog=%v", hold, bound, settle, watchdog)

	// No-decoy control first: the same producer mints when nothing pins the connection.
	// The control proves only that the producer and the real store mint end to end
	// with no decoy. It makes no claim about the derived bound, so it gets a generous
	// read bound: a deadline here would measure load, not the stimulus.
	ctrlP, err := evidence.NewProducer(testKey, decoyReader{s: s}, &countingWriter{s: s}, pinned, evidence.ProducerConfig{
		ObjectReadTimeout: controlBound, ObjectWriteTimeout: 3*time.Second + 2*hold, ExecTimeout: 2 * time.Minute, MaxOutputBytes: 1 << 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ref, err := ctrlP.GenerateProof(context.Background(), srcRef, required); err != nil || ref.IsZero() {
		t.Fatalf("the no-decoy control must mint end to end with no decoy: ref=%s err=%v", ref, err)
	}

	d := decoyStimulus{t: t, s: s, decoy: decoyRef, settle: settle}
	for attempt := 1; attempt <= 5; attempt++ {
		rec := &attemptRecord{}
		writer := &countingWriter{s: s}
		p, err := evidence.NewProducer(testKey, decoyReader{s: s, d: d, rec: rec}, writer, pinned, evidence.ProducerConfig{
			ObjectReadTimeout: bound, ObjectWriteTimeout: 3*time.Second + 2*hold, ExecTimeout: 2 * time.Minute, MaxOutputBytes: 1 << 20,
		})
		if err != nil {
			t.Fatal(err)
		}
		ref, err := generateWithin(t, p, srcRef, required, watchdog)
		rec.join()
		switch {
		case rec.instrument != "":
			t.Logf("attempt %d: instrument retry: %s", attempt, rec.instrument)
		case err != nil:
			requireDeadlineRefusal(t, ref, err)
			if writer.calls.Load() != 0 {
				t.Fatalf("a blocked source read must never reach the write: writes=%d", writer.calls.Load())
			}
			t.Logf("blocked read refused after store call %v (derived bound %v)", rec.callElapsed, bound)
			return
		case rec.callElapsed < 2*bound:
			t.Logf("attempt %d: the store call won the freed connection in %v; retrying", attempt, rec.callElapsed)
		default:
			t.Fatalf("blocked source read minted after exceeding ObjectReadTimeout: store call=%v; want a deadline refusal", rec.callElapsed)
		}
	}
	t.Fatal("instrument failure: the store call won the freed connection in all 5 attempts")
}

// TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout pins the envelope
// write leg against a real store whose one connection a decoy holds. The bound is
// derived from the decoy's measured hold, so the stimulus holds for more than 20x
// the bound on any machine. The stimulus holds NO SQLite write lock and claims
// nothing about lock contention: it only pins the store's pooled connection.
func TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout(t *testing.T) {
	// The stimulus needs at least two Ps, and that is a measurement rather than
	// a precaution: this driver's SQLite is pure Go, so at GOMAXPROCS=1 the decoy
	// monopolizes the single P for its whole ~53 ms read and the blocked reader
	// is never scheduled to join the pool queue at all -- it runs only after the
	// connection is free, wins it, and exhausts the retry budget. Measured on
	// unmutated code: 2/10 retry-exhaustion reds at GOMAXPROCS=1, 0/10 at the
	// default. Declare the instrument's precondition and enforce it instead of
	// inheriting whatever the harness happens to set.
	if prev := runtime.GOMAXPROCS(0); prev < 2 {
		// Announce it. Every other precondition in this file is loud, and a
		// precondition that silently corrects itself is one nobody can audit.
		t.Logf("raising GOMAXPROCS %d -> 2 for the duration of this test: the stimulus needs two Ps to exist concurrently", prev)
		defer runtime.GOMAXPROCS(prev)
		runtime.GOMAXPROCS(2)
	}
	s := openEvidenceStoreWithoutBusyWait(t)
	pinned := pinnedChecker(t)
	base, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	decoyRef := putEvidenceObject(t, s, []byte(strings.Repeat("d", 256<<20)))
	hold := measureGetObject(t, s, decoyRef)
	if hold < minDecoyHold {
		t.Fatalf("instrument failure: decoy held the connection for only %v; want at least %v for this stimulus to exist", hold, minDecoyHold)
	}
	bound := hold / decoyHoldRatio
	settle := hold / 40
	watchdog := 20*hold + 60*time.Second
	readTimeout := 3*time.Second + 2*hold
	required := []string{"twicePlusOne"}
	t.Logf("decoy hold=%v -> derived ObjectWriteTimeout=%v settle=%v watchdog=%v", hold, bound, settle, watchdog)

	// No-decoy control first: the same producer mints when nothing pins the connection.
	controlSrc := putSource(t, s, base)
	// Generous write bound for the control: it proves only that the producer and the
	// real store mint end to end with no decoy, not anything about the derived bound.
	ctrlP, err := evidence.NewProducer(testKey, s, decoyWriter{s: s}, pinned, evidence.ProducerConfig{
		ObjectReadTimeout: readTimeout, ObjectWriteTimeout: controlBound, ExecTimeout: 2 * time.Minute, MaxOutputBytes: 1 << 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ref, err := ctrlP.GenerateProof(context.Background(), controlSrc, required); err != nil || ref.IsZero() {
		t.Fatalf("the no-decoy control must mint end to end with no decoy: ref=%s err=%v", ref, err)
	}

	d := decoyStimulus{t: t, s: s, decoy: decoyRef, settle: settle}
	for attempt := 1; attempt <= 5; attempt++ {
		// Each armed attempt uses a source no earlier attempt stored: a won-connection
		// mint leaves its envelope committed, and a repeat would hide the next result.
		src := append(append([]byte(nil), base...), strings.Repeat("\n", attempt)...)
		srcRef := putSource(t, s, src)
		rec := &attemptRecord{}
		p, err := evidence.NewProducer(testKey, s, decoyWriter{s: s, d: d, rec: rec}, pinned, evidence.ProducerConfig{
			ObjectReadTimeout: readTimeout, ObjectWriteTimeout: bound, ExecTimeout: 2 * time.Minute, MaxOutputBytes: 1 << 20,
		})
		if err != nil {
			t.Fatal(err)
		}
		ref, err := generateWithin(t, p, srcRef, required, watchdog)
		rec.join()
		switch {
		case rec.instrument != "":
			t.Logf("attempt %d: instrument retry: %s", attempt, rec.instrument)
		case err != nil:
			requireDeadlineRefusal(t, ref, err)
			if rec.wouldBe.IsZero() {
				t.Fatal("instrument failure: the armed write never reached the store")
			}
			if _, ok, gerr := s.GetObject(boundedTestContext(t), rec.wouldBe); gerr != nil || ok {
				t.Fatalf("a refused envelope write must leave no object behind: ok=%v err=%v", ok, gerr)
			}
			t.Logf("blocked write refused after store call %v (derived bound %v)", rec.callElapsed, bound)
			return
		case rec.callElapsed < 2*bound:
			t.Logf("attempt %d: the store call won the freed connection in %v; retrying", attempt, rec.callElapsed)
		default:
			t.Fatalf("blocked envelope write stored after exceeding ObjectWriteTimeout: store call=%v; want a deadline refusal", rec.callElapsed)
		}
	}
	t.Fatal("instrument failure: the store call won the freed connection in all 5 attempts")
}
