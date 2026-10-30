package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

const episodeLaunchHelp = `usage: ailang-worldd episode launch <episode> --preset <name> [--presets-dir D] [--hosts-dir D] [--headless]
  --preset NAME    required; daemon and paths; must match the pinned start
  --presets-dir D  private preset directory
  --hosts-dir D    hosts init output directory
  --headless       run the pinned task with a bounded wall cap
Interactive launch is the default and uses this terminal and hosts init launcher.
Headless writes the stream and launch history
beside the session, then prints the finish step. The session remains live.
Example: ailang-worldd episode launch my-episode --preset my-repo --headless
`

var launchClient = newClient

func launchMode(p launchPreset, headless bool) (string, time.Duration, error) {
	mode := "interactive"
	deadline := 3600.0
	if p.Launch != nil {
		if p.Launch.Mode != "" {
			mode = p.Launch.Mode
		}
		if p.Launch.DeadlineS > 0 {
			deadline = p.Launch.DeadlineS
		}
	}
	if headless {
		mode = "headless"
	}
	if mode != "headless" {
		return mode, 0, nil
	}
	if deadline > float64((time.Duration(1<<63-1)-time.Minute)/time.Second) {
		return mode, 0, errors.New("launch deadline is too large")
	}
	wall := time.Duration(deadline * float64(time.Second))
	if wall <= 0 {
		return mode, 0, errors.New("launch deadline is below timer resolution")
	}
	return mode, wall, nil
}

// canonicalLaunchPath resolves existing ancestors too, so an absent child under
// a symlink cannot evade the workspace/sessions overlap check.
func canonicalLaunchPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(path)
	if err == nil {
		return real, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	parent := filepath.Dir(path)
	if parent == path {
		return "", err
	}
	real, err = canonicalLaunchPath(parent)
	return filepath.Join(real, filepath.Base(path)), err
}

func launchInside(a, b string) bool {
	rel, err := filepath.Rel(b, a)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func validateLaunchDir(dir, root, wt, sessions string) error {
	fail := func(what string) error {
		return presetRefusal(what+" (F8)", "the agent must start in an empty directory outside World data", "use an empty launch directory outside the workspace and sessions; the session is untouched")
	}
	real, err := canonicalLaunchPath(dir)
	if err != nil {
		return fail(err.Error())
	}
	for _, other := range []string{root, wt, sessions} {
		resolved, err := canonicalLaunchPath(other)
		if err != nil {
			return fail(err.Error())
		}
		if launchInside(real, resolved) || launchInside(resolved, real) {
			return fail(dir + " overlaps " + other)
		}
	}
	info, err := os.Lstat(dir)
	if errors.Is(err, os.ErrNotExist) {
		// Check the nearest existing ancestor is a usable directory, without
		// creating anything before validation and the mint fence.
		parent := filepath.Dir(dir)
		for {
			info, err = os.Stat(parent)
			if !errors.Is(err, os.ErrNotExist) || filepath.Dir(parent) == parent {
				break
			}
			parent = filepath.Dir(parent)
		}
		if err != nil || !info.IsDir() || info.Mode().Perm()&0o300 != 0o300 {
			return fail(dir + " has no usable parent directory")
		}
		return nil
	}
	if err != nil {
		return fail(err.Error())
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fail(dir + " is a symlink or not a directory")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fail(err.Error())
	}
	if len(entries) > 0 {
		return fail(dir + " is not empty")
	}
	if info.Mode().Perm()&0o300 != 0o300 {
		return fail(dir + " is not usable")
	}
	return nil
}

func launchBinary(agent string) (string, error) {
	if agent != "claude" && agent != "codex" {
		return "", presetRefusal("unknown agent "+agent, "headless supports claude and codex", "choose claude or codex in a new preset")
	}
	path, err := exec.LookPath(agent)
	if err == nil {
		path, err = filepath.Abs(path)
	}
	if err == nil {
		path, err = filepath.EvalSymlinks(path)
	}
	if err != nil {
		return "", presetRefusal(agent+" is missing from PATH (F10)", "launch requires the agent CLI", "install "+agent+" and add it to PATH")
	}
	return path, nil
}

func launchLockCheck(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		if _, e := os.Lstat(path); errors.Is(e, os.ErrNotExist) {
			return nil
		}
	}
	if err != nil {
		return presetRefusal("launch lock exists or cannot be read: "+path, "another launch may own it", "remove the lock if that launch has ended")
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	running := false
	if pid > 0 {
		e := syscall.Kill(pid, 0)
		running = e == nil || errors.Is(e, syscall.EPERM)
	}
	return presetRefusal(fmt.Sprintf("launch lock %s holds pid %d (running: %t)", path, pid, running), "only one launch may run per episode", "remove the lock if that launch has ended")
}

func launchRegularOrAbsent(path string) error {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	return nil
}

func launchCommands(h daemon.HealthResponse, profile string) ([]string, error) {
	if profile == "" {
		return nil, nil
	}
	for _, p := range h.ExecProfiles {
		if p.Project != profile {
			continue
		}
		data, err := os.ReadFile(p.File)
		if err != nil {
			return nil, err
		}
		parsed, err := broker.ParseExecProfile(data)
		if err != nil {
			return nil, err
		}
		if parsed.Digest != p.Digest || parsed.Project != profile {
			return nil, errors.New("served exec profile digest or project changed")
		}
		commands := make([]string, 0, len(parsed.Commands))
		for id := range parsed.Commands {
			commands = append(commands, id)
		}
		sort.Strings(commands)
		return commands, nil
	}
	return nil, fmt.Errorf("daemon does not serve start profile %s", profile)
}

func runEpisodeLaunch(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd episode launch", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	name := fs.String("preset", "", "required launch preset")
	dir := fs.String("presets-dir", defaultPresetsDir(), "private presets")
	headless := fs.Bool("headless", false, "headless launch")
	hostsDir := fs.String("hosts-dir", defaultHostsDir(), "hosts init output directory")
	pos, err := parseInterleaved(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(stdout, episodeLaunchHelp)
		return exitOK
	}
	fail := func(err error) int {
		fmt.Fprintln(stderr, presetRefusal(err.Error(), "launch needs a valid pinned episode", "ailang-worldd episode launch --help"))
		return exitUsage
	}
	if err != nil {
		return fail(err)
	}
	if len(pos) != 1 || !workspaceEpisodePattern.MatchString(pos[0]) || *name == "" || *name == "se-tools" {
		return fail(errors.New("episode and --preset <name> are required"))
	}
	return launchEpisode(pos[0], *name, *dir, *headless, stdout, stderr, *hostsDir)
}

// launchEpisode is shared by episode launch and the post-mint start path.
// Its preflight has no writes; only a fully validated launch takes the lock.
func launchEpisode(ep, name, presetsDir string, headless bool, stdout, stderr io.Writer, hostsDirs ...string) int {
	fail := func(err error) int {
		msg := err.Error()
		if !strings.Contains(msg, "\n  why:") {
			msg = presetRefusal(msg, "launch validation failed before starting the agent", "ailang-worldd episode launch "+ep+" --preset "+name+" --headless").Error()
		}
		fmt.Fprintln(stderr, msg)
		return exitUsage
	}
	p, err := finishPreset(name, presetsDir)
	if err != nil {
		return fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	h, err := fetchDaemonHealth(p.Addr)
	if err != nil {
		return fail(err)
	}
	if same, _ := store.SameDatabase(h.DBPath, p.DB); !same {
		return fail(presetRefusal(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, name, p.DB), "a preset must bind the served store", "fix the preset daemon, or run service status"))
	}
	c := launchClient(p.Addr)
	_, rep, err := fetchReport(ctx, c, ep)
	if err != nil && !errors.Is(err, errNoCalls) {
		return fail(err)
	}
	if rep.Start == nil {
		return fail(presetRefusal("episode "+ep+" has no start record", "launch reads the task and the launch settings pinned at start", "ailang-worldd episode start --preset "+name))
	}
	start := rep.Start
	if start.Preset.Name != name {
		return fail(presetRefusal("preset "+name+" does not match start preset "+start.Preset.Name, "launch uses the preset approved at start", "use --preset "+start.Preset.Name))
	}
	var pinned launchPreset
	if err := json.Unmarshal([]byte(start.Preset.Bytes), &pinned); err != nil {
		return fail(fmt.Errorf("pinned preset bytes: %w", err))
	}
	hostsDir := defaultHostsDir()
	if len(hostsDirs) > 0 {
		hostsDir = hostsDirs[0]
	}
	mode, wall, err := launchMode(pinned, headless)
	if err != nil {
		return fail(err)
	}
	var tty *os.File
	if mode == "interactive" {
		tty, err = openLaunchTerminal()
		if err != nil {
			return fail(presetRefusal("interactive launch needs a terminal ("+err.Error()+")", "the agent runs in this terminal", "run it from a terminal window, or add --headless"))
		}
		defer tty.Close()
	}
	var task []byte
	if start.Task == nil {
		if mode == "headless" {
			return fail(errors.New("start record has no task"))
		}
	} else {
		taskPath := filepath.Join(p.SessionsDir, ep+".task.md")
		if err := launchRegularOrAbsent(taskPath); err != nil {
			return fail(err)
		}
		task, err = os.ReadFile(taskPath)
		if err != nil {
			return fail(fmt.Errorf("pinned task copy %s is missing or unreadable: %w", taskPath, err))
		}
		if serviceHash(task) != start.Task.SHA256 {
			return fail(fmt.Errorf("pinned task copy %s digest does not match start task.sha256", taskPath))
		}
	}
	newRound := "ailang-worldd episode start --preset " + name + " --branch " + start.Branch
	token, err := resolveSession(filepath.Join(p.SessionsDir, ep+".session"), io.Discard)
	expired := func() int {
		if path, e := credentialsFor(p.DB, p.CredentialsDB); e == nil {
			if ro, e := store.OpenCredentialsReadOnly(path); e == nil {
				rows, e := ro.ListSessions(ctx, ep, store.MaxListSessions)
				_ = ro.Close()
				live := false
				for _, row := range rows {
					live = live || row.ExpiresAt >= time.Now().Unix()
				}
				if e == nil && len(rows) > 0 && !live {
					newRound = "ailang-worldd episode start --resume " + ep + " --preset " + name
				}
			}
		}
		return fail(presetRefusal("episode "+ep+"'s session has expired or been revoked (F9)", "launch requires a live session for this episode", newRound))
	}
	if err != nil {
		return expired()
	}
	status, _, body, err := c.doRequest(ctx, http.MethodGet, "/v1/session/self", nil, map[string]string{"Authorization": "Bearer " + token})
	var self daemon.SessionSelf
	if err != nil || status != http.StatusOK || json.Unmarshal(body, &self) != nil || self.Episode != ep || self.SecondsLeft <= 0 {
		return expired()
	}
	if mode == "interactive" {
		wall = time.Duration(self.SecondsLeft) * time.Second
	}
	contextText, err := launchContext(body, p.Addr, token)
	if err != nil {
		return fail(err)
	}
	commands, err := launchCommands(h, pinned.Profile)
	if err != nil {
		return fail(err)
	}
	prompt := ""
	if start.Task != nil {
		prompt = launchPrompt(string(task), commands)
	}
	if strings.Contains(prompt, token) {
		return fail(errors.New("pinned task would put the token in argv"))
	}
	host, err := launchBinary(start.Agent)
	if err != nil {
		return fail(err)
	}
	dir := filepath.Join(filepath.Dir(p.SessionsDir), "launch", ep)
	if err := validateLaunchDir(dir, p.WorkspaceRoot, filepath.Join(p.WorkspaceRoot, ep), p.SessionsDir); err != nil {
		return fail(err)
	}
	lock := filepath.Join(p.SessionsDir, ep+".launch.lock")
	if err := launchLockCheck(lock); err != nil {
		return fail(err)
	}
	if mode == "interactive" {
		host, err = launchLauncher(start.Agent, hostsDir, p.Addr)
		if err != nil {
			return fail(err)
		}
	}
	transcript := filepath.Join(p.SessionsDir, ep+".transcript.jsonl")
	errPath := filepath.Join(p.SessionsDir, ep+".agent.err")
	mcp := filepath.Join(p.SessionsDir, ep+".mcp.json")
	for _, path := range []string{transcript, errPath, mcp} {
		if err := launchRegularOrAbsent(path); err != nil {
			return fail(err)
		}
	}
	var grants []broker.Capability
	for _, g := range self.Grants {
		grants = append(grants, broker.Capability{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
	}
	mcpArg := ""
	if start.Agent == "claude" {
		mcpArg = mcp
	}
	argv := launchArgv(start.Agent, start.Model, prompt, dir, mcpArg, p.Addr, contextText, "", grants)
	sessionID := ""
	if mode == "interactive" {
		if start.Agent == "claude" {
			sessionID, err = launchSessionUUID()
			if err != nil {
				return fail(err)
			}
		}
		argv = interactiveLaunchArgv(host, start.Agent, start.Model, sessionID, contextText, prompt, grants)
	}
	for _, a := range append(append([]string{}, argv...), host) {
		if strings.Contains(a, token) {
			return fail(errors.New("launch argv would carry the session token"))
		}
	}
	credsPath, err := credentialsFor(p.DB, p.CredentialsDB)
	if err != nil {
		return fail(err)
	}
	// Check local history before starting; persistence uses the writer lock
	// only after exit, so the daemon can continue reading while the agent runs.
	ro, err := store.OpenCredentialsReadOnly(credsPath)
	if err != nil {
		return fail(err)
	}
	localStart, readErr := ro.EpisodeStart(ctx, ep)
	_ = ro.Close()
	if readErr != nil {
		return fail(readErr)
	}
	if localStart == nil {
		return fail(errors.New("local credentials store has no start record"))
	}
	f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fail(err)
	}
	defer os.Remove(lock)
	_, err = fmt.Fprintf(f, "%d\n", os.Getpid())
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		return fail(errors.Join(err, closeErr))
	}
	// A preceding launch may have finished between preflight and O_EXCL.
	// Reread its history under our episode lock before choosing the separator.
	ro, err = store.OpenCredentialsReadOnly(credsPath)
	if err != nil {
		return fail(err)
	}
	localStart, readErr = ro.EpisodeStart(ctx, ep)
	_ = ro.Close()
	if readErr != nil || localStart == nil {
		return fail(errors.Join(readErr, errors.New("cannot read launch history under the launch lock")))
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fail(err)
	}
	if err := os.Chmod(filepath.Dir(dir), 0o700); err != nil {
		return fail(err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return fail(err)
	}
	// Recheck under the launch lock immediately before writing credentials.
	if err := validateLaunchDir(dir, p.WorkspaceRoot, filepath.Join(p.WorkspaceRoot, ep), p.SessionsDir); err != nil {
		return fail(err)
	}
	if mode == "headless" && start.Agent == "claude" {
		if err := os.Remove(mcp); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fail(err)
		}
		if err := writeMCPConfig(mcp, p.Addr, token); err != nil {
			return fail(err)
		}
		defer os.Remove(mcp)
	} else {
		mcp = ""
	}
	outStage, err := os.CreateTemp(p.SessionsDir, "."+ep+".stdout-*")
	if err != nil {
		return fail(err)
	}
	defer os.Remove(outStage.Name())
	defer outStage.Close()
	errStage, err := os.CreateTemp(p.SessionsDir, "."+ep+".stderr-*")
	if err != nil {
		return fail(err)
	}
	defer os.Remove(errStage.Name())
	defer errStage.Close()
	childEnv := []string{}
	for _, kv := range childenv.Scrubbed(os.Environ()) {
		if !strings.HasPrefix(kv, "WORLD_SESSION=") && !strings.Contains(kv, token) {
			childEnv = append(childEnv, kv)
		}
	}
	if mode == "interactive" || start.Agent == "codex" {
		childEnv = append(childEnv, "WORLD_SESSION="+token)
	}
	redacted := redactLaunchArgv(argv, prompt, token)
	argv[0] = host
	runCtx, runCancel := context.WithTimeout(context.Background(), wall+time.Minute)
	defer runCancel()
	startTime := time.Now()
	started := startTime.UTC().Format(time.RFC3339)
	var result agentrun.Result
	var runErr error
	if mode == "interactive" {
		result, runErr = agentrun.RunInteractive(runCtx, argv, dir, childEnv, tty, wall)
	} else {
		result, runErr = agentrun.Run(runCtx, argv, dir, childEnv, outStage, errStage, wall)
	}
	endTime := time.Now()
	ended := endTime.UTC().Format(time.RFC3339)
	var deleted *bool
	var cleanupErr error
	if mcp != "" {
		e := os.Remove(mcp)
		ok := e == nil || errors.Is(e, os.ErrNotExist)
		deleted = &ok
		if !ok {
			cleanupErr = fmt.Errorf("delete private MCP config: %w", e)
		}
	}
	n := len(localStart.Launches) + 1
	var state, digest string
	var keepErr, stderrErr error
	var source hostTranscript
	if mode == "interactive" {
		if start.Agent == "claude" {
			source = claudeTranscriptPath(dir, sessionID, childEnv)
		} else {
			source = codexTranscriptPath(hostsDir, dir, startTime, endTime)
			sessionID = source.id
		}
		state = "missing"
		if source.path != "" {
			if e := copyHostTranscript(source, outStage); e != nil {
				source.reason = e.Error()
			} else {
				state, digest, keepErr = keepLaunchStream(outStage, transcript, token, n)
			}
		}
	} else {
		state, digest, keepErr = keepLaunchStream(outStage, transcript, token, n)
		_, _, stderrErr = keepLaunchStream(errStage, errPath, token, n)
	}
	entry := store.EpisodeLaunch{Mode: mode, StartedAt: started, EndedAt: ended, Exit: result.ExitCode, FinishReason: result.Reason,
		Transcript: transcript, TranscriptSHA256: digest, TranscriptState: state, Host: host, ArgvRedacted: redacted, MCPConfigDeleted: deleted, HostSessionID: sessionID, HostTranscript: source.path, TranscriptReason: source.reason}
	storeCtx, storeCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer storeCancel()
	creds, recordErr := store.OpenCredentials(credsPath)
	if recordErr == nil {
		recordErr = creds.AppendEpisodeLaunch(storeCtx, ep, entry)
		recordErr = errors.Join(recordErr, creds.Close())
	}
	lockErr := os.Remove(lock)
	fmt.Fprintf(stdout, "next: ailang-worldd episode finish %s --preset %s\nailang-worldd episode report %s\n", ep, name, ep)
	if err := errors.Join(runErr, keepErr, stderrErr, recordErr, cleanupErr, lockErr); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if result.ExitCode != 0 {
		fmt.Fprintf(stderr, "agent exit %d (%s)\n", result.ExitCode, result.Reason)
		return 1
	}
	return exitOK
}

// containsLaunchToken uses a sliding window; agent output never accumulates in
// memory. Digests in transcripts are legitimate, only the literal token is secret.
func containsLaunchToken(r io.Reader, token string) (bool, error) {
	if token == "" {
		return false, errors.New("empty session token")
	}
	buf := make([]byte, 32*1024)
	tail := []byte{}
	for {
		n, err := r.Read(buf)
		window := append(tail, buf[:n]...)
		if bytes.Contains(window, []byte(token)) {
			return true, nil
		}
		keep := min(len(token)-1, len(window))
		tail = append([]byte(nil), window[len(window)-keep:]...)
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
}

func keepLaunchStream(stage *os.File, path, token string, n int) (string, string, error) {
	if _, err := stage.Seek(0, io.SeekStart); err != nil {
		return "", "", err
	}
	hit, err := containsLaunchToken(stage, token)
	if err != nil {
		return "", "", err
	}
	if err := launchRegularOrAbsent(path); err != nil {
		return "", "", err
	}
	if previous, err := os.Open(path); err == nil {
		oldHit, scanErr := containsLaunchToken(previous, token)
		_ = previous.Close()
		if scanErr != nil {
			return "", "", scanErr
		}
		hit = hit || oldHit
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if hit || n == 1 {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o600)
	if err != nil {
		return "", "", err
	}
	if err := f.Chmod(0o600); err != nil {
		_ = f.Close()
		return "", "", err
	}
	state := "kept"
	if hit {
		state = "scrubbed"
		_, err = io.WriteString(f, "{\"world_refusal\":\"session token in agent stream; scrubbed\"}\n")
	} else {
		if n > 1 {
			// Ensure the separator is on its own line even after a partial line.
			if old, e := os.Open(path); e == nil {
				if size, e := old.Stat(); e == nil && size.Size() > 0 {
					last := []byte{0}
					_, _ = old.ReadAt(last, size.Size()-1)
					if last[0] != '\n' {
						_, err = io.WriteString(f, "\n")
					}
				}
				_ = old.Close()
			}
			if err == nil {
				_, err = fmt.Fprintf(f, "{\"world_launch\": %d}\n", n)
			}
		}
		if err == nil {
			_, err = stage.Seek(0, io.SeekStart)
		}
		if err == nil {
			_, err = io.Copy(f, stage)
		}
	}
	err = errors.Join(err, f.Close())
	if err != nil {
		return state, "", err
	}
	kept, err := os.Open(path)
	if err != nil {
		return state, "", err
	}
	defer kept.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, kept); err != nil {
		return state, "", err
	}
	return state, "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}
