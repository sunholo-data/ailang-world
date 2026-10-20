package surface

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// The rig fixture (testdata/rig) is synthetic, in the rig's formats, with
// placeholder paths (/rig/…): episode e1 (two rounds, a report with
// refusals) built #327's step, e2 (one round, no report) #328's; e3 names no
// PR, e4 another project's, e5 a PR outside the window; p1 is another
// project's record and is never read. Git is faked.

const rigGoldenPath = "testdata/index-rig.json"

const (
	stepE1  = "f36428abd9a72e5093b61cba14e716c37019f1ef" // #327, row-189 (attended)
	stepE2  = "4af9b24109b2f5cd376dbf91e32daf9774863c48" // #328, row-156 (loop)
	base0   = "0ee1945ffba795812f06dce95998b646adfac343"
	step316 = "59bbded417388861f3acae45c07eb9046f3cacaa" // #316's step
	hashA   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

// fakeGit answers from maps and counts its calls.
type fakeGit struct {
	nameStatus map[string]string // step → -z output
	tree       map[string]string // rev → -z output
	last       map[string]string // path → "%H %P"
	calls      map[string]int
}

func (g *fakeGit) NameStatus(base, step string) (string, error) {
	g.calls["name-status"]++
	return g.nameStatus[step], nil
}

func (g *fakeGit) Tree(rev string) (string, error) {
	g.calls["tree"]++
	return g.tree[rev], nil
}

func (g *fakeGit) LastChange(rev, path string) (string, error) {
	g.calls["last:"+path]++
	return g.last[path], nil
}

func newFakeGit() *fakeGit {
	return &fakeGit{
		nameStatus: map[string]string{
			stepE1: "M\x00CLAUDE.md\x00A\x00host/surface/surface.go\x00T\x00plugins/ailang-world/README.md\x00D\x00old/notes.md\x00",
			stepE2: "M\x00host/broker/exec.go\x00A\x00host/broker/exec_test.go\x00",
		},
		tree: map[string]string{
			base0: "CLAUDE.md\x00cmd/ailang-worldd/main.go\x00design_docs/planned/w-world-onboarding.md\x00",
		},
		last: map[string]string{
			// #316's merge: its second parent is a step in the window.
			"design_docs/planned/w-world-onboarding.md": "53862152280c2928a8a6c1fd78b2c362573f8450 7a679bfed34fc3fccfe37677c699ba02721b4ca1 " + step316,
			"cmd/ailang-worldd/main.go":                 hashA + " bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			// A merge before the window.
			"CLAUDE.md": "cccccccccccccccccccccccccccccccccccccccc dddddddddddddddddddddddddddddddddddddddd eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		},
		calls: map[string]int{},
	}
}

// rigFixture reads the fixture rig the way surface-index does: the
// project's round directories, prs.json and e1's report.
func rigFixture(t *testing.T) RigInput {
	t.Helper()
	var records []EpisodeRecord
	dirs, err := os.ReadDir("testdata/rig/evidence")
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range dirs {
		if !RoundDir("ailang-world", d.Name()) {
			continue
		}
		files, _ := filepath.Glob(filepath.Join("testdata/rig/evidence", d.Name(), "episodes", "*.json"))
		for _, f := range files {
			raw, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			r, err := ParseEpisodeRecord(raw)
			if err != nil {
				t.Fatalf("%s: %v", f, err)
			}
			records = append(records, r)
		}
	}
	prs, err := ParsePRs([]byte(readFixture(t, "rig/trial/prs.json")), "ailang-world")
	if err != nil {
		t.Fatal(err)
	}
	rep, err := ParseReport([]byte(readFixture(t, "rig/reports/e1.json")))
	if err != nil {
		t.Fatal(err)
	}
	return RigInput{Episodes: GroupRounds(records), PRs: prs, Reports: map[string]Report{"e1": rep}}
}

func rigIndex(t *testing.T) (Index, *fakeGit) {
	t.Helper()
	ix, err := Build(fixtureInput(t))
	if err != nil {
		t.Fatal(err)
	}
	g := newFakeGit()
	ix, err = AttachRig(ix, rigFixture(t), g)
	if err != nil {
		t.Fatal(err)
	}
	return ix, g
}

func TestParseEpisodeRecordErrors(t *testing.T) {
	good := `{"schema":"trial-episode/1","episode":"e1","started_at":"2026-10-20T07:00:00+00:00","finished_at":"2026-10-20T07:20:00+00:00","agent_run":{"wall_ms":1}}`
	if _, err := ParseEpisodeRecord([]byte(good)); err != nil {
		t.Fatalf("good record: %v", err)
	}
	for _, tc := range []struct{ from, to, want string }{
		{`"trial-episode/1"`, `"trial-episode/2"`, "schema"},
		{`"episode":"e1"`, `"episode":""`, "not an episode id"},
		{`"episode":"e1"`, `"episode":"../e1"`, "not an episode id"},
		{`"2026-10-20T07:00:00+00:00"`, `"7am"`, "started_at"},
		{`"2026-10-20T07:20:00+00:00"`, `""`, "finished_at"},
		{`"wall_ms":1`, `"wall_ms":-1`, "negative"},
		{`{"schema"`, `[{"schema"`, "does not decode"},
	} {
		_, err := ParseEpisodeRecord([]byte(strings.Replace(good, tc.from, tc.to, 1)))
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s → %s: err %v, want %q", tc.from, tc.to, err, tc.want)
		}
	}
}

func TestGroupRoundsAndRoundDir(t *testing.T) {
	in := rigFixture(t)
	var ids []string
	for _, ep := range in.Episodes {
		ids = append(ids, ep.ID)
	}
	// p1's directory is another project's: never read.
	if want := []string{"e1", "e2", "e3", "e4", "e5"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("episodes %v, want %v", ids, want)
	}
	e1 := in.Episodes[0]
	if len(e1.Rounds) != 2 || e1.Rounds[0].Task != "row-189-M5-plugin" || e1.Rounds[1].Task != "row-189-M5-plugin-round-2" {
		t.Errorf("e1 rounds: %+v", e1.Rounds)
	}
	// Rounds are ordered by started_at, whatever the input order.
	swapped := GroupRounds([]EpisodeRecord{e1.Rounds[1], e1.Rounds[0]})
	if swapped[0].Rounds[0].StartedAt != "2026-10-20T07:00:00+00:00" {
		t.Errorf("rounds not ordered by started_at: %+v", swapped[0].Rounds)
	}
	for name, want := range map[string]bool{
		"ailang-world": true, "ailang-world-r2": true, "ailang-world-r13": true,
		"ailang-world-r1": false, "ailang-world-r02": false, "ailang-world-rx": false, "ailang-world-r": false,
		"ailang-world-extra": false, "platform-source": false, "platform-source-r2": false,
	} {
		if got := RoundDir("ailang-world", name); got != want {
			t.Errorf("RoundDir(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestParsePRs(t *testing.T) {
	got, err := ParsePRs([]byte(readFixture(t, "rig/trial/prs.json")), "ailang-world")
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]int{"e1": 327, "e2": 328, "e5": 999}; !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePRs = %v, want %v", got, want)
	}
	got, _ = ParsePRs([]byte(`{"a":"ailang-world#0","b":"ailang-world#12x","c":7,"d":"xailang-world#3","e":"ailang-world#4"}`), "ailang-world")
	if want := map[string]int{"e": 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePRs edge cases = %v, want %v", got, want)
	}
	if _, err := ParsePRs([]byte(`[]`), "ailang-world"); err == nil {
		t.Error("a non-object prs.json parsed")
	}
}

func TestParseNameStatus(t *testing.T) {
	got, err := ParseNameStatus("M\x00b.go\x00A\x00a.go\x00T\x00c\x00D\x00d/e.md\x00")
	want := []Change{{"a.go", "A"}, {"b.go", "M"}, {"c", "M"}, {"d/e.md", "D"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ParseNameStatus = %v, %v; want %v", got, err, want)
	}
	if got, err := ParseNameStatus(""); err != nil || got == nil || len(got) != 0 {
		t.Errorf("empty diff = %v, %v", got, err)
	}
	for _, bad := range []string{"R100\x00a\x00", "M\x00a\x00M\x00", "M\x00/abs\x00", "M\x00../up\x00"} {
		if _, err := ParseNameStatus(bad); err == nil {
			t.Errorf("ParseNameStatus(%q) accepted", bad)
		}
	}
}

func TestClassifyLastChange(t *testing.T) {
	steps := map[string]bool{step316: true}
	for _, tc := range []struct{ line, via, from string }{
		{"", ViaOutside, ""},
		{hashA, ViaCommit, hashA}, // a root commit
		{hashA + " " + base0, ViaCommit, hashA},
		{hashA + " " + base0 + " " + step316, ViaStep, step316},
		{hashA + " " + base0 + " " + stepE1, ViaOutside, ""},
	} {
		via, from, err := ClassifyLastChange(tc.line, steps)
		if err != nil || via != tc.via || from != tc.from {
			t.Errorf("ClassifyLastChange(%q) = %s, %s, %v; want %s, %s", tc.line, via, from, err, tc.via, tc.from)
		}
	}
	if _, _, err := ClassifyLastChange("not-a-hash", steps); err == nil {
		t.Error("a malformed line classified")
	}
}

func TestAttachRigMatchesGolden(t *testing.T) {
	ix, _ := rigIndex(t)
	if err := Validate(ix); err != nil {
		t.Fatalf("the rig index is invalid: %v", err)
	}
	got, err := Marshal(ix)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(rigGoldenPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("AttachRig's output differs from %s:\n%s", rigGoldenPath, got)
	}
	again, _ := rigIndex(t)
	if b, _ := Marshal(again); !bytes.Equal(b, got) {
		t.Error("a second AttachRig differs from the first")
	}
	// The golden decodes strictly and validates.
	dec, err := Decode(want)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(dec); err != nil {
		t.Errorf("the rig golden is invalid: %v", err)
	}
	// Without a rig the output is M0a's, byte for byte.
	plain, _ := Build(fixtureInput(t))
	if b, _ := Marshal(plain); !bytes.Equal(b, []byte(readFixture(t, "index.json"))) {
		t.Error("an index built without a rig differs from the M0a golden")
	}
}

func TestAttachRigFacts(t *testing.T) {
	ix, g := rigIndex(t)
	if len(ix.Steps) != 2 {
		t.Fatalf("%d steps, want 2: %+v", len(ix.Steps), ix.Steps)
	}
	// Sorted by the step state's at: #328 (08:22Z) before #327 (08:31Z).
	e2, e1 := ix.Steps[0], ix.Steps[1]
	calls, entries := 42, [2]int64{1200, 1290}
	want1 := Step{Episode: "e1", State: stepE1, PR: 327, Task: "row-189-M5-plugin", Agent: "claude", Model: "claude-opus-5-5",
		Rounds: 2, Graders: Graders{4, 4}, Minutes: 13.5, Calls: &calls, Entries: &entries, Base: base0}
	if !reflect.DeepEqual(e1, want1) {
		t.Errorf("e1 = %+v, want %+v", e1, want1)
	}
	// 255000 ms is 4.25 minutes, rounded to 4.3; no report, so no calls or
	// entries; a task file keeps its base name only.
	want2 := Step{Episode: "e2", State: stepE2, PR: 328, Task: "row-156-m5b.md", Agent: "codex", Model: "gpt-x",
		Rounds: 1, Graders: Graders{2, 2}, Minutes: 4.3, Base: base0}
	if !reflect.DeepEqual(e2, want2) {
		t.Errorf("e2 = %+v, want %+v", e2, want2)
	}
	rig := ix.Source.Rig
	if rig == nil || rig.Episodes != 5 || rig.Attached != 2 || !reflect.DeepEqual(rig.Unattached, []string{"e3", "e4", "e5"}) ||
		ix.Source.ReadsMethod != ReadsMethod {
		t.Errorf("source = %+v, rig %+v", ix.Source, rig)
	}
	wantFiles := []Change{{"CLAUDE.md", "M"}, {"host/surface/surface.go", "A"}, {"old/notes.md", "D"}, {"plugins/ailang-world/README.md", "M"}}
	if fc := ix.Files[stepE1]; !reflect.DeepEqual(fc.Changes, wantFiles) || fc.Truncated != 0 {
		t.Errorf("files[e1] = %+v", fc)
	}
	// The report's order; "./x" cleans to a path already read and adds to
	// its count; "../x" is never copied.
	wantReads := []Read{
		{"design_docs/planned/w-world-onboarding.md", 3, ViaStep, step316},
		{"plugins/ailang-world/README.md", 1, ViaNew, ""},
		{"cmd/ailang-worldd/main.go", 3, ViaCommit, hashA},
		{"CLAUDE.md", 1, ViaOutside, ""},
	}
	if got := ix.Reads[stepE1]; !reflect.DeepEqual(got, wantReads) {
		t.Errorf("reads[e1] = %+v, want %+v", got, wantReads)
	}
	if _, ok := ix.Reads[stepE2]; ok {
		t.Error("reads for e2, which has no report")
	}
	// git: one tree per step with reads, a last change only for a path at base.
	if g.calls["tree"] != 1 || g.calls["last:plugins/ailang-world/README.md"] != 0 || g.calls["last:cmd/ailang-worldd/main.go"] != 1 || g.calls["name-status"] != 2 {
		t.Errorf("git calls %v", g.calls)
	}
	wantDecisions := []Decision{
		{stepE1, ActorWorld, DecisionGrade, "round 1: 3/4 graders", "2026-10-20T07:20:00Z"},
		{stepE1, ActorWorld, DecisionRefusal, "2 plan + 1 handler refusals", "2026-10-20T07:49:00Z"},
		{stepE1, ActorWorld, DecisionGrade, "round 2: 4/4 graders", "2026-10-20T07:50:00Z"},
		{stepE2, ActorWorld, DecisionGrade, "round 1: 2/2 graders", "2026-10-20T07:58:00Z"},
		{stepE2, ActorLoop, DecisionMerge, "merged #328", "2026-10-20T08:22:02Z"},
		{stepE1, ActorYou, DecisionMerge, "merged #327", "2026-10-20T08:31:33Z"},
	}
	if !reflect.DeepEqual(ix.Decisions, wantDecisions) {
		t.Errorf("decisions = %+v\nwant %+v", ix.Decisions, wantDecisions)
	}
}

func TestAttachRigCapsAndClaims(t *testing.T) {
	base, err := Build(fixtureInput(t))
	if err != nil {
		t.Fatal(err)
	}
	in := rigFixture(t)
	g := newFakeGit()
	var z strings.Builder
	for i := 0; i < MaxChanges+5; i++ {
		z.WriteString("A\x00f/" + strings.Repeat("0", 4-len(itoa(i))) + itoa(i) + "\x00")
	}
	g.nameStatus[stepE1] = z.String()
	// e3 claims #327 too, but started first: e3 is attached, e1 is not.
	in.PRs["e3"] = 327
	ix, err := AttachRig(base, in, g)
	if err != nil {
		t.Fatal(err)
	}
	if fc := ix.Files[stepE1]; len(fc.Changes) != MaxChanges || fc.Truncated != 5 || fc.Changes[0].Path != "f/0000" {
		t.Errorf("files[e1]: %d listed, %d truncated", len(fc.Changes), fc.Truncated)
	}
	if ix.Steps[1].Episode != "e3" || !reflect.DeepEqual(ix.Source.Rig.Unattached, []string{"e1", "e4", "e5"}) {
		t.Errorf("claims: steps %+v, unattached %v", ix.Steps, ix.Source.Rig.Unattached)
	}
	if err := Validate(ix); err != nil {
		t.Errorf("capped index refused: %v", err)
	}
	// No episode attaches: source.rig is present, the facts are absent.
	in.PRs = map[string]int{}
	ix, err = AttachRig(base, in, g)
	if err != nil || ix.Steps != nil || ix.Source.Rig.Attached != 0 || len(ix.Source.Rig.Unattached) != 5 {
		t.Errorf("no attachment: %+v, %v", ix.Source, err)
	}
	if err := Validate(ix); err != nil {
		t.Errorf("an index with an empty rig refused: %v", err)
	}
}

func itoa(i int) string { return strconv.Itoa(i) }

func rigGolden(t *testing.T) Index {
	t.Helper()
	ix, err := Decode([]byte(readFixture(t, "index-rig.json")))
	if err != nil {
		t.Fatal(err)
	}
	return ix
}

// One mutant per validateRig rule; each must be refused, naming the field.
func TestValidateRefusesRigMutants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(ix *Index)
		want   string
	}{
		{"facts without a rig", func(ix *Index) { ix.Source.Rig = nil }, "need source.rig"},
		{"reads method", func(ix *Index) { ix.Source.ReadsMethod = "creating-commit/v1" }, "source.reads_method"},
		{"attached count", func(ix *Index) { ix.Source.Rig.Attached = 3 }, "source.rig.attached 3"},
		{"duplicate episode", func(ix *Index) { ix.Steps[1].Episode = ix.Steps[0].Episode }, "is empty or a duplicate"},
		{"state is not a step", func(ix *Index) { ix.Steps[0].State = "7fb99dae21d0416b09eed95d2da1bfe04e95aa1d" }, "is not a step state"},
		{"state does not resolve", func(ix *Index) { ix.Steps[0].State = hashA }, "is not a step state"},
		{"state used twice", func(ix *Index) { ix.Steps[1].State = ix.Steps[0].State; ix.Steps[1].PR = ix.Steps[0].PR }, "is used twice"},
		{"pr mismatch", func(ix *Index) { ix.Steps[0].PR = 327 }, "pr 327, but state"},
		{"rounds", func(ix *Index) { ix.Steps[0].Rounds = 0 }, "rounds 0, want at least 1"},
		{"passed over total", func(ix *Index) { ix.Steps[0].Graders = Graders{3, 2} }, "graders 3/2"},
		{"files key", func(ix *Index) { ix.Files[step316] = FileChanges{Changes: []Change{}} }, "files[" + step316 + "] is not a step in steps"},
		{"reads key", func(ix *Index) { ix.Reads[step316] = []Read{} }, "reads[" + step316 + "] is not a step in steps"},
		{"decision step", func(ix *Index) { ix.Decisions[0].Step = step316 }, "decisions[0].step"},
		{"status", func(ix *Index) { ix.Files[stepE1].Changes[0].Status = "R" }, `status "R" is not A, M or D`},
		{"files path", func(ix *Index) { ix.Files[stepE1].Changes[0].Path = "/rig/ws/e1/CLAUDE.md" }, "is not local"},
		{"reads path", func(ix *Index) { ix.Reads[stepE1][0].Path = "../CLAUDE.md" }, `path "../CLAUDE.md" is not local`},
		{"via", func(ix *Index) { ix.Reads[stepE1][0].Via = "creator" }, `via "creator"`},
		{"from of via new", func(ix *Index) { ix.Reads[stepE1][1].From = hashA }, "via new has from"},
		{"from of via commit", func(ix *Index) { ix.Reads[stepE1][2].From = "abc" }, "via commit: from"},
		{"from of via step", func(ix *Index) { ix.Reads[stepE1][0].From = "7fb99dae21d0416b09eed95d2da1bfe04e95aa1d" }, "via step: from"},
		{"actor", func(ix *Index) { ix.Decisions[0].Actor = "mark" }, `actor "mark"`},
		{"kind", func(ix *Index) { ix.Decisions[0].Kind = "mint" }, `kind "mint"`},
		{"decision at", func(ix *Index) { ix.Decisions[0].At = "today" }, `at "today"`},
		{"truncated", func(ix *Index) { fc := ix.Files[stepE1]; fc.Truncated = -1; ix.Files[stepE1] = fc }, "-1 truncated"},
	} {
		ix := rigGolden(t)
		tc.mutate(&ix)
		err := Validate(ix)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err %v, want %q", tc.name, err, tc.want)
		}
	}
}

// The rig golden names no filesystem path, home directory or email, though
// every record it came from holds paths.
func TestRigGoldenHoldsNoPaths(t *testing.T) {
	raw := readFixture(t, "index-rig.json")
	for _, bad := range []string{"/rig/", "/Users/", "/home/", "@", "~/", "tasks/", "stdout_path", "worktree", "transcript"} {
		if strings.Contains(raw, bad) {
			t.Errorf("%s holds %q", rigGoldenPath, bad)
		}
	}
}
