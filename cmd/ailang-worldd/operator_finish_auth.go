package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const operatorFinishRemedy = "repair PATH, GitHub login or Git credential helper in attended service-runner setup and re-run readiness"

var operatorFinishSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

type operatorGitHubAuth int

func (e operatorGitHubAuth) Error() string {
	return fmt.Sprintf("finish blocked: github-auth (readiness step %d); %s", int(e), operatorFinishRemedy)
}

type operatorFinishAuth struct {
	GH, Git, Repo, Remote, Home, Socket, Path, Helper string
}

// Require the first executable on the approved PATH to match its pin.
func operatorFinishTool(c runnerConfig, dirs []string, name string) string {
	for _, dir := range dirs {
		path := filepath.Join(dir, name)
		st, err := os.Stat(path)
		if err != nil || !st.Mode().IsRegular() || st.Mode().Perm()&0111 == 0 {
			continue
		}
		real, err := filepath.EvalSymlinks(path)
		if err == nil && runnerPath(c.Home, real, c.UID, false, 0) == nil &&
			c.Executables[real] != "" && serviceFileHash(real) == c.Executables[real] {
			return real
		}
		return ""
	}
	return ""
}

func planOperatorFinishAuth(c runnerConfig, remote string) (operatorFinishAuth, error) {
	repo, err := finishRepository(remote)
	if err != nil {
		return operatorFinishAuth{}, operatorGitHubAuth(6)
	}
	dirs := append(append([]string(nil), c.ToolDirs...), filepath.SplitList(runnerSystemPath)...)
	gh, git := operatorFinishTool(c, dirs, "gh"), operatorFinishTool(c, dirs, "git")
	if gh == "" || git == "" {
		return operatorFinishAuth{}, operatorGitHubAuth(2)
	}
	socket, helper := c.FinishSSHSocket, ""
	if strings.HasPrefix(remote, "https://") {
		socket = ""
		helper = "!" + shellQuote(gh) + " auth git-credential"
		if c.FinishHTTPSHelper == "osxkeychain" {
			path := operatorFinishTool(c, dirs, "git-credential-osxkeychain")
			if path == "" {
				return operatorFinishAuth{}, operatorGitHubAuth(6)
			}
			helper = "!" + shellQuote(path)
		}
	}
	if socket != "" {
		st, err := os.Lstat(socket)
		if err != nil || st.Mode()&os.ModeSocket == 0 || int(st.Sys().(*syscall.Stat_t).Uid) != c.UID {
			return operatorFinishAuth{}, operatorGitHubAuth(6)
		}
	}
	return operatorFinishAuth{gh, git, repo, remote, c.Home, socket, strings.Join(dirs, ":"), helper}, nil
}

func (a operatorFinishAuth) prepare(ctx context.Context) error {
	env := []string{"HOME=" + a.Home, "PATH=" + a.Path, "GH_PROMPT_DISABLED=1", "GIT_TERMINAL_PROMPT=0"}
	gh := func(args ...string) (string, bool) {
		r, err := checkrun.RunIsolated(ctx, append([]string{a.GH}, args...), a.Home, 10*time.Second, env)
		return strings.TrimSpace(r.Stdout), err == nil && r.ExitCode == 0
	}
	if _, ok := gh("auth", "status", "--hostname", "github.com"); !ok {
		return operatorGitHubAuth(3)
	}
	if _, ok := gh("label", "list", "-R", a.Repo, "--search", "ailang-world", "--json", "name"); !ok {
		return operatorGitHubAuth(4)
	}
	permission, ok := gh("repo", "view", a.Repo, "--json", "viewerPermission", "--jq", ".viewerPermission")
	if !ok || (permission != "WRITE" && permission != "MAINTAIN" && permission != "ADMIN") {
		return operatorGitHubAuth(5)
	}

	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
	args := []string{"--git-dir=/dev/null", "-c", "credential.interactive=never", "-c", "credential.helper="}
	if strings.HasPrefix(a.Remote, "https://") {
		args = append(args, "-c", "credential.helper="+a.Helper)
	} else {
		if a.Socket != "" {
			conn, err := (&net.Dialer{Timeout: time.Second}).DialContext(ctx, "unix", a.Socket)
			if err != nil {
				return operatorGitHubAuth(6)
			}
			conn.Close()
			env = append(env, "SSH_AUTH_SOCK="+a.Socket)
		}
		args = append(args, "-c", "core.sshCommand=/usr/bin/ssh -F /dev/null -o BatchMode=yes -o StrictHostKeyChecking=yes")
	}
	out, err := worktree.RunIsolated(ctx, a.Git, append(args, "ls-remote", "--symref", "--", a.Remote, "HEAD"), env, nil)
	fields := strings.Fields(out)
	if err != nil || len(fields) != 5 || fields[0] != "ref:" || !strings.HasPrefix(fields[1], "refs/heads/") ||
		fields[2] != "HEAD" || !operatorFinishSHA.MatchString(fields[3]) || fields[4] != "HEAD" {
		return operatorGitHubAuth(6)
	}
	return nil
}

func prepareOperatorFinish(ctx context.Context, c runnerConfig, remote string) error {
	a, err := planOperatorFinishAuth(c, remote)
	if err != nil {
		return err
	}
	return a.prepare(ctx)
}
