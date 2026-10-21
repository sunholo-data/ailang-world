// Package surface is the surface index (row 188, w-world-os-surface S1 M0a):
// the document `GET /surface` reads, built outside the daemon from git
// history, and with a rig the episode facts (M0c, episodes.go).
//
//   - The shared world is the first-parent merge commits on a branch, one
//     shared state each.
//   - Each merge's second parent is a step: the PR's side.
//   - A step's brane comes from its PR branch's name (BraneOf), and the
//     brane's title from the charter's queue row (ChartRowTitles).
//
// Everything here is pure: the caller runs git (ailang-worldd surface-index)
// and hands Build the parsed log and each merge's merge-base. The package
// uses the standard library only, so the daemon can import it (M1).
package surface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Schema is the index's document type.
const Schema = "world/surface-index/v1"

// MaxBytes bounds index reads by the builder and daemon (4 MiB).
const MaxBytes = 4 << 20

// MaxStates bounds an index: Validate refuses more.
const MaxStates = 5000

// State kinds.
const (
	KindShared = "shared"
	KindStep   = "step"
)

// Brane kinds: who took the brane's steps.
const (
	BraneAttended = "attended"
	BraneLoop     = "loop"
	BraneMixed    = "mixed"
)

// Index is the world/surface-index/v1 document.
type Index struct {
	Schema  string  `json:"schema"`
	BuiltAt string  `json:"built_at"`
	Source  Source  `json:"source"`
	Branes  []Brane `json:"branes"`
	States  []State `json:"states"`
	// The episode facts (M0c, AttachRig): present only when the index was
	// built with a rig, so an index from git alone is unchanged.
	Steps     []Step                 `json:"steps,omitempty"`
	Files     map[string]FileChanges `json:"files,omitempty"` // step id → its changes
	Reads     map[string][]Read      `json:"reads,omitempty"` // step id → its reads
	Decisions []Decision             `json:"decisions,omitempty"`
}

// Source says what the index was built from. It never holds a filesystem
// path.
type Source struct {
	Rev         string     `json:"rev"`    // the revision asked for, e.g. dev
	Head        string     `json:"head"`   // the commit it resolved to
	Merges      int        `json:"merges"` // first-parent merges read
	Rig         *RigSource `json:"rig,omitempty"`
	ReadsMethod string     `json:"reads_method,omitempty"` // ReadsMethod, with a rig
}

// Brane is a named strand of work: one charter row.
type Brane struct {
	Key   string `json:"key"` // row-189
	Row   int    `json:"row"`
	Title string `json:"title"` // the row's slug, "" when the charter has no such row
	Kind  string `json:"kind"`  // attended, loop or mixed
}

// State is one world state. States are listed oldest first, each step just
// before the shared state that merged it.
type State struct {
	ID   string `json:"id"`   // the full commit hash
	Kind string `json:"kind"` // shared or step
	At   string `json:"at"`   // RFC 3339; a step's is its merge's
	PR   int    `json:"pr,omitempty"`
	// Parents: a shared state's are the previous shared state, then the
	// step it merged (the oldest shared state's first parent is its git
	// first parent, outside the window). A step's is its base when the base
	// is a state in the index, else none.
	Parents []string `json:"parents"`
	Brane   string   `json:"brane"`            // a step's brane key; "" on the shared lane
	Base    string   `json:"base,omitempty"`   // step: the merge-base of its merge's parents
	Merge   string   `json:"merge,omitempty"`  // step: the shared state that merged it
	Branch  string   `json:"branch,omitempty"` // the PR branch, owner prefix dropped
}

// Merge is one line of `git log --first-parent --merges
// --format=%H%x09%P%x09%cI%x09%s`.
type Merge struct {
	Hash    string
	Parents []string // exactly two: the first parent, then the PR's side
	At      time.Time
	Subject string
	PR      int    // from "Merge pull request #N from owner/branch"; 0 otherwise
	Branch  string // the branch without its owner; "" when the subject is not a PR merge
}

var (
	commitHash = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)
	prSubject  = regexp.MustCompile(`^Merge pull request #([0-9]+) from [^/\s]+/(\S+)$`)
)

// ParseMerges parses git log output, newest first as git prints it. Blank
// lines are skipped; any other malformed line is an error naming its line.
func ParseMerges(text string) ([]Merge, error) {
	var out []Merge
	for i, line := range strings.Split(text, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		m, err := parseMerge(line)
		if err != nil {
			return nil, fmt.Errorf("surface: merges line %d: %v", i+1, err)
		}
		out = append(out, m)
	}
	return out, nil
}

func parseMerge(line string) (Merge, error) {
	f := strings.SplitN(line, "\t", 4)
	if len(f) != 4 {
		return Merge{}, fmt.Errorf("%d tab-separated fields, want 4 (hash, parents, time, subject)", len(f))
	}
	if !commitHash.MatchString(f[0]) {
		return Merge{}, fmt.Errorf("hash %q is not a full commit hash", f[0])
	}
	parents := strings.Fields(f[1])
	if len(parents) != 2 {
		return Merge{}, fmt.Errorf("%s has %d parents, want 2", f[0], len(parents))
	}
	for _, p := range parents {
		if !commitHash.MatchString(p) {
			return Merge{}, fmt.Errorf("parent %q is not a full commit hash", p)
		}
	}
	at, err := time.Parse(time.RFC3339, f[2])
	if err != nil {
		return Merge{}, fmt.Errorf("time %q is not RFC 3339", f[2])
	}
	m := Merge{Hash: f[0], Parents: parents, At: at, Subject: f[3]}
	if sm := prSubject.FindStringSubmatch(f[3]); sm != nil {
		m.PR, _ = strconv.Atoi(sm[1])
		m.Branch = sm[2]
	}
	return m, nil
}

var branchRow = regexp.MustCompile(`(?:^|[/-])rows?([0-9]+)(?:$|[-/_.])`)

// BraneOf names the brane of a PR branch: the first row the branch names
// (row156, or rows190-191-197 → 190). sprint/ branches are the mission
// loop's; every other prefix (trial/, design/, docs/) is attended. A branch
// that names no row has no brane.
func BraneOf(branch string) (key string, row int, kind string, ok bool) {
	m := branchRow.FindStringSubmatch(branch)
	if m == nil {
		return "", 0, "", false
	}
	row, err := strconv.Atoi(m[1])
	if err != nil || row <= 0 {
		return "", 0, "", false
	}
	kind = BraneAttended
	if strings.HasPrefix(branch, "sprint/") {
		kind = BraneLoop
	}
	return "row-" + strconv.Itoa(row), row, kind, true
}

var queueRow = regexp.MustCompile(`^([0-9]+)\. \*\*`)

// ChartRowTitles reads the charter's queue rows ("NNN. **…") and returns
// each row's slug: the last space-separated token before the line's first
// "** · ", with "**" stripped. Bracketed tags before the slug are skipped
// that way. A row line with no "** · " is left out.
func ChartRowTitles(charter string) map[int]string {
	out := map[int]string{}
	for _, line := range strings.Split(charter, "\n") {
		m := queueRow.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		head, _, found := strings.Cut(line, "** · ")
		if !found {
			continue
		}
		fields := strings.Fields(head)
		slug := strings.ReplaceAll(fields[len(fields)-1], "**", "")
		row, err := strconv.Atoi(m[1])
		if err != nil || slug == "" {
			continue
		}
		if _, dup := out[row]; !dup {
			out[row] = slug
		}
	}
	return out
}

// BuildInput is what Build needs; the caller measures all of it.
type BuildInput struct {
	Merges  []Merge           // newest first, as ParseMerges returns them
	Bases   map[string]string // merge hash → merge-base of its two parents
	Titles  map[int]string    // ChartRowTitles
	BuiltAt time.Time
	Rev     string
	Head    string // the commit Rev resolved to; "" means the newest merge
}

// Build projects the merges into an index: per merge, oldest first, the step
// (its second parent) and then the shared state (the merge). The output
// depends only on the input, never on map order.
func Build(in BuildInput) (Index, error) {
	ix := Index{Schema: Schema, BuiltAt: in.BuiltAt.UTC().Format(time.RFC3339),
		Source: Source{Rev: in.Rev, Head: in.Head, Merges: len(in.Merges)},
		Branes: []Brane{}, States: []State{}}
	if ix.Source.Head == "" && len(in.Merges) > 0 {
		ix.Source.Head = in.Merges[0].Hash
	}
	if 2*len(in.Merges) > MaxStates {
		return Index{}, fmt.Errorf("surface: %d merges make %d states, over the bound of %d", len(in.Merges), 2*len(in.Merges), MaxStates)
	}
	ids := map[string]bool{}
	for _, m := range in.Merges {
		if len(m.Parents) != 2 {
			return Index{}, fmt.Errorf("surface: merge %s has %d parents, want 2", m.Hash, len(m.Parents))
		}
		ids[m.Hash], ids[m.Parents[1]] = true, true
	}
	branes := map[string]*Brane{}
	prevShared := ""
	for i := len(in.Merges) - 1; i >= 0; i-- {
		m := in.Merges[i]
		base, ok := in.Bases[m.Hash]
		if !ok || base == "" {
			return Index{}, fmt.Errorf("surface: merge %s has no merge-base", m.Hash)
		}
		at := m.At.Format(time.RFC3339)
		step := State{ID: m.Parents[1], Kind: KindStep, At: at, PR: m.PR, Parents: []string{},
			Base: base, Merge: m.Hash, Branch: m.Branch}
		if ids[base] {
			step.Parents = []string{base}
		}
		if key, row, kind, ok := BraneOf(m.Branch); ok {
			step.Brane = key
			b := branes[key]
			if b == nil {
				b = &Brane{Key: key, Row: row, Title: in.Titles[row], Kind: kind}
				branes[key] = b
			} else if b.Kind != kind {
				b.Kind = BraneMixed
			}
		}
		first := m.Parents[0]
		if prevShared != "" {
			first = prevShared
		}
		shared := State{ID: m.Hash, Kind: KindShared, At: at, PR: m.PR, Parents: []string{first, m.Parents[1]}}
		ix.States = append(ix.States, step, shared)
		prevShared = m.Hash
	}
	for _, b := range branes {
		ix.Branes = append(ix.Branes, *b)
	}
	sort.Slice(ix.Branes, func(i, j int) bool {
		if ix.Branes[i].Row != ix.Branes[j].Row {
			return ix.Branes[i].Row < ix.Branes[j].Row
		}
		return ix.Branes[i].Key < ix.Branes[j].Key
	})
	return ix, nil
}

// Marshal is the index's file form: indented JSON ending in a newline.
func Marshal(ix Index) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ix); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode reads an index strictly: an unknown field or trailing data is an
// error. It does not validate; call Validate.
func Decode(data []byte) (Index, error) {
	var ix Index
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&ix); err != nil {
		return Index{}, fmt.Errorf("surface: index does not decode: %v", err)
	}
	if dec.More() {
		return Index{}, fmt.Errorf("surface: index has data after the document")
	}
	return ix, nil
}
