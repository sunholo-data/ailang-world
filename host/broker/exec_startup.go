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
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

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
