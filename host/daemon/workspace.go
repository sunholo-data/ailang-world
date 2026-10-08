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
// leaves the AILANG names unbound, so the coordinator refuses a declared
// AILANG effect (R8) before the plan or any effect runs. Row 152: the two
// families are bound independently. Workspace.Exec needs only the episode
// worktree, so it is bound whenever that resolves, even when the AILANG
// handler cannot be built (module root, package cache, policy, summary).
// No confinement is added here:
// path, symlink, .git and deny-glob confinement stay AILANG's own (§4.3).

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
// and the run outcome shapes were measured on (V32–V39, V49, V52–V63 on
// v0.51.0; re-proved on v0.52.1 by the row-134 design §13, V69–V80, in row 135
// M0). A different release changes those facts, so startup refuses it (R-SE-8:
// never widen; stay on the last binary that passed AC4.1).
const ToolBinaryRelease = "AILANG v0.52.1"

// workspaceHandlerBudget bounds one episode handler's construction: the
// policy write plus the one `policy-tool` summary subprocess (~0.1 s, M4a).
// It is a var only as a test seam (production value 3 s; tests set it with
// withHandlerBudget, never concurrently with a build).
var workspaceHandlerBudget = 3 * time.Second

// episodeIDPattern is the §4.3 episode grammar: an episode id names one
// directory directly under the workspace root, never a path.
var episodeIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// ailangToolEffects are the eight effect names an episode's AILANG tool
// handler serves: row 134's six plus row 135's Ailang.RunEnv and
// Ailang.RunNet. All eight are bound together whenever the episode's AILANG
// handler is built (R8 refuses a plan whose declared effects lack a handler,
// and ailang-run declares all three run effects); the operator's run
// allowlist is enforced inside the handler, never by leaving a name
// unregistered (§4.2). When the handler cannot be built (row 152) none of the
// eight is bound and R8 refuses only the AILANG transitions.
var ailangToolEffects = []string{
	broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite, broker.EffectAilangCheck,
	broker.EffectAilangRun, broker.EffectAilangRunEnv, broker.EffectAilangRunNet,
	broker.EffectAilangDiscover, broker.EffectAilangCLI,
}

// workspaceEffects are the nine names an episode's registry binds: the eight
// above plus row 140's Workspace.Exec, bound whenever the episode worktree
// resolves, independently of the AILANG handler (row 152; workspace-exec
// declares only it). With no exec profile configured its handler is
// broker.ExecUnconfiguredHandler: a typed refusal that spawns nothing, never
// an unregistered name; with one, the episode's broker.ExecHandler.
var workspaceEffects = append(append([]string(nil), ailangToolEffects...), broker.EffectWorkspaceExec)

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
	// runCaps is the operator's ailang-run allowlist (row 135 §4.6), handed
	// to every episode's handler; the zero value admits IO/FS runs only.
	runCaps broker.RunCapsConfig
	// exec is the Workspace.Exec configuration (row 140): nil, the default,
	// binds the typed unconfigured refusal (broker.ExecUnconfiguredHandler).
	// The `serve --exec-*` flags fill it after their startup checks and the
	// startup probe (configureExec).
	exec *workspaceExec
	// moduleRoot and episodeModuleRoot are the --workspace-module-root and
	// --workspace-episode-module-root values (row 141 M1), validated at
	// startup; empty means the episode worktree itself is the sandbox.
	moduleRoot        string
	episodeModuleRoot map[string]string
	// packageCache and packageCacheDigest are the canonical
	// --workspace-package-cache and its startup digest (row 141 M2); "" when
	// unset.
	packageCache       string
	packageCacheDigest string
	errLog             io.Writer

	mu      sync.Mutex
	handler map[string]episodeTool // episode id -> constructed handler (guarded by mu)

	// execMu guards execH only. It is never held with mu, so exec
	// construction never waits behind an AILANG build (row 152).
	execMu sync.Mutex
	execH  map[string]episodeTool // episode id -> constructed exec handler
}

// workspaceExec is the verified exec profiles and the archived sandbox that
// runs them (row 140 §4.6).
type workspaceExec struct {
	// profile is the episode default: the only profile, when exactly one is
	// configured (nil with several).
	profile *broker.ExecProfile
	// byProject and episodeProject are the --exec-episode-project map.
	byProject      map[string]*broker.ExecProfile
	episodeProject map[string]string
	sandbox        *broker.ExecSandbox
	maxOutputBytes int64  // 0 = broker.DefaultExecMaxOutputBytes
	operatorHome   string // "" = os.UserHomeDir, as the probe resolved it
}

// profileFor selects an episode's profile: its --exec-episode-project
// mapping, else the only profile; nil when several are configured and the
// episode is unmapped.
func (e *workspaceExec) profileFor(episodeID string) *broker.ExecProfile {
	if project, ok := e.episodeProject[episodeID]; ok {
		return e.byProject[project]
	}
	return e.profile
}

// execStartupHook adjusts the exec startup config before it runs. It is nil
// in production; tests use it for the srt shims' pin and the probe client.
var execStartupHook func(*broker.ExecStartupConfig)

// configureExec applies the --exec-* flags (row 140 §4.3, §4.6): they need
// the workspace tools, and broker.StartExec verifies the profiles, the
// placement of the confinement stack and the pin, and runs the startup probe
// for every profile. Any refusal stops startup; nothing is enabled partly.
func (w *workspaceTools) configureExec(ctx context.Context, cfg Config) error {
	sc := broker.ExecStartupConfig{ProfileFiles: cfg.ExecProfiles, EpisodeProjects: cfg.ExecEpisodeProjects,
		SandboxDir: cfg.ExecSandbox, Node: cfg.ExecNode, MaxOutputBytes: cfg.ExecMaxOutputBytes}
	if len(sc.ProfileFiles)+len(sc.EpisodeProjects) == 0 && sc.SandboxDir == "" && sc.Node == "" && sc.MaxOutputBytes == 0 {
		return nil
	}
	if !w.enabled() {
		return fmt.Errorf("the --exec-* flags need the workspace tools (--workspace-root and --tool-ailang-bin): " +
			"workspace-exec runs in an episode worktree")
	}
	sc.Paths = broker.ExecHostPaths{StateDir: w.stateDir, WorkspaceRoot: w.root}
	if execStartupHook != nil {
		execStartupHook(&sc)
	}
	st, err := broker.StartExec(ctx, sc)
	if err != nil {
		return err
	}
	w.exec = &workspaceExec{profile: st.Only, byProject: st.Profiles, episodeProject: st.EpisodeProjects,
		sandbox: st.Sandbox, maxOutputBytes: st.MaxOutputBytes, operatorHome: st.OperatorHome}
	return nil
}

// execFailedHandler answers every call with the error that kept the
// episode's exec handler from being built: a recorded failure, never a run
// and never an unbound name (R8).
type execFailedHandler struct{ err error }

func (h execFailedHandler) Execute(context.Context, broker.EffectRequest, []byte) ([]byte, error) {
	return nil, h.err
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

// registry is workspaceRegistry(episodeID) of §4.3: Workspace.Exec bound to
// the episode's own handler whenever the worktree resolves, plus the eight
// AILANG tool effect names bound to the episode's AILANG handler when it can
// be built. A flag missing or an unsafe worktree yields an empty registry; an
// AILANG build failure prints its operator line and leaves only Workspace.Exec.
func (w *workspaceTools) registry(episodeID string) broker.Registry {
	if !w.enabled() {
		return broker.Registry{}
	}
	epRoot, ok := w.episodeRoot(episodeID)
	if !ok {
		return broker.Registry{}
	}
	reg := make(broker.Registry, len(workspaceEffects))
	// The AILANG family first, exec second: today's operator-line order.
	if h, err := w.ailangFamily(episodeID, epRoot); err != nil {
		var refusal *episodeRefusal
		if errors.As(err, &refusal) {
			fmt.Fprintln(w.errLog, refusal.line)
		} else {
			fmt.Fprintf(w.errLog, "ailang-worldd: workspace tools unavailable for episode %q: %v\n", episodeID, err)
		}
	} else {
		for _, name := range ailangToolEffects {
			reg[name] = h
		}
	}
	reg[broker.EffectWorkspaceExec] = w.execHandler(episodeID, epRoot)
	return reg
}

// ailangFamily resolves the episode's module root and builds (or fetches)
// its AILANG tool handler: the inputs only the eight AILANG names need.
func (w *workspaceTools) ailangFamily(episodeID, epRoot string) (broker.Handler, error) {
	sandbox, err := w.sandboxRoot(episodeID, epRoot)
	if err != nil {
		return nil, err
	}
	return w.episodeHandler(episodeID, sandbox)
}

// execHandler is the episode's Workspace.Exec handler: the typed refusal
// when no profile is configured, else a broker.ExecHandler bound to the
// episode worktree, built once (it renders and verifies the episode's srt
// settings) and cached like the AILANG tool handler.
func (w *workspaceTools) execHandler(episodeID, epRoot string) broker.Handler {
	if w.exec == nil {
		return broker.ExecUnconfiguredHandler{}
	}
	w.execMu.Lock()
	defer w.execMu.Unlock()
	if cached, ok := w.execH[episodeID]; ok && cached.root == epRoot {
		return cached.h
	}
	profile := w.exec.profileFor(episodeID)
	if profile == nil {
		return broker.ExecRefusalHandler{Why: broker.ExecNoEpisodeProfileRefusal(episodeID)}
	}
	var h broker.Handler
	eh, err := broker.NewExecHandler(broker.ExecHandlerConfig{Profile: profile, Sandbox: w.exec.sandbox,
		Episode: episodeID, Worktree: epRoot, WorkspaceRoot: w.root, StateDir: w.stateDir,
		OperatorHome: w.exec.operatorHome, MaxOutputBytes: w.exec.maxOutputBytes})
	if err != nil {
		fmt.Fprintf(w.errLog, "ailang-worldd: workspace-exec unavailable for episode %q: %v\n", episodeID, err)
		return execFailedHandler{err: err}
	}
	h = eh
	if w.execH == nil {
		w.execH = map[string]episodeTool{}
	}
	w.execH[episodeID] = episodeTool{root: epRoot, h: h}
	return h
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

// episodeHandler returns the episode's cached handler (sandbox is the module
// root: <worktree>/REL, the worktree itself by default), constructing it on
// first use (the constructor runs one summary subprocess, so it is built
// once per episode rather than per Dispatch).
func (w *workspaceTools) episodeHandler(episodeID, sandbox string) (broker.Handler, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if cached, ok := w.handler[episodeID]; ok && cached.root == sandbox {
		return cached.h, nil
	}
	policyDir := filepath.Join(w.stateDir, "policies")
	cacheDir := filepath.Join(w.stateDir, "cache", episodeID)
	for _, dir := range []string{policyDir, cacheDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
	}
	if err := w.linkPackageCache(episodeID, cacheDir); err != nil {
		return nil, err
	}
	if err := w.checkLockCoverage(episodeID, sandbox); err != nil {
		return nil, err
	}
	policyPath := filepath.Join(policyDir, episodeID+".toml")
	policy, err := broker.RenderEpisodePolicy(sandbox)
	if err != nil {
		return nil, err
	}
	if err := writePolicy(policyPath, policy); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), workspaceHandlerBudget)
	defer cancel()
	h, err := broker.NewAilangToolHandler(ctx, broker.AilangToolConfig{
		Bin: w.bin, BinRef: w.binRef, PolicyPath: policyPath, Root: sandbox, CacheDir: cacheDir,
		ExamplesDir: w.examplesDir, RunCaps: w.runCaps, PackageCacheDigest: w.packageCacheDigest,
	})
	if err != nil {
		return nil, err
	}
	tool := verifiedTool{path: w.bin, ref: w.binRef, h: h}
	if w.handler == nil {
		w.handler = map[string]episodeTool{}
	}
	w.handler[episodeID] = episodeTool{root: sandbox, h: tool}
	return tool, nil
}

// runCapsStartupBudget bounds the startup check of the run variants: one
// policy-tool summary per extra capability class (~0.1 s each).
const runCapsStartupBudget = 10 * time.Second

// configureRunCaps applies the --run-* flags (row 135 §4.3 gate 2, §4.6):
// the allowlist must validate, needs the workspace tools, and every variant
// it can produce is rendered once against the workspace root and verified by
// the tool binary's own `policy-tool summary` — a variant the policy layer
// refuses at load (V35) is a startup refusal, not a first-call failure.
func (w *workspaceTools) configureRunCaps(parent context.Context, cfg Config) error {
	op := broker.RunCapsConfig{Allow: cfg.RunAllowCaps, NetAllow: cfg.RunNetAllow, NetAllowHTTP: cfg.RunNetAllowHTTP}
	if len(op.Allow) == 0 && len(op.NetAllow) == 0 && !op.NetAllowHTTP {
		return nil
	}
	if err := op.Validate(); err != nil {
		return err
	}
	if !w.enabled() {
		return fmt.Errorf("--run-allow-caps, --run-net-allow and --run-net-allow-http need the workspace tools " +
			"(--workspace-root and --tool-ailang-bin)")
	}
	checkDir := filepath.Join(w.stateDir, "policies", ".run-caps-check")
	cacheDir := filepath.Join(w.stateDir, "cache", ".run-caps-check")
	for _, dir := range []string{checkDir, cacheDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	defer func() { _ = os.RemoveAll(checkDir) }()
	ctx, cancel := context.WithTimeout(parent, runCapsStartupBudget)
	defer cancel()
	if err := broker.VerifyRunCaps(ctx, broker.AilangToolConfig{
		Bin: w.bin, BinRef: w.binRef, PolicyPath: filepath.Join(checkDir, "root.toml"), Root: w.root, CacheDir: cacheDir,
		RunCaps: op,
	}); err != nil {
		return err
	}
	w.runCaps = op
	return nil
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
