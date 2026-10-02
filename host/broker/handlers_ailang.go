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
	"strings"
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

// episodeDenyWrite is the rendered fs_deny_write list (§4.3): CI config, agent
// harness config, the git metadata files a worktree does not protect by itself
// (V54), and the compile cache `run --policy` writes inside the worktree (V57).
var episodeDenyWrite = []string{
	".github/**", ".pi/**", ".claude/**", ".gitmodules", ".gitattributes", ".ailang/**",
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

// RenderEpisodePolicy renders the per-episode AILANG operator policy with
// exactly the §4.3 key set (measured admitted on v0.51.0, V39). root is the
// episode worktree, rendered verbatim as fs_sandbox; it must be an absolute,
// clean path with no character that would need TOML escaping.
func RenderEpisodePolicy(root string) ([]byte, error) {
	if !filepath.IsAbs(root) || filepath.Clean(root) != root {
		return nil, fmt.Errorf("broker: episode policy root %q is not a clean absolute path", root)
	}
	for _, r := range root {
		if r == '"' || r == '\\' || r == unicode.ReplacementChar || !unicode.IsPrint(r) {
			return nil, fmt.Errorf("broker: episode policy root %q has a character TOML would need escaped", root)
		}
	}
	quoted := make([]string, len(episodeDenyWrite))
	for i, glob := range episodeDenyWrite {
		quoted[i] = `"` + glob + `"`
	}
	var b strings.Builder
	b.WriteString(`security_mode = "restricted"` + "\n")
	b.WriteString(`allowed_caps = ["IO", "FS"]` + "\n")
	b.WriteString(`fs_sandbox = "` + root + `"` + "\n")
	fmt.Fprintf(&b, "timeout_ms = %d\n", episodePolicyTimeoutMS)
	b.WriteString("fs_deny_write = [" + strings.Join(quoted, ", ") + "]\n")
	b.WriteString(`entry = "main"` + "\n")
	b.WriteString("\n[budgets]\nFS = 1000\n")
	return []byte(b.String()), nil
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
type AilangToolConfig struct {
	Bin            string
	BinRef         hashref.HashRef
	PolicyPath     string
	Root           string
	CacheDir       string
	ExecTimeout    time.Duration
	MaxOutputBytes int64
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
	bounds       handlerBounds
	cliOps       map[string]bool
	policyDigest string
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
	h := &AilangToolHandler{
		bin: cfg.Bin, binRef: cfg.BinRef, policyPath: cfg.PolicyPath, root: root,
		cacheDir: cfg.CacheDir,
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
	stdout, err := h.policyTool(ctx, []byte(`{"op":"summary"}`))
	if err != nil {
		return fmt.Errorf("broker: policy summary: %w", err)
	}
	var resp struct {
		OK      bool `json:"ok"`
		Summary struct {
			SecurityMode string   `json:"security_mode"`
			PolicyDigest string   `json:"policy_digest"`
			FSSandbox    string   `json:"fs_sandbox"`
			CLI          []string `json:"cli"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(stdout, &resp); err != nil || !resp.OK {
		return fmt.Errorf("broker: policy summary is not an ok summary: %q", stdout)
	}
	s := resp.Summary
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
	case EffectAilangRun:
		return h.executeRun(ctx, payload)
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
		// case-insensitive: an "OP" key would reach it as the op. Only the
		// exact key may name the op that the allowlist checked.
		if key != "op" && strings.EqualFold(key, "op") {
			return nil, &AilangToolRefusalError{Effect: req.Effect, Why: fmt.Sprintf("key %q aliases \"op\"", key)}
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
	for _, reserved := range []string{"tool", "policy_digest", "world"} {
		if _, clash := resp[reserved]; clash {
			return nil, fmt.Errorf("broker: policy-tool response carries the reserved key %q", reserved)
		}
	}
	resp["tool"], _ = json.Marshal(h.binRef.String())
	resp["policy_digest"], _ = json.Marshal(h.policyDigest)
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
func (h *AilangToolHandler) childEnv() []string {
	return []string{
		"HOME=" + h.cacheDir,
		"PATH=/usr/bin:/bin",
		"LANG=C",
		"LC_ALL=C",
		"AILANG_CACHE_DIR=" + h.cacheDir,
	}
}

// policyTool is execution branch 1: the request on stdin, cwd = the worktree.
// policy-tool answers refusals as ok:false JSON with rc 0; any non-zero rc is
// a handler failure.
func (h *AilangToolHandler) policyTool(ctx context.Context, request []byte) ([]byte, error) {
	stderr := newBoundedSink(h.bounds.maxOutputBytes)
	stdout, err := runBounded(ctx, h.bounds, handlerCommand{
		path: h.bin, args: []string{"policy-tool", "--policy", h.policyPath},
		dir: h.root, env: h.childEnv(), stdin: request, stderr: stderr,
	})
	if err != nil {
		return nil, fmt.Errorf("%w (stderr %q)", err, stderr.Bytes())
	}
	return stdout, nil
}

// executeRun is execution branch 2: `ailang run --policy P [--args-json J] --
// <path>`. The `--` makes a hyphen-leading path a file, never a flag (V58).
func (h *AilangToolHandler) executeRun(ctx context.Context, payload []byte) ([]byte, error) {
	fields, err := decodePayloadObject(payload)
	if err != nil {
		return nil, &AilangToolRefusalError{Effect: EffectAilangRun, Why: err.Error()}
	}
	var path, argsJSON string
	for key, value := range fields {
		switch key {
		case "path":
			if json.Unmarshal(value, &path) != nil || path == "" {
				return nil, &AilangToolRefusalError{Effect: EffectAilangRun, Why: "path must be a non-empty string"}
			}
		case "args_json":
			if json.Unmarshal(value, &argsJSON) != nil || !json.Valid([]byte(argsJSON)) {
				return nil, &AilangToolRefusalError{Effect: EffectAilangRun, Why: "args_json must be a string holding JSON"}
			}
		default:
			return nil, &AilangToolRefusalError{Effect: EffectAilangRun, Why: fmt.Sprintf("unknown payload key %q", key)}
		}
	}
	if path == "" {
		return nil, &AilangToolRefusalError{Effect: EffectAilangRun, Why: "path is required"}
	}
	args := []string{"run", "--policy", h.policyPath}
	if argsJSON != "" {
		args = append(args, "--args-json", argsJSON)
	}
	args = append(args, "--", path)

	stderr := newBoundedSink(h.bounds.maxOutputBytes)
	stdout, err := runBounded(ctx, h.bounds, handlerCommand{
		path: h.bin, args: args, dir: h.root, env: h.childEnv(), stderr: stderr,
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
	return composeRunResult(exitCode, stdout, stderr.Bytes())
}

// runResult is the Ailang.Run handler output (§4.3).
type runResult struct {
	Admitted bool            `json:"admitted"`
	ExitCode int             `json:"exit_code"`
	Decision json.RawMessage `json:"decision"`
	Limit    json.RawMessage `json:"limit"`
	Stdout   string          `json:"stdout"`
	Stderr   string          `json:"stderr"`
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
			return nil, unrecognised
		}
		ok, valid := decisionOK(line.Decision)
		if !valid || !ok {
			return nil, unrecognised
		}
		if len(resultLines) == 1 {
			// (d): only with rc 3 and a parseable reason.
			var limit struct {
				Reason string `json:"reason"`
				Stage  string `json:"stage"`
			}
			if exitCode != 3 || json.Unmarshal([]byte(resultLines[0]), &limit) != nil || limit.Reason == "" || limit.Stage == "" {
				return nil, unrecognised
			}
			result.Limit = json.RawMessage(resultLines[0])
		}
		result.Admitted, result.Decision = true, line.Decision
	case len(policyLines) == 0 && len(resultLines) == 0 && exitCode != 0 && len(stderr) == 0:
		var static struct {
			Decision json.RawMessage `json:"decision"`
		}
		if json.Unmarshal(stdout, &static) != nil {
			return nil, unrecognised
		}
		ok, valid := decisionOK(static.Decision)
		if !valid || ok {
			return nil, unrecognised
		}
		result.Admitted, result.Decision = false, static.Decision
	default:
		return nil, unrecognised
	}
	return json.Marshal(result)
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
