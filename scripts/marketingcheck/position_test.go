package main

// Row 231 M1 (design_docs/planned/w-marketing-position-checker.md §2.2–§2.4, §3 M1, §4 P1–P17):
// `marketingcheck position` on the V64/V65 corpora and ten in-test arms, its parity with the
// banked oracle, its exit codes, and the deadline-enforcing subprocess runner both parity tests use.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// ---- the deadline-enforcing subprocess runner (design §2.2 [r2], §4 P17) ----

const (
	subprocBuildLimit = 120 * time.Second // oracle and CLI builds
	subprocRunLimit   = 30 * time.Second  // one checker run
	subprocDrainLimit = 5 * time.Second   // output draining after exit or kill
)

type subprocResult struct {
	Stdout, Stderr string
	Code           int
}

// boundedRun runs name with args in dir (empty = the package directory) under limit. The child
// starts in its own process group; on expiry the whole group is killed and the error is
// `subprocess-timeout: <command> exceeded <limit>`. Output draining is bounded by
// subprocDrainLimit, and the group is reaped after every run. It never skips or substitutes a
// result: a non-zero exit is returned as Code, anything else as an error.
func boundedRun(limit time.Duration, dir string, env []string, name string, args ...string) (subprocResult, error) {
	what := strings.Join(append([]string{name}, args...), " ")
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = subprocDrainLimit
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	res := subprocResult{Stdout: stdout.String(), Stderr: stderr.String(), Code: -1}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return res, fmt.Errorf("subprocess-timeout: %s exceeded %s", what, limit)
	}
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.Code = ee.ExitCode()
			return res, nil
		}
		return res, fmt.Errorf("subprocess %s: %v", what, err)
	}
	res.Code = 0
	return res, nil
}

// mustRun is boundedRun that fails the test on a timeout or a start/drain error.
func mustRun(t *testing.T, limit time.Duration, dir string, env []string, name string, args ...string) subprocResult {
	t.Helper()
	res, err := boundedRun(limit, dir, env, name, args...)
	if err != nil {
		t.Fatalf("%v\nstdout: %s\nstderr: %s", err, res.Stdout, res.Stderr)
	}
	return res
}

// procAlive reports whether pid is a live (non-zombie) process.
func procAlive(t *testing.T, pid int) bool {
	t.Helper()
	if syscall.Kill(pid, 0) != nil {
		return false
	}
	res, err := boundedRun(10*time.Second, "", nil, "ps", "-o", "stat=", "-p", strconv.Itoa(pid))
	if err != nil {
		t.Fatalf("ps: %v", err)
	}
	stat := strings.TrimSpace(res.Stdout)
	return stat != "" && !strings.HasPrefix(stat, "Z")
}

// TestMarketingSubprocessDeadline is the runner's own control (design §4 P17, AC-4): a fake
// checker that sleeps past the bound, with a grandchild in its process group.
func TestMarketingSubprocessDeadline(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "hang.sh")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nsleep 20 &\necho $! > \"$1\"\nwait\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	pidFile := filepath.Join(dir, "grandchild.pid")
	const limit = 2 * time.Second
	start := time.Now()
	_, err := boundedRun(limit, dir, nil, fake, pidFile)
	elapsed := time.Since(start)
	if err == nil || !strings.Contains(err.Error(), "subprocess-timeout: "+fake) {
		t.Fatalf("want a subprocess-timeout failure naming %s, got %v after %s", fake, err, elapsed)
	}
	if elapsed > limit+subprocDrainLimit+3*time.Second {
		t.Fatalf("the runner returned after %s, bound %s", elapsed, limit+subprocDrainLimit)
	}
	b, rerr := os.ReadFile(pidFile)
	pid, perr := strconv.Atoi(strings.TrimSpace(string(b)))
	if rerr != nil || perr != nil || pid <= 0 {
		t.Fatalf("grandchild pid not recorded: %v %v %q", rerr, perr, b)
	}
	deadline := time.Now().Add(3 * time.Second)
	for procAlive(t, pid) {
		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("grandchild %d survived the process-group kill", pid)
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Logf("subprocess-timeout after %s; grandchild %d reaped", elapsed.Round(time.Millisecond), pid)
}

// ---- the banked oracles and their pins (design §2.2 [r2]) ----

const (
	oracleDir            = "../../design_docs/verification/world-iter315"
	positionOracleSource = "v64-p0check-main.go.txt"
	positionOraclePin    = "0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c"
	phoneOracleSource    = "v66-phonecheck-main.go.txt"
	phoneOraclePin       = "31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c"
	oraclePinRegistry    = "../../design_docs/planned/w-marketing-capture-orchestration.md"
	positionBaseline     = "../../design_docs/verification/world-iter312/r1-rects.json"
	positionBaselinePin  = "ba7c83ce4e31c86a79ac54c2189f00605d05d757186421371681c6d78cc53333"
	positionTestdata     = "testdata/position"
)

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// ---- the 44 position inputs (design §2.3, §3 M1) ----

func TestMarketingPositionCLI(t *testing.T) {
	dir := t.TempDir()
	td := func(rel string) string { return filepath.Join(positionTestdata, filepath.FromSlash(rel)) }
	r1, err := os.ReadFile(positionBaseline)
	if err != nil {
		t.Fatal(err)
	}
	write := func(name string, mut func(m map[string]any)) string {
		var m map[string]any
		if err := json.Unmarshal(r1, &m); err != nil {
			t.Fatal(err)
		}
		mut(m)
		b, _ := json.Marshal(m)
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, b, 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	shots := func(m map[string]any) map[string]any { return m["shots"].(map[string]any) }
	emptyGraph := write("empty-graph.json", func(m map[string]any) { shots(m)["graph.png"].(map[string]any)["content"] = []any{} })
	noRegions := write("no-regions.json", func(m map[string]any) {
		for _, s := range positionShots {
			shots(m)[s].(map[string]any)["content"] = []any{}
		}
	})
	wrongSchema := write("wrong-schema.json", func(m map[string]any) { m["schema"] = "world/demo/capture-rectangles/v0" })
	unparsable := filepath.Join(dir, "unparsable.json")
	if err := os.WriteFile(unparsable, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name   string
		args   []string
		rc     int
		stdout []string // regexps, each must match a line
		stderr string   // substring
	}{
		{"accept-rc", []string{"--baseline", positionBaseline, "accept-seed5=" + td("stores/r1c-seed5.json"), "accept-eval-p1=" + td("stores/eval-p1.json")}, 0,
			[]string{`^accept-seed5 +ACCEPT regions=17$`, `^accept-eval-p1 +ACCEPT regions=17$`, `^TOTAL candidates=2 accept=2 refuse=0$`}, ""},
		{"refuse-rc", []string{"--baseline", positionBaseline, "refuse=" + td("v65/neg-attested+14.json")}, 1,
			[]string{`^refuse +REFUSE regions=17$`, `^    packet\.png/roster-attested: \(a3\) gap after roster-tested 29 vs R1 15$`, `^TOTAL candidates=1 accept=0 refuse=1$`}, ""},
		{"baseline-empty-content", []string{"--baseline", emptyGraph, "a=" + td("stores/r1c-seed5.json")}, 1, nil, "marketingcheck: bad-rectangles graph.png: content rectangles required"},
		{"candidate-zero-regions", []string{"--baseline", positionBaseline, "a=" + noRegions}, 1, nil, "marketingcheck: zero-regions " + noRegions},
		{"candidate-wrong-schema", []string{"--baseline", positionBaseline, "a=" + wrongSchema}, 1, nil, "marketingcheck: bad-rectangles " + wrongSchema + ": wrong schema"},
		{"candidate-unparsable", []string{"--baseline", positionBaseline, "a=" + unparsable}, 1, nil, "marketingcheck: bad-rectangles " + unparsable},
		{"candidate-missing-file", []string{"--baseline", positionBaseline, "a=" + filepath.Join(dir, "absent.json")}, 1, nil, "marketingcheck: bad-rectangles"},
		{"no-candidates", []string{"--baseline", positionBaseline}, 2, nil, "no candidates"},
		{"missing-baseline", []string{"a=" + td("stores/eval-p1.json")}, 2, nil, "--baseline is required"},
		{"bad-token", []string{"--baseline", positionBaseline, td("stores/eval-p1.json")}, 2, nil, "want LABEL=PATH"},
		{"duplicate-label", []string{"--baseline", positionBaseline, "a=" + td("stores/eval-p1.json"), "a=" + td("stores/eval-p2.json")}, 2, nil, "want LABEL=PATH"},
		{"empty-label", []string{"--baseline", positionBaseline, "=" + td("stores/eval-p1.json")}, 2, nil, "want LABEL=PATH"},
		{"empty-path", []string{"--baseline", positionBaseline, "a="}, 2, nil, "want LABEL=PATH"},
		{"bad-flag", []string{"--bogus"}, 2, nil, ""},
		{"help", []string{"--help"}, 0, nil, positionUsage},
	} {
		t.Run(c.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			rc := runPosition(c.args, &out, &errb)
			if rc != c.rc {
				t.Fatalf("rc %d, want %d\nstdout: %s\nstderr: %s", rc, c.rc, out.String(), errb.String())
			}
			for _, re := range c.stdout {
				if !regexp.MustCompile(`(?m)` + re).MatchString(out.String()) {
					t.Errorf("stdout lacks %s:\n%s", re, out.String())
				}
			}
			if !strings.Contains(errb.String(), c.stderr) {
				t.Errorf("stderr lacks %q:\n%s", c.stderr, errb.String())
			}
		})
	}
	t.Run("dispatch", func(t *testing.T) {
		if rc := run([]string{"position"}); rc != 2 {
			t.Fatalf("marketingcheck position with no arguments: rc %d, want 2", rc)
		}
		if !strings.Contains(usageText, "marketingcheck position --baseline RECTS LABEL=RECTS...") {
			t.Fatal("usageText lacks the position line")
		}
	})
}

func TestMarketingPositionGeometryMatchesOracle(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(oracleDir, positionOracleSource))
	if err != nil {
		t.Fatal(err)
	}
	hm := regexp.MustCompile(`(?m)^var heights = map\[string\]int\{(.*)\}$`).FindSubmatch(src)
	om := regexp.MustCompile(`(?m)^var order = \[\]string\{(.*)\}$`).FindSubmatch(src)
	wm := regexp.MustCompile(`c\.X\+c\.W <= ([0-9]+) &&`).FindAllSubmatch(src, -1)
	if hm == nil || om == nil || len(wm) != 1 {
		t.Fatalf("oracle geometry not found: heights %v order %v width matches %d", hm != nil, om != nil, len(wm))
	}
	heights := map[string]int{}
	for _, m := range regexp.MustCompile(`"([a-z-]+\.png)": ([0-9]+)`).FindAllSubmatch(hm[1], -1) {
		h, _ := strconv.Atoi(string(m[2]))
		heights[string(m[1])] = h
	}
	var order []string
	for _, m := range regexp.MustCompile(`"([a-z-]+\.png)"`).FindAllSubmatch(om[1], -1) {
		order = append(order, string(m[1]))
	}
	width, _ := strconv.Atoi(string(wm[0][1]))
	if len(heights) != 5 || len(order) != 5 {
		t.Fatalf("oracle lists %d heights and %d shots, want 5 and 5", len(heights), len(order))
	}
	if len(positionShots) != len(order) {
		t.Fatalf("position judges %d shots, oracle %d", len(positionShots), len(order))
	}
	for i, name := range order {
		if positionShots[i] != name {
			t.Errorf("shot %d: position %s, oracle %s", i, positionShots[i], name)
		}
		s, ok := shotByName(name)
		if !ok {
			t.Errorf("%s: absent from Shots", name)
			continue
		}
		if s.Width != width || s.Height != heights[name] {
			t.Errorf("%s: Shots %dx%d, oracle %dx%d", name, s.Width, s.Height, width, heights[name])
		}
	}
}
