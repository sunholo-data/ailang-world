package broker

// Row 140 M3 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.3,
// §4.6): the `serve --exec-*` startup. StartExec turns the operator's flags
// into verified exec configuration or refuses, naming the flag and the
// field: the profile files are read once, the read_roots rules that need the
// host's paths are applied to the literal path and its realpath, the
// confinement stack (srt, node, libnode, profiles, cache seeds) must lie out
// of every agent-writable root, srt is archived and checked against its pin,
// and every profile passes the startup probe (exec_probe.go). Nothing here
// runs an agent command.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ExecMinNode is srt's `engines` floor (§4.6: node ≥ 20.11).
var ExecMinNode = [2]int{20, 11}

// ExecStartupConfig is the `serve --exec-*` flags (§4.6).
type ExecStartupConfig struct {
	ProfileFiles    []string // --exec-profile, repeatable
	EpisodeProjects []string // --exec-episode-project EP=PROJECT, repeatable
	SandboxDir      string   // --exec-sandbox: srt's node_modules (or its package dir)
	Node            string   // --exec-node; "" = node on the operator's PATH
	MaxOutputBytes  int64    // --exec-max-output-bytes; 0 = DefaultExecMaxOutputBytes
	// Paths are the canonical state dir and workspace root; an empty
	// OperatorHome is the canonical os.UserHomeDir.
	Paths ExecHostPaths
	// Pin is the srt release required; the zero value is the production pin.
	Pin ExecSandboxPin
	// ProbeClient overrides the probe client (a test seam; "" = the node).
	ProbeClient string
}

// configured reports whether any exec flag was given.
func (c ExecStartupConfig) configured() bool {
	return len(c.ProfileFiles) > 0 || len(c.EpisodeProjects) > 0 || c.SandboxDir != "" || c.Node != "" || c.MaxOutputBytes != 0
}

// ExecStartup is the verified exec configuration `serve` binds.
type ExecStartup struct {
	Profiles        map[string]*ExecProfile // by project
	Only            *ExecProfile            // the profile, when exactly one is configured
	EpisodeProjects map[string]string       // episode -> project
	Sandbox         *ExecSandbox
	MaxOutputBytes  int64
	OperatorHome    string
}

// ExecStartupError is a startup refusal naming the flag (and its value).
type ExecStartupError struct {
	Flag  string
	Value string
	Err   error
}

func (e *ExecStartupError) Error() string {
	if e.Value == "" {
		return fmt.Sprintf("%s: %v", e.Flag, e.Err)
	}
	return fmt.Sprintf("%s %s: %v", e.Flag, e.Value, e.Err)
}

func (e *ExecStartupError) Unwrap() error { return e.Err }

func startupErr(flag, value, format string, args ...any) error {
	return &ExecStartupError{Flag: flag, Value: value, Err: fmt.Errorf(format, args...)}
}

// execStartupBudget bounds the node version check.
const execNodeVersionBudget = 5 * time.Second

// StartExec verifies the exec flags and probes every profile (§4.3, §4.6).
// It returns nil, nil when no exec flag is set.
func StartExec(ctx context.Context, cfg ExecStartupConfig) (*ExecStartup, error) {
	if !cfg.configured() {
		return nil, nil
	}
	if len(cfg.ProfileFiles) == 0 {
		for _, f := range []struct {
			flag string
			set  bool
		}{{"--exec-sandbox", cfg.SandboxDir != ""}, {"--exec-node", cfg.Node != ""},
			{"--exec-episode-project", len(cfg.EpisodeProjects) > 0}, {"--exec-max-output-bytes", cfg.MaxOutputBytes != 0}} {
			if f.set {
				return nil, startupErr(f.flag, "", "needs --exec-profile (it configures workspace-exec, which is off without a profile)")
			}
		}
	}
	if cfg.SandboxDir == "" {
		return nil, startupErr("--exec-sandbox", "", "is required with --exec-profile: an installed %s@%s (node_modules)",
			ExecSandboxPackage, ExecSandboxRelease)
	}
	maxOut := cfg.MaxOutputBytes
	if maxOut == 0 {
		maxOut = DefaultExecMaxOutputBytes
	}
	if maxOut < 1 {
		return nil, startupErr("--exec-max-output-bytes", strconv.FormatInt(cfg.MaxOutputBytes, 10), "must be at least 1")
	}
	paths := cfg.Paths
	if paths.OperatorHome == "" {
		h, err := os.UserHomeDir()
		if err == nil {
			h, err = filepath.EvalSymlinks(h)
		}
		if err != nil {
			return nil, fmt.Errorf("exec: the operator HOME: %w", err)
		}
		paths.OperatorHome = h
	}
	forbidden := []namedRoot{{"an exec cache", filepath.Join(paths.StateDir, "exec-cache")},
		{"the state dir", paths.StateDir}, {"the workspace root", paths.WorkspaceRoot}}

	node, err := resolveExecNode(ctx, cfg.Node, forbidden)
	if err != nil {
		return nil, err
	}
	nodeModules, err := resolveExecSandboxDir(cfg.SandboxDir, forbidden)
	if err != nil {
		return nil, err
	}
	st := &ExecStartup{Profiles: map[string]*ExecProfile{}, EpisodeProjects: map[string]string{}, MaxOutputBytes: maxOut,
		OperatorHome: paths.OperatorHome}
	var order []*ExecProfile
	for _, file := range cfg.ProfileFiles {
		p, err := loadExecProfileFile(file, paths, forbidden)
		if err != nil {
			return nil, err
		}
		if _, dup := st.Profiles[p.Project]; dup {
			return nil, startupErr("--exec-profile", file, "project %q is configured twice", p.Project)
		}
		st.Profiles[p.Project] = p
		order = append(order, p)
	}
	if len(order) == 1 {
		st.Only = order[0]
	}
	for _, pair := range cfg.EpisodeProjects {
		ep, project, ok := strings.Cut(pair, "=")
		switch {
		case !ok:
			return nil, startupErr("--exec-episode-project", pair, "is not EP=PROJECT")
		case !execEpisodePattern.MatchString(ep):
			return nil, startupErr("--exec-episode-project", pair, "episode %q is outside the episode grammar %s", ep, execEpisodePattern)
		case st.Profiles[project] == nil:
			return nil, startupErr("--exec-episode-project", pair, "names no configured profile (projects: %s)",
				strings.Join(sortedKeys(st.Profiles), ", "))
		case st.EpisodeProjects[ep] != "":
			return nil, startupErr("--exec-episode-project", pair, "maps episode %q a second time", ep)
		}
		st.EpisodeProjects[ep] = project
	}
	pin := cfg.Pin
	if pin == (ExecSandboxPin{}) {
		pin = ExecSandboxPin{Version: ExecSandboxRelease, CLISHA256: ExecSandboxCLISHA256}
	}
	st.Sandbox, err = ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nodeModules, Node: node, StateDir: paths.StateDir, Pin: pin})
	if err != nil {
		return nil, &ExecStartupError{Flag: "--exec-sandbox", Value: cfg.SandboxDir, Err: err}
	}
	if err := ProbeExecProfiles(ctx, ExecProbeConfig{Sandbox: st.Sandbox, Paths: paths, MaxOutputBytes: maxOut,
		Client: cfg.ProbeClient}, order); err != nil {
		return nil, err
	}
	return st, nil
}

// execEpisodePattern is the daemon's episode grammar (§4.3).
var execEpisodePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

type namedRoot struct{ name, path string }

// checkPlacement refuses a confinement-stack path whose realpath lies inside
// an agent-writable root (§4.6 Integrity: placement).
func checkPlacement(flag, value, real string, forbidden []namedRoot) error {
	for _, r := range forbidden {
		if r.path != "" && pathWithin(real, r.path) {
			return startupErr(flag, value, "resolves to %s, inside %s %s: the confinement stack (srt, node, libnode, profiles, "+
				"cache seeds) must lie outside the workspace root, the state dir and every exec cache", real, r.name, r.path)
		}
	}
	return nil
}

// resolveExecNode resolves --exec-node once (default: node on PATH), checks
// its placement (the binary, its directory and its libnode) and its version.
func resolveExecNode(ctx context.Context, flagValue string, forbidden []namedRoot) (string, error) {
	name := flagValue
	if name == "" {
		name = "node"
	}
	found := name
	if !strings.Contains(name, "/") {
		p, err := exec.LookPath(name)
		if err != nil {
			return "", startupErr("--exec-node", flagValue, "node does not resolve on PATH: %v", err)
		}
		found = p
	}
	abs, err := filepath.Abs(found)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return "", startupErr("--exec-node", flagValue, "does not resolve: %v", err)
	}
	if info, err := os.Stat(abs); err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
		return "", startupErr("--exec-node", flagValue, "%s is not an executable file", abs)
	}
	for _, p := range []string{abs, filepath.Dir(abs)} {
		if err := checkPlacement("--exec-node", flagValue, p, forbidden); err != nil {
			return "", err
		}
	}
	libs, err := libnodeDigests(abs)
	if err != nil {
		return "", &ExecStartupError{Flag: "--exec-node", Value: flagValue, Err: err}
	}
	for _, l := range libs {
		real, err := filepath.EvalSymlinks(l.Path)
		if err != nil {
			return "", &ExecStartupError{Flag: "--exec-node", Value: flagValue, Err: err}
		}
		if err := checkPlacement("--exec-node", flagValue, real, forbidden); err != nil {
			return "", err
		}
	}
	vctx, cancel := context.WithTimeout(ctx, execNodeVersionBudget)
	defer cancel()
	out, err := runBounded(vctx, handlerBounds{execTimeout: execNodeVersionBudget}, handlerCommand{
		path: abs, args: []string{"--version"}, dir: "/", env: execHostEnv()})
	if err != nil {
		return "", startupErr("--exec-node", flagValue, "%s --version: %v", abs, err)
	}
	version := strings.TrimSpace(string(out))
	major, minor, ok := parseNodeVersion(version)
	if !ok {
		return "", startupErr("--exec-node", flagValue, "%s --version printed %q, not vMAJOR.MINOR.PATCH", abs, version)
	}
	if major < ExecMinNode[0] || (major == ExecMinNode[0] && minor < ExecMinNode[1]) {
		return "", startupErr("--exec-node", flagValue, "version %s is below %d.%d (srt's engines)", version, ExecMinNode[0], ExecMinNode[1])
	}
	return abs, nil
}

func parseNodeVersion(v string) (int, int, bool) {
	parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
	if !strings.HasPrefix(v, "v") || len(parts) != 3 {
		return 0, 0, false
	}
	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	return major, minor, err1 == nil && err2 == nil
}

// resolveExecSandboxDir accepts srt's node_modules dir, or the package dir
// inside one, and returns the node_modules dir World archives whole (V42).
func resolveExecSandboxDir(dir string, forbidden []namedRoot) (string, error) {
	abs, err := filepath.Abs(dir)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return "", startupErr("--exec-sandbox", dir, "does not resolve: %v", err)
	}
	nm := abs
	if _, err := os.Stat(filepath.Join(abs, "@anthropic-ai", "sandbox-runtime", "package.json")); err != nil {
		var pkg struct{ Name string }
		raw, rerr := os.ReadFile(filepath.Join(abs, "package.json"))
		if rerr != nil || json.Unmarshal(raw, &pkg) != nil || pkg.Name != ExecSandboxPackage ||
			filepath.Base(filepath.Dir(abs)) != "@anthropic-ai" {
			return "", startupErr("--exec-sandbox", dir, "is neither a node_modules holding %s nor that package's directory "+
				"(install it with npm install %s@%s)", ExecSandboxPackage, ExecSandboxPackage, ExecSandboxRelease)
		}
		nm = filepath.Dir(filepath.Dir(abs))
	}
	if err := checkPlacement("--exec-sandbox", dir, nm, forbidden); err != nil {
		return "", err
	}
	return nm, nil
}

// loadExecProfileFile reads one --exec-profile once (§4.6 "Profile"):
// placement, the in-memory loader, the read_roots host rules and each cache
// seed's placement.
func loadExecProfileFile(file string, h ExecHostPaths, forbidden []namedRoot) (*ExecProfile, error) {
	abs, err := filepath.Abs(file)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, startupErr("--exec-profile", file, "does not exist")
		}
		return nil, startupErr("--exec-profile", file, "does not resolve: %v", err)
	}
	if err := checkPlacement("--exec-profile", file, abs, forbidden); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, startupErr("--exec-profile", file, "%v", err)
	}
	p, err := ParseExecProfile(data)
	if err != nil {
		return nil, &ExecStartupError{Flag: "--exec-profile", Value: file, Err: err}
	}
	if err := CheckExecReadRoots(p.ReadRoots, h); err != nil {
		return nil, &ExecStartupError{Flag: "--exec-profile", Value: file, Err: err}
	}
	for _, name := range sortedKeys(p.Caches) {
		seed := p.Caches[name].Seed
		if seed == "" {
			continue
		}
		real, err := filepath.EvalSymlinks(seed)
		if err != nil {
			return nil, startupErr("--exec-profile", file, "exec profile: caches.%s.seed: %q does not resolve: %v", name, seed, err)
		}
		if info, err := os.Stat(real); err != nil || !info.IsDir() {
			return nil, startupErr("--exec-profile", file, "exec profile: caches.%s.seed: %q is not a directory", name, seed)
		}
		if err := checkPlacement("--exec-profile", file+" caches."+name+".seed", real, forbidden); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// ExecHostPaths are the canonical host paths the exec fence is drawn
// around: all three are denied for reading inside the sandbox (§4.4).
type ExecHostPaths struct {
	OperatorHome  string
	StateDir      string
	WorkspaceRoot string
}

// ExecProbeDecoys are the startup probe's three fixed decoys (§4.6 arm 3):
// World plants them unsandboxed before the probe and removes them after;
// they do not depend on any profile.
func ExecProbeDecoys(h ExecHostPaths) []string {
	return []string{
		filepath.Join(h.OperatorHome, ".ailang-worldd-exec-probe-decoy"),
		filepath.Join(h.StateDir, "exec-probe", "decoy"),
		filepath.Join(h.WorkspaceRoot, ".exec-probe-sibling", "decoy"),
	}
}

// pathFold reports whether the host's filesystem compares names without
// case: APFS does (V6), so on darwin every containment test folds case and
// refuses the wider class.
var pathFold = runtime.GOOS == "darwin"

func pathEq(a, b string) bool {
	if pathFold {
		return strings.EqualFold(a, b)
	}
	return a == b
}

// pathWithin reports p equal to or lexically inside dir.
func pathWithin(p, dir string) bool {
	if pathEq(p, dir) || dir == "/" {
		return true
	}
	prefix := strings.TrimSuffix(dir, "/") + "/"
	return len(p) > len(prefix) && pathEq(p[:len(prefix)], prefix)
}

// CheckExecReadRoots applies §4.3's read_roots rules [REVISED r1]: an entry
// is refused when its literal path or its realpath equals or is an ancestor
// of the operator HOME, the state dir or the workspace root; lies inside the
// state dir or the workspace root; equals or contains a probe decoy; or is
// /. srt's own precedence is narrower and platform-specific (V36), and World
// adds each root's realpath to allowRead, so the whole class is refused.
func CheckExecReadRoots(roots []string, h ExecHostPaths) error {
	for _, r := range roots {
		if err := readRootOk(r, r, h); err != nil {
			return err
		}
		real, err := filepath.EvalSymlinks(r)
		if err != nil {
			return profileErr("read_roots", "%q does not resolve: %v", r, err)
		}
		if err := readRootOk(r, real, h); err != nil {
			return err
		}
	}
	return nil
}

func readRootOk(entry, p string, h ExecHostPaths) error {
	via := ""
	if p != entry {
		via = fmt.Sprintf(" (its realpath %q)", p)
	}
	if p == "/" {
		return profileErr("read_roots", "%q%s is \"/\"", entry, via)
	}
	for _, d := range []struct{ name, path string }{
		{"the operator HOME", h.OperatorHome}, {"the state dir", h.StateDir}, {"the workspace root", h.WorkspaceRoot},
	} {
		if pathWithin(d.path, p) {
			return profileErr("read_roots", "%q%s equals or is an ancestor of %s %s, which it would re-open", entry, via, d.name, d.path)
		}
	}
	for _, d := range []struct{ name, path string }{{"the state dir", h.StateDir}, {"the workspace root", h.WorkspaceRoot}} {
		if pathWithin(p, d.path) {
			return profileErr("read_roots", "%q%s lies inside %s %s", entry, via, d.name, d.path)
		}
	}
	for _, decoy := range ExecProbeDecoys(h) {
		if pathWithin(decoy, p) {
			return profileErr("read_roots", "%q%s equals or contains the probe decoy %s", entry, via, decoy)
		}
	}
	return nil
}
