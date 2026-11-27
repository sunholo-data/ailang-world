package verifygate

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Row 181 (attended 2026-10-16): the -race leg runs in its own CI job. These checks keep the
// split from silently dropping -race: the race job must exist and run the race legs, the main
// job must run the plain legs, and verify_go.sh must keep the race leg behind a loud bound.
func TestRaceLegRunsInItsOwnJob(t *testing.T) {
	job := ciRaceJobReader(t)
	plain, race := job("go-verify"), job("go-race-verify")
	for _, c := range []struct{ body, want, why string }{
		{plain, "VERIFY_GO_LEGS: plain", "the main Go job must skip only the race leg"},
		{race, "VERIFY_GO_LEGS: race", "the race job must run the race leg"},
		{race, "VERIFY_GO_RACE_TIMEOUT_S:", "the race job must set its measured bound"},
		{race, "run: ./scripts/verify_go.sh", "the race job must run the real gate"},
		{race, "AILANG_BIN=", "the race job must install the pinned interpreter"},
	} {
		if !strings.Contains(c.body, c.want) {
			t.Errorf("ci.yml: %q missing — %s", c.want, c.why)
		}
	}
	if strings.Contains(race, "VERIFY_GO_LEGS: plain") || strings.Contains(plain, "VERIFY_GO_LEGS: race") {
		t.Error("ci.yml: the legs are crossed between go-verify and go-race-verify")
	}

	gate, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
	if err != nil {
		t.Fatal(err)
	}
	g := string(gate)
	for _, want := range []string{
		`if [ "$VERIFY_GO_LEGS" != plain ]; then`,
		raceLegCall,
		`FATAL: -race leg timed out after {bound}s`,
		`FATAL: the race detector is not armed`,
	} {
		if !strings.Contains(g, want) {
			t.Errorf("verify_go.sh: %q missing — the race leg or its control/bound was removed", want)
		}
	}
}

const raceLegCall = `run_race_leg "$VERIFY_GO_RACE_TIMEOUT_S" go test ./... -count=1 -race -timeout "${VERIFY_GO_RACE_PKG_TIMEOUT_S}s"`

func ciRaceJobReader(t *testing.T) func(string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	return func(name string) string {
		t.Helper()
		i := strings.Index(src, "\n  "+name+":\n")
		if i < 0 {
			t.Fatalf("ci.yml has no job %q", name)
		}
		rest := src[i+1:]
		next := regexp.MustCompile(`\n  [a-z][a-z0-9-]*:\n`).FindStringIndex(rest[1:])
		if next == nil {
			return rest
		}
		return rest[:next[0]+1]
	}
}

func TestRaceBoundsAreOrdered(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
	if err != nil {
		t.Fatal(err)
	}
	gate := string(raw)
	race := ciRaceJobReader(t)("go-race-verify")
	one := func(src, pattern string) []string {
		t.Helper()
		matches := regexp.MustCompile(pattern).FindAllStringSubmatch(src, -1)
		if len(matches) != 1 {
			t.Fatalf("race bounds: %q matches %d times, want 1", pattern, len(matches))
		}
		return matches[0]
	}
	seconds := func(src, pattern string) int {
		t.Helper()
		n, err := strconv.Atoi(one(src, pattern)[1])
		if err != nil {
			t.Fatalf("race bounds: %q: %v", pattern, err)
		}
		return n
	}
	p := seconds(gate, `(?m)^VERIFY_GO_RACE_PKG_TIMEOUT_S=(\d+)$`)
	a := seconds(gate, `(?m)^VERIFY_GO_RACE_START_ALLOWANCE_S=(\d+)$`)
	l := seconds(gate, `(?m)^VERIFY_GO_RACE_TIMEOUT_S="\$\{VERIFY_GO_RACE_TIMEOUT_S:-(\d+)\}"$`)
	one(race, `(?m)^    runs-on: ubuntu-24\.04$`)
	if strings.Contains(race, "continue-on-error") {
		t.Error("race job must have zero continue-on-error directives")
	}
	// Match the step containing the gate invocation. Stop at the next step.
	step := one(race, `(?m)^      - [^\n]+\n(?:^ {8}[^\n]*\n)*?^        run: \./scripts/verify_go\.sh$`)[0]
	j := seconds(step, `(?m)^        timeout-minutes: (\d+)$`)
	ciLeg := seconds(step, `(?m)^          VERIFY_GO_RACE_TIMEOUT_S: "(\d+)"$`)
	call := one(gate, `(?m)^run_race_leg "\$VERIFY_GO_RACE_TIMEOUT_S" .*$`)[0]
	if call != raceLegCall {
		t.Errorf("race call must include -race and use the package bound with no -short, -run, -skip or -failfast=false: %s", call)
	}

	// Row 267 design §2 and §0. The class-invariant daemon extension was 166s.
	const slowNeedS, daemonInflatedNeedS = 646, 828
	const latestHeavyStartS, worlddExtensionS, daemonExtensionInflatedS = 554, 93, 347
	const backupKillS, preLegMaxS, stepMarginS = 60, 55, 300
	if p < (slowNeedS*5+3)/4 || p <= daemonInflatedNeedS {
		t.Errorf("package bound %ds must be at least 1.25 × %ds and above inflated need %ds", p, slowNeedS, daemonInflatedNeedS)
	}
	if a < latestHeavyStartS+worlddExtensionS+daemonExtensionInflatedS {
		t.Errorf("start allowance %ds must be at least %ds + %ds + %ds", a, latestHeavyStartS, worlddExtensionS, daemonExtensionInflatedS)
	}
	if ciLeg != l {
		t.Errorf("CI leg bound %ds must equal script default %ds", ciLeg, l)
	}
	if l < p+backupKillS+a || ciLeg < p+backupKillS+a {
		t.Errorf("leg bounds (default %ds, CI %ds) must be at least package %ds + backup %ds + start allowance %ds", l, ciLeg, p, backupKillS, a)
	}
	if j*60 < ciLeg+stepMarginS || stepMarginS < 5*preLegMaxS {
		t.Errorf("step %dmin must cover leg %ds + %ds, and the margin must cover 5 × %ds", j, ciLeg, stepMarginS, preLegMaxS)
	}
}

// TestVerifyGoRefusesUnknownLegs runs the real script: a bad VERIFY_GO_LEGS refuses (rc 2)
// instead of silently running fewer legs. It needs the pinned AILANG_BIN the gate checks first.
func TestVerifyGoRefusesUnknownLegs(t *testing.T) {
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the gate checks the pin before the legs; never skip")
	}
	cmd := exec.Command(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "VERIFY_GO_LEGS=both")
	out, err := cmd.CombinedOutput()
	ee, ok := err.(*exec.ExitError)
	if !ok || ee.ExitCode() != 2 || !strings.Contains(string(out), "VERIFY_GO_LEGS must be all, plain or race") {
		t.Fatalf("VERIFY_GO_LEGS=both: err=%v out=%s", err, out)
	}
}
