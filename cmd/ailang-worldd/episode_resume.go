package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func runEpisodeResume(fs *flag.FlagSet, pos []string, ep, name, presetsDir, profilesDir string, out, errOut io.Writer, env sessionEnv) int {
	next := "ailang-worldd episode start --resume " + ep + " --preset " + name
	fail := func(err error) int {
		fmt.Fprintln(errOut, err)
		return exitUsage
	}
	bad := len(pos) != 0
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "resume", "preset", "presets-dir", "exec-profiles-dir":
		default:
			bad = true
		}
	})
	if bad {
		return fail(presetRefusal("--resume takes only --preset, --presets-dir and --exec-profiles-dir", "resume reuses the pinned start", next))
	}
	if name == "" || name == "se-tools" || !workspaceEpisodePattern.MatchString(ep) {
		return fail(presetRefusal("--resume requires an episode and --preset", "resume reuses the pinned start", next))
	}
	p, err := finishPreset(name, presetsDir)
	if err != nil {
		return fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	path, err := credentialsFor(p.DB, p.CredentialsDB)
	if err != nil {
		return fail(err)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fail(presetRefusal("episode "+ep+" has no start record", "resume reads the pinned start", "ailang-worldd episode start --preset "+name))
	}
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		return fail(err)
	}
	start, err := ro.EpisodeStart(ctx, ep)
	if err != nil {
		_ = ro.Close()
		return fail(err)
	}
	rows, err := ro.ListSessions(ctx, ep, store.MaxListSessions)
	_ = ro.Close()
	if err != nil {
		return fail(err)
	}
	if start == nil {
		return fail(presetRefusal("episode "+ep+" has no start record", "resume reads the pinned start", "ailang-worldd episode start --preset "+name))
	}
	if name != start.Preset.Name {
		return fail(presetRefusal("preset "+name+" does not match start preset "+start.Preset.Name, "resume reuses the pinned start", "use --preset "+start.Preset.Name))
	}
	now := env.now()
	for _, row := range rows {
		if row.ExpiresAt > now {
			left := row.ExpiresAt - now
			return fail(presetRefusal(fmt.Sprintf("episode %s has a live session (%dh%dm left) (F12)", ep, left/3600, left/60%60), "resume re-mints only an expired session", "ailang-worldd episode launch "+ep+" --preset "+name))
		}
	}
	if len(rows) == 0 {
		return fail(presetRefusal("episode "+ep+" has no credential to resume (revoked or finished)", "revoke ends an episode", "ailang-worldd episode start --preset "+name+" --branch "+start.Branch))
	}
	old := rows[len(rows)-1]
	var grants []broker.Capability
	if err := json.Unmarshal([]byte(old.GrantsJSON), &grants); err != nil {
		return fail(err)
	}
	// Check the digest before prepareMint's compatibility check so drift gets
	// the resume-specific remedy even when the real clock trails the test clock.
	if old.Bound() {
		prof, _, err := loadMintProfile(profilesDir, old.Project)
		if err != nil {
			return fail(err)
		}
		if prof.Digest != old.ProfileDigest {
			return fail(presetRefusal("the exec profile "+old.Project+" changed since episode "+ep+" started", "resume keeps the pinned profile", "ailang-worldd episode start --preset "+name))
		}
	}
	req, err := prepareMint(mintOptions{verb: "episode start", episode: ep, db: p.DB, credsDB: p.CredentialsDB,
		wsRoot: p.WorkspaceRoot, repo: start.Repo, branch: start.Branch, grants: grants, ttl: start.TTLS,
		profile: old.Project, profilesDir: profilesDir})
	if err != nil {
		return fail(err)
	}
	if !req.reuse || req.degraded {
		return fail(errors.New("the worktree is gone or isn't git; start a new episode"))
	}
	if old.Bound() && (req.prof == nil || req.prof.Digest != old.ProfileDigest) {
		return fail(presetRefusal("the exec profile "+old.Project+" changed since episode "+ep+" started", "resume keeps the pinned profile", "ailang-worldd episode start --preset "+name))
	}
	git, err := sessionGit()
	if err != nil {
		return fail(err)
	}
	branch, err := worktree.Output(ctx, git, []string{"-C", req.dir, "symbolic-ref", "--short", "HEAD"})
	if err != nil || strings.TrimSpace(branch) != start.Branch {
		return fail(fmt.Errorf("worktree %s is not on start branch %s", req.dir, start.Branch))
	}
	suffix := ".expired-" + strconv.FormatInt(now, 10)
	var oldFiles []string
	for _, ext := range []string{".session", ".mcp.json"} {
		file := filepath.Join(p.SessionsDir, ep+ext)
		if err := launchRegularOrAbsent(file); err != nil {
			return fail(err)
		}
		if _, err := os.Lstat(file + suffix); !errors.Is(err, os.ErrNotExist) {
			return fail(fmt.Errorf("expired target already exists or cannot be inspected: %s", file+suffix))
		}
		if _, err := os.Lstat(file); err == nil {
			oldFiles = append(oldFiles, file)
		}
	}
	if err := checkOutDir(p.SessionsDir, true); err != nil {
		return fail(err)
	}
	if err := launchLockCheck(filepath.Join(p.SessionsDir, ep+".launch.lock")); err != nil {
		return fail(err)
	}
	h, err := fetchDaemonHealth(p.Addr)
	if err != nil {
		return fail(err)
	}
	if same, _ := store.SameDatabase(h.DBPath, p.DB); !same {
		return fail(presetRefusal(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, name, p.DB), "a preset must bind the served store", "fix the preset daemon, or run service status"))
	}
	if err := checkDaemon(h, req, p.Addr, episodeSelfBuild()); err != nil {
		return fail(err)
	}
	presetPath, err := expandPresetPath(filepath.Join(presetsDir, name+".json"))
	if err != nil {
		return fail(err)
	}
	bytes, err := os.ReadFile(presetPath)
	if err != nil {
		return fail(err)
	}
	digest := serviceHash(bytes)
	presetStatus := presetPath + " unchanged"
	if digest != start.Preset.SHA256 {
		presetStatus = presetPath + " changed since start (now " + digest + "); resume keeps the pinned bytes, a new episode start uses the new ones"
	}
	req.presetPrompt = fmt.Sprintf("resume episode %s (started %s; preset %s, pinned %s)\n  preset    %s\n  expired   credential %s expired %s; %s.session and %s.mcp.json are kept as %s.session%s and %s.mcp.json%s\n  worktree  %s on existing branch %s (adopted)\n  agent     %s %s\n  grants    %s\n  session   %ds\n",
		ep, start.StartedAt, name, start.Preset.SHA256, presetStatus, old.CredentialID, time.Unix(old.ExpiresAt, 0).Format(time.RFC3339), ep, ep, ep, suffix, ep, suffix, req.dir, start.Branch, start.Agent, start.Model, grantList(grants), start.TTLS)
	req.resume = &store.EpisodeResume{ResumedFrom: old.CredentialID, ResumedAt: time.Unix(now, 0).Format(time.RFC3339), TTLS: start.TTLS, PresetFileSHA256: digest}
	req.outPath = filepath.Join(p.SessionsDir, ep+".session")
	req.beforeToken = func() error {
		for _, file := range oldFiles {
			if err := os.Rename(file, file+suffix); err != nil {
				return err
			}
		}
		return nil
	}
	if code := provisionEpisode(req, out, errOut, env); code != exitOK {
		return code
	}
	fmt.Fprintf(out, "next: ailang-worldd episode launch %s --preset %s\n", ep, name)
	return exitOK
}
