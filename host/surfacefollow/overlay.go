package surfacefollow

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Row 220 M3: the open-episode overlay. An EpisodeSource (the daemon's
// read-only credentials and finish-record reads) returns facts; BuildOverlay
// is pure and joins them with the partitions of the same publication, so a
// merge that lands a finished episode replaces its open state in that
// publication.

const (
	// RetainEnded keeps a state whose session ended for seven days after
	// its last expiry or observed end. A PR-open state is kept longer only
	// while it started inside its repo's merge window (BuildOverlay).
	RetainEnded = 7 * 24 * time.Hour
	// Read bounds for one reconcile. Reaching one marks the overlay
	// incomplete with an unknown remainder; nothing is silently dropped.
	MaxOverlaySessions = 10000
	MaxOverlayStarts   = 10000
	MaxFinishLookups   = 2000
	// OverlayTimeout bounds one reconcile's reads.
	OverlayTimeout = 10 * time.Second
)

// Fixed, path-free overlay reasons.
const (
	ReasonUnavailable   = "running status unavailable: the session read failed; next: ailang-worldd session list --help"
	ReasonStaleFacts    = "running status unavailable: the latest session read failed, showing the previous read; next: ailang-worldd session list --help"
	ReasonSessionBound  = "session read bound reached; the remainder is unknown"
	ReasonStartBound    = "start record read bound reached; the remainder is unknown"
	ReasonFinishBound   = "finish record lookup bound reached; the remainder is unknown"
	ReasonCap           = "open episode cap of 1000 reached; omitted states are counted"
	ReasonInvalidIDs    = "some sessions have episode ids outside the grammar and are counted as omitted"
	ReasonOverlayRefuse = "open episode overlay refused by validation; next: ailang-worldd service status"
)

// SessionFact is one credential row: no token, grants or binding.
type SessionFact struct {
	Episode   string
	CreatedAt int64
	ExpiresAt int64
}

// StartFact is the whitelisted part of an immutable start record. RepoID is
// the configured repo its canonical repo maps to, or "" when none does; the
// source path itself never reaches this package.
type StartFact struct {
	Episode   string
	RepoID    string
	StartedAt string
	BaseSHA   string
	Branch    string
	Agent     string
	Model     string
	TaskTitle string
	// LastExpiry is the latest expiry the start's mint and resumes scheduled.
	LastExpiry int64
}

// FinishFact is a trusted row-218 finish record: operator-authenticated,
// read from the world store, never from a workspace file.
type FinishFact struct {
	Outcome string // done, push-failed, pr-failed or ambiguous
	Branch  string
	Commit  string
	PR      int
}

// EpisodeFacts is one successful read.
type EpisodeFacts struct {
	ReadAt   time.Time
	Sessions []SessionFact
	Starts   []StartFact
	Finishes map[string]FinishFact
	// Incomplete lists the fixed reasons of read bounds reached.
	Incomplete []string
	// EndObserved is when the follower first saw each episode not live. The
	// follower fills it; a source leaves it nil.
	EndObserved map[string]time.Time
}

// EpisodeSource reads facts at the host boundary. It must not write.
type EpisodeSource interface {
	ReadEpisodes(ctx context.Context, now time.Time, repos []Repo) (EpisodeFacts, error)
}

// EpisodeSink receives each reconcile: facts on success, or failed=true.
type EpisodeSink func(facts EpisodeFacts, failed bool)

type candidate struct {
	open  surface.OpenEpisode
	group int // 0 live, 1 PR open, 2 other
}

func matchesStep(s surface.State, f FinishFact) bool {
	if s.Kind != surface.KindStep {
		return false
	}
	if f.Commit != "" && s.ID == f.Commit {
		return true
	}
	// A PR number alone is never identity: the branch must agree too.
	return f.PR > 0 && s.PR == f.PR && s.Branch != "" && s.Branch == f.Branch
}

func finishState(f FinishFact) string {
	switch f.Outcome {
	case "done":
		if f.PR <= 0 {
			return surface.OpenAmbiguous
		}
		return surface.OpenPR
	case "push-failed":
		return surface.OpenPushFailed
	case "pr-failed":
		return surface.OpenPRFailed
	default:
		// Unknown outcomes are treated as ambiguous, never as a PR.
		return surface.OpenAmbiguous
	}
}

// BuildOverlay joins facts with the publication's accepted partitions. facts
// is the last successful read (nil if none); failed says the latest read
// failed. It returns nil only when no read has been attempted.
func BuildOverlay(parts []surface.FollowPartition, facts *EpisodeFacts, failed bool) *surface.FollowOpen {
	if facts == nil {
		if !failed {
			return nil
		}
		return &surface.FollowOpen{Status: surface.OpenUnavailable, Reason: ReasonUnavailable, Cap: surface.MaxOpenEpisodes,
			Episodes: []surface.OpenEpisode{}, Aliases: []surface.EpisodeAlias{}, Running: []surface.RunningEntry{}}
	}
	now := facts.ReadAt.UTC()
	o := &surface.FollowOpen{Status: surface.OpenComplete, CheckedAt: now.Format(time.RFC3339), Cap: surface.MaxOpenEpisodes,
		Episodes: []surface.OpenEpisode{}, Aliases: []surface.EpisodeAlias{}, Running: []surface.RunningEntry{}}
	reasons := append([]string(nil), facts.Incomplete...)

	states := map[string]map[string]surface.State{}
	steps := map[string][]surface.State{}
	// floor is each repo's oldest retained merge time: a PR-open episode
	// started before it could have merged outside the window unseen.
	floor := map[string]time.Time{}
	for _, p := range parts {
		m := map[string]surface.State{}
		for _, s := range p.Index.States {
			if t, err := time.Parse(time.RFC3339, s.At); err == nil {
				if f, ok := floor[p.ID]; !ok || t.Before(f) {
					floor[p.ID] = t
				}
			}
			m[s.ID] = s
			if s.Kind == surface.KindStep {
				steps[p.ID] = append(steps[p.ID], s)
			}
		}
		states[p.ID] = m
	}

	// Deduplicate credentials by episode: any unexpired row makes it live.
	live, lastExpiry := map[string]bool{}, map[string]int64{}
	episodes := map[string]bool{}
	for _, s := range facts.Sessions {
		episodes[s.Episode] = true
		if now.Unix() < s.ExpiresAt {
			live[s.Episode] = true
		}
		if s.ExpiresAt > lastExpiry[s.Episode] {
			lastExpiry[s.Episode] = s.ExpiresAt
		}
	}
	starts := map[string]StartFact{}
	for _, s := range facts.Starts {
		starts[s.Episode] = s
		episodes[s.Episode] = true
	}
	names := make([]string, 0, len(episodes))
	for ep := range episodes {
		names = append(names, ep)
	}
	sort.Strings(names)

	invalid := 0
	var cands []candidate
	for _, ep := range names {
		start, hasStart := starts[ep]
		if !surface.EpisodeGrammar.MatchString(ep) {
			if live[ep] || hasStart {
				invalid++
			}
			continue
		}
		if !hasStart {
			if live[ep] {
				o.Running = append(o.Running, surface.RunningEntry{Episode: ep, Reason: surface.ReasonNoStart})
			}
			continue
		}
		end := time.Unix(max(lastExpiry[ep], start.LastExpiry), 0).UTC()
		if seen, ok := facts.EndObserved[ep]; ok && seen.Before(end) {
			end = seen.UTC()
		}
		if end.After(now) {
			end = now
		}
		finish, finished := facts.Finishes[ep]
		prOpen := finished && finish.Outcome == "done" && finish.PR > 0
		// PR open outlives RetainEnded only while its merge would still be
		// visible: the episode started inside its repo's merge window.
		// Otherwise a merge or close outside the window would leave it
		// "PR open" forever.
		started, startErr := time.Parse(time.RFC3339, start.StartedAt)
		f, hasFloor := floor[start.RepoID]
		inWindow := startErr == nil && hasFloor && !started.Before(f)
		retained := live[ep] || (prOpen && inWindow) || now.Sub(end) <= RetainEnded
		if !retained {
			continue
		}
		if start.RepoID == "" {
			o.Unconfigured++
			if live[ep] {
				o.Running = append(o.Running, surface.RunningEntry{Episode: ep, Reason: surface.ReasonNoRepo})
			}
			continue
		}
		id, err := surface.QualifyID(start.RepoID, "episode", ep)
		if err != nil {
			invalid++
			continue
		}
		if prOpen {
			landed := ""
			for _, s := range steps[start.RepoID] {
				if matchesStep(s, finish) {
					landed = start.RepoID + ":commit:" + s.ID
					break
				}
			}
			if landed != "" {
				o.Aliases = append(o.Aliases, surface.EpisodeAlias{ID: id, Target: landed, SessionLive: live[ep]})
				continue
			}
		}
		e := surface.OpenEpisode{ID: id, Repo: start.RepoID, Episode: ep, SessionLive: live[ep],
			Branch: surface.SanitizeOpenText(start.Branch, 120), Task: surface.SanitizeOpenText(start.TaskTitle, 120),
			Agent: surface.SanitizeOpenText(start.Agent, 64), Model: surface.SanitizeOpenText(start.Model, 64)}
		if t, err := time.Parse(time.RFC3339, start.StartedAt); err == nil {
			e.StartedAt = t.UTC().Format(time.RFC3339)
		} else {
			invalid++
			continue
		}
		if _, err := surface.QualifyID(start.RepoID, "commit", start.BaseSHA); err == nil {
			e.Base = start.BaseSHA
			if _, ok := states[start.RepoID][e.Base]; ok {
				e.Anchor = start.RepoID + ":commit:" + e.Base
			}
		}
		e.BaseOutsideWindow = e.Anchor == ""
		group := 2
		switch {
		case finished:
			e.State = finishState(finish)
			if e.State == surface.OpenPR {
				e.PR, group = finish.PR, 1
			}
		case live[ep]:
			e.State = surface.OpenRunning
		default:
			e.State = surface.OpenEnded
		}
		if live[ep] {
			group = 0
		} else {
			e.EndedAt = end.Format(time.RFC3339)
		}
		e.Label = surface.OpenLabel(e.State)
		cands = append(cands, candidate{open: e, group: group})
	}
	sort.SliceStable(cands, func(i, j int) bool {
		a, b := cands[i], cands[j]
		if a.group != b.group {
			return a.group < b.group
		}
		if a.open.StartedAt != b.open.StartedAt {
			return a.open.StartedAt > b.open.StartedAt
		}
		return a.open.ID < b.open.ID
	})
	if len(cands) > surface.MaxOpenEpisodes {
		o.Omitted += len(cands) - surface.MaxOpenEpisodes
		cands = cands[:surface.MaxOpenEpisodes]
		reasons = append(reasons, ReasonCap)
	}
	for _, c := range cands {
		o.Episodes = append(o.Episodes, c.open)
	}
	if len(o.Running) > surface.MaxOpenEpisodes {
		o.Omitted += len(o.Running) - surface.MaxOpenEpisodes
		o.Running = o.Running[:surface.MaxOpenEpisodes]
		if len(reasons) == 0 || reasons[len(reasons)-1] != ReasonCap {
			reasons = append(reasons, ReasonCap)
		}
	}
	if invalid > 0 {
		o.Omitted += invalid
		reasons = append(reasons, ReasonInvalidIDs)
	}
	for _, r := range facts.Incomplete {
		if r != "" {
			o.RemainderUnknown = true
		}
	}
	if failed {
		o.Status, reasons = surface.OpenUnavailable, append([]string{ReasonStaleFacts}, reasons...)
	} else if len(reasons) > 0 {
		o.Status = surface.OpenIncomplete
	}
	o.Reason = strings.Join(reasons, "; ")
	if len(o.Reason) > 200 {
		o.Reason = reasons[0]
	}
	return o
}

// refusedOverlay replaces an overlay that failed validation, so a defect in
// episode facts never blocks merge partitions.
func refusedOverlay() *surface.FollowOpen {
	return &surface.FollowOpen{Status: surface.OpenUnavailable, Reason: ReasonOverlayRefuse, Cap: surface.MaxOpenEpisodes,
		Episodes: []surface.OpenEpisode{}, Aliases: []surface.EpisodeAlias{}, Running: []surface.RunningEntry{}}
}
