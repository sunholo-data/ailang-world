package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const surfaceGitHubChecksSchema = "world/surface-github-checks/v1"

type surfaceGitHubChecksReceipt struct {
	Schema     string                    `json:"schema"`
	Repo       string                    `json:"repo"`
	Slug       string                    `json:"slug"`
	PR         int                       `json:"pr"`
	Head       string                    `json:"head"`
	CapturedAt string                    `json:"capturedAt"`
	Checks     []surfacepage.GitHubCheck `json:"checks"`
}

var githubSlugPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,38}/[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,99}$`)

func githubOriginSlug(origin string) string {
	for _, prefix := range []string{"https://github.com/", "git@github.com:", "ssh://git@github.com/"} {
		if slug, ok := strings.CutPrefix(strings.TrimSpace(origin), prefix); ok {
			slug = strings.TrimSuffix(slug, ".git")
			if githubSlugPattern.MatchString(slug) {
				return strings.ToLower(slug)
			}
		}
	}
	return ""
}

// readSurfaceGitHubChecks uses only the accepted follow graph, a bounded local
// receipt, and git config. The daemon never contacts GitHub.
func (d *Daemon) readSurfaceGitHubChecks(ctx context.Context, env surface.FollowIndex, st surfacepage.FollowStack) *surfacepage.GitHubChecks {
	repo, kind, selected, err := surface.ParseQualifiedID(st.Front.ID)
	if err != nil || repo != st.Repo || kind != "commit" || st.Front.PR <= 0 || d.cfg.DBPath == "" {
		return nil
	}
	head := surfacePRHead(env, repo, selected, st.Front.PR)
	if head == "" {
		return nil
	}
	var repoPath string
	for _, config := range d.follow.config.Repos {
		if config.ID == repo {
			repoPath = config.Repo
			break
		}
	}
	if repoPath == "" {
		return nil
	}

	root := openSurfaceDir(filepath.Join(filepath.Dir(d.cfg.DBPath), "surface-checks"))
	if root == nil {
		return nil
	}
	defer root.Close()
	info, err := root.Lstat(repo)
	if err != nil || !info.IsDir() {
		return nil
	}
	repoDir, err := root.OpenRoot(repo)
	if err != nil {
		return nil
	}
	defer repoDir.Close()
	opened, err := repoDir.Stat(".")
	if err != nil || !os.SameFile(info, opened) {
		return nil
	}
	data := readSurfaceFile(repoDir, head+".json", 128<<10)
	receipt := parseSurfaceGitHubChecks(data, repo, head, st.Front.PR)
	if receipt == nil {
		return nil
	}
	origin, err := (surfacefollow.ReadGit{}).Git(ctx, []string{"-C", repoPath, "config", "--local", "--get", "remote.origin.url"})
	if err != nil || githubOriginSlug(origin) != receipt.Slug {
		return nil
	}
	return &surfacepage.GitHubChecks{PR: receipt.PR, Head: head,
		URL:        "https://github.com/" + receipt.Slug + "/pull/" + strconv.Itoa(receipt.PR),
		CapturedAt: receipt.CapturedAt, Checks: receipt.Checks}
}

// A recognized PR pair shares the same number, and the landed head is the
// selected step itself or the merge's actual observed second parent.
func surfacePRHead(env surface.FollowIndex, repo, selected string, pr int) string {
	for _, part := range env.Repos {
		if part.ID != repo {
			continue
		}
		var head, merge string
		for _, s := range part.Index.States {
			if s.ID == selected && s.PR == pr {
				switch {
				case s.Kind == surface.KindStep:
					head, merge = s.ID, s.Merge
				case s.Kind == surface.KindShared && len(s.Parents) == 2:
					head, merge = s.Parents[1], s.ID
				}
				break
			}
		}
		if head == "" || merge == "" {
			return ""
		}
		for _, s := range part.Index.States {
			if s.ID != merge || s.Kind != surface.KindShared || s.PR != pr || len(s.Parents) != 2 || s.Parents[1] != head {
				continue
			}
			if _, err := surface.QualifyID(repo, "commit", head); err == nil {
				return head
			}
		}
	}
	return ""
}

func parseSurfaceGitHubChecks(data []byte, repo, head string, pr int) *surfaceGitHubChecksReceipt {
	if len(data) == 0 || !utf8.Valid(data) {
		return nil
	}
	r := new(surfaceGitHubChecksReceipt)
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if dec.Decode(r) != nil || dec.Decode(new(any)) != io.EOF ||
		r.Schema != surfaceGitHubChecksSchema || r.Repo != repo || r.Head != head || r.PR != pr ||
		!githubSlugPattern.MatchString(r.Slug) || r.Slug != strings.ToLower(r.Slug) ||
		r.Checks == nil || len(r.Checks) > 250 {
		return nil
	}
	captured, err := time.Parse(time.RFC3339, r.CapturedAt)
	if err != nil {
		return nil
	}
	r.CapturedAt = captured.UTC().Format(time.RFC3339)
	for i, c := range r.Checks {
		if !surfaceGitHubCheckState(c) {
			return nil
		}
		r.Checks[i].Name = surface.SanitizeOpenText(c.Name, 160)
		if r.Checks[i].Name == "" {
			return nil
		}
	}
	return r
}

func surfaceGitHubCheckState(c surfacepage.GitHubCheck) bool {
	switch c.Status {
	case "EXPECTED", "ERROR", "FAILURE", "PENDING", "SUCCESS", "IN_PROGRESS", "QUEUED", "REQUESTED", "WAITING":
		return c.Conclusion == ""
	case "COMPLETED":
		switch c.Conclusion {
		case "", "ACTION_REQUIRED", "CANCELLED", "FAILURE", "NEUTRAL", "SKIPPED", "STALE", "STARTUP_FAILURE", "SUCCESS", "TIMED_OUT":
			return true
		}
	}
	return false
}
