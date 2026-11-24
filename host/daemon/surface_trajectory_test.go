package daemon

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceTrajectoryPrivateTranscriptAndNodeScope(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	injectFollow(d, followSnapshotFor(t, "trajectory", []string{"repo0", "repo1"}, nil))
	start := overlayStart("wa1", d.follow.config.Repos[0].Repo, time.Now())
	const secret = "raw-host-private-prompt"
	const raw = `{"type":"system","subtype":"init","cwd":"/private/work"}
{"type":"assistant","message":{"content":"` + secret + `"}}
{"type":"result","result":"ghp_fakeprivatevalue123","subtype":"success"}
`
	digest := hashref.SumSHA256([]byte(raw))
	start.Launches = []store.EpisodeLaunch{{Mode: "interactive", TranscriptState: "kept",
		Transcript: "/must/not/use/transcript.jsonl", TranscriptSHA256: digest.String()}}
	mintSurfaceTaskStart(t, d, start)
	dir := filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "wa1.transcript.jsonl")
	if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
		t.Fatal(err)
	}

	get := func(brane, at string) (int, string) {
		t.Helper()
		q := url.Values{"brane": {brane}, "at": {at}}
		res := followRequest(d, "/surface/trajectory.json?"+q.Encode())
		return res.Code, res.Body.String()
	}
	status, body := get("repo0:brane:shared", "repo0:episode:wa1")
	if status != 200 || strings.Contains(body, secret) || strings.Contains(body, "ghp_") ||
		strings.Contains(body, "/private") || strings.Contains(body, "/must/") {
		t.Fatalf("projection status %d or unsanitized response %s", status, body)
	}
	var out surfaceTrajectory
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatal(err)
	}
	if out.SourceState != "repo0:episode:wa1" || out.Episode != "wa1" ||
		!out.Harness.Available || len(out.Harness.Events) != 3 || out.Harness.Events[1].Kind != "assistant-event" ||
		!out.World.Available || len(out.World.Calls) != 0 ||
		out.RawTranscript == nil || out.RawTranscript.Filename != "wa1.transcript.jsonl" || out.RawTranscript.SHA256 != digest.Digest() {
		t.Fatalf("unexpected trajectory: %+v", out)
	}
	if code, _ := get("repo1:brane:shared", "repo1:episode:wa1"); code != 404 {
		t.Fatalf("wrong repo: status %d", code)
	}
	if code, _ := get("repo1:brane:shared", "repo0:episode:wa1"); code != 400 {
		t.Fatalf("mixed repos: status %d", code)
	}
	if err := os.WriteFile(path, []byte(raw+"tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	status, body = get("repo0:brane:shared", "repo0:episode:wa1")
	if status != 200 || strings.Contains(body, "rawTranscript") || !strings.Contains(body, `"harness":{"available":false,"events":[]}`) {
		t.Fatalf("tampered file was offered: %d %s", status, body)
	}
}

func TestSurfaceTrajectoryMergeResolvesOnlyVerifiedLinkedStep(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	snapshot := followSnapshotFor(t, "trajectory-landed", []string{"repo0"}, nil)
	injectFollow(d, snapshot)
	var step surface.State
	for _, s := range snapshot.Envelope.Repos[0].Index.States {
		if s.Kind == surface.KindStep && s.PR > 0 && s.Brane != "" {
			step = s
			break
		}
	}
	if step.ID == "" {
		t.Fatal("fixture lacks a landed PR step")
	}
	start := overlayStart("wa1", d.follow.config.Repos[0].Repo, time.Now())
	mintSurfaceTaskStart(t, d, start)
	q := url.Values{"brane": {"repo0:brane:shared"}, "at": {"repo0:commit:" + step.Merge}}
	endpoint := "/surface/trajectory.json?" + q.Encode()
	front := func() surfacepage.StackFront {
		t.Helper()
		res := followRequest(d, "/surface/stack.json?"+q.Encode())
		var st surfacepage.FollowStack
		if res.Code != 200 || json.Unmarshal(res.Body.Bytes(), &st) != nil {
			t.Fatalf("merge stack: %d %s", res.Code, res.Body)
		}
		return st.Front
	}
	if got := front(); got.TrajectoryEpisode != "" {
		t.Fatalf("merge before trusted alias: %+v", got)
	}
	if got := followRequest(d, endpoint); got.Code != 404 {
		t.Fatalf("merge without alias: %d %s", got.Code, got.Body)
	}
	snapshot.Envelope.Open.Episodes = nil
	snapshot.Envelope.Open.Aliases = []surface.EpisodeAlias{{ID: "repo0:episode:wa1", Target: "repo0:commit:" + step.ID}}
	if err := surface.ValidateFollow(snapshot.Envelope); err != nil {
		t.Fatal(err)
	}
	if got := front(); got.TrajectoryEpisode != "wa1" || got.Run != nil || got.TaskBody != "" {
		t.Fatalf("linked episode may expose trajectory only: %+v", got)
	}
	got := followRequest(d, endpoint)
	var data surfaceTrajectory
	if got.Code != 200 || json.Unmarshal(got.Body.Bytes(), &data) != nil ||
		data.Episode != "wa1" || data.SourceState != "repo0:commit:"+step.ID {
		t.Fatalf("linked merge: %d %s", got.Code, got.Body)
	}
}

func TestSurfaceTrajectoryRetainedScrubbedTranscriptMetadataOnly(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	dir := filepath.Join(filepath.Dir(cfg.DBPath), "sessions")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data := []byte(`{"world_refusal":"session token in agent stream; scrubbed"}` + "\n")
	digest := hashref.SumSHA256(data)
	if err := os.WriteFile(filepath.Join(dir, "e1.transcript.jsonl"), data, 0600); err != nil {
		t.Fatal(err)
	}
	start := &store.EpisodeStart{Episode: "e1", Agent: "codex", Launches: []store.EpisodeLaunch{
		{Mode: "interactive", TranscriptState: "scrubbed", TranscriptSHA256: digest.String(), EndedAt: "2026-10-21T09:01:00Z"},
	}}
	ctx := boundedTestContext(t)
	h, raw := d.readTrajectoryHarness(ctx, start, time.Time{})
	if h.Available || raw == nil || raw.SHA256 != digest.Digest() {
		t.Fatalf("scrubbed metadata h=%+v raw=%+v", h, raw)
	}
	h, raw = d.readTrajectoryHarness(ctx, start, time.Date(2026, 10, 21, 9, 0, 0, 0, time.UTC))
	if h.Available || raw != nil {
		t.Fatalf("later launch leaked into history h=%+v raw=%+v", h, raw)
	}
}

func TestSurfaceTrajectoryCommittedCallsStayInSelectedEpisodeHistory(t *testing.T) {
	seed := seedTwoEpisodes(t)
	ctx := boundedTestContext(t)
	got := readTrajectoryWorld(ctx, seed.d.intents, seed.d.reads, "ep-a", time.Time{})
	if !got.Available || got.Truncated || len(got.Calls) != 7 || got.Calls[0].Entry != 1 ||
		got.Calls[6].Entry != 10 || got.Calls[6].Surface != "a2a" || got.Calls[6].Tool != "ailang-edit" {
		t.Fatalf("unexpected committed evidence: %+v", got)
	}
	got = readTrajectoryWorld(ctx, seed.d.intents, seed.d.reads, "ep-a", time.Unix(1_800_000_007, 0))
	for _, call := range got.Calls {
		if call.Entry > 7 || call.Entry == 2 || call.Entry == 5 {
			t.Fatalf("later or other episode call leaked: %+v", call)
		}
	}
	if len(got.Calls) == 0 {
		t.Fatalf("no calls before cutoff: %+v", got)
	}
}

func TestSurfaceTrajectoryHarnessCapsProjectedEvents(t *testing.T) {
	const line = `{"type":"response_item","payload":{"type":"function_call","arguments":"private"}}` + "\n"
	data := []byte(strings.Repeat(line, maxTrajectoryEvents+1))
	start := &store.EpisodeStart{Launches: []store.EpisodeLaunch{{}}}
	got := projectTrajectoryHarness(boundedTestContext(t), start, time.Time{}, data, []hashref.HashRef{hashref.SumSHA256(data)})
	if !got.Available || !got.Truncated || len(got.Events) != maxTrajectoryEvents ||
		got.Events[0] != (trajectoryEvent{Kind: "tool-called", Launch: 1}) {
		t.Fatalf("harness count %d, available %v, truncated %v", len(got.Events), got.Available, got.Truncated)
	}
}

func TestSurfaceTrajectoryAgyTranscriptProjectsWithoutPrivateContent(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	injectFollow(d, followSnapshotFor(t, "trajectory-agy", []string{"repo0"}, nil))
	start := overlayStart("wa1", d.follow.config.Repos[0].Repo, time.Now())
	start.Agent = "agy"
	const secret = "agy-private-thinking-secret"
	const raw = `{"step_index":0,"type":"USER_INPUT","created_at":"2026-11-24T10:00:00Z","content":"` + secret + `"}
{"step_index":1,"type":"PLANNER_RESPONSE","created_at":"2026-11-24T10:00:01Z","thinking":"` + secret + `","tool_calls":[{"name":"call_mcp_tool"}]}
{"step_index":2,"type":"GENERIC","created_at":"2026-11-24T10:00:02Z","content":"tool output"}
{"step_index":3,"type":"PLANNER_RESPONSE","created_at":"2026-11-24T10:00:03Z","content":"done"}
`
	digest := hashref.SumSHA256([]byte(raw))
	start.Launches = []store.EpisodeLaunch{{Mode: "interactive", TranscriptState: "kept",
		Transcript: "/must/not/leak/transcript.jsonl", TranscriptSHA256: digest.String()}}
	mintSurfaceTaskStart(t, d, start)
	dir := filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "wa1.transcript.jsonl"), []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	q := url.Values{"brane": {"repo0:brane:shared"}, "at": {"repo0:episode:wa1"}}
	res := followRequest(d, "/surface/trajectory.json?"+q.Encode())
	if res.Code != 200 || strings.Contains(res.Body.String(), secret) || strings.Contains(res.Body.String(), "/must/") {
		t.Fatalf("status %d or unsanitized body: %s", res.Code, res.Body.String())
	}
	var out surfaceTrajectory
	if err := json.Unmarshal(res.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	wantKinds := []string{"turn-started", "tool-called", "tool-result", "assistant-event"}
	if !out.Harness.Available || len(out.Harness.Events) != len(wantKinds) {
		t.Fatalf("unexpected agy trajectory: %+v", out.Harness)
	}
	for i, want := range wantKinds {
		if out.Harness.Events[i].Kind != want || out.Harness.Events[i].At == "" {
			t.Errorf("event[%d] = %+v, want kind %q with timestamp", i, out.Harness.Events[i], want)
		}
	}
}
