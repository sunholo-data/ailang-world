package surfacepage

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestSurfaceDisplayTitle(t *testing.T) {
	for _, tc := range []struct{ raw, want string }{
		// The operator's rule, by example.
		{"Row 114 (attended): Q3 + K5 Q4–Q6 — 2 of ≥3 qualifying walks (Q4, Q6); F-7 why hides edit paths",
			"Q3 + K5 Q4–Q6 — 2 of ≥3 qualifying walks (Q4, Q6); F-7 why hides edit paths"},
		{"Row 156: [r4] Unicode-preserving detection redesign + M1c/M1d (iteration 282)",
			"[r4] Unicode-preserving detection redesign + M1c/M1d"},
		{"row 156 M2a+M2b: protected profile key + startup refusal rows (iteration 283)",
			"Protected profile key + startup refusal rows"},
		{"docs(row 189 M6): the Use World on your repo page (episode w189e)", "The Use World on your repo page"},
		{"Merge dev into feature", "Merge dev into feature"},
		// Edge cases.
		{"", ""},
		{"feat: surface model route", "Surface model route"},
		{"feat!: drop the old route", "Drop the old route"},
		{"fix(daemon)!: bound the read", "Bound the read"},
		{"docs(row-188):   plan", "Plan"},
		{"row 189 M1: credentials store", "Credentials store"},
		{"exec detection (episode w156a)", "Exec detection"},
		{"exec detection (iter 288)", "Exec detection"},
		{"exec detection (Iteration 3)  ", "Exec detection"},
		{"feat(surface): row 188 M3: the stack (episode w188f)", "The stack"},
		{"élan vital", "Élan vital"},
		// Only the forms the rule names are stripped.
		{"Feat: capitalised prefix stays", "Feat: capitalised prefix stays"},
		{"rowing 5: a word, not a row", "Rowing 5: a word, not a row"},
		{"row 156abc: no boundary stays", "Row 156abc: no boundary stays"},
		{"row 7 more than forty characters of milestone text here: stays", "Row 7 more than forty characters of milestone text here: stays"},
		{"a note (episodes are rounds)", "A note (episodes are rounds)"},
		{"exec (iter 3) detection", "Exec (iter 3) detection"},
	} {
		if got := DisplayTitle(tc.raw); got != tc.want {
			t.Errorf("DisplayTitle(%q) = %q, want %q", tc.raw, got, tc.want)
		}
	}
}

// golden compares got with testdata/name, writing a missing golden and
// failing, so the first run is never a pass.
func golden(t *testing.T, name string, v any) {
	t.Helper()
	got, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')
	path := filepath.Join("testdata", name)
	want, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		t.Fatalf("wrote %s: review it and run again", path)
	}
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s differs:\n%s", path, got)
	}
}

var modelNow = Now{At: fixedNow, Healthy: true, Revision: "0123456789abcdef0123456789abcdef01234567",
	Sessions: []Session{{Episode: "w188i", Started: fixedNow.Add(-6 * time.Minute)}}}

func TestSurfaceModelGolden(t *testing.T) {
	golden(t, "model-rig.golden.json", BuildModel(rigIndex(t), modelNow))
}

// TestSurfaceStackModelGolden pins five stacks: w189d's step (reads, files,
// decisions, checks), the two-ancestor merge 9003555, 0ee1945 (three columns
// ahead), a first step on its brane (a2730e4) and a brane-less step
// (bf0002b). Every state's stack is checked against LayoutStack in
// TestSurfaceStackModelRules.
func TestSurfaceStackModelGolden(t *testing.T) {
	ix := rigIndex(t)
	var all []StackModel
	for _, s := range ix.States {
		switch short(s.ID) {
		case "f36428a", "9003555", "0ee1945", "a2730e4", "bf0002b":
		default:
			continue
		}
		m, err := BuildStack(ix, stateBrane(s), s.ID, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		all = append(all, m)
	}
	golden(t, "stack-rig.golden.json", all)
}

func TestSurfaceModelDeterministic(t *testing.T) {
	ix := rigIndex(t)
	a, _ := json.Marshal(BuildModel(ix, modelNow))
	shuffled := ix
	shuffled.Branes = append([]surface.Brane(nil), ix.Branes...)
	shuffled.Steps = append([]surface.Step(nil), ix.Steps...)
	shuffled.Decisions = append([]surface.Decision(nil), ix.Decisions...)
	reverse(shuffled.Branes)
	reverse(shuffled.Steps)
	reverse(shuffled.Decisions)
	for i := 0; i < 20; i++ {
		if b, _ := json.Marshal(BuildModel(shuffled, modelNow)); !bytes.Equal(a, b) {
			t.Fatalf("run %d differs", i)
		}
	}
	// A stack keeps the decisions' recorded order (LayoutStack), and ranks
	// break time ties in state order (LayoutTimeline), so only branes and
	// steps move here.
	shuffled.Decisions = ix.Decisions
	for _, s := range ix.States {
		x, err := BuildStack(ix, stateBrane(s), s.ID, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		y, err := BuildStack(shuffled, stateBrane(s), s.ID, fixedNow)
		if err != nil || !reflect.DeepEqual(x, y) {
			t.Fatalf("stack %s differs", s.ID)
		}
	}
}

// TestSurfaceModelRules checks the model against LayoutTimeline and
// LayoutStack directly, so the goldens pin values the rules produced.
func TestSurfaceModelRules(t *testing.T) {
	ix := rigIndex(t)
	m := BuildModel(ix, modelNow)
	tl := LayoutTimeline(ix, fixedNow)
	if len(m.Lanes) != len(tl.Lanes) {
		t.Fatalf("%d lanes, want %d", len(m.Lanes), len(tl.Lanes))
	}
	total := 0
	for i, l := range m.Lanes {
		want := tl.Lanes[i]
		if l.Key != want.Key || l.Glyph != want.Glyph || l.Colour != want.Colour || len(l.States) != len(want.Nodes) {
			t.Errorf("lane %d = %s %s %s", i, l.Key, l.Glyph, l.Colour)
		}
		steps := 0
		for j, s := range l.States {
			if j > 0 && s.Rank <= l.States[j-1].Rank {
				t.Errorf("lane %s out of rank order", l.Key)
			}
			if s.Kind == surface.KindStep {
				steps++
			}
			if n := node(t, tl, s.ID); n.Rank != s.Rank || n.Badge != s.Badge {
				t.Errorf("%s rank/badge %d/%d, want %d/%d", s.ID, s.Rank, s.Badge, n.Rank, n.Badge)
			}
			if !strings.HasSuffix(s.At, "Z") {
				t.Errorf("%s at %q is not UTC", s.ID, s.At)
			}
		}
		total += len(l.States)
		if l.StepsLanded != steps {
			t.Errorf("lane %s stepsLanded %d, want %d", l.Key, l.StepsLanded, steps)
		}
	}
	if total != len(ix.States) || len(m.Links) != len(tl.Links) {
		t.Fatalf("%d states, %d links", total, len(m.Links))
	}
	shared := m.Lanes[0]
	if shared.Key != SharedKey || shared.Name != "Shared world" || shared.Title != "Shared world" || shared.Row != nil || shared.Goal != SharedGoal {
		t.Errorf("shared lane = %+v", shared)
	}
	byID := map[string]ModelState{}
	for _, l := range m.Lanes {
		for _, s := range l.States {
			byID[s.ID] = s
		}
	}
	// Labels: episode, then #PR, then the short hash.
	w189d := byID["f36428abd9a72e5093b61cba14e716c37019f1ef"]
	if w189d.Label != "e1" || w189d.Graders != "4/4" || w189d.Rounds != 2 || *w189d.Minutes != 13.5 || w189d.Agent != "claude" || w189d.FirstRoundPass != nil {
		t.Errorf("e1 = %+v", w189d)
	}
	if s := byID["652c341d5072c5f49fefe29f9409c93a5119a5be"]; s.Label != "#329" || s.Graders != "" || s.Minutes != nil {
		t.Errorf("unattached step = %+v", s)
	}
	if s := byID["90035557b2012895551867c0ebe2f36b853aac09"]; s.Label != "9003555" || s.Kind != surface.KindShared {
		t.Errorf("shared state = %+v", s)
	}
	edge := m.NowEdge
	if edge.Health != "World healthy" || edge.Checked != "09:00" || edge.Daemon != "0123456" || !edge.Clean ||
		len(edge.Sessions) != 1 || edge.Sessions[0] != (ModelSession{"w188i", "6 min ago"}) || len(edge.NeedsYou) != 0 {
		t.Errorf("now edge = %+v", edge)
	}
	unstamped := BuildModel(ix, Now{At: fixedNow, SessionsUnavailable: "live sessions could not be read",
		Sessions: []Session{{Episode: "x", Started: fixedNow}}})
	if e := unstamped.NowEdge; e.Daemon != "unstamped" || e.Clean || len(e.Sessions) != 0 || e.SessionsUnavailable == "" {
		t.Errorf("unstamped now edge = %+v", e)
	}
}

// TestSurfaceModelTitles: a step's display title is its merge title's, and
// its shared state takes the same; a state with no title has none.
func TestSurfaceModelTitles(t *testing.T) {
	ix := rigIndex(t)
	step := "f36428abd9a72e5093b61cba14e716c37019f1ef"
	for i := range ix.States {
		if ix.States[i].ID == step {
			ix.States[i].Title = "feat(plugin): row 189 M5: the World plugin (episode w189d)"
		}
	}
	m := BuildModel(ix, modelNow)
	titles := map[string]string{}
	for _, l := range m.Lanes {
		for _, s := range l.States {
			titles[s.ID] = s.Title
		}
	}
	if titles[step] != "The World plugin" || titles["90035557b2012895551867c0ebe2f36b853aac09"] != "The World plugin" {
		t.Errorf("titles = %q, %q", titles[step], titles["90035557b2012895551867c0ebe2f36b853aac09"])
	}
	if titles["7fb99dae21d0416b09eed95d2da1bfe04e95aa1d"] != "" {
		t.Error("a title was invented")
	}
	s, err := BuildStack(ix, SharedKey, "90035557b2012895551867c0ebe2f36b853aac09", fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if s.Front.Title != "The World plugin" || len(s.Front.Parents) != 2 || s.Front.Parents[1].Title != "The World plugin" {
		t.Errorf("shared front = %+v", s.Front)
	}
}

func TestSurfaceModelToday(t *testing.T) {
	ix := rigIndex(t)
	// Both attached steps merged on 20 Oct UTC: e2 (codex, 1 round, 2/2,
	// 4.3 min) and e1 (claude, 2 rounds, 4/4, 13.5 min).
	now := time.Date(2026, 10, 20, 23, 59, 0, 0, time.UTC)
	got := BuildModel(ix, Now{At: now}).NowEdge.Today
	want := Today{Landed: 2, ByAgent: map[string]int{"codex": 1, "claude": 1}, FirstRound: 1, FixRounds: 1, AgentMinutes: 17.8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("today = %+v, want %+v", got, want)
	}
	// The day is now's UTC date, whatever now's zone: 21 Oct 01:00+02:00
	// is still 20 Oct, and 21 Oct 00:00 UTC has nothing.
	got = BuildModel(ix, Now{At: time.Date(2026, 10, 21, 1, 0, 0, 0, time.FixedZone("CEST", 2*3600))}).NowEdge.Today
	if !reflect.DeepEqual(got, want) {
		t.Errorf("zoned now = %+v", got)
	}
	got = BuildModel(ix, Now{At: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)}).NowEdge.Today
	if want := (Today{ByAgent: map[string]int{}}); !reflect.DeepEqual(got, want) {
		t.Errorf("next day = %+v", got)
	}
	// The day is the merge's UTC date: 21 Oct 00:30+02:00 is 20 Oct. A
	// delivered verdict short of its total is not a first-round pass.
	for i := range ix.States {
		if ix.States[i].ID == "90035557b2012895551867c0ebe2f36b853aac09" {
			ix.States[i].At = "2026-10-21T00:30:00+02:00"
		}
	}
	ix.Steps[0].Delivered = &surface.Delivered{Passed: 1, Total: 2, Regraded: []string{}}
	got = BuildModel(ix, Now{At: now}).NowEdge.Today
	if got.Landed != 2 || got.FirstRound != 0 || got.FixRounds != 1 {
		t.Errorf("edited = %+v", got)
	}
}

func TestSurfaceStackModelRules(t *testing.T) {
	ix := rigIndex(t)
	for _, x := range ix.States {
		s, err := BuildStack(ix, stateBrane(x), x.ID, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		ls, _ := LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
		ids := func(cols []StackColumn) (out []string) {
			for _, c := range cols {
				for _, card := range c.Cards {
					out = append(out, c.Brane+":"+card.ID)
				}
			}
			return out
		}
		mids := func(cols []StackModelColumn) (out []string) {
			for _, c := range cols {
				for _, card := range c.Cards {
					out = append(out, c.Brane+":"+card.ID)
				}
			}
			return out
		}
		if !reflect.DeepEqual(ids(ls.Earlier), mids(s.Earlier)) || !reflect.DeepEqual(ids(ls.Later), mids(s.Later)) || len(s.Rail) != len(ls.Rail) {
			t.Errorf("%s: stack differs from LayoutStack", x.ID)
		}
		if s.Front.ID != x.ID || s.Front.Exact != ls.Front.Exact || !strings.HasSuffix(s.Front.Position, " / "+strconv.Itoa(len(ls.Rail))) {
			t.Errorf("%s: front = %+v", x.ID, s.Front)
		}
	}
	// w189d: built on 0ee1945, merged into 9003555, reads grouped by source.
	s, err := BuildStack(ix, "row-189", "f36428abd9a72e5093b61cba14e716c37019f1ef", fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Front
	if f.BuiltOn == nil || f.BuiltOn.ID != "0ee1945ffba795812f06dce95998b646adfac343" || f.MergedInto == nil || f.MergedInto.ID != "90035557b2012895551867c0ebe2f36b853aac09" {
		t.Errorf("builtOn/mergedInto = %+v / %+v", f.BuiltOn, f.MergedInto)
	}
	if f.Name != "World onboarding" || *f.Row != 189 || f.Position != "5 / 6" || f.Task != "row-189-M5-plugin" {
		t.Errorf("front = %s %d %s %s", f.Name, *f.Row, f.Position, f.Task)
	}
	if f.Reads == nil || len(*f.Reads) != 4 || *f.ReadCalls != 8 || *f.ReadFiles != 4 {
		t.Fatalf("reads = %+v", f.Reads)
	}
	if g := (*f.Reads)[0]; g.Via != surface.ViaStep || g.Label != "#316" || g.Glyph != GlyphTriangle || g.Files[0].Count != 3 {
		t.Errorf("first read group = %+v", g)
	}
	if f.Files == nil || len(*f.Files) != 4 || *f.FilesTruncated != 0 || len(f.Decisions) != 4 {
		t.Errorf("files/decisions = %v / %v", f.Files, f.Decisions)
	}
	if c := f.Checks; c.Graders != "4/4" || c.Rounds != 2 || *c.Calls != 42 || *c.Minutes != 13.5 {
		t.Errorf("checks = %+v", c)
	}
	// A shared state has parents and no step facts; absent stays absent.
	s, err = BuildStack(ix, SharedKey, "90035557b2012895551867c0ebe2f36b853aac09", fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	f = s.Front
	if len(f.Parents) != 2 || f.BuiltOn != nil || f.Reads != nil || f.Files != nil || f.Task != "" || f.Checks != (Checks{}) || f.Row != nil {
		t.Errorf("shared front = %+v", f)
	}
	f.Parents = nil // the merged step's own facts are its, not the shared state's
	body, _ := json.Marshal(f)
	for _, k := range []string{`"reads"`, `"readCalls"`, `"files"`, `"task"`, `"builtOn"`, `"graders"`, `"firstRoundPass"`, `"grants"`, `"ranOn"`} {
		if strings.Contains(string(body), k) {
			t.Errorf("shared front has %s", k)
		}
	}
	if _, err := BuildStack(ix, "unknown", ix.States[0].ID, fixedNow); err == nil {
		t.Error("unknown brane accepted")
	}
	if _, err := BuildStack(ix, "row-189", ix.States[0].ID, fixedNow); err == nil {
		t.Error("wrong lane accepted")
	}
}
