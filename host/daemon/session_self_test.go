package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type selfCredentials struct {
	row store.SessionRow
	id  string
}

func (s selfCredentials) ResolveSession(_ context.Context, id string) (store.SessionRow, bool, error) {
	return s.row, id == s.id, nil
}

func TestSessionSelf(t *testing.T) {
	seed := newReportSeed(t)
	d := seed.d
	token := strings.Repeat("a", 64)
	sum := sha256.Sum256([]byte(token))
	id := hex.EncodeToString(sum[:])
	now := time.Now().Unix()
	row := store.SessionRow{EpisodeID: "mine", ExpiresAt: now + 3600,
		GrantsJSON: `[{"effect":"Workspace.Read","scope":"*","budget":10}]`}
	d.resolver = authority.New(selfCredentials{row: row, id: id})
	d.workspace = &workspaceTools{root: "/work"}
	entry := seed.call("mcp", "mine", "read", "ailang-read", nil, map[string]any{"ok": true})
	_, _, err := d.store.AppendNextEffectIntent(boundedTestContext(t), "mine", store.EffectIntent{
		EpisodeID: "mine", Effect: "Workspace.Read", Scope: "*", Cost: 3,
		RequestRef: hashref.SumSHA256([]byte("self-test")), LogicalTime: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "http://127.0.0.1:47330/v1/session/self", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("%d: %s", w.Code, w.Body.String())
	}
	var got SessionSelf
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Episode != "mine" || got.Worktree != "/work/mine" || got.SecondsLeft < 3598 || got.SecondsLeft > 3600 ||
		got.Entries == nil || got.Entries.First != entry || got.Grants[0].Budget != 10 ||
		got.Grants[0].Used != 3 || got.Grants[0].Left != 7 || got.Surface != got.Addr+"/surface" {
		t.Fatalf("%+v", got)
	}
	if strings.Contains(w.Body.String(), token) || strings.Contains(w.Body.String(), id) {
		t.Fatal("credential leaked")
	}
	for _, tc := range []struct {
		name, header     string
		expired, revoked bool
	}{
		{"absent", "", false, false},
		{"malformed", "Bearer oops", false, false},
		{"expired", "Bearer " + token, true, false},
		{"revoked", "Bearer " + token, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := row
			key := id
			if tc.expired {
				r.ExpiresAt = now - 1
			}
			if tc.revoked {
				key = ""
			}
			d.resolver = authority.New(selfCredentials{row: r, id: key})
			req := httptest.NewRequest("GET", "/v1/session/self", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()
			d.Handler().ServeHTTP(w, req)
			if w.Code != 401 || strings.Contains(w.Body.String(), token) || strings.Contains(w.Body.String(), id) {
				t.Fatalf("%d: %s", w.Code, w.Body.String())
			}
		})
	}
}
