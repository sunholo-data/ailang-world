package broker

// Row 140 M3 (w-workspace-exec-toolchain-effect §4.3, §4.6): the profile's
// `probe` key and the read_roots rules that need the host's paths (HOME, the
// state dir, the workspace root), which ParseExecProfile cannot see.

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// §4.6 arm 7: the profile's `probe` argv, resolved once at load like a
// command's argv[0]; the default is the first command's argv[0] --version.
func TestParseExecProfileProbe(t *testing.T) {
	bin := fakeToolDir(t, "go", "sh")
	p, err := ParseExecProfile(profileJSON(t, bin, nil))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.Probe.Argv, []string{"go", "--version"}) || p.Probe.Argv0 != filepath.Join(bin, "go") || p.Probe.Argv0SHA256 == "" {
		t.Fatalf("default probe = %+v, want go --version resolved on path", p.Probe)
	}
	p, err = ParseExecProfile(profileJSON(t, bin, func(m map[string]any) { m["probe"] = []any{"go", "version"} }))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.Probe.Argv, []string{"go", "version"}) || p.Probe.Argv0 != filepath.Join(bin, "go") {
		t.Fatalf("probe = %+v, want go version resolved on path", p.Probe)
	}
	for _, tc := range []struct {
		name  string
		probe any
		want  string
	}{
		{"empty", []any{}, "probe"},
		{"not on path", []any{"cargo", "--version"}, "does not resolve"},
		{"empty item", []any{"go", ""}, "probe"},
		{"not an array", "go version", "probe"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseExecProfile(profileJSON(t, bin, func(m map[string]any) { m["probe"] = tc.probe }))
			var perr *ExecProfileError
			if !errors.As(err, &perr) || !strings.Contains(err.Error(), tc.want) || perr.Field != "probe" {
				t.Fatalf("err = %v, want an *ExecProfileError on field probe naming %q", err, tc.want)
			}
		})
	}
}

// hostPathsFixture lays HOME, the state dir and the workspace root in three
// different parents, so an ancestor of one is not an ancestor of another.
func hostPathsFixture(t *testing.T) (base string, h ExecHostPaths) {
	t.Helper()
	base = canonicalTempDir(t)
	h = ExecHostPaths{OperatorHome: filepath.Join(base, "h", "home"), StateDir: filepath.Join(base, "s", "state"),
		WorkspaceRoot: filepath.Join(base, "w", "ws")}
	mustMkdir(t, h.OperatorHome, h.StateDir, h.WorkspaceRoot, filepath.Join(h.WorkspaceRoot, "ep1"),
		filepath.Join(h.StateDir, "x"), filepath.Join(h.OperatorHome, "go", "pkg", "mod"), filepath.Join(base, "opt", "mod"))
	return base, h
}

// §4.3 [REVISED r1]: a read_roots entry is refused when its literal path OR
// its realpath equals or is an ancestor of HOME, the state dir or the
// workspace root, lies inside the state dir or the workspace root, equals or
// contains a probe decoy, or is / (MUT-READROOT-ANCESTOR).
func TestCheckExecReadRoots(t *testing.T) {
	base, h := hostPathsFixture(t)
	link := func(name, target string) string {
		p := filepath.Join(base, "links", name)
		mustMkdir(t, filepath.Dir(p))
		if err := os.Symlink(target, p); err != nil {
			t.Fatal(err)
		}
		return p
	}
	decoy := ExecProbeDecoys(h)[0]
	if decoy != filepath.Join(h.OperatorHome, ".ailang-worldd-exec-probe-decoy") {
		t.Fatalf("decoy 0 = %s", decoy)
	}
	for _, tc := range []struct {
		name, root, want string
	}{
		{"HOME itself", h.OperatorHome, "operator HOME"},
		{"an ancestor of HOME", filepath.Join(base, "h"), "operator HOME"},
		{"a symlink to HOME", link("home", h.OperatorHome), "operator HOME"},
		{"a symlink to an ancestor of HOME", link("h", filepath.Join(base, "h")), "operator HOME"},
		{"the state dir", h.StateDir, "state dir"},
		{"an ancestor of the state dir", filepath.Join(base, "s"), "state dir"},
		{"inside the state dir", filepath.Join(h.StateDir, "x"), "state dir"},
		{"the workspace root", h.WorkspaceRoot, "workspace root"},
		{"an ancestor of the workspace root", filepath.Join(base, "w"), "workspace root"},
		{"inside the workspace root", filepath.Join(h.WorkspaceRoot, "ep1"), "workspace root"},
		{"a symlink into the workspace root", link("ep1", filepath.Join(h.WorkspaceRoot, "ep1")), "workspace root"},
		{"the HOME probe decoy", decoy, "probe decoy"},
		{"/", "/", `"/"`},
		{"absent", filepath.Join(base, "nope"), "does not resolve"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckExecReadRoots([]string{filepath.Join(base, "opt", "mod"), tc.root}, h)
			var perr *ExecProfileError
			if !errors.As(err, &perr) || perr.Field != "read_roots" || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("read_roots [%s] -> %v, want a read_roots refusal naming %q", tc.root, err, tc.want)
			}
		})
	}
	// Inside HOME (the Go module cache) and outside every denied root: kept.
	if err := CheckExecReadRoots([]string{filepath.Join(h.OperatorHome, "go", "pkg", "mod"), filepath.Join(base, "opt", "mod")}, h); err != nil {
		t.Fatalf("legitimate read_roots refused: %v", err)
	}
	// APFS is case-insensitive: /…/H is the ancestor h of HOME there.
	if runtime.GOOS == "darwin" {
		upper := filepath.Join(base, "H")
		if err := CheckExecReadRoots([]string{upper}, h); err == nil || !strings.Contains(err.Error(), "operator HOME") {
			t.Fatalf("darwin read_roots [%s] -> %v, want the case-folded ancestor refused", upper, err)
		}
	}
}
