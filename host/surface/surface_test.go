package surface

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// The fixtures are a real slice of ailang-world dev (24 first-parent merges
// ending at adffe4c1, newest first), each merge's merge-base, and five
// charter rows copied verbatim from design_docs/world-mission.md.

const goldenPath = "testdata/index.json"

func readFixture(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// fixtureInput is the golden's BuildInput.
func fixtureInput(t *testing.T) BuildInput {
	t.Helper()
	merges, err := ParseMerges(readFixture(t, "merges.txt"))
	if err != nil {
		t.Fatal(err)
	}
	bases := map[string]string{}
	for i, line := range strings.Split(strings.TrimSpace(readFixture(t, "bases.txt")), "\n") {
		f := strings.Fields(line)
		if len(f) != 2 {
			t.Fatalf("bases.txt line %d: %q", i+1, line)
		}
		bases[f[0]] = f[1]
	}
	return BuildInput{Merges: merges, Bases: bases, Titles: ChartRowTitles(readFixture(t, "charter-rows.md")),
		BuiltAt: time.Date(2026, 10, 20, 15, 0, 0, 0, time.UTC), Rev: "dev",
		Head: "adffe4c197de2c2af3ed39b9c7f272b784f6f0e9"}
}

func TestParseMergesFixture(t *testing.T) {
	merges, err := ParseMerges(readFixture(t, "merges.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(merges) != 24 {
		t.Fatalf("%d merges, want 24", len(merges))
	}
	m := merges[3]
	if m.Hash != "90035557b2012895551867c0ebe2f36b853aac09" || m.PR != 327 || m.Branch != "trial/row189-m5-plugin" ||
		!reflect.DeepEqual(m.Parents, []string{"7fb99dae21d0416b09eed95d2da1bfe04e95aa1d", "f36428abd9a72e5093b61cba14e716c37019f1ef"}) ||
		m.At.Format(time.RFC3339) != "2026-10-20T10:31:33+02:00" {
		t.Errorf("merges[3] = %+v", m)
	}
}

func TestParseMergesErrors(t *testing.T) {
	const h1, h2, h3 = "1111111111111111111111111111111111111111", "2222222222222222222222222222222222222222", "3333333333333333333333333333333333333333"
	good := h1 + "\t" + h2 + " " + h3 + "\t2026-10-20T10:00:00+02:00\tMerge pull request #1 from o/trial/row1-x"
	for _, tc := range []struct{ bad, want string }{
		{h1 + "\t" + h2 + " " + h3 + "\t2026-10-20T10:00:00+02:00", "3 tab-separated fields"},
		{"abc123\t" + h2 + " " + h3 + "\t2026-10-20T10:00:00+02:00\ts", `hash "abc123"`},
		{h1 + "\t" + h2 + "\t2026-10-20T10:00:00+02:00\ts", "has 1 parents"},
		{h1 + "\t" + h2 + " " + h3 + " " + h1 + "\t2026-10-20T10:00:00+02:00\ts", "has 3 parents"},
		{h1 + "\t" + h2 + " XYZ\t2026-10-20T10:00:00+02:00\ts", `parent "XYZ"`},
		{h1 + "\t" + h2 + " " + h3 + "\t20 Oct 2026\ts", "not RFC 3339"},
	} {
		// The bad line is line 3: a good line and a blank line come first.
		_, err := ParseMerges(good + "\n\n" + tc.bad + "\n")
		if err == nil || !strings.Contains(err.Error(), "line 3:") || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%q: err %v, want line 3 and %q", tc.bad, err, tc.want)
		}
	}
	// A merge that is not a PR merge parses, with no PR and no branch.
	ms, err := ParseMerges(h1 + "\t" + h2 + " " + h3 + "\t2026-10-20T10:00:00Z\tMerge branch 'x' into dev\n")
	if err != nil || len(ms) != 1 || ms[0].PR != 0 || ms[0].Branch != "" {
		t.Errorf("non-PR merge: %+v, %v", ms, err)
	}
}

func TestBraneOf(t *testing.T) {
	type res struct {
		key  string
		row  int
		kind string
		ok   bool
	}
	for _, tc := range []struct {
		branch string
		want   res
	}{
		{"trial/row189-m5-plugin", res{"row-189", 189, BraneAttended, true}},
		{"trial/row189-m2-m3-episode-start", res{"row-189", 189, BraneAttended, true}},
		{"sprint/world-iter288-row156-m5b", res{"row-156", 156, BraneLoop, true}},
		{"sprint/world-iter282-row156-detection", res{"row-156", 156, BraneLoop, true}},
		{"design/row189-world-onboarding", res{"row-189", 189, BraneAttended, true}},
		{"docs/row198-episode-report", res{"row-198", 198, BraneAttended, true}},
		{"trial/rows190-191-197-read-ranges-delete-tmpdir", res{"row-190", 190, BraneAttended, true}},
		{"trial/row139-tool-pin-v0.54.15", res{"row-139", 139, BraneAttended, true}},
		{"docs/rulings-d84-d85", res{}},
		{"design/world-os-surface", res{}},
		{"docs/evidence-loop-d82-rows", res{}},
		{"trial/ledger-arms", res{}},
		{"", res{}},
	} {
		key, row, kind, ok := BraneOf(tc.branch)
		if got := (res{key, row, kind, ok}); got != tc.want {
			t.Errorf("BraneOf(%q) = %+v, want %+v", tc.branch, got, tc.want)
		}
	}
}

func TestChartRowTitles(t *testing.T) {
	got := ChartRowTitles(readFixture(t, "charter-rows.md"))
	want := map[int]string{
		101: "w-gofmt-drift-with-no-ci-gate", // no tags
		156: "w-exec-protected-test-paths",
		188: "w-desktop-shell", // two bracketed tags, one holding a markdown link
		189: "w-world-onboarding",
		198: "w-episode-report",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ChartRowTitles = %v, want %v", got, want)
	}
	// A row line with no "** · " and prose lines are left out.
	if got := ChartRowTitles("7. **no separator here\nsee row 8. **x** · y\n"); len(got) != 0 {
		t.Errorf("non-rows: %v", got)
	}
}

func TestBuildMatchesGolden(t *testing.T) {
	ix, err := Build(fixtureInput(t))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Marshal(ix)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Build's output differs from %s:\n%s", goldenPath, got)
	}
	// Deterministic: a second build is byte-identical.
	ix2, _ := Build(fixtureInput(t))
	if again, _ := Marshal(ix2); !bytes.Equal(again, got) {
		t.Error("a second Build differs from the first")
	}
}

func stateByID(ix Index, prefix string) (State, bool) {
	for _, s := range ix.States {
		if strings.HasPrefix(s.ID, prefix) {
			return s, true
		}
	}
	return State{}, false
}

func TestBuildFacts(t *testing.T) {
	ix, err := Build(fixtureInput(t))
	if err != nil {
		t.Fatal(err)
	}
	var shared, steps int
	for _, s := range ix.States {
		if s.Kind == KindShared {
			shared++
		} else {
			steps++
		}
	}
	if shared != 24 || steps != 24 || len(ix.Branes) != 6 || ix.Source.Merges != 24 {
		t.Errorf("%d shared, %d steps, %d branes, source %+v; want 24, 24, 6", shared, steps, len(ix.Branes), ix.Source)
	}
	// §4's merge: two ancestors from two strands.
	if s, _ := stateByID(ix, "9003555"); s.Kind != KindShared ||
		!reflect.DeepEqual(s.Parents, []string{"7fb99dae21d0416b09eed95d2da1bfe04e95aa1d", "f36428abd9a72e5093b61cba14e716c37019f1ef"}) {
		t.Errorf("9003555 = %+v", s)
	}
	if s, _ := stateByID(ix, "f36428a"); s.Kind != KindStep || s.Brane != "row-189" || s.PR != 327 ||
		s.Base != "0ee1945ffba795812f06dce95998b646adfac343" || s.Merge != "90035557b2012895551867c0ebe2f36b853aac09" {
		t.Errorf("f36428a = %+v", s)
	}
	if s, _ := stateByID(ix, "4af9b24"); s.Brane != "row-156" || s.PR != 328 || s.Base != "0ee1945ffba795812f06dce95998b646adfac343" {
		t.Errorf("#328's step = %+v", s)
	}
	// #317 names no row: no brane, on the shared lane. Its base is #311's
	// step, which is in the window, so that is its parent.
	if s, _ := stateByID(ix, "bf0002b"); s.Brane != "" || s.PR != 317 ||
		!reflect.DeepEqual(s.Parents, []string{"de79d9868fc88faead0fa40dc7fd2f8e282de07e"}) {
		t.Errorf("#317's step = %+v", s)
	}
	// A base outside the window is kept, but is no parent.
	if s, _ := stateByID(ix, "dfb9f4d"); s.Base != "6eaa997f0bfc126c55765aa15337515b6342e633" || len(s.Parents) != 0 {
		t.Errorf("#318's step = %+v", s)
	}
	// The window's oldest shared state keeps its git first parent; the next
	// one's first parent is the previous shared state, across direct commits.
	if s, _ := stateByID(ix, "1172ced"); s.Parents[0] != "987512d7941ddc5a1a0d9be04299906511e78c17" {
		t.Errorf("oldest shared = %+v", s)
	}
	if s, _ := stateByID(ix, "02bd141"); s.Parents[0] != "a0c73053f10092fd511ad22fcf6b95d6a3f896aa" {
		t.Errorf("02bd141 = %+v", s)
	}
	want := map[string]Brane{
		"row-189": {Key: "row-189", Row: 189, Title: "w-world-onboarding", Name: "World onboarding", Kind: BraneAttended},
		"row-156": {Key: "row-156", Row: 156, Title: "w-exec-protected-test-paths", Name: "Exec protected test paths", Kind: BraneLoop},
		"row-198": {Key: "row-198", Row: 198, Title: "w-episode-report", Name: "Episode report", Kind: BraneAttended},
		"row-190": {Key: "row-190", Row: 190, Kind: BraneAttended},
	}
	for _, b := range ix.Branes {
		if w, ok := want[b.Key]; ok && b != w {
			t.Errorf("brane %+v, want %+v", b, w)
		}
	}
}

func TestBuildMixedBraneAndRefusals(t *testing.T) {
	in := fixtureInput(t)
	// Make #329 (row 189, trial/) a loop step: row 189 becomes mixed.
	in.Merges = append([]Merge(nil), in.Merges...)
	in.Merges[2].Branch = "sprint/world-iter300-row189-x"
	ix, err := Build(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range ix.Branes {
		if b.Key == "row-189" && b.Kind != BraneMixed {
			t.Errorf("row-189 kind %q, want mixed", b.Kind)
		}
	}
	if err := Validate(ix); err != nil {
		t.Errorf("mixed index refused: %v", err)
	}
	in = fixtureInput(t)
	delete(in.Bases, "90035557b2012895551867c0ebe2f36b853aac09")
	if _, err := Build(in); err == nil || !strings.Contains(err.Error(), "90035557b2012895551867c0ebe2f36b853aac09 has no merge-base") {
		t.Errorf("missing base: %v", err)
	}
}

func goldenIndex(t *testing.T) Index {
	t.Helper()
	raw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	ix, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	return ix
}

func indexOf(t *testing.T, ix Index, prefix string) int {
	t.Helper()
	for i, s := range ix.States {
		if strings.HasPrefix(s.ID, prefix) {
			return i
		}
	}
	t.Fatalf("no state %s", prefix)
	return -1
}

func TestValidateAcceptsGolden(t *testing.T) {
	if err := Validate(goldenIndex(t)); err != nil {
		t.Fatal(err)
	}
}

// Each mutant breaks one rule; Validate must refuse it, naming the field and
// the id.
func TestValidateRefusesMutants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, ix *Index)
		want   string
	}{
		{"schema", func(t *testing.T, ix *Index) { ix.Schema = "world/surface-index/v2" }, "schema"},
		{"built_at", func(t *testing.T, ix *Index) { ix.BuiltAt = "yesterday" }, "built_at"},
		{"duplicate id", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "f36428a")].ID = "90035557b2012895551867c0ebe2f36b853aac09"
		}, "id 90035557b2012895551867c0ebe2f36b853aac09 duplicates"},
		{"dangling parent mid-chain", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "0ee1945")].Parents[0] = "ffffffffffffffffffffffffffffffffffffffff"
		}, `0ee1945ffba795812f06dce95998b646adfac343: parents[0] "ffffffffffffffffffffffffffffffffffffffff" does not resolve`},
		{"dangling step parent", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "bf0002b")].Parents[0] = "ffffffffffffffffffffffffffffffffffffffff"
		}, `parents[0] "ffffffffffffffffffffffffffffffffffffffff" does not resolve`},
		{"unknown brane", func(t *testing.T, ix *Index) { ix.States[indexOf(t, *ix, "f36428a")].Brane = "row-999" },
			`f36428abd9a72e5093b61cba14e716c37019f1ef: brane "row-999" is not in branes`},
		{"brane on a shared state", func(t *testing.T, ix *Index) { ix.States[indexOf(t, *ix, "9003555")].Brane = "row-189" },
			`90035557b2012895551867c0ebe2f36b853aac09: brane "row-189" on a shared state`},
		{"bad state kind", func(t *testing.T, ix *Index) { ix.States[indexOf(t, *ix, "f36428a")].Kind = "open" },
			`f36428abd9a72e5093b61cba14e716c37019f1ef: kind "open"`},
		{"bad brane kind", func(t *testing.T, ix *Index) { ix.Branes[0].Kind = "shared" }, `kind "shared" is not attended, loop or mixed`},
		{"duplicate brane", func(t *testing.T, ix *Index) { ix.Branes[1].Key = ix.Branes[0].Key }, "is a duplicate"},
		{"time backwards", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "9003555")].At = "2026-10-20T10:00:00+02:00"
		}, "90035557b2012895551867c0ebe2f36b853aac09: at 2026-10-20T10:00:00+02:00 is before its previous shared state"},
		{"step after its merge", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "f36428a")].At = "2026-10-20T11:00:00+02:00"
		}, "f36428abd9a72e5093b61cba14e716c37019f1ef: at 2026-10-20T11:00:00+02:00 is after its merge"},
		{"bad time", func(t *testing.T, ix *Index) { ix.States[0].At = "soon" }, `at "soon" is not RFC 3339`},
		{"merge does not point back", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "f36428a")].Merge = "0ffbb9c9d386f83aaa45b98e4d4b7763a5c5b845"
		}, "f36428abd9a72e5093b61cba14e716c37019f1ef: merge 0ffbb9c9d386f83aaa45b98e4d4b7763a5c5b845 is not a shared state whose second parent is this step"},
		{"merge does not resolve", func(t *testing.T, ix *Index) { ix.States[indexOf(t, *ix, "f36428a")].Merge = "" },
			`f36428abd9a72e5093b61cba14e716c37019f1ef: merge "" does not resolve`},
		{"shared first parent is a step", func(t *testing.T, ix *Index) {
			ix.States[indexOf(t, *ix, "9003555")].Parents[0] = "4af9b24109b2f5cd376dbf91e32daf9774863c48"
		}, "parents[0] 4af9b24109b2f5cd376dbf91e32daf9774863c48 is a step"},
		{"short id", func(t *testing.T, ix *Index) { ix.States[0].ID = "9003555" }, `id "9003555" is not a full commit hash`},
		{"too many states", func(t *testing.T, ix *Index) {
			for len(ix.States) <= MaxStates {
				ix.States = append(ix.States, ix.States[0])
			}
		}, "over the bound of 5000"},
	} {
		ix := goldenIndex(t)
		tc.mutate(t, &ix)
		err := Validate(ix)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err %v, want %q", tc.name, err, tc.want)
		}
	}
	// Control: the oldest shared state's first parent lies outside the
	// window, and that alone is not a refusal.
	ix := goldenIndex(t)
	if p := ix.States[indexOf(t, ix, "1172ced")].Parents[0]; ix.States[0].ID == p || indexOfOK(ix, p) {
		t.Errorf("control: the oldest first parent %s is in the window", p)
	}
}

func indexOfOK(ix Index, id string) bool {
	for _, s := range ix.States {
		if s.ID == id {
			return true
		}
	}
	return false
}

func TestDecodeIsStrict(t *testing.T) {
	if _, err := Decode([]byte(`{"schema":"world/surface-index/v1","extra":1}`)); err == nil {
		t.Error("an unknown field decoded")
	}
	if _, err := Decode([]byte(`{"schema":"world/surface-index/v1"} {}`)); err == nil {
		t.Error("trailing data decoded")
	}
}

// The index names no filesystem path, home directory or email.
func TestGoldenHoldsNoPaths(t *testing.T) {
	raw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"/Users/", "/home/", "@", "~/"} {
		if bytes.Contains(raw, []byte(bad)) {
			t.Errorf("%s holds %q", goldenPath, bad)
		}
	}
}
