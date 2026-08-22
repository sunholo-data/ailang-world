package workbench

import "errors"

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
