package daemon

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 220 M4 route tests over an injected published snapshot: the rig index
// enrolled as repo0 and repo1 (every row, PR and hash collides), with an
// overlay built by the real BuildOverlay from facts carrying secrets.
// Identities are fakes.

func rigFollowIndex(t *testing.T) surface.Index {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	ix, err := surface.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	return ix
}

const leakyTask = "Fix login for operator@example.invalid via https://x-token:secret@github.example.invalid/acme/repo.git in /Users/operator/private/notes ghp_abcdefghijklmnopqrstuv"

// followSnapshotFor publishes the given repos' partitions (nil status means
// fresh) with an overlay: one live episode per published repo, anchored at
// the same base hash, its task text full of secrets.
func followSnapshotFor(t *testing.T, gen string, ids []string, statuses map[string]surfacefollow.RepoStatus) *surfacefollow.Snapshot {
	t.Helper()
	ix := rigFollowIndex(t)
	now := time.Date(2026, 10, 21, 9, 0, 0, 0, time.UTC)
	s := &surfacefollow.Snapshot{Envelope: surface.FollowIndex{Schema: surface.FollowSchema, Generation: gen},
		Status: surfacefollow.StatusDocument{Schema: surfacefollow.StatusSchema, Mode: "follow", Generation: gen}}
	facts := &surfacefollow.EpisodeFacts{ReadAt: now, Finishes: map[string]surfacefollow.FinishFact{}}
	for i, id := range ids {
		s.Envelope.Repos = append(s.Envelope.Repos, surface.FollowPartition{ID: id, Head: ix.Source.Head, BuiltAt: "2026-10-21T08:58:00Z", CheckedAt: "2026-10-21T08:59:00Z", Index: ix})
		ep := "w" + string(rune('a'+i)) + "1"
		facts.Sessions = append(facts.Sessions, surfacefollow.SessionFact{Episode: ep, CreatedAt: now.Unix() - 60, ExpiresAt: now.Unix() + 3600})
		facts.Starts = append(facts.Starts, surfacefollow.StartFact{Episode: ep, RepoID: id, StartedAt: now.Add(-time.Minute).Format(time.RFC3339),
			BaseSHA: ix.States[len(ix.States)-1].ID, Branch: "trial/row189-" + ep, Agent: "claude", Model: "claude-opus-5-5", TaskTitle: leakyTask, LastExpiry: now.Unix() + 3600})
	}
	s.Envelope.Open = surfacefollow.BuildOverlay(s.Envelope.Repos, facts, false)
	if err := surface.ValidateFollow(s.Envelope); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"repo0", "repo1"} {
		st, ok := statuses[id]
		if !ok {
			st = surfacefollow.RepoStatus{ID: id, State: "fresh", CheckedAt: "2026-10-21T08:59:00Z", BuiltAt: "2026-10-21T08:58:00Z"}
		}
		s.Status.Repos = append(s.Status.Repos, st)
	}
	b, err := json.Marshal(s.Envelope)
	if err != nil {
		t.Fatal(err)
	}
	s.Bytes = b
	return s
}

func injectFollow(d *Daemon, s *surfacefollow.Snapshot) {
	d.follow.snapshot = func() *surfacefollow.Snapshot { return s }
}

func decodeFollowModel(t *testing.T, w interface{ Bytes() []byte }) surfacepage.FollowModel {
	t.Helper()
	var m surfacepage.FollowModel
	if err := json.Unmarshal(w.Bytes(), &m); err != nil {
		t.Fatalf("model: %v\n%s", err, w.Bytes())
	}
	return m
}

// AC12: one repo failing leaves the other's partition visible, each with
// its own status; partial data is a 200, never a 503.
func TestSurfaceFollowPartialReposStayVisible(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	refused := surfacefollow.RepoStatus{ID: "repo1", State: "refused",
		LastFailure: &surface.FollowFailure{Class: "source", Text: "Local branch or charter unavailable; next: git -C <repo> show <branch>:<charter>"}}
	injectFollow(d, followSnapshotFor(t, "g1", []string{"repo0"}, map[string]surfacefollow.RepoStatus{"repo1": refused}))
	w := followRequest(d, "/surface/model.json")
	if w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	m := decodeFollowModel(t, w.Body)
	if len(m.Repos) != 2 || m.Repos[0].State != "fresh" || m.Repos[0].Warn || m.Repos[1].State != "refused" || !m.Repos[1].Warn || m.Repos[1].Cause != "source" ||
		!strings.Contains(m.Repos[1].Line, "next: git -C <repo> show <branch>:<charter>") {
		t.Fatalf("statuses %+v", m.Repos)
	}
	for _, l := range m.Lanes {
		if l.Repo != "repo0" {
			t.Fatalf("lane %s from the refused repo", l.Key)
		}
	}
	// The failed repo later accepted a graph and then went stale: both
	// partitions are visible and the stale one keeps its last good check.
	stale := surfacefollow.RepoStatus{ID: "repo1", State: "stale", CheckedAt: "2026-10-21T08:40:00Z", BuiltAt: "2026-10-21T08:30:00Z",
		LastFailure: &surface.FollowFailure{Class: "timeout", Text: "Git call exceeded 2 s or queued job exceeded 20 s; next: ailang-worldd service status"}}
	injectFollow(d, followSnapshotFor(t, "g2", []string{"repo0", "repo1"}, map[string]surfacefollow.RepoStatus{"repo1": stale}))
	m = decodeFollowModel(t, followRequest(d, "/surface/model.json").Body)
	repos := map[string]bool{}
	for _, l := range m.Lanes {
		repos[l.Repo] = true
	}
	if !repos["repo0"] || !repos["repo1"] || !strings.Contains(m.Repos[1].Line, "stale · last successful local check ") || !strings.Contains(m.Repos[1].Line, "08:40 · Git call exceeded") ||
		!strings.HasPrefix(m.NowEdge.Health, "World serving · repo1 stale") {
		t.Fatalf("stale repo: %+v %q", m.Repos[1], m.NowEdge.Health)
	}
	// The static pages carry the same meaning.
	page := followRequest(d, "/surface/static")
	if page.Code != 200 || !strings.Contains(page.Body.String(), `data-warn="true">repo1 · dev · stale · last successful local check `) ||
		!strings.Contains(page.Body.String(), `<p class="source">repo0 · dev · local branch checked `) {
		t.Fatalf("static: %d", page.Code)
	}
	head := rigFollowIndex(t).Source.Head
	st := followRequest(d, "/surface/static?brane="+url.QueryEscape("repo1:brane:shared")+"&at="+url.QueryEscape("repo1:commit:"+head))
	if st.Code != 200 || !strings.Contains(st.Body.String(), `<p class="status">repo1 · dev · stale`) {
		t.Fatalf("static stack: %d", st.Code)
	}
	// No accepted graph anywhere: a structured 503 naming each repo.
	injectFollow(d, &surfacefollow.Snapshot{Envelope: surface.FollowIndex{Schema: surface.FollowSchema},
		Status: surfacefollow.StatusDocument{Repos: []surfacefollow.RepoStatus{{ID: "repo0", State: "building"}, refused}}})
	for _, path := range []string{"/surface/model.json", "/surface/stack.json?brane=a&at=b"} {
		w = followRequest(d, path)
		var r surfaceReason
		if w.Code != 503 || json.Unmarshal(w.Body.Bytes(), &r) != nil || !r.Refused || len(r.Repos) != 2 || r.Repos[1].State != "refused" ||
			!strings.Contains(r.Reason, "repo0: building; repo1: refused: source") {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body)
		}
	}
	if w := followRequest(d, "/surface/static"); w.Code != 503 || !strings.Contains(w.Body.String(), "repo1 · dev · refused") {
		t.Fatalf("static refusal: %d", w.Code)
	}
	if w := followRequest(d, "/surface"); w.Code != 503 {
		t.Fatalf("shell: %d", w.Code)
	}
}

// AC13 (server side): the model and the stack each carry the generation of
// the one snapshot their request read, so the page can tell a coherent pair
// from a pair that straddles a publication.
func TestSurfaceFollowModelAndStackGenerations(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	injectFollow(d, followSnapshotFor(t, "gA", []string{"repo0", "repo1"}, nil))
	head := rigFollowIndex(t).Source.Head
	stackPath := "/surface/stack.json?brane=" + url.QueryEscape("repo1:brane:shared") + "&at=" + url.QueryEscape("repo1:commit:"+head)
	m := decodeFollowModel(t, followRequest(d, "/surface/model.json").Body)
	var st surfacepage.FollowStack
	w := followRequest(d, stackPath)
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &st) != nil || st.Generation != m.Generation || m.Generation != "gA" || st.Repo != "repo1" {
		t.Fatalf("coherent pair: %d %s %s", w.Code, m.Generation, st.Generation)
	}
	injectFollow(d, followSnapshotFor(t, "gB", []string{"repo0", "repo1"}, nil))
	w = followRequest(d, stackPath)
	if json.Unmarshal(w.Body.Bytes(), &st) != nil || st.Generation != "gB" {
		t.Fatalf("a stack read after a publication must carry the new generation: %s", w.Body)
	}
	// The open episode's stack, and its address after the episode lands.
	open := "/surface/stack.json?brane=" + url.QueryEscape("repo1:brane:shared") + "&at=" + url.QueryEscape("repo1:episode:wb1")
	w = followRequest(d, open)
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &st) != nil || st.Front.Kind != surfacepage.KindOpen || st.Front.Brane != "repo1:brane:row-189" {
		t.Fatalf("open stack: %d %s", w.Code, w.Body)
	}
}

// AC14: no follow route serializes a token, preset bytes, an e-mail
// address, a source filesystem path or a remote credential URL.
func TestSurfaceFollowRoutesLeakNothing(t *testing.T) {
	cfg, root := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	injectFollow(d, followSnapshotFor(t, "gL", []string{"repo0", "repo1"}, nil))
	head := rigFollowIndex(t).Source.Head
	q := func(brane, at string) string {
		return "?brane=" + url.QueryEscape(brane) + "&at=" + url.QueryEscape(at)
	}
	paths := []string{"/surface/model.json", "/surface/index.json", "/surface/follow-status.json", "/surface/static", "/surface",
		"/surface/stack.json" + q("repo0:brane:shared", "repo0:commit:"+head), "/surface/stack.json" + q("repo0:brane:shared", "repo0:episode:wa1"),
		"/surface/static" + q("repo0:brane:shared", "repo0:commit:"+head), "/surface/static" + q("repo0:brane:shared", "repo0:episode:wa1")}
	secrets := []string{root, cfg.SurfaceFollow, d.follow.config.Repos[0].Repo, d.follow.config.Repos[1].Repo, "operator@example.invalid", "x-token",
		"secret@", "github.example.invalid", "/Users/operator", "ghp_", "https://"}
	sawTask := false
	for _, path := range paths {
		w := followRequest(d, path)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body)
		}
		body := w.Body.String()
		for _, s := range secrets {
			if strings.Contains(body, s) {
				t.Fatalf("%s leaks %q", path, s)
			}
		}
		if strings.Contains(body, "Fix login for [email]") {
			sawTask = true
		}
	}
	if !sawTask {
		t.Fatal("the sanitized task never reached a route: the leak check would be vacuous")
	}
}
