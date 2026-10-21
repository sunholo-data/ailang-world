package surfacepage

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// logoSVG is the AILANG mark (website/static/img/ailang-logo.svg), inlined:
// the page's CSP loads no images, but inline SVG is markup.
//
//go:embed ailang-logo.svg
var logoSVG string

// Now is the now edge's input, measured by the daemon per request.
type Now struct {
	At       time.Time // when the daemon checked
	Healthy  bool      // the daemon is serving and the index validated
	Revision string    // the daemon's vcs.revision; "" when unstamped
	Modified bool      // vcs.modified
	// Sessions are the live sessions in the credentials store, oldest first.
	// SessionsUnavailable, when set, says they could not be read.
	Sessions            []Session
	SessionsUnavailable string
}

// Session is one live session credential.
type Session struct {
	Episode string
	Started time.Time
}

type page struct {
	T       Timeline
	Logo    template.HTML
	States  int
	PlotX   int // the plot's viewBox starts here, so its x stays the layout's
	PlotW   int
	LabelX  int
	LabelW  int // the label column, sized from the lane names
	Health  string
	Count   string
	Running []string
	Note    string
}

// The label column is sized from the longest lane name at labelCharW (7.4 px
// a character at 12 px, a little over Montserrat's average), never narrower
// than PlotX. A name past maxLabelW is ellipsized; its lane's <title> keeps
// the whole name.
const (
	labelPad  = 12
	maxLabelW = 480
)

func labelCharsW(n int) int { return (n*37 + 4) / 5 }

var maxLabelRunes = (maxLabelW - LabelX - labelPad) * 5 / 37

// LabelText is a lane name as its label shows it.
func LabelText(title string) string {
	r := []rune(title)
	if len(r) <= maxLabelRunes {
		return title
	}
	return string(r[:maxLabelRunes-1]) + "…"
}

// LabelWidth is the label column's width for these lanes.
func LabelWidth(lanes []Lane) int {
	w := PlotX
	for _, l := range lanes {
		w = max(w, LabelX+labelCharsW(utf8.RuneCountInString(LabelText(l.Title)))+labelPad)
	}
	return w
}

// RenderTimeline renders the page. It refuses a timeline past
// surface.MaxStates, which Validate already bounds.
func RenderTimeline(tl Timeline, now Now) ([]byte, error) {
	p := page{T: tl, Logo: template.HTML(logoSVG), PlotX: PlotX, PlotW: tl.Width - PlotX, LabelX: LabelX, LabelW: LabelWidth(tl.Lanes), Note: now.SessionsUnavailable}
	for _, l := range tl.Lanes {
		p.States += len(l.Nodes)
	}
	if p.States > surface.MaxStates {
		return nil, fmt.Errorf("surfacepage: %d states, over the bound of %d", p.States, surface.MaxStates)
	}
	p.Health = healthLine(now)
	if now.SessionsUnavailable == "" {
		p.Count = strconv.Itoa(len(now.Sessions)) + " live sessions"
		if len(now.Sessions) == 1 {
			p.Count = "1 live session"
		}
		for _, s := range now.Sessions {
			p.Running = append(p.Running, "running "+s.Episode+" · started "+RelativeTime(s.Started, now.At))
		}
	}
	var buf bytes.Buffer
	if err := pageTemplate.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// healthLine is the now edge's positive line. "Healthy" claims only what was
// measured: the daemon answered and the index validated.
func healthLine(now Now) string {
	status := "World healthy"
	if !now.Healthy {
		status = "World serving, no index shown"
	}
	build := "daemon unstamped (no vcs.revision)"
	if now.Revision != "" {
		build = "daemon " + short(now.Revision) + ", clean"
		if now.Modified {
			build = "daemon " + short(now.Revision) + ", modified"
		}
	}
	return status + " · checked " + now.At.Format("15:04") + " · " + build
}

var pageTemplate = template.Must(template.Must(template.New("surface").Funcs(template.FuncMap{"label": LabelText}).Parse(chromeTemplate)).Parse(`<!doctype html>
<html lang="en" data-theme="paper">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>AILANG World</title>
<style>
:root[data-theme="paper"] { --bg: #FFFFFF; --band: #F6F8FA; --edge: #8496A5; --ink: #2C3E50; --sub: #526677; --line: #DCE3E8; --you: #BD3012; --tint: #FDF1EE; --header: #314352; --headink: #FFFFFF; --now: #FFFFFF; --ok: #2C3E50; }
:root[data-theme="slate"] { --bg: #2C3E50; --band: #314352; --edge: #8496A5; --ink: #FFFFFF; --sub: #DCE3E8; --line: #526677; --you: #BD3012; --tint: #314352; --header: #263544; --headink: #FFFFFF; --now: #314352; --ok: #FFFFFF; }
:root[data-theme="night"] { --bg: #0F1014; --band: transparent; --edge: #3C4252; --ink: #ECEBE6; --sub: #A9A8A2; --line: #23242A; --you: #F0C06A; --tint: #221C10; --header: #121317; --headink: #ECEBE6; --now: transparent; --ok: #7FD4B0; }
* { box-sizing: border-box; }
body { margin: 0; background: var(--bg); color: var(--ink); font: 14px/1.45 'Montserrat', 'Segoe UI', sans-serif; }
header.top { display: flex; align-items: center; gap: 10px; padding: 8px 16px; background: var(--header); color: var(--headink); }
.tile { display: inline-flex; padding: 3px; background: #FFFFFF; border-radius: 3px; }
.tile svg { display: block; width: 22px; height: 22px; }
.name { font-size: 16px; font-weight: 600; }
main { display: grid; grid-template-columns: minmax(0, 1fr) 300px; min-height: calc(100vh - 44px); }
.world { min-width: 0; padding: 12px 0; }
.frame { display: flex; align-items: flex-start; }
.labels { flex: none; display: block; }
.timeline { flex: 0 1 auto; min-width: 0; overflow-x: auto; direction: rtl; }
.timeline > svg { display: block; direction: ltr; }
.now { padding: 12px 16px; background: var(--now); border-left: 1px solid var(--edge); }
.now h2 { margin: 16px 0 4px; font-size: 12px; font-weight: 600; text-transform: uppercase; color: var(--sub); }
.health { color: var(--ok); font-weight: 600; }
.running { margin: 0; padding-left: 18px; color: var(--sub); }
.needs { padding: 6px 8px; border-radius: 3px; background: var(--tint); }
.reason { margin: 16px; color: var(--sub); }
.band { fill: var(--band); }
.rail, .tickline { stroke: var(--line); stroke-width: 1; }
.lane-label { fill: var(--ink); font-size: 12px; }
.tick { fill: var(--sub); font-size: 11px; }
.link { fill: none; stroke-width: 1.5; opacity: 0.8; }
.ring { fill: var(--bg); stroke-width: 2; }
.state:focus-visible, .state:hover { outline: none; }
.state:hover polygon, .state:hover circle { stroke: var(--ink); stroke-width: 2; }
.state:focus polygon, .state:focus circle, .state:focus-visible polygon, .state:focus-visible circle { stroke: var(--you); stroke-width: 3; }
.badge { fill: var(--you); font-size: 10px; font-weight: 700; }
.legend { display: flex; flex-wrap: wrap; gap: 4px 16px; margin: 8px 16px; padding: 0; list-style: none; color: var(--sub); font-size: 12px; }
.legend li { display: inline-flex; align-items: center; gap: 4px; }
{{template "chrome-style"}}
</style>
<script src="/surface/surface.js" defer></script>
</head>
<body>
<header class="top"><span class="tile" aria-hidden="true">{{.Logo}}</span><span class="name">World</span>{{template "chrome-head" "world"}}</header>
<main>
<section class="world" aria-label="world timeline">
{{- if .T.Reason}}
<p class="reason">{{.T.Reason}}</p>
{{- else if not .States}}
<p class="reason">The surface index holds no states yet.</p>
{{- else}}
<div class="frame">
<svg class="labels" width="{{.LabelW}}" height="{{.T.Height}}" viewBox="0 0 {{.LabelW}} {{.T.Height}}" role="img" aria-label="lanes">
<g class="bands">{{range .T.Lanes}}{{if .Band}}<rect class="band" x="0" y="{{.Top}}" width="{{$.LabelW}}" height="40"/>{{end}}{{end}}</g>
{{- range $lane := .T.Lanes}}
<g class="lane" data-brane="{{$lane.Key}}" data-glyph="{{$lane.Glyph}}"><title>{{$lane.Title}}</title>
{{- if $lane.Mark}}<polygon points="{{$lane.Mark}}" fill="{{$lane.Colour}}"/>{{else}}<circle class="ring" cx="16" cy="{{$lane.Y}}" r="6" stroke="{{$lane.Colour}}"/>{{end -}}
<text class="lane-label" x="{{$.LabelX}}" y="{{$lane.Y}}" dy="4">{{label $lane.Title}}</text></g>
{{- end}}
</svg>
<div class="timeline">
<svg class="plot" width="{{.PlotW}}" height="{{.T.Height}}" viewBox="{{.PlotX}} 0 {{.PlotW}} {{.T.Height}}" role="img" aria-label="states through time">
<g class="bands">{{range .T.Lanes}}{{if .Band}}<rect class="band" x="{{$.PlotX}}" y="{{.Top}}" width="{{$.PlotW}}" height="40"/>{{end}}{{end}}</g>
<g class="ticks">{{range .T.Ticks}}<line class="tickline" x1="{{.X}}" y1="20" x2="{{.X}}" y2="{{$.T.Height}}"/><text class="tick" x="{{.X}}" y="14">{{.Label}}</text>{{end}}</g>
{{- range .T.Lanes}}
<line class="rail" x1="{{$.PlotX}}" y1="{{.Y}}" x2="{{$.T.Width}}" y2="{{.Y}}"/>
{{- end}}
<g class="links">{{range .T.Links}}<path class="link" data-from="{{.From}}" data-to="{{.To}}" d="{{.D}}" stroke="{{.Colour}}"/>{{end}}</g>
{{- range $lane := .T.Lanes}}
<g class="nodes">
{{- range $lane.Nodes}}
<a class="state" data-state="{{.ID}}" data-kind="{{.Kind}}" data-brane="{{.Brane}}" href="{{.Href}}"><title>{{.Title}}</title>
{{- if .Points}}<polygon points="{{.Points}}" fill="{{$lane.Colour}}"/>{{else}}<circle class="ring" cx="{{.X}}" cy="{{.Y}}" r="6" stroke="{{$lane.Colour}}"/>{{end}}</a>
{{- if .Badge}}<text class="badge" data-state="{{.ID}}" x="{{.X}}" y="{{.Y}}" dx="7" dy="-8">{{.Badge}}</text>{{end}}
{{- end}}
</g>
{{- end}}
</svg>
</div>
</div>
<ul class="legend" aria-label="legend">
{{- range .T.Lanes}}
<li class="legend-item"><svg width="20" height="20" viewBox="0 0 20 20" aria-hidden="true">{{if .Legend}}<polygon points="{{.Legend}}" fill="{{.Colour}}"/>{{else}}<circle class="ring" cx="10" cy="10" r="6" stroke="{{.Colour}}"/>{{end}}</svg>{{.Title}}</li>
{{- end}}
</ul>
{{- end}}
</section>
<section class="now" aria-label="now">
<p class="health">{{.Health}}</p>
{{- if .T.Built}}
<p class="built">index built {{.T.Built}}</p>
{{- end}}
{{- if .Note}}
<p class="sessions">{{.Note}}</p>
{{- else}}
<p class="sessions">{{.Count}}</p>
{{- if .Running}}
<ul class="running">{{range .Running}}<li>{{.}}</li>{{end}}</ul>
{{- end}}
{{- end}}
<h2>Needs you</h2>
<p class="needs">Nothing needs you</p>
</section>
</main>
</body>
</html>
`))
