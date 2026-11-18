package verbs

// Row 243 V4 S1 (design §4.1, §5.1): the bar's admitted operator subset — the episode start and
// episode finish lines the app may later show a review for (S7). It is pure and, in S1, has no caller:
// it reads Lookup's rows and runs ParseInterleaved over a flag set that mirrors the CLI's own
// (cmd/ailang-worldd TestOperatorFlagParity holds the mirror to the CLI's real flag sets). The ledger
// answers it depends on arrive as an input (OperatorPolicy); nothing unanswered is admitted.

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"regexp"
)

// OperatorFlag is one flag of the CLI's episode start or episode finish flag set: its name, whether
// it is a boolean (it takes no value word), and whether the app's subset admits it.
type OperatorFlag struct {
	Name  string
	Bool  bool
	Admit bool
}

// EpisodeStartFlags and EpisodeFinishFlags are the CLI's two flag sets, every flag in the CLI's order
// (runEpisodeStart, runEpisodeFinishWith); Admit marks design §4.1's subset.
var (
	EpisodeStartFlags = []OperatorFlag{
		{"db", false, false}, {"credentials-db", false, false}, {"workspace-root", false, false}, {"repo", false, false},
		{"branch", false, false}, {"profile", false, false}, {"exec-profiles-dir", false, false}, {"task", false, false},
		{"agent", false, false}, {"model", false, false}, {"ttl", false, false}, {"preset", false, true},
		{"presets-dir", false, false}, {"resume", false, false}, {"hosts-dir", false, false}, {"launch", true, true},
		{"headless", true, true}, {"request", true, false}, {"grant", false, false}, {"budget", false, false},
		{"sessions-dir", false, false}, {"addr", false, false},
	}
	EpisodeFinishFlags = []OperatorFlag{
		{"preset", false, true}, {"presets-dir", false, false}, {"title", false, false}, {"include-unattributed", true, false},
		{"ready", true, true}, {"exec-profiles-dir", false, false}, {"dry-run", true, false}, {"confirm-on-surface", true, true},
		{"allow-failing", false, false},
	}
)

// EpisodeIDPattern is the CLI's and the daemon's episode grammar (cmd workspaceEpisodePattern, host/daemon
// episodeIDPattern): 1 to 64 characters. MaxPresetName is the CLI's preset-name cap (cmd maxPresetNameLen).
// OperatorActionIDPattern is design §3.1's NEW action id, oa- and 32 lowercase hex digits; S1 only defines it.
var (
	EpisodeIDPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	OperatorActionIDPattern = regexp.MustCompile(`^oa-[0-9a-f]{32}$`)
)

const MaxPresetName = 51

// ValidOperatorActionID reports whether id is an operator action id.
func ValidOperatorActionID(id string) bool { return OperatorActionIDPattern.MatchString(id) }

// Answer is the ledger's answer to one decision; only "A", "B" and "C" are answers.
type Answer string

// OperatorPolicy is the caller's answers, a pure input: Start is D-WORLD-119's, Finish D-WORLD-120's,
// HandoffProven whether design §11 E4 has passed. The zero value — nothing answered — admits nothing.
type OperatorPolicy struct {
	Start, Finish Answer
	HandoffProven bool
}

func (p OperatorPolicy) headless() bool    { return p.Start == "A" || p.Start == "B" }
func (p OperatorPolicy) interactive() bool { return p.Start == "A" && p.HandoffProven }
func (p OperatorPolicy) finish() bool      { return p.headless() && (p.Finish == "A" || p.Finish == "B") }
func (p OperatorPolicy) ready() bool       { return p.Finish == "A" } // asked only after finish()

// OperatorClass is AdmitOperator's classification, the typed outcome a later GET /surface/verb.json
// renders (S7): review a supported tuple; a usage refusal; a terminal line (outside the subset, or a mode
// the answers do not allow); a never family; or another line whose Lookup classification stands.
type OperatorClass string

const (
	OperatorReview   OperatorClass = "review"
	OperatorUsage    OperatorClass = "usage"
	OperatorTerminal OperatorClass = "terminal"
	OperatorNever    OperatorClass = "never"
	OperatorOther    OperatorClass = "other"
)

// OperatorIntent is one admitted tuple; Argv is its canonical CLI line's words.
type OperatorIntent struct {
	Kind                    string // "start" or "finish"
	Episode, Preset         string
	Headless                bool // start: --launch --headless, else --launch (interactive)
	ConfirmOnSurface, Ready bool // finish
	Argv                    []string
}

// OperatorAnswer is AdmitOperator's answer. Intent is set exactly when Class is OperatorReview; Field
// names the word a refusal is about ("" for none), Reason says why, CLI is the line as typed.
type OperatorAnswer struct {
	Class              OperatorClass
	Intent             *OperatorIntent
	Field, Reason, CLI string
}

// boolWords are the flag package's boolean spellings (strconv.ParseBool's), each with its value; host/verbs
// imports no strconv (host/daemon TestSurfaceVerbIsReadOnly pins its imports).
var boolWords = map[string]bool{"1": true, "t": true, "T": true, "TRUE": true, "true": true, "True": true,
	"0": false, "f": false, "F": false, "FALSE": false, "false": false, "False": false}

// operatorValue counts a flag's occurrences and keeps its last value (a bool's as on).
type operatorValue struct {
	isBool, on bool
	n          int
	s          string
}

func (v *operatorValue) String() string   { return v.s }
func (v *operatorValue) IsBoolFlag() bool { return v.isBool }
func (v *operatorValue) Set(s string) error {
	if v.isBool {
		b, ok := boolWords[s]
		if !ok {
			return errors.New("parse error") // the flag package's own bool error, word for word
		}
		v.on = b
	}
	v.n++
	v.s = s
	return nil
}

// AdmitOperator classifies one bar line's words (Split's) for the operator subset under policy p
// (design §4.1). It never runs, reads or writes anything.
func AdmitOperator(argv []string, p OperatorPolicy) OperatorAnswer {
	if len(argv) > 0 && argv[0] == "world-publish" { // a separate binary, outside Lookup's word sets
		return OperatorAnswer{Class: OperatorNever, Field: "world-publish", CLI: Quote(argv),
			Reason: "world-publish is never run from the app: publishing is typed at a terminal; at a terminal: " + Quote(argv)}
	}
	v, rest := Lookup(argv)
	var mirror []OperatorFlag
	switch v.Words {
	case "episode start":
		mirror = EpisodeStartFlags
	case "episode finish":
		mirror = EpisodeFinishFlags
	default:
		if v.Class == Never {
			return OperatorAnswer{Class: OperatorNever, Reason: v.Reason, CLI: v.CLI}
		}
		return OperatorAnswer{Class: OperatorOther, Reason: v.Reason, CLI: v.CLI}
	}
	refuse := func(c OperatorClass, field, why string) OperatorAnswer {
		return OperatorAnswer{Class: c, Field: field, Reason: Binary + " " + v.Words + ": " + why, CLI: v.CLI}
	}
	at := "; at a terminal: " + v.CLI
	fs := flag.NewFlagSet(Binary+" "+v.Words, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	got := map[string]*operatorValue{}
	for _, f := range mirror {
		got[f.Name] = &operatorValue{isBool: f.Bool}
		fs.Var(got[f.Name], f.Name, "")
	}
	pos, err := ParseInterleaved(fs, rest)
	if errors.Is(err, flag.ErrHelp) {
		return refuse(OperatorUsage, "--help", "see "+Binary+" "+v.Words+" --help at a terminal")
	}
	if err != nil {
		return refuse(OperatorUsage, "", err.Error())
	}
	for _, f := range mirror {
		if got[f.Name].n > 1 {
			return refuse(OperatorUsage, "--"+f.Name, "--"+f.Name+" is given more than once")
		}
	}
	for _, f := range mirror {
		if g := got[f.Name]; g.n == 1 && !f.Admit {
			return refuse(OperatorTerminal, "--"+f.Name, "--"+f.Name+" is not run from the app"+at)
		}
		if g := got[f.Name]; g.n == 1 && f.Bool && !g.on {
			return refuse(OperatorTerminal, "--"+f.Name, "--"+f.Name+" is given false; the app runs only its true form"+at)
		}
	}
	if len(pos) != 1 {
		return refuse(OperatorUsage, "<episode>", "needs exactly one <episode>")
	}
	if !EpisodeIDPattern.MatchString(pos[0]) {
		return refuse(OperatorUsage, "<episode>", fmt.Sprintf("episode %q does not match %s", pos[0], EpisodeIDPattern))
	}
	preset := got["preset"].s
	if preset == "" {
		return refuse(OperatorUsage, "--preset", "--preset <name> is required")
	}
	if preset == "se-tools" || len(preset) > MaxPresetName || !EpisodeIDPattern.MatchString(preset) {
		return refuse(OperatorUsage, "--preset", fmt.Sprintf("preset %q is not a launch preset name (the episode grammar, at most %d characters)", preset, MaxPresetName))
	}
	in := &OperatorIntent{Episode: pos[0], Preset: preset}
	if v.Words == "episode start" {
		launch, headless := got["launch"].on, got["headless"].on
		switch {
		case headless && !launch:
			return refuse(OperatorUsage, "--headless", "--headless requires --launch")
		case !launch:
			return refuse(OperatorTerminal, "--launch", "the app starts an episode only with --launch"+at)
		case headless && !p.headless():
			return refuse(OperatorTerminal, "--headless", "starting from the app waits for D-WORLD-119 (A or B)"+at)
		case !headless && !p.interactive():
			return refuse(OperatorTerminal, "--launch", "an interactive launch from the app needs D-WORLD-119 A and the herdr hand-off proven (design §11 E4)"+at)
		}
		in.Kind, in.Headless = "start", headless
		in.Argv = []string{Binary, "episode", "start", in.Episode, "--preset", preset, "--launch"}
		if headless {
			in.Argv = append(in.Argv, "--headless")
		}
		return OperatorAnswer{Class: OperatorReview, Intent: in, CLI: v.CLI}
	}
	in.Kind, in.ConfirmOnSurface, in.Ready = "finish", got["confirm-on-surface"].on, got["ready"].on
	if !p.finish() {
		return refuse(OperatorTerminal, "", "finishing from the app waits for D-WORLD-119 and D-WORLD-120 (A or B)"+at)
	}
	if in.Ready && !p.ready() {
		return refuse(OperatorTerminal, "--ready", "--ready from the app needs D-WORLD-120 A"+at)
	}
	in.Argv = []string{Binary, "episode", "finish", in.Episode, "--preset", preset}
	if in.ConfirmOnSurface {
		in.Argv = append(in.Argv, "--confirm-on-surface")
	}
	if in.Ready {
		in.Argv = append(in.Argv, "--ready")
	}
	return OperatorAnswer{Class: OperatorReview, Intent: in, CLI: v.CLI}
}
