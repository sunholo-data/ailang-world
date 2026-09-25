package pkgproj

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func pinnedAilang(t *testing.T) string {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: never skip")
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil || !strings.HasPrefix(string(out), "AILANG v0.41.0") {
		t.Fatalf("pinned binary unusable: %v %q", err, out)
	}
	return bin
}

var worldCore = Manifest{
	Package: Package{Name: "world/core", Edition: "1", AILANG: ">=0.30.0", Version: "0.1.0"},
	Exports: Exports{Modules: []string{"world/types", "world/contracts", "world/transitions", "world/logepoch"}},
	Effects: Effects{Max: []string{}},
}

func copyPackage(t *testing.T) string {
	t.Helper()
	dst := filepath.Join(t.TempDir(), "world-core")
	if out, err := exec.Command("cp", "-R", "../../packages/world-core", dst).CombinedOutput(); err != nil {
		t.Fatalf("cp: %v %s", err, out)
	}
	return dst
}

func editTypes(t *testing.T, dir string, f func(string) string) {
	t.Helper()
	p := filepath.Join(dir, "world", "types.ail")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	after := f(string(b))
	if after == string(b) {
		t.Fatal("instrument failure: edit changed nothing")
	}
	if err := os.WriteFile(p, []byte(after), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Re-break test: the row's exact stimulus.
func TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor(t *testing.T) {
	bin := pinnedAilang(t)
	base, err := QueryInterface(context.Background(), copyPackage(t), worldCore, bin)
	if err != nil {
		t.Fatal(err)
	}
	dir := copyPackage(t)
	editTypes(t, dir, func(s string) string {
		return strings.Replace(s, "  | ProofReceipt(HashRef)\n", "  | ProofReceipt(HashRef)\n  | ProbeArm(HashRef)\n", 1)
	})
	got, err := QueryInterface(context.Background(), dir, worldCore, bin)
	if err != nil {
		t.Fatal(err)
	}
	if got.V2 == base.V2 {
		t.Fatalf("interface v2 did not move on a new constructor: %s", got.V2)
	}
	if got.V1 != base.V1 {
		t.Fatalf("manifest-coverage v1 moved (%s -> %s); it is documented as manifest-only", base.V1, got.V1)
	}
	if got.Signatures != base.Signatures+1 {
		t.Fatalf("signatures %d -> %d, want +1", base.Signatures, got.Signatures)
	}
}

// Negative control: a comment moves neither.
func TestInterfaceV2IgnoresACommentOnlyEdit(t *testing.T) {
	bin := pinnedAilang(t)
	base, err := QueryInterface(context.Background(), copyPackage(t), worldCore, bin)
	if err != nil {
		t.Fatal(err)
	}
	dir := copyPackage(t)
	editTypes(t, dir, func(s string) string { return s + "-- comment-only control\n" })
	got, err := QueryInterface(context.Background(), dir, worldCore, bin)
	if err != nil {
		t.Fatal(err)
	}
	if got != base {
		t.Fatalf("comment-only edit moved the identity: %+v -> %+v", base, got)
	}
}

func TestQueryInterfaceRefusesAV1Disagreement(t *testing.T) {
	bin := pinnedAilang(t)
	wrong := worldCore
	wrong.Package.Edition = "2"
	_, err := QueryInterface(context.Background(), copyPackage(t), wrong, bin)
	if err == nil || !strings.Contains(err.Error(), "disagrees with pkgproj.InterfaceHash") {
		t.Fatalf("want v1 cross-check refusal, got %v", err)
	}
}

func TestParseQualityInterfaceFixtures(t *testing.T) {
	read := func(n string) []byte {
		b, err := os.ReadFile(filepath.Join("testdata", n))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	p, err := ParseQualityInterface(read("quality_world_core_pristine.json"))
	if err != nil || p.V2 != "sha256:ifacev2:b25fe03155db0c7bf595cf730295b945d6ac64ec415a1998fae8a693d621e8d8" || p.Signatures != 53 {
		t.Fatalf("pristine: %+v %v", p, err)
	}
	c, err := ParseQualityInterface(read("quality_world_core_probe_ctor.json"))
	if err != nil || c.V2 == p.V2 || c.V1 != p.V1 {
		t.Fatalf("ctor: %+v %v", c, err)
	}
	if _, err := ParseQualityInterface(read("quality_world_core_v2_unbuilt.json")); err == nil ||
		!strings.Contains(err.Error(), "interface identity v2 absent") {
		t.Fatalf("v2-unbuilt fixture: want the v2-absent refusal, got %v", err)
	}
}

func TestQueryInterfaceRefusesAnInfraExitEvenWithValidJSON(t *testing.T) {
	fixture, err := filepath.Abs(filepath.Join("testdata", "quality_world_core_pristine.json"))
	if err != nil {
		t.Fatal(err)
	}
	fake := filepath.Join(t.TempDir(), "ailang")
	script := "#!/bin/sh\ncat '" + fixture + "'\nexit 1\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := QueryInterface(context.Background(), t.TempDir(), worldCore, fake); err == nil {
		t.Fatal("exit 1 (usage/infra) accepted because stdout parsed")
	}
}

// Derived negatives: each is a recorded jq transform of the REAL pristine
// fixture, and each asserts the refusal text of exactly one guard.
func TestParseQualityInterfaceDerivedNegatives(t *testing.T) {
	for file, want := range map[string]string{
		"derived_quality_v2_holds_v1.json": "interface identity v2 absent or malformed",
		"derived_quality_v1_absent.json":   "interface hash v1 absent or malformed",
		"derived_quality_schema_v2.json":   `schema "ailang.package-quality/v2"`,
	} {
		b, err := os.ReadFile(filepath.Join("testdata", file))
		if err != nil {
			t.Fatal(err)
		}
		_, perr := ParseQualityInterface(b)
		if perr == nil || !strings.Contains(perr.Error(), want) {
			t.Errorf("%s: want %q, got %v", file, want, perr)
		}
	}
}

func fakeBinary(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "ailang")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

// tinyBounds must absorb fake-binary spawn latency under full-suite load: at
// timeout 1 s the descendant and cap tests failed 3/3 full `go test ./...`
// runs (the shell had not produced output or its pid within 1 s) while
// passing in isolation. 3 s / 1 s is measured green under load.
var tinyBounds = queryBounds{timeout: 3 * time.Second, waitDelay: time.Second, maxOutput: 4096}

// guard is every runWithin limit: well above timeout+waitDelay, well below the
// 30 s the hang/descendant doubles sleep, so a mutant REDS instead of hanging.
const guard = 10 * time.Second

// runWithin fails the test (instead of hanging it) if the call overruns.
func runWithin(t *testing.T, limit time.Duration, f func() error) (error, time.Duration) {
	t.Helper()
	done := make(chan error, 1)
	start := time.Now()
	go func() { done <- f() }()
	select {
	case err := <-done:
		return err, time.Since(start)
	case <-time.After(limit):
		t.Fatalf("QueryInterface did not return within %s", limit)
		return nil, 0
	}
}

func TestQueryInterfaceTimesOutOnAHangingBinary(t *testing.T) {
	bin := fakeBinary(t, "exec sleep 30\n")
	err, took := runWithin(t, guard, func() error {
		_, err := queryInterface(context.Background(), t.TempDir(), worldCore, bin, tinyBounds)
		return err
	})
	var te *QueryTimeoutError
	if !errors.As(err, &te) || te.Timeout != tinyBounds.timeout {
		t.Fatalf("want QueryTimeoutError(%s), got %v", tinyBounds.timeout, err)
	}
	t.Logf("hang: named timeout after %s", took)
}

// The direct child (sh) is killed at the deadline; its descendant keeps the
// inherited stdout open. WaitDelay must release Wait anyway. The descendant is
// NOT killed by this code (no process group) - asserted, then reaped here.
func TestQueryInterfaceReturnsWhileADescendantHoldsStdout(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "pid")
	bin := fakeBinary(t, "sleep 30 &\necho $! > '"+pidFile+"'\nwait\n")
	err, took := runWithin(t, guard, func() error {
		_, err := queryInterface(context.Background(), t.TempDir(), worldCore, bin, tinyBounds)
		return err
	})
	var te *QueryTimeoutError
	if !errors.As(err, &te) {
		t.Fatalf("want QueryTimeoutError, got %v", err)
	}
	raw, rerr := os.ReadFile(pidFile)
	if rerr != nil {
		t.Fatalf("instrument: descendant pid not recorded: %v", rerr)
	}
	pid := strings.TrimSpace(string(raw))
	alive := exec.Command("kill", "-0", pid).Run() == nil
	t.Logf("descendant: named timeout after %s; descendant pid %s alive after return: %v", took, pid, alive)
	_ = exec.Command("kill", "-9", pid).Run()
}

func TestQueryInterfaceCapsStdout(t *testing.T) {
	bin := fakeBinary(t, "head -c 100000 /dev/zero\n")
	// A long internal timeout: the cap, not the deadline, must produce the error.
	capOnly := queryBounds{timeout: 60 * time.Second, waitDelay: time.Second, maxOutput: 4096}
	err, _ := runWithin(t, guard, func() error {
		_, err := queryInterface(context.Background(), t.TempDir(), worldCore, bin, capOnly)
		return err
	})
	if !errors.Is(err, ErrQualityOutputOverflow) {
		t.Fatalf("want ErrQualityOutputOverflow, got %v", err)
	}
}

func TestQueryInterfaceUsesFiniteProductionBounds(t *testing.T) {
	if defaultQueryBounds.timeout <= 0 || defaultQueryBounds.timeout >= 120*time.Second ||
		defaultQueryBounds.waitDelay <= 0 || defaultQueryBounds.maxOutput <= 0 {
		t.Fatalf("production bounds not finite or not inside run_bounded 120: %+v", defaultQueryBounds)
	}
}

func TestQueryInterfaceHonoursCallerCancellation(t *testing.T) {
	bin := fakeBinary(t, "exec sleep 30\n")
	long := queryBounds{timeout: 6 * time.Second, waitDelay: time.Second, maxOutput: 4096}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	err, took := runWithin(t, 5*time.Second, func() error {
		_, err := queryInterface(ctx, t.TempDir(), worldCore, bin, long)
		return err
	})
	if !errors.Is(err, context.DeadlineExceeded) || errors.As(err, new(*QueryTimeoutError)) {
		t.Fatalf("want the caller's context error, not the internal bound: %v", err)
	}
	t.Logf("caller cancel: returned after %s", took)
}

// fixtureBinary is a fake ailang that prints a checked-in fixture and exits rc.
func fixtureBinary(t *testing.T, fixture string, rc int) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("testdata", fixture))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("instrument: fixture %s: %v", fixture, err)
	}
	return fakeBinary(t, "cat '"+abs+"'\nexit "+strconv.Itoa(rc)+"\n")
}

// T19 (r2, glm): exit 2 is accepted ONLY with the spurious PUB015 gate; exit 0
// only with no gate. Positive controls first, so a rule that refuses
// everything is also red.
func TestQueryInterfaceRefusesExit2WithANonSpuriousGate(t *testing.T) {
	for _, c := range []struct {
		fixture string
		rc      int
	}{
		{"quality_world_core_pristine.json", 2},     // real --no-run: [PUB015 "_smoke.ail failed"]
		{"quality_world_core_pristine_run.json", 0}, // real run mode: []
	} {
		id, err := QueryInterface(context.Background(), t.TempDir(), worldCore, fixtureBinary(t, c.fixture, c.rc))
		if err != nil || id.Signatures != 53 {
			t.Fatalf("control %s exit %d: want accepted, got %+v %v", c.fixture, c.rc, id, err)
		}
	}
	for _, c := range []struct {
		fixture string
		rc      int
		want    string
	}{
		// generated: compile failed, gates [PUB000, PUB015], no hash_v2 — must be
		// refused BY THE GATE RULE, not by the later v2-absent guard.
		{"quality_world_core_v2_unbuilt.json", 2, "exit 2 with gates [PUB000 PUB015]"},
		// derived: keeps hash_v2 (so nothing downstream would refuse) + a PUB000 gate
		{"derived_quality_extra_gate.json", 2, "exit 2 with gates [PUB015 PUB000]"},
		// exit 0 must carry no gate: the real --no-run doc carries PUB015
		{"quality_world_core_pristine.json", 0, "exit 0 with gates [PUB015]"},
		// derived: PUB015 but not the spurious message — the code alone is not enough
		{"derived_quality_pub015_other_msg.json", 2, "exit 2 with gates [PUB015]"},
	} {
		_, err := QueryInterface(context.Background(), t.TempDir(), worldCore, fixtureBinary(t, c.fixture, c.rc))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s exit %d: want refusal naming %q, got %v", c.fixture, c.rc, c.want, err)
		}
	}
}

// T20 (r2, astra): overflow must CANCEL the child, not merely error the
// writer. The double overflows stdout then sleeps; the internal timeout
// (60 s) is far beyond the 10 s return guard, so only overflow-triggered
// cancellation returns in time.
func TestQueryInterfaceOverflowCancelsTheChild(t *testing.T) {
	b := queryBounds{timeout: 60 * time.Second, waitDelay: time.Second, maxOutput: 4096}
	for stream, redirect := range map[string]string{"stdout": "", "stderr": " >&2"} {
		// stderr's cap is the fixed 64 KiB; 100000 bytes crosses both caps.
		bin := fakeBinary(t, "head -c 100000 /dev/zero"+redirect+"\nexec sleep 30\n")
		err, took := runWithin(t, guard, func() error {
			_, err := queryInterface(context.Background(), t.TempDir(), worldCore, bin, b)
			return err
		})
		if !errors.Is(err, ErrQualityOutputOverflow) {
			t.Fatalf("%s: want ErrQualityOutputOverflow, got %v", stream, err)
		}
		t.Logf("overflow on %s: returned after %s (internal timeout %s)", stream, took, b.timeout)
	}
}
