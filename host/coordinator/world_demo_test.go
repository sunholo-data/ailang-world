package coordinator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Match the fail-never-skip policy of newSeToolsRig and evidence.pinnedChecker.
// Every subprocess, including version discovery, has the plan's 40s outer bound.
func worldDemoCommand(t *testing.T, bin, dir string, args ...string) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("pinned %v: %v (context=%v)\nstdout: %s\nstderr: %s", args, err, ctx.Err(), out, stderr.String())
	}
	return out
}

func worldDemoPin(t *testing.T) string {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: world-demo requires the pinned v0.41.0 release; never skip")
	}
	absolute, err := filepath.Abs(bin)
	if err != nil {
		t.Fatal(err)
	}
	out := worldDemoCommand(t, absolute, seToolsRepoRoot(t), "--version")
	if line := strings.SplitN(string(out), "\n", 2)[0]; line != "AILANG v0.41.0" {
		t.Fatalf("AILANG_BIN reports %q; want AILANG v0.41.0", line)
	}
	return absolute
}

type worldDemoResult struct {
	Function string `json:"function"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	SkipKind string `json:"skip_kind"`
}

type worldDemoAIReport struct {
	Check struct {
		Passed     bool              `json:"passed"`
		ErrorCount int               `json:"error_count"`
		Errors     []json.RawMessage `json:"errors"`
	} `json:"check"`
	Verify struct {
		Available      bool              `json:"available"`
		Verified       int               `json:"verified"`
		Counterexample int               `json:"counterexample"`
		Skipped        int               `json:"skipped"`
		Errors         int               `json:"errors"`
		Results        []worldDemoResult `json:"results"`
	} `json:"verify"`
}

type worldDemoTestReport struct {
	Success    bool              `json:"success"`
	Total      int               `json:"total_tests"`
	Passed     int               `json:"passed_tests"`
	Failed     int               `json:"failed_tests"`
	Skipped    int               `json:"skipped_tests"`
	Vacuous    int               `json:"vacuous_skips"`
	Tests      []worldDemoResult `json:"tests"`
	Properties []worldDemoResult `json:"properties"`
}

// Same stdout-only JSON extraction as host/verifygate's ai-check probe; stderr
// must never become a JSON prefix. Typed result fields follow the evidence tests.
func worldDemoDecode(raw []byte, into any) error {
	start := bytes.IndexByte(raw, '{')
	if start < 0 {
		return fmt.Errorf("no JSON object in stdout: %q", raw)
	}
	return json.Unmarshal(raw[start:], into)
}

// Validate report health separately from the independent identity manifest.
// The null controls call THIS parser so deleting its empty guard is observable.
func worldDemoParseAI(raw []byte) (worldDemoAIReport, error) {
	var r worldDemoAIReport
	if err := worldDemoDecode(raw, &r); err != nil {
		return r, err
	}
	if !r.Check.Passed || r.Check.ErrorCount != 0 || len(r.Check.Errors) != 0 ||
		!r.Verify.Available || r.Verify.Counterexample != 0 || r.Verify.Skipped != 0 || r.Verify.Errors != 0 {
		return r, fmt.Errorf("unhealthy check/verify report: %+v", r)
	}
	if len(r.Verify.Results) == 0 {
		return r, fmt.Errorf("empty verify results refused")
	}
	for _, item := range r.Verify.Results {
		if item.Function == "" || item.Status != "verified" || item.SkipKind != "" {
			return r, fmt.Errorf("unverified result: %+v", item)
		}
	}
	return r, nil
}

func worldDemoParseTests(raw []byte) (worldDemoTestReport, error) {
	var r worldDemoTestReport
	if err := worldDemoDecode(raw, &r); err != nil {
		return r, err
	}
	if !r.Success || r.Failed != 0 || r.Skipped != 0 || r.Vacuous != 0 {
		return r, fmt.Errorf("unhealthy test report: %+v", r)
	}
	if len(r.Tests) == 0 {
		return r, fmt.Errorf("empty inline tests refused")
	}
	for _, item := range append(append([]worldDemoResult{}, r.Tests...), r.Properties...) {
		if item.Name == "" || item.Status != "pass" || item.SkipKind != "" {
			return r, fmt.Errorf("nonpassing test/property: %+v", item)
		}
	}
	if r.Total != len(r.Tests)+len(r.Properties) || r.Passed != r.Total {
		return r, fmt.Errorf("test totals disagree with named results")
	}
	return r, nil
}

func worldDemoNames(t *testing.T, label string, got, want []string) {
	t.Helper()
	seen := map[string]bool{}
	for _, name := range got {
		if seen[name] {
			t.Fatalf("%s: duplicate identity %q", label, name)
		}
		seen[name] = true
	}
	for _, name := range want {
		if !seen[name] {
			t.Errorf("%s: missing required identity %q", label, name)
		}
	}
	if len(got) != len(want) {
		t.Errorf("%s: got %d identities, want %d", label, len(got), len(want))
	}
}

func TestWorldDemoEvidencePinnedManifest(t *testing.T) {
	bin := worldDemoPin(t)
	dir := filepath.Join(seToolsRepoRoot(t), "packages", "world-demo")
	worldDemoCommand(t, bin, dir, "check", "evidence.ail")
	rawTests := worldDemoCommand(t, bin, dir, "test", "--format", "json", "evidence.ail")
	rawAI := worldDemoCommand(t, bin, dir, "ai-check", "-timeout", "5s", "evidence.ail")
	ai, err := worldDemoParseAI(rawAI)
	if err != nil {
		t.Fatal(err)
	}
	tests, err := worldDemoParseTests(rawTests)
	if err != nil {
		t.Fatal(err)
	}
	var functions, names []string
	for _, r := range ai.Verify.Results {
		functions = append(functions, r.Function)
	}
	for _, r := range tests.Tests {
		names = append(names, r.Name)
	}
	worldDemoNames(t, "Z3 manifest", functions, []string{"legacyGradeCode", "authoredClaimGradeCode", "rosterKindAllowed"})
	if ai.Verify.Verified != 3 {
		t.Errorf("verified total = %d, want 3", ai.Verify.Verified)
	}
	// Independent hardcoded identities, never generated from package source.
	worldDemoNames(t, "inline manifest", names, []string{
		"legacyGradeCode_test_1", "legacyGradeCode_test_2", "legacyGradeCode_test_3",
		"legacyGradeCode_test_4", "legacyGradeCode_test_5", "legacyGradeCode_test_6",
		"legacyGradeCode_test_7", "legacyGradeCode_test_8",
		"authoredClaimGradeCode_test_1", "authoredClaimGradeCode_test_2", "authoredClaimGradeCode_test_3",
		"authoredClaimGradeCode_test_4", "authoredClaimGradeCode_test_5",
		"rosterKindAllowed_test_1", "rosterKindAllowed_test_2", "rosterKindAllowed_test_3", "rosterKindAllowed_test_4",
		"rosterKindAllowed_test_5", "rosterKindAllowed_test_6", "rosterKindAllowed_test_7", "rosterKindAllowed_test_8",
	})
	// A healthy-looking null report must fail in the parser itself, before the
	// independent manifest. This is the M1a-empty mutation's named assertion.
	for _, empty := range [][]worldDemoResult{nil, {}} {
		null := ai
		null.Verify.Results = empty
		null.Verify.Verified = 0
		raw, err := json.Marshal(null)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := worldDemoParseAI(raw); err == nil {
			t.Fatal("null report must be refused: empty verify results accepted")
		}
		nullTests := tests
		nullTests.Tests = empty
		nullTests.Properties = nil
		nullTests.Total, nullTests.Passed = 0, 0
		raw, err = json.Marshal(nullTests)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := worldDemoParseTests(raw); err == nil {
			t.Fatal("null test report must be refused: empty inline tests accepted")
		}
	}
	t.Log("check.passed=true; Z3: legacyGradeCode, authoredClaimGradeCode, rosterKindAllowed VERIFIED; 21 named inline tests; zero errors/counterexamples/skips; null controls refused")
}

func TestWorldDemoLegacyGradeDrift(t *testing.T) {
	bin := worldDemoPin(t)
	root, scratch := seToolsRepoRoot(t), t.TempDir()
	// Stage byte-identical real modules in an isolated source root. Never rewrite
	// kernel imports or substitute a transcribed Go mapping for gradeOf.
	for _, rel := range []string{"world/types.ail", "world/logepoch.ail", "packages/world-demo/evidence.ail"} {
		raw, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		dest := filepath.Join(scratch, rel)
		if strings.HasPrefix(rel, "packages/") {
			dest = filepath.Join(scratch, "evidence.ail")
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dest, raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct {
		kind, value string
		want        int
	}{
		{"CompilerOutput", "CompilerOutput(ref)", 2},
		{"TestReport", "TestReport(ref, true)", 3},
		{"TestReport", "TestReport(ref, false)", 3},
		{"HumanApproval", "HumanApproval(ref)", 2},
		{"AiReview", "AiReview(ref, 0.8)", 1},
		{"RecordedEffect", "RecordedEffect(ref)", 2},
		{"ProofReceipt", "ProofReceipt(ref)", 1},
	}
	var source strings.Builder
	source.WriteString(`module probe
import evidence (legacyGradeCode, authoredClaimGradeCode, rosterKindAllowed)
import world/types (gradeOf, CompilerOutput, TestReport, HumanApproval, AiReview, RecordedEffect, ProofReceipt, PROVEN, TESTED, ATTESTED, CLAIMED)
export func main() -> () ! {IO} {
 let ref = {algo: "sha256", digest: "actual-input"};
`)
	for _, c := range cases {
		fmt.Fprintf(&source, " println(show(legacyGradeCode(%q)));\n println(show(match gradeOf(%s) { PROVEN => 4, TESTED => 3, ATTESTED => 2, CLAIMED => 1 }));\n", c.kind, c.value)
	}
	source.WriteString(` println(show(authoredClaimGradeCode("authored-claim")));
 println(show(authoredClaimGradeCode("AiReview")));
 println(show(legacyGradeCode("authored-claim")));
 println(show(rosterKindAllowed("authored-claim")));
 println(show(rosterKindAllowed("AiReview")))
}
`)
	if err := os.WriteFile(filepath.Join(scratch, "probe.ail"), []byte(source.String()), 0600); err != nil {
		t.Fatal(err)
	}
	worldDemoCommand(t, bin, scratch, "check", "probe.ail")
	out := worldDemoCommand(t, bin, scratch, "run", "--quiet", "--caps", "IO", "--entry", "main", "probe.ail")
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2*len(cases)+5 {
		t.Fatalf("runtime probe: want %d lines, got %q", 2*len(cases)+5, out)
	}
	for i, c := range cases {
		pkg, err := strconv.Atoi(lines[2*i])
		if err != nil {
			t.Fatal(err)
		}
		kernel, err := strconv.Atoi(lines[2*i+1])
		if err != nil {
			t.Fatal(err)
		}
		if pkg != kernel {
			t.Errorf("%s grade differs from executable kernel law: package=%d kernel=%d", c.kind, pkg, kernel)
		}
		if pkg != c.want || kernel != c.want || pkg == 4 || kernel == 4 {
			t.Errorf("%s: package=%d kernel=%d; want %d, never PROVEN", c.kind, pkg, kernel, c.want)
		}
		t.Logf("%s: package=%d executable kernel=%d", c.kind, pkg, kernel)
	}
	tail := strings.Join(lines[2*len(cases):], ",")
	if tail != "1,0,0,true,false" {
		t.Fatalf("authored claim must be CLAIMED, distinct from AiReview and legacy kinds: got %s", tail)
	}
	// Compare the new kind to an executed kernel CLAIMED result, without casting
	// authored text into an AiReview. Raw ProofReceipt still has no proof authority.
	if lines[2*len(cases)] != lines[2*(len(cases)-1)+1] {
		t.Fatal("authored claim differs from executable kernel CLAIMED mapping")
	}
	t.Log("authored-claim=CLAIMED; AiReview impersonation refused; raw ProofReceipt=CLAIMED; no mapping produces PROVEN")
}
