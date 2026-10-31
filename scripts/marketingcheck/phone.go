package main

// Row 231 (design_docs/planned/w-marketing-position-checker.md §2.1–§2.2): the P0/P3 phone
// predicate (orchestration §4.4 phone (i), (ii), (iii′)), ported from the banked oracle
// design_docs/verification/world-iter315/v66-phonecheck-main.go.txt, r3 only. The V58 scan
// functions are copied unchanged. The one deliberate change: (ii) is CheckRegion (the M7b region
// gate), not the oracle's std-dev stand-in, so a 15 px header band is refused (the parity
// ledger's one entry, header-15px).

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
)

const phoneUsage = "usage: marketingcheck phone LABEL=PNG..."

type phoneRect struct{ X0, Y0, X1, Y1 int }

func phoneLuma(im image.Image, x, y int) float64 {
	r, g, b, _ := im.At(x, y).RGBA()
	return 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
}

// phoneBands is the V58 scan: maximal runs of rows with luma < 110 ink in x < 370, above the
// bottom 20 px.
func phoneBands(im image.Image) []phoneRect {
	bd := im.Bounds()
	var out []phoneRect
	in := false
	var cur phoneRect
	for y := bd.Min.Y; y < bd.Max.Y-20; y++ {
		x0, x1 := -1, -1
		for x := bd.Min.X; x < bd.Min.X+370 && x < bd.Max.X; x++ {
			if phoneLuma(im, x, y) < 110 {
				if x0 < 0 {
					x0 = x
				}
				x1 = x
			}
		}
		if x0 >= 0 {
			if !in {
				cur = phoneRect{x0, y, x1 + 1, y + 1}
				in = true
			} else {
				cur.X0 = min(cur.X0, x0)
				cur.X1 = max(cur.X1, x1+1)
				cur.Y1 = y + 1
			}
		} else if in {
			out = append(out, cur)
			in = false
		}
	}
	if in {
		out = append(out, cur)
	}
	return out
}

// phoneScrollbar finds the thumb: rows in the bottom 20 px whose longest contiguous luma
// 120..200 run is at least 100 px long.
func phoneScrollbar(im image.Image) phoneRect {
	bd := im.Bounds()
	r := phoneRect{-1, -1, -1, -1}
	for y := bd.Max.Y - 20; y < bd.Max.Y; y++ {
		x0, x1 := -1, -1
		for x, start := 0, -1; x <= bd.Max.X; x++ {
			in := false
			if x < bd.Max.X {
				l := phoneLuma(im, x, y)
				in = l >= 120 && l <= 200
			}
			if in && start < 0 {
				start = x
			}
			if !in && start >= 0 {
				if x0 < 0 || (x-1)-start > x1-x0 {
					x0, x1 = start, x-1
				}
				start = -1
			}
		}
		if x0 >= 0 && x1-x0 >= 100 {
			if r.Y0 < 0 {
				r = phoneRect{x0, y, x1 + 1, y + 1}
			} else {
				r.Y1 = y + 1
			}
		}
	}
	return r
}

func phoneTitleRight(im image.Image, h phoneRect) int {
	right := -1
	for y := h.Y0; y < h.Y1; y++ {
		for x := 0; x < im.Bounds().Max.X; x++ {
			if phoneLuma(im, x, y) < 110 && x > right {
				right = x
			}
		}
	}
	return right
}

// PhoneScan is the V58 scan's reading of one phone image.
type PhoneScan struct {
	W, H, Bands        int
	Header, FirstEntry phoneRect
	Scrollbar          phoneRect
	TitleRight         int
}

const phoneEntryStart = 9

func ScanPhone(im image.Image) PhoneScan {
	b := phoneBands(im)
	s := PhoneScan{W: im.Bounds().Dx(), H: im.Bounds().Dy(), Bands: len(b), Scrollbar: phoneScrollbar(im), TitleRight: -1}
	if len(b) > 0 {
		s.Header = b[0]
		s.TitleRight = phoneTitleRight(im, b[0])
	}
	if len(b) >= phoneEntryStart+6 {
		u := b[phoneEntryStart]
		for _, r := range b[phoneEntryStart : phoneEntryStart+6] {
			u = phoneRect{min(u.X0, r.X0), min(u.Y0, r.Y0), max(u.X1, r.X1), max(u.Y1, r.Y1)}
		}
		s.FirstEntry = u
	}
	return s
}

func phoneNear(a, b phoneRect, tol int) bool {
	d := func(x, y int) bool { return x-y <= tol && y-x <= tol }
	return d(a.X0, b.X0) && d(a.Y0, b.Y0) && d(a.X1, b.X1) && d(a.Y1, b.Y1)
}

// The R1 phone baseline (r1-rects.json phone.png content, x,y,w,h -> X0,Y0,X1,Y1).
var phoneR1Header = phoneRect{25, 45, 25 + 345, 45 + 23}
var phoneR1Entry = phoneRect{41, 428, 41 + 323, 428 + 182}

func toRect(name string, r phoneRect) Rect { return Rect{name, r.X0, r.Y0, r.X1 - r.X0, r.Y1 - r.Y0} }

// PhoneJudge returns the r3 refusal reasons (empty = ACCEPT).
func PhoneJudge(im image.Image, s PhoneScan) []string {
	var why []string
	if s.Bands == 0 {
		why = append(why, "zero-bands: no ink rows read")
	}
	if !(phoneNear(s.Header, phoneR1Header, 4) && phoneNear(s.FirstEntry, phoneR1Entry, 4)) {
		why = append(why, fmt.Sprintf("(i) header %v entry %v vs R1 %v %v ±4", s.Header, s.FirstEntry, phoneR1Header, phoneR1Entry))
	}
	for _, r := range []Rect{toRect("header", s.Header), toRect("timeline-first", s.FirstEntry)} {
		if err := CheckRegion("phone.png", im, r); err != nil {
			why = append(why, "(ii) "+err.Error())
		}
	}
	if !(s.W == 390 && s.H == 844) {
		why = append(why, fmt.Sprintf("(iii′) dims %dx%d", s.W, s.H))
	}
	if !(s.Scrollbar.Y0 >= 0 && s.Scrollbar.X0 <= 8 && s.Scrollbar.X1 <= 384) {
		why = append(why, fmt.Sprintf("(iii′) thumb %v", s.Scrollbar))
	}
	if !(s.TitleRight >= 387) {
		why = append(why, fmt.Sprintf("(iii′) title-right %d", s.TitleRight))
	}
	return why
}

// runPhone is `marketingcheck phone`: rc 0 when every label is ACCEPT, 1 when any is REFUSE or
// an input fails to decode, 2 on a usage error (design §2.1).
func runPhone(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("phone", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprintln(stderr, phoneUsage) }
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	cands, err := parseLabelled(fs.Args())
	if err != nil {
		fmt.Fprintf(stderr, "marketingcheck: %v\n", err)
		fs.Usage()
		return 2
	}
	refused := 0
	for _, c := range cands {
		b, err := os.ReadFile(c[1])
		if err != nil {
			fmt.Fprintf(stderr, "marketingcheck: decode-failed %s: %s\n", c[1], err)
			return 1
		}
		im, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			fmt.Fprintf(stderr, "marketingcheck: decode-failed %s: %s\n", c[1], err)
			return 1
		}
		s := ScanPhone(im)
		why := PhoneJudge(im, s)
		v := "ACCEPT"
		if len(why) > 0 {
			v = "REFUSE"
			refused++
		}
		fmt.Fprintf(stdout, "%-28s %s %dx%d bands=%d header=%v entry=%v thumb=%v titleRight=%d\n", c[0], v, s.W, s.H, s.Bands, s.Header, s.FirstEntry, s.Scrollbar, s.TitleRight)
		for _, w := range why {
			fmt.Fprintln(stdout, "    "+w)
		}
	}
	fmt.Fprintf(stdout, "TOTAL candidates=%d accept=%d refuse=%d\n", len(cands), len(cands)-refused, refused)
	if refused > 0 {
		return 1
	}
	return 0
}
