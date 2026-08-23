package workbench

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"
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
	})
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
