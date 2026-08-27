package verifygate

import (
	"go/version"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
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
	for _, control := range []string{"ailang-verify:", "go-verify:", "uses: actions/setup-go@v5", "./scripts/verify_go.sh"} {
		if !strings.Contains(src, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", workflowPath, control)
		}
	}

	jobLine := regexp.MustCompile(`^  ([a-z0-9-]+):$`)
	seenJobs := false
	var jobs []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "jobs:" {
			seenJobs = true
			continue
		}
		if seenJobs {
			if match := jobLine.FindStringSubmatch(line); match != nil {
				jobs = append(jobs, match[1])
			}
		}
	}
	slices.Sort(jobs)
	// A third job moves this hand-maintained set in the same edit or this test reds.
	wantJobs := []string{"ailang-verify", "go-verify"}
	if !slices.Equal(jobs, wantJobs) {
		t.Errorf("ci.yml: enumerated jobs=%v, want %v; GOTOOLCHAIN pins=%d go-version pins=%d",
			jobs, wantJobs, len(keyedValues(lines, "GOTOOLCHAIN")), len(keyedValues(lines, "go-version")))
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
	if len(goToolchains) != len(wantJobs) {
		t.Errorf("ci.yml: GOTOOLCHAIN keyed-line count=%d, want %d (one per enumerated expected job)", len(goToolchains), len(wantJobs))
	}
	if len(goVersions) != len(wantJobs) {
		t.Errorf("ci.yml: go-version keyed-line count=%d, want %d (one per enumerated expected job)", len(goVersions), len(wantJobs))
	}
	if setupGoUses != len(wantJobs) {
		t.Errorf("ci.yml: actions/setup-go use count=%d, want %d (one per enumerated expected job)", setupGoUses, len(wantJobs))
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
	for _, control := range []string{"KNOWN_BAD=", "KNOWN_GOOD=", "PINNED="} {
		if !strings.Contains(src, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", scriptPath, control)
		}
	}
	shebangs := 0
	for _, line := range lines {
		if line == "#!/usr/bin/env bash" {
			shebangs++
		}
	}
	if shebangs != 1 {
		t.Fatalf("instrument failure: %s exact shebang count=%d, want 1", scriptPath, shebangs)
	}

	goodAssignments := shellAssignmentValues(lines, "KNOWN_GOOD")
	badAssignments := shellAssignmentValues(lines, "KNOWN_BAD")
	pinnedAssignments := shellAssignmentValues(lines, "PINNED")
	if len(goodAssignments) != 1 {
		t.Errorf("%s: KNOWN_GOOD assignment count=%d, want 1", scriptPath, len(goodAssignments))
	}
	if len(badAssignments) != 1 {
		t.Errorf("%s: KNOWN_BAD assignment count=%d, want 1", scriptPath, len(badAssignments))
	}
	if len(pinnedAssignments) != 1 {
		t.Errorf("%s: PINNED assignment count=%d, want 1", scriptPath, len(pinnedAssignments))
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
		t.Errorf("%s: KNOWN_BAD must contain at least one toolchain", scriptPath)
	}

	floor := moduleGoFloor(t, filepath.Join(repoRoot, "go.mod"))
	if !slices.Contains(good, floor) {
		t.Errorf("KNOWN_GOOD=%v does not probe the pinned toolchain %s from go.mod", good, floor)
	}
	pinned := ""
	if len(pinnedAssignments) == 1 {
		pinned = requireToolchainNamePin(t, scriptPath, "PINNED", pinnedAssignments[0])
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
}

func TestReproModuleFloorStaysBelowKnownBadToolchains(t *testing.T) {
	reproFloor := moduleGoFloor(t, filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "repro", "go.mod"))
	if !version.IsValid(reproFloor) {
		t.Fatalf("instrument failure: repro module floor %q is not a valid Go version", reproFloor)
	}

	scriptPath := filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "run.sh")
	raw, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	badAssignments := shellAssignmentValues(strings.Split(string(raw), "\n"), "KNOWN_BAD")
	if len(badAssignments) != 1 {
		t.Fatalf("instrument failure: %s: KNOWN_BAD assignment count=%d, want 1", scriptPath, len(badAssignments))
	}
	bad := strings.Fields(badAssignments[0])
	if len(bad) == 0 {
		t.Fatalf("instrument failure: %s: KNOWN_BAD must contain at least one toolchain", scriptPath)
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

func TestCanaryDeclaresPositiveArmOnly(t *testing.T) {
	canaryPath := filepath.Join(repoRoot, "host", "store", "toolchain_canary_test.go")
	raw, err := os.ReadFile(canaryPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if count := strings.Count(src, "stateRoot"); count < 2 {
		t.Fatalf("instrument failure: %s count(%q)=%d, want at least 2 before checking comment fences", canaryPath, "stateRoot", count)
	}
	if count := strings.Count(src, "GOTOOLCHAIN"); count != 0 {
		t.Errorf("%s count(%q)=%d, want 0; known-bad arms belong in the nested repro module", canaryPath, "GOTOOLCHAIN", count)
	}
	if !strings.Contains(src, "POSITIVE ARM ONLY") {
		t.Errorf("%s: required %q marker is absent", canaryPath, "POSITIVE ARM ONLY")
	}
}

const miscompileReproducerPath = "design_docs/verification/w-race-gate-blindspot/run.sh"

// TestMiscompileInstrumentStepIsGatedInCI pins the row-44 wiring on two channels that
// must not silently return. (1) `continue-on-error: true` converts an instrument's
// loudest possible output into silence, so it is forbidden in the miscompile step's
// own block. A flag on an unrelated step remains that step's business. (2) The
// instrument's platform polarity reads the KERNEL (`uname`); `go env` honours the
// env-var form of the platform tokens (measured in the design doc, P16), so executable
// uses of that overridable channel are forbidden in run.sh, and both kernel reads are
// asserted present so the probe cannot quietly revert.
// DECLARED RESIDUAL: a step-level `if:` that never evaluates true disables the step
// with this text intact (no actionlint runs in this repo — P41 V18); and these are
// byte-substring pins — a computed assignment (`eval`, string concatenation) evades
// them; the mechanism's runtime immunity is why that gap is acceptable (design doc,
// residuals 2 and 3).
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
	start := -1
	for i, l := range lines {
		if strings.Contains(l, miscompileReproducerPath) {
			for j := i; j >= 0; j-- {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "- name:") {
					start = j
					break
				}
			}
			break
		}
	}
	if start < 0 {
		t.Fatalf("instrument failure: could not locate the miscompile step block in ci.yml")
	}
	end := len(lines)
	for j := start + 1; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "- name:") {
			end = j
			break
		}
	}
	for i := start; i < end; i++ {
		if strings.Contains(lines[i], "continue-on-error") {
			t.Errorf("ci.yml:%d re-introduces %q in the miscompile step — row 44: a swallowed refusal is how this instrument died the first time", i+1, "continue-on-error")
		}
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
