package main

// Finish is an attended local operation. It reads authority (checks, branch,
// base) from the immutable start record, never from mutable preset checks.
import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const episodeFinishHelp = `usage: ailang-worldd episode finish <episode> --preset <name> [--presets-dir D]
    [--ready] [--exec-profiles-dir D] [--title <text>] [--allow-failing <check>]... [--include-unattributed] [--dry-run]

ATTENDED: requires a controlling terminal before reading files or contacting
the daemon. Runs the checks pinned at episode start, records each outcome,
attributes changed paths and makes a local git commit with World provenance.
Each unsandboxed check needs its own typed y at that terminal.
  --preset NAME           required; daemon and paths; must match start
  --presets-dir D         private preset directory
  --title TEXT            commit title (default: start task title)
  --allow-failing CHECK   allow this named failing check; repeatable
  --include-unattributed  include hand edits, listed as Not made through World
  --ready                 open a ready PR (default: draft)
  --exec-profiles-dir D    pinned profile lookup (default: private exec-profiles)
  --dry-run               print git show --stat and body, then reset HEAD

Without --dry-run, writes <sessions>/<episode>.pr-body.md (0600) and prints
an exact-byte approval at the terminal. Typed y records approval, pushes and
opens the PR, records its outcome, then revokes the episode credentials.
Example: ailang-worldd episode finish my-episode --preset my-repo --dry-run
`

var finishClient = newClient

func runEpisodeFinish(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	return runEpisodeFinishWith(args, stdout, stderr, env, commitOneObject)
}

func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEnv, record outcomeCommitter) (code int) {
	fs := flag.NewFlagSet("ailang-worldd episode finish", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	preset := fs.String("preset", "", "required launch preset")
	dir := fs.String("presets-dir", defaultPresetsDir(), "private presets")
	title := fs.String("title", "", "commit title")
	include := fs.Bool("include-unattributed", false, "include hand edits")
	ready := fs.Bool("ready", false, "open a ready PR")
	profilesDir := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "private exec profiles")
	dry := fs.Bool("dry-run", false, "print then reset local commit")
	var allowed multiFlag
	fs.Var(&allowed, "allow-failing", "allow this failing check")
	pos, err := parseInterleaved(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(stdout, episodeFinishHelp)
		return exitOK
	}
	if err != nil {
		return exitUsage
	}
	fail := func(what, why, next string) int {
		fmt.Fprintln(stderr, presetRefusal(what, why, next))
		return exitUsage
	}
	// F1 precedes even bad preset/episode validation, and all IO beyond flags.
	tty, err := env.openTerminal()
	if err != nil {
		return fail("refusing: no controlling terminal", "finish is an attended act", "run it in a terminal (an agent or CI can't finish an episode)")
	}
	defer tty.Close()
	if len(pos) != 1 || !workspaceEpisodePattern.MatchString(pos[0]) || *preset == "" {
		return fail("episode and --preset <name> are required", "finish needs a pinned start record", "ailang-worldd episode finish --help")
	}
	ep := pos[0]
	// One bound for the whole command: each check also has its own cap, and the
	// prompt waits on a person, so this is generous rather than tight.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	p, err := finishPreset(*preset, *dir)
	if err != nil {
		return fail(err.Error(), "finish needs the preset's daemon and paths", "ailang-worldd preset check "+*preset)
	}
	c := finishClient(p.Addr)
	h, err := fetchDaemonHealth(p.Addr)
	if err != nil {
		return fail(err.Error(), "finish needs a serving daemon", "ailang-worldd service status")
	}
	if same, _ := store.SameDatabase(h.DBPath, p.DB); !same {
		return fail(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, *preset, p.DB), "a preset must bind the served store", "fix the preset daemon, or run service status")
	}
	credentials, err := credentialsFor(p.DB, p.CredentialsDB)
	if err != nil {
		return fail(err.Error(), "finish needs the preset credentials store", "fix the preset")
	}
	if same, _ := store.SameDatabase(h.CredentialsDB, credentials); !same {
		return fail(fmt.Sprintf("the daemon serves credentials store %s; preset %s names %s (F6)", h.CredentialsDB, *preset, credentials), "the operator proof must be beside the served credentials store", "fix the preset daemon")
	}
	_, rep, err := fetchReport(ctx, c, ep)
	if err != nil && !errors.Is(err, errNoCalls) {
		return fail(err.Error(), "finish needs the episode report", "ailang-worldd episode report "+ep)
	}
	if rep.Start == nil {
		return fail("episode "+ep+" has no start record", "finish runs the checks pinned at start", "start the next episode with ailang-worldd episode start --preset <name>")
	}
	start := rep.Start
	if start.Preset.Name != *preset {
		return fail(fmt.Sprintf("preset %s does not match start preset %s", *preset, start.Preset.Name), "finish must use the start preset", "use --preset "+start.Preset.Name)
	}
	newRound := "ailang-worldd episode start --preset " + *preset + " --branch " + start.Branch
	token, err := resolveSession(filepath.Join(p.SessionsDir, ep+".session"), stderr)
	if err != nil {
		return fail("episode "+ep+"'s session has expired: finish records under it", "a live episode session is required", newRound)
	}
	status, _, body, err := c.doRequest(ctx, http.MethodGet, "/v1/session/self", nil, map[string]string{"Authorization": "Bearer " + token})
	var self daemon.SessionSelf
	if err != nil || status != http.StatusOK || json.Unmarshal(body, &self) != nil || self.Episode != ep || self.SecondsLeft <= 0 {
		return fail("episode "+ep+"'s session has expired: finish records under it", "a live episode session is required", newRound)
	}
	wt := filepath.Join(p.WorkspaceRoot, ep)
	degraded := !worktree.IsLinked(wt)
	for _, e := range h.DegradedEpisodes {
		degraded = degraded || e == ep
	}
	if degraded {
		return fail("episode "+ep+" is degraded", "finish requires a linked git worktree", newRound)
	}
	if rep.Uncommitted != 0 {
		return fail("an effect of "+ep+" is still in flight", "uncommitted_intents must be zero", "wait, or ailang-worldd episode report "+ep+" --follow")
	}
	git, err := sessionGit()
	if err != nil {
		return fail(err.Error(), "finish requires git", "install git")
	}
	query := func(args ...string) (string, error) {
		return worktree.Output(ctx, git, append([]string{"-C", wt, "-c", "core.fsmonitor=false"}, args...))
	}
	branch, err := query("symbolic-ref", "--short", "HEAD")
	if err != nil || strings.TrimSpace(branch) != start.Branch {
		return fail("HEAD is not on "+start.Branch, "finish uses the start record's branch", "restore the recorded branch")
	}
	previous, err := query("rev-parse", "HEAD")
	if err != nil {
		return fail(err.Error(), "cannot read HEAD", "repair the worktree")
	}
	previous = strings.TrimSpace(previous)
	bodyFile := filepath.Join(p.SessionsDir, ep+".pr-body.md")
	steps := func(sha string) string {
		return finishHandSteps(wt, start, sha, bodyFile)
	}
	message, err := query("log", "-1", "--format=%B")
	if err != nil {
		return fail(err.Error(), "cannot read HEAD message", "repair the worktree")
	}
	for _, line := range strings.Split(message, "\n") {
		if strings.HasPrefix(line, "World-Provenance: ") && strings.Contains(line, " episode="+ep+" ") {
			return fail("already committed at "+previous, "re-runs are row 218 M3", steps(previous))
		}
	}
	changed, err := finishChanged(query, start.Base.SHA)
	if err != nil {
		return fail(err.Error(), "cannot inspect changed paths", "repair the worktree")
	}
	if len(changed) == 0 {
		return fail("nothing to finish", "the worktree has no changes against "+start.Base.SHA, newRound)
	}
	if err := validateFinishChecks(start.Checks, allowed, wt); err != nil {
		return fail(err.Error(), "all checks must be valid before any runs", "fix checks in the next episode's preset")
	}
	reader := bufio.NewReader(tty)
	var results []finishCheckResult
	for _, check := range start.Checks {
		var r checkrun.Result
		var runErr error
		if check.Unsandboxed {
			cwd := filepath.Join(wt, check.Cwd)
			fmt.Fprintf(tty, "check %s runs the agent's code outside the sandbox: %s in %s\nRun this check? [y/N] ", check.Name, strings.Join(check.Argv, " "), cwd)
			answer, _ := reader.ReadString('\n')
			if strings.TrimSpace(answer) != "y" {
				return fail("check "+check.Name+" declined", "unsandboxed code requires its own typed y", "run finish again when ready")
			}
			r, runErr = checkrun.Run(ctx, check.Argv, cwd, time.Duration(check.TimeoutS*float64(time.Second)))
		} else {
			r, runErr = finishSandboxCheck(ctx, c, token, check)
		}
		if runErr != nil {
			r.ExitCode = -1
			r.Stderr = runErr.Error() + "\n" + r.Stderr
		}
		passed := r.ExitCode == 0 && runErr == nil
		payload := finishOutcome(rep, check, r, passed)
		raw, _ := json.Marshal(payload)
		canonical, _, err := parseTaskOutcome(raw)
		if err != nil {
			return fail(err.Error(), "check outcome cannot be recorded", "repair the check definition")
		}
		entry, err := record(ctx, c, token, taskOutcomeObject(canonical), taskOutcomeWrittenBy, "", "outcome record")
		if err != nil {
			return fail(err.Error(), "the check verdict was not confirmed recorded; do not retry blindly", "inspect the episode report")
		}
		results = append(results, finishCheckResult{check.Name, check.Unsandboxed, passed, entry, taskOutcomeObject(canonical).Hash})
		if !passed && !containsString(allowed, check.Name) {
			return fail(fmt.Sprintf("check %s failed (rc %d), outcome %d\n%s", check.Name, r.ExitCode, entry, clipOutcomeHead(r.Stderr)), "a required check failed", "fix it in a new round, or --allow-failing "+check.Name)
		}
	}
	// Refresh both evidence and paths after checks, which may produce files.
	_, rep, err = fetchReport(ctx, c, ep)
	if err != nil {
		return fail(err.Error(), "cannot refresh check evidence", "inspect the episode report")
	}
	if rep.Uncommitted != 0 {
		return fail("an effect of "+ep+" is still in flight", "uncommitted_intents must be zero", "wait for it")
	}
	changed, err = finishChanged(query, start.Base.SHA)
	if err != nil {
		return fail(err.Error(), "cannot refresh changed paths", "repair the worktree")
	}
	unattributed := finishUnattributed(changed, rep)
	if len(unattributed) > 0 && !*include {
		gaps, _ := json.Marshal(struct {
			Truncated []int64                   `json:"deltas_truncated"`
			Errors    []daemon.ReportDeltaError `json:"delta_errors"`
		}{rep.Writes.DeltasTruncated, rep.Writes.DeltaErrors})
		return fail("unattributed paths (F5): "+strings.Join(unattributed, ", "), "no World write or delta proves these paths; "+string(gaps), "remove them, or --include-unattributed")
	}
	actualTitle := *title
	if actualTitle == "" && start.Task != nil {
		actualTitle = start.Task.Title
	}
	if actualTitle == "" {
		actualTitle = "World episode " + ep
	}
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}
	protected, err := finishProtected(ctx, start, credentials, *profilesDir, changed)
	if err != nil {
		return fail(err.Error(), "finish must use the start profile binding", "restore the pinned exec profile")
	}
	paragraph := fmt.Sprintf("%s %s; %d calls; %d/%d checks passed; %d files changed.", start.Agent, start.Model, rep.TotalCalls, passed, len(results), len(changed))
	if len(protected) > 0 {
		paragraph += " Protected paths changed: " + strings.Join(protected, ", ") + "."
	}
	trailer := reportTrailer(rep, h.DBPath)
	add := reportAdditions{Checks: results, Unattributed: unattributed, Protected: protected, Finish: true}
	// Refuse unsafe body bytes before creating the local commit.
	if _, err := renderReportMarkdown(rep, h.DBPath, add); err != nil {
		return fail(err.Error(), "the body must pass the trial scrub", "inspect the episode report")
	}
	for _, part := range []string{actualTitle, paragraph, trailer} {
		if err := scrubMarkdown(part); err != nil {
			return fail(err.Error(), "the commit must pass the trial scrub", "choose safe commit text")
		}
	}
	sha, err := worktree.Commit(ctx, git, wt, start.Agent, changed, actualTitle, paragraph, trailer)
	if err != nil {
		return fail(err.Error(), "local git commit failed", "inspect git status before retrying")
	}
	add.Commit = sha
	rendered, renderErr := renderReportMarkdown(rep, h.DBPath, add)
	if *dry || renderErr != nil {
		// Always restore HEAD on a dry-run error too.
		defer func() {
			if _, err := query("reset", "-q", "--mixed", previous); err != nil {
				fmt.Fprintln(stderr, "dry-run reset failed: "+err.Error())
				code = exitUsage
			}
		}()
	}
	if renderErr != nil {
		return fail(renderErr.Error(), "body render failed", "inspect the episode report")
	}
	if *dry {
		stat, err := query("show", "--stat", "--format=medium", "HEAD")
		if err != nil {
			return fail(err.Error(), "cannot show dry-run commit", "inspect the worktree")
		}
		fmt.Fprintln(stdout, stat)
		fmt.Fprint(stdout, rendered)
		return exitOK
	}
	if err := writeFinishBody(bodyFile, rendered); err != nil {
		return fail(err.Error(), "commit kept at "+sha, steps(sha))
	}
	fmt.Fprintf(stdout, "committed at %s\nbody: %s\n", sha, bodyFile)
	repo, err := finishRepository(start.Remote.URL)
	if err != nil {
		return fail(err.Error(), "finish needs the start record's GitHub repository", steps(sha))
	}
	seam := finishDispatchFactory()
	hasLabel, err := seam.Prepare(ctx, wt, repo)
	if err != nil {
		return fail("gh can't reach GitHub as you (F8): "+err.Error(), "finish needs gh and a label lookup before approval", "gh auth login; "+steps(sha))
	}
	scope := broker.FinishApprovalScope{
		Repository: repo, Remote: start.Remote.URL, Head: start.Branch, SHA: sha, Base: start.Base.Branch,
		TitleSHA256: finishSHA([]byte(actualTitle)), BodySHA256: finishSHA([]byte(rendered)),
		Draft: !*ready, CreateLabel: !hasLabel, ExpiresAt: time.Now().Unix() + 900,
	}
	scopeText := broker.FormatFinishApprovalScope(scope)
	if _, err := broker.ParseFinishApprovalScope(scopeText); err != nil {
		return fail(err.Error(), "approval must bind every byte", steps(sha))
	}
	finishPrompt(tty, scope, actualTitle, rendered, protected)
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(answer) != "y" {
		return fail("pull request declined", "opening a PR requires typed y at the terminal; commit kept at "+sha, steps(sha))
	}
	var approval daemon.FinishApproval
	if err := finishOperatorPost(ctx, c, credentials, ep, daemon.FinishRequest{Phase: "approve", Scope: scopeText}, &approval); err != nil {
		return fail(err.Error(), "approval was not confirmed recorded", steps(sha))
	}
	if finishBeforeDispatch != nil {
		finishBeforeDispatch()
	}
	head, headErr := query("rev-parse", "HEAD")
	remote, remoteErr := query("remote", "get-url", start.Remote.Name)
	currentBody, bodyErr := os.ReadFile(bodyFile)
	if headErr != nil || remoteErr != nil || bodyErr != nil || strings.TrimSpace(head) != scope.SHA ||
		strings.TrimSpace(remote) != scope.Remote || finishSHA(currentBody) != scope.BodySHA256 ||
		finishSHA([]byte(actualTitle)) != scope.TitleSHA256 {
		return fail("the commit or body changed since approval (F6)", "approval binds HEAD, remote, title and exact body bytes", steps(sha))
	}
	outcome, pr, dispatchErr := finishDispatchOutcome(ctx, seam, git, wt, start.Remote.Name, repo, scope.Head,
		scope.SHA, scope.Base, actualTitle, bodyFile, scope.Draft, scope.CreateLabel)
	checkRefs := make([]string, 0, len(results))
	for _, result := range results {
		checkRefs = append(checkRefs, result.Ref)
	}
	finish := daemon.EpisodeFinish{
		Episode: ep, Branch: scope.Head, Base: start.Base.SHA, Commit: scope.SHA, Remote: scope.Remote,
		PR: pr, Approval: approval.Approval, Checks: checkRefs, Unattributed: unattributed,
		ReportSHA256: scope.BodySHA256, Trailer: trailer, Outcome: outcome, FinishedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := finishOperatorPost(ctx, c, credentials, ep, daemon.FinishRequest{Phase: "record", Record: &finish}, nil); err != nil {
		fmt.Fprintln(stderr, presetRefusal(err.Error(), "finish outcome was not confirmed recorded; session kept", "inspect the episode report; "+steps(sha)))
		return exitFatal
	}
	creds, err := store.OpenCredentials(credentials)
	if err == nil {
		err = revokeEpisodeCredentials(ctx, creds, ep, stdout)
		creds.Close()
	}
	if err != nil {
		fmt.Fprintln(stderr, presetRefusal("PR recorded, session NOT closed (F10)", err.Error(), "ailang-worldd session revoke --episode "+ep+" --db "+p.DB))
		return exitFatal
	}
	if dispatchErr != nil {
		what := "gh can't reach GitHub as you (F8)"
		if outcome == "push-failed" {
			what = "push failed (F7)"
		}
		if outcome == "ambiguous" {
			what = "the push or PR may have happened (F9)"
		}
		remaining := steps(sha)
		if strings.Contains(dispatchErr.Error(), "not an ancestor") {
			remaining = "resolve the diverged branch by hand; never force; " + remaining
		}
		if outcome == "ambiguous" {
			remaining = "inspect git ls-remote and gh pr view before any write; " + remaining
		}
		if outcome == "pr-failed" {
			parts := strings.SplitN(remaining, "\n", 2)
			if len(parts) == 2 {
				remaining = parts[1]
			}
		}
		fmt.Fprintln(stderr, presetRefusal(what+": "+dispatchErr.Error(), "recorded outcome "+outcome, remaining))
		return exitFatal
	}
	fmt.Fprintln(stdout, pr.URL)
	return exitOK
}

// finishPreset intentionally does not re-pin or validate mutable checks.
func finishPreset(name, dir string) (presetDaemon, error) {
	if !workspaceEpisodePattern.MatchString(name) || len(name) > maxPresetNameLen {
		return presetDaemon{}, fmt.Errorf("invalid preset name %q", name)
	}
	path, err := expandPresetPath(filepath.Join(dir, name+".json"))
	if err != nil {
		return presetDaemon{}, err
	}
	if err := presetLayout(path, defaultWorldPaths().workspaceRoot); err != nil {
		return presetDaemon{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return presetDaemon{}, err
	}
	var p launchPreset
	if err := strictPresetJSON(data, &p); err != nil {
		return presetDaemon{}, err
	}
	if p.Schema != "world/preset/v1" || p.Name != name {
		return presetDaemon{}, fmt.Errorf("invalid preset schema or name")
	}
	d, err := resolvePresetDaemon(p.Daemon)
	if err == nil {
		err = presetLayout(path, d.WorkspaceRoot)
	}
	return d, err
}

func containsString(values []string, value string) bool {
	for _, s := range values {
		if s == value {
			return true
		}
	}
	return false
}

func validateFinishChecks(checks []store.StartCheck, allowed []string, wt string) error {
	names := map[string]bool{}
	for _, c := range checks {
		if c.Name == "" || names[c.Name] {
			return fmt.Errorf("empty or duplicate check name %q", c.Name)
		}
		names[c.Name] = true
		if len(c.Argv) == 0 || len(c.Argv) > maxOutcomeArgv {
			return fmt.Errorf("check %s: argv needs 1..%d items", c.Name, maxOutcomeArgv)
		}
		for _, field := range append([]string{c.Name}, c.Argv...) {
			if field == "" || len(field) > maxOutcomeFieldBytes {
				return fmt.Errorf("check %s: field must have 1..%d bytes", c.Name, maxOutcomeFieldBytes)
			}
		}
		if !filepath.IsLocal(c.Cwd) || c.TimeoutS <= 0 || c.TimeoutS > 86400 {
			return fmt.Errorf("check %s: invalid cwd or timeout_s", c.Name)
		}
		// Resolve symlinks too: a lexically local cwd may point outside.
		root, rootErr := filepath.EvalSymlinks(wt)
		dir, err := filepath.EvalSymlinks(filepath.Join(wt, c.Cwd))
		if rootErr != nil || err != nil {
			return fmt.Errorf("check %s: cannot resolve cwd", c.Name)
		}
		rel, err := filepath.Rel(root, dir)
		if err != nil || !filepath.IsLocal(rel) {
			return fmt.Errorf("check %s: cwd leaves worktree", c.Name)
		}
		if !c.Unsandboxed && c.Exec == "" {
			return fmt.Errorf("check %s: missing pinned exec", c.Name)
		}
	}
	for _, name := range allowed {
		if !names[name] {
			return fmt.Errorf("--allow-failing %s is not a check", name)
		}
	}
	return nil
}

// finishChanged uses NUL records so spaces, newlines and rename sources
// cannot be confused with porcelain syntax.
func finishChanged(query func(...string) (string, error), base string) ([]string, error) {
	status, err := query("status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	diff, err := query("diff", "--name-only", "-z", "--no-renames", base, "HEAD", "--")
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	add := func(p string) error {
		if !filepath.IsLocal(p) || filepath.Clean(p) == "." {
			return fmt.Errorf("invalid changed path %q", p)
		}
		set[p] = true
		return nil
	}
	records := strings.Split(status, "\x00")
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if rec == "" {
			continue
		}
		if len(rec) < 4 || rec[2] != ' ' {
			return nil, fmt.Errorf("malformed git status")
		}
		if err := add(rec[3:]); err != nil {
			return nil, err
		}
		if strings.ContainsAny(rec[:2], "RC") {
			i++
			if i >= len(records) || records[i] == "" {
				return nil, fmt.Errorf("malformed git rename")
			}
			if err := add(records[i]); err != nil {
				return nil, err
			}
		}
	}
	for _, p := range strings.Split(diff, "\x00") {
		if p != "" {
			if err := add(p); err != nil {
				return nil, err
			}
		}
	}
	out := make([]string, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Strings(out)
	return out, nil
}

func finishUnattributed(paths []string, r daemon.EpisodeReport) []string {
	known := map[string]bool{}
	for _, p := range r.Writes.Written {
		known[filepath.Clean(p.Path)] = true
	}
	for _, p := range r.Writes.Produced {
		known[filepath.Clean(p.Path)] = true
	}
	var out []string
	for _, p := range paths {
		if !known[filepath.Clean(p)] {
			out = append(out, p)
		}
	}
	return out
}

func finishSandboxCheck(ctx context.Context, c *client, token string, check store.StartCheck) (checkrun.Result, error) {
	r := checkrun.Result{ExitCode: -1}
	out, err := (mcpClient{c: c, token: token}).post(ctx, "tools/call", 1, map[string]any{
		"name": "workspace-exec", "arguments": map[string]any{"command": check.Exec, "args": check.Args},
	})
	if err != nil {
		return r, err
	}
	if out.Kind != mcpResult {
		return r, fmt.Errorf("workspace-exec failed: %s %v", out.Body, out.Error)
	}
	var envelope struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if json.Unmarshal(out.Result, &envelope) != nil || len(envelope.Content) == 0 {
		return r, fmt.Errorf("malformed workspace-exec result")
	}
	var result struct {
		ExitCode *int   `json:"exit_code"`
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
		TimedOut bool   `json:"timed_out"`
		Refused  string `json:"refused"`
	}
	if err := json.Unmarshal([]byte(envelope.Content[0].Text), &result); err != nil {
		return r, err
	}
	r.Stdout, r.Stderr = result.Stdout, result.Stderr
	if envelope.IsError || result.ExitCode == nil || result.TimedOut {
		return r, fmt.Errorf("workspace-exec refused or timed out: %s", result.Refused)
	}
	r.ExitCode = *result.ExitCode
	return r, nil
}

func clipOutcomeHead(s string) string {
	if len(s) > maxOutcomeHeadBytes {
		s = s[:maxOutcomeHeadBytes]
	}
	return strings.ToValidUTF8(s, "")
}

func finishOutcome(rep daemon.EpisodeReport, c store.StartCheck, r checkrun.Result, passed bool) taskOutcome {
	class, category := "PASS", ""
	if !passed {
		class, category = "FAIL", "grader_failed"
	}
	where := "sandboxed"
	if c.Unsandboxed {
		where = "unsandboxed"
	}
	argv := make([]string, 0, len(c.Argv))
	for _, a := range c.Argv {
		argv = append(argv, shellQuote(a))
	}
	rule := "the operator runs " + strings.Join(argv, " ") + " " + where + " in the episode worktree"
	if c.Cwd != "." {
		rule += " (cwd " + c.Cwd + ")"
	}
	rule += " after the agent stops"
	// Keep the writer's 1024-byte field bound even for a long command.
	if len(rule) > maxOutcomeFieldBytes {
		rule = strings.ToValidUTF8(rule[:maxOutcomeFieldBytes], "")
	}
	from := max(int64(1), rep.Entries[0])
	return taskOutcome{
		Agent: rep.Start.Agent, Arm: "world", Entries: &outcomeEntries{From: from, To: max(from-1, rep.Entries[1])},
		EpisodeID: rep.Episode, Run: 1, Task: c.Name, Label: "episode-finish",
		Grading: outcomeGrading{Argv: c.Argv, InputFiles: []string{}, Rule: rule, GraderVersion: "trial-grader/1"},
		Outcome: outcomeResult{Class: class, ErrorCategory: category, CompileOk: passed, RuntimeOk: passed, StdoutOk: passed,
			StdoutHead: clipOutcomeHead(r.Stdout), StderrHead: clipOutcomeHead(r.Stderr)},
	}
}

func finishHandSteps(wt string, start *store.EpisodeStart, sha, file string) string {
	return fmt.Sprintf("git -C %s push %s %s\nnext: gh pr create --base %s --head %s --body-file %s --draft",
		shellQuote(wt), shellQuote(start.Remote.Name), shellQuote(sha+":refs/heads/"+start.Branch),
		shellQuote(start.Base.Branch), shellQuote(start.Branch), shellQuote(file))
}

func writeFinishBody(path, body string) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".pr-body-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := io.WriteString(f, body); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
