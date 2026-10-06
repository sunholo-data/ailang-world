package daemon

// Row 141 (w-workspace-project-layouts): project layouts the one-directory-
// per-episode default does not fit.
//
// M1 — a MODULE ROOT: `--workspace-module-root REL` (and the per-episode
// `--workspace-episode-module-root EP=REL`) makes <root>/<episode>/REL the
// AILANG sandbox, because policy-tool runs every CLI child with cwd = the
// sandbox and AILANG resolves bare imports against the cwd. The grammar is
// row 140's exec-profile path grammar; resolution is projectDir's rule
// (symlink-free, an existing directory) and NOTHING is ever created.
//
// M2 — a read-only REGISTRY PACKAGE CACHE: v0.52.1 resolves `pkg/...` imports
// under $HOME/.ailang/cache/registry, and the episode HOME is the empty
// per-episode cache dir. `--workspace-package-cache DIR` links an operator
// snapshot there, after checking at startup that the snapshot is a tree
// nothing can write (and that is not reachable through a symlink), so the
// unconfined `pkg_docs` child can no longer store a fetched package in it.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// checkModuleRoot refuses a REL outside the shared exec path grammar, an
// unclean one, and one renderPolicy would refuse. Every error names the flag
// and the raw value.
func checkModuleRoot(flag, value, rel, root string) error {
	switch {
	case !broker.ExecPathOk(rel):
		return fmt.Errorf("%s %q: %q must be relative, with no leading / or -, and no .. segment", flag, value, rel)
	case filepath.Clean(rel) != rel:
		return fmt.Errorf("%s %q: %q is not clean", flag, value, rel)
	}
	if _, err := broker.RenderEpisodePolicy(filepath.Join(root, rel)); err != nil {
		return fmt.Errorf("%s %q: %v", flag, value, err)
	}
	return nil
}

// resolveModuleRoots validates the default and the per-episode overrides
// without touching the filesystem.
func resolveModuleRoots(def string, pairs []string, root string) (string, map[string]string, error) {
	if def != "" {
		if err := checkModuleRoot("--workspace-module-root", def, def, root); err != nil {
			return "", nil, err
		}
	}
	byEp := map[string]string{}
	for _, pair := range pairs {
		const flag = "--workspace-episode-module-root"
		ep, rel, ok := strings.Cut(pair, "=")
		switch {
		case !ok:
			return "", nil, fmt.Errorf("%s %q: is not EP=REL", flag, pair)
		case !episodeIDPattern.MatchString(ep):
			return "", nil, fmt.Errorf("%s %q: episode %q is outside the episode grammar %s", flag, pair, ep, episodeIDPattern)
		}
		if _, dup := byEp[ep]; dup {
			return "", nil, fmt.Errorf("%s %q: maps episode %q a second time", flag, pair, ep)
		}
		if err := checkModuleRoot(flag, pair, rel, root); err != nil {
			return "", nil, err
		}
		byEp[ep] = rel
	}
	return def, byEp, nil
}

// sandboxRoot is the episode's AILANG sandbox: <epRoot>/REL, REL being the
// episode's override, else the default, else ".". It must be a real directory
// reached through no symlink (projectDir's rule). Never MkdirAll.
func (w *workspaceTools) sandboxRoot(episodeID, epRoot string) (string, error) {
	rel := w.moduleRoot
	if o, ok := w.episodeModuleRoot[episodeID]; ok {
		rel = o
	}
	if rel == "" {
		rel = "."
	}
	want := filepath.Join(epRoot, rel)
	got, err := filepath.EvalSymlinks(want)
	if err != nil {
		return "", fmt.Errorf("module root %q: %v", rel, err)
	}
	if got != want {
		return "", fmt.Errorf("module root %q: resolves through a symlink to %q", rel, got)
	}
	if info, err := os.Stat(want); err != nil || !info.IsDir() {
		return "", fmt.Errorf("module root %q: not a directory", rel)
	}
	return want, nil
}

// episodeRefusal is an episode-construction refusal that carries its own
// complete operator line, which registry() prints verbatim (once) instead of
// the generic "workspace tools unavailable" line.
type episodeRefusal struct{ line string }

func (e *episodeRefusal) Error() string { return e.line }

// geteuid is a seam: the uid-0 refusal is unreachable from a non-root test.
var geteuid = os.Geteuid

// withinDir reports whether path is base or lies under it (lexically).
func withinDir(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// resolvePackageCache checks --workspace-package-cache and returns its
// canonical path, the snapshot digest and the package count. It refuses, in
// order: not a directory; inside the workspace root; inside the state
// directory; uid 0 (mode bits do not stop root); any symlink, special file or
// writable entry (DIR itself included); fewer than one
// <ns>/<name>/<ver>/ailang.toml package.
func resolvePackageCache(dir, canonicalRoot, stateDir string) (canonical, digest string, n int, err error) {
	notDir := fmt.Errorf("--workspace-package-cache %q is not a directory", dir)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", "", 0, notDir
	}
	if canonical, err = filepath.EvalSymlinks(abs); err != nil {
		return "", "", 0, notDir
	}
	if info, err := os.Stat(canonical); err != nil || !info.IsDir() {
		return "", "", 0, notDir
	}
	if err := broker.CheckPolicyOutsideRoot(canonical, canonicalRoot); err != nil {
		return "", "", 0, &WorkspaceStateError{What: "the package cache is inside --workspace-root", Err: err}
	}
	if withinDir(canonical, stateDir) {
		return "", "", 0, &WorkspaceStateError{What: "the package cache is inside the state directory",
			Err: fmt.Errorf("%s is under %s", canonical, stateDir)}
	}
	if geteuid() == 0 {
		return "", "", 0, errors.New("the daemon runs as root (uid 0), where mode bits do not stop writes; run it as an unprivileged user")
	}
	var records []string
	err = filepath.WalkDir(canonical, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info() // Lstat semantics: a symlink is not followed
		if err != nil {
			return err
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			return fmt.Errorf("%s is a symlink; the package cache must hold none", path)
		case !info.IsDir() && !info.Mode().IsRegular():
			return fmt.Errorf("%s is neither a directory nor a regular file", path)
		case info.Mode().Perm()&0o222 != 0:
			return fmt.Errorf("%s is writable (mode %v); the package cache must be read-only (chmod -R a-w)", path, info.Mode().Perm())
		}
		rel, err := filepath.Rel(canonical, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if info.IsDir() {
			records = append(records, rel+"\x00d")
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(content)
		records = append(records, rel+"\x00f\x00"+hex.EncodeToString(sum[:]))
		if segs := strings.Split(rel, "/"); len(segs) == 4 && segs[3] == "ailang.toml" {
			n++
		}
		return nil
	})
	if err != nil {
		return "", "", 0, err
	}
	if n == 0 {
		return "", "", 0, fmt.Errorf("%s holds no package (<ns>/<name>/<ver>/ailang.toml)", canonical)
	}
	return canonical, packageCacheDigest(records), n, nil
}

// packageCacheDigest is "sha256:" + hex(sha256) over the records sorted in
// byte order of the whole record, each terminated by a newline. A record is
// relpath NUL "d" for a directory and relpath NUL "f" NUL hex(sha256(content))
// for a regular file; DIR itself has none. Sorted records, not walk order:
// "a/b" sorts after "a-c" although the walk visits it first.
func packageCacheDigest(records []string) string {
	sort.Strings(records)
	h := sha256.New()
	for _, r := range records {
		_, _ = io.WriteString(h, r+"\n")
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

// linkPackageCache makes <cacheDir>/.ailang/cache/registry the symlink to the
// snapshot, never deleting anything that holds content. With no snapshot
// configured it only removes a stale symlink (an earlier run with the flag
// set made it); a real directory there is today's behaviour and left alone.
func (w *workspaceTools) linkPackageCache(episodeID, cacheDir string) error {
	link := filepath.Join(cacheDir, ".ailang", "cache", "registry")
	if w.packageCache == "" {
		fi, err := os.Lstat(link)
		if err == nil && fi.Mode()&os.ModeSymlink != 0 {
			return os.Remove(link)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(link), 0o700); err != nil {
		return err
	}
	fi, err := os.Lstat(link)
	switch {
	case os.IsNotExist(err):
	case err != nil:
		return err
	case fi.Mode()&os.ModeSymlink != 0 && readlinkIs(link, w.packageCache):
		return nil
	case fi.IsDir() && dirEmpty(link):
		if err := os.Remove(link); err != nil {
			return err
		}
	case fi.Mode()&os.ModeSymlink != 0:
		old, _ := os.Readlink(link)
		return &episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache refused for episode %q: %s is a symlink to %s, "+
			"not to --workspace-package-cache %s; remove it if the snapshot moved", episodeID, link, old, w.packageCache)}
	default:
		return &episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache refused for episode %q: %s is not empty "+
			"and was not provisioned by --workspace-package-cache; clear it manually", episodeID, link)}
	}
	return os.Symlink(w.packageCache, link)
}

func readlinkIs(link, want string) bool {
	got, err := os.Readlink(link)
	return err == nil && got == want
}

func dirEmpty(dir string) bool {
	entries, err := os.ReadDir(dir)
	return err == nil && len(entries) == 0
}

// lockMaxBytes bounds the lock read: ailang.lock is agent-writable.
const lockMaxBytes = 1 << 20

// validPkgRel is true for exactly <ns>/<name>/<ver>: three non-empty
// segments, none "." or "..", so a lock entry cannot walk out of the cache.
func validPkgRel(rel string) bool {
	segs := strings.Split(rel, "/")
	if len(segs) != 3 {
		return false
	}
	for _, s := range segs {
		if s == "" || s == "." || s == ".." {
			return false
		}
	}
	return true
}

// checkLockCoverage refuses an episode whose <sandbox>/ailang.lock requires a
// registry package absent from the snapshot: such a snapshot passes every
// startup check yet fails the agent with `cache not found`. nil when no
// snapshot is configured or there is no lock. The lock is agent-writable and
// read with operator privileges, so it fails closed: it must be a regular
// file of at most 1 MiB holding valid JSON, and a registry entry whose
// name/version is not a plain <ns>/<name>/<ver> counts as not covered.
func (w *workspaceTools) checkLockCoverage(episodeID, sandbox string) error {
	if w.packageCache == "" {
		return nil
	}
	lock := filepath.Join(sandbox, "ailang.lock")
	bad := func(why string) error {
		return &episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache does not cover episode %q: lock %s %s; "+
			"rebuild the snapshot per QUICKSTART §11", episodeID, lock, why)}
	}
	// One open, then Fstat on it: no Lstat/Open window in which the agent can
	// swap a FIFO in (a blocking open would hold w.mu). O_NONBLOCK makes a
	// FIFO open return at once and O_NOFOLLOW refuses a symlink (ELOOP).
	f, err := os.OpenFile(lock, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	switch {
	case os.IsNotExist(err):
		return nil
	case errors.Is(err, syscall.ELOOP):
		return bad("is not a regular file")
	case err != nil:
		return bad("cannot be read: " + err.Error())
	}
	defer func() { _ = f.Close() }()
	if fi, err := f.Stat(); err != nil {
		return bad("cannot be inspected: " + err.Error())
	} else if !fi.Mode().IsRegular() {
		return bad("is not a regular file")
	}
	data, err := io.ReadAll(io.LimitReader(f, lockMaxBytes+1))
	switch {
	case err != nil:
		return bad("cannot be read: " + err.Error())
	case len(data) > lockMaxBytes:
		return bad("is larger than 1 MiB")
	}
	var doc struct {
		Packages []struct{ Name, Version, Source string } `json:"packages"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return bad("is not valid JSON: " + err.Error())
	}
	for _, p := range doc.Packages {
		if p.Source != "registry" {
			continue
		}
		rel := p.Name + "/" + p.Version
		covered := false
		if validPkgRel(rel) {
			st, err := os.Stat(filepath.Join(w.packageCache, rel, "ailang.toml"))
			covered = err == nil && st.Mode().IsRegular()
		}
		if !covered {
			return &episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache does not cover episode %q: "+
				"lock requires %s@%s, absent from %s; rebuild the snapshot per QUICKSTART §11", episodeID, p.Name, p.Version, w.packageCache)}
		}
	}
	return nil
}
