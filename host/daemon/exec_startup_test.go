package daemon

// Row 140 M3 (w-workspace-exec-toolchain-effect §4.3, §4.6, §6 AC3.1, §7):
// the `serve --exec-*` startup refusal table, bound to QUICKSTART §10 row by
// row. Each documented row is driven here through New (what `serve` runs)
// and must refuse startup naming what the row's last column quotes; every
// case here is a documented row and every documented row is a case.
//
// The probe rows drive each arm red with an srt SHIM run by World's real
// node: pass-through (runs the command unsandboxed, M0's control), no-op
// (runs nothing: MUT-PROBE-SKIPPED), zero (a signal death reported as 0). The
// package-only row needs the real pinned srt and skips, by name, only when it
// is absent. The table needs a real node (the probe client) and skips only
// when none is found; CI runs it verbosely and requires --- PASS.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

const (
	execSrtEnv  = "WORLD_EXEC_SRT_NODE_MODULES"
	execNodeEnv = "WORLD_EXEC_NODE"
)

// execTestNode is the real node the shims and the probe client run on.
func execTestNode(t *testing.T) string {
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

// realExecSrt is the pinned srt install, or a skip naming why.
func realExecSrt(t *testing.T) string {
	t.Helper()
	nm := os.Getenv(execSrtEnv)
	if nm == "" {
		t.Skipf("SKIP %s: real srt absent: %s is unset (point it at a node_modules holding %s@%s)",
			t.Name(), execSrtEnv, broker.ExecSandboxPackage, broker.ExecSandboxRelease)
	}
	return nm
}

var execShimCLIs = map[string]string{
	"passthrough": `const {spawnSync}=require('child_process');const os=require('os');
const a=process.argv.slice(2);
if(a[0]!=='--settings'||a[2]!=='--'){console.error('SHIM: bad launch');process.exit(99)}
const r=spawnSync(a[3],a.slice(4),{stdio:'inherit'});
if(r.error){process.exit(127)}
process.exit(r.status===null?128+os.constants.signals[r.signal]:r.status);
`,
	"noop": "process.exit(0);\n",
	"zero": `const {spawnSync}=require('child_process');const a=process.argv.slice(2);
spawnSync(a[3],a.slice(4),{stdio:'inherit'});process.exit(0);
`,
}

// writeSrtTree writes a node_modules holding an srt of version with cli.js
// cli, plus one dependency, and returns it with the pin it matches.
func writeSrtTree(t *testing.T, nm, version, cli string) broker.ExecSandboxPin {
	t.Helper()
	pkg := filepath.Join(nm, "@anthropic-ai", "sandbox-runtime")
	for p, body := range map[string]string{
		filepath.Join(pkg, "package.json"):                                `{"name":"@anthropic-ai/sandbox-runtime","version":"` + version + `"}`,
		filepath.Join(pkg, "dist", "cli.js"):                              cli,
		filepath.Join(pkg, "vendor", "seccomp", "x64", "apply-seccomp"):   "#!/bin/sh\n",
		filepath.Join(pkg, "vendor", "seccomp", "arm64", "apply-seccomp"): "#!/bin/sh\n",
		filepath.Join(nm, "commander", "package.json"):                    `{"name":"commander"}`,
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	sum := sha256.Sum256([]byte(cli))
	return broker.ExecSandboxPin{Version: version, CLISHA256: hex.EncodeToString(sum[:])}
}

// execRig is one startup case: a workspace fixture, a scratch HOME, a real
// node and a pass-through srt by default, and a profile written outside every
// agent-writable root.
type execRig struct {
	t       *testing.T
	f       wsFixture
	home    string
	node    string
	srt     string // the --exec-sandbox dir
	pin     *broker.ExecSandboxPin
	client  string // the probe client override ("" = World's node)
	profile map[string]any
	cfg     Config
}

func newExecRig(t *testing.T) *execRig {
	t.Helper()
	node := execTestNode(t)
	f := newWSFixture(t)
	r := &execRig{t: t, f: f, node: node, home: filepath.Join(f.base, "home")}
	if err := os.MkdirAll(r.home, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", r.home) // the operator HOME the fence denies
	r.srt = filepath.Join(f.outside, "srt", "node_modules")
	pin := writeSrtTree(t, r.srt, "0.0.0-test", execShimCLIs["passthrough"])
	r.pin = &pin
	r.profile = map[string]any{"profile": broker.ExecProfileVersion, "project": "proj",
		"path": []any{filepath.Dir(node), "/usr/bin", "/bin"}, "timeout_ms": 5000, "probe": []any{"sh", "-c", "echo toolchain-ok"},
		"commands": map[string]any{"test": map[string]any{"argv": []any{"sh", "-c", "echo ran"}}}}
	r.cfg = Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease),
		ExecSandbox: r.srt, ExecNode: node}
	return r
}

// writeProfile writes the rig's profile at path and names it --exec-profile.
func (r *execRig) writeProfile(path string) {
	r.t.Helper()
	data, err := json.Marshal(r.profile)
	if err != nil {
		r.t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		r.t.Fatal(err)
	}
	r.cfg.ExecProfiles = append(r.cfg.ExecProfiles, path)
}

// start runs New with the rig's config; the startup hook carries the shim's
// pin and the probe client override, the only two test seams.
func (r *execRig) start() (*Daemon, error) {
	r.t.Helper()
	if len(r.cfg.ExecProfiles) == 0 && r.profile != nil {
		r.writeProfile(filepath.Join(r.f.outside, "profiles", "proj.json"))
	}
	prev := execStartupHook
	execStartupHook = func(c *broker.ExecStartupConfig) {
		if r.pin != nil {
			c.Pin = *r.pin
		}
		c.ProbeClient = r.client
	}
	r.t.Cleanup(func() { execStartupHook = prev })
	return newWSDaemon(r.t, r.cfg)
}

// fakeNodeScript is a node stand-in that prints version (placement and the
// version check run before node ever runs a script).
func fakeNodeScript(t *testing.T, path, version string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho "+version+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// execRefusalCases are the AC3.1 rows by QUICKSTART §10 id.
var execRefusalCases = map[string]func(r *execRig){
	// The flags.
	"no-workspace-tools": func(r *execRig) { r.cfg.WorkspaceRoot, r.cfg.ToolAilangBin = "", "" },
	"sandbox-without-profile": func(r *execRig) {
		r.profile = nil
	},
	"node-without-profile": func(r *execRig) { r.profile, r.cfg.ExecSandbox = nil, "" },
	"episode-without-profile": func(r *execRig) {
		r.profile, r.cfg.ExecSandbox, r.cfg.ExecNode = nil, "", ""
		r.cfg.ExecEpisodeProjects = []string{"ep1=proj"}
	},
	"max-output-without-profile": func(r *execRig) {
		r.profile, r.cfg.ExecSandbox, r.cfg.ExecNode = nil, "", ""
		r.cfg.ExecMaxOutputBytes = 1 << 20
	},
	"no-sandbox":            func(r *execRig) { r.cfg.ExecSandbox = "" },
	"sandbox-not-srt":       func(r *execRig) { r.cfg.ExecSandbox = r.f.outside },
	"sandbox-other-version": func(r *execRig) { r.pin = nil; writeSrtTree(r.t, r.srt, "0.0.77", "// cli\n") },
	"sandbox-other-digest": func(r *execRig) {
		r.pin = nil
		writeSrtTree(r.t, r.srt, broker.ExecSandboxRelease, "// not the pinned cli.js\n")
	},
	"node-absent": func(r *execRig) { r.cfg.ExecNode = filepath.Join(r.f.outside, "no-such-node") },
	"node-too-old": func(r *execRig) {
		r.cfg.ExecNode = fakeNodeScript(r.t, filepath.Join(r.f.outside, "old", "node"), "v18.19.0")
	},
	"max-output-zero":    func(r *execRig) { r.cfg.ExecMaxOutputBytes = -1 },
	"episode-not-pair":   func(r *execRig) { r.cfg.ExecEpisodeProjects = []string{"ep1"} },
	"episode-bad-id":     func(r *execRig) { r.cfg.ExecEpisodeProjects = []string{"../x=proj"} },
	"episode-no-project": func(r *execRig) { r.cfg.ExecEpisodeProjects = []string{"ep1=other"} },
	"episode-twice":      func(r *execRig) { r.cfg.ExecEpisodeProjects = []string{"ep1=proj", "ep1=proj"} },
	"project-twice": func(r *execRig) {
		r.writeProfile(filepath.Join(r.f.outside, "profiles", "a.json"))
		r.writeProfile(filepath.Join(r.f.outside, "profiles", "b.json"))
	},
	"profile-absent": func(r *execRig) { r.cfg.ExecProfiles = []string{filepath.Join(r.f.outside, "nope.json")} },
	// §4.3 load-time refusals.
	"unknown-key":   func(r *execRig) { r.profile["network"] = []any{"example.com"} },
	"root-escapes":  func(r *execRig) { r.profile["root"] = "../other" },
	"argv0-missing": func(r *execRig) { cmdOf(r)["argv"] = []any{"cargo", "test"} },
	"env-world-owned": func(r *execRig) {
		r.profile["env"] = map[string]any{"HTTPS_PROXY": "http://example.com"}
	},
	"env-bad-name":  func(r *execRig) { r.profile["env"] = map[string]any{"A-B": "1"} },
	"timeout-cap":   func(r *execRig) { r.profile["timeout_ms"] = 9001 },
	"probe-missing": func(r *execRig) { r.profile["probe"] = []any{"cargo", "--version"} },
	// §4.3 read_roots (MUT-READROOT-ANCESTOR).
	"readroot-home": func(r *execRig) { r.profile["read_roots"] = []any{r.home} },
	"readroot-home-symlink": func(r *execRig) {
		link := filepath.Join(r.f.outside, "homelink")
		if err := os.Symlink(r.home, link); err != nil {
			r.t.Fatal(err)
		}
		r.profile["read_roots"] = []any{link}
	},
	"readroot-root-ancestor": func(r *execRig) { r.profile["read_roots"] = []any{filepath.Dir(r.f.root)} },
	"readroot-slash":         func(r *execRig) { r.profile["read_roots"] = []any{"/"} },
	"readroot-decoy": func(r *execRig) {
		r.profile["read_roots"] = []any{filepath.Join(r.home, ".ailang-worldd-exec-probe-decoy")}
	},
	"readroot-in-state": func(r *execRig) { r.profile["read_roots"] = []any{r.f.stateDir} },
	// §4.6 placement (MUT-PLACEMENT).
	"sandbox-in-root": func(r *execRig) { r.moveSrt(filepath.Join(r.f.root, "tools", "node_modules")) },
	"sandbox-in-state": func(r *execRig) {
		r.moveSrt(filepath.Join(r.f.stateDir, "tools", "node_modules"))
	},
	"sandbox-in-cache": func(r *execRig) {
		r.moveSrt(filepath.Join(r.f.stateDir, "exec-cache", "ep1", "proj", "node_modules"))
	},
	"node-in-root": func(r *execRig) { r.cfg.ExecNode = r.copyNode(filepath.Join(r.f.root, "bin", "node")) },
	"node-in-state": func(r *execRig) {
		r.cfg.ExecNode = r.copyNode(filepath.Join(r.f.stateDir, "bin", "node"))
	},
	"node-in-cache": func(r *execRig) {
		r.cfg.ExecNode = r.copyNode(filepath.Join(r.f.stateDir, "exec-cache", "ep1", "proj", "node"))
	},
	"profile-in-root":  func(r *execRig) { r.writeProfile(filepath.Join(r.f.root, "proj.json")) },
	"profile-in-state": func(r *execRig) { r.writeProfile(filepath.Join(r.f.stateDir, "proj.json")) },
	"profile-in-cache": func(r *execRig) {
		r.writeProfile(filepath.Join(r.f.stateDir, "exec-cache", "ep1", "proj", "proj.json"))
	},
	"seed-in-root": func(r *execRig) {
		seed := filepath.Join(r.f.root, "seed")
		if err := os.MkdirAll(seed, 0o755); err != nil {
			r.t.Fatal(err)
		}
		r.profile["caches"] = map[string]any{"GOCACHE": map[string]any{"seed": seed}}
	},
	// §4.6 the startup probe, each arm driven red.
	"probe-passthrough": func(*execRig) {}, // the rig's default srt is the pass-through shim
	"probe-noop": func(r *execRig) {
		pin := writeSrtTree(r.t, r.srt, "0.0.0-test", execShimCLIs["noop"])
		r.pin = &pin
	},
	"probe-zero": func(r *execRig) {
		pin := writeSrtTree(r.t, r.srt, "0.0.0-test", execShimCLIs["zero"])
		r.pin = &pin
	},
	"probe-toolchain": func(r *execRig) { r.profile["probe"] = []any{"sh", "-c", "exit 3"} },
	"probe-no-client": func(r *execRig) { r.client = filepath.Join(r.f.outside, "no-such-node") },
	"probe-package-only": func(r *execRig) {
		nm := realExecSrt(r.t)
		pkgOnly := filepath.Join(r.f.outside, "pkg-only", "node_modules", "@anthropic-ai", "sandbox-runtime")
		if err := os.MkdirAll(filepath.Dir(pkgOnly), 0o755); err != nil {
			r.t.Fatal(err)
		}
		if out, err := exec.Command("cp", "-R", filepath.Join(nm, "@anthropic-ai", "sandbox-runtime"), pkgOnly).CombinedOutput(); err != nil {
			r.t.Fatalf("cp: %v %s", err, out)
		}
		r.cfg.ExecSandbox, r.pin = filepath.Dir(filepath.Dir(pkgOnly)), nil
	},
}

func cmdOf(r *execRig) map[string]any {
	return r.profile["commands"].(map[string]any)["test"].(map[string]any)
}

// moveSrt rewrites the rig's srt tree at dst and names it --exec-sandbox.
func (r *execRig) moveSrt(dst string) {
	r.t.Helper()
	pin := writeSrtTree(r.t, dst, "0.0.0-test", execShimCLIs["passthrough"])
	r.cfg.ExecSandbox, r.pin = dst, &pin
}

// copyNode puts a node at dst. Placement is checked before node runs, so a
// stand-in that would pass the version check is enough.
func (r *execRig) copyNode(dst string) string {
	r.t.Helper()
	return fakeNodeScript(r.t, dst, "v22.0.0")
}

var execRefusalRow = regexp.MustCompile("^\\| `([a-z0-9-]+)` \\| (.+) \\| (.+) \\|$")

// quickstartExecRefusals parses QUICKSTART §10's refusal table: id -> the
// backticked texts its last column quotes.
func quickstartExecRefusals(t *testing.T) map[string][]string {
	t.Helper()
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 10. Build and test non-AILANG projects")
	if !ok {
		t.Fatal("QUICKSTART §10 is absent")
	}
	_, table, ok := strings.Cut(section, "<!-- exec-refusals:begin -->")
	table, _, ok2 := strings.Cut(table, "<!-- exec-refusals:end -->")
	if !ok || !ok2 {
		t.Fatal("QUICKSTART §10 has no exec-refusals table")
	}
	rows := map[string][]string{}
	quoted := regexp.MustCompile("`([^`]+)`")
	for _, line := range strings.Split(table, "\n") {
		m := execRefusalRow.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		if _, dup := rows[m[1]]; dup {
			t.Fatalf("QUICKSTART §10 lists %q twice", m[1])
		}
		for _, q := range quoted.FindAllStringSubmatch(m[3], -1) {
			rows[m[1]] = append(rows[m[1]], q[1])
		}
		if len(rows[m[1]]) == 0 {
			t.Fatalf("QUICKSTART §10 row %q quotes nothing serve prints", m[1])
		}
	}
	return rows
}

// TestExecStartupRefusalTable is AC3.1: every row of QUICKSTART §10's table
// refuses `serve` at startup (StageConfig, nothing listening) naming what the
// row quotes, and the table and these cases are the same set.
func TestExecStartupRefusalTable(t *testing.T) {
	rows := quickstartExecRefusals(t)
	var ids []string
	for id := range execRefusalCases {
		if _, ok := rows[id]; !ok {
			t.Errorf("case %q is not a row of QUICKSTART §10's refusal table", id)
		}
		ids = append(ids, id)
	}
	for id := range rows {
		if _, ok := execRefusalCases[id]; !ok {
			t.Errorf("QUICKSTART §10 row %q has no case here", id)
		}
	}
	sort.Strings(ids)
	for _, id := range ids {
		t.Run(id, func(t *testing.T) {
			r := newExecRig(t)
			execRefusalCases[id](r)
			d, err := r.start()
			if err == nil {
				_ = d
				t.Fatalf("serve started; want the refusal %q", rows[id])
			}
			var se *StartupError
			if !errors.As(err, &se) || se.Stage != StageConfig {
				t.Fatalf("err = %v, want a StartupError at %s", err, StageConfig)
			}
			for _, want := range rows[id] {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("refusal %q does not name %q (QUICKSTART §10 row %s)", err, want, id)
				}
			}
			t.Logf("%s: %v", id, err)
			// Fail closed: the probe's scratch and decoys are gone.
			for _, d := range broker.ExecProbeDecoys(broker.ExecHostPaths{OperatorHome: r.home, StateDir: r.f.stateDir,
				WorkspaceRoot: r.f.root}) {
				if _, err := os.Lstat(d); err == nil {
					t.Errorf("decoy %s left behind", d)
				}
			}
		})
	}
}
