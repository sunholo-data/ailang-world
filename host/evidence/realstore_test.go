package evidence_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const realStoreReadTimeout = 2 * time.Millisecond

type waitBoundWrapper struct {
	reader *store.Store
	busy   time.Duration
}

func (w waitBoundWrapper) ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (store.ObjectMeta, []byte, error) {
	return w.reader.ReadObject(ctx, ref, maxBytes)
}

func (w waitBoundWrapper) BusyTimeout() time.Duration { return w.busy }

func openEvidenceStore(t *testing.T) (*store.Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "evidence.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("close real store: %v", err)
		}
	})
	return s, path
}

func openEvidenceStoreWithoutBusyWait(t *testing.T) *store.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "evidence.db")
	dsn := "file:" + filepath.ToSlash(path) + "?" + url.Values{"_pragma": {"busy_timeout(0)"}}.Encode()
	s, err := store.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("close real store: %v", err)
		}
	})
	return s
}

// openEvidenceStoreWithBusyWindow opens a store whose lock-retry window is set
// by the CALLER's DSN, deliberately away from the production default, so an
// accessor that had drifted into the hardcoded 2000 ms literal reports a value
// the driver is not applying.
func openEvidenceStoreWithBusyWindow(t *testing.T, window time.Duration) (*store.Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "evidence.db")
	dsn := "file:" + filepath.ToSlash(path) + "?" + url.Values{
		"_pragma": {fmt.Sprintf("busy_timeout(%d)", window.Milliseconds())},
	}.Encode()
	s, err := store.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("close real store: %v", err)
		}
	})
	return s, dsn
}

func putEvidenceObject(t *testing.T, s *store.Store, payload []byte) hashref.HashRef {
	t.Helper()
	ref := hashref.SumSHA256(payload)
	if err := s.PutObject(store.Object{
		Hash: ref, InterfaceHash: evidence.InterfaceHashV1,
		SemanticID: evidence.ProofSemanticID, Provenance: "real-store-test", Payload: payload,
	}); err != nil {
		t.Fatal(err)
	}
	return ref
}

func validatorConfig(timeout time.Duration) evidence.CompilerConfig {
	return evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: timeout}
}

func realValidator(t *testing.T, s *store.Store, timeout time.Duration) *evidence.Validator {
	t.Helper()
	v, err := evidence.NewValidator(testKey, s, validatorConfig(timeout), []string{"world/types.ail"})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func measureGetObject(t *testing.T, s *store.Store, ref hashref.HashRef) time.Duration {
	t.Helper()
	start := time.Now()
	got, ok, err := s.GetObject(context.Background(), ref)
	elapsed := time.Since(start)
	if err != nil || !ok || got.Hash != ref {
		t.Fatalf("decoy GetObject: ok=%v hash=%v err=%v", ok, got.Hash, err)
	}
	return elapsed
}

func startDecoyRead(s *store.Store, ref hashref.HashRef) <-chan time.Duration {
	done := make(chan time.Duration, 1)
	go func() {
		start := time.Now()
		_, _, _ = s.GetObject(context.Background(), ref)
		done <- time.Since(start)
	}()
	return done
}

func TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout(t *testing.T) {
	// The stimulus needs at least two Ps, and that is a measurement rather than
	// a precaution: this driver's SQLite is pure Go, so at GOMAXPROCS=1 the decoy
	// monopolizes the single P for its whole ~53 ms read and the blocked reader
	// is never scheduled to join the pool queue at all -- it runs only after the
	// connection is free, wins it, and exhausts the retry budget. Measured on
	// unmutated code: 2/10 retry-exhaustion reds at GOMAXPROCS=1, 0/10 at the
	// default. Declare the instrument's precondition and enforce it instead of
	// inheriting whatever the harness happens to set.
	if prev := runtime.GOMAXPROCS(0); prev < 2 {
		defer runtime.GOMAXPROCS(prev)
		runtime.GOMAXPROCS(2)
	}
	s := openEvidenceStoreWithoutBusyWait(t)
	goodPayload := envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "good")
	goodRef := putEvidenceObject(t, s, goodPayload)
	v := realValidator(t, s, 3*time.Second)
	requireProvenControl(t, v, goodRef) // required no-decoy control runs first

	decoyRef := putEvidenceObject(t, s, []byte(strings.Repeat("d", 256<<20)))
	hold := measureGetObject(t, s, decoyRef)
	t.Logf("decoy hold=%v ObjectReadTimeout=%v ratio=%.1fx", hold, realStoreReadTimeout, float64(hold)/float64(realStoreReadTimeout))
	if hold <= 20*realStoreReadTimeout {
		t.Fatalf("instrument failure: decoy held connection for %v; want > 20x ObjectReadTimeout (%v)", hold, 20*realStoreReadTimeout)
	}
	v = realValidator(t, s, realStoreReadTimeout)

	for attempt := 1; attempt <= 5; attempt++ {
		decoyDone := startDecoyRead(s, decoyRef)
		time.Sleep(realStoreReadTimeout)
		resultDone := make(chan evidence.ValidationResult, 1)
		start := time.Now()
		go func() { resultDone <- v.ValidateProof(context.Background(), goodRef, testSubject) }()
		select {
		case result := <-resultDone:
			elapsed := time.Since(start)
			_, sealed := result.Validated()
			if sealed && elapsed < 2*realStoreReadTimeout {
				<-decoyDone
				continue // validation won the connection; retry honestly
			}
			if sealed {
				t.Fatalf("blocked read sealed after exceeding ObjectReadTimeout: elapsed=%v; want operational timeout error and no seal", elapsed)
			}
			if result.Err() == nil {
				t.Fatalf("blocked read returned no operational error: unsupported=%v", unsupportedOf(result))
			}
			// AC16 classifies SEALS, not refusals. A refusal is the correct
			// outcome however late the scheduler delivers it, and the 20x
			// watchdog above is what bounds that lateness; a 2x ceiling on the
			// REFUSAL path measures scheduler latency instead of the deadline.
			// Measured on unmutated, sha256-identical code: GOMAXPROCS=1 reds
			// 10/10 with "returned after 10-33ms", default GOMAXPROCS reds 0/10.
			t.Logf("blocked read refused after %v (bound %v, watchdog %v)", elapsed, realStoreReadTimeout, 20*realStoreReadTimeout)
			<-decoyDone
			return
		case <-time.After(20 * realStoreReadTimeout):
			t.Fatal("blocked read exceeded test-side 20x watchdog")
		}
	}
	t.Fatal("instrument failure: validator won the freed connection in all 5 attempts")
}

func TestOversizeProofReportIsRefused(t *testing.T) {
	s, _ := openEvidenceStore(t)
	report := reportFor(testSubject, testCompiler)
	report.Verified = make([]string, 256)
	for i := range report.Verified {
		report.Verified[i] = fmt.Sprintf("%03d-%s", i, strings.Repeat("x", 1015))
	}
	raw, err := evidence.EncodeProofReportV1(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) >= evidence.MaxBytes || len(raw) < 240<<10 {
		t.Fatalf("report size = %d; want approximately 250 KiB and below %d", len(raw), evidence.MaxBytes)
	}
	payload := envelopeFor(t, testKey, report, "good")
	if len(payload) <= evidence.MaxBytes {
		t.Fatalf("envelope size = %d; want above real-store read bound %d", len(payload), evidence.MaxBytes)
	}
	ref := putEvidenceObject(t, s, payload)
	v, err := evidence.NewValidator(testKey, s, validatorConfig(3*time.Second), []string{report.Verified[0]})
	if err != nil {
		t.Fatal(err)
	}
	r := v.ValidateProof(context.Background(), ref, testSubject)
	if _, sealed := r.Validated(); sealed {
		t.Fatal("oversize envelope sealed; want oversize")
	}
	if r.Err() != nil {
		t.Fatalf("got operational store-read error; want oversize: %v", r.Err())
	}
	requireReason(t, r, evidence.UnsupportedOversize, "real-store pre-materialization guard and typed mapping")
}

func TestConstructorPinsBusyTimeoutBelowObjectReadTimeout(t *testing.T) {
	s, _ := openEvidenceStore(t)
	good := putEvidenceObject(t, s, envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "good"))
	ordered := realValidator(t, s, 3*time.Second)
	requireProvenControl(t, ordered, good)

	// The round-11b cross-check, spelled so it can FAIL. Reading a PRAGMA back
	// out of a handle whose DSN this test itself set to the production default
	// compares the accessor against a literal the test supplied -- it survives an
	// accessor hardcoded to 2 s, which is exactly the stale-literal drift the
	// amendment exists to catch (measured: that mutant leaves this whole test
	// green). So cross-check against a window the CALLER chose and the driver
	// therefore actually applies, on a store whose connection is free.
	const callerWindow = 1500 * time.Millisecond
	custom, customDSN := openEvidenceStoreWithBusyWindow(t, callerWindow)
	if custom.BusyTimeout() == s.BusyTimeout() {
		t.Fatalf("instrument failure: caller window %v equals the default %v, so this arm could not see a stale literal",
			custom.BusyTimeout(), s.BusyTimeout())
	}
	db, err := sql.Open("sqlite", customDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var pragmaMS int64
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&pragmaMS); err != nil {
		t.Fatal(err)
	}
	if got, want := custom.BusyTimeout(), time.Duration(pragmaMS)*time.Millisecond; got != want {
		t.Fatalf("BusyTimeout() = %v; the window the driver applies for this DSN is %v", got, want)
	}

	_, err = evidence.NewValidator(testKey, s, validatorConfig(time.Second), []string{"world/types.ail"})
	if !errors.Is(err, evidence.ErrUnorderedTimeouts) || errors.Is(err, evidence.ErrInvalidValidatorConfig) {
		t.Fatalf("positive unordered timeout: got %v; want ErrUnorderedTimeouts and not ErrInvalidValidatorConfig", err)
	}

	decoyRef := putEvidenceObject(t, s, []byte(strings.Repeat("n", 256<<20)))
	hold := measureGetObject(t, s, decoyRef)
	t.Logf("nonblocking decoy hold=%v", hold)
	if hold <= 20*realStoreReadTimeout {
		t.Fatalf("instrument failure: nonblocking decoy hold %v; want > %v", hold, 20*realStoreReadTimeout)
	}
	done := startDecoyRead(s, decoyRef)
	time.Sleep(realStoreReadTimeout)
	start := time.Now()
	if got := s.BusyTimeout(); got != 2*time.Second {
		t.Fatalf("occupied BusyTimeout() = %v; want 2s", got)
	}
	_, err = evidence.NewValidator(testKey, s, validatorConfig(time.Second), []string{"world/types.ail"})
	elapsed := time.Since(start)
	t.Logf("occupied BusyTimeout/NewValidator elapsed=%v (decoy hold=%v)", elapsed, hold)
	if !errors.Is(err, evidence.ErrUnorderedTimeouts) {
		t.Fatalf("occupied NewValidator: got %v; want ErrUnorderedTimeouts", err)
	}
	if elapsed >= hold/4 {
		t.Fatalf("cached accessor/constructor blocked for %v; want far below decoy hold %v", elapsed, hold)
	}
	<-done
}

func TestReaderWaitBoundsCannotBeLostThroughWrapper(t *testing.T) {
	s, _ := openEvidenceStore(t)
	cfg := validatorConfig(time.Second)
	t.Run("unknown", func(t *testing.T) {
		unknown := waitBoundWrapper{reader: s, busy: -time.Nanosecond}
		if _, err := evidence.NewValidator(testKey, unknown, cfg, []string{"world/types.ail"}); !errors.Is(err, evidence.ErrInvalidValidatorConfig) {
			t.Fatalf("NewValidator accepted wrapper with unknown wait bound: %v; want ErrInvalidValidatorConfig", err)
		}
	})
	t.Run("forwarding-real-store", func(t *testing.T) {
		forwarding := waitBoundWrapper{reader: s, busy: s.BusyTimeout()}
		if _, err := evidence.NewValidator(testKey, forwarding, cfg, []string{"world/types.ail"}); !errors.Is(err, evidence.ErrUnorderedTimeouts) {
			t.Fatalf("forwarding wrapper lost real-store wait bound: got %v; want ErrUnorderedTimeouts", err)
		}
	})
	// The arm above passes for ANY wrapper reporting >= 1 s, including one that
	// forwards nothing and returns the 2 s default as a literal -- measured: it
	// survives its own precondition being removed. This arm cannot: 1800 ms is
	// ORDERED against a caller-chosen 1500 ms window and UNORDERED against the
	// 2000 ms default, so a wrapper that stopped forwarding refuses here.
	t.Run("forwarding-preserves-a-non-default-window", func(t *testing.T) {
		const callerWindow = 1500 * time.Millisecond
		custom, _ := openEvidenceStoreWithBusyWindow(t, callerWindow)
		if custom.BusyTimeout() != callerWindow {
			t.Fatalf("instrument failure: store reports %v; want the caller-chosen %v", custom.BusyTimeout(), callerWindow)
		}
		forwarding := waitBoundWrapper{reader: custom, busy: custom.BusyTimeout()}
		if _, err := evidence.NewValidator(testKey, forwarding, validatorConfig(1800*time.Millisecond), []string{"world/types.ail"}); err != nil {
			t.Fatalf("forwarding wrapper over a %v window refused an ORDERED 1800ms timeout: %v", callerWindow, err)
		}
	})
}
