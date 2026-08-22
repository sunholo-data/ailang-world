package workbench

import (
	"errors"
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
