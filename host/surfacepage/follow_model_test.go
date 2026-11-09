package surfacepage

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Row 220 M4 tests: the follow model, stack and static pages. The fixture is
// the rig index enrolled twice, as "world" and "tools", so every row, PR
// number and commit hash collides across the two repos (AC9). Identities are
// fakes.

func stateWithPrefix(t *testing.T, ix surface.Index, prefix string) surface.State {
	t.Helper()
	for _, s := range ix.States {
		if strings.HasPrefix(s.ID, prefix) {
			return s
		}
	}
	t.Fatalf("no state %s", prefix)
	return surface.State{}
}

// followEnv is a validated two-repo envelope with an open overlay: a live
// episode anchored in world, a PR-open one whose base is outside the window,
// an ended one on tools' episode lane, an alias with a live session, and a
// live session without a start.
func followEnv(t *testing.T) surface.FollowIndex {
	t.Helper()
	ix := rigIndex(t)
	part := func(id string) surface.FollowPartition {
		return surface.FollowPartition{ID: id, Head: ix.Source.Head, BuiltAt: "2026-10-21T08:58:00Z", CheckedAt: "2026-10-21T08:59:00Z", Index: ix}
	}
	base := stateWithPrefix(t, ix, "0ffbb9c").ID
	landed := stateWithPrefix(t, ix, "f36428a").ID
	env := surface.FollowIndex{Schema: surface.FollowSchema, Generation: "g0123456789abcdef", Repos: []surface.FollowPartition{part("world"), part("tools")},
		Open: &surface.FollowOpen{Status: surface.OpenComplete, CheckedAt: "2026-10-21T08:59:30Z", Cap: surface.MaxOpenEpisodes,
			Episodes: []surface.OpenEpisode{
				{ID: "world:episode:w189z", Repo: "world", Episode: "w189z", State: surface.OpenRunning, Label: surface.OpenLabel(surface.OpenRunning), SessionLive: true,
					StartedAt: "2026-10-21T08:30:00Z", Base: base, Anchor: "world:commit:" + base, Branch: "trial/row189-m7", Task: "Polish the plugin docs", Agent: "claude", Model: "claude-opus-5-5"},
				{ID: "world:episode:w190q", Repo: "world", Episode: "w190q", State: surface.OpenPR, Label: surface.OpenLabel(surface.OpenPR), PR: 340,
					StartedAt: "2026-10-21T07:00:00Z", EndedAt: "2026-10-21T08:00:00Z", Base: strings.Repeat("a", 40), BaseOutsideWindow: true, Branch: "trial/row190-fix", Agent: "codex", Model: "gpt-5-codex"},
				{ID: "tools:episode:t1", Repo: "tools", Episode: "t1", State: surface.OpenEnded, Label: surface.OpenLabel(surface.OpenEnded),
					StartedAt: "2026-10-20T07:00:00Z", EndedAt: "2026-10-20T09:00:00Z", Base: base, Anchor: "tools:commit:" + base, Branch: "scratch/tidy", Agent: "claude", Model: "claude-opus-5-5"},
			},
			Aliases: []surface.EpisodeAlias{{ID: "tools:episode:e1", Target: "tools:commit:" + landed, SessionLive: true}},
			Running: []surface.RunningEntry{{Episode: "w999", Reason: surface.ReasonNoStart}},
		}}
	if err := surface.ValidateFollow(env); err != nil {
		t.Fatalf("fixture envelope: %v", err)
	}
	return env
}

func followSources() []FollowSource {
	return []FollowSource{
		{ID: "world", Branch: "dev", State: SourceFresh, CheckedAt: "2026-10-21T08:59:00Z", BuiltAt: "2026-10-21T08:58:00Z"},
		{ID: "tools", Branch: "dev", State: SourceStale, CheckedAt: "2026-10-21T08:40:00Z", BuiltAt: "2026-10-21T08:30:00Z",
			Failure: &surface.FollowFailure{Class: "timeout", Text: "Git call exceeded 2 s or queued job exceeded 20 s; next: ailang-worldd service status"}},
	}
}

func repoOf(id string) string { r, _, _ := strings.Cut(id, ":"); return r }

// AC9: equal rows, PR numbers and hashes in two repos are distinct states,
// lanes and links; no link joins two repos.
func TestFollowModelKeepsCollidingReposDistinct(t *testing.T) {
	env := followEnv(t)
	m := BuildFollowModel(env, followSources(), Now{At: fixedNow})
	v1 := BuildModel(env.Repos[0].Index, Now{At: fixedNow})
	if m.Schema != FollowModelSchema || m.Generation != env.Generation {
		t.Fatalf("schema/generation: %s %s", m.Schema, m.Generation)
	}
	if len(m.Repos) != 2 || m.Repos[0].ID != "world" || m.Repos[1].ID != "tools" || m.Repos[0].Shared != "world:brane:shared" {
		t.Fatalf("repos not in config order: %+v", m.Repos)
	}
	ids, keys := map[string]bool{}, map[string]bool{}
	states, ranks := 0, map[int]bool{}
	for _, l := range m.Lanes {
		if keys[l.Key] || repoOf(l.Key) != l.Repo || !strings.HasPrefix(l.Key, l.Repo+":brane:") {
			t.Fatalf("lane %q (repo %q) is not uniquely qualified", l.Key, l.Repo)
		}
		keys[l.Key] = true
		for _, s := range l.States {
			if ids[s.ID] || repoOf(s.ID) != l.Repo {
				t.Fatalf("state %q on %q is not uniquely qualified", s.ID, l.Key)
			}
			ids[s.ID], ranks[s.Rank] = true, true
			states++
		}
	}
	v1States := 0
	for _, l := range v1.Lanes {
		v1States += len(l.States)
	}
	if states != 2*v1States+len(env.Open.Episodes) || len(ranks) != states {
		t.Fatalf("states %d (want %d), unique ranks %d", states, 2*v1States+len(env.Open.Episodes), len(ranks))
	}
	if !keys["world:brane:row-189"] || !keys["tools:brane:row-189"] {
		t.Fatal("row 189 is not a lane in both repos")
	}
	for _, k := range m.Links {
		if repoOf(k.From) != repoOf(k.To) || !ids[k.From] || !ids[k.To] {
			t.Fatalf("link %s → %s joins repos or names no state", k.From, k.To)
		}
	}
	if len(m.Links) != 2*len(v1.Links)+2 { // two anchored open states
		t.Fatalf("links %d, want %d", len(m.Links), 2*len(v1.Links)+2)
	}
	// Each repo keeps v1's order on every lane: global ranks only interleave.
	for _, l := range v1.Lanes {
		for _, repo := range []string{"world", "tools"} {
			var got []string
			for _, ml := range m.Lanes {
				if ml.Key == repo+":brane:"+l.Key {
					for _, s := range ml.States {
						if s.Kind != KindOpen {
							got = append(got, strings.TrimPrefix(s.ID, repo+":commit:"))
						}
					}
				}
			}
			var want []string
			for _, s := range l.States {
				want = append(want, s.ID)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("%s %s order differs from v1", repo, l.Key)
			}
		}
	}
	if m.NowEdge.Today.Landed != 2*v1.NowEdge.Today.Landed {
		t.Fatalf("today %d, want %d", m.NowEdge.Today.Landed, 2*v1.NowEdge.Today.Landed)
	}
}

func TestFollowModelOpenStates(t *testing.T) {
	env := followEnv(t)
	m := BuildFollowModel(env, followSources(), Now{At: fixedNow})
	find := func(id string) (ModelLane, ModelState) {
		for _, l := range m.Lanes {
			for _, s := range l.States {
				if s.ID == id {
					return l, s
				}
			}
		}
		t.Fatalf("no state %s", id)
		return ModelLane{}, ModelState{}
	}
	l, s := find("world:episode:w189z")
	if l.Key != "world:brane:row-189" || s.Kind != KindOpen || s.Open == nil || s.Open.Label != "session live" || !s.Open.SessionLive || s.Title != "Polish the plugin docs" {
		t.Fatalf("running episode on %s: %+v", l.Key, s)
	}
	anchored := false
	for _, k := range m.Links {
		if k.To == s.ID && k.From == env.Open.Episodes[0].Anchor {
			anchored = true
		}
	}
	if !anchored {
		t.Fatal("the running episode is not linked to its base")
	}
	l, s = find("world:episode:w190q")
	if l.Key != "world:brane:row-190" || s.PR != 340 || !s.Open.BaseOutsideWindow || s.Open.Label != "PR open" {
		t.Fatalf("PR-open episode: %s %+v", l.Key, s.Open)
	}
	l, s = find("tools:episode:t1")
	if l.Key != "tools:brane:open" || l.Name != "Open episodes" || l.Repo != "tools" || s.Open.Label != "session ended; outcome unknown" {
		t.Fatalf("branch without a row: %s %q %+v", l.Key, l.Name, s.Open)
	}
	// The added lane stays inside its repo's group.
	last := map[string]int{}
	for i, ml := range m.Lanes {
		if prev, ok := last[ml.Repo]; ok && prev != i-1 {
			t.Fatalf("lanes of %s are not contiguous", ml.Repo)
		}
		last[ml.Repo] = i
	}
	if len(m.Aliases) != 1 || m.Aliases[0].ID != "tools:episode:e1" {
		t.Fatalf("aliases %+v", m.Aliases)
	}
	_, landed := find(m.Aliases[0].Target)
	if !landed.SessionLive {
		t.Fatal("the landed step does not carry its live session")
	}
	o := m.NowEdge.Open
	if o == nil || o.Status != surface.OpenComplete || o.PROpen != 1 || o.Ended != 1 || len(o.Live) != 3 {
		t.Fatalf("running summary %+v", o)
	}
	if o.Live[0].ID != "world:episode:w189z" || o.Live[1].ID != m.Aliases[0].Target || o.Live[2].Reason != surface.ReasonNoStart {
		t.Fatalf("live entries %+v", o.Live)
	}
	// No overlay read yet is "not read yet", never zero running.
	env.Open = nil
	m = BuildFollowModel(env, followSources(), Now{At: fixedNow})
	if m.NowEdge.Open.Status != surface.OpenUnavailable || m.NowEdge.Open.Note == "" {
		t.Fatalf("missing overlay: %+v", m.NowEdge.Open)
	}
}

// AC12: one repo failing keeps the other's partition and both statuses.
func TestFollowModelPartialRepos(t *testing.T) {
	env := followEnv(t)
	env.Repos, env.Open = env.Repos[:1], nil
	sources := followSources()
	sources[1] = FollowSource{ID: "tools", Branch: "main", State: SourceRefused,
		Failure: &surface.FollowFailure{Class: "source", Text: "Local branch or charter unavailable; next: git -C <repo> show <branch>:<charter>"}}
	m := BuildFollowModel(env, sources, Now{At: fixedNow})
	if len(m.Repos) != 2 || m.Repos[1].Shared != "" || !m.Repos[1].Warn || !strings.Contains(m.Repos[1].Line, "tools · main · refused · Local branch or charter unavailable; next: git -C") {
		t.Fatalf("refused repo: %+v", m.Repos[1])
	}
	if m.Repos[0].Warn || m.Repos[0].Line != "world · dev · local branch checked 08:59" {
		t.Fatalf("fresh repo: %+v", m.Repos[0])
	}
	for _, l := range m.Lanes {
		if l.Repo != "world" {
			t.Fatalf("lane %s from a repo with no graph", l.Key)
		}
	}
	if m.NowEdge.Health != "World serving · tools refused" {
		t.Fatalf("health %q", m.NowEdge.Health)
	}
}

func TestFollowSourceLines(t *testing.T) {
	fail := &surface.FollowFailure{Class: "timeout", Text: "Git call exceeded 2 s or queued job exceeded 20 s; next: ailang-worldd service status"}
	for _, c := range []struct {
		src  FollowSource
		want string
		warn bool
	}{
		{FollowSource{ID: "world", Branch: "dev", State: SourceFresh, CheckedAt: "2026-10-21T08:59:00Z"}, "world · dev · local branch checked 08:59", false},
		{FollowSource{ID: "world", Branch: "dev", State: SourceUpdating, BuiltAt: "2026-10-21T08:58:00Z"}, "world · dev · updating · showing the graph built 08:58", false},
		{FollowSource{ID: "world", Branch: "dev", State: SourceStale, CheckedAt: "2026-10-20T08:40:00Z", Failure: fail},
			"world · dev · stale · last successful local check 20 Oct 08:40 · " + fail.Text, true},
		{FollowSource{ID: "world", Branch: "dev", State: SourceStale, Failure: fail}, "world · dev · stale · no successful local check yet · " + fail.Text, true},
		{FollowSource{ID: "world", Branch: "dev", State: SourceBuilding}, "world · dev · building · no graph accepted yet", false},
		{FollowSource{ID: "world", Branch: "dev", State: SourceRefused, Failure: fail}, "world · dev · refused · " + fail.Text, true},
	} {
		line, warn := SourceLine(c.src, fixedNow)
		if line != c.want || warn != c.warn {
			t.Errorf("%s: %q %v, want %q %v", c.src.State, line, warn, c.want, c.warn)
		}
	}
	if h := followHealth([]ModelRepo{{ID: "a", State: SourceFresh}, {ID: "b", State: SourceUpdating}}); h != "World healthy" {
		t.Errorf("updating is healthy: %q", h)
	}
	if h := followHealth([]ModelRepo{{ID: "a", State: SourceBuilding}, {ID: "b", State: SourceStale, Warn: true}}); h != "World serving · a building, b stale" {
		t.Errorf("health %q", h)
	}
}

func TestFollowStackResolution(t *testing.T) {
	env := followEnv(t)
	ix := env.Repos[0].Index
	step := stateWithPrefix(t, ix, "f36428a")
	// Qualified: one repo's stack, every id in it qualified with that repo.
	s, err := BuildFollowStack(env, "tools:brane:row-189", "tools:commit:"+step.ID, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if s.Schema != FollowStackSchema || s.Generation != env.Generation || s.Repo != "tools" || s.Front.ID != "tools:commit:"+step.ID || s.Front.Brane != "tools:brane:row-189" {
		t.Fatalf("front %s on %s (%s, %s)", s.Front.ID, s.Front.Brane, s.Schema, s.Generation)
	}
	b, _ := json.Marshal(s)
	if strings.Contains(string(b), "world:") {
		t.Fatal("a tools stack names a world id")
	}
	v1, _ := BuildStack(ix, "row-189", step.ID, fixedNow)
	if len(s.Rail) != len(v1.Rail) || len(s.Earlier) != len(v1.Earlier) || len(s.Later) != len(v1.Later) || len(*s.Front.Reads) != len(*v1.Front.Reads) {
		t.Fatal("qualification changed the stack's shape")
	}
	for _, r := range *s.Front.Reads {
		if (r.Via == surface.ViaStep || r.Via == surface.ViaCommit) && !strings.HasPrefix(r.FromID, "tools:commit:") {
			t.Fatalf("read source %s unqualified", r.FromID)
		}
	}
	// A shared stack is named for its repo and branch.
	merge, _ := BuildFollowStack(env, "world:brane:shared", "world:commit:"+step.Merge, fixedNow)
	if merge.Front.Name != "world · dev" {
		t.Fatalf("shared name %q", merge.Front.Name)
	}
	// Legacy (unqualified) addresses resolve only when exactly one repo holds them.
	_, err = BuildFollowStack(env, "row-189", step.ID, fixedNow)
	if err == nil || !strings.Contains(err.Error(), "world:row-189@f36428a") || !strings.Contains(err.Error(), "tools:row-189@f36428a") {
		t.Fatalf("ambiguous legacy address: %v", err)
	}
	one := env
	one.Repos, one.Open = env.Repos[:1], nil
	if s, err := BuildFollowStack(one, "row-189", step.ID, fixedNow); err != nil || s.Front.ID != "world:commit:"+step.ID {
		t.Fatalf("unique legacy address: %v", err)
	}
	// Episodes: an alias opens its landed step; an open episode its own card.
	if s, err := BuildFollowStack(env, "tools:brane:shared", "tools:episode:e1", fixedNow); err != nil || s.Front.ID != "tools:commit:"+step.ID || s.Front.Brane != "tools:brane:row-189" {
		t.Fatalf("alias: %v %s", err, s.Front.ID)
	}
	s, err = BuildFollowStack(env, "world:brane:shared", "world:episode:w189z", fixedNow)
	if err != nil || s.Front.Kind != KindOpen || s.Front.Brane != "world:brane:row-189" || s.Front.BuiltOn == nil || s.Front.BuiltOn.ID != env.Open.Episodes[0].Anchor ||
		s.Front.Position != "session live" || len(s.Rail) != 1 || s.Generation != env.Generation {
		t.Fatalf("open stack: %v %+v", err, s.Front)
	}
	// Missing selections are framed, never a guess.
	for _, c := range []struct{ brane, at, want string }{
		{"world:brane:row-189", "world:commit:" + strings.Repeat("b", 40), "outside world's retained window of 24 merges on dev"},
		{"other:brane:shared", "other:commit:" + step.ID, "repo other has no accepted graph"},
		{"world:brane:shared", "world:episode:gone", "neither open nor landed"},
		{"tools:brane:row-189", "world:commit:" + step.ID, "must be a brane of world"},
		{"world:brane:row-156", "world:commit:" + step.ID, "not on this brane"},
		{"row-189", "world:" + step.ID, "not a repo-qualified id"},
	} {
		if _, err := BuildFollowStack(env, c.brane, c.at, fixedNow); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %s: %v, want %q", c.brane, c.at, err, c.want)
		}
	}
}

// AC14 (static meaning): the no-script page carries each repo's timeline and
// the same freshness lines; its links are qualified.
func TestFollowStaticPages(t *testing.T) {
	env := followEnv(t)
	sources := followSources()
	m := BuildFollowModel(env, sources, Now{At: fixedNow})
	body, err := RenderFollowTimeline(env, m, Now{At: fixedNow, Healthy: true})
	if err != nil {
		t.Fatal(err)
	}
	html := string(body)
	wi, ti := strings.Index(html, `<h2 class="repo">world · dev</h2>`), strings.Index(html, `<h2 class="repo">tools · dev</h2>`)
	if wi < 0 || ti < wi {
		t.Fatal("repo groups missing or out of config order")
	}
	if !strings.Contains(html, `<p class="source" data-warn="true">tools · dev · stale · last successful local check 08:40 · Git call exceeded 2 s or queued job exceeded 20 s; next: ailang-worldd service status</p>`) ||
		!strings.Contains(html, `<p class="source">world · dev · local branch checked 08:59</p>`) {
		t.Fatal("sources missing from the static now edge")
	}
	if !strings.Contains(html, "World serving · tools stale") || !strings.Contains(html, "running w189z · world · session live") {
		t.Fatal("static health or running list missing")
	}
	if strings.Contains(html, `href="/surface/static?brane=row-`) || !strings.Contains(html, `href="/surface/static?brane=tools%3Abrane%3Arow-189&amp;at=tools%3Acommit%3A`) {
		t.Fatal("static links are not qualified")
	}
	step := stateWithPrefix(t, env.Repos[0].Index, "f36428a")
	st, err := LayoutFollowStack(env, "tools:brane:row-189", "tools:commit:"+step.ID, fixedNow, func(repo string) (string, bool) {
		for _, s := range sources {
			if s.ID == repo {
				return SourceLine(s, fixedNow)
			}
		}
		return "", false
	})
	if err != nil {
		t.Fatal(err)
	}
	page, err := RenderStack(st)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(page), `<p class="status">tools · dev · stale`) || !strings.Contains(string(page), `value="world tools:row-189@f36428a"`) ||
		strings.Contains(string(page), `brane=row-189&`) {
		t.Fatal("static stack lacks status, address or qualified links")
	}
	open, err := LayoutFollowStack(env, "world:brane:shared", "world:episode:w189z", fixedNow, nil)
	if err != nil || !strings.Contains(open.Reason, "w189z is an open episode in world: session live") {
		t.Fatalf("open episode page: %v %q", err, open.Reason)
	}
	// No accepted graph: the reason, and every repo's line.
	none := env
	none.Repos, none.Open = nil, nil
	body, err = RenderFollowTimeline(none, BuildFollowModel(none, sources, Now{At: fixedNow}), Now{At: fixedNow})
	if err != nil || !strings.Contains(string(body), "no repo has an accepted graph yet") || !strings.Contains(string(body), "tools · dev · stale") {
		t.Fatalf("empty follow page: %v", err)
	}
}

func TestFollowModelDeterministic(t *testing.T) {
	env := followEnv(t)
	a, _ := json.Marshal(BuildFollowModel(env, followSources(), Now{At: fixedNow}))
	for i := 0; i < 5; i++ {
		b, _ := json.Marshal(BuildFollowModel(env, followSources(), Now{At: fixedNow}))
		if string(a) != string(b) {
			t.Fatal("follow model is not deterministic")
		}
	}
}
