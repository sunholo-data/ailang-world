package broker

// Row 140 M2 against the REAL sandbox (w-workspace-exec-toolchain-effect §6
// AC2.2, AC2.3, AC2.4, AC2.5, AC2.7, AC2.8, AC2.10): ExecHandler runs the
// pinned srt from World's archive, launched exactly as §4.5 says. These tests
// need an srt install and node, and skip — by name, saying why — only when
// either is absent:
//
//	WORLD_EXEC_SRT_NODE_MODULES  a node_modules dir holding
//	                             @anthropic-ai/sandbox-runtime@0.0.78 (the
//	                             pin; its cli.js is checked against pin.json)
//	WORLD_EXEC_NODE              the node that runs srt (default: node on PATH)
//
// CI's go job installs both (the "Row 140 srt" step) and a verbose step fails
// unless every one of these tests prints --- PASS. On linux a read-denied
// region is a masked tmpfs (V45): a read there is ENOENT and a write may land
// in the throwaway mount, so the property every write arm asserts is "the
// host bytes are unchanged", plus a positive refusal token on darwin.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	execSrtEnv  = "WORLD_EXEC_SRT_NODE_MODULES"
	execNodeEnv = "WORLD_EXEC_NODE"
)

// realSrt returns the srt install and node, or skips the test by name.
func realSrt(t *testing.T) (nodeModules, node string) {
	t.Helper()
	nodeModules = os.Getenv(execSrtEnv)
	if nodeModules == "" {
		t.Skipf("SKIP %s: real srt absent: %s is unset (point it at a node_modules holding %s@%s)",
			t.Name(), execSrtEnv, ExecSandboxPackage, ExecSandboxRelease)
	}
	node = os.Getenv(execNodeEnv)
	if node == "" {
		var err error
		if node, err = exec.LookPath("node"); err != nil {
			t.Skipf("SKIP %s: node absent: %s is unset and no node is on PATH", t.Name(), execNodeEnv)
		}
	}
	abs, err := filepath.Abs(node)
	if err != nil {
		t.Fatal(err)
	}
	return nodeModules, abs
}

func realHome(t *testing.T) string {
	t.Helper()
	h, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(h)
	if err != nil {
		t.Fatal(err)
	}
	return real
}

// newSrtFixture is an exec fixture on the real, pinned srt with the operator
// HOME the real one (srt anchors its own default write paths there).
func newSrtFixture(t *testing.T, profile func(m map[string]any, f *execFixture), setup func(f *execFixture)) *execFixture {
	t.Helper()
	nm, node := realSrt(t)
	pin := ExecSandboxPin{Version: ExecSandboxRelease, CLISHA256: ExecSandboxCLISHA256}
	return newExecFixture(t, execFixtureOpts{node: node, tree: nm, pin: &pin, home: realHome(t), setup: setup,
		profileF: func(m map[string]any, f *execFixture) {
			m["path"] = []any{filepath.Dir(node), "/usr/bin", "/bin"}
			if profile != nil {
				profile(m, f)
			}
		}})
}

func addCommands(m map[string]any, cmds map[string]any) {
	c := m["commands"].(map[string]any)
	for k, v := range cmds {
		c[k] = v
	}
}

// Probe commands run World's own node inside the sandbox and print a
// positive token (§4.6 r2): WORLD-PROBE-WROTE / -READ, or -REFUSED <errno>.
func writeProbe(target string) map[string]any {
	return map[string]any{"argv": []any{"node", "-e", fmt.Sprintf(
		"try{require('fs').writeFileSync(%q,'x');console.log('WORLD-PROBE-WROTE')}catch(e){console.log('WORLD-PROBE-REFUSED '+e.code)}", target)}}
}

func readProbe(target string) map[string]any {
	return map[string]any{"argv": []any{"node", "-e", fmt.Sprintf(
		"try{require('fs').readFileSync(%q);console.log('WORLD-PROBE-READ')}catch(e){console.log('WORLD-PROBE-REFUSED '+e.code)}", target)}}
}

func probeToken(r execResult) string {
	for _, l := range strings.Split(r.Stdout, "\n") {
		if strings.HasPrefix(l, "WORLD-PROBE-") {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

var darwinRefusals = []string{"WORLD-PROBE-REFUSED EPERM", "WORLD-PROBE-REFUSED EACCES", "WORLD-PROBE-REFUSED EROFS"}

func isRefusal(tok string) bool {
	return strings.HasPrefix(tok, "WORLD-PROBE-REFUSED ")
}

// AC2.2: a real `go test` through the sandbox, from a seeded GOCACHE.
// AC2.8 (first half): the seed is byte-unchanged after a run that writes to
// GOCACHE, and the seed itself is not writable from inside (MUT-SEED-WRITABLE).
func TestExecSrtGoTest(t *testing.T) {
	_, _ = realSrt(t)
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("go is not on PATH: %v", err)
	}
	goBin, _ = filepath.Abs(goBin)
	out, err := exec.Command(goBin, "env", "GOROOT").Output()
	if err != nil {
		t.Fatal(err)
	}
	goroot, err := filepath.EvalSymlinks(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatal(err)
	}
	seed := filepath.Join(canonicalTempDir(t), "seeds", "go-build")
	goEnv := map[string]any{"GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "GOPROXY": "off", "GOTELEMETRY": "off", "GOFLAGS": "-mod=mod"}
	f := newSrtFixture(t, func(m map[string]any, _ *execFixture) {
		m["path"] = []any{filepath.Dir(goBin), "/usr/bin", "/bin"}
		m["env"] = goEnv
		m["read_roots"] = []any{goroot}
		m["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": seed}}
		addCommands(m, map[string]any{
			"go-test": map[string]any{"argv": []any{"go", "test", "-count=1"}, "flags": map[string]any{"-run": "regex", "-v": "bool"},
				"positional": "pkgpattern", "max_args": 8},
			"gocache-write": map[string]any{"argv": []any{"/bin/sh", "-c", `echo probe > "$GOCACHE/world-m2-probe" && echo WORLD-PROBE-WROTE`}},
			"w-seed":        map[string]any{"argv": []any{"/bin/sh", "-c", "echo x > " + shellQuote(filepath.Join(seed, "poison")) + " && echo WORLD-PROBE-WROTE || echo WORLD-PROBE-REFUSED"}},
		})
	}, func(f *execFixture) {
		files := map[string]string{
			"go.mod":          "module scratch\n\ngo 1.21\n",
			"ok/ok.go":        "package ok\n\nfunc Two() int { return 2 }\n",
			"ok/ok_test.go":   "package ok\n\nimport \"testing\"\n\nfunc TestTwo(t *testing.T) {\n\tif Two() != 2 {\n\t\tt.Fatal(\"two\")\n\t}\n}\n",
			"bad/bad_test.go": "package bad\n\nimport \"testing\"\n\nfunc TestBad(t *testing.T) { t.Fatal(\"WORLD-M2-EXPECTED-FAILURE\") }\n",
		}
		for rel, body := range files {
			p := filepath.Join(f.worktree, rel)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		// The operator prewarms the seed (V21): same toolchain, same env.
		if err := os.MkdirAll(seed, 0o755); err != nil {
			t.Fatal(err)
		}
		warm := exec.Command(goBin, "test", "-count=1", "./ok/", "./bad/")
		warm.Dir = f.worktree
		warm.Env = []string{"PATH=" + filepath.Dir(goBin) + ":/usr/bin:/bin", "HOME=" + filepath.Join(f.base, "warmhome"),
			"GOCACHE=" + seed, "TMPDIR=" + os.TempDir()}
		for k, v := range goEnv {
			warm.Env = append(warm.Env, k+"="+v.(string))
		}
		_, _ = warm.CombinedOutput() // ./bad fails by design
	})
	seedBefore, err := execManifest(seed)
	if err != nil || len(seedBefore) == 0 {
		t.Fatalf("seed manifest: %d entries, %v", len(seedBefore), err)
	}

	r := f.mustRun(t, "go-test", "./ok/")
	if r.ExitCode == nil || *r.ExitCode != 0 || !strings.Contains(r.Stdout, "ok") {
		t.Fatalf("go test ./ok/: exit %v timed_out %t\nstdout %s\nstderr %s", r.code(), r.TimedOut, r.Stdout, r.Stderr)
	}
	if !reflect.DeepEqual(r.Argv, []string{"go", "test", "-count=1", "./ok/"}) {
		t.Fatalf("argv = %q", r.Argv)
	}
	r = f.mustRun(t, "go-test", "-run", "TestBad", "./bad/")
	if r.ExitCode == nil || *r.ExitCode == 0 || !strings.Contains(r.Stdout+r.Stderr, "WORLD-M2-EXPECTED-FAILURE") {
		t.Fatalf("go test ./bad/: exit %v, want non-zero with the failure text\nstdout %s\nstderr %s", r.code(), r.Stdout, r.Stderr)
	}
	if r := f.mustRun(t, "gocache-write"); probeToken(r) != "WORLD-PROBE-WROTE" {
		t.Fatalf("the run could not write its GOCACHE: %+v", r)
	}
	if _, err := os.Stat(filepath.Join(f.cache, "GOCACHE", "world-m2-probe")); err != nil {
		t.Fatalf("the write did not land in the episode's GOCACHE clone: %v", err)
	}
	r = f.mustRun(t, "w-seed")
	if _, err := os.Stat(filepath.Join(seed, "poison")); err == nil || probeToken(r) == "WORLD-PROBE-WROTE" {
		t.Fatalf("the seed was writable from inside the sandbox (%s): it must never be in allowWrite", probeToken(r))
	}
	seedAfter, err := execManifest(seed)
	if err != nil || !reflect.DeepEqual(seedBefore, seedAfter) {
		t.Fatalf("the seed changed during the runs (%v)", err)
	}
}

// AC2.3: the trampoline maps a signal death to 128+n (MUT-NO-TRAMPOLINE,
// MUT-TRAMPOLINE-EXEC; bare srt reports a SIGTERM death as 0 on macOS, V11).
// MUT-SHELL-JOIN: argv reaches the command verbatim through srt.
func TestExecSrtExitStatusAndArgv(t *testing.T) {
	f := newSrtFixture(t, nil, nil)
	if r := f.mustRun(t, "term"); r.ExitCode == nil || *r.ExitCode != 143 {
		t.Fatalf("SIGTERM self-kill through srt -> %v, want 143 (stderr %q)", r.code(), r.Stderr)
	}
	if r := f.mustRun(t, "exit"); r.ExitCode == nil || *r.ExitCode != 3 {
		t.Fatalf("exit 3 through srt -> %v (stderr %q)", r.code(), r.Stderr)
	}
	r := f.mustRun(t, "test", "-k", "a b", "$(id)", "x;y", "X=Y")
	if r.Stdout != "[-k][a b][$(id)][x;y][X=Y]" {
		t.Fatalf("argv through srt = %q (stderr %q)", r.Stdout, r.Stderr)
	}
}

// AC2.4 on the real stack: every refusal is answered with nothing spawned.
func TestExecSrtRefusalsSpawnNothing(t *testing.T) {
	f := newSrtFixture(t, nil, nil)
	for _, args := range [][]string{{"-exec=sh"}, {"--"}, {"-c"}, {"-s=x.json"}, {"../x"}, {"/abs"}, {"-run=a\nb"}, {"--settings"}} {
		out, err := f.call(t, "test", args...)
		if err != nil || !strings.Contains(string(out), `"refused"`) {
			t.Fatalf("%q -> %s, %v; want a refusal", args, out, err)
		}
	}
	if out, err := f.call(t, "deploy"); err != nil || !strings.Contains(string(out), `"refused"`) {
		t.Fatalf("an unprofiled command -> %s, %v", out, err)
	}
	if n := f.h.spawns.Load(); n != 0 {
		t.Fatalf("exec spawn counter = %d, want 0", n)
	}
}

// fileState is a path's bytes, or "<absent>".
func fileState(p string) string {
	data, err := os.ReadFile(p)
	if err != nil {
		return "<absent>"
	}
	return string(data)
}

// AC2.5: the confinement matrix through the handler. Writes: a sibling
// episode, the operator HOME, /tmp/claude (MUT-NO-DENYWRITE-DEFAULTS), the
// worktree's .git pointer file and .GITHUB/x (MUT-NO-WORLD-DENY; .GITHUB is
// darwin-gating only, V49). Reads: a decoy under the operator HOME
// (MUT-READ-FENCE-OFF) and a sibling episode's file (MUT-ROOT-NOT-DENIED).
// Every write arm asserts the host bytes unchanged; darwin also demands a
// refusal token; every read arm demands a refusal token on both platforms.
func TestExecSrtConfinementMatrix(t *testing.T) {
	home := realHome(t)
	tag := fmt.Sprintf("%d-%d", os.Getpid(), time.Now().UnixNano())
	decoy := filepath.Join(home, ".world-m2-exec-decoy-"+tag)
	homeWrite := filepath.Join(home, ".world-m2-exec-home-"+tag)
	tmpClaude := filepath.Join("/tmp/claude", "world-m2-"+tag)
	if err := os.WriteFile(decoy, []byte("DECOY-NOT-A-SECRET\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(decoy); _ = os.Remove(homeWrite); _ = os.Remove(tmpClaude) })
	// /tmp/claude must EXIST, or a refusal there is a vacuous ENOENT (V49).
	if _, err := os.Stat("/tmp/claude"); os.IsNotExist(err) {
		if err := os.MkdirAll("/tmp/claude", 0o777); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Remove("/tmp/claude") })
	}
	type arm struct {
		name, target string
		write        bool
		darwinOnly   bool
	}
	var arms []arm
	f := newSrtFixture(t, func(m map[string]any, f *execFixture) {
		arms = []arm{
			{"w-sibling", filepath.Join(f.sibling, "pwn.txt"), true, false},
			{"w-home", homeWrite, true, false},
			{"w-tmpclaude", tmpClaude, true, false},
			{"w-gitptr", filepath.Join(f.worktree, ".git"), true, false},
			{"w-github-upper", filepath.Join(f.worktree, ".GITHUB", "x"), true, true},
			{"r-decoy", decoy, false, false},
			{"r-sibling", filepath.Join(f.sibling, "notes.txt"), false, false},
		}
		cmds := map[string]any{"w-inside": writeProbe(filepath.Join(f.worktree, "in.txt"))}
		for _, a := range arms {
			if a.write {
				cmds[a.name] = writeProbe(a.target)
			} else {
				cmds[a.name] = readProbe(a.target)
			}
		}
		addCommands(m, cmds)
	}, func(f *execFixture) {
		if err := os.WriteFile(filepath.Join(f.worktree, ".git"), []byte("gitdir: /nonexistent/.git/worktrees/ep1\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		// .github exists, so on a case-insensitive volume .GITHUB is it.
		if err := os.MkdirAll(filepath.Join(f.worktree, ".github"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(f.sibling, "notes.txt"), []byte("SIBLING-EPISODE-FILE\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	// Controls: inside is writable; the targets are reachable unsandboxed.
	if r := f.mustRun(t, "w-inside"); probeToken(r) != "WORLD-PROBE-WROTE" || fileState(filepath.Join(f.worktree, "in.txt")) != "x" {
		t.Fatalf("control: the worktree is not writable through the sandbox: %+v", r)
	}
	for _, p := range []string{decoy, filepath.Join(f.sibling, "notes.txt")} {
		if fileState(p) == "<absent>" {
			t.Fatalf("control: %s is not readable unsandboxed", p)
		}
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			before := fileState(a.target)
			r := f.mustRun(t, a.name)
			tok, after := probeToken(r), fileState(a.target)
			t.Logf("%s %s: token %q, host %q -> %q", runtime.GOOS, a.target, tok, before, after)
			if tok == "" {
				t.Fatalf("no probe token: the probe did not run (stdout %q stderr %q)", r.Stdout, r.Stderr)
			}
			if !a.write {
				if !isRefusal(tok) {
					t.Fatalf("read of %s through the sandbox: %s, want a refusal", a.target, tok)
				}
				return
			}
			if a.darwinOnly && runtime.GOOS != "darwin" {
				t.Logf("not gating on %s: a case-sensitive volume makes .GITHUB a different, non-special dir (V49)", runtime.GOOS)
				return
			}
			if before != after {
				t.Fatalf("write to %s changed the host bytes %q -> %q (token %s)", a.target, before, after, tok)
			}
			if runtime.GOOS == "darwin" && !slicesContains(darwinRefusals, tok) {
				t.Fatalf("write to %s: %s, want a refusal (EPERM)", a.target, tok)
			}
		})
	}
}

// AC2.5's .git/HEAD arm: an episode whose .git is a DIRECTORY. srt's own
// mandatory deny covers .git/config and .git/hooks but not .git/HEAD (V5),
// so only World's rendered deny closes it (MUT-NO-WORLD-DENY).
func TestExecSrtGitHeadIsNotWritable(t *testing.T) {
	var head string
	f := newSrtFixture(t, func(m map[string]any, f *execFixture) {
		head = filepath.Join(f.worktree, ".git", "HEAD")
		addCommands(m, map[string]any{"w-githead": writeProbe(head)})
	}, func(f *execFixture) {
		if err := os.MkdirAll(filepath.Join(f.worktree, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(f.worktree, ".git", "HEAD"), []byte("ref: refs/heads/dev\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	r := f.mustRun(t, "w-githead")
	tok := probeToken(r)
	if fileState(head) != "ref: refs/heads/dev\n" || !isRefusal(tok) {
		t.Fatalf(".git/HEAD through the sandbox: token %q, bytes %q", tok, fileState(head))
	}
}

// AC2.7: a command sleeping past timeout_ms is timed_out with its partial
// output, no member of its group is left, and the call returns inside the
// budget (MUT-NO-GROUP-KILL).
func TestExecSrtTimeoutKillsTheGroup(t *testing.T) {
	f := newSrtFixture(t, func(m map[string]any, _ *execFixture) {
		m["timeout_ms"] = 2500
		addCommands(m, map[string]any{"sleep": map[string]any{"argv": []any{"/bin/sh", "-c", "sleep 3017 & sleep 3017 & echo started; wait"}}})
	}, nil)
	start := time.Now()
	r := f.mustRun(t, "sleep")
	elapsed := time.Since(start)
	if !r.TimedOut || r.ExitCode != nil || !strings.Contains(r.Stdout, "started") {
		t.Fatalf("result = %+v, want timed_out with the partial head", r)
	}
	if elapsed > 2500*time.Millisecond+3*time.Second {
		t.Fatalf("returned after %s against a 2.5 s timeout_ms", elapsed)
	}
	if n := countProcs(t, "sleep 3017"); n != 0 {
		t.Fatalf("%d sleepers alive after the deadline: the group kill did not reach the sandbox tree", n)
	}
}

// AC2.8 (second half): PYTHONPYCACHEPREFIX keeps __pycache__ out of the
// worktree — the bytecode lands in the episode cache instead (V24).
func TestExecSrtPycacheStaysOutOfTheWorktree(t *testing.T) {
	_, _ = realSrt(t)
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Fatalf("python3 is not on PATH: %v", err)
	}
	py, _ = filepath.Abs(py)
	out, err := exec.Command(py, "-c", "import sys; print(sys.base_prefix)").Output()
	if err != nil {
		t.Fatal(err)
	}
	prefix, err := filepath.EvalSymlinks(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatal(err)
	}
	f := newSrtFixture(t, func(m map[string]any, _ *execFixture) {
		m["path"] = []any{filepath.Dir(py), "/usr/bin", "/bin"}
		m["read_roots"] = []any{prefix}
		addCommands(m, map[string]any{"py": map[string]any{"argv": []any{"python3", "-c", "import zzmod; print(zzmod.X)"}}})
	}, func(f *execFixture) {
		if err := os.WriteFile(filepath.Join(f.worktree, "zzmod.py"), []byte("X = 'ZZ-IMPORTED'\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	r := f.mustRun(t, "py")
	if r.ExitCode == nil || *r.ExitCode != 0 || strings.TrimSpace(r.Stdout) != "ZZ-IMPORTED" {
		t.Fatalf("python import through srt: exit %v stdout %q stderr %q", r.code(), r.Stdout, r.Stderr)
	}
	var inWorktree []string
	_ = filepath.Walk(f.worktree, func(p string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && info.Name() == "__pycache__" {
			inWorktree = append(inWorktree, p)
		}
		return nil
	})
	if len(inWorktree) != 0 {
		t.Fatalf("__pycache__ written into the worktree: %v", inWorktree)
	}
	var pyc []string
	_ = filepath.Walk(filepath.Join(f.cache, "pycache"), func(p string, info os.FileInfo, err error) error {
		if err == nil && strings.HasSuffix(p, ".pyc") {
			pyc = append(pyc, p)
		}
		return nil
	})
	if len(pyc) == 0 {
		t.Fatal("no bytecode in the episode's pycache: the arm is vacuous (python wrote none at all)")
	}
}

// AC2.10: a profile NODE_OPTIONS=--require=./zz.cjs must never reach the
// host-side node that runs srt (V40: it ran worktree code on the host and
// wrote outside every fence). zz.cjs runs only INSIDE the sandbox, where its
// write to a HOME-like scratch path outside allowWrite fails; the child's env
// is exactly the rendered K=V set (MUT-HOST-ENV).
func TestExecSrtProfileEnvNeverReachesTheHostNode(t *testing.T) {
	var marker string
	f := newSrtFixture(t, func(m map[string]any, f *execFixture) {
		marker = filepath.Join(f.base, "homelike", "zz-marker")
		m["env"] = map[string]any{"NODE_OPTIONS": "--require=./zz.cjs"}
		addCommands(m, map[string]any{"hello": map[string]any{"argv": []any{"node", "-e", "console.log('CHILD-RAN')"}}})
	}, func(f *execFixture) {
		if err := os.MkdirAll(filepath.Dir(marker), 0o755); err != nil {
			t.Fatal(err)
		}
		zz := fmt.Sprintf("try{require('fs').writeFileSync(%q,'escaped');console.error('ZZ-RAN wrote')}catch(e){console.error('ZZ-RAN '+e.code)}\n", marker)
		if err := os.WriteFile(filepath.Join(f.worktree, "zz.cjs"), []byte(zz), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	r := f.mustRun(t, "hello")
	if strings.TrimSpace(r.Stdout) != "CHILD-RAN" || !strings.Contains(r.Stderr, "ZZ-RAN") {
		t.Fatalf("stdout %q stderr %q: the child node did not run with NODE_OPTIONS (the arm would be vacuous)", r.Stdout, r.Stderr)
	}
	if strings.Contains(r.Stderr, "ZZ-RAN wrote") || fileState(marker) != "<absent>" {
		t.Fatalf("zz.cjs wrote %s (stderr %q): NODE_OPTIONS ran outside the sandbox", marker, r.Stderr)
	}
	r = f.mustRun(t, "env")
	child := strings.Split(strings.TrimSpace(r.Stdout), "\n")
	sort.Strings(child)
	want := []string{"HOME=" + f.cache + "/home", "TMPDIR=" + f.cache + "/tmp", "PATH=" + filepath.Dir(f.sandbox.Node) + ":/usr/bin:/bin",
		"LANG=C.UTF-8", "LC_ALL=C.UTF-8", "PYTHONPYCACHEPREFIX=" + f.cache + "/pycache", "npm_config_cache=" + f.cache + "/npm",
		"CI=1", "NO_COLOR=1", "NODE_OPTIONS=--require=./zz.cjs"}
	sort.Strings(want)
	if !reflect.DeepEqual(child, want) {
		t.Fatalf("child env through srt = %q\nwant exactly %q (no srt proxy vars, no TMPDIR=/tmp/claude)", child, want)
	}
}
