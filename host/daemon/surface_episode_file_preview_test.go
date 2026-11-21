package daemon

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

func TestSurfaceEpisodeFilePreviewsPinGitAndVerifyWriteDigest(t *testing.T) {
	seed := newReportSeed(t)
	d := seed.d
	repo, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write := func(path, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(repo, path), []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
		git("add", "--", path)
	}
	git("init", "-q", "-b", "dev")
	git("config", "user.name", "World Test")
	git("config", "user.email", "world@example.invalid")
	git("config", "commit.gpgSign", "false")
	write("read.txt", "contents at start\n")
	write("verified.txt", "recorded read\n")
	write("unlisted.txt", "unlisted but tracked\n")
	git("commit", "-qm", "base")
	base := git("rev-parse", "HEAD")
	write("write.txt", "verified final write\n")
	write("mismatch.txt", "different committed version\n")
	git("commit", "-qm", "finish")
	finish := git("rev-parse", "HEAD")
	write("read.txt", "mutable file that must stay private\n")
	write("write.txt", "mutable post-finish write\n")

	d.follow = &surfaceFollow{config: surfacefollow.Config{Repos: []surfacefollow.Repo{{ID: "repo0", Repo: repo, Branch: "dev"}}}}
	injectFollow(d, followSnapshotFor(t, "episode-file", []string{"repo0"}, nil))
	start := overlayStart("wa1", repo, time.Now())
	start.Base.SHA = base
	mintSurfaceTaskStart(t, d, start)

	seed.readCall("wa1", "read", map[string]any{"path": "read.txt", "offset": 1},
		read(50, map[string]any{"ok": true, "content": "archived private content", "total_lines": 1}), nil)
	commitEpisodeFileProof(t, seed, "file-read", "read", "verified.txt", "recorded read\n")
	commitEpisodeFileProof(t, seed, "old-write", "write", "write.txt", "private earlier write\n")
	commitEpisodeFileProof(t, seed, "write", "write", "write.txt", "verified final write\n")
	commitEpisodeFileProof(t, seed, "mismatch", "write", "mismatch.txt", "different private write\n")

	get := func(source, path string, wantCode int) surfacefollow.FilePreview {
		t.Helper()
		q := url.Values{"brane": {"repo0:brane:shared"}, "at": {"repo0:episode:wa1"}, "path": {path}, "source": {source}}
		res := followRequest(d, "/surface/file.json?"+q.Encode())
		if res.Code != wantCode {
			t.Fatalf("%s %s: HTTP %d, want %d: %s", source, path, res.Code, wantCode, res.Body)
		}
		var file surfacefollow.FilePreview
		if wantCode == 200 && json.Unmarshal(res.Body.Bytes(), &file) != nil {
			t.Fatalf("invalid file preview: %s", res.Body)
		}
		if strings.Contains(res.Body.String(), "private") || strings.Contains(res.Body.String(), "mutable") {
			t.Fatalf("response exposed noncommitted data: %s", res.Body)
		}
		return file
	}

	file := get("episode-read", "read.txt", 200)
	if file.SourceLabel != "At episode start" || file.Commit != "repo0:commit:"+base || file.Text == nil || *file.Text != "contents at start\n" {
		t.Fatalf("wrong read version: %+v", file)
	}
	file = get("episode-read", "verified.txt", 200)
	if file.SourceLabel != "Verified recorded read" || file.Text == nil || *file.Text != "recorded read\n" {
		t.Fatalf("verified read did not match immutable bytes: %+v", file)
	}
	get("episode-read", "unlisted.txt", 404)
	get("episode-read", "write.txt", 404)
	get("episode-write", "../read.txt", 404)
	get("episode-read", "read.txt&source=episode-write", 404)
	if file = get("episode-write", "write.txt", 200); file.Kind != "unavailable" || file.Text != nil {
		t.Fatalf("served a write absent a trusted finish: %+v", file)
	}
	record := EpisodeFinish{Episode: "wa1", Commit: finish, Outcome: "done"}
	putFinish(t, d, record, finishProvenance)
	file = get("episode-write", "write.txt", 200)
	if file.SourceLabel != "Verified written version" || file.Commit != "repo0:commit:"+finish ||
		file.Text == nil || *file.Text != "verified final write\n" {
		t.Fatalf("wrong matching latest write: %+v", file)
	}
	file = get("episode-write", "mismatch.txt", 200)
	if file.SourceLabel != "At episode finish" || file.Text == nil || *file.Text != "different committed version\n" {
		t.Fatalf("overclaimed unmatched write: %+v", file)
	}
}

// Install the production broker's object formats on top of the shared report
// fixture. The outer invocation args and broker request intentionally differ.
func commitEpisodeFileProof(t *testing.T, seed *reportSeed, task, op, path, content string) {
	t.Helper()
	effect := broker.EffectWorkspaceRead
	var args map[string]any
	result := map[string]any{"ok": true, "content": content}
	if op == "write" {
		effect = broker.EffectWorkspaceWrite
		args = map[string]any{"content": content}
		result = map[string]any{"ok": true}
	}
	if args == nil {
		args = map[string]any{}
	}
	args["op"], args["path"] = op, path
	seed.input = map[string]any{"path": path, "content": "outer private input that must never count"}
	c, intent := seed.build("mcp", "wa1", task, "ailang-"+op,
		[]seededEffect{{effect: effect, scope: broker.WorkspaceScope, cost: 1, before: 5, allowed: true, result: result}}, nil)
	seed.input = nil

	makeObject := func(kind string, bytes []byte) store.Object {
		obj := reportObject(kind, bytes)
		obj.Provenance = "host/broker"
		return obj
	}
	prefix := fmt.Sprintf("%d:%s%d:%s1:1800000000:", len(effect), effect, len(broker.WorkspaceScope), broker.WorkspaceScope)
	requestObj := makeObject(broker.EffectRequestV1, append([]byte(prefix), mustMarshal(t, args)...))
	resultObj := makeObject(broker.EffectResultV1, mustMarshal(t, result))
	effectObj := makeObject(broker.EffectRecordV1, broker.EncodeRecord(broker.EffectRecord{
		Effect: effect, Scope: broker.WorkspaceScope, Cost: 1, BudgetBefore: 5, BudgetAfter: 4, Allowed: true,
		RequestRef: requestObj.Hash, ResultRef: resultObj.Hash,
	}))

	var invocation map[string]any
	for _, obj := range c.Objects {
		if obj.SemanticID == coordinator.RecordV2 {
			if err := json.Unmarshal(obj.Payload, &invocation); err != nil {
				t.Fatal(err)
			}
		}
	}
	invocation["effects"] = []string{effectObj.Hash.String()}
	invocationObj := reportObject(coordinator.RecordV2, mustMarshal(t, invocation))
	intent.TransitionRef, c.Entry.TransitionRef = invocationObj.Hash, invocationObj.Hash
	c.Objects = slices.DeleteFunc(c.Objects, func(obj store.Object) bool {
		return obj.SemanticID == broker.EffectRequestV1 || obj.SemanticID == broker.EffectResultV1 || obj.SemanticID == broker.EffectRecordV1
	})
	c.Objects = append(c.Objects, requestObj, resultObj, effectObj, invocationObj)
	ctx := boundedTestContext(t)
	if _, _, err := seed.d.store.AppendIntent(ctx, intent.InvocationID, intent); err != nil {
		t.Fatal(err)
	}
	if err := seed.d.store.Commit(ctx, c); err != nil {
		t.Fatal(err)
	}
	seed.head = c.NextWorld
	if got := seed.d.episodeFileDigest(ctx, "wa1", path, op, seed.head.Revision); got != hashref.SumSHA256([]byte(content)) {
		t.Fatalf("%s did not use exact effect evidence", op)
	}
}
