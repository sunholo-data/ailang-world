package surfacepage

import (
	"bytes"
	"html/template"
)

// RenderStack renders the near view in the same frame and tokens as M2.
func RenderStack(s Stack) ([]byte, error) {
	var b bytes.Buffer
	err := stackTemplate.Execute(&b, struct {
		Stack
		Logo    template.HTML
		Address string
	}{s, template.HTML(logoSVG), StackAddress(s)})
	return b.Bytes(), err
}

// StackAddress is the stack's address bar text: "world KEY@SHORT", or as
// much of it as the page has.
func StackAddress(s Stack) string {
	addr := "world"
	if s.Lane.Key != "" {
		addr += " " + s.Lane.Key
		if s.Front.ID != "" {
			addr += "@" + short(s.Front.ID)
		}
	}
	return addr
}

var stackTemplate = template.Must(template.Must(template.New("stack").Parse(chromeTemplate)).Parse(`<!doctype html>
<html lang="en" data-theme="paper"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>AILANG World</title><style>
:root[data-theme="paper"] { --bg: #FFFFFF; --band: #F6F8FA; --edge: #8496A5; --ink: #2C3E50; --sub: #526677; --line: #DCE3E8; --you: #BD3012; --tint: #FDF1EE; --header: #314352; --headink: #FFFFFF; --now: #FFFFFF; --ok: #2C3E50; }
:root[data-theme="slate"] { --bg: #2C3E50; --band: #314352; --edge: #8496A5; --ink: #FFFFFF; --sub: #DCE3E8; --line: #526677; --you: #BD3012; --tint: #314352; --header: #263544; --headink: #FFFFFF; --now: #314352; --ok: #FFFFFF; }
:root[data-theme="night"] { --bg: #0F1014; --band: transparent; --edge: #3C4252; --ink: #ECEBE6; --sub: #A9A8A2; --line: #23242A; --you: #F0C06A; --tint: #221C10; --header: #121317; --headink: #ECEBE6; --now: transparent; --ok: #7FD4B0; }

* { box-sizing:border-box }
body { margin:0; background:var(--bg); color:var(--ink); font:14px/1.45 'Montserrat','Segoe UI',sans-serif }
header.top,.front header { display:flex; align-items:center; gap:12px; padding:10px 16px; background:var(--header); color:var(--headink) }
.tile { display:inline-flex; padding:3px; background:#FFFFFF; border-radius:3px }
.tile svg { width:22px; height:22px }
.name { font-size:16px; font-weight:600 }
a { color:inherit }
main { max-width:1500px; margin:auto; padding:16px 24px }
.back { display:inline-flex; align-items:center; gap:8px; margin-bottom:16px }
.glyph { width:20px; height:20px; vertical-align:middle; flex:none }
.earlier,.later { display:flex; gap:16px; justify-content:center; min-height:120px }
.column { flex:1; max-width:650px; display:flex; flex-direction:column; justify-content:flex-end; gap:5px; min-width:0 }
.card { display:flex; align-items:center; gap:8px; min-height:32px; padding:4px 12px; border:1px solid var(--edge); border-radius:3px; background:var(--band); text-decoration:none; overflow:hidden; white-space:nowrap }
.earlier .card:nth-last-child(3),.later .card:nth-child(3) { width:88%; margin:0 auto; opacity:.45 }
.earlier .card:nth-last-child(2),.later .card:nth-child(2) { width:94%; margin:0 auto; opacity:.7 }
.later .column { justify-content:flex-start }
.front { margin:12px 0; border:1px solid var(--edge); border-radius:3px; overflow:hidden }
h1 { font-size:18px; margin:0 }
.cols { display:grid; grid-template-columns:minmax(0,1.4fr) minmax(0,1fr) minmax(0,1.4fr) minmax(0,.6fr); gap:12px; padding:16px }
.col { min-width:0; border:1px solid var(--line); border-radius:3px; padding:12px; overflow-wrap:anywhere }
h2 { font-size:13px; margin:-12px -12px 12px; padding:8px 12px; background:var(--band); border-bottom:1px solid var(--line) }
ul { list-style:none; padding:0 }
li { margin:8px 0 }
.source-chip { display:inline-flex; align-items:center; gap:4px; white-space:nowrap; vertical-align:middle }
.read { display:flex; align-items:baseline; gap:4px }
.read-path { min-width:0; overflow-wrap:anywhere }
.read-detail { white-space:nowrap; flex:none }
.decision { display:flex; align-items:baseline; gap:6px }
.decision .actor { flex:none }
.decision-text { flex:1 1 auto; min-width:0; overflow-wrap:anywhere }
.decision small { flex:none; white-space:nowrap }
.actor { display:inline-block; padding:2px 7px; border:1px solid var(--edge); border-radius:3px }
[data-actor="you"] .actor { background:var(--you); color:#FFFFFF; border-color:var(--you) }
.check { display:block; font-size:24px; font-weight:600 }
small { color:var(--sub) }
.rail { display:flex; flex-wrap:wrap; gap:8px; margin-top:20px }
.rail a { padding:4px 8px; border:1px solid var(--line); border-radius:3px }
.rail [aria-current] { border-color:var(--ink); font-weight:700 }
a:focus-visible { outline:2px solid var(--you); outline-offset:2px }
@media(max-width:800px) { .cols { grid-template-columns:1fr 1fr } }
{{template "chrome-style"}}
</style><script src="/surface/surface.js" defer></script></head><body>
<header class="top"><span class="tile" aria-hidden="true">{{.Logo}}</span><span class="name">World</span>{{template "chrome-head" .Address}}</header>
<main><a class="back" href="/surface">{{template "laneGlyph" .Lane}} {{.Lane.Title}} · World timeline</a>
{{if .Reason}}<p class="reason">{{.Reason}}</p>{{else}}
<section class="earlier">{{range .Earlier}}<div class="column" data-brane="{{.Brane}}">{{range .Cards}}{{template "card" .}}{{end}}</div>{{end}}</section>
<article class="front" data-state="{{.Front.ID}}" data-brane="{{.Front.Brane}}">
<header>{{template "glyph" .Front}}<h1>{{.Front.Title}}</h1><span>{{.Lane.Title}}</span><time>{{.Front.Relative}} · {{.Front.Exact}}</time></header>
<div class="cols">
<div class="col" data-col="inputs"><h2>Inputs · came from</h2><p>{{if eq .Front.Kind "step"}}Built on{{else}}Parents{{end}}</p><ul>{{range .Inputs}}<li>{{template "target" .}}</li>{{end}}</ul>
{{if .Reads}}<h3>Read</h3>{{range .ReadSources}}<p class="read-source">{{if .Source.Href}}{{template "target" .Source}}{{else}}<span class="source-chip">{{.Label}}</span>{{end}} · {{.Files}} {{if eq .Files 1}}file{{else}}files{{end}}</p>{{end}}{{end}}
<ul>{{range .Reads}}<li class="read" data-path="{{.Path}}" data-count="{{.Count}}" data-via="{{.Via}}" data-from="{{.From}}"><span class="read-path">{{.Path}}</span><span class="read-detail">×{{.Count}} · {{if .Source.Href}}{{template "target" .Source}}{{else}}<span class="source-chip">{{.Label}}</span>{{end}}</span></li>{{end}}</ul>
{{if eq .Front.Kind "step"}}<p>{{.ReadSummary}}</p><p>Task · {{.Task}}</p><p>Agent · {{.Agent}}<br>Model · {{.Model}}</p>{{end}}</div>
<div class="col" data-col="outputs"><h2>Outputs · led to</h2><p>{{if eq .Front.Kind "step"}}Merged into{{else}}Led to{{end}}</p><ul>{{range .Outputs}}<li>{{template "target" .}}</li>{{end}}</ul>
<ul>{{range .Files}}<li class="file" data-path="{{.Path}}" data-status="{{.Status}}">{{if eq .Status "A"}}+{{else if eq .Status "M"}}~{{else}}−{{end}} {{.Path}}</li>{{end}}</ul>{{if .PR}}<p>PR #{{.PR}}</p>{{end}}</div>
<div class="col" data-col="decisions"><h2>Decisions</h2><ul>{{range .Decisions}}<li class="decision" data-actor="{{.Actor}}" data-kind="{{.Kind}}"><span class="actor">{{.ActorTitle}}</span><span class="decision-text">{{.Text}}</span><small>{{.Relative}}</small></li>{{end}}</ul></div>
<div class="col" data-col="checks"><h2>Checks</h2><p><span class="check" data-k="graders">{{.Graders}}</span>graders</p><p><span class="check" data-k="rounds">{{.Rounds}}</span>rounds</p><p><span class="check" data-k="calls">{{.Calls}}</span>calls</p><p><span class="check" data-k="minutes">{{.Minutes}}</span>agent minutes</p></div>
</div></article>
<section class="later">{{range .Later}}<div class="column" data-brane="{{.Brane}}">{{range .Cards}}{{template "card" .}}{{end}}</div>{{end}}</section>
<nav class="rail" aria-label="States on this brane">{{range .Rail}}<a data-state="{{.ID}}" href="{{.Href}}"{{if .Current}} aria-current="true"{{end}}>{{.Title}}</a>{{end}}</nav>
{{end}}</main></body></html>
{{define "glyph"}}<svg class="glyph" viewBox="0 0 20 20" aria-label="{{.Glyph}}">{{if .Points}}<polygon points="{{.Points}}" fill="{{.Colour}}"/>{{else}}<circle cx="10" cy="10" r="6" fill="none" stroke="{{.Colour}}" stroke-width="2"/>{{end}}</svg>{{end}}
{{define "laneGlyph"}}<svg class="glyph" viewBox="0 0 20 20" aria-label="{{.Glyph}}">{{if .Legend}}<polygon points="{{.Legend}}" fill="{{.Colour}}"/>{{else}}<circle cx="10" cy="10" r="6" fill="none" stroke="{{.Colour}}" stroke-width="2"/>{{end}}</svg>{{end}}
{{define "card"}}<a class="card" data-state="{{.ID}}" href="{{.Href}}"><span>{{template "glyph" .}}</span><span>{{.Title}}</span><span>{{.Relative}}</span></a>{{end}}
{{define "target"}}{{if .Href}}<a class="source-chip" href="{{.Href}}">{{template "glyph" .}}{{.Title}}</a>{{else}}{{.Title}}{{end}}{{end}}
`))
