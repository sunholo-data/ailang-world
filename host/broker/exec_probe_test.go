package broker

// Row 140 M3 (w-workspace-exec-toolchain-effect §4.6, §6 AC3.1, §7): the
// startup probe against srt SHIMS, a fake @anthropic-ai/sandbox-runtime whose
// cli.js World's own node runs. The pass-through shim runs the command
// UNSANDBOXED (M0's control, out/macos-passthrough-control.txt), so every
// confinement arm must fail by name; the no-op shim runs nothing
// (MUT-PROBE-SKIPPED); the zero shim reports every death as 0 (srt's own
// TERM -> 0 hazard, V11). These tests need a real node (the probe client is
// World's node) and skip, by name, only when none is found.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// shimCLIs are the cli.js bodies of the srt shims. Each takes srt's launch
// line (--settings S -- argv…) exactly as World passes it.
var shimCLIs = map[string]string{
	"passthrough": `const {spawnSync}=require('child_process');const os=require('os');
const a=process.argv.slice(2);
if(a[0]!=='--settings'||a[2]!=='--'){console.error('SHIM: bad launch '+JSON.stringify(a));process.exit(99)}
const r=spawnSync(a[3],a.slice(4),{stdio:'inherit'});
if(r.error){console.error('SHIM: '+r.error.code);process.exit(127)}
process.exit(r.status===null?128+os.constants.signals[r.signal]:r.status);
`,
	"noop": "process.exit(0);\n",
	"zero": `const {spawnSync}=require('child_process');
const a=process.argv.slice(2);
spawnSync(a[3],a.slice(4),{stdio:'inherit'});
process.exit(0);
`,
}

// probeNode is the real node the probe client and the shims run on.
func probeNode(t *testing.T) string {
	t.Helper()
	node := os.Getenv(execNodeEnv)
	if node == "" {
		var err error
		if node, err = exec.LookPath("node"); err != nil {
			t.Skipf("SKIP %s: node absent: %s is unset and no node is on PATH (the probe client is World's node)", t.Name(), execNodeEnv)
		}
	}
	real, err := filepath.EvalSymlinks(node)
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(real)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// shimSrtTree is a node_modules holding a fake srt of the given kind.
func shimSrtTree(t *testing.T, kind string) string {
	t.Helper()
	cli, ok := shimCLIs[kind]
	if !ok {
		t.Fatalf("no shim %q", kind)
	}
	return fakeSrtTree(t, "0.0.0-test", cli) // fakePin reads 0.0.0-test
}

type probeFixture struct {
	paths   ExecHostPaths
	sandbox *ExecSandbox
	profile *ExecProfile
	node    string
}

// newProbeFixture archives tree (fake pin unless pin is given) under a
// scratch state dir, with a scratch HOME and workspace root, and parses a
// minimal profile whose probe argv is probe.
func newProbeFixture(t *testing.T, tree string, pin *ExecSandboxPin, home string, probe []any) *probeFixture {
	t.Helper()
	node := probeNode(t)
	base := canonicalTempDir(t)
	f := &probeFixture{node: node, paths: ExecHostPaths{OperatorHome: home, StateDir: filepath.Join(base, "state"),
		WorkspaceRoot: filepath.Join(base, "workspace")}}
	if f.paths.OperatorHome == "" {
		f.paths.OperatorHome = filepath.Join(base, "home")
	}
	mustMkdir(t, f.paths.OperatorHome, f.paths.StateDir, f.paths.WorkspaceRoot, filepath.Join(f.paths.WorkspaceRoot, "ep1"))
	p := fakePin(t, tree)
	if pin != nil {
		p = *pin
	}
	archiveStart := time.Now()
	sb, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: tree, Node: node, StateDir: f.paths.StateDir, Pin: p})
	if err != nil {
		t.Fatalf("ArchiveExecSandbox: %v", err)
	}
	t.Logf("TIMING ArchiveExecSandbox took %s", time.Since(archiveStart))
	f.sandbox = sb
	m := map[string]any{"profile": ExecProfileVersion, "project": "proj", "path": []any{filepath.Dir(node), "/usr/bin", "/bin"},
		"timeout_ms": 9000, "probe": probe,
		"commands": map[string]any{"true": map[string]any{"argv": []any{"sh", "-c", "exit 0"}}}}
	f.profile = mustParseProfile(t, m)
	return f
}

func mustParseProfile(t *testing.T, m map[string]any) *ExecProfile {
	t.Helper()
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	p, err := ParseExecProfile(data)
	if err != nil {
		t.Fatalf("ParseExecProfile: %v", err)
	}
	return p
}

func (f *probeFixture) probe(t *testing.T, client string) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return ProbeExecProfiles(ctx, ExecProbeConfig{Sandbox: f.sandbox, Paths: f.paths, Client: client}, []*ExecProfile{f.profile})
}

// failedArms returns the arm names an *ExecProbeError reports.
func failedArms(t *testing.T, err error) []string {
	t.Helper()
	var perr *ExecProbeError
	if !errors.As(err, &perr) {
		t.Fatalf("probe error = %v (%T), want an *ExecProbeError", err, err)
	}
	if perr.Project != "proj" {
		t.Fatalf("probe error names project %q", perr.Project)
	}
	return perr.FailedArms()
}

// assertProbeCleaned: the probe leaves no decoy, scratch episode, settings
// file or exec cache behind, pass or fail.
func (f *probeFixture) assertProbeCleaned(t *testing.T) {
	t.Helper()
	for _, d := range ExecProbeDecoys(f.paths) {
		if _, err := os.Lstat(d); err == nil {
			t.Errorf("decoy %s left behind", d)
		}
	}
	for _, dir := range []string{f.paths.WorkspaceRoot, filepath.Join(f.paths.StateDir, "exec"), filepath.Join(f.paths.StateDir, "exec-cache")} {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "wep-") || e.Name() == ".exec-probe-sibling" {
				t.Errorf("probe scratch %s left in %s", e.Name(), dir)
			}
		}
	}
}

var execProbeArms = []string{ExecProbeArmWriteInside, ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmDefaultWrite,
	ExecProbeArmExitStatus, ExecProbeArmNetwork, ExecProbeArmToolchain}

func TestExecProbeArmNames(t *testing.T) {
	want := []string{"arm1-write-inside", "arm2-write-outside", "arm3-read-fence", "arm4-srt-default-write",
		"arm5-exit-status", "arm6-network", "arm7-toolchain"}
	if !slices.Equal(execProbeArms, want) {
		t.Fatalf("arm names = %v, want %v", execProbeArms, want)
	}
}

// AC3.1, MUT-PROBE-NET-VACUOUS: an UNSANDBOXED srt must be refused by every
// confinement arm — arm 6 included, not only arms 1–3 — while the arms a
// pass-through cannot break (write inside, exit status, toolchain) pass.
func TestExecProbePassThroughShimFailsEveryConfinementArm(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	err := f.probe(t, "")
	got := failedArms(t, err)
	t.Logf("pass-through: %v", err)
	want := []string{ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmDefaultWrite, ExecProbeArmNetwork}
	if !slices.Equal(got, want) {
		t.Fatalf("pass-through srt failed arms %v, want exactly %v", got, want)
	}
	// MUT-PROBE-NET-VACUOUS: each of arm 6's legs catches the unsandboxed
	// run on its own — the live counted listener (raw and via the proxy
	// environment) as well as the reverse leg.
	var perr *ExecProbeError
	errors.As(err, &perr)
	for _, f := range perr.Failed {
		if f.Arm != ExecProbeArmNetwork {
			continue
		}
		for _, leg := range []string{"a raw connect to World's live listener: token \"WORLD-PROBE-CONNECTED\", accepts 2",
			"a connect through srt's proxy environment: token \"WORLD-PROBE-CONNECTED\", accepts 3", "reverse leg: the host reached"} {
			if !strings.Contains(f.Why, leg) {
				t.Errorf("arm 6 under a pass-through srt does not report %q:\n%s", leg, f.Why)
			}
		}
	}
	f.assertProbeCleaned(t)
}

// MUT-PROBE-SKIPPED (unit half): a no-op srt runs nothing, so no arm can
// observe its positive token.
func TestExecProbeNoopShimFailsArm1(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "noop"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	got := failedArms(t, f.probe(t, ""))
	want := []string{ExecProbeArmWriteInside, ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmDefaultWrite,
		ExecProbeArmExitStatus, ExecProbeArmNetwork}
	if !slices.Equal(got, want) {
		t.Fatalf("no-op srt failed arms %v, want %v", got, want)
	}
	f.assertProbeCleaned(t)
}

// Arm 5: an srt that reports a signal death as 0 (bare srt's TERM -> 0 on
// macOS, V11) is refused by name.
func TestExecProbeZeroShimFailsArm5(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "zero"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	if got := failedArms(t, f.probe(t, "")); !slices.Contains(got, ExecProbeArmExitStatus) {
		t.Fatalf("zero srt failed arms %v, want %s among them", got, ExecProbeArmExitStatus)
	}
}

// Arm 7: a toolchain probe that exits non-zero is refused by name.
func TestExecProbeFailingToolchainFailsArm7(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "echo no-toolchain >&2; exit 3"})
	err := f.probe(t, "")
	if got := failedArms(t, err); !slices.Contains(got, ExecProbeArmToolchain) {
		t.Fatalf("failed arms %v, want %s", got, ExecProbeArmToolchain)
	}
	if !strings.Contains(err.Error(), "no-toolchain") {
		t.Fatalf("arm 7's refusal does not carry the toolchain's own words: %v", err)
	}
}

// MUT-PROBE-RC-ONLY: a missing probe client exits 127 with no token. That
// must FAIL arms 2, 3 and 6 — a non-zero rc is not a refusal by the sandbox.
func TestExecProbeMissingClientFailsTheRefusalArms(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	missing := filepath.Join(canonicalTempDir(t), "no-such-node")
	got := failedArms(t, f.probe(t, missing))
	for _, arm := range []string{ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmNetwork} {
		if !slices.Contains(got, arm) {
			t.Errorf("a missing probe client passed %s (failed arms %v): an rc without a token was read as a refusal", arm, got)
		}
	}
	f.assertProbeCleaned(t)
}

// arm8Fixture uses shell in place of both node and srt: no optional executables,
// sockets, external controls, or tool pins participate in the arm's verdict.
func arm8Fixture(t *testing.T, body string, patterns []string) *execProber {
	t.Helper()
	t.Setenv(execNodeEnv, "/bin/sh")
	f := newProbeFixture(t, fakeSrtTree(t, "0.0.0-test", "shift 3\nexec \"$@\"\n"), nil, "", []any{"sh", "-c", "exit 0"})
	f.profile.Protected = patterns
	pp := f.arm8Prober(t)
	client := filepath.Join(f.paths.OperatorHome, "client")
	if err := os.WriteFile(client, []byte("#!/bin/sh\nfor target; do :; done\n"+body), 0700); err != nil {
		t.Fatal(err)
	}
	pp.client = client
	return pp
}

func (f *probeFixture) arm8Prober(t *testing.T) *execProber {
	t.Helper()
	root := filepath.Join(f.paths.WorkspaceRoot, "ep1")
	h, err := NewExecHandler(ExecHandlerConfig{Profile: f.profile, Sandbox: f.sandbox, Episode: "ep1", Worktree: root,
		WorkspaceRoot: f.paths.WorkspaceRoot, StateDir: f.paths.StateDir, OperatorHome: f.paths.OperatorHome})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"home", "tmp", "pycache", "npm"} {
		mustMkdir(t, filepath.Join(h.cache, name))
	}
	return &execProber{h: h, cwd: root, client: f.node}
}

// Exercise startup's production entry point, so directly testing probeProtected
// cannot hide a disconnected arm. Like the other shim startup tests, this needs
// Node (probeNode skips by name when absent), but no real srt or tool pin.
func TestExecProbeArm8Wired(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "exit 0"})
	for _, row := range []struct {
		name      string
		protected []string
		total     string
	}{
		{"protected", []string{"tests/**"}, "of 8 arms"},
		{"without-protected", nil, "of 7 arms"},
	} {
		t.Run(row.name, func(t *testing.T) {
			f.profile.Protected = row.protected
			err := f.probe(t, "")
			got := failedArms(t, err)
			wantArm8 := len(row.protected) > 0
			if slices.Contains(got, ExecProbeArmProtectedWrite) != wantArm8 {
				t.Fatalf("startup failed arms %v: arm8-protected-write present must be %t", got, wantArm8)
			}
			if !strings.Contains(err.Error(), row.total) {
				t.Fatalf("startup error = %q, want %q", err, row.total)
			}
			t.Logf("startup failed arms %v; %s", got, row.total)
			f.assertProbeCleaned(t)
		})
	}
}

func TestExecProbeArm8FailsUnsandboxed(t *testing.T) {
	for _, row := range []struct {
		name, body, goos, why string
		patterns              []string
		wantPass              bool
	}{
		{"pass-through", "echo changed > \"$target\"; echo WORLD-PROBE-WROTE", "darwin", `token "WORLD-PROBE-WROTE"`, []string{"tests/**"}, false},
		{"missing-client", "", "darwin", `token ""`, []string{"tests/**"}, false},
		{"forged-token", "echo changed > \"$target\"; echo WORLD-PROBE-REFUSED EACCES", "darwin", "host unchanged=false", []string{"tests/**"}, false},
		{"linux-detect", "echo WORLD-PROBE-WROTE", "linux", "linux detection must report intact=false", []string{"spec?.go"}, false},
		{"linux-no-token", `echo changed > "$target"`, "linux", "missing success token", []string{"spec?.go"}, false},
		{"linux-checker-error", `echo changed > "$target"
echo WORLD-PROBE-WROTE
root=${target%/*/*}
printf '[filter "probe"]\n\tclean = cat\n' >> "$root/.git/config"`, "linux", "linux detection must report intact=false", []string{"spec?.go"}, false},
		{"linux-directory-write", `echo changed > "$target"; echo WORLD-PROBE-WROTE`, "linux", "protected modification", []string{"tests/**", "spec?.go"}, false},
		{"linux-ebusy-passes", `case "$target" in
*/tests/*) echo WORLD-PROBE-REFUSED EBUSY ;;
*) echo changed > "$target"; echo WORLD-PROBE-WROTE ;;
esac`, "linux", "", []string{"tests/**", "spec?.go"}, true},
		{"linux-exdev-passes", `case "$target" in
*/tests/*) echo WORLD-PROBE-REFUSED EXDEV ;;
*) echo changed > "$target"; echo WORLD-PROBE-WROTE ;;
esac`, "linux", "", []string{"tests/**", "spec?.go"}, true},
		{"darwin-create", `case "$target" in
*/new/*) echo changed > "$target"; echo WORLD-PROBE-WROTE ;;
*) echo WORLD-PROBE-REFUSED EPERM ;;
esac`, "darwin", "protected creation", []string{"spec?.go"}, false},
		{"darwin-create-forged", `case "$target" in
*/new/*) echo changed > "$target" ;;
esac
echo WORLD-PROBE-REFUSED EPERM`, "darwin", "protected creation", []string{"spec?.go"}, false},
		{"darwin-refuses-all", "echo WORLD-PROBE-REFUSED EPERM", "darwin", "", []string{"spec?.go"}, true},
	} {
		t.Run(row.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			// execFileRefusal consults the host, independently of probeProtected's goos.
			if (row.name == "linux-ebusy-passes" || row.name == "linux-exdev-passes") && hostGOOS() != "linux" {
				t.Skip("execFileRefusal uses the host GOOS; EBUSY/EXDEV rows run on linux hosts (CI)")
			}
			pp := arm8Fixture(t, row.body, row.patterns)
			if row.name == "missing-client" {
				pp.client += "-missing"
			}
			pp.probeProtected(ctx, row.goos)
			if row.wantPass {
				if len(pp.failed) != 0 {
					t.Fatalf("refusing both protected modification and creation must pass: %v", pp.failed)
				}
				return
			}
			err := &ExecProbeError{Project: "proj", Failed: pp.failed}
			if !slices.Contains(failedArms(t, err), ExecProbeArmProtectedWrite) {
				t.Fatalf("%s: arm8-protected-write must fail: %s", row.name, row.why)
			}
			for _, f := range err.Failed {
				if f.Arm == ExecProbeArmProtectedWrite && !strings.Contains(f.Why, row.why) {
					t.Fatalf("arm8 Why=%q, want %q", f.Why, row.why)
				}
				if f.Arm == ExecProbeArmProtectedWrite {
					var detail string
					switch row.name {
					case "linux-checker-error":
						detail = WorktreeFilterRefused
					case "darwin-create-forged":
						detail = "absent=false"
					}
					if !strings.Contains(f.Why, detail) {
						t.Fatalf("arm8 Why=%q, want %q", f.Why, detail)
					}
				}
			}
		})
	}
}

func TestExecProbeArm8Absent(t *testing.T) {
	for _, patterns := range [][]string{nil, {}} {
		pp := arm8Fixture(t, "exit 99", patterns)
		want := []string{"arm1-write-inside", "arm2-write-outside", "arm3-read-fence", "arm4-srt-default-write", "arm5-exit-status", "arm6-network", "arm7-toolchain"}
		if got := execProbeOrder(pp.h.profile); !slices.Equal(got, want) {
			t.Fatalf("legacy seven-arm order: got %v, want %v; no arm8 for absent set", got, want)
		}
		pp.probeProtected(context.Background(), hostGOOS())
		if pp.h.Spawns() != 0 || len(pp.failed) != 0 || !absent(filepath.Join(pp.h.worktree, ".git")) {
			t.Fatal("empty protection invoked arm8")
		}
	}
}

func TestExecProbeProtectedSample(t *testing.T) {
	for name, pattern := range map[string]string{"literal": "spec.go", "star": "*.spec", "question": "spec?.go", "class": "spec[0-9]?.go", "class-star": "spec[*].go", "class-question": "spec[?].go", "unicode": "[éΩ]?.GO", "negated": "[^a-z].spec", "mixed-case": "SpEc?.GO", "directory": "Mixed/tests/**", "nested": "sub/spec/**"} {
		t.Run(name, func(t *testing.T) {
			sample, err := execProbeProtectedSample(pattern)
			if err != nil || !matchProtected([]string{pattern}, sample) {
				t.Fatalf("probe specimen must match declared pattern %q: sample=%q error=%v", pattern, sample, err)
			}
		})
	}
	if _, err := execProbeProtectedSample("src/*.go"); err == nil {
		t.Fatal("unsupported construction must return explicit error")
	}
}
