package surfacepage

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func renderRig(t *testing.T, now Now) string {
	t.Helper()
	body, err := RenderTimeline(LayoutTimeline(rigIndex(t), fixedNow), now)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

var rigNow = Now{At: fixedNow, Healthy: true, Revision: "0ee1945ffba795812f06dce95998b646adfac343",
	Sessions: []Session{{Episode: "w188e", Started: fixedNow.Add(-6 * time.Minute)}}}

// TestSurfaceTokensThreeThemes is TestTokensBothThemes for /surface: the same
// twelve tokens in Paper, Slate and Night, and the human's colour in each.
func TestSurfaceTokensThreeThemes(t *testing.T) {
	body := renderRig(t, rigNow)
	style := regexp.MustCompile(`(?s)<style>(.*?)</style>`).FindStringSubmatch(body)
	if len(style) != 2 {
		t.Fatal("style block missing")
	}
	roots := regexp.MustCompile(`:root\[data-theme="([a-z]+)"\]\s*\{([^}]+)\}`).FindAllStringSubmatch(style[1], -1)
	if len(roots) != 3 {
		t.Fatalf("%d theme roots, want 3", len(roots))
	}
	want := []string{"--bg", "--band", "--edge", "--ink", "--sub", "--line", "--you", "--tint", "--header", "--headink", "--now", "--ok"}
	sort.Strings(want)
	themes := []struct{ name, bg, you string }{
		{"paper", "#FFFFFF", "#BD3012"},
		{"slate", "#2C3E50", "#BD3012"},
		{"night", "#0F1014", "#F0C06A"},
	}
	for i, root := range roots {
		th := themes[i]
		if root[1] != th.name {
			t.Errorf("theme %d is %s, want %s", i, root[1], th.name)
		}
		var names []string
		for _, m := range regexp.MustCompile(`(--[a-z]+)\s*:`).FindAllStringSubmatch(root[2], -1) {
			names = append(names, m[1])
		}
		sort.Strings(names)
		if !reflect.DeepEqual(names, want) {
			t.Errorf("%s tokens=%v want %v", th.name, names, want)
		}
		for token, value := range map[string]string{"--bg": th.bg, "--you": th.you} {
			if !regexp.MustCompile(`(?i)` + token + `\s*:\s*` + value + `\s*;`).MatchString(root[2]) {
				t.Errorf("%s %s is not %s", th.name, token, value)
			}
		}
	}
	if !strings.Contains(body, `<html lang="en" data-theme="paper">`) {
		t.Error("paper is not the default theme")
	}
	if strings.Contains(body, "<link") {
		t.Fatal("external stylesheet")
	}
}

func TestRenderTimeline(t *testing.T) {
	body := renderRig(t, rigNow)
	tl := LayoutTimeline(rigIndex(t), fixedNow)

	for _, l := range tl.Lanes {
		if !strings.Contains(body, `<g class="lane" data-brane="`+l.Key+`" data-glyph="`+l.Glyph+`"><title>`+l.Title+`</title>`) {
			t.Errorf("lane %s marker missing", l.Key)
		}
		if !strings.Contains(body, "</svg>"+l.Title+"</li>") {
			t.Errorf("legend entry %s missing", l.Title)
		}
		for _, n := range l.Nodes {
			marker := `<a class="state" data-state="` + n.ID + `" data-kind="` + n.Kind + `" data-brane="` + l.Key +
				`" href="/surface?brane=` + l.Key + `&amp;at=` + n.ID + `">`
			if strings.Count(body, marker) != 1 {
				t.Errorf("state %s marker missing", n.ID)
			}
		}
	}
	if got := strings.Count(body, `<a class="state" `); got != 48 {
		t.Errorf("%d states rendered, want 48", got)
	}
	if got := strings.Count(body, `<path class="link" `); got != len(tl.Links) {
		t.Errorf("%d links rendered, want %d", got, len(tl.Links))
	}
	for _, l := range tl.Links {
		if !strings.Contains(body, `<path class="link" data-from="`+l.From+`" data-to="`+l.To+`"`) {
			t.Errorf("link %s→%s missing", l.From, l.To)
		}
	}
	badges := regexp.MustCompile(`<text class="badge" data-state="([0-9a-f]+)"[^>]*>([0-9]+)</text>`).FindAllStringSubmatch(body, -1)
	if len(badges) != 1 || badges[0][1] != "f36428abd9a72e5093b61cba14e716c37019f1ef" || badges[0][2] != "1" {
		t.Errorf("badges = %v", badges)
	}
	for _, tick := range []string{"17 Oct", "19 Oct", "20 Oct"} {
		if !regexp.MustCompile(`<text class="tick"[^>]*>` + tick + `</text>`).MatchString(body) {
			t.Errorf("tick %s missing", tick)
		}
	}
	for _, want := range []string{
		`<span class="tile" aria-hidden="true"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">`,
		`<span class="name">World</span>`,
		`<section class="now" aria-label="now">`,
		"World healthy · checked 09:00 · daemon 0ee1945, clean",
		"1 live session",
		"running w188e · started 6 min ago",
		"Nothing needs you",
		"index built 18 h ago",
		"<title>e1 · w-world-onboarding · row 189 · yesterday 08:31 · 4/4 graders · 42 calls</title>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %q", want)
		}
	}
	checkInert(t, body)
}

// TestRenderOpensAtNow: the plot's scroller is right-to-left, so it opens at
// its right end (now) with no script, and its content is left-to-right, so
// time still runs oldest left. The lane labels sit in a fixed column outside
// it, row for row; the scroller holds only the plot.
func TestRenderOpensAtNow(t *testing.T) {
	body := renderRig(t, rigNow)
	tl := LayoutTimeline(rigIndex(t), fixedNow)

	style := regexp.MustCompile(`(?s)<style>(.*?)</style>`).FindStringSubmatch(body)[1]
	rule := func(selector string) string {
		t.Helper()
		m := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(selector) + `\s*\{([^}]*)\}`).FindStringSubmatch(style)
		if m == nil {
			t.Fatalf("no %s rule", selector)
		}
		return m[1]
	}
	scroller := rule(".timeline")
	for _, decl := range []string{`overflow-x:\s*auto;`, `direction:\s*rtl;`} {
		if !regexp.MustCompile(decl).MatchString(scroller) {
			t.Errorf(".timeline lacks %s: %q", decl, scroller)
		}
	}
	if content := rule(".timeline > svg"); !regexp.MustCompile(`direction:\s*ltr;`).MatchString(content) {
		t.Errorf(".timeline content is not ltr: %q", content)
	}

	frame := regexp.MustCompile(`(?s)<div class="frame">\n(<svg class="labels".*?</svg>)\n<div class="timeline">\n(<svg class="plot".*?</svg>)\n</div>\n</div>`).FindStringSubmatch(body)
	if frame == nil {
		t.Fatal("no label column beside a .timeline scroller")
	}
	labels, plot := frame[1], frame[2]
	if strings.Count(body, `<div class="timeline">`) != 1 || strings.Contains(plot, "<div") {
		t.Fatal("the scroller is not one element holding only the plot")
	}
	for _, l := range tl.Lanes {
		marker := `<g class="lane" data-brane="` + l.Key + `" data-glyph="` + l.Glyph + `"><title>` + l.Title + `</title>`
		if !strings.Contains(labels, marker) {
			t.Errorf("lane %s marker not in the label column", l.Key)
		}
		if !strings.Contains(labels, `<text class="lane-label" x="30" y="`+strconv.Itoa(l.Y)+`" dy="4">`+l.Title+`</text>`) {
			t.Errorf("lane %s label not on its row, y=%d", l.Key, l.Y)
		}
	}
	if strings.Contains(plot, `class="lane"`) || strings.Contains(plot, "lane-label") {
		t.Error("a lane marker or label is inside the scroller")
	}
	if strings.Contains(labels, `class="state"`) {
		t.Error("a state is in the label column")
	}
	if got := strings.Count(plot, `<a class="state" `); got != 48 {
		t.Errorf("%d states inside the scroller, want 48", got)
	}
	for _, m := range []string{`<path class="link" `, `<text class="badge" `, `<text class="tick" `} {
		if strings.Count(plot, m) != strings.Count(body, m) || !strings.Contains(plot, m) {
			t.Errorf("%s not all inside the scroller", m)
		}
	}

	// Same height and the same y user units: each label lines up with its
	// lane. The plot's viewBox starts at PlotX, so its x are the layout's.
	h := strconv.Itoa(tl.Height)
	lw := strconv.Itoa(LabelWidth(tl.Lanes))
	if !strings.HasPrefix(labels, `<svg class="labels" width="`+lw+`" height="`+h+`" viewBox="0 0 `+lw+` `+h+`"`) {
		t.Errorf("label column = %.120s", labels)
	}
	w := strconv.Itoa(tl.Width - PlotX)
	if !strings.HasPrefix(plot, `<svg class="plot" width="`+w+`" height="`+h+`" viewBox="`+strconv.Itoa(PlotX)+` 0 `+w+` `+h+`"`) {
		t.Errorf("plot = %.120s", plot)
	}
	// Oldest left: the first node starts inside the plot, ticks rise left to right.
	if PlotX+8 > Left {
		t.Errorf("rank 0 at %d is clipped by the plot edge %d", Left, PlotX)
	}
	for i := 1; i < len(tl.Ticks); i++ {
		if tl.Ticks[i].X <= tl.Ticks[i-1].X {
			t.Errorf("ticks out of time order: %v", tl.Ticks)
		}
	}
}

// checkInert is the page's script rule (M4a): exactly one script element,
// /surface/surface.js with no inline body, no handler attribute, and every
// link local.
func checkInert(t *testing.T, body string) {
	t.Helper()
	if n := strings.Count(strings.ToLower(body), "<script"); n != 1 || strings.Count(body, ScriptElement) != 1 {
		t.Errorf("page has %d <script elements, want exactly one %s", n, ScriptElement)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=`).FindString(body); m != "" {
		t.Errorf("page has a handler attribute %q", m)
	}
	hrefs := regexp.MustCompile(`(?i)\bhref\s*=\s*"([^"]*)"`).FindAllStringSubmatch(body, -1)
	if len(hrefs) == 0 && strings.Contains(body, `class="state"`) {
		t.Error("no hrefs found")
	}
	for _, h := range hrefs {
		if !(strings.HasPrefix(h[1], "/") || strings.HasPrefix(h[1], "#")) || strings.HasPrefix(h[1], "//") {
			t.Errorf("href %q is not local", h[1])
		}
	}
}

func TestRenderNowEdge(t *testing.T) {
	for _, tc := range []struct {
		name string
		now  Now
		want []string
		not  []string
	}{
		{"modified, two sessions", Now{At: fixedNow, Healthy: true, Revision: "abcdef0123", Modified: true,
			Sessions: []Session{{"w1", fixedNow.Add(-2 * time.Hour)}, {"w2", fixedNow}}},
			[]string{"daemon abcdef0, modified", "2 live sessions", "running w1 · started 2 h ago", "running w2 · started just now"}, nil},
		{"unstamped, none", Now{At: fixedNow, Healthy: true},
			[]string{"daemon unstamped", "0 live sessions"}, []string{"<li>running "}},
		{"sessions unreadable", Now{At: fixedNow, Healthy: true, SessionsUnavailable: "live sessions could not be read"},
			[]string{"live sessions could not be read"}, []string{"0 live sessions"}},
		{"no index", Now{At: fixedNow}, []string{"World serving, no index shown"}, []string{"World healthy"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := renderRig(t, tc.now)
			for _, w := range tc.want {
				if !strings.Contains(body, w) {
					t.Errorf("lacks %q", w)
				}
			}
			for _, n := range tc.not {
				if strings.Contains(body, n) {
					t.Errorf("has %q", n)
				}
			}
		})
	}
}

func TestRenderReasonAndEscaping(t *testing.T) {
	body, err := RenderTimeline(Timeline{Reason: `no index <b>"here"</b>`}, Now{At: fixedNow})
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	if !strings.Contains(s, `<p class="reason">no index &lt;b&gt;&#34;here&#34;&lt;/b&gt;</p>`) || strings.Contains(s, `class="lane"`) {
		t.Errorf("reason page = %s", s)
	}
	if !strings.Contains(s, `<span class="name">World</span>`) || !strings.Contains(s, "Nothing needs you") {
		t.Error("reason page lacks the header or the now edge")
	}
	checkInert(t, s)

	none := rigIndex(t)
	none.States, none.Steps, none.Decisions = nil, nil, nil
	empty, err := RenderTimeline(LayoutTimeline(none, fixedNow), Now{At: fixedNow, Healthy: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(empty), "The surface index holds no states yet.") {
		t.Error("an index with no states has no empty line")
	}

	hostile := Now{At: fixedNow, Healthy: true, Sessions: []Session{{Episode: `<img src=x onerror=alert(1)>`, Started: fixedNow}}}
	s = renderRig(t, hostile)
	if strings.Contains(s, "<img") || !strings.Contains(s, "running &lt;img src=x") {
		t.Error("an episode name was not escaped")
	}
}

// TestRenderLabelColumnFitsNames: the label column grows to the longest
// lane name (the real index's row 114 was clipped at PlotX), and a name past
// the cap is ellipsized with the whole name in its lane's <title>.
func TestRenderLabelColumnFitsNames(t *testing.T) {
	long := "w-provenance-teeth-unseeded-baselined-walk · row 114"
	huge := strings.Repeat("w-very-long-brane-name-", 5) + "end · row 999"
	tl := LayoutTimeline(rigIndex(t), fixedNow)
	tl.Lanes[1].Title, tl.Lanes[2].Title = long, huge
	body, err := RenderTimeline(tl, rigNow)
	if err != nil {
		t.Fatal(err)
	}
	w := LabelWidth(tl.Lanes)
	if w <= PlotX || w > maxLabelW || w < LabelX+len([]rune(long))*7 {
		t.Fatalf("label column %d px for a %d-character name", w, len([]rune(long)))
	}
	if !strings.Contains(string(body), `<svg class="labels" width="`+strconv.Itoa(w)+`"`) ||
		!strings.Contains(string(body), `dy="4">`+long+`</text>`) {
		t.Error("the long name is not shown whole in a column sized for it")
	}
	short := LabelText(huge)
	if !strings.HasSuffix(short, "…") || len([]rune(short)) != maxLabelRunes ||
		!strings.Contains(string(body), `dy="4">`+short+`</text>`) || !strings.Contains(string(body), `<title>`+huge+`</title>`) {
		t.Errorf("an over-long name is not ellipsized with its title kept: %q", short)
	}
	if got := LabelWidth(LayoutTimeline(surfaceIndexOf(t, "row 1"), fixedNow).Lanes); got != PlotX {
		t.Errorf("short names widened the column to %d", got)
	}
}

func surfaceIndexOf(t *testing.T, title string) surface.Index {
	t.Helper()
	return surface.Index{Branes: []surface.Brane{{Key: "row-1", Title: title, Kind: surface.BraneAttended}}}
}

func TestRenderBound(t *testing.T) {
	big := Timeline{Lanes: []Lane{{Key: SharedKey, Nodes: make([]Node, 5001)}}}
	if _, err := RenderTimeline(big, Now{At: fixedNow}); err == nil || !strings.Contains(err.Error(), "over the bound") {
		t.Fatalf("render past MaxStates = %v", err)
	}
}
