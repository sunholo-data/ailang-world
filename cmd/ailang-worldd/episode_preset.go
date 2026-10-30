package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func derivePresetEpisode(name, root, sessions string, now int64) (string, error) {
	base := name + "-" + time.Unix(now, 0).Local().Format("0601021504")
	for n := 1; n <= 9; n++ {
		id := base
		if n > 1 {
			id += "-" + strconv.Itoa(n)
		}
		free := true
		for _, path := range []string{filepath.Join(root, id), filepath.Join(sessions, id+".session")} {
			_, err := os.Lstat(path)
			if err == nil {
				free = false
			} else if !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
		}
		if free {
			if !workspaceEpisodePattern.MatchString(id) {
				return "", fmt.Errorf("derived episode %s exceeds episode grammar; pass an episode id", id)
			}
			return id, nil
		}
	}
	return "", presetRefusal(base+" … -9 all exist (F5)", "every derived id is occupied", "pass an episode id")
}

func checkPresetDaemon(h daemon.HealthResponse, req *mintRequest, p *resolvedPreset, addr string) error {
	if same, _ := store.SameDatabase(h.DBPath, p.Daemon.DB); !same {
		return presetRefusal(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, p.Name, p.Daemon.DB), "a preset must bind the served store", "fix the preset daemon, or run service status")
	}
	return checkDaemon(h, req, addr, episodeSelfBuild())
}

func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDir, task string, specs []string, out, errOut io.Writer, env sessionEnv) int {
	fail := func(err error) int {
		fmt.Fprintln(errOut, "ailang-worldd episode start: "+err.Error())
		return exitUsage
	}
	if len(specs) != 0 {
		return fail(presetRefusal("grants come from preset "+name+" (F7)", "grants and checks are closed", "edit the preset, or drop --preset"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	p, err := loadLaunchPreset(ctx, name, dir, profilesDir, defaultWorldPaths().workspaceRoot)
	if err != nil {
		return fail(err)
	}
	launching := fs.Lookup("launch") != nil && fs.Lookup("launch").Value.String() == "true"
	headless := fs.Lookup("headless") != nil && fs.Lookup("headless").Value.String() == "true"
	if launching {
		if _, _, err := launchMode(p.Preset, headless); err != nil {
			return fail(err)
		}
		mode, _, _ := launchMode(p.Preset, headless)
		if mode == "headless" && task == "" {
			return fail(presetRefusal("--launch needs --task", "headless launch runs the pinned task", "add --task <file>"))
		}
	}
	if len(pos) > 1 {
		return fail(errors.New("expected at most one episode id"))
	}
	ep := ""
	if len(pos) == 1 {
		ep = pos[0]
	} else {
		ep, err = derivePresetEpisode(name, p.Daemon.WorkspaceRoot, p.Daemon.SessionsDir, env.now())
		if err != nil {
			return fail(err)
		}
	}
	if !workspaceEpisodePattern.MatchString(ep) {
		return fail(fmt.Errorf("episode %q does not match %s", ep, workspaceEpisodePattern))
	}
	agent, model, ttl, branch := p.Preset.Agent, p.Preset.Model, p.Preset.TTLS, "world/"+ep
	overrides := []string{}
	overrideLines := []string{}
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "agent":
			overrideLines = append(overrideLines, "agent "+f.Value.String()+" (preset: "+agent+")")
			agent = f.Value.String()
		case "model":
			overrideLines = append(overrideLines, "model "+f.Value.String()+" (preset: "+model+")")
			model = f.Value.String()
		case "ttl":
			overrideLines = append(overrideLines, "ttl "+f.Value.String()+" (preset: "+strconv.FormatInt(ttl, 10)+")")
			ttl, _ = strconv.ParseInt(f.Value.String(), 10, 64)
		case "branch":
			overrideLines = append(overrideLines, "branch "+f.Value.String()+" (preset: "+branch+")")
			branch = f.Value.String()
		default:
			return
		}
		overrides = append(overrides, f.Name)
	})
	if agent != "claude" && agent != "codex" {
		return fail(errors.New("agent must be claude or codex"))
	}
	if launching {
		if _, err := launchBinary(agent); err != nil {
			return fail(err)
		}
		mode, _, _ := launchMode(p.Preset, headless)
		if mode == "interactive" {
			if _, err := launchLauncher(agent, fs.Lookup("hosts-dir").Value.String(), p.Daemon.Addr); err != nil {
				return fail(err)
			}
		}
		launchDir := filepath.Join(filepath.Dir(p.Daemon.SessionsDir), "launch", ep)
		if err := validateLaunchDir(launchDir, p.Daemon.WorkspaceRoot, filepath.Join(p.Daemon.WorkspaceRoot, ep), p.Daemon.SessionsDir); err != nil {
			return fail(err)
		}
		if err := launchLockCheck(filepath.Join(p.Daemon.SessionsDir, ep+".launch.lock")); err != nil {
			return fail(err)
		}
	}
	if ttl < 0 {
		return fail(errors.New("ttl must be non-negative"))
	}
	git, err := sessionGit()
	if err != nil {
		return fail(err)
	}
	query := func(args ...string) (string, error) {
		s, err := worktree.Output(ctx, git, append([]string{"-C", p.Preset.Repo}, args...))
		return strings.TrimSpace(s), err
	}
	if _, err := query("check-ref-format", "refs/heads/"+branch); err != nil {
		return fail(err)
	}
	existing := false
	if _, err := query("show-ref", "--verify", "refs/heads/"+branch); err == nil {
		existing = true
	}
	branchOverridden := false
	for _, o := range overrides {
		if o == "branch" {
			branchOverridden = true
		}
	}
	if existing && !branchOverridden {
		return fail(errors.New("derived branch already exists; pass --branch explicitly or choose another episode"))
	}
	if existing {
		list, err := query("worktree", "list", "--porcelain")
		if err != nil {
			return fail(err)
		}
		if strings.Contains("\n"+list+"\n", "\nbranch refs/heads/"+branch+"\n") {
			return fail(fmt.Errorf("existing branch %s is checked out in another worktree; choose a free branch", branch))
		}
	}
	sessionFile := filepath.Join(p.Daemon.SessionsDir, ep+".session")
	taskCopy := filepath.Join(p.Daemon.SessionsDir, ep+".task.md")
	mcp := ""
	if agent == "claude" {
		mcp = filepath.Join(p.Daemon.SessionsDir, ep+".mcp.json")
	}
	if launching {
		mode, _, _ := launchMode(p.Preset, headless)
		if mode == "interactive" {
			mcp = ""
		}
	}
	for _, path := range []string{sessionFile, taskCopy, mcp} {
		if path == "" {
			continue
		}
		if _, err := os.Lstat(path); err == nil {
			return fail(fmt.Errorf("%s already exists; choose another episode", path))
		} else if !errors.Is(err, os.ErrNotExist) {
			return fail(err)
		}
	}
	var taskData []byte
	var taskRecord *store.StartTask
	if task != "" {
		path, err := filepath.Abs(task)
		if err != nil {
			return fail(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fail(err)
		}
		if !info.Mode().IsRegular() {
			return fail(errors.New("task must be a regular file"))
		}
		taskData, err = os.ReadFile(path)
		if err != nil {
			return fail(err)
		}
		title := strings.TrimSpace(strings.SplitN(string(taskData), "\n", 2)[0])
		title = strings.TrimSpace(strings.TrimLeft(title, "#"))
		taskRecord = &store.StartTask{Path: path, SHA256: serviceHash(taskData), Title: title}
	}
	req, err := prepareMint(mintOptions{verb: "episode start", episode: ep, db: p.Daemon.DB, credsDB: p.Daemon.CredentialsDB,
		wsRoot: p.Daemon.WorkspaceRoot, repo: p.Preset.Repo, branch: branch, grants: p.Grants, ttl: ttl,
		outPath: sessionFile, outDirCreate: true, profile: p.Preset.Profile, profilesDir: profilesDir})
	if err != nil {
		return fail(err)
	}
	if req.prof != nil && req.prof.Digest != p.ProfileDigest {
		return fail(errors.New("exec profile changed during validation; run episode start again"))
	}
	if req.reuse {
		return fail(errors.New("launch preset episode already has a worktree; choose another episode"))
	}
	// A CLI never upgrades an older credentials store.
	if _, err := os.Stat(req.credPath); err == nil {
		ro, err := store.OpenCredentialsReadOnly(req.credPath)
		if err != nil {
			return fail(err)
		}
		version := ro.Version()
		_ = ro.Close()
		if version < store.CredentialsSchemaVersion {
			return fail(&store.UnboundCredentialsStoreError{Path: req.credPath})
		}
	}
	req.startPoint, req.existingBranch = p.BaseSHA, existing
	req.start = &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: time.Unix(env.now(), 0).Format(time.RFC3339),
		Preset: store.StartPreset{Name: name, Path: p.Path, SHA256: p.SHA256, Bytes: p.File}, Repo: p.Preset.Repo,
		Remote: store.StartRemote{Name: p.Preset.Remote, URL: p.RemoteURL}, Base: store.StartBase{Branch: p.Preset.Base, SHA: p.BaseSHA},
		Branch: branch, BranchExisted: existing, Agent: agent, Model: model, Overrides: overrides, Task: taskRecord, Checks: p.Checks, TTLS: ttl}
	var prompt strings.Builder
	fmt.Fprintf(&prompt, "episode %s from preset %s (%s, %s)\n  repo      %s, remote %s = %s\n  base      %s @ %s (%s)\n",
		ep, name, p.SHA256, p.Path, p.Preset.Repo, p.Preset.Remote, p.RemoteURL, p.Preset.Base, p.BaseSHA, p.Lag)
	class := "new branch"
	if existing {
		class = "existing branch"
	}
	fmt.Fprintf(&prompt, "  worktree  %s on %s %s\n  agent     %s %s\n", req.dir, class, branch, agent, model)
	if taskRecord != nil {
		fmt.Fprintf(&prompt, "  task      %s %q (%s)\n", taskRecord.Path, taskRecord.Title, taskRecord.SHA256)
	}
	fmt.Fprintf(&prompt, "  profile   %s %s\n", p.Preset.Profile, p.ProfileDigest)
	for _, c := range p.Checks {
		class := "UNSANDBOXED: finish asks again"
		if !c.Unsandboxed {
			class = "sandboxed: " + c.Exec
		}
		fmt.Fprintf(&prompt, "  check     %s (%s)\n", c.Name, class)
	}
	fmt.Fprintf(&prompt, "  grants    %s\n  session   %ds\n", grantList(p.Grants), ttl)
	for _, line := range overrideLines {
		fmt.Fprintln(&prompt, "  override  "+line)
	}
	req.presetPrompt = prompt.String()
	h, err := fetchDaemonHealth(p.Daemon.Addr)
	if err != nil {
		return fail(err)
	}
	if err := checkPresetDaemon(h, req, p, p.Daemon.Addr); err != nil {
		return fail(err)
	}
	req.afterToken = func(token string) error {
		if taskRecord != nil {
			f, err := os.OpenFile(taskCopy, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
			if err != nil {
				return err
			}
			_, err = f.Write(taskData)
			closeErr := f.Close()
			if err != nil {
				_ = os.Remove(taskCopy)
				return err
			}
			if closeErr != nil {
				_ = os.Remove(taskCopy)
				return closeErr
			}
		}
		if mcp != "" {
			if err := writeMCPConfig(mcp, p.Daemon.Addr, token); err != nil {
				_ = os.Remove(taskCopy)
				return err
			}
		}
		return nil
	}
	if code := provisionEpisode(req, out, errOut, env); code != exitOK {
		return code
	}
	if launching {
		return launchEpisode(ep, name, dir, headless, out, errOut, fs.Lookup("hosts-dir").Value.String())
	}
	taskPath := ""
	if taskRecord != nil {
		taskPath = taskRecord.Path
	}
	printLaunch(out, launchInfo{episode: ep, agent: agent, model: model, addr: p.Daemon.Addr, worktree: req.dir,
		branch: branch, sessionFile: sessionFile, mcpFile: mcp, task: taskPath, grants: p.Grants})
	return exitOK
}
