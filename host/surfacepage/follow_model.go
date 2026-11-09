package surfacepage

// Row 220 M4: the follow-mode view model. In follow mode the daemon serves
// world/surface-model/v2 and world/surface-stack/v2, built from one published
// envelope (surface.FollowIndex): every repo's unchanged v1 partition goes
// through BuildModel/BuildStack, then every state, brane and link is
// qualified with its repo id (world:commit:<sha>, world:brane:row-220,
// world:episode:<ep>). No edge ever joins two repos. The legacy v1 documents
// are untouched: BuildModel and BuildStack keep their contracts.

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// The follow documents' types.
const (
	FollowModelSchema = "world/surface-model/v2"
	FollowStackSchema = "world/surface-stack/v2"
)

// KindOpen is an open episode's state kind in the follow model. It is never
// a commit and never a merge.
const KindOpen = "open"

// Repo freshness states, as surfacefollow publishes them.
const (
	SourceFresh    = "fresh"
	SourceUpdating = "updating"
	SourceStale    = "stale"
	SourceBuilding = "building"
	SourceRefused  = "refused"
)

// FollowSource is one configured repo's freshness, read from the same
// snapshot as the envelope. It carries no path.
type FollowSource struct {
	ID, Branch string
	State      string // fresh, updating, stale, building or refused
	// CheckedAt is the last successful local check; BuiltAt the accepted
	// graph's build time. Both RFC 3339, "" when none.
	CheckedAt, BuiltAt string
	Failure            *surface.FollowFailure
}

// FollowModel is the world/surface-model/v2 document.
type FollowModel struct {
	Schema string `json:"schema"`
	// Generation identifies the published envelope: equal generations hold
	// equal graphs, so the page need not re-render.
	Generation string       `json:"generation"`
	Now        string       `json:"now"`
	NowEdge    NowEdge      `json:"nowEdge"`
	Repos      []ModelRepo  `json:"repos"`
	Lanes      []ModelLane  `json:"lanes"`
	Links      []ModelLink  `json:"links"`
	Aliases    []ModelAlias `json:"aliases"`
	parts      map[string]surface.FollowPartition
	lanes      map[string]int // qualified key → index in Lanes
}

// ModelRepo is one configured repo: its lane group and its freshness.
type ModelRepo struct {
	ID     string `json:"id"`
	Branch string `json:"branch"`
	// Shared is the key of the repo's shared lane; "" while it has no
	// accepted graph.
	Shared string `json:"shared,omitempty"`
	Merges int    `json:"merges"` // the retained window
	Head   string `json:"head,omitempty"`
	State  string `json:"state"`
	// Line is the whole status as the page shows it, with the next command
	// when one is needed; Warn marks a state that is not fresh or updating.
	Line    string `json:"line"`
	Warn    bool   `json:"warn"`
	Checked string `json:"checked,omitempty"` // last successful local check
	Cause   string `json:"cause,omitempty"`   // the failure class
}

// ModelAlias keeps an episode address on its landed step.
type ModelAlias struct {
	ID          string `json:"id"`
	Target      string `json:"target"`
	SessionLive bool   `json:"sessionLive"`
}

// ModelOpenState is an open state's lifecycle, from the overlay.
type ModelOpenState struct {
	State             string `json:"state"`
	Label             string `json:"label"`
	SessionLive       bool   `json:"sessionLive"`
	Base              string `json:"base,omitempty"` // short hash
	BaseOutsideWindow bool   `json:"baseOutsideWindow"`
	Branch            string `json:"branch,omitempty"`
	Ended             string `json:"ended,omitempty"` // RelativeTime
}

// ModelOpen is the Running section in follow mode: the overlay's
// completeness and the live episodes, never a count that could be read as
// zero after an incomplete read.
type ModelOpen struct {
	Status string         `json:"status"`
	Note   string         `json:"note,omitempty"`
	Live   []ModelRunning `json:"live"`
	PROpen int            `json:"prOpen"`
	Ended  int            `json:"ended"`
	Failed int            `json:"failed"`
}

// ModelRunning is one live episode: its open state, its landed step, or
// the reason it has neither.
type ModelRunning struct {
	Episode string `json:"episode"`
	ID      string `json:"id,omitempty"`
	Lane    string `json:"lane,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Label   string `json:"label,omitempty"`
	Started string `json:"started,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func qualify(repo, kind, local string) string {
	id, err := surface.QualifyID(repo, kind, local)
	if err != nil {
		// Validated partitions never reach here; keep the id unusable
		// rather than unqualified.
		return repo + ":" + kind + ":invalid"
	}
	return id
}

func clockOf(rfc string, now time.Time) string {
	t, err := time.Parse(time.RFC3339, rfc)
	if err != nil {
		return ""
	}
	t = t.UTC()
	if t.Format("2006-01-02") == now.UTC().Format("2006-01-02") {
		return t.Format("15:04")
	}
	return t.Format("2 Jan 15:04")
}

// SourceLine is a repo's freshness as a person reads it. It says "local
// branch checked": this process never fetches, so remote freshness is unknown.
func SourceLine(s FollowSource, now time.Time) (line string, warn bool) {
	head := s.ID + " · " + s.Branch
	checked := clockOf(s.CheckedAt, now)
	cause := ""
	if s.Failure != nil {
		cause = s.Failure.Text
	}
	switch s.State {
	case SourceFresh:
		return head + " · local branch checked " + checked, false
	case SourceUpdating:
		return head + " · updating · showing the graph built " + clockOf(s.BuiltAt, now), false
	case SourceStale:
		last := "last successful local check " + checked
		if checked == "" {
			last = "no successful local check yet"
		}
		return head + " · stale · " + last + " · " + cause, true
	case SourceBuilding:
		return head + " · building · no graph accepted yet", false
	case SourceRefused:
		return head + " · refused · " + cause, true
	}
	return head + " · status unknown", true
}

func followHealth(repos []ModelRepo) string {
	var bad []string
	for _, r := range repos {
		if r.Warn || r.State == SourceBuilding {
			bad = append(bad, r.ID+" "+r.State)
		}
	}
	if len(bad) == 0 {
		return "World healthy"
	}
	return "World serving · " + strings.Join(bad, ", ")
}

// BuildFollowModel is the world/surface-model/v2 document for one published
// envelope and its sources, in config order. It is pure and deterministic.
func BuildFollowModel(env surface.FollowIndex, sources []FollowSource, now Now) FollowModel {
	at := now.At.UTC()
	m := FollowModel{Schema: FollowModelSchema, Generation: env.Generation, Now: at.Format(time.RFC3339),
		Repos: []ModelRepo{}, Lanes: []ModelLane{}, Links: []ModelLink{}, Aliases: []ModelAlias{},
		parts: map[string]surface.FollowPartition{}, lanes: map[string]int{}}
	for _, p := range env.Repos {
		m.parts[p.ID] = p
	}
	m.NowEdge = NowEdge{Live: now.Live, LiveUnavailable: now.LiveUnavailable, Checked: at.Format("15:04"), Daemon: "unstamped",
		Sessions: []ModelSession{}, NeedsYou: []string{}, Today: Today{ByAgent: map[string]int{}}}
	m.NowEdge.Service = modelService(now)
	if now.Revision != "" {
		m.NowEdge.Daemon, m.NowEdge.Clean = short(now.Revision), !now.Modified
	}

	// Each repo's v1 model, qualified; global ranks merge them by time.
	type ranked struct {
		lane, state int
		at          time.Time
		repo, local int
	}
	var order []ranked
	minutes := 0.0
	for ri, src := range sources {
		line, warn := SourceLine(src, at)
		repo := ModelRepo{ID: src.ID, Branch: src.Branch, State: src.State, Line: line, Warn: warn, Checked: src.CheckedAt}
		if src.Failure != nil {
			repo.Cause = src.Failure.Class
		}
		part, ok := m.parts[src.ID]
		if ok {
			repo.Shared = qualify(src.ID, "brane", SharedKey)
			repo.Merges = part.Index.Source.Merges
			repo.Head = short(part.Head)
			v1 := BuildModel(part.Index, Now{At: at})
			for _, l := range v1.Lanes {
				l.Repo = src.ID
				if l.Key == SharedKey {
					l.Name, l.Title = src.ID+" · "+src.Branch, src.ID+" · "+src.Branch
					l.Goal = fmt.Sprintf("Everything merged into %s in %s: the world every brane there starts from and returns to.", src.Branch, src.ID)
				}
				l.Key = qualify(src.ID, "brane", l.Key)
				for si := range l.States {
					s := &l.States[si]
					s.ID = qualify(src.ID, "commit", s.ID)
					t, _ := time.Parse(time.RFC3339, s.At)
					order = append(order, ranked{lane: len(m.Lanes), state: si, at: t, repo: ri, local: s.Rank})
				}
				m.lanes[l.Key] = len(m.Lanes)
				m.Lanes = append(m.Lanes, l)
			}
			for _, k := range v1.Links {
				m.Links = append(m.Links, ModelLink{From: qualify(src.ID, "commit", k.From), To: qualify(src.ID, "commit", k.To), Colour: k.Colour})
			}
			td := v1.NowEdge.Today
			m.NowEdge.Today.Landed += td.Landed
			m.NowEdge.Today.FirstRound += td.FirstRound
			m.NowEdge.Today.FixRounds += td.FixRounds
			minutes += td.AgentMinutes
			for _, a := range sortedKeys(td.ByAgent) {
				m.NowEdge.Today.ByAgent[a] += td.ByAgent[a]
			}
		}
		m.Repos = append(m.Repos, repo)
	}
	m.NowEdge.Today.AgentMinutes = math.Round(minutes*10) / 10
	sort.SliceStable(order, func(i, j int) bool {
		a, b := order[i], order[j]
		if !a.at.Equal(b.at) {
			return a.at.Before(b.at)
		}
		if a.repo != b.repo {
			return a.repo < b.repo
		}
		return a.local < b.local
	})
	for r, o := range order {
		m.Lanes[o.lane].States[o.state].Rank = r
	}
	m.addOpen(env.Open, sources, at, len(order))
	m.NowEdge.Health = followHealth(m.Repos)
	return m
}

func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// openLane is the lane an open episode sits on: its branch's row when the
// branch names one, else its repo's episode lane. Lanes the window lacks are
// added after the repo's own, with the slate glyph.
func (m *FollowModel) openLane(repo, branch string) int {
	local, name := "open", "Open episodes"
	var row *int
	if k, r, _, ok := surface.BraneOf(branch); ok {
		local, name, row = k, fmt.Sprintf("row %d", r), &r
	}
	key := qualify(repo, "brane", local)
	if i, ok := m.lanes[key]; ok {
		return i
	}
	// Insert after the repo's last lane, keeping repos grouped.
	at := len(m.Lanes)
	for i, l := range m.Lanes {
		if l.Repo == repo {
			at = i + 1
		}
	}
	lane := ModelLane{Key: key, Repo: repo, Name: name, Row: row, Glyph: GlyphSquare, Colour: SlateColour, States: []ModelState{}}
	m.Lanes = append(m.Lanes[:at], append([]ModelLane{lane}, m.Lanes[at:]...)...)
	m.lanes = map[string]int{}
	for i, l := range m.Lanes {
		m.lanes[l.Key] = i
	}
	return at
}

// addOpen adds the overlay: open states on their lanes, anchor links, the
// aliases, and the Running summary. A missing overlay is "not read yet",
// never zero.
func (m *FollowModel) addOpen(o *surface.FollowOpen, sources []FollowSource, now time.Time, rank int) {
	summary := &ModelOpen{Status: surface.OpenUnavailable, Note: "running status not read yet", Live: []ModelRunning{}}
	m.NowEdge.Open = summary
	if o == nil {
		return
	}
	summary.Status, summary.Note = o.Status, o.Reason
	var notes []string
	if o.Reason != "" {
		notes = append(notes, o.Reason)
	}
	if o.Omitted > 0 {
		notes = append(notes, fmt.Sprintf("%d more not shown", o.Omitted))
	}
	if o.RemainderUnknown && o.Status == surface.OpenComplete {
		notes = append(notes, "the remainder is unknown")
	}
	if o.Unconfigured > 0 {
		notes = append(notes, fmt.Sprintf("%d from repos not configured", o.Unconfigured))
	}
	summary.Note = strings.Join(notes, " · ")
	branchOf := map[string]string{}
	for _, s := range sources {
		branchOf[s.ID] = s.Branch
	}
	eps := append([]surface.OpenEpisode(nil), o.Episodes...)
	repoIdx := map[string]int{}
	for i, s := range sources {
		repoIdx[s.ID] = i
	}
	sort.SliceStable(eps, func(i, j int) bool {
		if eps[i].StartedAt != eps[j].StartedAt {
			return eps[i].StartedAt < eps[j].StartedAt
		}
		if repoIdx[eps[i].Repo] != repoIdx[eps[j].Repo] {
			return repoIdx[eps[i].Repo] < repoIdx[eps[j].Repo]
		}
		return eps[i].ID < eps[j].ID
	})
	for _, e := range eps {
		if _, ok := m.parts[e.Repo]; !ok {
			continue // its partition is not published; validation forbids this
		}
		li := m.openLane(e.Repo, e.Branch)
		ms := openModelState(e, now)
		ms.Rank = rank
		rank++
		lane := &m.Lanes[li]
		lane.States = append(lane.States, ms)
		if e.Anchor != "" {
			m.Links = append(m.Links, ModelLink{From: e.Anchor, To: e.ID, Colour: lane.Colour})
		}
		switch e.State {
		case surface.OpenPR:
			summary.PROpen++
		case surface.OpenEnded:
			summary.Ended++
		case surface.OpenPushFailed, surface.OpenPRFailed, surface.OpenAmbiguous:
			summary.Failed++
		}
		if e.SessionLive {
			summary.Live = append(summary.Live, ModelRunning{Episode: e.Episode, ID: e.ID, Lane: lane.Key, Repo: e.Repo,
				Label: e.Label, Started: ms.Rel})
		}
	}
	for _, a := range o.Aliases {
		m.Aliases = append(m.Aliases, ModelAlias{ID: a.ID, Target: a.Target, SessionLive: a.SessionLive})
		if !a.SessionLive {
			continue
		}
		repo, _, ep, _ := surface.ParseQualifiedID(a.ID)
		r := ModelRunning{Episode: ep, ID: a.Target, Repo: repo, Label: "landed; session live"}
		for li := range m.Lanes {
			for si := range m.Lanes[li].States {
				if s := &m.Lanes[li].States[si]; s.ID == a.Target {
					s.SessionLive = true
					r.Lane = m.Lanes[li].Key
				}
			}
		}
		summary.Live = append(summary.Live, r)
	}
	for _, r := range o.Running {
		summary.Live = append(summary.Live, ModelRunning{Episode: r.Episode, Reason: r.Reason})
	}
}

func openModelState(e surface.OpenEpisode, now time.Time) ModelState {
	started, _ := time.Parse(time.RFC3339, e.StartedAt) // ValidateFollowOpen has checked it
	ms := ModelState{ID: e.ID, Kind: KindOpen, Label: e.Episode, Title: e.Task, At: started.UTC().Format(time.RFC3339),
		Rel: RelativeTime(started, now), Episode: e.Episode, Agent: e.Agent, Model: e.Model, PR: e.PR,
		Open: &ModelOpenState{State: e.State, Label: e.Label, SessionLive: e.SessionLive, Base: short(e.Base),
			BaseOutsideWindow: e.BaseOutsideWindow, Branch: e.Branch}}
	if e.EndedAt != "" {
		if ended, err := time.Parse(time.RFC3339, e.EndedAt); err == nil {
			ms.Open.Ended = RelativeTime(ended, now)
		}
	}
	return ms
}

// FollowStack is the world/surface-stack/v2 document: a v1 stack with every
// id qualified, its repo, and the generation it was read from.
type FollowStack struct {
	StackModel
	Generation string `json:"generation"`
	Repo       string `json:"repo"`
}

// followTarget is a resolved selection: a commit on a local brane, or an
// open episode.
type followTarget struct {
	repo      string
	part      surface.FollowPartition
	brane, at string // local
	open      *surface.OpenEpisode
}

// resolveFollow resolves a page selection. Qualified ids name exactly one
// repo. An episode id resolves to its open state or, once merged, its landed
// step. An unqualified (legacy) pair resolves only when exactly one repo
// holds it; otherwise the error lists the qualified choices.
func resolveFollow(env surface.FollowIndex, brane, at string) (followTarget, error) {
	parts := map[string]surface.FollowPartition{}
	for _, p := range env.Repos {
		parts[p.ID] = p
	}
	repo, kind, local, err := surface.ParseQualifiedID(at)
	if err != nil {
		if strings.Count(at, ":") > 0 || strings.Count(brane, ":") > 0 {
			return followTarget{}, fmt.Errorf("the selection is not a repo-qualified id (repo:commit:<sha> or repo:episode:<ep>)")
		}
		var hits []followTarget
		for _, p := range env.Repos {
			for _, s := range p.Index.States {
				if s.ID == at && stateBrane(s) == brane {
					hits = append(hits, followTarget{repo: p.ID, part: p, brane: brane, at: at})
				}
			}
		}
		switch len(hits) {
		case 0:
			return followTarget{}, fmt.Errorf("no configured repo holds this state on %s in its retained window", brane)
		case 1:
			return hits[0], nil
		}
		var choices []string
		for _, h := range hits {
			choices = append(choices, h.repo+":"+brane+"@"+short(at))
		}
		return followTarget{}, fmt.Errorf("this address names a state in %d repos; choose one of: %s", len(hits), strings.Join(choices, ", "))
	}
	p, ok := parts[repo]
	if !ok {
		return followTarget{}, fmt.Errorf("repo %s has no accepted graph in this generation", repo)
	}
	if kind == "episode" {
		if env.Open != nil {
			for _, a := range env.Open.Aliases {
				if a.ID == at {
					_, _, step, _ := surface.ParseQualifiedID(a.Target)
					for _, s := range p.Index.States {
						if s.ID == step {
							return followTarget{repo: repo, part: p, brane: stateBrane(s), at: step}, nil
						}
					}
				}
			}
			for i := range env.Open.Episodes {
				if e := &env.Open.Episodes[i]; e.ID == at {
					return followTarget{repo: repo, part: p, open: e}, nil
				}
			}
		}
		return followTarget{}, fmt.Errorf("episode %s is neither open nor landed in %s's retained window", local, repo)
	}
	if kind != "commit" {
		return followTarget{}, fmt.Errorf("the selection must be a commit or an episode")
	}
	brepo, bkind, blocal, err := surface.ParseQualifiedID(brane)
	if err != nil || brepo != repo || bkind != "brane" {
		return followTarget{}, fmt.Errorf("the brane must be a brane of %s (%s:brane:<key>)", repo, repo)
	}
	for _, s := range p.Index.States {
		if s.ID == local {
			if stateBrane(s) != blocal {
				return followTarget{}, fmt.Errorf("state is not on this brane")
			}
			return followTarget{repo: repo, part: p, brane: blocal, at: local}, nil
		}
	}
	return followTarget{}, fmt.Errorf("this state is outside %s's retained window of %d merges on %s", repo, p.Index.Source.Merges, p.Index.Source.Rev)
}

// BuildFollowStack is the world/surface-stack/v2 document for a selection
// in env as of now.
func BuildFollowStack(env surface.FollowIndex, brane, at string, now time.Time) (FollowStack, error) {
	t, err := resolveFollow(env, brane, at)
	if err != nil {
		return FollowStack{}, err
	}
	out := FollowStack{Generation: env.Generation, Repo: t.repo}
	if t.open != nil {
		sources := make([]FollowSource, 0, len(env.Repos))
		for _, p := range env.Repos {
			sources = append(sources, FollowSource{ID: p.ID, Branch: p.Index.Source.Rev, State: SourceFresh})
		}
		m := BuildFollowModel(env, sources, Now{At: now})
		var lane ModelLane
		var front ModelState
		for _, l := range m.Lanes {
			for _, s := range l.States {
				if s.ID == t.open.ID {
					lane, front = l, s
				}
			}
		}
		started, _ := time.Parse(time.RFC3339, t.open.StartedAt)
		f := StackFront{ModelState: front, Brane: lane.Key, Glyph: lane.Glyph, Colour: lane.Colour, Name: lane.Name, Row: lane.Row,
			Exact: started.UTC().Format("2 Jan 15:04 UTC"), Position: t.open.Label, Parents: []ModelState{}, Decisions: []FrontDecision{}, Task: t.open.Task}
		if t.open.Anchor != "" {
			for _, l := range m.Lanes {
				for _, s := range l.States {
					if s.ID == t.open.Anchor {
						s := s
						f.BuiltOn = &s
					}
				}
			}
		}
		out.StackModel = StackModel{Schema: FollowStackSchema, Now: now.UTC().Format(time.RFC3339), Front: f,
			Earlier: []StackModelColumn{}, Later: []StackModelColumn{}, Rail: []ModelState{front}}
		return out, nil
	}
	s, err := BuildStack(t.part.Index, t.brane, t.at, now)
	if err != nil {
		return FollowStack{}, err
	}
	out.StackModel = qualifyStack(s, t.repo, t.part.Index.Source.Rev)
	return out, nil
}

func qualifyStack(s StackModel, repo, branch string) StackModel {
	qs := func(ms ModelState) ModelState { ms.ID = qualify(repo, "commit", ms.ID); return ms }
	qp := func(ms *ModelState) *ModelState {
		if ms == nil {
			return nil
		}
		q := qs(*ms)
		return &q
	}
	qb := func(key string) string { return qualify(repo, "brane", key) }
	name := func(key, n string) string {
		if key == SharedKey {
			return repo + " · " + branch
		}
		return n
	}
	s.Schema = FollowStackSchema
	f := &s.Front
	f.Name = name(f.Brane, f.Name)
	f.ModelState, f.Brane = qs(f.ModelState), qb(f.Brane)
	f.BuiltOn, f.MergedInto = qp(f.BuiltOn), qp(f.MergedInto)
	for i := range f.Parents {
		f.Parents[i] = qs(f.Parents[i])
	}
	if f.Reads != nil {
		groups := append([]ReadGroup(nil), (*f.Reads)...)
		for i := range groups {
			if groups[i].Via == surface.ViaStep || groups[i].Via == surface.ViaCommit {
				groups[i].FromID = qualify(repo, "commit", groups[i].FromID)
			}
		}
		f.Reads = &groups
	}
	cols := func(cs []StackModelColumn) []StackModelColumn {
		for i := range cs {
			cs[i].Name = name(cs[i].Brane, cs[i].Name)
			cs[i].Brane = qb(cs[i].Brane)
			for j := range cs[i].Cards {
				cs[i].Cards[j] = qs(cs[i].Cards[j])
			}
		}
		return cs
	}
	s.Earlier, s.Later = cols(s.Earlier), cols(s.Later)
	for i := range s.Rail {
		s.Rail[i] = qs(s.Rail[i])
	}
	return s
}
