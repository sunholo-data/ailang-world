// Package surfacepage renders `GET /surface` (row 188, w-world-os-surface S1
// M2): the world timeline, laid out from the surface index by pure functions
// and rendered server-side with html/template, the way host/workbench lays out
// its graph (LayoutGraph).
//
//   - Lanes: the shared world first, then one lane per brane, the brane with
//     the newest state first.
//   - The x axis is ordinal: one fixed step per state, by rank of its time, so
//     a busy day stays as legible as a quiet one.
//   - Every brane has one glyph (shape and colour) everywhere it appears; no
//     lane uses the human's colour, which is kept for the badges.
//
// Maps serve only as lookups; nothing here iterates a map into the output.
package surfacepage

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Geometry of the timeline, in SVG user units (px).
const (
	XStep    = 28  // one rank
	Left     = 288 // lane labels sit left of the first state
	Top      = 28  // the tick row
	LaneH    = 40
	MarkX    = 16 // a lane's glyph, before its label
	LabelX   = 30
	PlotX    = Left - 16 // the plot's left edge; the label column is this wide
	padRight = 36
)

// SharedKey is the shared world's lane key.
const SharedKey = "shared"

// StaticPath is where the server-rendered pages live since N2b moved the app
// shell to /surface: the no-script fallback and the graders' structural
// fixture. Their links stay on it.
const StaticPath = "/surface/static"

// staticHref is a static stack page's link.
func staticHref(brane, id string) string {
	return StaticPath + "?brane=" + url.QueryEscape(brane) + "&at=" + url.QueryEscape(id)
}

// Glyph shapes.
const (
	GlyphCircle   = "circle"
	GlyphSquare   = "square"
	GlyphTriangle = "triangle"
	GlyphDiamond  = "diamond"
	GlyphHexagon  = "hexagon"
	GlyphPentagon = "pentagon"
	GlyphStar     = "star"
)

// Lane colours. The first six brane lanes take one each, in order; later
// lanes are slate. SharedColour is the shared world's grey.
const (
	SharedColour = "#8496A5"
	SlateColour  = "#526677"
	YouColour    = "#BD3012" // the human's; no lane uses it
)

var braneGlyphs = [6]struct{ Shape, Colour string }{
	{GlyphSquare, "#2F6690"},   // harbour
	{GlyphTriangle, "#3F7A5E"}, // sage
	{GlyphDiamond, "#6B4F8A"},  // plum
	{GlyphHexagon, "#8A5A2B"},
	{GlyphPentagon, "#2B7A8A"},
	{GlyphStar, "#7A2B5E"},
}

// shapeOffsets are each polygon glyph's vertices around its centre; a circle
// has none.
var shapeOffsets = map[string][][2]int{
	GlyphSquare:   {{-6, -6}, {6, -6}, {6, 6}, {-6, 6}},
	GlyphTriangle: {{0, -7}, {7, 6}, {-7, 6}},
	GlyphDiamond:  {{0, -8}, {8, 0}, {0, 8}, {-8, 0}},
	GlyphHexagon:  {{0, -7}, {6, -4}, {6, 4}, {0, 7}, {-6, 4}, {-6, -4}},
	GlyphPentagon: {{0, -7}, {7, -2}, {4, 6}, {-4, 6}, {-7, -2}},
	GlyphStar:     {{0, -8}, {2, -3}, {8, -2}, {3, 1}, {5, 6}, {0, 4}, {-5, 6}, {-3, 1}, {-8, -2}, {-2, -3}},
}

// Points is a glyph's SVG polygon points centred on (x, y); "" for a circle.
func Points(shape string, x, y int) string {
	offs := shapeOffsets[shape]
	parts := make([]string, 0, len(offs))
	for _, o := range offs {
		parts = append(parts, strconv.Itoa(x+o[0])+","+strconv.Itoa(y+o[1]))
	}
	return strings.Join(parts, " ")
}

// Timeline is the laid-out world timeline.
type Timeline struct {
	// Reason, when set, is why there is nothing to lay out (no index
	// configured, not built yet, or refused); the daemon sets it.
	Reason string `json:"reason,omitempty"`
	Built  string `json:"built,omitempty"` // when the index was built, in words
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Lanes  []Lane `json:"lanes"`
	Links  []Link `json:"links"`
	Ticks  []Tick `json:"ticks"`
}

// Lane is the shared world or one brane, with the states on it.
type Lane struct {
	Key    string `json:"key"` // SharedKey or a brane key
	Title  string `json:"title"`
	Row    int    `json:"row,omitempty"`
	Glyph  string `json:"glyph"`
	Colour string `json:"colour"`
	Top    int    `json:"top"`    // the lane's upper edge
	Y      int    `json:"y"`      // the lane's centre line
	Band   bool   `json:"band"`   // alternate lanes carry a band
	Mark   string `json:"mark"`   // the glyph's points beside the label
	Legend string `json:"legend"` // the glyph's points in a 20×20 legend box
	Nodes  []Node `json:"nodes"`
}

// Node is one state on its lane.
type Node struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`  // shared or step
	Brane  string `json:"brane"` // its lane's key
	Rank   int    `json:"rank"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Points string `json:"points"` // "" for a circle
	Href   string `json:"href"`
	Badge  int    `json:"badge"` // your decisions on it
	Title  string `json:"title"`
}

// Link is one line between two states: base → step, or step → merge.
type Link struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Colour string `json:"colour"` // the step's brane colour
	D      string `json:"d"`      // the SVG path
}

// Tick labels the first rank of each UTC calendar day that has a state.
type Tick struct {
	X     int    `json:"x"`
	Label string `json:"label"`
}

// LayoutTimeline lays out ix as of now. It is pure and deterministic: the
// output depends only on ix and now, never on map order.
func LayoutTimeline(ix surface.Index, now time.Time) Timeline {
	times := make([]time.Time, len(ix.States))
	for i, s := range ix.States {
		times[i], _ = time.Parse(time.RFC3339, s.At) // Validate has checked it
	}

	// Lanes: shared first, then branes by their newest state, then key.
	latest := map[string]time.Time{}
	for i, s := range ix.States {
		if s.Kind == surface.KindStep && s.Brane != "" && times[i].After(latest[s.Brane]) {
			latest[s.Brane] = times[i]
		}
	}
	branes := append([]surface.Brane(nil), ix.Branes...)
	sort.SliceStable(branes, func(i, j int) bool {
		a, b := latest[branes[i].Key], latest[branes[j].Key]
		if !a.Equal(b) {
			return a.After(b)
		}
		return branes[i].Key < branes[j].Key
	})
	lanes := []Lane{{Key: SharedKey, Title: "Shared world", Glyph: GlyphCircle, Colour: SharedColour}}
	for j, b := range branes {
		g := braneGlyphs[j%len(braneGlyphs)]
		colour := g.Colour
		if j >= len(braneGlyphs) {
			colour = SlateColour
		}
		lanes = append(lanes, Lane{Key: b.Key, Title: BraneTitle(b), Row: b.Row, Glyph: g.Shape, Colour: colour})
	}
	laneOf := map[string]int{}
	for i := range lanes {
		lanes[i].Top = Top + i*LaneH
		lanes[i].Y = lanes[i].Top + LaneH/2
		lanes[i].Band = i%2 == 1
		lanes[i].Mark = Points(lanes[i].Glyph, MarkX, lanes[i].Y)
		lanes[i].Legend = Points(lanes[i].Glyph, 10, 10)
		laneOf[lanes[i].Key] = i
	}

	// Ranks: every state by time, ties in index order (a step before the
	// shared state that merged it).
	order := make([]int, len(ix.States))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool { return times[order[a]].Before(times[order[b]]) })
	rank := make([]int, len(ix.States))
	var ticks []Tick
	lastDay := ""
	for r, i := range order {
		rank[i] = r
		if day := times[i].UTC().Format("2006-01-02"); day != lastDay {
			ticks = append(ticks, Tick{X: Left + r*XStep, Label: times[i].UTC().Format("2 Jan")})
			lastDay = day
		}
	}

	episodes := map[string]surface.Step{}
	for _, st := range ix.Steps {
		episodes[st.State] = st
	}
	yours := map[string]int{}
	for _, d := range ix.Decisions {
		if d.Actor == surface.ActorYou {
			yours[d.Step]++
		}
	}

	type at struct{ x, y, lane int }
	pos := map[string]at{}
	for i, s := range ix.States {
		li := 0
		if s.Kind == surface.KindStep {
			if k, ok := laneOf[s.Brane]; ok && s.Brane != "" {
				li = k
			}
		}
		lane := &lanes[li]
		n := Node{ID: s.ID, Kind: s.Kind, Brane: lane.Key, Rank: rank[i], X: Left + rank[i]*XStep, Y: lane.Y,
			Href: staticHref(lane.Key, s.ID)}
		n.Points = Points(lane.Glyph, n.X, n.Y)
		if s.Kind == surface.KindStep {
			n.Badge = yours[s.ID]
		}
		n.Title = nodeTitle(s, lane.Title, times[i], now, episodes)
		pos[s.ID] = at{n.X, n.Y, li}
		lane.Nodes = append(lane.Nodes, n)
	}

	links := []Link{}
	for _, s := range ix.States {
		if s.Kind != surface.KindStep {
			continue
		}
		to := pos[s.ID]
		colour := lanes[to.lane].Colour
		if from, ok := pos[s.Base]; ok && s.Base != "" {
			links = append(links, Link{From: s.Base, To: s.ID, Colour: colour, D: path(from.x, from.y, to.x, to.y)})
		}
		if merge, ok := pos[s.Merge]; ok && s.Merge != "" {
			links = append(links, Link{From: s.ID, To: s.Merge, Colour: colour, D: path(to.x, to.y, merge.x, merge.y)})
		}
	}

	for i := range lanes {
		if lanes[i].Nodes == nil {
			lanes[i].Nodes = []Node{}
		}
	}
	if ticks == nil {
		ticks = []Tick{}
	}
	built := ""
	if t, err := time.Parse(time.RFC3339, ix.BuiltAt); err == nil {
		built = RelativeTime(t, now)
	}
	return Timeline{
		Built:  built,
		Width:  Left + len(ix.States)*XStep + padRight,
		Height: Top + len(lanes)*LaneH + 8,
		Lanes:  lanes, Links: links, Ticks: ticks,
	}
}

// BraneTitle is a brane's label: its title and, for a row, "row N".
func BraneTitle(b surface.Brane) string {
	row := ""
	if b.Row > 0 {
		row = "row " + strconv.Itoa(b.Row)
	}
	switch {
	case b.Title != "" && row != "":
		return b.Title + " · " + row
	case b.Title != "":
		return b.Title
	case row != "":
		return row
	default:
		return b.Key
	}
}

// nodeTitle is a state's tooltip: its episode (or short id), its lane, the
// time in words, and for a step its graders and calls when measured.
func nodeTitle(s surface.State, lane string, t, now time.Time, episodes map[string]surface.Step) string {
	parts := []string{}
	ep, hasEp := episodes[s.ID]
	if hasEp && ep.Episode != "" {
		parts = append(parts, ep.Episode)
	} else {
		parts = append(parts, short(s.ID))
	}
	parts = append(parts, lane, RelativeTime(t, now))
	if s.Kind == surface.KindStep && hasEp {
		parts = append(parts, fmt.Sprintf("%d/%d graders", ep.Graders.Passed, ep.Graders.Total))
		if ep.Calls != nil {
			parts = append(parts, fmt.Sprintf("%d calls", *ep.Calls))
		}
	}
	return strings.Join(parts, " · ")
}

func short(id string) string {
	if len(id) > 7 {
		return id[:7]
	}
	return id
}

// path is a link's SVG path: straight along a lane, else an S-curve.
func path(x1, y1, x2, y2 int) string {
	if y1 == y2 {
		return fmt.Sprintf("M%d %d L%d %d", x1, y1, x2, y2)
	}
	mx := (x1 + x2) / 2
	return fmt.Sprintf("M%d %d C%d %d %d %d %d %d", x1, y1, mx, y1, mx, y2, x2, y2)
}

// RelativeTime says t the way people do, measured back from now, in now's
// location:
//
//   - under a minute (or later than now): "just now";
//   - under an hour: "N min ago"; under a day: "N h ago";
//   - the calendar day before now's: "yesterday HH:MM";
//   - up to 30 calendar days: "N days ago";
//   - older: "D Mon", with the year when it is not now's.
func RelativeTime(t, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return strconv.Itoa(int(d/time.Minute)) + " min ago"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d/time.Hour)) + " h ago"
	}
	t = t.In(now.Location())
	days := calendarDays(t, now)
	switch {
	case days <= 1:
		return "yesterday " + t.Format("15:04")
	case days <= 30:
		return strconv.Itoa(days) + " days ago"
	case t.Year() != now.Year():
		return t.Format("2 Jan 2006")
	default:
		return t.Format("2 Jan")
	}
}

// calendarDays counts the calendar days from t's date to now's.
func calendarDays(t, now time.Time) int {
	a := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	b := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return int(b.Sub(a) / (24 * time.Hour))
}
