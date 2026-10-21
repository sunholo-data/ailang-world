package surfacepage

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestSurfaceStackTruthTable(t *testing.T) {
	ix := rigIndex(t)
	state := func(id, kind, brane, base, merge string, minute int) surface.State {
		return surface.State{ID: id, Kind: kind, Brane: brane, Base: base, Merge: merge, At: fixedNow.Add(time.Duration(minute) * time.Minute).Format(time.RFC3339)}
	}
	// Hand traced through the fixture's base/merge edges. Cards are oldest
	// first above, nearest first below; own lane always precedes other lanes.
	for _, tt := range []struct {
		id            string
		before, after string
		states        []surface.State
	}{
		{"652c341", "row-189:5a8af7c,60ca0b3,f36428a;row-156:4af9b24", "shared:0ffbb9c,1998f5e,adffe4c;row-156:debb8ef", nil},
		{"f36428a", "row-189:ebe083b,5a8af7c,60ca0b3;row-156:985bee9,1bbf575,b0f989f;row-198:e42f030", "row-189:652c341;shared:9003555,0ffbb9c,1998f5e", nil},
		{"0ee1945", "shared:0adeb70,96d4901,fe2cc96;row-198:de79d98,e42f030", "shared:7fb99da,9003555,0ffbb9c;row-156:4af9b24,debb8ef;row-189:f36428a,652c341", nil},
		{"a2730e4", "shared:05efc0c", "shared:a0c7305,02bd141,78f2860", nil},
		{"9003555", "shared:fe2cc96,0ee1945,7fb99da;row-189:5a8af7c,60ca0b3,f36428a", "shared:0ffbb9c,1998f5e,49ac0d1;row-189:652c341", nil},
		{"sharedx", "shared:older00,nobrane", "", []surface.State{
			state("older00", surface.KindShared, "", "", "", 0),
			state("nobrane", surface.KindStep, "", "", "sharedx", 1),
			state("sharedx", surface.KindShared, "", "", "", 2),
		}},
		{"windowx", "other:other02,other03,other04", "", []surface.State{
			state("other01", surface.KindStep, "other", "", "base000", 0),
			state("other02", surface.KindStep, "other", "", "base000", 1),
			state("other03", surface.KindStep, "other", "", "base000", 2),
			state("other04", surface.KindStep, "other", "", "base000", 3),
			state("base000", surface.KindShared, "", "", "", 4),
			state("windowx", surface.KindStep, "own", "base000", "", 5),
		}},
		{"leaving", "", "shared:merge00,after00", []surface.State{
			state("leaving", surface.KindStep, "own", "", "merge00", 0),
			state("merge00", surface.KindShared, "", "", "", 1),
			state("nobrane", surface.KindStep, "", "merge00", "", 2),
			state("after00", surface.KindShared, "", "", "", 3),
		}},
	} {
		t.Run(tt.id, func(t *testing.T) {
			ix := ix
			if tt.states != nil {
				ix = surface.Index{States: tt.states, Branes: []surface.Brane{{Key: "own"}, {Key: "other"}}}
			}
			var x surface.State
			for _, s := range ix.States {
				if strings.HasPrefix(s.ID, tt.id) {
					x = s
				}
			}
			s, err := LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range []struct {
				got  []StackColumn
				want string
			}{{s.Earlier, tt.before}, {s.Later, tt.after}} {
				var cols []string
				for _, c := range v.got {
					var ids []string
					for _, card := range c.Cards {
						ids = append(ids, short(card.ID))
					}
					cols = append(cols, c.Brane+":"+strings.Join(ids, ","))
				}
				if got := strings.Join(cols, ";"); got != v.want {
					t.Errorf("got %s; want %s", got, v.want)
				}
			}
		})
	}
}

func TestSurfaceStackDeterminismAndRender(t *testing.T) {
	ix := rigIndex(t)
	for _, x := range ix.States {
		s, err := LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		shuffled := ix
		shuffled.States = append([]surface.State(nil), ix.States...)
		reverse(shuffled.States)
		b, err := LayoutStack(shuffled, stateBrane(x), x.ID, fixedNow)
		if err != nil || !reflect.DeepEqual(s, b) {
			t.Fatalf("not deterministic: %s", x.ID)
		}
		html, err := RenderStack(s)
		if err != nil {
			t.Fatal(err)
		}
		body := string(html)
		for _, col := range []string{"inputs", "outputs", "decisions", "checks"} {
			if !strings.Contains(body, `data-col="`+col+`"`) {
				t.Error("missing column", col)
			}
		}
		for _, k := range []string{"graders", "rounds", "calls", "minutes"} {
			if !strings.Contains(body, `class="check" data-k="`+k+`"`) {
				t.Error("missing check", k)
			}
		}
		if strings.Count(body, `aria-current="true"`) != 1 {
			t.Error("current rail state")
		}
		checkInert(t, body)
		for _, m := range regexp.MustCompile(`href="([^"]*)"`).FindAllStringSubmatch(body, -1) {
			if !strings.HasPrefix(m[1], "/surface") {
				t.Error("nonlocal href", m[1])
			}
		}
		for _, read := range s.Reads {
			if !strings.Contains(body, `class="read" data-path="`+read.Path+`"`) {
				t.Error("missing read")
			}
		}
		for _, file := range s.Files {
			if !strings.Contains(body, `class="file" data-path="`+file.Path+`" data-status="`+file.Status+`"`) {
				t.Error("missing file", file.Path)
			}
		}
		for _, decision := range s.Decisions {
			if !strings.Contains(body, `class="decision" data-actor="`+decision.Actor+`" data-kind="`+decision.Kind+`"`) {
				t.Error("missing decision")
			}
		}
		for _, cols := range [][]StackColumn{s.Earlier, s.Later} {
			seen := map[string]bool{}
			for _, col := range cols {
				if seen[col.Brane] {
					t.Errorf("state %s: duplicate column %s", x.ID, col.Brane)
				}
				seen[col.Brane] = true
				if len(col.Cards) > 3 {
					t.Error("column exceeds three cards")
				}
				for _, c := range col.Cards {
					if !strings.Contains(body, `class="card" data-state="`+c.ID+`"`) {
						t.Error("missing stack card")
					}
				}
			}
		}
		for _, column := range regexp.MustCompile(`(?s)<div class="column"[^>]*>(.*?)</div>`).FindAllStringSubmatch(body, -1) {
			if strings.Contains(column[1], "<div") {
				t.Error("nested div in column")
			}
		}
	}
	if _, err := LayoutStack(ix, "unknown", ix.States[0].ID, fixedNow); err == nil {
		t.Error("unknown brane accepted")
	}
	if _, err := LayoutStack(ix, "row-189", ix.States[0].ID, fixedNow); err == nil {
		t.Error("wrong lane accepted")
	}
}

func TestSurfaceStackMissingFacts(t *testing.T) {
	ix := rigIndex(t)
	// A missing base gives no cross-brane window; no rig facts means every
	// check is unknown rather than zero. The own-brane history still exists.
	x := ix.States[0]
	for _, s := range ix.States {
		if strings.HasPrefix(s.ID, "652c341") {
			x = s
		}
	}
	for i := range ix.States {
		if ix.States[i].ID == x.ID {
			ix.States[i].Base = "outside"
		}
	}
	ix.Steps = nil
	ix.Reads = nil
	s, err := LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Earlier) != 1 || s.Earlier[0].Brane != stateBrane(x) {
		t.Fatal("missing base made a window")
	}
	if s.Graders != "—" || s.Rounds != "—" || s.Calls != "—" || s.Minutes != "—" {
		t.Fatal("invented a missing check")
	}
	// Missing previous base falls back to the current base alone.
	ix = rigIndex(t)
	for i := range ix.States {
		if strings.HasPrefix(ix.States[i].ID, "f36428a") {
			ix.States[i].Base = "outside"
		}
	}
	s, err = LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Earlier) != 1 {
		t.Fatal("previous missing base did not restrict the window")
	}
}

func TestSurfaceStackW189eInlineReads(t *testing.T) {
	ix := rigIndex(t)
	// The compact rig attaches e1 to w189d's state and omits w189e's
	// report. Extend this local copy; retain the structure fixture unchanged.
	var previous string
	for i := range ix.Steps {
		if strings.HasPrefix(ix.Steps[i].State, "f36428a") {
			previous = ix.Steps[i].State
			ix.Steps[i].Episode = "w189d"
		}
	}
	for _, state := range ix.States {
		if strings.HasPrefix(state.ID, "652c341") {
			ix.Steps = append(ix.Steps, surface.Step{State: state.ID, Episode: "w189e"})
			ix.Reads[state.ID] = append([]surface.Read{
				{Path: "CLAUDE.md", Count: 1, Via: surface.ViaStep, From: previous},
			}, ix.Reads[previous]...)
		}
	}
	var id string
	for _, step := range ix.Steps {
		if step.Episode == "w189e" {
			id = step.State
		}
	}
	if id == "" {
		t.Fatal("missing w189e in index-rig.json")
	}
	stack, err := LayoutStack(ix, "row-189", id, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	html, err := RenderStack(stack)
	if err != nil {
		t.Fatal(err)
	}
	body := string(html)
	reads := regexp.MustCompile(`(?s)<li class="read"[^>]*>(.*?)</li>`).FindAllStringSubmatch(body, -1)
	if len(reads) != len(ix.Reads[id]) {
		t.Fatal("read count changed")
	}
	for i, read := range reads {
		if !strings.Contains(read[0], `data-path="`+ix.Reads[id][i].Path+`"`) {
			t.Error("read order changed")
		}
		if !strings.Contains(read[1], `class="source-chip"`) ||
			regexp.MustCompile(`<(div|p|ul|li|h[1-6]|br)\b`).MatchString(read[1]) {
			t.Error("read chip is missing or contains block markup", read[0])
		}
	}
	if !strings.Contains(body, "<h3>Read</h3>") {
		t.Error("missing Read heading")
	}
	var keys []string
	counts := map[string]int{}
	for _, read := range ix.Reads[id] {
		key := read.Via + ":" + read.From
		if counts[key] == 0 {
			keys = append(keys, key)
		}
		counts[key]++
	}
	if len(stack.ReadSources) != len(keys) {
		t.Fatal("source summary count")
	}
	for i, key := range keys {
		if stack.ReadSources[i].Files != counts[key] {
			t.Error("source file count")
		}
	}
	summaries := regexp.MustCompile(`(?s)<p class="read-source">(.*?)</p>`).FindAllStringSubmatch(body, -1)
	if len(summaries) != len(keys) {
		t.Error("missing source summary lines")
	}
	for _, step := range ix.Steps {
		if step.Episode != "w189d" {
			continue
		}
		for _, pattern := range []string{
			`(?s)<a class="card" data-state="` + step.State + `"[^>]*>.*?<span>w189d</span>`,
			`<a class="source-chip" href="[^"]*` + step.State + `">[^<]*<svg.*?</svg>w189d</a>`,
		} {
			if !regexp.MustCompile(pattern).MatchString(body) {
				t.Error("missing attached episode label", pattern)
			}
		}
	}
	// A decision's pill, text and time share one flex line: only the text
	// wraps, inside its own box, never away from its pill.
	if !strings.Contains(body, ".source-chip { display:inline-flex;") ||
		!strings.Contains(body, ".decision { display:flex; align-items:baseline;") ||
		!strings.Contains(body, ".decision .actor { flex:none }") ||
		!strings.Contains(body, ".decision-text { flex:1 1 auto; min-width:0; overflow-wrap:anywhere }") ||
		!strings.Contains(body, ".decision small { flex:none; white-space:nowrap }") {
		t.Error("missing inline layout rules")
	}
	for _, d := range regexp.MustCompile(`<li class="decision"[^>]*>(.*?)</li>`).FindAllStringSubmatch(body, -1) {
		if !regexp.MustCompile(`^<span class="actor">[^<]*</span><span class="decision-text">[^<]*</span><small>[^<]*</small>$`).MatchString(d[1]) {
			t.Errorf("decision is not pill, text, time: %s", d[1])
		}
	}
}
