package surfacepage

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Stack holds a selected state, its neighbouring history and its whole lane.
type Stack struct {
	Reason                          string
	Lane                            Lane
	Front                           StackCard
	Earlier, Later                  []StackColumn
	Rail                            []StackCard
	Inputs, Outputs                 []StackCard
	ReadSources                     []StackReadSource
	Reads                           []StackRead
	Files                           []surface.Change
	Decisions                       []StackDecision
	Task, Agent, Model, ReadSummary string
	Graders, Rounds, Calls, Minutes string
	PR                              int
}

// StackCard is a local navigation target with the lane's consistent glyph.
type StackCard struct {
	ID, Brane, Href, Title, Relative, Exact string
	Glyph, Colour, Points                   string
	Current                                 bool

	Kind string
}

// StackColumn contains at most three cards, in navigation order.
type StackColumn struct {
	Brane string
	Cards []StackCard
}

// StackRead retains the measured read and its navigable source, if indexed.
type StackRead struct {
	surface.Read
	Source StackCard
	Label  string
}

// StackReadSource summarizes files without reordering the report's reads.
type StackReadSource struct {
	Source StackCard
	Label  string
	Files  int
}

// StackDecision retains the recorded order and adds words for its actor/time.
type StackDecision struct {
	surface.Decision
	ActorTitle, Relative string
}

func stateBrane(s surface.State) string {
	if s.Kind == surface.KindStep && s.Brane != "" {
		return s.Brane
	}
	return SharedKey
}

// LayoutStack is a pure projection. Maps are lookups; traversal is time then ID.
func LayoutStack(ix surface.Index, brane, at string, now time.Time) (Stack, error) {
	out := Stack{Graders: "—", Rounds: "—", Calls: "—", Minutes: "—"}
	lanes := map[string]Lane{}
	for _, l := range LayoutTimeline(ix, now).Lanes {
		l.Nodes = nil
		lanes[l.Key] = l
	}
	l, ok := lanes[brane]
	if !ok {
		return out, fmt.Errorf("unknown brane")
	}
	out.Lane = l
	states := append([]surface.State(nil), ix.States...)
	times := map[string]time.Time{}
	byID := map[string]surface.State{}
	for _, s := range states {
		times[s.ID], _ = time.Parse(time.RFC3339, s.At)
		byID[s.ID] = s
	}
	sort.Slice(states, func(i, j int) bool {
		a, b := times[states[i].ID], times[states[j].ID]
		if !a.Equal(b) {
			return a.Before(b)
		}
		return states[i].ID < states[j].ID
	})
	x, ok := byID[at]
	if !ok || stateBrane(x) != brane {
		return out, fmt.Errorf("state is not on this brane")
	}
	episodes := map[string]string{}
	for _, step := range ix.Steps {
		episodes[step.State] = step.Episode
	}
	card := func(s surface.State) StackCard {
		title := short(s.ID)
		if s.Kind == surface.KindStep {
			if episode := episodes[s.ID]; episode != "" {
				title = episode
			} else if s.PR != 0 {
				title = fmt.Sprintf("#%d", s.PR)
			}
		}
		k := stateBrane(s)
		l := lanes[k]
		return StackCard{
			ID: s.ID, Brane: k, Kind: s.Kind,
			Href:     staticHref(k, s.ID),
			Title:    title,
			Relative: RelativeTime(times[s.ID], now),
			Exact:    times[s.ID].UTC().Format("2 Jan 15:04 UTC"),
			Glyph:    l.Glyph, Colour: l.Colour, Points: l.Legend, Current: s.ID == at,
		}
	}
	lane := func(k string) []surface.State {
		var ss []surface.State
		for _, s := range states {
			if stateBrane(s) == k {
				ss = append(ss, s)
			}
		}
		return ss
	}
	appendColumn := func(cols *[]StackColumn, k string, ss []surface.State) {
		if len(ss) == 0 {
			return
		}
		for _, c := range *cols {
			if c.Brane == k {
				return
			}
		}
		c := StackColumn{Brane: k}
		for _, s := range ss {
			c.Cards = append(c.Cards, card(s))
		}
		*cols = append(*cols, c)
	}
	own := lane(brane)
	pos := 0
	for i, s := range own {
		out.Rail = append(out.Rail, card(s))
		if s.ID == at {
			pos = i
		}
	}
	start := max(0, pos-3)
	appendColumn(&out.Earlier, brane, own[start:pos])
	end := min(len(own), pos+4)
	appendColumn(&out.Later, brane, own[pos+1:end])
	out.Front = card(x)

	// Earlier columns preserve the first appearance of each contributing brane.
	addEarlier := func(s surface.State) {
		k := stateBrane(s)
		for i := range out.Earlier {
			if out.Earlier[i].Brane == k {
				c := &out.Earlier[i]
				c.Cards = append(c.Cards, card(s))
				if len(c.Cards) > 3 {
					c.Cards = c.Cards[len(c.Cards)-3:]
				}
				return
			}
		}
		appendColumn(&out.Earlier, k, []surface.State{s})
	}
	leaving := func(base string) {
		for _, s := range states {
			if s.Kind != surface.KindStep || s.Base != base || stateBrane(s) == brane {
				continue
			}
			ss := lane(stateBrane(s))
			for i, v := range ss {
				if v.ID == s.ID {
					appendColumn(&out.Later, stateBrane(s), ss[i:min(len(ss), i+3)])
					break
				}
			}
		}
	}
	if x.Kind == surface.KindShared {
		for _, s := range states {
			if s.Kind != surface.KindStep || s.Merge != at {
				continue
			}
			ss := lane(stateBrane(s))
			for i, v := range ss {
				if v.ID == s.ID {
					appendColumn(&out.Earlier, stateBrane(s), ss[max(0, i-2):i+1])
					break
				}
			}
		}
		leaving(at)
		for _, id := range x.Parents {
			if s, ok := byID[id]; ok {
				out.Inputs = append(out.Inputs, card(s))
			} else {
				out.Inputs = append(out.Inputs, StackCard{Title: short(id)})
			}
		}
		for _, s := range states {
			if s.Kind == surface.KindStep && s.Base == at {
				out.Outputs = append(out.Outputs, card(s))
			}
		}
		for i, s := range states {
			if s.ID == at {
				for _, n := range states[i+1:] {
					if n.Kind == surface.KindShared {
						out.Outputs = append(out.Outputs, card(n))
						break
					}
				}
				break
			}
		}
	} else {
		base, hasBase := byID[x.Base]
		var prev surface.State
		for _, s := range own[:pos] {
			if s.Kind == surface.KindStep {
				prev = s
			}
		}
		pb, hasPrevBase := byID[prev.Base]
		if hasBase {
			for _, s := range states {
				inWindow := s.ID == base.ID
				if prev.ID != "" && hasPrevBase {
					inWindow = s.Kind == surface.KindShared && times[s.ID].After(times[pb.ID]) && !times[s.ID].After(times[base.ID])
				}
				if !inWindow {
					continue
				}
				for _, step := range states {
					if step.Kind == surface.KindStep && step.Merge == s.ID && stateBrane(step) != brane {
						addEarlier(step)
					}
				}
			}
			out.Inputs = append(out.Inputs, card(base))
		} else if x.Base != "" {
			out.Inputs = append(out.Inputs, StackCard{Title: short(x.Base)})
		}
		if m, ok := byID[x.Merge]; ok {
			var ss []surface.State
			for _, s := range states {
				if s.Kind == surface.KindShared && (times[s.ID].After(times[m.ID]) || times[s.ID].Equal(times[m.ID]) && s.ID >= m.ID) {
					ss = append(ss, s)
				}
			}
			appendColumn(&out.Later, SharedKey, ss[:min(len(ss), 3)])
			leaving(m.ID)
			out.Outputs = append(out.Outputs, card(m))
		}
	}
	if x.Kind == surface.KindStep {
		out.PR = x.PR
	}
	for _, s := range ix.Steps {
		if s.State == at {
			out.Front.Title = s.Episode
			out.Task, out.Agent, out.Model = s.Task, s.Agent, s.Model
			out.Graders = fmt.Sprintf("%d/%d", s.Graders.Passed, s.Graders.Total)
			out.Rounds = strconv.Itoa(s.Rounds)
			out.Minutes = strconv.FormatFloat(s.Minutes, 'f', -1, 64)
			if s.Calls != nil {
				out.Calls = strconv.Itoa(*s.Calls)
			}
		}
	}
	count := 0
	sources := map[string]int{}
	for _, r := range ix.Reads[at] {
		v := StackRead{Read: r}
		switch r.Via {
		case surface.ViaStep:
			if s, ok := byID[r.From]; ok {
				v.Source = card(s)
			}
		case surface.ViaCommit:
			v.Label = short(r.From)
		case surface.ViaNew:
			v.Label = "made here"
		case surface.ViaOutside:
			v.Label = "older"
		}
		key := r.Via + ":" + r.From
		if i, ok := sources[key]; ok {
			out.ReadSources[i].Files++
		} else {
			sources[key] = len(out.ReadSources)
			out.ReadSources = append(out.ReadSources, StackReadSource{Source: v.Source, Label: v.Label, Files: 1})
		}
		out.Reads = append(out.Reads, v)
		count += r.Count
	}
	if _, ok := ix.Reads[at]; ok {
		out.ReadSummary = fmt.Sprintf("%d reads over %d files", count, len(out.Reads))
	}
	out.Files = append([]surface.Change(nil), ix.Files[at].Changes...)
	for _, d := range ix.Decisions {
		if d.Step == at {
			t, _ := time.Parse(time.RFC3339, d.At)
			title := map[string]string{"you": "You", "world": "World", "loop": "Loop", "agent": "Agent"}[d.Actor]
			out.Decisions = append(out.Decisions, StackDecision{d, title, RelativeTime(t, now)})
		}
	}
	return out, nil
}
