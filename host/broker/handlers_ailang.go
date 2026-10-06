package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// Effect names of the software-engineering domain's tool handlers
// (w-software-engineering-domain §4.2/§4.3). Every one is scoped to the
// symbolic scope WorkspaceScope; the concrete worktree is bound per episode by
// the handler's root, never carried in the request.
const (
	EffectWorkspaceRead  = "Workspace.Read"
	EffectWorkspaceWrite = "Workspace.Write"
	EffectAilangCheck    = "Ailang.Check"
	EffectAilangRun      = "Ailang.Run"
	EffectAilangDiscover = "Ailang.Discover"
	EffectAilangCLI      = "Ailang.CLI"

	// The two extra run effects (row 135, D-135-3 = A): an ailang-run whose
	// caps hold Env is Ailang.RunEnv, one whose caps hold Net Ailang.RunNet;
	// any other run (Declassify included, D-135-5 = A) is Ailang.Run.
	EffectAilangRunEnv = "Ailang.RunEnv"
	EffectAilangRunNet = "Ailang.RunNet"

	// WorkspaceScope is the only scope an AilangToolHandler serves.
	WorkspaceScope = "worktree"
)

// ailangToolExecTimeout is the handler's own wall-clock cap (§4.3 deadline
// budget: the handler gets at most 10 s). episodePolicyTimeoutMS is rendered
// 2 s below it so the AILANG supervisor kills first and the run reports its
// structured limit (V61) instead of a Go-context kill.
const (
	ailangToolExecTimeout  = 10 * time.Second
	episodePolicyTimeoutMS = 8000
)

// episodeDenyDirs and episodeDenyFiles are the protected names (§4.3): CI
// config, agent harness config, the compile cache `run --policy` writes inside
// the worktree (V57), and the git metadata files a worktree does not protect
// by itself (V54).
var (
	episodeDenyDirs  = []string{".github", ".pi", ".claude", ".ailang"}
	episodeDenyFiles = []string{".gitmodules", ".gitattributes"}
)

// episodeDenyWrite is the rendered fs_deny_write list, case-folded in the
// policy because v0.51.0 matches it case-SENSITIVELY while macOS APFS (the
// default volume) resolves names case-insensitively (V68, R-SE-15): under the
// lowercase-only list `.CLAUDE/settings.json` and `.GITMODULES` were written
// through to `.claude/settings.json` and `.gitmodules`. Only `.git` is
// case-folded upstream. (v0.52.1 folds fs_deny_write upstream too — one
// matcher, fileguard.Protection; row-134 design §13 V73 — so on the pinned
// binary this rendering is redundant defence in depth, kept until a separate
// decision shrinks it.) The v0.51.0 matcher (effects.MatchDenyWrite) treats a
// pattern ending `/**` as a LITERAL prefix, so each directory is rendered as
// every fold variant of its name; any other pattern is a path.Match glob
// against the whole relative path and its base name, so each file is one
// character-class pattern. Variants follow Unicode simple case folding, which
// is what APFS was measured to fold onto these letters (V68: only U+017F ſ for
// `s`, besides ASCII case).
var episodeDenyWrite = renderEpisodeDenyWrite()

func renderEpisodeDenyWrite() []string {
	var out []string
	for _, dir := range episodeDenyDirs {
		for _, v := range foldVariants(dir) {
			out = append(out, v+"/**")
		}
	}
	for _, file := range episodeDenyFiles {
		out = append(out, foldClassPattern(file))
	}
	return out
}

// foldOrbit is r's Unicode simple case-folding orbit, starting with r itself.
func foldOrbit(r rune) []rune {
	orbit := []rune{r}
	for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
		orbit = append(orbit, f)
	}
	return orbit
}

// foldVariants is every spelling of name that case-folds to it, name first.
func foldVariants(name string) []string {
	variants := []string{""}
	for _, r := range name {
		orbit := foldOrbit(r)
		next := make([]string, 0, len(variants)*len(orbit))
		for _, prefix := range variants {
			for _, f := range orbit {
				next = append(next, prefix+string(f))
			}
		}
		variants = next
	}
	return variants
}

// foldClassPattern is a path.Match pattern matching exactly the fold variants
// of name (a name with no glob metacharacter).
func foldClassPattern(name string) string {
	var b strings.Builder
	for _, r := range name {
		orbit := foldOrbit(r)
		if len(orbit) == 1 {
			b.WriteRune(r)
			continue
		}
		b.WriteString("[" + string(orbit) + "]")
	}
	return b.String()
}

// fixedToolOps is the policy-tool op allowlist per effect name (§4.3).
// Ailang.CLI's allowlist is the policy summary's `cli` list (V40), read at
// construction; Ailang.Run never reaches policy-tool (it refuses `run`, V16).
var fixedToolOps = map[string][]string{
	EffectWorkspaceRead:  {"read"},
	EffectWorkspaceWrite: {"write", "edit"},
	EffectAilangCheck:    {"ai_check"},
	EffectAilangDiscover: {"builtins_list", "examples_search"},
}

// cliWriteFlags names, per op, the policy-tool flags that make an op write
// the worktree (V67, measured on v0.51.0). policy-tool runs `fmt --write`
// under the episode policy, and the formatter rewrites the file without the
// fs_deny_write check, so `.claude/**` and `.ailang/**` were writable through
// Ailang.CLI alone. The flag's value is ignored (`"write":"false"` still
// writes), so the key's presence is refused. (v0.52.1 refuses a deny-listed
// fmt --write and a "false" value itself, §13 V71/V72; `fmt --write` on an
// ordinary file still writes, so the refusal stays.) The other 19 cli ops admit no
// write-capable flag (the audited table in the tests). Writes go through
// Workspace.Write, never Ailang.CLI.
var cliWriteFlags = map[string][]string{
	"fmt": {"write"},
}

// policyToolRequestKeys are the policy-tool Request fields the handler reads
// (op) or polices (flags). policy-tool decodes into a Go struct, whose field
// match is case-insensitive, so any other spelling of these keys ("OP",
// "FLAGS", "flagſ") would reach it unseen and is refused.
var policyToolRequestKeys = []string{"op", "flags"}

// RenderEpisodePolicy renders the per-episode AILANG operator policy with
// exactly the §4.3 key set (measured admitted on v0.51.0, V39). root is the
// episode worktree, rendered verbatim as fs_sandbox; it must be an absolute,
// clean path with no character that would need TOML escaping.
func RenderEpisodePolicy(root string) ([]byte, error) {
	return renderPolicy(root, baseRunSpec)
}

// PolicyInsideWorkspaceError is the D4-style refusal (§4.3, AC4.5): a policy an
// agent can edit from inside its sandbox is no policy.
type PolicyInsideWorkspaceError struct {
	Path string
	Root string
}

func (e *PolicyInsideWorkspaceError) Error() string {
	return fmt.Sprintf("broker: %q resolves inside the workspace root %q; the policy must live outside it", e.Path, e.Root)
}

// CheckPolicyOutsideRoot refuses a policy (or cache) path that is, lexically or
// after symlink resolution, at or under root. root must exist; path need not.
func CheckPolicyOutsideRoot(path, root string) error {
	if !filepath.IsAbs(path) || !filepath.IsAbs(root) {
		return fmt.Errorf("broker: policy path %q and root %q must both be absolute", path, root)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("broker: resolve workspace root: %w", err)
	}
	candidates := []string{filepath.Clean(path), resolveExistingPrefix(path)}
	for _, candidate := range candidates {
		for _, base := range []string{filepath.Clean(root), canonicalRoot} {
			if within(candidate, base) {
				return &PolicyInsideWorkspaceError{Path: path, Root: root}
			}
		}
	}
	return nil
}

// resolveExistingPrefix resolves symlinks in the longest existing prefix of
// path and re-joins the missing remainder.
func resolveExistingPrefix(path string) string {
	path = filepath.Clean(path)
	rest := ""
	for cur := path; ; {
		if resolved, err := filepath.EvalSymlinks(cur); err == nil {
			return filepath.Join(resolved, rest)
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return path
		}
		rest = filepath.Join(filepath.Base(cur), rest)
		cur = parent
	}
}

func within(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// AilangToolConfig binds one episode's tool handler. Bin is the (archived)
// tool binary and BinRef its content hash; PolicyPath is the rendered policy,
// Root the episode worktree, CacheDir the AILANG_CACHE_DIR, both outside Root.
// ExamplesDir, when set, is the examples corpus `examples_search` reads (passed
// as AILANG_EXAMPLES; also outside Root). The corpus is NOT built into the
// v0.51.0 binary (V65): it reads AILANG_EXAMPLES, then ~/.ailang/examples,
// then CWD-relative examples/ dirs, so an unset ExamplesDir refuses the op
// rather than let it search the agent's own worktree or report nothing.
// v0.52.1 embeds a corpus and skips the CWD and ~/.ailang corpora under
// policy-tool (row-134 design §13 V76); the refusal is kept as the
// stricter-or-equal behaviour until serving the embedded corpus is decided.
//
// RunCaps (row 135) is the operator's run allowlist: the extra capabilities an
// ailang-run may request and where a Net run may connect. Its variants are
// rendered next to PolicyPath (<policy>.run-<caps>.toml), so also outside Root.
type AilangToolConfig struct {
	Bin            string
	BinRef         hashref.HashRef
	PolicyPath     string
	Root           string
	CacheDir       string
	ExamplesDir    string
	ExecTimeout    time.Duration
	MaxOutputBytes int64
	RunCaps        RunCapsConfig
	// PackageCacheDigest (row 141) is the startup digest of the read-only
	// registry package snapshot linked into the episode's HOME. When set it
	// is stamped on every result as "package_cache", next to policy_digest;
	// empty (no snapshot) leaves results byte-identical.
	PackageCacheDigest string
}

// AilangToolHandler serves the six §4.3 effect names through the hardened
// AILANG lane: policy-tool for every op, `run --policy` for Ailang.Run. All
// path, symlink, .git and deny-glob confinement is AILANG's own (V36, V52);
// this handler adds only the per-effect op allowlist and the envelope parse.
type AilangToolHandler struct {
	bin          string
	binRef       hashref.HashRef
	policyPath   string
	root         string
	cacheDir     string
	examplesDir  string // "" = no corpus configured: examples_search is refused
	bounds       handlerBounds
	cliOps       map[string]bool
	policyDigest string
	// packageCacheDigest is AilangToolConfig.PackageCacheDigest.
	packageCacheDigest string

	// Row 135: the operator's run allowlist, the base policy as a run
	// variant, and the per-cap-set variants rendered and verified so far.
	runCaps     RunCapsConfig
	baseVariant runVariant
	variantMu   sync.Mutex
	variants    map[string]runVariant
}

// NewAilangToolHandler validates cfg, refuses a policy or cache inside the
// root, and reads the policy summary once: it must report restricted mode and
// fs_sandbox == Root, and it supplies the Ailang.CLI op list and the digest.
func NewAilangToolHandler(ctx context.Context, cfg AilangToolConfig) (*AilangToolHandler, error) {
	switch {
	case cfg.Bin == "" || !filepath.IsAbs(cfg.Bin):
		return nil, fmt.Errorf("broker: the ailang tool handler requires an absolute tool binary path")
	case cfg.BinRef.IsZero():
		return nil, fmt.Errorf("broker: the ailang tool handler requires the tool binary's hash ref")
	case !filepath.IsAbs(cfg.PolicyPath) || !filepath.IsAbs(cfg.Root) || !filepath.IsAbs(cfg.CacheDir):
		return nil, fmt.Errorf("broker: policy path, root and cache dir must be absolute")
	}
	root, err := filepath.EvalSymlinks(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("broker: resolve workspace root: %w", err)
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("broker: workspace root %q is not a directory", cfg.Root)
	}
	if err := CheckPolicyOutsideRoot(cfg.PolicyPath, root); err != nil {
		return nil, err
	}
	if err := CheckPolicyOutsideRoot(cfg.CacheDir, root); err != nil {
		return nil, err
	}
	if cfg.ExamplesDir != "" {
		if !filepath.IsAbs(cfg.ExamplesDir) {
			return nil, fmt.Errorf("broker: the examples corpus dir must be absolute")
		}
		if err := CheckPolicyOutsideRoot(cfg.ExamplesDir, root); err != nil {
			return nil, err
		}
	}
	if err := cfg.RunCaps.Validate(); err != nil {
		return nil, err
	}
	h := &AilangToolHandler{
		bin: cfg.Bin, binRef: cfg.BinRef, policyPath: cfg.PolicyPath, root: root,
		cacheDir: cfg.CacheDir, examplesDir: cfg.ExamplesDir, runCaps: cfg.RunCaps, packageCacheDigest: cfg.PackageCacheDigest,
		bounds: handlerBounds{
			execTimeout: cfg.ExecTimeout, maxOutputBytes: cfg.MaxOutputBytes,
		},
	}
	if h.bounds.execTimeout <= 0 {
		h.bounds.execTimeout = ailangToolExecTimeout
	}
	h.bounds = h.bounds.normalized()
	if err := h.loadSummary(ctx); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *AilangToolHandler) loadSummary(ctx context.Context) error {
	s, stdout, err := h.summary(ctx, h.policyPath)
	if err != nil {
		if stdout == nil {
			return fmt.Errorf("broker: policy summary: %w", err)
		}
		return fmt.Errorf("broker: policy summary is not an ok summary: %q", stdout)
	}
	if s.SecurityMode != "restricted" {
		return fmt.Errorf("broker: policy security_mode is %q, want \"restricted\"", s.SecurityMode)
	}
	if sandbox, err := filepath.EvalSymlinks(s.FSSandbox); err != nil || sandbox != h.root {
		return fmt.Errorf("broker: policy fs_sandbox %q is not the workspace root %q", s.FSSandbox, h.root)
	}
	if s.PolicyDigest == "" || len(s.CLI) == 0 {
		return fmt.Errorf("broker: policy summary lacks a digest or a cli op list")
	}
	h.policyDigest = s.PolicyDigest
	caps := append([]string(nil), s.Caps...)
	sort.Strings(caps)
	h.baseVariant = runVariant{path: h.policyPath, info: runPolicyInfo{Digest: s.PolicyDigest,
		SecurityMode: s.SecurityMode, Caps: caps, NetAllow: []string{}}}
	h.cliOps = make(map[string]bool, len(s.CLI))
	for _, op := range s.CLI {
		h.cliOps[op] = true
	}
	return nil
}

// AilangToolRefusalError is a handler-level refusal: the request reached a
// handler that will not run it (wrong scope, op outside the effect's
// allowlist, malformed payload). The broker records it as `failed`.
type AilangToolRefusalError struct {
	Effect string
	Op     string
	Why    string
}

func (e *AilangToolRefusalError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("broker: %s refuses op %q: %s", e.Effect, e.Op, e.Why)
	}
	return fmt.Sprintf("broker: %s refused: %s", e.Effect, e.Why)
}

// AilangRunEnvelopeError reports a run whose output matches none of the four
// observed outcome shapes; it is never interpreted.
type AilangRunEnvelopeError struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func (e *AilangRunEnvelopeError) Error() string {
	return fmt.Sprintf("broker: ailang run produced no recognised outcome envelope (rc %d, stdout %q, stderr %q)",
		e.ExitCode, e.Stdout, e.Stderr)
}

func (h *AilangToolHandler) opsFor(effect string) map[string]bool {
	if effect == EffectAilangCLI {
		return h.cliOps
	}
	ops := map[string]bool{}
	for _, op := range fixedToolOps[effect] {
		ops[op] = true
	}
	return ops
}

func (h *AilangToolHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	if req.Scope != WorkspaceScope {
		return nil, &AilangToolRefusalError{Effect: req.Effect, Why: fmt.Sprintf("scope %q is not %q", req.Scope, WorkspaceScope)}
	}
	switch req.Effect {
	case EffectAilangRun, EffectAilangRunEnv, EffectAilangRunNet:
		return h.executeRun(ctx, req.Effect, payload)
	case EffectWorkspaceRead, EffectWorkspaceWrite, EffectAilangCheck, EffectAilangDiscover, EffectAilangCLI:
	default:
		return nil, fmt.Errorf("broker: ailang tool handler does not implement %q", req.Effect)
	}
	fields, err := decodePayloadObject(payload)
	if err != nil {
		return nil, &AilangToolRefusalError{Effect: req.Effect, Why: err.Error()}
	}
	var op string
	for key, value := range fields {
		// policy-tool decodes into a Go struct, whose field match is
		// case-insensitive: an "OP" key would reach it as the op, a "FLAGS"
		// key as the flags. Only the exact key may name what was checked.
		for _, name := range policyToolRequestKeys {
			if key != name && strings.EqualFold(key, name) {
				return nil, &AilangToolRefusalError{Effect: req.Effect, Why: fmt.Sprintf("key %q aliases %q", key, name)}
			}
		}
		if key == "op" {
			if err := json.Unmarshal(value, &op); err != nil {
				return nil, &AilangToolRefusalError{Effect: req.Effect, Why: "op is not a string"}
			}
		}
	}
	if !h.opsFor(req.Effect)[op] {
		return nil, &AilangToolRefusalError{Effect: req.Effect, Op: op, Why: "not in this effect's op allowlist"}
	}
	if why := writeFlagRefusal(op, fields["flags"]); why != "" {
		return nil, &AilangToolRefusalError{Effect: req.Effect, Op: op, Why: why}
	}
	if op == "examples_search" && h.examplesDir == "" {
		// Answered in policy-tool's own refusal shape (ok:false, rc 0), so the
		// caller sees why instead of an empty search.
		refused, _ := json.Marshal(NoExamplesCorpusRefusal)
		return h.stamp(map[string]json.RawMessage{"ok": json.RawMessage(`false`), "refused": refused})
	}
	// Re-encode what was checked, so the bytes policy-tool reads are the
	// bytes the allowlist saw (duplicate keys collapse here, not there).
	body, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("broker: re-encode policy-tool request: %w", err)
	}
	stdout, err := h.policyTool(ctx, body)
	if err != nil {
		return nil, err
	}
	resp, err := decodePayloadObject(stdout)
	if err != nil {
		return nil, fmt.Errorf("broker: policy-tool response: %w (stdout %q)", err, stdout)
	}
	var ok bool
	if raw, present := resp["ok"]; !present || json.Unmarshal(raw, &ok) != nil {
		return nil, fmt.Errorf("broker: policy-tool response has no boolean ok: %q", stdout)
	}
	for _, reserved := range []string{"tool", "policy_digest", "package_cache", "world"} {
		if _, clash := resp[reserved]; clash {
			return nil, fmt.Errorf("broker: policy-tool response carries the reserved key %q", reserved)
		}
	}
	return h.stamp(resp)
}

// writeFlagRefusal is the defence-in-depth twin of the ailang-cli plan's
// refusal (V67): "" when op's flags carry no write-capable flag. The flags
// object has already collapsed duplicate top-level "flags" keys (the map
// decode keeps the last, and the re-encode sends only that one); a flags
// value that is not an object is refused, since it cannot be inspected.
func writeFlagRefusal(op string, rawFlags json.RawMessage) string {
	if rawFlags == nil {
		return ""
	}
	var flags map[string]json.RawMessage
	if err := json.Unmarshal(rawFlags, &flags); err != nil {
		return "flags is not a JSON object"
	}
	for _, name := range cliWriteFlags[op] {
		if _, present := flags[name]; present {
			return fmt.Sprintf("flag %q writes the worktree; Ailang.CLI never writes (use Workspace.Write)", name)
		}
	}
	return ""
}

// NoExamplesCorpusRefusal is the examples_search answer when the daemon was
// started without an examples corpus (no --examples-dir and no
// ~/.ailang/examples at startup).
const NoExamplesCorpusRefusal = "no examples corpus configured: start ailang-worldd serve with --examples-dir DIR " +
	"(an AILANG examples corpus: manifest.json plus runnable/; `ailang examples download` makes one at ~/.ailang/examples)"

// stamp adds the handler's provenance keys to a response object.
func (h *AilangToolHandler) stamp(resp map[string]json.RawMessage) ([]byte, error) {
	resp["tool"], _ = json.Marshal(h.binRef.String())
	resp["policy_digest"], _ = json.Marshal(h.policyDigest)
	if h.packageCacheDigest != "" {
		resp["package_cache"], _ = json.Marshal(h.packageCacheDigest)
	}
	return json.Marshal(resp)
}

func decodePayloadObject(payload []byte) (map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(payload))
	if err := dec.Decode(&fields); err != nil || fields == nil {
		return nil, errors.New("payload is not a JSON object")
	}
	if dec.More() {
		return nil, errors.New("payload has trailing data after the JSON object")
	}
	return fields, nil
}

// childEnv is the whole child environment of both branches: minimal, never
// the parent's (so the registry credential cannot reach a tool), with the
// compile cache pointed outside the worktree (honoured by policy-tool; ignored
// by `run --policy`, V57 / R-SE-4, which the .ailang/** deny glob covers).
//
// AILANG_EXAMPLES (only when a corpus is configured) points examples_search at
// the operator's corpus: HOME is the per-episode cache, so the binary's
// ~/.ailang/examples fallback never finds the real one (V65).
func (h *AilangToolHandler) childEnv() []string {
	env := []string{
		"HOME=" + h.cacheDir,
		"PATH=/usr/bin:/bin",
		"LANG=C",
		"LC_ALL=C",
		"AILANG_CACHE_DIR=" + h.cacheDir,
	}
	if h.examplesDir != "" {
		env = append(env, "AILANG_EXAMPLES="+h.examplesDir)
	}
	return env
}

// policyTool is execution branch 1: the request on stdin, cwd = the worktree.
// policy-tool answers refusals as ok:false JSON with rc 0; any non-zero rc is
// a handler failure.
func (h *AilangToolHandler) policyTool(ctx context.Context, request []byte) ([]byte, error) {
	return h.policyToolWith(ctx, h.policyPath, request)
}

// policyToolWith runs policy-tool under policyPath (the base policy, or a run
// variant being verified by its summary).
func (h *AilangToolHandler) policyToolWith(ctx context.Context, policyPath string, request []byte) ([]byte, error) {
	stderr := newBoundedSink(h.bounds.maxOutputBytes)
	stdout, err := runBounded(ctx, h.bounds, handlerCommand{
		path: h.bin, args: []string{"policy-tool", "--policy", policyPath},
		dir: h.root, env: h.childEnv(), stdin: request, stderr: stderr,
	})
	if err != nil {
		return nil, fmt.Errorf("%w (stderr %q)", err, stderr.Bytes())
	}
	return stdout, nil
}

// executeRun is execution branch 2 (row 134 §4.3, row 135 §4.5): `ailang
// run --policy P [--args-json J] -- <path> [-- <argv…>]` with the payload's
// stdin piped in. P is the episode's base policy for an IO/FS run, otherwise
// the cap set's verified variant. The payload is checked again here (the
// plan's twin), the caps must select the effect requested, and every cap
// beyond IO/FS must be enabled by the operator — all before any subprocess.
func (h *AilangToolHandler) executeRun(ctx context.Context, effect string, payload []byte) ([]byte, error) {
	req, err := parseRunPayload(effect, payload)
	if err != nil {
		return nil, err
	}
	if why := h.operatorRefusal(req.caps); why != "" {
		return nil, &AilangToolRefusalError{Effect: effect, Why: why}
	}
	variant, err := h.runVariantFor(ctx, req.caps)
	if err != nil {
		return nil, err
	}
	var stdin []byte
	if req.stdin != nil {
		stdin = []byte(*req.stdin)
	}
	stderr := newBoundedSink(h.bounds.maxOutputBytes)
	stdout, err := runBounded(ctx, h.bounds, handlerCommand{
		path: h.bin, args: runArgv(variant.path, req), dir: h.root, env: h.childEnv(), stdin: stdin, stderr: stderr,
	})
	exitCode := 0
	if err != nil {
		var exitErr *HandlerExitError
		var procErr *exec.ExitError
		if !errors.As(err, &exitErr) || !errors.As(exitErr.Err, &procErr) || procErr.ExitCode() < 0 {
			return nil, fmt.Errorf("%w (stderr %q)", err, stderr.Bytes())
		}
		exitCode, stdout = procErr.ExitCode(), exitErr.Output
	}
	result, err := composeRunOutcome(exitCode, stdout, stderr.Bytes())
	if err != nil {
		return nil, err
	}
	info := variant.info
	result.Policy = &info
	result.PackageCache = h.packageCacheDigest
	return json.Marshal(result)
}

// runResult is the Ailang.Run handler output (§4.3).
type runResult struct {
	Admitted bool            `json:"admitted"`
	ExitCode int             `json:"exit_code"`
	Decision json.RawMessage `json:"decision"`
	Limit    json.RawMessage `json:"limit"`
	Stdout   string          `json:"stdout"`
	Stderr   string          `json:"stderr"`
	// Policy (row 135 §4.3) names the verified policy the run executed
	// under; set by executeRun, absent from the pure composition.
	Policy *runPolicyInfo `json:"policy,omitempty"`
	// PackageCache (row 141) is the read-only registry snapshot's startup
	// digest, as stamp adds it to the policy-tool results; omitted without one.
	PackageCache string `json:"package_cache,omitempty"`
}

const (
	policyLinePrefix       = "policy: "
	policyResultLinePrefix = "policy-result: "
)

// composeRunResult maps a finished run onto exactly one of the four outcome
// shapes measured on the v0.51.0 tool binary, with no default:
//
//	(a) admitted: any rc, exactly one `policy:` line on stderr (V33, V39)
//	(b) admitted then refused at runtime: rc 1, same stderr line (V34) — the
//	    same parse as (a); the program's own rc is reported, never judged
//	(c) refused before execution: any non-zero rc (2 for a policy violation
//	    or type error, 1 for an unreadable entry file), a not-ok decision
//	    JSON on stdout, stderr empty (V49, V63); the decision carries the
//	    binary's error_kind (e.g. read_failed) verbatim
//	(d) supervisor limit: rc 3, the `policy:` line plus one `policy-result:`
//	    line carrying `reason` and `stage` (V61)
//
// Anything else is an *AilangRunEnvelopeError. rc alone never signals success:
// a refused inner effect can exit 0 (V35), so output is returned verbatim.
func composeRunResult(exitCode int, stdout, stderr []byte) ([]byte, error) {
	result, err := composeRunOutcome(exitCode, stdout, stderr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func composeRunOutcome(exitCode int, stdout, stderr []byte) (runResult, error) {
	unrecognised := &AilangRunEnvelopeError{ExitCode: exitCode, Stdout: string(stdout), Stderr: string(stderr)}
	var policyLines, resultLines []string
	for _, line := range strings.Split(string(stderr), "\n") {
		switch {
		case strings.HasPrefix(line, policyLinePrefix):
			policyLines = append(policyLines, strings.TrimPrefix(line, policyLinePrefix))
		case strings.HasPrefix(line, policyResultLinePrefix):
			resultLines = append(resultLines, strings.TrimPrefix(line, policyResultLinePrefix))
		}
	}
	result := runResult{ExitCode: exitCode, Stdout: string(stdout), Stderr: string(stderr), Limit: json.RawMessage("null")}
	switch {
	case len(policyLines) == 1 && len(resultLines) <= 1:
		var line struct {
			OK       bool            `json:"ok"`
			Decision json.RawMessage `json:"decision"`
		}
		if json.Unmarshal([]byte(policyLines[0]), &line) != nil || !line.OK {
			return runResult{}, unrecognised
		}
		ok, valid := decisionOK(line.Decision)
		if !valid || !ok {
			return runResult{}, unrecognised
		}
		if len(resultLines) == 1 {
			// (d): only with rc 3 and a parseable reason.
			var limit struct {
				Reason string `json:"reason"`
				Stage  string `json:"stage"`
			}
			if exitCode != 3 || json.Unmarshal([]byte(resultLines[0]), &limit) != nil || limit.Reason == "" || limit.Stage == "" {
				return runResult{}, unrecognised
			}
			result.Limit = json.RawMessage(resultLines[0])
		}
		result.Admitted, result.Decision = true, line.Decision
	case len(policyLines) == 0 && len(resultLines) == 0 && exitCode != 0 && len(stderr) == 0:
		var static struct {
			Decision json.RawMessage `json:"decision"`
		}
		if json.Unmarshal(stdout, &static) != nil {
			return runResult{}, unrecognised
		}
		ok, valid := decisionOK(static.Decision)
		if !valid || ok {
			return runResult{}, unrecognised
		}
		result.Admitted, result.Decision = false, static.Decision
	default:
		return runResult{}, unrecognised
	}
	return result, nil
}

// decisionOK reads a decision object's boolean ok; valid is false when the
// decision is absent, not an object, or carries no boolean ok.
func decisionOK(raw json.RawMessage) (ok bool, valid bool) {
	var d map[string]json.RawMessage
	if json.Unmarshal(raw, &d) != nil || d == nil {
		return false, false
	}
	value, present := d["ok"]
	if !present || json.Unmarshal(value, &ok) != nil {
		return false, false
	}
	return ok, true
}
