package surfacepage

// Row 188 N1: the surface's view model, the JSON the app (N2) renders from
// GET /surface/model.json and GET /surface/stack.json. The rules stay here,
// in tested Go: lanes, glyphs, ranks, links and badges are LayoutTimeline's,
// the stack is LayoutStack's, and this file only projects them, with the
// records' human layer (display titles, summaries, delivered verdicts, Today).
//
// A field the index does not hold is absent from the JSON, never filled in.

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// The two documents' types.
const (
	ModelSchema = "world/surface-model/v1"
	StackSchema = "world/surface-stack/v1"
)

// SharedGoal is the shared lane's goal; the index has no charter row for it.
const SharedGoal = "Everything merged into dev: the one world every brane starts from and returns to."

// Model is the world/surface-model/v1 document.
type Model struct {
	Schema  string      `json:"schema"`
	Now     string      `json:"now"`
	NowEdge NowEdge     `json:"nowEdge"`
	Lanes   []ModelLane `json:"lanes"`
	Links   []ModelLink `json:"links"`
}

// NowEdge is the daemon's live facts and Today.
type NowEdge struct {
	Live            *LiveCursor `json:"live,omitempty"`
	LiveUnavailable string      `json:"liveUnavailable,omitempty"`

	Health  string `json:"health"`
	Checked string `json:"checked"` // HH:MM UTC
	Daemon  string `json:"daemon"`  // the build's 7-character revision, or "unstamped"
	Clean   bool   `json:"clean"`   // stamped and not vcs.modified
	// Sessions are the credentials store's live rows, oldest first;
	// SessionsUnavailable, when set, says they could not be read.
	Sessions            []ModelSession `json:"sessions"`
	SessionsUnavailable string         `json:"sessionsUnavailable,omitempty"`
	// NeedsYou is empty in S1: nothing the daemon reads records a wait on you.
	NeedsYou []string `json:"needsYou"`
	Today    Today    `json:"today"`
}

// LiveCursor is the newest log entry read with the model.
type LiveCursor struct {
	LastIndex int64   `json:"lastIndex"`
	LogHead   *string `json:"logHead"`
}

// ModelSession is one live session.
type ModelSession struct {
	Episode string `json:"episode"`
	Started string `json:"started"` // RelativeTime
}

// Today counts the attached steps whose merge landed on now's UTC date.
type Today struct {
	Landed       int            `json:"landed"`
	ByAgent      map[string]int `json:"byAgent"`
	FirstRound   int            `json:"firstRound"`   // one round, delivered all passed
	FixRounds    int            `json:"fixRounds"`    // the sum of rounds - 1
	AgentMinutes float64        `json:"agentMinutes"` // rounded to 0.1
}

// ModelLane is the shared world or one brane.
type ModelLane struct {
	Key         string       `json:"key"`
	Name        string       `json:"name"`
	Title       string       `json:"title"` // the charter slug; "" when the charter has none
	Row         *int         `json:"row"`   // null on the shared lane
	Goal        string       `json:"goal,omitempty"`
	Glyph       string       `json:"glyph"`
	Colour      string       `json:"colour"`
	StepsLanded int          `json:"stepsLanded"`
	States      []ModelState `json:"states"`
}

// ModelLink is one LayoutTimeline link.
type ModelLink struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Colour string `json:"colour"`
}

// ModelState is one state as a person reads it.
type ModelState struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Label   string `json:"label"`
	Title   string `json:"title,omitempty"`
	At      string `json:"at"` // RFC 3339 UTC
	Rel     string `json:"rel"`
	Rank    int    `json:"rank"`
	Badge   int    `json:"badge"`
	PR      int    `json:"pr,omitempty"`
	Episode string `json:"episode,omitempty"`
	Summary string `json:"summary,omitempty"`
	Agent   string `json:"agent,omitempty"`
	Model   string `json:"model,omitempty"`
	Graders string `json:"graders,omitempty"` // the delivered verdict, "P/T"
	Rounds  int    `json:"rounds,omitempty"`
	// Minutes is the agent wall time of every round.
	Minutes *float64 `json:"minutes,omitempty"`
	// FirstRoundPass needs round 1's graders, which the index does not
	// record; it stays absent until it does.
	FirstRoundPass *bool `json:"firstRoundPass,omitempty"`
}

var (
	conventionalPrefix = regexp.MustCompile(`^[a-z]+(\([^)]*\))?!?:\s*`)
	rowPrefix          = regexp.MustCompile(`^(?i:row)\s+\d+\b[^:]{0,40}:\s*`)
	runSuffix          = regexp.MustCompile(`\s*\((?i:episode|iter|iteration)\b[^)]*\)\s*$`)
)

// DisplayTitle is a raw merge title as the surface shows it, in four steps:
// strip a leading conventional prefix ("feat(surface): ", "fix!: "), then a
// leading row reference ("Row 114 (attended): ", "row 156 M5b: ",
// "Row 156: "), then a trailing "(episode …)", "(iter …)" or
// "(iteration …)", and upper-case the first letter. "" stays "".
func DisplayTitle(raw string) string {
	t := conventionalPrefix.ReplaceAllString(raw, "")
	t = rowPrefix.ReplaceAllString(t, "")
	t = runSuffix.ReplaceAllString(t, "")
	if t == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(t)
	return string(unicode.ToUpper(r)) + t[n:]
}

// modelBuilder holds one LayoutTimeline and every state's ModelState, so the
// model and each stack share one projection.
type modelBuilder struct {
	ix     surface.Index
	now    time.Time
	tl     Timeline
	lanes  map[string]ModelLane // without states
	states map[string]ModelState
	byID   map[string]surface.State
	steps  map[string]surface.Step // by step state
}

func newModelBuilder(ix surface.Index, now time.Time) *modelBuilder {
	b := &modelBuilder{ix: ix, now: now.UTC(), tl: LayoutTimeline(ix, now.UTC()),
		lanes: map[string]ModelLane{}, states: map[string]ModelState{},
		byID: map[string]surface.State{}, steps: map[string]surface.Step{}}
	for _, s := range ix.States {
		b.byID[s.ID] = s
	}
	for _, st := range ix.Steps {
		b.steps[st.State] = st
	}
	branes := map[string]surface.Brane{}
	for _, br := range ix.Branes {
		branes[br.Key] = br
	}
	for _, l := range b.tl.Lanes {
		ml := ModelLane{Key: l.Key, Glyph: l.Glyph, Colour: l.Colour, States: []ModelState{}}
		if l.Key == SharedKey {
			ml.Name, ml.Title, ml.Goal = l.Title, l.Title, SharedGoal
		} else {
			br := branes[l.Key]
			row := br.Row
			ml.Name, ml.Title, ml.Row, ml.Goal = br.Name, br.Title, &row, br.Goal
			if ml.Name == "" {
				ml.Name = l.Title // M2's label: "row 139" when the charter has no name
			}
		}
		b.lanes[l.Key] = ml
		for _, n := range l.Nodes {
			b.states[n.ID] = b.state(b.byID[n.ID], n)
		}
	}
	return b
}

// state projects one state; n is its LayoutTimeline node.
func (b *modelBuilder) state(s surface.State, n Node) ModelState {
	at, _ := time.Parse(time.RFC3339, s.At) // Validate has checked it
	ms := ModelState{ID: s.ID, Kind: s.Kind, Label: short(s.ID), At: at.UTC().Format(time.RFC3339),
		Rel: RelativeTime(at, b.now), Rank: n.Rank, Badge: n.Badge, PR: s.PR}
	switch s.Kind {
	case surface.KindStep:
		ms.Title = DisplayTitle(s.Title)
		if s.PR != 0 {
			ms.Label = fmt.Sprintf("#%d", s.PR)
		}
		if st, ok := b.steps[s.ID]; ok {
			if st.Episode != "" {
				ms.Label = st.Episode
			}
			ms.Episode, ms.Summary, ms.Agent, ms.Model = st.Episode, st.Summary, st.Agent, st.Model
			if st.Delivered != nil {
				ms.Graders = fmt.Sprintf("%d/%d", st.Delivered.Passed, st.Delivered.Total)
			}
			ms.Rounds = st.Rounds
			minutes := st.Minutes
			ms.Minutes = &minutes
		}
	case surface.KindShared:
		// A shared state reads as the step it merged.
		for _, p := range s.Parents {
			if step, ok := b.byID[p]; ok && step.Kind == surface.KindStep && step.Merge == s.ID {
				ms.Title = DisplayTitle(step.Title)
			}
		}
	}
	return ms
}

// BuildModel is the world/surface-model/v1 document for ix as of now. It is
// pure and deterministic; map order never reaches the output.
func BuildModel(ix surface.Index, now Now) Model {
	b := newModelBuilder(ix, now.At)
	m := Model{Schema: ModelSchema, Now: b.now.Format(time.RFC3339), Lanes: []ModelLane{}, Links: []ModelLink{}}
	m.NowEdge = NowEdge{Live: now.Live, LiveUnavailable: now.LiveUnavailable, Health: "World healthy", Checked: b.now.Format("15:04"), Daemon: "unstamped",
		Sessions: []ModelSession{}, SessionsUnavailable: now.SessionsUnavailable, NeedsYou: []string{}, Today: b.today()}
	if now.Revision != "" {
		m.NowEdge.Daemon, m.NowEdge.Clean = short(now.Revision), !now.Modified
	}
	if now.SessionsUnavailable == "" {
		for _, s := range now.Sessions {
			m.NowEdge.Sessions = append(m.NowEdge.Sessions, ModelSession{Episode: s.Episode, Started: RelativeTime(s.Started, b.now)})
		}
	}
	for _, l := range b.tl.Lanes {
		ml := b.lanes[l.Key]
		for _, n := range l.Nodes {
			ml.States = append(ml.States, b.states[n.ID])
			if n.Kind == surface.KindStep {
				ml.StepsLanded++
			}
		}
		sort.SliceStable(ml.States, func(i, j int) bool { return ml.States[i].Rank < ml.States[j].Rank })
		m.Lanes = append(m.Lanes, ml)
	}
	for _, l := range b.tl.Links {
		m.Links = append(m.Links, ModelLink{From: l.From, To: l.To, Colour: l.Colour})
	}
	return m
}

// today counts steps[] whose merge landed on now's UTC date.
func (b *modelBuilder) today() Today {
	t := Today{ByAgent: map[string]int{}}
	day := b.now.Format("2006-01-02")
	minutes := 0.0
	for _, st := range b.ix.Steps {
		s, ok := b.byID[st.State]
		if !ok {
			continue
		}
		landed := s.At
		if m, ok := b.byID[s.Merge]; ok {
			landed = m.At
		}
		at, err := time.Parse(time.RFC3339, landed)
		if err != nil || at.UTC().Format("2006-01-02") != day {
			continue
		}
		t.Landed++
		t.ByAgent[st.Agent]++
		if st.Rounds == 1 && st.Delivered != nil && st.Delivered.Passed == st.Delivered.Total {
			t.FirstRound++
		}
		if st.Rounds > 1 {
			t.FixRounds += st.Rounds - 1
		}
		minutes += st.Minutes
	}
	t.AgentMinutes = math.Round(minutes*10) / 10
	return t
}

// StackModel is the world/surface-stack/v1 document.
type StackModel struct {
	Schema  string             `json:"schema"`
	Now     string             `json:"now"`
	Front   StackFront         `json:"front"`
	Earlier []StackModelColumn `json:"earlier"`
	Later   []StackModelColumn `json:"later"`
	Rail    []ModelState       `json:"rail"`
}

// StackFront is the selected state and its card.
type StackFront struct {
	ModelState
	Brane      string       `json:"brane"`
	Glyph      string       `json:"glyph"`
	Colour     string       `json:"colour"`
	Name       string       `json:"name"`
	Row        *int         `json:"row"`
	Exact      string       `json:"exact"`    // "20 Oct 10:49 UTC"
	Position   string       `json:"position"` // "N / M" on its lane
	BuiltOn    *ModelState  `json:"builtOn,omitempty"`
	MergedInto *ModelState  `json:"mergedInto,omitempty"`
	Parents    []ModelState `json:"parents"` // a shared state's parents in the index
	// Reads, ReadCalls and ReadFiles are absent when the index holds no
	// reads for the state; Files and FilesTruncated when it holds no files.
	Reads          *[]ReadGroup      `json:"reads,omitempty"`
	ReadCalls      *int              `json:"readCalls,omitempty"`
	ReadFiles      *int              `json:"readFiles,omitempty"`
	Task           string            `json:"task,omitempty"`
	Files          *[]surface.Change `json:"files,omitempty"`
	FilesTruncated *int              `json:"filesTruncated,omitempty"`
	Decisions      []FrontDecision   `json:"decisions"`
	Checks         Checks            `json:"checks"`
}

// ReadGroup is the files read from one source, in first-read order.
type ReadGroup struct {
	Via    string     `json:"via"`
	FromID string     `json:"fromId"`
	Label  string     `json:"label"`
	Glyph  string     `json:"glyph,omitempty"`
	Colour string     `json:"colour,omitempty"`
	Title  string     `json:"title,omitempty"`
	Files  []ReadFile `json:"files"`
}

// ReadFile is one path read and how many times.
type ReadFile struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// FrontDecision is one recorded decision on the state.
type FrontDecision struct {
	Actor string `json:"actor"`
	Text  string `json:"text"`
	Rel   string `json:"rel"`
}

// Checks are the episode's measured checks; each is absent when unmeasured.
type Checks struct {
	Graders string   `json:"graders,omitempty"`
	Rounds  int      `json:"rounds,omitempty"`
	Calls   *int     `json:"calls,omitempty"`
	Minutes *float64 `json:"minutes,omitempty"`
}

// StackModelColumn is one LayoutStack column.
type StackModelColumn struct {
	Brane  string       `json:"brane"`
	Glyph  string       `json:"glyph"`
	Colour string       `json:"colour"`
	Name   string       `json:"name"`
	Cards  []ModelState `json:"cards"`
}

// BuildStack is the world/surface-stack/v1 document for the state at on
// brane. Its error is LayoutStack's: an unknown brane, or at not on it.
func BuildStack(ix surface.Index, brane, at string, now time.Time) (StackModel, error) {
	b := newModelBuilder(ix, now)
	s, err := LayoutStack(ix, brane, at, b.now)
	if err != nil {
		return StackModel{}, err
	}
	out := StackModel{Schema: StackSchema, Now: b.now.Format(time.RFC3339),
		Earlier: b.columns(s.Earlier), Later: b.columns(s.Later), Rail: []ModelState{}}
	pos := 0
	for i, c := range s.Rail {
		out.Rail = append(out.Rail, b.states[c.ID])
		if c.ID == at {
			pos = i + 1
		}
	}
	lane := b.lanes[brane]
	x := b.byID[at]
	f := StackFront{ModelState: b.states[at], Brane: brane, Glyph: lane.Glyph, Colour: lane.Colour, Name: lane.Name,
		Row: lane.Row, Exact: s.Front.Exact, Position: fmt.Sprintf("%d / %d", pos, len(s.Rail)),
		Parents: []ModelState{}, Decisions: []FrontDecision{}}
	ref := func(id string) *ModelState {
		if ms, ok := b.states[id]; ok && id != "" {
			return &ms
		}
		return nil
	}
	if x.Kind == surface.KindStep {
		f.BuiltOn, f.MergedInto = ref(x.Base), ref(x.Merge)
	} else {
		for _, p := range x.Parents {
			if ms := ref(p); ms != nil {
				f.Parents = append(f.Parents, *ms)
			}
		}
	}
	if reads, ok := ix.Reads[at]; ok {
		groups, calls := []ReadGroup{}, 0
		index := map[string]int{}
		for _, r := range reads {
			calls += r.Count
			key := r.Via + ":" + r.From
			i, seen := index[key]
			if !seen {
				i = len(groups)
				index[key] = i
				groups = append(groups, b.readGroup(r))
			}
			groups[i].Files = append(groups[i].Files, ReadFile{Path: r.Path, Count: r.Count})
		}
		files := len(reads)
		f.Reads, f.ReadCalls, f.ReadFiles = &groups, &calls, &files
	}
	if fc, ok := ix.Files[at]; ok {
		changes := append([]surface.Change{}, fc.Changes...)
		truncated := fc.Truncated
		f.Files, f.FilesTruncated = &changes, &truncated
	}
	for _, d := range s.Decisions {
		f.Decisions = append(f.Decisions, FrontDecision{Actor: d.Actor, Text: d.Text, Rel: d.Relative})
	}
	if st, ok := b.steps[at]; ok && x.Kind == surface.KindStep {
		f.Task = st.Task
		f.Checks = Checks{Graders: f.Graders, Rounds: st.Rounds, Calls: st.Calls, Minutes: f.Minutes}
	}
	out.Front = f
	return out, nil
}

// readGroup is a read source's header: a step in the index is named by its
// state and lane; the others by M3's words.
func (b *modelBuilder) readGroup(r surface.Read) ReadGroup {
	g := ReadGroup{Via: r.Via, FromID: r.From, Files: []ReadFile{}}
	switch r.Via {
	case surface.ViaStep:
		if ms, ok := b.states[r.From]; ok {
			lane := b.lanes[stateBrane(b.byID[r.From])]
			g.Label, g.Glyph, g.Colour, g.Title = ms.Label, lane.Glyph, lane.Colour, ms.Title
		} else {
			g.Label = short(r.From)
		}
	case surface.ViaCommit:
		g.Label = short(r.From)
	case surface.ViaNew:
		g.Label = "made here"
	case surface.ViaOutside:
		g.Label = "older"
	}
	return g
}

func (b *modelBuilder) columns(cols []StackColumn) []StackModelColumn {
	out := []StackModelColumn{}
	for _, c := range cols {
		lane := b.lanes[c.Brane]
		mc := StackModelColumn{Brane: c.Brane, Glyph: lane.Glyph, Colour: lane.Colour, Name: lane.Name, Cards: []ModelState{}}
		for _, card := range c.Cards {
			mc.Cards = append(mc.Cards, b.states[card.ID])
		}
		out = append(out, mc)
	}
	return out
}
