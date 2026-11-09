package daemon

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sort"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

// finishReader is the world-store subset the trusted finish read needs.
type finishReader interface {
	GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error)
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
}

// finishProvenance is the only provenance handleEpisodeFinish writes: the
// operator route, after its proof and approval checks (D-WORLD-100 R6).
const finishProvenance = "operator:episode-finish"

// followEpisodes is the row 220 M3 episode source: read-only pages of the
// credentials store and point reads of row 218's finish records. Calls never
// overlap (the follower serializes them), so its caches need no lock.
type followEpisodes struct {
	credentialsPath string
	finishes        finishReader
	canonical       func(string) string
	canonCache      map[string]string
	finishCache     map[hashref.HashRef]*surfacefollow.FinishFact
}

func newFollowEpisodes(credentialsPath string, finishes finishReader) *followEpisodes {
	return &followEpisodes{credentialsPath: credentialsPath, finishes: finishes, canonical: canonicalRepo,
		canonCache: map[string]string{}, finishCache: map[hashref.HashRef]*surfacefollow.FinishFact{}}
}

// canonicalRepo resolves a recorded start repo the way LoadConfig resolves
// configured repos, so a symlinked spelling maps to the same repo.
func canonicalRepo(p string) string {
	if c, err := filepath.EvalSymlinks(p); err == nil {
		return c
	}
	return filepath.Clean(p)
}

func lastScheduledExpiry(s store.EpisodeStartSummary) int64 {
	var last int64
	if t, err := time.Parse(time.RFC3339, s.StartedAt); err == nil {
		last = t.Unix() + s.TTLS
	}
	for _, r := range s.Resumes {
		if r.Failed {
			continue
		}
		if t, err := time.Parse(time.RFC3339, r.ResumedAt); err == nil && t.Unix()+r.TTLS > last {
			last = t.Unix() + r.TTLS
		}
	}
	return last
}

func (e *followEpisodes) ReadEpisodes(ctx context.Context, now time.Time, repos []surfacefollow.Repo) (surfacefollow.EpisodeFacts, error) {
	facts := surfacefollow.EpisodeFacts{ReadAt: now, Finishes: map[string]surfacefollow.FinishFact{}}
	ro, err := store.OpenCredentialsReadOnly(e.credentialsPath)
	if err != nil {
		return facts, err
	}
	defer func() { _ = ro.Close() }()

	var cursor store.SessionCursor
	for {
		page, next, err := ro.ListSessionsAfter(ctx, cursor, store.MaxOverlayPage)
		if err != nil {
			return facts, err
		}
		for _, r := range page {
			facts.Sessions = append(facts.Sessions, surfacefollow.SessionFact{Episode: r.EpisodeID, CreatedAt: r.CreatedAt, ExpiresAt: r.ExpiresAt})
		}
		cursor = next
		if len(page) < store.MaxOverlayPage {
			break
		}
		if len(facts.Sessions) >= surfacefollow.MaxOverlaySessions {
			facts.Incomplete = append(facts.Incomplete, surfacefollow.ReasonSessionBound)
			break
		}
	}

	byPath := map[string]string{}
	for _, r := range repos {
		byPath[r.Repo] = r.ID
	}
	after := ""
	for {
		page, err := ro.ListEpisodeStartSummaries(ctx, after, store.MaxOverlayPage)
		if err != nil {
			return facts, err
		}
		for _, s := range page {
			canon, ok := e.canonCache[s.Repo]
			if !ok {
				if len(e.canonCache) >= 4096 {
					e.canonCache = map[string]string{}
				}
				canon = e.canonical(s.Repo)
				e.canonCache[s.Repo] = canon
			}
			facts.Starts = append(facts.Starts, surfacefollow.StartFact{Episode: s.Episode, RepoID: byPath[canon],
				StartedAt: s.StartedAt, BaseSHA: s.BaseSHA, Branch: s.Branch, Agent: s.Agent, Model: s.Model,
				TaskTitle: s.TaskTitle, LastExpiry: lastScheduledExpiry(s)})
			after = s.Episode
		}
		if len(page) < store.MaxOverlayPage {
			break
		}
		if len(facts.Starts) >= surfacefollow.MaxOverlayStarts {
			facts.Incomplete = append(facts.Incomplete, surfacefollow.ReasonStartBound)
			break
		}
	}

	// Finish lookups: configured repos only, live sessions first, newest next.
	live := map[string]bool{}
	for _, s := range facts.Sessions {
		if now.Unix() < s.ExpiresAt {
			live[s.Episode] = true
		}
	}
	var lookups []surfacefollow.StartFact
	for _, s := range facts.Starts {
		if s.RepoID != "" {
			lookups = append(lookups, s)
		}
	}
	sort.SliceStable(lookups, func(i, j int) bool {
		if live[lookups[i].Episode] != live[lookups[j].Episode] {
			return live[lookups[i].Episode]
		}
		return lookups[i].StartedAt > lookups[j].StartedAt
	})
	if len(lookups) > surfacefollow.MaxFinishLookups {
		lookups = lookups[:surfacefollow.MaxFinishLookups]
		facts.Incomplete = append(facts.Incomplete, surfacefollow.ReasonFinishBound)
	}
	for _, s := range lookups {
		f, err := e.trustedFinish(ctx, s.Episode)
		if err != nil {
			return facts, err
		}
		if f != nil {
			facts.Finishes[s.Episode] = *f
		}
	}
	return facts, nil
}

// trustedFinish returns the episode's current finish record only when the
// operator route wrote it. Any other object under the head is not a finish.
func (e *followEpisodes) trustedFinish(ctx context.Context, episode string) (*surfacefollow.FinishFact, error) {
	ref, ok, err := e.finishes.GetRegistryHead(ctx, finishHead(episode))
	if err != nil || !ok {
		return nil, err
	}
	if f, ok := e.finishCache[ref]; ok {
		return f, nil
	}
	obj, ok, err := e.finishes.GetObject(ctx, ref)
	if err != nil {
		return nil, err
	}
	var fact *surfacefollow.FinishFact
	var record EpisodeFinish
	if ok && obj.SemanticID == EpisodeFinishV1 && obj.Provenance == finishProvenance &&
		json.Unmarshal(obj.Payload, &record) == nil && record.Episode == episode {
		switch record.Outcome {
		case "done", "push-failed", "pr-failed", "ambiguous":
			fact = &surfacefollow.FinishFact{Outcome: record.Outcome, Branch: record.Branch, Commit: record.Commit, PR: record.PR.Number}
		}
	}
	if !ok {
		return nil, nil
	}
	if len(e.finishCache) >= 4096 {
		e.finishCache = map[hashref.HashRef]*surfacefollow.FinishFact{}
	}
	// Objects are content-addressed, so a ref's verdict never changes.
	e.finishCache[ref] = fact
	return fact, nil
}
