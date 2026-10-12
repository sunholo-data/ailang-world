package evidence

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/proctest"
	"github.com/sunholo-data/ailang-world/host/store"
)

var internalTestKey = [32]byte{9, 8, 7}

// nopSeams is an inert ObjectReader and ObjectWriter for stage tests that never
// touch an object: it reports no busy window, so no ordering pin applies.
type nopSeams struct{}

func (nopSeams) ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error) {
	return store.ObjectMeta{}, nil, nil
}
func (nopSeams) PutObject(context.Context, store.Object) error { return nil }
func (nopSeams) BusyTimeout() time.Duration                    { return 0 }

// scriptedProducer writes an executable #!/bin/sh checker into its own temp dir,
// pins it by its own bytes, and constructs a producer through NewProducer. The
// script answers --version with a fixed line and otherwise runs body.
func scriptedProducer(t *testing.T, body string, execTimeout time.Duration, maxOut int64) (*Producer, string) {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo \"AILANG scripted\"; exit 0; fi\n" + body + "\n"
	path := filepath.Join(dir, "ailang")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := NewProducer(internalTestKey, nopSeams{}, nopSeams{}, CheckerConfig{
		ExecPath:        path,
		Compiler:        hashref.SumSHA256([]byte(script)),
		CompilerVersion: "AILANG scripted",
	}, ProducerConfig{
		ObjectReadTimeout:  time.Second,
		ObjectWriteTimeout: time.Second,
		ExecTimeout:        execTimeout,
		MaxOutputBytes:     maxOut,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p, dir
}

func runEnv(dir string) []string {
	return []string{"PATH=" + os.Getenv("PATH"), "PIDDIR=" + dir}
}

// killRecorder observes every process-group kill, forwarding to the real kill.
type killRecorder struct {
	mu    sync.Mutex
	pgids []int
	times []time.Time
}

func (k *killRecorder) record(pgid int, at time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.pgids = append(k.pgids, pgid)
	k.times = append(k.times, at)
}

func (k *killRecorder) snapshot() ([]int, []time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	return append([]int(nil), k.pgids...), append([]time.Time(nil), k.times...)
}

func swapKillGroup(t *testing.T) *killRecorder {
	t.Helper()
	rec := &killRecorder{}
	orig := killGroup
	killGroup = func(pgid int) error {
		rec.record(pgid, time.Now())
		return orig(pgid)
	}
	t.Cleanup(func() { killGroup = orig })
	return rec
}

func TestProducerRefusesStderrOverflowIndependently(t *testing.T) {
	t.Run("overflow", func(t *testing.T) {
		rec := swapKillGroup(t)
		p, dir := scriptedProducer(t, `printf '{}'; head -c 64 /dev/zero | tr '\0' e >&2; exec sleep 30`, 10*time.Second, 16)
		start := time.Now()
		out, err := p.runBounded(context.Background(), dir, runEnv(dir), 16)
		elapsed := time.Since(start)
		if !errors.Is(err, ErrOutputLimit) {
			t.Fatalf("64 stderr bytes over a 16-byte cap: want ErrOutputLimit; got %v", err)
		}
		if errors.Is(err, ErrExecTimeout) {
			t.Fatalf("an overflow must not be reported as a timeout: %v", err)
		}
		if !strings.Contains(err.Error(), "stderr") {
			t.Fatalf("detail must name the stderr stream: %v", err)
		}
		if out != nil {
			t.Fatalf("overflow must return nil bytes; got %q", out)
		}
		if elapsed >= 5*time.Second {
			t.Fatalf("the overflow kill must end the child; it ran %s against a 10s timeout", elapsed)
		}
		if pgids, _ := rec.snapshot(); len(pgids) != 1 {
			t.Fatalf("overflow must kill the group exactly once; killGroup saw %v", pgids)
		}
	})

	t.Run("independence-control", func(t *testing.T) {
		p, dir := scriptedProducer(t, `printf '{}'; head -c 64 /dev/zero | tr '\0' e >&2`, 10*time.Second, 1024)
		out, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
		if err != nil || string(out) != "{}" {
			t.Fatalf("64 stderr bytes under a 1024-byte cap must pass: out=%q err=%v", out, err)
		}
	})
}

func TestProducerExecTimeoutKillsTheWholeProcessGroup(t *testing.T) {
	// execTimeout is far above the shell's startup and pid writes, so the stimulus
	// completes inside the bound even when the full parallel suite is loaded.
	const execTimeout = 5 * time.Second
	rec := swapKillGroup(t)
	p, dir := scriptedProducer(t, `sleep 30 & echo $! > "$PIDDIR/grandchild.pid"; echo $$ > "$PIDDIR/child.pid"; wait`, execTimeout, 1024)
	start := time.Now()
	done := make(chan error, 1)
	go func() {
		_, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
		done <- err
	}()
	var err error
	select {
	case err = <-done:
	case <-time.After(30 * time.Second):
		t.Fatalf("runBounded did not return within 30s; want ErrExecTimeout after ExecTimeout=%s", execTimeout)
	}
	grandchild := proctest.ReadPid(t, filepath.Join(dir, "grandchild.pid"))
	child := proctest.ReadPid(t, filepath.Join(dir, "child.pid"))
	proctest.ReapOnCleanup(t, grandchild)
	proctest.ReapOnCleanup(t, child)

	if !errors.Is(err, ErrExecTimeout) {
		t.Fatalf("want ErrExecTimeout; got %v", err)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("the group kill must complete the cleanup; got %v", err)
	}
	pgids, times := rec.snapshot()
	if len(pgids) != 1 || pgids[0] != child {
		t.Fatalf("killGroup must see exactly the child's group %d; saw %v", child, pgids)
	}
	if lag := times[0].Sub(start); lag < execTimeout {
		t.Fatalf("first kill at %s after the call; want at least the %s ExecTimeout", lag, execTimeout)
	}
	if !proctest.DeadWithin(t, grandchild, 2*time.Second) {
		t.Fatalf("grandchild %d survived the process-group kill", grandchild)
	}
}

func TestProducerAdmissionRefusesBacklog(t *testing.T) {
	p, dir := scriptedProducer(t, `touch "$PIDDIR/ran"`, 10*time.Second, 1024)
	before := procbound.Outstanding()
	var releases []func()
	for {
		release, err := procbound.Admit()
		if err != nil {
			if !errors.Is(err, procbound.ErrCleanupBacklog) {
				t.Fatalf("filling the procbound slots: want ErrCleanupBacklog at the cap; got %v", err)
			}
			break
		}
		releases = append(releases, release)
	}
	defer func() {
		for _, release := range releases {
			release()
		}
	}()
	_, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
	if !errors.Is(err, procbound.ErrCleanupBacklog) {
		t.Fatalf("a full backlog must refuse the spawn: want ErrCleanupBacklog; got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "ran")); statErr == nil {
		t.Fatalf("the refused checker ran anyway")
	}
	for _, release := range releases {
		release()
	}
	releases = nil
	if got := procbound.Outstanding(); got != before {
		t.Fatalf("after releasing, outstanding = %d; want the pre-test %d", got, before)
	}
}

// pinnedChecker pins the real AILANG_BIN by its bytes and its version line. A
// binary-dependent arm fails loudly when the variable is unset and never skips.
func pinnedChecker(t *testing.T) CheckerConfig {
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
	return CheckerConfig{ExecPath: bin, Compiler: hashref.SumSHA256(data), CompilerVersion: "AILANG v0.41.0"}
}

func pinnedProducer(t *testing.T, execTimeout time.Duration, maxOut int64) *Producer {
	t.Helper()
	p, err := NewProducer(internalTestKey, nopSeams{}, nopSeams{}, pinnedChecker(t), ProducerConfig{
		ObjectReadTimeout:  time.Second,
		ObjectWriteTimeout: time.Second,
		ExecTimeout:        execTimeout,
		MaxOutputBytes:     maxOut,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestProducerRefusesStdoutOverflowIndependently(t *testing.T) {
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t.Run("real-binary-overflow", func(t *testing.T) {
		t.Setenv("TMPDIR", t.TempDir())
		out, err := pinnedProducer(t, 2*time.Minute, 16).runCheck(ctx, source)
		if !errors.Is(err, ErrOutputLimit) || out != nil {
			t.Fatalf("real ai-check output over a 16-byte cap: want ErrOutputLimit and nil bytes; got %d bytes, err=%v", len(out), err)
		}
	})
	t.Run("real-binary-control", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("TMPDIR", tmp)
		out, err := pinnedProducer(t, 2*time.Minute, 1<<20).runCheck(ctx, source)
		if err != nil {
			t.Fatalf("under the 1 MiB cap the real run must return its JSON: %v", err)
		}
		var doc struct {
			Verify struct {
				Results []struct {
					Function string `json:"function"`
					Status   string `json:"status"`
				} `json:"results"`
			} `json:"verify"`
		}
		if err := json.Unmarshal(out, &doc); err != nil {
			t.Fatalf("control output is not ai-check JSON: %v", err)
		}
		verified := false
		for _, r := range doc.Verify.Results {
			verified = verified || (r.Function == "twicePlusOne" && r.Status == "verified")
		}
		if !verified {
			t.Fatalf("control output lacks twicePlusOne verified: %s", out)
		}
		entries, err := os.ReadDir(tmp)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("the producer left %d entries in its TMPDIR; scratch must be removed", len(entries))
		}
	})
	t.Run("scripted-stdout-overflow", func(t *testing.T) {
		p, _ := scriptedProducer(t, `head -c 64 /dev/zero | tr '\0' o`, 10*time.Second, 16)
		out, err := p.runCheck(ctx, source)
		if !errors.Is(err, ErrOutputLimit) || out != nil {
			t.Fatalf("64 stdout bytes over a 16-byte cap: want ErrOutputLimit and nil bytes; got %d bytes, err=%v", len(out), err)
		}
		if !strings.Contains(err.Error(), "stdout") {
			t.Fatalf("detail must name the stdout stream: %v", err)
		}
	})
}

// producerFor builds a producer around an explicit checker with default bounds.
func producerFor(t *testing.T, checker CheckerConfig) *Producer {
	t.Helper()
	p, err := NewProducer(internalTestKey, nopSeams{}, nopSeams{}, checker, ProducerConfig{
		ObjectReadTimeout:  time.Second,
		ObjectWriteTimeout: time.Second,
		ExecTimeout:        time.Minute,
		MaxOutputBytes:     1 << 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// writeScript writes an executable #!/bin/sh file into a fresh temp dir.
func writeScript(t *testing.T, script string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ailang")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestProducerRefutesTamperedExecutableBytes(t *testing.T) {
	ctx := context.Background()
	t.Run("real-binary-wrong-pin", func(t *testing.T) {
		checker := pinnedChecker(t)
		real := checker.Compiler
		checker.Compiler = hashref.SumSHA256([]byte("not the pinned binary"))
		err := producerFor(t, checker).verifyTool(ctx)
		if !errors.Is(err, ErrUntrustedTool) {
			t.Fatalf("a binary whose bytes miss the pin must be refused: want ErrUntrustedTool; got %v", err)
		}
		if !strings.Contains(err.Error(), "hash to") || !strings.Contains(err.Error(), real.String()) {
			t.Fatalf("detail must name the hash to the real bytes (%s): %v", real, err)
		}
	})
	t.Run("scripted-executed-before-check", func(t *testing.T) {
		// The touch sits BEFORE the version branch, so any execution leaves the marker.
		body := "#!/bin/sh\ntouch \"$(dirname \"$0\")/executed\"\nif [ \"$1\" = \"--version\" ]; then echo \"AILANG scripted\"; exit 0; fi\n"
		path := writeScript(t, body)
		wrong := CheckerConfig{ExecPath: path, Compiler: hashref.SumSHA256([]byte("wrong bytes")), CompilerVersion: "AILANG scripted"}
		if err := producerFor(t, wrong).verifyTool(ctx); !errors.Is(err, ErrUntrustedTool) {
			t.Fatalf("a script pinned to the wrong bytes must be refused: want ErrUntrustedTool; got %v", err)
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(path), "executed")); err == nil {
			t.Fatal("the tampered executable ran before its bytes were checked")
		}
		right := CheckerConfig{ExecPath: path, Compiler: hashref.SumSHA256([]byte(body)), CompilerVersion: "AILANG scripted"}
		if err := producerFor(t, right).verifyTool(ctx); err != nil {
			t.Fatalf("control: the pinned bytes and version must pass: %v", err)
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(path), "executed")); err != nil {
			t.Fatalf("control did not execute the probe, so the marker proves nothing: %v", err)
		}
	})
}

func TestProducerRefutesWrongCompilerVersion(t *testing.T) {
	ctx := context.Background()
	checker := pinnedChecker(t)
	wrong := checker
	wrong.CompilerVersion = "AILANG v9.9.9"
	err := producerFor(t, wrong).verifyTool(ctx)
	if !errors.Is(err, ErrUntrustedTool) || !strings.Contains(err.Error(), "AILANG v0.41.0") || !strings.Contains(err.Error(), "AILANG v9.9.9") {
		t.Fatalf("a wrong version pin must name both lines: want ErrUntrustedTool with v0.41.0 and v9.9.9; got %v", err)
	}
	if err := producerFor(t, checker).verifyTool(ctx); err != nil {
		t.Fatalf("control: the correct version must pass: %v", err)
	}
}

func TestProducerRefusesMissingExecutable(t *testing.T) {
	checker := pinnedChecker(t)
	checker.ExecPath = filepath.Join(t.TempDir(), "absent")
	err := producerFor(t, checker).verifyTool(context.Background())
	if !errors.Is(err, ErrUntrustedTool) || !strings.Contains(err.Error(), "read executable") {
		t.Fatalf("a missing executable must be refused by name: want ErrUntrustedTool with read executable; got %v", err)
	}
}

// verdictClass names an interpretCheck outcome for comparison: the sentinel it
// wraps, "ok" for a verified set, or the raw text of anything unexpected.
func verdictClass(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, ErrUnverifiedRun):
		return "unverified"
	case errors.Is(err, ErrVacuousCheck):
		return "vacuous"
	default:
		return "other: " + err.Error()
	}
}

// fixtureBytes reads one committed ai-check fixture. Fixtures are read-only
// inputs: no test writes them (regen.sh is the only writer).
func fixtureBytes(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "aicheck", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// flippedVerifiedPass derives an ai-check output from the committed verifiedpass
// fixture by changing exactly one field. The derived arms pin branches that no
// real pinned run can reach, so each is labelled a derivation, never tool output.
func flippedVerifiedPass(t *testing.T, mutate func(doc map[string]any)) []byte {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(fixtureBytes(t, "verifiedpass.json"), &doc); err != nil {
		t.Fatal(err)
	}
	mutate(doc)
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestProducerInterpretationMatchesEveryFixtureShape(t *testing.T) {
	type arm struct {
		name     string
		stdout   []byte
		want     error // nil: the run verifies and verified is the expected set
		detail   string
		verified []string
	}
	arms := []arm{
		{name: "verifiedpass", stdout: fixtureBytes(t, "verifiedpass.json"), verified: []string{"twicePlusOne"}},
		{name: "refuted", stdout: fixtureBytes(t, "refuted.json"), want: ErrUnverifiedRun, detail: "verify.counterexample=1"},
		{name: "encerror", stdout: fixtureBytes(t, "encerror.json"), want: ErrUnverifiedRun, detail: "verify.errors=1"},
		{name: "mixed", stdout: fixtureBytes(t, "mixed.json"), want: ErrUnverifiedRun, detail: "verify.errors=1"},
		{name: "skippedshape", stdout: fixtureBytes(t, "skippedshape.json"), want: ErrUnverifiedRun, detail: "verify.skipped=1"},
		{name: "loaderror", stdout: fixtureBytes(t, "loaderror.json"), want: ErrUnverifiedRun, detail: "check.passed is not true"},
		{name: "nocontracts", stdout: fixtureBytes(t, "nocontracts.json"), want: ErrVacuousCheck, detail: "0 results records from non-empty output"},
		{name: "literal-empty", stdout: []byte(""), want: ErrVacuousCheck, detail: "empty ai-check stdout"},
		{name: "literal-unparseable", stdout: []byte("{"), want: ErrUnverifiedRun, detail: "unparseable"},
		{name: "derived-check-passed-false", want: ErrUnverifiedRun, detail: "check.passed",
			stdout: flippedVerifiedPass(t, func(d map[string]any) { d["check"].(map[string]any)["passed"] = false })},
		{name: "derived-verify-unavailable", want: ErrUnverifiedRun, detail: "verify.available",
			stdout: flippedVerifiedPass(t, func(d map[string]any) { d["verify"].(map[string]any)["available"] = false })},
		{name: "derived-status-unknown", want: ErrVacuousCheck, detail: "0 verified among 1 results records",
			stdout: flippedVerifiedPass(t, func(d map[string]any) {
				d["verify"].(map[string]any)["results"].([]any)[0].(map[string]any)["status"] = "unknown"
			})},
	}
	for _, a := range arms {
		a := a
		t.Run(a.name, func(t *testing.T) {
			got, err := interpretCheck(a.stdout)
			if a.want == nil {
				if err != nil || strings.Join(got, ",") != strings.Join(a.verified, ",") {
					t.Fatalf("want verified %v with nil error; got %v, err=%v", a.verified, got, err)
				}
				return
			}
			if !errors.Is(err, a.want) {
				t.Fatalf("want %v; got %v (verified %v)", a.want, err, got)
			}
			if !strings.Contains(err.Error(), a.detail) {
				t.Fatalf("detail %q missing from %v", a.detail, err)
			}
			if len(got) != 0 {
				t.Fatalf("a refused run must return no verified set; got %v", got)
			}
		})
	}
}

// TestFixturesRevalidateAgainstFreshBinaryOutput re-runs every committed module
// through the pinned binary. A pin move that changes the JSON shape reds here:
// the fixtures are never trusted as frozen truth.
func TestFixturesRevalidateAgainstFreshBinaryOutput(t *testing.T) {
	p := pinnedProducer(t, 2*time.Minute, 1<<20)
	for _, name := range []string{"verifiedpass", "refuted", "encerror", "mixed", "skippedshape", "loaderror", "nocontracts"} {
		name := name
		t.Run(name, func(t *testing.T) {
			fresh, err := p.runCheck(context.Background(), fixtureBytes(t, name+".ail"))
			if err != nil {
				t.Fatalf("fresh ai-check run: %v", err)
			}
			gotFresh, errFresh := interpretCheck(fresh)
			gotCommitted, errCommitted := interpretCheck(fixtureBytes(t, name+".json"))
			if verdictClass(errFresh) != verdictClass(errCommitted) || strings.Join(gotFresh, ",") != strings.Join(gotCommitted, ",") {
				t.Fatalf("fresh verdict %s %v differs from committed verdict %s %v: the pinned binary's output moved",
					verdictClass(errFresh), gotFresh, verdictClass(errCommitted), gotCommitted)
			}
		})
	}
}

// mintTestChecker pins the checker identity for mint tests. mint reads only the
// pin and the version line, never the binary, so the path is never executed.
func mintTestChecker() CheckerConfig {
	return CheckerConfig{
		ExecPath:        "/nonexistent/ailang",
		Compiler:        hashref.SumSHA256([]byte("pinned checker")),
		CompilerVersion: "AILANG v0.41.0",
	}
}

var mintTestSubject = hashref.SumSHA256([]byte("subject source bytes"))

// identitySet returns count sorted, unique identities of exactly size bytes each,
// so a test can place the envelope on either side of the validator's read cap
// through the real codec.
func identitySet(count, size int) []string {
	ids := make([]string, count)
	for i := range ids {
		ids[i] = fmt.Sprintf("id%03d", i) + strings.Repeat("x", size-5)
	}
	return ids
}

func TestProducerRequiresCallerNamedIdentities(t *testing.T) {
	checker := mintTestChecker()
	p := producerFor(t, checker)
	verified, err := interpretCheck(fixtureBytes(t, "verifiedpass.json"))
	if err != nil {
		t.Fatalf("the committed verifiedpass fixture must interpret: %v", err)
	}
	t.Run("control", func(t *testing.T) {
		if err := requireIdentities(verified, []string{"twicePlusOne"}); err != nil {
			t.Fatalf("a required identity that is verified must pass: %v", err)
		}
		env, ref, err := p.mint(mintTestSubject, verified)
		if err != nil {
			t.Fatalf("mint: %v", err)
		}
		decoded, err := DecodeAuthenticatedEnvelope(env)
		if err != nil {
			t.Fatalf("minted envelope must decode: %v", err)
		}
		mac := hmac.New(sha256.New, internalTestKey[:])
		mac.Write(decoded.Report)
		if !hmac.Equal(mac.Sum(nil), decoded.MAC) {
			t.Fatal("the MAC does not authenticate the report under the producer key")
		}
		report, err := DecodeProofReportV1(decoded.Report)
		if err != nil {
			t.Fatalf("minted report must decode: %v", err)
		}
		if strings.Join(report.Verified, ",") != "twicePlusOne" || report.Subject != mintTestSubject ||
			report.Compiler != checker.Compiler || report.CompilerVersion != checker.CompilerVersion {
			t.Fatalf("report fields do not match the configured subject and pin: %+v", report)
		}
		if report.Errors != 0 || report.Counterexamples != 0 || !report.CheckPassed || !report.ProofSucceeded {
			t.Fatalf("a minted report must be clean: %+v", report)
		}
		if ref != hashref.SumSHA256(env) {
			t.Fatalf("the returned ref %s must address the envelope bytes", ref)
		}
	})
	t.Run("absent", func(t *testing.T) {
		err := requireIdentities(verified, []string{"absentIdentity", "twicePlusOne"})
		if !errors.Is(err, ErrRequiredIdentityMissing) || !strings.Contains(err.Error(), "absentIdentity") {
			t.Fatalf("an unverified required identity must be refused by name: %v", err)
		}
	})
	t.Run("empty", func(t *testing.T) {
		for _, required := range [][]string{nil, {}} {
			err := requireIdentities(verified, required)
			if !errors.Is(err, ErrRequiredIdentityMissing) || !strings.Contains(err.Error(), "no required identities named") {
				t.Fatalf("an empty requirement must be refused: want ErrRequiredIdentityMissing; got %v", err)
			}
		}
	})
}

func TestProducerRefusesEnvelopeBeyondValidatorReadCap(t *testing.T) {
	p := producerFor(t, mintTestChecker())
	t.Run("over", func(t *testing.T) {
		env, ref, err := p.mint(mintTestSubject, identitySet(200, 1000))
		if !errors.Is(err, ErrEnvelopeTooLarge) {
			t.Fatalf("an envelope past the validator read cap must be refused: want ErrEnvelopeTooLarge; got %v", err)
		}
		if env != nil || !ref.IsZero() {
			t.Fatalf("a refused mint must return no bytes and a zero ref; got %d bytes, ref %s", len(env), ref)
		}
	})
	t.Run("near-bound-control", func(t *testing.T) {
		env, _, err := p.mint(mintTestSubject, identitySet(190, 1024))
		if err != nil {
			t.Fatalf("an envelope under the cap must mint: %v", err)
		}
		if len(env) > MaxBytes || len(env) <= MaxBytes-4096 {
			t.Fatalf("the control must straddle the cap near its bound: envelope %d bytes, cap %d", len(env), MaxBytes)
		}
	})
}
