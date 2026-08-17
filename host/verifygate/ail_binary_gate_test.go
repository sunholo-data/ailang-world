package verifygate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// pinned is the released delegate every shim arm runs the real gate against. It MUST come from
// AILANG_BIN, never from a literal: the pinned binary's tmp-rooted location is a convention of one
// dev rig and exists on no runner. CI's go-verify job exports AILANG_BIN (ci.yml) after installing and
// sha256-verifying the pinned v0.30.0 release, which is the same contract verify_go.sh asserts
// before it will run this package at all.
//
// The sibling packages read AILANG_BIN and `t.Skip` when it is unset. This package must NOT: a
// silent skip is the exact false-green class this whole milestone exists to close, so an unset or
// wrong AILANG_BIN is a loud t.Fatal (see requirePinned).
var (
	repoRoot = findRepoRoot()
	pinned   = os.Getenv("AILANG_BIN")
)

func findRepoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func requirePinned(t *testing.T) {
	t.Helper()
	if pinned == "" {
		t.Fatal("AILANG_BIN is unset — the shim arms need the pinned released delegate to run the real gate. " +
			"Never skip: a silent skip here is the false-green class this milestone closes. " +
			"verify_go.sh already refuses to run without it — export the pinned released binary named in CLAUDE.md.")
	}
	// Output(), never CombinedOutput(): the binary writes operational warnings to stderr, and
	// merging them prefixes the banner so this HasPrefix check fails on a correct pinned
	// release (measured 2026-08-17 — an `Observatory: …MB` warning failed all 17 arms here).
	out, err := exec.Command(pinned, "--version").Output()
	if err != nil || !strings.HasPrefix(string(out), "AILANG v0.30.0") {
		t.Fatalf("pinned delegate %q unavailable or wrong (never skip): err=%v output=%q", pinned, err, out)
	}
}

func runGate(t *testing.T, env map[string]string) (int, string) {
	t.Helper()
	cmd := exec.Command(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	cmd.Dir = repoRoot
	blocked := map[string]bool{
		"AILANG_BIN": true, "WORLD_PKG_AILANG_BIN": true,
		"AILANG_SHIM_VERSION_LINE": true, "AILANG_SHIM_DELEGATE": true,
		// Ambient AILANG_Z3_PATH is stripped so every arm resolves the solver the way the
		// LANE does (the binary's own search: PATH + /usr/local/bin, /usr/bin, /snap/bin,
		// /opt/homebrew/bin). `blocked` filters os.Environ() only — an explicit entry in the
		// `env` map below still reaches the child, which is how a Z3-absent control is armed
		// deterministically instead of relying on os/exec's dedup order.
		"AILANG_Z3_PATH": true,
	}
	cmd.Env = make([]string, 0, len(os.Environ())+len(env))
	for _, item := range os.Environ() {
		if !blocked[strings.SplitN(item, "=", 2)[0]] {
			cmd.Env = append(cmd.Env, item)
		}
	}
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	if err == nil {
		return 0, output.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String()
	}
	t.Fatalf("start verify gate: %v", err)
	return -1, output.String()
}

func shimEnv(line string) map[string]string {
	return map[string]string{
		"AILANG_BIN":               filepath.Join(repoRoot, "scripts", "testdata", "ailang_version_shim.sh"),
		"AILANG_SHIM_VERSION_LINE": line,
		"AILANG_SHIM_DELEGATE":     pinned,
		"WORLD_PKG_AILANG_BIN":     pinned,
	}
}

func requireRefusal(t *testing.T, line, code string) {
	t.Helper()
	requirePinned(t)
	rc, out := runGate(t, shimEnv(line))
	if rc != 1 || !strings.Contains(out, "["+code+"]") {
		t.Fatalf("want rc=1 [%s], got rc=%d\n%s", code, rc, out)
	}
}

// Markers the accept-arms assert on. The first three are emitted BEFORE anything Z3-dependent; the
// fourth (passedMarker) is the gate's terminal success line and needs a solver in the lane.
const (
	refusalMarker = "AILANG_BIN refused" // the version block refused; emitted only by _refuse
	leg1Marker    = "── Leg 1"           // the version block accepted and the gate proceeded
	// The signature a NON-delegating shim produces: ai-check returns empty output at the first
	// module, so the gate cannot parse its JSON. Measured to differ from the Z3-absent signature
	// (`required identity … MISSING from verify.results[]`), which is what makes delegation
	// provable without Z3.
	noDelegateMarker = "could not parse ai-check JSON"
	// The gate's terminal success line (verify_ail.sh:304). Unambiguous: the only other "gate
	// PASSED" line is `world package gate PASSED`, which does not contain this substring.
	passedMarker = "verify gate PASSED"
)

// requireProceeded asserts the version block ACCEPTED, the real gate ran on the real delegate, and
// THE GATE PASSED.
//
// The fourth assertion is what 9/OD-11 bought. It was absent because these tests run inside
// `go test ./...`, i.e. CI's go-verify job, which installed no Z3 — and without a solver `ai-check`
// reports verify.available=false and every contract vanishes from verify.results[] (V20/V27), so
// the gate reds at leg 1 for a reason with nothing to do with the version predicate. The three
// remaining markers were re-aimed at that lane and are all emitted before Z3 matters. Measured, the
// cost of that portability was exact and total: with AILANG_Z3_PATH pointed at a missing file, a
// PASSING gate (rc=0) and a FAILING gate (rc=1) produce IDENTICAL values for all three — refused=0,
// leg1=1, noDelegate=0 — so the accept-arms were green in the lane where the gate they drive was
// red. `verify gate PASSED` is the only marker that separates them.
//
// Mark ratified the install ("Yes you can install z3 on cicd", 9/OD-11); ci.yml now installs the
// same pinned Z3 in BOTH jobs, so the assertion is reachable everywhere this package runs. If a
// lane loses its solver, TestSolverAvailableInThisLane says so by name before these arms confuse
// anyone.
// checkProceeded is the accept contract as a PURE PREDICATE, returning an error instead of failing
// a test. The split exists so the contract can be pointed at output the suite chooses — which is
// the only way to prove the passedMarker assertion is load-bearing. In a solver-bearing lane every
// accept-arm passes whether or not that assertion is present, so a mutation neutering it survives
// against every arm above; TestAcceptContractRejectsASolverlessGate feeds this function the output
// of a deliberately solverless gate and requires it to REJECT, which no healthy-lane arm can do.
func checkProceeded(label, out string) error {
	if strings.Contains(out, refusalMarker) {
		return fmt.Errorf("%s: version block refused a release token\n%s", label, out)
	}
	if !strings.Contains(out, leg1Marker) {
		return fmt.Errorf("%s: gate never reached leg 1 — it did not proceed past the version block\n%s", label, out)
	}
	if strings.Contains(out, noDelegateMarker) {
		return fmt.Errorf("%s: shim did not delegate — ai-check produced no parseable output\n%s", label, out)
	}
	if !strings.Contains(out, passedMarker) {
		return fmt.Errorf("%s: gate proceeded but did NOT pass — %q absent. If this lane has no Z3 the whole "+
			"gate reds at leg 1; see TestSolverAvailableInThisLane and the Z3 install steps in "+
			".github/workflows/ci.yml\n%s", label, passedMarker, out)
	}
	return nil
}

func requireProceeded(t *testing.T, label, out string) {
	t.Helper()
	if err := checkProceeded(label, out); err != nil {
		t.Fatal(err)
	}
}

// TestAcceptContractRejectsASolverlessGate is the committed, always-running proof that the fourth
// marker earns its place — the non-vacuity claim of 9/OD-11 made permanent rather than left to a
// one-shot drill on a tree that no longer exists.
//
// It arms the Z3-absent control the way runGate's `blocked` map intends: ambient AILANG_Z3_PATH is
// stripped so arms resolve the solver the way their LANE does, but an explicit entry in the env map
// still reaches the child. So this is the ONLY way the suite can observe a failing gate, and it is
// the exact lane CI job 2 was before Mark ratified the install.
//
// The assertion is two-sided on purpose. First it re-measures the defect: all THREE legacy markers
// are satisfied by this failing gate, so the pre-9/OD-11 contract would have accepted it. Then it
// requires the current contract to reject it. Neuter the passedMarker check and this test is the
// one that reds.
func TestAcceptContractRejectsASolverlessGate(t *testing.T) {
	requirePinned(t)
	env := shimEnv("AILANG v0.33.0")
	// Assembled, never a rig-absolute literal (TestNoRigAbsolutePaths).
	env["AILANG_Z3_PATH"] = filepath.Join(t.TempDir(), "nosuch", "z3")

	rc, out := runGate(t, env)
	if rc == 0 {
		t.Fatalf("instrument failure: the solverless control did not disarm the gate — rc=0, so "+
			"there is no failing gate to test the contract against\n%s", out)
	}

	// The defect 9/OD-11 names, re-measured every run: the legacy three cannot see this failure.
	if strings.Contains(out, refusalMarker) {
		t.Fatalf("solverless control: gate refused the token instead of failing later — wrong failure\n%s", out)
	}
	if !strings.Contains(out, leg1Marker) {
		t.Fatalf("solverless control: gate never reached leg 1 — wrong failure\n%s", out)
	}
	if strings.Contains(out, noDelegateMarker) {
		t.Fatalf("solverless control: shim failed to delegate — wrong failure\n%s", out)
	}

	// And the current contract must REJECT it. This is the kill arm for the passedMarker assertion.
	if err := checkProceeded("solverless control", out); err == nil {
		t.Fatal("the accept contract ACCEPTED a gate that FAILED (rc!=0): every marker it checks was " +
			"satisfied by a solverless run, which is exactly the state 9/OD-11 closed. The " +
			"`verify gate PASSED` assertion in checkProceeded is not doing its job.")
	}
}

// TestSolverAvailableInThisLane is the self-naming instrument check for 9/OD-11. Every accept-arm
// above now asserts `verify gate PASSED`, which is unreachable without a solver; if a lane loses
// Z3 this test says exactly that, by name, instead of four arms failing with a leg-1 contract error.
//
// It is SELF-CONTROLLED: arm B points AILANG_Z3_PATH at a path that does not exist and requires
// verify.available to be FALSE. That is not decoration — it is the anti-vacuity proof, and unlike a
// source mutation it runs on every CI run forever. Measured on the rig: the binary does not link
// libz3 (`otool -L` names none) and does not fall back to PATH when AILANG_Z3_PATH is set and
// missing (it carries the string "AILANG_Z3_PATH set to %q but file not found"), so arm B is a
// faithful model of a solverless runner. PATH manipulation is NOT a valid control here: removing
// the toolchain directory also removes `go`, and the gate then reds for an unrelated reason.
//
// DECLARED RESIDUAL (rule 3j — a named gap is cheap, an assumed one is not): arm A's
// `verify.available` assertion is a DIAGNOSTIC, not a guard. It fires exactly when the lane has no
// solver, which is a lane this suite deliberately cannot create for itself — probe() strips ambient
// AILANG_Z3_PATH so each arm resolves the solver the way its LANE does. So a mutation neutering
// arm A SURVIVES here by construction, and it was measured surviving rather than assumed away. What
// is protected is the pair: arm B (a missing solver must produce available=false) and the
// arms-differ check together prove the probe reads the solver at all, so arm A cannot be green
// against a probe that reads nothing.
func TestSolverAvailableInThisLane(t *testing.T) {
	requirePinned(t)

	type aiCheck struct {
		Check struct {
			Passed bool `json:"passed"`
		} `json:"check"`
		Verify struct {
			Available bool `json:"available"`
			Verified  int  `json:"verified"`
		} `json:"verify"`
	}

	probe := func(t *testing.T, z3 string) aiCheck {
		t.Helper()
		cmd := exec.Command(pinned, "ai-check", "-timeout", "5s", "world/contracts.ail")
		cmd.Dir = repoRoot
		env := make([]string, 0, len(os.Environ())+1)
		for _, item := range os.Environ() {
			if strings.SplitN(item, "=", 2)[0] != "AILANG_Z3_PATH" {
				env = append(env, item)
			}
		}
		if z3 != "" {
			env = append(env, "AILANG_Z3_PATH="+z3)
		}
		cmd.Env = env
		// Exit status is ADVISORY here (V10/V20) exactly as in verify_ail.sh; the JSON is
		// authoritative. Only an unparseable body is instrument failure.
		raw, _ := cmd.Output()
		start := bytes.IndexByte(raw, '{')
		if start < 0 {
			t.Fatalf("instrument failure: ai-check produced no JSON object (%d bytes): %q", len(raw), raw)
		}
		var parsed aiCheck
		if err := json.Unmarshal(raw[start:], &parsed); err != nil {
			t.Fatalf("instrument failure: ai-check JSON did not parse: %v\n%s", err, raw)
		}
		return parsed
	}

	// Arm A — this lane, as the gate will see it.
	armA := probe(t, "")
	// Known-positive control FIRST: if check.passed is false the run is broken and an
	// available=false verdict would be attributed to the wrong cause.
	if !armA.Check.Passed {
		t.Fatalf("instrument failure: ai-check on world/contracts.ail did not pass its CHECK phase; "+
			"the solver verdict below would be meaningless (%+v)", armA)
	}
	if !armA.Verify.Available {
		t.Fatalf("NO SMT SOLVER IN THIS LANE: ai-check reports verify.available=false, so every " +
			"contract vanishes from verify.results[] (V20/V27) and scripts/verify_ail.sh cannot " +
			"reach `verify gate PASSED`. Install the pinned Z3 (see the workflow-level Z3_VER/Z3_SHA " +
			"and the two `Install Z3` steps in .github/workflows/ci.yml), or locally: " +
			"brew install z3 / apt install z3, or set AILANG_Z3_PATH.")
	}
	if armA.Verify.Verified < 1 {
		t.Fatalf("solver is available but verified %d contracts in world/contracts.ail, want >=1 (%+v)",
			armA.Verify.Verified, armA)
	}

	// Arm B — the negative control. Assembled, never a rig-absolute literal (TestNoRigAbsolutePaths).
	missing := filepath.Join(t.TempDir(), "nosuch", "z3")
	armB := probe(t, missing)
	if armB.Verify.Available {
		t.Fatalf("instrument is NOT discriminating: verify.available stayed true with AILANG_Z3_PATH "+
			"pointed at %q, so arm A's green proves nothing about this lane", missing)
	}
	if armA.Verify.Available == armB.Verify.Available {
		t.Fatalf("arms do not differ (A=%v B=%v) — the probe is not reading the solver",
			armA.Verify.Available, armB.Verify.Available)
	}
}

func TestKnownPositiveDelegates(t *testing.T) {
	requirePinned(t)
	_, out := runGate(t, shimEnv("AILANG v0.33.0"))
	requireProceeded(t, "delegating release arm", out)
}

func TestDirtyBuildRefused(t *testing.T) {
	requireRefusal(t, "AILANG v0.33.0-105-g38e119db1-dirty", "DEV_BUILD")
}

func TestPlainDirtyRefused(t *testing.T) {
	requireRefusal(t, "AILANG v0.33.0-dirty", "DEV_BUILD")
}

func TestGitDescribeRefused(t *testing.T) {
	requireRefusal(t, "AILANG v0.30.0-205-g54d6bd191", "DEV_BUILD")
}

func writeExecutable(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ailang")
	if err := os.WriteFile(path, []byte("#!/usr/bin/env bash\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNoVersionOutput(t *testing.T) {
	bin := writeExecutable(t, `if [ "${1:-}" = --version ]; then exit 0; fi; exec "${AILANG_SHIM_DELEGATE:?}" "$@"`)
	env := shimEnv("")
	env["AILANG_BIN"] = bin
	rc, out := runGate(t, env)
	if rc != 1 || !strings.Contains(out, "[NO_VERSION_OUTPUT]") {
		t.Fatalf("rc=%d\n%s", rc, out)
	}
}

func TestNoVersionToken(t *testing.T) { requireRefusal(t, "AILANG", "NO_VERSION_TOKEN") }
func TestNotARelease(t *testing.T)    { requireRefusal(t, "AILANG 0.33.0", "NOT_A_RELEASE") }

func TestUnresolvable(t *testing.T) {
	for _, tc := range []struct{ name, path string }{
		{"Missing", filepath.Join(t.TempDir(), "missing")},
		{"NoExec", func() string {
			p := filepath.Join(t.TempDir(), "ailang")
			if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
			return p
		}()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env := shimEnv("")
			env["AILANG_BIN"] = tc.path
			rc, out := runGate(t, env)
			if rc != 1 || !strings.Contains(out, "[UNRESOLVABLE]") {
				t.Fatalf("rc=%d\n%s", rc, out)
			}
		})
	}
}

func releaseVerdict(tok string) string {
	dev := regexp.MustCompile(`(-dirty$|-[0-9]+-g[0-9a-f]+)`)
	shape := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.]+)?$`)
	if tok == "" {
		return "NO_VERSION_TOKEN"
	}
	if dev.MatchString(tok) {
		return "DEV_BUILD"
	}
	if !shape.MatchString(tok) {
		return "NOT_A_RELEASE"
	}
	return "OK"
}

func TestUpstreamReleaseCorpus(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "testdata", "upstream_release_tags.txt"))
	if err != nil {
		t.Fatal(err)
	}
	tags := strings.Fields(string(data))
	if len(tags) < 60 {
		t.Fatalf("release corpus must be non-empty and >=60; got %d", len(tags))
	}
	for _, tag := range tags {
		if got := releaseVerdict(tag); got != "OK" {
			t.Errorf("upstream release %q rejected as %s", tag, got)
		}
	}
	for _, tok := range []string{"", "line", "AILANG", "0.33.0", "v0.33", "vX.Y.Z", "v0.33.0-dirty", "v0.33.0.0", "v0.30.0-205-g54d6bd191"} {
		if got := releaseVerdict(tok); got == "OK" {
			t.Errorf("negative control %q accepted", tok)
		}
	}
}

func TestInScriptControl(t *testing.T) {
	requirePinned(t)
	src, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	if err != nil {
		t.Fatal(err)
	}
	old, replacement := "v0.30.0|OK", "v0.30.0|DEV_BUILD"
	if strings.Count(string(src), old) != 1 {
		t.Fatalf("control anchor count=%d", strings.Count(string(src), old))
	}
	mutant := strings.Replace(string(src), old, replacement, 1)
	path, err := os.CreateTemp(filepath.Join(repoRoot, "scripts"), ".verify-control-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	name := path.Name()
	defer os.Remove(name)
	if _, err := path.WriteString(mutant); err != nil {
		t.Fatal(err)
	}
	if err := path.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(name, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(name)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "AILANG_BIN="+pinned, "WORLD_PKG_AILANG_BIN="+pinned)
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "release-verdict instrument broken") || strings.Contains(string(out), "── Leg 1") {
		t.Fatalf("inverted control was not an early loud failure: err=%v\n%s", err, out)
	}
}

// TestReleaseChangeNotice pins 9/CF-A-2 with THREE arms, not two.
//
// An observability feature cannot be proven by a mutation that reds, so non-vacuity here is a claim
// about arms DIFFERING: the quiet arms and the firing arm must disagree, which is what proves the
// line reads the resolved binary rather than a constant.
//
// The two quiet arms are not redundant. Legs 1-2 resolve a different release in each real lane —
// CI's job 1 installs `releases/latest` on PATH and exports no AILANG_BIN, while a local operator
// exports the documented v0.30.0 pin. An equality test against a single recorded observation is
// quiet in whichever lane it was written for and fires on EVERY run of the other, which is the
// always-firing defect 9/CF-A-2 exists to remove, merely relocated. ArmLocalPin is the regression
// pin for that: it was measured firing `moved from 'v0.33.0' to 'v0.30.0'` before the membership fix.
func TestReleaseChangeNotice(t *testing.T) {
	requirePinned(t)
	const notice = "UNRECOGNISED RELEASE"

	// Slash-free names: `go test -run` splits its pattern on `/`.
	quiet := []struct {
		name    string
		version string
	}{
		{"ArmCILatest", "AILANG v0.33.0"}, // upstream releases/latest — CI's legs 1-2
		{"ArmLocalPin", "AILANG v0.30.0"}, // the documented local pin — CLAUDE.md
	}
	quietCounts := map[string]int{}
	for _, arm := range quiet {
		_, out := runGate(t, shimEnv(arm.version))
		requireProceeded(t, arm.name, out)
		n := strings.Count(out, notice)
		quietCounts[arm.name] = n
		if n != 0 {
			t.Fatalf("%s: expected 0 notices, got %d\n%s", arm.name, n, out)
		}
	}

	_, outB := runGate(t, shimEnv("AILANG v0.34.0"))
	requireProceeded(t, "ArmUnrecognised", outB)
	countB := strings.Count(outB, notice)
	if countB != 1 {
		t.Fatalf("ArmUnrecognised: expected exactly 1 notice, got %d\n%s", countB, outB)
	}
	// The notice must name the token it refused AND the expected set it checked against.
	for _, want := range []string{"v0.34.0", "v0.33.0", "v0.30.0"} {
		if !strings.Contains(outB, want) {
			t.Fatalf("ArmUnrecognised: notice omits %q\n%s", want, outB)
		}
	}
	// The load-bearing assertion: the arms DIFFER. Without this a constant-folded notice
	// (MUT-VL-NOTICE-CONST) silences every arm and each arm above still reads self-consistently.
	for name, n := range quietCounts {
		if n == countB {
			t.Fatalf("arms do not differ: %s=%d, ArmUnrecognised=%d — the notice is not reading the binary", name, n, countB)
		}
	}
}

// TestEmptyExpectedReleaseSetFailsLoudly pins the refusal branch guarding the fixture. An empty set
// makes `grep -qxF` match nothing, so WITHOUT this branch the notice would fire on every run rather
// than never — but the branch is what turns a silently-emptied fixture into a loud stop instead of
// noise nobody reads. Reachable only by pointing the script at an empty file, so the drill is
// TestInScriptControl's: run a temp COPY of the real script with the fixture path redirected. That
// keeps it a test of the artifact rather than of a re-derivation, and never mutates a tracked file.
func TestEmptyExpectedReleaseSetFailsLoudly(t *testing.T) {
	requirePinned(t)
	src, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	if err != nil {
		t.Fatal(err)
	}
	// Anchor on the READ, not the bare path: the path also appears in the fatal message and in the
	// notice's remediation line, so a bare-path redirect would be ambiguous (measured: 3 occurrences).
	const fixtureRead = `grep -vE '^[[:space:]]*(#|$)' scripts/testdata/ailang_release_observed.txt`
	if n := strings.Count(string(src), fixtureRead); n != 1 {
		t.Fatalf("fixture-read anchor count=%d, want 1 (the redirect below would be ambiguous)", n)
	}
	empty := filepath.Join(t.TempDir(), "empty.txt")
	if err := os.WriteFile(empty, []byte("# only a comment\n\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mutant := strings.Replace(string(src), fixtureRead,
		`grep -vE '^[[:space:]]*(#|$)' `+empty, 1)

	f, err := os.CreateTemp(filepath.Join(repoRoot, "scripts"), ".verify-emptyset-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err := f.WriteString(mutant); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(name, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(name)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "AILANG_BIN="+pinned, "WORLD_PKG_AILANG_BIN="+pinned)
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "lists no releases") || strings.Contains(string(out), "── Leg 1") {
		t.Fatalf("empty expected-release set was not an early loud failure: err=%v\n%s", err, out)
	}
}

// TestNoRigAbsolutePaths is the regression guard for the defect that red CI on this milestone's
// first push: `pinned` was a hardcoded tmp-rooted literal, a convention of one dev rig that exists
// on no runner, so all ten shim-driven tests t.Fatal'd in the go-verify job while every local
// gate was rc=0. The local green was real and answered the wrong question. A rig-absolute literal in
// this package is always wrong — the binary's location is CI's to choose and AILANG_BIN's to carry.
func TestNoRigAbsolutePaths(t *testing.T) {
	entries, err := filepath.Glob(filepath.Join(repoRoot, "host", "verifygate", "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("instrument failure: globbed zero .go files in host/verifygate")
	}
	// Known-positive control: the scan must be able to see a string that IS present.
	var sawControl bool
	for _, path := range entries {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(src), "AILANG_BIN") {
			sawControl = true
		}
		// The needles are ASSEMBLED, never written literally: a scanner that also scans its own
		// source matches its own pattern list and reports three findings against a clean file.
		for _, bad := range []string{"/tmp" + "/ailang", "/Us" + "ers/", "/home" + "/runner/"} {
			if strings.Contains(string(src), bad) {
				t.Errorf("%s contains rig-absolute path %q — resolve the binary from AILANG_BIN instead", filepath.Base(path), bad)
			}
		}
	}
	if !sawControl {
		t.Fatal("instrument failure: known-positive control \"AILANG_BIN\" not found in any scanned file")
	}
}

// TestExpectedReleaseSetIsNonEmpty is the anti-vacuity guard on the fixture itself: an empty or
// all-comment file would silence the notice forever while every arm above still passed.
func TestExpectedReleaseSetIsNonEmpty(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "testdata", "ailang_release_observed.txt"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []string
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entries = append(entries, line)
	}
	if len(entries) < 2 {
		t.Fatalf("expected-release set has %d entries (%v); both real lanes must be listed or the notice fires on every run of one of them", len(entries), entries)
	}
	// Known-positive control: the file is not merely non-empty, it contains the two lane tokens.
	for _, want := range []string{"v0.33.0", "v0.30.0"} {
		if !slices.Contains(entries, want) {
			t.Fatalf("expected-release set %v is missing %q", entries, want)
		}
	}
}

// TestFixtureDiscrimination is the meta-control the handoff demands: omitting WORLD_PKG_AILANG_BIN
// makes it fall back to AILANG_BIN, so the gate reds at the world-package leg on compiler BYTES —
// rc=1 in both arms, which would let a careless test conclude "the version assertion fired" when it
// never ran. The essential claim is therefore that the failure is NOT the version predicate's, and
// that claim is Z3-independent. The precise `wrong compiler version` reason lives at leg 3; since
// 9/OD-11 both CI jobs install Z3, so leg 3 is reachable everywhere this package runs and the reason
// is asserted unconditionally.
func TestFixtureDiscrimination(t *testing.T) {
	requirePinned(t)
	env := shimEnv("AILANG v0.33.0")
	delete(env, "WORLD_PKG_AILANG_BIN")
	rc, out := runGate(t, env)
	if rc != 1 {
		t.Fatalf("fixture discrimination: want rc=1, got rc=%d\n%s", rc, out)
	}
	for _, marker := range []string{refusalMarker, "[DEV_BUILD]", "[NOT_A_RELEASE]", noDelegateMarker} {
		if strings.Contains(out, marker) {
			t.Fatalf("fixture discrimination: rc=1 came from %q, not from the compiler-identity check — "+
				"a test wired this way would credit the version predicate for a red it did not cause\n%s", marker, out)
		}
	}
	if !strings.Contains(out, leg1Marker) {
		t.Fatalf("fixture discrimination: gate never proceeded past the version block\n%s", out)
	}
	// Unconditional since 9/OD-11. Before the Z3 install this was guarded by `Contains(…, "── Leg 3")`
	// because the solverless go-verify job never got past leg 1 — measured: `── Leg 3` count 0 there
	// vs 1 locally, so in CI the guarded body had NEVER executed. With a solver in both jobs the leg
	// is always reached and the guard would only hide a regression.
	if !strings.Contains(out, "── Leg 3") {
		t.Fatalf("fixture discrimination: gate never reached leg 3 — it stopped earlier, so the "+
			"compiler-identity check below was never exercised (a Z3-less lane looks exactly like "+
			"this; see TestSolverAvailableInThisLane)\n%s", out)
	}
	if !strings.Contains(out, "wrong compiler version") {
		t.Fatalf("fixture discrimination: reached leg 3 but not via the compiler-version check\n%s", out)
	}
}

// TestZ3PinDeclaredOnceAndInstalledInBothJobs guards the structure 9/OD-11 bought.
//
// Two claims, and the second is the one with teeth: (1) the Z3 version and sha256 are declared
// EXACTLY ONCE, at workflow scope, so the two jobs cannot install different solvers — a job-local
// re-declaration would shadow the pin for that job only, which is 9/CF-A-2's "two lanes resolve
// different things" all over again; (2) BOTH jobs actually install it. Claim 2 is why this test is
// not decoration: if job 2's install is dropped, every accept-arm's `verify gate PASSED` becomes
// unreachable and 10 tests red with a leg-1 contract error that names no cause. This names the cause.
//
// DECLARED RESIDUAL (found by the evaluator, reproduced first-party): this is a STATIC text scan, so
// it sees the install command's TEXT, never whether the step RUNS. A step-level `if:` whose
// expression is always false at runtime — e.g. `contains(github.event.head_commit.message, '<a
// marker nobody writes>')` — disables job 2's install while leaving every byte this test counts
// intact, and `actionlint` is green too (it flags a literal `if: false`, not a non-constant
// always-false expression). The bar the milestone actually has to clear still holds, by two
// DYNAMIC backstops that run in the job itself: TestSolverAvailableInThisLane reds by name
// (`NO SMT SOLVER IN THIS LANE`), and every requireProceeded accept-arm reds on the absent
// passedMarker. So this test's claim (2) is narrower than "job 2 will have a solver" — it is "the
// pin is declared once and both install steps are PRESENT" — and that narrowing is stated here
// rather than left for a reader to discover.
func TestZ3PinDeclaredOnceAndInstalledInBothJobs(t *testing.T) {
	path := filepath.Join(repoRoot, ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	// Known-positive control: the scan must be able to see strings that ARE present, or a
	// zero-count assertion below could be satisfied by reading the wrong file.
	for _, control := range []string{"ailang-verify:", "go-verify:", "./scripts/verify_go.sh"} {
		if !strings.Contains(src, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", path, control)
		}
	}
	for _, tc := range []struct {
		needle string
		want   int
		why    string
	}{
		{"Z3_VER:", 1, "the version must be declared exactly once, at workflow scope"},
		{"Z3_SHA:", 1, "the sha256 must be declared exactly once, at workflow scope"},
		{"Z3_VER=", 0, "a job-local Z3_VER= shadows the workflow pin for that job only"},
		{"Z3_SHA=", 0, "a job-local Z3_SHA= shadows the workflow pin for that job only"},
	} {
		if got := strings.Count(src, tc.needle); got != tc.want {
			t.Errorf("ci.yml: count(%q)=%d, want %d — %s", tc.needle, got, tc.want, tc.why)
		}
	}

	// The install assertion is LINE-EXACT, not a substring count, and that is a repair rather than
	// a style choice: `strings.Count` over the install command SURVIVED a mutation redirecting one
	// job's solver to /usr/local/bin/z3x, because the mutated line still CONTAINS the needle. A
	// prefix-shaped needle cannot detect a suffix-shaped mutation. Counting whole trimmed lines can.
	const installLine = `sudo install -m 0755 "${Z3_DIR}/bin/z3" /usr/local/bin/z3`
	installs := 0
	for _, line := range strings.Split(src, "\n") {
		if strings.TrimSpace(line) == installLine {
			installs++
		}
	}
	if installs != 2 {
		t.Errorf("ci.yml: %d job(s) install the pinned solver to /usr/local/bin/z3, want 2 — BOTH jobs "+
			"must install it. Job 1 runs ai-check directly; job 2 runs it THROUGH host/verifygate's "+
			"shim arms, whose `verify gate PASSED` assertion is unreachable without a solver.", installs)
	}
}

func TestCorpusPredicateMatchesShellAnchors(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	if err != nil {
		t.Fatal(err)
	}
	for _, anchor := range []string{`(-dirty$|-[0-9]+-g[0-9a-f]+)`, `^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.]+)?$`} {
		if count := strings.Count(string(src), anchor); count != 1 {
			t.Errorf("shell predicate anchor %q count=%d", anchor, count)
		}
	}
}

func Example_reasonCodes() {
	fmt.Println(releaseVerdict("v0.33.0-dirty"))
	// Output: DEV_BUILD
}
