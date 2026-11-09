package surfacepage

// Row 220 M4: the static (no-script) pages in follow mode carry the same
// meaning as the app: one timeline per configured repo, in config order,
// every id repo-qualified, and each repo's freshness on the now edge.

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// QualifyTimeline qualifies every lane key, state id, link end and href of
// one repo's timeline.
func QualifyTimeline(tl Timeline, repo string) Timeline {
	out := tl
	out.Lanes = make([]Lane, len(tl.Lanes))
	for i, l := range tl.Lanes {
		l.Key = qualify(repo, "brane", l.Key)
		nodes := make([]Node, len(l.Nodes))
		for j, n := range l.Nodes {
			n.ID, n.Brane = qualify(repo, "commit", n.ID), l.Key
			n.Href = staticHref(n.Brane, n.ID)
			nodes[j] = n
		}
		l.Nodes = nodes
		out.Lanes[i] = l
	}
	out.Links = make([]Link, len(tl.Links))
	for i, k := range tl.Links {
		k.From, k.To = qualify(repo, "commit", k.From), qualify(repo, "commit", k.To)
		out.Links[i] = k
	}
	return out
}

// RenderFollowTimeline renders /surface/static in follow mode from the same
// envelope and model the app reads.
func RenderFollowTimeline(env surface.FollowIndex, m FollowModel, now Now) ([]byte, error) {
	p := page{Logo: template.HTML(logoSVG)}
	parts := map[string]surface.FollowPartition{}
	for _, part := range env.Repos {
		parts[part.ID] = part
	}
	for _, r := range m.Repos {
		p.Sources = append(p.Sources, pageSource{Line: r.Line, Warn: r.Warn})
		heading := r.ID + " · " + r.Branch
		part, ok := parts[r.ID]
		if !ok {
			p.Groups = append(p.Groups, pageGroup{Heading: heading, Reason: r.Line})
			continue
		}
		tl := QualifyTimeline(LayoutTimeline(part.Index, now.At), r.ID)
		for _, l := range tl.Lanes {
			p.States += len(l.Nodes)
		}
		p.Groups = append(p.Groups, groupOf(tl, heading))
	}
	if p.States > surface.MaxStates {
		return nil, fmt.Errorf("surfacepage: %d states, over the bound of %d", p.States, surface.MaxStates)
	}
	if len(env.Repos) == 0 {
		p.T.Reason = "no repo has an accepted graph yet"
	}
	now.Healthy = true
	p.Health = healthLine(now)
	if m.NowEdge.Health != "World healthy" {
		p.Health = m.NowEdge.Health + p.Health[len("World healthy"):]
	}
	open := m.NowEdge.Open
	if open == nil || open.Status == surface.OpenUnavailable {
		p.Note = "running status unavailable"
		if open != nil && open.Note != "" {
			p.Note = open.Note
		}
	} else {
		p.Count = strconv.Itoa(len(open.Live)) + " live sessions"
		if len(open.Live) == 1 {
			p.Count = "1 live session"
		}
		for _, l := range open.Live {
			what := l.Label
			if l.Reason != "" {
				what = l.Reason
			}
			line := "running " + l.Episode
			if l.Repo != "" {
				line += " · " + l.Repo
			}
			p.Running = append(p.Running, line+" · "+what)
		}
		if open.Note != "" {
			p.Running = append(p.Running, open.Note)
		}
	}
	var buf bytes.Buffer
	if err := pageTemplate.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func qualifyCard(c StackCard, repo string) StackCard {
	if c.Brane != "" {
		c.Brane = qualify(repo, "brane", c.Brane)
	}
	if c.ID != "" {
		c.ID = qualify(repo, "commit", c.ID)
	}
	if c.Href != "" {
		c.Href = staticHref(c.Brane, c.ID)
	}
	return c
}

// LayoutFollowStack is the static stack page for a selection in follow
// mode. An open episode has no landed stack yet; its page says what it is.
// status is the repo's freshness line, shown when warn is set.
func LayoutFollowStack(env surface.FollowIndex, brane, at string, now time.Time, status func(repo string) (string, bool)) (Stack, error) {
	t, err := resolveFollow(env, brane, at)
	if err != nil {
		return Stack{}, err
	}
	var out Stack
	if t.open != nil {
		e := t.open
		out = Stack{Reason: fmt.Sprintf("%s is an open episode in %s: %s. Nothing of it has landed yet.", e.Episode, e.Repo, e.Label),
			Address: "world " + e.Repo + ":" + e.Episode}
		if e.Task != "" {
			out.Reason += " Task: " + e.Task + "."
		}
		out.Lane.Title = e.Repo + " · " + t.part.Index.Source.Rev
	} else {
		out, err = LayoutStack(t.part.Index, t.brane, t.at, now)
		if err != nil {
			return Stack{}, err
		}
		if out.Lane.Key == SharedKey {
			out.Lane.Title = t.repo + " · " + t.part.Index.Source.Rev
		}
		out.Address = "world " + t.repo + ":" + out.Lane.Key + "@" + short(t.at)
		out.Lane.Key = qualify(t.repo, "brane", out.Lane.Key)
		out.Front = qualifyCard(out.Front, t.repo)
		cards := func(cs []StackCard) []StackCard {
			for i := range cs {
				cs[i] = qualifyCard(cs[i], t.repo)
			}
			return cs
		}
		out.Rail, out.Inputs, out.Outputs = cards(out.Rail), cards(out.Inputs), cards(out.Outputs)
		for _, cols := range [][]StackColumn{out.Earlier, out.Later} {
			for i := range cols {
				cols[i].Brane = qualify(t.repo, "brane", cols[i].Brane)
				cols[i].Cards = cards(cols[i].Cards)
			}
		}
		for i := range out.ReadSources {
			out.ReadSources[i].Source = qualifyCard(out.ReadSources[i].Source, t.repo)
		}
		for i := range out.Reads {
			out.Reads[i].Source = qualifyCard(out.Reads[i].Source, t.repo)
		}
	}
	if status != nil {
		if line, warn := status(t.repo); warn {
			out.Status = line
		}
	}
	return out, nil
}
