package workbench

import (
	"errors"
	"html/template"
	"io"
)

const WorkbenchPageLimit = 100
const MaxPayloadPreview = 64 << 10 // 64 KiB

type GradeLabel string

const (
	GradePROVEN   GradeLabel = "PROVEN"
	GradeTESTED   GradeLabel = "TESTED"
	GradeATTESTED GradeLabel = "ATTESTED"
	GradeCLAIMED  GradeLabel = "CLAIMED"
)

type Verdict string

const (
	VerdictPass Verdict = "PASS"
	VerdictFail Verdict = "FAIL"
)

const KindTestReport = "TestReport"

type GradeView struct {
	Available   bool
	Label       GradeLabel
	HasVerdict  bool
	Verdict     Verdict
	Unavailable string
}

type EdgeView struct {
	Relation  string
	Available bool
	Target    string
	Href      string
	Missing   string
}

type ObjectView struct {
	Hash             string
	InterfaceHash    string
	SemanticID       string
	Provenance       string
	PayloadShown     bool
	PayloadPreview   string
	PayloadTruncated bool
	Grade            GradeView
	Edges            []EdgeView
	Commits          *CommitView
	References       *ReferenceView
}

type CommitView struct {
	Edges     []EdgeView
	Truncated bool
	Continued bool
	NextHref  string
	FirstHref string
}

type ReferenceView struct {
	Edges     []EdgeView
	Truncated bool
	Continued bool
	NextHref  string
	FirstHref string
}

type EntryView struct {
	EntryIndex     int64
	EntryHash      string
	PrevEntryHash  string
	SemanticsEpoch int64
	WrittenBy      string
	SelectHref     string
	Edges          []EdgeView
}

type TimelineView struct {
	Entries  []EntryView
	NextHref string
	PrevHref string
}

type WorldView struct {
	Ref         string
	Revision    int64
	StateRoot   EdgeView
	LogHead     string
	Available   bool
	Unavailable string
}

// LiveView is the newest-first log position at render time.
type LiveView struct {
	Cursor int64
	Head   string
	Recent []EntryView
}

type DecisionRow struct {
	RequestRef string
	Href       string
	Effect     string
	Scope      string
	Requester  string
	Cost       int64
	Status     string
	DecidedBy  string
}
type DecisionsView struct {
	Rows        []DecisionRow
	Cursor      int64
	Truncated   bool
	Unavailable string
}

type Page struct {
	Decisions DecisionsView
	Graph     GraphView
	Live      LiveView
	Title     string
	World     WorldView
	Timeline  TimelineView
	Selected  *EntryView
	Object    *ObjectView
	Notice    string
}

var ErrInvalidGrade = errors.New("workbench: grade label is not one of PROVEN|TESTED|ATTESTED|CLAIMED")
var ErrMissingVerdict = errors.New("workbench: TestReport grade requires a verdict")

func workbenchHref(query string) string {
	if query != "" && query[0] == '?' {
		return "/workbench" + query
	}
	return "/workbench"
}

func edgeUnavailable(edge EdgeView) bool {
	if !edge.Available {
		return true
	}
	return false
}

const pageHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}}</title>
<style>
:root{--bg:#faf9f5;--card:#fff;--ink:#1a1a17;--mut:#6b6a64;--line:#dedcd3;--acc:#185fa5;--accbg:#e6f1fb;--ok:#0f6e56;--okbg:#e1f5ee;--warn:#854f0b;--warnbg:#faeeda}
@media (prefers-color-scheme: dark){:root{--bg:#1f1e1b;--card:#292824;--ink:#eceae3;--mut:#a3a199;--line:#3d3c36;--acc:#85b7eb;--accbg:#0c447c;--ok:#5dcaa5;--okbg:#085041;--warn:#fac775;--warnbg:#633806}}
body{font-family:system-ui,sans-serif;line-height:1.5;margin:0;color:var(--ink);background:var(--bg)}
header,nav,main,footer{padding:1rem 1.5rem}header{border-bottom:1px solid var(--line)}nav{background:var(--accbg)}
main{display:grid;gap:1rem;max-width:80rem;margin:auto}section{background:var(--card);border:1px solid var(--line);border-radius:.5rem;padding:1rem}
a{color:var(--acc)}h1,h2,h3{line-height:1.2}footer{color:var(--mut)}
.hash{display:inline-block;max-width:14ch;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom;white-space:nowrap}
.unavailable{font-weight:600;color:var(--warn)}.payload{overflow:auto;white-space:pre-wrap}.verdict-fail{font-weight:700}.verdict-pass{font-weight:700}
dl{display:grid;grid-template-columns:max-content 1fr;gap:.25rem 1rem}dt{font-weight:700}
.live-list{list-style:none;padding:0}.live-list li{padding:.4rem 0;border-bottom:1px solid var(--line)}
.fresh{animation:fresh 1.5s ease-out}@keyframes fresh{from{background:var(--okbg)}to{background:transparent}}
</style>
<script src="/workbench/live.js" defer></script>
</head>
<body>
<header><h1>{{.Title}}</h1><p role="status">{{.Notice}}</p></header>
<nav aria-label="world browser">
<a href="{{workbenchHref ""}}">workbench</a>
{{if .World.Available}}
<dl>
<dt>world</dt><dd><span class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</span></dd>
<dt>revision</dt><dd>{{.World.Revision}}</dd>
<dt>log head</dt><dd><span class="hash" title="{{.World.LogHead}}" aria-label="{{.World.LogHead}}">{{.World.LogHead}}</span></dd>
</dl>
{{template "edge" .World.StateRoot}}
{{else}}<span class="unavailable" role="note">UNAVAILABLE: {{.World.Unavailable}}</span>{{end}}
</nav>
<main>
<section aria-label="timeline">
<h2>Timeline</h2>
{{if .Timeline.PrevHref}}<a href="{{workbenchHref .Timeline.PrevHref}}">previous</a>{{end}}
{{if .Timeline.NextHref}}<a href="{{workbenchHref .Timeline.NextHref}}">next</a>{{end}}
{{range .Timeline.Entries}}
<article><h3>entry {{.EntryIndex}}</h3><dl>
{{template "entryFields" .}}
</dl><p><a href="{{workbenchHref .SelectHref}}">select entry {{.EntryIndex}}</a></p></article>
{{end}}
</section>
<section aria-label="inspector">
<h2>Inspector</h2>
{{with .Selected}}<article aria-label="selected entry"><h3>selected entry {{.EntryIndex}}</h3><dl>
{{template "entryFields" .}}
</dl>{{range .Edges}}{{template "edge" .}}{{end}}</article>{{end}}
{{with .Object}}
<dl><dt>object</dt><dd><span class="hash" title="{{.Hash}}" aria-label="{{.Hash}}">{{.Hash}}</span></dd>
<dt>interface</dt><dd><span class="hash" title="{{.InterfaceHash}}" aria-label="{{.InterfaceHash}}">{{.InterfaceHash}}</span></dd>
<dt>semantic ID</dt><dd>{{.SemanticID}}</dd><dt>provenance</dt><dd>{{.Provenance}}</dd></dl>
{{if .Grade.Available}}<p><span>{{.Grade.Label}}</span>{{if .Grade.HasVerdict}} {{if eq .Grade.Verdict "FAIL"}}<span class="verdict-fail" aria-label="test verdict FAIL">✗ verdict: {{.Grade.Verdict}}</span>{{else}}<span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: {{.Grade.Verdict}}</span>{{end}}{{end}}</p>{{else}}<p>GRADE UNAVAILABLE — {{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}</p>{{end}}
{{if .PayloadShown}}<p id="payload-label">raw bytes, not interpreted HTML</p><pre class="payload" aria-labelledby="payload-label">{{.PayloadPreview}}</pre>{{if .PayloadTruncated}}<p>truncated</p>{{end}}{{end}}
{{end}}
</section>
<section aria-label="live" aria-live="polite" data-live-cursor="{{.Live.Cursor}}">
<h2>Live log</h2>
{{if .Live.Recent}}<p>Log at entry {{.Live.Cursor}} (head {{.Live.Head}}) when this page was rendered</p>
<ul class="live-list">{{range .Live.Recent}}<li><a href="{{workbenchHref .SelectHref}}">select entry {{.EntryIndex}}</a> <span class="hash" title="{{.EntryHash}}">{{.EntryHash}}</span></li>{{end}}</ul>
{{else}}<p>The log is empty</p>{{end}}
</section>
<section aria-label="world graph"><h2>World graph</h2>
{{if .Graph.Unavailable}}<p>UNAVAILABLE: {{.Graph.Unavailable}}</p>{{end}}{{if .Graph.Nodes}}{{else}}<p>graph: select an entry</p>{{end}}
<svg role="img" aria-label="checked entry edges" viewBox="0 0 640 {{.Graph.Height}}">
{{range .Graph.Edges}}<line x1="{{.X1}}" y1="{{.Y1}}" x2="{{.X2}}" y2="{{.Y2}}" stroke="currentColor"><title>{{.Relation}}</title></line>{{end}}
{{range .Graph.Nodes}}<g>{{if .Available}}<a href="{{workbenchHref .Href}}">{{end}}<rect x="{{.X}}" y="{{.Y}}" width="280" height="28" fill="none" stroke="currentColor"{{if .Available}}{{else}} stroke-dasharray="4 3"{{end}}></rect><text x="{{.X}}" y="{{.Y}}" dy="20">{{if .Available}}{{else}}UNAVAILABLE: {{end}}{{.Label}}</text>{{if .Available}}</a>{{end}}</g>{{end}}
</svg></section>
<section aria-label="decisions"><h2>Decisions</h2>
<p>as of entry {{.Decisions.Cursor}}</p>
{{if .Decisions.Unavailable}}<p>UNAVAILABLE: {{.Decisions.Unavailable}}</p>{{else}}
{{range .Decisions.Rows}}<article><p><a href="{{workbenchHref .Href}}" class="hash" title="{{.RequestRef}}">{{.RequestRef}}</a></p><p>{{.Effect}}; scope: {{.Scope}}; requester: {{.Requester}}; cost: {{.Cost}}</p><p>{{.Status}}{{if .DecidedBy}} by {{.DecidedBy}}{{end}}</p></article>{{else}}<p>No approval requests recorded</p>{{end}}
{{if .Decisions.Truncated}}<p>Showing recent approvals; more recorded</p>{{end}}{{end}}
</section>
<section aria-label="provenance walk">
<h2>Provenance walk</h2>
{{with .Object}}{{range .Edges}}{{template "edge" .}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>{{end}}{{with .Commits}}<section aria-label="committedBy"><h3>committedBy</h3><p>Commits whose object set carried this object, oldest first. An object stored by PutObject or the journal before a commit carried it is attributed only to the commits that carried it.</p>{{range .Edges}}{{template "edge" .}}{{else}}{{if .Continued}}<p>no further commits carried this object</p>{{else}}<p>no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.</p>{{end}}{{end}}{{if .Truncated}}<p>Showing 100 commits; more recorded</p>{{end}}{{if .Continued}}<p><a href="{{workbenchHref .FirstHref}}">first page</a></p>{{end}}{{if .NextHref}}<a href="{{workbenchHref .NextHref}}">next commits</a>{{end}}</section>{{end}}{{with .References}}<section aria-label="referencedBy"><h3>referencedBy</h3><p>Entries/worlds only: transitionRef, transitionFn, interpreter, stateRoot. Registry, journal and object-interface inbound references are not included.</p>{{range .Edges}}{{template "edge" .}}{{else}}{{if .Continued}}<p>no further references recorded in these entry/world fields</p>{{else}}<p>none recorded in these entry/world fields</p>{{end}}{{end}}{{if .Truncated}}<p>Showing 100 references; more recorded</p>{{end}}{{if .Continued}}<p>New references before this cursor require restarting the walk. <a href="{{workbenchHref .FirstHref}}">first page</a></p>{{end}}{{if .NextHref}}<a href="{{workbenchHref .NextHref}}">next references</a>{{end}}</section>{{end}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>{{end}}
</section>
</main>
<footer aria-label="live status" role="status"><span data-live-status>Live updates: off. Reload to refresh.</span></footer>
</body>
</html>`

const partialsHTML = `{{define "edge"}}<p>{{.Relation}}: {{if edgeUnavailable .}}<span class="unavailable" role="note">UNAVAILABLE: {{.Missing}}</span>{{else}}<a href="{{workbenchHref .Href}}" class="hash" title="{{.Target}}" aria-label="{{.Target}}">{{.Target}}</a>{{end}}</p>{{end}}
{{define "entryFields"}}<dt>entry hash</dt><dd><span class="hash" title="{{.EntryHash}}" aria-label="{{.EntryHash}}">{{.EntryHash}}</span></dd>
<dt>previous entry</dt><dd><span class="hash" title="{{.PrevEntryHash}}" aria-label="{{.PrevEntryHash}}">{{.PrevEntryHash}}</span></dd>
<dt>semantics epoch</dt><dd>{{.SemanticsEpoch}}</dd><dt>written by</dt><dd>{{.WrittenBy}}</dd>{{end}}`

var pageTemplate = template.Must(template.Must(template.New("workbench").Funcs(template.FuncMap{
	"edgeUnavailable": edgeUnavailable,
	"workbenchHref":   workbenchHref,
}).Parse(pageHTML)).Parse(partialsHTML))

func Render(w io.Writer, p Page) error { return pageTemplate.Execute(w, p) }

func validGrade(l GradeLabel) bool {
	switch l {
	case GradePROVEN, GradeTESTED, GradeATTESTED, GradeCLAIMED:
		return true
	default:
		return false
	}
}

func NewGradeUnavailable(reason string) GradeView {
	return GradeView{Available: false, Unavailable: reason}
}

func NewGradeView(label GradeLabel, kind string, verdict *Verdict) (GradeView, error) {
	if !validGrade(label) {
		return GradeView{}, ErrInvalidGrade
	}
	if kind == KindTestReport && verdict == nil {
		return GradeView{}, ErrMissingVerdict
	}

	view := GradeView{
		Available:  true,
		Label:      label,
		HasVerdict: verdict != nil,
	}
	if verdict != nil {
		view.Verdict = *verdict
	}
	return view, nil
}
