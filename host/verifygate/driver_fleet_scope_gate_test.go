package verifygate

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	checkedSetDisclosure = "checked set: World-tracked paths under tools/launchd and exact file scripts/mission_decisions.sh, plus explicit REQUIRED_FLEET_PATHS"
	requiredPathLine     = "explicit required paths (checked separately): tools/launchd/lib/pin-root.sh"
	emptyRequiredLine    = "explicit required paths (checked separately): (none)"
	residualDisclosure   = "phase-3 residual enumeration only: tools/launchd and exact file scripts/mission_decisions.sh; files outside this boundary are unenumerated (not zero), not certified by phase 3"
)

type fleetScopeFixture struct {
	world string
	fleet string
	home  string
}

func TestDriverFleetResidualScope(t *testing.T) {
	fixture := newFleetScopeFixture(t, true)
	writeFixtureFile(t, fixture.fleet, "scripts/mission_decisions_v2.sh", "#!/bin/sh\nexit 0\n")
	fixture.commitAll(t, fixture.fleet, "outside-boundary addition")

	t.Run("outside_boundary_unenumerated", func(t *testing.T) {
		requireHeadPath(t, fixture, "scripts/mission_decisions_v2.sh", true)
		rc, out := fixture.runGate(t)
		if rc != 0 {
			t.Fatalf("outside-boundary command lost its rc=0 control: rc=%d\n%s", rc, out)
		}
		// This is the primary regression oracle. Keep it before assertions about the
		// revised success line so an old-wording mutant fails for the intended claim.
		if !strings.Contains(out, residualDisclosure) {
			t.Fatalf("missing phase-3 boundary/unenumerated disclosure:\n%s", out)
		}
		if !strings.Contains(out, "3 files match fleet HEAD") {
			t.Fatalf("outside-boundary command lost its three-match control:\n%s", out)
		}
		if strings.Contains(out, "scripts/mission_decisions_v2.sh") {
			t.Fatalf("outside-boundary sibling was individually reported:\n%s", out)
		}
		if strings.Count(out, "unclassified fleet-only path (not tracked, not required):") != 0 {
			t.Fatalf("outside-boundary sibling was counted as a residual:\n%s", out)
		}
		if !strings.Contains(out, checkedSetDisclosure) {
			t.Fatalf("missing checked-set disclosure:\n%s", out)
		}
	})

	t.Run("baseline_disclosure", func(t *testing.T) {
		baseline := newFleetScopeFixture(t, true)
		rc, out := baseline.runGate(t)
		requireFleetScopeSuccess(t, rc, out)
		if !strings.Contains(out, checkedSetDisclosure) {
			t.Fatalf("missing checked-set disclosure:\n%s", out)
		}
		if countTrimmedLine(out, requiredPathLine) != 1 {
			t.Fatalf("required-path disclosure count = %d, want 1:\n%s", countTrimmedLine(out, requiredPathLine), out)
		}
	})

	writeFixtureFile(t, fixture.fleet, "tools/launchd/new-helper.sh", "#!/bin/sh\nexit 0\n")
	fixture.commitAll(t, fixture.fleet, "inside-boundary addition")

	t.Run("inside_boundary_positive_control", func(t *testing.T) {
		requireHeadPath(t, fixture, "scripts/mission_decisions_v2.sh", true)
		requireHeadPath(t, fixture, "tools/launchd/new-helper.sh", true)
		rc, out := fixture.runGate(t)
		if rc != 0 || !strings.Contains(out, "3 files match fleet HEAD") {
			t.Fatalf("inside-boundary command lost its rc=0/three-match control: rc=%d\n%s", rc, out)
		}
		if strings.Contains(out, "scripts/mission_decisions_v2.sh") {
			t.Fatalf("outside-boundary sibling was individually reported:\n%s", out)
		}
		warning := "unclassified fleet-only path (not tracked, not required): tools/launchd/new-helper.sh"
		if strings.Count(out, warning) != 1 {
			t.Fatalf("inside-boundary warning count = %d, want 1:\n%s", strings.Count(out, warning), out)
		}
		summary := "1 unclassified fleet-only paths within the phase-3 boundary not certified (see above)"
		if strings.Count(out, summary) != 1 {
			t.Fatalf("qualified summary count = %d, want 1:\n%s", strings.Count(out, summary), out)
		}
		if !strings.Contains(out, residualDisclosure) {
			t.Fatalf("inside-boundary run lost residual disclosure:\n%s", out)
		}
	})

	t.Run("required_path_missing", func(t *testing.T) {
		missing := newFleetScopeFixture(t, false)
		requireHeadPath(t, missing, "tools/launchd/lib/pin-root.sh", false)
		requireHeadPath(t, missing, "tools/launchd/lib/pin-root.sh", true)
		rc, out := missing.runGate(t)
		if rc != 1 || !strings.Contains(out, "REQUIRED fleet paths MISSING LOCALLY:") ||
			!strings.Contains(out, "tools/launchd/lib/pin-root.sh (REQUIRED by World, absent locally)") {
			t.Fatalf("missing required path was not refused precisely: rc=%d\n%s", rc, out)
		}
		for _, forbidden := range []string{"files match fleet HEAD", "SKIPPED", "AILANG_BIN is unset", "enumerated 0 comparable", "MISSING IN FLEET", "committed copy differs", "accounting mismatch"} {
			if strings.Contains(out, forbidden) {
				t.Fatalf("required-path refusal contained incompatible output %q:\n%s", forbidden, out)
			}
		}
	})

	t.Run("empty_required_list", func(t *testing.T) {
		empty := newFleetScopeFixture(t, true)
		path := filepath.Join(empty.world, "scripts", "verify_go.sh")
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		old := "REQUIRED_FLEET_PATHS=(\n  \"tools/launchd/lib/pin-root.sh\"\n)"
		if bytes.Count(raw, []byte(old)) != 1 {
			t.Fatalf("required-array mutation matched %d blocks, want 1", bytes.Count(raw, []byte(old)))
		}
		if err := os.WriteFile(path, bytes.Replace(raw, []byte(old), []byte("REQUIRED_FLEET_PATHS=()"), 1), 0o755); err != nil {
			t.Fatal(err)
		}
		rc, out := empty.runBash(t, "-n", "scripts/verify_go.sh")
		if rc != 0 {
			t.Fatalf("empty-array verifier failed bash -n: rc=%d\n%s", rc, out)
		}
		rc, out = empty.runGate(t)
		if rc != 0 || !strings.Contains(out, "3 files match fleet HEAD") || countTrimmedLine(out, emptyRequiredLine) != 1 || strings.Contains(out, "unbound variable") {
			t.Fatalf("empty required list was not safely disclosed: rc=%d\n%s", rc, out)
		}
	})
}

func newFleetScopeFixture(t *testing.T, worldHasRequired bool) fleetScopeFixture {
	t.Helper()
	base := t.TempDir()
	f := fleetScopeFixture{
		world: filepath.Join(base, "world"),
		fleet: filepath.Join(base, "fleet"),
		home:  filepath.Join(base, "home"),
	}
	for _, root := range []string{f.world, f.fleet, f.home} {
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(f.home, "empty-gitconfig"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(f.home, "empty-hooks"), 0o755); err != nil {
		t.Fatal(err)
	}
	copyGateFile(t, f.world, "scripts/verify_go.sh", 0o755)
	common := map[string]string{
		"tools/launchd/control.sh":     "#!/bin/sh\nexit 0\n",
		"scripts/mission_decisions.sh": "#!/bin/sh\nexit 0\n",
	}
	for path, body := range common {
		writeFixtureFile(t, f.world, path, body)
		writeFixtureFile(t, f.fleet, path, body)
	}
	writeFixtureFile(t, f.fleet, "tools/launchd/lib/pin-root.sh", "#!/bin/sh\nexit 0\n")
	if worldHasRequired {
		writeFixtureFile(t, f.world, "tools/launchd/lib/pin-root.sh", "#!/bin/sh\nexit 0\n")
	}
	f.initAndCommit(t, f.world)
	f.initAndCommit(t, f.fleet)
	return f
}

func writeFixtureFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func (f fleetScopeFixture) initAndCommit(t *testing.T, root string) {
	t.Helper()
	f.runGit(t, root, "init", "-q")
	f.runGit(t, root, "config", "user.name", "Fleet Scope Test")
	f.runGit(t, root, "config", "user.email", "fleet-scope@example.invalid")
	f.runGit(t, root, "config", "commit.gpgsign", "false")
	f.runGit(t, root, "config", "core.hooksPath", filepath.Join(f.home, "empty-hooks"))
	f.commitAll(t, root, "fixture baseline")
}

func (f fleetScopeFixture) commitAll(t *testing.T, root, message string) {
	t.Helper()
	f.runGit(t, root, "add", "--all")
	f.runGit(t, root, "commit", "-q", "-m", message)
}

func (f fleetScopeFixture) runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	rc, out := runFleetScopeCommand(t, root, f.home, "git", args...)
	if rc != 0 {
		t.Fatalf("git %s failed: rc=%d\n%s", strings.Join(args, " "), rc, out)
	}
	return out
}

func (f fleetScopeFixture) runGate(t *testing.T) (int, string) {
	t.Helper()
	return f.runBash(t, "scripts/verify_go.sh", "--driver-fleet-check")
}

func (f fleetScopeFixture) runBash(t *testing.T, args ...string) (int, string) {
	t.Helper()
	return runFleetScopeCommand(t, f.world, f.home, requireBash(t), args...)
}

func runFleetScopeCommand(t *testing.T, dir, home, name string, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	blocked := map[string]bool{"CI": true, "AILANG_BIN": true, "AILANG_FLEET_REPO": true}
	cmd.Env = make([]string, 0, len(os.Environ())+6)
	for _, item := range os.Environ() {
		key := strings.SplitN(item, "=", 2)[0]
		if !blocked[key] && !strings.HasPrefix(key, "GIT_") && key != "HOME" {
			cmd.Env = append(cmd.Env, item)
		}
	}
	cmd.Env = append(cmd.Env,
		"HOME="+home,
		"AILANG_FLEET_REPO="+filepath.Dir(home)+"/fleet",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+filepath.Join(home, "empty-gitconfig"),
		"GIT_TERMINAL_PROMPT=0",
	)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("bounded command timed out: %s %s", name, strings.Join(args, " "))
	}
	if err == nil {
		return 0, output.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String()
	}
	t.Fatalf("start bounded command %s: %v", name, err)
	return -1, ""
}

func requireHeadPath(t *testing.T, f fleetScopeFixture, rel string, inFleet bool) {
	t.Helper()
	root := f.world
	want := false
	if inFleet {
		root, want = f.fleet, true
	}
	out := f.runGit(t, root, "ls-tree", "-r", "--name-only", "HEAD")
	found := false
	for _, line := range strings.Split(out, "\n") {
		found = found || strings.TrimSpace(line) == rel
	}
	if found != want {
		t.Fatalf("HEAD membership for %s in %s = %t, want %t\n%s", rel, root, found, want, out)
	}
}

func requireFleetScopeSuccess(t *testing.T, rc int, out string) {
	t.Helper()
	if rc != 0 || !strings.Contains(out, "3 files match fleet HEAD") || !strings.Contains(out, residualDisclosure) {
		t.Fatalf("fleet-scope success/disclosure missing: rc=%d\n%s", rc, out)
	}
	for _, forbidden := range []string{"SKIPPED", "AILANG_BIN is unset", "enumerated 0 comparable"} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("successful fleet-scope output contained %q:\n%s", forbidden, out)
		}
	}
}

func countTrimmedLine(out, want string) int {
	count := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == want {
			count++
		}
	}
	return count
}
