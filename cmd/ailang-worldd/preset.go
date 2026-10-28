package main

// Presets are operator-local authority configuration, so validation and git
// inspection belong at this host boundary, outside the pure World kernel.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const presetHelp = `usage: ailang-worldd preset list|show|check|init [name] [--presets-dir D]
  show name [--json]                file, digest and derived values
  check name [--no-daemon] [--exec-profiles-dir D]
                                   read-only validation; never runs checks
  init name --repo D --base B [--remote origin] [--agent codex] [--profile P]
                                   creates a private file once; never overwrites

Presets default to ~/.ailang-world/presets. Checks use trial.py's grader
shape; timeout_s defaults to 900. Sandboxed checks must fit the exec profile.
Example:
  ailang-worldd preset init my-repo --repo /path/to/repo --base dev
  ailang-worldd preset check my-repo --no-daemon
  ailang-worldd episode start --preset my-repo --task /path/to/task.md
The last command requires a human y at the controlling terminal.
`

type presetGrant struct {
	Effect string `json:"effect,omitempty"`
	Scope  string `json:"scope,omitempty"`
	Budget int64  `json:"budget"`
	Set    string `json:"set,omitempty"`
}
type presetDaemon struct {
	Addr          string `json:"addr"`
	DB            string `json:"db"`
	CredentialsDB string `json:"credentials_db"`
	WorkspaceRoot string `json:"workspace_root"`
	SessionsDir   string `json:"sessions_dir"`
}
type presetLaunch struct {
	Mode      string  `json:"mode"`
	DeadlineS float64 `json:"deadline_s,omitempty"`
}
type launchPreset struct {
	Schema  string            `json:"schema"`
	Name    string            `json:"name"`
	Repo    string            `json:"repo"`
	Remote  string            `json:"remote"`
	Base    string            `json:"base"`
	Profile string            `json:"profile,omitempty"`
	Agent   string            `json:"agent"`
	Model   string            `json:"model,omitempty"`
	TTLS    int64             `json:"ttl_s"`
	Grants  []presetGrant     `json:"grants"`
	Daemon  json.RawMessage   `json:"daemon,omitempty"`
	Launch  *presetLaunch     `json:"launch,omitempty"`
	Checks  []json.RawMessage `json:"checks,omitempty"`
}
type resolvedPreset struct {
	Name          string              `json:"name"`
	SHA256        string              `json:"sha256"`
	Path          string              `json:"path"`
	File          string              `json:"file"`
	Preset        launchPreset        `json:"preset"`
	BaseSHA       string              `json:"base_sha"`
	RemoteURL     string              `json:"remote_url"`
	Lag           string              `json:"lag"`
	Daemon        presetDaemon        `json:"daemon"`
	Grants        []broker.Capability `json:"grants"`
	Checks        []store.StartCheck  `json:"checks"`
	ProfileDigest string              `json:"profile_digest,omitempty"`
}

// UnmarshalJSON enforces the two closed grant forms, including a budget.
func (g *presetGrant) UnmarshalJSON(data []byte) error {
	type plain presetGrant
	var value plain
	if err := strictPresetJSON(data, &value); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if b, ok := fields["budget"]; !ok || string(b) == "null" {
		return errors.New("grant needs budget")
	}
	if _, set := fields["set"]; set {
		if len(fields) != 2 || value.Set != "se-tools" {
			return errors.New("set grant needs exactly set se-tools and budget")
		}
	} else if len(fields) != 3 || value.Effect == "" || value.Scope == "" {
		return errors.New("grant needs exactly effect, scope and budget")
	}
	*g = presetGrant(value)
	return nil
}

func defaultPresetsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ailang-world", "presets")
}
func presetRefusal(what, why, next string) error {
	return fmt.Errorf("%s\n  why: %s\n  next: %s", what, why, next)
}
func strictPresetJSON(data []byte, dst any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(dst); err != nil {
		return err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return errors.New("expected one JSON value")
	}
	return nil
}
func expandPresetPath(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return filepath.Abs(path)
}

// presetLayout precedes parsing and never consults a daemon.
func presetLayout(path, root string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return presetRefusal("preset "+path+" is a symlink or not a regular file", "grants and checks must be operator-owned", "use a regular private preset file")
	}
	if info.Mode().Perm()&0o022 != 0 {
		return presetRefusal(fmt.Sprintf("preset %s is group- or other-writable (%04o)", path, info.Mode().Perm()), "another user could widen its authority", "chmod 600 "+path)
	}
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	r, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if canonical, err := filepath.EvalSymlinks(r); err == nil {
		r = canonical
	}
	if within(r, p) {
		return presetRefusal("preset "+path+" is inside the workspace root", "an agent could change the next mint's authority", "move it outside "+r)
	}
	return nil
}
func presetNames(dir string) []string {
	entries, _ := os.ReadDir(dir)
	names := []string{}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	sort.Strings(names)
	return names
}

// validateGraderChecks mirrors trial.episode_spec's grader validation. The
// shared JSON vectors pin defaults and refusals in both languages.
func validateGraderChecks(raw []json.RawMessage, episode string) ([]store.StartCheck, error) {
	out := []store.StartCheck{}
	names := map[string]bool{}
	for _, data := range raw {
		var g struct {
			Name        string   `json:"name"`
			Argv        []string `json:"argv"`
			Cwd         string   `json:"cwd"`
			TimeoutS    float64  `json:"timeout_s"`
			Unsandboxed bool     `json:"unsandboxed"`
		}
		if err := strictPresetJSON(data, &g); err != nil {
			var typed *json.UnmarshalTypeError
			if errors.As(err, &typed) && typed.Field == "argv" {
				return nil, fmt.Errorf("episode %s: grader argv must be a non-empty list of strings", episode)
			}
			return nil, err
		}
		valid := len(g.Argv) > 0
		for _, a := range g.Argv {
			if a == "" {
				valid = false
			}
		}
		if !valid {
			return nil, fmt.Errorf("episode %s: grader argv must be a non-empty list of strings", episode)
		}
		if g.Name == "" {
			g.Name = g.Argv[0][strings.LastIndex(g.Argv[0], "/")+1:]
		}
		if names[g.Name] {
			return nil, fmt.Errorf("episode %s: duplicate grader name '%s'", episode, g.Name)
		}
		names[g.Name] = true
		if g.Cwd == "" {
			g.Cwd = "."
		}
		clean := filepath.Clean(g.Cwd)
		if filepath.IsAbs(g.Cwd) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
			return nil, fmt.Errorf("episode %s: grader cwd '%s' must be relative inside the worktree", episode, g.Cwd)
		}
		if g.TimeoutS == 0 {
			g.TimeoutS = 900
		}
		class := "sandboxed"
		if g.Unsandboxed {
			class = "UNSANDBOXED"
		}
		out = append(out, store.StartCheck{Name: g.Name, Argv: g.Argv, Cwd: g.Cwd, TimeoutS: g.TimeoutS, Unsandboxed: g.Unsandboxed, Sandbox: class})
	}
	return out, nil
}

func pinPresetChecks(checks []store.StartCheck, p *broker.ExecProfile) error {
	for i := range checks {
		c := &checks[i]
		if c.Unsandboxed {
			continue
		}
		nearest, best := "(no profile)", -1
		if p != nil {
			ids := []string{}
			for id := range p.Commands {
				ids = append(ids, id)
			}
			sort.Strings(ids)
			for _, id := range ids {
				cmd := p.Commands[id]
				n := 0
				for n < len(cmd.Argv) && n < len(c.Argv) && cmd.Argv[n] == c.Argv[n] {
					n++
				}
				if n > best {
					nearest, best = id, n
				}
				if n != len(cmd.Argv) {
					continue
				}
				args, err := broker.MatchExecArgs(cmd.ExecCommand, c.Argv[n:])
				if err != nil || filepath.Clean(c.Cwd) != filepath.Clean(p.Root) || c.TimeoutS <= 0 || c.TimeoutS*1000 > float64(p.TimeoutMS) {
					continue
				}
				c.Exec, c.Args = id, args
				break
			}
		}
		if c.Exec == "" {
			return presetRefusal(fmt.Sprintf("check %s: no matching profile command for %q (nearest command: %s)", c.Name, strings.Join(c.Argv, " "), nearest),
				"the argv, cwd, grammar and timeout must fit one sandbox command", "add \"unsandboxed\": true, or a profile command; row 145 covers long runs")
		}
	}
	return nil
}

func resolvePresetDaemon(raw json.RawMessage) (presetDaemon, error) {
	def := defaultWorldPaths()
	d := presetDaemon{Addr: defaultEpisodeAddr, DB: def.db, WorkspaceRoot: def.workspaceRoot, SessionsDir: defaultSessionsDir()}
	if len(raw) == 0 || string(raw) == "\"service\"" {
		home, _ := os.UserHomeDir()
		r, err := readServiceRecord(homeServicePaths(home))
		if err != nil {
			return d, presetRefusal("cannot read preset daemon service: "+err.Error(), "daemon defaults come from service.json", "ailang-worldd service status")
		}
		d.Addr, d.DB, d.CredentialsDB, d.WorkspaceRoot = r.Address, r.Store, r.Credentials, r.Workspace
		d.SessionsDir = filepath.Join(filepath.Dir(r.Store), "sessions")
	} else if err := strictPresetJSON(raw, &d); err != nil {
		return d, err
	}
	var err error
	for _, path := range []*string{&d.DB, &d.WorkspaceRoot, &d.SessionsDir} {
		*path, err = expandPresetPath(*path)
		if err != nil {
			return d, err
		}
	}
	if d.CredentialsDB != "" {
		d.CredentialsDB, err = expandPresetPath(d.CredentialsDB)
	}
	if err != nil {
		return d, err
	}
	d.Addr, err = loopbackAddr(d.Addr)
	return d, err
}

func loadLaunchPreset(ctx context.Context, name, dir, profilesDir, workspaceRoot string) (resolved *resolvedPreset, loadErr error) {
	defer func() {
		if loadErr != nil && !strings.Contains(loadErr.Error(), "\n  why:") {
			loadErr = presetRefusal("preset "+name+": "+loadErr.Error(), "the preset cannot safely provision an episode", "fix the field, then preset check "+name)
		}
	}()
	if name == "se-tools" {
		return nil, presetRefusal("preset se-tools is reserved", "it names the existing grant set", "choose another preset name")
	}
	if !workspaceEpisodePattern.MatchString(name) || len(name) > 53 {
		return nil, presetRefusal("invalid preset name "+name, "use the episode grammar, at most 53 characters", "choose a shorter lowercase name")
	}
	path, err := expandPresetPath(filepath.Join(dir, name+".json"))
	if err != nil {
		return nil, err
	}
	if err := presetLayout(path, workspaceRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, presetRefusal(fmt.Sprintf("no preset %s in %s (have: %s)", name, dir, strings.Join(presetNames(dir), ", ")), "a launch preset must exist", "ailang-worldd preset list --presets-dir "+dir)
		}
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p launchPreset
	if err := strictPresetJSON(data, &p); err != nil {
		return nil, presetRefusal("preset "+name+": "+err.Error(), "preset fields are closed", "fix the field, then preset check "+name)
	}
	if p.Name != name || p.Schema != "world/preset/v1" {
		return nil, presetRefusal("preset "+name+": name must equal the file basename and schema must be world/preset/v1", "the file identifies its authority", "fix the name and schema")
	}
	// An explicit workspace root is inspected before any other field.
	var layoutDaemon struct {
		WorkspaceRoot string `json:"workspace_root"`
	}
	if len(p.Daemon) > 0 && p.Daemon[0] == '{' {
		_ = json.Unmarshal(p.Daemon, &layoutDaemon)
		if layoutDaemon.WorkspaceRoot != "" {
			root, err := expandPresetPath(layoutDaemon.WorkspaceRoot)
			if err != nil {
				return nil, err
			}
			if err := presetLayout(path, root); err != nil {
				return nil, err
			}
		}
	}
	d, err := resolvePresetDaemon(p.Daemon)
	if err != nil {
		return nil, err
	}
	if err := presetLayout(path, d.WorkspaceRoot); err != nil {
		return nil, err
	}
	if p.Agent != "claude" && p.Agent != "codex" {
		return nil, errors.New("preset agent must be claude or codex")
	}
	if p.TTLS <= 0 || len(p.Grants) == 0 {
		return nil, errors.New("preset needs positive ttl_s and a non-empty grants list")
	}
	if p.Launch != nil {
		if p.Launch.Mode == "" {
			p.Launch.Mode = "interactive"
		}
		if p.Launch.Mode != "interactive" && p.Launch.Mode != "headless" {
			return nil, errors.New("launch mode must be interactive or headless")
		}
		if p.Launch.DeadlineS < 0 {
			return nil, errors.New("launch deadline_s must be non-negative")
		}
	}
	r := &resolvedPreset{Name: name, SHA256: serviceHash(data), Path: path, File: string(data), Preset: p, Daemon: d}
	for _, g := range p.Grants {
		if g.Budget < 0 {
			return nil, errors.New("grant budget must be non-negative")
		}
		var grants []broker.Capability
		if g.Set != "" {
			if g.Set != "se-tools" || g.Effect != "" || g.Scope != "" {
				return nil, errors.New("grant set must be se-tools, without effect or scope")
			}
			grants, err = sessionNewGrants("se-tools", nil, g.Budget)
		} else {
			if g.Effect == "" || g.Scope == "" {
				return nil, errors.New("grant needs effect, scope and budget")
			}
			grants, err = sessionNewGrants("", []string{fmt.Sprintf("%s=%s:%d", g.Effect, g.Scope, g.Budget)}, 0)
		}
		if err != nil {
			return nil, err
		}
		r.Grants = append(r.Grants, grants...)
	}
	r.Checks, err = validateGraderChecks(p.Checks, name)
	if err != nil {
		return nil, err
	}
	var profile *broker.ExecProfile
	if p.Profile != "" {
		profile, _, err = loadMintProfile(profilesDir, p.Profile)
		if err != nil {
			return nil, err
		}
		r.ProfileDigest = profile.Digest
	} else if hasGrant(r.Grants, broker.EffectWorkspaceExec) {
		return nil, errors.New("Workspace.Exec requires a preset profile")
	}
	if err := pinPresetChecks(r.Checks, profile); err != nil {
		return nil, err
	}
	r.Preset.Repo, err = expandPresetPath(p.Repo)
	if err != nil || p.Repo == "" {
		return nil, errors.New("preset repo must name a git repository")
	}
	git, err := sessionGit()
	if err != nil {
		return nil, err
	}
	query := func(args ...string) (string, error) {
		out, err := worktree.Output(ctx, git, append([]string{"-C", r.Preset.Repo}, args...))
		return strings.TrimSpace(out), err
	}
	if worktree.IsLinked(r.Preset.Repo) {
		return nil, errors.New("preset repo must not be a linked worktree")
	}
	if _, err := query("rev-parse", "--git-dir"); err != nil {
		return nil, err
	}
	if top, err := query("rev-parse", "--show-toplevel"); err == nil && worktree.IsLinked(top) {
		return nil, errors.New("preset repo must not be a linked worktree")
	}
	if _, err := query("check-ref-format", "refs/heads/"+p.Base); err != nil || p.Base == "" {
		return nil, errors.New("preset base must name a local branch")
	}
	r.BaseSHA, err = query("rev-parse", "--verify", "refs/heads/"+p.Base+"^{commit}")
	if err != nil {
		return nil, fmt.Errorf("preset %s: base %q is not a branch of %s", name, p.Base, p.Repo)
	}
	remotes, err := query("remote")
	if err != nil {
		return nil, err
	}
	found := false
	for _, remote := range strings.Split(remotes, "\n") {
		if remote == p.Remote && remote != "" {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("preset remote %q is not a remote in %s", p.Remote, p.Repo)
	}
	r.RemoteURL, err = query("remote", "get-url", p.Remote)
	if err != nil {
		return nil, err
	}
	ref := "refs/remotes/" + p.Remote + "/" + p.Base
	r.Lag = "never fetched"
	if _, err := query("rev-parse", "--verify", ref+"^{commit}"); err == nil {
		lag, err := query("rev-list", "--count", "refs/heads/"+p.Base+".."+ref)
		if err != nil {
			return nil, err
		}
		r.Lag = lag + " behind " + p.Remote + "/" + p.Base + " as last fetched"
	}
	return r, nil
}

func presetDaemonCheck(ctx context.Context, p *resolvedPreset, profilesDir string) error {
	req, err := prepareMint(mintOptions{verb: "preset check", episode: p.Name, db: p.Daemon.DB, credsDB: p.Daemon.CredentialsDB, wsRoot: p.Daemon.WorkspaceRoot,
		repo: p.Preset.Repo, grants: p.Grants, ttl: p.Preset.TTLS, profile: p.Preset.Profile, profilesDir: profilesDir})
	if err != nil {
		return err
	}
	h, err := fetchDaemonHealth(p.Daemon.Addr)
	if err != nil {
		return err
	}
	return checkPresetDaemon(h, req, p, p.Daemon.Addr)
}

func runPreset(args []string, out, errOut io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(out, presetHelp)
		return exitOK
	}
	verb := args[0]
	fs := flag.NewFlagSet("preset "+verb, flag.ContinueOnError)
	fs.SetOutput(errOut)
	dir := fs.String("presets-dir", defaultPresetsDir(), "preset files")
	profiles := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "exec profiles")
	noDaemon := fs.Bool("no-daemon", false, "skip daemon check")
	asJSON := fs.Bool("json", false, "JSON output")
	repo := fs.String("repo", "", "main checkout")
	base := fs.String("base", "", "local base branch")
	remote := fs.String("remote", "origin", "remote name")
	agent := fs.String("agent", "codex", "agent")
	profile := fs.String("profile", "", "exec profile")
	pos, err := parseInterleaved(fs, args[1:])
	fail := func(err error) int { fmt.Fprintln(errOut, err); return exitUsage }
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(out, presetHelp)
		return exitOK
	}
	if err != nil {
		return fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if verb == "init" {
		if len(pos) != 1 || *repo == "" || *base == "" {
			fmt.Fprint(errOut, presetHelp)
			return exitUsage
		}
		name := pos[0]
		if name == "se-tools" || !workspaceEpisodePattern.MatchString(name) || len(name) > 53 {
			return fail(errors.New("invalid or reserved preset name"))
		}
		p := launchPreset{Schema: "world/preset/v1", Name: name, Repo: *repo, Base: *base, Remote: *remote, Agent: *agent, Profile: *profile, TTLS: 43200,
			Grants: []presetGrant{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 5000}}, Daemon: json.RawMessage("\"service\"")}
		if *agent != "codex" && *agent != "claude" {
			return fail(errors.New("agent must be codex or claude"))
		}
		resolvedRepo, err := expandPresetPath(*repo)
		if err != nil {
			return fail(err)
		}
		if worktree.IsLinked(resolvedRepo) {
			return fail(errors.New("preset repo must not be a linked worktree"))
		}
		git, err := sessionGit()
		if err != nil {
			return fail(err)
		}
		for _, args := range [][]string{{"rev-parse", "--git-dir"}, {"check-ref-format", "refs/heads/" + *base}, {"rev-parse", "--verify", "refs/heads/" + *base + "^{commit}"}, {"remote", "get-url", *remote}} {
			if _, err := worktree.Output(ctx, git, append([]string{"-C", resolvedRepo}, args...)); err != nil {
				return fail(err)
			}
		}
		if *profile != "" {
			if _, _, err := loadMintProfile(*profiles, *profile); err != nil {
				return fail(err)
			}
		}
		absDir, err := expandPresetPath(*dir)
		if err != nil {
			return fail(err)
		}
		root := defaultWorldPaths().workspaceRoot
		home, _ := os.UserHomeDir()
		if service, err := readServiceRecord(homeServicePaths(home)); err == nil {
			root = service.Workspace
		}
		canonicalDir := absDir
		if canonical, err := filepath.EvalSymlinks(absDir); err == nil {
			canonicalDir = canonical
		}
		if within(root, canonicalDir) {
			return fail(presetRefusal("preset directory is inside the workspace root", "an agent could change its authority", "move presets outside "+root))
		}
		*dir = absDir
		data, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return fail(err)
		}
		if err := os.MkdirAll(*dir, 0o700); err != nil {
			return fail(err)
		}
		path := filepath.Join(*dir, name+".json")
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return fail(err)
		}
		_, err = f.Write(append(data, '\n'))
		closeErr := f.Close()
		if err != nil {
			return fail(err)
		}
		if closeErr != nil {
			return fail(closeErr)
		}
		fmt.Fprintln(out, string(data))
		return exitOK
	}
	if verb == "list" {
		for _, name := range presetNames(*dir) {
			p, err := loadLaunchPreset(ctx, name, *dir, *profiles, defaultWorldPaths().workspaceRoot)
			if err != nil {
				return fail(err)
			}
			fmt.Fprintf(out, "%s\t%s\t%s\t%s\n", p.Name, p.Preset.Repo, p.Preset.Agent, strings.TrimPrefix(p.SHA256, "sha256:")[:12])
		}
		return exitOK
	}
	if len(pos) != 1 || (verb != "show" && verb != "check") {
		fmt.Fprint(errOut, presetHelp)
		return exitUsage
	}
	p, err := loadLaunchPreset(ctx, pos[0], *dir, *profiles, defaultWorldPaths().workspaceRoot)
	if err != nil {
		return fail(err)
	}
	if verb == "show" {
		if *asJSON {
			data, _ := json.MarshalIndent(p, "", "  ")
			fmt.Fprintln(out, string(data))
		} else {
			fmt.Fprintf(out, "%s\nsha256 %s\nbase %s @ %s\n", p.File, p.SHA256, p.Preset.Base, p.BaseSHA)
			data, _ := json.MarshalIndent(p, "", "  ")
			fmt.Fprintln(out, string(data))
		}
		return exitOK
	}
	if !*noDaemon {
		if err := presetDaemonCheck(ctx, p, *profiles); err != nil {
			return fail(err)
		}
	} else {
		fmt.Fprintln(out, "daemon match skipped (--no-daemon)")
	}
	fmt.Fprintf(out, "base %s @ %.7s\n", p.Preset.Base, p.BaseSHA)
	for _, c := range p.Checks {
		class := "UNSANDBOXED"
		if !c.Unsandboxed {
			class = "sandboxed: " + c.Exec
		}
		fmt.Fprintf(out, "check %s: %s\n", c.Name, class)
	}
	return exitOK
}
