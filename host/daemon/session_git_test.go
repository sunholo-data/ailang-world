package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func TestSessionSelfGit(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	ws := filepath.Join(root, "ws")
	dir := filepath.Join(ws, "mine")
	for _, p := range []string{repo, ws} {
		if err := os.MkdirAll(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	git := func(args ...string) string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		out, err := worktree.Output(ctx, "git", append([]string{"-C", repo}, args...))
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(out)
	}
	git("init", "-b", "dev")
	commit := func(subject string) {
		git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-m", subject)
	}
	commit("base fixture")
	base := git("rev-parse", "HEAD")
	git("worktree", "add", "-b", "world/mine", dir)
	for _, subject := range []string{"one", "two", "three"} {
		commit(subject)
	}
	seed := newReportSeed(t)
	d := seed.d
	token := strings.Repeat("a", 64)
	sum := sha256.Sum256([]byte(token))
	id := hex.EncodeToString(sum[:])
	d.resolver = authority.New(selfCredentials{row: store.SessionRow{
		EpisodeID: "mine", ExpiresAt: time.Now().Unix() + 3600, GrantsJSON: "[]",
	}, id: id})
	d.workspace = &workspaceTools{root: ws}
	request := func() (SessionSelf, string) {
		t.Helper()
		req := httptest.NewRequest("GET", "http://127.0.0.1/v1/session/self", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		d.Handler().ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("%d: %s", w.Code, w.Body.String())
		}
		var self SessionSelf
		if err := json.Unmarshal(w.Body.Bytes(), &self); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(w.Body.String(), token) || strings.Contains(w.Body.String(), id) {
			t.Fatal("credential leaked")
		}
		return self, w.Body.String()
	}
	self, raw := request()
	g := self.Git
	if g == nil || g.Error != "" || g.Base != base || g.Head != base || g.Ref != "dev" || g.Behind != 3 || g.Ahead != 0 || g.Dirty || g.Branch != "world/mine" || g.RefHead != git("rev-parse", "dev") {
		t.Fatalf("git: %s", raw)
	}
	for _, field := range []string{`"ahead":0`, `"dirty":false`, `"behind":3`} {
		if !strings.Contains(raw, field) {
			t.Fatal(raw)
		}
	}
	if !strings.Contains(g.BaseLine(), "dev is 3 commits ahead of it") {
		t.Fatal(g.BaseLine())
	}
	encoded, _ := json.Marshal(g)
	t.Logf("git=%s\n%s", encoded, g.BaseLine())
	if err := os.WriteFile(filepath.Join(dir, "dirty"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	self, _ = request()
	if !self.Git.Dirty || !strings.Contains(self.Git.BaseLine(), "worktree has uncommitted changes") {
		t.Fatal(self.Git)
	}
	// AC7: a child inheriting git's pipes must not extend the endpoint budget.
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte("#!/bin/sh\n/bin/sleep 5\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	start := time.Now()
	self, raw = request()
	if elapsed := time.Since(start); elapsed >= 3*time.Second || self.Git == nil || self.Git.Error == "" {
		t.Fatalf("elapsed=%v response=%s", elapsed, raw)
	}
	encoded, _ = json.Marshal(self.Git)
	if strings.Count(string(encoded), ":") != 1 || !strings.HasPrefix(self.Git.BaseLine(), "Base unknown:") {
		t.Fatal(string(encoded))
	}
	d.workspace.root = t.TempDir()
	self, raw = request()
	if self.Git != nil || strings.Contains(raw, `"git"`) {
		t.Fatalf("degraded: %s", raw)
	}
}
