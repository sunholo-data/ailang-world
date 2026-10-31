package main

// Row 231 M1 (design_docs/planned/w-marketing-position-checker.md §2.2–§2.4, §3 M1, §4 P1–P17):
// `marketingcheck position` on the V64/V65 corpora and ten in-test arms, its parity with the
// banked oracle, its exit codes, and the deadline-enforcing subprocess runner both parity tests use.

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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

type labelledPath struct{ Label, Path string }

var positionStores = []string{"r1b-seed3", "r1b-seed4", "r1c-seed5", "r1c-seed6", "eval-p1", "eval-p2", "eval-p3"}

func posIdx(sh ShotRects, name string) int {
	for i, c := range sh.Content {
		if c.Name == name {
			return i
		}
	}
	panic(name)
}

// positionNewArms is the Go port of design_docs/verification/world-iter328/v13-isoarms-main.go.txt,
// applied to stores/r1c-seed5.json.
var positionNewArms = []struct {
	Name string
	Mut  func(d *Rectangles)
}{
	{"iso-a1-tested-w-3", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "roster-tested")].W -= 3 }},
	{"iso-freegap-validation-h+240", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "validation")].H += 240 }},
	{"iso-minh-approvals-h-1", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "approvals-chain")].H = 1 }},
	{"iso-nonpara-graph-svg-h+24", func(d *Rectangles) { d.Shots["graph.png"].Content[0].H += 24 }},
	{"iso-nonpara-payload-h+24", func(d *Rectangles) { d.Shots["provenance.png"].Content[0].H += 24 }},
	{"acc-tested-w-2", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "roster-tested")].W -= 2 }},
	{"perm-approvals-y+48", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "approvals-chain")].Y += 48 }},
	{"perm-approvals-y-48", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "approvals-chain")].Y -= 48 }},
	{"perm-validation-h+216", func(d *Rectangles) { sh := d.Shots["packet.png"]; sh.Content[posIdx(sh, "validation")].H += 216 }},
	{"perm-roster-drift-2px", func(d *Rectangles) {
		sh := d.Shots["packet.png"]
		shift := 0
		for _, n := range []string{"roster-tested", "roster-attested", "roster-claimed", "roster-proven"} {
			i := posIdx(sh, n)
			shift += 2
			sh.Content[i].Y += shift
			sh.Content[i].H += 2
			shift += 2
		}
	}},
}

// positionNewArmSHA256 is the banked generator's output (V13), so the port is byte-pinned.
var positionNewArmSHA256 = map[string]string{
	"acc-tested-w-2":               "d471e623c1fdc088789a91b912f42d82b8fc8861a8582e5bfa4a2e13a6fa99cc",
	"iso-a1-tested-w-3":            "4b1bb2420e4d611f14c1e5abf8b066f63cf061cbf466111ba316b1a3fb96ee45",
	"iso-freegap-validation-h+240": "8574e69cca8cd546e3653897232b13cca06214f017480d37a522c1f7a65ec69e",
	"iso-minh-approvals-h-1":       "4b6f9498077e64d94e544cc3b95e370a3c5df6a597f85cb32b90a4941092f930",
	"iso-nonpara-graph-svg-h+24":   "5b384a4aec473be76b1e62a32a4f0e0a1a8bbe10e0e87c07842566af5ec312c4",
	"iso-nonpara-payload-h+24":     "6da98f40004129e8c868812f385c7c8f4b65ee2d15ec0788d61a1a66d09e247b",
	"perm-approvals-y+48":          "f96fb392cfa657ebeea35641274de832f8428230839751dd28722a6c8c2ccba7",
	"perm-approvals-y-48":          "bc7a9fb2512d52418967d29e5619d901928bb8a621e6d817191fed1222aac08d",
	"perm-roster-drift-2px":        "7e9a35bba5a59c36168113a971e4058ea061d6b9a938d232190e47785af73a70",
	"perm-validation-h+216":        "6e4f9772de397ee47f7e3b953f1c69fdb044ee368cf15bb05060f1f15a2cd8bc",
}

// writePositionNewArms writes the ten in-test arms into dir and returns name -> path.
func writePositionNewArms(t *testing.T, dir string) map[string]string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join(positionTestdata, "stores", "r1c-seed5.json"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, a := range positionNewArms {
		var d Rectangles
		if err := json.Unmarshal(src, &d); err != nil {
			t.Fatal(err)
		}
		a.Mut(&d)
		b, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(dir, a.Name+".json")
		if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		out[a.Name] = p
	}
	return out
}

// positionInputs returns the 44 labelled inputs: V64 (11), V65 (19 + 4), new (10).
func positionInputs(t *testing.T) []labelledPath {
	t.Helper()
	var in []labelledPath
	for _, s := range []string{"r1-seed1", "r1-seed2", "r1-seed1-blockconverted", "r1-seed2-blockconverted"} {
		in = append(in, labelledPath{"stores/" + s, filepath.Join(positionTestdata, "stores", s+".json")})
	}
	for _, s := range positionStores {
		in = append(in, labelledPath{"stores/" + s, filepath.Join(positionTestdata, "stores", s+".json")})
	}
	for _, d := range []string{"v65", "v65-p1"} {
		files, err := filepath.Glob(filepath.Join(positionTestdata, d, "*.json"))
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(files)
		for _, f := range files {
			in = append(in, labelledPath{d + "/" + strings.TrimSuffix(filepath.Base(f), ".json"), f})
		}
	}
	arms := writePositionNewArms(t, t.TempDir())
	for _, a := range positionNewArms {
		in = append(in, labelledPath{"new/" + a.Name, arms[a.Name]})
	}
	if len(in) != 44 {
		t.Fatalf("position inputs: %d, want 44", len(in))
	}
	return in
}

func loadPositionBaseline(t *testing.T) Rectangles {
	t.Helper()
	base, err := LoadRectangles(positionBaseline)
	if err != nil {
		t.Fatal(err)
	}
	return base
}

// positionTag is a reason's clause tag: (a1)…(a4), (b), missing or background.
func positionTag(reason string) string {
	_, rest, _ := strings.Cut(reason, ": ")
	switch {
	case strings.HasPrefix(rest, "("):
		return rest[:strings.Index(rest, ")")+1]
	case strings.Contains(reason, "missing"):
		return "missing"
	case strings.Contains(rest, "edge >2px"):
		return "background"
	}
	return "?" + reason
}

// positionVerdict is "ACCEPT" or "REFUSE <tag>[,<tag>...]" with sorted unique tags.
func positionVerdict(reasons []string) string {
	if len(reasons) == 0 {
		return "ACCEPT"
	}
	set := map[string]bool{}
	for _, r := range reasons {
		set[positionTag(r)] = true
	}
	var tags []string
	for k := range set {
		tags = append(tags, k)
	}
	sort.Strings(tags)
	return "REFUSE " + strings.Join(tags, ",")
}

// positionWant is the 44-row verdict table (design §3 M1; V10, V13): 19 ACCEPT.
var positionWant = map[string]string{
	"stores/r1-seed1":                            "REFUSE (a1),(a2),(a3),(a4)",
	"stores/r1-seed2":                            "REFUSE (a1),(a2),(a3),(a4)",
	"stores/r1-seed1-blockconverted":             "ACCEPT",
	"stores/r1-seed2-blockconverted":             "ACCEPT",
	"stores/r1b-seed3":                           "ACCEPT",
	"stores/r1b-seed4":                           "ACCEPT",
	"stores/r1c-seed5":                           "ACCEPT",
	"stores/r1c-seed6":                           "ACCEPT",
	"stores/eval-p1":                             "ACCEPT",
	"stores/eval-p2":                             "ACCEPT",
	"stores/eval-p3":                             "ACCEPT",
	"v65/neg-attested+14":                        "REFUSE (a3)",
	"v65/neg-attested-dropped":                   "REFUSE missing",
	"v65/neg-background-over-content":            "REFUSE background",
	"v65/neg-dark-timeline-first+24":             "REFUSE (a3)",
	"v65/neg-graph-svg-at-timeline":              "REFUSE (a2),(a4)",
	"v65/neg-light-header+40":                    "REFUSE (a2),(a3)",
	"v65/neg-light-header-whole-view":            "REFUSE (a1),(a2),(a3),(a4),(b)",
	"v65/neg-packet-all+14":                      "REFUSE (a2)",
	"v65/neg-packet-stack-down-6-lines":          "REFUSE (a2),(b)",
	"v65/neg-payload+14":                         "REFUSE (a2)",
	"v65/neg-proven-grows-6-lines-below-fold":    "REFUSE (b)",
	"v65/neg-rows-reordered-tested-claimed":      "REFUSE (a3)",
	"v65/neg-tested-64x16-glyphs":                "REFUSE (a1),(a3),(a4)",
	"v65/neg-tested-at-attested-row":             "REFUSE (a3)",
	"v65/neg-tested-at-proven-row":               "REFUSE (a3)",
	"v65/neg-tested-height+14":                   "REFUSE (a4)",
	"v65/syn-pos-createdAt-wraps":                "ACCEPT",
	"v65/syn-pos-proven-5-lines":                 "ACCEPT",
	"v65/syn-pos-tested-wraps-3-lines":           "ACCEPT",
	"v65-p1/neg-packet-all+14":                   "REFUSE (a2)",
	"v65-p1/neg-proven-grows-6-lines-below-fold": "ACCEPT",
	"v65-p1/neg-tested-at-attested-row":          "REFUSE (a3)",
	"v65-p1/syn-pos-tested-wraps-3-lines":        "ACCEPT",
	"new/iso-a1-tested-w-3":                      "REFUSE (a1)",
	"new/iso-freegap-validation-h+240":           "REFUSE (a3)",
	"new/iso-minh-approvals-h-1":                 "REFUSE (a4)",
	"new/iso-nonpara-graph-svg-h+24":             "REFUSE (a4)",
	"new/iso-nonpara-payload-h+24":               "REFUSE (a4)",
	"new/acc-tested-w-2":                         "ACCEPT",
	"new/perm-approvals-y+48":                    "ACCEPT",
	"new/perm-approvals-y-48":                    "ACCEPT",
	"new/perm-validation-h+216":                  "ACCEPT",
	"new/perm-roster-drift-2px":                  "ACCEPT",
}

func TestMarketingPositionFixtures(t *testing.T) {
	sums := map[string]string{}
	f, err := os.Open(filepath.Join(positionTestdata, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			t.Fatalf("SHA256SUMS: malformed line %q", line)
		}
		sums[fields[1]] = fields[0]
	}
	if len(sums) != 34 {
		t.Fatalf("SHA256SUMS lists %d files, want 34", len(sums))
	}
	seen := 0
	err = filepath.WalkDir(positionTestdata, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Base(p) == "SHA256SUMS" {
			return err
		}
		rel, _ := filepath.Rel(positionTestdata, p)
		rel = filepath.ToSlash(rel)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		want, ok := sums[rel]
		if !ok {
			t.Errorf("%s: not listed in SHA256SUMS", rel)
		} else if got := sha256Hex(b); got != want {
			t.Errorf("%s: sha256 %s, SHA256SUMS %s", rel, got, want)
		}
		seen++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen != 34 {
		t.Errorf("testdata/position holds %d files besides SHA256SUMS, want 34", seen)
	}
	inputs, err := os.ReadFile(filepath.Join(oracleDir, "v64-inputs-sha256.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range append([]string{"r1-seed1", "r1-seed2"}, positionStores...) {
		want := sums["stores/"+s+".json"] + "  /tmp/world-iter312/" + s + "/rects.json"
		if !strings.Contains(string(inputs), want+"\n") {
			t.Errorf("stores/%s.json: sha256 %s is not the store's rects.json line in v64-inputs-sha256.txt", s, sums["stores/"+s+".json"])
		}
	}
	for _, s := range []string{"r1-seed1", "r1-seed2"} {
		got, err1 := os.ReadFile(filepath.Join(positionTestdata, "stores", s+"-blockconverted.json"))
		want, err2 := os.ReadFile(filepath.Join(oracleDir, "v64-rects-"+s+"-blockconverted.json"))
		if err1 != nil || err2 != nil || !bytes.Equal(got, want) {
			t.Errorf("stores/%s-blockconverted.json differs from the banked file (%v %v)", s, err1, err2)
		}
	}
	base, err := os.ReadFile(positionBaseline)
	if err != nil {
		t.Fatal(err)
	}
	if got := sha256Hex(base); got != positionBaselinePin {
		t.Errorf("baseline r1-rects.json sha256 %s, want %s", got, positionBaselinePin)
	}
	arms := writePositionNewArms(t, t.TempDir())
	for name, p := range arms {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if got := sha256Hex(b); got != positionNewArmSHA256[name] {
			t.Errorf("new arm %s: sha256 %s, banked generator (V13) %s", name, got, positionNewArmSHA256[name])
		}
	}
}

func TestMarketingPositionCorpus(t *testing.T) {
	base := loadPositionBaseline(t)
	accept := 0
	for _, in := range positionInputs(t) {
		t.Run(in.Label, func(t *testing.T) {
			cand, err := loadPositionCandidate(in.Path)
			if err != nil {
				t.Fatal(err)
			}
			reasons, regions := PositionJudge(base, cand)
			if regions != 17 {
				t.Errorf("regions %d, want 17", regions)
			}
			got := positionVerdict(reasons)
			if want, ok := positionWant[in.Label]; !ok || got != want {
				t.Errorf("want %q, got %q (reasons %q)\n\t%q: %q,", want, got, reasons, in.Label, got)
			}
		})
		if positionWant[in.Label] == "ACCEPT" {
			accept++
		}
	}
	if len(positionWant) != 44 || accept != 19 {
		t.Errorf("verdict table: %d rows, %d ACCEPT; want 44 and 19", len(positionWant), accept)
	}
}

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
