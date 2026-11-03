package main

// Row 240 (design_docs/planned/w-marketing-checker-tolerance-pins.md §2): one boundary pair per
// row-231 tolerance that no other test pinned, the R1 premise under which removing `gc >= 0`
// is equivalent, and the main.go routes for position and phone. Inputs are the committed seed-5
// store and skeleton; production tolerances are unchanged.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// positionGapPremise checks that every non-free gap of base is at least 2 px, the condition under
// which `gc >= 0` in (a3)'s non-free branch is implied by `posAbs(e) <= 2` (design §2.5). It also
// requires at least one non-free gap, so it cannot pass on an empty baseline.
func positionGapPremise(base Rectangles) error {
	n := 0
	for _, s := range positionShots {
		bs := base.Shots[s].Content
		for j := 1; j < len(bs); j++ {
			if posFreeGap(s, bs[j].Name) {
				continue
			}
			n++
			if gb := bs[j].Y - (bs[j-1].Y + bs[j-1].H); gb < 2 {
				return fmt.Errorf("%s/%s: non-free baseline gap %d < 2; `gc >= 0` is no longer implied (re-derive design §2.5)", s, bs[j].Name, gb)
			}
		}
	}
	if n == 0 {
		return fmt.Errorf("baseline has no non-free gap")
	}
	return nil
}

func TestMarketingPositionTolerancePins(t *testing.T) {
	base := loadPositionBaseline(t)
	t.Run("gc-a3-conjunct-source-pin", func(t *testing.T) {
		// T18 (removing `gc >= 0 &&` from the non-free (a3) gap) is banked as
		// equivalent on R1 (design §2.5). Editing that clause re-opens the
		// classification, so the clause's exact text is pinned here.
		src, err := os.ReadFile("position.go")
		if err != nil {
			t.Fatal(err)
		}
		const clause = "okGap, gRes := gc >= 0 && posAbs(e) <= 2, e"
		if n := strings.Count(string(src), clause); n != 1 {
			t.Errorf("position.go holds %d copies of the (a3) clause %q, want 1: T18's equivalence (design §2.5) must be re-derived", n, clause)
		}
	})
	t.Run("gc-premise-nonfree-gaps-ge-2", func(t *testing.T) {
		if err := positionGapPremise(base); err != nil {
			t.Error(err)
		}
	})
	src, err := os.ReadFile(filepath.Join(positionTestdata, "stores", "r1c-seed5.json"))
	if err != nil {
		t.Fatal(err)
	}
	bg := func(f func(r *Rect)) func(d *Rectangles) {
		return func(d *Rectangles) {
			sh := d.Shots["workbench-light.png"]
			f(&sh.Background)
			d.Shots["workbench-light.png"] = sh
		}
	}
	const bgRefuse = "workbench-light.png/background: edge >2px"
	for _, c := range []struct {
		name string
		mut  func(d *Rectangles)
		want []string
	}{
		{"bg-x0+2", bg(func(r *Rect) { r.X += 2; r.W -= 2 }), nil},
		{"bg-x0+3", bg(func(r *Rect) { r.X += 3; r.W -= 3 }), []string{bgRefuse}},
		{"bg-y+2", bg(func(r *Rect) { r.Y += 2; r.H -= 2 }), nil},
		{"bg-y+3", bg(func(r *Rect) { r.Y += 3; r.H -= 3 }), []string{bgRefuse}},
		{"bg-right-2", bg(func(r *Rect) { r.W -= 2 }), nil},
		{"bg-right-3", bg(func(r *Rect) { r.W -= 3 }), []string{bgRefuse}},
		{"bg-bottom-2", bg(func(r *Rect) { r.H -= 2 }), nil},
		{"bg-bottom-3", bg(func(r *Rect) { r.H -= 3 }), []string{bgRefuse}},
		{"a4-graph-svg-h+2", func(d *Rectangles) { d.Shots["graph.png"].Content[0].H += 2 }, nil},
		{"a4-graph-svg-h+3", func(d *Rectangles) { d.Shots["graph.png"].Content[0].H += 3 },
			[]string{"graph.png/graph-svg: (a4) height 277 vs 274"}},
		{"a3-timeline-gap+2", func(d *Rectangles) {
			sh := d.Shots["workbench-light.png"]
			sh.Content[posIdx(sh, "header")].H -= 1
			sh.Content[posIdx(sh, "timeline-first")].Y += 1
		}, nil},
		{"a3-timeline-gap+3", func(d *Rectangles) {
			sh := d.Shots["workbench-light.png"]
			sh.Content[posIdx(sh, "header")].H -= 1
			sh.Content[posIdx(sh, "timeline-first")].Y += 2
		}, []string{"workbench-light.png/timeline-first: (a3) gap after header 103 vs R1 100"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			var d Rectangles
			if err := json.Unmarshal(src, &d); err != nil {
				t.Fatal(err)
			}
			c.mut(&d)
			reasons, regions := PositionJudge(base, d)
			if regions != 17 || strings.Join(reasons, "\n") != strings.Join(c.want, "\n") {
				t.Errorf("regions %d reasons %q, want 17 and %q", regions, reasons, c.want)
			}
		})
	}
}

func TestMarketingPhoneEntryTolerancePins(t *testing.T) {
	dir := t.TempDir()
	for _, c := range []struct {
		dy   int
		want string
	}{{4, "ACCEPT"}, {5, "REFUSE (i)"}} {
		t.Run(fmt.Sprintf("entry+%d", c.dy), func(t *testing.T) {
			s := loadPhoneSkeleton(t, filepath.Join(phoneTestdata, "skeleton-r1c-seed5-pass-a-phone.json"))
			for i := phoneEntryStart; i < phoneEntryStart+6; i++ {
				s.Bands[i].Y0 += c.dy
				s.Bands[i].Y1 += c.dy
			}
			p := filepath.Join(dir, fmt.Sprintf("entry+%d.png", c.dy))
			writePhonePNG(t, p, renderPhoneSkeleton(s))
			reasons, sc := judgePhoneFile(t, p)
			want := phoneRect{phoneR1Entry.X0, phoneR1Entry.Y0 + c.dy, phoneR1Entry.X1, phoneR1Entry.Y1 + c.dy}
			if sc.Header != phoneR1Header || sc.FirstEntry != want || phoneVerdict(reasons) != c.want {
				t.Errorf("header %v entry %v reasons %q; want header %v entry %v %s", sc.Header, sc.FirstEntry, reasons, phoneR1Header, want, c.want)
			}
			t.Logf("entry %v %s %q", sc.FirstEntry, phoneVerdict(reasons), reasons)
		})
	}
}

// TestMarketingRunRoutesPositionPhone pins main.go's position and phone cases: a removed case
// falls through to usage (rc 2), and a swapped one rejects the other subcommand's arguments.
// run writes to os.Stdout/os.Stderr; it takes no writers, so the verdict lines go to test output.
func TestMarketingRunRoutesPositionPhone(t *testing.T) {
	if rc := run([]string{"position", "--baseline", positionBaseline, "a=" + filepath.Join(positionTestdata, "stores", "r1c-seed5.json")}); rc != 0 {
		t.Errorf("run position on the seed-5 store: rc %d, want 0", rc)
	}
	p := filepath.Join(t.TempDir(), "seed5.png")
	writePhonePNG(t, p, renderPhoneSkeleton(loadPhoneSkeleton(t, filepath.Join(phoneTestdata, "skeleton-r1c-seed5-pass-a-phone.json"))))
	if rc := run([]string{"phone", "a=" + p}); rc != 0 {
		t.Errorf("run phone on the seed-5 render: rc %d, want 0", rc)
	}
}
