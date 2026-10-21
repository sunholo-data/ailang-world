package surface

// Row 188 M0c: the episode facts. A step state gains the rig's episode that
// built it (steps[]), the files it changed (files{}), the files it read
// mapped to where they last changed (reads{}), and the decisions taken on it
// (decisions[]).
//
// The inputs are the rig's episode records (trial-episode/1, one per round),
// its prs.json, and the daemon's episode report; the caller reads them and
// runs git through the Git seam. Nothing here touches the filesystem or the
// network, and no path field of a record is ever copied into the index.

import (
	"encoding/json"
	"fmt"
	"math"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RecordSchema is the rig's episode record type.
const RecordSchema = "trial-episode/1"

// ReadsMethod names how reads{} maps a path to a source.
const ReadsMethod = "last-change-at-base/v1"

// MaxChanges bounds one step's listed changes; the rest count in truncated.
const MaxChanges = 1000

// Read vias: where a read path last changed before the step's base.
const (
	ViaNew     = "new"     // not at base: the step made it
	ViaStep    = "step"    // a merge of a step state in the index
	ViaCommit  = "commit"  // a non-merge commit
	ViaOutside = "outside" // a merge older than the index window
)

// Decision actors and kinds.
const (
	ActorWorld = "world"
	ActorYou   = "you"
	ActorLoop  = "loop"

	DecisionGrade   = "grade"
	DecisionRefusal = "refusal"
	DecisionMerge   = "merge"
)

// Delivered is the last round with successful operator re-runs applied.
type Delivered struct {
	Passed   int      `json:"passed"`
	Total    int      `json:"total"`
	Regraded []string `json:"regraded"`
}

// Step is one episode that built a step state.
//
// State is the step state's id. Task, Agent, Model and Base (solution.base)
// are round 1's, Graders the last round's, and Minutes the agent wall time
// of every round. Calls (total_calls) and Entries ([first, last] log
// entries) are the report's, absent when the daemon has no report.
type Step struct {
	Summary     string     `json:"summary,omitempty"`
	SummaryFull string     `json:"summary_full,omitempty"`
	Failed      []string   `json:"failed,omitempty"`
	Delivered   *Delivered `json:"delivered,omitempty"`
	Episode     string     `json:"episode"`
	State       string     `json:"state"`
	PR          int        `json:"pr"`
	Task        string     `json:"task"`
	Agent       string     `json:"agent"`
	Model       string     `json:"model"`
	Rounds      int        `json:"rounds"`
	Graders     Graders    `json:"graders"`
	Minutes     float64    `json:"minutes"`
	Calls       *int       `json:"calls,omitempty"`
	Entries     *[2]int64  `json:"entries,omitempty"`
	Base        string     `json:"base"`
}

// Graders counts a round's graders.
type Graders struct {
	Passed int `json:"passed"`
	Total  int `json:"total"`
}

// FileChanges is `git diff --name-status --no-renames <base> <step>`.
type FileChanges struct {
	Changes   []Change `json:"changes"`
	Truncated int      `json:"truncated"`
}

// Change is one changed path; Status is A, M or D.
type Change struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

// Read is one path a step's episode read, and where it last changed.
type Read struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
	Via   string `json:"via"`
	From  string `json:"from"` // a step id (via step), a commit (via commit), else ""
}

// Decision is one act on a step.
type Decision struct {
	Step  string `json:"step"`
	Actor string `json:"actor"`
	Kind  string `json:"kind"`
	Text  string `json:"text"`
	At    string `json:"at"` // UTC RFC 3339
}

// RigSource says what the rig gave: its episodes of the project, how many
// were attached to a step, and the ids of the rest.
type RigSource struct {
	Episodes   int      `json:"episodes"`
	Attached   int      `json:"attached"`
	Unattached []string `json:"unattached"`
}

// EpisodeRecord is the part of one trial-episode/1 record the index uses.
// Its path fields (worktree, transcript, report, graders' outputs) are not
// decoded at all.
type EpisodeRecord struct {
	AgentSummary string          `json:"agent_summary"`
	Regraded     map[string]bool `json:"-"`
	Schema       string          `json:"schema"`
	Episode      string          `json:"episode"`
	Project      string          `json:"project"`
	Agent        string          `json:"agent"`
	Model        string          `json:"model"`
	Task         string          `json:"task"`
	StartedAt    string          `json:"started_at"`
	FinishedAt   string          `json:"finished_at"`
	AgentRun     recordRun       `json:"agent_run"`
	Solution     recordSolution  `json:"solution"`
	Graders      []recordGrader  `json:"graders"`

	Started, Finished time.Time `json:"-"`
}

type recordRun struct {
	WallMS int64 `json:"wall_ms"`
}

type recordSolution struct {
	Base string `json:"base"`
}

type recordGrader struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

// ParseEpisodeRecord decodes one record and checks the fields the index uses.
func ParseEpisodeRecord(data []byte) (EpisodeRecord, error) {
	var r EpisodeRecord
	if err := json.Unmarshal(data, &r); err != nil {
		return EpisodeRecord{}, fmt.Errorf("surface: episode record does not decode: %v", err)
	}
	if r.Schema != RecordSchema {
		return EpisodeRecord{}, fmt.Errorf("surface: episode record schema %q, want %q", r.Schema, RecordSchema)
	}
	if r.Episode == "" || strings.ContainsAny(r.Episode, `/\`) {
		return EpisodeRecord{}, fmt.Errorf("surface: episode record id %q is not an episode id", r.Episode)
	}
	var err error
	if r.Started, err = time.Parse(time.RFC3339, r.StartedAt); err != nil {
		return EpisodeRecord{}, fmt.Errorf("surface: episode %s: started_at %q is not RFC 3339", r.Episode, r.StartedAt)
	}
	if r.Finished, err = time.Parse(time.RFC3339, r.FinishedAt); err != nil {
		return EpisodeRecord{}, fmt.Errorf("surface: episode %s: finished_at %q is not RFC 3339", r.Episode, r.FinishedAt)
	}
	if r.AgentRun.WallMS < 0 {
		return EpisodeRecord{}, fmt.Errorf("surface: episode %s: agent_run.wall_ms %d is negative", r.Episode, r.AgentRun.WallMS)
	}
	return r, nil
}

// Episode is one episode's rounds, ordered by started_at.
type Episode struct {
	ID     string
	Rounds []EpisodeRecord
}

// GroupRounds groups records by episode: each record is one round, and the
// rounds are ordered by started_at (then finished_at). Episodes are ordered
// by id.
func GroupRounds(records []EpisodeRecord) []Episode {
	by := map[string][]EpisodeRecord{}
	for _, r := range records {
		by[r.Episode] = append(by[r.Episode], r)
	}
	out := make([]Episode, 0, len(by))
	for id, rounds := range by {
		sort.SliceStable(rounds, func(i, j int) bool {
			if !rounds[i].Started.Equal(rounds[j].Started) {
				return rounds[i].Started.Before(rounds[j].Started)
			}
			return rounds[i].Finished.Before(rounds[j].Finished)
		})
		out = append(out, Episode{ID: id, Rounds: rounds})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// RoundDir reports whether an evidence directory holds project's records:
// <project> for round 1, <project>-r2, <project>-r3, … for later rounds.
func RoundDir(project, name string) bool {
	if name == project {
		return true
	}
	n, ok := strings.CutPrefix(name, project+"-r")
	if !ok || n == "" || n[0] == '0' {
		return false
	}
	v, err := strconv.Atoi(n)
	return err == nil && v >= 2
}

// ParsePRs reads prs.json (episode → "<project>#<N>", or prose) and returns
// the PR number of each episode whose value names project's PR exactly.
func ParsePRs(data []byte, project string) (map[string]int, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("surface: prs.json does not decode: %v", err)
	}
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(project) + `#([1-9][0-9]{0,8})$`)
	out := map[string]int{}
	for ep, v := range raw {
		var s string
		if json.Unmarshal(v, &s) != nil {
			continue
		}
		if m := re.FindStringSubmatch(s); m != nil {
			out[ep], _ = strconv.Atoi(m[1])
		}
	}
	return out, nil
}

// Report is the part of the daemon's EpisodeReport (GET
// /v1/episodes/<ep>/report) the index uses.
type Report struct {
	Entries         [2]int64    `json:"entries"`
	TotalCalls      int         `json:"total_calls"`
	PlanRefusals    reportCount `json:"plan_refusals"`
	HandlerRefusals reportCount `json:"handler_refusals"`
	Reads           reportReads `json:"reads"`
	LastAt          string      `json:"last_at"`
}

type reportCount struct {
	Count int `json:"count"`
}

type reportReads struct {
	Files []reportRead `json:"files"`
}

type reportRead struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// ParseReport decodes a report body.
func ParseReport(data []byte) (Report, error) {
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return Report{}, fmt.Errorf("surface: episode report does not decode: %v", err)
	}
	if r.PlanRefusals.Count+r.HandlerRefusals.Count > 0 {
		if _, err := time.Parse(time.RFC3339, r.LastAt); err != nil {
			return Report{}, fmt.Errorf("surface: episode report last_at %q is not RFC 3339", r.LastAt)
		}
	}
	return r, nil
}

// Git is the read-only git the episode facts need; the caller runs it in
// the repository (ailang-worldd surface-index), a test fakes it.
type Git interface {
	// NameStatus is `git diff --name-status --no-renames -z <base> <step>`.
	NameStatus(base, step string) (string, error)
	// Tree is `git ls-tree -r -z --name-only <rev>`.
	Tree(rev string) (string, error)
	// LastChange is `git log --first-parent -1 --format=%H%x20%P <rev> -- <path>`:
	// the last commit on rev's first-parent line that changed path, and its
	// parents ("" when none did).
	LastChange(rev, path string) (string, error)
}

// ParseNameStatus parses -z name-status output into changes sorted by path.
// A type change (T) is a modification; any other status is an error.
func ParseNameStatus(z string) ([]Change, error) {
	tok := strings.Split(strings.TrimSuffix(z, "\x00"), "\x00")
	if len(tok) == 1 && tok[0] == "" {
		return []Change{}, nil
	}
	if len(tok)%2 != 0 {
		return nil, fmt.Errorf("surface: name-status has %d fields, want status and path pairs", len(tok))
	}
	out := make([]Change, 0, len(tok)/2)
	for i := 0; i < len(tok); i += 2 {
		st, p := strings.TrimSpace(tok[i]), tok[i+1]
		switch st {
		case "A", "M", "D":
		case "T":
			st = "M"
		default:
			return nil, fmt.Errorf("surface: name-status %q for %q is not A, M, D or T", st, p)
		}
		if !filepath.IsLocal(p) {
			return nil, fmt.Errorf("surface: name-status path %q is not local", p)
		}
		out = append(out, Change{Path: p, Status: st})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// ParseTree parses -z ls-tree --name-only output into a set of paths.
func ParseTree(z string) map[string]bool {
	out := map[string]bool{}
	for _, p := range strings.Split(z, "\x00") {
		if p != "" {
			out[p] = true
		}
	}
	return out
}

// ClassifyLastChange maps LastChange's line to a read's via and from: a
// merge whose second parent is one of steps is via step; another merge is
// via outside; a non-merge commit is via commit; no line is via outside.
func ClassifyLastChange(line string, steps map[string]bool) (via, from string, err error) {
	f := strings.Fields(line)
	if len(f) == 0 {
		return ViaOutside, "", nil
	}
	for _, h := range f {
		if !commitHash.MatchString(h) {
			return "", "", fmt.Errorf("surface: last change %q is not commit hashes", line)
		}
	}
	switch {
	case len(f) <= 2:
		return ViaCommit, f[0], nil
	case steps[f[2]]:
		return ViaStep, f[2], nil
	default:
		return ViaOutside, "", nil
	}
}

// RigInput is what AttachRig needs from the rig and the daemon: the
// project's episodes (GroupRounds), its PRs (ParsePRs), and each episode's
// report, absent when the daemon has none.
type RigInput struct {
	Episodes []Episode
	PRs      map[string]int
	Reports  map[string]Report
}

// AttachRig adds the episode facts to an index Build made. An episode is
// attached when its PR is a step state's PR in the index; when two claim one
// step, the one whose first round started first is attached. Every other
// episode is listed in source.rig.unattached. The output depends only on the
// input and git's answers, never on map order.
func AttachRig(ix Index, in RigInput, git Git) (Index, error) {
	byID := map[string]int{}
	stepByPR := map[int]int{}
	stepIDs := map[string]bool{}
	for i, s := range ix.States {
		byID[s.ID] = i
		if s.Kind == KindStep {
			stepIDs[s.ID] = true
			if s.PR > 0 {
				stepByPR[s.PR] = i
			}
		}
	}
	braneKind := map[string]string{}
	for _, b := range ix.Branes {
		braneKind[b.Key] = b.Kind
	}

	eps := append([]Episode(nil), in.Episodes...)
	sort.SliceStable(eps, func(i, j int) bool {
		a, b := eps[i].Rounds[0].Started, eps[j].Rounds[0].Started
		if !a.Equal(b) {
			return a.Before(b)
		}
		return eps[i].ID < eps[j].ID
	})
	type attached struct {
		ep Episode
		st State
		at time.Time
	}
	var atts []attached
	claimed := map[string]bool{}
	unattached := []string{}
	for _, ep := range eps {
		if len(ep.Rounds) == 0 {
			return Index{}, fmt.Errorf("surface: episode %s has no rounds", ep.ID)
		}
		pr, ok := in.PRs[ep.ID]
		i, inIndex := stepByPR[pr]
		if !ok || !inIndex || claimed[ix.States[i].ID] {
			unattached = append(unattached, ep.ID)
			continue
		}
		st := ix.States[i]
		claimed[st.ID] = true
		at, err := time.Parse(time.RFC3339, st.At)
		if err != nil {
			return Index{}, fmt.Errorf("surface: step %s: at %q is not RFC 3339", st.ID, st.At)
		}
		atts = append(atts, attached{ep, st, at})
	}
	sort.SliceStable(atts, func(i, j int) bool {
		if !atts[i].at.Equal(atts[j].at) {
			return atts[i].at.Before(atts[j].at)
		}
		return atts[i].ep.ID < atts[j].ep.ID
	})
	sort.Strings(unattached)

	out := ix
	out.Source.Rig = &RigSource{Episodes: len(in.Episodes), Attached: len(atts), Unattached: unattached}
	out.Source.ReadsMethod = ReadsMethod
	out.Steps, out.Files, out.Reads, out.Decisions = nil, nil, nil, nil
	if len(atts) > 0 {
		out.Steps = []Step{}
		out.Files = map[string]FileChanges{}
		out.Reads = map[string][]Read{}
		out.Decisions = []Decision{}
	}
	for _, a := range atts {
		rep, hasReport := in.Reports[a.ep.ID]
		out.Steps = append(out.Steps, stepOf(a.ep, a.st, rep, hasReport))

		fc, err := filesOf(git, a.st)
		if err != nil {
			return Index{}, err
		}
		out.Files[a.st.ID] = fc
		if hasReport {
			reads, err := readsOf(git, a.st, rep, stepIDs)
			if err != nil {
				return Index{}, err
			}
			out.Reads[a.st.ID] = reads
		}
		merge := ix.States[byID[a.st.Merge]]
		ds, err := decisionsOf(a.ep, a.st, merge, braneKind[a.st.Brane], rep, hasReport)
		if err != nil {
			return Index{}, err
		}
		out.Decisions = append(out.Decisions, ds...)
	}
	sort.SliceStable(out.Decisions, func(i, j int) bool {
		x, y := out.Decisions[i], out.Decisions[j]
		if x.At != y.At { // all UTC RFC 3339 with no fraction, so text order is time order
			return x.At < y.At
		}
		if x.Kind != y.Kind {
			return x.Kind < y.Kind
		}
		return x.Step < y.Step
	})
	return out, nil
}

func stepOf(ep Episode, st State, rep Report, hasReport bool) Step {
	first, last := ep.Rounds[0], ep.Rounds[len(ep.Rounds)-1]
	var wall int64
	for _, r := range ep.Rounds {
		wall += r.AgentRun.WallMS
	}
	s := Step{Episode: ep.ID, State: st.ID, PR: st.PR, Task: taskName(first.Task), Agent: first.Agent,
		Model: first.Model, Rounds: len(ep.Rounds), Graders: gradersOf(last),
		Minutes: math.Round(float64(wall)/6000) / 10}
	s.Summary, s.SummaryFull = summaries(first.AgentSummary)
	s.Delivered = &Delivered{Passed: s.Graders.Passed, Total: s.Graders.Total, Regraded: []string{}}
	for _, g := range last.Graders {
		if !g.Passed && g.Name != "" {
			s.Failed = append(s.Failed, g.Name)
			if last.Regraded[g.Name] {
				s.Delivered.Passed++
				s.Delivered.Regraded = append(s.Delivered.Regraded, g.Name)
			}
		}
	}
	sort.Strings(s.Failed)
	sort.Strings(s.Delivered.Regraded)
	if commitHash.MatchString(first.Solution.Base) {
		s.Base = first.Solution.Base
	}
	if hasReport {
		calls, entries := rep.TotalCalls, rep.Entries
		s.Calls, s.Entries = &calls, &entries
	}
	return s
}

// taskName is a task's name: a task given as a file keeps its base name only,
// so no directory reaches the index.
func taskName(task string) string {
	if strings.ContainsAny(task, `/\`) {
		return path.Base(filepath.ToSlash(task))
	}
	return task
}

func gradersOf(r EpisodeRecord) Graders {
	g := Graders{Total: len(r.Graders)}
	for _, x := range r.Graders {
		if x.Passed {
			g.Passed++
		}
	}
	return g
}

func filesOf(git Git, st State) (FileChanges, error) {
	z, err := git.NameStatus(st.Base, st.ID)
	if err != nil {
		return FileChanges{}, fmt.Errorf("surface: step %s: name-status: %v", st.ID, err)
	}
	changes, err := ParseNameStatus(z)
	if err != nil {
		return FileChanges{}, fmt.Errorf("step %s: %v", st.ID, err)
	}
	fc := FileChanges{Changes: changes}
	if len(changes) > MaxChanges {
		fc.Truncated = len(changes) - MaxChanges
		fc.Changes = changes[:MaxChanges]
	}
	return fc, nil
}

// readsOf classifies each path the report lists, in its order (first read
// first). A path is cleaned before it is looked up; one that is not local is
// skipped, never copied, and two that clean to one path are counted as one.
func readsOf(git Git, st State, rep Report, stepIDs map[string]bool) ([]Read, error) {
	out := []Read{}
	if len(rep.Reads.Files) == 0 {
		return out, nil
	}
	treeZ, err := git.Tree(st.Base)
	if err != nil {
		return nil, fmt.Errorf("surface: step %s: ls-tree %s: %v", st.ID, st.Base, err)
	}
	tree := ParseTree(treeZ)
	seen := map[string]int{}
	for _, f := range rep.Reads.Files {
		p := path.Clean(filepath.ToSlash(f.Path))
		if !filepath.IsLocal(f.Path) || !filepath.IsLocal(p) {
			continue
		}
		if i, dup := seen[p]; dup {
			out[i].Count += f.Count
			continue
		}
		r := Read{Path: p, Count: f.Count, Via: ViaNew}
		if tree[p] {
			line, err := git.LastChange(st.Base, p)
			if err != nil {
				return nil, fmt.Errorf("surface: step %s: last change of %s: %v", st.ID, p, err)
			}
			if r.Via, r.From, err = ClassifyLastChange(line, stepIDs); err != nil {
				return nil, fmt.Errorf("step %s: %s: %v", st.ID, p, err)
			}
		}
		seen[p] = len(out)
		out = append(out, r)
	}
	return out, nil
}

func decisionsOf(ep Episode, st, merge State, braneKind string, rep Report, hasReport bool) ([]Decision, error) {
	var out []Decision
	for n, r := range ep.Rounds {
		g := gradersOf(r)
		out = append(out, Decision{Step: st.ID, Actor: ActorWorld, Kind: DecisionGrade,
			Text: fmt.Sprintf("round %d: %d/%d graders", n+1, g.Passed, g.Total), At: utc(r.Finished)})
	}
	if hasReport && rep.PlanRefusals.Count+rep.HandlerRefusals.Count > 0 {
		t, err := time.Parse(time.RFC3339, rep.LastAt)
		if err != nil {
			return nil, fmt.Errorf("surface: episode %s: report last_at %q is not RFC 3339", ep.ID, rep.LastAt)
		}
		out = append(out, Decision{Step: st.ID, Actor: ActorWorld, Kind: DecisionRefusal,
			Text: fmt.Sprintf("%d plan + %d handler refusals", rep.PlanRefusals.Count, rep.HandlerRefusals.Count), At: utc(t)})
	}
	t, err := time.Parse(time.RFC3339, merge.At)
	if err != nil {
		return nil, fmt.Errorf("surface: merge %s: at %q is not RFC 3339", merge.ID, merge.At)
	}
	actor := ActorYou
	if braneKind == BraneLoop {
		actor = ActorLoop
	}
	out = append(out, Decision{Step: st.ID, Actor: actor, Kind: DecisionMerge,
		Text: fmt.Sprintf("merged #%d", st.PR), At: utc(t)})
	return out, nil
}

func utc(t time.Time) string { return t.UTC().Format(time.RFC3339) }
