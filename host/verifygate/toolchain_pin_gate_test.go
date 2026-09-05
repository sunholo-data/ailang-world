package verifygate

import (
	"context"
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"go/version"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
)

func stripPinQuotes(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '\'' && value[len(value)-1] == '\'') ||
		(value[0] == '"' && value[len(value)-1] == '"')) {
		value = value[1 : len(value)-1]
	}
	return value
}

// canonicalizeVersionPin serves ONLY the conventions whose native spelling omits the
// "go" prefix: actions/setup-go `go-version:` values and go.mod `go` directives
// (design doc P8). Prepending is correct here and nowhere else.
func canonicalizeVersionPin(value string) string {
	value = stripPinQuotes(value)
	if value != "" && !strings.HasPrefix(value, "go") {
		value = "go" + value
	}
	return value
}

// requireToolchainNamePin serves the conventions that carry a Go toolchain NAME:
// ci.yml GOTOOLCHAIN pins and run.sh's PINNED=. Per https://go.dev/doc/toolchain
// (where `go help environment` sends GOTOOLCHAIN readers), the setting's grammar is
// <name>, <name>+auto, <name>+path, or the shorthands auto/local/path — but only the
// bare <name> form PINS. The mode words and +suffix forms are VALID to Go — measured:
// GOTOOLCHAIN=go1.26.6+auto is rc=0 and runs go1.26.6 (design doc P19) — so rejecting
// them is a pin-stability POLICY choice, not a validity claim: they are selection
// channels under which the resolved toolchain can move without this file changing.
// The second arm is this REPOSITORY'S PIN POLICY, not a claim of equivalence with
// the runtime's own name grammar (round-2 R1; custom `goV-suffix` names measured in
// design doc P22). Both arms are Fatalf, and neither is an instrument-class floor — the
// instrument is fine; the INPUT is not a pin. Stopping is still correct, because
// every comparison downstream (the agreement loop, the go.mod-floor check) takes
// this value as its operand: grading agreement over a non-pin is the E0
// misattribution reproduced one call later (design doc P21). One defect, one
// attributed message. Nothing is auto-corrected, because a value Go itself refuses
// (`go: invalid GOTOOLCHAIN "1.26.6"`, measured) must never be repaired into a
// value it would accept.
func requireToolchainNamePin(t *testing.T, source, key, raw string) string {
	t.Helper()
	value := stripPinQuotes(raw)
	switch {
	case value == "auto" || value == "local" || value == "path" || strings.Contains(value, "+"):
		t.Fatalf("%s: %s=%q is a toolchain-selection mode, not a pin; only a bare toolchain name (e.g. go1.26.6) pins", source, key, value)
	case !version.IsValid(value):
		t.Fatalf("%s: %s=%q is not an allowed standard Go toolchain pin; this repository requires a bare standard toolchain version accepted by go/version.IsValid (for example go1.26.6)", source, key, value)
	}
	return value
}

func keyedValues(lines []string, key string) []string {
	var values []string
	for _, line := range lines {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) == 2 && parts[0] == key {
			values = append(values, stripPinQuotes(parts[1]))
		}
	}
	return values
}

func moduleGoFloor(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var floors []string
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(line, "go ") {
			floors = append(floors, canonicalizeVersionPin(strings.TrimPrefix(line, "go ")))
		}
	}
	if len(floors) != 1 {
		t.Fatalf("%s: found %d line(s) beginning with %q, want exactly 1", path, len(floors), "go ")
	}
	return floors[0]
}

// TestGoToolchainPinsAgreeAndMatchJobList guards both workflow pin kinds and the module floor.
//
// DECLARED RESIDUAL (mirrors the Z3 precedent at ail_binary_gate_test.go: this is a STATIC
// text scan over YAML. It sees the pin TEXT, never whether the setup-go step RUNS: a
// step-level `if:` whose expression is always false at runtime — e.g. one keyed on a commit
// marker nobody writes — disables the install with every counted byte intact, and no
// actionlint runs anywhere in this repo (measured: 0 hits over *.yml/*.sh, V18), so not even
// a literal `if: false` is flagged. It cannot see what version the runner actually INSTALLED
// (setup-go cache, image drift): the job's `go version` step prints without asserting, and
// verify_go.sh's deny-list observes only the ACTIVE toolchain, which the surviving
// GOTOOLCHAIN pin makes go1.26.6 anyway. And it parses lines, not YAML: a pin folded into a
// flow-style `with: {…}` drops the keyed count (a RED here, not a silent pass), but an exotic
// form that preserved the counts AND smuggled a value past quote-stripping would pass — the
// hand-maintained constants above, not this parser, are the bar.
func TestGoToolchainPinsAgreeAndMatchJobList(t *testing.T) {
	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	lines := strings.Split(src, "\n")
	for _, control := range []string{"ailang-verify:", "go-verify:", "launchd-drivers:", "uses: actions/setup-go@v5", "./scripts/verify_go.sh"} {
		if !strings.Contains(src, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", workflowPath, control)
		}
	}

	// Jobs are enumerated WITH THEIR BODIES, because "how many jobs are there" and
	// "which jobs must carry a Go pin" stopped being the same question on 2026-09-02.
	// The `launchd-drivers` job (bash 3.2 on macos-latest) needs no Go toolchain at
	// all, so counting pins against the JOB count demanded a third pin that must not
	// exist. That conflation is what took dev red on the merge commit 68403ea: two
	// textually non-conflicting branches — one adding the job, one carrying this
	// gate — were each green alone and jointly red, and the gate was RIGHT to red,
	// it just could not say which of the two facts it objected to.
	jobLine := regexp.MustCompile(`^  ([a-z0-9-]+):$`)
	seenJobs := false
	var jobs []string
	jobBody := map[string][]string{}
	current := ""
	for _, line := range lines {
		if !seenJobs {
			if strings.TrimSpace(line) == "jobs:" {
				seenJobs = true
			}
			continue
		}
		if match := jobLine.FindStringSubmatch(line); match != nil {
			current = match[1]
			jobs = append(jobs, current)
			continue
		}
		if current != "" {
			jobBody[current] = append(jobBody[current], line)
		}
	}
	slices.Sort(jobs)
	// TWO hand-maintained sets, and a new job must be classified in BOTH edits or this
	// test reds: wantJobs says the job exists on purpose, wantGoPinnedJobs says whether
	// it is a Go job. A non-Go job added to wantJobs alone is still asserted to carry
	// ZERO Go pins, so "classified out" is a claim the gate checks rather than a hole.
	wantJobs := []string{"ailang-verify", "go-verify", "launchd-drivers"}
	wantGoPinnedJobs := []string{"ailang-verify", "go-verify"}
	if !slices.Equal(jobs, wantJobs) {
		t.Errorf("ci.yml: enumerated jobs=%v, want %v; GOTOOLCHAIN pins=%d go-version pins=%d",
			jobs, wantJobs, len(keyedValues(lines, "GOTOOLCHAIN")), len(keyedValues(lines, "go-version")))
	}
	for _, job := range wantGoPinnedJobs {
		if !slices.Contains(wantJobs, job) {
			t.Fatalf("instrument failure: Go-pinned job %q is not in the enumerated job set %v", job, wantJobs)
		}
	}

	// PER-JOB ATTRIBUTION, not a repo-wide count. A count over the whole file cannot
	// tell WHICH job carries a pin, so two GOTOOLCHAIN lines in one job and none in
	// the other satisfied the old arms exactly as well as one each — a fail-open the
	// gate has always had and which the job-kind split would otherwise have widened.
	pinKinds := []string{"GOTOOLCHAIN", "go-version", "actions/setup-go"}
	perJob := map[string][3]int{}
	var summed [3]int
	for _, job := range jobs {
		body := jobBody[job]
		counts := [3]int{len(keyedValues(body, "GOTOOLCHAIN")), len(keyedValues(body, "go-version")), 0}
		for _, line := range body {
			if strings.Contains(line, "uses: actions/setup-go@") {
				counts[2]++
			}
		}
		perJob[job] = counts
		for i := range summed {
			summed[i] += counts[i]
		}
	}
	if len(jobBody) == 0 {
		t.Fatalf("instrument failure: parsed %d job name(s) from %s but zero job bodies", len(jobs), workflowPath)
	}
	for _, job := range jobs {
		want := 0
		if slices.Contains(wantGoPinnedJobs, job) {
			want = 1
		}
		counts := perJob[job]
		for i, kind := range pinKinds {
			if counts[i] != want {
				t.Errorf("ci.yml: job %q carries %d %s line(s), want %d (Go-pinned jobs=%v)", job, counts[i], kind, want, wantGoPinnedJobs)
			}
		}
	}

	goToolchains := keyedValues(lines, "GOTOOLCHAIN")
	for i, raw := range goToolchains {
		goToolchains[i] = requireToolchainNamePin(t, "ci.yml", "GOTOOLCHAIN", raw)
	}
	goVersions := keyedValues(lines, "go-version")
	for i, raw := range goVersions {
		goVersions[i] = canonicalizeVersionPin(raw)
	}
	setupGoUses := strings.Count(src, "uses: actions/setup-go@")
	// Whole-file counts must equal the per-job sums, or a pin sits OUTSIDE every
	// enumerated job — a workflow-level `env: GOTOOLCHAIN`, or a job this parser
	// failed to attribute. Either way the per-job arms above were reading a set
	// smaller than the file, which is the one way they can be vacuously satisfied.
	for i, total := range []int{len(goToolchains), len(goVersions), setupGoUses} {
		if total != summed[i] {
			t.Errorf("ci.yml: %s whole-file count=%d but per-job sum=%d; a pin sits outside every enumerated job", pinKinds[i], total, summed[i])
		}
	}
	if len(goToolchains) != len(wantGoPinnedJobs) {
		t.Errorf("ci.yml: GOTOOLCHAIN keyed-line count=%d, want %d (one per Go-pinned job)", len(goToolchains), len(wantGoPinnedJobs))
	}
	if len(goVersions) != len(wantGoPinnedJobs) {
		t.Errorf("ci.yml: go-version keyed-line count=%d, want %d (one per Go-pinned job)", len(goVersions), len(wantGoPinnedJobs))
	}
	if setupGoUses != len(wantGoPinnedJobs) {
		t.Errorf("ci.yml: actions/setup-go use count=%d, want %d (one per Go-pinned job)", setupGoUses, len(wantGoPinnedJobs))
	}
	allPins := append(append([]string{}, goToolchains...), goVersions...)
	if len(allPins) > 0 {
		for _, pin := range allPins[1:] {
			if pin != allPins[0] {
				t.Errorf("ci.yml: toolchain pins disagree: GOTOOLCHAIN=%v go-version=%v", goToolchains, goVersions)
				break
			}
		}
	}

	goVersionNeedleCount := strings.Count(src, "go-version:")
	if goVersionNeedleCount < 2 {
		t.Fatalf("instrument failure: ci.yml count(%q)=%d, want at least 2 before checking the alternative-form zero", "go-version:", goVersionNeedleCount)
	}
	if count := strings.Count(src, "go-version-file"); count != 0 {
		t.Errorf("ci.yml: count(%q)=%d, want 0; alternative setup-go pin form bypasses keyed extraction", "go-version-file", count)
	}

	modulePath := filepath.Join(repoRoot, "go.mod")
	moduleRaw, err := os.ReadFile(modulePath)
	if err != nil {
		t.Fatal(err)
	}
	toolchainDirectives := 0
	for _, line := range strings.Split(string(moduleRaw), "\n") {
		if strings.HasPrefix(line, "toolchain ") {
			toolchainDirectives++
		}
	}
	if toolchainDirectives != 0 {
		t.Errorf("go.mod: toolchain directive count=%d, want 0; it is a hidden floor override", toolchainDirectives)
	}
	floor := moduleGoFloor(t, modulePath)
	if len(allPins) > 0 && floor != allPins[0] {
		t.Errorf("go.mod floor=%q disagrees with ci.yml toolchain pin=%q", floor, allPins[0])
	}

	matches, err := filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*"))
	if err != nil {
		t.Fatal(err)
	}
	workflowFiles := make([]string, 0, len(matches))
	for _, match := range matches {
		workflowFiles = append(workflowFiles, filepath.Base(match))
	}
	slices.Sort(workflowFiles)
	if !slices.Equal(workflowFiles, []string{"ci.yml"}) {
		t.Errorf("workflow files=%v, want exactly [ci.yml]; a second workflow may carry unscanned toolchain pins", workflowFiles)
	}
}

func shellAssignmentValues(lines []string, name string) []string {
	prefix := name + "=\""
	var values []string
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			rest := strings.TrimPrefix(line, prefix)
			if end := strings.IndexByte(rest, '"'); end >= 0 {
				values = append(values, rest[:end])
			}
		}
	}
	return values
}

// TestMiscompileInstrumentProbesPinnedToolchain binds the reproducer to the module floor.
//
// DECLARED RESIDUAL: this is a STATIC text scan. It proves the LIST and PINNED carry the
// floor token and that the pinned-guard machinery EXISTS as text (the saw_pinned_ok sites
// and the failure message). The guard's firing is proven at sprint time by AC5's guard-trip
// and exercised on every gated CI invocation of run.sh. The round-1 SKIPPED hole (a banner
// printed with the pin unprobed) is closed in run.sh itself, not merely narrowed here; the
// row-44 wiring test below binds that a loud refusal remains gating.
func TestMiscompileInstrumentProbesPinnedToolchain(t *testing.T) {
	scriptPath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "run.sh")
	raw, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	lines := strings.Split(src, "\n")
	fixturePath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "toolchain_pins.conf")
	fixtureRaw, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	fixtureSrc := string(fixtureRaw)
	fixtureLines := strings.Split(fixtureSrc, "\n")
	shebangs := 0
	for _, line := range lines {
		if line == "#!/usr/bin/env bash" {
			shebangs++
		}
	}
	if shebangs != 1 {
		t.Fatalf("instrument failure: %s exact shebang count=%d, want 1", scriptPath, shebangs)
	}

	goodAssignments := shellAssignmentValues(fixtureLines, "KNOWN_GOOD")
	badAssignments := shellAssignmentValues(fixtureLines, "KNOWN_BAD")
	pinnedAssignments := shellAssignmentValues(fixtureLines, "PINNED")
	if len(goodAssignments) != 1 {
		t.Errorf("%s: KNOWN_GOOD assignment count=%d, want 1", fixturePath, len(goodAssignments))
	}
	if len(badAssignments) != 1 {
		t.Errorf("%s: KNOWN_BAD assignment count=%d, want 1", fixturePath, len(badAssignments))
	}
	if len(pinnedAssignments) != 1 {
		t.Errorf("%s: PINNED assignment count=%d, want 1", fixturePath, len(pinnedAssignments))
	}
	for _, control := range []string{"KNOWN_BAD=", "KNOWN_GOOD=", "PINNED="} {
		if !strings.Contains(fixtureSrc, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", fixturePath, control)
		}
	}
	var good, bad []string
	if len(goodAssignments) == 1 {
		good = strings.Fields(goodAssignments[0])
	}
	if len(badAssignments) == 1 {
		bad = strings.Fields(badAssignments[0])
	}
	if len(good) == 0 {
		t.Errorf("%s: KNOWN_GOOD must contain at least one toolchain", scriptPath)
	}
	if len(bad) == 0 {
		t.Errorf("%s: KNOWN_BAD must contain at least one toolchain", fixturePath)
	}

	floor := moduleGoFloor(t, filepath.Join(repoRoot, "go.mod"))
	if !slices.Contains(good, floor) {
		t.Errorf("KNOWN_GOOD=%v does not probe the pinned toolchain %s from go.mod", good, floor)
	}
	pinned := ""
	if len(pinnedAssignments) == 1 {
		pinned = requireToolchainNamePin(t, fixturePath, "PINNED", pinnedAssignments[0])
	}
	if pinned != floor {
		t.Errorf("PINNED=%q, want go.mod floor %q", pinned, floor)
	}
	if pinned != "" && !slices.Contains(good, pinned) {
		t.Errorf("PINNED=%q is absent from KNOWN_GOOD=%v, so the probe loop cannot set its OK flag", pinned, good)
	}
	if slices.Contains(bad, floor) {
		t.Errorf("KNOWN_BAD=%v incorrectly labels the pinned toolchain %s as affected", bad, floor)
	}

	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Errorf("%s is not executable; CI invokes it directly", scriptPath)
	}
	if count := strings.Count(src, "saw_pinned_ok"); count < 3 {
		t.Errorf("%s: saw_pinned_ok site count=%d, want at least 3 (declaration, OK set, guard)", scriptPath, count)
	}
	if !strings.Contains(src, "INSTRUMENT FAILURE: the PINNED toolchain") {
		t.Errorf("%s: pinned-toolchain fail-loud guard message is absent", scriptPath)
	}

	// Direction pin (row 45, finding C): the site COUNT above cannot see the guard's
	// comparison direction — flipping `-eq 0` to `-ne 0` preserves all 3 sites and
	// every committed assertion (controller, iteration 134). Count comment-stripped
	// CODE lines only: the guard and the prose explaining it must not compete for one
	// namespace (iteration-133 lesson) — a comment may quote the literal freely.
	guardLines := 0
	for _, line := range lines {
		code := line
		if idx := strings.Index(code, "#"); idx >= 0 {
			code = code[:idx]
		}
		if strings.Contains(code, `[ "$saw_pinned_ok" -eq 0 ]`) {
			guardLines++
		}
	}
	if guardLines != 1 {
		t.Errorf("%s: executable pinned-OK guard-line count=%d, want exactly 1 — the floor must test ABSENCE of the OK flag (`-eq 0`); a flipped or duplicated guard is a different instrument", scriptPath, guardLines)
	}
}

func TestReproModuleFloorStaysBelowKnownBadToolchains(t *testing.T) {
	reproFloor := moduleGoFloor(t, filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "repro", "go.mod"))
	if !version.IsValid(reproFloor) {
		t.Fatalf("instrument failure: repro module floor %q is not a valid Go version", reproFloor)
	}

	fixturePath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "toolchain_pins.conf")
	raw, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	badAssignments := shellAssignmentValues(strings.Split(string(raw), "\n"), "KNOWN_BAD")
	if len(badAssignments) != 1 {
		t.Fatalf("instrument failure: %s: KNOWN_BAD assignment count=%d, want 1", fixturePath, len(badAssignments))
	}
	bad := strings.Fields(badAssignments[0])
	if len(bad) == 0 {
		t.Fatalf("instrument failure: %s: KNOWN_BAD must contain at least one toolchain", fixturePath)
	}
	oldest := bad[0]
	for _, tc := range bad {
		if !version.IsValid(tc) {
			t.Fatalf("instrument failure: KNOWN_BAD token %q is not a valid Go version; version.Compare would misorder it", tc)
		}
		if version.Compare(tc, oldest) < 0 {
			oldest = tc
		}
	}
	if version.Compare(reproFloor, oldest) > 0 {
		t.Fatalf("repro module floor %q is above the oldest KNOWN_BAD toolchain %q: every deny-listed probe SKIPs, saw_bad stays 0, and run.sh reds for the wrong reason (the V10 rehearsal)", reproFloor, oldest)
	}
}

const toolchainPinParserEnd = "# toolchain_pins.conf parser ends here."
const sentinelToolchainFixture = "KNOWN_BAD=\"sentinel-bad\"\nKNOWN_GOOD=\"sentinel-good\"\nPINNED=\"sentinel-pinned\"\n"

func toolchainPinParserPrologue(t *testing.T) string {
	t.Helper()
	scriptPath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "run.sh")
	raw, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(raw), "\n")
	end := -1
	for i, line := range lines {
		if line == toolchainPinParserEnd {
			if end != -1 {
				t.Fatalf("instrument failure: %s contains the parser-end marker more than once", scriptPath)
			}
			end = i
		}
	}
	if end == -1 {
		t.Fatalf("instrument failure: could not locate %q in %s", toolchainPinParserEnd, scriptPath)
	}
	return strings.Join(lines[:end+1], "\n") + "\n"
}

func requireBash(t *testing.T) string {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Fatalf("bash is required to execute the toolchain pin fixture: %v", err)
	}
	return bash
}

func writeExecutableAt(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatal(err)
	}
}

func runBoundedBash(t *testing.T, bash, script string, env []string) ([]byte, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bash, script)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("run.sh did not execute toolchain_pins.conf: %v", ctx.Err())
	}
	return out, err
}

func writeToolchainFixture(t *testing.T, dir, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "toolchain_pins.conf"), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(raw))
}

func requireChildExitCode(t *testing.T, err error) int {
	t.Helper()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("scratch parser did not report an exit status: %v", err)
	}
	return exitErr.ExitCode()
}

func assertWritableSentinel(t *testing.T, sentinel string) {
	t.Helper()
	file, err := os.Create(sentinel)
	if err != nil {
		t.Fatalf("known-positive writable-directory control could not create %s: %v", sentinel, err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("known-positive writable-directory control did not create %s: %v", sentinel, err)
	}
	if err := os.Remove(sentinel); err != nil {
		t.Fatal(err)
	}
}

// TestToolchainPinFixtureIsDataOnly makes the fixture's column-zero scan complete
// by construction. The bash -n arm is instrument health for the bounded parser;
// the fixture mutations are the load-bearing proof of the data-only claim.
func TestToolchainPinFixtureIsDataOnly(t *testing.T) {
	fixturePath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "toolchain_pins.conf")
	raw, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read toolchain pin fixture: %v", err)
	}
	record := regexp.MustCompile(`^([A-Z_]+)="([^"]*)"([[:space:]]+#.*)?$`)
	counts := map[string]int{}
	records := 0
	for lineNumber, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		match := record.FindStringSubmatch(line)
		if match == nil {
			t.Errorf("%s:%d line %q does not match the anchored assignment grammar", fixturePath, lineNumber+1, line)
			continue
		}
		records++
		counts[match[1]]++
	}
	if records == 0 {
		t.Fatalf("instrument failure: %s contains zero assignment records", fixturePath)
	}
	for _, name := range []string{"KNOWN_BAD", "KNOWN_GOOD", "PINNED"} {
		if counts[name] != 1 {
			t.Errorf("%s: %s assignment count=%d, want exactly 1", fixturePath, name, counts[name])
		}
	}
	if len(counts) != 3 {
		t.Errorf("%s: assignment names=%v, want exactly KNOWN_BAD, KNOWN_GOOD, PINNED", fixturePath, counts)
	}

	scriptPath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "run.sh")
	scriptRaw, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("read run.sh: %v", err)
	}
	const fixtureReference = `conf="$(dirname "$0")/toolchain_pins.conf"`
	if count := strings.Count(string(scriptRaw), fixtureReference); count != 1 {
		t.Errorf("%s fixture-reference count=%d, want exactly 1", scriptPath, count)
	}

	bash := requireBash(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bash, "-n", scriptPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("instrument failure: bash could not parse %s: %v: %s", scriptPath, err, output)
	}
	if ctx.Err() != nil {
		t.Fatalf("instrument failure: bash syntax check timed out for %s: %v", scriptPath, ctx.Err())
	}
}

// TestRunShExecutesToolchainPinFixture executes only a scratch copy of run.sh's
// prologue through the bounded parser. Bash syntax checks are instrument health;
// the observed sentinel values and rejection/non-execution arms are the runtime proof.
func TestRunShExecutesToolchainPinFixture(t *testing.T) {
	bash := requireBash(t)
	prologue := toolchainPinParserPrologue(t)

	t.Run("loads sentinel values", func(t *testing.T) {
		dir := t.TempDir()
		script := filepath.Join(dir, "run.sh")
		writeExecutableAt(t, script, prologue+`printf '%s|%s|%s\n' "$KNOWN_BAD" "$KNOWN_GOOD" "$PINNED"`+"\n")
		writeToolchainFixture(t, dir, sentinelToolchainFixture)

		if syntax, err := exec.Command(bash, "-n", script).CombinedOutput(); err != nil {
			t.Fatalf("run.sh did not execute toolchain_pins.conf: scratch syntax check failed: %v: %s", err, syntax)
		}
		out, err := runBoundedBash(t, bash, script, os.Environ())
		if err != nil {
			t.Fatalf("run.sh did not execute toolchain_pins.conf: %v: %s", err, out)
		}
		if got, want := strings.TrimSpace(string(out)), "sentinel-bad|sentinel-good|sentinel-pinned"; got != want {
			t.Fatalf("run.sh did not execute toolchain_pins.conf: output=%q, want %q", got, want)
		}
	})

	for _, tc := range []struct {
		name     string
		sentinel string
		badLine  func(string) string
	}{
		{"rejects command substitution without execution", "sentinel_m11", func(path string) string { return `KNOWN_BAD="$(touch ` + path + `)"` }},
		{"rejects backticks without execution", "sentinel_m12", func(path string) string { return "KNOWN_BAD=\"`touch " + path + "`\"" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			sentinel := filepath.Join(dir, tc.sentinel)
			assertWritableSentinel(t, sentinel)
			script := filepath.Join(dir, "run.sh")
			writeExecutableAt(t, script, prologue)
			fixture := filepath.Join(dir, "toolchain_pins.conf")
			writeToolchainFixture(t, dir, sentinelToolchainFixture)
			before := fileSHA256(t, fixture)
			writeToolchainFixture(t, dir, tc.badLine(sentinel)+"\nKNOWN_GOOD=\"sentinel-good\"\nPINNED=\"sentinel-pinned\"\n")
			mutant := fileSHA256(t, fixture)
			if mutant == before {
				t.Fatal("mutation did not change the scratch fixture")
			}

			out, err := runBoundedBash(t, bash, script, os.Environ())
			if err == nil {
				t.Fatalf("run.sh accepted executable fixture value; output=%q", out)
			}
			childRC := requireChildExitCode(t, err)
			const rejection = "value for 'KNOWN_BAD' contains a disallowed character"
			if !strings.Contains(string(out), rejection) {
				t.Fatalf("run.sh rejection=%q, want substring %q", out, rejection)
			}
			if _, statErr := os.Stat(sentinel); !os.IsNotExist(statErr) {
				t.Fatalf("fixture value executed: sentinel %s stat error=%v, want os.IsNotExist", sentinel, statErr)
			}
			writeToolchainFixture(t, dir, sentinelToolchainFixture)
			after := fileSHA256(t, fixture)
			if after != before {
				t.Fatalf("scratch fixture restore sha256=%s, want %s", after, before)
			}
			t.Logf("observed parser rejection: %s; child_rc=%d sha256 before=%s mutant=%s after=%s", rejection, childRC, before, mutant, after)
		})
	}

	t.Run("rejects PATH without clobbering it", func(t *testing.T) {
		dir := t.TempDir()
		pathOut := filepath.Join(dir, "path.out")
		const childPath = "/usr/bin:/bin:/known-child-path"
		observer := `trap 'printf "%s" "$PATH" > "$PATHOUT"' EXIT` + "\n"
		observedPrologue := strings.Replace(prologue, "set -uo pipefail\n", "set -uo pipefail\n"+observer, 1)
		if observedPrologue == prologue {
			t.Fatal("instrument failure: could not install PATH observer before the parser")
		}
		script := filepath.Join(dir, "run.sh")
		writeExecutableAt(t, script, observedPrologue)
		fixture := filepath.Join(dir, "toolchain_pins.conf")
		writeToolchainFixture(t, dir, sentinelToolchainFixture)
		before := fileSHA256(t, fixture)
		writeToolchainFixture(t, dir, "KNOWN_BAD=\"sentinel-bad\"\nPATH=\"/nonsense\"\nKNOWN_GOOD=\"sentinel-good\"\nPINNED=\"sentinel-pinned\"\n")
		mutant := fileSHA256(t, fixture)
		if mutant == before {
			t.Fatal("mutation did not change the scratch fixture")
		}
		env := append(os.Environ(), "PATH="+childPath, "PATHOUT="+pathOut)
		out, err := runBoundedBash(t, bash, script, env)
		if err == nil {
			t.Fatalf("run.sh accepted unknown fixture name; output=%q", out)
		}
		childRC := requireChildExitCode(t, err)
		const rejection = "unknown name 'PATH' (only KNOWN_BAD, KNOWN_GOOD, PINNED allowed)"
		if !strings.Contains(string(out), rejection) {
			t.Fatalf("run.sh rejection=%q, want substring %q", out, rejection)
		}
		recorded, readErr := os.ReadFile(pathOut)
		if readErr != nil {
			t.Fatalf("PATH observer did not record the refusing shell's PATH: %v", readErr)
		}
		if got := string(recorded); got != childPath {
			t.Fatalf("parser clobbered PATH before refusal: got %q, want %q", got, childPath)
		}

		controlOut := filepath.Join(dir, "path-control.out")
		control := filepath.Join(dir, "path-control.sh")
		writeExecutableAt(t, control, "#!/usr/bin/env bash\nset -uo pipefail\nPATH=/nonsense\ntrap 'printf \"%s\" \"$PATH\" > \"$PATHOUT\"' EXIT\n")
		if output, controlErr := runBoundedBash(t, bash, control, append(os.Environ(), "PATHOUT="+controlOut)); controlErr != nil {
			t.Fatalf("known-positive PATH observer control failed: %v: %s", controlErr, output)
		}
		controlRecorded, readErr := os.ReadFile(controlOut)
		if readErr != nil {
			t.Fatalf("known-positive PATH observer control wrote no observation: %v", readErr)
		}
		if got := string(controlRecorded); got != "/nonsense" {
			t.Fatalf("known-positive PATH observer control recorded %q, want /nonsense", got)
		}
		writeToolchainFixture(t, dir, sentinelToolchainFixture)
		after := fileSHA256(t, fixture)
		if after != before {
			t.Fatalf("scratch fixture restore sha256=%s, want %s", after, before)
		}
		t.Logf("observed parser rejection: %s; child_rc=%d sha256 before=%s mutant=%s after=%s", rejection, childRC, before, mutant, after)
	})
}

// canaryAssertionShapeProblems parses the canary test source and reports any deviation from the
// required assertion shape: TestToolchainCanary must exist exactly once and contain exactly one
// top-level `if` whose condition is rows[0].field != "stateRoot". That if-body must contain
// exactly one direct t.Fatalf expression statement, rather than merely a descendant call.
func canaryAssertionShapeProblems(src string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return []string{fmt.Sprintf("parse error: %v", err)}
	}

	var funcs []*ast.FuncDecl
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "TestToolchainCanary" {
			funcs = append(funcs, fn)
		}
	}
	if len(funcs) != 1 {
		return []string{fmt.Sprintf("TestToolchainCanary func decl count=%d, want exactly 1", len(funcs))}
	}

	var assertions []*ast.IfStmt
	for _, stmt := range funcs[0].Body.List {
		ifStmt, ok := stmt.(*ast.IfStmt)
		if !ok {
			continue
		}
		binary, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if ok && binary.Op == token.NEQ && isRowsField(binary.X) && isStateRootLit(binary.Y) {
			assertions = append(assertions, ifStmt)
		}
	}
	if len(assertions) != 1 {
		return []string{fmt.Sprintf("top-level `rows[0].field != \"stateRoot\"` assertion if-stmt count=%d, want exactly 1", len(assertions))}
	}

	directFatalf := 0
	for _, stmt := range assertions[0].Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if ok && isTFatalfCall(call) {
			directFatalf++
		}
	}
	if directFatalf != 1 {
		return []string{fmt.Sprintf("direct t.Fatalf expression statement in assertion body count=%d, want exactly 1", directFatalf)}
	}
	// A5 — reachability: no t.Skip / t.Skipf / t.SkipNow call anywhere in the body, INCLUDING
	// inside nested func literals. A t.Skip on the outer `t` inside a closure still Goexits
	// TestToolchainCanary (V30: the assertion after it never runs), so descending is not merely
	// conservative — it is correct. The selector set is CLOSED by construction against
	// `go doc testing.T`'s three skip methods (Skip, Skipf, SkipNow) — see V27.
	var problems []string
	skipCalls := 0
	ast.Inspect(funcs[0].Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "t" &&
					(sel.Sel.Name == "Skip" || sel.Sel.Name == "Skipf" || sel.Sel.Name == "SkipNow") {
					skipCalls++
				}
			}
		}
		return true
	})
	if skipCalls != 0 {
		problems = append(problems, fmt.Sprintf("t.Skip/t.Skipf/t.SkipNow call count=%d, want 0 (a skipped canary asserts nothing)", skipCalls))
	}

	// A6 — reachability: no early return in the body. The assertion is the last statement,
	// so any return before it neuters it. Traversal STOPS at *ast.FuncLit: a return inside a
	// nested func literal exits the literal, not TestToolchainCanary (V30), so counting it
	// would false-red a canary whose assertion still runs (V31, the row-55 class).
	returns := 0
	ast.Inspect(funcs[0].Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if _, ok := n.(*ast.ReturnStmt); ok {
			returns++
		}
		return true
	})
	if returns != 0 {
		problems = append(problems, fmt.Sprintf("return statement count=%d, want 0 (an early return neuters the assertion)", returns))
	}

	// A8 — reachability: no goto in the body. A `goto` over the assertion is dead-code
	// elision that leaves every shape check byte-identical AND leaves the canary reporting
	// `--- PASS`, so it is strictly QUIETER than a skip, which at least prints `--- SKIP`
	// (V35: mutant sha256 f7d6d640257d9f61, canary `--- PASS`, extended fence rc=0 before
	// this check existed). It is statically visible — `go vet` flags it `unreachable code`
	// — but `scripts/verify_go.sh` does not run vet, so nothing in the enforced gate caught
	// it. Traversal STOPS at *ast.FuncLit, like A6 and for the same reason: Go forbids a
	// goto crossing a function boundary, so a goto inside a closure cannot jump past the
	// outer assertion, and counting it would be a row-55-class false red.
	gotos := 0
	ast.Inspect(funcs[0].Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if b, ok := n.(*ast.BranchStmt); ok && b.Tok == token.GOTO {
			gotos++
		}
		return true
	})
	if gotos != 0 {
		problems = append(problems, fmt.Sprintf("goto statement count=%d, want 0 (a goto can jump past the assertion)", gotos))
	}

	// A7 — reachability: no build constraint on the file. A build tag can exclude the
	// canary from the build entirely, so the assertion never runs.
	//
	// ONE NARROWING, proven load-bearing by its own differ-test. It exists to kill a
	// row-55-class FALSE RED, not to weaken the check — a build tag Go would actually honour
	// still reds, in both the modern (M-BUILDTAG) and legacy (M-BUILDTAG-PLUS) forms.
	//
	// GRAMMAR — match with go/build/constraint, Go's OWN parser, never a byte prefix.
	// `strings.HasPrefix(c.Text, "// +build")` reds on `// +buildAlerts is a codename`, which
	// Go reads as ordinary prose: IsGoBuild and IsPlusBuild are both false, `go vet` is silent,
	// and the file is NOT excluded from the build — yet the prefix form redded the fence over a
	// canary that runs and passes (V36: mutant sha256 50f83fc840f7f68f, canary `--- PASS`,
	// prefix-A7 rc=1). `//go:buildFoo bar` is the same shape one constraint over. A prefix match
	// is strictly cruder than the grammar Go itself applies, and the gap is a false red.
	//
	// A POSITION narrowing was written, measured and DELETED rather than shipped (V37). The
	// theory was that only a comment group closing before the package clause can be a
	// constraint, so a grammar-valid `//go:build` quoted later in the file should not red.
	// It is unpinnable: Go REJECTS a misplaced `//go:build` outright — `go vet` rc=1
	// `misplaced //go:build comment` and the compile fence rc=1 — both after the package
	// clause and inside a function body. So the guard can only change the verdict on a tree
	// that does not build, and a gate's reading on an unbuildable tree is not a verdict. Its
	// differ-test "passes" for exactly that vacuous reason. Shipping it would have added a
	// branch no arm can ever reach: the anti-vacuity-floor class this charter tracks.
	buildTags := 0
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if constraint.IsGoBuild(c.Text) || constraint.IsPlusBuild(c.Text) {
				buildTags++
			}
		}
	}
	if buildTags != 0 {
		problems = append(problems, fmt.Sprintf("build constraint count=%d, want 0 (a build tag can exclude the canary from the build)", buildTags))
	}

	return problems
}

// isRowsField reports whether expr is structurally rows[0].field.
func isRowsField(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "field" {
		return false
	}
	index, ok := selector.X.(*ast.IndexExpr)
	if !ok {
		return false
	}
	rows, ok := index.X.(*ast.Ident)
	if !ok || rows.Name != "rows" {
		return false
	}
	literal, ok := index.Index.(*ast.BasicLit)
	return ok && literal.Kind == token.INT && literal.Value == "0"
}

func isStateRootLit(expr ast.Expr) bool {
	literal, ok := expr.(*ast.BasicLit)
	return ok && literal.Kind == token.STRING && literal.Value == `"stateRoot"`
}

func isTFatalfCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Fatalf" {
		return false
	}
	receiver, ok := selector.X.(*ast.Ident)
	return ok && receiver.Name == "t"
}

func TestCanaryDeclaresPositiveArmOnly(t *testing.T) {
	canaryPath := filepath.Join(repoRoot, "host", "store", "toolchain_canary_test.go")
	raw, err := os.ReadFile(canaryPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if problems := canaryAssertionShapeProblems(src); len(problems) > 0 {
		for _, problem := range problems {
			t.Errorf("%s: %s", canaryPath, problem)
		}
		t.Fatalf("instrument failure: %s no longer asserts the miscompile shape", canaryPath)
	}
	if count := strings.Count(src, "GOTOOLCHAIN"); count != 0 {
		t.Errorf("%s count(%q)=%d, want 0; known-bad arms belong in the nested repro module", canaryPath, "GOTOOLCHAIN", count)
	}
	if !strings.Contains(src, "POSITIVE ARM ONLY") {
		t.Errorf("%s: required %q marker is absent", canaryPath, "POSITIVE ARM ONLY")
	}
}

const miscompileReproducerPath = "design_docs/verification/w-race-gate-blindspot/run.sh"
const miscompileStepName = "Measure compiler reproducer (platform-conditional, gated)"
const expectedStepCol = 6

func indentOf(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }

// stripYAMLInlineComment removes a trailing `# …` comment from a scalar value. A value that
// opens with a quote is read to its closing quote first, so a `#` inside the quotes stays data.
func stripYAMLInlineComment(v string) string {
	v = strings.TrimLeft(v, " ")
	if len(v) > 0 && (v[0] == '"' || v[0] == '\'') {
		if k := strings.IndexByte(v[1:], v[0]); k >= 0 {
			return v[:k+2]
		}
		return v
	}
	if k := strings.Index(v, " #"); k >= 0 {
		return v[:k]
	}
	if strings.HasPrefix(v, "#") {
		return ""
	}
	return v
}

// stepMappingValue returns the value of `key` when `line` is one of the step's OWN block-mapping
// keys, and reports whether it matched. A step's keys are exactly the ones at stepCol+2: either
// riding the block-sequence dash (`<stepCol>- key: v`, which is where a step's FIRST key sits and
// which an indentation rule written only for sibling lines misses) or on a sibling line
// (`<stepCol+2>key: v`). Anything deeper belongs to a nested mapping or to a block scalar — a
// `run:` script — and is therefore NOT the step's key, which is the whole of queue row 62(b).
func stepMappingValue(line string, stepCol int, key string) (string, bool) {
	ind := indentOf(line)
	rest := line[ind:]
	switch {
	case ind == stepCol && strings.HasPrefix(rest, "- "):
		rest = rest[2:]
	case ind == stepCol+2:
	default:
		return "", false
	}
	if !strings.HasPrefix(rest, key+":") {
		return "", false
	}
	return strings.TrimSpace(stripYAMLInlineComment(rest[len(key)+1:])), true
}

// continueOnErrorIsCompliant reports whether an EXPLICIT continue-on-error value is a static
// opt-OUT. Only a literal false is: `true`, an empty value and a `${{ … }}` expression are all
// refusals, the last because its run-time value is not decidable here and may be true. The scan
// therefore fails CLOSED on any value it cannot read as false — queue row 62 is about false
// POSITIVES, and widening it into a fail-open would trade this instrument for its opposite.
func continueOnErrorIsCompliant(v string) bool {
	return strings.EqualFold(strings.Trim(strings.TrimSpace(v), `"'`), "false")
}

// continueOnErrorRefusalsIn scans the step block [start,end) for the step's own
// `continue-on-error` key and returns one refusal message per non-compliant setting. Queue row
// 62: the predecessor was `strings.Contains(line, "continue-on-error")` over the same range, so
// an explicit `continue-on-error: false` (a legitimate opt-OUT that swallows nothing) and a
// comment or `echo` merely NAMING the flag inside the step's own `run:` scalar both read as
// re-introductions — measured first-party at ci.yml:175 and ci.yml:177 before this fix, with the
// identical `echo` in an unrelated step staying green as the discriminating control.
//
// The neighbouring run.sh check one screen below already draws exactly this line, rejecting only
// EXECUTABLE uses of `go env GOOS` while letting a comment name the channel it warns about. This
// is that rule applied to the call site it missed.
//
// A second return value carries an instrument failure: the scan reads BLOCK mappings, so a step
// written as a YAML flow mapping (`- {continue-on-error: true, …}`) is refused loudly rather
// than passed over silently.
func continueOnErrorRefusalsIn(lines []string, start, end, stepCol int) (refusals []string, instrumentErr string) {
	for i := start; i < end && i < len(lines); i++ {
		if indentOf(lines[i]) == stepCol && strings.HasPrefix(strings.TrimSpace(lines[i]), "- {") {
			return nil, fmt.Sprintf("ci.yml:%d expresses the step as a YAML flow mapping; this scan reads block mappings only", i+1)
		}
		v, ok := stepMappingValue(lines[i], stepCol, "continue-on-error")
		if !ok {
			continue
		}
		if continueOnErrorIsCompliant(v) {
			continue
		}
		refusals = append(refusals, fmt.Sprintf("ci.yml:%d sets %q to %q in the miscompile step — row 44: a swallowed refusal is how this instrument died the first time; only a literal false is an accepted opt-out", i+1, "continue-on-error", v))
	}
	return refusals, ""
}

// TestMiscompileInstrumentStepIsGatedInCI pins the row-44 wiring on two channels that
// must not silently return. (1) `continue-on-error: true` converts an instrument's
// loudest possible output into silence, so it is forbidden in the miscompile step's
// own block. A flag on an unrelated step remains that step's business. (2) The
// instrument's platform polarity reads the KERNEL (`uname`); `go env` honours the
// env-var form of the platform tokens (measured in the design doc, P16), so executable
// uses of that overridable channel are forbidden in run.sh, and both kernel reads are
// asserted present so the probe cannot quietly revert.
// The step scope is an indentation-aware line scan anchored at the shallowest
// enclosing `steps:` key, as ratified by D-WORLD-30. DECLARED RESIDUAL: the scan is
// text-level and unbounded to the enclosing job; a future re-indent therefore fails
// loudly until expectedStepCol is intentionally updated (this replaces the row-44
// design doc's V19 scoping assumption). A step-level `if:` that never evaluates true
// also disables the step with this text intact (no actionlint runs in this repo — P41
// V18); and these are byte-substring pins — a computed assignment (`eval`, string
// concatenation) evades them; the mechanism's runtime immunity is why that gap is
// acceptable (design doc, residuals 2 and 3).
func TestMiscompileInstrumentStepIsGatedInCI(t *testing.T) {
	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if count := strings.Count(src, miscompileReproducerPath); count != 1 {
		t.Fatalf("instrument failure: ci.yml count(%q)=%d, want exactly 1 — the step this test pins must exist", miscompileReproducerPath, count)
	}
	// Scope the flag check to the miscompile STEP's own block (quorum round-2 R1:
	// "inspect only the miscompile step for continue-on-error … rather than banning
	// legitimate GOOS/GOARCH use globally"). A flag on an unrelated step is that
	// step's business; V19's scoping control proves this boundary holds.
	lines := strings.Split(src, "\n")
	identifyingLine := -1
	for i, l := range lines {
		if strings.Contains(l, miscompileReproducerPath) {
			identifyingLine = i
			break
		}
	}
	anchor := -1
	anchorCol := len(lines) + 1
	for j := identifyingLine; j >= 0; j-- {
		if strings.TrimSpace(lines[j]) == "steps:" && indentOf(lines[j]) < anchorCol {
			anchor = j
			anchorCol = indentOf(lines[j])
		}
	}
	if anchor < 0 {
		t.Fatalf("instrument failure: could not locate a steps: anchor above the miscompile identifying line in ci.yml")
	}
	stepCol := -1
	for j := anchor + 1; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "- ") {
			stepCol = indentOf(lines[j])
			break
		}
	}
	if stepCol < 0 {
		t.Fatalf("instrument failure: could not derive the step column below ci.yml:%d", anchor+1)
	}
	if stepCol != expectedStepCol {
		t.Fatalf("instrument failure: derived step column %d; update expectedStepCol after an intentional ci.yml re-indent", stepCol)
	}
	start := -1
	for j := identifyingLine; j >= 0; j-- {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "- ") && indentOf(lines[j]) == stepCol {
			start = j
			break
		}
	}
	if start < 0 {
		t.Fatalf("instrument failure: could not locate the miscompile step block in ci.yml")
	}
	end := len(lines)
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if (strings.HasPrefix(trimmed, "- ") && indentOf(lines[j]) == stepCol) ||
			(trimmed != "" && !strings.HasPrefix(trimmed, "#") && indentOf(lines[j]) < stepCol) {
			end = j
			break
		}
	}
	if !(start <= identifyingLine && identifyingLine < end) {
		t.Fatalf("instrument failure: located step block [%d,%d) does not contain the identifying line %d", start, end, identifyingLine)
	}
	foundName := false
	for j := start; j < end; j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "- name: "+miscompileStepName || trimmed == "name: "+miscompileStepName {
			foundName = true
			break
		}
	}
	if !foundName {
		t.Fatalf("instrument failure: located block is not the miscompile step %q", miscompileStepName)
	}
	refusals, instrumentErr := continueOnErrorRefusalsIn(lines, start, end, stepCol)
	if instrumentErr != "" {
		t.Fatalf("instrument failure: %s", instrumentErr)
	}
	for _, r := range refusals {
		t.Error(r)
	}
	runRaw, err := os.ReadFile(filepath.Join(repoRoot, miscompileReproducerPath))
	if err != nil {
		t.Fatal(err)
	}
	runSrc := string(runRaw)
	// Require the kernel read (quorum R3), and reject only EXECUTABLE uses of the
	// overridable channel — a comment may NAME the channel it exists to warn about.
	// The round-1 form banned the bare tokens repo-wide and therefore redded against
	// its own documentation on arrival (V19 arm A); this is round-2 R1's fix verbatim.
	for _, need := range []string{"uname -s", "uname -m"} {
		if !strings.Contains(runSrc, need) {
			t.Errorf("run.sh no longer reads the kernel via %q — quorum R3's fix reverted; the polarity must not come from an overridable variable", need)
		}
	}
	for i, line := range strings.Split(runSrc, "\n") {
		code := line
		if idx := strings.Index(code, "#"); idx >= 0 {
			code = code[:idx]
		}
		for _, bad := range []string{"go env GOOS", "go env GOARCH"} {
			if strings.Contains(code, bad) {
				t.Errorf("run.sh:%d executable use of %q — the polarity must not read an overridable channel", i+1, bad)
			}
		}
		if strings.Contains(code, "host_pair") && strings.Contains(code, "go env") {
			t.Errorf("run.sh:%d derives host_pair from `go env` — quorum R3", i+1)
		}
	}
}

// readVerifyGoSh returns verify_go.sh's text; a read failure is an instrument failure, never
// a verdict about the gate.
func readVerifyGoSh(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
	if err != nil {
		t.Fatalf("instrument failure: cannot read verify_go.sh: %v", err)
	}
	return string(raw)
}

// shellAssignmentCount counts line-anchored occurrences of a literal shell assignment, so a
// mention inside a diagnostic string (`... ACTIVE_GO=$ACTIVE_GO ...`) is not miscounted as one.
func shellAssignmentCount(src, assignment string) int {
	return len(regexp.MustCompile(`(?m)^[ \t]*`+regexp.QuoteMeta(assignment)).FindAllString(src, -1))
}

// p1ComparatorCallRe matches the P1 floor gate's comparator call with its two operands captured
// by NAME rather than pinned by spelling. Queue row 60.
//
// Queue row 61: the call must be the `if` CONDITION, so the verdict travels from the comparator
// to the branch with nothing in between that could rewrite it. The predecessor matched the call
// wherever it appeared, which left the consumption shape unbound -- and the consumption shape is
// the whole defect: one inserted `floor_rc=0` after `floor_rc=$?` reopened every refusal branch
// while all six static needles stayed green.
var p1ComparatorCallRe = regexp.MustCompile(`if go_version_ge "\$([A-Za-z_][A-Za-z0-9_]*)" "\$([A-Za-z_][A-Za-z0-9_]*)"; then`)

// p1VerdictLaunderRe matches any assignment of `$?` -- the shape queue row 61 showed is
// re-breakable by a single inserted line.
var p1VerdictLaunderRe = regexp.MustCompile(`(?m)^[ \t]*[A-Za-z_][A-Za-z0-9_]*=\$\?`)

// p1ComparatorBinding is P1d: the D-WORLD-28 floor gate's comparator must take the OBSERVED
// active toolchain first and the go.mod-DERIVED root floor second.
//
// Queue row 60: the contract is the operand ORDER and each operand's DERIVATION, never the
// variables' spelling. The predecessor assertion pinned the literal
// `go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"`, so a consistent rename -- semantically inert,
// `bash -n` clean, the gate behaving identically -- dropped the count to 0 and redded CI. It
// also bound LESS than it read: nothing tied the second operand to the go.mod floor read, so an
// operand reversal achieved by swapping the two ASSIGNMENTS passed it (measured: the predecessor
// literal count stayed 1). Binding the derivation closes both.
//
// Returns the two operand names, or a non-nil error naming the failing conjunct.
func p1ComparatorBinding(src, gate string) (activeName, floorName string, err error) {
	calls := p1ComparatorCallRe.FindAllStringSubmatch(gate, -1)
	if len(calls) != 1 {
		return "", "", fmt.Errorf("P1 comparator call `if go_version_ge \"$X\" \"$Y\"; then` count=%d, want 1: the toolchain floor comparison was removed, duplicated, or its verdict is no longer consumed directly by the branch (queue row 61)", len(calls))
	}
	activeName, floorName = calls[0][1], calls[0][2]
	if activeName == floorName {
		return "", "", fmt.Errorf("P1 comparator compares $%s with itself: the floor gate is vacuously true", activeName)
	}
	floorLit := floorName + `="go$(awk '/^go /{print $2; exit}' go.mod)"`
	if n := len(regexp.MustCompile(`(?m)^[ \t]*`+regexp.QuoteMeta(floorLit)).FindAllString(gate, -1)); n != 1 {
		return "", "", fmt.Errorf("P1 comparator's SECOND operand $%s is assigned from the go.mod floor read %d times inside the P1 block, want 1: the operands are reversed, or the second operand is not the derived root floor", floorName, n)
	}
	if n := shellAssignmentCount(src, activeName+`=$(go env GOVERSION)`); n != 1 {
		return "", "", fmt.Errorf("P1 comparator's FIRST operand $%s is assigned from `go env GOVERSION` %d times in verify_go.sh, want 1: the operands are reversed, or the first operand is not the observed active toolchain", activeName, n)
	}
	if n := len(regexp.MustCompile(`(?m)^[ \t]*`+regexp.QuoteMeta(activeName)+`=`).FindAllString(gate, -1)); n != 0 {
		return "", "", fmt.Errorf("P1 block reassigns $%s %d times, want 0: the observed active toolchain must not be shadowed inside the gate", activeName, n)
	}
	if n := len(p1VerdictLaunderRe.FindAllString(gate, -1)); n != 0 {
		return "", "", fmt.Errorf("P1 gate launders the comparator verdict through %d `<var>=$?` assignment(s), want 0: queue row 61 -- a verdict held in a reassignable variable is reopened by a single inserted line, and the gate then prints an arithmetically false success line", n)
	}
	return activeName, floorName, nil
}

// p1NeedleBindings runs the three identifier-bound needles as one unit: P1d's comparator
// binding, P2's execution binding, and the cross-binding that ties them together. The control
// test below evaluates mutants through this same function, so no conjunct is left without a
// red arm.
func p1NeedleBindings(src, gate string) (activeName string, err error) {
	activeName, _, err = p1ComparatorBinding(src, gate)
	if err != nil {
		return "", err
	}
	_, p2Name, err := p2ExecutionBinding(src)
	if err != nil {
		return "", err
	}
	if activeName != p2Name {
		return "", fmt.Errorf("P1 floor gate vets $%s while the race control runs under $%s: the gate does not vet the toolchain the control it protects actually uses", activeName, p2Name)
	}
	return activeName, nil
}

// p2RaceExecRe matches the race control's execution binding with the GOTOOLCHAIN variable
// captured by NAME rather than pinned by spelling. Queue row 60, second call site.
var p2RaceExecRe = regexp.MustCompile(`GOTOOLCHAIN="\$([A-Za-z_][A-Za-z0-9_]*)" go run -race \.`)

// p2ExecutionBinding is P2: the race-detector known-positive control must run under the OBSERVED
// active toolchain, never under nested auto-selection. Like P1d it binds the variable's
// DERIVATION rather than its spelling -- measured this iteration, the predecessor literal was the
// SECOND needle in this file that redded on a consistent, semantically inert rename, one call
// site over from the one queue row 60 named.
func p2ExecutionBinding(src string) (at int, name string, err error) {
	locs := p2RaceExecRe.FindAllStringSubmatchIndex(src, -1)
	if len(locs) != 1 {
		return -1, "", fmt.Errorf("P2 execution-binding needle `GOTOOLCHAIN=\"$X\" go run -race .` count=%d, want 1: the race control can silently return to nested toolchain auto-selection", len(locs))
	}
	name = src[locs[0][2]:locs[0][3]]
	if n := shellAssignmentCount(src, name+`=$(go env GOVERSION)`); n != 1 {
		return -1, "", fmt.Errorf("P2 pins GOTOOLCHAIN to $%s, which is assigned from `go env GOVERSION` %d times in verify_go.sh, want 1: the race control is not bound to the OBSERVED active toolchain", name, n)
	}
	return locs[0][0], name, nil
}

func p1FloorGateRegion(t *testing.T, src string) (comparator, gate string, start int) {
	t.Helper()
	const sentinel = "# --- P1 (queue row 48"
	const success = `echo "   ✓ toolchain floor gate:`
	if count := strings.Count(src, sentinel); count != 1 {
		t.Fatalf("P1 block sentinel %q count=%d, want 1: the D-WORLD-28 fail-closed block is absent or duplicated", sentinel, count)
	}
	if count := strings.Count(src, success); count != 1 {
		t.Fatalf("P1 success-line sentinel %q count=%d, want 1: the D-WORLD-28 fail-closed block cannot be delimited", success, count)
	}
	start = strings.Index(src, sentinel)
	end := strings.Index(src, success)
	if end < start {
		t.Fatalf("P1 success-line sentinel precedes the block sentinel: the D-WORLD-28 floor gate cannot be delimited")
	}
	end += len(success)
	region := src[start:end]
	closeAt := strings.Index(region, "\n}\n")
	if closeAt < 0 {
		t.Fatalf("P1 comparator function has no closing line: the three-way comparison contract cannot be inspected")
	}
	closeAt += len("\n}\n")
	return region[:closeAt], region[closeAt:], start
}

func TestRaceControlFloorStaysBelowRootToolchain(t *testing.T) {
	raceFloor := moduleGoFloor(t, filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "racecontrol", "go.mod"))
	if !version.IsValid(raceFloor) {
		t.Fatalf("instrument failure: racecontrol module floor %q is not a valid Go version", raceFloor)
	}
	rootFloor := moduleGoFloor(t, filepath.Join(repoRoot, "go.mod"))
	if !version.IsValid(rootFloor) {
		t.Fatalf("instrument failure: root module floor %q is not a valid Go version", rootFloor)
	}
	if version.Compare(raceFloor, rootFloor) > 0 {
		t.Fatalf("racecontrol module floor %q is above the root module floor %q: the control refuses before it can fire and verify_go.sh FATALs that the race detector is not armed for the wrong reason", raceFloor, rootFloor)
	}

	p1NeedleSet(t, readVerifyGoSh(t))
}

// p1NeedleSet is the whole P1/P2 needle set as a function of verify_go.sh's TEXT, so the
// anti-brittleness green control below can run the SAME assertions against a mutated source.
// Queue row 60: the predecessor set was reachable only through a disk read, so its one green
// control (M17, a reworded comment) could not be widened to cover an identifier rename -- and
// THREE separate needles turned out to pin an identifier, each one invisible until the one
// before it was relaxed.
func p1NeedleSet(t *testing.T, src string) {
	t.Helper()
	const verifyPath = "scripts/verify_go.sh"
	p2At, _, err := p2ExecutionBinding(src)
	if err != nil {
		t.Fatalf("%s: %v", verifyPath, err)
	}

	comparator, gate, p1At := p1FloorGateRegion(t, src)
	for _, needle := range []string{
		`[ "$root_go_lines" -ne 1 ]`,
		"is BELOW the root module floor",
		"cannot order toolchain tokens",
	} {
		if count := strings.Count(gate, needle); count != 1 {
			t.Fatalf("P1 floor-read/refusal branch %q count=%d, want 1: a refusal branch of the D-WORLD-28 floor gate was removed or reworded", needle, count)
		}
	}
	if count := strings.Count(gate, "verify_go.sh: FATAL:"); count != 3 {
		t.Fatalf("P1 floor gate FATAL refusal count=%d, want 3: all three failure modes must remain separately attributed", count)
	}
	if count := strings.Count(gate, "exit 1"); count != 3 {
		t.Fatalf("P1 block has %d `exit 1` statements, want 3 (one per refusal)", count)
	}
	if count := strings.Count(gate, "exit 0"); count != 0 {
		t.Fatalf("P1 refusal gate has %d `exit 0` statements, want 0: a refusal must not report success", count)
	}
	for _, exit := range []string{"exit 0", "exit 1", "exit 2"} {
		if !strings.Contains(comparator, exit) {
			t.Fatalf("P1 comparator is missing %q: its >= / < / malformed three-way contract was gutted", exit)
		}
	}
	p1Active, err := p1NeedleBindings(src, gate)
	if err != nil {
		t.Fatalf("%s: %v", verifyPath, err)
	}
	if versions := regexp.MustCompile(`go1\.[0-9]+\.[0-9]+`).FindAllString(comparator+gate, -1); len(versions) != 0 {
		t.Fatalf("P1 block contains hardcoded Go version literal(s) %v: the root floor must be read from go.mod", versions)
	}
	const floorRead = `awk '/^go /{print $2; exit}' go.mod`
	if count := strings.Count(gate, floorRead); count != 1 {
		t.Fatalf("P1 derived root-floor read %q count=%d, want 1: the floor must be read, never hardcoded", floorRead, count)
	}
	// Row 60, THIRD call site: the deny-list anchor is located by the same derived toolchain
	// variable P1d and P2 bind, so it too is inert under a consistent rename.
	denyLit := `case "$` + p1Active + `" in`
	denyAt := strings.Index(src, denyLit)
	if denyAt < 0 {
		t.Fatalf("instrument failure: verify_go.sh deny-list case %q is absent, so P1 ordering cannot be checked", denyLit)
	}
	if !(denyAt < p1At && p1At < p2At) {
		t.Fatalf("P1 floor gate is out of order (deny-list@%d, P1@%d, race leg@%d): it must vet after the deny-list and before the race control", denyAt, p1At, p2At)
	}
}

// assertRed requires a mutant to be rejected AND to be rejected for the stated reason.
//
// Measured this iteration: the needle set's conjuncts overlap, so four of them have no SOLE
// killer -- neutering any one leaves every arm green because another conjunct catches the same
// mutant. Pinning the message is what stops an arm from silently re-attributing to a different
// conjunct when one is weakened, which a bare `err != nil` cannot see.
func assertRed(t *testing.T, err error, wantSubstring, why string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: the binding accepted the mutant", why)
		return
	}
	if !strings.Contains(err.Error(), wantSubstring) {
		t.Errorf("%s: rejected, but for the wrong reason -- want a message containing %q, got: %v", why, wantSubstring, err)
	}
}

// TestP1NeedleSetIsInertUnderRenameAndCatchesReversal is the DURABLE green/red control pair for
// the P1/P2 needle set (queue row 60).
//
// Row 48 shipped the set with one green control, M17, which rewords a COMMENT word inside the P1
// block. That supports "the set is not any-edit-reds" for comments and leaves it UNMEASURED for
// identifiers -- and the identifier case is exactly where the set redded on a semantically inert
// edit. A drill-only green arm also decays the moment the drill stops being re-run, so the rename
// arm lives here as a test.
//
// The GREEN arm runs `p1NeedleSet` itself rather than the two binding helpers, so a needle added
// to the set later is covered without anyone remembering to widen this test. That matters: fixing
// the one needle row 60 named surfaced a second, and fixing that one surfaced a third.
//
// A green arm alone would be satisfiable by an assertion that never fails, so it is paired with
// five red arms, each naming a distinct fail-OPEN reshaping of the comparison.
func TestP1NeedleSetIsInertUnderRenameAndCatchesReversal(t *testing.T) {
	src := readVerifyGoSh(t)
	_, gate, _ := p1FloorGateRegion(t, src)

	activeName, floorName, err := p1ComparatorBinding(src, gate)
	if err != nil {
		t.Fatalf("instrument failure: the P1d binding must hold on the pristine verify_go.sh, got: %v", err)
	}
	if _, _, err := p2ExecutionBinding(src); err != nil {
		t.Fatalf("instrument failure: the P2 binding must hold on the pristine verify_go.sh, got: %v", err)
	}
	floorLit := floorName + `="go$(awk '/^go /{print $2; exit}' go.mod)"`
	callLit := `go_version_ge "$` + activeName + `" "$` + floorName + `"`
	raceLit := `GOTOOLCHAIN="$` + activeName + `" go run -race .`

	// bind re-derives the P1 region from a mutated source and evaluates BOTH identifier-bound
	// needles, so a red arm cannot pass by escaping the one it was aimed at. A mutation that
	// breaks the region delimiters is an instrument failure, not a scored red.
	bind := func(t *testing.T, mutated string) error {
		t.Helper()
		if mutated == src {
			t.Fatal("instrument failure: the mutation changed nothing")
		}
		_, mutGate, _ := p1FloorGateRegion(t, mutated)
		_, err := p1NeedleBindings(mutated, mutGate)
		return err
	}

	t.Run("GREEN/consistent rename is inert across the WHOLE needle set", func(t *testing.T) {
		renamed := strings.ReplaceAll(src, floorName, "P1RenamedFloor")
		renamed = strings.ReplaceAll(renamed, activeName, "P1RenamedActive")
		if renamed == src {
			t.Fatal("instrument failure: the rename changed nothing")
		}
		if strings.Contains(renamed, activeName) || strings.Contains(renamed, floorName) {
			t.Fatalf("instrument failure: the rename is not consistent -- %q or %q survives", activeName, floorName)
		}
		p1NeedleSet(t, renamed)
	})

	t.Run("RED/operands swapped in the call", func(t *testing.T) {
		swapped := strings.Replace(src, callLit, `go_version_ge "$`+floorName+`" "$`+activeName+`"`, 1)
		assertRed(t, bind(t, swapped), "is assigned from the go.mod floor read 0 times",
			"swapping the comparator's operands inverts the floor gate and must red")
	})

	t.Run("RED/operands swapped by reassignment", func(t *testing.T) {
		// The same inversion reached without touching the call: shadow the observed toolchain
		// inside the block and hand the derived floor to the FIRST operand. Measured this
		// iteration, the predecessor literal needle read GREEN on exactly this.
		reassigned := strings.Replace(src, floorLit,
			floorName+`=$(go env GOVERSION)`+"\n"+activeName+`="go$(awk '/^go /{print $2; exit}' go.mod)"`, 1)
		assertRed(t, bind(t, reassigned), "is assigned from the go.mod floor read 0 times",
			"inverting the gate by swapping the two assignments must red")
	})

	t.Run("RED/active toolchain shadowed inside the block", func(t *testing.T) {
		// Sole killer for the shadow conjunct: the floor assignment is untouched, so only the
		// shadow check can see that the vetted value is no longer the observed one.
		shadowed := strings.Replace(src, floorLit, activeName+`=$(cat /dev/null)`+"\n"+floorLit, 1)
		assertRed(t, bind(t, shadowed), "must not be shadowed inside the gate",
			"reassigning the observed toolchain inside the block makes the gate vet a value it did not observe")
	})

	t.Run("RED/active toolchain not read from go env GOVERSION", func(t *testing.T) {
		// Sole killer for the active-derivation conjunct: a plausible-looking substitute source
		// that is not the authoritative one.
		unobserved := strings.Replace(src, activeName+`=$(go env GOVERSION)`,
			activeName+`=$(go version | awk '{print $3}')`, 1)
		assertRed(t, bind(t, unobserved), "P1 comparator's FIRST operand",
			"the vetted toolchain must be read from `go env GOVERSION`, not from a parsed substitute")
	})

	t.Run("RED/comparator compares one operand with itself", func(t *testing.T) {
		vacuous := strings.Replace(src, callLit, `go_version_ge "$`+activeName+`" "$`+activeName+`"`, 1)
		assertRed(t, bind(t, vacuous), "with itself: the floor gate is vacuously true",
			"a self-comparison makes the floor gate vacuously true and must red")
	})

	t.Run("RED/comparator verdict not consumed directly by the branch", func(t *testing.T) {
		// Queue row 61, sole killer for the consumption-shape conjunct: the call, both operands
		// and every refusal branch survive verbatim, and the gate is nonetheless fail-OPEN
		// because the disjunction makes the condition unconditionally true. It carries no
		// `<var>=$?`, so the laundering conjunct cannot reach it.
		disjoined := strings.Replace(src, "if "+callLit+"; then", "if "+callLit+" || true; then", 1)
		assertRed(t, bind(t, disjoined), "no longer consumed directly by the branch",
			"the comparator's verdict must reach the branch unmediated; a disjunction opens the gate while every needle-visible byte survives")
	})

	t.Run("RED/comparator verdict re-laundered through a reassignable variable", func(t *testing.T) {
		// Queue row 61, sole killer for the laundering conjunct: the `if` form is untouched, so
		// only the `<var>=$?` check can see that the verdict has been copied somewhere an
		// inserted line could rewrite. This is the shape the row measured as fail-open.
		relaundered := strings.Replace(src, "\nfi\n"+`echo "   ✓ toolchain floor gate:`,
			"\nfi\nfloor_rc=$?\n"+`echo "   ✓ toolchain floor gate:`, 1)
		assertRed(t, bind(t, relaundered), "launders the comparator verdict through",
			"copying the verdict into a reassignable variable restores the row 61 fail-open shape and must red")
	})

	t.Run("RED/race control unbound from the observed toolchain", func(t *testing.T) {
		unbound := strings.Replace(src, raceLit, `go run -race .`, 1)
		assertRed(t, bind(t, unbound), "count=0, want 1: the race control can silently return",
			"dropping GOTOOLCHAIN returns the race control to nested auto-selection and must red")
	})

	t.Run("RED/floor gate vets a different toolchain than the race control runs", func(t *testing.T) {
		// A second, genuinely observed variable: each needle passes on its own, and only the
		// cross-binding notices that the gate no longer vets the toolchain it protects.
		divergent := strings.Replace(src, raceLit, `GOTOOLCHAIN="$P1OtherGo" go run -race .`, 1)
		divergent = strings.Replace(divergent, activeName+`=$(go env GOVERSION)`,
			activeName+`=$(go env GOVERSION)`+"\n"+`P1OtherGo=$(go env GOVERSION)`, 1)
		assertRed(t, bind(t, divergent), "the gate does not vet the toolchain the control it protects actually uses",
			"the floor gate must vet the same toolchain variable the race control runs under")
	})

	t.Run("RED/race control bound to an underived variable", func(t *testing.T) {
		hijacked := strings.Replace(src, raceLit, `GOTOOLCHAIN="$P1Unobserved" go run -race .`, 1)
		assertRed(t, bind(t, hijacked), "the race control is not bound to the OBSERVED active toolchain",
			"binding GOTOOLCHAIN to a variable never read from `go env GOVERSION` must red")
	})
}

// stepScanFixture renders a two-step block-sequence job whose FIRST step is the guarded one, so
// every arm below is scoped exactly as the production scan is: `- ` at column 6, the step's own
// keys at column 8, a `run:` block scalar at column 10, and an unrelated sibling step after it.
// guarded is spliced verbatim between the dash line and the run scalar; extra rides the dash.
func stepScanFixture(dashKey, guarded, runBody string) (lines []string, start, end, stepCol int) {
	dash := "      - name: guarded"
	if dashKey != "" {
		dash = "      - " + dashKey
	}
	src := []string{
		"jobs:",
		"  j:",
		"    steps:",
		dash,
	}
	if dashKey != "" {
		src = append(src, "        name: guarded")
	}
	if guarded != "" {
		src = append(src, "        "+guarded)
	}
	src = append(src,
		"        run: |",
		"          "+runBody,
		"      - name: unrelated",
		"        continue-on-error: true",
		"        run: echo other",
	)
	return src, 3, len(src) - 3, 6
}

// TestContinueOnErrorStepScanReadsTheKeyNotTheText is queue row 62's non-vacuous guard. Each arm
// names the conjunct it is the sole killer for; the mutation drill in the iteration record lands
// each neutering separately and confirms only the arm below goes red.
func TestContinueOnErrorStepScanReadsTheKeyNotTheText(t *testing.T) {
	cases := []struct {
		name     string
		dashKey  string
		guarded  string
		runBody  string
		wantRed  bool
		wantsVal string
	}{
		// The row-44 property this instrument exists for. Killer for: rejecting a true value.
		{name: "explicit true on a sibling key is refused", guarded: "continue-on-error: true", runBody: "./run.sh", wantRed: true, wantsVal: "true"},
		// Queue row 62(a): a legitimate explicit opt-OUT swallows nothing.
		{name: "explicit false on a sibling key is compliant", guarded: "continue-on-error: false", runBody: "./run.sh"},
		{name: "quoted and cased false is compliant", guarded: `continue-on-error: "False"`, runBody: "./run.sh"},
		{name: "false with a trailing comment is compliant", guarded: "continue-on-error: false # deliberate, row 62", runBody: "./run.sh"},
		// Queue row 62(b): script TEXT inside the step's own run: scalar is not a key.
		{name: "the flag named inside the run scalar is not a key", runBody: `echo "note: never add continue-on-error to this step"`},
		{name: "the flag named in a comment inside the block is not a key", guarded: "# never add continue-on-error here", runBody: "./run.sh"},
		// The hole in row 62's OWN proposed remedy: a step's first key rides the dash, two
		// columns left of every sibling key, so an indentation rule written only for siblings
		// is fail-OPEN on the shape an author is most likely to reach for.
		{name: "true riding the block-sequence dash is refused", dashKey: "continue-on-error: true", runBody: "./run.sh", wantRed: true, wantsVal: "true"},
		{name: "false riding the block-sequence dash is compliant", dashKey: "continue-on-error: false", runBody: "./run.sh"},
		// Not statically decidable, so it fails CLOSED rather than passing as "not true".
		{name: "an expression value is refused", guarded: "continue-on-error: ${{ github.event_name == 'push' }}", runBody: "./run.sh", wantRed: true, wantsVal: "${{ github.event_name == 'push' }}"},
		{name: "an empty value is refused", guarded: "continue-on-error:", runBody: "./run.sh", wantRed: true, wantsVal: ""},
		// A different key that merely starts with the same text is not this key.
		{name: "a longer key with the same prefix is not this key", guarded: "continue-on-error-policy: strict", runBody: "./run.sh"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines, start, end, stepCol := stepScanFixture(c.dashKey, c.guarded, c.runBody)
			refusals, instrumentErr := continueOnErrorRefusalsIn(lines, start, end, stepCol)
			if instrumentErr != "" {
				t.Fatalf("instrument failure: %s", instrumentErr)
			}
			if c.wantRed {
				if len(refusals) != 1 {
					t.Fatalf("want exactly 1 refusal, got %d: %v", len(refusals), refusals)
				}
				if !strings.Contains(refusals[0], `to "`+c.wantsVal+`"`) {
					t.Errorf("refusal does not report the value %q it read: %s", c.wantsVal, refusals[0])
				}
				return
			}
			if len(refusals) != 0 {
				t.Errorf("want no refusal, got %v", refusals)
			}
		})
	}
}

// TestContinueOnErrorStepScanIsBoundedToItsOwnStep is the boundary control the row-52 sprint's
// V19 arm established, re-asserted against the new scan: the fixture's UNRELATED step sets
// `continue-on-error: true`, and every arm above is green on it. This is the known-POSITIVE half
// — move the range over that step and the same scan must refuse — so a green above proves the
// scope holds rather than that the scan stopped working.
func TestContinueOnErrorStepScanIsBoundedToItsOwnStep(t *testing.T) {
	lines, _, _, stepCol := stepScanFixture("", "continue-on-error: false", "./run.sh")
	unrelated := len(lines) - 3
	if got := strings.TrimSpace(lines[unrelated]); got != "- name: unrelated" {
		t.Fatalf("instrument failure: fixture's unrelated step moved to %q", got)
	}
	refusals, instrumentErr := continueOnErrorRefusalsIn(lines, unrelated, len(lines), stepCol)
	if instrumentErr != "" {
		t.Fatalf("instrument failure: %s", instrumentErr)
	}
	if len(refusals) != 1 {
		t.Fatalf("control did not fire: scanning the unrelated step must refuse its `continue-on-error: true`, got %v", refusals)
	}
}

// TestContinueOnErrorStepScanRefusesAFlowMapping pins the fail-CLOSED disposition for the one
// YAML shape the scan cannot read: a step written as a flow mapping hides its keys from an
// indentation rule, so the scan reports an instrument failure instead of a green.
func TestContinueOnErrorStepScanRefusesAFlowMapping(t *testing.T) {
	lines := []string{
		"    steps:",
		"      - {name: guarded, continue-on-error: true, run: ./run.sh}",
		"      - name: unrelated",
	}
	refusals, instrumentErr := continueOnErrorRefusalsIn(lines, 1, 2, 6)
	if instrumentErr == "" {
		t.Fatalf("a flow-mapping step must be an instrument failure, not a silent pass (refusals=%v)", refusals)
	}
	if !strings.Contains(instrumentErr, "flow mapping") {
		t.Errorf("instrument failure does not name the shape it could not read: %s", instrumentErr)
	}
}
