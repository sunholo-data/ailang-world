package broker

// Row 140 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.1,
// §4.2): the Workspace.Exec effect, which runs one command of the operator's
// exec profile in the episode worktree under the sandbox runtime. M1 bound
// the name (R8: workspace-exec declares it, so it is always bound) to the
// typed unconfigured refusal below, which stays the behaviour when no
// profile is configured. M2 adds ExecHandler (§4.5): it matches the call
// against the profile (MatchExecArgs), re-verifies the confinement stack and
// runs the command through srt by way of runBounded — the one process
// lifecycle; this file owns no exec.Cmd, kill or procbound code
// (TestExecHandlerOwnsNoProcessLifecycle).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// EffectWorkspaceExec is the effect name of the workspace-exec tool, scoped
// to WorkspaceScope like every software-engineering effect.
const EffectWorkspaceExec = "Workspace.Exec"

// NoExecProfileRefusal is the Workspace.Exec answer when the daemon has no
// exec profile (§4.2 gate 2), in examples_search's no-corpus shape (V31).
const NoExecProfileRefusal = "no exec profile configured: start ailang-worldd serve with --exec-profile FILE"

// ExecUnconfiguredHandler serves Workspace.Exec when no exec profile is
// configured: it answers {ok:false, refused: NoExecProfileRefusal} and runs
// nothing. Any other effect or scope is a handler failure.
type ExecUnconfiguredHandler struct{}

func (ExecUnconfiguredHandler) Execute(_ context.Context, req EffectRequest, _ []byte) ([]byte, error) {
	if req.Effect != EffectWorkspaceExec {
		return nil, fmt.Errorf("broker: exec handler does not implement %q", req.Effect)
	}
	if req.Scope != WorkspaceScope {
		return nil, fmt.Errorf("broker: exec handler: scope %q is not %q", req.Scope, WorkspaceScope)
	}
	return execRefusal(NoExecProfileRefusal)
}

// ExecNoEpisodeProfileRefusal is the Workspace.Exec answer for an episode
// whose credential binds no profile and that no `--exec-episode-project`
// names while several profiles are configured (row 140 M3, §4.6): World never
// guesses a project. Row 189 M2: the fix that needs no restart is a credential
// bound at the mint.
func ExecNoEpisodeProfileRefusal(episode string) string {
	return fmt.Sprintf("no exec profile for episode %q: its session credential binds none and several profiles are configured; "+
		"mint the episode's session with --profile PROJECT (ailang-worldd episode start or session new; no restart needed), "+
		"or start ailang-worldd serve with --exec-episode-project %s=PROJECT", episode, episode)
}

// ExecProfileDigestMismatchRefusal is the Workspace.Exec answer for a
// credential bound at the mint to profile bytes other than the ones the
// daemon loaded (row 189 M2, D-189-2 A): the human confirmed the minted
// digest and no other, so every exec refuses until the session is re-minted.
func ExecProfileDigestMismatchRefusal(episode, project, minted, loaded string) string {
	return fmt.Sprintf("exec profile %q changed since episode %q was minted: the session credential is bound to %s, "+
		"the daemon loaded %s; re-mint the session against the current profile (ailang-worldd session revoke --episode %s, "+
		"then episode start or session new --profile %s)", project, episode, minted, loaded, episode, project)
}

// ExecProfileNotLoadedRefusal is the Workspace.Exec answer for an episode
// bound to a project the daemon has no profile for. why is the reason the
// profiles directory left it out, or "" when it was never configured.
func ExecProfileNotLoadedRefusal(episode, project, why string) string {
	if why != "" {
		return fmt.Sprintf("exec profile %q of episode %q was left out at serve startup: %s", project, episode, why)
	}
	return fmt.Sprintf("episode %q is bound to exec profile %q, which this daemon has not loaded; "+
		"put %s.json in the serve --exec-profiles-dir and restart serve once", episode, project, project)
}

// ExecRefusalHandler answers every Workspace.Exec call {ok:false, refused:
// Why} and runs nothing, like ExecUnconfiguredHandler.
type ExecRefusalHandler struct{ Why string }

func (r ExecRefusalHandler) Execute(_ context.Context, req EffectRequest, _ []byte) ([]byte, error) {
	if req.Effect != EffectWorkspaceExec {
		return nil, fmt.Errorf("broker: exec handler does not implement %q", req.Effect)
	}
	if req.Scope != WorkspaceScope {
		return nil, fmt.Errorf("broker: exec handler: scope %q is not %q", req.Scope, WorkspaceScope)
	}
	return execRefusal(r.Why)
}

func execRefusal(why string) ([]byte, error) {
	return json.Marshal(map[string]any{"ok": false, "refused": why})
}

// DefaultExecMaxOutputBytes is the per-stream runaway-output kill (§4.5,
// `--exec-max-output-bytes`, M3): V27 saw 1.47 MB from a full vitest run.
const DefaultExecMaxOutputBytes = 64 << 20

// execTrampoline runs inside the sandbox as `/bin/sh -c <it> world-exec
// K=V… argv…`: `env -i` replaces srt's environment with exactly World's
// K=V set (V41), and the shell, not exec'd away, turns a signal death into
// 128+n (V11; bare srt reports SIGTERM as 0 on macOS).
const execTrampoline = `/usr/bin/env -i "$@"; exit $?`

// execHostEnv is the whole environment of the host-side node that runs srt:
// it is outside the sandbox, so nothing of the profile's may reach it (V40).
func execHostEnv() []string { return []string{"PATH=/usr/bin:/bin"} }

// ExecHandlerConfig binds one episode's Workspace.Exec to one profile.
type ExecHandlerConfig struct {
	Tools         *ExecToolState // nil retains the fixed stack (also used by probe scratch handlers)
	Profile       *ExecProfile
	Sandbox       *ExecSandbox
	Episode       string
	Worktree      string // the canonical episode worktree
	WorkspaceRoot string // canonical; denied for reading, so siblings are unreadable
	StateDir      string // canonical World state dir; denied for reading
	// OperatorHome is the canonical operator $HOME ("" = os.UserHomeDir):
	// denied for reading, and where srt anchors two default write paths.
	OperatorHome string
	// MaxOutputBytes is the per-stream hard cap; 0 = DefaultExecMaxOutputBytes.
	MaxOutputBytes int64
}

// ExecHandler serves Workspace.Exec for one episode × project (§4.5).
type ExecHandler struct {
	tools          *ExecToolState
	profile        *ExecProfile
	sandbox        *ExecSandbox
	worktree       string
	cache          string // <state>/exec-cache/<ep>/<project>
	settingsPath   string // <state>/exec/<ep>.<project>.srt.json
	settings       []byte // the canonical rendering
	settingsDigest string
	maxOutput      int64
	seedMu         sync.Mutex
	baseline       *protectedBaseline
	// spawns counts cache seeding and srt, excluding baseline and phase Git calls.
	// A refusal or a drifted stack leaves it unchanged (AC2.4, AC2.9).
	spawns atomic.Int64
}

// NewExecHandler renders the episode's srt settings and verifies the file it
// wrote (§4.4 "Verification at render").
func NewExecHandler(cfg ExecHandlerConfig) (*ExecHandler, error) {
	if cfg.Profile == nil || cfg.Sandbox == nil {
		return nil, fmt.Errorf("broker: exec handler needs a profile and a verified sandbox")
	}
	// The daemon's episode grammar, which the project grammar shares: both
	// name one path component of the cache and settings paths.
	if !execProjectPattern.MatchString(cfg.Episode) {
		return nil, fmt.Errorf("broker: exec handler: episode id %q is outside the grammar", cfg.Episode)
	}
	home := cfg.OperatorHome
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("broker: exec handler: operator HOME: %w", err)
		}
		if home, err = filepath.EvalSymlinks(h); err != nil {
			return nil, fmt.Errorf("broker: exec handler: operator HOME: %w", err)
		}
	}
	readRoots := []string{}
	for _, r := range cfg.Profile.ReadRoots {
		real, err := filepath.EvalSymlinks(r)
		if err != nil {
			return nil, fmt.Errorf("broker: exec handler: read_root %q: %w", r, err)
		}
		readRoots = append(readRoots, r, real)
	}
	h := &ExecHandler{tools: cfg.Tools, profile: cfg.Profile, sandbox: cfg.Sandbox, worktree: cfg.Worktree, maxOutput: cfg.MaxOutputBytes,
		cache:        filepath.Join(cfg.StateDir, "exec-cache", cfg.Episode, cfg.Profile.Project),
		settingsPath: filepath.Join(cfg.StateDir, "exec", cfg.Episode+"."+cfg.Profile.Project+".srt.json")}
	if h.maxOutput <= 0 {
		h.maxOutput = DefaultExecMaxOutputBytes
	}
	data, err := renderExecSettings(execSettingsSpec{Worktree: cfg.Worktree, Cache: h.cache, OperatorHome: home,
		StateDir: cfg.StateDir, WorkspaceRoot: cfg.WorkspaceRoot, ReadRoots: dedupe(readRoots), Protected: cfg.Profile.Protected,
		SeccompDir: cfg.Sandbox.SeccompDir(hostGOOS(), hostGOARCH()), GOOS: hostGOOS()})
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(h.settingsPath), 0o700); err != nil {
		return nil, err
	}
	if h.settingsDigest, err = writeExecSettings(h.settingsPath, data); err != nil {
		return nil, err
	}
	h.settings = data
	if len(cfg.Profile.Protected) > 0 {
		// The handler is built inside the ctx-less BinderFor; the capture's
		// one construction bound is fed directly from this root.
		ctx, cancel := context.WithTimeout(context.Background(), protectedBaselineTimeout)
		baseline := captureProtectedBaseline(ctx, "", cfg.Worktree, cfg.Profile.Protected)
		cancel()
		h.baseline = &baseline
	}
	return h, nil
}

// execResult is the handler output and the record's `output` (§4.5).
type execResult struct {
	ExitCode        *int              `json:"exit_code"`
	TimedOut        bool              `json:"timed_out"`
	Limit           *string           `json:"limit"`
	DurationMS      int64             `json:"duration_ms"`
	Stdout          string            `json:"stdout"`
	StdoutBytes     int64             `json:"stdout_bytes"`
	StdoutTruncated bool              `json:"stdout_truncated"`
	StdoutSHA256    string            `json:"stdout_sha256"`
	Stderr          string            `json:"stderr"`
	StderrBytes     int64             `json:"stderr_bytes"`
	StderrTruncated bool              `json:"stderr_truncated"`
	StderrSHA256    string            `json:"stderr_sha256"`
	Argv            []string          `json:"argv"`
	Profile         execResultProfile `json:"profile"`
	Sandbox         execResultSandbox `json:"sandbox"`
	// WorktreeDelta (row 184) is what the command changed in the episode
	// worktree; omitted when nothing changed (byte-identical results).
	WorktreeDelta *WorktreeDelta   `json:"worktree_delta,omitempty"`
	Protected     *protectedRecord `json:"protected,omitempty"`
}

type execResultProfile struct {
	Argv0       string `json:"argv0"`
	Argv0Target string `json:"argv0_target"`
	Project     string `json:"project"`
	Digest      string `json:"digest"`
	Command     string `json:"command"`
	Argv0SHA256 string `json:"argv0_sha256"`
}

type execResultSandbox struct {
	Node           string `json:"node"`
	NodeTarget     string `json:"node_target"`
	NodeSHA256     string `json:"node_sha256"`
	Runtime        string `json:"runtime"`
	Version        string `json:"version"`
	CLISHA256      string `json:"cli_sha256"`
	SettingsDigest string `json:"settings_digest"`
	ReadFence      bool   `json:"read_fence"`
	Network        string `json:"network"`
}

func (h *ExecHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	if req.Effect != EffectWorkspaceExec {
		return nil, fmt.Errorf("broker: exec handler does not implement %q", req.Effect)
	}
	if req.Scope != WorkspaceScope {
		return nil, fmt.Errorf("broker: exec handler: scope %q is not %q", req.Scope, WorkspaceScope)
	}
	// Gate 2 (§4.2): the profile and the grammar, before anything runs.
	id, args, why := parseExecPayload(payload)
	if why != "" {
		return execRefusal(why)
	}
	cmd, ok := h.profile.Commands[id]
	if !ok {
		return execRefusal(fmt.Sprintf("command %q is not a command of exec profile %q", id, h.profile.Project))
	}
	normalized, err := MatchExecArgs(cmd.ExecCommand, args)
	if err != nil {
		return execRefusal(fmt.Sprintf("command %q: %v", id, err))
	}
	sandbox := h.sandbox
	if h.tools != nil {
		var profile *ExecProfile
		var why string
		sandbox, profile, why, err = h.tools.resolve(ctx, h.profile.Project, h.worktree)
		if err != nil {
			return nil, err
		}
		if why != "" {
			return execRefusal(why)
		}
		cmd = profile.Commands[id]
	}
	emitted := EmitExecArgs(cmd.ExecCommand, normalized)
	// §4.6: the stack that sets up the fence is re-verified before EACH call;
	// a mismatch is a failure naming the file, with nothing spawned.
	if err := sandbox.Verify(); err != nil {
		return nil, err
	}
	if err := verifyFileDigest(cmd.Argv0, cmd.Argv0SHA256); err != nil {
		return nil, err
	}
	if err := verifyExecSettings(h.settingsPath, h.settings); err != nil {
		return nil, err
	}
	cwd, err := h.projectDir()
	if err != nil {
		return nil, err
	}
	if err := h.prepareCache(ctx); err != nil {
		return nil, err
	}
	deadline, hasDeadline := ctx.Deadline()
	var before protectedCheck
	var reserve time.Duration
	if h.baseline != nil && h.baseline.err == nil {
		before = h.checkProtected(ctx)
		reserve = protectedCheckTimeout
	}
	budget := execChildBudget(time.Duration(h.profile.TimeoutMS)*time.Millisecond, deadline, hasDeadline, time.Now(), reserve)
	if budget <= 0 {
		return nil, fmt.Errorf("broker: exec handler: the call's budget is spent before %q could start", id)
	}
	launch := h.launch(append(append([]string{cmd.Argv0}, cmd.Argv[1:]...), emitted...), true)
	capture := newCaptureHeadTail(execHeadBytes, execTailBytes, h.maxOutput)
	// Row 184 D1: the worktree before the command, outside its budget.
	probe := beginWorktreeDelta(ctx, h.worktree)
	start := time.Now()
	h.spawns.Add(1)
	_, runErr := runBounded(ctx, handlerBounds{execTimeout: budget, maxOutputBytes: h.maxOutput}, handlerCommand{
		path: sandbox.Node, args: launch, dir: cwd, env: execHostEnv(), capture: capture,
	})
	if code, exited := handlerExitStatus(runErr); exited && code == 134 && !cleanupIncomplete(runErr) {
		if err := checkExecNodeVersion(ctx, sandbox.Node); err != nil {
			if ctx.Err() != nil {
				return execRefusal("the re-check after node moved was interrupted; retry")
			}
			why := err.Error()
			if h.tools != nil {
				return execRefusal(h.tools.nodeFailed(sandbox, why))
			}
			t := execAcceptedTool{name: "node", configured: sandbox.Node}
			next := execToolNext(t.configured, sandbox.Node, "", "", false)
			return execRefusal(execToolF15(t, sandbox.Node, why, "--version failed", next))
		}
	}
	configuredArgv0 := cmd.ConfiguredArgv0
	if configuredArgv0 == "" {
		configuredArgv0 = cmd.Argv0
	}
	configuredNode := sandbox.NodeConfigured
	if configuredNode == "" {
		configuredNode = sandbox.Node
	}
	res := execResult{DurationMS: time.Since(start).Milliseconds(),
		Argv: append(append([]string(nil), cmd.Argv...), emitted...),
		Profile: execResultProfile{Project: h.profile.Project, Digest: h.profile.Digest, Command: id,
			Argv0: configuredArgv0, Argv0Target: cmd.Argv0, Argv0SHA256: cmd.Argv0SHA256},
		Sandbox: execResultSandbox{Node: configuredNode, NodeTarget: sandbox.Node, NodeSHA256: sandbox.NodeSHA256,
			Runtime: ExecSandboxPackage, Version: sandbox.Version, CLISHA256: sandbox.CLISHA256,
			SettingsDigest: h.settingsDigest, ReadFence: true, Network: "none"}}
	if cleanupIncomplete(runErr) {
		return nil, runErr
	}
	switch code, exited := handlerExitStatus(runErr); {
	case runErr == nil:
		zero := 0
		res.ExitCode = &zero
	case errors.Is(runErr, ErrHandlerTimeout):
		res.TimedOut = true
	case errors.Is(runErr, ErrHandlerOverflow):
		limit := "output"
		res.Limit = &limit
	case exited:
		res.ExitCode = &code
	default:
		return nil, runErr
	}
	out, errOut := capture.stdout.snapshot(), capture.stderr.snapshot()
	res.Stdout, res.StdoutBytes, res.StdoutTruncated, res.StdoutSHA256 = out.Text, out.Bytes, out.Truncated, out.SHA256
	res.Stderr, res.StderrBytes, res.StderrTruncated, res.StderrSHA256 = errOut.Text, errOut.Bytes, errOut.Truncated, errOut.SHA256
	if h.baseline != nil {
		var after protectedCheck
		if h.baseline.err == nil {
			afterCtx := context.WithoutCancel(ctx)
			if hasDeadline {
				var cancel context.CancelFunc
				afterCtx, cancel = context.WithDeadline(afterCtx, deadline)
				defer cancel()
			}
			after = h.checkProtected(afterCtx)
		}
		record := buildProtectedRecord(*h.baseline, h.profile.Protected, protectedEnforced(hostGOOS()), before, after)
		res.Protected = &record
	}
	res.WorktreeDelta = probe.finish(ctx)
	return json.Marshal(res)
}

// protectedEnforced names what the sandbox enforces on goos (design §4.5):
// srt denies both pattern forms on darwin; Linux renders only the directory
// form, so base-name matches are detected, not prevented (D-156-3).
func protectedEnforced(goos string) string {
	if goos == "linux" {
		return "dirs-only"
	}
	return "glob"
}

// execChildBudget reserves the after-phase inside the original call deadline.
func execChildBudget(timeout time.Duration, deadline time.Time, hasDeadline bool, now time.Time, reserve time.Duration) time.Duration {
	if hasDeadline {
		return min(timeout, deadline.Sub(now)-reserve)
	}
	return timeout
}

// checkProtected resolves Git at call time; resolution failures belong to this phase.
func (h *ExecHandler) checkProtected(ctx context.Context) protectedCheck {
	git, err := worktreeGit()
	if err != nil {
		return protectedCheck{err: err}
	}
	return checkProtected(ctx, git, h.worktree, h.baseline.base, h.profile.Protected, h.baseline.ignored)
}

// launch is the §4.5 launch line the host node runs for argv: srt's own
// options closed by "--" (V12), then, with trampoline, the in-sandbox `env
// -i` trampoline carrying World's K=V set (V41). Only the startup probe's
// proxy leg (exec_probe.go) runs without the trampoline, to keep srt's
// proxy environment.
func (h *ExecHandler) launch(argv []string, trampoline bool) []string {
	line := []string{h.sandbox.CLI(), "--settings", h.settingsPath, "--"}
	if trampoline {
		line = append(append(line, "/bin/sh", "-c", execTrampoline, "world-exec"), h.childEnv()...)
	}
	return append(line, argv...)
}

// Spawns is how many subprocesses this handler has started (cache seeding
// and srt): a refusal, a denial or a drifted stack leaves it unchanged.
func (h *ExecHandler) Spawns() int64 { return h.spawns.Load() }

var execPayloadKeys = []string{"command", "args"}

// parseExecPayload is the handler's defence-in-depth twin of exec.ail's plan
// checks (§4.1): exact keys, the command id grammar and the args bounds.
func parseExecPayload(payload []byte) (string, []string, string) {
	fields, err := decodePayloadObject(payload)
	if err != nil {
		return "", nil, err.Error()
	}
	var id string
	var args []string
	for key, value := range fields {
		switch key {
		case "command":
			if json.Unmarshal(value, &id) != nil || !execCommandIDPattern.MatchString(id) {
				return "", nil, fmt.Sprintf("command must be a string matching %s", execCommandIDPattern)
			}
		case "args":
			if json.Unmarshal(value, &args) != nil || args == nil {
				return "", nil, "args must be an array of strings"
			}
			if len(args) > execMaxArgs {
				return "", nil, fmt.Sprintf("args has %d items; the limit is %d", len(args), execMaxArgs)
			}
			for _, a := range args {
				if len(a) < 1 || len(a) > 512 || strings.ContainsRune(a, 0) {
					return "", nil, "every arg must be 1–512 bytes with no NUL"
				}
			}
		default:
			for _, name := range execPayloadKeys {
				if strings.EqualFold(key, name) {
					return "", nil, fmt.Sprintf("key %q aliases %q", key, name)
				}
			}
			return "", nil, fmt.Sprintf("unknown payload key %q", key)
		}
	}
	if id == "" {
		return "", nil, "command is required"
	}
	return id, args, ""
}

// projectDir is <worktree>/<root>: an existing directory reached without a
// symlink, so a root cannot lead the cwd out of the worktree.
func (h *ExecHandler) projectDir() (string, error) {
	want := filepath.Join(h.worktree, h.profile.Root)
	got, err := filepath.EvalSymlinks(want)
	if err != nil || got != want {
		return "", fmt.Errorf("broker: exec handler: project root %q is not a directory inside the worktree (%v)", want, err)
	}
	if info, err := os.Stat(got); err != nil || !info.IsDir() {
		return "", fmt.Errorf("broker: exec handler: project root %q is not a directory", want)
	}
	return got, nil
}

// childEnv is the K=V set the trampoline's `env -i` installs (§4.5), built
// fresh for each call: World's set, then the profile env and the caches
// variables, each in byte order.
func (h *ExecHandler) childEnv() []string {
	kv := []string{
		"HOME=" + filepath.Join(h.cache, "home"),
		"TMPDIR=" + filepath.Join(h.cache, "tmp"),
		"PATH=" + strings.Join(h.profile.Path, ":"),
		"LANG=C.UTF-8",
		"LC_ALL=C.UTF-8",
		"PYTHONPYCACHEPREFIX=" + filepath.Join(h.cache, "pycache"),
		"npm_config_cache=" + filepath.Join(h.cache, "npm"),
		"CI=1",
		"NO_COLOR=1",
	}
	for _, name := range sortedKeys(h.profile.Env) {
		kv = append(kv, name+"="+h.profile.Env[name])
	}
	for _, name := range sortedKeys(h.profile.Caches) {
		kv = append(kv, name+"="+filepath.Join(h.cache, name))
	}
	return kv
}

// prepareCache makes the episode's exec cache (D-140-3 = A): World's own
// dirs, then each caches entry, seeded on first use by a clonefile copy of
// the operator's read-only seed (V21). The seed is never written.
func (h *ExecHandler) prepareCache(ctx context.Context) error {
	h.seedMu.Lock()
	defer h.seedMu.Unlock()
	for _, d := range []string{"home", "tmp", "pycache", "npm"} {
		if err := os.MkdirAll(filepath.Join(h.cache, d), 0o700); err != nil {
			return err
		}
	}
	for _, name := range sortedKeys(h.profile.Caches) {
		dst := filepath.Join(h.cache, name)
		if _, err := os.Lstat(dst); err == nil {
			continue
		}
		seed := h.profile.Caches[name].Seed
		if seed == "" {
			if err := os.MkdirAll(dst, 0o700); err != nil {
				return err
			}
			continue
		}
		args := []string{"-cR", seed, dst} // APFS clonefile
		if hostGOOS() != "darwin" {
			args = []string{"-a", "--reflink=auto", seed, dst}
		}
		h.spawns.Add(1)
		if _, err := runBounded(ctx, handlerBounds{execTimeout: execSeedBudget(ctx)}, handlerCommand{
			path: "/bin/cp", args: args, dir: h.cache, env: execHostEnv(),
		}); err != nil {
			_ = os.RemoveAll(dst)
			return fmt.Errorf("broker: exec handler: seed cache %s from %s: %w", name, seed, err)
		}
	}
	return nil
}

// execSeedBudget bounds a seed copy by the call's remaining budget.
func execSeedBudget(ctx context.Context) time.Duration {
	if dl, ok := ctx.Deadline(); ok {
		if r := time.Until(dl); r > 0 {
			return r
		}
		return time.Millisecond
	}
	return handlerExecTimeout
}
