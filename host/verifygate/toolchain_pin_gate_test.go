package verifygate

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func normalizeToolchainPin(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '\'' && value[len(value)-1] == '\'') ||
		(value[0] == '"' && value[len(value)-1] == '"')) {
		value = value[1 : len(value)-1]
	}
	if value != "" && !strings.HasPrefix(value, "go") {
		value = "go" + value
	}
	return value
}

func pinValues(lines []string, key string) []string {
	var values []string
	for _, line := range lines {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) == 2 && parts[0] == key {
			values = append(values, normalizeToolchainPin(parts[1]))
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
			floors = append(floors, normalizeToolchainPin(strings.TrimPrefix(line, "go ")))
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
			jobs, wantJobs, len(pinValues(lines, "GOTOOLCHAIN")), len(pinValues(lines, "go-version")))
	}

	goToolchains := pinValues(lines, "GOTOOLCHAIN")
	goVersions := pinValues(lines, "go-version")
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
// and the failure message). The guard's firing is proven at sprint time by AC6's guard-trip
// run and exercised on every CI invocation of run.sh (non-gating, continue-on-error: true,
// ci.yml:172) — the round-1 SKIPPED hole (a banner printed with the pin unprobed) is closed
// in run.sh itself, not merely narrowed here. What remains open by scope: nothing WATCHES
// the non-gating log — a loud failure nobody reads is loud only on inspection; flipping
// ci.yml:172 to gating is the named follow-up in Deferred Scope, paired with OD-1.
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
			t.Errorf("instrument failure: %s does not contain known-positive control %q", scriptPath, control)
		}
	}
	shebangs := 0
	for _, line := range lines {
		if line == "#!/usr/bin/env bash" {
			shebangs++
		}
	}
	if shebangs != 1 {
		t.Errorf("instrument failure: %s exact shebang count=%d, want 1", scriptPath, shebangs)
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
		pinned = normalizeToolchainPin(pinnedAssignments[0])
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
