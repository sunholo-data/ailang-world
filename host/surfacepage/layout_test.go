package surfacepage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

var fixedNow = time.Date(2026, 10, 21, 9, 0, 0, 0, time.UTC)

func rigIndex(t *testing.T) surface.Index {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	ix, err := surface.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := surface.Validate(ix); err != nil {
		t.Fatal(err)
	}
	return ix
}

func lane(t *testing.T, tl Timeline, key string) Lane {
	t.Helper()
	for _, l := range tl.Lanes {
		if l.Key == key {
			return l
		}
	}
	t.Fatalf("no lane %s", key)
	return Lane{}
}

func node(t *testing.T, tl Timeline, id string) Node {
	t.Helper()
	for _, l := range tl.Lanes {
		for _, n := range l.Nodes {
			if n.ID == id {
				return n
			}
		}
	}
	t.Fatalf("no node %s", id)
	return Node{}
}

// TestLayoutGolden pins the layout of the M0c rig fixture. A missing golden
// is written and the test fails, so the first run is never a pass.
func TestLayoutGolden(t *testing.T) {
	got, err := json.MarshalIndent(LayoutTimeline(rigIndex(t), fixedNow), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')
	golden := filepath.Join("testdata", "timeline-rig.golden.json")
	want, err := os.ReadFile(golden)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatal(err)
		}
		t.Fatalf("wrote %s: review it and run again", golden)
	}
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("layout differs from %s:\n%s", golden, got)
	}
}

func TestLayoutDeterministic(t *testing.T) {
	ix := rigIndex(t)
	a := LayoutTimeline(ix, fixedNow)
	for i := 0; i < 20; i++ {
		if b := LayoutTimeline(ix, fixedNow); !reflect.DeepEqual(a, b) {
			t.Fatalf("run %d differs", i)
		}
	}
	// Input order that carries no meaning does not move anything either.
	shuffled := ix
	shuffled.Branes = append([]surface.Brane(nil), ix.Branes...)
	shuffled.Steps = append([]surface.Step(nil), ix.Steps...)
	shuffled.Decisions = append([]surface.Decision(nil), ix.Decisions...)
	reverse(shuffled.Branes)
	reverse(shuffled.Steps)
	reverse(shuffled.Decisions)
	if b := LayoutTimeline(shuffled, fixedNow); !reflect.DeepEqual(a, b) {
		t.Fatal("brane, step or decision order changed the layout")
	}
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func TestLayoutRigFacts(t *testing.T) {
	ix := rigIndex(t)
	tl := LayoutTimeline(ix, fixedNow)

	var keys []string
	for _, l := range tl.Lanes {
		keys = append(keys, l.Key)
	}
	// Newest step first: row-156 20 Oct 14:25, row-189 12:49, row-198 08:25,
	// row-145 19 Oct 17:29, row-190 17:25, row-139 17 Oct 19:46.
	want := []string{"shared", "row-156", "row-189", "row-198", "row-145", "row-190", "row-139"}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("lanes = %v, want %v", keys, want)
	}

	// Every state once, by rank, 28 px apart, oldest left.
	seen := map[int]string{}
	total := 0
	for _, l := range tl.Lanes {
		for _, n := range l.Nodes {
			total++
			if prev, dup := seen[n.Rank]; dup {
				t.Errorf("rank %d held by %s and %s", n.Rank, prev, n.ID)
			}
			seen[n.Rank] = n.ID
			if n.X != Left+n.Rank*XStep || n.Y != l.Y || n.Brane != l.Key {
				t.Errorf("%s at (%d,%d) on %s", n.ID, n.X, n.Y, n.Brane)
			}
		}
	}
	if total != len(ix.States) {
		t.Fatalf("%d nodes for %d states", total, len(ix.States))
	}
	if first := node(t, tl, ix.States[0].ID); first.Rank != 0 {
		t.Errorf("oldest state rank %d", first.Rank)
	}

	// A step with no brane (#317, docs/rulings-d84-d85) sits on the shared lane.
	if n := node(t, tl, "bf0002b1a35371d6984dce6827c652c885599588"); n.Brane != SharedKey || n.Kind != surface.KindStep {
		t.Errorf("brane-less step on %s as %s", n.Brane, n.Kind)
	}
	if n := node(t, tl, "f36428abd9a72e5093b61cba14e716c37019f1ef"); n.Brane != "row-189" ||
		n.Href != "/surface/static?brane=row-189&at=f36428abd9a72e5093b61cba14e716c37019f1ef" {
		t.Errorf("w189d node = %+v", n)
	}

	// Ticks: one per UTC day, at that day's first rank, "D Mon".
	var labels []string
	for _, tk := range tl.Ticks {
		labels = append(labels, tk.Label)
	}
	if !reflect.DeepEqual(labels, []string{"17 Oct", "19 Oct", "20 Oct"}) {
		t.Fatalf("ticks = %v", labels)
	}
	if tl.Ticks[0].X != Left {
		t.Errorf("first tick at %d", tl.Ticks[0].X)
	}
	// 19 Oct 22:31+02:00 is still 19 Oct in UTC; 20 Oct starts at #325's
	// step, 02:29+02:00 (00:29 UTC).
	if want := node(t, tl, "1bbf57514e13f28f30cb3dc4042353d5d309452a").X; tl.Ticks[2].X != want {
		t.Errorf("20 Oct tick at %d, want %d", tl.Ticks[2].X, want)
	}

	// Titles: episode, lane, time in words, graders and calls.
	if got := node(t, tl, "f36428abd9a72e5093b61cba14e716c37019f1ef").Title; got != "e1 · w-world-onboarding · row 189 · yesterday 08:31 · 4/4 graders · 42 calls" {
		t.Errorf("w189d title = %q", got)
	}
	if got := node(t, tl, "4af9b24109b2f5cd376dbf91e32daf9774863c48").Title; got != "e2 · w-exec-protected-test-paths · row 156 · yesterday 08:22 · 2/2 graders" {
		t.Errorf("e2 title (no calls measured) = %q", got)
	}
	if got := node(t, tl, "adffe4c197de2c2af3ed39b9c7f272b784f6f0e9").Title; got != "adffe4c · Shared world · 18 h ago" {
		t.Errorf("shared title = %q", got)
	}
	if got := lane(t, tl, "row-139").Title; got != "row 139" {
		t.Errorf("untitled brane = %q", got)
	}
}

func TestGlyphsSevenBranes(t *testing.T) {
	ix := surface.Index{}
	base := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		key := fmt.Sprintf("row-%d", 100+i)
		ix.Branes = append(ix.Branes, surface.Brane{Key: key, Row: 100 + i, Kind: surface.BraneAttended})
		// row-100 is the newest, so lanes run row-100 … row-106.
		at := base.Add(time.Duration(10-i) * time.Hour).Format(time.RFC3339)
		ix.States = append(ix.States, surface.State{ID: fmt.Sprintf("%040x", i+1), Kind: surface.KindStep, At: at, Brane: key})
	}
	tl := LayoutTimeline(ix, base.Add(48*time.Hour))
	want := []struct{ key, glyph, colour string }{
		{"shared", GlyphCircle, SharedColour},
		{"row-100", GlyphSquare, "#2F6690"},
		{"row-101", GlyphTriangle, "#3F7A5E"},
		{"row-102", GlyphDiamond, "#6B4F8A"},
		{"row-103", GlyphHexagon, "#8A5A2B"},
		{"row-104", GlyphPentagon, "#2B7A8A"},
		{"row-105", GlyphStar, "#7A2B5E"},
		{"row-106", GlyphSquare, SlateColour},
	}
	if len(tl.Lanes) != len(want) {
		t.Fatalf("%d lanes", len(tl.Lanes))
	}
	for i, w := range want {
		l := tl.Lanes[i]
		if l.Key != w.key || l.Glyph != w.glyph || l.Colour != w.colour {
			t.Errorf("lane %d = %s %s %s, want %s %s %s", i, l.Key, l.Glyph, l.Colour, w.key, w.glyph, w.colour)
		}
		if strings.EqualFold(l.Colour, YouColour) {
			t.Errorf("lane %s uses the human's colour", l.Key)
		}
		if (l.Glyph == GlyphCircle) != (l.Mark == "") {
			t.Errorf("lane %s mark %q for %s", l.Key, l.Mark, l.Glyph)
		}
	}
	// Equal newest times fall back to key order.
	ix.States[1].At, ix.States[2].At = ix.States[0].At, ix.States[0].At
	tl = LayoutTimeline(ix, base)
	if tl.Lanes[1].Key != "row-100" || tl.Lanes[2].Key != "row-101" || tl.Lanes[3].Key != "row-102" {
		t.Errorf("tie order = %s %s %s", tl.Lanes[1].Key, tl.Lanes[2].Key, tl.Lanes[3].Key)
	}
}

func TestLinksAndBadges(t *testing.T) {
	ix := rigIndex(t)
	tl := LayoutTimeline(ix, fixedNow)
	colour := map[string]string{}
	for _, l := range tl.Lanes {
		for _, n := range l.Nodes {
			colour[n.ID] = l.Colour
		}
	}
	ids := map[string]bool{}
	for _, s := range ix.States {
		ids[s.ID] = true
	}
	// Each step: base → step when the base is in the index, then step →
	// merge, both in the step's lane colour.
	type pair struct{ from, to, colour string }
	var want []pair
	for _, s := range ix.States {
		if s.Kind != surface.KindStep {
			continue
		}
		if ids[s.Base] {
			want = append(want, pair{s.Base, s.ID, colour[s.ID]})
		}
		want = append(want, pair{s.ID, s.Merge, colour[s.ID]})
	}
	var got []pair
	for _, l := range tl.Links {
		got = append(got, pair{l.From, l.To, l.Colour})
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("links = %v\nwant %v", got, want)
	}
	// #303's base 6731d87 is outside the window: only its merge link.
	for _, l := range tl.Links {
		if l.To == "4af2f51df3d2bd10480ba3a41adbf4055b7df5df" {
			t.Errorf("link from outside the index: %+v", l)
		}
	}
	// The two-ancestor merge 9003555: w189d's link and the shared chain's.
	w189d := "f36428abd9a72e5093b61cba14e716c37019f1ef"
	found := false
	for _, l := range tl.Links {
		if l.From == w189d && l.To == "90035557b2012895551867c0ebe2f36b853aac09" && l.Colour == "#3F7A5E" {
			found = true
		}
	}
	if !found {
		t.Error("w189d's merge link missing or not in row-189's colour")
	}

	// Badges: your decisions only, and 0 shows nothing.
	for _, l := range tl.Lanes {
		for _, n := range l.Nodes {
			want := 0
			if n.ID == w189d {
				want = 1 // "merged #327", actor you; the loop's merge of #328 is not yours
			}
			if n.Badge != want {
				t.Errorf("%s badge %d, want %d", n.ID, n.Badge, want)
			}
		}
	}
	ix.Decisions = append(ix.Decisions, surface.Decision{Step: w189d, Actor: surface.ActorYou, Kind: surface.DecisionMerge},
		surface.Decision{Step: w189d, Actor: surface.ActorYou, Kind: surface.DecisionMerge})
	if n := node(t, LayoutTimeline(ix, fixedNow), w189d); n.Badge != 3 {
		t.Errorf("badge %d, want 3", n.Badge)
	}
}

func TestRelativeTimeBoundaries(t *testing.T) {
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		ago  time.Duration
		at   time.Time
		want string
	}{
		{ago: -time.Hour, want: "just now"},
		{ago: 0, want: "just now"},
		{ago: 59 * time.Second, want: "just now"},
		{ago: time.Minute, want: "1 min ago"},
		{ago: 59*time.Minute + 59*time.Second, want: "59 min ago"},
		{ago: time.Hour, want: "1 h ago"},
		{ago: 23*time.Hour + 59*time.Minute, want: "23 h ago"},
		{ago: 24 * time.Hour, want: "yesterday 12:00"},
		{at: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), want: "yesterday 00:00"},
		{at: time.Date(2026, 10, 19, 23, 59, 59, 0, time.UTC), want: "2 days ago"},
		{at: time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC), want: "30 days ago"},
		{at: time.Date(2026, 9, 20, 23, 59, 0, 0, time.UTC), want: "20 Sep"},
		{at: time.Date(2026, 1, 5, 8, 0, 0, 0, time.UTC), want: "5 Jan"},
		{at: time.Date(2025, 10, 20, 8, 0, 0, 0, time.UTC), want: "20 Oct 2025"},
	} {
		at := tc.at
		if at.IsZero() {
			at = now.Add(-tc.ago)
		}
		if got := RelativeTime(at, now); got != tc.want {
			t.Errorf("RelativeTime(%s) = %q, want %q", at.Format(time.RFC3339), got, tc.want)
		}
	}
	// A day is the calendar day in now's location, not UTC's.
	cest := time.FixedZone("CEST", 2*3600)
	if got := RelativeTime(time.Date(2026, 10, 20, 22, 30, 0, 0, time.UTC), time.Date(2026, 10, 22, 1, 0, 0, 0, cest)); got != "yesterday 00:30" {
		t.Errorf("zoned = %q", got)
	}
}
