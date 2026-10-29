package daemon

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// AC7: every route that accepts a session credential is exercised with the
// minted token. None can update the attended row, even with a forged start.
func TestEpisodeStartReadProjectionsAndSessionRoutesCannotWrite(t *testing.T) {
	seed := newReportSeed(t)
	d := seed.d
	path := d.credentialsPath
	c := testCreds(t, d)
	start := &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: "mine", Base: store.StartBase{Branch: "dev", SHA: strings.Repeat("a", 40)},
		Branch: "world/mine", Preset: store.StartPreset{Name: "p", SHA256: "sha256:p", Bytes: "{}"}, Agent: "codex", Model: "fixture"}
	now := time.Now().Unix()
	token, _, _, err := authority.MintStarted(boundedTestContext(t), c, "mine", authority.ProfileBinding{},
		[]broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 1, ExpiresAt: now + 3600}}, 3600, now, nil, start)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ method, path, body string }{
		{"GET", "/.well-known/agent.json", ""},
		{"GET", "/v1/session/self", ""},
		{"GET", "/v1/episodes/mine/report", ""},
		{"POST", "/v1/commit", `{"start":{"base":"forged"}}`},
		{"POST", "/mcp/", `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`},
		{"POST", "/a2a/", `{"jsonrpc":"2.0","id":1,"method":"tasks/get","params":{"id":"absent"}}`},
		{"POST", "/agui/", `{"start":{"base":"forged"}}`},
	} {
		t.Run(tc.path, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			r.Header.Set("Authorization", "Bearer "+token)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept", "application/json, text/event-stream")
			w := httptest.NewRecorder()
			d.Handler().ServeHTTP(w, r)
			if strings.Contains(w.Body.String(), token) {
				t.Fatal("token leaked")
			}
			if w.Code == 401 {
				t.Fatalf("session token refused: %s", w.Body.String())
			}
			if tc.path == "/v1/session/self" || tc.path == "/v1/episodes/mine/report" {
				if w.Code != 200 {
					t.Fatalf("%d %s", w.Code, w.Body.String())
				}
				var obj map[string]json.RawMessage
				if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
					t.Fatal(err)
				}
				if len(obj["start"]) == 0 || !bytes.Contains(obj["start"], []byte("dev")) {
					t.Fatal(w.Body.String())
				}
			}
			ro, err := store.OpenCredentialsReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := ro.EpisodeStart(boundedTestContext(t), "mine")
			_ = ro.Close()
			if err != nil || !reflect.DeepEqual(start, got) {
				t.Fatalf("session route changed start: %+v %v", got, err)
			}
		})
	}
	// Ordinary mints expose no start field.
	plain, _, _, err := authority.Mint(boundedTestContext(t), c, "plain", []broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 1}}, 3600, now, nil)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/v1/session/self", nil)
	r.Header.Set("Authorization", "Bearer "+plain)
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, r)
	if w.Code != 200 || strings.Contains(w.Body.String(), "\"start\"") {
		t.Fatal(w.Body.String())
	}
	// A report with calls but no launch preset also omits start.
	seed.call("mcp", "plain", "read", "ailang-read", nil, map[string]any{"ok": true})
	w = httptest.NewRecorder()
	d.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/v1/episodes/plain/report", nil))
	if w.Code != 200 || strings.Contains(w.Body.String(), "\"start\"") {
		t.Fatal(w.Body.String())
	}
}
