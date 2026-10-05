package broker

// Row 140 M2 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.6
// "Integrity of the confinement stack"): World runs srt from its own archive
// of the WHOLE install tree — the package alone does not load
// (ERR_MODULE_NOT_FOUND: commander, V42) — under
// <state>/archive/exec-sandbox/<digest>/node_modules/, and re-verifies that
// tree, node and libnode before every call (≈ 56 ms, V39/V42). A mismatch is
// a handler failure naming the file, with nothing spawned.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// The srt pin (M0, design_docs/verification/world-row140-m0/pin.json, bound
// by TestExecSandboxPinIsM0s). A different release changes the measured
// confinement facts, so it is refused (R-140-2).
const (
	ExecSandboxPackage   = "@anthropic-ai/sandbox-runtime"
	ExecSandboxRelease   = "0.0.78"
	ExecSandboxCLISHA256 = "3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e"
)

// ExecSandboxPin is the version and cli.js digest an install must carry.
type ExecSandboxPin struct {
	Version   string
	CLISHA256 string
}

// ExecSandboxConfig names an installed srt and the node that runs it.
type ExecSandboxConfig struct {
	// NodeModules is the install's node_modules dir (holding
	// @anthropic-ai/sandbox-runtime and its runtime dependencies).
	NodeModules string
	// Node is the absolute path of the node binary that runs srt.
	Node string
	// StateDir is World's canonical state dir; the archive lives under it.
	StateDir string
	// Pin is the required release; the zero value is the production pin.
	Pin ExecSandboxPin
}

// execFileDigest is one hashed file outside the archive (node, libnode).
type execFileDigest struct {
	Path   string
	SHA256 string
}

// execManifestEntry is one archived path: kind f (file), x (executable
// file), l (symlink) or d (directory); Sum is a file's sha256 or a link's
// target.
type execManifestEntry struct {
	Rel  string
	Kind string
	Sum  string
}

// ExecSandbox is an archived, verified srt plus the node that runs it.
type ExecSandbox struct {
	Dir        string // <state>/archive/exec-sandbox/<Digest>/node_modules
	Digest     string // sha256 of the manifest
	Version    string
	CLISHA256  string
	Node       string
	NodeSHA256 string
	Libnode    []execFileDigest
	manifest   []execManifestEntry
}

// CLI is the archived srt entry point World runs.
func (s *ExecSandbox) CLI() string {
	return filepath.Join(s.Dir, "@anthropic-ai", "sandbox-runtime", "dist", "cli.js")
}

// SeccompDir is the archived seccomp helper dir srt runs INSIDE the sandbox
// on linux, which the rendering re-allows for reading [M0, V43/V44]; "" on
// other platforms.
func (s *ExecSandbox) SeccompDir(goos, goarch string) string {
	if goos != "linux" {
		return ""
	}
	arch := map[string]string{"amd64": "x64", "arm64": "arm64"}[goarch]
	if arch == "" {
		arch = goarch
	}
	return filepath.Join(s.Dir, "@anthropic-ai", "sandbox-runtime", "vendor", "seccomp", arch)
}

// ExecStackDriftError is a per-call re-verification failure: the named file
// of the confinement stack is not what was verified at startup.
type ExecStackDriftError struct {
	Path string
	Why  string
}

func (e *ExecStackDriftError) Error() string {
	return fmt.Sprintf("broker: exec confinement stack drifted: %s: %s", e.Path, e.Why)
}

// ArchiveExecSandbox checks the install against the pin, copies the whole
// tree into the state dir's archive (read-only files) and records the
// sha256 of node and any libnode beside it.
func ArchiveExecSandbox(cfg ExecSandboxConfig) (*ExecSandbox, error) {
	pin := cfg.Pin
	if pin == (ExecSandboxPin{}) {
		pin = ExecSandboxPin{Version: ExecSandboxRelease, CLISHA256: ExecSandboxCLISHA256}
	}
	if !filepath.IsAbs(cfg.Node) || !cleanAbs(cfg.StateDir) {
		return nil, fmt.Errorf("broker: exec sandbox: node %q and state dir %q must be absolute", cfg.Node, cfg.StateDir)
	}
	src, err := filepath.EvalSymlinks(cfg.NodeModules)
	if err != nil {
		return nil, fmt.Errorf("broker: exec sandbox: %w", err)
	}
	pkgDir := filepath.Join(src, "@anthropic-ai", "sandbox-runtime")
	var pkg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	raw, err := os.ReadFile(filepath.Join(pkgDir, "package.json"))
	if err != nil || json.Unmarshal(raw, &pkg) != nil {
		return nil, fmt.Errorf("broker: exec sandbox: %s has no readable %s package.json", src, ExecSandboxPackage)
	}
	if pkg.Version != pin.Version {
		return nil, fmt.Errorf("broker: exec sandbox: %s is %s; only %s is measured (M0)", ExecSandboxPackage, pkg.Version, pin.Version)
	}
	cliSum, err := sha256File(filepath.Join(pkgDir, "dist", "cli.js"))
	if err != nil {
		return nil, fmt.Errorf("broker: exec sandbox: %w", err)
	}
	if cliSum != pin.CLISHA256 {
		return nil, fmt.Errorf("broker: exec sandbox: dist/cli.js sha256 %s is not the pinned %s", cliSum, pin.CLISHA256)
	}
	manifest, err := execManifest(src)
	if err != nil {
		return nil, err
	}
	digest := manifestDigest(manifest)
	base := filepath.Join(cfg.StateDir, "archive", "exec-sandbox")
	final := filepath.Join(base, digest)
	if _, err := os.Lstat(final); err != nil {
		if err := os.MkdirAll(base, 0o700); err != nil {
			return nil, err
		}
		tmp, err := os.MkdirTemp(base, ".tmp-")
		if err != nil {
			return nil, err
		}
		defer func() { _ = os.RemoveAll(tmp) }()
		if err := copyTree(src, filepath.Join(tmp, "node_modules")); err != nil {
			return nil, fmt.Errorf("broker: exec sandbox: archive: %w", err)
		}
		if err := os.Rename(tmp, final); err != nil {
			return nil, err
		}
	}
	sb := &ExecSandbox{Dir: filepath.Join(final, "node_modules"), Digest: digest, Version: pkg.Version, CLISHA256: cliSum,
		Node: cfg.Node, manifest: manifest}
	if sb.NodeSHA256, err = sha256File(cfg.Node); err != nil {
		return nil, fmt.Errorf("broker: exec sandbox: node: %w", err)
	}
	if sb.Libnode, err = libnodeDigests(cfg.Node); err != nil {
		return nil, err
	}
	if err := sb.Verify(); err != nil {
		return nil, err
	}
	return sb, nil
}

// libnodeDigests hashes node's shared library where node's install has one
// (<prefix>/lib/libnode.*, Homebrew's layout); an official linux build is
// static and has none. node's other dylibs are trusted by placement
// (R-140-11).
func libnodeDigests(node string) ([]execFileDigest, error) {
	real, err := filepath.EvalSymlinks(node)
	if err != nil {
		return nil, fmt.Errorf("broker: exec sandbox: node: %w", err)
	}
	matches, _ := filepath.Glob(filepath.Join(filepath.Dir(filepath.Dir(real)), "lib", "libnode.*"))
	sort.Strings(matches)
	var out []execFileDigest
	for _, m := range matches {
		sum, err := sha256File(m)
		if err != nil {
			return nil, fmt.Errorf("broker: exec sandbox: libnode: %w", err)
		}
		out = append(out, execFileDigest{Path: m, SHA256: sum})
	}
	return out, nil
}

// Verify re-checks the archived tree entry by entry, node and libnode, and
// names the first file that drifted.
func (s *ExecSandbox) Verify() error {
	got, err := execManifest(s.Dir)
	if err != nil {
		return &ExecStackDriftError{Path: s.Dir, Why: err.Error()}
	}
	want := map[string]execManifestEntry{}
	for _, e := range s.manifest {
		want[e.Rel] = e
	}
	seen := map[string]bool{}
	for _, e := range got {
		seen[e.Rel] = true
		w, ok := want[e.Rel]
		switch {
		case !ok:
			return &ExecStackDriftError{Path: filepath.Join(s.Dir, e.Rel), Why: "not in the archived manifest"}
		case w != e:
			return &ExecStackDriftError{Path: filepath.Join(s.Dir, e.Rel), Why: "changed since it was archived"}
		}
	}
	for _, e := range s.manifest {
		if !seen[e.Rel] {
			return &ExecStackDriftError{Path: filepath.Join(s.Dir, e.Rel), Why: "missing from the archive"}
		}
	}
	for _, f := range append([]execFileDigest{{Path: s.Node, SHA256: s.NodeSHA256}}, s.Libnode...) {
		if err := verifyFileDigest(f.Path, f.SHA256); err != nil {
			return err
		}
	}
	return nil
}

func verifyFileDigest(path, want string) error {
	sum, err := sha256File(path)
	if err != nil {
		return &ExecStackDriftError{Path: path, Why: err.Error()}
	}
	if sum != want {
		return &ExecStackDriftError{Path: path, Why: fmt.Sprintf("sha256 %s, verified as %s", sum, want)}
	}
	return nil
}

// execManifest walks root in lexical order.
func execManifest(root string) ([]execManifestEntry, error) {
	var out []execManifestEntry
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == root {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		switch t := d.Type(); {
		case t&fs.ModeSymlink != 0:
			target, err := os.Readlink(p)
			if err != nil {
				return err
			}
			out = append(out, execManifestEntry{Rel: rel, Kind: "l", Sum: target})
		case t.IsDir():
			out = append(out, execManifestEntry{Rel: rel, Kind: "d"})
		case t.IsRegular():
			info, err := d.Info()
			if err != nil {
				return err
			}
			sum, err := sha256File(p)
			if err != nil {
				return err
			}
			kind := "f"
			if info.Mode().Perm()&0o111 != 0 {
				kind = "x"
			}
			out = append(out, execManifestEntry{Rel: rel, Kind: kind, Sum: sum})
		default:
			return fmt.Errorf("%s is not a file, directory or symlink", p)
		}
		return nil
	})
	return out, err
}

func manifestDigest(m []execManifestEntry) string {
	h := sha256.New()
	for _, e := range m {
		fmt.Fprintf(h, "%s %s %s\n", e.Kind, e.Sum, e.Rel)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// copyTree copies src to dst: directories 0755, files read-only (0444, or
// 0555 when any exec bit was set), symlinks as links.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		switch t := d.Type(); {
		case t&fs.ModeSymlink != 0:
			link, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		case t.IsDir():
			return os.MkdirAll(target, 0o755)
		case t.IsRegular():
			info, err := d.Info()
			if err != nil {
				return err
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			mode := os.FileMode(0o444)
			if info.Mode().Perm()&0o111 != 0 {
				mode = 0o555
			}
			return os.WriteFile(target, data, mode)
		default:
			return fmt.Errorf("%s is not a file, directory or symlink", p)
		}
	})
}

// hostGOARCH is the platform the seccomp helper dir is chosen for.
func hostGOARCH() string { return runtime.GOARCH }
