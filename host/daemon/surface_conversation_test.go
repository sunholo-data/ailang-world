package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestSurfaceConversationVerifiedAndLiveAgy(t *testing.T) {
	cfg, _ := followFixture(t, 1)
	// Place DB inside a temp <home>/.ailang/world/world.db layout so live agy path jailing can be tested cleanly.
	fakeHome, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	worldDir := filepath.Join(fakeHome, ".ailang", "world")
	if err := os.MkdirAll(worldDir, 0o700); err != nil {
		t.Fatal(err)
	}
	cfg.DBPath = filepath.Join(worldDir, "world.db")
	cfg.SurfaceOrigins = []string{"http://localhost:7644"}

	d := newFollowDaemon(t, cfg)
	injectFollow(d, followSnapshotFor(t, "conversation", []string{"repo0"}, nil))

	// 1. Live agy session (no closed launches yet; reads <launch>/.agents/.world-session-meta.json -> brain transcript.jsonl).
	liveEp := "wa1"
	startLive := overlayStart(liveEp, d.follow.config.Repos[0].Repo, time.Now())
	startLive.Agent = "agy"
	mintSurfaceTaskStart(t, d, startLive)

	brainLogDir := filepath.Join(fakeHome, ".gemini", "antigravity-cli", "brain", "conv-live", ".system_generated", "logs")
	if err := os.MkdirAll(brainLogDir, 0o700); err != nil {
		t.Fatal(err)
	}
	transcriptFile := filepath.Join(brainLogDir, "transcript.jsonl")
	liveRaw := `{"step_index":0,"type":"USER_INPUT","created_at":"2026-11-26T16:00:00Z","content":"<USER_REQUEST>\nshow me the live transcript\n</USER_REQUEST>\n<ADDITIONAL_METADATA>ignored</ADDITIONAL_METADATA>"}
{"step_index":1,"type":"PLANNER_RESPONSE","created_at":"2026-11-26T16:00:02Z","content":"","tool_calls":[{"name":"call_mcp_tool","arguments":{"ToolName":"ailang-read"}}]}
{"step_index":2,"type":"PLANNER_RESPONSE","created_at":"2026-11-26T16:00:05Z","content":"Here is the live stream! Token ghp_1234567890abcdefghijklmno is scrubbed."}
`
	if err := os.WriteFile(transcriptFile, []byte(liveRaw), 0o600); err != nil {
		t.Fatal(err)
	}
	launchDir := filepath.Join(worldDir, "launch", liveEp, ".agents")
	if err := os.MkdirAll(launchDir, 0o700); err != nil {
		t.Fatal(err)
	}
	metaJSON, _ := json.Marshal(map[string]string{
		"conversationId": "conv-live",
		"transcriptPath": transcriptFile,
	})
	if err := os.WriteFile(filepath.Join(launchDir, ".world-session-meta.json"), metaJSON, 0o600); err != nil {
		t.Fatal(err)
	}

	mux := d.Handler()
	q := url.Values{"brane": {"repo0:brane:shared"}, "at": {"repo0:episode:" + liveEp}}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:7644/surface/conversation.json?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("live conversation status = %d: %s", rec.Code, rec.Body.String())
	}
	var out surfaceConversation
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Status != "live" || out.Total != 2 || len(out.Messages) != 2 {
		t.Fatalf("unexpected live conversation: %+v", out)
	}
	if out.Messages[0].Role != "user" || out.Messages[0].Text != "show me the live transcript" {
		t.Fatalf("unexpected user turn: %+v", out.Messages[0])
	}
	if out.Messages[1].Role != "assistant" || !strings.Contains(out.Messages[1].Text, "[redacted]") || strings.Contains(out.Messages[1].Text, "ghp_") || len(out.Messages[1].Tools) != 1 || out.Messages[1].Tools[0] != "ailang-read" {
		t.Fatalf("unexpected assistant turn: %+v", out.Messages[1])
	}

	// Foreign Host header is refused by approvalRoute (421).
	badReq := httptest.NewRequest(http.MethodGet, "http://evil.example/surface/conversation.json?"+q.Encode(), nil)
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusMisdirectedRequest {
		t.Fatalf("foreign Host status = %d, want 421", badRec.Code)
	}

	// 2. Sealed verified session transcript on a separate daemon fixture.
	cfg2, _ := followFixture(t, 1)
	cfg2.SurfaceOrigins = []string{"http://localhost:7644"}
	d2 := newFollowDaemon(t, cfg2)
	injectFollow(d2, followSnapshotFor(t, "conversation-sealed", []string{"repo0"}, nil))
	sealedEp := "wa1"
	claudeRaw := `{"type":"user","timestamp":"2026-11-26T15:00:00Z","message":{"role":"user","content":"Fix the header layout"}}
{"type":"assistant","timestamp":"2026-11-26T15:00:04Z","message":{"role":"assistant","content":[{"type":"tool_use","name":"mcp__world__ailang-edit"},{"type":"text","text":"Updated app.css header tokens."}]}}
`
	digest := hashref.SumSHA256([]byte(claudeRaw))
	startSealed := overlayStart(sealedEp, d2.follow.config.Repos[0].Repo, time.Now())
	startSealed.Agent = "claude"
	startSealed.Launches = []store.EpisodeLaunch{{
		Mode:             "interactive",
		TranscriptState:  "kept",
		TranscriptSHA256: digest.String(),
	}}
	mintSurfaceTaskStart(t, d2, startSealed)
	sessionsDir := filepath.Join(filepath.Dir(d2.cfg.DBPath), "sessions")
	if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionsDir, sealedEp+".transcript.jsonl"), []byte(claudeRaw), 0o600); err != nil {
		t.Fatal(err)
	}

	q2 := url.Values{"brane": {"repo0:brane:shared"}, "at": {"repo0:episode:" + sealedEp}}
	req2 := httptest.NewRequest(http.MethodGet, "http://localhost:7644/surface/conversation.json?"+q2.Encode(), nil)
	rec2 := httptest.NewRecorder()
	d2.Handler().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("sealed conversation status = %d: %s", rec2.Code, rec2.Body.String())
	}
	var sealedOut surfaceConversation
	if err := json.Unmarshal(rec2.Body.Bytes(), &sealedOut); err != nil {
		t.Fatal(err)
	}
	if sealedOut.Status != "verified" || sealedOut.SHA256 != digest.Digest() || len(sealedOut.Messages) != 2 {
		t.Fatalf("unexpected sealed conversation: %+v", sealedOut)
	}
	if sealedOut.Messages[1].Tools[0] != "ailang-edit" || sealedOut.Messages[1].Text != "Updated app.css header tokens." {
		t.Fatalf("unexpected claude assistant turn: %+v", sealedOut.Messages[1])
	}
}
