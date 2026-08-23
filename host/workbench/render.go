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
}

type EntryView struct {
	EntryIndex     int64
	EntryHash      string
	PrevEntryHash  string
	SemanticsEpoch int64
	TransitionFn   EdgeView
	Interpreter    EdgeView
	TransitionRef  EdgeView
	WrittenBy      string
	Edges          []EdgeView
}

type TimelineView struct {
	From      int64
	Limit     int
	Entries   []EntryView
	Truncated bool
	NextHref  string
	PrevHref  string
}

type WorldView struct {
	Ref         string
	Revision    int64
	StateRoot   EdgeView
	LogHead     string
	Available   bool
	Unavailable string
}

type Page struct {
	Title    string
	World    WorldView
	Timeline TimelineView
	Selected *EntryView
	Object   *ObjectView
	Notice   string
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
body{font-family:system-ui,sans-serif;line-height:1.5;margin:0;color:#17202a;background:#f7f8fa}
header,nav,main{padding:1rem 1.5rem}header{background:#17202a;color:#fff}nav{background:#e8edf2}
main{display:grid;gap:1rem}section{background:#fff;border:1px solid #ccd4dc;border-radius:.3rem;padding:1rem}
.hash{display:inline-block;max-width:14ch;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom;white-space:nowrap}
.unavailable{font-weight:600}.payload{overflow:auto;white-space:pre-wrap}.verdict-fail{font-weight:700}.verdict-pass{font-weight:700}
dl{display:grid;grid-template-columns:max-content 1fr;gap:.25rem 1rem}dt{font-weight:700}
</style>
</head>
<body>
<header><h1>{{.Title}}</h1><p role="status">{{.Notice}}</p></header>
<nav aria-label="world browser">
<a href="{{workbenchHref ""}}">workbench</a>
{{if .World.Available}}
<a href="{{workbenchHref .World.StateRoot.Href}}" class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</a>
{{else}}<span class="unavailable" role="note">UNAVAILABLE: {{.World.Unavailable}}</span>{{end}}
</nav>
<main>
<section aria-label="timeline">
<h2>Timeline</h2>
{{if .Timeline.PrevHref}}<a href="{{workbenchHref .Timeline.PrevHref}}">previous</a>{{end}}
{{if .Timeline.NextHref}}<a href="{{workbenchHref .Timeline.NextHref}}">next</a>{{end}}
{{range .Timeline.Entries}}
<article><h3>entry {{.EntryIndex}}</h3><dl>
<dt>entry hash</dt><dd><span class="hash" title="{{.EntryHash}}" aria-label="{{.EntryHash}}">{{.EntryHash}}</span></dd>
<dt>previous entry</dt><dd><span class="hash" title="{{.PrevEntryHash}}" aria-label="{{.PrevEntryHash}}">{{.PrevEntryHash}}</span></dd>
<dt>semantics epoch</dt><dd>{{.SemanticsEpoch}}</dd><dt>written by</dt><dd>{{.WrittenBy}}</dd>
</dl></article>
{{end}}
</section>
<section aria-label="inspector">
<h2>Inspector</h2>
{{with .Object}}
<dl><dt>object</dt><dd><span class="hash" title="{{.Hash}}" aria-label="{{.Hash}}">{{.Hash}}</span></dd>
<dt>interface</dt><dd><span class="hash" title="{{.InterfaceHash}}" aria-label="{{.InterfaceHash}}">{{.InterfaceHash}}</span></dd>
<dt>semantic ID</dt><dd>{{.SemanticID}}</dd><dt>provenance</dt><dd>{{.Provenance}}</dd></dl>
{{if .Grade.Available}}<p><span>{{.Grade.Label}}</span>{{if .Grade.HasVerdict}} {{if eq .Grade.Verdict "FAIL"}}<span class="verdict-fail" aria-label="test verdict FAIL">✗ verdict: {{.Grade.Verdict}}</span>{{else}}<span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: {{.Grade.Verdict}}</span>{{end}}{{end}}</p>{{else}}<p>GRADE UNAVAILABLE — {{.Grade.Unavailable}}</p>{{end}}
{{if .PayloadShown}}<p id="payload-label">raw bytes, not interpreted HTML</p><pre class="payload" aria-labelledby="payload-label">{{.PayloadPreview}}</pre>{{if .PayloadTruncated}}<p>truncated</p>{{end}}{{end}}
{{end}}
</section>
<section aria-label="provenance walk">
<h2>Provenance walk</h2>
{{with .Object}}{{range .Edges}}<p>{{.Relation}}: {{if edgeUnavailable .}}<span class="unavailable" role="note">UNAVAILABLE: {{.Missing}}</span>{{else}}<a href="{{workbenchHref .Href}}" class="hash" title="{{.Target}}" aria-label="{{.Target}}">{{.Target}}</a>{{end}}</p>{{end}}{{end}}
</section>
</main>
</body>
</html>`

var pageTemplate = template.Must(template.New("workbench").Funcs(template.FuncMap{
	"edgeUnavailable": edgeUnavailable,
	"workbenchHref":   workbenchHref,
}).Parse(pageHTML))

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
