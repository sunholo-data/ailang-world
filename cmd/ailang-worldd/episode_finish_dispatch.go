package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const finishLabelDescription = "Changes made by an agent working through AILANG World (provenance recorded)"

// finishDispatch is the only outward dispatch seam. Prepare resolves gh once
// and queries the label; the other methods are called only after approval.
type finishDispatch interface {
	Prepare(context.Context, string, string) (bool, error)
	RemoteHead(context.Context, string, string, string, string) (string, error)
	Ancestor(context.Context, string, string, string, string) (bool, error)
	Push(context.Context, string, string, string, string, string) error
	CreateLabel(context.Context, string, string) error
	ViewPR(context.Context, string, string, string) (daemon.FinishPR, bool, error)
	CreatePR(context.Context, string, string, string, string, string, string, bool) (daemon.FinishPR, error)
}

var finishDispatchFactory = func() finishDispatch { return &productionFinishDispatch{} }
var finishBeforeDispatch func()

type productionFinishDispatch struct {
	gh string
}

func (p *productionFinishDispatch) ghRun(ctx context.Context, wt string, args ...string) (string, error) {
	result, err := checkrun.Run(ctx, append([]string{p.gh}, args...), wt, 30*time.Second)
	if err != nil {
		return "", err
	}
	if result.ExitCode != 0 {
		return "", fmt.Errorf("gh failed (rc %d): %s", result.ExitCode, result.Stderr)
	}
	return strings.TrimSpace(result.Stdout), nil
}

func (p *productionFinishDispatch) Prepare(ctx context.Context, wt, repo string) (bool, error) {
	var err error
	p.gh, err = exec.LookPath("gh")
	if err != nil {
		return false, err
	}
	out, err := p.ghRun(ctx, wt, "label", "list", "-R", repo, "--search", "ailang-world", "--json", "name", "--jq", `any(.[]; .name == "ailang-world")`)
	if err != nil {
		return false, err
	}
	if out != "true" && out != "false" {
		return false, fmt.Errorf("invalid label lookup result")
	}
	return out == "true", nil
}

func (p *productionFinishDispatch) RemoteHead(ctx context.Context, git, wt, remote, branch string) (string, error) {
	out, err := worktree.Output(ctx, git, []string{"-C", wt, "ls-remote", "--heads", remote, "refs/heads/" + branch})
	if err != nil {
		return "", err
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", nil
	}
	if len(fields) != 2 || fields[1] != "refs/heads/"+branch || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(fields[0]) {
		return "", fmt.Errorf("invalid ls-remote result")
	}
	return fields[0], nil
}

func (p *productionFinishDispatch) Ancestor(ctx context.Context, git, wt, older, sha string) (bool, error) {
	_, err := worktree.Output(ctx, git, []string{"-C", wt, "merge-base", "--is-ancestor", older, sha})
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return err == nil, err
}

func (p *productionFinishDispatch) Push(ctx context.Context, git, wt, remote, sha, branch string) error {
	return worktree.Push(ctx, git, wt, remote, sha, branch)
}

func (p *productionFinishDispatch) CreateLabel(ctx context.Context, wt, repo string) error {
	_, err := p.ghRun(ctx, wt, "label", "create", "ailang-world", "-R", repo, "--color", "e73c17", "--description", finishLabelDescription)
	return err
}

func (p *productionFinishDispatch) ViewPR(ctx context.Context, wt, repo, branch string) (daemon.FinishPR, bool, error) {
	out, err := p.ghRun(ctx, wt, "pr", "view", branch, "-R", repo, "--json", "url,number,isDraft,state", "--jq", "{url,number,isDraft,state}")
	if err != nil {
		if strings.Contains(err.Error(), "no pull requests found") {
			return daemon.FinishPR{}, false, nil
		}
		return daemon.FinishPR{}, false, err
	}
	if out == "null" {
		return daemon.FinishPR{}, false, nil
	}
	var wire struct {
		URL     string `json:"url"`
		Number  int    `json:"number"`
		IsDraft bool   `json:"isDraft"`
		State   string `json:"state"`
	}
	if json.Unmarshal([]byte(out), &wire) != nil || wire.URL == "" || wire.Number <= 0 {
		return daemon.FinishPR{}, false, fmt.Errorf("invalid PR lookup result")
	}
	return daemon.FinishPR{URL: wire.URL, Number: wire.Number, Draft: wire.IsDraft}, wire.State == "OPEN", nil
}

func (p *productionFinishDispatch) CreatePR(ctx context.Context, wt, repo, base, branch, title, body string, draft bool) (daemon.FinishPR, error) {
	args := []string{"pr", "create", "-R", repo, "--base", base, "--head", branch, "--title", title, "--body-file", body, "--label", "ailang-world"}
	if draft {
		args = append(args, "--draft")
	}
	// gh pr create and label create have no --json flag. Create prints one URL.
	out, err := p.ghRun(ctx, wt, args...)
	if err != nil {
		return daemon.FinishPR{}, err
	}
	prefix := "https://github.com/" + repo + "/pull/"
	number, parseErr := strconv.Atoi(strings.TrimPrefix(out, prefix))
	if !strings.HasPrefix(out, prefix) || parseErr != nil || number <= 0 {
		return daemon.FinishPR{}, fmt.Errorf("invalid PR create result")
	}
	return daemon.FinishPR{URL: out, Number: number, Draft: draft}, nil
}

var finishGitHubRemote = regexp.MustCompile(`^(https://github\.com/|git@github\.com:|ssh://git@github\.com/)([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)$`)

func finishRepository(remote string) (string, error) {
	m := finishGitHubRemote.FindStringSubmatch(remote)
	if m == nil {
		return "", fmt.Errorf("finish opens PRs on github.com only")
	}
	repo := strings.TrimSuffix(m[3], ".git")
	if repo == "" || repo == "." || repo == ".." || m[2] == "." || m[2] == ".." {
		return "", fmt.Errorf("invalid GitHub repository")
	}
	return m[2] + "/" + repo, nil
}

func finishDispatchOutcome(ctx context.Context, seam finishDispatch, git, wt, remote, repo, branch, sha, base, title, body string, draft, createLabel bool) (string, daemon.FinishPR, error) {
	class := "push-failed"
	fail := func(err error) (string, daemon.FinishPR, error) {
		var timeout interface{ Timeout() bool }
		if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &timeout) && timeout.Timeout()) {
			class = "ambiguous"
		}
		return class, daemon.FinishPR{}, err
	}
	head, err := seam.RemoteHead(ctx, git, wt, remote, branch)
	if err != nil {
		return fail(err)
	}
	if head != "" && head != sha {
		ancestor, err := seam.Ancestor(ctx, git, wt, head, sha)
		if err != nil {
			return fail(err)
		}
		if !ancestor {
			return fail(fmt.Errorf("%s on %s is at %s, not an ancestor (F7)", branch, remote, head))
		}
	}
	if head != sha {
		if err := seam.Push(ctx, git, wt, remote, sha, branch); err != nil {
			return fail(err)
		}
	}
	class = "pr-failed"
	if createLabel {
		if err := seam.CreateLabel(ctx, wt, repo); err != nil {
			return fail(err)
		}
	}
	pr, found, err := seam.ViewPR(ctx, wt, repo, branch)
	if err != nil {
		return fail(err)
	}
	if !found {
		pr, err = seam.CreatePR(ctx, wt, repo, base, branch, title, body, draft)
		if err != nil {
			return fail(err)
		}
	}
	return "done", pr, nil
}
