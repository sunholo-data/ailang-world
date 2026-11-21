package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

func TestSurfaceFileReadsOnlyPublishedLandedOutputs(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	source := d.follow.config.Repos[0]
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", source.Repo}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write := func(path, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(source.Repo, path), []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
		git("add", "--", path)
	}
	git("init", "-q", "-b", source.Branch)
	git("config", "user.name", "World Test")
	git("config", "user.email", "world@example.invalid")
	git("config", "commit.gpgSign", "false")
	write(source.Charter, "220. **surface-follow** · goal\n")
	write("deleted", "first-parent deletion\n")
	write("unchanged", "untouched\n")
	path := "new [world]\tfile"
	write(path, "first parent\n")
	git("commit", "-qm", "base")
	git("switch", "-qc", "topic")
	write(path, "side branch\n")
	git("commit", "-qm", "topic")
	step := git("rev-parse", "HEAD")
	git("switch", "-q", source.Branch)
	git("merge", "--no-ff", "--no-commit", "topic")
	write(path, "<script>text from merge</script>\n")
	write("empty", "")
	write("binary", "\x00binary")
	write("large", strings.Repeat("a", 129<<10))
	if err := os.Symlink("/etc/passwd", filepath.Join(source.Repo, "link")); err != nil {
		t.Fatal(err)
	}
	git("add", "--", "link")
	git("rm", "-q", "--", "deleted")
	git("commit", "-qm", "Merge pull request #7 from owner/topic/row220")
	merge := git("rev-parse", "HEAD")
	part, err := d.follow.collector.BuildPinned(boundedTestContext(t), source, merge)
	if err != nil {
		t.Fatal(err)
	}
	snap := &surfacefollow.Snapshot{Envelope: surface.FollowIndex{Schema: surface.FollowSchema, Generation: "g-file", Repos: []surface.FollowPartition{part}}}
	d.follow.snapshot = func() *surfacefollow.Snapshot { return snap }
	// The worktree and branch may move after the index was published.
	write(path, "next head\n")
	git("commit", "-qm", "after publication")

	endpoint := func(repo, lane, id, path string) string {
		q := url.Values{"brane": {repo + ":brane:" + lane}, "at": {repo + ":commit:" + id}, "path": {path}}
		return "/surface/file.json?" + q.Encode()
	}
	check := func(path string, status int) *httptest.ResponseRecorder {
		t.Helper()
		w := followRequest(d, path)
		if w.Code != status {
			t.Fatalf("GET %s: %d, want %d: %s", path, w.Code, status, w.Body)
		}
		return w
	}
	view := func(lane, id, path, kind, status string) surfacefollow.FilePreview {
		t.Helper()
		w := check(endpoint(source.ID, lane, id, path), http.StatusOK)
		if w.Header().Get("Content-Type") != "application/json" || w.Header().Get("X-Content-Type-Options") != "nosniff" ||
			w.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("file response headers: %+v", w.Header())
		}
		var file surfacefollow.FilePreview
		if err := json.Unmarshal(w.Body.Bytes(), &file); err != nil {
			t.Fatal(err)
		}
		if file.Kind != kind || file.Status != status || file.Path != path || file.Merge != source.ID+":commit:"+merge {
			t.Fatalf("wrong file provenance or metadata: %+v", file)
		}
		return file
	}
	for _, selection := range [][2]string{{"shared", merge}, {"row-220", step}} {
		file := view(selection[0], selection[1], path, "text", "M")
		if file.Text == nil || *file.Text != "<script>text from merge</script>\n" ||
			!strings.Contains(file.Diff, "-first parent\n+<script>text from merge</script>\n") {
			t.Fatalf("wrong merge or diff: %+v", file)
		}
	}
	deleted := view("shared", merge, "deleted", "text", "D")
	if deleted.Text == nil || *deleted.Text != "first-parent deletion\n" || deleted.Size.After != nil ||
		!strings.Contains(deleted.Diff, "-first-parent deletion\n") {
		t.Fatalf("wrong deletion preview: %+v", deleted)
	}
	if empty := view("shared", merge, "empty", "text", "A"); empty.Text == nil || *empty.Text != "" {
		t.Fatalf("empty file became unknown: %+v", empty)
	}
	for filePath, kind := range map[string]string{"binary": "binary", "link": "symlink", "large": "too-large"} {
		file := view("shared", merge, filePath, kind, "A")
		if file.Text != nil || file.Diff != "" {
			t.Fatalf("%s exposed raw content", filePath)
		}
	}
	check(endpoint(source.ID, "shared", merge, "unchanged"), http.StatusNotFound)
	check(endpoint(source.ID, "shared", merge, "../etc/passwd"), http.StatusNotFound)
	check(endpoint(source.ID, "shared", git("rev-parse", "HEAD"), path), http.StatusNotFound)
	check(endpoint(source.ID, "row-220", merge, path), http.StatusNotFound)
	check(endpoint("repo1", "shared", merge, path), http.StatusNotFound)
	check(endpoint(source.ID, "shared", merge, "\x00"), http.StatusBadRequest)
	check(endpoint(source.ID, "shared", merge, path)+"&path=other", http.StatusBadRequest)
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, endpoint(source.ID, "shared", merge, path), nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST file status %d", w.Code)
	}

	// A step with no charter brane is rendered on the shared lane.
	unassigned := part
	unassigned.Index.States = append([]surface.State(nil), part.Index.States...)
	for i := range unassigned.Index.States {
		if unassigned.Index.States[i].ID == step {
			unassigned.Index.States[i].Brane = ""
		}
	}
	unassignedSnapshot := &surfacefollow.Snapshot{Envelope: surface.FollowIndex{
		Schema: surface.FollowSchema, Generation: "g-file-shared-step", Repos: []surface.FollowPartition{unassigned}}}
	d.follow.snapshot = func() *surfacefollow.Snapshot { return unassignedSnapshot }
	view("shared", step, path, "text", "M")
	check(endpoint(source.ID, "row-220", step, path), http.StatusNotFound)
}
