package surfacefollow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Row 220 M3 overlay tests (ACs 6, 7, 8, 11 and 14 at the data level). All
// identities are fakes.

var overlayNow = time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)

func liveSession(ep string) SessionFact {
	return SessionFact{Episode: ep, CreatedAt: overlayNow.Unix() - 60, ExpiresAt: overlayNow.Unix() + 3600}
}

func startFact(ep, repo string) StartFact {
	return StartFact{Episode: ep, RepoID: repo, StartedAt: overlayNow.Add(-time.Minute).Format(time.RFC3339),
		BaseSHA: sha(3), Branch: "trial/" + ep, Agent: "claude", Model: "opus", TaskTitle: "Follow open episodes",
		LastExpiry: overlayNow.Unix() + 3600}
}

func facts(sessions []SessionFact, starts []StartFact, finishes map[string]FinishFact) *EpisodeFacts {
	if finishes == nil {
		finishes = map[string]FinishFact{}
	}
	return &EpisodeFacts{ReadAt: overlayNow, Sessions: sessions, Starts: starts, Finishes: finishes}
}

func worldParts(t *testing.T) []surface.FollowPartition {
	t.Helper()
	return []surface.FollowPartition{partitionFor(t, testRepo(), 3)}
}

// build runs BuildOverlay and insists the result validates, so every test
// below also exercises the envelope contract.
func build(t *testing.T, parts []surface.FollowPartition, f *EpisodeFacts, failed bool) *surface.FollowOpen {
	t.Helper()
	o := BuildOverlay(parts, f, failed)
	if o == nil {
		t.Fatal("no overlay")
	}
	if err := surface.ValidateFollowOpen(*o, parts); err != nil {
		t.Fatalf("overlay does not validate: %v\n%+v", err, o)
	}
	return o
}

func TestOverlayDedupesCredentialsByEpisode(t *testing.T) {
	expired := SessionFact{Episode: "w1", CreatedAt: overlayNow.Unix() - 7200, ExpiresAt: overlayNow.Unix() - 3600}
	o := build(t, worldParts(t), facts([]SessionFact{liveSession("w1"), liveSession("w1"), expired}, []StartFact{startFact("w1", "world")}, nil), false)
	if len(o.Episodes) != 1 || o.Status != surface.OpenComplete {
		t.Fatalf("want one glyph for one episode, got %+v", o)
	}
	e := o.Episodes[0]
	if e.ID != "world:episode:w1" || e.State != surface.OpenRunning || !e.SessionLive || e.Label != "session live" || e.EndedAt != "" {
		t.Fatalf("%+v", e)
	}
	// Whitelisted start facts match the immutable start.
	if e.Base != sha(3) || e.Branch != "trial/w1" || e.Task != "Follow open episodes" || e.Agent != "claude" || e.Model != "opus" ||
		e.StartedAt != overlayNow.Add(-time.Minute).Format(time.RFC3339) {
		t.Fatalf("start facts: %+v", e)
	}
}

func TestOverlayAnchorsBaseOrMarksOutsideWindow(t *testing.T) {
	in := startFact("w1", "world")
	out := startFact("w2", "world")
	out.BaseSHA = sha(99)
	o := build(t, worldParts(t), facts([]SessionFact{liveSession("w1"), liveSession("w2")}, []StartFact{in, out}, nil), false)
	got := map[string]surface.OpenEpisode{}
	for _, e := range o.Episodes {
		got[e.Episode] = e
	}
	if got["w1"].Anchor != "world:commit:"+sha(3) || got["w1"].BaseOutsideWindow {
		t.Fatalf("in-window base not anchored: %+v", got["w1"])
	}
	if got["w2"].Anchor != "" || !got["w2"].BaseOutsideWindow || got["w2"].Base != sha(99) {
		t.Fatalf("outside base not marked: %+v", got["w2"])
	}
	// With no accepted partition yet, nothing can be anchored.
	o = build(t, nil, facts([]SessionFact{liveSession("w1")}, []StartFact{in}, nil), false)
	if o.Episodes[0].Anchor != "" || !o.Episodes[0].BaseOutsideWindow {
		t.Fatalf("anchored without a partition: %+v", o.Episodes[0])
	}
}

func TestOverlayExpiryAndRevokeWithoutFinish(t *testing.T) {
	parts := worldParts(t)
	expiredAt := overlayNow.Add(-time.Hour)
	expired := SessionFact{Episode: "w1", CreatedAt: expiredAt.Unix() - 3600, ExpiresAt: expiredAt.Unix()}
	st := startFact("w1", "world")
	st.LastExpiry = expiredAt.Unix()
	o := build(t, parts, facts([]SessionFact{expired}, []StartFact{st}, nil), false)
	if len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenEnded || o.Episodes[0].SessionLive ||
		o.Episodes[0].Label != "session ended; outcome unknown" || o.Episodes[0].EndedAt != expiredAt.Format(time.RFC3339) {
		t.Fatalf("expiry: %+v", o.Episodes)
	}
	// Revoke deletes every credential row; the follower's observation dates it.
	revoked := startFact("w2", "world")
	f := facts(nil, []StartFact{revoked}, nil)
	f.EndObserved = map[string]time.Time{"w2": overlayNow.Add(-2 * time.Minute)}
	o = build(t, parts, f, false)
	if len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenEnded || o.Episodes[0].EndedAt != overlayNow.Add(-2*time.Minute).Format(time.RFC3339) {
		t.Fatalf("revoke: %+v", o.Episodes)
	}
	// Seven days after the last expiry it is retained; one second later it is not.
	for _, tc := range []struct {
		age  time.Duration
		kept bool
	}{{RetainEnded, true}, {RetainEnded + time.Second, false}} {
		st := startFact("w3", "world")
		st.LastExpiry = overlayNow.Add(-tc.age).Unix()
		o := build(t, parts, facts(nil, []StartFact{st}, nil), false)
		if (len(o.Episodes) == 1) != tc.kept {
			t.Fatalf("age %s: kept=%v, want %v", tc.age, len(o.Episodes) == 1, tc.kept)
		}
	}
}

func TestOverlayMissingStartAndUnconfiguredRepo(t *testing.T) {
	gone := SessionFact{Episode: "legacy-old", CreatedAt: 1, ExpiresAt: 2}
	o := build(t, worldParts(t), facts([]SessionFact{liveSession("legacy"), gone, liveSession("elsewhere")},
		[]StartFact{startFact("elsewhere", "")}, nil), false)
	if len(o.Episodes) != 0 {
		t.Fatalf("guessed a repo: %+v", o.Episodes)
	}
	want := []surface.RunningEntry{{Episode: "elsewhere", Reason: "repo not configured"}, {Episode: "legacy", Reason: "start record unavailable"}}
	if fmt.Sprint(o.Running) != fmt.Sprint(want) || o.Unconfigured != 1 {
		t.Fatalf("running %+v unconfigured %d", o.Running, o.Unconfigured)
	}
}

func TestOverlayCapAndReadBoundsAreExplicit(t *testing.T) {
	var sessions []SessionFact
	var starts []StartFact
	for i := 0; i < surface.MaxOpenEpisodes+1; i++ {
		ep := fmt.Sprintf("e%04d", i)
		sessions = append(sessions, liveSession(ep))
		starts = append(starts, startFact(ep, "world"))
	}
	o := build(t, worldParts(t), facts(sessions, starts, nil), false)
	if len(o.Episodes) != surface.MaxOpenEpisodes || o.Omitted != 1 || o.Status != surface.OpenIncomplete ||
		!strings.Contains(o.Reason, "cap of 1000") {
		t.Fatalf("cap: %d episodes, omitted %d, %s %q", len(o.Episodes), o.Omitted, o.Status, o.Reason)
	}
	// A read bound never yields a complete zero.
	f := facts(nil, nil, nil)
	f.Incomplete = []string{ReasonSessionBound}
	o = build(t, nil, f, false)
	if o.Status != surface.OpenIncomplete || !o.RemainderUnknown || len(o.Episodes) != 0 {
		t.Fatalf("bounded read claimed completeness: %+v", o)
	}
	// The envelope's state ceiling counts open states too.
	parts := worldParts(t)
	env := surface.FollowIndex{Schema: surface.FollowSchema, Generation: "g", Repos: parts, Open: build(t, parts, facts(sessions, starts, nil), false)}
	if err := surface.ValidateFollow(env); err != nil {
		t.Fatalf("capped envelope refused: %v", err)
	}
}

func TestOverlayUnavailableIsNeverZero(t *testing.T) {
	o := build(t, nil, nil, true)
	if o.Status != surface.OpenUnavailable || !strings.Contains(o.Reason, "running status unavailable") {
		t.Fatalf("%+v", o)
	}
	if BuildOverlay(nil, nil, false) != nil {
		t.Fatal("overlay before any read")
	}
	// A failed read keeps the previous facts, labelled unavailable.
	prev := facts([]SessionFact{liveSession("w1")}, []StartFact{startFact("w1", "world")}, nil)
	o = build(t, worldParts(t), prev, true)
	if o.Status != surface.OpenUnavailable || len(o.Episodes) != 1 || o.CheckedAt != overlayNow.Format(time.RFC3339) {
		t.Fatalf("%+v", o)
	}
}

func TestOverlayFinishBeforeRevokeAndFailureOutcomes(t *testing.T) {
	parts := worldParts(t)
	done := FinishFact{Outcome: "done", Branch: "trial/w1", Commit: sha(50), PR: 12}
	st := startFact("w1", "world")
	// Finish lands while the credential is still live.
	o := build(t, parts, facts([]SessionFact{liveSession("w1")}, []StartFact{st}, map[string]FinishFact{"w1": done}), false)
	if len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenPR || o.Episodes[0].PR != 12 || !o.Episodes[0].SessionLive ||
		o.Episodes[0].Label != "PR open" {
		t.Fatalf("live finish: %+v", o.Episodes)
	}
	// Revoke removes the credential: still PR open, and kept past seven days
	// while the episode started inside the repo's merge window.
	floor := time.Time{}
	for _, s := range parts[0].Index.States {
		if at, err := time.Parse(time.RFC3339, s.At); err == nil && (floor.IsZero() || at.Before(floor)) {
			floor = at
		}
	}
	if floor.IsZero() {
		t.Fatal("fixture has no merge times")
	}
	old := startFact("w1", "world")
	old.LastExpiry = overlayNow.Add(-30 * 24 * time.Hour).Unix()
	old.StartedAt = floor.Format(time.RFC3339)
	o = build(t, parts, facts(nil, []StartFact{old}, map[string]FinishFact{"w1": done}), false)
	if len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenPR || o.Episodes[0].SessionLive {
		t.Fatalf("revoked finish inside the window: %+v", o.Episodes)
	}
	// Started before the window: its merge or close could be unseen, so it
	// is not claimed as PR open past seven days.
	old.StartedAt = floor.Add(-time.Second).Format(time.RFC3339)
	o = build(t, parts, facts(nil, []StartFact{old}, map[string]FinishFact{"w1": done}), false)
	if len(o.Episodes) != 0 || len(o.Aliases) != 0 {
		t.Fatalf("PR open outside the window kept forever: %+v", o.Episodes)
	}
	// Within seven days it is still shown, window or not.
	old.LastExpiry = overlayNow.Add(-time.Hour).Unix()
	o = build(t, parts, facts(nil, []StartFact{old}, map[string]FinishFact{"w1": done}), false)
	if len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenPR {
		t.Fatalf("recent PR open outside the window: %+v", o.Episodes)
	}
	for outcome, want := range map[string]string{"push-failed": surface.OpenPushFailed, "pr-failed": surface.OpenPRFailed,
		"ambiguous": surface.OpenAmbiguous, "invented": surface.OpenAmbiguous, "done-without-pr": surface.OpenAmbiguous} {
		f := FinishFact{Outcome: outcome, Branch: "trial/w1", Commit: sha(2), PR: 7}
		if outcome == "done-without-pr" {
			f = FinishFact{Outcome: "done", Branch: "trial/w1", Commit: sha(2)}
		}
		st := startFact("w1", "world")
		o := build(t, parts, facts([]SessionFact{liveSession("w1")}, []StartFact{st}, map[string]FinishFact{"w1": f}), false)
		if len(o.Episodes) != 1 || o.Episodes[0].State != want || o.Episodes[0].PR != 0 || len(o.Aliases) != 0 {
			t.Fatalf("%s became %+v aliases %+v", outcome, o.Episodes, o.Aliases)
		}
	}
}

func TestOverlayMergeAliasesLandedStep(t *testing.T) {
	parts := worldParts(t) // step sha(2), PR 7 from trial/row220
	st := startFact("w1", "world")
	for name, tc := range map[string]struct {
		finish FinishFact
		live   bool
		alias  bool
	}{
		"commit":         {FinishFact{Outcome: "done", Branch: "trial/w1", Commit: sha(2), PR: 99}, false, true},
		"pr-and-branch":  {FinishFact{Outcome: "done", Branch: "trial/row220", Commit: sha(77), PR: 7}, false, true},
		"live-session":   {FinishFact{Outcome: "done", Branch: "trial/w1", Commit: sha(2), PR: 7}, true, true},
		"pr-number-only": {FinishFact{Outcome: "done", Branch: "trial/other", Commit: sha(77), PR: 7}, false, false},
	} {
		t.Run(name, func(t *testing.T) {
			var sessions []SessionFact
			if tc.live {
				sessions = []SessionFact{liveSession("w1")}
			}
			o := build(t, parts, facts(sessions, []StartFact{st}, map[string]FinishFact{"w1": tc.finish}), false)
			if !tc.alias {
				if len(o.Aliases) != 0 || len(o.Episodes) != 1 || o.Episodes[0].State != surface.OpenPR {
					t.Fatalf("attached by PR number alone: %+v %+v", o.Aliases, o.Episodes)
				}
				return
			}
			if len(o.Episodes) != 0 || len(o.Aliases) != 1 || o.Aliases[0].ID != "world:episode:w1" ||
				o.Aliases[0].Target != "world:commit:"+sha(2) || o.Aliases[0].SessionLive != tc.live {
				t.Fatalf("episodes %+v aliases %+v", o.Episodes, o.Aliases)
			}
		})
	}
	// A matching step in another repo is not this episode's merge.
	tools := testRepo()
	tools.ID = "tools"
	other := []surface.FollowPartition{partitionFor(t, tools, 3)}
	o := build(t, other, facts(nil, []StartFact{st}, map[string]FinishFact{"w1": {Outcome: "done", Commit: sha(2), PR: 7}}), false)
	if len(o.Aliases) != 0 || len(o.Episodes) != 1 {
		t.Fatalf("cross-repo alias: %+v", o)
	}
}

func TestObserveEnds(t *testing.T) {
	seen := map[string]time.Time{}
	st := []StartFact{startFact("w1", "world"), startFact("w2", "world")}
	t0 := overlayNow
	got := observeEnds(seen, EpisodeFacts{Sessions: []SessionFact{liveSession("w1")}, Starts: st}, t0)
	if _, ok := got["w1"]; ok || !got["w2"].Equal(t0) {
		t.Fatalf("%v", got)
	}
	t1 := t0.Add(2 * time.Hour) // w1's credential has expired by now
	got = observeEnds(seen, EpisodeFacts{Sessions: []SessionFact{liveSession("w1")}, Starts: st}, t1)
	if !got["w1"].Equal(t1) || !got["w2"].Equal(t0) {
		t.Fatalf("first observation must stick: %v", got)
	}
	got = observeEnds(seen, EpisodeFacts{Starts: st[:1]}, t1.Add(time.Minute))
	if _, ok := got["w2"]; ok || len(seen) != 1 {
		t.Fatalf("unread episode retained: %v", got)
	}
	got["w1"] = time.Time{}
	if seen["w1"].IsZero() {
		t.Fatal("published map aliases follower state")
	}
}

// AC 14 at the data level: what a start can hold never reaches the bytes.
func TestOverlaySerializationLeaksNothing(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	st := startFact("w1", "world")
	st.TaskTitle = "Fix login for operator@example.invalid via https://x-token:ghp_abcdefgh12345678@github.com/o/r and /Users/operator/secret/path"
	st.Branch = "trial/w1"
	st.Agent = "claude /Users/operator/bin"
	st.Model = "opus sk-abcdefghijklmnop"
	p.SinkEpisodes(*facts([]SessionFact{liveSession("w1")}, []StartFact{st}, nil), false)
	deliver(p, partitionFor(t, r, 3))
	s := p.Snapshot()
	status, _ := json.Marshal(s.Status)
	if s.Envelope.Open == nil || len(s.Envelope.Open.Episodes) != 1 {
		t.Fatalf("no open state published: %+v", s.Envelope.Open)
	}
	for _, out := range [][]byte{s.Bytes, status} {
		for _, secret := range []string{"operator@example.invalid", "x-token", "ghp_abcdefgh12345678", "github.com", "/Users/operator", "sk-abcdefghijklmnop", p.config.Repos[0].Repo} {
			if strings.Contains(string(out), secret) {
				t.Fatalf("leaked %q in %s", secret, out)
			}
		}
	}
	if task := s.Envelope.Open.Episodes[0].Task; !strings.Contains(task, "[email]") || !strings.Contains(task, "[url]") || !strings.Contains(task, "[path]") {
		t.Fatalf("redaction markers missing: %q", task)
	}
	if _, err := surface.DecodeFollow(s.Bytes); err != nil {
		t.Fatal(err)
	}
}

// AC 8: the merge and the alias land in the same publication.
func TestPublisherAliasesInTheSamePublication(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	st := startFact("w1", "world")
	st.BaseSHA = sha(3)
	done := FinishFact{Outcome: "done", Branch: "trial/w1", Commit: sha(5), PR: 8}
	p.SinkEpisodes(*facts(nil, []StartFact{st}, map[string]FinishFact{"w1": done}), false)
	deliver(p, partitionFor(t, r, 3))
	before := p.Snapshot()
	if len(before.Envelope.Open.Episodes) != 1 || before.Envelope.Open.Episodes[0].State != surface.OpenPR {
		t.Fatalf("before merge: %+v", before.Envelope.Open)
	}
	// The next head merges sha(5) as its step: the fake log's second parent.
	g := newGit()
	g.setHook(func(_ context.Context, a []string) (string, error, bool) {
		if a[3] == "log" {
			return sha(6) + "\t" + sha(3) + " " + sha(5) + "\t2026-10-29T11:00:00Z\tMerge pull request #8 from owner/trial/w1\n", nil, true
		}
		if a[3] == "merge-base" {
			return sha(3), nil, true
		}
		return "", nil, false
	})
	next, err := NewCollector(g, newClock()).BuildPinned(context.Background(), r, sha(6))
	if err != nil {
		t.Fatal(err)
	}
	deliver(p, next)
	after := p.Snapshot()
	ix, err := surface.DecodeFollow(after.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if ix.Repos[0].Head != sha(6) || len(ix.Open.Episodes) != 0 || len(ix.Open.Aliases) != 1 ||
		ix.Open.Aliases[0].Target != "world:commit:"+sha(5) || ix.Generation == before.Envelope.Generation {
		t.Fatalf("merge and alias not in one publication: %+v", ix.Open)
	}
	if after.Status.Open == nil || after.Status.Open.Episodes != 0 || after.Status.Generation != ix.Generation {
		t.Fatalf("status disagrees: %+v", after.Status.Open)
	}
}

// A reconcile that changes only checked_at is not a new generation.
func TestOverlayCheckedAtIsNotAGeneration(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	deliver(p, partitionFor(t, r, 3))
	f := facts([]SessionFact{liveSession("w1")}, []StartFact{startFact("w1", "world")}, nil)
	p.SinkEpisodes(*f, false)
	p.publish()
	g := p.Snapshot().Status.Generation
	f.ReadAt = f.ReadAt.Add(PollInterval)
	p.SinkEpisodes(*f, false)
	p.publish()
	if s := p.Snapshot(); s.Status.Generation != g || s.Envelope.Open.CheckedAt != f.ReadAt.Format(time.RFC3339) {
		t.Fatalf("generation %s→%s checked %s", g, s.Status.Generation, s.Envelope.Open.CheckedAt)
	}
	// A failed read changes the published meaning, so it is a new generation.
	p.SinkEpisodes(EpisodeFacts{}, true)
	p.publish()
	if s := p.Snapshot(); s.Status.Generation == g || s.Envelope.Open.Status != surface.OpenUnavailable || len(s.Envelope.Open.Episodes) != 1 {
		t.Fatalf("failed read: %+v", s.Envelope.Open)
	}
}

type fakeEpisodes struct {
	mu    sync.Mutex
	calls int
	facts EpisodeFacts
	err   error
	gate  chan struct{}
	// deaf ignores cancellation, like a read stuck in the driver.
	deaf bool
}

func (f *fakeEpisodes) ReadEpisodes(ctx context.Context, now time.Time, _ []Repo) (EpisodeFacts, error) {
	f.mu.Lock()
	f.calls++
	gate, out, err, deaf := f.gate, f.facts, f.err, f.deaf
	f.mu.Unlock()
	if gate != nil && deaf {
		<-gate
	} else if gate != nil {
		select {
		case <-gate:
		case <-ctx.Done():
			return EpisodeFacts{}, ctx.Err()
		}
	}
	return out, err
}

func (f *fakeEpisodes) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func (f *fakeEpisodes) set(mutate func(*fakeEpisodes)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	mutate(f)
}

func runWithEpisodes(t *testing.T, p *Publisher, g *fakeGit, clock *fakeClock, src EpisodeSource) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	fd, pd := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(pd)
		p.Run(ctx)
	}()
	go func() {
		defer close(fd)
		_ = (&Follower{Config: p.config, Collector: NewCollector(g, clock), Sink: p.Sink, Episodes: src, EpisodeSink: p.SinkEpisodes}).Run(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		for _, done := range []chan struct{}{fd, pd} {
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Error("follower or publisher did not stop")
			}
		}
	})
}

// AC 6 with the scheduler: reconcile every 15 s with git unchanged, and a new
// live episode appears within 30 s without an agent-side call.
func TestReconcileEveryPollWithGitUnchanged(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	r = p.config.Repos[0]
	g, clock := newGit(), newClock()
	g.move(r.Repo, sha(3))
	src := &fakeEpisodes{}
	runWithEpisodes(t, p, g, clock, src)
	until(t, func() bool {
		s := p.Snapshot()
		return len(s.Envelope.Repos) == 1 && s.Envelope.Open != nil && src.count() == 1
	})
	if s := p.Snapshot(); s.Envelope.Open.Status != surface.OpenComplete || len(s.Envelope.Open.Episodes) != 0 {
		t.Fatalf("%+v", s.Envelope.Open)
	}
	logs := countCommand(g, "log")
	minted := clock.Now().Add(time.Second)
	src.set(func(f *fakeEpisodes) {
		f.facts = EpisodeFacts{Sessions: []SessionFact{{Episode: "w9", CreatedAt: minted.Unix(), ExpiresAt: minted.Unix() + 3600}},
			Starts: []StartFact{{Episode: "w9", RepoID: "world", StartedAt: minted.Format(time.RFC3339), BaseSHA: sha(3),
				Branch: "trial/w9", Agent: "codex", Model: "gpt", LastExpiry: minted.Unix() + 3600}}}
	})
	clock.advance(PollInterval)
	until(t, func() bool {
		o := p.Snapshot().Envelope.Open
		return o != nil && len(o.Episodes) == 1 && o.Episodes[0].ID == "world:episode:w9"
	})
	if src.count() != 2 || countCommand(g, "log") != logs {
		t.Fatalf("reconciles %d, logs %d→%d", src.count(), logs, countCommand(g, "log"))
	}
	if e := p.Snapshot().Envelope.Open.Episodes[0]; e.Anchor != "world:commit:"+sha(3) || e.Agent != "codex" || e.State != surface.OpenRunning {
		t.Fatalf("%+v", e)
	}
	// A failing read is published as unavailable, never as zero running.
	src.set(func(f *fakeEpisodes) { f.err = errors.New("credentials store at /private/path is locked") })
	clock.advance(PollInterval)
	until(t, func() bool { return p.Snapshot().Envelope.Open.Status == surface.OpenUnavailable })
	if s := p.Snapshot(); len(s.Envelope.Open.Episodes) != 1 || strings.Contains(string(s.Bytes), "/private/path") {
		t.Fatalf("%+v", s.Envelope.Open)
	}
}

func TestReconcilesNeverOverlap(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	g, clock := newGit(), newClock()
	g.move(p.config.Repos[0].Repo, sha(3))
	gate := make(chan struct{})
	src := &fakeEpisodes{gate: gate, deaf: true}
	runWithEpisodes(t, p, g, clock, src)
	until(t, func() bool { return src.count() == 1 })
	for i := 0; i < 3; i++ {
		clock.advance(PollInterval)
		until(t, func() bool { return countCommand(g, "rev-parse") >= 2+i })
	}
	if src.count() != 1 {
		t.Fatalf("overlapping reconciles: %d", src.count())
	}
	src.set(func(f *fakeEpisodes) { f.gate = nil })
	close(gate)
	// Ticks during the blocked read coalesce into exactly one follow-up.
	until(t, func() bool { return src.count() == 2 })
	until(t, func() bool { return p.Snapshot().Envelope.Open != nil })
	if src.count() != 2 {
		t.Fatalf("reconciles %d, want 2", src.count())
	}
}

func TestReconcileDeadline(t *testing.T) {
	r := testRepo()
	p := publisherFor(t, r)
	g, clock := newGit(), newClock()
	g.move(p.config.Repos[0].Repo, sha(3))
	src := &fakeEpisodes{gate: make(chan struct{})}
	runWithEpisodes(t, p, g, clock, src)
	until(t, func() bool { return src.count() == 1 })
	clock.advance(OverlayTimeout + time.Nanosecond)
	until(t, func() bool { o := p.Snapshot().Envelope.Open; return o != nil && o.Status == surface.OpenUnavailable })
}
