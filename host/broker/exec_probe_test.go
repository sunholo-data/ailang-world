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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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

// T1 (design §7): the reverse leg's whole decision for one attempt, §4.2's order.
func TestExecProbeReverseLegJudge(t *testing.T) {
	const bound, eperm, inuse = "WORLD-PROBE-BOUND", "WORLD-PROBE-REFUSED EPERM", "WORLD-PROBE-REFUSED EADDRINUSE"
	child, foreign, silent := reverseSeen{child: true}, reverseSeen{foreign: true}, reverseSeen{silent: true}
	for _, row := range []struct {
		name string
		tok  string
		seen reverseSeen
		want revVerdict
	}{
		{"bound-none", bound, reverseSeen{}, revPass},
		{"eperm-none", eperm, reverseSeen{}, revPass},
		{"bound-foreign", bound, foreign, revInconclusive},
		{"eperm-foreign", eperm, foreign, revInconclusive},
		{"eaddrinuse-none", inuse, reverseSeen{}, revInconclusive},
		{"bound-child", bound, child, revReached},
		{"bound-child-foreign", bound, reverseSeen{child: true, foreign: true}, revReached},
		{"bound-child-silent", bound, reverseSeen{child: true, silent: true}, revReached},
		{"notoken-child", "", child, revReached},
		{"bound-silent", bound, silent, revSilent},
		{"bound-silent-foreign", bound, reverseSeen{silent: true, foreign: true}, revSilent},
		{"eperm-silent", eperm, silent, revSilent},
		{"notoken-none", "", reverseSeen{}, revNoToken},
		{"notoken-foreign", "", foreign, revNoToken},
	} {
		if got := judgeReverse(row.tok, row.seen); got != row.want {
			t.Errorf("%s: judgeReverse(%q, %+v) = %d, want %d", row.name, row.tok, row.seen, got, row.want)
		}
	}
}

// revHolderKinds are the test-owned holders of a reverse-leg port (design
// §6): chatty ones drain World's hello first (closing with unread input
// would reset) and then answer; "silent" accepts and never writes;
// "reset" accepts and closes with SO_LINGER 0.
var revHolderKinds = []string{"http-400", "forged-nonce", "garbage", "clean-close", "silent", "reset"}

func revHolder(t *testing.T, kind string) int {
	t.Helper()
	if !slices.Contains(revHolderKinds, kind) {
		t.Fatalf("revHolder: unknown kind %q", kind)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var held []net.Conn
	t.Cleanup(func() {
		_ = ln.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, c := range held {
			_ = c.Close()
		}
	})
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			switch kind {
			case "silent":
				mu.Lock()
				held = append(held, c)
				mu.Unlock()
				continue
			case "reset":
				_ = c.(*net.TCPConn).SetLinger(0)
				_ = c.Close()
				continue
			}
			go func() {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(2 * time.Second))
				_, _ = bufio.NewReader(c).ReadString('\n') // World's hello
				switch kind {
				case "http-400":
					_, _ = io.WriteString(c, "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n")
				case "forged-nonce":
					_, _ = io.WriteString(c, "WORLD-PROBE-NONCE "+strings.Repeat("0", 32)+"\n")
				case "garbage":
					_, _ = io.WriteString(c, "garbage\n")
				} // clean-close: zero bytes, then FIN
			}()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

// revPicker is the reverse leg's port seam in tests: call n returns
// steps[n-1] ("free" = a port nothing holds, else a revHolder kind) and
// records the port; a call past the list fails, so an unbounded retry
// ends at once instead of hanging.
type revPicker struct {
	t     *testing.T
	steps []string
	calls atomic.Int32
	mu    sync.Mutex
	ports []int
}

func newRevPicker(t *testing.T, steps ...string) *revPicker { return &revPicker{t: t, steps: steps} }

func (p *revPicker) pick() (int, error) {
	n := int(p.calls.Add(1))
	if n > len(p.steps) {
		return 0, fmt.Errorf("test picker exhausted after %d ports", len(p.steps))
	}
	var port int
	if p.steps[n-1] == "free" {
		var err error
		if port, err = freeLoopbackPort(); err != nil {
			return 0, err
		}
	} else {
		port = revHolder(p.t, p.steps[n-1])
	}
	p.mu.Lock()
	p.ports = append(p.ports, port)
	p.mu.Unlock()
	return port, nil
}

func (p *revPicker) port(t *testing.T, i int) int {
	t.Helper()
	p.mu.Lock()
	defer p.mu.Unlock()
	if i >= len(p.ports) {
		t.Fatalf("picker handed out %d ports, want at least %d", len(p.ports), i+1)
	}
	return p.ports[i]
}

// revFixture is arm8Fixture's /bin/sh node, srt and client (V24: no node,
// no real srt): the client answers the proxy op with a refusal and the
// listen op with listenTok ("" = no token), and never binds.
func revFixture(t *testing.T, listenTok string, p *revPicker) *execProber {
	t.Helper()
	body := "for a in \"$@\"; do [ \"$a\" = proxy ] && echo 'WORLD-PROBE-REFUSED ECONNREFUSED'; done\n"
	if listenTok != "" {
		body += "echo '" + listenTok + "'\n"
	}
	pp := arm8Fixture(t, body, nil)
	pp.pickPort = p.pick
	return pp
}

func runReverseLeg(t *testing.T, pp *execProber) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, why := pp.reverseLeg(ctx, 1) // port 1: nothing listens there
	return why
}

const revReachedText = "the host reached the listener the child bound"

// T2 (design §7): a foreign answer on attempt 1 is retried; attempt 2 passes.
func TestExecProbeReverseLegChattyForeignHolderIsRetried(t *testing.T) {
	for _, plat := range []struct{ name, tok string }{{"linux-bound", "WORLD-PROBE-BOUND"}, {"darwin-eperm", "WORLD-PROBE-REFUSED EPERM"}} {
		for _, kind := range []string{"http-400", "forged-nonce", "garbage", "clean-close"} {
			t.Run(plat.name+"-"+kind, func(t *testing.T) {
				p := newRevPicker(t, kind, "free", "free")
				why := runReverseLeg(t, revFixture(t, plat.tok, p))
				if why != "" || p.calls.Load() != 2 {
					t.Fatalf("Why = %q, picker calls %d; want a pass on attempt 2 (calls 2)", why, p.calls.Load())
				}
			})
		}
	}
}

// T3 (design §7): exhaustion refuses distinctly and never as a reach.
func TestExecProbeReverseLegExhaustionRefusesDistinctly(t *testing.T) {
	for _, row := range []struct {
		name, tok string
		steps     []string
	}{
		{"chatty-every-time", "WORLD-PROBE-BOUND", []string{"http-400", "http-400", "http-400"}},
		{"child-eaddrinuse", "WORLD-PROBE-REFUSED EADDRINUSE", []string{"free", "free", "free"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := newRevPicker(t, row.steps...)
			why := runReverseLeg(t, revFixture(t, row.tok, p))
			if !strings.HasPrefix(why, "reverse leg: inconclusive after 3 attempts") || strings.Contains(why, revReachedText) || p.calls.Load() != 3 {
				t.Fatalf("Why = %q, picker calls %d; want inconclusive after 3 attempts, no reach text, calls 3", why, p.calls.Load())
			}
		})
	}
}

// T5 (design §7): a child with no token, or a clean pass, is judged on attempt 1.
func TestExecProbeReverseLegNoTokenIsNotRetried(t *testing.T) {
	for _, row := range []struct {
		name, tok string
		steps     []string
	}{
		{"silent-client", "", []string{"free", "free", "free"}},
		{"clean", "WORLD-PROBE-BOUND", []string{"free", "free", "free"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := newRevPicker(t, row.steps...)
			why := runReverseLeg(t, revFixture(t, row.tok, p))
			if row.name == "clean" {
				if why != "" || p.calls.Load() != 1 {
					t.Fatalf("Why = %q, picker calls %d; want a pass on attempt 1", why, p.calls.Load())
				}
				return
			}
			if !strings.Contains(why, "reverse leg: no token from the child's listener") || p.calls.Load() != 1 {
				t.Fatalf("Why = %q, picker calls %d; want no-token on attempt 1", why, p.calls.Load())
			}
		})
	}
}

// T7 (design §7): a silent holder is a sticky refusal, on any attempt.
func TestExecProbeReverseLegSilentHolderRefusesSticky(t *testing.T) {
	silentText := func(port int) string {
		return fmt.Sprintf("reverse leg: a listener on 127.0.0.1:%d accepted the host's dial and answered nothing", port)
	}
	for _, row := range []struct {
		name  string
		steps []string
		at    int // the attempt index (0-based) that must refuse
		calls int32
	}{
		{"silent-attempt-1", []string{"silent", "free", "free"}, 0, 1},
		{"silent-after-chatty", []string{"http-400", "silent", "free"}, 1, 2},
		{"reset-holder", []string{"reset", "free", "free"}, 0, 1},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := newRevPicker(t, row.steps...)
			why := runReverseLeg(t, revFixture(t, "WORLD-PROBE-BOUND", p))
			if !strings.HasPrefix(why, silentText(p.port(t, row.at))) || !strings.Contains(why, "answered nothing to World's hello") ||
				strings.Contains(why, revReachedText) || p.calls.Load() != row.calls {
				t.Fatalf("Why = %q, picker calls %d; want the sticky silent refusal on attempt %d, calls %d", why, p.calls.Load(), row.at+1, row.calls)
			}
		})
	}
}
