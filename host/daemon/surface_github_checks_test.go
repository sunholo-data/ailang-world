package daemon

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceGitHubChecksUseExactLandedPR(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	snapshot := followSnapshotFor(t, "checks", []string{"repo0", "repo1"}, nil)
	injectFollow(d, snapshot)
	for _, config := range d.follow.config.Repos {
		for _, args := range [][]string{{"init", "--quiet"}, {"config", "--local", "remote.origin.url", "git@github.com:Owner/Repo.git"}} {
			cmd := exec.Command("git", append([]string{"-C", config.Repo}, args...)...)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("fixture git: %v: %s", err, output)
			}
		}
	}
	var step surface.State
	for _, state := range snapshot.Envelope.Repos[0].Index.States {
		if state.Kind == surface.KindStep && state.PR > 0 && state.Brane != "" {
			step = state
			break
		}
	}
	if step.ID == "" {
		t.Fatal("fixture lacks a landed PR")
	}
	dir := filepath.Join(filepath.Dir(cfg.DBPath), "surface-checks", "repo0")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	receipt := surfaceGitHubChecksReceipt{
		Schema: surfaceGitHubChecksSchema, Repo: "repo0", Slug: "owner/repo", PR: step.PR, Head: step.ID,
		CapturedAt: "2026-10-21T08:00:00Z", Checks: []surfacepage.GitHubCheck{{Name: "Build", Status: "COMPLETED", Conclusion: "SUCCESS"}},
	}
	save := func(r surfaceGitHubChecksReceipt) {
		t.Helper()
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, step.ID+".json"), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	selected := func(repo, lane, id string) *surfacepage.GitHubChecks {
		t.Helper()
		q := url.Values{"brane": {repo + ":brane:" + lane}, "at": {repo + ":commit:" + id}}
		rec := followRequest(d, "/surface/stack.json?"+q.Encode())
		var stack surfacepage.FollowStack
		if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &stack) != nil {
			t.Fatalf("stack: %d %s", rec.Code, rec.Body)
		}
		return stack.Front.GitHubChecks
	}
	for _, test := range []struct {
		name string
		edit func(*surfaceGitHubChecksReceipt)
	}{
		{"missing", nil},
		{"wrong repo", func(r *surfaceGitHubChecksReceipt) { r.Repo = "repo1" }},
		{"wrong origin", func(r *surfaceGitHubChecksReceipt) { r.Slug = "owner/private" }},
		{"wrong PR", func(r *surfaceGitHubChecksReceipt) { r.PR++ }},
		{"stale head", func(r *surfaceGitHubChecksReceipt) { r.Head = step.Merge }},
		{"unknown status", func(r *surfaceGitHubChecksReceipt) { r.Checks[0].Status = "NEW" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.edit != nil {
				bad := receipt
				bad.Checks = append([]surfacepage.GitHubCheck(nil), receipt.Checks...)
				test.edit(&bad)
				save(bad)
			}
			if got := selected("repo0", step.Brane, step.ID); got != nil {
				t.Fatalf("unverified receipt: %+v", got)
			}
		})
	}
	save(receipt)
	for _, got := range []*surfacepage.GitHubChecks{
		selected("repo0", step.Brane, step.ID),
		selected("repo0", "shared", step.Merge),
	} {
		if got == nil || got.Head != step.ID || got.PR != step.PR || got.CapturedAt != receipt.CapturedAt ||
			!strings.HasPrefix(got.URL, "https://github.com/owner/repo/pull/") || len(got.Checks) != 1 || got.Checks[0] != receipt.Checks[0] {
			t.Fatalf("landed receipt: %+v", got)
		}
	}
	if got := selected("repo1", step.Brane, step.ID); got != nil {
		t.Fatalf("receipt crossed a repo with a colliding commit: %+v", got)
	}
	for _, path := range []string{"/surface/model.json", "/surface/index.json", "/surface/follow-status.json"} {
		if rec := followRequest(d, path); rec.Code != 200 || strings.Contains(rec.Body.String(), "githubChecks") {
			t.Fatalf("%s overview receipt: %d", path, rec.Code)
		}
	}
	receipt.CapturedAt = "2026-10-21T09:00:00Z"
	receipt.Checks = []surfacepage.GitHubCheck{{Name: "Build", Status: "COMPLETED", Conclusion: "FAILURE"}}
	save(receipt)
	if got := selected("repo0", step.Brane, step.ID); got == nil || got.CapturedAt != receipt.CapturedAt || got.Checks[0].Conclusion != "FAILURE" {
		t.Fatalf("reimport receipt: %+v", got)
	}
}

func TestGitHubOriginSlugRejectsCredentialsAndOtherHosts(t *testing.T) {
	for _, origin := range []string{
		"https://x-token:secret@github.com/owner/repo",
		"https://github.com.example.invalid/owner/repo",
		"https://github.com/owner/repo?token=secret",
		"https://github.com/owner/repo\nhttps://github.com/other/repo",
		"https://github.com/owner/../../secret",
	} {
		if got := githubOriginSlug(origin); got != "" {
			t.Fatal("unverified origin accepted")
		}
	}
	for _, origin := range []string{
		"https://github.com/Owner/Repo.git",
		"git@github.com:Owner/Repo.git",
		"ssh://git@github.com/Owner/Repo.git",
	} {
		if got := githubOriginSlug(origin); got != "owner/repo" {
			t.Fatalf("configured GitHub origin slug = %q", got)
		}
	}
}
