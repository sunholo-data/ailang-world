package daemon

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surface"
)

func episodeFileEffect(effect string, result map[string]any) []seededEffect {
	return []seededEffect{{effect: effect, scope: "worktree", cost: 1, before: 50, allowed: true, result: result}}
}

func TestSurfaceEpisodeFilesProjectsRecordedActivityAndGaps(t *testing.T) {
	seed := newReportSeed(t)
	const episode = "selected"
	ctx := boundedTestContext(t)
	before := seed.readCall(episode, "first", map[string]any{"path": "host/worker.go"}, read(50, readOK), nil)
	seed.readCall("selected:other", "other", map[string]any{"path": "private-other.txt"}, read(50, readOK), nil)
	seed.readCall(episode, "second", map[string]any{"path": "host/worker.go"}, read(50, readOK), nil)
	seed.readCall(episode, "absolute", map[string]any{"path": "/private/task.txt"}, read(50, readOK), nil)
	seed.readCall(episode, "unlisted", map[string]any{"path": "repo/ghp_abcdefgh12345678"}, read(50, readOK), nil)
	seed.readCall(episode, "refused", map[string]any{"path": "refused.txt"}, read(50, map[string]any{"ok": false, "refused": "private error"}), nil)
	seed.pending("mcp", episode, "pending", "ailang-read")
	seed.input = map[string]any{"path": "host/worker.go", "content": "private source text"}
	write := seed.call("mcp", episode, "write", "ailang-write", episodeFileEffect(broker.EffectWorkspaceWrite, map[string]any{"ok": true}), nil)
	edit := seed.call("mcp", episode, "edit", "ailang-edit", episodeFileEffect(broker.EffectWorkspaceWrite, map[string]any{"ok": true}), nil)
	seed.input = nil
	produced := seed.call("a2a", episode, "produced", "workspace-exec", episodeFileEffect(broker.EffectWorkspaceExec,
		map[string]any{"worktree_delta": map[string]any{"changes": []map[string]any{{"path": "generated.txt"}}, "truncated": true, "error": "private host path"}}), nil)

	got := buildEpisodeFileActivity(ctx, seed.d.intents, seed.d.reads, episode, time.Time{})
	if got.Status != "partial" || got.Omitted == nil || got.Omitted.Unlisted != 2 || got.DeltaGaps != 2 ||
		!reflect.DeepEqual(got.Reads, []surfaceFileUse{{Path: "host/worker.go", Count: 2, FirstEntry: before}}) ||
		!reflect.DeepEqual(got.Writes, []surfaceFileUse{{Path: "host/worker.go", Op: "edit", Count: 2, FirstEntry: write, lastEntry: edit}}) ||
		!reflect.DeepEqual(got.Produced, []surfaceFileUse{{Path: "generated.txt", Count: 1, FirstEntry: produced}}) {
		t.Fatalf("unexpected file activity: %+v", got)
	}
	body := string(mustMarshal(t, got))
	for _, secret := range []string{"private-other", "/private/", "ghp_", "refused.txt", "private error", "private source text", "private host path", "lastEntry"} {
		if strings.Contains(body, secret) {
			t.Fatalf("file projection includes %q: %s", secret, body)
		}
	}

	history := buildEpisodeFileActivity(ctx, seed.d.intents, seed.d.reads, episode, time.Unix(1_800_000_000+before, 0))
	if history.Status != "complete" || len(history.Reads) != 1 || history.Reads[0].Count != 1 ||
		len(history.Writes)+len(history.Produced) != 0 {
		t.Fatalf("historical file activity: %+v", history)
	}
	empty := buildEpisodeFileActivity(ctx, seed.d.intents, seed.d.reads, "other", time.Time{})
	if empty.Status != "complete" || len(empty.Reads)+len(empty.Writes)+len(empty.Produced) != 0 {
		t.Fatalf("verified lack of file activity: %+v", empty)
	}
	if missing := buildEpisodeFileActivity(ctx, nil, nil, episode, time.Time{}); missing.Status != "unavailable" {
		t.Fatalf("unavailable file activity reported as empty: %+v", missing)
	}
}

func TestSurfaceEpisodeFilesIsBoundedPerClass(t *testing.T) {
	seed := newReportSeed(t)
	delta := make([]map[string]any, 0, surfaceEpisodeFilePaths+1)
	for i := 0; i <= surfaceEpisodeFilePaths; i++ {
		delta = append(delta, map[string]any{"path": fmt.Sprintf("generated-%03d.txt", i)})
	}
	seed.call("mcp", "selected", "exec", "workspace-exec", episodeFileEffect(broker.EffectWorkspaceExec,
		map[string]any{"worktree_delta": map[string]any{"changes": delta}}), nil)
	got := buildEpisodeFileActivity(boundedTestContext(t), seed.d.intents, seed.d.reads, "selected", time.Time{})
	if got.Status != "partial" || len(got.Produced) != surfaceEpisodeFilePaths || got.Omitted == nil || got.Omitted.Produced != 1 {
		t.Fatalf("unbounded or unreported cap: %+v", got)
	}
}

func TestSurfaceEpisodeFilesValidatesQualifiedSelectionAndLinkedMerge(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	snapshot := followSnapshotFor(t, "files", []string{"repo0", "repo1"}, nil)
	injectFollow(d, snapshot)
	start := overlayStart("wa1", d.follow.config.Repos[0].Repo, time.Now())
	mintSurfaceTaskStart(t, d, start)
	var step surface.State
	for _, s := range snapshot.Envelope.Repos[0].Index.States {
		if s.Kind == surface.KindStep && s.PR > 0 && s.Brane != "" {
			step = s
			break
		}
	}
	if step.ID == "" {
		t.Fatal("fixture lacks a recognized landed PR")
	}
	get := func(brane, at string) (int, surfaceEpisodeFiles) {
		t.Helper()
		q := url.Values{"brane": {brane}, "at": {at}}
		rec := followRequest(d, "/surface/episode-files.json?"+q.Encode())
		var out surfaceEpisodeFiles
		if rec.Code == 200 && json.Unmarshal(rec.Body.Bytes(), &out) != nil {
			t.Fatalf("invalid file response: %s", rec.Body)
		}
		return rec.Code, out
	}
	if code, out := get("repo0:brane:shared", "repo0:episode:wa1"); code != 200 || out.Episode != "wa1" || out.SourceState != "repo0:episode:wa1" || out.Status != "complete" {
		t.Fatalf("open selection: %d %+v", code, out)
	}
	if code, _ := get("repo0:brane:shared", "repo1:episode:wa1"); code != 400 {
		t.Fatalf("mixed repository accepted: %d", code)
	}
	if code, _ := get("repo0:brane:shared", "repo0:commit:"+step.Merge); code != 404 {
		t.Fatalf("unassociated merge accepted: %d", code)
	}

	snapshot.Envelope.Open.Episodes = snapshot.Envelope.Open.Episodes[1:]
	snapshot.Envelope.Open.Aliases = []surface.EpisodeAlias{{ID: "repo0:episode:wa1", Target: "repo0:commit:" + step.ID}}
	if err := surface.ValidateFollow(snapshot.Envelope); err != nil {
		t.Fatal(err)
	}
	if code, out := get("repo0:brane:shared", "repo0:commit:"+step.Merge); code != 200 || out.Episode != "wa1" || out.SourceState != "repo0:commit:"+step.ID {
		t.Fatalf("linked merge: %d %+v", code, out)
	}
}
