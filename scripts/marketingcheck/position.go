package main

// Row 231 (design_docs/planned/w-marketing-position-checker.md §2.1–§2.2): the P0/P3 position
// predicate, ported from the banked oracle design_docs/verification/world-iter315/v64-p0check-main.go.txt
// (`check` mode, r3 only). The reason strings are byte-identical to the oracle's `r3:` lines;
// TestMarketingPositionParity holds them to it. Image geometry comes from the Shots table, which
// TestMarketingPositionGeometryMatchesOracle pins to the oracle's constants.

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

const positionUsage = "usage: marketingcheck position --baseline RECTS LABEL=RECTS..."

const positionLine = 24

var positionShots = []string{"workbench-light.png", "workbench-dark.png", "packet.png", "graph.png", "provenance.png"}

func posAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func posLineStep(d int) bool {
	k := int(math.Round(float64(d) / positionLine))
	return posAbs(d-k*positionLine) <= 2
}

func posPara(shot string) bool { return shot == "packet.png" }

func posFreeGap(shot, name string) bool {
	return shot == "packet.png" && (name == "approvals-chain" || name == "roster-tested")
}

func posInside(c Rect, w, h int) bool {
	return c.X >= 0 && c.Y >= 0 && c.X+c.W <= w && c.Y+c.H <= h-40
}

// PositionJudge returns the r3 refusal reasons (empty = ACCEPT) and the number of baseline
// regions it compared (content plus one background per shot).
func PositionJudge(base, cand Rectangles) (reasons []string, regions int) {
	for _, s := range positionShots {
		bs := base.Shots[s]
		shot, _ := shotByName(s)
		H, W := shot.Height, shot.Width
		cs, ok := cand.Shots[s]
		if !ok {
			reasons = append(reasons, s+": shot missing")
			continue
		}
		cm := map[string]Rect{}
		for _, c := range cs.Content {
			cm[c.Name] = c
		}
		for i, b := range bs.Content {
			regions++
			c, ok := cm[b.Name]
			if !ok {
				reasons = append(reasons, s+"/"+b.Name+": missing")
				continue
			}
			if !posInside(c, W, H) {
				reasons = append(reasons, fmt.Sprintf("%s/%s: (b) bottom %d > H-40 %d or outside", s, b.Name, c.Y+c.H, H-40))
			}
			if !(posAbs(c.X-b.X) <= 2 && posAbs(c.X+c.W-b.X-b.W) <= 2) {
				reasons = append(reasons, fmt.Sprintf("%s/%s: (a1) left/right %d..%d vs %d..%d", s, b.Name, c.X, c.X+c.W, b.X, b.X+b.W))
			}
			dh := c.H - b.H
			if posPara(s) {
				if !(posLineStep(dh) && c.H >= positionLine) {
					reasons = append(reasons, fmt.Sprintf("%s/%s: (a4) height %d vs %d not whole lines", s, b.Name, c.H, b.H))
				}
			} else if posAbs(dh) > 2 {
				reasons = append(reasons, fmt.Sprintf("%s/%s: (a4) height %d vs %d", s, b.Name, c.H, b.H))
			}
			if i == 0 {
				if posAbs(c.Y-b.Y) > 2 {
					reasons = append(reasons, fmt.Sprintf("%s/%s: (a2) stack-head top %d vs %d", s, b.Name, c.Y, b.Y))
				}
				continue
			}
			pb := bs.Content[i-1]
			pc, ok := cm[pb.Name]
			if !ok {
				continue
			}
			gb := b.Y - (pb.Y + pb.H)
			gc := c.Y - (pc.Y + pc.H)
			e := gc - gb
			okGap := gc >= 0 && posAbs(e) <= 2
			if posFreeGap(s, b.Name) {
				okGap = gc >= 0 && posLineStep(e)
			}
			if !okGap {
				reasons = append(reasons, fmt.Sprintf("%s/%s: (a3) gap after %s %d vs R1 %d", s, b.Name, pb.Name, gc, gb))
			}
		}
		regions++
		b, c := bs.Background, cs.Background
		if !(posAbs(c.X-b.X) <= 2 && posAbs(c.Y-b.Y) <= 2 && posAbs(c.X+c.W-b.X-b.W) <= 2 && posAbs(c.Y+c.H-b.Y-b.H) <= 2) {
			reasons = append(reasons, fmt.Sprintf("%s/background: edge >2px", s))
		}
		if !posInside(c, W, H) {
			reasons = append(reasons, fmt.Sprintf("%s/background: (b)", s))
		}
	}
	return reasons, regions
}

// loadPositionCandidate reads one candidate rects.json. A wrong schema or unparsable file is
// bad-rectangles; a file with no content rectangle in the five 1440 px shots is zero-regions.
func loadPositionCandidate(path string) (Rectangles, error) {
	var r Rectangles
	b, err := os.ReadFile(path)
	if err != nil {
		return r, refuse("bad-rectangles", path, err.Error())
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return r, refuse("bad-rectangles", path, err.Error())
	}
	if r.Schema != "world/demo/capture-rectangles/v1" {
		return r, refuse("bad-rectangles", path, "wrong schema")
	}
	n := 0
	for _, s := range positionShots {
		n += len(r.Shots[s].Content)
	}
	if n == 0 {
		return r, refuse("zero-regions", path, "no content rectangles in the five 1440 px shots")
	}
	return r, nil
}

// parseLabelled parses LABEL=PATH tokens: each label non-empty and unique, each path non-empty,
// at least one token.
func parseLabelled(args []string) ([][2]string, error) {
	var out [][2]string
	seen := map[string]bool{}
	for _, a := range args {
		label, p, ok := strings.Cut(a, "=")
		if !ok || label == "" || p == "" || seen[label] {
			return nil, fmt.Errorf("want LABEL=PATH with a unique label, got %q", a)
		}
		seen[label] = true
		out = append(out, [2]string{label, p})
	}
	if len(out) == 0 {
		return nil, errors.New("no candidates")
	}
	return out, nil
}

// printRefusal writes `marketingcheck: <code> <name>: <detail>` to stderr.
func printRefusal(stderr io.Writer, err error) {
	var r *Refusal
	if errors.As(err, &r) {
		fmt.Fprintf(stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		return
	}
	fmt.Fprintf(stderr, "marketingcheck: %v\n", err)
}

// runPosition is `marketingcheck position`: rc 0 when every label is ACCEPT, 1 when any is
// REFUSE or an input is refused, 2 on a usage error (design §2.1).
func runPosition(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("position", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, positionUsage)
		fs.PrintDefaults()
	}
	baseline := fs.String("baseline", "", "R1 baseline rectangles JSON (design_docs/verification/world-iter312/r1-rects.json)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	cands, err := parseLabelled(fs.Args())
	if *baseline == "" {
		err = errors.New("--baseline is required")
	}
	if err != nil {
		fmt.Fprintf(stderr, "marketingcheck: %v\n", err)
		fs.Usage()
		return 2
	}
	base, err := LoadRectangles(*baseline)
	if err != nil {
		printRefusal(stderr, err)
		return 1
	}
	refused := 0
	for _, c := range cands {
		cand, err := loadPositionCandidate(c[1])
		if err != nil {
			printRefusal(stderr, err)
			return 1
		}
		reasons, regions := PositionJudge(base, cand)
		if regions == 0 {
			fmt.Fprintf(stderr, "marketingcheck: zero-regions %s: compared 0 regions\n", c[0])
			return 1
		}
		v := "ACCEPT"
		if len(reasons) > 0 {
			v = "REFUSE"
			refused++
		}
		fmt.Fprintf(stdout, "%-38s %s regions=%d\n", c[0], v, regions)
		for _, r := range reasons {
			fmt.Fprintln(stdout, "    "+r)
		}
	}
	fmt.Fprintf(stdout, "TOTAL candidates=%d accept=%d refuse=%d\n", len(cands), len(cands)-refused, refused)
	if refused > 0 {
		return 1
	}
	return 0
}
