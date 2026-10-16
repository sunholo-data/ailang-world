package verifygate

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Row 181 (attended 2026-10-16): the -race leg runs in its own CI job. These checks keep the
// split from silently dropping -race: the race job must exist and run the race legs, the main
// job must run the plain legs, and verify_go.sh must keep the race leg behind a loud bound.
func TestRaceLegRunsInItsOwnJob(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	job := func(name string) string {
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
		`"go", "test", "./...", "-count=1", "-race"`,
		`FATAL: -race leg timed out after {bound}s`,
		`FATAL: the race detector is not armed`,
	} {
		if !strings.Contains(g, want) {
			t.Errorf("verify_go.sh: %q missing — the race leg or its control/bound was removed", want)
		}
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
