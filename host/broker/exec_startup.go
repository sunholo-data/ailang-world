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
	ProfileFiles []string // --exec-profile, repeatable
	// ProfilesDir is --exec-profiles-dir (row 189 M2, D-189-2 A): every
	// <project>.json in it except *.proposed.json is loaded. Unlike a
	// --exec-profile, a directory profile that fails to load or probe is
	// left out (ExecStartup.Excluded) instead of refusing startup.
	ProfilesDir     string
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
	return len(c.ProfileFiles) > 0 || c.ProfilesDir != "" || len(c.EpisodeProjects) > 0 || c.SandboxDir != "" || c.Node != "" || c.MaxOutputBytes != 0
}

// ExecStartup is the verified exec configuration `serve` binds.
type ExecStartup struct {
	Profiles map[string]*ExecProfile // by project
	// Only is the profile when exactly one is configured: loaded, and none
	// left out (a left-out sibling must never make another project's profile
	// every episode's default).
	Only            *ExecProfile
	EpisodeProjects map[string]string // episode -> project
	// Sources is each loaded profile's file (its realpath), by project.
	Sources map[string]string
	// Excluded are the profiles directory's left-out profiles, in file order.
	Excluded       []ExecProfileExclusion
	Sandbox        *ExecSandbox
	Tools          *ExecToolState
	MaxOutputBytes int64
	OperatorHome   string
}

// ExecProfileExclusion is one profile of --exec-profiles-dir that `serve` left
// out: it failed to load (parse, placement, read_roots, cache seeds) or the
// startup probe, or it named a project another profile already configures.
// Project is the file's stem, which a loaded directory profile must equal, so
// a mint or an exec naming it can be refused with Reason.
type ExecProfileExclusion struct {
	File    string
	Project string
	Reason  string
}

// ExecProfileSuffix and ExecProposedSuffix name the profiles directory's
// files: <project>.json is active; <project>.proposed.json (D-189-5, written
// by `profile init`) is inert until the operator renames it.
const (
	ExecProfileSuffix  = ".json"
	ExecProposedSuffix = ".proposed.json"
)

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
	if len(cfg.ProfileFiles) == 0 && cfg.ProfilesDir == "" {
		for _, f := range []struct {
			flag string
			set  bool
		}{{"--exec-sandbox", cfg.SandboxDir != ""}, {"--exec-node", cfg.Node != ""},
			{"--exec-episode-project", len(cfg.EpisodeProjects) > 0}, {"--exec-max-output-bytes", cfg.MaxOutputBytes != 0}} {
			if f.set {
				return nil, startupErr(f.flag, "", "needs --exec-profile or --exec-profiles-dir (it configures workspace-exec, which is off without a profile)")
			}
		}
	}
	if cfg.SandboxDir == "" {
		return nil, startupErr("--exec-sandbox", "", "is required with --exec-profile or --exec-profiles-dir: an installed %s@%s (node_modules)",
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

	configuredNode := cfg.Node
	if configuredNode == "" {
		var err error
		configuredNode, err = execNodePath("")
		if err != nil {
			return nil, err
		}
	}
	node, err := resolveExecNode(ctx, configuredNode, forbidden)
	if err != nil {
		return nil, err
	}
	nodeModules, err := resolveExecSandboxDir(cfg.SandboxDir, forbidden)
	if err != nil {
		return nil, err
	}
	st := &ExecStartup{Profiles: map[string]*ExecProfile{}, EpisodeProjects: map[string]string{}, Sources: map[string]string{},
		MaxOutputBytes: maxOut, OperatorHome: paths.OperatorHome}
	var order []*ExecProfile
	flaggedFiles := map[string]string{} // realpath -> the flag value
	for _, file := range cfg.ProfileFiles {
		p, err := loadExecProfileFile(file, paths, forbidden)
		if err != nil {
			return nil, err
		}
		if _, dup := st.Profiles[p.Project]; dup {
			return nil, startupErr("--exec-profile", file, "project %q is configured twice", p.Project)
		}
		st.Profiles[p.Project] = p
		real := realpathOr(file)
		st.Sources[p.Project] = real
		flaggedFiles[real] = file
		order = append(order, p)
	}
	flagged := len(order)
	if cfg.ProfilesDir != "" {
		loaded, err := loadExecProfilesDir(cfg.ProfilesDir, paths, forbidden, st, flaggedFiles)
		if err != nil {
			return nil, err
		}
		order = append(order, loaded...)
	}
	excludedProject := func(project string) bool {
		for _, x := range st.Excluded {
			if x.Project == project {
				return true
			}
		}
		return false
	}
	for _, pair := range cfg.EpisodeProjects {
		ep, project, ok := strings.Cut(pair, "=")
		switch {
		case !ok:
			return nil, startupErr("--exec-episode-project", pair, "is not EP=PROJECT")
		case !execEpisodePattern.MatchString(ep):
			return nil, startupErr("--exec-episode-project", pair, "episode %q is outside the episode grammar %s", ep, execEpisodePattern)
		case st.Profiles[project] == nil && !excludedProject(project):
			// A left-out directory profile is still a configured project:
			// the mapping stands, and that episode's exec refuses with the
			// reason instead of taking every other episode down with it.
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
	errs, err := ProbeExecProfilesEach(ctx, ExecProbeConfig{Sandbox: st.Sandbox, Paths: paths, MaxOutputBytes: maxOut,
		Client: cfg.ProbeClient}, order)
	if err != nil {
		return nil, err
	}
	if err := st.applyProbe(order, flagged, errs); err != nil {
		return nil, err
	}
	st.Tools, err = NewExecToolState(ExecToolStateConfig{Node: configuredNode, Sandbox: st.Sandbox, Profiles: st.Profiles,
		ProbeConfig: ExecProbeConfig{Paths: paths, MaxOutputBytes: maxOut, Client: cfg.ProbeClient}})
	if err != nil {
		return nil, err
	}
	return st, nil
}

// applyProbe applies the startup probe's per-profile results (errs[i] is
// order[i]'s): a failed --exec-profile (the first flagged of order) refuses
// startup, as before; a failed directory profile is left out with the probe's
// reason and every other profile serves. Only is set last, from what remains.
func (st *ExecStartup) applyProbe(order []*ExecProfile, flagged int, errs []error) error {
	if err := errors.Join(errs[:flagged]...); err != nil {
		return err
	}
	for i := flagged; i < len(order); i++ {
		if errs[i] == nil {
			continue
		}
		p := order[i]
		st.Excluded = append(st.Excluded, ExecProfileExclusion{File: st.Sources[p.Project], Project: p.Project, Reason: errs[i].Error()})
		delete(st.Profiles, p.Project)
		delete(st.Sources, p.Project)
	}
	if len(st.Profiles) == 1 && len(st.Excluded) == 0 {
		for _, p := range st.Profiles {
			st.Only = p
		}
	}
	return nil
}

// realpathOr is path's absolute, symlink-free form, or path itself when it
// does not resolve (the loader has already refused that case).
func realpathOr(path string) string {
	abs, err := filepath.Abs(path)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return path
	}
	return abs
}

// loadExecProfilesDir loads the profiles directory (row 189 M2): every
// <project>.json in it, in name order, except *.proposed.json. The directory
// itself must resolve to a directory outside every agent-writable root, or
// startup refuses. A file that is also a --exec-profile is skipped (it is
// already loaded); any other file that fails to load, whose project is not
// its file stem, or whose project is already configured is left out in
// st.Excluded, and the rest load. It returns the loaded profiles in order.
func loadExecProfilesDir(dir string, paths ExecHostPaths, forbidden []namedRoot, st *ExecStartup, flaggedFiles map[string]string) ([]*ExecProfile, error) {
	abs, err := filepath.Abs(dir)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return nil, startupErr("--exec-profiles-dir", dir, "does not resolve: %v", err)
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return nil, startupErr("--exec-profiles-dir", dir, "is not a directory")
	}
	if err := checkPlacement("--exec-profiles-dir", dir, abs, forbidden); err != nil {
		return nil, err
	}
	ents, err := os.ReadDir(abs)
	if err != nil {
		return nil, startupErr("--exec-profiles-dir", dir, "%v", err)
	}
	var loaded []*ExecProfile
	for _, e := range ents {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ExecProfileSuffix) || strings.HasSuffix(name, ExecProposedSuffix) {
			continue
		}
		file := filepath.Join(abs, name)
		if _, dup := flaggedFiles[realpathOr(file)]; dup {
			continue
		}
		stem := strings.TrimSuffix(name, ExecProfileSuffix)
		exclude := func(why string) {
			st.Excluded = append(st.Excluded, ExecProfileExclusion{File: file, Project: stem, Reason: why})
		}
		p, err := loadExecProfileFile(file, paths, forbidden)
		switch {
		case err != nil:
			exclude(err.Error())
		case p.Project != stem:
			exclude(fmt.Sprintf("its project is %q, not its file name %q: a directory profile is <project>.json", p.Project, stem))
		case st.Profiles[p.Project] != nil:
			exclude(fmt.Sprintf("project %q is already configured by %s", p.Project, st.Sources[p.Project]))
		default:
			st.Profiles[p.Project] = p
			st.Sources[p.Project] = file
			loaded = append(loaded, p)
		}
	}
	return loaded, nil
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

// checkExecChainPlacement checks the configured spelling and every component
// before following a link, including links in parent directories. Checking
// only EvalSymlinks' final result would trust an agent-writable steering link.
func checkExecChainPlacement(flag, value, path string, forbidden []namedRoot) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := checkPlacement(flag, value, abs, forbidden); err != nil {
		return err
	}
	pending := strings.Split(abs, string(filepath.Separator))
	prefix := string(filepath.Separator)
	links := 0
	for len(pending) > 0 {
		part := pending[0]
		pending = pending[1:]
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			prefix = filepath.Dir(prefix)
			continue
		}
		candidate := filepath.Join(prefix, part)
		if err := checkPlacement(flag, value, candidate, forbidden); err != nil {
			return err
		}
		info, err := os.Lstat(candidate)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			prefix = candidate
			continue
		}
		links++
		if links > 255 {
			return startupErr(flag, value, "too many symlinks")
		}
		target, err := os.Readlink(candidate)
		if err != nil {
			return err
		}
		if filepath.IsAbs(target) {
			prefix = string(filepath.Separator)
		}
		// Do not Clean a link's target: a/../b must visit a first,
		// because a itself may be a symlink into a forbidden root.
		pending = append(strings.Split(target, string(filepath.Separator)), pending...)
	}
	return nil
}

// resolveExecNode resolves --exec-node once (default: node on PATH), checks
// its placement (the binary, its directory and its libnode) and its version.
func resolveExecNode(ctx context.Context, flagValue string, forbidden []namedRoot) (string, error) {
	found, err := execNodePath(flagValue)
	if err != nil {
		return "", err
	}
	if err := checkExecChainPlacement("--exec-node", flagValue, found, forbidden); err != nil {
		return "", err
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
		if err := checkExecChainPlacement("--exec-node", flagValue, l.Path, forbidden); err != nil {
			return "", err
		}
	}
	if err := checkExecNodeVersion(ctx, abs); err != nil {
		return "", startupErr("--exec-node", flagValue, "%v", err)
	}
	return abs, nil
}

func execNodePath(flagValue string) (string, error) {
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
	return filepath.Abs(found)
}

func checkExecNodeVersion(ctx context.Context, abs string) error {
	vctx, cancel := context.WithTimeout(ctx, execNodeVersionBudget)
	defer cancel()
	out, err := runBounded(vctx, handlerBounds{execTimeout: execNodeVersionBudget}, handlerCommand{
		path: abs, args: []string{"--version"}, dir: "/", env: execHostEnv()})
	if err != nil {
		var exited *HandlerExitError
		if errors.As(err, &exited) && len(exited.Output) > 0 {
			return fmt.Errorf("%s (--version failed)", firstExecFailure(string(exited.Output)))
		}
		return fmt.Errorf("%s --version: %w", abs, err)
	}
	version := strings.TrimSpace(string(out))
	major, minor, ok := parseNodeVersion(version)
	if !ok {
		return fmt.Errorf("%s --version printed %q, not vMAJOR.MINOR.PATCH", abs, version)
	}
	if major < ExecMinNode[0] || (major == ExecMinNode[0] && minor < ExecMinNode[1]) {
		return fmt.Errorf("version %s is below %d.%d (srt's engines)", version, ExecMinNode[0], ExecMinNode[1])
	}
	return nil
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
	commands := []ExecProfileCommand{p.Probe}
	for _, c := range p.Commands {
		commands = append(commands, c)
	}
	for _, c := range commands {
		real, err := filepath.EvalSymlinks(c.Argv0)
		if err != nil {
			return nil, startupErr("--exec-profile", file, "argv[0] %s: %v", c.Argv0, err)
		}
		if _, _, err := resolveArgv0(real, nil); err != nil {
			return nil, startupErr("--exec-profile", file, "argv[0] %s: %v", c.Argv0, err)
		}
		if err := checkExecChainPlacement("--exec-profile", file+" argv[0] "+c.Argv0, c.Argv0, forbidden); err != nil {
			return nil, err
		}
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
