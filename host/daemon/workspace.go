package daemon

// Row 134 M4b (w-software-engineering-domain §4.3, D-SE-1 A, D-SE-4 A): the
// production handler registry for the software-engineering tools.
//
// `serve --workspace-root DIR --tool-ailang-bin PATH` binds episode E to the
// worktree DIR/E. For each episode the daemon renders the §4.3 policy at
// <db-dir>/policies/E.toml (mode 0600), points AILANG_CACHE_DIR at
// <db-dir>/cache/E, and binds the six Workspace.*/Ailang.* effect names to
// one broker.AilangToolHandler running the ARCHIVED tool binary. Anything
// short of that — a flag missing, an episode id outside the grammar, a
// worktree that is absent, not a directory, or reached through a symlink —
// yields an EMPTY registry, so the coordinator refuses a declared effect
// (R8) before the plan or any effect runs. No confinement is added here:
// path, symlink, .git and deny-glob confinement stay AILANG's own (§4.3).

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// ToolBinaryRelease is the only tool-binary release the confinement matrix
// and the run outcome shapes were measured on (V32–V39, V49, V52–V63). A
// different release changes those facts, so startup refuses it (R-SE-8:
// never widen; stay on the last binary that passed AC4.1).
const ToolBinaryRelease = "AILANG v0.51.0"

// workspaceHandlerBudget bounds one episode handler's construction: the
// policy write plus the one `policy-tool` summary subprocess (~0.1 s, M4a).
const workspaceHandlerBudget = 3 * time.Second

// episodeIDPattern is the §4.3 episode grammar: an episode id names one
// directory directly under the workspace root, never a path.
var episodeIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// workspaceEffects are the six effect names an episode's handler serves.
var workspaceEffects = []string{
	broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite, broker.EffectAilangCheck,
	broker.EffectAilangRun, broker.EffectAilangDiscover, broker.EffectAilangCLI,
}

// workspaceTools is the resolved --workspace-root/--tool-ailang-bin pair.
// The zero value (either half missing) serves no effect.
type workspaceTools struct {
	root     string // canonical workspace root, "" when --workspace-root is unset
	bin      string // archived tool binary path, "" when --tool-ailang-bin is unset
	binRef   hashref.HashRef
	stateDir string // canonical <db-dir>: policies/ and cache/ live here
	// examplesDir is the canonical --examples-dir, "" when none is configured
	// (examples-search then refuses, broker.NoExamplesCorpusRefusal).
	examplesDir string
	errLog      io.Writer

	mu      sync.Mutex
	handler map[string]episodeTool // episode id -> constructed handler
}

type episodeTool struct {
	root string // the canonical worktree the handler is bound to
	h    broker.Handler
}

// WorkspaceStateError is the AC4.5 daemon-half refusal: the daemon's own
// state (the store, the archived binaries, the rendered policies, the tool
// cache) must not be reachable from inside the workspace root, where an
// agent's writes land — an agent that can edit its own policy has no policy.
type WorkspaceStateError struct {
	What string
	Err  error
}

func (e *WorkspaceStateError) Error() string { return fmt.Sprintf("%s: %v", e.What, e.Err) }
func (e *WorkspaceStateError) Unwrap() error { return e.Err }

// resolveWorkspaceRoot canonicalises --workspace-root and refuses a root
// that is not a directory or that contains the daemon's state directory or
// any path under it (D4).
func resolveWorkspaceRoot(root, dbPath string) (canonicalRoot, stateDir string, err error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	canonicalRoot, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return "", "", fmt.Errorf("resolve --workspace-root: %w", err)
	}
	if info, err := os.Stat(canonicalRoot); err != nil || !info.IsDir() {
		return "", "", fmt.Errorf("--workspace-root %q is not a directory", root)
	}
	stateDir, err = canonicalStateDir(dbPath)
	if err != nil {
		return "", "", err
	}
	db, err := filepath.Abs(dbPath)
	if err != nil {
		return "", "", err
	}
	for _, p := range []struct{ what, path string }{
		{"the store directory", stateDir},
		{"the store", db},
		{"the artifact archive", archive.New(db).Root()},
		{"the policy directory", filepath.Join(stateDir, "policies")},
		{"the tool cache directory", filepath.Join(stateDir, "cache")},
	} {
		if err := broker.CheckPolicyOutsideRoot(p.path, canonicalRoot); err != nil {
			return "", "", &WorkspaceStateError{What: p.what + " is inside --workspace-root", Err: err}
		}
	}
	return canonicalRoot, stateDir, nil
}

// resolveExamplesDir canonicalises --examples-dir and refuses one that is not
// a directory or that lies inside the canonical workspace root: a corpus an
// agent can write is a corpus it can poison for every later search.
func resolveExamplesDir(dir, canonicalRoot string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("resolve --examples-dir: %w", err)
	}
	if info, err := os.Stat(resolved); err != nil || !info.IsDir() {
		return "", fmt.Errorf("--examples-dir %q is not a directory", dir)
	}
	if err := broker.CheckPolicyOutsideRoot(resolved, canonicalRoot); err != nil {
		return "", &WorkspaceStateError{What: "the examples corpus is inside --workspace-root", Err: err}
	}
	return resolved, nil
}

// canonicalStateDir is <db-dir> with symlinks resolved (it exists once the
// store is open; before that its parent may be all that exists).
func canonicalStateDir(dbPath string) (string, error) {
	abs, err := filepath.Abs(dbPath)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(abs)
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		return resolved, nil
	}
	return dir, nil
}

// archiveToolBinary archives --tool-ailang-bin exactly as --ailang-bin is
// archived (content-addressed, read-only, version probed from the archived
// bytes) and refuses any release but ToolBinaryRelease. It returns the
// archived path: the handler never runs the configured path or a PATH lookup.
func archiveToolBinary(a *archive.Archive, bin string) (string, hashref.HashRef, error) {
	ref, err := a.Archive(bin)
	if err != nil {
		return "", hashref.HashRef{}, err
	}
	m, err := a.ReadManifest(ref)
	if err != nil {
		return "", hashref.HashRef{}, err
	}
	if got := releaseFromVersion(m.Version); got != ToolBinaryRelease {
		return "", hashref.HashRef{}, fmt.Errorf("tool binary is %q; only %q is measured for the workspace tools", got, ToolBinaryRelease)
	}
	path, err := a.Resolve(ref)
	if err != nil {
		return "", hashref.HashRef{}, err
	}
	if err := verifyArchivedTool(path, ref); err != nil {
		return "", hashref.HashRef{}, err
	}
	return path, ref, nil
}

// verifyArchivedTool re-hashes the archived bytes against their address, as
// the capsule does for the interpreter before every run.
func verifyArchivedTool(path string, ref hashref.HashRef) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if got := hex.EncodeToString(h.Sum(nil)); ref.Algo() != hashref.AlgoSHA256 || got != ref.Digest() {
		return fmt.Errorf("archived tool binary %s hashes to sha256:%s, not %s", path, got, ref)
	}
	return nil
}

// enabled reports whether both flags were given.
func (w *workspaceTools) enabled() bool { return w != nil && w.root != "" && w.bin != "" }

// registry is workspaceRegistry(episodeID) of §4.3: the six effect names
// bound to the episode's handler, or an empty registry.
func (w *workspaceTools) registry(episodeID string) broker.Registry {
	if !w.enabled() {
		return broker.Registry{}
	}
	epRoot, ok := w.episodeRoot(episodeID)
	if !ok {
		return broker.Registry{}
	}
	h, err := w.episodeHandler(episodeID, epRoot)
	if err != nil {
		fmt.Fprintf(w.errLog, "ailang-worldd: workspace tools unavailable for episode %q: %v\n", episodeID, err)
		return broker.Registry{}
	}
	reg := make(broker.Registry, len(workspaceEffects))
	for _, name := range workspaceEffects {
		reg[name] = h
	}
	return reg
}

// episodeRoot applies the episode grammar and requires root/episode to be a
// real directory: EvalSymlinks(root/episode) must be exactly root/episode, so
// a symlink can redirect an episode neither out of the root nor onto a
// sibling episode's worktree (a strictly narrower reading of §4.3's "a
// directory inside EvalSymlinks(root)").
func (w *workspaceTools) episodeRoot(episodeID string) (string, bool) {
	if !episodeIDPattern.MatchString(episodeID) {
		return "", false
	}
	want := filepath.Join(w.root, episodeID)
	resolved, err := filepath.EvalSymlinks(want)
	if err != nil || resolved != want {
		return "", false
	}
	if info, err := os.Stat(resolved); err != nil || !info.IsDir() {
		return "", false
	}
	return resolved, true
}

// episodeHandler returns the episode's cached handler, constructing it on
// first use (the constructor runs one summary subprocess, so it is built
// once per episode rather than per Dispatch).
func (w *workspaceTools) episodeHandler(episodeID, epRoot string) (broker.Handler, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if cached, ok := w.handler[episodeID]; ok && cached.root == epRoot {
		return cached.h, nil
	}
	policyDir := filepath.Join(w.stateDir, "policies")
	cacheDir := filepath.Join(w.stateDir, "cache", episodeID)
	for _, dir := range []string{policyDir, cacheDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
	}
	policyPath := filepath.Join(policyDir, episodeID+".toml")
	policy, err := broker.RenderEpisodePolicy(epRoot)
	if err != nil {
		return nil, err
	}
	if err := writePolicy(policyPath, policy); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), workspaceHandlerBudget)
	defer cancel()
	h, err := broker.NewAilangToolHandler(ctx, broker.AilangToolConfig{
		Bin: w.bin, BinRef: w.binRef, PolicyPath: policyPath, Root: epRoot, CacheDir: cacheDir,
		ExamplesDir: w.examplesDir,
	})
	if err != nil {
		return nil, err
	}
	tool := verifiedTool{path: w.bin, ref: w.binRef, h: h}
	if w.handler == nil {
		w.handler = map[string]episodeTool{}
	}
	w.handler[episodeID] = episodeTool{root: epRoot, h: tool}
	return tool, nil
}

// writePolicy replaces path with data at mode 0600, atomically.
func writePolicy(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".policy-*.tmp")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// verifiedTool re-verifies the archived tool bytes before every execution,
// as the capsule does for the interpreter: a binary changed on disk after
// startup is a handler failure (recorded `failed`), never a silent run.
type verifiedTool struct {
	path string
	ref  hashref.HashRef
	h    broker.Handler
}

func (v verifiedTool) Execute(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, error) {
	if err := verifyArchivedTool(v.path, v.ref); err != nil {
		return nil, err
	}
	return v.h.Execute(ctx, req, payload)
}
