package workbench

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	"text/template/parse"
)

func TestGradeViewRejectsUnavailableOrInvalidInput(t *testing.T) {
	t.Run("accepts-four-valid-labels", func(t *testing.T) {
		pass := VerdictPass
		labels := []GradeLabel{GradePROVEN, GradeTESTED, GradeATTESTED, GradeCLAIMED}
		for _, label := range labels {
			kind := ""
			var verdict *Verdict
			if label == GradeTESTED {
				kind = KindTestReport
				verdict = &pass
			}
			view, err := NewGradeView(label, kind, verdict)
			if err != nil {
				t.Fatalf("NewGradeView(%q): %v", label, err)
			}
			if !view.Available || view.Label != label {
				t.Fatalf("NewGradeView(%q) = %+v", label, view)
			}
		}
	})

	t.Run("rejects-invalid-label", func(t *testing.T) {
		for _, label := range []GradeLabel{"UNSUPPORTED", "proven", "", "CLAIMED "} {
			_, err := NewGradeView(label, "", nil)
			if !errors.Is(err, ErrInvalidGrade) {
				t.Errorf("NewGradeView(%q) error = %v, want ErrInvalidGrade", label, err)
			}
		}
	})

	t.Run("unavailable-is-not-claimed", func(t *testing.T) {
		view := NewGradeUnavailable("no canonical host projection")
		if view.Label == GradeCLAIMED {
			t.Error("unavailable grade was downgraded to CLAIMED")
		}
		if view.Available {
			t.Error("unavailable grade reports Available=true")
		}
	})

	t.Run("no-proven-inference", func(t *testing.T) {
		view := NewGradeUnavailable("no canonical host projection")
		if view.Label == GradePROVEN {
			t.Error("unavailable grade inferred PROVEN")
		}
	})
}

func TestGradeViewRequiresTestVerdict(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		_, err := NewGradeView(GradeTESTED, KindTestReport, nil)
		if !errors.Is(err, ErrMissingVerdict) {
			t.Fatalf("error = %v, want ErrMissingVerdict", err)
		}
	})

	t.Run("fail", func(t *testing.T) {
		fail := VerdictFail
		view, err := NewGradeView(GradeTESTED, KindTestReport, &fail)
		if err != nil {
			t.Fatal(err)
		}
		if view.Label != GradeTESTED || view.Verdict != VerdictFail || !view.HasVerdict {
			t.Fatalf("view = %+v", view)
		}
		var body bytes.Buffer
		if err := Render(&body, Page{Title: "test", Object: &ObjectView{Grade: view}}); err != nil {
			t.Fatal(err)
		}
		if rendered := body.String(); !strings.Contains(rendered, "TESTED") || !strings.Contains(rendered, "verdict: FAIL") || !strings.Contains(rendered, `aria-label="test verdict FAIL"`) {
			t.Fatalf("rendered test grade and verdict = %q", rendered)
		}
		want := `<p><span>TESTED</span> <span class="verdict-fail" aria-label="test verdict FAIL">✗ verdict: FAIL</span></p>`
		if rendered := body.String(); !strings.Contains(rendered, want) || strings.Contains(rendered, `aria-label="test verdict PASS"`) {
			t.Fatalf("rendered FAIL verdict, want %q in %q", want, rendered)
		}
	})

	t.Run("pass", func(t *testing.T) {
		pass := VerdictPass
		view, err := NewGradeView(GradeTESTED, KindTestReport, &pass)
		if err != nil {
			t.Fatal(err)
		}
		if view.Label != GradeTESTED || view.Verdict != VerdictPass || !view.HasVerdict {
			t.Fatalf("view = %+v", view)
		}
		var body bytes.Buffer
		if err := Render(&body, Page{Title: "test", Object: &ObjectView{Grade: view}}); err != nil {
			t.Fatal(err)
		}
		want := `<p><span>TESTED</span> <span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: PASS</span></p>`
		if rendered := body.String(); !strings.Contains(rendered, want) || strings.Contains(rendered, `aria-label="test verdict FAIL"`) {
			t.Fatalf("rendered PASS verdict, want %q in %q", want, rendered)
		}
	})
}

func TestRenderGradeWithoutVerdictClaim(t *testing.T) {
	proven, err := NewGradeView(GradePROVEN, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		grade GradeView
		want  string
	}{
		{"no-verdict", proven, `<p><span>PROVEN</span></p>`},
		{"unavailable", NewGradeUnavailable("no canonical host projection"), `<p>GRADE UNAVAILABLE — no canonical host projection</p>`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			if err := Render(&body, Page{Title: "test", Object: &ObjectView{Grade: tc.grade}}); err != nil {
				t.Fatal(err)
			}
			rendered := body.String()
			if !strings.Contains(rendered, tc.want) {
				t.Fatalf("rendered grade, want %q in %q", tc.want, rendered)
			}
			for _, claim := range []string{`class="verdict-`, `aria-label="test verdict`} {
				if strings.Contains(rendered, claim) {
					t.Fatalf("grade without a verdict rendered a verdict claim %q in %q", claim, rendered)
				}
			}
		})
	}
}

func TestWorkbenchHrefAppendsOnlyQueryStrings(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", "/workbench"},
		{"?object=abc", "/workbench?object=abc"},
		{"object=abc", "/workbench"},
		{"//evil.example/x", "/workbench"},
		{"javascript:alert(1)", "/workbench"},
	} {
		if got := workbenchHref(tc.in); got != tc.want {
			t.Errorf("workbenchHref(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRenderEscapesAllObjectText(t *testing.T) {
	hostile := `<script>alert(1)</script>" onmouseover=x&`
	p := Page{
		Title:  "SAFE-LABEL-9f2",
		Notice: hostile,
		Timeline: TimelineView{Entries: []EntryView{{
			WrittenBy: hostile,
		}}},
		Object: &ObjectView{
			SemanticID:     hostile,
			Provenance:     hostile,
			PayloadShown:   true,
			PayloadPreview: hostile,
			Grade:          NewGradeUnavailable(hostile),
			Edges: []EdgeView{
				{Available: true, Target: hostile, Href: "?object=safe"},
				{Missing: hostile},
			},
		},
	}
	var body bytes.Buffer
	if err := Render(&body, p); err != nil {
		t.Fatal(err)
	}
	rendered := body.String()
	if !strings.Contains(rendered, "&lt;script&gt;") {
		t.Fatalf("escaped script marker missing from %q", rendered)
	}
	if strings.Contains(rendered, "<script>alert") {
		t.Fatalf("hostile script rendered as markup: %q", rendered)
	}
	if !strings.Contains(rendered, "SAFE-LABEL-9f2") {
		t.Fatalf("positive control missing from %q", rendered)
	}
}

func TestRenderEmitsOnlyLocalLinks(t *testing.T) {
	p := Page{
		Title: "links",
		Timeline: TimelineView{
			PrevHref: "?from=0",
			NextHref: "?from=100",
		},
		Object: &ObjectView{Edges: []EdgeView{{Available: true, Target: "object", Href: "?object=abc"}}},
	}
	var body bytes.Buffer
	if err := Render(&body, p); err != nil {
		t.Fatal(err)
	}
	rendered := body.String()
	matches := regexp.MustCompile(`href="([^"]*)"`).FindAllStringSubmatch(rendered, -1)
	if len(matches) == 0 {
		t.Fatal("rendered page contains no links")
	}
	for _, match := range matches {
		if match[1] != "/workbench" && !strings.HasPrefix(match[1], "/workbench?") {
			t.Errorf("non-local workbench link %q", match[1])
		}
	}
	for _, forbidden := range []string{"http://", "https://", `href="//`, "data:", "javascript:"} {
		if strings.Contains(rendered, forbidden) {
			t.Errorf("rendered page contains forbidden link form %q", forbidden)
		}
	}
}

func TestRenderUnavailableProvenanceEdge(t *testing.T) {
	p := Page{
		Title: "edges",
		Object: &ObjectView{Edges: []EdgeView{
			{Relation: "transition", Available: true, Target: "abc", Href: "?object=abc"},
			{Relation: "predecessor", Missing: "predecessor world relation is not stored by this projection"},
		}},
	}
	var body bytes.Buffer
	if err := Render(&body, p); err != nil {
		t.Fatal(err)
	}
	rendered := body.String()
	for _, want := range []string{`<a href="/workbench?object=`, "UNAVAILABLE:", "predecessor world relation is not stored by this projection"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered page missing %q: %q", want, rendered)
		}
	}
}

func renderPage(t *testing.T, p Page) string {
	t.Helper()
	var body bytes.Buffer
	if err := Render(&body, p); err != nil {
		t.Fatal(err)
	}
	return body.String()
}

func worldNav(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, `<nav aria-label="world browser">`)
	if start < 0 {
		t.Fatal("world nav missing")
	}
	end := strings.Index(body[start:], "</nav>")
	if end < 0 {
		t.Fatal("world nav unclosed")
	}
	return body[start : start+end]
}

func TestRenderWorldPane(t *testing.T) {
	root := EdgeView{Relation: "stateRoot", Available: true, Target: "root-205", Href: "?object=root-205"}
	world := WorldView{Ref: "world-205", Revision: 17, LogHead: "log-205", StateRoot: root, Available: true}
	nav := func(p WorldView) string {
		t.Helper()
		got := worldNav(t, renderPage(t, Page{Title: "world", World: p}))
		if !strings.Contains(got, `<a href="/workbench">workbench</a>`) {
			t.Fatalf("home link missing: %s", got)
		}
		return got
	}
	t.Run("metadata", func(t *testing.T) {
		got := nav(world)
		for _, want := range []string{`<dt>world</dt><dd><span class="hash" title="world-205" aria-label="world-205">world-205</span></dd>`, `<dt>revision</dt><dd>17</dd>`, `<dt>log head</dt><dd><span class="hash" title="log-205" aria-label="log-205">log-205</span></dd>`} {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q in %s", want, got)
			}
		}
		if strings.Contains(got, `aria-label="world-205">world-205</a>`) || strings.Contains(got, `aria-label="log-205">log-205</a>`) {
			t.Error("world or log head became an anchor")
		}
	})
	t.Run("stored-root", func(t *testing.T) {
		got := nav(world)
		if want := `stateRoot: <a href="/workbench?object=root-205" class="hash" title="root-205" aria-label="root-205">root-205</a>`; !strings.Contains(got, want) {
			t.Errorf("missing root edge %q in %s", want, got)
		}
		if strings.Count(got, `?object=`) != 1 {
			t.Errorf("object anchors: %s", got)
		}
	})
	t.Run("missing-root", func(t *testing.T) {
		p := world
		p.StateRoot = EdgeView{Relation: "stateRoot", Target: "root-205", Missing: "object root-205 is not stored"}
		got := nav(p)
		if !strings.Contains(got, `stateRoot: <span class="unavailable" role="note">UNAVAILABLE: object root-205 is not stored</span>`) || strings.Contains(got, `?object=`) {
			t.Errorf("missing root: %s", got)
		}
	})
	t.Run("no-selected-world", func(t *testing.T) {
		got := nav(WorldView{Unavailable: "no world selected"})
		if !strings.Contains(got, `UNAVAILABLE: no world selected`) || strings.Contains(got, `?object=`) {
			t.Errorf("no world: %s", got)
		}
	})
	t.Run("escaping", func(t *testing.T) {
		p := world
		p.Ref = `<script>alert(1)</script>`
		p.LogHead = `" onmouseover="x`
		p.StateRoot.Target = `<root&>`
		got := nav(p)
		if strings.Contains(got, `<script>`) || strings.Contains(got, `onmouseover="x`) || !strings.Contains(got, `&lt;script&gt;`) || !strings.Contains(got, `&lt;root&amp;&gt;`) {
			t.Errorf("unescaped nav: %s", got)
		}
	})
	t.Run("revision-zero", func(t *testing.T) {
		p := world
		p.Revision = 0
		if got := nav(p); !strings.Contains(got, `<dt>revision</dt><dd>0</dd>`) {
			t.Errorf("zero revision: %s", got)
		}
	})
}

// selectedArticle returns the substring from the selected-entry marker to the
// first </article> after it, or "" when the marker is absent.
func selectedArticle(body string) string {
	const marker = `<article aria-label="selected entry">`
	start := strings.Index(body, marker)
	if start < 0 {
		return ""
	}
	end := strings.Index(body[start:], "</article>")
	if end < 0 {
		return body[start:]
	}
	return body[start : start+end]
}

func TestRenderSelectedEntry(t *testing.T) {
	p := Page{Title: "selected", Selected: &EntryView{
		EntryIndex: 7,
		EntryHash:  "SEL-HASH-7",
		Edges: []EdgeView{
			{Relation: "transitionRef", Available: true, Target: "abc", Href: "?object=abc"},
			{Relation: "interpreter", Target: "def", Missing: "object def is not stored"},
		},
	}}
	article := selectedArticle(renderPage(t, p))
	if article == "" {
		t.Fatal(`rendered page has no <article aria-label="selected entry">`)
	}
	for _, want := range []string{
		`<h3>selected entry 7</h3>`,
		`title="SEL-HASH-7"`,
		`transitionRef: <a href="/workbench?object=abc"`,
		`interpreter: <span class="unavailable" role="note">UNAVAILABLE: object def is not stored</span>`,
	} {
		if !strings.Contains(article, want) {
			t.Errorf("selected-entry article missing %q: %q", want, article)
		}
	}
	// CONTROL: the article is Selected's and nothing else's.
	if body := renderPage(t, Page{Title: "unselected"}); strings.Contains(body, `aria-label="selected entry"`) {
		t.Errorf("Selected=nil still rendered a selected-entry article: %q", body)
	}
}

func TestRenderTimelineRowSelectLink(t *testing.T) {
	body := renderPage(t, Page{Title: "rows", Timeline: TimelineView{Entries: []EntryView{{EntryIndex: 3, SelectHref: "?from=0&entry=3"}}}})
	if want := `<a href="/workbench?from=0&amp;entry=3">select entry 3</a>`; !strings.Contains(body, want) {
		t.Errorf("row select link %q missing: %q", want, body)
	}
}

func TestRenderTimelinePagingLinks(t *testing.T) {
	t.Run("prev", func(t *testing.T) {
		body := renderPage(t, Page{Title: "prev", Timeline: TimelineView{PrevHref: "?from=0&entry=0"}})
		if want := `<a href="/workbench?from=0&amp;entry=0">previous</a>`; !strings.Contains(body, want) {
			t.Errorf("previous link %q missing: %q", want, body)
		}
	})
	t.Run("next", func(t *testing.T) {
		body := renderPage(t, Page{Title: "next", Timeline: TimelineView{NextHref: "?from=100&entry=100"}})
		if want := `<a href="/workbench?from=100&amp;entry=100">next</a>`; !strings.Contains(body, want) {
			t.Errorf("next link %q missing: %q", want, body)
		}
	})
	t.Run("neither", func(t *testing.T) {
		body := renderPage(t, Page{Title: "neither"})
		for _, unwanted := range []string{">previous</a>", ">next</a>"} {
			if strings.Contains(body, unwanted) {
				t.Errorf("zero TimelineView rendered %q", unwanted)
			}
		}
	})
}

var workbenchViewTypes = []reflect.Type{
	reflect.TypeOf(Page{}), reflect.TypeOf(TimelineView{}), reflect.TypeOf(EntryView{}),
	reflect.TypeOf(WorldView{}), reflect.TypeOf(ObjectView{}), reflect.TypeOf(GradeView{}),
	reflect.TypeOf(EdgeView{}), reflect.TypeOf(CommitView{}), reflect.TypeOf(ReferenceView{}),
}

type fieldCensus struct {
	templates *template.Template
	seen      map[string]bool
	visited   map[string]bool
	active    map[string]bool
}

func censusType(typ reflect.Type) reflect.Type {
	for typ != nil && typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func (c *fieldCensus) visitType(typ reflect.Type) {
	typ = censusType(typ)
	if typ != nil && typ.Kind() == reflect.Struct {
		c.visited[typ.Name()] = true
	}
}

func (c *fieldCensus) field(dot reflect.Type, names []string) (reflect.Type, error) {
	typ := dot
	for _, name := range names {
		typ = censusType(typ)
		if typ == nil || typ.Kind() != reflect.Struct {
			return nil, fmt.Errorf("field %s on non-struct %v", name, typ)
		}
		field, ok := typ.FieldByName(name)
		if !ok || field.PkgPath != "" {
			return nil, fmt.Errorf("unresolved %s.%s", typ.Name(), name)
		}
		c.visitType(typ)
		key := typ.Name() + "." + field.Name
		c.seen[key] = true
		typ = field.Type
		c.visitType(typ)
	}
	return typ, nil
}

func (c *fieldCensus) expr(node parse.Node, dot reflect.Type) (reflect.Type, error) {
	switch n := node.(type) {
	case *parse.DotNode:
		c.visitType(dot)
		return dot, nil
	case *parse.FieldNode:
		return c.field(dot, n.Ident)
	case *parse.StringNode:
		return reflect.TypeOf(""), nil
	case *parse.BoolNode:
		return reflect.TypeOf(true), nil
	case *parse.NumberNode:
		return reflect.TypeOf(int64(0)), nil
	case *parse.PipeNode:
		if len(n.Decl) != 0 || n.IsAssign || len(n.Cmds) == 0 {
			return nil, fmt.Errorf("unsupported pipeline declaration or empty pipeline: %s", n)
		}
		if len(n.Cmds) != 1 {
			return nil, fmt.Errorf("unsupported chained pipeline: %s", n)
		}
		var result reflect.Type
		for _, cmd := range n.Cmds {
			var err error
			result, err = c.expr(cmd, dot)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case *parse.CommandNode:
		if len(n.Args) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		if fun, ok := n.Args[0].(*parse.IdentifierNode); ok {
			args := make([]reflect.Type, 0, len(n.Args)-1)
			for _, arg := range n.Args[1:] {
				typ, err := c.expr(arg, dot)
				if err != nil {
					return nil, err
				}
				args = append(args, typ)
			}
			switch fun.Ident {
			case "edgeUnavailable":
				if len(args) != 1 || censusType(args[0]) != reflect.TypeOf(EdgeView{}) {
					return nil, fmt.Errorf("edgeUnavailable requires EdgeView, got %v", args)
				}
				c.seen["EdgeView.Available"] = true
				return reflect.TypeOf(true), nil
			case "workbenchHref":
				if len(args) != 1 || args[0] != reflect.TypeOf("") {
					return nil, fmt.Errorf("workbenchHref requires string, got %v", args)
				}
				return reflect.TypeOf(""), nil
			case "eq":
				if len(args) < 2 {
					return nil, fmt.Errorf("eq requires two arguments")
				}
				return reflect.TypeOf(true), nil
			default:
				return nil, fmt.Errorf("unknown template function %q", fun.Ident)
			}
		}
		if len(n.Args) != 1 {
			return nil, fmt.Errorf("unsupported command arguments: %s", n)
		}
		return c.expr(n.Args[0], dot)
	case *parse.VariableNode, *parse.ChainNode:
		return nil, fmt.Errorf("unsupported variable or chain: %s", node)
	default:
		return nil, fmt.Errorf("unsupported expression %T", node)
	}
}

func (c *fieldCensus) walk(node parse.Node, dot reflect.Type) error {
	if node == nil {
		return nil
	}
	c.visitType(dot)
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return nil
		}
		for _, child := range n.Nodes {
			if err := c.walk(child, dot); err != nil {
				return err
			}
		}
	case *parse.TextNode, *parse.CommentNode:
		return nil
	case *parse.ActionNode:
		_, err := c.expr(n.Pipe, dot)
		return err
	case *parse.IfNode:
		if _, err := c.expr(n.Pipe, dot); err != nil {
			return err
		}
		if err := c.walk(n.List, dot); err != nil {
			return err
		}
		return c.walk(n.ElseList, dot)
	case *parse.WithNode:
		next, err := c.expr(n.Pipe, dot)
		if err != nil {
			return err
		}
		if err := c.walk(n.List, next); err != nil {
			return err
		}
		return c.walk(n.ElseList, dot)
	case *parse.RangeNode:
		collection, err := c.expr(n.Pipe, dot)
		if err != nil {
			return err
		}
		collection = censusType(collection)
		if collection == nil {
			return fmt.Errorf("range has no type: %s", n.Pipe)
		}
		if collection.Kind() != reflect.Slice && collection.Kind() != reflect.Array && collection.Kind() != reflect.Map {
			return fmt.Errorf("unsupported range type %v", collection)
		}
		if err := c.walk(n.List, collection.Elem()); err != nil {
			return err
		}
		return c.walk(n.ElseList, dot)
	case *parse.TemplateNode:
		next, err := c.expr(n.Pipe, dot)
		if err != nil {
			return err
		}
		if next == nil {
			return fmt.Errorf("template %q has unresolved argument type", n.Name)
		}
		key := n.Name + "/" + next.String()
		if c.active[key] {
			return nil
		}
		child := c.templates.Lookup(n.Name)
		if child == nil || child.Tree == nil {
			return fmt.Errorf("unknown template %q", n.Name)
		}
		c.active[key] = true
		defer delete(c.active, key)
		return c.walk(child.Tree.Root, next)
	default:
		return fmt.Errorf("unsupported template node %T", node)
	}
	return nil
}

func missingWorkbenchFields(page, partials string) ([]string, error) {
	tm, err := template.New("workbench").Funcs(template.FuncMap{"edgeUnavailable": edgeUnavailable, "workbenchHref": workbenchHref}).Parse(page)
	if err != nil {
		return nil, err
	}
	tm, err = tm.Parse(partials)
	if err != nil {
		return nil, err
	}
	c := &fieldCensus{templates: tm, seen: map[string]bool{}, visited: map[string]bool{}, active: map[string]bool{}}
	if err := c.walk(tm.Tree.Root, reflect.TypeOf(Page{})); err != nil {
		return nil, err
	}
	var missing []string
	seen := c.seen
	if len(seen) == 0 {
		return nil, fmt.Errorf("template census observed no fields")
	}
	for _, typ := range workbenchViewTypes {
		if !c.visited[typ.Name()] {
			return nil, fmt.Errorf("view type %s was not visited", typ.Name())
		}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				continue
			}
			if !seen[typ.Name()+"."+field.Name] {
				missing = append(missing, typ.Name()+"."+field.Name)
			}
		}
	}
	sort.Strings(missing)
	return missing, nil
}

func TestWorkbenchViewFieldsAllRender(t *testing.T) {
	missing, err := missingWorkbenchFields(pageHTML, partialsHTML)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 0 {
		t.Fatalf("unrendered fields: %v", missing)
	}
}

func TestWorkbenchFieldCensusMutations(t *testing.T) {
	good, err := missingWorkbenchFields(pageHTML, partialsHTML)
	if err != nil || len(good) != 0 {
		t.Fatalf("known-good template: missing=%v err=%v", good, err)
	}
	for _, tc := range []struct{ name, old, replacement, want string }{
		{"world-ref-decoy", `{{.World.Ref}}`, `literal-world`, "WorldView.Ref"},
		{"world-unavailable", `{{.World.Unavailable}}`, `literal-unavailable`, "WorldView.Unavailable"},
		{"object-edges", `{{with .Object}}{{range .Edges}}`, `{{with .Object}}{{range .Commits.Edges}}`, "ObjectView.Edges"},
		{"commit-truncated", `{{if .Truncated}}<p>Showing 100 commits; more recorded</p>`, `{{if false}}<p>Showing 100 commits; more recorded</p>`, "CommitView.Truncated"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			page := pageHTML
			if tc.name == "world-ref-decoy" {
				if strings.Count(page, tc.old) == 0 {
					t.Fatal("world ref mutation anchor absent")
				}
				page = strings.ReplaceAll(page, tc.old, tc.replacement) + `<!-- .World.Ref decoy -->`
			} else {
				if strings.Count(page, tc.old) != 1 {
					t.Fatalf("mutation anchor count=%d", strings.Count(page, tc.old))
				}
				page = strings.Replace(page, tc.old, tc.replacement, 1)
			}
			missing, err := missingWorkbenchFields(page, partialsHTML)
			if err != nil {
				t.Fatal(err)
			}
			if len(missing) != 1 || missing[0] != tc.want {
				t.Fatalf("missing=%v, want [%s]", missing, tc.want)
			}
			good, err := missingWorkbenchFields(pageHTML, partialsHTML)
			if err != nil || len(good) != 0 {
				t.Fatalf("known-good control: %v %v", good, err)
			}
		})
	}
}

func TestEdgeUnavailableFunctionContract(t *testing.T) {
	if !edgeUnavailable(EdgeView{Available: false}) || edgeUnavailable(EdgeView{Available: true}) {
		t.Fatal("edgeUnavailable disagrees with Available")
	}
}

// provenanceWalkBody returns the trimmed text between the provenance-walk
// section's </h2> and its </section>, failing if either marker is missing.
func provenanceWalkBody(t *testing.T, body string) string {
	t.Helper()
	const start = `<h2>Provenance walk</h2>`
	i := strings.Index(body, start)
	if i < 0 {
		t.Fatalf("no provenance-walk heading in %q", body)
	}
	rest := body[i+len(start):]
	j := strings.Index(rest, "</section>")
	if j < 0 {
		t.Fatalf("provenance-walk section is not closed in %q", body)
	}
	return strings.TrimSpace(rest[:j])
}

func TestRenderProvenanceWalkNeverBlank(t *testing.T) {
	for _, tc := range []struct {
		name, want, unwanted string
		page                 Page
	}{
		{"no-object", `<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>`, "no provenance edges were supplied", Page{}},
		{"object-without-edges", `<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>`, "no object selected", Page{Object: &ObjectView{}}},
		{"supplied-edges", `<p>interface: <a href="/workbench?object=abc"`, "UNAVAILABLE", Page{Object: &ObjectView{Edges: []EdgeView{{Relation: "interface", Available: true, Target: "abc", Href: "?object=abc"}}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			walk := provenanceWalkBody(t, renderPage(t, tc.page))
			if walk == "" {
				t.Fatal("provenance walk rendered a heading followed by nothing")
			}
			if !strings.Contains(walk, tc.want) {
				t.Errorf("provenance walk missing %q: %q", tc.want, walk)
			}
			if strings.Contains(walk, tc.unwanted) {
				t.Errorf("provenance walk contains %q: %q", tc.unwanted, walk)
			}
		})
	}
}

func TestRenderGradeReasonNeverEmpty(t *testing.T) {
	for _, tc := range []struct {
		name, want, unwanted string
		grade                GradeView
	}{
		{"zero-grade", `<p>GRADE UNAVAILABLE — no grade reason was supplied</p>`, `GRADE UNAVAILABLE — </p>`, GradeView{}},
		{"supplied-reason", `<p>GRADE UNAVAILABLE — named reason</p>`, "no grade reason was supplied", NewGradeUnavailable("named reason")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := renderPage(t, Page{Title: "grade", Object: &ObjectView{Grade: tc.grade}})
			if !strings.Contains(body, tc.want) {
				t.Errorf("grade line missing %q: %q", tc.want, body)
			}
			if strings.Contains(body, tc.unwanted) {
				t.Errorf("grade line contains %q: %q", tc.unwanted, body)
			}
		})
	}
}

func TestReferenceView(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{}})
		if strings.Contains(body, `aria-label="referencedBy"`) {
			t.Fatal("nil references rendered")
		}
	})
	t.Run("empty", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{References: &ReferenceView{}}})
		for _, want := range []string{"none recorded in these entry/world fields", "Entries/worlds only: transitionRef, transitionFn, interpreter, stateRoot. Registry, journal and object-interface inbound references are not included."} {
			if !strings.Contains(body, want) {
				t.Fatalf("missing %q", want)
			}
		}
	})
	t.Run("checked", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{References: &ReferenceView{Edges: []EdgeView{{Relation: "transitionRef", Target: "entry 3", Href: "?from=3&entry=3", Available: true}}, Truncated: true}}})
		if !strings.Contains(body, "Showing 100 references; more recorded") || !strings.Contains(body, "transitionRef: <a") {
			t.Fatal(body)
		}
	})
}

func TestCommitView(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if body := renderPage(t, Page{Object: &ObjectView{}}); strings.Contains(body, `<h3>committedBy</h3>`) {
			t.Fatal("nil commits rendered")
		}
	})
	t.Run("none", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{Commits: &CommitView{}}})
		for _, want := range []string{`<h3>committedBy</h3>`, `Commits whose object set carried this object, oldest first. An object stored by PutObject or the journal before a commit carried it is attributed only to the commits that carried it.`, `no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.`} {
			if !strings.Contains(body, want) {
				t.Fatalf("missing %q: %s", want, body)
			}
		}
	})
	t.Run("continued-empty", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{Commits: &CommitView{Continued: true, FirstHref: "?object=a&refsAfter=b"}}})
		if !strings.Contains(body, "no further commits carried this object") || !strings.Contains(body, `first page</a>`) || strings.Contains(body, "no commit carried this object:") {
			t.Fatal(body)
		}
	})
	t.Run("available-and-truncated", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{Edges: []EdgeView{{Relation: "interface", Available: true, Target: "if", Href: "?object=if"}}, Commits: &CommitView{Edges: []EdgeView{{Relation: "committedBy", Available: true, Target: "entry 3", Href: "?from=3&entry=3"}}, Truncated: true, NextHref: "?object=a&commitsAfter=3"}, References: &ReferenceView{}}})
		for _, want := range []string{`committedBy: <a href="/workbench?from=3&amp;entry=3"`, `Showing 100 commits; more recorded`, `href="/workbench?object=a&amp;commitsAfter=3"`, `next commits</a>`} {
			if !strings.Contains(body, want) {
				t.Fatalf("missing %q: %s", want, body)
			}
		}
		if !(strings.Index(body, "interface: <a") < strings.Index(body, `<h3>committedBy</h3>`) && strings.Index(body, `<h3>committedBy</h3>`) < strings.Index(body, `<h3>referencedBy</h3>`)) {
			t.Fatal(body)
		}
	})
	t.Run("html-safe-link", func(t *testing.T) {
		body := renderPage(t, Page{Object: &ObjectView{Commits: &CommitView{Edges: []EdgeView{{Relation: "committedBy", Available: true, Target: "entry 1", Href: `?from=1&entry=1&x=<script>`}}}}})
		if strings.Contains(body, `<script>`) || !strings.Contains(body, `entry 1</a>`) {
			t.Fatal(body)
		}
	})
}
