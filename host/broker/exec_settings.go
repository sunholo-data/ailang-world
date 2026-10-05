package broker

// Row 140 M2 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.4):
// the srt settings World renders per episode × project at
// <state>/exec/<ep>.<project>.srt.json (0600). The shape is M0's measured one
// (design_docs/verification/world-row140-m0/m0_matrix.py): srt 0.0.78
// requires network.deniedDomains, and enableWeakerNetworkIsolation and
// allowAppleEvents are TOP-LEVEL keys. srt has no `summary` command, so
// World verifies the file itself: re-read, byte-compare with its canonical
// rendering, digest, and refuse a file it did not write.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// execSettingsSpec is everything one rendering depends on; every path is
// canonical (symlinks resolved) and absolute.
type execSettingsSpec struct {
	Worktree      string   // the episode worktree: writable, readable
	Cache         string   // <state>/exec-cache/<ep>/<project>: writable, readable
	OperatorHome  string   // denied for reading; anchors srt's two $HOME write defaults
	StateDir      string   // World's state dir: denied for reading
	WorkspaceRoot string   // every episode's parent: denied, so siblings are unreadable
	ReadRoots     []string // the profile's read_roots plus each one's realpath
	SeccompDir    string   // linux only: the archived seccomp helper dir [M0, V43/V44]
	GOOS          string
}

type execSettingsFS struct {
	DenyRead   []string `json:"denyRead"`
	AllowRead  []string `json:"allowRead"`
	AllowWrite []string `json:"allowWrite"`
	DenyWrite  []string `json:"denyWrite"`
}

type execSettingsNet struct {
	AllowedDomains      []string `json:"allowedDomains"`
	DeniedDomains       []string `json:"deniedDomains"`
	AllowLocalBinding   bool     `json:"allowLocalBinding"`
	AllowAllUnixSockets bool     `json:"allowAllUnixSockets"`
}

type execSettings struct {
	Filesystem                   execSettingsFS  `json:"filesystem"`
	Network                      execSettingsNet `json:"network"`
	EnableWeakerNetworkIsolation bool            `json:"enableWeakerNetworkIsolation"`
	AllowAppleEvents             bool            `json:"allowAppleEvents"`
}

// execDenyNames are World's episode deny names (§4.4: the AILANG policy's
// names, plus `.git`, which srt's own list covers only in part, V5): each is
// denied as an absolute path under the worktree, a dir or a file alike.
func execDenyNames() []string {
	names := []string{".git"}
	names = append(names, episodeDenyDirs...)
	return append(names, episodeDenyFiles...)
}

// srtDefaultWritePaths are the paths srt always allows writing (V4); World
// denies each. /private/tmp/claude is darwin's alias of /tmp/claude.
func srtDefaultWritePaths(home, goos string) []string {
	paths := []string{"/tmp/claude"}
	if goos == "darwin" {
		paths = append(paths, "/private/tmp/claude")
	}
	return append(paths, filepath.Join(home, ".npm", "_logs"), filepath.Join(home, ".claude", "debug"))
}

func hostGOOS() string { return runtime.GOOS }

// renderExecSettings is the canonical rendering (§4.4): reads fenced by
// path list (D-140-1 = A), no network (D-140-2 = A), writes only to the
// worktree and the episode's exec cache, World's deny list and srt's four
// defaults closed.
func renderExecSettings(spec execSettingsSpec) ([]byte, error) {
	for _, p := range append([]string{spec.Worktree, spec.Cache, spec.OperatorHome, spec.StateDir, spec.WorkspaceRoot}, spec.ReadRoots...) {
		if !cleanAbs(p) {
			return nil, fmt.Errorf("broker: exec settings path %q is not clean and absolute", p)
		}
	}
	var denyWrite []string
	for _, n := range execDenyNames() {
		denyWrite = append(denyWrite, filepath.Join(spec.Worktree, n))
	}
	denyWrite = append(denyWrite, srtDefaultWritePaths(spec.OperatorHome, spec.GOOS)...)
	allowRead := dedupe(append([]string{spec.Worktree, spec.Cache}, spec.ReadRoots...))
	if spec.GOOS == "linux" {
		if !cleanAbs(spec.SeccompDir) {
			return nil, fmt.Errorf("broker: exec settings: linux needs the archived seccomp helper dir, got %q", spec.SeccompDir)
		}
		allowRead = append(allowRead, spec.SeccompDir)
	}
	s := execSettings{
		Filesystem: execSettingsFS{
			DenyRead:   []string{spec.OperatorHome, spec.StateDir, spec.WorkspaceRoot},
			AllowRead:  allowRead,
			AllowWrite: []string{spec.Worktree, spec.Cache},
			DenyWrite:  denyWrite,
		},
		Network: execSettingsNet{AllowedDomains: []string{}, DeniedDomains: []string{}},
	}
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func dedupe(xs []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

// ExecSettingsError reports a settings file that is not the one World wrote:
// the call fails (recorded `failed`) and nothing runs under it.
type ExecSettingsError struct {
	Path string
	Why  string
}

func (e *ExecSettingsError) Error() string {
	return fmt.Sprintf("broker: exec settings %s is not World's rendering: %s", e.Path, e.Why)
}

// writeExecSettings writes data at path (0600, atomically), then verifies it
// and returns the settings_digest.
func writeExecSettings(path string, data []byte) (string, error) {
	if err := writeFileAtomic(path, data); err != nil {
		return "", fmt.Errorf("broker: write exec settings: %w", err)
	}
	if err := verifyExecSettings(path, data); err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// verifyExecSettings refuses anything but a regular 0600 file holding
// exactly World's canonical bytes.
func verifyExecSettings(path string, canonical []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return &ExecSettingsError{Path: path, Why: err.Error()}
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		return &ExecSettingsError{Path: path, Why: fmt.Sprintf("mode %v, want a regular 0600 file", info.Mode())}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		return &ExecSettingsError{Path: path, Why: err.Error()}
	}
	if !bytes.Equal(got, canonical) {
		return &ExecSettingsError{Path: path, Why: "its bytes differ from the canonical rendering"}
	}
	return nil
}
