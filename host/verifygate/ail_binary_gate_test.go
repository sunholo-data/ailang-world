package verifygate

import (
	"bytes"
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

var (
	repoRoot = findRepoRoot()
	pinned   = "/tmp/ailang-v0300/ailang"
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
	out, err := exec.Command(pinned, "--version").CombinedOutput()
	if err != nil || !strings.HasPrefix(string(out), "AILANG v0.30.0") {
		t.Fatalf("pinned delegate unavailable or wrong (never skip): err=%v output=%q", err, out)
	}
}

func runGate(t *testing.T, env map[string]string) (int, string) {
	t.Helper()
	cmd := exec.Command(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	cmd.Dir = repoRoot
	blocked := map[string]bool{
		"AILANG_BIN": true, "WORLD_PKG_AILANG_BIN": true,
		"AILANG_SHIM_VERSION_LINE": true, "AILANG_SHIM_DELEGATE": true,
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

func TestKnownPositiveDelegates(t *testing.T) {
	requirePinned(t)
	rc, out := runGate(t, shimEnv("AILANG v0.33.0"))
	if rc != 0 || !strings.Contains(out, "verify gate PASSED") {
		t.Fatalf("delegating release arm: rc=%d\n%s", rc, out)
	}
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
		rc, out := runGate(t, shimEnv(arm.version))
		n := strings.Count(out, notice)
		quietCounts[arm.name] = n
		if rc != 0 || n != 0 {
			t.Fatalf("%s: expected rc=0 with 0 notices, got rc=%d n=%d\n%s", arm.name, rc, n, out)
		}
	}

	rcB, outB := runGate(t, shimEnv("AILANG v0.34.0"))
	countB := strings.Count(outB, notice)
	if rcB != 0 || countB != 1 {
		t.Fatalf("ArmUnrecognised: expected rc=0 with exactly 1 notice, got rc=%d n=%d\n%s", rcB, countB, outB)
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

func TestFixtureDiscrimination(t *testing.T) {
	requirePinned(t)
	env := shimEnv("AILANG v0.33.0")
	delete(env, "WORLD_PKG_AILANG_BIN")
	rc, out := runGate(t, env)
	if rc != 1 || !strings.Contains(out, "wrong compiler version") ||
		strings.Contains(out, "[DEV_BUILD]") || strings.Contains(out, "[NOT_A_RELEASE]") {
		t.Fatalf("fixture discrimination: rc=%d\n%s", rc, out)
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
